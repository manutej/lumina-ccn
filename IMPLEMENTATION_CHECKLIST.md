# Implementation Checklist: Custom Keybindings, Copy, & Colors

## 📋 What's New

Three powerful features ready to integrate:

1. **✅ Custom Keybindings** (`keybindings.go`) - One-handed navigation
2. **✅ Copy/Selection** (`clipboard.go`) - Text copy to clipboard
3. **✅ Better Colors** (UI update) - Clear pane distinction

---

## 🚀 Implementation Steps

### Step 1: Add Clipboard Dependency
```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
go get github.com/atotto/clipboard
go mod tidy
```
**Time**: 1 minute ⏱️

---

### Step 2: Update model.go
Add to `AppModel` struct (around line 45):
```go
keyBindings *KeyBindings
clipboard   *ClipboardManager
```

Add to `NewAppModel()` function:
```go
kb := LoadKeyBindings()
clipboard := NewClipboardManager()
```

Then in the return statement:
```go
return AppModel{
    // ... existing ...
    keyBindings: &kb,
    clipboard:   clipboard,
}
```

**Time**: 5 minutes ⏱️
**Files**: `model.go`

---

### Step 3: Replace Key Handler in main.go
Replace the entire switch statement (lines 52-154) with the action-based handler.

See `INTEGRATION_GUIDE.md` for exact code.

**Key changes**:
- Use `m.keyBindings.FindAction()` instead of hardcoded switch
- Map actions to behavior (scroll_down, page_up, copy, etc.)
- Clear selection when switching views
- Handle copy action

**Time**: 20 minutes ⏱️
**Files**: `main.go` (Update function)

---

### Step 4: Update View() Function for Colors
Replace pane color definitions (around lines 173-184):

**From**:
```go
paneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD"))      // Purple

activePaneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#FF79C6"))      // Pink
```

**To**:
```go
paneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#666666"))      // Dark gray

activePaneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#00D084")).     // Bright teal
    Bold(true)
```

**Time**: 5 minutes ⏱️
**Files**: `main.go` (View function)

---

### Step 5: Update Status Bar Text
Update status messages to show new actions (around line 238):

**Viewer status** (was):
```
[VIEWER] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | ?: help | q: quit
```

**Viewer status** (now):
```
[VIEWER] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | y: copy | ?: help | q: quit
```

**Time**: 2 minutes ⏱️
**Files**: `main.go` (View function)

---

### Step 6: Test Build
```bash
cd ccn
go build -o lumina
./lumina ~/.config  # Test it!
```

**Test checklist**:
- [ ] Builds without errors
- [ ] Help overlay shows (press ?)
- [ ] Navigation works
- [ ] Active pane color is clearly visible (bright teal)
- [ ] Inactive pane color is clearly visible (dark gray)

**Time**: 10 minutes ⏱️

---

### Step 7: Test Keybindings (Bonus)
Create `~/.config/lumina/keybindings.json`:

```json
{
  "navigation": [
    {"key": "j", "action": "down", "view": "filetree"},
    {"key": "down", "action": "down", "view": "filetree"},
    {"key": "k", "action": "up", "view": "filetree"},
    {"key": "up", "action": "up", "view": "filetree"},
    {"key": "h", "action": "back", "view": "filetree"},
    {"key": "esc", "action": "back", "view": "filetree"},
    {"key": "enter", "action": "open", "view": "filetree"}
  ],
  "scrolling": [
    {"key": "j", "action": "scroll_down", "view": "viewer"},
    {"key": "k", "action": "scroll_up", "view": "viewer"},
    {"key": "d", "action": "page_down", "view": "viewer"},
    {"key": "e", "action": "page_up", "view": "viewer"},
    {"key": "g", "action": "top", "view": "viewer"},
    {"key": "G", "action": "bottom", "view": "viewer"}
  ],
  "actions": [
    {"key": "/", "action": "filter", "view": "filetree"},
    {"key": "?", "action": "help", "view": "any"},
    {"key": "tab", "action": "switch_view", "view": "any"},
    {"key": "y", "action": "copy", "view": "viewer"}
  ],
  "app_controls": [
    {"key": "q", "action": "quit", "view": "any"},
    {"key": "ctrl+c", "action": "quit", "view": "any"}
  ]
}
```

**Test**:
```bash
# Restart lumina
lumina
# Try: d (down), e (up) - should work!
# Try: y (copy) - should copy selected text
```

**Time**: 5 minutes ⏱️

---

## 📊 Total Implementation Time

| Step | Time |
|------|------|
| 1. Add dependency | 1 min |
| 2. Update model.go | 5 min |
| 3. Replace key handler | 20 min |
| 4. Update colors | 5 min |
| 5. Status bar text | 2 min |
| 6. Test build | 10 min |
| 7. Test keybindings | 5 min |
| **TOTAL** | **48 minutes** |

