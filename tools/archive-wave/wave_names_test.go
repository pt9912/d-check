package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Die Tests dieser Datei halten die Koexistenz der beiden Schreibweisen
// (MR-082): `wave-<NN>` und `archive.zip` fuer Neues, `welle-<NN>` und
// `archiv.zip` fuer den Bestand. Ein Kommentar "Rot, wenn" nennt die Aenderung,
// die den Test faerbt; ein Test ohne ihn ist als Gegenprobe oder Regressions-
// waechter benannt und bleibt dabei gruen.

// Rot, wenn ArchiveName fuer wave- `archiv.zip` oder fuer welle- `archive.zip`
// liefert.
func TestArchiveName(t *testing.T) {
	cases := map[string]string{
		"wave-98":  "archive.zip",
		"wave-100": "archive.zip",
		"welle-97": "archiv.zip",
		"welle-1":  "archiv.zip",
	}
	for id, want := range cases {
		if got := ArchiveName(id); got != want {
			t.Errorf("ArchiveName(%q) = %q, erwartet %q", id, got, want)
		}
	}
}

// Rot, wenn welleIDInFieldRE `wave-` nicht liest (kein Treffer fuer wave-98)
// oder wenn der Vergleich in CollectSlices nur die Ziffern vergleicht (dann
// sammelt welle-98 die wave-98-Slices mit ein).
func TestCollectSlices_WaveUndWelleGleicheZiffernGetrennt(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-701-a.md"),
		"# Slice slice-701: A\n\n**Welle:** [wave-98](../wave-98-x.md).\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-702-b.md"),
		"# Slice slice-702: B\n\n**Welle:** wave-98.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-703-c.md"),
		"# Slice slice-703: C\n\n**Welle:** welle-98.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-98/slice-704-d.md"),
		"# Slice slice-704: D\n\n**Welle:** wave-98.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-705-e.md"),
		"# Slice slice-705: E\n\n**Welle:** — wellenlos.\n")

	wave, err := CollectSlices(root, "wave-98")
	if err != nil {
		t.Fatal(err)
	}
	if got := baseNames(wave); !reflect.DeepEqual(got, []string{"slice-701-a.md", "slice-702-b.md", "slice-704-d.md"}) {
		t.Errorf("wave-98 sammelt %v", got)
	}
	welle, err := CollectSlices(root, "welle-98")
	if err != nil {
		t.Fatal(err)
	}
	if got := baseNames(welle); !reflect.DeepEqual(got, []string{"slice-703-c.md"}) {
		t.Errorf("welle-98 sammelt %v", got)
	}
}

func baseNames(paths []string) []string {
	var out []string
	for _, p := range paths {
		out = append(out, filepath.Base(p))
	}
	return out
}

// Rot, wenn FindWellePlan fuer `wave-` nichts findet oder die Ergebnisnotiz
// `wave-99-results.md` als Plan zaehlt.
func TestFindWellePlan_Wave(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-99-testwelle.md"), "# Welle wave-99: Testwelle\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-99-results.md"), "# Ergebnis wave-99\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-99-bestand.md"), "# Welle welle-99: Bestand\n")

	got, err := FindWellePlan(root, "wave-99")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "wave-99-testwelle.md" {
		t.Fatalf("erwartet wave-99-testwelle.md, got %q", got)
	}
	got, err = FindWellePlan(root, "welle-99")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "welle-99-bestand.md" {
		t.Fatalf("erwartet welle-99-bestand.md, got %q", got)
	}
}

// Rot, wenn h1RE die Ueberschrift "# Welle wave-98: Titel" nicht als
// Kennung plus Titel liest.
func TestExtractTitle_Wave(t *testing.T) {
	if got := ExtractTitle("# Welle wave-98: Baseline heben\n\nText.\n"); got != "Baseline heben" {
		t.Fatalf("got %q", got)
	}
}

func buildWaveFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-99-testwelle.md"),
		"# Welle wave-99: Testwelle\n\nInhalt des Plans.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-99-results.md"),
		"# Ergebnis wave-99\n\nBleibt unangetastet.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-801-eins.md"),
		"# Slice slice-801: Eins\n\n**Welle:** wave-99.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-802-zwei.md"),
		"# Slice slice-802: Zwei\n\n**Welle:** [wave-99](../wave-99-testwelle.md).\n")
	writeFile(t, filepath.Join(root, "docs/reviews/2026-11-26-slice-801-r1.md"),
		"# Review slice-801\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/observations.md"),
		"Beleg: [slice-801](done/slice-801-eins.md).\n")
	return root
}

func runWaveApply(t *testing.T, root, id string) []Move {
	t.Helper()
	plan, err := FindWellePlan(root, id)
	if err != nil {
		t.Fatal(err)
	}
	slices, err := CollectSlices(root, id)
	if err != nil {
		t.Fatal(err)
	}
	reviews, err := CollectReviews(root, slices)
	if err != nil {
		t.Fatal(err)
	}
	moves, err := Apply(root, Plan{WelleID: id, WellePlan: plan, Slices: slices, Reviews: reviews})
	if err != nil {
		t.Fatal(err)
	}
	return moves
}

// Rot, wenn Apply fuer wave-99 `archiv.zip` schreibt, wenn ein Stub den alten
// Namen im unzip-Zeiger traegt oder wenn das Archiv nicht alle vier Eintraege
// haelt.
func TestFixture_EndToEnd_Wave(t *testing.T) {
	root := buildWaveFixture(t)
	moves := runWaveApply(t, root, "wave-99")
	if len(moves) != 3 {
		t.Fatalf("erwartet 3 Moves (Plan, zwei Slices), got %d: %v", len(moves), moves)
	}
	dir := filepath.Join(root, "docs/plan/planning/done/wave-99")

	if _, err := os.Stat(filepath.Join(dir, "archiv.zip")); !os.IsNotExist(err) {
		t.Errorf("wave-99 haette kein archiv.zip schreiben duerfen (err=%v)", err)
	}
	zr, err := zip.OpenReader(filepath.Join(dir, "archive.zip"))
	if err != nil {
		t.Fatalf("archive.zip fehlt: %v", err)
	}
	defer zr.Close()
	if len(zr.File) != 4 {
		t.Fatalf("erwartet 4 Zip-Eintraege, got %d", len(zr.File))
	}

	for _, name := range []string{"wave-99-testwelle.md", "slice-801-eins.md", "slice-802-zwei.md"} {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("Stub %s fehlt: %v", name, err)
		}
		if !strings.Contains(string(b), "unzip -p done/wave-99/archive.zip") {
			t.Errorf("Stub %s zeigt nicht auf archive.zip: %q", name, string(b))
		}
		if strings.Contains(string(b), "archiv.zip") {
			t.Errorf("Stub %s traegt den alten Archiv-Namen: %q", name, string(b))
		}
	}
	plan, err := os.ReadFile(filepath.Join(dir, "wave-99-testwelle.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(plan), "**Ergebnisnotiz:** wave-99-results.md") {
		t.Errorf("Welle-Stub nennt die Ergebnisnotiz nicht: %q", string(plan))
	}

	b, err := os.ReadFile(filepath.Join(root, "docs/plan/planning/done/wave-99-results.md"))
	if err != nil || string(b) != "# Ergebnis wave-99\n\nBleibt unangetastet.\n" {
		t.Errorf("Ergebnisnotiz veraendert oder geloescht: %v %q", err, string(b))
	}

	hits, err := RewriteRepo(root, moves)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("erwartet 1 betroffene Datei, got %d: %v", len(hits), hits)
	}
	b, err = os.ReadFile(filepath.Join(root, "docs/plan/planning/observations.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "Beleg: [slice-801](done/wave-99/slice-801-eins.md).\n"; string(b) != want {
		t.Fatalf("got %q, want %q", string(b), want)
	}
}

