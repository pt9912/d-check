package rules

import "strings"

// Line ist eine vorverarbeitete Markdown-Zeile: Fenced-Code-Blöcke
// sind entfernt (ganze Zeilen), Inline-Code-Spans geleert
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritte 1–2).
type Line struct {
	No   int // 1-basiert
	Text string
}

// proseLine ist eine rohe Zeile außerhalb von Fenced-Code-Blöcken
// (1-basierte Zeilennummer) — der gemeinsame Fence-Automat von
// PreprocessMarkdown, extractHeadingLines und statusOf.
type proseLine struct {
	no  int
	raw string
}

// proseLines liefert alle Zeilen außerhalb von Fenced-Code-Blöcken;
// Fence-Zeilen selbst entfallen (spec/spezifikation.md
// §DC-FA-LINK-001.a Schritt 1).
func proseLines(content []byte) []proseLine {
	var out []proseLine
	inFence := false
	for i, raw := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimLeft(raw, " \t")
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		out = append(out, proseLine{no: i + 1, raw: raw})
	}
	return out
}

// fenceInfo liefert die Info-String-Sprache einer Fence-Öffner-Zeile:
// das erste Token nach der ```- bzw. ~~~-Folge, kleingeschrieben
// (z. B. "```mermaid" → "mermaid").
func fenceInfo(trimmed string) string {
	ch := trimmed[0]
	j := 0
	for j < len(trimmed) && trimmed[j] == ch {
		j++
	}
	info := strings.TrimSpace(trimmed[j:])
	if sp := strings.IndexAny(info, " \t"); sp != -1 {
		info = info[:sp]
	}
	return strings.ToLower(info)
}

// diagramFenceLines liefert die rohen Zeilen (1-basiert) INNERHALB von
// Fenced-Code-Blöcken, deren Öffner-Sprache in langs steht — das
// Gegenstück zu proseLines für das Modul diagrams (DC-FA-DIAG-001.a).
// Genau diese sonst opaken Diagramm-Fences werden gelesen; alle übrigen
// Fences und der Prosa-Text bleiben außen vor. Derselbe einfache
// Fence-Automat wie proseLines (Umschalten bei jeder ```/~~~-Zeile).
func diagramFenceLines(content []byte, langs map[string]bool) []proseLine {
	var out []proseLine
	inFence, open := false, false
	for i, raw := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimLeft(raw, " \t")
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			if inFence {
				inFence, open = false, false
			} else {
				inFence, open = true, langs[fenceInfo(trimmed)]
			}
			continue
		}
		if inFence && open {
			out = append(out, proseLine{no: i + 1, raw: raw})
		}
	}
	return out
}

// PreprocessMarkdown wendet Fence- und Inline-Code-Behandlung an.
// Fence-Zeilen selbst und Zeilen im Fence-Zustand entfallen komplett.
// Inline-Code-Spans werden absatzweise erkannt (CommonMark: ein Span
// darf Zeilenumbrüche enthalten) und positionserhaltend geleert
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritt 2).
func PreprocessMarkdown(content []byte) []Line {
	prose := proseLines(content)
	stripped := stripInlineCodeByLine(prose)
	out := make([]Line, 0, len(prose))
	for _, pl := range prose {
		text, ok := stripped[pl.no]
		if !ok {
			text = pl.raw // Leerzeile (kein Absatz-Mitglied)
		}
		out = append(out, Line{No: pl.no, Text: text})
	}
	return out
}

// proseParagraphs gruppiert Prosa-Zeilen zu Absätzen: Leerzeilen und
// Lücken in der Zeilennummerierung (übersprungene Fences) beenden
// einen Absatz — über solche Grenzen hinweg existiert kein
// Inline-Code-Span (spec/spezifikation.md §DC-FA-LINK-001.a Schritt 2).
func proseParagraphs(prose []proseLine) [][]proseLine {
	var groups [][]proseLine
	var cur []proseLine
	prevNo := 0
	for _, pl := range prose {
		blank := strings.TrimSpace(pl.raw) == ""
		if len(cur) > 0 && (blank || pl.no != prevNo+1) {
			groups = append(groups, cur)
			cur = nil
		}
		if !blank {
			cur = append(cur, pl)
		}
		prevNo = pl.no
	}
	if len(cur) > 0 {
		groups = append(groups, cur)
	}
	return groups
}

