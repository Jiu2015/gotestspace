package testmain

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGitLog will execute git command using current environment,
// and show abbrev commit in 7 digits.
func TestGitLog(t *testing.T) {
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