// Rot, wenn Apply fuer den Bestand `archive.zip` statt `archiv.zip`
// schreibt -- die Gegenseite von TestFixture_EndToEnd_Wave.
func TestFixture_Bestand_SchreibtWeiterArchiv(t *testing.T) {
	root := buildFixture(t)
	runWaveApply(t, root, "welle-99")
	dir := filepath.Join(root, "docs/plan/planning/done/welle-99")
	if _, err := os.Stat(filepath.Join(dir, "archiv.zip")); err != nil {
		t.Fatalf("archiv.zip fehlt: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "archive.zip")); !os.IsNotExist(err) {
		t.Errorf("welle-99 haette kein archive.zip schreiben duerfen (err=%v)", err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "slice-501-eins.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "unzip -p done/welle-99/archiv.zip") {
		t.Errorf("Bestands-Stub zeigt nicht auf archiv.zip: %q", string(b))
	}
}


// Der Fall, der nur unter Koexistenz entsteht: eine Bestands-Welle und eine
// neue Welle im selben Repo. Rot, wenn das Archivieren der einen die andere
// beruehrt (Slices, Plan, Ergebnisnotiz, Archiv-Verzeichnis) oder wenn eine
// der beiden die Schreibweise der anderen uebernimmt.
func TestFixture_Gemischt_WelleUndWaveImSelbenRepo(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-97-bestand.md"), "# Welle welle-97: Bestand\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/welle-97-results.md"), "# Ergebnis welle-97\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-901-alt.md"),
		"# Slice slice-901: Alt\n\n**Welle:** welle-97.\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-98-neu.md"), "# Welle wave-98: Neu\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-98-results.md"), "# Ergebnis wave-98\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-902-neu.md"),
		"# Slice slice-902: Neu\n\n**Welle:** wave-98.\n")
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")

	runWaveApply(t, root, "wave-98")

	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/wave-98/slice-902-neu.md")); err != nil {
		t.Errorf("der Slice der neuen Welle wurde nicht archiviert: %v", err)
	}

	for _, p := range []string{
		"docs/plan/planning/done/welle-97-bestand.md",
		"docs/plan/planning/done/welle-97-results.md",
		"docs/plan/planning/done/slice-901-alt.md",
	} {
		if _, err := os.Stat(filepath.Join(root, p)); err != nil {
			t.Errorf("%s gehoert der Bestands-Welle und darf nicht angefasst sein: %v", p, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/welle-97")); !os.IsNotExist(err) {
		t.Errorf("das Verzeichnis der Bestands-Welle haette nicht entstehen duerfen (err=%v)", err)
	}

	runWaveApply(t, root, "welle-97")
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/welle-97/archiv.zip")); err != nil {
		t.Errorf("archiv.zip der Bestands-Welle fehlt: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/wave-98/archive.zip")); err != nil {
		t.Errorf("archive.zip der neuen Welle fehlt: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/wave-98/archiv.zip")); !os.IsNotExist(err) {
		t.Errorf("die neue Welle haette kein archiv.zip bekommen duerfen (err=%v)", err)
	}
}

// Rot, wenn runWelle fuer wave-99 im Dry-Run schreibt.
func TestRun_DryRun_NoWrites_Wave(t *testing.T) {
	root := buildWaveFixture(t)
	before := snapshot(t, root)
	if err := runWelle(root, "wave-99", false); err != nil {
		t.Fatalf("Dry-Run schlug fehl: %v", err)
	}
	if after := snapshot(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("Dry-Run hat den Baum veraendert")
	}
}

// Regressionswaechter, gruen auch ohne Praefix-Erkennung in collect.go (der
// Plan wird ueber das Praefix gefunden, der Fail-Closed-Zweig "weder
// Welle-Plan noch Slices" nicht erreicht). Rot, wenn runWelle fuer wave-
// `archiv.zip` schreibt.
func TestRun_PlanOhneEigeneSlices_Wave(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-73-x.md"), "# Welle wave-73: X\n")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/wave-73-results.md"), "# Ergebnis wave-73\n")
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")

	if err := runWelle(root, "wave-73", true); err != nil {
		t.Fatalf("erwartet Erfolg, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/plan/planning/done/wave-73/archive.zip")); err != nil {
		t.Errorf("archive.zip fehlt: %v", err)
	}
}

