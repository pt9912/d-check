package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveLink(t *testing.T) {
	cases := []struct {
		source, target, wantResolved, wantFragment string
	}{
		{"docs/plan/planning/observations.md", "done/slice-190-a.md", "docs/plan/planning/done/slice-190-a.md", ""},
		{"docs/plan/planning/done/slice-185.md", "slice-184-a.md", "docs/plan/planning/done/slice-184-a.md", ""},
		{"docs/plan/planning/in-progress/roadmap.md", "../done/slice-190-a.md", "docs/plan/planning/done/slice-190-a.md", ""},
		{"harness/conventions/done/MR-057.md", "../../../docs/plan/planning/done/slice-190-a.md", "docs/plan/planning/done/slice-190-a.md", ""},
		{"docs/plan/planning/observations.md", "done/slice-190-a.md#anker", "docs/plan/planning/done/slice-190-a.md", "#anker"},
	}
	for _, c := range cases {
		gotResolved, gotFragment := resolveLink(c.source, c.target)
		if gotResolved != c.wantResolved || gotFragment != c.wantFragment {
			t.Errorf("resolveLink(%q, %q) = (%q, %q), want (%q, %q)",
				c.source, c.target, gotResolved, gotFragment, c.wantResolved, c.wantFragment)
		}
	}
}

// TestRewriteFile ist die Umkehr-Probe fuer den Verweis-Nachzug (BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt):
// ein Link, der NICHT im Move-Set steht, muss unveraendert bleiben -- kehrt
// man die ok-Pruefung in RewriteFile um (jeder Link "trifft"), wird dieser
// Test rot, weil der unbeteiligte Link dann faelschlich umgeschrieben wuerde.
func TestRewriteFile(t *testing.T) {
	moves := []Move{
		{Old: "docs/plan/planning/done/slice-190-a.md", New: "docs/plan/planning/done/welle-87/slice-190-a.md"},
	}
	content := "Siehe [slice-190](done/slice-190-a.md) und [slice-999](done/slice-999-b.md)."
	got, n := RewriteFile("docs/plan/planning/observations.md", content, moves)
	if n != 1 {
		t.Fatalf("erwartet genau 1 Ersetzung, got %d", n)
	}
	if !containsAll(got, "(done/welle-87/slice-190-a.md)", "(done/slice-999-b.md)") {
		t.Fatalf("unerwarteter Inhalt: %q", got)
	}
}

func TestRewriteFile_MitAnker(t *testing.T) {
	moves := []Move{
		{Old: "docs/plan/planning/done/slice-190-a.md", New: "docs/plan/planning/done/welle-87/slice-190-a.md"},
	}
	content := "[Text](done/slice-190-a.md#abschnitt)"
	got, n := RewriteFile("docs/plan/planning/observations.md", content, moves)
	if n != 1 {
		t.Fatalf("erwartet 1 Ersetzung, got %d", n)
	}
	if got != "[Text](done/welle-87/slice-190-a.md#abschnitt)" {
		t.Fatalf("Fragment nicht erhalten: %q", got)
	}
}

func TestRewriteRepo(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs/plan/planning/observations.md"),
		"[slice-190](done/slice-190-a.md)")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-185-b.md"),
		"[slice-190](slice-190-a.md)")
	writeFile(t, filepath.Join(root, "docs/plan/planning/done/slice-999-c.md"),
		"kein betroffener Link hier")

	moves := []Move{
		{Old: "docs/plan/planning/done/slice-190-a.md", New: "docs/plan/planning/done/welle-87/slice-190-a.md"},
	}
	hits, err := RewriteRepo(root, moves)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 2 {
		t.Fatalf("erwartet 2 betroffene Dateien, got %d: %v", len(hits), hits)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
