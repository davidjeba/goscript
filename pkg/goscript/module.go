package goscript

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// ModuleSpec describes a unit of GoScript code or capability.
type ModuleSpec struct {
	Name         string            `json:"name"`
	Path         string            `json:"path,omitempty"`
	Main         string            `json:"main,omitempty"`
	Version      string            `json:"version,omitempty"`
	Description  string            `json:"description,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Exports      []string          `json:"exports,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Normalize fills safe defaults for a module spec.
func (m *ModuleSpec) Normalize() {
	if m.Name == "" {
		m.Name = "module"
	}

	if m.Path == "" {
		m.Path = "."
	}

	if m.Main == "" {
		m.Main = m.Path
	}

	if m.Dependencies == nil {
		m.Dependencies = map[string]string{}
	}

	if m.Metadata == nil {
		m.Metadata = map[string]string{}
	}
}

// Validate ensures the module is fit for registration.
func (m ModuleSpec) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("module name is required")
	}

	if strings.TrimSpace(m.Path) == "" {
		return fmt.Errorf("module path is required")
	}

	return nil
}

// ModuleRegistry stores named module specs.
type ModuleRegistry struct {
	mu      sync.RWMutex
	modules map[string]ModuleSpec
}

// NewModuleRegistry creates a new registry.
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]ModuleSpec),
	}
}

// Register adds or updates a module specification.
func (r *ModuleRegistry) Register(spec ModuleSpec) error {
	spec.Normalize()
	if err := spec.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.modules[spec.Name] = spec
	return nil
}

// Get returns a registered module by name.
func (r *ModuleRegistry) Get(name string) (ModuleSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spec, ok := r.modules[name]
	return spec, ok
}

// List returns all registered modules.
func (r *ModuleRegistry) List() []ModuleSpec {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]ModuleSpec, 0, len(r.modules))
	for _, spec := range r.modules {
		out = append(out, spec)
	}
	return out
}

// Resolve tries to locate a module by name or by path fragment.
func (r *ModuleRegistry) Resolve(target string) (ModuleSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if spec, ok := r.modules[target]; ok {
		return spec, true
	}

	for _, spec := range r.modules {
		if spec.Path == target || spec.Main == target {
			return spec, true
		}

		if strings.Contains(filepath.ToSlash(spec.Path), filepath.ToSlash(target)) {
			return spec, true
		}
	}

	return ModuleSpec{}, false
}

// GlobalModuleRegistry is a shared registry for runtime discovery.
var GlobalModuleRegistry = NewModuleRegistry()

