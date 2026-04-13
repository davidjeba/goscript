package goscript

import (
        "context"
        "fmt"
        "strings"
        "sync"
)

// ErrorBoundary catches panics and errors that occur during child component
// rendering and falls back to a user-provided error component instead of
// crashing the entire page.
type ErrorBoundary struct {
        *BaseComponent
        fallback Component
        children []Component
        errored  bool
        errMsg   string
        mutex    sync.RWMutex
}

// NewErrorBoundary creates a new ErrorBoundary with the given fallback component
// and child components. If any child component panics or returns an error during
// rendering, the fallback is displayed instead.
func NewErrorBoundary(fallback Component, children ...Component) *ErrorBoundary {
        base := NewBaseComponent(nil, nil)
        return &ErrorBoundary{
                BaseComponent: base,
                fallback:      fallback,
                children:      children,
                errored:       false,
        }
}

// Render safely renders all child components. If any child panics, the panic
// is recovered and the fallback component is rendered instead. The error message
// is captured and can be inspected via the Error() method.
func (eb *ErrorBoundary) Render() string {
        eb.mutex.Lock()
        defer eb.mutex.Unlock()

        if eb.errored {
                if eb.fallback != nil {
                        props := Props{
                                "errorMessage": eb.errMsg,
                        }
                        _ = NewBaseComponent(props, nil)
                        return eb.fallback.Render()
                }
                return fmt.Sprintf(`<div class="error-boundary" style="padding:1rem;border:2px solid #e53e3e;border-radius:8px;background:#fff5f5;color:#e53e3e;">
<strong>Error:</strong> %s
</div>`, eb.errMsg)
        }

        var html strings.Builder
        for _, child := range eb.children {
                func() {
                        defer func() {
                                if r := recover(); r != nil {
                                        eb.errored = true
                                        eb.errMsg = fmt.Sprintf("%v", r)
                                }
                        }()
                        html.WriteString(child.Render())
                }()
        }

        if eb.errored {
                if eb.fallback != nil {
                        return eb.fallback.Render()
                }
                return fmt.Sprintf(`<div class="error-boundary" style="padding:1rem;border:2px solid #e53e3e;border-radius:8px;background:#fff5f5;color:#e53e3e;">
<strong>Error:</strong> %s
</div>`, eb.errMsg)
        }

        return html.String()
}

// Error returns the error message captured by the boundary, or an empty string
// if no error has occurred.
func (eb *ErrorBoundary) Error() string {
        eb.mutex.RLock()
        defer eb.mutex.RUnlock()
        return eb.errMsg
}

// Reset clears the error state, allowing the children to attempt rendering again.
func (eb *ErrorBoundary) Reset() {
        eb.mutex.Lock()
        defer eb.mutex.Unlock()
        eb.errored = false
        eb.errMsg = ""
}

// HasError returns true if the boundary has caught an error.
func (eb *ErrorBoundary) HasError() bool {
        eb.mutex.RLock()
        defer eb.mutex.RUnlock()
        return eb.errored
}

// LoadingBoundary wraps an asynchronous component loader with a skeleton component.
// While the async content is loading, the skeleton is displayed. Once the loader
// completes, the resolved component is rendered.
type LoadingBoundary struct {
        *BaseComponent
        skeleton Component
        loader   func(context.Context) (Component, error)
        loaded   Component
        loadErr  error
        ready    bool
        mutex    sync.RWMutex
}

// NewLoadingBoundary creates a new LoadingBoundary with the given skeleton
// component and asynchronous loader function. The loader receives a context
// and returns a Component or an error.
func NewLoadingBoundary(skeleton Component, loader func(context.Context) (Component, error)) *LoadingBoundary {
        base := NewBaseComponent(nil, nil)
        return &LoadingBoundary{
                BaseComponent: base,
                skeleton:      skeleton,
                loader:       loader,
                ready:         false,
        }
}

// Render returns either the skeleton placeholder (while loading) or the resolved
// component (once loading is complete). Call Load to initiate the asynchronous
// loading process before rendering.
func (lb *LoadingBoundary) Render() string {
        return lb.RenderCtx(context.Background())
}

