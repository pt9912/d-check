package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// listYAML ist die Listen-Form aus DC-FA-REQ-001: Überschriften für FA, eine
// Tabelle für QA, beide aus derselben Datei.
const listYAML = `trace:
  requirements:
    - source: spec/lastenheft.md
      id-pattern: 'SFC-FA-\d{2}'
      format: headings
    - source: spec/lastenheft.md
      id-pattern: 'SFC-QA-\d{2}'
      format: table
      table:
        id-column: ID
        text-column: Anforderung
        modality-column: Priorität
`

// Happy: die Liste wird je Eintrag übernommen; die Felder der Objekt-Form
// bleiben unbelegt.
func TestDecode_RequirementsListHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(listYAML))
	if err != nil {
		t.Fatalf("gültige Liste abgelehnt: %v", err)
	}
	tr := cfg.Trace
	if len(tr.Sources) != 2 {
		t.Fatalf("Sources = %d Einträge, want 2", len(tr.Sources))
	}
	if tr.Source != "" || tr.ReqPattern != nil || tr.Format != "" || tr.Table != nil || tr.Modality != nil {
		t.Fatalf("Listen-Form belegt die Felder der Objekt-Form: %+v", tr)
	}
	fa, qa := tr.Sources[0], tr.Sources[1]
	if fa.Format != model.TraceFormatHeadings || !fa.ReqPattern.MatchString("SFC-FA-01") || fa.ReqPattern.MatchString("SFC-QA-01") {
		t.Fatalf("Eintrag 0 falsch: %+v", fa)
	}
	if qa.Format != model.TraceFormatTable || qa.Table == nil || qa.Table.IDColumn != "ID" || qa.Table.ModalityColumn != "Priorität" {
		t.Fatalf("Eintrag 1 falsch: %+v", qa)
	}
}

// Die Objekt-Form bleibt gültig und belegt Sources nicht (DC-QA-02).
func TestDecode_RequirementsObjectFormUnchanged(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("trace:\n  requirements:\n    source: spec/lastenheft.md\n    id-pattern: 'R-\\d+'\n"))
	if err != nil {
		t.Fatalf("Objekt-Form abgelehnt: %v", err)
	}
	if cfg.Trace.Sources != nil || cfg.Trace.Source != "spec/lastenheft.md" || cfg.Trace.ReqPattern == nil {
		t.Fatalf("Objekt-Form falsch übernommen: %+v", cfg.Trace)
	}
}

// Negative: jede Form, die weder Objekt noch nichtleere Liste von Objekten ist,
// und jeder Fehler eines Eintrags endet fail-closed; der Fehler nennt den Index.
func TestDecode_RequirementsListErrors(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want string
	}{
		{"leere liste (AK-7)", "trace:\n  requirements: []\n", "darf als Liste nicht leer sein"},
		{"null-eintrag", "trace:\n  requirements:\n    - source: a.md\n    -\n", "trace.requirements[1] ist leer"},
		{"skalar statt objekt", "trace:\n  requirements: foo\n", "cannot unmarshal"},
		{"skalar als eintrag", "trace:\n  requirements:\n    - foo\n", "cannot unmarshal"},
		{"unbekannter schlüssel im eintrag", "trace:\n  requirements:\n    - source: a.md\n    - sorce: b.md\n", "field sorce not found"},
		{"unbekannter schlüssel in der objekt-form", "trace:\n  requirements:\n    sorce: b.md\n", "field sorce not found"},
		{"unbekannter schlüssel im table-block", "trace:\n  requirements:\n    - format: table\n      table:\n        id-column: ID\n        text-colum: X\n", "field text-colum not found"},
		{"regex im zweiten eintrag", "trace:\n  requirements:\n    - id-pattern: 'A-\\d'\n    - id-pattern: '('\n", "trace.requirements[1].id-pattern"},
		{"format im zweiten eintrag", "trace:\n  requirements:\n    - format: headings\n    - format: csv\n", "trace.requirements[1].format"},
		{"table ohne table-format im ersten eintrag", "trace:\n  requirements:\n    - table:\n        id-column: ID\n        text-column: X\n", "trace.requirements[0].table"},
		{"pfad im zweiten eintrag", "trace:\n  requirements:\n    - source: a.md\n    - source: ../b.md\n", "trace.requirements[1].source"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(tc.yaml))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Decode-Fehler = %v, want enthält %q", err, tc.want)
			}
		})
	}
}

// Die Zeilennummer eines unbekannten Schlüssels zeigt auf das Dokument, nicht
// auf den Teilbaum des Eintrags.
func TestDecode_RequirementsListErrorLine(t *testing.T) {
	yamlDoc := "scan:\n  roots: [docs]\ntrace:\n  requirements:\n    - source: a.md\n    - source: b.md\n      sorce: c.md\n"
	_, err := configyaml.Decode([]byte(yamlDoc))
	if err == nil || !strings.Contains(err.Error(), "line 7: field sorce not found") {
		t.Fatalf("Decode-Fehler = %v, want line 7", err)
	}
	_, err = configyaml.Decode([]byte("trace:\n  requirements:\n    source: a.md\n    sorce: c.md\n"))
	if err == nil || !strings.Contains(err.Error(), "line 4: field sorce not found") {
		t.Fatalf("Objekt-Form: Decode-Fehler = %v, want line 4", err)
	}
}
