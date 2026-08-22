// Package configyaml ist der driven Config-Adapter: striktes Decoding
// und statische Validierung von .d-check.yml (DC-FA-CONF-001;
// Strategie gemäß ADR-0003). Den Datei-Inhalt beschafft das CLI über
// den Filesystem-Adapter — dieser Adapter macht kein I/O
// (spec/architecture.md §2). Die Existenz-Prüfung deklarierter
// Scan-Wurzeln erfolgt im Kern beim Lauf-Start (ebenfalls Exit 2).
package configyaml

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"gopkg.in/yaml.v3"
)

// FileName ist der konventionelle Name der Konfigurationsdatei.
const FileName = ".d-check.yml"

// RetargetFileName schreibt den Datei-Präfix einer Decoder-Fehlermeldung auf die
// tatsächlich geladene Datei um. Alle Meldungen dieses Pakets beginnen mit
// `FileName + ":"`, weil der Decoder kein I/O macht und den Pfad nicht kennt
// (spec/architecture.md §2). Wird die Konfiguration per --config aus einer
// anderen Datei gelesen (DC-FA-CLI-012), zeigte die Meldung sonst auf eine
// Datei, die der Lauf nie gelesen hat — dieselbe Provenance-Ehrlichkeit, die
// die Config-Pin-Befunde des Moduls sources tragen. Nur das erste Vorkommen
// wird ersetzt: der Präfix, nicht ein zufällig zitierter Dateiname im Rest.
func RetargetFileName(err error, name string) error {
	if err == nil || name == "" || name == FileName {
		return err
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, FileName+":") {
		return err
	}
	return errors.New(name + strings.TrimPrefix(msg, FileName))
}

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
// Konfigurations-Schluessel haben (anchors, spans, pins).
type rawScopeOnly struct {
	Scope *rawScope `yaml:"scope"`
}

// rawLinks traegt scope und die resolve-from-Gruppen (DC-FA-LINK-001
// §Ortsfeste Verweise).
type rawLinks struct {
	Scope       *rawScope         `yaml:"scope"`
	ResolveFrom []rawResolveGroup `yaml:"resolve-from"`
}

// rawResolveGroup ist eine Gruppe ortsfester Verweise: dirs sind die
// wandernden Verzeichnisse (Quellen), fixed-dirs die ortsfesten Ziele.
type rawResolveGroup struct {
	Dirs      []string `yaml:"dirs"`
	FixedDirs []string `yaml:"fixed-dirs"`
}

// rawCodepaths traegt scope, roots und das exempt-paths-Ventil
// (DC-FA-CODE-001; benannt, damit der Typ in raw und scopeOfCodepaths
// identisch ist).
type rawCodepaths struct {
	Scope       *rawScope `yaml:"scope"`
	Roots       []string  `yaml:"roots"`
	ExemptPaths []string  `yaml:"exempt-paths"`
	IgnoreRefs  []string  `yaml:"ignore-refs"`
	CheckLines  bool      `yaml:"check-lines"`
}

// rawIgnoreRef traegt einen Eintrag des geteilten Referenz-Ventils
// ignore-refs (DC-FA-REF-001): in (Quell-Skopus), refs (Ziel-Globs),
// keep (Ausnahmen).
type rawIgnoreRef struct {
	In   string   `yaml:"in"`
	Refs []string `yaml:"refs"`
	Keep []string `yaml:"keep"`
}

// rawDiagrams trägt scope, fences, die Muster und das Datei-Ventil des Moduls
// diagrams (DC-FA-DIAG-001); das zweite Ventil ist der Zeilen-Marker und lebt
// im Kern.
type rawDiagrams struct {
	Scope       *rawScope           `yaml:"scope"`
	Fences      []string            `yaml:"fences"`
	Patterns    []rawDiagramPattern `yaml:"patterns"`
	ExemptPaths []string            `yaml:"exempt-paths"`
}

type rawDiagramPattern struct {
	Regex     string `yaml:"regex"`
	DefinedIn string `yaml:"defined-in"`
}

// rawVersions trägt scope und die Muster-Quellen-Paare des Moduls versions
// (DC-FA-VER-001) in beiden Schreibweisen: die Kurzform direkt (pin-pattern,
// current-from, exempt-paths) und die Liste unter patterns. Beide zugleich
// sind ein Nutzungsfehler.
type rawVersions struct {
	Scope *rawScope `yaml:"scope"`
	// Zeiger, weil die Mischform an der ANWESENHEIT eines Schlüssels hängt,
	// nicht an seinem Wert: `pin-pattern: ""` neben patterns ist eine
	// Mischform, ein fehlender Schlüssel nicht.
	PinPattern  *string             `yaml:"pin-pattern"`
	CurrentFrom *string             `yaml:"current-from"`
	ExemptPaths *[]string           `yaml:"exempt-paths"`
	Patterns    []rawVersionPattern `yaml:"patterns"`
}

// rawVersionPattern ist ein Paar der versions.patterns-Liste: eigenes Muster,
// eigene Quelle, eigene Ausnahmen (DC-FA-VER-001).
type rawVersionPattern struct {
	PinPattern  string   `yaml:"pin-pattern"`
	CurrentFrom string   `yaml:"current-from"`
	ExemptPaths []string `yaml:"exempt-paths"`
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
	Roadmap   string       `yaml:"roadmap"`
	Heading   string       `yaml:"heading"`
	Marker    string       `yaml:"marker"`
	SliceGlob string       `yaml:"slice-glob"`
	Closure   *rawClosure  `yaml:"closure"`
	Waves     *rawWaves    `yaml:"waves"`
}

// rawWaves trägt die dritte planning-Fähigkeit (DC-FA-PLAN-001
// §Wellen-Invariante): dir ist der Aktivierungs-Schalter (leer/abwesend ⇒
// inert), alles Weitere hat Konventions-Defaults im Kern.
// Die vier optionalen Felder sind Zeiger, damit ein ABWESENDER Schlüssel von
// einem explizit gesetzten unterscheidbar bleibt: abwesend heißt Kern-Default,
// explizit leer ist eine Null-Aussage und bricht mit Exit 2 ab (W1).
type rawWaves struct {
	Dir           string  `yaml:"dir"`
	DoneDir       string  `yaml:"done-dir"`
	Glob          *string `yaml:"glob"`
	ResultsGlob   *string `yaml:"results-glob"`
	NextHeading   *string `yaml:"next-heading"`
	ClosedHeading *string `yaml:"closed-heading"`
	Mode          *string `yaml:"mode"`
}

// rawClosure trägt die zweite planning-Fähigkeit (DC-FA-PLAN-001
// §Closure-Note-Struktur): dir ist der Aktivierungs-Schalter (leer/abwesend ⇒
// inert), heading-pattern/min-sentences/boilerplate sind optional mit
// Konventions-Defaults im Kern.
// MinSentences und Glob sind Zeiger, damit ein ABWESENDER Schlüssel von einem
// explizit gesetzten Wert unterscheidbar bleibt: abwesend heißt Kern-Default
// (4 bzw. der Verweis auf slice-glob), explizit `0` bzw. `""` wäre eine
// Null-Aussage und damit ein stilles Grün, das der Vertrag mit Exit 2 ablehnt.
type rawClosure struct {
	Dir            string   `yaml:"dir"`
	Glob           *string  `yaml:"glob"`
	HeadingPattern string   `yaml:"heading-pattern"`
	MinSentences   *int     `yaml:"min-sentences"`
	Boilerplate    []string `yaml:"boilerplate"`
	Placeholder    bool     `yaml:"placeholder"`
}