// RenderCtx renders the loading boundary with an explicit context. If the
// component has not been loaded yet, the skeleton is displayed. The caller
// should use the Load method to trigger asynchronous loading.
func (lb *LoadingBoundary) RenderCtx(ctx context.Context) string {
        lb.mutex.RLock()
        defer lb.mutex.RUnlock()

        if !lb.ready {
                if lb.skeleton != nil {
                        return lb.skeleton.Render()
                }
                return renderDefaultSkeleton(3)
        }

        if lb.loadErr != nil {
                return fmt.Sprintf(`<div class="loading-error" style="color:#e53e3e;padding:1rem;">Failed to load content: %s</div>`, lb.loadErr.Error())
        }

        if lb.loaded != nil {
                return lb.loaded.Render()
        }

        return renderDefaultSkeleton(3)
}

// Load initiates the asynchronous loading of the component in a background
// goroutine. The result is stored internally and becomes available on the
// next Render call.
func (lb *LoadingBoundary) Load(ctx context.Context) {
        go func() {
                component, err := lb.loader(ctx)

                lb.mutex.Lock()
                defer lb.mutex.Unlock()
                lb.ready = true
                lb.loadErr = err
                lb.loaded = component
        }()
}

// LoadSync loads the component synchronously, blocking until the loader completes.
func (lb *LoadingBoundary) LoadSync(ctx context.Context) {
        component, err := lb.loader(ctx)

        lb.mutex.Lock()
        defer lb.mutex.Unlock()
        lb.ready = true
        lb.loadErr = err
        lb.loaded = component
}

// IsReady returns true if the component has been loaded.
func (lb *LoadingBoundary) IsReady() bool {
        lb.mutex.RLock()
        defer lb.mutex.RUnlock()
        return lb.ready
}

// SkeletonComponent generates animated skeleton placeholder HTML for use during
// loading states. The skeleton consists of pulsing gray bars that simulate
// content loading.
type SkeletonComponent struct {
        *BaseComponent
        lines int
        width string
}

// NewSkeletonLoader creates a new SkeletonComponent with the specified number of
// placeholder lines. Each line has a slightly randomized width to create a more
// natural loading appearance.
func NewSkeletonLoader(lines int) *SkeletonComponent {
        if lines <= 0 {
                lines = 3
        }
        base := NewBaseComponent(nil, nil)
        return &SkeletonComponent{
                BaseComponent: base,
                lines:         lines,
                width:         "100%",
        }
}

// Render generates the skeleton HTML with animated placeholder bars.
func (s *SkeletonComponent) Render() string {
        return renderDefaultSkeleton(s.lines)
}

// renderDefaultSkeleton generates a skeleton loading placeholder with the given
// number of lines. Each line uses a CSS shimmer animation.
func renderDefaultSkeleton(lines int) string {
        var sb strings.Builder
        sb.WriteString(`<div class="goscript-skeleton-loader" style="padding:1rem;">`)

        widths := []string{"100%", "85%", "92%", "78%", "95%", "88%", "70%"}
        for i := 0; i < lines; i++ {
                width := "100%"
                if i < len(widths) {
                        width = widths[i]
                }
                height := "14px"
                if i == lines-1 {
                        height = "12px"
                        width = "60%"
                }
                sb.WriteString(fmt.Sprintf(
                        `<div class="goscript-skeleton" style="width:%s;height:%s;margin-bottom:12px;border-radius:4px;"></div>`,
                        width, height,
                ))
        }

        sb.WriteString(`</div>`)
        sb.WriteString(`<style>
.goscript-skeleton {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: goscript-shimmer 1.5s ease-in-out infinite;
}
@keyframes goscript-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>`)

        return sb.String()
}

// NewErrorComponent creates a simple error display Component from an error message.
// This is a convenience function for creating fallback components for ErrorBoundary.
func NewErrorComponent(message string) Component {
        base := NewBaseComponent(nil, nil)
        return &errorDisplayComponent{
                BaseComponent: base,
                message:       message,
        }
}

// errorDisplayComponent is an internal component that renders an error message.
type errorDisplayComponent struct {
        *BaseComponent
        message string
}

// Render renders the error display HTML.
func (e *errorDisplayComponent) Render() string {
        if e.message == "" {
                e.message = "An unexpected error occurred"
        }
        return fmt.Sprintf(`<div class="error-display" style="padding:1.5rem;border:2px solid #e53e3e;border-radius:8px;background:#fff5f5;color:#c53030;font-family:system-ui,sans-serif;">
<div style="font-size:1.25rem;font-weight:600;margin-bottom:0.5rem;">Something went wrong</div>
<div style="color:#718096;">%s</div>
</div>`, e.message)
}
