package security

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/davidjeba/goscript/pkg/jetpack/core"
)

// SecurityLevel defines the security level
type SecurityLevel string

const (
	SecurityLevelLow    SecurityLevel = "low"
	SecurityLevelMedium SecurityLevel = "medium"
	SecurityLevelHigh   SecurityLevel = "high"
	SecurityLevelCritical SecurityLevel = "critical"
)

// VulnerabilityType defines the type of vulnerability
type VulnerabilityType string

const (
	VulnXSS             VulnerabilityType = "xss"
	VulnSQLInjection    VulnerabilityType = "sql_injection"
	VulnCSRF            VulnerabilityType = "csrf"
	VulnAuthBypass      VulnerabilityType = "auth_bypass"
	VulnInsecureCrypto  VulnerabilityType = "insecure_crypto"
	VulnWeakPassword    VulnerabilityType = "weak_password"
	VulnMissingHeaders  VulnerabilityType = "missing_headers"
	VulnOutdatedLibrary VulnerabilityType = "outdated_library"
	VulnMisconfiguration VulnerabilityType = "misconfiguration"
	VulnDataExposure    VulnerabilityType = "data_exposure"
)

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string            `json:"id"`
	Type        VulnerabilityType `json:"type"`
	Level       SecurityLevel     `json:"level"`
	Description string            `json:"description"`
	Location    string            `json:"location"`
	Timestamp   time.Time         `json:"timestamp"`
	Remediation string            `json:"remediation"`
	References  []string          `json:"references"`
	Fixed       bool              `json:"fixed"`
}

// SecurityConfig represents the configuration for security monitoring
type SecurityConfig struct {
	Enabled                bool          `json:"enabled"`
	VulnerabilityScanEnabled bool        `json:"vulnerability_scan_enabled"`
	AuthTrackingEnabled    bool          `json:"auth_tracking_enabled"`
	AnomalyDetectionEnabled bool         `json:"anomaly_detection_enabled"`
	ComplianceCheckEnabled bool          `json:"compliance_check_enabled"`
	ScanInterval           time.Duration `json:"scan_interval"`
	AlertThreshold         SecurityLevel `json:"alert_threshold"`
	AutoFix                bool          `json:"auto_fix"`
	ReportPath             string        `json:"report_path"`
	ExcludePaths           []string      `json:"exclude_paths"`
}

// SecurityMonitor monitors security vulnerabilities and issues
type SecurityMonitor struct {
	Jetpack         *core.Jetpack
	Config          SecurityConfig
	Vulnerabilities map[string]*Vulnerability
	AuthFailures    map[string]int
	SuspiciousActivities []string
	LastScanTime    time.Time
	ScanCount       int
	mutex           sync.RWMutex
}

// NewSecurityMonitor creates a new security monitor
func NewSecurityMonitor(jetpack *core.Jetpack) *SecurityMonitor {
	return &SecurityMonitor{
		Jetpack: jetpack,
		Config: SecurityConfig{
			Enabled:                true,
			VulnerabilityScanEnabled: true,
			AuthTrackingEnabled:    true,
			AnomalyDetectionEnabled: true,
			ComplianceCheckEnabled: true,
			ScanInterval:           time.Hour,
			AlertThreshold:         SecurityLevelMedium,
			AutoFix:                false,
			ReportPath:             "security_report.json",
			ExcludePaths:           []string{"/assets/", "/public/"},
		},
		Vulnerabilities:     make(map[string]*Vulnerability),
		AuthFailures:        make(map[string]int),
		SuspiciousActivities: make([]string, 0),
		LastScanTime:        time.Time{},
		ScanCount:           0,
	}
}

