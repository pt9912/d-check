package configyaml_test

import (
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
// Fehler, ein leerer Wert gilt wie ein fehlender (DC-FA-LLEN-001).
func TestDecode_LinelengthOhneBlockUndOhneGrenze(t *testing.T) {
	for _, in := range []string{"modules: [links]\n", "linelength:\n  exempt-tables: false\n", "linelength:\n  max-chars:\n"} {
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
		"Grenze als Liste":        {"linelength:\n  max-chars: [80]\n", "linelength.max-chars"},
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
