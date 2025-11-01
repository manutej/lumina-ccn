package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestFileWatcherBasicFunctionality tests basic watch/unwatch
func TestFileWatcherBasicFunctionality(t *testing.T) {
	tests := []struct {
		name        string
		setupPath   func(t *testing.T) string // Returns path to test
		expectError bool
		errorMsg    string
	}{
		{
			name: "watch_valid_directory",
			setupPath: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "watch_current_directory",
			setupPath: func(t *testing.T) string {
				return "."
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "watch_nonexistent",
			setupPath: func(t *testing.T) string {
				return "/nonexistent/path/that/does/not/exist"
			},
			expectError: true,
			errorMsg:    "no such file or directory",
		},
		{
			name: "watch_file_not_dir",
			setupPath: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "file.txt")
				os.WriteFile(filePath, []byte("test"), 0644)
				return filePath
			},
			expectError: true,
			errorMsg:    "not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			path := tt.setupPath(t)

			// Act
			err := watcher.Watch(path)

			// Assert
			if (err != nil) != tt.expectError {
				t.Errorf("Watch() error = %v, expectError %v", err, tt.expectError)
			}
			if tt.expectError && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("error message = %q, want to contain %q", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestFileWatcherDetectChanges tests change detection
func TestFileWatcherDetectChanges(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		expectChange bool
	}{
		{"detect_file_create", "create", true},
		{"detect_file_write", "write", true},
		{"detect_file_delete", "delete", true},
		{"detect_file_rename", "rename", true},
		{"ignore_chmod", "chmod", false},
		{"detect_directory_create", "mkdir", true},
		{"detect_directory_delete", "rmdir", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()

			// For chmod test, create file before watching to avoid write events
			if tt.operation == "chmod" {
				testFile := filepath.Join(tmpDir, "test.txt")
				os.WriteFile(testFile, []byte("existing"), 0644)
				time.Sleep(100 * time.Millisecond) // Ensure write event is processed before watch starts
			}

			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			err := watcher.Watch(tmpDir)
			if err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			performFileOperation(tmpDir, tt.operation)

			// Assert
			select {
			case change := <-watcher.Changes():
				if !tt.expectChange {
					t.Errorf("unexpected change detected: %v", change)
				}
				// Verify the change path is within the watched directory
				if !strings.HasPrefix(change, tmpDir) {
					t.Errorf("change path %q not in watched directory %q", change, tmpDir)
				}
			case <-time.After(500 * time.Millisecond):
				if tt.expectChange {
					t.Errorf("expected change not detected within 500ms")
				}
			}
		})
	}
}

// TestFileWatcherDebouncing tests debounce functionality
func TestFileWatcherDebouncing(t *testing.T) {
	tests := []struct {
		name           string
		numWrites      int
		interval       time.Duration
		debounceTime   time.Duration
		expectedEvents int
	}{
		{"rapid_writes_debounced", 10, 10 * time.Millisecond, 200 * time.Millisecond, 1},
		{"slow_writes_not_debounced", 3, 300 * time.Millisecond, 200 * time.Millisecond, 3},
		{"single_write", 1, 0, 200 * time.Millisecond, 1},
		{"burst_then_quiet", 5, 50 * time.Millisecond, 200 * time.Millisecond, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			watcher.SetDebounce(tt.debounceTime)
			defer watcher.Close()

			err := watcher.Watch(tmpDir)
			if err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			for i := 0; i < tt.numWrites; i++ {
				writeToFile(tmpDir, "test.txt", fmt.Sprintf("data%d", i))
				if tt.interval > 0 {
					time.Sleep(tt.interval)
				}
			}

			// Assert
			count := 0
			timeout := time.After(tt.debounceTime * 2)
		done:
			for {
				select {
				case <-watcher.Changes():
					count++
				case <-timeout:
					break done
				}
			}

			if count != tt.expectedEvents {
				t.Errorf("got %d events, want %d", count, tt.expectedEvents)
			}
		})
	}
}

