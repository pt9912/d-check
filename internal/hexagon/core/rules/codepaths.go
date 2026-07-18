package rules

import (
	"fmt"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonCodepathMissing — Grund-Code des Moduls codepaths
// (spec/spezifikation.md §4).
const ReasonCodepathMissing = "codepath-missing"

// ignoreMarker nimmt eine Zeile von der codepaths- und der
// ids-`always`-Prüfung aus — nur von diesen beiden: deterministische
// Befunde anderer Module werden behoben, nicht unterdrückt
// (DC-FA-CODE-001; Geltungsbereich auf ids erweitert mit
// DC-FA-ID-001 0.8.0, slice-018 — illustrative Beispiel-IDs).
const ignoreMarker = "d-check:ignore"

// CheckCodepaths prüft explizite Pfade in Inline-Code-Spans
// (spec/spezifikation.md §DC-FA-CODE-001.a). Arbeitet auf den rohen
// Prosa-Zeilen, weil die Vorverarbeitung Inline-Code für die übrigen
// Module gerade entfernt; teilt den Slug-Cache mit anchors.
func CheckCodepaths(fsys driven.Filesystem, file string, content []byte, cfg model.CodepathsConfig, slugCache map[string]map[string]bool, ignoreRefs []model.IgnoreRef) []model.Finding {
	var findings []model.Finding
	// exempt-paths: ganze Datei von der codepaths-Prüfung ausnehmen
	// (datei-weit, Glob wie scan.ignore; Vorbild ids-Ventil). Wirkt
	// unabhängig von cfg.Roots (§DC-FA-CODE-001.a).
	if ignored(file, cfg.ExemptPaths) {
		return nil
	}
	// Geteiltes Referenz-Ventil: Top-Level ignore-refs (DC-FA-REF-001)
	// plus der modul-lokale Alias codepaths.ignore-refs (ohne in/keep,
	// skopiert auf codepaths — byte-identisch zur bisherigen Fassung).
	refs := ignoreRefs
	if len(cfg.IgnoreRefs) > 0 {
		refs = append(append([]model.IgnoreRef(nil), ignoreRefs...), model.IgnoreRef{Refs: cfg.IgnoreRefs})
	}
	prose := proseLines(content)
	spans := inlineSpansByLine(prose)
	for _, pl := range prose {
		if strings.Contains(pl.raw, ignoreMarker) {
			continue
		}
		// Headings sind Titel, keine Prosa-Referenzen — gleiche
		// Ausnahme wie DC-FA-ID-001 (und: ein Marker im Heading
		// würde dessen Anker-Slug verändern).
		if _, _, ok := parseATXHeading(pl.raw); ok {
			continue
		}
		for _, sp := range spans[pl.no] {
			// Span als Linktext ([`…`](ziel)): das Ziel prüft das
			// Modul links — der Text ist Beschriftung, kein Pfad-
			// Anspruch (§DC-FA-CODE-001.a Schritt 2).
			if sp.start > 0 && pl.raw[sp.start-1] == '[' && sp.end < len(pl.raw) && pl.raw[sp.end] == ']' {
				continue
			}
			value, from, to, hasRange := codepathValueAndRange(pl.raw[sp.valStart:sp.valEnd], cfg.CheckLines)
			rootRel, ok := classifyCodepath(value, cfg.Roots)
			if !ok {
				continue
			}
			findings = append(findings,
				checkCodepathTarget(fsys, file, pl.no, value, rootRel, hasRange, from, to, refs, slugCache)...)
		}
	}
	return findings
}

// normalizeCodepath: Whitespace trimmen, umschließende
// Anführungszeichen und schließende Satzzeichen entfernen —
// iterativ bis stabil, damit auch `"./a.md",` sauber wird
// (§DC-FA-CODE-001.a Schritt 3).
func normalizeCodepath(v string) string {
	for {
		w := strings.TrimSpace(v)
		w = trimLineSuffix(w)
		w = strings.TrimRight(w, ".,;:")
		w = strings.Trim(w, `"'`)
		if w == v {
			return v
		}
		v = w
	}
}

// codepathValueAndRange normalisiert den Span-Wert und trennt bei aktivem
// check-lines eine Zeilen-Referenz (`:<von>`/`:<von>-<bis>`) ab, die für
// den Zeilen-Check (Schritt 6) gemerkt wird. Ohne check-lines oder ohne
// Bereich bleibt value die reine normalizeCodepath-Ausgabe (byte-identisch;
// §DC-FA-CODE-001.a Schritt 3).
func codepathValueAndRange(raw string, checkLines bool) (value string, from, to int, hasRange bool) {
	if checkLines {
		if p, f, t, has := splitCitationRange(normalizeCodepathBase(raw)); has {
			return p, f, t, true
		}
	}
	return normalizeCodepath(raw), 0, 0, false
}

// trimLineSuffix trennt ein Zeilen-Suffix `:NNN` ab — die etablierte
// Datei:Zeile-Konvention (z. B. `spec/lastenheft.md:290`) referenziert
// die Datei, nicht einen Pfad mit Doppelpunkt.
func trimLineSuffix(v string) string {
	i := strings.LastIndexByte(v, ':')
	if i <= 0 || i == len(v)-1 {
		return v
	}
	for _, r := range v[i+1:] {
		if r < '0' || r > '9' {
			return v
		}
	}
	return v[:i]
}

// normalizeCodepathBase entfernt Whitespace, umschließende
// Anführungszeichen und schließende Satzzeichen iterativ — wie
// normalizeCodepath, aber OHNE das Zeilen-Suffix abzutrennen (das löst
// splitCitationRange bei aktivem check-lines separat, um den Bereich zu
// erhalten; §DC-FA-CODE-001.a Schritt 3).
func normalizeCodepathBase(v string) string {
	for {
		w := strings.TrimSpace(v)
		w = strings.TrimRight(w, ".,;:")
		w = strings.Trim(w, `"'`)
		if w == v {
			return v
		}
		v = w
	}
}

// splitCitationRange trennt ein Zeilen-Suffix `:<von>` oder
// `:<von>-<bis>` (beide Teile rein numerisch) vom Pfad ab. base = Pfad
// ohne Suffix; ohne `-` ist to == from; has=false, wenn kein solches
// Suffix vorliegt (§DC-FA-CODE-001.a Schritt 3).
func splitCitationRange(v string) (base string, from, to int, has bool) {
	i := strings.LastIndexByte(v, ':')
	if i <= 0 || i == len(v)-1 {
		return v, 0, 0, false
	}
	suffix := v[i+1:]
	fromStr, toStr := suffix, suffix
	if d := strings.IndexByte(suffix, '-'); d >= 0 {
		fromStr, toStr = suffix[:d], suffix[d+1:]
	}
	f, ok1 := allDigits(fromStr)
	t, ok2 := allDigits(toStr)
	if !ok1 || !ok2 {
		return v, 0, 0, false
	}
	return v[:i], f, t, true
}

// allDigits parst einen nicht-leeren rein numerischen String; ok=false
// sonst (leer oder mit Nicht-Ziffern).
func allDigits(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
	}
	return n, true
}

