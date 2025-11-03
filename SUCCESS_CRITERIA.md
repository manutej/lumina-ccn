# SUCCESS CRITERIA - Lumina CCN Features

**Date**: November 3, 2025
**Status**: Test-Driven Development Reference
**Purpose**: Define clear, testable success criteria BEFORE implementation

---

## Philosophy

> "Success criteria are very important!"
> — User requirement, 2025-11-03

All features must have:
1. **Clear success criteria** (what does "done" mean?)
2. **Test plan** (how do we verify it works?)
3. **Failure scenarios** (what are the edge cases?)
4. **User acceptance criteria** (when is the user happy?)

---

## Feature 1: Copy/Paste Functionality

### Current Status: ❌ NOT WORKING PROPERLY

**User Report**: "I still have issues with select > copy > paste"

### Success Criteria

#### SC1.1: Copy Entire Document (No Selection)
**Given**: User is viewing a markdown file
**When**: User presses 'y' key
**Then**:
- ✅ Entire document content is copied to system clipboard
- ✅ User can paste into external app (VS Code, Notes, etc.)
- ✅ Content preserves line breaks and formatting
- ✅ No error messages shown

**Test Command**:
```bash
# 1. Open lumina
./ccn docs/

# 2. Navigate to any .md file
# 3. Press 'y'
# 4. Open external editor
# 5. Cmd+V (or Ctrl+V on Linux)
# Expected: Full document content appears
```

#### SC1.2: Mouse Selection → Copy
**Given**: User is viewing a markdown file
**When**:
1. User clicks and drags to select text
2. User presses 'y' key
**Then**:
- ✅ Only selected text is copied to clipboard
- ✅ Selection is visually highlighted during drag
- ✅ User can paste selection into external app
- ✅ Selection clears after copy (visual feedback)

**Test Command**:
```bash
# 1. Open lumina
./ccn docs/

# 2. Navigate to any .md file
# 3. Click and drag to select 2-3 lines
# 4. Press 'y'
# 5. Open external editor
# 6. Cmd+V (or Ctrl+V)
# Expected: Only selected lines appear
```

#### SC1.3: Copy Empty Content (Edge Case)
**Given**: User is viewing an empty file
**When**: User presses 'y'
**Then**:
- ✅ No crash or hang
- ✅ Clipboard is not modified (or empty string is copied)
- ✅ User sees status message: "No content to copy"

#### SC1.4: Copy Very Large File (Performance)
**Given**: User is viewing a 10,000+ line markdown file
**When**: User presses 'y'
**Then**:
- ✅ Copy completes in < 500ms
- ✅ No UI freeze or lag
- ✅ Full content is in clipboard

#### SC1.5: Copy with Special Characters
**Given**: Document contains code blocks, emojis, unicode
**When**: User presses 'y'
**Then**:
- ✅ All special characters preserved
- ✅ Code blocks maintain formatting
- ✅ No character corruption

#### SC1.6: Selection Boundaries (Edge Cases)
**Given**: User selects text at document boundaries
**When**: User drags from:
- First character to last character
- Middle of line to middle of another line
- Single word (double-click)
**Then**:
- ✅ Selection respects boundaries (no overflow)
- ✅ Copy works correctly for all selection types

### Failure Scenarios

❌ **Current Issues to Fix**:
1. Selection not visually highlighted → Cannot see what's selected
2. Copy fails silently → No feedback to user
3. Clipboard empty after paste → Data not reaching system clipboard
4. Selection state persists → Confusion about what's selected

### Test Plan

#### Manual Tests

**Test 1: Basic Copy (No Selection)**
```
Steps:
1. Build: go build -o ccn .
2. Run: ./ccn docs/
3. Navigate to CHANGELOG.md
4. Press 'y'
5. Open external editor
6. Paste (Cmd+V)

Expected: Full CHANGELOG.md content appears
Pass/Fail: ___
```

**Test 2: Mouse Selection Copy**
```
Steps:
1. Open ccn with docs/
2. Navigate to any markdown file
3. Click at line 5, drag to line 8
4. Press 'y'
5. Paste into external editor

Expected: Lines 5-8 appear
Pass/Fail: ___
```

**Test 3: Empty Selection**
```
Steps:
1. Open ccn with docs/
2. Navigate to any markdown file
3. Click without dragging (single click)
4. Press 'y'
5. Paste into external editor

Expected: Either full document OR previous clipboard content (not crash)
Pass/Fail: ___
```

**Test 4: Copy Then Navigate**
```
Steps:
1. Open ccn with docs/
2. Navigate to file A
3. Press 'y'
4. Navigate to file B
5. Paste into external editor

Expected: File A content (not File B)
Pass/Fail: ___
```

