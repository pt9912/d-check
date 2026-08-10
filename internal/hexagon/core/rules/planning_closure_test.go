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

// Die Fence-Erkennung trimmt nur Space und Tab, nicht unicode-weit — identisch
// zur Vorverarbeitung (TrimFenceIndent). Sonst kippt eine mit U+00A0
// eingerückte Zeile die Fence-Parität NUR hier, und alles dahinter faellt aus
// der Messung: die Floskel waere unsichtbar, und weil dieselbe Zeile fuer das
// Modul spans kein Toggle ist, meldete auch dessen fence-unclosed nicht.
func TestClosureUnicodeWhitespaceIstKeineFenceEinrueckung(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\n\u00a0```\n\nwie besprochen umgesetzt.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	cfg := closureCfg()
	cfg.Closure.Boilerplate = []string{"wie besprochen"}
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteBoilerplate {
		t.Fatalf("U+00A0 darf die Fence-Paritaet nicht kippen ⇒ boilerplate erwartet, got %+v", f)
	}
}

// Dieselbe Trimmung gilt bei der Suche NACH der Ueberschrift: eine U+00A0-Zeile
// VOR ihr duerfte sie nicht in einen Fence-Zustand ziehen — sonst faende die
// Messung gar keinen Abschnitt und meldete missing statt zu messen.
func TestClosureUnicodeWhitespaceVorDerUeberschrift(t *testing.T) {
	body := "# Slice\n\n\u00a0```\n\n## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": body}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("U+00A0 vor der Ueberschrift darf sie nicht verdecken, got %+v", f)
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

// fail-closed: ein Closure-Verzeichnis OHNE Kandidaten ist kein „nichts zu tun",
// sondern ein leer laufendes Gate. `closure.dir` zu setzen ist die Behauptung,
// dass dort Notizen liegen; stimmt sie nicht mehr (Bestand umgezogen), wäre
// stilles Grün die schlechteste Antwort (R1-F-5).
func TestClosureOhneKandidatenFailClosed(t *testing.T) {
	files := map[string]string{closureDir + "/.keep": "", closureDir + "/README.md": "# Index\n"}
	f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("kandidatenfreies Closure-Verzeichnis muss fail-closed melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "leer") {
		t.Errorf("Meldung muss das Leerlaufen benennen, got %q", f[0].Message)
	}
}

// Gegenprobe zu F-1 (HIGH): eine Zeile wie `#1 …` ist Fließtext, keine H1 — sie
// darf den Abschnitt NICHT vorzeitig beenden. Vor dem Fix schnitt sie den Rest
// weg und ein Floskel-Treffer dahinter blieb unsichtbar (stilles Grün).
func TestClosureRautenTextBeendetAbschnittNicht(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nEin Satz. Zwei. Drei. Vier.\n\n#1 war ein Thema\n\nalles gut\n"
	files := map[string]string{closureDir + "/slice-001-a.md": "# Slice\n\n" + note}
	cfg := closureCfg()
	cfg.Closure.Boilerplate = []string{"alles gut"}
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteBoilerplate {
		t.Fatalf("`#1 …` ist keine Überschrift ⇒ Floskel dahinter muss gefunden werden, got %+v", f)
	}
}

// Und die Gegenrichtung: `#1 …` ist auch als Abschnitts-ERÖFFNER keine
// Überschrift (sonst matchte ein Muster darauf).
func TestClosureRautenTextIstKeinAbschnittsAnfang(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n#1 Closure-Notiz ist hier nicht\n\n## 7. Closure-Notiz\n\nEin Satz. Zwei. Drei. Vier.\n",
	}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("`#1 …` darf den Abschnitt nicht eröffnen, got %+v", f)
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

// reverseListFS liefert Verzeichnis-Einträge in UMGEKEHRTER Reihenfolge. Ohne
// so ein Double wäre der Determinismus-Test tautologisch: MemFS sortiert selbst,
// und das Entfernen der Sortierung im Kern bliebe grün (R2-F-7).
type reverseListFS struct{ driven.Filesystem }

func (r reverseListFS) List(dir string) ([]driven.DirEntry, error) {
	entries, err := r.Filesystem.List(dir)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries, nil
}

// Determinismus (DC-QA-02): die Befund-Reihenfolge hängt an der Sortierung im
// Kern, nicht an der Laune des Dateisystems — belegt mit einem Port, der
// absichtlich verkehrt herum liefert.
func TestClosureDeterministischeReihenfolge(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-003-c.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
		closureDir + "/slice-001-a.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
		closureDir + "/slice-002-b.md": "# S\n\n## 7. Closure-Notiz\n\nx.\n",
	}
	got := CheckPlanningClosure(reverseListFS{coretest.NewMemFS(files)}, closureCfg())
	want := []string{
		closureDir + "/slice-001-a.md",
		closureDir + "/slice-002-b.md",
		closureDir + "/slice-003-c.md",
	}
	if len(got) != len(want) {
		t.Fatalf("erwartet %d Befunde, got %+v", len(want), got)
	}
	for i := range want {
		if got[i].File != want[i] {
			t.Fatalf("Befund %d = %s, want %s (Kern muss sortieren, nicht der Port)", i, got[i].File, want[i])
		}
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

// Ohne closure.glob ist der Befundsatz byte-identisch zu einem Lauf, der
// ausdrücklich den slice-glob als Kandidaten-Filter setzt — der Default ist ein
// Verweis, kein zweites Muster.
func TestClosureGlobDefaultByteIdentisch(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote,
		closureDir + "/welle-01-x.md":  "# Welle\n\nOhne Notiz.\n",
	}
	ohne := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg())
	cfg := closureCfg()
	cfg.Closure.Glob = cfg.EffectiveSliceGlob()
	mit := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(ohne) != len(mit) {
		t.Fatalf("Befundsatz weicht ab: ohne=%+v mit=%+v", ohne, mit)
	}
	for i := range ohne {
		if ohne[i] != mit[i] {
			t.Fatalf("Befund %d weicht ab: %+v vs %+v", i, ohne[i], mit[i])
		}
	}
	if len(ohne) != 0 {
		t.Fatalf("die welle-Datei darf unter dem Slice-Glob nicht gesehen werden: %+v", ohne)
	}
}

