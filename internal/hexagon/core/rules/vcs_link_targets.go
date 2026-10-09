package rules

import (
	"bytes"
	"path"
	"regexp"
	"sort"
	"strings"
)

// normalizedLinkTargetLines liefert je Zeile, in der ein Link-Ziel normiert
// wurde, ihre Fassung mit normierten Zielen (DC-FA-VCS-001.a Schritt 4).
// Normiert wird ein Ziel, das resolve auflöst (gegen die Vereinigung der
// Pfad-Bäume von BASE und HEAD): es wird durch Marke, Dateiname und Anker
// ersetzt, sodass ein reiner Pfad-Nachzug denselben Text ergibt; ein Ziel, das
// nicht auflöst, bleibt roh. Was ein Link ist,
// beantwortet dieselbe Erkennung wie das Modul links: PreprocessMarkdown
// (Fenced-Code entfällt, Code-Spans absatzweise positionserhaltend geleert),
// ExtractLinkSpans und definitionRe. Filter engen das ein, jeder nur in
// Richtung Drift: opaqueLines nimmt ganze Zeilen aus, linkTargetCuts einzelne
// Links.
// GRENZE: Ein Pfad-Nachzug eines auflösenden Ziels in Code, der hinter einer
// Listenmarke beginnt (Fence, HTML-Block, eingerückter Code), wird normiert
// und geht durch -- ein Inhaltswort dort nicht, es löst nicht auf. Fail-safe
// bleibt Drift: ein Ziel auf der Folgezeile, ein Nachzug auf eine Datei mit
// anderem Namen, eine Referenz-Definition auf einer Zeile mit CRLF-Ende, ein
// Link, den ein zeilenlokal falsch gepaarter Code-Span verdeckt.
func normalizedLinkTargetLines(content []byte, resolve func(string) (string, bool)) map[int]string {
	raw := splitLines(content)
	opaque := opaqueLines(raw)
	out := make(map[int]string)
	for _, ln := range PreprocessMarkdown(content) {
		r := raw[ln.No-1]
		if opaque[ln.No] || len(ln.Text) != len(r) {
			continue
		}
		var repls []targetRepl
		for _, c := range linkTargetCuts(ln.Text, lineCodeSpans(r)) {
			if norm, ok := resolve(r[c[0]:c[1]]); ok {
				repls = append(repls, targetRepl{c[0], c[1], norm})
			}
		}
		if len(repls) > 0 {
			out[ln.No] = replaceRanges(r, repls)
		}
	}
	return out
}

// normalizedTargetMark steht vor jedem normierten Ziel. Normiert wird nur,
// wenn keine der beiden Fassungen das Zeichen trägt (markFreeTree) -- ein
// Ziel, das nur wie eine normierte Form aussieht, gleicht ihr deshalb nie.
const normalizedTargetMark = "\x00"

// markFreeTree liefert tree, wenn keine der beiden Fassungen die Marke trägt,
// sonst nil: dann wird nicht normiert, und eine rohe Marke im Text kann einer
// normierten Form nicht gleichen (fail-safe -- der Nachzug bleibt Drift).
func markFreeTree(tree map[string]bool, base, head []byte) map[string]bool {
	if bytes.Contains(base, []byte(normalizedTargetMark)) || bytes.Contains(head, []byte(normalizedTargetMark)) {
		return nil
	}
	return tree
}

// linkTargetResolver liefert für die Datei file die Auflösung ihrer Link-Ziele
// gegen tree, nil ohne tree. tree ist die Vereinigung der Pfad-Bäume von BASE
// und HEAD (pathTree): beide Seiten lösen gleich auf, ein unveränderter Link
// auf eine gelöschte Datei bleibt gleich. Ein relatives Ziel, das auf einen
// Eintrag von tree zeigt, wird zu Marke, Dateiname und Anker; ein absolutes,
// ein externes (`://`), eines mit Query, ein reiner Anker und eines außerhalb
// des Repos lösen nicht auf.
// GRENZE: Ein Nachzug auf eine Datei, die nur noch in BASE existiert, geht
// durch -- den toten Link meldet das Modul links.
func linkTargetResolver(file string, tree map[string]bool) func(string) (string, bool) {
	if tree == nil {
		return nil
	}
	return func(target string) (string, bool) {
		t := strings.TrimSuffix(strings.TrimPrefix(target, "<"), ">")
		anchor := ""
		if i := strings.IndexByte(t, '#'); i >= 0 {
			t, anchor = t[:i], t[i:]
		}
		if t == "" || strings.HasPrefix(t, "/") || strings.Contains(t, "://") || strings.Contains(t, "?") {
			return "", false
		}
		p := path.Clean(path.Join(path.Dir(file), t))
		if p == ".." || strings.HasPrefix(p, "../") || !tree[p] {
			return "", false
		}
		return normalizedTargetMark + path.Base(p) + anchor, true
	}
}

// pathTree liefert die Menge der Pfade und aller ihrer Verzeichnisse.
func pathTree(paths []string) map[string]bool {
	out := make(map[string]bool, len(paths))
	for _, p := range paths {
		for ; p != "." && p != "/" && !out[p]; p = path.Dir(p) {
			out[p] = true
		}
	}
	return out
}

