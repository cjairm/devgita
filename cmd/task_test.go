package cmd

import (
	"fmt"
	"os"
	"testing"
)

// mockTaskRunner records calls to each task method.
type mockTaskRunner struct {
	refreshBranchArg    string
	refreshBranchCalled bool
	refreshBranchErr    error

	resetMainCalled bool
	resetMainErr    error

	reinstallLibsCalled bool
	reinstallLibsErr    error

	reinstallLibArg    string
	reinstallLibCalled bool
	reinstallLibErr    error

	deleteBranchArg    string
	deleteBranchCalled bool
	deleteBranchErr    error

	reviewScopeCalled bool
	reviewScopeRet    string
	reviewScopeErr    error

	branchDiffArg    string
	branchDiffCalled bool
	branchDiffRet    string
	branchDiffErr    error
}

func (m *mockTaskRunner) RefreshBranch(target string) error {
	m.refreshBranchCalled = true
	m.refreshBranchArg = target
	return m.refreshBranchErr
}

func (m *mockTaskRunner) ResetMainBranch() error {
	m.resetMainCalled = true
	return m.resetMainErr
}

func (m *mockTaskRunner) ReinstallLibraries() error {
	m.reinstallLibsCalled = true
	return m.reinstallLibsErr
}

func (m *mockTaskRunner) ReinstallLibrary(name string) error {
	m.reinstallLibCalled = true
	m.reinstallLibArg = name
	return m.reinstallLibErr
}

func (m *mockTaskRunner) DeleteBranch(target string) error {
	m.deleteBranchCalled = true
	m.deleteBranchArg = target
	return m.deleteBranchErr
}

func (m *mockTaskRunner) ReviewScope() (string, error) {
	m.reviewScopeCalled = true
	return m.reviewScopeRet, m.reviewScopeErr
}

func (m *mockTaskRunner) BranchDiff(file string) (string, error) {
	m.branchDiffCalled = true
	m.branchDiffArg = file
	return m.branchDiffRet, m.branchDiffErr
}

func setupTaskMock(t *testing.T, mock taskRunner) func() {
	t.Helper()
	orig := newTaskManager
	newTaskManager = func() taskRunner { return mock }
	return func() { newTaskManager = orig }
}

