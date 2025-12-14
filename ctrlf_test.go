package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestCtrlFActivatesSearchMode(t *testing.T) {
	// Create a minimal model
	m := AppModel{
		currentMode: NormalMode,
		width:       120,
		height:      40,
		ready:       true,
	}

	t.Logf("Initial mode: %d (NormalMode=%d, SearchMode=%d)", m.currentMode, NormalMode, SearchMode)

	// Simulate Ctrl+F
	ctrlF := tea.KeyMsg{Type: tea.KeyCtrlF}
	t.Logf("Sending key: %q", ctrlF.String())

	result, _ := m.Update(ctrlF)
	newModel := result.(AppModel)

	t.Logf("After Ctrl+F: mode=%d", newModel.currentMode)

	if newModel.currentMode != SearchMode {
		t.Errorf("Expected SearchMode (%d), got %d", SearchMode, newModel.currentMode)
	}
}

func TestSlashActivatesFinderMode(t *testing.T) {
	// Create a minimal model
	m := AppModel{
		currentMode:   NormalMode,
		markdownFiles: []string{"test.md"}, // Pre-load so it doesn't go to LoadingMode
		width:         120,
		height:        40,
		ready:         true,
	}

	t.Logf("Initial mode: %d", m.currentMode)

	// Simulate /
	slash := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	t.Logf("Sending key: %q", slash.String())

	result, _ := m.Update(slash)
	newModel := result.(AppModel)

	t.Logf("After /: mode=%d", newModel.currentMode)

	if newModel.currentMode != FinderMode {
		t.Errorf("Expected FinderMode (%d), got %d", FinderMode, newModel.currentMode)
	}
}
