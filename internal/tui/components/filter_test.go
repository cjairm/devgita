package tuicomponents

import (
	"strings"
	"testing"
)

func TestFilterFieldHandleKey(t *testing.T) {
	t.Run("printable characters append and report change", func(t *testing.T) {
		f := FilterField{Active: true}
		if !f.HandleKey("a") || !f.HandleKey("b") {
			t.Error("typing should report a change")
		}
		if f.Text != "ab" {
			t.Errorf("expected text %q, got %q", "ab", f.Text)
		}
	})

	t.Run("backspace deletes", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if !f.HandleKey("backspace") {
			t.Error("backspace should report a change")
		}
		if f.Text != "a" {
			t.Errorf("expected text %q, got %q", "a", f.Text)
		}
	})

	t.Run("backspace on empty text is a no-op", func(t *testing.T) {
		f := FilterField{Active: true}
		if f.HandleKey("backspace") {
			t.Error("backspace on empty text should not report a change")
		}
	})

	t.Run("esc clears and deactivates", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if !f.HandleKey("esc") {
			t.Error("esc with text should report a change")
		}
		if f.Active || f.Text != "" {
			t.Errorf("esc should clear and deactivate, got %+v", f)
		}
	})

	t.Run("enter keeps text and deactivates", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if f.HandleKey("enter") {
			t.Error("enter should not report a change")
		}
		if f.Active || f.Text != "ab" {
			t.Errorf("enter should keep text and deactivate, got %+v", f)
		}
	})

	t.Run("non-printable keys are ignored", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if f.HandleKey("ctrl+c") {
			t.Error("multi-char keys should not change the filter")
		}
		if f.Text != "ab" {
			t.Errorf("expected text unchanged, got %q", f.Text)
		}
	})
}

func TestFilterFieldInsertText(t *testing.T) {
	t.Run("inserts pasted text in one shot and reports a change", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if !f.InsertText("cd/ef") {
			t.Error("expected InsertText to report a change")
		}
		if f.Text != "abcd/ef" {
			t.Errorf("expected text %q, got %q", "abcd/ef", f.Text)
		}
	})

	t.Run("strips control characters from the pasted content", func(t *testing.T) {
		f := FilterField{Active: true}
		f.InsertText("feat\nname\r")
		if f.Text != "featname" {
			t.Errorf("expected control chars stripped, got %q", f.Text)
		}
	})

	t.Run("all-control paste reports no change", func(t *testing.T) {
		f := FilterField{Active: true, Text: "ab"}
		if f.InsertText("\n\r") {
			t.Error("expected no change when pasted content is entirely control chars")
		}
		if f.Text != "ab" {
			t.Errorf("expected text unchanged, got %q", f.Text)
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
