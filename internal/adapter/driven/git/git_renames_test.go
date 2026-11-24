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
	"github.com/go-git/go-git/v5/plumbing/format/index"
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
func exe(content string) tfile { return tfile{content: content, mode: filemode.Executable} }
func lnk(target string) tfile  { return tfile{content: target, mode: filemode.Symlink} }
func gitlinkTo(h string) tfile { return tfile{mode: filemode.Submodule, hash: plumbing.NewHash(h)} }

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
// liefert Name → Kennung; salt wird an jede Nachricht gehängt und ändert damit
// alle Kennungen, nicht die Form.
func graph(t *testing.T, repo *gogit.Repository, specs []commitSpec, salt string) map[string]plumbing.Hash {
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
		ids[s.name] = commitIn(t, repo, s.name+salt, s.files, parents...)
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
func fD() tfile { return reg("delta\n") }

type renameCase struct {
	name       string
	commits    []commitSpec
	base, head string
	want       m
}

// files ist Kurzschrift für den Baum eines Commits.
type files = map[string]tfile

// runRenameCases fährt jeden Fall mit drei verschiedenen Nachrichten-Salzen:
// dieselbe Topologie mit anderen Commit-Kennungen ergibt dieselbe Zuordnung
// (DC-QA-02), die Aufzählung nach Kennung darf nichts verschieben.
func runRenameCases(t *testing.T, cases []renameCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.want
			if want == nil {
				want = m{}
			}
			for _, salt := range []string{"", " zweite Kennungen", " dritte Kennungen"} {
				dir, repo := newRepo(t)
				ids := graph(t, repo, tc.commits, salt)
				ad, err := gitadapter.Open(dir)
				if err != nil {
					t.Fatal(err)
				}
				got, err := ad.Renames(ids[tc.base].String(), ids[tc.head].String())
				if err != nil {
					t.Fatalf("unerwarteter Fehler: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("Salz %q: Zuordnung = %v, want %v", salt, got, want)
				}
			}
		})
	}
}

// Umbenennungen innerhalb einer Linie (DC-FA-VCS-001): ein Commit, getrennte
// Commits, Kette; Umbenennen und im selben Commit ändern ist keine.
func TestRenamesLinear(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"ein Commit", []commitSpec{
			{"base", nil, files{"spec/a.md": fA(), "x.md": fC()}},
			{"head", []string{"base"}, files{"spec/b.md": fA(), "x.md": fC()}},
		}, "base", "head", m{"spec/a.md": "spec/b.md"}},
		{"Move und Inhaltsänderung in getrennten Commits", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"move", []string{"base"}, files{"b.md": fA()}},
			{"head", []string{"move"}, files{"b.md": reg("alpha neu\n")}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Inhaltsänderung vor dem Move", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"edit", []string{"base"}, files{"a.md": reg("alpha neu\n"), "n.md": fC()}},
			{"head", []string{"edit"}, files{"b.md": reg("alpha neu\n"), "n.md": fC()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Kette über drei Commits", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"c1", []string{"base"}, files{"b.md": fA()}},
			{"c2", []string{"c1"}, files{"b.md": reg("beta neu\n")}},
			{"head", []string{"c2"}, files{"c.md": reg("beta neu\n")}},
		}, "base", "head", m{"a.md": "c.md"}},
		{"Umbenennen und im selben Commit ändern", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"head", []string{"base"}, files{"b.md": reg("alpha neu\n")}},
		}, "base", "head", nil},
		{"gelöscht und anderes angelegt", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"head", []string{"base"}, files{"b.md": fB(), "n.md": fC()}},
		}, "base", "head", nil},
		{"Ende in einem Commit gelöscht und angelegt: Inhaltsänderung an Ort und Stelle", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"c1", []string{"base"}, files{"b.md": fA()}},
			{"head", []string{"c1"}, files{"b.md": fB()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Kette: Ende in HEAD gelöscht", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"c1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"head", []string{"c1"}, files{"n.md": fC()}},
		}, "base", "head", nil},
		{"Ende gelöscht und in einem späteren Commit als andere Datei neu angelegt", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"c1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"c2", []string{"c1"}, files{"n.md": fC()}},
			{"head", []string{"c2"}, files{"b.md": fB(), "n.md": fC()}},
		}, "base", "head", nil},
		{"Zwischenpfad einer Kette neu belegt", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"c1", []string{"base"}, files{"b.md": fA()}},
			{"c2", []string{"c1"}, files{"d.md": fA()}},
			{"head", []string{"c2"}, files{"d.md": fA(), "b.md": fB()}},
		}, "base", "head", m{"a.md": "d.md"}},
	})
}

