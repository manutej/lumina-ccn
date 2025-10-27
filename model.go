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
	path  string
	name  string
	isDir bool
	size  int64
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

	// State
	currentView    ViewMode
	selectionMode  SelectionMode // Track selection mode for keyboard-based selection
	selectedFile   string
	markdownFiles  []string
	showHelp       bool // Toggle for help overlay
	selectionStart int  // Character position where selection starts

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

	// NEW: Custom keybindings, clipboard, and color management
	keyBindings     *KeyBindings
	clipboard       *ClipboardManager
	colorManager    *ColorManager
	tableOfContents *TableOfContents // TOC for current file
}

// NewAppModel creates a new application model
func NewAppModel(rootPath string) AppModel {
	// Initialize file list
	items := loadDirectory(rootPath)
	fileList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	fileList.Title = "Files"
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
	toc := NewTableOfContents()

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
		tableOfContents:  toc,
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
			path:  parent,
			name:  "..",
			isDir: true,
			size:  0,
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
				path:  fullPath,
				name:  entry.Name(),
				isDir: entry.IsDir(),
				size:  info.Size(),
			})
		}
	}

	return items
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

	// If it's a markdown file, render it with Glamour and parse TOC
	if strings.HasSuffix(path, ".md") && m.markdownRenderer != nil {
		rendered, err := m.markdownRenderer.Render(m.viewerContent)
		if err != nil {
			// Fallback to plain text if rendering fails
			m.renderedContent = m.viewerContent
		} else {
			m.renderedContent = rendered
		}
		m.viewer.SetContent(m.renderedContent)

		// Parse markdown headers for table of contents
		if m.tableOfContents != nil {
			m.tableOfContents.ParseMarkdown(m.viewerContent)
		}
	} else {
		// For non-markdown files, show plain text
		m.renderedContent = m.viewerContent
		m.viewer.SetContent(m.viewerContent)

		// Clear TOC for non-markdown files
		if m.tableOfContents != nil {
			m.tableOfContents.ParseMarkdown("")
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
		items := loadDirectory(fileItem.path)
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
	items := loadDirectory(parent)
	m.fileList.SetItems(items)
}
