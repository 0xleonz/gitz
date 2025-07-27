package utils

import (
	"os"
	"os/exec"
)

// ejecuta otro subcomando de gitz (ie,"message", "push").
func Call(args ...string) error {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
