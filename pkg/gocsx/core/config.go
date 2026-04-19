package core

// Config represents the configuration for Gocsx
type Config struct {
	// Theme configuration
	Theme *ThemeConfig

	// Platform-specific configuration
	Platform PlatformConfig

	// Breakpoints for responsive design
	Breakpoints map[string]int

	// Whether to enable dark mode
	DarkMode bool

	// Whether to enable RTL support
	RTL bool

	// Whether to enable animations
	Animations bool

	// Custom variants
	Variants map[string]VariantConfig

	// Prefix for all classes
	Prefix string
}

// ThemeConfig represents the theme configuration
type ThemeConfig struct {
	// Color palette
	Colors map[string]map[string]string

	// Spacing scale
	Spacing map[string]string

	// Typography scale
	Typography map[string]map[string]string

	// Border radius scale
	BorderRadius map[string]string

	// Shadow scale
	Shadows map[string]string

	// Z-index scale
	ZIndex map[string]int

	// Animation durations
	Durations map[string]string

	// Animation easings
	Easings map[string]string

	// Custom theme values
	Custom map[string]interface{}
}

// PlatformConfig represents platform-specific configuration
type PlatformConfig struct {
	// Target platform: "web", "mobile", "ar", "vr"
	Target string

	// Platform-specific overrides
	Overrides map[string]map[string]interface{}

	// Platform-specific features
	Features map[string]bool
}

// VariantConfig represents a custom variant configuration
type VariantConfig struct {
	// Variant name
	Name string

	// Variant description
	Description string

	// CSS property to modify
	Property string

	// Values for the variant
	Values map[string]string
}

