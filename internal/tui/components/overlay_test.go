package tuicomponents

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestOverlay_PlainBackground(t *testing.T) {
	background := strings.Join([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	}, "\n")
	popup := strings.Join([]string{
		"XXXX",
		"XXXX",
	}, "\n")

	got := Overlay(background, popup, 10, 6)
	lines := strings.Split(got, "\n")

	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d: %q", len(lines), got)
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 10 {
			t.Errorf("row %d: expected display width 10, got %d in %q", i, w, l)
		}
	}

	// popup (2 rows x 4 cols) is centered inside a 10x6 screen: rows 2-3,
	// columns 3-6.
	for _, i := range []int{2, 3} {
		stripped := ansi.Strip(lines[i])
		if !strings.Contains(stripped, "XXXX") {
			t.Errorf("row %d missing popup content: %q", i, stripped)
		}
		if !strings.HasPrefix(stripped, "aaa") {
			t.Errorf("row %d should keep background on the left: %q", i, stripped)
		}
		if !strings.HasSuffix(stripped, "aaa") {
			t.Errorf("row %d should keep background on the right: %q", i, stripped)
		}
	}
	for _, i := range []int{0, 1, 4, 5} {
		if lines[i] != "aaaaaaaaaa" {
			t.Errorf("row %d should be untouched background, got %q", i, lines[i])
		}
	}
}

func TestOverlay_AnsiStyledBackground(t *testing.T) {
	red := "\x1b[31m"
	row := red + "aaaaaaaaaa" + ansiReset
	background := strings.Join([]string{row, row, row, row, row, row}, "\n")
	popup := strings.Join([]string{"XXXX", "XXXX"}, "\n")

	got := Overlay(background, popup, 10, 6)
	lines := strings.Split(got, "\n")

	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d: %q", len(lines), got)
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 10 {
			t.Errorf("row %d: expected display width 10, got %d in %q", i, w, l)
		}
		if stripped := ansi.Strip(l); ansi.StringWidth(stripped) != 10 {
			t.Errorf("row %d: stripped width mismatch: %q", i, stripped)
		}
	}

	// Rows under the popup must reset styling right before the popup segment
	// so the background's red doesn't bleed into the popup content.
	for _, i := range []int{2, 3} {
		stripped := ansi.Strip(lines[i])
		if !strings.Contains(stripped, "XXXX") {
			t.Errorf("row %d missing popup content: %q", i, stripped)
		}
		idx := strings.Index(lines[i], "XXXX")
		if idx < len(ansiReset) || lines[i][idx-len(ansiReset):idx] != ansiReset {
			t.Errorf("row %d: expected reset immediately before popup content: %q", i, lines[i])
		}
	}
	// Untouched rows keep their original styling verbatim.
	for _, i := range []int{0, 1, 4, 5} {
		if lines[i] != row {
			t.Errorf("row %d should be untouched styled background, got %q", i, lines[i])
		}
	}
}

func TestOverlay_PopupTallerThanScreen(t *testing.T) {
	background := strings.Repeat("b", 20)
	var popupRows []string
	for range 50 {
		popupRows = append(popupRows, "P")
	}
	popup := strings.Join(popupRows, "\n")

	got := Overlay(background, popup, 20, 5)
	lines := strings.Split(got, "\n")

	if len(lines) != 5 {
		t.Fatalf("expected output clipped to screen height 5, got %d lines", len(lines))
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 20 {
			t.Errorf("row %d: expected display width 20, got %d in %q", i, w, l)
		}
	}
	// Popup fills every row since it's taller than the screen.
	for i, l := range lines {
		if !strings.Contains(l, "P") {
			t.Errorf("row %d: expected popup content, got %q", i, l)
		}
	}
}

func TestOverlay_PopupWiderThanScreen(t *testing.T) {
	background := strings.Repeat("b", 5)
	popup := strings.Repeat("P", 50)

	got := Overlay(background, popup, 5, 3)
	lines := strings.Split(got, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 5 {
			t.Errorf("row %d: expected display width 5 (clipped), got %d in %q", i, w, l)
		}
	}
	// Popup fills the entire width since it's wider than the screen.
	middle := ansi.Strip(lines[1])
	if !strings.Contains(middle, "PPPPP") {
		t.Errorf("middle row should be entirely popup content: %q", middle)
	}
}

func TestOverlay_TinyAndZeroSizes(t *testing.T) {
	cases := []struct {
		name          string
		width, height int
	}{
		{"zero width", 0, 5},
		{"zero height", 5, 0},
		{"negative width", -3, 5},
		{"negative height", 5, -3},
		{"both zero", 0, 0},
		{"1x1", 1, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Overlay panicked for %s: %v", tc.name, r)
				}
			}()
			got := Overlay("background text", "popup text", tc.width, tc.height)
			if tc.width <= 0 || tc.height <= 0 {
				if got != "" {
					t.Errorf("expected empty output for %s, got %q", tc.name, got)
				}
				return
			}
			lines := strings.Split(got, "\n")
			if len(lines) != tc.height {
				t.Errorf("%s: expected %d lines, got %d", tc.name, tc.height, len(lines))
			}
			for _, l := range lines {
				if w := ansi.StringWidth(l); w != tc.width {
					t.Errorf("%s: expected width %d, got %d in %q", tc.name, tc.width, w, l)
				}
			}
		})
	}
}

func TestOverlay_BackgroundShorterThanScreen(t *testing.T) {
	background := "short"
	popup := "X"

	got := Overlay(background, popup, 10, 4)
	lines := strings.Split(got, "\n")

	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), got)
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 10 {
			t.Errorf("row %d: expected width 10, got %d in %q", i, w, l)
		}
	}
	// Row 0 is the one real background line and sits outside the popup's
	// vertical span (offsetY = (4-1)/2 = 1), so it must be untouched.
	if stripped := ansi.Strip(lines[0]); stripped != "short     " {
		t.Errorf("row 0 should be untouched background, got %q", stripped)
	}
	// Row 1 has no background line to draw from; it should be blank padding
	// with the popup composited on top, not a panic or garbage.
	if !strings.Contains(ansi.Strip(lines[1]), "X") {
		t.Errorf("row 1 should still show popup content over blank padding: %q", lines[1])
	}
}
