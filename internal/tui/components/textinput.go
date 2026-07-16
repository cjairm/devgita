// Shared helpers for the toolkit's single-line text-editing components
// (FilterField, FuzzyPicker's query, and the worktree TUI's name prompt):
// sanitizing bracketed-paste content before insertion, and deleting the
// last rune (rather than the last byte) on backspace.

package tuicomponents

import (
	"strings"
	"unicode/utf8"
)

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

// TrimLastRune removes the last rune of s. Plain byte-slice truncation
// (s[:len(s)-1]) splits a multi-byte UTF-8 rune (e.g. an accented letter or
// emoji pasted in) and leaves an invalid trailing byte sequence, so
// backspace must decode the last rune's width first.
func TrimLastRune(s string) string {
	if s == "" {
		return s
	}
	_, size := utf8.DecodeLastRuneInString(s)
	return s[:len(s)-size]
}
