# Phase 2 API Reference

**Version**: v1.1.0-alpha
**Date**: 2025-10-27
**Status**: Implementation Complete

This document provides comprehensive API documentation for the 5 core components implemented in Phase 2 of the LUMINA project.

---

## Table of Contents

1. [FileWatcher](#1-filewatcher)
2. [RipgrepExecutor](#2-ripgrepexecutor)
3. [FuzzyFinder](#3-fuzzyfinder)
4. [KeybindingHandler](#4-keybindinghandler)
5. [GlamourRenderer](#5-glamourrenderer)

---

## 1. FileWatcher

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/file_watcher.go`

### Interface Definition

```go
type FileWatcher interface {
    // Watch adds a path to the watch list
    Watch(path string) error

    // Unwatch removes a path from the watch list
    Unwatch(path string) error

    // Changes returns a read-only channel for change notifications
    Changes() <-chan string

    // Close stops the file watcher and releases resources
    Close() error

    // SetDebounce configures the debounce duration
    SetDebounce(duration time.Duration)

    // SetEventFilter specifies which event types to filter out
    SetEventFilter(filterOut []string)

    // SetRecursive enables/disables recursive directory watching
    SetRecursive(recursive bool)
}
```

### Implementation

```go
type FileWatcherImpl struct {
    watcher       *fsnotify.Watcher
    changes       chan string
    debounce      time.Duration
    debounceTimer *time.Timer
    debounceMu    sync.Mutex
    pendingPath   string
    closed        bool
    closeMu       sync.Mutex
    wg            sync.WaitGroup
    eventFilter   map[string]bool
    filterMu      sync.RWMutex
    recursive     bool
}
```

### Constructor

```go
func NewFileWatcher() FileWatcher
```

**Defaults**:
- Debounce: 200ms
- Event filter: `chmod` events filtered by default
- Recursive: `true`
- Channel buffer: 100 events

### Usage Example

```go
// Create watcher
watcher := NewFileWatcher()
defer watcher.Close()

// Configure
watcher.SetDebounce(300 * time.Millisecond)
watcher.SetEventFilter([]string{"chmod", "rename"})

// Watch directory
if err := watcher.Watch("/path/to/docs"); err != nil {
    log.Fatal(err)
}

// Listen for changes
for change := range watcher.Changes() {
    fmt.Printf("File changed: %s\n", change)
    // Reload content
}
```

### Methods

#### Watch(path string) error

Adds a path to the watch list. Resolves symlinks and validates directory.

**Parameters**:
- `path`: Absolute or relative path to directory

**Returns**:
- `error`: Returns error if path doesn't exist, is not a directory, or permission denied

**Errors**:
- `"no such file or directory"` - Path doesn't exist
- `"not a directory"` - Path points to a file
- `"permission denied"` - Insufficient permissions

**Behavior**:
- Resolves symlinks via `filepath.EvalSymlinks()`
- Validates path is a directory
- Recursively watches subdirectories if `recursive = true`

#### Unwatch(path string) error

Removes a path from the watch list.

**Parameters**:
- `path`: Path to stop watching

**Returns**:
- `error`: Error from fsnotify.Remove()

#### Changes() <-chan string

Returns a read-only channel that receives file change notifications.

**Returns**:
- `<-chan string`: Channel with changed file paths (debounced)

**Channel Behavior**:
- Buffered (100 events)
- Closed when `Close()` is called
- Debounced (default 200ms)
- Drops events if buffer is full

#### Close() error

Stops the file watcher and releases all resources.

**Returns**:
- `error`: Error from closing fsnotify watcher

**Behavior**:
- Safe to call multiple times
- Waits for event loop goroutine to finish
- Closes the changes channel

#### SetDebounce(duration time.Duration)

Configures the debounce duration for file events.

**Parameters**:
- `duration`: Time to wait before emitting event (e.g., `200 * time.Millisecond`)

**Thread Safety**: Safe for concurrent use

#### SetEventFilter(filterOut []string)

Specifies which event types to filter out.

**Parameters**:
- `filterOut`: Slice of event types to ignore

**Valid Event Types**:
- `"chmod"` - Permission changes
- `"write"` - File modifications
- `"create"` - File/directory creation
- `"remove"` - File/directory deletion
- `"rename"` - File/directory rename

**Example**:
```go
watcher.SetEventFilter([]string{"chmod", "rename"})
```

#### SetRecursive(recursive bool)

Enables/disables recursive directory watching.

**Parameters**:
- `recursive`: `true` to watch subdirectories, `false` for top-level only

### Error Conditions

| Error | Cause | Resolution |
|-------|-------|------------|
| "no such file or directory" | Path doesn't exist | Check path validity |
| "not a directory" | Path points to file | Use directory path |
| "permission denied" | Insufficient permissions | Check file permissions |

### Performance Characteristics

- **Debouncing**: Reduces event flood by 10-100x
- **Channel buffering**: Handles burst of 100 events
- **Recursive watching**: Automatically watches new subdirectories
- **Thread-safe**: All methods safe for concurrent use

### Advanced Usage

**Watching Multiple Paths**:
```go
watcher := NewFileWatcher()
watcher.Watch("/path/to/docs")
watcher.Watch("/path/to/config")

for change := range watcher.Changes() {
    // Handle changes from either path
}
```

**Custom Debounce for Rapid Changes**:
```go
watcher.SetDebounce(500 * time.Millisecond) // Wait 500ms
```

**Filtering Out Noise**:
```go
watcher.SetEventFilter([]string{"chmod"}) // Ignore permission changes
```

---

## 2. RipgrepExecutor

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/ripgrep_executor.go`

### Interface Definition

```go
type MockRipgrepExecutor interface {
    // Execute runs ripgrep and returns results via channel
    Execute(ctx context.Context, query string) (<-chan json.RawMessage, error)

    // ExecuteWithLimit runs ripgrep with max result count
    ExecuteWithLimit(ctx context.Context, query string, maxResults int) (<-chan json.RawMessage, error)

    // SetBufferSize configures buffer size for streaming
    SetBufferSize(size int)

    // SetContextLines configures before/after context
    SetContextLines(before, after int)

    // SetFileTypes sets file type filters
    SetFileTypes(include, exclude []string)

    // SetCaseSensitive configures case sensitivity
    SetCaseSensitive(sensitive bool)
}
```

### Implementation

```go
type RipgrepExecutorImpl struct {
    maxConcurrent int
    sem           chan struct{} // Semaphore for concurrency
    bufferSize    int
    contextBefore int
    contextAfter  int
    includeTypes  []string
    excludeTypes  []string
    caseSensitive bool
}
```

### Constructor

```go
func NewRipgrepExecutor() MockRipgrepExecutor
```

**Defaults**:
- Max concurrent: 4 searches
- Buffer size: 32KB
- Case sensitive: `false`
- Context lines: 0 before, 0 after

### Usage Example

```go
// Create executor
rg := NewRipgrepExecutor()

// Configure
rg.SetCaseSensitive(false)
rg.SetContextLines(2, 2)
rg.SetFileTypes([]string{"md", "txt"}, []string{"log"})

// Execute search
ctx := context.Background()
results, err := rg.Execute(ctx, "TODO")
if err != nil {
    log.Fatal(err)
}

// Process results
for rawMsg := range results {
    var match RipgrepMatch
    json.Unmarshal(rawMsg, &match)
    fmt.Printf("%s:%d: %s\n", match.Path, match.LineNumber, match.Text)
}
```

### Methods

#### Execute(ctx context.Context, query string) (<-chan json.RawMessage, error)

Runs ripgrep with the given query and streams results.

**Parameters**:
- `ctx`: Context for cancellation
- `query`: Search pattern (regex supported)

**Returns**:
- `<-chan json.RawMessage`: Channel with JSON results
- `error`: Validation error if query is invalid

**Behavior**:
- Validates query for shell injection
- Acquires semaphore (max 4 concurrent)
- Streams results asynchronously
- Closes channel when search completes or context canceled

**Query Validation**:
- Empty queries rejected
- Shell metacharacters (`;`, `|`, `` ` ``, `$`) rejected

#### ExecuteWithLimit(ctx context.Context, query string, maxResults int) (<-chan json.RawMessage, error)

Runs ripgrep with a maximum result count.

**Parameters**:
- `ctx`: Context for cancellation
- `query`: Search pattern
- `maxResults`: Maximum number of results to return

**Returns**:
- `<-chan json.RawMessage`: Limited result channel
- `error`: Validation error

**Behavior**:
- Automatically stops consuming after `maxResults`
- More efficient than manual limiting

**Example**:
```go
results, _ := rg.ExecuteWithLimit(ctx, "error", 100) // First 100 matches only
```

#### SetBufferSize(size int)

Configures the buffer size for streaming results.

**Parameters**:
- `size`: Buffer size in bytes (default: 32KB)

**Recommendation**:
- Small files: 16KB
- Large files: 64KB - 128KB

#### SetContextLines(before, after int)

Configures context lines before/after matches.

**Parameters**:
- `before`: Lines before match
- `after`: Lines after match

**Example**:
```go
rg.SetContextLines(2, 2) // Show 2 lines before and after each match
```

#### SetFileTypes(include, exclude []string)

Sets file type filters using ripgrep's type system.

**Parameters**:
- `include`: File types to include (e.g., `[]string{"md", "txt"}`)
- `exclude`: File types to exclude (e.g., `[]string{"log", "tmp"}`)

**Common Types**:
- `"md"` - Markdown
- `"go"` - Go source
- `"py"` - Python
- `"js"` - JavaScript
- `"txt"` - Text files

**Example**:
```go
rg.SetFileTypes([]string{"md"}, []string{}) // Only markdown
rg.SetFileTypes([]string{}, []string{"log"}) // Exclude logs
```

#### SetCaseSensitive(sensitive bool)

Configures case sensitivity.

**Parameters**:
- `sensitive`: `true` for case-sensitive, `false` for case-insensitive

**Default**: `false` (case-insensitive)

### RipgrepManager

Higher-level manager for parsed search results.

```go
type RipgrepManager interface {
    Search(ctx context.Context, query string) (<-chan RipgrepMatch, error)
    SetMaxConcurrent(max int)
}

type RipgrepMatch struct {
    Path       string
    LineNumber int
    Text       string
}
```

**Usage**:
```go
manager := NewRipgrepManager(4)
matches, _ := manager.Search(ctx, "TODO")

for match := range matches {
    fmt.Printf("%s:%d: %s\n", match.Path, match.LineNumber, match.Text)
}
```

### Error Conditions

| Error | Cause | Resolution |
|-------|-------|------------|
| "query cannot be empty" | Empty query string | Provide non-empty query |
| "invalid query - shell metacharacters not allowed" | Shell injection attempt | Remove `;`, `|`, `` ` ``, `$` |

### Performance Characteristics

- **Concurrency**: Max 4 concurrent searches (configurable)
- **Streaming**: Results stream as they arrive
- **Buffer**: 32KB default (adjustable)
- **Cancellation**: Context-aware, stops on cancel

### Security

**Query Validation**:
- Prevents shell injection
- Rejects metacharacters: `;`, `|`, `` ` ``, `$`
- Uses `exec.CommandContext` for safe execution

---

## 3. FuzzyFinder

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/fuzzy_finder_impl.go`

### Implementation

```go
type FuzzyFinderImpl struct {
    items    []string // Original items
    cursor   int      // Current cursor position
    filter   string   // Current filter query
    filtered []string // Filtered results
    mode     string   // Interaction mode
}
```

### Constructor

```go
func NewFuzzyFinderImpl(items []string) *FuzzyFinderImpl
```

**Parameters**:
- `items`: List of items to search/navigate

**Initial State**:
- Cursor: 0
- Filter: "" (empty)
- Filtered: All items
- Mode: "normal"

### Usage Example

```go
// Create fuzzy finder
items := []string{"README.md", "main.go", "config.yaml"}
ff := NewFuzzyFinderImpl(items)

// Navigate
ff.HandleKey("down") // Move to next
ff.HandleKey("j")    // Same as down

// Filter
ff.SetFilter("read")
results := ff.FilteredResults() // ["README.md"]

// Jump
ff.JumpToLetter("m") // Jump to "main.go"

// Get selection
selected := ff.SelectedItem() // Current item
```

### Methods

#### SetCursor(index int)

Sets cursor position if within valid bounds.

**Parameters**:
- `index`: Target cursor position

**Behavior**:
- Validates index is within `[0, len(filtered))`
- Silently ignores out-of-bounds values

#### Cursor() int

Returns current cursor position.

**Returns**:
- `int`: Current cursor index

#### SelectedItem() string

Returns the currently selected item.

**Returns**:
- `string`: Selected item or empty string if no selection

#### HandleKey(key string)

Processes navigation keys with wrapping support.

**Parameters**:
- `key`: Key name

**Supported Keys**:
- `"down"`, `"j"` - Move down (wraps at end)
- `"up"`, `"k"` - Move up (wraps at start)
- `"pagedown"` - Jump down 10 items
- `"pageup"` - Jump up 10 items
- `"home"` - Go to first item
- `"end"` - Go to last item

**Wrapping Behavior**:
```go
ff.HandleKey("down") // At end: wraps to start
ff.HandleKey("up")   // At start: wraps to end
```

#### SetFilter(query string)

Applies a filter query and updates filtered results.

**Parameters**:
- `query`: Search string (case-insensitive)

**Behavior**:
- Resets cursor to 0
- Empty query shows all items
- Case-insensitive substring matching
- Updates `filtered` slice

**Example**:
```go
ff.SetFilter("read")   // Matches "README.md"
ff.SetFilter("READ")   // Also matches (case-insensitive)
ff.SetFilter("")       // Shows all items
```

#### Filter() string

Returns current filter query.

**Returns**:
- `string`: Active filter query

#### FilteredResults() []string

Returns current filtered item list.

**Returns**:
- `[]string`: Slice of items matching filter

#### SetMode(mode string)

Sets current interaction mode.

**Parameters**:
- `mode`: Mode name (e.g., "normal", "search", "help")

#### Mode() string

Returns current interaction mode.

**Returns**:
- `string`: Current mode

#### JumpToLetter(letter string) error

Moves cursor to first item starting with the given letter.

**Parameters**:
- `letter`: Letter to jump to (case-insensitive)

**Returns**:
- `error`: Always `nil` (reserved for future use)

**Behavior**:
- Case-insensitive matching
- Searches filtered results, not original items
- No movement if no match found

**Example**:
```go
// Items: ["apple", "banana", "cherry"]
ff.JumpToLetter("b") // Cursor moves to "banana"
ff.JumpToLetter("c") // Cursor moves to "cherry"
ff.JumpToLetter("z") // No movement (no match)
```

### Performance Characteristics

- **Filtering**: O(n) where n = number of items
- **Navigation**: O(1) constant time
- **Jump**: O(n) linear search
- **Memory**: O(n) for filtered results

### Edge Cases

**Empty List**:
```go
ff := NewFuzzyFinderImpl([]string{})
ff.HandleKey("down") // No panic, no movement
```

**Single Item**:
```go
ff := NewFuzzyFinderImpl([]string{"only.txt"})
ff.HandleKey("down") // Stays on item (wraps to self)
```

**No Filter Matches**:
```go
ff.SetFilter("nonexistent")
ff.FilteredResults() // Returns empty slice []
```

---

## 4. KeybindingHandler

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/keybinding_impl.go`

### Implementation

```go
type KeybindingHandler struct {
    mode string
}
```

### Constructor

```go
func NewKeybindingHandler() *KeybindingHandler
```

**Initial State**:
- Mode: "normal"

### Usage Example

```go
// Create handler
kh := NewKeybindingHandler()

// Handle key events
model := &AppModel{}
result := kh.HandleKey("q", model)
if result == nil {
    // Exit signal
    os.Exit(0)
}

// Mode transitions
kh.HandleKey("?", model) // Enter help mode
fmt.Println(kh.GetMode()) // "help"

kh.HandleKey("escape", model) // Return to normal
fmt.Println(kh.GetMode()) // "normal"
```

### Methods

#### HandleKey(key string, model interface{}) interface{}

Processes a key event based on current mode.

**Parameters**:
- `key`: Key name or character
- `model`: Application model (passed through)

**Returns**:
- `interface{}`: Updated model, or `nil` to signal exit

**Exit Keys**:
- `"q"` - Quit
- `"ctrl+c"` - Force quit

**Mode Transitions**:
- `"?"` → Help mode
- `"/"` → Search mode
- `"escape"` → Normal mode

**Example**:
```go
result := kh.HandleKey("q", model)
if result == nil {
    // User wants to quit
}
```

#### SetMode(mode string)

Updates the current interaction mode.

**Parameters**:
- `mode`: New mode name

**Common Modes**:
- `"normal"` - Default navigation
- `"search"` - Search input
- `"help"` - Help overlay
- `"filtering"` - Filter input

#### GetMode() string

Returns current interaction mode.

**Returns**:
- `string`: Current mode name

### Mode-Aware Behavior

The handler supports context-aware key bindings:

```go
kh := NewKeybindingHandler()

// In normal mode
kh.SetMode("normal")
kh.HandleKey("/", model) // Enters search mode

// In search mode
kh.SetMode("search")
kh.HandleKey("/", model) // Types "/" character
```

### Exit Detection Pattern

```go
for {
    key := readKey()
    result := kh.HandleKey(key, model)

    if result == nil {
        // Exit signal received
        break
    }

    model = result // Continue with updated model
}
```

---

## 5. GlamourRenderer

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/glamour_impl.go`

### Implementation

```go
type GlamourRendererImpl struct {
    renderer *glamour.TermRenderer // Lazy-initialized
    theme    string                // Detected or manual
    width    int                   // Rendering width
}
```

### Constructor

```go
func NewGlamourRendererImpl() *GlamourRendererImpl
```

**Defaults**:
- Theme: Auto-detected (light/dark)
- Width: 80 columns
- Renderer: Lazy-initialized on first `Render()`

### Usage Example

```go
// Create renderer
gr := NewGlamourRendererImpl()

// Render markdown
markdown := "# Hello World\n\nThis is **bold** text."
rendered, err := gr.Render(markdown)
if err != nil {
    log.Fatal(err)
}
fmt.Print(rendered)

// Custom width
gr.SetWidth(120)

// Manual theme
gr.SetTheme("dark")
```

### Methods

#### Render(markdown string) (string, error)

Converts markdown to terminal-formatted output.

**Parameters**:
- `markdown`: Markdown source text

**Returns**:
- `string`: Formatted terminal output
- `error`: Rendering error

**Behavior**:
- Lazy-initializes renderer on first call
- Auto-detects theme if not set
- Word-wraps at configured width
- Applies syntax highlighting

**Example**:
```go
markdown := `
# Title

- Item 1
- Item 2

\`\`\`go
func main() {
    fmt.Println("Hello")
}
\`\`\`
`
rendered, _ := gr.Render(markdown)
```

#### DetectedTheme() string

Returns current theme setting.

**Returns**:
- `string`: `"light"` or `"dark"`

**Theme Detection**:
- Reads `COLORFGBG` environment variable
- Format: `"foreground;background"`
- Background > 7 → light theme
- Default: dark theme

#### SetTheme(theme string)

Manually overrides the theme.

**Parameters**:
- `theme`: `"light"` or `"dark"`

**Behavior**:
- Forces re-initialization on next `Render()`

**Example**:
```go
gr.SetTheme("light") // Force light theme
```

#### SetWidth(width int)

Updates the rendering width.

**Parameters**:
- `width`: Column width (e.g., 80, 120)

**Behavior**:
- Forces re-initialization on next `Render()`
- Affects word wrapping

**Example**:
```go
gr.SetWidth(120) // Wider output
```

#### Width() int

Returns current rendering width.

**Returns**:
- `int`: Width in columns

#### WrapEnabled() bool

Returns whether word wrapping is enabled.

**Returns**:
- `bool`: Always `true` (wrapping always enabled)

#### RenderSearchResult(result SearchResult) string

Formats a search result for display.

**Parameters**:
- `result`: Search result with path and highlight

**Returns**:
- `string`: Styled output with padding

**SearchResult Structure**:
```go
type SearchResult struct {
    Path      string
    Highlight string
}
```

#### GetElementStyle(element string) Style

Returns styling for a specific markdown element.

**Parameters**:
- `element`: Element name (e.g., "heading", "code")

**Returns**:
- `Style`: Style configuration

#### NewViewport(height int) Viewport

Creates a viewport for scrollable content.

**Parameters**:
- `height`: Viewport height in lines

**Returns**:
- `Viewport`: Scrollable content manager

**Viewport Interface**:
```go
type Viewport interface {
    SetContent(content string)
    SetYOffset(offset int)
    YOffset() int
    Width() int
    Height() int
    NeedsReflow() bool
}
```

### Theme Auto-Detection

**COLORFGBG Environment Variable**:
```bash
# Light terminal
export COLORFGBG="0;15"  # Background 15 > 7 → light

# Dark terminal
export COLORFGBG="15;0"  # Background 0 < 7 → dark
```

**Manual Override**:
```go
gr.SetTheme("dark") // Ignore auto-detection
```

### Performance Characteristics

- **Lazy Initialization**: Renderer created only when needed
- **Caching**: Renderer instance reused until theme/width changes
- **Word Wrapping**: Automatic at configured width
- **Memory**: Minimal overhead (~100KB per renderer)

### Advanced Usage

**Custom Width for Terminal Size**:
```go
width, _ := term.GetSize(0)
gr.SetWidth(width - 10) // Leave margin
```

**Conditional Theme**:
```go
if os.Getenv("TERM") == "xterm-256color" {
    gr.SetTheme("dark")
}
```

---

## Common Patterns

### 1. File Watching + Auto-Reload

```go
watcher := NewFileWatcher()
watcher.Watch("/docs")

for change := range watcher.Changes() {
    content, _ := os.ReadFile(change)
    gr := NewGlamourRendererImpl()
    rendered, _ := gr.Render(string(content))
    fmt.Print(rendered)
}
```

### 2. Search + Filter + Display

```go
rg := NewRipgrepExecutor()
results, _ := rg.Execute(ctx, "TODO")

var paths []string
for msg := range results {
    var match RipgrepMatch
    json.Unmarshal(msg, &match)
    paths = append(paths, match.Path)
}

ff := NewFuzzyFinderImpl(paths)
ff.SetFilter("main")
selected := ff.SelectedItem()
```

### 3. Keybinding + Mode Management

```go
kh := NewKeybindingHandler()

for {
    key := readKey()

    switch kh.GetMode() {
    case "normal":
        kh.HandleKey(key, model)
    case "search":
        if key == "escape" {
            kh.SetMode("normal")
        }
    }
}
```

---

## Error Handling

All components follow consistent error handling:

1. **Validation Errors**: Return immediately (e.g., invalid query)
2. **Resource Errors**: Return with context (e.g., file not found)
3. **Runtime Errors**: Log and continue (e.g., watcher errors)

**Example**:
```go
if err := watcher.Watch(path); err != nil {
    if strings.Contains(err.Error(), "permission denied") {
        log.Printf("Cannot watch %s: insufficient permissions", path)
        return
    }
    log.Fatal(err)
}
```

---

## Thread Safety

| Component | Thread-Safe | Notes |
|-----------|-------------|-------|
| FileWatcher | ✅ Yes | All methods use mutexes |
| RipgrepExecutor | ✅ Yes | Semaphore-based concurrency |
| FuzzyFinder | ❌ No | Single-threaded UI component |
| KeybindingHandler | ❌ No | Single-threaded input handler |
| GlamourRenderer | ⚠️ Partial | Lazy init not thread-safe |

**Recommendation**: Use single-threaded event loop for UI components, dedicate goroutines for I/O (watcher, ripgrep).

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| v1.1.0-alpha | 2025-10-27 | Phase 2 implementation complete |
| v1.0.1-alpha | 2025-10-21 | Phase 1.5 (TOC, Search, Keybindings) |

---

**Documentation Generated**: 2025-10-27
**Maintained By**: docs-generator agent
**Next Review**: Phase 3 Planning