#### Automated Tests

**clipboard_test.go** (to be created):
```go
// TestCopyEntireDocument verifies copying without selection
// TestCopySelection verifies copying selected text
// TestCopyEmpty verifies handling empty content
// TestSelectionBoundaries verifies edge cases
// TestClipboardIntegration verifies system clipboard access
```

### User Acceptance Criteria

✅ **Feature is DONE when**:
1. User can press 'y' and paste document anywhere
2. User can select text with mouse and copy selection
3. User sees visual feedback during selection
4. Copy works 100% of the time (no silent failures)
5. Status bar shows "Copied!" confirmation message
6. All edge cases handled gracefully (empty, large files, etc.)

---

## Feature 2: Fuzzy Finder (Phase 3 Week 1)

### Current Status: ✅ IMPLEMENTED, NEEDS TESTING

### Success Criteria

#### SC2.1: Activate Fuzzy Finder
**Given**: User is viewing any pane
**When**: User presses '/'
**Then**:
- ✅ Modal appears centered on screen
- ✅ Input field is focused and ready
- ✅ File list shows all markdown files
- ✅ Base UI is still visible (dimmed background)

#### SC2.2: Filter Files in Real-Time
**Given**: Fuzzy finder is open
**When**: User types "phase"
**Then**:
- ✅ Results filter to match "phase"
- ✅ Filtering happens instantly (< 50ms)
- ✅ Match count updates in UI
- ✅ First match is auto-selected

#### SC2.3: Navigate Results
**Given**: Fuzzy finder shows 10 results
**When**: User presses ↓ five times
**Then**:
- ✅ Selection moves down one item per press
- ✅ Selected item is visually highlighted
- ✅ Selection wraps at boundaries

#### SC2.4: Select and Load File
**Given**: Fuzzy finder is showing results
**When**: User presses Enter
**Then**:
- ✅ Modal closes
- ✅ Selected file loads in viewer
- ✅ Viewer shows rendered markdown
- ✅ Normal mode resumes

#### SC2.5: Cancel Finder
**Given**: Fuzzy finder is open
**When**: User presses Esc
**Then**:
- ✅ Modal closes
- ✅ No file is loaded
- ✅ Previous state preserved
- ✅ Normal mode resumes

#### SC2.6: Performance with Large File Lists
**Given**: Directory has 1000+ markdown files
**When**: User types query
**Then**:
- ✅ Filtering completes in < 100ms
- ✅ UI remains responsive
- ✅ No lag or freeze

### Test Plan

#### Manual Tests

**Test 1: Basic Activation**
```
Steps:
1. Build and run ccn
2. Press '/'

Expected: Modal appears with file list
Pass/Fail: ___
```

**Test 2: Real-Time Filtering**
```
Steps:
1. Press '/'
2. Type "phase"
3. Observe results update

Expected: Files matching "phase" shown instantly
Pass/Fail: ___
```

**Test 3: Navigation and Selection**
```
Steps:
1. Press '/'
2. Type "doc"
3. Press ↓ three times
4. Press Enter

Expected: Fourth matching file loads
Pass/Fail: ___
```

**Test 4: Cancel Without Selection**
```
Steps:
1. Press '/'
2. Type "test"
3. Press Esc

Expected: Modal closes, no file loaded
Pass/Fail: ___
```

#### Automated Tests

**fuzzy_finder_test.go** (already exists):
- ✅ 310 lines of tests
- ✅ Tests filtering, navigation, selection
- ✅ Tests performance benchmarks

### User Acceptance Criteria

✅ **Feature is DONE when**:
1. User can open finder with '/' from any view
2. User can type and see instant filtering
3. User can navigate with ↑/↓ and select with Enter
4. User can cancel with Esc
5. Performance is smooth even with 1000+ files
6. Help text is clear and visible

---

## Feature 3: Ripgrep Content Search (Phase 3 Week 2)

### Current Status: ⏳ BACKEND COMPLETE, UI PENDING

### Success Criteria

#### SC3.1: Activate Content Search
**Given**: User is viewing any pane
**When**: User presses Ctrl+F (or defined keybinding)
**Then**:
- ✅ Search modal appears
- ✅ Input field focused
- ✅ Status shows "Search: "
- ✅ Results pane ready

#### SC3.2: Stream Search Results
**Given**: Search modal is open
**When**: User types search query and presses Enter
**Then**:
- ✅ Results stream in real-time (not blocking)
- ✅ Each result shows: filename, line number, match context
- ✅ Match text is highlighted
- ✅ Results are scrollable

