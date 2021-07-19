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
	output, _, err := ExecuteCommand(context.Background(), "", nil, "env")
	assert.NoError(err)
	assert.Contains(output, "GOSHELLHELPERTEST=Just_for_test")

	// Read set env
	output, _, err = ExecuteCommand(context.Background(), "", []string{"GOSHELLHELPERTEST=Just_for_test123"}, "env")
	assert.NoError(err)
	assert.Contains(output, "GOSHELLHELPERTEST=Just_for_test123")

	_, _, err = ExecuteCommand(context.Background(), "", nil, "evn1111")
	if err == nil {
		assert.Error(err)
	}
}
