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
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// vcs.follow-renames end-to-end gegen ein echtes On-Disk-git-Repo
// (DC-FA-VCS-001, ADR-0115): eine Layout-Migration benennt eine verlinkte Datei
// um und zieht den Link einer Accepted-ADR nach.

// vcsRenamesConfig ist vcsConfig mit den Schlüsseln aus keys, je eine Zeile.
func vcsRenamesConfig(keys ...string) string {
	cfg := vcsConfig
	for _, k := range keys {
		cfg += "  " + k + "\n"
	}
	return cfg
}

// mvAt benennt eine Datei um: der neue Pfad entsteht mit demselben Inhalt, der
// alte wird aus Index und Arbeitsbaum entfernt.
func mvAt(t *testing.T, dir string, wt *gogit.Worktree, from, to string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, from))
	if err != nil {
		t.Fatal(err)
	}
	writeAt(t, dir, to, string(b))
	if _, err := wt.Remove(from); err != nil {
		t.Fatal(err)
	}
}

func adrLinking(target string) string {
	return adrText("Accepted", "Siehe [Spezifikation]("+target+").")
}

const specBody = "# Spezifikation\n\n## Teil\n\nInhalt.\n"

// runVCSRange fährt das Modul vcs über base..head und liefert Exit-Code und beide Ströme.
func runVCSRange(dir, base, head string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "vcs", "--range", base + ".." + head, dir}, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// TestVCS_FollowRenames deckt (DC-FA-VCS-001) Umbenennen mit Nachzug, Löschen und
// Neuanlegen, Umbiegen sowie geänderten Anker und Linktext ab, dazu
// den Nachzug in getrennten Commits. Je Fall: mit dem Schlüssel das erwartete
// Ergebnis, ohne ihn bleibt der Nachzug Drift und jeder andere Befund bestehen.
func TestVCS_FollowRenames(t *testing.T) {
	type repoCase struct {
		name string
		// build legt die Commits an und liefert BASE und HEAD.
		build           func(t *testing.T, dir string, wt *gogit.Worktree) (string, string)
		withKey, noKey  int
	}
	cases := []repoCase{
		{name: "umbenennen und nachziehen in einem commit", withKey: 0, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				c1 := commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
				return c1, commitAll(t, wt, "c2")
			}},
		{name: "reiner move und nachzug in getrennten commits", withKey: 0, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				c1 := commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				commitAll(t, wt, "c2 move")
				writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
				return c1, commitAll(t, wt, "c3 nachzug")
			}},
		{name: "gelöscht und mit anderem inhalt neu angelegt", withKey: 1, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				c1 := commitAll(t, wt, "c1")
				if _, err := wt.Remove("spec.md"); err != nil {
					t.Fatal(err)
				}
				writeAt(t, dir, "specification.md", "# Specification\n\n## Part\n\nAnderer Inhalt.\n")
				writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
				return c1, commitAll(t, wt, "c2")
			}},
		{name: "umbiegen auf eine bestehende, nicht umbenannte datei", withKey: 1, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "other.md", "# Anderes\n\nDokument.\n")
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				c1 := commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				writeAt(t, dir, "adr-x.md", adrLinking("other.md"))
				return c1, commitAll(t, wt, "c2")
			}},
		{name: "nachzug mit geändertem anker", withKey: 1, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md#teil"))
				c1 := commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				writeAt(t, dir, "adr-x.md", adrLinking("specification.md#part"))
				return c1, commitAll(t, wt, "c2")
			}},
		{name: "nachzug mit geändertem linktext", withKey: 1, noKey: 1,
			build: func(t *testing.T, dir string, wt *gogit.Worktree) (string, string) {
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				c1 := commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				writeAt(t, dir, "adr-x.md", strings.Replace(adrLinking("specification.md"), "[Spezifikation]", "[Specification]", 1))
				return c1, commitAll(t, wt, "c2")
			}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				label string
				cfg   string
				want  int
			}{
				{"mit follow-renames", vcsRenamesConfig("follow-renames: true"), c.withKey},
				{"mit beiden Schlüsseln", vcsRenamesConfig("ignore-link-targets: true", "follow-renames: true"), c.withKey},
				{"ohne den Schlüssel", vcsConfig, c.noKey},
			} {
				dir := t.TempDir()
				wt := initVCSRepo(t, dir)
				writeAt(t, dir, ".d-check.yml", mode.cfg)
				base, head := c.build(t, dir, wt)
				code, stdout, stderr := runVCSRange(dir, base, head)
				if code != mode.want {
					t.Fatalf("%s: Exit = %d, want %d\nstdout=%s\nstderr=%s", mode.label, code, mode.want, stdout, stderr)
				}
				if mode.want == 1 && !strings.Contains(stdout, "core-drift-vcs") {
					t.Fatalf("%s: stdout ohne core-drift-vcs: %s", mode.label, stdout)
				}
			}
		})
	}
}

