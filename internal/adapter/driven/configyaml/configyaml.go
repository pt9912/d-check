// Package configyaml ist der driven Config-Adapter: striktes Decoding
// und statische Validierung von .d-check.yml (DC-FA-CONF-001;
// Strategie gemäß ADR-0003). Den Datei-Inhalt beschafft das CLI über
// den Filesystem-Adapter — dieser Adapter macht kein I/O
// (spec/architecture.md §2). Die Existenz-Prüfung deklarierter
// Scan-Wurzeln erfolgt im Kern beim Lauf-Start (ebenfalls Exit 2).
package configyaml

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

)

// FileName ist der feste Name der Konfigurationsdatei.
const FileName = ".d-check.yml"

type rawIDPattern struct {
	Regex       string   `yaml:"regex"`
	Target      string   `yaml:"target"`
	LinkPolicy  string   `yaml:"link-policy"`
	ExemptPaths []string `yaml:"exempt-paths"`
}

type rawIDs struct {
	Scope    *rawScope      `yaml:"scope"`
	Patterns []rawIDPattern `yaml:"patterns"`
}

// rawScope ist der modul-lokale Scan-Scope (DC-FA-CONF-002). Roots
// ist Pflicht, wenn scope gesetzt ist — nil (fehlend) ist von der
// expliziten leeren Liste unterscheidbar.
type rawScope struct {
	Roots  []string `yaml:"roots"`
	Ignore []string `yaml:"ignore"`
}

// rawScopeOnly traegt Module, die ausser scope keine eigenen
// Konfigurations-Schluessel haben (links, anchors).
type rawScopeOnly struct {
	Scope *rawScope `yaml:"scope"`
}

// rawCodepaths traegt scope, roots und das exempt-paths-Ventil
// (DC-FA-CODE-001; benannt, damit der Typ in raw und scopeOfCodepaths
// identisch ist).
type rawCodepaths struct {
	Scope       *rawScope `yaml:"scope"`
	Roots       []string  `yaml:"roots"`
	ExemptPaths []string  `yaml:"exempt-paths"`
	IgnoreRefs  []string  `yaml:"ignore-refs"`
}

// rawDiagrams trägt scope, fences und die Muster des Moduls diagrams
// (DC-FA-DIAG-001).
type rawDiagrams struct {
	Scope    *rawScope           `yaml:"scope"`
	Fences   []string            `yaml:"fences"`
	Patterns []rawDiagramPattern `yaml:"patterns"`
}

type rawDiagramPattern struct {
	Regex     string `yaml:"regex"`
	DefinedIn string `yaml:"defined-in"`
}

// rawVersions trägt scope, das Pin-Muster, die Quelle der aktuellen Version
// und das exempt-paths-Ventil des Moduls versions (DC-FA-VER-001).
type rawVersions struct {
	Scope       *rawScope `yaml:"scope"`
	PinPattern  string    `yaml:"pin-pattern"`
	CurrentFrom string    `yaml:"current-from"`
	ExemptPaths []string  `yaml:"exempt-paths"`
}

// rawImmutable trägt scope und exclude-sections des Moduls immutable
// (DC-FA-IMM-001); es gibt keine Pflichtfelder (opt-in pro Datei über den
// Pin-Marker).
type rawImmutable struct {
	Scope           *rawScope `yaml:"scope"`
	ExcludeSections []string  `yaml:"exclude-sections"`
}

// rawVCS trägt scope und die Parameter des Moduls vcs (DC-FA-VCS-001):
// paths (Glob-Klasse der geschützten Dateien), immutable-when (Pflicht-Regex
// der BASE-Immutabilität), exclude-sections, status-line (Kopf-Status-Zeile)
// und head-allow (erlaubter Status-Übergang).
type rawVCS struct {
	Scope           *rawScope `yaml:"scope"`
	Paths           []string  `yaml:"paths"`
	ImmutableWhen   string    `yaml:"immutable-when"`
	ExcludeSections []string  `yaml:"exclude-sections"`
	StatusLine      string    `yaml:"status-line"`
	HeadAllow       string    `yaml:"head-allow"`
}

