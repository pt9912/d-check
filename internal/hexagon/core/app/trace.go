package app

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// reqIDFull ist die **Default**-Gestalt einer Anforderungs-Kennung
// (`<PREFIX>-FA-<BEREICH>-NNN` bzw. `<PREFIX>-QA-NN`) — präfix-agnostisch
// (DC, AC, …), deckungsgleich mit dem ids-/matrix-Anforderungs-Muster
// (slice-036, DC-FA-CLI-009); per `trace.requirements.id-pattern` überschreibbar
// (slice-066).
var reqIDFull = regexp.MustCompile(`[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?`)

// adrFileShape: `0013-foo.md` → Capture `0013` (ADR-Kennung über den
// Dateinamen). sliceFileShape: `slice-039-foo.md` → `slice-039`. Die **Default**-
// Datei-Gestalten; per `trace.adrs.file-pattern`/`trace.slices.file-pattern`
// überschreibbar (slice-066).
var adrFileShape = regexp.MustCompile(`^(\d{4})-.*\.md$`)
var sliceFileShape = regexp.MustCompile(`^(slice-\d+)-.*\.md$`)

// TraceRow ist eine RTM-Zeile: eine Anforderung mit den sie referenzierenden
// ADRs/Slices, den kuratierten Coverage-Labels (DC-FA-COV-001; `omitempty` ⇒
// ohne Coverage byte-identisch) und der Waisen-Markierung (weder Slice noch
// Coverage).
type TraceRow struct {
	ID       string   `json:"id" yaml:"id"`
	Title    string   `json:"title" yaml:"title"`
	ADRs     []string `json:"adrs" yaml:"adrs"`
	Slices   []string `json:"slices" yaml:"slices"`
	Coverage []string `json:"coverage,omitempty" yaml:"coverage,omitempty"`
	// Modality: RFC-2119-Stufe (DC-FA-MOD-001; `omitempty` ⇒ ohne aktives
	// modality byte-identisch).
	Modality string `json:"modality,omitempty" yaml:"modality,omitempty"`
	Orphan   bool   `json:"orphan" yaml:"orphan"`
}

// TraceMatrix ist die vollständige Requirements Traceability Matrix,
// deterministisch sortiert (DC-QA-02).
type TraceMatrix struct {
	Requirements []TraceRow `json:"requirements" yaml:"requirements"`
	Total        int        `json:"total" yaml:"total"`
	Orphans      int        `json:"orphans" yaml:"orphans"`
	// CrossConsistency: Befunde des opt-in Kreuzverweis-Abgleichs
	// (DC-FA-XREF-001) — additiv **neben** der RTM, keine RTM-Spalte;
	// `omitempty` ⇒ ohne Block byte-identisch (DC-QA-02).
	CrossConsistency []CrossFinding `json:"crossConsistency,omitempty" yaml:"crossConsistency,omitempty"`
	// CoverageActive: ≥1 trace.coverage-Quelle konfiguriert ⇒ der Reporter
	// rendert die Coverage-Spalte (DC-FA-COV-001). Nicht serialisiert — reine
	// Reporter-Steuerung; ohne Quelle false ⇒ RTM byte-identisch.
	CoverageActive bool `json:"-" yaml:"-"`
	// ModalityActive: trace.requirements.modality präsent ⇒ Modality-Spalte
	// (DC-FA-MOD-001). GatingOrphans: Zahl der Waisen, die --require-complete
	// gaten — ohne modality == Orphans (alle), mit modality nur die
	// require-levels-Stufen. Nicht serialisiert.
	ModalityActive bool `json:"-" yaml:"-"`
	GatingOrphans  int  `json:"-" yaml:"-"`
	// CrossActive: trace.cross-consistency präsent ⇒ der Reporter rendert den
	// Abgleich-Abschnitt auch bei 0 Differenzen (Beleg statt Schweigen). Nicht
	// serialisiert — ohne Block false ⇒ RTM byte-identisch.
	CrossActive bool `json:"-" yaml:"-"`
}

// Kanonische Default-Quellen der RTM (d-check-Konvention; Doku-Domäne). Per
// trace-Block überschreibbar (slice-066, DC-FA-CLI-009).
const (
	traceReqSource = "spec/lastenheft.md"
	traceADRDir    = "docs/plan/adr"
	traceSliceDir  = "docs/plan/planning"
)

