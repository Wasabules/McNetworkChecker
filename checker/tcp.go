package checker

import (
	"context"
	"fmt"
	"net"
	"time"
)

type TCPResult struct {
	Connected  bool   `json:"connected"`
	LatencyMs  int64  `json:"latencyMs"`
	LocalAddr  string `json:"localAddr"`
	RemoteAddr string `json:"remoteAddr"`
	Duration   int64  `json:"duration"`
	Error      string `json:"error,omitempty"`
	Success    bool   `json:"success"`
}

func CheckTCP(ctx context.Context, host string, port int, logFn LogFunc) TCPResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	result := TCPResult{}
	start := time.Now()
	defer func() { result.Duration = time.Since(start).Milliseconds() }()

	address := fmt.Sprintf("%s:%d", host, port)
	logFn(T("tcp.connecting", address))

	dialer := net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	result.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		logFn(T("tcp.fail", err.Error()))
		return result
	}
	defer conn.Close()

	result.Connected = true
	result.Success = true
	result.LocalAddr = conn.LocalAddr().String()
	result.RemoteAddr = conn.RemoteAddr().String()
	logFn(T("tcp.connected", result.LatencyMs))
	logFn(T("tcp.local", result.LocalAddr))
	logFn(T("tcp.remote", result.RemoteAddr))
	return result
}
