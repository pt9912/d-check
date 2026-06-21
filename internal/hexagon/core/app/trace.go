package app

import (
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// reqIDFull erkennt eine vollständige Anforderungs-Kennung
// (`<PREFIX>-FA-<BEREICH>-NNN` bzw. `<PREFIX>-QA-NN`) — präfix-agnostisch
// (DC, AC, …), deckungsgleich mit dem ids-/matrix-Anforderungs-Muster
// (slice-036, DC-FA-CLI-009).
var reqIDFull = regexp.MustCompile(`[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?`)

// adrFileShape: `0013-foo.md` → Capture `0013` (ADR-Kennung über den
// Dateinamen). sliceFileShape: `slice-039-foo.md` → `slice-039`.
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

// Kanonische Quellen der RTM (d-check-Konvention; Doku-Domäne).
const (
	traceReqSource = "spec/lastenheft.md"
	traceADRDir    = "docs/plan/adr"
	traceSliceDir  = "docs/plan/planning"
)

// BuildTraceMatrix leitet die RTM aus den kanonischen Quellen ab
// (DC-FA-CLI-009, slice-036): Anforderungen aus dem Lastenheft, ihre
// Referenzen aus ADR- und Slice-Dateien (Doku-only). Reiner Lese-Pfad
// (DC-QA-03), deterministisch (DC-QA-02).
func BuildTraceMatrix(fsys driven.Filesystem) (TraceMatrix, error) {
	titles, order, err := traceRequirements(fsys)
	if err != nil {
		return TraceMatrix{}, err
	}
	adrRefs, err := traceRefs(fsys, traceADRDir, adrFileShape, "ADR-")
	if err != nil {
		return TraceMatrix{}, err
	}
	sliceRefs, err := traceRefs(fsys, traceSliceDir, sliceFileShape, "")
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

// traceRequirements sammelt die im Lastenheft definierten Anforderungen
// (Heading mit führender Anforderungs-Kennung) als ID→Titel plus die
// byteweise sortierte Reihenfolge.
func traceRequirements(fsys driven.Filesystem) (map[string]string, []string, error) {
	if !pathExists(fsys, traceReqSource) {
		return map[string]string{}, nil, nil
	}
	content, err := fsys.ReadFile(traceReqSource)
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
		if !isFullReqID(tok) {
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

// isFullReqID prüft, ob tok GANZ eine Anforderungs-Kennung ist (nicht nur
// als Teilstring) — verhindert, dass `DC-FA-CLI-006.a` o. Ä. zählt.
func isFullReqID(tok string) bool {
	return tok != "" && reqIDFull.FindString(tok) == tok
}

// traceTitle entfernt die führende Kennung und einen Trenner (—/-/:) aus
// dem Heading-Klartext.
func traceTitle(plain, id string) string {
	rest := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(plain), id))
	return strings.TrimSpace(strings.TrimLeft(rest, "—-:· "))
}

// traceRefs scannt die Markdown-Dateien unter dir, leitet je Datei ihre
// eigene Kennung über den Dateinamen ab (fileShape + prefix) und sammelt
// pro referenzierter Anforderungs-Kennung die referenzierenden
// Datei-Kennungen (dedupliziert, sortiert). Dateien ohne passende Kennung
// (z. B. README.md) werden übersprungen.
func traceRefs(fsys driven.Filesystem, dir string, fileShape *regexp.Regexp, prefix string) (map[string][]string, error) {
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
		for _, ref := range reqIDFull.FindAllString(string(content), -1) {
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
