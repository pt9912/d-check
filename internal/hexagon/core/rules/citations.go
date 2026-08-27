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
//
// Der Pfad beginnt mit einem NICHT-Whitespace-Zeichen. Sonst wäre ein Pfad
// aus lauter Leerzeichen ein gültiger Parse mit leerem Ziel — und genau der
// entstünde, wenn der Pfad in Backticks stünde: das positionserhaltende
// Strippen ersetzt die Spanne durch Leerzeichen. Ein fehlender Pfad ist ein
// malformter Span und damit fail-closed (DC-FA-CITE-001.a Schritt 1), kein
// Befund gegen die leere Zeichenkette.
var citeDirectiveRe = regexp.MustCompile(`<!--\s*d-check:cite\s+(\S.*?):(\d+)(?:-(\d+))?\s*-->`)

// wsRun kollabiert jeden Whitespace-Lauf (inkl. Zeilenumbruch) zu einem
// Leerzeichen (DC-FA-CITE-001.a Schritt 4).
var wsRun = regexp.MustCompile(`\s+`)

// CheckCitations prüft per `<!-- d-check:cite <pfad>:<von>-<bis> -->`
// ausgezeichnete Zitate (spec/spezifikation.md §DC-FA-CITE-001.a):
// der whitespace-normalisierte Zitattext muss ein zusammenhängender
// Teilstring der normalisierten Quell-Spanne sein. Eine strukturell
// unbrauchbare Direktive (malformter Span, kein folgendes Zitat) ist
// fail-closed → error (der Aufrufer mappt auf Exit 2); Zitat-Fäule und
// Zitat-Abweichung sind Befunde (Exit 1).
//
// Die Frage „ist diese Zeile eine Direktive" bekommt die GETEILTE Antwort
// (ADR-0054): Marker-Suche und Direktiven-Parse laufen auf dem
// fence-bewussten, inline-code-gestrippten Text — eine in Inline-Code
// geschriebene Direktiv-Syntax ist eine Erwähnung, keine Direktive. Der
// ZITATTEXT dagegen wird roh gelesen: dort sind die Bytes der Gegenstand,
// nicht ihre Prosa-Rolle.
func CheckCitations(fsys driven.Filesystem, file string, content []byte) ([]model.Finding, error) {
	prose := proseLines(content)
	stripped := stripInlineCodeByLine(prose)
	var findings []model.Finding
	for i, pl := range prose {
		if !citeMarkerRe.MatchString(stripped[pl.no]) {
			continue
		}
		fs, err := citationForDirective(fsys, file, prose, stripped, i)
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
func citationForDirective(fsys driven.Filesystem, file string, prose []proseLine, stripped map[int]string, i int) ([]model.Finding, error) {
	pl := prose[i]
	m := citeDirectiveRe.FindStringSubmatch(stripped[pl.no])
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
//
// Ein Fenced-Block trennt: er begrenzt den Absatz wie eine Leerzeile
// (fencedBlockBetween, dieselbe Grenze wie proseParagraphs). Liegt einer
// zwischen Direktive und Kandidat, folgt der Direktive kein Zitat.
func citationQuote(prose []proseLine, i int) (string, bool) {
	j := i + 1
	for j < len(prose) && strings.TrimSpace(prose[j].raw) == "" {
		j++
	}
	if j >= len(prose) || fencedBlockWithin(prose, i, j) {
		return "", false
	}
	if strings.HasPrefix(strings.TrimLeft(prose[j].raw, " \t"), ">") {
		return citationBlockquote(prose, j), true
	}
	return inlineQuoteSpan(citationParagraph(prose, j))
}

// fencedBlockWithin meldet, ob zwischen prose[i] und prose[j] ein Fenced-Block
// lag — geprüft wird jeder Schritt, nicht nur der letzte: die Lücke entsteht
// beim Überspringen der Leerzeilen.
func fencedBlockWithin(prose []proseLine, i, j int) bool {
	for k := i + 1; k <= j; k++ {
		if fencedBlockBetween(prose[k-1].no, prose[k].no) {
			return true
		}
	}
	return false
}

// citationBlockquote sammelt die zusammenhängenden `>`-Zeilen ab prose[j],
// jeweils ohne ihr `> `-Präfix. Eine Leer-, Nicht-`>`- oder Fence-Grenze
// beendet den Block.
func citationBlockquote(prose []proseLine, j int) string {
	var b []string
	for k := j; k < len(prose); k++ {
		if fencedBlockBetween(prose[k-1].no, prose[k].no) {
			break
		}
		t := strings.TrimLeft(prose[k].raw, " \t")
		if !strings.HasPrefix(t, ">") {
			break
		}
		b = append(b, strings.TrimPrefix(strings.TrimPrefix(t, ">"), " "))
	}
	return strings.Join(b, "\n")
}

// citationParagraph liefert den Absatz ab prose[j] — begrenzt durch eine
// Leerzeile oder einen Fenced-Block, wie proseParagraphs.
func citationParagraph(prose []proseLine, j int) string {
	var para []string
	for k := j; k < len(prose); k++ {
		if strings.TrimSpace(prose[k].raw) == "" {
			break
		}
		if fencedBlockBetween(prose[k-1].no, prose[k].no) {
			break
		}
		para = append(para, prose[k].raw)
	}
	return strings.Join(para, "\n")
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
