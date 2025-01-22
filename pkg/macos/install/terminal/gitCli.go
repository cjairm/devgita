// -------------------------
// NOTE: Write documentation how to use this
// - What's github cli? https://www.youtube.com/watch?v=uy_PEGgUF4U
// - Some useful commands: https://www.youtube.com/watch?v=in_H8MbiHpw
// - Full documentation: https://cli.github.com/
// - Install documentation: https://github.com/cli/cli?tab=readme-ov-file#installation
// -------------------------

package macos

import (
	"github.com/cjairm/devgita/pkg/common"
)

func InstallGitHubCli() error {
	if common.IsCommandInstalled("gh") {
		return upgradeGitHubCLI()
	}
	return installGitHubCLI()
}

// installGitHubCLI installs GitHub CLI using Homebrew.
func installGitHubCLI() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing GitHub CLI",
		PostExecutionMessage: "GitHub CLI installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"install",
			"gh",
		},
	}
	return common.ExecCommand(cmd)
}

// upgradeGitHubCLI upgrades GitHub CLI using Homebrew.
func upgradeGitHubCLI() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Upgrading GitHub CLI",
		PostExecutionMessage: "GitHub CLI upgraded ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"upgrade",
			"gh",
		},
	}
	return common.ExecCommand(cmd)
}
