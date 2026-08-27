package rules

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
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

// §DC-FA-CODE-001.a — ignore-refs nimmt einen aufgelösten Ziel-Pfad
// REFERENZ-WEIT von der Existenz-Prüfung aus (Tombstone-Register,
// ADR-0025): per-Pfad, nicht datei-weit — der Eintrag stellt genau
// diesen Pfad in JEDER Datei still, übrige Verweise bleiben geprüft.
// Ohne Eintrag byte-identisch (DC-QA-02).
func TestCodepathsIgnoreRefs(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/adr.md":   "Frozen: `tools/weg.sh` und lebt: `../fehlt.md`.\n",
		"docs/slice.md": "Auch frozen: `tools/weg.sh`.\n",
	})
	roots := []string{"docs", "tools"}
	// Default-leer: der entfernte Pfad feuert in BEIDEN Dateien, dazu
	// der zweite fehlende Verweis → 3 Befunde (byte-identisch zum
	// Verhalten vor ignore-refs).
	cfgPlain := model.Config{Codepaths: model.CodepathsConfig{Roots: roots}}
	plain, err := Run(m, nil, cfgPlain, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plain.Findings) != 3 {
		t.Fatalf("ohne ignore-refs: %d Befunde, want 3", len(plain.Findings))
	}
	// Mit ignore-refs tools/weg.sh: der Tombstone-Pfad ist in beiden
	// Dateien still; nur der nicht-ignorierte ../fehlt.md bleibt
	// (Negative: per-Pfad, kein klassenweites Loch).
	cfgIgnore := model.Config{Codepaths: model.CodepathsConfig{Roots: roots, IgnoreRefs: []string{"tools/weg.sh"}}}
	ign, err := Run(m, nil, cfgIgnore, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range ign.Findings {
		got = append(got, fmt.Sprintf("%s %s %s", f.File, f.Target, f.Reason))
	}
	want := []string{"docs/adr.md ../fehlt.md codepath-missing"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("mit ignore-refs tools/weg.sh: Befunde = %v\nwant %v", got, want)
	}
	// Glob-Form: tools/*.sh deckt denselben Pfad ab → wieder nur der
	// lebende Verweis bleibt.
	cfgGlob := model.Config{Codepaths: model.CodepathsConfig{Roots: roots, IgnoreRefs: []string{"tools/*.sh"}}}
	g, err := Run(m, nil, cfgGlob, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Findings) != 1 {
		t.Fatalf("mit ignore-refs-Glob tools/*.sh: %d Befunde, want 1", len(g.Findings))
	}
}

