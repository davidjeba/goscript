# Jetpack - Performance Monitoring and Optimization

Jetpack is a comprehensive performance monitoring and optimization toolkit for web applications, providing real-time metrics, Google Lighthouse integration, and security monitoring.

## Features

- **Real-time Performance Monitoring**: Track FPS, memory usage, API latency, and more
- **Google Lighthouse Integration**: Run Lighthouse audits and track Core Web Vitals
- **Performance Panel**: Floating translucent panel for real-time metrics visualization
- **Chrome DevTools Extension**: Advanced performance monitoring in Chrome DevTools
- **Security Monitoring**: Track vulnerabilities, suspicious activities, and security compliance
- **Full-Stack Monitoring**: Frontend, backend, middleware, and database performance metrics
- **Exportable Reports**: Generate comprehensive performance and security reports

## Installation

```bash
# Install Jetpack using GOPM
gopm jetpack init
```

## Usage

### Basic Monitoring

```bash
# Start monitoring a web application
gopm jetpack monitor http://localhost:3000

# Run a Lighthouse audit
gopm jetpack lighthouse https://example.com
```

### Performance Panel

```bash
# Show the performance panel
gopm jetpack panel show

# Hide the performance panel
gopm jetpack panel hide

# Configure the performance panel
gopm jetpack panel config --position=bottom-right --opacity=0.8 --theme=dark
```

### Metrics Management

```bash
# List available metrics
gopm jetpack metrics list

# Track specific metrics
gopm jetpack metrics track fps memory_usage api_latency

# Stop tracking specific metrics
gopm jetpack metrics untrack memory_usage
```

### Security Monitoring

```bash
# Scan for security vulnerabilities
gopm jetpack security scan

# Check security headers
gopm jetpack security headers https://example.com

# Check TLS configuration
gopm jetpack security tls example.com:443
```

### Exporting and Reporting

```bash
# Export metrics to JSON
gopm jetpack export json --output=metrics.json

# Generate a performance report
gopm jetpack report performance --output=performance-report.html

# Generate a security report
gopm jetpack report security --output=security-report.html

# Generate a full report
gopm jetpack report full --output=full-report.html
```

### Chrome Extension

```bash
# Build the Chrome extension
gopm jetpack chrome build

# Install the Chrome extension
gopm jetpack chrome install

# Update the Chrome extension
gopm jetpack chrome update
```

## Performance Panel

The Jetpack Performance Panel is a floating translucent panel that displays real-time performance metrics on your web page. It helps developers monitor performance metrics without having to open DevTools.

### Features

- **Real-time Metrics**: FPS, memory usage, API latency, and more
- **Customizable**: Position, opacity, theme, and metrics
- **Draggable**: Move the panel anywhere on the screen
- **Collapsible**: Minimize the panel when not needed
- **Alerts**: Visual indicators when metrics exceed thresholds

### Configuration

You can configure the performance panel using the `gopm jetpack panel config` command or through the panel's settings tab.

```bash
# Configure the performance panel
gopm jetpack panel config --position=bottom-right --opacity=0.8 --theme=dark --refresh-rate=1000
```

## Chrome DevTools Extension

The Jetpack Chrome DevTools Extension provides advanced performance monitoring capabilities directly in Chrome DevTools.

### Features

- **Dedicated Panel**: A dedicated panel in Chrome DevTools for Jetpack
- **Advanced Metrics**: More detailed metrics than the floating panel
- **Lighthouse Integration**: Run Lighthouse audits directly from the panel
- **Network Monitoring**: Detailed network request analysis
- **Security Analysis**: Security vulnerability scanning and reporting

### Installation

```bash
# Build and install the Chrome extension
gopm jetpack chrome build
gopm jetpack chrome install
```

## Metrics

Jetpack tracks a wide range of performance metrics across the entire stack:

### Frontend Metrics

- **FPS**: Frames per second
- **Page Load**: Total page load time
- **First Paint**: Time to first paint
- **First Contentful Paint**: Time to first contentful paint
- **Largest Contentful Paint**: Time to largest contentful paint
- **Time to Interactive**: Time until the page is interactive
- **Total Blocking Time**: Total time the main thread was blocked
- **Cumulative Layout Shift**: Measure of visual stability
- **Memory Usage**: JavaScript memory usage
- **Network Requests**: Number and size of network requests
- **Resource Size**: Size of resources (JS, CSS, images, etc.)
- **JS Execution Time**: JavaScript execution time
- **DOM Size**: Number of DOM elements

### Backend Metrics

- **API Latency**: Response time for API requests
- **API Throughput**: Number of requests per second
- **Error Rate**: Percentage of failed requests
- **CPU Usage**: Server CPU usage
- **Memory Usage**: Server memory usage
- **Goroutines**: Number of active goroutines
- **GC Pause**: Garbage collection pause time

### Database Metrics

- **Query Time**: Time to execute database queries
- **Query Count**: Number of database queries
- **Connection Pool**: Database connection pool usage
- **Index Usage**: Database index usage
- **Table Size**: Database table size

### Security Metrics

- **Security Score**: Overall security score
- **Vulnerabilities**: Number of detected vulnerabilities
- **Auth Failures**: Number of authentication failures
- **Suspicious Activity**: Number of suspicious activities

## Security Monitoring

Jetpack includes comprehensive security monitoring capabilities:

### Features

- **Vulnerability Scanning**: Scan for common vulnerabilities (XSS, SQL injection, etc.)
- **Security Headers**: Check for proper security headers
- **TLS Configuration**: Verify secure TLS configuration
- **Authentication Monitoring**: Track authentication failures and brute force attempts
- **Anomaly Detection**: Detect suspicious activities
- **Compliance Checking**: Check compliance with security standards

### Usage

```bash
# Scan for security vulnerabilities
gopm jetpack security scan

# Check security headers
gopm jetpack security headers https://example.com

# Check TLS configuration
gopm jetpack security tls example.com:443
```

## Integration with GoScript Ecosystem

Jetpack integrates seamlessly with the GoScript ecosystem:

- **Gocsx**: Monitor CSS performance and optimize styles
- **WebGPU**: Track WebGPU performance and optimize shaders
- **GoUIX**: Monitor component rendering performance
- **GoScale API**: Track API performance and optimize endpoints
- **GoScale DB**: Monitor database performance and optimize queries

## Configuration

Jetpack uses a configuration file located at `~/.jetpack/config.json` or in the project directory as `.jetpackrc.json`.

Example configuration:

```json
{
  "monitoring": {
    "enabled": true,
    "target": "http://localhost:3000",
    "refresh_rate": 1000,
    "metrics": ["fps", "memory_usage", "api_latency", "page_load"]
  },
  "panel": {
    "enabled": true,
    "position": "bottom-right",
    "opacity": 0.8,
    "theme": "dark",
    "show_charts": true,
    "show_alerts": true
  },
  "lighthouse": {
    "enabled": true,
    "categories": ["performance", "accessibility", "best-practices", "seo", "pwa"],
    "throttling": true
  },
  "security": {
    "enabled": true,
    "vulnerability_scan_enabled": true,
    "auth_tracking_enabled": true,
    "anomaly_detection_enabled": true,
    "compliance_check_enabled": true,
    "scan_interval": 3600
  },
  "export": {
    "enabled": true,
    "format": "json",
    "interval": 3600,
    "path": "./jetpack-metrics.json"
  }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License