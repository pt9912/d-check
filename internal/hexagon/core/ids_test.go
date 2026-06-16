package core

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// DC-FA-ID-001 Happy/Boundary/Negative über Run: verlinkt → kein
// Befund; Inline-Code/Fence → linkpflichtfrei; nackt → id-unlinked;
// Ziel-Teil eines Links und Bildreferenzen → kein Fließtext.
func TestIDsModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "[ADR-0042](plan/adr/0042-beispiel.md)\n" + // verlinkt
			"`ADR-0042` in Inline-Code\n" + // linkpflichtfrei
			"nacktes ADR-0042 im Fließtext\n" + // Befund
			"```\nADR-0042 im Fence\n```\n" + // linkpflichtfrei
			"[Entscheidung](plan/adr/ADR-0042-beispiel.md)\n" + // Ziel-Teil
			"![Diagramm ADR-0042](plan/adr/bild.png)\n", // Bildreferenz
		"docs/plan/adr/0042-beispiel.md": "x",
	})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"},
	}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s", f.Line, f.Target, f.Reason))
	}
	want := []string{"3 ADR-0042 id-unlinked"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
}

// DC-FA-ID-001.a — Muster-Präzedenz: Deklarationsreihenfolge, das
// erste passende Muster gewinnt pro Vorkommen (überlappende Treffer).
func TestIDsMusterPraezedenz(t *testing.T) {
	lines := []Line{{No: 1, Text: "siehe ADR-0042 hier"}}
	long := IDPattern{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "x"}
	short := IDPattern{Regex: regexp.MustCompile(`ADR-\d{2}`), Target: "x"}

	got := checkIDs("f.md", nil, lines, []IDPattern{long, short})
	if len(got) != 1 || got[0].Target != "ADR-0042" {
		t.Fatalf("lang zuerst: %+v", got)
	}
	got = checkIDs("f.md", nil, lines, []IDPattern{short, long})
	if len(got) != 1 || got[0].Target != "ADR-00" {
		t.Fatalf("kurz zuerst: %+v", got)
	}
}

// „Verlinkt" = Vorkommen liegt im Linktext: mehrere Links pro Zeile,
// Vorkommen zwischen den Links bleibt ein Befund.
func TestIDsLinktextSpannen(t *testing.T) {
	lines := []Line{{No: 1, Text: "[ADR-0001](a.md) und ADR-0002 sowie [x ADR-0003 y](b.md)"}}
	p := []IDPattern{{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "x"}}
	got := checkIDs("f.md", nil, lines, p)
	if len(got) != 1 || got[0].Target != "ADR-0002" {
		t.Fatalf("Befunde = %+v (genau ADR-0002 erwartet)", got)
	}
}

// DC-FA-CONF-001: nicht existierendes ids-Target → Fehler (Exit 2),
// unabhängig davon, ob das Modul ids aktiv ist; Datei- und
// Verzeichnis-Targets (auch mit Slash) sind gültig.
func TestIDsTargetMussExistieren(t *testing.T) {
	m := newMemFS(map[string]string{"docs/a.md": "x"})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`X-\d`), Target: "gibt/es/nicht"},
	}}
	if _, err := Run(m, nil, cfg, []string{"ids"}); err == nil {
		t.Fatal("Fehler erwartet (Target fehlt)")
	}
	if _, err := Run(m, nil, cfg, []string{"links"}); err == nil {
		t.Fatal("Fehler auch ohne aktives ids-Modul erwartet (Config-Constraint)")
	}
	cfgOK := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`X-\d`), Target: "docs/a.md"},
		{Regex: regexp.MustCompile(`Y-\d`), Target: "docs/"},
		{Regex: regexp.MustCompile(`Z-\d`), Target: "/"}, // Repo-Wurzel
	}}
	if _, err := Run(m, nil, cfgOK, []string{"ids"}); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// fileInTarget: Datei-/Verzeichnis-/Wurzel-Targets und Escape-Fälle.
