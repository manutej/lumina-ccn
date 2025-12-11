# LUMINA CCN Phase 3: UI Integration Implementation Spec

**Branch**: `feature/phase-3-ui-integration`
**Baseline**: `master` (stable, 135+ tests passing)
**Created**: 2025-12-11 via 6-round meta-prompting

---

## Quick Reference

| Item | Value |
|------|-------|
| **WHAT** | Integrate Phase 2 backends (Finder, Search, Watcher) into UI |
| **WHY** | Enable quick file finding, content search, and live reload |
| **HOW** | 4 milestones, each with own commit, additive changes only |

### Milestones

| # | Feature | Trigger | File | Commit |
|---|---------|---------|------|--------|
| M1 | Fuzzy Finder | `/` key | finder_ui.go | `feat(finder): ...` |
| M2 | Ripgrep Search | `Ctrl+F` | search_ui.go | `feat(search): ...` |
| M3 | File Watcher | automatic | watcher_ui.go | `feat(watcher): ...` |
| M4 | Polish & Docs | - | README, help | `docs(phase3): ...` |

### Files

**Create**: `finder_ui.go`, `search_ui.go`, `watcher_ui.go`, `messages.go`
**Modify**: `main.go` (Update, View), `help.go`, `README.md`
**DO NOT TOUCH**: `fuzzy_finder_impl.go`, `ripgrep_executor.go`, `file_watcher.go`

---

## Critical Constraints

1. **STABILITY**: Current app is stable. All changes MUST be additive.
   - DO NOT refactor existing code unless absolutely necessary
   - DO NOT modify Phase 2 backend implementations
   - Each milestone = 1 commit (easy rollback with git revert)

2. **STATE MACHINE**: Use existing UIMode enum and transitionTo() helper
   - UIMode: NormalMode, FinderMode, SearchMode, HelpMode, LoadingMode
   - Always clean up state when exiting a mode
   - Only ONE mode active at a time (enforced by transitionTo)

3. **BUBBLE TEA PATTERN**: All async operations use tea.Cmd
   - Never block in Update() - return a Cmd instead
   - Use message types for async results (SearchResultMsg, etc.)
   - Cancel contexts properly to avoid goroutine leaks

4. **TESTING**: Each milestone must have tests before commit

---

## Anti-Patterns (What NOT To Do)

- DON'T modify fuzzy_finder_impl.go, ripgrep_executor.go, or file_watcher.go
- DON'T add global variables - all state on AppModel
- DON'T block in Update() with time.Sleep or channel waits
- DON'T forget to cancel contexts (causes goroutine leaks)
- DON'T mix mode-specific key handlers (isolate per mode)
- DON'T skip tests - each milestone needs passing tests
- DON'T make large commits - one feature per commit
- DON'T ignore edge cases - handle empty results, errors gracefully

---

## Milestone 1: Fuzzy Finder UI

**Trigger**: `/` key | **Priority**: P0

### Visual Mockup

