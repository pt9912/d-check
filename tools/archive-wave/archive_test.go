package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildFixture legt eine konstruierte Test-Welle mit zwei Slices und einem
// Review-Report an -- das im Slice-Plan verlangte Fixture, nicht der echte
// Bestand.
func buildFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-99-testwelle.md"),
		"# Welle welle-99: Testwelle\n\nInhalt des Plans.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-99-results.md"),
		"# Ergebnis welle-99\n\nBleibt unangetastet.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-501-eins.md"),
		"# Slice slice-501: Eins\n\n**Welle:** welle-99.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-502-zwei.md"),
		"# Slice slice-502: Zwei\n\n**Welle:** [welle-99](../welle-99-testwelle.md).\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-01-slice-501-r1.md"),
		"# Review slice-501\n")
	// Fremd-Slice, gehoert NICHT zur Testwelle -- Gegenprobe.
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-999-fremd.md"),
		"# Slice slice-999: Fremd\n\n**Welle:** welle-1.\n")
	// Verweis von aussen auf einen der beiden Slices -- fuer den
	// Nachzug-Teil des Fixture-Tests.
	writeFile(t, filepath.Join(root, "docs/plan/planning/observations.md"),
		"Beleg: [slice-501](done/slice-501-eins.md).\n")
	return root
}

func TestFixture_EndToEnd(t *testing.T) {
	root := buildFixture(t)

	wellePlan, err := FindWellePlan(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	slices, err := CollectSlices(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	if len(slices) != 2 {
		t.Fatalf("erwartet 2 Slices, got %d: %v", len(slices), slices)
	}
	reviews, err := CollectReviews(root, slices)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 {
		t.Fatalf("erwartet 1 Review, got %d: %v", len(reviews), reviews)
	}

	p := Plan{WelleID: "welle-99", WellePlan: wellePlan, Slices: slices, Reviews: reviews}
	moves, err := Apply(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(moves) != 3 { // Welle-Plan + 2 Slices
		t.Fatalf("erwartet 3 Moves, got %d: %v", len(moves), moves)
	}

	archiveDir := filepath.Join(root, "docs/plan/planning/done/welle-99")

	// Archiv vollstaendig? Umkehr-Probe: fehlt hier ein Eintrag (z. B. weil
	// buildZip vorzeitig abbricht), schlaegt dieser Block fehl.
	zr, err := zip.OpenReader(filepath.Join(archiveDir, "archiv.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	wantEntries := map[string]bool{
		"docs/plan/planning/done/welle-99-testwelle.md": false,
		"docs/plan/planning/done/slice-501-eins.md":     false,
		"docs/plan/planning/done/slice-502-zwei.md":     false,
		"docs/reviews/2026-09-01-slice-501-r1.md":       false,
	}
	if len(zr.File) != len(wantEntries) {
		t.Fatalf("erwartet %d Zip-Eintraege, got %d", len(wantEntries), len(zr.File))
	}
	for _, f := range zr.File {
		if _, ok := wantEntries[f.Name]; !ok {
			t.Errorf("unerwarteter Zip-Eintrag: %s", f.Name)
		}
		wantEntries[f.Name] = true
	}
	for name, seen := range wantEntries {
		if !seen {
			t.Errorf("Zip-Eintrag fehlt: %s", name)
		}
	}

	// Stubs am neuen Ort, alte Volltexte weg.
	if _, err := os.Stat(filepath.Join(archiveDir, "slice-501-eins.md")); err != nil {
		t.Errorf("Slice-Stub fehlt: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/slice-501-eins.md")); !os.IsNotExist(err) {
		t.Error("alter Slice-Volltext haette geloescht sein muessen")
	}
	if _, err := os.Stat(filepath.Join(root, "docs/reviews/2026-09-01-slice-501-r1.md")); !os.IsNotExist(err) {
		t.Error("Review-Report haette geloescht sein muessen (kein Stub fuer Reviews)")
	}
	if _, err := os.Stat(filepath.Join(archiveDir, "2026-09-01-slice-501-r1.md")); !os.IsNotExist(err) {
		t.Error("Review-Report haette KEINEN Stub bekommen duerfen")
	}

	// Ergebnisnotiz bleibt vollstaendig und flach.
	resultsPath := filepath.Join(root, "docs/plan/planning/done/welle-99-results.md")
	b, err := os.ReadFile(resultsPath)
	if err != nil {
		t.Fatalf("Ergebnisnotiz wurde angefasst oder geloescht: %v", err)
	}
	if string(b) != "# Ergebnis welle-99\n\nBleibt unangetastet.\n" {
		t.Error("Ergebnisnotiz-Inhalt veraendert")
	}

	// Fremd-Slice unangetastet.
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/slice-999-fremd.md")); err != nil {
		t.Errorf("Fremd-Slice haette liegen bleiben muessen: %v", err)
	}

	// Verweis-Nachzug -- betrifft ZWEI Dateien: den externen Beleg UND den
	// frisch geschriebenen Stub von slice-502 selbst, dessen unveraendert
	// uebernommenes **Welle:**-Feld noch auf den alten Welle-Plan-Pfad zeigt.
	hits, err := RewriteRepo(root, moves)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 2 {
		t.Fatalf("erwartet 2 betroffene Dateien, got %d: %v", len(hits), hits)
	}
	stubB, err := os.ReadFile(filepath.Join(archiveDir, "slice-502-zwei.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stubB), "welle-99-testwelle.md)") {
		t.Fatalf("Stub-eigener Welle-Verweis wurde nicht nachgezogen: %q", string(stubB))
	}

	b, err = os.ReadFile(filepath.Join(root, "docs/plan/planning/observations.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "Beleg: [slice-501](done/welle-99/slice-501-eins.md).\n"
	if string(b) != want {
		t.Fatalf("got %q, want %q", string(b), want)
	}
}
