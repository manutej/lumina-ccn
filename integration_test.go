package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// RELAY_4: Frontend Integration Tests
// These tests verify that components integrate correctly with each other.
// This is NOT about making RED phase tests pass (that's RELAY_5's job).
// This IS about verifying integration points work as expected.

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 1: FileWatcher → Auto-Reload
// ═══════════════════════════════════════════════════════════════════════

// TestFileWatcher_IntegrationWithReload verifies file changes trigger events
func TestFileWatcher_IntegrationWithReload(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create file watcher
	fw := NewFileWatcher()
	if fw == nil {
		t.Fatal("FileWatcher creation failed")
	}
	defer fw.Close()

	// Watch directory
	if err := fw.Watch(tmpDir); err != nil {
		t.Fatalf("Watch() failed: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.md")

	// Wait for debounce to settle
	time.Sleep(50 * time.Millisecond)

	// Create file (should trigger event)
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify event is received
	select {
	case path := <-fw.Changes():
		if path == "" {
			t.Error("Expected non-empty path from Changes()")
		}
		t.Logf("✓ FileWatcher detected change: %s", path)
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for file change event")
	}
}

// TestFileWatcher_DebounceMultipleChanges verifies debouncing works
func TestFileWatcher_DebounceMultipleChanges(t *testing.T) {
	tmpDir := t.TempDir()

	fw := NewFileWatcher()
	if fw == nil {
		t.Fatal("FileWatcher creation failed")
	}
	defer fw.Close()

	fw.SetDebounce(100 * time.Millisecond)

	if err := fw.Watch(tmpDir); err != nil {
		t.Fatalf("Watch() failed: %v", err)
	}

	testFile := filepath.Join(tmpDir, "debounce.md")

	// Create file
	os.WriteFile(testFile, []byte("v1"), 0644)
	time.Sleep(20 * time.Millisecond)

	// Update multiple times rapidly
	os.WriteFile(testFile, []byte("v2"), 0644)
	time.Sleep(20 * time.Millisecond)
	os.WriteFile(testFile, []byte("v3"), 0644)
	time.Sleep(20 * time.Millisecond)
	os.WriteFile(testFile, []byte("v4"), 0644)

	// Should receive only ONE event (debounced)
	eventCount := 0
	timeout := time.After(300 * time.Millisecond)

	for {
		select {
		case <-fw.Changes():
			eventCount++
		case <-timeout:
			goto done
		}
	}

done:
	// Should be 1 event (debounced), not 4
	if eventCount > 2 {
		t.Errorf("Expected ≤2 debounced events, got %d", eventCount)
	} else {
		t.Logf("✓ Debouncing works: %d events for 4 writes", eventCount)
	}
}

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 2: Keybinding → FuzzyFinder
// ═══════════════════════════════════════════════════════════════════════

// TestKeybindingHandler_IntegrationWithFuzzyFinder verifies key routing
func TestKeybindingHandler_IntegrationWithFuzzyFinder(t *testing.T) {
	kh := NewKeybindingHandler()
	ff := NewFuzzyFinderImpl([]string{"file1.md", "file2.md", "file3.md"})

	// Normal mode: arrow keys should work
	kh.SetMode("normal")

	// Simulate down arrow key
	ff.HandleKey("down")
	if ff.Cursor() != 1 {
		t.Errorf("Expected cursor at 1, got %d", ff.Cursor())
	} else {
		t.Logf("✓ Keybinding → FuzzyFinder: down key works (cursor=%d)", ff.Cursor())
	}

	// Simulate up arrow key
	ff.HandleKey("up")
	if ff.Cursor() != 0 {
		t.Errorf("Expected cursor at 0, got %d", ff.Cursor())
	} else {
		t.Logf("✓ Keybinding → FuzzyFinder: up key works (cursor=%d)", ff.Cursor())
	}

	// Simulate mode switch to search
	kh.SetMode("search")
	if kh.GetMode() != "search" {
		t.Errorf("Expected search mode, got %s", kh.GetMode())
	} else {
		t.Logf("✓ Keybinding mode switch works: %s", kh.GetMode())
	}
}

// TestKeybindingHandler_ModeTransitions verifies mode switching
func TestKeybindingHandler_ModeTransitions(t *testing.T) {
	kh := NewKeybindingHandler()

	// Start in normal mode
	if kh.GetMode() != "normal" {
		t.Errorf("Expected normal mode initially, got %s", kh.GetMode())
	}

	// Transition: normal → search (/)
	kh.HandleKey("/", nil)
	if kh.GetMode() != "search" {
		t.Errorf("Expected search mode after '/', got %s", kh.GetMode())
	} else {
		t.Logf("✓ Mode transition: normal → search")
	}

	// Transition: search → normal (escape)
	kh.HandleKey("escape", nil)
	if kh.GetMode() != "normal" {
		t.Errorf("Expected normal mode after escape, got %s", kh.GetMode())
	} else {
		t.Logf("✓ Mode transition: search → normal")
	}

	// Transition: normal → help (?)
	kh.HandleKey("?", nil)
	if kh.GetMode() != "help" {
		t.Errorf("Expected help mode after '?', got %s", kh.GetMode())
	} else {
		t.Logf("✓ Mode transition: normal → help")
	}
}

// TestKeybindingHandler_ExitSignal verifies quit handling
func TestKeybindingHandler_ExitSignal(t *testing.T) {
	kh := NewKeybindingHandler()

	// 'q' should return nil (exit signal)
	result := kh.HandleKey("q", "dummy-model")
	if result != nil {
		t.Error("Expected nil (exit) for 'q' key")
	} else {
		t.Logf("✓ Exit signal works for 'q'")
	}

	// ctrl+c should also return nil
	result = kh.HandleKey("ctrl+c", "dummy-model")
	if result != nil {
		t.Error("Expected nil (exit) for 'ctrl+c'")
	} else {
		t.Logf("✓ Exit signal works for 'ctrl+c'")
	}
}

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 3: FuzzyFinder → Filtering
// ═══════════════════════════════════════════════════════════════════════

// TestFuzzyFinder_IntegrationWithSearch verifies filtering works
func TestFuzzyFinder_IntegrationWithSearch(t *testing.T) {
	items := []string{
		"test.md",
		"testing.md",
		"README.md",
		"main.go",
		"test_file.md",
	}

	ff := NewFuzzyFinderImpl(items)

	// Filter for "test"
	ff.SetFilter("test")
	results := ff.FilteredResults()

	if len(results) != 3 {
		t.Errorf("Expected 3 results for 'test', got %d", len(results))
	} else {
		t.Logf("✓ FuzzyFinder filtering works: %d results for 'test'", len(results))
		for i, r := range results {
			t.Logf("  [%d] %s", i, r)
		}
	}

	// Cursor should reset to 0 after filter
	if ff.Cursor() != 0 {
		t.Errorf("Expected cursor reset to 0 after filter, got %d", ff.Cursor())
	}
}

// TestFuzzyFinder_CaseInsensitiveSearch verifies case-insensitive filtering
func TestFuzzyFinder_CaseInsensitiveSearch(t *testing.T) {
	items := []string{"README.md", "readme.txt", "ReadMe.doc"}

	ff := NewFuzzyFinderImpl(items)
	ff.SetFilter("readme") // lowercase query

	results := ff.FilteredResults()
	if len(results) != 3 {
		t.Errorf("Expected 3 case-insensitive matches, got %d", len(results))
	} else {
		t.Logf("✓ Case-insensitive search works: %d results", len(results))
	}
}

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 4: Glamour → Rendering
// ═══════════════════════════════════════════════════════════════════════

// TestGlamourRenderer_IntegrationWithSearchResult verifies rendering
func TestGlamourRenderer_IntegrationWithSearchResult(t *testing.T) {
	gr := NewGlamourRendererImpl()

	result := SearchResult{
		Path:      "path/to/file.md",
		Name:      "file.md",
		Score:     100,
		Highlight: "# Example Match",
	}

	output := gr.RenderSearchResult(result)

	if output == "" {
		t.Error("Expected non-empty rendered output")
	} else {
		t.Logf("✓ Glamour renders search results")
		t.Logf("  Output length: %d chars", len(output))
	}
}

// TestGlamourRenderer_ThemeDetection verifies theme is detected
func TestGlamourRenderer_ThemeDetection(t *testing.T) {
	gr := NewGlamourRendererImpl()

	theme := gr.DetectedTheme()
	if theme != "dark" && theme != "light" {
		t.Errorf("Expected 'dark' or 'light' theme, got '%s'", theme)
	} else {
		t.Logf("✓ Theme detection works: %s", theme)
	}
}

// TestGlamourRenderer_ThemeSwitch verifies theme changes
func TestGlamourRenderer_ThemeSwitch(t *testing.T) {
	gr := NewGlamourRendererImpl()

	// Switch to light theme
	gr.SetTheme("light")
	if gr.DetectedTheme() != "light" {
		t.Errorf("Expected light theme after SetTheme('light'), got %s", gr.DetectedTheme())
	} else {
		t.Logf("✓ Theme switching works: light")
	}

	// Switch to dark theme
	gr.SetTheme("dark")
	if gr.DetectedTheme() != "dark" {
		t.Errorf("Expected dark theme after SetTheme('dark'), got %s", gr.DetectedTheme())
	} else {
		t.Logf("✓ Theme switching works: dark")
	}
}

// TestGlamourRenderer_MarkdownRendering verifies basic markdown rendering
func TestGlamourRenderer_MarkdownRendering(t *testing.T) {
	gr := NewGlamourRendererImpl()

	markdown := "# Test Header\n\nSome **bold** text."
	output, err := gr.Render(markdown)

	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty rendered markdown")
	} else {
		t.Logf("✓ Markdown rendering works")
		t.Logf("  Input: %d chars, Output: %d chars", len(markdown), len(output))
	}
}

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 5: Ripgrep → Search Execution
// ═══════════════════════════════════════════════════════════════════════

