package checker

import (
	"fmt"
	"strings"
	"time"
)

func GenerateReport(diag FullDiagnostic) string {
	var sb strings.Builder
	w := func(format string, args ...interface{}) { sb.WriteString(fmt.Sprintf(format, args...)) }

	w("%s\n", T("rpt.header"))
	w("%s\n", T("rpt.server", diag.Address))
	w("%s\n", T("rpt.date", time.Now().Format("2006-01-02 15:04:05")))
	w("%s\n\n", T("rpt.resolved", diag.ResolvedHost, diag.ResolvedPort))

	// DNS
	w("%s\n", T("rpt.dns"))
	if diag.IsSkipped("dns") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.DNS.Skipped {
		w("%s\n", T("rpt.skipIp"))
	} else if diag.DNS.Success {
		w("%s\n", T("rpt.ok", diag.DNS.Duration))
	} else {
		w("%s\n", T("rpt.failDur", diag.DNS.Duration))
		if diag.DNS.Error != "" {
			w("%s\n", T("rpt.error", diag.DNS.Error))
		}
	}
	if len(diag.DNS.SRVRecords) > 0 {
		w("SRV :\n")
		for _, srv := range diag.DNS.SRVRecords {
			w("  %s:%d (prio=%d, w=%d)\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
		}
	}
	if diag.DNS.CNAME != "" {
		w("CNAME : %s\n", diag.DNS.CNAME)
	}
	if len(diag.DNS.ARecords) > 0 {
		w("A     : %s\n", strings.Join(diag.DNS.ARecords, ", "))
	}
	if len(diag.DNS.AAAARecords) > 0 {
		w("AAAA  : %s\n", strings.Join(diag.DNS.AAAARecords, ", "))
	}
	w("\n")

	// Ping
	w("%s\n", T("rpt.ping"))
	if diag.IsSkipped("ping") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.Ping.Success {
		w("%s\n", T("rpt.pingOk", diag.Ping.Received, diag.Ping.Sent))
		w("%s\n", T("rpt.latency", diag.Ping.MinMs, diag.Ping.AvgMs, diag.Ping.MaxMs))
	} else {
		w("%s\n", T("rpt.fail"))
		if diag.Ping.Error != "" {
			w("%s\n", T("rpt.error", diag.Ping.Error))
		}
	}
	w("\n")

	// Traceroute ICMP
	w("%s\n", T("rpt.traceIcmp"))
	if diag.IsSkipped("traceroute") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.Traceroute.Success {
		w("%s\n", T("rpt.hops", diag.Traceroute.Hops))
		if diag.Traceroute.Output != "" {
			w("%s\n", diag.Traceroute.Output)
		}
	} else {
		w("%s\n", T("rpt.fail"))
		if diag.Traceroute.Error != "" {
			w("%s\n", T("rpt.error", diag.Traceroute.Error))
		}
	}
	w("\n")

	// TCP
	w("%s\n", T("rpt.tcp", diag.ResolvedPort))
	if diag.IsSkipped("tcp") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.TCP.Success {
		w("%s\n", T("rpt.tcpOk", diag.TCP.LatencyMs))
		w("%s\n", T("rpt.local", diag.TCP.LocalAddr))
		w("%s\n", T("rpt.remote", diag.TCP.RemoteAddr))
	} else {
		w("%s\n", T("rpt.fail"))
		if diag.TCP.Error != "" {
			w("%s\n", T("rpt.error", diag.TCP.Error))
		}
	}
	w("\n")

	// TCP Traceroute
	w("%s\n", T("rpt.traceTcp", diag.ResolvedPort))
	if diag.IsSkipped("tcpTraceroute") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.TCPTraceroute.Success {
		w("%s\n", T("rpt.hops", diag.TCPTraceroute.Hops))
		if diag.TCPTraceroute.Output != "" {
			w("%s\n", diag.TCPTraceroute.Output)
		}
	} else if diag.TCPTraceroute.Error != "" {
		w("%s\n", T("rpt.notAvail"))
		w("%s\n", T("rpt.note", diag.TCPTraceroute.Error))
	} else {
		w("%s\n", T("rpt.fail"))
	}
	w("\n")

	// Minecraft
	w("%s\n", T("rpt.minecraft"))
	if diag.IsSkipped("minecraft") {
		w("%s\n", T("rpt.skipUser"))
	} else if diag.Minecraft.Success {
		w("%s\n", T("rpt.okDuration", diag.Minecraft.Duration))
		w("%s\n", T("rpt.mcVersion", diag.Minecraft.Version, diag.Minecraft.ProtocolVer))
		w("%s\n", T("rpt.mcPlayers", diag.Minecraft.PlayersOnline, diag.Minecraft.PlayersMax))
		w("%s\n", T("rpt.mcMotd", diag.Minecraft.MOTD))
		if len(diag.Minecraft.PlayersSample) > 0 {
			w("%s\n", T("rpt.mcOnline"))
			for _, p := range diag.Minecraft.PlayersSample {
				w("  - %s (%s)\n", p.Name, p.ID)
			}
		}
		w("%s\n", T("rpt.mcLatency", diag.Minecraft.LatencyMs))
	} else {
		w("%s\n", T("rpt.fail"))
		if diag.Minecraft.Error != "" {
			w("%s\n", T("rpt.error", diag.Minecraft.Error))
		}
	}
	w("\n%s\n", T("rpt.footer"))
	return sb.String()
}
