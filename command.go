package testspace

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"sync"
	"syscall"
)

type command struct {
	reader   io.Reader
	writer   io.WriteCloser
	stderr   *stdErr
	cmd      *exec.Cmd
	context  context.Context
	waitOnce sync.Once
	waitErr  error
}

func (c *command) Read(p []byte) (int, error) {
	if c.reader == nil {
		panic("command invalid with reader")
	}

	return c.reader.Read(p)
}

func (c *command) Write(p []byte) (int, error) {
	if c.writer == nil {
		panic("command invalid with writer")
	}

	return c.writer.Write(p)
}

func (c *command) wait() {
	if c.writer != nil {
		c.writer.Close()
	}

	if c.reader != nil {
		io.Copy(ioutil.Discard, c.reader)
	}

	c.waitErr = c.cmd.Wait()
}

func (c *command) Wait() error {
	c.waitOnce.Do(c.wait)

	return c.waitErr
}

// The command use to check to set stdin type
var setStdinType io.Reader = setStdin{}

type setStdin struct{}

func (setStdin) Read([]byte) (int, error) {
	return 0, errors.New("just used for tag if or not use command stdin")
}

type stdErr struct {
	errContents []byte
}

func (s *stdErr) Write(p []byte) (int, error) {
	if len(p) > 0 {
		s.errContents = append(s.errContents, p...)
		return len(p), nil
	}

	return 0, nil
}

func (s *stdErr) GetStderr() []byte {
	return s.errContents
}

func new(ctx context.Context, cmd *exec.Cmd, stdin io.Reader, stdout, stderr io.Writer, env ...string) (*command, error) {
	if ctx.Done() == nil {
		panic("the context must have Done() method, use WithCancel, WithTimeout etc.")
	}

	resultCommand := &command{
		cmd:     cmd,
		context: ctx,
		stderr:  &stdErr{},
	}

	if len(env) > 0 {
		cmd.Env = append(cmd.Env, env...)
	}

	if stdin == setStdinType {
		pip, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		resultCommand.writer = pip
	}

	if resultCommand.writer == nil && stdin != nil {
		cmd.Stdin = stdin
	}

	if stdout == nil {
		pip, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		resultCommand.reader = pip
	} else {
		cmd.Stdout = stdout
	}

	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = resultCommand.stderr
	}

	// Set the child process have a new process group
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		if process := cmd.Process; process != nil && process.Pid > 0 {
			// Kill the the process and it's child process
			syscall.Kill(-process.Pid, syscall.SIGKILL)
		}
	}()

	return resultCommand, nil
}
