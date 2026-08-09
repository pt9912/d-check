package rules

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

const closureDir = "docs/plan/planning/done"

// closureCfg baut eine planning-Config mit aktiver Closure-Fähigkeit.
func closureCfg() model.PlanningConfig {
	return model.PlanningConfig{
		Roadmap: planRoadmap,
		Closure: model.ClosureConfig{Dir: closureDir},
	}
}

// richNote ist eine Closure-Notiz über der Default-Schwelle (4 Satzende-Zeichen).
const richNote = "## 7. Closure-Notiz (nach `done/`)\n\n" +
	"Umgesetzt: das Modul liegt vor. Der Test war rot, weil die Grenze fehlte.\n" +
	"Folge-Slice: slice-999. Architektur-Beobachtung: der Port trägt zwei Rollen.\n"

// listErrFS ist ein Filesystem, dessen List **immer** scheitert — MemFS kann
// „Verzeichnis fehlt" nicht ausdrücken (es liefert eine leere Liste), und genau
// die Unterscheidung trägt den fail-closed Vertrag (C2).
type listErrFS struct{ driven.Filesystem }

func (listErrFS) List(string) ([]driven.DirEntry, error) { return nil, errors.New("kein Verzeichnis") }

// reasonsOf verdichtet Befunde auf ihre Grund-Codes.
func reasonsOf(fs []model.Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Reason)
	}
	return out
}

// Happy Path: Notizen über der Schwelle ⇒ kein Befund.
func TestClosureHappyPath(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote,
		closureDir + "/slice-002-b.md": "# Slice\n\n" + richNote,
	}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}

// Boundary: ohne closure.dir ist die Fähigkeit inert — auch wenn im Baum eine
// leere Closure-Notiz liegt, die sonst rot wäre (byte-identisch, DC-QA-02).
func TestClosureInertOhneDir(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n## 7. Closure-Notiz\n\n_Ausstehend._\n",
	}
	cfg := model.PlanningConfig{Roadmap: planRoadmap} // kein Closure.Dir
	if f := CheckPlanningClosure(coretest.NewMemFS(files), cfg); f != nil {
		t.Fatalf("ohne closure.dir muss die Fähigkeit inert sein, got %+v", f)
	}
}

// Negative: der zurückgelassene Platzhalter — der reale Fehlerfall.
func TestClosureThinPlatzhalter(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n## 9. Closure-Notiz (nach `done/`)\n\n_Ausstehend._\n",
	}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("erwartet genau ein closure-note-thin, got %+v", f)
	}
	if f[0].Line != 3 {
		t.Errorf("Zeile muss die Abschnitts-Überschrift sein (3), got %d", f[0].Line)
	}
	if f[0].Target != closureDir {
		t.Errorf("target muss das Closure-Verzeichnis sein, got %q", f[0].Target)
	}
	// Die Meldung nennt Ist und Soll — sonst weiß der Autor nicht, wie weit er fehlt.
	if !strings.Contains(f[0].Message, "1") || !strings.Contains(f[0].Message, "4") {
		t.Errorf("Meldung muss Ist- und Soll-Zahl nennen, got %q", f[0].Message)
	}
}

// Negative: kein passender Abschnitt ⇒ missing, und missing schließt die
// beiden Mess-Codes aus (ohne Abschnitt gibt es nichts zu messen).
func TestClosureMissingSchliesstMessungAus(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n## 7. Ergebnis\n\nKein Closure-Abschnitt hier.\n",
	}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("erwartet genau ein closure-note-missing, got %+v", f)
	}
}

// Negative: Floskel, case-insensitiv.
func TestClosureBoilerplateCaseInsensitiv(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote + "\nWIE GEPLANT UMGESETZT.\n",
	}
	cfg := closureCfg()
	cfg.Closure.Boilerplate = []string{"wie geplant umgesetzt"}
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteBoilerplate {
		t.Fatalf("erwartet genau ein closure-note-boilerplate, got %+v", f)
	}
}

// Boundary: Satzende-Zeichen INNERHALB eines Fenced-Blocks zählen nicht — die
// Bereinigung ist wirksam, nicht dekorativ. Ohne sie wäre dieser Fall grün.
func TestClosureCodeBlockZaehltNicht(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nFertig.\n\n```\nsatz. satz. satz. satz. satz.\n```\n"
	files := map[string]string{closureDir + "/slice-001-a.md": "# Slice\n\n" + note}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("Fence-Inhalt darf nicht zählen ⇒ thin erwartet, got %+v", f)
	}
}

// Boundary: eine '#'-Zeile IM Fence ist keine Überschrift — sonst würde ein
// Beispielblock den geprüften Abschnitt eröffnen.
func TestClosureHeadingImFenceZaehltNicht(t *testing.T) {
	body := "# Slice\n\n```\n## 7. Closure-Notiz\n```\n\n## 7. Ergebnis\n\nNichts.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": body}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("Fence-Überschrift darf nicht zählen ⇒ missing erwartet, got %+v", f)
	}
}

