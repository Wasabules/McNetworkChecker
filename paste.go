package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PasteResult struct {
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

// UploadPaste uploads text content to the chosen paste service.
// Returns the paste URL or an error.
func (a *App) UploadPaste(service, content string) PasteResult {
	switch service {
	case "mclo.gs":
		return pasteMcLogs(content)
	case "pastes.dev":
		return pastePastesDev(content)
	case "dpaste.org":
		return pasteDpasteOrg(content)
	case "paste.rs":
		return pastePasteRs(content)
	case "0x0.st":
		return paste0x0(content)
	case "dpaste.com":
		return pasteDpasteCom(content)
	default:
		return PasteResult{Error: "Unknown service: " + service}
	}
}

// GetPasteServices returns the list of available paste services for the frontend.
func (a *App) GetPasteServices() []map[string]string {
	return []map[string]string{
		{"id": "mclo.gs", "name": "mclo.gs", "desc": "Minecraft Log Sharing"},
		{"id": "pastes.dev", "name": "pastes.dev", "desc": "Lucko Paste (MC)"},
		{"id": "dpaste.org", "name": "dpaste.org", "desc": "dpaste.org"},
		{"id": "paste.rs", "name": "paste.rs", "desc": "Rocket Paste"},
		{"id": "0x0.st", "name": "0x0.st", "desc": "Null Pointer"},
		{"id": "dpaste.com", "name": "dpaste.com", "desc": "dpaste.com"},
	}
}

// ---- Service implementations ----

// mclo.gs — POST form content=... → JSON {success, url}
func pasteMcLogs(content string) PasteResult {
	resp, err := httpClient.PostForm("https://api.mclo.gs/1/log", url.Values{"content": {content}})
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()
	var r struct {
		Success bool   `json:"success"`
		URL     string `json:"url"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return PasteResult{Error: "JSON decode: " + err.Error()}
	}
	if !r.Success {
		return PasteResult{Error: "mclo.gs: " + r.Error}
	}
	return PasteResult{URL: r.URL}
}

// pastes.dev — POST raw body → JSON {key} → https://pastes.dev/{key}
func pastePastesDev(content string) PasteResult {
	req, err := http.NewRequest("POST", "https://api.pastes.dev/post", strings.NewReader(content))
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("User-Agent", "McNetworkChecker (github.com/Wasabules/McNetworkChecker)")
	resp, err := httpClient.Do(req)
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()
	var r struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return PasteResult{Error: "JSON decode: " + err.Error()}
	}
	if r.Key == "" {
		return PasteResult{Error: "pastes.dev: empty key"}
	}
	return PasteResult{URL: "https://pastes.dev/" + r.Key}
}

// dpaste.org — POST form content=...&format=url → plain URL
func pasteDpasteOrg(content string) PasteResult {
	resp, err := httpClient.PostForm("https://dpaste.org/api/", url.Values{
		"content": {content}, "format": {"url"}, "expires": {"2592000"},
	})
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	u := strings.TrimSpace(string(body))
	if !strings.HasPrefix(u, "http") {
		return PasteResult{Error: "dpaste.org: " + u}
	}
	return PasteResult{URL: u}
}

// paste.rs — POST raw body → plain URL in response body
func pastePasteRs(content string) PasteResult {
	req, err := http.NewRequest("POST", "https://paste.rs/", strings.NewReader(content))
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	u := strings.TrimSpace(string(body))
	if !strings.HasPrefix(u, "http") {
		return PasteResult{Error: "paste.rs: " + u}
	}
	return PasteResult{URL: u}
}

// 0x0.st — POST multipart file=@content → plain URL
func paste0x0(content string) PasteResult {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "diagnostic.txt")
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	fw.Write([]byte(content))
	w.Close()

	req, err := http.NewRequest("POST", "https://0x0.st", &buf)
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", "McNetworkChecker")
	resp, err := httpClient.Do(req)
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	u := strings.TrimSpace(string(body))
	if !strings.HasPrefix(u, "http") {
		return PasteResult{Error: fmt.Sprintf("0x0.st (%d): %s", resp.StatusCode, u)}
	}
	return PasteResult{URL: u}
}

// dpaste.com — POST form content=... → URL in Location header or body
func pasteDpasteCom(content string) PasteResult {
	data := url.Values{"content": {content}, "syntax": {"text"}, "expiry_days": {"30"}}
	req, err := http.NewRequest("POST", "https://dpaste.com/api/v2/", strings.NewReader(data.Encode()))
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "McNetworkChecker")

	// Don't follow redirects — we want the Location header
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	if err != nil {
		return PasteResult{Error: err.Error()}
	}
	defer resp.Body.Close()

	if loc := resp.Header.Get("Location"); loc != "" {
		return PasteResult{URL: strings.TrimSpace(loc)}
	}
	body, _ := io.ReadAll(resp.Body)
	u := strings.TrimSpace(string(body))
	if strings.HasPrefix(u, "http") {
		return PasteResult{URL: u}
	}
	return PasteResult{Error: fmt.Sprintf("dpaste.com (%d): %s", resp.StatusCode, u)}
}
