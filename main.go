package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Phase 3 Message Types for Fuzzy Finder

// FinderActivatedMsg indicates the fuzzy finder modal should be shown
type FinderActivatedMsg struct{}

// FinderInputMsg updates the finder input and triggers filtering
type FinderInputMsg struct {
	input string
}

// FinderSelectionMsg indicates a file was selected from the finder
type FinderSelectionMsg struct {
	filePath string
}

// FinderCanceledMsg closes the finder modal without selection
type FinderCanceledMsg struct{}

// Init initializes the model
func (m AppModel) Init() tea.Cmd {
	return nil
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

	case tea.KeyMsg:
		// Handle help overlay toggle
		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			return m, nil
		}

		// If help is shown, ESC or ? closes it, other keys are ignored
		if m.showHelp {
			if msg.String() == "esc" || msg.String() == "?" {
				m.showHelp = false
				return m, nil
			}
			// Ignore other keys when help is shown (except quit)
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}

		// Phase 3: Handle fuzzy finder modal input
		if m.finderActive {
			switch msg.String() {
			case "esc":
				// Cancel finder
				m.finderActive = false
				m.finderInput = ""
				m.currentMode = NormalMode
				return m, nil

			case "enter":
				// Select item and load file
				selected := m.fuzzyFinder.SelectedItem()
				if selected != "" {
					m.finderActive = false
					m.finderInput = ""
					m.currentMode = NormalMode
					// Load the selected file
					if err := m.loadFileContent(selected); err == nil {
						// Successfully loaded file
					}
				}
				return m, nil

			case "up", "ctrl+p":
				// Navigate up in results
				m.fuzzyFinder.HandleKey("up")
				return m, nil

			case "down", "ctrl+n":
				// Navigate down in results
				m.fuzzyFinder.HandleKey("down")
				return m, nil

			case "backspace":
				// Delete character from input
				if len(m.finderInput) > 0 {
					m.finderInput = m.finderInput[:len(m.finderInput)-1]
					m.fuzzyFinder.SetFilter(m.finderInput)
				}
				return m, nil

			default:
				// Add character to input if it's a single printable character
				if len(msg.String()) == 1 {
					m.finderInput += msg.String()
					m.fuzzyFinder.SetFilter(m.finderInput)
				}
				return m, nil
			}
		}

		// Phase 3: Activate finder with "/" key
		if msg.String() == "/" && !m.finderActive {
			debugFile, _ := os.OpenFile("/tmp/lumina_key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if debugFile != nil {
				fmt.Fprintf(debugFile, "FUZZY FINDER ACTIVATED: finderActive=%v, markdownFiles=%d\n", m.finderActive, len(m.markdownFiles))
				debugFile.Close()
			}
			// Collect all markdown files
			if len(m.markdownFiles) == 0 {
				m.markdownFiles = findMarkdownFiles(m.rootPath)
			}
			// Initialize fuzzy finder with file list
			m.fuzzyFinder = NewFuzzyFinderImpl(m.markdownFiles)
			m.finderActive = true
			m.finderInput = ""
			m.currentMode = FinderMode
			return m, nil
		}

		// NEW: Use configurable keybindings
		viewName := []string{"filetree", "viewer", "preview"}[m.currentView]
		action := m.keyBindings.FindAction(msg.String(), viewName)

		// If no view-specific action, try "any" view
		if action == "" {
			action = m.keyBindings.FindAction(msg.String(), "any")
		}

		// DEBUG: Log key presses and actions
		debugFile, _ := os.OpenFile("/tmp/lumina_key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "KEY: %q, VIEW: %s, ACTION: %q\n", msg.String(), viewName, action)
			debugFile.Close()
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
			} else if m.currentView == ViewerView && m.selectionMode != SelectionInactive {
				// Cancel selection in viewer mode
				m.selectionMode = SelectionInactive
				m.clipboard.ClearSelection()
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
				// Extend selection if in selection mode
				if m.selectionMode != SelectionInactive {
					currentLine := m.viewer.YOffset + m.viewer.Height - 1
					m.clipboard.ExtendSelection(currentLine, 0)
				}
			}

		case "scroll_up":
			if m.currentView == ViewerView {
				m.viewer.LineUp(1)
				// Extend selection if in selection mode
				if m.selectionMode != SelectionInactive {
					currentLine := m.viewer.YOffset
					m.clipboard.ExtendSelection(currentLine, 0)
				}
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

		// Visual selection mode (vim-style)
		case "start_selection":
			if m.currentView == ViewerView {
				m.selectionMode = SelectionCharacter
				// Start selection at current viewer position (top of visible area)
				currentLine := m.viewer.YOffset
				m.clipboard.StartSelection(currentLine, 0, false)
				// Extend by one character to match vim v behavior
				m.clipboard.ExtendSelection(currentLine, 1)
			}

		case "start_line_selection":
			if m.currentView == ViewerView {
				m.selectionMode = SelectionLine
				// Start line selection at current viewer position
				currentLine := m.viewer.YOffset
				m.clipboard.StartSelection(currentLine, 0, false)
				// Immediately extend to end of line (vim V behavior)
				m.clipboard.ExtendSelection(currentLine, 9999)
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
					m.selectionMode = SelectionInactive
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
			m.showHelp = !m.showHelp
		}
	}

	// Update components based on current view
	if m.currentView == FileTreeView {
		m.fileList, cmd = m.fileList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
				// Button press - start selection
				m.startMouseSelection(textLine, textCol)
				return m, nil
			}
			if msg.Action == 2 {
				// Drag - extend selection (or start if not already started)
				if !m.mouseDragActive {
					// First drag event - start selection
					m.startMouseSelection(textLine, textCol)
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
				// First motion - start selection
				m.startMouseSelection(textLine, textCol)
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

// startMouseSelection initializes mouse-based text selection
func (m *AppModel) startMouseSelection(line, col int) {
	m.mouseDragActive = true
	m.mouseDragStartLine = line
	m.mouseDragStartCol = col
	m.selectionMode = SelectionCharacter
	m.clipboard.StartSelection(line, col, false)
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

// renderFinderModal renders the fuzzy finder modal overlay
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

	// Results
	results := m.fuzzyFinder.FilteredResults()
	if len(results) == 0 {
		content.WriteString(itemStyle.Render("No matches found"))
	} else {
		// Show up to 15 results
		maxResults := min(15, len(results))
		cursor := m.fuzzyFinder.Cursor()

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
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | v/V: select | y: copy | /: fuzzy | ?: help | q: quit", viewName)
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

	// If help overlay is active, render it on top
	if m.showHelp {
		helpOverlay := getHelpOverlay(m.width, m.height, m.colorManager)

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			helpOverlay,
			lipgloss.WithWhitespaceChars(" "),
		)
	}

	// Phase 3: If finder modal is active, render it on top
	if m.finderActive {
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