// rawPlanning trägt die Parameter des Moduls planning (DC-FA-PLAN-001): roadmap
// (Pflicht zum Aktivieren, leer ⇒ inert) plus die Konventions-Defaults heading/
// marker/slice-glob. **Keine** scope — planning ist ein Post-Pass ohne Datei-Scan
// (ein `planning.scope` wäre wirkungslos; der strikte Decoder lehnt es ab).
type rawPlanning struct {
	Roadmap   string `yaml:"roadmap"`
	Heading   string `yaml:"heading"`
	Marker    string `yaml:"marker"`
	SliceGlob string `yaml:"slice-glob"`
}

// rawCommits trägt scope und die Parameter des Moduls commits
// (DC-FA-COMMITS-001): id-patterns (Regex-Liste gültiger Traceability-Kennungen,
// leer ⇒ inert) und exempt-pattern (Betreff-Ausnahme, z. B. `^(Merge |Revert )`).
type rawCommits struct {
	Scope         *rawScope `yaml:"scope"`
	IDPatterns    []string  `yaml:"id-patterns"`
	ExemptPattern string    `yaml:"exempt-pattern"`
}

type rawMatrix struct {
	Scope   *rawScope `yaml:"scope"`
	Classes []struct {
		Name      string   `yaml:"name"`
		Paths     []string `yaml:"paths"`
		Order     []string `yaml:"order"`
		Direction string   `yaml:"direction"`
		Token     string   `yaml:"token"`
	} `yaml:"classes"`
	Rules []struct {
		From  string `yaml:"from"`
		To    string `yaml:"to"`
		Allow bool   `yaml:"allow"`
	} `yaml:"rules"`
	Status *struct {
		Forbidden             []string `yaml:"forbidden"`
		AllowSupersedeLineage bool     `yaml:"allow-supersede-lineage"`
		SupersedeFields       []string `yaml:"supersede-fields"`
	} `yaml:"status"`
	ExcludeSections []string `yaml:"exclude-sections"`
	ExemptPaths     []string `yaml:"exempt-paths"`
}

// rawExternal nutzt Pointer, damit ein explizit gesetzter Wert 0 vom
// Nicht-Setzen unterscheidbar bleibt (Constraint 1–300 bzw. 1–16 —
// auch 0 ist ein Konfigurationsfehler, kein stiller Default).
type rawExternal struct {
	Scope          *rawScope `yaml:"scope"`
	TimeoutSeconds *int      `yaml:"timeout-seconds"`
	Parallel       *int      `yaml:"parallel"`
}

