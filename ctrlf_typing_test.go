package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestSearchModeTyping(t *testing.T) {
	// Create model and enter search mode
	m := AppModel{
		currentMode: NormalMode,
		width:       120,
		height:      40,
		ready:       true,
	}

	// Step 1: Press Ctrl+F
	ctrlF := tea.KeyMsg{Type: tea.KeyCtrlF}
	result, _ := m.Update(ctrlF)
	m = result.(AppModel)

	t.Logf("After Ctrl+F: mode=%d searchPhase=%d searchQuery=%q", m.currentMode, m.searchPhase, m.searchQuery)

	if m.currentMode != SearchMode {
		t.Fatalf("Expected SearchMode, got %d", m.currentMode)
	}

	// Step 2: Type 't'
	keyT := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	t.Logf("Sending key: %q (Type=%d)", keyT.String(), keyT.Type)

	result, _ = m.Update(keyT)
	m = result.(AppModel)

	t.Logf("After 't': mode=%d searchPhase=%d searchQuery=%q", m.currentMode, m.searchPhase, m.searchQuery)

	if m.searchQuery != "t" {
		t.Errorf("Expected searchQuery='t', got %q", m.searchQuery)
	}

	// Step 3: Type 'e'
	keyE := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	result, _ = m.Update(keyE)
	m = result.(AppModel)

	t.Logf("After 'e': searchQuery=%q", m.searchQuery)

	if m.searchQuery != "te" {
		t.Errorf("Expected searchQuery='te', got %q", m.searchQuery)
	}

	// Step 4: Type 's'
	keyS := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	result, _ = m.Update(keyS)
	m = result.(AppModel)

	t.Logf("After 's': searchQuery=%q", m.searchQuery)

	if m.searchQuery != "tes" {
		t.Errorf("Expected searchQuery='tes', got %q", m.searchQuery)
	}
}