// rawStructure ist eine Regel des Moduls structure (DC-FA-STRUCT-001).
// MinSentences/MaxTasks sind Zeiger, damit ein ABWESENDER Schlüssel
// (Bedingung aus) von einem explizit gesetzten Wert unterscheidbar bleibt.
type rawStructure struct {
	Files          string   `yaml:"files"`
	Section        string   `yaml:"section"`
	SectionPattern string   `yaml:"section-pattern"`
	Sections       string   `yaml:"sections"`
	NonEmpty       bool     `yaml:"non-empty"`
	MinSentences   *int     `yaml:"min-sentences"`
	MaxTasks       *int     `yaml:"max-tasks"`
	ForbidPattern  string   `yaml:"forbid-pattern"`
	RequirePattern string   `yaml:"require-pattern"`
	RequireAll     []string `yaml:"require-all"`
	TableOrder     string   `yaml:"table-order"`
	TableColumn    *int     `yaml:"table-column"`
	HeadingsMatch  string   `yaml:"headings-match"`
	HeadingsLevel  *int     `yaml:"headings-level"`
	ExemptPaths    []string `yaml:"exempt-paths"`
}

// applyStructure validiert die Regel-Liste am Config-Rand (§DC-FA-STRUCT-001.a
// Schritt 1): was der Kern nur schlucken könnte, bricht hier laut ab. Eine
// doppelte Regel-Identität ist ebenfalls Exit 2 — sonst fielen zwei leer
// laufende Regeln unter der Befund-Deduplikation zusammen.
func applyStructure(rs []rawStructure) ([]model.StructureRule, error) {
	out := make([]model.StructureRule, 0, len(rs))
	seen := map[string]bool{}
	for i, r := range rs {
		rule, err := applyStructureRule(i, r)
		if err != nil {
			return nil, err
		}
		if seen[rule.Identity()] {
			return nil, fmt.Errorf("%s: structure[%d]: Regel-Identität %q kommt doppelt vor", FileName, i, rule.Identity())
		}
		seen[rule.Identity()] = true
		out = append(out, rule)
	}
	return out, nil
}

// structureBedingungsFehler prüft die Bedingungs-Parameter einer Regel und
// liefert die Fehlermeldung (leer ⇒ gültig). Ausgelagert, damit
// applyStructureRule unter der gocyclo-Schwelle bleibt.
func structureBedingungsFehler(r rawStructure) string {
	// Geordneter Slice statt Map: bei zwei ungueltigen Mustern in derselben
	// Regel entschiede sonst die Iterationsreihenfolge, welche Meldung der
	// Nutzer sieht — dieselbe Eingabe muss dieselbe Ausgabe liefern (DC-QA-02).
	for _, m := range []struct{ name, pat string }{
		{"section-pattern", r.SectionPattern}, {"forbid-pattern", r.ForbidPattern},
		{"require-pattern", r.RequirePattern}, {"headings-match", r.HeadingsMatch},
	} {
		if m.pat == "" {
			continue
		}
		if _, err := regexp.Compile(m.pat); err != nil {
			return fmt.Sprintf("%s %q ist kein gültiges Regex: %v", m.name, m.pat, err)
		}
	}
	if r.MinSentences != nil && *r.MinSentences < 1 {
		return fmt.Sprintf("min-sentences %d muss >= 1 sein", *r.MinSentences)
	}
	if r.MaxTasks != nil && *r.MaxTasks < 0 {
		return fmt.Sprintf("max-tasks %d muss >= 0 sein", *r.MaxTasks)
	}
	for _, m := range r.RequireAll {
		if strings.TrimSpace(m) == "" {
			return "require-all enthält einen leeren Eintrag"
		}
	}
	return structureChronologieFehler(r)
}

// structureChronologieFehler prüft die drei Config-Ränder der
// Chronologie-Bedingung (ADR-0057): eine halbe Aktivierung ist ein
// Config-Fehler, kein Zustand — table-column ohne table-order bricht laut.
func structureChronologieFehler(r rawStructure) string {
	switch r.TableOrder {
	case "", "asc", "desc":
	default:
		return fmt.Sprintf("table-order %q muss asc oder desc sein", r.TableOrder)
	}
	if r.TableColumn != nil && *r.TableColumn < 1 {
		return fmt.Sprintf("table-column %d muss >= 1 sein", *r.TableColumn)
	}
	if r.TableColumn != nil && r.TableOrder == "" {
		return "table-column ist ohne table-order wirkungslos (halbe Aktivierung)"
	}
	return structureUeberschriftFehler(r)
}

// structureUeberschriftFehler prüft die beiden Config-Ränder der
// Überschriften-Bedingung: eine Ebene außerhalb des ATX-Bereichs, und —
// analog zu table-column — die halbe Aktivierung.
func structureUeberschriftFehler(r rawStructure) string {
	if r.HeadingsLevel != nil && (*r.HeadingsLevel < 1 || *r.HeadingsLevel > 6) {
		return fmt.Sprintf("headings-level %d muss zwischen 1 und 6 liegen", *r.HeadingsLevel)
	}
	if r.HeadingsLevel != nil && r.HeadingsMatch == "" {
		return "headings-level ist ohne headings-match wirkungslos (halbe Aktivierung)"
	}
	return ""
}

func applyStructureRule(i int, r rawStructure) (model.StructureRule, error) {
	fail := func(f string, a ...any) (model.StructureRule, error) {
		return model.StructureRule{}, fmt.Errorf("%s: structure[%d]: "+f, append([]any{FileName, i}, a...)...)
	}
	if r.Files == "" {
		return fail("files ist Pflicht")
	}
	if _, err := path.Match(r.Files, "probe"); err != nil {
		return fail("files %q ist kein gültiges Glob: %v", r.Files, err)
	}
	for _, g := range r.ExemptPaths {
		if _, err := path.Match(g, "probe"); err != nil {
			return fail("exempt-paths %q ist kein gültiges Glob: %v", g, err)
		}
	}
	if (r.Section == "") == (r.SectionPattern == "") {
		return fail("genau eines von section oder section-pattern ist Pflicht")
	}
	switch r.Sections {
	case "", "one", "each":
	default:
		return fail("sections %q muss one oder each sein", r.Sections)
	}
	if msg := structureBedingungsFehler(r); msg != "" {
		return fail("%s", msg)
	}
	return model.StructureRule{
		Files: r.Files, Section: r.Section, SectionPattern: r.SectionPattern,
		Sections: r.Sections, NonEmpty: r.NonEmpty, MinSentences: r.MinSentences,
		MaxTasks: r.MaxTasks, ForbidPattern: r.ForbidPattern,
		RequirePattern: r.RequirePattern, RequireAll: r.RequireAll,
		TableOrder: r.TableOrder, TableColumn: r.TableColumn,
		HeadingsMatch: r.HeadingsMatch, HeadingsLevel: r.HeadingsLevel,
		ExemptPaths: r.ExemptPaths,
	}, nil
}

// rawTracked trägt scope und die Parameter des Moduls tracked
// (DC-FA-TRK-001): exempt-targets (Globs über aufgelöste Ziel-Pfade,
// referenz-weit — absichtlich untrackte Ziele). scope gilt für die
// Quell-Dateien (DC-FA-CONF-002) wie bei jedem scannenden Modul.
type rawTracked struct {
	Scope         *rawScope `yaml:"scope"`
	ExemptTargets []string  `yaml:"exempt-targets"`
}

