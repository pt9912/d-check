package core

import "strings"

// checkSpans prüft Markdown-Span-Artefakte (DC-FA-SPAN-001,
// spec/spezifikation.md §DC-FA-SPAN-001.a): ungeschlossene
// Code-Span-Öffnungen, die an Nicht-Whitespace kleben
// (span-unclosed), und Link-Syntax im Linktext eines weiteren Links
// (span-nested-link). Es gibt keinen Opt-out-Marker — Befunde werden
// behoben, nicht unterdrückt.
func checkSpans(file string, content []byte, lines []Line) []Finding {
	var findings []Finding
	findings = append(findings, checkUnclosedSpans(file, content)...)
	findings = append(findings, checkNestedLinks(file, lines)...)
	return findings
}

// checkUnclosedSpans arbeitet absatzweise (gleiche Absatz-Semantik
// wie die Vorverarbeitung): Backtick-Folgen werden links nach rechts
// gepaart; eine ungeschlossene Folge, auf die unmittelbar
// Nicht-Whitespace folgt, kippt die Parität des Absatzes und ist ein
// Befund (§DC-FA-SPAN-001.a Schritt 1).
func checkUnclosedSpans(file string, content []byte) []Finding {
	var findings []Finding
	for _, grp := range proseParagraphs(proseLines(content)) {
		raws := make([]string, len(grp))
		for i, pl := range grp {
			raws[i] = pl.raw
		}
		joined := strings.Join(raws, "\n")
		for _, open := range unclosedRuns(joined) {
			if !sticksToText(joined, open) {
				continue
			}
			line, target := unclosedTarget(grp, joined, open)
			findings = append(findings, Finding{
				File: file, Line: line, Rule: "spans",
				Target: target, Reason: ReasonSpanUnclosed,
			})
		}
	}
	return findings
}

// runSpan ist eine Backtick-Folge [start,end) im gejointen Absatz.
type runSpan struct{ start, end int }

// unclosedRuns liefert alle Backtick-Folgen ohne gleich lange
// schließende Folge — Pairing wie forEachInlineCodeSpan (die
// öffnende Folge bestimmt die schließende, ungeschlossene Folgen
// sind literal und die Suche läuft dahinter weiter).
func unclosedRuns(s string) []runSpan {
	var open []runSpan
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
			open = append(open, runSpan{start: i, end: j})
			i = j
			continue
		}
		i = closeAt
	}
	return open
}

// sticksToText: unmittelbar nach der Folge steht Nicht-Whitespace
// (Leerzeichen, Tab, Zeilenumbruch und Absatz-Ende zählen als
// Whitespace — solche Folgen sind beabsichtigt literal).
func sticksToText(s string, r runSpan) bool {
	if r.end >= len(s) {
		return false
	}
	switch s[r.end] {
	case ' ', '\t', '\n':
		return false
	}
	return true
}

// unclosedTarget liefert Zeilennummer der Folge und das gemeldete
// Ziel: Backtick-Folge samt folgender Nicht-Whitespace-Zeichen,
// gekappt auf 30 Zeichen (§DC-FA-SPAN-001.a Schritt 1).
func unclosedTarget(grp []proseLine, joined string, r runSpan) (int, string) {
	line := grp[0].no
	off := 0
	for i, pl := range grp {
		lineEnd := off + len(pl.raw)
		if r.start >= off && r.start <= lineEnd {
			line = grp[i].no
			break
		}
		off = lineEnd + 1
	}
	end := r.end
	for end < len(joined) && end-r.start < 30 {
		c := joined[end]
		if c == ' ' || c == '\t' || c == '\n' {
			break
		}
		end++
	}
	return line, joined[r.start:end]
}

// checkNestedLinks meldet Link-Syntax im Linktext eines weiteren
// Links — lexikalisch `](`…`)` direkt gefolgt von `](` auf den
// vorverarbeiteten Zeilen (§DC-FA-SPAN-001.a Schritt 2).
func checkNestedLinks(file string, lines []Line) []Finding {
	var findings []Finding
	for _, ln := range lines {
		for _, m := range nestedLinkMatches(ln.Text) {
			findings = append(findings, Finding{
				File: file, Line: ln.No, Rule: "spans",
				Target: m, Reason: ReasonSpanNestedLink,
			})
		}
	}
	return findings
}

// nestedLinkMatches findet alle `](…)](`-Vorkommen einer Zeile;
// das gemeldete Ziel ist der Treffer, gekappt auf 40 Zeichen.
// Bildreferenzen als Linktext (`[![…](…)](…)` — Badge-Muster) sind
// legales Markdown und kein Treffer (§DC-FA-SPAN-001.a Schritt 2).
func nestedLinkMatches(text string) []string {
	var out []string
	for i := 0; i+1 < len(text); i++ {
		if text[i] != ']' || text[i+1] != '(' {
			continue
		}
		closeIdx, ok := matchBracket(text, i+1, '(', ')')
		if !ok {
			continue
		}
		if closeIdx+2 < len(text) && text[closeIdx+1] == ']' && text[closeIdx+2] == '(' {
			if !innerIsImage(text, i) {
				end := closeIdx + 3
				if end-i > 40 {
					end = i + 40
				}
				out = append(out, text[i:end])
			}
		}
		i = closeIdx
	}
	return out
}

// innerIsImage prüft, ob der innere Link, dessen Text an Position
// closeBracket endet, eine Bildreferenz ist (öffnende Klammer per
// Rückwärts-Balance gesucht, `!` unmittelbar davor).
func innerIsImage(text string, closeBracket int) bool {
	depth := 1
	for p := closeBracket - 1; p >= 0; p-- {
		switch text[p] {
		case ']':
			depth++
		case '[':
			depth--
			if depth == 0 {
				return p > 0 && text[p-1] == '!'
			}
		}
	}
	return false
}
