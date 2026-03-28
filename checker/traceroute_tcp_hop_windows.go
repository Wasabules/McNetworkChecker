//go:build windows

package checker

import (
	"encoding/binary"
	"net"
	"time"
	"unsafe"
)

// ipOptionInfo matches the Windows IP_OPTION_INFORMATION struct layout.
type ipOptionInfo struct {
	Ttl         uint8
	Tos         uint8
	Flags       uint8
	OptionsSize uint8
	OptionsData uintptr
}

// icmpHop sends a single ICMP echo with the given TTL using the Windows IcmpSendEcho API.
// Returns the responding host IP (intermediate router or destination), whether the
// destination was reached, and the RTT. Works without admin privileges.
func icmpHop(destIP net.IP, ttl int, timeout time.Duration) (routerIP string, reached bool, rttMs int64) {
	ip4 := destIP.To4()
	if ip4 == nil {
		return "", false, 0
	}

	handle, _, _ := procIcmpCreateFile.Call()
	if handle == 0 || handle == ^uintptr(0) {
		return "", false, 0
	}
	defer procIcmpCloseHandle.Call(handle)

	destAddr := binary.LittleEndian.Uint32(ip4)
	sendData := []byte("MCTRACE\x00")
	replyBuf := make([]byte, 256)
	timeoutMs := uint32(timeout.Milliseconds())
	opts := ipOptionInfo{Ttl: uint8(ttl)}

	ret, _, _ := procIcmpSendEcho.Call(
		handle,
		uintptr(destAddr),
		uintptr(unsafe.Pointer(&sendData[0])),
		uintptr(len(sendData)),
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&replyBuf[0])),
		uintptr(len(replyBuf)),
		uintptr(timeoutMs),
	)

	if ret == 0 {
		return "", false, 0
	}

	// Reply layout: Address(4) + Status(4) + RoundTripTime(4)
	replyIP := net.IPv4(replyBuf[0], replyBuf[1], replyBuf[2], replyBuf[3])
	status := binary.LittleEndian.Uint32(replyBuf[4:8])
	rtt := binary.LittleEndian.Uint32(replyBuf[8:12])

	switch status {
	case icmpStatusSuccess:
		return replyIP.String(), true, int64(rtt)
	case icmpStatusTTLExpired:
		return replyIP.String(), false, int64(rtt)
	default:
		return "", false, 0
	}
}
