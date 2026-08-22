package rules

// Abschnitts-Mechanik, geteilt von der Closure-Fähigkeit des Moduls planning
// (§DC-FA-PLAN-001.a Schritte C3–C4) und vom Modul structure
// (§DC-FA-STRUCT-001.a Schritte 3–5). Die Spezifikation weist die
// Closure-Fähigkeit als Preset der allgemeinen Struktur-Semantik aus; zwei
// Kopien könnten auseinanderlaufen, ohne dass ein Test es merkt.

// SectionHead ist eine gefundene Abschnitts-Überschrift: 1-basierte Zeile und
// ATX-Ebene.
type SectionHead struct {
	Line  int
	Level int
}

// FindSectionHeads liefert alle echten ATX-Überschriften außerhalb von
// Fenced-Code-Blöcken, auf die match zutrifft, in Dokument-Reihenfolge. match
// bekommt die rohe Zeile.
func FindSectionHeads(lines []string, match func(raw string) bool) []SectionHead {
	var out []SectionHead
	inFence := false
	for i, raw := range lines {
		if FenceToggle(TrimFenceIndent(raw)) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lvl, _, ok := parseATXHeading(raw)
		if !ok || !match(raw) {
			continue
		}
		out = append(out, SectionHead{Line: i + 1, Level: lvl})
	}
	return out
}

// SectionEnd liefert die 1-basierte Zeilennummer der Überschrift, die den bei
// headingNo beginnenden Abschnitt beendet (0 ⇒ Dateiende): die nächste echte
// ATX-Überschrift gleicher oder höherer Ebene außerhalb von Fenced-Code.
func SectionEnd(lines []string, headingNo, level int) int {
	inFence := false
	for i := headingNo; i < len(lines); i++ {
		if FenceToggle(TrimFenceIndent(lines[i])) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if lvl, _, ok := parseATXHeading(lines[i]); ok && lvl <= level {
			return i + 1
		}
	}
	return 0
}

// SectionProse liefert den bereinigten Abschnitts-Text: Zeilen zwischen der
// Überschrift und dem Abschnitts-Ende, mit entfernten Fenced-Code-Blöcken und
// geleerten Inline-Code-Spans (dieselbe Vorverarbeitung wie im übrigen
// Scanner). Jede Bedingung beider Konsumenten arbeitet auf diesem Text.
func SectionProse(content []byte, lines []string, headingNo, level int) string {
	end := SectionEnd(lines, headingNo, level)
	var b []byte
	for _, ln := range PreprocessMarkdown(content) {
		if ln.No <= headingNo || (end != 0 && ln.No >= end) {
			continue
		}
		b = append(b, ln.Text...)
		b = append(b, '\n')
	}
	return string(b)
}

// SectionHeading ist eine Überschrift innerhalb eines Abschnitts: 1-basierte
// Zeile, ATX-Ebene und der getrimmte Überschriften-Text (ohne `#`-Folge).
type SectionHeading struct {
	Line  int
	Level int
	Text  string
}

// SectionHeadings liefert die Überschriften der Ebene level INNERHALB des bei
// headingNo beginnenden Abschnitts (§DC-FA-STRUCT-001.a Schritt 6) — mit
// derselben Erkennung, die den Abschnitt findet: echte ATX-Überschriften
// außerhalb von Fenced-Code, beliebig viel führender Weißraum, Leerzeichen
// oder Tab als Trenner. Der Bereich endet an SectionEnd; eine Ebene flacher
// als der Abschnitt kann darin nicht vorkommen.
func SectionHeadings(lines []string, headingNo, sectionLevel, level int) []SectionHeading {
	end := SectionEnd(lines, headingNo, sectionLevel)
	var out []SectionHeading
	inFence := false
	for i := headingNo; i < len(lines); i++ {
		if end != 0 && i+1 >= end {
			break
		}
		if FenceToggle(TrimFenceIndent(lines[i])) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lvl, text, ok := parseATXHeading(lines[i])
		if !ok || lvl != level {
			continue
		}
		out = append(out, SectionHeading{Line: i + 1, Level: lvl, Text: text})
	}
	return out
}
