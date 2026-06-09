package tuicomponents

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// ToastKind controls the border/title color of a notification toast.
type ToastKind int

const (
	ToastNeedsReview ToastKind = iota
	ToastInfo
	ToastError
)

// Toast is the data model for a notification toast.
type Toast struct {
	Kind   ToastKind
	Title  string
	Body   string // omitted if empty
	Action string // e.g. "⏎ to attach"; omitted if empty
}

// Notification renders a bordered toast box for top-right placement.
// maxWidth < 6 → returns "".
func (p *Palette) Notification(t Toast, maxWidth int) string {
	if maxWidth < 6 {
		return ""
	}
	inner := maxWidth - 2
	b := p.ToastBorder.Render
	sep := strings.Repeat("─", inner)

	row := func(content string, visLen int) string {
		pad := strings.Repeat(" ", max(0, inner-visLen))
		return b("│") + content + pad + b("│")
	}

	lines := []string{b("┌" + sep + "┐")}

	title := ansi.Truncate(t.Title, inner, "")
	lines = append(lines, row(p.ToastTitle.Render(title), ansi.StringWidth(title)))

	if t.Body != "" {
		body := ansi.Truncate(t.Body, inner, "")
		lines = append(lines, row(p.ToastBody.Render(body), ansi.StringWidth(body)))
	}
	if t.Action != "" {
		action := ansi.Truncate(t.Action, inner, "")
		lines = append(lines, row(p.ToastAction.Render(action), ansi.StringWidth(action)))
	}

	lines = append(lines, b("└"+sep+"┘"))
	return strings.Join(lines, "\n")
}
