package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// crossYAML ist eine vollständige, gültige trace.cross-consistency-Config; die
// Tests hängen gezielte Abweichungen an bzw. lassen Felder weg.
const crossYAML = `trace:
  cross-consistency:
    forward:
      file: docs/traceability.md
      req-column: Anforderung
      design-column: Design-Artefakte
      design-pattern: 'GG-AR-[A-Z0-9-]+'
    backward:
      file: spec/architecture.md
      edge-column: Bezug
      req-pattern: 'GG-[A-Z]+-\d{3}'
`

// TestDecode_CrossConsistencyHappy (DC-FA-XREF-001): eine gültige Config wird
// übernommen, die Muster kompiliert und die Defaults aufgelöst — mode `equal`,
// ranges `true`, artifact-id-column der positionelle Sentinel `first`.
func TestDecode_CrossConsistencyHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(crossYAML))
	if err != nil {
		t.Fatalf("gültige cross-consistency-Config abgelehnt: %v", err)
	}
	cc := cfg.Trace.CrossConsistency
	if cc == nil {
		t.Fatal("trace.cross-consistency nicht übernommen")
	}
	if cc.Mode != model.TraceCrossModeEqual {
		t.Fatalf("mode-Default = %q, want %q", cc.Mode, model.TraceCrossModeEqual)
	}
	if !cc.Forward.Ranges || !cc.Backward.Ranges {
		t.Fatalf("ranges-Default nicht true: forward=%v backward=%v", cc.Forward.Ranges, cc.Backward.Ranges)
	}
	if cc.Backward.ArtifactIDColumn != model.TraceCrossArtifactFirst {
		t.Fatalf("artifact-id-column-Default = %q, want %q", cc.Backward.ArtifactIDColumn, model.TraceCrossArtifactFirst)
	}
	if cc.ExcludeReq != nil {
		t.Fatalf("exclude-req ohne Angabe gesetzt: %v", cc.ExcludeReq)
	}
	if cc.Forward.DesignPattern == nil || !cc.Forward.DesignPattern.MatchString("GG-AR-COMP-CORE") {
		t.Fatalf("design-pattern nicht kompiliert/wirksam: %v", cc.Forward.DesignPattern)
	}
	if cc.Backward.ReqPattern == nil || !cc.Backward.ReqPattern.MatchString("GG-ARCH-006") {
		t.Fatalf("req-pattern nicht kompiliert/wirksam: %v", cc.Backward.ReqPattern)
	}
}

// TestDecode_CrossConsistencyOhneBlock (DC-QA-02): ohne den Block bleibt das Feld
// nil ⇒ kein Abgleich, RTM byte-identisch.
func TestDecode_CrossConsistencyOhneBlock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("trace:\n  requirements:\n    source: spec/lastenheft.md\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if cfg.Trace.CrossConsistency != nil {
		t.Fatalf("cross-consistency ohne Block gesetzt: %+v", cfg.Trace.CrossConsistency)
	}
}

// TestDecode_CrossConsistencyExplizit: mode/ranges/exclude-req/artifact-id-column
// sind explizit setzbar — insbesondere `ranges: false` gegen den true-Default.
func TestDecode_CrossConsistencyExplizit(t *testing.T) {
	in := crossYAML + "    mode: superset\n    exclude-req: '^GG-SPEC-'\n"
	in = strings.Replace(in, "      design-pattern: 'GG-AR-[A-Z0-9-]+'\n",
		"      design-pattern: 'GG-AR-[A-Z0-9-]+'\n      ranges: false\n", 1)
	in = strings.Replace(in, "      edge-column: Bezug\n",
		"      edge-column: Bezug\n      artifact-id-column: Kennung\n", 1)
	cfg, err := configyaml.Decode([]byte(in))
	if err != nil {
		t.Fatalf("gültige Config abgelehnt: %v", err)
	}
	cc := cfg.Trace.CrossConsistency
	if cc.Mode != model.TraceCrossModeSuperset {
		t.Fatalf("mode = %q, want superset", cc.Mode)
	}
	if cc.Forward.Ranges {
		t.Fatal("explizites ranges: false wurde nicht übernommen")
	}
	if cc.Backward.ArtifactIDColumn != "Kennung" {
		t.Fatalf("artifact-id-column = %q, want Kennung", cc.Backward.ArtifactIDColumn)
	}
	if cc.ExcludeReq == nil || !cc.ExcludeReq.MatchString("GG-SPEC-042") {
		t.Fatalf("exclude-req nicht kompiliert/wirksam: %v", cc.ExcludeReq)
	}
}

// TestDecode_CrossConsistencyFehler (DC-FA-XREF-001, fail-closed): jeder
// Config-Defekt ist Exit 2 statt eines stillen Nicht-Abgleichs. Ohne diese
// Negativtests machte ein Refactor die Validierung lautlos nachgiebig.
func TestDecode_CrossConsistencyFehler(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "forward fehlt",
			in: "trace:\n  cross-consistency:\n    backward:\n      file: spec/architecture.md\n" +
				"      edge-column: Bezug\n      req-pattern: 'GG-[A-Z]+-\\d{3}'\n",
			want: "braucht forward und backward",
		},
		{
			name: "backward fehlt",
			in: "trace:\n  cross-consistency:\n    forward:\n      file: docs/traceability.md\n" +
				"      req-column: Anforderung\n      design-column: Design-Artefakte\n      design-pattern: 'GG-AR-.+'\n",
			want: "braucht forward und backward",
		},
		{name: "unbekannter mode", in: crossYAML + "    mode: gleich\n", want: "mode"},
		{
			name: "ungültiges design-pattern",
			in:   strings.Replace(crossYAML, "'GG-AR-[A-Z0-9-]+'", "'[unclosed'", 1),
			want: "design-pattern",
		},
		{
			name: "ungültiges req-pattern",
			in:   strings.Replace(crossYAML, `'GG-[A-Z]+-\d{3}'`, "'(unclosed'", 1),
			want: "req-pattern",
		},
		{name: "ungültiges exclude-req", in: crossYAML + "    exclude-req: '(unclosed'\n", want: "exclude-req"},
		{
			name: "leere req-column",
			in:   strings.Replace(crossYAML, "      req-column: Anforderung\n", "      req-column: ''\n", 1),
			want: "req-column ist leer",
		},
		{
			name: "leere edge-column",
			in:   strings.Replace(crossYAML, "      edge-column: Bezug\n", "      edge-column: ''\n", 1),
			want: "edge-column ist leer",
		},
		{
			name: "forward.file außerhalb der Wurzel",
			in:   strings.Replace(crossYAML, "      file: docs/traceability.md\n", "      file: ../fremd/t.md\n", 1),
			want: "muss relativ zur Repo-Wurzel liegen",
		},
		{
			name: "unbekannter Schlüssel im Block",
			in:   crossYAML + "    generator: true\n",
			want: "field generator not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(tc.in))
			if err == nil {
				t.Fatal("Konfigurationsfehler erwartet, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Fehlertext %q enthält %q nicht", err.Error(), tc.want)
			}
		})
	}
}
