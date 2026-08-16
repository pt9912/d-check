package rules

import (
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const wavesDir = "docs/plan/planning"

// wavesCfg aktiviert die dritte Faehigkeit auf dem Standard-Layout.
func wavesCfg() model.PlanningConfig {
	c := planCfg()
	c.Waves = model.WavesConfig{Dir: wavesDir}
	return c
}

// wavesRoadmap baut eine Roadmap aus den drei Wellen-Abschnitten.
func wavesRoadmap(aktiv, vorschau, abgeschlossen string) string {
	return "# Roadmap\n\n## Aktuelle Welle\n\n" + aktiv +
		"\n\n## Nächste Wellen\n\n| Welle | Trigger |\n|---|---|\n" + vorschau +
		"\n## Abgeschlossene Wellen\n\n| Welle | Abschluss |\n|---|---|\n" + abgeschlossen
}

// wavesFindings verdichtet den Lauf auf die Grund-Codes der Wellen-Faehigkeit.
func wavesFindings(t *testing.T, files map[string]string, cfg model.PlanningConfig) []string {
	t.Helper()
	var out []string
	for _, f := range CheckPlanningWaves(coretest.NewMemFS(files), cfg) {
		out = append(out, f.Reason)
	}
	return out
}

// Happy Path: aktive Welle mit genau einem flachen Dokument, Vorschau ohne
// Kennungen, Register und Ergebnisnotizen decken sich beidseitig.
func TestWavesHappyPath(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.",
			"| Ein geplanter Name | Freigabe |\n",
			"| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md":         "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
	}
	if got := wavesFindings(t, files, wavesCfg()); got != nil {
		t.Fatalf("konsistente Wellen-Lage → 0 Befunde, got %v", got)
	}
}

// W3, beide Richtungen und der Mehrfach-Fall.
func TestWavesDrift(t *testing.T) {
	faelle := map[string]map[string]string{
		"aktiv ohne Dokument": {
			planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.", "", ""),
		},
		"Ruhe mit Dokument": {
			planRoadmap:                  wavesRoadmap("Keine aktive Welle.", "", ""),
			wavesDir + "/welle-9-neu.md": "# Welle 9\n",
		},
		"aktiv mit zwei Dokumenten": {
			planRoadmap:                  wavesRoadmap("welle-9-neu in Arbeit.", "", ""),
			wavesDir + "/welle-9-neu.md": "# Welle 9\n",
			wavesDir + "/welle-10-x.md":  "# Welle 10\n",
		},
	}
	for name, files := range faelle {
		t.Run(name, func(t *testing.T) {
			got := wavesFindings(t, files, wavesCfg())
			if len(got) != 1 || got[0] != model.ReasonWaveDrift {
				t.Fatalf("%s → genau ein wave-drift erwartet, got %v", name, got)
			}
		})
	}
}

// W5a: eine Vorschau-Zeile nennt eine Welle, die bereits eine Datei hat.
func TestWavesVorschauMitDatei(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.",
			"| welle-8-alt | Freigabe |\n", "| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWavePreviewExists {
		t.Fatalf("Vorschau-Zeile mit vorhandener Datei → wave-preview-exists, got %v", got)
	}
}

// W5b: eine Abschluss-Zeile ohne Ergebnisnotiz.
func TestWavesAbschlussOhneNotiz(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.", "",
			"| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveResultsMissing {
		t.Fatalf("Abschluss-Zeile ohne Notiz → wave-results-missing, got %v", got)
	}
}

// W5c: die Gegenrichtung — Ergebnisnotiz ohne Registerzeile. Genau der Fall,
// der in einer realen Wellen-Closure dreimal eingetreten ist (BEO-001).
func TestWavesNotizOhneRegisterzeile(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.", "",
			"| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
		wavesDir + "/done/welle-7-results.md": "# Ergebnis 7\n",
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveUnregistered {
		t.Fatalf("Notiz ohne Registerzeile → wave-unregistered, got %v", got)
	}
}

// Rollen-Trennung: eine Ergebnisnotiz, die im WELLEN-Verzeichnis liegt, zaehlt
// nicht als flaches Plan-Dokument — sonst waere die Aktiv-Status-Aussage von
// der Ablage der Notizen abhaengig.
func TestWavesErgebnisnotizIstKeinPlanDokument(t *testing.T) {
	files := map[string]string{
		planRoadmap:                          wavesRoadmap("Keine aktive Welle.", "", ""),
		wavesDir + "/welle-8-results.md":     "# Ergebnis 8 (faelschlich flach)\n",
	}
	if got := wavesFindings(t, files, wavesCfg()); got != nil {
		t.Fatalf("Ergebnisnotiz ist kein Plan-Dokument → 0 Befunde, got %v", got)
	}
}

// Namens-Vorschau: die erste Spalte traegt einen Namen, die TRIGGER-Spalte
// nennt eine laufende Welle. Gelesen wird die erste Spalte — ein Zeilen-Scan
// meldete hier falsch.
func TestWavesVorschauNameMitWelleImTrigger(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.",
			"| Ein geplanter Name | Freigabe; nach welle-9-neu geschnitten |\n", ""),
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	if got := wavesFindings(t, files, wavesCfg()); got != nil {
		t.Fatalf("Kennung nur in der Trigger-Spalte → 0 Befunde, got %v", got)
	}
}

// Ohne waves.dir wird kein Wellendokument geoeffnet.
func TestWavesInertOhneDir(t *testing.T) {
	files := map[string]string{
		planRoadmap:                  wavesRoadmap("Keine aktive Welle.", "", ""),
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	if got := wavesFindings(t, files, planCfg()); got != nil {
		t.Fatalf("ohne waves.dir inert → 0 Befunde, got %v", got)
	}
}

// Eine Tabellenzeile IM Fence ist keine — dieselbe Lexik wie im Modul targets.
func TestWavesTabellenzeileImFenceZaehltNicht(t *testing.T) {
	fence := "```"
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.",
			fence+"markdown\n| welle-8-alt | Beispiel |\n"+fence+"\n",
			"| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
	}
	if got := wavesFindings(t, files, wavesCfg()); got != nil {
		t.Fatalf("Beispiel-Tabelle im Fence darf nicht zaehlen, got %v", got)
	}
}
