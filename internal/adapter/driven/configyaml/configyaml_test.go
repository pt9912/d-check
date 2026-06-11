package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

func TestDecode_Vollbeispiel(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(`
scan:
  roots: [docs, spec]
  ignore: ["docs/archive/**"]
modules: [links, anchors, ids]
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
matrix:
  classes:
    - name: contract
      paths: [spec/lastenheft.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
  rules:
    - {from: contract, to: adr, allow: false}
  status:
    forbidden: [superseded, deprecated]
  exclude-sections: [Historie]
external:
  timeout-seconds: 10
  parallel: 4
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Roots) != 2 || len(cfg.IDPatterns) != 1 || cfg.IDPatterns[0].Target != "docs/plan/adr/" {
		t.Fatalf("cfg = %+v", cfg)
	}
	// matrix-Konfiguration wird in den Kern durchgereicht (DC-FA-MTX-001)
	m := cfg.Matrix
	if len(m.Classes) != 2 || m.Classes[0].Name != "contract" ||
		len(m.Rules) != 1 || m.Rules[0].From != "contract" || m.Rules[0].Allow ||
		len(m.StatusForbidden) != 2 || len(m.ExcludeSections) != 1 {
		t.Fatalf("Matrix = %+v", m)
	}
}

// matrix.status fehlt → Default [superseded, deprecated]
// (spec/spezifikation.md §2).
func TestDecode_MatrixStatusDefault(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x.md]\n"))
	if err != nil {
		t.Fatal(err)
	}
	got := cfg.Matrix.StatusForbidden
	if len(got) != 2 || got[0] != "superseded" || got[1] != "deprecated" {
		t.Fatalf("StatusForbidden = %v", got)
	}
}

func TestDecode_UnbekannterSchluessel(t *testing.T) {
	_, err := configyaml.Decode([]byte("modul: x\n"))
	if err == nil || !strings.Contains(err.Error(), "line") {
		t.Fatalf("err = %v (Zeilenangabe erwartet)", err)
	}
	if strings.Contains(err.Error(), "configyaml.raw") {
		t.Fatalf("interner Typname leakt: %v", err)
	}
}

func TestDecode_LeereUndKommentarDatei(t *testing.T) {
	for _, content := range []string{"", "# nur kommentar\n", "\n\n"} {
		cfg, err := configyaml.Decode([]byte(content))
		if err != nil || cfg.Roots != nil || cfg.Modules != nil {
			t.Fatalf("configyaml.Decode(%q): cfg=%+v err=%v (Defaults erwartet)", content, cfg, err)
		}
	}
}

func TestDecode_UngueltigerRegex(t *testing.T) {
	_, err := configyaml.Decode([]byte("ids:\n  patterns:\n    - regex: '[unclosed'\n      target: docs/\n"))
	if err == nil || !strings.Contains(err.Error(), "regex") {
		t.Fatalf("err = %v", err)
	}
}

// ids-Regexe, die den Leerstring matchen (X*, (ADR)?), würden an
// jeder Position Leer-Matches und damit unbrauchbare Befunde mit
// leerem Target erzeugen — Konfigurationsfehler statt Befund-Flut.
func TestDecode_LeerstringRegex(t *testing.T) {
	for _, re := range []string{"X*", "(ADR)?", ""} {
		_, err := configyaml.Decode([]byte("ids:\n  patterns:\n    - regex: '" + re + "'\n      target: docs/\n"))
		if err == nil || !strings.Contains(err.Error(), "Leerstring") {
			t.Fatalf("regex %q: err = %v (Leerstring-Ablehnung erwartet)", re, err)
		}
	}
}

func TestDecode_MatrixUndExternalConstraints(t *testing.T) {
	if _, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x]\n  rules:\n    - {from: a, to: fehlt, allow: false}\n")); err == nil {
		t.Fatal("undeklarierte Klasse: Fehler erwartet")
	}
	if _, err := configyaml.Decode([]byte("external:\n  timeout-seconds: 999\n")); err == nil {
		t.Fatal("timeout außerhalb 1–300: Fehler erwartet")
	}
	// explizit gesetzte 0 ist KEIN stiller Default, sondern Fehler
	// (Review R1 zu slice-008; Constraint 1–300 bzw. 1–16)
	if _, err := configyaml.Decode([]byte("external:\n  timeout-seconds: 0\n")); err == nil {
		t.Fatal("timeout-seconds: 0 muss Konfigurationsfehler sein")
	}
	if _, err := configyaml.Decode([]byte("external:\n  parallel: 0\n")); err == nil {
		t.Fatal("parallel: 0 muss Konfigurationsfehler sein")
	}
}

func TestDecode_OhneDatei(t *testing.T) {
	cfg, err := configyaml.Decode(nil)
	if err != nil || cfg.Roots != nil || cfg.Modules != nil {
		t.Fatalf("Defaults erwartet: %+v (%v)", cfg, err)
	}
}
