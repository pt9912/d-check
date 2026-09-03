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
