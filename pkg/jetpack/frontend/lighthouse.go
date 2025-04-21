package frontend

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davidjeba/goscript/pkg/jetpack/core"
)

// LighthouseConfig represents the configuration for Lighthouse integration
type LighthouseConfig struct {
	Enabled       bool     `json:"enabled"`
	Categories    []string `json:"categories"`
	Locale        string   `json:"locale"`
	MaxWaitTime   int      `json:"max_wait_time"`
	FormFactor    string   `json:"form_factor"`
	Throttling    bool     `json:"throttling"`
	OnlyAudits    []string `json:"only_audits,omitempty"`
	SkipAudits    []string `json:"skip_audits,omitempty"`
	EmulatedDevice string  `json:"emulated_device,omitempty"`
}

// LighthouseResult represents the result of a Lighthouse audit
type LighthouseResult struct {
	URL           string                 `json:"url"`
	Categories    map[string]float64     `json:"categories"`
	Audits        map[string]interface{} `json:"audits"`
	Timestamp     time.Time              `json:"timestamp"`
	LighthouseVersion string             `json:"lighthouse_version"`
	UserAgent     string                 `json:"user_agent"`
	Environment   map[string]interface{} `json:"environment"`
}

// LighthouseMonitor integrates with Google Lighthouse for web performance analysis
type LighthouseMonitor struct {
	Jetpack       *core.Jetpack
	Config        LighthouseConfig
	Results       []*LighthouseResult
	LastRunTime   time.Time
	RunCount      int
	AutoRunEnabled bool
	AutoRunInterval time.Duration
}

// NewLighthouseMonitor creates a new Lighthouse monitor
func NewLighthouseMonitor(jetpack *core.Jetpack) *LighthouseMonitor {
	return &LighthouseMonitor{
		Jetpack: jetpack,
		Config: LighthouseConfig{
			Enabled:       true,
			Categories:    []string{"performance", "accessibility", "best-practices", "seo", "pwa"},
			Locale:        "en-US",
			MaxWaitTime:   45,
			FormFactor:    "desktop",
			Throttling:    true,
			EmulatedDevice: "Nexus 5X",
		},
		Results:       make([]*LighthouseResult, 0),
		AutoRunEnabled: false,
		AutoRunInterval: time.Hour,
	}
}

// RunAudit runs a Lighthouse audit on the specified URL
func (lm *LighthouseMonitor) RunAudit(url string) (*LighthouseResult, error) {
	// In a real implementation, this would call the Lighthouse API
	// For now, we'll simulate a result
	
	result := &LighthouseResult{
		URL: url,
		Categories: map[string]float64{
			"performance":    0.85,
			"accessibility":  0.92,
			"best-practices": 0.87,
			"seo":            0.95,
			"pwa":            0.65,
		},
		Audits: map[string]interface{}{
			"first-contentful-paint": map[string]interface{}{
				"id":    "first-contentful-paint",
				"title": "First Contentful Paint",
				"description": "First Contentful Paint marks the time at which the first text or image is painted",
				"score": 0.89,
				"displayValue": "1.2 s",
				"numericValue": 1243,
			},
			"speed-index": map[string]interface{}{
				"id":    "speed-index",
				"title": "Speed Index",
				"description": "Speed Index shows how quickly the contents of a page are visibly populated",
				"score": 0.87,
				"displayValue": "1.8 s",
				"numericValue": 1823,
			},
			"largest-contentful-paint": map[string]interface{}{
				"id":    "largest-contentful-paint",
				"title": "Largest Contentful Paint",
				"description": "Largest Contentful Paint marks the time at which the largest text or image is painted",
				"score": 0.82,
				"displayValue": "2.1 s",
				"numericValue": 2134,
			},
			"total-blocking-time": map[string]interface{}{
				"id":    "total-blocking-time",
				"title": "Total Blocking Time",
				"description": "Sum of all time periods between FCP and Time to Interactive, when task length exceeded 50ms",
				"score": 0.75,
				"displayValue": "120 ms",
				"numericValue": 120,
			},
			"cumulative-layout-shift": map[string]interface{}{
				"id":    "cumulative-layout-shift",
				"title": "Cumulative Layout Shift",
				"description": "Cumulative Layout Shift measures the movement of visible elements within the viewport",
				"score": 0.92,
				"displayValue": "0.05",
				"numericValue": 0.05,
			},
		},
		Timestamp: time.Now(),
		LighthouseVersion: "9.6.0",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		Environment: map[string]interface{}{
			"networkUserAgent": "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36",
			"hostUserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"benchmarkIndex": 1000,
		},
	}
	
	lm.Results = append(lm.Results, result)
	lm.LastRunTime = time.Now()
	lm.RunCount++
	
	// Record metrics in Jetpack
	lm.recordMetricsFromResult(result)
	
	return result, nil
}