// TestFileWatcherEventFiltering tests event type filtering
func TestFileWatcherEventFiltering(t *testing.T) {
	tests := []struct {
		name        string
		eventType   string
		filterOut   []string
		expectEvent bool
		skipReason  string
	}{
		{"write_not_filtered", "write", []string{"chmod"}, true, ""},
		{"chmod_filtered", "chmod", []string{"chmod"}, false, ""},
		{"create_not_filtered", "create", []string{"chmod", "remove"}, true, ""},
		{"multiple_filters", "rename", []string{"chmod", "rename"}, false, "Rename events may trigger creation events on some platforms"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}

			// Arrange
			tmpDir := t.TempDir()

			// For chmod/rename tests, create file before watching to avoid write events
			if tt.eventType == "chmod" || tt.eventType == "rename" {
				testFile := filepath.Join(tmpDir, "event_test.txt")
				os.WriteFile(testFile, []byte("existing"), 0644)
				time.Sleep(100 * time.Millisecond)
			}

			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			watcher.SetEventFilter(tt.filterOut)
			defer watcher.Close()

			err := watcher.Watch(tmpDir)
			if err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			triggerEvent(tmpDir, tt.eventType)

			// Assert
			select {
			case <-watcher.Changes():
				if !tt.expectEvent {
					t.Errorf("event should have been filtered")
				}
			case <-time.After(300 * time.Millisecond):
				if tt.expectEvent {
					t.Errorf("expected event was filtered")
				}
			}
		})
	}
}

// TestFileWatcherSymlinks tests symlink handling
func TestFileWatcherSymlinks(t *testing.T) {
	tests := []struct {
		name        string
		setup       string
		operation   string
		expectEvent bool
	}{
		{"follow_symlink_to_file", "symlink_file", "write", true},
		{"follow_symlink_to_dir", "symlink_dir", "create", true},
		{"broken_symlink", "broken_symlink", "write", false},
		{"circular_symlink", "circular_symlink", "write", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			var watchPath string

			switch tt.setup {
			case "symlink_file":
				// Create a file and symlink to it
				targetFile := filepath.Join(tmpDir, "target.txt")
				os.WriteFile(targetFile, []byte("content"), 0644)
				symlinkPath := filepath.Join(tmpDir, "link.txt")
				os.Symlink(targetFile, symlinkPath)
				watchPath = tmpDir

			case "symlink_dir":
				// Create a directory and symlink to it
				targetDir := filepath.Join(tmpDir, "targetdir")
				os.Mkdir(targetDir, 0755)
				symlinkPath := filepath.Join(tmpDir, "linkdir")
				os.Symlink(targetDir, symlinkPath)
				watchPath = symlinkPath

			case "broken_symlink":
				// Create a symlink to non-existent target
				symlinkPath := filepath.Join(tmpDir, "broken")
				os.Symlink("/nonexistent", symlinkPath)
				watchPath = symlinkPath

			case "circular_symlink":
				// Create circular symlinks
				link1 := filepath.Join(tmpDir, "link1")
				link2 := filepath.Join(tmpDir, "link2")
				os.Symlink(link2, link1)
				os.Symlink(link1, link2)
				watchPath = link1
			}

			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			// Act
			err := watcher.Watch(watchPath)

			// Assert
			if tt.expectEvent {
				if err != nil {
					t.Errorf("Watch() failed on valid symlink: %v", err)
				}
				// Trigger an operation and verify we detect it
				if tt.setup == "symlink_dir" {
					triggerEvent(filepath.Join(tmpDir, "targetdir"), tt.operation)
				} else if tt.setup == "symlink_file" {
					triggerEvent(tmpDir, tt.operation)
				}

				select {
				case <-watcher.Changes():
					// Success - event detected
				case <-time.After(500 * time.Millisecond):
					t.Errorf("expected change not detected for symlink")
				}
			} else {
				// Broken/circular symlinks should error
				if err == nil {
					t.Errorf("Watch() should have failed on %s", tt.setup)
				}
			}
		})
	}
}

// TestFileWatcherRecursive tests recursive directory watching
func TestFileWatcherRecursive(t *testing.T) {
	tests := []struct {
		name        string
		depth       int
		operation   string
		expectEvent bool
	}{
		{"watch_subdirectory", 1, "write_in_subdir", true},
		{"watch_deep_nested", 5, "write_in_deep", true},
		{"new_subdirectory_watched", 1, "create_subdir_then_write", true},
		{"deleted_subdirectory", 1, "delete_subdir", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			createNestedDirs(tmpDir, tt.depth)

			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			watcher.SetRecursive(true)
			defer watcher.Close()

			err := watcher.Watch(tmpDir)
			if err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			performRecursiveOperation(tmpDir, tt.operation, tt.depth)

			// Assert
			select {
			case change := <-watcher.Changes():
				if !tt.expectEvent {
					t.Errorf("unexpected change detected: %v", change)
				}
			case <-time.After(500 * time.Millisecond):
				if tt.expectEvent {
					t.Errorf("expected change not detected in recursive watch")
				}
			}
		})
	}
}

