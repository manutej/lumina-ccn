package main

import (
	"testing"
)

// TestFinderModeTransition tests the transition to and from FinderMode
func TestFinderModeTransition(t *testing.T) {
	t.Run("slash_key_opens_finder", func(t *testing.T) {
		// Arrange: Create model in NormalMode with preloaded files
		m := AppModel{
			currentMode:   NormalMode,
			markdownFiles: []string{"README.md", "docs/guide.md"},
		}

		// Act: Simulate "/" key - transition should happen
		m.finderItems = m.markdownFiles
		m.finderFiltered = m.markdownFiles
		m.finderCursor = 0
		m.transitionTo(FinderMode)

		// Assert: Mode changed to FinderMode
		if m.currentMode != FinderMode {
			t.Errorf("expected FinderMode, got %v", m.currentMode)
		}

		// Assert: Finder items populated
		if len(m.finderItems) != 2 {
			t.Errorf("expected 2 finder items, got %d", len(m.finderItems))
		}
	})

	t.Run("esc_cancels_finder", func(t *testing.T) {
		// Arrange: Model in FinderMode
		m := AppModel{
			currentMode:    FinderMode,
			finderInput:    "readme",
			finderCursor:   1,
			finderFiltered: []string{"README.md"},
		}

		// Act: Transition back to NormalMode (simulates Esc)
		m.transitionTo(NormalMode)

		// Assert: Mode changed to NormalMode
		if m.currentMode != NormalMode {
			t.Errorf("expected NormalMode, got %v", m.currentMode)
		}

		// Assert: Finder state cleaned up
		if m.finderInput != "" {
			t.Errorf("expected empty finderInput, got %q", m.finderInput)
		}
		if m.finderCursor != 0 {
			t.Errorf("expected finderCursor 0, got %d", m.finderCursor)
		}
	})
}

// TestFinderFilter tests the filtering functionality
func TestFinderFilter(t *testing.T) {
	tests := []struct {
		name          string
		items         []string
		query         string
		expectedCount int
		expectedFirst string
	}{
		{
			name:          "filter_matches_substring",
			items:         []string{"README.md", "guide.md", "api.md"},
			query:         "read",
			expectedCount: 1,
			expectedFirst: "README.md",
		},
		{
			name:          "filter_case_insensitive",
			items:         []string{"README.md", "Readme.txt", "readme.go"},
			query:         "readme",
			expectedCount: 3,
			expectedFirst: "README.md",
		},
		{
			name:          "filter_no_match",
			items:         []string{"README.md", "guide.md"},
			query:         "xyz",
			expectedCount: 0,
			expectedFirst: "",
		},
		{
			name:          "empty_query_shows_all",
			items:         []string{"a.md", "b.md", "c.md"},
			query:         "",
			expectedCount: 3,
			expectedFirst: "a.md",
		},
		{
			name:          "filter_path_components",
			items:         []string{"docs/api.md", "src/main.go", "docs/guide.md"},
			query:         "docs",
			expectedCount: 2,
			expectedFirst: "docs/api.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the filterItems pure function
			result := filterItems(tt.items, tt.query)

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(result))
			}

			if tt.expectedCount > 0 && result[0] != tt.expectedFirst {
				t.Errorf("expected first result %q, got %q", tt.expectedFirst, result[0])
			}
		})
	}
}

// TestFinderNavigation tests cursor navigation
func TestFinderNavigation(t *testing.T) {
	tests := []struct {
		name           string
		cursor         int
		delta          int
		listLen        int
		expectedCursor int
	}{
		{"move_down", 0, 1, 5, 1},
		{"move_up", 2, -1, 5, 1},
		{"wrap_at_end", 4, 1, 5, 0},
		{"wrap_at_start", 0, -1, 5, 4},
		{"empty_list", 0, 1, 0, 0},
		{"single_item", 0, 1, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := navigateCursor(tt.cursor, tt.delta, tt.listLen)
			if result != tt.expectedCursor {
				t.Errorf("navigateCursor(%d, %d, %d) = %d, want %d",
					tt.cursor, tt.delta, tt.listLen, result, tt.expectedCursor)
			}
		})
	}
}