// countLines zählt die Zeilen eines Datei-Inhalts; ein finaler
// Zeilenumbruch zählt nicht als eigene Zeile (geteilt vom codepaths-
// Zeilen-Check und citations).
func countLines(content []byte) int {
	s := string(content)
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

// classifyCodepath entscheidet konservativ, ob value ein expliziter
// Pfad ist (§DC-FA-CODE-001.a Schritt 4); rootRel meldet einen
// Präfix-Treffer (Auflösung relativ zur Repo-Wurzel).
func classifyCodepath(v string, roots []string) (rootRel, ok bool) {
	if v == "" || strings.ContainsAny(v, " \t{}<>|*?=") {
		return false, false
	}
	if strings.Contains(v, "…") || strings.Contains(v, "->") || strings.Contains(v, "→") {
		return false, false
	}
	if strings.HasPrefix(v, "//") || strings.HasPrefix(v, "#") || IsExternalScheme(v) {
		return false, false
	}
	if strings.HasPrefix(v, "./") || strings.HasPrefix(v, "../") {
		return false, true
	}
	for _, r := range roots {
		if strings.HasPrefix(v, strings.TrimSuffix(r, "/")+"/") {
			return true, true
		}
	}
	return false, false
}

// checkCodepathTarget löst auf (Datei- bzw. Wurzel-relativ), prüft
// Escape und Existenz; bei Markdown-Zielen mit Fragment zusätzlich
// den Anker (§DC-FA-CODE-001.a Schritt 5). Ein aufgelöster Pfad, der
// ein ignore-refs-Glob matcht, wird vor allen Prüfungen übersprungen
// (referenz-weites Tombstone-Ventil, ADR-0025).
func checkCodepathTarget(fsys driven.Filesystem, file string, line int, value string, rootRel, hasRange bool, from, to int, refs []model.IgnoreRef, slugCache map[string]map[string]bool) []model.Finding {
	pathPart, frag := value, ""
	if idx := strings.IndexByte(value, '#'); idx != -1 {
		pathPart, frag = value[:idx], value[idx+1:]
	}
	finding := func(reason, message string) []model.Finding {
		return []model.Finding{{
			File: file, Line: line, Rule: "codepaths",
			Target: value, Reason: reason, Message: message,
		}}
	}
	var rel string
	var escaped bool
	if rootRel {
		rel, escaped = ResolveConfigPath(pathPart)
	} else {
		var ok bool
		rel, escaped, ok = ResolveTarget(file, pathPart)
		if !ok {
			return nil
		}
	}
	if refIgnored(refs, file, rel) {
		return nil
	}
	if escaped {
		return finding(model.ReasonRepoEscape, "aufgelöstes Ziel verlässt die Repository-Wurzel")
	}
	kind, err := fsys.Kind(rel)
	if err != nil || kind == driven.KindMissing {
		return finding(ReasonCodepathMissing, "Ziel des Inline-Code-Pfads existiert nicht")
	}
	// Zeilen-Check (§DC-FA-CODE-001.a Schritt 6, nur bei check-lines):
	// das Ziel existiert bereits — jetzt den gemerkten Bereich prüfen.
	if hasRange {
		if fs := checkCodepathLineRange(fsys, file, line, value, rel, from, to); fs != nil {
			return fs
		}
	}
	if frag == "" || !strings.HasSuffix(rel, ".md") {
		return nil
	}
	slugs := codepathSlugs(fsys, slugCache, rel)
	if slugs == nil || slugs[frag] {
		return nil // nicht lesbar → schweigen; Anker ok → kein Befund
	}
	return finding(ReasonAnchorMissing, "Anker entspricht keinem Heading-Slug und keinem HTML-Anker der Zieldatei")
}

// checkCodepathLineRange prüft den gemerkten Zeilen-Bereich gegen das
// bereits als existierend bestätigte Ziel (§DC-FA-CODE-001.a Schritt 6):
// `von > bis` ⇒ citation-inverted-range; hat die Datei weniger als `bis`
// Zeilen ⇒ citation-out-of-range. Der Befund-target trägt den Bereich.
func checkCodepathLineRange(fsys driven.Filesystem, file string, line int, value, rel string, from, to int) []model.Finding {
	target := fmt.Sprintf("%s:%d-%d", value, from, to)
	mk := func(reason, msg string) []model.Finding {
		return []model.Finding{{
			File: file, Line: line, Rule: "codepaths",
			Target: target, Reason: reason, Message: msg,
		}}
	}
	if from < 1 || from > to {
		return mk(ReasonCitationInvertedRange, "Zeilen-Referenz ungültig (von < 1 oder von > bis)")
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		return mk(ReasonCitationOutOfRange, "Zeilen-Referenz hinter dem Datei-Ende")
	}
	if to > countLines(content) {
		return mk(ReasonCitationOutOfRange, "Zeilen-Referenz hinter dem Datei-Ende")
	}
	return nil
}

// codepathSlugs liest die gültige Anker-Menge der Zieldatei über
// denselben Cache wie das Modul anchors (nil = nicht lesbar).
func codepathSlugs(fsys driven.Filesystem, cache map[string]map[string]bool, rel string) map[string]bool {
	if s, ok := cache[rel]; ok {
		return s
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		cache[rel] = nil
		return nil
	}
	s := AnchorSet(content)
	cache[rel] = s
	return s
}
