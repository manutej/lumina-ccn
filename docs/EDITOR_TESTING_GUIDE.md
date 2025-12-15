# $EDITOR Integration Testing Guide

## Feature Overview

**Keybinding**: `e` (in Viewer pane)
**Purpose**: Open current file in external editor
**Commit**: `0a907fd feat(editor): Add $EDITOR integration (Phase 4)`

---

## Quick Test

```bash
# 1. Build and run LUMINA
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
go build -o ccn
./ccn

# 2. Navigate to a markdown file
# Use j/k to select, Enter to open

# 3. Press Tab to switch to Viewer pane

# 4. Press 'e' to open in editor
# Your $EDITOR will open (vim by default)

# 5. Make a change, save, and exit
# LUMINA will resume with file reloaded
```

---

## Test Cases

### Test 1: Default Editor (vim)

**Setup**: Ensure $EDITOR is not set or set to vim
```bash
unset EDITOR
./ccn
```

**Steps**:
1. Open any markdown file
2. Switch to Viewer pane (Tab)
3. Press `e`

**Expected**: vim opens with the file

---

### Test 2: Custom Editor (VS Code)

**Setup**:
```bash
export EDITOR="code --wait"
./ccn
```

**Steps**:
1. Open any markdown file
2. Switch to Viewer pane (Tab)
3. Press `e`

**Expected**: VS Code opens with the file, LUMINA waits until VS Code window closes

---

### Test 3: Custom Editor (neovim)

**Setup**:
```bash
export EDITOR="nvim"
./ccn
```

**Steps**:
1. Open any markdown file
2. Switch to Viewer pane (Tab)
3. Press `e`

**Expected**: neovim opens with the file

---

### Test 4: File Reload After Edit

**Steps**:
1. Open a markdown file in LUMINA
2. Note the content
3. Press `e` to open in editor
4. Add some text at the top
5. Save and exit editor

**Expected**:
- LUMINA resumes
- File shows updated content
- `[RELOADED]` indicator appears briefly
- Scroll position is preserved (or near original)

---

### Test 5: Scroll Position Preservation

**Steps**:
1. Open a long markdown file (100+ lines)
2. Scroll to middle (press `G` then `u` a few times)
3. Note your position
4. Press `e` to edit
5. Make no changes, just exit

**Expected**: Scroll position remains at same location

---

### Test 6: No File Selected

**Steps**:
1. Start LUMINA
2. Stay in File Tree pane (don't open a file)
3. Switch to Viewer pane (Tab)
4. Press `e`

**Expected**: Nothing happens (no crash, no action)

---

### Test 7: Edit JSON File

**Steps**:
1. Open a `.json` file in LUMINA
2. Switch to Viewer pane
3. Press `e`

**Expected**: Editor opens with JSON file

---

## Verification Commands

```bash
# Check your current $EDITOR
echo $EDITOR

# Set to vim
export EDITOR=vim

# Set to VS Code (must use --wait for TUI to resume properly)
export EDITOR="code --wait"

# Set to neovim
export EDITOR=nvim

# Set to nano
export EDITOR=nano
```

---

## Common Issues

### Issue: VS Code Opens But LUMINA Resumes Immediately

**Cause**: VS Code needs `--wait` flag to block until window closes

**Fix**:
```bash
export EDITOR="code --wait"
```

### Issue: Editor Not Found

**Cause**: Editor command not in PATH

**Fix**: Use full path
```bash
export EDITOR="/usr/local/bin/nvim"
```

### Issue: Terminal Garbled After Editor

**Cause**: Editor didn't restore terminal state properly

**Fix**: Run `reset` command or restart terminal

---

## Implementation Details

### Code Flow

```
User presses 'e'
       ↓
handleNormalMode() detects "edit" action
       ↓
openInEditorCmd(filePath) called
       ↓
tea.ExecProcess suspends TUI, runs editor
       ↓
Editor closes → EditorClosedMsg sent
       ↓
Update() handles EditorClosedMsg
       ↓
File reloaded, scroll position restored
```

### Key Files

| File | Function |
|------|----------|
| `main.go:138-148` | `openInEditorCmd()` implementation |
| `main.go:279-293` | `EditorClosedMsg` handler |
| `main.go:693-697` | `edit` action case |
| `keybindings.go:59` | `e` keybinding definition |

---

## Automated Tests

Currently no automated tests for editor integration (requires interactive terminal). Manual testing recommended.

Future: Could add mock tests for:
- `EditorClosedMsg` handling
- File reload after edit
- Scroll position preservation

---

**Last Updated**: 2025-12-15
**Version**: v1.4.2-alpha
**Linear Issue**: CET-348
