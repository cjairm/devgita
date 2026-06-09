package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// KeyHint is a single key+description pair for the hint bar.
type KeyHint struct {
	Key, Desc string
}

// HintBar renders the K1 persistent hint bar as "key desc · key desc · ..." truncated to width.
// Empty hints slice returns "". Width <= 0 returns "".
// String width is measured with ansi.StringWidth (display width, not byte length).
func (p *Palette) HintBar(hints []KeyHint, width int) string {
	if len(hints) == 0 || width <= 0 {
		return ""
	}
	sep := p.HintSep.Render(" · ")
	parts := make([]string, len(hints))
	for i, h := range hints {
		if h.Key != "" {
			parts[i] = p.HintKey.Render(h.Key) + " " + p.HintDesc.Render(h.Desc)
		} else {
			parts[i] = p.HintDesc.Render(h.Desc)
		}
	}
	return ansi.Truncate(strings.Join(parts, sep), width, "")
}
