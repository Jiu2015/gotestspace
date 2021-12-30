package testspace

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
)

// SimpleExecuteCommand to execute cmd sample method.
// The error may be testspace.Error type, so please use type assertions or use errors.As
func SimpleExecuteCommand(ctx context.Context, path string, env []string, commandName string,
	args ...string) (output string, outErr string, err error) {
	cmd := exec.Command(commandName, args...)

	if len(path) > 0 {
		cmd.Dir = path
	}

	stderr := &bytes.Buffer{}

	spaceCommand, err := new(ctx, cmd, nil, nil, stderr, env...)
	if err != nil {
		return "", stderr.String(), err
	}

	cmdStdout, err := ioutil.ReadAll(spaceCommand)
	if err != nil {
		return "", stderr.String(), err
	}

	output = string(cmdStdout)
	outErr = spaceCommand.GetStderr()

	if err = cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitError.ExitCode()
			// Wrap error to testspace.Error
			err = NewGoTestSpaceError(exitError.ExitCode(), output, stderr.String(), err)
		}

		return output, stderr.String(), err
	}

	return output, outErr, err
}

// NewTestSpaceCommand will return space-command for advantage use,
// You must get stdin, stdout and stderr before spaceCommand.Wait(), and do not miss spaceCommand.Wait()
func NewTestSpaceCommand(ctx context.Context, path string, env []string, enableStdin bool, stdout, stderr io.Writer,
	commandName string, args ...string) (Commander, error) {
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
