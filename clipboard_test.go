package main

import (
	"strings"
	"testing"

	"github.com/atotto/clipboard"
)

// TestCopySelection_WithoutSelection tests copying entire document when no selection exists
func TestCopySelection_WithoutSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "# Header\n\nThis is a test document.\n\nWith multiple lines."

	// Clear selection to ensure we're testing no-selection case
	cm.ClearSelection()

	// Act: Copy without selection
	err := cm.CopySelection(content)

	// Assert
	if err != nil {
		t.Errorf("CopySelection should not error when copying entire document: %v", err)
	}

	// Verify clipboard content matches full document
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		t.Errorf("Failed to read clipboard: %v", err)
	}

	if clipboardContent != content {
		t.Errorf("Clipboard content mismatch.\nExpected: %q\nGot: %q", content, clipboardContent)
	}

	// Verify internal state
	if cm.lastCopy != content {
		t.Errorf("lastCopy state mismatch.\nExpected: %q\nGot: %q", content, cm.lastCopy)
	}
}

// TestCopySelection_EmptyContent tests copying empty content
func TestCopySelection_EmptyContent(t *testing.T) {
	cm := NewClipboardManager()

	// Act: Try to copy empty content
	err := cm.CopySelection("")

	// Assert: Should return error
	if err == nil {
		t.Error("CopySelection should return error for empty content")
	}

	expectedErr := "no content to copy"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error containing %q, got: %v", expectedErr, err)
	}
}

// TestCopySelection_WithValidSelection tests copying selected text
func TestCopySelection_WithValidSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"

	// Set up selection: Lines 2-3 (indices 1-2)
	cm.StartSelection(1, 0, false)  // Start at line 1, col 0
	cm.ExtendSelection(2, 6)        // End at line 2, col 6

	// Act: Copy selection
	err := cm.CopySelection(content)

	// Assert
	if err != nil {
		t.Errorf("CopySelection should not error with valid selection: %v", err)
	}

	// Verify clipboard contains selected lines
	clipboardContent, _ := clipboard.ReadAll()
	expectedSelection := "Line 2\nLine 3"

	if clipboardContent != expectedSelection {
		t.Errorf("Selection copy mismatch.\nExpected: %q\nGot: %q", expectedSelection, clipboardContent)
	}
}

// TestCopySelection_SingleLineSelection tests copying single line selection
func TestCopySelection_SingleLineSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "Line 1\nLine 2\nLine 3"

	// Select middle word in line 2 (characters 5-10 in "Line 2")
	cm.StartSelection(1, 5, false)  // Start at line 1, col 5 ("2")
	cm.ExtendSelection(1, 6)        // End at line 1, col 6

	// Act
	err := cm.CopySelection(content)

	// Assert
	if err != nil {
		t.Errorf("CopySelection should not error: %v", err)
	}

	clipboardContent, _ := clipboard.ReadAll()
	expectedSelection := "2" // Just the "2" character

	if clipboardContent != expectedSelection {
		t.Errorf("Single line selection mismatch.\nExpected: %q\nGot: %q", expectedSelection, clipboardContent)
	}
}

// TestCopySelection_MultiLineSelection tests copying multi-line selection
func TestCopySelection_MultiLineSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "First line\nSecond line\nThird line\nFourth line"

	// Select from middle of line 1 to middle of line 3
	cm.StartSelection(1, 7, false)  // Start at "line" in "Second line"
	cm.ExtendSelection(2, 6)        // End at "line" in "Third line"

	// Act
	err := cm.CopySelection(content)

	// Assert
	if err != nil {
		t.Errorf("CopySelection should not error: %v", err)
	}

	clipboardContent, _ := clipboard.ReadAll()
	expectedSelection := "line\nThird " // From "line" through "Third " (including space before "line")

	if clipboardContent != expectedSelection {
		t.Errorf("Multi-line selection mismatch.\nExpected: %q\nGot: %q", expectedSelection, clipboardContent)
	}
}

// TestCopySelection_EmptySelection tests copying when selection is empty (same start and end)
func TestCopySelection_EmptySelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "Test content"

	// Start selection but don't extend (same start/end point)
	cm.StartSelection(0, 0, false)
	// Don't call ExtendSelection - selection is empty

	// Act
	err := cm.CopySelection(content)

	// Assert: With empty selection, should copy entire document
	if err != nil {
		t.Errorf("CopySelection should not error with empty selection: %v", err)
	}

	// Should copy entire content since selection is effectively empty
	clipboardContent, _ := clipboard.ReadAll()
	// Current implementation may fail here if it doesn't handle empty selection properly
	// This test documents the expected behavior
	if clipboardContent == "" {
		t.Error("Empty selection should copy entire document, got empty clipboard")
	}
}

