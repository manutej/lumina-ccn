package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/lumina/ccn/utils"
)

// ViewMode represents the current view/pane focus
type ViewMode int

const (
	FileTreeView ViewMode = iota
	ViewerView
	PreviewView
)

// SortMode represents the file sorting mode
type SortMode int

const (
	SortAlphabetical SortMode = iota
	SortModifiedTime
	SortSize
)

// UIMode represents the current UI interaction mode (State Machine)
type UIMode int

const (
	NormalMode UIMode = iota
	FinderMode
	SearchMode
	HelpMode
	LoadingMode
)

// SelectionMode represents text selection state in viewer
type SelectionMode int

const (
	SelectionInactive SelectionMode = iota
	SelectionCharacter
	SelectionLine
	SelectionBlock
)

// FileItem represents a file or directory in the file tree
type FileItem struct {
	path    string
	name    string
	isDir   bool
	size    int64
	modTime int64 // Unix timestamp for modification time
}

// Implement list.Item interface for FileItem
func (f FileItem) FilterValue() string { return f.name }
func (f FileItem) Title() string       { return f.name }
func (f FileItem) Description() string {
	if f.name == ".." {
		return "⬆️  Parent directory"
	}
	if f.isDir {
		return "📁 Directory"
	}
	return "📄 File"
}

// AppModel represents the complete application state
type AppModel struct {
	// Window dimensions
	width  int
	height int
	ready  bool

	// Current working directory
	rootPath    string
	currentPath string

	// UI Components
	fileList         list.Model
	viewer           viewport.Model
	preview          viewport.Model
	markdownRenderer *utils.MarkdownRenderer

	// State Machine - Only ONE state active at a time
	currentMode UIMode

	// View State
	currentView    ViewMode
	selectionMode  SelectionMode // Track selection mode for keyboard-based selection
	selectedFile   string
	markdownFiles  []string
	selectionStart int // Character position where selection starts

	// Mouse drag state
	mouseDragActive    bool // Whether mouse drag is currently happening
	mouseDragStartLine int  // Where the drag started (line)
	mouseDragStartCol  int  // Where the drag started (column)

	// Content
	viewerContent   string
	renderedContent string

	// Layout dimensions (percentages)
	fileTreeWidth int // 20% of width
	viewerWidth   int // 60% of width
	previewWidth  int // 20% of width

	// Custom keybindings, clipboard, and color management
	keyBindings  *KeyBindings
	clipboard    *ClipboardManager
	colorManager *ColorManager
	contextPanel *ContextPanel // Multi-mode right panel (TOC, File Info, Stats, etc.)

	// Phase 3: Fuzzy Finder State (flat, no helper object)
	finderInput    string
	finderCursor   int
	finderItems    []string
	finderFiltered []string

	// Phase 3: Ripgrep Search State (flat)
	searchQuery      string
	searchResults    []RipgrepResult
	searchCursor     int
	searchInProgress bool

	// Phase 3: File Watcher State
	fileWatcher             *FileWatcher
	watcherActive           bool
	fileChangedNotification bool
	loadingMessage          string

	// Sorting
	currentSortMode SortMode
}

