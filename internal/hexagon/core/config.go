package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ValidModules sind die vertraglich gültigen Regelmodul-Namen
// (DC-FA-CLI-002).
var ValidModules = []string{"links", "anchors", "ids", "matrix", "external"}

// ImplementedModules sind die in diesem Stand lauffähigen Module.
// `ids`/`matrix`/`external` folgen mit der Regelmodul-Welle; bis
// dahin werden aktivierte, aber nicht implementierte Module mit
// stderr-Hinweis übersprungen.
var ImplementedModules = map[string]bool{"links": true, "anchors": true}

// DefaultModules ist der Default-Modulsatz (DC-FA-CLI-002).
var DefaultModules = []string{"links", "anchors"}

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
}

// IDPattern ist ein deklariertes Kennungs-Muster (DC-FA-ID-001).
type IDPattern struct {
	Regex  *regexp.Regexp
	Target string
}

// EffectiveModules wendet die Modul-Auflösung an
// (spec/spezifikation.md §DC-FA-CLI-002.a): (Config oder Default)
// ∪ enable ∖ disable; CLI nach Config. Unbekannte Namen → Fehler mit
// Liste der gültigen Namen.
func EffectiveModules(cfg Config, enable, disable []string) ([]string, error) {
	valid := map[string]bool{}
	for _, m := range ValidModules {
		valid[m] = true
	}
	for _, m := range append(append([]string{}, enable...), disable...) {
		if !valid[m] {
			return nil, fmt.Errorf("unbekanntes Modul %q — gültig: %s",
				m, strings.Join(ValidModules, ", "))
		}
	}
	base := cfg.Modules
	if base == nil {
		base = DefaultModules
	}
	set := map[string]bool{}
	for _, m := range base {
		if !valid[m] {
			return nil, fmt.Errorf("unbekanntes Modul %q in der Konfiguration — gültig: %s",
				m, strings.Join(ValidModules, ", "))
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
