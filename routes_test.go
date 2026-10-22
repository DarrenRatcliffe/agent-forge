package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
