package rules

import (
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Dieser Test ist die Verkoerperung von BEO-ALL/shared-lexicon-drifts-at-edges: eine geteilte Lexik driftet an
// den Raendern, weil jeder Konsument sie selbst vorbereitet. Er prueft nicht,
// dass eine bestimmte Funktion aufgerufen wird, sondern dass alle Konsumenten
// derselben Frage bei DERSELBEN Eingabe DASSELBE antworten — die Kopplung ist
// beobachtbar statt behauptet.
//
// Die Frage hier: "ist das ein Anker?" Konsumenten sind ALLE DREI, die sie
// stellen: anchors (Befund anchor-missing), versions (current-from,
// fail-closed) und pins (aufloesbarer Ziel-Span). Ein Fall, in dem einer den
// Anker kennt und ein anderer nicht, ist die Klasse — dreimal eingetreten.
//
// Der Test misst die AUFLOESBARKEIT des Ankers, nicht den Zuschnitt des
// adressierten Spans; letzterer ist je Modul vertraglich verschieden
// (ADR-0019/ADR-0020) und gehoert in die Modul-Tests.

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
	cfg.Versions.Patterns[0].CurrentFrom = "docs/ziel.md#" + anker
	_, err := Run(m, nil, cfg, []string{"versions"})
	return err == nil
}

// pinsKenntAnker fragt das Modul pins: es hasht den Ziel-Span nur, wenn es den
// Anker aufloest — sonst schweigt es (kein Ersatz-Befund). Der Pin traegt einen
// absichtlich falschen Hash: kommt link-stale, war der Anker aufloesbar.
func pinsKenntAnker(t *testing.T, ziel, anker string) bool {
	t.Helper()
	m := coretest.NewMemFS(map[string]string{
		"docs/ziel.md": ziel,
		"docs/src.md": "Siehe [Ziel](ziel.md#" + anker + ") <!-- dpin: sha256:" +
			strings.Repeat("a", 64) + " -->.\n",
	})
	res, err := Run(m, nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatalf("pins-Lauf unerwartet fehlgeschlagen: %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == model.ReasonLinkStale {
			return true
		}
	}
	return false
}

// Zweite Kopplung, seit die Chronologie-Bedingung der dritte Konsument ist
// (ADR-0057, Form aus ADR-0054): die Frage "zaehlt diese Zeile als
// Tabellenzeile?" (tableRowLine) muss von targets (gate-phantom aus
// Doku-Tabellen), planning.waves (Registerzeilen) und structure (Datenzeilen
// der Chronologie-Bedingung) GLEICH beantwortet werden. Die Schreibweisen
// teilen sich eine gemeinsame Antwort — eine Zeile im Fence ist ein Beispiel,
// eine eingerueckte Zeile (ein bis vier Leerzeichen, Tab) ist keine
// Spalte-0-Tabellenzeile.
//
// linelength (Ausnahme (b)) stellt eine ANDERE Frage — "ist die Zeile Teil
// einer Tabelle und deshalb nicht umbrechbar?" — und beantwortet sie mit
// unbreakableTableRow (DC-FA-LLEN-001, ADR-0116, ADR-0054 Entscheidung 2); sie
// schliesst tableRowLine ein und erweitert sie um bis drei Leerzeichen. Ihre
// Kopplung ist die zweite Assertion unten, ueber ihre eigenen Faelle.
//
// Reichweiten-Notiz (Review F-3): die Kopf-/Trennzeilen-Antwort
// (tableHeaderOrSeparator) und die Zell-Antwort (tableCells) haben heute je
// ZWEI Konsumenten (planning.waves, structure) und sind ueber die gemeinsame
// Funktion gebunden, nicht ueber eine Kopplung. Nach der ADR-0054-Schwelle
// bekommt jede dieser Fragen ihren Kopplungs-Test mit dem DRITTEN
// Konsumenten — wer einen hinzufuegt, schreibt ihn hier daneben.

// tabellenFall verpackt die zu pruefende Zeile in ihren Kontext.
type tabellenFall struct {
	name    string
	wrap    func(row string) string
	gueltig bool
}

func tabellenFaelle() []tabellenFall {
	fence := "```"
	return []tabellenFall{
		{"ausserhalb Fence", func(r string) string { return r + "\n" }, true},
		{"im Fence", func(r string) string { return fence + "\n" + r + "\n" + fence + "\n" }, false},
		{"ein Leerzeichen", func(r string) string { return " " + r + "\n" }, false},
		{"eingerueckt", func(r string) string { return "  " + r + "\n" }, false},
		{"drei Leerzeichen", func(r string) string { return "   " + r + "\n" }, false},
		{"vier Leerzeichen", func(r string) string { return "    " + r + "\n" }, false},
		{"Tab", func(r string) string { return "\t" + r + "\n" }, false},
	}
}

