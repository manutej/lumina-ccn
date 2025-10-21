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

## Color Customization & Theme Management

### Overview

Lumina now supports **automatic theme detection** that intelligently switches between dark and light color schemes based on your terminal background. The theme system includes:

- **Auto-detection** (default) - Automatically detects light/dark terminal backgrounds
- **Manual themes** - Dark, Light, Dracula, Tokyo Night, and more
- **Environment variable override** - Set `GLAMOUR_STYLE` for markdown rendering
- **Configuration file** - Customize colors in `~/.config/lumina/colors.json`

### Theme Configuration

#### Location
- **macOS/Linux**: `~/.config/lumina/colors.json`
- **Windows**: `%APPDATA%\lumina\colors.json`

#### Default Configuration

The default configuration uses auto-detection:

```json
{
  "theme": "auto",
  "custom": {}
}
```

**Theme Options**:
- `"auto"` - Auto-detect based on terminal background (recommended)
- `"dark"` - Force dark theme
- `"light"` - Force light theme (optimized for light terminal backgrounds)

#### Light Theme Features

When running in a light terminal background, Lumina automatically switches to a light-optimized palette:

**Light Theme Colors**:
```
Title:                Dark Purple (#5A4A8A)
Active Pane Border:   Dark Teal (#007A5E)
Inactive Border:      Light Gray (#CCCCCC)
Status Bar:           Dark Text (#333333)
Help Command Text:    Dark Blue (#0066CC)
Help Border:          Dark Pink (#C1426B)
Error Text:           Dark Red (#CC0000)
Success Text:         Dark Green (#007A5E)
```

**Dark Theme Colors** (preserved for backwards compatibility):
```
Title:                Purple (#7D56F4)
Active Pane Border:   Bright Teal (#00D084)
Inactive Border:      Dark Gray (#666666)
Status Bar:           Dark Gray (#666666)
Help Command Text:    Cyan (#8BE9FD)
Help Border:          Pink (#FF79C6)
Error Text:           Red (#FF5555)
Success Text:         Green (#00D084)
```

### Force a Specific Theme

To force Lumina to always use a specific theme, edit `~/.config/lumina/colors.json`:

```json
{
  "theme": "dark",
  "custom": {}
}
```

Options:
- `"dark"` - Always use dark theme
- `"light"` - Always use light theme (for light terminal backgrounds)
- `"auto"` - Let Lumina detect the best theme

### Environment Variable Override for Markdown

For fine-grained control over markdown rendering style (independent of UI theme), set the `GLAMOUR_STYLE` environment variable:

```bash
# Use light markdown style with dark UI
GLAMOUR_STYLE=light lumina

# Use dark markdown style (default)
GLAMOUR_STYLE=dark lumina

# Use Dracula markdown theme
GLAMOUR_STYLE=dracula lumina

# Use Tokyo Night markdown theme
GLAMOUR_STYLE=tokyo-night lumina
```

Available markdown styles:
- `auto` - Auto-detect (Glamour's built-in detection)
- `dark` - Dark theme
- `light` - Light theme
- `dracula` - Dracula theme
- `tokyo-night` - Tokyo Night theme
- `pink` - Pink theme
- `notty` - No terminal colors

### Real-World Examples

#### Light Terminal with Light Content

If you use a light terminal background and want everything optimized for light:

```bash
# Terminal settings: light background
# Then run:
lumina
# Lumina will auto-detect and use light theme ✅
```

No configuration needed! Lumina handles it automatically.

#### Dark Terminal, Light Markdown

If you prefer dark UI but light markdown rendering:

```bash
GLAMOUR_STYLE=light lumina
# UI stays dark, markdown uses light colors
```

#### Forcing Dark Mode (Dark Terminal)

```json
{
  "theme": "dark",
  "custom": {}
}
```

Then run: `lumina`

#### Forcing Light Mode (Any Terminal)

```json
{
  "theme": "light",
  "custom": {}
}
```

Then run: `lumina`

### Auto-Detection Details

Lumina uses multiple methods to detect terminal background:

1. **COLORFGBG environment variable** - Set by most terminal emulators
   - Light backgrounds: 7, 15
   - Dark backgrounds: 0, 8

2. **Terminal capabilities** - Queries terminal for background color support

3. **Fallback** - Defaults to dark if detection fails

### Custom Colors (Future Enhancement)

The `custom` field in `colors.json` is reserved for user-defined color overrides (coming in Phase 2):

```json
{
  "theme": "light",
  "custom": {
    "title": "#FF6B6B",
    "active_border": "#4ECDC4"
  }
}
```

### Recommended Color Schemes

**For Light Terminal Backgrounds**:
- Use theme: `"auto"` (recommended)
- Or force: `"light"`
- Status bar becomes dark, highly readable on light backgrounds

**For Dark Terminal Backgrounds**:
- Use theme: `"auto"` (recommended)
- Or force: `"dark"`
- Bright colors pop against dark backgrounds

**Mixed Environment** (laptop + monitor):
- Use `"auto"` - Lumina adapts when starting in different terminals
- Create symlink to same config: `ln -s ~/.config/lumina/colors.json ~/.config/lumina/.backup`

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
