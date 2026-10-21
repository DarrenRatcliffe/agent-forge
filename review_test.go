package main

import "testing"

func TestReviewerFor(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"main.go", "go-reviewer"},
		{"app/models.py", "python-reviewer"},
		{"SCRIPT.PY", "python-reviewer"},
		{"src/lib.rs", "rust-reviewer"},
		{"server.js", "node-reviewer"},
		{"server.ts", "node-reviewer"},
		{"config.cjs", "node-reviewer"},
		{"components/Button.tsx", "frontend-reviewer"},
		{"Button.jsx", "frontend-reviewer"},
		{"styles/site.scss", "frontend-reviewer"},
		{"migrations/001_init.sql", "db-reviewer"},
		{"infra/main.tf", "devops-reviewer"},
		{"infra/main.bicep", "devops-reviewer"},
		{"Dockerfile", "devops-reviewer"},
		{"docker/Dockerfile.prod", "devops-reviewer"},
		{"docker-compose.yml", "devops-reviewer"},
		{"compose.yaml", "devops-reviewer"},
		{".github/workflows/build.yml", "devops-reviewer"},
		{`C:\repo\.github\workflows\ci.yaml`, "devops-reviewer"},
		{"azure-pipelines.yml", "devops-reviewer"},
		{"compose.go", "go-reviewer"},
		{"config.yml", "code-reviewer"},
		{"Main.java", "code-reviewer"},
		{"README", "code-reviewer"},
	}
	for _, tt := range tests {
		if got := reviewerFor(tt.file); got != tt.want {
			t.Errorf("reviewerFor(%q) = %q, want %q", tt.file, got, tt.want)
		}
	}
}

func TestArchiveURL(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"main", "https://github.com/owner/repo/archive/main.zip"},
		{"v1.0.1", "https://github.com/owner/repo/archive/v1.0.1.zip"},
		{"feature/x", "https://github.com/owner/repo/archive/feature/x.zip"},
		{"v2-wip", "https://github.com/owner/repo/archive/v2-wip.zip"},
	}
	for _, tt := range tests {
		if got := archiveURL("owner/repo", tt.version); got != tt.want {
			t.Errorf("archiveURL(%q) = %q, want %q", tt.version, got, tt.want)
		}
	}
}