// tabellenFaelleLinelength sind die Faelle der Frage "Teil einer Tabelle und
// deshalb nicht umbrechbar?" (unbreakableTableRow): gueltig heisst, die Zeile
// ist ausgenommen. Bis drei Leerzeichen davor und Spalte 0 ausgenommen; vier
// Leerzeichen, Tab, U+00A0, `>`, eine Listenmarke davor und der Fence
// gemessen.
func tabellenFaelleLinelength() []tabellenFall {
	fence := "```"
	return []tabellenFall{
		{"ausserhalb Fence", func(r string) string { return r + "\n" }, true},
		{"im Fence", func(r string) string { return fence + "\n" + r + "\n" + fence + "\n" }, false},
		{"ein Leerzeichen", func(r string) string { return " " + r + "\n" }, true},
		{"zwei Leerzeichen", func(r string) string { return "  " + r + "\n" }, true},
		{"drei Leerzeichen", func(r string) string { return "   " + r + "\n" }, true},
		{"drei Leerzeichen im Fence", func(r string) string { return fence + "\n   " + r + "\n" + fence + "\n" }, false},
		{"vier Leerzeichen", func(r string) string { return "    " + r + "\n" }, false},
		{"Tab", func(r string) string { return "\t" + r + "\n" }, false},
		{"geschuetztes Leerzeichen", func(r string) string { return "\u00a0" + r + "\n" }, false},
		{"Zitat", func(r string) string { return "> " + r + "\n" }, false},
		{"Listenmarke", func(r string) string { return "- " + r + "\n" }, false},
	}
}

// targetsKenntZeile: ein `make ghost` in der Zeile erzeugt gate-phantom genau
// dann, wenn die Zeile als Tabellenzeile zaehlt (das Makefile kennt ghost nicht).
func targetsKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	files := map[string]string{
		"Makefile":  "",
		"docs/x.md": "# D\n\n" + wrap("| `make ghost` | x |"),
	}
	cfg := model.TargetsConfig{
		Makefiles: []string{"Makefile"}, DocTables: []string{"docs/x.md"}, Authority: []string{"docs/x.md"},
	}
	f, err := CheckTargets(coretest.NewMemFS(files), cfg)
	if err != nil {
		t.Fatalf("targets-Lauf unerwartet fehlgeschlagen: %v", err)
	}
	for _, x := range f {
		if x.Reason == ReasonGatePhantom {
			return true
		}
	}
	return false
}

// wavesKenntZeile: eine Registerzeile ohne Ergebnisnotiz erzeugt
// wave-results-missing genau dann, wenn die Zeile als Tabellenzeile zaehlt.
func wavesKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	roadmap := "# R\n\n## Aktuelle Welle\n\nKeine aktive Welle.\n\n" +
		"## Nächste Wellen\n\n| Welle | Trigger |\n|---|---|\n\n" +
		"## Abgeschlossene Wellen\n\n| Welle | Abschluss |\n|---|---|\n" +
		wrap("| welle-8-alt | 2026-01-01 |")
	files := map[string]string{planRoadmap: roadmap}
	for _, f := range CheckPlanningWaves(coretest.NewMemFS(files), wavesCfg()) {
		if f.Reason == model.ReasonWaveResultsMissing {
			return true
		}
	}
	return false
}

// structureKenntZeile: eine untypisierbare Datenzeile erzeugt
// section-cell-untyped genau dann, wenn die Zeile als Tabellenzeile zaehlt
// (sonst meldet der Abschnitt den Leerlauf — ein ANDERER Code).
func structureKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	files := map[string]string{
		"docs/a.md": "# T\n\n## H\n\n" + wrap("| kaputt |"),
	}
	r := model.StructureRule{Files: "docs/*.md", Section: "## H", Table: model.TableRule{Order: "desc"}}
	for _, f := range CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r}) {
		if f.Reason == model.ReasonSectionCellUntyped {
			return true
		}
	}
	return false
}

// linelengthKenntZeile: eine überlange Zeile ist ausgenommen genau dann, wenn
// sie Teil einer Tabelle und deshalb nicht umbrechbar ist (Antwort von
// unbreakableTableRow). Ausnahme (a) und (c) sind aus, damit nur (b)
// antwortet: eine Zeile im Zaun ist dann keine Tabellenzeile.
func linelengthKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	cfg := llCfg(10)
	cfg.Linelength.ExemptCodeBlocks = llOff()
	cfg.Linelength.ExemptUnbreakable = llOff()
	res := llRun(t, map[string]string{"docs/a.md": "# T\n\n" + wrap("| zelle | zweite zelle |")}, cfg)
	return len(llLines(res, "docs/a.md")) == 0
}

