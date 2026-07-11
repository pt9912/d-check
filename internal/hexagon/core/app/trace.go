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
	Orphan   bool     `json:"orphan" yaml:"orphan"`
}

// TraceMatrix ist die vollständige Requirements Traceability Matrix,
// deterministisch sortiert (DC-QA-02).
type TraceMatrix struct {
	Requirements []TraceRow `json:"requirements" yaml:"requirements"`
	Total        int        `json:"total" yaml:"total"`
	Orphans      int        `json:"orphans" yaml:"orphans"`
	// CoverageActive: ≥1 trace.coverage-Quelle konfiguriert ⇒ der Reporter
	// rendert die Coverage-Spalte (DC-FA-COV-001). Nicht serialisiert — reine
	// Reporter-Steuerung; ohne Quelle false ⇒ RTM byte-identisch.
	CoverageActive bool `json:"-" yaml:"-"`
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
		adrDir: traceADRDir, adrFile: adrFileShape, adrPrefix: "ADR-",
		sliceDir: traceSliceDir, sliceFile: sliceFileShape, slicePrefix: tc.SlicePrefix,
		coverage: tc.Coverage,
	}
	if tc.Source != "" {
		rt.source = tc.Source
	}
	if tc.ReqPattern != nil {
		rt.reqPat = tc.ReqPattern
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
	titles, order, err := traceRequirements(fsys, rt.source, rt.reqPat)
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
	m := TraceMatrix{
		Requirements:   make([]TraceRow, 0, len(order)),
		CoverageActive: len(rt.coverage) > 0,
	}
	for _, id := range order {
		row := TraceRow{ID: id, Title: titles[id], ADRs: adrRefs[id], Slices: sliceRefs[id], Coverage: covRefs[id]}
		if row.ADRs == nil {
			row.ADRs = []string{}
		}
		if row.Slices == nil {
			row.Slices = []string{}
		}
		// Waise = weder Slice- noch Coverage-Referenz (DC-FA-CLI-011, slice-067).
		if len(row.Slices) == 0 && len(row.Coverage) == 0 {
			row.Orphan = true
			m.Orphans++
		}
		m.Requirements = append(m.Requirements, row)
	}
	m.Total = len(m.Requirements)
	return m, nil
}

// traceRequirements sammelt die im Anforderungs-Quelldokument (source)
// definierten Anforderungen (Heading mit führender Anforderungs-Kennung nach
// reqPat) als ID→Titel plus die byteweise sortierte Reihenfolge.
func traceRequirements(fsys driven.Filesystem, source string, reqPat *regexp.Regexp) (map[string]string, []string, error) {
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
		if _, seen := titles[tok]; !seen {
			titles[tok] = traceTitle(plain, tok)
		}
	}
	order := make([]string, 0, len(titles))
	for id := range titles {
		order = append(order, id)
	}
	sort.Strings(order)
	return titles, order, nil
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
)

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
		ids, err := coverageIDs(text, reqPat, src.Ranges)
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

// coverageIDs extrahiert die abgedeckten Anforderungs-Kennungen aus text: exakte
// reqPat-Treffer plus (bei ranges) die range-/enum-expandierten IDs (DC-FA-COV-001).
func coverageIDs(text string, reqPat *regexp.Regexp, ranges bool) ([]string, error) {
	out := append([]string{}, reqPat.FindAllString(text, -1)...)
	if !ranges {
		return out, nil
	}
	for _, loc := range reqPat.FindAllStringIndex(text, -1) {
		exp, err := expandRange(text[loc[0]:loc[1]], text[loc[1]:], reqPat)
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
func expandRange(id, rest string, reqPat *regexp.Regexp) ([]string, error) {
	d := trailingDigits.FindStringIndex(id)
	if d == nil {
		return nil, nil // Kennung endet nicht auf Ziffern → keine Range
	}
	fam, startStr := id[:d[0]], id[d[0]:]
	width := len(startStr)
	if rm := rangeSuffix.FindStringSubmatch(rest); rm != nil {
		endStr := rm[1]
		if len(endStr) != width {
			return nil, fmt.Errorf("trace.coverage: Range %s..%s mit abweichender Ziffern-Breite", id, endStr)
		}
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		if start > end {
			return nil, fmt.Errorf("trace.coverage: Range %s..%s mit AAA>BBB", id, endStr)
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