// Ohne den Schlüssel — fehlend oder ausdrücklich aus — ist die Ausgabe (DC-FA-VCS-001,
// DC-QA-02) byte-identisch; der Lauf liest dann keinen Commit zwischen BASE und
// HEAD (siehe TestVCS_FollowRenamesFailClosed).
func TestVCS_FollowRenamesAusByteIdentisch(t *testing.T) {
	var outputs [2][3]string
	for i, cfg := range []string{vcsConfig, vcsRenamesConfig("follow-renames: false")} {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", cfg)
		writeAt(t, dir, "spec.md", specBody)
		writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
		c1 := commitAll(t, wt, "c1")
		mvAt(t, dir, wt, "spec.md", "specification.md")
		writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
		c2 := commitAll(t, wt, "c2")
		code, stdout, stderr := runVCSRange(dir, c1, c2)
		outputs[i] = [3]string{string(rune('0' + code)), stdout, stderr}
	}
	if outputs[0][0] != "1" || !strings.Contains(outputs[0][1], "core-drift-vcs") {
		t.Fatalf("ohne Schlüssel bleibt der Nachzug Drift: %v", outputs[0])
	}
	if outputs[0] != outputs[1] {
		t.Fatalf("Ausgabe ohne Schlüssel und mit follow-renames: false weicht ab:\n%v\n%v", outputs[0], outputs[1])
	}
}

// Ein Commit zwischen BASE und HEAD, der unlesbar ist (DC-FA-VCS-001), bricht den Lauf nur mit
// dem Schlüssel (Exit 2, Hinweis auf stderr); ohne ihn wird er nie gelesen
// (Exit 0). BASE und HEAD sind lesbar, und zwischen ihnen liegt keine Drift.
func TestVCS_FollowRenamesFailClosed(t *testing.T) {
	for _, c := range []struct {
		label string
		cfg   string
		want  int
	}{
		{"mit follow-renames", vcsRenamesConfig("follow-renames: true"), 2},
		{"ohne den Schlüssel", vcsConfig, 0},
	} {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", c.cfg)
		writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
		writeAt(t, dir, "spec.md", specBody)
		writeAt(t, dir, "n.txt", "eins\n")
		c1 := commitAll(t, wt, "c1")
		writeAt(t, dir, "n.txt", "zwei zwei\n")
		c2 := commitAll(t, wt, "c2")
		writeAt(t, dir, "n.txt", "drei drei drei\n")
		c3 := commitAll(t, wt, "c3")
		removeLooseObject(t, dir, plumbing.NewHash(c2))

		code, stdout, stderr := runVCSRange(dir, c1, c3)
		if code != c.want {
			t.Fatalf("%s: Exit = %d, want %d\nstdout=%s\nstderr=%s", c.label, code, c.want, stdout, stderr)
		}
		if c.want == 2 && !strings.Contains(stderr, "nicht lesbar") {
			t.Fatalf("%s: stderr nennt die Ursache nicht: %s", c.label, stderr)
		}
	}
}