// resolvedTrace hält die aufgelösten RTM-Quellen (konfigurierter Wert über
// Konventions-Default). Der Nullwert einer model.TraceConfig ⇒ alle Defaults ⇒
// byte-identisch (DC-QA-02).
type resolvedTrace struct {
	source      string
	reqPat      *regexp.Regexp
	format      string
	table       *model.TraceTableConfig
	strictReqs  bool
	adrDir      string
	adrFile     *regexp.Regexp
	adrPrefix   string
	sliceDir    string
	sliceFile   *regexp.Regexp
	slicePrefix string
	coverage    []model.TraceCoverage
}

// resolveTrace legt die Konventions-Defaults an und überschreibt sie mit jedem
// gesetzten trace-Feld. Ein leerer ADR-Präfix bleibt bewusst Default `ADR-`
// (nicht ausdrückbar, DC-FA-CLI-009 Out-of-Scope); der Slice-Präfix ist per
// Default leer und wird direkt übernommen.
func resolveTrace(tc model.TraceConfig) resolvedTrace {
	rt := resolvedTrace{
		source: traceReqSource, reqPat: reqIDFull,
		format: model.TraceFormatHeadings,
		adrDir: traceADRDir, adrFile: adrFileShape, adrPrefix: "ADR-",
		sliceDir: traceSliceDir, sliceFile: sliceFileShape, slicePrefix: tc.SlicePrefix,
		coverage: tc.Coverage,
	}
	if tc.Source != "" {
		rt.source = tc.Source
		rt.strictReqs = true
	}
	if tc.ReqPattern != nil {
		rt.reqPat = tc.ReqPattern
	}
	if tc.Format != "" {
		rt.format = tc.Format
	}
	rt.table = tc.Table
	if rt.format == model.TraceFormatTable {
		rt.strictReqs = true
	}
	if tc.ADRDir != "" {
		rt.adrDir = tc.ADRDir
	}
	if tc.ADRFile != nil {
		rt.adrFile = tc.ADRFile
	}
	if tc.ADRPrefix != "" {
		rt.adrPrefix = tc.ADRPrefix
	}
	if tc.SliceDir != "" {
		rt.sliceDir = tc.SliceDir
	}
	if tc.SliceFile != nil {
		rt.sliceFile = tc.SliceFile
	}
	return rt
}

// BuildTraceMatrix leitet die RTM aus den (per trace-Block konfigurierbaren)
// Quellen ab (DC-FA-CLI-009, slice-036/066): Anforderungen aus dem Lastenheft,
// ihre Referenzen aus ADR- und Slice-Dateien (Doku-only). Der Nullwert von tc ⇒
// Konventions-Default ⇒ byte-identisch (DC-QA-02). Reiner Lese-Pfad (DC-QA-03).
func BuildTraceMatrix(fsys driven.Filesystem, tc model.TraceConfig) (TraceMatrix, error) {
	rt := resolveTrace(tc)
	titles, order, tableModality, err := loadTraceRequirements(fsys, rt)
	if err != nil {
		return TraceMatrix{}, err
	}
	adrRefs, err := traceRefs(fsys, rt.adrDir, rt.adrFile, rt.adrPrefix, rt.reqPat)
	if err != nil {
		return TraceMatrix{}, err
	}
	sliceRefs, err := traceRefs(fsys, rt.sliceDir, rt.sliceFile, rt.slicePrefix, rt.reqPat)
	if err != nil {
		return TraceMatrix{}, err
	}
	covRefs, err := coverageRefs(fsys, rt.coverage, rt.reqPat)
	if err != nil {
		return TraceMatrix{}, err
	}
	// Der Kreuzverweis-Abgleich läuft nach der RTM-Verrechnung und ändert die
	// RTM selbst nicht (DC-FA-XREF-001.a Schritt 1).
	cross, err := crossConsistency(fsys, tc.CrossConsistency, rt.reqPat)
	if err != nil {
		return TraceMatrix{}, err
	}
	m := TraceMatrix{
		Requirements:     make([]TraceRow, 0, len(order)),
		CoverageActive:   len(rt.coverage) > 0,
		ModalityActive:   tc.Modality != nil,
		CrossConsistency: cross,
		CrossActive:      tc.CrossConsistency != nil,
	}
	mm, modByID, err := traceModalities(fsys, tc.Modality, rt, tableModality)
	if err != nil {
		return TraceMatrix{}, err
	}
	for _, id := range order {
		row := traceRow(id, titles, adrRefs, sliceRefs, covRefs, modByID, m.ModalityActive)
		// Waise = weder Slice- noch Coverage-Referenz (DC-FA-CLI-011, slice-067);
		// gatend nur, wenn ohne modality (alle) oder die Stufe in require-levels
		// liegt (DC-FA-MOD-001, slice-068).
		if row.Orphan {
			m.Orphans++
			if !m.ModalityActive || mm.requireLevels[row.Modality] {
				m.GatingOrphans++
			}
		}
		m.Requirements = append(m.Requirements, row)
	}
	m.Total = len(m.Requirements)
	return m, nil
}

