package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const (
	rvDone    = "docs/plan/planning/done"
	rvBestand = "## 2. Definition of Done\n\n" +
		"- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor\n" +
		"      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8\n"
	rvMuster = `Review durchgeführt`
)

func rvRunCfg(files map[string]string, cfg model.ReviewsConfig) []model.Finding {
	return CheckReviews(coretest.NewMemFS(files), cfg)
}

// Die Abnahme des CR: die DoD-Zeile aus dem Bestand des Absenders ist mit dem
// Default-Muster keine Zusage — mit promise-pattern ist sie eine, rot ohne
// Report, grün mit dem Report samt Datums-Präfix und Suffix.
func TestReviewsPromisePattern_AbnahmeDesCR(t *testing.T) {
	slice := rvDone + "/slice-039-abholzustand-und-lesepfad-publication.md"
	files := map[string]string{slice: rvBestand}
	if f := rvRunCfg(files, rvCfg()); f != nil {
		t.Fatalf("Default-Muster: die Zeile ist keine Zusage, kein Befund erwartet, got %+v", f)
	}
	cfg := rvCfg()
	cfg.PromisePattern = rvMuster
	f := rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != slice || f[0].Reason != ReasonReviewMissing || f[0].Line != 3 {
		t.Fatalf("rot: review-missing auf Zeile 3 von %s erwartet, got %+v", slice, f)
	}
	files["docs/reviews/2026-10-08-slice-039-abholzustand-und-lesepfad-publication-review.md"] = "# Review\n"
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("grün: der Report deckt die Zusage, got %+v", f)
	}
}

// Negativfälle des Musters: ein anderes Review-Wort im DoD-Punkt und dieselbe
// Wortfolge in einem Absatz außerhalb eines DoD-Punkts sind keine Zusage.
func TestReviewsPromisePattern_Negativfaelle(t *testing.T) {
	cfg := rvCfg()
	cfg.PromisePattern = rvMuster
	files := map[string]string{
		rvDone + "/slice-001-a.md": "## 2. Definition of Done\n\n- [x] Adaptions-Review notiert\n",
		rvDone + "/slice-002-b.md": "## 7. Notiz\n\nDas Review durchgeführt zu haben, half.\n",
		rvDone + "/slice-003-c.md": "## 2. Definition of Done\n\n- [ ] Review-Report liegt vor\n",
	}
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("keine der drei Dateien trägt eine Zusage nach dem Muster, got %+v", f)
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

// Die benannte Grenze von match: name: ein Basisname, der Präfix eines anderen
// ist, wird auch von dessen Report gedeckt.
func TestReviewsMatch_NamePraefixGrenze(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	files := map[string]string{
		rvDone + "/slice-a-foo.md":                  "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-10-09-slice-a-foo-bar.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("die benannte Grenze: der Report von slice-a-foo-bar deckt slice-a-foo, got %+v", f)
	}
}

// require-promises: Kandidaten ohne eine einzige Zusage melden — ohne den
// Schlüssel ist derselbe Bestand grün.
func TestReviewsRequirePromises_LeerlaufRot(t *testing.T) {
	files := map[string]string{rvDone + "/slice-039-x.md": rvBestand}
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
		full:                             "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
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

// Ein unlesbares Unterverzeichnis unter recursive ist kein stilles Grün.
func TestReviewsRecursive_UnlesbaresUnterverzeichnis(t *testing.T) {
	files := map[string]string{
		rvDone + "/slice-001-a.md":          "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
		"docs/reviews/2026-01-01-slice-001-a.md": "# Review\n",
		rvDone + "/wellenlos/slice-240-x.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review\n",
	}
	cfg := rvCfg()
	cfg.Recursive = true
	fsys := reviewsListErrFS{MemFS: coretest.NewMemFS(files), errDir: rvDone + "/wellenlos"}
	f := CheckReviews(fsys, cfg)
	if len(f) != 1 || !strings.Contains(f[0].Message, rvDone+"/wellenlos") {
		t.Fatalf("unlesbares Unterverzeichnis ⇒ Befund mit seinem Pfad, got %+v", f)
	}
}
