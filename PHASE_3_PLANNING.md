# Phase 3: UI Component Integration - Planning Document

**Date**: November 1, 2025
**Status**: Planning
**Target**: Integrate Phase 2 backend modules into Bubble Tea UI

---

## 🎯 Overview

Phase 3 integrates the fully-tested Phase 2 backend modules (fuzzy finder, ripgrep search, file watcher) into the existing Bubble Tea UI architecture, creating a cohesive, production-ready TUI experience.

---

## 📋 Current Architecture Analysis

### Existing Bubble Tea Model (`main.go`)

**Current Structure**:
```go
type AppModel struct {
    // Core state
    currentPath    string
    selectedFile   string
    currentView    ViewType  // FileTree, Viewer, Preview

    // UI components
    fileList       list.Model      // Bubbles list for file tree
    viewer         viewport.Model  // Bubbles viewport for content

    // Dimensions
    width, height  int
    fileTreeWidth  int
    viewerWidth    int
    previewWidth   int

    // Features
    clipboard      ClipboardManager
    colorManager   ColorManager
    keyBindings    KeyBindings
    showHelp       bool

    // Mouse selection
    mouseDragActive     bool
    mouseDragStartLine  int
    mouseDragStartCol   int
    selectionMode       SelectionMode
}
```

**Current Update() Flow**:
1. Handle window resize messages
2. Handle mouse events (selection, scrolling)
3. Handle keyboard events via keybindings
4. Update active component (fileList or viewer)
5. Return updated model

**Current View() Flow**:
1. Render 3-pane layout (File Tree | Viewer | Preview)
2. Apply active pane styling (pink border)
3. Render status bar with context-aware hints
4. Optionally overlay help modal

---

## 🎯 Phase 3 Objectives

### Primary Goals

1. **Fuzzy Finder Integration**
   - Add fuzzy finder modal overlay
   - Trigger with `/` key from any pane
   - Real-time filtering as user types
   - Navigate with j/k, select with Enter, cancel with Esc

2. **Ripgrep Search Integration**
   - Add search results pane (split viewer or separate modal)
   - Trigger with `Ctrl+F` or custom key
   - Display results with file:line:snippet format
   - Jump to result location in viewer

3. **File Watcher Integration**
   - Auto-reload currently open file when changed
   - Visual indicator for file changes (status bar notification)
   - Debounced updates to prevent UI flicker
   - Preserve scroll position on reload

4. **UI State Management**
   - Add new state modes: `FinderMode`, `SearchMode`, `NormalMode`
   - Clean state transitions between modes
   - Proper cleanup when switching modes

---

## 🏗️ Implementation Architecture

### New Model Fields

```go
type AppModel struct {
    // ... existing fields ...

    // Phase 3: Fuzzy Finder
    finderActive    bool
    finderInput     string
    finderResults   []FileItem
    finderSelected  int
    fuzzyFinder     *FuzzyFinder  // from fuzzy_finder_impl.go

    // Phase 3: Ripgrep Search
    searchActive    bool
    searchInput     string
    searchResults   []RipgrepResult
    searchSelected  int
    ripgrepExecutor *RipgrepExecutor  // from ripgrep_executor.go

    // Phase 3: File Watcher
    fileWatcher     *FileWatcher  // from file_watcher.go
    watcherActive   bool
    fileChangedNotification bool

    // Phase 3: Mode management
    currentMode     UIMode  // Normal, Finder, Search
}

type UIMode int
const (
    NormalMode UIMode = iota
    FinderMode
    SearchMode
)
```

### Message Types (Bubble Tea)

```go
// Fuzzy Finder messages
type FinderOpenMsg struct{}
type FinderCloseMsg struct{}
type FinderInputMsg string
type FinderSelectMsg int

// Ripgrep Search messages
type SearchStartMsg struct{ Query string }
type SearchResultMsg struct{ Result RipgrepResult }
type SearchCompleteMsg struct{}
type SearchSelectMsg int

// File Watcher messages
type FileChangedMsg struct{ Path string }
type FileReloadMsg struct{ Content string }
```

---

## 📐 UI Layout Design

### Mode 1: Normal Mode (Current)
```
┌─────────────────────────────────────────────────────────────┐
│ Claude Code Navigator (CCN) - /path/to/project             │
├──────────────┬───────────────────────────┬──────────────────┤
│ File Tree    │ Viewer                    │ Preview          │
│ (20%)        │ (60%)                     │ (20%)            │
│              │                           │                  │
│ ✓ README.md  │ # README                  │ Coming soon      │
│   src/       │                           │                  │
│   docs/      │ This is a project...      │                  │
│              │                           │                  │
├──────────────┴───────────────────────────┴──────────────────┤
│ [FILE TREE] Tab: switch | j/k: nav | /: finder | Ctrl+F...  │
└─────────────────────────────────────────────────────────────┘
```

