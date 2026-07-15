package tuicomponents

// MoveCursor moves a list cursor by delta within the navigable row indices,
// wrapping at both ends. When the cursor is not currently on a navigable row
// (e.g. it sits on a group header), it jumps to the nearest navigable row in
// the direction of travel, wrapping if none exists. Returns cursor unchanged
// when indices is empty.
func MoveCursor(indices []int, cursor, delta int) int {
	if len(indices) == 0 {
		return cursor
	}
	cur := -1
	for i, idx := range indices {
		if idx == cursor {
			cur = i
			break
		}
	}
	if cur == -1 {
		if delta > 0 {
			for _, idx := range indices {
				if idx > cursor {
					return idx
				}
			}
			return indices[0]
		}
		for i := len(indices) - 1; i >= 0; i-- {
			if indices[i] < cursor {
				return indices[i]
			}
		}
		return indices[len(indices)-1]
	}
	cur = ((cur+delta)%len(indices) + len(indices)) % len(indices)
	return indices[cur]
}

// ClampCursor keeps a list cursor on a valid row after the rows are rebuilt:
// it returns the first navigable index at or after cursor, falling back to
// the last one. Returns 0 when indices is empty.
func ClampCursor(indices []int, cursor int) int {
	if len(indices) == 0 {
		return 0
	}
	for _, i := range indices {
		if i >= cursor {
			return i
		}
	}
	return indices[len(indices)-1]
}
