package tuicomponents

import (
	"strings"
	"testing"
)

// seededFilter builds an active filter whose text is text with the caret at
// the end, matching how a user would have typed it.
func seededFilter(text string) FilterField {
	f := FilterField{Active: true}
	f.input.SetValue(text)
	return f
}

func TestFilterFieldHandleKey(t *testing.T) {
	t.Run("printable characters append and report change", func(t *testing.T) {
		f := FilterField{Active: true}
		if !f.HandleKey("a") || !f.HandleKey("b") {
			t.Error("typing should report a change")
		}
		if f.Value() != "ab" {
			t.Errorf("expected text %q, got %q", "ab", f.Value())
		}
	})

	t.Run("backspace deletes", func(t *testing.T) {
		f := seededFilter("ab")
		if !f.HandleKey("backspace") {
			t.Error("backspace should report a change")
		}
		if f.Value() != "a" {
			t.Errorf("expected text %q, got %q", "a", f.Value())
		}
	})

	t.Run("backspace on empty text is a no-op", func(t *testing.T) {
		f := FilterField{Active: true}
		if f.HandleKey("backspace") {
			t.Error("backspace on empty text should not report a change")
		}
	})

	t.Run("caret movement edits mid-string without reporting a change", func(t *testing.T) {
		f := seededFilter("ab")
		if f.HandleKey("left") {
			t.Error("caret movement should not report a text change")
		}
		if !f.HandleKey("X") {
			t.Error("typing should report a change")
		}
		if f.Value() != "aXb" {
			t.Errorf("expected caret insert to give %q, got %q", "aXb", f.Value())
		}
	})

	t.Run("esc clears and deactivates", func(t *testing.T) {
		f := seededFilter("ab")
		if !f.HandleKey("esc") {
			t.Error("esc with text should report a change")
		}
		if f.Active || f.Value() != "" {
			t.Errorf("esc should clear and deactivate, got %+v", f)
		}
	})

	t.Run("enter keeps text and deactivates", func(t *testing.T) {
		f := seededFilter("ab")
		if f.HandleKey("enter") {
			t.Error("enter should not report a change")
		}
		if f.Active || f.Value() != "ab" {
			t.Errorf("enter should keep text and deactivate, got %+v", f)
		}
	})

	t.Run("non-printable keys are ignored", func(t *testing.T) {
		f := seededFilter("ab")
		if f.HandleKey("ctrl+c") {
			t.Error("multi-char keys should not change the filter")
		}
		if f.Value() != "ab" {
			t.Errorf("expected text unchanged, got %q", f.Value())
		}
	})
}

func TestFilterFieldInsertText(t *testing.T) {
	t.Run("inserts pasted text in one shot and reports a change", func(t *testing.T) {
		f := seededFilter("ab")
		if !f.InsertText("cd/ef") {
			t.Error("expected InsertText to report a change")
		}
		if f.Value() != "abcd/ef" {
			t.Errorf("expected text %q, got %q", "abcd/ef", f.Value())
		}
	})

	t.Run("strips control characters from the pasted content", func(t *testing.T) {
		f := FilterField{Active: true}
		f.InsertText("feat\nname\r")
		if f.Value() != "featname" {
			t.Errorf("expected control chars stripped, got %q", f.Value())
		}
	})

	t.Run("all-control paste reports no change", func(t *testing.T) {
		f := seededFilter("ab")
		if f.InsertText("\n\r") {
			t.Error("expected no change when pasted content is entirely control chars")
		}
		if f.Value() != "ab" {
			t.Errorf("expected text unchanged, got %q", f.Value())
		}
	})
}

func TestHelpOverlayContainsEntriesAndCloseHint(t *testing.T) {
	p := NewPalette()
	out := p.HelpOverlay("Keys", []WhichKeyEntry{{Key: "q", Desc: "quit"}}, 80, 24)
	if !strings.Contains(out, "quit") {
		t.Error("overlay should contain the entry description")
	}
	if !strings.Contains(out, "press any key to close") {
		t.Error("overlay should contain the close hint")
	}
}
