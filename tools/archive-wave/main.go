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
	"path/filepath"
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
	if wellePlan == "" {
		fmt.Println("  Welle-Plan: (keiner -- Welle vor der Plan-Datei-Konvention, nur -results.md)")
	} else {
		fmt.Printf("  Welle-Plan: %s\n", RelPath(root, wellePlan))
	}
	fmt.Printf("  Slices (%d):\n", len(slices))
	for _, s := range slices {
		fmt.Printf("    %s\n", RelPath(root, s))
	}
	fmt.Printf("  Review-Reports (%d):\n", len(reviews))
	for _, r := range reviews {
		fmt.Printf("    %s\n", RelPath(root, r))
	}

	// Fail-closed nur, wenn BEIDE Signale fehlen -- das ist der Tippfehler-
	// Fall (Welle existiert gar nicht). Ein Welle-Plan ohne eigene Slices ist
	// dagegen legitim (gemessen an welle-73: ihr Closure-Trigger ist ein
	// Slice aus welle-69, sie liefert selbst keinen): dann bleibt etwas zu
	// archivieren -- der Plan.
	if len(slices) == 0 && wellePlan == "" {
		return fmt.Errorf("weder Welle-Plan noch Slices fuer %s gefunden -- Abbruch, fail-closed statt leeres Archiv", welleID)
	}

	p := Plan{WelleID: welleID, WellePlan: wellePlan, Slices: slices, Reviews: reviews}

	if !apply {
		moves := previewMoves(root, p)
		hits, err := PreviewRewrites(root, moves)
		if err != nil {
			return err
		}
		// Review-Reports bekommen bei -apply keinen Stub -- sie sind zum
		// Zeitpunkt des Verweis-Nachzugs bereits geloescht und koennen selbst
		// nie eine nachgezogene Referenz tragen. Ohne diesen Filter zeigte die
		// Vorschau einen Treffer, der bei -apply nie eintritt (gemessen an
		// welle-60: Review-Reports verweisen auf andere gesammelte Slices).
		dropReviewSelfHits(hits, root, p.Reviews)
		fmt.Println("  Geplante Verweis-Fixes (ohne -apply wird nichts geschrieben):")
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
// fuer die Verweis-Vorschau ohne -apply.
func previewMoves(root string, p Plan) []Move {
	archiveDir := "docs/plan/planning/done/" + p.WelleID
	var moves []Move
	if p.WellePlan != "" {
		moves = append(moves, Move{
			Old: RelPath(root, p.WellePlan),
			New: archiveDir + "/" + filepath.Base(p.WellePlan),
		})
	}
	for _, s := range p.Slices {
		moves = append(moves, Move{
			Old: RelPath(root, s),
			New: archiveDir + "/" + filepath.Base(s),
		})
	}
	return moves
}

// dropReviewSelfHits entfernt Review-Report-Pfade aus einer Treffer-Map --
// sie werden bei -apply ohne Stub geloescht, bevor der Verweis-Nachzug
// laeuft, und koennen deshalb selbst nie eine nachgezogene Referenz tragen.
func dropReviewSelfHits(hits map[string]int, root string, reviews []string) {
	for _, r := range reviews {
		delete(hits, RelPath(root, r))
	}
}
