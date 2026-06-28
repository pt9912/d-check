package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

func matrixTestConfig() model.MatrixConfig {
	return model.MatrixConfig{
		Classes: []model.MatrixClass{
			{Name: "contract", Paths: []string{"spec/lastenheft.md"}},
			{Name: "adr", Paths: []string{"docs/plan/adr/[0-9]*.md"}},
			{Name: "slice", Paths: []string{"docs/plan/planning/**/slice-*.md"}},
		},
		Rules:           []model.MatrixRule{{From: "contract", To: "adr", Allow: false}},
		StatusForbidden: []string{"superseded", "deprecated"},
	}
}

// DC-FA-MTX-001 Happy/Boundary/Negative über Run: Slice → aktives ADR
// ok; Referenz auf superseded ADR → matrix-inactive; Lastenheft → ADR
// → matrix-forbidden mit beiden Klassen in der Meldung.
func TestMatrixModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md": "# LH\n[verboten](../docs/plan/adr/0001-x.md)\n" +
			"[doppelt](../docs/plan/adr/0002-y.md)\n",
		"docs/plan/adr/0001-x.md": "# ADR-0001 — X\n\n**Status:** Accepted\n",
		"docs/plan/adr/0002-y.md": "# ADR-0002 — Y\n\n**Status:** Superseded by ADR-0042\n",
		"docs/plan/planning/done/slice-001-a.md": "# S\n[ok](../../adr/0001-x.md)\n[inaktiv](../../adr/0002-y.md)\n",
		"docs/notiz.md": "[unklassifiziert](plan/adr/0002-y.md)\n",
	})
	res, err := Run(m, nil, model.Config{Matrix: matrixTestConfig()}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	// Zeile 3 des Lastenhefts verletzt BEIDE Bedingungen — Regel- und
	// Status-Prüfung sind unabhängig, zwei Befunde
	// (spec/spezifikation.md §DC-FA-MTX-001.a).
	want := []string{
		"docs/plan/planning/done/slice-001-a.md:3 matrix-inactive",
		"spec/lastenheft.md:2 matrix-forbidden",
		"spec/lastenheft.md:3 matrix-forbidden",
		"spec/lastenheft.md:3 matrix-inactive",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
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
	m := coretest.NewMemFS(map[string]string{
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
	res, err := Run(m, nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 2 || res.Findings[0].Line != 2 || res.Findings[1].Line != 8 {
		t.Fatalf("Befunde = %+v (Zeilen 2 und 8 erwartet)", res.Findings)
	}

	// Präzedenz: dieselbe Datei matcht zwei Klassen — die erste gewinnt
	classes := []model.MatrixClass{
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
		// Fence-Inhalt in der Status-Sektion ist kein Statuswert
		"## Status\n```\nSuperseded\n```\nAktiv\n": "Aktiv",
		"## Status\n```\nnur Fence\n```\n":         "",
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

// Status wird nur aus Markdown-Zielen extrahiert — Nicht-Markdown wird
// weder gelesen noch gecacht (kein Voll-Read von Binärdateien).
func TestCachedStatusNurMarkdown(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{"docs/bild.png": "**Status:** Superseded\n"})
	cache := map[string]*string{}
	if got := cachedStatus(m, cache, "docs/bild.png"); got != "" {
		t.Fatalf("Status = %q, want \"\" (Nicht-Markdown)", got)
	}
	if len(cache) != 0 {
		t.Fatal("Nicht-Markdown-Ziel darf nicht gelesen/gecacht werden")
	}
}

// Cache-Treffer: vorhandene Einträge (Wert und nil-Sentinel für
// unlesbare Ziele) werden ohne erneuten Read beantwortet.
func TestCachedStatusCacheTreffer(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{})
	wert := "Accepted"
	cache := map[string]*string{
		"docs/a.md": &wert,
		"docs/b.md": nil, // unlesbar → matrix schweigt
	}
	if got := cachedStatus(m, cache, "docs/a.md"); got != "Accepted" {
		t.Fatalf("Cache-Wert = %q", got)
	}
	if got := cachedStatus(m, cache, "docs/b.md"); got != "" {
		t.Fatalf("nil-Sentinel = %q, want \"\"", got)
	}
	// unlesbares Ziel (fehlt im memFS) wird als nil gecacht
	if got := cachedStatus(m, cache, "docs/fehlt.md"); got != "" {
		t.Fatalf("unlesbares Ziel = %q, want \"\"", got)
	}
	if s, ok := cache["docs/fehlt.md"]; !ok || s != nil {
		t.Fatal("unlesbares Ziel muss als nil gecacht werden")
	}
}

// DC-FA-MTX-001 Supersede-Lineage (0.14.0): opt-in nimmt die deklarierte
// Lineage-Kante von matrix-inactive aus (Happy); fremde Quellen bleiben
// inaktiv (Boundary); Default aus ⇒ Befundsatz byte-identisch (Negative).
func TestMatrixSupersedeLineage(t *testing.T) {
	files := map[string]string{
		// X löst Y ab (Feld in fetter Header-Form) und verweist darauf.
		"docs/plan/adr/0006-x.md": "# ADR-0006 — X\n\n" +
			"**Status:** Accepted\n" +
			"**Aenderungstyp:** Supersedes ADR 0003\n\n" +
			"**Bezug:** [ADR 0003](0003-lifecycle.md)\n",
		// Y ist abgelöst (inaktiv).
		"docs/plan/adr/0003-lifecycle.md": "# ADR-0003 — Lifecycle\n\n**Status:** Superseded by ADR-0006\n",
		// Fremde Quelle ohne Supersede-Feld verweist ebenfalls auf Y.
		"docs/plan/planning/done/slice-001-a.md": "# S\n[alt](../../adr/0003-lifecycle.md)\n",
	}

	// Happy + Boundary — Lineage AN: nur die fremde Quelle bleibt inaktiv.
	cfg := matrixTestConfig()
	cfg.AllowSupersedeLineage = true
	cfg.SupersedeFields = []string{"Supersedes", "Aenderungstyp"}
	res, err := Run(coretest.NewMemFS(files), nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	if got := inactiveFiles(res.Findings); len(got) != 1 || got[0] != "docs/plan/planning/done/slice-001-a.md" {
		t.Fatalf("Lineage AN: matrix-inactive = %v, want nur die fremde Quelle", got)
	}

	// Negative (Default aus): die Lineage-Kante erzeugt ebenfalls inactive.
	resOff, err := Run(coretest.NewMemFS(files), nil, model.Config{Matrix: matrixTestConfig()}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	if got := inactiveFiles(resOff.Findings); len(got) != 2 {
		t.Fatalf("Lineage AUS: matrix-inactive = %v, want 2 (byte-identisch zum Alt-Verhalten)", got)
	}
}

func inactiveFiles(fs []model.Finding) []string {
	var out []string
	for _, f := range fs {
		if f.Reason == ReasonMatrixInactive {
			out = append(out, f.File)
		}
	}
	return out
}

// Feld-Extraktion (fette + Frontmatter-Form, case-insensitiv) und die
// Match-Formen des Lineage-Vergleichs (Linktext, Basename, leere Werte).
func TestSupersedeFieldValueUndLineage(t *testing.T) {
	fields := []string{"Supersedes"}
	if v, ok := supersedeFieldValue("Supersedes: ADR 0003", fields); !ok || v != "ADR 0003" {
		t.Fatalf("plain-Form = (%q,%v)", v, ok)
	}
	if v, ok := supersedeFieldValue("**supersedes:** 0003-x.md", fields); !ok || v != "0003-x.md" {
		t.Fatalf("fette Form case-insensitiv = (%q,%v)", v, ok)
	}
	if _, ok := supersedeFieldValue("**Bezug:** etwas", fields); ok {
		t.Fatal("Nicht-Feld-Zeile darf nicht matchen")
	}
	if lineageExempt(nil, "ADR 0003", "docs/x.md") {
		t.Fatal("leere values dürfen nicht ausnehmen")
	}
	if !lineageExempt([]string{"loest 0003-lifecycle ab"}, "", "docs/plan/adr/0003-lifecycle.md") {
		t.Fatal("Basename-Match (ohne Linktext) erwartet")
	}
	if lineageExempt([]string{"verweist auf etwas anderes"}, "ADR 0099", "docs/plan/adr/0099-y.md") {
		t.Fatal("ohne Nennung des Ziels keine Ausnahme")
	}
}

// Eine explizit erlaubte Regel (allow: true) erzeugt keinen Befund.
func TestMatrixErlaubteRegel(t *testing.T) {
	cfg := matrixTestConfig()
	cfg.Rules = []model.MatrixRule{{From: "slice", To: "adr", Allow: true}}
	m := coretest.NewMemFS(map[string]string{
		"docs/plan/planning/done/slice-001-a.md": "[ok](../../adr/0001-x.md)\n",
		"docs/plan/adr/0001-x.md":                "**Status:** Accepted\n",
	})
	res, err := Run(m, nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("res = %+v, err = %v (erlaubte Regel)", res, err)
	}
}

// DC-FA-MTX-002 Happy/Boundary: klasseninterne Richtung über order/direction.
// Aufwärts (höherer Rang → kleinerer Index) ist ok; abwärts — auch transitiv —
// erzeugt matrix-downward; ein rangfreies Klassen-Mitglied nimmt nicht teil.
func TestMatrixDownwardRichtung(t *testing.T) {
	cfg := model.MatrixConfig{
		Classes: []model.MatrixClass{{
			Name: "spec-straten",
			Paths: []string{
				"spec/lastenheft.md", "spec/spezifikation.md",
				"spec/architecture.md", "spec/notiz.md",
			},
			Order: []string{
				"spec/lastenheft.md", "spec/spezifikation.md", "spec/architecture.md",
			},
			Direction: model.DirectionNoDownward,
		}},
	}
	m := coretest.NewMemFS(map[string]string{
		// architecture (Rang 2) → lastenheft (Rang 0): aufwärts, kein Befund
		"spec/architecture.md": "# A\n[hoch](lastenheft.md)\n",
		// spezifikation (Rang 1) → architecture (Rang 2): abwärts, Befund Zeile 2
		"spec/spezifikation.md": "# S\n[runter](architecture.md)\n",
		// lastenheft (Rang 0) → architecture (Rang 2): abwärts transitiv, Befund Zeile 2
		"spec/lastenheft.md": "# L\n[transitiv](architecture.md)\n",
		// notiz.md ist Klassen-Mitglied, aber rangfrei (kein order-Treffer) → kein Befund
		"spec/notiz.md": "# N\n[egal](architecture.md)\n",
	})
	res, err := Run(m, nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	want := []string{
		"spec/lastenheft.md:2 matrix-downward",
		"spec/spezifikation.md:2 matrix-downward",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
	// Boundary-Meldung nennt beide Ränge.
	for _, f := range res.Findings {
		if !strings.Contains(f.Message, "Rang") {
			t.Fatalf("matrix-downward-Meldung ohne Ränge: %q", f.Message)
		}
	}
}

// DC-FA-MTX-002 Default-aus: ohne order/direction ist der Befundsatz
// byte-identisch (keine matrix-downward-Befunde, DC-QA-02).
func TestMatrixDownwardDefaultAus(t *testing.T) {
	cfg := model.MatrixConfig{
		Classes: []model.MatrixClass{
			{Name: "spec-straten", Paths: []string{"spec/lastenheft.md", "spec/architecture.md"}},
		},
	}
	m := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md":   "# L\n[runter](architecture.md)\n",
		"spec/architecture.md": "# A\n",
	})
	res, err := Run(m, nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("res = %+v, err = %v (ohne order/direction kein Befund)", res, err)
	}
}

// DC-FA-MTX-002 Kanten (Review MEDIUM-1): rangfreies Ziel, Gleichrang und
// Selbstverweis lösen kein matrix-downward aus.
func TestMatrixDownwardKanten(t *testing.T) {
	noFinding := func(t *testing.T, name string, cfg model.MatrixConfig, files map[string]string) {
		t.Helper()
		res, err := Run(coretest.NewMemFS(files), nil, model.Config{Matrix: cfg}, []string{"matrix"})
		if err != nil || len(res.Findings) != 0 {
			t.Fatalf("%s: res=%+v err=%v (kein Befund erwartet)", name, res, err)
		}
	}
	// (a) rangfreies Ziel: Quelle rangbehaftet, Ziel Klassen-Mitglied ohne order-Treffer.
	noFinding(t, "rangfreies Ziel",
		model.MatrixConfig{Classes: []model.MatrixClass{{
			Name: "c", Paths: []string{"spec/lastenheft.md", "spec/notiz.md"},
			Order: []string{"spec/lastenheft.md"}, Direction: model.DirectionNoDownward,
		}}},
		map[string]string{"spec/lastenheft.md": "# L\n[zu rangfrei](notiz.md)\n", "spec/notiz.md": "# N\n"})
	// (b) Gleichrang: ein order-Glob matcht beide Dateien → gleicher Rang.
	noFinding(t, "Gleichrang",
		model.MatrixConfig{Classes: []model.MatrixClass{{
			Name: "c", Paths: []string{"spec/a.md", "spec/b.md"},
			Order: []string{"spec/*.md"}, Direction: model.DirectionNoDownward,
		}}},
		map[string]string{"spec/a.md": "# A\n[gleichrang](b.md)\n", "spec/b.md": "# B\n"})
	// (c) Selbstverweis: Datei verweist auf sich selbst (si == di).
	noFinding(t, "Selbstverweis",
		model.MatrixConfig{Classes: []model.MatrixClass{{
			Name: "c", Paths: []string{"spec/lastenheft.md", "spec/architecture.md"},
			Order: []string{"spec/lastenheft.md", "spec/architecture.md"}, Direction: model.DirectionNoDownward,
		}}},
		map[string]string{"spec/lastenheft.md": "# L\n[selbst](lastenheft.md)\n", "spec/architecture.md": "# A\n"})
}

// DC-FA-MTX-003: token-basierte Referenz-Richtung. Unmarkierter Slice-Token im
// ADR-Körper → matrix-forbidden; mit Provenance-Marker auf der Zeile → kein
// Befund; in exempt-paths-Datei → grandfathered; Token nur im Markdown-Link →
// kein Doppelbefund (Link-Pfad); Selbst-Klasse zählt nicht.
func TestMatrixTokenReferenz(t *testing.T) {
	cfg := model.MatrixConfig{
		Classes: []model.MatrixClass{
			{Name: "adr", Paths: []string{"docs/plan/adr/[0-9]*.md"}},
			{Name: "slice", Paths: []string{"docs/plan/planning/**/slice-*.md"},
				Token: regexp.MustCompile(`slice-\d{3}`)},
		},
		Rules:           []model.MatrixRule{{From: "adr", To: "slice", Allow: false}},
		ExemptPaths:     []string{"docs/plan/adr/0001-*.md"},
		ExcludeSections: []string{"Geschichte"},
	}
	m := coretest.NewMemFS(map[string]string{
		"docs/plan/adr/0050-x.md": "# A\nEntsteht mit slice-042 als Grundlage.\n",                  // unmarkiert → Befund Z.2
		"docs/plan/adr/0051-y.md": "# A\nVerifiziert in slice-043. <!-- d-check:status-provenance -->\n", // markiert → frei
		"docs/plan/adr/0001-z.md": "# A\nWegen slice-044 entschieden.\n",                            // grandfathered → frei
		"docs/plan/adr/0052-l.md": "# A\nSiehe [slice-045](../planning/done/slice-045-x.md).\n",     // nur Link → einmal (Link), kein Token-Doppel
		"docs/plan/adr/0053-f.md": "# A\n```\nslice-077\n```\n",                                     // Token in Fenced-Code → zählt nicht
		"docs/plan/adr/0054-h.md": "# A\n## Geschichte\nWegen slice-088 (Provenance).\n",            // Token unter exclude-sections → zählt nicht
		"docs/plan/adr/0055-m.md": "# A\nFolgt slice-061 und slice-062.\n",                          // zwei Token in einer Zeile → zwei Befunde (FindAll)
		"docs/plan/planning/done/slice-045-x.md": "# S\nFolgt slice-046.\n",                         // slice→slice: kein Befund (self-class)
	})
	res, err := Run(m, nil, model.Config{Matrix: cfg}, []string{"matrix"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Reason))
	}
	want := []string{
		"docs/plan/adr/0050-x.md:2 matrix-forbidden",
		"docs/plan/adr/0052-l.md:2 matrix-forbidden",
		"docs/plan/adr/0055-m.md:2 matrix-forbidden", // slice-061
		"docs/plan/adr/0055-m.md:2 matrix-forbidden", // slice-062
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
	// Token-Befund nennt beide Klassen.
	for _, f := range res.Findings {
		if f.File == "docs/plan/adr/0050-x.md" && (!strings.Contains(f.Message, "adr") || !strings.Contains(f.Message, "slice")) {
			t.Fatalf("Token-Meldung ohne beide Klassen: %q", f.Message)
		}
	}
}
