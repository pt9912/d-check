package git_test

import (
	"testing"

	gitadapter "github.com/pt9912/d-check/internal/adapter/driven/git"
)

// TestTrackedPaths: der Index ist die Wahrheit (DC-FA-TRK-001) — committete
// UND frisch gestagte Dateien sind getrackt, eine nur im Arbeitsbaum
// liegende Datei nicht.
func TestTrackedPaths(t *testing.T) {
	dir, wt := repoAt(t)
	put(t, dir, "committed.md", "x")
	snapshot(t, wt, "c1")
	put(t, dir, "staged.md", "x")
	if _, err := wt.Add("staged.md"); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "untracked.md", "x")

	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	set, err := a.TrackedPaths()
	if err != nil {
		t.Fatal(err)
	}
	if !set["committed.md"] {
		t.Fatalf("committed.md fehlt im Index-Set: %v", set)
	}
	if !set["staged.md"] {
		t.Fatalf("staged.md (gestagt, nie committet) fehlt im Index-Set — Index-Wahrheit verletzt: %v", set)
	}
	if set["untracked.md"] {
		t.Fatalf("untracked.md fälschlich im Index-Set: %v", set)
	}
}

// TestTrackedPathsFrischesRepo: ein Repo ohne jeden Eintrag liefert eine
// leere Menge (kein Fehler) — fail-closed bleibt dem unlesbaren Index
// vorbehalten.
func TestTrackedPathsFrischesRepo(t *testing.T) {
	dir, _ := repoAt(t)
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	set, err := a.TrackedPaths()
	if err != nil {
		t.Fatalf("err = %v, want nil (leerer Index ist kein Fehler)", err)
	}
	if len(set) != 0 {
		t.Fatalf("Set = %v, want leer", set)
	}
}
