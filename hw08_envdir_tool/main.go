package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s /path/to/env/dir command [args...]\n", os.Args[0])
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading environment directory: %s\n", err)
		os.Exit(1)
	}

	code := RunCmd(command, env)
	os.Exit(code)
}
