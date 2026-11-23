package configyaml_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

// Der Block wird gelesen: Grenze, drei Schalter, Datei-Ventil und scope
// (DC-FA-LLEN-001, DC-FA-CONF-001).
func TestDecode_LinelengthVollblock(t *testing.T) {
	in := "linelength:\n" +
		"  scope: {roots: [docs], ignore: [\"docs/alt/**\"]}\n" +
		"  max-chars: 120\n" +
		"  exempt-code-blocks: false\n" +
		"  exempt-tables: true\n" +
		"  exempt-unbreakable: false\n" +
		"  exempt-paths: [CHANGELOG.md, \"docs/frei/**\"]\n"
	cfg, err := configyaml.Decode([]byte(in))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	ll := cfg.Linelength
	if ll.MaxChars == nil || *ll.MaxChars != 120 {
		t.Fatalf("max-chars nicht gelesen: %+v", ll)
	}
	if ll.CodeBlocksExempt() || !ll.TablesExempt() || ll.UnbreakableExempt() {
		t.Fatalf("Schalter nicht gelesen: %+v", ll)
	}
	if len(ll.ExemptPaths) != 2 || ll.ExemptPaths[1] != "docs/frei/**" {
		t.Fatalf("exempt-paths nicht gelesen: %+v", ll.ExemptPaths)
	}
	sc := cfg.Scopes["linelength"]
	if sc == nil || len(sc.Roots) != 1 || sc.Roots[0] != "docs" || len(sc.Ignore) != 1 {
		t.Fatalf("scope nicht gelesen: %+v", sc)
	}
}

// Ohne Block bleibt die Konfiguration leer, die Schalter stehen im Default
// „ausgenommen“; ein Block ohne max-chars ist erst im aktiven Modul ein
// Fehler; ein fehlender Schalter steht im Default (DC-FA-LLEN-001).
func TestDecode_LinelengthOhneBlockUndOhneGrenze(t *testing.T) {
	for _, in := range []string{"modules: [links]\n", "linelength:\n  exempt-tables: false\n"} {
		cfg, err := configyaml.Decode([]byte(in))
		if err != nil || cfg.Linelength.MaxChars != nil {
			t.Fatalf("%q: err=%v cfg=%+v", in, err, cfg.Linelength)
		}
		if !cfg.Linelength.CodeBlocksExempt() || !cfg.Linelength.UnbreakableExempt() {
			t.Fatalf("%q: Default der Schalter verloren: %+v", in, cfg.Linelength)
		}
	}
}

// Jede Fehlkonfiguration des Blocks endet am Rand mit einem Fehler (Exit 2),
// der Schlüssel und Ursache nennt — auch ohne aktives Modul (DC-FA-LLEN-001).
func TestDecode_LinelengthEingabefehler(t *testing.T) {
	cases := map[string]struct{ in, want string }{
		"Grenze 0":                {"linelength:\n  max-chars: 0\n", "linelength.max-chars 0 muss >= 1"},
		"Grenze negativ":          {"linelength:\n  max-chars: -3\n", "linelength.max-chars -3 muss >= 1"},
		"Grenze als Zeichenkette": {"linelength:\n  max-chars: \"80\"\n", "linelength.max-chars \"80\" ist keine ganze Zahl"},
		"Grenze als Kommazahl":    {"linelength:\n  max-chars: 80.5\n", "linelength.max-chars \"80.5\" ist keine ganze Zahl"},
		"Grenze als Liste":        {"linelength:\n  max-chars: [80]\n", "linelength.max-chars ist eine Liste, keine ganze Zahl"},
		"Grenze als Abbildung":    {"linelength:\n  max-chars: {n: 80}\n", "linelength.max-chars ist eine Abbildung, keine ganze Zahl"},
		"Grenze 80.0":             {"linelength:\n  max-chars: 80.0\n", "linelength.max-chars \"80.0\" ist keine ganze Zahl"},
		"Schalter nein":           {"linelength:\n  max-chars: 80\n  exempt-tables: \"nein\"\n", "linelength.exempt-tables \"nein\" ist kein Wahrheitswert"},
		"Schalter yes":            {"linelength:\n  max-chars: 80\n  exempt-code-blocks: yes\n", "linelength.exempt-code-blocks \"yes\" ist kein Wahrheitswert"},
		"Schalter als Zahl":       {"linelength:\n  max-chars: 80\n  exempt-unbreakable: 1\n", "linelength.exempt-unbreakable \"1\" ist kein Wahrheitswert"},
		"unbekannter Schlüssel":   {"linelength:\n  max-chars: 80\n  max-lines: 3\n", "max-lines"},
		"leeres Glob":             {"linelength:\n  max-chars: 80\n  exempt-paths: [\"\"]\n", "linelength.exempt-paths enthält ein leeres Glob"},
		"ungültiges Segment":      {"linelength:\n  max-chars: 80\n  exempt-paths: [\"docs/[a\"]\n", "Segment \"[a\""},
		"scope ohne roots":        {"linelength:\n  max-chars: 80\n  scope: {ignore: [x]}\n", "linelength.scope.roots fehlt"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(c.in))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, erwartet %q", err, c.want)
			}
		})
	}
}