// DefaultConfig returns the default configuration for Gocsx
func DefaultConfig() *Config {
	return &Config{
		Theme: &ThemeConfig{
			Colors: map[string]map[string]string{
				"primary": {
					"50":  "#f0f9ff",
					"100": "#e0f2fe",
					"200": "#bae6fd",
					"300": "#7dd3fc",
					"400": "#38bdf8",
					"500": "#0ea5e9",
					"600": "#0284c7",
					"700": "#0369a1",
					"800": "#075985",
					"900": "#0c4a6e",
					"950": "#082f49",
				},
				"neutral": {
					"50":  "#f9fafb",
					"100": "#f3f4f6",
					"200": "#e5e7eb",
					"300": "#d1d5db",
					"400": "#9ca3af",
					"500": "#6b7280",
					"600": "#4b5563",
					"700": "#374151",
					"800": "#1f2937",
					"900": "#111827",
					"950": "#030712",
				},
				// Add more color palettes here
			},
			Spacing: map[string]string{
				"0":   "0",
				"px":  "1px",
				"0.5": "0.125rem",
				"1":   "0.25rem",
				"1.5": "0.375rem",
				"2":   "0.5rem",
				"2.5": "0.625rem",
				"3":   "0.75rem",
				"3.5": "0.875rem",
				"4":   "1rem",
				"5":   "1.25rem",
				"6":   "1.5rem",
				"7":   "1.75rem",
				"8":   "2rem",
				"9":   "2.25rem",
				"10":  "2.5rem",
				"11":  "2.75rem",
				"12":  "3rem",
				"14":  "3.5rem",
				"16":  "4rem",
				"20":  "5rem",
				"24":  "6rem",
				"28":  "7rem",
				"32":  "8rem",
				"36":  "9rem",
				"40":  "10rem",
				"44":  "11rem",
				"48":  "12rem",
				"52":  "13rem",
				"56":  "14rem",
				"60":  "15rem",
				"64":  "16rem",
				"72":  "18rem",
				"80":  "20rem",
				"96":  "24rem",
			},
			Typography: map[string]map[string]string{
				"fontFamily": {
					"sans":  `ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"`,
					"serif": `ui-serif, Georgia, Cambria, "Times New Roman", Times, serif`,
					"mono":  `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace`,
				},
				"fontSize": {
					"xs":   "0.75rem",
					"sm":   "0.875rem",
					"base": "1rem",
					"lg":   "1.125rem",
					"xl":   "1.25rem",
					"2xl":  "1.5rem",
					"3xl":  "1.875rem",
					"4xl":  "2.25rem",
					"5xl":  "3rem",
					"6xl":  "3.75rem",
					"7xl":  "4.5rem",
					"8xl":  "6rem",
					"9xl":  "8rem",
				},
				"fontWeight": {
					"thin":       "100",
					"extralight": "200",
					"light":      "300",
					"normal":     "400",
					"medium":     "500",
					"semibold":   "600",
					"bold":       "700",
					"extrabold":  "800",
					"black":      "900",
				},
				"lineHeight": {
					"none":   "1",
					"tight":  "1.25",
					"snug":   "1.375",
					"normal": "1.5",
					"relaxed": "1.625",
					"loose":  "2",
				},
			},
			BorderRadius: map[string]string{
				"none":  "0",
				"sm":    "0.125rem",
				"md":    "0.375rem",
				"lg":    "0.5rem",
				"xl":    "0.75rem",
				"2xl":   "1rem",
				"3xl":   "1.5rem",
				"full":  "9999px",
			},
			Shadows: map[string]string{
				"sm":     "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
				"md":     "0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)",
				"lg":     "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)",
				"xl":     "0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)",
				"2xl":    "0 25px 50px -12px rgba(0, 0, 0, 0.25)",
				"inner":  "inset 0 2px 4px 0 rgba(0, 0, 0, 0.06)",
				"none":   "none",
			},
			ZIndex: map[string]int{
				"0":    0,
				"10":   10,
				"20":   20,
				"30":   30,
				"40":   40,
				"50":   50,
				"auto": 0,
			},
			Durations: map[string]string{
				"75":   "75ms",
				"100":  "100ms",
				"150":  "150ms",
				"200":  "200ms",
				"300":  "300ms",
				"500":  "500ms",
				"700":  "700ms",
				"1000": "1000ms",
			},
			Easings: map[string]string{
				"linear":      "linear",
				"in":          "cubic-bezier(0.4, 0, 1, 1)",
				"out":         "cubic-bezier(0, 0, 0.2, 1)",
				"in-out":      "cubic-bezier(0.4, 0, 0.2, 1)",
			},
			Custom: make(map[string]interface{}),
		},
		Platform: PlatformConfig{
			Target: "web",
			Overrides: make(map[string]map[string]interface{}),
			Features: map[string]bool{
				"darkMode":   true,
				"rtl":        false,
				"animations": true,
			},
		},
		Breakpoints: map[string]int{
			"sm": 640,
			"md": 768,
			"lg": 1024,
			"xl": 1280,
			"2xl": 1536,
		},
		DarkMode:   true,
		RTL:        false,
		Animations: true,
		Variants:   make(map[string]VariantConfig),
		Prefix:     "",
	}
}

// NewConfig creates a new configuration with custom options
func NewConfig(options ...func(*Config)) *Config {
	config := DefaultConfig()
	
	for _, option := range options {
		option(config)
	}
	
	return config
}

// WithTheme sets the theme configuration
func WithTheme(theme *ThemeConfig) func(*Config) {
	return func(c *Config) {
		c.Theme = theme
	}
}

// WithPlatform sets the platform configuration
func WithPlatform(platform PlatformConfig) func(*Config) {
	return func(c *Config) {
		c.Platform = platform
	}
}

// WithBreakpoints sets the breakpoints configuration
func WithBreakpoints(breakpoints map[string]int) func(*Config) {
	return func(c *Config) {
		c.Breakpoints = breakpoints
	}
}

// WithPrefix sets the class prefix
func WithPrefix(prefix string) func(*Config) {
	return func(c *Config) {
		c.Prefix = prefix
	}
}