# LUMINA Testing Strategy

## Overview

Comprehensive testing strategy for LUMINA covering unit tests, integration tests, and manual testing procedures.

## Testing Pyramid

```
                  /\
                 /  \
                /    \
               /Manual\           Small - Exploratory testing
              /--------\
             /Integration\        Medium - Feature workflows
            /------------\
           /  Unit Tests  \       Large - Individual functions
          /----------------\
```

## Test Types

### 1. Unit Tests

Test individual functions and components in isolation.

**Coverage Targets:**
- Core logic: 80%+
- Utilities: 90%+
- UI components: 60%+ (harder to test)

**Example Test Files:**

```
ccn/
├── clipboard_test.go          # Clipboard operations
├── keybindings_test.go        # Keybinding parsing
├── toc_test.go                # TOC extraction
├── search_test.go             # Search algorithms
└── version_test.go            # Version info
```

**Writing Unit Tests:**

```go
// clipboard_test.go
package main

import (
    "testing"
)

func TestClipboardManager_CopyText(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid text",
            input:   "Hello, World!",
            wantErr: false,
        },
        {
            name:    "empty text",
            input:   "",
            wantErr: false,
        },
        {
            name:    "unicode text",
            input:   "こんにちは",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cm := NewClipboardManager()
            err := cm.CopyText(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("CopyText() error = %v, wantErr %v", err, tt.wantErr)
            }

            if err == nil {
                // Verify clipboard content
                got, _ := cm.GetText()
                if got != tt.input {
                    t.Errorf("CopyText() = %v, want %v", got, tt.input)
                }
            }
        })
    }
}

func BenchmarkClipboardCopy(b *testing.B) {
    cm := NewClipboardManager()
    text := "Sample markdown content"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = cm.CopyText(text)
    }
}
```

### 2. Integration Tests

Test feature workflows and component interactions.

```go
// integration_test.go
package main

import (
    "os"
    "path/filepath"
    "testing"
)

func TestFullWorkflow_FileNavigation(t *testing.T) {
    // Setup: Create temp directory with markdown files
    tmpDir := t.TempDir()
    createTestFiles(t, tmpDir)

    // Initialize app model
    model := NewAppModel(tmpDir)

    // Test: Navigate through files
    model.Update(tea.KeyMsg{Type: tea.KeyDown})
    model.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Verify: File loaded correctly
    if model.currentFile == "" {
        t.Error("Expected file to be loaded")
    }

    // Verify: Content rendered
    if model.viewportContent == "" {
        t.Error("Expected viewport to have content")
    }
}

func TestFullWorkflow_CopyAndSearch(t *testing.T) {
    tmpDir := t.TempDir()
    createTestFiles(t, tmpDir)

    model := NewAppModel(tmpDir)

    // Navigate to file
    model.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Search for text
    model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
    model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("TODO")})
    model.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Copy selection
    model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

    // Verify clipboard
    cm := model.clipboard
    text, err := cm.GetText()
    if err != nil {
        t.Fatalf("Failed to get clipboard: %v", err)
    }

    if text != "TODO" {
        t.Errorf("Expected 'TODO', got %v", text)
    }
}
```

### 3. Table-Driven Tests

Use table-driven tests for comprehensive coverage:

```go
// toc_test.go
func TestExtractHeadings(t *testing.T) {
    tests := []struct {
        name     string
        markdown string
        want     []Heading
    }{
        {
            name: "simple headings",
            markdown: `# Title
