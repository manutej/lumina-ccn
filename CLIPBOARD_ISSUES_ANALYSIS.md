# Clipboard Copy/Paste Issues Analysis

**Date**: November 3, 2025
**Status**: Code Review and Root Cause Analysis
**User Report**: "I still have issues with select > copy > paste"

---

## Current Implementation Review

### File: clipboard.go (236 lines)

**Components**:
1. `SelectionState` struct - tracks selection boundaries
2. `ClipboardManager` - manages copy operations
3. `extractSelection()` - extracts text from selection
4. Helper functions for selection management

---

## Identified Issues

### Issue 1: Empty Selection Handling ❌

**Location**: `CopySelection()` function (lines 64-84)

**Problem**:
```go
// If selection is enabled, use selected text
if cm.selection.Enabled {
    selected := cm.GetSelection(content)
    if selected == "" {
        return fmt.Errorf("selection is empty")
    }
    textToCopy = selected
}
```

**Why It's a Problem**:
- When user clicks without dragging, selection is "Enabled" but Start == End
- `GetSelection()` returns empty string
- Function returns error instead of gracefully fallback to copying entire document
- **User Experience**: User sees no feedback, clipboard unchanged

**Expected Behavior**:
- If selection is empty, copy entire document as fallback
- OR: Provide clear visual feedback that selection is empty

---

### Issue 2: Selection Visual Feedback Missing ❌

**Location**: main.go, handleMouseEvent()

**Problem**:
- Selection state is tracked internally
- NO visual highlighting of selected text in viewer
- User has no idea what they've selected

**User Experience Impact**:
- User drags mouse → No visual feedback
- User presses 'y' → No idea what was copied
- This breaks the fundamental UX principle: "Show, don't guess"

**Current Code** (main.go:239-248):
```go
if err := m.clipboard.CopySelection(m.viewerContent); err == nil {
    // Copied! Could show status message
    // For now, selection clears
    m.clipboard.ClearSelection()
}
```

**Issues**:
1. Comment says "could show status message" - but doesn't
2. Selection clears silently
3. No confirmation to user that copy succeeded

---

### Issue 3: Selection Boundary Edge Cases ⚠️

**Location**: `extractSelection()` function (lines 124-185)

**Problem Areas**:

#### 3a: Line 177 - EndCol == 0 Edge Case
```go
if endCol > 0 {
    if endCol > len(line) {
        endCol = len(line)
    }
    result += line[:endCol]
}
```

**Issue**: If `endCol == 0`, nothing is added from the last line. This might be intentional, but it's inconsistent with user expectation when selecting at the start of a line.

#### 3b: Lines 136-139 - Boundary Normalization
```go
if startLine > endLine || (startLine == endLine && startCol > endCol) {
    startLine, endLine = endLine, startLine
    startCol, endCol = endCol, startCol
}
```

**Issue**: This normalizes reversed selections, but doesn't check if the resulting selection is empty (startLine == endLine && startCol == endCol).

---

### Issue 4: No Clipboard Access Feedback ❌

**Location**: Multiple locations

**Problem**:
- Clipboard operations can fail (no xclip on Linux, permissions, etc.)
- Failures are logged to debug file but not shown to user
- User thinks copy worked, but clipboard is empty

**Example** (main.go:250-254):
```go
} else {
    debugFile, _ := os.OpenFile("/tmp/lumina_mouse_debug.log", ...)
    if debugFile != nil {
        fmt.Fprintf(debugFile, "  → COPY FAILED: %v\n", err)
        debugFile.Close()
    }
}
```

**User Impact**:
- Silent failure
- No error message
- No way to know copy didn't work until trying to paste

---

### Issue 5: Selection Persistence Confusion ⚠️

**Location**: main.go:173-174

**Code**:
```go
case "switch_view":
    m.currentView = (m.currentView + 1) % 3
    m.clipboard.ClearSelection() // Clear selection when switching views
```

**Potential Issue**:
- Selection clears when switching views
- But what if user selected text, switched view, then pressed 'y'?
- Expected: Copy what was selected
- Actual: Copy fails or copies nothing

---

### Issue 6: Mouse Coordinate Transform Accuracy ⚠️

**Location**: main.go:365-379, `screenToDocCoords()`

**Potential Issue**:
```go
func (m *AppModel) screenToDocCoords(screenX, screenY int) (line, col int) {
    // ...
    line = paneRelativeY + m.viewer.YOffset
    col = paneRelativeX
```

**Concerns**:
- Does not account for wrapped lines in rendered markdown
- Character-level selection might not align with displayed text
- Rendered markdown (with formatting) vs raw markdown (stored in viewerContent)

**Example**:
- User sees: "**Bold Text**" (rendered bold)
- Raw content: "**Bold Text**"
- Selection coordinates might be off by 4 chars (the ** markers)

---

## Root Cause Summary

The primary issue is **lack of visual feedback and graceful error handling**:

1. ❌ **No Selection Highlighting** - Users can't see what they've selected
2. ❌ **Silent Failures** - Errors don't show in UI
3. ❌ **Empty Selection Not Handled** - Clicking without dragging breaks copy
4. ❌ **No Copy Confirmation** - Users don't know if copy worked

---

## Proposed Fixes

### Fix 1: Handle Empty Selection Gracefully

**File**: clipboard.go, `CopySelection()`

**Change**:
```go
func (cm *ClipboardManager) CopySelection(content string) error {
    var textToCopy string

    // If selection is enabled, use selected text
    if cm.selection.Enabled {
        selected := cm.GetSelection(content)
        if selected == "" {
            // Empty selection - fall back to copying entire content
            textToCopy = content
        } else {
            textToCopy = selected
        }
    } else {
        // No selection - copy entire content
        if content == "" {
            return fmt.Errorf("no content to copy")
        }
        textToCopy = content
    }

    cm.lastCopy = textToCopy
    return clipboard.WriteAll(textToCopy)
}
```

