package tuiworktree

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cjairm/devgita/internal/tooling/worktree"
)

// --- Create flow (n) ---

func TestNewWorktreePicksCursorRepoFirst(t *testing.T) {
	m := makeTestModel(testStatuses()) // cursor starts on repo-a/feature-a
	var gotSlug string
	m.repoCandidatesFn = func(cursorRepoSlug string) ([]string, error) {
		gotSlug = cursorRepoSlug
		return []string{"/repos/repo-a", "/repos/other"}, nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)

	if gotSlug != "repo-a" {
		t.Errorf("expected repoCandidatesFn called with cursor repo 'repo-a', got %q", gotSlug)
	}
	if m3.createMode != createRepoPick {
		t.Fatalf("n should enter repo-pick mode, got mode=%d", m3.createMode)
	}
	if m3.repoPicker == nil {
		t.Fatal("expected a repo picker to be built")
	}
	sel, ok := m3.repoPicker.Selected()
	if !ok || sel.Command != "/repos/repo-a" {
		t.Errorf("expected first candidate pre-selected, got %+v ok=%v", sel, ok)
	}
}

func TestNewWorktreeIgnoredWhileCreating(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.creating = true
	candidatesCalled := false
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		candidatesCalled = true
		return []string{"/repos/alpha"}, nil
	}

	m2, cmd := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)

	if cmd != nil {
		t.Error("n while a create is in flight should return no command")
	}
	if candidatesCalled {
		t.Error("n while a create is in flight must not rebuild the repo picker")
	}
	if m3.createMode != createNone {
		t.Errorf("n while a create is in flight must not change createMode, got %d", m3.createMode)
	}
	if m3.repoPicker != nil {
		t.Error("n while a create is in flight must not build a repo picker")
	}
}

func TestRepoPickFilterSelectsMatch(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha", "/repos/beta"}, nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)

	for _, ch := range "beta" {
		m2, _ = m3.Update(tea.KeyPressMsg{Code: ch})
		m3 = m2.(Model)
	}
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)

	if m5.createMode != createNameInput {
		t.Fatalf("expected name-input mode after selecting a match, got mode=%d", m5.createMode)
	}
	if m5.createRepo != "/repos/beta" {
		t.Errorf("expected createRepo '/repos/beta', got %q", m5.createRepo)
	}
}

func TestRepoPickSelectedValidatesAndResolves(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha"}, nil
	}
	var validated string
	m.validateRepoPathFn = func(path string) (string, error) {
		validated = path
		// Return a different resolved root than the raw candidate to prove
		// the resolution is actually applied, not just a pass-through (e.g.
		// a zoxide candidate that is a non-root subdirectory of a repo).
		return "/resolved/alpha-root", nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)

	if validated != "/repos/alpha" {
		t.Errorf(
			"expected validateRepoPathFn called with the selected candidate, got %q",
			validated,
		)
	}
	if m5.createMode != createNameInput {
		t.Fatalf(
			"expected name-input mode after selecting a valid candidate, got mode=%d",
			m5.createMode,
		)
	}
	if m5.createRepo != "/resolved/alpha-root" {
		t.Errorf(
			"expected createRepo to be validateRepoPathFn's resolved root, got %q",
			m5.createRepo,
		)
	}
}

func TestRepoPickSelectedInvalidStaysInPicker(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha"}, nil
	}
	m.validateRepoPathFn = func(path string) (string, error) {
		return "", fmt.Errorf("not a git repository: %s", path)
	}
	createCalled := false
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)

	if m5.createMode != createRepoPick {
		t.Fatalf(
			"selecting a candidate that fails validation should stay in repo-pick mode, got mode=%d",
			m5.createMode,
		)
	}
	if m5.status == "" {
		t.Error("expected a status message describing the invalid candidate")
	}
	if createCalled {
		t.Error("createFn must not be called when the selected candidate fails validation")
	}
}

