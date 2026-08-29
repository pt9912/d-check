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
	return RunWithVCS(fsys, httpc, nil, nil, "", "", cfg, modules)
}

// RunWithVCS führt den Prüflauf aus (spec/spezifikation.md
// §DC-FA-CLI-001.a Schritte 3–5). Umgebungsfehler (nicht lesbare
// Wurzel/Datei) führen zu error — der Aufrufer mappt auf Exit 2.
// httpc wird ausschließlich vom Modul external benutzt; vcs ausschließlich
// von den opt-in git-Modulen (vcs/commits: mit vcsBase/vcsHead-Range,
// DC-FA-VCS-001/DC-FA-COMMITS-001; tracked: nur der Index, ohne Range,
// DC-FA-TRK-001) — nicht-hermetische Läufe, die nur das lokale read-only
// `.git` lesen (DC-QA-03: keine Netzwerkzugriffe; ohne aktives git-Modul
// bleibt vcs nil und ungenutzt).
func RunWithVCS(fsys driven.Filesystem, httpc driven.HTTPChecker, vcs driven.VCS, wp driven.WorkflowParser, vcsBase, vcsHead string, cfg model.Config, modules []string) (Result, error) {
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
	// Modul versions: je Paar die aktuelle Version aus seiner current-from-
	// Quelle auflösen (DC-FA-VER-001.a; fail-closed → Exit 2). Nur wenn aktiv;
	// ohne Paar bleibt die Liste leer und das Modul wirkungslos.
	var versionsPatterns []ResolvedVersionPattern
	if active["versions"] {
		for i, p := range cfg.Versions.Patterns {
			cur, from, verr := resolveCurrentVersion(fsys, p.EffectiveCurrentFrom(), versionSource(cfg.Versions, i))
			if verr != nil {
				return res, verr
			}
			versionsPatterns = append(versionsPatterns, ResolvedVersionPattern{Pattern: p, Current: cur, FromFile: from, Index: i})
		}
	}

	// Modul tracked (DC-FA-TRK-001.a Schritte 1–2): die Index-Menge wird
	// einmal je Lauf über den VCS-Port geladen — fail-closed: aktives
	// tracked ohne verdrahteten Port (kein lesbares .git) oder mit
	// unlesbarem Index ⇒ error, der Aufrufer mappt auf Exit 2.
	var trackedSet map[string]bool
	if active["tracked"] {
		if vcs == nil {
			return res, fmt.Errorf("das Modul tracked braucht ein lesbares .git unter der Scan-Wurzel (DC-FA-TRK-001, fail-closed)")
		}
		set, terr := vcs.TrackedPaths()
		if terr != nil {
			return res, terr
		}
		trackedSet = set
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
		slugCache:        map[string]map[string]bool{},
		statusCache:      map[string]*string{},
		diagCache:        map[string]map[string]bool{},
		spanCache:        map[string]string{},
		versionsPatterns: versionsPatterns,
		tracked:          trackedSet,
	}
	for _, file := range files {
		if err := st.checkFile(file); err != nil {
			return res, err
		}
	}
	res.Findings = st.findings
	res.Findings = append(res.Findings, st.netFindings(httpc)...)
	// Post-Pässe (vcs/commits über den VCS-Port, planning hermetisch über fsys):
	// arbeiten NACH dem Datei-Scan auf git-Historie bzw. Roadmap-Layout, nicht auf
	// den gescannten Dateien. Ein Port-Fehler (fehlendes .git/Range) ist fail-closed.
	post, perr := runPostPasses(fsys, vcs, wp, vcsBase, vcsHead, cfg, active)
	if perr != nil {
		return res, perr
	}
	res.Findings = append(res.Findings, post...)
	res.Findings = model.SortFindings(res.Findings)
	return res, nil
}

// netFindings führt die Netz-Post-Pässe aus (geteilter httpc, DC-QA-03; beide
// strikt opt-in): external prüft Erreichbarkeit über die gesammelten Refs,
// sources holt/hasht/vergleicht die Marker- + Config-Pins (DC-FA-SRC-001.a).
func (st *runState) netFindings(httpc driven.HTTPChecker) []model.Finding {
	var out []model.Finding
	if st.active["external"] {
		out = append(out, CheckExternal(httpc, st.extRefs, st.cfg.External.EffectiveParallel())...)
	}
	if st.active["sources"] {
		out = append(out, CheckSources(httpc, st.srcRefs, st.cfg.Sources, st.cfg.ConfigFile)...)
	}
	return out
}

