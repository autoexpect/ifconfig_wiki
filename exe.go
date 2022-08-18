package main

import (
	"bytes"
	"os/exec"
)

// Exe ...
func Exe(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	return stdout.String()
}
