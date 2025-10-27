# Mouse Selection Fixes - 2025-10-21

## Investigation Summary

The issue investigation revealed three critical problems with mouse selection functionality:

### 1. **Mouse Selection Not Triggering Handler**
- **Problem**: Mouse events were arriving (Type=11, Action=2) but selection handler code wasn't executing
- **Root Cause**: The handler only checked `Type == 11`, but modern terminals also send `Type == 1` (classic button events)
- **Evidence**: Debug logs showed events being logged but handler code never reached (no output from selection start/extend code)

### 2. **TOC Display Cramped/Scrunched**
- **Problem**: TOC entries were heavily truncated with ellipsis appearing too early
- **Root Cause**: Width calculation was too aggressive: `maxLen := width - len(indent) - 4` left minimal display space
- **Effect**: Entries like "LUXOR - Implementation" became "L..." or "LUXOR - I..."

### 3. **Handler Execution Order Issues**
- **Problem**: File tree click handler might run after viewer selection, potentially intercepting events
- **Analysis**: Handler structure was correct but diagnostic logging was sparse

---

## Fixes Applied

### Fix #1: Enhanced Mouse Event Detection

**File**: `main.go` (lines 242-323)

**Changes**:
1. Check BOTH `Type == 1` AND `Type == 11` for selection
2. For `Type == 1` (classic events): Handle Button=1 with Action 1/2/3
3. For `Type == 11` (modern events): Handle Action=2 drag motion
4. Added detailed debug logging at each step

**Code Pattern**:
```go
// Type 1 = button press (classic events), Type 11 = motion/drag (newer events)
if (msg.Type == 1 || msg.Type == 11) && m.currentView == ViewerView && m.selectedFile != "" {
    // ... selection code

    // Type 1, Button 1 = left button press/drag
    if msg.Type == 1 && msg.Button == 1 {
        if msg.Action == 1 { /* Start selection */ }
        if msg.Action == 2 { /* Extend selection */ }
        if msg.Action == 3 { /* End selection */ }
    }

    // Type 11, Action 2 = modern drag motion events
    if msg.Type == 11 && msg.Action == 2 {
        // ... motion-based selection
    }
}
```

**Debug Logging**:
- `TEXT_SELECTION`: Entry point showing Type, Action, position, and drag state
- `SELECTION_STARTED`: When drag begins at specific coordinates
- `SELECTION_EXTENDED`: When drag extends to new coordinates
- `SELECTION_ENDED`: When drag releases
- `MOTION_SELECTION_STARTED`: For Type=11 motion events
- `MOTION_SELECTION_EXTENDED`: For Type=11 motion drag continuation

### Fix #2: Improved TOC Width Calculation

**File**: `toc.go` (lines 80-112)

**Changes**:
1. Calculate `availableSpace` more carefully: `width - len(prefix) - len(indent) - 1`
2. Only truncate when absolutely necessary (`availableSpace > 4`)
3. Leave proper margin for ellipsis
4. Handle very narrow terminals gracefully

**Before**:
```go
maxLen := width - len(indent) - 4  // Too aggressive
if maxLen > 1 && len(title) > maxLen {
    title = title[:maxLen-1] + "…"
}
```

**After**:
```go
availableSpace := width - len(prefix) - len(indent) - 1
if availableSpace > 4 && len(title) > availableSpace {
    maxLen := availableSpace - 1
    if maxLen > 0 {
        title = title[:maxLen] + "…"
    }
} else if availableSpace <= 4 {
    // Handle very narrow terminals
    if len(title) > availableSpace && availableSpace > 1 {
        title = title[:availableSpace-1] + "…"
    }
}
```

**Impact**:
- TOC entries now display more fully before truncating
- No unnecessary ellipsis when there's space
- Better use of available terminal width

### Fix #3: Enhanced Diagnostic Logging

**File**: `main.go` (line 222-224)

**Changes**:
- Added `currentView` ID (0, 1, or 2) to logging for clarity
- Clarifies which view is active during mouse events
- Helps diagnose execution flow issues