// unbreakableTableRowKenntZeile fragt die Lexik direkt: dieselbe Zeile, dieselbe
// Prosa-Menge, die Antwort des Helfers statt der des Moduls.
func unbreakableTableRowKenntZeile(wrap func(string) string) bool {
	content := []byte("# T\n\n" + wrap("| zelle | zweite zelle |"))
	lines := textLines(content)
	prose := proseLineSet(content)
	for i := range lines {
		if strings.Contains(lines[i], "zelle") {
			return unbreakableTableRow(lines, prose, i)
		}
	}
	return false
}

// tableRowLineKenntZeile ist die Spalte-0-Antwort derselben Zeile.
func tableRowLineKenntZeile(wrap func(string) string) bool {
	content := []byte("# T\n\n" + wrap("| zelle | zweite zelle |"))
	lines := textLines(content)
	prose := proseLineSet(content)
	for i := range lines {
		if strings.Contains(lines[i], "zelle") {
			return tableRowLine(lines, prose, i)
		}
	}
	return false
}

// Erste Assertion der Tabellen-Kopplung (ADR-0054 Entscheidung 4, ADR-0116
// Entscheid 4): targets, planning.waves und structure beantworten die Frage
// "Tabellenzeile?" (tableRowLine) auf EINER Antwort — auch fuer die
// eingerueckte Zeile, die sie weiter nicht lesen.
func TestTabellenzeilenFrageHatEineAntwort(t *testing.T) {
	for _, c := range tabellenFaelle() {
		t.Run(c.name, func(t *testing.T) {
			antworten := map[string]bool{
				"targets":        targetsKenntZeile(t, c.wrap),
				"planning.waves": wavesKenntZeile(t, c.wrap),
				"structure":      structureKenntZeile(t, c.wrap),
				"tableRowLine":   tableRowLineKenntZeile(c.wrap),
			}
			for modul, antwort := range antworten {
				if antwort != c.gueltig {
					t.Fatalf("%s antwortet %v, gemeinsam erwartet %v (Antworten: %v)",
						modul, antwort, c.gueltig, antworten)
				}
			}
		})
	}
}

// Zweite Assertion (ADR-0116 Entscheid 4): linelength beantwortet "Teil einer
// Tabelle und deshalb nicht umbrechbar?" auf der Antwort von
// unbreakableTableRow, in beiden Richtungen — bis drei Leerzeichen davor
// ausgenommen; vier Leerzeichen, Tab, U+00A0, `>`, eine Listenmarke und der
// Fence gemessen. Der Helfer schliesst tableRowLine ein: wo tableRowLine "ja"
// sagt, sagt er es auch.
func TestTabellenzeilenFrageVonLinelengthHatEineAntwort(t *testing.T) {
	for _, c := range tabellenFaelleLinelength() {
		t.Run(c.name, func(t *testing.T) {
			antworten := map[string]bool{
				"linelength":          linelengthKenntZeile(t, c.wrap),
				"unbreakableTableRow": unbreakableTableRowKenntZeile(c.wrap),
			}
			for modul, antwort := range antworten {
				if antwort != c.gueltig {
					t.Fatalf("%s antwortet %v, erwartet %v (Antworten: %v)", modul, antwort, c.gueltig, antworten)
				}
			}
			if tableRowLineKenntZeile(c.wrap) && !unbreakableTableRowKenntZeile(c.wrap) {
				t.Fatalf("tableRowLine sagt ja, unbreakableTableRow nein: der Helfer schliesst sie nicht ein")
			}
		})
	}
}

// Die Kopplung selbst: fuer jede Schreibweise muessen ALLE Konsumenten dieselbe
// Antwort geben. Weicht einer ab, ist die Lexik an ihrem Rand gedriftet.
func TestAnkerFrageHatEineAntwort(t *testing.T) {
	for _, f := range ankerFaelle() {
		t.Run(f.name, func(t *testing.T) {
			antworten := map[string]bool{
				"anchors":  anchorsKenntAnker(t, f.ziel, f.anker),
				"versions": versionsKenntAnker(t, f.ziel, f.anker),
				"pins":     pinsKenntAnker(t, f.ziel, f.anker),
			}
			for modul, antwort := range antworten {
				if antwort != f.gueltig {
					t.Fatalf("%s antwortet %v, gemeinsam erwartet %v (Antworten: %v)",
						modul, antwort, f.gueltig, antworten)
				}
			}
		})
	}
}