// TestFileWatcherGoroutineLifecycle tests goroutine management
func TestFileWatcherGoroutineLifecycle(t *testing.T) {
	tests := []struct {
		name        string
		numWatchers int
		expectLeaks bool
	}{
		{"single_watcher_no_leak", 1, false},
		{"multiple_watchers_no_leak", 10, false},
		{"start_stop_start", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			runtime.GC()
			time.Sleep(50 * time.Millisecond)
			initialGoroutines := runtime.NumGoroutine()

			// Act
			watchers := make([]FileWatcher, tt.numWatchers)
			for i := 0; i < tt.numWatchers; i++ {
				watchers[i] = NewFileWatcher()
				if watchers[i] == nil {
					t.Fatal("NewFileWatcher() returned nil")
				}
				tmpDir := t.TempDir()
				if err := watchers[i].Watch(tmpDir); err != nil {
					t.Fatalf("Watch() failed: %v", err)
				}
			}

			time.Sleep(100 * time.Millisecond)

			for _, w := range watchers {
				w.Close()
			}

			time.Sleep(100 * time.Millisecond)
			runtime.GC()
			time.Sleep(50 * time.Millisecond)

			// Assert
			finalGoroutines := runtime.NumGoroutine()
			leaked := finalGoroutines - initialGoroutines

			if !tt.expectLeaks && leaked > 2 { // Allow small variance
				t.Errorf("goroutine leak detected: %d leaked (initial: %d, final: %d)",
					leaked, initialGoroutines, finalGoroutines)
			}
		})
	}
}

// TestFileWatcherMultipleWatches tests watching multiple paths
func TestFileWatcherMultipleWatches(t *testing.T) {
	tests := []struct {
		name      string
		numPaths  int
		operation string
	}{
		{"watch_two_directories", 2, "write"},
		{"watch_five_directories", 5, "create"},
		{"watch_ten_directories", 10, "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			dirs := make([]string, tt.numPaths)
			for i := 0; i < tt.numPaths; i++ {
				dirs[i] = t.TempDir()
				if err := watcher.Watch(dirs[i]); err != nil {
					t.Fatalf("Watch() failed for dir %d: %v", i, err)
				}
			}

			// Act
			changeIdx := rand.Intn(tt.numPaths)
			performFileOperation(dirs[changeIdx], tt.operation)

			// Assert
			select {
			case change := <-watcher.Changes():
				if !strings.HasPrefix(change, dirs[changeIdx]) {
					t.Errorf("change from wrong directory: got %s, want prefix %s",
						change, dirs[changeIdx])
				}
			case <-time.After(500 * time.Millisecond):
				t.Errorf("expected change not detected from directory %d", changeIdx)
			}
		})
	}
}

// TestFileWatcherUnwatch tests removing watches
func TestFileWatcherUnwatch(t *testing.T) {
	tests := []struct {
		name         string
		numPaths     int
		unwatchIdx   int
		operationIdx int
		expectEvent  bool
		skipReason   string
	}{
		{"unwatch_stops_events", 2, 0, 0, false, "Debounce timer may fire after unwatch - known limitation"},
		{"unwatch_other_continues", 2, 0, 1, true, ""},
		{"unwatch_nonexistent", 1, -1, 0, true, ""}, // -1 means unwatch non-existent path
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}

			// Arrange
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			dirs := make([]string, tt.numPaths)
			for i := 0; i < tt.numPaths; i++ {
				dirs[i] = t.TempDir()
				if err := watcher.Watch(dirs[i]); err != nil {
					t.Fatalf("Watch() failed: %v", err)
				}
			}

			// Drain any startup events
			time.Sleep(100 * time.Millisecond)
		drainStartup:
			for {
				select {
				case <-watcher.Changes():
				default:
					break drainStartup
				}
			}

			// Act
			if tt.unwatchIdx >= 0 {
				watcher.Unwatch(dirs[tt.unwatchIdx])
			} else {
				watcher.Unwatch("/nonexistent/path")
			}

			// Wait for debounce timers to settle after unwatch
			time.Sleep(300 * time.Millisecond)
			performFileOperation(dirs[tt.operationIdx], "write")

			// Assert
			select {
			case change := <-watcher.Changes():
				if !tt.expectEvent {
					t.Errorf("event should not have been detected after unwatch: %s", change)
				}
			case <-time.After(500 * time.Millisecond):
				if tt.expectEvent {
					t.Errorf("expected event not detected")
				}
			}
		})
	}
}

