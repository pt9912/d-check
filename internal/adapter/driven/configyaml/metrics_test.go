package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

const mtYAML = "metrics:\n" +
	"  plans-dir: docs/plan/planning/in-progress\n" +
	"  block: groessenmarke\n" +
	"  fields: {path: pfad, base: basis, estimate: schaetzung, limit: marke, exception: ueberschreitung, cut-reason: schnitt-grund}\n" +
	"  surcharge: {numerator: 135, denominator: 10000}\n" +
	"  rounding: 100\n" +
	"  cut-threshold: 300000\n"

// Der Voll-Block wird gelesen (DC-FA-MET-001, DC-FA-CONF-001).
func TestDecode_MetricsVollblock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(mtYAML))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	m := cfg.Metrics
	if m.PlansDir != "docs/plan/planning/in-progress" || m.Block != "groessenmarke" {
		t.Fatalf("Verzeichnis oder Block nicht gelesen: %+v", m)
	}
	if m.SurchargeNumerator != 135 || m.SurchargeDenominator != 10000 || m.Rounding != 100 {
		t.Fatalf("Formel nicht gelesen: %+v", m)
	}
	if !m.CutEnabled || m.CutThreshold != 300000 {
		t.Fatalf("Schwelle nicht gelesen: %+v", m)
	}
	if m.Fields.Limit != "marke" || m.Fields.CutReason != "schnitt-grund" {
		t.Fatalf("Feldnamen nicht gelesen: %+v", m.Fields)
	}
}

// Ohne Schwelle ist die Schwellenregel aus; Feldnamen fallen auf die Defaults.
func TestDecode_MetricsOhneSchwelleMitDefaults(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("metrics:\n  plans-dir: p\n  block: b\n  surcharge: {numerator: 1, denominator: 1}\n  rounding: 1\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if cfg.Metrics.CutEnabled {
		t.Fatalf("ohne cut-threshold keine Schwellenregel: %+v", cfg.Metrics)
	}
	if f := cfg.Metrics.EffectiveFields(); f.Path != "path" || f.Limit != "limit" {
		t.Fatalf("Defaults falsch: %+v", f)
	}
}

// Ohne Block bleibt die Konfiguration leer.
func TestDecode_MetricsOhneBlock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("modules: [links]\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if cfg.Metrics.PlansDir != "" || cfg.Metrics.Block != "" {
		t.Fatalf("erwartet leer, got %+v", cfg.Metrics)
	}
}

// Config-Rand (Exit 2): Pflichtschlüssel, Wertebereiche, Feldnamen und
// unbekannte Schlüssel.
func TestDecode_MetricsRand(t *testing.T) {
	base := func(replace, with string) string { return strings.Replace(mtYAML, replace, with, 1) }
	cases := map[string]string{
		"plans-dir fehlt":                  base("  plans-dir: docs/plan/planning/in-progress\n", ""),
		"plans-dir leer":                   base("docs/plan/planning/in-progress", "' '"),
		"block fehlt":                      base("  block: groessenmarke\n", ""),
		"surcharge fehlt":                  base("  surcharge: {numerator: 135, denominator: 10000}\n", ""),
		"numerator fehlt":                  base("{numerator: 135, denominator: 10000}", "{denominator: 10000}"),
		"rounding fehlt":                   base("  rounding: 100\n", ""),
		"denominator null":                 base("denominator: 10000", "denominator: 0"),
		"rounding null":                    base("rounding: 100", "rounding: 0"),
		"numerator negativ":                base("numerator: 135", "numerator: -1"),
		"cut-threshold negativ":            base("cut-threshold: 300000", "cut-threshold: -5"),
		"Feldnamen gleich":                 base("limit: marke", "limit: pfad"),
		"unbekannter Schlüssel":            mtYAML + "  unbekannt: 1\n",
		"unbekannter Schlüssel in fields":  base("cut-reason: schnitt-grund", "cut-reasons: x"),
		"rounding keine Zahl":              base("rounding: 100", "rounding: hundert"),
		"excluded mit Glob":                mtYAML + "  excluded: [\"a/*\"]\n",
		"excluded absolut":                 mtYAML + "  excluded: [/a]\n",
		"excluded mit ..":                  mtYAML + "  excluded: [a/../b]\n",
		"excluded mit Punkt":               mtYAML + "  excluded: [./a]\n",
		"excluded leer":                    mtYAML + "  excluded: [\"\"]\n",
	}
	for name, doc := range cases {
		if _, err := configyaml.Decode([]byte(doc)); err == nil {
			t.Errorf("%s: erwartet Fehler (Exit 2)", name)
		}
	}
}

// Die ausgenommenen Pfade werden gelesen.
func TestDecode_MetricsExcluded(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(mtYAML + "  excluded: [spec/vertraege/, docs/x]\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if got := cfg.Metrics.Excluded; len(got) != 2 || got[0] != "spec/vertraege/" || got[1] != "docs/x" {
		t.Fatalf("excluded nicht gelesen: %v", got)
	}
}
