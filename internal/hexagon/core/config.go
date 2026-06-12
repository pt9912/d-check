package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// validModules sind die vertraglich gültigen Regelmodul-Namen
// (DC-FA-CLI-002).
func validModules() []string {
	return []string{"links", "anchors", "ids", "matrix", "external", "codepaths", "spans"}
}

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
}

// MatrixClass ist eine über Pfad-Globs deklarierte Dokumentklasse
// (DC-FA-MTX-001; Deklarationsreihenfolge = Präzedenz).
type MatrixClass struct {
	Name  string
	Paths []string
}

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
	ExcludeSections []string
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
