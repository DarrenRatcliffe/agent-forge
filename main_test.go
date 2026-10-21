package main

import (
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
