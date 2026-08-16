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

// CheckVersions ist das Regelmodul `versions` (DC-FA-VER-001): jeder
// Versions-Pin (cfg.PinPattern) muss die aktuelle Version (current) tragen,
// sonst Befund version-stale. Anders als die übrigen Module liest es die Pins
// AUCH innerhalb von Fenced-Code (Pins leben in Kommando-Beispielen) — eine
// gescopte Ausnahme von der Fence-Opazität, reiner Muster-Scan über die
// Rohzeilen (kein Parser). Ventile wie ids/codepaths: exempt-paths (datei-weit)
// und der Zeilen-Marker d-check:ignore; die current-from-Datei (fromFile) ist
// selbst ausgenommen (ihr Verlauf listet bewusst alle Versionen). Ohne
// PinPattern wirkungslos (byte-identisch, DC-QA-02).
func CheckVersions(file string, content []byte, cfg model.VersionsConfig, current, fromFile string) []model.Finding {
	if cfg.PinPattern == nil {
		return nil
	}
	if file == fromFile || ignored(file, cfg.ExemptPaths) {
		return nil
	}
	var findings []model.Finding
	for i, raw := range strings.Split(string(content), "\n") {
		if strings.Contains(raw, ignoreMarker) {
			continue // d-check:ignore — Zeile von der versions-Prüfung frei
		}
		for _, m := range cfg.PinPattern.FindAllStringSubmatchIndex(raw, -1) {
			ver := pinVersion(raw, m)
			if ver == current {
				continue
			}
			findings = append(findings, model.Finding{
				File: file, Line: i + 1, Rule: "versions",
				Target: ver, Reason: model.ReasonVersionStale,
				Message: fmt.Sprintf("Versions-Pin trägt %s, erwartet %s (versions.current-from)", ver, current),
			})
		}
	}
	return findings
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

// resolveCurrentVersion liest die aktuelle Version aus der current-from-Quelle
// (spec/spezifikation.md §DC-FA-VER-001.a Schritt 1): Datei links vom '#',
// Anker rechts; der Span ist die Heading-Section des Ankers bzw. die ganze
// Datei ohne Anker. Fehlt die Datei/der Anker oder trägt der Span keine
// Version → Fehler (Exit 2, fail-closed). Liefert zusätzlich den aufgelösten
// Datei-Pfad (Selbst-Ausnahme der current-from-Datei).
func resolveCurrentVersion(fsys driven.Filesystem, currentFrom string) (version, fromFile string, err error) {
	filePart, anchor := currentFrom, ""
	if i := strings.IndexByte(currentFrom, '#'); i != -1 {
		filePart, anchor = currentFrom[:i], currentFrom[i+1:]
	}
	rel, escaped := ResolveConfigPath(filePart)
	if escaped {
		return "", "", fmt.Errorf("versions.current-from verlässt die Repository-Wurzel: %s", currentFrom)
	}
	if rel == "" {
		return "", "", fmt.Errorf("versions.current-from muss eine Datei benennen: %s", currentFrom)
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		return "", "", fmt.Errorf("versions.current-from nicht lesbar (%s): %w", currentFrom, err)
	}
	span := string(content)
	if anchor != "" {
		s, ok := headingSection(content, anchor)
		if !ok {
			return "", "", fmt.Errorf("versions.current-from: Anker #%s nicht auflösbar in %s", anchor, filePart)
		}
		span = s
	}
	ver := versionRE.FindString(span)
	if ver == "" {
		return "", "", fmt.Errorf("versions.current-from: keine Version im adressierten Span von %s", currentFrom)
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
	for idx, h := range headings {
		if Slugify(h.text) != anchor {
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
// (id=/name=) bis zur nächsten Überschrift. Gezählt wird der Anker nur
// außerhalb von Fences (DC-FA-ANCH-001.b); der Span selbst bleibt roh.
func htmlAnchorSection(content []byte, lines []string, headings []headingLine, anchor string) (string, bool) {
	htmlRE := regexp.MustCompile(`(?i)(?:id|name)\s*=\s*["']` + regexp.QuoteMeta(anchor) + `["']`)
	prose := make(map[int]bool)
	for _, pl := range proseLines(content) {
		prose[pl.no] = true
	}
	for i, raw := range lines {
		if !prose[i+1] || !htmlRE.MatchString(raw) {
			continue
		}
		start := i + 1
		end := len(lines)
		for _, h2 := range headings {
			if h2.line > start {
				end = h2.line - 1
				break
			}
		}
		return strings.Join(lines[start-1:end], "\n"), true
	}
	return "", false
}