func TestRepoPickFreeTypedPathValidated(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha"}, nil
	}
	var validated string
	m.validateRepoPathFn = func(path string) (string, error) {
		validated = path
		return "/resolved/root", nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	for _, ch := range "/nonexistent/typed/path" {
		m2, _ = m3.Update(tea.KeyPressMsg{Code: ch})
		m3 = m2.(Model)
	}
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)

	if validated != "/nonexistent/typed/path" {
		t.Errorf("expected validateRepoPathFn called with the typed query, got %q", validated)
	}
	if m5.createMode != createNameInput {
		t.Fatalf(
			"expected name-input mode after a valid free-typed path, got mode=%d",
			m5.createMode,
		)
	}
	if m5.createRepo != "/resolved/root" {
		t.Errorf(
			"expected createRepo to be validateRepoPathFn's resolved root, got %q",
			m5.createRepo,
		)
	}
}

func TestRepoPickFreeTypedPathInvalidStaysInPicker(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha"}, nil
	}
	m.validateRepoPathFn = func(path string) (string, error) {
		return "", fmt.Errorf("not a git repository: %s", path)
	}
	createCalled := false
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	for _, ch := range "zzz-bad-path" {
		m2, _ = m3.Update(tea.KeyPressMsg{Code: ch})
		m3 = m2.(Model)
	}
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)

	if m5.createMode != createRepoPick {
		t.Fatalf(
			"invalid free-typed path should stay in repo-pick mode, got mode=%d",
			m5.createMode,
		)
	}
	if m5.status == "" {
		t.Error("expected a status message describing the invalid path")
	}
	if createCalled {
		t.Error("createFn must not be called when the typed path fails validation")
	}
}

func TestRepoPickEscReturnsToNormal(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(_ string) ([]string, error) {
		return []string{"/repos/alpha"}, nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	if m3.createMode != createRepoPick {
		t.Fatal("expected repo-pick mode after n")
	}

	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m5 := m4.(Model)
	if m5.createMode != createNone {
		t.Error("esc should return to normal mode")
	}
	if m5.repoPicker != nil {
		t.Error("esc should clear the repo picker")
	}
}

func TestNameInputTypingAccumulates(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"

	for _, ch := range "feat" {
		m2, _ := m.Update(tea.KeyPressMsg{Code: ch})
		m = m2.(Model)
	}
	if m.createInput.Value != "feat" {
		t.Errorf("expected createName 'feat', got %q", m.createInput.Value)
	}
}

func TestNameInputPasteInsertsInOneShot(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"

	m2, _ := m.Update(tea.PasteMsg{Content: "feat/pasted-branch"})
	m3 := m2.(Model)
	if m3.createInput.Value != "feat/pasted-branch" {
		t.Errorf("expected createName %q, got %q", "feat/pasted-branch", m3.createInput.Value)
	}
}

func TestNameInputPasteStripsControlChars(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")

	m2, _ := m.Update(tea.PasteMsg{Content: "/branch\n"})
	m3 := m2.(Model)
	if m3.createInput.Value != "feat/branch" {
		t.Errorf("expected createName %q, got %q", "feat/branch", m3.createInput.Value)
	}
}

func TestNameInputBackspaceRemovesLastRuneNotLastByte(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("café")

	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	m3 := m2.(Model)
	if m3.createInput.Value != "caf" {
		t.Errorf("expected createName %q, got %q", "caf", m3.createInput.Value)
	}
}

func TestRepoPickPasteInsertsIntoQuery(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.repoCandidatesFn = func(string) ([]string, error) {
		return []string{"/repos/repo-a", "/repos/other"}, nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'n'})
	m3 := m2.(Model)
	if m3.createMode != createRepoPick {
		t.Fatalf("expected repo-pick mode, got mode=%d", m3.createMode)
	}

	m4, _ := m3.Update(tea.PasteMsg{Content: "other"})
	m5 := m4.(Model)
	if m5.repoPicker.Query() != "other" {
		t.Errorf("expected query %q, got %q", "other", m5.repoPicker.Query())
	}
}

func TestNameInputEnterEmptyIsNoop(t *testing.T) {
	createCalled := false
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if cmd != nil {
		cmd()
	}
	if createCalled {
		t.Error("enter with an empty name must not call createFn")
	}
	if m3.createMode != createNameInput {
		t.Error("enter with an empty name should stay in name-input mode")
	}
}

