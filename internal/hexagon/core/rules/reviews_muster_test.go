package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const (
	rvDone = "docs/plan/planning/done"
	// rvVorlage ist der Review-Punkt der Baseline-Slice-Vorlage.
	rvVorlage = "## 2. Definition of Done\n\n" +
		"- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor\n" +
		"      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8\n"
	// rvEigen ist ein Review-Punkt, den nur ein eigenes Muster erkennt.
	rvEigen  = "## 2. Definition of Done\n\n- [x] Code-Review erledigt, Report liegt vor\n"
	rvMuster = `Code-Review erledigt`
)

func rvRunCfg(files map[string]string, cfg model.ReviewsConfig) []model.Finding {
	return CheckReviews(coretest.NewMemFS(files), cfg)
}

// Der Review-Punkt der Baseline-Vorlage ist mit dem Default eine Zusage: rot
// ohne Report, grün mit einem Report samt Datums-Präfix und Suffix.
func TestReviewsDefault_ErkenntVorlagenForm(t *testing.T) {
	slice := rvDone + "/slice-039-abholzustand-und-lesepfad-publication.md"
	files := map[string]string{slice: rvVorlage}
	f := rvRunCfg(files, rvCfg())
	if len(f) != 1 || f[0].File != slice || f[0].Reason != ReasonReviewMissing || f[0].Line != 3 {
		t.Fatalf("rot: review-missing auf Zeile 3 von %s erwartet, got %+v", slice, f)
	}
	files["docs/reviews/2026-10-08-slice-039-abholzustand-und-lesepfad-publication-review.md"] = "# Review\n"
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("grün: der Report deckt die Zusage, got %+v", f)
	}
}

// Ein eigenes Muster erkennt eine Form, die der Default nicht kennt.
func TestReviewsPromisePattern_EigeneForm(t *testing.T) {
	slice := rvDone + "/slice-040-x.md"
	files := map[string]string{slice: rvEigen}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("Default: „Code-Review erledigt“ ist keine Zusage, got %+v", f)
	}
	cfg := rvCfg()
	cfg.PromisePattern = rvMuster
	f := rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != slice || f[0].Reason != ReasonReviewMissing {
		t.Fatalf("eigenes Muster: review-missing auf %s erwartet, got %+v", slice, f)
	}
}

// Negativfälle des Defaults: ein anderes Review-Wort im DoD-Punkt und
// dieselbe Wortfolge außerhalb eines Checkbox-Punkts sind keine Zusage.
func TestReviewsDefault_Negativfaelle(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001-a.md": "## 2. Definition of Done\n\n- [x] Adaptions-Review notiert\n",
		rvDone + "/slice-002-b.md": "## 7. Notiz\n\nDas Review durchgeführt zu haben, half.\n",
		rvDone + "/slice-003-c.md": "## 2. Definition of Done\n\n- [ ] Review-Report liegt vor\n",
		rvDone + "/slice-004-d.md": "## 7. Notiz\n\nEin unabhängiger Review fand drei Befunde.\n",
	}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("keine der vier Dateien trägt eine Zusage, got %+v", f)
	}
}

// Die benannte Grenze: jeder Checkbox-Punkt der Datei zählt, auch einer in
// einem Codeblock.
func TestReviewsDefault_CheckboxImCodeblockZaehlt(t *testing.T) {
	slice := rvDone + "/slice-005-e.md"
	files := map[string]string{slice: "## 3. Beispiel\n\n```\n- [ ] Review durchgeführt\n```\n"}
	f := rvRunCfg(files, rvCfg())
	if len(f) != 1 || f[0].File != slice || f[0].Line != 4 {
		t.Fatalf("die benannte Grenze: der Punkt im Codeblock zählt, got %+v", f)
	}
}

// Eine Zusage in einer Datei ohne slice-<NNN>-Kennung fällt unter match: id
// nicht still aus, sondern meldet; unter match: name deckt der Report sie.
func TestReviewsMatch_BenannteKennung(t *testing.T) {
	slice := rvDone + "/slice-w3-login-flow.md"
	files := map[string]string{
		slice: "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-10-09-slice-w3-login-flow-r1.md": "# Review\n",
	}
	f := rvRunCfg(files, rvCfg())
	if len(f) != 1 || f[0].File != slice || !strings.Contains(f[0].Message, "match: name") {
		t.Fatalf("match: id — Befund mit Hinweis auf match: name erwartet, got %+v", f)
	}
	cfg := rvCfg()
	cfg.Match = "name"
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("match: name — der Report deckt die Zusage, got %+v", f)
	}
	delete(files, "docs/reviews/2026-10-09-slice-w3-login-flow-r1.md")
	files["docs/reviews/2026-10-09-slice-w3-anderes-r1.md"] = "# Review\n"
	f = rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != slice || f[0].Reason != ReasonReviewMissing {
		t.Fatalf("match: name — ein fremder Report deckt nicht, got %+v", f)
	}
}

// match: name verlangt nach dem Basisnamen ein Zeichen, das weder Buchstabe
// noch Ziffer ist: slice-x1 wird nicht vom Report zu slice-x12-y gedeckt. Die
// benannte Grenze: mit Bindestrich deckt der Report zu slice-a-foo-bar auch
// slice-a-foo.
func TestReviewsMatch_NameWortgrenze(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	ziffer := map[string]string{
		rvDone + "/slice-x1.md":                     "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-10-09-slice-x12-y-r1.md": "# Review\n",
	}
	if f := rvRunCfg(ziffer, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-x1.md" {
		t.Fatalf("slice-x1 darf nicht vom Report zu slice-x12-y gedeckt sein, got %+v", f)
	}
	strich := map[string]string{
		rvDone + "/slice-a-foo.md":                  "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-10-09-slice-a-foo-bar.md": "# Review\n",
	}
	if f := rvRunCfg(strich, cfg); f != nil {
		t.Fatalf("die benannte Grenze: der Report zu slice-a-foo-bar deckt slice-a-foo, got %+v", f)
	}
}

