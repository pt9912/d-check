package rules

import (
	"maps"
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
// der in einer realen Wellen-Closure dreimal eingetreten ist (BEO-ALL/registry-vs-authority-table-drift).
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

// F-1: ein unlesbares waves.dir ist fail-closed — auch im Ruhe-Zustand, in dem
// die leere Menge konsistent waere und ein Tippfehler die Faehigkeit sonst
// dauerhaft still abschaltete.
func TestWavesUnlesbaresVerzeichnisFailClosed(t *testing.T) {
	// MemFS kann "Verzeichnis fehlt" nicht ausdruecken (leere Liste statt
	// Fehler) — der listErrFS-Doppelgaenger traegt die Unterscheidung.
	for name, aktiv := range map[string]string{"Ruhe": "Keine aktive Welle.", "aktiv": "welle-9-neu in Arbeit."} {
		t.Run(name, func(t *testing.T) {
			fs := listErrFS{coretest.NewMemFS(map[string]string{planRoadmap: wavesRoadmap(aktiv, "", "")})}
			cfg := planCfg()
			cfg.Waves = model.WavesConfig{Dir: "docs/plan/tippfehler"}
			var got []string
			for _, f := range CheckPlanningWaves(fs, cfg) {
				got = append(got, f.Reason)
			}
			if len(got) == 0 || got[0] != model.ReasonWaveDrift {
				t.Fatalf("%s: unlesbares waves.dir muss fail-closed melden, got %v", name, got)
			}
		})
	}
}

// F-2: eine fehlende Register-Ueberschrift schaltet die Register-Aussagen nicht
// wortlos ab — die Aktivierung IST die Behauptung, dass die Roadmap beide
// Register fuehrt.
func TestWavesFehlendeRegisterUeberschriftFailClosed(t *testing.T) {
	files := map[string]string{
		// Abschluss-Register unter fremdem Namen: die Zeile ohne Notiz und die
		// Notiz ohne Zeile blieben sonst beide unsichtbar.
		planRoadmap: "# Roadmap\n\n## Aktuelle Welle\n\nwelle-9-neu in Arbeit.\n\n" +
			"## Nächste Wellen\n\n| Welle | Trigger |\n|---|---|\n\n" +
			"## Fertige Wellen\n\n| Welle | Abschluss |\n|---|---|\n| welle-8-alt | x |\n",
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-7-results.md": "# Ergebnis 7\n",
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("fehlende Register-Ueberschrift muss fail-closed melden, got %v", got)
	}
}

// F-5 (Kopplung): die Wellen-Faehigkeit liest DENSELBEN Aktiv-Status wie die
// Slice-Invariante. Beobachtbar am Fence-Fall: ein Ruhe-Marker, der NUR in
// einem Beispiel-Block steht, zaehlt nicht — eine eigene naive Bestimmung
// (Teilstring ueber alle Zeilen) hielte die Roadmap fuer ruhend und
// verschluckte das wave-drift.
func TestWavesAktivStatusIstDieGeteilteAntwort(t *testing.T) {
	fence := "```"
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			fence+"text\nKeine aktive Welle.\n"+fence+"\n\nwelle-9-neu in Arbeit.", "", ""),
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("Marker nur im Fence ⇒ aktiv ohne Dokument ⇒ wave-drift, got %v", got)
	}
}

// F-6: die flach-Haelfte von W5a — die Vorschau nennt eine Welle, deren
// PLAN-Dokument flach liegt (nicht im Ruheort).
func TestWavesVorschauMitFlachemDokument(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.",
			"| welle-9-neu | laeuft schon |\n", ""),
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWavePreviewExists {
		t.Fatalf("Vorschau-Zeile mit flachem Dokument → wave-preview-exists, got %v", got)
	}
}

// F-4: das target der Gegenrichtung ist die Ergebnisnotiz (Ventil-Paritaet),
// nicht die Kennung.
func TestWavesUnregisteredTargetIstDieNotiz(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.", "",
			"| welle-8-alt | 2026-01-01 |\n"),
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
		wavesDir + "/done/welle-7-results.md": "# Ergebnis 7\n",
	}
	fs := CheckPlanningWaves(coretest.NewMemFS(files), wavesCfg())
	if len(fs) != 1 || fs[0].Reason != model.ReasonWaveUnregistered {
		t.Fatalf("erwartet genau ein wave-unregistered, got %+v", fs)
	}
	if want := wavesDir + "/done/welle-7-results.md"; fs[0].Target != want {
		t.Fatalf("target = %q, want die Notiz %q", fs[0].Target, want)
	}
}

