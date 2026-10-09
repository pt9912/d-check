package rules

import (
	"regexp"
	"sort"
	"strings"
)

// blankedLinkTargetLines liefert je Zeile, in der ein Link-Ziel geleert wurde,
// ihre Fassung ohne dieses Ziel (DC-FA-VCS-001.a Schritt 4). Was ein Link ist,
// beantwortet dieselbe Erkennung wie das Modul links: PreprocessMarkdown
// (Fenced-Code entfällt, Code-Spans absatzweise positionserhaltend geleert),
// ExtractLinkSpans und definitionRe. Vier Filter engen das ein, jeder nur in
// Richtung Drift: eine Zeile mit mindestens vier Spalten Einzug oder in einem
// HTML-Block bleibt unverändert; ein Link mit escapter Klammer ist keiner; eine
// Fußnote (`[^…]:`) und eine Referenz-Definition ohne pfadartiges Ziel bleiben
// stehen.
// GRENZE: Ein Ziel auf der Folgezeile wird nicht geleert, sein Nachzug bleibt
// Drift. Eine Absatz-Folgezeile in der Form einer Referenz-Definition wird
// geleert, wie die links-Erkennung sie liest, auch wo Markdown sie als Text
// rendert.
func blankedLinkTargetLines(content []byte) map[int]string {
	raw := splitLines(content)
	html := htmlBlockLines(raw)
	out := make(map[int]string)
	for _, ln := range PreprocessMarkdown(content) {
		r := raw[ln.No-1]
		if html[ln.No] || indentColumns(r) >= 4 || len(ln.Text) != len(r) {
			continue
		}
		if cuts := linkTargetCuts(ln.Text); len(cuts) > 0 {
			out[ln.No] = cutRanges(r, cuts)
		}
	}
	return out
}

// linkTargetCuts liefert die Byte-Bereiche der Link-Ziele einer vorverarbeiteten
// Zeile: das Ziel-Token jedes Inline-Links und Bilds (ohne Titel) und das Ziel
// einer Referenz-Definition.
func linkTargetCuts(text string) [][2]int {
	var cuts [][2]int
	for _, sp := range ExtractLinkSpans(text) {
		if escapedAt(text, sp.TextStart-1) || escapedAt(text, sp.TextEnd) || escapedAt(text, sp.End-1) {
			continue
		}
		if c, ok := targetToken(text, sp.TextEnd+2, sp.End-1); ok {
			cuts = append(cuts, c)
		}
	}
	if m := definitionRe.FindStringSubmatchIndex(text); m != nil {
		label, target := text[m[2]:m[3]], text[m[4]:m[5]]
		if !strings.HasPrefix(label, "^") && (strings.HasPrefix(target, "<") || strings.ContainsAny(target, "./#:")) {
			cuts = append(cuts, [2]int{m[4], m[5]})
		}
	}
	return cuts
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

// cutRanges entfernt die Bereiche aus line; überlappende werden übersprungen.
func cutRanges(line string, cuts [][2]int) string {
	sort.Slice(cuts, func(a, b int) bool { return cuts[a][0] < cuts[b][0] })
	var b strings.Builder
	last := 0
	for _, c := range cuts {
		if c[0] < last {
			continue
		}
		b.WriteString(line[last:c[0]])
		last = c[1]
	}
	b.WriteString(line[last:])
	return b.String()
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

// htmlBlockLines liefert die 1-basierten Zeilen, die zu einem HTML-Block
// gehören: ab einer Zeile, die mit `<` und einem Buchstaben, `/`, `!` oder `?`
// beginnt, bis zur nächsten Leerzeile — bei `<pre>`, `<script>`, `<style>` und
// `<textarea>` bis zur Zeile mit dem schließenden Tag.
// GRENZE: weiter als CommonMark, das nicht jedes Tag als Blockanfang liest --
// eine Absatzzeile, die mit Inline-HTML beginnt, gilt ebenfalls als Block, ihr
// Nachzug bleibt Drift.
func htmlBlockLines(lines []string) map[int]bool {
	out := make(map[int]bool)
	inBlock, raw := false, false
	for i, l := range lines {
		if !inBlock && htmlBlockStartRE.MatchString(l) {
			inBlock, raw = true, htmlRawStartRE.MatchString(l)
		}
		if !inBlock {
			continue
		}
		if !raw && strings.TrimSpace(l) == "" {
			inBlock = false
			continue
		}
		out[i+1] = true
		if raw && htmlRawEndRE.MatchString(l) {
			inBlock = false
		}
	}
	return out
}
