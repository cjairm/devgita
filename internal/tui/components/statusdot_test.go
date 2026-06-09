package tuicomponents_test

import (
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/tooling/worktree"
	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

func TestSessionStateFromWorktree(t *testing.T) {
	cases := []struct {
		name        string
		status      worktree.WorktreeStatus
		needsReview bool
		dirtyCount  int
		want        tuicomponents.SessionState
	}{
		{
			name:   "active window no review no dirty",
			status: worktree.WorktreeStatus{WindowActive: true},
			want:   tuicomponents.StateRunning,
		},
		{
			name:        "active window needs review",
			status:      worktree.WorktreeStatus{WindowActive: true},
			needsReview: true,
			want:        tuicomponents.StateNeedsReview,
		},
		{
			name:       "inactive window dirty",
			status:     worktree.WorktreeStatus{WindowActive: false},
			dirtyCount: 3,
			want:       tuicomponents.StateDirty,
		},
		{
			name:   "inactive window no dirty",
			status: worktree.WorktreeStatus{WindowActive: false},
			want:   tuicomponents.StateNoSession,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tuicomponents.SessionStateFromWorktree(tc.status, tc.needsReview, tc.dirtyCount)
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestStatusGlyphNoANSI(t *testing.T) {
	p := tuicomponents.NewPalette()
	cases := []struct {
		state tuicomponents.SessionState
		glyph string
	}{
		{tuicomponents.StateRunning, "●"},
		{tuicomponents.StateNeedsReview, "◆"},
		{tuicomponents.StateDirty, "●"},
		{tuicomponents.StateNoSession, "○"},
	}
	for _, tc := range cases {
		got := p.StatusGlyph(tc.state)
		if got != tc.glyph {
			t.Errorf("state %d: got %q want %q", tc.state, got, tc.glyph)
		}
		if strings.ContainsRune(got, '\x1b') {
			t.Errorf("state %d: StatusGlyph must not contain ANSI escape bytes", tc.state)
		}
	}
}

func TestStatusDotContainsGlyph(t *testing.T) {
	p := tuicomponents.NewPalette()
	cases := []struct {
		state tuicomponents.SessionState
		glyph string
	}{
		{tuicomponents.StateRunning, "●"},
		{tuicomponents.StateNeedsReview, "◆"},
		{tuicomponents.StateDirty, "●"},
		{tuicomponents.StateNoSession, "○"},
	}
	for _, tc := range cases {
		got := p.StatusDot(tc.state)
		if !strings.Contains(got, tc.glyph) {
			t.Errorf("state %d: StatusDot %q does not contain glyph %q", tc.state, got, tc.glyph)
		}
	}
}
