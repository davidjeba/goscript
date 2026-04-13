package gscompiler

import "github.com/davidjeba/goscript/pkg/gslib"

// RuntimeJS contains the embedded goscript client runtime.
// This is re-exported from pkg/gslib for convenient access within
// the gscompiler package. The actual embedding happens in
// pkg/gslib/runtime_embed.go via go:embed.
//
// Use WithRuntime() bundle option to prepend it to the output.
var RuntimeJS = gslib.RuntimeJS
