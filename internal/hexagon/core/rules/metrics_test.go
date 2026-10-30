package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// mtCfg ist die Konfiguration des Anlass-Repos laut Angaben des Issues:
// Zuschlag 135/10000, Einheit 100, Schwelle 300000 (DC-FA-MET-001).
func mtCfg() model.MetricsConfig {
	return model.MetricsConfig{
		PlansDir: "plans", Block: "groessenmarke",
		SurchargeNumerator: 135, SurchargeDenominator: 10000, Rounding: 100,
		CutThreshold: 300000, CutEnabled: true,
	}
}

// mtPlan baut einen Plan mit einem Block; body ist der Block-Inhalt.
func mtPlan(body string) string {
	return "# Plan\n\n```groessenmarke\n" + body + "```\n\nText.\n"
}

func mtRun(t *testing.T, files map[string]string, cfg model.MetricsConfig) (MetricsResult, error) {
	t.Helper()
	return CheckMetrics(coretest.NewMemFS(files), cfg)
}

func mtOne(t *testing.T, content string) (MetricsResult, error) {
	t.Helper()
	return mtRun(t, map[string]string{"plans/slice-1.md": content}, mtCfg())
}

func mtComplete(est, limit, extra string) string {
	return "path: src/\nbase: 9a0efa5\nestimate: " + est + "\nlimit: " + limit + "\n" + extra
}

// Happy Path und belegte Zahlen (DC-FA-MET-001): Schätzung ⇒ Marke nach
// 100 · ⌈135 · S / 10000⌉; die Grenzfälle um die Schwelle.
func TestMetricsFormelBelegteZahlen(t *testing.T) {
	cases := []struct {
		est, limit string
		cutNeeded  bool
	}{
		{"6790", "9200", true}, {"5036", "6800", true}, {"5430", "7400", true},
		{"9215", "12500", true}, {"4670", "6400", true},
		{"2000", "2700", false}, {"2001", "2800", false}, {"2222", "3000", false}, {"2223", "3100", true},
		{"0", "0", false},
	}
	for _, c := range cases {
		extra := ""
		if c.cutNeeded {
			extra = "cut-reason: Teile teilen das Format\n"
		}
		res, err := mtOne(t, mtPlan(mtComplete(c.est, c.limit, extra)))
		if err != nil || len(res.Findings) != 0 {
			t.Fatalf("S=%s: erwartet befundfrei, got %v / %+v", c.est, err, res.Findings)
		}
		if res.Plans != 1 || res.WithBlock != 1 {
			t.Fatalf("Bezugsmenge falsch: %+v", res)
		}
	}
}

// Negativ (Prüfung 1): eine Marke, die nicht der Formel entspricht.
func TestMetricsLimitMismatch(t *testing.T) {
	res, err := mtOne(t, mtPlan(mtComplete("1000", "1500", "")))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != ReasonMetricsLimitMismatch {
		t.Fatalf("erwartet limit-mismatch, got %+v", res.Findings)
	}
	f := res.Findings[0]
	if f.File != "plans/slice-1.md" || f.Rule != "metrics" || f.Target != "limit" || f.Line != 7 {
		t.Fatalf("Befund-Felder falsch: %+v", f)
	}
}

// Negativ (Prüfung 4): leerer Grund, auch nur aus Leerraum; unabhängig von
// jeder Messung.
func TestMetricsReasonEmpty(t *testing.T) {
	body := mtComplete("1000", "1400", "exception: 8000 |\nexception: 9000 |   \nexception: 9500 | wegen Review\n")
	res, err := mtOne(t, mtPlan(body))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 2 {
		t.Fatalf("erwartet zwei Befunde, got %+v", res.Findings)
	}
	for _, f := range res.Findings {
		if f.Reason != ReasonMetricsReasonEmpty || f.Target != "exception" {
			t.Fatalf("falscher Befund: %+v", f)
		}
	}
}

