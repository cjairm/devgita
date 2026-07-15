package tuicomponents

import "testing"

func TestMoveCursor(t *testing.T) {
	indices := []int{1, 2, 4}

	t.Run("moves within indices", func(t *testing.T) {
		if got := MoveCursor(indices, 1, 1); got != 2 {
			t.Errorf("expected 2, got %d", got)
		}
	})

	t.Run("wraps forward past the end", func(t *testing.T) {
		if got := MoveCursor(indices, 4, 1); got != 1 {
			t.Errorf("expected wrap to 1, got %d", got)
		}
	})

	t.Run("wraps backward past the start", func(t *testing.T) {
		if got := MoveCursor(indices, 1, -1); got != 4 {
			t.Errorf("expected wrap to 4, got %d", got)
		}
	})

	t.Run("jumps forward from a non-navigable row", func(t *testing.T) {
		if got := MoveCursor(indices, 3, 1); got != 4 {
			t.Errorf("expected 4, got %d", got)
		}
	})

	t.Run("jumps backward from a non-navigable row", func(t *testing.T) {
		if got := MoveCursor(indices, 3, -1); got != 2 {
			t.Errorf("expected 2, got %d", got)
		}
	})

	t.Run("wraps when no navigable row exists in the direction", func(t *testing.T) {
		if got := MoveCursor(indices, 5, 1); got != 1 {
			t.Errorf("expected wrap to 1, got %d", got)
		}
		if got := MoveCursor(indices, 0, -1); got != 4 {
			t.Errorf("expected wrap to 4, got %d", got)
		}
	})

	t.Run("empty indices returns cursor unchanged", func(t *testing.T) {
		if got := MoveCursor(nil, 3, 1); got != 3 {
			t.Errorf("expected 3, got %d", got)
		}
	})
}

func TestClampCursor(t *testing.T) {
	indices := []int{1, 2, 4}

	t.Run("keeps cursor already on a row", func(t *testing.T) {
		if got := ClampCursor(indices, 2); got != 2 {
			t.Errorf("expected 2, got %d", got)
		}
	})

	t.Run("advances to the next row", func(t *testing.T) {
		if got := ClampCursor(indices, 3); got != 4 {
			t.Errorf("expected 4, got %d", got)
		}
	})

	t.Run("falls back to the last row", func(t *testing.T) {
		if got := ClampCursor(indices, 9); got != 4 {
			t.Errorf("expected 4, got %d", got)
		}
	})

	t.Run("empty indices returns 0", func(t *testing.T) {
		if got := ClampCursor(nil, 7); got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})
}
