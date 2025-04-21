package engine

import (
	"sync"
	"time"
)

// RenderingContext represents the type of rendering context
type RenderingContext string

const (
	// Context2D represents a 2D canvas rendering context
	Context2D RenderingContext = "2d"
	
	// Context3D represents a 3D WebGPU rendering context
	Context3D RenderingContext = "3d"
	
	// ContextHybrid represents a hybrid rendering context (both 2D and 3D)
	ContextHybrid RenderingContext = "hybrid"
)

// PerformanceLevel represents the engine performance level
type PerformanceLevel string

const (
	// PerformanceLow is for battery saving mode
	PerformanceLow PerformanceLevel = "low"
	
	// PerformanceMedium is the default performance level
	PerformanceMedium PerformanceLevel = "medium"
	
	// PerformanceHigh is for maximum performance
	PerformanceHigh PerformanceLevel = "high"
	
	// PerformanceAdaptive automatically adjusts based on device capabilities
	PerformanceAdaptive PerformanceLevel = "adaptive"
)

// EngineConfig represents the configuration for the engine
type EngineConfig struct {
	// Context is the rendering context type
	Context RenderingContext
	
	// TargetFPS is the target frames per second
	TargetFPS int
	
	// AutoDetect automatically detects the appropriate context
	AutoDetect bool
	
	// PerformanceLevel sets the performance level
	PerformanceLevel PerformanceLevel
	
	// EnableDebug enables debug mode
	EnableDebug bool
	
	// EnableStats enables performance statistics
	EnableStats bool
}

// DefaultEngineConfig returns the default engine configuration
func DefaultEngineConfig() *EngineConfig {
	return &EngineConfig{
		Context:         Context2D,
		TargetFPS:       60,
		AutoDetect:      true,
		PerformanceLevel: PerformanceAdaptive,
		EnableDebug:     false,
		EnableStats:     false,
	}
}

// Engine represents the core engine system
type Engine struct {
	// Configuration
	Config *EngineConfig
	
	// Current context
	currentContext RenderingContext
	
	// Performance metrics
	fps            float64
	frameTime      float64
	lastFrameTime  time.Time
	frameCount     int
	fpsUpdateTime  time.Time
	
	// Engine state
	running        bool
	paused         bool
	
	// Render loop
	renderCallback func(float64)
	
	// Performance throttling
	throttleLevel  float64
	
	// Mutex for thread safety
	mutex          sync.RWMutex
	
	// Stats
	stats          *EngineStats
}

// EngineStats represents engine performance statistics
type EngineStats struct {
	FPS            float64
	FrameTime      float64
	DrawCalls      int
	Triangles      int
	Textures       int
	ShaderSwitches int
	MemoryUsage    float64
}

// NewEngine creates a new engine instance
func NewEngine(config *EngineConfig) *Engine {
	if config == nil {
		config = DefaultEngineConfig()
	}
	
	engine := &Engine{
		Config:        config,
		currentContext: config.Context,
		lastFrameTime: time.Now(),
		fpsUpdateTime: time.Now(),
		throttleLevel: 1.0,
		stats:         &EngineStats{},
	}
	
	if config.AutoDetect {
		engine.detectContext()
	}
	
	return engine
}

// detectContext automatically detects the appropriate rendering context
func (e *Engine) detectContext() {
	// This would normally check for WebGPU support, 3D requirements, etc.
	// For now, we'll default to 2D
	e.currentContext = Context2D
}

// Start starts the engine
func (e *Engine) Start() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	if e.running {
		return
	}
	
	e.running = true
	e.paused = false
	e.lastFrameTime = time.Now()
	e.fpsUpdateTime = time.Now()
	
	// Start the render loop
	go e.renderLoop()
}

// Stop stops the engine
func (e *Engine) Stop() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.running = false
}

// Pause pauses the engine
func (e *Engine) Pause() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.paused = true
}

