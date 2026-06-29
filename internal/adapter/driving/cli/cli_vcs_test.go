package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// Modul vcs (DC-FA-VCS-001) end-to-end gegen ein echtes On-Disk-git-Repo —
// genau der Pfad, den `make adr-check` später fährt (Image + read-only .git).

const vcsConfig = `modules: [links]
vcs:
  paths: ["adr-*.md"]
  immutable-when: '^\*\*Status:\*\* Accepted'
  exclude-sections: [Geschichte]
  status-line: '^\*\*Status:\*\*'
  head-allow: '^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})'
`

func adrText(status, decision string) string {
	return "# ADR-0099 — X\n\n**Status:** " + status + "\n\n## Entscheidung\n\n" +
		decision + "\n\n## Geschichte\n\n| Datum | Ereignis |\n|---|---|\n| 2026-01-01 | x |\n"
}

// initVCSRepo legt ein git-Repo unter dir an und committet die Dateien.
func initVCSRepo(t *testing.T, dir string) *gogit.Worktree {
	t.Helper()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	return wt
}

func writeAt(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitAll(t *testing.T, wt *gogit.Worktree, msg string) string {
	t.Helper()
	if err := wt.AddWithOptions(&gogit.AddOptions{All: true}); err != nil {
		t.Fatal(err)
	}
	h, err := wt.Commit(msg, &gogit.CommitOptions{
		Author: &object.Signature{Name: "T", Email: "t@example.com", When: time.Unix(1700000000, 0)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return h.String()
}

// TestVCS_RangeDrift: ein Körper-Edit an einer Accepted-ADR über die Range ⇒
// core-drift-vcs, Exit 1 (deckt setupVCS/Open/ChangedPaths/FileAt/CheckVCS ab).
func TestVCS_RangeDrift(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", vcsConfig)
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
	c1 := commitAll(t, wt, "c1")
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue B."))
	c2 := commitAll(t, wt, "c2")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "core-drift-vcs") {
		t.Fatalf("stdout ohne core-drift-vcs: %s", stdout.String())
	}
}

// TestVCS_RangeClean: nur ein Anhang unter ## Geschichte ⇒ kein Befund, Exit 0.
func TestVCS_RangeClean(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", vcsConfig)
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
	c1 := commitAll(t, wt, "c1")
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A.")+"| 2026-02-02 | Notiz |\n")
	c2 := commitAll(t, wt, "c2")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// TestVCS_Staged: ein staged Körper-Edit ⇒ core-drift-vcs, Exit 1.
func TestVCS_Staged(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", vcsConfig)
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
	commitAll(t, wt, "c1")
	// staged Änderung ohne Commit.
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue B."))
	if _, err := wt.Add("adr-x.md"); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "vcs", "--staged", dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// TestVCS_FailClosed deckt die fail-closed-Pfade ab: vcs aktiv ohne Range,
// mit kaputter Range und ohne .git ⇒ jeweils Exit 2.
func TestVCS_FailClosed(t *testing.T) {
	withGit := t.TempDir()
	wt := initVCSRepo(t, withGit)
	writeAt(t, withGit, ".d-check.yml", vcsConfig)
	writeAt(t, withGit, "adr-x.md", adrText("Accepted", "Tue A."))
	commitAll(t, wt, "c1")

	noGit := t.TempDir()
	writeAt(t, noGit, ".d-check.yml", vcsConfig)
	writeAt(t, noGit, "adr-x.md", adrText("Accepted", "Tue A."))

	cases := []struct {
		name string
		args []string
		root string
	}{
		{"ohne range/staged", []string{"--enable", "vcs"}, withGit},
		{"kaputte range ohne ..", []string{"--enable", "vcs", "--range", "deadbeef"}, withGit},
		{"unauflösbare basis", []string{"--enable", "vcs", "--range", "0000000..HEAD"}, withGit},
		{"ohne .git", []string{"--enable", "vcs", "--staged"}, noGit},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(append(c.args, c.root), &stdout, &stderr)
			if code != 2 {
				t.Fatalf("Exit = %d, want 2 (fail-closed)\nstderr=%s", code, stderr.String())
			}
		})
	}
}

// TestVCS_RangeStagedExklusiv: --range und --staged zusammen ⇒ Nutzungsfehler.
func TestVCS_RangeStagedExklusiv(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "vcs", "--range", "a..b", "--staged", t.TempDir()}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("Exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "nicht kombinierbar") {
		t.Fatalf("stderr ohne Kombinations-Fehler: %s", stderr.String())
	}
}
