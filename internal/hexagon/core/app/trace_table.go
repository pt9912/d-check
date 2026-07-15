package app

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

var tableDelimiterCell = regexp.MustCompile(`^:?-{3,}:?$`)

// traceTableRequirements extrahiert Anforderungen aus allen relevanten
// Markdown-Pipe-Tabellen der Quelle (DC-FA-REQ-001). Relevanz entsteht allein
// durch die exakt konfigurierten Header-Namen; Positionen werden nie geraten.
func traceTableRequirements(fsys driven.Filesystem, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig) (map[string]string, []string, map[string]string, error) {
	content, err := fsys.ReadFile(source)
	if err != nil {
		return nil, nil, nil, err
	}
	lines := markdownTableLines(content)
	extracted := tableExtraction{
		titles:          map[string]string{},
		modalityTexts:   map[string]string{},
		usedTextHeaders: map[string]bool{},
	}

	for i := 0; i+1 < len(lines); i++ {
		next, found, err := extractTableAt(lines, i, source, reqPat, cfg, &extracted)
		if err != nil {
			return nil, nil, nil, err
		}
		if found {
			i = next
		}
	}

	if !extracted.foundTable {
		return nil, nil, nil, fmt.Errorf("trace.requirements: keine Tabelle mit den konfigurierten Headern %s in %q", configuredTableHeaders(cfg), source)
	}
	for _, name := range cfg.TextColumns {
		if !extracted.usedTextHeaders[name] {
			return nil, nil, nil, fmt.Errorf("trace.requirements: konfigurierter Text-Header %q kommt in keiner Tabelle in %q vor", name, source)
		}
	}
	if extracted.dupErr != nil {
		return nil, nil, nil, extracted.dupErr
	}
	order := make([]string, 0, len(extracted.titles))
	for id := range extracted.titles {
		order = append(order, id)
	}
	sort.Strings(order)
	return extracted.titles, order, extracted.modalityTexts, nil
}

type tableExtraction struct {
	titles          map[string]string
	modalityTexts   map[string]string
	usedTextHeaders map[string]bool
	foundTable      bool
	dupErr          error
}

func extractTableAt(lines []markdownTableLine, i int, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig, out *tableExtraction) (int, bool, error) {
	header, ok := tableHeaderAt(lines, i)
	if !ok {
		return i, false, nil
	}
	cols, relevant, err := bindTableColumns(header, cfg)
	if err != nil {
		return i, false, fmt.Errorf("trace.requirements: Tabelle ab Zeile %d: %w", lines[i].no, err)
	}
	if relevant {
		out.foundTable = true
		out.usedTextHeaders[cols.textName] = true
	}
	next, err := consumeTableRows(lines, i+2, len(header), cols, relevant, source, reqPat, cfg, out)
	return next - 1, true, err
}

func tableHeaderAt(lines []markdownTableLine, i int) ([]string, bool) {
	if !lines[i].prose || !lines[i+1].prose {
		return nil, false
	}
	header, ok := splitPipeTableLine(lines[i].text)
	if !ok {
		return nil, false
	}
	delimiter, ok := splitPipeTableLine(lines[i+1].text)
	if !ok || len(delimiter) != len(header) || !isTableDelimiter(delimiter) {
		return nil, false
	}
	return header, true
}

func consumeTableRows(lines []markdownTableLine, start, width int, cols tableColumns, relevant bool, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig, out *tableExtraction) (int, error) {
	for j := start; j < len(lines); j++ {
		if !lines[j].prose || strings.TrimSpace(lines[j].text) == "" {
			return j, nil
		}
		cells, row := splitPipeTableLine(lines[j].text)
		if !row {
			return j, nil
		}
		if len(cells) != width {
			if relevant {
				return j, fmt.Errorf("trace.requirements: Tabellenzeile %d hat %d statt %d Zellen", lines[j].no, len(cells), width)
			}
			return j, nil
		}
		if relevant {
			addTableRequirement(cells, cols, source, reqPat, cfg, out)
		}
	}
	return len(lines), nil
}

func addTableRequirement(cells []string, cols tableColumns, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig, out *tableExtraction) {
	id := strings.TrimSpace(cells[cols.id])
	if !isFullReqID(reqPat, id) {
		return
	}
	if _, duplicate := out.titles[id]; duplicate {
		switch cfg.DuplicateIDs {
		case model.TraceDuplicateFirst:
			return
		case model.TraceDuplicateLast:
			// Überschreiben ist der explizite Brownfield-Vertrag.
		default:
			// Fehler zurückstellen, damit Header-/Struktur-Präzedenz (unbenutzte
			// Text-Alternative) vor der Duplicate-ID-Meldung greift
			// (DC-FA-REQ-001.a). Bis dahin gewinnt die erste ID (kein Überschreiben).
			if out.dupErr == nil {
				out.dupErr = fmt.Errorf("trace.requirements: doppelte Anforderungs-ID %q in Tabellenquelle %q", id, source)
			}
			return
		}
	}
	out.titles[id] = strings.TrimSpace(cells[cols.text])
	modalityColumn := cols.text
	if cols.modality >= 0 {
		modalityColumn = cols.modality
	}
	out.modalityTexts[id] = strings.TrimSpace(cells[modalityColumn])
}