// TestGetSelection_NoSelection tests getting selection when none exists
func TestGetSelection_NoSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "Test content"

	// Act: Get selection without setting one
	selected := cm.GetSelection(content)

	// Assert: Should return empty string
	if selected != "" {
		t.Errorf("GetSelection should return empty string when no selection. Got: %q", selected)
	}
}

// TestSelectionBoundaries_StartAfterEnd tests selection with reversed boundaries
func TestSelectionBoundaries_StartAfterEnd(t *testing.T) {
	cm := NewClipboardManager()
	content := "Line 1\nLine 2\nLine 3"

	// Set selection with start after end (should normalize)
	cm.StartSelection(2, 0, false)   // Start at line 2
	cm.ExtendSelection(1, 0)         // End at line 1 (before start)

	// Act
	err := cm.CopySelection(content)

	// Assert: Should handle reversed selection gracefully
	if err != nil {
		t.Errorf("CopySelection should handle reversed selection: %v", err)
	}

	// Should normalize and copy line 1 to line 2
	clipboardContent, _ := clipboard.ReadAll()
	expectedSelection := "Line 2\nLine 3"

	if clipboardContent != expectedSelection {
		t.Errorf("Reversed selection mismatch.\nExpected: %q\nGot: %q", expectedSelection, clipboardContent)
	}
}

// TestSelectionBoundaries_OutOfBounds tests selection beyond document boundaries
func TestSelectionBoundaries_OutOfBounds(t *testing.T) {
	cm := NewClipboardManager()
	content := "Line 1\nLine 2\nLine 3"

	// Select beyond document boundaries
	cm.StartSelection(0, 0, false)
	cm.ExtendSelection(100, 100)  // Way beyond document end

	// Act
	err := cm.CopySelection(content)

	// Assert: Should not crash, should handle gracefully
	if err != nil {
		t.Errorf("CopySelection should handle out-of-bounds selection: %v", err)
	}

	// Should copy up to document end
	clipboardContent, _ := clipboard.ReadAll()
	if clipboardContent == "" {
		t.Error("Out-of-bounds selection should copy available content, got empty")
	}
}

// TestClearSelection tests clearing selection
func TestClearSelection(t *testing.T) {
	cm := NewClipboardManager()
	content := "Test content"

	// Set up selection
	cm.StartSelection(0, 0, false)
	cm.ExtendSelection(0, 5)

	// Verify selection exists
	if !cm.HasSelection() {
		t.Error("Selection should be active after StartSelection/ExtendSelection")
	}

	// Act: Clear selection
	cm.ClearSelection()

	// Assert: Selection should be cleared
	if cm.HasSelection() {
		t.Error("Selection should be cleared after ClearSelection")
	}

	// Copying after clear should copy entire document
	err := cm.CopySelection(content)
	if err != nil {
		t.Errorf("CopySelection after clear should work: %v", err)
	}

	clipboardContent, _ := clipboard.ReadAll()
	if clipboardContent != content {
		t.Errorf("After clear, should copy entire document.\nExpected: %q\nGot: %q", content, clipboardContent)
	}
}

// TestHasSelection tests selection state tracking
func TestHasSelection(t *testing.T) {
	cm := NewClipboardManager()

	// Initially no selection
	if cm.HasSelection() {
		t.Error("New ClipboardManager should have no selection")
	}

	// Start selection
	cm.StartSelection(0, 0, false)
	if !cm.HasSelection() {
		t.Error("HasSelection should return true after StartSelection")
	}

	// Clear selection
	cm.ClearSelection()
	if cm.HasSelection() {
		t.Error("HasSelection should return false after ClearSelection")
	}
}

// TestGetSelectionBounds tests getting selection boundaries
func TestGetSelectionBounds(t *testing.T) {
	cm := NewClipboardManager()

	// Set selection bounds
	startLine, startCol := 5, 10
	endLine, endCol := 8, 15

	cm.StartSelection(startLine, startCol, false)
	cm.ExtendSelection(endLine, endCol)

	// Act: Get bounds
	gotStartLine, gotStartCol, gotEndLine, gotEndCol := cm.GetSelectionBounds()

	// Assert: Should match what we set
	if gotStartLine != startLine || gotStartCol != startCol ||
		gotEndLine != endLine || gotEndCol != endCol {
		t.Errorf("Selection bounds mismatch.\nExpected: (%d,%d) to (%d,%d)\nGot: (%d,%d) to (%d,%d)",
			startLine, startCol, endLine, endCol,
			gotStartLine, gotStartCol, gotEndLine, gotEndCol)
	}
}

