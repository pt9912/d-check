package core

import (
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls matrix (spec/spezifikation.md §4).
const (
	ReasonMatrixForbidden = "matrix-forbidden"
	ReasonMatrixInactive  = "matrix-inactive"
)

// checkMatrix ist das Regelmodul `matrix` (DC-FA-MTX-001):
// Referenzregeln zwischen Dokumentklassen plus Status-Bedingungen.
// Dateien ohne Klasse nehmen nicht an der Prüfung teil; Links in
// ausgenommenen Sektionen (exclude-sections) werden übersprungen
// (spec/spezifikation.md §DC-FA-MTX-001.a).
func checkMatrix(fsys driven.Filesystem, file string, content []byte, lines []Line,
	cfg MatrixConfig, statusCache map[string]*string) []Finding {
	srcClass, ok := classOf(cfg.Classes, file)
	if !ok {
		return nil
	}
	excluded := excludedRanges(content, cfg.ExcludeSections)
	var findings []Finding
	for _, ref := range ExtractLinks(lines) {
		if inRanges(excluded, ref.Line) {
			continue // Provenance-Ausnahme
		}
		rel, escaped, ok := localTarget(file, ref.Target)
		if !ok || escaped { // Repo-Escape meldet `links`
			continue
		}
		dstClass, ok := classOf(cfg.Classes, rel)
		if !ok {
			continue
		}
		if rule, found := ruleFor(cfg.Rules, srcClass, dstClass); found && !rule.Allow {
			findings = append(findings, Finding{
				File: file, Line: ref.Line, Rule: "matrix",
				Target: ref.Target, Reason: ReasonMatrixForbidden,
				Message: "Referenz " + srcClass + " → " + dstClass + " ist nicht erlaubt",
			})
		}
		if status := cachedStatus(fsys, statusCache, rel); statusForbidden(status, cfg.StatusForbidden) {
			findings = append(findings, Finding{
				File: file, Line: ref.Line, Rule: "matrix",
				Target: ref.Target, Reason: ReasonMatrixInactive,
				Message: "Referenz auf inaktives Dokument (Status: " + status + ")",
			})
		}
	}
	return findings
}

// classOf liefert die erste deklarierte Klasse, deren Glob den Pfad
// matcht (Deklarationsreihenfolge = Präzedenz).
func classOf(classes []MatrixClass, rel string) (string, bool) {
	for _, c := range classes {
		for _, g := range c.Paths {
			if matchGlob(g, rel) {
				return c.Name, true
			}
		}
	}
	return "", false
}

// ruleFor liefert die erste deklarierte Regel für das Klassen-Paar.
func ruleFor(rules []MatrixRule, from, to string) (MatrixRule, bool) {
	for _, r := range rules {
		if r.From == from && r.To == to {
			return r, true
		}
	}
	return MatrixRule{}, false
}

// lineRange ist ein 1-basierter Zeilenbereich [from, to);
// to == 0 bedeutet: bis Dateiende.
type lineRange struct{ from, to int }

// excludedRanges liefert die Zeilenbereiche der per exclude-sections
// ausgenommenen Sektionen: vom Heading bis zum nächsten Heading
// gleicher oder höherer Ebene.
func excludedRanges(content []byte, sections []string) []lineRange {
	if len(sections) == 0 {
		return nil
	}
	ex := map[string]bool{}
	for _, s := range sections {
		ex[s] = true
	}
	headings := extractHeadingLines(content)
	var out []lineRange
	for i, h := range headings {
		if !ex[plainHeadingText(h.text)] {
			continue
		}
		r := lineRange{from: h.line}
		for _, h2 := range headings[i+1:] {
			if h2.level <= h.level {
				r.to = h2.line
				break
			}
		}
		out = append(out, r)
	}
	return out
}

func inRanges(rs []lineRange, line int) bool {
	for _, r := range rs {
		if line >= r.from && (r.to == 0 || line < r.to) {
			return true
		}
	}
	return false
}

// plainHeadingText liefert den getrimmten Heading-Text ohne
// Markdown-Auszeichnung (Links → Linktext, Backticks/Sterne entfernt)
// — Vergleichsbasis für exclude-sections (case-sensitiv) und das
// Status-Heading (case-insensitiv).
func plainHeadingText(s string) string {
	t := stripHeadingLinks(s)
	t = strings.Map(func(r rune) rune {
		if r == '`' || r == '*' {
			return -1
		}
		return r
	}, t)
	return strings.TrimSpace(t)
}

// statusOf extrahiert das Status-Feld eines Dokuments in fester
// Reihenfolge (spec/spezifikation.md §DC-FA-MTX-001.a Schritt 2):
// (1) erste Zeile, die mit `**Status:**` beginnt; (2) sonst erste
// nicht-leere Textzeile unter einem Status-Heading (beliebige Ebene,
// case-insensitiv). Beide Formen lesen nur Prosa-Zeilen — Inhalte in
// Fenced-Code-Blöcken sind keine Statuswerte. "" = kein Status
// (Dokument gilt als aktiv).
func statusOf(content []byte) string {
	prose := proseLines(content)
	for _, pl := range prose {
		trimmed := strings.TrimLeft(pl.raw, " \t")
		if strings.HasPrefix(trimmed, "**Status:**") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "**Status:**"))
		}
	}
	return statusUnderHeading(prose, extractHeadingLines(content))
}

// statusUnderHeading liefert die erste nicht-leere Prosa-Zeile unter
// dem ersten Status-Heading (Form 2 der Status-Extraktion).
func statusUnderHeading(prose []proseLine, headings []headingLine) string {
	for i, h := range headings {
		if !strings.EqualFold(plainHeadingText(h.text), "Status") {
			continue
		}
		end := 0 // 0 = bis Dateiende
		if i+1 < len(headings) {
			end = headings[i+1].line
		}
		for _, pl := range prose {
			if pl.no <= h.line || (end != 0 && pl.no >= end) {
				continue
			}
			if t := strings.TrimSpace(pl.raw); t != "" {
				return t
			}
		}
		return ""
	}
	return ""
}

// statusForbidden prüft den Status-Wert case-insensitiv als
// Präfix-Match gegen die verbotenen Werte (so matcht
// "Superseded by ADR-0042" den Wert "superseded").
func statusForbidden(status string, forbidden []string) bool {
	if status == "" {
		return false
	}
	for _, f := range forbidden {
		if f != "" && len(status) >= len(f) && strings.EqualFold(status[:len(f)], f) {
			return true
		}
	}
	return false
}

// cachedStatus liefert den Status der Zieldatei aus dem Cache bzw.
// liest ihn nach (nil = nicht lesbar → matrix schweigt, `links`
// meldet ein fehlendes Ziel). Status wird nur aus Markdown-Zielen
// extrahiert — andere Ziele gelten als aktiv (kein Voll-Read von
// Binärdateien; spec/spezifikation.md §DC-FA-MTX-001.a).
func cachedStatus(fsys driven.Filesystem, cache map[string]*string, rel string) string {
	if !strings.HasSuffix(rel, ".md") {
		return ""
	}
	if s, ok := cache[rel]; ok {
		if s == nil {
			return ""
		}
		return *s
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		cache[rel] = nil
		return ""
	}
	s := statusOf(content)
	cache[rel] = &s
	return s
}
