package git_test

import (
	"os"
	"path/filepath"
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
	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

// Black-Box-Adapter-Tests gegen ein echtes On-Disk-git-Repo über die exportierte
// Open()-Tür (DC-FA-VCS-001) — granularer als der CLI-Integrationstest.

func repoAt(t *testing.T) (string, *gogit.Worktree) {
	t.Helper()
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	return dir, wt
}

// vorLaenge liefert die Länge des vorhandenen Inhalts oder -1, wenn es keinen
// gibt — die E/A-Hälfte des Fixture-Wächters, dessen Entscheidung im Kern liegt.
func vorLaenge(full string) int {
	b, err := os.ReadFile(full)
	if err != nil {
		return -1
	}
	return len(b)
}

func put(t *testing.T, dir, name, content string) {
	t.Helper()
	full := filepath.Join(dir, name)
	if err := coretest.GitFixtureRewriteHazard(name, vorLaenge(full), len(content)); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func snapshot(t *testing.T, wt *gogit.Worktree, msg string) string {
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

func contains(all []string, path string) bool {
	for _, p := range all {
		if p == path {
			return true
		}
	}
	return false
}

// TestAllPathsAndFileAt: AllPaths an zwei Ständen (Added/Modified/Deleted
// zeigen sich als Mengen-Differenz) + FileAt an Refs.
func TestAllPathsAndFileAt(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "keep.md", "v1\n")
	put(t, dir, "sub/gone.md", "alt\n")
	first := snapshot(t, wt, "first")

	put(t, dir, "keep.md", "fassung zwei\n")
	if _, err := wt.Remove("sub/gone.md"); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "fresh.md", "neu\n")
	second := snapshot(t, wt, "second")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := a.AllPaths(first, second)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(baseAll, "keep.md") || !contains(baseAll, "sub/gone.md") || contains(baseAll, "fresh.md") {
		t.Fatalf("BASE-Menge falsch: %v", baseAll)
	}
	if !contains(headAll, "keep.md") || contains(headAll, "sub/gone.md") || !contains(headAll, "fresh.md") {
		t.Fatalf("HEAD-Menge falsch: %v", headAll)
	}

	if b, ok, err := a.FileAt(first, "keep.md"); err != nil || !ok || string(b) != "v1\n" {
		t.Fatalf("FileAt(first,keep) = %q,%v,%v", b, ok, err)
	}
	if _, ok, err := a.FileAt(second, "sub/gone.md"); err != nil || ok {
		t.Fatalf("FileAt(second, gelöscht) erwartet ok=false, got ok=%v err=%v", ok, err)
	}
	if _, _, err := a.AllPaths("0000000000000000000000000000000000000000", second); err == nil {
		t.Fatal("unauflösbare Basis hätte fail-closed liefern müssen")
	}
	if _, _, err := a.FileAt("0000000000000000000000000000000000000000", "keep.md"); err == nil {
		t.Fatal("unauflösbares Ref hätte fail-closed liefern müssen")
	}
}

// TestStaged: staged-Menge (base-Tree gegen den Index) + FileAt(Index) + ohne
// HEAD leer.
func TestStaged(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "adr.md", "Accepted A\n")
	snapshot(t, wt, "first")
	put(t, dir, "adr.md", "Accepted B, zweite Fassung\n")
	if _, err := wt.Add("adr.md"); err != nil {
		t.Fatal(err)
	}

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := a.AllPaths("HEAD", driven.IndexRef)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(baseAll, "adr.md") || !contains(headAll, "adr.md") {
		t.Fatalf("staged: adr.md an beiden Enden erwartet, base=%v head=%v", baseAll, headAll)
	}
	if b, ok, err := a.FileAt(driven.IndexRef, "adr.md"); err != nil || !ok || string(b) != "Accepted B, zweite Fassung\n" {
		t.Fatalf("FileAt(Index) = %q,%v,%v", b, ok, err)
	}
	if _, ok, _ := a.FileAt(driven.IndexRef, "weg.md"); ok {
		t.Fatal("FileAt(Index, fehlend) erwartet ok=false")
	}
}

func TestStagedNoHead(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "adr.md", "x\n")
	if _, err := wt.Add("adr.md"); err != nil {
		t.Fatal(err)
	}
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := a.AllPaths("HEAD", driven.IndexRef)
	if err != nil || baseAll != nil || headAll != nil {
		t.Fatalf("ohne HEAD: leer/kein Fehler erwartet, got %v / %v / %v", baseAll, headAll, err)
	}
}

