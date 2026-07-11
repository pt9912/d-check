package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls matrix (spec/spezifikation.md §4).
const (
	ReasonMatrixForbidden = "matrix-forbidden"
	ReasonMatrixInactive  = "matrix-inactive"
	ReasonMatrixDownward  = "matrix-downward"
)

// provenanceMarker nimmt eine verbotene Token-Referenz aus (DC-FA-MTX-003):
// die Autor-Deklaration „dies ist Provenance, keine Entscheidungsgrundlage".
const provenanceMarker = "d-check:status-provenance"

// linkSpanRe entfernt Markdown-Link-Spans `[text](ziel)` vor der Token-Suche —
// Token in Links deckt die Link-Prüfung ab (kein Doppelbefund).
var linkSpanRe = regexp.MustCompile(`\[[^\]]*\]\([^)]*\)`)

// CheckMatrix ist das Regelmodul `matrix` (DC-FA-MTX-001):
// Referenzregeln zwischen Dokumentklassen plus Status-Bedingungen.
// Dateien ohne Klasse nehmen nicht an der Prüfung teil; Links in
// ausgenommenen Sektionen (exclude-sections) werden übersprungen
// (spec/spezifikation.md §DC-FA-MTX-001.a).
func CheckMatrix(fsys driven.Filesystem, file string, content []byte, lines []Line,
	cfg model.MatrixConfig, statusCache map[string]*string) []model.Finding {
	if ignored(file, cfg.ExemptPaths) {
		return nil // grandfathered: ganz übersprungen (DC-FA-MTX-003)
	}
	srcClass, ok := classOf(cfg.Classes, file)
	if !ok {
		return nil
	}
	// srcOrdered != nil ⇒ die Quell-Klasse trägt eine aktive klasseninterne
	// Richtung (DC-FA-MTX-002); dann wird der Rang gegen das Ziel geprüft.
	srcOrdered := orderedClass(cfg.Classes, srcClass)
	excluded := excludedRanges(content, cfg.ExcludeSections)
	// Supersede-Lineage: einmal pro Quelldatei die Feld-Werte gewinnen
	// (No-op ohne Flag, Befundsatz byte-identisch, DC-QA-02).
	supersedeValues := lineageValues(cfg, content)
	var findings []model.Finding
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
			findings = append(findings, model.Finding{
				File: file, Line: ref.Line, Rule: "matrix",
				Target: ref.Target, Reason: ReasonMatrixForbidden,
				Message: "Referenz " + srcClass + " → " + dstClass + " ist nicht erlaubt",
			})
		}
		// Klasseninterne Verweisrichtung (DC-FA-MTX-002), unabhängig von
		// matrix-forbidden/-inactive.
		if f, ok := downwardFinding(srcOrdered, srcClass, dstClass, file, rel, ref.Line, ref.Target); ok {
			findings = append(findings, f)
		}
		if status := cachedStatus(fsys, statusCache, rel); statusForbidden(status, cfg.StatusForbidden) {
			// Supersede-Lineage: die ablösende Datei darf auf das von
			// ihr abgelöste (inaktive) Ziel verweisen — die Ausnahme
			// gilt nur für matrix-inactive, nicht für matrix-forbidden.
			if !lineageExempt(supersedeValues, ref.Text, rel) {
				findings = append(findings, model.Finding{
					File: file, Line: ref.Line, Rule: "matrix",
					Target: ref.Target, Reason: ReasonMatrixInactive,
					Message: "Referenz auf inaktives Dokument (Status: " + status + ")",
				})
			}
		}
	}
	// Token-Referenzen (DC-FA-MTX-003): verbotene Referenzen als bare ID-Token
	// im Prosa-Körper, sofern nicht per Provenance-Marker deklariert.
	findings = append(findings, tokenFindings(srcClass, cfg, content, excluded, file)...)
	return findings
}

