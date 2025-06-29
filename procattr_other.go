//go:build !windows
// +build !windows

package main

import "os/exec"

// SetNoWindow is a no-op on non-Windows systems
func SetNoWindow(cmd *exec.Cmd) {
	// Do nothing
}
