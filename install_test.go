package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureStdout runs f and returns what it printed to stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	defer func() { os.Stdout = old }()
	f()
	w.Close()
	os.Stdout = old
	return <-done
}

func testFS() *memFS {
	return &memFS{files: map[string][]byte{
		"agents/api-dev.md":                []byte("agent v2"),
		"skills/testing-patterns/SKILL.md": []byte("skill v2"),
	}}
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRunInstallDryRunWritesNothing(t *testing.T) {
	root := t.TempDir()
	out := captureStdout(t, func() {
		if err := runInstall(testFS(), root, false); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "would install") || !strings.Contains(out, "2 new") {
		t.Errorf("dry run output %q, want it to report 2 new files it would install", out)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("dry run wrote %d entries into the target", len(entries))
	}
}

func TestRunInstallNewSameOverwrite(t *testing.T) {
	root := t.TempDir()
	agent := filepath.Join(root, "agents", "api-dev.md")
	skill := filepath.Join(root, "skills", "testing-patterns", "SKILL.md")

	// First install: everything is new.
	out := captureStdout(t, func() {
		if err := runInstall(testFS(), root, true); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "2 new, 0 overwritten, 0 unchanged") {
		t.Errorf("first install output %q", out)
	}
	if got := readFile(t, agent); got != "agent v2" {
		t.Errorf("installed agent = %q", got)
	}
	if got := readFile(t, skill); got != "skill v2" {
		t.Errorf("installed skill = %q", got)
	}

	// Same files again: nothing to do, and no backups.
	out = captureStdout(t, func() {
		if err := runInstall(testFS(), root, true); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "0 new, 0 overwritten, 2 unchanged") {
		t.Errorf("second install output %q", out)
	}
	if _, err := os.Stat(agent + ".bak"); !os.IsNotExist(err) {
		t.Errorf("unchanged file got a backup (err %v)", err)
	}

	// A locally edited file is overwritten and the edit kept as .bak.
	if err := os.WriteFile(agent, []byte("my local edit"), 0o644); err != nil {
		t.Fatal(err)
	}
	out = captureStdout(t, func() {
		if err := runInstall(testFS(), root, true); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "0 new, 1 overwritten, 1 unchanged") {
		t.Errorf("third install output %q", out)
	}
	if got := readFile(t, agent); got != "agent v2" {
		t.Errorf("overwritten agent = %q", got)
	}
	if got := readFile(t, agent+".bak"); got != "my local edit" {
		t.Errorf("backup = %q, want the local edit", got)
	}
}

func TestRunContents(t *testing.T) {
	out := captureStdout(t, func() {
		if err := runContents(testFS()); err != nil {
			t.Fatal(err)
		}
	})
	for _, want := range []string{"agents/api-dev.md", "skills/testing-patterns/SKILL.md", "2 files"} {
		if !strings.Contains(out, want) {
			t.Errorf("contents output %q, want it to contain %q", out, want)
		}
	}
}

func TestTargetRoot(t *testing.T) {
	dir := t.TempDir()
	if got, err := targetRoot(dir); err != nil || got != dir {
		t.Errorf("targetRoot(%q) = %q, %v", dir, got, err)
	}
	got, err := targetRoot("")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != ".claude" || !filepath.IsAbs(got) {
		t.Errorf("default targetRoot = %q, want an absolute path ending in .claude", got)
	}
}

func TestCountActions(t *testing.T) {
	tests := []struct {
		set  []bool
		want int
	}{
		{nil, 0},
		{[]bool{false, false}, 0},
		{[]bool{true, false, false}, 1},
		{[]bool{false, true, true}, 2},
	}
	for _, tt := range tests {
		if got := countActions(tt.set...); got != tt.want {
			t.Errorf("countActions(%v) = %d, want %d", tt.set, got, tt.want)
		}
	}
}
