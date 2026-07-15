package cmd

import "testing"

func TestIsInteractiveTerminal_FalseUnderTestRunner(t *testing.T) {
	// `go test` never attaches a real TTY to stdout, so this must be false —
	// this is also exactly the condition dg list relies on to decide
	// whether to fall back to plain output in CI.
	if isInteractiveTerminal() {
		t.Error("expected isInteractiveTerminal() to be false under the test runner")
	}
}
