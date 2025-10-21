# How to Change LUMINA & Terminal Theme

## Summary

Your LUMINA app now has **dark mode forced as default**. Here's how to control themes:

---

## 🎨 Changing LUMINA Theme

### Current State
Dark theme is forced: `~/.config/lumina/colors.json`
```json
{
  "theme": "dark",
  "custom": {}
}
```

### To Change LUMINA Theme

**Step 1: Edit configuration file**
```bash
# Open the config file
nano ~/.config/lumina/colors.json
```

**Step 2: Choose a theme**

Option A - Force Dark (current):
```json
{
  "theme": "dark",
  "custom": {}
}
```

Option B - Force Light (for light terminals):
```json
{
  "theme": "light",
  "custom": {}
}
```

Option C - Auto-detect (intelligent):
```json
{
  "theme": "auto",
  "custom": {}
}
```

**Step 3: Save and restart**
```bash
# Save file (Ctrl+X, Y, Enter in nano)
# Then restart LUMINA:
lumina
```

---

## 🖥️ Changing Your Terminal Theme (Reduce Blue Light)

Your terminal itself can have a dark theme. This is the main way to reduce blue light pollution.

### macOS

**Using Terminal.app (built-in):**
```bash
# Open Terminal → Preferences (Cmd + ,)
# Go to Profiles → Select "Pro" (dark)
# Click Default
```

**Using iTerm2 (recommended for dark mode):**
```bash
# Install: brew install --cask iterm2

# In iTerm2:
# Preferences → Profiles → Colors
# Color Presets → Select "Solarized Dark"

# Or install more themes:
git clone https://github.com/mbadolato/iTerm2-Color-Schemes
# Import .itermcolors files via iTerm2 GUI
```

### Ubuntu/Debian

**GNOME Terminal with Solarized Dark:**
```bash
# Install Solarized theme
git clone https://github.com/aruhier/gnome-terminal-colors-solarized.git
cd gnome-terminal-colors-solarized
./install.sh -s dark

# Then in GNOME Terminal:
# Preferences → Profiles → Edit → Colors
# Choose "Solarized Dark"
```

### Fedora/RedHat

**Konsole (KDE Terminal):**
```bash
# Settings → Manage Profiles
# Color scheme → Select "Solarized Dark Background"
# Or "Dracula"
```

### Windows

**Windows Terminal:**
```powershell
# Open Settings (Ctrl + ,)
# Appearance → Color scheme
# Select: "Dracula" or "Nord"
# Or set Theme to "Dark"
```

---

## 🌙 Blue Light Reduction Setup

For comfortable evening coding:

### 1. Terminal Dark Theme
✓ Use Solarized Dark or Dracula
✓ Reduces blue light from terminal
✓ Easier on eyes during long sessions

### 2. System Night Light
- **macOS**: System Preferences → Displays → Night Shift (Enable)
- **Windows**: Settings → System → Display → Night Light (Enable)
- **Linux**: `sudo apt install redshift` then `redshift-gtk`

### 3. LUMINA Configuration
✓ Dark theme (already forced)
✓ All UI colors optimized for dark terminal

### 4. Display Settings
- Reduce brightness to 60-70%
- Use warm desk lighting
- Font size: 13pt
- Line height: 1.4

---

## 📋 Quick Reference

### LUMINA Theme Options

| Theme | Config Value | Use Case |
|-------|------------|----------|
| **Dark** | `"dark"` | Dark terminal (current, recommended) |
| **Light** | `"light"` | Light terminal, light background |
| **Auto** | `"auto"` | Intelligent detection (future) |

### Color Palettes

**Dark Theme:**
```
Title: Purple (#7D56F4)
Active Border: Bright Teal (#00D084)
Status Bar: Dark Gray (#666666)
```

**Light Theme:**
```
Title: Dark Purple (#5A4A8A)
Active Border: Dark Teal (#007A5E)
Status Bar: Dark Text (#333333)
```

### Terminal Color Schemes

| Scheme | Platform | Blue Light | Notes |
|--------|----------|-----------|-------|
| **Solarized Dark** | All | ⭐⭐⭐ | Scientifically designed |
| **Dracula** | All | ⭐⭐ | Popular, good contrast |
| **Nord** | All | ⭐⭐⭐ | Cool, gentle colors |
| **Gruvbox** | Unix | ⭐⭐⭐⭐ | Warm, best for evening |

---

## 🔍 Verify Changes

**After changing terminal theme:**
```bash
# Run LUMINA
lumina

# You should see:
# ✓ Terminal has dark background
# ✓ LUMINA shows dark theme colors
# ✓ No eye strain or glare
# ✓ Text is clearly readable
```

**Run verification test:**
```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
./test-dark-theme.sh
```

Expected result: ✅ All 14 tests PASS

---

## 📚 Additional Resources

- **LUMINA Docs**: `CONFIG.md`, `THEME_QUICK_START.md`
- **Terminal Setup**: `TERMINAL_DARK_MODE_SETUP.md`
- **Solarized**: https://ethanschoonover.com/solarized/
- **Dracula**: https://draculatheme.com/
- **Blue Light Info**: https://en.wikipedia.org/wiki/Blue_light_and_health

---

## Troubleshooting

**Q: Changes not applying?**
A: Restart LUMINA after editing config file.

**Q: Can't find config file?**
A: Location: `~/.config/lumina/colors.json`
   Create if missing: `mkdir -p ~/.config/lumina`

**Q: Terminal doesn't support color schemes?**
A: Most modern terminals do. Update your terminal or try iTerm2/Windows Terminal.

**Q: Status bar text not visible?**
A: Use light theme if terminal background is light:
   Change config to: `"theme": "light"`

---

## Summary

✅ **Dark theme**: Forced by default in LUMINA
✅ **Terminal dark theme**: Set your terminal to Solarized Dark or Dracula
✅ **Blue light reduction**: Enabled via Night Light + dark colors
✅ **Ready to use**: Just run `lumina`

**Happy coding with reduced blue light! 🌙👀**
