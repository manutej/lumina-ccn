# LUMINA CCN Phase 3: Meta-Prompt v6.0 (FINAL)

**Created**: 2025-12-11 via 6-round meta-prompting framework
**Purpose**: Reference prompt for Phase 3 UI Integration implementation

---

## Meta-Prompting Evolution Summary

| Round | Focus                  | Key Additions                                             |
|-------|------------------------|-----------------------------------------------------------|
| v1.0  | Initial Draft          | Basic task description (~100 words)                       |
| v2.0  | Context & Constraints  | Architecture details, backend APIs, constraints (~400 words)|
| v3.0  | Implementation Details | Message types, code patterns, file structure (~800 words) |
| v4.0  | Examples & Edge Cases  | Visual mockups, error handling, concrete examples (~1500 words)|
| v5.0  | Actionability          | Milestones, acceptance criteria, checklists (~2000 words) |
| v6.0  | Final Polish           | Quick reference, anti-patterns, troubleshooting (~2500 words)|

---

## PROMPT v6.0 (Production Ready)

```
╔═══════════════════════════════════════════════════════════════════════╗
║ LUMINA CCN PHASE 3: UI INTEGRATION IMPLEMENTATION SPEC                ║
║ Branch: feature/phase-3-ui-integration                                ║
║ Baseline: master (stable, 135+ tests passing)                         ║
╚═══════════════════════════════════════════════════════════════════════╝

┌─────────────────────────────────────────────────────────────────────────┐
│ QUICK REFERENCE                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│ WHAT: Integrate Phase 2 backends (Finder, Search, Watcher) into UI     │
│ WHY:  Enable quick file finding, content search, and live reload       │
│ HOW:  4 milestones, each with own commit, additive changes only        │
│                                                                         │
│ MILESTONES:                                                             │
│ M1: Fuzzy Finder (/ key)     → finder_ui.go    → "feat(finder): ..."   │
│ M2: Ripgrep Search (Ctrl+F)  → search_ui.go    → "feat(search): ..."   │
│ M3: File Watcher (auto)      → watcher_ui.go   → "feat(watcher): ..."  │
│ M4: Polish & Docs            → README, help    → "docs(phase3): ..."   │
│                                                                         │
│ FILES TO CREATE: finder_ui.go, search_ui.go, watcher_ui.go, messages.go│
│ FILES TO MODIFY: main.go (Update, View), help.go, README.md            │
│ FILES TO NOT TOUCH: fuzzy_finder_impl.go, ripgrep_executor.go,         │
│                     file_watcher.go (Phase 2 backends - tested & stable)│
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ CRITICAL CONSTRAINTS                                                    │
├─────────────────────────────────────────────────────────────────────────┤
│ 1. STABILITY: Current app is stable. All changes MUST be additive.     │
│    - DO NOT refactor existing code unless absolutely necessary         │
│    - DO NOT modify Phase 2 backend implementations                     │
│    - Each milestone = 1 commit (easy rollback with git revert)         │
│                                                                         │
│ 2. STATE MACHINE: Use existing UIMode enum and transitionTo() helper   │
│    - UIMode: NormalMode, FinderMode, SearchMode, HelpMode, LoadingMode │
│    - Always clean up state when exiting a mode                         │
│    - Only ONE mode active at a time (enforced by transitionTo)         │
│                                                                         │
│ 3. BUBBLE TEA PATTERN: All async operations use tea.Cmd                │
│    - Never block in Update() - return a Cmd instead                    │
│    - Use message types for async results (SearchResultMsg, etc.)       │
│    - Cancel contexts properly to avoid goroutine leaks                 │
│                                                                         │
│ 4. TESTING: Each milestone must have tests before commit               │
│    - Unit tests for state transitions                                  │
│    - Integration tests for full workflows                              │
│    - Manual verification checklist                                     │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ ANTI-PATTERNS (What NOT To Do)                                          │
├─────────────────────────────────────────────────────────────────────────┤
│ DON'T: Modify fuzzy_finder_impl.go, ripgrep_executor.go, or            │
│        file_watcher.go - these are tested backends                     │
│ DON'T: Add global variables - all state on AppModel                    │
│ DON'T: Block in Update() with time.Sleep or channel waits              │
│ DON'T: Forget to cancel contexts (causes goroutine leaks)              │
│ DON'T: Mix mode-specific key handlers (isolate per mode)               │
│ DON'T: Skip tests - each milestone needs passing tests                 │
│ DON'T: Make large commits - one feature per commit                     │
│ DON'T: Ignore edge cases - handle empty results, errors gracefully     │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## MILESTONE 1: FUZZY FINDER UI

**Trigger**: `/` key | **Priority**: P0 | **Commit**: `feat(finder): ...`

### Visual Mockup

```
┌─────────────────────────────────────────────────────────────────────┐
│ LUMINA - /path/to/project                                           │
├─────────────────────────────────────────────────────────────────────┤
│░░░░░░░│  ╭───────────────────────────────╮  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │ Find File                     │  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  ├───────────────────────────────┤  │░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░│  │ > readme█                     │  │░░░░░░░░░░░░░░░░░░░░░░│
│░(dim)░│  ├───────────────────────────────┤  │░░░░░(dimmed)░░░░░░░░░│
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

### Code Pattern: Finder Key Handler

