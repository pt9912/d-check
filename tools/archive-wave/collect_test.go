package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFindWellePlan(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-87-archivierung.md"), "# Welle welle-87: X\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-87-results.md"), "# Ergebnis\n")

	got, err := FindWellePlan(root, "welle-87")
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if filepath.Base(got) != "welle-87-archivierung.md" {
		t.Fatalf("erwartet den Plan, nicht die Ergebnisnotiz, got %s", got)
	}
}

func TestFindWellePlan_Mehrdeutig(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-87-a.md"), "x")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-87-b.md"), "x")

	if _, err := FindWellePlan(root, "welle-87"); err == nil {
		t.Fatal("erwartet Fehler bei mehrdeutigem Bestand")
	}
}

func TestFindWellePlan_Keine(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/.gitkeep"), "")

	if _, err := FindWellePlan(root, "welle-87"); err == nil {
		t.Fatal("erwartet Fehler ohne Kandidat")
	}
}

// TestCollectSlices ist zugleich die Umkehr-Probe fuer die Sammeln-Bedingung
// (BEO-023): ein Slice mit ABWEICHENDER Welle (welle-86) und ein wellenloser
// Slice muessen draussen bleiben -- dreht man den Vergleich in
// CollectSlices auf "IMMER wahr", wird genau dieser Test rot.
func TestCollectSlices(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-190-a.md"),
		"# Slice slice-190: A\n\n**Welle:** [welle-87](../welle-87-x.md).\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-191-b.md"),
		"# Slice slice-191: B\n\n**Welle:** welle-87.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-999-fremd.md"),
		"# Slice slice-999: Fremd\n\n**Welle:** welle-86.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-998-wellenlos.md"),
		"# Slice slice-998: Wellenlos\n\n**Welle:** — wellenlos.\n")

	got, err := CollectSlices(root, "welle-87")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("erwartet 2 Slices, got %d: %v", len(got), got)
	}
	for _, p := range got {
		base := filepath.Base(p)
		if base != "slice-190-a.md" && base != "slice-191-b.md" {
			t.Errorf("unerwarteter Treffer: %s", base)
		}
	}
}

func TestCollectReviews(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-01-slice-190-r1.md"), "x")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-01-slice-190-r2.md"), "x")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-09-01-slice-999-fremd.md"), "x")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-08-01-cr-ohne-slice.md"), "x")

	got, err := CollectReviews(root, []string{filepath.Join(root, "docs/plan/planning/done/slice-190-a.md")})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("erwartet 2 Reviews (1:N zulaessig), got %d: %v", len(got), got)
	}
}

func TestReadWelleField(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "slice.md")
	writeFile(t, path, "# X\n\n**Welle:** [welle-87](../welle-87-x.md).\n")

	got, err := ReadWelleField(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "[welle-87](../welle-87-x.md)."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSliceIDFromPath(t *testing.T) {
	if got := SliceIDFromPath("/a/b/slice-190-titel.md"); got != "slice-190" {
		t.Fatalf("got %q", got)
	}
	if got := SliceIDFromPath("/a/b/welle-87.md"); got != "" {
		t.Fatalf("erwartet leeren String ohne slice-Kennung, got %q", got)
	}
}
