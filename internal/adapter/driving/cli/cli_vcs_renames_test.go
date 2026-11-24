package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

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
