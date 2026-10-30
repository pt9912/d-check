package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// mtCfg: Zuschlag 135/10000, Einheit 100, Schwelle 300000 (DC-FA-MET-001).
func mtCfg() model.MetricsConfig {
	return model.MetricsConfig{
		PlansDir: "plans", Block: "groessenmarke",
		SurchargeNumerator: 135, SurchargeDenominator: 10000, Rounding: 100,
		CutThreshold: 300000, CutEnabled: true,
		Table: model.MetricsTable{
			Header:    []string{"Nr.", "Zusage", "Fänger", "Mutation", "Klasse"},
			KeyColumn: "Nr.", ClassColumn: "Klasse",
			Classes: map[string]int{"x": 1000, "nebenläufig": 85, "andere": 80},
		},
	}
}

// mtTable ist die Tabelle eines Plans: eine Zeile der Klasse x (Faktor 1000).
const mtTable = "\n| Nr. | Zusage | Fänger | Mutation | Klasse |\n|---|---|---|---|---|\n| 1 | a | b | c | x |\n"

// mtPlan baut einen Plan mit einem Block; body ist der Block-Inhalt.
func mtPlan(body string) string {
	return "# Plan\n\n```groessenmarke\n" + body + "```\n\nText.\n" + mtTable
}

func mtRun(t *testing.T, files map[string]string, cfg model.MetricsConfig) (MetricsResult, error) {
	t.Helper()
	return CheckMetrics(coretest.NewMemFS(files), &fakeVCS{lineStat: driven.LineStat{Lines: 1}}, cfg)
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
		"Pfad mit Punkt-Segment":    mtPlan("path: ./src\n"),
		"Pfad mit leerem Segment":   mtPlan("path: src//x\n"),
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
	tilde := "~~~groessenmarke\n" + mtComplete("1000", "1400", "") + "~~~\n" + mtTable
	if res, err = mtOne(t, tilde); err != nil || res.WithBlock != 1 || len(res.Findings) != 0 {
		t.Fatalf("Tilden-Block: %v / %+v", err, res)
	}
	indented := "    ```groessenmarke\n" + mtComplete("1000", "1500", "") + "    ```\n" + mtTable
	if res, err = mtOne(t, indented); err != nil || len(res.Findings) != 1 {
		t.Fatalf("eingerückter Block: %v / %+v", err, res)
	}
	long := "````groessenmarke\n" + mtComplete("1000", "1400", "") + "```\n````\n"
	if _, err = mtOne(t, long); err == nil {
		t.Fatal("eine kürzere Zaun-Zeile schließt nicht; sie ist Block-Inhalt ohne ':' und damit ein Fehler")
	}
	longOK := "````groessenmarke\n" + mtComplete("1000", "1400", "") + "````\n" + mtTable
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

// Ein Plan, der ein Symlink ist, wird nicht still übergangen (fail-closed).
func TestMetricsPlanSymlinkIstEinFehler(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"plans/slice-1.md": "x\n"})
	fs.AddSymlink("plans/slice-1.md")
	if _, err := CheckMetrics(fs, &fakeVCS{}, mtCfg()); err == nil || !strings.Contains(err.Error(), "Symlink") {
		t.Fatalf("erwartet Symlink-Fehler, got %v", err)
	}
}

// Die Schwelle gilt strikt: genau auf der Schwelle ist kein Befund, einen
// Schritt darüber schon.
func TestMetricsSchwelleGenauAufDerGrenze(t *testing.T) {
	cfg := mtCfg()
	cfg.CutThreshold = 135 * 2000
	r, err := mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(mtComplete("2000", "2700", ""))}, cfg)
	if err != nil || len(r.Findings) != 0 {
		t.Fatalf("auf der Schwelle kein Befund, got %v / %+v", err, r.Findings)
	}
	r, err = mtRun(t, map[string]string{"plans/slice-1.md": mtPlan(mtComplete("2001", "2800", ""))}, cfg)
	if err != nil || len(r.Findings) != 1 || r.Findings[0].Reason != ReasonMetricsCutReasonMissing {
		t.Fatalf("über der Schwelle ein Befund, got %v / %+v", err, r.Findings)
	}
}

