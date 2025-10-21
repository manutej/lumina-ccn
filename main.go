package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

		// NEW: Use configurable keybindings
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
			if m.currentView == ViewerView {
				if err := m.clipboard.CopySelection(m.viewerContent); err == nil {
					// Copied! Could show status message
					// For now, selection clears
					m.clipboard.ClearSelection()
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
	viewerContent := m.viewer.View()
	if m.selectedFile == "" {
		viewerContent = "No file selected\n\nNavigate in the file tree and press Enter to view a file."
	}
	viewerPane := viewerStyle.
		Width(m.viewerWidth).
		Render(viewerContent)

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
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: nav | Enter: open | h/Esc: back | /: filter | ?: help | q: quit", viewName)
	case ViewerView:
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | y: copy | ?: help | q: quit", viewName)
	case PreviewView:
		statusText = fmt.Sprintf("[%s] Tab: switch | Coming soon | ?: help | q: quit", viewName)
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
