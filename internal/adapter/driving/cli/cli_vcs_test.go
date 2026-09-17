package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
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

// vcsConfigSub ist dieselbe Klasse wie vcsConfig, aber mit "**/" davor — die
// ADR liegt in einem Unterverzeichnis, dessen Tree unlesbar gemacht wird
// (TestVCS_UnlesbareUnterbaeume).
const vcsConfigSub = `modules: [links]
vcs:
  paths: ["**/adr-*.md"]
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
	full := filepath.Join(dir, name)
	vor := -1
	if b, err := os.ReadFile(full); err == nil {
		vor = len(b)
	}
	if err := coretest.GitFixtureRewriteHazard(name, vor, len(content)); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
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
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue B, zweite Fassung."))
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
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue B, zweite Fassung."))
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

// treeEntryHash löst den Tree-Eintrag path an ref auf und liefert seinen Hash.
func treeEntryHash(t *testing.T, dir, ref, path string) plumbing.Hash {
	t.Helper()
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	h, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		t.Fatal(err)
	}
	c, err := repo.CommitObject(*h)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := c.Tree()
	if err != nil {
		t.Fatal(err)
	}
	e, err := tree.FindEntry(path)
	if err != nil {
		t.Fatal(err)
	}
	return e.Hash
}

// removeLooseObject entfernt das lose Objekt zu h — dieselbe Wirkung wie ein
// Pack unter unkanonischem Namen (CO-001), nur ohne Pack-Mechanik im Test.
func removeLooseObject(t *testing.T, dir string, h plumbing.Hash) {
	t.Helper()
	s := h.String()
	if err := os.Remove(filepath.Join(dir, ".git", "objects", s[:2], s[2:])); err != nil {
		t.Fatal(err)
	}
}

// TestVCS_UnlesbareObjekte deckt CO-001s vier Ausprägungen End-to-End gegen
// dasselbe Muster von Probe-Repo ab (slice-220, DoD 2): BASE-Blob, BASE-Tree
// mit Pendant, BASE-Tree ohne Pendant und HEAD-Tree unlesbar. Vor slice-220
// meldeten die ersten beiden zwar schon Exit 2 (slice-218), die dritte blieb
// `0 Befund(e)`/Exit 0, und die vierte meldete fälschlich `core-drift-vcs`
// mit Exit 1 statt eines Umgebungsfehlers — der neue Entwurf (die geschützte
// Klasse direkt gegen beide Trees aufgelöst) schließt alle vier auf demselben
// Codepfad.
func TestVCS_UnlesbareObjekte(t *testing.T) {
	t.Run("BASE-Blob unlesbar (M)", func(t *testing.T) {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", vcsConfig)
		writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
		c1 := commitAll(t, wt, "c1")
		writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue B, zweite Fassung."))
		c2 := commitAll(t, wt, "c2")

		removeLooseObject(t, dir, treeEntryHash(t, dir, c1, "adr-x.md"))

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
		if code != 2 {
			t.Fatalf("Exit = %d, want 2\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
		}
	})

	t.Run("BASE-Tree unlesbar, mit Pendant (A)", func(t *testing.T) {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", vcsConfigSub)
		writeAt(t, dir, "sub/other.md", "x\n")
		c1 := commitAll(t, wt, "c1")
		writeAt(t, dir, "sub/adr-x.md", adrText("Accepted", "Tue A."))
		c2 := commitAll(t, wt, "c2")

		removeLooseObject(t, dir, treeEntryHash(t, dir, c1, "sub"))

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
		if code != 2 {
			t.Fatalf("Exit = %d, want 2\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
		}
	})

	t.Run("BASE-Tree unlesbar, ohne Pendant (Verzeichnis geloescht)", func(t *testing.T) {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", vcsConfigSub)
		writeAt(t, dir, "sub/adr-x.md", adrText("Accepted", "Tue A."))
		c1 := commitAll(t, wt, "c1")
		if _, err := wt.Remove("sub/adr-x.md"); err != nil {
			t.Fatal(err)
		}
		c2 := commitAll(t, wt, "c2 - verzeichnis geloescht")

		removeLooseObject(t, dir, treeEntryHash(t, dir, c1, "sub"))

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
		if code != 2 {
			t.Fatalf("Exit = %d, want 2 (vor slice-220: still 0 Befund(e)/Exit 0)\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
		}
	})

	t.Run("HEAD-Tree unlesbar (vierte Ausprägung, vorher Fehldiagnose)", func(t *testing.T) {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", vcsConfigSub)
		writeAt(t, dir, "sub/adr-x.md", adrText("Accepted", "Tue A."))
		writeAt(t, dir, "sub/other.md", "x\n")
		c1 := commitAll(t, wt, "c1")
		writeAt(t, dir, "sub/other.md", "geaendert\n")
		c2 := commitAll(t, wt, "c2")

		removeLooseObject(t, dir, treeEntryHash(t, dir, c2, "sub"))

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "vcs", "--range", c1 + ".." + c2, dir}, &stdout, &stderr)
		if code != 2 {
			t.Fatalf("Exit = %d, want 2 (vor slice-220: faelschlich core-drift-vcs, Exit 1)\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
		}
		if strings.Contains(stdout.String(), "core-drift-vcs") {
			t.Fatalf("Fehldiagnose als core-drift-vcs, statt als Umgebungsfehler zu brechen: %s", stdout.String())
		}
	})
}
