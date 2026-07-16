package tuicomponents

import (
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
)

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
			f.Text = TrimLastRune(f.Text)
			changed = true
		}
	default:
		// Counting runes (not bytes) so a non-ASCII character typed directly
		// (e.g. "é") still falls into this case instead of being dropped by a
		// byte-length check — every control-key string this switch matches
		// ("esc", "backspace", ...) has more than one rune either way.
		if utf8.RuneCountInString(key) == 1 && key >= " " {
			f.Text += key
			changed = true
		}
	}
	return changed
}

// InsertText inserts pasted text into the filter in one shot and reports
// whether it changed anything, the paste counterpart to HandleKey: a
// tea.PasteMsg carries the whole clipboard content as one string, and
// HandleKey's default case only accepts single-rune keys, so a multi-rune
// paste would otherwise be silently dropped rune-by-rune.
func (f *FilterField) InsertText(text string) (changed bool) {
	text = SanitizePaste(text)
	if text == "" {
		return false
	}
	f.Text += text
	return true
}

// FilterHint renders the hint-bar line shown while a filter field is active.
func (p *Palette) FilterHint(f FilterField, width int) string {
	hint := "filter: " + f.Text + "█  · esc: clear · enter: keep"
	return p.HintDesc.Render(ansi.Truncate(hint, width, ""))
}
