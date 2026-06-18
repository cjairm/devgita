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