### Mode 2: Finder Mode (New)
```
┌─────────────────────────────────────────────────────────────┐
│ Claude Code Navigator (CCN) - /path/to/project             │
├──────────────┬───────────────────────────┬──────────────────┤
│              │                           │                  │
│              │  ┌─────────────────────┐  │                  │
│              │  │ Fuzzy Finder        │  │                  │
│ (Dimmed)     │  ├─────────────────────┤  │ (Dimmed)         │
│              │  │ Query: src          │  │                  │
│              │  ├─────────────────────┤  │                  │
│              │  │ ▸ src/main.go       │  │                  │
│              │  │   src/utils.go      │  │                  │
│              │  │   src/model.go      │  │                  │
│              │  │ 42/1000 matches     │  │                  │
│              │  └─────────────────────┘  │                  │
├──────────────┴───────────────────────────┴──────────────────┤
│ [FINDER] j/k: navigate | Enter: select | Esc: cancel        │
└─────────────────────────────────────────────────────────────┘
```

### Mode 3: Search Mode (New)
```
┌─────────────────────────────────────────────────────────────┐
│ Claude Code Navigator (CCN) - /path/to/project             │
├──────────────┬───────────────────────────┬──────────────────┤
│ File Tree    │ Search Results            │ Viewer           │
│ (20%)        │ (40%)                     │ (40%)            │
│              │                           │                  │
│ ✓ README.md  │ Query: "TODO"             │ # README         │
│   src/       │ ─────────────────         │                  │
│   docs/      │ ▸ main.go:42              │ Line 42:         │
│              │   // TODO: fix bug        │ // TODO: fix bug │
│              │   src/utils.go:18         │                  │
│              │   # TODO: refactor        │ (Preview)        │
│              │                           │                  │
├──────────────┴───────────────────────────┴──────────────────┤
│ [SEARCH] j/k: navigate | Enter: jump | n/N: next/prev | Esc │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 State Transition Diagram

```
                  ┌──────────────┐
                  │ Normal Mode  │
                  │              │
                  │ - File Tree  │
                  │ - Viewer     │
                  │ - Preview    │
                  └──────┬───────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
       / key          Ctrl+F        File Changed
         │               │               │
         ▼               ▼               ▼
  ┌──────────┐    ┌──────────┐   ┌──────────────┐
  │  Finder  │    │  Search  │   │  Reload      │
  │  Mode    │    │  Mode    │   │  Notification│
  └─────┬────┘    └─────┬────┘   └──────────────┘
        │               │
    Enter/Esc       Enter/Esc
        │               │
        └───────┬───────┘
                │
                ▼
         ┌──────────────┐
         │ Normal Mode  │
         │ (Updated)    │
         └──────────────┘
