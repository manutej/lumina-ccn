package main

import (
	"context"
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
	query  string
	ch     <-chan RipgrepMatch
	cancel context.CancelFunc
}
type SearchResultMsg struct {
	result RipgrepResult
}
type SearchCompletedMsg struct{}
type SearchErrorMsg struct {
	err error
}

// SearchInputPhase represents whether user is typing query or viewing results
type SearchInputPhase int

const (
	SearchPhaseInput   SearchInputPhase = iota // User is typing search query
	SearchPhaseResults                         // User is browsing results
)

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

// startSearchCmd initiates a ripgrep search (Milestone 2)
func startSearchCmd(rootPath, query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		manager := NewRipgrepManager(4)
		ch, err := manager.Search(ctx, query)
		if err != nil {
			cancel()
			return SearchErrorMsg{err: err}
		}
		return SearchStartedMsg{query: query, ch: ch, cancel: cancel}
	}
}

// listenSearchResultsCmd listens for search results from ripgrep (Milestone 2)
func listenSearchResultsCmd(ch <-chan RipgrepMatch) tea.Cmd {
	return func() tea.Msg {
		match, ok := <-ch
		if !ok {
			return SearchCompletedMsg{}
		}
		return SearchResultMsg{
			result: RipgrepResult{
				FilePath: match.Path,
				Line:     match.LineNumber,
				Text:     match.Text,
			},
		}
	}
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

	case SearchStartedMsg:
		// Search initiated - store cancel function, channel, and start listening
		m.searchCancel = msg.cancel
		m.searchResultsChan = msg.ch
		m.searchInProgress = true
		m.searchPhase = 1 // Results phase
		return m, listenSearchResultsCmd(msg.ch)

	case SearchResultMsg:
		// Accumulate search results and continue listening
		m.searchResults = append(m.searchResults, msg.result)
		// Continue listening for more results from stored channel
		if m.searchResultsChan != nil && m.searchInProgress {
			return m, listenSearchResultsCmd(m.searchResultsChan)
		}
		return m, nil

	case SearchCompletedMsg:
		m.searchInProgress = false
		m.searchResultsChan = nil
		return m, nil

	case SearchErrorMsg:
		m.searchInProgress = false
		m.searchResultsChan = nil
		m.searchErrorMessage = msg.err.Error()
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

	case "up", "ctrl+p", "k":
		// Navigate up in results (Blocker 2: Pure function)
		m.finderCursor = navigateCursor(m.finderCursor, -1, len(m.finderFiltered))
		return m, nil

	case "down", "ctrl+n", "j":
		// Navigate down in results (Blocker 2: Pure function)
		m.finderCursor = navigateCursor(m.finderCursor, 1, len(m.finderFiltered))
		return m, nil

	case "backspace":
		// Delete character from input
		if len(m.finderInput) > 0 {
			m.finderInput = m.finderInput[:len(m.finderInput)-1]
			m.finderFiltered = filterItems(m.finderItems, m.finderInput) // Blocker 2: Pure function
			m.finderCursor = 0                                           // Reset cursor to top
		}
		return m, nil

	default:
		// Add character to input if it's a single printable character
		if len(msg.String()) == 1 {
			m.finderInput += msg.String()
			m.finderFiltered = filterItems(m.finderItems, m.finderInput) // Blocker 2: Pure function
			m.finderCursor = 0                                           // Reset cursor to top
		}
		return m, nil
	}
}

