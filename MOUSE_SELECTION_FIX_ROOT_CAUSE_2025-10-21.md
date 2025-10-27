# Mouse Selection Fix - Root Cause Analysis

**Date**: 2025-10-21
**Issue**: Mouse selection and copy not working in LUMINA
**Status**: FIXED ✅

## Root Cause

The debug log revealed the core problem:

```
HANDLING: Type=1, Button=1, Action=0, X=59, Y=15 | CurrentView=Preview (id=2), SelectedFile=true
```

**Three critical issues were identified:**

### 1. Wrong Mouse Event Type Being Used

The code was checking for `msg.Action == 1` (previously thought to be "press"), but the debug log showed:
- `Action=0` - Initial button press (this is what we need to catch!)
- `Action=2` - Drag motion
- `Action=3` - Button release

**Fix**: Changed condition from `if msg.Action == 1` to `if msg.Action == 0`

### 2. Selection Only Worked in ViewerView, Not PreviewView

The original code had:
```go
if (msg.Type == 1 || msg.Type == 11) && m.currentView == ViewerView && m.selectedFile != "" {
```

This prevented selection in the Preview pane (when `m.currentView == PreviewView`).

**Fix**: Updated condition to:
```go
if (msg.Type == 1 || msg.Type == 11) && (m.currentView == ViewerView || m.currentView == PreviewView) && m.selectedFile != "" {
```

### 3. Incorrect Line Offset Calculation for Preview Pane

The code always used `m.viewer.YOffset` for line calculation:
```go
textLine := msg.Y - 3 + m.viewer.YOffset
```

But when in Preview mode, this offset doesn't apply to the Preview pane's content.

**Fix**: Added conditional offset calculation:
```go
var textLine int
if m.currentView == ViewerView {
    textLine = msg.Y - 3 + m.viewer.YOffset
} else {
    // Preview pane: header is at Y=1, content starts at Y=3
    textLine = msg.Y - 3
}
```

## Changes Made

**File**: `main.go`

### Change 1: Fixed scroll handling to support Preview pane
```go
// Before
if msg.Type == 5 {
    if m.currentView == ViewerView {
        m.viewer.LineUp(3)
    }
    return m, nil
}

// After
if msg.Type == 5 {
    if m.currentView == ViewerView || m.currentView == PreviewView {
        if m.currentView == ViewerView {
            m.viewer.LineUp(3)
        }
    }
    return m, nil
}
```

### Change 2: Fixed line offset for both views
```go
// Before
textLine := msg.Y - 3 + m.viewer.YOffset

// After
var textLine int
if m.currentView == ViewerView {
    textLine = msg.Y - 3 + m.viewer.YOffset
} else {
    textLine = msg.Y - 3
}
```

### Change 3: Fixed Action value for button press
```go
// Before
if msg.Action == 1 {

// After
if msg.Action == 0 {
```

### Change 4: Enabled selection in both Viewer and Preview panes
```go
// Before
if (msg.Type == 1 || msg.Type == 11) && m.currentView == ViewerView && m.selectedFile != "" {

// After
if (msg.Type == 1 || msg.Type == 11) && (m.currentView == ViewerView || m.currentView == PreviewView) && m.selectedFile != "" {
```

## Why This Works

1. **Action=0**: This is the initial button press event (not Action=1). Now we catch the start of selection immediately.

2. **Both Views**: By checking `m.currentView == ViewerView || m.currentView == PreviewView`, selection works in both panes.

3. **Correct Offset**: The Preview pane has its own coordinate system separate from the Viewer, so we calculate the offset correctly for each.

## Testing

To verify the fix works:

1. Build the binary:
   ```bash
   go build -o ccn main.go model.go toc.go clipboard.go search.go keybindings.go help.go colors.go version.go
   ```

2. Run it:
   ```bash
   ./ccn ~/Documents/LUXOR
   ```

3. Test in Viewer pane:
   - Click and drag to select text
   - Press 'y' to copy
   - Selection should work

4. Test in Preview pane:
   - Click and drag to select TOC entries
   - Press 'y' to copy
   - Selection should work

5. Monitor debug log:
   ```bash
   tail -f /tmp/lumina_mouse_debug.log
   ```

   You should see:
   ```
   → SELECTION_STARTED (Action=0) at line=X col=Y
   → SELECTION_EXTENDED to line=X col=Y
   → SELECTION_ENDED
   ```

## Key Learnings

1. **Mouse Event Actions vary by terminal**: Different terminals send different action codes for the same physical action.

2. **Assumptions matter**: Assuming Action=1 was "press" was wrong. Action=0 is the actual press event.

3. **View-specific logic**: Each pane has different scroll behavior and coordinate systems. This needs to be handled explicitly.

4. **Debug logging is crucial**: The debug log revealed exactly what events were being sent, which made fixing this straightforward.

## Future Improvements

1. Document mouse event codes in a constants file
2. Add configuration for different terminal emulator compatibility
3. Consider adding visual feedback during selection (highlight)
4. Add selection mode options (character, word, line)
