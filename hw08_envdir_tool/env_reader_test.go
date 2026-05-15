package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read valid env directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "envdir_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		files := map[string]string{
			"VAR1":  "value1",
			"VAR2":  "value2  \t",
			"VAR3":  "line1\nline2",
			"VAR4":  "val\x00with\x00nulls",
			"EMPTY": "",
		}

		for name, content := range files {
			err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o644)
			require.NoError(t, err)
		}

		// File with = should be ignored
		err = os.WriteFile(filepath.Join(tmpDir, "INVALID=VAR"), []byte("val"), 0o644)
		require.NoError(t, err)

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		require.Equal(t, EnvValue{Value: "value1", NeedRemove: false}, env["VAR1"])
		require.Equal(t, EnvValue{Value: "value2", NeedRemove: false}, env["VAR2"])           // trimmed
		require.Equal(t, EnvValue{Value: "line1", NeedRemove: false}, env["VAR3"])            // first line
		require.Equal(t, EnvValue{Value: "val\nwith\nnulls", NeedRemove: false}, env["VAR4"]) // nulls replaced
		require.Equal(t, EnvValue{NeedRemove: true}, env["EMPTY"])
		_, exists := env["INVALID=VAR"]
		require.False(t, exists)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := ReadDir("/non/existent/path")
		require.Error(t, err)
	})
}