```go
func handleFinderKey(m *AppModel, keyMsg tea.KeyMsg) tea.Cmd {
    switch keyMsg.String() {
    case "esc":
        m.transitionTo(NormalMode)
        return nil
    case "enter":
        if len(m.finderFiltered) > 0 && m.finderCursor < len(m.finderFiltered) {
            selectedPath := m.finderFiltered[m.finderCursor]
            m.transitionTo(NormalMode)
            return func() tea.Msg {
                return FileSelectedMsg{Path: selectedPath}
            }
        }
    case "j", "down":
        m.finderCursor = navigateCursor(m.finderCursor, 1,
                                         len(m.finderFiltered))
    case "k", "up":
        m.finderCursor = navigateCursor(m.finderCursor, -1,
                                         len(m.finderFiltered))
    case "backspace":
        if len(m.finderInput) > 0 {
            m.finderInput = m.finderInput[:len(m.finderInput)-1]
            updateFinderFilter(m, m.finderInput)
        }
    default:
        // Append printable characters to filter
        if len(keyMsg.String()) == 1 {
            m.finderInput += keyMsg.String()
            updateFinderFilter(m, m.finderInput)
        }
    }
    return nil
}
```

---

## MILESTONE 2: RIPGREP SEARCH UI

**Trigger**: `Ctrl+F` | **Priority**: P0 | **Commit**: `feat(search): ...`

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

### Code Pattern: Async Search with tea.Cmd

```go
func startSearchCmd(query string) tea.Cmd {
    return func() tea.Msg {
        executor := NewRipgrepExecutor()
        ctx, cancel := context.WithTimeout(context.Background(),
                                           30*time.Second)
        ch, err := executor.Execute(ctx, query)
        if err != nil {
            cancel()
            return SearchErrorMsg{Err: err}
        }
        return SearchStartedMsg{Ch: ch, Cancel: cancel}
    }
}

func listenSearchResultCmd(ch <-chan json.RawMessage) tea.Cmd {
    return func() tea.Msg {
        result, ok := <-ch
        if !ok {
            return SearchCompleteMsg{}
        }
        var msg RipgrepMessage
        if err := json.Unmarshal(result, &msg); err != nil {
            return listenSearchResultCmd(ch)() // Skip bad result
        }
        if msg.Type == "match" {
            return SearchResultMsg{Match: ParseRipgrepMatch(msg)}
        }
        return listenSearchResultCmd(ch)() // Skip non-match
    }
}
```

---

## MILESTONE 3: FILE WATCHER UI

**Trigger**: Automatic | **Priority**: P1 | **Commit**: `feat(watcher): ...`

### Acceptance Criteria

- [ ] AC1: Watcher starts when markdown file is opened
- [ ] AC2: Watcher stops when different file is opened
- [ ] AC3: External file change triggers automatic reload
- [ ] AC4: Scroll position preserved after reload
- [ ] AC5: Status bar shows notification: "File reloaded"
- [ ] AC6: Notification auto-clears after 3 seconds
- [ ] AC7: No goroutine leaks on app exit

### Code Pattern: Scroll Preservation

```go
func (m *AppModel) reloadFileWithScrollPreservation(path string) {
    // Save scroll position
    savedYOffset := m.viewer.YOffset

    // Reload content
    m.loadFileContent(path)

    // Restore scroll position (clamped to new content length)
    maxOffset := m.viewer.TotalLineCount() - m.viewer.Height
    if savedYOffset > maxOffset {
        savedYOffset = maxOffset
    }
    if savedYOffset < 0 {
        savedYOffset = 0
    }
    m.viewer.YOffset = savedYOffset
}
```

---

## MILESTONE 4: POLISH & DOCUMENTATION

**Priority**: P2 | **Commits**: `docs(phase3): ...` + `chore(release): v1.1.0-beta`

### Checklist

- [ ] Update help.go with new keybindings section
- [ ] Update README.md with:
  - [ ] New keybindings table (/, Ctrl+F)
  - [ ] Feature descriptions (finder, search, auto-reload)
  - [ ] Screenshots of new UI modes
- [ ] Create CHANGELOG.md entry for Phase 3
- [ ] Run full test suite: `go test ./...`
- [ ] Performance validation:
  - [ ] Fuzzy filter <50ms (test with 1000 files)
  - [ ] Search first result <100ms
  - [ ] File reload <200ms (1MB file)
- [ ] Memory leak check: Run app 10 minutes, check memory stable
- [ ] Cross-mode testing: Rapid mode switching doesn't cause issues
- [ ] Tag release: git tag v1.1.0-beta

---

## TROUBLESHOOTING

| Issue | Fix |
|-------|-----|
| Modal doesn't appear centered | Use lipgloss.Place() with explicit width/height from model |
| Search results arrive but UI doesn't update | Ensure listenSearchResultCmd returns another Cmd to continue listening |
| Goroutine leak warnings | Store cancel function in model, call it on mode exit |
| Scroll position not preserved | Save YOffset BEFORE SetContent(), restore AFTER |
| Keys don't work in FinderMode | Check Update() switch order - mode-specific handlers before defaults |
| Test failures after changes | Ensure tests use mocks for FileWatcher and RipgrepExecutor |

---

## ROLLBACK PROCEDURE

If any milestone causes issues:

```bash
# Identify problematic commit
git log --oneline -5

# Revert single milestone
git revert <commit-hash> --no-edit

# Or reset to last known good state
git reset --hard <good-commit>

# Force push only if necessary (feature branch only!)
git push --force-with-lease origin feature/phase-3-ui-integration
```

---

## START HERE

1. Verify branch: `git branch` -> should show feature/phase-3-ui-integration
2. Start with Milestone 1: Create finder_ui.go
3. Follow implementation checklist step by step
4. Run tests before committing
5. Commit with message: "feat(finder): Add fuzzy finder modal UI"
6. Proceed to Milestone 2

---

**Generated via 6-round meta-prompting framework**