// rawTargets trägt die Parameter des Moduls targets (DC-FA-TGT-001): makefiles
// (Regelnamen-Quellen), doc-tables (make-X-Tabellen für Richtung 1), authority
// (Vollständigkeits-Quelle für Richtung 2), exempt-targets (Utility-Regeln ohne
// Doku-Pflicht). **Keine** scope — targets ist ein Post-Pass ohne Datei-Scan
// (wie planning; ein targets.scope wäre wirkungslos, der strikte Decoder lehnt
// es ab).
type rawTargets struct {
	Makefiles     []string `yaml:"makefiles"`
	DocTables     []string `yaml:"doc-tables"`
	Authority     string   `yaml:"authority"`
	ExemptTargets []string `yaml:"exempt-targets"`
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

// rawSource ist ein Config-Pin des Moduls sources (DC-FA-SRC-001): url als
// yaml.Node, damit die Befund-Zeile (Node-Zeile des url-Feldes) erhalten bleibt.
type rawSource struct {
	URL    yaml.Node `yaml:"url"`
	Sha256 string    `yaml:"sha256"`
	Unpack string    `yaml:"unpack"`
}

// rawTrace bildet den opt-in trace-Block ab (DC-FA-CLI-009): konfigurierbare
// RTM-Quellen. Jeder Unter-Block/jedes Feld ist optional; abwesend ⇒ Default im
// Kern (byte-identisch, DC-QA-02).
type rawTrace struct {
	Requirements *rawTraceRequirements `yaml:"requirements"`
	ADRs         *struct {
		Dir         string `yaml:"dir"`
		FilePattern string `yaml:"file-pattern"`
		IDPrefix    string `yaml:"id-prefix"`
	} `yaml:"adrs"`
	Slices *struct {
		Dir         string `yaml:"dir"`
		FilePattern string `yaml:"file-pattern"`
		IDPrefix    string `yaml:"id-prefix"`
	} `yaml:"slices"`
	Coverage         []rawCoverage        `yaml:"coverage"`
	CrossConsistency *rawCrossConsistency `yaml:"cross-consistency"`
}

// rawCrossConsistency bildet den opt-in Kreuzverweis-Abgleich ab
// (DC-FA-XREF-001); forward und backward sind beide Pflicht.
type rawCrossConsistency struct {
	Forward    *rawCrossForward  `yaml:"forward"`
	Backward   *rawCrossBackward `yaml:"backward"`
	Mode       string            `yaml:"mode"`
	ExcludeReq string            `yaml:"exclude-req"`
}

// rawCrossForward ist die Vorwärts-Sicht (Anforderung → Design).
type rawCrossForward struct {
	File            string   `yaml:"file"`
	Sections        []string `yaml:"sections"`
	ExcludeSections []string `yaml:"exclude-sections"`
	ReqColumn       string   `yaml:"req-column"`
	DesignColumn    string   `yaml:"design-column"`
	DesignPattern   string   `yaml:"design-pattern"`
	ReqPattern      string   `yaml:"req-pattern"`
	Ranges          *bool    `yaml:"ranges"`
}

// rawCrossBackward ist die Rück-Kanten-Sicht (Design → Anforderung).
type rawCrossBackward struct {
	File             string   `yaml:"file"`
	Sections         []string `yaml:"sections"`
	ArtifactIDColumn string   `yaml:"artifact-id-column"`
	EdgeColumn       string   `yaml:"edge-column"`
	ReqPattern       string   `yaml:"req-pattern"`
	Ranges           *bool    `yaml:"ranges"`
}

type rawTraceRequirements struct {
	Source    string         `yaml:"source"`
	IDPattern string         `yaml:"id-pattern"`
	Format    string         `yaml:"format"`
	Table     *rawTraceTable `yaml:"table"`
	Modality  *rawModality   `yaml:"modality"`
}

// rawTraceTable bindet tabellarische Anforderungsquellen an exakte
// Header-Namen (DC-FA-REQ-001).
type rawTraceTable struct {
	IDColumn       string    `yaml:"id-column"`
	TextColumn     *string   `yaml:"text-column"`
	TextColumns    *[]string `yaml:"text-columns"`
	ModalityColumn string    `yaml:"modality-column"`
	DuplicateIDs   string    `yaml:"duplicate-ids"`
}

// rawModality bildet den opt-in Modalitäts-Block ab (DC-FA-MOD-001).
type rawModality struct {
	Levels        map[string][]string `yaml:"levels"`
	RequireLevels []string            `yaml:"require-levels"`
}

// rawCoverage bildet eine kuratierte Coverage-Quelle ab (DC-FA-COV-001).
type rawCoverage struct {
	Files           []string `yaml:"files"`
	Label           string   `yaml:"label"`
	Ranges          *bool    `yaml:"ranges"`
	Sections        []string `yaml:"sections"`
	ExcludeSections []string `yaml:"exclude-sections"`
}

// raw bildet das Voll-Schema von .d-check.yml ab
// (spec/spezifikation.md §2). Unbekannte Schlüssel sind durch
// KnownFields(true) Fehler.
type raw struct {
	Scan *struct {
		Roots  []string `yaml:"roots"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"scan"`
	IgnoreRefs []rawIgnoreRef `yaml:"ignore-refs"`
	Modules    []string       `yaml:"modules"`
	Links     *rawLinks     `yaml:"links"`
	Anchors   *rawScopeOnly `yaml:"anchors"`
	Spans     *rawScopeOnly `yaml:"spans"`
	Pins      *rawScopeOnly `yaml:"pins"`
	Hostpaths *struct {
		Scope    *rawScope `yaml:"scope"`
		Prefixes []string  `yaml:"prefixes"`
	} `yaml:"hostpaths"`
	IDs       *rawIDs       `yaml:"ids"`
	Matrix    *rawMatrix    `yaml:"matrix"`
	External  *rawExternal  `yaml:"external"`
	Codepaths *rawCodepaths `yaml:"codepaths"`
	Diagrams  *rawDiagrams  `yaml:"diagrams"`
	Versions  *rawVersions  `yaml:"versions"`
	Immutable *rawImmutable `yaml:"immutable"`
	Vcs       *rawVCS       `yaml:"vcs"`
	Commits   *rawCommits   `yaml:"commits"`
	Planning  *rawPlanning  `yaml:"planning"`
	Tracked   *rawTracked   `yaml:"tracked"`
	Structure []rawStructure `yaml:"structure"`
	Targets   *rawTargets   `yaml:"targets"`
	// Sources ist eine bare Liste `sources[]` (spec/spezifikation.md §2;
	// KEIN Map mit scope — das Modul nutzt den globalen Scan-Scope).
	Sources []rawSource `yaml:"sources"`
	Trace   *rawTrace   `yaml:"trace"`
	// Citations ist parameterlos (direktiven-getrieben) — nur scope.
	Citations *rawScopeOnly `yaml:"citations"`
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

// applyResolveFrom validiert die resolve-from-Gruppen (DC-FA-LINK-001
// §Ortsfeste Verweise, Schritt 6). Config-Rand-Disziplin: eine Gruppe aus
// weniger als zwei wandernden Orten prüft nichts, ein absoluter oder
// `..`-Pfad verlässt die Wurzel, und ein Verzeichnis in mehreren Gruppen
// machte die Quellen-Zuordnung mehrdeutig — alles Exit 2 statt stillem Grün.
func applyResolveFrom(l *rawLinks, cfg *model.Config) error {
	if l == nil || len(l.ResolveFrom) == 0 {
		return nil
	}
	gesehen := map[string]bool{}
	for i, g := range l.ResolveFrom {
		if len(g.Dirs) < 2 {
			return fmt.Errorf(
				"%s: links.resolve-from[%d].dirs braucht mindestens zwei wandernde Verzeichnisse (eine Gruppe aus einem Ort prüft nichts)", FileName, i)
		}
		for _, d := range append(append([]string{}, g.Dirs...), g.FixedDirs...) {
			if d == "" || strings.HasPrefix(d, "/") || strings.Contains(d, "..") {
				return fmt.Errorf(
					"%s: links.resolve-from[%d]: Verzeichnis %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..', nicht leer)", FileName, i, d)
			}
		}
		for _, d := range g.Dirs {
			c := path.Clean(d)
			if gesehen[c] {
				return fmt.Errorf(
					"%s: links.resolve-from: Verzeichnis %q ist dirs-Mitglied mehrerer Gruppen — die Quellen-Zuordnung wäre mehrdeutig", FileName, d)
			}
			gesehen[c] = true
		}
		cfg.ResolveFrom = append(cfg.ResolveFrom, model.ResolveFromGroup{
			Dirs: g.Dirs, FixedDirs: g.FixedDirs,
		})
	}
	return nil
}

// applyModules wendet die modul-spezifischen Validierungen/Kompilierungen in
// fester Reihenfolge an und liefert den ersten Fehler (Exit 2). Ausgelagert aus
// Decode, damit dessen gocyclo-Komplexität unter der Schwelle bleibt.
func applyModules(r *raw, cfg *model.Config) error {
	if err := applyIgnoreRefs(r, cfg); err != nil {
		return err
	}
	if err := applyResolveFrom(r.Links, cfg); err != nil {
		return err
	}
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
	return applyRemainingModules(r, cfg)
}

// applyRemainingModules setzt die Modul-Validierung fort (ausgelagert, damit
// applyModules unter der gocyclo-Schwelle bleibt).
func applyRemainingModules(r *raw, cfg *model.Config) error {
	if err := applyVCS(r, cfg); err != nil {
		return err
	}
	if err := applyCommits(r, cfg); err != nil {
		return err
	}
	st, serr := applyStructure(r.Structure)
	if serr != nil {
		return serr
	}
	cfg.Structure = st
	if err := applyPlanning(r, cfg); err != nil {
		return err
	}
	if err := applyTracked(r, cfg); err != nil {
		return err
	}
	if err := applyTargets(r, cfg); err != nil {
		return err
	}
	if err := applySources(r, cfg); err != nil {
		return err
	}
	if err := applyTrace(r, cfg); err != nil {
		return err
	}
	return applyScopes(r, cfg)
}

// compileTracePattern kompiliert ein trace-Muster (DC-FA-CLI-009). Leer ⇒ nil
// (Default im Kern). Ein nicht kompilierbares Muster ist ein Fehler (Exit 2);
// requireCapture ⇒ das Muster muss mindestens eine Capture-Gruppe tragen, sonst
// wäre die Owner-Kennung (m[1]) undefiniert (fail-closed, verhindert ein
// Laufzeit-Panic).
func compileTracePattern(field, pattern string, requireCapture bool) (*regexp.Regexp, error) {
	if strings.TrimSpace(pattern) == "" {
		return nil, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %v", FileName, field, err)
	}
	if requireCapture && re.NumSubexp() < 1 {
		return nil, fmt.Errorf("%s: %s braucht mindestens eine Capture-Gruppe '(…)' für die Owner-Kennung", FileName, field)
	}
	return re, nil
}

// validateTracePath prüft, dass ein trace-Pfad relativ zur Repo-Wurzel liegt
// (kein führendes '/', kein '..') — analog planning.roadmap. Leer ⇒ zulässig
// (Default im Kern). Die Existenz prüft der Kern beim Lauf; eine nichtleer
// explizite Requirements-Quelle ist dabei fail-closed (DC-FA-REQ-001.a).
func validateTracePath(field, p string) error {
	if strings.HasPrefix(p, "/") || strings.Contains(p, "..") {
		return fmt.Errorf("%s: %s %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, field, p)
	}
	return nil
}

// applyTrace validiert und kompiliert den opt-in trace-Block (DC-FA-CLI-009):
// je Unter-Block Pfad-Validierung + Regex-Kompilierung; die file-pattern der
// Referenzklassen brauchen eine Capture-Gruppe. Alles fail-closed (Exit 2). Ohne
// trace-Block bleibt cfg.Trace der Nullwert ⇒ Konventions-Default im Kern
// (byte-identisch, DC-QA-02).
func applyTrace(r *raw, cfg *model.Config) error {
	if r.Trace == nil {
		return nil
	}
	t := r.Trace
	var tc model.TraceConfig
	if req := t.Requirements; req != nil {
		if err := applyTraceRequirements(req, &tc); err != nil {
			return err
		}
	}
	if adr := t.ADRs; adr != nil {
		if err := validateTracePath("trace.adrs.dir", adr.Dir); err != nil {
			return err
		}
		re, err := compileTracePattern("trace.adrs.file-pattern", adr.FilePattern, true)
		if err != nil {
			return err
		}
		tc.ADRDir, tc.ADRFile, tc.ADRPrefix = adr.Dir, re, adr.IDPrefix
	}
	if sl := t.Slices; sl != nil {
		if err := validateTracePath("trace.slices.dir", sl.Dir); err != nil {
			return err
		}
		re, err := compileTracePattern("trace.slices.file-pattern", sl.FilePattern, true)
		if err != nil {
			return err
		}
		tc.SliceDir, tc.SliceFile, tc.SlicePrefix = sl.Dir, re, sl.IDPrefix
	}
	if err := applyTraceCoverage(t.Coverage, &tc); err != nil {
		return err
	}
	if err := applyCrossConsistency(t.CrossConsistency, &tc); err != nil {
		return err
	}
	cfg.Trace = tc
	return nil
}

// applyCrossConsistency validiert den opt-in Kreuzverweis-Block (DC-FA-XREF-001)
// und kompiliert seine Muster. Fail-closed (Exit 2): fehlender Pflichtblock,
// unbekannter mode, nicht kompilierendes Regex, leeres Pflichtfeld. Abwesend ⇒
// kein Abgleich, RTM byte-identisch (DC-QA-02). Datei-Existenz und Header-Bindung
// prüft der Kern beim Lauf (config-zeitig erfolgt kein I/O).
func applyCrossConsistency(raw *rawCrossConsistency, tc *model.TraceConfig) error {
	if raw == nil {
		return nil
	}
	if raw.Forward == nil || raw.Backward == nil {
		return fmt.Errorf("%s: trace.cross-consistency braucht forward und backward", FileName)
	}
	mode := raw.Mode
	if mode == "" {
		mode = model.TraceCrossModeEqual
	}
	if mode != model.TraceCrossModeEqual && mode != model.TraceCrossModeSuperset {
		return fmt.Errorf("%s: trace.cross-consistency.mode %q ist ungültig (erwartet: equal oder superset)", FileName, mode)
	}
	excludeReq, err := compileTracePattern("trace.cross-consistency.exclude-req", raw.ExcludeReq, false)
	if err != nil {
		return err
	}
	forward, err := applyCrossForward(raw.Forward)
	if err != nil {
		return err
	}
	backward, err := applyCrossBackward(raw.Backward)
	if err != nil {
		return err
	}
	tc.CrossConsistency = &model.TraceCrossConsistency{
		Forward: forward, Backward: backward, Mode: mode, ExcludeReq: excludeReq,
	}
	return nil
}

// applyCrossForward validiert die Vorwärts-Sicht: alle vier Pflichtfelder gesetzt,
// Datei innerhalb der Repo-Wurzel, design-pattern kompilierbar; ranges-Default true.
func applyCrossForward(raw *rawCrossForward) (model.TraceCrossForward, error) {
	var out model.TraceCrossForward
	if err := requireCrossFields(map[string]string{
		"trace.cross-consistency.forward.file":           raw.File,
		"trace.cross-consistency.forward.req-column":     raw.ReqColumn,
		"trace.cross-consistency.forward.design-column":  raw.DesignColumn,
		"trace.cross-consistency.forward.design-pattern": raw.DesignPattern,
	}); err != nil {
		return out, err
	}
	if err := validateTracePath("trace.cross-consistency.forward.file", raw.File); err != nil {
		return out, err
	}
	designPat, err := compileTracePattern("trace.cross-consistency.forward.design-pattern", raw.DesignPattern, false)
	if err != nil {
		return out, err
	}
	// Leer ⇒ nil ⇒ Default `requirements.id-pattern` im Kern (DC-FA-XREF-001).
	reqPat, err := compileTracePattern("trace.cross-consistency.forward.req-pattern", raw.ReqPattern, false)
	if err != nil {
		return out, err
	}
	return model.TraceCrossForward{
		File: raw.File, Sections: raw.Sections, ExcludeSections: raw.ExcludeSections,
		ReqColumn: raw.ReqColumn, DesignColumn: raw.DesignColumn,
		DesignPattern: designPat, ReqPattern: reqPat, Ranges: crossRanges(raw.Ranges),
	}, nil
}

// applyCrossBackward validiert die Rück-Sicht; artifact-id-column ist per Default
// der positionelle Sentinel `first` (heterogene ID-Header, ADR-0038).
func applyCrossBackward(raw *rawCrossBackward) (model.TraceCrossBackward, error) {
	var out model.TraceCrossBackward
	if err := requireCrossFields(map[string]string{
		"trace.cross-consistency.backward.file":        raw.File,
		"trace.cross-consistency.backward.edge-column": raw.EdgeColumn,
		"trace.cross-consistency.backward.req-pattern": raw.ReqPattern,
	}); err != nil {
		return out, err
	}
	if err := validateTracePath("trace.cross-consistency.backward.file", raw.File); err != nil {
		return out, err
	}
	reqPat, err := compileTracePattern("trace.cross-consistency.backward.req-pattern", raw.ReqPattern, false)
	if err != nil {
		return out, err
	}
	idColumn := raw.ArtifactIDColumn
	if idColumn == "" {
		idColumn = model.TraceCrossArtifactFirst
	}
	return model.TraceCrossBackward{
		File: raw.File, Sections: raw.Sections, ArtifactIDColumn: idColumn,
		EdgeColumn: raw.EdgeColumn, ReqPattern: reqPat, Ranges: crossRanges(raw.Ranges),
	}, nil
}

// requireCrossFields meldet das erste leere Pflichtfeld — deterministisch über
// die sortierten Schlüssel (Map-Iteration ist es nicht).
func requireCrossFields(fields map[string]string) error {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if strings.TrimSpace(fields[name]) == "" {
			return fmt.Errorf("%s: %s ist leer", FileName, name)
		}
	}
	return nil
}

// crossRanges löst den ranges-Default auf (true, DC-FA-XREF-001); der Pointer
// unterscheidet nicht-gesetzt von explizit false.
func crossRanges(set *bool) bool {
	if set == nil {
		return true
	}
	return *set
}

func applyTraceRequirements(req *rawTraceRequirements, tc *model.TraceConfig) error {
	if err := validateTracePath("trace.requirements.source", req.Source); err != nil {
		return err
	}
	re, err := compileTracePattern("trace.requirements.id-pattern", req.IDPattern, false)
	if err != nil {
		return err
	}
	tc.Source, tc.ReqPattern = req.Source, re
	if err := applyTraceRequirementsFormat(req.Format, req.Table, tc); err != nil {
		return err
	}
	mod, err := applyModality(req.Modality)
	if err != nil {
		return err
	}
	tc.Modality = mod
	return nil
}

// applyTraceRequirementsFormat validiert die Format-/Block-Konsistenz vor
// jedem I/O (DC-FA-REQ-001): table braucht eindeutig benannte Rollen,
// headings akzeptiert keinen wirkungslosen table-Block.
func applyTraceRequirementsFormat(format string, table *rawTraceTable, tc *model.TraceConfig) error {
	switch format {
	case "", model.TraceFormatHeadings:
		if table != nil {
			return fmt.Errorf("%s: trace.requirements.table ist nur mit format: table zulässig", FileName)
		}
		tc.Format = format
		return nil
	case model.TraceFormatTable:
		if table == nil {
			return fmt.Errorf("%s: trace.requirements.format: table braucht einen table-Block", FileName)
		}
	default:
		return fmt.Errorf("%s: trace.requirements.format %q ist ungültig (erwartet: headings oder table)", FileName, format)
	}
	return applyTraceTableConfig(table, tc)
}

func applyTraceTableConfig(table *rawTraceTable, tc *model.TraceConfig) error {
	id := strings.TrimSpace(table.IDColumn)
	mod := strings.TrimSpace(table.ModalityColumn)
	if id == "" {
		return fmt.Errorf("%s: trace.requirements.table.id-column darf nicht leer sein", FileName)
	}
	texts, err := traceTextColumns(table)
	if err != nil {
		return err
	}
	if mod != "" && mod == id {
		return fmt.Errorf("%s: trace.requirements.table-Spaltennamen müssen verschieden sein", FileName)
	}
	for _, txt := range texts {
		if txt == id || mod != "" && txt == mod {
			return fmt.Errorf("%s: trace.requirements.table-Spaltennamen müssen verschieden sein", FileName)
		}
	}
	duplicateIDs := strings.TrimSpace(table.DuplicateIDs)
	if duplicateIDs == "" {
		duplicateIDs = model.TraceDuplicateError
	}
	if !validDuplicatePolicy(duplicateIDs) {
		return fmt.Errorf("%s: trace.requirements.table.duplicate-ids %q ist ungültig (erwartet: error, first oder last)", FileName, duplicateIDs)
	}
	tc.Format = model.TraceFormatTable
	tc.Table = &model.TraceTableConfig{IDColumn: id, TextColumns: texts, ModalityColumn: mod, DuplicateIDs: duplicateIDs}
	return nil
}

func validDuplicatePolicy(policy string) bool {
	return policy == model.TraceDuplicateError || policy == model.TraceDuplicateFirst || policy == model.TraceDuplicateLast
}

func traceTextColumns(table *rawTraceTable) ([]string, error) {
	singleSet := table.TextColumn != nil
	multipleSet := table.TextColumns != nil
	if singleSet && multipleSet {
		return nil, fmt.Errorf("%s: trace.requirements.table.text-column und text-columns sind alternativ, nicht gleichzeitig zulässig", FileName)
	}
	if singleSet {
		single := strings.TrimSpace(*table.TextColumn)
		if single == "" {
			return nil, fmt.Errorf("%s: trace.requirements.table.text-column darf nicht leer sein", FileName)
		}
		return []string{single}, nil
	}
	if !multipleSet {
		return nil, fmt.Errorf("%s: trace.requirements.table braucht text-column oder text-columns", FileName)
	}
	if len(*table.TextColumns) == 0 {
		return nil, fmt.Errorf("%s: trace.requirements.table.text-columns darf nicht leer sein", FileName)
	}
	texts := make([]string, 0, len(*table.TextColumns))
	seen := map[string]bool{}
	for _, raw := range *table.TextColumns {
		text := strings.TrimSpace(raw)
		if text == "" {
			return nil, fmt.Errorf("%s: trace.requirements.table.text-columns darf keinen leeren Header enthalten", FileName)
		}
		if seen[text] {
			return nil, fmt.Errorf("%s: trace.requirements.table.text-columns enthält %q doppelt", FileName, text)
		}
		seen[text] = true
		texts = append(texts, text)
	}
	return texts, nil
}

// applyModality validiert den opt-in Modalitäts-Block (DC-FA-MOD-001) und liefert
// ihn (nil ⇒ abwesend, aus). Präsenz — auch leere Map — ist aktiv. Fail-closed
// (Exit 2): leerer Stufen-Name, reservierter Name `unknown`, leeres Keyword,
// dasselbe Keyword in mehreren Stufen (Nondeterminismus), oder ein require-levels-
// Eintrag, der weder deklarierte (bzw. Default-)Stufe noch `unknown` ist.
func applyModality(raw *rawModality) (*model.TraceModality, error) {
	if raw == nil {
		return nil, nil
	}
	if err := validateModalityLevels(raw.Levels); err != nil {
		return nil, err
	}
	effective := raw.Levels
	if len(effective) == 0 {
		effective = model.DefaultModalityLevels()
	}
	valid := map[string]bool{"unknown": true}
	for name := range effective {
		valid[name] = true
	}
	for _, rl := range raw.RequireLevels {
		if !valid[rl] {
			return nil, fmt.Errorf("%s: trace.requirements.modality.require-levels: %q ist keine deklarierte Stufe (noch \"unknown\")", FileName, rl)
		}
	}
	return &model.TraceModality{Levels: raw.Levels, RequireLevels: raw.RequireLevels}, nil
}

// validateModalityLevels prüft eine explizit gesetzte levels-Map fail-closed:
// leerer Stufen-Name, reservierter Name `unknown`, leeres Keyword, dasselbe
// Keyword in mehreren Stufen (Nondeterminismus). Leere Map ⇒ Defaults, ok.
func validateModalityLevels(levels map[string][]string) error {
	seenKW := map[string]string{}
	for name, kws := range levels {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("%s: trace.requirements.modality.levels: leerer Stufen-Name", FileName)
		}
		if name == "unknown" {
			return fmt.Errorf("%s: trace.requirements.modality.levels: Stufen-Name \"unknown\" ist reserviert", FileName)
		}
		for _, kw := range kws {
			if strings.TrimSpace(kw) == "" {
				return fmt.Errorf("%s: trace.requirements.modality.levels[%q]: leeres Keyword", FileName, name)
			}
			key := strings.ToLower(strings.TrimSpace(kw))
			if other, dup := seenKW[key]; dup && other != name {
				return fmt.Errorf("%s: trace.requirements.modality: Keyword %q in mehreren Stufen (%q, %q) — nondeterministisch", FileName, kw, other, name)
			}
			seenKW[key] = name
		}
	}
	return nil
}

// applyTraceCoverage validiert die opt-in Coverage-Quellen (DC-FA-COV-001):
// nicht-leere files (jede innerhalb der Repo-Wurzel), nicht-leeres label,
// ranges-Default true. Fehlende Datei / Sektions-/Range-Fehler sind Lauf-zeitig
// (Kern), da config-zeitig kein I/O erfolgt.
func applyTraceCoverage(coverage []rawCoverage, tc *model.TraceConfig) error {
	for i, cov := range coverage {
		if len(cov.Files) == 0 {
			return fmt.Errorf("%s: trace.coverage[%d].files ist leer", FileName, i)
		}
		for _, f := range cov.Files {
			if strings.TrimSpace(f) == "" {
				return fmt.Errorf("%s: trace.coverage[%d].files enthält einen leeren Pfad", FileName, i)
			}
			if err := validateTracePath(fmt.Sprintf("trace.coverage[%d].files", i), f); err != nil {
				return err
			}
		}
		if strings.TrimSpace(cov.Label) == "" {
			return fmt.Errorf("%s: trace.coverage[%d].label ist leer", FileName, i)
		}
		ranges := true // Default true (DC-FA-COV-001); der Pointer unterscheidet nicht-gesetzt.
		if cov.Ranges != nil {
			ranges = *cov.Ranges
		}
		tc.Coverage = append(tc.Coverage, model.TraceCoverage{
			Files: cov.Files, Label: cov.Label, Ranges: ranges,
			Sections: cov.Sections, ExcludeSections: cov.ExcludeSections,
		})
	}
	return nil
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
	closure, err := applyClosure(p.Closure)
	if err != nil {
		return err
	}
	waves, err := applyWaves(p.Waves)
	if err != nil {
		return err
	}
	cfg.Planning = model.PlanningConfig{
		Roadmap: p.Roadmap, Heading: p.Heading, Marker: p.Marker, SliceGlob: p.SliceGlob,
		Closure: closure, Waves: waves,
	}
	return nil
}

// applyWaves validiert die Wellen-Invariante-Parameter (DC-FA-PLAN-001,
// Spez-Schritt W1). Dieselbe Config-Rand-Disziplin: ein Pfad, der die Wurzel
// verlässt, ein ungültiges Glob und eine leere Register-Überschrift sind
// jeweils ein stilles Grün, wenn man sie durchlässt.
func applyWaves(w *rawWaves) (model.WavesConfig, error) {
	if w == nil {
		return model.WavesConfig{}, nil
	}
	for name, dir := range map[string]string{"dir": w.Dir, "done-dir": w.DoneDir} {
		if strings.HasPrefix(dir, "/") || strings.Contains(dir, "..") {
			return model.WavesConfig{}, fmt.Errorf(
				"%s: planning.waves.%s %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')",
				FileName, name, dir)
		}
	}
	glob, err := wavesGlob(w.Glob, "glob")
	if err != nil {
		return model.WavesConfig{}, err
	}
	resultsGlob, err := wavesGlob(w.ResultsGlob, "results-glob")
	if err != nil {
		return model.WavesConfig{}, err
	}
	// An der Kennung hängt jede Zuordnung: beide Globs müssen mit demselben
	// nicht-leeren literalen Präfix beginnen — sonst laufen Datei-Zuordnung und
	// Zeilen-Erkennung auseinander (falsches wave-drift bei führendem
	// Platzhalter, unerreichbares wave-unregistered bei fremdem
	// results-Präfix).
	pfx := globPrefix(effectiveOr(glob, "welle-*.md"))
	if pfx == "" {
		return model.WavesConfig{}, fmt.Errorf(
			"%s: planning.waves.glob %q beginnt mit einem Platzhalter — die Wellen-Kennung braucht ein literales Präfix", FileName, *w.Glob)
	}
	if rp := globPrefix(effectiveOr(resultsGlob, "welle-*-results.md")); !strings.HasPrefix(rp, pfx) {
		return model.WavesConfig{}, fmt.Errorf(
			"%s: planning.waves.results-glob %q trägt nicht das Kennungs-Präfix %q von planning.waves.glob", FileName, rp, pfx)
	}
	next, err := wavesHeading(w.NextHeading, "next-heading")
	if err != nil {
		return model.WavesConfig{}, err
	}
	closed, err := wavesHeading(w.ClosedHeading, "closed-heading")
	if err != nil {
		return model.WavesConfig{}, err
	}
	mode, err := wavesMode(w.Mode)
	if err != nil {
		return model.WavesConfig{}, err
	}
	return model.WavesConfig{
		Dir: w.Dir, DoneDir: w.DoneDir, Glob: glob, ResultsGlob: resultsGlob,
		NextHeading: next, ClosedHeading: closed, Mode: mode,
	}, nil
}

// wavesMode prüft das Kardinalitäts-Modell der ersten Wellen-Aussage
// (Spez-Schritt W1): abwesend ⇒ "" (Kern-Default one); jeder andere Wert als
// one/many — der explizit leere String eingeschlossen — bricht mit Exit 2 ab,
// dieselbe Zeiger-Disziplin wie bei den übrigen waves-Schlüsseln.
func wavesMode(m *string) (string, error) {
	if m == nil {
		return "", nil
	}
	if *m != "one" && *m != "many" {
		return "", fmt.Errorf(
			"%s: planning.waves.mode %q ist ungültig — erlaubt sind \"one\" (Default) und \"many\" (weglassen ⇒ one)",
			FileName, *m)
	}
	return *m, nil
}

// wavesGlob prüft einen explizit gesetzten Glob: leer ist eine Null-Aussage
// (Exit 2), ungültig bricht ab; abwesend liefert "" (Kern-Default).
func wavesGlob(g *string, name string) (string, error) {
	if g == nil {
		return "", nil
	}
	if *g == "" {
		return "", fmt.Errorf(
			"%s: planning.waves.%s ist leer — kein Basisname könnte matchen (weglassen ⇒ Default)", FileName, name)
	}
	if _, err := path.Match(*g, "probe"); err != nil {
		return "", fmt.Errorf("%s: planning.waves.%s %q ist kein gültiges Glob: %v", FileName, name, *g, err)
	}
	return *g, nil
}

// wavesHeading prüft eine explizit gesetzte Register-Überschrift: leer oder
// nur Whitespace ist eine Null-Aussage (Exit 2); abwesend liefert "".
func wavesHeading(h *string, name string) (string, error) {
	if h == nil {
		return "", nil
	}
	if strings.TrimSpace(*h) == "" {
		return "", fmt.Errorf(
			"%s: planning.waves.%s ist leer — keine Zeile könnte die Register-Überschrift sein (weglassen ⇒ Default)", FileName, name)
	}
	return *h, nil
}

// globPrefix liefert den literalen Teil eines Globs vor dem ersten
// Platzhalter.
func globPrefix(glob string) string {
	if i := strings.IndexAny(glob, "*?["); i >= 0 {
		return glob[:i]
	}
	return glob
}

// effectiveOr liefert den gesetzten Wert oder den Default.
func effectiveOr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// applyClosure validiert die Closure-Note-Struktur-Parameter (DC-FA-PLAN-001).
// Dieselbe Config-Rand-Disziplin wie bei planning.slice-glob: was der Kern nur
// schlucken könnte, bricht hier laut ab (Exit 2) — ein nicht kompilierendes
// Muster, eine Null-/Negativ-Schwelle und ein leerer Floskel-Eintrag sind
// jeweils ein stilles Grün, wenn man sie durchlässt.
func applyClosure(c *rawClosure) (model.ClosureConfig, error) {
	if c == nil {
		return model.ClosureConfig{}, nil
	}
	if strings.HasPrefix(c.Dir, "/") || strings.Contains(c.Dir, "..") {
		return model.ClosureConfig{}, fmt.Errorf(
			"%s: planning.closure.dir %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, c.Dir)
	}
	if c.HeadingPattern != "" {
		if _, err := regexp.Compile(c.HeadingPattern); err != nil {
			return model.ClosureConfig{}, fmt.Errorf(
				"%s: planning.closure.heading-pattern %q ist kein gültiges Regex: %v", FileName, c.HeadingPattern, err)
		}
	}
	// Ein EXPLIZIT gesetzter Glob muss etwas treffen können; abwesend ⇒ Verweis
	// auf slice-glob im Kern. Leer oder ungültig bräche die Aussage, die das
	// Setzen des Schlüssels darstellt (§DC-FA-PLAN-001.a Schritt C2).
	glob := ""
	if c.Glob != nil {
		if *c.Glob == "" {
			return model.ClosureConfig{}, fmt.Errorf(
				"%s: planning.closure.glob ist leer — kein Kandidat könnte matchen (weglassen ⇒ planning.slice-glob)", FileName)
		}
		if _, err := path.Match(*c.Glob, "probe"); err != nil {
			return model.ClosureConfig{}, fmt.Errorf(
				"%s: planning.closure.glob %q ist kein gültiges Glob: %v", FileName, *c.Glob, err)
		}
		glob = *c.Glob
	}
	// Nur ein EXPLIZIT gesetzter Wert wird geprüft; abwesend ⇒ Kern-Default (4).
	// 0 wäre „jede Notiz genügt", negativ ist sinnlos — beides ist ein stilles
	// Grün am Config-Rand und bricht daher laut ab.
	minSentences := 0
	if c.MinSentences != nil {
		if *c.MinSentences < 1 {
			return model.ClosureConfig{}, fmt.Errorf(
				"%s: planning.closure.min-sentences %d muss >= 1 sein", FileName, *c.MinSentences)
		}
		minSentences = *c.MinSentences
	}
	for _, b := range c.Boilerplate {
		if strings.TrimSpace(b) == "" {
			return model.ClosureConfig{}, fmt.Errorf(
				"%s: planning.closure.boilerplate enthält einen leeren Eintrag (träfe jeden Text)", FileName)
		}
	}
	return model.ClosureConfig{
		Dir: c.Dir, Glob: glob, HeadingPattern: c.HeadingPattern,
		MinSentences: minSentences, Boilerplate: c.Boilerplate,
		Placeholder: c.Placeholder,
	}, nil
}

// applyTracked validiert die Parameter des Moduls tracked (DC-FA-TRK-001):
// jedes exempt-targets-Glob muss ein gültiges path.Match-Muster sein — sonst
// schluckte der Kern den ErrBadPattern und das Ventil wäre still wirkungslos
// bzw. die Prüfung fail-open (dieselbe Config-Rand-Disziplin wie
// planning.slice-glob, seit slice-057).
func applyTracked(r *raw, cfg *model.Config) error {
	if r.Tracked == nil {
		return nil
	}
	for _, g := range r.Tracked.ExemptTargets {
		if g == "" {
			return fmt.Errorf("%s: tracked.exempt-targets enthält ein leeres Glob", FileName)
		}
		// Segmentweise validieren — die Laufzeit (matchGlob) matcht je
		// '/'-Segment; ein nur als Ganzes gültiges Muster wäre still
		// wirkungslos (seit slice-059).
		for _, seg := range strings.Split(g, "/") {
			if seg == "**" {
				continue
			}
			if _, err := path.Match(seg, "probe"); err != nil {
				return fmt.Errorf("%s: tracked.exempt-targets %q ist kein gültiges Glob (Segment %q): %v", FileName, g, seg, err)
			}
		}
	}
	cfg.Tracked = model.TrackedConfig{ExemptTargets: r.Tracked.ExemptTargets}
	return nil
}

// applyTargets validiert die Parameter des Moduls targets (DC-FA-TGT-001): jeder
// konfigurierte Pfad (makefiles/doc-tables/authority) muss relativ zur
// Repo-Wurzel liegen (die Existenz prüft der Kern beim Lauf, fail-closed →
// Exit 2). Leere makefiles ⇒ Modul inert.
func applyTargets(r *raw, cfg *model.Config) error {
	if r.Targets == nil {
		return nil
	}
	t := r.Targets
	paths := append(append(append([]string{}, t.Makefiles...), t.DocTables...), t.Authority)
	for _, p := range paths {
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "/") || strings.Contains(p, "..") {
			return fmt.Errorf("%s: targets-Pfad %q muss relativ zur Repo-Wurzel liegen (kein '/', kein '..')", FileName, p)
		}
	}
	cfg.Targets = model.TargetsConfig{
		Makefiles: t.Makefiles, DocTables: t.DocTables,
		Authority: t.Authority, ExemptTargets: t.ExemptTargets,
	}
	return nil
}

// applySources validiert die Config-Pins des Moduls sources (DC-FA-SRC-001):
// jeder Eintrag braucht eine absolute http(s)-url und genau 64 Hex-`sha256`;
// `unpack` ist `none` (Default) oder `zip`. Jede Verletzung ⇒ Exit 2 (fail-closed).
// `sha256` wird case-insensitiv geführt (zu Kleinbuchstaben normalisiert). Die
// Befund-Zeile ist die yaml-Node-Zeile des url-Feldes (Fallback 1).
func applySources(r *raw, cfg *model.Config) error {
	if len(r.Sources) == 0 {
		return nil
	}
	pins := make([]model.SourcePin, 0, len(r.Sources))
	for i, s := range r.Sources {
		pin, err := validateSource(i, s)
		if err != nil {
			return err
		}
		pins = append(pins, pin)
	}
	cfg.Sources = model.SourcesConfig{Pins: pins}
	return nil
}

// validateSource prüft einen einzelnen sources[]-Eintrag (fail-closed).
func validateSource(i int, s rawSource) (model.SourcePin, error) {
	url := strings.TrimSpace(s.URL.Value)
	switch {
	case url == "":
		return model.SourcePin{}, fmt.Errorf("%s: sources[%d].url fehlt", FileName, i)
	case !isHTTPURL(url):
		return model.SourcePin{}, fmt.Errorf("%s: sources[%d].url %q ist keine absolute http(s)-URL", FileName, i, url)
	}
	sum := strings.TrimSpace(s.Sha256)
	if !is64Hex(sum) {
		return model.SourcePin{}, fmt.Errorf("%s: sources[%d].sha256 muss genau 64 Hex-Zeichen sein", FileName, i)
	}
	unpack := s.Unpack
	if unpack == "" {
		unpack = "none"
	}
	if unpack != "none" && unpack != model.SourceUnpackZip {
		return model.SourcePin{}, fmt.Errorf("%s: sources[%d].unpack %q ungültig (nur none oder zip)", FileName, i, unpack)
	}
	line := s.URL.Line
	if line <= 0 {
		line = 1
	}
	return model.SourcePin{URL: url, Sha256: strings.ToLower(sum), Unpack: unpack, Line: line}, nil
}

// isHTTPURL prüft case-insensitiv auf ein absolutes http(s)-Schema (wie das
// Modul external — sources holt nur http/https, spec/spezifikation.md §2).
func isHTTPURL(u string) bool {
	return (len(u) >= 7 && strings.EqualFold(u[:7], "http://")) ||
		(len(u) >= 8 && strings.EqualFold(u[:8], "https://"))
}

// is64Hex prüft, ob s genau 64 Hex-Zeichen (Groß/Klein) ist (DC-FA-SRC-001
// sha256). encoding/hex akzeptiert beide Schreibweisen und gerade Längen.
func is64Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
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
		CheckLines:  r.Codepaths.CheckLines,
	}
	return nil
}

