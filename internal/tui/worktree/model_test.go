package tuiworktree

import (
	"fmt"
	"os"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/internal/tooling/task"
	"github.com/cjairm/devgita/internal/tooling/worktree"
	tuicomponents "github.com/cjairm/devgita/internal/tui/components"
)

func init() { testutil.InitLogger() }

func makeTestModel(statuses []worktree.WorktreeStatus) Model {
	m := Model{
		collapsed:     map[string]bool{},
		palette:       tuicomponents.NewPalette(),
		leftPaneWidth: minLeftPaneWidth,
		width:         120,
		height:        40,
	}
	m.diffFn = func(_ string) (task.BranchDiffResult, error) {
		return task.BranchDiffResult{Content: "diff content", Files: 1, Added: 5, Removed: 2}, nil
	}
	m.attachFn = func(_, _ string) error { return nil }
	m.removeFn = func(_, _ string, _ bool) error { return nil }
	m.removeSessionFn = func(_, _ string) error { return nil }
	m.repairFn = func(_, _ string, _ worktree.AICoder) error { return nil }
	m.windowSessionFn = func(_ string) (string, bool) { return "", false }
	m.repoCandidatesFn = func(_ string) ([]string, error) { return nil, nil }
	m.validateRepoPathFn = func(path string) (string, error) { return path, nil }
	m.checkHookCompatibilityFn = func(_ string) []string { return nil }
	m.createFn = func(_, _ string) (string, error) { return "", nil }
	m.statuses = statuses
	m.rebuildRows()
	return m
}

// flattenCmd runs cmd and recursively unwraps any tea.BatchMsg it produces,
// so tests can assert on the individual messages a tea.Batch yields without
// needing the full bubbletea runtime to fan them back into Update.
func flattenCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		var out []tea.Msg
		for _, c := range batch {
			out = append(out, flattenCmd(c)...)
		}
		return out
	}
	return []tea.Msg{msg}
}

func testStatuses() []worktree.WorktreeStatus {
	return []worktree.WorktreeStatus{
		{
			Name:         "feature-a",
			Repo:         "repo-a",
			Path:         "/tmp/a",
			TmuxWindow:   "wt-feature-a",
			WindowActive: true,
		},
		{
			Name:         "feature-b",
			Repo:         "repo-a",
			Path:         "/tmp/b",
			TmuxWindow:   "wt-feature-b",
			WindowActive: false,
		},
		{
			Name:         "feature-x",
			Repo:         "repo-b",
			Path:         "/tmp/x",
			TmuxWindow:   "wt-feature-x",
			WindowActive: true,
		},
	}
}

func TestBuildRowsGrouping(t *testing.T) {
	statuses := testStatuses()
	rows := buildRows(statuses, map[string]bool{}, "")
	// Should have: repo-a header, feature-a, feature-b, repo-b header, feature-x
	if len(rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(rows))
	}
	if rows[0].kind != rowRepo || rows[0].repo != "repo-a" {
		t.Error("expected first row to be repo-a header")
	}
	if rows[1].kind != rowWorktree || rows[1].status.Name != "feature-a" {
		t.Error("expected second row to be feature-a")
	}
	if rows[3].kind != rowRepo || rows[3].repo != "repo-b" {
		t.Error("expected fourth row to be repo-b header")
	}
}

func TestCursorSkipsRepoHeaders(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Initial cursor should be on a worktree row
	if m.rows[m.cursor].kind != rowWorktree {
		t.Error("initial cursor should be on a worktree row")
	}
	// Move down
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'j'})
	m3 := m2.(Model)
	if m3.rows[m3.cursor].kind != rowWorktree {
		t.Error("after j, cursor should still be on a worktree row")
	}
}

func TestFoldHide(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Collapse repo-a
	m.collapsed["repo-a"] = true
	m.rebuildRows()
	// Should have: repo-a header, repo-b header, feature-x
	if len(m.rows) != 3 {
		t.Fatalf("expected 3 rows after collapse, got %d", len(m.rows))
	}
}

