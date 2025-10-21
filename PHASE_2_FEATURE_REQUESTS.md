# Phase 2 Feature Requests
**Date**: 2025-10-21
**Priority**: HIGH - Core UX Improvements
**Status**: Documented & Ready for Implementation

---

## Feature 1: Mouse Scroll in Panels ⭐ HIGH PRIORITY

### User Request
"Mouse scroll option in each both panels depending on current panel selected"

### What It Should Do
- Enable mouse wheel scrolling in the viewer pane when it's active (currently selected)
- Enable mouse wheel scrolling in the file tree when it's active
- Disable mouse interactions in inactive panes
- Provide smooth, intuitive scrolling behavior

### Implementation Notes
- Bubble Tea supports mouse events
- Check `tea.MouseMsg` in event handler
- Need to track which pane is active: `m.currentView`
- Adjust scroll amount: `LineDown()`, `HalfViewDown()`, `ViewDown()`
- Consider scroll sensitivity (lines vs pages)

### Expected Behavior
```
Active Pane (Viewer):
- Scroll up: m.viewer.LineUp(3)    [3-line scroll]
- Scroll down: m.viewer.LineDown(3) [3-line scroll]
- Hold shift: m.viewer.HalfViewUp/Down() [page scroll]

Active Pane (File Tree):
- Scroll up/down: Navigate up/down in list
```

### Files to Modify
- `main.go` - Add mouse event handler
- `model.go` - Track mouse state if needed

### Estimated Effort: 2-3 hours

---

## Feature 2: SHIFT + Letter Jump ⭐ VERY HIGH PRIORITY

### User Request
"SHIFT + <letter> when in file directory to automatically jump to a certain letter in the list"

### What It Should Do
- When in file tree pane, press SHIFT+A to jump to first file starting with 'A'
- Press SHIFT+D to jump to first directory starting with 'D'
- Case-insensitive matching
- Wrap around at end of list (if no match below, search from top)
- Similar to Vim's `f`/`F` commands or file browser behavior

### Implementation Notes
- Listen for `shift+a` through `shift+z` key combinations
- Get current file list from `m.fileList`
- Search for first item matching the letter
- Move selection with `m.fileList.SetCursor(index)`
- Only works in file tree pane

### Expected Behavior
```
File tree current state:
- README.md
- config.json
- docs/
- helpers.go
- main.go

User presses: SHIFT+H
Result: Cursor jumps to "helpers.go"

User presses: SHIFT+D
Result: Cursor jumps to "docs/"
```

### Files to Modify
- `main.go` - Add shift+letter detection and jump logic
- `model.go` - Add jump-to-letter function

### Estimated Effort: 2 hours

---

## Feature 3: Sort Toggle (S Key) ⭐ HIGH PRIORITY

### User Request
"S key command binding to toggle between alphabetical and recent sort. Useful for recently created/modified files."

### What It Should Do
- Press 's' (or 'S'?) to toggle sort mode
- Two sort modes:
  1. **Alphabetical** (default) - A-Z, then directories
  2. **Recent** - Most recently modified first
- Visual indicator showing current sort mode in status bar
- Remember sort preference during session

### Sort Algorithm

**Alphabetical**:
```
Directories first (sorted A-Z)
Then files (sorted A-Z)
```

**Recent** (by modification time):
```
Most recently modified items first
Mix of files and directories
```

### Implementation Notes
- Add `sortMode` to `AppModel` struct
- Values: `const SortAlphabetical, SortRecent`
- Modify `loadDirectory()` to sort based on mode
- When 's' pressed: toggle mode and reload file list
- Update status bar to show: `[Viewer] sort: ABC` or `[Viewer] sort: RECENT`

### Expected Behavior
```
File Tree (Alphabetical):
- docs/
- helpers.go
- main.go
- README.md

After pressing 's':

File Tree (Recent):
- helpers.go (modified 5 min ago)
- config.json (modified 1 hour ago)
- docs/ (modified 2 hours ago)
- main.go (modified yesterday)
```

### Files to Modify
- `main.go` - Add 's' key binding handler
- `model.go` - Add `sortMode`, `SortType` enum, sort logic
- `help.go` - Document 's' key in help text

### Estimated Effort: 2-3 hours

---

## Feature 4: Global File Search (/) ⭐ VERY HIGH PRIORITY

