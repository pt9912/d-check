package git_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func statusOf(changes []driven.VCSChange, path string) (driven.VCSStatus, bool) {
	for _, c := range changes {
		if c.Path == path {
			return c.Status, true
		}
	}
	return 0, false
}

// TestRangeAndFileAt: Range-Diff (Added/Modified/Deleted) + FileAt an Refs.
func TestRangeAndFileAt(t *testing.T) {
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
	changes, err := a.ChangedPaths(first, second)
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]driven.VCSStatus{
		"keep.md": driven.VCSModified, "sub/gone.md": driven.VCSDeleted, "fresh.md": driven.VCSAdded,
	} {
		if got, ok := statusOf(changes, path); !ok || got != want {
			t.Fatalf("%s: %c erwartet, got %c ok=%v (alle %v)", path, want, got, ok, changes)
		}
	}

	if b, ok, err := a.FileAt(first, "keep.md"); err != nil || !ok || string(b) != "v1\n" {
		t.Fatalf("FileAt(first,keep) = %q,%v,%v", b, ok, err)
	}
	if _, ok, err := a.FileAt(second, "sub/gone.md"); err != nil || ok {
		t.Fatalf("FileAt(second, gelöscht) erwartet ok=false, got ok=%v err=%v", ok, err)
	}
	if _, err := a.ChangedPaths("0000000000000000000000000000000000000000", second); err == nil {
		t.Fatal("unauflösbare Basis hätte fail-closed liefern müssen")
	}
	if _, _, err := a.FileAt("0000000000000000000000000000000000000000", "keep.md"); err == nil {
		t.Fatal("unauflösbares Ref hätte fail-closed liefern müssen")
	}
}

// TestStaged: staged-Diff (HEAD-Tree vs. Index) + FileAt(Index) + ohne HEAD leer.
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
	changes, err := a.ChangedPaths("HEAD", driven.IndexRef)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := statusOf(changes, "adr.md"); !ok || got != driven.VCSModified {
		t.Fatalf("staged: adr.md Modified erwartet, got %v", changes)
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
	changes, err := a.ChangedPaths("HEAD", driven.IndexRef)
	if err != nil || changes != nil {
		t.Fatalf("ohne HEAD: leer/kein Fehler erwartet, got %v / %v", changes, err)
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