// raw bildet das Voll-Schema von .d-check.yml ab
// (spec/spezifikation.md §2). Unbekannte Schlüssel sind durch
// KnownFields(true) Fehler.
type raw struct {
	Scan *struct {
		Roots  []string `yaml:"roots"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"scan"`
	Modules  []string      `yaml:"modules"`
	Links    *rawScopeOnly `yaml:"links"`
	Anchors  *rawScopeOnly `yaml:"anchors"`
	Spans    *rawScopeOnly `yaml:"spans"`
	Pins     *rawScopeOnly `yaml:"pins"`
	Hostpaths *struct {
		Scope    *rawScope `yaml:"scope"`
		Prefixes []string  `yaml:"prefixes"`
	} `yaml:"hostpaths"`
	IDs      *rawIDs       `yaml:"ids"`
	Matrix   *rawMatrix    `yaml:"matrix"`
	External *rawExternal  `yaml:"external"`
	Codepaths *rawCodepaths `yaml:"codepaths"`
	Diagrams  *rawDiagrams  `yaml:"diagrams"`
	Versions  *rawVersions  `yaml:"versions"`
	Immutable *rawImmutable `yaml:"immutable"`
	Vcs       *rawVCS       `yaml:"vcs"`
	Commits   *rawCommits   `yaml:"commits"`
	Planning  *rawPlanning  `yaml:"planning"`
}

// Decode parst und validiert den Datei-Inhalt vollständig — Syntax
// und statische Semantik; jeder Fehler bedeutet Exit 2 ohne Prüfung.
// content == nil bedeutet: keine Konfigurationsdatei → Defaults.
func Decode(content []byte) (model.Config, error) {
	var cfg model.Config
	if content == nil {
		return cfg, nil
	}
	r, err := decodeStrict(content)
	if err != nil {
		return cfg, err
	}
	if r == nil {
		return cfg, nil // leeres YAML-Null-Dokument → Defaults
	}

	if r.Scan != nil {
		cfg.Roots = r.Scan.Roots
		cfg.Ignore = r.Scan.Ignore
	}
	cfg.Modules = r.Modules

	if err := applyModules(r, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// applyModules wendet die modul-spezifischen Validierungen/Kompilierungen in
// fester Reihenfolge an und liefert den ersten Fehler (Exit 2). Ausgelagert aus
// Decode, damit dessen gocyclo-Komplexität unter der Schwelle bleibt.
func applyModules(r *raw, cfg *model.Config) error {
	if err := applyIDs(r.IDs, cfg); err != nil {
		return err
	}
	if err := applyMatrix(r.Matrix, cfg); err != nil {
		return err
	}
	if err := applyExternal(r.External, cfg); err != nil {
		return err
	}
	if err := applyCodepaths(r, cfg); err != nil {
		return err
	}
	if err := applyHostpaths(r, cfg); err != nil {
		return err
	}
	if err := applyDiagrams(r, cfg); err != nil {
		return err
	}
	if err := applyVersions(r, cfg); err != nil {
		return err
	}
	applyImmutable(r, cfg)
	if err := applyVCS(r, cfg); err != nil {
		return err
	}
	if err := applyCommits(r, cfg); err != nil {
		return err
	}
	if err := applyPlanning(r, cfg); err != nil {
		return err
	}
	return applyScopes(r, cfg)
}

// compileConfigRegex kompiliert ein Zeilen-Regex einer Modul-Konfiguration
// (vcs, commits). required ⇒ ein leeres Muster ist ein Konfigurationsfehler;
// sonst ⇒ leer ergibt nil (Prüfung aus). Ein nicht kompilierbares Muster ist
// immer ein Fehler (Exit 2).
func compileConfigRegex(field, pattern string, required bool) (*regexp.Regexp, error) {
	if strings.TrimSpace(pattern) == "" {
		if required {
			return nil, fmt.Errorf("%s: %s fehlt", FileName, field)
		}
		return nil, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %v", FileName, field, err)
	}
	return re, nil
}

// applyVCS validiert und kompiliert die Parameter des Moduls vcs
// (DC-FA-VCS-001): immutable-when ist Pflicht, status-line/head-allow sind
// optional; alle müssen kompilieren. Ohne paths ist das Modul inert
// (byte-identisch zum Lauf ohne das Modul, DC-QA-02) — paths ist nicht Pflicht.
func applyVCS(r *raw, cfg *model.Config) error {
	if r.Vcs == nil {
		return nil
	}
	v := r.Vcs
	when, err := compileConfigRegex("vcs.immutable-when", v.ImmutableWhen, true)
	if err != nil {
		return err
	}
	statusLine, err := compileConfigRegex("vcs.status-line", v.StatusLine, false)
	if err != nil {
		return err
	}
	headAllow, err := compileConfigRegex("vcs.head-allow", v.HeadAllow, false)
	if err != nil {
		return err
	}
	cfg.VCS = model.VCSConfig{
		Paths: v.Paths, ImmutableWhen: when, ExcludeSections: v.ExcludeSections,
		StatusLine: statusLine, HeadAllow: headAllow,
	}
	return nil
}

// applyCommits validiert und kompiliert die Parameter des Moduls commits
// (DC-FA-COMMITS-001): jedes id-patterns-Muster muss kompilieren und darf den
// Leerstring nicht matchen (sonst gälte jede Message als getraced — Silent-Grün,
// Exit 2); exempt-pattern ist optional. Ohne id-patterns ist das Modul inert
// (byte-identisch zum Lauf ohne das Modul, DC-QA-02) — id-patterns ist nicht Pflicht.
func applyCommits(r *raw, cfg *model.Config) error {
	if r.Commits == nil {
		return nil
	}
	c := r.Commits
	var pats []*regexp.Regexp
	for i, p := range c.IDPatterns {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("%s: commits.id-patterns[%d] ist leer", FileName, i)
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("%s: commits.id-patterns[%d]: %v", FileName, i, err)
		}
		if re.MatchString("") {
			return fmt.Errorf("%s: commits.id-patterns[%d] matcht den Leerstring", FileName, i)
		}
		pats = append(pats, re)
	}
	exempt, err := compileConfigRegex("commits.exempt-pattern", c.ExemptPattern, false)
	if err != nil {
		return err
	}
	cfg.Commits = model.CommitsConfig{IDPatterns: pats, ExemptPattern: exempt}
	return nil
}

// applyPlanning validiert die Parameter des Moduls planning (DC-FA-PLAN-001):
// roadmap muss (wenn gesetzt) relativ zur Repo-Wurzel liegen (kein `/`, kein `..`);
// die Existenz prüft der Kern beim Lauf (fail-closed → planning-drift). Leere
// roadmap ⇒ Modul inert. heading/marker/slice-glob sind optional (Konventions-
// Defaults im Kern).
func applyPlanning(r *raw, cfg *model.Config) error {
	if r.Planning == nil {
		return nil
	}
	p := r.Planning
	if strings.HasPrefix(p.Roadmap, "/") || strings.Contains(p.Roadmap, "..") {
		return fmt.Errorf("%s: planning.roadmap %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, p.Roadmap)
	}
	// slice-glob muss ein gültiges path.Match-Muster sein — sonst schluckte der Kern
	// den ErrBadPattern und meldete fälschlich „keine Slices" (fail-open Silent-Grün am
	// Config-Rand; R2-MEDIUM slice-057). Nur ein explizit gesetztes Muster prüfen; den
	// Default `slice-*.md` kennt der Kern.
	if p.SliceGlob != "" {
		if _, err := path.Match(p.SliceGlob, "probe"); err != nil {
			return fmt.Errorf("%s: planning.slice-glob %q ist kein gültiges Glob: %v", FileName, p.SliceGlob, err)
		}
	}
	cfg.Planning = model.PlanningConfig{
		Roadmap: p.Roadmap, Heading: p.Heading, Marker: p.Marker, SliceGlob: p.SliceGlob,
	}
	return nil
}

// applyCodepaths validiert die codepaths-Präfixe (DC-FA-CODE-001).
func applyCodepaths(r *raw, cfg *model.Config) error {
	if r.Codepaths == nil {
		return nil
	}
	for _, root := range r.Codepaths.Roots {
		if strings.TrimSpace(root) == "" {
			return fmt.Errorf("%s: codepaths.roots enthält ein leeres Präfix", FileName)
		}
		if strings.HasPrefix(root, "/") || root == ".." || strings.Contains(root, "..") {
			return fmt.Errorf("%s: codepaths.roots-Präfix %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, root)
		}
	}
	cfg.Codepaths = model.CodepathsConfig{
		Roots:       r.Codepaths.Roots,
		ExemptPaths: r.Codepaths.ExemptPaths,
		IgnoreRefs:  r.Codepaths.IgnoreRefs,
	}
	return nil
}

// applyHostpaths validiert die hostpaths-Präfixliste (DC-FA-HOST-001:
// nicht-leere Verzeichnisnamen ohne '/').
func applyHostpaths(r *raw, cfg *model.Config) error {
	if r.Hostpaths == nil {
		return nil
	}
	for _, p := range r.Hostpaths.Prefixes {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("%s: hostpaths.prefixes enthält einen leeren Namen", FileName)
		}
		if strings.Contains(p, "/") {
			return fmt.Errorf("%s: hostpaths.prefixes-Eintrag %q muss ein Verzeichnisname ohne '/' sein", FileName, p)
		}
	}
	cfg.Hostpaths = model.HostpathsConfig{Prefixes: r.Hostpaths.Prefixes}
	return nil
}

// applyDiagrams validiert und kompiliert die diagrams-Muster
// (DC-FA-DIAG-001): nicht-leere fences-Einträge, kompilierbarer Regex
// (nicht den Leerstring matchend) und Pflicht-defined-in.
func applyDiagrams(r *raw, cfg *model.Config) error {
	if r.Diagrams == nil {
		return nil
	}
	for _, f := range r.Diagrams.Fences {
		if strings.TrimSpace(f) == "" {
			return fmt.Errorf("%s: diagrams.fences enthält einen leeren Eintrag", FileName)
		}
	}
	for i, p := range r.Diagrams.Patterns {
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			return fmt.Errorf("%s: diagrams.patterns[%d].regex: %v", FileName, i, err)
		}
		if re.MatchString("") {
			return fmt.Errorf("%s: diagrams.patterns[%d].regex matcht den Leerstring", FileName, i)
		}
		if p.DefinedIn == "" {
			return fmt.Errorf("%s: diagrams.patterns[%d].defined-in fehlt", FileName, i)
		}
		cfg.Diagrams.Patterns = append(cfg.Diagrams.Patterns, model.DiagramPattern{Regex: re, DefinedIn: p.DefinedIn})
	}
	cfg.Diagrams.Fences = r.Diagrams.Fences
	return nil
}

// applyVersions validiert und kompiliert das versions-Pin-Muster
// (DC-FA-VER-001): Pflicht-pin-pattern (kompilierbar, nicht den Leerstring
// matchend); current-from/exempt-paths werden im Kern beim Lauf-Start bzw. je
// Datei ausgewertet.
func applyVersions(r *raw, cfg *model.Config) error {
	if r.Versions == nil {
		return nil
	}
	v := r.Versions
	if strings.TrimSpace(v.PinPattern) == "" {
		return fmt.Errorf("%s: versions.pin-pattern fehlt", FileName)
	}
	re, err := regexp.Compile(v.PinPattern)
	if err != nil {
		return fmt.Errorf("%s: versions.pin-pattern: %v", FileName, err)
	}
	if re.MatchString("") {
		return fmt.Errorf("%s: versions.pin-pattern matcht den Leerstring", FileName)
	}
	cfg.Versions = model.VersionsConfig{PinPattern: re, CurrentFrom: v.CurrentFrom, ExemptPaths: v.ExemptPaths}
	return nil
}

// applyImmutable übernimmt exclude-sections des Moduls immutable
// (DC-FA-IMM-001). Keine Validierung/keine Pflichtfelder — ohne Pin-Marker in
// einer Datei ist das Modul wirkungslos (opt-in pro Datei); die Marker-Erkennung
// läuft im Kern.
func applyImmutable(r *raw, cfg *model.Config) {
	if r.Immutable == nil {
		return
	}
	cfg.Immutable = model.ImmutableConfig{ExcludeSections: r.Immutable.ExcludeSections}
}

// applyScopes übernimmt die modul-lokalen Scan-Scopes
// (DC-FA-CONF-002): scope ersetzt den globalen Scan für genau dieses
// Modul; roots ist Pflicht (keine stille Vererbung).
func applyScopes(r *raw, cfg *model.Config) error {
	scopes := []struct {
		module string
		scope  *rawScope
	}{
		{"links", scopeOf(r.Links)},
		{"anchors", scopeOf(r.Anchors)},
		{"spans", scopeOf(r.Spans)},
		{"hostpaths", scopeOfHostpaths(r.Hostpaths)},
		{"ids", scopeOfIDs(r.IDs)},
		{"matrix", scopeOfMatrix(r.Matrix)},
		{"external", scopeOfExternal(r.External)},
		{"codepaths", scopeOfCodepaths(r.Codepaths)},
		{"diagrams", scopeOfDiagrams(r.Diagrams)},
		{"versions", scopeOfVersions(r.Versions)},
		{"pins", scopeOf(r.Pins)},
		{"immutable", scopeOfImmutable(r.Immutable)},
		{"vcs", scopeOfVcs(r.Vcs)},
		{"commits", scopeOfCommits(r.Commits)},
	}
	for _, sc := range scopes {
		if sc.scope == nil {
			continue
		}
		if sc.scope.Roots == nil {
			return fmt.Errorf("%s: %s.scope.roots fehlt — scope ersetzt den globalen Scan und braucht explizite Wurzeln (leere Liste prüft nichts)", FileName, sc.module)
		}
		if cfg.Scopes == nil {
			cfg.Scopes = map[string]*model.ScopeConfig{}
		}
		cfg.Scopes[sc.module] = &model.ScopeConfig{Roots: sc.scope.Roots, Ignore: sc.scope.Ignore}
	}
	return nil
}

// scopeOf-Helfer: nil-sichere Extraktion des scope-Schluessels der
// jeweiligen Modul-Sektion.
func scopeOf(v *rawScopeOnly) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfIDs(v *rawIDs) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfMatrix(v *rawMatrix) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfExternal(v *rawExternal) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfHostpaths(v *struct {
	Scope    *rawScope `yaml:"scope"`
	Prefixes []string  `yaml:"prefixes"`
}) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfCodepaths(v *rawCodepaths) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfDiagrams(v *rawDiagrams) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfVersions(v *rawVersions) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfImmutable(v *rawImmutable) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfVcs(v *rawVCS) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

func scopeOfCommits(v *rawCommits) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

// decodeStrict dekodiert mit KnownFields; nil-raw bei leerem Dokument
// (Review R1/B2).
func decodeStrict(content []byte) (*raw, error) {
	var r raw
	dec := yaml.NewDecoder(bytes.NewReader(content))
	dec.KnownFields(true)
	if err := dec.Decode(&r); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		// interne Go-Typnamen aus der yaml-Meldung entfernen
		// (Review R1/B5).
		typeLeak := regexp.MustCompile(` in type \S+`)
		return nil, fmt.Errorf("%s: %s", FileName, typeLeak.ReplaceAllString(err.Error(), ""))
	}
	return &r, nil
}

func applyIDs(ids *rawIDs, cfg *model.Config) error {
	if ids == nil {
		return nil
	}
	for i, p := range ids.Patterns {
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			return fmt.Errorf("%s: ids.patterns[%d].regex: %v", FileName, i, err)
		}
		if re.MatchString("") {
			return fmt.Errorf("%s: ids.patterns[%d].regex matcht den Leerstring", FileName, i)
		}
		if p.Target == "" {
			return fmt.Errorf("%s: ids.patterns[%d].target fehlt", FileName, i)
		}
		if p.LinkPolicy != "" && p.LinkPolicy != "prose" && p.LinkPolicy != model.AlwaysPolicy {
			return fmt.Errorf("%s: ids.patterns[%d].link-policy: ungültig %q (erlaubt: prose, always)", FileName, i, p.LinkPolicy)
		}
		cfg.IDPatterns = append(cfg.IDPatterns, model.IDPattern{
			Regex: re, Target: p.Target,
			LinkPolicy: p.LinkPolicy, ExemptPaths: p.ExemptPaths,
		})
	}
	return nil
}

