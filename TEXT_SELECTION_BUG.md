# Text Selection Bug Report

## Current Implementation Status

### Existing Code
- Selection state tracking in `clipboard.go`
- Keyboard selection modes (Character, Line, Block)
- Copy functionality implemented
- No visual text highlighting

### Identified Issues
1. **No Visual Highlighting**
   - Selection state exists, but no rendering mechanism
   - Users cannot see which text is currently selected
   - Breaks expected text selection UX

2. **Partial Selection Extraction**
   - `extractSelection()` function handles multi-line selection
   - Potential edge cases with line ending/start of selection
   - Needs comprehensive testing

3. **Limited Selection Modes**
   - Currently supports Character, Line, Block modes
   - No mouse-based selection
   - No shift+arrow key selection

### Recommended Fixes
1. Implement visual text highlighting
   - Use Lip Gloss or Bubbles for selection styling
   - Add background color or inverse text for selected regions
2. Enhance keyboard selection
   - Support shift+arrow key selection
   - Improve word/line selection accuracy
3. Add mouse selection support
4. Create comprehensive selection unit tests

### Code Changes Needed
- Modify `model.go` to add selection rendering
- Update `clipboard.go` with more robust selection logic
- Create selection highlight utility in a new file

## Test Scenarios
- Single line selection
- Multi-line selection
- Block selection
- Edge of document selections
- Unicode and special character selections
- Performance with large documents

## Estimated Effort
- High Priority
- 2-3 days of implementation
- Comprehensive testing required

## Potential Performance Considerations
- Optimize selection extraction for large files
- Minimize memory overhead during selection