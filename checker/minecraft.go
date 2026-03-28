package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type MinecraftResult struct {
	Version       string            `json:"version"`
	ProtocolVer   int               `json:"protocolVersion"`
	MOTD          string            `json:"motd"`
	MOTDSpans     []MOTDSpan        `json:"motdSpans"`
	PlayersOnline int               `json:"playersOnline"`
	PlayersMax    int               `json:"playersMax"`
	PlayersSample []MinecraftPlayer `json:"playersSample"`
	Favicon       string            `json:"favicon"`
	LatencyMs     int64             `json:"latencyMs"`
	Duration      int64             `json:"duration"`
	RawJSON       string            `json:"rawJson"`
	Error         string            `json:"error,omitempty"`
	Success       bool              `json:"success"`
}

type MinecraftPlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type serverStatus struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description json.RawMessage `json:"description"`
	Favicon     string          `json:"favicon"`
}

func QueryMinecraft(ctx context.Context, host string, port int, logFn LogFunc) MinecraftResult {
	if logFn == nil {
		logFn = func(_ string) {}
	}
	result := MinecraftResult{}
	start := time.Now()
	defer func() { result.Duration = time.Since(start).Milliseconds() }()

	address := fmt.Sprintf("%s:%d", host, port)
	logFn(T("mc.connecting", address))

	dialer := net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		result.Error = T("mc.connFail", err.Error())
		logFn(T("tcp.fail", err.Error()))
		return result
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	logFn(T("mc.tcpOk"))

	logFn(T("mc.handshake"))
	if _, err := conn.Write(buildHandshakePacket(host, port)); err != nil {
		result.Error = T("mc.handshakeFail", err.Error())
		logFn(T("tcp.fail", err.Error()))
		return result
	}

	logFn(T("mc.statusReq"))
	if _, err := conn.Write(buildStatusRequestPacket()); err != nil {
		result.Error = T("mc.statusFail", err.Error())
		logFn(T("tcp.fail", err.Error()))
		return result
	}

	logFn(T("mc.reading"))
	jsonStr, err := readStatusResponse(conn)
	if err != nil {
		result.Error = T("mc.readFail", err.Error())
		logFn(T("tcp.fail", err.Error()))
		return result
	}
	result.RawJSON = jsonStr
	logFn(T("mc.jsonOk"))

	var status serverStatus
	if err := json.Unmarshal([]byte(jsonStr), &status); err != nil {
		result.Error = T("mc.jsonFail", err.Error())
		logFn(T("tcp.fail", err.Error()))
		return result
	}

	result.Version = status.Version.Name
	result.ProtocolVer = status.Version.Protocol
	result.PlayersOnline = status.Players.Online
	result.PlayersMax = status.Players.Max
	result.Favicon = status.Favicon
	result.MOTD, result.MOTDSpans = parseDescription(status.Description)

	logFn(T("mc.version", result.Version, result.ProtocolVer))
	logFn(T("mc.players", result.PlayersOnline, result.PlayersMax))
	logFn(T("mc.motd", result.MOTD))

	for _, p := range status.Players.Sample {
		result.PlayersSample = append(result.PlayersSample, MinecraftPlayer{Name: p.Name, ID: p.ID})
	}
	if len(result.PlayersSample) > 0 {
		names := make([]string, len(result.PlayersSample))
		for i, p := range result.PlayersSample {
			names[i] = p.Name
		}
		logFn(T("mc.online", strings.Join(names, ", ")))
	}

	logFn(T("mc.latency"))
	pingStart := time.Now()
	payload := time.Now().UnixMilli()
	if _, err := conn.Write(buildPingPacket(payload)); err == nil {
		if pong, err := readPongResponse(conn); err == nil && pong == payload {
			result.LatencyMs = time.Since(pingStart).Milliseconds()
			logFn(T("mc.latencyOk", result.LatencyMs))
		} else {
			logFn(T("mc.pongFail"))
		}
	}

	if result.Favicon != "" {
		logFn(T("mc.favicon"))
	}
	result.Success = true
	logFn(T("mc.done"))
	return result
}
