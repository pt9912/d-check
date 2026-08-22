package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// versionRE erkennt einen Versionsstring (optionales führendes v) — daraus
// wird die erwartete Version aus dem current-from-Span gezogen
// (spec/spezifikation.md §DC-FA-VER-001.a Schritt 1).
var versionRE = regexp.MustCompile(`v?\d+\.\d+\.\d+`)

// ResolvedVersionPattern ist ein Muster-Quellen-Paar samt der am Lauf-Start
// aufgelösten erwarteten Version und dem Pfad ihrer Quell-Datei
// (spec/spezifikation.md §DC-FA-VER-001.a). Index ist die Position in der
// Konfiguration — sie benennt das Paar in Meldungen, sobald es mehr als eines
// gibt.
type ResolvedVersionPattern struct {
	Pattern  model.VersionPattern
	Current  string
	FromFile string
	Index    int
}

// CheckVersions ist das Regelmodul `versions` (DC-FA-VER-001): jeder
// Versions-Pin muss die aktuelle Version seines Paares tragen, sonst Befund
// version-stale. Anders als die übrigen Module liest es die Pins AUCH
// innerhalb von Fenced-Code (Pins leben in Kommando-Beispielen) — eine
// gescopte Ausnahme von der Fence-Opazität, reiner Muster-Scan über die
// Rohzeilen (kein Parser).
//
// Je Zeile entsteht HÖCHSTENS EIN Befund pro Pin-Wert, weil die Befund-Adresse
// (Datei, Zeile, Regel, Ziel, Grund) zwei Befunde an derselben Stelle nicht
// unterscheiden kann — die geteilte Nachrunde würde den zweiten verwerfen und
// mit ihm seine Erwartung. Treffen mehrere Paare denselben Wert, nennt die
// Nachricht deshalb ALLE Erwartungen in Deklarationsreihenfolge.
//
// Die beiden DATEI-Ventile sind paar-lokal — exempt-paths gilt für sein Paar,
// die Selbst-Ausnahme der Quell-Datei nur für deren eigene Reihe; der
// ZEILEN-Marker d-check:ignore nimmt die Zeile allen Paaren aus. Ohne Paar
// wirkungslos (byte-identisch, DC-QA-02).
func CheckVersions(file string, content []byte, patterns []ResolvedVersionPattern) []model.Finding {
	active := applicableVersionPatterns(file, patterns)
	if len(active) == 0 {
		return nil
	}
	var findings []model.Finding
	for i, raw := range strings.Split(string(content), "\n") {
		if strings.Contains(raw, ignoreMarker) {
			continue // d-check:ignore — Zeile von der versions-Prüfung frei
		}
		for _, st := range staleVersionsOfLine(raw, active) {
			findings = append(findings, model.Finding{
				File: file, Line: i + 1, Rule: "versions",
				Target: st.target, Reason: model.ReasonVersionStale,
				Message: staleVersionMessage(st, len(patterns) > 1),
			})
		}
	}
	return findings
}

// staleVersion ist ein veralteter Pin-Wert einer Zeile samt den Erwartungen
// aller Paare, die ihn melden — in Deklarationsreihenfolge.
type staleVersion struct {
	target string
	expect []versionExpectation
}

// versionExpectation ist die erwartete Version eines Paares samt dessen
// Deklarations-Index.
type versionExpectation struct {
	index int
	want  string
}

// staleVersionsOfLine gruppiert die veralteten Pin-Werte einer Zeile nach dem
// gefundenen Wert; die Gruppen stehen in der Reihenfolge ihres ersten
// Auftretens, die Erwartungen darin in Deklarationsreihenfolge der Paare.
// Dasselbe Paar trägt zu einem Wert höchstens eine Erwartung bei, auch wenn es
// ihn mehrfach auf der Zeile trifft.
func staleVersionsOfLine(raw string, active []ResolvedVersionPattern) []staleVersion {
	var groups []staleVersion
	at := make(map[string]int)
	for _, p := range active {
		for _, m := range p.Pattern.PinPattern.FindAllStringSubmatchIndex(raw, -1) {
			ver := pinVersion(raw, m)
			if ver == p.Current {
				continue
			}
			idx, ok := at[ver]
			if !ok {
				at[ver] = len(groups)
				groups = append(groups, staleVersion{target: ver})
				idx = len(groups) - 1
			}
			g := &groups[idx]
			if len(g.expect) > 0 && g.expect[len(g.expect)-1].index == p.Index {
				continue // dasselbe Paar, derselbe Wert, zweiter Treffer
			}
			g.expect = append(g.expect, versionExpectation{index: p.Index, want: p.Current})
		}
	}
	return groups
}

