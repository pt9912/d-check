package git_test

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"

	gitadapter "github.com/pt9912/d-check/internal/adapter/driven/git"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Tests der Umbenennungs-Abbildung (driven.VCS.Renames, DC-FA-VCS-001.a Schritt 2,
// ADR-0115) gegen echte git-Objekte. Die Fixtures legen Blobs, Bäume und Commits
// direkt im Objekt-Speicher an: Form und Abstammung jedes Commits stehen in der
// Tabelle, kein Index und kein Arbeitsbaum redet dazwischen.

// tfile ist ein Eintrag eines Fixture-Baums: Inhalt und Art; ein Gitlink trägt
// stattdessen eine Kennung.
type tfile struct {
	content string
	mode    filemode.FileMode
	hash    plumbing.Hash
}

func reg(content string) tfile { return tfile{content: content, mode: filemode.Regular} }

// blobIn legt ein Blob mit dem Inhalt content an.
func blobIn(t *testing.T, repo *gogit.Repository, content string) plumbing.Hash {
	t.Helper()
	obj := repo.Storer.NewEncodedObject()
	obj.SetType(plumbing.BlobObject)
	w, err := obj.Writer()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(w, content); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	h, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// treeIn legt den Baum zu files samt Unterbäumen an und liefert seine Kennung.
func treeIn(t *testing.T, repo *gogit.Repository, files map[string]tfile) plumbing.Hash {
	t.Helper()
	dirs := map[string]map[string]tfile{}
	var entries []object.TreeEntry
	for p, f := range files {
		head, rest, nested := strings.Cut(p, "/")
		if nested {
			if dirs[head] == nil {
				dirs[head] = map[string]tfile{}
			}
			dirs[head][rest] = f
			continue
		}
		h := f.hash
		if f.mode != filemode.Submodule {
			h = blobIn(t, repo, f.content)
		}
		entries = append(entries, object.TreeEntry{Name: head, Mode: f.mode, Hash: h})
	}
	for name, sub := range dirs {
		entries = append(entries, object.TreeEntry{Name: name, Mode: filemode.Dir, Hash: treeIn(t, repo, sub)})
	}
	// git ordnet die Einträge nach dem Namen, ein Verzeichnis als Name samt "/".
	key := func(e object.TreeEntry) string {
		if e.Mode == filemode.Dir {
			return e.Name + "/"
		}
		return e.Name
	}
	sort.Slice(entries, func(i, j int) bool { return key(entries[i]) < key(entries[j]) })
	obj := repo.Storer.NewEncodedObject()
	if err := (&object.Tree{Entries: entries}).Encode(obj); err != nil {
		t.Fatal(err)
	}
	h, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// commitIn legt einen Commit mit dem Baum zu files und den Eltern parents an.
func commitIn(t *testing.T, repo *gogit.Repository, msg string, files map[string]tfile, parents ...plumbing.Hash) plumbing.Hash {
	t.Helper()
	sig := object.Signature{Name: "T", Email: "t@example.com", When: time.Unix(1700000000, 0)}
	c := &object.Commit{
		Author: sig, Committer: sig, Message: msg,
		TreeHash: treeIn(t, repo, files), ParentHashes: parents,
	}
	obj := repo.Storer.NewEncodedObject()
	if err := c.Encode(obj); err != nil {
		t.Fatal(err)
	}
	h, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// commitSpec ist ein Commit der Tabelle: Name, Namen der Eltern (leer: Wurzel),
// Baum.
type commitSpec struct {
	name    string
	parents []string
	files   map[string]tfile
}

// graph legt die Commits in Tabellen-Reihenfolge an (Eltern vor Kindern) und
// liefert Name → Kennung.
func graph(t *testing.T, repo *gogit.Repository, specs []commitSpec) map[string]plumbing.Hash {
	t.Helper()
	ids := map[string]plumbing.Hash{}
	for _, s := range specs {
		var parents []plumbing.Hash
		for _, p := range s.parents {
			h, ok := ids[p]
			if !ok {
				t.Fatalf("Fixture: Eltern-Commit %q von %q fehlt", p, s.name)
			}
			parents = append(parents, h)
		}
		ids[s.name] = commitIn(t, repo, s.name, s.files, parents...)
	}
	return ids
}

// newRepo legt ein leeres Repository an.
func newRepo(t *testing.T) (string, *gogit.Repository) {
	t.Helper()
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	return dir, repo
}

// removeObject entfernt das lose Objekt zu h — das Objekt wird unlesbar, ohne
// dass ein anderer Zustand sich ändert.
func removeObject(t *testing.T, dir string, h plumbing.Hash) {
	t.Helper()
	s := h.String()
	if err := os.Remove(filepath.Join(dir, ".git", "objects", s[:2], s[2:])); err != nil {
		t.Fatal(err)
	}
}

// m ist Kurzschrift für die erwartete Zuordnung Quelle → Ende.
type m = map[string]string

// Dateien, die in den Fällen wiederkehren: ein Inhalt je Name, damit gleiche
// Namen gleiche Blobs tragen.
func fA() tfile { return reg("alpha\n") }
func fB() tfile { return reg("beta\n") }
func fC() tfile { return reg("gamma\n") }

type renameCase struct {
	name       string
	commits    []commitSpec
	base, head string
	want       m
}

func runRenameCases(t *testing.T, cases []renameCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir, repo := newRepo(t)
			ids := graph(t, repo, tc.commits)
			ad, err := gitadapter.Open(dir)
			if err != nil {
				t.Fatal(err)
			}
			got, err := ad.Renames(ids[tc.base].String(), ids[tc.head].String())
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			want := tc.want
			if want == nil {
				want = m{}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Zuordnung = %v, want %v", got, want)
			}
		})
	}
}

// Umbenennungen innerhalb einer Linie (DC-FA-VCS-001): ein Commit, getrennte
// Commits, Kette; Umbenennen und im selben Commit ändern ist keine.
func TestRenamesLinear(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"ein Commit", []commitSpec{
			{"base", nil, map[string]tfile{"spec/a.md": fA(), "x.md": fC()}},
			{"head", []string{"base"}, map[string]tfile{"spec/b.md": fA(), "x.md": fC()}},
		}, "base", "head", m{"spec/a.md": "spec/b.md"}},
		{"Move und Inhaltsänderung in getrennten Commits", []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA()}},
			{"move", []string{"base"}, map[string]tfile{"b.md": fA()}},
			{"head", []string{"move"}, map[string]tfile{"b.md": reg("alpha neu\n")}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Inhaltsänderung vor dem Move", []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA(), "n.md": fC()}},
			{"edit", []string{"base"}, map[string]tfile{"a.md": reg("alpha neu\n"), "n.md": fC()}},
			{"head", []string{"edit"}, map[string]tfile{"b.md": reg("alpha neu\n"), "n.md": fC()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Kette über drei Commits", []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA()}},
			{"c1", []string{"base"}, map[string]tfile{"b.md": fA()}},
			{"c2", []string{"c1"}, map[string]tfile{"b.md": reg("beta neu\n")}},
			{"head", []string{"c2"}, map[string]tfile{"c.md": reg("beta neu\n")}},
		}, "base", "head", m{"a.md": "c.md"}},
		{"Umbenennen und im selben Commit ändern", []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA()}},
			{"head", []string{"base"}, map[string]tfile{"b.md": reg("alpha neu\n")}},
		}, "base", "head", nil},
		{"gelöscht und anderes angelegt", []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA(), "n.md": fC()}},
			{"head", []string{"base"}, map[string]tfile{"b.md": fB(), "n.md": fC()}},
		}, "base", "head", nil},
	})
}