// validateMatrixDirection erzwingt die fail-closed-Kopplung von
// order/direction (DC-FA-MTX-002): direction nur DirectionNoDownward,
// direction und order ausschließlich gemeinsam — eine Richtungs-Deklaration
// darf nicht still wirkungslos sein.
func validateMatrixDirection(i int, direction string, order []string) error {
	switch {
	case direction == "" && len(order) == 0:
		return nil
	case direction == "":
		return fmt.Errorf("%s: matrix.classes[%d].order ohne direction", FileName, i)
	case direction != model.DirectionNoDownward:
		return fmt.Errorf("%s: matrix.classes[%d].direction %q unbekannt — nur %q",
			FileName, i, direction, model.DirectionNoDownward)
	case len(order) == 0:
		return fmt.Errorf("%s: matrix.classes[%d].direction ohne order", FileName, i)
	default:
		return nil
	}
}

// compileMatrixToken kompiliert das optionale token-Regex einer Klasse
// (DC-FA-MTX-003): leer ⇒ nil (nur Link-Erkennung); sonst muss es kompilieren
// und darf den Leerstring nicht matchen (sonst still jede Zeile, Exit 2).
func compileMatrixToken(i int, token string) (*regexp.Regexp, error) {
	if token == "" {
		return nil, nil
	}
	re, err := regexp.Compile(token)
	if err != nil {
		return nil, fmt.Errorf("%s: matrix.classes[%d].token kompiliert nicht: %w", FileName, i, err)
	}
	if re.MatchString("") {
		return nil, fmt.Errorf("%s: matrix.classes[%d].token matcht den Leerstring", FileName, i)
	}
	return re, nil
}

