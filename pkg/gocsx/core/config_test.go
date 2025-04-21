package core

import (
        "testing"
)

func TestNewConfig(t *testing.T) {
        config := NewConfig()
        
        // Test default values
        AssertNotNil(t, config)
        AssertNotNil(t, config.Theme)
        AssertNotNil(t, config.Breakpoints)
        AssertEqual(t, "light", config.Theme.Mode)
        AssertEqual(t, 4, len(config.Breakpoints))
}

func TestConfigSetTheme(t *testing.T) {
        config := NewConfig()
        
        // Test setting theme
        theme := &Theme{
                Mode: "dark",
                Colors: map[string]string{
                        "primary":   "#0f766e",
                        "secondary": "#7c3aed",
                        "accent":    "#f59e0b",
                        "background": "#1f2937",
                        "text":      "#f3f4f6",
                },
        }
        
        config.SetTheme(theme)
        
        AssertEqual(t, "dark", config.Theme.Mode)
        AssertEqual(t, "#0f766e", config.Theme.Colors["primary"])
        AssertEqual(t, "#7c3aed", config.Theme.Colors["secondary"])
        AssertEqual(t, "#f59e0b", config.Theme.Colors["accent"])
        AssertEqual(t, "#1f2937", config.Theme.Colors["background"])
        AssertEqual(t, "#f3f4f6", config.Theme.Colors["text"])
}

func TestConfigAddBreakpoint(t *testing.T) {
        config := NewConfig()
        
        // Test adding a breakpoint
        config.AddBreakpoint("ultrawide", 2560)
        
        AssertEqual(t, 5, len(config.Breakpoints))
        AssertEqual(t, 2560, config.Breakpoints["ultrawide"])
}

func TestConfigRemoveBreakpoint(t *testing.T) {
        config := NewConfig()
        
        // Test removing a breakpoint
        config.RemoveBreakpoint("md")
        
        AssertEqual(t, 3, len(config.Breakpoints))
        _, exists := config.Breakpoints["md"]
        AssertFalse(t, exists)
}

func TestConfigGetBreakpoint(t *testing.T) {
        config := NewConfig()
        
        // Test getting a breakpoint
        bp, exists := config.GetBreakpoint("lg")
        
        AssertTrue(t, exists)
        AssertEqual(t, 1024, bp)
        
        // Test getting a non-existent breakpoint
        bp, exists = config.GetBreakpoint("nonexistent")
        
        AssertFalse(t, exists)
        AssertEqual(t, 0, bp)
}

func TestConfigGetMediaQuery(t *testing.T) {
        config := NewConfig()
        
        // Test getting a media query
        query := config.GetMediaQuery("md")
        
        AssertEqual(t, "@media (min-width: 768px)", query)
        
        // Test getting a non-existent media query
        query = config.GetMediaQuery("nonexistent")
        
        AssertEqual(t, "", query)
}

func TestConfigSetPlatformConfig(t *testing.T) {
        config := NewConfig()
        
        // Test setting platform config
        platformConfig := map[string]interface{}{
                "mobile": map[string]interface{}{
                        "touchEnabled": true,
                        "deviceType":   "phone",
                },
                "web": map[string]interface{}{
                        "darkModeSupport": true,
                        "cssPrefix":       "gocsx",
                },
        }
        
        config.SetPlatformConfig(platformConfig)
        
        AssertNotNil(t, config.PlatformConfig)
        AssertEqual(t, platformConfig, config.PlatformConfig)
        
        // Test getting platform config
        mobileConfig, exists := config.GetPlatformConfig("mobile")
        
        AssertTrue(t, exists)
        AssertEqual(t, platformConfig["mobile"], mobileConfig)
        
        // Test getting a non-existent platform config
        nonexistentConfig, exists := config.GetPlatformConfig("nonexistent")
        
        AssertFalse(t, exists)
        AssertNil(t, nonexistentConfig)
}

func TestConfigSetVariants(t *testing.T) {
        config := NewConfig()
        
        // Test setting variants
        variants := map[string]map[string]string{
                "hover": {
                        "scale": "transform: scale(1.05);",
                        "glow":  "box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);",
                },
                "active": {
                        "pressed": "transform: scale(0.95);",
                        "outline": "outline: 2px solid #0f766e;",
                },
        }
        
        config.SetVariants(variants)
        
        AssertNotNil(t, config.Variants)
        AssertEqual(t, variants, config.Variants)
        
        // Test getting variants
        hoverVariants, exists := config.GetVariant("hover")
        
        AssertTrue(t, exists)
        AssertEqual(t, variants["hover"], hoverVariants)
        
        // Test getting a non-existent variant
        nonexistentVariant, exists := config.GetVariant("nonexistent")
        
        AssertFalse(t, exists)
        AssertNil(t, nonexistentVariant)
}

func TestConfigToJSON(t *testing.T) {
        config := NewConfig()
        
        // Test converting to JSON
        json, err := config.ToJSON()
        
        AssertNil(t, err)
        AssertNotNil(t, json)
        AssertContains(t, json, "\"theme\":")
        AssertContains(t, json, "\"breakpoints\":")
}

func TestConfigFromJSON(t *testing.T) {
        // Test creating from JSON
        json := `{
                "theme": {
                        "mode": "dark",
                        "colors": {
                                "primary": "#0f766e",
                                "secondary": "#7c3aed",
                                "accent": "#f59e0b",
                                "background": "#1f2937",
                                "text": "#f3f4f6"
                        }
                },
                "breakpoints": {
                        "sm": 640,
                        "md": 768,
                        "lg": 1024,
                        "xl": 1280
                }
        }`
        
        config, err := ConfigFromJSON(json)
        
        AssertNil(t, err)
        AssertNotNil(t, config)
        AssertEqual(t, "dark", config.Theme.Mode)
        AssertEqual(t, "#0f766e", config.Theme.Colors["primary"])
        AssertEqual(t, 4, len(config.Breakpoints))
        AssertEqual(t, 640, config.Breakpoints["sm"])
}

func TestConfigClone(t *testing.T) {
        config := NewConfig()
        
        // Modify the config
        config.Theme.Mode = "dark"
        config.Theme.Colors["primary"] = "#0f766e"
        config.AddBreakpoint("ultrawide", 2560)
        
        // Clone the config
        clone := config.Clone()
        
        // Verify the clone has the same values
        AssertEqual(t, "dark", clone.Theme.Mode)
        AssertEqual(t, "#0f766e", clone.Theme.Colors["primary"])
        AssertEqual(t, 5, len(clone.Breakpoints))
        AssertEqual(t, 2560, clone.Breakpoints["ultrawide"])
        
        // Modify the clone and verify the original is unchanged
        clone.Theme.Mode = "light"
        clone.Theme.Colors["primary"] = "#047857"
        clone.RemoveBreakpoint("ultrawide")
        
        AssertEqual(t, "dark", config.Theme.Mode)
        AssertEqual(t, "#0f766e", config.Theme.Colors["primary"])
        AssertEqual(t, 5, len(config.Breakpoints))
        AssertEqual(t, 2560, config.Breakpoints["ultrawide"])
}