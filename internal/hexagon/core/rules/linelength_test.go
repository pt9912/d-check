package rules

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// llCfg baut eine Konfiguration mit der Grenze limit und den Defaults der
// drei Ausnahmen.
func llCfg(limit int) model.Config {
	m := limit
	return model.Config{Linelength: model.LinelengthConfig{MaxChars: &m}}
}

func llOff() *bool { b := false; return &b }

// llRun fährt den Lauf mit dem Modul linelength allein.
func llRun(t *testing.T, files map[string]string, cfg model.Config) Result {
	t.Helper()
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"linelength"})
	if err != nil {
		t.Fatalf("unerwarteter Eingabefehler: %v", err)
	}
	return res
}

// llLines liefert die Zeilen der Befunde einer Datei.
func llLines(res Result, file string) []int {
	var out []int
	for _, f := range res.Findings {
		if f.Reason == ReasonLineCharsExceeded && f.File == file {
			out = append(out, f.Line)
		}
	}
	return out
}

// llWords baut eine Zeile aus Wörtern der Länge vier, getrennt durch ein
// Space, mit genau n Codepunkten; kein Wort ist länger als vier.
func llWords(n int) string {
	var b strings.Builder
	for b.Len() < n {
		if b.Len()%5 == 4 {
			b.WriteByte(' ')
		} else {
			b.WriteByte('a')
		}
	}
	return b.String()
}

func llWord(n int) string { return strings.Repeat("w", n) }

func llWantLines(t *testing.T, res Result, file string, want ...int) {
	t.Helper()
	got := llLines(res, file)
	if len(want) == 0 && len(got) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: Befund-Zeilen %v, erwartet %v (Befunde: %+v)", file, got, want, res.Findings)
	}
}

func llNote(t *testing.T, res Result) string {
	t.Helper()
	for _, n := range res.Notes {
		if strings.HasPrefix(n, "linelength: ") {
			return n
		}
	}
	t.Fatalf("Bezugsmengen-Zeile fehlt: %v", res.Notes)
	return ""
}

// Eine Zeile genau auf der Grenze ist befundfrei, eine darüber ein Befund mit
// Länge, Grenze und Ziel (DC-FA-LLEN-001).
func TestLinelengthGrenzeGenauAufDerSchwelle(t *testing.T) {
	res := llRun(t, map[string]string{"docs/a.md": llWords(80) + "\n" + llWords(81) + "\n"}, llCfg(80))
	llWantLines(t, res, "docs/a.md", 2)
	f := res.Findings[0]
	if f.Rule != "linelength" || f.Target != "max-chars=80" || f.Message != "Zeile hat 81 Zeichen, erlaubt sind 80" {
		t.Fatalf("Befund: %+v", f)
	}
}

// Ein einziges Wort mit max-chars + 1 Zeichen nimmt (c) im Default aus; mit
// exempt-unbreakable: false ist es ein Befund (DC-FA-LLEN-001).
func TestLinelengthEinWortUeberDerGrenze(t *testing.T) {
	files := map[string]string{"docs/a.md": llWord(81) + "\n"}
	res := llRun(t, files, llCfg(80))
	llWantLines(t, res, "docs/a.md")
	if n := llNote(t, res); !strings.HasSuffix(n, "1 davon mit einem Wort über der Grenze ausgenommen") {
		t.Fatalf("Bezugsmenge: %q", n)
	}
	cfg := llCfg(80)
	cfg.Linelength.ExemptUnbreakable = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 1)
}

