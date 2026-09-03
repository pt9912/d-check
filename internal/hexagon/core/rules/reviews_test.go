package rules

import (
	"fmt"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

func rvCfg() model.ReviewsConfig {
	return model.ReviewsConfig{DoneDir: "docs/plan/planning/done", ReviewsDir: "docs/reviews"}
}

func rvRun(files map[string]string) []model.Finding {
	return CheckReviews(coretest.NewMemFS(files), rvCfg())
}

// Inert: kein DoneDir konfiguriert ⇒ keine Datei wird geoeffnet.
func TestReviewsInert(t *testing.T) {
	if f := CheckReviews(coretest.NewMemFS(map[string]string{
		"docs/plan/planning/done/slice-001-x.md": "## 2. Definition of Done\n\n- [x] Review\n",
	}), model.ReviewsConfig{}); f != nil {
		t.Fatalf("erwartet nil (inert), got %+v", f)
	}
}

// Happy Path: Review-Zusage UND passender Report ⇒ kein Befund.
func TestReviewsHappyPath(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/slice-100-x.md": "## 2. Definition of Done\n\n" +
			"- [x] `make gates` grün; unabhängiger Review.\n",
		"docs/reviews/2026-01-01-slice-100-x-review.md": "# Review\n",
	}
	if f := rvRun(files); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", reasons(f))
	}
}

// Negativ: Review-Zusage ohne passenden Report ⇒ review-missing, Zeile der
// Zusage-Zeile im Dokument.
func TestReviewsMissing(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/slice-101-y.md": "# Slice slice-101\n\n## 2. Definition of Done\n\n" +
			"- [ ] unabhängiger Review.\n",
		"docs/reviews/2026-01-01-slice-999-unrelated.md": "# Review\n",
	}
	f := rvRun(files)
	if !hasReason(f, ReasonReviewMissing) {
		t.Fatalf("erwartet review-missing, got %+v", reasons(f))
	}
	var hit *model.Finding
	for i := range f {
		if f[i].File == "docs/plan/planning/done/slice-101-y.md" {
			hit = &f[i]
		}
	}
	if hit == nil {
		t.Fatalf("kein Befund auf der Slice-Datei, got %+v", f)
	}
	if hit.Line != 5 {
		t.Errorf("Line = %d, want 5 (die Review-Zusage-Zeile)", hit.Line)
	}
}

// Alle drei CommonMark-Bullet-Formen tragen die Zusage, unabhängig vom
// Haken-Zustand -- dieselbe Toleranz wie taskItemRE/ADR-0074.
func TestReviewsBulletForms(t *testing.T) {
	for _, bullet := range []string{"-", "*", "+"} {
		content := "## 2. Definition of Done\n\n" + bullet + " [ ] unabhängiger Review.\n"
		files := map[string]string{"docs/plan/planning/done/slice-102-z.md": content}
		f := rvRun(files)
		if !hasReason(f, ReasonReviewMissing) {
			t.Fatalf("bullet %q: erwartet review-missing, got %+v", bullet, reasons(f))
		}
	}
}

// Ohne Review-Zeile ist der Slice kein Kandidat -- kein Befund auf IHM, aber
// die leere Zusage-Menge ist selbst fail-closed, wenn sie repo-weit auf null
// faellt.
func TestReviewsNoPromiseNoFinding(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/slice-103-a.md": "## 2. Definition of Done\n\n- [x] `make gates` grün.\n",
		"docs/reviews/x.md":                      "irrelevant",
	}
	f := rvRun(files)
	if len(f) != 0 {
		t.Fatalf("erwartet befundfrei (keine Zusage, kein Fail-Closed-Trigger da anderswo Kandidaten fehlen), got %+v", reasons(f))
	}
}

// Fail-closed: DoneDir aktiv, aber keine Kandidaten -- sieht sonst identisch
// aus wie "alles gedeckt".
func TestReviewsEmptyScopeFailsClosed(t *testing.T) {
	f := rvRun(map[string]string{"docs/plan/planning/done/welle-01-x.md": "kein Slice"})
	if !hasReason(f, ReasonReviewMissing) {
		t.Fatalf("erwartet review-missing (leere Pruefmenge), got %+v", reasons(f))
	}
}

// exempt-paths hebt einen Kandidaten aus der Pruefmenge -- der einzige Slice
// verschwindet damit aus candidates, und der Leerlauf-Befund bleibt stehen
// (Konfig-Kommentar: "hebt den Leerlauf-Befund NICHT auf") statt eines
// per-Slice-Befunds auf der ausgenommenen Datei.
func TestReviewsExemptPaths(t *testing.T) {
	cfg := model.ReviewsConfig{
		DoneDir: "docs/plan/planning/done", ReviewsDir: "docs/reviews",
		ExemptPaths: []string{"docs/plan/planning/done/slice-104-*.md"},
	}
	f := CheckReviews(coretest.NewMemFS(map[string]string{
		"docs/plan/planning/done/slice-104-b.md": "## 2. Definition of Done\n\n- [ ] unabhängiger Review.\n",
	}), cfg)
	if len(f) != 1 || f[0].File != "docs/plan/planning/done" {
		t.Fatalf("erwartet genau den Leerlauf-Befund auf DoneDir, got %+v", f)
	}
}

