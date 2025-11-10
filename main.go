package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Phase 3 Message Types

// Async File Loading (Blocker 1 Fix)
type MarkdownFilesLoadingMsg struct{}
type MarkdownFilesLoadedMsg struct {
	files []string
}

// Mouse Drag Timeout (Blocker 3 Fix)
type MouseDragTimeoutMsg struct{}

// Fuzzy Finder Messages
type FinderActivatedMsg struct{}
type FinderInputMsg struct {
	input string
}
type FinderSelectionMsg struct {
	filePath string
}
type FinderCanceledMsg struct{}

// Ripgrep Search Messages (Phase 3 Week 2)
type SearchStartedMsg struct {
	query string
}
type SearchResultMsg struct {
	result RipgrepResult
}
type SearchCompletedMsg struct{}
type SearchErrorMsg struct {
	err error
}

// File Watcher Messages (Phase 3 Week 3)
type FileChangedMsg struct {
	path string
}
type FileReloadMsg struct {
	content string
}

// Init initializes the model
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Command Functions (Blocker 1 & 3 Fixes)

// loadMarkdownFilesCmd asynchronously loads markdown files (Blocker 1 Fix)
func loadMarkdownFilesCmd(rootPath string) tea.Cmd {
	return func() tea.Msg {
		files := findMarkdownFiles(rootPath)
		return MarkdownFilesLoadedMsg{files: files}
	}
}

// dragTimeoutCmd creates a timeout command for mouse drag recovery (Blocker 3 Fix)
func dragTimeoutCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return MouseDragTimeoutMsg{}
	})
}

// Update handles messages and updates the model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateDimensions()
		return m, nil

	case tea.MouseMsg:
		// Handle mouse events
		return m.handleMouseEvent(msg)

	// Handle async message types
	case MarkdownFilesLoadedMsg:
		m.markdownFiles = msg.files
		m.finderItems = msg.files
		m.finderFiltered = msg.files
		m.transitionTo(FinderMode)
		return m, nil

	case MouseDragTimeoutMsg:
		// Auto-recover from stuck drag state (Blocker 3 Fix)
		if m.mouseDragActive {
			m.mouseDragActive = false
			m.clipboard.ClearSelection()
		}
		return m, nil

	case SearchResultMsg:
		// Accumulate search results
		m.searchResults = append(m.searchResults, msg.result)
		return m, nil

	case SearchCompletedMsg:
		m.searchInProgress = false
		return m, nil

	case FileChangedMsg:
		// Reload file when changed
		if m.selectedFile == msg.path {
			m.loadFileContent(msg.path)
		}
		return m, nil

	case tea.KeyMsg:
		// State Machine: Route based on current mode (Blocker 4 Fix)
		switch m.currentMode {
		case HelpMode:
			return m.handleHelpMode(msg)
		case FinderMode:
			return m.handleFinderMode(msg)
		case SearchMode:
			return m.handleSearchMode(msg)
		case LoadingMode:
			return m.handleLoadingMode(msg)
		case TOCMode:
			return m.handleTOCMode(msg)
		case NormalMode:
			return m.handleNormalMode(msg)
		}
	}

	// Update components based on current view
	if m.currentView == FileTreeView {
		m.fileList, cmd = m.fileList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Mode Handlers (State Machine Implementation)

func (m AppModel) handleHelpMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" || msg.String() == "?" {
		m.transitionTo(NormalMode)
		return m, nil
	}
	// Allow quit from help
	if msg.String() == "q" || msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	// Ignore other keys in help mode
	return m, nil
}

