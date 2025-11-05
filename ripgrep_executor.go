package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// RipgrepExecutorImpl executes ripgrep searches with concurrency control.
// Thread-safe implementation using semaphore for limiting concurrent searches.
type RipgrepExecutorImpl struct {
	maxConcurrent     int
	sem               chan struct{} // Semaphore for concurrency control
	bufferSize        int
	contextBefore     int
	contextAfter      int
	includeTypes      []string
	excludeTypes      []string
	caseSensitive     bool
	ripgrepCheckOnce  sync.Once // Thread-safe initialization
	ripgrepAvailable  bool
	ripgrepCheckError error
}

// NewRipgrepExecutor creates a new ripgrep executor instance.
// Default configuration: 4 concurrent searches, 32KB buffer, case-insensitive.
func NewRipgrepExecutor() MockRipgrepExecutor {
	const defaultConcurrency = 4
	return &RipgrepExecutorImpl{
		maxConcurrent: defaultConcurrency,
		sem:           make(chan struct{}, defaultConcurrency),
		bufferSize:    32 * 1024,
		caseSensitive: false,
	}
}

// Execute runs ripgrep with the given query and returns results via channel.
// Results are streamed asynchronously. The channel is closed when search completes or context is canceled.
// Returns error if query validation fails or ripgrep is not installed.
func (r *RipgrepExecutorImpl) Execute(ctx context.Context, query string) (<-chan json.RawMessage, error) {
	// Validate query FIRST (security check - happens before ripgrep availability)
	if err := r.validateQuery(query); err != nil {
		return nil, err
	}

	// Then check ripgrep availability
	if err := r.checkRipgrepAvailable(); err != nil {
		return nil, err
	}

	out := make(chan json.RawMessage, 100)

	go func() {
		defer close(out)

		// Acquire semaphore
		select {
		case r.sem <- struct{}{}:
			defer func() { <-r.sem }()
		case <-ctx.Done():
			return
		}

		// Build command arguments
		args := []string{"--json"}

		// Case sensitivity
		if !r.caseSensitive {
			args = append(args, "-i")
		}

		// Context lines
		if r.contextBefore > 0 {
			args = append(args, fmt.Sprintf("-B%d", r.contextBefore))
		}
		if r.contextAfter > 0 {
			args = append(args, fmt.Sprintf("-A%d", r.contextAfter))
		}

		// File type filters
		for _, ft := range r.includeTypes {
			args = append(args, "-t", ft)
		}
		for _, ft := range r.excludeTypes {
			args = append(args, "-T", ft)
		}

		args = append(args, query)

		cmd := exec.CommandContext(ctx, "rg", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		if err := cmd.Start(); err != nil {
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, r.bufferSize), 1024*1024)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				// Make a copy of scanner.Bytes() to avoid race condition
				// scanner.Bytes() returns a slice that's reused on next scan
				bytes := make([]byte, len(scanner.Bytes()))
				copy(bytes, scanner.Bytes())
				out <- json.RawMessage(bytes)
			}
		}

		cmd.Wait() // Wait for process completion (ignoring exit code)
	}()

	return out, nil
}

// checkRipgrepAvailable verifies ripgrep is installed and available in PATH.
// Result is cached to avoid repeated checks. Returns helpful error message if missing.
// Thread-safe using sync.Once.
func (r *RipgrepExecutorImpl) checkRipgrepAvailable() error {
	r.ripgrepCheckOnce.Do(func() {
		// Check if ripgrep is in PATH
		_, err := exec.LookPath("rg")
		if err != nil {
			r.ripgrepCheckError = fmt.Errorf("ripgrep not found in PATH: %w\n\nInstall ripgrep:\n  macOS:  brew install ripgrep\n  Linux:  apt-get install ripgrep (Debian/Ubuntu) or yum install ripgrep (RedHat/Fedora)\n  More:   https://github.com/BurntSushi/ripgrep#installation", err)
			r.ripgrepAvailable = false
		} else {
			r.ripgrepAvailable = true
			r.ripgrepCheckError = nil
		}
	})

	return r.ripgrepCheckError
}

