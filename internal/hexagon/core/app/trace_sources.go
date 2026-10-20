package app

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// reqEntry ist eine aufgelöste Anforderungsquelle (DC-FA-REQ-001.a Schritt 8).
// index ist -1 in der Objekt-Form: ihre Meldungen tragen keinen Index und
// bleiben byte-identisch (DC-QA-02).
type reqEntry struct {
	index    int
	source   string
	reqPat   *regexp.Regexp
	format   string
	table    *model.TraceTableConfig
	strict   bool
	modality *model.TraceModality
}

// newReqEntry legt die Konventions-Defaults an und überschreibt sie mit jedem
// gesetzten Feld; strict ist die Absicht „nichtleer explizite source oder
// format: table" (DC-FA-REQ-001.a Schritt 6), je Eintrag gebildet.
func newReqEntry(index int, s model.TraceSource) reqEntry {
	e := reqEntry{
		index: index, source: traceReqSource, reqPat: reqIDFull,
		format: model.TraceFormatHeadings, table: s.Table, modality: s.Modality,
	}
	if s.Source != "" {
		e.source = s.Source
		e.strict = true
	}
	if s.ReqPattern != nil {
		e.reqPat = s.ReqPattern
	}
	if s.Format != "" {
		e.format = s.Format
	}
	if e.format == model.TraceFormatTable {
		e.strict = true
	}
	return e
}

// resolveEntries bildet die Objekt-Form als Liste mit einem Eintrag ab.
func resolveEntries(tc model.TraceConfig) []reqEntry {
	if len(tc.Sources) == 0 {
		return []reqEntry{newReqEntry(-1, model.TraceSource{
			Source: tc.Source, ReqPattern: tc.ReqPattern, Format: tc.Format,
			Table: tc.Table, Modality: tc.Modality,
		})}
	}
	out := make([]reqEntry, len(tc.Sources))
	for i, s := range tc.Sources {
		out[i] = newReqEntry(i, s)
	}
	return out
}

// referencePattern ist das Muster der Referenzseite: bei einem Eintrag dessen
// Muster selbst, sonst die Vereinigung aller Muster mit längstem Treffer —
// eine Kennung gilt als erkannt, wenn sie auf das Muster mindestens eines
// Eintrags passt (DC-FA-REQ-001.a Schritt 8).
func referencePattern(entries []reqEntry) (*regexp.Regexp, error) {
	if len(entries) == 1 {
		return entries[0].reqPat, nil
	}
	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = "(?:" + e.reqPat.String() + ")"
	}
	re, err := regexp.Compile(strings.Join(parts, "|"))
	if err != nil {
		return nil, fmt.Errorf("trace.requirements: die id-pattern der Einträge lassen sich nicht vereinigen: %v", err)
	}
	re.Longest()
	return re, nil
}

// indexed setzt den Index des Eintrags in die Meldung; in der Objekt-Form
// bleibt sie unverändert.
func (e reqEntry) indexed(err error) error {
	if e.index < 0 || err == nil {
		return err
	}
	label := fmt.Sprintf("trace.requirements[%d]", e.index)
	msg := err.Error()
	if strings.Contains(msg, "trace.requirements") {
		return fmt.Errorf("%s", strings.Replace(msg, "trace.requirements", label, 1))
	}
	return fmt.Errorf("%s: %s", label, msg)
}

// loadedReqs ist die Vereinigung der Anforderungen aller Einträge: Titel,
// Kennungen in Listenreihenfolge (innerhalb eines Eintrags sortiert), der
// Eigentümer-Eintrag je Kennung und der Modalitäts-Eingabetext der Tabellen.
type loadedReqs struct {
	titles     map[string]string
	order      []string
	owner      map[string]int
	tableTexts map[string]string
}

func loadTraceRequirements(fsys driven.Filesystem, entries []reqEntry) (loadedReqs, error) {
	out := loadedReqs{
		titles: map[string]string{}, owner: map[string]int{}, tableTexts: map[string]string{},
	}
	for i, e := range entries {
		titles, order, texts, err := loadEntryRequirements(fsys, e)
		if err != nil {
			return loadedReqs{}, e.indexed(err)
		}
		for _, id := range order {
			if prev, dup := out.owner[id]; dup {
				return loadedReqs{}, fmt.Errorf(
					"trace.requirements: Kennung %q kommt in Eintrag %d und in Eintrag %d vor (doppelt definiert oder überlappende id-pattern)",
					id, prev, i)
			}
			out.owner[id] = i
			out.titles[id] = titles[id]
			if t, ok := texts[id]; ok {
				out.tableTexts[id] = t
			}
		}
		out.order = append(out.order, order...)
	}
	return out, nil
}

func loadEntryRequirements(fsys driven.Filesystem, e reqEntry) (map[string]string, []string, map[string]string, error) {
	if e.format != model.TraceFormatHeadings && e.format != model.TraceFormatTable {
		return nil, nil, nil, fmt.Errorf("trace.requirements: unbekanntes Format %q", e.format)
	}
	if e.format == model.TraceFormatTable && e.table == nil {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Format table braucht eine Tabellenkonfiguration")
	}
	if e.strict && !pathExists(fsys, e.source) {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Quelle %q fehlt (Format %s)", e.source, e.format)
	}
	titles, order, texts, err := entryRequirements(fsys, e)
	if err != nil {
		return nil, nil, nil, err
	}
	if e.strict && len(order) == 0 {
		return nil, nil, nil, fmt.Errorf("trace.requirements: Quelle %q im Format %s ergab 0 Anforderungen", e.source, e.format)
	}
	return titles, order, texts, nil
}

// entryRequirements wählt den konfigurierten Extraktor und liefert Titel,
// sortierte Kennungen und (nur für table) den Modalitäts-Eingabetext.
func entryRequirements(fsys driven.Filesystem, e reqEntry) (map[string]string, []string, map[string]string, error) {
	if e.format == model.TraceFormatTable {
		return traceTableRequirements(fsys, e.source, e.reqPat, e.table)
	}
	titles, order, err := traceHeadingRequirements(fsys, e.source, e.reqPat, e.strict)
	return titles, order, nil, err
}

// entryModality ist die Klassifikation eines Eintrags. Ein Eintrag ohne
// modality klassifiziert nicht (active == false): seine Waisen gaten alle
// (DC-FA-REQ-001.a Schritt 8).
type entryModality struct {
	active bool
	mm     modalityMatcher
	byID   map[string]string
}

func traceModalities(fsys driven.Filesystem, entries []reqEntry, reqs loadedReqs) ([]entryModality, error) {
	out := make([]entryModality, len(entries))
	for i, e := range entries {
		if e.modality == nil {
			continue
		}
		mm := resolveModality(e.modality)
		var byID map[string]string
		if e.format == model.TraceFormatTable {
			byID = classifyRequirementTexts(reqs.tableTexts, mm)
		} else {
			var err error
			if byID, err = requirementModality(fsys, e.source, e.reqPat, mm); err != nil {
				return nil, e.indexed(err)
			}
		}
		out[i] = entryModality{active: true, mm: mm, byID: byID}
	}
	return out, nil
}
