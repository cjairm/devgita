package tuicomponents

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// StatusBarMode controls the mode badge rendered on the left of the status bar.
type StatusBarMode int

const (
	ModeNormal  StatusBarMode = iota
	ModeCommand               // command mode badge
	ModeInsert                // insert mode badge
)

// StatusBarModel is the data model for StatusBar.
type StatusBarModel struct {
	Mode       StatusBarMode
	Breadcrumb string // e.g. "repo-a > branch-1"
	State      SessionState
	StateLabel string // e.g. "running" or "2 running · 1 done"
	Added      int
	Removed    int
	Index      int // 1-based position; ignored when Total == 0
	Total      int
}

// StatusBar renders the layout-B style bottom bar across the full width.
// Layout: [mode badge + breadcrumb] ... [dot + state label] ... [diff stat + position]
func (p *Palette) StatusBar(m StatusBarModel, width int) string {
	var badge string
	switch m.Mode {
	case ModeCommand:
		badge = p.ModeCommand.Render("COMMAND")
	case ModeInsert:
		badge = p.ModeInsert.Render("INSERT")
	default:
		badge = p.ModeNormal.Render("NORMAL")
	}

	left := badge + " " + m.Breadcrumb
	mid := p.StatusDot(m.State) + " " + m.StateLabel
	right := p.DiffStat(m.Added, m.Removed)
	if m.Total > 0 {
		right += " " + fmt.Sprintf("%d/%d", m.Index, m.Total)
	}

	leftW := ansi.StringWidth(left)
	midW := ansi.StringWidth(mid)
	rightW := ansi.StringWidth(right)
	total := leftW + midW + rightW

	if total >= width {
		return ansi.Truncate(left, width, "")
	}

	gap := width - total
	space1 := gap / 2
	space2 := gap - space1
	return left + strings.Repeat(" ", space1) + mid + strings.Repeat(" ", space2) + right
}
