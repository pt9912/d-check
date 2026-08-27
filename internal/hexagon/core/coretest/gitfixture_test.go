package coretest_test

import (
	"errors"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

// Negativ-Selbsttest des Fixture-Wächters: ohne ihn bliebe sein Ausfall still —
// der meldende Zweig wird von keinem anderen Test erreicht.
func TestGitFixtureRewriteHazard(t *testing.T) {
	for _, c := range []struct {
		name     string
		oldLen   int
		newLen   int
		erwartet bool
	}{
		{"keine Vorgänger-Datei", -1, 4, false},
		{"gleiche Länge", 4, 4, true},
		{"länger", 4, 9, false},
		{"kürzer", 9, 4, false},
		{"beide leer", 0, 0, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := coretest.GitFixtureRewriteHazard("f.md", c.oldLen, c.newLen)
			if got := errors.Is(err, coretest.ErrGleichLangeNeuschreibung); got != c.erwartet {
				t.Fatalf("Befund = %v, erwartet %v (err = %v)", got, c.erwartet, err)
			}
		})
	}
}
