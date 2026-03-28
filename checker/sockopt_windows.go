//go:build windows

package checker

import "syscall"

func setSockTTL(fd uintptr, ttl int) error {
	// Windows: IPPROTO_IP=0, IP_TTL=4
	return syscall.SetsockoptInt(syscall.Handle(fd), 0, 4, ttl)
}
