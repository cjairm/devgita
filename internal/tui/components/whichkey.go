package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// WhichKeyEntry is a single keybinding entry in the which-key popup.
type WhichKeyEntry struct {
	Key, Desc string
}

// WhichKeyPopup renders the K2 which-key popup in a multi-column bordered box.
// cols < 1 → treated as 1. maxWidth < 6 → returns "".
func (p *Palette) WhichKeyPopup(title string, entries []WhichKeyEntry, cols, maxWidth int) string {
	if maxWidth < 6 {
		return ""
	}
	if cols < 1 {
		cols = 1
	}

	b := p.PaletteBorder.Render
	inner := maxWidth - 2
	sep := strings.Repeat("─", inner)
	colW := inner / cols

	titleText := ansi.Truncate(title, inner, "")
	lpad := max((inner-ansi.StringWidth(titleText))/2, 0)
	rpad := max(inner-lpad-ansi.StringWidth(titleText), 0)

	lines := []string{
		b("┌" + sep + "┐"),
		b(
			"│",
		) + strings.Repeat(
			" ",
			lpad,
		) + p.RepoHeader.Render(
			titleText,
		) + strings.Repeat(
			" ",
			rpad,
		) + b(
			"│",
		),
		b("├" + sep + "┤"),
	}

	for i := 0; i < len(entries); i += cols {
		var cells []string
		for c := 0; c < cols; c++ {
			if i+c >= len(entries) {
				cells = append(cells, strings.Repeat(" ", colW))
				continue
			}
			e := entries[i+c]
			keyS := p.HintKey.Render(e.Key)
			cell := " " + keyS + " " + e.Desc
			visW := 1 + ansi.StringWidth(e.Key) + 1 + ansi.StringWidth(e.Desc)
			pad := strings.Repeat(" ", max(0, colW-visW))
			cells = append(cells, cell+pad)
		}
		lines = append(lines, b("│")+strings.Join(cells, b("│"))+b("│"))
	}

	lines = append(lines, b("└"+sep+"┘"))
	return strings.Join(lines, "\n")
}
