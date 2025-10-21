# Integration Guide: Custom Keybindings, Copy, and Better Colors

## Quick Summary

Three new files have been created:
1. **`keybindings.go`** - Configurable keybinding system
2. **`clipboard.go`** - Copy/selection functionality
3. **`CONFIG.md`** - User configuration guide

These integrate into the existing code with minimal changes.

---

## Step 1: Add Clipboard Dependency

```bash
cd ccn
go get github.com/atotto/clipboard
```

This provides system clipboard access across all platforms.

---

## Step 2: Initialize in NewAppModel()

In `model.go`, add to `AppModel` struct:

```go
type AppModel struct {
	// ... existing fields ...

	// NEW: Configuration and clipboard
	keyBindings *KeyBindings
	clipboard   *ClipboardManager
}
```

Then in `NewAppModel()`:

```go
func NewAppModel(rootPath string) AppModel {
	// ... existing code ...

	// NEW: Load keybindings and initialize clipboard
	kb := LoadKeyBindings()
	clipboard := NewClipboardManager()

	return AppModel{
		// ... existing fields ...
		keyBindings: &kb,
		clipboard:   clipboard,
	}
}
```

---

## Step 3: Update Key Handling in main.go

**Replace the hardcoded switch statement** (lines 52-154) with this action-based approach:

```go
case tea.KeyMsg:
	// Handle help overlay toggle
	if msg.String() == "?" {
		m.showHelp = !m.showHelp
		return m, nil
	}

	// If help is shown, ESC or ? closes it
	if m.showHelp {
		if msg.String() == "esc" || msg.String() == "?" {
			m.showHelp = false
			return m, nil
		}
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	}

	// NEW: Use configurable keybindings
	viewName := []string{"filetree", "viewer", "preview"}[m.currentView]
	action := m.keyBindings.FindAction(msg.String(), viewName)

	// If no view-specific action, try "any" view
	if action == "" {
		action = m.keyBindings.FindAction(msg.String(), "any")
	}

	// Execute action
	switch action {
	// App controls
	case "quit":
		return m, tea.Quit

	// Navigation
	case "back":
		if m.currentView == FileTreeView {
			if m.fileList.FilterState() == list.Filtering {
				m.fileList.ResetFilter()
			} else {
				m.navigateUp()
			}
		}

	case "open":
		if m.currentView == FileTreeView {
			m.navigateToSelectedFile()
		}

	case "switch_view":
		m.currentView = (m.currentView + 1) % 3
		m.clipboard.ClearSelection() // Clear selection when switching views

	// File tree navigation
	case "down":
		if m.currentView == FileTreeView {
			m.fileList, cmd = m.fileList.Update(msg)
		}
		return m, cmd

	case "up":
		if m.currentView == FileTreeView {
			m.fileList, cmd = m.fileList.Update(msg)
		}
		return m, cmd

	// Viewer scrolling
	case "scroll_down":
		if m.currentView == ViewerView {
			m.viewer.LineDown(1)
		}

	case "scroll_up":
		if m.currentView == ViewerView {
			m.viewer.LineUp(1)
		}

	case "page_down":
		if m.currentView == ViewerView {
			m.viewer.HalfViewDown()
		}

	case "page_up":
		if m.currentView == ViewerView {
			m.viewer.HalfViewUp()
		}

	case "view_down":
		if m.currentView == ViewerView {
			m.viewer.ViewDown()
		}

	case "view_up":
		if m.currentView == ViewerView {
			m.viewer.ViewUp()
		}

	case "top":
		if m.currentView == ViewerView {
			m.viewer.GotoTop()
		}

	case "bottom":
		if m.currentView == ViewerView {
			m.viewer.GotoBottom()
		}

	// NEW: Copy functionality
	case "copy":
		if m.currentView == ViewerView {
			if err := m.clipboard.CopySelection(m.viewerContent); err == nil {
				// Copied! Could show status message
				// For now, selection clears
				m.clipboard.ClearSelection()
			}
		}

	// Filter
	case "filter":
		if m.currentView == FileTreeView {
			// Toggle filter - let bubbles handle it
			if m.fileList.FilterState() != list.Filtering {
				m.fileList.SetFilteringEnabled(true)
			}
		}

	case "help":
		m.showHelp = !m.showHelp
	}

	return m, tea.Batch(cmds...)
```

