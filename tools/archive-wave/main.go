// archive-wave setzt Modul 6, Schritt 4 des Planning-Harness-Regelwerks um
// ("Zeitdokumente der Welle archivieren"): sammelt die Slices und
// Review-Reports einer geschlossenen Welle, baut ein archiv.zip, ersetzt die
// Volltexte durch gekuerzte Stubs und zieht repo-weite Verweise nach.
//
// Portabel gehalten: kein Import aus einem d-check-internen Paket. Jedes
// Repo mit demselben Planning-Layout (docs/plan/planning/, docs/reviews/,
// **Welle:**-Feld) kann dieses Verzeichnis unveraendert uebernehmen.
//
// Sicherer Default: ohne -apply wird NICHTS geschrieben, nur der geplante
// Umfang ausgegeben.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	welle := flag.String("welle", "", "Wellen-Kennung, z. B. welle-87 (Pflicht)")
	root := flag.String("root", ".", "Repo-Wurzel")
	apply := flag.Bool("apply", false, "Aenderungen tatsaechlich schreiben (Default: nur anzeigen)")
	flag.Parse()

	if *welle == "" {
		fmt.Fprintln(os.Stderr, "archive-wave: -welle ist Pflicht (z. B. -welle=welle-87)")
		os.Exit(2)
	}

	if err := run(*root, *welle, *apply); err != nil {
		fmt.Fprintln(os.Stderr, "archive-wave: "+err.Error())
		os.Exit(1)
	}
}

func run(root, welleID string, apply bool) error {
	wellePlan, err := FindWellePlan(root, welleID)
	if err != nil {
		return err
	}
	slices, err := CollectSlices(root, welleID)
	if err != nil {
		return err
	}
	reviews, err := CollectReviews(root, slices)
	if err != nil {
		return err
	}

	fmt.Printf("archive-wave: %s\n", welleID)
	fmt.Printf("  Welle-Plan: %s\n", RelPath(root, wellePlan))
	fmt.Printf("  Slices (%d):\n", len(slices))
	for _, s := range slices {
		fmt.Printf("    %s\n", RelPath(root, s))
	}
	fmt.Printf("  Review-Reports (%d):\n", len(reviews))
	for _, r := range reviews {
		fmt.Printf("    %s\n", RelPath(root, r))
	}

	if len(slices) == 0 {
		return fmt.Errorf("keine Slices fuer %s gefunden -- Abbruch, fail-closed statt leeres Archiv", welleID)
	}

	p := Plan{WelleID: welleID, WellePlan: wellePlan, Slices: slices, Reviews: reviews}

	if !apply {
		moves := previewMoves(root, p)
		hits, err := PreviewRewrites(root, moves)
		if err != nil {
			return err
		}
		fmt.Println("  Geplante Verweis-Fixes (--apply schreibt nichts ohne dieses Flag):")
		if len(hits) == 0 {
			fmt.Println("    (keine)")
		}
		for _, f := range SortedKeys(hits) {
			fmt.Printf("    %s: %d\n", f, hits[f])
		}
		return nil
	}

	moves, err := Apply(root, p)
	if err != nil {
		return err
	}
	hits, err := RewriteRepo(root, moves)
	if err != nil {
		return err
	}
	fmt.Println("  Verweise nachgezogen:")
	for _, f := range SortedKeys(hits) {
		fmt.Printf("    %s: %d\n", f, hits[f])
	}
	fmt.Println("  Fertig -- git status pruefen, dann git mv + Commit wie im Plan vorgesehen.")
	return nil
}

// previewMoves berechnet dieselben Move-Ziele wie Apply, ohne zu schreiben --
// fuer den --dry-run-Verweis-Vorschau-Schritt.
func previewMoves(root string, p Plan) []Move {
	archiveDir := "docs/plan/planning/done/" + p.WelleID
	var moves []Move
	moves = append(moves, Move{
		Old: RelPath(root, p.WellePlan),
		New: archiveDir + "/" + baseName(p.WellePlan),
	})
	for _, s := range p.Slices {
		moves = append(moves, Move{
			Old: RelPath(root, s),
			New: archiveDir + "/" + baseName(s),
		})
	}
	return moves
}

func baseName(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}
