package core

import (
	"reflect"
	"testing"
)

// DC-FA-HOST-001: Happy/Boundary/Negative gegen In-Memory-FS.
// (Die Beispiel-Pfade leben in Test-Strings, nicht in Doku — das
// Modul selbst prüft nur Markdown.)
func TestHostpathsModul(t *testing.T) {
	m := newMemFS(map[string]string{
		// Happy: relative und Repo-Wurzel-absolute Angaben
		"docs/ok.md": "siehe [a](b.md) und `docs/x.md` sowie (/docs/y.md)\nhttps://example.org/home/x bleibt URL",
		"docs/b.md":  "x",
		// Boundary: Host-Pfad im Fence
		"docs/fence.md": "```\ncd /" + "home/alice/repo\n```\ndanach prosa",
		// Negative: Prosa und Inline-Code
		"docs/leak.md": "liegt unter /" + "home/alice/repo, fertig\nim Code: `/" + "mnt/data/token.db` und C:\\Users\\a\\x sowie \\\\srv\\share\\y",
	})
	res, err := Run(m, nil, Config{Roots: []string{"docs"}}, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		if f.Reason != ReasonHostpathForbidden {
			t.Fatalf("unerwarteter Reason: %+v", f)
		}
		got = append(got, f.Target)
	}
	want := []string{
		"/" + "home/alice/repo",
		"/" + "mnt/data/token.db", // Inline-Code wird mitgeprüft, der Span-Backtick gehört nicht zum Pfad
		"C:\\Users\\a\\x",
		"\\\\srv\\share\\y",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
	for _, f := range res.Findings {
		if f.File != "docs/leak.md" {
			t.Fatalf("Befund außerhalb leak.md: %+v", f)
		}
	}
}

// Konfigurierbare Präfixliste ersetzt den Default.
func TestHostpathsPrefixesKonfigurierbar(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "unter /" + "srv/data/x liegt es; /" + "home/y ist hier erlaubt",
	})
	cfg := Config{
		Roots:     []string{"docs"},
		Hostpaths: HostpathsConfig{Prefixes: []string{"srv"}},
	}
	res, err := Run(m, nil, cfg, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "/"+"srv/data/x" {
		t.Fatalf("Findings = %+v", res.Findings)
	}
}
