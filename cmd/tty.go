/*
 * Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import "os"

// isInteractiveTerminal reports whether stdout is attached to a real terminal.
// Used by dg list / dg validate to decide between the interactive dashboard and
// plain-text output (piped, redirected, or CI contexts always get plain text).
func isInteractiveTerminal() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
