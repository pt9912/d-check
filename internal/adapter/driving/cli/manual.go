package cli

import (
	"fmt"
	"io"
	"strings"

	dcheck "github.com/pt9912/d-check"
	"github.com/pt9912/d-check/internal/hexagon/core/app"
)

// manualComboError prüft --manual (DC-FA-CLI-013.a Schritt 1): ein leerer
// Begriff und jede Kombination mit einem anderen Modus sind Nutzungsfehler.
func manualComboError(o options) string {
	if !o.manualSet {
		return ""
	}
	if strings.TrimSpace(o.manual) == "" {
		return "--manual braucht einen Begriff"
	}
	if o.json || o.yaml || o.doctor || o.repair || o.trace || o.requireComplete || o.printConfig ||
		o.printMK || o.suggestConfig != "" || o.commitMsg != "" || o.vcsRange != "" || o.vcsStaged {
		return "--manual ist ein eigener Modus (mit keiner anderen Modus-Option kombinierbar)"
	}
	return ""
}

// runManual gibt die Abschnitte der mitgelieferten Dokumente aus, deren
// Überschrift den Begriff nennt (DC-FA-CLI-013.a Schritte 2–6) — repo-frei
// und netzlos; kein Treffer ⇒ Exit 2 mit den Dokument-Titeln als Hinweis.
func runManual(term string, stdout, stderr io.Writer) int {
	docs := []app.ManualDoc{
		{Path: dcheck.HandbuchPfad, Content: dcheck.Handbuch},
		{Path: dcheck.SpezifikationPfad, Content: dcheck.Spezifikation},
	}
	out, found := app.ManualSections(docs, term)
	if !found {
		fmt.Fprintf(stderr, "d-check: error: --manual: keine Überschrift nennt %q — der Titel eines Dokuments liefert es ganz: %q, %q\n",
			strings.TrimSpace(term), titel(dcheck.Handbuch), titel(dcheck.Spezifikation))
		return 2
	}
	fmt.Fprint(stdout, out)
	return 0
}

// titel ist der Text der ersten Überschrift erster Ebene eines Dokuments.
func titel(doc string) string {
	for _, l := range strings.Split(doc, "\n") {
		if strings.HasPrefix(l, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(l, "# "))
		}
	}
	return ""
}
