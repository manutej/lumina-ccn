package main

import (
	"regexp"
	"strings"
)

// TOCEntry represents a heading in the table of contents
type TOCEntry struct {
	Level   int    // 1-6 (# to ######)
	Title   string // The heading text
	LineNum int    // Line number in document (for jumping)
}

// TableOfContents represents the complete TOC structure
type TableOfContents struct {
	entries  []TOCEntry
	selected int // Currently selected entry (for navigation)
}

// NewTableOfContents creates a new TOC
func NewTableOfContents() *TableOfContents {
	return &TableOfContents{
		entries:  []TOCEntry{},
		selected: 0,
	}
}

// ParseMarkdown extracts headers from markdown content
func (toc *TableOfContents) ParseMarkdown(content string) {
	toc.entries = []TOCEntry{}
	toc.selected = 0

	lines := strings.Split(content, "\n")
	headingRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	for lineNum, line := range lines {
		matches := headingRegex.FindStringSubmatch(line)
		if matches != nil {
			level := len(matches[1]) // Number of # symbols
			title := strings.TrimSpace(matches[2])

			// Skip if title is empty
			if title == "" {
				continue
			}

			toc.entries = append(toc.entries, TOCEntry{
				Level:   level,
				Title:   title,
				LineNum: lineNum,
			})
		}
	}
}

// View renders the table of contents
func (toc *TableOfContents) View(width int, height int) string {
	if len(toc.entries) == 0 {
		return "📋 Table of Contents\n\n(No headings found)"
	}

	// Build the TOC display
	lines := []string{"📋 Table of Contents\n"}

	// Calculate how many entries we can show
	maxEntries := height - 3 // Leave room for header and status

	// Determine start index (keep selected item visible)
	startIdx := 0
	if toc.selected >= maxEntries {
		startIdx = toc.selected - maxEntries + 1
	}

	endIdx := startIdx + maxEntries
	if endIdx > len(toc.entries) {
		endIdx = len(toc.entries)
	}

	for i := startIdx; i < endIdx; i++ {
		entry := toc.entries[i]

		// Indent based on level (2 spaces per level)
		indent := strings.Repeat("  ", entry.Level-1)

		// Create the display string
		prefix := "  "
		if i == toc.selected {
			prefix = "→ " // Selected indicator (2 chars)
		}

		// Calculate available space for title
		// prefix (2) + indent + title should fit in width
		availableSpace := width - len(prefix) - len(indent) - 1 // -1 for margin

		title := entry.Title
		// Only truncate if we absolutely must
		if availableSpace > 4 && len(title) > availableSpace {
			// Leave room for ellipsis
			maxLen := availableSpace - 1
			if maxLen > 0 {
				title = title[:maxLen] + "…"
			}
		} else if availableSpace <= 4 {
			// Very narrow space - try to show at least some chars
			if len(title) > availableSpace && availableSpace > 1 {
				title = title[:availableSpace-1] + "…"
			}
		}

		lines = append(lines, prefix+indent+title)
	}

	// Add navigation hint
	if len(toc.entries) > 1 {
		lines = append(lines, "\n(↑/↓ to navigate, Enter to jump)")
	}

	return strings.Join(lines, "\n")
}

// SelectNext moves to the next TOC entry
func (toc *TableOfContents) SelectNext() {
	if toc.selected < len(toc.entries)-1 {
		toc.selected++
	}
}

// SelectPrevious moves to the previous TOC entry
func (toc *TableOfContents) SelectPrevious() {
	if toc.selected > 0 {
		toc.selected--
	}
}

// SelectFirst jumps to the first TOC entry
func (toc *TableOfContents) SelectFirst() {
	toc.selected = 0
}

// SelectLast jumps to the last TOC entry
func (toc *TableOfContents) SelectLast() {
	if len(toc.entries) > 0 {
		toc.selected = len(toc.entries) - 1
	}
}

// GetSelectedEntry returns the currently selected TOC entry
func (toc *TableOfContents) GetSelectedEntry() *TOCEntry {
	if toc.selected >= 0 && toc.selected < len(toc.entries) {
		return &toc.entries[toc.selected]
	}
	return nil
}

// GetSelectedLineNum returns the line number of the selected entry
func (toc *TableOfContents) GetSelectedLineNum() int {
	if entry := toc.GetSelectedEntry(); entry != nil {
		return entry.LineNum
	}
	return 0
}

// HasEntries returns whether the TOC has any entries
func (toc *TableOfContents) HasEntries() bool {
	return len(toc.entries) > 0
}

// GetEntryCount returns the number of TOC entries
func (toc *TableOfContents) GetEntryCount() int {
	return len(toc.entries)
}

// GetSelectedIndex returns the current selection index
func (toc *TableOfContents) GetSelectedIndex() int {
	return toc.selected
}

// SelectByIndex selects a specific entry by index
func (toc *TableOfContents) SelectByIndex(index int) {
	if index >= 0 && index < len(toc.entries) {
		toc.selected = index
	}
}

// SelectByLineNum selects the entry closest to a given line number
// Useful for syncing viewer scroll with TOC
func (toc *TableOfContents) SelectByLineNum(lineNum int) {
	closestIdx := 0
	closestDiff := 999999

	for i, entry := range toc.entries {
		diff := abs(entry.LineNum - lineNum)
		if diff < closestDiff {
			closestDiff = diff
			closestIdx = i
		}
	}

	toc.selected = closestIdx
}

// JumpToLevel shows entries of a specific level or higher (collapsing deeper levels)
func (toc *TableOfContents) FilterByLevel(maxLevel int) []TOCEntry {
	result := []TOCEntry{}
	for _, entry := range toc.entries {
		if entry.Level <= maxLevel {
			result = append(result, entry)
		}
	}
	return result
}

// GetHierarchy returns a tree-like structure of the TOC
func (toc *TableOfContents) GetHierarchy() string {
	if len(toc.entries) == 0 {
		return "No headings found"
	}

	var result []string
	for _, entry := range toc.entries {
		indent := strings.Repeat("  ", entry.Level-1)
		result = append(result, indent+"- "+entry.Title)
	}

	return strings.Join(result, "\n")
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SearchTOC finds entries matching a search term
func (toc *TableOfContents) SearchTOC(searchTerm string) []TOCEntry {
	searchTerm = strings.ToLower(searchTerm)
	var matches []TOCEntry

	for _, entry := range toc.entries {
		if strings.Contains(strings.ToLower(entry.Title), searchTerm) {
			matches = append(matches, entry)
		}
	}

	return matches
}
