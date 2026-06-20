package rules

import (
	"strings"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// CheckIDs ist das Regelmodul `ids` (DC-FA-ID-001): nackte Kennungen
// im Fließtext — außerhalb von Inline-Code und Fenced-Code-Blöcken —
// müssen als Markdown-Link auf ihre Definition ausgeführt sein. Die
// Muster werden in Deklarationsreihenfolge gematcht; das erste
// passende Muster gewinnt pro Vorkommen. Linkpflichtfrei sind
// Heading-Zeilen (Struktur, kein Fließtext), Vorkommen im deklarierten
// Target des Musters (Definitions-Ort), Dateien aus den `exempt-paths`
// des Musters und Zeilen mit dem d-check:ignore-Marker — die beiden
// Ventile sind ein Ganzdatei- bzw. Ganzzeilen-Carve-out und gelten für
// nackte wie für Inline-Code-Vorkommen (spec/spezifikation.md
// §DC-FA-ID-001.a).
func CheckIDs(file string, content []byte, lines []Line, patterns []model.IDPattern) []model.Finding {
	if len(patterns) == 0 {
		return nil
	}
	// Datei-Invarianten je Muster einmal vorberechnen (wie inTarget):
	// Definitions-Ort und exempt-paths hängen nur von (Datei, Muster) ab,
	// nicht vom einzelnen Vorkommen.
	inTarget := make([]bool, len(patterns))
	exemptFile := make([]bool, len(patterns))
	for i, p := range patterns {
		inTarget[i] = fileInTarget(file, p.Target)
		exemptFile[i] = Ignored(file, p.ExemptPaths)
	}
	// Rohe Prosa-Zeilen einmal gewinnen; daraus die d-check:ignore-Zeilen.
	// Beide Prüfpfade (nackt + Inline-Code) konsultieren denselben Satz,
	// damit die Zeilen-Ausnahme nicht divergieren kann.
	prose := proseLines(content)
	ignoreLines := markerLines(prose)
	var findings []model.Finding
	for _, ln := range lines {
		if _, _, isHeading := parseATXHeading(ln.Text); isHeading {
			continue // Headings sind linkpflichtfrei
		}
		if ignoreLines[ln.No] {
			continue // d-check:ignore — Zeile von der ids-Prüfung frei
		}
		findings = append(findings, CheckIDLine(file, ln, patterns, inTarget, exemptFile)...)
	}
	// link-policy: always — zusätzlich Vorkommen INNERHALB von
	// Inline-Code-Spans (DC-FA-ID-001.a). Additiv zur prose-Prüfung;
	// Code-Span- und Fließtext-Bereiche sind disjunkt.
	findings = append(findings, checkIDsAlways(file, prose, ignoreLines, patterns, inTarget, exemptFile)...)
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
func checkIDsAlways(file string, prose []proseLine, ignoreLines map[int]bool, patterns []model.IDPattern, inTarget, exemptFile []bool) []model.Finding {
	var active []int
	for i, p := range patterns {
		if p.LinkPolicy == model.AlwaysPolicy {
			active = append(active, i)
		}
	}
	if len(active) == 0 {
		return nil
	}
	spans := inlineSpansByLine(prose)
	var findings []model.Finding
	for _, pl := range prose {
		findings = append(findings, alwaysLineFindings(file, pl, spans[pl.no], ignoreLines, patterns, active, inTarget, exemptFile)...)
	}
	return findings
}

// alwaysLineFindings prüft die Inline-Code-Spans einer Prosa-Zeile gegen
// die always-Muster (Deklarationsreihenfolge = Präzedenz über das
// gemeinsame claimed). Linkpflichtfrei: Heading-Zeile, d-check:ignore.
func alwaysLineFindings(file string, pl proseLine, codeSpans []inlineSpan, ignoreLines map[int]bool, patterns []model.IDPattern, active []int, inTarget, exemptFile []bool) []model.Finding {
	if len(codeSpans) == 0 || ignoreLines[pl.no] {
		return nil
	}
	if _, _, ok := parseATXHeading(pl.raw); ok {
		return nil
	}
	linkSpans := ExtractLinkSpans(pl.raw)
	var claimed [][2]int
	var findings []model.Finding
	for _, pi := range active {
		findings = append(findings,
			alwaysPatternFindings(file, pl, codeSpans, patterns[pi], inTarget[pi], exemptFile[pi], linkSpans, &claimed)...)
	}
	return findings
}

// alwaysPatternFindings prüft ein einzelnes always-Muster gegen alle
// Code-Span-Werte der Zeile; ein Vorkommen ist frei im target, in
// exempt-paths oder im Linktext.
func alwaysPatternFindings(file string, pl proseLine, codeSpans []inlineSpan, p model.IDPattern, inTgt, exempt bool, linkSpans []LinkSpan, claimed *[][2]int) []model.Finding {
	var findings []model.Finding
	for _, sp := range codeSpans {
		val := pl.raw[sp.valStart:sp.valEnd]
		for _, m := range p.Regex.FindAllStringIndex(val, -1) {
			start, end := sp.valStart+m[0], sp.valStart+m[1]
			if overlapsClaimed(*claimed, start, end) {
				continue // Vorkommen gehört einem früheren Muster
			}
			*claimed = append(*claimed, [2]int{start, end})
			if inTgt || exempt || IDOccurrenceExempt(linkSpans, start, end) {
				continue
			}
			findings = append(findings, model.Finding{
				File: file, Line: pl.no, Rule: "ids",
				Target: val[m[0]:m[1]], Reason: model.ReasonIDUnlinked,
				Message: "Kennung ohne Link auf ihre Definition",
			})
		}
	}
	return findings
}

// CheckIDLine prüft eine Fließtext-Zeile gegen alle Muster
// (Deklarationsreihenfolge; überlappende Vorkommen gehören dem
// früheren Muster).
func CheckIDLine(file string, ln Line, patterns []model.IDPattern, inTarget, exemptFile []bool) []model.Finding {
	spans := ExtractLinkSpans(ln.Text)
	var claimed [][2]int
	var findings []model.Finding
	for pi, p := range patterns {
		for _, m := range p.Regex.FindAllStringIndex(ln.Text, -1) {
			if overlapsClaimed(claimed, m[0], m[1]) {
				continue // Vorkommen gehört einem früheren Muster
			}
			claimed = append(claimed, [2]int{m[0], m[1]})
			if inTarget[pi] || exemptFile[pi] || IDOccurrenceExempt(spans, m[0], m[1]) {
				continue // Definitions-Ort, exempt-paths bzw. verlinkt/kein Fließtext
			}
			findings = append(findings, model.Finding{
				File: file, Line: ln.No, Rule: "ids",
				Target: ln.Text[m[0]:m[1]], Reason: model.ReasonIDUnlinked,
				Message: "Kennung ohne Link auf ihre Definition",
			})
		}
	}
	return findings
}

// markerLines liefert die 1-basierten Nummern der Prosa-Zeilen, die den
// d-check:ignore-Marker tragen — geprüft auf der rohen Zeile. Der eine
// zurückgegebene Satz wird von beiden ids-Prüfpfaden konsultiert (nackte
// Fließtext-Vorkommen in CheckIDs, Inline-Code-Vorkommen in
// alwaysLineFindings), damit die Zeilen-Ausnahme nicht divergieren kann
// (spec/spezifikation.md §DC-FA-ID-001.a, Ventil für nackte wie
// Inline-Code-Vorkommen).
func markerLines(prose []proseLine) map[int]bool {
	var out map[int]bool
	for _, pl := range prose {
		if strings.Contains(pl.raw, ignoreMarker) {
			if out == nil {
				out = make(map[int]bool)
			}
			out[pl.no] = true
		}
	}
	return out
}

// fileInTarget: liegt die geprüfte Datei im deklarierten Target des
// Musters (die Target-Datei selbst bzw. unterhalb des
// Target-Verzeichnisses)? Dort ist die Kennung „zu Hause" — die
// Definitions-Stelle muss nicht auf sich selbst verlinken.
func fileInTarget(file, target string) bool {
	rel, escaped := ResolveConfigPath(target)
	if escaped {
		return false
	}
	if rel == "" {
		return true // Repo-Wurzel als Target: alles ist Definitions-Ort
	}
	return file == rel || strings.HasPrefix(file, rel+"/")
}

// IDOccurrenceExempt prüft, ob ein Vorkommen [start,end) linkpflichtfrei
// ist: es liegt im Linktext eines Nicht-Bild-Links (= verlinkt) oder
// innerhalb der Link-/Bild-Syntax außerhalb des Linktexts
// (Ziel-Klammer, Alt-Text — kein Fließtext;
// spec/spezifikation.md §DC-FA-ID-001.a).
func IDOccurrenceExempt(spans []LinkSpan, start, end int) bool {
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
