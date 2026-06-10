package core

import "strings"

// Line ist eine vorverarbeitete Markdown-Zeile: Fenced-Code-Blöcke
// sind entfernt (ganze Zeilen), Inline-Code-Spans geleert
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritte 1–2).
type Line struct {
	No   int // 1-basiert
	Text string
}

// PreprocessMarkdown wendet Fence- und Inline-Code-Behandlung an.
// Fence-Zeilen selbst und Zeilen im Fence-Zustand entfallen komplett.
func PreprocessMarkdown(content []byte) []Line {
	var out []Line
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
		out = append(out, Line{No: i + 1, Text: stripInlineCode(raw)})
	}
	return out
}

// stripInlineCode leert Backtick-Spans. Ein Span wird von zwei gleich
// langen Backtick-Folgen begrenzt (die öffnende Folge bestimmt die
// schließende).
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
		i = closeAt // Span (inkl. Backticks) entfällt
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

// ExtractLinks findet Inline-Links [text](ziel) und Bilder
// ![alt](ziel); mehrere pro Zeile werden alle erfasst
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritt 3).
func ExtractLinks(lines []Line) []LinkRef {
	var refs []LinkRef
	for _, ln := range lines {
		for i := 0; i < len(ln.Text); i++ {
			ref, next, ok := parseLinkAt(ln.Text, i)
			if !ok {
				continue
			}
			ref.Line = ln.No
			refs = append(refs, ref)
			i = next
		}
	}
	return refs
}

// parseLinkAt liest an Position i einen Inline-Link (Linktext
// klammer-balanciert, Ziel mit balancierten Klammern).
func parseLinkAt(s string, i int) (LinkRef, int, bool) {
	isImage := false
	start := i
	switch {
	case s[i] == '!' && i+1 < len(s) && s[i+1] == '[':
		isImage = true
		start = i + 1
	case s[i] != '[':
		return LinkRef{}, i, false
	}
	textEnd, ok := matchBracket(s, start, '[', ']')
	if !ok || textEnd+1 >= len(s) || s[textEnd+1] != '(' {
		return LinkRef{}, i, false
	}
	destEnd, ok := matchBracket(s, textEnd+1, '(', ')')
	if !ok {
		return LinkRef{}, i, false
	}
	return LinkRef{Target: normalizeTarget(s[textEnd+2 : destEnd]), IsImage: isImage}, destEnd, true
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
