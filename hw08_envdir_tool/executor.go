package main

import (
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for eVar, eVal := range env {
		if eVal.NeedRemove {
			os.Unsetenv(eVar)
			continue
		}
		os.Setenv(eVar, eVal.Value)
	}
	proc := exec.Command(cmd[0], cmd[1:]...)
	setSTD(proc, os.Stdout, os.Stderr, os.Stdin)

	if err := proc.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return
}

func setSTD(proc *exec.Cmd, out io.Writer, err io.Writer, rd io.Reader) {
	proc.Stdout = out
	proc.Stderr = err
	proc.Stdin = rd
}