// NewAppModel creates a new application model
func NewAppModel(rootPath string) AppModel {
	// Initialize file list with default sort mode
	items := loadDirectorySorted(rootPath, SortAlphabetical)
	fileList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	fileList.Title = "Files [Alphabetical]"
	fileList.SetShowStatusBar(false)
	fileList.SetFilteringEnabled(true)

	// Initialize viewports
	viewer := viewport.New(0, 0)
	preview := viewport.New(0, 0)

	// Initialize markdown renderer (will be updated with actual width later)
	renderer, _ := utils.NewMarkdownRenderer(80)

	// NEW: Initialize keybindings, clipboard, and color manager
	keyBindings := LoadKeyBindings()
	clipboard := NewClipboardManager()
	colorManager, _ := NewColorManager()
	contextPanel := NewContextPanel()

	m := AppModel{
		rootPath:         rootPath,
		currentPath:      rootPath,
		fileList:         fileList,
		viewer:           viewer,
		preview:          preview,
		markdownRenderer: renderer,
		currentView:      FileTreeView,
		markdownFiles:    []string{}, // Lazy-load on first use (for global search)
		keyBindings:      &keyBindings,
		clipboard:        clipboard,
		colorManager:     colorManager,
		contextPanel:     contextPanel,
		// Phase 3: State Machine
		currentMode: NormalMode,
		// Phase 3: Fuzzy Finder State (flat)
		finderInput:    "",
		finderCursor:   0,
		finderItems:    []string{},
		finderFiltered: []string{},
		// Phase 3: Search State
		searchQuery:      "",
		searchResults:    []RipgrepResult{},
		searchCursor:     0,
		searchInProgress: false,
		// Phase 3: File Watcher
		fileWatcher:             nil,
		watcherActive:           false,
		fileChangedNotification: false,
		loadingMessage:          "",
		// Sorting
		currentSortMode: SortAlphabetical,
		// Set default dimensions for immediate display
		width:  120,  // Default width
		height: 40,   // Default height
		ready:  true, // Mark ready immediately
	}

	// Apply default dimensions
	m.updateDimensions()

	return m
}

// loadDirectory loads files from a directory
func loadDirectory(path string) []list.Item {
	var items []list.Item

	// Add parent directory entry ".." if not at filesystem root
	parent := filepath.Dir(path)
	if parent != path {
		items = append(items, FileItem{
			path:    parent,
			name:    "..",
			isDir:   true,
			size:    0,
			modTime: 0,
		})
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		// Skip hidden files (but not ".." which we already added)
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		// Only show markdown files and directories
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".md") {
			items = append(items, FileItem{
				path:    fullPath,
				name:    entry.Name(),
				isDir:   entry.IsDir(),
				size:    info.Size(),
				modTime: info.ModTime().Unix(),
			})
		}
	}

	return items
}

// loadDirectorySorted loads files from a directory and applies sorting
func loadDirectorySorted(path string, sortMode SortMode) []list.Item {
	items := loadDirectory(path)
	return sortFileItems(items, sortMode)
}

// sortFileItems sorts a list of FileItems based on the specified sort mode
// Returns the sorted slice (needed because slice modifications don't persist with local reassignment)
func sortFileItems(items []list.Item, sortMode SortMode) []list.Item {
	if len(items) == 0 {
		return items
	}

	// Keep ".." parent directory at the top
	parentIdx := -1
	for i, item := range items {
		if fileItem, ok := item.(FileItem); ok && fileItem.name == ".." {
			parentIdx = i
			break
		}
	}

	// Extract the parent item if it exists
	var parentItem list.Item
	var itemsToSort []list.Item
	if parentIdx >= 0 {
		parentItem = items[parentIdx]
		// Create new slice without parent
		itemsToSort = make([]list.Item, 0, len(items)-1)
		itemsToSort = append(itemsToSort, items[:parentIdx]...)
		itemsToSort = append(itemsToSort, items[parentIdx+1:]...)
	} else {
		itemsToSort = items
	}

	// Sort based on mode
	switch sortMode {
	case SortAlphabetical:
		sortAlphabetically(itemsToSort)
	case SortModifiedTime:
		sortByModifiedTime(itemsToSort)
	case SortSize:
		sortBySize(itemsToSort)
	}

	// Re-insert parent item at the top
	if parentIdx >= 0 {
		result := make([]list.Item, 0, len(items))
		result = append(result, parentItem)
		result = append(result, itemsToSort...)
		return result
	}

	return itemsToSort
}

// sortAlphabetically sorts items alphabetically (directories first, then files)
func sortAlphabetically(items []list.Item) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			item1, ok1 := items[i].(FileItem)
			item2, ok2 := items[j].(FileItem)
			if !ok1 || !ok2 {
				continue
			}

			// Directories before files
			if item1.isDir != item2.isDir {
				if item2.isDir {
					items[i], items[j] = items[j], items[i]
				}
				continue
			}

			// Alphabetical within same type
			if strings.ToLower(item2.name) < strings.ToLower(item1.name) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// sortByModifiedTime sorts items by modification time (newest first, directories first)