// Mehrdeutigkeit bleibt Drift (DC-FA-VCS-001): ein gleicher Eintrag auf mehr als
// einem verschwundenen oder neuen Pfad, zwei Enden, ein Ende in BASE, zwei
// Quellen auf ein Ende.
func TestRenamesMehrdeutigkeit(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"zwei verschwundene Pfade mit gleichem Blob", []commitSpec{
			{"base", nil, files{"a.md": fA(), "a2.md": fA()}},
			{"head", []string{"base"}, files{"b.md": fA()}},
		}, "base", "head", nil},
		{"zwei neue Pfade mit gleichem Blob", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"head", []string{"base"}, files{"b.md": fA(), "b2.md": fA()}},
		}, "base", "head", nil},
		{"Kopie: die Quelle besteht fort", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"head", []string{"base"}, files{"a.md": fA(), "b.md": fA()}},
		}, "base", "head", nil},
		{"zwei Linien führen zu zwei Enden", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"s2", []string{"base"}, files{"c.md": fA(), "n.md": fC()}},
			{"head", []string{"s1", "s2"}, files{"b.md": fA(), "c.md": fA(), "n.md": fC()}},
		}, "base", "head", nil},
		{"zwei Quellen auf ein Ende (Merge)", []commitSpec{
			{"base", nil, files{"a.md": fA(), "b.md": fB(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"e.md": fA(), "b.md": fB(), "n.md": fC()}},
			{"s2", []string{"base"}, files{"a.md": fA(), "e.md": fB(), "n.md": fC()}},
			{"head", []string{"s1", "s2"}, files{"e.md": fA(), "n.md": fC()}},
		}, "base", "head", nil},
		{"beide Seiten benennen gleich um", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"s2", []string{"base"}, files{"b.md": fA(), "n.md": fD()}},
			{"head", []string{"s1", "s2"}, files{"b.md": fA(), "n.md": fD()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Tausch über einen Zwischennamen", []commitSpec{
			{"base", nil, files{"a.md": fA(), "b.md": fB()}},
			{"c1", []string{"base"}, files{"t.md": fA(), "b.md": fB()}},
			{"c2", []string{"c1"}, files{"t.md": fA(), "a.md": fB()}},
			{"head", []string{"c2"}, files{"b.md": fA(), "a.md": fB()}},
		}, "base", "head", nil},
		{"Tausch der Inhalte in einem Commit", []commitSpec{
			{"base", nil, files{"a.md": fA(), "b.md": fB()}},
			{"head", []string{"base"}, files{"a.md": fB(), "b.md": fA()}},
		}, "base", "head", nil},
		{"Name freigemacht und neu belegt", []commitSpec{
			{"base", nil, files{"a.md": fA(), "c.md": fB()}},
			{"c1", []string{"base"}, files{"c.md": fB()}},
			{"head", []string{"c1"}, files{"a.md": fB()}},
		}, "base", "head", nil},
		{"Rotation: A nimmt als Quelle teil, C nicht", []commitSpec{
			{"base", nil, files{"a.md": fA(), "c.md": fB()}},
			{"c1", []string{"base"}, files{"x.md": fA(), "c.md": fB()}},
			{"head", []string{"c1"}, files{"x.md": fA(), "a.md": fB()}},
		}, "base", "head", m{"a.md": "x.md"}},
	})
}

