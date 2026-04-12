// Package goscript — types.go provides version information, framework metadata,
// and serves as the central registration point for key GoScript v2 types.
//
// This file does NOT define new types that belong in their respective feature files.
// Instead, it re-exports version constants and framework identification used
// throughout the GoScript ecosystem.
package goscript

// Version is the semantic version of the GoScript framework.
const Version = "2.0.0"

// FrameworkName is the canonical name of the GoScript framework.
const FrameworkName = "GoScript"

// GoVersion is the minimum required Go version for this framework release.
const GoVersion = "go1.22"

// FeatureFlags tracks which v2 features are available at runtime.
var FeatureFlags = map[string]bool{
	"app-router":           true,
	"streaming-ssr":        true,
	"server-components":    true,
	"api-routes":           true,
	"middleware-pipeline":   true,
	"ssg-isr":              true,
	"error-loading":        true,
	"metadata-seo":         true,
	"hmr":                  true,
	"convention-routes":    true,
}

// V2Features lists all v2 feature names in their canonical order.
var V2Features = []string{
	"App Router (file-system conventions)",
	"Streaming SSR with Suspense",
	"Server & Client Components",
	"Convention-based API Routes",
	"Composable Middleware Pipeline",
	"Static Site Generation & ISR",
	"Error & Loading Boundaries",
	"Metadata & SEO API",
	"Hot Module Replacement Dev Server",
	"Route Groups & Dynamic Segments",
}

// V1SubModules lists the original v1 sub-modules that ship with GoScript.
var V1SubModules = []string{
	"Gocsx — CSS-in-Go Framework",
	"GoScale API — GraphQL/gRPC Hybrid API System",
	"GoScale DB — PostgreSQL-compatible Database",
	"GoScale Edge — Edge Computing Network",
	"GOPM — Go Package Manager",
	"GoUIX — Interactive Canvas UI Framework",
	"Jetpack Core — Performance Monitoring",
	"Jetpack Frontend — Lighthouse & Performance Panel",
	"Jetpack Security — Vulnerability Scanning",
}

// IsFeatureEnabled checks whether a specific v2 feature flag is enabled.
func IsFeatureEnabled(feature string) bool {
	return FeatureFlags[feature]
}

// GetFrameworkInfo returns a map of framework metadata useful for diagnostics
// and about pages.
func GetFrameworkInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":          FrameworkName,
		"version":       Version,
		"go_version":    GoVersion,
		"features":      V2Features,
		"sub_modules":   V1SubModules,
		"feature_flags": FeatureFlags,
	}
}
