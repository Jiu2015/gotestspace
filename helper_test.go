package testspace

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommand(t *testing.T) {
	assert := assert.New(t)

	// Read system env
	os.Setenv("GOSHELLHELPERTEST", "Just_for_test")
	defer os.Unsetenv("GOSHELLHELPERTEST")

	cancelCtx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	output, _, err := SimpleExecuteCommand(cancelCtx1, "", nil, "env")
	assert.NoError(err)
	assert.Contains(output, "GOSHELLHELPERTEST=Just_for_test")

	output, _, err = SimpleExecuteCommand(cancelCtx1, "",
		[]string{"GOSHELLHELPERTEST=Just_for_test123"},
		"env")
	assert.NoError(err)
	assert.Contains(output, "GOSHELLHELPERTEST=Just_for_test123")

	_, _, err = SimpleExecuteCommand(cancelCtx1, "", nil, "evn1111")
	if err == nil {
		assert.Error(err)
	}

	_, stderr, err := SimpleExecuteCommand(cancelCtx1, "", nil, "git", "aaa")
	assert.Error(err)
	assert.Contains(stderr, "git --help")
}
