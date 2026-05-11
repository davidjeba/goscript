package vibe

// TransitionType describes the motion curve used for an animation.
type TransitionType string

const (
	TransitionTween   TransitionType = "tween"
	TransitionSpring  TransitionType = "spring"
	TransitionInertia TransitionType = "inertia"
)

// RepeatType controls how repeated animations loop.
type RepeatType string

const (
	RepeatLoop    RepeatType = "loop"
	RepeatReverse RepeatType = "reverse"
	RepeatMirror  RepeatType = "mirror"
)

// Transition describes the timing and physics for an animation.
type Transition struct {
	Type       TransitionType `json:"type,omitempty"`
	Duration   float64        `json:"duration,omitempty"`
	Delay      float64        `json:"delay,omitempty"`
	Ease       string         `json:"ease,omitempty"`
	Stiffness  float64        `json:"stiffness,omitempty"`
	Damping    float64        `json:"damping,omitempty"`
	Mass       float64        `json:"mass,omitempty"`
	Bounce     float64        `json:"bounce,omitempty"`
	Repeat     int            `json:"repeat,omitempty"`
	RepeatType RepeatType     `json:"repeatType,omitempty"`
}

// DefaultTransition returns a sensible default that feels close to motion-first UI work.
func DefaultTransition() Transition {
	return Transition{
		Type:       TransitionSpring,
		Duration:   0.3,
		Ease:       "easeOut",
		Stiffness:  220,
		Damping:    22,
		Mass:       1,
		Bounce:     0,
		RepeatType: RepeatLoop,
	}
}

// Normalize fills in transition defaults without clobbering explicit values.
func (t Transition) Normalize() Transition {
	defaults := DefaultTransition()

	if t.Type == "" {
		t.Type = defaults.Type
	}
	if t.Duration <= 0 {
		t.Duration = defaults.Duration
	}
	if t.Ease == "" {
		t.Ease = defaults.Ease
	}
	if t.Stiffness <= 0 {
		t.Stiffness = defaults.Stiffness
	}
	if t.Damping <= 0 {
		t.Damping = defaults.Damping
	}
	if t.Mass <= 0 {
		t.Mass = defaults.Mass
	}
	if t.RepeatType == "" {
		t.RepeatType = defaults.RepeatType
	}

	return t
}
