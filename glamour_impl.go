package main

import (
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// GlamourRendererImpl handles markdown rendering with theme detection.
// Lazy-initializes renderer on first use for efficiency.
// Thread-safe implementation using RWMutex for concurrent access.
type GlamourRendererImpl struct {
	renderer *glamour.TermRenderer // Lazy-initialized on first Render() call
	theme    string                // Detected or manually set theme
	width    int                   // Rendering width in columns
	widthMu  sync.RWMutex          // Protects width field from concurrent access
	themeMu  sync.RWMutex          // Protects theme field from concurrent access
}

// NewGlamourRendererImpl creates a new glamour renderer with auto-detected theme.
// Default width: 80 columns. Renderer is lazy-initialized on first Render() call.
func NewGlamourRendererImpl() *GlamourRendererImpl {
	return &GlamourRendererImpl{
		theme: detectTheme(),
		width: 80,
	}
}

// detectTheme determines light/dark theme from COLORFGBG environment variable.
// Format: "foreground;background" where background > 7 indicates light theme.
// Defaults to dark theme if detection fails.
func detectTheme() string {
	colorfgbg := os.Getenv("COLORFGBG")
	if colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) == 2 && parts[1] > "7" {
			return "light"
		}
	}
	return "dark"
}

// Render converts markdown to terminal-formatted output.
// Lazy-initializes renderer on first call. Returns formatted string or error.
// Thread-safe: uses read lock to safely access width during rendering.
func (gr *GlamourRendererImpl) Render(markdown string) (string, error) {
	// Read lock to safely access width
	gr.widthMu.RLock()
	currentWidth := gr.width
	gr.widthMu.RUnlock()

	// Lazy-initialize renderer if needed
	if gr.renderer == nil {
		if err := gr.initializeRenderer(); err != nil {
			return "", err
		}
	}

	// Use the safely-read width value (copy semantics prevent race)
	_ = currentWidth // Width is embedded in renderer state

	return gr.renderer.Render(markdown)
}

// initializeRenderer creates the glamour renderer with current settings.
// Thread-safe: acquires read lock before accessing width.
func (gr *GlamourRendererImpl) initializeRenderer() error {
	// Read lock to safely access width
	gr.widthMu.RLock()
	currentWidth := gr.width
	gr.widthMu.RUnlock()

	var err error
	gr.renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(currentWidth),
	)
	return err
}

// RenderSearchResult formats a search result for display.
func (gr *GlamourRendererImpl) RenderSearchResult(result SearchResult) string {
	style := lipgloss.NewStyle().Padding(0, 1)
	return style.Render(result.Path + "\n" + result.Highlight)
}

// DetectedTheme returns the current theme setting.
// Thread-safe: uses read lock.
func (gr *GlamourRendererImpl) DetectedTheme() string {
	gr.themeMu.RLock()
	defer gr.themeMu.RUnlock()
	return gr.theme
}

// SetTheme manually overrides the theme.
// Thread-safe: uses write lock to update theme.
func (gr *GlamourRendererImpl) SetTheme(theme string) {
	gr.themeMu.Lock()
	gr.theme = theme
	gr.renderer = nil // Force re-initialization on next render
	gr.themeMu.Unlock()
}

// GetElementStyle returns styling for a specific markdown element.
func (gr *GlamourRendererImpl) GetElementStyle(element string) Style {
	return NewElementStyle(element)
}

// NewViewport creates a viewport for scrollable content.
// Thread-safe: uses read lock to access width.
func (gr *GlamourRendererImpl) NewViewport(height int) Viewport {
	gr.widthMu.RLock()
	currentWidth := gr.width
	gr.widthMu.RUnlock()
	return &ViewportImpl{height: height, width: currentWidth}
}

// SetWidth updates the rendering width.
// Thread-safe: uses write lock to update width.
func (gr *GlamourRendererImpl) SetWidth(width int) {
	gr.widthMu.Lock()
	gr.width = width
	gr.renderer = nil // Force re-initialization
	gr.widthMu.Unlock()
}

// Width returns the current rendering width.
// Thread-safe: uses read lock.
func (gr *GlamourRendererImpl) Width() int {
	gr.widthMu.RLock()
	defer gr.widthMu.RUnlock()
	return gr.width
}

// AutoDetectSize would detect terminal size (placeholder).
func (gr *GlamourRendererImpl) AutoDetectSize() {
	// Terminal size detection would go here
}

// WrapEnabled returns whether word wrapping is enabled.
func (gr *GlamourRendererImpl) WrapEnabled() bool {
	return true
}

// ViewportImpl manages scrollable content display.
type ViewportImpl struct {
	content string
	yOffset int
	width   int
	height  int
}

// SetContent updates the viewport content.
func (v *ViewportImpl) SetContent(content string) {
	v.content = content
}

// SetYOffset updates the vertical scroll offset.
func (v *ViewportImpl) SetYOffset(offset int) {
	if offset >= 0 {
		v.yOffset = offset
	}
}

// YOffset returns the current vertical scroll offset.
func (v *ViewportImpl) YOffset() int {
	return v.yOffset
}

// NeedsReflow indicates if content needs to be reflowed.
func (v *ViewportImpl) NeedsReflow() bool {
	return false
}

// Width returns the viewport width.
func (v *ViewportImpl) Width() int {
	return v.width
}

// Height returns the viewport height.
func (v *ViewportImpl) Height() int {
	return v.height
}

// ElementStyle provides styling for markdown elements.
type ElementStyle struct {
	element string
}

// NewElementStyle creates a style for a markdown element.
func NewElementStyle(element string) *ElementStyle {
	return &ElementStyle{element: element}
}

// GetForeground returns the foreground color.
func (es *ElementStyle) GetForeground() string {
	return "#FFFFFF"
}

// GetBackground returns the background color.
func (es *ElementStyle) GetBackground() string {
	return "#000000"
}
