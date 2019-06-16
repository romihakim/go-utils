package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Exec — Execute an external program
// returnVar, 0: succ; 1: fail
// Return the last line from the result of the command.
func Exec(command string, output *[]string, returnVar *int) string {
	r, _ := regexp.Compile(`[ ]+`)
	parts := r.Split(command, -1)

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	cmd := exec.Command(parts[0], args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		*returnVar = 1
		return ""
	}

	*returnVar = 0
	*output = strings.Split(strings.TrimRight(string(out), "\n"), "\n")

	if l := len(*output); l > 0 {
		return (*output)[l-1]
	}

	return ""
}

// System — Execute an external program and display the output
// returnVar, 0: succ; 1: fail
// Returns the last line of the command output on success, and "" on failure.
func System(command string, returnVar *int) string {
	*returnVar = 0

	var stdBuf bytes.Buffer
	var err, err1, err2, err3 error

	r, _ := regexp.Compile(`[ ]+`)
	parts := r.Split(command, -1)

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	cmd := exec.Command(parts[0], args...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := io.MultiWriter(os.Stdout, &stdBuf)
	stderr := io.MultiWriter(os.Stderr, &stdBuf)

	err = cmd.Start()
	if err != nil {
		*returnVar = 1
		return ""
	}

	go func() {
		_, err1 = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, err2 = io.Copy(stderr, stderrIn)
	}()

	err3 = cmd.Wait()
	if err1 != nil || err2 != nil || err3 != nil {
		if err1 != nil {
			fmt.Println(err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		if err3 != nil {
			fmt.Println(err3)
		}
		*returnVar = 1
		return ""
	}

	if output := strings.TrimRight(stdBuf.String(), "\n"); output != "" {
		pos := strings.LastIndex(output, "\n")
		if pos == -1 {
			return output
		}
		return output[pos+1:]
	}

	return ""
}

// Passthru — Execute an external program and display raw output
// returnVar, 0: succ; 1: fail
func Passthru(command string, returnVar *int) {
	r, _ := regexp.Compile(`[ ]+`)
	parts := r.Split(command, -1)

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	cmd := exec.Command(parts[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		*returnVar = 1
		fmt.Println(err)
	} else {
		*returnVar = 0
	}
}
