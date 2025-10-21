# Phase 2 Implementation Plan
**Date**: 2025-10-21
**Version**: v1.1.0-alpha (Target)
**Duration**: 1-2 weeks (12-16 hours estimated)
**Status**: Ready for Implementation

---

## Overview

Phase 2 focuses on **user accessibility, search capabilities, and input flexibility**. The goal is to make LUMINA accessible to both terminal power users and casual users while maintaining vim keybinding support.

### Phase 2 Goals
- ✅ Enable mouse interaction (scroll)
- ✅ Improve navigation (jump, sort)
- ✅ Enhance search capabilities (global search)
- ✅ Increase accessibility (arrow keys, shift modifiers)
- ✅ Maintain backward compatibility (all existing keybindings work)

---

## Feature Implementation Order

### Recommended Priority (Based on Impact + Effort)

```
Week 1:
┌─────────────────────────────────────────────┐
│ Feature 4: Global Search (/)                │ ← CRITICAL (highest impact)
│ Estimated: 4-5 hours                        │
│ Dependencies: search.go (already exists)    │
│ User Impact: Very High                      │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Feature 2: SHIFT+Letter Jump                │ ← HIGH (very high impact)
│ Estimated: 2 hours                          │
│ Dependencies: keybindings.go                │
│ User Impact: High                           │
└─────────────────────────────────────────────┘

Week 2:
┌─────────────────────────────────────────────┐
│ Feature 3: Sort Toggle (S Key)              │ ← HIGH (improves usability)
│ Estimated: 2-3 hours                        │
│ Dependencies: model.go                      │
│ User Impact: High                           │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Feature 1: Mouse Scroll                     │ ← MEDIUM (nice to have)
│ Estimated: 2-3 hours                        │
│ Dependencies: main.go event handler         │
│ User Impact: Medium                         │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Feature 5: Scroll Improvements              │ ← MEDIUM (accessibility)
│ Estimated: 2-3 hours                        │
│ Dependencies: keybindings.go                │
│ User Impact: Medium                         │
└─────────────────────────────────────────────┘
```

---

## Feature 1: Global File Search (/)

### Complexity: HIGH | Impact: VERY HIGH | Effort: 4-5 hours

#### Current State
```go
// Currently in main.go lines 151-157:
case "filter":
    if m.currentView == FileTreeView {
        if m.fileList.FilterState() != list.Filtering {
            m.fileList.SetFilteringEnabled(true)  // ← Only filters current dir
        }
    }
```

#### Desired State
```go
// NEW: Distinguish between filter and search
case "filter":        // "/" → search current directory (existing)
case "global_search": // Different key or mode
```

#### Implementation Steps

**Step 1: Update keybindings.go**
```go
// Add new action for global search
type KeyBindings struct {
    // ... existing ...
    Actions: []KeyBinding{
        {Key: "/", Action: "filter", View: "filetree"},      // Current dir filter
        {Key: "shift+/", Action: "global_search", View: "filetree"},  // OR
        {Key: "ctrl+f", Action: "global_search", View: "any"},        // Alt approach
        // ... rest ...
    }
}
```

**Step 2: Create search state in model.go**
```go
type AppModel struct {
    // ... existing fields ...

    // NEW: Search state
    searchMode      bool      // Is search active?
    searchQuery     string    // Current search term
    searchResults   []string  // Matching file paths
    searchIndex     int       // Current result position
    showSearchPane  bool      // Show search results
}

// Helper function:
func (m *AppModel) SearchMarkdownFiles(query string) []string {
    var results []string

    // Search using existing findMarkdownFiles()
    allFiles := findMarkdownFiles(m.rootPath)

    // Filter by query
    for _, file := range allFiles {
        if strings.Contains(file, query) {
            results = append(results, file)
        }
    }

    return results
}
```

**Step 3: Add search handler in main.go**
```go
case "global_search":
    if m.currentView == FileTreeView {
        m.searchMode = true
        m.searchQuery = ""
        m.searchResults = []string{}
        // Show search prompt (similar to filter)
    }

// When in search mode, handle input differently
if m.searchMode {
    switch msg.String() {
    case "enter":
        // Execute search
        m.searchResults = m.SearchMarkdownFiles(m.searchQuery)
        m.searchIndex = 0
        m.showSearchPane = true
    case "esc":
        m.searchMode = false
        m.searchQuery = ""
    case "backspace":
        if len(m.searchQuery) > 0 {
            m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
        }
    default:
        m.searchQuery += msg.String()
    }
}
```

