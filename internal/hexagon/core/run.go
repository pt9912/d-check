package core

import (
	"fmt"
	"sort"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Result ist das Ergebnis eines Prüflaufs.
type Result struct {
	Findings     []Finding
	FilesChecked int
}

// Run führt den Prüflauf aus (spec/spezifikation.md
// §DC-FA-CLI-001.a Schritte 3–5). Umgebungsfehler (nicht lesbare
// Wurzel/Datei) führen zu error — der Aufrufer mappt auf Exit 2.
// httpc wird ausschließlich vom explizit aktivierten Modul external
// benutzt (DC-QA-03: keine Netzwerkzugriffe im Default).
func Run(fsys driven.Filesystem, httpc driven.HTTPChecker, cfg Config, modules []string) (Result, error) {
	var res Result
	active := map[string]bool{}
	for _, m := range modules {
		active[m] = true
	}

	// Config-Constraint (DC-FA-CONF-001, spec/spezifikation.md §2):
	// deklarierte ids-Targets müssen existieren — unabhängig von den
	// aktiven Modulen; Verletzung bedeutet Exit 2 ohne Prüfung.
	if err := ensureIDTargetsExist(fsys, cfg.IDPatterns); err != nil {
		return res, err
	}

	// Effektiver Scan-Scope pro Modul (DC-FA-CONF-002,
	// spec/spezifikation.md §DC-FA-CONF-002.a): Module mit
	// deklariertem scope laufen über einen eigenen Discover-Lauf,
	// alle übrigen über den globalen Scan; geprüft wird die
	// Vereinigungsmenge, jede Datei wird genau einmal gelesen.
	inScope, files, err := discoverScopes(fsys, cfg, modules)
	if err != nil {
		return res, err
	}
	res.FilesChecked = len(files)

	st := &runState{
		fsys: fsys, cfg: cfg, active: active, inScope: inScope,
		slugCache:   map[string]map[string]bool{},
		statusCache: map[string]*string{},
	}
	for _, file := range files {
		if err := st.checkFile(file); err != nil {
			return res, err
		}
	}
	res.Findings = st.findings
	if active["external"] {
		res.Findings = append(res.Findings,
			CheckExternal(httpc, st.extRefs, cfg.External.EffectiveParallel())...)
	}
	res.Findings = SortFindings(res.Findings)
	return res, nil
}

// runState bündelt den Lauf-Zustand (Caches, Befunde) — der
// Datei-Dispatch lebt in checkFile, damit Run flach bleibt.
type runState struct {
	fsys        driven.Filesystem
	cfg         Config
	active      map[string]bool
	inScope     map[string]map[string]bool
	slugCache   map[string]map[string]bool
	statusCache map[string]*string
	extRefs     []ExternalRef
	findings    []Finding
}

// applies: Modul aktiv und Datei in dessen effektivem Scope
// (DC-FA-CONF-002.a).
func (st *runState) applies(module, file string) bool {
	return st.active[module] && st.inScope[module][file]
}

// checkFile liest die Datei genau einmal und lässt alle anwendbaren
// Module prüfen (spec/spezifikation.md §DC-FA-CONF-002.a Schritt 2).
func (st *runState) checkFile(file string) error {
	content, err := st.fsys.ReadFile(file)
	if err != nil {
		return fmt.Errorf("%s nicht lesbar: %w", file, err)
	}
	lines := PreprocessMarkdown(content)
	if st.applies("links", file) {
		st.findings = append(st.findings, CheckLinks(st.fsys, file, lines)...)
	}
	if st.applies("anchors", file) {
		st.findings = append(st.findings, CheckAnchors(st.fsys, file, content, lines, st.slugCache)...)
	}
	if st.applies("ids", file) {
		st.findings = append(st.findings, CheckIDs(file, content, lines, st.cfg.IDPatterns)...)
	}
	if st.applies("matrix", file) {
		st.findings = append(st.findings, CheckMatrix(st.fsys, file, content, lines, st.cfg.Matrix, st.statusCache)...)
	}
	if st.applies("codepaths", file) {
		st.findings = append(st.findings, CheckCodepaths(st.fsys, file, content, st.cfg.Codepaths, st.slugCache)...)
	}
	if st.applies("hostpaths", file) {
		st.findings = append(st.findings, CheckHostpaths(file, content, st.cfg.Hostpaths)...)
	}
	if st.applies("spans", file) {
		st.findings = append(st.findings, CheckSpans(file, content, lines)...)
	}
	if st.applies("external", file) {
		st.extRefs = append(st.extRefs, CollectExternalURLs(file, lines)...)
	}
	return nil
}

// discoverScopes ermittelt pro aktivem Modul die Datei-Menge seines
// effektiven Scan-Scopes und die bytewise sortierte Vereinigungsmenge
// (DC-FA-CONF-002.a). Der globale Scan läuft höchstens einmal; jeder
// Modul-Scope ist ein eigener Discover-Lauf mit denselben Regeln
// (Existenz, Repo-Escape, Pruning, SKIP_DIRS — DC-FA-SCAN-001).
func discoverScopes(fsys driven.Filesystem, cfg Config, modules []string) (map[string]map[string]bool, []string, error) {
	inScope := map[string]map[string]bool{}
	union := map[string]bool{}

	var global map[string]bool
	// Leerer Modulsatz (modules: [] — bewusste Setzung): der globale
	// Scan läuft trotzdem, damit die Datei-Zählung byte-identisch zum
	// Verhalten vor DC-FA-CONF-002 bleibt (Boundary-AK).
	if len(modules) == 0 {
		if _, err := discoverInto(fsys, cfg.Roots, cfg.Ignore, union); err != nil {
			return nil, nil, err
		}
	}
	for _, m := range modules {
		sc := cfg.Scopes[m]
		if sc == nil {
			if global == nil {
				g, err := discoverInto(fsys, cfg.Roots, cfg.Ignore, union)
				if err != nil {
					return nil, nil, err
				}
				global = g
			}
			inScope[m] = global
			continue
		}
		set, err := discoverInto(fsys, sc.Roots, sc.Ignore, union)
		if err != nil {
			return nil, nil, fmt.Errorf("%s.scope: %w", m, err)
		}
		inScope[m] = set
	}

	files := make([]string, 0, len(union))
	for f := range union {
		files = append(files, f)
	}
	sort.Strings(files)
	return inScope, files, nil
}

// discoverInto führt einen Discover-Lauf aus und sammelt die Treffer
// als Menge sowie additiv in die Vereinigungsmenge.
func discoverInto(fsys driven.Filesystem, roots, ignore []string, union map[string]bool) (map[string]bool, error) {
	files, err := DiscoverFiles(fsys, roots, ignore)
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(files))
	for _, f := range files {
		set[f] = true
		union[f] = true
	}
	return set, nil
}

// ensureIDTargetsExist prüft die in ids.patterns[] deklarierten
// Targets (Datei oder Verzeichnis, relativ zur Repo-Wurzel): sie
// müssen existieren und innerhalb der Repo-Wurzel liegen.
func ensureIDTargetsExist(fsys driven.Filesystem, patterns []IDPattern) error {
	for _, p := range patterns {
		rel, escaped := ResolveConfigPath(p.Target)
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