```
┌─────────────────────────────────────────────────────────────────────┐
│ LUMINA - /path/to/project                                           │
├─────────────────────────────────────────────────────────────────────┤
│░░░░░░░│  ╭───────────────────────────────╮  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │ Find File                     │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  ├───────────────────────────────┤  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │ > readme█                     │  │░░░░░░░░░░░░░░░░░░░░░░│
│░(dim)░│  ├───────────────────────────────┤  │░░░(dimmed)░░░░░░░░░░░│
│░░░░░░░│  │ > README.md                   │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │   docs/README.md              │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │   archive/README.md           │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  ├───────────────────────────────┤  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │ 3 of 42 files                 │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  ╰───────────────────────────────╯  │░░░░░░░░░░░░░░░░░░░░░░│
├─────────────────────────────────────────────────────────────────────┤
│ [FINDER] j/k:navigate | Enter:open | Esc:cancel                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Acceptance Criteria

- [ ] AC1: `/` from NormalMode -> FinderMode with modal overlay
- [ ] AC2: Typing updates filter in real-time (<50ms)
- [ ] AC3: j/k moves cursor with wrap-around
- [ ] AC4: Enter opens file and returns to NormalMode
- [ ] AC5: Esc returns to NormalMode with no side effects
- [ ] AC6: Empty results shows "No matches for '{query}'"
- [ ] AC7: Background panes visually dimmed

### Implementation Checklist

- [ ] Create finder_ui.go with:
  - [ ] renderFinderModal(m AppModel) string
  - [ ] handleFinderKeypress(m *AppModel, msg tea.KeyMsg) (tea.Model, tea.Cmd)
  - [ ] updateFinderFilter(m *AppModel, query string)
- [ ] Modify main.go Update():
  - [ ] Add "/" key handler -> m.transitionTo(FinderMode)
  - [ ] Add FinderMode case -> delegate to handleFinderKeypress
  - [ ] Initialize m.finderItems with findMarkdownFiles() on first open
- [ ] Modify main.go View():
  - [ ] When FinderMode, call renderFinderModal() as overlay
  - [ ] Dim background (reduce color brightness or add overlay)
- [ ] Create finder_ui_test.go with:
  - [ ] TestFinderOpen
  - [ ] TestFinderFilter
  - [ ] TestFinderNavigation
  - [ ] TestFinderSelect
  - [ ] TestFinderCancel
- [ ] Manual verification:
  - [ ] Run `./lumina`, press `/`, type query, select, verify file opens
  - [ ] Test Esc cancels cleanly
  - [ ] Test empty query shows all files

---

## Milestone 2: Ripgrep Search UI

**Trigger**: `Ctrl+F` | **Priority**: P0

### Visual Mockup

```
┌─────────────────────────────────────────────────────────────────────┐
│ LUMINA - /path/to/project                                           │
├──────────────┬────────────────────────┬─────────────────────────────┤
│ Files [A-Z]  │ Search: "TODO"         │ # main.go                   │
│              │ ───────────────────────│                             │
│ > README.md  │ > main.go:42           │ Line 42:                    │
│   src/       │   // TODO: fix nav     │ // TODO: fix navigation     │
│   docs/      │   model.go:156         │                             │
│              │   // TODO: cleanup     │ (context around match)      │
│              │   search.go:89         │                             │
│              │   // TODO: async       │                             │
│              │ ───────────────────────│                             │
│              │ 3 matches (0.05s)      │                             │
├──────────────┴────────────────────────┴─────────────────────────────┤
│ [SEARCH] j/k:navigate | Enter:jump | n/N:next/prev | Esc:close      │
└─────────────────────────────────────────────────────────────────────┘
```

### Acceptance Criteria

- [ ] AC1: Ctrl+F -> SearchMode with input prompt
- [ ] AC2: Enter executes search, results stream in asynchronously
- [ ] AC3: j/k navigates results
- [ ] AC4: Enter on result jumps to file:line in viewer
- [ ] AC5: n/N cycles through matches
- [ ] AC6: Esc cancels and returns to NormalMode
- [ ] AC7: "ripgrep not found" shows install instructions
- [ ] AC8: Results limited to 500, shows "X of Y+" indicator
- [ ] AC9: Empty query shows "Enter search query"

### Implementation Checklist

- [ ] Create messages.go with message types
- [ ] Create search_ui.go with rendering and key handling
- [ ] Modify main.go for SearchMode integration
- [ ] Create search_ui_test.go with tests

---

## Milestone 3: File Watcher UI

**Trigger**: Automatic | **Priority**: P1

### Acceptance Criteria

- [ ] AC1: Watcher starts when markdown file is opened
- [ ] AC2: Watcher stops when different file is opened
- [ ] AC3: External file change triggers automatic reload
- [ ] AC4: Scroll position preserved after reload
- [ ] AC5: Status bar shows notification: "File reloaded"
- [ ] AC6: Notification auto-clears after 3 seconds
- [ ] AC7: No goroutine leaks on app exit

### Implementation Checklist

- [ ] Create watcher_ui.go with lifecycle management
- [ ] Modify model.go for watcher integration
- [ ] Modify main.go for FileChangedMsg handling
- [ ] Create watcher_ui_test.go with tests

---

## Milestone 4: Polish & Documentation

**Priority**: P2

### Checklist

- [ ] Update help.go with new keybindings section
- [ ] Update README.md with new features
- [ ] Create CHANGELOG.md entry for Phase 3
- [ ] Run full test suite: `go test ./...`
- [ ] Performance validation
- [ ] Memory leak check
- [ ] Tag release: git tag v1.1.0-beta

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Modal doesn't appear centered | Use lipgloss.Place() with explicit width/height |
| Search results don't update UI | Ensure listenSearchResultCmd returns another Cmd |
| Goroutine leak warnings | Store cancel function in model, call on mode exit |
| Scroll position not preserved | Save YOffset BEFORE SetContent(), restore AFTER |
| Keys don't work in FinderMode | Check Update() switch order - mode handlers first |

---

## Rollback Procedure

```bash
# Identify problematic commit
git log --oneline -5

# Revert single milestone
git revert <commit-hash> --no-edit

# Or reset to last known good state
git reset --hard <good-commit>
```

---

**Generated via 6-round meta-prompting framework**
