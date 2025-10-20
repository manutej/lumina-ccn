package utils

import (
	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer handles markdown rendering
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer(width int) (*MarkdownRenderer, error) {
	// Create a glamour renderer with dark style
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		renderer: r,
		width:    width,
	}, nil
}

// Render renders markdown content
func (m *MarkdownRenderer) Render(content string) (string, error) {
	return m.renderer.Render(content)
}

// UpdateWidth updates the renderer width
func (m *MarkdownRenderer) UpdateWidth(width int) error {
	// Create a new renderer with updated width
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}

	m.renderer = r
	m.width = width
	return nil
}