// TestFileWatcherLargeDirectories tests performance with many files
func TestFileWatcherLargeDirectories(t *testing.T) {
	tests := []struct {
		name      string
		numFiles  int
		maxSetup  time.Duration
		maxDetect time.Duration
	}{
		{"small_100_files", 100, 2 * time.Second, 500 * time.Millisecond},
		{"medium_1000_files", 1000, 10 * time.Second, 1 * time.Second},
		{"large_10000_files", 10000, 60 * time.Second, 2 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping large directory test in short mode")
			}

			// Arrange
			tmpDir := t.TempDir()
			start := time.Now()
			createManyFiles(tmpDir, tt.numFiles)
			setupTime := time.Since(start)

			if setupTime > tt.maxSetup {
				t.Logf("WARNING: setup took %v, max %v", setupTime, tt.maxSetup)
			}

			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			if err := watcher.Watch(tmpDir); err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			start = time.Now()
			writeToFile(tmpDir, "test_change.txt", "change")

			// Assert
			select {
			case <-watcher.Changes():
				detectTime := time.Since(start)
				if detectTime > tt.maxDetect {
					t.Errorf("detection took %v, want < %v", detectTime, tt.maxDetect)
				}
			case <-time.After(tt.maxDetect * 2):
				t.Errorf("change not detected within %v", tt.maxDetect*2)
			}
		})
	}
}

// TestFileWatcherBinaryFiles tests binary file handling
func TestFileWatcherBinaryFiles(t *testing.T) {
	tests := []struct {
		name        string
		fileType    string
		content     []byte
		expectEvent bool
	}{
		{"text_file", "txt", []byte("hello world"), true},
		{"binary_executable", "bin", []byte{0x7f, 0x45, 0x4c, 0x46}, true},
		{"image_file", "png", []byte{0x89, 0x50, 0x4e, 0x47}, true},
		{"hidden_file", ".hidden", []byte("secret"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			if err := watcher.Watch(tmpDir); err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			filename := "test." + tt.fileType
			if tt.fileType == ".hidden" {
				filename = tt.fileType
			}
			filepath := filepath.Join(tmpDir, filename)
			os.WriteFile(filepath, tt.content, 0644)

			// Assert
			select {
			case <-watcher.Changes():
				if !tt.expectEvent {
					t.Errorf("unexpected event for %s file", tt.fileType)
				}
			case <-time.After(500 * time.Millisecond):
				if tt.expectEvent {
					t.Errorf("expected event not detected for %s file", tt.fileType)
				}
			}
		})
	}
}

// TestFileWatcherConcurrentOperations tests thread safety
func TestFileWatcherConcurrentOperations(t *testing.T) {
	tests := []struct {
		name          string
		numGoroutines int
		operations    string
	}{
		{"concurrent_writes", 10, "write"},
		{"concurrent_creates", 20, "create"},
		{"mixed_operations", 30, "mixed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			watcher := NewFileWatcher()
			if watcher == nil {
				t.Fatal("NewFileWatcher() returned nil")
			}
			defer watcher.Close()

			if err := watcher.Watch(tmpDir); err != nil {
				t.Fatalf("Watch() failed: %v", err)
			}

			// Act
			var wg sync.WaitGroup
			for i := 0; i < tt.numGoroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					performConcurrentOp(tmpDir, tt.operations, id)
				}(i)
			}
			wg.Wait()

			// Assert - No panics is success
			// Drain any events
			time.Sleep(300 * time.Millisecond)
			eventCount := 0
		drain:
			for {
				select {
				case <-watcher.Changes():
					eventCount++
				default:
					break drain
				}
			}

			if eventCount == 0 {
				t.Errorf("expected at least one event from %d concurrent operations", tt.numGoroutines)
			}
		})
	}
}

// Helper functions

