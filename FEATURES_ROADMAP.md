# Lumina Features Roadmap

## Current Status (v1.0.0-alpha)
✅ 3-pane layout (File Tree | Viewer | Preview)
✅ Glamour markdown rendering
✅ Vim-style keybindings
✅ Help overlay
✅ Global command (lumina)

---

## NEW: Phase 1.5 - Quality & UX (THIS WEEK)

### Completed ✅
- [x] Navigation bug fixes
- [x] Option+arrow keybindings
- [x] ESC key behavior fixed
- [x] Parent directory navigation (..)

### In Progress 🔄

#### 1. **Custom Keybindings** [NEW FILES]
- ✅ `keybindings.go` - Configuration system
- ✅ `CONFIG.md` - User guide
- ✅ `INTEGRATION_GUIDE.md` - Dev guide

**Enables**:
- One-handed navigation (d/e instead of d/u)
- Custom key combinations
- User preferences saved in config file
- Auto-generated default config

**User Story**: "I want to use `d` for down and `e` for up, keeping my hand in one place"

**Effort**: 4 hours integration + testing

---

#### 2. **Copy/Selection Capability** [NEW FILES]
- ✅ `clipboard.go` - Selection and copy system
- ✅ Support for system clipboard (macOS, Linux, Windows)
- ✅ Keyboard and mouse selection modes

**Enables**:
- Select text with keyboard (Shift+hjkl or arrow keys)
- Select text with mouse (click & drag)
- Copy to system clipboard with `y` key
- Paste anywhere with system paste (Ctrl+V / Cmd+V)

**User Story**: "I need to copy text from markdown files - currently I can't select anything!"

**Effort**: 6 hours integration + testing

---

#### 3. **Better Visual Distinction** [QUICK WIN]
- Pane border colors improved
- Active pane: `#00D084` (bright teal) instead of `#FF79C6` (pink)
- Inactive pane: `#666666` (dark gray) instead of `#874BFD` (purple)
- Optional: Add bold styling to active pane

**User Story**: "When I switch tabs, I can't clearly see which pane is active"

**Effort**: 1 hour change

---

### Summary: Phase 1.5 Additions
| Feature | Effort | Dependency | Priority |
|---------|--------|------------|----------|
| Custom Keybindings | 4h | None | HIGH |
| Copy/Selection | 6h | clipboard lib | CRITICAL |
| Better Colors | 1h | None | HIGH |
| **Total** | **11h** | - | - |

**Timeline**: 1-2 weeks with proper testing

---

## Phase 2 - Core Features (Planned)

### Fuzzy Finder (3-5 days)
- Real-time fuzzy search across file tree
- Telescope-style overlay
- Jump to file instantly
- Related: `sahilm/fuzzy` package

### Ripgrep Integration (3-5 days)
- Content search across all files
- Results with context/line numbers
- Jump to match in viewer
- Performance optimized

### File Watching (2-3 days)
- Auto-reload when files change
- Don't lose scroll position
- Visual indicator for updated files
- fsnotify integration (already in go.mod)

### Enhanced Vim Keybindings (2 days)
- Search within document (/)
- Next/previous match (n/N)
- Bookmarks (ma to mark, 'a to jump)
- Visual mode selections

---

## Future: Phase 3+ - Workflow Integration

### Editor Integration
- Open files in $EDITOR
- Live preview while editing
- Auto-reload detection
- Git integration

### Claude Code Integration
- `/moe` command shortcuts
- `/workflows` quick access
- Session bookmarks
- MCP server awareness

### Agentic Features (Using anthropic-sdk-go)
- AI-powered suggestions
- Code snippet explanations
- Intelligent navigation
- Crush-inspired AI coding agent

---

## Library Recommendations

### ✅ Already Using
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/glamour` - Markdown rendering
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/fsnotify/fsnotify` - File watching

### 📦 Recommended for Future
- `github.com/charmbracelet/anthropic-sdk-go` - Claude SDK (Phase 3+)
- `github.com/charmbracelet/crush` - Code agent reference
- `github.com/atotto/clipboard` - Clipboard access [NEW]
- `github.com/sahilm/fuzzy` - Fuzzy search (Phase 2)

---

## Color Scheme Options

Currently using:
- Inactive: `#666666` (dark gray)
- Active: `#00D084` (bright teal)
- Title: `#7D56F4` (purple)

### Alternative Themes

**High Contrast (Recommended)**
```json
{
  "inactive": "#555555",
  "active": "#00FF00",
  "accent": "#7D56F4"
}
```

**Cyberpunk**
```json
{
  "inactive": "#0F0F0F",
  "active": "#00FF00",
  "accent": "#FF00FF"
}
```

**Monokai**
```json
{
  "inactive": "#3E3D32",
  "active": "#A1EFD3",
  "accent": "#F92672"
}
```

---

## Testing Checklist

### Phase 1.5 Testing
- [ ] Custom keybindings load from config
- [ ] One-handed navigation works (d/e example)
- [ ] Copy functionality works on all platforms
- [ ] Selected text highlights visually
- [ ] Color change obvious in all terminals
- [ ] Help text updated
- [ ] Status bar reflects new actions

### Platform Testing
- [ ] macOS - Copy works with clipboard
- [ ] Linux - Requires xclip/xsel (document requirement)
- [ ] Windows - Copy works via clipboard

---

## Success Metrics

### v1.0.1 (After Phase 1.5)
- [ ] 40% test coverage (from 0%)
- [ ] All P0 bugs fixed
- [ ] Custom keybindings working
- [ ] Copy functionality working
- [ ] Clear visual distinction between panes

### v1.0.2 (After Phase 2)
- [ ] Fuzzy finder implemented
- [ ] Ripgrep integration done
- [ ] File watching active
- [ ] Enhanced vim keybindings

### v1.1.0 (After Phase 3)
- [ ] Editor integration
- [ ] Claude Code shortcuts
- [ ] agentic AI features
- [ ] Production-ready

---

## Quick Reference: New Features

### One-Handed Navigation
```bash
# Edit: ~/.config/lumina/keybindings.json
# Change "d" action to "page_down"
# Change "e" action to "page_up"
# Then: d = down half-page, e = up half-page
```

### Copy Text
```bash
# Select text (keyboard: Shift+hjkl or arrow keys)
# Press: y
# System clipboard updated!
# Paste anywhere: Ctrl+V (Linux/Windows) or Cmd+V (macOS)
```

### Better Colors
```bash
# Automatic on upgrade
# Active pane border: Bright teal (#00D084)
# Inactive pane border: Dark gray (#666666)
# Much clearer visual distinction!
```

---

## Resources

- **Config Location**: `~/.config/lumina/keybindings.json`
- **Charm Ecosystem**: https://charm.sh
- **Bubble Tea Docs**: https://github.com/charmbracelet/bubbletea
- **Glamour Styling**: https://github.com/charmbracelet/glamour
- **Anthropic SDK**: https://github.com/charmbracelet/anthropic-sdk-go

---

## Questions?

See:
- `CONFIG.md` - Configuration guide for users
- `INTEGRATION_GUIDE.md` - Integration guide for developers
- `ccn/README.md` - General project overview