// Jeder Schlüssel, der höchstens einmal stehen darf, wird als Duplikat
// erkannt; Pfad und Ausnahme-Eintrag dürfen wiederholt werden.
func TestMetricsDuplikatePruefungJeSchluessel(t *testing.T) {
	for _, dup := range []string{"base: 9a0efa5\n", "limit: 1400\n", "estimate: 1000\n", "cut-reason: x\ncut-reason: y\n"} {
		if _, err := mtOne(t, mtPlan(mtComplete("1000", "1400", "cut-reason: x\n")+dup)); err == nil {
			t.Fatalf("Duplikat %q muss ein Fehler sein", dup)
		}
	}
	body := mtComplete("1000", "1400", "path: more/\nexception: 2000 | a\nexception: 3000 | b\n")
	if r, err := mtOne(t, mtPlan(body)); err != nil || len(r.Findings) != 0 {
		t.Fatalf("Pfad und Ausnahme dürfen mehrfach stehen: %v / %+v", err, r.Findings)
	}
}

// Eine negative Zahl und ein Vorzeichen sind Eingabefehler.
func TestMetricsNegativeZahl(t *testing.T) {
	for _, n := range []string{"-5", "+5"} {
		if _, err := mtOne(t, mtPlan(mtComplete(n, "100", ""))); err == nil {
			t.Fatalf("Zahl %q muss ein Fehler sein", n)
		}
	}
}

// Zitat-Ebene und Backtick-Zaun mit Backtick im Info-String sind keine Zäune
// des Blocks (der gemeinsame Öffner-Test des Produkts entscheidet).
func TestMetricsZaunAbgrenzung(t *testing.T) {
	quoted := "> ```groessenmarke\n> path: x\n> ```\n"
	if r, err := mtOne(t, quoted); err != nil || r.WithBlock != 0 {
		t.Fatalf("Zitat-Ebene ist kein Block: %v / %+v", err, r)
	}
	inline := "```groessenmarke `x`\npath: x\n```\n"
	if r, err := mtOne(t, inline); err != nil || r.WithBlock != 0 {
		t.Fatalf("Backtick im Info-String öffnet keinen Zaun: %v / %+v", err, r)
	}
}

func mtMeasured(t *testing.T, lines int, body string, v *fakeVCS) (MetricsResult, error) {
	t.Helper()
	v.lineStat.Lines = lines
	return CheckMetrics(coretest.NewMemFS(map[string]string{"plans/slice-1.md": mtPlan(body)}), v, mtCfg())
}

// Prüfung 3: Messung auf oder unter der Marke ist grün; darüber ohne
// Ausnahme-Eintrag ein Befund; mit N ≥ Messung gedeckt; mit N < Messung nicht.
func TestMetricsMessungGegenMarke(t *testing.T) {
	body := mtComplete("1000", "1400", "")
	for _, lines := range []int{0, 1399, 1400} {
		if r, err := mtMeasured(t, lines, body, &fakeVCS{}); err != nil || len(r.Findings) != 0 {
			t.Fatalf("Messung %d ≤ Marke: %v / %+v", lines, err, r.Findings)
		}
	}
	r, err := mtMeasured(t, 1401, body, &fakeVCS{})
	if err != nil || len(r.Findings) != 1 || r.Findings[0].Reason != ReasonMetricsOverLimit {
		t.Fatalf("Messung 1401 > Marke: %v / %+v", err, r.Findings)
	}
	if f := r.Findings[0]; f.Line != 7 || f.Target != "limit" || f.Rule != "metrics" {
		t.Fatalf("Befund-Felder falsch: %+v", f)
	}
	covered := mtComplete("1000", "1400", "exception: 1600 | Fix-Runden nach Review\n")
	if r, err = mtMeasured(t, 1600, covered, &fakeVCS{}); err != nil || len(r.Findings) != 0 {
		t.Fatalf("N = Messung deckt: %v / %+v", err, r.Findings)
	}
	if r, err = mtMeasured(t, 1601, covered, &fakeVCS{}); err != nil || len(r.Findings) != 1 {
		t.Fatalf("N < Messung deckt nicht: %v / %+v", err, r.Findings)
	}
}

