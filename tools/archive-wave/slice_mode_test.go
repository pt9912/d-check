package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildSliceFixture legt einen wellenlosen Slice mit zwei Review-Reports und
// einer externen Referenz an -- die konstruierte Test-Grundlage fuer den
// Einzel-Slice-Modus, nicht der echte Bestand.
func buildSliceFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-601-eins.md"),
		"# Slice slice-601: Eins\n\n**Welle:** — wellenlos.\n\nInhalt.\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-04-slice-601-r1.md"),
		"# Review slice-601 r1\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-04-slice-601-verifikation.md"),
		"# Verifikation slice-601\n")
	// Fremd-Slice, gehoert NICHT zum Ziel -- Gegenprobe.
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-602-fremd.md"),
		"# Slice slice-602: Fremd\n\n**Welle:** — wellenlos.\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-04-slice-602-r1.md"),
		"# Review slice-602\n")
	// Externer Verweis auf einen der beiden geloeschten Review-Reports --
	// fuer den Dangling-Check-Teil des Tests.
	writeFile(t, filepath.Join(root, "docs/plan/planning/in-progress/roadmap.md"),
		"Beleg: [r1](../../../reviews/2026-09-04-slice-601-r1.md).\n")
	return root
}

func TestRunSlice_DryRun_NoWrites(t *testing.T) {
	root := buildSliceFixture(t)
	before := snapshot(t, root)

	if err := runSlice(root, "slice-601", false); err != nil {
		t.Fatalf("Dry-Run schlug fehl: %v", err)
	}

	after := snapshot(t, root)
	for k, v := range before {
		if after[k] != v {
			t.Fatalf("Dry-Run hat %s veraendert -- ohne -apply darf nichts geschrieben werden", k)
		}
	}
}

func TestRunSlice_Apply(t *testing.T) {
	root := buildSliceFixture(t)

	if err := runSlice(root, "slice-601", true); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	// Die Slice-Datei bleibt an ihrem Pfad -- kein Move, nur Inhaltsersatz.
	stubB, err := os.ReadFile(filepath.Join(root, "docs/plan/planning/done/slice-601-eins.md"))
	if err != nil {
		t.Fatalf("Stub fehlt am unveraenderten Pfad: %v", err)
	}
	stub := string(stubB)
	if strings.Contains(stub, "Inhalt.") {
		t.Fatalf("Stub traegt noch den Volltext: %q", stub)
	}
	if !strings.Contains(stub, "unzip -p done/slice-601-archiv.zip") {
		t.Fatalf("Stub nennt nicht den flachen Archiv-Pfad: %q", stub)
	}
	if !strings.Contains(stub, "**Welle:** — wellenlos.") {
		t.Fatalf("Stub hat das urspruengliche Welle-Feld nicht uebernommen: %q", stub)
	}
	if !strings.Contains(stub, "eigene Closure") {
		t.Fatalf("Stub nennt nicht die eigene Closure als Einsammlung: %q", stub)
	}

	// Archiv liegt flach neben dem Stub, nicht in einem Unterverzeichnis.
	zipPath := filepath.Join(root, "docs/plan/planning/done/slice-601-archiv.zip")
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("archiv.zip fehlt oder unlesbar: %v", err)
	}
	defer zr.Close()
	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	want := []string{
		"docs/plan/planning/done/slice-601-eins.md",
		"docs/reviews/2026-09-04-slice-601-r1.md",
		"docs/reviews/2026-09-04-slice-601-verifikation.md",
	}
	for _, w := range want {
		if !names[w] {
			t.Errorf("archiv.zip fehlt Eintrag %s (hat: %v)", w, names)
		}
	}

	// Review-Reports sind geloescht, ohne Stub.
	for _, p := range []string{
		"docs/reviews/2026-09-04-slice-601-r1.md",
		"docs/reviews/2026-09-04-slice-601-verifikation.md",
	} {
		if _, err := os.Stat(filepath.Join(root, p)); !os.IsNotExist(err) {
			t.Errorf("%s haette geloescht sein muessen", p)
		}
	}

	// Der Fremd-Slice und sein Review bleiben unangetastet.
	if _, err := os.Stat(filepath.Join(root, "docs/reviews/2026-09-04-slice-602-r1.md")); err != nil {
		t.Errorf("Fremd-Review haette unangetastet bleiben muessen: %v", err)
	}
}

