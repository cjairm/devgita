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
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Lazygit",
		PostExecutionMessage: "Lazygit installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"install",
			"lazygit",
		},
	}
	return common.ExecCommand(cmd)
}

// upgradeLazygit upgrades Lazygit using Homebrew.
func upgradeLazygit() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Upgrading Lazygit",
		PostExecutionMessage: "Lazygit upgraded ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"upgrade",
			"lazygit",
		},
	}
	return common.ExecCommand(cmd)
}