// Ein Ausnahme-Eintrag mit leerem Grund deckt die Überschreitung trotzdem und
// ist selbst der Befund der Prüfung 4.
func TestMetricsAusnahmeOhneGrundDecktUndIstBefund(t *testing.T) {
	body := mtComplete("1000", "1400", "exception: 1600 |\n")
	r, err := mtMeasured(t, 1500, body, &fakeVCS{})
	if err != nil || len(r.Findings) != 1 || r.Findings[0].Reason != ReasonMetricsReasonEmpty {
		t.Fatalf("erwartet nur reason-empty: %v / %+v", err, r.Findings)
	}
}

// Die Messung bekommt Basis, HEAD, Pfade und ausgenommene Pfade; Hinweise für
// Binärdateien, Pfade ohne Vorkommen und einen Eintrag ohne Überschreitung.
func TestMetricsMessungArgumenteUndHinweise(t *testing.T) {
	cfg := mtCfg()
	cfg.Excluded = []string{"vertraege/"}
	v := &fakeVCS{lineStat: driven.LineStat{Lines: 10, Binary: []string{"a.bin", "b.bin"}, Missing: []string{"src/tippfehler"}}}
	body := "path: src/\npath: test/\nbase: 9a0efa5\nestimate: 1000\nlimit: 1400\nexception: 2000 | Grund\n"
	r, err := CheckMetrics(coretest.NewMemFS(map[string]string{"plans/slice-1.md": mtPlan(body)}), v, cfg)
	if err != nil || len(r.Findings) != 0 {
		t.Fatalf("erwartet befundfrei: %v / %+v", err, r.Findings)
	}
	c := v.lineCall
	if c == nil || c.base != "9a0efa5" || c.head != "HEAD" || len(c.paths) != 2 || c.paths[1] != "test/" || len(c.excluded) != 1 || c.excluded[0] != "vertraege/" {
		t.Fatalf("Aufruf falsch: %+v", c)
	}
	notes := strings.Join(r.Notes(), "\n")
	for _, want := range []string{"2 Datei(en) nicht gezählt (Binärdatei oder Submodul)", "Pfad src/tippfehler existiert weder an der Basis noch an HEAD", "exception ohne Überschreitung (Messung 10, limit 1400)"} {
		if !strings.Contains(notes, want) {
			t.Fatalf("Hinweis %q fehlt in %q", want, notes)
		}
	}
}

// Fail-closed: kein git-Port, ein Fehler des Ports (Basis nicht auflösbar,
// kein Vorfahr, kein .git) ⇒ Fehler; im Zustand Marke ausstehend keine Messung.
func TestMetricsMessungFailClosed(t *testing.T) {
	files := map[string]string{"plans/slice-1.md": mtPlan(mtComplete("1000", "1400", ""))}
	if _, err := CheckMetrics(coretest.NewMemFS(files), nil, mtCfg()); err == nil {
		t.Fatal("ohne git-Port erwartet Fehler")
	}
	_, err := CheckMetrics(coretest.NewMemFS(files), &fakeVCS{err: fmt.Errorf("Basis nicht auflösbar")}, mtCfg())
	if err == nil || !strings.Contains(err.Error(), "plans/slice-1.md") {
		t.Fatalf("Port-Fehler muss den Plan nennen: %v", err)
	}
	pending := map[string]string{"plans/slice-1.md": mtPlan("path: src/\n")}
	if _, err := CheckMetrics(coretest.NewMemFS(pending), nil, mtCfg()); err != nil {
		t.Fatalf("Marke ausstehend misst nicht: %v", err)
	}
}

const mtHead = "| Nr. | Zusage | Fänger | Mutation | Klasse |\n|---|---|---|---|---|\n"

// mtRows baut Datenzeilen `| n | a | b | c | klasse |` aus Klassen.
func mtRows(classes ...string) string {
	var sb strings.Builder
	for i, c := range classes {
		fmt.Fprintf(&sb, "| %d | a | b | c | %s |\n", i+1, c)
	}
	return sb.String()
}

func mtWithTable(body, table string) string {
	return "# Plan\n\n```groessenmarke\n" + body + "```\n\n" + table
}

