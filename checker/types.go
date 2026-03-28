package checker

// FullDiagnostic aggregates all step results for report generation.
type FullDiagnostic struct {
	Address       string           `json:"address"`
	ResolvedHost  string           `json:"resolvedHost"`
	ResolvedPort  int              `json:"resolvedPort"`
	DNS           DNSResult        `json:"dns"`
	Ping          PingResult       `json:"ping"`
	Traceroute    TracerouteResult `json:"traceroute"`
	TCP           TCPResult        `json:"tcp"`
	TCPTraceroute TracerouteResult `json:"tcpTraceroute"`
	Minecraft     MinecraftResult  `json:"minecraft"`
	SkippedSteps  []string         `json:"skippedSteps"`
}

func (d FullDiagnostic) IsSkipped(step string) bool {
	for _, s := range d.SkippedSteps {
		if s == step {
			return true
		}
	}
	return false
}
