package configyaml

import (
	"strings"
	"testing"
)

func TestDecode_Vollbeispiel(t *testing.T) {
	cfg, err := Decode([]byte(`
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
}

func TestDecode_UnbekannterSchluessel(t *testing.T) {
	_, err := Decode([]byte("modul: x\n"))
	if err == nil || !strings.Contains(err.Error(), "line") {
		t.Fatalf("err = %v (Zeilenangabe erwartet)", err)
	}
}

func TestDecode_UngueltigerRegex(t *testing.T) {
	_, err := Decode([]byte("ids:\n  patterns:\n    - regex: '[unclosed'\n      target: docs/\n"))
	if err == nil || !strings.Contains(err.Error(), "regex") {
		t.Fatalf("err = %v", err)
	}
}

func TestDecode_MatrixUndExternalConstraints(t *testing.T) {
	if _, err := Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x]\n  rules:\n    - {from: a, to: fehlt, allow: false}\n")); err == nil {
		t.Fatal("undeklarierte Klasse: Fehler erwartet")
	}
	if _, err := Decode([]byte("external:\n  timeout-seconds: 999\n")); err == nil {
		t.Fatal("timeout außerhalb 1–300: Fehler erwartet")
	}
}

func TestDecode_OhneDatei(t *testing.T) {
	cfg, err := Decode(nil)
	if err != nil || cfg.Roots != nil || cfg.Modules != nil {
		t.Fatalf("Defaults erwartet: %+v (%v)", cfg, err)
	}
}
