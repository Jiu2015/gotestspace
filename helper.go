package testspace

import (
	"context"
	"io"
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

// NewTestSpaceCommand will return space-command for advantage use,
// You must get stdin, stdout and stderr before spaceCommand.Wait(), and do not miss spaceCommand.Wait()
func NewTestSpaceCommand(ctx context.Context, path string, env []string, enableStdin bool, stdout, stderr io.Writer,
	commandName string, args ...string) (*command, error) {
	var tempStdin io.Reader
	cmd := exec.Command(commandName, args...)

	if len(path) > 0 {
		cmd.Dir = path
	}

	if enableStdin {
		tempStdin = setStdinType
	}

	spaceCommand, err := new(ctx, cmd, tempStdin, stdout, stderr, env...)
	if err != nil {
		return nil, err
	}

	return spaceCommand, nil
}