// Gezählt wird in Codepunkten, nicht in Bytes; ein Tab ist ein Zeichen, eine
// ungültige UTF-8-Folge zählt je Byte (DC-FA-LLEN-001).
func TestLinelengthZaehltCodepunkte(t *testing.T) {
	// 20 Wörter aus "äö😀" (3 Codepunkte, 8 Bytes) mit Space: 80 Codepunkte.
	multi := strings.Repeat("äö😀 ", 20)
	if got := len([]rune(multi)); got != 80 {
		t.Fatalf("Fixture trägt %d Codepunkte", got)
	}
	tabs := strings.Repeat("abc\t", 20)
	invalid := llWords(78) + "\xff\xfe"
	files := map[string]string{"docs/a.md": strings.Join([]string{
		multi,            // 1: 80 Codepunkte, 180 Bytes
		multi + "x",      // 2: 81
		tabs,             // 3: 80 mit Tabs
		tabs + "x",       // 4: 81
		invalid,          // 5: 80, zwei ungültige Bytes
		invalid + "\xff", // 6: 81
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 2, 4, 6)
}

// Ein \r am Zeilenende gehört nicht zur Zeile, auch am Ende einer letzten
// Zeile ohne \n; ein \r an anderer Stelle ist ein Zeichen (DC-FA-LLEN-001).
func TestLinelengthZeilenendeOhneCR(t *testing.T) {
	files := map[string]string{
		"docs/crlf.md": llWords(80) + "\r\n" + llWords(79) + "\r" + "x\r\n",
		"docs/last.md": llWords(80) + "\r",
	}
	res := llRun(t, files, llCfg(80))
	llWantLines(t, res, "docs/crlf.md", 2)
	llWantLines(t, res, "docs/last.md")
}

// Nur ein einzelnes \r am Zeilenende fällt weg: vor \r\r\n und am Ende \r\r
// bleibt das erste \r ein Zeichen der Zeile (DC-FA-LLEN-001).
func TestLinelengthNurEinCRFaelltWeg(t *testing.T) {
	files := map[string]string{
		"docs/a.md": llWords(80) + "\r\r\n" + llWords(79) + "\r\r\n" + llWords(80) + "\r\r",
	}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 1, 3)
}

// Ausnahme (a): Zeilen in Backtick-, Tilde- und eingerücktem Zaun samt der
// Zaun-Zeile sind im Default ausgenommen, mit exempt-code-blocks: false je ein
// Befund (DC-FA-LLEN-001).
func TestLinelengthAusnahmeFencedCode(t *testing.T) {
	long := llWords(90)
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"# T",
		"```" + llWords(90), // 2: Zaun-Zeile über der Grenze
		long,                // 3
		"```",
		"~~~",
		long, // 6
		"~~~",
		"  ```",
		long, // 9
		"  ```",
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md")
	cfg := llCfg(80)
	cfg.Linelength.ExemptCodeBlocks = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 2, 3, 6, 9)
}

// Grenzen der geteilten Lexik: eine ~~~-Zeile schließt einen Backtick-Block
// früh, ein eingerückter Code-Block ohne Zaun und ein Zaun hinter > sind kein
// Fenced-Code (DC-FA-LLEN-001).
func TestLinelengthFenceGrenzenDerLexik(t *testing.T) {
	long := llWords(90)
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"```",
		long, // 2: im Block
		"~~~",
		long, // 4: nach dem frühen Schluss gemessen
		"",
		"    " + long, // 6: eingerückter Code-Block
		"",
		"> ```",
		"> " + long, // 9: Zaun im Zitat schaltet nicht
		"> ```",
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 4, 6, 9)
}

// Ein nie geschlossener Zaun nimmt den Rest der Datei aus; linelength meldet
// ihn nicht (DC-FA-LLEN-001).
func TestLinelengthOffenerZaunNimmtDenRestAus(t *testing.T) {
	long := llWords(90)
	files := map[string]string{"docs/a.md": long + "\n```\n" + long + "\n" + long + "\n"}
	res := llRun(t, files, llCfg(80))
	llWantLines(t, res, "docs/a.md", 1)
	for _, f := range res.Findings {
		if f.Reason == model.ReasonFenceUnclosed {
			t.Fatalf("linelength meldet den offenen Zaun: %+v", f)
		}
	}
	if n := llNote(t, res); !strings.Contains(n, "1 Zeile(n) geprüft") {
		t.Fatalf("Bezugsmenge zählt Zeilen hinter dem offenen Zaun: %q", n)
	}
}