// Dritte Kopplung (ADR-0054 Entscheidung 4): die Frage "liegt die Zeile
// ausserhalb von Fenced-Code?" beantwortet proseLineSet fuer targets,
// structure und planning.waves (ueber die Tabellenzeile, die nur in der
// Prosa-Menge zaehlt), planning (Ruhe-Marker), vcs (immutable-when und
// Kopf-Status-Zeile) und linelength (Ausnahme (a)). Die Faelle sind die
// Raender des Fence-Automaten: Umschalten bei jeder Zaun-Zeile, jede
// Einrueckung vor dem Zaun, keine Infozeile mit Backtick, ein Zaun hinter >
// schaltet nicht.
func fenceFaelle() []tabellenFall {
	bt, tl := "```", "~~~"
	return []tabellenFall{
		{"Backtick-Zaun", func(r string) string { return bt + "\n" + r + "\n" + bt + "\n" }, false},
		{"Tilde-Zaun", func(r string) string { return tl + "\n" + r + "\n" + tl + "\n" }, false},
		{"eingerueckter Zaun", func(r string) string { return "    " + bt + "\n" + r + "\n    " + bt + "\n" }, false},
		{"Tilde schliesst Backtick-Block", func(r string) string { return bt + "\nx\n" + tl + "\n" + r + "\n" }, true},
		{"Backtick in der Infozeile", func(r string) string { return bt + "a`b\n" + r + "\n" }, true},
		{"offener Zaun", func(r string) string { return bt + "\n" + r + "\n" }, false},
		{"Zaun im Zitat", func(r string) string { return "> " + bt + "\n\n" + r + "\n\n> " + bt + "\n" }, true},
	}
}

// linelengthSiehtProsa: mit Ausnahme (a) an und (b), (c) aus wird eine
// ueberlange Zeile genau dann gemessen, wenn sie ausserhalb von Fenced-Code
// liegt.
func linelengthSiehtProsa(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	cfg := llCfg(10)
	cfg.Linelength.ExemptTables = llOff()
	cfg.Linelength.ExemptUnbreakable = llOff()
	res := llRun(t, map[string]string{"docs/a.md": "# T\n\n" + wrap("| zelle | zweite zelle |")}, cfg)
	return len(llLines(res, "docs/a.md")) > 0
}

// planningKenntZeile: der Ruhe-Marker im Block der Überschrift zählt genau
// dann, wenn er außerhalb von Fenced-Code steht; ohne Slice meldet planning
// sonst planning-drift.
func planningKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	files := map[string]string{planRoadmap: "# Roadmap\n\n## Aktuelle Welle\n\n" + wrap("Keine aktive Welle.")}
	return len(CheckPlanning(coretest.NewMemFS(files), planCfg())) == 0
}

// vcsKenntZeile: immutable-when trifft eine Zeile genau dann, wenn sie
// außerhalb von Fenced-Code steht.
func vcsKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	return baseImmutable([]byte("# T\n\n"+wrap("Status: Accepted")), regexp.MustCompile(`^Status: Accepted`))
}

// vcsStatusKenntZeile: die Kopf-Status-Zeile vor der ersten H2 zählt genau
// dann, wenn sie außerhalb von Fenced-Code steht.
func vcsStatusKenntZeile(t *testing.T, wrap func(string) string) bool {
	t.Helper()
	content := []byte("# T\n\n" + wrap("Status: Accepted"))
	return vcsHeadStatusLineNo(content, splitLines(content), regexp.MustCompile(`^Status: Accepted`)) != 0
}

// Die Fence-Kopplung selbst: dieselbe Eingabe, eine Antwort in allen Modulen,
// die proseLineSet lesen.
func TestFencedCodeFrageHatEineAntwort(t *testing.T) {
	for _, c := range fenceFaelle() {
		t.Run(c.name, func(t *testing.T) {
			antworten := map[string]bool{
				"targets":        targetsKenntZeile(t, c.wrap),
				"structure":      structureKenntZeile(t, c.wrap),
				"planning.waves": wavesKenntZeile(t, c.wrap),
				"planning":       planningKenntZeile(t, c.wrap),
				"vcs":            vcsKenntZeile(t, c.wrap),
				"vcs.status":     vcsStatusKenntZeile(t, c.wrap),
				"linelength":     linelengthSiehtProsa(t, c.wrap),
			}
			for modul, antwort := range antworten {
				if antwort != c.gueltig {
					t.Fatalf("%s antwortet %v, gemeinsam erwartet %v (Antworten: %v)",
						modul, antwort, c.gueltig, antworten)
				}
			}
		})
	}
}
