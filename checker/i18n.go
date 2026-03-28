package checker

import (
	"fmt"
	"sync"
)

const (
	_FR = iota
	_EN
	_ES
	_DE
	_PT
	_NL
)

type tr [_NL]string

var (
	currentLocale int = _FR
	localeMu      sync.RWMutex
)

func SetLocale(code string) {
	localeMu.Lock()
	defer localeMu.Unlock()
	switch code {
	case "fr":
		currentLocale = _FR
	case "en":
		currentLocale = _EN
	case "es":
		currentLocale = _ES
	case "de":
		currentLocale = _DE
	case "pt":
		currentLocale = _PT
	default:
		currentLocale = _EN
	}
}

func GetLocale() string {
	localeMu.RLock()
	defer localeMu.RUnlock()
	return [...]string{"fr", "en", "es", "de", "pt"}[currentLocale]
}

// T translates a key with optional fmt.Sprintf arguments.
func T(key string, args ...interface{}) string {
	localeMu.RLock()
	idx := currentLocale
	localeMu.RUnlock()

	if m, ok := msgs[key]; ok {
		s := m[idx]
		if s == "" {
			s = m[_EN]
		}
		if len(args) > 0 {
			return fmt.Sprintf(s, args...)
		}
		return s
	}
	return key
}

