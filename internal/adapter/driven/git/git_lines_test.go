package git_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	gitadapter "github.com/pt9912/d-check/internal/adapter/driven/git"
)

// linesFixture legt zwei Commits an, deren Unterschied alle Randformen der
// Zeilen-Summe trägt (DC-FA-MET-001): Ändern, Löschen, Hinzufügen, leere Datei,
// Umbenennung mit Inhaltsänderung, Binärdatei, Datei außerhalb der Pfade und
// Datei im ausgenommenen Teilbaum.
func linesFixture(t *testing.T) (dir, first, second string) {
	t.Helper()
	dir, wt := repoAt(t)
	put(t, dir, "src/a.txt", "l1\nl2\nl3\nl4\nl5\n")
	put(t, dir, "src/del.txt", "d1\nd2\nd3\n")
	put(t, dir, "src/same.txt", "s1\ns2\n")
	put(t, dir, "src/empty.txt", "")
	put(t, dir, "src/old.txt", "r1\nr2\nr3\nr4\n")
	put(t, dir, "src/bin.dat", "AB\x00CD\n")
	put(t, dir, "other/o.txt", "o1\no2\n")
	put(t, dir, "docs/vertraege/v.txt", "v1\nv2\n")
	put(t, dir, "src/gone-empty.txt", "")
	put(t, dir, "src/crlf.txt", "a\r\nb\r\nc")
	put(t, dir, "src/mode.sh", "echo\n")
	first = snapshot(t, wt, "base")

	put(t, dir, "src/a.txt", "l1\nL2\nl3\nl4\nl5\nl6\n")
	if _, err := wt.Remove("src/del.txt"); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "src/added.txt", "n1\nn2\nn3\nn4\n")
	put(t, dir, "src/empty.txt", "x\n")
	if err := os.Remove(filepath.Join(dir, "src/old.txt")); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "src/new.txt", "r1\nr2\nr3\nr4\nr5\n")
	put(t, dir, "src/bin.dat", "AB\x00EF\x00\n")
	put(t, dir, "other/o.txt", "o1\nO22\n")
	put(t, dir, "docs/vertraege/v.txt", "v1\nV22\n")
	put(t, dir, "src/größe leer.txt", "")
	if _, err := wt.Remove("src/gone-empty.txt"); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "src/crlf.txt", "a\r\nB\r\nc\r\nd\r\n")
	put(t, dir, "src/ü ber.txt", "u1\nu2\n")
	if err := os.Chmod(filepath.Join(dir, "src/mode.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	second = snapshot(t, wt, "head")
	return dir, first, second
}

// Die Zeilen-Summe gleicht `git diff --numstat --no-renames` über dieselben
// Pfade. Die erwartete Summe 27 ist die Summe des Kommandos am gleichen Inhalt:
// a.txt 2+1, added.txt 4+0, bin.dat binär, del.txt 0+3, empty.txt 1+0,
// new.txt 5+0, old.txt 0+4 (Rename-Erkennung aus), crlf.txt 3+2 (CRLF, ohne
// Schluss-Umbruch), „ü ber.txt“ 2+0; leere Datei neu und gelöscht, ein reiner
// Moduswechsel und die Binärdatei tragen 0 bei.
func TestLineChanges_GleichNumstat(t *testing.T) {
	dir, first, second := linesFixture(t)
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := a.LineChanges(first, second, []string{"src/", "docs/", "src/a.txt", "nope/"}, []string{"docs/vertraege/"})
	if err != nil {
		t.Fatal(err)
	}
	if stat.Lines != 27 {
		t.Fatalf("Summe = %d, erwartet 27 (git diff --numstat --no-renames)", stat.Lines)
	}
	if !reflect.DeepEqual(stat.Binary, []string{"src/bin.dat"}) {
		t.Fatalf("Binärdateien = %v", stat.Binary)
	}
	if !reflect.DeepEqual(stat.Missing, []string{"nope/"}) {
		t.Fatalf("Pfade ohne Vorkommen = %v", stat.Missing)
	}
}

// Ein Pfad ist literal: ein Präfix ohne Schrägstrich trifft keinen längeren
// Namen, ein Glob-Zeichen ist ein Zeichen des Namens, und dieselbe Datei über
// Datei- und Verzeichnispfad zählt einmal.
func TestLineChanges_LiteralePfade(t *testing.T) {
	dir, first, second := linesFixture(t)
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"src/a", "src/*", "src/a.tx"} {
		st, err := a.LineChanges(first, second, []string{p}, nil)
		if err != nil || st.Lines != 0 || len(st.Missing) != 1 {
			t.Fatalf("Pfad %q darf nichts treffen: %+v, %v", p, st, err)
		}
	}
	one, _ := a.LineChanges(first, second, []string{"src/a.txt"}, nil)
	dup, _ := a.LineChanges(first, second, []string{"src/a.txt", "src/a.txt"}, nil)
	under, _ := a.LineChanges(first, second, []string{"src/"}, nil)
	both, _ := a.LineChanges(first, second, []string{"src/a.txt", "src/"}, nil)
	if one.Lines != 3 || dup.Lines != 3 || both.Lines != under.Lines {
		t.Fatalf("dieselbe Datei zählt einmal: Datei %d, doppelt %d, Verzeichnis %d, beides %d", one.Lines, dup.Lines, under.Lines, both.Lines)
	}
	same, _ := a.LineChanges(first, second, []string{"src/same.txt"}, nil)
	if same.Lines != 0 || len(same.Missing) != 0 {
		t.Fatalf("vorhandener, unveränderter Pfad: Summe 0, nicht fehlend: %+v", same)
	}
}