func TestTask_RefreshBranch(t *testing.T) {
	t.Run("no args defaults to empty target", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.refreshBranchCalled {
			t.Error("expected RefreshBranch to be called")
		}
		if mock.refreshBranchArg != "" {
			t.Errorf("expected empty target, got %q", mock.refreshBranchArg)
		}
	})

	t.Run("passes target arg", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{"develop"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.refreshBranchArg != "develop" {
			t.Errorf("expected target 'develop', got %q", mock.refreshBranchArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{refreshBranchErr: fmt.Errorf("git failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ResetMainBranch(t *testing.T) {
	t.Run("calls ResetMainBranch", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskResetMainBranchCmd.RunE(taskResetMainBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.resetMainCalled {
			t.Error("expected ResetMainBranch to be called")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{resetMainErr: fmt.Errorf("reset failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskResetMainBranchCmd.RunE(taskResetMainBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ReinstallLibraries(t *testing.T) {
	t.Run("calls ReinstallLibraries", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibrariesCmd.RunE(taskReinstallLibrariesCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.reinstallLibsCalled {
			t.Error("expected ReinstallLibraries to be called")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{reinstallLibsErr: fmt.Errorf("clean failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibrariesCmd.RunE(taskReinstallLibrariesCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ReinstallLibrary(t *testing.T) {
	t.Run("passes library name", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibraryCmd.RunE(taskReinstallLibraryCmd, []string{"lodash"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.reinstallLibCalled {
			t.Error("expected ReinstallLibrary to be called")
		}
		if mock.reinstallLibArg != "lodash" {
			t.Errorf("expected 'lodash', got %q", mock.reinstallLibArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{reinstallLibErr: fmt.Errorf("rm failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibraryCmd.RunE(taskReinstallLibraryCmd, []string{"lodash"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_DeleteBranch(t *testing.T) {
	t.Run("no args passes empty target", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.deleteBranchCalled {
			t.Error("expected DeleteBranch to be called")
		}
		if mock.deleteBranchArg != "" {
			t.Errorf("expected empty target, got %q", mock.deleteBranchArg)
		}
	})

	t.Run("passes target arg", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{"develop"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.deleteBranchArg != "develop" {
			t.Errorf("expected 'develop', got %q", mock.deleteBranchArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{deleteBranchErr: fmt.Errorf("delete failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ReviewScope(t *testing.T) {
	t.Run("calls ReviewScope and prints its output", func(t *testing.T) {
		mock := &mockTaskRunner{
			reviewScopeRet: "branch: feat/x -> main (default)  [ahead 1, behind 0]",
		}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReviewScopeCmd.RunE(taskReviewScopeCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.reviewScopeCalled {
			t.Error("expected ReviewScope to be called")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{reviewScopeErr: fmt.Errorf("git failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReviewScopeCmd.RunE(taskReviewScopeCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_BranchDiff(t *testing.T) {
	t.Run("no --file passes empty string", func(t *testing.T) {
		mock := &mockTaskRunner{branchDiffRet: "diff --git a/x b/x"}
		restore := setupTaskMock(t, mock)
		defer restore()
		taskBranchDiffFileFlag = ""

		err := taskBranchDiffCmd.RunE(taskBranchDiffCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.branchDiffCalled {
			t.Error("expected BranchDiff to be called")
		}
		if mock.branchDiffArg != "" {
			t.Errorf("expected empty file arg, got %q", mock.branchDiffArg)
		}
	})

	t.Run("passes --file flag", func(t *testing.T) {
		mock := &mockTaskRunner{branchDiffRet: "diff --git a/go.sum b/go.sum"}
		restore := setupTaskMock(t, mock)
		defer restore()
		taskBranchDiffFileFlag = "go.sum"
		defer func() { taskBranchDiffFileFlag = "" }()

		err := taskBranchDiffCmd.RunE(taskBranchDiffCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.branchDiffArg != "go.sum" {
			t.Errorf("expected file arg 'go.sum', got %q", mock.branchDiffArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{branchDiffErr: fmt.Errorf("diff failed")}
		restore := setupTaskMock(t, mock)
		defer restore()
		taskBranchDiffFileFlag = ""

		err := taskBranchDiffCmd.RunE(taskBranchDiffCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

// mockPRRunner records calls to the PR task methods.
type mockPRRunner struct {
	calls   []string
	lastArg map[string]string
	ret     string
	err     error
}

func newMockPRRunner() *mockPRRunner {
	return &mockPRRunner{lastArg: map[string]string{}, ret: "ok"}
}

func (m *mockPRRunner) record(name string, kv ...string) (string, error) {
	m.calls = append(m.calls, name)
	for i := 0; i+1 < len(kv); i += 2 {
		m.lastArg[kv[i]] = kv[i+1]
	}
	return m.ret, m.err
}

func (m *mockPRRunner) ReviewThreads(pr, state string) (string, error) {
	return m.record("ReviewThreads", "pr", pr, "state", state)
}

func (m *mockPRRunner) ResolveThread(id string) (string, error) {
	return m.record("ResolveThread", "id", id)
}

func (m *mockPRRunner) UnresolveThread(id string) (string, error) {
	return m.record("UnresolveThread", "id", id)
}

func (m *mockPRRunner) ReplyThread(id, body string) (string, error) {
	return m.record("ReplyThread", "id", id, "body", body)
}

func (m *mockPRRunner) SubmitReview(pr, verdict, body, comments string) (string, error) {
	return m.record(
		"SubmitReview",
		"pr", pr,
		"verdict", verdict,
		"body", body,
		"comments", comments,
	)
}

func (m *mockPRRunner) CreatePR(title, body, base string) (string, error) {
	return m.record("CreatePR", "title", title, "body", body, "base", base)
}

func (m *mockPRRunner) UpdatePRDescription(pr, body string) (string, error) {
	return m.record("UpdatePRDescription", "pr", pr, "body", body)
}

func (m *mockPRRunner) ApprovePR(pr, body string) (string, error) {
	return m.record("ApprovePR", "pr", pr, "body", body)
}

func (m *mockPRRunner) RequestChangesPR(pr, body string) (string, error) {
	return m.record("RequestChangesPR", "pr", pr, "body", body)
}

func (m *mockPRRunner) CommentPR(pr, body string) (string, error) {
	return m.record("CommentPR", "pr", pr, "body", body)
}

func (m *mockPRRunner) MergePR(pr, method string) (string, error) {
	return m.record("MergePR", "pr", pr, "method", method)
}
func (m *mockPRRunner) PRView(pr string) (string, error)   { return m.record("PRView", "pr", pr) }
func (m *mockPRRunner) PRChecks(pr string) (string, error) { return m.record("PRChecks", "pr", pr) }
func (m *mockPRRunner) CurrentPR() (string, error)         { return m.record("CurrentPR") }
func (m *mockPRRunner) CurrentRepo() (string, error)       { return m.record("CurrentRepo") }

func setupPRMock(t *testing.T, mock prRunner) func() {
	t.Helper()
	orig := newPRTasks
	newPRTasks = func() prRunner { return mock }
	return func() { newPRTasks = orig }
}

func TestPRTask_Dispatch(t *testing.T) {
	t.Run("review-threads passes flags", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()
		prFlag, prStateFlag = "7", "all"
		defer func() { prFlag, prStateFlag = "", "unresolved" }()

		if err := taskReviewThreadsCmd.RunE(taskReviewThreadsCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(mock.calls) != 1 || mock.calls[0] != "ReviewThreads" {
			t.Fatalf("expected ReviewThreads call, got %v", mock.calls)
		}
		if mock.lastArg["pr"] != "7" || mock.lastArg["state"] != "all" {
			t.Fatalf("unexpected args: %v", mock.lastArg)
		}
	})

	t.Run("resolve-thread passes id", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		if err := taskResolveThreadCmd.RunE(taskResolveThreadCmd, []string{"PRRT_x"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.calls[0] != "ResolveThread" || mock.lastArg["id"] != "PRRT_x" {
			t.Fatalf("unexpected dispatch: %v / %v", mock.calls, mock.lastArg)
		}
	})

	t.Run("reply-thread passes id and body", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		if err := taskReplyThreadCmd.RunE(
			taskReplyThreadCmd,
			[]string{"PRRT_x", "fixed"},
		); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.lastArg["id"] != "PRRT_x" || mock.lastArg["body"] != "fixed" {
			t.Fatalf("unexpected args: %v", mock.lastArg)
		}
	})

	t.Run("create-pr passes title/body/base", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()
		prTitleFlag, prBodyFlag, prBaseFlag = "T", "B", "main"
		defer func() { prTitleFlag, prBodyFlag, prBaseFlag = "", "", "" }()

		if err := taskCreatePRCmd.RunE(taskCreatePRCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.lastArg["title"] != "T" || mock.lastArg["body"] != "B" ||
			mock.lastArg["base"] != "main" {
			t.Fatalf("unexpected args: %v", mock.lastArg)
		}
	})

	t.Run("merge-pr passes method", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()
		prFlag, prMethodFlag = "9", "rebase"
		defer func() { prFlag, prMethodFlag = "", "squash" }()

		if err := taskMergePRCmd.RunE(taskMergePRCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.lastArg["pr"] != "9" || mock.lastArg["method"] != "rebase" {
			t.Fatalf("unexpected args: %v", mock.lastArg)
		}
	})

	t.Run("pr-view dispatches", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()
		prFlag = ""

		if err := taskPRViewCmd.RunE(taskPRViewCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.calls[0] != "PRView" {
			t.Fatalf("expected PRView, got %v", mock.calls)
		}
	})

	t.Run("current-pr dispatches", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		if err := taskCurrentPRCmd.RunE(taskCurrentPRCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.calls[0] != "CurrentPR" {
			t.Fatalf("expected CurrentPR, got %v", mock.calls)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := newMockPRRunner()
		mock.err = fmt.Errorf("boom")
		defer setupPRMock(t, mock)()

		if err := taskPRChecksCmd.RunE(taskPRChecksCmd, nil); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestPRTask_BodyFile(t *testing.T) {
	t.Run("create-pr reads markdown body from file verbatim", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		md := "## Summary\n\n- Adds `dg task`\n- **Bold** and a list\n\n```go\nfmt.Println(\"hi\")\n```\n"
		dir := t.TempDir()
		path := dir + "/body.md"
		if err := os.WriteFile(path, []byte(md), 0o644); err != nil {
			t.Fatalf("write temp: %v", err)
		}

		prTitleFlag, prBodyFlag, prBodyFileFlag, prBaseFlag = "T", "", path, ""
		defer func() { prTitleFlag, prBodyFlag, prBodyFileFlag, prBaseFlag = "", "", "", "" }()

		if err := taskCreatePRCmd.RunE(taskCreatePRCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.lastArg["body"] != md {
			t.Fatalf(
				"markdown body not passed verbatim.\nwant: %q\ngot:  %q",
				md,
				mock.lastArg["body"],
			)
		}
	})

	t.Run("both --body and --body-file is an error", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		dir := t.TempDir()
		path := dir + "/b.md"
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("write temp: %v", err)
		}
		prTitleFlag, prBodyFlag, prBodyFileFlag = "T", "inline", path
		defer func() { prTitleFlag, prBodyFlag, prBodyFileFlag = "", "", "" }()

		if err := taskCreatePRCmd.RunE(taskCreatePRCmd, nil); err == nil {
			t.Fatal("expected error when both --body and --body-file are set")
		}
		if len(mock.calls) != 0 {
			t.Fatal("expected no dispatch when flags conflict")
		}
	})

	t.Run("missing body-file surfaces a read error", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		prTitleFlag, prBodyFlag, prBodyFileFlag = "T", "", "/no/such/file.md"
		defer func() { prTitleFlag, prBodyFlag, prBodyFileFlag = "", "", "" }()

		if err := taskCreatePRCmd.RunE(taskCreatePRCmd, nil); err == nil {
			t.Fatal("expected error for unreadable --body-file")
		}
	})

	t.Run("submit-review reads event, body-file, and comments-file", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		dir := t.TempDir()
		bodyPath := dir + "/review.md"
		commentsPath := dir + "/comments.json"
		comments := `[{"path":"a.go","line":1,"body":"x"}]`
		if err := os.WriteFile(bodyPath, []byte("## Review"), 0o644); err != nil {
			t.Fatalf("write body: %v", err)
		}
		if err := os.WriteFile(commentsPath, []byte(comments), 0o644); err != nil {
			t.Fatalf("write comments: %v", err)
		}

		prFlag, prEventFlag, prBodyFileFlag, prCommentsFile = "42", "request-changes", bodyPath, commentsPath
		defer func() {
			prFlag, prEventFlag, prBodyFlag, prBodyFileFlag, prCommentsFile = "", "", "", "", ""
		}()

		if err := taskSubmitReviewCmd.RunE(taskSubmitReviewCmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.calls[0] != "SubmitReview" {
			t.Fatalf("expected SubmitReview, got %v", mock.calls)
		}
		if mock.lastArg["pr"] != "42" || mock.lastArg["verdict"] != "request-changes" ||
			mock.lastArg["body"] != "## Review" || mock.lastArg["comments"] != comments {
			t.Fatalf("unexpected args: %v", mock.lastArg)
		}
	})

	t.Run("submit-review surfaces an unreadable comments-file", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		prEventFlag, prCommentsFile = "comment", "/no/such/comments.json"
		defer func() { prEventFlag, prCommentsFile, prBodyFlag = "", "", "" }()

		if err := taskSubmitReviewCmd.RunE(taskSubmitReviewCmd, nil); err == nil {
			t.Fatal("expected error for unreadable --comments-file")
		}
		if len(mock.calls) != 0 {
			t.Fatal("expected no dispatch when comments-file is unreadable")
		}
	})

	t.Run("reply-thread accepts --body-file", func(t *testing.T) {
		mock := newMockPRRunner()
		defer setupPRMock(t, mock)()

		dir := t.TempDir()
		path := dir + "/reply.md"
		if err := os.WriteFile(path, []byte("Fixed in **abc123**"), 0o644); err != nil {
			t.Fatalf("write temp: %v", err)
		}
		prBodyFileFlag = path
		defer func() { prBodyFileFlag = "" }()

		if err := taskReplyThreadCmd.RunE(taskReplyThreadCmd, []string{"PRRT_x"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.lastArg["body"] != "Fixed in **abc123**" {
			t.Fatalf("unexpected body: %q", mock.lastArg["body"])
		}
	})
}
