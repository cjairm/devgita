// -------------------------
// NOTE: Write documentation how to use this
// - Some useful commands: https://www.youtube.com/watch?v=CPLdltN7wgE
// - Full documentation: https://github.com/jesseduffield/lazygit
// -------------------------

package macos

import (
	"github.com/cjairm/devgita/pkg/common"
)

func InstallLazyGit() error {
	if common.IsCommandInstalled("lazygit") {
		return upgradeLazygit()
	}
	return installLazygit()
}

// installLazygit installs Lazygit using Homebrew.
func installLazygit() error {
	return common.ExecCommand(
		"Installing Lazygit",
		"Lazygit installed ✔",
		"brew",
		"install",
		"lazygit",
	)
}

// upgradeLazygit upgrades Lazygit using Homebrew.
func upgradeLazygit() error {
	return common.ExecCommand(
		"Upgrading Lazygit",
		"Lazygit upgraded ✔",
		"brew",
		"upgrade",
		"lazygit",
	)
}