func TestOpenMissing(t *testing.T) {
	if _, err := gitadapter.Open(t.TempDir()); err == nil {
		t.Fatal("Open ohne .git hätte fail-closed liefern müssen")
	}
}

// TestCommitMessages: die Range base..head liefert die Nicht-Merge-Messages
// (ohne die Basis selbst) — die git-Eingabe des Moduls commits (DC-FA-COMMITS-001).
func TestCommitMessages(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "f.md", "1\n")
	base := snapshot(t, wt, "A: feat ADR-0001")
	put(t, dir, "f.md", "zwei\n")
	snapshot(t, wt, "B: chore ohne id")
	put(t, dir, "f.md", "drei drei\n")
	head := snapshot(t, wt, "C: docs slice-056")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	metas, err := a.CommitMessages(base, head)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, m := range metas {
		got[strings.TrimSpace(m.Message)] = true
		if len(m.ShortSHA) != 7 {
			t.Errorf("ShortSHA %q nicht 7-stellig", m.ShortSHA)
		}
	}
	if !got["B: chore ohne id"] || !got["C: docs slice-056"] {
		t.Fatalf("Range A..C sollte B und C enthalten: %v", got)
	}
	if got["A: feat ADR-0001"] {
		t.Fatal("Basis A darf nicht in der Range erscheinen")
	}
}

// TestCommitMessagesSkipsMerges: ein Merge-Commit (2 Parents) wird übersprungen
// (git rev-list --no-merges-Parität).
func TestCommitMessagesSkipsMerges(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "f.md", "1\n")
	base := snapshot(t, wt, "A: ADR-0001")
	put(t, dir, "f.md", "zwei\n")
	b := snapshot(t, wt, "B: slice-056")
	put(t, dir, "f.md", "drei drei\n")
	if err := wt.AddWithOptions(&gogit.AddOptions{All: true}); err != nil {
		t.Fatal(err)
	}
	mh, err := wt.Commit("Merge: kein bezug", &gogit.CommitOptions{
		Author:  &object.Signature{Name: "T", Email: "t@example.com", When: time.Unix(1700000000, 0)},
		Parents: []plumbing.Hash{plumbing.NewHash(b), plumbing.NewHash(base)},
	})
	if err != nil {
		t.Fatal(err)
	}
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	metas, err := a.CommitMessages(base, mh.String())
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range metas {
		if strings.HasPrefix(m.Message, "Merge:") {
			t.Fatalf("Merge-Commit darf nicht erscheinen: %v", metas)
		}
	}
}

// TestCommitMessagesFailClosed: staged (IndexRef) und eine unauflösbare Basis
// brechen laut ab (fail-closed, Exit 2).
func TestCommitMessagesFailClosed(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "f.md", "1\n")
	base := snapshot(t, wt, "A: ADR-0001")
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.CommitMessages(base, driven.IndexRef); err == nil {
		t.Fatal("IndexRef (staged) muss fail-closed liefern")
	}
	if _, err := a.CommitMessages("0000000000000000000000000000000000000000", base); err == nil {
		t.Fatal("unauflösbare Basis muss fail-closed liefern")
	}
}

// Ein REINER Rename (Inhalt byte-identisch) muss auch über AllPaths die
// ALTE Menge tragen (Delete-Hälfte) und den alten Pfad aus der NEUEN Menge
// verlieren: DC-FA-VCS-001 sagt den Befund für die umbenannte immutable Datei
// zu, ohne einen Eingabe-Modus einzuschränken oder Inhalts-Ähnlichkeit zu
// messen — der Kern (CheckVCS) bildet die Mengen-Differenz selbst, ohne dass
// der Adapter eine Rename-Erkennung bräuchte.
func TestAllPathsPureRenameYieldsDelete(t *testing.T) {
	dir, wt := repoAt(t)
	const body = "# ADR-0001 — Kern\n\n**Status:** Accepted\n\nEin Text, der unverändert bleibt.\n"
	put(t, dir, "adr/0001-kern.md", body)
	first := snapshot(t, wt, "first")

	put(t, dir, "adr/0002-kern.md", body) // byte-identisch — der Rename ist rein
	if _, err := wt.Remove("adr/0001-kern.md"); err != nil {
		t.Fatal(err)
	}
	second := snapshot(t, wt, "reiner Rename")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := a.AllPaths(first, second)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(baseAll, "adr/0001-kern.md") {
		t.Fatalf("reiner Rename: alter Pfad fehlt in BASE-Menge: %v", baseAll)
	}
	if contains(headAll, "adr/0001-kern.md") {
		t.Fatalf("reiner Rename: alter Pfad taucht noch in HEAD-Menge auf: %v", headAll)
	}
	if !contains(headAll, "adr/0002-kern.md") {
		t.Fatalf("reiner Rename: neuer Pfad fehlt in HEAD-Menge: %v", headAll)
	}
}

