package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SearchResult represents a single search result
type SearchResult struct {
	Path      string // Full path to file
	Name      string // File name only
	Score     int    // Fuzzy match score (higher = better)
	LineNum   int    // For content search (future)
	Highlight string // For display
}

// FileSearcher performs fast file search
type FileSearcher struct {
	rootPath    string
	allFiles    []string        // Cached file list
	results     []SearchResult
	selected    int
	searchTerm  string
	isSearching bool
	mu          sync.RWMutex
}

// NewFileSearcher creates a new file searcher
func NewFileSearcher(rootPath string) *FileSearcher {
	return &FileSearcher{
		rootPath: rootPath,
		allFiles: []string{},
		results:  []SearchResult{},
		selected: 0,
	}
}

// BuildFileIndex builds a complete index of all files (one-time, then cached)
func (fs *FileSearcher) BuildFileIndex() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.allFiles = []string{}

	err := filepath.Walk(fs.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip inaccessible files
		}

		// Skip hidden directories and common ignore patterns
		if info.IsDir() {
			if shouldIgnoreDir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Add file to index
		relPath, _ := filepath.Rel(fs.rootPath, path)
		fs.allFiles = append(fs.allFiles, relPath)

		return nil
	})

	return err
}

// Search performs a fuzzy search on the file index
func (fs *FileSearcher) Search(term string) []SearchResult {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if term == "" {
		return []SearchResult{}
	}

	fs.searchTerm = term
	fs.results = []SearchResult{}
	fs.selected = 0

	// Fuzzy match against all files
	for _, file := range fs.allFiles {
		score := fuzzyMatch(file, term)
		if score > 0 {
			fileName := filepath.Base(file)
			fs.results = append(fs.results, SearchResult{
				Path:    file,
				Name:    fileName,
				Score:   score,
				Highlight: highlightMatch(file, term),
			})
		}
	}

	// Sort by score (highest first)
	sortByScore(fs.results)

	return fs.results
}

// GetResults returns current search results
func (fs *FileSearcher) GetResults() []SearchResult {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.results
}

// SelectNext moves to next result
func (fs *FileSearcher) SelectNext() {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.selected < len(fs.results)-1 {
		fs.selected++
	}
}

// SelectPrevious moves to previous result
func (fs *FileSearcher) SelectPrevious() {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.selected > 0 {
		fs.selected--
	}
}

// GetSelectedResult returns the currently selected result
func (fs *FileSearcher) GetSelectedResult() *SearchResult {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if fs.selected >= 0 && fs.selected < len(fs.results) {
		return &fs.results[fs.selected]
	}
	return nil
}

// GetSelectedIndex returns current selection index
func (fs *FileSearcher) GetSelectedIndex() int {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.selected
}

// GetResultCount returns number of results
func (fs *FileSearcher) GetResultCount() int {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return len(fs.results)
}

// View renders the search results as a string
func (fs *FileSearcher) View(width int, height int) string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if len(fs.results) == 0 {
		return "🔍 Search Results\n\nNo matches found"
	}

	lines := []string{"🔍 Search Results"}
	lines = append(lines, "")

	// Calculate how many results to show
	maxResults := height - 4
	startIdx := 0

	if fs.selected >= maxResults {
		startIdx = fs.selected - maxResults + 1
	}

	endIdx := startIdx + maxResults
	if endIdx > len(fs.results) {
		endIdx = len(fs.results)
	}

	for i := startIdx; i < endIdx; i++ {
		result := fs.results[i]
		prefix := "  "
		if i == fs.selected {
			prefix = "→ "
		}

		// Truncate to fit width
		display := result.Highlight
		if len(display) > width-4 {
			display = display[:width-5] + "…"
		}

		lines = append(lines, prefix+display)
	}

	// Status line
	if len(fs.results) > 0 {
		lines = append(lines, "")
		lines = append(lines, "["+fs.searchTerm+"] ↑/↓ navigate | Enter: open | Esc: cancel")
		lines = append(lines, "Results: "+string(rune(len(fs.results))))
	}

	return strings.Join(lines, "\n")
}

// Helper functions

// shouldIgnoreDir returns true if directory should be skipped
func shouldIgnoreDir(path string) bool {
	name := filepath.Base(path)

	ignorePatterns := []string{
		".git", ".hg", ".svn",
		"node_modules", ".node_modules",
		"vendor", ".vendor",
		"build", "dist", ".build",
		"__pycache__", ".pytest_cache",
		".venv", "venv",
		".DS_Store",
		".cache",
		"target", ".cargo",
	}

	for _, pattern := range ignorePatterns {
		if name == pattern {
			return true
		}
	}

	return false
}

// fuzzyMatch performs fuzzy matching and returns a score
// Higher score = better match
func fuzzyMatch(text, pattern string) int {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)

	if pattern == "" {
		return 0
	}

	if !strings.Contains(text, pattern) {
		return 0
	}

	score := 0
	textIdx := 0
	patternIdx := 0

	// Bonus for matching at word boundaries
	for patternIdx < len(pattern) && textIdx < len(text) {
		if text[textIdx] == pattern[patternIdx] {
			// Exact character match
			score += 10

			// Bonus if it's at the start of a word or filename
			if textIdx == 0 || text[textIdx-1] == '/' || text[textIdx-1] == '_' || text[textIdx-1] == '-' {
				score += 15
			}

			// Bonus if it's immediately after previous match
			if patternIdx > 0 && textIdx > 0 && text[textIdx-1] == pattern[patternIdx-1] {
				score += 5
			}

			patternIdx++
		}
		textIdx++
	}

	// Exact match bonus
	if text == pattern {
		score += 1000
	}

	// Match in filename bonus
	filename := strings.ToLower(filepath.Base(text))
	if filename == pattern {
		score += 500
	}

	// Early match bonus (matching at start)
	if strings.HasPrefix(text, pattern) {
		score += 100
	}

	// Penalize matches that are incomplete
	if patternIdx < len(pattern) {
		return 0 // Pattern not fully matched
	}

	return score
}

// sortByScore sorts results by fuzzy match score (highest first)
func sortByScore(results []SearchResult) {
	// Simple bubble sort for small result sets
	// In production, use more efficient algorithm for large sets
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

// highlightMatch creates a highlighted display string
func highlightMatch(text, pattern string) string {
	// Simple implementation - just show the path
	// Advanced: could use ANSI color codes for highlighting matches
	return text
}

// QuickSearch is a simpler, faster search using string containment
type QuickSearch struct {
	term    string
	results []string
	index   int
}

// NewQuickSearch creates a quick search for immediate results
func NewQuickSearch() *QuickSearch {
	return &QuickSearch{
		results: []string{},
		index:   0,
	}
}

// QuickSearchFiles performs a simple contains-based search
// Useful for immediate feedback while building full index
func (qs *QuickSearch) QuickSearchFiles(files []string, term string) []string {
	term = strings.ToLower(term)
	qs.term = term
	qs.results = []string{}
	qs.index = 0

	for _, file := range files {
		if strings.Contains(strings.ToLower(file), term) {
			qs.results = append(qs.results, file)
		}
	}

	return qs.results
}