// Der Fehlerpfad (DC-FA-VCS-001): ein Objekt, das nur die Umbenennungs-Erkennung liest und das
// unlesbar ist, bricht ab — auch wenn BASE und HEAD lesbar sind.
func TestRenamesUnlesbaresObjekt(t *testing.T) {
	// base → mid → head; mid legt d/f.md an, head nimmt es zurück: der Unterbaum d
	// von mid unterscheidet sich von dem an base und dem an head.
	build := func(t *testing.T) (string, map[string]plumbing.Hash, *gogit.Repository) {
		dir, repo := newRepo(t)
		ids := graph(t, repo, []commitSpec{
			{"base", nil, map[string]tfile{"a.md": fA(), "d/e.md": fC()}},
			{"mid", []string{"base"}, map[string]tfile{"b.md": fA(), "d/e.md": fC(), "d/f.md": fB()}},
			{"head", []string{"mid"}, map[string]tfile{"b.md": fA(), "d/e.md": fC()}},
		})
		return dir, ids, repo
	}
	treeOf := func(t *testing.T, repo *gogit.Repository, h plumbing.Hash) *object.Tree {
		t.Helper()
		cm, err := repo.CommitObject(h)
		if err != nil {
			t.Fatal(err)
		}
		tr, err := cm.Tree()
		if err != nil {
			t.Fatal(err)
		}
		return tr
	}
	cases := []struct {
		name string
		// victim wählt das Objekt, das unlesbar wird; es gehört weder zum Baum von
		// base noch zu dem von head.
		victim func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash
	}{
		{"Commit-Objekt eines Zwischen-Commits", func(_ *testing.T, _ *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
			return ids["mid"]
		}},
		{"Baum eines Schritts", func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
			cm, err := repo.CommitObject(ids["mid"])
			if err != nil {
				t.Fatal(err)
			}
			return cm.TreeHash
		}},
		{"Teilbaum, der sich im Schritt ändert", func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
			e, err := treeOf(t, repo, ids["mid"]).FindEntry("d")
			if err != nil {
				t.Fatal(err)
			}
			return e.Hash
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir, ids, repo := build(t)
			victim := tc.victim(t, repo, ids)
			for _, end := range []string{"base", "head"} {
				if cm, err := repo.CommitObject(ids[end]); err != nil || cm.TreeHash == victim {
					t.Fatalf("Fixture: %s darf das Opfer nicht tragen (err=%v)", end, err)
				}
			}
			removeObject(t, dir, victim)
			ad, err := gitadapter.Open(dir)
			if err != nil {
				t.Fatal(err)
			}
			// Vorbedingung: BASE und HEAD sind lesbar, die Enumeration besteht.
			if _, _, err := ad.AllPaths(ids["base"].String(), ids["head"].String()); err != nil {
				t.Fatalf("Vorbedingung verfehlt, AllPaths: %v", err)
			}
			if _, err := ad.Renames(ids["base"].String(), ids["head"].String()); err == nil {
				t.Fatal("unlesbares Objekt still übergangen — erwartet war ein Fehler")
			}
		})
	}
}