// StartMonitoring starts security monitoring
func (sm *SecurityMonitor) StartMonitoring() {
	if !sm.Config.Enabled {
		return
	}
	
	// Register security metrics
	sm.registerSecurityMetrics()
	
	// Start vulnerability scanning
	if sm.Config.VulnerabilityScanEnabled {
		go sm.startVulnerabilityScan()
	}
	
	// Start anomaly detection
	if sm.Config.AnomalyDetectionEnabled {
		go sm.startAnomalyDetection()
	}
	
	// Start compliance checking
	if sm.Config.ComplianceCheckEnabled {
		go sm.startComplianceCheck()
	}
}

// registerSecurityMetrics registers security metrics with Jetpack
func (sm *SecurityMonitor) registerSecurityMetrics() {
	// Security score
	threshold := 80.0
	sm.Jetpack.RegisterMetric(
		core.MetricSecurityScore,
		"security_score",
		"Overall security score",
		"score",
		&threshold,
		[]string{"security"},
	)
	
	// Vulnerabilities
	threshold = 5.0
	sm.Jetpack.RegisterMetric(
		core.MetricVulnerabilities,
		"vulnerabilities",
		"Number of detected vulnerabilities",
		"count",
		&threshold,
		[]string{"security"},
	)
	
	// Auth failures
	threshold = 10.0
	sm.Jetpack.RegisterMetric(
		core.MetricAuthFailures,
		"auth_failures",
		"Number of authentication failures",
		"count",
		&threshold,
		[]string{"security"},
	)
	
	// Suspicious activity
	threshold = 5.0
	sm.Jetpack.RegisterMetric(
		core.MetricSuspiciousActivity,
		"suspicious_activity",
		"Number of suspicious activities",
		"count",
		&threshold,
		[]string{"security"},
	)
}

// startVulnerabilityScan starts periodic vulnerability scanning
func (sm *SecurityMonitor) startVulnerabilityScan() {
	ticker := time.NewTicker(sm.Config.ScanInterval)
	defer ticker.Stop()
	
	// Run initial scan
	sm.ScanVulnerabilities()
	
	for range ticker.C {
		sm.ScanVulnerabilities()
	}
}

// startAnomalyDetection starts anomaly detection
func (sm *SecurityMonitor) startAnomalyDetection() {
	ticker := time.NewTicker(sm.Config.ScanInterval / 2)
	defer ticker.Stop()
	
	for range ticker.C {
		sm.DetectAnomalies()
	}
}

// startComplianceCheck starts compliance checking
func (sm *SecurityMonitor) startComplianceCheck() {
	ticker := time.NewTicker(sm.Config.ScanInterval * 2)
	defer ticker.Stop()
	
	// Run initial check
	sm.CheckCompliance()
	
	for range ticker.C {
		sm.CheckCompliance()
	}
}

// ScanVulnerabilities scans for vulnerabilities
func (sm *SecurityMonitor) ScanVulnerabilities() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.LastScanTime = time.Now()
	sm.ScanCount++
	
	// In a real implementation, this would scan for actual vulnerabilities
	// For now, we'll simulate finding some vulnerabilities
	
	// Clear fixed vulnerabilities
	for id, vuln := range sm.Vulnerabilities {
		if vuln.Fixed {
			delete(sm.Vulnerabilities, id)
		}
	}
	
	// Simulate finding vulnerabilities
	vulnerabilities := []*Vulnerability{
		{
			ID:          "VULN-001",
			Type:        VulnMissingHeaders,
			Level:       SecurityLevelMedium,
			Description: "Missing Content-Security-Policy header",
			Location:    "HTTP Response Headers",
			Timestamp:   time.Now(),
			Remediation: "Add Content-Security-Policy header to all responses",
			References:  []string{"https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP"},
			Fixed:       false,
		},
		{
			ID:          "VULN-002",
			Type:        VulnOutdatedLibrary,
			Level:       SecurityLevelHigh,
			Description: "Using outdated library with known vulnerabilities",
			Location:    "package.json",
			Timestamp:   time.Now(),
			Remediation: "Update the library to the latest version",
			References:  []string{"https://nvd.nist.gov/vuln/detail/CVE-2021-12345"},
			Fixed:       false,
		},
	}
	
	// Add vulnerabilities to the map
	for _, vuln := range vulnerabilities {
		if _, ok := sm.Vulnerabilities[vuln.ID]; !ok {
			sm.Vulnerabilities[vuln.ID] = vuln
		}
	}
	
	// Update metrics
	sm.updateSecurityMetrics()
}