// applyIgnoreRefs validiert das geteilte Referenz-Ventil ignore-refs
// (DC-FA-REF-001): jedes in/refs/keep-Glob muss segmentweise ein gültiges
// path.Match-Muster sein — dieselbe Rand-Disziplin wie tracked.exempt-targets,
// sonst schluckte der Kern den ErrBadPattern und das Ventil wäre still
// wirkungslos. Leeres in = repo-weit (erlaubt); leere refs/keep-Elemente sind
// ein Fehler. Ein leeres refs (Liste) macht den Eintrag inert.
func applyIgnoreRefs(r *raw, cfg *model.Config) error {
	if len(r.IgnoreRefs) == 0 {
		return nil
	}
	out := make([]model.IgnoreRef, 0, len(r.IgnoreRefs))
	for _, e := range r.IgnoreRefs {
		if err := validIgnoreRefEntry(e); err != nil {
			return err
		}
		out = append(out, model.IgnoreRef{In: e.In, Refs: e.Refs, Keep: e.Keep})
	}
	cfg.IgnoreRefs = out
	return nil
}

// validIgnoreRefEntry validiert einen ignore-refs-Eintrag: leeres in =
// repo-weit (erlaubt), leere refs/keep-Elemente sind ein Fehler; jedes
// nichtleere Glob wird segmentweise geprüft.
func validIgnoreRefEntry(e rawIgnoreRef) error {
	if e.In != "" {
		if err := validRefGlob(e.In); err != nil {
			return err
		}
	}
	for _, list := range [][]string{e.Refs, e.Keep} {
		for _, g := range list {
			if g == "" {
				return fmt.Errorf("%s: ignore-refs enthält ein leeres refs/keep-Glob", FileName)
			}
			if err := validRefGlob(g); err != nil {
				return err
			}
		}
	}
	return nil
}

