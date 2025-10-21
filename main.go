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

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			// ESC as universal cancel/back operation
			if m.currentView == FileTreeView && m.fileList.FilterState() == list.Filtering {
				// Cancel filtering if active
				m.fileList.ResetFilter()
			} else if m.currentView == FileTreeView {
				// Navigate up directory when in file tree
				m.navigateUp()
			}
			return m, nil

		// Navigation keys
		case "tab":
			// Cycle through views
			m.currentView = (m.currentView + 1) % 3
			return m, nil

		case "enter":
			// Navigate into selected file/directory
			if m.currentView == FileTreeView {
				m.navigateToSelectedFile()
			}
			return m, nil

		case "backspace", "h":
			// Navigate up directory
			if m.currentView == FileTreeView {
				m.navigateUp()
			}
			return m, nil

		// Vim-style navigation
		case "j", "down":
			if m.currentView == FileTreeView {
				m.fileList, cmd = m.fileList.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.currentView == ViewerView {
				m.viewer.LineDown(1)
			}
			return m, tea.Batch(cmds...)

		case "k", "up":
			if m.currentView == FileTreeView {
				m.fileList, cmd = m.fileList.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.currentView == ViewerView {
				m.viewer.LineUp(1)
			}
			return m, tea.Batch(cmds...)

		case "g":
			if m.currentView == ViewerView {
				m.viewer.GotoTop()
			}
			return m, nil

		case "G":
			if m.currentView == ViewerView {
				m.viewer.GotoBottom()
			}
			return m, nil

		case "d":
			if m.currentView == ViewerView {
				m.viewer.HalfViewDown()
			}
			return m, nil

		case "u":
			if m.currentView == ViewerView {
				m.viewer.HalfViewUp()
			}
			return m, nil

		// Option+Arrow keybindings for pagination
		case "alt+right", "alt+l":
			if m.currentView == ViewerView {
				m.viewer.ViewDown()
			}
			return m, nil

		case "alt+left", "alt+h":
			if m.currentView == ViewerView {
				m.viewer.ViewUp()
			}
			return m, nil

		case "alt+down":
			if m.currentView == ViewerView {
				m.viewer.HalfViewDown()
			}
			return m, nil

		case "alt+up":
			if m.currentView == ViewerView {
				m.viewer.HalfViewUp()
			}
			return m, nil
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

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 2)

	paneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD"))

	activePaneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF79C6"))

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

	// Status bar (context-aware based on active pane)
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(0, 1)

	viewName := []string{"FILE TREE", "VIEWER", "PREVIEW"}[m.currentView]
	var statusText string

	switch m.currentView {
	case FileTreeView:
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: nav | Enter: open | h/Esc: back | /: filter | ?: help | q: quit", viewName)
	case ViewerView:
		statusText = fmt.Sprintf("[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | ?: help | q: quit", viewName)
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
		helpOverlay := getHelpOverlay(m.width, m.height)

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
		showHelp     bool
		showVersion  bool
		showKeys     bool
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
