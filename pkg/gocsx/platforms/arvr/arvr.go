package arvr

import (
        "fmt"
        "strings"

        "github.com/davidjeba/goscript/pkg/gocsx/core"
)

// ARVRPlatform implements the Platform interface for AR/VR devices
type ARVRPlatform struct {
        // Configuration
        Config *core.Config

        // Device type (ar, vr, mr)
        DeviceType string

        // Headset type (quest, hololens, etc.)
        HeadsetType string

        // Controller type (6dof, 3dof, hand-tracking)
        ControllerType string

        // Display resolution
        ResolutionWidth  int
        ResolutionHeight int

        // Field of view
        FOV float64

        // Degrees of freedom (3 or 6)
        DOF int

        // Room scale support
        RoomScaleSupported bool

        // Hand tracking support
        HandTrackingSupported bool

        // Eye tracking support
        EyeTrackingSupported bool
}

// NewARVRPlatform creates a new AR/VR platform adapter
func NewARVRPlatform(config *core.Config) *ARVRPlatform {
        return &ARVRPlatform{
                Config:                config,
                DeviceType:            "vr",
                HeadsetType:           "generic",
                ControllerType:        "6dof",
                ResolutionWidth:       2880,
                ResolutionHeight:      1600,
                FOV:                   110.0,
                DOF:                   6,
                RoomScaleSupported:    true,
                HandTrackingSupported: false,
                EyeTrackingSupported:  false,
        }
}

// TransformCSS transforms CSS for AR/VR platforms
func (p *ARVRPlatform) TransformCSS(css string) string {
        // Add AR/VR-specific meta tags
        if strings.Contains(css, "<head>") {
                metaTags := `<head>
                <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
                <meta name="apple-mobile-web-app-capable" content="yes">
                <meta name="webxr-device" content="true">`
                
                css = strings.Replace(css, "<head>", metaTags, 1)
        }

        // Add AR/VR-specific styles
        css = css + `
        /* AR/VR-specific styles */
        body {
                background-color: transparent;
                overflow: hidden;
                user-select: none;
                -webkit-user-select: none;
        }
        
        /* Cursor styles for VR */
        .vr-cursor {
                position: fixed;
                width: 10px;
                height: 10px;
                border-radius: 50%;
                background-color: white;
                box-shadow: 0 0 10px rgba(255, 255, 255, 0.8);
                pointer-events: none;
                z-index: 9999;
                transform: translate(-50%, -50%);
        }
        
        /* AR-specific styles */
        .ar-only {
                display: none;
        }
        
        /* VR-specific styles */
        .vr-only {
                display: none;
        }
        
        /* MR-specific styles */
        .mr-only {
                display: none;
        }
        
        /* 3DOF controller styles */
        .dof3-only {
                display: none;
        }
        
        /* 6DOF controller styles */
        .dof6-only {
                display: none;
        }
        
        /* Hand tracking styles */
        .hand-tracking-only {
                display: none;
        }
        `

        // Add device-specific class
        if p.DeviceType == "ar" {
                css = strings.Replace(css, ".ar-only { display: none; }", ".ar-only { display: block; }", 1)
        } else if p.DeviceType == "vr" {
                css = strings.Replace(css, ".vr-only { display: none; }", ".vr-only { display: block; }", 1)
        } else if p.DeviceType == "mr" {
                css = strings.Replace(css, ".mr-only { display: none; }", ".mr-only { display: block; }", 1)
        }

        // Add controller-specific class
        if p.ControllerType == "3dof" {
                css = strings.Replace(css, ".dof3-only { display: none; }", ".dof3-only { display: block; }", 1)
        } else if p.ControllerType == "6dof" {
                css = strings.Replace(css, ".dof6-only { display: none; }", ".dof6-only { display: block; }", 1)
        } else if p.ControllerType == "hand-tracking" {
                css = strings.Replace(css, ".hand-tracking-only { display: none; }", ".hand-tracking-only { display: block; }", 1)
        }

        return css
}

// GenerateUtilityClasses generates AR/VR-specific utility classes
func (p *ARVRPlatform) GenerateUtilityClasses() map[string]string {
        utilities := make(map[string]string)

        // Spatial positioning utilities
        utilities["spatial-fixed"] = "position: fixed; transform-style: preserve-3d;"
        utilities["spatial-relative"] = "position: relative; transform-style: preserve-3d;"
        utilities["spatial-absolute"] = "position: absolute; transform-style: preserve-3d;"

        // Depth utilities
        utilities["depth-near"] = "z-index: 100; transform: translateZ(0.5m);"
        utilities["depth-mid"] = "z-index: 0; transform: translateZ(0m);"
        utilities["depth-far"] = "z-index: -100; transform: translateZ(-0.5m);"

        // Gaze interaction utilities
        utilities["gaze-target"] = "cursor: pointer;"
        utilities["gaze-hover"] = "transition: all 0.3s ease-in-out;"
        utilities["gaze-active"] = "transform: scale(1.1);"

        // Controller interaction utilities
        utilities["controller-target"] = "cursor: pointer;"
        utilities["controller-hover"] = "transition: all 0.2s ease-in-out;"
        utilities["controller-active"] = "transform: scale(1.05);"

        // Hand interaction utilities
        utilities["hand-target"] = "cursor: grab;"
        utilities["hand-hover"] = "transition: all 0.2s ease-in-out; box-shadow: 0 0 15px rgba(255, 255, 255, 0.5);"
        utilities["hand-active"] = "cursor: grabbing; transform: scale(1.05);"

        // Visibility utilities
        utilities["visible-in-ar"] = ".ar-only { display: block; }"
        utilities["visible-in-vr"] = ".vr-only { display: block; }"
        utilities["visible-in-mr"] = ".mr-only { display: block; }"

        // Billboarding (always face user)
        utilities["billboard"] = "transform-style: preserve-3d; transform: rotateY(var(--user-rotation-y));"

        return utilities
}