// Die Zählung je Klasse: Schätzung = Klassenfaktor · Zeilen + Feldfaktor · Feld;
// Prüfung 2 hält die Schätzung des Blocks gegen sie -- als Befund nur bei
// Messung 0, danach als Hinweis.
func TestMetricsSchaetzungAusTabelle(t *testing.T) {
	cfg := mtCfg()
	cfg.EstimateInputs = map[string]int{"fakes": 200, "produktcode": 1}
	table := mtHead + mtRows("nebenläufig", "andere", "andere", "nebenläufig", "andere")
	// 85·2 + 80·3 + 200·1 + 1700 = 2310
	body := func(est string) string {
		return "path: src/\nbase: 9a0efa5\nestimate: " + est + "\nlimit: 3200\nfakes: 1\nproduktcode: 1700\ncut-reason: x\n"
	}
	run := func(est string, lines int) (MetricsResult, error) {
		v := &fakeVCS{lineStat: driven.LineStat{Lines: lines}}
		return CheckMetrics(coretest.NewMemFS(map[string]string{"plans/slice-1.md": mtWithTable(body(est), table)}), v, cfg)
	}
	if r, err := run("2310", 0); err != nil || len(r.Findings) != 0 {
		t.Fatalf("Schätzung stimmt: %v / %+v", err, r.Findings)
	}
	r, err := run("2311", 0)
	if err != nil || len(r.Findings) != 1 || r.Findings[0].Reason != ReasonMetricsEstimateMismatch {
		t.Fatalf("Schätzung weicht bei Messung 0 ab: %v / %+v", err, r.Findings)
	}
	if f := r.Findings[0]; f.Line != 6 || f.Target != "estimate" {
		t.Fatalf("Befund-Felder falsch: %+v", f)
	}
	r, err = run("2311", 5)
	if err != nil || len(r.Findings) != 0 || !strings.Contains(strings.Join(r.Notes(), "\n"), "estimate 2311 weicht von der Tabelle ab (erwartet 2310; Messung 5, nur ein Hinweis)") {
		t.Fatalf("bei Messung ≥ 1 nur ein Hinweis: %v / %+v / %v", err, r.Findings, r.Notes())
	}
}