func sortByModifiedTime(items []list.Item) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			item1, ok1 := items[i].(FileItem)
			item2, ok2 := items[j].(FileItem)
			if !ok1 || !ok2 {
				continue
			}

			// Directories before files
			if item1.isDir != item2.isDir {
				if item2.isDir {
					items[i], items[j] = items[j], items[i]
				}
				continue
			}

			// Newer items first (higher timestamp = more recent)
			if item2.modTime > item1.modTime {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// sortBySize sorts items by size (largest first, directories first)
func sortBySize(items []list.Item) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			item1, ok1 := items[i].(FileItem)
			item2, ok2 := items[j].(FileItem)
			if !ok1 || !ok2 {
				continue
			}

			// Directories before files
			if item1.isDir != item2.isDir {
				if item2.isDir {
					items[i], items[j] = items[j], items[i]
				}
				continue
			}

			// Larger files first
			if item2.size > item1.size {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// getSortModeName returns a human-readable name for the sort mode
func getSortModeName(sortMode SortMode) string {
	switch sortMode {
	case SortAlphabetical:
		return "Alphabetical"
	case SortModifiedTime:
		return "Modified Time"
	case SortSize:
		return "Size"
	default:
		return "Unknown"
	}
}

// cycleSortMode cycles to the next sort mode
func (m *AppModel) cycleSortMode() {
	m.currentSortMode = (m.currentSortMode + 1) % 3
	// Reload directory with new sort mode
	items := loadDirectorySorted(m.currentPath, m.currentSortMode)
	m.fileList.SetItems(items)
	// Update title to show current sort mode
	m.fileList.Title = "Files [" + getSortModeName(m.currentSortMode) + "]"
}

// findMarkdownFiles recursively finds all markdown files in a directory
func findMarkdownFiles(rootPath string) []string {
	var files []string

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		// Add markdown files
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			files = append(files, path)
		}

		return nil
	})

	return files
}

// updateDimensions recalculates layout dimensions based on window size
func (m *AppModel) updateDimensions() {
	if !m.ready {
		return
	}

	// Calculate pane widths (with borders)
	m.fileTreeWidth = int(float64(m.width) * 0.20)
	m.viewerWidth = int(float64(m.width) * 0.60)
	m.previewWidth = m.width - m.fileTreeWidth - m.viewerWidth - 4 // account for borders

	// Update component dimensions
	headerHeight := 3
	statusBarHeight := 2
	contentHeight := m.height - headerHeight - statusBarHeight

	m.fileList.SetSize(m.fileTreeWidth-2, contentHeight)
	m.viewer.Width = m.viewerWidth - 2
	m.viewer.Height = contentHeight
	m.preview.Width = m.previewWidth - 2
	m.preview.Height = contentHeight

	// Update markdown renderer width
	if m.markdownRenderer != nil {
		m.markdownRenderer.UpdateWidth(m.viewerWidth - 4)
	}

	// Re-render current content if a file is loaded
	if m.selectedFile != "" && strings.HasSuffix(m.selectedFile, ".md") {
		if rendered, err := m.markdownRenderer.Render(m.viewerContent); err == nil {
			m.renderedContent = rendered
			m.viewer.SetContent(m.renderedContent)
		}
	}
}

// loadFileContent loads the content of a file
func (m *AppModel) loadFileContent(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	m.viewerContent = string(content)
	m.selectedFile = path

	// If it's a markdown file, render it with Glamour and update context panel
	if strings.HasSuffix(path, ".md") && m.markdownRenderer != nil {
		rendered, err := m.markdownRenderer.Render(m.viewerContent)
		if err != nil {
			// Fallback to plain text if rendering fails
			m.renderedContent = m.viewerContent
		} else {
			m.renderedContent = rendered
		}
		m.viewer.SetContent(m.renderedContent)

		// Update context panel with file analysis (TOC, stats, metadata)
		if m.contextPanel != nil {
			m.contextPanel.UpdateContent(path, m.viewerContent)
		}
	} else {
		// For non-markdown files, show plain text
		m.renderedContent = m.viewerContent
		m.viewer.SetContent(m.viewerContent)

		// Clear context panel for non-markdown files
		if m.contextPanel != nil {
			m.contextPanel.Clear()
		}
	}

	m.viewer.GotoTop()
	return nil
}

