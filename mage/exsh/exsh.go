package exsh

import (
	"os"
	"os/exec"
)

func IsCmdAvail(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func Command(cs string, args ...string) *exec.Cmd {
	cmd := exec.Command(cs, args...)
	cmd.Stdout = os.Stdout
	return cmd
}
