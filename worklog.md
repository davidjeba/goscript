---
Task ID: 2
Agent: fullstack-builder
Task: Build complete goscript v2 project with all v1 files + 10 v2 improvements

Work Log:
- Copied all 48 v1 source files from /home/z/my-project/v1-sources/ to /home/z/my-project/goscript/
- Read and analyzed all v1 types: Component, BaseComponent, Props, Store, Router, RouteHandler, SSREngine, Context, AssetManager, LifecycleComponent, ComponentRegistry
- Updated go.mod from Go 1.17 to Go 1.22
- Created 10 new v2 improvement files in pkg/goscript/:
  1. app-router.go — AppRouter with file-system conventions, route groups, dynamic/catch-all segments, nested layouts
  2. streaming-ssr.go — StreamSSREngine with SuspenseBoundary, chunked HTML streaming, skeleton fallbacks
  3. server-client-components.go — ServerComponent (zero JS) and ClientComponent (hydration markers, events)
  4. api-routes.go — APIRouter with rate limiting (sliding window), auth, CORS, automatic JSON responses
  5. middleware-pipeline.go — Pipeline with Gzip, RequestID, CORS, Security, Recovery, Logging, RateLimit, Session middleware
  6. ssg-isr.go — SSGEngine with RenderSSG, RenderSSR, RenderISR modes and background revalidation
  7. error-loading-boundaries.go — ErrorBoundary (panic recovery), LoadingBoundary (async loading), SkeletonComponent
  8. metadata-seo.go — MetadataBuilder with OpenGraph, Twitter Cards, JSON-LD, robots directives
  9. hmr.go — DevServer with WebSocket HMR, file watching, gorilla/websocket integration
  10. types.go — Version constants, FeatureFlags, framework metadata
- Rewrote cmd/server/main.go to showcase all v2 features while maintaining v1 compatibility
- Wrote comprehensive README.md with feature comparison table, quick start, project structure, code examples, and v1 sub-module documentation

Stage Summary:
- Complete goscript v2 project at /home/z/my-project/goscript/
- 59 total files (48 v1 + 10 v2 + 1 updated main.go)
- All v2 files integrate with existing v1 types (Props, Component, Store, Router, RouteHandler, SSREngine, Context, AssetManager, LifecycleComponent, ComponentRegistry)
- go.mod updated to Go 1.22
- README.md with comprehensive documentation
