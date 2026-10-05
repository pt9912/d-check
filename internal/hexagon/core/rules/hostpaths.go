package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
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

// tildeRE erkennt Home-relative Pfade: `~/` + erstes Segment, das nicht
// mit `.` beginnt (Werkzeug-Konventionen wie ~/.config bleiben still).
// Die Vorbedingung ist die des Unix-Musters plus `~` (eine Strikethrough-
// Tilde `~~` ist kein Home-Verweis); `~user/` trifft nicht, weil `/`
// unmittelbar auf die Tilde folgen muss.
var tildeRE = regexp.MustCompile(`(^|[^A-Za-z0-9_.:/~-])(~/[^.\s<>)\]"'/` + "`" + `][^\s<>)\]"'` + "`" + `]*)`)

// CheckHostpaths meldet host-lokale absolute und Home-relative Pfade in
// Prosa und Inline-Code (DC-FA-HOST-001, spec/spezifikation.md
// §DC-FA-HOST-001.a). Fenced-Code-Blöcke sind ausgenommen — dort
// gehören bewusste Beispiel-Pfade hin; es gibt keinen Zeilen-Marker.
// Ein Unix-Treffer, der innerhalb eines Home-relativen Treffers beginnt,
// entfällt (genau ein Befund, in voller Form). exempt-targets greift für
// Unix- und Tilde-Funde, nicht für die festen Windows-/UNC-Muster.
func CheckHostpaths(file string, content []byte, cfg model.HostpathsConfig) []model.Finding {
	unixRE := unixHostpathRE(cfg)
	var findings []model.Finding
	add := func(pl proseLine, raw string, exemptable bool) {
		path := strings.TrimRight(raw, ".,;:")
		if path == "" || (exemptable && exemptTarget(path, cfg.ExemptTargets)) {
			return
		}
		findings = append(findings, model.Finding{
			File: file, Line: pl.no, Rule: "hostpaths",
			Target: path, Reason: model.ReasonHostpathForbidden,
		})
	}
	for _, pl := range proseLines(content) {
		tildes := tildeRE.FindAllStringSubmatchIndex(pl.raw, -1)
		for _, m := range tildes {
			add(pl, pl.raw[m[4]:m[5]], true)
		}
		for _, m := range unixRE.FindAllStringSubmatchIndex(pl.raw, -1) {
			if !insideSpan(m[4], tildes) {
				add(pl, pl.raw[m[4]:m[5]], true)
			}
		}
		for _, re := range []*regexp.Regexp{windowsDriveRE, windowsUNCRE} {
			for _, m := range re.FindAllStringSubmatch(pl.raw, -1) {
				add(pl, m[2], false)
			}
		}
	}
	return findings
}

// insideSpan prüft, ob pos in der Pfad-Gruppe eines der Treffer liegt.
func insideSpan(pos int, matches [][]int) bool {
	for _, m := range matches {
		if pos >= m[4] && pos < m[5] {
			return true
		}
	}
	return false
}

// unixHostpathRE baut das Unix-Muster aus der (konfigurierbaren)
// Präfixliste; die Wortgrenzen-Vorbedingung schließt Buchstaben,
// Ziffern, `_`, `.`, `:`, `/`, `-` aus — URL-Pfade hinter Schemata
// matchen damit nicht.
func unixHostpathRE(cfg model.HostpathsConfig) *regexp.Regexp {
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
