package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// BorderedPane renders a rounded-corner box with the title embedded in the top
// border line (e.g. "╭─ inventory ─────────╮"), matching the T1 wireframe page.
// lines are body rows (may already carry ANSI styling); each is padded or
// truncated to fit the interior width so every returned line has exactly
// `width` display columns.
func (p *Palette) BorderedPane(title string, width int, lines []string) string {
	if width < 6 {
		width = 6
	}
	inner := width - 2
	border := p.PaletteBorder.Render

	// Top border: ╭─ title ─────╮
	// Total width = 1 (╭) + 1 (─) + 1 (space) + len(title) + 1 (space) + dashes + 1 (╮)
	// The space after "╭─" is required so the rendered width matches `width`
	// exactly — omitting it leaves the top border one column short.
	// title must also be truncated to fit, otherwise a title wider than
	// width-5 pushes the top border past `width` (dashCount alone can't fix it).
	maxTitleWidth := width - 5
	if maxTitleWidth < 0 {
		maxTitleWidth = 0
	}
	title = ansi.Truncate(title, maxTitleWidth, "")
	dashCount := width - 5 - ansi.StringWidth(title)
	if dashCount < 0 {
		dashCount = 0
	}
	top := border(
		"╭─ ",
	) + p.RepoHeader.Render(
		title,
	) + border(
		" "+strings.Repeat("─", dashCount)+"╮",
	)

	var sb strings.Builder
	sb.WriteString(top)
	for _, line := range lines {
		trimmed := ansi.Truncate(line, inner, "")
		pad := inner - ansi.StringWidth(trimmed)
		if pad < 0 {
			pad = 0
		}
		sb.WriteString("\n")
		sb.WriteString(border("│") + trimmed + strings.Repeat(" ", pad) + border("│"))
	}
	sb.WriteString("\n")
	sb.WriteString(border("╰" + strings.Repeat("─", inner) + "╯"))
	return sb.String()
}
