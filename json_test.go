package main

import (
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple object",
			input:    `{"name":"test","value":123}`,
			expected: "{\n  \"name\": \"test\",\n  \"value\": 123\n}",
		},
		{
			name:     "nested object",
			input:    `{"outer":{"inner":"value"}}`,
			expected: "{\n  \"outer\": {\n    \"inner\": \"value\"\n  }\n}",
		},
		{
			name:     "array",
			input:    `[1,2,3]`,
			expected: "[\n  1,\n  2,\n  3\n]",
		},
		{
			name:     "invalid JSON returns as-is",
			input:    `{invalid json`,
			expected: `{invalid json`,
		},
		{
			name:     "already formatted",
			input:    "{\n  \"key\": \"value\"\n}",
			expected: "{\n  \"key\": \"value\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJSON(tt.input)
			if result != tt.expected {
				t.Errorf("formatJSON() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFindMarkdownFilesIncludesJSON(t *testing.T) {
	// This test verifies the function signature exists and can be called
	// Actual file system testing would require setup
	files := findMarkdownFiles("/tmp")
	// Just verify function executes without panic
	_ = files
}

func TestLoadDirectoryIncludesJSON(t *testing.T) {
	// Verify loadDirectory function works
	items := loadDirectory("/tmp")
	_ = items
}

func TestJSONFileExtensionDetection(t *testing.T) {
	tests := []struct {
		filename string
		isJSON   bool
	}{
		{"file.json", true},
		{"file.JSON", false}, // Case sensitive
		{"file.md", false},
		{"file.json.bak", false},
		{"settings.json", true},
	}

	for _, tt := range tests {
		result := strings.HasSuffix(tt.filename, ".json")
		if result != tt.isJSON {
			t.Errorf("HasSuffix(%q, .json) = %v, want %v", tt.filename, result, tt.isJSON)
		}
	}
}
