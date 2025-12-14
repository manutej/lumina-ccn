# LUMINA - Claude Code Navigator
## Project-Specific Claude Code Configuration

**Project Path**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/`
**Inherits From**: `~/.claude/CLAUDE.md` (Global Configuration)
**Last Updated**: 2025-12-11
**Version**: v1.4.1-alpha
**Status**: Phase 3 In Progress (Milestones 1-2 Complete)
**Current Branch**: `feature/phase-3-ui-integration`

---

## Vision & Purpose

**LUMINA** (Claude Code Navigator) is a production-ready TUI application for navigating and viewing Claude Code workflows with beautiful markdown rendering. Built with the Charm ecosystem (Bubble Tea, Glamour, Lip Gloss), it solves a critical workflow problem: helping developers quickly navigate through extensive documentation, specifications, and workflow files.

### Core Problem It Solves
Users of Claude Code need a fast, beautiful way to browse project documentation with:
- Quick file navigation across large codebases
- Rendered markdown viewing with syntax highlighting
- Content search across files (both filename and content)
- Visual context about documents (TOC, stats, metadata)

### Design Philosophy
- **Pragmatic Technology Choice**: Leverage battle-tested Charm components over building from scratch
- **Elm Architecture (TEA)**: Single source of truth, pure functions, async via tea.Cmd
- **Backend-First Development**: Implement and test backends before UI integration
- **Incremental Phase-Based Delivery**: Each phase adds user value

---

## Critical Build Information

### Binary Alias Issue (LEARNED 2025-12-11)

The `lumina` shell alias points to `ccn`, NOT `lumina`:
```bash
alias lumina='/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/ccn'
```

**ALWAYS rebuild `ccn` after code changes:**
```bash
go build -o ccn
```

**Verification:**
```bash
which lumina          # Shows: aliased to .../ccn
ls -la ccn            # Check timestamp is current
```

---

## Quick Facts

| Property | Value |
|----------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Bubble Tea (TUI), Charm Stack |
| **Build Command** | `go build -o ccn` |
| **Run Command** | `lumina` (alias to `ccn`) |
| **Binary Size** | ~14MB (arm64) |
| **Test Count** | 55+ tests passing |
| **Codebase** | ~5,900 lines across 21 Go files |

---

## Key Bindings

| Key | Mode | Action | Description |
|-----|------|--------|-------------|
| `/` | Normal | Fuzzy Finder | Search files by NAME across codebase |
| `Ctrl+F` | Normal | Ripgrep Search | Search file CONTENT across codebase |
| `?` | Normal | Help | Show help overlay |
| `Tab` | Any | Switch View | Cycle between file tree, viewer, preview |
| `j/k` | Any | Navigate | Move down/up |
| `d/u` | Viewer | Page | Half-page down/up |
| `g/G` | Viewer | Jump | Top/bottom of document |
| `Enter` | Finder/Search | Select | Open file or jump to result |
| `Esc` | Modal | Cancel | Exit current modal |
| `q` | Normal | Quit | Exit application |
| `y` | Viewer | Copy | Copy document to clipboard |
| `m` | Preview | Mode | Cycle context panel modes |
| `Shift+A-Z` | File Tree | Jump | Jump to file starting with letter |

### Search Features Clarification

- **`/` (Fuzzy Finder)**: Searches **file names** - find files by path/name pattern
- **`Ctrl+F` (Ripgrep Search)**: Searches **file contents** - find text patterns within files

Both are **cross-file** searches across the entire codebase.

### Supported File Types

| Type | Extension | Features |
|------|-----------|----------|
| **Markdown** | `.md` | Glamour rendering, TOC, document stats |
| **JSON** | `.json` | Pretty-print formatting, syntax highlighting, JSON stats |

---

## Architecture

### Technology Stack
- **TUI Framework**: Bubble Tea (MVU/Elm architecture)
- **UI Components**: Bubbles (list, viewport)
- **Markdown Rendering**: Glamour (from Glow)
- **Terminal Styling**: Lip Gloss
- **File Watching**: fsnotify
- **Fuzzy Matching**: sahilm/fuzzy
- **Clipboard**: atotto/clipboard

### State Machine (UIMode)

```
NormalMode (0) ──┬── "/" ──────> FinderMode (1)
                 ├── "Ctrl+F" ─> SearchMode (2)
                 ├── "?" ──────> HelpMode (3)
                 └── async ────> LoadingMode (4)