func (m AppModel) handleFinderMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel finder
		m.transitionTo(NormalMode)
		return m, nil

	case "enter":
		// Select item and load file
		if len(m.finderFiltered) > 0 && m.finderCursor < len(m.finderFiltered) {
			selected := m.finderFiltered[m.finderCursor]
			m.transitionTo(NormalMode)
			// Load the selected file
			if err := m.loadFileContent(selected); err == nil {
				// Successfully loaded file
			}
		}
		return m, nil

	case "up", "ctrl+p":
		// Navigate up in results (Blocker 2: Pure function)
		m.finderCursor = navigateCursor(m.finderCursor, -1, len(m.finderFiltered))
		return m, nil

	case "down", "ctrl+n":
		// Navigate down in results (Blocker 2: Pure function)
		m.finderCursor = navigateCursor(m.finderCursor, 1, len(m.finderFiltered))
		return m, nil

	case "backspace":
		// Delete character from input
		if len(m.finderInput) > 0 {
			m.finderInput = m.finderInput[:len(m.finderInput)-1]
			m.finderFiltered = filterItems(m.finderItems, m.finderInput) // Blocker 2: Pure function
			m.finderCursor = 0 // Reset cursor to top
		}
		return m, nil

	default:
		// Add character to input if it's a single printable character
		if len(msg.String()) == 1 {
			m.finderInput += msg.String()
			m.finderFiltered = filterItems(m.finderItems, m.finderInput) // Blocker 2: Pure function
			m.finderCursor = 0 // Reset cursor to top
		}
		return m, nil
	}
}

func (m AppModel) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.transitionTo(NormalMode)
		return m, nil

	case "enter":
		// Jump to selected result
		if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
			result := m.searchResults[m.searchCursor]
			m.loadFileContent(result.FilePath)
			// TODO: Scroll to line
			m.transitionTo(NormalMode)
		}
		return m, nil

	case "up", "ctrl+p", "k":
		m.searchCursor = navigateCursor(m.searchCursor, -1, len(m.searchResults))
		return m, nil

	case "down", "ctrl+n", "j":
		m.searchCursor = navigateCursor(m.searchCursor, 1, len(m.searchResults))
		return m, nil

	case "n":
		// Next match
		m.searchCursor = navigateCursor(m.searchCursor, 1, len(m.searchResults))
		return m, nil

	case "N":
		// Previous match
		m.searchCursor = navigateCursor(m.searchCursor, -1, len(m.searchResults))
		return m, nil
	}
	return m, nil
}

func (m AppModel) handleLoadingMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Allow quit during loading
	if msg.String() == "q" || msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	// Ignore other keys while loading
	return m, nil
}

