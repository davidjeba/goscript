/**
 * Goscript Client Runtime v2.0
 * ===============================
 * The bridge between Go server and browser DOM.
 * The GS compiler produces JavaScript that CALLS these functions.
 * This file is embedded in the Go binary via go:embed and served
 * automatically at /__goscript/runtime.js.
 *
 * Architecture:
 *   - State layer mirrors Go's Store on the client
 *   - DOM helpers mirror Go's CreateElement for consistency
 *   - Event bus enables server-push via GS-Trigger response headers
 *   - Reactive attributes (gs-trigger, gs-target, gs-swap) drive
 *     HTML-over-the-wire updates without writing any JavaScript
 *
 * Usage (from compiled .gs code):
 *   const [count, setCount] = __gs.useState(0);
 *   const el = __gs.h('div', { className: 'counter' },
 *     __gs.h('button', { onclick: () => setCount(n => n + 1) }, 'Click me'),
 *     __gs.h('span', null, count())
 *   );
 *   __gs.mount('#app', el);
 */
const __gs = (() => {
    'use strict';

    // =========================================================================
    //  SECTION 1 — STATE MANAGEMENT
    // =========================================================================
    //  Mirrors Go's Store on the client side. Each useState call allocates
    //  a unique numeric ID. Subscribers are notified on every change so that
    //  dependent DOM nodes can re-render.
    // =========================================================================

    /** @type {Map<number, *>} Internal state store keyed by numeric ID */
    const _state = new Map();

    /**
     * @type {Map<number, Array<function(*, *): void>>}
     * Subscribers keyed by state ID. Each callback receives (newVal, oldVal).
     */
    const _subscribers = new Map();

    /** Auto-incrementing counter for state IDs */
    let _stateId = 0;

    /**
     * useState creates a reactive state slot, mirroring Go's GlobalStore.
     * Returns [getter, setter] — the getter reads the current value and the
     * setter accepts either a direct value or an updater function(prev).
     *
     * @param {*} initial — the initial value or a thunk () => initial
     * @returns {[function(): *, function(*|function(*): *): void]}
     *
     * @example
     *   const [count, setCount] = __gs.useState(0);
     *   setCount(n => n + 1);
     *   console.log(count()); // 1
     */
    function useState(initial) {
        const id = ++_stateId;
        const initialVal = typeof initial === 'function' ? initial() : initial;
        _state.set(id, initialVal);

        // Getter — returns the current value for this state slot
        const get = () => _state.get(id);

        // Setter — accepts a value or updater function(prev) => next
        const set = (valOrFn) => {
            const old = _state.get(id);
            const newVal = typeof valOrFn === 'function' ? valOrFn(old) : valOrFn;
            if (old !== newVal) {
                _state.set(id, newVal);
                notify(id, newVal, old);
            }
        };

        return [get, set];
    }

    /**
     * useReducer wraps useState with a reducer pattern (like React useReducer).
     * The reducer receives (prevState, action) and must return the new state.
     *
     * @param {function(*, *): *} reducer — (state, action) => newState
     * @param {*} initial — the initial state value
     * @returns {[function(): *, function(*): void]} [state getter, dispatch]
     */
    function useReducer(reducer, initial) {
        const [state, setState] = useState(initial);
        return [state, (action) => setState(prev => reducer(prev, action))];
    }

    /**
     * Subscribe to state changes for a given state ID.
     * Returns an unsubscribe function.
     *
     * @param {number} stateId — the state slot to watch
     * @param {function(*, *): void} callback — (newVal, oldVal) => void
     * @returns {function(): void} unsubscribe function
     */
    function subscribe(stateId, callback) {
        if (!_subscribers.has(stateId)) _subscribers.set(stateId, []);
        _subscribers.get(stateId).push(callback);
        return () => {
            const subs = _subscribers.get(stateId);
            if (!subs) return;
            const idx = subs.indexOf(callback);
            if (idx >= 0) subs.splice(idx, 1);
        };
    }

    /**
     * Internal: notify all subscribers of a state change.
     * @param {number} stateId
     * @param {*} newVal
     * @param {*} oldVal
     */
    function notify(stateId, newVal, oldVal) {
        if (_subscribers.has(stateId)) {
            _subscribers.get(stateId).forEach(cb => cb(newVal, oldVal));
        }
    }

    /**
     * Returns a plain snapshot of all state as an object.
     * Useful for serialization or debugging.
     *
     * @returns {Object.<string, *>}
     */
    function getState() {
        return Object.fromEntries(_state);
    }

    /**
     * Hydrates the client state from server-provided data.
     * Called automatically during initialization from window.__GOSCRIPT_STATE__.
     *
     * @param {Object.<string, *>} serverState — key-value pairs from Go server
     */
    function hydrate(serverState) {
        if (serverState) {
            for (const [key, value] of Object.entries(serverState)) {
                const numKey = typeof key === 'string' ? parseInt(key, 10) : key;
                if (!isNaN(numKey)) {
                    _state.set(numKey, value);
                }
            }
        }
    }

    // =========================================================================
    //  SECTION 2 — DOM CREATION
    // =========================================================================
    //  Mirrors Go's CreateElement on the client side. Supports:
    //    - className, style (object), dangerouslySetInnerHTML
    //    - Event handlers via on* props
    //    - ref objects ({ current: Element })
    //    - Automatic kebab-case conversion for data-* and aria-* attributes
    // =========================================================================

    /**
     * insertChild appends a single child node (text, element, fragment, array)
     * to a parent element. Used internally by createElement and fragment.
     *
     * @param {Element|DocumentFragment} parent
     * @param {*} child
     */
    function insertChild(parent, child) {
        if (child == null || child === false) return;
        if (typeof child === 'string' || typeof child === 'number' || typeof child === 'boolean') {
            parent.appendChild(document.createTextNode(String(child)));
        } else if (child instanceof HTMLElement || child instanceof DocumentFragment) {
            parent.appendChild(child);
        } else if (Array.isArray(child)) {
            child.flat(Infinity).forEach(c => insertChild(parent, c));
        }
    }

    /**
     * createElement creates an HTML element from a tag string, optional props,
     * and any number of children. This is the client-side counterpart of Go's
     * goscript.CreateElement().
     *
     * Supported props:
     *   - className / class  → element.className
     *   - style (object)     → Object.assign(el.style, style)
     *   - dangerouslySetInnerHTML → el.innerHTML
     *   - on* (function)     → el.addEventListener(...)
     *   - ref ({ current })  → sets ref.current = el
     *   - All others         → el.setAttribute(kebab-case, value)
     *
     * @param {string} tag — HTML tag name (e.g. 'div', 'span')
     * @param {Object|null} [props] — attribute / event map
     * @param {...*} children — text, elements, arrays, falsy values
     * @returns {HTMLElement}
     */
    function createElement(tag, props) {
        const children = Array.prototype.slice.call(arguments, 2);
        const el = document.createElement(tag);

        // Apply props / attributes
        if (props) {
            Object.entries(props).forEach(([key, val]) => {
                // Skip null, undefined, and false
                if (val === false || val == null) return;

                if (key === 'className' || key === 'class') {
                    // CSS class list
                    el.className = String(val);
                } else if (key === 'style' && typeof val === 'object') {
                    // Inline style object
                    Object.assign(el.style, val);
                } else if (key === 'dangerouslySetInnerHTML') {
                    // Raw HTML injection
                    el.innerHTML = val.__html != null ? val.__html : String(val);
                } else if (key.startsWith('on') && typeof val === 'function') {
                    // Event handler — convert onClick → click
                    const event = key.slice(2).toLowerCase();
                    el.addEventListener(event, val);
                    // Prevent default form submission for onsubmit
                    if (event === 'submit') {
                        el.setAttribute('onsubmit', 'return false;');
                    }
                } else if (key === 'ref') {
                    // Ref object pattern: { current: Element }
                    if (typeof val === 'object' && val !== null) {
                        val.current = el;
                    }
                } else if (key === 'htmlFor') {
                    // htmlFor → for attribute
                    el.setAttribute('for', String(val));
                } else {
                    // Default: convert camelCase → kebab-case and set attribute
                    const attrName = key.replace(/([A-Z])/g, '-$1').toLowerCase();
                    el.setAttribute(attrName, String(val));
                }
            });
        }

        // Append children
        children.flat(Infinity).forEach(child => {
            if (child == null || child === false) return;
            if (typeof child === 'string' || typeof child === 'number' || typeof child === 'boolean') {
                el.appendChild(document.createTextNode(String(child)));
            } else if (child instanceof HTMLElement || child instanceof DocumentFragment) {
                el.appendChild(child);
            } else if (Array.isArray(child)) {
                child.flat(Infinity).forEach(c => insertChild(el, c));
            }
        });

        return el;
    }

    /**
     * fragment creates a DocumentFragment from any number of children.
     * Useful for returning multiple elements from a render function
     * without a wrapper div.
     *
     * @param {...*} children
     * @returns {DocumentFragment}
     */
    function fragment() {
        const children = Array.prototype.slice.call(arguments).flat(Infinity);
        const frag = document.createDocumentFragment();
        children.forEach(c => insertChild(frag, c));
        return frag;
    }

    /**
     * mount replaces or populates a target element with the given node.
     * If the target is a DocumentFragment, it clears the target and appends.
     * Otherwise, the target element is replaced in-place.
     *
     * @param {string|Element} selector — CSS selector or DOM element
     * @param {HTMLElement|DocumentFragment} element — the element to mount
     */
    function mount(selector, element) {
        const target = typeof selector === 'string'
            ? document.querySelector(selector)
            : selector;
        if (target) {
            if (element instanceof DocumentFragment) {
                target.innerHTML = '';
                target.appendChild(element);
            } else {
                target.replaceWith(element);
            }
        }
    }

    // =========================================================================
    //  SECTION 3 — COMPONENT SYSTEM
    // =========================================================================
    //  Lightweight component model. Components are simply render functions.
    //  The registry tracks them for debugging and hot-reload support.
    // =========================================================================

    /** @type {Map<number, {name: string, renderFn: function}>} */
    const _componentRenders = new Map();

    /** Auto-incrementing component ID */
    let _componentId = 0;

    /**
     * component executes a render function and returns its result.
     * This is the simplest form of a component — just a render thunk.
     *
     * @param {function(): HTMLElement|DocumentFragment|string} renderFn
     * @returns {HTMLElement|DocumentFragment|string}
     */
    function component(renderFn) {
        return renderFn();
    }

    /**
     * createComponent registers a named component and immediately renders it.
     * Used for dev-tools integration and future hot-reload support.
     *
     * @param {string} name — human-readable component name
     * @param {function(): HTMLElement|DocumentFragment|string} renderFn
     * @returns {HTMLElement|DocumentFragment|string}
     */
    function createComponent(name, renderFn) {
        const id = ++_componentId;
        _componentRenders.set(id, { name, renderFn });
        return renderFn();
    }

    // =========================================================================
    //  SECTION 4 — EFFECTS & HOOKS
    // =========================================================================
    //  Client-side hooks that mirror Go's hook functions.
    //  These are simplified for the initial runtime; full dependency tracking
    //  will be added in a future version.
    // =========================================================================

    /**
     * useEffect runs a side-effect function. For v1, the effect runs
     * immediately when called. Returns the cleanup function if provided.
     *
     * @param {function(function(): void): void|function(): void} fn
     * @param {Array<*>} [deps] — dependency array (reserved for future use)
     * @returns {function(): void|undefined} cleanup function
     */
    function useEffect(fn, deps) {
        // Execute the effect; if it returns a cleanup function, return it
        const cleanup = fn();
        return typeof cleanup === 'function' ? cleanup : undefined;
    }

    /**
     * useRef creates a mutable ref object with a .current property.
     * Equivalent to { current: initialValue }.
     *
     * @param {*} [initial] — initial value for .current
     * @returns {{ current: * }}
     */
    function useRef(initial) {
        return { current: initial };
    }

    /**
     * useMemo computes a value and caches it. For v1, the function is
     * called immediately on every invocation. Dependency-based caching
     * will be added in a future version.
     *
     * @param {function(): *} fn — memoized computation
     * @param {Array<*>} [deps] — dependency array (reserved)
     * @returns {*}
     */
    function useMemo(fn, deps) {
        return fn();
    }

    /**
     * useCallback returns the function unchanged. For v1, no memoization
     * is performed. Dependency-based caching will be added later.
     *
     * @param {function} fn
     * @param {Array<*>} [deps] — dependency array (reserved)
     * @returns {function}
     */
    function useCallback(fn, deps) {
        return fn;
    }

    // =========================================================================
    //  SECTION 5 — API HELPERS
    // =========================================================================
    //  Thin wrappers around fetch() that:
    //    1. Set the GS-Request header so the Go server identifies them
    //    2. Throw descriptive errors on non-OK responses
    //    3. Intercept GS-Trigger / GS-State response headers automatically
    // =========================================================================

    /**
     * Performs a GET request and parses the JSON response.
     * Sets the GS-Request header for server-side identification.
     *
     * @param {string} url — relative or absolute URL
     * @returns {Promise<Object>} parsed JSON
     * @throws {Error} on non-OK responses
     */
    async function getJSON(url) {
        const res = await fetch(url, {
            headers: { 'GS-Request': 'true' },
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        return res.json();
    }

    /**
     * Performs a POST request with a JSON body.
     *
     * @param {string} url
     * @param {Object} data — request body (will be JSON-stringified)
     * @returns {Promise<Object>} parsed JSON response
     */
    async function postJSON(url, data) {
        const res = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'GS-Request': 'true',
            },
            body: JSON.stringify(data),
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        return res.json();
    }

    /**
     * Performs a PUT request with a JSON body.
     *
     * @param {string} url
     * @param {Object} data
     * @returns {Promise<Object>} parsed JSON response
     */
    async function putJSON(url, data) {
        const res = await fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'GS-Request': 'true',
            },
            body: JSON.stringify(data),
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        return res.json();
    }

    /**
     * Performs a DELETE request.
     *
     * @param {string} url
     * @returns {Promise<Object>} parsed JSON response
     */
    async function deleteJSON(url) {
        const res = await fetch(url, {
            method: 'DELETE',
            headers: { 'GS-Request': 'true' },
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        return res.json();
    }

    /**
     * Performs a POST request and returns the raw HTML/text response.
     * Used when the server returns an HTML fragment for DOM swapping.
     *
     * @param {string} url
     * @param {Object} data
     * @returns {Promise<string>} HTML response text
     */
    async function postHTML(url, data) {
        const res = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'GS-Request': 'true',
            },
            body: JSON.stringify(data),
        });
        if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        return res.text();
    }

    // =========================================================================
    //  SECTION 6 — CLIENT-SIDE ROUTER
    // =========================================================================
    //  Simple hash-free SPA router. Uses the History API (pushState /
    //  popstate) so URLs remain clean. Navigate with __gs.navigate('/path').
    // =========================================================================

    /**
     * navigate pushes a new URL onto the browser history stack
     * and dispatches a popstate event so route listeners can react.
     *
     * @param {string} path — the new path (e.g. '/about')
     */
    function navigate(path) {
        window.history.pushState(null, '', path);
        window.dispatchEvent(new PopStateEvent('popstate'));
    }

    /**
     * usePathname returns the current URL pathname.
     * @returns {string}
     */
    function usePathname() {
        return window.location.pathname;
    }

    /**
     * useParams returns the URL search params as a plain object.
     * @returns {Object.<string, string>}
     */
    function useParams() {
        const params = {};
        new URLSearchParams(window.location.search).forEach((v, k) => {
            params[k] = v;
        });
        return params;
    }

    /**
     * useQuery is an alias for useParams for familiarity.
     * @returns {Object.<string, string>}
     */
    function useQuery() {
        return Object.fromEntries(new URLSearchParams(window.location.search));
    }

    /**
     * Link creates an <a> element that uses the client-side router
     * instead of causing a full page navigation.
     *
     * @param {string} href — the target path
     * @param {Object} [props] — additional anchor attributes
     * @returns {HTMLElement}
     */
    function Link(href, props) {
        return createElement('a', Object.assign({}, props, {
            href: href,
            onclick: function (e) {
                e.preventDefault();
                navigate(href);
            },
        }));
    }

    // =========================================================================
    //  SECTION 7 — EVENT BUS
    // =========================================================================
    //  A publish/subscribe event system. The Go server can trigger client-side
    //  events by setting the GS-Trigger response header. The runtime intercepts
    //  all fetch responses and fires matching events automatically.
    // =========================================================================

    /** @type {Map<string, Array<function(*): void>>} */
    const _eventBus = new Map();

    /**
     * Subscribe to a named event. Returns an unsubscribe function.
     *
     * @param {string} event — event name (e.g. 'cart:updated')
     * @param {function(*): void} handler — callback receiving event detail
     * @returns {function(): void} unsubscribe
     */
    function on(event, handler) {
        if (!_eventBus.has(event)) _eventBus.set(event, []);
        _eventBus.get(event).push(handler);
        return () => off(event, handler);
    }

    /**
     * Emit a named event to all subscribers. Errors in handlers are caught
     * and logged so one bad listener doesn't break the chain.
     *
     * @param {string} event — event name
     * @param {*} [detail] — event payload
     */
    function emit(event, detail) {
        if (_eventBus.has(event)) {
            _eventBus.get(event).forEach(h => {
                try {
                    h(detail);
                } catch (e) {
                    console.error('[goscript] event error on "' + event + '":', e);
                }
            });
        }
    }

    /**
     * Unsubscribe a handler from a named event.
     *
     * @param {string} event
     * @param {function(*): void} handler
     */
    function off(event, handler) {
        if (_eventBus.has(event)) {
            const handlers = _eventBus.get(event);
            if (handlers) {
                const idx = handlers.indexOf(handler);
                if (idx >= 0) handlers.splice(idx, 1);
            }
        }
    }

    // =========================================================================
    //  SECTION 8 — REALTIME (SSE & WebSocket)
    // =========================================================================
    //  Helpers for Server-Sent Events and WebSocket connections.
    //  SSE is auto-reconnecting; WebSocket requires manual reconnect logic.
    // =========================================================================

    /**
     * sse opens a Server-Sent Events connection and calls the handler
     * with each parsed JSON message. Returns a control object with close().
     *
     * @param {string} url — SSE endpoint URL
     * @param {function(Object): void} handler — message callback
     * @returns {{ close: function(): void }}
     */
    function sse(url, handler) {
        const source = new EventSource(url);
        source.onmessage = (e) => {
            try {
                handler(JSON.parse(e.data));
            } catch (err) {
                // If parsing fails, pass raw string
                handler(e.data);
            }
        };
        source.onerror = () => {
            // Browser auto-reconnects; we close on error to avoid spam
            source.close();
        };
        return { close: () => source.close() };
    }

    /**
     * ws opens a WebSocket connection with structured message handling.
     * Returns an object with send(data) and close() methods.
     *
     * @param {string} url — WebSocket endpoint URL
     * @param {Object} handlers — { onopen, onmessage, onclose, onerror }
     * @returns {{ send: function(*): void, close: function(): void }}
     */
    function ws(url, handlers) {
        const socket = new WebSocket(url);

        socket.onopen = handlers.onopen || (() => {});
        socket.onmessage = (e) => {
            try {
                handlers.onmessage(JSON.parse(e.data));
            } catch (err) {
                // Pass raw data if JSON parsing fails
                if (handlers.onmessage) handlers.onmessage(e.data);
            }
        };
        socket.onclose = handlers.onclose || (() => {});
        socket.onerror = handlers.onerror || (() => {});

        return {
            /**
             * Send data over the WebSocket. Objects are JSON-stringified.
             * @param {*} data
             */
            send: (data) => socket.send(typeof data === 'string' ? data : JSON.stringify(data)),

            /** Close the WebSocket connection */
            close: () => socket.close(),
        };
    }

    // =========================================================================
    //  SECTION 9 — STRING HELPERS
    // =========================================================================
    //  Lightweight Go-style string formatting, primarily sprintf.
    // =========================================================================

    /**
     * sprintf performs simple Go-style string formatting.
     * Supports: %s (string), %d (number), %v (any value).
     * Verb count must match argument count.
     *
     * @param {string} format — format string with %s, %d, %v placeholders
     * @param {...*} args — values to interpolate
     * @returns {string}
     *
     * @example
     *   __gs.sprintf('Hello, %s! You have %d messages.', 'Alice', 5);
     *   // → "Hello, Alice! You have 5 messages."
     */
    function sprintf(format) {
        const args = Array.prototype.slice.call(arguments, 1);
        let i = 0;
        return format.replace(/%[sdv%]/g, (match) => {
            if (match === '%%') return '%';
            const val = args[i++];
            return val != null ? String(val) : 'null';
        });
    }

    // =========================================================================
    //  SECTION 10 — REACTIVE ATTRIBUTE ENGINE
    // =========================================================================
    //  Processes gs-trigger / gs-target / gs-swap attributes on DOM elements.
    //  When an element with gs-trigger is activated (e.g. clicked), the runtime
    //  sends a request to the specified URL and swaps the response HTML into
    //  the target element. This is the HTML-over-the-wire pattern.
    // =========================================================================

    /**
     * _swapHTML performs a DOM swap operation based on the strategy string.
     *
     * @param {Element} target — the element to swap into / around
     * @param {string} content — raw HTML string
     * @param {string} strategy — one of: innerHTML, outerHTML, beforeend,
     *                           afterend, beforebegin, afterbegin, delete, morph
     */
    function _swapHTML(target, content, strategy) {
        if (!target) return;

        switch (strategy) {
            case 'innerHTML':
                target.innerHTML = content;
                break;

            case 'outerHTML': {
                const temp = document.createElement('div');
                temp.innerHTML = content;
                const parent = target.parentNode;
                if (parent) {
                    while (temp.firstChild) {
                        parent.insertBefore(temp.firstChild, target);
                    }
                    parent.removeChild(target);
                }
                break;
            }

            case 'beforeend':
                target.insertAdjacentHTML('beforeend', content);
                break;

            case 'afterend':
                target.insertAdjacentHTML('afterend', content);
                break;

            case 'beforebegin':
                target.insertAdjacentHTML('beforebegin', content);
                break;

            case 'afterbegin':
                target.insertAdjacentHTML('afterbegin', content);
                break;

            case 'delete':
                target.remove();
                break;

            case 'morph':
                // Simple morph: diff innerHTML and update changed nodes
                if (target.innerHTML !== content) {
                    target.innerHTML = content;
                }
                break;

            case 'none':
            default:
                // Do nothing
                break;
        }
    }

    /**
     * _getTriggerURL extracts the URL from a gs-trigger value.
     * Supports format: "event url" (e.g. "click /api/data")
     *
     * @param {string} triggerValue — the gs-trigger attribute value
     * @returns {{ event: string, url: string, modifier: string|null }}
     */
    function _parseTrigger(triggerValue) {
        const parts = triggerValue.trim().split(/\s+/);
        const event = parts[0] || 'click';
        const url = parts[1] || '';
        // Third part can be a modifier like "changed" for filtered triggers
        const modifier = parts[2] || null;
        return { event, url, modifier };
    }

    /**
     * _initReactiveEngine scans the DOM for elements with gs-trigger
     * and attaches event listeners for the HTML-over-the-wire pattern.
     * Called once during runtime initialization.
     */
    function _initReactiveEngine() {
        document.querySelectorAll('[gs-trigger]').forEach(el => {
            const triggerAttr = el.getAttribute('gs-trigger');
            if (!triggerAttr) return;

            const parsed = _parseTrigger(triggerAttr);
            const eventNames = parsed.event.split(',');

            eventNames.forEach(eventName => {
                eventName = eventName.trim();

                // Handle special triggers
                if (eventName === 'load') {
                    // Execute immediately on initialization
                    _executeReactiveRequest(el, parsed);
                    return;
                }

                if (eventName === 'every') {
                    // Polling: gs-trigger="every 2s /api/data"
                    const intervalMs = _parseDuration(parsed.url);
                    if (intervalMs > 0) {
                        setInterval(() => _executeReactiveRequest(el, {
                            event: 'poll',
                            url: parsed.modifier || '',
                        }), intervalMs);
                    }
                    return;
                }

                if (eventName === 'intersect') {
                    // Intersection Observer: load when visible
                    const observer = new IntersectionObserver((entries) => {
                        entries.forEach(entry => {
                            if (entry.isIntersecting) {
                                _executeReactiveRequest(el, {
                                    event: 'intersect',
                                    url: parsed.url,
                                });
                                observer.unobserve(el);
                            }
                        });
                    });
                    observer.observe(el);
                    return;
                }

                // Standard DOM events (click, submit, change, etc.)
                el.addEventListener(eventName, (e) => {
                    // Check gs-confirm attribute
                    const confirmMsg = el.getAttribute('gs-confirm');
                    if (confirmMsg && !window.confirm(confirmMsg)) {
                        e.preventDefault();
                        return;
                    }

                    // Handle form submission
                    if (eventName === 'submit') {
                        e.preventDefault();
                    }

                    _executeReactiveRequest(el, parsed);
                });
            });
        });

        // Handle gs-boost on <a> and <form> elements
        document.querySelectorAll('[gs-boost]').forEach(el => {
            const tagName = el.tagName.toLowerCase();

            if (tagName === 'a') {
                el.addEventListener('click', (e) => {
                    e.preventDefault();
                    const href = el.getAttribute('href');
                    if (href) {
                        _executeReactiveRequest(el, {
                            event: 'click',
                            url: href,
                        });
                    }
                });
            } else if (tagName === 'form') {
                el.addEventListener('submit', (e) => {
                    e.preventDefault();
                    const action = el.getAttribute('action');
                    if (action) {
                        const formData = new FormData(el);
                        const method = (el.getAttribute('method') || 'GET').toUpperCase();

                        const target = el.getAttribute('gs-target');
                        const swap = el.getAttribute('gs-swap') || 'innerHTML';

                        fetch(action, {
                            method: method,
                            body: formData,
                            headers: { 'GS-Request': 'true' },
                        })
                            .then(res => _processResponseHeaders(res))
                            .then(res => res.text())
                            .then(html => {
                                if (target) {
                                    const targetEl = document.querySelector(target);
                                    if (targetEl) _swapHTML(targetEl, html, swap);
                                }
                                _reprocessReactiveElements();
                            })
                            .catch(err => {
                                console.error('[goscript] boost request error:', err);
                            });
                    }
                });
            }
        });
    }

    /**
     * _executeReactiveRequest sends a request for a reactive element
     * and performs the DOM swap with the response.
     *
     * @param {Element} el — the element with gs-trigger
     * @param {{ event: string, url: string, modifier: string|null }} parsed
     */
    function _executeReactiveRequest(el, parsed) {
        const target = el.getAttribute('gs-target') || el.getAttribute('gs-target');
        const swap = el.getAttribute('gs-swap') || 'innerHTML';
        const indicator = el.getAttribute('gs-indicator');
        const disabledEls = el.getAttribute('gs-disabled');
        const pushUrl = el.getAttribute('gs-push-url');

        // Show loading indicator if specified
        let indicatorEl = null;
        if (indicator) {
            indicatorEl = document.querySelector(indicator);
            if (indicatorEl) indicatorEl.style.display = '';
        }

        // Disable specified elements during request
        let disabledElements = [];
        if (disabledEls) {
            disabledEls.split(',').forEach(sel => {
                const target = document.querySelector(sel.trim());
                if (target) {
                    target.disabled = true;
                    disabledElements.push(target);
                }
            });
        }

        fetch(parsed.url, {
            headers: { 'GS-Request': 'true' },
        })
            .then(res => _processResponseHeaders(res))
            .then(res => res.text())
            .then(html => {
                // Perform the DOM swap
                if (target) {
                    const targetEl = document.querySelector(target);
                    if (targetEl) _swapHTML(targetEl, html, swap);
                }

                // Push URL if specified
                if (pushUrl) {
                    window.history.pushState(null, '', pushUrl);
                }

                // Re-process any new reactive elements in the swapped content
                _reprocessReactiveElements();
            })
            .catch(err => {
                console.error('[goscript] reactive request error:', err);
            })
            .finally(() => {
                // Hide loading indicator
                if (indicatorEl) indicatorEl.style.display = 'none';

                // Re-enable disabled elements
                disabledElements.forEach(target => {
                    target.disabled = false;
                });
            });
    }

    /**
     * _processResponseHeaders reads GS-Trigger and GS-State headers
     * from a fetch response and dispatches events / hydrates state.
     *
     * @param {Response} res — fetch Response object
     * @returns {Response} the original response (for chaining)
     */
    function _processResponseHeaders(res) {
        // Handle GS-Trigger: fire client-side events from server
        const trigger = res.headers.get('GS-Trigger');
        if (trigger) {
            try {
                const events = JSON.parse(trigger);
                if (Array.isArray(events)) {
                    events.forEach(e => emit(e.name, e.detail));
                } else {
                    emit(events.name, events.detail);
                }
            } catch (err) {
                // Invalid JSON — ignore
            }
        }

        // Handle GS-State: sync server state to client
        const stateHeader = res.headers.get('GS-State');
        if (stateHeader) {
            try {
                hydrate(JSON.parse(stateHeader));
            } catch (err) {
                // Invalid JSON — ignore
            }
        }

        return res;
    }

    /**
     * _reprocessReactiveElements re-scans the DOM for new elements
     * with gs-trigger that were added via swaps.
     */
    function _reprocessReactiveElements() {
        document.querySelectorAll('[gs-trigger]:not([gs-initialized])').forEach(el => {
            el.setAttribute('gs-initialized', '');
            const triggerAttr = el.getAttribute('gs-trigger');
            if (!triggerAttr) return;

            const parsed = _parseTrigger(triggerAttr);
            const eventNames = parsed.event.split(',');

            eventNames.forEach(eventName => {
                eventName = eventName.trim();

                if (eventName === 'load') {
                    _executeReactiveRequest(el, parsed);
                    return;
                }

                el.addEventListener(eventName, (e) => {
                    if (eventName === 'submit') e.preventDefault();
                    _executeReactiveRequest(el, parsed);
                });
            });
        });
    }

    /**
     * _parseDuration converts a duration string like "2s", "500ms", "1m"
     * into milliseconds.
     *
     * @param {string} duration — duration string
     * @returns {number} milliseconds, or 0 if unparseable
     */
    function _parseDuration(duration) {
        const match = duration.match(/^(\d+)\s*(ms|s|m)$/);
        if (!match) return 0;
        const value = parseInt(match[1], 10);
        switch (match[2]) {
            case 'ms': return value;
            case 's': return value * 1000;
            case 'm': return value * 60 * 1000;
            default: return 0;
        }
    }

    // =========================================================================
    //  SECTION 11 — INITIALIZATION
    // =========================================================================
    //  Bootstraps the runtime: hydrates state from the server, patches
    //  fetch to intercept response headers, and starts the reactive engine.
    // =========================================================================

    // Hydrate from server state embedded in the HTML
    hydrate(window.__GOSCRIPT_STATE__);

    // Patch window.fetch to intercept GS-* response headers on ALL requests
    const _origFetch = window.fetch;
    window.fetch = function () {
        return _origFetch.apply(this, arguments).then(res => {
            return _processResponseHeaders(res);
        });
    };

    // Start the reactive attribute engine once the DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', _initReactiveEngine);
    } else {
        _initReactiveEngine();
    }

    // =========================================================================
    //  SECTION 12 — PUBLIC API
    // =========================================================================

    return {
        // --- State Management ---
        useState:    useState,
        useReducer:  useReducer,
        subscribe:   subscribe,
        getState:    getState,
        hydrate:     hydrate,

        // --- DOM Creation ---
        createElement:   createElement,   // full name
        el:              createElement,   // short alias
        h:               createElement,   // JSX-like alias
        mount:           mount,
        fragment:        fragment,
        insertChild:     insertChild,
        createComponent: createComponent,
        component:       component,

        // --- Effects & Hooks ---
        useEffect:  useEffect,
        useRef:     useRef,
        useMemo:    useMemo,
        useCallback: useCallback,

        // --- API Helpers ---
        getJSON:    getJSON,
        postJSON:   postJSON,
        putJSON:    putJSON,
        deleteJSON: deleteJSON,
        postHTML:   postHTML,

        // --- Router ---
        navigate:    navigate,
        usePathname: usePathname,
        useParams:   useParams,
        useQuery:    useQuery,
        Link:        Link,

        // --- Event Bus ---
        on:   on,
        off:  off,
        emit: emit,

        // --- Realtime ---
        sse: sse,
        ws:  ws,

        // --- String Helpers ---
        sprintf: sprintf,

        // --- Internal (exposed for advanced use) ---
        /** @type {string} Runtime version */
        version: '2.0.0',
    };
})();

// Expose globally so compiled .gs code can access the runtime
window.__gs = __gs;