**Impact**: Clicking without dragging will copy entire document instead of failing.

---

### Fix 2: Add Visual Selection Highlighting

**File**: main.go, `View()` function

**Change**: Add selection highlighting to viewer viewport

**Pseudocode**:
```go
if m.clipboard.HasSelection() {
    startLine, startCol, endLine, endCol := m.clipboard.GetSelectionBounds()

    // Highlight selected text with different background color
    // Use lipgloss to style selected lines
    highlightStyle := lipgloss.NewStyle().Background(lipgloss.Color("226")) // Yellow

    // Apply highlighting to viewer content
    // ... (implementation details)
}
```

**Impact**: Users can see what they're selecting in real-time.

---

### Fix 3: Add Status Bar Feedback

**File**: model.go, add `statusMessage` field

**Changes**:
1. Add `statusMessage string` to AppModel
2. After copy, set `statusMessage = "Copied!"`
3. Display in status bar
4. Clear after 2 seconds

**Example**:
```go
// In Update() after copy
case "copy":
    if m.currentView == ViewerView {
        if err := m.clipboard.CopySelection(m.viewerContent); err == nil {
            m.statusMessage = "✓ Copied to clipboard!"
            return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
                return clearStatusMsg{}
            })
        } else {
            m.statusMessage = "✗ Copy failed: " + err.Error()
        }
    }
```

**Impact**: Clear feedback on copy success/failure.

---

### Fix 4: Improve Selection Boundary Checks

**File**: clipboard.go, `extractSelection()`

**Change**: Add check for empty selection after normalization

```go
// After normalization (line 139)
// Check for empty selection
if startLine == endLine && startCol == endCol {
    return "" // Empty selection
}
```

**Impact**: Consistent behavior for edge cases.

---

### Fix 5: Add Keyboard Status Bar Hints

**File**: main.go, status bar rendering

**Change**: Update hints based on selection state

```go
if m.clipboard.HasSelection() {
    statusHints = "y: copy selection | Esc: clear selection | ?: help | q: quit"
} else {
    statusHints = "y: copy all | ?: help | q: quit"
}
```

**Impact**: Users understand what 'y' will do based on context.

---

## Testing Checklist

### Manual Tests (Must Pass)

**Test 1: Copy Without Selection**
- [ ] Open ccn
- [ ] Navigate to any markdown file
- [ ] Press 'y' (no mouse selection)
- [ ] Paste in external editor
- [ ] **Expected**: Full document content pastes
- [ ] **Expected**: Status bar shows "✓ Copied to clipboard!"

**Test 2: Mouse Selection → Copy**
- [ ] Open ccn
- [ ] Navigate to any markdown file
- [ ] Click and drag to select 2-3 lines
- [ ] **Expected**: Selection is visually highlighted
- [ ] Press 'y'
- [ ] **Expected**: Status bar shows "✓ Copied to clipboard!"
- [ ] Paste in external editor
- [ ] **Expected**: Only selected lines paste

**Test 3: Empty Selection (Click Without Drag)**
- [ ] Open ccn
- [ ] Click once in viewer (don't drag)
- [ ] Press 'y'
- [ ] Paste in external editor
- [ ] **Expected**: Full document OR previous clipboard (no crash)

**Test 4: Copy Then Switch Views**
- [ ] Select text
- [ ] Switch to file tree view (Tab)
- [ ] Press 'y'
- [ ] **Expected**: Either copy works or clear error message

**Test 5: Copy Large File**
- [ ] Open file with 1000+ lines
- [ ] Press 'y'
- [ ] **Expected**: Copy completes in < 500ms
- [ ] **Expected**: Status confirmation shown

**Test 6: Copy Special Characters**
- [ ] Open file with emojis, unicode, code blocks
- [ ] Press 'y'
- [ ] Paste
- [ ] **Expected**: All characters preserved

---

## Priority Order

### P0 (Critical - Must Fix)
1. ✅ Fix empty selection handling (Fix 1)
2. ✅ Add status bar feedback (Fix 3)
3. ✅ Add selection visual highlighting (Fix 2)

### P1 (High - Should Fix)
4. ⏳ Improve keyboard hints (Fix 5)
5. ⏳ Better boundary checks (Fix 4)

### P2 (Medium - Nice to Have)
6. ⏳ Handle mouse coordinate transforms better (Issue 6)
7. ⏳ Add copy history

---

## Success Metrics

✅ **Feature is DONE when**:
1. User presses 'y' → Always copies something (never fails silently)
2. User sees visual feedback for selection (highlighted text)
3. User sees confirmation message after copy ("✓ Copied!")
4. All manual tests pass 100% of the time
5. Edge cases handled gracefully (empty selection, large files, etc.)

---

## Implementation Order

1. **Phase 1**: Fix empty selection handling (clipboard.go)
2. **Phase 2**: Add status message system (model.go, main.go)
3. **Phase 3**: Add selection visual highlighting (main.go, View())
4. **Phase 4**: Update keyboard hints (main.go)
5. **Phase 5**: Manual testing and verification

---

## Estimated Time

- Analysis: ✅ Complete
- Implementation: ~2-3 hours
- Testing: ~1 hour
- **Total**: 3-4 hours

---

**Next Action**: Begin Phase 1 implementation (Fix empty selection handling)

**Last Updated**: 2025-11-03
**Status**: Ready for Implementation
