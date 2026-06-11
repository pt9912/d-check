package core

// checkIDs ist das Regelmodul `ids` (DC-FA-ID-001): nackte Kennungen
// im Fließtext — außerhalb von Inline-Code und Fenced-Code-Blöcken —
// müssen als Markdown-Link auf ihre Definition ausgeführt sein. Die
// Muster werden in Deklarationsreihenfolge gematcht; das erste
// passende Muster gewinnt pro Vorkommen (spec/spezifikation.md
// §DC-FA-ID-001.a).
func checkIDs(file string, lines []Line, patterns []IDPattern) []Finding {
	if len(patterns) == 0 {
		return nil
	}
	var findings []Finding
	for _, ln := range lines {
		spans := ExtractLinkSpans(ln.Text)
		var claimed [][2]int
		for _, p := range patterns {
			for _, m := range p.Regex.FindAllStringIndex(ln.Text, -1) {
				if overlapsClaimed(claimed, m[0], m[1]) {
					continue // Vorkommen gehört einem früheren Muster
				}
				claimed = append(claimed, [2]int{m[0], m[1]})
				if idOccurrenceExempt(spans, m[0], m[1]) {
					continue
				}
				findings = append(findings, Finding{
					File: file, Line: ln.No, Rule: "ids",
					Target: ln.Text[m[0]:m[1]], Reason: ReasonIDUnlinked,
					Message: "Kennung ohne Link auf ihre Definition",
				})
			}
		}
	}
	return findings
}

// idOccurrenceExempt prüft, ob ein Vorkommen [start,end) linkpflichtfrei
// ist: es liegt im Linktext eines Nicht-Bild-Links (= verlinkt) oder
// innerhalb der Link-/Bild-Syntax außerhalb des Linktexts
// (Ziel-Klammer, Alt-Text — kein Fließtext;
// spec/spezifikation.md §DC-FA-ID-001.a).
func idOccurrenceExempt(spans []LinkSpan, start, end int) bool {
	for _, sp := range spans {
		if start >= sp.TextStart && end <= sp.TextEnd && !sp.IsImage {
			return true // verlinkt: im Linktext eines Markdown-Links
		}
		if start >= sp.Start && end <= sp.End {
			return true // Link-Syntax (Ziel) oder Bildreferenz: kein Fließtext
		}
	}
	return false
}

// overlapsClaimed prüft Überlappung mit bereits zugeordneten
// Vorkommen (Muster-Präzedenz: erstes Match gewinnt).
func overlapsClaimed(claimed [][2]int, start, end int) bool {
	for _, c := range claimed {
		if start < c[1] && c[0] < end {
			return true
		}
	}
	return false
}
