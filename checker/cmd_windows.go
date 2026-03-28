//go:build windows

package checker

import (
	"os/exec"
	"syscall"
)

// hideConsole prevents a console window from flashing when running external commands.
func hideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
