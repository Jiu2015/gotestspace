package testspace

import "fmt"

const (
	// GenerateError the generate error, the exit code is 1
	GenerateError = 1

	// Misuse the misuse command, the exit code is 2
	Misuse = 2

	// NotExecute the not execute command error, the exit code is 126
	NotExecute = 126

	// NotFound the command not found error, the exit code is 127
	NotFound = 127

	// InvalidArgument the command has invalid argument, the exit code is 128
	InvalidArgument = 128
)

// Error gotestspace custom error
type Error struct {
	// The exit code with command
	exitCode int

	// The exit with signal, if exit but not through system singal, it will be 0
	signal int

	// The original error
	err error

	// The standard output
	stdout string

	// The standard error output
	stderr string
}

// Error print the generate error
func (e Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return ""
}

// Unwrap get the origin error
func (e Error) Unwrap() error {
	return e.err
}

// GetMessage get the message from custom error
// it will print the error from exitcode and signal
func (e *Error) GetMessage() string {
	if e.err == nil {
		return ""
	}

	if e.signal > 0 {
		return fmt.Sprintf("command interrupt by singal: %d", e.signal)
	}

	switch e.exitCode {
	case GenerateError:
		return "catch shell generate error"
	case Misuse:
		return "misuse of shell build-in"
	case NotExecute:
		return "command cannot execute"
	case NotFound:
		return "command cannot found"
	case InvalidArgument:
		return "invalid argument"
	default:
		return "unknown error"
	}
}

// GetStdout get the standard output
func (e *Error) GetStdout() string {
	return e.stdout
}

// GetStderr get the standard error output
func (e *Error) GetStderr() string {
	return e.stderr
}

// GetExitCode get the exit code
func (e *Error) GetExitCode() int {
	return e.exitCode
}

// NewGoTestSpaceError create new gotestspace error
func NewGoTestSpaceError(exitCode int, stdout, stderr string, err error) Error {
	signal := splitCodeWithSignal(exitCode)

	return Error{
		exitCode: exitCode,
		signal:   signal,
		stdout:   stdout,
		stderr:   stderr,
		err:      err,
	}
}

// splitCodeWithSignal used to check exit code and get the signal number
func splitCodeWithSignal(exitCode int) (signal int) {
	if exitCode > 128 && exitCode < 255 {
		return exitCode - 128
	}

	return 0
}
