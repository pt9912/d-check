package core

import (
	"fmt"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Result ist das Ergebnis eines Prüflaufs.
type Result struct {
	Findings     []Finding
	FilesChecked int
	// SkippedModules: aktivierte, aber (noch) nicht implementierte
	// Module — der Aufrufer meldet sie auf stderr.
	SkippedModules []string
}

// Run führt den Prüflauf aus (spec/spezifikation.md
// §DC-FA-CLI-001.a Schritte 3–5). Umgebungsfehler (nicht lesbare
// Wurzel/Datei) führen zu error — der Aufrufer mappt auf Exit 2.
func Run(fsys driven.Filesystem, cfg Config, modules []string) (Result, error) {
	var res Result
	active := map[string]bool{}
	for _, m := range modules {
		if isImplemented(m) {
			active[m] = true
		} else {
			res.SkippedModules = append(res.SkippedModules, m)
		}
	}

	// Config-Constraint (DC-FA-CONF-001, spec/spezifikation.md §2):
	// deklarierte ids-Targets müssen existieren — unabhängig von den
	// aktiven Modulen; Verletzung bedeutet Exit 2 ohne Prüfung.
	if err := ensureIDTargetsExist(fsys, cfg.IDPatterns); err != nil {
		return res, err
	}

	files, err := DiscoverFiles(fsys, cfg.Roots, cfg.Ignore)
	if err != nil {
		return res, err
	}
	res.FilesChecked = len(files)

	slugCache := map[string]map[string]bool{}
	for _, file := range files {
		content, err := fsys.ReadFile(file)
		if err != nil {
			return res, fmt.Errorf("%s nicht lesbar: %w", file, err)
		}
		lines := PreprocessMarkdown(content)
		if active["links"] {
			res.Findings = append(res.Findings, checkLinks(fsys, file, lines)...)
		}
		if active["anchors"] {
			res.Findings = append(res.Findings, checkAnchors(fsys, file, content, lines, slugCache)...)
		}
		if active["ids"] {
			res.Findings = append(res.Findings, checkIDs(file, lines, cfg.IDPatterns)...)
		}
	}
	res.Findings = SortFindings(res.Findings)
	return res, nil
}

// ensureIDTargetsExist prüft die in ids.patterns[] deklarierten
// Targets (Datei oder Verzeichnis, relativ zur Repo-Wurzel): sie
// müssen existieren und innerhalb der Repo-Wurzel liegen.
func ensureIDTargetsExist(fsys driven.Filesystem, patterns []IDPattern) error {
	for _, p := range patterns {
		rel, escaped := resolveConfigPath(p.Target)
		if escaped {
			return fmt.Errorf("konfiguriertes ids-Target verlässt die Repository-Wurzel: %s", p.Target)
		}
		if rel == "" {
			continue // Repo-Wurzel existiert per Definition
		}
		kind, err := fsys.Kind(rel)
		if err != nil {
			return err
		}
		if kind == driven.KindMissing {
			return fmt.Errorf("konfiguriertes ids-Target existiert nicht: %s", p.Target)
		}
	}
	return nil
}
