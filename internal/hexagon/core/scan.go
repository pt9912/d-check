package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// SkipDirs werden immer übersprungen (DC-FA-SCAN-001,
// spec/spezifikation.md §3 SKIP_DIRS).
var SkipDirs = map[string]bool{
	".git": true, "node_modules": true, "build": true, "target": true,
	"dist": true, "vendor": true, ".venv": true, "__pycache__": true,
	".idea": true, ".vscode": true,
}

// DefaultRoots sind die optionalen Default-Scan-Wurzeln
// (DC-FA-SCAN-001).
var DefaultRoots = []string{"docs", "spec"}

// DiscoverFiles ermittelt die zu prüfenden Markdown-Dateien in
// deterministischer Reihenfolge (Pfade bytewise sortiert, DC-QA-02).
// explicitRoots stammen aus der Konfiguration und müssen existieren;
// ohne explizite Wurzeln gelten DefaultRoots (optional) plus *.md der
// Repo-Wurzel (DC-FA-SCAN-001).
func DiscoverFiles(fsys driven.Filesystem, explicitRoots, ignore []string) ([]string, error) {
	var files []string
	if explicitRoots != nil {
		for _, r := range explicitRoots {
			r = strings.Trim(r, "/")
			// "." steht für die Repo-Wurzel (gesamter Baum)
			if r == "." || r == "" {
				if err := walkMarkdown(fsys, "", ignore, &files); err != nil {
					return nil, err
				}
				continue
			}
			kind, err := fsys.Kind(r)
			if err != nil {
				return nil, err
			}
			if kind != driven.KindDir {
				return nil, fmt.Errorf("konfigurierte Scan-Wurzel existiert nicht: %s", r)
			}
			if err := walkMarkdown(fsys, r, ignore, &files); err != nil {
				return nil, err
			}
		}
	} else {
		for _, r := range DefaultRoots {
			kind, err := fsys.Kind(r)
			if err != nil {
				return nil, err
			}
			if kind != driven.KindDir {
				continue // Default-Wurzeln sind optional
			}
			if err := walkMarkdown(fsys, r, ignore, &files); err != nil {
				return nil, err
			}
		}
		// *.md direkt in der Repo-Wurzel
		entries, err := fsys.List("")
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.Kind == driven.KindFile && strings.HasSuffix(e.Name, ".md") && !ignored(e.Name, ignore) {
				files = append(files, e.Name)
			}
		}
	}
	sort.Strings(files)
	return files, nil
}

func walkMarkdown(fsys driven.Filesystem, dir string, ignore []string, out *[]string) error {
	entries, err := fsys.List(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		rel := e.Name
		if dir != "" {
			rel = dir + "/" + e.Name
		}
		switch e.Kind {
		case driven.KindDir:
			if SkipDirs[e.Name] {
				continue
			}
			if err := walkMarkdown(fsys, rel, ignore, out); err != nil {
				return err
			}
		case driven.KindFile:
			if strings.HasSuffix(e.Name, ".md") && !ignored(rel, ignore) {
				*out = append(*out, rel)
			}
		}
	}
	return nil
}

func ignored(rel string, ignore []string) bool {
	for _, pat := range ignore {
		if matchGlob(pat, rel) {
			return true
		}
	}
	return false
}
