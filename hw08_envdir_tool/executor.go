package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// Let's use a map for easier manipulation
	resEnv := make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			resEnv[kv[0]] = kv[1]
		}
	}

	for k, v := range env {
		if v.NeedRemove {
			delete(resEnv, k)
		} else {
			resEnv[k] = v.Value
		}
	}

	newEnv := make([]string, 0, len(resEnv))
	for k, v := range resEnv {
		newEnv = append(newEnv, k+"="+v)
	}
	c.Env = newEnv

	err := c.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}
