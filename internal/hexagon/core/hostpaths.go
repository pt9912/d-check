package core

import (
	"regexp"
	"strings"
)

// defaultHostPrefixes ist die Default-Präfixliste
// (spec/spezifikation.md §DC-FA-HOST-001.a Schritt 2).
func defaultHostPrefixes() []string {
	return []string{"Development", "home", "Users", "Volumes", "mnt", "media"}
}

// Windows-Laufwerks- und UNC-Muster sind fest (nicht konfigurierbar);
// die Wortgrenzen-Vorbedingung steckt in der ersten Gruppe (RE2 kennt
// kein Lookbehind), Gruppe 2 ist der Pfad.
var (
	windowsDriveRE = regexp.MustCompile(`(^|[^A-Za-z0-9_])([A-Za-z]:\\[^\s<>)\]"'` + "`" + `]*)`)
	windowsUNCRE   = regexp.MustCompile(`(^|[^A-Za-z0-9_])(\\\\[A-Za-z0-9][A-Za-z0-9_.-]*\\[^\s<>)\]"'` + "`" + `]*)`)
)

// CheckHostpaths meldet host-lokale absolute Pfade in Prosa und
// Inline-Code (DC-FA-HOST-001, spec/spezifikation.md
// §DC-FA-HOST-001.a). Fenced-Code-Blöcke sind ausgenommen — dort
// gehören bewusste Beispiel-Pfade hin; es gibt keinen Opt-out-Marker.
func CheckHostpaths(file string, content []byte, cfg HostpathsConfig) []Finding {
	unixRE := unixHostpathRE(cfg)
	var findings []Finding
	for _, pl := range proseLines(content) {
		for _, re := range []*regexp.Regexp{unixRE, windowsDriveRE, windowsUNCRE} {
			for _, m := range re.FindAllStringSubmatch(pl.raw, -1) {
				path := strings.TrimRight(m[2], ".,;:")
				if path == "" {
					continue
				}
				findings = append(findings, Finding{
					File: file, Line: pl.no, Rule: "hostpaths",
					Target: path, Reason: ReasonHostpathForbidden,
				})
			}
		}
	}
	return findings
}

// unixHostpathRE baut das Unix-Muster aus der (konfigurierbaren)
// Präfixliste; die Wortgrenzen-Vorbedingung schließt Buchstaben,
// Ziffern, `_`, `.`, `:`, `/`, `-` aus — URL-Pfade hinter Schemata
// matchen damit nicht.
func unixHostpathRE(cfg HostpathsConfig) *regexp.Regexp {
	prefixes := cfg.Prefixes
	if prefixes == nil {
		prefixes = defaultHostPrefixes()
	}
	quoted := make([]string, len(prefixes))
	for i, p := range prefixes {
		quoted[i] = regexp.QuoteMeta(p)
	}
	return regexp.MustCompile(
		`(^|[^A-Za-z0-9_.:/-])(/(?:` + strings.Join(quoted, "|") + `)/[^\s<>)\]"'` + "`" + `]*)`)
}
