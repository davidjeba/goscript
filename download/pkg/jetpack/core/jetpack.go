package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MetricType defines the type of metric being tracked
type MetricType string

const (
	// Frontend metric types
	MetricFPS            MetricType = "fps"
	MetricPageLoad       MetricType = "page_load"
	MetricFirstPaint     MetricType = "first_paint"
	MetricFirstContentful MetricType = "first_contentful_paint"
	MetricLargestContentful MetricType = "largest_contentful_paint"
	MetricTTI            MetricType = "time_to_interactive"
	MetricTBT            MetricType = "total_blocking_time"
	MetricCLS            MetricType = "cumulative_layout_shift"
	MetricMemoryUsage    MetricType = "memory_usage"
	MetricNetworkRequests MetricType = "network_requests"
	MetricResourceSize   MetricType = "resource_size"
	MetricJSExecution    MetricType = "js_execution_time"
	MetricDOMSize        MetricType = "dom_size"
	
	// Backend metric types
	MetricAPILatency     MetricType = "api_latency"
	MetricAPIThroughput  MetricType = "api_throughput"
	MetricErrorRate      MetricType = "error_rate"
	MetricCPUUsage       MetricType = "cpu_usage"
	MetricMemoryUsageServer MetricType = "memory_usage_server"
	MetricGoroutines     MetricType = "goroutines"
	MetricGCPause        MetricType = "gc_pause"
	
	// Database metric types
	MetricQueryTime      MetricType = "query_time"
	MetricQueryCount     MetricType = "query_count"
	MetricConnectionPool MetricType = "connection_pool"
	MetricIndexUsage     MetricType = "index_usage"
	MetricTableSize      MetricType = "table_size"
	
	// Security metric types
	MetricSecurityScore  MetricType = "security_score"
	MetricVulnerabilities MetricType = "vulnerabilities"
	MetricAuthFailures   MetricType = "auth_failures"
	MetricSuspiciousActivity MetricType = "suspicious_activity"
)