```

---

## 🛠️ Implementation Plan

### Week 1: Fuzzy Finder Integration

**Day 1-2: State Management**
- Add finder-related fields to AppModel
- Implement FinderMode state
- Create message types for finder events
- Write state transition logic

**Day 3-4: UI Integration**
- Design finder modal overlay with Lip Gloss
- Integrate FuzzyFinder backend (fuzzy_finder_impl.go)
- Handle input events (typing, filtering)
- Implement result navigation (j/k keys)

**Day 5: Testing & Polish**
- Test fuzzy matching performance (1000+ files)
- Test keybinding conflicts
- Polish modal styling
- Write integration tests

**Deliverables**:
- ✅ Fuzzy finder modal renders correctly
- ✅ Real-time filtering works (<50ms)
- ✅ Navigation and selection work smoothly
- ✅ Clean state transitions (Normal ↔ Finder)

---

### Week 2: Ripgrep Search Integration

**Day 1-2: Search State**
- Add search-related fields to AppModel
- Implement SearchMode state
- Create ripgrep message types
- Handle async search results (channels)

**Day 3-4: UI Integration**
- Design search results pane layout
- Integrate RipgrepExecutor backend
- Stream results to UI as they arrive
- Implement result navigation

**Day 5: Jump-to-Location**
- Implement jump to file:line in viewer
- Highlight matched text in viewer
- Scroll to match location
- Handle n/N for next/previous match

**Deliverables**:
- ✅ Search UI with streaming results
- ✅ Jump to match location works
- ✅ Match highlighting in viewer
- ✅ Handle large result sets (1000+ matches)

---

### Week 3: File Watcher Integration

**Day 1-2: Watcher Setup**
- Integrate FileWatcher backend
- Add watcher state to AppModel
- Create file change message types
- Handle watcher lifecycle (start/stop)

**Day 3-4: Auto-reload**
- Implement file reload on change detection
- Preserve scroll position during reload
- Add debouncing (500ms)
- Visual notification in status bar

**Day 5: Testing & Edge Cases**
- Test rapid file changes
- Test large file reloads
- Test watcher with symlinks
- Handle errors gracefully

**Deliverables**:
- ✅ Auto-reload on file change
- ✅ Scroll position preserved
- ✅ Visual change indicator
- ✅ No UI flicker or lag

---

### Week 4: Polish & Integration

**Day 1-2: UI Refinements**
- Consistent styling across modes
- Smooth transitions and animations
- Error message handling
- Loading indicators

**Day 3-4: Testing**
- End-to-end integration tests
- Performance testing (large codebases)
- Memory leak testing
- Cross-mode interaction testing

**Day 5: Documentation**
- Update README with new features
- Create user guide for finder/search
- Document keybindings
- Update PROJECT-STATUS.md

**Deliverables**:
- ✅ All modes work seamlessly together
- ✅ Comprehensive test coverage
- ✅ Updated documentation
- ✅ Phase 3 complete!

---

## 📊 Success Metrics

### Performance Targets

| Metric | Target | Test Method |
|--------|--------|-------------|
| Fuzzy finder latency | <50ms | 1000 files, real-time typing |
| Search result streaming | <100ms first result | ripgrep 10,000 files |
| File reload latency | <200ms | 1MB markdown file |
| Memory usage | <50MB | All features active |
| UI responsiveness | 60 FPS | Visual smoothness check |

### Functional Requirements

**Fuzzy Finder**:
- ✅ Opens with `/` key
- ✅ Filters in real-time
- ✅ Shows match score ranking
- ✅ Handles 10,000+ files
- ✅ Esc cancels without side effects

**Ripgrep Search**:
- ✅ Opens with `Ctrl+F` or custom key
- ✅ Streams results as they arrive
- ✅ Jump to match preserves context
- ✅ Highlights matched text
- ✅ Navigate with n/N

**File Watcher**:
- ✅ Detects file changes within 500ms
- ✅ Auto-reloads current file
- ✅ Visual notification in status bar
- ✅ Scroll position preserved
- ✅ No watcher for non-markdown files

---

## 🔍 Technical Challenges & Solutions

### Challenge 1: Async Message Handling

**Problem**: Ripgrep results arrive asynchronously via channels. How to integrate with Bubble Tea's synchronous Update() model?

**Solution**: Use Bubble Tea's `Cmd` pattern:
```go
func executeSearch(query string) tea.Cmd {
    return func() tea.Msg {
        resultChan, _ := ripgrep.Execute(ctx, query)
        for result := range resultChan {
            // Send each result as a message
            return SearchResultMsg{Result: result}
        }
        return SearchCompleteMsg{}
    }
}
```

### Challenge 2: Modal Overlay Rendering

**Problem**: How to render finder modal on top of existing UI without breaking layout?

**Solution**: Use `lipgloss.Place()` with layering:
```go
if m.finderActive {
    finderOverlay := renderFinderModal(m)
    baseView := renderNormalView(m)

    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        finderOverlay,
        lipgloss.WithWhitespaceForeground(baseView),
    )
}
```

### Challenge 3: State Transition Cleanup

**Problem**: Switching between modes might leave stale state (e.g., active searches, watcher goroutines).

**Solution**: Explicit cleanup on mode exit:
```go
func (m AppModel) exitFinderMode() AppModel {
    m.finderActive = false
    m.finderInput = ""
    m.finderResults = nil
    m.finderSelected = 0
    return m
}
```

### Challenge 4: Scroll Position Preservation

**Problem**: File reload resets viewer scroll to top.

**Solution**: Save and restore viewport offset:
```go
func (m AppModel) reloadFile(newContent string) AppModel {
    oldOffset := m.viewer.YOffset
    m.viewer.SetContent(newContent)
    m.viewer.YOffset = oldOffset  // Restore position
    return m
}
```

---

## 🧪 Testing Strategy

### Unit Tests (Already Complete from Phase 2)

- ✅ FuzzyFinder filtering and scoring
- ✅ RipgrepExecutor execution and parsing
- ✅ FileWatcher change detection and debouncing

### Integration Tests (New for Phase 3)

**Test Suite 1: Fuzzy Finder UI**
```go
func TestFinderUIIntegration(t *testing.T) {
    // Test: Open finder with /
    // Test: Type query updates results
    // Test: j/k navigation
    // Test: Enter selects file
    // Test: Esc cancels
}
```

**Test Suite 2: Search UI**
```go
func TestSearchUIIntegration(t *testing.T) {
    // Test: Open search with Ctrl+F
    // Test: Results stream to UI
    // Test: Jump to match location
    // Test: n/N navigation
    // Test: Highlight in viewer
}
```

**Test Suite 3: File Watcher UI**
```go
func TestFileWatcherUIIntegration(t *testing.T) {
    // Test: File change detected
    // Test: Auto-reload triggered
    // Test: Scroll position preserved
    // Test: Notification displayed
    // Test: Debouncing works
}
```

### Manual Testing Checklist

**Fuzzy Finder**:
- [ ] Open with `/` from file tree
- [ ] Open with `/` from viewer
- [ ] Real-time filtering feels instant
- [ ] Match highlighting is clear
- [ ] Works with 1000+ files
- [ ] Esc returns to previous state

**Search**:
- [ ] Open with `Ctrl+F`
- [ ] Results stream smoothly
- [ ] Jump to match is accurate
- [ ] Highlight is visible
- [ ] n/N navigation cycles correctly
- [ ] Works with 10,000+ matches

**File Watcher**:
- [ ] Change file externally (vim, vscode)
- [ ] Reload happens within 1 second
- [ ] Scroll position maintained
- [ ] Notification appears in status bar
- [ ] No flicker or lag
- [ ] Watcher stops when file closed

---

## 📚 Documentation Updates

### README.md Updates

Add new sections:
- **Fuzzy Finder**: How to use `/` key for quick file navigation
- **Content Search**: How to use `Ctrl+F` for ripgrep search
- **Auto-reload**: Explain file watcher behavior

### Keybinding Reference

Update keybinding table:
```
Normal Mode:
  /         Open fuzzy finder
  Ctrl+F    Open content search