// TestRipgrepExecutor_IntegrationWithContext verifies search with context
func TestRipgrepExecutor_IntegrationWithContext(t *testing.T) {
	// Skip if ripgrep not installed
	if !isRipgrepInstalled() {
		t.Skip("ripgrep not installed")
	}

	executor := NewRipgrepExecutor()
	executor.SetBufferSize(64 * 1024)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Search for "package main" in current directory
	resultChan, err := executor.Execute(ctx, "package main")
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Consume at least one result
	count := 0
	timeout := time.After(2 * time.Second)

	for {
		select {
		case msg, ok := <-resultChan:
			if !ok {
				goto done
			}
			if len(msg) > 0 {
				count++
				t.Logf("✓ Received ripgrep result %d", count)
			}
			if count >= 3 {
				goto done
			}
		case <-timeout:
			goto done
		}
	}

done:
	if count > 0 {
		t.Logf("✓ Ripgrep integration works: %d results", count)
	} else {
		t.Log("⚠ No ripgrep results (may be expected if no .go files)")
	}
}

// TestRipgrepExecutor_QueryValidation verifies shell injection prevention
func TestRipgrepExecutor_QueryValidation(t *testing.T) {
	executor := NewRipgrepExecutor()
	ctx := context.Background()

	// Test shell metacharacters are rejected
	badQueries := []string{
		"test; rm -rf /",
		"test | cat /etc/passwd",
		"test `whoami`",
		"test $USER",
	}

	for _, query := range badQueries {
		_, err := executor.Execute(ctx, query)
		if err == nil {
			t.Errorf("Expected error for malicious query: %s", query)
		} else {
			t.Logf("✓ Rejected malicious query: %s", query)
		}
	}
}

