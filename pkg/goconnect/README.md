# GoConnect: Backend Integration System for GoScript

GoConnect provides a high-performance, secure integration layer between GoScript frontend applications and various backend frameworks. It offers optimized connectors for popular backend systems like Django, WordPress, and Beego, with a focus on security and performance.

## Core Features

### Plugin Framework

GoConnect's extensible plugin system allows for seamless integration with any backend framework.

```go
// Register a backend plugin
goconnect.RegisterBackend("django", &goconnect.BackendPlugin{
    Name: "Django",
    Version: "1.0.0",
    Initialize: func(config map[string]interface{}) error {
        // Initialize Django-specific configuration
        return nil
    },
    Adapters: map[string]goconnect.Adapter{
        "rest": djangoRestAdapter,
        "auth": djangoAuthAdapter,
        "forms": djangoFormsAdapter,
    },
})

// Use the plugin
client := goconnect.New("django", map[string]interface{}{
    "baseURL": "https://api.example.com",
    "csrfTokenName": "csrftoken",
})
```

### Framework-specific Adapters

Optimized connectors for popular backend frameworks with framework-specific features.

#### Django Integration

```go
// Django REST API integration
response, err := client.API().Get("/api/users/", map[string]interface{}{
    "active": true,
    "role": "admin",
})

// Django form submission with CSRF protection
success, errors := client.Forms().Submit("contact", map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "message": "Hello, world!",
})

// Django authentication
user, err := client.Auth().Login("username", "password")
```

#### WordPress Integration

```go
// WordPress REST API integration
posts, err := client.API().Get("/wp-json/wp/v2/posts", map[string]interface{}{
    "per_page": 10,
    "categories": 5,
})

// WordPress form submission
success, errors := client.Forms().Submit("comment", map[string]interface{}{
    "post_id": 123,
    "author": "John Doe",
    "email": "john@example.com",
    "content": "Great article!",
})

// WordPress authentication
user, err := client.Auth().Login("username", "password")
```

#### Beego Integration

```go
// Beego API integration
response, err := client.API().Get("/api/v1/users", map[string]interface{}{
    "active": true,
})

// Beego form submission
success, errors := client.Forms().Submit("registration", map[string]interface{}{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword",
})

// Beego authentication
user, err := client.Auth().Login("username", "password")
```

### Authentication System

Comprehensive authentication system with support for various authentication methods and granular RBAC.

```go
// Configure authentication
auth := goconnect.NewAuth(goconnect.AuthOptions{
    Providers: []goconnect.AuthProvider{
        goconnect.JWTProvider{
            TokenStorage: goconnect.LocalStorage,
            TokenKey: "auth_token",
            RefreshTokenKey: "refresh_token",
            AutoRefresh: true,
        },
        goconnect.OAuth2Provider{
            ClientID: "client_id",
            AuthorizationEndpoint: "https://auth.example.com/oauth/authorize",
            TokenEndpoint: "https://auth.example.com/oauth/token",
            RedirectURI: "https://myapp.com/callback",
            Scope: "profile email",
        },
        goconnect.CookieProvider{
            CookieName: "session_id",
            Secure: true,
            HttpOnly: true,
        },
    },
})

// Login with different providers
user, err := auth.Login("jwt", map[string]interface{}{
    "username": "johndoe",
    "password": "password123",
})

// Check authentication status
if auth.IsAuthenticated() {
    // User is authenticated
    currentUser := auth.GetCurrentUser()
}

// Logout
auth.Logout()
```

### Access Control System

Granular RBAC with dynamic access control capabilities.

```go
// Configure access control
rbac := goconnect.NewRBAC(goconnect.RBACOptions{
    RolesEndpoint: "/api/roles",
    PermissionsEndpoint: "/api/permissions",
    CacheExpiration: 5 * time.Minute,
})

// Check permissions
if rbac.Can("edit", "document", documentId) {
    // User can edit this document
}

// Cell-level access control
if rbac.CanAccessCell("users", userId, "salary") {
    // User can access the salary cell for this user
}

// Time-based access
if rbac.CanWithTime("view", "financial_report", reportId, time.Now()) {
    // User can view this financial report at the current time
}

// Request/approval workflow
requestId := rbac.RequestAccess("delete", "user", userId, map[string]interface{}{
    "reason": "Account cleanup",
    "notify": "admin@example.com",
})

// Check request status
status := rbac.GetRequestStatus(requestId)
```

### Security Layer

Built-in protection against common vulnerabilities.

```go
// Configure security
security := goconnect.NewSecurity(goconnect.SecurityOptions{
    CSRF: goconnect.CSRFOptions{
        Enabled: true,
        TokenName: "csrf_token",
        HeaderName: "X-CSRF-Token",
    },
    XSS: goconnect.XSSOptions{
        Enabled: true,
        SanitizeResponses: true,
    },
    RateLimiting: goconnect.RateLimitOptions{
        Enabled: true,
        MaxRequests: 100,
        Window: 60 * time.Second,
    },
})

// Security is automatically applied to all requests
```

