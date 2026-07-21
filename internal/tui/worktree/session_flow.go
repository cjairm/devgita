// The s → name-prompt → CreateSession flow for standalone tmux sessions.
// Deliberately separate from create_flow.go's n/N worktree flow: a plain
// session has no repo pick, no hook-compatibility check, and no layout pick,
// so reusing createMode/createInput here would conflate two flows whose
// invariants (createRepo, dispatchCreate's layout handling) don't apply to a
// session at all.

package tuiworktree

import (
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/cjairm/devgita/pkg/paths"
)

// sessionCreatedMsg reports a successful createSessionFn call made outside
// tmux, where there is no attached client to switch - so Update sets the
// "session created:" confirmation and refreshes the session rows instead.
// Inside tmux, a successful create instead switches-and-quits directly from
// dispatchSessionCreate's tea.Cmd, and this message is never produced.
type sessionCreatedMsg struct {
	name string
}

// sessionKilledMsg reports a successful killSessionFn call so Update can drop
// the killed session from m.sessions, rebuild rows, and set a "removed:"
// confirmation - mirroring deletedMsg's shape for worktree deletes.
type sessionKilledMsg struct {
	name string
}

// handleNewSession opens the floating name prompt for the s keybinding, from
// any row (a plain session isn't scoped to whatever row the cursor happens to
// be on). No re-entry guard is needed: while m.creatingSession is true,
// handleKey's early interception routes every key (including another "s") to
// handleSessionNameInputKey instead of back through this dispatch.
func (m Model) handleNewSession() (tea.Model, tea.Cmd) {
	m.creatingSession = true
	m.sessionNameInput.Reset()
	return m, nil
}

// handleSessionNameInputKey drives the floating single-line session-name
// prompt: esc cancels, enter with a non-empty name kicks off
// dispatchSessionCreate, and every other key is a name edit delegated to the
// shared TextInput - the same shape as handleNameInputKey, minus the hook-
// compatibility confirm and layout-pick branching that don't apply here.
func (m Model) handleSessionNameInputKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.clearSessionState()
		return m, nil

	case "enter":
		if m.sessionNameInput.Value == "" {
			return m, nil
		}
		return m.dispatchSessionCreate()

	default:
		m.sessionNameInput.HandleKey(key)
	}
	return m, nil
}

// handleSessionNameInputPaste inserts pasted text into the session-name field
// in one shot, the paste counterpart to handleSessionNameInputKey.
func (m Model) handleSessionNameInputPaste(text string) (tea.Model, tea.Cmd) {
	m.sessionNameInput.InsertText(text)
	return m, nil
}

// clearSessionState resets the s → name-prompt flow back to normal mode, used
// on cancellation and once a create has been dispatched.
func (m *Model) clearSessionState() {
	m.creatingSession = false
	m.sessionNameInput.Reset()
}

// dispatchSessionCreate captures the entered name, closes the prompt, and
// kicks off the async createSessionFn call. Workdir is always the user's home
// directory (paths.Paths.Home.Root): a plain session is deliberately not tied
// to any repo, so a selected-row-derived path would contradict the concept,
// and the TUI's own cwd is just wherever dg ws happened to launch from.
//
// Duplicate session names are not pre-checked here (Q2 in the cycle plan):
// tmux's own "new-session -s <name>" already fails when the name exists, and
// a prompt-layer HasSession check would both duplicate that enforcement and
// add a TOCTOU gap between the check and the create. The failure just
// surfaces via statusMsg, same as every other action failure in this file.
func (m Model) dispatchSessionCreate() (tea.Model, tea.Cmd) {
	name := m.sessionNameInput.Value
	m.clearSessionState()
	m.status = actionStatus("creating session", name)

	createFn := m.createSessionFn
	switchFn := m.switchToSessionFn
	workdir := paths.Paths.Home.Root
	insideTmux := os.Getenv("TMUX") != ""

	return m, func() tea.Msg {
		if err := createFn(name, workdir); err != nil {
			return statusMsg("create session failed: " + err.Error())
		}
		if !insideTmux {
			// tmux new-session -d works without an attached client, so
			// create-outside-tmux is allowed through - only the switch is
			// skipped, mirroring the worktree create-outside-tmux path
			// (handleCreateSuccess).
			return sessionCreatedMsg{name: name}
		}
		if err := switchFn(name); err != nil {
			return statusMsg("switch failed: " + err.Error())
		}
		return tea.QuitMsg{}
	}
}

// handleSwitchToSession is enter's rowSession counterpart to handleAttach:
// switches the attached client to the selected session and quits. Guarded by
// the same $TMUX check and message as handleAttach, since moving the client
// requires one to already be attached.
func (m Model) handleSwitchToSession() (tea.Model, tea.Cmd) {
	sel, ok := m.selectedSession()
	if !ok {
		return m, nil
	}

	if os.Getenv("TMUX") == "" {
		m.status = notInsideTmuxStatus
		return m, nil
	}

	switchFn := m.switchToSessionFn
	name := sel.Name
	return m, func() tea.Msg {
		if err := switchFn(name); err != nil {
			return statusMsg("switch failed: " + err.Error())
		}
		return tea.QuitMsg{}
	}
}

// handleKillSession is d's rowSession counterpart to handleDelete: a
// two-press kill confirmation, armed/confirmed the same way
// confirmThenRemove is (arm on first press, clear on any other key, confirm
// on second press - see handleKey's pendingKillSession clearing block), but
// keyed by session name alone since sessions have no repo.
func (m Model) handleKillSession() (tea.Model, tea.Cmd) {
	sel, ok := m.selectedSession()
	if !ok {
		return m, nil
	}
	name := sel.Name

	// First press (or cursor moved to another session): arm
	if m.pendingKillSession != name {
		m.pendingKillSession = name
		return m, nil
	}

	// Second press: execute
	m.pendingKillSession = ""
	killFn := m.killSessionFn
	m.status = actionStatus("killing session", name)
	return m, func() tea.Msg {
		if err := killFn(name); err != nil {
			return statusMsg("kill session failed: " + err.Error())
		}
		return sessionKilledMsg{name: name}
	}
}

// renderSessionNameInputPopup builds the raw (uncentered) session-name-prompt
// popup content; composited over the dashboard background via Overlay,
// mirroring renderNameInputPopup's shape.
func (m Model) renderSessionNameInputPopup() string {
	maxW := min(m.width-2, 64)
	lines := []string{"> " + m.sessionNameInput.RenderPlain()}
	return m.palette.BorderedPane("New session — name", maxW, lines)
}
