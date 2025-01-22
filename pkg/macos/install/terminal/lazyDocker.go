// -------------------------
// NOTE: Write documentation how to use this
// - Some useful commands: https://www.youtube.com/watch?v=lu8edvTDUvI
// - Full documentation: https://github.com/jesseduffield/lazydocker
// -------------------------

package macos

import (
	"github.com/cjairm/devgita/pkg/common"
)

func InstallLazyDocker() error {
	if common.IsCommandInstalled("lazydocker") {
		return upgradeLazyDocker()
	}
	return installLazyDocker()
}

// installLazyDocker installs LazyDocker using Homebrew.
func installLazyDocker() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing LazyDocker",
		PostExecutionMessage: "Lazydocker installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"install",
			"jesseduffield/lazydocker/lazydocker",
		},
	}
	return common.ExecCommand(cmd)
}

// upgradeLazyDocker upgrades LazyDocker using Homebrew.
func upgradeLazyDocker() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Upgrading LazyDocker",
		PostExecutionMessage: "Lazydocker upgraded ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"upgrade",
			"jesseduffield/lazydocker/lazydocker",
		},
	}
	return common.ExecCommand(cmd)
}