// Negativ (Prüfung 5): Schwelle überschritten ohne Pflichtfeld, mit leerem
// Feld; genau auf der Schwelle kein Befund.
func TestMetricsCutReasonMissing(t *testing.T) {
	res, err := mtOne(t, mtPlan(mtComplete("2223", "3100", "")))
	if err != nil || len(res.Findings) != 1 || res.Findings[0].Reason != ReasonMetricsCutReasonMissing {
		t.Fatalf("erwartet cut-reason-missing, got %v / %+v", err, res.Findings)
	}
	if res.Findings[0].Line != 6 || res.Findings[0].Target != "cut-reason" {
		t.Fatalf("Befund-Felder falsch: %+v", res.Findings[0])
	}
	res, _ = mtOne(t, mtPlan(mtComplete("2223", "3100", "cut-reason:   \n")))
	if len(res.Findings) != 1 {
		t.Fatalf("leeres Feld zählt als fehlend, got %+v", res.Findings)
	}
	res, _ = mtOne(t, mtPlan(mtComplete("2222", "3000", "")))
	if len(res.Findings) != 0 {
		t.Fatalf("2222 liegt unter der Schwelle, got %+v", res.Findings)
	}
	cfg := mtCfg()
	cfg.CutEnabled = false
	r, _ := mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(mtComplete("9215", "12500", ""))}, cfg)
	if len(r.Findings) != 0 {
		t.Fatalf("ohne Schwellenregel kein Befund, got %+v", r.Findings)
	}
}

// Zustände: kein Block und Marke ausstehend sind Hinweise, keine Befunde.
func TestMetricsZustaende(t *testing.T) {
	files := map[string]string{
		"plans/slice-1.md": "# Plan ohne Block\n",
		"plans/slice-2.md": mtPlan("path: src/\nfakes: 2\n"),
		"plans/slice-3.md": mtPlan(mtComplete("1000", "1400", "")),
	}
	res, err := mtRun(t, files, mtCfg())
	if err != nil || len(res.Findings) != 0 {
		t.Fatalf("erwartet befundfrei, got %v / %+v", err, res.Findings)
	}
	if res.Plans != 3 || res.WithBlock != 2 {
		t.Fatalf("Bezugsmenge falsch: %+v", res)
	}
	notes := strings.Join(res.Notes(), "\n")
	for _, want := range []string{"3 Plan/Pläne geprüft, 2 mit Block", "plans/slice-1.md: kein Block", "plans/slice-2.md: Marke ausstehend"} {
		if !strings.Contains(notes, want) {
			t.Fatalf("Hinweis %q fehlt in %q", want, notes)
		}
	}
}

// Plan-Suche: nur slice-*.md unmittelbar im Verzeichnis; ohne Treffer grün mit
// Hinweis; ein fehlendes oder dateiförmiges Verzeichnis ist ein Fehler.
func TestMetricsPlanSuche(t *testing.T) {
	files := map[string]string{
		"plans/notes.md":        "```groessenmarke\npath: x\n```\n",
		"plans/sub/slice-9.md":  mtPlan(mtComplete("1", "999", "")),
		"plans/slice-1.md.bak":  "x",
	}
	res, err := mtRun(t, files, mtCfg())
	if err != nil || len(res.Findings) != 0 || res.Plans != 0 {
		t.Fatalf("erwartet grün ohne Plan, got %v / %+v", err, res)
	}
	if got := strings.Join(res.Notes(), "\n"); !strings.Contains(got, "kein Plan in plans") {
		t.Fatalf("Hinweis 'kein Plan' fehlt: %q", got)
	}
	if _, err := mtRun(t, map[string]string{"x.md": "x"}, mtCfg()); err == nil {
		t.Fatal("fehlendes Plan-Verzeichnis muss ein Fehler sein")
	}
	if _, err := mtRun(t, map[string]string{"plans": "datei"}, mtCfg()); err == nil {
		t.Fatal("Plan-Verzeichnis als Datei muss ein Fehler sein")
	}
	cfg := mtCfg()
	cfg.PlansDir = ""
	if _, err := mtRun(t, files, cfg); err == nil {
		t.Fatal("leeres plans-dir muss ein Fehler sein")
	}
}