func (m AppModel) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Phase 0: Input phase - user is typing query
	if m.searchPhase == 0 {
		switch msg.String() {
		case "esc":
			m.transitionTo(NormalMode)
			return m, nil

		case "enter":
			// Execute search if query is not empty
			if len(m.searchQuery) > 0 {
				m.searchResults = []RipgrepResult{} // Clear previous results
				m.searchCursor = 0
				m.searchInProgress = true
				m.searchErrorMessage = ""
				return m, startSearchCmd(m.rootPath, m.searchQuery)
			}
			return m, nil

		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
			return m, nil

		default:
			// Add printable characters to query
			if len(msg.String()) == 1 {
				m.searchQuery += msg.String()
			}
			return m, nil
		}
	}

	// Phase 1: Results phase - user is browsing results
	switch msg.String() {
	case "esc":
		m.transitionTo(NormalMode)
		return m, nil

	case "enter":
		// Jump to selected result
		if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
			result := m.searchResults[m.searchCursor]
			m.loadFileContent(result.FilePath)
			// Scroll to the matching line
			for i := 0; i < result.Line && i < 1000; i++ {
				m.viewer.LineDown(1)
			}
			m.transitionTo(NormalMode)
		}
		return m, nil

	case "backspace":
		// Go back to input phase
		m.searchPhase = 0
		m.searchResults = []RipgrepResult{}
		m.searchCursor = 0
		if m.searchCancel != nil {
			m.searchCancel()
			m.searchCancel = nil
		}
		m.searchInProgress = false
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
	}

	// Letter jump navigation: Shift+Letter to jump to next file starting with that letter
	if m.currentView == FileTreeView && len(msg.String()) == 1 {
		char := msg.String()
		// Check if it's an uppercase letter (Shift+letter)
		if char >= "A" && char <= "Z" {
			m.jumpToNextFileStartingWith(char)
			return m, nil
		}
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

	case "switch_view":
		m.currentView = (m.currentView + 1) % 3
		m.clipboard.ClearSelection() // Clear selection when switching views
		m.mouseDragActive = false    // Reset drag state when switching views

	// File tree navigation
	case "down":
		if m.currentView == FileTreeView {
			m.fileList, cmd = m.fileList.Update(msg)
			return m, cmd
		} else if m.currentView == PreviewView && m.contextPanel.currentMode == TOCMode {
			m.contextPanel.GetTOC().SelectNext()
		}

	case "up":
		if m.currentView == FileTreeView {
			m.fileList, cmd = m.fileList.Update(msg)
			return m, cmd
		} else if m.currentView == PreviewView && m.contextPanel.currentMode == TOCMode {
			m.contextPanel.GetTOC().SelectPrevious()
		}

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
		} else if m.currentView == PreviewView && m.contextPanel.currentMode == TOCMode {
			m.contextPanel.GetTOC().SelectFirst()
		}

	case "bottom":
		if m.currentView == ViewerView {
			m.viewer.GotoBottom()
		} else if m.currentView == PreviewView && m.contextPanel.currentMode == TOCMode {
			m.contextPanel.GetTOC().SelectLast()
		}

	// Context Panel (PreviewView) - Mode switching
	case "m":
		// Cycle through panel modes: TOC → File Info → Stats → Quick Actions
		if m.currentView == PreviewView {
			m.contextPanel.CycleMode()
		}

	case "open":
		// File tree: navigate into file/directory
		if m.currentView == FileTreeView {
			m.navigateToSelectedFile()
			// Context Panel: Jump to selected TOC entry in viewer
		} else if m.currentView == PreviewView && m.contextPanel.currentMode == TOCMode {
			if entry := m.contextPanel.GetTOC().GetSelectedEntry(); entry != nil {
				m.viewer.GotoTop()
				// Jump to the line number
				for i := 0; i < entry.LineNum; i++ {
					m.viewer.LineDown(1)
				}
				// Switch to viewer pane to show the jumped location
				m.currentView = ViewerView
			}
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

	case "sort":
		// Cycle through sort modes
		m.cycleSortMode()
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

// renderLoadingModal renders a loading indicator
func (m AppModel) renderLoadingModal() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("active-border"))).
		Padding(2, 4)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colorManager.GetColor("title"))).
		Bold(true)

	content := textStyle.Render(m.loadingMessage + "...")
	return modalStyle.Render(content)
}

