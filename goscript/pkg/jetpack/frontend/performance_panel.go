package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/davidjeba/goscript/pkg/jetpack/core"
)

// PanelConfig represents the configuration for the performance panel
type PanelConfig struct {
	Enabled       bool     `json:"enabled"`
	Position      string   `json:"position"`
	Opacity       float64  `json:"opacity"`
	Theme         string   `json:"theme"`
	RefreshRate   int      `json:"refresh_rate"`
	MaxMetrics    int      `json:"max_metrics"`
	ShowCharts    bool     `json:"show_charts"`
	ShowAlerts    bool     `json:"show_alerts"`
	Draggable     bool     `json:"draggable"`
	Resizable     bool     `json:"resizable"`
	Collapsible   bool     `json:"collapsible"`
	DefaultMetrics []string `json:"default_metrics"`
}

// PerformancePanel represents the floating performance panel
type PerformancePanel struct {
	Jetpack       *core.Jetpack
	Config        PanelConfig
	Visible       bool
	Collapsed     bool
	Width         int
	Height        int
	SelectedTab   string
	SelectedMetrics []string
	LastUpdate    time.Time
}

// NewPerformancePanel creates a new performance panel
func NewPerformancePanel(jetpack *core.Jetpack) *PerformancePanel {
	return &PerformancePanel{
		Jetpack: jetpack,
		Config: PanelConfig{
			Enabled:       true,
			Position:      "bottom-right",
			Opacity:       0.8,
			Theme:         "dark",
			RefreshRate:   1000,
			MaxMetrics:    10,
			ShowCharts:    true,
			ShowAlerts:    true,
			Draggable:     true,
			Resizable:     true,
			Collapsible:   true,
			DefaultMetrics: []string{
				"fps",
				"memory_usage",
				"page_load",
				"first_contentful_paint",
				"largest_contentful_paint",
				"cumulative_layout_shift",
				"api_latency",
				"error_rate",
			},
		},
		Visible:       true,
		Collapsed:     false,
		Width:         350,
		Height:        500,
		SelectedTab:   "overview",
		SelectedMetrics: []string{
			"fps",
			"memory_usage",
			"page_load",
			"first_contentful_paint",
		},
		LastUpdate:    time.Now(),
	}
}

// Show shows the performance panel
func (pp *PerformancePanel) Show() {
	pp.Visible = true
}

// Hide hides the performance panel
func (pp *PerformancePanel) Hide() {
	pp.Visible = false
}

// Toggle toggles the visibility of the performance panel
func (pp *PerformancePanel) Toggle() {
	pp.Visible = !pp.Visible
}

// Collapse collapses the performance panel
func (pp *PerformancePanel) Collapse() {
	pp.Collapsed = true
}

// Expand expands the performance panel
func (pp *PerformancePanel) Expand() {
	pp.Collapsed = false
}

// ToggleCollapse toggles the collapsed state of the performance panel
func (pp *PerformancePanel) ToggleCollapse() {
	pp.Collapsed = !pp.Collapsed
}

// SetPosition sets the position of the performance panel
func (pp *PerformancePanel) SetPosition(position string) {
	pp.Config.Position = position
}

// SetOpacity sets the opacity of the performance panel
func (pp *PerformancePanel) SetOpacity(opacity float64) {
	if opacity < 0 {
		opacity = 0
	} else if opacity > 1 {
		opacity = 1
	}
	
	pp.Config.Opacity = opacity
}

// SetTheme sets the theme of the performance panel
func (pp *PerformancePanel) SetTheme(theme string) {
	pp.Config.Theme = theme
}

// SetRefreshRate sets the refresh rate of the performance panel
func (pp *PerformancePanel) SetRefreshRate(refreshRate int) {
	pp.Config.RefreshRate = refreshRate
}

// SelectTab selects a tab in the performance panel
func (pp *PerformancePanel) SelectTab(tab string) {
	pp.SelectedTab = tab
}

// SelectMetric adds a metric to the selected metrics
func (pp *PerformancePanel) SelectMetric(metric string) {
	// Check if metric is already selected
	for _, m := range pp.SelectedMetrics {
		if m == metric {
			return
		}
	}
	
	// Check if we've reached the maximum number of metrics
	if len(pp.SelectedMetrics) >= pp.Config.MaxMetrics {
		// Remove the first metric
		pp.SelectedMetrics = pp.SelectedMetrics[1:]
	}
	
	pp.SelectedMetrics = append(pp.SelectedMetrics, metric)
}

// UnselectMetric removes a metric from the selected metrics
func (pp *PerformancePanel) UnselectMetric(metric string) {
	for i, m := range pp.SelectedMetrics {
		if m == metric {
			pp.SelectedMetrics = append(pp.SelectedMetrics[:i], pp.SelectedMetrics[i+1:]...)
			return
		}
	}
}

// ResetMetrics resets the selected metrics to the default
func (pp *PerformancePanel) ResetMetrics() {
	pp.SelectedMetrics = make([]string, len(pp.Config.DefaultMetrics))
	copy(pp.SelectedMetrics, pp.Config.DefaultMetrics)
}

