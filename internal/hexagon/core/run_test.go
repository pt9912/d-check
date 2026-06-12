package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// DC-FA-LINK-001: Happy/Boundary/Negative gegen In-Memory-FS.
func TestLinksModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md":   "[ok](b.md)\n[fehlt](nicht-da.md)\n```\n[imfence](auch-weg.md)\n```\n[escape](../../etc/passwd)\n[anker](#nur-anker)\n[extern](https://example.org)\n[frag](weg.md#x)",
		"docs/b.md":   "ziel",
		"docs/sub.md": "[leer](mit%20leer.md)",
	})
	m.files["docs/mit leer.md"] = "x"

	res, err := Run(m, nil, Config{}, []string{"links"})
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

	res, err := Run(m, nil, Config{}, []string{"links"})
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
	// explizite Wurzel darf die Repo-Wurzel nicht verlassen
	if _, err := DiscoverFiles(m, []string{"../anderes-repo"}, nil); err == nil {
		t.Fatal("Repo-Escape-Wurzel: erwarteter Fehler blieb aus")
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

// unreadableFS simuliert ein nicht lesbares Verzeichnis (z. B.
// root-eigene Build-Reste wie .gradle/): List darauf schlägt fehl.
type unreadableFS struct {
	*memFS
	deny string
}

func (u unreadableFS) List(relDir string) ([]driven.DirEntry, error) {
	if relDir == u.deny || strings.HasPrefix(relDir, u.deny+"/") {
		return nil, fmt.Errorf("open %s: permission denied", relDir)
	}
	return u.memFS.List(relDir)
}

// Ignore-Muster prunen den Verzeichnis-Abstieg: ein vollständig
// ignorierter, unlesbarer Teilbaum ist kein Laufzeitfehler
// (Adoptions-Befund pkcs11-course, slice-014; §scan.ignore).
func TestDiscoverFiles_IgnorePruntAbstieg(t *testing.T) {
	m := newMemFS(map[string]string{
		"README.md":         "x",
		"kaputt/krams.md":   "x",
		"kaputt/tief/t.md":  "x",
		"docs/a.md":         "x",
	})
	fs := unreadableFS{memFS: m, deny: "kaputt"}

	// ohne Ignore: der unlesbare Teilbaum bricht den Lauf ab
	if _, err := DiscoverFiles(fs, []string{"."}, nil); err == nil {
		t.Fatal("unlesbares Verzeichnis ohne Ignore: erwarteter Fehler blieb aus")
	}
	// mit `pfad/**`-Ignore wird der Teilbaum nicht betreten
	files, err := DiscoverFiles(fs, []string{"."}, []string{"kaputt/**"})
	if err != nil {
		t.Fatalf("Ignore-Muster prunt nicht: %v", err)
	}
	want := []string{"README.md", "docs/a.md"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("files = %v, want %v", files, want)
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
		res, err := Run(m, nil, Config{}, []string{"links"})
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

// Aktives external ohne verdrahteten Checker (nil) bleibt ein
// No-op — kein Fehler, keine Befunde (Kern-Tests laufen netzlos).
func TestExternalOhneChecker(t *testing.T) {
	m := newMemFS(map[string]string{"docs/a.md": "[x](https://example.org)"})
	res, err := Run(m, nil, Config{}, []string{"external", "links"})
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("res = %+v, err = %v", res, err)
	}
}

// DC-FA-CONF-002: Modul-lokaler Scan-Scope — die vier
// Akzeptanzkriterien gegen das In-Memory-FS.
func TestRun_ModulScope(t *testing.T) {
	m := newMemFS(map[string]string{
		"README.md":          "ADR-0042 nackt und [kaputt](fehlt.md)",
		"spec/s.md":          "ADR-0042 nackt in spec",
		"docs/adr/0042-x.md": "# Titel",
	})
	pattern := []IDPattern{{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/adr/"}}

	// Happy Path: ids nur auf spec/, links global — id-unlinked
	// stammt ausschließlich aus spec/, links prüft weiter alles.
	cfg := Config{
		Roots:      []string{"."},
		IDPatterns: pattern,
		Scopes:     map[string]*ScopeConfig{"ids": {Roots: []string{"spec"}}},
	}
	res, err := Run(m, nil, cfg, []string{"links", "ids"})
	if err != nil {
		t.Fatal(err)
	}
	var idFiles, linkFiles []string
	for _, f := range res.Findings {
		switch f.Rule {
		case "ids":
			idFiles = append(idFiles, f.File)
		case "links":
			linkFiles = append(linkFiles, f.File)
		}
	}
	if !reflect.DeepEqual(idFiles, []string{"spec/s.md"}) {
		t.Fatalf("id-unlinked außerhalb des Modul-Scopes: %v", idFiles)
	}
	if !reflect.DeepEqual(linkFiles, []string{"README.md"}) {
		t.Fatalf("links nicht mehr global: %v", linkFiles)
	}
	if res.FilesChecked != 3 {
		t.Fatalf("Vereinigungsmenge erwartet 3, got %d", res.FilesChecked)
	}

	// Boundary: ohne scope byte-identisch zum globalen Lauf.
	plain := Config{Roots: []string{"."}, IDPatterns: pattern}
	resA, err := Run(m, nil, plain, []string{"links", "ids"})
	if err != nil {
		t.Fatal(err)
	}
	withEmptyMap := plain
	withEmptyMap.Scopes = map[string]*ScopeConfig{}
	resB, err := Run(m, nil, withEmptyMap, []string{"links", "ids"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resA, resB) {
		t.Fatalf("ohne scope nicht identisch:\nA=%+v\nB=%+v", resA, resB)
	}
	if len(resA.Findings) != 3 { // 2× id-unlinked (README, spec) + 1× links
		t.Fatalf("globaler Lauf erwartet 3 Befunde, got %d", len(resA.Findings))
	}

	// Boundary: explizit leere roots-Liste prüft nichts (für ids),
	// Vereinigungsmenge bleibt der globale links-Scope.
	cfg.Scopes = map[string]*ScopeConfig{"ids": {Roots: []string{}}}
	res, err = Run(m, nil, cfg, []string{"links", "ids"})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if f.Rule == "ids" {
			t.Fatalf("leerer Scope erzeugte ids-Befund: %+v", f)
		}
	}
	if res.FilesChecked != 3 {
		t.Fatalf("Union bei leerem ids-Scope erwartet 3, got %d", res.FilesChecked)
	}

	// Negative: nicht existente bzw. Repo-Escape-Wurzel → Fehler
	// (Exit 2 im CLI), mit Modul-Präfix in der Meldung.
	for _, bad := range [][]string{{"handbuch"}, {"../raus"}} {
		cfg.Scopes = map[string]*ScopeConfig{"ids": {Roots: bad}}
		if _, err := Run(m, nil, cfg, []string{"links", "ids"}); err == nil {
			t.Fatalf("ungültige scope-Wurzel %v akzeptiert", bad)
		} else if !strings.Contains(err.Error(), "ids.scope") {
			t.Fatalf("Fehlermeldung ohne Modul-Kontext: %v", err)
		}
	}
}

// DC-FA-CONF-002: ein Modul-Scope kann Dateien umfassen, die der
// globale Scan nicht enthält (eigener Discover-Lauf, kein Filter).
func TestRun_ModulScopeAusserhalbGlobal(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "sauber",
		"extra/e.md": "[kaputt](fehlt.md)",
	})
	cfg := Config{
		Roots:  []string{"docs"},
		Scopes: map[string]*ScopeConfig{"links": {Roots: []string{"extra"}}},
	}
	res, err := Run(m, nil, cfg, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].File != "extra/e.md" {
		t.Fatalf("Scope außerhalb des globalen Scans nicht geprüft: %+v", res.Findings)
	}
	if res.FilesChecked != 2 { // docs/a.md (anchors global) + extra/e.md (links)
		t.Fatalf("Union erwartet 2, got %d", res.FilesChecked)
	}
}
