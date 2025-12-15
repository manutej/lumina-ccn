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

func TestSearchModePaste(t *testing.T) {
	// Create model and enter search mode
	m := AppModel{
		currentMode: NormalMode,
		width:       120,
		height:      40,
		ready:       true,
	}

	// Enter search mode
	ctrlF := tea.KeyMsg{Type: tea.KeyCtrlF}
	result, _ := m.Update(ctrlF)
	m = result.(AppModel)

	if m.currentMode != SearchMode {
		t.Fatalf("Expected SearchMode, got %d", m.currentMode)
	}

	// Simulate paste - multiple runes in a single KeyMsg
	pasteMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello world")}
	result, _ = m.Update(pasteMsg)
	m = result.(AppModel)

	if m.searchQuery != "hello world" {
		t.Errorf("Expected searchQuery='hello world', got %q", m.searchQuery)
	}

	// Type another character to verify normal typing still works after paste
	keyX := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}}
	result, _ = m.Update(keyX)
	m = result.(AppModel)

	if m.searchQuery != "hello world!" {
		t.Errorf("Expected searchQuery='hello world!', got %q", m.searchQuery)
	}
}

func TestFinderModePaste(t *testing.T) {
	// Create model already in finder mode with files loaded
	// (simulating after async file collection completes)
	m := AppModel{
		currentMode:    FinderMode,
		width:          120,
		height:         40,
		ready:          true,
		finderItems:    []string{"README.md", "main.go", "model.go", "test_file.go"},
		finderFiltered: []string{"README.md", "main.go", "model.go", "test_file.go"},
		finderInput:    "",
	}

	// Simulate paste - multiple runes
	pasteMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("main")}
	result, _ := m.Update(pasteMsg)
	m = result.(AppModel)

	if m.finderInput != "main" {
		t.Errorf("Expected finderInput='main', got %q", m.finderInput)
	}

	// Verify filtering worked
	if len(m.finderFiltered) != 1 || m.finderFiltered[0] != "main.go" {
		t.Errorf("Expected filtered to contain only 'main.go', got %v", m.finderFiltered)
	}
}