// Der eigene Glob weitet die Kandidaten-Menge der Closure-Fähigkeit, OHNE die
// Grundmenge der Lifecycle-Invariante zu berühren — das ist der ganze Zweck der
// Entkopplung (ADR-0051).
func TestClosureGlobWeitetNurDieEigeneMenge(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote,
		closureDir + "/welle-01-x.md":  "# Welle\n\nOhne Notiz.\n",
	}
	cfg := closureCfg()
	cfg.Closure.Glob = "*.md"
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing ||
		f[0].File != closureDir+"/welle-01-x.md" {
		t.Fatalf("erwartet genau ein missing auf der welle-Datei, got %+v", f)
	}
	if cfg.EffectiveSliceGlob() != "slice-*.md" {
		t.Fatalf("der Slice-Glob wurde mitgezogen: %q", cfg.EffectiveSliceGlob())
	}
}

// Null Kandidaten unter dem GESETZTEN Glob ist fail-closed wie unter dem
// geerbten: den Filter zu setzen ist die Behauptung, dass er etwas trifft.
func TestClosureGlobNullKandidatenFailClosed(t *testing.T) {
	files := map[string]string{closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote}
	cfg := closureCfg()
	cfg.Closure.Glob = "paket-*.md"
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing || f[0].File != closureDir {
		t.Fatalf("erwartet fail-closed auf dem Verzeichnis, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "paket-*.md") {
		t.Errorf("die Meldung nennt den gesetzten Glob nicht: %q", f[0].Message)
	}
}

// placeholderCfg: Closure-Fähigkeit mit aktiver vierter Bedingung.
func placeholderCfg() model.PlanningConfig {
	c := closureCfg()
	c.Closure.Placeholder = true
	return c
}

// Ohne den Schalter ist der Befundsatz byte-identisch — der Template-Rumpf
// passiert alle drei bestehenden Bedingungen.
func TestClosurePlaceholderAusIstByteIdentisch(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nErgebnis: <ergebnis>. Belege: <belege>. Offen: <offen>. Ende: <ende>.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), closureCfg()); f != nil {
		t.Fatalf("ohne Schalter erwartet befundfrei, got %+v", f)
	}
}

// Der Template-Rumpf meldet GENAU EINEN Befund — den ersten Treffer, an seiner
// Zeile. Mehrere Platzhalter derselben Notiz sind dieselbe Reparatur.
func TestClosurePlaceholderErsterTreffer(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\nErgebnis: <ergebnis>. Belege: <belege>.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNotePlaceholder {
		t.Fatalf("erwartet genau ein placeholder, got %+v", f)
	}
	if f[0].Line != 5 {
		t.Errorf("Zeile = %d, want 5 (die Zeile des ersten Treffers)", f[0].Line)
	}
	if !strings.Contains(f[0].Message, "<ergebnis>") {
		t.Errorf("die Meldung nennt den Treffer nicht: %q", f[0].Message)
	}
}

// Die Falsch-Positiv-Klassen EINZELN, nicht als Sammel-Fixture: jede muss für
// sich grün sein, sonst verdeckt eine die andere.
func TestClosurePlaceholderFalschPositivKlassen(t *testing.T) {
	cases := map[string]string{
		"Vergleichszeichen":       "p95 < 1 s und Recall > 0,9 gemessen.",
		"Generic":                 "Der Puffer ist ein vector<float> geblieben.",
		"Autolink":                "Quelle: <https://example.org/a> geprüft.",
		"Mail-Adresse":            "Kontakt war <a@example.org> im Team.",
		"HTML-Tag mit Attribut":  "Der Anker <a id=\"x\"></a> steht im Register.",
		"Vergleich ohne Leerzeichen": "Die Latenz blieb <1 s und der Recall >0,9 im Median.",
		"Tabellenzelle":           "| Wert | <1 s | >0,9 |",
		"Winkelklammer-Linkziel":  "Siehe [dort](<docs/x.md>) im Register.",
		"Nicht-ASCII davor":       "Ergebnis — <2 s und >1 GB blieben stabil.",
		"HTML-Tag selbstschliessend": "Ein <br/> trennt die Zeilen.",
		"Inline-Code zeigt Syntax": "Die Vorlage nennt `<PREFIX>` und `<a id>` als Muster.",
		"Kommentar":               "Der Marker <!-- d-check:ignore --> bleibt stehen.",
		"schliessendes Tag":       "Das </div> gehört zum Markup.",
	}
	for name, satz := range cases {
		t.Run(name, func(t *testing.T) {
			note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\n" + satz + "\n"
			files := map[string]string{closureDir + "/slice-001-a.md": note}
			if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
				t.Fatalf("Falsch-Positiv: %+v", f)
			}
		})
	}
}

