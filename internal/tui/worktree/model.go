// Package tuiworktree provides the Bubble Tea TUI dashboard for dg wt ui.
package tuiworktree

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/tooling/worktree"
)

const (
	minLeftPaneWidth     = 20
	defaultLeftPaneWidth = 35
	maxLeftPaneWidthPct  = 0.60
	dividerWidth         = 1
	refreshInterval      = 1500 * time.Millisecond
	navPauseThreshold    = 500 * time.Millisecond
	maxDiffBytes         = 64 * 1024
)

type tabKind int

const (
	tabAgent tabKind = iota
	tabDiff
)

// --- Messages ---

type (
	statusesMsg     []worktree.WorktreeStatus
	agentContentMsg string
	agentOfflineMsg struct{}
	diffMsg         struct {
		content               string
		files, added, removed int
	}
)

type (
	tickMsg   time.Time
	statusMsg string
)

// --- Model ---

// Model is the Bubble Tea model for the worktree TUI dashboard.
type Model struct {
	mgr     *worktree.WorktreeManager
	tmuxApp *tmux.Tmux
	gitApp  *git.Git
	gc      *config.GlobalConfig

	statuses     []worktree.WorktreeStatus
	rows         []row
	cursor       int // index into rows (always a rowWorktree)
	collapsed    map[string]bool
	allCollapsed bool

	activeTab    tabKind
	agentContent string
	agentHash    string
	agentOffline bool

	diffContent string
	diffFiles   int
	diffAdded   int
	diffRemoved int
	diffScroll  int

	filtering bool
	filter    string

	status string

	width  int
	height int
	styles Styles

	lastNavTime time.Time

	leftPaneWidth int

	dragging   bool
	dragStartX int

	pendingDelete string // "repo/name" or ""
	showHelp      bool

	// Injected I/O seams (overridable in tests)
	captureFn  func(session, window string) (string, error)
	diffFn     func(path string) (string, error)
	diffStatFn func(path string) (files, added, removed int, err error)
	attachFn   func(session, window string) error
	removeFn   func(repo, name string, force bool) error
	repairFn   func(repo, name string, coder worktree.AICoder) error
}

func newModel(
	mgr *worktree.WorktreeManager,
	tmuxApp *tmux.Tmux,
	gitApp *git.Git,
	gc *config.GlobalConfig,
) Model {
	m := Model{
		mgr:           mgr,
		tmuxApp:       tmuxApp,
		gitApp:        gitApp,
		gc:            gc,
		collapsed:     map[string]bool{},
		styles:        newStyles(),
		leftPaneWidth: defaultLeftPaneWidth,
	}
	m.captureFn = func(session, window string) (string, error) {
		return tmuxApp.CapturePane(session, window)
	}
	m.diffFn = func(path string) (string, error) {
		return gitApp.Diff(path)
	}
	m.diffStatFn = func(path string) (files, added, removed int, err error) {
		return gitApp.DiffStat(path)
	}
	m.attachFn = func(session, window string) error {
		return tmuxApp.SwitchToWindow(session, window)
	}
	m.removeFn = func(repo, name string, force bool) error {
		return mgr.RemoveInRepo(repo, name, force)
	}
	m.repairFn = func(repo, name string, coder worktree.AICoder) error {
		return mgr.RepairInRepo(repo, name, coder)
	}
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadCmd(),
		m.tickCmd(),
	)
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		statuses, err := m.mgr.List()
		if err != nil {
			return statusMsg("failed to list worktrees: " + err.Error())
		}
		return statusesMsg(statuses)
	}
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) captureAgentCmd(s worktree.WorktreeStatus) tea.Cmd {
	window := worktree.GetWindowName(s.Name)
	if m.tmuxApp == nil {
		return func() tea.Msg { return agentOfflineMsg{} }
	}
	session, ok := m.tmuxApp.WindowSession(window)
	if !ok {
		return func() tea.Msg { return agentOfflineMsg{} }
	}
	capFn := m.captureFn
	return func() tea.Msg {
		content, err := capFn(session, window)
		if err != nil {
			return agentOfflineMsg{}
		}
		return agentContentMsg(content)
	}
}

func (m Model) computeDiffCmd(s worktree.WorktreeStatus) tea.Cmd {
	df := m.diffFn
	dsf := m.diffStatFn
	path := s.Path
	return func() tea.Msg {
		content, err := df(path)
		if err != nil {
			content = "(diff unavailable: " + err.Error() + ")"
		}
		if len(content) > maxDiffBytes {
			content = content[:maxDiffBytes] + "\n... (truncated)"
		}
		files, added, removed, _ := dsf(path)
		return diffMsg{content: content, files: files, added: added, removed: removed}
	}
}

