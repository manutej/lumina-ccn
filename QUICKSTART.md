# Quick Start Guide - Lumina

Get started with Lumina (Claude Code Navigator) in 60 seconds!

## Installation Status ✅

**Lumina is already installed and ready to use!**

- ✅ Binary compiled at: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/ccn`
- ✅ Global command: `lumina` (symlinked to `~/bin/lumina`)
- ✅ Git repository initialized with `main` branch
- ✅ Phase 2 development branch: `phase-2-development`

## Quick Start

### 0. Check Installation (5 seconds)

```bash
# Verify lumina is installed
lumina --version
# Output: Lumina (CCN) Version: 1.0.0-alpha

# Show help
lumina --help

# Show keyboard shortcuts
lumina --keys
```

### 1. Run Lumina

```bash
# Navigate current directory
lumina

# Navigate specific directory
lumina ~/Documents/LUXOR/PROJECTS/LUMINA

# Navigate any markdown documentation
lumina /path/to/docs
```

### 2. Learn the Basics (30 seconds)

**General**:
- `?` - Show help overlay (keyboard shortcuts)
- `Tab` - Switch between panes
- `q` - Quit

**File Tree (Left Pane)**:
- `j` / `k` - Navigate up/down
- `Enter` - Open file or enter directory
- `Backspace` / `h` - Go to parent directory
- `/` - Filter files

**Viewer (Center Pane)**:
- `j` / `k` - Scroll up/down
- `d` / `u` - Half page down/up
- `g` / `G` - Jump to top/bottom

**Quick Tip**: Press `?` inside the app for full keyboard reference!

### 3. Try It Now!

```bash
# Navigate the LUMINA project
lumina ~/Documents/LUXOR/PROJECTS/LUMINA

# Once inside:
# 1. Press 'j' to navigate file list
# 2. Press 'Enter' on README.md
# 3. Press 'Tab' to switch to viewer
# 4. Press 'j'/'k' to scroll markdown
# 5. Press 'q' to quit
```

## Visual Guide

```
┌─────────────────────────────────────────────────────────────────┐
│ Claude Code Navigator (CCN) - LUMINA                            │
├─────────────────┬───────────────────────────────────────────────┤
│  📁 File Tree   │         📄 Glamour Rendered Markdown          │
│  ╭─────────────╮│  ╭───────────────────────────────────────╮   │
│  │►README.md   ││  │ # Lumina Project                       │   │
│  │ MOE-FINAL   ││  │                                        │   │
│  │ CHARM-ECO   ││  │ A next-generation markdown editor...   │   │
│  │ expert-*.md ││  │                                        │   │
│  ╰─────────────╯│  │ ## Features                            │   │
│                 │  │ - Beautiful rendering with Glamour     │   │
│  [FILE TREE]    │  │ - Vim-style keybindings                │   │
│                 │  ╰───────────────────────────────────────╯   │
│                 │                                               │
│                 │         📋 Preview (Phase 2)                  │
└─────────────────┴───────────────────────────────────────────────┘
│ Tab: switch | hjkl: nav | Enter: open | Backspace: up | q: quit│
└─────────────────────────────────────────────────────────────────┘
```

## Common Use Cases

### 1. Review Project Documentation

```bash
lumina ~/Documents/LUXOR/PROJECTS/my-project
```

Navigate through README, API docs, and markdown files with beautiful rendering.

### 2. Navigate Claude Code Workflows

```bash
lumina ~/Documents/LUXOR/PROJECTS
```

Browse multiple project documentation sets in one session.

### 3. Quick File Preview

```bash
# Open lumina in current directory
lumina

# Navigate to file with j/k
# Press Enter to view with Glamour rendering
# Press q to exit
```

## Keybinding Cheat Sheet

| Key | Action | Where |
|-----|--------|-------|
| `j` | Down | File Tree, Viewer |
| `k` | Up | File Tree, Viewer |
| `h` | Parent directory | File Tree |
| `Enter` | Open/Enter | File Tree |
| `Tab` | Next pane | All |
| `d` | Page down | Viewer |
| `u` | Page up | Viewer |
| `g` | Top | Viewer |
| `G` | Bottom | Viewer |
| `q` | Quit | All |
| `/` | Filter | File Tree |

## Tips & Tricks

### 1. Fast Navigation
- Use `Tab` to quickly switch between file tree and viewer
- Use `gg` / `G` in viewer to jump to top/bottom of documents

### 2. Filter Files
- Press `/` in file tree to filter by filename
- Start typing to narrow down results
- Clear filter with backspace

### 3. Vim Muscle Memory
- All vim navigation keys work: `hjkl`, `gg`, `G`, `d`, `u`
- More vim features coming in Phase 2!

## Next Steps

### Try These Features

1. **Navigate Multiple Files**
   ```bash
   lumina ~/Documents/LUXOR/PROJECTS/LUMINA
   # Open different .md files to see Glamour rendering
   ```

2. **Test Responsive Layout**
   ```bash
   lumina .
   # Resize your terminal - watch panes adapt!
   ```

3. **Vim Navigation**
   ```bash
   lumina ~/Documents/LUXOR/PROJECTS/LUMINA
   # Open README.md
   # Try: gg (top), G (bottom), d (page down), u (page up)
   ```

## Coming in Phase 2

🚀 **Exciting features on the way**:
- Fuzzy file finder (`/` to search all files)
- Ripgrep content search across files
- File watching with auto-reload
- Enhanced vim keybindings (marks, search)
- Split panes (horizontal/vertical)

## Troubleshooting

### Command Not Found

```bash
# Check symlink
ls -la ~/bin/lumina

# If missing, recreate:
ln -sf /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/ccn ~/bin/lumina

# Verify PATH includes ~/bin
echo $PATH | grep "$HOME/bin"
```

### Binary Not Working

```bash
# Rebuild
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
go build -o ccn

# Test locally
./ccn .
```

### Need Help?

- Check `README.md` for full documentation
- Check `PROJECT-STATUS.md` for implementation details
- Review `PHASE-1-COMPLETE.md` for features overview

## Feedback

This is an active development project. Phase 1 is complete and stable for daily use. Phase 2 development is on the `phase-2-development` branch.

---

**Enjoy exploring your documentation with Lumina!** 🌟

Run `lumina --help` to see available options (coming soon).
