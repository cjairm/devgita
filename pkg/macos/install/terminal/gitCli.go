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
	return common.ExecCommand(
		"Installing GitHub CLI",
		"GitHub CLI installed ✔",
		"brew",
		"install",
		"gh",
	)
}

// upgradeGitHubCLI upgrades GitHub CLI using Homebrew.
func upgradeGitHubCLI() error {
	return common.ExecCommand(
		"Upgrading GitHub CLI",
		"GitHub CLI upgraded ✔",
		"brew",
		"upgrade",
		"gh",
	)
}
