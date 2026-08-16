package rules

import (
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const (
	planRoadmap = "docs/plan/planning/in-progress/roadmap.md"
	planSlice   = "docs/plan/planning/in-progress/slice-001-x.md"
)

func planCfg() model.PlanningConfig { return model.PlanningConfig{Roadmap: planRoadmap} }

// roadmapText baut eine Roadmap mit gegebenem Aktuelle-Welle-Block. Der
// Ruhe-Marker steht zusätzlich in einer **späteren** Sektion — so verifiziert
// jeder Fall das Block-Scoping (der Marker außerhalb des Blocks darf den Status
// nicht verfälschen).
func roadmapText(heading, block string) string {
	return "# Roadmap\n\n" + heading + "\n\n" + block +
		"\n\n## Nächste Wellen\n\nKeine aktive Welle (nur erklärende Prosa hier).\n"
}

// TestCheckPlanning deckt die Selbsttest-Klassen des abgelösten
// tools/planning-consistency.sh als Modul-Tests ab (ADR-0028): beide
// Drift-Richtungen, beide Gegenproben, Heading-Guard, Block-Scoping.
func TestCheckPlanning(t *testing.T) {
	active := "welle-x-modul in Arbeit."
	idle := "Keine aktive Welle — wartet auf Trigger."
	cases := []struct {
		name    string
		heading string
		block   string
		slice   bool
		drift   bool
	}{
		{"aktiv+slice konsistent (Block-Scoping: Marker in Nächste-Wellen zählt nicht)", "## Aktuelle Welle", active, true, false},
		{"idle+kein-slice konsistent", "## Aktuelle Welle", idle, false, false},
		{"slice+idle drift (Richtung A)", "## Aktuelle Welle", idle, true, true},
		{"kein-slice+aktiv drift (Richtung B)", "## Aktuelle Welle", active, false, true},
		{"kaputte Überschrift fail-closed (Richtung C)", "## Aktuelle Wellen", idle, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{planRoadmap: roadmapText(tc.heading, tc.block)}
			if tc.slice {
				files[planSlice] = "# slice\n"
			}
			f := CheckPlanning(coretest.NewMemFS(files), planCfg())
			if (len(f) > 0) != tc.drift {
				t.Fatalf("drift=%v, erwartet %v (findings %+v)", len(f) > 0, tc.drift, f)
			}
			if tc.drift {
				if f[0].Reason != model.ReasonPlanningDrift || f[0].Rule != "planning" || f[0].File != planRoadmap {
					t.Errorf("Befund falsch: %+v", f[0])
				}
			}
		})
	}
}

