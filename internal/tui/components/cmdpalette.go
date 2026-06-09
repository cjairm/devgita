package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// PaletteItem is a single entry in the command palette.
type PaletteItem struct {
	Command, Hint string
}

// CommandPalette renders the K3 command palette overlay.
// Layout: ": query█" input line followed by items with right-aligned hints.
// selectedIdx out of range → no item highlighted. width < 6 → returns "".
func (p *Palette) CommandPalette(query string, items []PaletteItem, selectedIdx, width int) string {
	if width < 6 {
		return ""
	}

	inputLine := ansi.Truncate(": "+p.PaletteInput.Render(query)+"█", width, "")

	lines := []string{inputLine}
	for i, item := range items {
		hintW := ansi.StringWidth(item.Hint)
		cmdMax := width - hintW - 1
		if cmdMax < 1 {
			cmdMax = 1
		}
		cmd := ansi.Truncate(item.Command, cmdMax, "")
		pad := strings.Repeat(" ", max(0, width-ansi.StringWidth(cmd)-hintW))
		hint := p.PaletteHint.Render(item.Hint)

		if i == selectedIdx {
			// Render the whole row with selection highlight (plain text for bg color)
			plain := ansi.Truncate(item.Command+pad+item.Hint, width, "")
			lines = append(lines, p.PaletteSelected.Render(plain))
		} else {
			lines = append(lines, cmd+pad+hint)
		}
	}
	return strings.Join(lines, "\n")
}