## Subtitle
### Section`,
            want: []Heading{
                {Level: 1, Text: "Title", Line: 0},
                {Level: 2, Text: "Subtitle", Line: 1},
                {Level: 3, Text: "Section", Line: 2},
            },
        },
        {
            name: "headings with formatting",
            markdown: `# **Bold** Title
## _Italic_ Subtitle`,
            want: []Heading{
                {Level: 1, Text: "Bold Title", Line: 0},
                {Level: 2, Text: "Italic Subtitle", Line: 1},
            },
        },
        {
            name: "no headings",
            markdown: "Just plain text\nNo headings here",
            want: []Heading{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ExtractHeadings(tt.markdown)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ExtractHeadings() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 4. Benchmark Tests

Performance benchmarks for critical paths:

```go
// search_test.go
func BenchmarkFuzzySearch(b *testing.B) {
    files := generateTestFiles(1000) // 1000 files

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = FuzzySearch(files, "readme")
    }
}

func BenchmarkMarkdownRendering(b *testing.B) {
    content := loadLargeMarkdownFile() // ~1MB file

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = RenderMarkdown(content)
    }
}
```

## Running Tests

### Basic Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...

# Run specific test
go test -run TestClipboardManager_CopyText

# Run benchmarks
go test -bench=. ./...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./...
```

### Advanced Test Commands

```bash
# Verbose output
go test -v ./...

# Show coverage by function
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Generate coverage badge
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}'

# Continuous testing (requires entr or similar)
find . -name "*.go" | entr -c go test ./...

# Test with timeout
go test -timeout 30s ./...

# Parallel testing
go test -parallel 4 ./...
```

## Test Organization

### Directory Structure

```
ccn/
├── clipboard.go
├── clipboard_test.go           # Unit tests for clipboard
├── keybindings.go
├── keybindings_test.go         # Unit tests for keybindings
├── toc.go
├── toc_test.go                 # Unit tests for TOC
├── search.go
├── search_test.go              # Unit tests for search
├── testdata/                   # Test fixtures
│   ├── sample.md
│   ├── large.md
│   └── keybindings.json
└── integration_test.go         # Integration tests
```

### Test Fixtures

```
testdata/
├── markdown/
│   ├── simple.md               # Basic markdown
│   ├── complex.md              # Complex formatting
│   ├── large.md                # Performance testing
│   └── unicode.md              # Unicode handling
├── config/
│   ├── valid_keybindings.json
│   ├── invalid_keybindings.json
│   └── minimal_keybindings.json
└── expected/
    ├── toc_simple.json         # Expected TOC output
    └── search_results.json     # Expected search results
```

## Testing Checklist

### Before Commit (Automated via Git Hook)

- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Code formatted: `gofmt -s -w .`
- [ ] Linting passes: `go vet ./...`
- [ ] Build successful: `go build -o lumina`

### Before Release

- [ ] Full test suite passes on all platforms
- [ ] Coverage meets targets (>70%)
- [ ] Benchmarks show acceptable performance
- [ ] Manual testing completed (see below)
- [ ] Integration tests pass
- [ ] No memory leaks (profiling)
- [ ] Documentation updated

## Manual Testing Procedures

### Smoke Tests

**Quick validation after build:**

```bash
# 1. Version check
./lumina --version
# Expected: Version 1.0.1-alpha, Phase 1.5, Date 2025-10-21

# 2. Help output
./lumina --help
# Expected: Usage information displayed

# 3. Basic navigation
./lumina ~/.config
# Expected: File tree appears, can navigate with j/k

# 4. File viewing
# Press Enter on a markdown file
# Expected: File content rendered with Glamour

# 5. Help overlay
# Press ?
# Expected: Help overlay appears with keybindings

# 6. Quit
# Press q
# Expected: Application exits cleanly
```

### Feature Testing

**Phase 1.5 Features:**

```bash
# 1. Custom Keybindings
# - Edit ~/.config/lumina/keybindings.json
# - Change "down": ["j"] to "down": ["x"]
# - Restart lumina
# - Verify 'x' moves down instead of 'j'

# 2. Copy/Selection
# - Navigate to markdown file
# - Select text (implementation pending)
# - Press 'y' to copy
# - Paste into another app
# - Verify correct content copied

# 3. Pane Color Distinction
# - Open lumina with multiple panes
# - Verify active pane is bright teal (#00D084)
# - Verify inactive panes are dark gray (#666666)
# - Switch tabs, verify colors update

# 4. Table of Contents
# - Open markdown file with headings
# - Verify TOC appears (implementation pending)
# - Navigate TOC with arrow keys
# - Press Enter to jump to section
```

### Cross-Platform Testing

**Test on each platform:**

- [ ] macOS (arm64, amd64)
- [ ] Linux (amd64, arm64)
- [ ] Windows (amd64)

**Platform-specific checks:**

```bash
# macOS
./lumina ~/.config
# Test: Native clipboard integration
# Test: Option+Arrow keys work

# Linux
./lumina ~/.config
# Test: xclip/xsel clipboard support
# Test: Terminal color support

# Windows
lumina.exe C:\Users\User\.config
# Test: Windows clipboard API
# Test: Terminal color support (Windows Terminal)
```

### Performance Testing

```bash
# 1. Large file rendering
./lumina ~/large-markdown-file.md
# Expected: Renders in <500ms

# 2. Deep directory tree
./lumina ~/projects
# Expected: File tree builds in <1s

# 3. Rapid navigation
# Quickly press j/k repeatedly
# Expected: No lag, smooth scrolling

# 4. Memory usage
# Open lumina, monitor with top/htop
# Expected: <50MB memory usage
```

### Regression Testing

**Test all previous features still work:**

- [ ] File tree navigation (Phase 1)
- [ ] Markdown rendering (Phase 1)
- [ ] Vim keybindings (Phase 1)
- [ ] Help overlay (Phase 1)
- [ ] Custom keybindings (Phase 1.5)
- [ ] Clipboard operations (Phase 1.5)
- [ ] Pane colors (Phase 1.5)

## Continuous Integration

### GitHub Actions Workflow

Tests automatically run on:
- Push to `main` or `develop`
- Pull requests to `main` or `develop`
- Git tags matching `v*`

**Test matrix:**
- Go 1.25.3
- Platforms: Ubuntu, macOS, Windows
- Test types: Unit, Integration, Benchmarks

### Coverage Requirements

**Required coverage before merge:**
- New code: 70%+ coverage
- Modified code: Maintain or improve coverage
- Overall project: 60%+ coverage

**Coverage gates:**
```bash
# Fail if coverage drops below 60%
go test -coverprofile=coverage.out ./...
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$coverage < 60" | bc -l) )); then
    echo "Coverage $coverage% is below 60%"
    exit 1
fi
```

## Test Maintenance

### Updating Tests

When adding features:
1. Write test first (TDD approach)
2. Implement feature
3. Verify test passes
4. Add integration test if needed
5. Update manual test checklist

When fixing bugs:
1. Write test that reproduces bug
2. Verify test fails
3. Fix bug
4. Verify test passes
5. Add to regression test suite

### Flaky Tests

If tests are flaky:
```bash
# Run test 100 times to verify stability
go test -count=100 -run TestProblematicTest

# Run with race detector
go test -race -count=10 -run TestProblematicTest
```

Fix flaky tests by:
- Removing time dependencies (use mocks)
- Avoiding global state
- Ensuring proper cleanup in tests
- Using table-driven tests

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Best Practices: Testing](https://go.dev/doc/effective_go#testing)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testing Best Practices](https://golang.org/doc/code.html#Testing)

---

**Last Updated**: 2025-10-21
