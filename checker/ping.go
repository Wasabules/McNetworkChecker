package checker

import (
	"context"
	"strings"
	"time"
)

type PingResult struct {
	Sent     int     `json:"sent"`
	Received int     `json:"received"`
	Lost     int     `json:"lost"`
	MinMs    float64 `json:"minMs"`
	AvgMs    float64 `json:"avgMs"`
	MaxMs    float64 `json:"maxMs"`
	Output   string  `json:"output"`
	Duration int64   `json:"duration"`
	Error    string  `json:"error,omitempty"`
	Success  bool    `json:"success"`
}

func RunPing(ctx context.Context, host string, logFn LogFunc) PingResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	start := time.Now()
	logFn(T("ping.sending", host))

	var lines []string
	capture := func(line string) {
		lines = append(lines, line)
		logFn(line)
	}

	result := icmpPing(ctx, host, 4, 3*time.Second, capture)
	result.Duration = time.Since(start).Milliseconds()
	result.Output = strings.Join(lines, "\n")

	if result.Success {
		logFn(T("ping.result", result.Received, result.Sent, result.MinMs, result.AvgMs, result.MaxMs))
	} else if result.Error != "" {
		logFn(T("ping.fail", result.Error))
	} else {
		logFn(T("ping.noReply"))
	}
	return result
}
