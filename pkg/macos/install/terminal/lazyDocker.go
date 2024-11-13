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
	return common.ExecCommand(
		"Installing LazyDocker",
		"Lazydocker installed ✔",
		"brew",
		"install",
		"jesseduffield/lazydocker/lazydocker",
	)
}

// upgradeLazyDocker upgrades LazyDocker using Homebrew.
func upgradeLazyDocker() error {
	return common.ExecCommand(
		"Upgrading LazyDocker",
		"Lazydocker upgraded ✔",
		"brew",
		"upgrade",
		"jesseduffield/lazydocker/lazydocker",
	)
}
