package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher monitors file system changes with debouncing and filtering.
// All methods are safe for concurrent use from multiple goroutines.
type FileWatcher interface {
	// Watch adds a path to the watch list. Returns error if path doesn't exist or is not a directory.
	Watch(path string) error

	// Unwatch removes a path from the watch list.
	Unwatch(path string) error

	// Changes returns a read-only channel that receives change notifications.
	Changes() <-chan string

	// Close stops the file watcher and releases all resources.
	Close() error

	// SetDebounce configures the debounce duration for file events.
	SetDebounce(duration time.Duration)

	// SetEventFilter specifies which event types to filter out (e.g., "chmod", "write").
	SetEventFilter(filterOut []string)

	// SetRecursive enables/disables recursive directory watching.
	SetRecursive(recursive bool)
}

// FileWatcherImpl implements FileWatcher using fsnotify.
// Thread-safe implementation with debouncing and event filtering.
type FileWatcherImpl struct {
	watcher       *fsnotify.Watcher
	changes       chan string
	debounce      time.Duration
	debounceTimer *time.Timer
	debounceMu    sync.Mutex
	pendingPath   string
	closed        bool
	closeMu       sync.Mutex
	wg            sync.WaitGroup
	eventFilter   map[string]bool
	filterMu      sync.RWMutex // Use RWMutex for read-heavy operations
	recursive     bool
}

// NewFileWatcher creates a new FileWatcher instance.
// Returns nil if fsnotify initialization fails.
// Default configuration: 200ms debounce, chmod events filtered, recursive watching enabled.
func NewFileWatcher() FileWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}

	fw := &FileWatcherImpl{
		watcher:     watcher,
		changes:     make(chan string, 100),
		debounce:    200 * time.Millisecond,
		eventFilter: make(map[string]bool),
		recursive:   true,
	}

	// Filter chmod events by default to reduce noise
	fw.eventFilter["chmod"] = true

	fw.wg.Add(1)
	go fw.eventLoop()

	return fw
}

// Watch adds a path to the watch list
func (fw *FileWatcherImpl) Watch(path string) error {
	// Resolve symlinks with loop detection
	resolved, err := fw.followSymlinksWithLoopDetection(path)
	if err != nil {
		// Check if it's a broken symlink or doesn't exist
		if os.IsNotExist(err) {
			return fmt.Errorf("no such file or directory")
		}
		return err
	}

	// Check if path exists and is a directory
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

	// Add to watcher
	if err := fw.watcher.Add(resolved); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied")
		}
		return err
	}

	// Recursively watch subdirectories if enabled
	if fw.recursive {
		return filepath.Walk(resolved, func(walkPath string, walkInfo os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return nil // Skip paths with errors
			}
			if walkInfo.IsDir() && walkPath != resolved {
				fw.watcher.Add(walkPath) // Ignore errors for subdirs
			}
			return nil
		})
	}

	return nil
}

// Unwatch removes a path from the watch list
func (fw *FileWatcherImpl) Unwatch(path string) error {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		resolved = path
	}
	return fw.watcher.Remove(resolved)
}

// Changes returns a channel that receives change notifications
func (fw *FileWatcherImpl) Changes() <-chan string {
	return fw.changes
}

// Close stops the file watcher and releases resources
func (fw *FileWatcherImpl) Close() error {
	fw.closeMu.Lock()
	if fw.closed {
		fw.closeMu.Unlock()
		return nil
	}
	fw.closed = true
	fw.closeMu.Unlock()

	err := fw.watcher.Close()
	fw.wg.Wait()
	close(fw.changes)
	return err
}

// SetDebounce sets the debounce duration
func (fw *FileWatcherImpl) SetDebounce(duration time.Duration) {
	fw.debounceMu.Lock()
	defer fw.debounceMu.Unlock()
	fw.debounce = duration
}

