# Test Golden Files

This directory contains golden files for snapshot testing during the RED and GREEN phases of TDD.

## Golden File Organization

### Search Results
- `search_result_simple.golden` - Simple search result rendering
- `search_result_long_path.golden` - Search result with long file path
- `search_result_multiline.golden` - Multiline snippet rendering
- `search_result_narrow.golden` - Narrow terminal width rendering

### UI Rendering
- `ui_render_initial.golden` - Initial UI state
- `ui_render_search.golden` - UI in search mode
- `ui_render_filtered.golden` - UI with active filter

### Syntax Highlighting
- `syntax_go.golden` - Go code block syntax highlighting
- `syntax_python.golden` - Python code block syntax highlighting
- `syntax_json.golden` - JSON code block syntax highlighting
- `syntax_plain.golden` - Plain code block (no language)
- `inline_code.golden` - Inline code rendering

### Markdown Elements
- `heading_h1.golden` - H1 heading rendering
- `heading_h2.golden` - H2 heading rendering
- `bold.golden` - Bold text rendering
- `italic.golden` - Italic text rendering
- `link.golden` - Link rendering
- `list_unordered.golden` - Unordered list rendering
- `list_ordered.golden` - Ordered list rendering
- `blockquote.golden` - Blockquote rendering
- `hr.golden` - Horizontal rule rendering
- `table.golden` - Table rendering

## Usage in Tests

Golden files are used in tests like this:

```go
func TestRender(t *testing.T) {
    renderer := NewGlamourRenderer()
    output := renderer.Render("# Heading")

    golden := readGoldenFile("heading_h1.golden")
    if output != golden {
        t.Errorf("output differs from golden file")
    }
}
```

## RED Phase

During the RED phase, these files are created as **empty placeholders**. Tests will fail because:
1. The rendering functions don't exist yet
2. The golden files are empty

## GREEN Phase

During the GREEN phase:
1. Implement the rendering functions
2. Generate actual output
3. Fill in the golden files with correct expected output
4. Tests should pass when output matches golden files

## Updating Golden Files

When expected output changes (e.g., improving formatting), update golden files by:
1. Run tests to see current output
2. Review output for correctness
3. Update golden files with new expected output
4. Re-run tests to verify they pass