// TestVCS_FollowRenamesZehnLaeufe: dieselbe Eingabe ergibt zehnmal dieselbe
// Ausgabe (DC-FA-VCS-001, DC-QA-02).
func TestVCS_FollowRenamesZehnLaeufe(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", vcsRenamesConfig("follow-renames: true"))
	writeAt(t, dir, "spec.md", specBody)
	writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
	c1 := commitAll(t, wt, "c1")
	mvAt(t, dir, wt, "spec.md", "specification.md")
	commitAll(t, wt, "c2 move")
	writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
	c3 := commitAll(t, wt, "c3 nachzug")
	code, stdout, stderr := runVCSRange(dir, c1, c3)
	for i := 0; i < 9; i++ {
		c, o, e := runVCSRange(dir, c1, c3)
		if c != code || o != stdout || e != stderr {
			t.Fatalf("Lauf %d weicht ab: %d/%q/%q gegen %d/%q/%q", i+2, c, o, e, code, stdout, stderr)
		}
	}
	if code != 0 {
		t.Fatalf("Exit = %d, want 0\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}
}

// TestVCS_FollowRenamesStaged: mit --staged ist der einzige Schritt HEAD →
// Index (DC-FA-VCS-001): ein gestagter Move samt Nachzug folgt, ein Move in
// einem früheren Commit mit nur gestagtem Nachzug bleibt Drift.
func TestVCS_FollowRenamesStaged(t *testing.T) {
	for _, c := range []struct {
		name       string
		moveStaged bool
		withKey    int
		noKey      int
	}{
		{"git mv und Nachzug gestagt", true, 0, 1},
		{"Move in einem früheren Commit, nur der Nachzug gestagt", false, 1, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				label string
				cfg   string
				want  int
			}{
				{"mit follow-renames", vcsRenamesConfig("follow-renames: true"), c.withKey},
				{"mit beiden Schlüsseln", vcsRenamesConfig("ignore-link-targets: true", "follow-renames: true"), c.withKey},
				{"ohne den Schlüssel", vcsConfig, c.noKey},
			} {
				dir := t.TempDir()
				wt := initVCSRepo(t, dir)
				writeAt(t, dir, ".d-check.yml", mode.cfg)
				writeAt(t, dir, "spec.md", specBody)
				writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
				commitAll(t, wt, "c1")
				mvAt(t, dir, wt, "spec.md", "specification.md")
				if c.moveStaged {
					if _, err := wt.Add("specification.md"); err != nil {
						t.Fatal(err)
					}
				} else {
					commitAll(t, wt, "c2 move")
				}
				writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
				if _, err := wt.Add("adr-x.md"); err != nil {
					t.Fatal(err)
				}
				var stdout, stderr bytes.Buffer
				code := cli.Run([]string{"--enable", "vcs", "--staged", dir}, &stdout, &stderr)
				if code != mode.want {
					t.Fatalf("%s: Exit = %d, want %d\nstdout=%s\nstderr=%s", mode.label, code, mode.want, stdout.String(), stderr.String())
				}
			}
		})
	}
}

// TestVCS_FollowRenamesBasisKeinVorfahr: ist BASE kein Vorfahr von HEAD, trägt
// die Range keine Umbenennung — der Lauf urteilt wie ohne den Schlüssel, ohne
// Exit 2 (DC-FA-VCS-001, ADR-0115 Entscheidung 12).
func TestVCS_FollowRenamesBasisKeinVorfahr(t *testing.T) {
	var outputs [2][3]string
	for i, cfg := range []string{vcsConfig, vcsRenamesConfig("follow-renames: true")} {
		dir := t.TempDir()
		wt := initVCSRepo(t, dir)
		writeAt(t, dir, ".d-check.yml", cfg)
		writeAt(t, dir, "spec.md", specBody)
		writeAt(t, dir, "adr-x.md", adrLinking("spec.md"))
		c0 := commitAll(t, wt, "c0")
		mvAt(t, dir, wt, "spec.md", "specification.md")
		writeAt(t, dir, "adr-x.md", adrLinking("specification.md"))
		head := commitAll(t, wt, "head")
		if err := wt.Checkout(&gogit.CheckoutOptions{Hash: plumbing.NewHash(c0), Force: true}); err != nil {
			t.Fatal(err)
		}
		writeAt(t, dir, "seitlich.txt", "ein anderer Zweig\n")
		base := commitAll(t, wt, "seitlich")
		code, stdout, stderr := runVCSRange(dir, base, head)
		outputs[i] = [3]string{string(rune('0' + code)), stdout, stderr}
	}
	if outputs[0][0] != "1" {
		t.Fatalf("ohne Schlüssel: Exit 1 erwartet, got %v", outputs[0])
	}
	if outputs[0] != outputs[1] {
		t.Fatalf("Basis kein Vorfahr: Ausgabe mit dem Schlüssel weicht ab:\n%v\n%v", outputs[0], outputs[1])
	}
}

