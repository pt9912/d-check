package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// --manual gibt die Abschnitte beider Dokumente aus, deren Überschrift den
// Begriff nennt — repo-frei, auch hinter dem Pfad-Argument des Images
// (DC-FA-CLI-013).
func TestManual_GibtAbschnitteAus(t *testing.T) {
	var so, se bytes.Buffer
	code := cli.Run([]string{"/repo-gibt-es-nicht", "--manual", "planning"}, &so, &se)
	out := so.String()
	if code != 0 || se.Len() != 0 {
		t.Fatalf("Exit %d, stderr %q", code, se.String())
	}
	for _, want := range []string{"==> docs/user/benutzerhandbuch.md:", "==> spec/spezifikation.md:", "planning"} {
		if !strings.Contains(out, want) {
			t.Errorf("Ausgabe ohne %q", want)
		}
	}
}

// Kein Treffer, leerer Begriff und eine Kombination mit einem anderen Modus
// sind Nutzungsfehler (DC-FA-CLI-013).
func TestManual_Nutzungsfehler(t *testing.T) {
	for name, args := range map[string][]string{
		"kein Treffer": {"--manual", "zzz-kein-ueberschriftsbegriff-zzz"},
		"leer":         {"--manual", "  "},
		"mit --json":   {"--manual", "reviews", "--json"},
		"mit --trace":  {"--manual", "reviews", "--trace"},
	} {
		var so, se bytes.Buffer
		if code := cli.Run(args, &so, &se); code != 2 || so.Len() != 0 {
			t.Errorf("%s: Exit %d, stdout %q — Exit 2 ohne Ausgabe erwartet", name, code, so.String())
		}
	}
	var so, se bytes.Buffer
	cli.Run([]string{"--manual", "zzz-kein-ueberschriftsbegriff-zzz"}, &so, &se)
	if !strings.Contains(se.String(), "Benutzerhandbuch") {
		t.Errorf("der Hinweis nennt die Dokument-Titel, got %q", se.String())
	}
}
