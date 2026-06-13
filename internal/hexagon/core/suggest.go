package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// idShape ist die allgemeine Kennungs-Gestalt für die Extraktion aus
// Headings (DC-FA-CLI-006.a Schritt 2): Präfix aus Großbuchstaben-Segmenten,
// dann -NNN, optional ein Suffix-Buchstabe.
var idShape = regexp.MustCompile(`^([A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*)-(\d+)([A-Za-z]?)$`)

// suggestedPattern ist ein aus einer Autoritäts-Quelle abgeleitetes
// ids-Muster samt der Quell-Kennungen (für den Nachweis-Kommentar).
type suggestedPattern struct {
	target string
	regex  string
	ids    []string
}

// SuggestConfig liest die benannten Autoritäts-Quellen, leitet je Quelle
// ein ids-Muster aus den dort definierten Kennungen ab und liefert ein
// kommentiertes .d-check.yml-Gerüst (DC-FA-CLI-006.a). Reiner Lese-Pfad —
// es wird nie geschrieben (DC-QA-03). Eine fehlende oder die Repo-Wurzel
// verlassende Quelle ist ein Fehler (CLI: Exit 2).
func SuggestConfig(fsys driven.Filesystem, sources []string) (string, error) {
	var patterns []suggestedPattern
	for _, src := range sources {
		rel, escaped := resolveConfigPath(src)
		if escaped {
			return "", fmt.Errorf("Autoritäts-Quelle verlässt die Repository-Wurzel: %s", src)
		}
		kind, err := fsys.Kind(rel)
		if err != nil || kind == driven.KindMissing {
			return "", fmt.Errorf("Autoritäts-Quelle existiert nicht: %s", src)
		}
		ids, err := extractDefinedIDs(fsys, rel, kind)
		if err != nil {
			return "", err
		}
		patterns = append(patterns, suggestedPattern{target: src, regex: deriveRegex(ids), ids: ids})
	}
	return renderSuggestion(patterns, probeOptInModules(fsys)), nil
}

// extractDefinedIDs sammelt die in den Headings einer Quelle (Datei oder
// Verzeichnis) definierten Kennungen — das führende Token, das idShape
// erfüllt (kein Fließtext-Mining). Ergebnis bytewise sortiert.
func extractDefinedIDs(fsys driven.Filesystem, rel string, kind driven.EntryKind) ([]string, error) {
	files := []string{rel}
	if kind == driven.KindDir {
		discovered, err := DiscoverFiles(fsys, []string{rel}, nil)
		if err != nil {
			return nil, err
		}
		files = discovered
	}
	seen := map[string]bool{}
	for _, f := range files {
		content, err := fsys.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, h := range ExtractHeadings(content) {
			fields := strings.Fields(h)
			if len(fields) > 0 && idShape.MatchString(fields[0]) {
				seen[fields[0]] = true
			}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

// deriveRegex bildet aus den Kennungen einer Quelle die Präfix-Alternation
// (?:p₁|p₂|…)-\d+[A-Za-z]? (DC-FA-CLI-006.a Schritt 3). Round-Trip-Invariante:
// das Ergebnis matcht jede Eingabe-Kennung. Leer, wenn keine Kennung vorliegt.
func deriveRegex(ids []string) string {
	prefixes := map[string]bool{}
	hasLetter := false
	for _, id := range ids {
		if m := idShape.FindStringSubmatch(id); m != nil {
			prefixes[m[1]] = true
			hasLetter = hasLetter || m[3] != ""
		}
	}
	if len(prefixes) == 0 {
		return ""
	}
	ps := make([]string, 0, len(prefixes))
	for p := range prefixes {
		ps = append(ps, regexp.QuoteMeta(p))
	}
	sort.Strings(ps)
	body := ps[0]
	if len(ps) > 1 {
		body = "(?:" + strings.Join(ps, "|") + ")"
	}
	body += `-\d+`
	if hasLetter {
		body += `[A-Za-z]?`
	}
	return body
}

// probeOptInModules lässt die opt-in-Module probeweise laufen und liefert
// (in fester Reihenfolge) jene mit ≥1 Befund (DC-FA-CLI-006.a Schritt 4).
func probeOptInModules(fsys driven.Filesystem) []string {
	optIn := []string{"codepaths", "spans", "hostpaths"}
	res, err := Run(fsys, nil, Config{Modules: optIn}, optIn)
	if err != nil {
		return nil
	}
	active := map[string]bool{}
	for _, f := range res.Findings {
		active[f.Rule] = true
	}
	var out []string
	for _, m := range optIn {
		if active[m] {
			out = append(out, m)
		}
	}
	return out
}

// renderSuggestion baut das kommentierte, dekodierbare Gerüst.
func renderSuggestion(patterns []suggestedPattern, probed []string) string {
	var b strings.Builder
	b.WriteString("# .d-check.yml — Vorschlag aus `d-check --suggest-config` (advisory).\n")
	b.WriteString("# Prüfen und verengen: die ids-Muster sind eine Best-Guess-Ableitung\n")
	b.WriteString("# aus den definierten Kennungen der benannten Quellen.\n\n")
	b.WriteString("scan:\n  roots: [\".\"]\n\n")

	modules := append([]string{"links", "anchors"}, probed...)
	b.WriteString("modules: [" + strings.Join(modules, ", ") + "]\n")

	var withIDs []suggestedPattern
	for _, p := range patterns {
		if p.regex != "" {
			withIDs = append(withIDs, p)
		} else {
			fmt.Fprintf(&b, "# Hinweis: in %q keine definierten Kennungen gefunden.\n", p.target)
		}
	}
	if len(withIDs) == 0 {
		return b.String()
	}
	b.WriteString("\nids:\n  patterns:\n")
	for _, p := range withIDs {
		fmt.Fprintf(&b, "    # abgeleitet aus %d Kennung(en) in %s: %s\n",
			len(p.ids), p.target, strings.Join(p.ids, ", "))
		fmt.Fprintf(&b, "    - regex: '%s'\n      target: %s\n", p.regex, p.target)
		b.WriteString("      # link-policy: always   # einkommentieren für strenge Linkpflicht\n")
	}
	return b.String()
}
