// TextInput is the single-line editor shared by the toolkit's text-editing
// components (FilterField, FuzzyPicker's query, and the worktree TUI's name
// prompt). It tracks the edit text plus a caret position, so every field
// built on it gets the same left/right/home/end movement and mid-string
// insert/delete — not just append-and-backspace — from one implementation.
// Centralizing the rune math here also keeps the "split a multi-byte rune on
// delete" class of bug impossible to reintroduce per-field.

package tuicomponents

import (
	"strings"
	"unicode/utf8"

	lip "charm.land/lipgloss/v2"
)

// caretStyle draws the block cursor as a reverse-video cell so it stays
// visible both at the end of the text and in the middle of it.
var caretStyle = lip.NewStyle().Reverse(true)

// plainStyle renders text with no color, for fields (filter, name prompt)
// that draw their edit text unstyled.
var plainStyle = lip.NewStyle()

// TextInput holds a single line of editable text and a caret. cursor is a
// rune offset in [0, runeCount(Value)] — Value's byte length is never a valid
// caret position because a multi-byte rune would split it.
type TextInput struct {
	Value  string
	cursor int
}

// SanitizePaste strips control characters (including newlines) from
// bracketed-paste content before it's inserted into a single-line text
// field. Bubble Tea delivers a paste as one tea.PasteMsg carrying the full
// clipboard content, which may include a trailing newline or other control
// bytes a terminal/clipboard tends to add — left in, they'd corrupt a field
// that renders as one line.
func SanitizePaste(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
}

// Cursor returns the caret's rune offset.
func (t TextInput) Cursor() int { return t.cursor }

// SetValue replaces the text and moves the caret to the end. Used when a field
// seeds an initial value rather than building it up keystroke by keystroke.
func (t *TextInput) SetValue(s string) {
	t.Value = s
	t.cursor = utf8.RuneCountInString(s)
}

// Reset clears the text and caret.
func (t *TextInput) Reset() {
	t.Value = ""
	t.cursor = 0
}

func (t TextInput) runeLen() int { return utf8.RuneCountInString(t.Value) }

// MoveLeft / MoveRight / Home / End reposition the caret without editing.
func (t *TextInput) MoveLeft() {
	if t.cursor > 0 {
		t.cursor--
	}
}

func (t *TextInput) MoveRight() {
	if t.cursor < t.runeLen() {
		t.cursor++
	}
}

func (t *TextInput) Home() { t.cursor = 0 }

func (t *TextInput) End() { t.cursor = t.runeLen() }

// Backspace deletes the rune before the caret and reports whether anything
// changed (nothing to delete at the start of the line).
func (t *TextInput) Backspace() bool {
	if t.cursor == 0 {
		return false
	}
	runes := []rune(t.Value)
	t.Value = string(runes[:t.cursor-1]) + string(runes[t.cursor:])
	t.cursor--
	return true
}

// Delete removes the rune at the caret (forward delete) and reports whether
// anything changed (nothing to delete at the end of the line).
func (t *TextInput) Delete() bool {
	runes := []rune(t.Value)
	if t.cursor >= len(runes) {
		return false
	}
	t.Value = string(runes[:t.cursor]) + string(runes[t.cursor+1:])
	return true
}

// Insert splices s in at the caret and advances the caret past it. s is
// inserted verbatim; callers that take raw clipboard/paste content should use
// InsertText so control bytes are stripped first.
func (t *TextInput) Insert(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(t.Value)
	t.Value = string(runes[:t.cursor]) + s + string(runes[t.cursor:])
	t.cursor += utf8.RuneCountInString(s)
	return true
}

// InsertText inserts sanitized paste content at the caret — the paste
// counterpart to a run of HandleKey calls.
func (t *TextInput) InsertText(text string) bool {
	return t.Insert(SanitizePaste(text))
}

// HandleKey applies one editing keypress. It reports handled=false for keys it
// doesn't own (esc, enter, up/down, ctrl+j/ctrl+k, …) so a caller can layer
// its own mode/navigation handling on top, and changed=true only when the
// text itself changed (caret movement is handled but not a change), so callers
// that re-filter on edits don't re-run on a bare arrow key.
func (t *TextInput) HandleKey(key string) (handled, changed bool) {
	switch key {
	case "left", "ctrl+b":
		t.MoveLeft()
		return true, false
	case "right", "ctrl+f":
		t.MoveRight()
		return true, false
	case "home", "ctrl+a":
		t.Home()
		return true, false
	case "end", "ctrl+e":
		t.End()
		return true, false
	case "backspace":
		return true, t.Backspace()
	case "delete":
		return true, t.Delete()
	default:
		// Count runes (not bytes) so a non-ASCII character typed directly (e.g.
		// "é", "日") still inserts instead of being dropped by a byte-length
		// check — every control-key string above has more than one rune.
		if utf8.RuneCountInString(key) == 1 && key >= " " {
			return true, t.Insert(key)
		}
	}
	return false, false
}

// RenderPlain renders the text with an unstyled block caret, for fields that
// draw their content without color.
func (t TextInput) RenderPlain() string {
	return renderCaret(t.Value, t.cursor, plainStyle)
}

// renderCaret draws value with a block caret at rune offset cursor: the rune
// under the caret is shown reverse-video, and a caret sitting at the end
// appends a reverse-video space so the line always shows where typing lands.
// textStyle styles the non-caret text.
func renderCaret(value string, cursor int, textStyle lip.Style) string {
	runes := []rune(value)
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(runes) {
		cursor = len(runes)
	}
	if cursor == len(runes) {
		return styleOrEmpty(textStyle, value) + caretStyle.Render(" ")
	}
	before := styleOrEmpty(textStyle, string(runes[:cursor]))
	under := caretStyle.Render(string(runes[cursor]))
	after := styleOrEmpty(textStyle, string(runes[cursor+1:]))
	return before + under + after
}

// styleOrEmpty avoids emitting a style's escape sequences around an empty
// segment (which would render as a stray reset), returning "" as-is.
func styleOrEmpty(style lip.Style, s string) string {
	if s == "" {
		return ""
	}
	return style.Render(s)
}
