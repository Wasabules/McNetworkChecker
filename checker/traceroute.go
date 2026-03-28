package checker

import (
	"context"
	"runtime"
	"strings"
	"time"
)

type TracerouteResult struct {
	Output   string `json:"output"`
	Hops     int    `json:"hops"`
	Duration int64  `json:"duration"`
	Error    string `json:"error,omitempty"`
	Success  bool   `json:"success"`
}

func RunTraceroute(ctx context.Context, host string, logFn LogFunc) TracerouteResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	result := TracerouteResult{}
	start := time.Now()
	defer func() { result.Duration = time.Since(start).Milliseconds() }()

	logFn(T("trace.sysCmd"))
	logFn(T("trace.start", host))

	var name string
	var args []string
	switch runtime.GOOS {
	case "windows":
		name = "tracert"
		args = []string{"-d", "-w", "2000", "-h", "20", host}
	default:
		name = "traceroute"
		args = []string{"-n", "-w", "2", "-m", "20", host}
	}

	output, err := runCmdStreaming(ctx, logFn, name, args...)
	result.Output = output
	if err != nil && output == "" {
		result.Error = err.Error()
		logFn(T("trace.fail", err.Error()))
		return result
	}
	result.Hops = countHops(result.Output)
	result.Success = result.Hops > 0
	if result.Success {
		logFn(T("trace.done", result.Hops))
	}
	return result
}

func countHops(output string) int {
	hops := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 1 && trimmed[0] >= '1' && trimmed[0] <= '9' {
			hops++
		}
	}
	return hops
}

func ParseHops(output string) map[int]string {
	hops := make(map[int]string)
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 3 || trimmed[0] < '1' || trimmed[0] > '9' {
			continue
		}
		i := 0
		for i < len(trimmed) && trimmed[i] >= '0' && trimmed[i] <= '9' {
			i++
		}
		hopNum := 0
		for _, c := range trimmed[:i] {
			hopNum = hopNum*10 + int(c-'0')
		}
		if ip := findIPv4(trimmed[i:]); ip != "" {
			hops[hopNum] = ip
		}
	}
	return hops
}

func findIPv4(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			end := i
			dots := 0
			for end < len(s) && (s[end] == '.' || (s[end] >= '0' && s[end] <= '9')) {
				if s[end] == '.' {
					dots++
				}
				end++
			}
			if dots == 3 && end-i >= 7 {
				parts := strings.Split(s[i:end], ".")
				if len(parts) == 4 {
					valid := true
					for _, p := range parts {
						if len(p) == 0 || len(p) > 3 {
							valid = false
							break
						}
						n := 0
						for _, c := range p {
							n = n*10 + int(c-'0')
						}
						if n > 255 {
							valid = false
							break
						}
					}
					if valid {
						return s[i:end]
					}
				}
			}
			i = end - 1
		}
	}
	return ""
}
