package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// --config (DC-FA-CLI-012) end-to-end gegen ein echtes Temp-Repo: der Schalter
// verschiebt NUR die Herkunft der Konfiguration — er ersetzt, ergänzt nicht,
// bleibt innerhalb der Scan-Wurzel und fällt nie still zurück.

// configPathRepo legt ein Repo mit einem kaputten Link an, dazu zwei Profile:
// die konventionelle Datei schaltet `links` AUS (Lauf wäre grün), die
// alternative schaltet `links` AN (Lauf wäre rot). So ist an Exit-Code
// ablesbar, welches Profil tatsächlich galt.
func configPathRepo(t *testing.T, dir string) {
	t.Helper()
	docs := filepath.Join(dir, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(docs, "a.md"), "# A\n\n[kaputt](fehlt.md)\n")
	mustWrite(t, filepath.Join(dir, ".d-check.yml"), "modules: []\n")
	mustWrite(t, filepath.Join(dir, "profil.yml"), "modules: [links]\n")
}

// Happy Path: die alternative Datei gilt, die konventionelle wird NICHT gelesen
// (sonst bliebe der Lauf grün).
func TestConfigPath_Ersetzt(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--config", "profil.yml", dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1 (alternatives Profil aktiviert links)\nstdout=%s\nstderr=%s",
			code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), "target-missing") {
		t.Errorf("erwartet target-missing aus dem alternativen Profil, stdout=%s", so.String())
	}
}

// Gegenprobe: ohne den Schalter gilt die konventionelle Datei (Lauf grün).
// Zusammen mit dem Fall darüber belegt das, dass der Schalter wirkt.
func TestConfigPath_OhneFlagUnveraendert(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{dir}, &so, &se); code != 0 {
		t.Fatalf("Exit = %d, want 0 (konventionelle Datei schaltet links aus)\nstdout=%s\nstderr=%s",
			code, so.String(), se.String())
	}
}

// fail-closed: fehlende Datei ⇒ Exit 2, KEIN Rückfall auf die konventionelle
// Datei (die hier grün wäre) und keiner auf die Defaults.
func TestConfigPath_FehltFailClosed(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	var so, se bytes.Buffer
	code := cli.Run([]string{"--config", "gibtsnicht.yml", dir}, &so, &se)
	if code != 2 {
		t.Fatalf("Exit = %d, want 2\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(se.String(), "gibtsnicht.yml") {
		t.Errorf("stderr muss den getippten Pfad nennen, got %s", se.String())
	}
}

// fail-closed: ein Pfad außerhalb der Scan-Wurzel ⇒ Exit 2 (der read-only-Mount
// ist die Grenze).
func TestConfigPath_AusserhalbDerWurzel(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	for _, bad := range []string{"../raus.yml", "../../etc/passwd"} {
		var so, se bytes.Buffer
		if code := cli.Run([]string{"--config", bad, dir}, &so, &se); code != 2 {
			t.Errorf("--config %q: Exit = %d, want 2\nstderr=%s", bad, code, se.String())
		}
	}
}

// Ein absoluter Pfad INNERHALB der Wurzel ist zulässig (auf die Wurzel bezogen).
func TestConfigPath_AbsolutInnerhalbErlaubt(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	var so, se bytes.Buffer
	abs := filepath.Join(dir, "profil.yml")
	if code := cli.Run([]string{"--config", abs, dir}, &so, &se); code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
}

// Ein Verzeichnis ist keine Konfigurationsdatei ⇒ Exit 2.
func TestConfigPath_VerzeichnisFailClosed(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--config", "docs", dir}, &so, &se); code != 2 {
		t.Fatalf("Exit = %d, want 2\nstderr=%s", code, se.String())
	}
}

// Ungültiges Schema in der alternativen Datei ⇒ Exit 2 wie bei der
// konventionellen (DC-FA-CONF-001; die Validierung ist unverändert).
func TestConfigPath_UngueltigesSchema(t *testing.T) {
	dir := t.TempDir()
	configPathRepo(t, dir)
	mustWrite(t, filepath.Join(dir, "kaputt.yml"), "modules: [gibtsnicht]\n")
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--config", "kaputt.yml", dir}, &so, &se); code != 2 {
		t.Fatalf("Exit = %d, want 2\nstderr=%s", code, se.String())
	}
}

// Das Closure-Profil end-to-end: genau der Pfad, den `make verify-closure-notes`
// fährt — eigene Config-Datei, nur `planning`, Befund über den done/-Bestand.
func TestConfigPath_ClosureProfil(t *testing.T) {
	dir := t.TempDir()
	done := filepath.Join(dir, "docs/plan/planning/done")
	progress := filepath.Join(dir, "docs/plan/planning/in-progress")
	if err := os.MkdirAll(done, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(progress, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(progress, "roadmap.md"), roadmapMD("Keine aktive Welle."))
	mustWrite(t, filepath.Join(done, "slice-001-x.md"),
		"# Slice\n\n## 7. Closure-Notiz\n\n_Ausstehend._\n")
	mustWrite(t, filepath.Join(dir, ".d-check.closure.yml"), `modules: []
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
  closure:
    dir: docs/plan/planning/done
`)
	var so, se bytes.Buffer
	code := cli.Run([]string{"--config", ".d-check.closure.yml", "--enable", "planning", dir}, &so, &se)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), "closure-note-thin") {
		t.Errorf("erwartet closure-note-thin, stdout=%s", so.String())
	}
}
