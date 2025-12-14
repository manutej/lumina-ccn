package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSearchModeRenderShowsTypedText(t *testing.T) {
	// Create model with ColorManager needed for rendering
	colorManager, _ := NewColorManager()

	m := AppModel{
		currentMode:  NormalMode,
		width:        120,
		height:       40,
		ready:        true,
		colorManager: colorManager,
	}

	// Step 1: Press Ctrl+F
	ctrlF := tea.KeyMsg{Type: tea.KeyCtrlF}
	result, _ := m.Update(ctrlF)
	m = result.(AppModel)

	if m.currentMode != SearchMode {
		t.Fatalf("Expected SearchMode, got %d", m.currentMode)
	}

	// Step 2: Type 'test'
	for _, char := range "test" {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		result, _ = m.Update(keyMsg)
		m = result.(AppModel)
	}

	t.Logf("After typing 'test': searchQuery=%q", m.searchQuery)

	if m.searchQuery != "test" {
		t.Fatalf("Expected searchQuery='test', got %q", m.searchQuery)
	}

	// Step 3: Render the View and check if "test" appears
	rendered := m.View()
	t.Logf("Rendered view length: %d", len(rendered))

	if !strings.Contains(rendered, "test") {
		t.Errorf("Rendered view does not contain 'test'")
		// Show first 500 chars of rendered view for debugging
		if len(rendered) > 500 {
			t.Logf("First 500 chars: %s", rendered[:500])
		} else {
			t.Logf("Full render: %s", rendered)
		}
	} else {
		t.Log("SUCCESS: 'test' found in rendered view")
	}
}

func TestSearchAndFinderInputConsistency(t *testing.T) {
	colorManager, _ := NewColorManager()

	// Test Search Mode typing
	t.Run("Search Mode", func(t *testing.T) {
		m := AppModel{
			currentMode:  SearchMode,
			searchPhase:  0,
			searchQuery:  "",
			width:        120,
			height:       40,
			ready:        true,
			colorManager: colorManager,
		}

		// Type 'a'
		keyA := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		result, _ := m.handleSearchMode(keyA)
		m = result.(AppModel)

		if m.searchQuery != "a" {
			t.Errorf("Search: expected query='a', got %q", m.searchQuery)
		}
	})

	// Test Finder Mode typing
	t.Run("Finder Mode", func(t *testing.T) {
		m := AppModel{
			currentMode:    FinderMode,
			finderInput:    "",
			finderItems:    []string{"test.md"},
			finderFiltered: []string{"test.md"},
			width:          120,
			height:         40,
			ready:          true,
			colorManager:   colorManager,
		}

		// Type 'a'
		keyA := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		result, _ := m.handleFinderMode(keyA)
		m = result.(AppModel)

		if m.finderInput != "a" {
			t.Errorf("Finder: expected input='a', got %q", m.finderInput)
		}
	})
}