Finder Mode:
  j/k       Navigate results
  Enter     Select file
  Esc       Cancel

Search Mode:
  j/k       Navigate results
  Enter     Jump to match
  n/N       Next/previous match
  Esc       Close search
```

### Architecture Documentation

Create `ARCHITECTURE.md`:
- Bubble Tea MVU pattern explanation
- State management strategy
- Message flow diagrams
- Extension points for future features

---

## 🎯 Phase 3 Milestones

### Milestone 1: Fuzzy Finder Complete (Week 1)
- ✅ Modal UI rendered
- ✅ Real-time filtering works
- ✅ Navigation and selection functional
- ✅ Integration tests passing

### Milestone 2: Search Integration Complete (Week 2)
- ✅ Search UI with results pane
- ✅ Streaming results working
- ✅ Jump-to-location functional
- ✅ Match highlighting working

### Milestone 3: File Watcher Complete (Week 3)
- ✅ Auto-reload on file change
- ✅ Scroll position preserved
- ✅ Visual notifications working
- ✅ Watcher lifecycle managed

### Milestone 4: Phase 3 Complete (Week 4)
- ✅ All features integrated
- ✅ Tests passing (unit + integration)
- ✅ Documentation updated
- ✅ Performance targets met
- ✅ Ready for Phase 4 (Claude Code integration)

---

## 🔮 Future Considerations (Phase 4+)

### Phase 4: Claude Code Integration
- Terminal pane split for Claude interaction
- `/moe` and `/workflows` shortcuts
- Git status awareness in file tree
- MCP server detection and integration

### Phase 5: Advanced Features
- Multiple file tabs
- Split pane support (horizontal/vertical)
- Session persistence (save/restore state)
- Bookmark system (vim marks)
- Custom themes and configuration

---

## 📊 Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Performance degradation with large repos | HIGH | Implement lazy loading, result limiting |
| State management complexity | MEDIUM | Clear state machine, documented transitions |
| Memory leaks from goroutines | HIGH | Strict lifecycle management, cleanup on exit |
| Keybinding conflicts | LOW | Configurable keybindings, mode isolation |
| Async message ordering issues | MEDIUM | Use message queues, proper synchronization |

---

## ✅ Definition of Done

Phase 3 is complete when:

1. **All features working**:
   - ✅ Fuzzy finder opens, filters, and selects files
   - ✅ Ripgrep search finds matches and jumps to location
   - ✅ File watcher auto-reloads changed files

2. **Tests passing**:
   - ✅ 98 UI integration tests implemented and passing
   - ✅ End-to-end workflow tests passing
   - ✅ Performance benchmarks met

3. **Documentation complete**:
   - ✅ README updated with new features
   - ✅ Keybinding reference updated
   - ✅ Architecture documentation written

4. **Code quality**:
   - ✅ No goroutine leaks
   - ✅ No memory leaks
   - ✅ Clean state transitions
   - ✅ Error handling comprehensive

5. **User experience**:
   - ✅ Smooth transitions between modes
   - ✅ No UI flicker or lag
   - ✅ Intuitive keybindings
   - ✅ Helpful error messages

---

**Next Steps**: Begin Week 1 implementation - Fuzzy Finder Integration

**Estimated Completion**: 4 weeks from start date