---

## Step 4: Improve Pane Color Distinction

In the `View()` function, replace the color definitions (around line 173):

**Current (low contrast)**:
```go
paneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD"))      // Purple

activePaneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#FF79C6"))      // Pink (hard to see difference)
```

**Improved (high contrast)**:
```go
paneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#666666"))      // Dark gray (inactive)

activePaneStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#00D084")).     // Bright teal/green (ACTIVE)
    Bold(true)                                        // Add bold for more emphasis
```

**Alternative color schemes**:

**Option A: Neon Green (High Contrast)**
```go
BorderForeground(lipgloss.Color("#00FF00"))         // Bright neon green
```

**Option B: Electric Blue**
```go
BorderForeground(lipgloss.Color("#00BFFF"))         // Deep sky blue
```

**Option C: Modern Purple**
```go
// Inactive
BorderForeground(lipgloss.Color("#444444"))         // Charcoal
// Active
BorderForeground(lipgloss.Color("#A78BFA"))         // Bright purple
```

Choose the one that feels right for you!

---

## Step 5: Add Selection Highlighting to Viewer

To make selected text visually distinct, add to the `View()` function:

```go
// Viewer pane
viewerStyle := paneStyle
if m.currentView == ViewerView {
    viewerStyle = activePaneStyle
}

viewerContent := m.viewer.View()
if m.selectedFile == "" {
    viewerContent = "No file selected\n\nNavigate in the file tree and press Enter to view a file."
} else if m.clipboard.HasSelection() {
    // NEW: Add visual indicator when text is selected
    selectedStart, selectedCol, selectedEnd, selectedEol := m.clipboard.GetSelectionBounds()
    // Will implement visual highlighting in next phase
    _ = selectedStart // Using these for now
    _ = selectedCol
    _ = selectedEnd
    _ = selectedEol
}

viewerPane := viewerStyle.
    Width(m.viewerWidth).
    Render(viewerContent)
```

---

## Step 6: Update Status Bar

Update the status bar text to show copy capability (around line 238):

```go
case ViewerView:
    statusText = fmt.Sprintf(
        "[%s] Tab: switch | j/k: scroll | d/u: page | g/G: top/bottom | y: copy | ?: help | q: quit",
        viewName,
    )
```

---

## Step 7: Add Go Module Dependency

Update `go.mod`:

```bash
go get github.com/atotto/clipboard
```

---

## Step 8: Test Integration

```bash
# Build with new features
go build -o lumina

# Run and test
./lumina ~/.config

# Test copy functionality:
# 1. Navigate to a markdown file
# 2. Try selecting text (will be enhanced)
# 3. Press 'y' to copy
# 4. Paste elsewhere to verify
```

---

## What You Get

✅ **Custom keybindings** loaded from `~/.config/lumina/keybindings.json`
✅ **One-handed navigation** (e.g., `d` for down, `e` for up)
✅ **Copy to clipboard** with `y` key
✅ **Clear pane indication** with new green/blue borders
✅ **Backward compatible** - works with existing code
✅ **User friendly** - auto-creates config on first run

---

## Usage Example: One-Handed Config

Create `~/.config/lumina/keybindings.json`:

```json
{
  "scrolling": [
    {"key": "j", "action": "scroll_down", "view": "viewer"},
    {"key": "k", "action": "scroll_up", "view": "viewer"},
    {"key": "d", "action": "page_down", "view": "viewer"},
    {"key": "e", "action": "page_up", "view": "viewer"}
  ]
}
```

Now `d` = down, `e` = up - perfect for one-handed use! 🎯

---

## Next Phase: Enhanced Selection

Future work can add:
- Mouse click + drag selection
- Visual highlighting of selected text
- Multiple selection modes (character, line, block)
- Selection status in status bar
- Copy feedback ("Copied N characters")

But the infrastructure is ready now!
