package rules

import (
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Dieser Test ist die Verkoerperung von BEO-003: eine geteilte Lexik driftet an
// den Raendern, weil jeder Konsument sie selbst vorbereitet. Er prueft nicht,
// dass eine bestimmte Funktion aufgerufen wird, sondern dass alle Konsumenten
// derselben Frage bei DERSELBEN Eingabe DASSELBE antworten — die Kopplung ist
// beobachtbar statt behauptet.
//
// Die Frage hier: "ist das ein Anker?" Konsumenten: anchors (Befund
// anchor-missing) und versions (current-from, fail-closed). Ein Fall, in dem
// anchors den Anker kennt und current-from nicht — oder umgekehrt — ist die
// Klasse, dreimal eingetreten.

// ankerFall ist eine Anker-Schreibweise samt der Zieldatei, die sie enthaelt.
type ankerFall struct {
	name   string
	ziel   string // Inhalt der Zieldatei
	anker  string // das Fragment, mit dem beide Module fragen
	gueltig bool  // erwartete gemeinsame Antwort
}

func ankerFaelle() []ankerFall {
	fence := "```"
	return []ankerFall{
		{"Heading-Slug", "# Ziel\n\n## Aktuell\n\nv0.27.0\n", "aktuell", true},
		{"Duplikat-Slug", "# Ziel\n\n## Alt\n\nv0.1.0\n\n## Alt\n\nv0.27.0\n", "alt-1", true},
		{"HTML-Anker", "# Ziel\n\n<a id=\"aktuell\"></a> v0.27.0\n", "aktuell", true},
		{"kodiertes Fragment", "# Ziel\n\n<a id=\"a b\"></a> v0.27.0\n", "a%20b", true},
		{"Anker im Fence", "# Ziel\n\n" + fence + "html\n<a id=\"aktuell\"></a> v0.27.0\n" + fence + "\n", "aktuell", false},
		{"Anker in Inline-Code", "# Ziel\n\nBeispiel: `<a id=\"aktuell\"></a>` v0.27.0\n", "aktuell", false},
		{"data-id", "# Ziel\n\n<span data-id=\"aktuell\"></span> v0.27.0\n", "aktuell", false},
		{"name an Nicht-a", "# Ziel\n\n<area name=\"aktuell\"> v0.27.0\n", "aktuell", false},
		{"ohne Tag", "# Ziel\n\nid=\"aktuell\" v0.27.0\n", "aktuell", false},
		{"andere Schreibweise", "# Ziel\n\n<a id=\"Aktuell\"></a> v0.27.0\n", "aktuell", false},
	}
}

// anchorsKenntAnker fragt das Modul anchors: es meldet anchor-missing genau
// dann, wenn es den Anker NICHT kennt.
func anchorsKenntAnker(t *testing.T, ziel, anker string) bool {
	t.Helper()
	m := coretest.NewMemFS(map[string]string{
		"docs/ziel.md": ziel,
		"docs/src.md":  "Siehe [Ziel](ziel.md#" + anker + ").\n",
	})
	res, err := Run(m, nil, model.Config{}, []string{"links", "anchors"})
	if err != nil {
		t.Fatalf("anchors-Lauf unerwartet fehlgeschlagen: %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == ReasonAnchorMissing {
			return false
		}
	}
	return true
}

// versionsKenntAnker fragt das Modul versions: current-from bricht fail-closed
// ab, wenn es den Anker nicht aufloest.
func versionsKenntAnker(t *testing.T, ziel, anker string) bool {
	t.Helper()
	m := coretest.NewMemFS(map[string]string{
		"docs/ziel.md": ziel,
		"docs/use.md":  "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	cfg := versionsCfg()
	cfg.Versions.CurrentFrom = "docs/ziel.md#" + anker
	_, err := Run(m, nil, cfg, []string{"versions"})
	return err == nil
}

// Die Kopplung selbst: fuer jede Schreibweise muessen beide Module dieselbe
// Antwort geben. Weicht eine ab, ist die Lexik an ihrem Rand gedriftet.
func TestAnkerFrageHatEineAntwort(t *testing.T) {
	for _, f := range ankerFaelle() {
		t.Run(f.name, func(t *testing.T) {
			a := anchorsKenntAnker(t, f.ziel, f.anker)
			v := versionsKenntAnker(t, f.ziel, f.anker)
			if a != v {
				t.Fatalf("zwei Antworten auf dieselbe Frage: anchors=%v, versions=%v", a, v)
			}
			if a != f.gueltig {
				t.Fatalf("gemeinsame Antwort %v, erwartet %v", a, f.gueltig)
			}
		})
	}
}