func loadTraceRequirements(fsys driven.Filesystem, rt resolvedTrace) (map[string]string, []string, map[string]string, error) {
	if rt.format != model.TraceFormatHeadings && rt.format != model.TraceFormatTable {
		return nil, nil, nil, fmt.Errorf("trace.requirements: unbekanntes Format %q", rt.format)
	}
	if rt.format == model.TraceFormatTable && rt.table == nil {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Format table braucht eine Tabellenkonfiguration")
	}
	if rt.strictReqs && !pathExists(fsys, rt.source) {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Quelle %q fehlt (Format %s)", rt.source, rt.format)
	}
	titles, order, tableModality, err := traceRequirements(fsys, rt)
	if err != nil {
		return nil, nil, nil, err
	}
	if rt.strictReqs && len(order) == 0 {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Quelle %q im Format %s ergab 0 Anforderungen", rt.source, rt.format)
	}
	return titles, order, tableModality, nil
}

func traceModalities(fsys driven.Filesystem, cfg *model.TraceModality, rt resolvedTrace, tableTexts map[string]string) (modalityMatcher, map[string]string, error) {
	if cfg == nil {
		return modalityMatcher{}, nil, nil
	}
	matcher := resolveModality(cfg)
	if rt.format == model.TraceFormatTable {
		return matcher, classifyRequirementTexts(tableTexts, matcher), nil
	}
	modalities, err := requirementModality(fsys, rt.source, rt.reqPat, matcher)
	return matcher, modalities, err
}

// traceRow baut eine RTM-Zeile aus den Referenz-Maps und bestimmt die
// Waisen-Markierung (¬slice ∧ ¬coverage); bei aktivem modality die Stufe
// (Default `unknown`).
func traceRow(id string, titles map[string]string, adrRefs, sliceRefs, covRefs map[string][]string, modByID map[string]string, modalityActive bool) TraceRow {
	row := TraceRow{ID: id, Title: titles[id], ADRs: adrRefs[id], Slices: sliceRefs[id], Coverage: covRefs[id]}
	if row.ADRs == nil {
		row.ADRs = []string{}
	}
	if row.Slices == nil {
		row.Slices = []string{}
	}
	if modalityActive {
		if row.Modality = modByID[id]; row.Modality == "" {
			row.Modality = "unknown"
		}
	}
	row.Orphan = len(row.Slices) == 0 && len(row.Coverage) == 0
	return row
}

// traceRequirements wählt den konfigurierten Extraktor und liefert Titel,
// sortierte Kennungen und (nur für table) den Modalitäts-Eingabetext.
func traceRequirements(fsys driven.Filesystem, rt resolvedTrace) (map[string]string, []string, map[string]string, error) {
	if rt.format == model.TraceFormatTable {
		return traceTableRequirements(fsys, rt.source, rt.reqPat, rt.table)
	}
	titles, order, err := traceHeadingRequirements(fsys, rt.source, rt.reqPat, rt.strictReqs)
	return titles, order, nil, err
}