// Merge-Seiten und Abstammung (DC-FA-VCS-001): eine Kette setzt sich nur entlang
// der Abstammung fort, ein Pfad, der auf einer Linie verschwindet, hat dort kein
// Ende, eine Seite, die ihn bestehen lässt, trägt nichts bei.
func TestRenamesMerge(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"Merge: eine Seite benennt um, die andere ändert nur den Inhalt", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"s2", []string{"base"}, files{"a.md": reg("alpha neu\n"), "n.md": fC()}},
			{"head", []string{"s1", "s2"}, files{"b.md": reg("alpha neu\n"), "n.md": fC()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Merge: die andere Seite löscht den Pfad und legt ihn später neu an", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"s2a", []string{"base"}, files{"n.md": fC()}},
			{"s2b", []string{"s2a"}, files{"a.md": fB(), "n.md": fC()}},
			{"head", []string{"s1", "s2b"}, files{"b.md": fA(), "a.md": fB(), "n.md": fC()}},
		}, "base", "head", nil},
		{"Merge: X → Y und Y → Z auf zwei Seiten setzen sich nicht fort", []commitSpec{
			{"base", nil, files{"x.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"y.md": fA(), "n.md": fC()}},
			{"s2a", []string{"base"}, files{"x.md": fA(), "y.md": fB(), "n.md": fC()}},
			{"s2b", []string{"s2a"}, files{"x.md": fA(), "z.md": fB(), "n.md": fC()}},
			{"head", []string{"s1", "s2b"}, files{"y.md": fA(), "z.md": fB(), "n.md": fC()}},
		}, "base", "head", m{"x.md": "y.md"}},
		{"Merge: die Auflösung entfernt das Ende", []commitSpec{
			{"base", nil, files{"x.md": fA(), "n.md": fC()}},
			{"s1", []string{"base"}, files{"y.md": fA(), "n.md": fC()}},
			{"s2", []string{"base"}, files{"x.md": fA(), "n.md": fD()}},
			{"head", []string{"s1", "s2"}, files{"n.md": fD()}},
		}, "base", "head", nil},
		{"Merge: eine Seite zweigt vor BASE ab", []commitSpec{
			{"p0", nil, files{"a.md": fA(), "k.md": fC()}},
			{"base", []string{"p0"}, files{"a.md": fA(), "k.md": fC(), "n.md": fD()}},
			{"c1", []string{"base"}, files{"b.md": fA(), "k.md": fC(), "n.md": fD()}},
			{"s1", []string{"p0"}, files{"a.md": fA(), "k.md": fB()}},
			{"head", []string{"c1", "s1"}, files{"b.md": fA(), "k.md": fB(), "n.md": fD()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Wurzel-Commit in der Range trägt keine Umbenennung bei", []commitSpec{
			{"base", nil, files{"a.md": fA(), "n.md": fC()}},
			{"c1", []string{"base"}, files{"b.md": fA(), "n.md": fC()}},
			{"r", nil, files{"r.md": fD()}},
			{"r2", []string{"r"}, files{"s.md": fD()}},
			{"head", []string{"c1", "r2"}, files{"b.md": fA(), "n.md": fC(), "s.md": fD()}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Basis kein Vorfahr von HEAD", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"o", nil, files{"a.md": fA()}},
			{"head", []string{"o"}, files{"b.md": fA()}},
		}, "base", "head", nil},
	})
}

// Arten der Einträge (DC-FA-VCS-001): das Ausführbar-Bit zählt nicht, ein
// Symlink gleicht keiner Datei, ein Verzeichnis und ein Gitlink nehmen nicht
// teil; die Dateien eines umbenannten Verzeichnisses folgen einzeln.
func TestRenamesArtDesEintrags(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"nur das Ausführbar-Bit ändert sich", []commitSpec{
			{"base", nil, files{"a.sh": reg("x\n")}},
			{"head", []string{"base"}, files{"b.sh": exe("x\n")}},
		}, "base", "head", m{"a.sh": "b.sh"}},
		{"Symlink gegen Datei mit gleichem Blob", []commitSpec{
			{"base", nil, files{"a.md": reg("x")}},
			{"head", []string{"base"}, files{"b.md": lnk("x")}},
		}, "base", "head", nil},
		{"Symlink wird zum Symlink", []commitSpec{
			{"base", nil, files{"a.md": lnk("x")}},
			{"head", []string{"base"}, files{"b.md": lnk("x")}},
		}, "base", "head", m{"a.md": "b.md"}},
		{"Verzeichnis umbenannt: die Dateien folgen einzeln", []commitSpec{
			{"base", nil, files{"d1/a.md": fA(), "d1/b.md": fB()}},
			{"head", []string{"base"}, files{"d2/a.md": fA(), "d2/b.md": fB()}},
		}, "base", "head", m{"d1/a.md": "d2/a.md", "d1/b.md": "d2/b.md"}},
		{"Gitlink ersetzt die Datei", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"head", []string{"base"}, files{"a.md": gitlinkTo("cafebabecafebabecafebabecafebabecafebabe")}},
		}, "base", "head", nil},
		{"ein Gitlink ist kein Ende", []commitSpec{
			{"base", nil, files{"a.md": fA()}},
			{"head", []string{"base"}, files{"b.md": gitlinkTo("cafebabecafebabecafebabecafebabecafebabe")}},
		}, "base", "head", nil},
		{"Datei wird zum Verzeichnis gleichen Namens", []commitSpec{
			{"base", nil, files{"a": fA()}},
			{"head", []string{"base"}, files{"a/b.md": fA()}},
		}, "base", "head", m{"a": "a/b.md"}},
	})
}