// F-3: ein Konsument mit anderem Praefix (wave-*) funktioniert ohne
// welle-Rueckfall.
func TestWavesFremdesPraefix(t *testing.T) {
	cfg := planCfg()
	cfg.Waves = model.WavesConfig{Dir: wavesDir, Glob: "wave-*.md", ResultsGlob: "wave-*-results.md"}
	files := map[string]string{
		planRoadmap: wavesRoadmap("wave-3-x in Arbeit.",
			"| Ein Name | Freigabe |\n", "| wave-2-y | 2026-01-01 |\n"),
		wavesDir + "/wave-3-x.md":            "# Wave 3\n",
		wavesDir + "/done/wave-2-results.md": "# Ergebnis 2\n",
	}
	if got := wavesFindings(t, files, cfg); got != nil {
		t.Fatalf("wave-Praefix-Konsument konsistent → 0 Befunde, got %v", got)
	}
}

// N-1: wave-drift traegt drei Bedeutungen — die Befund-Ziele muessen sie
// trennen, sonst kollabiert die Deduplikation zwei Reparaturen zu einer.
// Der harte Fall: done-dir == dir und BEIDE Register-Ueberschriften fehlen.
func TestWavesDedupTrenntDieDriftBedeutungen(t *testing.T) {
	cfg := planCfg()
	cfg.Waves = model.WavesConfig{Dir: wavesDir, DoneDir: wavesDir}
	files := map[string]string{
		planRoadmap:                  "# Roadmap\n\n## Aktuelle Welle\n\nwelle-9-neu in Arbeit.\n",
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	fs := CheckPlanningWaves(coretest.NewMemFS(files), cfg)
	got := model.SortFindings(fs)
	if len(got) != 2 {
		t.Fatalf("beide fehlenden Register-Ueberschriften muessen einzeln melden, got %+v", got)
	}
	if got[0].Target == got[1].Target {
		t.Fatalf("die Ziele muessen die Bedeutungen trennen: %q", got[0].Target)
	}
}

// F-14: Kopf- und Trennzeile deklarieren keine Welle, auch wenn eine Kennung
// darin steht (W4: "ohne Kopf- und Trennzeile").
func TestWavesKopfzeileMitKennungZaehltNicht(t *testing.T) {
	files := map[string]string{
		planRoadmap: "# Roadmap\n\n## Aktuelle Welle\n\nwelle-9-neu in Arbeit.\n\n" +
			"## Nächste Wellen\n\n| Welle | Trigger |\n|---|---|\n\n" +
			"## Abgeschlossene Wellen\n\n| Welle (z. B. welle-1) | Abschluss |\n|---|---|\n",
		wavesDir + "/welle-9-neu.md": "# Welle 9\n",
	}
	if got := wavesFindings(t, files, wavesCfg()); got != nil {
		t.Fatalf("Kennung in der Kopfzeile darf nicht zaehlen, got %v", got)
	}
}

// wavesManyCfg aktiviert die Bijektions-Semantik (mode: many, Lastenheft 0.62.0).
func wavesManyCfg() model.PlanningConfig {
	c := wavesCfg()
	c.Waves.Mode = "many"
	return c
}

// W3 unter many: zwei gelistete Wellen, zwei Dateien, der Ruhe-Marker steht
// daneben — kein Befund (Lastenheft §Wellen-Happy-Path many).
func TestWavesManyBijektionHappyPath(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			"- [welle-1-basis](../welle-1-basis.md)\n- [welle-2-ausbau](../welle-2-ausbau.md)\n\nKeine aktive Welle.",
			"", ""),
		wavesDir + "/welle-1-basis.md":  "# Welle 1\n",
		wavesDir + "/welle-2-ausbau.md": "# Welle 2\n",
	}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("Kennungs-Mengen decken sich (Marker orthogonal) → 0 Befunde, got %v", got)
	}
}