// stripInlineCodeByLine ersetzt Backtick-Spans (inkl. Backticks)
// absatzweise durch Leerzeichen gleicher Länge — positionserhaltend,
// damit angrenzender Text nicht zu Schein-Vorkommen verschmilzt
// (DC-FA-ID-001). Ein Span wird von zwei gleich langen
// Backtick-Folgen begrenzt (die öffnende Folge bestimmt die
// schließende) und darf innerhalb des Absatzes Zeilenumbrüche
// enthalten; eine im Absatz ungeschlossene Folge ist literal.
// Ergebnis: gestrippter Text pro Zeilennummer (nur Absatz-Zeilen).
func stripInlineCodeByLine(prose []proseLine) map[int]string {
	out := make(map[int]string, len(prose))
	for _, grp := range proseParagraphs(prose) {
		raws := make([]string, len(grp))
		for i, pl := range grp {
			raws[i] = pl.raw
		}
		joined := []byte(strings.Join(raws, "\n"))
		forEachInlineCodeSpan(string(joined), func(start, end, _, _ int) {
			for k := start; k < end; k++ {
				if joined[k] != '\n' {
					joined[k] = ' '
				}
			}
		})
		for i, line := range strings.Split(string(joined), "\n") {
			out[grp[i].no] = line
		}
	}
	return out
}

// inlineSpan ist ein Inline-Code-Span, der vollständig innerhalb
// einer Zeile liegt (Offsets relativ zur Zeile, wie
// forEachInlineCodeSpan).
type inlineSpan struct {
	start, end, valStart, valEnd int
}

// inlineSpansByLine liefert pro Zeilennummer die vollständig in der
// Zeile liegenden Code-Spans — absatzweise Erkennung wie
// stripInlineCodeByLine. Mehrzeilige Spans liefern keine Einträge
// (ein Wert mit Zeilenumbruch ist nie ein Pfad, §DC-FA-CODE-001.a).
func inlineSpansByLine(prose []proseLine) map[int][]inlineSpan {
	out := make(map[int][]inlineSpan)
	for _, grp := range proseParagraphs(prose) {
		raws := make([]string, len(grp))
		for i, pl := range grp {
			raws[i] = pl.raw
		}
		joined := strings.Join(raws, "\n")
		// Zeilen-Offsets im gejointen Absatz
		starts := make([]int, len(grp))
		off := 0
		for i, r := range raws {
			starts[i] = off
			off += len(r) + 1
		}
		forEachInlineCodeSpan(joined, func(start, end, valStart, valEnd int) {
			for i := range grp {
				lineStart := starts[i]
				lineEnd := lineStart + len(raws[i])
				if start >= lineStart && end <= lineEnd {
					out[grp[i].no] = append(out[grp[i].no], inlineSpan{
						start:    start - lineStart,
						end:      end - lineStart,
						valStart: valStart - lineStart,
						valEnd:   valEnd - lineStart,
					})
					break
				}
			}
		})
	}
	return out
}

// forEachInlineCodeSpan ruft fn für jeden Inline-Code-Span auf:
// [start,end) umfasst den Span inkl. Backticks, [valStart,valEnd)
// den Inhalt. Gemeinsamer Scanner für das positionserhaltende
// Stripping (übrige Module) und das Lesen der Span-Werte
// (Modul codepaths, spec/spezifikation.md §DC-FA-CODE-001.a).
func forEachInlineCodeSpan(s string, fn func(start, end, valStart, valEnd int)) {
	i := 0
	for i < len(s) {
		if s[i] != '`' {
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] == '`' {
			j++
		}
		closeAt := findClosingRun(s, j, j-i)
		if closeAt == -1 {
			// keine schließende Folge: die öffnende ist literal,
			// die Suche läuft dahinter weiter (CommonMark)
			i = j
			continue
		}
		fn(i, closeAt, j, closeAt-(j-i))
		i = closeAt
	}
}

// findClosingRun sucht ab from eine Backtick-Folge exakt der Länge
// runLen und liefert die Position dahinter (-1 wenn keine existiert).
func findClosingRun(s string, from, runLen int) int {
	for k := from; k+runLen <= len(s); k++ {
		if s[k] != '`' {
			continue
		}
		l := k
		for l < len(s) && s[l] == '`' {
			l++
		}
		if l-k == runLen {
			return l
		}
		k = l - 1
	}
	return -1
}

// headingLine ist ein ATX-Heading mit Zeilennummer (1-basiert) und
// Ebene — gemeinsamer Scanner für anchors (Slugs) und matrix
// (Sektions-Ausnahmen, Status-Heading).
type headingLine struct {
	line  int
	level int
	text  string
}

// parseATXHeading erkennt eine ATX-Heading-Zeile (#–######) und
// liefert Ebene und getrimmten Text.
func parseATXHeading(raw string) (level int, text string, ok bool) {
	trimmed := strings.TrimLeft(raw, " \t")
	if !strings.HasPrefix(trimmed, "#") {
		return 0, "", false
	}
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level > 6 || level >= len(trimmed) || (trimmed[level] != ' ' && trimmed[level] != '\t') {
		return 0, "", false
	}
	return level, strings.TrimSpace(trimmed[level+1:]), true
}

