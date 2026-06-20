package core

import (
	"strings"

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

// checkCodepaths prüft explizite Pfade in Inline-Code-Spans
// (spec/spezifikation.md §DC-FA-CODE-001.a). Arbeitet auf den rohen
// Prosa-Zeilen, weil die Vorverarbeitung Inline-Code für die übrigen
// Module gerade entfernt; teilt den Slug-Cache mit anchors.
func checkCodepaths(fsys driven.Filesystem, file string, content []byte, cfg CodepathsConfig, slugCache map[string]map[string]bool) []Finding {
	var findings []Finding
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
			value := normalizeCodepath(pl.raw[sp.valStart:sp.valEnd])
			rootRel, ok := classifyCodepath(value, cfg.Roots)
			if !ok {
				continue
			}
			findings = append(findings,
				checkCodepathTarget(fsys, file, pl.no, value, rootRel, slugCache)...)
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
// den Anker (§DC-FA-CODE-001.a Schritt 5).
func checkCodepathTarget(fsys driven.Filesystem, file string, line int, value string, rootRel bool, slugCache map[string]map[string]bool) []Finding {
	pathPart, frag := value, ""
	if idx := strings.IndexByte(value, '#'); idx != -1 {
		pathPart, frag = value[:idx], value[idx+1:]
	}
	finding := func(reason, message string) []Finding {
		return []Finding{{
			File: file, Line: line, Rule: "codepaths",
			Target: value, Reason: reason, Message: message,
		}}
	}
	var rel string
	var escaped bool
	if rootRel {
		rel, escaped = resolveConfigPath(pathPart)
	} else {
		var ok bool
		rel, escaped, ok = ResolveTarget(file, pathPart)
		if !ok {
			return nil
		}
	}
	if escaped {
		return finding(ReasonRepoEscape, "aufgelöstes Ziel verlässt die Repository-Wurzel")
	}
	kind, err := fsys.Kind(rel)
	if err != nil || kind == driven.KindMissing {
		return finding(ReasonCodepathMissing, "Ziel des Inline-Code-Pfads existiert nicht")
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
