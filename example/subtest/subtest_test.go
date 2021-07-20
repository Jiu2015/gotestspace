package subtest_example

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	testspace "gitlab.alibaba-inc.com/agit/gotestspace"
)

var myTestSpace testspace.Space

func initTestSpace(t *testing.T) {
	var err error

	myTestSpace, err = testspace.Create(
		testspace.WithPathOption("testspace-*"),
		testspace.WithShellOption(`
			git config --global core.abbrev 10 &&
			git config --global init.defaultBranch master &&
			git init --bare repo.git &&
			git clone repo.git workdir &&
			(
				cd workdir &&
				printf "A\n" >A &&
				git add A &&
				test_tick &&
				git commit -m "A" &&
				printf "B\n" >B &&
				git add B &&
				test_tick &&
				git commit -m "B" &&
				git push -u origin HEAD
			)
		`),
	)
	assert.Nil(t, err)
}
func subTestGitRevParse(t *testing.T) {
	cmd := exec.Command("git", "rev-parse", "master")
	cmd.Dir = myTestSpace.GetPath("repo.git")
	actual, err := cmd.Output()
	assert.Nil(t, err)
	expect := "b475dff771ecf2c3d9f8baca56d436492cf915bc\n"
	assert.Equal(t, expect, string(actual))
}

// subTestGitLog will execute git command using current environment,
// and show abbrev commit in 7 digits.
func subTestGitLog(t *testing.T) {
	cmd := exec.Command("git", "log", "--oneline", "master", "--")
	cmd.Dir = myTestSpace.GetPath("workdir")
	actual, err := cmd.Output()
	assert.Nil(t, err)
	expect := strings.Join([]string{
		"b475dff B",
		"fe77af8 A",
		"",
	}, "\n")
	assert.Equal(t, expect, string(actual))
}

// subTestGitLogInTestSpaceEnv will run script in workspace environment,
// which has core.abbrev settings in global git config settings.
// Will show abbrev commit ID in 10 digits.
func subTestGitLogInTestSpaceEnv(t *testing.T) {
	stdout, _, err := myTestSpace.Execute("git -C workdir log --oneline master --")
	assert.Nil(t, err)
	expect := strings.Join([]string{
		"b475dff771 B",
		"fe77af8910 A",
		"",
	}, "\n")
	assert.Equal(t, expect, stdout)
}

func TestWithSubtest(t *testing.T) {
	// Setup
	initTestSpace(t)

	t.Run("T=1", subTestGitRevParse)
	t.Run("T=2", subTestGitLog)
	t.Run("T=3", subTestGitLogInTestSpaceEnv)

	// Tear-down
	myTestSpace.Cleanup()
}