// GetPanelData gets the data for the performance panel
func (pp *PerformancePanel) GetPanelData() map[string]interface{} {
	pp.LastUpdate = time.Now()
	
	data := make(map[string]interface{})
	
	// Add basic info
	data["visible"] = pp.Visible
	data["collapsed"] = pp.Collapsed
	data["position"] = pp.Config.Position
	data["opacity"] = pp.Config.Opacity
	data["theme"] = pp.Config.Theme
	data["width"] = pp.Width
	data["height"] = pp.Height
	data["selected_tab"] = pp.SelectedTab
	data["last_update"] = pp.LastUpdate
	
	// Add selected metrics data
	selectedMetricsData := make([]map[string]interface{}, 0, len(pp.SelectedMetrics))
	for _, metricName := range pp.SelectedMetrics {
		metric, err := pp.Jetpack.GetMetric(metricName)
		if err != nil {
			continue
		}
		
		metricData := map[string]interface{}{
			"name":        metricName,
			"type":        metric.Type,
			"description": metric.Description,
			"unit":        metric.Unit,
			"alert":       metric.Alert,
		}
		
		// Get latest value
		latestValue, err := pp.Jetpack.GetMetricLatest(metricName)
		if err == nil {
			metricData["latest_value"] = latestValue
		}
		
		// Get average value
		avgValue, err := pp.Jetpack.GetMetricAverage(metricName)
		if err == nil {
			metricData["average_value"] = avgValue
		}
		
		selectedMetricsData = append(selectedMetricsData, metricData)
	}
	
	data["selected_metrics"] = selectedMetricsData
	
	// Add all available metrics for selection
	availableMetrics := make([]string, 0)
	for name := range pp.Jetpack.Metrics {
		availableMetrics = append(availableMetrics, name)
	}
	
	data["available_metrics"] = availableMetrics
	
	return data
}

