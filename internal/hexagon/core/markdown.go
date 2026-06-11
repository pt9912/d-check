package core

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

// PreprocessMarkdown wendet Fence- und Inline-Code-Behandlung an.
// Fence-Zeilen selbst und Zeilen im Fence-Zustand entfallen komplett.
func PreprocessMarkdown(content []byte) []Line {
	var out []Line
	for _, pl := range proseLines(content) {
		out = append(out, Line{No: pl.no, Text: stripInlineCode(pl.raw)})
	}
	return out
}

// stripInlineCode ersetzt Backtick-Spans (inkl. Backticks) durch
// Leerzeichen gleicher Länge — positionserhaltend, damit angrenzender
// Text nicht zu Schein-Vorkommen verschmilzt (DC-FA-ID-001;
// spec/spezifikation.md §DC-FA-LINK-001.a Schritt 2). Ein Span wird
// von zwei gleich langen Backtick-Folgen begrenzt (die öffnende Folge
// bestimmt die schließende).
func stripInlineCode(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '`' {
			b.WriteByte(s[i])
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] == '`' {
			j++
		}
		closeAt := findClosingRun(s, j, j-i)
		if closeAt == -1 {
			b.WriteString(s[i:]) // keine schließende Folge
			break
		}
		for ; i < closeAt; i++ {
			b.WriteByte(' ')
		}
	}
	return b.String()
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
		forEachLink(ln.Text, func(ref LinkRef, _ LinkSpan) {
			ref.Line = no
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
	ref := LinkRef{Target: normalizeTarget(s[textEnd+2 : destEnd]), IsImage: isImage}
	span := LinkSpan{
		Start: i, End: destEnd + 1,
		TextStart: start + 1, TextEnd: textEnd,
		IsImage: isImage,
	}
	return ref, span, true
}

// normalizeTarget entquotet <…>-Ziele und trennt ein Titel-Suffix ab.
func normalizeTarget(t string) string {
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
