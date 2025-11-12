package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// printHelp displays the CLI help information
func printHelp() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF79C6"))

	cmdStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	help := strings.Builder{}

	// Title
	help.WriteString(titleStyle.Render("Lumina - Claude Code Navigator (CCN)"))
	help.WriteString("\n\n")

	// Description
	help.WriteString(descStyle.Render("A Glow-inspired TUI for navigating Claude Code workflows with beautiful markdown rendering."))
	help.WriteString("\n\n")

	// Usage
	help.WriteString(headerStyle.Render("USAGE:"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("lumina"))
	help.WriteString(" [directory] [flags]\n\n")

	// Options
	help.WriteString(headerStyle.Render("OPTIONS:"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("-h, --help"))
	help.WriteString("       Show this help message\n  ")
	help.WriteString(cmdStyle.Render("-v, --version"))
	help.WriteString("    Show version information\n  ")
	help.WriteString(cmdStyle.Render("-k, --keys"))
	help.WriteString("       Show keyboard shortcuts reference\n\n")

	// Arguments
	help.WriteString(headerStyle.Render("ARGUMENTS:"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("directory"))
	help.WriteString("       Path to navigate (default: current directory)\n\n")

	// Examples
	help.WriteString(headerStyle.Render("EXAMPLES:"))
	help.WriteString("\n  ")
	help.WriteString(descStyle.Render("# Navigate current directory"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("lumina"))
	help.WriteString("\n\n  ")
	help.WriteString(descStyle.Render("# Navigate specific directory"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("lumina ~/Documents/LUXOR/PROJECTS/LUMINA"))
	help.WriteString("\n\n  ")
	help.WriteString(descStyle.Render("# Show keyboard shortcuts"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("lumina --keys"))
	help.WriteString("\n\n  ")
	help.WriteString(descStyle.Render("# Show version"))
	help.WriteString("\n  ")
	help.WriteString(cmdStyle.Render("lumina --version"))
	help.WriteString("\n\n")

	// Links
	help.WriteString(headerStyle.Render("LEARN MORE:"))
	help.WriteString("\n  Documentation: ")
	help.WriteString(cmdStyle.Render("~/Documents/LUXOR/PROJECTS/LUMINA/ccn/README.md"))
	help.WriteString("\n  Quick Start:   ")
	help.WriteString(cmdStyle.Render("~/Documents/LUXOR/PROJECTS/LUMINA/ccn/QUICKSTART.md"))
	help.WriteString("\n  Linear:        ")
	help.WriteString(cmdStyle.Render("https://linear.app/ceti-luxor/project/lumina"))
	help.WriteString("\n\n")

	fmt.Print(help.String())
}

// printVersion displays version information
func printVersion() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF79C6"))

	version := strings.Builder{}

	version.WriteString(titleStyle.Render("Lumina (CCN)"))
	version.WriteString("\n\n")
	version.WriteString(labelStyle.Render("Version:     "))
	version.WriteString(infoStyle.Render(Version))
	version.WriteString("\n")
	version.WriteString(labelStyle.Render("Build Phase: "))
	version.WriteString(infoStyle.Render(BuildPhase))
	version.WriteString("\n")
	version.WriteString(labelStyle.Render("Build Date:  "))
	version.WriteString(infoStyle.Render(BuildDate))
	version.WriteString("\n\n")
	version.WriteString(infoStyle.Render("Built with the Charm Stack 🧙‍♂️✨"))
	version.WriteString("\n")

	fmt.Print(version.String())
}

// printKeyboardShortcuts displays keyboard shortcuts reference
func printKeyboardShortcuts() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF79C6"))

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8BE9FD"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	shortcuts := strings.Builder{}

	// Title
	shortcuts.WriteString(titleStyle.Render("Lumina Keyboard Shortcuts"))
	shortcuts.WriteString("\n\n")

	// General Navigation
	shortcuts.WriteString(headerStyle.Render("GENERAL"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Tab"))
	shortcuts.WriteString("           Cycle through panes (File Tree → Viewer → Preview)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("q, Ctrl+C"))
	shortcuts.WriteString("     Quit application\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("?"))
	shortcuts.WriteString("             Show this help (inside app)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Esc"))
	shortcuts.WriteString("           Close help overlay / Cancel / Go back\n\n")

	// File Tree Navigation
	shortcuts.WriteString(headerStyle.Render("FILE TREE (Left Pane)"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("j, ↓"))
	shortcuts.WriteString("          Navigate down\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("k, ↑"))
	shortcuts.WriteString("          Navigate up\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Shift+Letter"))
	shortcuts.WriteString("  Jump to next file starting with that letter (e.g. Shift+R)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Enter"))
	shortcuts.WriteString("         Open file or enter directory (use '..' to go up)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("h, Esc"))
	shortcuts.WriteString("        Go to parent directory\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("/"))
	shortcuts.WriteString("             Filter files (start typing)\n\n")

	// Viewer Navigation
	shortcuts.WriteString(headerStyle.Render("VIEWER (Center Pane)"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("j, ↓"))
	shortcuts.WriteString("          Scroll down one line\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("k, ↑"))
	shortcuts.WriteString("          Scroll up one line\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("d"))
	shortcuts.WriteString("             Scroll down half page\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("u"))
	shortcuts.WriteString("             Scroll up half page\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("g"))
	shortcuts.WriteString("             Go to top of file\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("G"))
	shortcuts.WriteString("             Go to bottom of file\n\n")

	// Context Panel (Right Pane)
	shortcuts.WriteString(headerStyle.Render("CONTEXT PANEL (Right Pane)"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("m"))
	shortcuts.WriteString("             Cycle modes: TOC → File Info → Stats → Actions\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("j, ↓"))
	shortcuts.WriteString("          Navigate TOC entries (in TOC mode)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("k, ↑"))
	shortcuts.WriteString("          Navigate TOC entries (in TOC mode)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Enter"))
	shortcuts.WriteString("         Jump to selected heading (in TOC mode)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("g/G"))
	shortcuts.WriteString("           First/last TOC entry\n\n")

	// Coming Soon
	shortcuts.WriteString(headerStyle.Render("COMING IN PHASE 2"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("/"))
	shortcuts.WriteString("             Fuzzy file finder (telescope-style)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("Ctrl+F"))
	shortcuts.WriteString("        Content search (ripgrep)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("n, N"))
	shortcuts.WriteString("          Next/previous search match\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("m{a-z}"))
	shortcuts.WriteString("        Set bookmark (vim marks)\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(keyStyle.Render("'{a-z}"))
	shortcuts.WriteString("        Jump to bookmark\n\n")

	// Tips
	shortcuts.WriteString(headerStyle.Render("TIPS"))
	shortcuts.WriteString("\n")
	shortcuts.WriteString("  ")
	shortcuts.WriteString(descStyle.Render("• Use Tab to quickly switch focus between panes"))
	shortcuts.WriteString("\n  ")
	shortcuts.WriteString(descStyle.Render("• Active pane has a pink border"))
	shortcuts.WriteString("\n  ")
	shortcuts.WriteString(descStyle.Render("• Markdown files render with beautiful Glamour styling"))
	shortcuts.WriteString("\n  ")
	shortcuts.WriteString(descStyle.Render("• Vim users: All standard hjkl navigation works!"))
	shortcuts.WriteString("\n\n")

	fmt.Print(shortcuts.String())
}

// getHelpOverlay returns a styled help overlay for in-app display with theme support
func getHelpOverlay(width, height int, colorManager *ColorManager) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorManager.GetColor("help-border"))).
		Padding(1, 2).
		Width(width - 4).
		Height(height - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorManager.GetColor("help-title"))).
		Align(lipgloss.Center)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorManager.GetColor("help-command")))

	content := strings.Builder{}

	content.WriteString(titleStyle.Render("KEYBOARD SHORTCUTS"))
	content.WriteString("\n\n")

	// Quick reference table
	content.WriteString("GENERAL\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("Tab"))
	content.WriteString("       Switch panes\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("q"))
	content.WriteString("         Quit\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("?"))
	content.WriteString("         Toggle this help\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("Esc"))
	content.WriteString("       Close help / Cancel / Go back\n\n")

	content.WriteString("FILE TREE\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("j/k"))
	content.WriteString("       Navigate up/down\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("Shift+X"))
	content.WriteString("   Jump to file starting with X\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("Enter"))
	content.WriteString("     Open file/directory\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("h/Esc"))
	content.WriteString("     Parent directory\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("/"))
	content.WriteString("         Filter files\n\n")

	content.WriteString("VIEWER\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("j/k"))
	content.WriteString("       Scroll up/down\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("d/u"))
	content.WriteString("       Half page down/up\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("g/G"))
	content.WriteString("       Top/bottom\n\n")

	content.WriteString("CONTEXT PANEL\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("m"))
	content.WriteString("         Cycle modes (TOC/Info/Stats)\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("j/k"))
	content.WriteString("       Navigate (TOC mode)\n")
	content.WriteString("  ")
	content.WriteString(keyStyle.Render("Enter"))
	content.WriteString("     Jump to heading\n\n")

	content.WriteString("Press ")
	content.WriteString(keyStyle.Render("? or Esc"))
	content.WriteString(" to close")

	return boxStyle.Render(content.String())
}