// validateQuery checks query for potential shell injection attacks.
func (r *RipgrepExecutorImpl) validateQuery(query string) error {
	if query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	// Prevent shell metacharacters that could enable command injection
	if strings.ContainsAny(query, ";|`$") {
		return fmt.Errorf("invalid query - shell metacharacters not allowed")
	}
	return nil
}

// ExecuteWithLimit runs ripgrep with a maximum result count.
// Automatically stops consuming results after maxResults are reached.
func (r *RipgrepExecutorImpl) ExecuteWithLimit(ctx context.Context, query string, maxResults int) (<-chan json.RawMessage, error) {
	resultChan, err := r.Execute(ctx, query)
	if err != nil {
		return nil, err
	}

	limitedChan := make(chan json.RawMessage, 100)

	go func() {
		defer close(limitedChan)
		count := 0
		for msg := range resultChan {
			if count >= maxResults {
				break // Stop consuming after limit reached
			}
			select {
			case limitedChan <- msg:
				count++
			case <-ctx.Done():
				return
			}
		}
	}()

	return limitedChan, nil
}

// SetBufferSize configures the buffer size for streaming
func (r *RipgrepExecutorImpl) SetBufferSize(size int) {
	r.bufferSize = size
}

// SetContextLines configures before/after context lines
func (r *RipgrepExecutorImpl) SetContextLines(before, after int) {
	r.contextBefore = before
	r.contextAfter = after
}

// SetFileTypes sets file type filters
func (r *RipgrepExecutorImpl) SetFileTypes(include, exclude []string) {
	r.includeTypes = include
	r.excludeTypes = exclude
}

// SetCaseSensitive configures case sensitivity
func (r *RipgrepExecutorImpl) SetCaseSensitive(sensitive bool) {
	r.caseSensitive = sensitive
}

// RipgrepManagerImpl implements RipgrepManager
type RipgrepManagerImpl struct {
	executor      MockRipgrepExecutor
	maxConcurrent int
}

// NewRipgrepManager creates a new ripgrep manager with concurrency limits
func NewRipgrepManager(maxConcurrent int) RipgrepManager {
	executor := NewRipgrepExecutor()
	return &RipgrepManagerImpl{
		executor:      executor,
		maxConcurrent: maxConcurrent,
	}
}

// Search executes a search query and returns parsed matches
func (m *RipgrepManagerImpl) Search(ctx context.Context, query string) (<-chan RipgrepMatch, error) {
	rawChan, err := m.executor.Execute(ctx, query)
	if err != nil {
		return nil, err
	}

	matchChan := make(chan RipgrepMatch, 100)

	go func() {
		defer close(matchChan)
		for raw := range rawChan {
			var msg RipgrepMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}

			if msg.Type == "match" || msg.Type == "context" {
				match := ParseRipgrepMatch(msg)
				select {
				case matchChan <- match:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return matchChan, nil
}

// SetMaxConcurrent sets the maximum concurrent searches
func (m *RipgrepManagerImpl) SetMaxConcurrent(max int) {
	m.maxConcurrent = max
}

// ParseRipgrepMatch parses a ripgrep message into a match
func ParseRipgrepMatch(msg RipgrepMessage) RipgrepMatch {
	var result RipgrepMatch

	// Parse the data field
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return result
	}

	// Extract path
	if pathData, ok := data["path"].(map[string]interface{}); ok {
		if text, ok := pathData["text"].(string); ok {
			result.Path = text
		}
	}

	// Extract line number
	if lineNum, ok := data["line_number"].(float64); ok {
		result.LineNumber = int(lineNum)
	}

	// Extract text
	if linesData, ok := data["lines"].(map[string]interface{}); ok {
		if text, ok := linesData["text"].(string); ok {
			result.Text = text
		}
	}

	return result
}
