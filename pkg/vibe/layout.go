package vibe

// LayoutSnapshot represents an element box at a point in time.
type LayoutSnapshot struct {
	ID     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// LayoutDelta represents the transform needed to move between two layouts.
type LayoutDelta struct {
	TranslateX float64 `json:"translateX"`
	TranslateY float64 `json:"translateY"`
	ScaleX     float64 `json:"scaleX"`
	ScaleY     float64 `json:"scaleY"`
}

// ComputeLayoutDelta measures a transform-based delta between two layout snapshots.
func ComputeLayoutDelta(from, to LayoutSnapshot) LayoutDelta {
	delta := LayoutDelta{
		TranslateX: to.X - from.X,
		TranslateY: to.Y - from.Y,
		ScaleX:     1,
		ScaleY:     1,
	}

	if from.Width > 0 && to.Width > 0 {
		delta.ScaleX = to.Width / from.Width
	}
	if from.Height > 0 && to.Height > 0 {
		delta.ScaleY = to.Height / from.Height
	}

	return delta
}