func TestNameInputEscCancels(t *testing.T) {
	createCalled := false
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m3 := m2.(Model)
	if m3.createMode != createNone {
		t.Error("esc should return to normal mode")
	}
	if m3.createRepo != "" || m3.createInput.Value != "" {
		t.Error("esc should clear create state")
	}
	if createCalled {
		t.Error("esc must not call createFn")
	}
}

func TestCreateSuccessAttachesAndQuits(t *testing.T) {
	t.Setenv("TMUX", "test-session")

	var gotRepo, gotName string
	m := makeTestModel(testStatuses())
	m.mgr = &worktree.WorktreeManager{}
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.createFn = func(repoPath, name string) (string, error) {
		gotRepo = repoPath
		gotName = name
		return "", nil
	}
	attachCalled := false
	m.windowSessionFn = func(_ string) (string, bool) { return "sess", true }
	m.attachFn = func(_, _ string) error {
		attachCalled = true
		return nil
	}

	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if m3.createMode != createNone {
		t.Error("enter should leave name-input mode immediately, before the async create resolves")
	}
	if !m3.creating {
		t.Error("enter should set creating=true before the async create resolves")
	}
	if cmd == nil {
		t.Fatal("expected a command after enter")
	}

	msg := cmd()
	created, ok := msg.(createdMsg)
	if !ok {
		t.Fatalf("expected createdMsg after a successful createFn, got %T", msg)
	}
	if gotRepo != "/repos/alpha" || gotName != "feat" {
		t.Errorf("createFn called with wrong args: repo=%q name=%q", gotRepo, gotName)
	}

	m3b, cmd2 := m3.Update(created)
	msgs := flattenCmd(cmd2)

	if m3b.(Model).creating {
		t.Error("creating should be cleared once createdMsg is processed")
	}
	if !attachCalled {
		t.Error("expected attach to be attempted after a successful create")
	}
	foundQuit := false
	for _, mm := range msgs {
		if _, ok := mm.(tea.QuitMsg); ok {
			foundQuit = true
		}
	}
	if !foundQuit {
		t.Error("expected tea.QuitMsg among the resulting messages after a successful attach")
	}
}

func TestCreateFnFailureSetsStatusNoAttachNoQuit(t *testing.T) {
	t.Setenv("TMUX", "test-session")

	attachCalled := false
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.createFn = func(_, _ string) (string, error) {
		return "", fmt.Errorf("worktree already exists")
	}
	m.attachFn = func(_, _ string) error {
		attachCalled = true
		return nil
	}

	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if m3.createMode != createNone {
		t.Error(
			"enter should leave name-input mode immediately regardless of the eventual create result",
		)
	}
	if !m3.creating {
		t.Error("enter should set creating=true before the async create resolves")
	}
	if cmd == nil {
		t.Fatal("expected a command after enter")
	}

	msg := cmd()
	fm, ok := msg.(createFailedMsg)
	if !ok {
		t.Fatalf("expected createFailedMsg on create failure, got %T", msg)
	}
	if !strings.Contains(fm.err.Error(), "worktree already exists") {
		t.Errorf("expected the createFn error to be preserved, got %q", fm.err)
	}

	m4, _ := m3.Update(fm)
	m5 := m4.(Model)
	if m5.status == "" {
		t.Error("expected status set on the model after create failure")
	}
	if !strings.Contains(m5.status, "create failed") {
		t.Errorf("expected 'create failed' in status, got %q", m5.status)
	}
	if m5.creating {
		t.Error("creating should be cleared once createFailedMsg is processed")
	}
	if attachCalled {
		t.Error("attach must not be attempted when create fails")
	}
}

func TestCreateSuccessAttachFailureTriggersRefresh(t *testing.T) {
	t.Setenv("TMUX", "test-session")

	m := makeTestModel(testStatuses())
	m.mgr = &worktree.WorktreeManager{}
	m.createFn = func(_, _ string) (string, error) { return "", nil }
	m.windowSessionFn = func(_ string) (string, bool) { return "", false }
	m.repairFn = func(_, _ string, _ worktree.AICoder) error {
		return fmt.Errorf("repair unavailable")
	}

	created := createdMsg{repoPath: "/repos/alpha", name: "feat"}
	_, cmd := m.Update(created)
	msgs := flattenCmd(cmd)

	var gotStatus, gotStatuses bool
	for _, mm := range msgs {
		switch mm.(type) {
		case statusMsg:
			gotStatus = true
		case statusesMsg:
			gotStatuses = true
		}
	}
	if !gotStatus {
		t.Error("expected a failure statusMsg when attach and repair both fail")
	}
	if !gotStatuses {
		t.Error("expected a statusesMsg refresh even when attach fails, per the cycle plan")
	}
}

