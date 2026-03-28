//go:build !windows

package checker

import (
	"net"
	"time"
)

// icmpHop is a no-op on Unix. Getting intermediate router IPs requires
// raw ICMP sockets (root/CAP_NET_RAW). The TCP traceroute falls back to
// TCP-only probing which still detects when the destination is reached.
func icmpHop(_ net.IP, _ int, _ time.Duration) (routerIP string, reached bool, rttMs int64) {
	return "", false, 0
}
