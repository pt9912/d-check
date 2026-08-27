package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// Modul commits (DC-FA-COMMITS-001) end-to-end gegen ein echtes On-Disk-git-Repo —
// genau der Pfad, den `make trace-check` fährt (Image + read-only .git / stdin).
// Nutzt die Helfer aus cli_vcs_test.go (initVCSRepo/writeAt/commitAll).

const commitsConfigYAML = `modules: [links]
commits:
  id-patterns:
    - 'ADR-\d{4}'
    - 'slice-\d+'
  exempt-pattern: '^(Merge |Revert )'
`

// TestCommits_RangeUntraceable: ein Nicht-Merge-Commit ohne Kennung in der Range
// ⇒ commit-untraceable, Exit 1 (deckt resolveVCS/CommitMessages/CheckCommits ab).
func TestCommits_RangeUntraceable(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "eins\n")
	c1 := commitAll(t, wt, "c1: feat ADR-0001")
	writeAt(t, dir, "f.md", "zweite fassung\n")
	c2 := commitAll(t, wt, "c2: chore ohne bezug")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "commits", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "commit-untraceable") {
		t.Fatalf("stdout ohne commit-untraceable: %s", stdout.String())
	}
}

// TestCommits_RangeClean: trägt der Range-Commit eine Kennung ⇒ kein Befund, Exit 0.
func TestCommits_RangeClean(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "eins\n")
	c1 := commitAll(t, wt, "c1: init")
	writeAt(t, dir, "f.md", "zweite fassung\n")
	c2 := commitAll(t, wt, "c2: docs slice-056 Body")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "commits", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// TestCommits_CommitMsgFile: --commit-msg <datei> — der Kurzschluss-Modus des Hooks
// (ohne Range/Scan). Kennungslose Message ⇒ Exit 1, Message mit Kennung ⇒ Exit 0.
func TestCommits_CommitMsgFile(t *testing.T) {
	dir := t.TempDir()
	initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "x\n")

	bad := filepath.Join(dir, "MSG_BAD")
	if err := os.WriteFile(bad, []byte("chore: kein bezug\n# ADR-0001 nur im Kommentar\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--commit-msg", bad, dir}, &so, &se); code != 1 {
		t.Fatalf("bad: Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), "commit-untraceable") {
		t.Fatalf("bad: stdout ohne commit-untraceable: %s", so.String())
	}

	ok := filepath.Join(dir, "MSG_OK")
	if err := os.WriteFile(ok, []byte("feat: DC-FA-COMMITS-001? nein, slice-056\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	so.Reset()
	se.Reset()
	if code := cli.Run([]string{"--commit-msg", ok, dir}, &so, &se); code != 0 {
		t.Fatalf("ok: Exit = %d, want 0\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
}

// TestCommits_CommitMsgExempt: ein Merge-Betreff ist kennungsfrei erlaubt (Exit 0).
func TestCommits_CommitMsgExempt(t *testing.T) {
	dir := t.TempDir()
	initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "x\n")
	msg := filepath.Join(dir, "MSG")
	if err := os.WriteFile(msg, []byte("Merge branch 'feature'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var so, se bytes.Buffer
	if code := cli.Run([]string{"--commit-msg", msg, dir}, &so, &se); code != 0 {
		t.Fatalf("Merge: Exit = %d, want 0\nstderr=%s", code, se.String())
	}
}

// TestCommits_CommitMsgStdin: --commit-msg - liest die Pending-Message von stdin
// (der Pfad, den der commit-msg-Hook über `docker run -i` pipet).
func TestCommits_CommitMsgStdin(t *testing.T) {
	dir := t.TempDir()
	initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "x\n")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.WriteString("chore: kein bezug\n"); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	var so, se bytes.Buffer
	if code := cli.Run([]string{"--commit-msg", "-", dir}, &so, &se); code != 1 {
		t.Fatalf("stdin: Exit = %d, want 1\nstdout=%s\nstderr=%s", code, so.String(), se.String())
	}
	if !strings.Contains(so.String(), "commit-untraceable") {
		t.Fatalf("stdin: stdout ohne commit-untraceable: %s", so.String())
	}
}

// TestCommits_FailClosed: commits aktiv ohne Range ⇒ Exit 2; --commit-msg ohne
// konfigurierte id-patterns ⇒ Exit 2; --commit-msg + --range ⇒ Nutzungsfehler.
func TestCommits_FailClosed(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", commitsConfigYAML)
	writeAt(t, dir, "f.md", "x\n")
	commitAll(t, wt, "c1: init ADR-0001")

	noCommits := t.TempDir()
	writeAt(t, noCommits, ".d-check.yml", "modules: [links]\n")
	writeAt(t, noCommits, "f.md", "x\n")
	msg := filepath.Join(noCommits, "MSG")
	if err := os.WriteFile(msg, []byte("chore: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		args []string
	}{
		{"commits ohne range", []string{"--enable", "commits", dir}},
		{"commit-msg ohne id-patterns", []string{"--commit-msg", msg, noCommits}},
		{"commit-msg + range exklusiv", []string{"--commit-msg", msg, "--range", "a..b", dir}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var so, se bytes.Buffer
			if code := cli.Run(c.args, &so, &se); code != 2 {
				t.Fatalf("Exit = %d, want 2 (fail-closed)\nstderr=%s", code, se.String())
			}
		})
	}
}
