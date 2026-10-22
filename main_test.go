package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDownloadConfigTimesOut(t *testing.T) {
	// A server that accepts the request but never sends a response.
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer srv.Close()
	defer close(release)

	oldURL, oldClient := githubURL, downloadClient
	githubURL, downloadClient = srv.URL, &http.Client{Timeout: 100 * time.Millisecond}
	defer func() { githubURL, downloadClient = oldURL, oldClient }()

	done := make(chan error, 1)
	go func() {
		_, err := downloadConfig("owner/repo", "main")
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("downloadConfig succeeded against a server that never responds")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("downloadConfig did not time out")
	}
}

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

func TestServeMuxListsAreJSON(t *testing.T) {
	full := httptest.NewServer(newServeMux(newTestConfig(t)))
	defer full.Close()
	empty := httptest.NewServer(newServeMux(t.TempDir()))
	defer empty.Close()

	tests := []struct {
		url  string
		key  string
		want []string
	}{
		{full.URL + "/agents", "agents", []string{"api-dev"}},
		{full.URL + "/skills", "skills", []string{"testing-patterns"}},
		{empty.URL + "/agents", "agents", []string{}},
		{empty.URL + "/skills", "skills", []string{}},
	}
	for _, tt := range tests {
		resp, body := get(t, tt.url)
		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("GET %s: Content-Type %q, want application/json", tt.url, ct)
		}
		var got map[string]json.RawMessage
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Errorf("GET %s: body %q is not valid JSON: %v", tt.url, body, err)
			continue
		}
		var names []string
		if err := json.Unmarshal(got[tt.key], &names); err != nil || names == nil {
			t.Errorf("GET %s: %q = %s, want a JSON array", tt.url, tt.key, got[tt.key])
			continue
		}
		if strings.Join(names, ",") != strings.Join(tt.want, ",") {
			t.Errorf("GET %s: %q = %v, want %v", tt.url, tt.key, names, tt.want)
		}
		var count int
		if err := json.Unmarshal(got["count"], &count); err != nil || count != len(tt.want) {
			t.Errorf("GET %s: count = %s, want %d", tt.url, got["count"], len(tt.want))
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

func TestQueryActionsUseTarget(t *testing.T) {
	// Keep the printed agent and skill text out of the test output.
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()
	oldStdout := os.Stdout
	os.Stdout = devNull
	defer func() { os.Stdout = oldStdout }()

	target := newTestConfig(t)
	if err := runAgent(target, "api-dev"); err != nil {
		t.Errorf("runAgent with -target: %v", err)
	}
	if err := runSkill(target, "testing-patterns"); err != nil {
		t.Errorf("runSkill with -target: %v", err)
	}
	if err := runAgent(t.TempDir(), "api-dev"); err == nil {
		t.Error("runAgent found api-dev in an empty -target")
	}
}
