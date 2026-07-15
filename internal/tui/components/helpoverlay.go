package tuicomponents

import "strings"

// HelpOverlay renders a centered which-key popup with a "press any key to
// close" hint, positioned for a width×height screen. It is the standard help
// modal for all devgita TUIs; callers show it on `?` and dismiss it on any
// key press.
func (p *Palette) HelpOverlay(title string, entries []WhichKeyEntry, width, height int) string {
	maxW := min(width-2, 64)
	popup := p.WhichKeyPopup(title, entries, 1, maxW) + "\n" +
		p.HintDesc.Render("press any key to close")

	lines := strings.Split(popup, "\n")
	topPad := max((height-len(lines))/2, 0)
	indent := strings.Repeat(" ", max((width-maxW)/2, 0))

	out := make([]string, 0, topPad+len(lines))
	for range topPad {
		out = append(out, "")
	}
	for _, l := range lines {
		out = append(out, indent+l)
	}
	return strings.Join(out, "\n")
}