```

### Elm Architecture (TEA)
- **Model**: Single source of truth (`AppModel` struct)
- **Update**: Pure functions processing messages, returns `(Model, Cmd)`
- **View**: Pure rendering from state

### Three-Pane Layout
```
┌─────────────┬──────────────────────┬─────────────┐
│  File Tree  │       Viewer         │   Preview   │
│    (1/4)    │       (2/4)          │    (1/4)    │
│             │                      │             │
│  Navigate   │  Markdown Rendered   │  TOC/Stats  │
│  j/k/Enter  │  with Glamour        │  4 modes    │
└─────────────┴──────────────────────┴─────────────┘
```

### Context Panel Modes (Right Pane)
1. **Table of Contents**: Navigate by heading, jump with Enter
2. **File Info**: Metadata (size, word count, reading time)
3. **Document Stats**: Structure analysis (headings, code blocks, links)
4. **Quick Actions**: Foundation for future features

---

## Phase History & Current Status

### Phase 1 (MVP) - COMPLETE
**v1.0.0-alpha** (2025-10-20)
- 3-pane layout (file tree | viewer | preview)
- File navigation with directory browsing
- Glamour markdown rendering
- Vim-style keybindings
- Help overlay system

### Phase 1.5 (UX Quality) - COMPLETE
**v1.0.1-alpha** (2025-10-21)
- Custom keybindings system (JSON config)
- Copy/Selection capability (system clipboard)
- Improved pane color distinction
- Table of Contents foundation

### Phase 2 (Backend Systems) - COMPLETE
**v1.3.0-alpha** (2025-11-11)
- Fuzzy file finder backend (sahilm/fuzzy)
- Ripgrep integration with streaming results
- File watching with fsnotify
- Shift+Letter quick jump navigation

### Phase 3 (UI Integration) - IN PROGRESS
**v1.4.1-alpha** (2025-12-11)

| Milestone | Status | Tests | Key Features |
|-----------|--------|-------|--------------|
| **M1: Fuzzy Finder UI** | Complete | 18 | `/` trigger, async loading, fuzzy filter |
| **M2: Ripgrep Search UI** | Complete | 15 | `Ctrl+F` trigger, streaming results |
| **M3: File Watcher UI** | Pending | - | Auto-reload, change detection |
| **M4: Polish & Docs** | Pending | - | Help updates, documentation |

---

## Project Structure

```
ccn/
├── .claude/
│   └── CLAUDE.md                    # This file
│
├── main.go                          # Entry point, UI, message handling (~1200 lines)
├── model.go                         # AppModel state definition (~650 lines)
├── keybindings.go                   # Configurable keybindings (175 lines)
├── context_panel.go                 # Multi-mode right pane (385 lines)
├── toc.go                           # Table of Contents (237 lines)
├── clipboard.go                     # Copy functionality (221 lines)
├── ripgrep.go                       # Ripgrep search manager (~200 lines)
├── colors.go                        # Color management
├── help.go                          # Help overlay
├── version.go                       # Version info
│
├── *_test.go                        # Test files (55+ tests)
│
├── ccn                              # PRIMARY BINARY (alias target)
├── lumina                           # Secondary binary
│
├── docs/
│   ├── PHASE_3_META_PROMPT_v6.md    # Phase 3 specification
│   ├── ELM_ARCHITECTURE_FOR_CCN.md  # TEA pattern explanation
│   ├── GIT_WORKFLOW.md              # Git workflow guide
│   ├── CONFIG.md                    # Keybindings configuration
│   ├── TESTING_STRATEGY.md          # Test design
│   └── ...
│
├── PHASE_3_IMPLEMENTATION_SPEC.md   # Implementation details
├── PHASE_3_BLOCKERS_FIXED.md        # Blocker resolutions
├── CHANGELOG.md                     # Version history
├── README.md                        # Project overview
└── go.mod, go.sum                   # Dependencies
```

---

## Key Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Full feature overview, installation, usage |
| `CHANGELOG.md` | Detailed version history |
| `PHASE_3_IMPLEMENTATION_SPEC.md` | Current phase milestones |
| `docs/PHASE_3_META_PROMPT_v6.md` | Implementation reference |
| `docs/ELM_ARCHITECTURE_FOR_CCN.md` | Architecture patterns |
| `docs/CONFIG.md` | Keybindings customization |

---

## Configuration

### Keybindings
**Location**: `~/.config/lumina/keybindings.json` (auto-created on first run)

**Categories**:
- **Navigation**: File tree movement (j, k, Enter, Backspace)
- **Scrolling**: Viewer paging (d, u, g, G)
- **Actions**: Search, help, copy (/, ?, y)
- **App Controls**: Quit (q, Ctrl+C)

**Example - One-Handed Navigation**:
```json
{
  "scrolling": [
    {"key": "d", "action": "page_down", "view": "viewer"},
    {"key": "e", "action": "page_up", "view": "viewer"}
  ]
}
```

---

## Development Workflow

### Build & Test
```bash
# Build (ALWAYS use ccn for alias compatibility)
go build -o ccn

