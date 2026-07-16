package tuicomponents

import (
	"sort"
	"strings"
	"unicode/utf8"
)

// FuzzyPickerAction reports what happened after a keypress, letting a caller
// (a bubbletea model) tell "still open" apart from "user picked something"
// and "user cancelled".
type FuzzyPickerAction int

const (
	FuzzyPickerNone FuzzyPickerAction = iota
	FuzzyPickerSelected
	FuzzyPickerCancelled
)

// FuzzyPickerResult is returned by HandleKey after processing one keypress.
type FuzzyPickerResult struct {
	Action FuzzyPickerAction
	Item   PaletteItem // valid only when Action == FuzzyPickerSelected
}

// FuzzyPicker is a searchable single-select list: typed characters fuzzy-
// filter and rank the item set (see FuzzyMatch), up/down/j/k move the
// cursor, enter selects, esc cancels.
type FuzzyPicker struct {
	Title string

	palette  *Palette
	items    []PaletteItem
	query    string
	filtered []int // indices into items, ranked best match first
	cursor   int   // index into filtered, not items
}

// NewFuzzyPicker builds a picker over items, titled for its bordered box.
func NewFuzzyPicker(title string, items []PaletteItem) *FuzzyPicker {
	p := &FuzzyPicker{Title: title, palette: NewPalette()}
	p.SetItems(items)
	return p
}

// SetItems replaces the full (unfiltered) item set and re-filters against the
// current query.
func (p *FuzzyPicker) SetItems(items []PaletteItem) {
	p.items = items
	p.refilter()
}

// Query returns the current search text.
func (p *FuzzyPicker) Query() string {
	return p.query
}

// InsertText inserts pasted text into the query in one shot and refilters
// once, the paste counterpart to HandleKey: a tea.PasteMsg carries the whole
// clipboard content as one string, and HandleKey's default case only accepts
// single-rune keys, so a multi-rune paste would otherwise be silently
// dropped rune-by-rune.
func (p *FuzzyPicker) InsertText(text string) {
	text = SanitizePaste(text)
	if text == "" {
		return
	}
	p.query += text
	p.refilter()
}

// Selected returns the item under the cursor in the filtered list, if any.
func (p *FuzzyPicker) Selected() (PaletteItem, bool) {
	if p.cursor < 0 || p.cursor >= len(p.filtered) {
		return PaletteItem{}, false
	}
	return p.items[p.filtered[p.cursor]], true
}

// HandleKey processes one keypress and reports what happened.
//
// Unlike the dashboard, this is a type-to-filter picker: every printable
// character, including "j" and "k", must be available for the query text, or
// names starting with those letters ("kafka", "json-tool") could never be
// filtered down. So cursor nav is arrow keys only (plus ctrl+j/ctrl+k as
// chorded vim-style equivalents, which don't collide with typing since they
// carry the ctrl modifier), and bare "j"/"k" fall through to the default case
// like any other character.
func (p *FuzzyPicker) HandleKey(key string) FuzzyPickerResult {
	switch key {
	case "esc":
		return FuzzyPickerResult{Action: FuzzyPickerCancelled}
	case "enter":
		if item, ok := p.Selected(); ok {
			return FuzzyPickerResult{Action: FuzzyPickerSelected, Item: item}
		}
	case "backspace":
		if len(p.query) > 0 {
			p.query = TrimLastRune(p.query)
			p.refilter()
		}
	case "up", "ctrl+k":
		p.cursor = MoveCursor(p.navIndices(), p.cursor, -1)
	case "down", "ctrl+j":
		p.cursor = MoveCursor(p.navIndices(), p.cursor, 1)
	default:
		// len() counts bytes, so a non-ASCII rune (e.g. "é", "日") arriving as a
		// multi-byte UTF-8 string would fail this check and be silently
		// dropped. Every control-key string this switch matches ("esc",
		// "backspace", "ctrl+j", ...) has more than one rune, so counting
		// runes still excludes them from falling into the default case.
		if utf8.RuneCountInString(key) == 1 && key >= " " {
			p.query += key
			p.refilter()
		}
	}
	return FuzzyPickerResult{Action: FuzzyPickerNone}
}

// navIndices is the identity index set [0, len(filtered)) — every filtered
// row is navigable, so it's handed to MoveCursor as-is (unlike the dashboard,
// which skips group-header rows).
func (p *FuzzyPicker) navIndices() []int {
	idx := make([]int, len(p.filtered))
	for i := range idx {
		idx[i] = i
	}
	return idx
}

// refilter re-ranks items against the current query and resets the cursor to
// the top match. Resetting (rather than trying to preserve the old cursor
// position) is simplest and matches user expectation: the best match for the
// new query should be highlighted, not whatever row happened to be at that
// position before.
func (p *FuzzyPicker) refilter() {
	type ranked struct {
		idx  int
		rank FuzzyRank
	}
	matches := make([]ranked, 0, len(p.items))
	for i, item := range p.items {
		if rank := FuzzyMatch(p.query, item.Command); rank != FuzzyNoMatch {
			matches = append(matches, ranked{idx: i, rank: rank})
		}
	}
	sort.SliceStable(matches, func(a, b int) bool {
		return matches[a].rank > matches[b].rank
	})

	p.filtered = make([]int, len(matches))
	for i, m := range matches {
		p.filtered[i] = m.idx
	}
	p.cursor = 0
}

// View renders the picker as a bordered box (title + CommandPalette body),
// sized to width.
func (p *FuzzyPicker) View(width int) string {
	items := make([]PaletteItem, len(p.filtered))
	for i, idx := range p.filtered {
		items[i] = p.items[idx]
	}

	// Mirror BorderedPane's own width clamp exactly: it floors its outer
	// width to 6 before computing inner := width-2. If this used the raw
	// (unclamped) width instead, CommandPalette could be handed a width
	// smaller than the interior BorderedPane actually draws, silently
	// truncating or blanking content that would otherwise fit.
	innerWidth := width
	if innerWidth < 6 {
		innerWidth = 6
	}
	innerWidth -= 2

	body := p.palette.CommandPalette(p.query, items, p.cursor, innerWidth)
	if len(items) == 0 {
		body += "\n" + p.palette.Inactive.Render("(no matches)")
	}

	lines := strings.Split(body, "\n")
	return p.palette.BorderedPane(p.Title, width, lines)
}
