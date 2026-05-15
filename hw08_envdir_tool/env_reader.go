package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			continue
		}

		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		info, err := f.Stat()
		if err != nil {
			f.Close()
			continue
		}

		if info.Size() == 0 {
			env[name] = EnvValue{NeedRemove: true}
			f.Close()
			continue
		}

		content, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}

		// Use the first line
		lines := bytes.SplitN(content, []byte("\n"), 2)
		valBytes := lines[0]

		// Replace null characters with \n
		valBytes = bytes.ReplaceAll(valBytes, []byte{0x00}, []byte("\n"))

		// Trim right spaces and tabs
		val := strings.TrimRight(string(valBytes), " \t")

		env[name] = EnvValue{Value: val, NeedRemove: false}
	}

	return env, nil
}