// traceHeadingRequirements sammelt führende Heading-Kennungen. Im Strict-
// Modus ist eine doppelte ID ein Quelldatenfehler; im Default gewinnt wie
// bisher der erste Treffer (byte-identisch, DC-QA-02).
func traceHeadingRequirements(fsys driven.Filesystem, source string, reqPat *regexp.Regexp, strict bool) (map[string]string, []string, error) {
	if !pathExists(fsys, source) {
		return map[string]string{}, nil, nil
	}
	content, err := fsys.ReadFile(source)
	if err != nil {
		return nil, nil, err
	}
	titles := map[string]string{}
	for _, h := range rules.ExtractHeadings(content) {
		plain := rules.StripHeadingLinks(h)
		fields := strings.Fields(plain)
		if len(fields) == 0 {
			continue
		}
		tok := strings.Trim(fields[0], "`.,:;")
		if !isFullReqID(reqPat, tok) {
			continue
		}
		if _, seen := titles[tok]; seen {
			if strict {
				return nil, nil, fmt.Errorf("trace.requirements: doppelte Anforderungs-ID %q in Quelle %q", tok, source)
			}
			continue
		}
		titles[tok] = traceTitle(plain, tok)
	}
	order := make([]string, 0, len(titles))
	for id := range titles {
		order = append(order, id)
	}
	sort.Strings(order)
	return titles, order, nil
}

// classifyRequirementTexts wendet denselben Modalitätsmatcher auf die vom
// Tabellenextraktor gelieferte Spalte an (DC-FA-REQ-001/DC-FA-MOD-001).
func classifyRequirementTexts(texts map[string]string, mm modalityMatcher) map[string]string {
	res := make(map[string]string, len(texts))
	for id, body := range texts {
		res[id] = mm.classify(body)
	}
	return res
}

// isFullReqID prüft, ob tok GANZ eine Anforderungs-Kennung nach reqPat ist
// (nicht nur als Teilstring) — verhindert, dass ein Sub-ID-Suffix o. Ä. zählt.
func isFullReqID(reqPat *regexp.Regexp, tok string) bool {
	return tok != "" && reqPat.FindString(tok) == tok
}

// traceTitle entfernt die führende Kennung und einen Trenner (—/-/:) aus
// dem Heading-Klartext.
func traceTitle(plain, id string) string {
	// Optional backtick-umschlossene Kennung am Anfang entfernen — die
	// Backticks NUR als ID-Wrapper (drei gezielte TrimPrefix), nicht in der
	// abschließenden Trenner-Klasse: sonst verlöre ein Titel-initialer
	// Code-Span seinen führenden Backtick (Review R1 LOW-1 / R2 LOW-2).
	p := strings.TrimPrefix(strings.TrimSpace(plain), "`")
	p = strings.TrimPrefix(p, id)
	p = strings.TrimPrefix(p, "`")
	return strings.TrimSpace(strings.TrimLeft(p, "—-:· "))
}

// traceRefs scannt die Markdown-Dateien unter dir, leitet je Datei ihre
// eigene Kennung über den Dateinamen ab (fileShape-Capture 1 + prefix) und
// sammelt pro referenzierter Anforderungs-Kennung (reqPat) die referenzierenden
// Datei-Kennungen (dedupliziert, sortiert). Dateien ohne passende Kennung
// (z. B. README.md) werden übersprungen.
func traceRefs(fsys driven.Filesystem, dir string, fileShape *regexp.Regexp, prefix string, reqPat *regexp.Regexp) (map[string][]string, error) {
	if !pathExists(fsys, dir) {
		return map[string][]string{}, nil
	}
	files, err := rules.DiscoverFiles(fsys, []string{dir}, nil)
	if err != nil {
		return nil, err
	}
	seen := map[string]map[string]bool{}
	for _, f := range files {
		m := fileShape.FindStringSubmatch(path.Base(f))
		if m == nil {
			continue
		}
		owner := prefix + m[1]
		content, err := fsys.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, ref := range reqPat.FindAllString(string(content), -1) {
			if seen[ref] == nil {
				seen[ref] = map[string]bool{}
			}
			seen[ref][owner] = true
		}
	}
	return sortedSets(seen), nil
}