// §DC-FA-CODE-001.a Schritt 5 / ADR-0025 — ein ignore-refs-Treffer
// unterdrückt ALLE drei Grund-Codes (codepath-missing, repo-escape,
// anchor-missing), weil das Match VOR Escape/Existenz/Anker greift.
// Verriegelt die Schritt-Reihenfolge: verschöbe man den ignored()-Aufruf
// hinter den Escape- oder Anker-Block, bräche genau dieser Test.
func TestCodepathsIgnoreRefsUnterdruecktEscapeUndAnker(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		// existierende Zieldatei OHNE den Anker 'fehlt'
		"docs/real.md": "# Titel\n",
		// a.md an der Wurzel: `../oben.md` eskaliert (rohes Ziel ==
		// aufgelöstes rel), `docs/real.md#fehlt` trägt einen toten Anker.
		"a.md": "Escape: `../oben.md` und toter Anker: `docs/real.md#fehlt`.\n",
	})
	roots := []string{"docs"}
	// Ohne ignore-refs: genau repo-escape + anchor-missing.
	plain, err := Run(m, nil, model.Config{Codepaths: model.CodepathsConfig{Roots: roots}}, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	var reasons []string
	for _, f := range plain.Findings {
		reasons = append(reasons, f.Reason)
	}
	if len(plain.Findings) != 2 {
		t.Fatalf("ohne ignore-refs: %d Befunde, want 2 (repo-escape + anchor-missing): %v", len(plain.Findings), reasons)
	}
	// Mit ignore-refs auf BEIDE aufgelösten Pfade: kein Befund — der
	// Treffer unterdrückt repo-escape UND anchor-missing, nicht nur
	// codepath-missing. Das Glob matcht den aufgelösten Wurzel-relativen
	// Pfad (Escape: `../oben.md`; Anker: `docs/real.md` nach Fragment-Abtrennung).
	cfg := model.Config{Codepaths: model.CodepathsConfig{
		Roots:      roots,
		IgnoreRefs: []string{"../oben.md", "docs/real.md"},
	}}
	ign, err := Run(m, nil, cfg, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ign.Findings) != 0 {
		t.Fatalf("mit ignore-refs: %d Befunde, want 0 (repo-escape + anchor-missing unterdrückt)", len(ign.Findings))
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

// DC-FA-CODE-001.a / DC-FA-ID-001.a: der Ventil-Marker in Inline-Code ist eine
// ERWÄHNUNG, keine Direktive — dieselbe geteilte Prosa-Antwort, die citations
// für seine Direktive bekommt. Ohne das nimmt eine Zeile, die das Ventil
// beschreibt, sich selbst aus dem Gate; im Bestand traf das fünf Zeilen des
// Lastenhefts, die den Marker dokumentieren.
//
// Beide Prüfpfade sind gedeckt — der Pfad in Inline-Code (codepaths) und die
// Kennung (ids); sie konsultieren denselben Satz, damit die Ausnahme nicht
// divergieren kann. Je Fall die Gegenprobe mit freiem Marker.
func TestIgnoreMarkerInInlineCodeIstErwaehnung(t *testing.T) {
	adr := model.IDPattern{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/", LinkPolicy: "always"}
	cases := []struct {
		name   string
		body   string
		module string
		want   []string
	}{
		{"codepaths erwähnt", "Der Marker `d-check:ignore` nimmt `./gibtsnicht.md` aus.\n",
			"codepaths", []string{"codepath-missing"}},
		{"codepaths frei", "Der Marker <!-- d-check:ignore --> nimmt `./gibtsnicht.md` aus.\n",
			"codepaths", nil},
		{"ids erwähnt", "Der Marker `d-check:ignore` nimmt `ADR-0042` aus.\n",
			"ids", []string{"id-unlinked"}},
		{"ids frei", "Der Marker <!-- d-check:ignore --> nimmt `ADR-0042` aus.\n",
			"ids", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := coretest.NewMemFS(map[string]string{
				"docs/a.md":               tc.body,
				"docs/plan/adr/0001-x.md": "Definitions-Ort\n",
			})
			cfg := model.Config{IDPatterns: []model.IDPattern{adr}}
			res, err := Run(m, nil, cfg, []string{tc.module})
			if err != nil {
				t.Fatalf("Run(%s): %v", tc.module, err)
			}
			var got []string
			for _, f := range res.Findings {
				got = append(got, f.Reason)
			}
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Fatalf("%s: %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// DC-FA-VER-001.a: bei versions bleibt die Erkennung ROH — und das ist eine
// BENANNTE GRENZE, keine andere Frage (ADR-0062). Seine Eingabe sind ALLE
// Zeilen, also eine Obermenge, die die Prosa-Zeilen einschliesst: auf der Zeile
// dieses Fixtures — freie Prosa, kein Fence — antwortet das Produkt damit
// zweifach, weil codepaths dieselbe Zeile gestrippt liest. Der Test nagelt die
// Divergenz fest, damit ein spaeterer Angleich eine Entscheidung ist und kein
// Nebeneffekt.
//
// diagrams ist hier NICHT gedeckt: seine Eingabe (Zeilen innerhalb eines Fence
// und die Oeffnungszeile) kennt gar kein Inline-Code, die Konstellation ist
// dort nicht konstruierbar.
func TestIgnoreMarkerBleibtRohBeiVersionsUndDiagrams(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": "Stand: v2.0.0\n",
		"docs/v.md":   "Pin `v1.0.0` mit `d-check:ignore` in Backticks.\n",
	})
	cfg := model.Config{Versions: model.VersionsConfig{Patterns: []model.VersionPattern{{
		PinPattern:  regexp.MustCompile(`v\d+\.\d+\.\d+`),
		CurrentFrom: "docs/src.md",
	}}}}
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatalf("Run(versions): %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("versions: %d Befund(e), want 0 — der Marker wirkt dort weiter roh (%v)",
			len(res.Findings), res.Findings)
	}
}

// DC-FA-CODE-001.a / DC-FA-ID-001.a: der Marker muss bei diesen zwei
// Konsumenten in einem HTML-KOMMENTAR stehen — das ist die Kommentar-Lexik
// von Markdown, und Markdown ist ihre Eingabe (ADR-0063). Eine blanke
// Erwähnung in Prosa wirkt nicht mehr; sie tat es bisher, und damit schaltete
// jeder Satz ueber das Ventil die Pruefung ab, ueber die er schreibt.
//
// Die Bedingung ist KONSERVATIV: ein `>` vor dem Marker im Kommentar laesst
// ihn nicht gelten. Ein verpasster Marker ist Falsch-Rot — laut; ein
// erfundener waere stilles Gruen.
func TestIgnoreMarkerBrauchtDieKommentarForm(t *testing.T) {
	adr := model.IDPattern{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/", LinkPolicy: "always"}
	cases := []struct {
		name   string
		body   string
		module string
		want   []string
	}{
		{"blank wirkt nicht (codepaths)", "Blanker Marker d-check:ignore und `./weg.md`.\n",
			"codepaths", []string{"codepath-missing"}},
		{"Kommentar wirkt (codepaths)", "Marker <!-- d-check:ignore --> und `./weg.md`.\n",
			"codepaths", nil},
		{"blank wirkt nicht (ids)", "Blanker Marker d-check:ignore und `ADR-0042`.\n",
			"ids", []string{"id-unlinked"}},
		{"Kommentar wirkt (ids)", "Marker <!-- d-check:ignore --> und `ADR-0042`.\n",
			"ids", nil},
		{"Kommentar mit > davor wirkt nicht", "Marker <!-- a > b d-check:ignore --> und `./weg.md`.\n",
			"codepaths", []string{"codepath-missing"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := coretest.NewMemFS(map[string]string{
				"docs/a.md":               tc.body,
				"docs/plan/adr/0001-x.md": "Definitions-Ort\n",
			})
			cfg := model.Config{IDPatterns: []model.IDPattern{adr}}
			res, err := Run(m, nil, cfg, []string{tc.module})
			if err != nil {
				t.Fatalf("Run(%s): %v", tc.module, err)
			}
			var got []string
			for _, f := range res.Findings {
				got = append(got, f.Reason)
			}
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Fatalf("%s: %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// DC-FA-VER-001.a: bei versions bleibt der Marker ein TOKEN — seine Eingabe
// sind ALLE Zeilen, auch die in Fences fremder Sprachen; eine
// Markdown-Kommentar-Form gaebe es dort nicht zu fordern. Zusammen mit der
// festgelegten Token-Form von diagrams (spec/spezifikation.md
// §DC-FA-DIAG-001.a) ist die Form damit EINGABE-abhaengig, nicht beliebig.
func TestIgnoreMarkerBleibtTokenBeiVersions(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": "Stand: v2.0.0\n",
		"docs/v.md":   "Blanker Marker d-check:ignore, Pin v1.0.0.\n",
	})
	cfg := model.Config{Versions: model.VersionsConfig{Patterns: []model.VersionPattern{{
		PinPattern:  regexp.MustCompile(`v\d+\.\d+\.\d+`),
		CurrentFrom: "docs/src.md",
	}}}}
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatalf("Run(versions): %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("versions: %d Befund(e), want 0 — der blanke Token wirkt dort weiter (%v)",
			len(res.Findings), res.Findings)
	}
}
