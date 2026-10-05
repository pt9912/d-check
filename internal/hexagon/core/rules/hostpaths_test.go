package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"reflect"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

// DC-FA-HOST-001: Happy/Boundary/Negative gegen In-Memory-FS.
// (Die Beispiel-Pfade leben in Test-Strings, nicht in Doku — das
// Modul selbst prüft nur Markdown.)
func TestHostpathsModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		// Happy: relative und Repo-Wurzel-absolute Angaben
		"docs/ok.md": "siehe [a](b.md) und `docs/x.md` sowie (/docs/y.md)\nhttps://example.org/home/x bleibt URL",
		"docs/b.md":  "x",
		// Boundary: Host-Pfad im Fence
		"docs/fence.md": "```\ncd /" + "home/alice/repo\n```\ndanach prosa",
		// Negative: Prosa und Inline-Code
		"docs/leak.md": "liegt unter /" + "home/alice/repo, fertig\nim Code: `/" + "mnt/data/token.db` und C:\\Users\\a\\x sowie \\\\srv\\share\\y",
	})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		if f.Reason != model.ReasonHostpathForbidden {
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
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "unter /" + "srv/data/x liegt es; /" + "home/y ist hier erlaubt",
	})
	cfg := model.Config{
		Roots:     []string{"docs"},
		Hostpaths: model.HostpathsConfig{Prefixes: []string{"srv"}},
	}
	res, err := Run(m, nil, cfg, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "/"+"srv/data/x" {
		t.Fatalf("Findings = %+v", res.Findings)
	}
}

// DC-FA-HOST-001 (Tilde): ein Home-relativer Pfad mit Nicht-Punkt-Segment
// wird in voller Form und genau einmal gemeldet; Werkzeug-Konventionen
// unter ~/.<name>, die nackte Tilde, ~user/, URL-Tilden und Fences bleiben still.
func TestHostpathsTilde(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/t.md": "Konfig in ~/.claude/settings.json bleibt still\n" +
			"Repo unter ~/src/projekt/x.md, fertig\n" +
			"auch ~/" + "Development/d-check/README.md\n" +
			"ungefähr ~5 % und ~/ allein\n" +
			"Fremdnutzer ~alice/x bleibt\n" +
			"URL https://example.org/~/x und https://example.org/~bob/y\n" +
			"im Code: `~/code/z`\n" +
			"Pfad~/x klebt\n" +
			"```\ncd ~/src/fence\n```\n",
	})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.Target)
	}
	want := []string{
		"~/src/projekt/x.md",
		"~/" + "Development/d-check/README.md",
		"~/code/z",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// DC-FA-HOST-001 (Ventil): exempt-targets nimmt Unix- und Tilde-Funde aus,
// deren normalisierter Pfad ein Glob trifft; Windows-/UNC-Funde bleiben fest.
func TestHostpathsExemptTargets(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "~/Library/Mobile/x und ~/src/y sowie /" + "mnt/data/z und /" +
			"mnt/other/q und C:\\Users\\a\\x",
	})
	cfg := model.Config{
		Roots: []string{"docs"},
		Hostpaths: model.HostpathsConfig{
			ExemptTargets: []string{"~/Library/**", "/" + "mnt/data/**", "C:*"},
		},
	}
	res, err := Run(m, nil, cfg, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.Target)
	}
	want := []string{"/" + "mnt/other/q", "C:\\Users\\a\\x", "~/src/y"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}
