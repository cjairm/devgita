package tuicomponents

import "testing"

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

func TestTrimLastRune(t *testing.T) {
	t.Run("empty string is a no-op", func(t *testing.T) {
		if got := TrimLastRune(""); got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("removes one ASCII rune", func(t *testing.T) {
		if got := TrimLastRune("feat"); got != "fea" {
			t.Errorf("expected %q, got %q", "fea", got)
		}
	})

	t.Run("removes one multi-byte rune without corrupting the rest", func(t *testing.T) {
		if got := TrimLastRune("café"); got != "caf" {
			t.Errorf("expected %q, got %q", "caf", got)
		}
	})

	t.Run("removes one emoji rune", func(t *testing.T) {
		if got := TrimLastRune("hi🎉"); got != "hi" {
			t.Errorf("expected %q, got %q", "hi", got)
		}
	})
}
