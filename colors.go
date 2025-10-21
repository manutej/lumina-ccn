package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ThemeMode represents the color theme
type ThemeMode string

const (
	ThemeDark  ThemeMode = "dark"
	ThemeLight ThemeMode = "light"
	ThemeAuto  ThemeMode = "auto"
)

// ColorPalette defines all colors used in the application
type ColorPalette struct {
	// UI Elements
	Title          string
	InactiveBorder string
	ActiveBorder   string
	StatusBar      string
	Background     string
	Foreground     string

	// Help Overlay
	HelpTitle       string
	HelpHeader      string
	HelpCommand     string
	HelpDescription string
	HelpBorder      string

	// Additional elements
	PaneIndicator   string
	SearchHighlight string
	ErrorText       string
	SuccessText     string
}

// ThemeConfig represents the saved theme configuration
type ThemeConfig struct {
	Theme  ThemeMode         `json:"theme"`
	Custom map[string]string `json:"custom,omitempty"`
}

var (
	// Dark theme palette
	darkPalette = ColorPalette{
		Title:          "#7D56F4", // Purple
		InactiveBorder: "#666666", // Dark gray
		ActiveBorder:   "#00D084", // Bright teal
		StatusBar:      "#666666", // Dark gray
		Background:     "#1e1e1e", // Near black
		Foreground:     "#F8F8F2", // Off-white

		HelpTitle:       "#7D56F4", // Purple
		HelpHeader:      "#FF79C6", // Pink
		HelpCommand:     "#8BE9FD", // Cyan
		HelpDescription: "#F8F8F2", // Off-white
		HelpBorder:      "#FF79C6", // Pink

		PaneIndicator:   "#00D084", // Bright teal
		SearchHighlight: "#FFD700", // Gold
		ErrorText:       "#FF5555", // Red
		SuccessText:     "#00D084", // Green
	}

	// Light theme palette - optimized for light backgrounds
	lightPalette = ColorPalette{
		Title:          "#5A4A8A", // Dark purple
		InactiveBorder: "#CCCCCC", // Light gray
		ActiveBorder:   "#007A5E", // Dark teal
		StatusBar:      "#333333", // Dark text
		Background:     "#FFFFFF", // White
		Foreground:     "#1A1A1A", // Near black

		HelpTitle:       "#5A4A8A", // Dark purple
		HelpHeader:      "#C1426B", // Dark pink
		HelpCommand:     "#0066CC", // Dark blue
		HelpDescription: "#2A2A2A", // Dark gray
		HelpBorder:      "#C1426B", // Dark pink

		PaneIndicator:   "#007A5E", // Dark teal
		SearchHighlight: "#CC8800", // Dark gold
		ErrorText:       "#CC0000", // Dark red
		SuccessText:     "#007A5E", // Dark green
	}
)

// ColorManager manages the application's color theme
type ColorManager struct {
	CurrentTheme ThemeMode
	Palette      ColorPalette
	ConfigPath   string
}

// NewColorManager creates a new color manager with automatic theme detection
func NewColorManager() (*ColorManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	cm := &ColorManager{
		ConfigPath: filepath.Join(home, ".config", "lumina", "colors.json"),
	}

	// Try to load saved configuration
	config, err := cm.loadConfig()
	if err == nil && config.Theme != "" {
		cm.CurrentTheme = config.Theme
	} else {
		// Default to auto-detection
		cm.CurrentTheme = ThemeAuto
	}

	// Apply theme
	cm.applyTheme()

	return cm, nil
}

