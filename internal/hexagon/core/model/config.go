package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// validModules sind die vertraglich gültigen Regelmodul-Namen
// (DC-FA-CLI-002).
func validModules() []string {
	return []string{"links", "anchors", "ids", "matrix", "external", "codepaths", "spans", "hostpaths", "diagrams", "versions", "pins", "immutable", "vcs", "commits", "planning", "tracked", "targets"}
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
	// Tracked: Parameter des Moduls tracked (DC-FA-TRK-001).
	Tracked TrackedConfig
	// Targets: Parameter des Moduls targets (DC-FA-TGT-001).
	Targets TargetsConfig
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
}

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
}

// EffectiveFences liefert die zu öffnenden Fence-Sprachen (Default
// "mermaid", DC-FA-DIAG-001.a).
func (d DiagramsConfig) EffectiveFences() []string {
	if len(d.Fences) == 0 {
		return []string{"mermaid"}
	}
	return d.Fences
}

// VersionsConfig sind die Parameter des Moduls versions (DC-FA-VER-001):
// PinPattern erkennt einen Versions-Pin (Version in Capture-Gruppe 1, sonst
// der ganze Treffer), CurrentFrom adressiert den Span mit der aktuellen
// Version (Default `version.md#aktuell`), ExemptPaths nimmt ganze Dateien aus.
// Ohne PinPattern ist das Modul wirkungslos (byte-identisch zum Lauf ohne das
// Modul, DC-QA-02).
type VersionsConfig struct {
	PinPattern  *regexp.Regexp
	CurrentFrom string
	ExemptPaths []string
}

// EffectiveCurrentFrom liefert die Quelle der aktuellen Version (Default
// `version.md#aktuell`, spec/spezifikation.md §DC-FA-VER-001.a).
func (v VersionsConfig) EffectiveCurrentFrom() string {
	if v.CurrentFrom == "" {
		return "version.md#aktuell"
	}
	return v.CurrentFrom
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
