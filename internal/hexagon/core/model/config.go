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
	return []string{"links", "anchors", "ids", "matrix", "external", "codepaths", "spans", "hostpaths", "diagrams", "versions", "pins", "immutable", "vcs"}
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
	// Scopes: modul-lokale Scan-Scopes (DC-FA-CONF-002); Schlüssel
	// ist der Modulname, nil-Eintrag/fehlender Schlüssel = globaler
	// Scope.
	Scopes map[string]*ScopeConfig
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