**Before**:
```go
fmt.Fprintf(debugFile, "HANDLING: Type=%d, Button=%d, Action=%d, X=%d, Y=%d | CurrentView=%s, SelectedFile=%v\n",...)
```

**After**:
```go
fmt.Fprintf(debugFile, "HANDLING: Type=%d, Button=%d, Action=%d, X=%d, Y=%d | CurrentView=%s (id=%d), SelectedFile=%v\n",...)
```

---

## Mouse Event Type Reference

Based on investigation findings:

| Type | Event | Button | Action | Usage |
|------|-------|--------|--------|-------|
| 1 | Button press/release | 1 (left) | 1=press, 2=drag, 3=release | Classic event style |
| 5 | Scroll wheel | N/A | N/A | Scroll up |
| 6 | Scroll wheel | N/A | N/A | Scroll down |
| 11 | Motion/drag | 0 | 2=drag motion | Modern drag events |

**Key Insight**: Different terminal emulators send different Type values for the same user action. By checking both Type 1 and Type 11, we support both old and new event styles.

---

## Testing Checklist

After building, verify the following:

### Selection Tests
- [ ] Click and drag in viewer - selection should start immediately
- [ ] Text should highlight during drag (visual feedback from clipboard)
- [ ] Release should end selection and prepare for copy
- [ ] `y` key should copy selected text to clipboard
- [ ] Selection should work on multiple lines
- [ ] Switching views should clear selection

### TOC Display Tests
- [ ] TOC entries should show full text when space allows
- [ ] No premature ellipsis when text fits
- [ ] Deep nested entries (Level 4+) should still be readable
- [ ] Very long headings should truncate gracefully
- [ ] Narrow terminals should handle TOC properly

### Debug Logging Tests
- [ ] Check `/tmp/lumina_mouse_debug.log` for selection events
- [ ] Should see `TEXT_SELECTION` for each event
- [ ] Should see `SELECTION_STARTED`, `SELECTION_EXTENDED`, `SELECTION_ENDED` entries
- [ ] CurrentView ID should match active view

---

## Build Instructions

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
go build -o ccn main.go model.go toc.go clipboard.go search.go \
    keybindings.go help.go colors.go version.go
```

---

## Debug Log Location

```
/tmp/lumina_mouse_debug.log
```

To monitor in real-time:
```bash
tail -f /tmp/lumina_mouse_debug.log
```

Clear before testing:
```bash
rm /tmp/lumina_mouse_debug.log
```

---

## Expected Behavior After Fixes

### Mouse Selection Flow

1. **Click in viewer** → `TEXT_SELECTION` logged, type/action detected
2. **Check Type** → If Type=1 (Button=1) or Type=11, proceed
3. **Action=1 or first motion** → `SELECTION_STARTED` logged
4. **Drag** → Multiple `SELECTION_EXTENDED` entries logged
5. **Release** → `SELECTION_ENDED` logged
6. **Copy with 'y'** → Text copied to clipboard

### TOC Display Behavior

1. **Markdown file loaded** → Headings parsed
2. **TOC pane populated** → Full heading text visible without truncation
3. **Deep nesting** → Proper indentation, text not consumed by indent
4. **Long headings** → Truncate only if necessary, with ellipsis
5. **Navigation** → Up/Down arrows move through TOC smoothly

---

## Related Issues Fixed

- Mouse selection handler not executing despite events arriving
- TOC entries appearing as single characters or heavily truncated
- Handler execution order concerns (verified as not an issue)
- Incomplete diagnostic logging

---

## Version

- **Date**: 2025-10-21
- **Build**: Post-fix (after applying both mouse event and TOC width changes)
- **Status**: Ready for testing

---

## Next Steps

1. Test mouse selection in viewer with both left-click drag and selection
2. Verify TOC entries display properly without truncation
3. Test in narrow and wide terminals
4. Check clipboard functionality with `y` key
5. Verify no performance regression
6. Document final working behavior in COMPLETION_SUMMARY

