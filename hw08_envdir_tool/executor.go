package main

import (
	"errors"
	"os"
	"os/exec"
)

var commandError *exec.ExitError

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(envName)
			continue
		}
		os.Setenv(envName, envValue.Value)
	}

	commandName := cmd[0]
	args := cmd[1:]
	command := exec.Command(commandName, args...)
	command.Env = os.Environ()
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if errors.As(err, &commandError) {
			return commandError.ExitCode()
		}
	}

	return 0
}