// renderSearchModal renders the ripgrep search UI (Phase 3 Week 2)
func (m AppModel) renderSearchModal() string {
	modalWidth := min(m.width-10, 100)
	modalHeight := min(m.height-10, 30)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("active-border"))).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colorManager.GetColor("title"))).
		Bold(true)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("11")).
		Bold(true)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")). // Red
		Bold(true)

	var content strings.Builder

	// Search query input with cursor
	if m.searchPhase == 0 {
		content.WriteString(titleStyle.Render("🔍 Search: " + m.searchQuery + "█"))
	} else {
		content.WriteString(titleStyle.Render("🔍 Search: " + m.searchQuery))
	}
	content.WriteString("\n\n")

	// Show error message if any
	if m.searchErrorMessage != "" {
		content.WriteString(errorStyle.Render("Error: " + m.searchErrorMessage))
		content.WriteString("\n\n")
	}

	// Phase 0: Input phase
	if m.searchPhase == 0 {
		if m.searchQuery == "" {
			content.WriteString(itemStyle.Render("Type a search query and press Enter"))
		} else {
			content.WriteString(itemStyle.Render("Press Enter to search"))
		}
	} else {
		// Phase 1: Results phase
		if m.searchInProgress {
			content.WriteString(itemStyle.Render(fmt.Sprintf("Searching... (%d results)", len(m.searchResults))))
		} else if len(m.searchResults) == 0 {
			content.WriteString(itemStyle.Render("No results found for \"" + m.searchQuery + "\""))
		} else {
			// Show up to 20 results
			maxResults := min(20, len(m.searchResults))
			for i := 0; i < maxResults; i++ {
				result := m.searchResults[i]
				// Truncate text if too long
				text := strings.TrimSpace(result.Text)
				if len(text) > 60 {
					text = text[:57] + "..."
				}
				line := fmt.Sprintf("%s:%d: %s", filepath.Base(result.FilePath), result.Line, text)

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
	}

	content.WriteString("\n\n")
	if m.searchPhase == 0 {
		content.WriteString(itemStyle.Render("Enter: search | Esc: cancel"))
	} else {
		content.WriteString(itemStyle.Render("↑/↓: navigate | Enter: jump | Backspace: new search | Esc: cancel"))
	}

	return modalStyle.Render(content.String())
}

// renderFinderModal renders the fuzzy finder modal overlay using flat state (Blocker 2 Fix)
func (m AppModel) renderFinderModal() string {
	// Modal dimensions
	modalWidth := min(80, m.width-10)
	modalHeight := min(20, m.height-10)

	// Modal styles
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colorManager.GetColor("active-border"))).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colorManager.GetColor("title"))).
		Bold(true)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // Light gray

	selectedItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).  // Black
		Background(lipgloss.Color("11")). // Bright yellow
		Bold(true)

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

	return modalStyle.Render(content.String())
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
	// Determine viewer content with proper rendering order
	var viewerContent string
	if m.selectedFile == "" {
		viewerContent = "No file selected\n\nNavigate in the file tree and press Enter to view a file."
	} else if m.clipboard.HasSelection() {
		// CRITICAL: Apply highlighting to ORIGINAL content BEFORE rendering
		// This ensures selection highlighting targets correct coordinates
		highlighted := m.highlightSelection(m.viewerContent)
		// THEN render with glamour (adds ANSI codes after highlighting)
		rendered, err := m.markdownRenderer.Render(highlighted)
		if err != nil {
			viewerContent = highlighted // Fallback to highlighted plain text
		} else {
			viewerContent = rendered
		}
	} else {
		// No selection, use cached rendered version (performance optimization)
		viewerContent = m.renderedContent
	}

	// Update viewer with final content
	m.viewer.SetContent(viewerContent)

	viewerPane := viewerStyle.
		Width(m.viewerWidth).
		Render(m.viewer.View())

	// Preview pane - Context Panel (TOC, File Info, Stats, Quick Actions)
	previewStyle := paneStyle
	if m.currentView == PreviewView {
		previewStyle = activePaneStyle
	}

	previewContent := "📋 Context Panel\n\n(Open a markdown file to see content)"
	if m.selectedFile != "" {
		previewContent = m.contextPanel.Render(m.previewWidth, m.height, m.colorManager)
	}

	previewPane := previewStyle.
		Width(m.previewWidth).
		Render(previewContent)

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
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | /: fuzzy find | y: copy | ?: help | q: quit", viewName)
	case PreviewView:
		modeName := m.contextPanel.GetModeName()
		if m.contextPanel.GetTOC().HasEntries() && m.contextPanel.currentMode == TOCMode {
			statusText = fmt.Sprintf("[%s: %s] j/k: nav | Enter: jump | m: mode | ?: help | q: quit", viewName, modeName)
		} else {
			statusText = fmt.Sprintf("[%s: %s] m: mode | ?: help | q: quit", viewName, modeName)
		}
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
		// Search UI (Milestone 2 complete)
		searchOverlay := m.renderSearchModal()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			searchOverlay,
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
