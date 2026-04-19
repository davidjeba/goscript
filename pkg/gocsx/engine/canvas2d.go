package engine

import (
	"fmt"
	"sync"
)

// Canvas2D represents a 2D canvas
type Canvas2D struct {
	// Canvas ID
	ID string
	
	// Canvas width
	Width int
	
	// Canvas height
	Height int
	
	// Canvas context
	Context *Canvas2DContext
	
	// Engine
	Engine *Engine
	
	// Render callback
	RenderCallback func(context *Canvas2DContext, deltaTime float64)
	
	// Mutex for thread safety
	mutex sync.RWMutex
}

// Canvas2DContext represents a 2D canvas context
type Canvas2DContext struct {
	// Fill style
	FillStyle string
	
	// Stroke style
	StrokeStyle string
	
	// Line width
	LineWidth float64
	
	// Line cap
	LineCap string
	
	// Line join
	LineJoin string
	
	// Miter limit
	MiterLimit float64
	
	// Global alpha
	GlobalAlpha float64
	
	// Global composite operation
	GlobalCompositeOperation string
	
	// Font
	Font string
	
	// Text align
	TextAlign string
	
	// Text baseline
	TextBaseline string
	
	// Shadow color
	ShadowColor string
	
	// Shadow blur
	ShadowBlur float64
	
	// Shadow offset X
	ShadowOffsetX float64
	
	// Shadow offset Y
	ShadowOffsetY float64
	
	// Transform matrix
	Transform [6]float64
	
	// Clip region
	ClipRegion []Path2D
	
	// Stats
	Stats *Canvas2DStats
}

// Canvas2DStats represents 2D canvas statistics
type Canvas2DStats struct {
	// Draw calls
	DrawCalls int
	
	// Fill calls
	FillCalls int
	
	// Stroke calls
	StrokeCalls int
	
	// Clear calls
	ClearCalls int
	
	// Text calls
	TextCalls int
	
	// Image calls
	ImageCalls int
	
	// Path calls
	PathCalls int
	
	// Transform calls
	TransformCalls int
	
	// Memory usage
	MemoryUsage float64
}

// Path2D represents a 2D path
type Path2D struct {
	// Path ID
	ID string
	
	// Path commands
	Commands []Path2DCommand
}

// Path2DCommand represents a 2D path command
type Path2DCommand struct {
	// Command type
	Type string
	
	// Command parameters
	Params []float64
}

// NewCanvas2D creates a new 2D canvas
func NewCanvas2D(id string, width, height int, engine *Engine) *Canvas2D {
	// Create a context
	context := &Canvas2DContext{
		FillStyle:               "#000000",
		StrokeStyle:             "#000000",
		LineWidth:               1,
		LineCap:                 "butt",
		LineJoin:                "miter",
		MiterLimit:              10,
		GlobalAlpha:             1,
		GlobalCompositeOperation: "source-over",
		Font:                    "10px sans-serif",
		TextAlign:               "start",
		TextBaseline:            "alphabetic",
		ShadowColor:             "rgba(0, 0, 0, 0)",
		ShadowBlur:              0,
		ShadowOffsetX:           0,
		ShadowOffsetY:           0,
		Transform:               [6]float64{1, 0, 0, 1, 0, 0},
		ClipRegion:              []Path2D{},
		Stats:                   &Canvas2DStats{},
	}
	
	// Create a canvas
	canvas := &Canvas2D{
		ID:      id,
		Width:   width,
		Height:  height,
		Context: context,
		Engine:  engine,
	}
	
	// Set render callback
	engine.SetRenderCallback(canvas.Render)
	
	return canvas
}

// Render renders the canvas
func (c *Canvas2D) Render(deltaTime float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Reset stats
	c.Context.Stats = &Canvas2DStats{}
	
	// Call render callback
	if c.RenderCallback != nil {
		c.RenderCallback(c.Context, deltaTime)
	}
	
	// Update engine stats
	c.Engine.UpdateStats(
		c.Context.Stats.DrawCalls,
		0, // No triangles in 2D
		c.Context.Stats.ImageCalls,
		0, // No shader switches in 2D
		c.Context.Stats.MemoryUsage,
	)
}

