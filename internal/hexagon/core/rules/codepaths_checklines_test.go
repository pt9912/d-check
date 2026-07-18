package rules

import (
	"fmt"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const clTarget = "L1\nL2\nL3\n" // drei Zeilen

func clRun(t *testing.T, citing string, checkLines bool) []string {
	t.Helper()
	m := coretest.NewMemFS(map[string]string{
		"docs/target.md": clTarget,
		"docs/a.md":      citing,
	})
	cfg := model.Config{Codepaths: model.CodepathsConfig{Roots: []string{"docs"}, CheckLines: checkLines}}
	res, err := Run(m, nil, cfg, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s", f.Target, f.Reason))
	}
	return got
}

// DC-FA-CODE-001.a Schritt 6: check-lines verifiziert die Zeilen-Referenz
// eines Inline-Code-Pfads (`:<von>`/`:<von>-<bis>`) gegen das Ziel.
func TestCodepathsCheckLines(t *testing.T) {
	cases := []struct {
		name, citing string
		want         []string
	}{
		{"in-range einzeln", "Siehe `docs/target.md:2`.\n", nil},
		{"in-range Bereich", "Siehe `docs/target.md:1-3`.\n", nil},
		{"hinter Datei-Ende", "Siehe `docs/target.md:9`.\n", []string{"docs/target.md:9-9 citation-out-of-range"}},
		{"Bereich zu weit", "Siehe `docs/target.md:2-9`.\n", []string{"docs/target.md:2-9 citation-out-of-range"}},
		{"invertiert", "Siehe `docs/target.md:3-1`.\n", []string{"docs/target.md:3-1 citation-inverted-range"}},
		{"von null (F-1)", "Siehe `docs/target.md:0-2`.\n", []string{"docs/target.md:0-2 citation-inverted-range"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := clRun(t, c.citing, true)
			if fmt.Sprint(got) != fmt.Sprint(c.want) {
				t.Fatalf("%s (check-lines): %v, want %v", c.name, got, c.want)
			}
		})
	}
}

// DC-QA-02: ohne check-lines wird ein einzelnes :<zahl> wie bisher
// abgetrennt und die Zeile NICHT geprüft — byte-identisch, kein Befund,
// unabhängig davon, ob die Zeile existiert.
func TestCodepathsCheckLinesDefaultAus(t *testing.T) {
	for _, citing := range []string{"Siehe `docs/target.md:2`.\n", "Siehe `docs/target.md:9`.\n"} {
		if got := clRun(t, citing, false); len(got) != 0 {
			t.Fatalf("default aus %q: %v, want 0 Befunde (byte-identisch)", citing, got)
		}
	}
}

// Ein FEHLENDES Ziel bleibt codepath-missing (der Zeilen-Check setzt
// Existenz voraus, nicht umgekehrt) — auch mit check-lines.
func TestCodepathsCheckLinesFehlendesZielBleibtMissing(t *testing.T) {
	got := clRun(t, "Weg: `docs/fehlt.md:2`.\n", true)
	want := []string{"docs/fehlt.md codepath-missing"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("fehlendes Ziel mit check-lines: %v, want %v", got, want)
	}
}
