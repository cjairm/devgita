package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// ansiReset closes any open SGR styling. Overlay stamps this at both edges of
// the popup so a background color never bleeds into the popup box, and a
// popup color never bleeds back into the background that resumes after it.
const ansiReset = "\x1b[0m"

// Overlay composites popup on top of background, centered within a
// width×height screen. Background lines directly under the popup are cut at
// the popup's left edge and resumed after its right edge; everything else in
// the background is left untouched. This is the building block that lets a
// modal (e.g. the help popup) float over a live screen instead of replacing
// it outright.
func Overlay(background, popup string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	popupLines, popupWidth, popupHeight := clipPopup(popup, width, height)
	bgLines := canvasLines(background, width, height)

	offsetX := max((width-popupWidth)/2, 0)
	offsetY := max((height-popupHeight)/2, 0)

	out := make([]string, height)
	for i, bg := range bgLines {
		if i < offsetY || i >= offsetY+popupHeight {
			out[i] = bg
			continue
		}
		left := ansi.Truncate(bg, offsetX, "")
		right := ansi.Cut(bg, offsetX+popupWidth, width)
		out[i] = left + ansiReset + popupLines[i-offsetY] + ansiReset + right
	}
	return strings.Join(out, "\n")
}

// clipPopup splits popup into lines, clips it to fit within width×height (a
// popup larger than the screen would otherwise blow past the background
// canvas and misalign the seams), and pads every line to a uniform width so
// the seam math in Overlay always cuts at exact column boundaries.
func clipPopup(popup string, width, height int) (lines []string, w, h int) {
	raw := strings.Split(popup, "\n")
	h = min(len(raw), height)
	raw = raw[:h]

	for _, l := range raw {
		w = max(w, ansi.StringWidth(l))
	}
	w = min(w, width)

	lines = make([]string, h)
	for i, l := range raw {
		lines[i] = padToWidth(l, w)
	}
	return lines, w, h
}

// canvasLines renders background as exactly height lines of exactly width
// display columns each, padding short/missing lines and truncating long ones,
// so the rest of Overlay can index into it without bounds checks.
func canvasLines(background string, width, height int) []string {
	raw := strings.Split(background, "\n")
	lines := make([]string, height)
	for i := range height {
		if i < len(raw) {
			lines[i] = padToWidth(raw[i], width)
		} else {
			lines[i] = strings.Repeat(" ", width)
		}
	}
	return lines
}

// padToWidth truncates or space-pads line to exactly width display columns,
// using ansi.StringWidth so escape codes aren't counted as visible width.
func padToWidth(line string, width int) string {
	if width <= 0 {
		return ""
	}
	t := ansi.Truncate(line, width, "")
	if pad := width - ansi.StringWidth(t); pad > 0 {
		t += strings.Repeat(" ", pad)
	}
	return t
}