// DetectAnomalies detects security anomalies
func (sm *SecurityMonitor) DetectAnomalies() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// In a real implementation, this would analyze logs and detect anomalies
	// For now, we'll simulate finding some anomalies
	
	// Clear old activities
	sm.SuspiciousActivities = make([]string, 0)
	
	// Simulate finding anomalies
	sm.SuspiciousActivities = append(sm.SuspiciousActivities,
		fmt.Sprintf("Unusual login pattern detected from IP 192.168.1.100 at %s", time.Now().Format(time.RFC3339)),
		fmt.Sprintf("Multiple failed API requests from IP 192.168.1.200 at %s", time.Now().Format(time.RFC3339)),
	)
	
	// Update metrics
	sm.updateSecurityMetrics()
}

// CheckCompliance checks security compliance
func (sm *SecurityMonitor) CheckCompliance() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// In a real implementation, this would check for compliance with security standards
	// For now, we'll simulate a compliance check
	
	// Simulate compliance check
	complianceIssues := []string{
		"Missing HTTPS redirection",
		"Weak TLS configuration",
		"Missing rate limiting",
	}
	
	// Log compliance issues
	for _, issue := range complianceIssues {
		fmt.Printf("Compliance issue: %s\n", issue)
	}
	
	// Update metrics
	securityScore := 100.0 - (float64(len(complianceIssues)) * 10.0)
	if securityScore < 0 {
		securityScore = 0
	}
	
	sm.Jetpack.RecordMetric("security_score", securityScore)
}

// TrackAuthFailure tracks authentication failures
func (sm *SecurityMonitor) TrackAuthFailure(username, ipAddress string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	key := fmt.Sprintf("%s:%s", username, ipAddress)
	sm.AuthFailures[key]++
	
	// Check for brute force attempts
	if sm.AuthFailures[key] >= 5 {
		sm.SuspiciousActivities = append(sm.SuspiciousActivities,
			fmt.Sprintf("Possible brute force attack: %d failed login attempts for user %s from IP %s", 
				sm.AuthFailures[key], username, ipAddress))
	}
	
	// Update metrics
	sm.updateSecurityMetrics()
}

// ResetAuthFailures resets authentication failures for a user
func (sm *SecurityMonitor) ResetAuthFailures(username, ipAddress string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	key := fmt.Sprintf("%s:%s", username, ipAddress)
	delete(sm.AuthFailures, key)
	
	// Update metrics
	sm.updateSecurityMetrics()
}

// updateSecurityMetrics updates security metrics
func (sm *SecurityMonitor) updateSecurityMetrics() {
	// Count vulnerabilities by severity
	vulnCount := len(sm.Vulnerabilities)
	highVulnCount := 0
	for _, vuln := range sm.Vulnerabilities {
		if vuln.Level == SecurityLevelHigh || vuln.Level == SecurityLevelCritical {
			highVulnCount++
		}
	}
	
	// Calculate security score
	securityScore := 100.0
	securityScore -= float64(vulnCount) * 5.0
	securityScore -= float64(highVulnCount) * 10.0
	securityScore -= float64(len(sm.SuspiciousActivities)) * 5.0
	
	// Count auth failures
	authFailureCount := 0
	for _, count := range sm.AuthFailures {
		authFailureCount += count
	}
	securityScore -= float64(authFailureCount) * 2.0
	
	if securityScore < 0 {
		securityScore = 0
	}
	
	// Update metrics
	sm.Jetpack.RecordMetric("security_score", securityScore)
	sm.Jetpack.RecordMetric("vulnerabilities", float64(vulnCount))
	sm.Jetpack.RecordMetric("auth_failures", float64(authFailureCount))
	sm.Jetpack.RecordMetric("suspicious_activity", float64(len(sm.SuspiciousActivities)))
}

