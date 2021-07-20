package testspace

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// createOption the create workspace option
type createOption struct {
	workspacePath string
	template      string
	environments  []string
	customShell   string
}

// CreateOption the option for create shell workspace
type CreateOption struct {
	setOpt func(opt *createOption)
}

// WithPathOption set the path
func WithPathOption(workPath string) CreateOption {
	var (
		fullPath string
		err      error
	)
	workPath = filepath.Clean(workPath)
	basePath := filepath.Base(workPath)
	dirPath := filepath.Dir(workPath)
	if strings.Contains(basePath, "*") {
		if filepath.IsAbs(dirPath) {
			fullPath, err = ioutil.TempDir(dirPath, basePath)
		} else {
			fullPath, err = ioutil.TempDir("", workPath)
		}
	} else if !filepath.IsAbs(workPath) {
		fullPath, err = os.Getwd()
		if err == nil {
			fullPath = path.Join(fullPath, workPath)
		}
	} else {
		fullPath = workPath
	}
	if err != nil {
		panic(err)
	}

	return CreateOption{setOpt: func(opt *createOption) {
		opt.workspacePath = fullPath
	}}
}

// WithTemplateOption set custom template
func WithTemplateOption(customTemplate string) CreateOption {
	return CreateOption{setOpt: func(opt *createOption) {
		// The template will append the default template which init test_tick function
		opt.template = opt.template + "\n" + customTemplate
	}}
}

// WithEnvironmentsOption set environments
func WithEnvironmentsOption(environments ...string) CreateOption {
	return CreateOption{setOpt: func(opt *createOption) {
		opt.environments = append(opt.environments, environments...)
	}}
}

// WithShellOption set custom shell from user
func WithShellOption(customShell string) CreateOption {
	return CreateOption{setOpt: func(opt *createOption) {
		opt.customShell = customShell
	}}
}

func mergeOptions(options []CreateOption) *createOption {
	currentPath, _ := os.Getwd()
	tempWorkspacePath := path.Join(currentPath, time.Now().Format("20060102150405"))
	o := &createOption{
		// Default path
		workspacePath: tempWorkspacePath,
		template: `
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}
`,
		environments: []string{
			"GIT_AUTHOR_EMAIL=author@example.com",
			"GIT_AUTHOR_NAME='A U Thor'",
			"GIT_COMMITTER_EMAIL=committer@example.com",
			"GIT_COMMITTER_NAME='C O Mitter'",
		},
		customShell: "",
	}

	for _, opt := range options {
		opt.setOpt(o)
	}

	o.environments = append([]string{fmt.Sprintf("HOME=%s", o.workspacePath)},
		o.environments...)

	return o
}