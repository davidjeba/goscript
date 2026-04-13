package gslib

import _ "embed"

// RuntimeJS contains the embedded goscript client runtime.
// This is the full __gs runtime (state management, DOM creation,
// component system, effects/hooks, API helpers, router, event bus,
// realtime, string helpers, reactive attribute engine).
//
// The gscompiler bundler uses WithRuntime() to prepend this to the
// output bundle.
//
//go:embed runtime.js
var RuntimeJS string
