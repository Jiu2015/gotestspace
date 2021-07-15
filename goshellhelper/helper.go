package goshellhelper

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
)

// ExecuteCommand the execute cmd sample method
func ExecuteCommand(ctx context.Context, path string, env []string, command string, args ...string) (output string, err error) {
	var (
		stdout,
		stderr bytes.Buffer
	)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = path
	if len(env) != 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	err = cmd.Run()
	if err != nil {
		return "", err
	}
	output = string(stdout.Bytes())
	errStr := string(stderr.Bytes())
	if len(errStr) > 0 {
		err = errors.New(errStr)
	}
	return
}