// navigateToSelectedFile navigates to the currently selected file
func (m *AppModel) navigateToSelectedFile() error {
	item := m.fileList.SelectedItem()
	if item == nil {
		return nil
	}

	fileItem := item.(FileItem)

	if fileItem.isDir {
		// Navigate into directory
		m.currentPath = fileItem.path
		items := loadDirectorySorted(fileItem.path, m.currentSortMode)
		m.fileList.SetItems(items)
		return nil
	}

	// Load file content
	return m.loadFileContent(fileItem.path)
}

// navigateUp navigates to the parent directory
func (m *AppModel) navigateUp() {
	parent := filepath.Dir(m.currentPath)

	// Only block if we're already at filesystem root
	if parent == m.currentPath {
		return // Already at filesystem root (e.g., "/" on Unix)
	}

	m.currentPath = parent
	items := loadDirectorySorted(parent, m.currentSortMode)
	m.fileList.SetItems(items)
}

// jumpToNextFileStartingWith jumps to the next file/directory starting with the given letter
// This enables quick navigation via Shift+Letter (e.g., Shift+R to jump to files starting with 'R')
func (m *AppModel) jumpToNextFileStartingWith(letter string) {
	items := m.fileList.Items()
	if len(items) == 0 {
		return
	}

	// Get current index
	currentIndex := m.fileList.Index()

	// Convert search letter to lowercase for case-insensitive matching
	searchLetter := strings.ToLower(letter)

	// Search from current+1 to end
	for i := currentIndex + 1; i < len(items); i++ {
		if item, ok := items[i].(FileItem); ok {
			// Skip ".." parent directory entry
			if item.name == ".." {
				continue
			}
			// Check if name starts with the letter (case-insensitive)
			if strings.HasPrefix(strings.ToLower(item.name), searchLetter) {
				m.fileList.Select(i)
				return
			}
		}
	}

	// Wrap around: search from beginning to current
	for i := 0; i <= currentIndex; i++ {
		if item, ok := items[i].(FileItem); ok {
			// Skip ".." parent directory entry
			if item.name == ".." {
				continue
			}
			// Check if name starts with the letter (case-insensitive)
			if strings.HasPrefix(strings.ToLower(item.name), searchLetter) {
				m.fileList.Select(i)
				return
			}
		}
	}

	// If no match found, do nothing (stay at current position)
}

// RipgrepResult represents a single search result from ripgrep
type RipgrepResult struct {
	FilePath string
	Line     int
	Column   int
	Text     string
}

// State Transition Functions (Blocker 4 Fix)

// transitionTo cleanly transitions between UI modes
func (m *AppModel) transitionTo(newMode UIMode) {
	// Clean up old state
	switch m.currentMode {
	case FinderMode:
		m.finderInput = ""
		m.finderCursor = 0
		m.finderFiltered = []string{}
	case SearchMode:
		m.searchQuery = ""
		m.searchCursor = 0
	case HelpMode:
		// No cleanup needed
	case LoadingMode:
		m.loadingMessage = ""
	}

	m.currentMode = newMode
}

// Pure Helper Functions (Blocker 2 Fix)

// filterItems performs case-insensitive substring filtering
func filterItems(items []string, query string) []string {
	if query == "" {
		return items
	}

	var filtered []string
	lowerQuery := strings.ToLower(query)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), lowerQuery) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// navigateCursor wraps cursor navigation with boundary checking
func navigateCursor(cursor, delta, listLen int) int {
	if listLen == 0 {
		return 0
	}
	newCursor := cursor + delta
	if newCursor < 0 {
		return listLen - 1 // Wrap to end
	}
	if newCursor >= listLen {
		return 0 // Wrap to start
	}
	return newCursor
}