// Negativliste der Block-Erkennung: jeder Eingabefehler beendet den Lauf mit
// einem Fehler, der Plan und Ursache nennt; kein Befund erscheint.
func TestMetricsEingabefehler(t *testing.T) {
	cases := map[string]string{
		"zwei Blöcke":               mtPlan("path: a\n") + "\n```groessenmarke\npath: b\n```\n",
		"Block ohne Schluss-Zaun":   "```groessenmarke\npath: a\n",
		"Zeile ohne Doppelpunkt":    mtPlan("path a\n"),
		"Schlüssel doppelt":         mtPlan(mtComplete("1", "100", "estimate: 2\n")),
		"Zahl mit führender Null":   mtPlan("path: a\nbase: 9a0efa5\nestimate: 05\nlimit: 100\n"),
		"Zahl mit Tausendertrenner": mtPlan("path: a\nbase: 9a0efa5\nestimate: 1.000\nlimit: 100\n"),
		"Zahl leer":                 mtPlan("path: a\nbase: 9a0efa5\nestimate:\nlimit: 100\n"),
		"Zahl läuft über":           mtPlan("path: a\nbase: 9a0efa5\nestimate: 99999999999999999999\nlimit: 100\n"),
		"Basis zu kurz":             mtPlan("path: a\nbase: 9a0ef\nestimate: 1\nlimit: 100\n"),
		"Basis mit Großbuchstaben":  mtPlan("path: a\nbase: 9A0EFA5\nestimate: 1\nlimit: 100\n"),
		"Pfad mit Glob":             mtPlan("path: src/*\n"),
		"Pfad mit ..":               mtPlan("path: ../x\n"),
		"Pfad absolut":              mtPlan("path: /x\n"),
		"Pfad mit Steuerzeichen":    mtPlan("path: a\tb\x01\n"),
		"Pfad leer":                 mtPlan("path:\n"),
		"Block ohne Pfad":           mtPlan("base: 9a0efa5\nestimate: 1\nlimit: 100\n"),
		"Basis ohne Schätzung":      mtPlan("path: a\nbase: 9a0efa5\n"),
		"Ausnahme ohne Strich":      mtPlan(mtComplete("1000", "1400", "exception: 8000 Grund\n")),
		"Ausnahme ohne Zahl":        mtPlan(mtComplete("1000", "1400", "exception: viel | Grund\n")),
	}
	for name, content := range cases {
		res, err := mtOne(t, content)
		if err == nil {
			t.Fatalf("%s: erwartet Fehler, got %+v", name, res)
		}
		if !strings.Contains(err.Error(), "plans/slice-1.md") {
			t.Fatalf("%s: Fehler nennt den Plan nicht: %v", name, err)
		}
		if len(res.Findings) != 0 {
			t.Fatalf("%s: bei einem Eingabefehler erscheint kein Befund: %+v", name, res.Findings)
		}
	}
}

