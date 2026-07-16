// The n → repo-pick → name-input → create flow: picking a repo (from a
// FuzzyPicker list or a free-typed path), naming the new worktree, running
// the hook-compatibility confirm, and following a successful create with the
// same attach-and-quit path as pressing enter on an existing worktree row.

package tuiworktree

import (
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

// createMode tracks progress through the n → repo-pick → name-input → create
// flow, layered on top of the dashboard's existing mode flags (showHelp,
// filter.Active, diffFocused) the same way those are: a small explicit field
// checked early in handleKey, not a deep state machine.
type createMode int

const (
	createNone createMode = iota
	createRepoPick
	createNameInput
)

// createdMsg reports a successful createFn call so Update can trigger the
// attach-and-quit (or refresh-on-cancel) follow-up outside the create cmd
// itself, matching how every other async result in this model flows back
// through a typed message rather than a callback. warning carries a non-fatal
// WarnFn message captured during the create (e.g. recent-repos record
// failure), if any, so it can be surfaced without blocking attach/quit.
type createdMsg struct {
	repoPath string
	name     string
	warning  string
}

// createFailedMsg reports a failed createFn call as its own typed message
// (rather than the generic statusMsg every other failure in this model
// uses) so Update can clear m.creating specifically when a create's result
// is what's being processed, instead of clearing it on every unrelated
// statusMsg (delete/session-delete/attach failures use statusMsg too, and
// none of those should touch the create-in-flight guard).
type createFailedMsg struct {
	err error
}

// handleNewWorktree opens the repo picker for the n keybinding, offering the
// cursor row's repo first. RepoCandidates already orders that repo first
// when present, matching FuzzyPicker's initial cursor at index 0, so no
// extra pre-selection logic is needed here.
//
// Guarded by m.creating: clearCreateState runs synchronously in
// handleNameInputKey's enter branch, before the async createFn tea.Cmd it
// returns has run, so without this guard a user could press n again while a
// create is still in flight and start a second concurrent createFn call —
// racing the mgr.WarnFn swap/restore in createFn (see newModel). m.creating
// stays true for exactly that window, so a second n is a no-op instead.
func (m Model) handleNewWorktree() (tea.Model, tea.Cmd) {
	if m.creating {
		return m, nil
	}

	var cursorRepoSlug string
	if sel, ok := m.selectedStatus(); ok {
		cursorRepoSlug = sel.Repo
	} else if m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].kind == rowRepo {
		cursorRepoSlug = m.rows[m.cursor].repo
	}

	candidates, err := m.repoCandidatesFn(cursorRepoSlug)
	if err != nil {
		m.status = "failed to list repos: " + err.Error()
		return m, nil
	}

	items := make([]tuicomponents.PaletteItem, len(candidates))
	for i, c := range candidates {
		items[i] = tuicomponents.PaletteItem{Command: c}
	}

	m.repoPicker = tuicomponents.NewFuzzyPicker("New worktree — pick a repo", items)
	m.createMode = createRepoPick
	return m, nil
}

// resolveAndValidateRepoPath runs candidate through validateRepoPathFn and,
// on success, transitions into name-input mode with the resolved root
// recorded as createRepo. On failure it reports the error via m.status and
// leaves createMode untouched, so the caller stays in the picker.
//
// Both handleRepoPickKey branches (a list selection and a free-typed query)
// call this: a FuzzyPicker candidate isn't guaranteed to be a git repo
// root — RepoCandidates' zoxide source only canonicalizes a path
// (config.CanonicalRepoPath), it never confirms the path is actually a repo
// root, so a zoxide-tracked subdirectory could otherwise slip through
// unvalidated the way only free-typed paths used to be checked. Validating
// both the same way also means a resolved root that differs from the raw
// candidate (e.g. a subdirectory resolving to its repo's root) is what
// downstream code (checkHookCompatibilityFn, createFn, and
// handleCreateSuccess's repo-slug derivation) actually sees.
func (m *Model) resolveAndValidateRepoPath(candidate string) {
	root, err := m.validateRepoPathFn(candidate)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.createRepo = root
	m.createMode = createNameInput
	m.createName = ""
}

// handleRepoPickKey delegates to the FuzzyPicker for everything except the
// free-typed-path case: FuzzyPicker only knows about the candidate items it
// was given, so a query that matches none of them makes enter report
// FuzzyPickerNone (by design) instead of selecting anything. The model
// validates that query as a path itself in that case.
func (m Model) handleRepoPickKey(key string) (tea.Model, tea.Cmd) {
	result := m.repoPicker.HandleKey(key)
	switch result.Action {
	case tuicomponents.FuzzyPickerCancelled:
		m.clearCreateState()

	case tuicomponents.FuzzyPickerSelected:
		m.resolveAndValidateRepoPath(result.Item.Command)

	case tuicomponents.FuzzyPickerNone:
		if key != "enter" {
			break
		}
		query := m.repoPicker.Query()
		if query == "" {
			break
		}
		m.resolveAndValidateRepoPath(query)
	}
	return m, nil
}

