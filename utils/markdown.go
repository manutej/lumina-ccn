package utils

import (
	"os"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer handles markdown rendering
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	width    int
	style    string
}

// MarkdownStyle represents available rendering styles
type MarkdownStyle string

const (
	StyleAuto       MarkdownStyle = "auto"        // Auto-detect based on terminal
	StyleDark       MarkdownStyle = "dark"        // Dark theme (default)
	StyleLight      MarkdownStyle = "light"       // Light theme
	StyleDracula    MarkdownStyle = "dracula"     // Dracula theme
	StyleTokyoNight MarkdownStyle = "tokyo-night" // Tokyo Night theme
	StylePink       MarkdownStyle = "pink"        // Pink theme
	StyleNotty      MarkdownStyle = "notty"       // No TTY theme
)

// NewMarkdownRenderer creates a new markdown renderer with auto-detected theme
func NewMarkdownRenderer(width int) (*MarkdownRenderer, error) {
	// Use auto-detection for intelligent theme selection based on terminal
	return NewMarkdownRendererWithStyle(width, StyleAuto)
}

// NewMarkdownRendererWithStyle creates a renderer with a specific style
func NewMarkdownRendererWithStyle(width int, style MarkdownStyle) (*MarkdownRenderer, error) {
	var opts []glamour.TermRendererOption

	// First, check if GLAMOUR_STYLE environment variable is set
	// This takes precedence over programmatic settings
	if envStyle := os.Getenv("GLAMOUR_STYLE"); envStyle != "" {
		opts = append(opts, glamour.WithStylePath(envStyle))
		opts = append(opts, glamour.WithWordWrap(width))

		r, err := glamour.NewTermRenderer(opts...)
		if err != nil {
			return nil, err
		}

		return &MarkdownRenderer{
			renderer: r,
			width:    width,
			style:    "env:" + envStyle,
		}, nil
	}

	// Apply style option
	switch style {
	case StyleAuto:
		opts = append(opts, glamour.WithAutoStyle())
	case StyleDark:
		opts = append(opts, glamour.WithStylePath("dark"))
	case StyleLight:
		opts = append(opts, glamour.WithStylePath("light"))
	case StyleDracula:
		opts = append(opts, glamour.WithStylePath("dracula"))
	case StyleTokyoNight:
		opts = append(opts, glamour.WithStylePath("tokyo-night"))
	case StylePink:
		opts = append(opts, glamour.WithStylePath("pink"))
	case StyleNotty:
		opts = append(opts, glamour.WithStylePath("notty"))
	default:
		// Default to auto-detection
		opts = append(opts, glamour.WithAutoStyle())
	}

	// Word wrapping
	opts = append(opts, glamour.WithWordWrap(width))

	// Create renderer
	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		renderer: r,
		width:    width,
		style:    string(style),
	}, nil
}

// Render renders markdown content
func (m *MarkdownRenderer) Render(content string) (string, error) {
	return m.renderer.Render(content)
}

// UpdateWidth updates the renderer width
func (m *MarkdownRenderer) UpdateWidth(width int) error {
	// Create a new renderer with updated width, preserving style
	style := MarkdownStyle(m.style)
	renderer, err := NewMarkdownRendererWithStyle(width, style)
	if err != nil {
		return err
	}

	m.renderer = renderer.renderer
	m.width = width
	return nil
}

// SetStyle changes the rendering style
func (m *MarkdownRenderer) SetStyle(style MarkdownStyle) error {
	renderer, err := NewMarkdownRendererWithStyle(m.width, style)
	if err != nil {
		return err
	}

	m.renderer = renderer.renderer
	m.style = string(style)
	return nil
}

// GetStyle returns the current style
func (m *MarkdownRenderer) GetStyle() string {
	return m.style
}
