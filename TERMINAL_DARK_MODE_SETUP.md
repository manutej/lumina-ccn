# Terminal Dark Mode Setup Guide

## Reduce Blue Light Pollution with Dark Terminal Theme

This guide shows how to configure various terminal emulators to use dark themes, which reduces blue light and is easier on your eyes during evening work sessions.

---

## macOS Terminals

### Terminal.app (Built-in)

1. **Open Terminal → Preferences** (Cmd + ,)
2. Go to **Profiles** tab
3. Select a dark profile or create custom:
   - **Pro**: Dark theme built-in, minimal
   - **Homebrew**: Green on black, hacker aesthetic
   - **Novel**: Off-white on dark, excellent for reading
4. Click **Default** button to make it default
5. Close and reopen Terminal for changes

**Recommended profiles for blue light reduction**:
- ✅ **Pro** (built-in) - Minimal colors, easy on eyes
- ✅ **Homebrew** - Low color intensity
- ✅ **Solarized Dark** (install separately)

### iTerm2 (Advanced, Recommended)

1. **Download**: https://iterm2.com
2. **Install** via Homebrew: `brew install --cask iterm2`
3. **Configure dark theme**:
   - Preferences → Profiles → Colors
   - Select **Solarized Dark** (built-in)
   - Or: Import custom color scheme
4. **Blue light reduction**:
   - Preferences → Profiles → Colors → Color Presets
   - Choose: **Solarized Dark** or **Dracula**

**Best iTerm2 color schemes for eyes**:
- ✅ **Solarized Dark** - Scientifically designed for readability
- ✅ **Dracula** - Dark with good contrast
- ✅ **Nord** - Cool blues, minimal eye strain
- ✅ **Gruvbox Dark** - Warm tones, gentle on eyes

**Install custom themes**:
```bash
# Clone iTerm2 color schemes
git clone https://github.com/mbadolato/iTerm2-Color-Schemes.git
cd iTerm2-Color-Schemes
# Import .itermcolors files via iTerm2 GUI
```

### VS Code Integrated Terminal

1. **Settings** (Cmd + ,)
2. Search: `terminal.integrated.theme`
3. Set to: `dark`
4. Alternatively, use VS Code theme settings:
   - **Command Palette** (Cmd + Shift + P)
   - Type: `Preferences: Color Theme`
   - Choose dark theme:
     - ✅ **Dracula Official** - Popular dark theme
     - ✅ **One Dark Pro** - Minimal, easy on eyes
     - ✅ **Nord** - Cool, reduced blue light

### tmux (Terminal Multiplexer)

Add to `~/.tmux.conf`:

```bash
# Enable 256 colors
set -g default-terminal "screen-256color"

# Dark background (forces dark mode)
set-option -g terminal-overrides ',xterm*:colors=256'

# Optional: Set specific colors
set -g status-bg black
set -g status-fg white
set -g window-status-current-bg white
set -g window-status-current-fg black
```

Reload: `tmux source-file ~/.tmux.conf`

---

## Linux Terminals

### GNOME Terminal

1. **Open Terminal → Preferences** (Ctrl + ,)
2. Go to **Profiles**
3. Create new profile or edit default:
   - Set **Text color**: Light gray
   - Set **Background**: Black
   - Disable: **Use colors from system theme**
4. Select as default

**Install dark theme**:
```bash
# Download Solarized Dark
git clone https://github.com/aruhier/gnome-terminal-colors-solarized.git
cd gnome-terminal-colors-solarized
./install.sh
```

### Konsole (KDE)

1. **Settings → Manage Profiles**
2. **Create New Profile** or edit default
3. **Appearance**: Select dark color scheme
4. **Color Schemes**:
   - ✅ **Solarized Dark Background**
   - ✅ **Dracula**
   - ✅ **Nord**

### xterm

Add to `~/.Xresources`:

```
xterm*background: #000000
xterm*foreground: #EEEEEE
xterm*cursorColor: #00FF00
xterm*colorBD: #FFFFFF
```

Reload: `xrdb -merge ~/.Xresources`

### Alacritty (GPU-Accelerated)

Edit `~/.config/alacritty/alacritty.yml`:

```yaml
colors:
  # Solarized Dark
  primary:
    background: '#002b36'
    foreground: '#839496'
  normal:
    black:   '#073642'
    red:     '#dc322f'
    green:   '#859900'
    yellow:  '#b58900'
    blue:    '#268bd2'
    magenta: '#2aa198'
    cyan:    '#2aa198'
    white:   '#eee8d5'
```

---

## Windows Terminals

### Windows Terminal (Recommended)

1. **Open Settings** (Ctrl + ,)
2. Go to **Appearance**
3. **Color scheme**: Select dark theme
   - ✅ **Dracula** (dark, reduces blue light)
   - ✅ **Nord** (cool dark theme)
   - ✅ **One Half Dark** (minimal)
4. **Theme**: Set to **Dark**