// Ausnahme (b): eine |-Zeile ohne Trennzeile ist im Default ausgenommen, mit
// exempt-tables: false ein Befund (DC-FA-LLEN-001).
func TestLinelengthAusnahmeTabellenzeile(t *testing.T) {
	files := map[string]string{"docs/a.md": "| " + llWords(90) + " |\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md")
	cfg := llCfg(80)
	cfg.Linelength.ExemptTables = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 1)
}

// llRow liefert eine überlange Tabellenzeile ohne Wort über der Grenze 80.
func llRow() string { return "| " + llWords(90) + " |" }

// Ausnahme (b), eingerückt: eine |-Zeile mit einem, zwei und drei Leerzeichen
// (U+0020) davor ist im Default ausgenommen, mit exempt-tables: false je ein
// Befund (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthEingerueckteTabellenzeile(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		" " + llRow(),
		"  " + llRow(),
		"   " + llRow(),
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md")
	cfg := llCfg(80)
	cfg.Linelength.ExemptTables = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 1, 2, 3)
}

// Grenze 3 gegen 4: drei Leerzeichen vor | sind ausgenommen, vier der Befund;
// jeder Fall steht allein, damit nicht nur eine Seite der Bedingung variiert
// (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabellenzeileDreiGegenVierLeerzeichen(t *testing.T) {
	cases := []struct {
		name string
		line string
		want []int
	}{
		{"drei Leerzeichen", "   " + llRow(), nil},
		{"vier Leerzeichen", "    " + llRow(), []int{1}},
		{"fuenf Leerzeichen", "     " + llRow(), []int{1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			llWantLines(t, llRun(t, map[string]string{"docs/a.md": c.line + "\n"}, llCfg(80)), "docs/a.md", c.want...)
		})
	}
}

// Grenzen von (b) je Fall allein: Tab, geschütztes Leerzeichen (U+00A0), ein
// Leerzeichen und dann ein Tab, Zitat-Zeile, GFM-Zeile ohne führende Pipe und
// Text mit Einrückung ohne | sind keine Tabellenzeilen und werden gemessen
// (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabellenzeileGrenzen(t *testing.T) {
	cases := []struct{ name, line string }{
		{"Tab", "\t" + llRow()},
		{"geschuetztes Leerzeichen", " " + llRow()},
		{"Leerzeichen dann Tab", " \t" + llRow()},
		{"drei Leerzeichen dann Tab", "   \t" + llRow()},
		{"Zitat", "> " + llRow()},
		{"Zitat mit Einrueckung", "  > " + llRow()},
		{"ohne fuehrende Pipe", "a | " + llWords(90)},
		{"Einrueckung ohne Pipe", "   " + llWords(90)},
		{"Pipe nicht erstes Zeichen", "   x" + llRow()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			llWantLines(t, llRun(t, map[string]string{"docs/a.md": c.line + "\n"}, llCfg(80)), "docs/a.md", 1)
		})
	}
}

// Grenze 1: eine Prosa-Zeile, die nach drei Leerzeichen mit | beginnt, ohne
// dass eine Tabelle dasteht, ist ausgenommen, wie in Spalte 0 — eine
// Trennzeile oder eine zweite |-Zeile wird nicht verlangt
// (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthProsaZeileMitPipeOhneTabelleIstAusgenommen(t *testing.T) {
	prose := "| " + llWords(90)
	for _, line := range []string{prose, "   " + prose} {
		files := map[string]string{"docs/a.md": "Ein Satz.\n\n" + line + "\n\nNoch ein Satz.\n"}
		llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md")
	}
}