// Kombinierbar: eine duenne Notiz MIT Platzhalter meldet beides.
func TestClosurePlaceholderNebenThin(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nErgebnis: <ergebnis>.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	got := reasonsOf(CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()))
	if len(got) != 2 {
		t.Fatalf("erwartet thin UND placeholder, got %v", got)
	}
}

// Ein Fenced-Code-Block im Abschnitt zaehlt nicht — die Vorverarbeitung
// entfernt ihn, also kann er keinen Platzhalter beisteuern.
func TestClosurePlaceholderFenceZaehltNicht(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\n```text\nErgebnis: <ergebnis>.\n```\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
		t.Fatalf("Fence-Inhalt darf nicht zaehlen: %+v", f)
	}
}

// Der Ein-Zeichen-Fall und zwei benachbarte Platzhalter (deren konsumierte
// Vorzeichen sich ueberlappen) — gemeldet wird der erste.
func TestClosurePlaceholderRandformen(t *testing.T) {
	for name, satz := range map[string]string{
		"ein Zeichen": "Offen: <x>.",
		// Der erste Treffer wird vom HTML-Nachfilter verworfen; sein konsumiertes
		// Vorzeichen darf den direkt folgenden echten Platzhalter nicht verdecken.
		"verworfener Treffer verdeckt den naechsten nicht": "Offen: <a><ergebnis>.",
	} {
		t.Run(name, func(t *testing.T) {
			note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\n" + satz + "\n"
			files := map[string]string{closureDir + "/slice-001-a.md": note}
			f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg())
			if len(f) != 1 || f[0].Reason != model.ReasonClosureNotePlaceholder {
				t.Fatalf("erwartet genau ein placeholder, got %+v", f)
			}
		})
	}
}

// Ein Platzhalter AUSSERHALB des Abschnitts (davor wie danach) zaehlt nicht.
func TestClosurePlaceholderNurImAbschnitt(t *testing.T) {
	body := "# Slice\n\nDavor: <vorher>.\n\n## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\n" +
		"## 8. Danach\n\nDanach: <nachher>.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": body}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
		t.Fatalf("ausserhalb des Abschnitts darf nichts melden: %+v", f)
	}
}

// Jeder Eintrag der HTML-Tag-Liste muss WIRKEN: ohne diesen Test liessen sich
// fast alle streichen, ohne dass etwas rot wird (Review-Befund).
func TestClosurePlaceholderHtmlListeVollstaendigWirksam(t *testing.T) {
	for tag := range htmlTagNames() {
		t.Run(tag, func(t *testing.T) {
			note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\nDer <" + tag + "> blieb stehen.\n"
			files := map[string]string{closureDir + "/slice-001-a.md": note}
			if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
				t.Fatalf("<%s> ist ein HTML-Tag und darf nicht melden: %+v", tag, f)
			}
		})
	}
}

// Das gemeldete Ziel ist der Treffer, auf 40 Runen gekappt, und die
// Tag-Erkennung ist case-insensitiv.
func TestClosurePlaceholderZielUndSchreibweise(t *testing.T) {
	lang := "<" + strings.Repeat("f", 60) + ">"
	note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\nOffen: " + lang + ".\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg())
	if len(f) != 1 {
		t.Fatalf("erwartet genau einen Befund, got %+v", f)
	}
	if n := len([]rune(f[0].Message)); !strings.Contains(f[0].Message, "<"+strings.Repeat("f", 39)) || n > 200 {
		t.Errorf("Kappung auf 40 Runen nicht wirksam: %q", f[0].Message)
	}
	note = "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\nDer <SECTION> blieb stehen.\n"
	files = map[string]string{closureDir + "/slice-001-a.md": note}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
		t.Fatalf("Tag-Erkennung muss case-insensitiv sein: %+v", f)
	}
}

// Der Schraegstrich im Vorzeichen: ein Pfad-Fragment eroeffnet keinen
// Platzhalter.
func TestClosurePlaceholderSchraegstrichDavor(t *testing.T) {
	note := "## 7. Closure-Notiz\n\nEins. Zwei. Drei. Vier.\n\nDer Pfad a/<b> bleibt Pfad.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": note}
	if f := CheckPlanningClosure(coretest.NewMemFS(files), placeholderCfg()); f != nil {
		t.Fatalf("Schraegstrich davor darf nicht melden: %+v", f)
	}
}