// TestAllPathsUnlesbarerUnterbaum: den Kern der Regression von slice-220 —
// ein Unterbaum, dessen Objekt fehlt, macht die gemeinsame git-Wartungsaufgabe
// `git maintenance run --task=loose-objects` real (CO-001, dritte Ausprägung).
// Hier ohne Pack-Mechanik: dieselbe Wirkung entsteht, wenn das lose Objekt
// fehlt. Anders als der geteilte Walker hinter Tree.Files() (der das still
// abschneidet) muss AllPaths hier abbrechen.
func TestAllPathsUnlesbarerUnterbaum(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "keep.md", "x\n")
	put(t, dir, "sub/a.md", "a\n")
	put(t, dir, "sub/b.md", "b\n")
	base := snapshot(t, wt, "base")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Vorbedingung: alle drei Pfade lesbar.
	baseAll, _, err := a.AllPaths(base, base)
	if err != nil {
		t.Fatal(err)
	}
	if len(baseAll) != 3 {
		t.Fatalf("Vorbedingung verfehlt: erwartet 3 Pfade, got %v", baseAll)
	}

	// Den Unterbaum "sub" unlesbar machen.
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(plumbing.NewHash(base))
	if err != nil {
		t.Fatal(err)
	}
	tree, err := commit.Tree()
	if err != nil {
		t.Fatal(err)
	}
	entry, err := tree.FindEntry("sub")
	if err != nil {
		t.Fatal(err)
	}
	h := entry.Hash.String()
	if err := os.Remove(filepath.Join(dir, ".git", "objects", h[:2], h[2:])); err != nil {
		t.Fatal(err)
	}

	// Kern der Regression: AllPaths bricht ab, statt den Unterbaum wortlos zu
	// überspringen.
	if _, _, err := a.AllPaths(base, base); err == nil {
		t.Fatal("unlesbarer Unterbaum still übersprungen — erwartet war ein Fehler")
	}
}

// TestAllPathsGitlinkWirdUebersprungen: ein Gitlink (Submodul) trägt den
// Commit-Hash eines FREMDEN Repos, keinen Blob oder Tree dieses Repos — der
// Walker darf ihn weder als Datei-Pfad melden noch versuchen, ihn als
// Unterbaum zu lesen (das wäre ein Fehlalarm ohne Ursache).
func TestAllPathsGitlinkWirdUebersprungen(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "keep.md", "x\n")
	first := snapshot(t, wt, "first")

	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(plumbing.NewHash(first))
	if err != nil {
		t.Fatal(err)
	}
	baseTree, err := commit.Tree()
	if err != nil {
		t.Fatal(err)
	}

	gitlinkHash := plumbing.NewHash("cafebabecafebabecafebabecafebabecafebabe")
	entries := append([]object.TreeEntry{
		{Name: "mod", Mode: filemode.Submodule, Hash: gitlinkHash},
	}, baseTree.Entries...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	newTree := &object.Tree{Entries: entries}
	obj := repo.Storer.NewEncodedObject()
	if err := newTree.Encode(obj); err != nil {
		t.Fatal(err)
	}
	treeHash, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		t.Fatal(err)
	}

	second := commitWithTree(t, repo, treeHash, plumbing.NewHash(first), "gitlink hinzugefügt")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	baseAll, headAll, err := a.AllPaths(first, second.String())
	if err != nil {
		t.Fatalf("Gitlink hätte nicht scheitern dürfen: %v", err)
	}
	if contains(headAll, "mod") {
		t.Fatal("Gitlink als Datei-Pfad gemeldet")
	}
	if len(headAll) != 1 || headAll[0] != "keep.md" {
		t.Fatalf("erwartet nur keep.md in HEAD-Menge, got %v", headAll)
	}
	if len(baseAll) != 1 || baseAll[0] != "keep.md" {
		t.Fatalf("erwartet nur keep.md in BASE-Menge, got %v", baseAll)
	}
}

