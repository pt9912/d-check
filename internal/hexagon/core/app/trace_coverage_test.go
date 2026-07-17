package app

import (
	"regexp"
	"testing"
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
		{"Klammer in der URL bricht das Suffix", "`](a(1).md)..003", nil},
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
