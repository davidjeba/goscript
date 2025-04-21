package mobile

import (
        "fmt"
        "strings"

        "github.com/davidjeba/goscript/pkg/gocsx/core"
)

// MobilePlatform implements the Platform interface for mobile devices
type MobilePlatform struct {
        // Configuration
        Config *core.Config

        // Device type (phone, tablet)
        DeviceType string

        // Operating system (iOS, Android)
        OS string

        // Screen dimensions
        ScreenWidth  int
        ScreenHeight int

        // Pixel density
        PixelDensity float64

        // Touch enabled
        TouchEnabled bool

        // Orientation (portrait, landscape)
        Orientation string
}

// NewMobilePlatform creates a new mobile platform adapter
func NewMobilePlatform(config *core.Config) *MobilePlatform {
        return &MobilePlatform{
                Config:       config,
                DeviceType:   "phone",
                OS:           "android",
                ScreenWidth:  375,
                ScreenHeight: 667,
                PixelDensity: 2.0,
                TouchEnabled: true,
                Orientation:  "portrait",
        }
}

// TransformCSS transforms CSS for mobile platforms
func (p *MobilePlatform) TransformCSS(css string) string {
        // Add mobile-specific viewport meta tag
        if strings.Contains(css, "<head>") {
                css = strings.Replace(css, "<head>", `<head>
                <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">`, 1)
        }

        // Add touch-specific styles
        css = css + `
        /* Mobile-specific styles */
        button, a, .clickable {
                cursor: pointer;
                -webkit-tap-highlight-color: rgba(0, 0, 0, 0);
                tap-highlight-color: rgba(0, 0, 0, 0);
                touch-action: manipulation;
        }
        
        /* Improve form elements on mobile */
        input, select, textarea {
                font-size: 16px; /* Prevent zoom on focus in iOS */
        }
        
        /* Orientation-specific styles */
        @media (orientation: portrait) {
                .landscape-only {
                        display: none !important;
                }
        }
        
        @media (orientation: landscape) {
                .portrait-only {
                        display: none !important;
                }
        }
        
        /* Device-specific styles */
        @media (max-width: 480px) {
                .tablet-only {
                        display: none !important;
                }
        }
        
        @media (min-width: 481px) {
                .phone-only {
                        display: none !important;
                }
        }
        
        /* iOS-specific styles */
        .ios-only {
                display: none;
        }
        
        /* Android-specific styles */
        .android-only {
                display: none;
        }
        `

        // Add OS-specific class
        if p.OS == "ios" {
                css = strings.Replace(css, ".ios-only { display: none; }", ".ios-only { display: block; }", 1)
        } else if p.OS == "android" {
                css = strings.Replace(css, ".android-only { display: none; }", ".android-only { display: block; }", 1)
        }

        return css
}

// GenerateUtilityClasses generates mobile-specific utility classes
func (p *MobilePlatform) GenerateUtilityClasses() map[string]string {
        utilities := make(map[string]string)

        // Touch utilities
        utilities["touch-none"] = "touch-action: none;"
        utilities["touch-pan-x"] = "touch-action: pan-x;"
        utilities["touch-pan-y"] = "touch-action: pan-y;"
        utilities["touch-manipulation"] = "touch-action: manipulation;"

        // Safe area utilities for notched devices
        utilities["safe-top"] = "padding-top: env(safe-area-inset-top);"
        utilities["safe-bottom"] = "padding-bottom: env(safe-area-inset-bottom);"
        utilities["safe-left"] = "padding-left: env(safe-area-inset-left);"
        utilities["safe-right"] = "padding-right: env(safe-area-inset-right);"
        utilities["safe-all"] = "padding: env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom) env(safe-area-inset-left);"

        // Mobile-specific display utilities
        utilities["hide-on-phone"] = "@media (max-width: 480px) { display: none !important; }"
        utilities["hide-on-tablet"] = "@media (min-width: 481px) { display: none !important; }"
        utilities["hide-on-portrait"] = "@media (orientation: portrait) { display: none !important; }"
        utilities["hide-on-landscape"] = "@media (orientation: landscape) { display: none !important; }"

        // Mobile-specific interaction utilities
        utilities["no-tap-highlight"] = "-webkit-tap-highlight-color: rgba(0, 0, 0, 0); tap-highlight-color: rgba(0, 0, 0, 0);"
        utilities["prevent-zoom"] = "font-size: 16px;"
        utilities["momentum-scroll"] = "-webkit-overflow-scrolling: touch; overflow-scrolling: touch;"

        return utilities
}

// DetectDevice detects the mobile device type and capabilities
func (p *MobilePlatform) DetectDevice(userAgent string) {
        // Detect iOS
        if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") || strings.Contains(userAgent, "iPod") {
                p.OS = "ios"
                
                // Detect iPad
                if strings.Contains(userAgent, "iPad") {
                        p.DeviceType = "tablet"
                        p.ScreenWidth = 768
                        p.ScreenHeight = 1024
                } else {
                        p.DeviceType = "phone"
                        p.ScreenWidth = 375
                        p.ScreenHeight = 667
                }
        } else if strings.Contains(userAgent, "Android") {
                p.OS = "android"
                
                // Detect Android tablet
                if strings.Contains(userAgent, "Android") && !strings.Contains(userAgent, "Mobile") {
                        p.DeviceType = "tablet"
                        p.ScreenWidth = 800
                        p.ScreenHeight = 1280
                } else {
                        p.DeviceType = "phone"
                        p.ScreenWidth = 360
                        p.ScreenHeight = 640
                }
        }

        // Detect orientation from dimensions
        if p.ScreenWidth > p.ScreenHeight {
                p.Orientation = "landscape"
        } else {
                p.Orientation = "portrait"
        }

        fmt.Printf("Detected mobile device: %s %s, %dx%d, %s orientation\n", 
                p.OS, p.DeviceType, p.ScreenWidth, p.ScreenHeight, p.Orientation)
}

// GetMediaQuery returns a media query for the current device
func (p *MobilePlatform) GetMediaQuery() string {
        if p.DeviceType == "phone" {
                return "@media (max-width: 480px)"
        } else if p.DeviceType == "tablet" {
                return "@media (min-width: 481px) and (max-width: 1024px)"
        }
        return ""
}

// GetTouchEvents returns touch event handlers for mobile interactions
func (p *MobilePlatform) GetTouchEvents() map[string]string {
        events := make(map[string]string)
        
        events["tap"] = "touchstart"
        events["doubleTap"] = "touchstart" // With custom detection
        events["longPress"] = "touchstart" // With custom detection
        events["swipeLeft"] = "touchmove"  // With custom detection
        events["swipeRight"] = "touchmove" // With custom detection
        events["swipeUp"] = "touchmove"    // With custom detection
        events["swipeDown"] = "touchmove"  // With custom detection
        events["pinch"] = "touchmove"      // With custom detection
        events["rotate"] = "touchmove"     // With custom detection
        
        return events
}

// Register registers the mobile platform with the core
func Register(g *core.Gocsx) {
        platform := NewMobilePlatform(g.Config)
        g.RegisterPlatform("mobile", platform)
}