// tokenFindings erkennt verbotene Referenzen als bare ID-Token (DC-FA-MTX-003):
// je Prosa-Zeile außerhalb exclude-sections und ohne Provenance-Marker wird der
// Link-befreite Text gegen das token-Regex jeder anderen Klasse geprüft; eine
// verbotene Kante (`rules`) erzeugt matrix-forbidden. Token in Markdown-Links
// (linkSpanRe-entfernt) und in Fences (proseLines) zählen nicht.
func tokenFindings(srcClass string, cfg model.MatrixConfig, content []byte,
	excluded []lineRange, file string) []model.Finding {
	var findings []model.Finding
	for _, pl := range proseLines(content) {
		if inRanges(excluded, pl.no) || strings.Contains(pl.raw, provenanceMarker) {
			continue
		}
		stripped := linkSpanRe.ReplaceAllString(pl.raw, " ")
		for _, c := range cfg.Classes {
			if c.Token == nil || c.Name == srcClass {
				continue
			}
			if rule, found := ruleFor(cfg.Rules, srcClass, c.Name); !found || rule.Allow {
				continue
			}
			for _, loc := range c.Token.FindAllStringIndex(stripped, -1) {
				tok := stripped[loc[0]:loc[1]]
				findings = append(findings, model.Finding{
					File: file, Line: pl.no, Rule: "matrix",
					Target: tok, Reason: ReasonMatrixForbidden,
					Message: "Token-Referenz " + srcClass + " → " + c.Name + " (" + tok +
						") ist nicht erlaubt — Provenance via <!-- d-check:status-provenance --> deklarieren",
				})
			}
		}
	}
	return findings
}

// lineageValues liefert die Supersede-Feldwerte der Quelldatei, falls das
// Flag aktiv ist — sonst nil (DC-QA-02: ohne Flag byte-identisch).
func lineageValues(cfg model.MatrixConfig, content []byte) []string {
	if cfg.AllowSupersedeLineage {
		return supersedeFieldValues(content, cfg.SupersedeFields)
	}
	return nil
}

// supersedeFieldValues sammelt die Werte aller Felder aus fields in der
// Quelldatei. Eine Feld-Zeile hat die Form `**Feld:** Wert` oder
// `Feld: Wert` (Feldname case-insensitiv); gelesen werden nur
// Prosa-Zeilen außerhalb von Fenced-Code (spec/spezifikation.md
// §DC-FA-MTX-001.a Schritt 4).
func supersedeFieldValues(content []byte, fields []string) []string {
	if len(fields) == 0 {
		return nil
	}
	var vals []string
	for _, pl := range proseLines(content) {
		if v, ok := supersedeFieldValue(pl.raw, fields); ok {
			vals = append(vals, v)
		}
	}
	return vals
}

// supersedeFieldValue liefert den Wert, wenn die getrimmte Zeile mit
// einem der Felder als `**Feld:**`- oder `Feld:`-Präfix beginnt
// (case-insensitiv); die fette Form hat Vorrang.
func supersedeFieldValue(raw string, fields []string) (string, bool) {
	t := strings.TrimSpace(raw)
	for _, f := range fields {
		if bold := "**" + f + ":**"; len(t) >= len(bold) && strings.EqualFold(t[:len(bold)], bold) {
			return strings.TrimSpace(t[len(bold):]), true
		}
		if plain := f + ":"; len(t) >= len(plain) && strings.EqualFold(t[:len(plain)], plain) {
			return strings.TrimSpace(t[len(plain):]), true
		}
	}
	return "", false
}

// lineageExempt prüft, ob ein Supersede-Feldwert das Ziel der Referenz
// nennt: normalisierter Teilzeichenketten-Vergleich gegen den Linktext
// (falls nicht leer) bzw. den Zielpfad (rel, Basename, Basename ohne
// Endung). Leere values ⇒ keine Ausnahme (Default byte-identisch).
func lineageExempt(values []string, linkText, rel string) bool {
	if len(values) == 0 {
		return false
	}
	cands := lineageCandidates(linkText, rel)
	for _, v := range values {
		nv := normalizeLineage(v)
		for _, c := range cands {
			if strings.Contains(nv, c) {
				return true
			}
		}
	}
	return false
}

// lineageCandidates liefert die normalisierten Erkennungsformen des
// Referenzziels (Linktext, Pfad, Basename, Basename ohne Endung).
func lineageCandidates(linkText, rel string) []string {
	var out []string
	add := func(s string) {
		if n := normalizeLineage(s); n != "" {
			out = append(out, n)
		}
	}
	add(linkText)
	add(rel)
	base := rel
	if i := strings.LastIndex(base, "/"); i >= 0 {
		base = base[i+1:]
	}
	add(base)
	if i := strings.LastIndex(base, "."); i > 0 {
		add(base[:i])
	}
	return out
}