// Marker-Orthogonalität (Lastenheft §Wellen-Boundary): Welle offen, nichts
// beansprucht, Marker ZUSÄTZLICH zum Zeiger — unter many meldet weder
// wave-drift noch planning-drift; unter one meldet derselbe Zustand
// wave-drift (die Abgrenzung der beiden Modelle, kein Fehler).
func TestWavesManyMarkerOrthogonal(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			"- [welle-1-basis](../welle-1-basis.md)\n\nKeine aktive Welle.", "", ""),
		wavesDir + "/welle-1-basis.md": "# Welle 1\n",
	}
	if got := CheckPlanning(coretest.NewMemFS(files), wavesManyCfg()); got != nil {
		t.Fatalf("Marker bei leerem Slice-Verzeichnis → kein planning-drift, got %+v", got)
	}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("many: das Eröffnungs-Fenster ist legitim → 0 Befunde, got %v", got)
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("one: derselbe Zustand bleibt wave-drift (Abgrenzung, Default byte-identisch), got %v", got)
	}
}

// W3 unter many, beide Richtungen: das Befund-target ist die KENNUNG — zwei
// Richtungen an derselben Blockzeile bleiben im Dedup-Tupel verschieden
// (Ventil-Parität wie bei den Register-Aussagen).
func TestWavesManyBijektionBeideRichtungen(t *testing.T) {
	files := map[string]string{
		planRoadmap:                    wavesRoadmap("- welle-3-geist ist offen\n", "", ""),
		wavesDir + "/welle-4-still.md": "# Welle 4\n",
	}
	fnd := CheckPlanningWaves(coretest.NewMemFS(files), wavesManyCfg())
	if len(fnd) != 2 || fnd[0].Reason != model.ReasonWaveDrift || fnd[1].Reason != model.ReasonWaveDrift {
		t.Fatalf("Zeiger ohne Datei + Datei ohne Zeiger → 2× wave-drift, got %+v", fnd)
	}
	if fnd[0].Target != "welle-3" || fnd[1].Target != "welle-4" {
		t.Fatalf("target = betroffene Kennung (welle-3, welle-4), got %q, %q",
			fnd[0].Target, fnd[1].Target)
	}
}

// Fence-Inhalte im Block zählen nicht als Zeiger, Mehrfachnennung derselben
// Kennung zählt einmal (Erkennungs-Verfahren der Register: Prosa-Zeilen).
func TestWavesManyFenceUndMehrfachnennung(t *testing.T) {
	aktiv := "- [welle-5-echt](../welle-5-echt.md) — Zweitnennung: welle-5-echt\n\n" +
		"```\nwelle-6-beispiel\n```\n"
	files := map[string]string{
		planRoadmap:                   wavesRoadmap(aktiv, "", ""),
		wavesDir + "/welle-5-echt.md": "# Welle 5\n",
	}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("Fence zählt nicht, Mehrfachnennung einmal → 0 Befunde, got %v", got)
	}
}

// Kardinalität null: leerer Block (nur Marker) und kein flaches Dokument —
// die leere Menge gleicht der leeren Menge.
func TestWavesManyNullKardinalitaet(t *testing.T) {
	files := map[string]string{planRoadmap: wavesRoadmap("Keine aktive Welle.", "", "")}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("beidseitig leer → 0 Befunde, got %v", got)
	}
}

// Mehr-Wellen-Betrieb bleibt unter one der Mehrfach-Fall, unter many ist er
// der Normalfall — dieselben Dateien, zwei Modelle.
func TestWavesManyMehrWellenGegenprobe(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			"- [welle-1-basis](../welle-1-basis.md)\n- [welle-2-ausbau](../welle-2-ausbau.md)\n", "", ""),
		wavesDir + "/welle-1-basis.md":  "# Welle 1\n",
		wavesDir + "/welle-2-ausbau.md": "# Welle 2\n",
	}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("many: n Zeiger, n Dateien → 0 Befunde, got %v", got)
	}
	got := wavesFindings(t, files, wavesCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("one: zwei flache Dokumente bleiben wave-drift, got %v", got)
	}
}

// fail-closed bleibt modus-unabhängig: ein unlesbares waves.dir meldet auch
// unter many wave-drift statt still zu grünen (der listErrFS-Doppelgänger
// trägt die Unterscheidung, die MemFS nicht ausdrücken kann).
func TestWavesManyFailClosedBleibt(t *testing.T) {
	fs := listErrFS{coretest.NewMemFS(map[string]string{
		planRoadmap: wavesRoadmap("Keine aktive Welle.", "", ""),
	})}
	cfg := wavesManyCfg()
	cfg.Waves.Dir = "docs/plan/tippfehler"
	var got []string
	for _, f := range CheckPlanningWaves(fs, cfg) {
		got = append(got, f.Reason)
	}
	if len(got) == 0 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("unlesbares Verzeichnis → fail-closed wave-drift, got %v", got)
	}
}