func (m AppModel) handleTOCMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.tableOfContents == nil {
		m.transitionTo(NormalMode)
		return m, nil
	}

	switch msg.String() {
	case "esc", "q", "t":
		// Close TOC modal
		m.transitionTo(NormalMode)
		return m, nil

	case "enter":
		// Jump to selected heading
		if m.tableOfContents.HasEntries() {
			lineNum := m.tableOfContents.GetSelectedLineNum()
			m.gotoLine(lineNum)
			m.transitionTo(NormalMode)
		}
		return m, nil

	case "up", "k":
		m.tableOfContents.SelectPrevious()
		return m, nil

	case "down", "j":
		m.tableOfContents.SelectNext()
		return m, nil

	case "g":
		m.tableOfContents.SelectFirst()
		return m, nil

	case "G":
		m.tableOfContents.SelectLast()
		return m, nil

	case "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m AppModel) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Global keybindings that work in normal mode
	switch msg.String() {
	case "?":
		m.transitionTo(HelpMode)
		return m, nil

	case "/":
		// Activate finder
		if len(m.markdownFiles) == 0 {
			// Async load files (Blocker 1 Fix)
			m.transitionTo(LoadingMode)
			m.loadingMessage = "Loading files..."
			return m, loadMarkdownFilesCmd(m.rootPath)
		}
		// Files already loaded
		m.finderItems = m.markdownFiles
		m.finderFiltered = m.markdownFiles
		m.finderCursor = 0
		m.transitionTo(FinderMode)
		return m, nil

	case "ctrl+f":
		// Activate search
		m.transitionTo(SearchMode)
		return m, nil

	case "t":
		// Activate TOC (Table of Contents)
		if m.selectedFile != "" && m.tableOfContents != nil && m.tableOfContents.HasEntries() {
			m.transitionTo(TOCMode)
		}
		return m, nil
	}

	// Use configurable keybindings for other actions
	viewName := []string{"filetree", "viewer", "preview"}[m.currentView]
	action := m.keyBindings.FindAction(msg.String(), viewName)

	// If no view-specific action, try "any" view
	if action == "" {
		action = m.keyBindings.FindAction(msg.String(), "any")
	}

	// Execute action
	switch action {
		// App controls
		case "quit":
			return m, tea.Quit

		// Navigation
		case "back":
			if m.currentView == FileTreeView {
				if m.fileList.FilterState() == list.Filtering {
					m.fileList.ResetFilter()
				} else {
					m.navigateUp()
				}
			}

		case "open":
			if m.currentView == FileTreeView {
				m.navigateToSelectedFile()
			}

		case "switch_view":
			m.currentView = (m.currentView + 1) % 3
			m.clipboard.ClearSelection() // Clear selection when switching views
			m.mouseDragActive = false    // Reset drag state when switching views

		// File tree navigation
		case "down":
			if m.currentView == FileTreeView {
				m.fileList, cmd = m.fileList.Update(msg)
			}
			return m, cmd

		case "up":
			if m.currentView == FileTreeView {
				m.fileList, cmd = m.fileList.Update(msg)
			}
			return m, cmd

		// Viewer scrolling
		case "scroll_down":
			if m.currentView == ViewerView {
				m.viewer.LineDown(1)
			}

		case "scroll_up":
			if m.currentView == ViewerView {
				m.viewer.LineUp(1)
			}

		case "page_down":
			if m.currentView == ViewerView {
				m.viewer.HalfViewDown()
			}

		case "page_up":
			if m.currentView == ViewerView {
				m.viewer.HalfViewUp()
			}

		case "view_down":
			if m.currentView == ViewerView {
				m.viewer.ViewDown()
			}

		case "view_up":
			if m.currentView == ViewerView {
				m.viewer.ViewUp()
			}

		case "top":
			if m.currentView == ViewerView {
				m.viewer.GotoTop()
			}

		case "bottom":
			if m.currentView == ViewerView {
				m.viewer.GotoBottom()
			}

		// NEW: Copy functionality
		case "copy":
			debugFile, _ := os.OpenFile("/tmp/lumina_mouse_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if debugFile != nil {
				sl, sc, el, ec := m.clipboard.GetSelectionBounds()
				fmt.Fprintf(debugFile, "COPY KEY PRESSED: Selection bounds: %d:%d to %d:%d, HasSelection=%v\n",
					sl, sc, el, ec, m.clipboard.HasSelection())
				debugFile.Close()
			}
			if m.currentView == ViewerView {
				if err := m.clipboard.CopySelection(m.viewerContent); err == nil {
					debugFile, _ := os.OpenFile("/tmp/lumina_mouse_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if debugFile != nil {
						fmt.Fprintf(debugFile, "  → COPY SUCCESSFUL\n")
						debugFile.Close()
					}
					// Copied! Could show status message
					// For now, selection clears
					m.clipboard.ClearSelection()
				} else {
					debugFile, _ := os.OpenFile("/tmp/lumina_mouse_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if debugFile != nil {
						fmt.Fprintf(debugFile, "  → COPY FAILED: %v\n", err)
						debugFile.Close()
					}
				}
			}

		// Filter
		case "filter":
			if m.currentView == FileTreeView {
				// Toggle filter - let bubbles handle it
				if m.fileList.FilterState() != list.Filtering {
					m.fileList.SetFilteringEnabled(true)
				}
			}

		case "help":
			m.transitionTo(HelpMode)
			return m, nil
		}

	// Update components based on current view
	if m.currentView == FileTreeView {
		m.fileList, cmd = m.fileList.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleMouseEvent processes mouse events for text selection and scrolling
func (m AppModel) handleMouseEvent(msg tea.MouseMsg) (AppModel, tea.Cmd) {
	// Handle mouse wheel scrolling in Viewer pane
	if msg.Type == 5 && m.currentView == ViewerView {
		// Scroll up
		m.viewer.LineUp(3)
		return m, nil
	}
	if msg.Type == 6 && m.currentView == ViewerView {
		// Scroll down
		m.viewer.LineDown(3)
		return m, nil
	}

	// Handle text selection in Viewer pane
	if (msg.Type == 1 || msg.Type == 11) && m.currentView == ViewerView && m.selectedFile != "" {
		// Transform screen coordinates to document coordinates
		textLine, textCol := m.screenToDocCoords(msg.X, msg.Y)

		// Bounds check - click outside viewer pane (negative coords from screen transform)
		paneRelativeX := msg.X - (m.fileTreeWidth + 2)
		paneRelativeY := msg.Y - 3
		if paneRelativeX < 0 || paneRelativeY < 0 {
			return m, nil
		}

		// Type 1 = button press events
		if msg.Type == 1 && msg.Button == 1 {
			if msg.Action == 0 {
				// Button press - start selection with timeout recovery
				cmd := m.startMouseSelection(textLine, textCol)
				return m, cmd
			}
			if msg.Action == 2 {
				// Drag - extend selection (or start if not already started)
				if !m.mouseDragActive {
					// First drag event - start selection
					cmd := m.startMouseSelection(textLine, textCol)
					return m, cmd
				}
				// Extend selection
				m.clipboard.StartSelection(m.mouseDragStartLine, m.mouseDragStartCol, false)
				m.clipboard.ExtendSelection(textLine, textCol)
				return m, nil
			}
			if msg.Action == 3 {
				// Release - end selection
				m.mouseDragActive = false
				return m, nil
			}
		}

		// Type 11 = motion events (pure mouse movement, no button info)
		if msg.Type == 11 && msg.Action == 2 {
			if !m.mouseDragActive {
				// First motion - start selection with timeout recovery
				cmd := m.startMouseSelection(textLine, textCol)
				return m, cmd
			} else {
				// Continue selection
				m.clipboard.StartSelection(m.mouseDragStartLine, m.mouseDragStartCol, false)
				m.clipboard.ExtendSelection(textLine, textCol)
			}
			return m, nil
		}
	}

	// Handle file tree clicks
	if m.currentView == FileTreeView {
		var cmd tea.Cmd
		m.fileList, cmd = m.fileList.Update(msg)
		return m, cmd
	}

	return m, nil
}

// startMouseSelection initializes mouse-based text selection with timeout recovery
func (m *AppModel) startMouseSelection(line, col int) tea.Cmd {
	m.mouseDragActive = true
	m.mouseDragStartLine = line
	m.mouseDragStartCol = col
	m.selectionMode = SelectionCharacter
	m.clipboard.StartSelection(line, col, false)
	return dragTimeoutCmd() // Blocker 3 Fix: Auto-recover after 5 seconds
}

// gotoLine scrolls the viewer to a specific line number
func (m *AppModel) gotoLine(lineNum int) {
	if lineNum < 0 {
		lineNum = 0
	}

	// Set the viewport's Y offset to the target line
	m.viewer.YOffset = lineNum

	// Ensure we don't scroll past the end
	maxOffset := len(strings.Split(m.viewer.View(), "\n")) - m.viewer.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.viewer.YOffset > maxOffset {
		m.viewer.YOffset = maxOffset
	}
}

// screenToDocCoords transforms screen coordinates to document coordinates
func (m *AppModel) screenToDocCoords(screenX, screenY int) (line, col int) {
	// Account for viewer pane's actual position
	// File tree (m.fileTreeWidth) + left border (1) + file tree border (1)
	viewerStartX := m.fileTreeWidth + 2
	viewerStartY := 3 // header height

	// Transform screen coords to pane-relative coords
	paneRelativeX := screenX - viewerStartX
	paneRelativeY := screenY - viewerStartY

	// Transform to document coordinates
	line = paneRelativeY + m.viewer.YOffset
	col = paneRelativeX

	return
}

// highlightSelection applies visual highlighting to selected text
func (m AppModel) highlightSelection(content string) string {
	sl, sc, el, ec := m.clipboard.GetSelectionBounds()

	// Normalize selection bounds (handle backwards selection)
	if sl > el || (sl == el && sc > ec) {
		sl, sc, el, ec = el, ec, sl, sc
	}

	// Split content into lines
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return content
	}

	// Highlight selected lines
	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("11")). // Bright yellow background
		Foreground(lipgloss.Color("0"))   // Black text

	for i := range lines {
		// Use ABSOLUTE line numbers from the full document
		if i >= sl && i <= el {
			line := lines[i]

			if i == sl && i == el {
				// Single line selection
				if sc >= 0 && ec <= len(line) && sc <= ec {
					before := line[:sc]
					selected := line[sc:ec]
					after := line[ec:]
					lines[i] = before + highlightStyle.Render(selected) + after
				}
			} else if i == sl {
				// First line of multi-line selection
				if sc >= 0 && sc <= len(line) {
					before := line[:sc]
					selected := line[sc:]
					lines[i] = before + highlightStyle.Render(selected)
				}
			} else if i == el {
				// Last line of multi-line selection
				if ec >= 0 && ec <= len(line) {
					selected := line[:ec]
					after := line[ec:]
					lines[i] = highlightStyle.Render(selected) + after
				}
			} else {
				// Middle lines - highlight entire line
				lines[i] = highlightStyle.Render(line)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// Modal Rendering Helpers (DRY Principle)

// createModalStyles returns common modal styles
func (m AppModel) createModalStyles() (modalBase, title, item, selectedItem lipgloss.Style) {
	modalBase = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("active-border"))).
		Padding(1, 2)

	title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colorManager.GetColor("title"))).
		Bold(true)

	item = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	selectedItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("11")).
		Bold(true)

	return
}

