package main

import (
	"McNetworkChecker/checker"
	"context"
	"net"
	"strconv"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// beginStep creates a per-step cancellable context and resets the skip flag.
func (a *App) beginStep(ctx context.Context) (context.Context, context.CancelFunc) {
	stepCtx, cancel := context.WithCancel(ctx)
	a.stepMu.Lock()
	a.cancelStep = cancel
	a.stepSkipped = false
	a.stepMu.Unlock()
	return stepCtx, cancel
}

// endStep cleans up and reports whether the diagnostic was fully stopped or
// just the step was skipped. Uses the explicit stepSkipped flag — NOT stepCtx.Err()
// (which is always set after cancel() cleanup).
func (a *App) endStep(ctx context.Context, cancel context.CancelFunc) (stopped, skipped bool) {
	cancel() // always cleanup the child context
	a.stepMu.Lock()
	wasSkipped := a.stepSkipped
	a.cancelStep = nil
	a.stepSkipped = false
	a.stepMu.Unlock()

	// Parent context cancelled = full stop (even if skip was also pressed)
	if ctx.Err() != nil {
		return true, false
	}
	return false, wasSkipped
}

// runStep executes a single diagnostic step with skip/stop support.
// Returns true if the whole diagnostic should stop.
func (a *App) runStep(ctx context.Context, diag *checker.FullDiagnostic, stepID string, fn func(context.Context)) bool {
	stepCtx, stepCancel := a.beginStep(ctx)
	a.emitStep(stepID, "running", nil)
	fn(stepCtx)
	stopped, skipped := a.endStep(ctx, stepCancel)
	if stopped {
		return true
	}
	if skipped {
		a.logFor(stepID)(checker.T("diag.userSkip"))
		a.emitSkipped(stepID, "user")
		diag.SkippedSteps = append(diag.SkippedSteps, stepID)
	}
	return false
}

func (a *App) runDiagnostic(ctx context.Context, address string) {
	// Always notify the frontend when this goroutine exits (normal, stopped, or error)
	defer wailsRuntime.EventsEmit(a.ctx, "check:finished", nil)

	host, port := parseAddress(address)
	diag := checker.FullDiagnostic{Address: address}

	// Step 1: DNS (auto-skip if direct IP)
	var dnsResult checker.DNSResult
	if net.ParseIP(host) != nil {
		log := a.logFor("dns")
		log(checker.T("diag.ipDirect"))
		dnsResult = checker.DNSResult{
			ResolvedHost: host, ResolvedPort: port, Skipped: true, Success: true,
		}
		if net.ParseIP(host).To4() != nil {
			dnsResult.ARecords = []string{host}
		} else {
			dnsResult.AAAARecords = []string{host}
		}
		a.emitSkipped("dns", "auto")
	} else {
		if a.runStep(ctx, &diag, "dns", func(sCtx context.Context) {
			dnsResult = checker.ResolveDNS(sCtx, host, port, a.logFor("dns"))
		}) {
			return
		}
		if !diag.IsSkipped("dns") {
			a.emitStep("dns", "done", dnsResult)
		}
	}
	diag.DNS = dnsResult
	if dnsResult.Success {
		host = dnsResult.ResolvedHost
		port = dnsResult.ResolvedPort
	}
	diag.ResolvedHost = host
	diag.ResolvedPort = port

	pingTarget := host
	if len(dnsResult.ARecords) > 0 {
		pingTarget = dnsResult.ARecords[0]
	}

	// Step 2: Ping
	var pingResult checker.PingResult
	if a.runStep(ctx, &diag, "ping", func(sCtx context.Context) {
		pingResult = checker.RunPing(sCtx, pingTarget, a.logFor("ping"))
	}) {
		return
	}
	if !diag.IsSkipped("ping") {
		diag.Ping = pingResult
		a.emitStep("ping", "done", pingResult)
	}

	// Step 3: Traceroute ICMP
	var traceResult checker.TracerouteResult
	if a.runStep(ctx, &diag, "traceroute", func(sCtx context.Context) {
		traceResult = checker.RunTraceroute(sCtx, pingTarget, a.logFor("traceroute"))
	}) {
		return
	}
	if !diag.IsSkipped("traceroute") {
		diag.Traceroute = traceResult
		a.emitStep("traceroute", "done", traceResult)
	}

	// Parse ICMP traceroute output → hop-to-IP map for TCP traceroute enrichment
	var icmpHops map[int]string
	if traceResult.Success {
		icmpHops = checker.ParseHops(traceResult.Output)
	}

	// Step 4: TCP Check
	var tcpResult checker.TCPResult
	if a.runStep(ctx, &diag, "tcp", func(sCtx context.Context) {
		tcpResult = checker.CheckTCP(sCtx, host, port, a.logFor("tcp"))
	}) {
		return
	}
	if !diag.IsSkipped("tcp") {
		diag.TCP = tcpResult
		a.emitStep("tcp", "done", tcpResult)
	}

	// Step 5: TCP Traceroute
	var tcpTraceResult checker.TracerouteResult
	if a.runStep(ctx, &diag, "tcpTraceroute", func(sCtx context.Context) {
		tcpTraceResult = checker.RunTCPTraceroute(sCtx, pingTarget, port, icmpHops, a.logFor("tcpTraceroute"))
	}) {
		return
	}
	if !diag.IsSkipped("tcpTraceroute") {
		diag.TCPTraceroute = tcpTraceResult
		a.emitStep("tcpTraceroute", "done", tcpTraceResult)
	}

	// Step 6: Minecraft SLP
	var mcResult checker.MinecraftResult
	if a.runStep(ctx, &diag, "minecraft", func(sCtx context.Context) {
		mcResult = checker.QueryMinecraft(sCtx, host, port, a.logFor("minecraft"))
	}) {
		return
	}
	if !diag.IsSkipped("minecraft") {
		diag.Minecraft = mcResult
		a.emitStep("minecraft", "done", mcResult)
	}

	wailsRuntime.EventsEmit(a.ctx, "check:report", checker.GenerateReport(diag))
}

func parseAddress(address string) (string, int) {
	address = strings.TrimSpace(address)
	address = strings.TrimSuffix(address, ".")
	host, portStr, err := net.SplitHostPort(address)
	if err == nil {
		if p, err := strconv.Atoi(portStr); err == nil && p > 0 && p < 65536 {
			return host, p
		}
	}
	return address, 25565
}