// TestRunSlice_RejectsWelleSlice belegt: ein Slice mit gesetztem, echtem
// Welle-Feld gehoert in den -welle-Modus und wird hier abgelehnt statt
// still falsch archiviert.
func TestRunSlice_RejectsWelleSlice(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-603-mit-welle.md"),
		"# Slice slice-603: Mit Welle\n\n**Welle:** welle-42.\n")

	if err := runSlice(root, "slice-603", true); err == nil {
		t.Fatal("erwartet Fehler fuer einen Slice mit echter Wellen-Zugehoerigkeit")
	}
}

// TestRunSlice_DanglingReviewReference belegt den Melde-Teil (kein
// automatischer Nachzug moeglich, da kein Move-Ziel existiert): ein
// externer Verweis auf einen geloeschten Review-Report wird als toter
// Verweis erkannt.
func TestRunSlice_DanglingReviewReference(t *testing.T) {
	root := buildSliceFixture(t)
	slicePath, err := FindSlice(root, "slice-601")
	if err != nil {
		t.Fatal(err)
	}
	reviews, err := CollectReviews(root, []string{slicePath})
	if err != nil {
		t.Fatal(err)
	}
	hits, err := FindReferencesToPaths(root, reviews, append([]string{slicePath}, reviews...))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"docs/plan/planning/in-progress/roadmap.md": 1}
	if len(hits) != len(want) || hits["docs/plan/planning/in-progress/roadmap.md"] != 1 {
		t.Fatalf("got %v, want %v", hits, want)
	}
}

// TestFindReferencesToPaths_ExcludesSelfAndReviews belegt die exclude-Menge:
// ein Review-Report, der auf einen ANDEREN, ebenfalls zu loeschenden Review
// verweist (dieselbe Eigenheit, die runWelle als "dropReviewSelfHits" fuer
// den Wellen-Modus behandelt), darf sich nicht selbst als toten Verweis
// zaehlen. Ohne die exclude-Pruefung wuerde dieser Test rot.
func TestFindReferencesToPaths_ExcludesSelfAndReviews(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-604-eins.md"),
		"# Slice slice-604: Eins\n\n**Welle:** — wellenlos.\n")
	r1 := filepath.Join(root, "docs/reviews/2026-09-04-slice-604-r1.md")
	writeFile(t, r1, "# Review r1\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-04-slice-604-verifikation.md"),
		"# Verifikation\n\nBezieht sich auf [r1](2026-09-04-slice-604-r1.md).\n")

	reviews := []string{
		r1,
		filepath.Join(root, "docs/reviews/2026-09-04-slice-604-verifikation.md"),
	}
	slicePath := filepath.Join(root, "docs/plan/planning/done/slice-604-eins.md")
	hits, err := FindReferencesToPaths(root, reviews, append([]string{slicePath}, reviews...))
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 0 {
		t.Fatalf("erwartet keine Treffer (Verweis liegt in einer ausgeschlossenen Datei), got %v", hits)
	}
}

func TestValidateModeFlags(t *testing.T) {
	cases := []struct {
		welle, sliceID string
		wantErr        bool
	}{
		{"welle-87", "", false},
		{"", "slice-137", false},
		{"", "", true},
		{"welle-87", "slice-137", true},
	}
	for _, c := range cases {
		err := validateModeFlags(c.welle, c.sliceID)
		if (err != nil) != c.wantErr {
			t.Errorf("validateModeFlags(%q, %q): got err=%v, wantErr=%v", c.welle, c.sliceID, err, c.wantErr)
		}
	}
}
