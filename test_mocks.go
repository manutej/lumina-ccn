package main

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

// Mock interfaces for testing
// These are defined here but NOT implemented - implementation will happen in GREEN phase

// MockFileSystemProvider provides filesystem abstraction for testing
type MockFileSystemProvider interface {
	ReadFile(path string) ([]byte, error)
	Stat(path string) (os.FileInfo, error)
	ListDir(path string) ([]string, error)
	Watch(path string) error
	Unwatch(path string) error
}

// MockRipgrepExecutor provides ripgrep command execution abstraction
type MockRipgrepExecutor interface {
	// Execute runs ripgrep with the given query and returns results via channel
	Execute(ctx context.Context, query string) (<-chan json.RawMessage, error)

	// ExecuteWithLimit runs ripgrep with a maximum result count
	ExecuteWithLimit(ctx context.Context, query string, maxResults int) (<-chan json.RawMessage, error)

	// SetBufferSize configures the buffer size for streaming
	SetBufferSize(size int)

	// SetContextLines configures before/after context lines
	SetContextLines(before, after int)

	// SetFileTypes sets file type filters
	SetFileTypes(include, exclude []string)

	// SetCaseSensitive configures case sensitivity
	SetCaseSensitive(sensitive bool)
}

// MockFilesystemWatcher provides file watching abstraction with fsnotify
type MockFilesystemWatcher interface {
	// Watch starts watching a path for changes
	Watch(path string) error

	// Unwatch stops watching a path
	Unwatch(path string) error

	// Changes returns a channel of file change events
	Changes() <-chan string

	// Close stops the watcher and releases resources
	Close() error

	// SetDebounce sets the debounce duration for rapid changes
	SetDebounce(duration time.Duration)

	// SetEventFilter sets which event types to filter out
	SetEventFilter(filterOut []string)

	// SetRecursive enables/disables recursive directory watching
	SetRecursive(recursive bool)
}

// RipgrepMessage represents a JSON message from ripgrep
type RipgrepMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// RipgrepMatch represents a parsed match result
type RipgrepMatch struct {
	Path       string
	LineNumber int
	Text       string
	Context    []string
}

// FuzzyFinder interface (to be implemented)
type FuzzyFinder interface {
	// Navigation
	SetCursor(index int)
	Cursor() int
	SelectedItem() string
	HandleKey(key string)

	// Filtering
	SetFilter(query string)
	Filter() string
	FilteredResults() []string

	// Mode management
	SetMode(mode string)
	Mode() string

	// Jump functionality
	JumpToLetter(letter string) error
}

// GlamourRenderer interface (to be implemented)
type GlamourRenderer interface {
	// Rendering
	Render(markdown string) (string, error)
	RenderSearchResult(result SearchResult) string

	// Theming
	DetectedTheme() string
	SetTheme(theme string)
	GetElementStyle(element string) Style

	// Viewport
	NewViewport(height int) Viewport

	// Configuration
	SetWidth(width int)
	Width() int
	AutoDetectSize()
	WrapEnabled() bool
}

// Viewport interface (to be implemented)
type Viewport interface {
	SetContent(content string)
	SetYOffset(offset int)
	YOffset() int
	NeedsReflow() bool
	Width() int
	Height() int
}

// Style interface (to be implemented)
type Style interface {
	GetForeground() string
	GetBackground() string
}

// RipgrepManager interface for concurrent search management
type RipgrepManager interface {
	Search(ctx context.Context, query string) (<-chan RipgrepMatch, error)
	SetMaxConcurrent(max int)
}

// Constructor functions (declarations only - no implementation)
// These will be implemented in the GREEN phase

// NewFuzzyFinder creates a new fuzzy finder instance
func NewFuzzyFinder(items []string) FuzzyFinder {
	panic("not implemented - RED phase")
}

// NewGlamourRenderer creates a new Glamour renderer instance
func NewGlamourRenderer() GlamourRenderer {
	panic("not implemented - RED phase")
}

// NewRipgrepExecutor, NewRipgrepManager, ParseRipgrepMatch implemented in ripgrep_executor.go
