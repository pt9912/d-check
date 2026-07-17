package app

import (
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// ggPat spiegelt eine Fremd-Konvention (grid-gym): GG-<FAMILIE>-NNN.
var ggPat = regexp.MustCompile(`GG-[A-Z][A-Z0-9]*-\d{3}`)

// slice-067 (DC-FA-COV-001): der Range-Parser expandiert `..`-Ranges
// breiten-erhaltend und inklusiv, `/`-Aufzählungen, verwirft Nicht-id-pattern-
// Treffer und ist fail-closed bei AAA>BBB / Breiten-Mismatch.
func TestExpandRange(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		rest    string
		want    []string
		wantErr bool
	}{
		{"inklusiv breiten-erhaltend", "GG-QA-001", "..006", []string{"GG-QA-001", "GG-QA-002", "GG-QA-003", "GG-QA-004", "GG-QA-005", "GG-QA-006"}, false},
		{"einzelner Wert (AAA==BBB)", "GG-QA-003", "..003", []string{"GG-QA-003"}, false},
		{"aufzählung /", "GG-RT-004", "/005", []string{"GG-RT-005"}, false},
		{"aufzählung mehrfach", "GG-RT-004", "/005/006", []string{"GG-RT-005", "GG-RT-006"}, false},
		{"familie mit ziffer (DNP3)", "GG-DNP3-001", "..003", []string{"GG-DNP3-001", "GG-DNP3-002", "GG-DNP3-003"}, false},
		{"kein suffix", "GG-QA-001", " und Text", nil, false},
		{"AAA>BBB fail-closed", "GG-RT-009", "..003", nil, true},
		{"breiten-mismatch fail-closed", "GG-QA-001", "..0010", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandRange("trace.coverage", tc.id, tc.rest, ggPat)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("erwartete Fehler, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if !equalStrings(got, tc.want) {
				t.Fatalf("expandRange(%q,%q) = %v, want %v", tc.id, tc.rest, got, tc.want)
			}
		})
	}
}

