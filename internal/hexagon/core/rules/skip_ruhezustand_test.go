package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Ein Repo, dessen Slices alle archiviert sind, ist im Ruhezustand: die Menge,
// die erst skip-pattern leert, meldet in reviews, planning und structure
// nichts, wenn skip-allows-empty es erklärt (DC-FA-RVW-001, DC-FA-PLAN-001,
// DC-FA-STRUCT-001). Ohne den Schlüssel bleibt sie fail-closed, ebenso die
// Menge, aus der skip-pattern nichts genommen hat.

// rvArchiviert ist der Stub eines archivierten Slice, wie ihn das
// Archiv-Werkzeug schreibt: Marker, keine DoD.
const rvArchiviert = "# slice-001 — x\n\n" + stubMarker

// Mit Wellen: nur ein Stub unter done/<welle-id>/. Ohne recursive ist die
// Menge leer; recursive und skip-pattern allein lassen sie leer und rot; erst
// skip-allows-empty erklärt sie zum Ruhezustand.
func TestReviewsRuhezustand_StubUnterWelle(t *testing.T) {
	files := map[string]string{
		rvDone + "/welle-1/slice-001.md": rvArchiviert,
		rvDone + "/welle-1-results.md":   "# welle-1\n",
		"docs/reviews/.gitkeep":          "",
	}
	cfg := rvCfg()
	cfg.Match = "name"
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "leere Pruefmenge") {
		t.Fatalf("ohne recursive ist die Menge leer und ein Befund, got %+v", f)
	}
	cfg.Recursive = true
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "leere Pruefmenge") {
		t.Fatalf("ohne skip-allows-empty bleibt die geleerte Menge ein Befund, got %+v", f)
	}
	cfg.SkipAllowsEmpty = true
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("erklärter Ruhezustand ⇒ kein Befund, got %+v", f)
	}
}

// Ohne Wellen: ein flacher Stub mit require-promises — skip-pattern nimmt ihn
// aus, skip-allows-empty erklärt die Leere.
func TestReviewsRuhezustand_FlacherStubMitRequirePromises(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001.md":         rvArchiviert,
		rvDone + "/slice-001-archiv.zip": "zip",
		"docs/reviews/.gitkeep":          "",
	}
	cfg := rvCfg()
	cfg.Match = "name"
	cfg.RequirePromises = true
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "require-promises") {
		t.Fatalf("ohne skip-pattern zählt der Stub als Kandidat ohne Zusage, got %+v", f)
	}
	cfg.SkipPattern = stubPattern
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "leere Pruefmenge") {
		t.Fatalf("skip-pattern ohne skip-allows-empty ⇒ leere Menge als Befund, got %+v", f)
	}
	cfg.SkipAllowsEmpty = true
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("erklärter Ruhezustand ⇒ kein Befund, auch mit require-promises, got %+v", f)
	}
}

// require-promises urteilt nur über die übrig gelassenen Kandidaten: ein
// Volltext ohne Zusage neben einem Stub bleibt ein Befund, auch mit dem
// Schlüssel.
func TestReviewsRuhezustand_RequirePromisesZaehltNurUebrige(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001.md": rvArchiviert + "\n- [x] Review durchgeführt\n",
		rvDone + "/slice-002.md": "## 2. Definition of Done\n\n- [x] fertig\n",
	}
	cfg := rvCfg()
	cfg.RequirePromises = true
	cfg.SkipPattern = stubPattern
	cfg.SkipAllowsEmpty = true
	if f := rvRunCfg(files, cfg); len(f) != 1 || !strings.Contains(f[0].Message, "unter 1 Kandidat(en)") {
		t.Fatalf("der Volltext ohne Zusage meldet, die Zusage im Stub zählt nicht, got %+v", f)
	}
}

// Mit dem Schlüssel: skip-pattern gesetzt, aber nichts übersprungen ⇒ die
// leere Menge bleibt fail-closed, ebenso ein unlesbares reviews-dir neben
// lauter Stubs.
func TestReviewsRuhezustand_GegenprobenBleibenRot(t *testing.T) {
	cfg := rvCfg()
	cfg.SkipPattern = stubPattern
	cfg.SkipAllowsEmpty = true
	if f := rvRunCfg(map[string]string{rvDone + "/README.md": "# done\n"}, cfg); len(f) != 1 {
		t.Fatalf("keine Slice-Datei ⇒ Befund, auch mit skip-allows-empty, got %+v", f)
	}
	stubs := map[string]string{rvDone + "/slice-001.md": rvArchiviert}
	f := CheckReviews(reviewsListErrFS{MemFS: coretest.NewMemFS(stubs), errDir: "docs/reviews"}, cfg)
	if len(f) != 1 || !strings.Contains(f[0].Message, "reviews-dir lesbar: false") {
		t.Fatalf("unlesbares reviews-dir bleibt ein Befund, got %+v", f)
	}
}

