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

var fencedBlockRE = regexp.MustCompile("(?s)```.*?```")
var inlineCodeRE = regexp.MustCompile("`[^`\n]*`")

// stripStandaloneInlineCode entfernt Inline-Code-Spannen, LASSEN aber die
// weitaus haeufigere, echte Zitat-Form dieses Repos stehen: eine Kennung als
// Markdown-Link-Label, "[`DC-FA-XXX-001`](ziel)" -- der Span endet dort
// unmittelbar vor "](", nicht vor Prosa. Go's RE2 kennt kein Lookahead,
// deshalb der manuelle Nachbar-Zeichen-Check statt eines Regex-Ausschlusses.
func stripStandaloneInlineCode(content string) string {
	idx := inlineCodeRE.FindAllStringIndex(content, -1)
	if idx == nil {
		return content
	}
	var b strings.Builder
	last := 0
	for _, m := range idx {
		start, end := m[0], m[1]
		b.WriteString(content[last:start])
		if end+1 < len(content) && content[end] == ']' && content[end+1] == '(' {
			b.WriteString(content[start:end]) // Link-Label erhalten
		}
		last = end
	}
	b.WriteString(content[last:])
	return b.String()
}

// ExtractSurvivingIDs liest die DC-*/ADR-*-Kennungen aus einem Slice-Volltext
// -- die "Kennungen, die den Vorgang ueberlebt haben" aus dem Baseline-Stub-
// Template. Gemessen: ohne das wird eine Anforderung, deren EINZIGER
// zitierender Slice archiviert wird, zur Trace-Waise (der Stub traegt den
// Original-Text nicht mehr, nur d-checks eigener --require-complete-Lauf hat
// das ans Licht gebracht, kein anderes Gate). Sortiert, dedupliziert.
//
// Fenced Blocks werden vollstaendig entfernt; ein reiner Inline-Code-Span
// (kein Link-Label) ebenso -- ein Slice kann eine erfundene, aehnlich
// geformte Kennung als Illustration eines Parsing-Grenzfalls zeigen
// (gemessen an slice-075: "`GG-QA-001, 007 Sekunden`"), ohne Ausschluss
// haette der Stub eine Anforderung behauptet, die es nie gab. Ein
// Link-Label ("[`DC-FA-XXX-001`](ziel)") bleibt dagegen erhalten -- das ist
// die weitaus haeufigere, echte Zitat-Form dieses Repos; ein blindes
// Wegwerfen aller Inline-Code-Spannen haette echte Belege mitgeloescht
// (gemessen: erste Fassung dieser Funktion strich 3 von 4 echten Kennungen
// in slice-073).
func ExtractSurvivingIDs(content string) []string {
	content = fencedBlockRE.ReplaceAllString(content, "")
	content = stripStandaloneInlineCode(content)
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

// SliceStubStandalone erzeugt den Slice-Stub-Text fuer den wellenlosen
// Einzel-Slice-Modus: Stub und Archiv liegen zusammen in
// `WellenlosArchiveDir` (archive.go), und es gibt keine einsammelnde Welle --
// an die Stelle von "Archiviert mit: <welle-id>" tritt "Archiviert: <Datum>
// (eigene Closure)". Dasselbe vendorte Template
// (`archiv-stub-slice.template.md`) beschreibt beide Faelle: sein
// "Archiviert mit"-Feld benennt die *Einsammlung*, hier eben durch die
// eigene Closure statt durch eine Welle -- keine neue Norm, nur ein neuer
// Feldwert. Der d-check:ignore-Marker auf der Hervorgegangen-Zeile exempt sie
// von der Linkpflicht: die Kennungen im archivierten Original-Volltext waren
// zum Teil verlinkt, der Stub traegt nur noch den blossen Text.
func SliceStubStandalone(id, title, welleField, hervorgegangen string) string {
	return fmt.Sprintf(
		"# %s — %s\n\n"+
			"> **ARCHIVIERT** — Volltext:\n"+
			"> `unzip -p done/wellenlos/%s-archiv.zip <pfad-im-archiv>`\n\n"+
			"**Welle:** %s\n"+
			"**Archiviert:** %s (eigene Closure)\n"+
			"**Hervorgegangen:** %s <!-- d-check:ignore (Kennungen aus dem archivierten Volltext, absichtlich unverlinkt) -->\n",
		id, title, id, welleField, placeholder, hervorgegangen,
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
