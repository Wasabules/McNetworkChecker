//go:build windows

package checker

import (
	"context"
	"encoding/binary"
	"net"
	"syscall"
	"time"
	"unsafe"
)

var (
	modIphlpapi         = syscall.NewLazyDLL("iphlpapi.dll")
	procIcmpCreateFile  = modIphlpapi.NewProc("IcmpCreateFile")
	procIcmpSendEcho    = modIphlpapi.NewProc("IcmpSendEcho")
	procIcmpCloseHandle = modIphlpapi.NewProc("IcmpCloseHandle")
)

const (
	icmpStatusSuccess    = 0
	icmpStatusBufSmall   = 11001
	icmpStatusDestUnrch  = 11003
	icmpStatusTimedOut   = 11010
	icmpStatusTTLExpired = 11013
)

func icmpPing(ctx context.Context, host string, count int, timeout time.Duration, logFn LogFunc) PingResult {
	result := PingResult{Sent: count}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupHost(host)
		if err != nil || len(ips) == 0 {
			result.Error = T("ping.resolveErr", err)
			return result
		}
		ip = net.ParseIP(ips[0])
	}
	ip4 := ip.To4()
	if ip4 == nil {
		result.Error = T("ping.ipv4Only")
		return result
	}

	handle, _, err := procIcmpCreateFile.Call()
	if handle == 0 || handle == ^uintptr(0) {
		result.Error = "IcmpCreateFile: " + err.Error()
		return result
	}
	defer procIcmpCloseHandle.Call(handle)

	destAddr := binary.LittleEndian.Uint32(ip4)
	sendData := []byte("McNetworkChecker")
	replyBuf := make([]byte, 256)
	timeoutMs := uint32(timeout.Milliseconds())
	var totalMs float64
	result.MinMs = -1

	for i := 0; i < count; i++ {
		if ctx.Err() != nil {
			break
		}
		ret, _, _ := procIcmpSendEcho.Call(
			handle, uintptr(destAddr),
			uintptr(unsafe.Pointer(&sendData[0])), uintptr(len(sendData)),
			0,
			uintptr(unsafe.Pointer(&replyBuf[0])), uintptr(len(replyBuf)),
			uintptr(timeoutMs),
		)
		if ret > 0 {
			replyIP := net.IPv4(replyBuf[0], replyBuf[1], replyBuf[2], replyBuf[3])
			status := binary.LittleEndian.Uint32(replyBuf[4:8])
			rtt := binary.LittleEndian.Uint32(replyBuf[8:12])
			switch status {
			case icmpStatusSuccess:
				rttF := float64(rtt)
				result.Received++
				totalMs += rttF
				if result.MinMs < 0 || rttF < result.MinMs {
					result.MinMs = rttF
				}
				if rttF > result.MaxMs {
					result.MaxMs = rttF
				}
				logFn(T("ping.reply", replyIP, rtt))
			case icmpStatusTimedOut:
				logFn(T("ping.timeout"))
			case icmpStatusDestUnrch:
				logFn(T("ping.unreachable", replyIP))
			case icmpStatusTTLExpired:
				logFn(T("ping.ttlExpired", replyIP))
			default:
				logFn(T("ping.icmpErr", status))
			}
		} else {
			logFn(T("ping.timeout"))
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