// DetectDevice detects the AR/VR device type and capabilities
func (p *ARVRPlatform) DetectDevice(userAgent string, xrSupport bool) {
        // Detect device type based on user agent and XR support
        if strings.Contains(userAgent, "OculusBrowser") || strings.Contains(userAgent, "Quest") {
                p.DeviceType = "vr"
                p.HeadsetType = "quest"
                p.ControllerType = "6dof"
                p.ResolutionWidth = 3664
                p.ResolutionHeight = 1920
                p.FOV = 104.0
                p.DOF = 6
                p.RoomScaleSupported = true
                p.HandTrackingSupported = true
        } else if strings.Contains(userAgent, "HoloLens") {
                p.DeviceType = "mr"
                p.HeadsetType = "hololens"
                p.ControllerType = "hand-tracking"
                p.ResolutionWidth = 2048
                p.ResolutionHeight = 1080
                p.FOV = 52.0
                p.DOF = 6
                p.RoomScaleSupported = true
                p.HandTrackingSupported = true
                p.EyeTrackingSupported = true
        } else if strings.Contains(userAgent, "iPhone") && xrSupport {
                p.DeviceType = "ar"
                p.HeadsetType = "arkit"
                p.ControllerType = "3dof"
                p.ResolutionWidth = 1920
                p.ResolutionHeight = 1080
                p.FOV = 70.0
                p.DOF = 6
                p.RoomScaleSupported = true
                p.HandTrackingSupported = false
        } else if strings.Contains(userAgent, "Android") && xrSupport {
                p.DeviceType = "ar"
                p.HeadsetType = "arcore"
                p.ControllerType = "3dof"
                p.ResolutionWidth = 1920
                p.ResolutionHeight = 1080
                p.FOV = 70.0
                p.DOF = 6
                p.RoomScaleSupported = true
                p.HandTrackingSupported = false
        } else {
                // Generic VR fallback
                p.DeviceType = "vr"
                p.HeadsetType = "generic"
                p.ControllerType = "6dof"
                p.ResolutionWidth = 2880
                p.ResolutionHeight = 1600
                p.FOV = 110.0
                p.DOF = 6
                p.RoomScaleSupported = true
                p.HandTrackingSupported = false
        }

        fmt.Printf("Detected AR/VR device: %s %s, %s controller, %dx%d resolution, %.1f° FOV\n",
                p.DeviceType, p.HeadsetType, p.ControllerType, p.ResolutionWidth, p.ResolutionHeight, p.FOV)
}

// GetInteractionEvents returns event handlers for AR/VR interactions
func (p *ARVRPlatform) GetInteractionEvents() map[string]string {
        events := make(map[string]string)
        
        if p.DeviceType == "vr" {
                if p.ControllerType == "6dof" || p.ControllerType == "3dof" {
                        events["select"] = "xrselect"
                        events["squeeze"] = "xrsqueeze"
                        events["hover"] = "xrhover"
                } else if p.ControllerType == "hand-tracking" {
                        events["pinch"] = "xrpinch"
                        events["grab"] = "xrgrab"
                        events["hover"] = "xrhover"
                }
        } else if p.DeviceType == "ar" || p.DeviceType == "mr" {
                events["tap"] = "xrtap"
                events["place"] = "xrplace"
                events["hover"] = "xrhover"
        }
        
        // Common events
        events["enter"] = "xrenter"
        events["exit"] = "xrexit"
        
        return events
}

// GetSpatialProperties returns spatial properties for AR/VR positioning
func (p *ARVRPlatform) GetSpatialProperties() map[string]float64 {
        properties := make(map[string]float64)
        
        properties["defaultDistance"] = 2.0 // 2 meters
        properties["minDistance"] = 0.5     // 0.5 meters
        properties["maxDistance"] = 100.0   // 100 meters
        properties["interactionDistance"] = 1.5 // 1.5 meters
        properties["defaultHeight"] = 1.6   // 1.6 meters (average eye height)
        properties["defaultScale"] = 1.0    // 1:1 scale
        
        return properties
}

// Register registers the AR/VR platform with the core
func Register(g *core.Gocsx) {
        platform := NewARVRPlatform(g.Config)
        g.RegisterPlatform("arvr", platform)
}