package tuicomponents

import "strings"

// Tab is a single entry in a tab bar.
type Tab struct {
	Label string
	Extra string // optional suffix, e.g. a diff stat
}

// TabBar renders an underline-style tab bar.
// If activeIdx < 0 or >= len(tabs), no tab is marked active.
// Empty tabs slice returns "".
func (p *Palette) TabBar(tabs []Tab, activeIdx int) string {
	if len(tabs) == 0 {
		return ""
	}
	parts := make([]string, len(tabs))
	for i, t := range tabs {
		label := t.Label
		if t.Extra != "" {
			label += " " + t.Extra
		}
		if i == activeIdx {
			parts[i] = p.TabActive.Render(label)
		} else {
			parts[i] = p.TabInactive.Render(label)
		}
	}
	return strings.Join(parts, "  ")
}
