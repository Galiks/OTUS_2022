package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

var ErrEmptyCmd = errors.New("commands is empty")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Println(ErrEmptyCmd)
		return 1
	}
	for eVar, eVal := range env {
		if eVal.NeedRemove {
			os.Unsetenv(eVar)
			continue
		}
		os.Setenv(eVar, eVal.Value)
	}
	commandName := cmd[0]
	args := cmd[1:]
	proc := exec.Command(commandName, args...)
	setSTD(proc, os.Stdout, os.Stderr, os.Stdin)

	if err := proc.Run(); err != nil {
		errExitCode := &exec.ExitError{}
		if errors.As(err, &errExitCode) {
			return errExitCode.ExitCode()
		}
	}

	return
}

func setSTD(proc *exec.Cmd, out io.Writer, err io.Writer, rd io.Reader) {
	proc.Stdout = out
	proc.Stderr = err
	proc.Stdin = rd
}