// TestSplitLines tests the internal line splitting function
func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single line",
			input:    "Hello",
			expected: []string{"Hello"},
		},
		{
			name:     "Multiple lines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:     "Trailing newline",
			input:    "Line 1\nLine 2\n",
			expected: []string{"Line 1", "Line 2"},
		},
		{
			name:     "Empty lines",
			input:    "Line 1\n\nLine 3",
			expected: []string{"Line 1", "", "Line 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("splitLines() length mismatch.\nExpected %d lines, got %d lines\nExpected: %v\nGot: %v",
					len(tt.expected), len(result), tt.expected, result)
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("splitLines() line %d mismatch.\nExpected: %q\nGot: %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

// TestExtractSelection tests the internal selection extraction function
func TestExtractSelection(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		sel      SelectionState
		expected string
	}{
		{
			name:  "Disabled selection",
			lines: []string{"Line 1", "Line 2"},
			sel:   SelectionState{Enabled: false},
			expected: "",
		},
		{
			name:  "Single character",
			lines: []string{"Hello"},
			sel: SelectionState{
				Enabled:   true,
				StartLine: 0,
				StartCol:  1,
				EndLine:   0,
				EndCol:    2,
			},
			expected: "e", // "Hello"[1:2] = "e"
		},
		{
			name:  "Full line",
			lines: []string{"Line 1", "Line 2"},
			sel: SelectionState{
				Enabled:   true,
				StartLine: 0,
				StartCol:  0,
				EndLine:   0,
				EndCol:    6,
			},
			expected: "Line 1",
		},
		{
			name:  "Multi-line selection",
			lines: []string{"Line 1", "Line 2", "Line 3"},
			sel: SelectionState{
				Enabled:   true,
				StartLine: 0,
				StartCol:  7, // Start after "Line 1"
				EndLine:   2,
				EndCol:    4, // End at "Line" in "Line 3"
			},
			expected: "\nLine 2\nLine", // From end of Line 1 through "Line" in Line 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSelection(tt.lines, tt.sel)

			if result != tt.expected {
				t.Errorf("extractSelection() mismatch.\nExpected: %q\nGot: %q", tt.expected, result)
			}
		})
	}
}

// TestCopySelection_SpecialCharacters tests copying content with special characters
func TestCopySelection_SpecialCharacters(t *testing.T) {
	cm := NewClipboardManager()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Emojis",
			content: "Hello 👋 World 🌍",
		},
		{
			name:    "Unicode",
			content: "日本語 テスト",
		},
		{
			name:    "Code block with backticks",
			content: "```go\nfunc main() {}\n```",
		},
		{
			name:    "Special symbols",
			content: "© ® ™ € £ ¥ ¢",
		},
		{
			name:    "Tab characters",
			content: "Line\twith\ttabs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear selection to copy entire content
			cm.ClearSelection()

			err := cm.CopySelection(tt.content)
			if err != nil {
				t.Errorf("CopySelection should handle special characters: %v", err)
			}

			clipboardContent, _ := clipboard.ReadAll()
			if clipboardContent != tt.content {
				t.Errorf("Special character handling failed.\nExpected: %q\nGot: %q", tt.content, clipboardContent)
			}
		})
	}
}

// TestCopySelection_LargeContent tests performance with large content
func TestCopySelection_LargeContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	cm := NewClipboardManager()

	// Generate 10,000 lines of content
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("This is line number ")
		sb.WriteString(strings.Repeat("x", 100)) // Make each line 100+ chars
		sb.WriteString("\n")
	}
	largeContent := sb.String()

	// Act: Copy large content
	err := cm.CopySelection(largeContent)

	// Assert: Should complete without error
	if err != nil {
		t.Errorf("CopySelection should handle large content: %v", err)
	}

	// Verify clipboard has the content
	clipboardContent, _ := clipboard.ReadAll()
	if len(clipboardContent) != len(largeContent) {
		t.Errorf("Large content copy size mismatch.\nExpected: %d bytes\nGot: %d bytes",
			len(largeContent), len(clipboardContent))
	}
}

// Benchmark tests

func BenchmarkCopySelection_SmallDocument(b *testing.B) {
	cm := NewClipboardManager()
	content := strings.Repeat("Line of text\n", 100) // 100 lines

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cm.CopySelection(content)
	}
}

func BenchmarkCopySelection_LargeDocument(b *testing.B) {
	cm := NewClipboardManager()
	content := strings.Repeat("Line of text\n", 10000) // 10,000 lines

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cm.CopySelection(content)
	}
}

func BenchmarkCopySelection_WithSelection(b *testing.B) {
	cm := NewClipboardManager()
	content := strings.Repeat("Line of text\n", 1000)

	// Set up a selection
	cm.StartSelection(10, 0, false)
	cm.ExtendSelection(20, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cm.CopySelection(content)
	}
}