// GenerateHTML generates the HTML for the performance panel
func (pp *PerformancePanel) GenerateHTML() (string, error) {
	if !pp.Visible {
		return "", nil
	}
	
	data := pp.GetPanelData()
	
	// Convert data to JSON for JavaScript
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	// Define the HTML template
	tmplStr := `
<div id="jetpack-performance-panel" class="jetpack-panel jetpack-theme-{{.theme}}" style="
	position: fixed;
	{{if eq .position "top-left"}}top: 10px; left: 10px;{{end}}
	{{if eq .position "top-right"}}top: 10px; right: 10px;{{end}}
	{{if eq .position "bottom-left"}}bottom: 10px; left: 10px;{{end}}
	{{if eq .position "bottom-right"}}bottom: 10px; right: 10px;{{end}}
	width: {{.width}}px;
	{{if .collapsed}}
		height: 30px;
		overflow: hidden;
	{{else}}
		height: {{.height}}px;
	{{end}}
	background-color: {{if eq .theme "dark"}}rgba(30, 30, 30, {{.opacity}}){{else}}rgba(240, 240, 240, {{.opacity}}){{end}};
	border-radius: 5px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
	font-family: 'Arial', sans-serif;
	font-size: 12px;
	color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
	z-index: 9999;
	transition: all 0.3s ease;
	overflow: auto;
">
	<div class="jetpack-panel-header" style="
		padding: 5px 10px;
		background-color: {{if eq .theme "dark"}}rgba(20, 20, 20, 0.8){{else}}rgba(220, 220, 220, 0.8){{end}};
		border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
		display: flex;
		justify-content: space-between;
		align-items: center;
		cursor: move;
	">
		<div class="jetpack-panel-title">
			Jetpack Performance Monitor
		</div>
		<div class="jetpack-panel-controls">
			<button onclick="jetpackToggleCollapse()" style="
				background: none;
				border: none;
				color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
				cursor: pointer;
				font-size: 14px;
				margin-right: 5px;
			">
				{{if .collapsed}}▼{{else}}▲{{end}}
			</button>
			<button onclick="jetpackHidePanel()" style="
				background: none;
				border: none;
				color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
				cursor: pointer;
				font-size: 14px;
			">
				✕
			</button>
		</div>
	</div>
	
	{{if not .collapsed}}
	<div class="jetpack-panel-tabs" style="
		display: flex;
		background-color: {{if eq .theme "dark"}}rgba(40, 40, 40, 0.8){{else}}rgba(230, 230, 230, 0.8){{end}};
		border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
	">
		<div class="jetpack-panel-tab {{if eq .selected_tab "overview"}}active{{end}}" 
			onclick="jetpackSelectTab('overview')" 
			style="
				padding: 5px 10px;
				cursor: pointer;
				{{if eq .selected_tab "overview"}}
					background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(250, 250, 250, 0.8){{end}};
					border-bottom: 2px solid #4285f4;
				{{end}}
			">
			Overview
		</div>
		<div class="jetpack-panel-tab {{if eq .selected_tab "metrics"}}active{{end}}" 
			onclick="jetpackSelectTab('metrics')" 
			style="
				padding: 5px 10px;
				cursor: pointer;
				{{if eq .selected_tab "metrics"}}
					background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(250, 250, 250, 0.8){{end}};
					border-bottom: 2px solid #4285f4;
				{{end}}
			">
			Metrics
		</div>
		<div class="jetpack-panel-tab {{if eq .selected_tab "lighthouse"}}active{{end}}" 
			onclick="jetpackSelectTab('lighthouse')" 
			style="
				padding: 5px 10px;
				cursor: pointer;
				{{if eq .selected_tab "lighthouse"}}
					background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(250, 250, 250, 0.8){{end}};
					border-bottom: 2px solid #4285f4;
				{{end}}
			">
			Lighthouse
		</div>
		<div class="jetpack-panel-tab {{if eq .selected_tab "settings"}}active{{end}}" 
			onclick="jetpackSelectTab('settings')" 
			style="
				padding: 5px 10px;
				cursor: pointer;
				{{if eq .selected_tab "settings"}}
					background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(250, 250, 250, 0.8){{end}};
					border-bottom: 2px solid #4285f4;
				{{end}}
			">
			Settings
		</div>
	</div>
	
	<div class="jetpack-panel-content" style="
		padding: 10px;
		overflow: auto;
		height: calc(100% - 70px);
	">
		{{if eq .selected_tab "overview"}}
			<div class="jetpack-panel-section">
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">Key Metrics</h3>
				<div class="jetpack-metrics-grid" style="
					display: grid;
					grid-template-columns: repeat(2, 1fr);
					gap: 10px;
				">
					{{range .selected_metrics}}
						<div class="jetpack-metric-card" style="
							background-color: {{if eq $.theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
							border-radius: 4px;
							padding: 8px;
							{{if .alert}}
								border-left: 3px solid #f44336;
							{{else}}
								border-left: 3px solid #4caf50;
							{{end}}
						">
							<div class="jetpack-metric-name" style="
								font-weight: bold;
								margin-bottom: 5px;
							">{{.name}}</div>
							<div class="jetpack-metric-value" style="
								font-size: 18px;
								font-weight: bold;
								color: {{if .alert}}#f44336{{else}}{{if eq $.theme "dark"}}#fff{{else}}#333{{end}}{{end}};
							">
								{{.latest_value}} {{.unit}}
							</div>
							<div class="jetpack-metric-avg" style="
								font-size: 10px;
								color: {{if eq $.theme "dark"}}#aaa{{else}}#777{{end}};
							">
								Avg: {{.average_value}} {{.unit}}
							</div>
						</div>
					{{end}}
				</div>
			</div>
			
			<div class="jetpack-panel-section" style="margin-top: 15px;">
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">Alerts</h3>
				<div class="jetpack-alerts-list" style="
					background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
					border-radius: 4px;
					padding: 8px;
				">
					{{$hasAlerts := false}}
					{{range .selected_metrics}}
						{{if .alert}}
							{{$hasAlerts = true}}
							<div class="jetpack-alert-item" style="
								margin-bottom: 5px;
								padding-bottom: 5px;
								border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#ddd{{end}};
							">
								<div style="display: flex; justify-content: space-between;">
									<span style="font-weight: bold; color: #f44336;">{{.name}}</span>
									<span>{{.latest_value}} {{.unit}}</span>
								</div>
								<div style="font-size: 10px; color: {{if eq $.theme "dark"}}#aaa{{else}}#777{{end}};">
									{{.description}}
								</div>
							</div>
						{{end}}
					{{end}}
					{{if not $hasAlerts}}
						<div style="text-align: center; padding: 10px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
							No alerts at this time
						</div>
					{{end}}
				</div>
			</div>
		{{end}}
		
		{{if eq .selected_tab "metrics"}}
			<div class="jetpack-panel-section">
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">All Metrics</h3>
				<div class="jetpack-metrics-filter" style="
					margin-bottom: 10px;
				">
					<input type="text" id="jetpack-metrics-filter" placeholder="Filter metrics..." style="
						width: 100%;
						padding: 5px;
						border: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
						border-radius: 3px;
						background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(255, 255, 255, 0.8){{end}};
						color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
					">
				</div>
				<div class="jetpack-metrics-list" style="
					max-height: 300px;
					overflow-y: auto;
					background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
					border-radius: 4px;
					padding: 8px;
				">
					{{range .available_metrics}}
						<div class="jetpack-metric-item" style="
							margin-bottom: 5px;
							padding: 5px;
							border-radius: 3px;
							cursor: pointer;
							background-color: {{if eq $.theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(255, 255, 255, 0.8){{end}};
						" onclick="jetpackToggleMetric('{{.}}')">
							<div style="display: flex; justify-content: space-between; align-items: center;">
								<span>{{.}}</span>
								<input type="checkbox" {{range $.selected_metrics}}{{if eq . $.}}checked{{end}}{{end}}>
							</div>
						</div>
					{{end}}
				</div>
			</div>
		{{end}}
		
		{{if eq .selected_tab "lighthouse"}}
			<div class="jetpack-panel-section">
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">Lighthouse Scores</h3>
				<div class="jetpack-lighthouse-scores" style="
					display: grid;
					grid-template-columns: repeat(2, 1fr);
					gap: 10px;
					margin-bottom: 15px;
				">
					<div class="jetpack-lighthouse-score" style="
						background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
						border-radius: 4px;
						padding: 8px;
						text-align: center;
					">
						<div style="font-weight: bold; margin-bottom: 5px;">Performance</div>
						<div style="
							font-size: 24px;
							font-weight: bold;
							color: #4caf50;
						">85</div>
					</div>
					<div class="jetpack-lighthouse-score" style="
						background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
						border-radius: 4px;
						padding: 8px;
						text-align: center;
					">
						<div style="font-weight: bold; margin-bottom: 5px;">Accessibility</div>
						<div style="
							font-size: 24px;
							font-weight: bold;
							color: #4caf50;
						">92</div>
					</div>
					<div class="jetpack-lighthouse-score" style="
						background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
						border-radius: 4px;
						padding: 8px;
						text-align: center;
					">
						<div style="font-weight: bold; margin-bottom: 5px;">Best Practices</div>
						<div style="
							font-size: 24px;
							font-weight: bold;
							color: #4caf50;
						">87</div>
					</div>
					<div class="jetpack-lighthouse-score" style="
						background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
						border-radius: 4px;
						padding: 8px;
						text-align: center;
					">
						<div style="font-weight: bold; margin-bottom: 5px;">SEO</div>
						<div style="
							font-size: 24px;
							font-weight: bold;
							color: #4caf50;
						">95</div>
					</div>
				</div>
				
				<button onclick="jetpackRunLighthouse()" style="
					background-color: #4285f4;
					color: white;
					border: none;
					border-radius: 3px;
					padding: 8px 12px;
					cursor: pointer;
					width: 100%;
					margin-bottom: 15px;
				">
					Run Lighthouse Audit
				</button>
				
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">Core Web Vitals</h3>
				<div class="jetpack-web-vitals" style="
					background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
					border-radius: 4px;
					padding: 8px;
				">
					<div class="jetpack-web-vital" style="
						margin-bottom: 8px;
						padding-bottom: 8px;
						border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
					">
						<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
							<span style="font-weight: bold;">LCP</span>
							<span style="color: #4caf50;">2.1s</span>
						</div>
						<div style="
							height: 6px;
							background-color: {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
							border-radius: 3px;
							overflow: hidden;
						">
							<div style="
								width: 70%;
								height: 100%;
								background-color: #4caf50;
								border-radius: 3px;
							"></div>
						</div>
					</div>
					<div class="jetpack-web-vital" style="
						margin-bottom: 8px;
						padding-bottom: 8px;
						border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
					">
						<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
							<span style="font-weight: bold;">FID</span>
							<span style="color: #4caf50;">15ms</span>
						</div>
						<div style="
							height: 6px;
							background-color: {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
							border-radius: 3px;
							overflow: hidden;
						">
							<div style="
								width: 90%;
								height: 100%;
								background-color: #4caf50;
								border-radius: 3px;
							"></div>
						</div>
					</div>
					<div class="jetpack-web-vital">
						<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
							<span style="font-weight: bold;">CLS</span>
							<span style="color: #4caf50;">0.05</span>
						</div>
						<div style="
							height: 6px;
							background-color: {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
							border-radius: 3px;
							overflow: hidden;
						">
							<div style="
								width: 85%;
								height: 100%;
								background-color: #4caf50;
								border-radius: 3px;
							"></div>
						</div>
					</div>
				</div>
			</div>
		{{end}}
		
		{{if eq .selected_tab "settings"}}
			<div class="jetpack-panel-section">
				<h3 style="margin: 0 0 10px 0; font-size: 14px;">Panel Settings</h3>
				<div class="jetpack-settings-form" style="
					background-color: {{if eq .theme "dark"}}rgba(50, 50, 50, 0.8){{else}}rgba(245, 245, 245, 0.8){{end}};
					border-radius: 4px;
					padding: 8px;
				">
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: block; margin-bottom: 5px;">Position</label>
						<select id="jetpack-position-setting" onchange="jetpackUpdateSetting('position', this.value)" style="
							width: 100%;
							padding: 5px;
							border: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
							border-radius: 3px;
							background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(255, 255, 255, 0.8){{end}};
							color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
						">
							<option value="top-left" {{if eq .position "top-left"}}selected{{end}}>Top Left</option>
							<option value="top-right" {{if eq .position "top-right"}}selected{{end}}>Top Right</option>
							<option value="bottom-left" {{if eq .position "bottom-left"}}selected{{end}}>Bottom Left</option>
							<option value="bottom-right" {{if eq .position "bottom-right"}}selected{{end}}>Bottom Right</option>
						</select>
					</div>
					
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: block; margin-bottom: 5px;">Theme</label>
						<select id="jetpack-theme-setting" onchange="jetpackUpdateSetting('theme', this.value)" style="
							width: 100%;
							padding: 5px;
							border: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
							border-radius: 3px;
							background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(255, 255, 255, 0.8){{end}};
							color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
						">
							<option value="dark" {{if eq .theme "dark"}}selected{{end}}>Dark</option>
							<option value="light" {{if eq .theme "light"}}selected{{end}}>Light</option>
						</select>
					</div>
					
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: block; margin-bottom: 5px;">Opacity: {{.opacity}}</label>
						<input type="range" id="jetpack-opacity-setting" min="0.1" max="1" step="0.1" value="{{.opacity}}" 
							oninput="jetpackUpdateSetting('opacity', this.value)" style="
							width: 100%;
						">
					</div>
					
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: block; margin-bottom: 5px;">Refresh Rate (ms)</label>
						<input type="number" id="jetpack-refresh-setting" min="100" max="10000" step="100" value="{{.Config.refresh_rate}}" 
							onchange="jetpackUpdateSetting('refresh_rate', this.value)" style="
							width: 100%;
							padding: 5px;
							border: 1px solid {{if eq .theme "dark"}}#444{{else}}#ccc{{end}};
							border-radius: 3px;
							background-color: {{if eq .theme "dark"}}rgba(60, 60, 60, 0.8){{else}}rgba(255, 255, 255, 0.8){{end}};
							color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
						">
					</div>
					
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: flex; align-items: center;">
							<input type="checkbox" id="jetpack-charts-setting" {{if .Config.show_charts}}checked{{end}} 
								onchange="jetpackUpdateSetting('show_charts', this.checked)" style="
								margin-right: 5px;
							">
							Show Charts
						</label>
					</div>
					
					<div class="jetpack-setting-item" style="margin-bottom: 10px;">
						<label style="display: flex; align-items: center;">
							<input type="checkbox" id="jetpack-alerts-setting" {{if .Config.show_alerts}}checked{{end}} 
								onchange="jetpackUpdateSetting('show_alerts', this.checked)" style="
								margin-right: 5px;
							">
							Show Alerts
						</label>
					</div>
					
					<button onclick="jetpackResetSettings()" style="
						background-color: #f44336;
						color: white;
						border: none;
						border-radius: 3px;
						padding: 8px 12px;
						cursor: pointer;
						width: 100%;
					">
						Reset Settings
					</button>
				</div>
			</div>
		{{end}}
	</div>
	{{end}}
</div>

<script>
	// Store panel data
	const jetpackPanelData = {{.dataJSON}};
	
	// Panel functions
	function jetpackHidePanel() {
		document.getElementById('jetpack-performance-panel').style.display = 'none';
	}
	
	function jetpackToggleCollapse() {
		const panel = document.getElementById('jetpack-performance-panel');
		const isCollapsed = panel.style.height === '30px';
		
		if (isCollapsed) {
			panel.style.height = '{{.height}}px';
			panel.style.overflow = 'auto';
		} else {
			panel.style.height = '30px';
			panel.style.overflow = 'hidden';
		}
		
		// Update button text
		const button = panel.querySelector('.jetpack-panel-controls button:first-child');
		button.textContent = isCollapsed ? '▲' : '▼';
	}
	
	function jetpackSelectTab(tab) {
		// This would be handled by the backend in a real implementation
		console.log('Selected tab:', tab);
	}
	
	function jetpackToggleMetric(metric) {
		// This would be handled by the backend in a real implementation
		console.log('Toggled metric:', metric);
	}
	
	function jetpackUpdateSetting(setting, value) {
		// This would be handled by the backend in a real implementation
		console.log('Updated setting:', setting, value);
	}
	
	function jetpackResetSettings() {
		// This would be handled by the backend in a real implementation
		console.log('Reset settings');
	}
	
	function jetpackRunLighthouse() {
		// This would be handled by the backend in a real implementation
		console.log('Running Lighthouse audit');
	}
	
	// Make panel draggable
	const panel = document.getElementById('jetpack-performance-panel');
	const header = panel.querySelector('.jetpack-panel-header');
	
	let isDragging = false;
	let offsetX, offsetY;
	
	header.addEventListener('mousedown', (e) => {
		isDragging = true;
		offsetX = e.clientX - panel.getBoundingClientRect().left;
		offsetY = e.clientY - panel.getBoundingClientRect().top;
	});
	
	document.addEventListener('mousemove', (e) => {
		if (!isDragging) return;
		
		const x = e.clientX - offsetX;
		const y = e.clientY - offsetY;
		
		panel.style.left = x + 'px';
		panel.style.top = y + 'px';
		panel.style.right = 'auto';
		panel.style.bottom = 'auto';
	});
	
	document.addEventListener('mouseup', () => {
		isDragging = false;
	});
	
	// Filter metrics
	const filterInput = document.getElementById('jetpack-metrics-filter');
	if (filterInput) {
		filterInput.addEventListener('input', (e) => {
			const filter = e.target.value.toLowerCase();
			const metricItems = document.querySelectorAll('.jetpack-metric-item');
			
			metricItems.forEach(item => {
				const metricName = item.textContent.trim().toLowerCase();
				if (metricName.includes(filter)) {
					item.style.display = 'block';
				} else {
					item.style.display = 'none';
				}
			});
		});
	}
</script>
`
	
	// Create template
	tmpl, err := template.New("panel").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	
	// Execute template
	var html strings.Builder
	err = tmpl.Execute(&html, map[string]interface{}{
		"theme":            pp.Config.Theme,
		"position":         pp.Config.Position,
		"opacity":          pp.Config.Opacity,
		"width":            pp.Width,
		"height":           pp.Height,
		"collapsed":        pp.Collapsed,
		"selected_tab":     pp.SelectedTab,
		"selected_metrics": data["selected_metrics"],
		"available_metrics": data["available_metrics"],
		"Config":           pp.Config,
		"dataJSON":         template.JS(string(dataJSON)),
	})
	if err != nil {
		return "", err
	}
	
	return html.String(), nil
}

