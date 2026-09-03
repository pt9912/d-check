package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// snapshot liest jede Datei unter root in eine Pfad->Inhalt-Map -- Grundlage
// der Umkehr-Probe "ohne -apply wird nichts geschrieben".
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		out[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TestRun_DryRun_NoWrites ist die Umkehr-Probe zum DoD-Punkt "ohne -apply
// schreibt das Werkzeug nichts": eine Mutation, die apply faelschlich
// schreiben liesse (z. B. ein vertauschtes `if !apply` in run), macht diesen
// Test rot, weil der Baum-Vergleich danach nicht mehr leer waere.
func TestRun_DryRun_NoWrites(t *testing.T) {
	root := buildFixture(t)
	before := snapshot(t, root)

	if err := run(root, "welle-99", false); err != nil {
		t.Fatalf("Dry-Run schlug fehl: %v", err)
	}

	after := snapshot(t, root)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("Dry-Run hat den Baum veraendert -- ohne -apply darf nichts geschrieben werden")
	}
}

func TestPreviewMoves(t *testing.T) {
	root := buildFixture(t)
	wellePlan, err := FindWellePlan(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	slices, err := CollectSlices(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	p := Plan{WelleID: "welle-99", WellePlan: wellePlan, Slices: slices}

	moves := previewMoves(root, p)
	if len(moves) != 3 { // Welle-Plan + 2 Slices
		t.Fatalf("erwartet 3 Moves, got %d: %v", len(moves), moves)
	}
	want := map[string]string{
		"docs/plan/planning/done/welle-99-testwelle.md": "docs/plan/planning/done/welle-99/welle-99-testwelle.md",
		"docs/plan/planning/done/slice-501-eins.md":     "docs/plan/planning/done/welle-99/slice-501-eins.md",
		"docs/plan/planning/done/slice-502-zwei.md":     "docs/plan/planning/done/welle-99/slice-502-zwei.md",
	}
	for _, m := range moves {
		if want[m.Old] != m.New {
			t.Errorf("Move %s -> %s, erwartet -> %s", m.Old, m.New, want[m.Old])
		}
	}
}

// TestPreviewRewrites_OnlyExistingFiles belegt die dokumentierte Grenze
// zwischen Vorschau und Anwendung: der Stub von slice-502 existiert im
// Dry-Run noch nicht, sein eigener (bei -apply nachgezogener) Selbstverweis
// zaehlt hier also nicht mit -- nur der externe Beleg in observations.md.
func TestPreviewRewrites_OnlyExistingFiles(t *testing.T) {
	root := buildFixture(t)
	wellePlan, err := FindWellePlan(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	slices, err := CollectSlices(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	p := Plan{WelleID: "welle-99", WellePlan: wellePlan, Slices: slices}
	moves := previewMoves(root, p)

	hits, err := PreviewRewrites(root, moves)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"docs/plan/planning/observations.md": 1}
	if !reflect.DeepEqual(hits, want) {
		t.Fatalf("got %v, want %v", hits, want)
	}
}

// TestDropReviewSelfHits belegt eine an welle-60 gemessene Eigenheit: ein
// Review-Report kann selbst auf einen mit-gesammelten Slice verweisen, wird
// bei -apply aber ohne Stub geloescht, bevor der Verweis-Nachzug laeuft --
// die Vorschau darf diesen Treffer nicht als tatsaechlich eintretend zeigen.
func TestDropReviewSelfHits(t *testing.T) {
	root := "/repo"
	hits := map[string]int{
		"docs/plan/planning/observations.md": 1,
		"docs/reviews/r1.md":                 2,
	}
	dropReviewSelfHits(hits, root, []string{filepath.Join(root, "docs/reviews/r1.md")})
	want := map[string]int{"docs/plan/planning/observations.md": 1}
	if !reflect.DeepEqual(hits, want) {
		t.Fatalf("got %v, want %v", hits, want)
	}
}

// TestRun_PlanOhneEigeneSlices belegt die an welle-73 gemessene Eigenheit:
// ein Welle-Plan ohne eigene Slices (ihr Closure-Trigger ist ein Slice
// einer ANDEREN Welle) ist legitim und archiviert den Plan allein -- nur
// wenn WEDER Plan NOCH Slices existieren, ist es der Tippfehler-Fall.
func TestRun_PlanOhneEigeneSlices(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-73-x.md"),
		"# Welle welle-73: X\n\nLiefert nichts Eigenes.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-73-results.md"),
		"# Ergebnis welle-73\n")
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")

	if err := run(root, "welle-73", true); err != nil {
		t.Fatalf("erwartet Erfolg (Plan ohne eigene Slices ist legitim), got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/welle-73/archiv.zip")); err != nil {
		t.Errorf("archiv.zip fehlt: %v", err)
	}
}

// TestRun_WederPlanNochSlices belegt die Gegenprobe: existiert die Welle
// gar nicht (Tippfehler), bleibt der Fail-Closed-Guard scharf.
func TestRun_WederPlanNochSlices(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/.gitkeep"), "")

	if err := run(root, "welle-404", true); err == nil {
		t.Fatal("erwartet Fehler ohne Plan und ohne Slices")
	}
}
