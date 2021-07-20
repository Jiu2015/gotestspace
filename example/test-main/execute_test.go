package test_main_example

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGitLogInTestSpaceEnv will run script in workspace environment,
// which has core.abbrev settings in global git config settings.
// Will show abbrev commit ID in 10 digits.
func TestGitLogInTestSpaceEnv(t *testing.T) {
	stdout, _, err := myTestSpace.Execute("git -C workdir log --oneline master --")
	assert.Nil(t, err)
	expect := strings.Join([]string{
		"b475dff771 B",
		"fe77af8910 A",
		"",
	}, "\n")
	assert.Equal(t, expect, stdout)
}
