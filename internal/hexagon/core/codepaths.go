package core

import (
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonCodepathMissing — Grund-Code des Moduls codepaths
// (spec/spezifikation.md §4).
const ReasonCodepathMissing = "codepath-missing"

// ignoreMarker nimmt eine Zeile von der codepaths-Prüfung aus — und
// nur von dieser: deterministische Befunde anderer Module werden
// behoben, nicht unterdrückt (DC-FA-CODE-001).
const ignoreMarker = "d-check:ignore"

// CodepathsConfig trägt die Präfixe für Wurzel-relative
// Inline-Code-Pfade (spec/spezifikation.md §2; `./`/`../` werden
// immer erkannt).
type CodepathsConfig struct {
	Roots []string
}

// checkCodepaths prüft explizite Pfade in Inline-Code-Spans
// (spec/spezifikation.md §DC-FA-CODE-001.a). Arbeitet auf den rohen
// Prosa-Zeilen, weil die Vorverarbeitung Inline-Code für die übrigen
// Module gerade entfernt; teilt den Slug-Cache mit anchors.
func checkCodepaths(fsys driven.Filesystem, file string, content []byte, cfg CodepathsConfig, slugCache map[string]map[string]bool) []Finding {
	var findings []Finding
	for _, pl := range proseLines(content) {
		if strings.Contains(pl.raw, ignoreMarker) {
			continue
		}
		// Headings sind Titel, keine Prosa-Referenzen — gleiche
		// Ausnahme wie DC-FA-ID-001 (und: ein Marker im Heading
		// würde dessen Anker-Slug verändern).
		if _, _, ok := parseATXHeading(pl.raw); ok {
			continue
		}
		forEachInlineCodeSpan(pl.raw, func(_, _, valStart, valEnd int) {
			value := normalizeCodepath(pl.raw[valStart:valEnd])
			rootRel, ok := classifyCodepath(value, cfg.Roots)
			if !ok {
				return
			}
			findings = append(findings,
				checkCodepathTarget(fsys, file, pl.no, value, rootRel, slugCache)...)
		})
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
	return finding(ReasonAnchorMissing, "Anker entspricht keinem Heading-Slug der Zieldatei")
}

// codepathSlugs liest die Slug-Menge der Zieldatei über denselben
// Cache wie das Modul anchors (nil = nicht lesbar).
func codepathSlugs(fsys driven.Filesystem, cache map[string]map[string]bool, rel string) map[string]bool {
	if s, ok := cache[rel]; ok {
		return s
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		cache[rel] = nil
		return nil
	}
	s := HeadingSlugs(content)
	cache[rel] = s
	return s
}
