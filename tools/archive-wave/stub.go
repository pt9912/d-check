// Erzeugt die gekuerzten Stub-Texte nach den beiden Baseline-Templates
// (archiv-stub-slice.template.md, archiv-stub-welle.template.md). Ein Stub
// traegt keine Abschnittsueberschriften -- das ist die Kuerzung, die die
// Templates form-pruefbar machen.
package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var h1RE = regexp.MustCompile(`(?m)^#\s+(?:Slice\s+|Welle\s+)?[\w-]+:?\s*(.*)$`)

// ExtractTitle liest den Titel aus der ersten Ueberschrift eines Slice- oder
// Welle-Plans ("# Slice slice-190: Titel" bzw. "# Welle welle-87: Titel").
// Liefert einen leeren String, wenn keine passende Ueberschrift gefunden
// wird -- der Aufrufer entscheidet, ob das ein Fehler ist.
func ExtractTitle(content string) string {
	m := h1RE.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

// survivingIDRE erfasst genau die beiden Kennungs-Formen, die das Produkt
// selbst fuer Requirements-Traceability liest (trace.requirements.id-pattern
// aus .d-check.yml) und die ADR-Form (trace.adrs.file-pattern) -- keine
// erfundene, eigene Klasse.
var survivingIDRE = regexp.MustCompile(`\b(?:[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?|ADR-\d{4})\b`)

// ExtractSurvivingIDs liest die DC-*/ADR-*-Kennungen aus einem Slice-Volltext
// -- die "Kennungen, die den Vorgang ueberlebt haben" aus dem Baseline-Stub-
// Template. Gemessen: ohne das wird eine Anforderung, deren EINZIGER
// zitierender Slice archiviert wird, zur Trace-Waise (der Stub traegt den
// Original-Text nicht mehr, nur d-checks eigener --require-complete-Lauf hat
// das ans Licht gebracht, kein anderes Gate). Sortiert, dedupliziert.
func ExtractSurvivingIDs(content string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range survivingIDRE.FindAllString(content, -1) {
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	sort.Strings(out)
	return out
}

// FormatHervorgegangen liefert den Feldwert nach Template-Vorgabe: die
// gefundenen Kennungen kommagetrennt, sonst "— keine —" (der Template-eigene
// Leerwert, nicht der generische Platzhalter).
func FormatHervorgegangen(ids []string) string {
	if len(ids) == 0 {
		return "— keine —"
	}
	return strings.Join(ids, ", ")
}

const placeholder = "<manuell auszufuellen>"

// SliceStub erzeugt den Slice-Stub-Text. welleField ist der Wert des
// urspruenglichen **Welle:**-Feldes (z. B. "welle-87" oder "ohne Welle") --
// er bleibt unveraendert, waehrend "Archiviert mit" die einsammelnde Welle
// nennt (fuer wellenlose Slices sind das zwei verschiedene Tatsachen).
// hervorgegangen ist das fertig formatierte Feld (FormatHervorgegangen).
func SliceStub(id, title, welleField, welleID, hervorgegangen string) string {
	return fmt.Sprintf(
		"# %s — %s\n\n"+
			"> **ARCHIVIERT** — Volltext:\n"+
			"> `unzip -p done/%s/archiv.zip <pfad-im-archiv>`\n\n"+
			"**Welle:** %s\n"+
			"**Archiviert mit:** %s · **Geschlossen:** %s\n"+
			"**Hervorgegangen:** %s\n",
		id, title, welleID, welleField, welleID, placeholder, hervorgegangen,
	)
}

// WelleStub erzeugt den Welle-Plan-Stub-Text. nSlices/nReviews sind gezaehlt,
// nicht geschaetzt -- die Zahl, gegen die sich die Vollstaendigkeit des
// Archivs abzaehlen laesst.
func WelleStub(id, title, resultsFile string, nSlices, nReviews int) string {
	return fmt.Sprintf(
		"# %s — %s\n\n"+
			"> **ARCHIVIERT** — Volltext:\n"+
			"> `unzip -p done/%s/archiv.zip <pfad-im-archiv>`\n\n"+
			"**Geschlossen:** %s · **Ergebnisnotiz:** %s\n"+
			"**Archivierte Vorgänge:** %d Slices, %d Reviews\n",
		id, title, id, placeholder, resultsFile, nSlices, nReviews,
	)
}
