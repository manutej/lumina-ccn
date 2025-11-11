package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ContextPanelMode represents different display modes for the right panel
type ContextPanelMode int

const (
	TOCMode ContextPanelMode = iota
	FileInfoMode
	DocumentStatsMode
	QuickActionsMode
)

// ContextPanel manages the right panel display with multiple modes
type ContextPanel struct {
	currentMode ContextPanelMode
	toc         *TableOfContents
	fileInfo    *FileInfo
	docStats    *DocumentStats
}

// FileInfo holds metadata about the current file
type FileInfo struct {
	Path         string
	Name         string
	Size         int64
	LineCount    int
	WordCount    int
	CharCount    int
	ReadingTime  int // minutes
	LastModified time.Time
}

// DocumentStats holds statistics about markdown content
type DocumentStats struct {
	HeadingsByLevel map[int]int // Level (1-6) -> count
	CodeBlockCount  int
	LinkCount       int
	ImageCount      int
	ListItemCount   int
	TableCount      int
}

// NewContextPanel creates a new context panel
func NewContextPanel() *ContextPanel {
	return &ContextPanel{
		currentMode: TOCMode,
		toc:         NewTableOfContents(),
		fileInfo:    &FileInfo{},
		docStats:    &DocumentStats{HeadingsByLevel: make(map[int]int)},
	}
}

// CycleMode switches to the next display mode
func (cp *ContextPanel) CycleMode() {
	cp.currentMode = (cp.currentMode + 1) % 4
}

// SetMode sets a specific mode
func (cp *ContextPanel) SetMode(mode ContextPanelMode) {
	cp.currentMode = mode
}

// GetModeName returns the current mode name
func (cp *ContextPanel) GetModeName() string {
	switch cp.currentMode {
	case TOCMode:
		return "Table of Contents"
	case FileInfoMode:
		return "File Info"
	case DocumentStatsMode:
		return "Document Stats"
	case QuickActionsMode:
		return "Quick Actions"
	default:
		return "Unknown"
	}
}

// UpdateContent analyzes content and updates all modes
func (cp *ContextPanel) UpdateContent(filePath string, content string) {
	// Update TOC
	cp.toc.ParseMarkdown(content)
	cp.toc.SelectFirst()

	// Update File Info
	cp.updateFileInfo(filePath, content)

	// Update Document Stats
	cp.updateDocumentStats(content)
}

// updateFileInfo calculates file metadata
func (cp *ContextPanel) updateFileInfo(filePath string, content string) {
	info, err := os.Stat(filePath)
	if err != nil {
		return
	}

	lines := strings.Split(content, "\n")
	words := len(strings.Fields(content))
	readingTime := words / 200 // Assuming 200 words per minute

	cp.fileInfo = &FileInfo{
		Path:         filePath,
		Name:         filepath.Base(filePath),
		Size:         info.Size(),
		LineCount:    len(lines),
		WordCount:    words,
		CharCount:    len(content),
		ReadingTime:  readingTime,
		LastModified: info.ModTime(),
	}
}

// updateDocumentStats analyzes markdown structure
func (cp *ContextPanel) updateDocumentStats(content string) {
	stats := &DocumentStats{
		HeadingsByLevel: make(map[int]int),
	}

	lines := strings.Split(content, "\n")

	// Regexes
	headingRegex := regexp.MustCompile(`^(#{1,6})\s+.+$`)
	codeBlockRegex := regexp.MustCompile("^```")
	linkRegex := regexp.MustCompile(`\[.*?\]\(.*?\)`)
	imageRegex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	listRegex := regexp.MustCompile(`^\s*[-*+]\s+`)
	tableRegex := regexp.MustCompile(`^\|.*\|$`)

	inCodeBlock := false

	for _, line := range lines {
		// Code blocks
		if codeBlockRegex.MatchString(line) {
			if !inCodeBlock {
				stats.CodeBlockCount++
			}
			inCodeBlock = !inCodeBlock
		}

		// Headings
		if matches := headingRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			stats.HeadingsByLevel[level]++
		}

		// Links (excluding images)
		stats.LinkCount += len(linkRegex.FindAllString(line, -1))

		// Images
		stats.ImageCount += len(imageRegex.FindAllString(line, -1))

		// List items
		if listRegex.MatchString(line) {
			stats.ListItemCount++
		}

		// Tables (simple detection)
		if tableRegex.MatchString(line) {
			stats.TableCount++
		}
	}

	// Adjust link count (subtract images)
	stats.LinkCount -= stats.ImageCount

	cp.docStats = stats
}

// Render renders the context panel based on current mode
func (cp *ContextPanel) Render(width, height int, colorManager *ColorManager) string {
	switch cp.currentMode {
	case TOCMode:
		return cp.renderTOC(width, height, colorManager)
	case FileInfoMode:
		return cp.renderFileInfo(width, height, colorManager)
	case DocumentStatsMode:
		return cp.renderDocumentStats(width, height, colorManager)
	case QuickActionsMode:
		return cp.renderQuickActions(width, height, colorManager)
	default:
		return "Unknown Mode"
	}
}

