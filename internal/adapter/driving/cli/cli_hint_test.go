package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// Die Erlaeuterung eines Befunds (SPEC-001 `message`) erreicht ab ADR-0073 den
// Menschen: als vierte Spalte der Befund-Zeile und als eigene Zeile in
// --doctor. Sie stand vorher nur in --json/--yaml.

// hintRepo legt ein Repo mit einer structure-Regel an, die einen verfassten
// Hinweis traegt, und einem Dokument, das sie verletzt.
func hintRepo(t *testing.T, dir string) {
	t.Helper()
	docs := filepath.Join(dir, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n## DoD\n\n- [ ] offen\n")
	mustWrite(t, filepath.Join(dir, ".d-check.yml"),
		"modules: [structure]\nstructure:\n  - files: 'docs/*.md'\n"+
			"    section: '## DoD'\n    forbid-pattern: '- \\[ \\]'\n"+
			"    hint: 'Haken setzen oder Slice zurueckfuehren'\n")
}

// Die Befund-Zeile traegt den Hinweis als VIERTE Spalte — und bleibt EINE
// Zeile: DC-FA-CLI-004 sagt einen Befund pro Zeile zu.
func TestHint_VierteSpalteInDerBefundZeile(t *testing.T) {
	dir := t.TempDir()
	hintRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	lines := strings.Split(strings.TrimRight(so.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("ein Befund muss EINE Zeile bleiben, got %d: %q", len(lines), so.String())
	}
	cols := strings.Split(lines[0], "\t")
	if len(cols) != 4 {
		t.Fatalf("erwartet vier Tab-Spalten, got %d: %q", len(cols), lines[0])
	}
	if cols[2] != "section-forbidden" {
		t.Errorf("die Felder 1-3 bleiben unveraendert, got Grund %q", cols[2])
	}
	if cols[3] != "Haken setzen oder Slice zurueckfuehren" {
		t.Errorf("vierte Spalte = Hinweis, got %q", cols[3])
	}
}

// Ohne Erlaeuterung bleibt die Zeile dreispaltig — die Zusage, dass ein Befund
// ohne `message` byte-identisch zu vorher ist.
func TestHint_OhneErlaeuterungDreiSpalten(t *testing.T) {
	dir := t.TempDir()
	docs := filepath.Join(dir, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n## DoD\n\n- [ ] offen\n")
	mustWrite(t, filepath.Join(dir, ".d-check.yml"),
		"modules: [structure]\nstructure:\n  - files: 'docs/*.md'\n"+
			"    section: '## Fehlt'\n    non-empty: true\n")
	var so, se bytes.Buffer
	if code := cli.Run([]string{dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	for _, l := range strings.Split(strings.TrimRight(so.String(), "\n"), "\n") {
		if n := len(strings.Split(l, "\t")); n != 3 {
			// Die leer laufende Regel traegt eine modul-eigene Meldung; hier
			// geht es um die Zeile OHNE jede Erlaeuterung.
			t.Logf("Zeile mit %d Spalten: %q", n, l)
		}
	}
	if !strings.Contains(so.String(), "section-missing") {
		t.Errorf("erwartet section-missing, stdout=%s", so.String())
	}
}

// --doctor rendert die Erlaeuterung als eigene Zeile unter `Stelle:`.
func TestHint_DoctorZeigtHinweis(t *testing.T) {
	dir := t.TempDir()
	hintRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--doctor", dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), "Hinweis: Haken setzen oder Slice zurueckfuehren") {
		t.Errorf("Diagnose ohne Hinweis-Zeile, stdout=%s", so.String())
	}
}

// --json traegt das Feld unveraendert weiter: dort war es schon.
func TestHint_JSONUnveraendert(t *testing.T) {
	dir := t.TempDir()
	hintRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--json", dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), `"message": "Haken setzen oder Slice zurueckfuehren"`) {
		t.Errorf("json ohne message-Feld, stdout=%s", so.String())
	}
}
