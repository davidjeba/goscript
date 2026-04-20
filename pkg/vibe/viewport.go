package vibe

// ScrollState represents scroll-linked values.
type ScrollState struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	XProgress float64 `json:"xProgress"`
	YProgress float64 `json:"yProgress"`
}

// InViewOptions configures scroll-triggered observation.
type InViewOptions struct {
	Margin string  `json:"margin,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Once   bool    `json:"once,omitempty"`
}

// DefaultInViewOptions returns defaults inspired by motion's in-view controls.
func DefaultInViewOptions() InViewOptions {
	return InViewOptions{
		Margin: "0px",
		Amount: 0.25,
	}
}

// NormalizeScroll converts an absolute position into a 0..1 progress value.
func NormalizeScroll(position, total float64) float64 {
	if total <= 0 {
		return 0
	}
	if position <= 0 {
		return 0
	}
	if position >= total {
		return 1
	}
	return position / total
}