// rangeAwareIDs: exakte Treffer plus (bei ranges) die Expansion; range-Fehler
// propagiert.
func TestRangeAwareIDs(t *testing.T) {
	text := "Abdeckung: GG-QA-001..003 und GG-RT-004/005.\n"
	got, err := rangeAwareIDs("trace.coverage", text, ggPat, true)
	if err != nil {
		t.Fatalf("Fehler: %v", err)
	}
	// exakt: GG-QA-001, GG-RT-004; expandiert: GG-QA-001..003, GG-RT-005.
	for _, want := range []string{"GG-QA-001", "GG-QA-002", "GG-QA-003", "GG-RT-004", "GG-RT-005"} {
		if !contains(got, want) {
			t.Fatalf("rangeAwareIDs ohne %q: %v", want, got)
		}
	}
	// ranges:false ⇒ nur exakte Treffer, keine Expansion.
	got2, _ := rangeAwareIDs("trace.coverage", text, ggPat, false)
	if contains(got2, "GG-QA-003") {
		t.Fatalf("ranges:false expandierte trotzdem: %v", got2)
	}
	// ungültige Range propagiert als Fehler.
	if _, err := rangeAwareIDs("trace.coverage", "GG-RT-009..003", ggPat, true); err == nil {
		t.Fatal("erwartete Fehler bei AAA>BBB")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// slice-073 (ADR-0039): Link-Transparenz. Steht die Kennung unter Linkpflicht
// (DC-FA-ID-001), folgt ihre Fortsetzung nicht mehr unmittelbar — genau EIN
// Link-Suffix darf dazwischen stehen. Ohne diese Regel bricht die
// unqualifizierte Range-Zusage des Lastenhefts dort, wo d-checks eigene
// Linkpflicht greift.
func TestExpandRangeLinkTransparent(t *testing.T) {
	tests := []struct {
		name string
		rest string
		want []string
	}{
		// Positiv: die belegte Realform (Code-Span im Linktext) und ohne Span.
		{"range hinter Link mit Code-Span", "`](../spec/x.md)..003",
			[]string{"GG-QA-001", "GG-QA-002", "GG-QA-003"}},
		{"range hinter Link ohne Code-Span", "](../spec/x.md)..003",
			[]string{"GG-QA-001", "GG-QA-002", "GG-QA-003"}},
		{"enum hinter Link", "`](x.md)/002/003", []string{"GG-QA-002", "GG-QA-003"}},
		{"unverlinkt bleibt unverändert", "..003",
			[]string{"GG-QA-001", "GG-QA-002", "GG-QA-003"}},

		// Negativ — gegen das Raten (ADR-0039: genau eines, sonst nichts).
		{"zwei Link-Suffixe", "`](a.md)](b.md)..003", nil},
		{"Zeichen zwischen ) und ..", "`](a.md)x..003", nil},
		{"Whitespace zwischen ) und ..", "`](a.md) ..003", nil},
		{"Link ohne Fortsetzung", "`](a.md)", nil},
		// Ein Klammer-Ziel ist ein LEGITIMER Link (der kanonische Reader grenzt
		// balanciert ab) — die Fortsetzung dahinter wird gelesen. Die erste Fassung
		// erwartete hier nil und schrieb damit den Defekt fest (R1-F-1/F-2).
		{"balanciertes Klammer-Ziel, danach Range", "`](a(1).md)..003",
			[]string{"GG-QA-001", "GG-QA-002", "GG-QA-003"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandRange("trace.coverage", "GG-QA-001", tc.rest, ggPat)
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if !equalStrings(got, tc.want) {
				t.Fatalf("expandRange(GG-QA-001, %q) = %v, want %v", tc.rest, got, tc.want)
			}
		})
	}
}

// slice-073: die Fail-closed-Fälle gelten auch hinter einem Link — sonst
// entkäme eine ungültige Range der Prüfung, sobald sie verlinkt ist.
func TestExpandRangeLinkTransparentFailClosed(t *testing.T) {
	if _, err := expandRange("trace.coverage", "GG-RT-009", "`](x.md)..003", ggPat); err == nil {
		t.Fatal("AAA>BBB hinter einem Link muss fail-closed bleiben")
	}
	if _, err := expandRange("trace.coverage", "GG-QA-001", "`](x.md)..0010", ggPat); err == nil {
		t.Fatal("Breiten-Mismatch hinter einem Link muss fail-closed bleiben")
	}
}

// slice-073 R1-F-1 (HIGH): Klammern im Linkziel. Die erste Fassung grenzte das
// Ziel per Regex an der ERSTEN `)` ab — eine zweite Link-Definition neben dem
// kanonischen Reader. Folge: der URL-REST landete im Range-Parser, und eine Zelle
// GANZ OHNE Range-Notation expandierte (`/002/003` aus dem Pfad) und versteckte
// damit Waisen. Ein stiller Falschbefund, gefährlicher als der behobene laute.
func TestExpandRangeBalancedLinkTarget(t *testing.T) {
	tests := []struct {
		name string
		rest string
		want []string
	}{
		// KEINE Range-Notation in der Zelle — es darf nichts expandieren.
		{"Pfadsegmente hinter Klammer-URL sind kein Enum", "`](../specs/Rev(2)/002/003.md)", nil},
		{"Range-artiger Text in der URL", "`](../a(1)..003.md)", nil},
		// Balanciertes Ziel: das Suffix endet korrekt, die Fortsetzung wird gelesen.
		{"Klammer-URL, danach echte Range", "`](https://x.org/A_(b))..003",
			[]string{"GG-QA-001", "GG-QA-002", "GG-QA-003"}},
		{"Klammer-URL, danach echtes Enum", "`](https://x.org/A_(b))/002/003",
			[]string{"GG-QA-002", "GG-QA-003"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandRange("trace.coverage", "GG-QA-001", tc.rest, ggPat)
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if !equalStrings(got, tc.want) {
				t.Fatalf("expandRange(GG-QA-001, %q) = %v, want %v", tc.rest, got, tc.want)
			}
		})
	}
}

// slice-073 R1-F-4: die DoD hakte Tests auf KONSUMENTEN-Ebene ab, die nicht
// existierten — alle Tests riefen expandRange direkt. Das hier ist die
// Fitness-Aussage von ADR-0039, an ihrer echten Wirkungsstelle gemessen:
// verlinkt und unverlinkt decken identisch ab.
func TestCoverageRefsLinkTransparent(t *testing.T) {
	reqs := func(body string) map[string]string {
		return map[string]string{"docs/traceability.md": "# T\n\n" + body + "\n"}
	}
	src := []model.TraceCoverage{{Files: []string{"docs/traceability.md"}, Label: "Trace", Ranges: true}}
	for _, tc := range []struct {
		name string
		body string
		want int
	}{
		{"unverlinkt", "Abdeckung: GG-QA-001..003", 3},
		{"verlinkt", "Abdeckung: [`GG-QA-001`](../spec/x.md)..003", 3},
		{"verlinkt mit Klammer-URL, ohne Range", "Abdeckung: [`GG-QA-001`](../a(2)/002/003.md)", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := coverageRefs(coretest.NewMemFS(reqs(tc.body)), src, ggPat)
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if len(got) != tc.want {
				t.Fatalf("abgedeckte Anforderungen = %d, want %d (%v)", len(got), tc.want, got)
			}
		})
	}
}

// slice-074 (ADR-0040): der Splitter ist ein REINER Splitter — er kennt den
// Unterschied zwischen Header und Datenzeile nicht und darf deshalb nichts
// entfernen. Die erste Fassung strippte hier und wandte damit eine Body-Regel auf
// den Header an; die Tabelle fiel unter die Trennzeilen-Breite und wurde wortlos
// übersprungen (Review R1-F-1). Dieser Test pinnt die Nicht-Zuständigkeit.
func TestSplitPipeTableLineLaesstKommentarZelleStehen(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{"Kommentar am Ende ist eine Zelle", "| a | b | <!-- d-check:ignore (x) -->",
			[]string{"a", "b", "<!-- d-check:ignore (x) -->"}},
		{"auch mit Schluss-Pipe", "| a | b | <!-- d-check:ignore (x) --> |",
			[]string{"a", "b", "<!-- d-check:ignore (x) -->"}},
		{"Kommentar in einer Zelle bleibt Inhalt", "| a <!-- x --> | b |", []string{"a <!-- x -->", "b"}},
		{"ohne Kommentar unverändert", "| a | b |", []string{"a", "b"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := splitPipeTableLine(tc.line)
			if !ok {
				t.Fatalf("Zeile nicht als Tabellenzeile erkannt: %q", tc.line)
			}
			if !equalStrings(got, tc.want) {
				t.Fatalf("splitPipeTableLine(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}

// slice-074 (ADR-0040): die Toleranz sitzt dort, wo der Header-Kontext bekannt ist.
// Genau eine überzählige Zelle, und nur wenn sie ganzzellig ein Kommentar ist.
func TestCellCountOK(t *testing.T) {
	hdr := []string{"Kennung", "Anforderung"}
	tests := []struct {
		name  string
		cells []string
		want  bool
	}{
		{"exakt", []string{"F-1", "Text"}, true},
		{"eine Kommentar-Zelle zu viel — toleriert", []string{"F-1", "Text", "<!-- d-check:ignore (x) -->"}, true},
		{"eine ECHTE Zelle zu viel — fail-closed", []string{"F-1", "Text", "x"}, false},
		{"zwei zu viel, letzte Kommentar — fail-closed", []string{"F-1", "Text", "x", "<!-- y -->"}, false},
		{"eine zu wenig — fail-closed", []string{"F-1"}, false},
		{"zwei Kommentare in einer Zelle — fail-closed", []string{"F-1", "Text", "<!-- x --> <!-- y -->"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := cellCountOK(tc.cells, hdr); got != tc.want {
				t.Fatalf("cellCountOK(%v) = %v, want %v", tc.cells, got, tc.want)
			}
		})
	}
}
