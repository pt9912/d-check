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

	got := checkIDs("f.md", lines, []IDPattern{long, short})
	if len(got) != 1 || got[0].Target != "ADR-0042" {
		t.Fatalf("lang zuerst: %+v", got)
	}
	got = checkIDs("f.md", lines, []IDPattern{short, long})
	if len(got) != 1 || got[0].Target != "ADR-00" {
		t.Fatalf("kurz zuerst: %+v", got)
	}
}

// „Verlinkt" = Vorkommen liegt im Linktext: mehrere Links pro Zeile,
// Vorkommen zwischen den Links bleibt ein Befund.
func TestIDsLinktextSpannen(t *testing.T) {
	lines := []Line{{No: 1, Text: "[ADR-0001](a.md) und ADR-0002 sowie [x ADR-0003 y](b.md)"}}
	p := []IDPattern{{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "x"}}
	got := checkIDs("f.md", lines, p)
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
	}}
	if _, err := Run(m, nil, cfgOK, []string{"ids"}); err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
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
	lines := PreprocessMarkdown([]byte("AD`x`R-0042\n"))
	p := []IDPattern{{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "x"}}
	if got := checkIDs("f.md", lines, p); len(got) != 0 {
		t.Fatalf("Phantom-Befund: %+v", got)
	}
}
