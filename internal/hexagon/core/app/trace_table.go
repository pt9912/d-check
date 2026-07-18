package app

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// tableDelimiterCell folgt GFM: eine Trennzelle braucht **einen** Bindestrich
// (`^:?-+:?$`), optionale Doppelpunkte bleiben. Rein erweiternd gegenüber dem
// früheren `-{3,}` — jede bisher erkannte Trennzeile bleibt erkannt
// (spec/spezifikation.md §DC-FA-REQ-001.a Schritt 3, ADR-0042).
var tableDelimiterCell = regexp.MustCompile(`^:?-+:?$`)

// traceTableRequirements extrahiert Anforderungen aus allen relevanten
// Markdown-Pipe-Tabellen der Quelle (DC-FA-REQ-001). Relevanz entsteht allein
// durch die exakt konfigurierten Header-Namen; Positionen werden nie geraten.
func traceTableRequirements(fsys driven.Filesystem, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig) (map[string]string, []string, map[string]string, error) {
	content, err := fsys.ReadFile(source)
	if err != nil {
		return nil, nil, nil, err
	}
	extracted := tableExtraction{
		titles:          map[string]string{},
		modalityTexts:   map[string]string{},
		usedTextHeaders: map[string]bool{},
	}

	isRelevant := func(header []string) bool {
		_, ok, err := bindTableColumns(header, cfg)
		return err == nil && ok
	}
	for _, t := range markdownTables(content, nil, isRelevant) {
		if err := extractTable(t, source, reqPat, cfg, &extracted); err != nil {
			return nil, nil, nil, err
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

// extractTable wertet eine Tabelle aus, sofern ihre Header alle konfigurierten
// Rollen tragen. Ein Zellenzahl-Bruch in einer **relevanten** Tabelle ist
// fail-closed (DC-FA-REQ-001) — sonst risse die Tabelle still ab und
// Anforderungen verschwänden lautlos.
func extractTable(t mdTable, source string, reqPat *regexp.Regexp, cfg *model.TraceTableConfig, out *tableExtraction) error {
	cols, relevant, err := bindTableColumns(t.header, cfg)
	if err != nil {
		return fmt.Errorf("trace.requirements: Tabelle ab Zeile %d: %w", t.line, err)
	}
	if !relevant {
		return nil
	}
	out.foundTable = true
	out.usedTextHeaders[cols.textName] = true
	if t.badLine != 0 {
		return fmt.Errorf("trace.requirements: Tabellenzeile %d hat %d statt %d Zellen", t.badLine, t.badCells, len(t.header))
	}
	for _, row := range t.rows {
		addTableRequirement(row.cells, cols, source, reqPat, cfg, out)
	}
	return nil
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

// mdTableRow ist eine Datenzeile samt **Original**-Zeilennummer.
type mdTableRow struct {
	cells []string
	line  int
}

// mdTable ist eine erkannte Markdown-Pipe-Tabelle: Header-Zellen, Datenzeilen und
// die erste Pipe-Zeile mit abweichender Zellenzahl (badLine 0 ⇒ keine; der
// Konsument entscheidet, ob das ein Fehler ist — für eine nicht-relevante Tabelle
// ist es bloß das Tabellenende). Der gemeinsame Lese-Kern der header-gebundenen
// Tabellen-Konsumenten (DC-FA-REQ-001, DC-FA-XREF-001).
type mdTable struct {
	header   []string
	line     int
	rows     []mdTableRow
	badLine  int
	badCells int
}

// markdownTables liefert alle Pipe-Tabellen von content in Dokument-Reihenfolge.
// mask (nil ⇒ ganze Datei) blendet Zeilen außerhalb der gewählten Abschnitte aus
// — eine ausgeblendete Zeile beendet die laufende Tabelle wie Fließtext, und die
// Zeilennummern bleiben die der Original-Datei (rules.SectionMask). isRelevantHeader
// (nil ⇒ inaktiv) meldet, ob eine Headerzeile eine konfigurierte Rolle bindet — der
// Konsument reicht sein Relevanz-Prädikat durch, damit ein relevanter Header die
// laufende Tabelle beendet (DC-FA-REQ-001.a Schritt 5, ADR-0043).
func markdownTables(content []byte, mask []bool, isRelevantHeader func([]string) bool) []mdTable {
	lines := markdownTableLines(content)
	var out []mdTable
	for i := 0; i+1 < len(lines); i++ {
		if !maskAllows(mask, lines[i].no) || !maskAllows(mask, lines[i+1].no) {
			continue
		}
		header, ok := tableHeaderAt(lines, i)
		if !ok {
			continue
		}
		t := mdTable{header: header, line: lines[i].no}
		next := consumeTableRows(lines, i+2, mask, &t, isRelevantHeader)
		out = append(out, t)
		i = next - 1
	}
	return out
}

// consumeTableRows sammelt die Datenzeilen ab start und liefert den Index der
// ersten Zeile, die nicht mehr zur Tabelle gehört.
func consumeTableRows(lines []markdownTableLine, start int, mask []bool, t *mdTable, isRelevantHeader func([]string) bool) int {
	for j := start; j < len(lines); j++ {
		if !lines[j].prose || !maskAllows(mask, lines[j].no) {
			return j
		}
		// Grenze am relevanten Header (DC-FA-REQ-001.a Schritt 5, ADR-0043):
		// bildet diese Zeile mit ihrer Folgezeile einen gültigen Header + Trennzeile
		// und bindet ihr Header eine Rolle, ist sie der Header einer NEUEN Tabelle
		// und beendet die laufende — der Re-Scan (markdownTables, i = next-1) erkennt
		// sie ab hier. Ein rollenloser Header (z. B. all-dashes) beendet nicht. So
		// verschwindet keine relevante Tabelle mehr still in einer vorangehenden.
		if isRelevantHeader != nil && j+1 < len(lines) && maskAllows(mask, lines[j+1].no) {
			if hdr, ok := tableHeaderAt(lines, j); ok && isRelevantHeader(hdr) {
				return j
			}
		}
		cells, row := splitPipeTableLine(lines[j].text)
		if !row {
			return j
		}
		if len(cells) != len(t.header) {
			t.badLine, t.badCells = lines[j].no, len(cells)
			return j
		}
		t.rows = append(t.rows, mdTableRow{cells: cells, line: lines[j].no})
	}
	return len(lines)
}

// maskAllows prüft die 1-basierte Zeilennummer gegen die Abschnitts-Maske.
func maskAllows(mask []bool, no int) bool {
	return mask == nil || no-1 < len(mask) && mask[no-1]
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
				// Öffner-Entscheidung inkl. CommonMark-Infozeilen-Regel über das
				// GETEILTE Prädikat rules.FenceToggle — dieselbe Definition wie in
				// proseLines, damit trace und links dasselbe Dokument sehen
				// (Review R-F-1). Nur der längenabgeglichene Schluss (unten) bleibt
				// lokal (ADR-0042: naiver-Toggle-vs-strikter-Schluss bewusst offen).
				if rules.FenceToggle(trimmed) {
					fenceChar, fenceLen = char, run
					continue
				}
			} else if char == fenceChar && run >= fenceLen && strings.TrimSpace(trimmed[run:]) == "" {
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
