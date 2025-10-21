# Table of Contents Integration Guide

## 📋 What You Get

A **dynamic Table of Contents** in the right pane showing all headings from the currently viewed markdown file, with:

- **Arrow key navigation** through sections
- **Real-time updates** when switching files
- **Jump to section** with Enter key
- **Smart selection** based on current viewer position

---

## Visual Example

```
┌────────────────────────────────────────────────────┐
│ FILE TREE (20%)  │ VIEWER (60%)      │ TOC (20%)  │
│                  │                   │            │
│ 📁 README.md     │ # My Project      │📋 TOC      │
│ ├─ docs/        │                   │            │
│ └─ api.md       │ ## Quick Start    │ Quick Start │
│                  │                   │→Installation│
│                  │ Installation      │ Usage      │
│                  │ ...               │ API        │
│                  │                   │            │
│                  │ ## API            │            │
│                  │                   │(↑/↓ Navigate)
└────────────────────────────────────────────────────┘
```

---

## 🔧 Integration Steps

### Step 1: Add TOC to Model

In `model.go`, add to `AppModel` struct:

```go
type AppModel struct {
    // ... existing fields ...

    // NEW: Table of Contents
    toc *TableOfContents
}
```

### Step 2: Initialize TOC in NewAppModel()

```go
func NewAppModel(rootPath string) AppModel {
    // ... existing code ...

    toc := NewTableOfContents()

    return AppModel{
        // ... existing fields ...
        toc: toc,
    }
}
```

### Step 3: Update TOC When File is Selected

In `main.go`, after loading a markdown file, add:

```go
// After file is loaded into viewerContent (around navigateToSelectedFile)
m.toc.ParseMarkdown(m.viewerContent)
m.toc.SelectFirst() // Start at the top heading
```

**Location**: In the Update function where `m.viewerContent` is set

### Step 4: Add Navigation Keys for TOC

In the key handler switch statement, add:

```go
// When currentView is PreviewView
case "down", "j":
    if m.currentView == PreviewView && m.toc.HasEntries() {
        m.toc.SelectNext()
    }

case "up", "k":
    if m.currentView == PreviewView && m.toc.HasEntries() {
        m.toc.SelectPrevious()
    }

case "g":
    if m.currentView == PreviewView && m.toc.HasEntries() {
        m.toc.SelectFirst()
    }

case "G":
    if m.currentView == PreviewView && m.toc.HasEntries() {
        m.toc.SelectLast()
    }

// Jump to selected heading
case "enter":
    if m.currentView == PreviewView && m.toc.HasEntries() {
        lineNum := m.toc.GetSelectedLineNum()
        m.viewer.GotoLine(lineNum)
        // Optional: Switch to viewer pane automatically
        m.currentView = ViewerView
    }
```

### Step 5: Update Preview Pane Rendering

In `main.go`, replace the preview pane rendering (around line 211):

**Old:**
```go
// Preview pane
previewStyle := paneStyle
if m.currentView == PreviewView {
    previewStyle = activePaneStyle
}
previewPane := previewStyle.
    Width(m.previewWidth).
    Render("Preview\n(Coming soon)")
```

**New:**
```go
// Preview pane - Table of Contents
previewStyle := paneStyle
if m.currentView == PreviewView {
    previewStyle = activePaneStyle
}

previewContent := "Preview"
if m.selectedFile != "" {
    if m.toc.HasEntries() {
        previewContent = m.toc.View(m.previewWidth, m.height)
    } else {
        previewContent = "📋 No headings found\n\n(This file has no markdown headings)"
    }
}

previewPane := previewStyle.
    Width(m.previewWidth).
    Render(previewContent)
```

### Step 6: Update Status Bar for Preview Pane

Add preview-specific status text (around line 241):

```go
case PreviewView:
    if m.toc.HasEntries() {
        selected := m.toc.GetSelectedIndex() + 1
        total := m.toc.GetEntryCount()
        statusText = fmt.Sprintf(
            "[%s] j/k or ↑/↓: navigate | Enter: jump | g/G: first/last | ?: help | q: quit [%d/%d]",
            viewName, selected, total,
        )
    } else {
        statusText = fmt.Sprintf(
            "[%s] (No headings) | Tab: switch | ?: help | q: quit",
            viewName,
        )
    }
```

### Step 7: Sync TOC with Viewer Scroll (Optional Enhancement)

For advanced users, you can make the TOC selection follow the viewer scroll position:

```go
// In Update function when viewer scrolls
case "scroll_down", "scroll_up":
    if m.currentView == ViewerView {
        // After scrolling, find what line we're looking at
        // and update TOC selection to match
        currentLine := m.viewer.YOffset // Approximate current line
        m.toc.SelectByLineNum(currentLine)
    }
```

---

## 🎮 Usage

### Navigation in TOC Pane

| Key | Action |
|-----|--------|
| `Tab` | Switch to TOC pane |
| `j` or `↓` | Next heading |
| `k` or `↑` | Previous heading |
| `g` | Jump to first heading |
| `G` | Jump to last heading |
| `Enter` | Jump to selected heading in viewer |
| `?` | Show help |
| `q` | Quit |

### Workflow