// ========================
//  ALL BACKEND TRANSLATIONS
// ========================
var msgs = map[string]tr{
	// ---- DNS ----
	"dns.ipDetected":  {"Adresse IP detectee : %s", "IP address detected: %s", "Dirección IP detectada: %s", "IP-Adresse erkannt: %s", "Endereço IP detectado: %s"},
	"dns.searchSrv":   {"Recherche SRV _minecraft._tcp.%s ...", "SRV lookup _minecraft._tcp.%s ...", "Búsqueda SRV _minecraft._tcp.%s ...", "SRV-Abfrage _minecraft._tcp.%s ...", "Consulta SRV _minecraft._tcp.%s ..."},
	"dns.srvFound":    {"  SRV -> %s:%d (prio=%d, poids=%d)", "  SRV -> %s:%d (prio=%d, weight=%d)", "  SRV -> %s:%d (prio=%d, peso=%d)", "  SRV -> %s:%d (Prio=%d, Gewicht=%d)", "  SRV -> %s:%d (prio=%d, peso=%d)"},
	"dns.noSrv":       {"  Aucun enregistrement SRV", "  No SRV record found", "  Sin registro SRV", "  Kein SRV-Eintrag gefunden", "  Nenhum registro SRV encontrado"},
	"dns.searchCname":  {"Recherche CNAME pour %s ...", "CNAME lookup for %s ...", "Búsqueda CNAME para %s ...", "CNAME-Abfrage für %s ...", "Consulta CNAME para %s ..."},
	"dns.cnameFound":   {"  CNAME -> %s", "  CNAME -> %s", "  CNAME -> %s", "  CNAME -> %s", "  CNAME -> %s"},
	"dns.noCname":      {"  Pas de CNAME", "  No CNAME", "  Sin CNAME", "  Kein CNAME", "  Sem CNAME"},
	"dns.searchA":      {"Resolution A/AAAA pour %s ...", "A/AAAA resolution for %s ...", "Resolución A/AAAA para %s ...", "A/AAAA-Auflösung für %s ...", "Resolução A/AAAA para %s ..."},
	"dns.fail":         {"  ECHEC : %s", "  FAILED: %s", "  FALLO: %s", "  FEHLER: %s", "  FALHA: %s"},
	"dns.resolved":     {"Hote resolu : %s:%d", "Host resolved: %s:%d", "Host resuelto: %s:%d", "Host aufgelöst: %s:%d", "Host resolvido: %s:%d"},

	// ---- Ping ----
	"ping.sending":  {"Envoi de 4 paquets ICMP vers %s (ping natif) ...", "Sending 4 ICMP packets to %s (native ping) ...", "Enviando 4 paquetes ICMP a %s (ping nativo) ...", "Sende 4 ICMP-Pakete an %s (nativer Ping) ...", "Enviando 4 pacotes ICMP para %s (ping nativo) ..."},
	"ping.result":   {"Resultat : %d/%d paquets, min=%.0fms moy=%.0fms max=%.0fms", "Result: %d/%d packets, min=%.0fms avg=%.0fms max=%.0fms", "Resultado: %d/%d paquetes, min=%.0fms med=%.0fms max=%.0fms", "Ergebnis: %d/%d Pakete, Min=%.0fms Avg=%.0fms Max=%.0fms", "Resultado: %d/%d pacotes, min=%.0fms med=%.0fms max=%.0fms"},
	"ping.noReply":  {"Aucun paquet recu", "No packets received", "Ningún paquete recibido", "Keine Pakete empfangen", "Nenhum pacote recebido"},
	"ping.fail":     {"ECHEC : %s", "FAILED: %s", "FALLO: %s", "FEHLER: %s", "FALHA: %s"},

	// ---- TCP ----
	"tcp.connecting": {"Connexion TCP vers %s ...", "TCP connection to %s ...", "Conexión TCP a %s ...", "TCP-Verbindung zu %s ...", "Conexão TCP para %s ..."},
	"tcp.fail":       {"ECHEC : %s", "FAILED: %s", "FALLO: %s", "FEHLER: %s", "FALHA: %s"},
	"tcp.connected":  {"Connecte en %dms", "Connected in %dms", "Conectado en %dms", "Verbunden in %dms", "Conectado em %dms"},
	"tcp.local":      {"  Local  : %s", "  Local : %s", "  Local : %s", "  Lokal : %s", "  Local : %s"},
	"tcp.remote":     {"  Distant: %s", "  Remote: %s", "  Remoto: %s", "  Remote: %s", "  Remoto: %s"},

	// ---- Traceroute ICMP ----
	"trace.sysCmd":  {"Utilisation de la commande systeme (tracert/traceroute)", "Using system command (tracert/traceroute)", "Usando comando del sistema (tracert/traceroute)", "Verwende Systembefehl (tracert/traceroute)", "Usando comando do sistema (tracert/traceroute)"},
	"trace.start":   {"Traceroute ICMP vers %s (max 20 sauts) ...", "ICMP traceroute to %s (max 20 hops) ...", "Traceroute ICMP a %s (max 20 saltos) ...", "ICMP-Traceroute zu %s (max 20 Hops) ...", "Traceroute ICMP para %s (max 20 saltos) ..."},
	"trace.fail":    {"ECHEC : %s", "FAILED: %s", "FALLO: %s", "FEHLER: %s", "FALHA: %s"},
	"trace.done":    {"Termine : %d sauts", "Done: %d hops", "Terminado: %d saltos", "Fertig: %d Hops", "Terminado: %d saltos"},

	// ---- Traceroute TCP ----
	"tcptrace.start":    {"Traceroute TCP vers %s (max %d sauts) ...", "TCP traceroute to %s (max %d hops) ...", "Traceroute TCP a %s (max %d saltos) ...", "TCP-Traceroute zu %s (max %d Hops) ...", "Traceroute TCP para %s (max %d saltos) ..."},
	"tcptrace.open":     {"Port %d ouvert, destination atteinte en %d sauts", "Port %d open, destination reached in %d hops", "Puerto %d abierto, destino alcanzado en %d saltos", "Port %d offen, Ziel in %d Hops erreicht", "Porta %d aberta, destino alcançado em %d saltos"},
	"tcptrace.closed":   {"Port %d ferme, hote atteint en %d sauts", "Port %d closed, host reached in %d hops", "Puerto %d cerrado, host alcanzado en %d saltos", "Port %d geschlossen, Host in %d Hops erreicht", "Porta %d fechada, host alcançado em %d saltos"},
	"tcptrace.filtered": {"Hote atteint mais port %d filtre", "Host reached but port %d filtered", "Host alcanzado pero puerto %d filtrado", "Host erreicht aber Port %d gefiltert", "Host alcançado mas porta %d filtrada"},
	"tcptrace.noReach":  {"Destination non atteinte", "Destination not reached", "Destino no alcanzado", "Ziel nicht erreicht", "Destino não alcançado"},

	// ---- Minecraft ----
	"mc.connecting":  {"Connexion au serveur Minecraft %s ...", "Connecting to Minecraft server %s ...", "Conectando al servidor Minecraft %s ...", "Verbinde mit Minecraft-Server %s ...", "Conectando ao servidor Minecraft %s ..."},
	"mc.connFail":    {"Connexion echouee : %s", "Connection failed: %s", "Conexión fallida: %s", "Verbindung fehlgeschlagen: %s", "Conexão falhou: %s"},
	"mc.tcpOk":       {"Connexion TCP etablie", "TCP connection established", "Conexión TCP establecida", "TCP-Verbindung hergestellt", "Conexão TCP estabelecida"},
	"mc.handshake":   {"Envoi du handshake (protocole 767) ...", "Sending handshake (protocol 767) ...", "Enviando handshake (protocolo 767) ...", "Sende Handshake (Protokoll 767) ...", "Enviando handshake (protocolo 767) ..."},
	"mc.handshakeFail": {"Erreur envoi handshake : %s", "Handshake send error: %s", "Error al enviar handshake: %s", "Fehler beim Senden des Handshakes: %s", "Erro ao enviar handshake: %s"},
	"mc.statusReq":   {"Envoi de la requete de statut ...", "Sending status request ...", "Enviando solicitud de estado ...", "Sende Statusanfrage ...", "Enviando pedido de estado ..."},
	"mc.statusFail":  {"Erreur envoi status request : %s", "Status request error: %s", "Error en solicitud de estado: %s", "Fehler bei Statusanfrage: %s", "Erro no pedido de estado: %s"},
	"mc.reading":     {"Lecture de la reponse ...", "Reading response ...", "Leyendo respuesta ...", "Lese Antwort ...", "Lendo resposta ..."},
	"mc.readFail":    {"Erreur lecture reponse : %s", "Response read error: %s", "Error al leer respuesta: %s", "Fehler beim Lesen der Antwort: %s", "Erro ao ler resposta: %s"},
	"mc.jsonOk":      {"Reponse JSON recue", "JSON response received", "Respuesta JSON recibida", "JSON-Antwort empfangen", "Resposta JSON recebida"},
	"mc.jsonFail":    {"Erreur parsing JSON : %s", "JSON parse error: %s", "Error al parsear JSON: %s", "JSON-Parsing-Fehler: %s", "Erro ao parsear JSON: %s"},
	"mc.version":     {"Version  : %s (protocole %d)", "Version : %s (protocol %d)", "Versión : %s (protocolo %d)", "Version : %s (Protokoll %d)", "Versão  : %s (protocolo %d)"},
	"mc.players":     {"Joueurs  : %d/%d", "Players : %d/%d", "Jugadores: %d/%d", "Spieler : %d/%d", "Jogadores: %d/%d"},
	"mc.motd":        {"MOTD     : %s", "MOTD    : %s", "MOTD    : %s", "MOTD    : %s", "MOTD    : %s"},
	"mc.online":      {"En ligne : %s", "Online  : %s", "En línea: %s", "Online  : %s", "Online  : %s"},
	"mc.latency":     {"Mesure de la latence ...", "Measuring latency ...", "Midiendo latencia ...", "Messe Latenz ...", "Medindo latência ..."},
	"mc.latencyOk":   {"Latence protocole : %dms", "Protocol latency: %dms", "Latencia protocolo: %dms", "Protokoll-Latenz: %dms", "Latência protocolo: %dms"},
	"mc.pongFail":    {"Pong invalide ou timeout", "Invalid pong or timeout", "Pong inválido o timeout", "Ungültiger Pong oder Timeout", "Pong inválido ou timeout"},
	"mc.favicon":     {"Favicon recu", "Favicon received", "Favicon recibido", "Favicon empfangen", "Favicon recebido"},
	"mc.done":        {"Diagnostic Minecraft termine avec succes", "Minecraft diagnostic completed", "Diagnóstico Minecraft completado", "Minecraft-Diagnose abgeschlossen", "Diagnóstico Minecraft concluído"},

	// ---- Diagnostic flow ----
	"diag.ipDirect":   {"Adresse IP directe detectee, resolution DNS non necessaire", "Direct IP address detected, DNS resolution not needed", "Dirección IP directa, resolución DNS no necesaria", "Direkte IP-Adresse, DNS-Auflösung nicht nötig", "Endereço IP direto, resolução DNS não necessária"},
	"diag.userSkip":   {"Etape passee par l'utilisateur", "Step skipped by user", "Paso omitido por el usuario", "Schritt vom Benutzer übersprungen", "Etapa ignorada pelo utilizador"},

	// ---- Report ----
	"rpt.header":       {"=== Rapport Diagnostic McNetworkChecker ===", "=== McNetworkChecker Diagnostic Report ===", "=== Informe Diagnóstico McNetworkChecker ===", "=== McNetworkChecker Diagnosebericht ===", "=== Relatório Diagnóstico McNetworkChecker ==="},
	"rpt.server":       {"Serveur : %s", "Server : %s", "Servidor: %s", "Server : %s", "Servidor: %s"},
	"rpt.date":         {"Date    : %s", "Date   : %s", "Fecha  : %s", "Datum  : %s", "Data   : %s"},
	"rpt.resolved":     {"Resolu  : %s:%d", "Resolved: %s:%d", "Resuelto: %s:%d", "Aufgelöst: %s:%d", "Resolvido: %s:%d"},
	"rpt.dns":          {"--- Resolution DNS ---", "--- DNS Resolution ---", "--- Resolución DNS ---", "--- DNS-Auflösung ---", "--- Resolução DNS ---"},
	"rpt.ping":         {"--- Ping ICMP ---", "--- ICMP Ping ---", "--- Ping ICMP ---", "--- ICMP-Ping ---", "--- Ping ICMP ---"},
	"rpt.traceIcmp":    {"--- Traceroute ICMP ---", "--- ICMP Traceroute ---", "--- Traceroute ICMP ---", "--- ICMP-Traceroute ---", "--- Traceroute ICMP ---"},
	"rpt.tcp":          {"--- Connexion TCP (port %d) ---", "--- TCP Connection (port %d) ---", "--- Conexión TCP (puerto %d) ---", "--- TCP-Verbindung (Port %d) ---", "--- Conexão TCP (porta %d) ---"},
	"rpt.traceTcp":     {"--- Traceroute TCP (port %d) ---", "--- TCP Traceroute (port %d) ---", "--- Traceroute TCP (puerto %d) ---", "--- TCP-Traceroute (Port %d) ---", "--- Traceroute TCP (porta %d) ---"},
	"rpt.minecraft":    {"--- Serveur Minecraft ---", "--- Minecraft Server ---", "--- Servidor Minecraft ---", "--- Minecraft-Server ---", "--- Servidor Minecraft ---"},
	"rpt.footer":       {"=== Fin du rapport ===", "=== End of report ===", "=== Fin del informe ===", "=== Ende des Berichts ===", "=== Fim do relatório ==="},
	"rpt.skipUser":     {"Statut : SKIP (passe par l'utilisateur)", "Status: SKIP (skipped by user)", "Estado: SKIP (omitido por el usuario)", "Status: SKIP (vom Benutzer übersprungen)", "Estado: SKIP (ignorado pelo utilizador)"},
	"rpt.skipIp":       {"Statut : SKIP (adresse IP directe)", "Status: SKIP (direct IP address)", "Estado: SKIP (dirección IP directa)", "Status: SKIP (direkte IP-Adresse)", "Estado: SKIP (endereço IP direto)"},
	"rpt.ok":           {"Statut : OK (%d ms)", "Status: OK (%d ms)", "Estado: OK (%d ms)", "Status: OK (%d ms)", "Estado: OK (%d ms)"},
	"rpt.okDuration":   {"Statut   : OK (%d ms)", "Status  : OK (%d ms)", "Estado  : OK (%d ms)", "Status  : OK (%d ms)", "Estado  : OK (%d ms)"},
	"rpt.fail":         {"Statut : ECHEC", "Status: FAILED", "Estado: FALLO", "Status: FEHLGESCHLAGEN", "Estado: FALHA"},
	"rpt.failDur":      {"Statut : ECHEC (%d ms)", "Status: FAILED (%d ms)", "Estado: FALLO (%d ms)", "Status: FEHLGESCHLAGEN (%d ms)", "Estado: FALHA (%d ms)"},
	"rpt.error":        {"Erreur : %s", "Error : %s", "Error : %s", "Fehler: %s", "Erro  : %s"},
	"rpt.pingOk":       {"Statut  : OK (%d/%d paquets recus)", "Status : OK (%d/%d packets received)", "Estado : OK (%d/%d paquetes recibidos)", "Status : OK (%d/%d Pakete empfangen)", "Estado : OK (%d/%d pacotes recebidos)"},
	"rpt.latency":      {"Latence : min=%.1fms, moy=%.1fms, max=%.1fms", "Latency: min=%.1fms, avg=%.1fms, max=%.1fms", "Latencia: min=%.1fms, med=%.1fms, max=%.1fms", "Latenz: Min=%.1fms, Avg=%.1fms, Max=%.1fms", "Latência: min=%.1fms, med=%.1fms, max=%.1fms"},
	"rpt.hops":         {"Statut : OK (%d sauts)", "Status: OK (%d hops)", "Estado: OK (%d saltos)", "Status: OK (%d Hops)", "Estado: OK (%d saltos)"},
	"rpt.tcpOk":        {"Statut : OK (%d ms)", "Status: OK (%d ms)", "Estado: OK (%d ms)", "Status: OK (%d ms)", "Estado: OK (%d ms)"},
	"rpt.local":        {"Locale   : %s", "Local  : %s", "Local  : %s", "Lokal  : %s", "Local  : %s"},
	"rpt.remote":       {"Distante : %s", "Remote : %s", "Remoto : %s", "Remote : %s", "Remoto : %s"},
	"rpt.notAvail":     {"Statut : NON DISPONIBLE", "Status: NOT AVAILABLE", "Estado: NO DISPONIBLE", "Status: NICHT VERFÜGBAR", "Estado: NÃO DISPONÍVEL"},
	"rpt.note":         {"Note   : %s", "Note  : %s", "Nota  : %s", "Notiz : %s", "Nota  : %s"},
	"rpt.mcVersion":    {"Version  : %s (protocole %d)", "Version : %s (protocol %d)", "Versión : %s (protocolo %d)", "Version : %s (Protokoll %d)", "Versão  : %s (protocolo %d)"},
	"rpt.mcPlayers":    {"Joueurs  : %d/%d", "Players : %d/%d", "Jugadores: %d/%d", "Spieler : %d/%d", "Jogadores: %d/%d"},
	"rpt.mcMotd":       {"MOTD     : %s", "MOTD    : %s", "MOTD    : %s", "MOTD    : %s", "MOTD    : %s"},
	"rpt.mcOnline":     {"En ligne :", "Online:", "En línea:", "Online:", "Online:"},
	"rpt.mcLatency":    {"Latence protocole : %d ms", "Protocol latency: %d ms", "Latencia protocolo: %d ms", "Protokoll-Latenz: %d ms", "Latência protocolo: %d ms"},

	// ---- Ping platform-specific ----
	"ping.reply":      {"Reponse de %s : temps=%dms", "Reply from %s: time=%dms", "Respuesta de %s: tiempo=%dms", "Antwort von %s: Zeit=%dms", "Resposta de %s: tempo=%dms"},
	"ping.timeout":    {"Delai d'attente depasse", "Request timed out", "Tiempo de espera agotado", "Zeitüberschreitung", "Tempo limite excedido"},
	"ping.unreachable": {"Hote %s injoignable", "Host %s unreachable", "Host %s inalcanzable", "Host %s nicht erreichbar", "Host %s inacessível"},
	"ping.ttlExpired": {"TTL expire depuis %s", "TTL expired from %s", "TTL expirado desde %s", "TTL abgelaufen von %s", "TTL expirado de %s"},
	"ping.icmpErr":    {"Erreur ICMP (code=%d)", "ICMP error (code=%d)", "Error ICMP (código=%d)", "ICMP-Fehler (Code=%d)", "Erro ICMP (código=%d)"},
	"ping.replySeq":   {"Reponse de %s : seq=%d temps=%.1fms", "Reply from %s: seq=%d time=%.1fms", "Respuesta de %s: seq=%d tiempo=%.1fms", "Antwort von %s: Seq=%d Zeit=%.1fms", "Resposta de %s: seq=%d tempo=%.1fms"},
	"ping.sendErr":    {"Erreur envoi : %s", "Send error: %s", "Error de envío: %s", "Sendefehler: %s", "Erro de envio: %s"},
	"ping.sockErr":    {"Impossible d'ouvrir un socket ICMP : %s", "Cannot open ICMP socket: %s", "No se puede abrir socket ICMP: %s", "ICMP-Socket kann nicht geöffnet werden: %s", "Impossível abrir socket ICMP: %s"},
	"ping.resolveErr": {"Resolution echouee : %v", "Resolution failed: %v", "Resolución fallida: %v", "Auflösung fehlgeschlagen: %v", "Resolução falhou: %v"},
	"ping.ipv4Only":   {"Seul IPv4 est supporte pour le ping natif Windows", "Only IPv4 is supported for native Windows ping", "Solo IPv4 es soportado para ping nativo Windows", "Nur IPv4 wird für nativen Windows-Ping unterstützt", "Apenas IPv4 é suportado para ping nativo Windows"},
}
