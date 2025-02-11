package bash

import (
	cmd "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/files"
)

type Bash struct {
	Cmd cmd.Command
}

func NewBash() *Bash {
	osCmd := cmd.NewCommand()
	return &Bash{Cmd: osCmd}
}

func (b *Bash) CopyCustomConfig() error {
	err := files.MoveFromConfigsToLocalConfig([]string{"bash"}, []string{"devgita"})
	if err != nil {
		return err
	}
	return nil
}