// TestCheckPlanningDuplicateHeading: zwei kanonische Überschriften ⇒ mehrdeutiger
// Aktiv-Status ⇒ planning-drift (fail-closed; verhindert das Silent-Grün, das ein
// nur den ersten Block prüfendes Modul bei Marker im zweiten Block hätte).
// **Testaufbau so, dass NUR der Duplikat-Guard die Drift erzeugt** (erster Block
// aktiv + Slice ⇒ per Erst-Block konsistent): ohne den Guard bliebe der Lauf grün,
// darum killt dieser Test die Guard-Mutation (R3-MEDIUM slice-057).
func TestCheckPlanningDuplicateHeading(t *testing.T) {
	rm := "# R\n\n## Aktuelle Welle\n\nwelle-x aktiv.\n\n## Aktuelle Welle\n\nKeine aktive Welle.\n\n## Ende\n"
	files := map[string]string{planRoadmap: rm, planSlice: "# s\n"} // Erst-Block aktiv + Slice ⇒ ohne Guard konsistent
	f := CheckPlanning(coretest.NewMemFS(files), planCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonPlanningDrift {
		t.Fatalf("mehrfache Überschrift (Erst-Block aktiv + Slice) ⇒ planning-drift (fail-closed) erwartet, bekam %+v", f)
	}
}

// TestCheckPlanningFailClosedMissingFile: fehlende Roadmap-Datei ⇒ planning-drift
// (kein stilles Grün).
func TestCheckPlanningFailClosedMissingFile(t *testing.T) {
	f := CheckPlanning(coretest.NewMemFS(map[string]string{}), planCfg())
	if len(f) != 1 || f[0].Reason != model.ReasonPlanningDrift {
		t.Fatalf("fehlende Roadmap ⇒ planning-drift erwartet, bekam %+v", f)
	}
}

// TestCheckPlanningInert: leere Roadmap bzw. nil-fs ⇒ inert (byte-identisch).
func TestCheckPlanningInert(t *testing.T) {
	if f := CheckPlanning(coretest.NewMemFS(map[string]string{planRoadmap: "x"}), model.PlanningConfig{}); f != nil {
		t.Fatalf("leere Roadmap ⇒ inert, bekam %+v", f)
	}
	if f := CheckPlanning(nil, planCfg()); f != nil {
		t.Fatalf("nil-fs ⇒ inert, bekam %+v", f)
	}
}

// TestCheckPlanningCustomConfig: heading/marker/slice-glob sind überschreibbar
// (parametrierbares Modul, nicht hartcodiert).
func TestCheckPlanningCustomConfig(t *testing.T) {
	cfg := model.PlanningConfig{Roadmap: planRoadmap, Heading: "## Current", Marker: "IDLE", SliceGlob: "wave-*.md"}
	rm := "# R\n\n## Current\n\nIDLE now.\n"
	// idle + kein wave-Slice ⇒ konsistent
	if f := CheckPlanning(coretest.NewMemFS(map[string]string{planRoadmap: rm}), cfg); f != nil {
		t.Fatalf("custom idle+kein-Slice konsistent, bekam %+v", f)
	}
	// idle + wave-Slice ⇒ drift (slice-*.md-Default würde NICHT matchen → prüft die Glob-Config)
	files := map[string]string{planRoadmap: rm, "docs/plan/planning/in-progress/wave-01.md": "x"}
	if f := CheckPlanning(coretest.NewMemFS(files), cfg); len(f) != 1 || f[0].Reason != model.ReasonPlanningDrift {
		t.Fatalf("custom idle+wave-Slice ⇒ drift erwartet, bekam %+v", f)
	}
}

// DC-FA-PLAN-001: eine Ueberschrift IM Fence ist keine. Eine Roadmap, die ihren
// eigenen Abschnitt in einem Beispiel-Block zeigt, galt sonst als "mehrfach
// vorhanden" und meldete planning-drift — ein Falsch-Rot.
func TestPlanningUeberschriftImFenceZaehltNicht(t *testing.T) {
	block := "```markdown\n## Aktuelle Welle\nBeispiel\n```\n\nKeine aktive Welle."
	m := coretest.NewMemFS(map[string]string{
		planRoadmap: roadmapText("## Aktuelle Welle", block),
	})
	if f := CheckPlanning(m, planCfg()); f != nil {
		t.Fatalf("Beispiel-Ueberschrift im Fence darf nicht zaehlen, got %+v", f)
	}
}

// Dieselbe Achse, andere Wirkung: eine Raute-Zeile IM Fence beendete den
// Aktiv-Block vorzeitig, der Ruhe-Marker dahinter ging verloren — ebenfalls
// Falsch-Rot.
func TestPlanningRautenzeileImFenceBeendetDenBlockNicht(t *testing.T) {
	block := "```sh\n## nur ein Kommentar im Beispiel\n```\n\nKeine aktive Welle."
	m := coretest.NewMemFS(map[string]string{
		planRoadmap: roadmapText("## Aktuelle Welle", block),
	})
	if f := CheckPlanning(m, planCfg()); f != nil {
		t.Fatalf("Raute-Zeile im Fence darf den Block nicht beenden, got %+v", f)
	}
}
