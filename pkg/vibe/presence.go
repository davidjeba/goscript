package vibe

import "sync"

// PresenceState tracks whether a keyed element is entering, present, or exiting.
type PresenceState struct {
	Key      string `json:"key"`
	Present  bool   `json:"present"`
	Entering bool   `json:"entering"`
	Exiting  bool   `json:"exiting"`
}

// PresenceController provides AnimatePresence-style bookkeeping.
type PresenceController struct {
	mutex  sync.RWMutex
	states map[string]PresenceState
}

// NewPresenceController creates a new presence tracker.
func NewPresenceController() *PresenceController {
	return &PresenceController{
		states: make(map[string]PresenceState),
	}
}

// Enter marks an element as entering/present.
func (pc *PresenceController) Enter(key string) PresenceState {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	state := PresenceState{
		Key:      key,
		Present:  true,
		Entering: true,
		Exiting:  false,
	}
	pc.states[key] = state
	return state
}

// MarkPresent clears the entering state after the first paint.
func (pc *PresenceController) MarkPresent(key string) PresenceState {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	state := pc.states[key]
	state.Key = key
	state.Present = true
	state.Entering = false
	state.Exiting = false
	pc.states[key] = state
	return state
}

// Exit marks an element as exiting but still tracked.
func (pc *PresenceController) Exit(key string) PresenceState {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	state := pc.states[key]
	state.Key = key
	state.Present = false
	state.Entering = false
	state.Exiting = true
	pc.states[key] = state
	return state
}

// SafeToRemove removes a previously exiting element from tracking.
func (pc *PresenceController) SafeToRemove(key string) {
	pc.mutex.Lock()
	delete(pc.states, key)
	pc.mutex.Unlock()
}

// Snapshot returns a copy of the current presence states.
func (pc *PresenceController) Snapshot() []PresenceState {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	result := make([]PresenceState, 0, len(pc.states))
	for _, state := range pc.states {
		result = append(result, state)
	}
	return result
}