// TestVCS_FollowRenamesImmutableDateiSelbst: die Umbenennung der immutablen
// Datei selbst bleibt Drift (DC-FA-VCS-001): der Schlüssel wirkt nur auf Ziele.
func TestVCS_FollowRenamesImmutableDateiSelbst(t *testing.T) {
	dir := t.TempDir()
	wt := initVCSRepo(t, dir)
	writeAt(t, dir, ".d-check.yml", vcsRenamesConfig("follow-renames: true"))
	writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
	c1 := commitAll(t, wt, "c1")
	mvAt(t, dir, wt, "adr-x.md", "adr-y.md")
	c2 := commitAll(t, wt, "c2")
	code, stdout, stderr := runVCSRange(dir, c1, c2)
	if code != 1 || !strings.Contains(stdout, "core-drift-vcs") {
		t.Fatalf("Exit = %d, want 1 mit core-drift-vcs\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}
}

// gitlinkCommit legt einen Commit an, dessen Baum dem von parent gleicht, nur
// dass path ein Gitlink ist.
func gitlinkCommit(t *testing.T, dir, parent, path string) string {
	t.Helper()
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := repo.CommitObject(plumbing.NewHash(parent))
	if err != nil {
		t.Fatal(err)
	}
	tree, err := pc.Tree()
	if err != nil {
		t.Fatal(err)
	}
	var entries []object.TreeEntry
	for _, e := range tree.Entries {
		if e.Name == path {
			e = object.TreeEntry{Name: path, Mode: filemode.Submodule, Hash: plumbing.NewHash("cafebabecafebabecafebabecafebabecafebabe")}
		}
		entries = append(entries, e)
	}
	obj := repo.Storer.NewEncodedObject()
	if err := (&object.Tree{Entries: entries}).Encode(obj); err != nil {
		t.Fatal(err)
	}
	th, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	sig := object.Signature{Name: "T", Email: "t@example.com", When: time.Unix(1700000001, 0)}
	cobj := repo.Storer.NewEncodedObject()
	c := &object.Commit{Author: sig, Committer: sig, Message: "gitlink", TreeHash: th, ParentHashes: []plumbing.Hash{plumbing.NewHash(parent)}}
	if err := c.Encode(cobj); err != nil {
		t.Fatal(err)
	}
	h, err := repo.Storer.SetEncodedObject(cobj)
	if err != nil {
		t.Fatal(err)
	}
	return h.String()
}

// TestVCS_ImmutableDateiDurchGitlinkErsetzt: ersetzt ein Gitlink eine immutable
// Datei an ihrem Pfad, meldet das Modul `core-drift-vcs` — mit --range wie mit
// --staged, mit und ohne vcs.follow-renames (DC-FA-VCS-001.a Schritt 2: der
// Index nimmt dieselben drei Arten wie der Baum).
func TestVCS_ImmutableDateiDurchGitlinkErsetzt(t *testing.T) {
	for _, withKey := range []bool{false, true} {
		cfg := vcsConfig
		if withKey {
			cfg = vcsRenamesConfig("follow-renames: true")
		}
		t.Run("--range", func(t *testing.T) {
			dir := t.TempDir()
			wt := initVCSRepo(t, dir)
			writeAt(t, dir, ".d-check.yml", cfg)
			writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
			c1 := commitAll(t, wt, "c1")
			c2 := gitlinkCommit(t, dir, c1, "adr-x.md")
			code, stdout, stderr := runVCSRange(dir, c1, c2)
			if code != 1 || !strings.Contains(stdout, "core-drift-vcs") {
				t.Fatalf("follow-renames=%v: Exit = %d, want 1 mit core-drift-vcs\nstdout=%s\nstderr=%s", withKey, code, stdout, stderr)
			}
		})
		t.Run("--staged", func(t *testing.T) {
			dir := t.TempDir()
			wt := initVCSRepo(t, dir)
			writeAt(t, dir, ".d-check.yml", cfg)
			writeAt(t, dir, "adr-x.md", adrText("Accepted", "Tue A."))
			commitAll(t, wt, "c1")
			repo, err := gogit.PlainOpen(dir)
			if err != nil {
				t.Fatal(err)
			}
			idx, err := repo.Storer.Index()
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range idx.Entries {
				if e.Name == "adr-x.md" {
					e.Mode = filemode.Submodule
					e.Hash = plumbing.NewHash("cafebabecafebabecafebabecafebabecafebabe")
				}
			}
			if err := repo.Storer.SetIndex(idx); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := cli.Run([]string{"--enable", "vcs", "--staged", dir}, &stdout, &stderr)
			if code != 1 || !strings.Contains(stdout.String(), "core-drift-vcs") {
				t.Fatalf("follow-renames=%v: Exit = %d, want 1 mit core-drift-vcs\nstdout=%s\nstderr=%s", withKey, code, stdout.String(), stderr.String())
			}
		})
	}
}
