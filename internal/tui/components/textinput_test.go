package tuicomponents

import (
	"strings"
	"testing"
)

func TestSanitizePaste(t *testing.T) {
	t.Run("strips embedded newlines and carriage returns", func(t *testing.T) {
		got := SanitizePaste("feat/new-thing\n")
		if got != "feat/new-thing" {
			t.Errorf("expected trailing newline stripped, got %q", got)
		}
		got = SanitizePaste("feat\r\nbranch")
		if got != "featbranch" {
			t.Errorf("expected CRLF stripped, got %q", got)
		}
	})

	t.Run("strips other control characters", func(t *testing.T) {
		got := SanitizePaste("a\tb\x1bc\x7fd")
		if got != "abcd" {
			t.Errorf("expected control chars stripped, got %q", got)
		}
	})

	t.Run("leaves printable unicode untouched", func(t *testing.T) {
		got := SanitizePaste("café 日本語")
		if got != "café 日本語" {
			t.Errorf("expected unicode preserved, got %q", got)
		}
	})

	t.Run("all-control input becomes empty", func(t *testing.T) {
		if got := SanitizePaste("\n\r\t"); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func TestTextInputEditing(t *testing.T) {
	t.Run("types characters and advances the caret", func(t *testing.T) {
		var ti TextInput
		for _, k := range []string{"f", "e", "a", "t"} {
			ti.HandleKey(k)
		}
		if ti.Value != "feat" || ti.Cursor() != 4 {
			t.Errorf("expected %q caret 4, got %q caret %d", "feat", ti.Value, ti.Cursor())
		}
	})

	t.Run("left/right move the caret and insert mid-string", func(t *testing.T) {
		var ti TextInput
		ti.SetValue("feat")
		ti.HandleKey("left")
		ti.HandleKey("left")
		if _, changed := ti.HandleKey("X"); !changed {
			t.Error("mid-string insert should report a change")
		}
		if ti.Value != "feXat" {
			t.Errorf("expected %q, got %q", "feXat", ti.Value)
		}
	})

	t.Run("backspace deletes before the caret, delete deletes at it", func(t *testing.T) {
		var ti TextInput
		ti.SetValue("feat")
		ti.Home()
		ti.MoveRight() // caret after 'f'
		if !ti.Backspace() {
			t.Error("backspace should report a change")
		}
		if ti.Value != "eat" { // 'f' removed
			t.Errorf("after backspace expected %q, got %q", "eat", ti.Value)
		}
		if !ti.Delete() {
			t.Error("delete should report a change")
		}
		if ti.Value != "at" { // 'e' at caret removed
			t.Errorf("after delete expected %q, got %q", "at", ti.Value)
		}
	})

	t.Run("caret clamps at both ends", func(t *testing.T) {
		var ti TextInput
		ti.SetValue("ab")
		ti.Home()
		if ti.Backspace() {
			t.Error("backspace at start should be a no-op")
		}
		ti.End()
		if ti.Delete() {
			t.Error("delete at end should be a no-op")
		}
	})

	t.Run("deletes a multi-byte rune without corrupting the rest", func(t *testing.T) {
		var ti TextInput
		ti.SetValue("café")
		if !ti.Backspace() {
			t.Error("backspace should report a change")
		}
		if ti.Value != "caf" {
			t.Errorf("expected %q, got %q", "caf", ti.Value)
		}
	})
}

func TestTextInputRenderCaret(t *testing.T) {
	// The caret cell must always be present so the user sees where typing
	// lands, both at the end of the text and in the middle of it.
	var ti TextInput
	ti.SetValue("ab")
	if end := ti.RenderPlain(); !strings.Contains(end, "a") || !strings.Contains(end, "b") {
		t.Errorf("end-caret render should contain the text, got %q", end)
	}
	ti.Home()
	mid := ti.RenderPlain()
	if !strings.Contains(mid, "a") || !strings.Contains(mid, "b") {
		t.Errorf("mid-caret render should contain the text, got %q", mid)
	}
}
