package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// SelectionState represents text selection in the viewer
type SelectionState struct {
	Enabled      bool   // Whether selection is active
	StartLine    int    // Line where selection starts
	StartCol     int    // Column where selection starts
	EndLine      int    // Line where selection ends
	EndCol       int    // Column where selection ends
	IsRectangular bool   // Block selection vs line selection
}

// ClipboardManager handles copy/paste operations
type ClipboardManager struct {
	selection SelectionState
	lastCopy  string // Last copied text
}

// NewClipboardManager creates a new clipboard manager
func NewClipboardManager() *ClipboardManager {
	return &ClipboardManager{
		selection: SelectionState{Enabled: false},
	}
}

// StartSelection marks the beginning of a selection
func (cm *ClipboardManager) StartSelection(line, col int, rectangular bool) {
	cm.selection.Enabled = true
	cm.selection.StartLine = line
	cm.selection.StartCol = col
	cm.selection.EndLine = line
	cm.selection.EndCol = col
	cm.selection.IsRectangular = rectangular
}

// ExtendSelection extends selection to a new position
func (cm *ClipboardManager) ExtendSelection(line, col int) {
	if !cm.selection.Enabled {
		return
	}
	cm.selection.EndLine = line
	cm.selection.EndCol = col
}

// GetSelection returns the currently selected text from content
func (cm *ClipboardManager) GetSelection(content string) string {
	if !cm.selection.Enabled {
		return ""
	}

	lines := splitLines(content)
	return extractSelection(lines, cm.selection)
}

// CopySelection copies selected text to system clipboard
func (cm *ClipboardManager) CopySelection(content string) error {
	selected := cm.GetSelection(content)
	if selected == "" {
		return fmt.Errorf("no text selected")
	}

	cm.lastCopy = selected
	return clipboard.WriteAll(selected)
}

// ClearSelection clears the current selection
func (cm *ClipboardManager) ClearSelection() {
	cm.selection.Enabled = false
}

// HasSelection returns whether text is currently selected
func (cm *ClipboardManager) HasSelection() bool {
	return cm.selection.Enabled
}

// GetSelectionBounds returns the selection start and end as (startLine, startCol, endLine, endCol)
func (cm *ClipboardManager) GetSelectionBounds() (int, int, int, int) {
	return cm.selection.StartLine, cm.selection.StartCol,
		cm.selection.EndLine, cm.selection.EndCol
}

// Helper functions

// splitLines splits content into lines
func splitLines(content string) []string {
	// This will be implemented when we need it
	// For now, returning basic split
	lines := make([]string, 0)
	current := ""
	for _, r := range content {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// extractSelection extracts text from lines based on selection state
func extractSelection(lines []string, sel SelectionState) string {
	if !sel.Enabled {
		return ""
	}

	// Normalize selection (start before end)
	startLine := sel.StartLine
	startCol := sel.StartCol
	endLine := sel.EndLine
	endCol := sel.EndCol

	if startLine > endLine || (startLine == endLine && startCol > endCol) {
		startLine, endLine = endLine, startLine
		startCol, endCol = endCol, startCol
	}

	// Handle single-line selection
	if startLine == endLine {
		if startLine >= len(lines) {
			return ""
		}
		line := lines[startLine]
		if endCol > len(line) {
			endCol = len(line)
		}
		if startCol < 0 {
			startCol = 0
		}
		return line[startCol:endCol]
	}

	// Handle multi-line selection
	result := ""

	// First line (partial)
	if startLine < len(lines) {
		line := lines[startLine]
		if startCol < len(line) {
			result += line[startCol:]
		}
		result += "\n"
	}

	// Middle lines (complete)
	for i := startLine + 1; i < endLine && i < len(lines); i++ {
		result += lines[i] + "\n"
	}

	// Last line (partial)
	if endLine < len(lines) {
		line := lines[endLine]
		if endCol > 0 {
			if endCol > len(line) {
				endCol = len(line)
			}
			result += line[:endCol]
		}
	}

	return result
}

// CopyMode represents different copy modes
type CopyMode int

const (
	CharacterMode CopyMode = iota
	LineMode
	BlockMode
)

// KeyboardSelection handles keyboard-based text selection
type KeyboardSelection struct {
	manager *ClipboardManager
	mode    CopyMode
}

// NewKeyboardSelection creates a keyboard selection handler
func NewKeyboardSelection(cm *ClipboardManager) *KeyboardSelection {
	return &KeyboardSelection{
		manager: cm,
		mode:    CharacterMode,
	}
}

// SelectLine selects the current line
func (ks *KeyboardSelection) SelectLine(lineNum int) {
	ks.manager.StartSelection(lineNum, 0, false)
	ks.manager.ExtendSelection(lineNum, 9999) // Extend to end of line
	ks.mode = LineMode
}

// SelectWord selects from current position to end of word
func (ks *KeyboardSelection) SelectWord(line string, startPos int) string {
	// Skip whitespace
	pos := startPos
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	// Select word
	wordStart := pos
	for pos < len(line) && line[pos] != ' ' {
		pos++
	}

	if pos > wordStart {
		return line[wordStart:pos]
	}
	return ""
}
