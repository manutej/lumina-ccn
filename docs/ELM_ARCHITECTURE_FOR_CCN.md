# Elm Architecture Guide for CCN (Claude Code Navigator)

**Date**: November 1, 2025
**Purpose**: Deep dive into Elm Architecture as implemented in Bubble Tea and applied to CCN
**Audience**: Developers working on Phase 3 UI integration

---

## Table of Contents

1. [What is The Elm Architecture (TEA)?](#what-is-the-elm-architecture)
2. [Core Components](#core-components)
3. [Bubble Tea Implementation](#bubble-tea-implementation)
4. [CCN Current Architecture](#ccn-current-architecture)
5. [Message Flow Analysis](#message-flow-analysis)
6. [Identified Blockers](#identified-blockers)
7. [Solutions & Best Practices](#solutions--best-practices)
8. [Phase 3 Week 1 Analysis](#phase-3-week-1-analysis)

---

## What is The Elm Architecture?

The Elm Architecture (TEA) is a pattern for building interactive applications with three core parts:

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│     ┌───────────┐        ┌──────────┐        ┌──────────┐ │
│     │           │        │          │        │          │ │
│     │   Model   │───────▶│   View   │───────▶│   HTML   │ │
│     │  (State)  │        │          │        │          │ │
│     │           │        │          │        │          │ │
│     └─────▲─────┘        └──────────┘        └────┬─────┘ │
│           │                                        │       │
│           │                                        │       │
│           │              ┌──────────┐             │       │
│           │              │          │             │       │
│           └──────────────│  Update  │◀────────────┘       │
│                          │          │                     │
│                          └──────────┘                     │
│                               ▲                            │
│                               │                            │
│                           Messages                         │
│                      (User Interactions)                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

1. **Unidirectional Data Flow**: Data flows in one direction (Model → View → Messages → Update → Model)
2. **Single Source of Truth**: All state lives in the Model
3. **Pure Functions**: Update and View are pure (no side effects)
4. **Explicit State**: All state transitions are explicit via messages
5. **Predictable**: Same input always produces same output

---

## Core Components

### 1. Model (State)

The Model represents **all** application state.

**Elm Example**:
```elm
type alias Model =
  { content : String
  , count : Int
  }
```

**CCN Example** (`model.go`):
```go
type AppModel struct {
    // Window state
    width  int
    height int
    ready  bool

    // Navigation state
    currentView ViewMode
    currentMode UIMode

    // File state
    selectedFile   string
    viewerContent  string

    // Fuzzy finder state
    finderActive   bool
    finderInput    string
    finderSelected int
    fuzzyFinder    *FuzzyFinderImpl
}
```

**Best Practices**:
- ✅ Store all UI state in Model
- ✅ Use explicit types (enums/constants)
- ✅ Make state serializable
- ❌ Don't store computed values (derive them in View)
- ❌ Don't store references to UI components

### 2. Messages (Events)

Messages represent **all possible state changes**.

**Elm Example**:
```elm
type Msg
  = Increment
  | Decrement
  | Change String
```

**CCN Example** (`main.go`):
```go
// Phase 3 Message Types for Fuzzy Finder

type FinderActivatedMsg struct{}

type FinderInputMsg struct {
    input string
}

type FinderSelectionMsg struct {
    filePath string
}

type FinderCanceledMsg struct{}
```

**Bubble Tea Built-in Messages**:
```go
tea.KeyMsg         // Keyboard input
tea.MouseMsg       // Mouse events
tea.WindowSizeMsg  // Terminal resize
```

**Best Practices**:
- ✅ Create custom message types for each logical action
- ✅ Use struct fields to carry data with messages
- ✅ Make messages self-documenting (clear names)
- ❌ Don't reuse messages for different purposes
- ❌ Don't put complex logic in message creation

### 3. Update (State Transitions)

The Update function processes messages and returns new state.

**Elm Example**:
```elm
update : Msg -> Model -> Model
update msg model =
  case msg of
    Increment ->
      model + 1

    Decrement ->
      model - 1

    Change newContent ->
      { model | content = newContent }
```

**CCN Example** (`main.go`):
```go
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        if m.finderActive {
            switch msg.String() {
            case "enter":
                // Select item
                selected := m.fuzzyFinder.SelectedItem()
                if selected != "" {
                    m.finderActive = false
                    m.loadFileContent(selected)
                }
                return m, nil
            }
        }
    }

    return m, cmd
}
```

**Best Practices**:
- ✅ Keep Update pure (no side effects directly)
- ✅ Return new Model (don't mutate existing)
- ✅ Use Commands (tea.Cmd) for side effects
- ✅ Handle all possible messages
- ❌ Don't make HTTP calls directly in Update
- ❌ Don't mutate global state

### 4. View (Rendering)

The View function renders the current state to UI.

**Elm Example**:
```elm
view : Model -> Html Msg
view model =
  div []
    [ button [ onClick Decrement ] [ text "-" ]
    , div [] [ text (String.fromInt model) ]
    , button [ onClick Increment ] [ text "+" ]
    ]
```

**CCN Example** (`main.go`):
```go
func (m AppModel) View() string {
    if !m.ready {
        return "Loading..."
    }

    // Render base view
    baseView := lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderTitle(),
        m.renderContent(),
        m.renderStatus(),
    )

    // Overlay finder modal if active
    if m.finderActive {
        finderOverlay := m.renderFinderModal()
        return lipgloss.Place(
            m.width, m.height,
            lipgloss.Center, lipgloss.Center,
            finderOverlay,
        )
    }

    return baseView
}
```

**Best Practices**:
- ✅ View is pure (same Model → same output)
- ✅ Derive computed values in View
- ✅ Use composition (break into smaller functions)
- ❌ Don't modify Model in View
- ❌ Don't call external APIs in View

---

## Bubble Tea Implementation

Bubble Tea is a Go implementation of The Elm Architecture for terminal UIs.

### Key Differences from Elm

| Aspect | Elm | Bubble Tea |
|--------|-----|------------|
| **Language** | Elm (functional) | Go (imperative) |
| **Rendering** | HTML/DOM | Terminal (ANSI) |
| **Commands** | `Cmd Msg` | `tea.Cmd` |
| **Subscriptions** | `Sub Msg` | `tea.Cmd` (polling) |
| **Type Safety** | Compile-time (strong) | Runtime (interface{}) |

### Bubble Tea Architecture

```go
// Init: Initialize model and return initial command
func (m Model) Init() tea.Cmd {
    return nil // or return a command
}

// Update: Handle messages, return new model and command
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard
    case tea.MouseMsg:
        // Handle mouse
    }
    return m, nil
}

// View: Render current state to string
func (m Model) View() string {
    return "Hello, World!"
}
```

### Commands in Bubble Tea

Commands represent **side effects** that return messages.

**Example: Loading a file**
```go
// Command that loads a file and returns a message
func loadFileCmd(path string) tea.Cmd {
    return func() tea.Msg {
        content, err := os.ReadFile(path)
        if err != nil {
            return FileLoadErrorMsg{err}
        }
        return FileLoadedMsg{content: string(content)}
    }
}

// Usage in Update
case tea.KeyMsg:
    if msg.String() == "enter" {
        selectedPath := m.getSelectedPath()
        return m, loadFileCmd(selectedPath)
    }
```

**Built-in Commands**:
```go
tea.Quit                    // Exit the program
tea.Batch(cmd1, cmd2, ...)  // Run multiple commands
tea.Sequence(cmd1, cmd2)    // Run commands in order
```

---

## CCN Current Architecture

### Model Structure Analysis

```go
type AppModel struct {
    // ✅ GOOD: Window dimensions (derived from tea.WindowSizeMsg)
    width  int
    height int
    ready  bool

    // ✅ GOOD: Explicit state enums
    currentView    ViewMode      // FileTreeView | ViewerView | PreviewView
    currentMode    UIMode         // NormalMode | FinderMode | SearchMode
    selectionMode  SelectionMode  // SelectionInactive | SelectionCharacter | ...

    // ✅ GOOD: Navigation state
    rootPath    string
    currentPath string

    // ✅ GOOD: UI components (Bubble Tea bubbles)
    fileList         list.Model
    viewer           viewport.Model
    preview          viewport.Model
    markdownRenderer *utils.MarkdownRenderer

    // ✅ GOOD: Content state
    selectedFile    string
    viewerContent   string
    renderedContent string

    // ⚠️ MIXED: Helper managers (custom components)
    keyBindings     *KeyBindings
    clipboard       *ClipboardManager
    colorManager    *ColorManager
    tableOfContents *TableOfContents

    // ✅ GOOD: Fuzzy finder state (Phase 3)
    finderActive   bool
    finderInput    string
    finderSelected int
    fuzzyFinder    *FuzzyFinderImpl

    // ⚠️ POTENTIAL BLOCKER: Mouse drag state
    mouseDragActive    bool
    mouseDragStartLine int
    mouseDragStartCol  int
}
```

### Message Flow

**Current Flow** (Phase 3 Week 1):
```
User presses "/"
    ↓
tea.KeyMsg{String: "/"} received
    ↓
Update() checks: msg.String() == "/" && !m.finderActive
    ↓
Initialize fuzzy finder:
    - Load markdown files (if not loaded)
    - Create FuzzyFinderImpl
    - Set finderActive = true
    - Set currentMode = FinderMode
    ↓
Return (m, nil)
    ↓
View() checks: if m.finderActive
    ↓
Render finder modal overlay
    ↓
User types "test"
    ↓
tea.KeyMsg{String: "t"} received
    ↓
Update() checks: if m.finderActive
    ↓
Append to finderInput: "t"
Call fuzzyFinder.SetFilter("t")
    ↓
Return (m, nil)
    ↓
View() re-renders with filtered results
```

### Update Function Flow Chart

```
Update(msg tea.Msg) receives message
    │
    ├─▶ tea.WindowSizeMsg?
    │   └─▶ Update dimensions, return (m, nil)
    │
    ├─▶ tea.MouseMsg?
    │   └─▶ handleMouseEvent(msg)
    │
    ├─▶ tea.KeyMsg?
        │
        ├─▶ "?" key?
        │   └─▶ Toggle help overlay, return (m, nil)
        │
        ├─▶ Help overlay active?
        │   └─▶ Handle help keys, return (m, nil)
        │
        ├─▶ Finder active?
        │   │
        │   ├─▶ "esc" → Close finder
        │   ├─▶ "enter" → Select item, load file
        │   ├─▶ "up"/"down" → Navigate
        │   └─▶ Character → Add to input, filter
        │
        ├─▶ "/" key && !finderActive?
        │   └─▶ Activate finder, return (m, nil)
        │
        └─▶ Use keybindings to find action
            └─▶ Execute action, return (m, cmd)
```

---

## Message Flow Analysis

### Good Patterns in CCN

✅ **Explicit State Modes**
```go
type UIMode int
const (
    NormalMode UIMode = iota
    FinderMode
    SearchMode
)
```
This follows Elm's pattern of using sum types to make states explicit.

✅ **Message-Driven State Changes**
```go
if msg.String() == "/" && !m.finderActive {
    m.finderActive = true
    m.currentMode = FinderMode
    return m, nil
}
```
State changes happen only in response to messages.

✅ **Pure View Function**
```go
func (m AppModel) View() string {
    // No side effects, just rendering
    if m.finderActive {
        return m.renderFinderModal()
    }
    return baseView
}
```

### Potential Issues

⚠️ **Issue 1: Direct State Mutation in Update**
```go
// Current code (anti-pattern in strict TEA)
m.finderInput += msg.String()
m.fuzzyFinder.SetFilter(m.finderInput)
```

**Problem**: We're mutating `m` directly instead of creating new model.

**Why it works in Go/Bubble Tea**: Go uses pass-by-value for structs, so `m` is already a copy.

**Better Pattern** (more explicit):
```go
newModel := m
newModel.finderInput += msg.String()
newModel.fuzzyFinder.SetFilter(newModel.finderInput)
return newModel, nil
```

⚠️ **Issue 2: Side Effects in Update**
```go
if selected != "" {
    m.finderActive = false
    m.loadFileContent(selected)  // ⚠️ Side effect!
}
```

**Problem**: `loadFileContent()` reads from disk (I/O side effect).

**Elm Way**: Use a Command
```go
// Define message type
type FileLoadedMsg struct {
    content string
    err     error
}

// Create command
func loadFileCmd(path string) tea.Cmd {
    return func() tea.Msg {
        content, err := os.ReadFile(path)
        return FileLoadedMsg{content: string(content), err: err}
    }
}

// In Update
case "enter":
    selected := m.fuzzyFinder.SelectedItem()
    if selected != "" {
        m.finderActive = false
        return m, loadFileCmd(selected)
    }
```

**Why current pattern works**: File I/O is fast enough in CCN context, but using Commands would be more robust for async operations.

⚠️ **Issue 3: Lazy Loading in Update**
```go
if len(m.markdownFiles) == 0 {
    m.markdownFiles = findMarkdownFiles(m.rootPath)
}
```

**Problem**: Blocking I/O during Update (could freeze UI for large directories).

**Better Pattern**: Use Command with loading state
```go
type MarkdownFilesLoadedMsg struct {
    files []string
}

func loadMarkdownFilesCmd(rootPath string) tea.Cmd {
    return func() tea.Msg {
        files := findMarkdownFiles(rootPath)
        return MarkdownFilesLoadedMsg{files: files}
    }
}

// In Update
case "/":
    if len(m.markdownFiles) == 0 {
        m.currentMode = LoadingMode
        return m, loadMarkdownFilesCmd(m.rootPath)
    }
    // ... activate finder
```

---

## Identified Blockers

### 🚧 Blocker 1: Synchronous File I/O in Update

**Symptom**: UI freezes when loading large directories or files.

**Root Cause**: `findMarkdownFiles()` and `loadFileContent()` block the Update function.

**Impact**: Violates TEA principle of keeping Update fast and pure.

**Solution**: Use Commands for all I/O operations.

**Example Fix**:
```go
// Message types
type MarkdownFilesLoadingMsg struct{}
type MarkdownFilesLoadedMsg struct { files []string }
type FileLoadingMsg struct { path string }
type FileLoadedMsg struct { content string; err error }

// Commands
func loadMarkdownFilesCmd(rootPath string) tea.Cmd {
    return func() tea.Msg {
        return MarkdownFilesLoadingMsg{}
        // In background, send MarkdownFilesLoadedMsg when done
    }
}

// Update handles loading states
case MarkdownFilesLoadingMsg:
    m.loadingFiles = true
    return m, nil

case MarkdownFilesLoadedMsg:
    m.markdownFiles = msg.files
    m.loadingFiles = false
    m.finderActive = true
    return m, nil
```

### 🚧 Blocker 2: Complex State in Helper Objects

**Symptom**: State split between AppModel and helper objects (FuzzyFinderImpl, KeyBindings, etc.)

**Root Cause**: Go's object-oriented style conflicts with TEA's single Model principle.

**Impact**: Harder to reason about state, potential for inconsistencies.

**Current Pattern**:
```go
type AppModel struct {
    fuzzyFinder *FuzzyFinderImpl  // Has its own state
}

// State is split
m.finderInput = "test"           // In AppModel
m.fuzzyFinder.SetFilter("test")  // In FuzzyFinderImpl
```

**Better Pattern**: Flatten state into AppModel
```go
type AppModel struct {
    // Finder state (all in one place)
    finderInput    string
    finderCursor   int
    finderItems    []string
    finderFiltered []string
}

// Logic as pure functions (no state)
func filterItems(items []string, query string) []string {
    // Pure filtering logic
}

// In Update
case tea.KeyMsg:
    m.finderInput += msg.String()
    m.finderFiltered = filterItems(m.finderItems, m.finderInput)
```

### 🚧 Blocker 3: Mouse State Management

**Symptom**: Mouse drag state (`mouseDragActive`, `mouseDragStartLine`, etc.) is complex.

**Root Cause**: Mouse interactions are inherently stateful (drag = multi-message sequence).

**Impact**: Drag state can get out of sync if messages are missed.

**Current Pattern**:
```go
case tea.MouseMsg:
    if msg.Type == tea.MouseLeft && msg.Action == tea.MouseActionPress {
        m.mouseDragActive = true
        m.mouseDragStartLine = msg.Y
    }
    if msg.Type == tea.MouseLeft && msg.Action == tea.MouseActionRelease {
        m.mouseDragActive = false
    }
```

**Potential Issue**: If Release message is missed, `mouseDragActive` stays true forever.

**Better Pattern**: Use message sequences with timeout
```go
type MouseDragStartedMsg struct { line, col int }
type MouseDragEndedMsg struct {}
type MouseDragTimeoutMsg struct {}

// Command that sends timeout message after 5 seconds
func dragTimeoutCmd() tea.Cmd {
    return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
        return MouseDragTimeoutMsg{}
    })
}

// In Update
case tea.MouseMsg:
    if msg.Action == tea.MouseActionPress {
        m.mouseDragActive = true
        return m, dragTimeoutCmd()  // Auto-cancel after 5s
    }

case MouseDragTimeoutMsg:
    m.mouseDragActive = false  // Recover from stuck state
```

### 🚧 Blocker 4: Complex Navigation State

**Symptom**: Multiple overlapping states (currentView, currentMode, finderActive, showHelp).

**Root Cause**: State explosion from combinations.

**Impact**: Hard to reason about which state takes precedence.

**Current Pattern**:
```go
// What happens if both are true?
if m.showHelp { /* render help */ }
if m.finderActive { /* render finder */ }
```

**Better Pattern**: Use a single state machine
```go
type AppState int
const (
    StateNormal AppState = iota
    StateHelp
    StateFinder
    StateSearch
    StateLoading
)

type AppModel struct {
    state AppState
    // ... other fields
}

// In View
switch m.state {
case StateNormal:
    return m.renderNormal()
case StateHelp:
    return m.renderHelp()
case StateFinder:
    return m.renderFinder()
}
```

---

## Solutions & Best Practices

### Solution 1: Embrace Commands for I/O

**Pattern**: Any I/O operation should return a Command.

**Example**:
```go
// ❌ Bad: Synchronous I/O in Update
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    content, _ := os.ReadFile(path)  // Blocks!
    m.viewerContent = string(content)
    return m, nil
}

// ✅ Good: Asynchronous I/O via Command
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, loadFileCmd(path)
}

func loadFileCmd(path string) tea.Cmd {
    return func() tea.Msg {
        content, err := os.ReadFile(path)
        return FileLoadedMsg{content, err}
    }
}
```

### Solution 2: Flatten State

**Pattern**: Keep all state in AppModel, use pure functions for logic.

**Example**:
```go
// ❌ Bad: State in helper object
type FuzzyFinder struct {
    items    []string
    cursor   int
    filtered []string
}

// ✅ Good: State in AppModel, logic in functions
type AppModel struct {
    finderItems    []string
    finderCursor   int
    finderFiltered []string
}

func filterItems(items []string, query string) []string {
    // Pure function
}
```

### Solution 3: Use State Machines

**Pattern**: Model complex states as explicit state machines.

**Example**:
```go
type AppState int
const (
    StateNormal AppState = iota
    StateLoading
    StateFinder
    StateSearch
    StateHelp
)

// State transitions are explicit
func (m AppModel) transitionTo(newState AppState) AppModel {
    m.state = newState
    // Clean up state from previous mode
    if newState != StateFinder {
        m.finderActive = false
    }
    return m
}
```

### Solution 4: Implement Message Batching

**Pattern**: When multiple state changes happen together, batch them.

**Example**:
```go
// ❌ Bad: Multiple sequential updates
m.finderActive = true
m.currentMode = FinderMode
m.finderInput = ""

// ✅ Good: Single update function
func (m AppModel) activateFinder() AppModel {
    return AppModel{
        ...m,  // Copy existing state
        finderActive: true,
        currentMode: FinderMode,
        finderInput: "",
    }
}
```

---

## Phase 3 Week 1 Analysis

### What We Did Right ✅

1. **Explicit State Modes**: Created `UIMode` enum for clear state representation
2. **Message Types**: Defined custom message types (FinderActivatedMsg, etc.)
3. **Pure View**: `renderFinderModal()` is a pure function
4. **Unidirectional Flow**: Messages → Update → Model → View works correctly

### What Could Be Improved ⚠️

1. **Synchronous I/O**: `findMarkdownFiles()` blocks Update
2. **Helper Object State**: `FuzzyFinderImpl` maintains its own state
3. **Direct Mutation**: `m.finderInput +=` instead of creating new model

### Recommended Refactoring

**Priority 1: Use Commands for File Loading**
```go
// Current (synchronous)
if len(m.markdownFiles) == 0 {
    m.markdownFiles = findMarkdownFiles(m.rootPath)
}

// Improved (asynchronous)
if len(m.markdownFiles) == 0 {
    return m, loadMarkdownFilesCmd(m.rootPath)
}
```

**Priority 2: Flatten Fuzzy Finder State**
```go
// Current (split state)
type AppModel struct {
    finderInput string
    fuzzyFinder *FuzzyFinderImpl  // Has cursor, filtered, etc.
}

// Improved (flat state)
type AppModel struct {
    finderInput    string
    finderCursor   int
    finderFiltered []string
}
```

**Priority 3: Implement State Machine**
```go
// Current (boolean flags)
if m.finderActive { ... }
if m.showHelp { ... }

// Improved (explicit states)
switch m.state {
case StateFinder: ...
case StateHelp: ...
}
```

### Impact on Week 2 (Ripgrep Search)

**Lessons Learned**:
1. Use Commands for Ripgrep execution (it's async by nature)
2. Stream results via messages, not callback functions
3. Keep search state in AppModel, not in RipgrepExecutor

**Recommended Pattern for Week 2**:
```go
// Message types
type SearchStartedMsg struct { query string }
type SearchResultMsg struct { result RipgrepResult }
type SearchCompletedMsg struct {}

// Command that streams results
func searchCmd(query string) tea.Cmd {
    return func() tea.Msg {
        // Start ripgrep
        // Send SearchResultMsg for each result
        // Send SearchCompletedMsg when done
    }
}

// Update handles streaming results
case SearchResultMsg:
    m.searchResults = append(m.searchResults, msg.result)
    return m, nil  // View re-renders with new result
```

---

## Best Practices Summary

### ✅ Do This

1. **Keep Update Pure**: No I/O, no global state mutations
2. **Use Commands for Side Effects**: File I/O, network, timers
3. **Flatten State**: All state in Model, logic in pure functions
4. **Explicit State Machines**: Use enums instead of boolean flags
5. **Message-Driven**: Every state change via a message
6. **Pure View**: Same Model → Same output, no side effects
7. **Type-Safe Messages**: Create custom message types
8. **Handle All Cases**: Pattern match on all possible messages

### ❌ Don't Do This

1. **Don't Block Update**: No synchronous I/O
2. **Don't Mutate Global State**: No package-level variables
3. **Don't Split State**: Keep state in Model, not helpers
4. **Don't Skip Messages**: Handle every message type
5. **Don't Put Logic in View**: View only renders
6. **Don't Use Callbacks**: Use messages instead
7. **Don't Ignore Errors**: Handle via error messages
8. **Don't Use Booleans for State**: Use enums

---

## Conclusion

The Elm Architecture is at the heart of Bubble Tea and provides a robust pattern for building interactive TUIs. CCN's Phase 3 Week 1 implementation follows TEA principles well, but there are opportunities to improve:

1. **Commands**: Use for all I/O operations
2. **Flat State**: Move helper object state into AppModel
3. **State Machines**: Replace boolean flags with enums
4. **Async Operations**: Embrace message-driven async patterns

These improvements will make Week 2 (Ripgrep Search) implementation more robust, especially for streaming search results and handling long-running operations.

---

## Further Reading

- [Elm Architecture Guide](https://guide.elm-lang.org/architecture/)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Charm Libraries](https://github.com/charmbracelet)

---

**Status**: Ready for Week 2 with deeper architectural understanding 🚀
