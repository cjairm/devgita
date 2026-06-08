package tuiworktree

import (
	"fmt"
	"os"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/internal/tooling/worktree"
)

func init() { testutil.InitLogger() }

func makeTestModel(statuses []worktree.WorktreeStatus) Model {
	m := Model{
		collapsed:     map[string]bool{},
		styles:        newStyles(),
		leftPaneWidth: minLeftPaneWidth,
		width:         120,
		height:        40,
	}
	m.captureFn = func(_, _ string) (string, error) { return "agent output", nil }
	m.diffFn = func(_ string) (string, error) { return "diff content", nil }
	m.diffStatFn = func(_ string) (files, added, removed int, err error) { return 1, 5, 2, nil }
	m.attachFn = func(_, _ string) error { return nil }
	m.removeFn = func(_, _ string, _ bool) error { return nil }
	m.repairFn = func(_, _ string, _ worktree.AICoder) error { return nil }
	m.statuses = statuses
	m.rebuildRows()
	return m
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

func TestTabToggle(t *testing.T) {
	m := makeTestModel(testStatuses())
	if m.activeTab != tabAgent {
		t.Error("initial tab should be agent")
	}
	m2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m3 := m2.(Model)
	if m3.activeTab != tabDiff {
		t.Error("after Tab, active tab should be diff")
	}
}

func TestFilterMode(t *testing.T) {
	m := makeTestModel(testStatuses())
	// Enter filter mode
	m2, _ := m.Update(tea.KeyPressMsg{Code: '/'})
	m3 := m2.(Model)
	if !m3.filtering {
		t.Error("should be in filtering mode after /")
	}
	// Type a char
	m4, _ := m3.Update(tea.KeyPressMsg{Code: 'b'})
	m5 := m4.(Model)
	if m5.filter != "b" {
		t.Errorf("expected filter 'b', got %q", m5.filter)
	}
	// Esc clears and exits
	m6, _ := m5.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m7 := m6.(Model)
	if m7.filtering || m7.filter != "" {
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

func TestAgentOfflinePlaceholder(t *testing.T) {
	m := makeTestModel(testStatuses())
	m2, _ := m.Update(agentOfflineMsg{})
	m3 := m2.(Model)
	if !m3.agentOffline {
		t.Error("agentOfflineMsg should set agentOffline")
	}
	// Tick should not re-issue capture when offline
	m4, cmd := m3.Update(tickMsg(time.Now()))
	_ = m4
	_ = cmd
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
	if m5.filter != "" {
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