func TestFileInTarget(t *testing.T) {
	cases := []struct {
		file, target string
		want         bool
	}{
		{"docs/a.md", "docs/a.md", true},    // Target-Datei selbst
		{"docs/sub/b.md", "docs", true},     // unterhalb des Verzeichnisses
		{"docs2/a.md", "docs", false},       // Präfix-Falle docs2
		{"docs/a.md", "/", true},            // Repo-Wurzel: alles Definitions-Ort
		{"docs/a.md", "../draussen", false}, // Escape zählt nie als Target
	}
	for _, c := range cases {
		if got := fileInTarget(c.file, c.target); got != c.want {
			t.Errorf("fileInTarget(%q, %q) = %v, want %v", c.file, c.target, got, c.want)
		}
	}
}

// DC-FA-CONF-001: ids-Targets sind relativ zur Repo-Wurzel — Pfade,
// die die Wurzel verlassen (..), sind Konfigurationsfehler, auch wenn
// das Ziel außerhalb existieren würde.
func TestIDsTargetDarfWurzelNichtVerlassen(t *testing.T) {
	m := newMemFS(map[string]string{"docs/a.md": "x"})
	for _, target := range []string{"../draussen", "docs/../../x", ".."} {
		cfg := Config{IDPatterns: []IDPattern{
			{Regex: regexp.MustCompile(`X-\d`), Target: target},
		}}
		_, err := Run(m, nil, cfg, []string{"ids"})
		if err == nil || !strings.Contains(err.Error(), "verlässt") {
			t.Fatalf("Target %q: err = %v (Repo-Escape-Fehler erwartet)", target, err)
		}
	}
}

// Definitions-Ort und Headings sind linkpflichtfrei
// (Spez-Fortschreibung slice-007): Vorkommen im deklarierten Target
// des Musters sowie in ATX-Heading-Zeilen erzeugen keinen Befund.
func TestIDsDefinitionsOrtUndHeadings(t *testing.T) {
	m := newMemFS(map[string]string{
		"spec/lastenheft.md":      "### DC-FA-CLI-001 — Titel\nDC-FA-CLI-001 im eigenen Dokument\n",
		"docs/plan/adr/0001-x.md": "# ADR-0001 — X\nersetzt ADR-0042 nicht\n",
		"docs/a.md":               "## Abschnitt zu DC-FA-CLI-001\nnacktes DC-FA-CLI-001 hier\n",
	})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`DC-(FA-[A-Z]+|QA)-\d+`), Target: "spec/lastenheft.md"},
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"},
	}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].File != "docs/a.md" || res.Findings[0].Line != 2 {
		t.Fatalf("Befunde = %+v (genau docs/a.md:2 erwartet)", res.Findings)
	}
}

// Positionserhaltendes Inline-Code-Stripping: angrenzender Text darf
// nicht zu Phantom-Kennungen verschmelzen (Review R1 zu slice-006).
func TestIDsKeinePhantomKennungDurchInlineCode(t *testing.T) {
	content := []byte("AD`x`R-0042\n")
	lines := PreprocessMarkdown(content)
	p := []IDPattern{{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "x"}}
	if got := checkIDs("f.md", content, lines, p); len(got) != 0 {
		t.Fatalf("Phantom-Befund: %+v", got)
	}
}

// DC-FA-ID-001 (0.8.0) link-policy: always — Inline-Code-Vorkommen sind
// linkpflichtig; die Ventile (Linktext, target, exempt-paths,
// d-check:ignore, Fence, Heading) bleiben frei.
func TestIDsLinkPolicyAlways(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "`ADR-0042` als Code-Span ohne Link\n" + // L1: Befund
			"[`ADR-0043`](plan/adr/0043.md) verlinkt\n" + // L2: Linktext → frei
			"`ADR-0044` <!-- d-check:ignore (Beispiel) -->\n" + // L3: Marker → frei
			"## Heading mit `ADR-0047`\n" + // L4: Heading → frei
			"```\n`ADR-0045` im Fence\n```\n", // L5-7: Fence → frei
		"docs/reviews/r.md":      "`ADR-0046` literal im Review\n", // exempt-paths → frei
		"docs/plan/adr/0043.md":  "x",                              // target-Dir
		"docs/plan/adr/0099-x.md": "`ADR-0099` im eigenen target\n", // target → frei
	})
	cfg := Config{IDPatterns: []IDPattern{{
		Regex:       regexp.MustCompile(`ADR-\d{4}`),
		Target:      "docs/plan/adr/",
		LinkPolicy:  "always",
		ExemptPaths: []string{"docs/reviews/**"},
	}}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s %s", f.File, f.Line, f.Target, f.Reason))
	}
	want := []string{"docs/a.md:1 ADR-0042 id-unlinked"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("always-Befunde = %v, want %v", got, want)
	}
}

