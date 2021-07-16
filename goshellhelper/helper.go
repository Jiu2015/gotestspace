package goshellhelper

import (
	"bytes"
	"context"
	"os/exec"
)

// ExecuteCommand the execute cmd sample method
func ExecuteCommand(ctx context.Context, path string, env []string, command string, args ...string) (output string, outerr string, err error) {
	var (
		stdout,
		stderr bytes.Buffer
	)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = path
	if len(env) != 0 {
		cmd.Env = env
	}

	err = cmd.Run()
	if err != nil {
		return "", string(stderr.Bytes()), err
	}
	output = string(stdout.Bytes())
	outerr = string(stderr.Bytes())

	return
}