// Abschnitts-Grenze: eine TIEFERE Überschrift gehört noch zum Abschnitt, eine
// gleich-/höherrangige beendet ihn.
func TestClosureAbschnittsGrenzeNachEbene(t *testing.T) {
	// Der Abschnitt endet bei '## 8.'; die Sätze danach dürfen nicht mitzählen.
	// Die Unter-Überschrift bewusst OHNE Punkt: eine Nummer wie „7.1" trüge
	// selbst ein Satzende-Zeichen (gewollt simpel, siehe countSentenceEnds) und
	// machte den Fall mehrdeutig.
	body := "# Slice\n\n## 7. Closure-Notiz\n\nEin Satz.\n\n### Vertiefung\n\nZwei. Drei.\n\n" +
		"## 8. Anhang\n\nVier. Fünf. Sechs. Sieben.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": body}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("nur der Abschnitt bis zur nächsten H2 zählt ⇒ thin (3 < 4), got %+v", f)
	}
}

// Nur Dateien, die dem slice-glob entsprechen, sind Kandidaten.
func TestClosureNurSliceGlobKandidaten(t *testing.T) {
	files := map[string]string{
		closureDir + "/welle-60-results.md": "# Welle\n\nOhne Closure-Abschnitt.\n",
		closureDir + "/README.md":           "# Index\n",
		closureDir + "/slice-001-a.md":      "# Slice\n\n" + richNote,
	}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("Nicht-Slice-Dateien sind keine Kandidaten, got %+v", f)
	}
}

// Ein leeres Closure-Verzeichnis ist befundfrei (kein abgeschlossener Slice ist
// kein Fehler) — Gegenstück zum fail-closed Fall darunter.
func TestClosureLeeresVerzeichnisBefundfrei(t *testing.T) {
	files := map[string]string{closureDir + "/.keep": ""}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("leeres Closure-Verzeichnis muss befundfrei sein, got %+v", f)
	}
}

// fail-closed: gesetztes, aber unlesbares Verzeichnis ⇒ Befund statt stillem
// Grün. Sonst schaltete ein Tippfehler im Pfad die ganze Prüfung ab.
func TestClosureVerzeichnisUnlesbarFailClosed(t *testing.T) {
	fsys := listErrFS{coretest.NewMemFS(map[string]string{})}
	f := CheckPlanningClosure(fsys, closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("unlesbares Closure-Verzeichnis muss fail-closed melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "fail-closed") {
		t.Errorf("Meldung muss den fail-closed-Grund nennen, got %q", f[0].Message)
	}
}

// Ein Kandidat kann thin UND boilerplate tragen (verschiedene Bedingungen).
func TestClosureThinUndBoilerplateNebeneinander(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n## 7. Closure-Notiz\n\nFertig.\n",
	}
	cfg := closureCfg()
	cfg.Closure.Boilerplate = []string{"Fertig."}
	got := reasonsOf(CheckPlanningClosure(coretest.NewMemFS(files), cfg))
	if len(got) != 2 {
		t.Fatalf("erwartet zwei Befunde (thin + boilerplate), got %v", got)
	}
}

// Determinismus: identischer Baum ⇒ identische Reihenfolge (DC-QA-02).
func TestClosureDeterministischeReihenfolge(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-003-c.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
		closureDir + "/slice-001-a.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
		closureDir + "/slice-002-b.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
	}
	first := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	for i := 0; i < 5; i++ {
		again := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
		for j := range first {
			if first[j].File != again[j].File {
				t.Fatalf("Reihenfolge nicht stabil: %s vs %s", first[j].File, again[j].File)
			}
		}
	}
	if len(first) != 3 || first[0].File != closureDir+"/slice-001-a.md" {
		t.Fatalf("erwartet drei Befunde, nach Namen sortiert, got %+v", first)
	}
}

// Ein eigenes heading-pattern schlägt den Default.
func TestClosureEigenesHeadingPattern(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n## Abschluss\n\nEin Satz.\n",
	}
	cfg := closureCfg()
	cfg.Closure.HeadingPattern = `^## Abschluss`
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("eigenes Muster muss greifen (thin statt missing), got %+v", f)
	}
}

// Die Aktiv-Status-Prüfung bleibt von der Closure-Fähigkeit unberührt und
// umgekehrt — beide können nebeneinander Befunde liefern.
func TestClosureUnabhaengigVonAktivStatus(t *testing.T) {
	files := map[string]string{
		planRoadmap:                   roadmapText("## Aktuelle Welle", "Keine aktive Welle."),
		planSlice:                     "# aktiver Slice\n",
		closureDir + "/slice-001-a.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
	}
	cfg := closureCfg()
	drift := CheckPlanning(coretest.NewMemFS(files), cfg)
	closure := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(drift) != 1 || drift[0].Reason != model.ReasonPlanningDrift {
		t.Fatalf("Aktiv-Status muss unabhängig driften, got %+v", drift)
	}
	if len(closure) != 1 || closure[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("Closure muss unabhängig melden, got %+v", closure)
	}
}
