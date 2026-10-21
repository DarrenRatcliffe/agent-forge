package main

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestConfig builds a fake ~/.claude with one agent, one skill and a
// secret file next to them that the server must never return.
func newTestConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"agents/api-dev.md":                    "agent api-dev",
		"skills/testing-patterns/SKILL.md":     "skill testing-patterns",
		"skills/testing-patterns/checklist.md": "checklist testing-patterns",
		"skills/testing-patterns/notes.txt":    "not served",
		".credentials.json":                    "SECRET",
		"secret.md":                            "SECRET",
	}
	for rel, body := range files {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestServeMuxNormalRequests(t *testing.T) {
	srv := httptest.NewServer(newServeMux(newTestConfig(t)))
	defer srv.Close()

	tests := []struct {
		path string
		want string
	}{
		{"/health", `"ok"`},
		{"/agents", "api-dev"},
		{"/agent/api-dev", "agent api-dev"},
		{"/skills", "testing-patterns"},
		{"/skill/testing-patterns", "skill testing-patterns"},
		{"/skill/testing-patterns/SKILL.md", "skill testing-patterns"},
		{"/skill/testing-patterns/checklist.md", "checklist testing-patterns"},
	}
	for _, tt := range tests {
		resp, body := get(t, srv.URL+tt.path)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s: status %d, want 200", tt.path, resp.StatusCode)
		}
		if !strings.Contains(body, tt.want) {
			t.Errorf("GET %s: body %q, want it to contain %q", tt.path, body, tt.want)
		}
	}
}

func TestServeMuxRejectsTraversal(t *testing.T) {
	srv := httptest.NewServer(newServeMux(newTestConfig(t)))
	defer srv.Close()

	paths := []string{
		"/skill/..%2F.credentials.json",
		"/skill/%2E%2E/.credentials.json",
		"/skill/%2e%2e%2f.credentials.json",
		"/skill/testing-patterns/..%2F..%2F.credentials.json",
		"/skill/testing-patterns/notes.txt",
		"/skill/testing-patterns/SKILL.md/extra",
		"/agent/..%2Fsecret",
		"/agent/%2E%2E/secret",
		"/agent/API-DEV",
		"/agent/-api",
	}
	for _, p := range paths {
		resp, body := get(t, srv.URL+p)
		if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: status %d, want 400 or 404", p, resp.StatusCode)
		}
		if strings.Contains(body, "SECRET") {
			t.Errorf("GET %s: leaked file outside the config directories", p)
		}
	}
}

func get(t *testing.T, url string) (*http.Response, string) {
	t.Helper()
	// Send the path as written; don't let the client clean or re-encode it.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var sb strings.Builder
	if _, err := io.Copy(&sb, resp.Body); err != nil {
		t.Fatal(err)
	}
	return resp, sb.String()
}

func TestRootHelpUsesRequestHost(t *testing.T) {
	srv := httptest.NewServer(newServeMux(newTestConfig(t)))
	defer srv.Close()

	_, body := get(t, srv.URL+"/")
	if want := "curl " + srv.URL + "/agents"; !strings.Contains(body, want) {
		t.Errorf("GET /: body %q, want it to contain %q", body, want)
	}
}

func TestBaseURL(t *testing.T) {
	tests := []struct {
		addr string
		want string
	}{
		{"127.0.0.1:8081", "http://127.0.0.1:8081"},
		{"127.0.0.1:9090", "http://127.0.0.1:9090"},
		{"0.0.0.0:8081", "http://localhost:8081"},
		{":9090", "http://localhost:9090"},
		{"[::]:8081", "http://localhost:8081"},
		{"[::1]:8081", "http://[::1]:8081"},
	}
	for _, tt := range tests {
		if got := baseURL(tt.addr); got != tt.want {
			t.Errorf("baseURL(%q) = %q, want %q", tt.addr, got, tt.want)
		}
	}
}

// zipOf builds a zip archive holding the given files.
func zipOf(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(w, body); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractConfig(t *testing.T) {
	cfg, err := extractConfig(zipOf(t, map[string]string{
		"agent-forge-main/configs/agents/api-dev.md":                "agent",
		"agent-forge-main/configs/skills/testing-patterns/SKILL.md": "skill",
		"agent-forge-main/configs/rules/go-dev.md":                  "rule",
		"agent-forge-main/configs/agents/.DS_Store":                 "junk",
		"agent-forge-main/configs/other/notes.md":                   "not config",
		"agent-forge-main/README.md":                                "not config",
		"README.md":                                                 "not config",
	}))
	if err != nil {
		t.Fatal(err)
	}
	want := "agents/api-dev.md,rules/go-dev.md,skills/testing-patterns/SKILL.md"
	if got := strings.Join(cfg.names(), ","); got != want {
		t.Errorf("names = %s, want %s", got, want)
	}
	if got := string(cfg.files["agents/api-dev.md"]); got != "agent" {
		t.Errorf("agents/api-dev.md = %q, want %q", got, "agent")
	}
}

func TestExtractConfigRejectsUnsafePaths(t *testing.T) {
	for _, name := range []string{
		"agent-forge-main/configs/agents/../../../.bashrc",
		"agent-forge-main/configs/skills/../../evil.md",
	} {
		if _, err := extractConfig(zipOf(t, map[string]string{name: "x"})); err == nil {
			t.Errorf("extractConfig accepted %s", name)
		}
	}
}

func TestServeMuxUnknownPathIsNotFound(t *testing.T) {
	srv := httptest.NewServer(newServeMux(newTestConfig(t)))
	defer srv.Close()

	if resp, _ := get(t, srv.URL+"/"); resp.StatusCode != http.StatusOK {
		t.Errorf("GET /: status %d, want 200", resp.StatusCode)
	}
	for _, p := range []string{"/agnets", "/nope", "/health/extra"} {
		if resp, _ := get(t, srv.URL+p); resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: status %d, want 404", p, resp.StatusCode)
		}
	}
}
