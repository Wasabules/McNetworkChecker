//go:build !windows

package checker

import "os/exec"

func hideConsole(_ *exec.Cmd) {}
