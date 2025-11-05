package main

import (
	"strings"
)

// FuzzyFinderImpl provides fuzzy finding and filtering for a list of items.
// Supports navigation with wrapping, case-insensitive filtering, and jump-to-letter.
// Thread-safe for concurrent reads (cursor, filter state).
type FuzzyFinderImpl struct {
	items    []string // Original item list (immutable after creation)
	cursor   int      // Current cursor position in filtered results
	filter   string   // Current filter query
	filtered []string // Filtered results based on current query
	mode     string   // Interaction mode (normal, search, help)
}

// NewFuzzyFinderImpl creates a new fuzzy finder with the given items.
func NewFuzzyFinderImpl(items []string) *FuzzyFinderImpl {
	return &FuzzyFinderImpl{
		items:    items,
		cursor:   0,
		filter:   "",
		filtered: items,
		mode:     "normal",
	}
}

// SetCursor sets the cursor position if within valid bounds.
func (ff *FuzzyFinderImpl) SetCursor(index int) {
	if index >= 0 && index < len(ff.filtered) {
		ff.cursor = index
	}
}

// Cursor returns the current cursor position.
func (ff *FuzzyFinderImpl) Cursor() int {
	return ff.cursor
}

// SelectedItem returns the currently selected item.
func (ff *FuzzyFinderImpl) SelectedItem() string {
	if ff.cursor >= 0 && ff.cursor < len(ff.filtered) {
		return ff.filtered[ff.cursor]
	}
	return ""
}

// HandleKey processes navigation keys with wrapping support.
func (ff *FuzzyFinderImpl) HandleKey(key string) {
	if len(ff.filtered) == 0 {
		return
	}

	switch key {
	case "down", "j":
		ff.cursor = (ff.cursor + 1) % len(ff.filtered)
	case "up", "k":
		ff.cursor = (ff.cursor - 1 + len(ff.filtered)) % len(ff.filtered)
	case "pagedown":
		ff.cursor = min(ff.cursor+10, len(ff.filtered)-1)
	case "pageup":
		ff.cursor = max(ff.cursor-10, 0)
	case "home":
		ff.cursor = 0
	case "end":
		ff.cursor = len(ff.filtered) - 1
	}
}

// SetFilter applies a filter query and updates filtered results.
func (ff *FuzzyFinderImpl) SetFilter(query string) {
	ff.filter = query
	ff.cursor = 0

	if query == "" {
		ff.filtered = ff.items
		return
	}

	ff.filtered = make([]string, 0, len(ff.items))
	lowerQuery := strings.ToLower(query)

	for _, item := range ff.items {
		if strings.Contains(strings.ToLower(item), lowerQuery) {
			ff.filtered = append(ff.filtered, item)
		}
	}
}

// Filter returns the current filter query.
func (ff *FuzzyFinderImpl) Filter() string {
	return ff.filter
}

// FilteredResults returns the current filtered item list.
func (ff *FuzzyFinderImpl) FilteredResults() []string {
	return ff.filtered
}

// SetMode sets the current interaction mode.
func (ff *FuzzyFinderImpl) SetMode(mode string) {
	ff.mode = mode
}

// Mode returns the current interaction mode.
func (ff *FuzzyFinderImpl) Mode() string {
	return ff.mode
}

// JumpToLetter moves cursor to the first item starting with the given letter.
func (ff *FuzzyFinderImpl) JumpToLetter(letter string) error {
	lowerLetter := strings.ToLower(letter)
	for i, item := range ff.filtered {
		if strings.HasPrefix(strings.ToLower(item), lowerLetter) {
			ff.cursor = i
			return nil
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