// renderTOC renders the table of contents
func (cp *ContextPanel) renderTOC(width, height int, colorManager *ColorManager) string {
	if !cp.toc.HasEntries() {
		return "📋 Table of Contents\n\n(No headings found)\n\nPress 'm' to switch modes"
	}

	content := cp.toc.View(width, height-3) // Leave room for mode switcher
	content += "\n\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Press 'm' to switch modes")

	return content
}

// renderFileInfo renders file metadata
func (cp *ContextPanel) renderFileInfo(width, height int, colorManager *ColorManager) string {
	if cp.fileInfo.Path == "" {
		return "📄 File Info\n\n(No file selected)\n\nPress 'm' to switch modes"
	}

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorManager.GetColor("help-command")))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	var lines []string
	lines = append(lines, "📄 File Info\n")
	lines = append(lines, labelStyle.Render("Name:"))
	lines = append(lines, "  "+valueStyle.Render(cp.fileInfo.Name))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Size:"))
	lines = append(lines, "  "+valueStyle.Render(formatBytes(cp.fileInfo.Size)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Lines:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.fileInfo.LineCount)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Words:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.fileInfo.WordCount)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Characters:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.fileInfo.CharCount)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Reading Time:"))
	if cp.fileInfo.ReadingTime == 0 {
		lines = append(lines, "  "+valueStyle.Render("< 1 min"))
	} else {
		lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d min", cp.fileInfo.ReadingTime)))
	}
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Modified:"))
	lines = append(lines, "  "+valueStyle.Render(cp.fileInfo.LastModified.Format("2006-01-02 15:04")))

	lines = append(lines, "")
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Press 'm' to switch modes"))

	return strings.Join(lines, "\n")
}

// renderDocumentStats renders document statistics
func (cp *ContextPanel) renderDocumentStats(width, height int, colorManager *ColorManager) string {
	if cp.fileInfo.Path == "" {
		return "📈 Document Stats\n\n(No file selected)\n\nPress 'm' to switch modes"
	}

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorManager.GetColor("help-command")))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	var lines []string
	lines = append(lines, "📈 Document Stats\n")

	// Headings
	lines = append(lines, labelStyle.Render("Headings:"))
	totalHeadings := 0
	for level := 1; level <= 6; level++ {
		count := cp.docStats.HeadingsByLevel[level]
		if count > 0 {
			totalHeadings += count
			lines = append(lines, fmt.Sprintf("  H%d: %s", level, valueStyle.Render(fmt.Sprintf("%d", count))))
		}
	}
	if totalHeadings == 0 {
		lines = append(lines, "  "+valueStyle.Render("None"))
	}
	lines = append(lines, "")

	// Other stats
	lines = append(lines, labelStyle.Render("Code Blocks:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.docStats.CodeBlockCount)))
	lines = append(lines, "")

	lines = append(lines, labelStyle.Render("Links:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.docStats.LinkCount)))
	lines = append(lines, "")

	lines = append(lines, labelStyle.Render("Images:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.docStats.ImageCount)))
	lines = append(lines, "")

	lines = append(lines, labelStyle.Render("List Items:"))
	lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.docStats.ListItemCount)))
	lines = append(lines, "")

	if cp.docStats.TableCount > 0 {
		lines = append(lines, labelStyle.Render("Table Rows:"))
		lines = append(lines, "  "+valueStyle.Render(fmt.Sprintf("%d", cp.docStats.TableCount)))
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Press 'm' to switch modes"))

	return strings.Join(lines, "\n")
}

// renderQuickActions renders quick actions menu
func (cp *ContextPanel) renderQuickActions(width, height int, colorManager *ColorManager) string {
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorManager.GetColor("help-command")))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	var lines []string
	lines = append(lines, "⚡ Quick Actions\n")
	lines = append(lines, keyStyle.Render("c")+" "+descStyle.Render("Copy file path"))
	lines = append(lines, keyStyle.Render("n")+" "+descStyle.Render("Copy filename"))
	lines = append(lines, keyStyle.Render("y")+" "+descStyle.Render("Copy content"))
	lines = append(lines, keyStyle.Render("r")+" "+descStyle.Render("Reload file"))
	lines = append(lines, "")
	lines = append(lines, descStyle.Render("(Coming soon)"))
	lines = append(lines, "  • Open in $EDITOR")
	lines = append(lines, "  • Show git status")
	lines = append(lines, "  • Reveal in tree")
	lines = append(lines, "")
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Press 'm' to switch modes"))

	return strings.Join(lines, "\n")
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetTOC returns the TOC component
func (cp *ContextPanel) GetTOC() *TableOfContents {
	return cp.toc
}

// Clear resets all panel data
func (cp *ContextPanel) Clear() {
	cp.toc = NewTableOfContents()
	cp.fileInfo = &FileInfo{}
	cp.docStats = &DocumentStats{HeadingsByLevel: make(map[int]int)}
}