// Basis gleich HEAD ist 0; eine Basis als Kurz-Hash löst auf; eine
// unauflösbare Basis und eine Basis, die kein Vorfahr ist, sind Fehler.
func TestLineChanges_BasisFailClosed(t *testing.T) {
	dir, first, second := linesFixture(t)
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st, err := a.LineChanges(second, second, []string{"src/"}, nil); err != nil || st.Lines != 0 {
		t.Fatalf("Basis gleich HEAD: %+v, %v", st, err)
	}
	if st, err := a.LineChanges(first[:7], second, []string{"src/"}, []string{"docs/"}); err != nil || st.Lines != 27 {
		t.Fatalf("Kurz-Hash: %+v, %v", st, err)
	}
	if _, err := a.LineChanges("0000000000000000000000000000000000000000", second, []string{"src/"}, nil); err == nil {
		t.Fatal("unauflösbare Basis ⇒ Fehler")
	}
	_, err = a.LineChanges(second, first, []string{"src/"}, nil)
	if err == nil || !strings.Contains(err.Error(), "kein vorfahr") {
		t.Fatalf("Basis ist kein Vorfahr ⇒ Fehler, got %v", err)
	}
	if _, err := a.LineChanges(first, "nicht-vorhanden", []string{"src/"}, nil); err == nil {
		t.Fatal("unauflösbare Spitze ⇒ Fehler")
	}
}

// Eine Basis auf einem Geschwister-Zweig ist kein Vorfahr der Spitze, auch
// wenn beide denselben Elternteil haben.
func TestLineChanges_BasisAufGeschwisterZweig(t *testing.T) {
	dir, first, second := linesFixture(t)
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if err := wt.Checkout(&gogit.CheckoutOptions{Hash: plumbing.NewHash(first), Branch: plumbing.NewBranchReferenceName("side"), Create: true}); err != nil {
		t.Fatal(err)
	}
	put(t, dir, "side.txt", "s\n")
	side := snapshot(t, wt, "side")
	a, err := gitadapter.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = a.LineChanges(side, second, []string{"src/"}, nil)
	if err == nil || !strings.Contains(err.Error(), "kein vorfahr") {
		t.Fatalf("Geschwister-Zweig ⇒ kein Vorfahr, got %v", err)
	}
}
