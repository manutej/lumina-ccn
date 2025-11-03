# CCN Blockers Analysis - Phase 3

**Date**: November 1, 2025
**Based on**: Elm Architecture deep dive research
**Status**: Blockers identified with solutions

---

## Executive Summary

After researching the Elm Architecture (TEA) that underlies Bubble Tea, I've identified **4 architectural blockers** in our Phase 3 Week 1 implementation that could cause issues in Week 2 and beyond.

**Good News**: All blockers have clear solutions and won't require major refactoring.

---

## 🚧 Blocker 1: Synchronous File I/O in Update

### Problem

```go
// Current code in Update()
if msg.String() == "/" && !m.finderActive {
    if len(m.markdownFiles) == 0 {
        m.markdownFiles = findMarkdownFiles(m.rootPath)  // ⚠️ Blocks UI!
    }
    m.fuzzyFinder = NewFuzzyFinderImpl(m.markdownFiles)
    m.finderActive = true
    return m, nil
}
```

**Issue**: `findMarkdownFiles()` walks the entire directory tree synchronously, blocking the UI thread.

**Impact**:
- UI freezes for large projects (1000+ markdown files)
- Violates TEA principle of keeping Update fast
- No loading indicator shown to user

### Solution: Use Commands

```go
// Define message types
type MarkdownFilesLoadingMsg struct{}
type MarkdownFilesLoadedMsg struct {
    files []string
}

// Create command
func loadMarkdownFilesCmd(rootPath string) tea.Cmd {
    return func() tea.Msg {
        files := findMarkdownFiles(rootPath)
        return MarkdownFilesLoadedMsg{files: files}
    }
}

// In Update
case "/":
    if len(m.markdownFiles) == 0 {
        m.loadingFiles = true
        return m, tea.Batch(
            loadMarkdownFilesCmd(m.rootPath),
            spinner.Tick,  // Show loading spinner
        )
    }
    m.finderActive = true
    return m, nil

case MarkdownFilesLoadedMsg:
    m.markdownFiles = msg.files
    m.loadingFiles = false
    m.finderActive = true
    m.fuzzyFinder = NewFuzzyFinderImpl(msg.files)
    return m, nil
```

**Benefits**:
- ✅ Non-blocking UI
- ✅ Loading indicator possible
- ✅ User can cancel operation
- ✅ Follows TEA principles

---

## 🚧 Blocker 2: State Split Between AppModel and Helper Objects

### Problem

```go
type AppModel struct {
    finderInput    string           // State in AppModel
    finderSelected int              // State in AppModel
    fuzzyFinder    *FuzzyFinderImpl // ⚠️ More state in helper!
}

type FuzzyFinderImpl struct {
    items    []string  // State duplicated
    cursor   int       // State duplicated
    filter   string    // State duplicated
    filtered []string  // Derived state
}
```

**Issue**: State is split between AppModel and FuzzyFinderImpl, violating TEA's "single Model" principle.

**Impact**:
- State can get out of sync
- Harder to serialize/deserialize state
- Confusing what's in Model vs helper
- Violates single source of truth

### Solution: Flatten State

```go
// Model: All state in one place
type AppModel struct {
    // Finder state (flat)
    finderActive   bool
    finderInput    string
    finderCursor   int
    finderItems    []string
}

// Pure helper functions (no state)
func filterItems(items []string, query string) []string {
    if query == "" {
        return items
    }

    var filtered []string
    lowerQuery := strings.ToLower(query)
    for _, item := range items {
        if strings.Contains(strings.ToLower(item), lowerQuery) {
            filtered = append(filtered, item)
        }
    }
    return filtered
}

func navigateCursor(cursor, delta, listLen int) int {
    newCursor := cursor + delta
    if newCursor < 0 {
        return listLen - 1  // Wrap to end
    }
    if newCursor >= listLen {
        return 0  // Wrap to start
    }
    return newCursor
}

// In Update
case tea.KeyMsg:
    if m.finderActive {
        switch msg.String() {
        case "up":
            filtered := filterItems(m.finderItems, m.finderInput)
            m.finderCursor = navigateCursor(m.finderCursor, -1, len(filtered))

        default:
            m.finderInput += msg.String()
        }
    }

// In View
func (m AppModel) renderFinderModal() string {
    filtered := filterItems(m.finderItems, m.finderInput)

    for i, item := range filtered {
        if i == m.finderCursor {
            // Render selected
        }
    }
}
```