# Run tests
go test ./...

# Run specific test
go test -v -run TestSearchModeTyping

# Run app
lumina
```

### Common Issues

#### Changes Not Appearing in App
```bash
# Check which binary the alias uses
which lumina
# Output: lumina: aliased to /Users/manu/.../ccn

# Rebuild the correct binary
go build -o ccn
```

---

## Go Patterns

### Value vs Pointer Receivers
The codebase uses value receivers `(m AppModel)` for Bubble Tea compatibility:

```go
// Value receiver - returns modified copy
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.someField = newValue  // Modifies local copy
    return m, nil           // Return modified copy
}

// Pointer receiver - modifies in place
func (m *AppModel) transitionTo(newMode UIMode) {
    m.currentMode = newMode  // Modifies original
}
```

### Async Pattern (tea.Cmd)
```go
// Non-blocking operations return commands
func startSearchCmd(rootPath, query string) tea.Cmd {
    return func() tea.Msg {
        // Async work here
        return SearchResultMsg{...}
    }
}
```

---

## Performance Targets

| Metric | Target | Actual |
|--------|--------|--------|
| Binary Size | <20MB | 14MB |
| Startup Time | <100ms | ~50ms |
| Memory Base | <20MB | ~10MB |
| File Loading | <100ms | <10ms |
| Markdown Render | <100ms | <50ms |
| Fuzzy Filter | <50ms/1000 files | <50ms |

---

## Recommended Agents

| Agent | Use Case |
|-------|----------|
| **git-genius** | Git operations, commits, branches |
| **practical-programmer** | Code quality, DRY/KISS principles |
| **debug-detective** | Root cause analysis |
| **test-engineer** | Writing comprehensive tests |

---

## Future Roadmap

### Phase 3 Remaining
- [ ] File watcher auto-reload UI (Milestone 3)
- [ ] Help overlay updates (Milestone 4)

### Phase 4+ (Strategic)
- [ ] Editor integration ($EDITOR)
- [ ] Claude Code integration (/workflows, /moe)
- [ ] Agentic features (anthropic-sdk-go)
- [ ] Git awareness
- [ ] Intra-file search (within single document)

---

## What Makes LUMINA Unique

1. **Pragmatic Technology**: Charm ecosystem saves ~15 days vs building from scratch
2. **Elm Architecture**: Proper TEA pattern enables confident refactoring
3. **Backend-First**: Phase 2 backends tested before Phase 3 UI integration
4. **Incremental Delivery**: Each phase adds user value
5. **Production Quality**: Comprehensive docs, tests, error handling
6. **Context-Aware Viewing**: Right pane transforms based on document

---

## Quick Reference

```bash
# Rebuild and test
go build -o ccn && go test ./... && lumina

# Run all tests verbose
go test -v ./... 2>&1 | head -100

# Check test count
go test -v ./... 2>&1 | grep -c "^--- PASS"

# View recent commits
git log --oneline -10
```

---

**Project Maintainer**: Manu Mulaveesala
**Created With**: Claude Code + Charm Ecosystem
**Last Updated**: 2025-12-11
**Test Coverage**: 55+ tests passing