// GetVulnerabilities gets all vulnerabilities
func (sm *SecurityMonitor) GetVulnerabilities() []*Vulnerability {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	vulnerabilities := make([]*Vulnerability, 0, len(sm.Vulnerabilities))
	for _, vuln := range sm.Vulnerabilities {
		vulnerabilities = append(vulnerabilities, vuln)
	}
	
	return vulnerabilities
}

// GetVulnerability gets a vulnerability by ID
func (sm *SecurityMonitor) GetVulnerability(id string) (*Vulnerability, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	vuln, ok := sm.Vulnerabilities[id]
	if !ok {
		return nil, fmt.Errorf("vulnerability %s not found", id)
	}
	
	return vuln, nil
}

// FixVulnerability marks a vulnerability as fixed
func (sm *SecurityMonitor) FixVulnerability(id string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	vuln, ok := sm.Vulnerabilities[id]
	if !ok {
		return fmt.Errorf("vulnerability %s not found", id)
	}
	
	vuln.Fixed = true
	
	// Update metrics
	sm.updateSecurityMetrics()
	
	return nil
}

// GetAuthFailures gets all authentication failures
func (sm *SecurityMonitor) GetAuthFailures() map[string]int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	// Create a copy to avoid race conditions
	authFailures := make(map[string]int)
	for key, count := range sm.AuthFailures {
		authFailures[key] = count
	}
	
	return authFailures
}

// GetSuspiciousActivities gets all suspicious activities
func (sm *SecurityMonitor) GetSuspiciousActivities() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	// Create a copy to avoid race conditions
	activities := make([]string, len(sm.SuspiciousActivities))
	copy(activities, sm.SuspiciousActivities)
	
	return activities
}

