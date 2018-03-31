package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("make", os.Args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Unfortunately, Go has not yet standardized the API for
		// querying a process's exit status. So we reduce failing
		// statuses to wau for now.
		os.Exit(1)
	}
}