// renderModal renders a centered modal with given content
func (m AppModel) renderModal(content string, width, height int) string {
	modalBase, _, _, _ := m.createModalStyles()
	modalStyle := modalBase.Width(width).Height(height)
	return modalStyle.Render(content)
}

// renderLoadingModal renders a loading indicator
func (m AppModel) renderLoadingModal() string {
	_, titleStyle, _, _ := m.createModalStyles()
	content := titleStyle.Render(m.loadingMessage + "...")
	return m.renderModal(content, 30, 5)
}

// renderSearchModal renders the ripgrep search UI (Phase 3 Week 2)
func (m AppModel) renderSearchModal() string {
	modalWidth := min(m.width-10, 100)
	modalHeight := min(m.height-10, 30)

	_, titleStyle, itemStyle, selectedItemStyle := m.createModalStyles()

	var content strings.Builder

	// Search query input
	content.WriteString(titleStyle.Render("🔍 Search: " + m.searchQuery))
	content.WriteString("\n\n")

	// Results
	if m.searchInProgress {
		content.WriteString(itemStyle.Render("Searching..."))
	} else if len(m.searchResults) == 0 {
		content.WriteString(itemStyle.Render("No results found"))
	} else {
		// Show up to 20 results
		maxResults := min(20, len(m.searchResults))
		for i := 0; i < maxResults; i++ {
			result := m.searchResults[i]
			line := fmt.Sprintf("%s:%d: %s", filepath.Base(result.FilePath), result.Line, result.Text)

			if i == m.searchCursor {
				content.WriteString(selectedItemStyle.Render("▶ " + line))
			} else {
				content.WriteString(itemStyle.Render("  " + line))
			}
			if i < maxResults-1 {
				content.WriteString("\n")
			}
		}

		if len(m.searchResults) > maxResults {
			content.WriteString("\n")
			content.WriteString(itemStyle.Render(fmt.Sprintf("  ... and %d more", len(m.searchResults)-maxResults)))
		}
	}

	content.WriteString("\n\n")
	content.WriteString(itemStyle.Render("↑/↓: navigate | Enter: jump | n/N: next/prev | Esc: cancel"))

	return m.renderModal(content.String(), modalWidth, modalHeight)
}