// GenerateReport generates a security report
func (sm *SecurityMonitor) GenerateReport() (string, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	report := map[string]interface{}{
		"timestamp":            time.Now(),
		"security_score":       0.0,
		"vulnerabilities":      sm.Vulnerabilities,
		"auth_failures":        sm.AuthFailures,
		"suspicious_activities": sm.SuspiciousActivities,
		"last_scan_time":       sm.LastScanTime,
		"scan_count":           sm.ScanCount,
	}
	
	// Get security score
	securityScore, err := sm.Jetpack.GetMetricLatest("security_score")
	if err == nil {
		report["security_score"] = securityScore
	}
	
	// Convert to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// CheckTLSConfiguration checks the TLS configuration of a server
func (sm *SecurityMonitor) CheckTLSConfiguration(host string) (map[string]interface{}, error) {
	// Connect to the server
	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	
	// Get the TLS connection state
	state := conn.ConnectionState()
	
	// Check TLS version
	var tlsVersion string
	switch state.Version {
	case tls.VersionTLS10:
		tlsVersion = "TLS 1.0"
	case tls.VersionTLS11:
		tlsVersion = "TLS 1.1"
	case tls.VersionTLS12:
		tlsVersion = "TLS 1.2"
	case tls.VersionTLS13:
		tlsVersion = "TLS 1.3"
	default:
		tlsVersion = fmt.Sprintf("Unknown (%d)", state.Version)
	}
	
	// Check cipher suite
	cipherSuite := tls.CipherSuiteName(state.CipherSuite)
	
	// Check certificate
	cert := state.PeerCertificates[0]
	
	// Build result
	result := map[string]interface{}{
		"host":            host,
		"tls_version":     tlsVersion,
		"cipher_suite":    cipherSuite,
		"certificate": map[string]interface{}{
			"subject":      cert.Subject.String(),
			"issuer":       cert.Issuer.String(),
			"not_before":   cert.NotBefore,
			"not_after":    cert.NotAfter,
			"dns_names":    cert.DNSNames,
			"serial_number": cert.SerialNumber.String(),
		},
		"secure": state.Version >= tls.VersionTLS12,
	}
	
	return result, nil
}

// CheckSecurityHeaders checks security headers in an HTTP response
func (sm *SecurityMonitor) CheckSecurityHeaders(url string) (map[string]interface{}, error) {
	// Make a request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Check security headers
	headers := resp.Header
	
	// Required security headers
	requiredHeaders := map[string]string{
		"Content-Security-Policy": "",
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options": "DENY",
		"X-XSS-Protection": "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Referrer-Policy": "no-referrer",
		"Permissions-Policy": "",
	}
	
	// Check each header
	headerResults := make(map[string]interface{})
	for header, expectedValue := range requiredHeaders {
		value := headers.Get(header)
		if value == "" {
			headerResults[header] = map[string]interface{}{
				"present": false,
				"value":   nil,
				"valid":   false,
			}
		} else {
			valid := true
			if expectedValue != "" && !strings.Contains(value, expectedValue) {
				valid = false
			}
			
			headerResults[header] = map[string]interface{}{
				"present": true,
				"value":   value,
				"valid":   valid,
			}
		}
	}
	
	// Calculate score
	score := 0.0
	maxScore := float64(len(requiredHeaders))
	for _, result := range headerResults {
		if r, ok := result.(map[string]interface{}); ok {
			if r["present"].(bool) {
				score += 0.5
				if r["valid"].(bool) {
					score += 0.5
				}
			}
		}
	}
	
	// Build result
	result := map[string]interface{}{
		"url":      url,
		"headers":  headerResults,
		"score":    score,
		"max_score": maxScore,
		"percentage": (score / maxScore) * 100.0,
	}
	
	return result, nil
}

// ScanForSQLInjection scans for SQL injection vulnerabilities
func (sm *SecurityMonitor) ScanForSQLInjection(url string, params map[string]string) (map[string]interface{}, error) {
	// SQL injection payloads
	payloads := []string{
		"' OR '1'='1",
		"1' OR '1'='1",
		"' OR 1=1--",
		"' OR 1=1#",
		"') OR 1=1--",
		"admin'--",
	}
	
	// Build result
	result := map[string]interface{}{
		"url":      url,
		"params":   params,
		"vulnerable": false,
		"details":  []map[string]interface{}{},
	}
	
	// In a real implementation, this would test each payload against each parameter
	// For now, we'll simulate the scan
	
	// Simulate finding a vulnerability
	if params["id"] != "" {
		result["vulnerable"] = true
		result["details"] = append(result["details"].([]map[string]interface{}), map[string]interface{}{
			"param":   "id",
			"payload": payloads[0],
			"response": "Database error occurred",
		})
	}
	
	return result, nil
}

// ScanForXSS scans for XSS vulnerabilities
func (sm *SecurityMonitor) ScanForXSS(url string, params map[string]string) (map[string]interface{}, error) {
	// XSS payloads
	payloads := []string{
		"<script>alert(1)</script>",
		"<img src=x onerror=alert(1)>",
		"<svg onload=alert(1)>",
		"\"><script>alert(1)</script>",
		"'><script>alert(1)</script>",
		"javascript:alert(1)",
	}
	
	// Build result
	result := map[string]interface{}{
		"url":      url,
		"params":   params,
		"vulnerable": false,
		"details":  []map[string]interface{}{},
	}
	
	// In a real implementation, this would test each payload against each parameter
	// For now, we'll simulate the scan
	
	// Simulate finding a vulnerability
	if params["search"] != "" {
		result["vulnerable"] = true
		result["details"] = append(result["details"].([]map[string]interface{}), map[string]interface{}{
			"param":   "search",
			"payload": payloads[0],
			"response": "Payload reflected in response",
		})
	}
	
	return result, nil
}