package goscript

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOptionAndResult(t *testing.T) {
	opt := Some("hello")
	if got := opt.UnwrapOr("fallback"); got != "hello" {
		t.Fatalf("expected hello, got %v", got)
	}

	res := Ok(42)
	if !res.IsOk() {
		t.Fatalf("expected result to be ok")
	}
	if got := res.UnwrapOr(0); got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}

	if Match("alpha",
		MatchCase{Equals: "beta", Then: func(v interface{}) interface{} { return "nope" }},
		MatchCase{Kind: 0, Then: func(v interface{}) interface{} { return "any" }},
	) != nil {
		t.Fatalf("unexpected match result")
	}
}

func TestScheduler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := NewScheduler(2)
	scheduler.Start(ctx)
	defer scheduler.Stop()

	if err := scheduler.Submit(Task{
		Name: "echo",
		Handler: func(ctx context.Context) (interface{}, error) {
			return "ok", nil
		},
	}); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	select {
	case result := <-scheduler.Results():
		if result.Err != nil || result.Value != "ok" {
			t.Fatalf("unexpected result: %+v", result)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for task result")
	}
}

func TestRealtimeHub(t *testing.T) {
	hub := NewRealtimeHub(2)
	sub := hub.Subscribe("system", "listener-1")
	defer hub.Unsubscribe("system", "listener-1")

	hub.Ping("system", "unit-test")

	select {
	case event := <-sub:
		if event.Kind != "ping" {
			t.Fatalf("expected ping event, got %q", event.Kind)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for ping event")
	}

	if events := hub.Poll("system", 10); len(events) == 0 {
		t.Fatalf("expected history to contain events")
	}
}

func TestInferenceRouter(t *testing.T) {
	router := NewInferenceRouter(
		fakeProvider{label: "local"},
		nil,
		fakeProvider{label: "fallback"},
	)

	response, err := router.Infer(context.Background(), InferenceRequest{Model: "tiny"})
	if err != nil {
		t.Fatalf("unexpected inference error: %v", err)
	}
	if response.Provider != "local" {
		t.Fatalf("expected local provider, got %q", response.Provider)
	}
}

type fakeProvider struct {
	label string
}

func (f fakeProvider) Infer(ctx context.Context, request InferenceRequest) (InferenceResponse, error) {
	if request.Model == "fail" {
		return InferenceResponse{}, errors.New("forced failure")
	}

	return InferenceResponse{
		Model:    request.Model,
		Output:   request.Input,
		Provider: f.label,
	}, nil
}