func applyMatrix(m *rawMatrix, cfg *model.Config) error {
	if m == nil {
		return nil
	}
	classes := map[string]bool{}
	for i, c := range m.Classes {
		if c.Name == "" || classes[c.Name] {
			return fmt.Errorf("%s: matrix.classes[%d].name fehlt oder doppelt", FileName, i)
		}
		if err := validateMatrixDirection(i, c.Direction, c.Order); err != nil {
			return err
		}
		token, err := compileMatrixToken(i, c.Token)
		if err != nil {
			return err
		}
		classes[c.Name] = true
		cfg.Matrix.Classes = append(cfg.Matrix.Classes, model.MatrixClass{
			Name: c.Name, Paths: c.Paths, Order: c.Order, Direction: c.Direction, Token: token,
		})
	}
	for i, rule := range m.Rules {
		if !classes[rule.From] || !classes[rule.To] {
			return fmt.Errorf("%s: matrix.rules[%d] referenziert undeklarierte Klasse", FileName, i)
		}
		cfg.Matrix.Rules = append(cfg.Matrix.Rules, model.MatrixRule{From: rule.From, To: rule.To, Allow: rule.Allow})
	}
	if m.Status != nil {
		cfg.Matrix.StatusForbidden = m.Status.Forbidden
		cfg.Matrix.AllowSupersedeLineage = m.Status.AllowSupersedeLineage
		for i, f := range m.Status.SupersedeFields {
			if strings.TrimSpace(f) == "" {
				return fmt.Errorf("%s: matrix.status.supersede-fields[%d] ist leer", FileName, i)
			}
		}
		cfg.Matrix.SupersedeFields = m.Status.SupersedeFields
	} else {
		cfg.Matrix.StatusForbidden = []string{"superseded", "deprecated"}
	}
	cfg.Matrix.ExcludeSections = m.ExcludeSections
	cfg.Matrix.ExemptPaths = m.ExemptPaths
	return nil
}

func applyExternal(e *rawExternal, cfg *model.Config) error {
	if e == nil {
		return nil
	}
	if t := e.TimeoutSeconds; t != nil {
		if *t < 1 || *t > 300 {
			return fmt.Errorf("%s: external.timeout-seconds außerhalb 1–300", FileName)
		}
		cfg.External.TimeoutSeconds = *t
	}
	if p := e.Parallel; p != nil {
		if *p < 1 || *p > 16 {
			return fmt.Errorf("%s: external.parallel außerhalb 1–16", FileName)
		}
		cfg.External.Parallel = *p
	}
	return nil
}