// Grenze 2: eine Tabelle auf der Zeile der Listenmarke beginnt nicht mit
// Leerzeichen und |, auch nicht unter einer Marke von vier Spalten; sie wird
// gemessen (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabellenzeileHinterListenmarkeWirdGemessen(t *testing.T) {
	cases := []struct{ name, line string }{
		{"Bindestrich", "- " + llRow()},
		{"Sternchen", "* " + llRow()},
		{"nummeriert", "1. " + llRow()},
		{"zweistellig nummeriert", "10. " + llRow()},
		{"Marke mit Einrueckung", "  - " + llRow()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			llWantLines(t, llRun(t, map[string]string{"docs/a.md": c.line + "\n"}, llCfg(80)), "docs/a.md", 1)
		})
	}
}

// Die Form des Anlasses: eine Tabelle unter einem nummerierten Listenpunkt
// mit drei Leerzeichen davor ist im Default befundfrei, mit exempt-tables:
// false ein Befund je überlanger Zeile, und der Listenpunkt selbst bleibt
// gemessen (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabelleUnterNummeriertemListenpunkt(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"1. Tabelle:",
		"",
		"   | a | b |",
		"   |---|---|",
		"   " + llRow(),
		"",
		"2. " + llWords(90),
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 7)
	cfg := llCfg(80)
	cfg.Linelength.ExemptTables = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 5, 7)
}

// Die Schalter sind unabhängig: mit (a) aus und (b) an liegt eine |-Zeile im
// Zaun — auch eine eingerückte, auch mit drei Leerzeichen — nicht außerhalb
// von Fenced-Code und wird gemessen (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabellenzeileImZaunWirdGemessen(t *testing.T) {
	cfg := llCfg(80)
	cfg.Linelength.ExemptCodeBlocks = llOff()
	for _, c := range []struct{ name, line string }{
		{"Spalte 0", llRow()},
		{"drei Leerzeichen", "   " + llRow()},
		{"ein Leerzeichen", " " + llRow()},
	} {
		t.Run(c.name, func(t *testing.T) {
			files := map[string]string{"docs/b.md": "```\n" + c.line + "\n```\n"}
			llWantLines(t, llRun(t, files, cfg), "docs/b.md", 2)
		})
	}
}

// Ausnahme (b), Rest unverändert: Dateien ohne eingerückte |-Zeile verhalten
// sich wie zuvor — Kopf-, Trenn- und Datenzeilen in Spalte 0 sind
// ausgenommen, die Zeilen darum gemessen, mit (b) aus alle gemessen
// (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthTabellenzeileSpalteNullUnveraendert(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"# T",
		llWords(90),
		llRow(),
		"|" + strings.Repeat(" --- |", 20),
		llRow(),
		llWords(90),
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 2, 6)
	cfg := llCfg(80)
	cfg.Linelength.ExemptTables = llOff()
	llWantLines(t, llRun(t, files, cfg), "docs/a.md", 2, 3, 4, 5, 6)
}

// Bezugsmenge: eine Zeile mit drei Leerzeichen vor | zählt im Default nicht
// zu den geprüften Zeilen, mit exempt-tables: false zählt sie
// (DC-FA-LLEN-001, ADR-0116).
func TestLinelengthBezugsmengeEingerueckteTabellenzeile(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n   | a | b |\nkurz\n"}
	if n := llNote(t, llRun(t, files, llCfg(80))); n != "linelength: 1 Datei(en), 2 Zeile(n) geprüft, 0 davon mit einem Wort über der Grenze ausgenommen" {
		t.Fatalf("Bezugsmenge im Default: %q", n)
	}
	cfg := llCfg(80)
	cfg.Linelength.ExemptTables = llOff()
	if n := llNote(t, llRun(t, files, cfg)); n != "linelength: 1 Datei(en), 3 Zeile(n) geprüft, 0 davon mit einem Wort über der Grenze ausgenommen" {
		t.Fatalf("Bezugsmenge mit exempt-tables: false: %q", n)
	}
}

