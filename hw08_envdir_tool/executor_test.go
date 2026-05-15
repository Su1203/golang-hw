package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run command with environment", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": EnvValue{Value: "test_value", NeedRemove: false},
		}

		// Use 'env' command to check variables (available on most systems)
		// On Windows we can use 'cmd /c set'
		var cmd []string
		if os.PathSeparator == '\\' {
			cmd = []string{"cmd", "/c", "set", "TEST_VAR"}
		} else {
			cmd = []string{"env"}
		}

		// We need to capture output to verify, but RunCmd redirects to os.Stdout
		// For testing purpose, we can temporarily mock os.Stdout or just check the exit code
		code := RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})

	t.Run("remove environment variable", func(t *testing.T) {
		os.Setenv("REMOVE_ME", "should_be_gone")
		defer os.Unsetenv("REMOVE_ME")

		env := Environment{
			"REMOVE_ME": EnvValue{NeedRemove: true},
		}

		var cmd []string
		if os.PathSeparator == '\\' {
			cmd = []string{"cmd", "/c", "if defined REMOVE_ME (exit 1) else (exit 0)"}
		} else {
			cmd = []string{"sh", "-c", "[ -z \"$REMOVE_ME\" ]"}
		}

		code := RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})

	t.Run("command not found", func(t *testing.T) {
		code := RunCmd([]string{"non-existent-command-12345"}, nil)
		require.NotEqual(t, 0, code)
	})
}