type tableColumns struct {
	id       int
	text     int
	modality int
	textName string
}

// bindTableColumns prüft, ob header alle konfigurierten Namen enthält. Erst
// dann ist die Tabelle relevant; doppelte Rollen-Header sind fail-closed.
func bindTableColumns(header []string, cfg *model.TraceTableConfig) (tableColumns, bool, error) {
	cols := tableColumns{id: -1, text: -1, modality: -1}
	counts := map[string]int{}
	indices := map[string]int{}
	for i, cell := range header {
		counts[cell]++
		if _, exists := indices[cell]; !exists {
			indices[cell] = i
		}
	}
	if counts[cfg.IDColumn] == 0 || cfg.ModalityColumn != "" && counts[cfg.ModalityColumn] == 0 {
		return cols, false, nil
	}
	for _, name := range []string{cfg.IDColumn, cfg.ModalityColumn} {
		if name != "" && counts[name] > 1 {
			return cols, false, fmt.Errorf("konfigurierter Header %q kommt mehrfach vor", name)
		}
	}
	textName := ""
	for _, name := range cfg.TextColumns {
		if counts[name] > 1 {
			return cols, false, fmt.Errorf("konfigurierter Header %q kommt mehrfach vor", name)
		}
		if counts[name] == 1 {
			if textName != "" {
				return cols, false, fmt.Errorf("alternative Text-Header %q und %q kommen in derselben Tabelle vor", textName, name)
			}
			textName = name
		}
	}
	if textName == "" {
		return cols, false, nil
	}
	cols.id, cols.text = indices[cfg.IDColumn], indices[textName]
	cols.textName = textName
	if cfg.ModalityColumn != "" {
		cols.modality = indices[cfg.ModalityColumn]
	}
	return cols, true, nil
}

func configuredTableHeaders(cfg *model.TraceTableConfig) string {
	textHeaders := make([]string, 0, len(cfg.TextColumns))
	for _, name := range cfg.TextColumns {
		textHeaders = append(textHeaders, fmt.Sprintf("%q", name))
	}
	headers := []string{fmt.Sprintf("%q", cfg.IDColumn), "einer von " + strings.Join(textHeaders, ", ")}
	if cfg.ModalityColumn != "" {
		headers = append(headers, fmt.Sprintf("%q", cfg.ModalityColumn))
	}
	return strings.Join(headers, ", ")
}

func isTableDelimiter(cells []string) bool {
	for _, cell := range cells {
		if !tableDelimiterCell.MatchString(strings.TrimSpace(cell)) {
			return false
		}
	}
	return len(cells) > 0
}

type markdownTableLine struct {
	no    int
	text  string
	prose bool
}

// markdownTableLines markiert Zeilen außerhalb von ```/~~~-Fences. Zeilen
// bleiben positionsgetreu erhalten, damit ein Fence niemals Header und
// Trennzeile künstlich benachbart macht.
func markdownTableLines(content []byte) []markdownTableLine {
	rawLines := strings.Split(string(content), "\n")
	out := make([]markdownTableLine, len(rawLines))
	var fenceChar byte
	fenceLen := 0
	for i, raw := range rawLines {
		out[i] = markdownTableLine{no: i + 1, text: raw}
		trimmed := strings.TrimLeft(raw, " \t")
		if char, run := fenceMarker(trimmed); run >= 3 {
			if fenceLen == 0 {
				fenceChar, fenceLen = char, run
				continue
			}
			if char == fenceChar && run >= fenceLen && strings.TrimSpace(trimmed[run:]) == "" {
				fenceChar, fenceLen = 0, 0
				continue
			}
		}
		out[i].prose = fenceLen == 0
	}
	return out
}

func fenceMarker(line string) (byte, int) {
	if line == "" || line[0] != '`' && line[0] != '~' {
		return 0, 0
	}
	char := line[0]
	i := 0
	for i < len(line) && line[i] == char {
		i++
	}
	return char, i
}

// splitPipeTableLine teilt an unescaped Pipes außerhalb korrekt geschlossener
// einzeiliger Backtick-Spans. Führende/abschließende Pipes sind optional.
func splitPipeTableLine(line string) ([]string, bool) {
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
