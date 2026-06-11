package core

import (
	"fmt"
	"strings"
	"testing"
)

func matrixTestConfig() MatrixConfig {
	return MatrixConfig{
		Classes: []MatrixClass{
			{Name: "contract", Paths: []string{"spec/lastenheft.md"}},
			{Name: "adr", Paths: []string{"docs/plan/adr/[0-9]*.md"}},
			{Name: "slice", Paths: []string{"docs/plan/planning/**/slice-*.md"}},
		},
		Rules:           []MatrixRule{{From: "contract", To: "adr", Allow: false}},
		StatusForbidden: []string{"superseded", "deprecated"},
	}
}

// DC-FA-MTX-001 Happy/Boundary/Negative über Run: Slice → aktives ADR
// ok; Referenz auf superseded ADR → matrix-inactive; Lastenheft → ADR
// → matrix-forbidden mit beiden Klassen in der Meldung.
func TestMatrixModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"spec/lastenheft.md":      "# LH\n[verboten](../docs/plan/adr/0001-x.md)\n",
		"docs/plan/adr/0001-x.md": "# ADR-0001 — X\n\n**Status:** Accepted\n",
		"docs/plan/adr/0002-y.md": "# ADR-0002 — Y\n\n**Status:** Superseded by ADR-0042\n",
		"docs/plan/planning/done/slice-001-a.md": "# S\n[ok](../../adr/0001-x.md)\n[inaktiv](../../adr/0002-y.md)\n",
		"docs/notiz.md": "[unklassifiziert](plan/adr/0002-y.md)\n",
	})
	res, err := Run(m, Config{Matrix: matrixTestConfig()}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	want := []string{
		"docs/plan/planning/done/slice-001-a.md:3 matrix-inactive",
		"spec/lastenheft.md:2 matrix-forbidden",
	}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
	// Negative verlangt beide Klassen in der Meldung
	for _, f := range res.Findings {
		if f.Reason == ReasonMatrixForbidden &&
			(!strings.Contains(f.Message, "contract") || !strings.Contains(f.Message, "adr")) {
			t.Fatalf("Meldung ohne beide Klassen: %q", f.Message)
		}
	}
}

// exclude-sections: Links in ausgenommenen Sektionen (inkl.
// Untersektionen) werden nicht geprüft; nach Sektions-Ende greift die
// Prüfung wieder. Klassen-Präzedenz: erste deklarierte Klasse gewinnt.
func TestMatrixExcludeSectionsUndPraezedenz(t *testing.T) {
	cfg := matrixTestConfig()
	cfg.ExcludeSections = []string{"Historie"}
	m := newMemFS(map[string]string{
		"spec/lastenheft.md": "# LH\n" +
			"[vorher](../docs/plan/adr/0001-x.md)\n" + // Zeile 2: Befund
			"## Historie\n" +
			"[ausgenommen](../docs/plan/adr/0001-x.md)\n" + // Zeile 4: frei
			"### Untersektion\n" +
			"[auch frei](../docs/plan/adr/0001-x.md)\n" + // Zeile 6: frei
			"## Danach\n" +
			"[wieder geprüft](../docs/plan/adr/0001-x.md)\n", // Zeile 8: Befund
		"docs/plan/adr/0001-x.md": "**Status:** Accepted\n",
	})
	res, err := Run(m, Config{Matrix: cfg}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 2 || res.Findings[0].Line != 2 || res.Findings[1].Line != 8 {
		t.Fatalf("Befunde = %+v (Zeilen 2 und 8 erwartet)", res.Findings)
	}

	// Präzedenz: dieselbe Datei matcht zwei Klassen — die erste gewinnt
	classes := []MatrixClass{
		{Name: "speziell", Paths: []string{"docs/plan/adr/0001-x.md"}},
		{Name: "adr", Paths: []string{"docs/plan/adr/[0-9]*.md"}},
	}
	if name, _ := classOf(classes, "docs/plan/adr/0001-x.md"); name != "speziell" {
		t.Fatalf("Klassen-Präzedenz verletzt: %q", name)
	}
}

// DC-FA-MTX-001.a Schritt 2 — Status-Extraktion in fester Reihenfolge.
func TestStatusOf(t *testing.T) {
	cases := map[string]string{
		"**Status:** Accepted\n":                          "Accepted",
		"# T\n## Status\n\nProposed — Review offen\n":     "Proposed — Review offen",
		"## status\nDeprecated\n":                         "Deprecated", // Heading case-insensitiv
		"# T\nnur Text\n":                                 "",
		"```\n**Status:** im Fence\n```\n## Status\nAktiv\n": "Aktiv",
		"## Status\n\n## Weiter\nText\n":                  "", // leere Status-Sektion
		"## Status\nHeading-Form\n\n**Status:** Zeile gewinnt\n": "Zeile gewinnt",
	}
	for in, want := range cases {
		if got := statusOf([]byte(in)); got != want {
			t.Errorf("statusOf(%q) = %q, want %q", in, got, want)
		}
	}
	if !statusForbidden("Superseded by ADR-0042", []string{"superseded", "deprecated"}) {
		t.Error("Präfix-Match case-insensitiv erwartet")
	}
	if statusForbidden("Accepted", []string{"superseded"}) || statusForbidden("", []string{"superseded"}) {
		t.Error("aktive Dokumente dürfen nicht matchen")
	}
}