// loadConfig loads the color configuration from disk
func (cm *ColorManager) loadConfig() (*ThemeConfig, error) {
	data, err := os.ReadFile(cm.ConfigPath)
	if err != nil {
		return nil, err
	}

	var config ThemeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// saveConfig saves the color configuration to disk
func (cm *ColorManager) saveConfig() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(cm.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	config := ThemeConfig{
		Theme: cm.CurrentTheme,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.ConfigPath, data, 0644)
}

// applyTheme applies the current theme
func (cm *ColorManager) applyTheme() {
	switch cm.CurrentTheme {
	case ThemeDark:
		cm.Palette = darkPalette
	case ThemeLight:
		cm.Palette = lightPalette
	case ThemeAuto:
		// Auto-detect based on terminal background
		cm.Palette = cm.detectTheme()
		cm.CurrentTheme = cm.getDetectedTheme()
	default:
		cm.Palette = darkPalette
		cm.CurrentTheme = ThemeDark
	}
}

// detectTheme detects the terminal background color and returns appropriate palette
func (cm *ColorManager) detectTheme() ColorPalette {
	// Try to detect terminal background using COLORFGBG environment variable
	// COLORFGBG format is "foreground;background"
	// Light backgrounds are typically: 7, 15, 0 (light)
	// Dark backgrounds are typically: 0, 8, etc.
	if bg := cm.detectBackgroundColor(); bg == "light" {
		return lightPalette
	}

	// Default to dark if detection fails
	return darkPalette
}

// detectBackgroundColor attempts to detect the terminal background color
func (cm *ColorManager) detectBackgroundColor() string {
	// Check COLORFGBG environment variable (set by many terminals)
	if colorBG := os.Getenv("COLORFGBG"); colorBG != "" {
		// Parse the background color (second number after semicolon)
		// Format is typically "foreground;background"
		parts := strings.Split(colorBG, ";")
		if len(parts) >= 2 {
			bg := parts[1]
			// Light backgrounds: 7 (white), 15 (bright white)
			// Dark backgrounds: 0 (black), 8 (bright black)
			if bg == "7" || bg == "15" {
				return "light"
			}
			return "dark"
		}
	}

	// Check for ITERM_SESSION_ID (iTerm on macOS often uses light backgrounds)
	if iterm := os.Getenv("ITERM_SESSION_ID"); iterm != "" {
		// iTerm users often prefer light backgrounds
		// This is a heuristic - could be improved with actual color detection
		if os.Getenv("COLORFGBG") == "" {
			// Without explicit COLORFGBG, assume dark for iTerm
			return "dark"
		}
	}

	// Check TERM_PROGRAM for terminal hints
	termProgram := os.Getenv("TERM_PROGRAM")
	if termProgram == "iTerm.app" {
		return "dark" // iTerm defaults to dark, but user can override
	}

	// Default to dark (most common for development terminals)
	return "dark"
}

// getDetectedTheme returns the theme that was detected
func (cm *ColorManager) getDetectedTheme() ThemeMode {
	if bg := cm.detectBackgroundColor(); bg == "light" {
		return ThemeLight
	}
	return ThemeDark
}

// SetTheme changes the current theme
func (cm *ColorManager) SetTheme(theme ThemeMode) error {
	cm.CurrentTheme = theme
	cm.applyTheme()
	return cm.saveConfig()
}

// GetStyle returns a styled text element
func (cm *ColorManager) GetStyle(name string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cm.GetColor(name)))
}

// GetColor returns a color by name from the current palette
func (cm *ColorManager) GetColor(name string) string {
	switch name {
	case "title":
		return cm.Palette.Title
	case "inactive-border":
		return cm.Palette.InactiveBorder
	case "active-border":
		return cm.Palette.ActiveBorder
	case "status-bar":
		return cm.Palette.StatusBar
	case "foreground":
		return cm.Palette.Foreground
	case "background":
		return cm.Palette.Background
	case "help-title":
		return cm.Palette.HelpTitle
	case "help-header":
		return cm.Palette.HelpHeader
	case "help-command":
		return cm.Palette.HelpCommand
	case "help-description":
		return cm.Palette.HelpDescription
	case "help-border":
		return cm.Palette.HelpBorder
	case "pane-indicator":
		return cm.Palette.PaneIndicator
	case "search-highlight":
		return cm.Palette.SearchHighlight
	case "error":
		return cm.Palette.ErrorText
	case "success":
		return cm.Palette.SuccessText
	default:
		return cm.Palette.Foreground
	}
}

// PrintThemeInfo prints information about the current theme
func (cm *ColorManager) PrintThemeInfo() {
	fmt.Printf("Current theme: %s\n", cm.CurrentTheme)
	fmt.Printf("Config file: %s\n", cm.ConfigPath)
	fmt.Printf("\nColor palette:\n")
	fmt.Printf("  Title: %s\n", cm.Palette.Title)
	fmt.Printf("  Active Border: %s\n", cm.Palette.ActiveBorder)
	fmt.Printf("  Inactive Border: %s\n", cm.Palette.InactiveBorder)
	fmt.Printf("  Status Bar: %s\n", cm.Palette.StatusBar)
	fmt.Printf("  Foreground: %s\n", cm.Palette.Foreground)
	fmt.Printf("  Background: %s\n", cm.Palette.Background)
}