func TestFoldUnfold(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Start on a worktree in repo-a
	if m.rows[m.cursor].kind != rowWorktree || m.rows[m.cursor].status.Repo != "repo-a" {
		t.Fatal("expected initial cursor on repo-a worktree")
	}

	// h collapses the repo and lands cursor on the repo header
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'h'})
	m3 := m2.(Model)
	if m3.rows[m3.cursor].kind != rowRepo {
		t.Fatalf("after h, cursor should be on repo header, got kind=%d", m3.rows[m3.cursor].kind)
	}
	if m3.rows[m3.cursor].repo != "repo-a" {
		t.Errorf("cursor should be on repo-a header, got %q", m3.rows[m3.cursor].repo)
	}

	// l expands it and returns cursor to a worktree inside repo-a
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'l'})
	m5 := m4.(Model)
	if m5.rows[m5.cursor].kind != rowWorktree {
		t.Fatal("after l, cursor should be on a worktree row")
	}
	if m5.rows[m5.cursor].status.Repo != "repo-a" {
		t.Errorf("after l, cursor should be in repo-a, got %q", m5.rows[m5.cursor].status.Repo)
	}
}

func TestCollapsedHeaderReachableAfterNavAway(t *testing.T) {
	m := makeTestModel(testStatuses())

	// Collapse repo-a — cursor lands on repo-a header
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'h'})
	m3 := m2.(Model)
	if m3.rows[m3.cursor].kind != rowRepo || m3.rows[m3.cursor].repo != "repo-a" {
		t.Fatal("after h, cursor should be on repo-a header")
	}

	// Navigate away with j — should reach feature-x in repo-b
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'j'})
	m5 := m4.(Model)
	if m5.rows[m5.cursor].kind != rowWorktree || m5.rows[m5.cursor].status.Repo != "repo-b" {
		t.Fatalf("after j from collapsed header, expected repo-b worktree, got kind=%d repo=%q",
			m5.rows[m5.cursor].kind, m5.rows[m5.cursor].status.Repo)
	}

	// Navigate back with j — should wrap to repo-a header (collapsed, navigable)
	m6, _ := m5.Update(tea.KeyPressMsg{Code: 'j'})
	m7 := m6.(Model)
	if m7.rows[m7.cursor].kind != rowRepo || m7.rows[m7.cursor].repo != "repo-a" {
		t.Fatalf("after wrap, expected repo-a header, got kind=%d repo=%q",
			m7.rows[m7.cursor].kind, m7.rows[m7.cursor].repo)
	}

	// l from repo-a header should expand and land on worktree in repo-a
	m8, _ := m7.Update(tea.KeyPressMsg{Code: 'l'})
	m9 := m8.(Model)
	if m9.rows[m9.cursor].kind != rowWorktree || m9.rows[m9.cursor].status.Repo != "repo-a" {
		t.Error("after l on re-reached header, cursor should be in repo-a worktree")
	}
}

func TestFilterMode(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Enter filter mode
	m2, _ := m.Update(tea.KeyPressMsg{Code: '/'})
	m3 := m2.(Model)
	if !m3.filter.Active {
		t.Error("should be in filtering mode after /")
	}
	// Type a char
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'b'})
	m5 := m4.(Model)
	if m5.filter.Text != "b" {
		t.Errorf("expected filter 'b', got %q", m5.filter.Text)
	}
	// Esc clears and exits
	m6, _ := m5.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m7 := m6.(Model)
	if m7.filter.Active || m7.filter.Text != "" {
		t.Error("esc should clear filter and exit filtering mode")
	}
}

func TestAttachOutsideTmux(t *testing.T) {
	os.Unsetenv("TMUX") //nolint:errcheck
	m := makeTestModel(testStatuses())
	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m3 := m2.(Model)
	if m3.status == "" {
		t.Error("should set status when not in tmux, not quit")
	}
}

