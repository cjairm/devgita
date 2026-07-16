// Package tuiworktree provides the Bubble Tea TUI dashboard for dg wt ui.
package tuiworktree

import (
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/tooling/task"
	"github.com/cjairm/devgita/internal/tooling/worktree"
	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

const (
	minLeftPaneWidth     = 20
	defaultLeftPaneWidth = 35
	maxLeftPaneWidthPct  = 0.60
	dividerWidth         = 1
	refreshInterval      = 3 * time.Second
	maxDiffBytes         = 64 * 1024
)

// --- Messages ---

type (
	statusesMsg []worktree.WorktreeStatus
	diffMsg     struct {
		content               string
		files, added, removed int
		fileLines             []int  // line indexes of per-file headers, for [ / ] jumps
		base                  string // comparison base label, e.g. "main @3e90667"
		branch                string // worktree branch the diff belongs to
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

	diffContent   string
	diffFiles     int
	diffAdded     int
	diffRemoved   int
	diffScroll    int
	diffFocused   bool
	diffFileLines []int
	diffBase      string
	diffBranch    string

	filter tuicomponents.FilterField

	status string

	width  int
	height int

	palette *tuicomponents.Palette

	leftPaneWidth int

	dragging   bool
	dragStartX int

	pendingDelete        string // "repo/name" or ""
	pendingSessionDelete string // "repo/name" or ""
	showHelp             bool

	createMode         createMode
	repoPicker         *tuicomponents.FuzzyPicker
	createRepo         string // resolved repo path chosen in repo-pick mode
	createName         string // in-progress name text in name-input mode
	pendingHookWarning bool   // armed by a first enter when CheckHookCompatibility found warnings; a second enter confirms, any other key (or edited name) de-arms it
	creating           bool   // true from the moment the create tea.Cmd is dispatched until its result (createdMsg/createFailedMsg) is processed; the ONLY thing that actually enforces "one create at a time" (see createFn's WarnFn-swap comment below) — handleNewWorktree checks this and ignores n while it's true

	// Injected I/O seams (overridable in tests)
	diffFn                   func(path string) (task.BranchDiffResult, error)
	attachFn                 func(session, window string) error
	removeFn                 func(repo, name string, force bool) error
	removeSessionFn          func(repo, name string) error
	repairFn                 func(repo, name string, coder worktree.AICoder) error
	windowSessionFn          func(window string) (string, bool)
	repoCandidatesFn         func(cursorRepoSlug string) ([]string, error)
	validateRepoPathFn       func(path string) (string, error)
	checkHookCompatibilityFn func(repoPath string) []string
	createFn                 func(repoPath, name string) (warning string, err error)
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
		palette:       tuicomponents.NewPalette(),
		leftPaneWidth: defaultLeftPaneWidth,
	}
	m.diffFn = func(path string) (task.BranchDiffResult, error) {
		return task.BranchDiffAt(gitApp, path)
	}
	m.attachFn = func(session, window string) error {
		return tmuxApp.SwitchToWindow(session, window)
	}
	m.removeFn = func(repo, name string, force bool) error {
		return mgr.RemoveInRepo(repo, name, force)
	}
	m.removeSessionFn = func(repo, name string) error {
		return mgr.RemoveWithSessionInRepo(repo, name)
	}
	m.repairFn = func(repo, name string, coder worktree.AICoder) error {
		return mgr.RepairInRepo(repo, name, coder)
	}
	m.windowSessionFn = tmuxApp.WindowSession
	m.repoCandidatesFn = mgr.RepoCandidates
	m.validateRepoPathFn = mgr.ValidateRepoPath
	// CheckHookCompatibility only stats/reads hook files (and reads
	// core.hooksPath via a read-only git config --get) — no prints, no
	// prompts, no writes — so it's safe to call directly from the TUI, unlike
	// worktree.go's own force=false path which raw-prints and blocks on
	// os.Stdin. The model calls this itself before create so the user still
	// gets the warning, just through a TUI-safe confirm (see
	// handleNameInputKey), instead of losing it to a hardcoded force=true.
	m.checkHookCompatibilityFn = gitApp.CheckHookCompatibility
	m.createFn = func(repoPath, name string) (string, error) {
		alias := worktree.ResolveAIAlias("", gc)
		coder, err := worktree.ResolveAICoder(alias)
		if err != nil {
			return "", err
		}
		// mgr.WarnFn fires synchronously from inside CreateAt below (e.g. the
		// recent-repos store failed to record this create). Swapping it to a
		// local closure and restoring it via defer right after CreateAt
		// returns is safe only because this TUI never runs two creates
		// concurrently — and that's actually true, not just assumed: this
		// tea.Cmd closure only ever runs while m.creating is true, and
		// handleNewWorktree (the only way to start another create) refuses
		// to open the picker while m.creating is true, so a second createFn
		// call can never be in flight to race this swap/restore.
		var warning string
		original := mgr.WarnFn
		mgr.WarnFn = func(msg string) { warning = msg }
		defer func() { mgr.WarnFn = original }()
		// force=true is safe here specifically because the model already ran
		// its own equivalent hook-compatibility check (checkHookCompatibilityFn)
		// and, when it found warnings, its own equivalent TUI-safe confirm
		// (handleNameInputKey's pendingHookWarning arm/second-enter) before
		// ever reaching this closure — so this isn't skipping the check,
		// worktree.go's raw fmt.Println/stdin-read version of it is just
		// redundant by this point and would otherwise corrupt or hang the
		// running bubbletea alt-screen display.
		if err := mgr.CreateAt(repoPath, name, coder, true); err != nil {
			return "", err
		}
		return warning, nil
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

func (m Model) computeDiffCmd(s worktree.WorktreeStatus) tea.Cmd {
	df := m.diffFn
	p := m.palette
	path := s.Path
	branch := s.Branch
	if branch == "" {
		branch = s.Name
	}
	return func() tea.Msg {
		res, err := df(path)
		content := res.Content
		if err != nil {
			content = "(diff unavailable: " + err.Error() + ")"
		} else {
			content = rewriteFileHeaders(content, res.FileStats, p)
		}
		if len(content) > maxDiffBytes {
			content = content[:maxDiffBytes] + "\n... (truncated)"
		}
		base := res.BaseBranch
		if base != "" && res.BaseSHA != "" {
			base += " @" + res.BaseSHA
		}
		return diffMsg{
			content:   content,
			files:     res.Files,
			added:     res.Added,
			removed:   res.Removed,
			fileLines: diffFileHeaderLines(content),
			base:      base,
			branch:    branch,
		}
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
	m.rows = buildRows(m.statuses, m.collapsed, m.filter.Text)
	// Keep cursor on a valid worktree row
	m.cursor = tuicomponents.ClampCursor(worktreeIndices(m.rows), m.cursor)
}

// navigableIndices returns row indices that j/k visit: all worktree rows plus
// collapsed repo header rows (so the user can reach a collapsed header and press l).
func (m *Model) navigableIndices() []int {
	var out []int
	for i, r := range m.rows {
		if r.kind == rowWorktree || (r.kind == rowRepo && m.collapsed[r.repo]) {
			out = append(out, i)
		}
	}
	return out
}

func (m *Model) moveCursor(delta int) {
	m.cursor = tuicomponents.MoveCursor(m.navigableIndices(), m.cursor, delta)
}

func (m Model) safeMaxLeft() int {
	return max(int(float64(m.width)*maxLeftPaneWidthPct), minLeftPaneWidth)
}

func (m Model) rightPaneWidth() int {
	return max(m.width-m.leftPaneWidth-dividerWidth, 0)
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
			return m, m.computeDiffCmd(sel)
		}
		return m, nil

	case diffMsg:
		m.diffContent = msg.content
		m.diffFiles = msg.files
		m.diffAdded = msg.added
		m.diffRemoved = msg.removed
		m.diffFileLines = msg.fileLines
		m.diffBase = msg.base
		m.diffBranch = msg.branch
		return m, nil

	case tickMsg:
		// Periodic refresh: reload statuses, which re-derives the selected
		// worktree's diff via the statusesMsg handler.
		return m, tea.Batch(m.loadCmd(), m.tickCmd())

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case createdMsg:
		return m.handleCreateSuccess(msg.repoPath, msg.name, msg.warning)

	case createFailedMsg:
		m.creating = false
		m.status = "create failed: " + msg.err.Error()
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

	// Repo-pick and name-input float over the dashboard and intercept every
	// key while active, the same way showHelp does above.
	if m.createMode == createRepoPick {
		return m.handleRepoPickKey(key)
	}
	if m.createMode == createNameInput {
		return m.handleNameInputKey(key)
	}

	if m.filter.Active {
		if m.filter.HandleKey(key) {
			m.rebuildRows()
		}
		return m, nil
	}

	if m.diffFocused {
		return m.handleDiffKey(key)
	}

	// Clear pending deletes on any key that doesn't continue the confirmation
	if key != "d" && m.pendingDelete != "" {
		m.pendingDelete = ""
	}
	if key != "D" && m.pendingSessionDelete != "" {
		m.pendingSessionDelete = ""
	}

	switch key {
	case "?":
		m.showHelp = true
		return m, nil

	case "q", "ctrl+c":
		return m, tea.Quit

	case "j":
		m.diffScroll = 0
		m.moveCursor(1)
		if sel, ok := m.selectedStatus(); ok {
			return m, m.computeDiffCmd(sel)
		}
		return m, nil

	case "k":
		m.diffScroll = 0
		m.moveCursor(-1)
		if sel, ok := m.selectedStatus(); ok {
			return m, m.computeDiffCmd(sel)
		}
		return m, nil

	case "h":
		var collapseRepo string
		if sel, ok := m.selectedStatus(); ok {
			collapseRepo = sel.Repo
		} else if m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].kind == rowRepo {
			collapseRepo = m.rows[m.cursor].repo
		}
		if collapseRepo != "" {
			m.collapsed[collapseRepo] = true
			m.rebuildRows()
			// Land cursor on the just-collapsed repo header so l can re-expand it.
			for i, r := range m.rows {
				if r.kind == rowRepo && r.repo == collapseRepo {
					m.cursor = i
					break
				}
			}
		}
		return m, nil

	case "l":
		var expandRepo string
		wasCollapsed := false
		if sel, ok := m.selectedStatus(); ok {
			expandRepo = sel.Repo
			wasCollapsed = m.collapsed[expandRepo]
		} else if m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].kind == rowRepo {
			expandRepo = m.rows[m.cursor].repo
			wasCollapsed = m.collapsed[expandRepo]
		}
		if expandRepo != "" {
			m.collapsed[expandRepo] = false
			m.rebuildRows()
			if wasCollapsed {
				// Move cursor to first worktree of the just-expanded repo.
				for i, r := range m.rows {
					if r.kind == rowWorktree && r.repo == expandRepo {
						m.cursor = i
						break
					}
				}
			}
		}
		return m, nil

	case "z":
		m.allCollapsed = !m.allCollapsed
		for _, s := range m.statuses {
			m.collapsed[s.Repo] = m.allCollapsed
		}
		m.rebuildRows()
		if m.allCollapsed && len(m.rows) > 0 {
			m.cursor = 0 // first visible row is a repo header when all collapsed
		}
		return m, nil

	case "/":
		m.filter.Active = true
		return m, nil

	case "space":
		if m.diffContent != "" {
			m.diffFocused = true
		}
		return m, nil

	case "enter":
		return m.handleAttach()

	case "d":
		return m.handleDelete()

	case "D":
		return m.handleSessionDelete()

	case "r":
		return m.handleRepair()

	case "n":
		return m.handleNewWorktree()

	case "ctrl+d":
		m.diffScroll = min(m.diffScroll+m.diffPageSize(), m.maxDiffScroll())
		return m, nil

	case "ctrl+u":
		m.diffScroll = max(m.diffScroll-m.diffPageSize(), 0)
		return m, nil
	}

	return m, nil
}

// handleDiffKey processes keys while the diff pane is focused: vim-style
// scrolling, [ / ] file jumps, and esc/space to return to the list.
func (m Model) handleDiffKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc", "space":
		m.diffFocused = false

	case "?":
		m.showHelp = true

	case "j", "down":
		m.diffScroll = min(m.diffScroll+1, m.maxDiffScroll())

	case "k", "up":
		m.diffScroll = max(m.diffScroll-1, 0)

	case "ctrl+d":
		m.diffScroll = min(m.diffScroll+m.diffPageSize(), m.maxDiffScroll())

	case "ctrl+u":
		m.diffScroll = max(m.diffScroll-m.diffPageSize(), 0)

	case "g":
		m.diffScroll = 0

	case "G":
		m.diffScroll = m.maxDiffScroll()

	case "]":
		for _, ln := range m.diffFileLines {
			if ln > m.diffScroll {
				m.diffScroll = min(ln, m.maxDiffScroll())
				break
			}
		}

	case "[":
		for i := len(m.diffFileLines) - 1; i >= 0; i-- {
			if m.diffFileLines[i] < m.diffScroll {
				m.diffScroll = m.diffFileLines[i]
				break
			}
		}
	}
	return m, nil
}

