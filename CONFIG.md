# Lumina Configuration Guide

## Keybindings Configuration

### Location
- **macOS/Linux**: `~/.config/lumina/keybindings.json`
- **Windows**: `%APPDATA%\lumina\keybindings.json`

### Default Keybindings
Created automatically on first run. To customize, edit the JSON file.

### Example: One-Handed Navigation

If you prefer one-handed navigation (keeping your hand in one place), you can customize keybindings:

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

### Available Actions

#### Navigation
- `down` - Move down in file tree or scroll down in viewer
- `up` - Move up in file tree or scroll up in viewer
- `back` - Go to parent directory (file tree only)
- `open` - Open selected file/directory (file tree only)

#### Scrolling (Viewer only)
- `scroll_down` - Single line down
- `scroll_up` - Single line up
- `page_down` - Half page down
- `page_up` - Half page up
- `view_down` - Full page down (Alt+Down)
- `view_up` - Full page up (Alt+Up)
- `top` - Jump to top of file
- `bottom` - Jump to bottom of file

#### Actions
- `filter` - Toggle search/filter in file tree
- `help` - Show help overlay
- `switch_view` - Cycle through panes
- `copy` - Copy selected text (NEW!)

#### App Controls
- `quit` - Exit Lumina
- `copy` - Copy selected text to clipboard

### View Targets
- `filetree` - File tree pane only
- `viewer` - Markdown viewer pane only
- `any` - Any pane (global)

### Supported Key Names

**Single Keys**: `a-z`, `0-9`, space, enter, tab, esc, backspace, delete

**Modifiers**:
- `shift+{key}` - Shift modifier
- `ctrl+{key}` - Control modifier
- `alt+{key}` - Alt/Option modifier
- `cmd+{key}` - Command key (macOS)

**Special Keys**:
- `up`, `down`, `left`, `right` - Arrow keys
- `pageup`, `pagedown` - Page keys
- `home`, `end` - Home/End keys
- `f1-f12` - Function keys

---

## Copy/Selection Feature

### Keyboard Selection

**Copy with Keyboard**:
```
y                   # Copy selected text (must select first)
shift+j/k/h/l      # Select text with shifted arrow keys
shift+d/u/g/G      # Select larger text blocks
```

**Selection Modes**:
1. **Character mode** (default): Select individual characters
2. **Line mode**: `V` to select entire lines
3. **Block mode**: `Ctrl+V` to select rectangular blocks

### Mouse Selection

- **Click and drag** to select text
- **Double-click** to select word
- **Triple-click** to select paragraph

### Copy to Clipboard

Once text is selected:
- Press `y` to copy to system clipboard
- Use standard paste (`Ctrl+V` / `Cmd+V`) anywhere

---

## Color Customization

### Current Pane Indicator Colors

In future versions, edit `~/.config/lumina/colors.json`:

```json
{
  "inactive_border": "#874BFD",      // Purple (inactive)
  "active_border": "#00D084",        // Green (active) - IMPROVED CONTRAST
  "accent": "#7D56F4",               // Title color
  "text": "#FFFFFF",                 // Default text
  "highlight": "#FF79C6",            // Highlight color
  "error": "#FF5555"                 // Error messages
}
```

### Recommended Color Schemes

**High Contrast (Recommended)**:
- Inactive: `#555555` (dark gray)
- Active: `#00FF00` (bright green)

**Modern**:
- Inactive: `#444444` (charcoal)
- Active: `#00D084` (teal/green)
- Accent: `#A1EFD3` (light teal)

**Cyberpunk**:
- Inactive: `#0F0F0F` (black)
- Active: `#00FF00` (neon green)
- Accent: `#FF00FF` (magenta)

---

## First-Time Setup

1. Run `lumina` normally (creates default config)
2. Edit `~/.config/lumina/keybindings.json` for your preferences
3. Restart `lumina`

Changes take effect immediately on next startup.

---

## Troubleshooting

### Config file not found
Lumina creates a default config on first run. If missing:
```bash
lumina --create-config
```

### Keybindings not working
1. Check JSON syntax (use `jq` to validate)
2. Verify key names match supported keys (see list above)
3. Check view target ("filetree", "viewer", or "any")
4. Restart Lumina

### Copy not working
1. Ensure text is selected (should see highlighted text)
2. Press `y` to copy
3. Paste with system paste command (`Ctrl+V` / `Cmd+V`)

### Need system-level permissions?
Clipboard access is handled by the OS automatically on:
- macOS: Supported ✅
- Linux: Requires `xclip` or `xsel` (usually pre-installed)
- Windows: Supported ✅

Install missing clipboard tool:
```bash
# Ubuntu/Debian
sudo apt install xclip

# Fedora
sudo dnf install xclip

# macOS (via Homebrew)
brew install pbcopy pbpaste
```