func TestDeleteDoubleConfirm(t *testing.T) {
	removeCalled := false
	var removedRepo, removedName string

	m := makeTestModel(testStatuses())
	m.removeFn = func(repo, name string, force bool) error {
		removeCalled = true
		removedRepo = repo
		removedName = name
		return nil
	}

	// First d: arm
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'd'})
	m3 := m2.(Model)
	if removeCalled {
		t.Error("first d should not delete")
	}
	if m3.pendingDelete == "" {
		t.Error("first d should arm pendingDelete")
	}

	// Non-d key clears arm
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'j'})
	m5 := m4.(Model)
	if m5.pendingDelete != "" {
		t.Error("j should clear pendingDelete")
	}
	if removeCalled {
		t.Error("no delete should have happened")
	}

	// Second d on same row deletes
	m6, _ := m.Update(tea.KeyPressMsg{Code: 'd'})
	m7 := m6.(Model)
	m8, cmd := m7.Update(tea.KeyPressMsg{Code: 'd'})
	_ = m8
	if cmd != nil {
		// Execute the command
		msg := cmd()
		_ = msg
	}
	if !removeCalled {
		t.Error("second d should call removeFn")
	}
	if removedRepo != "repo-a" {
		t.Errorf("expected repo 'repo-a', got %q", removedRepo)
	}
	if removedName != "feature-a" {
		t.Errorf("expected name 'feature-a', got %q", removedName)
	}
}

func TestSessionDeleteDoubleConfirm(t *testing.T) {
	removeCalled := false
	var removedRepo, removedName string

	m := makeTestModel(testStatuses())
	m.removeSessionFn = func(repo, name string) error {
		removeCalled = true
		removedRepo = repo
		removedName = name
		return nil
	}

	// First D: arm
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'D'})
	m3 := m2.(Model)
	if removeCalled {
		t.Error("first D should not delete")
	}
	if m3.pendingSessionDelete == "" {
		t.Error("first D should arm pendingSessionDelete")
	}

	// Non-D key clears arm
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'j'})
	m5 := m4.(Model)
	if m5.pendingSessionDelete != "" {
		t.Error("j should clear pendingSessionDelete")
	}
	if removeCalled {
		t.Error("no delete should have happened")
	}

	// d does not confirm a D-armed delete
	m6, _ := m3.Update(tea.KeyPressMsg{Code: 'd'})
	m7 := m6.(Model)
	if m7.pendingSessionDelete != "" {
		t.Error("d should clear pendingSessionDelete instead of confirming it")
	}
	if removeCalled {
		t.Error("d after D must not trigger the session delete")
	}

	// Second D on same row deletes worktree + session
	m8, _ := m.Update(tea.KeyPressMsg{Code: 'D'})
	m9 := m8.(Model)
	m10, cmd := m9.Update(tea.KeyPressMsg{Code: 'D'})
	_ = m10
	if cmd != nil {
		cmd()
	}
	if !removeCalled {
		t.Error("second D should call removeSessionFn")
	}
	if removedRepo != "repo-a" {
		t.Errorf("expected repo 'repo-a', got %q", removedRepo)
	}
	if removedName != "feature-a" {
		t.Errorf("expected name 'feature-a', got %q", removedName)
	}
}

func TestSessionDeleteErrorPropagation(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.removeSessionFn = func(_, _ string) error {
		return fmt.Errorf("no such session")
	}

	m2, _ := m.Update(tea.KeyPressMsg{Code: 'D'})
	m3 := m2.(Model)
	m4, cmd := m3.Update(tea.KeyPressMsg{Code: 'D'})
	if cmd == nil {
		t.Fatal("expected a command after second D")
	}
	msg := cmd()
	m5, _ := m4.(Model).Update(msg)
	m6 := m5.(Model)
	if m6.status == "" {
		t.Error("expected inline status message on session delete error")
	}
	if m6.pendingSessionDelete != "" {
		t.Error("pendingSessionDelete should be cleared after error")
	}
}

