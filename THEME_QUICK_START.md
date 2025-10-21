# LUMINA Theme Quick Start Guide

## TL;DR - Just Works! 🎉

Lumina now automatically detects your terminal background and switches themes:

- **Dark terminal** → Dark theme (your existing colors) ✅
- **Light terminal** → Light theme (new! optimized colors) ✅

**No configuration needed!** Just run:

```bash
lumina
```

---

## What Changed?

Previously, Lumina always used dark mode colors, making the status bar text invisible on light terminal backgrounds.

Now, Lumina intelligently adapts:

| Aspect | Dark Terminal | Light Terminal |
|--------|---------------|----------------|
| **Status Bar** | Dark gray text | Dark text (#333333) ✅ |
| **Title** | Purple | Dark purple ✅ |
| **Borders (active)** | Bright teal | Dark teal ✅ |
| **Borders (inactive)** | Dark gray | Light gray ✅ |
| **Readability** | ✅ Good | ✅ Good (NEW!) |

---

## Three Ways to Control Your Theme

### 1. Auto-Detect (Default - Recommended)
```bash
# Lumina detects your terminal background automatically
lumina
```

**Best for**: Most users, laptop + external monitor, working in different terminals

### 2. Force Dark Theme
Edit `~/.config/lumina/colors.json`:
```json
{
  "theme": "dark",
  "custom": {}
}
```

**Best for**: Prefer dark mode even on light terminals

### 3. Force Light Theme
Edit `~/.config/lumina/colors.json`:
```json
{
  "theme": "light",
  "custom": {}
}
```

**Best for**: Always work in light terminals

---

## Markdown Rendering Options

Want to control markdown colors separately from the UI?

```bash
# Light markdown, dark UI
GLAMOUR_STYLE=light lumina

# Dark markdown (default)
GLAMOUR_STYLE=dark lumina

# Dracula markdown theme
GLAMOUR_STYLE=dracula lumina

# Tokyo Night markdown theme
GLAMOUR_STYLE=tokyo-night lumina
```

Other options: `pink`, `notty`

---

## Configuration File Location

- **macOS/Linux**: `~/.config/lumina/colors.json`
- **Windows**: `%APPDATA%\lumina\colors.json`

Created automatically on first run!

---

## How Auto-Detection Works

Lumina checks:

1. **Terminal environment variable** (`COLORFGBG`)
   - Set by most terminal emulators
   - Values: 7, 15 = light | 0, 8 = dark

2. **Terminal program** (`TERM_PROGRAM`)
   - iTerm defaults to dark
   - Others default to dark (safe fallback)

3. **Fallback**: Defaults to dark if detection fails

---

## Testing It

### See Dark Theme
```bash
# Terminal settings: dark background
lumina
# Should show bright cyan borders, pink help text
```

### See Light Theme
```bash
# Terminal settings: light background
lumina
# Should show dark teal borders, dark blue help text
# Status bar text should be dark and readable ✅
```

### Force Dark Theme
```bash
# Terminal settings: light background
# Edit ~/.config/lumina/colors.json: "theme": "dark"
lumina
# Stays dark even on light background
```

---

## Color Palettes

### Dark Theme (Original)
- **Title**: Purple (#7D56F4)
- **Active Border**: Bright Teal (#00D084)
- **Inactive Border**: Dark Gray (#666666)
- **Status Bar**: Dark Gray (#666666)

### Light Theme (New)
- **Title**: Dark Purple (#5A4A8A)
- **Active Border**: Dark Teal (#007A5E)
- **Inactive Border**: Light Gray (#CCCCCC)
- **Status Bar**: Dark Text (#333333)

---

## Troubleshooting

### Colors not changing?

1. Check config file: `cat ~/.config/lumina/colors.json`
2. Try forcing a theme (edit JSON)
3. Restart Lumina (changes take effect on startup)

### Can't see status bar text?

This should be fixed! If it's still hard to read:
- Check your terminal colors are working
- Try forcing light theme: `"theme": "light"`
- Report the issue with your terminal name

### Want different colors?

Phase 2 will add custom color overrides! For now, you can:
- Use `GLAMOUR_STYLE` for markdown
- Choose between dark/light themes
- Report color issues

---

## FAQ

**Q: Do I need to do anything?**
A: Nope! Lumina detects your terminal automatically.

**Q: Will my dark terminal setup still work?**
A: Yes! Dark theme is preserved. Auto-detection respects dark terminals.

**Q: Can I use custom colors?**
A: Coming in Phase 2! For now, choose dark or light.

**Q: How do I report color issues?**
A: Include: terminal name, background color, theme used.

**Q: What if auto-detection doesn't work?**
A: Manually set theme in `~/.config/lumina/colors.json`

---

## Next in Phase 2

- 🎨 Custom color picker
- 🔄 Theme selector in-app
- 📚 More built-in themes
- 💾 Color scheme import/export

---

**Ready to use?** Just run:

```bash
lumina
```

Happy reading! 📚✨