// targetRepl ersetzt line[start:end] durch text.
type targetRepl struct {
	start, end int
	text       string
}

// replaceRanges wendet die Ersetzungen an; überlappende werden übersprungen.
func replaceRanges(line string, repls []targetRepl) string {
	sort.Slice(repls, func(a, b int) bool { return repls[a].start < repls[b].start })
	var b strings.Builder
	last := 0
	for _, r := range repls {
		if r.start < last {
			continue
		}
		b.WriteString(line[last:r.start])
		b.WriteString(r.text)
		last = r.end
	}
	b.WriteString(line[last:])
	return b.String()
}

// listCodeRE trifft einen Listenpunkt, dessen Inhalt mit eingerücktem Code
// beginnt: fünf Leerzeichen oder ein Tab hinter der Marke.
var listCodeRE = regexp.MustCompile(`^ {0,3}(?:[-+*]|[0-9]{1,9}[.)])(?: {5}|[ ]*\t)`)

// opaqueLines liefert die 1-basierten Zeilen, deren Links nicht normiert
// werden: eingerückter Code (mindestens vier Spalten, auch hinter einer
// Listenmarke), jede Zitatzeile, Fenced-Code nach einem strengen Automaten
// (zusätzlich zu dem der Vorverarbeitung, deren Fence-Erkennung jeden Einzug
// zulässt) und HTML-Blöcke.
// GRENZE: Auch ein Link in einem Zitat oder in einem eingerückten
// Listen-Folgeabsatz bleibt ungeleert, sein Nachzug bleibt Drift.
func opaqueLines(lines []string) map[int]bool {
	out := htmlBlockLines(lines)
	for no := range strictFenceLines(lines) {
		out[no] = true
	}
	for i, l := range lines {
		if indentColumns(l) >= 4 || listCodeRE.MatchString(l) || strings.HasPrefix(strings.TrimLeft(l, " "), ">") {
			out[i+1] = true
		}
	}
	return out
}

// strictFenceLines liefert die Zeilen eines Fenced-Code-Blocks samt Öffner und
// Schließer: Öffner mit höchstens drei Leerzeichen Einzug und mindestens drei
// gleichen Zeichen, Schließer aus demselben Zeichen in mindestens derselben
// Zahl; ohne Schließer reicht der Block bis zum Dateiende.
func strictFenceLines(lines []string) map[int]bool {
	out := make(map[int]bool)
	var char byte
	n := 0
	for i, l := range lines {
		lead := len(l) - len(strings.TrimLeft(l, " "))
		trimmed := strings.TrimRight(l[lead:], " \t\r")
		c, run := FenceRun(trimmed)
		if n == 0 {
			if lead <= 3 && run >= 3 && (c == '~' || strings.IndexByte(trimmed[run:], '`') == -1) {
				char, n = c, run
				out[i+1] = true
			}
			continue
		}
		out[i+1] = true
		if lead <= 3 && c == char && run >= n && run == len(trimmed) {
			n = 0
		}
	}
	return out
}

// lineCodeSpans liefert die Code-Spans, die ein Scan der Zeile allein findet;
// ein Span aus der Vorzeile kann die Paarung verschieben (fail-safe).
func lineCodeSpans(line string) [][2]int {
	var spans [][2]int
	forEachInlineCodeSpan(line, func(start, end, _, _ int) {
		spans = append(spans, [2]int{start, end})
	})
	return spans
}

// linkDestRE trifft einen gültigen Zielausdruck: ein Ziel (`<…>` oder ohne
// Leerraum) und ein optionaler Titel.
var linkDestRE = regexp.MustCompile(`^[ \t]*(?:<[^<>\n]*>|[^\s<>]+)(?:[ \t]+(?:"[^"]*"|'[^']*'|\([^()]*\)))?[ \t]*$`)

// linkTargetCuts liefert die Byte-Bereiche der Link-Ziele einer vorverarbeiteten
// Zeile: das Ziel-Token jedes Inline-Links und Bilds (ohne Titel) und das Ziel
// einer Referenz-Definition; auf einer Definitionszeile zählt nur ihr Ziel,
// nicht ein link-förmiger Titel. Kein Link ist ein Treffer mit escapter
// Linktext-Klammer, mit ungültigem Zielausdruck, dessen öffnende Klammer in
// einem zeilenlokalen Code-Span der Rohzeile liegt (code) oder dessen Linktext
// einen Link enthält, und keiner ist eine Fußnote (`[^…]:`). Eine escapte
// Zielklammer lässt das Ziel auf `\` enden; es löst nie auf und bleibt roh.
func linkTargetCuts(text string, code [][2]int) [][2]int {
	if m := definitionRe.FindStringSubmatchIndex(text); m != nil {
		if strings.HasPrefix(text[m[2]:m[3]], "^") {
			return nil
		}
		return [][2]int{{m[4], m[5]}}
	}
	var cuts [][2]int
	for _, sp := range ExtractLinkSpans(text) {
		if escapedAt(text, sp.TextStart-1) || escapedAt(text, sp.TextEnd) {
			continue
		}
		if overlapsAny(code, sp.Start, sp.Start+1) || containsLink(text[sp.TextStart:sp.TextEnd]) {
			continue
		}
		if !linkDestRE.MatchString(text[sp.TextEnd+2 : sp.End-1]) {
			continue
		}
		if c, ok := targetToken(text, sp.TextEnd+2, sp.End-1); ok {
			cuts = append(cuts, c)
		}
	}
	return cuts
}

