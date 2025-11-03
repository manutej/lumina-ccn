# Manual Test Checklist - Copy/Paste Fixes

**Date**: November 3, 2025
**Build**: `ccn` (15MB binary)
**Status**: Ready for Manual Testing

---

## 🎯 What Was Fixed

### Critical Fixes Implemented:
1. ✅ **Empty Selection Handling** - Clicking without dragging now copies entire document (no longer fails)
2. ✅ **Status Message System** - Visual feedback when copy succeeds or fails
3. ✅ **Graceful Error Handling** - No more silent failures

### Files Modified:
- `clipboard.go` - Fixed CopySelection() to handle empty selections gracefully
- `model.go` - Added statusMessage field
- `main.go` - Added status message display and timer system

---

## 📋 Pre-Test Setup

### 1. Verify Build
```bash
cd /home/user/lumina-ccn
ls -lh ccn
# Should show: ccn (15MB binary) with recent timestamp
```

### 2. Prepare Test Files
```bash
# Use existing documentation
cd /home/user/lumina-ccn
./ccn docs/
```

---

## ✅ Test Suite

### Test 1: Copy Entire Document (No Selection)
**Priority**: P0 (Critical)

**Steps**:
1. Run: `./ccn docs/`
2. Navigate to any `.md` file (use arrow keys)
3. Press `Enter` to open file
4. Press `Tab` to switch to VIEWER pane
5. **Press `y`** (do NOT select any text first)
6. **Check status bar** → Should show: `✓ Copied document to clipboard!`
7. Open external editor (VS Code, Notes, nano, etc.)
8. Paste (`Cmd+V` or `Ctrl+V`)

**Expected Result**:
- ✅ Status bar shows success message for 2 seconds
- ✅ Full document content appears in external editor
- ✅ All formatting, line breaks, special characters preserved
- ✅ No error messages

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 2: Copy with Mouse Selection
**Priority**: P0 (Critical)

**Steps**:
1. Run: `./ccn docs/`
2. Open any markdown file
3. Switch to VIEWER pane (`Tab`)
4. **Click and drag** to select 2-3 lines of text
5. Press `y`
6. **Check status bar** → Should show: `✓ Copied selection to clipboard!`
7. Paste in external editor

**Expected Result**:
- ✅ Status bar shows "Copied selection" message
- ✅ ONLY selected lines appear in external editor (not full document)
- ✅ Selection clears after copy

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 3: Empty Selection (Click Without Drag)
**Priority**: P0 (Critical)

