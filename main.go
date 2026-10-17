package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
)

//go:embed all:claudeconfig
var bundle embed.FS

const bundleRoot = "claudeconfig"

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
	contents := flag.Bool("contents", false, "list the bundled files")
	dryrun := flag.Bool("dryrun", false, "show how files would be installed, without writing")
	install := flag.Bool("install", false, "install the files")
	target := flag.String("target", "", "install directory (default ~/.claude)")
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

	var err error
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
	fmt.Fprintln(os.Stderr, `Usage: claude-install <-contents | -dryrun | -install> [-target dir]

  -contents   list the files bundled in this installer
  -dryrun     show how the files would be installed (writes nothing)
  -install    install the files into ~/.claude
  -target     install directory instead of ~/.claude`)
}

func collect() ([]string, error) {
	var files []string
	err := fs.WalkDir(bundle, bundleRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() == ".DS_Store" {
			return nil
		}
		files = append(files, p[len(bundleRoot)+1:])
		return nil
	})
	sort.Strings(files)
	return files, err
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

func runContents() error {
	files, err := collect()
	if err != nil {
		return err
	}
	for _, rel := range files {
		data, err := bundle.ReadFile(path.Join(bundleRoot, rel))
		if err != nil {
			return err
		}
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
		data, err := bundle.ReadFile(path.Join(bundleRoot, rel))
		if err != nil {
			return nil, err
		}
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