// renderTOCModal renders the Table of Contents modal
func (m AppModel) renderTOCModal() string {
	if m.tableOfContents == nil || !m.tableOfContents.HasEntries() {
		return ""
	}

	// Modal dimensions
	modalWidth := min(60, m.width-10)
	modalHeight := min(m.height-10, m.tableOfContents.GetEntryCount()+6)

	// Get common styles
	_, titleStyle, itemStyle, selectedItemStyle := m.createModalStyles()

	// Build content
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("📋 Table of Contents"))
	content.WriteString("\n\n")

	// TOC entries
	selectedIdx := m.tableOfContents.GetSelectedIndex()
	maxDisplay := modalHeight - 6 // Leave room for title and footer

	// Calculate scroll window
	startIdx := 0
	if selectedIdx >= maxDisplay {
		startIdx = selectedIdx - maxDisplay + 1
	}
	endIdx := startIdx + maxDisplay
	if endIdx > m.tableOfContents.GetEntryCount() {
		endIdx = m.tableOfContents.GetEntryCount()
	}

	// Render visible entries
	for i := startIdx; i < endIdx; i++ {
		entry := m.tableOfContents.entries[i]

		// Indent based on heading level
		indent := strings.Repeat("  ", entry.Level-1)

		// Truncate title if needed
		availableWidth := modalWidth - len(indent) - 6 // 6 for borders and arrow
		title := entry.Title
		if len(title) > availableWidth {
			title = title[:availableWidth-1] + "…"
		}

		line := indent + title

		if i == selectedIdx {
			content.WriteString(selectedItemStyle.Render("▶ " + line))
		} else {
			content.WriteString(itemStyle.Render("  " + line))
		}
		if i < endIdx-1 {
			content.WriteString("\n")
		}
	}

	// Navigation hint
	content.WriteString("\n\n")
	content.WriteString(itemStyle.Render("j/k: navigate | Enter: jump | t/Esc: close"))

	return m.renderModal(content.String(), modalWidth, modalHeight)
}

