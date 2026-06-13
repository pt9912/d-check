package core

import "strings"

// checkIDs ist das Regelmodul `ids` (DC-FA-ID-001): nackte Kennungen
// im Fließtext — außerhalb von Inline-Code und Fenced-Code-Blöcken —
// müssen als Markdown-Link auf ihre Definition ausgeführt sein. Die
// Muster werden in Deklarationsreihenfolge gematcht; das erste
// passende Muster gewinnt pro Vorkommen. Linkpflichtfrei sind
// Heading-Zeilen (Struktur, kein Fließtext) und Vorkommen im
// deklarierten Target des Musters (Definitions-Ort;
// spec/spezifikation.md §DC-FA-ID-001.a).
func checkIDs(file string, content []byte, lines []Line, patterns []IDPattern) []Finding {
	if len(patterns) == 0 {
		return nil
	}
	inTarget := make([]bool, len(patterns))
	for i, p := range patterns {
		inTarget[i] = fileInTarget(file, p.Target)
	}
	var findings []Finding
	for _, ln := range lines {
		if _, _, isHeading := parseATXHeading(ln.Text); isHeading {
			continue // Headings sind linkpflichtfrei
		}
		findings = append(findings, checkIDLine(file, ln, patterns, inTarget)...)
	}
	// link-policy: always — zusätzlich Vorkommen INNERHALB von
	// Inline-Code-Spans (DC-FA-ID-001.a). Additiv zur prose-Prüfung;
	// Code-Span- und Fließtext-Bereiche sind disjunkt.
	findings = append(findings, checkIDsAlways(file, content, patterns, inTarget)...)
	return findings
}

// checkIDsAlways prüft Kennungs-Vorkommen innerhalb von Inline-Code-Spans
// für Muster mit LinkPolicy "always" (spec/spezifikation.md
// §DC-FA-ID-001.a). Arbeitet — wie codepaths — auf den rohen
// Prosa-Zeilen, weil die übrige Vorverarbeitung Inline-Code gerade
// leert. Ein Code-Span-Vorkommen ist linkpflichtfrei, wenn es im
// Linktext liegt, im target des Musters steht, die Datei ein
// exempt-paths-Glob matcht, die Zeile d-check:ignore trägt oder es eine
// Heading-Zeile ist.
func checkIDsAlways(file string, content []byte, patterns []IDPattern, inTarget []bool) []Finding {
	var active []int
	for i, p := range patterns {
		if p.LinkPolicy == AlwaysPolicy {
			active = append(active, i)
		}
	}
	if len(active) == 0 {
		return nil
	}
	prose := proseLines(content)
	spans := inlineSpansByLine(prose)
	var findings []Finding
	for _, pl := range prose {
		findings = append(findings, alwaysLineFindings(file, pl, spans[pl.no], patterns, active, inTarget)...)
	}
	return findings
}

// alwaysLineFindings prüft die Inline-Code-Spans einer Prosa-Zeile gegen
// die always-Muster (Deklarationsreihenfolge = Präzedenz über das
// gemeinsame claimed). Linkpflichtfrei: Heading-Zeile, d-check:ignore.
func alwaysLineFindings(file string, pl proseLine, codeSpans []inlineSpan, patterns []IDPattern, active []int, inTarget []bool) []Finding {
	if len(codeSpans) == 0 || strings.Contains(pl.raw, ignoreMarker) {
		return nil
	}
	if _, _, ok := parseATXHeading(pl.raw); ok {
		return nil
	}
	linkSpans := ExtractLinkSpans(pl.raw)
	var claimed [][2]int
	var findings []Finding
	for _, pi := range active {
		findings = append(findings,
			alwaysPatternFindings(file, pl, codeSpans, patterns[pi], inTarget[pi], linkSpans, &claimed)...)
	}
	return findings
}

// alwaysPatternFindings prüft ein einzelnes always-Muster gegen alle
// Code-Span-Werte der Zeile; ein Vorkommen ist frei im target, in
// exempt-paths oder im Linktext.
func alwaysPatternFindings(file string, pl proseLine, codeSpans []inlineSpan, p IDPattern, inTgt bool, linkSpans []LinkSpan, claimed *[][2]int) []Finding {
	var findings []Finding
	for _, sp := range codeSpans {
		val := pl.raw[sp.valStart:sp.valEnd]
		for _, m := range p.Regex.FindAllStringIndex(val, -1) {
			start, end := sp.valStart+m[0], sp.valStart+m[1]
			if overlapsClaimed(*claimed, start, end) {
				continue // Vorkommen gehört einem früheren Muster
			}
			*claimed = append(*claimed, [2]int{start, end})
			if inTgt || ignored(file, p.ExemptPaths) || idOccurrenceExempt(linkSpans, start, end) {
				continue
			}
			findings = append(findings, Finding{
				File: file, Line: pl.no, Rule: "ids",
				Target: val[m[0]:m[1]], Reason: ReasonIDUnlinked,
				Message: "Kennung ohne Link auf ihre Definition",
			})
		}
	}
	return findings
}

// checkIDLine prüft eine Fließtext-Zeile gegen alle Muster
// (Deklarationsreihenfolge; überlappende Vorkommen gehören dem
// früheren Muster).
func checkIDLine(file string, ln Line, patterns []IDPattern, inTarget []bool) []Finding {
	spans := ExtractLinkSpans(ln.Text)
	var claimed [][2]int
	var findings []Finding
	for pi, p := range patterns {
		for _, m := range p.Regex.FindAllStringIndex(ln.Text, -1) {
			if overlapsClaimed(claimed, m[0], m[1]) {
				continue // Vorkommen gehört einem früheren Muster
			}
			claimed = append(claimed, [2]int{m[0], m[1]})
			if inTarget[pi] || idOccurrenceExempt(spans, m[0], m[1]) {
				continue // Definitions-Ort bzw. verlinkt/kein Fließtext
			}
			findings = append(findings, Finding{
				File: file, Line: ln.No, Rule: "ids",
				Target: ln.Text[m[0]:m[1]], Reason: ReasonIDUnlinked,
				Message: "Kennung ohne Link auf ihre Definition",
			})
		}
	}
	return findings
}

// fileInTarget: liegt die geprüfte Datei im deklarierten Target des
// Musters (die Target-Datei selbst bzw. unterhalb des
// Target-Verzeichnisses)? Dort ist die Kennung „zu Hause" — die
// Definitions-Stelle muss nicht auf sich selbst verlinken.
func fileInTarget(file, target string) bool {
	rel, escaped := resolveConfigPath(target)
	if escaped {
		return false
	}
	if rel == "" {
		return true // Repo-Wurzel als Target: alles ist Definitions-Ort
	}
	return file == rel || strings.HasPrefix(file, rel+"/")
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