// SetRenderCallback sets the render callback
func (c *Canvas2D) SetRenderCallback(callback func(context *Canvas2DContext, deltaTime float64)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.RenderCallback = callback
}

// SetSize sets the canvas size
func (c *Canvas2D) SetSize(width, height int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.Width = width
	c.Height = height
}

// GetContext gets the canvas context
func (c *Canvas2D) GetContext() *Canvas2DContext {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return c.Context
}

// GetStats gets the canvas statistics
func (c *Canvas2D) GetStats() *Canvas2DStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return c.Context.Stats
}

// ClearRect clears a rectangle
func (ctx *Canvas2DContext) ClearRect(x, y, width, height float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.ClearCalls++
	
	// This would normally clear the rectangle on the canvas
	// For now, we'll just update the stats
}

// FillRect fills a rectangle
func (ctx *Canvas2DContext) FillRect(x, y, width, height float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.FillCalls++
	
	// This would normally fill the rectangle on the canvas
	// For now, we'll just update the stats
}

// StrokeRect strokes a rectangle
func (ctx *Canvas2DContext) StrokeRect(x, y, width, height float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.StrokeCalls++
	
	// This would normally stroke the rectangle on the canvas
	// For now, we'll just update the stats
}

// FillText fills text
func (ctx *Canvas2DContext) FillText(text string, x, y float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.TextCalls++
	ctx.Stats.FillCalls++
	
	// This would normally fill the text on the canvas
	// For now, we'll just update the stats
}

// StrokeText strokes text
func (ctx *Canvas2DContext) StrokeText(text string, x, y float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.TextCalls++
	ctx.Stats.StrokeCalls++
	
	// This would normally stroke the text on the canvas
	// For now, we'll just update the stats
}

// MeasureText measures text
func (ctx *Canvas2DContext) MeasureText(text string) float64 {
	// This would normally measure the text
	// For now, we'll just return a dummy value
	return float64(len(text) * 10)
}

// BeginPath begins a path
func (ctx *Canvas2DContext) BeginPath() {
	ctx.Stats.PathCalls++
	
	// This would normally begin a path on the canvas
	// For now, we'll just update the stats
}

// ClosePath closes a path
func (ctx *Canvas2DContext) ClosePath() {
	ctx.Stats.PathCalls++
	
	// This would normally close a path on the canvas
	// For now, we'll just update the stats
}

// MoveTo moves to a point
func (ctx *Canvas2DContext) MoveTo(x, y float64) {
	ctx.Stats.PathCalls++
	
	// This would normally move to a point on the canvas
	// For now, we'll just update the stats
}

// LineTo draws a line to a point
func (ctx *Canvas2DContext) LineTo(x, y float64) {
	ctx.Stats.PathCalls++
	
	// This would normally draw a line to a point on the canvas
	// For now, we'll just update the stats
}

// BezierCurveTo draws a bezier curve
func (ctx *Canvas2DContext) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	ctx.Stats.PathCalls++
	
	// This would normally draw a bezier curve on the canvas
	// For now, we'll just update the stats
}

// QuadraticCurveTo draws a quadratic curve
func (ctx *Canvas2DContext) QuadraticCurveTo(cpx, cpy, x, y float64) {
	ctx.Stats.PathCalls++
	
	// This would normally draw a quadratic curve on the canvas
	// For now, we'll just update the stats
}

// Arc draws an arc
func (ctx *Canvas2DContext) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise bool) {
	ctx.Stats.PathCalls++
	
	// This would normally draw an arc on the canvas
	// For now, we'll just update the stats
}

// ArcTo draws an arc to a point
func (ctx *Canvas2DContext) ArcTo(x1, y1, x2, y2, radius float64) {
	ctx.Stats.PathCalls++
	
	// This would normally draw an arc to a point on the canvas
	// For now, we'll just update the stats
}

// Rect adds a rectangle to the path
func (ctx *Canvas2DContext) Rect(x, y, width, height float64) {
	ctx.Stats.PathCalls++
	
	// This would normally add a rectangle to the path on the canvas
	// For now, we'll just update the stats
}

// Fill fills the current path
func (ctx *Canvas2DContext) Fill() {
	ctx.Stats.DrawCalls++
	ctx.Stats.FillCalls++
	
	// This would normally fill the current path on the canvas
	// For now, we'll just update the stats
}