// require-promises: Kandidaten ohne eine einzige Zusage melden — ohne den
// Schlüssel ist derselbe Bestand grün.
func TestReviewsRequirePromises_LeerlaufRot(t *testing.T) {
	files := map[string]string{rvDone + "/slice-039-x.md": rvEigen}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("ohne require-promises: grün über null Zusagen, got %+v", f)
	}
	cfg := rvCfg()
	cfg.RequirePromises = true
	f := rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != rvDone || !strings.Contains(f[0].Message, "require-promises") {
		t.Fatalf("require-promises: ein Befund auf done-dir erwartet, got %+v", f)
	}
}

// recursive und skip-pattern: der Volltext unter wellenlos/ wird geprüft, der
// Stub daneben nicht — ohne recursive bleibt das Unterverzeichnis ungelesen.
func TestReviewsRecursiveUndSkipPattern(t *testing.T) {
	full := rvDone + "/wellenlos/slice-240-x.md"
	files := map[string]string{
		rvDone + "/slice-001-a.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-01-01-slice-001-a-review.md": "# Review\n",
		full: "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		rvDone + "/wellenlos/slice-083-y.md": "# slice-083\n\n" + stubMarker + "\n- [x] unabhängiger Review\n",
	}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("ohne recursive ist wellenlos/ ungelesen, got %+v", f)
	}
	cfg := rvCfg()
	cfg.Recursive = true
	if f := rvRunCfg(files, cfg); len(f) != 2 {
		t.Fatalf("recursive ohne skip-pattern: Volltext und Stub melden, got %+v", f)
	}
	cfg.SkipPattern = stubPattern
	f := rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != full {
		t.Fatalf("recursive mit skip-pattern: nur der Volltext %s meldet, got %+v", full, f)
	}
}

// Unter recursive bleiben die immer übersprungenen Verzeichnisse unbetreten.
func TestReviewsRecursive_SkipDirsBleibenUnbetreten(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001-a.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-01-01-slice-001-a.md": "# Review\n",
		rvDone + "/node_modules/slice-002-b.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
	}
	cfg := rvCfg()
	cfg.Recursive = true
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("node_modules/ darf nicht betreten werden, got %+v", f)
	}
}

// Ein unlesbares Unterverzeichnis meldet mit seinem Pfad; die übrigen
// Einträge werden weiter gelesen, und der Leerlauf-Befund entfällt neben ihm.
func TestReviewsRecursive_UnlesbaresUnterverzeichnis(t *testing.T) {
	missing := rvDone + "/slice-900-z.md"
	files := map[string]string{
		rvDone + "/archiv/slice-100-q.md":  "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		missing:                            "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		rvDone + "/zz/slice-950-w.md":      "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-01-01-slice-950-w.md": "# Review\n",
	}
	cfg := rvCfg()
	cfg.Recursive = true
	fsys := reviewsListErrFS{MemFS: coretest.NewMemFS(files), errDir: rvDone + "/archiv"}
	f := CheckReviews(fsys, cfg)
	if len(f) != 2 || !strings.Contains(f[0].Message, rvDone+"/archiv") || f[1].File != missing {
		t.Fatalf("Befund auf dem Verzeichnis und auf %s, kein Leerlauf-Befund erwartet, got %+v", missing, f)
	}
}

// Liegt der einzige Bestand im unlesbaren Unterverzeichnis, meldet nur das
// Verzeichnis — kein zweiter, widersprüchlicher Leerlauf-Befund.
func TestReviewsRecursive_UnlesbaresUnterverzeichnisOhneLeerlauf(t *testing.T) {
	files := map[string]string{
		rvDone + "/archiv/slice-100-q.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-01-01-x.md":    "# Review\n",
	}
	cfg := rvCfg()
	cfg.Recursive = true
	cfg.RequirePromises = true
	fsys := reviewsListErrFS{MemFS: coretest.NewMemFS(files), errDir: rvDone + "/archiv"}
	f := CheckReviews(fsys, cfg)
	if len(f) != 1 || !strings.Contains(f[0].Message, rvDone+"/archiv") {
		t.Fatalf("genau ein Befund auf dem unlesbaren Verzeichnis erwartet, got %+v", f)
	}
}

// Die Wortgrenze nach dem Basisnamen gilt für jeden Buchstaben, nicht nur
// ASCII: slice-a-grö wird nicht vom Report zu slice-a-größe gedeckt.
func TestReviewsMatch_NameWortgrenzeUnicode(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	files := map[string]string{
		rvDone + "/slice-a-grö.md":                    "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-10-09-slice-a-größe-r1.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-a-grö.md" {
		t.Fatalf("slice-a-grö darf nicht vom Report zu slice-a-größe gedeckt sein, got %+v", f)
	}
}

// Die Vorlagen-Form zählt nur, wo sie den Punkt direkt hinter der Task-Box
// eröffnet: zusammengesetzt, verneint oder mitten im Satz ist sie keine
// Zusage.
func TestReviewsDefault_VorlagenFormNurAmPunktanfang(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001-a.md": "## 2. Definition of Done\n\n- [x] Adaptions-Review durchgeführt\n",
		rvDone + "/slice-002-b.md": "## 6. Risiken\n\n- [ ] kein Review durchgeführt (entfallen)\n",
		rvDone + "/slice-003-c.md": "## 2. Definition of Done\n\n- [x] `make gates` grün, Review durchgeführt\n",
	}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("keine der drei Formen ist eine Zusage, got %+v", f)
	}
}