// Archivierte Stubs liegen in Unterverzeichnissen und sind nicht Kandidat --
// reviewCandidates scannt DoneDir nicht rekursiv.
func TestReviewsIgnoresArchivedSubdirs(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/welle-87/slice-190-x.md": "- [ ] unabhängiger Review.\n",
	}
	f := rvRun(files)
	if len(f) != 1 || f[0].File != "docs/plan/planning/done" {
		t.Fatalf("archivierter Stub wurde faelschlich als Kandidat gelesen, oder Fail-Closed griff nicht: %+v", f)
	}
}

// H1-Regression (Review-Fund, commit 85b1fce): die Phrase steht bei der
// Mehrheit des Bestands NICHT auf der Checkbox-Zeile selbst, sondern auf einer
// Folgezeile desselben DoD-Items (gemessen an slice-138 u. a.). Eine Regel, die
// nur die Checkbox-Zeile liest, uebersieht das und meldet faelschlich
// "keine Zusage".
func TestReviewsPromiseOnContinuationLine(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/slice-106-b.md": "## 2. Definition of Done\n\n" +
			"- [x] Lastenheft auf 0.65.4; `make gates` Exit 0 (zehn\n" +
			"      Glieder), unabhängiger Review\n" +
			"      ([Report](../../../reviews/x.md)).\n",
	}
	f := rvRun(files)
	if !hasReason(f, ReasonReviewMissing) {
		t.Fatalf("erwartet review-missing (Phrase auf Folgezeile), got %+v", reasons(f))
	}
	var hit *model.Finding
	for i := range f {
		if f[i].File == "docs/plan/planning/done/slice-106-b.md" {
			hit = &f[i]
		}
	}
	if hit == nil {
		t.Fatalf("kein Befund auf der Slice-Datei, got %+v", f)
	}
	if hit.Line != 3 {
		t.Errorf("Line = %d, want 3 (die Checkbox-Zeile des Items, nicht die Phrasen-Zeile)", hit.Line)
	}
}

// Gegenprobe zu TestReviewsPromiseOnContinuationLine: dieselbe Form, aber mit
// passendem Report -- die Zusammenfuehrung mehrerer Zeilen zu einem Item darf
// den Happy Path nicht brechen.
func TestReviewsPromiseOnContinuationLineCovered(t *testing.T) {
	files := map[string]string{
		"docs/plan/planning/done/slice-106-b.md": "## 2. Definition of Done\n\n" +
			"- [x] Lastenheft auf 0.65.4; `make gates` Exit 0 (zehn\n" +
			"      Glieder), unabhängiger Review\n" +
			"      ([Report](../../../reviews/x.md)).\n",
		"docs/reviews/2026-01-01-slice-106-review.md": "# Review\n",
	}
	if f := rvRun(files); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", reasons(f))
	}
}

// reviewsListErrFS erzwingt einen Lesefehler fuer EIN benanntes Verzeichnis --
// coretest.MemFS selbst kennt keine Fehler, ein unlesbares reviews-dir
// braucht deshalb eine eigene Fake.
type reviewsListErrFS struct {
	*coretest.MemFS
	errDir string
}

func (f reviewsListErrFS) List(rel string) ([]driven.DirEntry, error) {
	if rel == f.errDir {
		return nil, fmt.Errorf("simuliert: %s nicht lesbar", rel)
	}
	return f.MemFS.List(rel)
}

// H2-Regression (Review-Fund, commit 85b1fce): ein unlesbares reviews-dir MIT
// vorhandenen Review-Zusagen darf NICHT zusaetzlich zu den bereits erzeugten
// Pro-Kandidat-Befunden eine zweite, textlich widerspruechliche
// "leere Pruefmenge"-Meldung erzeugen -- die Menge ist dann gerade NICHT leer.
func TestReviewsUnreadableReviewsDirWithPromisesNoRedundantFinding(t *testing.T) {
	base := coretest.NewMemFS(map[string]string{
		"docs/plan/planning/done/slice-105-x.md": "## 2. Definition of Done\n\n- [x] unabhängiger Review.\n",
	})
	fsys := reviewsListErrFS{MemFS: base, errDir: "docs/reviews"}
	f := CheckReviews(fsys, rvCfg())
	if len(f) != 1 {
		t.Fatalf("erwartet genau einen Befund (der Kandidat selbst, keine zusaetzliche Leerlauf-Meldung), got %+v", f)
	}
	if f[0].File != "docs/plan/planning/done/slice-105-x.md" {
		t.Fatalf("erwartet Befund auf dem Kandidaten, got %+v", f)
	}
}

// Gegenprobe: ein unlesbares reviews-dir OHNE jede Kandidaten-Zusage
// braucht weiterhin die generische Fail-Closed-Meldung -- sonst waere ein
// leerer Bestand von einem unlesbaren Verzeichnis nicht zu unterscheiden.
func TestReviewsUnreadableReviewsDirNoPromisesStillFailsClosed(t *testing.T) {
	base := coretest.NewMemFS(map[string]string{
		"docs/plan/planning/done/slice-105-x.md": "## 2. Definition of Done\n\n- [x] `make gates` grün.\n",
	})
	fsys := reviewsListErrFS{MemFS: base, errDir: "docs/reviews"}
	f := CheckReviews(fsys, rvCfg())
	if len(f) != 1 || f[0].File != "docs/plan/planning/done" {
		t.Fatalf("erwartet genau die Leerlauf-Meldung auf DoneDir, got %+v", f)
	}
}