**Steps**:
1. Run: `./ccn docs/`
2. Open any markdown file
3. Switch to VIEWER pane
4. **Single click** in viewer (don't drag)
5. Press `y`
6. **Check status bar** → Should show success message
7. Paste in external editor

**Expected Result**:
- ✅ No crash or hang
- ✅ Status bar shows success message
- ✅ Either full document OR nothing is copied (no error)
- ✅ No silent failure

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 4: Copy Error Handling
**Priority**: P1 (High)

**Steps**:
1. Run: `./ccn docs/`
2. Navigate to empty file (or create one)
3. Try to copy empty file with `y`
4. **Check status bar** → Should show error message

**Expected Result**:
- ✅ Status bar shows: `✗ Copy failed: no content to copy`
- ✅ Error message displays for 3 seconds (longer than success)
- ✅ No crash

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 5: Copy Large File (Performance)
**Priority**: P1 (High)

**Steps**:
1. Run: `./ccn docs/`
2. Open large file (e.g., `CHANGELOG.md` or `SUCCESS_CRITERIA.md`)
3. Press `y`
4. **Time the response** - Should be instant

**Expected Result**:
- ✅ Copy completes in < 500ms
- ✅ No UI freeze or lag
- ✅ Status message appears immediately
- ✅ Full content in clipboard

**Pass/Fail**: ___________

**Timing**: ________ms

**Notes**:
____________________________________________
____________________________________________

---

### Test 6: Copy Special Characters
**Priority**: P1 (High)

**Steps**:
1. Create test file with special content:
```markdown
# Test Special Characters

**Bold Text**
*Italic Text*
`code snippet`

```go
func main() {
    fmt.Println("Hello 👋")
}
```

Unicode: 日本語 テスト
Emojis: 🎯 ✅ ❌ 🚀
Symbols: © ® ™ € £
```
2. Open file in ccn
3. Press `y` to copy
4. Paste in external editor

**Expected Result**:
- ✅ All special characters preserved
- ✅ Emojis display correctly
- ✅ Code blocks maintain formatting
- ✅ No character corruption

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 7: Status Message Auto-Clear
**Priority**: P2 (Medium)

**Steps**:
1. Run: `./ccn docs/`
2. Open any file
3. Press `y` to copy
4. **Wait and observe status bar**
5. Count how long message stays visible

**Expected Result**:
- ✅ Success message displays for ~2 seconds
- ✅ Message automatically clears
- ✅ Status bar returns to normal hints
- ✅ No flickering or UI glitches

**Pass/Fail**: ___________

**Duration**: ________seconds

**Notes**:
____________________________________________
____________________________________________

---

### Test 8: Multiple Copy Operations
**Priority**: P2 (Medium)

**Steps**:
1. Run: `./ccn docs/`
2. Open file A, press `y`
3. Navigate to file B, press `y`
4. Navigate to file C, press `y`
5. Paste in external editor

**Expected Result**:
- ✅ Each copy shows status message
- ✅ Clipboard contains LAST copied file (file C)
- ✅ No memory leaks or slowdown
- ✅ Status messages clear properly

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 9: Copy Then Switch Views
**Priority**: P2 (Medium)

**Steps**:
1. Run: `./ccn docs/`
2. Open file in VIEWER
3. Select text (optional)
4. Press `Tab` to switch to FILE TREE
5. Press `y`
6. **Check behavior**

**Expected Result**:
- ✅ Either: Copy still works (copies last viewed file)
- ✅ OR: Clear error message shown
- ✅ No crash or hang

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

### Test 10: Stress Test - Rapid Copy Operations
**Priority**: P2 (Medium)

**Steps**:
1. Run: `./ccn docs/`
2. Open any file
3. **Rapidly press `y` 10 times** (as fast as possible)
4. Observe UI

**Expected Result**:
- ✅ UI remains responsive
- ✅ Status messages queue properly (or latest wins)
- ✅ No crash or hang
- ✅ Clipboard has valid content

**Pass/Fail**: ___________

**Notes**:
____________________________________________
____________________________________________

---

## 🔍 Visual Inspection Tests

### Test V1: Status Bar Message Visibility
**Check**: Status messages are readable and prominent

- [ ] Success message (`✓`) is clearly visible
- [ ] Error message (`✗`) stands out
- [ ] Message doesn't get cut off on narrow terminals
- [ ] Color/styling is appropriate for theme

**Pass/Fail**: ___________

---

### Test V2: Keyboard Hints Accuracy
**Check**: Status bar shows correct hints in each view

- [ ] FILE TREE view: Shows correct keys
- [ ] VIEWER view: Shows `y: copy` hint
- [ ] PREVIEW view: Shows correct keys
- [ ] Hints update when status message clears

**Pass/Fail**: ___________

---

## 🐛 Known Issues Check

### Issue Check 1: Selection Not Highlighted
**Status**: ⚠️ Known limitation (not yet implemented)

**Test**: Click and drag to select text

**Expected**:
- ⚠️ Selection currently NOT visually highlighted
- ✅ But copy still works
- ✅ Status message indicates what was copied

**This is OK for now** - Visual highlighting is Phase 2 enhancement

**Pass/Fail**: ___________

---

## 📊 Test Results Summary

### Overall Results

| Test | Priority | Status | Notes |
|------|----------|--------|-------|
| 1. Copy entire document | P0 | ☐ Pass ☐ Fail | |
| 2. Copy with selection | P0 | ☐ Pass ☐ Fail | |
| 3. Empty selection | P0 | ☐ Pass ☐ Fail | |
| 4. Error handling | P1 | ☐ Pass ☐ Fail | |
| 5. Large file performance | P1 | ☐ Pass ☐ Fail | |
| 6. Special characters | P1 | ☐ Pass ☐ Fail | |
| 7. Status auto-clear | P2 | ☐ Pass ☐ Fail | |
| 8. Multiple copies | P2 | ☐ Pass ☐ Fail | |
| 9. Copy then switch | P2 | ☐ Pass ☐ Fail | |
| 10. Stress test | P2 | ☐ Pass ☐ Fail | |

### Blocking Issues Found:
_______________________________________________________
_______________________________________________________
_______________________________________________________

### Non-Blocking Issues Found:
_______________________________________________________
_______________________________________________________
_______________________________________________________

---

## ✅ Sign-Off Criteria

**Ready to commit when**:
- [ ] All P0 tests pass
- [ ] At least 80% of P1 tests pass
- [ ] No blocking bugs found
- [ ] Performance targets met (<500ms for large files)
- [ ] Status messages work correctly

**Tester Name**: _______________________
**Date**: _______________________
**Overall Status**: ☐ PASS ☐ FAIL ☐ NEEDS WORK

---

## 🚀 Next Steps After Testing

### If Tests Pass:
1. ✅ Commit changes with test results
2. ✅ Update CHANGELOG.md
3. ✅ Push to repository
4. ✅ Mark user issue as resolved

### If Tests Fail:
1. ❌ Document failures in detail
2. ❌ Identify root cause
3. ❌ Fix issues
4. ❌ Re-run tests

---

## 📝 Testing Environment

**OS**: ____________________________
**Terminal**: ____________________________
**Go Version**: ____________________________
**Binary Size**: 15MB
**Build Date**: November 3, 2025

---

## 🆘 Troubleshooting

### If Copy Doesn't Work:
1. Check clipboard tool is installed:
   - macOS: Built-in (`pbcopy`)
   - Linux: Install `xclip` or `xsel`
   - Windows: Built-in

2. Check debug log:
   ```bash
   tail -f /tmp/lumina_mouse_debug.log
   ```

3. Test clipboard directly:
   ```bash
   echo "test" | pbcopy  # macOS
   echo "test" | xclip   # Linux
   ```

### If Status Messages Don't Appear:
1. Check terminal supports color
2. Check status bar is visible (not hidden by small terminal)
3. Check colorManager is initialized

---

**Last Updated**: 2025-11-03
**Version**: v1.0.2-alpha (with copy/paste fixes)
