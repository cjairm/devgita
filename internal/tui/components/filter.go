package tuicomponents

import "github.com/charmbracelet/x/ansi"

// FilterField is the text-filter state shared by list TUIs: `/` activates it,
// then keys edit the text until esc (clear) or enter (keep) deactivates it.
type FilterField struct {
	Active bool
	Text   string
}

// HandleKey processes one key press while the filter is active and reports
// whether the filter text changed (callers rebuild their rows on change).
// esc clears the text and deactivates, enter keeps the text and deactivates,
// backspace deletes, and any printable character appends.
func (f *FilterField) HandleKey(key string) (changed bool) {
	switch key {
	case "esc":
		changed = f.Text != ""
		f.Text = ""
		f.Active = false
	case "enter":
		f.Active = false
	case "backspace":
		if len(f.Text) > 0 {
			f.Text = f.Text[:len(f.Text)-1]
			changed = true
		}
	default:
		if len(key) == 1 && key >= " " {
			f.Text += key
			changed = true
		}
	}
	return changed
}

// FilterHint renders the hint-bar line shown while a filter field is active.
func (p *Palette) FilterHint(f FilterField, width int) string {
	hint := "filter: " + f.Text + "█  · esc: clear · enter: keep"
	return p.HintDesc.Render(ansi.Truncate(hint, width, ""))
}