// containsLink meldet, ob ein Linktext selbst einen Link enthält -- dann ist
// der äußere Treffer nach CommonMark kein Link; ein Bild darin ist erlaubt.
func containsLink(linkText string) bool {
	for _, sp := range ExtractLinkSpans(linkText) {
		if !sp.IsImage {
			return true
		}
	}
	return false
}

// overlapsAny meldet, ob [start,end) einen der halboffenen Bereiche berührt.
func overlapsAny(spans [][2]int, start, end int) bool {
	for _, sp := range spans {
		if start < sp[1] && sp[0] < end {
			return true
		}
	}
	return false
}

// targetToken grenzt im Zielausdruck text[from:to] das Ziel ab, wie
// NormalizeTarget es liest: `<…>` samt Klammern, sonst bis zum ersten
// Leerraum; ein Titel dahinter bleibt.
func targetToken(text string, from, to int) ([2]int, bool) {
	i := from
	for i < to && (text[i] == ' ' || text[i] == '\t') {
		i++
	}
	if i >= to {
		return [2]int{}, false
	}
	if text[i] == '<' {
		if j := strings.IndexByte(text[i:to], '>'); j != -1 {
			return [2]int{i, i + j + 1}, true
		}
		return [2]int{}, false
	}
	j := i
	for j < to && text[j] != ' ' && text[j] != '\t' {
		j++
	}
	return [2]int{i, j}, true
}

// escapedAt meldet, ob vor s[i] eine ungerade Zahl Backslashes steht.
func escapedAt(s string, i int) bool {
	n := 0
	for k := i - 1; k >= 0 && s[k] == '\\'; k-- {
		n++
	}
	return n%2 == 1
}

// indentColumns zählt den Einzug einer Zeile in Spalten; ein Tab springt zum
// nächsten Vielfachen von vier (CommonMark).
func indentColumns(line string) int {
	col := 0
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case ' ':
			col++
		case '\t':
			col += 4 - col%4
		default:
			return col
		}
	}
	return col
}

var (
	htmlBlockStartRE = regexp.MustCompile(`^ {0,3}<[A-Za-z/!?]`)
	htmlRawStartRE   = regexp.MustCompile(`(?i)^ {0,3}<(?:pre|script|style|textarea)(?:[\s>]|$)`)
	htmlRawEndRE     = regexp.MustCompile(`(?i)</(?:pre|script|style|textarea)>`)
)

// htmlBlockEnd liefert für die Startzeile eines HTML-Blocks die Bedingung, an
// der er endet (CommonMark-Typen 1 bis 5), oder nil für einen Block, der an der
// nächsten Leerzeile endet.
func htmlBlockEnd(start string) func(string) bool {
	t := strings.TrimLeft(start, " ")
	contains := func(marker string) func(string) bool {
		return func(l string) bool { return strings.Contains(l, marker) }
	}
	switch {
	case htmlRawStartRE.MatchString(start):
		return htmlRawEndRE.MatchString
	case strings.HasPrefix(t, "<!--"):
		return contains("-->")
	case strings.HasPrefix(t, "<?"):
		return contains("?>")
	case strings.HasPrefix(t, "<![CDATA["):
		return contains("]]>")
	case strings.HasPrefix(t, "<!"):
		return contains(">")
	}
	return nil
}

// htmlBlockLines liefert die 1-basierten Zeilen, die zu einem HTML-Block
// gehören: ab einer Zeile, die mit `<` und einem Buchstaben, `/`, `!` oder `?`
// beginnt, bis zum Endmarker ihres Typs (`</pre>` u. a., `-->`, `?>`, `]]>`,
// `>`) oder sonst bis zur nächsten Leerzeile.
// GRENZE: weiter als CommonMark, das nicht jedes Tag als Blockanfang liest --
// eine Absatzzeile, die mit Inline-HTML beginnt, gilt ebenfalls als Block, ihr
// Nachzug bleibt Drift.
func htmlBlockLines(lines []string) map[int]bool {
	out := make(map[int]bool)
	inBlock := false
	var end func(string) bool
	for i, l := range lines {
		if !inBlock && htmlBlockStartRE.MatchString(l) {
			inBlock, end = true, htmlBlockEnd(l)
		}
		if !inBlock {
			continue
		}
		if end == nil && strings.TrimSpace(l) == "" {
			inBlock = false
			continue
		}
		out[i+1] = true
		if end != nil && end(l) {
			inBlock = false
		}
	}
	return out
}