// Range-/Enum-Suffixe hinter einer Anforderungs-Kennung (DC-FA-COV-001):
// trailingDigits schneidet die Familie ab, rangeSuffix erkennt `..BBB`,
// enumSuffix `/BBB/CCC`.
var (
	trailingDigits = regexp.MustCompile(`\d+$`)
	rangeSuffix    = regexp.MustCompile(`^\.\.(\d+)`)
	enumSuffix     = regexp.MustCompile(`^(?:/\d+)+`)
	// linkSuffix erkennt ein Markdown-Link-Suffix unmittelbar hinter der Kennung:
	// optionales schließendes Code-Span-Backtick, `]`, geklammertes Ziel
	// (DC-FA-COV-001.a Link-Transparenz, ADR-0039). Das Ziel endet an der ERSTEN
	// `)` — enthält eine URL selbst eine Klammer, greift die Regel nicht und es
	// wird nicht expandiert. Das ist die sichere Richtung: lieber keine Expansion
	// als eine geratene.
	linkSuffix = regexp.MustCompile("^`?\\]\\([^)]*\\)")
)

// skipLinkSuffix überspringt **höchstens ein** Markdown-Link-Suffix hinter einer
// Kennung (DC-FA-COV-001.a Link-Transparenz, ADR-0039). Steht die Kennung unter
// Linkpflicht (DC-FA-ID-001), folgt ihre Range-/Enum-Fortsetzung nicht mehr
// unmittelbar, sondern erst hinter `](…)` — ohne dieses Überspringen bräche die
// unqualifizierte Range-Zusage des Lastenhefts genau dort, wo d-checks eigene
// Linkpflicht greift. Bewusst NUR eines und sonst nichts (kein Whitespace, keine
// Emphasis, kein zweites Suffix): jede weitere Toleranz rät die Autor-Absicht.
func skipLinkSuffix(rest string) string {
	if m := linkSuffix.FindString(rest); m != "" {
		return rest[len(m):]
	}
	return rest
}

// coverageRefs scannt die trace.coverage-Quellen (DC-FA-COV-001) und liefert je
// Anforderungs-Kennung die abdeckenden Labels (dedupliziert, sortiert).
// Fail-closed: fehlende Datei, Sektionsname ohne Heading-Treffer oder ungültige
// Range ⇒ Fehler (Exit 2). Leere Liste ⇒ nil (byte-identisch, DC-QA-02).
func coverageRefs(fsys driven.Filesystem, sources []model.TraceCoverage, reqPat *regexp.Regexp) (map[string][]string, error) {
	seen := map[string]map[string]bool{}
	for _, src := range sources {
		if err := scanCoverageSource(fsys, src, reqPat, seen); err != nil {
			return nil, err
		}
	}
	return sortedSets(seen), nil
}

// scanCoverageSource liest die Dateien einer Coverage-Quelle, prüft die
// Sektionsnamen (fail-closed) und trägt je abgedeckter Kennung das Label in seen.
func scanCoverageSource(fsys driven.Filesystem, src model.TraceCoverage, reqPat *regexp.Regexp, seen map[string]map[string]bool) error {
	contents, err := readCoverageFiles(fsys, src.Files)
	if err != nil {
		return err
	}
	if err := checkSectionNames(src, contents); err != nil {
		return err
	}
	for _, c := range contents {
		text := string(rules.SelectSections(c, src.Sections, src.ExcludeSections))
		ids, err := rangeAwareIDs("trace.coverage", text, reqPat, src.Ranges)
		if err != nil {
			return err
		}
		for _, id := range ids {
			if seen[id] == nil {
				seen[id] = map[string]bool{}
			}
			seen[id][src.Label] = true
		}
	}
	return nil
}

