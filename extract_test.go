package main

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

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