// Der erste Fehler in Plan-Reihenfolge gewinnt; Befunde früherer Pläne
// erscheinen nicht.
func TestMetricsErsterFehlerBeendetDenLauf(t *testing.T) {
	files := map[string]string{
		"plans/slice-1.md": mtPlan(mtComplete("1000", "1500", "")),
		"plans/slice-2.md": mtPlan("path a\n"),
		"plans/slice-3.md": mtPlan("path b\n"),
	}
	res, err := mtRun(t, files, mtCfg())
	if err == nil || !strings.Contains(err.Error(), "slice-2.md") || strings.Contains(err.Error(), "slice-3.md") {
		t.Fatalf("erwartet Fehler an slice-2.md, got %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("kein Befund bei einem Fehler, got %+v", res.Findings)
	}
}

// Zäune: andere Info-Strings verbergen ihren Inhalt, Tilden und Einrückung
// zählen, ein längerer Zaun wird nur von einem gleich langen geschlossen, CRLF
// gilt wie LF, unbekannte Schlüssel werden ignoriert.
func TestMetricsZaeune(t *testing.T) {
	hidden := "```text\n```groessenmarke\npath: x\n```\n\nText\n"
	res, err := mtOne(t, hidden)
	if err != nil || res.WithBlock != 0 {
		t.Fatalf("Block im fremden Zaun ist kein Block: %v / %+v", err, res)
	}
	tilde := "~~~groessenmarke\n" + mtComplete("1000", "1400", "") + "~~~\n"
	if res, err = mtOne(t, tilde); err != nil || res.WithBlock != 1 || len(res.Findings) != 0 {
		t.Fatalf("Tilden-Block: %v / %+v", err, res)
	}
	indented := "    ```groessenmarke\n" + mtComplete("1000", "1500", "") + "    ```\n"
	if res, err = mtOne(t, indented); err != nil || len(res.Findings) != 1 {
		t.Fatalf("eingerückter Block: %v / %+v", err, res)
	}
	long := "````groessenmarke\n" + mtComplete("1000", "1400", "") + "```\n````\n"
	if _, err = mtOne(t, long); err == nil {
		t.Fatal("eine kürzere Zaun-Zeile schließt nicht; sie ist Block-Inhalt ohne ':' und damit ein Fehler")
	}
	longOK := "````groessenmarke\n" + mtComplete("1000", "1400", "") + "````\n"
	if res, err = mtOne(t, longOK); err != nil || res.WithBlock != 1 {
		t.Fatalf("längerer Zaun: %v / %+v", err, res)
	}
	crlf := strings.ReplaceAll(mtPlan(mtComplete("1000", "1400", "produktcode: 1500\nfakes: 2\n")), "\n", "\r\n")
	if res, err = mtOne(t, crlf); err != nil || res.WithBlock != 1 || len(res.Findings) != 0 {
		t.Fatalf("CRLF und unbekannte Schlüssel: %v / %+v", err, res)
	}
}

// Feldnamen sind Konfiguration; die Defaults greifen für leere Einträge.
func TestMetricsFeldnamenKonfigurierbar(t *testing.T) {
	cfg := mtCfg()
	cfg.Fields = model.MetricsFields{Path: "pfad", Base: "basis", Estimate: "schaetzung", Limit: "marke", Exception: "ueberschreitung", CutReason: "schnitt-grund"}
	body := "pfad: src/\nbasis: 9a0efa5\nschaetzung: 2223\nmarke: 3100\nueberschreitung: 4000 |\n"
	res, err := mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(body)}, cfg)
	if err != nil || len(res.Findings) != 2 {
		t.Fatalf("erwartet zwei Befunde (leerer Grund, Schnittfrage), got %v / %+v", err, res.Findings)
	}
	f := model.MetricsConfig{}.EffectiveFields()
	if f.Path != "path" || f.CutReason != "cut-reason" || f.Exception != "exception" {
		t.Fatalf("Defaults falsch: %+v", f)
	}
}

// Überlauf der Ganzzahl-Arithmetik ist ein Eingabefehler, kein falsches Grün.
func TestMetricsUeberlauf(t *testing.T) {
	big := "9223372036854775807"
	if _, err := mtOne(t, mtPlan(mtComplete(big, "100", ""))); err == nil {
		t.Fatal("Zähler · Schätzung läuft über, erwartet Fehler")
	}
	cfg := mtCfg()
	cfg.SurchargeNumerator, cfg.SurchargeDenominator = 1, 3
	if _, err := mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(mtComplete(big, "100", ""))}, cfg); err == nil {
		t.Fatal("Aufrundung läuft über, erwartet Fehler")
	}
	cfg.SurchargeDenominator, cfg.Rounding = 1, 1<<40
	if _, err := mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(mtComplete("9000000000000000", "100", ""))}, cfg); err == nil {
		t.Fatal("Obergrenze läuft über, erwartet Fehler")
	}
}

// Determinismus (DC-QA-02): dieselbe Eingabe, dieselbe Ausgabe.
func TestMetricsDeterminismus(t *testing.T) {
	files := map[string]string{
		"plans/slice-2.md": mtPlan(mtComplete("1000", "1500", "")),
		"plans/slice-1.md": "kein Block\n",
	}
	a, _ := mtRun(t, files, mtCfg())
	b, _ := mtRun(t, files, mtCfg())
	if strings.Join(a.Notes(), "|") != strings.Join(b.Notes(), "|") || len(a.Findings) != len(b.Findings) {
		t.Fatal("zwei Läufe weichen ab")
	}
}
