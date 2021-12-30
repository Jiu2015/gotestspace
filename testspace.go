package testspace

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CustomCleaner the customer cleaners for registration
type CustomCleaner func() error

// Space the repo action interface, contains repo action methods
type Space interface {
	Cleanup() error
	GetPath(subDirName string) string
	GetMultiPath(subDirNames ...string) string
	GetEnvStr() []string
	GetTemplateStr() string
	GetShellStr() string
	GetOutputStr() string
	GetOutErr() string

	// Execute the command.
	// The error may be testspace.Error type, so please use type assertions or use errors.As
	Execute(ctx context.Context, shell string) (stdout string, stderr string, _ error)

	// ExecuteWithStdin Will enable stdin on the command, you can do a lot of advanced things.
	// WARNING: You must call command.Wait() method after you operate command!
	ExecuteWithStdin(ctx context.Context, shell string) (*command, error)

	// RegistrationCustomCleaner used for registration the cleaners func for clean other resources
	// while the testing finished
	RegistrationCustomCleaner(cleaners ...CustomCleaner)
}

// workSpace the repo struct
type workSpace struct {
	path        string
	env         []string
	template    string
	customShell string
	output      string
	outErr      string
	cleaners    []CustomCleaner
}

// Cleanup destroy the workspace path
func (w *workSpace) Cleanup() error {
	if len(w.path) == 0 {
		return errors.New("the workspace path invalid, please check and delete it manually")
	}

	if err := os.RemoveAll(w.path); err != nil {
		return err
	}

	for _, c := range w.cleaners {
		if err := c(); err != nil {
			return err
		}
	}

	return nil
}

// GetPath get current workspace path
func (w *workSpace) GetPath(subDirName string) string {
	subDirName = filepath.Clean(subDirName)
	if strings.HasPrefix(subDirName, "../") {
		subDirName = ""
	}

	return path.Join(w.path, subDirName)
}

// GetMultiPath get the multiple path with path.Join
func (w *workSpace) GetMultiPath(subDirNames ...string) string {
	tmpPathNames := make([]string, 0, len(subDirNames)+1)
	tmpPathNames = append(tmpPathNames, w.path)

	for _, s := range subDirNames {
		if !strings.HasPrefix(s, "../") {
			tmpPathNames = append(tmpPathNames, s)
		}
	}

	return path.Join(tmpPathNames...)
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

// Execute the command
// The error may be testspace.Error type, so please use type assertions or use errors.As
func (w *workSpace) Execute(ctx context.Context, shell string) (stdout string, stderr string, _ error) {
	mixedShell := w.template + "\n" + shell
	output, outErr, err := SimpleExecuteCommand(ctx,
		w.path, w.env, "/bin/bash", "-c", mixedShell)

	return output, outErr, err
}

// ExecuteWithStdin execute shell with stdin
func (w *workSpace) ExecuteWithStdin(ctx context.Context, shell string) (*command, error) {
	mixedShell := w.template + "\n" + shell
	return NewTestSpaceCommand(ctx, w.path, w.env, true, nil, nil,
		"/bin/bash", "-c", mixedShell)
}

// RegistrationCustomCleaner register the custom cleaners for clean
func (w *workSpace) RegistrationCustomCleaner(cleaners ...CustomCleaner) {
	w.cleaners = append(w.cleaners, cleaners...)
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

	if err := initGitWorkspace(currentOption.workspacePath); err != nil {
		return nil, err
	}

	space := &workSpace{
		path:        currentOption.workspacePath,
		env:         currentOption.environments,
		template:    currentOption.template,
		customShell: currentOption.customShell,
		cleaners:    currentOption.cleaners,
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, fn, _, ok := runtime.Caller(1); ok {
		space.env = append(space.env, "CALLER="+fn)
		space.env = append(space.env, "CALLER_DIR="+filepath.Dir(fn))
	}

	if _, stderr, err := space.Execute(cancelCtx, currentOption.customShell); err != nil {
		// If the command got error, then cleanup the temporary folder
		space.Cleanup()
		return space, fmt.Errorf("err: %s, stderr: %s", err, stderr)
	}

	return space, nil
}

// Init the git workspace, to prevent source code edit by mistake
func initGitWorkspace(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, stderr, err := SimpleExecuteCommand(ctx, path, nil, "git", "init", ".")
	if err != nil {
		return fmt.Errorf("err: %s, stderr: %s", err, stderr)
	}

	return nil
}