// validRefGlob prüft segmentweise (wie die Laufzeit matchGlob), dass ein
// ignore-refs-Glob ein gültiges path.Match-Muster ist — sonst schluckte der
// Kern den ErrBadPattern und das Ventil wäre still wirkungslos.
func validRefGlob(g string) error {
	for _, seg := range strings.Split(g, "/") {
		if seg == "**" {
			continue
		}
		if _, err := path.Match(seg, "probe"); err != nil {
			return fmt.Errorf("%s: ignore-refs-Glob %q ist ungültig (Segment %q): %v", FileName, g, seg, err)
		}
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

// validateSegmentGlobs prüft eine Glob-Liste SEGMENTWEISE — so, wie die
// Laufzeit (matchGlob) sie matcht: ein nur als Ganzes gültiges Muster waere
// still wirkungslos. Ein leerer Eintrag ist ein Fehler, `**` ist zulässig.
func validateSegmentGlobs(name string, globs []string) error {
	for _, g := range globs {
		if g == "" {
			return fmt.Errorf("%s: %s enthält ein leeres Glob", FileName, name)
		}
		for _, seg := range strings.Split(g, "/") {
			if seg == "**" {
				continue
			}
			if _, err := path.Match(seg, "probe"); err != nil {
				return fmt.Errorf("%s: %s %q ist kein gültiges Glob (Segment %q): %v", FileName, name, g, seg, err)
			}
		}
	}
	return nil
}

// applyDiagrams validiert und kompiliert die diagrams-Muster
// (DC-FA-DIAG-001): nicht-leere fences-Einträge, kompilierbarer Regex
// (nicht den Leerstring matchend), Pflicht-defined-in und gültige
// exempt-paths-Globs.
func applyDiagrams(r *raw, cfg *model.Config) error {
	if r.Diagrams == nil {
		return nil
	}
	for _, f := range r.Diagrams.Fences {
		if strings.TrimSpace(f) == "" {
			return fmt.Errorf("%s: diagrams.fences enthält einen leeren Eintrag", FileName)
		}
	}
	if err := validateSegmentGlobs("diagrams.exempt-paths", r.Diagrams.ExemptPaths); err != nil {
		return err
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
	cfg.Diagrams.ExemptPaths = r.Diagrams.ExemptPaths
	return nil
}

// applyVersions validiert und kompiliert die versions-Muster-Quellen-Paare
// (DC-FA-VER-001). Die Kurzform IST die einelementige patterns-Liste und wird
// hier in sie übersetzt — der Kern kennt nur die Liste, es gibt genau einen
// Auswertungspfad. Beide Schreibweisen zugleich, eine explizit leere Liste
// oder ein Paar ohne pin-pattern sind Nutzungsfehler (Exit 2);
// current-from/exempt-paths werden im Kern beim Lauf-Start bzw. je Datei
// ausgewertet.
func applyVersions(r *raw, cfg *model.Config) error {
	if r.Versions == nil {
		return nil
	}
	v := r.Versions
	if v.Patterns != nil {
		if hasShortVersionForm(v) {
			return fmt.Errorf("%s: versions: Kurzform (pin-pattern/current-from/exempt-paths) und versions.patterns zugleich gesetzt — "+
				"welche Ausnahmen für welches Muster gälten, wäre nicht ablesbar", FileName)
		}
		if len(v.Patterns) == 0 {
			return fmt.Errorf("%s: versions.patterns ist leer", FileName)
		}
		patterns := make([]model.VersionPattern, 0, len(v.Patterns))
		for i, rp := range v.Patterns {
			p, err := compileVersionPattern(fmt.Sprintf("versions.patterns[%d]", i), rp)
			if err != nil {
				return err
			}
			patterns = append(patterns, p)
		}
		cfg.Versions = model.VersionsConfig{Patterns: patterns}
		return nil
	}
	p, err := compileVersionPattern("versions", shortVersionPattern(v))
	if err != nil {
		return err
	}
	cfg.Versions = model.VersionsConfig{Patterns: []model.VersionPattern{p}}
	return nil
}

// hasShortVersionForm meldet, ob einer der drei Kurzform-Schlüssel ANWESEND
// ist (DC-FA-VER-001) — zusammen mit patterns ein Nutzungsfehler. Grenze: ein
// Schlüssel ohne Wert (`pin-pattern:` ohne Muster) ist im YAML von einem
// fehlenden Schlüssel nicht unterscheidbar und zählt als fehlend.
func hasShortVersionForm(v *rawVersions) bool {
	return v.PinPattern != nil || v.CurrentFrom != nil || v.ExemptPaths != nil
}

// shortVersionPattern liest die drei Kurzform-Schlüssel als ein Paar; ein
// fehlender Schlüssel wird zum Leerwert, den compileVersionPattern bzw. der
// Kern-Default behandelt.
func shortVersionPattern(v *rawVersions) rawVersionPattern {
	rp := rawVersionPattern{}
	if v.PinPattern != nil {
		rp.PinPattern = *v.PinPattern
	}
	if v.CurrentFrom != nil {
		rp.CurrentFrom = *v.CurrentFrom
	}
	if v.ExemptPaths != nil {
		rp.ExemptPaths = *v.ExemptPaths
	}
	return rp
}

// compileVersionPattern validiert ein Muster-Quellen-Paar: Pflicht-pin-pattern,
// kompilierbar und nicht den Leerstring matchend (DC-FA-VER-001). prefix
// benennt die Fundstelle in der Konfiguration, damit die Meldung bei einer
// Liste auf das Paar zeigt.
func compileVersionPattern(prefix string, rp rawVersionPattern) (model.VersionPattern, error) {
	if strings.TrimSpace(rp.PinPattern) == "" {
		return model.VersionPattern{}, fmt.Errorf("%s: %s.pin-pattern fehlt", FileName, prefix)
	}
	re, err := regexp.Compile(rp.PinPattern)
	if err != nil {
		return model.VersionPattern{}, fmt.Errorf("%s: %s.pin-pattern: %v", FileName, prefix, err)
	}
	if re.MatchString("") {
		return model.VersionPattern{}, fmt.Errorf("%s: %s.pin-pattern matcht den Leerstring", FileName, prefix)
	}
	return model.VersionPattern{PinPattern: re, CurrentFrom: rp.CurrentFrom, ExemptPaths: rp.ExemptPaths}, nil
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
		{"links", scopeOfLinks(r.Links)},
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
		{"tracked", scopeOfTracked(r.Tracked)},
		{"citations", scopeOf(r.Citations)},
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
func scopeOfLinks(v *rawLinks) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

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

func scopeOfTracked(v *rawTracked) *rawScope {
	if v == nil {
		return nil
	}
	return v.Scope
}

// decodeStrict dekodiert mit KnownFields; nil-raw bei leerem Dokument.
func decodeStrict(content []byte) (*raw, error) {
	var r raw
	dec := yaml.NewDecoder(bytes.NewReader(content))
	dec.KnownFields(true)
	if err := dec.Decode(&r); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		// interne Go-Typnamen aus der yaml-Meldung entfernen
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
