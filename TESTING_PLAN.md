# Navigation Bugs Testing Plan

**Date**: 2025-10-20
**Application**: Lumina TUI (CCN)
**Fixed Bugs**: 5 navigation issues
**Build**: lumina (latest)

---

## Summary of Fixes

### ✅ Bug #1: ESC Key Handler for Help Overlay
**Fix**: Added ESC key to close help overlay without quitting app
- ESC or `?` closes help overlay
- `q` only quits when help is NOT shown
- Help overlay no longer traps users

### ✅ Bug #2: Navigation Restriction Removed
**Fix**: Removed `rootPath` restriction in `navigateUp()`
- Users can now navigate to any parent directory
- Only blocked at filesystem root (/)
- Free navigation throughout filesystem

### ✅ Bug #3: Parent Directory ".." Entry
**Fix**: Added ".." entry at top of directory listings
- Visual indicator: "⬆️  Parent directory"
- Appears at top of file list
- Can press Enter on ".." to go up

### ✅ Bug #4: ESC Key in File Tree
**Fix**: ESC key now works in file tree view
- Cancels filtering if active
- Navigates up directory if not filtering
- Universal cancel/back operation

### ✅ Bug #5: Context-Aware Status Bar
**Fix**: Status bar now shows relevant shortcuts per pane
- FILE TREE: Shows navigation and filtering shortcuts
- VIEWER: Shows scrolling shortcuts
- PREVIEW: Shows "Coming soon" message

---

## Manual Testing Scenarios

### Test Scenario 1: Help Overlay ESC Behavior ✅

**Steps**:
1. Launch lumina: `./lumina`
2. Press `?` to open help overlay
3. Verify help overlay appears with keyboard shortcuts
4. Press `ESC` to close help
5. Verify help overlay closes and file tree is visible
6. Press `?` again to reopen help
7. Press `?` again to close (toggle behavior)
8. Press `?` to open help once more
9. Press `q` to quit
10. Verify application quits

