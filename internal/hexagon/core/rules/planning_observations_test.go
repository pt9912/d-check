package rules_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
)

const obsRegister = `# Beobachtungs-Register

| Kennung | Beobachtung | Zähler | Belege |
|---|---|---|---|
| BEO-002 | Ränder bleiben stehen | 7 | slice-093 |
| BEO-024 | Kanal haengt an der Arbeitsweise | 2 | slice-176 |
`

func obsCfg(dirs ...string) model.PlanningConfig {
	return model.PlanningConfig{Observations: model.ObservationsConfig{
		Register: "plan/observations.md", Dirs: dirs,
	}}
}

// Eine zitierte Kennung ohne Registerzeile meldet — die maschinelle Haelfte
// der Register-Paarung. Die Umkehrung ist ausdruecklich NICHT geprueft:
// BEO-002 steht im Register und wird nirgends zitiert.
func TestObservationsCitedWithoutRow(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"plan/observations.md": obsRegister,
		"plan/done/slice-1.md": "Ausloeser: [`BEO-024`](../observations.md) und BEO-999 daneben.\n",
	})
	got := rules.CheckPlanningObservations(fs, obsCfg("plan/done"))
	if len(got) != 1 {
		t.Fatalf("ein Befund erwartet, got %d (%v)", len(got), got)
	}
	if got[0].Reason != model.ReasonObservationUnregistered {
		t.Fatalf("Grund = %q", got[0].Reason)
	}
	if !strings.Contains(got[0].Message, "BEO-999") {
		t.Fatalf("Meldung nennt die Kennung nicht: %q", got[0].Message)
	}
}

// Die tragende Unterscheidung, gemessen am eigenen Bestand: die verbreitete
// Zitier-Form fuehrt die Kennung IN Backticks ([`BEO-NNN`](…)). Ein reines
// Code-Span ist dagegen ein Beispiel. Eine Regel "Inline-Code zaehlt nicht"
// uebersaehe die Zitate; eine ohne die Ausnahme meldete die Beispiele.
func TestObservationsCodeSpanIsExampleButLinkTextCounts(t *testing.T) {
	cases := []struct {
		name, line string
		want       int
	}{
		{"reines Code-Span ist Beispiel", "ein erfundenes `BEO-999` faellt nicht auf\n", 0},
		{"Linktext zaehlt", "siehe [`BEO-999`](../observations.md)\n", 1},
		{"Prosa zaehlt", "siehe BEO-999 im Register\n", 1},
		{"Linktext ohne Backticks zaehlt", "siehe [BEO-999](../observations.md)\n", 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs := coretest.NewMemFS(map[string]string{
				"plan/observations.md": obsRegister,
				"plan/done/s.md":       c.line,
			})
			got := rules.CheckPlanningObservations(fs, obsCfg("plan/done"))
			if len(got) != c.want {
				t.Fatalf("Befunde = %d, want %d (%v)", len(got), c.want, got)
			}
		})
	}
}

// Eine Quer-Referenz IM Register deklariert nichts: nur die erste Zelle einer
// Tabellenzeile fuehrt eine Kennung. Sonst deklarierte sich jedes Beispiel
// selbst, und der Waechter faenge nie etwas.
func TestObservationsOnlyFirstCellDeclares(t *testing.T) {
	reg := obsRegister + "| BEO-003 | Text mit BEO-777 darin | 1 | slice-1 |\n"
	fs := coretest.NewMemFS(map[string]string{
		"plan/observations.md": reg,
		"plan/done/s.md":       "siehe BEO-777\n",
	})
	if got := rules.CheckPlanningObservations(fs, obsCfg("plan/done")); len(got) != 1 {
		t.Fatalf("BEO-777 ist nicht deklariert, ein Befund erwartet, got %d (%v)", len(got), got)
	}
}

// Fail-closed: ohne Register ist die Faehigkeit inert, mit unlesbarem Register
// oder Verzeichnis meldet sie — ein Tippfehler schaltete sie sonst still ab.
func TestObservationsFailClosedAndInert(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"plan/observations.md": obsRegister})
	if got := rules.CheckPlanningObservations(fs, model.PlanningConfig{}); got != nil {
		t.Fatalf("ohne Register inert erwartet, got %v", got)
	}
	miss := model.PlanningConfig{Observations: model.ObservationsConfig{Register: "plan/fehlt.md"}}
	if got := rules.CheckPlanningObservations(fs, miss); len(got) != 1 {
		t.Fatalf("unlesbares Register: ein Befund erwartet, got %v", got)
	}
	if got := rules.CheckPlanningObservations(fs, obsCfg("plan/gibtsnicht")); len(got) != 1 {
		t.Fatalf("unlesbares Verzeichnis: ein Befund erwartet, got %v", got)
	}
}

// F-1 der Review-Runde: zwei VERSCHIEDENE ungedeckte Kennungen auf EINER Zeile
// muessen zwei Befunde ergeben. Die Deduplikation laeuft ueber (Datei, Zeile,
// Regel, Ziel, Grund) OHNE die Meldung — traegt das Ziel die Datei statt der
// Kennung, faellt der zweite Befund still weg, und die Zahl ist falsch, ohne
// dass der Lauf gruen wuerde. ids loest das seit jeher ueber dasselbe Mittel.
func TestObservationsTwoIDsOnOneLineDoNotCollide(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"plan/observations.md": obsRegister,
		"plan/done/s.md":       "zwei: [BEO-777](../observations.md) und [BEO-888](../observations.md)\n",
	})
	got := model.SortFindings(rules.CheckPlanningObservations(fs, obsCfg("plan/done")))
	if len(got) != 2 {
		t.Fatalf("zwei Kennungen auf einer Zeile: 2 Befunde erwartet, got %d (%v)", len(got), got)
	}
	if got[0].Target == got[1].Target {
		t.Fatalf("beide Befunde tragen dasselbe Ziel %q — sie kollidieren in der Deduplikation", got[0].Target)
	}
}

// Ebenso zwei unlesbare Zitier-Verzeichnisse: sie teilen Datei, Zeile, Regel
// und Grund und faenden ohne eigenes Ziel zu EINEM Befund zusammen.
func TestObservationsTwoBadDirsDoNotCollide(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"plan/observations.md": obsRegister})
	got := model.SortFindings(rules.CheckPlanningObservations(fs, obsCfg("plan/weg-a", "plan/weg-b")))
	if len(got) != 2 {
		t.Fatalf("zwei fehlende Verzeichnisse: 2 Befunde erwartet, got %d (%v)", len(got), got)
	}
}

// F-5: eine erste Zelle darf die Backtick-Konvention des Dokuments tragen,
// ohne still aufzuhoeren zu deklarieren.
func TestObservationsFirstCellMayCarryBackticks(t *testing.T) {
	reg := obsRegister + "| `BEO-003` | dekoriert | 1 | slice-1 |\n"
	fs := coretest.NewMemFS(map[string]string{
		"plan/observations.md": reg,
		"plan/done/s.md":       "siehe BEO-003\n",
	})
	if got := rules.CheckPlanningObservations(fs, obsCfg("plan/done")); len(got) != 0 {
		t.Fatalf("dekorierte erste Zelle deklariert nicht: %v", got)
	}
}
