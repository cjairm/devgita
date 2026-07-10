package tuiinventory

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/cjairm/devgita/internal/inventory"
	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

// Options configures the dashboard's initial filter state.
type Options struct {
	ProblemsOnly bool   // pre-applied "problems only" filter (dg validate sets this)
	Category     string // pre-filter to a single category key (e.g. "fonts"); "" = all
}

type model struct {
	items []inventory.Item
	title string

	rows      []row
	cursor    int
	collapsed map[string]bool
	groupMode groupMode

	problemsOnly bool
	filtering    bool
	filter       string

	width, height int

	palette *tuicomponents.Palette
}

func newModel(items []inventory.Item, opts Options) model {
	title := "inventory"
	if opts.Category != "" {
		if label, ok := categoryLabels[opts.Category]; ok {
			title = label
		}
		var filtered []inventory.Item
		for _, it := range items {
			if it.Category == opts.Category {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

	m := model{
		items:        items,
		title:        title,
		collapsed:    map[string]bool{},
		groupMode:    groupByCategory,
		problemsOnly: opts.ProblemsOnly,
		palette:      tuicomponents.NewPalette(),
	}
	m.rebuildRows()
	return m
}

func (m *model) rebuildRows() {
	m.rows = buildRows(m.items, m.groupMode, m.collapsed, m.filter, m.problemsOnly)
	indices := itemIndices(m.rows)
	if len(indices) == 0 {
		m.cursor = 0
		return
	}
	for _, i := range indices {
		if i >= m.cursor {
			m.cursor = i
			return
		}
	}
	m.cursor = indices[len(indices)-1]
}

func (m *model) moveCursor(delta int) {
	indices := navigableIndices(m.rows, m.collapsed)
	if len(indices) == 0 {
		return
	}
	cur := -1
	for i, idx := range indices {
		if idx == m.cursor {
			cur = i
			break
		}
	}
	if cur == -1 {
		if delta > 0 {
			for _, idx := range indices {
				if idx > m.cursor {
					m.cursor = idx
					return
				}
			}
		} else {
			for i := len(indices) - 1; i >= 0; i-- {
				if indices[i] < m.cursor {
					m.cursor = indices[i]
					return
				}
			}
		}
		m.cursor = indices[0]
		return
	}
	cur = ((cur+delta)%len(indices) + len(indices)) % len(indices)
	m.cursor = indices[cur]
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.filtering {
		switch key {
		case "esc":
			m.filter = ""
			m.filtering = false
			m.rebuildRows()
		case "enter":
			m.filtering = false
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.rebuildRows()
			}
		default:
			if len(key) == 1 && key >= " " {
				m.filter += key
				m.rebuildRows()
			}
		}
		return m, nil
	}

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j":
		m.moveCursor(1)
	case "k":
		m.moveCursor(-1)
	case "h":
		if g := m.cursorGroup(); g != "" {
			m.collapsed[g] = true
			m.rebuildRows()
			m.landCursorOnGroup(g)
		}
	case "l":
		if g := m.cursorGroup(); g != "" {
			m.collapsed[g] = false
			m.rebuildRows()
		}
	case "/":
		m.filtering = true
	case "p":
		m.problemsOnly = !m.problemsOnly
		m.rebuildRows()
	case "g":
		if m.groupMode == groupByCategory {
			m.groupMode = groupByStatus
		} else {
			m.groupMode = groupByCategory
		}
		m.collapsed = map[string]bool{}
		m.rebuildRows()
	}
	return m, nil
}

func (m model) cursorGroup() string {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return ""
	}
	return m.rows[m.cursor].group
}

func (m *model) landCursorOnGroup(group string) {
	for i, r := range m.rows {
		if r.kind == rowGroup && r.group == group {
			m.cursor = i
			return
		}
	}
}

func (m model) View() tea.View {
	v := tea.NewView(m.renderContent())
	v.AltScreen = true
	return v
}

func (m model) renderContent() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	hint := m.renderHint(m.width)
	summary := m.renderSummary(m.width)

	// Reserve 1 line for hint, 1 for summary, 2 for the pane's own top/bottom border.
	viewportHeight := m.height - 4
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	start, end := visibleWindow(len(m.rows), m.cursor, viewportHeight)

	var lines []string
	for i := start; i < end; i++ {
		lines = append(lines, m.renderRow(i))
	}

	return m.palette.BorderedPane(m.title, m.width, lines) + "\n" + summary + "\n" + hint
}

// visibleWindow returns [start, end) into a rowsLen-length list such that the
// window has at most viewportHeight rows and always contains cursor.
func visibleWindow(rowsLen, cursor, viewportHeight int) (start, end int) {
	if rowsLen <= viewportHeight {
		return 0, rowsLen
	}
	start = cursor - viewportHeight/2
	if start < 0 {
		start = 0
	}
	end = start + viewportHeight
	if end > rowsLen {
		end = rowsLen
		start = end - viewportHeight
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

func (m model) renderRow(i int) string {
	r := m.rows[i]
	innerWidth := m.width - 2
	if r.kind == rowGroup {
		collapse := "▾"
		if m.collapsed[r.group] {
			collapse = "▸"
		}
		text := collapse + " " + r.group
		count := fmt.Sprintf("%d", r.count)
		pad := innerWidth - ansi.StringWidth(text) - ansi.StringWidth(count)
		if pad < 1 {
			pad = 1
		}
		plain := text + strings.Repeat(" ", pad) + count
		if i == m.cursor {
			return m.palette.Selected.Render(plain)
		}
		return m.palette.RepoHeader.Render(
			text,
		) + strings.Repeat(
			" ",
			pad,
		) + m.palette.SectionHead.Render(
			count,
		)
	}

	glyph := statusGlyph(r.item.State)
	name := "  " + glyph + " " + r.item.Name
	if i == m.cursor {
		plain := name
		if r.item.Source == "pre-existing" {
			plain += " (pre-existing)"
		}
		return m.palette.Selected.Render(plain)
	}
	line := "  " + statusDot(
		m.palette,
		r.item.State,
	) + " " + r.item.Name + sourceTag(
		m.palette,
		r.item.Source,
	)
	return line
}

func (m model) renderSummary(width int) string {
	categories := map[string]bool{}
	missing := 0
	for _, it := range m.items {
		categories[it.Category] = true
		if it.State == inventory.StateMissing {
			missing++
		}
	}
	text := fmt.Sprintf(
		"%d CATEGORIES · %d ITEMS · %d MISSING",
		len(categories),
		len(m.items),
		missing,
	)
	return m.palette.SectionHead.Render(ansi.Truncate(text, width, ""))
}

func (m model) renderHint(width int) string {
	if m.filtering {
		hint := "filter: " + m.filter + "█  · esc: clear · enter: keep"
		return m.palette.HintDesc.Render(ansi.Truncate(hint, width, ""))
	}
	hints := []tuicomponents.KeyHint{
		{Key: "j/k", Desc: "move"},
		{Key: "h/l", Desc: "collapse/expand"},
		{Key: "/", Desc: "filter"},
		{Key: "p", Desc: "problems"},
		{Key: "g", Desc: "group"},
		{Key: "q", Desc: "quit"},
	}
	return m.palette.HintBar(hints, width)
}