// extractHeadingLines liefert alle ATX-Headings (#–######) außerhalb
// von Fenced-Code-Blöcken, in Dokumentreihenfolge (Setext wird in 0.x
// nicht unterstützt — spec/spezifikation.md §DC-FA-ANCH-001.a).
func extractHeadingLines(content []byte) []headingLine {
	var out []headingLine
	for _, pl := range proseLines(content) {
		level, text, ok := parseATXHeading(pl.raw)
		if !ok {
			continue
		}
		out = append(out, headingLine{line: pl.no, level: level, text: text})
	}
	return out
}

// matchBracket liefert die Position der balanciert schließenden
// Klammer zur öffnenden an Position open.
func matchBracket(s string, open int, lo, hi byte) (int, bool) {
	depth := 0
	for j := open; j < len(s); j++ {
		switch s[j] {
		case lo:
			depth++
		case hi:
			depth--
			if depth == 0 {
				return j, true
			}
		}
	}
	return 0, false
}

// LinkRef ist ein extrahierter Markdown-Link.
type LinkRef struct {
	Line    int
	Target  string // roher Zielausdruck (ohne <>-Quoting und Titel)
	Text    string // Linktext zwischen [ und ] (für den matrix-Lineage-Match)
	IsImage bool
}

// LinkSpan beschreibt die Byte-Spannen eines Inline-Links innerhalb
// einer vorverarbeiteten Zeile: [Start,End) umfasst den ganzen Link,
// [TextStart,TextEnd) den Linktext — Grundlage der „im
// Linktext"-Erkennung des Moduls ids (spec/spezifikation.md
// §DC-FA-ID-001.a).
type LinkSpan struct {
	Start, End         int
	TextStart, TextEnd int
	IsImage            bool
}

// forEachLink ruft fn für jeden Inline-Link der Zeile auf —
// gemeinsamer Iterator von ExtractLinks und ExtractLinkSpans.
func forEachLink(text string, fn func(LinkRef, LinkSpan)) {
	for i := 0; i < len(text); i++ {
		ref, span, ok := parseLinkAt(text, i)
		if !ok {
			continue
		}
		fn(ref, span)
		i = span.End - 1
	}
}

// ExtractLinks findet Inline-Links [text](ziel) und Bilder
// ![alt](ziel); mehrere pro Zeile werden alle erfasst
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritt 3).
func ExtractLinks(lines []Line) []LinkRef {
	var refs []LinkRef
	for _, ln := range lines {
		no := ln.No
		text := ln.Text
		forEachLink(text, func(ref LinkRef, span LinkSpan) {
			ref.Line = no
			ref.Text = text[span.TextStart:span.TextEnd]
			refs = append(refs, ref)
		})
	}
	return refs
}

// ExtractLinkSpans liefert die Link-Spannen einer vorverarbeiteten
// Zeile in Vorkommens-Reihenfolge.
func ExtractLinkSpans(text string) []LinkSpan {
	var spans []LinkSpan
	forEachLink(text, func(_ LinkRef, span LinkSpan) {
		spans = append(spans, span)
	})
	return spans
}

// parseLinkAt liest an Position i einen Inline-Link (Linktext
// klammer-balanciert, Ziel mit balancierten Klammern).
func parseLinkAt(s string, i int) (LinkRef, LinkSpan, bool) {
	isImage := false
	start := i
	switch {
	case s[i] == '!' && i+1 < len(s) && s[i+1] == '[':
		isImage = true
		start = i + 1
	case s[i] != '[':
		return LinkRef{}, LinkSpan{}, false
	}
	textEnd, ok := matchBracket(s, start, '[', ']')
	if !ok || textEnd+1 >= len(s) || s[textEnd+1] != '(' {
		return LinkRef{}, LinkSpan{}, false
	}
	destEnd, ok := matchBracket(s, textEnd+1, '(', ')')
	if !ok {
		return LinkRef{}, LinkSpan{}, false
	}
	ref := LinkRef{Target: NormalizeTarget(s[textEnd+2 : destEnd]), IsImage: isImage}
	span := LinkSpan{
		Start: i, End: destEnd + 1,
		TextStart: start + 1, TextEnd: textEnd,
		IsImage: isImage,
	}
	return ref, span, true
}

// NormalizeTarget entquotet <…>-Ziele und trennt ein Titel-Suffix ab.
func NormalizeTarget(t string) string {
	t = strings.TrimSpace(t)
	if strings.HasPrefix(t, "<") {
		if end := strings.IndexByte(t, '>'); end != -1 {
			return t[1:end]
		}
		return strings.TrimPrefix(t, "<")
	}
	if sp := strings.IndexAny(t, " \t"); sp != -1 {
		return t[:sp]
	}
	return t
}