func (m Model) diffPageSize() int {
	return max(m.height-4, 1)
}

func (m Model) maxDiffScroll() int {
	return max(strings.Count(m.diffContent, "\n"), 0)
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

	return m, m.attachToWindowCmd(sel.Repo, sel.Name)
}

// attachToWindowCmd looks up repo/name's tmux window (auto-repairing it
// first if the window is missing) and attaches, quitting the TUI on
// success. Extracted from handleAttach so the create flow can reuse the
// identical retry/auto-repair logic for a worktree that was just created and
// isn't in m.statuses yet (selectedStatus can't provide it there), instead
// of duplicating this logic at a second call site.
func (m Model) attachToWindowCmd(repo, name string) tea.Cmd {
	window := worktree.GetWindowName(repo, name)
	attachFn := m.attachFn
	repairFn := m.repairFn
	windowSessionFn := m.windowSessionFn
	gc := m.gc

	return func() tea.Msg {
		session, ok := windowSessionFn(window)
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
		if err := repairFn(repo, name, coder); err != nil {
			return statusMsg("repair failed: " + err.Error())
		}
		session, ok = windowSessionFn(window)
		if !ok {
			return statusMsg("repair succeeded but window not found")
		}
		if err := attachFn(session, window); err != nil {
			return statusMsg("attach after repair failed: " + err.Error())
		}
		return tea.QuitMsg{}
	}
}