// Stroke strokes the current path
func (ctx *Canvas2DContext) Stroke() {
	ctx.Stats.DrawCalls++
	ctx.Stats.StrokeCalls++
	
	// This would normally stroke the current path on the canvas
	// For now, we'll just update the stats
}

// Clip clips the current path
func (ctx *Canvas2DContext) Clip() {
	ctx.Stats.PathCalls++
	
	// This would normally clip the current path on the canvas
	// For now, we'll just update the stats
}

// IsPointInPath checks if a point is in the current path
func (ctx *Canvas2DContext) IsPointInPath(x, y float64) bool {
	ctx.Stats.PathCalls++
	
	// This would normally check if a point is in the current path on the canvas
	// For now, we'll just return a dummy value
	return false
}

// DrawImage draws an image
func (ctx *Canvas2DContext) DrawImage(image string, x, y float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.ImageCalls++
	
	// This would normally draw an image on the canvas
	// For now, we'll just update the stats
}

// CreateLinearGradient creates a linear gradient
func (ctx *Canvas2DContext) CreateLinearGradient(x0, y0, x1, y1 float64) string {
	// This would normally create a linear gradient on the canvas
	// For now, we'll just return a dummy value
	return fmt.Sprintf("linear-gradient(%f,%f,%f,%f)", x0, y0, x1, y1)
}

// CreateRadialGradient creates a radial gradient
func (ctx *Canvas2DContext) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) string {
	// This would normally create a radial gradient on the canvas
	// For now, we'll just return a dummy value
	return fmt.Sprintf("radial-gradient(%f,%f,%f,%f,%f,%f)", x0, y0, r0, x1, y1, r1)
}

// CreatePattern creates a pattern
func (ctx *Canvas2DContext) CreatePattern(image, repetition string) string {
	// This would normally create a pattern on the canvas
	// For now, we'll just return a dummy value
	return fmt.Sprintf("pattern(%s,%s)", image, repetition)
}

// PutImageData puts image data
func (ctx *Canvas2DContext) PutImageData(imageData []byte, x, y float64) {
	ctx.Stats.DrawCalls++
	ctx.Stats.ImageCalls++
	
	// This would normally put image data on the canvas
	// For now, we'll just update the stats
}

// GetImageData gets image data
func (ctx *Canvas2DContext) GetImageData(x, y, width, height float64) []byte {
	ctx.Stats.ImageCalls++
	
	// This would normally get image data from the canvas
	// For now, we'll just return a dummy value
	return []byte{}
}

// Save saves the canvas state
func (ctx *Canvas2DContext) Save() {
	ctx.Stats.TransformCalls++
	
	// This would normally save the canvas state
	// For now, we'll just update the stats
}

// Restore restores the canvas state
func (ctx *Canvas2DContext) Restore() {
	ctx.Stats.TransformCalls++
	
	// This would normally restore the canvas state
	// For now, we'll just update the stats
}

// Scale scales the canvas
func (ctx *Canvas2DContext) Scale(x, y float64) {
	ctx.Stats.TransformCalls++
	
	// This would normally scale the canvas
	// For now, we'll just update the stats
}

// Rotate rotates the canvas
func (ctx *Canvas2DContext) Rotate(angle float64) {
	ctx.Stats.TransformCalls++
	
	// This would normally rotate the canvas
	// For now, we'll just update the stats
}

// Translate translates the canvas
func (ctx *Canvas2DContext) Translate(x, y float64) {
	ctx.Stats.TransformCalls++
	
	// This would normally translate the canvas
	// For now, we'll just update the stats
}

// Transform transforms the canvas
func (ctx *Canvas2DContext) Transform(a, b, c, d, e, f float64) {
	ctx.Stats.TransformCalls++
	
	// This would normally transform the canvas
	// For now, we'll just update the stats
}

// SetTransform sets the canvas transform
func (ctx *Canvas2DContext) SetTransform(a, b, c, d, e, f float64) {
	ctx.Stats.TransformCalls++
	
	// This would normally set the canvas transform
	// For now, we'll just update the stats
}

// ResetTransform resets the canvas transform
func (ctx *Canvas2DContext) ResetTransform() {
	ctx.Stats.TransformCalls++
	
	// This would normally reset the canvas transform
	// For now, we'll just update the stats
}