// handleNameInputKey drives the floating single-line name prompt: printable
// characters and backspace edit createName, esc cancels the whole create
// flow, and enter with a non-empty name kicks off createFn. Auto-naming on a
// blank name is out of scope for this cycle, so enter with no text is a
// no-op rather than falling back to a generated name.
//
// Enter also runs the hook-compatibility check (checkHookCompatibilityFn)
// before create: if it finds warnings, the same two-press arm/confirm
// pattern as delete (confirmThenRemove/pendingDelete) is used instead of
// worktree.go's own force=false prompt, which raw-prints and blocks on
// os.Stdin and would corrupt or hang the running bubbletea alt-screen
// program. The first enter arms pendingHookWarning and shows the warning as
// a status; a second enter proceeds. Editing the name (or any other key)
// de-arms it, same as confirmThenRemove clearing pendingDelete.
func (m Model) handleNameInputKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.clearCreateState()
		return m, nil

	case "enter":
		if m.createName == "" {
			return m, nil
		}
		if !m.pendingHookWarning {
			if warnings := m.checkHookCompatibilityFn(m.createRepo); len(warnings) > 0 {
				m.pendingHookWarning = true
				m.status = "hook warning: " + strings.Join(warnings, "; ") +
					" — press enter again to create anyway"
				return m, nil
			}
		}
		repoPath := m.createRepo
		name := m.createName
		createFn := m.createFn
		m.clearCreateState()
		// Armed here, before the tea.Cmd below ever runs: Bubble Tea runs
		// every returned tea.Cmd in its own goroutine, so without setting
		// this now (synchronously, in the same Update call that clears
		// createMode back to createNone) a user could press n again and
		// start a second concurrent create before this one resolves.
		// handleNewWorktree checks m.creating and refuses to open the
		// picker while it's true; cleared again once createdMsg or
		// createFailedMsg is processed (handleCreateSuccess / Update).
		m.creating = true
		return m, func() tea.Msg {
			warning, err := createFn(repoPath, name)
			if err != nil {
				return createFailedMsg{err: err}
			}
			return createdMsg{repoPath: repoPath, name: name, warning: warning}
		}

	case "backspace":
		m.pendingHookWarning = false
		if len(m.createName) > 0 {
			m.createName = m.createName[:len(m.createName)-1]
		}

	default:
		m.pendingHookWarning = false
		if utf8.RuneCountInString(key) == 1 && key >= " " {
			m.createName += key
		}
	}
	return m, nil
}

// clearCreateState resets the n → repo-pick → name-input flow back to
// normal mode, used on cancellation and once a create has been kicked off.
// It does not touch m.creating: that field's lifecycle is independent,
// spanning from just before the create tea.Cmd is dispatched until its
// result is processed, which outlives this reset (clearCreateState runs
// synchronously before the async createFn call even starts).
func (m *Model) clearCreateState() {
	m.createMode = createNone
	m.repoPicker = nil
	m.createRepo = ""
	m.createName = ""
	m.pendingHookWarning = false
}

// handleCreateSuccess follows a successful createFn with the same
// attach-and-quit path as pressing enter on an existing worktree row
// (attachToWindowCmd), batched with a statuses refresh so the just-created
// worktree shows up in the dashboard list even when attach doesn't end in
// quitting (e.g. the window couldn't be found and repair also failed). A
// non-empty warning (captured from mgr.WarnFn by createFn) is surfaced as
// status prefixed with "created, but:" so the user still knows the create
// itself succeeded, without blocking the attach/quit flow that follows.
func (m Model) handleCreateSuccess(repoPath, name, warning string) (tea.Model, tea.Cmd) {
	m.creating = false
	if warning != "" {
		m.status = "created, but: " + warning
	} else if os.Getenv("TMUX") == "" {
		m.status = "worktree created: " + name
	}
	if os.Getenv("TMUX") == "" {
		return m, m.loadCmd()
	}
	repoSlug := filepath.Base(repoPath)
	return m, tea.Batch(m.attachToWindowCmd(repoSlug, name), m.loadCmd())
}

// renderRepoPickPopup builds the raw (uncentered) repo-picker popup content;
// the caller composites it over the dashboard background via Overlay, same
// as the help popup.
func (m Model) renderRepoPickPopup() string {
	maxW := min(m.width-2, 64)
	return m.repoPicker.View(maxW)
}

// renderNameInputPopup builds the raw (uncentered) name-prompt popup
// content; composited over the dashboard background via Overlay.
func (m Model) renderNameInputPopup() string {
	maxW := min(m.width-2, 64)
	lines := []string{"> " + m.createName + "█"}
	return m.palette.BorderedPane("New worktree — name", maxW, lines)
}
