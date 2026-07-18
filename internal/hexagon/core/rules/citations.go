package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes der Zitat-Verifikation (spec/spezifikation.md §4).
// citation-out-of-range und citation-inverted-range teilt citations mit
// dem codepaths-Zeilen-Check (DC-FA-CODE-001.a Schritt 6); citation-mismatch
// ist citations-eigen (DC-FA-CITE-001.a Schritt 5).
const (
	ReasonCitationOutOfRange    = "citation-out-of-range"
	ReasonCitationInvertedRange = "citation-inverted-range"
	ReasonCitationMismatch      = "citation-mismatch"
)

// citationMinLen ist die Mindestlänge (in Runen) des normalisierten
// Zitattexts, unter der nicht verglichen wird — ein sehr kurzer Teilstring
// träfe zufällig (schwache Diskriminierung; DC-FA-CITE-001.a Schritt 4,
// dokumentierter Trade-off, keine Falsch-Rot-Gefahr).
const citationMinLen = 16

// citeMarkerRe erkennt die Kommentar-Form der Direktive (`<!-- d-check:cite`)
// — nur so wird eine Zeile als Direktive behandelt, nicht bei bloßer
// Erwähnung des Strings in Prosa.
var citeMarkerRe = regexp.MustCompile(`<!--\s*d-check:cite\b`)

// citeDirectiveRe parst die volle Direktive: Pfad (nicht-gierig bis zum
// Zeilen-Suffix), 1-basiertes <von> und optionales <bis>.
var citeDirectiveRe = regexp.MustCompile(`<!--\s*d-check:cite\s+(.+?):(\d+)(?:-(\d+))?\s*-->`)

// wsRun kollabiert jeden Whitespace-Lauf (inkl. Zeilenumbruch) zu einem
// Leerzeichen (DC-FA-CITE-001.a Schritt 4).
var wsRun = regexp.MustCompile(`\s+`)

// CheckCitations prüft per `<!-- d-check:cite <pfad>:<von>-<bis> -->`
// ausgezeichnete Zitate (spec/spezifikation.md §DC-FA-CITE-001.a):
// der whitespace-normalisierte Zitattext muss ein zusammenhängender
// Teilstring der normalisierten Quell-Spanne sein. Eine strukturell
// unbrauchbare Direktive (malformter Span, kein folgendes Zitat) ist
// fail-closed → error (der Aufrufer mappt auf Exit 2); Zitat-Fäule und
// Zitat-Abweichung sind Befunde (Exit 1). Arbeitet auf den rohen,
// fence-bewussten Zeilen.
func CheckCitations(fsys driven.Filesystem, file string, content []byte) ([]model.Finding, error) {
	prose := proseLines(content)
	var findings []model.Finding
	for i, pl := range prose {
		if !citeMarkerRe.MatchString(pl.raw) {
			continue
		}
		fs, err := citationForDirective(fsys, file, prose, i)
		if err != nil {
			return nil, err
		}
		findings = append(findings, fs...)
	}
	return findings, nil
}

// citationForDirective prüft eine einzelne d-check:cite-Direktive
// (prose[i]) und liefert höchstens einen Befund. Eine strukturell
// unbrauchbare Direktive (malformter Span, kein folgendes Zitat, nicht
// auflösbarer Pfad) ist fail-closed → error.
func citationForDirective(fsys driven.Filesystem, file string, prose []proseLine, i int) ([]model.Finding, error) {
	pl := prose[i]
	m := citeDirectiveRe.FindStringSubmatch(pl.raw)
	if m == nil {
		return nil, fmt.Errorf("%s:%d: malformte d-check:cite-Direktive — erwartet <!-- d-check:cite <pfad>:<von>-<bis> --> (DC-FA-CITE-001.a Schritt 1, fail-closed)", file, pl.no)
	}
	pathRaw := m[1]
	from, _ := strconv.Atoi(m[2])
	to := from
	if m[3] != "" {
		to, _ = strconv.Atoi(m[3])
	}
	target := fmt.Sprintf("%s:%d-%d", pathRaw, from, to)
	finding := func(reason, msg string) []model.Finding {
		return []model.Finding{{
			File: file, Line: pl.no, Rule: "citations",
			Target: target, Reason: reason, Message: msg,
		}}
	}
	// Schritt 2: das folgende Zitat (>-Block ODER inline). Fehlt es,
	// ist die Direktive unbrauchbar → fail-closed.
	quote, ok := citationQuote(prose, i)
	if !ok {
		return nil, fmt.Errorf("%s:%d: d-check:cite ohne folgendes Zitat (>-Block oder inline „…\"/\"…\") — DC-FA-CITE-001.a Schritt 2, fail-closed", file, pl.no)
	}
	// Schritt 3: Pfad auflösen; Repo-Escape und Zitat-Fäule sind Befunde.
	rel, escaped := resolveCitePath(file, pathRaw)
	if escaped {
		return finding(model.ReasonRepoEscape, "d-check:cite-Ziel verlässt die Repository-Wurzel"), nil
	}
	// Ungültiger Bereich (1-basiert): von < 1 oder von > bis. Fängt auch
	// die Untergrenze ab, bevor citationSpan mit von-1 indiziert.
	if from < 1 || from > to {
		return finding(ReasonCitationInvertedRange, "d-check:cite-Bereich ungültig (von < 1 oder von > bis)"), nil
	}
	span, ok := citationSpan(fsys, rel, from, to)
	if !ok {
		return finding(ReasonCitationOutOfRange, "d-check:cite-Ziel fehlt oder hat weniger als <bis> Zeilen"), nil
	}
	// Schritt 4/5: normalisieren, Mindestlänge, Teilstring-Vergleich.
	nq := normalizeWhitespace(quote)
	if len([]rune(nq)) < citationMinLen {
		return nil, nil // zu kurz, ungeprüft (dokumentierter Trade-off)
	}
	if strings.Contains(normalizeWhitespace(span), nq) {
		return nil, nil
	}
	return finding(ReasonCitationMismatch, "Zitattext ist kein zusammenhängender Teilstring der Quell-Spanne (Zitat-Fäule)"), nil
}

