package tuicomponents

import "strings"

// helpPopupContent builds the which-key popup body plus its close hint, sized
// to fit within maxWidth. Shared by HelpOverlay (self-centered, for callers
// that render it as the entire screen) and HelpPopup (raw, for callers that
// composite it over a live background via Overlay) so the popup content can't
// drift between the two.
func (p *Palette) helpPopupContent(title string, entries []WhichKeyEntry, maxWidth int) string {
	return p.WhichKeyPopup(title, entries, 1, maxWidth) + "\n" +
		p.HintDesc.Render("press any key to close")
}

// HelpPopup renders the which-key popup and close hint sized to fit within
// width, without centering it on a screen. Callers composite the result over
// a live background via Overlay instead of replacing the screen outright.
func (p *Palette) HelpPopup(title string, entries []WhichKeyEntry, width int) string {
	return p.helpPopupContent(title, entries, min(width-2, 64))
}

// HelpOverlay renders a centered which-key popup with a "press any key to
// close" hint, positioned for a width×height screen. It is the standard help
// modal for TUIs that render the popup as the entire screen content; callers
// show it on `?` and dismiss it on any key press.
func (p *Palette) HelpOverlay(title string, entries []WhichKeyEntry, width, height int) string {
	maxW := min(width-2, 64)
	popup := p.helpPopupContent(title, entries, maxW)

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