func (m Model) selectedStatus() (worktree.WorktreeStatus, bool) {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return worktree.WorktreeStatus{}, false
	}
	r := m.rows[m.cursor]
	if r.kind != rowWorktree {
		return worktree.WorktreeStatus{}, false
	}
	return r.status, true
}

func (m *Model) rebuildRows() {
	m.rows = buildRows(m.statuses, m.collapsed, m.filter)
	// Ensure cursor points to a worktree row
	indices := worktreeIndices(m.rows)
	if len(indices) == 0 {
		m.cursor = 0
		return
	}
	// Keep cursor on a valid worktree row
	found := false
	for _, i := range indices {
		if i >= m.cursor {
			m.cursor = i
			found = true
			break
		}
	}
	if !found {
		m.cursor = indices[len(indices)-1]
	}
}

func (m *Model) moveCursor(delta int) {
	indices := worktreeIndices(m.rows)
	if len(indices) == 0 {
		return
	}
	// Find current position in indices
	cur := 0
	for i, idx := range indices {
		if idx == m.cursor {
			cur = i
			break
		}
	}
	cur += delta
	// Wrap around
	cur = ((cur % len(indices)) + len(indices)) % len(indices)
	m.cursor = indices[cur]
}

func (m Model) safeMaxLeft() int {
	return max(int(float64(m.width)*maxLeftPaneWidthPct), minLeftPaneWidth)
}

func (m Model) rightPaneWidth() int {
	return max(m.width-m.leftPaneWidth-dividerWidth, 0)
}