// Einer der vier eigenen Werte (max-chars, drei Schalter) ohne Wert (YAML-null)
// ist gesetzt und ungültig: Exit 2 am Rand, auch ohne aktives Modul; abwesend
// ist nur ein fehlender Schlüssel (DC-FA-LLEN-001).
func TestDecode_LinelengthSchluesselOhneWert(t *testing.T) {
	cases := map[string]struct{ in, key string }{
		"Schalter leer":      {"linelength:\n  max-chars: 80\n  exempt-code-blocks:\n", "exempt-code-blocks"},
		"Schalter Tilde":     {"linelength:\n  max-chars: 80\n  exempt-tables: ~\n", "exempt-tables"},
		"Schalter null":      {"linelength:\n  max-chars: 80\n  exempt-unbreakable: null\n", "exempt-unbreakable"},
		"Grenze leer":        {"modules: [links]\nlinelength:\n  max-chars:\n", "max-chars"},
		"Grenze leer allein": {"linelength:\n  max-chars:\n  exempt-tables: true\n", "max-chars"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(c.in))
			if err == nil || !strings.Contains(err.Error(), "linelength."+c.key+" steht ohne Wert (YAML-null)") {
				t.Fatalf("err = %v, erwartet Exit 2 für linelength.%s ohne Wert", err, c.key)
			}
		})
	}
}

// Gegenprobe: bei versions (pin-pattern: neben patterns ist keine Mischform),
// planning (marker: ergibt den Default), structure (max-lines:) und ids
// (scope: ist die globale Scan-Menge) bleibt ein Schlüssel ohne Wert
// abwesend — die strenge Lesart gilt nur für die vier eigenen Werte von
// linelength (DC-FA-LLEN-001, DC-FA-VER-001, DC-FA-PLAN-001, DC-FA-STRUCT-001,
// DC-FA-CONF-002).
func TestDecode_SchluesselOhneWertBeiNachbarnAbwesend(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("versions:\n  pin-pattern:\n  patterns:\n    - {pin-pattern: 'x:(v\\d+\\.\\d+\\.\\d+)'}\n"))
	if err != nil || len(cfg.Versions.Patterns) != 1 {
		t.Fatalf("versions: err=%v cfg=%+v", err, cfg.Versions)
	}
	cfg, err = configyaml.Decode([]byte("planning:\n  roadmap: r.md\n  marker:\n"))
	if err != nil || cfg.Planning.Marker != "" || cfg.Planning.EffectiveMarker() != "Keine aktive Welle" {
		t.Fatalf("planning: err=%v cfg=%+v", err, cfg.Planning)
	}
	cfg, err = configyaml.Decode([]byte("structure:\n  - files: 'a/*.md'\n    section: '## H'\n    non-empty: true\n    max-lines:\n"))
	if err != nil || len(cfg.Structure) != 1 || cfg.Structure[0].MaxLines != nil {
		t.Fatalf("structure: err=%v cfg=%+v", err, cfg.Structure)
	}
	cfg, err = configyaml.Decode([]byte("ids:\n  scope:\n"))
	if err != nil || cfg.Scopes["ids"] != nil {
		t.Fatalf("ids: err=%v scopes=%+v", err, cfg.Scopes)
	}
}

// scope und exempt-paths ohne Wert gelten wie bei jedem scannenden Modul als
// abwesend: dieselbe Konfiguration wie ohne den Schlüssel, also die globale
// Scan-Menge und keine Ausnahme (DC-FA-LLEN-001, DC-FA-CONF-002).
func TestDecode_LinelengthScopeUndExemptPathsOhneWert(t *testing.T) {
	ohne, err := configyaml.Decode([]byte("linelength:\n  max-chars: 80\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []string{"linelength:\n  max-chars: 80\n  scope:\n", "linelength:\n  max-chars: 80\n  exempt-paths:\n"} {
		mit, err := configyaml.Decode([]byte(in))
		if err != nil || !reflect.DeepEqual(ohne, mit) {
			t.Fatalf("%q: err=%v\nohne %+v\nmit  %+v", in, err, ohne, mit)
		}
		if mit.Scopes["linelength"] != nil || len(mit.Linelength.ExemptPaths) != 0 {
			t.Fatalf("%q: scope oder Ausnahme gesetzt: %+v", in, mit)
		}
	}
}

// Ein YAML-Alias steht für den Wert seines Ankers und wird an ihm geprüft:
// eine ganze Zahl und ein Wahrheitswert sind gültig, eine Liste meldet ihre
// Art, ein Anker ohne Wert ist ein Schlüssel ohne Wert (DC-FA-LLEN-001).
func TestDecode_LinelengthAlias(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("file:\n  - files: '*.md'\n    max-lines: &g 20\nlinelength:\n  exempt-code-blocks: &b false\n  exempt-tables: *b\n  max-chars: *g\n"))
	if err != nil || cfg.Linelength.MaxChars == nil || *cfg.Linelength.MaxChars != 20 || cfg.Linelength.TablesExempt() {
		t.Fatalf("Alias auf Zahl und Wahrheitswert: err=%v cfg=%+v", err, cfg.Linelength)
	}
	cases := map[string]struct{ in, want string }{
		"Liste": {"scan:\n  roots: &r [docs]\nlinelength:\n  max-chars: *r\n", "linelength.max-chars ist eine Liste, keine ganze Zahl"},
		"null":  {"scan:\n  ignore: &n\nlinelength:\n  max-chars: 80\n  exempt-tables: *n\n", "linelength.exempt-tables steht ohne Wert (YAML-null)"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := configyaml.Decode([]byte(c.in)); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, erwartet %q", err, c.want)
			}
		})
	}
}
