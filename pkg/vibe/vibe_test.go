package vibe

import "testing"

func TestMotionValueSubscription(t *testing.T) {
	value := NewMotionValue(10)

	var last ChangeEvent
	unsubscribe := value.Subscribe(func(event ChangeEvent) {
		last = event
	})

	value.Set(25)
	unsubscribe()
	value.Set(40)

	if last.Previous != 10 || last.Current != 25 {
		t.Fatalf("unexpected change event: %+v", last)
	}

	if last.Velocity != 15 {
		t.Fatalf("expected velocity 15, got %v", last.Velocity)
	}
}

func TestResolveVariant(t *testing.T) {
	variants := VariantSet{
		"hidden": {
			Name:  "hidden",
			Style: StyleMap{"opacity": 0, "y": 24},
		},
		"visible": {
			Name:       "visible",
			Style:      StyleMap{"opacity": 1},
			Transition: Transition{Type: TransitionSpring, Duration: 0.4},
		},
	}

	style, transitions, err := ResolveVariant(variants, "hidden", "visible")
	if err != nil {
		t.Fatalf("ResolveVariant returned error: %v", err)
	}

	if style["opacity"] != 1 {
		t.Fatalf("expected merged opacity 1, got %v", style["opacity"])
	}

	if style["y"] != 24 {
		t.Fatalf("expected y 24, got %v", style["y"])
	}

	if len(transitions) != 1 || transitions[0].Type != TransitionSpring {
		t.Fatalf("unexpected transitions: %+v", transitions)
	}
}

func TestComputeLayoutDelta(t *testing.T) {
	from := LayoutSnapshot{ID: "card", X: 10, Y: 20, Width: 100, Height: 50}
	to := LayoutSnapshot{ID: "card", X: 40, Y: 35, Width: 200, Height: 100}

	delta := ComputeLayoutDelta(from, to)
	if delta.TranslateX != 30 || delta.TranslateY != 15 {
		t.Fatalf("unexpected translation delta: %+v", delta)
	}
	if delta.ScaleX != 2 || delta.ScaleY != 2 {
		t.Fatalf("unexpected scale delta: %+v", delta)
	}
}