### User Request
"Search command doesn't work with / - only applies filter. Want search across global files (at least .md for sure), similar to what glow does."

### Current Behavior
- `/` in file tree enters filter mode (searches current directory)
- Filters by filename only
- Doesn't search across all directories

### Desired Behavior
- `/` opens a search prompt (like current filter)
- Search across ALL markdown files in the project
- Display search results in a results pane
- Navigate through results
- Can open a result and jump to that file

### How Glow Does It
- `glow` has a search mode that:
  1. Scans all `.md` files recursively
  2. Shows matching results in a list
  3. Opens selected result

### Implementation Notes
- Use existing `findMarkdownFiles()` from model.go
- Search across filenames AND potentially file content (future enhancement)
- For now: filename search across all files

**Search UX**:
```
User presses: /
Shows: Search: _________________ (enter search term)

User types: "config"
Shows results:
  - docs/CONFIG.md
  - docs/CONFIGURATION_MANAGEMENT.md
  - .claude/settings.json (config file)

User presses: Enter or arrows to select, then Enter to open
```

### Two Implementation Approaches

**Approach A (Simpler - Phase 2)**:
- Filename search only
- Display results as overlay
- Uses existing `preview` pane or new search pane

**Approach B (Full - Phase 2.5)**:
- Filename search
- Content search (grep-like)
- Fuzzy matching
- Index for performance

**Recommendation**: Start with Approach A

### Files to Modify/Create
- `main.go` - Add search handler, differentiate `/` filter vs search
- `model.go` - Add search functionality, search state
- `search.go` - Already exists! Enhance it with real search logic
- `help.go` - Document search functionality

### Estimated Effort: 4-5 hours (Approach A)

---

## Feature 5: Scroll Improvements ⭐ HIGH PRIORITY

### User Request
Multiple improvements to make scrolling more accessible:

1. **"Scroll works with d/u rather than d/e"**
   - Currently: d=page_down, u=page_up (existing)
   - User confirms: Keep this! This is working.

2. **"j/k is good"**
   - Currently: j=scroll_down, k=scroll_up (one line at a time)
   - User confirms: Keep this! This is working.

3. **"Still like ability to use arrows"**
   - Currently: ↑ and ↓ work for scroll
   - Confirm: Keep working

4. **"SHIFT + up/down allows page movement"**
   - NEW: Shift+↑ for page up, Shift+↓ for page down
   - Alternative to 'u' and 'd' keys
   - More accessible for non-terminal users