func TestDeleteDuplicateNameAcrossRepos(t *testing.T) {
	statuses := []worktree.WorktreeStatus{
		{Name: "feature-x", Repo: "repo-a", Path: "/tmp/ax", WindowActive: false},
		{Name: "feature-x", Repo: "repo-b", Path: "/tmp/bx", WindowActive: false},
	}
	var calledRepo string
	m := makeTestModel(statuses)
	m.removeFn = func(repo, name string, force bool) error {
		calledRepo = repo
		return nil
	}
	// Navigate to the repo-b/feature-x row
	m.rebuildRows()
	for i, r := range m.rows {
		if r.kind == rowWorktree && r.status.Repo == "repo-b" {
			m.cursor = i
			break
		}
	}
	// Double d
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'd'})
	m3 := m2.(Model)
	m4, cmd := m3.Update(tea.KeyPressMsg{Code: 'd'})
	_ = m4
	if cmd != nil {
		cmd()
	}
	if calledRepo != "repo-b" {
		t.Errorf("expected repo-b to be deleted, got %q", calledRepo)
	}
}

func TestRepairCallsRepairFn(t *testing.T) {
	repairCalled := false
	var repairedRepo, repairedName string
	m := makeTestModel(testStatuses())
	m.repairFn = func(repo, name string, coder worktree.AICoder) error {
		repairCalled = true
		repairedRepo = repo
		repairedName = name
		return nil
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'r'})
	if cmd != nil {
		cmd()
	}
	if !repairCalled {
		t.Error("r should call repairFn")
	}
	_ = repairedRepo
	_ = repairedName
}

func TestFilterHidesNonMatchingRows(t *testing.T) {
	m := makeTestModel(testStatuses())
	totalBefore := len(m.rows)
	if totalBefore == 0 {
		t.Fatal("expected rows before filter")
	}

	// Enter filter mode and type "repo-b"
	m2, _ := m.Update(tea.KeyPressMsg{Code: '/'})
	for _, ch := range "repo-b" {
		m2, _ = m2.(Model).Update(tea.KeyPressMsg{Code: ch})
	}
	m3 := m2.(Model)

	// Only repo-b header + feature-x should remain
	if len(m3.rows) >= totalBefore {
		t.Errorf("filter should reduce rows: before=%d after=%d", totalBefore, len(m3.rows))
	}
	for _, r := range m3.rows {
		if r.kind == rowWorktree && r.status.Repo != "repo-b" {
			t.Errorf("filter left row from wrong repo: %s/%s", r.status.Repo, r.status.Name)
		}
	}

	// Esc clears filter and restores all rows
	m4, _ := m3.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m5 := m4.(Model)
	if m5.filter.Text != "" {
		t.Error("esc should clear filter string")
	}
	if len(m5.rows) != totalBefore {
		t.Errorf("after esc rows should be restored: want %d got %d", totalBefore, len(m5.rows))
	}
}

func TestCursorWrapsAtBoundaries(t *testing.T) {
	m := makeTestModel(testStatuses())
	indices := worktreeIndices(m.rows)
	if len(indices) < 2 {
		t.Skip("need at least 2 worktree rows")
	}

	// Start at first worktree row
	m.cursor = indices[0]

	// k on first row should wrap to last
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'k'})
	m3 := m2.(Model)
	if m3.cursor != indices[len(indices)-1] {
		t.Errorf("k on first row should wrap to last, got cursor=%d want=%d",
			m3.cursor, indices[len(indices)-1])
	}

	// j on last row should wrap to first
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'j'})
	m5 := m4.(Model)
	if m5.cursor != indices[0] {
		t.Errorf("j on last row should wrap to first, got cursor=%d want=%d",
			m5.cursor, indices[0])
	}
}

func TestDeleteErrorPropagation(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.removeFn = func(_, _ string, _ bool) error {
		return fmt.Errorf("disk full")
	}

	// Arm
	m2, _ := m.Update(tea.KeyPressMsg{Code: 'd'})
	m3 := m2.(Model)
	// Confirm
	m4, cmd := m3.Update(tea.KeyPressMsg{Code: 'd'})
	_ = m4
	if cmd == nil {
		t.Fatal("expected a command after second d")
	}
	msg := cmd()
	m5, _ := m4.(Model).Update(msg)
	m6 := m5.(Model)
	if m6.status == "" {
		t.Error("expected inline status message on delete error")
	}
	if m6.pendingDelete != "" {
		t.Error("pendingDelete should be cleared after error")
	}
}

