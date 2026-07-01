package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// Modul planning (DC-FA-PLAN-001) end-to-end gegen ein echtes Temp-Repo — genau
// der Pfad, den `make planning-check` fährt (Image + read-only-Mount, hermetisch).

const planningConfigYAML = `modules: [links]
planning:
  roadmap: docs/plan/planning/in-progress/roadmap.md
`

// planningRepo legt ein Repo mit Roadmap (gegebener roadmap-Inhalt) + optional
// einem Slice an und schreibt die planning-Config.
func planningRepo(t *testing.T, dir, roadmap string, withSlice bool) {
	t.Helper()
	progress := filepath.Join(dir, "docs/plan/planning/in-progress")
	if err := os.MkdirAll(progress, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(progress, "roadmap.md"), roadmap)
	if withSlice {
		mustWrite(t, filepath.Join(progress, "slice-001-x.md"), "# slice\n")
	}
	mustWrite(t, filepath.Join(dir, ".d-check.yml"), planningConfigYAML)
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func roadmapMD(block string) string {
	return "# Roadmap\n\n## Aktuelle Welle\n\n" + block + "\n\n## Nächste Wellen\n\nKeine aktive Welle (Prosa).\n"
}

// TestPlanning_Consistent: aktive Welle + Slice ⇒ Exit 0.
func TestPlanning_Consistent(t *testing.T) {
	dir := t.TempDir()
	planningRepo(t, dir, roadmapMD("welle-x in Arbeit."), true)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--enable", "planning", dir}, &so, &se); code != 0 {
		t.Fatalf("Exit = %d, want 0\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
}

// TestPlanning_Drift: Slice vorhanden, aber Roadmap trägt den Ruhe-Marker ⇒
// planning-drift, Exit 1.
func TestPlanning_Drift(t *testing.T) {
	dir := t.TempDir()
	planningRepo(t, dir, roadmapMD("Keine aktive Welle — wartet auf Trigger."), true)
	var so, se bytes.Buffer
	code := cli.Run([]string{"--enable", "planning", dir}, &so, &se)
	if code != 1 || !strings.Contains(so.String(), "planning-drift") {
		t.Fatalf("Exit = %d, stdout = %q (planning-drift/Exit 1 erwartet)", code, so.String())
	}
}

// TestPlanning_FailClosed: eine kaputte §Aktuelle-Welle-Überschrift ⇒
// planning-drift, Exit 1 (kein stilles Grün).
func TestPlanning_FailClosed(t *testing.T) {
	dir := t.TempDir()
	broken := "# Roadmap\n\n## Aktuelle Wellen\n\nKeine aktive Welle.\n"
	planningRepo(t, dir, broken, true)
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--enable", "planning", dir}, &so, &se); code != 1 {
		t.Fatalf("kaputte Überschrift ⇒ Exit 1 erwartet, bekam %d\nstdout=%s", code, so.String())
	}
}