// InjectIntoHTML injects the performance panel into an HTML page
func (pp *PerformancePanel) InjectIntoHTML(html string) (string, error) {
	if !pp.Visible {
		return html, nil
	}
	
	panelHTML, err := pp.GenerateHTML()
	if err != nil {
		return html, err
	}
	
	// Find the closing body tag
	bodyCloseIndex := strings.LastIndex(html, "</body>")
	if bodyCloseIndex == -1 {
		// If no body tag, append to the end
		return html + panelHTML, nil
	}
	
	// Insert the panel HTML before the closing body tag
	return html[:bodyCloseIndex] + panelHTML + html[bodyCloseIndex:], nil
}

// GenerateExtensionHTML generates the HTML for the Chrome extension
func (pp *PerformancePanel) GenerateExtensionHTML() (string, error) {
	data := pp.GetPanelData()
	
	// Convert data to JSON for JavaScript
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	// Define the HTML template for the extension
	tmplStr := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Jetpack Performance Monitor</title>
	<style>
		body {
			font-family: 'Arial', sans-serif;
			margin: 0;
			padding: 0;
			background-color: {{if eq .theme "dark"}}#1e1e1e{{else}}#f5f5f5{{end}};
			color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
		}
		
		.container {
			width: 800px;
			min-height: 600px;
			padding: 20px;
		}
		
		.header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 20px;
			padding-bottom: 10px;
			border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
		}
		
		.header h1 {
			margin: 0;
			font-size: 24px;
		}
		
		.tabs {
			display: flex;
			margin-bottom: 20px;
			border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};
		}
		
		.tab {
			padding: 10px 20px;
			cursor: pointer;
			border-bottom: 3px solid transparent;
		}
		
		.tab.active {
			border-bottom: 3px solid #4285f4;
			font-weight: bold;
		}
		
		.tab-content {
			display: none;
		}
		
		.tab-content.active {
			display: block;
		}
		
		.card {
			background-color: {{if eq .theme "dark"}}#2d2d2d{{else}}#fff{{end}};
			border-radius: 8px;
			padding: 15px;
			margin-bottom: 20px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		}
		
		.card h2 {
			margin-top: 0;
			margin-bottom: 15px;
			font-size: 18px;
			border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};
			padding-bottom: 10px;
		}
		
		.metrics-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 15px;
		}
		
		.metric-item {
			background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
			border-radius: 6px;
			padding: 12px;
			position: relative;
		}
		
		.metric-item.alert {
			border-left: 4px solid #f44336;
		}
		
		.metric-name {
			font-weight: bold;
			margin-bottom: 5px;
		}
		
		.metric-value {
			font-size: 24px;
			font-weight: bold;
			margin-bottom: 5px;
		}
		
		.metric-value.alert {
			color: #f44336;
		}
		
		.metric-description {
			font-size: 12px;
			color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};
		}
		
		.chart-container {
			height: 300px;
			margin-bottom: 20px;
		}
		
		.lighthouse-scores {
			display: flex;
			justify-content: space-between;
			margin-bottom: 20px;
		}
		
		.lighthouse-score {
			text-align: center;
			padding: 15px;
			background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
			border-radius: 6px;
			width: 22%;
		}
		
		.score-circle {
			width: 80px;
			height: 80px;
			border-radius: 50%;
			margin: 0 auto 10px;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 24px;
			font-weight: bold;
			color: white;
		}
		
		.score-label {
			font-weight: bold;
		}
		
		.good {
			background-color: #0cce6b;
		}
		
		.average {
			background-color: #ffa400;
		}
		
		.poor {
			background-color: #ff4e42;
		}
		
		.web-vitals {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 15px;
			margin-bottom: 20px;
		}
		
		.web-vital {
			background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
			border-radius: 6px;
			padding: 12px;
		}
		
		.web-vital-header {
			display: flex;
			justify-content: space-between;
			margin-bottom: 10px;
		}
		
		.web-vital-name {
			font-weight: bold;
		}
		
		.web-vital-value {
			font-weight: bold;
		}
		
		.web-vital-value.good {
			color: #0cce6b;
		}
		
		.web-vital-value.average {
			color: #ffa400;
		}
		
		.web-vital-value.poor {
			color: #ff4e42;
		}
		
		.progress-bar {
			height: 8px;
			background-color: {{if eq .theme "dark"}}#555{{else}}#eee{{end}};
			border-radius: 4px;
			overflow: hidden;
		}
		
		.progress-fill {
			height: 100%;
			border-radius: 4px;
		}
		
		.progress-fill.good {
			background-color: #0cce6b;
		}
		
		.progress-fill.average {
			background-color: #ffa400;
		}
		
		.progress-fill.poor {
			background-color: #ff4e42;
		}
		
		.settings-form {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 15px;
		}
		
		.form-group {
			margin-bottom: 15px;
		}
		
		.form-group label {
			display: block;
			margin-bottom: 5px;
			font-weight: bold;
		}
		
		.form-group input, .form-group select {
			width: 100%;
			padding: 8px;
			border: 1px solid {{if eq .theme "dark"}}#555{{else}}#ddd{{end}};
			border-radius: 4px;
			background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#fff{{end}};
			color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
		}
		
		.form-group input[type="checkbox"] {
			width: auto;
		}
		
		.button {
			padding: 10px 15px;
			border: none;
			border-radius: 4px;
			cursor: pointer;
			font-weight: bold;
		}
		
		.button-primary {
			background-color: #4285f4;
			color: white;
		}
		
		.button-secondary {
			background-color: {{if eq .theme "dark"}}#555{{else}}#eee{{end}};
			color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
		}
		
		.button-danger {
			background-color: #f44336;
			color: white;
		}
		
		.button-success {
			background-color: #0cce6b;
			color: white;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Jetpack Performance Monitor</h1>
			<div>
				<span>Last updated: {{.last_update}}</span>
			</div>
		</div>
		
		<div class="tabs">
			<div class="tab {{if eq .selected_tab "overview"}}active{{end}}" data-tab="overview">Overview</div>
			<div class="tab {{if eq .selected_tab "metrics"}}active{{end}}" data-tab="metrics">Metrics</div>
			<div class="tab {{if eq .selected_tab "lighthouse"}}active{{end}}" data-tab="lighthouse">Lighthouse</div>
			<div class="tab {{if eq .selected_tab "network"}}active{{end}}" data-tab="network">Network</div>
			<div class="tab {{if eq .selected_tab "settings"}}active{{end}}" data-tab="settings">Settings</div>
		</div>
		
		<div class="tab-content {{if eq .selected_tab "overview"}}active{{end}}" id="overview-tab">
			<div class="card">
				<h2>Key Metrics</h2>
				<div class="metrics-grid">
					{{range .selected_metrics}}
						<div class="metric-item {{if .alert}}alert{{end}}">
							<div class="metric-name">{{.name}}</div>
							<div class="metric-value {{if .alert}}alert{{end}}">{{.latest_value}} {{.unit}}</div>
							<div class="metric-description">{{.description}}</div>
						</div>
					{{end}}
				</div>
			</div>
			
			<div class="card">
				<h2>Performance Overview</h2>
				<div class="chart-container">
					<!-- Chart would be rendered here using a library like Chart.js -->
					<div style="text-align: center; padding: 100px 0; color: {{if eq $.theme "dark"}}#aaa{{else}}#777{{end}};">
						Performance chart will be displayed here
					</div>
				</div>
			</div>
			
			<div class="card">
				<h2>Alerts</h2>
				{{$hasAlerts := false}}
				{{range .selected_metrics}}
					{{if .alert}}
						{{$hasAlerts = true}}
						<div style="
							margin-bottom: 10px;
							padding: 10px;
							background-color: {{if eq $.theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
							border-left: 4px solid #f44336;
							border-radius: 4px;
						">
							<div style="font-weight: bold; margin-bottom: 5px;">{{.name}}: {{.latest_value}} {{.unit}}</div>
							<div style="font-size: 14px; color: {{if eq $.theme "dark"}}#aaa{{else}}#777{{end}};">{{.description}}</div>
						</div>
					{{end}}
				{{end}}
				{{if not $hasAlerts}}
					<div style="text-align: center; padding: 20px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
						No alerts at this time
					</div>
				{{end}}
			</div>
		</div>
		
		<div class="tab-content {{if eq .selected_tab "metrics"}}active{{end}}" id="metrics-tab">
			<div class="card">
				<h2>All Metrics</h2>
				<div style="margin-bottom: 15px;">
					<input type="text" id="metrics-filter" placeholder="Filter metrics..." style="
						width: 100%;
						padding: 8px;
						border: 1px solid {{if eq .theme "dark"}}#555{{else}}#ddd{{end}};
						border-radius: 4px;
						background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#fff{{end}};
						color: {{if eq .theme "dark"}}#fff{{else}}#333{{end}};
					">
				</div>
				
				<div style="
					max-height: 500px;
					overflow-y: auto;
					border: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};
					border-radius: 4px;
				">
					<table style="width: 100%; border-collapse: collapse;">
						<thead>
							<tr style="
								background-color: {{if eq .theme "dark"}}#2d2d2d{{else}}#f5f5f5{{end}};
								position: sticky;
								top: 0;
							">
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Name</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Value</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Unit</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Type</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Status</th>
							</tr>
						</thead>
						<tbody>
							{{range .selected_metrics}}
								<tr class="metric-row">
									<td style="padding: 10px; border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};">{{.name}}</td>
									<td style="padding: 10px; border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};">{{.latest_value}}</td>
									<td style="padding: 10px; border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};">{{.unit}}</td>
									<td style="padding: 10px; border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};">{{.type}}</td>
									<td style="padding: 10px; border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};">
										{{if .alert}}
											<span style="color: #f44336; font-weight: bold;">Alert</span>
										{{else}}
											<span style="color: #0cce6b;">Normal</span>
										{{end}}
									</td>
								</tr>
							{{end}}
						</tbody>
					</table>
				</div>
			</div>
			
			<div class="card">
				<h2>Metric Details</h2>
				<div style="
					padding: 20px;
					text-align: center;
					color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};
				">
					Select a metric from the table above to view details
				</div>
			</div>
		</div>
		
		<div class="tab-content {{if eq .selected_tab "lighthouse"}}active{{end}}" id="lighthouse-tab">
			<div class="card">
				<h2>Lighthouse Scores</h2>
				<div class="lighthouse-scores">
					<div class="lighthouse-score">
						<div class="score-circle good">85</div>
						<div class="score-label">Performance</div>
					</div>
					<div class="lighthouse-score">
						<div class="score-circle good">92</div>
						<div class="score-label">Accessibility</div>
					</div>
					<div class="lighthouse-score">
						<div class="score-circle good">87</div>
						<div class="score-label">Best Practices</div>
					</div>
					<div class="lighthouse-score">
						<div class="score-circle good">95</div>
						<div class="score-label">SEO</div>
					</div>
				</div>
				
				<div style="text-align: center; margin-bottom: 20px;">
					<button class="button button-primary">Run Lighthouse Audit</button>
				</div>
			</div>
			
			<div class="card">
				<h2>Core Web Vitals</h2>
				<div class="web-vitals">
					<div class="web-vital">
						<div class="web-vital-header">
							<div class="web-vital-name">LCP</div>
							<div class="web-vital-value good">2.1s</div>
						</div>
						<div class="progress-bar">
							<div class="progress-fill good" style="width: 70%;"></div>
						</div>
						<div style="font-size: 12px; margin-top: 5px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
							Largest Contentful Paint
						</div>
					</div>
					<div class="web-vital">
						<div class="web-vital-header">
							<div class="web-vital-name">FID</div>
							<div class="web-vital-value good">15ms</div>
						</div>
						<div class="progress-bar">
							<div class="progress-fill good" style="width: 90%;"></div>
						</div>
						<div style="font-size: 12px; margin-top: 5px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
							First Input Delay
						</div>
					</div>
					<div class="web-vital">
						<div class="web-vital-header">
							<div class="web-vital-name">CLS</div>
							<div class="web-vital-value good">0.05</div>
						</div>
						<div class="progress-bar">
							<div class="progress-fill good" style="width: 85%;"></div>
						</div>
						<div style="font-size: 12px; margin-top: 5px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
							Cumulative Layout Shift
						</div>
					</div>
				</div>
			</div>
			
			<div class="card">
				<h2>Opportunities</h2>
				<div style="
					margin-bottom: 10px;
					padding: 15px;
					background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
					border-radius: 6px;
				">
					<div style="font-weight: bold; margin-bottom: 5px;">Properly size images</div>
					<div style="font-size: 14px; margin-bottom: 10px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
						Serve images that are appropriately-sized to save cellular data and improve load time.
					</div>
					<div style="font-size: 14px; color: #ffa400;">Potential savings of 250KB</div>
				</div>
				
				<div style="
					margin-bottom: 10px;
					padding: 15px;
					background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
					border-radius: 6px;
				">
					<div style="font-weight: bold; margin-bottom: 5px;">Eliminate render-blocking resources</div>
					<div style="font-size: 14px; margin-bottom: 10px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">
						Resources are blocking the first paint of your page. Consider delivering critical JS/CSS inline and deferring all non-critical JS/styles.
					</div>
					<div style="font-size: 14px; color: #ffa400;">Potential savings of 500ms</div>
				</div>
			</div>
		</div>
		
		<div class="tab-content {{if eq .selected_tab "network"}}active{{end}}" id="network-tab">
			<div class="card">
				<h2>Network Requests</h2>
				<div style="
					max-height: 400px;
					overflow-y: auto;
					border: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};
					border-radius: 4px;
				">
					<table style="width: 100%; border-collapse: collapse;">
						<thead>
							<tr style="
								background-color: {{if eq .theme "dark"}}#2d2d2d{{else}}#f5f5f5{{end}};
								position: sticky;
								top: 0;
							">
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">URL</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Type</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Size</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Time</th>
								<th style="padding: 10px; text-align: left; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#ddd{{end}};">Status</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">https://example.com/main.js</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">script</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">125 KB</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">350 ms</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">200</td>
							</tr>
							<tr>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">https://example.com/styles.css</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">stylesheet</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">45 KB</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">120 ms</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">200</td>
							</tr>
							<tr>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">https://example.com/image.jpg</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">image</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">250 KB</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">180 ms</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">200</td>
							</tr>
							<tr>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">https://example.com/api/data</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">fetch</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">15 KB</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">220 ms</td>
								<td style="padding: 10px; border-bottom: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};">200</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			
			<div class="card">
				<h2>Network Summary</h2>
				<div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 15px;">
					<div style="
						background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
						border-radius: 6px;
						padding: 15px;
						text-align: center;
					">
						<div style="font-size: 12px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">Total Requests</div>
						<div style="font-size: 24px; font-weight: bold; margin: 10px 0;">24</div>
					</div>
					<div style="
						background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
						border-radius: 6px;
						padding: 15px;
						text-align: center;
					">
						<div style="font-size: 12px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">Total Size</div>
						<div style="font-size: 24px; font-weight: bold; margin: 10px 0;">1.2 MB</div>
					</div>
					<div style="
						background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
						border-radius: 6px;
						padding: 15px;
						text-align: center;
					">
						<div style="font-size: 12px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">Total Time</div>
						<div style="font-size: 24px; font-weight: bold; margin: 10px 0;">1.5s</div>
					</div>
					<div style="
						background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
						border-radius: 6px;
						padding: 15px;
						text-align: center;
					">
						<div style="font-size: 12px; color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};">Failed Requests</div>
						<div style="font-size: 24px; font-weight: bold; margin: 10px 0;">0</div>
					</div>
				</div>
			</div>
			
			<div class="card">
				<h2>Waterfall Chart</h2>
				<div style="
					height: 300px;
					background-color: {{if eq .theme "dark"}}#3d3d3d{{else}}#f9f9f9{{end}};
					border-radius: 6px;
					display: flex;
					align-items: center;
					justify-content: center;
					color: {{if eq .theme "dark"}}#aaa{{else}}#777{{end}};
				">
					Waterfall chart will be displayed here
				</div>
			</div>
		</div>
		
		<div class="tab-content {{if eq .selected_tab "settings"}}active{{end}}" id="settings-tab">
			<div class="card">
				<h2>Panel Settings</h2>
				<div class="settings-form">
					<div class="form-group">
						<label for="theme-setting">Theme</label>
						<select id="theme-setting">
							<option value="dark" {{if eq .theme "dark"}}selected{{end}}>Dark</option>
							<option value="light" {{if eq .theme "light"}}selected{{end}}>Light</option>
						</select>
					</div>
					
					<div class="form-group">
						<label for="position-setting">Position</label>
						<select id="position-setting">
							<option value="top-left" {{if eq .position "top-left"}}selected{{end}}>Top Left</option>
							<option value="top-right" {{if eq .position "top-right"}}selected{{end}}>Top Right</option>
							<option value="bottom-left" {{if eq .position "bottom-left"}}selected{{end}}>Bottom Left</option>
							<option value="bottom-right" {{if eq .position "bottom-right"}}selected{{end}}>Bottom Right</option>
						</select>
					</div>
					
					<div class="form-group">
						<label for="opacity-setting">Opacity: {{.opacity}}</label>
						<input type="range" id="opacity-setting" min="0.1" max="1" step="0.1" value="{{.opacity}}">
					</div>
					
					<div class="form-group">
						<label for="refresh-setting">Refresh Rate (ms)</label>
						<input type="number" id="refresh-setting" min="100" max="10000" step="100" value="{{.Config.refresh_rate}}">
					</div>
					
					<div class="form-group">
						<label>
							<input type="checkbox" id="charts-setting" {{if .Config.show_charts}}checked{{end}}>
							Show Charts
						</label>
					</div>
					
					<div class="form-group">
						<label>
							<input type="checkbox" id="alerts-setting" {{if .Config.show_alerts}}checked{{end}}>
							Show Alerts
						</label>
					</div>
				</div>
				
				<div style="margin-top: 20px; display: flex; gap: 10px;">
					<button class="button button-primary">Save Settings</button>
					<button class="button button-danger">Reset Settings</button>
				</div>
			</div>
			
			<div class="card">
				<h2>Metrics Settings</h2>
				<div style="margin-bottom: 15px;">
					<label style="display: block; margin-bottom: 5px; font-weight: bold;">Selected Metrics</label>
					<div style="
						max-height: 200px;
						overflow-y: auto;
						border: 1px solid {{if eq .theme "dark"}}#444{{else}}#eee{{end}};
						border-radius: 4px;
						padding: 10px;
					">
						{{range .selected_metrics}}
							<div style="
								display: flex;
								justify-content: space-between;
								align-items: center;
								padding: 5px 0;
								border-bottom: 1px solid {{if eq $.theme "dark"}}#444{{else}}#eee{{end}};
							">
								<span>{{.name}}</span>
								<button class="button button-secondary" style="padding: 3px 8px; font-size: 12px;">Remove</button>
							</div>
						{{end}}
					</div>
				</div>
				
				<div style="margin-top: 20px; display: flex; gap: 10px;">
					<button class="button button-primary">Add Metric</button>
					<button class="button button-secondary">Reset to Default</button>
				</div>
			</div>
		</div>
	</div>
	
	<script>
		// Store panel data
		const jetpackPanelData = {{.dataJSON}};
		
		// Tab switching
		document.querySelectorAll('.tab').forEach(tab => {
			tab.addEventListener('click', () => {
				const tabId = tab.getAttribute('data-tab');
				
				// Update active tab
				document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
				tab.classList.add('active');
				
				// Update active content
				document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
				document.getElementById(tabId + '-tab').classList.add('active');
			});
		});
		
		// Metrics filtering
		const metricsFilter = document.getElementById('metrics-filter');
		if (metricsFilter) {
			metricsFilter.addEventListener('input', (e) => {
				const filter = e.target.value.toLowerCase();
				const metricRows = document.querySelectorAll('.metric-row');
				
				metricRows.forEach(row => {
					const name = row.querySelector('td:first-child').textContent.toLowerCase();
					if (name.includes(filter)) {
						row.style.display = 'table-row';
					} else {
						row.style.display = 'none';
					}
				});
			});
		}
	</script>
</body>
</html>
`
	
	// Create template
	tmpl, err := template.New("extension").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	
	// Execute template
	var html strings.Builder
	err = tmpl.Execute(&html, map[string]interface{}{
		"theme":            pp.Config.Theme,
		"position":         pp.Config.Position,
		"opacity":          pp.Config.Opacity,
		"selected_tab":     pp.SelectedTab,
		"selected_metrics": data["selected_metrics"],
		"available_metrics": data["available_metrics"],
		"last_update":      pp.LastUpdate.Format(time.RFC1123),
		"Config":           pp.Config,
		"dataJSON":         template.JS(string(dataJSON)),
	})
	if err != nil {
		return "", err
	}
	
	return html.String(), nil
}