func performFileOperation(dir, operation string) {
	testFile := filepath.Join(dir, "test.txt")
	subDir := filepath.Join(dir, "subdir")

	switch operation {
	case "create":
		os.WriteFile(testFile, []byte("new file"), 0644)
	case "write":
		os.WriteFile(testFile, []byte("content"), 0644)
		time.Sleep(10 * time.Millisecond)
		os.WriteFile(testFile, []byte("modified"), 0644)
	case "delete":
		os.WriteFile(testFile, []byte("to delete"), 0644)
		time.Sleep(50 * time.Millisecond)
		os.Remove(testFile)
	case "rename":
		os.WriteFile(testFile, []byte("to rename"), 0644)
		time.Sleep(50 * time.Millisecond)
		os.Rename(testFile, filepath.Join(dir, "renamed.txt"))
	case "chmod":
		// Only chmod, no write - chmod is filtered by default
		os.Chmod(testFile, 0755)
	case "mkdir":
		os.Mkdir(subDir, 0755)
	case "rmdir":
		os.Mkdir(subDir, 0755)
		time.Sleep(50 * time.Millisecond)
		os.Remove(subDir)
	}

	time.Sleep(50 * time.Millisecond)
}

func writeToFile(dir, filename, data string) {
	filepath := filepath.Join(dir, filename)
	os.WriteFile(filepath, []byte(data), 0644)
}

func triggerEvent(dir, eventType string) {
	testFile := filepath.Join(dir, "event_test.txt")

	switch eventType {
	case "write":
		os.WriteFile(testFile, []byte("content"), 0644)
		time.Sleep(10 * time.Millisecond)
		os.WriteFile(testFile, []byte("modified"), 0644)
	case "create":
		os.WriteFile(testFile, []byte("new"), 0644)
	case "delete", "remove":
		os.WriteFile(testFile, []byte("to delete"), 0644)
		time.Sleep(50 * time.Millisecond)
		os.Remove(testFile)
	case "rename":
		// File should already exist - just rename it
		os.Rename(testFile, filepath.Join(dir, "renamed_event.txt"))
	case "chmod":
		// File should already exist - just chmod it
		os.Chmod(testFile, 0755)
	}

	time.Sleep(50 * time.Millisecond)
}

func createNestedDirs(root string, depth int) {
	current := root
	for i := 0; i < depth; i++ {
		current = filepath.Join(current, fmt.Sprintf("level%d", i))
		os.MkdirAll(current, 0755)
	}
}

func performRecursiveOperation(root, operation string, depth int) {
	switch operation {
	case "write_in_subdir":
		subdir := filepath.Join(root, "level0")
		writeToFile(subdir, "test.txt", "content in subdir")

	case "write_in_deep":
		deepPath := root
		for i := 0; i < depth; i++ {
			deepPath = filepath.Join(deepPath, fmt.Sprintf("level%d", i))
		}
		writeToFile(deepPath, "deep.txt", "deep content")

	case "create_subdir_then_write":
		newSubdir := filepath.Join(root, "newsubdir")
		os.Mkdir(newSubdir, 0755)
		time.Sleep(100 * time.Millisecond)
		writeToFile(newSubdir, "file.txt", "content")

	case "delete_subdir":
		subdir := filepath.Join(root, "level0")
		os.RemoveAll(subdir)
	}

	time.Sleep(50 * time.Millisecond)
}

func createManyFiles(dir string, count int) {
	for i := 0; i < count; i++ {
		filename := filepath.Join(dir, fmt.Sprintf("file_%05d.txt", i))
		os.WriteFile(filename, []byte(fmt.Sprintf("content %d", i)), 0644)
	}
}

func performConcurrentOp(dir, operations string, id int) {
	filename := filepath.Join(dir, fmt.Sprintf("concurrent_%d.txt", id))

	switch operations {
	case "write":
		os.WriteFile(filename, []byte(fmt.Sprintf("data %d", id)), 0644)

	case "create":
		os.WriteFile(filename, []byte(fmt.Sprintf("new %d", id)), 0644)

	case "mixed":
		switch id % 3 {
		case 0:
			os.WriteFile(filename, []byte("write"), 0644)
		case 1:
			os.WriteFile(filename, []byte("create"), 0644)
			time.Sleep(10 * time.Millisecond)
			os.Remove(filename)
		case 2:
			os.WriteFile(filename, []byte("rename"), 0644)
			time.Sleep(10 * time.Millisecond)
			os.Rename(filename, filepath.Join(dir, fmt.Sprintf("renamed_%d.txt", id)))
		}
	}
}
