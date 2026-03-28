package checker

import (
	"context"
	"net"
	"strings"
	"time"
)

type SRVRecord struct {
	Target   string `json:"target"`
	Port     int    `json:"port"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
}

type DNSResult struct {
	SRVRecords   []SRVRecord `json:"srvRecords"`
	ARecords     []string    `json:"aRecords"`
	AAAARecords  []string    `json:"aaaaRecords"`
	CNAME        string      `json:"cname"`
	ResolvedHost string      `json:"resolvedHost"`
	ResolvedPort int         `json:"resolvedPort"`
	Duration     int64       `json:"duration"`
	Error        string      `json:"error,omitempty"`
	Skipped      bool        `json:"skipped"`
	Success      bool        `json:"success"`
}

func ResolveDNS(ctx context.Context, host string, port int, logFn LogFunc) DNSResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	result := DNSResult{ResolvedHost: host, ResolvedPort: port}
	start := time.Now()
	defer func() { result.Duration = time.Since(start).Milliseconds() }()

	if ip := net.ParseIP(host); ip != nil {
		logFn(T("dns.ipDetected", host))
		if ip.To4() != nil {
			result.ARecords = []string{host}
		} else {
			result.AAAARecords = []string{host}
		}
		result.Success = true
		return result
	}

	logFn(T("dns.searchSrv", host))
	_, srvs, err := net.DefaultResolver.LookupSRV(ctx, "minecraft", "tcp", host)
	if err == nil && len(srvs) > 0 {
		for _, srv := range srvs {
			target := strings.TrimSuffix(srv.Target, ".")
			result.SRVRecords = append(result.SRVRecords, SRVRecord{Target: target, Port: int(srv.Port), Priority: int(srv.Priority), Weight: int(srv.Weight)})
			logFn(T("dns.srvFound", target, srv.Port, srv.Priority, srv.Weight))
		}
		result.ResolvedHost = strings.TrimSuffix(srvs[0].Target, ".")
		result.ResolvedPort = int(srvs[0].Port)
	} else {
		logFn(T("dns.noSrv"))
	}

	logFn(T("dns.searchCname", result.ResolvedHost))
	cname, err := net.DefaultResolver.LookupCNAME(ctx, result.ResolvedHost)
	if err == nil && cname != "" {
		cname = strings.TrimSuffix(cname, ".")
		if cname != result.ResolvedHost {
			result.CNAME = cname
			logFn(T("dns.cnameFound", cname))
		} else {
			logFn(T("dns.noCname"))
		}
	} else {
		logFn(T("dns.noCname"))
	}

	logFn(T("dns.searchA", result.ResolvedHost))
	ips, err := net.DefaultResolver.LookupHost(ctx, result.ResolvedHost)
	if err != nil {
		result.Error = err.Error()
		logFn(T("dns.fail", err.Error()))
		return result
	}
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed == nil {
			continue
		}
		if parsed.To4() != nil {
			result.ARecords = append(result.ARecords, ip)
			logFn("  A    -> " + ip)
		} else {
			result.AAAARecords = append(result.AAAARecords, ip)
			logFn("  AAAA -> " + ip)
		}
	}
	result.Success = len(result.ARecords) > 0 || len(result.AAAARecords) > 0
	logFn(T("dns.resolved", result.ResolvedHost, result.ResolvedPort))
	return result
}
