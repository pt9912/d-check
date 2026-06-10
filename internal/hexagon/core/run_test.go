package core

import (
	"fmt"
	"reflect"
	"testing"
)

// DC-FA-LINK-001: Happy/Boundary/Negative gegen In-Memory-FS.
func TestLinksModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md":   "[ok](b.md)\n[fehlt](nicht-da.md)\n```\n[imfence](auch-weg.md)\n```\n[escape](../../etc/passwd)\n[anker](#nur-anker)\n[extern](https://example.org)\n[frag](weg.md#x)",
		"docs/b.md":   "ziel",
		"docs/sub.md": "[leer](mit%20leer.md)",
	})
	m.files["docs/mit leer.md"] = "x"

	res, err := Run(m, Config{}, []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s %s", f.File, f.Line, f.Target, f.Reason))
	}
	want := []string{
		"docs/a.md:2 nicht-da.md target-missing",
		"docs/a.md:6 ../../etc/passwd repo-escape",
		// weg.md#x: Datei fehlt → genau EIN Befund aus links (Fragment
		// abgetrennt; anchors schweigt — DC-FA-ANCH-001 Boundary)
		"docs/a.md:9 weg.md#x target-missing",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
	if res.FilesChecked != 4 {
		t.Fatalf("FilesChecked = %d, want 4", res.FilesChecked)
	}
}

// DC-FA-LINK-002: Symlink-Vorrang, genau ein Befund pro Ziel.
func TestSymlinkAblehnung(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "[intern](link-intern.md)\n[außen](link-raus)\n[normal](b.md)",
		"docs/b.md": "x",
	})
	m.symlinks["docs/link-intern.md"] = true // zeigt intern — trotzdem Befund
	m.symlinks["docs/link-raus"] = true      // zeigt nach außen

	res, err := Run(m, Config{}, []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 2 {
		t.Fatalf("Befunde = %d, want 2: %+v", len(res.Findings), res.Findings)
	}
	for _, f := range res.Findings {
		if f.Reason != ReasonSymlink {
			t.Errorf("Reason = %s, want symlink (genau ein Befund pro Ziel)", f.Reason)
		}
	}
}

// DC-FA-SCAN-001: Default-Wurzeln, Skip-Dirs, Ignore, explizite Wurzeln.
func TestDiscoverFiles(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md":                   "x",
		"docs/archive/alt.md":         "x",
		"spec/s.md":                   "x",
		"node_modules/x/README.md":    "x",
		"docs/node_modules/tief.md":   "x",
		"README.md":                   "x",
		"irgendwo/anders.md":          "x",
		"docs/build/generated.md":     "x",
	})

	files, err := DiscoverFiles(m, nil, []string{"docs/archive/**"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"README.md", "docs/a.md", "spec/s.md"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("files = %v, want %v", files, want)
	}

	// explizite Wurzel muss existieren (Negative-Kriterium)
	if _, err := DiscoverFiles(m, []string{"handbuch"}, nil); err == nil {
		t.Fatal("fehlende explizite Wurzel: erwarteter Fehler blieb aus")
	}
	// explizite Wurzel ersetzt Defaults
	files, err = DiscoverFiles(m, []string{"irgendwo"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(files, []string{"irgendwo/anders.md"}) {
		t.Fatalf("files = %v", files)
	}
}

// DC-FA-CLI-002: Modul-Auflösung.
func TestEffectiveModules(t *testing.T) {
	mods, err := EffectiveModules(Config{}, nil, nil)
	if err != nil || !reflect.DeepEqual(mods, []string{"anchors", "links"}) {
		t.Fatalf("Default = %v (%v)", mods, err)
	}
	// CLI-Präzedenz: Config aktiviert, CLI deaktiviert
	mods, err = EffectiveModules(Config{Modules: []string{"links", "ids"}}, nil, []string{"ids"})
	if err != nil || !reflect.DeepEqual(mods, []string{"links"}) {
		t.Fatalf("Präzedenz = %v (%v)", mods, err)
	}
	// unbekanntes Modul → Fehler mit Liste
	if _, err := EffectiveModules(Config{}, []string{"foo"}, nil); err == nil {
		t.Fatal("unbekanntes Modul: erwarteter Fehler blieb aus")
	}
}

// DC-QA-02: identische Eingabe ⇒ identische, sortierte Ausgabe.
func TestDeterminismus(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/z.md": "[f1](f1.md)\n[f2](f2.md)",
		"docs/a.md": "[f3](f3.md)",
		"spec/m.md": "[f4](f4.md)",
	})
	var prev []Finding
	for i := 0; i < 10; i++ {
		res, err := Run(m, Config{}, []string{"links"})
		if err != nil {
			t.Fatal(err)
		}
		if prev != nil && !reflect.DeepEqual(prev, res.Findings) {
			t.Fatalf("Lauf %d weicht ab", i)
		}
		prev = res.Findings
	}
	if prev[0].File != "docs/a.md" || prev[len(prev)-1].File != "spec/m.md" {
		t.Fatalf("Sortierung verletzt: %+v", prev)
	}
}

func TestSortFindingsDedupe(t *testing.T) {
	f := Finding{File: "a.md", Line: 1, Rule: "links", Target: "x", Reason: "target-missing"}
	out := SortFindings([]Finding{f, f, {File: "a.md", Line: 1, Rule: "links", Target: "x", Reason: "symlink"}})
	if len(out) != 2 {
		t.Fatalf("Dedupe: %d Befunde, want 2", len(out))
	}
}

// Nicht implementierte Module werden gemeldet und übersprungen.
func TestSkippedModules(t *testing.T) {
	m := newMemFS(map[string]string{"docs/a.md": "x"})
	res, err := Run(m, Config{}, []string{"ids", "links"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res.SkippedModules, []string{"ids"}) {
		t.Fatalf("Skipped = %v", res.SkippedModules)
	}
}
