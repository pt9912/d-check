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
// langen Backtick-Folgen begrenzt (längste Folge zuerst, d. h. die
// öffnende Folge bestimmt die schließende).
func stripInlineCode(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '`' {
			b.WriteByte(s[i])
			i++
			continue
		}
		// öffnende Backtick-Folge vermessen
		j := i
		for j < len(s) && s[j] == '`' {
			j++
		}
		runLen := j - i
		// schließende Folge gleicher Länge suchen
		close := -1
		for k := j; k+runLen <= len(s); k++ {
			if s[k] != '`' {
				continue
			}
			l := k
			for l < len(s) && s[l] == '`' {
				l++
			}
			if l-k == runLen {
				close = l
				break
			}
			k = l - 1
		}
		if close == -1 {
			// keine schließende Folge: Rest unverändert übernehmen
			b.WriteString(s[i:])
			break
		}
		i = close // Span (inkl. Backticks) entfällt
	}
	return b.String()
}

// LinkRef ist ein extrahierter Markdown-Link.
type LinkRef struct {
	Line    int
	Target  string // roher Zielausdruck (ohne <>-Quoting und Titel)
	IsImage bool
}

// ExtractLinks findet Inline-Links [text](ziel) und Bilder
// ![alt](ziel); mehrere pro Zeile werden alle erfasst
// (spec/spezifikation.md §DC-FA-LINK-001.a Schritt 3). Der Linktext
// wird klammer-balanciert gelesen.
func ExtractLinks(lines []Line) []LinkRef {
	var refs []LinkRef
	for _, ln := range lines {
		s := ln.Text
		for i := 0; i < len(s); i++ {
			isImage := false
			start := i
			if s[i] == '!' && i+1 < len(s) && s[i+1] == '[' {
				isImage = true
				start = i + 1
			} else if s[i] != '[' {
				continue
			}
			// Linktext klammer-balanciert bis zur schließenden ]
			depth := 0
			j := start
			for ; j < len(s); j++ {
				if s[j] == '[' {
					depth++
				} else if s[j] == ']' {
					depth--
					if depth == 0 {
						break
					}
				}
			}
			if j >= len(s) || j+1 >= len(s) || s[j+1] != '(' {
				continue
			}
			// Ziel bis zur schließenden ) (Klammern im Ziel balanciert)
			k := j + 2
			pdepth := 1
			for ; k < len(s); k++ {
				if s[k] == '(' {
					pdepth++
				} else if s[k] == ')' {
					pdepth--
					if pdepth == 0 {
						break
					}
				}
			}
			if k >= len(s) {
				continue
			}
			target := normalizeTarget(s[j+2 : k])
			refs = append(refs, LinkRef{Line: ln.No, Target: target, IsImage: isImage})
			i = k
		}
	}
	return refs
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
