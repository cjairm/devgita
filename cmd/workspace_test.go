/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestWorkspaceCmd_Metadata checks workspaceCmd's Use/Aliases/Args without
// ever calling RunE/Execute() — see worktree_test.go's comment above
// TestWorktreeCreateCmd_AIAndLayoutMutuallyExclusive for why this codebase
// avoids Command.Execute() in these tests (it would reach RunE and, here,
// tuiworktree.Run()'s real Bubble Tea program).
func TestWorkspaceCmd_Metadata(t *testing.T) {
	if workspaceCmd.Use != "ws" {
		t.Errorf("expected Use %q, got %q", "ws", workspaceCmd.Use)
	}

	found := false
	for _, alias := range workspaceCmd.Aliases {
		if alias == "workspace" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected aliases %v to contain %q", workspaceCmd.Aliases, "workspace")
	}

	if workspaceCmd.RunE == nil {
		t.Fatal("expected workspaceCmd.RunE to be set")
	}
}

// TestWorkspaceCmd_NoArgs pins Args to cobra.NoArgs by exercising its
// validation function directly (comparing func values isn't possible in Go),
// without calling RunE/Execute().
func TestWorkspaceCmd_NoArgs(t *testing.T) {
	if workspaceCmd.Args == nil {
		t.Fatal("expected workspaceCmd.Args to be set")
	}
	if err := workspaceCmd.Args(workspaceCmd, nil); err != nil {
		t.Errorf("expected no args to be valid, got: %v", err)
	}
	if err := workspaceCmd.Args(workspaceCmd, []string{"extra"}); err == nil {
		t.Error("expected an error for unexpected args, got nil")
	}
}

// TestWorkspaceCmd_RegisteredOnRoot confirms `dg ws` is wired into rootCmd.
func TestWorkspaceCmd_RegisteredOnRoot(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c == workspaceCmd {
			found = c
			break
		}
	}
	if found == nil {
		t.Fatal("expected workspaceCmd to be registered as a child of rootCmd")
	}
}