// citationSpan liest die Quell-Zeilen von..bis (1-basiert) und verbindet
// sie; ok=false, wenn die Datei fehlt oder weniger als bis Zeilen hat
// (finaler Zeilenumbruch zählt nicht als eigene Zeile).
func citationSpan(fsys driven.Filesystem, rel string, from, to int) (string, bool) {
	src, err := fsys.ReadFile(rel)
	if err != nil {
		return "", false
	}
	lines := strings.Split(string(src), "\n")
	n := len(lines)
	if n > 0 && lines[n-1] == "" {
		n--
	}
	if to > n {
		return "", false
	}
	return strings.Join(lines[from-1:to], "\n"), true
}

// resolveCitePath löst den Direktiven-Pfad auf: Datei-relativ bei
// `./`/`../`, sonst Wurzel-relativ (DC-FA-CITE-001.a Schritt 1/3, wie
// DC-FA-CODE-001.a).
func resolveCitePath(file, p string) (rel string, escaped bool) {
	if strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
		rel, escaped, _ = ResolveTarget(file, p)
		return rel, escaped
	}
	return ResolveConfigPath(p)
}

// citationQuote liefert den Zitattext, der der Direktive (prose[i]) folgt.
// Ist die nächste nicht-leere Zeile ein `>`-Blockquote, gilt dieser Block;
// sonst der nächste inline-Zitat-Span im selben Absatz (DC-FA-CITE-001.a
// Schritt 2). ok=false, wenn keines von beidem gefunden wird.
func citationQuote(prose []proseLine, i int) (string, bool) {
	j := i + 1
	for j < len(prose) && strings.TrimSpace(prose[j].raw) == "" {
		j++
	}
	if j >= len(prose) {
		return "", false
	}
	if strings.HasPrefix(strings.TrimLeft(prose[j].raw, " \t"), ">") {
		var b []string
		for k := j; k < len(prose); k++ {
			t := strings.TrimLeft(prose[k].raw, " \t")
			if !strings.HasPrefix(t, ">") {
				break
			}
			t = strings.TrimPrefix(strings.TrimPrefix(t, ">"), " ")
			b = append(b, t)
		}
		return strings.Join(b, "\n"), true
	}
	// inline: den Absatz ab j bis zur nächsten Leerzeile sammeln.
	var para []string
	for k := j; k < len(prose); k++ {
		if strings.TrimSpace(prose[k].raw) == "" {
			break
		}
		para = append(para, prose[k].raw)
	}
	return inlineQuoteSpan(strings.Join(para, "\n"))
}

// inlineQuoteSpan extrahiert den ersten inline-Zitat-Span aus text:
// `„` öffnet + `"` schließt, oder das erste `"` öffnet + das nächste
// `"` schließt — der früheste Öffner gewinnt (DC-FA-CITE-001.a Schritt 2).
func inlineQuoteSpan(text string) (string, bool) {
	low := strings.Index(text, "„")     // U+201E
	straight := strings.IndexByte(text, '"') // U+0022
	switch {
	case low >= 0 && (straight < 0 || low < straight):
		rest := text[low+len("„"):]
		if c := strings.IndexByte(rest, '"'); c >= 0 {
			return rest[:c], true
		}
	case straight >= 0:
		rest := text[straight+1:]
		if c := strings.IndexByte(rest, '"'); c >= 0 {
			return rest[:c], true
		}
	}
	return "", false
}

// normalizeWhitespace kollabiert jeden Whitespace-Lauf zu einem Leerzeichen
// und trimmt (DC-FA-CITE-001.a Schritt 4).
func normalizeWhitespace(s string) string {
	return strings.TrimSpace(wsRun.ReplaceAllString(s, " "))
}