// MetricValue represents a single metric value
type MetricValue struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// Metric represents a performance metric
type Metric struct {
	Type        MetricType    `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Unit        string        `json:"unit"`
	Values      []MetricValue `json:"values"`
	Threshold   *float64      `json:"threshold,omitempty"`
	Alert       bool          `json:"alert"`
	Tags        []string      `json:"tags"`
	mutex       sync.RWMutex
}

// Jetpack is the main performance monitoring system
type Jetpack struct {
	Metrics        map[string]*Metric
	DevMode        bool
	PanelVisible   bool
	PanelPosition  string
	PanelOpacity   float64
	RefreshRate    time.Duration
	AlertThreshold float64
	AlertCallback  func(metric *Metric)
	ExportEnabled  bool
	ExportEndpoint string
	ExportInterval time.Duration
	mutex          sync.RWMutex
	
	// Components
	Frontend *FrontendMonitor
	Backend  *BackendMonitor
	Database *DatabaseMonitor
	Security *SecurityMonitor
}

// FrontendMonitor tracks frontend performance metrics
type FrontendMonitor struct {
	Jetpack        *Jetpack
	LighthouseEnabled bool
	WebVitalsEnabled bool
	ResourceTrackingEnabled bool
	NetworkTrackingEnabled bool
	MemoryTrackingEnabled bool
	FPSTrackingEnabled bool
	UserTimingEnabled bool
}

// BackendMonitor tracks backend performance metrics
type BackendMonitor struct {
	Jetpack        *Jetpack
	APITrackingEnabled bool
	SystemMetricsEnabled bool
	TraceEnabled bool
	ProfileEnabled bool
	LogAnalysisEnabled bool
}

// DatabaseMonitor tracks database performance metrics
type DatabaseMonitor struct {
	Jetpack        *Jetpack
	QueryTrackingEnabled bool
	ConnectionPoolTracking bool
	SchemaAnalysisEnabled bool
	IndexAnalysisEnabled bool
	SlowQueryThreshold time.Duration
}

// SecurityMonitor tracks security metrics and vulnerabilities
type SecurityMonitor struct {
	Jetpack        *Jetpack
	VulnerabilityScanEnabled bool
	AuthTrackingEnabled bool
	AnomalyDetectionEnabled bool
	ComplianceCheckEnabled bool
	ScanInterval time.Duration
}

// NewJetpack creates a new Jetpack instance
func NewJetpack() *Jetpack {
	jp := &Jetpack{
		Metrics:        make(map[string]*Metric),
		DevMode:        false,
		PanelVisible:   false,
		PanelPosition:  "bottom-right",
		PanelOpacity:   0.8,
		RefreshRate:    time.Second,
		AlertThreshold: 0.9,
		ExportEnabled:  false,
		ExportEndpoint: "",
		ExportInterval: time.Minute,
	}
	
	// Initialize components
	jp.Frontend = &FrontendMonitor{
		Jetpack:        jp,
		LighthouseEnabled: true,
		WebVitalsEnabled: true,
		ResourceTrackingEnabled: true,
		NetworkTrackingEnabled: true,
		MemoryTrackingEnabled: true,
		FPSTrackingEnabled: true,
		UserTimingEnabled: true,
	}
	
	jp.Backend = &BackendMonitor{
		Jetpack:        jp,
		APITrackingEnabled: true,
		SystemMetricsEnabled: true,
		TraceEnabled: true,
		ProfileEnabled: true,
		LogAnalysisEnabled: true,
	}
	
	jp.Database = &DatabaseMonitor{
		Jetpack:        jp,
		QueryTrackingEnabled: true,
		ConnectionPoolTracking: true,
		SchemaAnalysisEnabled: true,
		IndexAnalysisEnabled: true,
		SlowQueryThreshold: time.Millisecond * 100,
	}
	
	jp.Security = &SecurityMonitor{
		Jetpack:        jp,
		VulnerabilityScanEnabled: true,
		AuthTrackingEnabled: true,
		AnomalyDetectionEnabled: true,
		ComplianceCheckEnabled: true,
		ScanInterval: time.Hour,
	}
	
	return jp
}

// EnableDevMode enables developer mode with performance panel
func (jp *Jetpack) EnableDevMode() {
	jp.mutex.Lock()
	defer jp.mutex.Unlock()
	
	jp.DevMode = true
	jp.PanelVisible = true
}

// DisableDevMode disables developer mode
func (jp *Jetpack) DisableDevMode() {
	jp.mutex.Lock()
	defer jp.mutex.Unlock()
	
	jp.DevMode = false
	jp.PanelVisible = false
}

// SetPanelPosition sets the position of the performance panel
func (jp *Jetpack) SetPanelPosition(position string) {
	jp.mutex.Lock()
	defer jp.mutex.Unlock()
	
	jp.PanelPosition = position
}

// SetPanelOpacity sets the opacity of the performance panel
func (jp *Jetpack) SetPanelOpacity(opacity float64) {
	jp.mutex.Lock()
	defer jp.mutex.Unlock()
	
	if opacity < 0 {
		opacity = 0
	} else if opacity > 1 {
		opacity = 1
	}
	
	jp.PanelOpacity = opacity
}

// RegisterMetric registers a new metric
func (jp *Jetpack) RegisterMetric(metricType MetricType, name, description, unit string, threshold *float64, tags []string) *Metric {
	jp.mutex.Lock()
	defer jp.mutex.Unlock()
	
	metric := &Metric{
		Type:        metricType,
		Name:        name,
		Description: description,
		Unit:        unit,
		Values:      make([]MetricValue, 0),
		Threshold:   threshold,
		Alert:       false,
		Tags:        tags,
	}
	
	jp.Metrics[name] = metric
	
	return metric
}

// GetMetric gets a metric by name
func (jp *Jetpack) GetMetric(name string) (*Metric, error) {
	jp.mutex.RLock()
	defer jp.mutex.RUnlock()
	
	metric, ok := jp.Metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %s not found", name)
	}
	
	return metric, nil
}

// RecordMetric records a metric value
func (jp *Jetpack) RecordMetric(name string, value float64) error {
	metric, err := jp.GetMetric(name)
	if err != nil {
		return err
	}
	
	metric.mutex.Lock()
	defer metric.mutex.Unlock()
	
	metricValue := MetricValue{
		Value:     value,
		Timestamp: time.Now(),
	}
	
	metric.Values = append(metric.Values, metricValue)
	
	// Check threshold
	if metric.Threshold != nil && value >= *metric.Threshold {
		metric.Alert = true
		
		// Call alert callback if set
		if jp.AlertCallback != nil {
			go jp.AlertCallback(metric)
		}
	}
	
	return nil
}

// GetMetricAverage gets the average value of a metric
func (jp *Jetpack) GetMetricAverage(name string) (float64, error) {
	metric, err := jp.GetMetric(name)
	if err != nil {
		return 0, err
	}
	
	metric.mutex.RLock()
	defer metric.mutex.RUnlock()
	
	if len(metric.Values) == 0 {
		return 0, nil
	}
	
	var sum float64
	for _, value := range metric.Values {
		sum += value.Value
	}
	
	return sum / float64(len(metric.Values)), nil
}

// GetMetricLatest gets the latest value of a metric
func (jp *Jetpack) GetMetricLatest(name string) (float64, error) {
	metric, err := jp.GetMetric(name)
	if err != nil {
		return 0, err
	}
	
	metric.mutex.RLock()
	defer metric.mutex.RUnlock()
	
	if len(metric.Values) == 0 {
		return 0, nil
	}
	
	return metric.Values[len(metric.Values)-1].Value, nil
}

// ExportMetrics exports metrics to the configured endpoint
func (jp *Jetpack) ExportMetrics() error {
	if !jp.ExportEnabled || jp.ExportEndpoint == "" {
		return nil
	}
	
	jp.mutex.RLock()
	defer jp.mutex.RUnlock()
	
	// In a real implementation, this would send metrics to the endpoint
	// For now, we'll just serialize to JSON
	data, err := json.Marshal(jp.Metrics)
	if err != nil {
		return err
	}
	
	fmt.Printf("Exporting metrics to %s: %s\n", jp.ExportEndpoint, string(data))
	
	return nil
}

// StartExporting starts exporting metrics at the configured interval
func (jp *Jetpack) StartExporting() {
	if !jp.ExportEnabled || jp.ExportEndpoint == "" {
		return
	}
	
	go func() {
		ticker := time.NewTicker(jp.ExportInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			jp.ExportMetrics()
		}
	}()
}

// GetPanelData gets data for the performance panel
func (jp *Jetpack) GetPanelData() map[string]interface{} {
	jp.mutex.RLock()
	defer jp.mutex.RUnlock()
	
	data := make(map[string]interface{})
	
	// Add basic info
	data["dev_mode"] = jp.DevMode
	data["panel_visible"] = jp.PanelVisible
	data["panel_position"] = jp.PanelPosition
	data["panel_opacity"] = jp.PanelOpacity
	
	// Add metrics
	metrics := make(map[string]interface{})
	for name, metric := range jp.Metrics {
		metric.mutex.RLock()
		
		metricData := map[string]interface{}{
			"type":        metric.Type,
			"description": metric.Description,
			"unit":        metric.Unit,
			"alert":       metric.Alert,
			"tags":        metric.Tags,
		}
		
		if len(metric.Values) > 0 {
			latest := metric.Values[len(metric.Values)-1]
			metricData["latest_value"] = latest.Value
			metricData["latest_timestamp"] = latest.Timestamp
		}
		
		if metric.Threshold != nil {
			metricData["threshold"] = *metric.Threshold
		}
		
		metric.mutex.RUnlock()
		
		metrics[name] = metricData
	}
	
	data["metrics"] = metrics
	
	return data
}