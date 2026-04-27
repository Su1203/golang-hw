package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: go-envdir /path/to/env/dir command arg1 arg2")
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
