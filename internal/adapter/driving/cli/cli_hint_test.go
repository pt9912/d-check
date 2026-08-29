package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// Die Erlaeuterung eines Befunds (SPEC-001 `message`) erreicht den Menschen als
// vierte Spalte der Befund-Zeile und als eigene Zeile in --doctor; --json und
// --yaml tragen sie als Feld (ADR-0073).

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


// Ohne Erläuterung bleibt die Zeile DREISPALTIG. Das Modul `spans` meldet
// `fence-unclosed` ohne `message` — die Gegenprobe zur vierten Spalte.
func TestHint_OhneErlaeuterungDreiSpalten(t *testing.T) {
	dir := t.TempDir()
	docs := filepath.Join(dir, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n```go\nnie geschlossen\n")
	mustWrite(t, filepath.Join(dir, ".d-check.yml"), "modules: [spans]\n")
	var so, se bytes.Buffer
	if code := cli.Run([]string{dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	lines := strings.Split(strings.TrimRight(so.String(), "\n"), "\n")
	if len(lines) == 0 || so.String() == "" {
		t.Fatalf("erwartet mindestens einen Befund, stdout=%q", so.String())
	}
	for _, l := range lines {
		if n := len(strings.Split(l, "\t")); n != 3 {
			t.Errorf("Befund ohne Erläuterung muss dreispaltig sein, got %d: %q", n, l)
		}
	}
}

// Der `hint` erreicht `--json` — mit seinem WERT, nicht nur als Schlüssel: das
// Feld gab es schon, sein Inhalt ist jetzt verfasst.
func TestHint_JSONTraegtDenVerfasstenText(t *testing.T) {
	dir := t.TempDir()
	hintRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--json", dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), `"message": "Haken setzen oder Slice zurueckfuehren"`) {
		t.Errorf("json ohne verfasste Erläuterung, stdout=%s", so.String())
	}
}

// Ein `hint` mit Tab oder Zeilenumbruch bräche die Zeilen- und Feldgrenze des
// Ausgabeformats. Die Konfiguration weist ihn ab, statt ihn still zu
// begradigen (DC-FA-CLI-004, ADR-0073).
func TestHint_TabUndUmbruchWerdenAbgewiesen(t *testing.T) {
	for name, hint := range map[string]string{
		"Zeilenumbruch": "Zeile eins\\nZeile zwei",
		"Tabulator":     "Feld4\\tFeld5",
	} {
		dir := t.TempDir()
		docs := filepath.Join(dir, "docs")
		if err := os.MkdirAll(docs, 0o755); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n## DoD\n\n- [ ] offen\n")
		mustWrite(t, filepath.Join(dir, ".d-check.yml"),
			"modules: [structure]\nstructure:\n  - files: 'docs/*.md'\n"+
				"    section: '## DoD'\n    forbid-pattern: '- \\[ \\]'\n"+
				"    hint: \""+hint+"\"\n")
		var so, se bytes.Buffer
		if code := cli.Run([]string{dir}, &so, &se); code != 2 {
			t.Errorf("%s: Exit = %d, want 2\nstdout=%s\nstderr=%s", name, code, so.String(), se.String())
		}
	}
}

// Eine Erläuterung, die NICHT aus der Konfiguration kommt, kann der Config-Rand
// nicht abweisen — `commits` trägt den Commit-Betreff. Der Reporter macht
// daraus ein einzeiliges Feld, damit ein Befund eine Zeile bleibt.
func TestHint_ReporterMachtEinzeiligesFeld(t *testing.T) {
	dir := t.TempDir()
	docs := filepath.Join(dir, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n## DoD\n\n- [ ] offen\n")
	mustWrite(t, filepath.Join(dir, ".d-check.yml"),
		"modules: [structure]\nstructure:\n  - files: 'docs/*.md'\n"+
			"    section: '## DoD'\n    forbid-pattern: '- \\[ \\]'\n"+
			"    hint: 'Haken setzen'\n")
	var so, se bytes.Buffer
	if code := cli.Run([]string{dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	lines := strings.Split(strings.TrimRight(so.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("ein Befund muss EINE Zeile bleiben, got %d: %q", len(lines), so.String())
	}
	if n := len(strings.Split(lines[0], "\t")); n != 4 {
		t.Fatalf("erwartet vier Felder, got %d: %q", n, lines[0])
	}
}