// Eine eingerückte Tabellenzeile mit einem Wort über der Grenze nimmt (b)
// aus, bevor (c) sie sieht: sie zählt nicht in U (ADR-0116, Konsequenzen).
func TestLinelengthEingerueckteTabellenzeileZaehltNichtInU(t *testing.T) {
	files := map[string]string{"docs/a.md": "   | " + llWord(90) + " |\n"}
	n := llNote(t, llRun(t, files, llCfg(80)))
	if n != "linelength: 1 Datei(en), 0 Zeile(n) geprüft, 0 davon mit einem Wort über der Grenze ausgenommen" {
		t.Fatalf("Bezugsmenge: %q", n)
	}
}

// Ausnahme (c) prüft das längste Wort: ein Wort über der Grenze nimmt die
// ganze Zeile aus, auch mit umbrechbarem Text davor; ein Wort genau auf der
// Grenze oder ein kurzer Überhang ist ein Befund (DC-FA-LLEN-001).
func TestLinelengthAusnahmeWortUeberDerGrenze(t *testing.T) {
	url := "https://example.org/" + llWord(70) // 90 Zeichen, ein Wort
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"Siehe " + url,            // 1: ausgenommen
		llWord(80) + " x",         // 2: längstes Wort genau 80 ⇒ Befund
		llWord(81) + " x",         // 3: längstes Wort 81 ⇒ ausgenommen
		llWords(78) + " abcdefgh", // 4: kurzes letztes Wort ragt über ⇒ Befund
		url + " und mehr Text",    // 5: Wort vor umbrechbarem Text ⇒ ausgenommen
	}, "\n") + "\n"}
	res := llRun(t, files, llCfg(80))
	llWantLines(t, res, "docs/a.md", 2, 4)
	if n := llNote(t, res); n != "linelength: 1 Datei(en), 5 Zeile(n) geprüft, 3 davon mit einem Wort über der Grenze ausgenommen" {
		t.Fatalf("Bezugsmenge: %q", n)
	}
}

// Leerraum sind Space und Tab: ein Tab trennt Wörter, ein geschütztes
// Leerzeichen nicht (DC-FA-LLEN-001).
func TestLinelengthWortTrennerSindSpaceUndTab(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		llWord(50) + "\t" + llWord(50),     // 1: zwei Wörter ⇒ Befund
		llWord(50) + "\u00a0" + llWord(50), // 2: ein Wort mit 101 ⇒ ausgenommen
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 1)
}

// Ein Wort, das nur mit seiner Einrückung über die Grenze reicht, und eine
// überlange ATX-Überschrift sind Befunde (DC-FA-LLEN-001).
func TestLinelengthEinrueckungUndUeberschrift(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"  " + llWord(80),    // 1: Wort 80, Zeile 82
		"# " + llWords(90),   // 2: Überschrift
		"#### " + llWord(80), // 3: Überschrift, Wort genau auf der Grenze
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(80)), "docs/a.md", 1, 2, 3)
}

// Ein Präfix der Markdown-Form schiebt das erste Wort über die Grenze: die
// Zeile ist ein Befund, (c) liest die Form nicht (DC-FA-LLEN-001).
func TestLinelengthFormPraefix(t *testing.T) {
	files := map[string]string{"docs/a.md": strings.Join([]string{
		"- " + llWord(119),
		"> " + llWord(119),
		"- foo " + llWord(119),
		"-   " + llWord(117),
	}, "\n") + "\n"}
	llWantLines(t, llRun(t, files, llCfg(120)), "docs/a.md", 1, 2, 3, 4)
}

// Kein Zeilen-Marker: d-check:ignore wirkt nicht; bei einer Grenze unter 14
// ist das Marker-Wort selbst ein Wort über der Grenze (DC-FA-LLEN-001).
func TestLinelengthKeinZeilenMarker(t *testing.T) {
	marked := map[string]string{"docs/a.md": llWords(80) + " <!-- d-check:ignore -->\n"}
	llWantLines(t, llRun(t, marked, llCfg(80)), "docs/a.md", 1)
	llWantLines(t, llRun(t, marked, llCfg(14)), "docs/a.md", 1)

	small := map[string]string{"docs/b.md": "ab cd ef <!-- d-check:ignore -->\n"}
	res := llRun(t, small, llCfg(10))
	llWantLines(t, res, "docs/b.md")
	if n := llNote(t, res); !strings.HasSuffix(n, " 1 davon mit einem Wort über der Grenze ausgenommen") {
		t.Fatalf("Bezugsmenge: %q", n)
	}
}