**Benefits**:
- ✅ All state in Model (single source of truth)
- ✅ Easy to serialize entire state
- ✅ Pure functions are testable
- ✅ Follows TEA principles

---

## 🚧 Blocker 3: Complex Mouse Drag State

### Problem

```go
type AppModel struct {
    mouseDragActive    bool  // Can get stuck true
    mouseDragStartLine int
    mouseDragStartCol  int
}

// In Update
case tea.MouseMsg:
    if msg.Action == tea.MouseActionPress {
        m.mouseDragActive = true
        m.mouseDragStartLine = msg.Y
    }
    if msg.Action == tea.MouseActionRelease {
        m.mouseDragActive = false  // ⚠️ What if this message is missed?
    }
```

**Issue**: If the MouseActionRelease message is missed (terminal glitch, user switches windows), `mouseDragActive` stays true forever.

**Impact**:
- Stuck in drag mode
- Can't select text properly
- Requires app restart to fix

### Solution: Add Timeout Recovery

```go
// Message types
type MouseDragTimeoutMsg struct{}

// Command with timeout
func dragTimeoutCmd() tea.Cmd {
    return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
        return MouseDragTimeoutMsg{}
    })
}

// In Update
case tea.MouseMsg:
    if msg.Action == tea.MouseActionPress {
        m.mouseDragActive = true
        m.mouseDragStartLine = msg.Y
        return m, dragTimeoutCmd()  // Auto-cancel after 5s
    }
    if msg.Action == tea.MouseActionRelease {
        m.mouseDragActive = false
        return m, nil
    }

case MouseDragTimeoutMsg:
    if m.mouseDragActive {
        // Recover from stuck drag state
        m.mouseDragActive = false
        m.clipboard.ClearSelection()
    }
    return m, nil
```

**Benefits**:
- ✅ Auto-recovery from stuck state
- ✅ Better user experience
- ✅ Defensive programming
- ✅ No manual intervention needed

---

## 🚧 Blocker 4: State Explosion from Multiple Booleans

### Problem

```go
type AppModel struct {
    showHelp     bool
    finderActive bool
    searchActive bool  // Coming in Week 2
    // What if both are true? Which takes precedence?
}

// In View
if m.showHelp {
    return renderHelp()
}
if m.finderActive {
    return renderFinder()
}
if m.searchActive {
    return renderSearch()
}
// ⚠️ Ambiguous priority!
```

**Issue**: With N booleans, there are 2^N possible states. Most combinations don't make sense.

**Impact**:
- Confusing which state takes precedence
- Bugs from unexpected state combinations
- Hard to reason about transitions
- Code becomes if/else spaghetti

### Solution: Use State Machine

```go
// Explicit states (only one active at a time)
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

    // State-specific data
    finderInput    string
    searchQuery    string
    loadingMessage string
}

// State transitions are explicit
func (m AppModel) transitionTo(newState AppState) AppModel {
    // Clean up old state
    switch m.state {
    case StateFinder:
        m.finderInput = ""
    case StateSearch:
        m.searchQuery = ""
    }

    m.state = newState
    return m
}

// In Update
case tea.KeyMsg:
    switch msg.String() {
    case "?":
        return m.transitionTo(StateHelp), nil
    case "/":
        return m.transitionTo(StateFinder), loadMarkdownFilesCmd(m.rootPath)
    case "ctrl+f":
        return m.transitionTo(StateSearch), nil
    case "esc":
        return m.transitionTo(StateNormal), nil
    }

// In View (crystal clear)
switch m.state {
case StateNormal:
    return m.renderNormal()
case StateHelp:
    return m.renderHelp()
case StateFinder:
    return m.renderFinder()
case StateSearch:
    return m.renderSearch()
case StateLoading:
    return m.renderLoading()
}
```

**Benefits**:
- ✅ Only 1 state at a time (clear precedence)
- ✅ Explicit state transitions
- ✅ Impossible states are impossible
- ✅ Easy to add new states
- ✅ Clear in View which UI to show

---

## Impact on Week 2 (Ripgrep Search)

### Critical: Ripgrep Must Use Commands

Ripgrep is inherently asynchronous and streams results. We **must** use Commands, not callbacks.