// normalizeLineage faltet Groß-/Kleinschreibung und kollabiert
// Whitespace auf einzelne Leerzeichen (reihenfolgestabil, DC-QA-02).
func normalizeLineage(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// classOf liefert die erste deklarierte Klasse, deren Glob den Pfad
// matcht (Deklarationsreihenfolge = Präzedenz).
func classOf(classes []model.MatrixClass, rel string) (string, bool) {
	for _, c := range classes {
		for _, g := range c.Paths {
			if matchGlob(g, rel) {
				return c.Name, true
			}
		}
	}
	return "", false
}

// orderedClass liefert die Klasse `name`, falls sie eine aktive
// klasseninterne Richtung trägt (DC-FA-MTX-002: DirectionNoDownward +
// nicht-leeres Order), sonst nil. Die fail-closed-Kopplung von
// order/direction ist im Config-Adapter validiert.
func orderedClass(classes []model.MatrixClass, name string) *model.MatrixClass {
	for i := range classes {
		if c := &classes[i]; c.Name == name &&
			c.Direction == model.DirectionNoDownward && len(c.Order) > 0 {
			return c
		}
	}
	return nil
}

// rankOf liefert den Rang von rel = Index des ersten Order-Globs, den der
// Pfad matcht (First-Match wie classOf); false, wenn keiner matcht
// (rangfrei — nimmt an der Richtungsprüfung nicht teil).
func rankOf(order []string, rel string) (int, bool) {
	for i, g := range order {
		if matchGlob(g, rel) {
			return i, true
		}
	}
	return 0, false
}

// downwardFinding prüft die klasseninterne Verweisrichtung (DC-FA-MTX-002):
// Quelle und Ziel in derselben geordneten Klasse, beide rangbehaftet, Quell-Rang
// kleiner (autoritativer) als Ziel-Rang ⇒ matrix-downward. Sonst ok=false.
func downwardFinding(srcOrdered *model.MatrixClass, srcClass, dstClass, file, rel string,
	line int, target string) (model.Finding, bool) {
	if srcOrdered == nil || dstClass != srcClass {
		return model.Finding{}, false
	}
	si, ok1 := rankOf(srcOrdered.Order, file)
	di, ok2 := rankOf(srcOrdered.Order, rel)
	if !ok1 || !ok2 || si >= di {
		return model.Finding{}, false
	}
	return model.Finding{
		File: file, Line: line, Rule: "matrix",
		Target: target, Reason: ReasonMatrixDownward,
		Message: "Abwärtsverweis innerhalb " + srcClass + ": Rang " +
			strconv.Itoa(si) + " → " + strconv.Itoa(di) + " ist nicht erlaubt",
	}, true
}

// ruleFor liefert die erste deklarierte Regel für das Klassen-Paar.
func ruleFor(rules []model.MatrixRule, from, to string) (model.MatrixRule, bool) {
	for _, r := range rules {
		if r.From == from && r.To == to {
			return r, true
		}
	}
	return model.MatrixRule{}, false
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

// SelectSections schränkt content auf die per include benannten Abschnitte ein
// (leer ⇒ ganze Datei) und entfernt die per exclude benannten (DC-FA-COV-001):
// beide über dieselbe Heading-Span-Semantik wie exclude-sections (Überschrift
// bis zur nächsten gleich-/höherrangigen; Name = voller Heading-Klartext,
// exakter Vergleich via plainHeadingText). So ergibt include
// ["27.1 Anforderung zu Design"] + exclude ["27.1.1 Anforderungen ohne
// Design-Artefakt"] genau „§27.1 ohne §27.1.1". Zeilen-basiert (join mit "\n").
func SelectSections(content []byte, include, exclude []string) []byte {
	if len(include) == 0 && len(exclude) == 0 {
		return content
	}
	inc := excludedRanges(content, include)
	exc := excludedRanges(content, exclude)
	lines := strings.Split(string(content), "\n")
	keep := make([]string, 0, len(lines))
	for i, ln := range lines {
		no := i + 1
		if len(include) > 0 && !inRanges(inc, no) {
			continue
		}
		if inRanges(exc, no) {
			continue
		}
		keep = append(keep, ln)
	}
	return []byte(strings.Join(keep, "\n"))
}

// HeadingTexts liefert die Klartext-Überschriften von content (voller
// Heading-Text ohne Markdown-Auszeichnung, wie der Vergleich in exclude-sections)
// — für den fail-closed Sektionsnamen-Guard des Coverage-Scans (DC-FA-COV-001).
func HeadingTexts(content []byte) []string {
	hs := extractHeadingLines(content)
	out := make([]string, 0, len(hs))
	for _, h := range hs {
		out = append(out, plainHeadingText(h.text))
	}
	return out
}

// plainHeadingText liefert den getrimmten Heading-Text ohne
// Markdown-Auszeichnung (Links → Linktext, Backticks/Sterne entfernt)
// — Vergleichsbasis für exclude-sections (case-sensitiv) und das
// Status-Heading (case-insensitiv).
func plainHeadingText(s string) string {
	t := StripHeadingLinks(s)
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