5. **"Make it accessible to non-terminal junkies"**
   - Combine multiple scroll options
   - Don't force vim keys
   - Arrow keys + shift modifiers
   - Mouse scrolling (Feature #1)

### Current Scroll Bindings (Phase 1.5)
```
Viewer Pane:
j / ↓         → scroll down 1 line
k / ↑         → scroll up 1 line
d             → page down (half page)
u             → page up (half page)
g             → go to top
G             → go to bottom
Alt+↓        → page down (Option+Down on macOS)
Alt+↑        → page up (Option+Up on macOS)
```

### Proposed New Bindings (Phase 2)
```
Keep All Existing:
j / ↓         → scroll down 1 line
k / ↑         → scroll up 1 line
d             → page down (half page)
u             → page up (half page)
g             → go to top
G             → go to bottom

Add For Accessibility:
Shift+↓       → page down (like Alt+↓ but more discoverable)
Shift+↑       → page up (like Alt+↑ but more discoverable)

Keep Optional:
Alt+↓        → page down (power user option)
Alt+↑        → page up (power user option)
Mouse wheel   → scroll (see Feature #1)
```

### Implementation Notes
- Update `keybindings.go` default bindings
- Add Shift+ modifier handling in main.go key handler
- Update help text to show multiple scroll options
- Update status bar to show: `[VIEWER] ↑↓: scroll | Shift↑↓: page | j/k: line | d/u: page`

### Files to Modify
- `keybindings.go` - Add shift+up, shift+down bindings
- `main.go` - Handle shift modifier
- `help.go` - Update help text with accessibility info
- `version.go` - Increment to v1.1.0-alpha when complete

### Expected Behavior
```
New User (Non-Terminal):
- Presses: Arrow keys (↑↓) → scrolls line by line ✅
- Presses: Shift+Arrow (↑↓) → pages ✅
- Uses: Mouse wheel → scrolls ✅

Power User (Terminal):
- Uses: j/k for line scroll ✅
- Uses: d/u for page scroll ✅
- Uses: Arrow keys as backup ✅
- Uses: g/G for jump ✅
```

### Estimated Effort: 2-3 hours

---

## Priority Matrix

| Feature | Priority | Effort | Impact | Phase |
|---------|----------|--------|--------|-------|
| Mouse Scroll | HIGH | 2-3h | Medium | 2.0 |
| SHIFT+Letter Jump | VERY HIGH | 2h | High | 2.0 |
| Sort Toggle (S) | HIGH | 2-3h | High | 2.0 |
| Global Search (/) | VERY HIGH | 4-5h | Very High | 2.0 |
| Scroll Improvements | HIGH | 2-3h | High | 2.0 |

---

## Implementation Timeline

### Phase 2.0 (Priority Order)
1. **Global Search (/)** - 4-5h (highest impact)
2. **SHIFT+Letter Jump** - 2h (very high impact)
3. **Sort Toggle (S)** - 2-3h
4. **Mouse Scroll** - 2-3h
5. **Scroll Improvements** - 2-3h

**Total Estimated Time**: 12-16 hours
**Recommended Duration**: 1-2 weeks

### Suggested Breakdown
- **Week 1**: Global search + SHIFT+letter jump
- **Week 2**: Sort toggle + mouse scroll + scroll improvements + testing

---

## Implementation Checklist

### Feature: Global Search
- [ ] Differentiate `/` filter (current dir) vs search (all dirs)
- [ ] Implement search across all `.md` files
- [ ] Create search results view
- [ ] Handle search UX and result navigation
- [ ] Update help text
- [ ] Add search to keybindings config
- [ ] Test with multiple file structures

### Feature: SHIFT+Letter Jump
- [ ] Detect SHIFT+letter key combinations
- [ ] Implement jump logic with wrapping
- [ ] Update file list cursor position
- [ ] Test with various file lists
- [ ] Update keybindings config

### Feature: Sort Toggle
- [ ] Add sort mode to AppModel
- [ ] Implement sorting logic (alphabetical, recent)
- [ ] Handle 's' key binding
- [ ] Update file list when sort changes
- [ ] Show sort indicator in status bar
- [ ] Update help text

### Feature: Mouse Scroll
- [ ] Add mouse event handler
- [ ] Detect active pane
- [ ] Implement scroll in viewer
- [ ] Implement scroll in file tree
- [ ] Handle shift for different scroll amounts
- [ ] Test with different terminal emulators

### Feature: Scroll Improvements
- [ ] Add shift+up, shift+down bindings
- [ ] Update help text
- [ ] Document all scroll options
- [ ] Update status bar with accessibility info

---

## Testing Plan

### Unit Tests
- [ ] Test sort logic (alphabetical vs recent)
- [ ] Test letter jump with various lists
- [ ] Test search across file structures
- [ ] Test scroll boundary conditions

### Integration Tests
- [ ] Test sort toggle persistence
- [ ] Test search results navigation
- [ ] Test mouse + keyboard combinations
- [ ] Test help text accuracy

### User Testing
- [ ] Test with non-terminal users (arrow keys, mouse)
- [ ] Test with power users (vim keys)
- [ ] Test across different terminal emulators
- [ ] Verify all keybindings work

---

## Documentation Updates Needed

### Updated Files
- `README.md` - New features
- `help.go` - Help text
- `CONFIG.md` - New keybindings
- `CHANGELOG.md` - Feature list
- `docs/` - Feature guides

### New Documentation
- `docs/SEARCH_FEATURE.md` - Global search guide
- `docs/SORT_FEATURE.md` - Sort modes guide
- `docs/ACCESSIBILITY.md` - Non-terminal user guide

---

## Notes

### Design Decisions
1. Keep all existing vim keybindings (don't break muscle memory)
2. Add accessibility layers on top (shift modifiers, mouse)
3. Provide visual feedback for each feature (status bar, indicators)
4. Make features optional in config (future)

### Future Enhancements (Phase 3+)
- Content search across files (not just filenames)
- Regex support in search
- Saved searches / bookmarks
- Custom sort orders
- Search history
- Mouse selection and copy

---

**Created**: 2025-10-21
**Version**: Phase 2 Feature Specification
**Status**: Ready for Implementation
