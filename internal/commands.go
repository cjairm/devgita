package commands

import (
	"bufio"
	"os"
	"os/exec"

	"github.com/cjairm/devgita/pkg/utils"
)

type CommandParams struct {
	PreExecMsg  string
	PostExecMsg string
	IsSudo      bool
	Verbose     bool
	Command     string
	Args        []string
}

func ExecCommand(cmd CommandParams) error {
	if cmd.PreExecMsg != "" {
		utils.Print(cmd.PreExecMsg, "")
	}
	command := cmd.Command
	if cmd.IsSudo {
		command = "sudo " + command
	}
	execCommand := exec.Command(command, cmd.Args...)

	// Create pipes to capture standard output
	stdout, err := execCommand.StdoutPipe()
	if err != nil {
		return err
	}

	// Create pipes to capture standard error
	stderr, err := execCommand.StderrPipe()
	if err != nil {
		return err
	}

	// To allow user to interact with the command (password)
	execCommand.Stdin = os.Stdin

	// Start the command
	if err := execCommand.Start(); err != nil {
		return err
	}

	if cmd.Verbose == true {
		// Create a scanner for stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				utils.PrintSecondary(scanner.Text())
			}
		}()

		// Create a scanner for stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				utils.PrintError(scanner.Text())
			}
		}()
	}

	err = execCommand.Wait()
	if cmd.PostExecMsg != "" && err == nil {
		utils.Print(cmd.PostExecMsg, "")
	}
	return err
}

// TODO: Rewrite this to have cross functional code
func AddLineToFile(line, filePath string) error {
	execCommand := CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "sh",
		Args: []string{
			"-c",
			"echo \"" + line + "\" >> " + filePath,
		},
	}
	return ExecCommand(execCommand)
}
