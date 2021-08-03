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

	cancelCtx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	// Read set env
	output, _, err = SimpleExecuteCommand(cancelCtx2, "",
		[]string{"GOSHELLHELPERTEST=Just_for_test123"},
		"env")
	assert.NoError(err)
	assert.Contains(output, "GOSHELLHELPERTEST=Just_for_test123")

	cancelCtx3, cancel3 := context.WithCancel(context.Background())
	defer cancel3()
	_, _, err = SimpleExecuteCommand(cancelCtx3, "", nil, "evn1111")
	if err == nil {
		assert.Error(err)
	}
}
