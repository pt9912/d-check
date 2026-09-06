package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/report"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Der Reporter gibt eine Zusage über einen Text, den er nicht kontrolliert: die
// Erläuterung kommt aus der Konfiguration oder aus dem geprüften Material (das
// Modul `commits` trägt den Commit-Betreff). Tab und Zeilenumbruch würden die
// Feld- und Zeilengrenze verschieben — DC-FA-CLI-004 sagt EINEN Befund pro
// Zeile zu (ADR-0073).
func TestText_ErlaeuterungBleibtEinFeld(t *testing.T) {
	for name, msg := range map[string]string{
		"Zeilenumbruch": "Zeile eins\nZeile zwei",
		"Wagenruecklauf": "Zeile eins\r\nZeile zwei",
		"Tabulator":     "Feld4\tFeld5",
	} {
		var so, se bytes.Buffer
		f := model.Finding{
			File: "docs/a.md", Line: 3, Rule: "commits",
			Target: "abc1234", Reason: "commit-untraceable", Message: msg,
		}
		if err := report.Text(&so, &se, []model.Finding{f}, report.Summary{FilesChecked: 1, FindingCount: 1}); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		lines := strings.Split(strings.TrimRight(so.String(), "\n"), "\n")
		if len(lines) != 1 {
			t.Errorf("%s: ein Befund muss EINE Zeile bleiben, got %d: %q", name, len(lines), so.String())
			continue
		}
		if n := len(strings.Split(lines[0], "\t")); n != 4 {
			t.Errorf("%s: erwartet vier Felder, got %d: %q", name, n, lines[0])
		}
	}
}

// Ohne Erläuterung bleibt die Zeile dreispaltig — die Gegenprobe.
func TestText_OhneErlaeuterungDreiFelder(t *testing.T) {
	var so, se bytes.Buffer
	f := model.Finding{
		File: "docs/a.md", Line: 3, Rule: "spans",
		Target: "docs/a.md", Reason: "fence-unclosed",
	}
	if err := report.Text(&so, &se, []model.Finding{f}, report.Summary{FilesChecked: 1, FindingCount: 1}); err != nil {
		t.Fatal(err)
	}
	line := strings.TrimRight(so.String(), "\n")
	if n := len(strings.Split(line, "\t")); n != 3 {
		t.Fatalf("erwartet drei Felder, got %d: %q", n, line)
	}
}

// summary.notes ist ein NEUER, generischer Vertrag (DC-FA-MENT-001,
// DC-FA-CLI-004: summary traegt MINDESTENS filesChecked und findingCount,
// weitere Felder sind ausdruecklich zugelassen). Im Default-Modus erscheint
// jede Note als eigene stderr-Zeile VOR der Zaehlzeile.
func TestText_NotesAufStderrVorDerZaehlzeile(t *testing.T) {
	var out, errb bytes.Buffer
	sum := report.Summary{FilesChecked: 3, FindingCount: 0,
		Notes: []string{"mentions: 2 von 2 Artefakt(en) erwähnt, über 1 Dokument(e)"}}
	if err := report.Text(&out, &errb, nil, sum); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	lines := strings.Split(strings.TrimRight(errb.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("erwartet Note + Zaehlzeile, got %q", errb.String())
	}
	if !strings.Contains(lines[0], "2 von 2") {
		t.Fatalf("erste stderr-Zeile ist nicht die Note: %q", lines[0])
	}
	if !strings.Contains(lines[1], "3 Datei(en) geprüft") {
		t.Fatalf("zweite stderr-Zeile ist nicht die Zaehlzeile: %q", lines[1])
	}
	if out.Len() != 0 {
		t.Fatalf("die Note gehoert auf stderr, nicht auf stdout: %q", out.String())
	}
}

// Eine Note darf die Zeilengrenze nicht verschieben — dieselbe Zusage wie fuer
// die Erlaeuterung eines Befundes (ADR-0073).
func TestText_NoteBleibtEineZeile(t *testing.T) {
	var out, errb bytes.Buffer
	sum := report.Summary{Notes: []string{"erste\nzweite\tdritte"}}
	if err := report.Text(&out, &errb, nil, sum); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if n := strings.Count(strings.TrimRight(errb.String(), "\n"), "\n"); n != 1 {
		t.Fatalf("erwartet Note + Zaehlzeile = 2 Zeilen, got %q", errb.String())
	}
}

// RUECKWAERTSKOMPATIBILITAET, und sie haengt an einem einzigen `omitempty`:
// Ohne Note darf in der JSON-/YAML-Ausgabe KEIN notes-Schluessel erscheinen
// (DC-FA-MENT-001, Kriterium "Default byte-identisch"). Fiele das Tag weg,
// truege jede Ausgabe des Produkts "notes": null, und niemand meldete es.
func TestJSONYAML_OhneNoteKeinNotesSchluessel(t *testing.T) {
	for name, f := range map[string]func(*bytes.Buffer, []model.Finding, report.Summary) error{
		"JSON": func(b *bytes.Buffer, fs []model.Finding, s report.Summary) error { return report.JSON(b, fs, s, 0) },
		"YAML": func(b *bytes.Buffer, fs []model.Finding, s report.Summary) error { return report.YAML(b, fs, s, 0) },
	} {
		var b bytes.Buffer
		if err := f(&b, nil, report.Summary{FilesChecked: 1}); err != nil {
			t.Fatalf("%s: unerwarteter Fehler: %v", name, err)
		}
		if strings.Contains(b.String(), "notes") {
			t.Fatalf("%s: notes-Schluessel ohne Note vorhanden: %s", name, b.String())
		}
	}
}

// Und mit Note steht sie in BEIDEN maschinenlesbaren Formen (DC-FA-MENT-001:
// "in beiden Ausgabe-Formen", ausdrueckliche Festlegung).
func TestJSONYAML_MitNoteImSummary(t *testing.T) {
	sum := report.Summary{FilesChecked: 1, Notes: []string{"mentions: 1 von 1"}}
	var j, y bytes.Buffer
	if err := report.JSON(&j, nil, sum, 0); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if err := report.YAML(&y, nil, sum, 0); err != nil {
		t.Fatalf("YAML: %v", err)
	}
	for name, got := range map[string]string{"JSON": j.String(), "YAML": y.String()} {
		if !strings.Contains(got, "notes") || !strings.Contains(got, "1 von 1") {
			t.Fatalf("%s: Note fehlt im summary: %s", name, got)
		}
	}
}