// unreadableCase beschreibt ein Objekt, das nur die Umbenennungs-Erkennung liest.
type unreadableCase struct {
	name       string
	commits    []commitSpec
	base, head string
	// victim wählt das Objekt, das unlesbar wird; es gehört weder zum Baum von
	// BASE noch zu dem von HEAD (die Vorbedingung der Probe prüft das über
	// AllPaths).
	victim func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash
}

func commitTreeHash(t *testing.T, repo *gogit.Repository, h plumbing.Hash) plumbing.Hash {
	t.Helper()
	cm, err := repo.CommitObject(h)
	if err != nil {
		t.Fatal(err)
	}
	return cm.TreeHash
}

func subtreeHash(t *testing.T, repo *gogit.Repository, h plumbing.Hash, path string) plumbing.Hash {
	t.Helper()
	cm, err := repo.CommitObject(h)
	if err != nil {
		t.Fatal(err)
	}
	tr, err := cm.Tree()
	if err != nil {
		t.Fatal(err)
	}
	e, err := tr.FindEntry(path)
	if err != nil {
		t.Fatal(err)
	}
	return e.Hash
}

// Der Fehlerpfad (DC-FA-VCS-001): ein Objekt, das nur die Umbenennungs-Erkennung
// liest und das unlesbar ist, bricht ab — auch wenn BASE und HEAD lesbar sind.
func TestRenamesUnlesbaresObjekt(t *testing.T) {
	// base → mid → head; mid legt d/f.md an, head nimmt es zurück: der Unterbaum d
	// von mid unterscheidet sich von dem an base und dem an head.
	linear := []commitSpec{
		{"base", nil, files{"a.md": fA(), "d/e.md": fC()}},
		{"mid", []string{"base"}, files{"b.md": fA(), "d/e.md": fC(), "d/f.md": fB()}},
		{"head", []string{"mid"}, files{"b.md": fA(), "d/e.md": fC()}},
	}
	// Eine Merge-Seite zweigt vor BASE ab: p0 ist weder BASE noch HEAD.
	fork := []commitSpec{
		{"p0", nil, files{"a.md": fA(), "k.md": fC()}},
		{"base", []string{"p0"}, files{"a.md": fA(), "k.md": fC(), "n.md": fD()}},
		{"c1", []string{"base"}, files{"b.md": fA(), "k.md": fC(), "n.md": fD()}},
		{"s1", []string{"p0"}, files{"a.md": fA(), "k.md": fB()}},
		{"head", []string{"c1", "s1"}, files{"b.md": fA(), "k.md": fB(), "n.md": fD()}},
	}
	for _, tc := range []unreadableCase{
		{"Commit-Objekt eines Zwischen-Commits", linear, "base", "head",
			func(_ *testing.T, _ *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash { return ids["mid"] }},
		{"Baum eines Schritts", linear, "base", "head",
			func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
				return commitTreeHash(t, repo, ids["mid"])
			}},
		{"Teilbaum, der sich im Schritt ändert", linear, "base", "head",
			func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
				return subtreeHash(t, repo, ids["mid"], "d")
			}},
		{"Baum des Eltern-Commits vor BASE", fork, "base", "head",
			func(t *testing.T, repo *gogit.Repository, ids map[string]plumbing.Hash) plumbing.Hash {
				return commitTreeHash(t, repo, ids["p0"])
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, repo := newRepo(t)
			ids := graph(t, repo, tc.commits, "")
			removeObject(t, dir, tc.victim(t, repo, ids))
			ad, err := gitadapter.Open(dir)
			if err != nil {
				t.Fatal(err)
			}
			// Vorbedingung: BASE und HEAD sind lesbar, die Enumeration besteht.
			if _, _, err := ad.AllPaths(ids[tc.base].String(), ids[tc.head].String()); err != nil {
				t.Fatalf("Vorbedingung verfehlt, AllPaths: %v", err)
			}
			if _, err := ad.Renames(ids[tc.base].String(), ids[tc.head].String()); err == nil {
				t.Fatal("unlesbares Objekt still übergangen — erwartet war ein Fehler")
			}
		})
	}
}

