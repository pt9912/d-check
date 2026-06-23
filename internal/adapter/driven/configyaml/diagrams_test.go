package configyaml_test

import (
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

// DC-FA-DIAG-001: diagrams-Block mit fences + Mustern (regex + defined-in).
func TestDecodeDiagrams(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(`
diagrams:
  fences: [mermaid]
  patterns:
    - regex: 'ARC-\d{2}'
      defined-in: spec/arch.md
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Diagrams.Patterns) != 1 ||
		cfg.Diagrams.Patterns[0].DefinedIn != "spec/arch.md" ||
		len(cfg.Diagrams.Fences) != 1 || cfg.Diagrams.Fences[0] != "mermaid" {
		t.Fatalf("diagrams falsch dekodiert: %+v", cfg.Diagrams)
	}
}

// DC-FA-CONF-002: diagrams.scope ersetzt den globalen Scan-Scope.
func TestDecodeDiagramsScope(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(`
diagrams:
  scope:
    roots: [spec]
  patterns:
    - regex: 'ARC-\d{2}'
      defined-in: spec/arch.md
`))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Scopes["diagrams"] == nil || len(cfg.Scopes["diagrams"].Roots) != 1 {
		t.Fatalf("diagrams.scope nicht übernommen: %+v", cfg.Scopes)
	}
}

// DC-FA-DIAG-001: statische Validierungsfehler (Exit 2 ohne Prüfung).
func TestDecodeDiagramsErrors(t *testing.T) {
	cases := map[string]string{
		"bad-regex": `
diagrams:
  patterns:
    - regex: '['
      defined-in: x
`,
		"empty-match": `
diagrams:
  patterns:
    - regex: 'x*'
      defined-in: x
`,
		"missing-defined-in": `
diagrams:
  patterns:
    - regex: 'ARC-\d'
`,
		"empty-fence": `
diagrams:
  fences: ['']
`,
	}
	for name, y := range cases {
		if _, err := configyaml.Decode([]byte(y)); err == nil {
			t.Fatalf("%s: erwartete einen Fehler", name)
		}
	}
}