func TestDiffFocusMode(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.diffContent = "line0\nline1\nline2\nline3\nline4"
	m.diffFileLines = []int{0, 3}

	// space focuses the diff pane
	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m3 := m2.(Model)
	if !m3.diffFocused {
		t.Fatal("space should focus the diff pane")
	}

	// j scrolls the diff without moving the list cursor
	cursorBefore := m3.cursor
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'j'})
	m5 := m4.(Model)
	if m5.diffScroll != 1 {
		t.Errorf("j should scroll diff by 1, got %d", m5.diffScroll)
	}
	if m5.cursor != cursorBefore {
		t.Error("j while focused must not move the list cursor")
	}

	// ] jumps to the next file header
	m6, _ := m5.Update(tea.KeyPressMsg{Code: ']'})
	m7 := m6.(Model)
	if m7.diffScroll != 3 {
		t.Errorf("] should jump to next file header line 3, got %d", m7.diffScroll)
	}

	// [ jumps back to the previous file header
	m8, _ := m7.Update(tea.KeyPressMsg{Code: '['})
	m9 := m8.(Model)
	if m9.diffScroll != 0 {
		t.Errorf("[ should jump back to header line 0, got %d", m9.diffScroll)
	}

	// G/g hit bottom/top
	m10, _ := m9.Update(tea.KeyPressMsg{Code: 'G'})
	m11 := m10.(Model)
	if m11.diffScroll != 4 {
		t.Errorf("G should scroll to last line, got %d", m11.diffScroll)
	}

	// d while focused must not arm a delete
	m12, _ := m11.Update(tea.KeyPressMsg{Code: 'd'})
	m13 := m12.(Model)
	if m13.pendingDelete != "" {
		t.Error("d while focused must not arm pendingDelete")
	}

	// esc returns focus to the list
	m14, _ := m13.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m15 := m14.(Model)
	if m15.diffFocused {
		t.Error("esc should unfocus the diff pane")
	}
}

func TestDiffFocusRequiresContent(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.diffContent = ""
	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m2.(Model).diffFocused {
		t.Error("space must not focus an empty (still loading) diff pane")
	}
}

func TestDiffHeaderShowsComparison(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.diffBase = "main @3e90667"
	m.diffBranch = "feature-a"
	m.diffFiles, m.diffAdded, m.diffRemoved = 2, 96, 14
	m.diffContent = "x"
	right := m.renderRight(100)
	header := strings.Split(ansi.Strip(right), "\n")[0]
	if !strings.Contains(header, "main @3e90667 ← feature-a") {
		t.Errorf("expected comparison label in header, got %q", header)
	}
	if !strings.Contains(header, "±2 +96 -14") {
		t.Errorf("expected stat line in header, got %q", header)
	}
}

func TestNarrowTerminalNoPanic(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Very narrow terminal
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 30, Height: 20})
	m3 := m2.(Model)
	// Should not panic when rendering
	v := m3.View()
	_ = v
	if m3.rightPaneWidth() < 0 {
		t.Error("rightPaneWidth should not be negative")
	}
}

func TestHelpOverlayShowsDashboardBackground(t *testing.T) {
	m := makeTestModel(testStatuses())
	m.showHelp = true
	out := ansi.Strip(m.renderContent())

	if !strings.Contains(out, "press any key to close") {
		t.Error("expected the help popup to be rendered")
	}
	if !strings.Contains(out, "repo-a") {
		t.Error(
			"expected the dashboard background (repo-a) to remain visible behind the help popup",
		)
	}
	if !strings.Contains(out, "feature-a") {
		t.Error(
			"expected the dashboard background (feature-a) to remain visible behind the help popup",
		)
	}
}

// Create-flow (n → repo-pick → name-input → create) tests live in
// create_flow_test.go, mirroring the create_flow.go/model.go split.