**Install custom scheme**:
```powershell
# Download Windows Terminal themes
git clone https://github.com/mbadolato/iTerm2-Color-Schemes
# Use color scheme JSON in Windows Terminal settings
```

### PowerShell (Command Prompt Alternative)

1. Right-click title bar → **Properties**
2. **Colors** tab
3. Set:
   - **Screen Background**: Black
   - **Screen Text**: Light Gray
   - **Popup Background**: Black
   - **Popup Text**: White

### ConEmu

1. **Settings → Colors**
2. Select color scheme:
   - ✅ **Solarized**
   - ✅ **Dracula**
3. Save as default

---

## General Blue Light Reduction Tips

### 1. **Enable Night Light / Flux**

**macOS**:
- System Preferences → Displays → Night Shift
- Set: "Sunset to Sunrise" or custom hours

**Windows 10/11**:
- Settings → System → Display → Night Light
- Enable and set schedule

**Linux**:
```bash
# Install redshift (blue light filter)
sudo apt install redshift redshift-gtk
redshift-gtk  # Run with GUI
```

### 2. **Color Scheme Recommendations for Blue Light Reduction**

| Scheme | Platform | Blue Light | Readability | Notes |
|--------|----------|-----------|------------|-------|
| **Solarized Dark** | All | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Scientific design |
| **Dracula** | All | ⭐⭐ | ⭐⭐⭐⭐ | Popular, good contrast |
| **Nord** | All | ⭐⭐⭐ | ⭐⭐⭐⭐ | Cool, gentle on eyes |
| **Gruvbox Dark** | Unix/Linux | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | Warm tones, best evening |
| **One Dark** | All | ⭐⭐ | ⭐⭐⭐⭐⭐ | High contrast |
| **Pro (Terminal)** | macOS | ⭐⭐⭐ | ⭐⭐⭐ | Minimal, clean |

### 3. **Font Settings for Comfort**

For any terminal, adjust:
- **Font Size**: 12-14pt (easier reading)
- **Line Height**: 1.3-1.5 (reduces eye strain)
- **Font**: Monospace with good hinting
  - ✅ **SF Mono** (macOS, clean)
  - ✅ **Fira Code** (excellent ligatures)
  - ✅ **JetBrains Mono** (designed for coding)
  - ✅ **Cascadia Code** (modern, clear)

### 4. **Terminal-Specific Blue Light Settings**

**Reduce cursor blink rate**:
```bash
# Most terminals: Settings → Cursor
# Set blink rate to slow or off
# Reduces additional light pollution
```

**Dim terminal background slightly** (if possible):
```bash
# macOS Terminal: Profiles → Colors → Opacity
# Set to 95-98% (slightly dimmed, feels easier on eyes)
```

---

## LUMINA Integration

Your LUMINA app now supports dark theme perfectly:

```bash
# LUMINA will auto-detect dark terminal
lumina
# → Shows dark theme with bright colors
# → Perfect for dark terminal setup
```

**Forced dark theme in LUMINA**:
```json
// ~/.config/lumina/colors.json
{
  "theme": "dark",
  "custom": {}
}
```

---

## Evening Work Setup (Recommended)

For maximum comfort during evening work:

1. **Terminal**: Solarized Dark or Gruvbox Dark
2. **System**: Enable Night Light (macOS/Windows) or redshift (Linux)
3. **Editor/Terminal Font**: 13pt, line height 1.4
4. **LUMINA**: Set to dark theme
5. **Screen brightness**: Reduce to 60-70%
6. **Room lighting**: Dim or use desk lamp

This combination significantly reduces blue light and eye strain during long evening coding sessions.

---

## Quick Setup Commands

### macOS (iTerm2 + Solarized Dark)
```bash
brew install --cask iterm2
# Then in iTerm2:
# Preferences → Profiles → Colors → Color Presets → Solarized Dark
```

### Ubuntu/Debian (GNOME Terminal + Solarized)
```bash
git clone https://github.com/aruhier/gnome-terminal-colors-solarized.git
cd gnome-terminal-colors-solarized
./install.sh -s dark
```

### Fedora (Konsole + Dark Theme)
```bash
sudo dnf install konsole
# Konsole → Settings → Manage Profiles → Color Scheme: Solarized Dark
```

---

## Verification

After setup, verify everything works:

```bash
# Test dark terminal with LUMINA
lumina

# You should see:
# ✅ Dark background in terminal
# ✅ Bright colors (teal borders, purple title) in LUMINA
# ✅ No eye strain or blue light glare
# ✅ Status bar text clearly visible
```

---

## Additional Resources

- **Solarized**: https://ethanschoonover.com/solarized/
- **Dracula**: https://draculatheme.com/
- **Nord**: https://www.nordtheme.com/
- **Gruvbox**: https://github.com/morhetz/gruvbox
- **Blue Light**: https://en.wikipedia.org/wiki/Blue_light_and_health

---

**Happy coding with reduced blue light! 🌙👀**
