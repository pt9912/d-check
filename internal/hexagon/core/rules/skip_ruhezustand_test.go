package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Ein Repo, dessen Slices alle archiviert sind, ist im Ruhezustand: die Menge,
// die erst skip-pattern leert, meldet in reviews, planning und structure
// nichts (DC-FA-RVW-001, DC-FA-PLAN-001, DC-FA-STRUCT-001). Fail-closed bleibt
// die Menge, aus der skip-pattern nichts genommen hat.

// rvArchiviert ist der Stub eines archivierten Slice, wie ihn das
// Archiv-Werkzeug schreibt: Marker, keine DoD.
const rvArchiviert = "# slice-001 — x\n\n" + stubMarker

// Mit Wellen: nur ein Stub unter done/<welle-id>/, auch ohne require-promises;
// erst recursive mit skip-pattern liest ihn und nimmt ihn aus.
func TestReviewsRuhezustand_StubUnterWelle(t *testing.T) {
	files := map[string]string{
		rvDone + "/welle-1/slice-001.md": rvArchiviert,
		rvDone + "/welle-1-results.md":   "# welle-1\n",
		"docs/reviews/.gitkeep":          "",
	}
	cfg := rvCfg()
	cfg.Match = "name"
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "leere Pruefmenge") {
		t.Fatalf("ohne recursive und skip-pattern bleibt die leere Menge ein Befund, got %+v", f)
	}
	cfg.Recursive = true
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("alle Slices archiviert ⇒ kein Befund, got %+v", f)
	}
}

// Ohne Wellen: ein flacher Stub mit require-promises — erst skip-pattern nimmt
// ihn aus der Zählung der Zusagen.
func TestReviewsRuhezustand_FlacherStubMitRequirePromises(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001.md":           rvArchiviert,
		rvDone + "/slice-001-archiv.zip":   "zip",
		"docs/reviews/.gitkeep":            "",
	}
	cfg := rvCfg()
	cfg.Match = "name"
	cfg.RequirePromises = true
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "require-promises") {
		t.Fatalf("ohne skip-pattern zählt der Stub als Kandidat ohne Zusage, got %+v", f)
	}
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("alle Slices archiviert ⇒ kein Befund, auch mit require-promises, got %+v", f)
	}
}

// require-promises urteilt nur über die übrig gelassenen Kandidaten: ein
// Volltext ohne Zusage neben einem Stub bleibt ein Befund.
func TestReviewsRuhezustand_RequirePromisesZaehltNurUebrige(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001.md": rvArchiviert + "\n- [x] Review durchgeführt\n",
		rvDone + "/slice-002.md": "## 2. Definition of Done\n\n- [x] fertig\n",
	}
	cfg := rvCfg()
	cfg.RequirePromises = true
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "unter 1 Kandidat(en)") {
		t.Fatalf("der Volltext ohne Zusage meldet, die Zusage im Stub zählt nicht, got %+v", f)
	}
}

// skip-pattern gesetzt, aber nichts übersprungen: die leere Menge bleibt
// fail-closed, ebenso ein unlesbares reviews-dir neben lauter Stubs.
func TestReviewsRuhezustand_GegenprobenBleibenRot(t *testing.T) {
	cfg := rvCfg()
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(map[string]string{rvDone + "/README.md": "# done\n"}, cfg); len(f) != 1 {
		t.Fatalf("keine Slice-Datei ⇒ Befund, auch mit skip-pattern, got %+v", f)
	}
	stubs := map[string]string{rvDone + "/slice-001.md": rvArchiviert}
	f := CheckReviews(reviewsListErrFS{MemFS: coretest.NewMemFS(stubs), errDir: "docs/reviews"}, cfg)
	if len(f) != 1 || !strings.Contains(f[0].Message, "reviews-dir lesbar: false") {
		t.Fatalf("unlesbares reviews-dir bleibt ein Befund, got %+v", f)
	}
}

// planning.closure: ein Verzeichnis aus lauter Stubs meldet nichts; eines ohne
// passende Datei weiter closure-note-missing.
func TestClosureRuhezustand(t *testing.T) {
	cfg := closureCfg()
	cfg.Closure.SkipPattern = stubPattern
	cfg.Closure.Recursive = true
	stubs := map[string]string{
		closureDir + "/slice-001-a.md":         rvArchiviert,
		closureDir + "/welle-1/slice-002-b.md": rvArchiviert,
	}
	if f := CheckPlanningClosure(coretest.NewMemFS(stubs), cfg); f != nil {
		t.Fatalf("alle Slices archiviert ⇒ kein Befund, got %+v", f)
	}
	leer := map[string]string{closureDir + "/README.md": "# done\n"}
	f := CheckPlanningClosure(coretest.NewMemFS(leer), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("keine passende Datei ⇒ closure-note-missing, got %+v", f)
	}
}

// structure: eine Regel über lauter Stubs meldet nichts; leert exempt-paths
// die Menge, bleibt section-missing.
func TestStructureRuhezustand(t *testing.T) {
	stubs := map[string]string{"done/wellenlos/slice-001-a.md": rvArchiviert}
	r := skipRule()
	r.SkipPattern = stubPattern
	if f := CheckStructure(coretest.NewMemFS(stubs), []model.StructureRule{r}); f != nil {
		t.Fatalf("alle Dateien archiviert ⇒ kein Befund, got %+v", f)
	}
	voll := map[string]string{"done/wellenlos/slice-001-a.md": "# Slice\n\n## 2. Definition of Done\n\n- [x] fertig\n"}
	r.ExemptPaths = []string{"done/wellenlos/slice-001-a.md"}
	f := CheckStructure(coretest.NewMemFS(voll), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("exempt-paths leert die Menge ⇒ section-missing, got %+v", f)
	}
}
