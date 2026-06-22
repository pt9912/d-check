package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"fmt"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

// DC-FA-CODE-001 Happy/Negative über Run: existierender Backtick-Pfad
// (Präfix- und ./-Form) → kein Befund; fehlendes Ziel →
// codepath-missing; Repo-Escape → repo-escape; Fence wird nicht
// geprüft.
func TestCodepathsModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Siehe `docs/b.md` und `./b.md` sowie `../README.md`.\n" +
			"Kaputt: `../fehlt.md` und Escape: `../../etc/passwd`.\n" +
			"```\n`../auch-weg.md` im Fence\n```\n" +
			"Kein Pfad: `docs/{a,b}.md`, `a -> b`, `irgendwas`, `# anker`, `http://x`.\n" +
			"## Heading mit `../weg.md` ist Titel, keine Referenz\n" +
			"Linktext-Span: [`../auch-nicht.md`](b.md) prüft das Modul links\n",
		"docs/b.md": "x",
		"README.md": "x",
	})
	cfg := model.Config{Codepaths: model.CodepathsConfig{Roots: []string{"docs"}}}
	res, err := Run(m, nil, cfg, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s", f.Line, f.Target, f.Reason))
	}
	want := []string{ // sortiert: Target bytewise ('/' < 'f')
		"2 ../../etc/passwd repo-escape",
		"2 ../fehlt.md codepath-missing",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
}

// DC-FA-CODE-001 Boundary: der Marker d-check:ignore stellt NUR
// dieses Modul still — der Link-Befund derselben Zeile bleibt.
func TestCodepathsIgnoreMarkerNurDiesesModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Beispiel `../../etc/passwd` und [kaputt](fehlt.md) <!-- d-check:ignore (Angriffs-Beispiel) -->\n" +
			"Ohne Marker: `../fehlt.md`\n",
	})
	res, err := Run(m, nil, model.Config{}, []string{"codepaths", "links"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s", f.Line, f.Rule, f.Reason))
	}
	want := []string{
		"1 links target-missing",      // Marker wirkt NICHT auf links
		"2 codepaths codepath-missing", // ohne Marker greift codepaths
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
}

// §DC-FA-CODE-001.a — exempt-paths nimmt eine ganze Datei von der
// codepaths-Prüfung aus (Glob wie scan.ignore; Vorbild ids-Ventil,
// slice-043). Ohne Ventil byte-identisch (DC-QA-02).
func TestCodepathsExemptPaths(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/reviews/r1.md": "Report zitiert `../fehlt.md` und `docs/auch-weg.md`.\n",
	})
	cfgPlain := model.Config{Codepaths: model.CodepathsConfig{Roots: []string{"docs"}}}
	plain, err := Run(m, nil, cfgPlain, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plain.Findings) != 2 {
		t.Fatalf("ohne exempt-paths: %d Befunde, want 2", len(plain.Findings))
	}
	cfgExempt := model.Config{Codepaths: model.CodepathsConfig{Roots: []string{"docs"}, ExemptPaths: []string{"docs/reviews/**"}}}
	exempt, err := Run(m, nil, cfgExempt, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(exempt.Findings) != 0 {
		t.Fatalf("mit exempt-paths docs/reviews/**: %d Befunde, want 0", len(exempt.Findings))
	}
}

// §DC-FA-CODE-001.a Schritt 3+5 — Normalisierung und Anker-Prüfung
// (gleiches Slug-Verfahren wie anchors, geteilter Cache).
func TestCodepathsNormalisierungUndAnker(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Zitiert: `\"./b.md\",` und `./b.md#zweck` sowie `./b.md#gibt-es-nicht`.\n",
		"docs/b.md": "# Zweck\n",
	})
	res, err := Run(m, nil, model.Config{}, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != ReasonAnchorMissing ||
		res.Findings[0].Target != "./b.md#gibt-es-nicht" || res.Findings[0].Rule != "codepaths" {
		t.Fatalf("Befunde = %+v (genau ein anchor-missing erwartet)", res.Findings)
	}
}

// §DC-FA-CODE-001.a — codepaths prüft den Anker gegen die gemeinsame
// Anker-Menge inkl. Inline-HTML-Anker (slice-022): Treffer auf eine
// HTML-id → kein Befund, Fehlschlag → anchor-missing.
func TestCodepathsHTMLAnker(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Siehe `./b.md#html-id` und `./b.md#fehlt`.\n",
		"docs/b.md": "<div id=\"html-id\">Inhalt</div>\n",
	})
	res, err := Run(m, nil, model.Config{}, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != ReasonAnchorMissing ||
		res.Findings[0].Target != "./b.md#fehlt" || res.Findings[0].Rule != "codepaths" {
		t.Fatalf("Befunde = %+v (genau ein anchor-missing ./b.md#fehlt erwartet)", res.Findings)
	}
}

// classifyCodepath: konservative Erkennung im Detail.
func TestClassifyCodepath(t *testing.T) {
	roots := []string{"docs/", "tools"}
	cases := []struct {
		v       string
		rootRel bool
		ok      bool
	}{
		{"./x.md", false, true},
		{"../x.md", false, true},
		{"docs/x.md", true, true},
		{"tools/x.sh", true, true}, // Präfix ohne Slash deklariert
		{"docs2/x.md", false, false},
		{"", false, false},
		{"a b", false, false},
		{"docs/*.md", false, false},
		{"…/x", false, false},
		{"//server/x", false, false},
		{"#anker", false, false},
		{"mailto:a@b", false, false},
	}
	for _, c := range cases {
		rootRel, ok := classifyCodepath(c.v, roots)
		if rootRel != c.rootRel || ok != c.ok {
			t.Errorf("classifyCodepath(%q) = (%v,%v), want (%v,%v)", c.v, rootRel, ok, c.rootRel, c.ok)
		}
	}
	if got := normalizeCodepath(` "./a.md", `); got != "./a.md" {
		t.Errorf("normalizeCodepath = %q", got)
	}
	// Datei:Zeile-Konvention referenziert die Datei
	if got := normalizeCodepath("spec/x.md:290"); got != "spec/x.md" {
		t.Errorf("Zeilen-Suffix: %q", got)
	}
	if got := normalizeCodepath("a:b"); got != "a:b" {
		t.Errorf("kein Zeilen-Suffix: %q", got)
	}
}