// renderFinderModal renders the fuzzy finder modal overlay using flat state (Blocker 2 Fix)
func (m AppModel) renderFinderModal() string {
	// Modal dimensions
	modalWidth := min(80, m.width-10)
	modalHeight := min(20, m.height-10)

	// Get common styles
	_, inputStyle, itemStyle, selectedItemStyle := m.createModalStyles()

	// Build content
	var content strings.Builder

	// Input line
	content.WriteString(inputStyle.Render("🔍 Find: " + m.finderInput + "█"))
	content.WriteString("\n\n")

	// Results (using flat state - Blocker 2 Fix)
	results := m.finderFiltered
	if len(results) == 0 {
		content.WriteString(itemStyle.Render("No matches found"))
	} else {
		// Show up to 15 results
		maxResults := min(15, len(results))
		cursor := m.finderCursor

		for i := 0; i < maxResults; i++ {
			if i == cursor {
				content.WriteString(selectedItemStyle.Render("▶ " + results[i]))
			} else {
				content.WriteString(itemStyle.Render("  " + results[i]))
			}
			if i < maxResults-1 {
				content.WriteString("\n")
			}
		}

		// Show count if more results exist
		if len(results) > maxResults {
			content.WriteString("\n")
			content.WriteString(itemStyle.Render(fmt.Sprintf("  ... and %d more", len(results)-maxResults)))
		}
	}

	content.WriteString("\n\n")
	content.WriteString(itemStyle.Render("↑/↓: navigate | Enter: select | Esc: cancel"))

	return m.renderModal(content.String(), modalWidth, modalHeight)
}

