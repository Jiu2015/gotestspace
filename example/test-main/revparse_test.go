package test_main_example

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRevParse(t *testing.T) {
	cmd := exec.Command("git", "rev-parse", "master")
	cmd.Dir = myTestSpace.GetPath("repo.git")
	actual, err := cmd.Output()
	assert.Nil(t, err)
	expect := "b475dff771ecf2c3d9f8baca56d436492cf915bc\n"
	assert.Equal(t, expect, string(actual))
}
