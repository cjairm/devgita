package tuicomponents

import (
	"github.com/charmbracelet/x/ansi"
)

// FilterField is the text-filter state shared by list TUIs: `/` activates it,
// then keys edit the text until esc (clear) or enter (keep) deactivates it.
// The edit text and caret live in a shared TextInput, so the filter gets the
// same left/right/home/end and mid-string editing as every other text field.
type FilterField struct {
	Active bool
	input  TextInput
}

// Value returns the current filter text.
func (f *FilterField) Value() string { return f.input.Value }

// HandleKey processes one key press while the filter is active and reports
// whether the filter text changed (callers rebuild their rows on change).
// esc clears the text and deactivates, enter keeps the text and deactivates,
// and every other key is an edit delegated to the shared TextInput (backspace,
// caret movement, and printable insertion).
func (f *FilterField) HandleKey(key string) (changed bool) {
	switch key {
	case "esc":
		changed = f.input.Value != ""
		f.input.Reset()
		f.Active = false
		return changed
	case "enter":
		f.Active = false
		return false
	default:
		_, changed = f.input.HandleKey(key)
		return changed
	}
}

// InsertText inserts pasted text into the filter in one shot and reports
// whether it changed anything, the paste counterpart to HandleKey: a
// tea.PasteMsg carries the whole clipboard content as one string, and
// HandleKey's per-key editing only accepts single-rune keys, so a multi-rune
// paste would otherwise be silently dropped rune-by-rune.
func (f *FilterField) InsertText(text string) (changed bool) {
	return f.input.InsertText(text)
}

// FilterHint renders the hint-bar line shown while a filter field is active.
func (p *Palette) FilterHint(f FilterField, width int) string {
	hint := "filter: " + f.input.RenderPlain() + "  · esc: clear · enter: keep"
	return p.HintDesc.Render(ansi.Truncate(hint, width, ""))
}