func contentHash(s string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.leftPaneWidth = min(m.leftPaneWidth, m.safeMaxLeft())
		return m, nil

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Button == tea.MouseLeft && mouse.X == m.leftPaneWidth {
			m.dragging = true
			m.dragStartX = mouse.X
		}
		return m, nil

	case tea.MouseMotionMsg:
		if m.dragging {
			mouse := msg.Mouse()
			maxLeft := m.safeMaxLeft()
			newWidth := min(max(mouse.X, minLeftPaneWidth), maxLeft)
			m.leftPaneWidth = newWidth
		}
		return m, nil

	case tea.MouseReleaseMsg:
		m.dragging = false
		return m, nil

	case statusesMsg:
		m.statuses = []worktree.WorktreeStatus(msg)
		m.rebuildRows()
		if sel, ok := m.selectedStatus(); ok {
			return m, tea.Batch(m.captureAgentCmd(sel), m.computeDiffCmd(sel))
		}
		return m, nil

	case agentContentMsg:
		hash := contentHash(string(msg))
		if hash != m.agentHash {
			m.agentContent = string(msg)
			m.agentHash = hash
			m.agentOffline = false
		}
		return m, nil

	case agentOfflineMsg:
		m.agentOffline = true
		return m, nil

	case diffMsg:
		m.diffContent = msg.content
		m.diffFiles = msg.files
		m.diffAdded = msg.added
		m.diffRemoved = msg.removed
		return m, nil

	case tickMsg:
		sel, ok := m.selectedStatus()
		if !ok || m.agentOffline {
			return m, m.tickCmd()
		}
		if time.Since(m.lastNavTime) > navPauseThreshold {
			return m, tea.Batch(m.captureAgentCmd(sel), m.tickCmd())
		}
		return m, m.tickCmd()

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Help overlay: any key dismisses it
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}

	if m.filtering {
		switch key {
		case "esc":
			m.filter = ""
			m.filtering = false
			m.rebuildRows()
			return m, nil
		case "enter":
			m.filtering = false
			return m, nil
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.rebuildRows()
			}
			return m, nil
		default:
			if len(key) == 1 && key >= " " {
				m.filter += key
				m.rebuildRows()
			}
		}
		return m, nil
	}

	// Clear pending delete on any non-d key
	if key != "d" && m.pendingDelete != "" {
		m.pendingDelete = ""
	}

	switch key {
	case "?":
		m.showHelp = true
		return m, nil

	case "q", "ctrl+c":
		return m, tea.Quit

	case "j":
		m.lastNavTime = time.Now()
		m.pendingDelete = ""
		m.diffScroll = 0
		m.agentOffline = false
		m.moveCursor(1)
		if sel, ok := m.selectedStatus(); ok {
			return m, tea.Batch(m.captureAgentCmd(sel), m.computeDiffCmd(sel))
		}
		return m, nil

	case "k":
		m.lastNavTime = time.Now()
		m.pendingDelete = ""
		m.diffScroll = 0
		m.agentOffline = false
		m.moveCursor(-1)
		if sel, ok := m.selectedStatus(); ok {
			return m, tea.Batch(m.captureAgentCmd(sel), m.computeDiffCmd(sel))
		}
		return m, nil

	case "h":
		if sel, ok := m.selectedStatus(); ok {
			m.collapsed[sel.Repo] = true
			m.rebuildRows()
		}
		return m, nil

	case "l":
		if sel, ok := m.selectedStatus(); ok {
			m.collapsed[sel.Repo] = false
			m.rebuildRows()
		}
		return m, nil

	case "z":
		m.allCollapsed = !m.allCollapsed
		for _, s := range m.statuses {
			m.collapsed[s.Repo] = m.allCollapsed
		}
		m.rebuildRows()
		return m, nil

	case "tab":
		if m.activeTab == tabAgent {
			m.activeTab = tabDiff
		} else {
			m.activeTab = tabAgent
		}
		m.diffScroll = 0
		return m, nil

	case "/":
		m.filtering = true
		return m, nil

	case "enter":
		return m.handleAttach()

	case "d":
		return m.handleDelete()

	case "r":
		return m.handleRepair()

	case "ctrl+d":
		if m.activeTab == tabDiff {
			pageSize := m.height - 4
			m.diffScroll += pageSize
		}
		return m, nil

	case "ctrl+u":
		if m.activeTab == tabDiff {
			pageSize := m.height - 4
			m.diffScroll = max(m.diffScroll-pageSize, 0)
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleAttach() (tea.Model, tea.Cmd) {
	sel, ok := m.selectedStatus()
	if !ok {
		return m, nil
	}

	if os.Getenv("TMUX") == "" {
		m.status = "not inside tmux; run `dg wt ui` from a tmux session"
		return m, nil
	}

	window := worktree.GetWindowName(sel.Name)
	attachFn := m.attachFn
	repairFn := m.repairFn
	tmuxApp := m.tmuxApp
	gc := m.gc

	return m, func() tea.Msg {
		session, ok := tmuxApp.WindowSession(window)
		if ok {
			if err := attachFn(session, window); err != nil {
				return statusMsg("attach failed: " + err.Error())
			}
			return tea.QuitMsg{}
		}
		// Auto-repair: window missing
		alias := worktree.ResolveAIAlias("", gc)
		coder, err := worktree.ResolveAICoder(alias)
		if err != nil {
			return statusMsg("repair failed: " + err.Error())
		}
		if err := repairFn(sel.Repo, sel.Name, coder); err != nil {
			return statusMsg("repair failed: " + err.Error())
		}
		session, ok = tmuxApp.WindowSession(window)
		if !ok {
			return statusMsg("repair succeeded but window not found")
		}
		if err := attachFn(session, window); err != nil {
			return statusMsg("attach after repair failed: " + err.Error())
		}
		return tea.QuitMsg{}
	}
}

func (m Model) handleDelete() (tea.Model, tea.Cmd) {
	sel, ok := m.selectedStatus()
	if !ok {
		return m, nil
	}

	key := sel.Repo + "/" + sel.Name

	if m.pendingDelete == key {
		// Second d: execute delete
		m.pendingDelete = ""
		removeFn := m.removeFn
		repo := sel.Repo
		name := sel.Name
		statuses := m.statuses
		return m, func() tea.Msg {
			if err := removeFn(repo, name, true); err != nil {
				return statusMsg("delete failed: " + err.Error())
			}
			// Drop from statuses
			var updated []worktree.WorktreeStatus
			for _, s := range statuses {
				if s.Repo != repo || s.Name != name {
					updated = append(updated, s)
				}
			}
			return statusesMsg(updated)
		}
	}

	// First d: arm
	m.pendingDelete = key
	return m, nil
}

func (m Model) handleRepair() (tea.Model, tea.Cmd) {
	sel, ok := m.selectedStatus()
	if !ok {
		return m, nil
	}
	repairFn := m.repairFn
	gc := m.gc
	repo := sel.Repo
	name := sel.Name
	return m, func() tea.Msg {
		alias := worktree.ResolveAIAlias("", gc)
		coder, err := worktree.ResolveAICoder(alias)
		if err != nil {
			return statusMsg("repair failed: " + err.Error())
		}
		if err := repairFn(repo, name, coder); err != nil {
			return statusMsg("repair failed: " + err.Error())
		}
		return statusMsg("repaired: " + name)
	}
}

// View implements tea.Model.
func (m Model) View() tea.View {
	content := m.renderContent()
	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) renderContent() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.showHelp {
		return m.renderHelpOverlay()
	}

	rpw := m.rightPaneWidth()
	lpw := m.leftPaneWidth

	// Narrow terminal fallback
	if rpw <= 0 {
		left := m.renderLeft(m.width - 1)
		hint := m.renderHint(m.width)
		status := m.renderStatus(m.width)
		lines := strings.Split(left, "\n")
		// Trim to height
		maxLines := max(m.height-2, 0)
		if len(lines) > maxLines {
			lines = lines[:maxLines]
		}
		return strings.Join(lines, "\n") + "\n" + hint + "\n" + status
	}

	left := m.renderLeft(lpw)
	divider := m.renderDivider(m.height - 2)
	right := m.renderRight(rpw)
	hint := m.renderHint(m.width)
	status := m.renderStatus(m.width)

	// Join left + divider + right horizontally
	leftLines := padLines(strings.Split(left, "\n"), m.height-2, lpw)
	rightLines := padLines(strings.Split(right, "\n"), m.height-2, rpw)
	divLines := strings.Split(divider, "\n")

	var combined []string
	for i := range m.height - 2 {
		l := ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		d := ""
		if i < len(divLines) {
			d = divLines[i]
		}
		r := ""
		if i < len(rightLines) {
			r = rightLines[i]
		}
		combined = append(combined, l+d+r)
	}

	body := strings.Join(combined, "\n")
	return body + "\n" + hint + "\n" + status
}

func (m Model) renderLeft(width int) string {
	var sb strings.Builder
	for i, r := range m.rows {
		var line string
		if r.kind == rowRepo {
			collapse := "▼"
			if m.collapsed[r.repo] {
				collapse = "▶"
			}
			text := collapse + " " + r.repo
			line = m.styles.RepoHeader.Render(
				text,
			) + strings.Repeat(
				" ",
				max(0, width-ansi.StringWidth(text)),
			)
		} else {
			g := glyphFor(r.status)
			name := ansi.Truncate(r.status.Name, width-4, "")
			pendingKey := r.status.Repo + "/" + r.status.Name
			plainText := "  " + g + " " + name
			padding := strings.Repeat(" ", max(0, width-ansi.StringWidth(plainText)))

			if i == m.cursor {
				if m.pendingDelete == pendingKey {
					line = m.styles.ArmedRow.Render(plainText + padding)
				} else {
					line = m.styles.SelectedRow.Render(plainText + padding)
				}
			} else {
				glyphStyle := m.styles.InactiveGlyph
				if r.status.WindowActive {
					glyphStyle = m.styles.ActiveGlyph
				}
				line = "  " + glyphStyle.Render(g) + " " + name + padding
			}
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func (m Model) renderDivider(height int) string {
	divChar := m.styles.Divider.Render("│")
	lines := make([]string, height)
	for i := range lines {
		lines[i] = divChar
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderRight(width int) string {
	// Tab bar
	agentLabel := "Agent"
	diffLabel := "Diff"
	if m.activeTab == tabAgent {
		agentLabel = m.styles.ActiveTab.Render(agentLabel)
		diffLabel = m.styles.InactiveTab.Render(diffLabel)
	} else {
		agentLabel = m.styles.InactiveTab.Render(agentLabel)
		diffLabel = m.styles.ActiveTab.Render(diffLabel)
	}
	tabBar := agentLabel + "  " + diffLabel
	contentHeight := max(m.height-4, 0) // height minus hint, status, tabbar, blank line

	var content string
	if m.activeTab == tabAgent {
		if m.agentOffline {
			content = m.styles.OfflinePlaceholder.Render("⟂ window offline — press r to repair")
		} else {
			content = m.renderAgentContent(width, contentHeight)
		}
	} else {
		content = m.renderDiffContent(width, contentHeight)
	}

	return tabBar + "\n" + content
}

func (m Model) renderAgentContent(width, height int) string {
	if m.agentContent == "" {
		return m.styles.InactiveGlyph.Render("(loading...)")
	}
	lines := strings.Split(m.agentContent, "\n")
	var truncated []string
	for _, line := range lines {
		truncated = append(truncated, ansi.Truncate(line, width, ""))
	}
	// Show last `height` lines (most recent output)
	if len(truncated) > height {
		truncated = truncated[len(truncated)-height:]
	}
	return strings.Join(truncated, "\n")
}

func (m Model) renderDiffContent(width, height int) string {
	if m.diffContent == "" {
		return m.styles.InactiveGlyph.Render("(no changes)")
	}
	lines := strings.Split(m.diffContent, "\n")
	// Apply scroll
	start := min(max(m.diffScroll, 0), len(lines)-1)
	end := min(start+height, len(lines))
	visible := lines[start:end]
	var truncated []string
	for _, line := range visible {
		truncated = append(truncated, ansi.Truncate(line, width, ""))
	}
	return strings.Join(truncated, "\n")
}

func (m Model) renderHint(width int) string {
	if m.pendingDelete != "" {
		parts := strings.SplitN(m.pendingDelete, "/", 2)
		name := m.pendingDelete
		if len(parts) == 2 {
			name = parts[1]
		}
		hint := "press d again to delete " + name + " · any other key cancels"
		return m.styles.HintBar.Render(ansi.Truncate(hint, width, ""))
	}
	hint := "↵ attach · j/k move · h/l fold · z all · ⇥ tab · d del · r repair · / filter · ? help · q quit"
	if m.filtering {
		hint = "filter: " + m.filter + "█  · esc: clear · enter: keep"
	}
	return m.styles.HintBar.Render(ansi.Truncate(hint, width, ""))
}

func (m Model) renderStatus(width int) string {
	if m.status == "" {
		return ""
	}
	return m.styles.StatusMsg.Render(ansi.Truncate(m.status, width, ""))
}

func (m Model) renderHelpOverlay() string {
	type entry struct{ key, desc string }
	entries := []entry{
		{"enter", "attach (auto-repairs missing window)"},
		{"j / k", "move cursor down / up"},
		{"h / l", "collapse / expand repo"},
		{"z", "toggle collapse all repos"},
		{"tab", "switch Agent / Diff pane"},
		{"d d", "delete worktree (confirm twice)"},
		{"r", "repair (recreate window + relaunch AI)"},
		{"/", "filter  esc:clear  enter:keep"},
		{"ctrl+d / ctrl+u", "scroll Diff pane down / up"},
		{"?", "toggle this help"},
		{"q / ctrl+c", "quit"},
	}

	const keyColW = 16
	const descColW = 36
	// inner: 1 space + key + 1 space + │ + 1 space + desc + 1 space = keyColW+descColW+4
	boxInnerW := keyColW + descColW + 4
	boxW := boxInnerW + 2
	if boxW > m.width-2 {
		boxW = m.width - 2
		boxInnerW = boxW - 2
	}

	b := m.styles.HelpBorder.Render
	sep := strings.Repeat("─", boxInnerW)

	title := "Keybindings"
	titleStyled := m.styles.RepoHeader.Render(title)
	lpad := max((boxInnerW-len(title))/2, 0)
	rpad := max(boxInnerW-lpad-len(title), 0)

	var sb strings.Builder
	sb.WriteString(b("┌"+sep+"┐") + "\n")
	sb.WriteString(
		b(
			"│",
		) + strings.Repeat(
			" ",
			lpad,
		) + titleStyled + strings.Repeat(
			" ",
			rpad,
		) + b(
			"│",
		) + "\n",
	)
	sb.WriteString(b("├"+sep+"┤") + "\n")

	for _, e := range entries {
		keyStyled := m.styles.HelpKey.Render(e.key)
		keyPad := strings.Repeat(" ", max(0, keyColW-ansi.StringWidth(e.key)))
		desc := ansi.Truncate(e.desc, descColW, "")
		descPad := strings.Repeat(" ", max(0, descColW-ansi.StringWidth(desc)))
		sb.WriteString(
			b("│") + " " + keyStyled + keyPad + b(" │ ") + desc + descPad + " " + b("│") + "\n",
		)
	}

	sb.WriteString(b("└"+sep+"┘") + "\n")
	sb.WriteString(m.styles.HintBar.Render("press any key to close"))

	boxLines := strings.Split(strings.TrimRight(sb.String(), "\n"), "\n")
	topPad := max((m.height-len(boxLines))/2, 0)
	leftPad := max((m.width-boxW)/2, 0)
	indent := strings.Repeat(" ", leftPad)

	var out []string
	for range topPad {
		out = append(out, "")
	}
	for _, l := range boxLines {
		out = append(out, indent+l)
	}
	return strings.Join(out, "\n")
}

// padLines ensures a slice of lines has exactly n entries, each padded/truncated to w visible chars.
// Uses ansi.StringWidth so ANSI escape codes are not counted as visible characters.
func padLines(lines []string, n, w int) []string {
	blank := strings.Repeat(" ", w)
	result := make([]string, n)
	for i := range n {
		if i < len(lines) {
			t := ansi.Truncate(lines[i], w, "")
			vis := ansi.StringWidth(t)
			if vis < w {
				t += strings.Repeat(" ", w-vis)
			}
			result[i] = t
		} else {
			result[i] = blank
		}
	}
	return result
}