// staleVersionMessage nennt den gefundenen Wert und jede Erwartung mit ihrer
// Quelle. Die Fundstelle wird nur benannt, wenn es mehr als ein Paar gibt —
// eine Ein-Paar-Konfiguration (und damit die Kurzform) behält ihren Wortlaut
// byte-identisch (DC-QA-02).
func staleVersionMessage(st staleVersion, named bool) string {
	parts := make([]string, 0, len(st.expect))
	for _, e := range st.expect {
		key := "versions.current-from"
		if named {
			key = fmt.Sprintf("versions.patterns[%d].current-from", e.index)
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", e.want, key))
	}
	return fmt.Sprintf("Versions-Pin trägt %s, erwartet %s", st.target, strings.Join(parts, " sowie "))
}

// applicableVersionPatterns liefert die Paare, deren Datei-Ventile diese Datei
// nicht ausnehmen: die eigene Quell-Datei eines Paares und dessen exempt-paths
// schalten nur dieses Paar ab, nicht die übrigen
// (spec/spezifikation.md §DC-FA-VER-001.a Schritt 4).
func applicableVersionPatterns(file string, patterns []ResolvedVersionPattern) []ResolvedVersionPattern {
	active := make([]ResolvedVersionPattern, 0, len(patterns))
	for _, p := range patterns {
		if file == p.FromFile || ignored(file, p.Pattern.ExemptPaths) {
			continue
		}
		active = append(active, p)
	}
	return active
}

// pinVersion liefert die Version eines Pin-Treffers: Capture-Gruppe 1, falls
// das Muster eine hat und sie matchte, sonst der ganze Treffer
// (spec/spezifikation.md §DC-FA-VER-001.a Schritt 2).
func pinVersion(raw string, m []int) string {
	if len(m) >= 4 && m[2] >= 0 {
		return raw[m[2]:m[3]]
	}
	return raw[m[0]:m[1]]
}

// versionSource benennt die Fundstelle eines Paares in Meldungen: bei genau
// einem Paar den Kurzform-Schlüssel, sonst das Paar mit seinem Index — eine
// Ein-Paar-Konfiguration behält ihren Wortlaut byte-identisch (DC-QA-02).
func versionSource(cfg model.VersionsConfig, index int) string {
	if len(cfg.Patterns) < 2 {
		return "versions.current-from"
	}
	return fmt.Sprintf("versions.patterns[%d].current-from", index)
}

// resolveCurrentVersion liest die aktuelle Version aus der current-from-Quelle
// (spec/spezifikation.md §DC-FA-VER-001.a Schritt 1): Datei links vom '#',
// Anker rechts; der Span ist die Heading-Section des Ankers bzw. die ganze
// Datei ohne Anker. Fehlt die Datei/der Anker oder trägt der Span keine
// Version → Fehler (Exit 2, fail-closed). source benennt die Fundstelle in der
// Konfiguration. Liefert zusätzlich den aufgelösten Datei-Pfad
// (Selbst-Ausnahme der current-from-Datei).
func resolveCurrentVersion(fsys driven.Filesystem, currentFrom, source string) (version, fromFile string, err error) {
	filePart, anchor := currentFrom, ""
	if i := strings.IndexByte(currentFrom, '#'); i != -1 {
		filePart, anchor = currentFrom[:i], DecodeFragment(currentFrom[i+1:])
	}
	rel, escaped := ResolveConfigPath(filePart)
	if escaped {
		return "", "", fmt.Errorf("%s verlässt die Repository-Wurzel: %s", source, currentFrom)
	}
	if rel == "" {
		return "", "", fmt.Errorf("%s muss eine Datei benennen: %s", source, currentFrom)
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		return "", "", fmt.Errorf("%s nicht lesbar (%s): %w", source, currentFrom, err)
	}
	span := string(content)
	if anchor != "" {
		s, ok := headingSection(content, anchor)
		if !ok {
			return "", "", fmt.Errorf("%s: Anker #%s nicht auflösbar in %s", source, anchor, filePart)
		}
		span = s
	}
	ver := versionRE.FindString(span)
	if ver == "" {
		return "", "", fmt.Errorf("%s: keine Version im adressierten Span von %s", source, currentFrom)
	}
	return ver, rel, nil
}

// headingSection liefert den Text der durch anchor adressierten Sektion: die
// Heading-Section (Heading-Zeile bis zur nächsten gleich-/höherrangigen
// Überschrift) eines Headings, dessen GitHub-Slug == anchor; sonst ab der
// Zeile eines Inline-HTML-Ankers (id=/name=) bis zur nächsten Überschrift.
//
// Beide Anker-Formen zählen nur außerhalb von Fences (DC-FA-ANCH-001.b): der
// Heading-Zweig über den Prosa-Automaten, der HTML-Zweig über dieselbe
// Zeilen-Menge. Der zurückgegebene SPAN bleibt roh, einschließlich Fenced-Code
// (ADR-0019/ADR-0020).
func headingSection(content []byte, anchor string) (string, bool) {
	lines := strings.Split(string(content), "\n")
	headings := extractHeadingLines(content)
	if span, ok := slugSection(lines, headings, anchor); ok {
		return span, true
	}
	return htmlAnchorSection(content, lines, headings, anchor)
}

// slugSection liefert die Heading-Section des Headings, dessen GitHub-Slug auf
// anchor passt: von seiner Zeile bis zur nächsten gleich-/höherrangigen
// Überschrift.
func slugSection(lines []string, headings []headingLine, anchor string) (string, bool) {
	texts := make([]string, len(headings))
	for i, h := range headings {
		texts[i] = h.text
	}
	slugs := headingSlugsOrdered(texts)
	for idx, h := range headings {
		if slugs[idx] != anchor {
			continue
		}
		end := len(lines)
		for _, h2 := range headings[idx+1:] {
			if h2.level <= h.level {
				end = h2.line - 1
				break
			}
		}
		return strings.Join(lines[h.line-1:end], "\n"), true
	}
	return "", false
}

// htmlAnchorSection liefert den Span ab der Zeile eines Inline-HTML-Ankers
// bis zur nächsten Überschrift. Welche Zeichenfolge ein Anker IST, entscheidet
// die geteilte Erkennung (htmlAnchorLines, DC-FA-ANCH-001.b) — nicht eine
// eigene Regex: sonst gäbe derselbe Lauf zwei Antworten auf dieselbe Frage.
// Der zurückgegebene Span bleibt roh, einschließlich Fenced-Code.
func htmlAnchorSection(content []byte, lines []string, headings []headingLine, anchor string) (string, bool) {
	start, ok := htmlAnchorLines(content)[anchor]
	if !ok {
		return "", false
	}
	end := len(lines)
	for _, h2 := range headings {
		if h2.line > start {
			end = h2.line - 1
			break
		}
	}
	return strings.Join(lines[start-1:end], "\n"), true
}