**Step 4: Display search results**
```go
// In View() function, show search results overlay
if m.showSearchPane {
    // Display results in center pane or overlay
    // Let user select result with j/k
    // Press Enter to open selected file
    // Press Esc to close
}
```

**Step 5: Update help text**
```go
// In help.go, add:
// "  /             Search all .md files (global)"
// "  ctrl+f        Same as above (alternative)"
```

#### Testing Checklist
- [ ] `/` opens search prompt
- [ ] Type query: `config`
- [ ] Press Enter: Shows matching files
- [ ] Navigate results with j/k
- [ ] Press Enter on result: Opens file
- [ ] Press Esc: Closes search
- [ ] Search across directory structure works
- [ ] Case-insensitive matching

#### Success Criteria
- ✅ User can search all markdown files
- ✅ Results display clearly
- ✅ Can open any result
- ✅ Escape closes search without losing current file
- ✅ Search is fast (<1 second for typical project)

---

## Feature 2: SHIFT+Letter Jump

### Complexity: MEDIUM | Impact: HIGH | Effort: 2 hours

#### Current State
No jump-to-letter functionality

#### Desired State
```
File Tree Contents:
- README.md
- config.json
- docs/
- helpers.go
- main.go

User presses: Shift+H
Result: Cursor jumps to "helpers.go"
```

#### Implementation Steps

**Step 1: Add binding detection in main.go**
```go
// In Update() handler for KeyMsg:
case tea.KeyMsg:
    key := msg.String()

    // NEW: Detect Shift+letter
    if strings.HasPrefix(key, "shift+") && len(key) == 7 {  // shift+X format
        letter := strings.ToLower(key[6:])
        if len(letter) == 1 && letter >= "a" && letter <= "z" {
            m.jumpToLetter(letter)
        }
    }
```

**Step 2: Implement jump logic**
```go
// In model.go, add:
func (m *AppModel) jumpToLetter(letter string) error {
    if m.currentView != FileTreeView {
        return nil
    }

    // Get current file list
    items := m.fileList.Items()
    if len(items) == 0 {
        return nil
    }

    // Get current cursor position
    currentIndex := m.fileList.Cursor()

    // Search from current position forward
    for i := currentIndex + 1; i < len(items); i++ {
        item := items[i].(FileItem)
        if strings.HasPrefix(strings.ToLower(item.name), letter) {
            m.fileList.SetCursor(i)
            return nil
        }
    }

    // Wrap around: search from beginning
    for i := 0; i <= currentIndex; i++ {
        item := items[i].(FileItem)
        if strings.HasPrefix(strings.ToLower(item.name), letter) {
            m.fileList.SetCursor(i)
            return nil
        }
    }

    return nil  // No match found
}
```

**Step 3: Update keybindings config**
```go
// In keybindings.go, add to Navigation section:
Navigation: []KeyBinding{
    {Key: "j", Action: "down", View: "filetree"},
    {Key: "k", Action: "up", View: "filetree"},
    {Key: "h", Action: "back", View: "filetree"},
    // NEW:
    {Key: "shift+a", Action: "jump_a", View: "filetree"},
    {Key: "shift+b", Action: "jump_b", View: "filetree"},
    // ... shift+c through shift+z ...
}

// OR use a generic approach:
// Just detect "shift+X" pattern and handle in main.go
```

**Step 4: Update help text**
```go
// In help.go, add:
// "  Shift+A-Z    Jump to first file starting with letter"
```

#### Testing Checklist
- [ ] Shift+A jumps to first 'A' file
- [ ] Shift+M jumps to first 'M' directory
- [ ] Case-insensitive (shift+a same as shift+A)
- [ ] Wrapping works (jump from end to beginning)
- [ ] No jump if no match
- [ ] Works with hidden files if shown
- [ ] Jump position maintained when opening file