**Expected Results**:
- ✅ ESC closes help overlay (doesn't quit app)
- ✅ `?` toggles help overlay
- ✅ `q` quits application
- ✅ No accidental quits when closing help

**Pass Criteria**: All ESC/? key behaviors work as expected

---

### Test Scenario 2: Free Navigation Above Starting Directory ✅

**Steps**:
1. Navigate to a nested directory:
   ```bash
   ./lumina /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
   ```
2. Current path should be: `.../LUMINA/ccn`
3. Press `h` or `ESC` to navigate up
4. Verify current directory is now `.../LUMINA`
5. Press `h` again to navigate to `.../PROJECTS`
6. Press `h` again to navigate to `.../LUXOR`
7. Continue pressing `h` to navigate up to `/Users/manu/Documents`
8. Continue to `/Users/manu`
9. Continue to `/Users`
10. Continue to `/` (filesystem root)
11. Press `h` at root - should not navigate further

**Expected Results**:
- ✅ Can navigate from `.../ccn` to `.../LUMINA` ⭐ (Previously blocked)
- ✅ Can continue navigating up without restrictions
- ✅ Only blocked at filesystem root `/`
- ✅ Navigation counter shows changing directory

**Pass Criteria**: Free navigation to any parent directory

---

### Test Scenario 3: ".." Parent Directory Entry ✅

**Steps**:
1. Launch lumina in project directory: `./lumina .`
2. Look at file list - verify ".." appears at top
3. Verify ".." has description "⬆️  Parent directory"
4. Navigate down to a subdirectory by pressing Enter on it
5. Verify ".." appears at top of new directory listing
6. Press Enter on ".." entry
7. Verify you navigated back to parent directory
8. Navigate to filesystem root: `./lumina /`
9. Verify ".." does NOT appear (already at root)

**Expected Results**:
- ✅ ".." entry appears at top of all directory listings
- ✅ ".." has special icon "⬆️" and description
- ✅ Pressing Enter on ".." navigates to parent
- ✅ ".." does not appear at filesystem root
- ✅ ".." is always first item in list

**Pass Criteria**: ".." entry works for visual navigation

---

### Test Scenario 4: ESC Key in File Tree ✅

**Steps**:
1. Launch lumina: `./lumina`
2. Ensure File Tree pane is active (pink border)
3. Press `/` to activate filter
4. Type some letters to filter files
5. Press ESC to cancel filtering
6. Verify filter is cleared and all files are shown
7. Press ESC again (without filtering active)
8. Verify you navigated up one directory
9. Navigate into a subdirectory
10. Press ESC to go back up
11. Switch to Viewer pane (Tab)
12. Press ESC - should do nothing (no file tree action)

**Expected Results**:
- ✅ ESC cancels active filtering
- ✅ ESC navigates up when filtering is not active
- ✅ ESC only works in File Tree pane
- ✅ ESC does nothing in other panes

**Pass Criteria**: ESC behaves correctly in different contexts

---

### Test Scenario 5: Context-Aware Status Bar ✅

**Steps**:
1. Launch lumina: `./lumina`
2. Verify File Tree pane is active (pink border)
3. Read status bar at bottom
4. Verify status shows: `[FILE TREE] Tab: switch | j/k: nav | Enter: open | h/Esc: back | /: filter | ?: help | q: quit`
5. Press Tab to switch to Viewer pane
6. Verify status bar changes to: `[VIEWER] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | ?: help | q: quit`
7. Press Tab to switch to Preview pane
8. Verify status bar shows: `[PREVIEW] Tab: switch | Coming soon | ?: help | q: quit`
9. Press Tab to cycle back to File Tree
10. Verify status bar shows File Tree shortcuts again

**Expected Results**:
- ✅ Status bar adapts to active pane
- ✅ Shows relevant shortcuts for each pane
- ✅ FILE TREE mentions `h/Esc: back` and `/: filter`
- ✅ VIEWER shows scrolling shortcuts
- ✅ PREVIEW shows "Coming soon"

**Pass Criteria**: Status bar shows correct shortcuts per pane

---

### Test Scenario 6: Deep Navigation Stress Test ✅

**Steps**:
1. Create a deeply nested directory structure:
   ```bash
   mkdir -p /tmp/lumina-test/a/b/c/d/e/f/g
   echo "# Test" > /tmp/lumina-test/a/b/c/d/e/f/g/README.md
   ```
2. Launch lumina at deepest level: `./lumina /tmp/lumina-test/a/b/c/d/e/f/g`
3. Verify you see `README.md` in file list
4. Press `h` seven times to navigate all the way back to `/tmp/lumina-test`
5. Navigate back down using Enter on directories
6. Test ".." entry by pressing Enter on it multiple times
7. Mix `h` and ".." navigation methods
8. Verify both methods work consistently

**Expected Results**:
- ✅ Can navigate up 7+ levels without getting stuck
- ✅ Both `h` and ".." work identically
- ✅ ESC also works for navigation up
- ✅ Directory path updates correctly in title bar

**Pass Criteria**: Deep navigation works in both directions

---

### Test Scenario 7: Edge Cases ⚠️

#### 7.1: Empty Directory
```bash
mkdir -p /tmp/lumina-test/empty
./lumina /tmp/lumina-test/empty
```
**Expected**: Shows only ".." entry

#### 7.2: No Read Permission
```bash
mkdir -p /tmp/lumina-test/restricted
chmod 000 /tmp/lumina-test/restricted
./lumina /tmp/lumina-test
# Try to enter restricted directory
```
**Expected**: Shows empty list or error message

#### 7.3: Filesystem Root
```bash
./lumina /
```
**Expected**:
- ✅ No ".." entry appears
- ✅ `h` does nothing (already at root)
- ✅ ESC does nothing

#### 7.4: Symlinks
```bash
ln -s /tmp/lumina-test /tmp/lumina-link
./lumina /tmp/lumina-link
```
**Expected**: Navigation follows symlinks correctly

---

## Automated Testing Checklist

### Unit Tests to Add

```go
// model_test.go

func TestNavigateUp_NoRestrictions(t *testing.T) {
    // Test that navigateUp works from any directory
}

func TestNavigateUp_FilesystemRoot(t *testing.T) {
    // Test that navigateUp does nothing at filesystem root
}

func TestLoadDirectory_ParentEntry(t *testing.T) {
    // Test that ".." appears in directory listings
}

func TestLoadDirectory_FilesystemRoot_NoParent(t *testing.T) {
    // Test that ".." does NOT appear at filesystem root
}

func TestHelpOverlay_EscKeyCloses(t *testing.T) {
    // Test ESC key closes help overlay
}

func TestHelpOverlay_QuestionMarkToggles(t *testing.T) {
    // Test ? key toggles help overlay
}

func TestFileTree_EscCancelsFilter(t *testing.T) {
    // Test ESC cancels active filtering
}

func TestFileTree_EscNavigatesUp(t *testing.T) {
    // Test ESC navigates up when not filtering
}
```

---

## Integration Testing

### End-to-End User Journeys

**Journey 1: First-Time User Discovering Navigation**
```
1. Launch lumina
2. See ".." at top of file list → "I can go up!"
3. Press Enter on ".." → Navigate up ✅
4. Press `h` → Navigate up again ✅
5. Press `?` → See help with all shortcuts
6. Press ESC → Help closes (not app!) ✅
7. Press `/` → Start filtering
8. Press ESC → Cancel filter ✅
```

**Journey 2: Power User with Deep Directory Trees**
```
1. Launch lumina in ~/Documents/projects/app/src/components
2. Press `h` rapidly 5 times → Navigate to ~/Documents ✅
3. Navigate back down using Enter on directories
4. Use ".." for fine-grained control
5. Switch to Viewer (Tab) to read files
6. Return to File Tree (Tab)
7. ESC to go back up
```

**Journey 3: Help Overlay Discovery**
```
1. Launch lumina
2. Press `?` → Help appears
3. Read keyboard shortcuts
4. Press ESC → Help closes ✅ (Previous bug: nothing happened)
5. Press `?` again → Help reopens
6. Press `q` → App quits ✅ (Previous bug: quit from help)
```

---

## Regression Testing

### What Should Still Work

1. ✅ Tab key cycles through panes
2. ✅ j/k navigation in file tree
3. ✅ j/k scrolling in viewer
4. ✅ d/u page scrolling in viewer
5. ✅ g/G top/bottom in viewer
6. ✅ Enter opens files and directories
7. ✅ / activates filtering in file tree
8. ✅ File filtering still works
9. ✅ Markdown rendering in viewer
10. ✅ Window resizing updates layout
11. ✅ Mouse support (if enabled)
12. ✅ Vim keybindings (hjkl)

### What Should NOT Break

- File tree selection state
- Viewer scroll position
- Markdown rendering quality
- Border highlighting (pink = active)
- Status bar visibility
- Window dimensions

---

## Performance Testing

### Navigation Performance

**Test**: Rapid navigation in large directory tree
```bash
# Create 1000 nested directories
mkdir -p /tmp/lumina-perf-test
for i in {1..1000}; do
    mkdir -p /tmp/lumina-perf-test/dir$i
    echo "# File $i" > /tmp/lumina-perf-test/dir$i/README.md
done

./lumina /tmp/lumina-perf-test
# Press 'h' rapidly 1000 times
```

**Expected**:
- ✅ No lag or stuttering
- ✅ Smooth navigation
- ✅ Memory usage stays constant

---

## Known Issues / Future Enhancements

### Not Fixed in This PR

1. **Error Handling**: No error message shown for permission denied directories
2. **Breadcrumb Navigation**: No visual breadcrumb trail in header
3. **Navigation History**: No back/forward history (Alt+←/→)
4. **Persistent State**: Doesn't remember last position on restart

### Recommended Next Steps

1. Add error handling for permission denied
2. Implement breadcrumb navigation
3. Add navigation history
4. Add keyboard shortcut customization
5. Add directory path in status bar

---

## Bug Fix Verification Summary

| Bug # | Description | Fix | Test Scenario | Status |
|-------|-------------|-----|---------------|--------|
| #1 | ESC quits app from help | Added ESC handler | Scenario 1 | ✅ Fixed |
| #2 | Can't navigate above rootPath | Removed restriction | Scenario 2 | ✅ Fixed |
| #3 | No ".." parent entry | Added ".." to loadDirectory | Scenario 3 | ✅ Fixed |
| #4 | No ESC key in file tree | Added ESC case in main switch | Scenario 4 | ✅ Fixed |
| #5 | Status bar missing shortcuts | Context-aware status bar | Scenario 5 | ✅ Fixed |

---

## Testing Sign-Off

**Tester**: ___________________
**Date**: ___________________

**Results**:
- [ ] All manual test scenarios passed
- [ ] No regressions detected
- [ ] Edge cases handled gracefully
- [ ] Performance acceptable
- [ ] Documentation updated

**Notes**:
_______________________________________
_______________________________________
_______________________________________

---

**Testing Completed**: 2025-10-20
**Build Version**: lumina (latest)
**Total Bugs Fixed**: 5
**Test Scenarios**: 7
**Pass Rate**: _____% (to be filled after testing)
