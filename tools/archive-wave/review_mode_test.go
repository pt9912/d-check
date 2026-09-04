package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildReviewFixture legt einen eigenstaendigen Review (kein slice-<NNN> im
// Namen) und einen Fremd-Review mit slice-Bezug an -- Letzterer ist die
// Gegenprobe fuer die Zugehoerigkeits-Pruefung.
func buildReviewFixture(t *testing.T, heading string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/reviews/2026-06-11-cr-dc-fa-code-001.md"),
		heading+"\n\nInhalt des Reviews. Zitiert DC-FA-CODE-001 und ADR-0001.\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-08-25-slice-137-toolchain-freshness-review.md"),
		"# Review slice-137\n")
	return root
}

func TestRunReview_DryRun_NoWrites(t *testing.T) {
	root := buildReviewFixture(t, "# Review-Report: Change Request `DC-FA-CODE-001` — 2026-06-11")
	before := snapshot(t, root)

	if err := runReview(root, "2026-06-11-cr-dc-fa-code-001.md", false); err != nil {
		t.Fatalf("Dry-Run schlug fehl: %v", err)
	}

	after := snapshot(t, root)
	for k, v := range before {
		if after[k] != v {
			t.Fatalf("Dry-Run hat %s veraendert -- ohne -apply darf nichts geschrieben werden", k)
		}
	}
}

func TestRunReview_Apply(t *testing.T) {
	root := buildReviewFixture(t, "# Review-Report: Change Request `DC-FA-CODE-001` — 2026-06-11")

	if err := runReview(root, "2026-06-11-cr-dc-fa-code-001.md", true); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "docs/reviews/2026-06-11-cr-dc-fa-code-001.md")); !os.IsNotExist(err) {
		t.Fatalf("Review haette am alten Pfad nicht mehr existieren duerfen")
	}
	stubB, err := os.ReadFile(filepath.Join(root, "docs/reviews/archiv/2026-06-11-cr-dc-fa-code-001.md"))
	if err != nil {
		t.Fatalf("Stub fehlt im Archiv-Verzeichnis: %v", err)
	}
	stub := string(stubB)
	if !strings.Contains(stub, "Change Request `DC-FA-CODE-001` — 2026-06-11") {
		t.Fatalf("Stub traegt nicht die volle Ueberschrift (Praefix 'Review-Report:' haette nicht verschluckt werden duerfen): %q", stub)
	}
	if strings.Contains(stub, "Inhalt des Reviews") {
		t.Fatalf("Stub traegt noch den Volltext: %q", stub)
	}
	if !strings.Contains(stub, "unzip -p docs/reviews/archiv/2026-06-11-cr-dc-fa-code-001-archiv.zip") {
		t.Fatalf("Stub nennt nicht den Archiv-Pfad: %q", stub)
	}
	if !strings.Contains(stub, "DC-FA-CODE-001") || !strings.Contains(stub, "ADR-0001") {
		t.Fatalf("Stub uebernimmt die Hervorgegangen-Kennungen nicht: %q", stub)
	}

	zipPath := filepath.Join(root, "docs/reviews/archiv/2026-06-11-cr-dc-fa-code-001-archiv.zip")
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("archiv.zip fehlt oder unlesbar: %v", err)
	}
	defer zr.Close()
	if len(zr.File) != 1 || zr.File[0].Name != "docs/reviews/2026-06-11-cr-dc-fa-code-001.md" {
		t.Fatalf("archiv.zip hat unerwarteten Inhalt: %+v", zr.File)
	}

	if _, err := os.Stat(filepath.Join(root, "docs/reviews/2026-08-25-slice-137-toolchain-freshness-review.md")); err != nil {
		t.Errorf("Fremd-Review mit slice-Bezug haette unangetastet bleiben muessen: %v", err)
	}
}

// TestRunReview_RejectsSliceReview belegt: ein Review mit slice-<NNN> im
// Dateinamen gehoert in den -slice-Modus (dessen Sammel-Logik ihn findet,
// sobald der Slice archiviert wird) und wird hier abgelehnt.
func TestRunReview_RejectsSliceReview(t *testing.T) {
	root := buildReviewFixture(t, "# Review-Report: X")

	if err := runReview(root, "2026-08-25-slice-137-toolchain-freshness-review.md", true); err == nil {
		t.Fatal("erwartet Fehler fuer einen Review mit slice-<NNN> im Dateinamen")
	}
}

// TestExtractFullHeading_RealeUeberschriftenformen belegt: ExtractFullHeading
// darf das fuehrende Wort einer uneinheitlichen Review-Ueberschrift nicht
// verschlucken -- anders als ExtractTitle, das fuer das Slice-/Welle-
// Praefixschema gebaut ist. Gegen drei der elf realen Formen geprueft
// (Fixture-Kopie).
func TestExtractFullHeading_RealeUeberschriftenformen(t *testing.T) {
	cases := []struct{ heading, want string }{
		{"# Review-Report: Change Request `DC-FA-CODE-001` — 2026-06-11",
			"Review-Report: Change Request `DC-FA-CODE-001` — 2026-06-11"},
		{"# Review — Lastenheft-CR 0.15.0/0.16.0 (`--doctor` + `--repair`)",
			"Review — Lastenheft-CR 0.15.0/0.16.0 (`--doctor` + `--repair`)"},
		{"# Review Release-Prep v0.70.0 — `11d8a60`",
			"Review Release-Prep v0.70.0 — `11d8a60`"},
	}
	for _, c := range cases {
		got := ExtractFullHeading(c.heading + "\n\nRest.\n")
		if got != c.want {
			t.Errorf("ExtractFullHeading(%q) = %q, want %q", c.heading, got, c.want)
		}
	}
}

func TestFindReview_Keine(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")

	if _, err := FindReview(root, "nicht-vorhanden.md"); err == nil {
		t.Fatal("erwartet Fehler fuer einen nicht existierenden Review")
	}
}
