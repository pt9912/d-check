package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// DC-FA-TRK-001 Negative (E2E): ein existierendes, aber untracktes Linkziel
// ⇒ target-untracked, Exit 1 (links bleibt still — das Ziel existiert).
func TestTracked_UntrackedZiel(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, "doc.md", "# D\n\n[u](u.md)\n")
	commitAll(t, wt, "c1")
	writeAt(t, dir, "u.md", "# U\n")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "target-untracked") {
		t.Fatalf("stdout ohne target-untracked: %s", stdout.String())
	}
}

// DC-FA-TRK-001 Boundary (Index-Wahrheit, E2E): dasselbe Ziel frisch gestagt
// (nie committet) ⇒ kein Befund, Exit 0.
func TestTracked_StagedZielGruen(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, "doc.md", "# D\n\n[u](u.md)\n")
	commitAll(t, wt, "c1")
	writeAt(t, dir, "u.md", "# U\n")
	if _, err := wt.Add("u.md"); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (gestagt = getrackt)\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// DC-FA-TRK-001 Ventil (E2E): tracked.exempt-targets nimmt das aufgelöste
// Ziel referenz-weit aus.
func TestTracked_ExemptTargets(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", "tracked:\n  exempt-targets: [\"u.md\"]\n")
	writeAt(t, dir, "doc.md", "# D\n\n[u](u.md)\n")
	commitAll(t, wt, "c1")
	writeAt(t, dir, "u.md", "# U\n")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (u.md exempt)\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// DC-FA-CONF-001 Negative: ein ungültiges exempt-targets-Glob ⇒ Exit 2
// (config-zeitig, fail-closed — kein still wirkungsloses Ventil). Die
// Validierung läuft SEGMENTWEISE wie die Laufzeit (matchGlob) — ein nur
// als Ganzes gültiges Muster mit kaputtem Segment fällt ebenso durch.
func TestTracked_UngueltigesGlobExit2(t *testing.T) {
	for _, glob := range []string{"[", "[a/b].md"} {
		dir := t.TempDir()
		initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", "tracked:\n  exempt-targets: [\""+glob+"\"]\n")
		writeAt(t, dir, "doc.md", "x\n")

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr)
		if code != 2 || !strings.Contains(stderr.String(), "exempt-targets") {
			t.Fatalf("Glob %q: Exit = %d, stderr = %q, want 2 + exempt-targets-Hinweis", glob, code, stderr.String())
		}
	}
}

// DC-FA-TRK-001 fail-closed (E2E): aktives tracked ohne .git ⇒ Exit 2 mit
// Hinweis — kein stilles Grün.
func TestTracked_FailClosedOhneGit(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// write() aus cli_acceptance_test braucht root-relative Pfade; hier reicht
	// eine Datei über writeAt (kein git-Init!).
	writeAt(t, dir, "doc.md", "# D\n")
	code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "git") {
		t.Fatalf("Exit = %d, stderr = %q, want 2 + git-Hinweis", code, stderr.String())
	}
}

// tracked braucht KEINE Range: --enable tracked ohne --range/--staged ist
// gültig (anders als vcs/commits) — und umgekehrt erzwingt ein zusätzlich
// aktives vcs die Range weiterhin.
func TestTracked_OhneRangeGueltigVcsWeiterMitRange(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, "doc.md", "# D\n")
	commitAll(t, wt, "c1")

	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--enable", "tracked", dir}, &stdout, &stderr); code != 0 {
		t.Fatalf("tracked ohne Range: Exit = %d, want 0\nstderr=%s", code, stderr.String())
	}
	stderr.Reset()
	if code := cli.Run([]string{"--enable", "tracked", "--enable", "vcs", dir}, &stdout, &stderr); code != 2 ||
		!strings.Contains(stderr.String(), "--range") {
		t.Fatalf("vcs+tracked ohne Range: Exit/stderr = %d/%q, want 2 + Range-Hinweis", code, stderr.String())
	}
}

// DC-FA-CLI-010 (9→10): das --print-mk-Fragment trägt doc-tracked — mit
// --enable tracked, fokussierter --disable-Liste, OHNE Range und ohne
// --disable tracked.
func TestPrintMk_DocTracked(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := cli.Run([]string{"--print-mk"}, &stdout, &stderr); code != 0 {
		t.Fatalf("Exit = %d\nstderr=%s", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "\ndoc-tracked: ## ") {
		t.Fatalf("d-check.mk ohne doc-tracked-Target:\n%s", out)
	}
	var line string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "--enable tracked") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatalf("doc-tracked-Recipe ohne --enable tracked:\n%s", out)
	}
	for _, want := range []string{"--disable links", "--disable vcs", "--disable planning"} {
		if !strings.Contains(line, want) {
			t.Fatalf("doc-tracked-Recipe ohne %q:\n%s", want, line)
		}
	}
	if strings.Contains(line, "--range") || strings.Contains(line, "--staged") {
		t.Fatalf("doc-tracked trägt fälschlich eine Range:\n%s", line)
	}
	if strings.Contains(line, "--disable tracked") {
		t.Fatalf("doc-tracked wählt tracked fälschlich ab:\n%s", line)
	}
}