// Bezugsmenge: Zeilen wie bei file, mit und ohne abschließendes \n; die
// Tabellenzeile zählt nicht als geprüft, die über (c) ausgenommene als eigene
// Zahl (DC-FA-LLEN-001).
func TestLinelengthBezugsmenge(t *testing.T) {
	body := strings.Join([]string{"# T", "| a | b |", llWord(81), "kurz", "noch kurz"}, "\n")
	want := "linelength: 1 Datei(en), 4 Zeile(n) geprüft, 1 davon mit einem Wort über der Grenze ausgenommen"
	for _, content := range []string{body, body + "\n"} {
		res := llRun(t, map[string]string{"docs/a.md": content}, llCfg(80))
		if n := llNote(t, res); n != want {
			t.Fatalf("%q: Bezugsmenge %q", content, n)
		}
		if countLines([]byte(content)) != 5 {
			t.Fatalf("file zählt %d Zeilen", countLines([]byte(content)))
		}
	}
}

// Eine leere Markdown-Datei hat keine Zeile; die Bezugsmenge nennt sie mit
// 0 Zeilen, wie file (DC-FA-LLEN-001).
func TestLinelengthLeereDatei(t *testing.T) {
	res := llRun(t, map[string]string{"docs/a.md": ""}, llCfg(80))
	if len(res.Findings) != 0 {
		t.Fatalf("Befunde: %+v", res.Findings)
	}
	if n := llNote(t, res); n != "linelength: 1 Datei(en), 0 Zeile(n) geprüft, 0 davon mit einem Wort über der Grenze ausgenommen" {
		t.Fatalf("Bezugsmenge: %q", n)
	}
	if countLines(nil) != 0 {
		t.Fatal("file zählt die leere Datei anders")
	}
}

// Eingabemenge: eine Datei unter exempt-paths, eine außerhalb von scope und
// eine Nicht-Markdown-Datei werden nicht gemessen (DC-FA-LLEN-001,
// DC-FA-CONF-002).
func TestLinelengthEingabemenge(t *testing.T) {
	long := llWords(90) + "\n"
	files := map[string]string{
		"docs/a.md":       long,
		"docs/frei/b.md":  long,
		"spec/c.md":       long,
		"docs/notiz.txt":  long,
		"docs/gross.MD":   long,
		"docs/x.markdown": long,
	}
	cfg := llCfg(80)
	cfg.Linelength.ExemptPaths = []string{"docs/frei/**"}
	cfg.Scopes = map[string]*model.ScopeConfig{"linelength": {Roots: []string{"docs"}}}
	res := llRun(t, files, cfg)
	if len(res.Findings) != 1 || res.Findings[0].File != "docs/a.md" {
		t.Fatalf("Befunde: %+v", res.Findings)
	}
	if n := llNote(t, res); !strings.HasPrefix(n, "linelength: 1 Datei(en), 1 Zeile(n)") {
		t.Fatalf("Bezugsmenge: %q", n)
	}
}

// Eine leere Menge ist kein Fehler: ohne Markdown-Datei und mit explizit
// leerer scope.roots kein Befund und die Bezugsmenge 0 (DC-FA-LLEN-001).
func TestLinelengthLeereMenge(t *testing.T) {
	zero := "linelength: 0 Datei(en), 0 Zeile(n) geprüft, 0 davon mit einem Wort über der Grenze ausgenommen"
	res := llRun(t, map[string]string{"docs/notiz.txt": llWords(90) + "\n"}, llCfg(80))
	if len(res.Findings) != 0 || llNote(t, res) != zero {
		t.Fatalf("ohne Markdown: %+v %v", res.Findings, res.Notes)
	}
	cfg := llCfg(80)
	cfg.Scopes = map[string]*model.ScopeConfig{"linelength": {Roots: []string{}}}
	res = llRun(t, map[string]string{"docs/a.md": llWords(90) + "\n"}, cfg)
	if len(res.Findings) != 0 || llNote(t, res) != zero {
		t.Fatalf("leere scope.roots: %+v %v", res.Findings, res.Notes)
	}
}

