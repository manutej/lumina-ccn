# Mouse Selection & TOC Fixes - Verification Guide

## Quick Start

### Build Latest
```bash
cd ~/Documents/LUXOR/PROJECTS/LUMINA/ccn
go build -o ccn main.go model.go toc.go clipboard.go search.go \
    keybindings.go help.go colors.go version.go
```

### Monitor Debug Log
```bash
# In separate terminal
tail -f /tmp/lumina_mouse_debug.log
```

### Run CCN
```bash
./ccn ~/Documents/LUXOR
```

---

## Issue #1: Mouse Selection Not Triggering

### What Was Fixed
- Added support for **both** Type=1 and Type=11 mouse events
- Type=1 events: Classic button press/drag/release
- Type=11 events: Modern motion/drag events
- Added detailed logging at each step

### How to Verify

1. **Switch to Viewer** (press `Tab` until "VIEWER" is highlighted)
2. **Select a file** from the tree (press Enter)
3. **Click and drag** in the text area to select text
4. **Expected behavior**:
   - Text should highlight during drag
   - Selection should persist until cleared
   - Debug log should show `SELECTION_STARTED`, `SELECTION_EXTENDED`, `SELECTION_ENDED`

### Debug Output

```
HANDLING: Type=1, Button=1, Action=1, X=30, Y=10 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=1, textLine=15, textCol=30, dragActive=false
  → SELECTION_STARTED at line=15 col=30

HANDLING: Type=1, Button=1, Action=2, X=35, Y=12 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=2, textLine=17, textCol=35, dragActive=true
  → SELECTION_EXTENDED to line=17 col=35

HANDLING: Type=1, Button=1, Action=3, X=35, Y=12 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=3, textLine=17, textCol=35, dragActive=true
  → SELECTION_ENDED
```

### Pass Criteria
- [ ] Selection starts immediately on click
- [ ] Drag extends selection
- [ ] Release ends selection
- [ ] Debug log shows all three stages
- [ ] No error messages

---

## Issue #2: TOC Display Cramped/Scrunched

### What Was Fixed
- Improved width calculation for TOC entries
- Changed from `width - indent - 4` to `width - prefix - indent - 1`
- Only truncate when necessary
- Preserve full text when space allows

### How to Verify

