package jq

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
)

func init() { logger.Init(false) }

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Jq{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != "jq" {
		t.Fatalf("expected InstallPackage(jq), got %q", mc.InstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Jq{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != "jq" {
		t.Fatalf("expected MaybeInstallPackage(jq), got %q", mc.MaybeInstalled)
	}
}

func TestForceConfigure(t *testing.T) {
	app := &Jq{Cmd: commands.NewMockCommand()}
	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}
}

func TestSoftConfigure(t *testing.T) {
	app := &Jq{Cmd: commands.NewMockCommand()}
	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}
}

func TestExecuteCommand(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Jq{Cmd: mc, Base: mockBase}

	t.Run("successful execution", func(t *testing.T) {
		mockBase.SetExecCommandResult(`{"key":"value"}`, "", nil)

		if err := app.ExecuteCommand(".", "file.json"); err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("expected 1 ExecCommand call, got %d", mockBase.GetExecCommandCallCount())
		}
		lastCall := mockBase.GetLastExecCommandCall()
		if lastCall.Command != "jq" {
			t.Fatalf("expected command 'jq', got %q", lastCall.Command)
		}
		if len(lastCall.Args) != 2 || lastCall.Args[0] != "." || lastCall.Args[1] != "file.json" {
			t.Fatalf("unexpected args: %v", lastCall.Args)
		}
		if lastCall.IsSudo {
			t.Fatal("expected IsSudo to be false")
		}
	})

	t.Run("command execution error", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "parse error", fmt.Errorf("exit 2"))

		err := app.ExecuteCommand(".invalid")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "failed to run jq command") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func boolPtr(b bool) *bool { return &b }

func TestFormatReviewThreads(t *testing.T) {
	const sampleJSON = `{"data":{"repository":{"pullRequest":{"reviewThreads":{"nodes":[]}}}}}`

	// assertArgSeq verifies that wantSeq appears as a contiguous subsequence of args.
	assertArgSeq := func(t *testing.T, args []string, wantSeq ...string) {
		t.Helper()
		for i := 0; i+len(wantSeq) <= len(args); i++ {
			match := true
			for j, w := range wantSeq {
				if args[i+j] != w {
					match = false
					break
				}
			}
			if match {
				return
			}
		}
		t.Fatalf("expected args to contain sequence %v, got %v", wantSeq, args)
	}

	t.Run("unresolved passes resolved=false", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
		mockBase.SetExecCommandResult("## Location\nFile: `a.go`", "", nil)

		out, err := app.FormatReviewThreads(sampleJSON, boolPtr(false))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != "## Location\nFile: `a.go`" {
			t.Fatalf("unexpected output: %q", out)
		}

		call := mockBase.GetLastExecCommandCall()
		if call == nil {
			t.Fatal("no ExecCommand call recorded")
		}
		if call.Command != "jq" {
			t.Fatalf("expected command 'jq', got %q", call.Command)
		}
		if call.Args[0] != "-r" {
			t.Fatalf("expected first arg '-r', got %q", call.Args[0])
		}
		assertArgSeq(t, call.Args, "--argjson", "resolved", "false")

		// The filter must be present and the JSON must be passed as a temp file.
		if !strings.Contains(call.Args[4], "reviewThreads") {
			t.Fatalf("expected jq filter in args, got %q", call.Args[4])
		}
		last := call.Args[len(call.Args)-1]
		if !strings.HasSuffix(last, ".json") {
			t.Fatalf("expected last arg to be a .json temp file, got %q", last)
		}
	})

	t.Run("resolved passes resolved=true", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
		mockBase.SetExecCommandResult("", "", nil)

		if _, err := app.FormatReviewThreads(sampleJSON, boolPtr(true)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertArgSeq(t, mockBase.GetLastExecCommandCall().Args, "--argjson", "resolved", "true")
	})

	t.Run("nil resolved passes null", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
		mockBase.SetExecCommandResult("", "", nil)

		if _, err := app.FormatReviewThreads(sampleJSON, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertArgSeq(t, mockBase.GetLastExecCommandCall().Args, "--argjson", "resolved", "null")
	})

	t.Run("surfaces jq stderr on failure", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
		mockBase.SetExecCommandResult("", "jq: error: syntax", fmt.Errorf("exit 3"))

		_, err := app.FormatReviewThreads(sampleJSON, boolPtr(false))
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "jq: error: syntax") {
			t.Fatalf("expected stderr in error, got: %v", err)
		}
	})

	t.Run("wraps error when stderr empty", func(t *testing.T) {
		mockBase := commands.NewMockBaseCommand()
		app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
		mockBase.SetExecCommandResult("", "", fmt.Errorf("jq not found"))

		_, err := app.FormatReviewThreads(sampleJSON, nil)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "failed to run filter") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFormatPRView(t *testing.T) {
	mockBase := commands.NewMockBaseCommand()
	app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
	mockBase.SetExecCommandResult("PR #42: Title\nstate: OPEN", "", nil)

	out, err := app.FormatPRView(`{"number":42}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "PR #42: Title\nstate: OPEN" {
		t.Fatalf("unexpected output: %q", out)
	}
	call := mockBase.GetLastExecCommandCall()
	if call.Command != "jq" || call.Args[0] != "-r" {
		t.Fatalf("expected jq -r, got %v / %v", call.Command, call.Args)
	}
	if !strings.Contains(call.Args[1], "PR #") {
		t.Fatalf("expected prViewFilter, got %q", call.Args[1])
	}
	if !strings.HasSuffix(call.Args[len(call.Args)-1], ".json") {
		t.Fatalf("expected temp .json file last, got %q", call.Args[len(call.Args)-1])
	}
}

func TestFormatPRChecks(t *testing.T) {
	mockBase := commands.NewMockBaseCommand()
	app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
	mockBase.SetExecCommandResult("SUCCESS\tbuild", "", nil)

	out, err := app.FormatPRChecks(`[{"name":"build","state":"SUCCESS"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "SUCCESS\tbuild" {
		t.Fatalf("unexpected output: %q", out)
	}
	call := mockBase.GetLastExecCommandCall()
	if !strings.Contains(call.Args[1], "No checks") {
		t.Fatalf("expected prChecksFilter, got %q", call.Args[1])
	}
}

func TestFormatErrorSurfacing(t *testing.T) {
	mockBase := commands.NewMockBaseCommand()
	app := &Jq{Cmd: commands.NewMockCommand(), Base: mockBase}
	mockBase.SetExecCommandResult("", "jq: syntax error", fmt.Errorf("exit 3"))

	if _, err := app.FormatPRView("{}"); err == nil {
		t.Fatal("expected error")
	} else if !strings.Contains(err.Error(), "jq: syntax error") {
		t.Fatalf("expected stderr surfaced, got: %v", err)
	}
}