// View renders the UI
func (m AppModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Define styles using color manager
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colorManager.GetColor("title"))).
		Padding(0, 2)

	// Pane colors with theme support
	paneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("inactive-border")))

	activePaneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("active-border"))).
		Bold(true)

	// Header
	title := titleStyle.Render("Claude Code Navigator (CCN) - " + filepath.Base(m.currentPath))

	// File tree pane
	fileTreeStyle := paneStyle
	if m.currentView == FileTreeView {
		fileTreeStyle = activePaneStyle
	}
	fileTreePane := fileTreeStyle.
		Width(m.fileTreeWidth).
		Render(m.fileList.View())

	// Viewer pane
	viewerStyle := paneStyle
	if m.currentView == ViewerView {
		viewerStyle = activePaneStyle
	}
	// Determine viewer content - use cached rendered version
	// Selection highlighting is handled by clipboard/highlighting system separately
	var viewerContent string
	if m.selectedFile == "" {
		viewerContent = "No file selected\n\nNavigate in the file tree and press Enter to view a file."
	} else {
		// Always use the cached rendered version for consistency
		// Selection highlighting should be done on the rendered content, not the raw markdown
		viewerContent = m.renderedContent
	}

	// Update viewer with final content
	m.viewer.SetContent(viewerContent)

	viewerPane := viewerStyle.
		Width(m.viewerWidth).
		Render(m.viewer.View())

	// Preview pane
	previewStyle := paneStyle
	if m.currentView == PreviewView {
		previewStyle = activePaneStyle
	}
	previewPane := previewStyle.
		Width(m.previewWidth).
		Render("Preview\n(Coming soon)")

	// Join panes horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		fileTreePane,
		viewerPane,
		previewPane,
	)

	// Status bar with theme-aware colors
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colorManager.GetColor("status-bar"))).
		Padding(0, 1)

	viewName := []string{"FILE TREE", "VIEWER", "PREVIEW"}[m.currentView]
	var statusText string

	switch m.currentView {
	case FileTreeView:
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: nav | Enter: open | h/Esc: back | /: fuzzy find | ?: help | q: quit", viewName)
	case ViewerView:
		tocHint := ""
		if m.tableOfContents != nil && m.tableOfContents.HasEntries() {
			tocHint = " | t: TOC"
		}
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | /: find%s | y: copy | ?: help | q: quit", viewName, tocHint)
	case PreviewView:
		statusText = fmt.Sprintf("[%s] Tab: switch | /: fuzzy find | ?: help | q: quit", viewName)
	}

	status := statusStyle.Render(statusText)

	baseView := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		status,
	)

	// State Machine: Render overlays based on current mode (Blocker 4 Fix)
	switch m.currentMode {
	case HelpMode:
		helpOverlay := getHelpOverlay(m.width, m.height, m.colorManager)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			helpOverlay,
			lipgloss.WithWhitespaceChars(" "),
		)

	case FinderMode:
		finderOverlay := m.renderFinderModal()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			finderOverlay,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
		)

	case SearchMode:
		// TODO: Implement search UI (Phase 3 Week 2)
		searchOverlay := m.renderSearchModal()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			searchOverlay,
			lipgloss.WithWhitespaceChars(" "),
		)

	case TOCMode:
		// Table of Contents overlay
		tocOverlay := m.renderTOCModal()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			tocOverlay,
			lipgloss.WithWhitespaceChars(" "),
		)

	case LoadingMode:
		loadingOverlay := m.renderLoadingModal()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			loadingOverlay,
			lipgloss.WithWhitespaceChars(" "),
		)

	case NormalMode:
		return baseView
	}

	return baseView
}

func main() {
	// Define command-line flags
	var (
		showHelp    bool
		showVersion bool
		showKeys    bool
	)

	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.BoolVar(&showKeys, "keys", false, "Show keyboard shortcuts reference")
	flag.BoolVar(&showKeys, "k", false, "Show keyboard shortcuts reference (shorthand)")

	// Custom usage message
	flag.Usage = printHelp

	// Parse flags
	flag.Parse()

	// Handle flags
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if showKeys {
		printKeyboardShortcuts()
		os.Exit(0)
	}

	// Get the root path from remaining arguments or use current directory
	rootPath := "."
	if flag.NArg() > 0 {
		rootPath = flag.Arg(0)
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Verify path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Path does not exist: %s\n", absPath)
		os.Exit(1)
	}

	// Create the initial model
	m := NewAppModel(absPath)

	// Create the Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