**Wrong Approach** (will block UI):
```go
case "ctrl+f":
    results := ripgrep.Search(query)  // ⚠️ Blocks for seconds!
    m.searchResults = results
    return m, nil
```

**Right Approach** (non-blocking):
```go
// Message types
type SearchStartedMsg struct { query string }
type SearchResultMsg struct { result RipgrepResult }
type SearchCompletedMsg struct {}
type SearchErrorMsg struct { err error }

// Command that streams results
func searchCmd(query string, rootPath string) tea.Cmd {
    return func() tea.Msg {
        // Immediately return SearchStartedMsg
        go func() {
            executor := NewRipgrepExecutor()
            resultChan := make(chan RipgrepResult)

            go executor.Search(query, rootPath, resultChan)

            for result := range resultChan {
                // Send each result as a message
                program.Send(SearchResultMsg{result: result})
            }

            program.Send(SearchCompletedMsg{})
        }()

        return SearchStartedMsg{query: query}
    }
}

// In Update (handles streaming)
case SearchStartedMsg:
    m.state = StateSearch
    m.searchQuery = msg.query
    m.searchResults = []RipgrepResult{}  // Clear old results
    m.searchInProgress = true
    return m, nil

case SearchResultMsg:
    m.searchResults = append(m.searchResults, msg.result)
    return m, nil  // View re-renders with new result

case SearchCompletedMsg:
    m.searchInProgress = false
    return m, nil
```

### Recommended Week 2 Architecture

```go
type AppModel struct {
    // Use state machine
    state AppState  // Normal, Finder, Search, Help, Loading

    // Search state (flat, not in helper)
    searchQuery      string
    searchResults    []RipgrepResult
    searchCursor     int
    searchInProgress bool

    // No RipgrepExecutor instance stored!
    // Use it only in commands
}
```

---

## Action Plan

### Immediate (Week 1 Refactor)

1. ✅ **Document blockers** (this file)
2. **Optional**: Refactor finder to use Commands (can wait)

### Before Week 2 (Critical)

1. ✅ **Implement state machine** (StateNormal, StateFinder, StateSearch, StateHelp)
2. ✅ **Flatten finder state** (remove FuzzyFinderImpl dependency)
3. ✅ **Create search command pattern** (SearchStartedMsg, SearchResultMsg, etc.)

### Week 2 Implementation

1. ✅ **Use Commands for Ripgrep** (non-blocking, streaming)
2. ✅ **Stream results via messages** (not callbacks)
3. ✅ **Keep search state in AppModel** (not in RipgrepExecutor)

---

## Testing Strategy

### Test Pure Functions

```go
func TestFilterItems(t *testing.T) {
    items := []string{"apple", "banana", "cherry"}

    result := filterItems(items, "a")
    expected := []string{"apple", "banana"}

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Test State Transitions

```go
func TestStateTransitions(t *testing.T) {
    m := AppModel{state: StateNormal}

    m = m.transitionTo(StateFinder)
    if m.state != StateFinder {
        t.Errorf("Expected StateFinder, got %v", m.state)
    }

    // Verify old state cleaned up
    if m.finderInput != "" {
        t.Error("Finder input should be empty after transition")
    }
}
```

### Test Update Logic

```go
func TestFinderInputHandling(t *testing.T) {
    m := AppModel{
        state: StateFinder,
        finderInput: "",
    }

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
    newModel, _ := m.Update(msg)

    if newModel.(AppModel).finderInput != "t" {
        t.Error("Expected input to be 't'")
    }
}
```

---

## Conclusion

The Elm Architecture research revealed 4 blockers that could impact Week 2:

1. ✅ **Synchronous I/O** → Use Commands
2. ✅ **Split State** → Flatten into AppModel
3. ✅ **Mouse Drag State** → Add timeout recovery
4. ✅ **Boolean Flags** → Use state machine

**All blockers have clear solutions** and won't require major refactoring. The most critical change for Week 2 is using Commands for Ripgrep to enable streaming results without blocking the UI.

**Next Step**: Decide whether to refactor Week 1 implementation before starting Week 2, or proceed with Week 2 using the improved patterns.

---

**Status**: Blockers identified with solutions ✅
**Recommendation**: Proceed with Week 2 using improved patterns (Commands, state machine, flat state)