// stagedRepo legt ein Repository mit einem Commit aus first an und liefert den
// Adapter, den Arbeitsbaum und das Verzeichnis.
func stagedRepo(t *testing.T, first map[string]string) (*gitadapter.Adapter, *gogit.Worktree, string) {
	t.Helper()
	dir, wt := repoAt(t)
	for p, c := range first {
		put(t, dir, p, c)
	}
	snapshot(t, wt, "first")
	ad, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	return ad, wt, dir
}

// mvStaged benennt eine Datei um und legt beides in den Index.
func mvStaged(t *testing.T, dir string, wt *gogit.Worktree, from, to string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, to)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(dir, from), filepath.Join(dir, to)); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Remove(from); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Add(to); err != nil {
		t.Fatal(err)
	}
}

// Mit --staged ist der einzige Schritt HEAD → Index (DC-FA-VCS-001).
func TestRenamesStaged(t *testing.T) {
	t.Run("git mv im Index", func(t *testing.T) {
		ad, wt, dir := stagedRepo(t, map[string]string{"spec/a.md": "alpha\n", "x.md": "gamma\n"})
		mvStaged(t, dir, wt, "spec/a.md", "spec/b.md")
		got, err := ad.Renames("HEAD", driven.IndexRef)
		if err != nil {
			t.Fatal(err)
		}
		if want := (m{"spec/a.md": "spec/b.md"}); !reflect.DeepEqual(got, want) {
			t.Fatalf("Zuordnung = %v, want %v", got, want)
		}
	})
	t.Run("der Move liegt in einem früheren Commit", func(t *testing.T) {
		ad, wt, dir := stagedRepo(t, map[string]string{"a.md": "alpha\n"})
		mvStaged(t, dir, wt, "a.md", "b.md")
		snapshot(t, wt, "move")
		put(t, dir, "adr.md", "Nachzug\n")
		if _, err := wt.Add("adr.md"); err != nil {
			t.Fatal(err)
		}
		got, err := ad.Renames("HEAD", driven.IndexRef)
		if err != nil || len(got) != 0 {
			t.Fatalf("Zuordnung = %v, %v; erwartet leer ohne Fehler (benannte Grenze)", got, err)
		}
	})
	t.Run("zwei neue Pfade mit gleichem Blob sind mehrdeutig", func(t *testing.T) {
		ad, wt, dir := stagedRepo(t, map[string]string{"a.md": "alpha\n"})
		mvStaged(t, dir, wt, "a.md", "b.md")
		put(t, dir, "b2.md", "alpha\n")
		if _, err := wt.Add("b2.md"); err != nil {
			t.Fatal(err)
		}
		got, err := ad.Renames("HEAD", driven.IndexRef)
		if err != nil || len(got) != 0 {
			t.Fatalf("Zuordnung = %v, %v; erwartet leer ohne Fehler", got, err)
		}
	})
	t.Run("ohne HEAD", func(t *testing.T) {
		dir, wt := repoAt(t)
		put(t, dir, "a.md", "alpha\n")
		if _, err := wt.Add("a.md"); err != nil {
			t.Fatal(err)
		}
		ad, err := gitadapter.Open(dir)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ad.Renames("HEAD", driven.IndexRef)
		if err != nil || got == nil || len(got) != 0 {
			t.Fatalf("Zuordnung = %v, %v; erwartet leer ohne Fehler", got, err)
		}
	})
	t.Run("Basis nicht auflösbar", func(t *testing.T) {
		ad, _, _ := stagedRepo(t, map[string]string{"a.md": "alpha\n"})
		if _, err := ad.Renames(strings.Repeat("0", 40), driven.IndexRef); err == nil {
			t.Fatal("nicht auflösbare staged-Basis hätte fail-closed liefern müssen")
		}
	})
	t.Run("ein Gitlink im Index nimmt nicht teil", func(t *testing.T) {
		ad, _, dir := stagedRepo(t, map[string]string{"a.md": "alpha\n"})
		repo, err := gogit.PlainOpen(dir)
		if err != nil {
			t.Fatal(err)
		}
		idx, err := repo.Storer.Index()
		if err != nil {
			t.Fatal(err)
		}
		idx.Entries = append(idx.Entries[:0], &index.Entry{
			Name: "mod", Mode: filemode.Submodule, Hash: plumbing.NewHash("cafebabecafebabecafebabecafebabecafebabe"),
		})
		if err := repo.Storer.SetIndex(idx); err != nil {
			t.Fatal(err)
		}
		got, err := ad.Renames("HEAD", driven.IndexRef)
		if err != nil || len(got) != 0 {
			t.Fatalf("Zuordnung = %v, %v; erwartet leer ohne Fehler", got, err)
		}
	})
}

