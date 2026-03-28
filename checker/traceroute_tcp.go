package checker

import (
	"context"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
)

func RunTCPTraceroute(ctx context.Context, host string, port int, icmpHops map[int]string, logFn LogFunc) TracerouteResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	result := TracerouteResult{}
	start := time.Now()
	defer func() { result.Duration = time.Since(start).Milliseconds() }()

	maxHops := 20
	hopTimeout := 1 * time.Second
	address := fmt.Sprintf("%s:%d", host, port)

	destIP := net.ParseIP(host)
	if destIP == nil {
		if ips, err := net.LookupHost(host); err == nil && len(ips) > 0 {
			destIP = net.ParseIP(ips[0])
		}
	}

	logFn(T("tcptrace.start", address, maxHops))

	var lines []string
	for ttl := 1; ttl <= maxHops; ttl++ {
		if ctx.Err() != nil {
			break
		}

		type icmpRes struct {
			routerIP string
			reached  bool
			rttMs    int64
		}
		icmpCh := make(chan icmpRes, 1)
		go func(t int) {
			var r icmpRes
			if destIP != nil {
				r.routerIP, r.reached, r.rttMs = icmpHop(destIP, t, hopTimeout)
			}
			icmpCh <- r
		}(ttl)

		tcpReached, tcpRefused, tcpRtt, _ := tcpHop(ctx, address, ttl, hopTimeout)
		icmpR := <-icmpCh

		var line string

		if tcpReached {
			ip := icmpR.routerIP
			if ip == "" {
				ip = host
			}
			line = fmt.Sprintf("%2d  %s  %dms  [OPEN]", ttl, ip, tcpRtt.Milliseconds())
			logFn(line)
			lines = append(lines, line)
			result.Hops = ttl
			result.Success = true
			logFn(T("tcptrace.open", port, ttl))
			break
		}
		if tcpRefused {
			ip := icmpR.routerIP
			if ip == "" {
				ip = host
			}
			line = fmt.Sprintf("%2d  %s  %dms  [CLOSED]", ttl, ip, tcpRtt.Milliseconds())
			logFn(line)
			lines = append(lines, line)
			result.Hops = ttl
			result.Success = true
			logFn(T("tcptrace.closed", port, ttl))
			break
		}

		// Resolve router IP: live ICMP > cached ICMP traceroute > unknown
		routerIP := icmpR.routerIP
		if routerIP == "" && icmpHops != nil {
			routerIP = icmpHops[ttl]
		}

		if icmpR.reached {
			ip := routerIP
			if ip == "" {
				ip = host
			}
			line = fmt.Sprintf("%2d  %s  %dms  [FILTERED]", ttl, ip, icmpR.rttMs)
			logFn(line)
			lines = append(lines, line)
			result.Hops = ttl
			result.Success = true
			logFn(T("tcptrace.filtered", port))
			break
		}

		if routerIP != "" {
			rtt := icmpR.rttMs
			if rtt == 0 && tcpRtt < hopTimeout-200*time.Millisecond {
				rtt = tcpRtt.Milliseconds()
			}
			line = fmt.Sprintf("%2d  %s  %dms", ttl, routerIP, rtt)
		} else {
			if tcpRtt < hopTimeout-200*time.Millisecond {
				line = fmt.Sprintf("%2d  *  %dms", ttl, tcpRtt.Milliseconds())
			} else {
				line = fmt.Sprintf("%2d  *", ttl)
			}
		}
		logFn(line)
		lines = append(lines, line)
	}

	result.Output = strings.Join(lines, "\n")
	if !result.Success {
		logFn(T("tcptrace.noReach"))
	}
	return result
}

func tcpHop(ctx context.Context, address string, ttl int, timeout time.Duration) (reached, refused bool, rtt time.Duration, err error) {
	d := net.Dialer{
		Timeout: timeout,
		Control: func(network, addr string, c syscall.RawConn) error {
			var setErr error
			ctrlErr := c.Control(func(fd uintptr) { setErr = setSockTTL(fd, ttl) })
			if ctrlErr != nil {
				return ctrlErr
			}
			return setErr
		},
	}
	start := time.Now()
	conn, dialErr := d.DialContext(ctx, "tcp", address)
	rtt = time.Since(start)
	if conn != nil {
		conn.Close()
		return true, false, rtt, nil
	}
	if dialErr != nil {
		errStr := strings.ToLower(dialErr.Error())
		if strings.Contains(errStr, "refused") || strings.Contains(errStr, "reset") {
			return false, true, rtt, dialErr
		}
	}
	return false, false, rtt, dialErr
}