// recordMetricsFromResult records metrics from a Lighthouse result
func (lm *LighthouseMonitor) recordMetricsFromResult(result *LighthouseResult) {
	// Record category scores
	for category, score := range result.Categories {
		metricName := fmt.Sprintf("lighthouse_%s_score", category)
		
		// Check if metric exists, if not register it
		_, err := lm.Jetpack.GetMetric(metricName)
		if err != nil {
			threshold := 0.75 // 75% is a good threshold for Lighthouse scores
			lm.Jetpack.RegisterMetric(
				core.MetricType(fmt.Sprintf("lighthouse_%s", category)),
				metricName,
				fmt.Sprintf("Lighthouse %s score", category),
				"score",
				&threshold,
				[]string{"lighthouse", category},
			)
		}
		
		// Record the metric
		lm.Jetpack.RecordMetric(metricName, score)
	}
	
	// Record key performance metrics
	keyMetrics := map[string]string{
		"first-contentful-paint":   "first_contentful_paint",
		"largest-contentful-paint": "largest_contentful_paint",
		"speed-index":              "speed_index",
		"total-blocking-time":      "total_blocking_time",
		"cumulative-layout-shift":  "cumulative_layout_shift",
	}
	
	for auditID, jetpackMetric := range keyMetrics {
		if audit, ok := result.Audits[auditID].(map[string]interface{}); ok {
			if numericValue, ok := audit["numericValue"].(float64); ok {
				metricName := fmt.Sprintf("lighthouse_%s", jetpackMetric)
				
				// Check if metric exists, if not register it
				_, err := lm.Jetpack.GetMetric(metricName)
				if err != nil {
					var threshold *float64
					
					// Set appropriate thresholds based on metric
					switch auditID {
					case "first-contentful-paint":
						t := 2000.0 // 2 seconds
						threshold = &t
					case "largest-contentful-paint":
						t := 2500.0 // 2.5 seconds
						threshold = &t
					case "speed-index":
						t := 3000.0 // 3 seconds
						threshold = &t
					case "total-blocking-time":
						t := 200.0 // 200 milliseconds
						threshold = &t
					case "cumulative-layout-shift":
						t := 0.1 // 0.1 CLS
						threshold = &t
					}
					
					unit := "ms"
					if auditID == "cumulative-layout-shift" {
						unit = "score"
					}
					
					lm.Jetpack.RegisterMetric(
						core.MetricType(jetpackMetric),
						metricName,
						audit["title"].(string),
						unit,
						threshold,
						[]string{"lighthouse", "performance"},
					)
				}
				
				// Record the metric
				lm.Jetpack.RecordMetric(metricName, numericValue)
			}
		}
	}
}

// StartAutoRun starts automatically running Lighthouse audits
func (lm *LighthouseMonitor) StartAutoRun(url string) {
	if !lm.AutoRunEnabled {
		lm.AutoRunEnabled = true
		
		go func() {
			ticker := time.NewTicker(lm.AutoRunInterval)
			defer ticker.Stop()
			
			for range ticker.C {
				if !lm.AutoRunEnabled {
					break
				}
				
				lm.RunAudit(url)
			}
		}()
	}
}

// StopAutoRun stops automatically running Lighthouse audits
func (lm *LighthouseMonitor) StopAutoRun() {
	lm.AutoRunEnabled = false
}

// GetLatestResult gets the latest Lighthouse result
func (lm *LighthouseMonitor) GetLatestResult() *LighthouseResult {
	if len(lm.Results) == 0 {
		return nil
	}
	
	return lm.Results[len(lm.Results)-1]
}

// ExportResultToJSON exports a Lighthouse result to JSON
func (lm *LighthouseMonitor) ExportResultToJSON(result *LighthouseResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// GenerateReport generates a human-readable report from a Lighthouse result
func (lm *LighthouseMonitor) GenerateReport(result *LighthouseResult) string {
	if result == nil {
		return "No result available"
	}
	
	report := fmt.Sprintf("Lighthouse Report for %s\n", result.URL)
	report += fmt.Sprintf("Generated at: %s\n", result.Timestamp.Format(time.RFC1123))
	report += fmt.Sprintf("Lighthouse Version: %s\n\n", result.LighthouseVersion)
	
	report += "Category Scores:\n"
	for category, score := range result.Categories {
		report += fmt.Sprintf("  %s: %.2f\n", category, score)
	}
	
	report += "\nKey Metrics:\n"
	for auditID, audit := range result.Audits {
		if auditMap, ok := audit.(map[string]interface{}); ok {
			if title, ok := auditMap["title"].(string); ok {
				if displayValue, ok := auditMap["displayValue"].(string); ok {
					if score, ok := auditMap["score"].(float64); ok {
						report += fmt.Sprintf("  %s: %s (Score: %.2f)\n", title, displayValue, score)
					}
				}
			}
		}
	}
	
	return report
}