#### Success Criteria
- ✅ Shift+letter jumps to first matching file/directory
- ✅ Wrapping at end of list works
- ✅ Case-insensitive matching
- ✅ Quick navigation in large file trees

---

## Feature 3: Sort Toggle (S Key)

### Complexity: MEDIUM | Impact: HIGH | Effort: 2-3 hours

#### Implementation Steps

**Step 1: Add sort state to AppModel**
```go
// In model.go:
type SortMode int

const (
    SortAlphabetical SortMode = iota
    SortRecent
)

type AppModel struct {
    // ... existing fields ...
    sortMode SortMode  // Current sort mode
}
```

**Step 2: Implement sort logic**
```go
// In model.go, add:
func (m *AppModel) sortFiles(items []list.Item) []list.Item {
    // Convert to FileItem slice
    files := make([]FileItem, len(items))
    for i, item := range items {
        files[i] = item.(FileItem)
    }

    switch m.sortMode {
    case SortAlphabetical:
        sort.Slice(files, func(i, j int) bool {
            // Directories first
            if files[i].isDir != files[j].isDir {
                return files[i].isDir
            }
            // Then alphabetical
            return files[i].name < files[j].name
        })

    case SortRecent:
        sort.Slice(files, func(i, j int) bool {
            // By modification time (recent first)
            iInfo, _ := os.Stat(files[i].path)
            jInfo, _ := os.Stat(files[j].path)
            iTime := iInfo.ModTime()
            jTime := jInfo.ModTime()
            return iTime.After(jTime)  // Recent first
        })
    }

    // Convert back to list.Item
    result := make([]list.Item, len(files))
    for i, f := range files {
        result[i] = f
    }
    return result
}
```

**Step 3: Add toggle handler**
```go
// In main.go Update() handler:
case "sort_toggle":
    if m.currentView == FileTreeView {
        // Toggle sort mode
        if m.sortMode == SortAlphabetical {
            m.sortMode = SortRecent
        } else {
            m.sortMode = SortAlphabetical
        }

        // Re-sort current file list
        items := m.fileList.Items()
        sorted := m.sortFiles(items)
        m.fileList.SetItems(sorted)
    }
```

**Step 4: Update keybindings**
```go
// In keybindings.go:
Actions: []KeyBinding{
    // ... existing ...
    {Key: "s", Action: "sort_toggle", View: "filetree"},
}
```

**Step 5: Update status bar**
```go
// In View() function, update status text:
viewName := []string{"FILE TREE", "VIEWER", "PREVIEW"}[m.currentView]

// NEW: Add sort indicator
var sortIndicator string
if m.currentView == FileTreeView {
    if m.sortMode == SortAlphabetical {
        sortIndicator = " [sort: ABC]"
    } else {
        sortIndicator = " [sort: RECENT]"
    }
}

statusText := fmt.Sprintf(
    "[%s]%s Tab: switch | j/k: nav | s: sort | ...",
    viewName, sortIndicator,
)
```

#### Testing Checklist
- [ ] Press 's' toggles sort mode
- [ ] Status bar shows current sort
- [ ] Alphabetical: directories first, then files
- [ ] Recent: recently modified items first
- [ ] Sort persists while navigating
- [ ] Sort applies when opening new directory

#### Success Criteria
- ✅ Sort toggle works with 's' key
- ✅ Visual indicator in status bar
- ✅ Alphabetical and recent modes work correctly
- ✅ Smooth transition between sort modes

---

## Feature 4: Mouse Scroll

### Complexity: LOW | Impact: MEDIUM | Effort: 2-3 hours

#### Implementation Steps

**Step 1: Add mouse event handler in main.go**
```go
// In Update() handler:
case tea.MouseMsg:
    // Handle mouse wheel scroll
    if msg.Button == tea.MouseWheelUp || msg.Button == tea.MouseWheelDown {
        // Determine which pane contains mouse position
        // Scroll appropriately
    }
```

**Step 2: Implement pane detection**
```go
func (m *AppModel) getPaneAtMouse(x, y int) ViewMode {
    // Check x coordinate against pane widths
    if x < m.fileTreeWidth {
        return FileTreeView
    } else if x < m.fileTreeWidth + m.viewerWidth {
        return ViewerView
    } else {
        return PreviewView
    }
}
```