## Advanced Features

### Data Fetching

Optimized data fetching with caching, batching, and prefetching.

```go
// Configure data fetching
fetcher := goconnect.NewFetcher(goconnect.FetcherOptions{
    Cache: goconnect.CacheOptions{
        Enabled: true,
        TTL: 5 * time.Minute,
        MaxSize: 100,
    },
    Batch: goconnect.BatchOptions{
        Enabled: true,
        MaxBatchSize: 10,
        BatchWindow: 50 * time.Millisecond,
    },
    Prefetch: goconnect.PrefetchOptions{
        Enabled: true,
        PrefetchOnHover: true,
    },
})

// Simple fetch
user, err := fetcher.Fetch("/api/users/123")

// Fetch with query parameters
posts, err := fetcher.Fetch("/api/posts", map[string]interface{}{
    "author": "johndoe",
    "published": true,
    "limit": 10,
})

// Fetch with cache control
data, err := fetcher.Fetch("/api/data", nil, goconnect.FetchOptions{
    CacheControl: "no-cache",
})

// Prefetch data
fetcher.Prefetch("/api/users/123")
```

### Real-time Updates

Support for WebSockets and Server-Sent Events for real-time updates.

```go
// Configure real-time updates
realtime := goconnect.NewRealtime(goconnect.RealtimeOptions{
    Type: goconnect.WebSocket,
    URL: "wss://api.example.com/ws",
    Reconnect: true,
    MaxReconnectAttempts: 5,
    ReconnectInterval: 2 * time.Second,
})

// Subscribe to a channel
subscription := realtime.Subscribe("users", func(data interface{}) {
    // Handle real-time updates
    fmt.Printf("Received update: %v\n", data)
})

// Send a message
realtime.Send("chat", map[string]interface{}{
    "message": "Hello, world!",
    "sender": "johndoe",
})

// Unsubscribe
subscription.Unsubscribe()

// Close connection
realtime.Close()
```

### Form Handling

Advanced form handling with backend integration.

```go
// Create a form with backend integration
form := goconnect.NewForm("registration", goconnect.FormOptions{
    Endpoint: "/api/register",
    Method: "POST",
    ValidateOnChange: true,
    ValidateOnBlur: true,
    ValidateOnSubmit: true,
})

// Initialize form with backend schema
form.InitFromSchema("/api/forms/registration/schema")

// Handle form submission
form.OnSubmit(func(values map[string]interface{}, isValid bool) {
    if isValid {
        // Form is valid, submit to backend
        form.Submit(values, func(success bool, response interface{}) {
            if success {
                // Handle successful submission
            } else {
                // Handle submission error
            }
        })
    }
})

// Connect form to UI
gosky.Component("RegistrationForm", func(props map[string]interface{}) string {
    return `
        <form id="registration" onsubmit="{{.form.HandleSubmit}}">
            <!-- Form fields -->
        </form>
    `
}, map[string]interface{}{
    "form": form,
})
```

### File Upload

Optimized file upload with progress tracking and chunked uploads.

```go
// Configure file upload
uploader := goconnect.NewUploader(goconnect.UploaderOptions{
    Endpoint: "/api/upload",
    ChunkSize: 1024 * 1024, // 1MB chunks
    Concurrency: 3,
    AllowedTypes: []string{"image/*", "application/pdf"},
    MaxFileSize: 10 * 1024 * 1024, // 10MB
})

// Upload a file
upload := uploader.Upload(file, map[string]interface{}{
    "folder": "documents",
    "public": true,
})

// Track progress
upload.OnProgress(func(progress float64) {
    fmt.Printf("Upload progress: %.2f%%\n", progress * 100)
})

// Handle completion
upload.OnComplete(func(result interface{}) {
    fmt.Printf("Upload complete: %v\n", result)
})

// Handle error
upload.OnError(func(err error) {
    fmt.Printf("Upload error: %v\n", err)
})

// Cancel upload
upload.Cancel()
```

## Integration with GoScript Ecosystem

GoConnect is designed to work seamlessly with other GoScript components:

- **GoSky**: For data fetching during SSR and hydration
- **GoStore**: For managing API request state
- **Gocsx**: For styling based on API data
- **Jetpack**: For monitoring API performance

## Getting Started

```bash
# Install GoConnect using GOPM
gopm get github.com/davidjeba/goscript/pkg/goconnect

# Create a new project with GoConnect
gopm init myproject --template goconnect
```

## Performance Considerations

GoConnect is optimized for frontend-backend communication performance:

- **Connection Pooling**: Reuse connections for better performance
- **Request Batching**: Combine multiple requests into one
- **Response Caching**: Cache responses to reduce server load
- **Compression**: Automatically compress request and response data
- **Optimistic Updates**: Update UI before server confirmation
- **Prefetching**: Load data before it's needed
- **Incremental Loading**: Load data incrementally for better UX