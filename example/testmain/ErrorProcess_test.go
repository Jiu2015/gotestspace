package testmain

import (
	"context"
	"testing"

	testspace "github.com/Jiu2015/gotestspace"
	"github.com/stretchr/testify/assert"
)

// There is non-exist command, then use testspace.Error for get the error details
func TestCommandWithGoTestSpaceError(t *testing.T) {
	var testSpaceError testspace.Error
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, _, err := myTestSpace.Execute(cancelCtx, "gitnotexit log --oneline")
	if assert.ErrorAs(t, err, &testSpaceError, "failed while execute") {
		assert.Equal(t, 127, testSpaceError.GetExitCode(), "the exit code invalid")
		assert.Equal(t, "command cannot found", testSpaceError.GetMessage())
		assert.Equal(t, "/bin/bash: line 13: gitnotexit: command not found\n", testSpaceError.GetStderr())
	}
}
