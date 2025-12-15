package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// KeyBinding represents a single key action
type KeyBinding struct {
	Key    string `json:"key"`
	Action string `json:"action"`
	View   string `json:"view"` // "viewer", "filetree", "any"
}

// KeyBindings represents the complete keybinding configuration
type KeyBindings struct {
	Navigation  []KeyBinding `json:"navigation"`
	Scrolling   []KeyBinding `json:"scrolling"`
	Actions     []KeyBinding `json:"actions"`
	AppControls []KeyBinding `json:"app_controls"`
}

// DefaultKeyBindings returns the default vim-style keybindings
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Navigation: []KeyBinding{
			{Key: "j", Action: "down", View: "filetree"},
			{Key: "down", Action: "down", View: "filetree"},
			{Key: "k", Action: "up", View: "filetree"},
			{Key: "up", Action: "up", View: "filetree"},
			{Key: "h", Action: "back", View: "filetree"},
			{Key: "backspace", Action: "back", View: "filetree"},
			{Key: "esc", Action: "back", View: "filetree"},
			{Key: "enter", Action: "open", View: "filetree"},
		},
		Scrolling: []KeyBinding{
			{Key: "j", Action: "scroll_down", View: "viewer"},
			{Key: "down", Action: "scroll_down", View: "viewer"},
			{Key: "k", Action: "scroll_up", View: "viewer"},
			{Key: "up", Action: "scroll_up", View: "viewer"},
			{Key: "d", Action: "page_down", View: "viewer"},
			{Key: "u", Action: "page_up", View: "viewer"},
			{Key: "g", Action: "top", View: "viewer"},
			{Key: "G", Action: "bottom", View: "viewer"},
			{Key: "alt+down", Action: "page_down", View: "viewer"},
			{Key: "alt+up", Action: "page_up", View: "viewer"},
			{Key: "alt+right", Action: "view_down", View: "viewer"},
			{Key: "alt+left", Action: "view_up", View: "viewer"},
			{Key: "v", Action: "start_selection", View: "viewer"},
			{Key: "V", Action: "start_line_selection", View: "viewer"},
		},
		Actions: []KeyBinding{
			{Key: "/", Action: "filter", View: "filetree"},
			{Key: "s", Action: "sort", View: "filetree"},
			{Key: "?", Action: "help", View: "any"},
			{Key: "tab", Action: "switch_view", View: "any"},
			{Key: "y", Action: "copy", View: "viewer"},
			{Key: "e", Action: "edit", View: "viewer"},
		},
		AppControls: []KeyBinding{
			{Key: "q", Action: "quit", View: "any"},
			{Key: "ctrl+c", Action: "quit", View: "any"},
		},
	}
}

// LoadKeyBindings loads keybindings from config file, or returns defaults
func LoadKeyBindings() KeyBindings {
	configDir, err := getConfigDir()
	if err != nil {
		return DefaultKeyBindings()
	}

	configPath := filepath.Join(configDir, "keybindings.json")

	// If config doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		createDefaultConfig(configPath)
		return DefaultKeyBindings()
	}

	// Read and parse config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultKeyBindings()
	}

	var kb KeyBindings
	if err := json.Unmarshal(data, &kb); err != nil {
		return DefaultKeyBindings()
	}

	return kb
}

// FindAction returns the action for a given key in a view context
func (kb KeyBindings) FindAction(key string, view string) string {
	// Search through all categories
	allBindings := append(
		append(
			append(kb.Navigation, kb.Scrolling...),
			kb.Actions...,
		),
		kb.AppControls...,
	)

	for _, binding := range allBindings {
		if binding.Key == key && (binding.View == view || binding.View == "any") {
			return binding.Action
		}
	}
	return ""
}

// GetConfigDir returns the config directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "lumina"), nil
}

// createDefaultConfig creates the default config file
func createDefaultConfig(configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	defaults := DefaultKeyBindings()
	data, err := json.MarshalIndent(defaults, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ExampleCustomConfig returns an example custom configuration
func ExampleCustomConfig() string {
	config := KeyBindings{
		Navigation: []KeyBinding{
			{Key: "j", Action: "down", View: "filetree"},
			{Key: "down", Action: "down", View: "filetree"},
			{Key: "k", Action: "up", View: "filetree"},
			{Key: "up", Action: "up", View: "filetree"},
			{Key: "h", Action: "back", View: "filetree"},
			{Key: "esc", Action: "back", View: "filetree"},
			{Key: "enter", Action: "open", View: "filetree"},
		},
		Scrolling: []KeyBinding{
			// CUSTOM: d=down, e=up (one-handed navigation)
			{Key: "j", Action: "scroll_down", View: "viewer"},
			{Key: "down", Action: "scroll_down", View: "viewer"},
			{Key: "k", Action: "scroll_up", View: "viewer"},
			{Key: "up", Action: "scroll_up", View: "viewer"},
			{Key: "d", Action: "page_down", View: "viewer"},
			{Key: "e", Action: "page_up", View: "viewer"}, // CUSTOM!
			{Key: "g", Action: "top", View: "viewer"},
			{Key: "G", Action: "bottom", View: "viewer"},
		},
		Actions: []KeyBinding{
			{Key: "/", Action: "filter", View: "filetree"},
			{Key: "s", Action: "sort", View: "filetree"},
			{Key: "?", Action: "help", View: "any"},
			{Key: "tab", Action: "switch_view", View: "any"},
			{Key: "y", Action: "copy", View: "viewer"},
		},
		AppControls: []KeyBinding{
			{Key: "q", Action: "quit", View: "any"},
			{Key: "ctrl+c", Action: "quit", View: "any"},
		},
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	return string(data)
}