// Marke ausstehend: der Hinweis nennt die berechnete Schätzung und Marke; eine
// fehlende Eingabe des Blocks ist ein Fehler.
func TestMetricsMarkeAusstehendMitZahlen(t *testing.T) {
	cfg := mtCfg()
	cfg.EstimateInputs = map[string]int{"fakes": 200}
	table := mtHead + mtRows("andere", "andere")
	files := func(body string) map[string]string {
		return map[string]string{"plans/slice-1.md": mtWithTable(body, table)}
	}
	r, err := CheckMetrics(coretest.NewMemFS(files("path: src/\nfakes: 2\n")), nil, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// 80·2 + 200·2 = 560 ⇒ Marke 100 · ⌈135 · 560 / 10000⌉ = 800
	if got := strings.Join(r.Notes(), "\n"); !strings.Contains(got, "Marke ausstehend (estimate 560, limit 800)") {
		t.Fatalf("Hinweis mit Zahlen fehlt: %q", got)
	}
	if _, err := CheckMetrics(coretest.NewMemFS(files("path: src/\n")), nil, cfg); err == nil {
		t.Fatal("fehlendes Eingabe-Feld muss ein Fehler sein")
	}
	if _, err := CheckMetrics(coretest.NewMemFS(files("path: src/\nfakes: 2\nfakes: 3\n")), nil, cfg); err == nil {
		t.Fatal("doppeltes Eingabe-Feld muss ein Fehler sein")
	}
	if _, err := CheckMetrics(coretest.NewMemFS(files("path: src/\nfakes: zwei\n")), nil, cfg); err == nil {
		t.Fatal("Eingabe-Feld ohne Zahl muss ein Fehler sein")
	}
}

// Negativliste der Tabellen-Erkennung: jeder Tabellenfehler ist ein
// Eingabefehler des Plans.
func TestMetricsTabellenfehler(t *testing.T) {
	body := mtComplete("1000", "1400", "")
	cases := map[string]string{
		"Tabelle fehlt":              mtWithTable(body, "Kein Tabelle hier.\n"),
		"Tabelle mehrdeutig":         mtWithTable(body, mtHead+mtRows("x")+"\nText\n\n"+mtHead+mtRows("x")),
		"Trennzeile fehlt":           mtWithTable(body, "| Nr. | Zusage | Fänger | Mutation | Klasse |\n| 1 | a | b | c | x |\n"),
		"Trennzeile falsche Zellen":  mtWithTable(body, "| Nr. | Zusage | Fänger | Mutation | Klasse |\n|---|---|---|\n| 1 | a | b | c | x |\n"),
		"Trennzeile ohne Strich":     mtWithTable(body, "| Nr. | Zusage | Fänger | Mutation | Klasse |\n|---|---|---|---|::|\n| 1 | a | b | c | x |\n"),
		"Zeile mit vier Zellen":      mtWithTable(body, mtHead+"| 1 | a | b | x |\n"),
		"Zeile mit sechs Zellen":     mtWithTable(body, mtHead+"| 1 | a | b | c | x | y |\n"),
		"Nr. leer":                   mtWithTable(body, mtHead+"|  | a | b | c | x |\n"),
		"Nr. doppelt":                mtWithTable(body, mtHead+"| 1 | a | b | c | x |\n| 1 | a | b | c | x |\n"),
		"Klasse unbekannt":           mtWithTable(body, mtHead+"| 1 | a | b | c | y |\n"),
		"Klasse falsch geschrieben":  mtWithTable(body, mtHead+"| 1 | a | b | c | X |\n"),
		"Klasse leer":                mtWithTable(body, mtHead+"| 1 | a | b | c |  |\n"),
		"Tabelle ohne Zeilen":        mtWithTable(body, mtHead),
		"Tabelle nur im Zaun":        mtWithTable(body, "```\n"+mtHead+mtRows("x")+"```\n"),
		"ungeschützter Strich im Span": mtWithTable(body, mtHead+"| 1 | `a|b` | b | c | x | y |\n"),
	}
	for name, content := range cases {
		res, err := mtOne(t, content)
		if err == nil {
			t.Fatalf("%s: erwartet Fehler, got %+v", name, res)
		}
		if !strings.Contains(err.Error(), "plans/slice-1.md") || len(res.Findings) != 0 {
			t.Fatalf("%s: Fehler nennt den Plan nicht oder Befunde erscheinen: %v / %+v", name, err, res.Findings)
		}
	}
}

// Tolerante Formen, die der gemeinsame Zeilen-Teiler liest: Schlussstrich
// fehlt, kein Anfangsstrich, `\|` in einer Zelle, CRLF, Trennzeile mit
// Doppelpunkten, eine andere Tabelle mit ähnlichem Kopf wird nicht gelesen.
func TestMetricsTabellenFormen(t *testing.T) {
	body := mtComplete("1000", "1400", "")
	other := "| Nr. | Zusage | Klasse |\n|---|---|---|\n| 1 | a | zzz |\n\n"
	forms := map[string]string{
		"ohne Schlussstrich":     mtHead + "| 1 | a | b | c | x\n",
		"ohne Anfangsstrich":     "Nr. | Zusage | Fänger | Mutation | Klasse\n---|---|---|---|---\n1 | a | b | c | x\n",
		"Strich in Zelle":        mtHead + "| 1 | a \\| a | b | c | x |\n",
		"Doppelpunkte":           "| Nr. | Zusage | Fänger | Mutation | Klasse |\n|:--|:-:|--:|---|:---:|\n| 1 | a | b | c | x |\n",
		"andere Tabelle davor":   other + mtHead + "| 1 | a | b | c | x |\n",
		"CRLF":                   strings.ReplaceAll(mtHead+"| 1 | a | b | c | x |\n", "\n", "\r\n"),
		"Tabelle endet an Text":  mtHead + "| 1 | a | b | c | x |\nText mit | Strich\n| 2 | a | b | c | y |\n",
		"Tabelle endet am Zaun":  mtHead + "| 1 | a | b | c | x |\n```\n| 2 | a | b | c | y |\n```\n",
		"Leerzeile beendet":      mtHead + "| 1 | a | b | c | x |\n\n| 2 | a | b | c | y |\n",
	}
	for name, table := range forms {
		res, err := mtOne(t, mtWithTable(body, table))
		if err != nil || len(res.Findings) != 0 {
			t.Fatalf("%s: erwartet befundfrei (Schätzung 1000 = 1 Zeile x), got %v / %+v", name, err, res.Findings)
		}
	}
}
