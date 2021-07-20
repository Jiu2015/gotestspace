package testspace

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Space the repo action interface, contains repo action methods
type Space interface {
	Cleanup() error
	GetPath(subDirName string) string
	GetEnvStr() []string
	GetTemplateStr() string
	GetShellStr() string
	GetOutputStr() string
	GetOutErr() string
	Execute(shell string) (stdout string, stderr string, _ error)
}

// workSpace the repo struct
type workSpace struct {
	path        string
	env         []string
	template    string
	customShell string
	output      string
	outErr      string
}

// Cleanup destroy the workspace path
func (w *workSpace) Cleanup() error {
	if len(w.path) == 0 {
		return errors.New("the workspace path invalid, please check and delete it manually")
	}

	return os.RemoveAll(w.path)
}

// GetPath get current workspace path
func (w *workSpace) GetPath(subDirName string) string {
	subDirName = filepath.Clean(subDirName)
	if strings.HasPrefix(subDirName, "../") {
		subDirName = ""
	}
	return path.Join(w.path, subDirName)
}

// GetEnvStr get current environments string
func (w *workSpace) GetEnvStr() []string {
	return w.env
}

// GetTemplateStr get template string
func (w *workSpace) GetTemplateStr() string {
	return w.template
}

// GetShellStr get the shell which has been run
func (w *workSpace) GetShellStr() string {
	return w.customShell
}

// GetOutputStr get the shell output string
func (w *workSpace) GetOutputStr() string {
	return w.output
}

// GetOutErr get the error print
func (w *workSpace) GetOutErr() string {
	return w.outErr
}

func (w *workSpace) Execute(shell string) (stdout string, stderr string, _ error) {
	mixedShell := w.template + "\n" + shell
	output, outErr, err := ExecuteCommand(context.Background(),
		w.path, w.env, "/bin/bash", "-c", mixedShell)
	if err != nil {
		w.output = output
		w.outErr = outErr
		return output, outErr, err
	}

	w.output = output
	w.outErr = outErr

	return output, outErr, nil
}

// Create create repo object
func Create(options ...CreateOption) (Space, error) {
	currentOption := mergeOptions(options)

	// Check the dir is or not exist
	if _, err := os.Stat(currentOption.workspacePath); err != nil {
		// Create the workspace directory
		err = os.MkdirAll(currentOption.workspacePath, 0755)
		if err != nil {
			return nil, err
		}
	}

	initGitWorkspace(currentOption.workspacePath)

	space := &workSpace{
		path:        currentOption.workspacePath,
		env:         currentOption.environments,
		template:    currentOption.template,
		customShell: currentOption.customShell,
	}

	if _, _, err := space.Execute(currentOption.customShell); err != nil {
		// If the command got error, then cleanup the temporary folder
		space.Cleanup()
		return space, err
	}

	return space, nil
}

// Init the git workspace, to prevent source code edit by mistake
func initGitWorkspace(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _, err := ExecuteCommand(ctx, path, nil, "git", "init", ".")
	if err != nil {
		return err
	}

	return nil
}