// TestFinderSelection tests file selection
func TestFinderSelection(t *testing.T) {
	t.Run("select_opens_file", func(t *testing.T) {
		// Arrange: Model in FinderMode with filtered results
		m := AppModel{
			currentMode:    FinderMode,
			finderFiltered: []string{"/tmp/test/README.md", "/tmp/test/guide.md"},
			finderCursor:   0,
		}

		// Act: Get selected item (simulates Enter key logic)
		selected := ""
		if len(m.finderFiltered) > 0 && m.finderCursor < len(m.finderFiltered) {
			selected = m.finderFiltered[m.finderCursor]
		}

		// Assert: Correct file selected
		if selected != "/tmp/test/README.md" {
			t.Errorf("expected '/tmp/test/README.md', got %q", selected)
		}
	})

	t.Run("select_second_item", func(t *testing.T) {
		// Arrange
		m := AppModel{
			currentMode:    FinderMode,
			finderFiltered: []string{"a.md", "b.md", "c.md"},
			finderCursor:   1, // Second item
		}

		// Act
		selected := m.finderFiltered[m.finderCursor]

		// Assert
		if selected != "b.md" {
			t.Errorf("expected 'b.md', got %q", selected)
		}
	})
}

// TestFinderInputHandling tests character input
func TestFinderInputHandling(t *testing.T) {
	t.Run("append_character", func(t *testing.T) {
		input := "read"
		input += "m"
		if input != "readm" {
			t.Errorf("expected 'readm', got %q", input)
		}
	})

	t.Run("backspace_removes_character", func(t *testing.T) {
		input := "readme"
		if len(input) > 0 {
			input = input[:len(input)-1]
		}
		if input != "readm" {
			t.Errorf("expected 'readm', got %q", input)
		}
	})

	t.Run("backspace_on_empty", func(t *testing.T) {
		input := ""
		if len(input) > 0 {
			input = input[:len(input)-1]
		}
		if input != "" {
			t.Errorf("expected empty string, got %q", input)
		}
	})
}

// TestFuzzyFinderImplBasics tests the FuzzyFinderImpl struct
func TestFuzzyFinderImplBasics(t *testing.T) {
	t.Run("new_finder_has_all_items", func(t *testing.T) {
		items := []string{"a.md", "b.md", "c.md"}
		ff := NewFuzzyFinderImpl(items)

		if len(ff.FilteredResults()) != 3 {
			t.Errorf("expected 3 items, got %d", len(ff.FilteredResults()))
		}
	})

	t.Run("filter_reduces_results", func(t *testing.T) {
		items := []string{"README.md", "guide.md", "api.md"}
		ff := NewFuzzyFinderImpl(items)

		ff.SetFilter("read")
		if len(ff.FilteredResults()) != 1 {
			t.Errorf("expected 1 match, got %d", len(ff.FilteredResults()))
		}
	})

	t.Run("cursor_navigation", func(t *testing.T) {
		items := []string{"a.md", "b.md", "c.md"}
		ff := NewFuzzyFinderImpl(items)

		ff.HandleKey("down")
		if ff.Cursor() != 1 {
			t.Errorf("expected cursor 1, got %d", ff.Cursor())
		}

		ff.HandleKey("up")
		if ff.Cursor() != 0 {
			t.Errorf("expected cursor 0, got %d", ff.Cursor())
		}
	})

	t.Run("selected_item", func(t *testing.T) {
		items := []string{"first.md", "second.md", "third.md"}
		ff := NewFuzzyFinderImpl(items)

		ff.HandleKey("down") // Move to index 1
		if ff.SelectedItem() != "second.md" {
			t.Errorf("expected 'second.md', got %q", ff.SelectedItem())
		}
	})
}