// commitWithTree legt einen Commit mit gegebenem Tree und Parent direkt über
// den Storer an (ohne Worktree — für Bäume, die Konstrukte tragen, die ein
// Worktree-Checkout nicht herstellen kann, wie ein synthetischer Gitlink).
func commitWithTree(t *testing.T, repo *gogit.Repository, treeHash, parent plumbing.Hash, msg string) plumbing.Hash {
	t.Helper()
	sig := object.Signature{Name: "T", Email: "t@example.com", When: time.Unix(1700000001, 0)}
	c := &object.Commit{
		Author: sig, Committer: sig, Message: msg,
		TreeHash:     treeHash,
		ParentHashes: []plumbing.Hash{parent},
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

// TestFileAtUnlesbaresObjekt: Ein UNLESBARES Objekt darf nicht wie eine im
// Tree fehlende Datei gelesen werden. go-git meldet beide als
// object.ErrFileNotFound; nur der zweite Fall ist harmlos. Ohne die
// Unterscheidung im Adapter überspringt der Immutabilitäts-Vergleich die
// Datei still, und eine echte Core-Änderung verschwindet (DC-FA-VCS-001).
func TestFileAtUnlesbaresObjekt(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "docs/plan/adr/0001-a.md", "# ADR-0001\n\n**Status:** Accepted\n\nUrsprung.\n")
	base := snapshot(t, wt, "base")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Vorbedingung: lesbar, und der Inhalt stimmt.
	if _, ok, err := a.FileAt(base, "docs/plan/adr/0001-a.md"); err != nil || !ok {
		t.Fatalf("Vorbedingung verfehlt: ok=%v err=%v", ok, err)
	}

	// Das Blob unlesbar machen — dieselbe Wirkung wie ein Pack unter
	// unkanonischem Namen, nur ohne Pack-Mechanik im Test.
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(plumbing.NewHash(base))
	if err != nil {
		t.Fatal(err)
	}
	tree, err := commit.Tree()
	if err != nil {
		t.Fatal(err)
	}
	entry, err := tree.FindEntry("docs/plan/adr/0001-a.md")
	if err != nil {
		t.Fatal(err)
	}
	h := entry.Hash.String()
	if err := os.Remove(filepath.Join(dir, ".git", "objects", h[:2], h[2:])); err != nil {
		t.Fatal(err)
	}

	// Kern der Regression: Fehler, nicht (nil, false, nil).
	b, ok, err := a.FileAt(base, "docs/plan/adr/0001-a.md")
	if err == nil {
		t.Fatalf("unlesbares Objekt still übersprungen: ok=%v len=%d — erwartet war ein Fehler", ok, len(b))
	}
	if !strings.Contains(err.Error(), "nicht lesbar") {
		t.Errorf("Meldung nennt die Ursache nicht: %v", err)
	}

	// Gegenrichtung: eine im Tree WIRKLICH fehlende Datei bleibt (nil, false, nil).
	if _, ok, err := a.FileAt(base, "docs/plan/adr/0002-gibt-es-nicht.md"); err != nil || ok {
		t.Errorf("fehlende Datei falsch behandelt: ok=%v err=%v", ok, err)
	}
}

// TestFileAtEintragOhneBlob: Ein Tree-Eintrag, der keinen Datei-Inhalt
// bezeichnet — Verzeichnis oder Gitlink —, ist kein unlesbares Objekt. Wer
// beide gleich behandelt, tauscht ein stilles Übersehen gegen einen Fehlalarm
// mit irreführender Abhilfe (DC-FA-VCS-001).
func TestFileAtEintragOhneBlob(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "docs/plan/adr/0001-a.md", "# ADR-0001\n\n**Status:** Accepted\n\nUrsprung.\n")
	base := snapshot(t, wt, "base")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	// "docs" existiert im Tree, trägt aber keinen Datei-Inhalt.
	if _, ok, err := a.FileAt(base, "docs"); err != nil || ok {
		t.Errorf("Verzeichnis-Eintrag als unlesbares Objekt gelesen: ok=%v err=%v", ok, err)
	}

	// Und die Nachbarn bleiben, wie sie waren.
	if _, ok, err := a.FileAt(base, "docs/plan/adr/0001-a.md"); err != nil || !ok {
		t.Errorf("reguläre Datei falsch behandelt: ok=%v err=%v", ok, err)
	}
	if _, ok, err := a.FileAt(base, "gibt/es/nicht.md"); err != nil || ok {
		t.Errorf("fehlende Datei falsch behandelt: ok=%v err=%v", ok, err)
	}
}