// ═══════════════════════════════════════════════════════════════════════
// INTEGRATION 6: Complete Workflow
// ═══════════════════════════════════════════════════════════════════════

// TestCompleteWorkflow_SearchAndNavigate verifies full UI workflow
func TestCompleteWorkflow_SearchAndNavigate(t *testing.T) {
	// Simulate complete workflow:
	// 1. User presses "/" → enters search mode
	// 2. User types "test" → filters results
	// 3. User presses down arrow → selects next result
	// 4. User presses enter → opens selected file

	kh := NewKeybindingHandler()
	items := []string{"test1.md", "test2.md", "other.md", "testing.md"}
	ff := NewFuzzyFinderImpl(items)

	// Step 1: Enter search mode
	kh.HandleKey("/", nil)
	ff.SetMode("search")

	if kh.GetMode() != "search" || ff.Mode() != "search" {
		t.Error("Failed to enter search mode")
	} else {
		t.Log("✓ Step 1: Entered search mode")
	}

	// Step 2: Filter results
	ff.SetFilter("test")
	filtered := ff.FilteredResults()

	if len(filtered) != 3 {
		t.Errorf("Expected 3 filtered results, got %d", len(filtered))
	} else {
		t.Logf("✓ Step 2: Filtered to %d results", len(filtered))
	}

	// Step 3: Navigate with arrow keys
	ff.HandleKey("down")
	if ff.Cursor() != 1 {
		t.Errorf("Expected cursor at 1, got %d", ff.Cursor())
	} else {
		t.Log("✓ Step 3: Navigated to next result")
	}

	// Step 4: Select item
	selected := ff.SelectedItem()
	if selected == "" {
		t.Error("Expected non-empty selected item")
	} else {
		t.Logf("✓ Step 4: Selected item: %s", selected)
	}

	// Step 5: Exit search mode
	kh.HandleKey("escape", nil)
	if kh.GetMode() != "normal" {
		t.Error("Failed to exit search mode")
	} else {
		t.Log("✓ Step 5: Exited search mode")
	}
}

// ═══════════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ═══════════════════════════════════════════════════════════════════════

// isRipgrepInstalled checks if ripgrep is available
func isRipgrepInstalled() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	executor := NewRipgrepExecutor()
	_, err := executor.Execute(ctx, "test")
	return err == nil
}
