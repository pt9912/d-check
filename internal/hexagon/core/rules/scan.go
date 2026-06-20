package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// isSkipDir: diese Verzeichnisse werden immer übersprungen
// (DC-FA-SCAN-001, spec/spezifikation.md §3 SKIP_DIRS).
func isSkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "build", "target", "dist",
		"vendor", ".venv", "__pycache__", ".idea", ".vscode",
		".gradle":
		return true
	}
	return false
}

// defaultRoots sind die optionalen Default-Scan-Wurzeln
// (DC-FA-SCAN-001).
func defaultRoots() []string { return []string{"docs", "spec"} }

// DiscoverFiles ermittelt die zu prüfenden Markdown-Dateien in
// deterministischer Reihenfolge (Pfade bytewise sortiert, DC-QA-02).
// explicitRoots stammen aus der Konfiguration und müssen existieren
// ("." = gesamte Repo-Wurzel); ohne explizite Wurzeln gelten die
// Default-Wurzeln (optional) plus *.md der Repo-Wurzel
// (DC-FA-SCAN-001).
func DiscoverFiles(fsys driven.Filesystem, explicitRoots, ignore []string) ([]string, error) {
	var files []string
	var err error
	if explicitRoots != nil {
		err = discoverExplicit(fsys, explicitRoots, ignore, &files)
	} else {
		err = discoverDefaults(fsys, ignore, &files)
	}
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func discoverExplicit(fsys driven.Filesystem, roots, ignore []string, files *[]string) error {
	for _, r := range roots {
		rel, escaped := ResolveConfigPath(r)
		if escaped {
			return fmt.Errorf("konfigurierte Scan-Wurzel verlässt die Repository-Wurzel: %s", r)
		}
		if rel == "" {
			// "." steht für die gesamte Repo-Wurzel
			if err := walkMarkdown(fsys, "", ignore, files); err != nil {
				return err
			}
			continue
		}
		kind, err := fsys.Kind(rel)
		if err != nil {
			return err
		}
		if kind != driven.KindDir {
			return fmt.Errorf("konfigurierte Scan-Wurzel existiert nicht: %s", rel)
		}
		if err := walkMarkdown(fsys, rel, ignore, files); err != nil {
			return err
		}
	}
	return nil
}

func discoverDefaults(fsys driven.Filesystem, ignore []string, files *[]string) error {
	for _, r := range defaultRoots() {
		kind, err := fsys.Kind(r)
		if err != nil {
			return err
		}
		if kind != driven.KindDir {
			continue // Default-Wurzeln sind optional
		}
		if err := walkMarkdown(fsys, r, ignore, files); err != nil {
			return err
		}
	}
	// *.md direkt in der Repo-Wurzel
	entries, err := fsys.List("")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Kind == driven.KindFile && strings.HasSuffix(e.Name, ".md") && !Ignored(e.Name, ignore) {
			*files = append(*files, e.Name)
		}
	}
	return nil
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
			if isSkipDir(e.Name) || dirIgnored(rel, ignore) {
				continue
			}
			if err := walkMarkdown(fsys, rel, ignore, out); err != nil {
				return err
			}
		case driven.KindFile:
			if strings.HasSuffix(e.Name, ".md") && !Ignored(rel, ignore) {
				*out = append(*out, rel)
			}
		}
	}
	return nil
}

// dirIgnored prunt den Verzeichnis-Abstieg: ein Teilbaum, dessen
// Inhalte ein Ignore-Muster vollständig abdeckt (`pfad/**` bzw. ein
// direkt auf das Verzeichnis matchendes Muster), wird nicht betreten —
// unlesbare ignorierte Verzeichnisse sind dadurch kein Laufzeitfehler
// (spec/spezifikation.md §`.d-check.yml` scan.ignore).
func dirIgnored(rel string, ignore []string) bool {
	for _, pat := range ignore {
		if MatchGlob(pat, rel) {
			return true
		}
		if sub, ok := strings.CutSuffix(pat, "/**"); ok && MatchGlob(sub, rel) {
			return true
		}
	}
	return false
}
