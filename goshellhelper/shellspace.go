package goshellhelper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// ShellSpace the repo action interface, contains repo action methods
type ShellSpace interface {
	Cleanup() error
	GetPath(subDirName string) (string, error)
	GetEnvStr() []string
	GetTemplateStr() string
	GetShellStr() string
	GetOutputStr() string
}

// WorkSpace the repo struct
type WorkSpace struct {
	path        string
	env         []string
	template    string
	customShell string
	output      string
}

// Cleanup destroy the workspace path
func (w *WorkSpace) Cleanup() error {
	if len(w.path) == 0 {
		return errors.New("the workspace path invalid, please check and delete it manually")
	}

	return os.RemoveAll(w.path)
}

// GetPath get current workspace path
func (w *WorkSpace) GetPath(subDirName string) (string, error) {
	if strings.Contains(subDirName, ".") ||
		strings.Contains(subDirName, "&") ||
		strings.Contains(subDirName, "|") {
		return "", errors.New("the subDirName illegal")
	}

	return path.Join(w.path, subDirName), nil
}

// GetEnvStr get current environments string
func (w *WorkSpace) GetEnvStr() []string {
	return w.env
}

// GetTemplateStr get template string
func (w *WorkSpace) GetTemplateStr() string {
	return w.template
}

// GetShellStr get the shell which has been run
func (w *WorkSpace) GetShellStr() string {
	return w.customShell
}

// GetOutputStr get the shell output string
func (w *WorkSpace) GetOutputStr() string {
	return w.output
}

// NewShellSpace create repo object
func NewShellSpace(options ...CreateOption) (ShellSpace, error) {
	currentOption := mergeOptions(options)
	mixedShell := currentOption.template + currentOption.customShell

	// Check the dir is or not exist
	_, err := os.Stat(currentOption.workspacePath)
	if os.IsExist(err) {
		return nil, fmt.Errorf("the path had existed, path: %s", currentOption.workspacePath)
	}

	// Create the workspace directory
	err = os.MkdirAll(currentOption.workspacePath, 0755)
	if err != nil {
		return nil, err
	}

	initGitWorkspace(currentOption.workspacePath)

	output, err := ExecuteCommand(context.Background(), currentOption.workspacePath, currentOption.environments, "/bin/bash", "-c", mixedShell)
	if err != nil {
		return nil, err
	}

	return &WorkSpace{
		path:        currentOption.workspacePath,
		env:         currentOption.environments,
		template:    currentOption.template,
		customShell: currentOption.customShell,
		output:      output,
	}, nil

}

// Init the git workspace, to prevent source code edit by mistake
func initGitWorkspace(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := ExecuteCommand(ctx, path, nil, "git", "init", ".")
	if err != nil {
		return err
	}

	return nil
}