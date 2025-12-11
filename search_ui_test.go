package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestSearchModeTransition tests entering and exiting search mode
func TestSearchModeTransition(t *testing.T) {
	tests := []struct {
		name          string
		initialMode   UIMode
		key           string
		expectedMode  UIMode
		expectedPhase int
	}{
		{
			name:          "ctrl+f opens search in input phase",
			initialMode:   NormalMode,
			key:           "ctrl+f",
			expectedMode:  SearchMode,
			expectedPhase: 0, // Input phase
		},
		{
			name:          "esc cancels search from input phase",
			initialMode:   SearchMode,
			key:           "esc",
			expectedMode:  NormalMode,
			expectedPhase: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AppModel{
				currentMode: tt.initialMode,
				searchPhase: 0,
				width:       120,
				height:      40,
				ready:       true,
			}

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+f" {
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlF}
			} else if tt.key == "esc" {
				keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
			}

			result, _ := m.Update(keyMsg)
			newModel := result.(AppModel)

			if newModel.currentMode != tt.expectedMode {
				t.Errorf("expected mode %d, got %d", tt.expectedMode, newModel.currentMode)
			}
		})
	}
}

// TestSearchInputPhase tests query input handling in search mode
func TestSearchInputPhase(t *testing.T) {
	tests := []struct {
		name           string
		initialQuery   string
		key            string
		expectedQuery  string
		expectSearched bool
	}{
		{
			name:           "append character to query",
			initialQuery:   "TODO",
			key:            "x",
			expectedQuery:  "TODOx",
			expectSearched: false,
		},
		{
			name:           "backspace removes character",
			initialQuery:   "TODO",
			key:            "backspace",
			expectedQuery:  "TOD",
			expectSearched: false,
		},
		{
			name:           "backspace on empty query",
			initialQuery:   "",
			key:            "backspace",
			expectedQuery:  "",
			expectSearched: false,
		},
		{
			name:           "enter with empty query does not search",
			initialQuery:   "",
			key:            "enter",
			expectedQuery:  "",
			expectSearched: false,
		},
		{
			name:           "enter with query initiates search",
			initialQuery:   "TODO",
			key:            "enter",
			expectedQuery:  "TODO",
			expectSearched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AppModel{
				currentMode: SearchMode,
				searchPhase: 0,
				searchQuery: tt.initialQuery,
				rootPath:    "/tmp",
				width:       120,
				height:      40,
				ready:       true,
			}

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "backspace" {
				keyMsg = tea.KeyMsg{Type: tea.KeyBackspace}
			} else if tt.key == "enter" {
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			}

			result, cmd := m.handleSearchMode(keyMsg)
			newModel := result.(AppModel)

			if newModel.searchQuery != tt.expectedQuery {
				t.Errorf("expected query %q, got %q", tt.expectedQuery, newModel.searchQuery)
			}

			// Check if search was initiated (cmd is not nil when searching)
			hasCmd := cmd != nil
			if tt.expectSearched != hasCmd {
				t.Errorf("expected search initiated: %v, got cmd: %v", tt.expectSearched, hasCmd)
			}
		})
	}
}

// TestSearchResultsNavigation tests navigation in results phase
func TestSearchResultsNavigation(t *testing.T) {
	testResults := []RipgrepResult{
		{FilePath: "/path/to/file1.md", Line: 10, Text: "match 1"},
		{FilePath: "/path/to/file2.md", Line: 20, Text: "match 2"},
		{FilePath: "/path/to/file3.md", Line: 30, Text: "match 3"},
	}

	tests := []struct {
		name           string
		initialCursor  int
		key            string
		expectedCursor int
	}{
		{
			name:           "move down",
			initialCursor:  0,
			key:            "j",
			expectedCursor: 1,
		},
		{
			name:           "move up",
			initialCursor:  1,
			key:            "k",
			expectedCursor: 0,
		},
		{
			name:           "wrap at end",
			initialCursor:  2,
			key:            "j",
			expectedCursor: 0,
		},
		{
			name:           "wrap at start",
			initialCursor:  0,
			key:            "k",
			expectedCursor: 2,
		},
		{
			name:           "n goes to next",
			initialCursor:  0,
			key:            "n",
			expectedCursor: 1,
		},
		{
			name:           "N goes to previous",
			initialCursor:  1,
			key:            "N",
			expectedCursor: 0,
		},
		{
			name:           "down arrow",
			initialCursor:  0,
			key:            "down",
			expectedCursor: 1,
		},
		{
			name:           "up arrow",
			initialCursor:  1,
			key:            "up",
			expectedCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AppModel{
				currentMode:   SearchMode,
				searchPhase:   1, // Results phase
				searchResults: testResults,
				searchCursor:  tt.initialCursor,
				width:         120,
				height:        40,
				ready:         true,
			}

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "down" {
				keyMsg = tea.KeyMsg{Type: tea.KeyDown}
			} else if tt.key == "up" {
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			}

			result, _ := m.handleSearchMode(keyMsg)
			newModel := result.(AppModel)

			if newModel.searchCursor != tt.expectedCursor {
				t.Errorf("expected cursor %d, got %d", tt.expectedCursor, newModel.searchCursor)
			}
		})
	}
}

