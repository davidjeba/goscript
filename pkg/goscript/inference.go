package goscript

import (
	"context"
	"fmt"
	"sync"
)

// InferenceRequest describes an edge or remote inference job.
type InferenceRequest struct {
	Model    string                 `json:"model"`
	Input    interface{}            `json:"input"`
	Device   string                 `json:"device,omitempty"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// InferenceResponse describes the result of inference.
type InferenceResponse struct {
	Model    string                 `json:"model"`
	Output   interface{}            `json:"output"`
	Provider string                 `json:"provider"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// InferenceProvider performs inference for a request.
type InferenceProvider interface {
	Infer(context.Context, InferenceRequest) (InferenceResponse, error)
}

// InferenceRouter routes requests to the best provider.
type InferenceRouter struct {
	mu        sync.RWMutex
	local     InferenceProvider
	remote    InferenceProvider
	fallback  InferenceProvider
}

// NewInferenceRouter creates a router with optional providers.
func NewInferenceRouter(local, remote, fallback InferenceProvider) *InferenceRouter {
	return &InferenceRouter{
		local:    local,
		remote:   remote,
		fallback: fallback,
	}
}

// SetLocal updates the local provider.
func (r *InferenceRouter) SetLocal(provider InferenceProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.local = provider
}

// SetRemote updates the remote provider.
func (r *InferenceRouter) SetRemote(provider InferenceProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.remote = provider
}

// SetFallback updates the fallback provider.
func (r *InferenceRouter) SetFallback(provider InferenceProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = provider
}

// Infer executes the best available provider.
func (r *InferenceRouter) Infer(ctx context.Context, request InferenceRequest) (InferenceResponse, error) {
	r.mu.RLock()
	local := r.local
	remote := r.remote
	fallback := r.fallback
	r.mu.RUnlock()

	if local != nil {
		if response, err := local.Infer(ctx, request); err == nil {
			return response, nil
		}
	}

	if remote != nil {
		if response, err := remote.Infer(ctx, request); err == nil {
			return response, nil
		}
	}

	if fallback != nil {
		return fallback.Infer(ctx, request)
	}

	return InferenceResponse{}, fmt.Errorf("no inference provider available for model %q", request.Model)
}

