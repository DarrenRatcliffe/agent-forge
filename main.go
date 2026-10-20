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
	defaultRepo    = "ic-consulting/agent-forge"
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
	agent := flag.String("agent", "", "get guidance from an agent (e.g., api-dev, go-reviewer)")
	skill := flag.String("skill", "", "get skill details and checklist (e.g., testing-patterns)")
	serve := flag.Bool("serve", false, "start HTTP server for agents/skills on :8081")
	reviewFile := flag.String("review", "", "get review guidance for a file (auto-detects reviewer by extension)")
	target := flag.String("target", "", "install directory (default ~/.claude)")
	repo := flag.String("repo", defaultRepo, "GitHub repo (owner/repo)")
	version := flag.String("version", defaultVersion, "version tag or branch (default main)")
	flag.Usage = usage
	flag.Parse()

	// Count install/update actions (require download from GitHub)
	downloadActions := 0
	for _, b := range []bool{*contents, *dryrun, *install} {
		if b {
			downloadActions++
		}
	}

	// Local query actions (don't require GitHub download)
	localActions := 0
	if *agent != "" || *skill != "" || *serve || *reviewFile != "" {
		localActions++
	}

	// Must have exactly one action
	if downloadActions+localActions != 1 {
		usage()
		os.Exit(2)
	}

	var err error

	// Handle local query actions first (no GitHub download needed)
	switch {
	case *agent != "":
		err = runAgent(*agent)
	case *skill != "":
		err = runSkill(*skill)
	case *serve:
		err = runServer()
	case *reviewFile != "":
		err = runReview(*reviewFile)
	default:
		// Handle download/install actions
		fmt.Fprintf(os.Stderr, "Downloading config from https://github.com/%s@%s...\n", *repo, *version)
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
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: agent-install <ACTION> [options]

INSTALL/UPDATE ACTIONS (require GitHub download):
  -contents   list the files to be installed
  -dryrun     show how the files would be installed (writes nothing)
  -install    install the files into ~/.claude

LOCAL QUERY ACTIONS (use installed agents/skills):
  -agent NAME        get guidance from an agent (e.g., api-dev, go-reviewer)
  -skill NAME        get skill details and checklist (e.g., testing-patterns)
  -serve             start HTTP server for agents/skills on http://localhost:8081
  -review FILE       get review guidance for a file (auto-detects reviewer by extension)

OPTIONS:
  -target     install directory (default ~/.claude)
  -repo       GitHub repo owner/repo (default ic-consulting/agents)
  -version    version tag or branch (default main)

EXAMPLES:
  # Install agents
  agent-install -install
  agent-install -dryrun -version v1.0.0

  # Query agents (after install)
  agent-install -agent api-dev
  agent-install -agent go-reviewer
  agent-install -skill testing-patterns
  agent-install -skill frontend-accessibility
  agent-install -review mycode.py
  agent-install -serve`)
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

// getConfigPath returns the directory where agents/skills/rules are installed
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude"), nil
}

// runAgent loads and displays an agent file
func runAgent(agentName string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	agentFile := filepath.Join(configPath, "agents", agentName+".md")
	data, err := os.ReadFile(agentFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("agent not found: %s (install with: agent-install -install)", agentName)
		}
		return err
	}

	fmt.Println(string(data))
	return nil
}

// runSkill loads and displays a skill file with its checklist
func runSkill(skillName string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	skillFile := filepath.Join(configPath, "skills", skillName, "SKILL.md")
	data, err := os.ReadFile(skillFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill not found: %s (install with: agent-install -install)", skillName)
		}
		return err
	}

	fmt.Println("=== SKILL ===")
	fmt.Println(string(data))

	checklistFile := filepath.Join(configPath, "skills", skillName, "checklist.md")
	if checklistData, err := os.ReadFile(checklistFile); err == nil {
		fmt.Println("\n=== CHECKLIST ===")
		fmt.Println(string(checklistData))
	}

	return nil
}

// runServer starts an HTTP server to serve agents and skills
func runServer() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// List all agents
	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		files, _ := filepath.Glob(filepath.Join(configPath, "agents", "*.md"))
		var agents []string
		for _, f := range files {
			name := filepath.Base(f)
			agents = append(agents, name[:len(name)-3])
		}
		sort.Strings(agents)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"agents":%v,"count":%d}`+"\n", agents, len(agents))
	})

	// Get specific agent
	http.HandleFunc("/agent/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[7:] // Strip "/agent/"
		if name == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "agent name required")
			return
		}
		data, err := os.ReadFile(filepath.Join(configPath, "agents", name+".md"))
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "agent not found")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})

	// List all skills
	http.HandleFunc("/skills", func(w http.ResponseWriter, r *http.Request) {
		entries, _ := filepath.Glob(filepath.Join(configPath, "skills", "*"))
		var skills []string
		for _, e := range entries {
			if info, err := os.Stat(e); err == nil && info.IsDir() {
				skills = append(skills, filepath.Base(e))
			}
		}
		sort.Strings(skills)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"skills":%v,"count":%d}`+"\n", skills, len(skills))
	})

	// Get skill guidance
	http.HandleFunc("/skill/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path[7:], "/") // Strip "/skill/"
		if len(parts) < 1 || parts[0] == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "skill name required")
			return
		}
		skillName := parts[0]
		requestFile := "SKILL.md"
		if len(parts) > 1 && parts[1] != "" {
			requestFile = parts[1]
		}
		data, err := os.ReadFile(filepath.Join(configPath, "skills", skillName, requestFile))
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "skill or file not found")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Root help
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `agent-forge HTTP API

Endpoints:
  GET /health                    - health check
  GET /agents                    - list all agents (JSON)
  GET /agent/{name}              - get agent details (plain text)
  GET /skills                    - list all skills (JSON)
  GET /skill/{name}              - get skill SKILL.md (plain text)
  GET /skill/{name}/checklist.md - get skill checklist (plain text)

Examples:
  curl http://localhost:8081/agents
  curl http://localhost:8081/agent/api-dev
  curl http://localhost:8081/skill/testing-patterns
  curl http://localhost:8081/skill/testing-patterns/checklist.md
`)
	})

	fmt.Println("Starting agent-forge HTTP server...")
	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("  Health:     http://localhost:8081/health")
	fmt.Println("  Agents:     http://localhost:8081/agents")
	fmt.Println("  Agent:      http://localhost:8081/agent/{name}")
	fmt.Println("  Skills:     http://localhost:8081/skills")
	fmt.Println("  Skill:      http://localhost:8081/skill/{name}")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  curl http://localhost:8081/agent/api-dev")
	fmt.Println("  curl http://localhost:8081/skill/testing-patterns")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("")

	return http.ListenAndServe(":8081", nil)
}

// runReview detects reviewer agent based on file extension
func runReview(filePath string) error {
	ext := filepath.Ext(filePath)
	reviewerAgent := ""

	switch ext {
	case ".go":
		reviewerAgent = "go-reviewer"
	case ".py":
		reviewerAgent = "python-reviewer"
	case ".js", ".ts", ".jsx", ".tsx", ".mjs":
		reviewerAgent = "node-reviewer"
	case ".rs":
		reviewerAgent = "rust-reviewer"
	case ".java":
		reviewerAgent = "code-reviewer"
	default:
		reviewerAgent = "code-reviewer"
	}

	fmt.Printf("📋 Suggested reviewer: %s\n\n", reviewerAgent)
	return runAgent(reviewerAgent)
}
