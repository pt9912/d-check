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

// variableHostpathRE erkennt `$NAME/…`, `${NAME}/…` und `"$NAME"/…` für die
// konfigurierten Variablen-Namen (DC-FA-HOST-001): dasselbe erste Segment wie
// bei der Tilde (nicht `.`, nicht `/`), die Vorbedingung der Tilde plus `$` und
// `\` (kein `$$HOME`, keine maskierte Variable). Namen ohne Schrägstrich direkt
// dahinter (`$HOMEDIR/x`, `${HOME:-/tmp}/x`, `${HOME}x/y`) treffen nicht. Leere
// Liste ⇒ nil, das Modul bleibt dann auf dem Verhalten ohne Variablen.
func variableHostpathRE(names []string) *regexp.Regexp {
	if len(names) == 0 {
		return nil
	}
	quoted := make([]string, len(names))
	for i, n := range names {
		quoted[i] = regexp.QuoteMeta(n)
	}
	alt := strings.Join(quoted, "|")
	return regexp.MustCompile(
		`(^|[^A-Za-z0-9_.:/~$\\-])(\$(?:(?:` + alt + `)|\{(?:` + alt + `)\})"?/[^.\s<>)\]"'/` + "`" + `][^\s<>)\]"'` + "`" + `]*)`)
}

// CheckHostpaths meldet host-lokale absolute und Home-relative Pfade in
// Prosa und Inline-Code (DC-FA-HOST-001, spec/spezifikation.md
// §DC-FA-HOST-001.a). Fenced-Code-Blöcke sind ausgenommen — dort
// gehören bewusste Beispiel-Pfade hin; es gibt keinen Zeilen-Marker.
// Ein Unix-Treffer, der innerhalb eines Home-relativen oder Variablen-Treffers
// beginnt, entfällt (genau ein Befund, in voller Form). Variablen-Funde gibt es nur mit
// `variables`; das schließende `"` von `"$HOME"/x` gehört nicht zum gemeldeten
// Pfad. exempt-targets greift für
// Unix-, Tilde- und Variablen-Funde, nicht für die festen Windows-/UNC-Muster.
func CheckHostpaths(file string, content []byte, cfg model.HostpathsConfig) []model.Finding {
	unixRE := unixHostpathRE(cfg)
	varRE := variableHostpathRE(cfg.Variables)
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
		// Spannen, die einen Unix-Treffer verdrängen: Tilde und Variable.
		homes := tildeRE.FindAllStringSubmatchIndex(pl.raw, -1)
		for _, m := range homes {
			add(pl, pl.raw[m[4]:m[5]], true)
		}
		for _, m := range variableMatches(pl.raw, varRE) {
			homes = append(homes, m)
			add(pl, withoutClosingQuote(pl.raw[m[4]:m[5]]), true)
		}
		for _, m := range unixRE.FindAllStringSubmatchIndex(pl.raw, -1) {
			if !insideSpan(m[4], homes) {
				add(pl, pl.raw[m[4]:m[5]], true)
			}
		}
		for _, p := range windowsPaths(pl.raw) {
			add(pl, p, false)
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

// variableMatches liefert die Variablen-Treffer einer Zeile; ohne konfigurierte
// Variablen (re == nil) keinen.
func variableMatches(raw string, re *regexp.Regexp) [][]int {
	if re == nil {
		return nil
	}
	return re.FindAllStringSubmatchIndex(raw, -1)
}

// withoutClosingQuote entfernt das schließende `"` von `"$HOME"/x` aus dem
// gemeldeten Pfad; das Restpfad-Zeichen schließt `"` aus, es gibt höchstens eines.
func withoutClosingQuote(path string) string {
	return strings.Replace(path, `"/`, "/", 1)
}

// windowsPaths liefert die Windows-Laufwerks- und UNC-Pfade einer Zeile.
func windowsPaths(raw string) []string {
	var out []string
	for _, re := range []*regexp.Regexp{windowsDriveRE, windowsUNCRE} {
		for _, m := range re.FindAllStringSubmatch(raw, -1) {
			out = append(out, m[2])
		}
	}
	return out
}