// Resume resumes the engine
func (e *Engine) Resume() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.paused = false
	e.lastFrameTime = time.Now()
}

// SetRenderCallback sets the render callback function
func (e *Engine) SetRenderCallback(callback func(float64)) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.renderCallback = callback
}

// SetContext sets the rendering context
func (e *Engine) SetContext(context RenderingContext) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.currentContext = context
}

// GetContext gets the current rendering context
func (e *Engine) GetContext() RenderingContext {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.currentContext
}

// SetPerformanceLevel sets the performance level
func (e *Engine) SetPerformanceLevel(level PerformanceLevel) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.Config.PerformanceLevel = level
	
	// Adjust throttling based on performance level
	switch level {
	case PerformanceLow:
		e.throttleLevel = 0.5
		e.Config.TargetFPS = 30
	case PerformanceMedium:
		e.throttleLevel = 0.75
		e.Config.TargetFPS = 60
	case PerformanceHigh:
		e.throttleLevel = 1.0
		e.Config.TargetFPS = 120
	case PerformanceAdaptive:
		// This would normally detect device capabilities
		e.throttleLevel = 0.75
		e.Config.TargetFPS = 60
	}
}

// GetStats gets the current engine statistics
func (e *Engine) GetStats() *EngineStats {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return &EngineStats{
		FPS:            e.fps,
		FrameTime:      e.frameTime,
		DrawCalls:      e.stats.DrawCalls,
		Triangles:      e.stats.Triangles,
		Textures:       e.stats.Textures,
		ShaderSwitches: e.stats.ShaderSwitches,
		MemoryUsage:    e.stats.MemoryUsage,
	}
}

// renderLoop is the main render loop
func (e *Engine) renderLoop() {
	for {
		e.mutex.RLock()
		running := e.running
		paused := e.paused
		callback := e.renderCallback
		targetFPS := e.Config.TargetFPS
		e.mutex.RUnlock()
		
		if !running {
			break
		}
		
		if paused || callback == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		
		// Calculate delta time
		now := time.Now()
		deltaTime := now.Sub(e.lastFrameTime).Seconds()
		
		// Update FPS counter
		e.frameCount++
		if now.Sub(e.fpsUpdateTime).Seconds() >= 1.0 {
			e.mutex.Lock()
			e.fps = float64(e.frameCount) / now.Sub(e.fpsUpdateTime).Seconds()
			e.frameCount = 0
			e.fpsUpdateTime = now
			e.mutex.Unlock()
		}
		
		// Call render callback
		callback(deltaTime)
		
		// Update frame time
		frameEnd := time.Now()
		e.mutex.Lock()
		e.frameTime = frameEnd.Sub(now).Seconds() * 1000 // in milliseconds
		e.lastFrameTime = now
		e.mutex.Unlock()
		
		// Throttle to target FPS
		targetFrameTime := 1.0 / float64(targetFPS)
		actualFrameTime := frameEnd.Sub(now).Seconds()
		if actualFrameTime < targetFrameTime {
			sleepTime := time.Duration((targetFrameTime - actualFrameTime) * float64(time.Second))
			time.Sleep(sleepTime)
		}
	}
}

// IsRunning checks if the engine is running
func (e *Engine) IsRunning() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.running
}

// IsPaused checks if the engine is paused
func (e *Engine) IsPaused() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.paused
}

// GetFPS gets the current FPS
func (e *Engine) GetFPS() float64 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.fps
}

// GetFrameTime gets the current frame time in milliseconds
func (e *Engine) GetFrameTime() float64 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.frameTime
}

// UpdateStats updates the engine statistics
func (e *Engine) UpdateStats(drawCalls, triangles, textures, shaderSwitches int, memoryUsage float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	e.stats.DrawCalls = drawCalls
	e.stats.Triangles = triangles
	e.stats.Textures = textures
	e.stats.ShaderSwitches = shaderSwitches
	e.stats.MemoryUsage = memoryUsage
}