// Die Index-Enumeration nimmt dieselben drei Arten wie der Baum (DC-FA-VCS-001.a
// Schritt 2): ein Gitlink im Index ist kein Datei-Pfad.
func TestAllPathsStagedGitlinkIstKeinePfad(t *testing.T) {
	ad, _, dir := stagedRepo(t, map[string]string{"keep.md": "x\n"})
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := repo.Storer.Index()
	if err != nil {
		t.Fatal(err)
	}
	idx.Entries = append(idx.Entries, &index.Entry{
		Name: "mod", Mode: filemode.Submodule, Hash: plumbing.NewHash("cafebabecafebabecafebabecafebabecafebabe"),
	})
	if err := repo.Storer.SetIndex(idx); err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := ad.AllPaths("HEAD", driven.IndexRef)
	if err != nil {
		t.Fatal(err)
	}
	if contains(headAll, "mod") || !contains(headAll, "keep.md") || len(headAll) != 1 {
		t.Fatalf("Index-Menge = %v, erwartet nur keep.md (der Gitlink ist kein Datei-Pfad)", headAll)
	}
	if len(baseAll) != 1 || baseAll[0] != "keep.md" {
		t.Fatalf("BASE-Menge = %v", baseAll)
	}
}

// Nicht auflösbare Refs sind ein Fehler, eine leere Range keine Umbenennung
// (DC-FA-VCS-001).
func TestRenamesFehlerFormen(t *testing.T) {
	dir, repo := newRepo(t)
	ids := graph(t, repo, []commitSpec{
		{"base", nil, files{"a.md": fA()}},
		{"head", []string{"base"}, files{"b.md": fA()}},
	}, "")
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
	if got, err := ad.Renames(ids["head"].String(), ids["base"].String()); err != nil || len(got) != 0 {
		t.Errorf("Renames(head, base) = %v, %v; erwartet leer ohne Fehler (Basis kein Vorfahr)", got, err)
	}
}

// Die Commit-Menge ist „von HEAD erreichbar, von BASE nicht“ (DC-FA-VCS-001,
// ADR-0115 Entscheidung 2): ein Commit vor BASE trägt nichts bei. a.md wird vor
// BASE gelöscht und vor BASE mit anderem Inhalt neu angelegt; läse die Range
// diese Commits mit, trüge a.md die Linie „gelöscht“, und die Umbenennung in der
// Range ginge verloren.
func TestRenamesCommitMengeOhneVorBase(t *testing.T) {
	runRenameCases(t, []renameCase{
		{"Pfad vor BASE gelöscht und neu angelegt", []commitSpec{
			{"r", nil, files{"a.md": fA(), "n.md": fC()}},
			{"r2", []string{"r"}, files{"n.md": fC()}},
			{"base", []string{"r2"}, files{"a.md": fB(), "n.md": fC()}},
			{"head", []string{"base"}, files{"b.md": fB(), "n.md": fC()}},
		}, "base", "head", m{"a.md": "b.md"}},
	})
}