// Mit --staged ist der einzige Schritt HEAD → Index (DC-FA-VCS-001).
func TestRenamesStaged(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "spec/a.md", "alpha\n")
	put(t, dir, "x.md", "gamma\n")
	snapshot(t, wt, "first")
	if err := os.Rename(filepath.Join(dir, "spec/a.md"), filepath.Join(dir, "spec/b.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Remove("spec/a.md"); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Add("spec/b.md"); err != nil {
		t.Fatal(err)
	}
	ad, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ad.Renames("HEAD", driven.IndexRef)
	if err != nil {
		t.Fatal(err)
	}
	if want := (m{"spec/a.md": "spec/b.md"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Zuordnung = %v, want %v", got, want)
	}
}

// Nicht auflösbare Refs sind ein Fehler, eine leere Range keine Umbenennung
// (DC-FA-VCS-001).
func TestRenamesFehlerFormen(t *testing.T) {
	dir, repo := newRepo(t)
	ids := graph(t, repo, []commitSpec{
		{"base", nil, map[string]tfile{"a.md": fA()}},
		{"head", []string{"base"}, map[string]tfile{"b.md": fA()}},
	})
	ad, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	zero := strings.Repeat("0", 40)
	for _, c := range []struct{ base, head string }{
		{zero, ids["head"].String()},
		{ids["base"].String(), zero},
	} {
		if _, err := ad.Renames(c.base, c.head); err == nil {
			t.Errorf("Renames(%s, %s): nicht auflösbares Ref hätte fail-closed liefern müssen", c.base, c.head)
		}
	}
	if got, err := ad.Renames(ids["base"].String(), ids["base"].String()); err != nil || len(got) != 0 {
		t.Errorf("Renames(base, base) = %v, %v; erwartet leer ohne Fehler", got, err)
	}
}
