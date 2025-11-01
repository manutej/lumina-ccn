package main

// KeybindingHandler manages keyboard input routing and mode-aware key dispatch.
type KeybindingHandler struct {
	mode string
}

// NewKeybindingHandler creates a new keybinding handler in normal mode.
func NewKeybindingHandler() *KeybindingHandler {
	return &KeybindingHandler{mode: "normal"}
}

// HandleKey processes a key event based on the current mode.
// Returns nil to signal exit, otherwise returns the updated model.
func (kh *KeybindingHandler) HandleKey(key string, model interface{}) interface{} {
	switch key {
	case "q", "ctrl+c":
		return nil // Exit signal
	case "?":
		kh.SetMode("help")
	case "/":
		kh.SetMode("search")
	case "escape":
		kh.SetMode("normal")
	}

	return model
}

// SetMode updates the current interaction mode.
func (kh *KeybindingHandler) SetMode(mode string) {
	kh.mode = mode
}

// GetMode returns the current interaction mode.
func (kh *KeybindingHandler) GetMode() string {
	return kh.mode
}