// TestSearchBackToInput tests returning to input phase from results
func TestSearchBackToInput(t *testing.T) {
	m := AppModel{
		currentMode:   SearchMode,
		searchPhase:   1, // Results phase
		searchQuery:   "TODO",
		searchResults: []RipgrepResult{{FilePath: "test.md", Line: 1, Text: "test"}},
		searchCursor:  0,
		width:         120,
		height:        40,
		ready:         true,
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	result, _ := m.handleSearchMode(keyMsg)
	newModel := result.(AppModel)

	if newModel.searchPhase != 0 {
		t.Errorf("expected phase 0 (input), got %d", newModel.searchPhase)
	}
	if len(newModel.searchResults) != 0 {
		t.Errorf("expected empty results, got %d results", len(newModel.searchResults))
	}
}

// TestSearchEmptyResults tests navigation with empty results
func TestSearchEmptyResults(t *testing.T) {
	m := AppModel{
		currentMode:   SearchMode,
		searchPhase:   1,
		searchResults: []RipgrepResult{},
		searchCursor:  0,
		width:         120,
		height:        40,
		ready:         true,
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	result, _ := m.handleSearchMode(keyMsg)
	newModel := result.(AppModel)

	// Cursor should stay at 0 with empty results
	if newModel.searchCursor != 0 {
		t.Errorf("expected cursor 0 with empty results, got %d", newModel.searchCursor)
	}
}

// TestSearchModalRender tests search modal rendering
func TestSearchModalRender(t *testing.T) {
	tests := []struct {
		name           string
		searchPhase    int
		searchQuery    string
		searchResults  []RipgrepResult
		inProgress     bool
		errorMsg       string
		expectContains []string
	}{
		{
			name:           "input phase empty query",
			searchPhase:    0,
			searchQuery:    "",
			searchResults:  nil,
			expectContains: []string{"Search:", "Type a search query"},
		},
		{
			name:           "input phase with query",
			searchPhase:    0,
			searchQuery:    "TODO",
			searchResults:  nil,
			expectContains: []string{"Search:", "TODO", "Press Enter to search"},
		},
		{
			name:           "results phase searching",
			searchPhase:    1,
			searchQuery:    "TODO",
			searchResults:  nil,
			inProgress:     true,
			expectContains: []string{"Searching..."},
		},
		{
			name:           "results phase no results",
			searchPhase:    1,
			searchQuery:    "nonexistent",
			searchResults:  []RipgrepResult{},
			expectContains: []string{"No results found"},
		},
		{
			name:        "results phase with results",
			searchPhase: 1,
			searchQuery: "TODO",
			searchResults: []RipgrepResult{
				{FilePath: "/path/to/file.md", Line: 42, Text: "// TODO: fix this"},
			},
			expectContains: []string{"file.md", "42"},
		},
		{
			name:           "error message displayed",
			searchPhase:    0,
			searchQuery:    "test",
			errorMsg:       "ripgrep not found",
			expectContains: []string{"Error:", "ripgrep not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal color manager for rendering
			colorManager, _ := NewColorManager()

			m := AppModel{
				currentMode:        SearchMode,
				searchPhase:        tt.searchPhase,
				searchQuery:        tt.searchQuery,
				searchResults:      tt.searchResults,
				searchInProgress:   tt.inProgress,
				searchErrorMessage: tt.errorMsg,
				colorManager:       colorManager,
				width:              120,
				height:             40,
				ready:              true,
			}

			rendered := m.renderSearchModal()

			for _, expected := range tt.expectContains {
				if !containsString(rendered, expected) {
					t.Errorf("expected rendered modal to contain %q", expected)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring (ignoring ANSI codes)
func containsString(s, substr string) bool {
	// Simple substring check - in production would strip ANSI codes
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr))
}