// Rot, wenn Slice-Stub oder Welle-Stub fuer wave- den alten Namen zeigen.
func TestStubs_Wave(t *testing.T) {
	slice := SliceStub("slice-801", "Eins", "wave-99", "wave-99", "DC-FA-PLAN-001")
	if !strings.Contains(slice, "unzip -p done/wave-99/archive.zip") {
		t.Errorf("Slice-Stub: %q", slice)
	}
	wave := WelleStub("wave-99", "Testwelle", "wave-99-results.md", 2, 1)
	if !strings.Contains(wave, "unzip -p done/wave-99/archive.zip") {
		t.Errorf("Welle-Stub: %q", wave)
	}
}

// Rot, wenn der Slice-Modus `-archiv.zip` neben oder statt `-archive.zip`
// schreibt (Verdikt Frage 8: Slice und Review schreiben die neue Form).
func TestRunSlice_SchreibtNurDieNeueForm(t *testing.T) {
	root := buildSliceFixture(t)
	if err := runSlice(root, "slice-601", true); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Join(root, "docs/plan/planning/done/wellenlos"))
	if err != nil {
		t.Fatal(err)
	}
	var zips []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".zip") {
			zips = append(zips, e.Name())
		}
	}
	if !reflect.DeepEqual(zips, []string{"slice-601-archive.zip"}) {
		t.Fatalf("erwartet genau slice-601-archive.zip, got %v", zips)
	}
}

// Rot, wenn der Review-Modus `-archiv.zip` neben oder statt `-archive.zip`
// schreibt.
func TestRunReview_SchreibtNurDieNeueForm(t *testing.T) {
	root := buildReviewFixture(t, "# Review-Report: Change Request `DC-FA-CODE-001` — 2026-06-11")
	if err := runReview(root, "2026-06-11-cr-dc-fa-code-001.md", true); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Join(root, "docs/reviews/archiv"))
	if err != nil {
		t.Fatal(err)
	}
	var zips []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".zip") {
			zips = append(zips, e.Name())
		}
	}
	if !reflect.DeepEqual(zips, []string{"2026-06-11-cr-dc-fa-code-001-archive.zip"}) {
		t.Fatalf("got %v", zips)
	}
}

// Rot, wenn die Zugehoerigkeits-Pruefung des Slice-Modus `wave-` nicht
// erkennt: ein Slice einer neuen Welle wuerde dann als wellenlos archiviert.
func TestRunSlice_RejectsWaveSlice(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-603-mit-wave.md"),
		"# Slice slice-603: Mit Wave\n\n**Welle:** wave-42.\n")
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")
	err := runSlice(root, "slice-603", true)
	if err == nil || !strings.Contains(err.Error(), "gehoert zu einer Welle") {
		t.Fatalf("erwartet die Ablehnung wegen der Wellen-Zugehoerigkeit, got %v", err)
	}
}

// Die Gegenprobe: eine fremde neue Welle in einer spaeteren Zeile des Feldes
// ist Kontext, keine Zugehoerigkeit (wie bei welle-86).
func TestRunSlice_AcceptsContextualWaveMention(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-605-kontext.md"),
		"# Slice slice-605: Kontext\n\n**Welle:** — wellenlos, solange keine Closure-Bedingung ueber die eigene\n"+
			"DoD hinausgeht. Der Anlass liegt in [wave-86](wave-86/x.md).\n")
	writeFile(t, filepath.Join(root, "docs/reviews/.gitkeep"), "")
	if err := runSlice(root, "slice-605", false); err != nil {
		t.Fatalf("wave-86 ist hier nur Kontext: %v", err)
	}
}