1. **Switch to Preview** (press `Tab` twice, should see "PREVIEW")
2. **Select a markdown file** (switch to FileTree tab, select .md file, press Enter)
3. **Switch back to Preview** to see TOC
4. **Expected behavior**:
   - Headings should display in full (not like "LUXOR - I...")
   - Deep nested headings (####) should be readable
   - Long headings truncate gracefully with "…"
   - No unnecessary ellipsis when text fits

### Examples

**Before (cramped)**:
```
📋 Table of Contents

  LUXOR - I...
  Ov...
  Ar...
  Ge...
```

**After (readable)**:
```
📋 Table of Contents

  LUXOR - Implementation
  Overview
  Architecture
  Getting Started
    Introduction
    Prerequisites
```

### Pass Criteria
- [ ] Full heading text visible
- [ ] Proper indentation for nested items
- [ ] No premature truncation
- [ ] Ellipsis only when necessary
- [ ] Readable in narrow terminals

---

## Issue #3: Handler Execution Order

### What Was Fixed
- Added enhanced diagnostic logging showing CurrentView ID (0=FileTree, 1=Viewer, 2=Preview)
- Verified handler structure is correct
- Selection handling runs BEFORE file tree handling

### How to Verify

1. **In debug log**, look for "CurrentView=" entries
2. **When in Viewer**: Should show "id=1"
3. **When selecting in FileTree**: Should show "id=0"
4. **Verify no event interception**: Selection events in Viewer should NOT redirect to FileTree

### Debug Output Pattern

```
CurrentView=Viewer (id=1) - Selection events go to viewer
CurrentView=FileTree (id=0) - Click events go to file tree
CurrentView=Preview (id=2) - Focus on preview pane
```

### Pass Criteria
- [ ] Correct view ID logged for each event
- [ ] No event redirection between views
- [ ] Selection in viewer doesn't affect file tree
- [ ] View switching works smoothly

---

## Complete Test Scenario

### Setup
```bash
cd ~/Documents/LUXOR/PROJECTS/LUMINA/ccn
rm /tmp/lumina_mouse_debug.log  # Clear old logs
./ccn ~/Documents/LUXOR
```

### Test Steps

**Step 1: Navigate to a File**
1. (Should be on FileTree by default)
2. Use `j/k` to navigate down to a markdown file (.md)
3. Press `Enter` to view it

**Step 2: Test Selection**
1. Press `Tab` to switch to Viewer
2. Click and drag across some text
3. Expected: Text highlights during drag
4. Check `/tmp/lumina_mouse_debug.log`:
   ```bash
   tail -20 /tmp/lumina_mouse_debug.log
   ```

**Step 3: Test Copy**
1. With text selected from Step 2
2. Press `y` to copy
3. Try `Ctrl+Shift+V` or `Cmd+V` (platform dependent) to paste

**Step 4: Test TOC**
1. Press `Tab` twice to get to Preview
2. Should see table of contents for selected file
3. Verify headings display fully without premature truncation
4. Press `j/k` to navigate through TOC

**Step 5: Test Narrow Terminal**
1. Resize terminal window to ~80 characters wide
2. Switch between views
3. Verify TOC still displays readably
4. Verify selection still works

---

## Expected Logs Pattern

### Good Selection Log
```
HANDLING: Type=1, Button=1, Action=1, X=40, Y=8 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=1, textLine=5, textCol=40, dragActive=false
  → SELECTION_STARTED at line=5 col=40

HANDLING: Type=1, Button=1, Action=2, X=45, Y=8 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=2, textLine=5, textCol=45, dragActive=true
  → SELECTION_EXTENDED to line=5 col=45

HANDLING: Type=1, Button=1, Action=3, X=45, Y=8 | CurrentView=Viewer (id=1), SelectedFile=true
  → TEXT_SELECTION: Type=1, Action=3, textLine=5, textCol=45, dragActive=true
  → SELECTION_ENDED
```

### Bad Selection Log (Problem Indicator)
```
HANDLING: Type=11, Button=0, Action=2, X=40, Y=8 | CurrentView=FileTree (id=0), SelectedFile=true
  (No TEXT_SELECTION output - handler didn't execute)
```
> **Problem**: Type=11 but CurrentView=0 (FileTree, not Viewer). Selection handler requires Viewer.

---

## Troubleshooting

### Selection Not Working

1. **Check debug log exists**:
   ```bash
   ls -lh /tmp/lumina_mouse_debug.log
   ```

2. **Verify you're in Viewer**:
   - Status bar should show "[VIEWER]"
   - If not, press Tab to switch

3. **Verify file is selected**:
   - Must have a file selected before selection works
   - View should show file content

4. **Check if Type=11 events are being skipped**:
   ```bash
   grep "Type=11" /tmp/lumina_mouse_debug.log
   ```
   > If mostly Type=11, check for corresponding TEXT_SELECTION entries. If missing, the check for Type=11 isn't working.

### TOC Still Cramped

1. **Check terminal width**:
   ```bash
   echo $COLUMNS
   ```
   > Should be > 30 for reasonable TOC display

2. **Verify headings exist**:
   - TOC should show "No headings found" if none
   - Select a markdown file with #headers

3. **Check for very long titles**:
   - Headings over 100 chars might still truncate
   - This is correct behavior

### View Not Switching

1. **Press Tab** multiple times to cycle through views
2. **Status bar** should show active view
3. **Border** of active pane should be bright

---

## Success Criteria Summary

### All Fixes Working ✓
- [ ] Mouse selection starts on click in Viewer
- [ ] Selection extends through drag
- [ ] Selection ends on release
- [ ] Debug log shows all three stages
- [ ] TOC entries display fully without premature truncation
- [ ] CurrentView ID correctly identifies active pane
- [ ] No event interception between views
- [ ] Narrow terminals handled gracefully
- [ ] Copy functionality works (y key)

---

## Files Modified

1. **main.go** (lines 217-333)
   - Enhanced mouse event handling
   - Support for Type=1 and Type=11
   - Detailed logging

2. **toc.go** (lines 80-112)
   - Improved width calculation
   - Better truncation logic
   - Preserve text when space allows

---

## Next Steps After Verification

1. Document final results in COMPLETION_SUMMARY_2025-10-21.md
2. Commit changes if all tests pass
3. Update VERSION or CHANGELOG
4. Consider Phase 2 features