// planning.closure: ein Verzeichnis aus lauter Stubs meldet nur mit dem
// Schlüssel nichts; eines ohne passende Datei weiter closure-note-missing.
func TestClosureRuhezustand(t *testing.T) {
	cfg := closureCfg()
	cfg.Closure.SkipPattern = stubPattern
	cfg.Closure.Recursive = true
	stubs := map[string]string{
		closureDir + "/slice-001-a.md":         rvArchiviert,
		closureDir + "/welle-1/slice-002-b.md": rvArchiviert,
	}
	if f := CheckPlanningClosure(coretest.NewMemFS(stubs), cfg); len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("ohne skip-allows-empty ⇒ closure-note-missing, got %+v", f)
	}
	cfg.Closure.SkipAllowsEmpty = true
	if f := CheckPlanningClosure(coretest.NewMemFS(stubs), cfg); f != nil {
		t.Fatalf("erklärter Ruhezustand ⇒ kein Befund, got %+v", f)
	}
	leer := map[string]string{closureDir + "/README.md": "# done\n"}
	f := CheckPlanningClosure(coretest.NewMemFS(leer), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("keine passende Datei ⇒ closure-note-missing, got %+v", f)
	}
}

// structure: eine Regel über lauter Stubs meldet nur mit dem Schlüssel
// nichts; leert exempt-paths die Menge, bleibt section-missing.
func TestStructureRuhezustand(t *testing.T) {
	stubs := map[string]string{"done/wellenlos/slice-001-a.md": rvArchiviert}
	r := skipRule()
	r.SkipPattern = stubPattern
	if f := CheckStructure(coretest.NewMemFS(stubs), []model.StructureRule{r}); len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("ohne skip-allows-empty ⇒ section-missing, got %+v", f)
	}
	r.SkipAllowsEmpty = true
	if f := CheckStructure(coretest.NewMemFS(stubs), []model.StructureRule{r}); f != nil {
		t.Fatalf("erklärter Ruhezustand ⇒ kein Befund, got %+v", f)
	}
	voll := map[string]string{"done/wellenlos/slice-001-a.md": "# Slice\n\n## 2. Definition of Done\n\n- [x] fertig\n"}
	r.ExemptPaths = []string{"done/wellenlos/slice-001-a.md"}
	f := CheckStructure(coretest.NewMemFS(voll), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("exempt-paths leert die Menge ⇒ section-missing, auch mit skip-allows-empty, got %+v", f)
	}
}

// Gemischte Leere: exempt-paths nimmt eine Datei, skip-pattern die übrige ⇒
// mit dem Schlüssel still; eine Datei, die beide Ausnahmen trifft, zählt als
// exempt, nicht als übersprungen.
func TestStructureRuhezustand_GemischteAusnahmen(t *testing.T) {
	r := skipRule()
	r.SkipPattern = stubPattern
	r.SkipAllowsEmpty = true
	r.ExemptPaths = []string{"done/wellenlos/slice-001-a.md"}
	gemischt := map[string]string{
		"done/wellenlos/slice-001-a.md": "# Slice\n\n## 2. Definition of Done\n\n- [x] fertig\n",
		"done/wellenlos/slice-002-b.md": rvArchiviert,
	}
	if f := CheckStructure(coretest.NewMemFS(gemischt), []model.StructureRule{r}); f != nil {
		t.Fatalf("exempt + übersprungen, mit Schlüssel ⇒ kein Befund, got %+v", f)
	}
	beides := map[string]string{"done/wellenlos/slice-001-a.md": rvArchiviert}
	if f := CheckStructure(coretest.NewMemFS(beides), []model.StructureRule{r}); len(f) != 1 {
		t.Fatalf("Datei trifft exempt und skip ⇒ zählt als exempt, die Leere bleibt ein Befund, got %+v", f)
	}
}
