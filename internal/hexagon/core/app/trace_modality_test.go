package app

import (
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// slice-068 (DC-FA-MOD-001): der Klassifikator liest die Stufe aus dem ersten/
// längsten Keyword-Treffer im normalisierten Body; Negation (MUSS NICHT ≠ DARF
// NICHT), Wortgrenze (musste ≠ MUSS), Emphasis/Umbruch-Normalisierung, unknown.
func TestModalityClassify(t *testing.T) {
	mm := resolveModality(&model.TraceModality{}) // leer ⇒ Built-in-Defaults
	tests := []struct {
		name string
		body string
		want string
	}{
		{"must", "Das System MUSS verfügbar sein.", "must"},
		{"should", "Die Anbindung SOLLTE dokumentiert sein.", "should"},
		{"may", "Die Plattform KANN RL unterstuetzen.", "may"},
		{"darf-nicht = must", "Der MVP DARF NICHT externe Dienste verlangen.", "must"},
		{"muss-nicht = may (längster Treffer)", "Das Feld MUSS NICHT gesetzt sein.", "may"},
		{"emphasis normalisiert", "Das Feld **MUSS** NICHT gesetzt sein.", "may"},
		{"umbruch normalisiert", "Das Feld MUSS\nNICHT gesetzt sein.", "may"},
		{"wortgrenze: musste ≠ MUSS", "Das System musste getestet werden.", "unknown"},
		{"kein modal-verb ⇒ unknown", "Die Plattform ist KEIN vollstaendiges SCADA-System.", "unknown"},
		{"englisch", "The system MUST NOT allow it.", "must"},
		{"früheste position gewinnt", "Es SOLLTE geprüft werden; später MUSS es gelten.", "should"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mm.classify(tc.body); got != tc.want {
				t.Fatalf("classify(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}

// Override der Keywords + require-levels-Auflösung.
func TestModalityResolveOverride(t *testing.T) {
	mm := resolveModality(&model.TraceModality{
		Levels:        map[string][]string{"pflicht": {"ZWINGEND"}, "kür": {"OPTIONAL"}},
		RequireLevels: []string{"pflicht"},
	})
	if got := mm.classify("Das ist ZWINGEND."); got != "pflicht" {
		t.Fatalf("Override-Klassifikation: %q", got)
	}
	if got := mm.classify("Das ist OPTIONAL."); got != "kür" {
		t.Fatalf("Override-Klassifikation: %q", got)
	}
	if !mm.requireLevels["pflicht"] || mm.requireLevels["kür"] {
		t.Fatalf("require-levels falsch: %+v", mm.requireLevels)
	}
	// Default require-levels = [must].
	def := resolveModality(&model.TraceModality{})
	if !def.requireLevels["must"] || def.requireLevels["may"] {
		t.Fatalf("Default require-levels falsch: %+v", def.requireLevels)
	}
}

// normalizeBody: Emphasis raus, Whitespace/Umbruch zu einem Leerzeichen.
func TestNormalizeBody(t *testing.T) {
	if got := normalizeBody("**MUSS**\n  NICHT `x`"); got != "MUSS NICHT x" {
		t.Fatalf("normalizeBody = %q", got)
	}
}
