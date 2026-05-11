package goscript

type SSREngine struct {
	store *Store
}

func NewSSREngine(store *Store) *SSREngine {
	return &SSREngine{store: store}
}

func (ssr *SSREngine) RenderToString(component Component) (string, error) {
	rendered := component.Render()
	state := map[string]interface{}{}
	if ssr != nil && ssr.store != nil {
		state = ssr.store.state
	}

	renderedShell, err := RenderHydrationShell(rendered, HydrationPayload{
		AppID:   "app",
		Version: "go-script",
		State:   state,
	})
	if err != nil {
		return "", err
	}

	return renderedShell, nil
}