// Abwärtskompatibilität: ohne link-policy (Default prose) sind
// Code-Span-Vorkommen weiterhin linkpflichtfrei (DC-QA-02).
func TestIDsLinkPolicyProseDefault(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md":              "`ADR-0042` als Code-Span ohne Link\n", // prose → frei
		"docs/plan/adr/0001-x.md": "x",
	})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"}, // kein link-policy
	}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("prose-Default sollte Code-Spans frei lassen, got %+v", res.Findings)
	}
}

// DC-FA-ID-001 (0.13.0): die Ventile exempt-paths und d-check:ignore
// gelten auch für NACKTE Fließtext-Vorkommen (Ganzdatei-/Ganzzeilen-
// Carve-out), nicht nur für die always-Inline-Code-Vorkommen.
func TestIDsVentileNackteVorkommen(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/reviews/r.md": "Nackt im Review: ADR-0042\n", // exempt-paths → frei
		"docs/a.md":         "Nackt mit Marker: ADR-0043 <!-- d-check:ignore (Beispiel) -->\n", // ignore → frei
		"docs/b.md":         "Nackt ohne Schutz: ADR-0044\n", // Kontrolle → Befund
		"docs/plan/adr/0001-x.md": "x",
	})
	cfg := Config{IDPatterns: []IDPattern{{
		Regex:       regexp.MustCompile(`ADR-\d{4}`),
		Target:      "docs/plan/adr/",
		LinkPolicy:  "always",
		ExemptPaths: []string{"docs/reviews/**"},
	}}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s %s", f.File, f.Line, f.Target, f.Reason))
	}
	want := []string{"docs/b.md:1 ADR-0044 id-unlinked"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("Ventil-Befunde (nackt) = %v, want %v", got, want)
	}
}

// DC-FA-ID-001 (0.13.0): exempt-paths wirkt unabhängig von der
// link-policy — auch unter dem Default prose nimmt es nackte
// Vorkommen aus (Ganzdatei-Carve-out).
func TestIDsExemptPathsProseDefault(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/reviews/r.md":       "Nackt im Review: ADR-0042\n", // exempt → frei
		"docs/b.md":               "Nackt: ADR-0043\n",           // Kontrolle → Befund
		"docs/plan/adr/0001-x.md": "x",
	})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/",
			ExemptPaths: []string{"docs/reviews/**"}}, // kein link-policy (prose)
	}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	want := []string{"docs/b.md:1 id-unlinked"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("exempt-paths unter prose = %v, want %v", got, want)
	}
}

// DC-FA-ID-001 (0.13.0): der d-check:ignore-Marker nimmt eine NACKTE
// Prosa-ID auch unter der Default-Politik prose aus (politik-unabhängig,
// Gegenstück zu TestIDsExemptPathsProseDefault für das zweite Ventil).
func TestIDsIgnoreMarkerProseDefault(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md":               "Nackt mit Marker: ADR-0042 <!-- d-check:ignore (Beispiel) -->\n", // ignore → frei
		"docs/b.md":               "Nackt: ADR-0043\n",                                               // Kontrolle → Befund
		"docs/plan/adr/0001-x.md": "x",
	})
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"}, // kein link-policy (prose)
	}}
	res, err := Run(m, nil, cfg, []string{"ids"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	want := []string{"docs/b.md:1 id-unlinked"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("d-check:ignore unter prose = %v, want %v", got, want)
	}
}
