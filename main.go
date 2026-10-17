package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GitHub config repository details
const (
	defaultRepo    = "jmcampbell72/agents"
	defaultVersion = "main"
)

// In-memory file storage (replaces embed.FS)
type memFS struct {
	files map[string][]byte
}

var configFS *memFS

type action string

const (
	actNew       action = "NEW"
	actOverwrite action = "OVERWRITE"
	actSame      action = "SAME"
)

type item struct {
	rel  string // path relative to claudeconfig, slash-separated
	dest string // absolute destination path
	data []byte
	act  action
}

func main() {
	contents := flag.Bool("contents", false, "list the files to be installed")
	dryrun := flag.Bool("dryrun", false, "show how files would be installed, without writing")
	install := flag.Bool("install", false, "install the files")
	target := flag.String("target", "", "install directory (default ~/.claude)")
	repo := flag.String("repo", defaultRepo, "GitHub repo (owner/repo)")
	version := flag.String("version", defaultVersion, "version tag or branch (default main)")
	flag.Usage = usage
	flag.Parse()

	n := 0
	for _, b := range []bool{*contents, *dryrun, *install} {
		if b {
			n++
		}
	}
	if n != 1 {
		usage()
		os.Exit(2)
	}

	// Download config from GitHub
	fmt.Fprintf(os.Stderr, "Downloading config from https://github.com/%s@%s...\n", *repo, *version)
	var err error
	configFS, err = downloadConfig(*repo, *version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	switch {
	case *contents:
		err = runContents()
	case *dryrun:
		err = runInstall(*target, false)
	case *install:
		err = runInstall(*target, true)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: claude-install <-contents | -dryrun | -install> [options]

Actions (pick one):
  -contents   list the files to be installed
  -dryrun     show how the files would be installed (writes nothing)
  -install    install the files into ~/.claude

Options:
  -target     install directory (default ~/.claude)
  -repo       GitHub repo owner/repo (default jmcampbell72/agents)
  -version    version tag or branch (default main)

Examples:
  claude-install -contents
  claude-install -dryrun -version v1.0.0
  claude-install -install -target ~/.claude-dev`)
}

func collect() ([]string, error) {
	var files []string
	for rel := range configFS.files {
		files = append(files, rel)
	}
	sort.Strings(files)
	return files, nil
}

func targetRoot(override string) (string, error) {
	if override != "" {
		return filepath.Abs(override)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude"), nil
}

// downloadConfig fetches the config archive from GitHub and extracts it to memory.
// It returns a memFS containing all files from agents/, skills/, and rules/ directories.
func downloadConfig(repo, version string) (*memFS, error) {
	// Build GitHub archive URL
	// Format: https://github.com/owner/repo/archive/refs/tags/v1.0.0.zip
	// or for branches: https://github.com/owner/repo/archive/refs/heads/main.zip
	var url string
	// Versions starting with 'v' are tags; common branch names are branches
	isTag := strings.HasPrefix(version, "v") && len(version) > 1 && version[1] >= '0' && version[1] <= '9'
	if isTag {
		url = fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.zip", repo, version)
	} else {
		// Treat as branch (main, develop, feature/x, etc.)
		url = fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", repo, version)
	}

	// Download the zip
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}

	// Read the entire response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse as zip
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("parse zip: %w", err)
	}

	mfs := &memFS{files: make(map[string][]byte)}

	// Extract files from agents/, skills/, rules/ directories (under configs/)
	for _, f := range zr.File {
		// Skip directories and DS_Store
		if f.FileInfo().IsDir() || strings.Contains(f.Name, ".DS_Store") {
			continue
		}

		// Extract relative path after the repo-version directory
		// GitHub archives extract to "agents-main/" or "agents-v1.0.0/"
		// We want files from configs/agents/, configs/skills/, configs/rules/
		var rel string
		parts := strings.Split(f.Name, "/")
		if len(parts) > 2 {
			// Skip the first part (directory created by GitHub archive)
			// and keep path starting from "configs/..."
			rest := strings.Join(parts[1:], "/")
			if strings.HasPrefix(rest, "configs/") {
				rel = rest[8:] // Strip "configs/" prefix
			}
		}

		if rel == "" {
			continue
		}

		// Only include agents/, skills/, and rules/ files
		if !strings.HasPrefix(rel, "agents/") && !strings.HasPrefix(rel, "skills/") && !strings.HasPrefix(rel, "rules/") {
			continue
		}

		// Read file content
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f.Name, err)
		}

		mfs.files[rel] = data
	}

	if len(mfs.files) == 0 {
		return nil, fmt.Errorf("no config files found in archive")
	}

	return mfs, nil
}

func runContents() error {
	files, err := collect()
	if err != nil {
		return err
	}
	for _, rel := range files {
		data := configFS.files[rel]
		fmt.Printf("%-40s %6d bytes\n", rel, len(data))
	}
	fmt.Printf("%d files\n", len(files))
	return nil
}

func plan(root string) ([]item, error) {
	files, err := collect()
	if err != nil {
		return nil, err
	}
	items := make([]item, 0, len(files))
	for _, rel := range files {
		data := configFS.files[rel]
		it := item{rel: rel, dest: filepath.Join(root, filepath.FromSlash(rel)), data: data, act: actNew}
		if existing, err := os.ReadFile(it.dest); err == nil {
			if bytes.Equal(existing, data) {
				it.act = actSame
			} else {
				it.act = actOverwrite
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read %s: %w", it.dest, err)
		}
		items = append(items, it)
	}
	return items, nil
}

func runInstall(targetFlag string, write bool) error {
	root, err := targetRoot(targetFlag)
	if err != nil {
		return err
	}
	items, err := plan(root)
	if err != nil {
		return err
	}

	counts := map[action]int{}
	for _, it := range items {
		counts[it.act]++
		fmt.Printf("%-9s %s -> %s\n", it.act, it.rel, it.dest)
		if !write || it.act == actSame {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(it.dest), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(it.dest), err)
		}
		if it.act == actOverwrite {
			old, err := os.ReadFile(it.dest)
			if err != nil {
				return err
			}
			if err := os.WriteFile(it.dest+".bak", old, 0o644); err != nil {
				return fmt.Errorf("backup %s: %w", it.dest, err)
			}
		}
		if err := os.WriteFile(it.dest, it.data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", it.dest, err)
		}
	}

	verb := "would install"
	if write {
		verb = "installed"
	}
	fmt.Printf("\n%s into %s: %d new, %d overwritten, %d unchanged\n",
		verb, root, counts[actNew], counts[actOverwrite], counts[actSame])
	return nil
}