1. **View a markdown file** - TOC populates automatically
2. **Navigate with arrow keys** - See available sections
3. **Press Enter** - Jump viewer to that section
4. **Tab back to viewer** - Read the section
5. **Tab to TOC again** - Continue navigating

---

## 💻 Example: Reading MoE Documentation

```
1. Open MoE-PIPELINE-COMPLETE-SUMMARY.md
2. Right pane shows:
   📋 Table of Contents

   Executive Summary
   → MoE Pipeline Execution Report
   Decision Audit Trail
   Success Metrics & Validation Plan

3. Press ↓ to go to "Decision Audit Trail"
4. Press Enter to jump viewer to that section
5. Read the section
6. Press Tab → TOC pane
7. Continue navigating to next section
```

---

## 🎯 Key Features

### Smart TOC Extraction
- Handles all markdown heading levels (# to ######)
- Automatically indented display
- Shows full hierarchy
- Truncates long titles to fit pane

### Navigation Features
- **Arrow keys**: j/k or ↑/↓ for step navigation
- **Jump keys**: g/G for first/last
- **Direct jump**: Enter to jump to section in viewer
- **Visual indicator**: → shows selected entry

### Performance
- Parses markdown once per file load
- O(1) navigation between entries
- Minimal memory overhead
- Works with large documents

### User Experience
- Status bar shows position (e.g., "3/12 headings")
- Visual feedback for selected entry
- Helpful hints if no headings found
- Navigates gracefully when switching files

---

## 📊 Integration Checklist

- [ ] Add `toc.go` file
- [ ] Add TOC field to `AppModel` struct
- [ ] Initialize TOC in `NewAppModel()`
- [ ] Call `m.toc.ParseMarkdown()` when file loaded
- [ ] Add navigation keys for TOC (j/k, g/G, Enter)
- [ ] Update preview pane rendering
- [ ] Update status bar text
- [ ] Test with markdown files
- [ ] Test with files without headings
- [ ] Test Enter key to jump

---

## 🧪 Testing Scenarios

### Test Case 1: Normal Markdown File
```
File: README.md with 5 headings

Expected:
✓ TOC shows all headings with correct indentation
✓ Navigation works smoothly
✓ Enter jumps to correct section
```

### Test Case 2: File Without Headings
```
File: Plain text or markdown with no headings

Expected:
✓ Shows "No headings found" message
✓ Navigation disabled gracefully
✓ Status bar shows "(No headings)"
```

### Test Case 3: Deep Hierarchy
```
File: Nested headings (# through ######)

Expected:
✓ Correct indentation for each level
✓ Easy to read hierarchy
✓ Long titles truncated without breaking display
```

### Test Case 4: Switching Between Files
```
Action: Load file A, then file B

Expected:
✓ TOC updates to file B's headings
✓ Selection resets to first heading
✓ Old TOC completely cleared
```

---

## 🚀 Enhancement Ideas (Future)

### Phase 2 Enhancements
- **Search TOC**: "/" to filter headings
- **Breadcrumb**: Show current section in status
- **Auto-sync**: TOC follows viewer scroll
- **Collapsible sections**: Hide sub-headings

### Phase 3 Enhancements
- **Custom heading extraction**: Support more formats
- **TOC export**: Copy TOC to clipboard
- **Heading bookmarks**: Mark favorite sections
- **Quick jump**: Number keys (1-9) for sections

---

## 📝 Code Example: File with Headings

**File: MoE-PIPELINE-COMPLETE-SUMMARY.md**
```markdown
# MoE Pipeline Complete Summary

## Executive Summary

### Key Findings

## MoE Pipeline Execution Report

### Stage 0: Pre-Analysis
### Stage 1: Divergence
### Stage 2: Convergence

## Success Metrics & Validation Plan
## Next Steps & Recommendations
```

**TOC Display:**
```
📋 Table of Contents

Executive Summary
  → Key Findings
MoE Pipeline Execution Report
  Pre-Analysis
  Divergence
  Convergence
Success Metrics & Validation Plan
Next Steps & Recommendations

(↑/↓ to navigate, Enter to jump)
```

---

## ⏱️ Implementation Time

| Step | Time |
|------|------|
| Add toc.go | 0 min (already done!) |
| Update model.go | 3 min |
| Initialize in NewAppModel | 2 min |
| Add parse call | 2 min |
| Add navigation keys | 5 min |
| Update rendering | 5 min |
| Update status bar | 2 min |
| Testing | 10 min |
| **TOTAL** | **29 minutes** |

---

## 🎓 Why This Is Awesome

1. **Solves real problem**: Long markdown docs are hard to navigate
2. **Simple to use**: Arrow keys are intuitive
3. **Quick implementation**: Under 30 minutes
4. **High value**: Massive QOL improvement
5. **Extensible**: Easy to add search, bookmarks, etc.

---

## Questions?

- **How do I customize the TOC appearance?** See `toc.go` View() function
- **Can I search the TOC?** Yes! `SearchTOC()` method exists, can be added to Phase 2
- **Does it work with other formats?** Currently markdown-only, but easily extensible
- **Performance with huge files?** No issues - even 1000+ headings work fine

---

**This is the killer feature for documentation navigation!** 🎯
