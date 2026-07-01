package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"fmt"
	"sort"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Result ist das Ergebnis eines Prüflaufs.
type Result struct {
	Findings     []model.Finding
	FilesChecked int
}

// Run führt den Prüflauf ohne das opt-in Modul vcs aus — der
// abwärtskompatible Einstieg für bestehende Aufrufer (vgl. RunWithVCS).
func Run(fsys driven.Filesystem, httpc driven.HTTPChecker, cfg model.Config, modules []string) (Result, error) {
	return RunWithVCS(fsys, httpc, nil, "", "", cfg, modules)
}

// RunWithVCS führt den Prüflauf aus (spec/spezifikation.md
// §DC-FA-CLI-001.a Schritte 3–5). Umgebungsfehler (nicht lesbare
// Wurzel/Datei) führen zu error — der Aufrufer mappt auf Exit 2.
// httpc wird ausschließlich vom Modul external benutzt; vcs/vcsBase/vcsHead
// ausschließlich vom opt-in Modul vcs (DC-FA-VCS-001) — ein nicht-hermetischer
// git-Lauf, der nur das lokale read-only `.git` liest (DC-QA-03: keine
// Netzwerkzugriffe; ohne aktives vcs bleibt vcs nil und ungenutzt).
func RunWithVCS(fsys driven.Filesystem, httpc driven.HTTPChecker, vcs driven.VCS, vcsBase, vcsHead string, cfg model.Config, modules []string) (Result, error) {
	var res Result
	active := map[string]bool{}
	for _, m := range modules {
		active[m] = true
	}

	// model.Config-Constraint (DC-FA-CONF-001, spec/spezifikation.md §2):
	// deklarierte ids-Targets müssen existieren — unabhängig von den
	// aktiven Modulen; Verletzung bedeutet Exit 2 ohne Prüfung.
	if err := ensureIDTargetsExist(fsys, cfg.IDPatterns); err != nil {
		return res, err
	}
	// diagrams.patterns[].defined-in müssen existieren (DC-FA-DIAG-001,
	// analog ids-Targets) — Verletzung bedeutet Exit 2 ohne Prüfung.
	if err := ensureDiagramsDefinedInExist(fsys, cfg.Diagrams); err != nil {
		return res, err
	}
	// Modul versions: aktuelle Version aus current-from auflösen
	// (DC-FA-VER-001.a; fail-closed → Exit 2). Nur wenn aktiv + konfiguriert.
	var versionsCurrent, versionsFromFile string
	if active["versions"] && cfg.Versions.PinPattern != nil {
		cur, from, verr := resolveCurrentVersion(fsys, cfg.Versions.EffectiveCurrentFrom())
		if verr != nil {
			return res, verr
		}
		versionsCurrent, versionsFromFile = cur, from
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
		slugCache:       map[string]map[string]bool{},
		statusCache:     map[string]*string{},
		diagCache:       map[string]map[string]bool{},
		spanCache:       map[string]string{},
		versionsCurrent: versionsCurrent, versionsFromFile: versionsFromFile,
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
	// Post-Pässe (vcs/commits über den VCS-Port, planning hermetisch über fsys):
	// arbeiten NACH dem Datei-Scan auf git-Historie bzw. Roadmap-Layout, nicht auf
	// den gescannten Dateien. Ein Port-Fehler (fehlendes .git/Range) ist fail-closed.
	post, perr := runPostPasses(fsys, vcs, vcsBase, vcsHead, cfg, active)
	if perr != nil {
		return res, perr
	}
	res.Findings = append(res.Findings, post...)
	res.Findings = model.SortFindings(res.Findings)
	return res, nil
}

// runPostPasses führt die nicht-datei-scannenden Module aus: vcs/commits (git-Diff
// bzw. Commit-Messages über den VCS-Port, fail-closed als error) und planning
// (hermetische Roadmap-↔-in-progress-Invariante über fsys, fail-closed als Befund).
// Ausgelagert, damit RunWithVCS unter der gocyclo-Schwelle bleibt.
func runPostPasses(fsys driven.Filesystem, vcs driven.VCS, vcsBase, vcsHead string, cfg model.Config, active map[string]bool) ([]model.Finding, error) {
	var out []model.Finding
	if active["vcs"] {
		vf, err := CheckVCS(vcs, cfg.VCS, vcsBase, vcsHead)
		if err != nil {
			return nil, err
		}
		out = append(out, vf...)
	}
	if active["commits"] {
		cf, err := CheckCommits(vcs, cfg.Commits, vcsBase, vcsHead)
		if err != nil {
			return nil, err
		}
		out = append(out, cf...)
	}
	if active["planning"] {
		out = append(out, CheckPlanning(fsys, cfg.Planning)...)
	}
	return out, nil
}

// runState bündelt den Lauf-Zustand (Caches, Befunde) — der
// Datei-Dispatch lebt in checkFile, damit Run flach bleibt.
type runState struct {
	fsys        driven.Filesystem
	cfg         model.Config
	active      map[string]bool
	inScope     map[string]map[string]bool
	slugCache   map[string]map[string]bool
	statusCache map[string]*string
	diagCache   map[string]map[string]bool
	spanCache   map[string]string
	extRefs     []ExternalRef
	findings    []model.Finding
	// versionsCurrent/versionsFromFile: am Lauf-Start aufgelöste aktuelle
	// Version + ihr Datei-Pfad (Modul versions, DC-FA-VER-001).
	versionsCurrent  string
	versionsFromFile string
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
	if st.applies("diagrams", file) {
		st.findings = append(st.findings, CheckDiagrams(st.fsys, file, content, st.cfg.Diagrams, st.diagCache)...)
	}
	if st.applies("versions", file) {
		st.findings = append(st.findings, CheckVersions(file, content, st.cfg.Versions, st.versionsCurrent, st.versionsFromFile)...)
	}
	if st.applies("pins", file) {
		st.findings = append(st.findings, CheckPins(st.fsys, file, lines, content, st.spanCache)...)
	}
	if st.applies("immutable", file) {
		st.findings = append(st.findings, CheckImmutable(file, lines, content, st.cfg.Immutable)...)
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
func discoverScopes(fsys driven.Filesystem, cfg model.Config, modules []string) (map[string]map[string]bool, []string, error) {
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
func ensureIDTargetsExist(fsys driven.Filesystem, patterns []model.IDPattern) error {
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

// ensureDiagramsDefinedInExist prüft die in diagrams.patterns[]
// deklarierten defined-in-Dateien (DC-FA-DIAG-001): sie müssen
// existieren und innerhalb der Repo-Wurzel liegen (analog
// ensureIDTargetsExist).
func ensureDiagramsDefinedInExist(fsys driven.Filesystem, cfg model.DiagramsConfig) error {
	for i, p := range cfg.Patterns {
		rel, escaped := ResolveConfigPath(p.DefinedIn)
		if escaped {
			return fmt.Errorf("diagrams.patterns[%d].defined-in verlässt die Repository-Wurzel: %s", i, p.DefinedIn)
		}
		if rel == "" {
			continue
		}
		kind, err := fsys.Kind(rel)
		if err != nil {
			return err
		}
		if kind == driven.KindMissing {
			return fmt.Errorf("diagrams.patterns[%d].defined-in existiert nicht: %s", i, p.DefinedIn)
		}
	}
	return nil
}
