# Testing Guide - LUMINA Phase 2

**Version**: v1.1.0-alpha
**Date**: 2025-10-27
**Test Coverage**: 60%+ target

This guide explains how to run, write, and debug tests for the LUMINA project following Test-Driven Development (TDD) principles.

---

## Table of Contents

1. [Running Tests](#running-tests)
2. [Test Structure](#test-structure)
3. [TDD Workflow (RED-GREEN-REFACTOR)](#tdd-workflow)
4. [Writing New Tests](#writing-new-tests)
5. [Test Fixtures and Mocks](#test-fixtures-and-mocks)
6. [Debugging Failing Tests](#debugging-failing-tests)
7. [Integration Tests](#integration-tests)
8. [Performance Tests](#performance-tests)

---

## Running Tests

### Run All Tests

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn
go test -v
```

**Output**:
```
=== RUN   TestFileWatcherBasicFunctionality
=== RUN   TestFileWatcherDetectChanges
...
PASS
ok      github.com/lumina/ccn    2.456s
```

### Run Specific Test File

```bash
go test -v -run TestFileWatcher file_watcher_test.go file_watcher.go test_mocks.go
```

### Run Specific Test Function

```bash
go test -v -run TestFileWatcherBasicFunctionality
```

### Run with Coverage

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out  # View in browser
```

### Run Short Tests Only

```bash
go test -short  # Skips performance/integration tests
```

### Run Tests with Race Detector

```bash
go test -race  # Detect race conditions
```

### Watch Mode (Auto-run on Changes)

```bash
# Install entr
brew install entr

# Watch and auto-test
ls *.go | entr -c go test -v
```

---

## Test Structure

### File Organization

```
/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/
├── file_watcher.go              # Implementation
├── file_watcher_test.go         # Unit tests
├── ripgrep_executor.go          # Implementation
├── ripgrep_integration_test.go  # Integration tests
├── fuzzy_finder_impl.go         # Implementation
├── fuzzy_finder_test.go         # Unit tests
├── keybinding_impl.go           # Implementation
├── keybinding_test.go           # Unit tests
├── glamour_impl.go              # Implementation
├── glamour_integration_test.go  # Integration tests
├── integration_test.go          # Full integration tests
└── test_mocks.go                # Test mocks and fixtures
```

### Test Naming Convention

```go
// Unit tests
func TestFileWatcherBasicFunctionality(t *testing.T) { }
func TestFuzzyFinderNavigation(t *testing.T) { }

// Integration tests
func TestRipgrepIntegration(t *testing.T) { }
func TestGlamourIntegration(t *testing.T) { }

// Performance tests
func TestFileWatcherPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping performance test in short mode")
    }
}
```

### Table-Driven Tests

All tests follow table-driven pattern:

```go
func TestExample(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case_1", "input1", "output1"},
        {"case_2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## TDD Workflow

### Phase 1: RED (Write Failing Test)

**Goal**: Write a test that fails because feature doesn't exist yet.

**Example** (FileWatcher):

```go
func TestFileWatcherBasicFunctionality(t *testing.T) {
    tests := []struct {
        name        string
        path        string
        expectError bool
    }{
        {"watch_valid_directory", "/tmp/test", false},
        {"watch_nonexistent", "/nonexistent/path", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            watcher := NewFileWatcher()
            defer watcher.Close()

            // Act
            err := watcher.Watch(tt.path)

            // Assert
            if (err != nil) != tt.expectError {
                t.Errorf("Watch() error = %v, expectError %v", err, tt.expectError)
            }
        })
    }
}
```

**Run Test**:
```bash
go test -v -run TestFileWatcherBasicFunctionality
```

**Expected Output** (RED):
```
--- FAIL: TestFileWatcherBasicFunctionality (0.00s)
    --- FAIL: TestFileWatcherBasicFunctionality/watch_valid_directory (0.00s)
        file_watcher_test.go:37: Watch() error = <nil>, expectError false
```

### Phase 2: GREEN (Make It Pass)

**Goal**: Write minimal code to make the test pass.

**Implementation**:

```go
func (fw *FileWatcherImpl) Watch(path string) error {
    // Resolve symlinks
    resolved, err := filepath.EvalSymlinks(path)
    if err != nil {
        return fmt.Errorf("no such file or directory")
    }

    // Check if directory exists
    info, err := os.Stat(resolved)
    if err != nil {
        return fmt.Errorf("no such file or directory")
    }

    if !info.IsDir() {
        return fmt.Errorf("not a directory")
    }

    // Add to watcher
    return fw.watcher.Add(resolved)
}
```

**Run Test Again**:
```bash
go test -v -run TestFileWatcherBasicFunctionality
```

**Expected Output** (GREEN):
```
--- PASS: TestFileWatcherBasicFunctionality (0.00s)
    --- PASS: TestFileWatcherBasicFunctionality/watch_valid_directory (0.00s)
    --- PASS: TestFileWatcherBasicFunctionality/watch_nonexistent (0.00s)
PASS
```

### Phase 3: REFACTOR (Clean Up)

**Goal**: Improve code quality without changing behavior.

**Before**:
```go
func (fw *FileWatcherImpl) Watch(path string) error {
    resolved, err := filepath.EvalSymlinks(path)
    if err != nil {
        return fmt.Errorf("no such file or directory")
    }
    info, err := os.Stat(resolved)
    if err != nil {
        return fmt.Errorf("no such file or directory")
    }
    if !info.IsDir() {
        return fmt.Errorf("not a directory")
    }
    return fw.watcher.Add(resolved)
}
```

**After** (Refactored):
```go
func (fw *FileWatcherImpl) Watch(path string) error {
    resolved, err := filepath.EvalSymlinks(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("no such file or directory")
        }
        return err
    }

    info, err := os.Stat(resolved)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("no such file or directory")
        }
        if os.IsPermission(err) {
            return fmt.Errorf("permission denied")
        }
        return err
    }

    if !info.IsDir() {
        return fmt.Errorf("not a directory")
    }

    return fw.watcher.Add(resolved)
}
```

**Run Tests Again** (Ensure GREEN still):
```bash
go test -v -run TestFileWatcherBasicFunctionality
```

---

## Writing New Tests

### Step-by-Step Guide

**1. Identify Feature to Test**

Example: Add "jump to letter" feature to FuzzyFinder

**2. Write Test First (RED)**

```go
func TestFuzzyFinderJumpToLetter(t *testing.T) {
    tests := []struct {
        name         string
        items        []string
        jumpLetter   string
        expectedPos  int
    }{
        {"jump_to_a", []string{"apple", "banana", "cherry"}, "a", 0},
        {"jump_to_b", []string{"apple", "banana", "cherry"}, "b", 1},
        {"jump_no_match", []string{"apple", "banana"}, "z", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ff := NewFuzzyFinderImpl(tt.items)

            // Act
            ff.JumpToLetter(tt.jumpLetter)

            // Assert
            if got := ff.Cursor(); got != tt.expectedPos {
                t.Errorf("got position %d, want %d", got, tt.expectedPos)
            }
        })
    }
}
```

**3. Run Test (Expect Failure)**

```bash
go test -v -run TestFuzzyFinderJumpToLetter
```

**4. Implement Minimal Code (GREEN)**

```go
func (ff *FuzzyFinderImpl) JumpToLetter(letter string) error {
    lowerLetter := strings.ToLower(letter)
    for i, item := range ff.filtered {
        if strings.HasPrefix(strings.ToLower(item), lowerLetter) {
            ff.cursor = i
            return nil
        }
    }
    return nil
}
```

**5. Run Test Again (Expect Pass)**

**6. Refactor if Needed**

**7. Add Edge Cases**

```go
{"jump_case_insensitive", []string{"Apple"}, "a", 0},
{"jump_from_middle", []string{"apple", "banana", "cherry"}, "c", 2},
```

### Test Template

```go
func TestComponentFeature(t *testing.T) {
    tests := []struct {
        name     string
        // Input fields
        input    string
        // Expected fields
        expected string
    }{
        {"descriptive_name", "input", "expected"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange: Set up test data
            component := NewComponent()

            // Act: Execute the function
            result := component.Method(tt.input)

            // Assert: Verify the result
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

---

## Test Fixtures and Mocks

### Test Fixtures

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/test_mocks.go`

```go
// RipgrepMatch represents a search result
type RipgrepMatch struct {
    Path       string
    LineNumber int
    Text       string
}

// RipgrepMessage represents raw ripgrep JSON
type RipgrepMessage struct {
    Type string
    Data json.RawMessage
}

// SearchResult represents a formatted result
type SearchResult struct {
    Path      string
    Highlight string
}

// Style represents element styling
type Style interface {
    GetForeground() string
    GetBackground() string
}

// Viewport represents scrollable content
type Viewport interface {
    SetContent(content string)
    SetYOffset(offset int)
    YOffset() int
    Width() int
    Height() int
    NeedsReflow() bool
}
```

### Using Temporary Directories

```go
func TestFileWatcherDetectChanges(t *testing.T) {
    // Create temporary directory
    tmpDir := t.TempDir()  // Auto-cleaned up after test

    watcher := NewFileWatcher()
    defer watcher.Close()

    watcher.Watch(tmpDir)

    // Create test file
    testFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(testFile, []byte("content"), 0644)

    // Wait for change event
    select {
    case change := <-watcher.Changes():
        if change != testFile {
            t.Errorf("expected %s, got %s", testFile, change)
        }
    case <-time.After(500 * time.Millisecond):
        t.Error("change not detected")
    }
}
```

### Mock Interfaces

```go
// MockRipgrepExecutor for testing without ripgrep
type MockRipgrepExecutor interface {
    Execute(ctx context.Context, query string) (<-chan json.RawMessage, error)
    ExecuteWithLimit(ctx context.Context, query string, max int) (<-chan json.RawMessage, error)
    SetBufferSize(size int)
    SetContextLines(before, after int)
    SetFileTypes(include, exclude []string)
    SetCaseSensitive(sensitive bool)
}

// Usage in tests
func TestSearchWithMock(t *testing.T) {
    mock := &MockRipgrepExecutorImpl{
        results: []RipgrepMatch{
            {Path: "file.txt", LineNumber: 10, Text: "TODO: fix"},
        },
    }

    results, _ := mock.Execute(context.Background(), "TODO")
    // Verify results...
}
```

---

## Debugging Failing Tests

### 1. Read the Error Message

```
--- FAIL: TestFileWatcherDetectChanges (0.50s)
    --- FAIL: TestFileWatcherDetectChanges/detect_file_create (0.50s)
        file_watcher_test.go:81: expected change not detected
```

**Analysis**:
- Test: `TestFileWatcherDetectChanges/detect_file_create`
- Line: `file_watcher_test.go:81`
- Issue: "expected change not detected" (timeout)

### 2. Add Debug Output

```go
func TestFileWatcherDetectChanges(t *testing.T) {
    tmpDir := t.TempDir()
    t.Logf("Watching directory: %s", tmpDir)  // Debug output

    watcher := NewFileWatcher()
    defer watcher.Close()

    watcher.Watch(tmpDir)

    // Create file
    testFile := filepath.Join(tmpDir, "test.txt")
    t.Logf("Creating file: %s", testFile)
    os.WriteFile(testFile, []byte("content"), 0644)

    select {
    case change := <-watcher.Changes():
        t.Logf("Received change: %s", change)
        // Assert...
    case <-time.After(500 * time.Millisecond):
        t.Error("change not detected (timeout after 500ms)")
    }
}
```

### 3. Run Single Test with Verbose Output

```bash
go test -v -run TestFileWatcherDetectChanges/detect_file_create
```

**Output**:
```
=== RUN   TestFileWatcherDetectChanges/detect_file_create
    file_watcher_test.go:60: Watching directory: /tmp/test123
    file_watcher_test.go:67: Creating file: /tmp/test123/test.txt
    file_watcher_test.go:71: Received change: /tmp/test123/test.txt
--- PASS: TestFileWatcherDetectChanges/detect_file_create (0.25s)
```

### 4. Check Goroutine Leaks

```go
func TestFileWatcherClose(t *testing.T) {
    before := runtime.NumGoroutine()
    t.Logf("Goroutines before: %d", before)

    watcher := NewFileWatcher()
    watcher.Watch(t.TempDir())
    watcher.Close()

    time.Sleep(100 * time.Millisecond)  // Wait for cleanup

    after := runtime.NumGoroutine()
    t.Logf("Goroutines after: %d", after)

    if after > before+2 {  // Allow small variance
        t.Errorf("goroutine leak: %d leaked", after-before)
    }
}
```

### 5. Use Race Detector

```bash
go test -race -run TestFileWatcherConcurrent
```

**Output if race detected**:
```
==================
WARNING: DATA RACE
Write at 0x00c0001a0080 by goroutine 7:
  main.(*FileWatcherImpl).handleEvent()
      file_watcher.go:220 +0x123

Previous write at 0x00c0001a0080 by goroutine 6:
  main.(*FileWatcherImpl).SetDebounce()
      file_watcher.go:165 +0x89
```

**Fix**: Add mutex locking.

### 6. Common Debugging Patterns

**Timeout Issues**:
```go
// Increase timeout if test is flaky
case <-time.After(1 * time.Second):  // Was 500ms
```

**Channel Blocking**:
```go
// Add default case to prevent blocking
select {
case ch <- value:
default:
    t.Log("Channel full, dropping event")
}
```

**Filesystem Delays**:
```go
// Wait for filesystem to sync
time.Sleep(50 * time.Millisecond)
```

---

## Integration Tests

### Full Integration Test Structure

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/integration_test.go`

```go
func TestFullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // 1. Watch directory
    watcher := NewFileWatcher()
    defer watcher.Close()

    tmpDir := t.TempDir()
    watcher.Watch(tmpDir)

    // 2. Create markdown file
    mdFile := filepath.Join(tmpDir, "README.md")
    content := "# Test\n\nThis is **bold**."
    os.WriteFile(mdFile, []byte(content), 0644)

    // 3. Wait for change event
    select {
    case change := <-watcher.Changes():
        t.Logf("Detected change: %s", change)
    case <-time.After(1 * time.Second):
        t.Fatal("change not detected")
    }

    // 4. Render with Glamour
    gr := NewGlamourRendererImpl()
    rendered, err := gr.Render(content)
    if err != nil {
        t.Fatalf("render failed: %v", err)
    }

    // 5. Verify output
    if !strings.Contains(rendered, "Test") {
        t.Error("rendered output missing title")
    }
}
```

### Ripgrep Integration Test

**File**: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/ripgrep_integration_test.go`

```go
func TestRipgrepIntegration(t *testing.T) {
    // Create test files
    tmpDir := t.TempDir()
    os.WriteFile(filepath.Join(tmpDir, "file1.md"), []byte("TODO: fix bug"), 0644)
    os.WriteFile(filepath.Join(tmpDir, "file2.md"), []byte("DONE: feature complete"), 0644)

    // Execute search
    rg := NewRipgrepExecutor()
    ctx := context.Background()
    results, err := rg.Execute(ctx, "TODO")
    if err != nil {
        t.Fatal(err)
    }

    // Count matches
    count := 0
    for range results {
        count++
    }

    if count != 1 {
        t.Errorf("expected 1 match, got %d", count)
    }
}
```

---

## Performance Tests

### Benchmark Tests

```go
func BenchmarkFuzzyFinderFilter(b *testing.B) {
    items := make([]string, 10000)
    for i := 0; i < 10000; i++ {
        items[i] = fmt.Sprintf("file_%d.txt", i)
    }

    ff := NewFuzzyFinderImpl(items)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ff.SetFilter("file_5")
    }
}
```

**Run Benchmarks**:
```bash
go test -bench=. -benchmem
```

**Output**:
```
BenchmarkFuzzyFinderFilter-8    10000    115234 ns/op    32768 B/op    5 allocs/op
```

### Performance Constraints

```go
func TestFileWatcherPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping performance test")
    }

    tmpDir := t.TempDir()
    watcher := NewFileWatcher()
    defer watcher.Close()

    watcher.Watch(tmpDir)

    // Measure detection time
    start := time.Now()
    os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("data"), 0644)

    select {
    case <-watcher.Changes():
        duration := time.Since(start)
        if duration > 100*time.Millisecond {
            t.Errorf("detection took %v, want < 100ms", duration)
        }
    case <-time.After(500 * time.Millisecond):
        t.Fatal("detection timeout")
    }
}
```

---

## Test Coverage

### Generate Coverage Report

```bash
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

**Output**:
```
github.com/lumina/ccn/file_watcher.go:82:      Watch           85.7%
github.com/lumina/ccn/file_watcher.go:134:     Unwatch         100.0%
github.com/lumina/ccn/fuzzy_finder_impl.go:30: SetCursor       100.0%
...
total:                                         (statements)    62.3%
```

### View HTML Coverage

```bash
go tool cover -html=coverage.out
```

**Coverage Goals**:
- **Overall**: 60%+
- **Core Logic**: 80%+
- **Edge Cases**: 50%+

---

## Quick Reference

### Common Commands

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestFileWatcher

# Run with coverage
go test -cover

# Run with race detector
go test -race

# Skip slow tests
go test -short

# Run benchmarks
go test -bench=.

# Watch mode
ls *.go | entr -c go test -v
```

### Assertion Patterns

```go
// Equality
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// Error checking
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Boolean
if !condition {
    t.Error("condition should be true")
}

// Substring
if !strings.Contains(output, "expected") {
    t.Errorf("output missing expected substring")
}
```

---

**Guide Version**: 1.0
**Last Updated**: 2025-10-27
**Maintained By**: docs-generator agent