// readCoverageFiles liest die explizit benannten Coverage-Dateien; eine fehlende
// Datei ⇒ Fehler (fail-closed, DC-FA-COV-001 — anders als adrs/slices).
func readCoverageFiles(fsys driven.Filesystem, files []string) ([][]byte, error) {
	contents := make([][]byte, 0, len(files))
	for _, f := range files {
		if !pathExists(fsys, f) {
			return nil, fmt.Errorf("trace.coverage: Datei %q fehlt", f)
		}
		c, err := fsys.ReadFile(f)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

// sortedSets wandelt map[key]set in map[key]sortierte-Liste (deterministisch).
func sortedSets(seen map[string]map[string]bool) map[string][]string {
	res := make(map[string][]string, len(seen))
	for k, set := range seen {
		list := make([]string, 0, len(set))
		for v := range set {
			list = append(list, v)
		}
		sort.Strings(list)
		res[k] = list
	}
	return res
}

// checkSectionNames stellt fail-closed sicher, dass jeder konfigurierte
// Sektionsname (sections/exclude-sections) über die Quell-Dateien ≥1 Überschrift
// trifft — Tippfehler-/Kurzform-Guard (DC-FA-COV-001). Sonst blankte eine
// Whitelist ohne Treffer still die ganze Datei bzw. griffe eine Blacklist nicht.
func checkSectionNames(src model.TraceCoverage, contents [][]byte) error {
	if len(src.Sections) == 0 && len(src.ExcludeSections) == 0 {
		return nil
	}
	headings := map[string]bool{}
	for _, c := range contents {
		for _, h := range rules.HeadingTexts(c) {
			headings[h] = true
		}
	}
	for _, name := range append(append([]string{}, src.Sections...), src.ExcludeSections...) {
		if !headings[name] {
			return fmt.Errorf("trace.coverage %q: Abschnitt %q trifft keine Überschrift (voller Heading-Text erwartet)", src.Label, name)
		}
	}
	return nil
}

// rangeAwareIDs extrahiert die Anforderungs-Kennungen aus text: exakte
// reqPat-Treffer plus (bei ranges) die range-/enum-expandierten IDs
// (DC-FA-COV-001). field benennt den Config-Schlüssel für Fehlermeldungen — die
// Range-Semantik teilen sich die Coverage-Quellen und beide Sichten des
// Kreuzverweis-Abgleichs (DC-FA-XREF-001).
func rangeAwareIDs(field, text string, reqPat *regexp.Regexp, ranges bool) ([]string, error) {
	out := append([]string{}, reqPat.FindAllString(text, -1)...)
	if !ranges {
		return out, nil
	}
	for _, loc := range reqPat.FindAllStringIndex(text, -1) {
		exp, err := expandRange(field, text[loc[0]:loc[1]], text[loc[1]:], reqPat)
		if err != nil {
			return nil, err
		}
		out = append(out, exp...)
	}
	return out, nil
}

// expandRange expandiert das Range-/Enum-Suffix unmittelbar hinter der Kennung id
// (rest = Text danach, DC-FA-COV-001). Familie = id ohne Trailing-Ziffern;
// `..BBB` breiten-erhaltend inklusiv, `/BBB/CCC` als Aufzählung; jede expandierte
// ID gegen reqPat geprüft (Nicht-Treffer verworfen). Fail-closed: AAA>BBB oder
// abweichende Ziffern-Breite ⇒ Fehler.
func expandRange(field, id, rest string, reqPat *regexp.Regexp) ([]string, error) {
	d := trailingDigits.FindStringIndex(id)
	if d == nil {
		return nil, nil // Kennung endet nicht auf Ziffern → keine Range
	}
	fam, startStr := id[:d[0]], id[d[0]:]
	width := len(startStr)
	rest = skipLinkSuffix(rest)
	if rm := rangeSuffix.FindStringSubmatch(rest); rm != nil {
		endStr := rm[1]
		if len(endStr) != width {
			return nil, fmt.Errorf("%s: Range %s..%s mit abweichender Ziffern-Breite", field, id, endStr)
		}
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		if start > end {
			return nil, fmt.Errorf("%s: Range %s..%s mit AAA>BBB", field, id, endStr)
		}
		var out []string
		for i := start; i <= end; i++ {
			if cand := fmt.Sprintf("%s%0*d", fam, width, i); reqPat.FindString(cand) == cand {
				out = append(out, cand)
			}
		}
		return out, nil
	}
	if em := enumSuffix.FindString(rest); em != "" {
		var out []string
		for _, part := range strings.Split(strings.TrimPrefix(em, "/"), "/") {
			if cand := fam + part; reqPat.FindString(cand) == cand {
				out = append(out, cand)
			}
		}
		return out, nil
	}
	return nil, nil
}

// mdEmphasis matcht Markdown-Emphasis-Zeichen (Sternchen/Backtick) für die
// Body-Normalisierung der Modalitäts-Klassifikation (DC-FA-MOD-001).
var mdEmphasis = regexp.MustCompile("[*`]+")

// modalityMatcher hält den kompilierten Keyword-Matcher (längster Treffer via
// longest-first-Alternation), die Rückabbildung Keyword→Stufe und die gatenden
// Stufen (DC-FA-MOD-001).
type modalityMatcher struct {
	re            *regexp.Regexp
	kwToLevel     map[string]string // lowercased Keyword → Stufe
	requireLevels map[string]bool
}

// resolveModality baut den Matcher aus der (ggf. Default-)Keyword-Menge:
// Alternation längster-zuerst (deterministisch, bei Gleichstand byte-sortiert),
// `(?i)\b…\b`; require-levels Default `[must]`.
func resolveModality(m *model.TraceModality) modalityMatcher {
	levels := m.Levels
	if len(levels) == 0 {
		levels = model.DefaultModalityLevels()
	}
	kwToLevel := map[string]string{}
	var kws []string
	for level, list := range levels {
		for _, kw := range list {
			key := strings.ToLower(strings.TrimSpace(kw))
			if key == "" {
				continue
			}
			kwToLevel[key] = level
			kws = append(kws, strings.TrimSpace(kw))
		}
	}
	sort.Slice(kws, func(i, j int) bool {
		if len(kws[i]) != len(kws[j]) {
			return len(kws[i]) > len(kws[j])
		}
		return kws[i] < kws[j]
	})
	quoted := make([]string, len(kws))
	for i, kw := range kws {
		quoted[i] = regexp.QuoteMeta(kw)
	}
	req := m.RequireLevels
	if len(req) == 0 {
		req = []string{"must"}
	}
	rl := map[string]bool{}
	for _, r := range req {
		rl[r] = true
	}
	mm := modalityMatcher{kwToLevel: kwToLevel, requireLevels: rl}
	if len(quoted) > 0 {
		mm.re = regexp.MustCompile(`(?i)\b(?:` + strings.Join(quoted, "|") + `)\b`)
	}
	return mm
}

// classify liefert die Stufe der Anforderung anhand des ersten (frühesten)/
// längsten Keyword-Treffers im normalisierten Body; kein Treffer ⇒ "unknown".
func (mm modalityMatcher) classify(body string) string {
	if mm.re == nil {
		return "unknown"
	}
	if hit := mm.re.FindString(normalizeBody(body)); hit != "" {
		if level, ok := mm.kwToLevel[strings.ToLower(hit)]; ok {
			return level
		}
	}
	return "unknown"
}

// normalizeBody entfernt Markdown-Emphasis und zieht Whitespace/Umbrüche zu je
// einem Leerzeichen zusammen (DC-FA-MOD-001.a Schritt 2) — sonst matchte
// `**MUSS** NICHT` oder `MUSS\nNICHT` die Phrase nicht.
func normalizeBody(body string) string {
	return strings.Join(strings.Fields(mdEmphasis.ReplaceAllString(body, "")), " ")
}

// requirementModality klassifiziert je Anforderung (source) ihre Stufe über den
// Body-Abschnitt (rules.HeadingSections); nur Headings, deren erstes Token eine
// Anforderungs-Kennung (reqPat) ist (DC-FA-MOD-001).
func requirementModality(fsys driven.Filesystem, source string, reqPat *regexp.Regexp, mm modalityMatcher) (map[string]string, error) {
	res := map[string]string{}
	if !pathExists(fsys, source) {
		return res, nil
	}
	content, err := fsys.ReadFile(source)
	if err != nil {
		return nil, err
	}
	for _, hs := range rules.HeadingSections(content) {
		fields := strings.Fields(rules.StripHeadingLinks(hs.Text))
		if len(fields) == 0 {
			continue
		}
		tok := strings.Trim(fields[0], "`.,:;")
		if !isFullReqID(reqPat, tok) {
			continue
		}
		if _, seen := res[tok]; !seen {
			res[tok] = mm.classify(hs.Body)
		}
	}
	return res, nil
}
