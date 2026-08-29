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
