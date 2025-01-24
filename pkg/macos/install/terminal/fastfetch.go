// -------------------------
// NOTE: Write documentation or create icon to open and get information of this Mac
// - Documentation: https://github.com/fastfetch-cli/fastfetch
// -------------------------

package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
)

func InstallFastFetch(devgitaPath string) error {
	if err := common.InstallOrUpdateBrewPackage("fastfetch"); err != nil {
		return err
	}
	if err := configureFastFetch(devgitaPath); err != nil {
		return fmt.Errorf("Error copying fastfetch config: %v", err)
	}
	return nil

}

func configureFastFetch(devgitaPath string) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "fastfetch")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"configs",
		"fastfetch",
	)
	if err := common.MoveContents(devgitaConfig, configDir); err != nil {
		return fmt.Errorf("error setting up fastfetch: %w", err)
	}
	return nil
}