// Explizites mode: one verhält sich wie der Default — dieselbe Eingabe,
// derselbe Befund (die Durchreichung prüft der Config-Rand-Test).
func TestWavesExplizitesOneBleibtSingleton(t *testing.T) {
	cfg := wavesCfg()
	cfg.Waves.Mode = "one"
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			"- [welle-1-basis](../welle-1-basis.md)\n\nKeine aktive Welle.", "", ""),
		wavesDir + "/welle-1-basis.md": "# Welle 1\n",
	}
	got := wavesFindings(t, files, cfg)
	if len(got) != 1 || got[0] != model.ReasonWaveDrift {
		t.Fatalf("mode: one = Default-Singleton, got %v", got)
	}
}

// Fehlt die kanonische Überschrift, meldet die Slice-Invariante fail-closed —
// die Wellen-Fähigkeit schweigt auch unter many (kein Doppelbefund; der
// fail-Pfad liegt vor der Modus-Verzweigung).
func TestWavesManyUnbestimmbarerAktivStatusSchweigt(t *testing.T) {
	files := map[string]string{
		planRoadmap:                    "# Roadmap\n\nOhne kanonische Überschrift.\n",
		wavesDir + "/welle-1-basis.md": "# Welle 1\n",
	}
	if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
		t.Fatalf("unbestimmbarer Status → Schweigen der Wellen-Fähigkeit, got %v", got)
	}
}

// Die Register-Aussagen (W4/W5) greifen unter many unverändert: gefüllte
// Register bleiben grün, eine Abschluss-Zeile ohne Notiz meldet weiter.
func TestWavesManyRegisterAussagenUnveraendert(t *testing.T) {
	files := map[string]string{
		planRoadmap: wavesRoadmap(
			"- [welle-9-neu](../welle-9-neu.md)\n",
			"| Ein geplanter Name | Freigabe |\n",
			"| welle-8-alt | 2026-01-01 |\n| welle-7-lueck | 2026-01-02 |\n"),
		wavesDir + "/welle-9-neu.md":          "# Welle 9\n",
		wavesDir + "/done/welle-8-results.md": "# Ergebnis 8\n",
	}
	got := wavesFindings(t, files, wavesManyCfg())
	if len(got) != 1 || got[0] != model.ReasonWaveResultsMissing {
		t.Fatalf("Register-Aussagen bleiben unter many scharf (welle-7 ohne Notiz), got %v", got)
	}
}

// Kennungs-Menge statt Datei-Menge (DC-FA-PLAN-001 §Wellen-Boundary gleiche
// Kennung): zwei flache Dokumente derselben Kennung sind EIN Element — unter
// one genügt der genannte Aktiv-Status, unter many ein Zeiger dieser Kennung.
func TestWavesGleicheKennungZaehltEinmal(t *testing.T) {
	doppel := map[string]string{
		wavesDir + "/welle-9-neu.md":      "# Welle 9\n",
		wavesDir + "/welle-9-nachtrag.md": "# Welle 9, zweites Dokument\n",
	}
	t.Run("one", func(t *testing.T) {
		files := map[string]string{planRoadmap: wavesRoadmap("welle-9-neu in Arbeit.", "", "")}
		maps.Copy(files, doppel)
		if got := wavesFindings(t, files, wavesCfg()); got != nil {
			t.Fatalf("one: zwei Dateien derselben Kennung = genau eine Kennung → 0 Befunde, got %v", got)
		}
	})
	t.Run("many", func(t *testing.T) {
		files := map[string]string{planRoadmap: wavesRoadmap(
			"- [welle-9-neu](../welle-9-neu.md)\n\nKeine aktive Welle.", "", "")}
		maps.Copy(files, doppel)
		if got := wavesFindings(t, files, wavesManyCfg()); got != nil {
			t.Fatalf("many: ein Zeiger deckt zwei Dateien derselben Kennung → 0 Befunde, got %v", got)
		}
	})
}