**Step 3: Scroll in appropriate pane**
```go
case tea.MouseMsg:
    if msg.Button == tea.MouseWheelUp {
        pane := m.getPaneAtMouse(msg.X, msg.Y)

        switch pane {
        case ViewerView:
            m.viewer.LineUp(3)  // Scroll up 3 lines
        case FileTreeView:
            m.fileList.Update(tea.KeyMsg{...})  // Move up
        }
    } else if msg.Button == tea.MouseWheelDown {
        // Similar for down
    }
}
```

**Step 4: Update help text**
```go
// In help.go, add:
// "  Mouse wheel   Scroll active pane"
```

#### Testing Checklist
- [ ] Mouse wheel scrolls active pane
- [ ] Scroll up works
- [ ] Scroll down works
- [ ] Different scroll amounts per pane (lines vs items)
- [ ] Non-active panes don't scroll
- [ ] Works with Shift (faster scroll?)

#### Success Criteria
- ✅ Mouse wheel scrolls active pane
- ✅ Does not affect inactive panes
- ✅ Smooth scrolling experience

---

## Feature 5: Scroll Improvements

### Complexity: LOW | Impact: MEDIUM | Effort: 2-3 hours

#### Current Bindings
```
j/k       → line scroll ✅
d/u       → page scroll ✅
↑/↓       → line scroll ✅
Alt+↑/↓   → page scroll ✅
```

#### New Bindings to Add
```
Shift+↑   → page scroll up (alternative to u)
Shift+↓   → page scroll down (alternative to d)
```

#### Implementation Steps

**Step 1: Add Shift modifier detection**
```go
// In keybindings.go, add:
Scrolling: []KeyBinding{
    {Key: "shift+up", Action: "page_up", View: "viewer"},
    {Key: "shift+down", Action: "page_down", View: "viewer"},
    // ... keep existing bindings ...
}
```

**Step 2: No changes needed in main.go**
- The action-based dispatch already handles "page_up" and "page_down"
- Just need the keybinding to exist

**Step 3: Update help text**
```go
// In help.go, reorganize scroll help:
VIEWER (Scroll & Navigation)
  j, ↓         Scroll down (1 line)
  k, ↑         Scroll up (1 line)
  d            Page down (half page)
  u            Page up (half page)
  Shift+↓      Page down (alternative)
  Shift+↑      Page up (alternative)
  Alt+↓        Page down (power user)
  Alt+↑        Page up (power user)
  g            Go to top
  G            Go to bottom
```

**Step 4: Update status bar**
```go
// More comprehensive status message:
statusText = fmt.Sprintf(
    "[%s] ↑↓: scroll | Shift↑↓: page | j/k: line | d/u: page | g/G: top/bottom | y: copy | ?: help | q: quit",
    viewName,
)
```

#### Testing Checklist
- [ ] Shift+↑ pages up
- [ ] Shift+↓ pages down
- [ ] All existing scroll works
- [ ] Help text shows all options
- [ ] Non-terminal users can navigate with arrows alone
- [ ] Power users can use vim keys
- [ ] Both work simultaneously

#### Success Criteria
- ✅ Multiple scroll options available
- ✅ Accessible to non-terminal users
- ✅ Power user options still work
- ✅ Clear documentation

---

## Implementation Schedule

### Week 1 - Search & Navigation
```
Day 1-2: Global File Search (/)
- Monday: Implementation + basic testing
- Tuesday: Integration + help text

Day 3: SHIFT+Letter Jump
- Wednesday: Implementation + testing
```

### Week 2 - UX Improvements
```
Day 1: Sort Toggle (S Key)
- Thursday: Implementation + testing

Day 2: Mouse Scroll + Scroll Improvements
- Friday: Both features + comprehensive testing
```

### Buffer Time
- Weekend: Testing, documentation, bugfixes

---

## Testing Strategy

