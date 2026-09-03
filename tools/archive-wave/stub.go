// Erzeugt die gekuerzten Stub-Texte nach den beiden Baseline-Templates
// (archiv-stub-slice.template.md, archiv-stub-welle.template.md). Ein Stub
// traegt keine Abschnittsueberschriften -- das ist die Kuerzung, die die
// Templates form-pruefbar machen.
package main

import (
	"fmt"
	"regexp"
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

const placeholder = "<manuell auszufuellen>"

// SliceStub erzeugt den Slice-Stub-Text. welleField ist der Wert des
// urspruenglichen **Welle:**-Feldes (z. B. "welle-87" oder "ohne Welle") --
// er bleibt unveraendert, waehrend "Archiviert mit" die einsammelnde Welle
// nennt (fuer wellenlose Slices sind das zwei verschiedene Tatsachen).
func SliceStub(id, title, welleField, welleID string) string {
	return fmt.Sprintf(
		"# %s — %s\n\n"+
			"> **ARCHIVIERT** — Volltext:\n"+
			"> `unzip -p done/%s/archiv.zip <pfad-im-archiv>`\n\n"+
			"**Welle:** %s\n"+
			"**Archiviert mit:** %s · **Geschlossen:** %s\n"+
			"**Hervorgegangen:** %s\n",
		id, title, welleID, welleField, welleID, placeholder, placeholder,
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