✅ **Less than 1 hour to get all three features working!**

---

## 🧪 Testing Checklist

### Visual Appearance
- [ ] Active pane border is bright teal/green (clearly visible)
- [ ] Inactive pane border is dark gray (clearly distinguishable)
- [ ] Switching tabs (Tab key) shows clear visual change
- [ ] Help overlay still looks good

### Keybindings
- [ ] Default keybindings work (j/k, d/u, g/G)
- [ ] Config file is created at `~/.config/lumina/keybindings.json`
- [ ] Custom keybindings load from config
- [ ] One-handed example works (d=down, e=up)

### Copy Functionality
- [ ] Copy option appears in help text
- [ ] Pressing `y` copies selected text (or shows error if nothing selected)
- [ ] Clipboard integration works on macOS
- [ ] Clipboard integration works on Linux (with xclip)
- [ ] Clipboard integration works on Windows

### User Experience
- [ ] Clear indication which pane is active
- [ ] Can navigate with custom keys
- [ ] Can copy markdown content to paste elsewhere
- [ ] Help text reflects all changes

---

## 📚 Documentation Files

After implementation, users have access to:

1. **`CONFIG.md`** - How to customize keybindings
   - Config file location
   - Example configurations
   - All available actions and keys
   - Troubleshooting

2. **`FEATURES_ROADMAP.md`** - What's new and what's coming
   - Phase 1.5 additions
   - Planned Phase 2+ features
   - Future integration with Claude Code

3. **`INTEGRATION_GUIDE.md`** - Developer integration guide
   - Step-by-step code changes
   - Color scheme options
   - Technical details

---

## 🎯 What Users Get

### Immediately (After Implementation)
✅ Custom keybindings in `~/.config/lumina/keybindings.json`
✅ One-handed navigation support (d/e instead of d/u)
✅ Copy text to clipboard with `y` key
✅ Clear visual distinction between active/inactive panes
✅ Help text updated with new features

### Before Editing (Important!)
✅ Can select and copy markdown content
✅ No editing capabilities needed for this feature
✅ Works with existing Glamour rendering

---

## 🚀 Next Phase

Once implemented and tested:

1. **Phase 2 Features** (2-3 weeks)
   - Fuzzy finder for files
   - Ripgrep for content search
   - File watching for auto-reload

2. **Phase 3 Features** (following weeks)
   - Editor integration
   - Claude Code shortcuts
   - AI-powered agent features

---

## ⚠️ Known Issues / Considerations

### Clipboard on Linux
- Requires `xclip` or `xsel` to be installed
- Document in README for Linux users
- Auto-check if available and warn if missing

### Color Customization
- Currently hard-coded in view.go
- Future: Could add `~/.config/lumina/colors.json`
- For now: Users can edit source if they want different colors

### Selection Highlighting
- Current implementation stores selection state
- Visual highlighting could be added in Phase 2
- Infrastructure is ready for it

---

## 📝 Commit Messages

When committing these changes:

```bash
git add keybindings.go clipboard.go CONFIG.md INTEGRATION_GUIDE.md FEATURES_ROADMAP.md IMPLEMENTATION_CHECKLIST.md

git commit -m "feat: Add custom keybindings, copy functionality, and improved pane colors

- Implement configurable keybinding system via ~/.config/lumina/keybindings.json
- Add copy/selection capability with system clipboard integration
- Improve pane color distinction (bright teal active, dark gray inactive)
- Support one-handed navigation (customizable d/e keys)
- Add comprehensive documentation and guides

This addresses user requests for:
1. Custom keybinding support (keep hand in one place)
2. Copy capability (currently can't select/copy markdown)
3. Better visual indication of active pane

All features backward compatible, no breaking changes."
```

---

## 🎓 Learning Resources

- **Keybindings Pattern**: Studied from vim, emacs, and other TUIs
- **Clipboard Pattern**: Standard approach across all platforms
- **Color Contrast**: WCAG accessibility guidelines
- **Terminal Colors**: 256-color ANSI standard + hex support

---

## Final Notes

**This is one of the most valuable features you requested!**

The current code has:
- ❌ No selection capability (users can't copy anything!)
- ❌ Hardcoded keybindings (can't customize)
- ❌ Unclear pane indication (users get lost)

After implementation:
- ✅ Copy markdown content easily
- ✅ Customize keybindings exactly how you want
- ✅ Always know which pane is active

This positions LUMINA for the next phase while directly addressing real user pain points.

---

**Questions?** See the related documentation files:
- `CONFIG.md` - User configuration
- `INTEGRATION_GUIDE.md` - Technical integration details
- `FEATURES_ROADMAP.md` - Future plans and alternatives
