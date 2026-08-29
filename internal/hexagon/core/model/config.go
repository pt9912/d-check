package model

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ResolveFromGroup ist eine Gruppe ortsfester Verweise (DC-FA-LINK-001
// §Ortsfeste Verweise, Schritt 6): Dateien in den WANDERNDEN Verzeichnissen
// (Dirs) sind Quellen, deren relative Ziele von jedem Ort der Gruppe
// (Dirs ∪ FixedDirs) auflösen müssen — und überall auf dasselbe Ziel.
// Dateien in FixedDirs sind am Endzustand und keine Quellen.
type ResolveFromGroup struct {
	Dirs      []string
	FixedDirs []string
}

// validModules sind die vertraglich gültigen Regelmodul-Namen
// (DC-FA-CLI-002).
func validModules() []string {
	return []string{"links", "anchors", "ids", "matrix", "external", "codepaths", "spans", "hostpaths", "diagrams", "versions", "pins", "immutable", "vcs", "commits", "planning", "tracked", "targets", "citations", "sources", "structure"}
}

// ValidModules ist die exportierte Sicht auf validModules (DC-FA-CLI-002) —
// genutzt vom --print-mk-Generator (DC-FA-CLI-010), damit das doc-immutable-
// Target seine vcs-Fokus-Disable-Liste aus dem Modulsatz ableitet statt sie zu
// duplizieren. Liefert bei jedem Aufruf eine frische Liste in Deklarations-
// reihenfolge (deterministisch, DC-QA-02).
func ValidModules() []string { return validModules() }

// defaultModules ist der Default-Modulsatz (DC-FA-CLI-002).
func defaultModules() []string { return []string{"links", "anchors"} }

// Config ist die validierte Konfiguration (Teilmenge, die der
// aktuelle Stand auswertet; das Voll-Schema inkl. ids/matrix/external
// wird vom Config-Adapter bereits strikt geparst und statisch
// validiert — spec/spezifikation.md §2).
type Config struct {
	// Roots: nil = Default-Wurzeln; sonst explizit (müssen existieren).
	Roots  []string
	Ignore []string
	// IgnoreRefs: geteiltes Referenz-Ventil (DC-FA-REF-001), honoriert
	// von links/anchors/codepaths; nil = kein Eintrag (byte-identisch).
	IgnoreRefs []IgnoreRef
	// ResolveFrom: Gruppen ortsfester Verweise (DC-FA-LINK-001 Schritt 6);
	// nil = kein Eintrag (byte-identisch).
	ResolveFrom []ResolveFromGroup
	// Modules: nil = DefaultModules.
	Modules []string
	// IDPatterns: bereits kompilierte Kennungs-Muster (Modul ids).
	IDPatterns []IDPattern
	// Matrix: Dokumentklassen und Referenzregeln (Modul matrix).
	Matrix MatrixConfig
	// External: Parameter des Moduls external.
	External ExternalConfig
	// Codepaths: Präfixe des Moduls codepaths.
	Codepaths CodepathsConfig
	// Citations: das Modul citations ist direktiven-getrieben und
	// parameterlos (DC-FA-CITE-001) — die Struktur trägt nur den
	// Scope-Platzhalter, damit Enable/Scope wie bei jedem Modul greifen.
	Citations CitationsConfig
	// Hostpaths: Parameter des Moduls hostpaths.
	Hostpaths HostpathsConfig
	// Diagrams: Parameter des Moduls diagrams (DC-FA-DIAG-001).
	Diagrams DiagramsConfig
	// Versions: Parameter des Moduls versions (DC-FA-VER-001).
	Versions VersionsConfig
	// Immutable: Parameter des Moduls immutable (DC-FA-IMM-001).
	Immutable ImmutableConfig
	// VCS: Parameter des Moduls vcs (DC-FA-VCS-001).
	VCS VCSConfig
	// Commits: Parameter des Moduls commits (DC-FA-COMMITS-001).
	Commits CommitsConfig
	// Planning: Parameter des Moduls planning (DC-FA-PLAN-001).
	Planning PlanningConfig
	Structure []StructureRule
	// Tracked: Parameter des Moduls tracked (DC-FA-TRK-001).
	Tracked TrackedConfig
	// Targets: Parameter des Moduls targets (DC-FA-TGT-001).
	Targets TargetsConfig
	// ConfigFile: der Wurzel-relative Pfad der TATSÄCHLICH geladenen
	// Konfigurationsdatei — konventionell `.d-check.yml`, mit `--config`
	// (DC-FA-CLI-012) die dort genannte. Befunde, die die Konfiguration als
	// Herkunft nennen, müssen diesen Pfad tragen; sonst zeigte die Meldung auf
	// eine Datei, die der Lauf nie gelesen hat. Leer ⇒ es wurde keine
	// Konfiguration geladen (reiner Default-Lauf).
	ConfigFile string
	// Sources: die Config-Pins des Moduls sources (DC-FA-SRC-001); leer ⇒
	// nur Marker-Pins (bzw. gar keine) werden geprüft (byte-identisch).
	Sources SourcesConfig
	// Trace: konfigurierbare Quellen der RTM (DC-FA-CLI-009); Nullwert =
	// Konventions-Default (byte-identisch).
	Trace TraceConfig
	// Scopes: modul-lokale Scan-Scopes (DC-FA-CONF-002); Schlüssel
	// ist der Modulname, nil-Eintrag/fehlender Schlüssel = globaler
	// Scope.
	Scopes map[string]*ScopeConfig
}

// IgnoreRef ist ein Eintrag des geteilten Referenz-Ventils `ignore-refs`
// (DC-FA-REF-001; Ziel-Achsen-Pendant zu scan.ignore). Ein Ziel wird
// ignoriert, wenn Refs matcht UND Keep nicht (Keep reihenfolge-unabhängig).
type IgnoreRef struct {
	// In: optionaler Glob (Syntax wie scan.ignore) auf die Quelldatei
	// (die Datei, in der die Referenz steht); "" = repo-weit.
	In string
	// Refs: Globs auf den aufgelösten Ziel-Pfad; leer = Eintrag inert.
	Refs []string
	// Keep: Globs; Ausnahmen — ein von Refs getroffenes Ziel bleibt
	// geprüft, wenn ein Keep-Glob es trifft (Keep gewinnt unbedingt).
	Keep []string
}

// ScopeConfig ersetzt für genau ein Modul den globalen Scan-Scope
// (DC-FA-CONF-002; Roots ist im Adapter als Pflichtfeld validiert —
// nie nil, eine leere Liste prüft nichts).
type ScopeConfig struct {
	Roots  []string
	Ignore []string
}

// IDPattern ist ein deklariertes Kennungs-Muster (DC-FA-ID-001).
type IDPattern struct {
	Regex  *regexp.Regexp
	Target string
	// LinkPolicy: "" = prose (Default, abwärtskompatibel); "always"
	// macht auch Inline-Code-Vorkommen linkpflichtig
	// (spec/spezifikation.md §DC-FA-ID-001.a).
	LinkPolicy string
	// ExemptPaths: Globs (Syntax wie scan.ignore) für Dateien ohne
	// Linkpflicht des Musters — gilt für nackte Fließtext- wie für
	// Inline-Code-Vorkommen, unabhängig von der LinkPolicy
	// (spec/spezifikation.md §DC-FA-ID-001.a).
	ExemptPaths []string
}

// AlwaysPolicy ist der link-policy-Wert, der auch Inline-Code-Vorkommen
// linkpflichtig macht (DC-FA-ID-001.a).
const AlwaysPolicy = "always"

// MatrixClass ist eine über Pfad-Globs deklarierte Dokumentklasse
// (DC-FA-MTX-001; Deklarationsreihenfolge = Präzedenz).
type MatrixClass struct {
	Name  string
	Paths []string
	// Order: optionale Rang-Globs der klasseninternen Schichtung
	// (DC-FA-MTX-002), autoritativste Schicht zuerst. Der Rang einer Datei
	// ist der Index des ersten matchenden Globs (First-Match wie die
	// Klassenzuordnung); ohne Treffer ist die Datei rangfrei. Nur zusammen
	// mit Direction wirksam (fail-closed im Config-Adapter validiert).
	Order []string
	// Direction: Richtungspolitik innerhalb der Klasse (DC-FA-MTX-002).
	// "" = aus (byte-identisch zu DC-FA-MTX-001); DirectionNoDownward
	// verbietet klasseninterne Verweise von höher- auf niederrangig.
	Direction string
	// Token: optionales Regex, das Referenzen auf diese Klasse als bare
	// ID-Token im Fließtext erkennt (DC-FA-MTX-003). nil = nur Link-Erkennung.
	Token *regexp.Regexp
}

// DirectionNoDownward ist die einzige unterstützte Richtungspolitik
// (DC-FA-MTX-002): klasseninterne Referenzen nur aufwärts zur
// autoritativeren Schicht; ein Abwärtsverweis erzeugt matrix-downward.
const DirectionNoDownward = "no-downward"

// MatrixRule deklariert, ob Referenzen von From nach To erlaubt sind.
type MatrixRule struct {
	From, To string
	Allow    bool
}

// MatrixConfig ist die validierte matrix-Konfiguration
// (DC-FA-MTX-001; Defaults gemäß spec/spezifikation.md §2).
type MatrixConfig struct {
	Classes         []MatrixClass
	Rules           []MatrixRule
	StatusForbidden []string
	// AllowSupersedeLineage (Default false) nimmt die deklarierte
	// Supersede-Lineage-Kante von der Status-Prüfung aus; ohne true
	// wird SupersedeFields nie konsultiert (byte-identisch, DC-QA-02).
	AllowSupersedeLineage bool
	// SupersedeFields: Feldnamen (z. B. "Supersedes", "Aenderungstyp"),
	// aus deren Werten die Ablösung gelesen wird (nur wirksam bei
	// AllowSupersedeLineage; spec/spezifikation.md §DC-FA-MTX-001.a).
	SupersedeFields []string
	ExcludeSections []string
	// ExemptPaths: Globs (Syntax wie scan.ignore); Dateien ganz ohne
	// matrix-Prüfung — Grandfathering immutabler Dokumente (DC-FA-MTX-003).
	ExemptPaths []string
}

// ExternalConfig sind die Parameter des Moduls external
// (DC-FA-EXT-001); 0 = Default gemäß spec/spezifikation.md §3.
type ExternalConfig struct {
	TimeoutSeconds int
	Parallel       int
}

// EffectiveTimeoutSeconds liefert das Timeout (EXTERNAL_TIMEOUT = 10 s).
func (e ExternalConfig) EffectiveTimeoutSeconds() int {
	if e.TimeoutSeconds == 0 {
		return 10
	}
	return e.TimeoutSeconds
}

// EffectiveParallel liefert die Parallelität (EXTERNAL_PARALLEL = 4).
func (e ExternalConfig) EffectiveParallel() int {
	if e.Parallel == 0 {
		return 4
	}
	return e.Parallel
}

// CodepathsConfig trägt die Präfixe für Wurzel-relative
// Inline-Code-Pfade (spec/spezifikation.md §2; `./`/`../` werden
// immer erkannt).
type CodepathsConfig struct {
	Roots []string
	// ExemptPaths: Globs (Syntax wie scan.ignore) für Dateien, die das
	// Modul codepaths nicht prüft — datei-weit, unabhängig von Roots
	// (spec/spezifikation.md §DC-FA-CODE-001.a). Vorbild: das
	// gleichnamige Ventil von IDPattern (DC-FA-ID-001).
	ExemptPaths []string
	// IgnoreRefs: Globs (Syntax wie scan.ignore) für aufgelöste
	// Ziel-Pfade, die nicht existenz-/anker-geprüft werden — referenz-weit
	// (datei-/zeilen-unabhängig), Tombstone-Register bewusst entfernter
	// Artefakte (spec/spezifikation.md §DC-FA-CODE-001.a; ADR-0025).
	IgnoreRefs []string
	// CheckLines: Zeilen-Referenz (`datei:<von>-<bis>`) eines
	// Inline-Code-Pfads verifizieren statt sie zu verwerfen
	// (opt-in; Default aus ⇒ byte-identisch, DC-FA-CODE-001.a Schritt 6).
	CheckLines bool
}

// CitationsConfig sind die Parameter des Moduls citations
// (DC-FA-CITE-001) — direktiven-getrieben und parameterlos; die Struktur
// existiert, damit das Modul über die gemeinsame Config-/Scope-Maschinerie
// enable-/skopierbar ist (Vorbild: das parameterlose Modul spans).
type CitationsConfig struct{}

// HostpathsConfig sind die Parameter des Moduls hostpaths
// (DC-FA-HOST-001); Prefixes nil = Default-Liste.
type HostpathsConfig struct {
	Prefixes []string
}

// DiagramPattern ist ein Kennungs-Muster des Moduls diagrams
// (DC-FA-DIAG-001): Regex erkennt das Token in der Diagramm-Fence,
// DefinedIn ist die Datei, in der das Token als eigenständiges
// Kennungs-Token außerhalb von Fences vorkommen muss (Existenz, nicht
// Linkpflicht).
type DiagramPattern struct {
	Regex     *regexp.Regexp
	DefinedIn string
}

// DiagramsConfig sind die Parameter des Moduls diagrams (DC-FA-DIAG-001);
// Fences nil = Default ["mermaid"]. Ohne Patterns ist das Modul
// wirkungslos (byte-identisch zum Lauf ohne das Modul, DC-QA-02).
type DiagramsConfig struct {
	Fences   []string
	Patterns []DiagramPattern
	// ExemptPaths nimmt ganze Dateien aus der Pruefung (datei-weit, Glob wie
	// scan.ignore) — Ventil-Paritaet zu ids/codepaths; der Zeilen-Marker ist
	// das zweite Ventil und lebt im Modul.
	ExemptPaths []string
}

// EffectiveFences liefert die zu öffnenden Fence-Sprachen (Default
// "mermaid", DC-FA-DIAG-001.a).
func (d DiagramsConfig) EffectiveFences() []string {
	if len(d.Fences) == 0 {
		return []string{"mermaid"}
	}
	return d.Fences
}

// VersionPattern ist ein Paar aus Pin-Muster, Versions-Quelle und eigenen
// Ausnahmen (DC-FA-VER-001): PinPattern erkennt einen Versions-Pin (Version in
// Capture-Gruppe 1, sonst der ganze Treffer), CurrentFrom adressiert den Span
// mit der aktuellen Version (Default `version.md#aktuell`), ExemptPaths nimmt
// ganze Dateien aus — nur für dieses Paar.
type VersionPattern struct {
	PinPattern  *regexp.Regexp
	CurrentFrom string
	ExemptPaths []string
}

// EffectiveCurrentFrom liefert die Quelle der aktuellen Version (Default
// `version.md#aktuell`, spec/spezifikation.md §DC-FA-VER-001.a).
func (v VersionPattern) EffectiveCurrentFrom() string {
	if v.CurrentFrom == "" {
		return "version.md#aktuell"
	}
	return v.CurrentFrom
}

// VersionsConfig sind die Parameter des Moduls versions (DC-FA-VER-001): eine
// Liste von Muster-Quellen-Paaren, die je für sich ausgewertet werden. Die
// Kurzform der Konfiguration (die drei Schlüssel direkt unter `versions`) IST
// die einelementige Liste und wird am Config-Rand in sie übersetzt — der Kern
// kennt nur die Liste, es gibt genau einen Auswertungspfad. Ohne Paar ist das
// Modul wirkungslos (byte-identisch zum Lauf ohne das Modul, DC-QA-02).
type VersionsConfig struct {
	Patterns []VersionPattern
}

// ImmutableConfig sind die Parameter des Moduls immutable (DC-FA-IMM-001):
// ExcludeSections nennt Heading-Titel, deren Abschnitte nicht zum gehashten
// Core zählen (für ADRs typisch ["Geschichte"]; Vergleich wie
// MatrixConfig.ExcludeSections). Ohne Pin-Marker in einer Datei ist das Modul
// für diese Datei wirkungslos (opt-in pro Datei) — die Konfiguration ist
// daher optional.
type ImmutableConfig struct {
	ExcludeSections []string
}

// VCSConfig sind die Parameter des Moduls vcs (DC-FA-VCS-001): Paths ist die
// Glob-Klasse der zu schützenden Dateien (leer ⇒ Modul inert); ImmutableWhen
// entscheidet, ob die BASE-Version immutabel ist (irgendeine Zeile matcht);
// ExcludeSections nennt Heading-Titel, deren Abschnitte nicht zum Core zählen
// (wie ImmutableConfig.ExcludeSections); StatusLine erkennt die Kopf-Status-Zeile
// (erstes Vorkommen vor der ersten H2 — wird aus dem Core entfernt und für die
// HeadAllow-Prüfung gelesen); HeadAllow ist das erlaubte Muster der HEAD-Status-
// Zeile (Übergangs-Prüfung). Die git-historienbasierte Hälfte der Immutabilität —
// die hermetische Pin-Hälfte ist ImmutableConfig.
type VCSConfig struct {
	Paths           []string
	ImmutableWhen   *regexp.Regexp
	ExcludeSections []string
	StatusLine      *regexp.Regexp
	HeadAllow       *regexp.Regexp
}

// CommitsConfig sind die Parameter des Moduls commits (DC-FA-COMMITS-001):
// IDPatterns sind die kompilierten Muster gültiger Traceability-Kennungen — eine
// bereinigte Commit-Message ohne Match auf **irgendein** Muster erzeugt
// `commit-untraceable` (leere Liste ⇒ Modul inert). ExemptPattern nimmt eine
// Message aus, deren Betreff (erste Zeile) matcht (Selbstkonfig `^(Merge |Revert )`);
// nil ⇒ keine Ausnahme. Die Portierung des abgelösten `tools/trace-check.sh`.
type CommitsConfig struct {
	IDPatterns    []*regexp.Regexp
	ExemptPattern *regexp.Regexp
}

// PlanningConfig sind die Parameter des Moduls planning (DC-FA-PLAN-001):
// Roadmap ist die Roadmap-Datei mit dem Aktiv-Status-Abschnitt (leer ⇒ Modul
// inert; ihr Verzeichnis ist das Slice-Verzeichnis); Heading die kanonische
// H2-Überschrift, Marker der literale Ruhe-Marker, SliceGlob das Basisnamen-Glob
// der Slice-Dateien. Heading/Marker/SliceGlob leer ⇒ Konventions-Default
// (`## Aktuelle Welle` / „Keine aktive Welle" / `slice-*.md`). Hermetisch — nur
// der Filesystem-Port, kein git.
type PlanningConfig struct {
	Roadmap   string
	Heading   string
	Marker    string
	SliceGlob string
	Closure   ClosureConfig
	Waves     WavesConfig
}

// WavesConfig ist die dritte planning-Fähigkeit (DC-FA-PLAN-001
// §Wellen-Invariante, Spez-Schritte W1–W5): dieselbe Lifecycle-Invariante eine
// Ebene höher — die Wellen-Abschnitte der Roadmap gegen die Wellen-Dateien.
// Opt-in INNERHALB des opt-in Moduls: ohne Dir wird kein Wellendokument
// geöffnet und der Befundsatz ist byte-identisch (DC-QA-02).
//
// Plan-Dokument und Ergebnisnotiz sind zwei ROLLEN: die Aktiv-Status-Aussage
// fragt nach dem Plan-Dokument, die Abschluss-Aussage nach der Notiz. Gegen das
// Plan-Dokument gemessen meldete die Abschluss-Aussage über zwei reale
// Planungs-Bäume 19-mal falsch (ADR-0055).
type WavesConfig struct {
	Dir           string
	DoneDir       string
	Glob          string
	ResultsGlob   string
	NextHeading   string
	ClosedHeading string
	Mode          string
}

// EffectiveDoneDir liefert den Ruheort der Ergebnisnotizen (Default: `done`
// unterhalb von Dir).
func (w WavesConfig) EffectiveDoneDir() string {
	if w.DoneDir == "" {
		return path.Join(w.Dir, "done")
	}
	return w.DoneDir
}

// EffectiveGlob liefert den Basisnamen-Glob der Plan-Dokumente.
func (w WavesConfig) EffectiveGlob() string {
	if w.Glob == "" {
		return "welle-*.md"
	}
	return w.Glob
}

// EffectiveResultsGlob liefert den Basisnamen-Glob der Ergebnisnotizen; diese
// Menge wird von den Plan-Dokumenten ABGEZOGEN.
func (w WavesConfig) EffectiveResultsGlob() string {
	if w.ResultsGlob == "" {
		return "welle-*-results.md"
	}
	return w.ResultsGlob
}

// EffectiveNextHeading liefert die H2 des Vorschau-Registers.
func (w WavesConfig) EffectiveNextHeading() string {
	if w.NextHeading == "" {
		return "## Nächste Wellen"
	}
	return w.NextHeading
}

// EffectiveClosedHeading liefert die H2 des Abschluss-Registers.
func (w WavesConfig) EffectiveClosedHeading() string {
	if w.ClosedHeading == "" {
		return "## Abgeschlossene Wellen"
	}
	return w.ClosedHeading
}

// EffectiveMode liefert das Kardinalitäts-Modell der ersten Wellen-Aussage
// (§DC-FA-PLAN-001.a Schritt W3, Lastenheft 0.62.0): "one" = Singleton
// (Default), "many" = Kennungs-Mengen-Bijektion. Dass nur diese beiden Werte
// ankommen — der explizit leere String eingeschlossen —, stellt der
// Config-Rand sicher (Exit 2, Schritt W1).
func (w WavesConfig) EffectiveMode() string {
	if w.Mode == "" {
		return "one"
	}
	return w.Mode
}

// StructureRule ist eine Regel des Moduls structure (DC-FA-STRUCT-001): eine
// Dokumentklasse über eigene Globs, ein Abschnitts-Selektor und die optionalen
// Bedingungen. MinSentences/MaxTasks sind Zeiger, damit ein ABWESENDER
// Schlüssel (Bedingung aus) von einem explizit gesetzten Wert unterscheidbar
// bleibt — sonst wäre die Null-Schwelle unerreichbar prüfbar.
type StructureRule struct {
	Files          string
	Section        string
	SectionPattern string
	Sections       string
	NonEmpty       bool
	MinSentences   *int
	MaxTasks       *int
	ForbidPattern  string
	RequirePattern string
	RequireAll     []string
	// HeadingsMatch schaltet die Ueberschriften-Bedingung scharf; HeadingsLevel
	// ist die geprueften ATX-Ebene — Zeiger, damit der Default (Ebene des
	// Abschnitts + 1) von einem explizit gesetzten Wert unterscheidbar bleibt.
	HeadingsMatch string
	HeadingsLevel *int
	// Table buendelt die beiden TABELLEN-Bedingungen unter einem Schluessel
	// (ADR-0070). Beide sprechen ueber dieselbe Tabelle, adressieren ihre
	// Spalte aber verschieden — die Klammer macht das sichtbar, statt es fuenf
	// flachen Schluesseln zu ueberlassen.
	Table       TableRule
	ExemptPaths []string
}

// TableRule sind die tabellenbezogenen Bedingungen einer StructureRule
// (ADR-0070, Konfigurations-Schluessel `table`): die Chronologie-Monotonie
// (ADR-0057) und die Zellengrenzen je benannter Spalte (ADR-0069).
//
// ABGRENZUNG der beiden Spalten-Begriffe, die hier nebeneinander stehen:
// OrderColumn adressiert ueber die POSITION (typisiert, ein Griff daneben
// faellt als section-cell-untyped auf), Columns ueber den KOPFZEILEN-NAMEN
// (eine Laengen-Messung hat diese Selbstkontrolle nicht, deshalb meldet ein
// umbenannter Kopf laut). Der Praefix in `order-column` traegt die Kopplung an
// Order, die der Config-Rand ohnehin erzwingt.
type TableRule struct {
	// Order (asc/desc) schaltet die Chronologie-Monotonie scharf; OrderColumn
	// ist die 1-basierte Schluesselspalte — Zeiger, damit ein abwesender
	// Schluessel (Default 1) von einem explizit gesetzten unterscheidbar
	// bleibt (Exit-2-Rand: explizit < 1).
	Order       string
	OrderColumn *int
	// Columns traegt je Eintrag EINE Spalte und ihre Grenzen. Die Liste ist
	// der Grund fuer die Klammer: vor ADR-0070 kostete jede weitere Spalte
	// desselben Abschnitts eine eigene Regel samt wiederholtem Selektor.
	Columns []TableColumnRule
}

// TableColumnRule ist die Zellengrenzen-Zusage ueber EINE Spalte, adressiert
// ueber ihren Kopfzeilen-Namen (ADR-0069). CellMaxChars/CellMinChars sind Ober-
// und Untergrenze in ZEICHEN — Zeiger, damit ein abwesender Schluessel von
// einem explizit gesetzten Wert unterscheidbar bleibt. Mindestens eine der
// beiden ist Pflicht; eine Obergrenze allein laesst die LEERE Zelle passieren.
type TableColumnRule struct {
	Name         string
	CellMaxChars *int
	CellMinChars *int
}

// EffectiveHeadingsLevel liefert die geprueften ATX-Ebene der
// Ueberschriften-Bedingung: der explizit gesetzte Wert, sonst die Ebene des
// Abschnitts + 1 (spec/spezifikation.md §DC-FA-STRUCT-001.a Schritt 6).
func (r StructureRule) EffectiveHeadingsLevel(sectionLevel int) int {
	if r.HeadingsLevel != nil {
		return *r.HeadingsLevel
	}
	return sectionLevel + 1
}

// EffectiveOrderColumn liefert die 1-basierte Schluesselspalte der
// Chronologie-Bedingung (Default 1).
func (t TableRule) EffectiveOrderColumn() int {
	if t.OrderColumn == nil {
		return 1
	}
	return *t.OrderColumn
}

// Identity ist die Regel-Identität aus Glob und Abschnitts-Selektor. Sie steht
// im Befund-Ziel, weil die Deduplikation (Datei, Zeile, Regel, Ziel, Grund)
// sonst je Datei den Befund der zweiten Regel verlöre.
func (r StructureRule) Identity() string {
	sel := r.Section
	if sel == "" {
		sel = r.SectionPattern
	}
	// Die AUSDRUECKLICH gesetzte Chronologie-Spalte gehoert dazu: zwei
	// Chronologie-Zusagen ueber denselben Abschnitt sind verschiedene Zusagen
	// und brauchen verschiedene Identitaeten (ADR-0069). Nicht dazu gehoert
	// eine unbenannte Default-Spalte: sonst aenderte sich das target jeder
	// Bestandsregel, die keine nennt. Die ZELLEN-Spalten stehen NICHT hier —
	// sie leben seit ADR-0070 als Liste INNERHALB einer Regel, koennen also
	// keine zwei Regeln mehr kollidieren lassen; ihre Trennung leistet
	// ColumnTarget je Befund.
	if r.Table.OrderColumn != nil {
		sel += " :: Spalte " + strconv.Itoa(*r.Table.OrderColumn)
	}
	return r.Files + " :: " + sel
}

// ColumnTarget ist das Befund-Ziel einer SPALTEN-bezogenen Bedingung: die
// Regel-Identitaet plus die benannte Spalte. Ohne diesen Zusatz fielen zwei zu
// lange Zellen DERSELBEN Zeile unter die Deduplikation (Datei, Zeile, Regel,
// Ziel, Grund) zusammen — seit die Spalten einer Regel eine Liste sind, traegt
// die Identitaet sie nicht mehr (ADR-0070).
func (r StructureRule) ColumnTarget(name string) string {
	return r.Identity() + " :: Spalte " + name
}

// EffectiveSections liefert den Kardinalitäts-Modus (Default `one`).
func (r StructureRule) EffectiveSections() string {
	if r.Sections == "" {
		return "one"
	}
	return r.Sections
}

// ClosureConfig sind die Parameter der zweiten planning-Fähigkeit
// (DC-FA-PLAN-001 §Closure-Note-Struktur): die strukturelle Prüfung der
// Closure-Notizen abgeschlossener Slices. Dir ist der Aktivierungs-Schalter —
// leer ⇒ Fähigkeit inert, es wird KEINE Slice-Datei geöffnet (byte-identisch,
// DC-QA-02). Geprüft wird Struktur, nicht Bedeutung; die Floskel-Semantik
// bleibt dem inferentiellen Nachlauf überlassen.
type ClosureConfig struct {
	Dir            string
	Glob           string
	HeadingPattern string
	MinSentences   int
	Boilerplate    []string
	Placeholder    bool
}

// EffectiveHeadingPattern liefert das RE2-Muster der Closure-Notiz-Überschrift
// (Default deckt die gelebten Formen: H2/H3, beliebige Nummer, beliebiges Suffix).
func (c ClosureConfig) EffectiveHeadingPattern() string {
	if c.HeadingPattern == "" {
		return `^#{2,3} .*Closure-Notiz`
	}
	return c.HeadingPattern
}

// EffectiveMinSentences liefert die Mindestzahl der Satzende-Zeichen außerhalb
// der Fenced-Code-Blöcke (Default 4). Die Null steht für „nicht gesetzt", nicht
// für „Schwelle 0" — eine Null-Schwelle wäre ein stilles Grün und lehnt der
// Config-Decoder ab (Exit 2).
func (c ClosureConfig) EffectiveMinSentences() int {
	if c.MinSentences == 0 {
		return 4
	}
	return c.MinSentences
}

// EffectiveHeading liefert die kanonische H2-Überschrift (Default `## Aktuelle Welle`).
func (p PlanningConfig) EffectiveHeading() string {
	if p.Heading == "" {
		return "## Aktuelle Welle"
	}
	return p.Heading
}

// EffectiveMarker liefert den Ruhe-Marker (Default „Keine aktive Welle").
func (p PlanningConfig) EffectiveMarker() string {
	if p.Marker == "" {
		return "Keine aktive Welle"
	}
	return p.Marker
}

// EffectiveClosureGlob liefert das Basisnamen-Glob der Closure-Kandidaten
// (§DC-FA-PLAN-001.a Schritt C2). Nicht gesetzt ⇒ EffectiveSliceGlob: der
// Default ist ein Verweis, kein kopiertes Literal — es gibt genau ein Muster zu
// pflegen, solange niemand die beiden Grundmengen trennt.
func (p PlanningConfig) EffectiveClosureGlob() string {
	if p.Closure.Glob == "" {
		return p.EffectiveSliceGlob()
	}
	return p.Closure.Glob
}

// EffectiveSliceGlob liefert das Slice-Basisnamen-Glob (Default `slice-*.md`).
func (p PlanningConfig) EffectiveSliceGlob() string {
	if p.SliceGlob == "" {
		return "slice-*.md"
	}
	return p.SliceGlob
}

// TrackedConfig sind die Parameter des Moduls tracked (DC-FA-TRK-001):
// ExemptTargets nimmt aufgelöste Ziel-Pfade referenz-weit von der
// Getrackt-Prüfung aus (Glob wie scan.ignore; analog codepaths.ignore-refs) —
// für absichtlich untrackte Ziele. Ohne Einträge byte-identisch.
type TrackedConfig struct {
	ExemptTargets []string
}

// TargetsConfig sind die Parameter des Moduls targets (DC-FA-TGT-001):
// Makefiles sind die Wurzel-relativen Makefile-Quellen, aus denen Regelnamen
// per statischer Zeilen-Heuristik extrahiert werden (leer ⇒ Modul inert);
// DocTables sind die Doku-Dateien, deren `make X`-Tabellenzeilen gegen die
// Regelmenge geprüft werden (Richtung 1, gate-phantom; leer ⇒ Richtung 1
// entfällt); Authority ist die Doku-Datei, in der jede nicht-exempte Regel als
// `make X`-Tabellenzeile stehen muss (Richtung 2, gate-undocumented; leer ⇒
// Richtung 2 entfällt); ExemptTargets nimmt Regelnamen (exakt) von der
// Doku-Pflicht aus (Utility-Targets). Hermetisch — nur der Filesystem-Port.
type TargetsConfig struct {
	Makefiles     []string
	DocTables     []string
	Authority     string
	ExemptTargets []string
}

// SourcePin ist ein Config-Pin des Moduls sources (DC-FA-SRC-001): eine auf
// Sha256 gepinnte externe http(s)-Quelle. Unpack ist "none" (Roh-Byte-Hash)
// oder "zip" (Content-Manifest-Hash). Sha256 wird case-insensitiv geführt
// (im Config-Adapter zu Kleinbuchstaben normalisiert). Line ist die Zeile des
// url-Feldes in .d-check.yml für den Befund (Fallback 1).
type SourcePin struct {
	URL    string
	Sha256 string
	Unpack string
	Line   int
}

// SourcesConfig sind die Config-Pins des Moduls sources (DC-FA-SRC-001).
// Die Marker-Pins am Link werden separat beim Scan gesammelt; ohne aktives
// sources ist der Befundsatz byte-identisch (DC-QA-02) und es wird keine
// Netzverbindung geöffnet (DC-QA-03 — zweite Netz-Tür neben external).
type SourcesConfig struct {
	Pins []SourcePin
}

// SourceUnpackZip ist der unpack-Wert für Archiv-Ziele (DC-FA-SRC-001).
const SourceUnpackZip = "zip"

// TraceConfig sind die konfigurierbaren Quellen der Requirements Traceability
// Matrix (DC-FA-CLI-009/DC-FA-REQ-001): Source ist die Anforderungs-Quelldatei;
// Format/Table wählen Heading- oder Tabellenextraktion; ReqPattern erkennt eine
// Anforderungs-Kennung (Ganz-Token/-Zelle UND als Referenz in ADR/Slice-Dateien);
// ADRDir/SliceDir sind die Referenz-Verzeichnisse;
// ADRFile/SliceFile leiten je Datei die Owner-Kennung über den Basisnamen ab
// (Capture-Gruppe 1), ADRPrefix/SlicePrefix werden ihr vorangestellt. Reiner
// Daten-Struct: der **Nullwert** (alle Felder leer/nil) bedeutet „Konventions-
// Default"; die Auflösung der Defaults liegt im app-Kern (trace.go), wo die
// Default-Regex/-Pfade als Konstanten leben. Ohne trace-Block ⇒ RTM
// byte-identisch (DC-QA-02).
type TraceConfig struct {
	Source      string
	ReqPattern  *regexp.Regexp
	Format      string
	Table       *TraceTableConfig
	ADRDir      string
	ADRFile     *regexp.Regexp
	ADRPrefix   string
	SliceDir    string
	SliceFile   *regexp.Regexp
	SlicePrefix string
	// Coverage: opt-in kuratierte Coverage-Quellen (DC-FA-COV-001); leer ⇒
	// keine Coverage-Spalte, RTM byte-identisch.
	Coverage []TraceCoverage
	// Modality: opt-in Modalitäts-Klassifikation (DC-FA-MOD-001); nil ⇒ aus
	// (byte-identisch). Präsenz (auch leere Map `modality: {}`) ⇒ aktiv, dann
	// greifen die Default-Keywords.
	Modality *TraceModality
	// CrossConsistency: opt-in Kreuzverweis-Abgleich zweier Traceability-Sichten
	// (DC-FA-XREF-001); nil ⇒ kein Abgleich, RTM byte-identisch.
	CrossConsistency *TraceCrossConsistency
}

const (
	// TraceCrossModeEqual verlangt F = B und meldet beide Richtungsdifferenzen
	// (Default, DC-FA-XREF-001).
	TraceCrossModeEqual = "equal"
	// TraceCrossModeSuperset verlangt nur F ⊇ B — allein B\F ist ein Befund.
	TraceCrossModeSuperset = "superset"
	// TraceCrossArtifactFirst ist der positionelle Sentinel der Rückwärts-ID-Spalte
	// (erste Spalte statt Header-Name) — begründet durch die über die Tabellen
	// heterogenen ID-Header (`Kennung`/`Port-ID`/…), während die Kanten-Spalte
	// einheitlich heißt (ADR-0038).
	TraceCrossArtifactFirst = "first"
)

// TraceCrossConsistency ist der opt-in Mengenabgleich zweier unabhängig
// gepflegter Sichten derselben Anforderung→Design-Relation (DC-FA-XREF-001):
// Forward ist die kuratierte Vorwärts-RTM-Tabelle (Anforderung → Design-Menge),
// Backward sind die Rück-Kanten (Design → Anforderung) und die **Quelle der
// Wahrheit**. Mode ist TraceCrossModeEqual oder TraceCrossModeSuperset;
// ExcludeReq nimmt Ableitungssprung-Kennungen (Mittelschicht-Familien) vor dem
// Diff aus beiden Sichten (nil ⇒ keine Ausnahme).
type TraceCrossConsistency struct {
	Forward    TraceCrossForward
	Backward   TraceCrossBackward
	Mode       string
	ExcludeReq *regexp.Regexp
}

// TraceCrossForward beschreibt die Vorwärts-Sicht (DC-FA-XREF-001): File ist die
// Wurzel-relative Tabellendatei, Sections/ExcludeSections spannen sie über die
// Heading-Span-Semantik (wie matrix.exclude-sections). ReqColumn und DesignColumn
// sind header-gebundene Spaltennamen; DesignPattern extrahiert die
// Design-Artefakt-Kennungen aus der Design-Zelle und wird bewusst mit der
// Rückwärts-Sicht geteilt (gemeinsamer Namensraum — sonst wäre der Diff
// bedeutungslos). Ranges aktiviert die Range-/Enum-Expansion der Anforderungs-IDs.
type TraceCrossForward struct {
	File            string
	Sections        []string
	ExcludeSections []string
	ReqColumn       string
	DesignColumn    string
	DesignPattern   *regexp.Regexp
	// ReqPattern erkennt die Anforderungs-IDs in der ReqColumn-Zelle — symmetrisch
	// zu TraceCrossBackward.ReqPattern. nil ⇒ Default `requirements.id-pattern`
	// (aufgelöst im app-Kern). Eigenständig, nicht abgeleitet: welche Anforderungen
	// der Abgleich vergleicht, entscheidet das Muster — **nicht** die
	// RTM-Mitgliedschaft (DC-FA-XREF-001).
	ReqPattern *regexp.Regexp
	Ranges     bool
}

// TraceCrossBackward beschreibt die Rück-Kanten-Sicht (DC-FA-XREF-001): File ist
// die Wurzel-relative Quelle, Sections spannt sie (Whitelist). ArtifactIDColumn
// ist ein Header-Name oder der Sentinel TraceCrossArtifactFirst (erste Spalte);
// EdgeColumn ist der header-gebundene Name der Kanten-Spalte (z. B. `Bezug`),
// ReqPattern erkennt die Anforderungs-Kennungen in deren Zelle. Ranges aktiviert
// die Range-/Enum-Expansion.
type TraceCrossBackward struct {
	File             string
	Sections         []string
	ArtifactIDColumn string
	EdgeColumn       string
	ReqPattern       *regexp.Regexp
	Ranges           bool
}

const (
	// TraceFormatHeadings ist die bestehende ATX-Heading-Grammatik und der
	// Default bei leerem Format (byte-identisch, DC-QA-02).
	TraceFormatHeadings = "headings"
	// TraceFormatTable aktiviert Markdown-Pipe-Tabellen (DC-FA-REQ-001).
	TraceFormatTable = "table"
	// TraceDuplicateError ist die sichere Default-Politik für doppelte
	// Tabellen-IDs; First/Last sind explizite Brownfield-Overrides.
	TraceDuplicateError = "error"
	// TraceDuplicateFirst behält bei Mehrfachdefinitionen die erste Zeile.
	TraceDuplicateFirst = "first"
	// TraceDuplicateLast behält bei Mehrfachdefinitionen die letzte Zeile.
	TraceDuplicateLast = "last"
)

// TraceTableConfig bindet die Rollen einer Requirement-Tabelle an exakte
// Header-Namen. TextColumns enthält mindestens einen alternativen Header;
// ModalityColumn ist optional. DuplicateIDs ist error, first oder last.
type TraceTableConfig struct {
	IDColumn       string
	TextColumns    []string
	ModalityColumn string
	DuplicateIDs   string
}

// TraceModality sind die Parameter der Modalitäts-Klassifikation (DC-FA-MOD-001):
// Levels bildet Stufen-Name → Modal-Verb-Keywords ab (leer ⇒ Built-in DE+EN-
// RFC-2119-Default im app-Kern); RequireLevels nennt die Stufen, deren Waisen
// `--require-complete` gaten (leer ⇒ Default `[must]`).
type TraceModality struct {
	Levels        map[string][]string
	RequireLevels []string
}

// DefaultModalityLevels ist die kanonische Built-in-DE+EN-RFC-2119-Keyword-Menge
// (DC-FA-MOD-001) — genutzt, wenn `modality.levels` leer ist. Deterministisch
// (frische Map je Aufruf); dup-frei über die Stufen. Die DE-Negations-Asymmetrie
// ist bewusst: `DARF NICHT` = Verbot (must), `MUSS NICHT` = braucht-nicht (may).
func DefaultModalityLevels() map[string][]string {
	return map[string][]string{
		"must":   {"MUSS", "MUESSEN", "MÜSSEN", "DARF NICHT", "DÜRFEN NICHT", "MUST", "SHALL", "MUST NOT", "SHALL NOT"},
		"should": {"SOLLTE", "SOLLTEN", "SOLLTE NICHT", "SOLLTEN NICHT", "SHOULD", "SHOULD NOT"},
		"may":    {"KANN", "KÖNNEN", "MUSS NICHT", "MÜSSEN NICHT", "MAY", "OPTIONAL"},
	}
}

// TraceCoverage ist eine kuratierte Coverage-Quelle der RTM (DC-FA-COV-001):
// Files sind die explizit benannten Quell-Dateien (keine dir/pattern-Ableitung —
// gegen ADR-Kontamination); Label ist die feste Owner-Kennung in der
// Coverage-Spalte; Ranges (Default true) aktiviert die `<FAM>-AAA..BBB`-/
// `/`-Enum-Expansion; Sections/ExcludeSections scopen den gescannten Text über
// die Heading-Span-Semantik (voller Heading-Klartext, wie matrix.exclude-sections).
type TraceCoverage struct {
	Files           []string
	Label           string
	Ranges          bool
	Sections        []string
	ExcludeSections []string
}

// EffectiveModules wendet die Modul-Auflösung an
// (spec/spezifikation.md §DC-FA-CLI-002.a): (Config oder Default)
// ∪ enable ∖ disable; CLI nach Config. Unbekannte Namen → Fehler mit
// Liste der gültigen Namen.
func EffectiveModules(cfg Config, enable, disable []string) ([]string, error) {
	valid := map[string]bool{}
	for _, m := range validModules() {
		valid[m] = true
	}
	for _, m := range append(append([]string{}, enable...), disable...) {
		if !valid[m] {
			return nil, fmt.Errorf("unbekanntes Modul %q — gültig: %s",
				m, strings.Join(validModules(), ", "))
		}
	}
	base := cfg.Modules
	if base == nil {
		base = defaultModules()
	}
	set := map[string]bool{}
	for _, m := range base {
		if !valid[m] {
			return nil, fmt.Errorf("unbekanntes Modul %q in der Konfiguration — gültig: %s",
				m, strings.Join(validModules(), ", "))
		}
		set[m] = true
	}
	for _, m := range enable {
		set[m] = true
	}
	for _, m := range disable {
		delete(set, m)
	}
	var out []string
	for m := range set {
		out = append(out, m)
	}
	sort.Strings(out)
	return out, nil
}
