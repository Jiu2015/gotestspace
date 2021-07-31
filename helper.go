package testspace

import (
	"context"
	"io/ioutil"
	"os/exec"
)

// SimpleExecuteCommand the execute cmd sample method
func SimpleExecuteCommand(ctx context.Context, path string, env []string, commandName string,
	args ...string) (output string, outErr string, err error) {
	cmd := exec.Command(commandName, args...)

	if len(path) > 0 {
		cmd.Dir = path
	}

	spaceCommand, err := new(ctx, cmd, nil, nil, nil, env...)
	if err != nil {
		return "", "", err
	}

	cmdStdout, err := ioutil.ReadAll(spaceCommand)
	if err != nil {
		return "", "", err
	}
	output = string(cmdStdout)

	if spaceCommand.stderr != nil {
		outErr = string(spaceCommand.stderr.GetStderr())
	}

	if err = cmd.Wait(); err != nil {
		return "", "", err
	}

	return output, outErr, err
}
