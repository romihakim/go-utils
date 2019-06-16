package utils

import (
	"testing"
)

func TestProgramExecution(t *testing.T) {
	var output []string
	var retVal int

	tExec := Exec("ls -l", &output, &retVal)
	gt(t, float64(len(tExec)), 0)
	equal(t, 0, retVal)

	tSystem := System("ls -l", &retVal)
	equal(t, 0, retVal)
	gt(t, float64(len(tSystem)), 0)

	Passthru("echo hello", &retVal)
	equal(t, 0, retVal)
}