func TestNameInputEnterWithHookWarningRequiresSecondEnter(t *testing.T) {
	createCalled := false
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.checkHookCompatibilityFn = func(repoPath string) []string {
		if repoPath != "/repos/alpha" {
			t.Errorf("expected checkHookCompatibilityFn called with createRepo, got %q", repoPath)
		}
		return []string{"pre-commit (contains \"[ -d .git\")"}
	}
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	// First enter: arms the warning, must not call createFn or leave name-input mode.
	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if cmd != nil {
		cmd()
	}
	if createCalled {
		t.Error("createFn must not be called on the first enter when hook warnings exist")
	}
	if m3.createMode != createNameInput {
		t.Error("first enter with hook warnings should stay in name-input mode")
	}
	if !m3.pendingHookWarning {
		t.Error("first enter with hook warnings should arm pendingHookWarning")
	}
	if !strings.Contains(m3.status, "hook warning") {
		t.Errorf("expected status to describe the hook warning, got %q", m3.status)
	}

	// Second enter: confirms, calls createFn, and leaves name-input mode.
	m4, cmd2 := m3.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m5 := m4.(Model)
	if cmd2 == nil {
		t.Fatal("expected a command after the confirming second enter")
	}
	cmd2()
	if !createCalled {
		t.Error("createFn must be called on the second enter once the hook warning is confirmed")
	}
	if m5.createMode != createNone {
		t.Error("second enter should leave name-input mode")
	}
	if m5.pendingHookWarning {
		t.Error("pendingHookWarning should be cleared once the create is kicked off")
	}
}

func TestNameInputEnterWithoutHookWarningCreatesImmediately(t *testing.T) {
	createCalled := false
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.checkHookCompatibilityFn = func(_ string) []string { return nil }
	m.createFn = func(_, _ string) (string, error) {
		createCalled = true
		return "", nil
	}

	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if cmd == nil {
		t.Fatal("expected a command on the first enter when there are no hook warnings")
	}
	cmd()
	if !createCalled {
		t.Error("createFn must be called on the first enter when there are no hook warnings")
	}
	if m3.createMode != createNone {
		t.Error("enter with no hook warnings should leave name-input mode immediately")
	}
}

func TestNameInputEditingNameDearmsHookWarning(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.createMode = createNameInput
	m.createRepo = "/repos/alpha"
	m.createInput.SetValue("feat")
	m.checkHookCompatibilityFn = func(_ string) []string {
		return []string{"pre-commit (contains \"[ -d .git\")"}
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if !m3.pendingHookWarning {
		t.Fatal("expected pendingHookWarning armed after first enter")
	}

	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'x'})
	m5 := m4.(Model)
	if m5.pendingHookWarning {
		t.Error("typing after the warning was armed should de-arm it, like confirmThenRemove")
	}
}

func TestCreateSuccessWarningSurfacesAsStatus(t *testing.T) {
	t.Setenv("TMUX", "test-session")

	m := makeTestModel(testStatuses())
	m.mgr = &worktree.WorktreeManager{}
	m.windowSessionFn = func(_ string) (string, bool) { return "sess", true }
	m.attachFn = func(_, _ string) error { return nil }

	created := createdMsg{
		repoPath: "/repos/alpha",
		name:     "feat",
		warning:  "recent-repos store: disk full",
	}
	m2, cmd := m.Update(created)
	m3 := m2.(Model)
	if !strings.Contains(m3.status, "created, but:") || !strings.Contains(m3.status, "disk full") {
		t.Errorf("expected the warning to surface via status, got %q", m3.status)
	}

	msgs := flattenCmd(cmd)
	foundQuit := false
	for _, mm := range msgs {
		if _, ok := mm.(tea.QuitMsg); ok {
			foundQuit = true
		}
	}
	if !foundQuit {
		t.Error("a create warning must not prevent the attach/quit flow from proceeding")
	}
}
