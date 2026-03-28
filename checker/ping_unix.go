//go:build !windows

package checker

import (
	"context"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func icmpPing(ctx context.Context, host string, count int, timeout time.Duration, logFn LogFunc) PingResult {
	result := PingResult{Sent: count}

	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		result.Error = T("ping.resolveErr", err)
		return result
	}

	conn, network, err := listenICMP()
	if err != nil {
		result.Error = T("ping.sockErr", err)
		return result
	}
	defer conn.Close()

	id := os.Getpid() & 0xffff
	var totalMs float64
	result.MinMs = -1

	for seq := 0; seq < count; seq++ {
		if ctx.Err() != nil {
			break
		}
		msg := icmp.Message{Type: ipv4.ICMPTypeEcho, Code: 0, Body: &icmp.Echo{ID: id, Seq: seq, Data: []byte("McNetworkChecker")}}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			continue
		}
		conn.SetDeadline(time.Now().Add(timeout))

		var writeAddr net.Addr
		if network == "udp4" {
			writeAddr = &net.UDPAddr{IP: dst.IP}
		} else {
			writeAddr = dst
		}

		sendTime := time.Now()
		if _, err := conn.WriteTo(msgBytes, writeAddr); err != nil {
			logFn(T("ping.sendErr", err))
			continue
		}

		buf := make([]byte, 1500)
		n, _, err := conn.ReadFrom(buf)
		rtt := time.Since(sendTime)
		if err != nil {
			logFn(T("ping.timeout"))
			continue
		}

		rm, err := icmp.ParseMessage(1, buf[:n])
		if err != nil {
			continue
		}
		if rm.Type == ipv4.ICMPTypeEchoReply {
			rttMs := float64(rtt.Microseconds()) / 1000.0
			result.Received++
			totalMs += rttMs
			if result.MinMs < 0 || rttMs < result.MinMs {
				result.MinMs = rttMs
			}
			if rttMs > result.MaxMs {
				result.MaxMs = rttMs
			}
			logFn(T("ping.replySeq", dst.IP, seq+1, rttMs))
		}
	}

	result.Lost = result.Sent - result.Received
	if result.Received > 0 {
		result.AvgMs = totalMs / float64(result.Received)
		result.Success = true
	}
	if result.MinMs < 0 {
		result.MinMs = 0
	}
	return result
}

func listenICMP() (*icmp.PacketConn, string, error) {
	conn, err := icmp.ListenPacket("udp4", "")
	if err == nil {
		return conn, "udp4", nil
	}
	conn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err == nil {
		return conn, "ip4:icmp", nil
	}
	return nil, "", err
}
