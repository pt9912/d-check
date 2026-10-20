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
		{"leere liste", "trace:\n  requirements: []\n", "darf als Liste nicht leer sein"},
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

// Aliase und Merge-Keys bleiben in der Listen-Form wie in der Objekt-Form
// auflösbar: zwei Einträge mit gleicher source teilen ihren Rahmen über einen
// Anker.
func TestDecode_RequirementsListAliasAndMerge(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(`trace:
  requirements:
    - &fa
      source: spec/lastenheft.md
      id-pattern: 'SFC-FA-\d{2}'
    - <<: *fa
      id-pattern: 'SFC-QA-\d{2}'
      format: table
      table:
        id-column: ID
        text-column: Anforderung
    - *fa
`))
	if err != nil {
		t.Fatalf("Anker und Merge-Key abgelehnt: %v", err)
	}
	s := cfg.Trace.Sources
	if len(s) != 3 {
		t.Fatalf("Sources = %d Einträge, want 3", len(s))
	}
	if s[1].Source != "spec/lastenheft.md" || !s[1].ReqPattern.MatchString("SFC-QA-01") || s[1].Format != model.TraceFormatTable {
		t.Fatalf("Merge-Key nicht aufgelöst: %+v", s[1])
	}
	if s[2].Source != "spec/lastenheft.md" || !s[2].ReqPattern.MatchString("SFC-FA-01") {
		t.Fatalf("Alias nicht aufgelöst: %+v", s[2])
	}
}

// Ein unbekannter Schlüssel hinter einem Anker oder Merge-Key bleibt ein Fehler,
// mit der Zeile seiner Definition.
func TestDecode_RequirementsListUnknownKeyBehindAlias(t *testing.T) {
	for name, doc := range map[string]string{
		"alias":     "trace:\n  requirements:\n    - &a\n      sorce: x.md\n    - *a\n",
		"merge-key": "trace:\n  requirements:\n    - &a\n      sorce: x.md\n    - <<: *a\n      source: y.md\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(doc))
			if err == nil || !strings.Contains(err.Error(), "line 4: field sorce not found") {
				t.Fatalf("Decode-Fehler = %v, want line 4", err)
			}
		})
	}
}

// Die Zeile bleibt die des Dokuments, auch mit Kommentaren und Leerzeilen im
// Block — in der Objekt-Form wie in der Liste.
func TestDecode_RequirementsErrorLineWithCommentsAndBlankLines(t *testing.T) {
	for name, c := range map[string]struct{ doc, want string }{
		"objekt": {"trace:\n  requirements:\n    # Kommentar\n\n    source: a.md\n    # noch einer\n    sorce: c.md\n", "line 7: field sorce not found"},
		"liste":  {"trace:\n  requirements:\n    # Kommentar\n    - source: a.md\n\n    - source: b.md\n      # zwischen\n      sorce: c.md\n", "line 8: field sorce not found"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(c.doc))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("Decode-Fehler = %v, want enthält %q", err, c.want)
			}
		})
	}
}