// runPostPasses führt die nicht-datei-scannenden Module aus: vcs/commits (git-Diff
// bzw. Commit-Messages über den VCS-Port, fail-closed als error) und planning
// (hermetische Roadmap-↔-in-progress-Invariante über fsys, fail-closed als Befund).
// Ausgelagert, damit RunWithVCS unter der gocyclo-Schwelle bleibt.
func runPostPasses(fsys driven.Filesystem, vcs driven.VCS, wp driven.WorkflowParser, vcsBase, vcsHead string, cfg model.Config, active map[string]bool) ([]model.Finding, error) {
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
	if active["links"] {
		out = append(out, CheckResolveFromDirs(fsys, cfg.ResolveFrom)...)
	}
	if active["planning"] {
		out = append(out, CheckPlanning(fsys, cfg.Planning)...)
		out = append(out, CheckPlanningClosure(fsys, cfg.Planning)...)
		out = append(out, CheckPlanningWaves(fsys, cfg.Planning)...)
	}
	if active["structure"] {
		out = append(out, CheckStructure(fsys, cfg.Structure)...)
	}
	if active["workflows"] {
		out = append(out, CheckWorkflows(fsys, wp, cfg.Workflows)...)
	}
	if active["targets"] {
		tf, err := CheckTargets(fsys, cfg.Targets)
		if err != nil {
			return nil, err
		}
		out = append(out, tf...)
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
	srcRefs     []SourceRef
	findings    []model.Finding
	// versionsPatterns: die Muster-Quellen-Paare samt der am Lauf-Start
	// aufgelösten erwarteten Version und ihrem Datei-Pfad (Modul versions,
	// DC-FA-VER-001).
	versionsPatterns []ResolvedVersionPattern
	// tracked: die einmal je Lauf geladene git-Index-Menge (Modul tracked,
	// DC-FA-TRK-001.a Schritt 2); nil, wenn das Modul inaktiv ist.
	tracked map[string]bool
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
		st.findings = append(st.findings, CheckLinks(st.fsys, file, lines, st.cfg.IgnoreRefs)...)
		st.findings = append(st.findings, CheckResolveFrom(st.fsys, file, lines, st.cfg.ResolveFrom, st.cfg.IgnoreRefs)...)
	}
	if st.applies("anchors", file) {
		st.findings = append(st.findings, CheckAnchors(st.fsys, file, content, lines, st.slugCache, st.cfg.IgnoreRefs)...)
	}
	if st.applies("ids", file) {
		st.findings = append(st.findings, CheckIDs(file, content, lines, st.cfg.IDPatterns)...)
	}
	if st.applies("matrix", file) {
		st.findings = append(st.findings, CheckMatrix(st.fsys, file, content, lines, st.cfg.Matrix, st.statusCache)...)
	}
	if st.applies("codepaths", file) {
		st.findings = append(st.findings, CheckCodepaths(st.fsys, file, content, st.cfg.Codepaths, st.slugCache, st.cfg.IgnoreRefs)...)
	}
	if err := st.appendCitationFindings(file, content); err != nil {
		return err
	}
	if err := st.appendSourceRefs(file, content, lines); err != nil {
		return err
	}
	st.appendTailFindings(file, content, lines)
	return nil
}

// appendSourceRefs sammelt die Marker-Pins des Moduls sources; eine malformte
// source-pin-Direktive ist fail-closed → error (Aufrufer mappt auf Exit 2,
// DC-FA-SRC-001.a Schritt 2). Config-Pins verarbeitet der Post-Pass separat.
func (st *runState) appendSourceRefs(file string, content []byte, lines []Line) error {
	if !st.applies("sources", file) {
		return nil
	}
	refs, err := CollectSourcePins(file, lines, content)
	if err != nil {
		return err
	}
	st.srcRefs = append(st.srcRefs, refs...)
	return nil
}

// appendCitationFindings führt das Modul citations aus; eine strukturell
// unbrauchbare d-check:cite-Direktive ist fail-closed → error (Aufrufer
// mappt auf Exit 2, DC-FA-CITE-001.a).
func (st *runState) appendCitationFindings(file string, content []byte) error {
	if !st.applies("citations", file) {
		return nil
	}
	cf, err := CheckCitations(st.fsys, file, content)
	if err != nil {
		return err
	}
	st.findings = append(st.findings, cf...)
	return nil
}

// appendTailFindings führt die restlichen inhalts-basierten Module aus
// (keines meldet fail-closed); ausgelagert, damit checkFile flach bleibt.
func (st *runState) appendTailFindings(file string, content []byte, lines []Line) {
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
		st.findings = append(st.findings, CheckVersions(file, content, st.versionsPatterns)...)
	}
	if st.applies("pins", file) {
		st.findings = append(st.findings, CheckPins(st.fsys, file, lines, content, st.spanCache)...)
	}
	if st.applies("immutable", file) {
		st.findings = append(st.findings, CheckImmutable(file, lines, content, st.cfg.Immutable)...)
	}
	if st.applies("tracked", file) {
		st.findings = append(st.findings, CheckTracked(st.fsys, file, lines, st.cfg.Tracked, st.tracked)...)
	}
	if st.applies("external", file) {
		st.extRefs = append(st.extRefs, CollectExternalURLs(file, lines)...)
	}
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
		// Anders als ein ids-Target darf defined-in KEIN Verzeichnis sein: das
		// Modul liest die Datei, um die Definitionsmenge zu bilden. Ein
		// Verzeichnis liefert eine leere Menge — jede Kennung des Diagramms
		// waere dann "undefiniert", ein Befund-Sturm ohne Hinweis auf die
		// Ursache. Fail-closed statt still falsch.
		if kind != driven.KindFile {
			return fmt.Errorf("diagrams.patterns[%d].defined-in ist keine Datei: %s", i, p.DefinedIn)
		}
	}
	return nil
}