// Ein aktives Modul ohne max-chars ist ein Eingabefehler vor dem Lauf; ohne
// aktives Modul bleibt derselbe Block inert, ohne Zeile in den Hinweisen
// (DC-FA-LLEN-001).
func TestLinelengthOhneGrenze(t *testing.T) {
	files := map[string]string{"docs/a.md": "[x](fehlt.md)\n"}
	cfg := model.Config{Linelength: model.LinelengthConfig{ExemptTables: llOff()}}
	_, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links", "linelength"})
	if err == nil || !strings.Contains(err.Error(), "linelength.max-chars") {
		t.Fatalf("aktiv ohne Grenze: err=%v", err)
	}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links"})
	if err != nil || len(res.Notes) != 0 || len(res.Findings) != 1 {
		t.Fatalf("inaktiv: err=%v notes=%v findings=%+v", err, res.Notes, res.Findings)
	}
}

// Ohne aktives Modul ist ein Lauf mit gültigem linelength-Block gleich dem
// Lauf ohne Block (DC-FA-LLEN-001).
func TestLinelengthInaktivIstByteIdentisch(t *testing.T) {
	files := map[string]string{"docs/a.md": llWords(90) + "\n[x](fehlt.md)\n"}
	without, err := Run(coretest.NewMemFS(files), nil, model.Config{}, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	with, err := Run(coretest.NewMemFS(files), nil, llCfg(10), []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(without, with) {
		t.Fatalf("ohne Block %+v, mit Block %+v", without, with)
	}
}

// Befunde zweier Dateien stehen nach Pfad, dann Zeile; derselbe Lauf liefert
// dasselbe Ergebnis (DC-FA-LLEN-001, DC-QA-02).
func TestLinelengthReihenfolgeUndDeterminismus(t *testing.T) {
	long := llWords(90)
	files := map[string]string{
		"docs/b.md": "kurz\n" + long + "\n" + long + "\n",
		"docs/a.md": long + "\nkurz\n" + long + "\n",
	}
	first := llRun(t, files, llCfg(80))
	var got []string
	for _, f := range first.Findings {
		got = append(got, f.File+":"+strconv.Itoa(f.Line))
	}
	want := []string{"docs/a.md:1", "docs/a.md:3", "docs/b.md:2", "docs/b.md:3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Reihenfolge %v, erwartet %v", got, want)
	}
	for i := 0; i < 10; i++ {
		if again := llRun(t, files, llCfg(80)); !reflect.DeepEqual(first, again) {
			t.Fatalf("Lauf %d weicht ab: %+v", i, again)
		}
	}
}

// Die Zeilen von pinsync und linelength (textLines) zählen wie file
// (countLines), über die Randformen: leer, nur Zeilenende, abschließendes und
// fehlendes \n, \r am Ende, doppeltes \r (DC-FA-LLEN-001, DC-FA-FILE-001).
func TestTextLinesZaehltWieFile(t *testing.T) {
	for _, c := range []string{"", "\n", "a", "a\n", "a\n\n", "\n\n", "a\r", "a\r\n", "a\r\r\n", "a\nb", "a\nb\n", "\r"} {
		if got, want := len(textLines([]byte(c))), countLines([]byte(c)); got != want {
			t.Errorf("%q: textLines %d Zeile(n), countLines %d", c, got, want)
		}
	}
	if got := textLines([]byte("a\r\r\nb\r")); !reflect.DeepEqual(got, []string{"a\r", "b"}) {
		t.Errorf("Zerlegung %q", got)
	}
}