// confirmThenRemove implements the shared two-press delete confirmation.
// pending is the currently armed "repo/name" key (or ""); remove performs the
// removal on the second press. It returns the new pending value and, on
// confirmation, a command that runs the removal and refreshes the list.
func (m Model) confirmThenRemove(
	pending string,
	remove func(repo, name string) error,
) (string, tea.Cmd) {
	sel, ok := m.selectedStatus()
	if !ok {
		return pending, nil
	}

	key := sel.Repo + "/" + sel.Name

	// First press (or cursor moved to another row): arm
	if pending != key {
		return key, nil
	}

	// Second press: execute
	repo := sel.Repo
	name := sel.Name
	statuses := m.statuses
	return "", func() tea.Msg {
		if err := remove(repo, name); err != nil {
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

func (m Model) handleDelete() (tea.Model, tea.Cmd) {
	removeFn := m.removeFn
	pending, cmd := m.confirmThenRemove(m.pendingDelete, func(repo, name string) error {
		return removeFn(repo, name, true)
	})
	m.pendingDelete = pending
	return m, cmd
}

func (m Model) handleSessionDelete() (tea.Model, tea.Cmd) {
	pending, cmd := m.confirmThenRemove(m.pendingSessionDelete, m.removeSessionFn)
	m.pendingSessionDelete = pending
	return m, cmd
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

	background := m.renderDashboard()
	if m.showHelp {
		return tuicomponents.Overlay(background, m.renderHelpPopup(), m.width, m.height)
	}
	if m.createMode == createRepoPick {
		return tuicomponents.Overlay(background, m.renderRepoPickPopup(), m.width, m.height)
	}
	if m.createMode == createNameInput {
		return tuicomponents.Overlay(background, m.renderNameInputPopup(), m.width, m.height)
	}
	return background
}

// renderDashboard renders the normal (non-help) dashboard: narrow-terminal
// fallback or the left+divider+right layout, plus hint and status lines. It
// always runs, even while the help popup is shown, so renderContent has a
// live background to composite the popup over instead of blanking the screen.
func (m Model) renderDashboard() string {
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
	// isLastChild reports whether row i is the last worktree under its repo header.
	isLastChild := func(i int) bool {
		repo := m.rows[i].status.Repo
		for j := i + 1; j < len(m.rows); j++ {
			if m.rows[j].kind == rowRepo {
				return true
			}
			if m.rows[j].kind == rowWorktree && m.rows[j].status.Repo == repo {
				return false
			}
		}
		return true
	}

	const branchChar = "∕" // U+2215 DIVISION SLASH — branch glyph (1 display col)

	var sb strings.Builder
	for i, r := range m.rows {
		var line string
		if r.kind == rowRepo {
			collapse := "▼"
			if m.collapsed[r.repo] {
				collapse = "▶"
			}
			text := collapse + " " + r.repo
			pad := strings.Repeat(" ", max(0, width-ansi.StringWidth(text)))
			if i == m.cursor {
				// Cursor landed here after h — show repo header with selection highlight.
				line = m.palette.Selected.Render(text + pad)
			} else {
				line = m.palette.RepoHeader.Render(text) + pad
			}
		} else {
			state := tuicomponents.SessionStateFromWorktree(r.status, false, 0)
			// Tree connector: "└ " for last child, "  " otherwise (both 2 display cols).
			connectorRaw := "  "
			connectorStyled := "  "
			if isLastChild(i) {
				connectorRaw = "└ "
				connectorStyled = m.palette.Divider.Render("└") + " "
			}
			// prefix = connector(2) + dot(1) + branchChar(1) + space(1) = 5 display cols
			name := ansi.Truncate(r.status.Name, max(0, width-5), "")
			pendingKey := r.status.Repo + "/" + r.status.Name
			padding := strings.Repeat(" ", max(0, width-5-ansi.StringWidth(name)))

			if i == m.cursor {
				g := m.palette.StatusGlyph(state)
				plainText := connectorRaw + g + branchChar + " " + name
				if m.pendingDelete == pendingKey || m.pendingSessionDelete == pendingKey {
					line = m.palette.Armed.Render(plainText + padding)
				} else {
					line = m.palette.Selected.Render(plainText + padding)
				}
			} else {
				line = connectorStyled + m.palette.StatusDot(
					state,
				) + m.palette.BranchLabel() + " " + name + padding
			}
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func (m Model) renderDivider(height int) string {
	style := m.palette.Divider
	if m.diffFocused {
		// Brighter divider signals the diff pane holds keyboard focus.
		style = m.palette.HintKey
	}
	divChar := style.Render("│")
	lines := make([]string, height)
	for i := range lines {
		lines[i] = divChar
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderRight(width int) string {
	header := m.palette.DiffStatLine(m.diffFiles, m.diffAdded, m.diffRemoved)
	// GitHub-style "base ← compare" label, shown once for the whole diff.
	if m.diffBase != "" && m.diffBranch != "" {
		header = m.palette.RepoHeader.Render(m.diffBase) +
			m.palette.Divider.Render(" ← ") +
			m.palette.DiffFileHeader.Render(m.diffBranch) +
			"  " + header
	}
	contentHeight := max(m.height-4, 0) // height minus hint, status, header, blank line
	content := m.renderDiffContent(width, contentHeight)
	return ansi.Truncate(header, width, "") + "\n" + content
}

func (m Model) renderDiffContent(width, height int) string {
	if m.diffContent == "" {
		return m.palette.Inactive.Render("(loading...)")
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

// armedDeleteHint renders the confirmation hint for a two-press delete.
// pending is the armed "repo/name" key, key the keypress to repeat, and
// suffix an optional description of extra effects beyond the worktree delete.
func (m Model) armedDeleteHint(pending, key, suffix string, width int) string {
	parts := strings.SplitN(pending, "/", 2)
	name := pending
	if len(parts) == 2 {
		name = parts[1]
	}
	hint := "press " + key + " again to delete " + name + suffix + " · any other key cancels"
	return m.palette.HintDesc.Render(ansi.Truncate(hint, width, ""))
}

func (m Model) renderHint(width int) string {
	if m.pendingSessionDelete != "" {
		return m.armedDeleteHint(m.pendingSessionDelete, "D", " and kill its session", width)
	}
	if m.pendingDelete != "" {
		return m.armedDeleteHint(m.pendingDelete, "d", "", width)
	}
	if m.createMode == createRepoPick {
		hints := []tuicomponents.KeyHint{
			{Key: "esc", Desc: "cancel"},
			{Key: "enter", Desc: "select"},
			{Key: "↑/↓", Desc: "move"},
		}
		return m.palette.HintBar(hints, width)
	}
	if m.createMode == createNameInput {
		hints := []tuicomponents.KeyHint{
			{Key: "esc", Desc: "cancel"},
			{Key: "enter", Desc: "create"},
		}
		return m.palette.HintBar(hints, width)
	}
	if m.filter.Active {
		return m.palette.FilterHint(m.filter, width)
	}
	if m.diffFocused {
		hints := []tuicomponents.KeyHint{
			{Key: "esc", Desc: "back"},
			{Key: "j/k", Desc: "scroll"},
			{Key: "^d/^u", Desc: "page"},
			{Key: "[/]", Desc: "file"},
			{Key: "g/G", Desc: "top/end"},
			{Key: "?", Desc: "help"},
			{Key: "q", Desc: "quit"},
		}
		return m.palette.HintBar(hints, width)
	}
	hints := []tuicomponents.KeyHint{
		{Key: "↵", Desc: "attach"},
		{Key: "n", Desc: "new"},
		{Key: "spc", Desc: "diff"},
		{Key: "j/k", Desc: "move"},
		{Key: "h/l", Desc: "fold"},
		{Key: "z", Desc: "all"},
		{Key: "d", Desc: "del"},
		{Key: "D", Desc: "del+sess"},
		{Key: "r", Desc: "repair"},
		{Key: "/", Desc: "filter"},
		{Key: "?", Desc: "help"},
		{Key: "q", Desc: "quit"},
	}
	return m.palette.HintBar(hints, width)
}

func (m Model) renderStatus(width int) string {
	if m.status == "" {
		return ""
	}
	return m.palette.StatusMsg.Render(ansi.Truncate(m.status, width, ""))
}

// renderHelpPopup builds the raw (uncentered) help popup content; the
// caller composites it over the dashboard background via Overlay.
func (m Model) renderHelpPopup() string {
	entries := []tuicomponents.WhichKeyEntry{
		{Key: "enter", Desc: "attach (auto-repairs missing window)"},
		{Key: "n", Desc: "create a new worktree (repo picker → name prompt)"},
		{Key: "j / k", Desc: "move cursor down / up"},
		{Key: "h / l", Desc: "collapse / expand repo"},
		{Key: "z", Desc: "toggle collapse all repos"},
		{Key: "d d", Desc: "delete worktree (confirm twice)"},
		{Key: "D D", Desc: "delete worktree + kill its session"},
		{Key: "r", Desc: "repair (recreate window + relaunch AI)"},
		{Key: "/", Desc: "filter  esc:clear  enter:keep"},
		{Key: "space", Desc: "focus diff pane (esc returns to the list)"},
		{Key: "ctrl+d / ctrl+u", Desc: "scroll diff down / up"},
		{Key: "[ / ]", Desc: "previous / next file (diff focused)"},
		{Key: "g / G", Desc: "diff top / bottom (diff focused)"},
		{Key: "?", Desc: "toggle this help"},
		{Key: "q / ctrl+c", Desc: "quit"},
	}
	return m.palette.HelpPopup("Keybindings", entries, m.width)
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
