package rules

// Die EINE Zell-Antwort des Produkts (DC-FA-STRUCT-001.a Schritt 6): eine
// Tabellenzeile in ihre Zellen zerlegen. Sie liegt hier, weil die übrige
// Markdown-Lexik hier liegt (markdown.go) und weil der Kern-Schnitt nur diese
// Richtung erlaubt — die header-gebundenen Tabellen-Leser in `app` rufen
// herunter, nicht umgekehrt (ADR-0012).
//
// Grenze: die Erkennung ist ZEILEN-lokal. Ein Backtick-Span über zwei Zeilen
// ist keiner; eine Tabellenzeile ist in GFM ohnehin eine Zeile.

import "strings"

// SplitPipeTableLine teilt eine Tabellenzeile an unescapten Pipes außerhalb
// korrekt geschlossener einzeiliger Backtick-Spans. Führende und abschließende
// Pipes sind optional. Zellen kommen getrimmt zurück; `\|` steht als `|` im
// Ergebnis — gemessen wird damit die Zelle, die der Leser sieht, nicht ihre
// Quell-Schreibweise. Zweiter Rückgabewert false ⇒ die Zeile trägt keinen
// Zelltrenner und ist keine Tabellenzeile.
func SplitPipeTableLine(line string) ([]string, bool) {
	s := strings.TrimSpace(line)
	if s == "" {
		return nil, false
	}
	leading := strings.HasPrefix(s, "|")
	trailing := hasUnescapedTrailingPipe(s)
	splitter := pipeLineSplitter{}
	for i := 0; i < len(s); {
		i = splitter.consume(s, i)
	}
	if splitter.separators == 0 {
		return nil, false
	}
	cells := append(splitter.cells, strings.TrimSpace(splitter.cell.String()))
	if leading && len(cells) > 0 {
		cells = cells[1:]
	}
	if trailing && len(cells) > 0 {
		cells = cells[:len(cells)-1]
	}
	return cells, true
}

type pipeLineSplitter struct {
	cells      []string
	cell       strings.Builder
	separators int
	codeRun    int
}

func (p *pipeLineSplitter) consume(s string, i int) int {
	switch {
	case s[i] == '\\' && i+1 < len(s) && s[i+1] == '|':
		p.cell.WriteByte('|')
		return i + 2
	case s[i] == '`':
		return p.consumeBackticks(s, i)
	case s[i] == '|' && p.codeRun == 0:
		p.cells = append(p.cells, strings.TrimSpace(p.cell.String()))
		p.cell.Reset()
		p.separators++
		return i + 1
	default:
		p.cell.WriteByte(s[i])
		return i + 1
	}
}

func (p *pipeLineSplitter) consumeBackticks(s string, i int) int {
	j := i
	for j < len(s) && s[j] == '`' {
		j++
	}
	run := j - i
	switch p.codeRun {
	case 0:
		if hasClosingBacktickRun(s, j, run) {
			p.codeRun = run
		}
	case run:
		p.codeRun = 0
	}
	p.cell.WriteString(s[i:j])
	return j
}

// hasClosingBacktickRun prüft vor dem Öffnen, ob auf derselben Zeile eine
// exakt gleich lange schließende Folge existiert. Ohne Abschluss ist die
// aktuelle Folge literal und nachfolgende Pipes bleiben Zelltrenner.
func hasClosingBacktickRun(s string, from, runLen int) bool {
	for i := from; i < len(s); {
		if s[i] != '`' {
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] == '`' {
			j++
		}
		if j-i == runLen {
			return true
		}
		i = j
	}
	return false
}

func hasUnescapedTrailingPipe(s string) bool {
	if !strings.HasSuffix(s, "|") {
		return false
	}
	backslashes := 0
	for i := len(s) - 2; i >= 0 && s[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 0
}