#### SC3.3: Jump to Match
**Given**: Search results are displayed
**When**: User presses Enter on a result
**Then**:
- ✅ File loads in viewer
- ✅ Viewer scrolls to matching line
- ✅ Match is visually highlighted
- ✅ Search modal closes

#### SC3.4: Search Performance
**Given**: Repository has 100+ markdown files
**When**: User searches for common term
**Then**:
- ✅ First results appear in < 200ms
- ✅ Full search completes in < 2 seconds
- ✅ UI remains responsive during search
- ✅ User can cancel search mid-execution

#### SC3.5: Search Options
**Given**: User is in search modal
**When**: User toggles options:
- Case sensitive / insensitive
- Whole word / partial
- File type filters
**Then**:
- ✅ Options are clearly indicated in UI
- ✅ Results update based on options
- ✅ Options persist across searches

### Test Plan (To Be Defined)

**NOTE**: This feature is NOT YET implemented in UI.
Tests will be defined during Week 2 implementation.

---

## Feature 4: File Watching & Auto-Reload

### Current Status: ⏳ BACKEND COMPLETE, UI PENDING

### Success Criteria

#### SC4.1: Detect File Changes
**Given**: User is viewing a markdown file
**When**: File is modified externally (e.g., in VS Code)
**Then**:
- ✅ Change detected within 500ms
- ✅ UI shows "File changed" indicator
- ✅ User is prompted to reload

#### SC4.2: Auto-Reload on User Confirmation
**Given**: File change detected
**When**: User confirms reload
**Then**:
- ✅ Viewer updates with new content
- ✅ Scroll position preserved (if possible)
- ✅ No UI flicker or flash

#### SC4.3: Watch Multiple Files
**Given**: User navigates between multiple files
**When**: Any watched file changes
**Then**:
- ✅ Changes are tracked per file
- ✅ Reload happens for correct file
- ✅ No memory leaks from watchers

### Test Plan (To Be Defined)

**NOTE**: This feature is NOT YET implemented in UI.
Tests will be defined during Week 2+ implementation.

---

## Cross-Feature Success Criteria

### Performance
- ✅ App starts in < 500ms
- ✅ File navigation responds in < 50ms
- ✅ Copy operations complete in < 100ms
- ✅ Fuzzy search filters in < 50ms
- ✅ Content search returns first results in < 200ms

### Stability
- ✅ No crashes under normal use
- ✅ No memory leaks during extended sessions
- ✅ Graceful handling of all edge cases
- ✅ Clear error messages for failures

### User Experience
- ✅ Visual feedback for all actions
- ✅ Status bar shows current context
- ✅ Help overlay accessible with '?'
- ✅ Keybindings are discoverable and consistent

---

## Testing Checklist

### Pre-Release Testing

**Phase 1: Core Functionality**
- [ ] Copy entire document (no selection)
- [ ] Copy with mouse selection
- [ ] Fuzzy finder activation and filtering
- [ ] Fuzzy finder navigation and selection
- [ ] File tree navigation
- [ ] Markdown rendering

**Phase 2: Edge Cases**
- [ ] Copy empty file
- [ ] Copy very large file (10,000+ lines)
- [ ] Fuzzy finder with 1000+ files
- [ ] Selection at document boundaries
- [ ] Cancel operations (Esc key)

**Phase 3: Cross-Platform**
- [ ] macOS: Copy/paste works with pbcopy
- [ ] Linux: Copy/paste works with xclip/xsel
- [ ] Terminal compatibility (iTerm, Terminal.app, Alacritty)

**Phase 4: Regression Testing**
- [ ] All Phase 1 features still work
- [ ] Keybindings configuration still loads
- [ ] Colors and themes work correctly
- [ ] Help overlay displays correctly

---

## Definition of "Done"

A feature is **DONE** when:

1. ✅ All success criteria are met
2. ✅ All manual tests pass
3. ✅ Automated tests exist and pass
4. ✅ Edge cases are handled
5. ✅ Performance targets are met
6. ✅ User acceptance criteria satisfied
7. ✅ Documentation is updated
8. ✅ No known bugs remain

---

## Current Priority: Fix Copy/Paste

**Why**: User explicitly stated "I still have issues with select > copy > paste"

**Next Steps**:
1. ✅ Define success criteria (this document)
2. ⏳ Write comprehensive tests
3. ⏳ Debug and fix implementation
4. ⏳ Manual verification
5. ⏳ Commit with proof of testing

**Success Metric**: User can copy/paste text 100% of the time with zero issues.

---

**Last Updated**: 2025-11-03
**Next Review**: After copy/paste fix is complete
