package app

import (
	"path"
	"regexp"
	"sort"
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

// TraceRow ist eine RTM-Zeile: eine Anforderung mit den sie
// referenzierenden ADRs/Slices und der Waisen-Markierung (kein Slice).
type TraceRow struct {
	ID     string   `json:"id" yaml:"id"`
	Title  string   `json:"title" yaml:"title"`
	ADRs   []string `json:"adrs" yaml:"adrs"`
	Slices []string `json:"slices" yaml:"slices"`
	Orphan bool     `json:"orphan" yaml:"orphan"`
}

// TraceMatrix ist die vollständige Requirements Traceability Matrix,
// deterministisch sortiert (DC-QA-02).
type TraceMatrix struct {
	Requirements []TraceRow `json:"requirements" yaml:"requirements"`
	Total        int        `json:"total" yaml:"total"`
	Orphans      int        `json:"orphans" yaml:"orphans"`
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
	m := TraceMatrix{Requirements: make([]TraceRow, 0, len(order))}
	for _, id := range order {
		row := TraceRow{ID: id, Title: titles[id], ADRs: adrRefs[id], Slices: sliceRefs[id]}
		if row.ADRs == nil {
			row.ADRs = []string{}
		}
		if len(row.Slices) == 0 {
			row.Slices = []string{}
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
	res := make(map[string][]string, len(seen))
	for req, owners := range seen {
		list := make([]string, 0, len(owners))
		for o := range owners {
			list = append(list, o)
		}
		sort.Strings(list)
		res[req] = list
	}
	return res, nil
}