### Unit Tests to Add
```go
// search_test.go
TestSearchMarkdownFiles()
TestSearchCaseSensitivity()
TestSearchResults()

// jump_test.go
TestJumpToLetter()
TestJumpWrapping()
TestJumpNearEnd()

// sort_test.go
TestSortAlphabetical()
TestSortRecent()
TestSortToggle()

// scroll_test.go
TestMouseScrollDetection()
TestScrollBoundaries()
```

### Integration Tests
```bash
# Start app, run through feature sequence
./test-phase-2.sh

# Test each feature in combination with others
# Verify git hooks still pass
# Build successful on all changes
```

### Manual Testing Checklist
See PHASE_2_FEATURE_REQUESTS.md for detailed manual testing procedures

---

## Documentation Updates

### Files to Update
- `README.md` - Add Phase 2 features section
- `help.go` - Update all help text
- `CONFIG.md` - Document new keybindings
- `CHANGELOG.md` - Record Phase 2 changes
- `docs/ACCESSIBILITY.md` - NEW - non-terminal user guide

### New Documentation
- `PHASE_2_COMPLETE.md` - Completion summary
- `docs/SEARCH_GUIDE.md` - Global search user guide
- `docs/NAVIGATION_GUIDE.md` - Jump and sort guide

---

## Version Management

### Current Version
```
v1.0.1-alpha (Phase 1.5)
```

### Phase 2 Version
```
v1.1.0-alpha (Phase 2)
```

### Version Increment
```go
// Update version.go when Phase 2 complete:
Version = "1.1.0-alpha"
BuildPhase = "Phase 2"
BuildDate = "2025-11-04"  // Expected completion
```

---

## Rollback Points

### Before Phase 2 Implementation
```bash
git tag v1.0.1-alpha-final

# If issues occur:
git checkout v1.0.1-alpha-final
go build -o lumina
```

### After Each Feature
```bash
# After Global Search:
git tag v1.1.0-alpha-search

# After SHIFT+Letter Jump:
git tag v1.1.0-alpha-jump

# ... etc
```

---

## Success Criteria for Phase 2

✅ **All 5 Features Implemented**
- [ ] Global search works across all files
- [ ] SHIFT+letter jumps to files
- [ ] 's' toggles between sort modes
- [ ] Mouse wheel scrolls active pane
- [ ] Shift+arrows page scroll

✅ **All Tests Pass**
- [ ] Unit tests: 60%+ coverage target
- [ ] Integration tests: All features work together
- [ ] Manual testing: All procedures pass
- [ ] Build: No errors or warnings

✅ **Documentation Complete**
- [ ] Help text updated
- [ ] User guides written
- [ ] CHANGELOG updated
- [ ] Code comments clear

✅ **Production Ready**
- [ ] Git hooks active
- [ ] All changes committed
- [ ] Tagged with v1.1.0-alpha
- [ ] Ready for user feedback

---

## Next Steps (After Phase 2)

### Phase 3 (Future)
- [ ] Editor integration ($EDITOR)
- [ ] Claude Code integration (/workflows, /moe)
- [ ] Agentic features
- [ ] Git awareness

### Long-term Enhancements
- [ ] Fuzzy file search
- [ ] Content search (grep)
- [ ] Bookmarks / favorites
- [ ] Custom themes
- [ ] Plugin system

---

## Risk Assessment

### Low Risk
- ✅ Scroll improvements (minimal changes)
- ✅ Mouse scroll (isolated feature)
- ✅ SHIFT+letter jump (local to file tree)

### Medium Risk
- ⚠️ Sort toggle (affects file list state)
- ⚠️ Global search (new code path)

### Mitigation
- Comprehensive testing
- Git rollback capability
- Feature flags (if needed)

---

## Conclusion

**Phase 2 is designed to significantly improve LUMINA's usability and accessibility while maintaining all existing functionality.**

Target:
- ✅ Implement all 5 features
- ✅ Maintain 100% backward compatibility
- ✅ Add comprehensive tests
- ✅ Complete in 1-2 weeks

Result:
- ✅ v1.1.0-alpha with search, navigation, and accessibility improvements
- ✅ Ready for broader user testing
- ✅ Foundation for Phase 3 features

---

**Plan Created**: 2025-10-21
**Target Completion**: 2025-11-04
**Status**: Ready for Implementation