// SetEventFilter sets which event types to filter out
func (fw *FileWatcherImpl) SetEventFilter(filterOut []string) {
	fw.filterMu.Lock()
	defer fw.filterMu.Unlock()
	fw.eventFilter = make(map[string]bool)
	for _, eventType := range filterOut {
		fw.eventFilter[eventType] = true
	}
}

// SetRecursive enables/disables recursive directory watching
func (fw *FileWatcherImpl) SetRecursive(recursive bool) {
	fw.recursive = recursive
}

// eventLoop processes fsnotify events
func (fw *FileWatcherImpl) eventLoop() {
	defer fw.wg.Done()

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			fw.handleEvent(event)
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			// Log errors but continue
			_ = err
		}
	}
}

// handleEvent processes a single fsnotify event
func (fw *FileWatcherImpl) handleEvent(event fsnotify.Event) {
	// Check event filter (use RLock for read-only access)
	if fw.shouldFilterEvent(event) {
		return
	}

	// Handle directory creation - add to watch list if recursive
	if fw.recursive && event.Op&fsnotify.Create == fsnotify.Create {
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			fw.watcher.Add(event.Name)
		}
	}

	// Debounce the event
	fw.debounceMu.Lock()
	defer fw.debounceMu.Unlock()

	fw.pendingPath = event.Name

	if fw.debounceTimer != nil {
		fw.debounceTimer.Stop()
	}

	fw.debounceTimer = time.AfterFunc(fw.debounce, func() {
		fw.debounceMu.Lock()
		path := fw.pendingPath
		fw.pendingPath = ""
		fw.debounceMu.Unlock()

		if path != "" {
			// Check if watcher is closed before sending
			fw.closeMu.Lock()
			closed := fw.closed
			fw.closeMu.Unlock()

			if !closed {
				select {
				case fw.changes <- path:
				default:
					// Channel full, skip this event
				}
			}
		}
	})
}

// shouldFilterEvent checks if an event should be filtered based on event type.
// Uses RLock for efficient concurrent reads.
func (fw *FileWatcherImpl) shouldFilterEvent(event fsnotify.Event) bool {
	fw.filterMu.RLock()
	defer fw.filterMu.RUnlock()

	if event.Op&fsnotify.Chmod == fsnotify.Chmod && fw.eventFilter["chmod"] {
		return true
	}
	if event.Op&fsnotify.Write == fsnotify.Write && fw.eventFilter["write"] {
		return true
	}
	if event.Op&fsnotify.Create == fsnotify.Create && fw.eventFilter["create"] {
		return true
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove && fw.eventFilter["remove"] {
		return true
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename && fw.eventFilter["rename"] {
		return true
	}
	return false
}

// followSymlinksWithLoopDetection resolves symlinks while detecting circular references.
// Returns error if symlink loop detected or max depth (40) exceeded.
// Uses inode tracking to detect cycles efficiently.
func (fw *FileWatcherImpl) followSymlinksWithLoopDetection(path string) (string, error) {
	const maxHops = 40 // Linux MAXSYMLINKS is typically 40
	seenInodes := make(map[uint64]bool)
	currentPath := path

	for i := 0; i < maxHops; i++ {
		// Get file info without following symlink
		linkInfo, err := os.Lstat(currentPath)
		if err != nil {
			return "", err
		}

		// If not a symlink, we're done
		if linkInfo.Mode()&os.ModeSymlink == 0 {
			return currentPath, nil
		}

		// Track inode to detect loops
		if stat, ok := linkInfo.Sys().(*syscall.Stat_t); ok {
			inode := stat.Ino
			if seenInodes[inode] {
				return "", fmt.Errorf("symlink loop detected at %s", currentPath)
			}
			seenInodes[inode] = true
		}

		// Read the symlink target
		target, err := os.Readlink(currentPath)
		if err != nil {
			return "", err
		}

		// Handle relative vs absolute symlink targets
		if !filepath.IsAbs(target) {
			currentPath = filepath.Join(filepath.Dir(currentPath), target)
		} else {
			currentPath = target
		}
	}

	return "", fmt.Errorf("too many symlink levels (>%d) for %s", maxHops, path)
}
