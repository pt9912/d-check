package app

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Richtungslabel der Kreuzverweis-Befunde (DC-FA-XREF-001).
const (
	// CrossDirForwardOnly: das Artefakt steht in der Vorwärts-RTM, aber keine
	// Rück-Kante nennt die Anforderung.
	CrossDirForwardOnly = "in RTM, ohne Rück-Kante"
	// CrossDirBackwardOnly: eine Rück-Kante nennt die Anforderung, die
	// Vorwärts-RTM kennt das Artefakt nicht.
	CrossDirBackwardOnly = "Rück-Kante, ohne RTM-Eintrag"
)

// Config-Schlüssel als Fehler-/Meldungs-Präfix beider Sichten.
const (
	crossForwardField  = "trace.cross-consistency.forward"
	crossBackwardField = "trace.cross-consistency.backward"
)

// CrossFinding ist **eine** Mengendifferenz zwischen den beiden Sichten
// (DC-FA-XREF-001): Artifact ist das Design-Artefakt, das nur einer Sicht bekannt
// ist; File/Line zeigen auf die Quell-Zeile der Sicht, die es nennt (die andere
// Sicht hat für dieses Artefakt naturgemäß keine Zeile).
type CrossFinding struct {
	Requirement string `json:"requirement" yaml:"requirement"`
	Artifact    string `json:"artifact" yaml:"artifact"`
	Direction   string `json:"direction" yaml:"direction"`
	File        string `json:"file" yaml:"file"`
	Line        int    `json:"line" yaml:"line"`
}

// crossEdge ist die Fundstelle einer Anforderung→Artefakt-Kante.
type crossEdge struct {
	file string
	line int
}

// crossView ist eine Sicht der Relation: je Anforderung die genannten Artefakte
// mit ihrer Fundstelle. Beide Sichten münden in dieselbe Struktur — genau das
// macht den Diff zu einem reinen Mengenvergleich.
type crossView map[string]map[string]crossEdge

// add nimmt eine Kante auf; die **erste** Fundstelle gewinnt (deterministisch,
// DC-QA-02).
func (v crossView) add(req, artifact string, e crossEdge) {
	if v[req] == nil {
		v[req] = map[string]crossEdge{}
	}
	if _, seen := v[req][artifact]; !seen {
		v[req][artifact] = e
	}
}

// crossConsistency führt den opt-in Kreuzverweis-Abgleich aus (DC-FA-XREF-001.a):
// beide Quellen lesen, Sichten extrahieren, invertieren, Mengen diffen. cc nil ⇒
// kein Abgleich (RTM byte-identisch, DC-QA-02). reqPat ist das
// `trace.requirements.id-pattern` — es erkennt die Anforderungen der
// Vorwärts-ID-Spalte, während die Rück-Sicht ihr eigenes `req-pattern` für die
// Kanten-Zelle mitbringt (freier Prosa-Text statt kuratierter ID-Spalte).
// Die Fehlerpräzedenz ist die des Vertrags: Quellen lesen → Header-Bindung →
// Range-Expansion → Diff. Reiner Lese-Pfad (DC-QA-03).
func crossConsistency(fsys driven.Filesystem, cc *model.TraceCrossConsistency, reqPat *regexp.Regexp) ([]CrossFinding, error) {
	if cc == nil {
		return nil, nil
	}
	fwdContent, err := readCrossFile(fsys, crossForwardField, cc.Forward.File)
	if err != nil {
		return nil, err
	}
	bwdContent, err := readCrossFile(fsys, crossBackwardField, cc.Backward.File)
	if err != nil {
		return nil, err
	}
	fwd, err := forwardView(fwdContent, cc.Forward, reqPat)
	if err != nil {
		return nil, err
	}
	bwd, err := backwardView(bwdContent, cc.Backward, cc.Forward.DesignPattern)
	if err != nil {
		return nil, err
	}
	return diffViews(fwd, bwd, cc), nil
}

// readCrossFile liest eine Sicht-Quelle; eine fehlende Datei ist fail-closed
// (Exit 2) — eine stillschweigend leere Sicht meldete sonst jede Kante der
// Gegenseite als Differenz.
func readCrossFile(fsys driven.Filesystem, field, file string) ([]byte, error) {
	if !pathExists(fsys, file) {
		return nil, fmt.Errorf("%s.file: Datei %q fehlt", field, file)
	}
	return fsys.ReadFile(file)
}

// forwardView liest die Vorwärts-Sicht: je Anforderung der ID-Spalte (range-aware)
// alle Design-Artefakte der Design-Zelle (DC-FA-XREF-001.a Schritt 2).
func forwardView(content []byte, fc model.TraceCrossForward, reqPat *regexp.Regexp) (crossView, error) {
	if err := checkCrossSections(crossForwardField, content, fc.Sections, fc.ExcludeSections); err != nil {
		return nil, err
	}
	mask := rules.SectionMask(content, fc.Sections, fc.ExcludeSections)
	view := crossView{}
	found := false
	for _, t := range markdownTables(content, mask) {
		idx, relevant, err := bindCrossColumns(t.header, fc.ReqColumn, fc.DesignColumn)
		if err != nil {
			return nil, fmt.Errorf("%s: Tabelle ab Zeile %d: %w", crossForwardField, t.line, err)
		}
		if !relevant {
			continue
		}
		found = true
		if err := crossBadRow(crossForwardField, t); err != nil {
			return nil, err
		}
		if err := forwardRows(t, idx[0], idx[1], fc, reqPat, view); err != nil {
			return nil, err
		}
	}
	if !found {
		return nil, fmt.Errorf("%s: keine Tabelle mit den konfigurierten Headern %q, %q in %q",
			crossForwardField, fc.ReqColumn, fc.DesignColumn, fc.File)
	}
	return view, nil
}

// forwardRows trägt die Kanten einer relevanten Vorwärts-Tabelle ein: je
// Anforderung der ID-Zelle (range-aware) jedes Artefakt der Design-Zelle.
func forwardRows(t mdTable, reqIdx, designIdx int, fc model.TraceCrossForward, reqPat *regexp.Regexp, view crossView) error {
	for _, row := range t.rows {
		reqs, err := rangeAwareIDs(crossForwardField, row.cells[reqIdx], reqPat, fc.Ranges)
		if err != nil {
			return err
		}
		for _, req := range reqs {
			for _, artifact := range fc.DesignPattern.FindAllString(row.cells[designIdx], -1) {
				view.add(req, artifact, crossEdge{file: fc.File, line: row.line})
			}
		}
	}
	return nil
}

// backwardView liest die Rück-Kanten und **invertiert** sie: je in der
// Kanten-Zelle genannter Anforderung wird die Artefakt-ID aufgenommen
// (DC-FA-XREF-001.a Schritt 3). Die Artefakt-ID kommt bewusst über
// designPat — dasselbe Muster wie die Vorwärts-Sicht, denn nur ein
// gemeinsamer Namensraum macht den Mengen-Diff bedeutungsvoll.
func backwardView(content []byte, bc model.TraceCrossBackward, designPat *regexp.Regexp) (crossView, error) {
	if err := checkCrossSections(crossBackwardField, content, bc.Sections, nil); err != nil {
		return nil, err
	}
	mask := rules.SectionMask(content, bc.Sections, nil)
	view := crossView{}
	found := false
	for _, t := range markdownTables(content, mask) {
		idIdx, edgeIdx, relevant, err := bindBackwardColumns(t.header, bc)
		if err != nil {
			return nil, fmt.Errorf("%s: Tabelle ab Zeile %d: %w", crossBackwardField, t.line, err)
		}
		if !relevant {
			continue
		}
		found = true
		if err := crossBadRow(crossBackwardField, t); err != nil {
			return nil, err
		}
		if err := backwardRows(t, idIdx, edgeIdx, bc, designPat, view); err != nil {
			return nil, err
		}
	}
	if !found {
		return nil, fmt.Errorf("%s: keine Tabelle mit dem konfigurierten Header %q in %q",
			crossBackwardField, bc.EdgeColumn, bc.File)
	}
	return view, nil
}

// backwardRows trägt die invertierten Kanten einer relevanten Rück-Tabelle ein:
// je Anforderung der Kanten-Zelle (range-aware) die Artefakt-ID der ID-Spalte.
func backwardRows(t mdTable, idIdx, edgeIdx int, bc model.TraceCrossBackward, designPat *regexp.Regexp, view crossView) error {
	for _, row := range t.rows {
		artifact := designPat.FindString(row.cells[idIdx])
		if artifact == "" {
			continue // Zeile ohne Artefakt-Kennung (Zwischenüberschrift o. Ä.)
		}
		reqs, err := rangeAwareIDs(crossBackwardField, row.cells[edgeIdx], bc.ReqPattern, bc.Ranges)
		if err != nil {
			return err
		}
		for _, req := range reqs {
			view.add(req, artifact, crossEdge{file: bc.File, line: row.line})
		}
	}
	return nil
}

// crossBadRow macht einen Zellenzahl-Bruch in einer relevanten Tabelle
// fail-closed — dieselbe Reader-Semantik wie DC-FA-REQ-001: eine still
// abreißende Tabelle verlöre Kanten und erzeugte Phantom-Differenzen.
func crossBadRow(field string, t mdTable) error {
	if t.badLine == 0 {
		return nil
	}
	return fmt.Errorf("%s: Tabellenzeile %d hat %d statt %d Zellen", field, t.badLine, t.badCells, len(t.header))
}

// bindCrossColumns bindet die konfigurierten Header an ihre Spaltenindizes.
// Relevanz entsteht — wie in DC-FA-REQ-001 — allein durch die Präsenz **aller**
// Namen; ein mehrfach vorkommender Rollen-Header ist fail-closed (die Bindung
// wäre sonst geraten).
func bindCrossColumns(header []string, names ...string) ([]int, bool, error) {
	counts := map[string]int{}
	indices := map[string]int{}
	for i, cell := range header {
		counts[cell]++
		if _, exists := indices[cell]; !exists {
			indices[cell] = i
		}
	}
	for _, n := range names {
		if counts[n] == 0 {
			return nil, false, nil
		}
	}
	idx := make([]int, 0, len(names))
	for _, n := range names {
		if counts[n] > 1 {
			return nil, false, fmt.Errorf("konfigurierter Header %q kommt mehrfach vor", n)
		}
		idx = append(idx, indices[n])
	}
	return idx, true, nil
}

// bindBackwardColumns bindet die Kanten-Spalte über ihren Header und die
// Artefakt-ID-Spalte entweder positionell (Sentinel `first`, da deren Header über
// die Tabellen variiert — ADR-0038) oder ebenfalls über den Header-Namen.
func bindBackwardColumns(header []string, bc model.TraceCrossBackward) (int, int, bool, error) {
	names := []string{bc.EdgeColumn}
	if bc.ArtifactIDColumn != model.TraceCrossArtifactFirst {
		names = append(names, bc.ArtifactIDColumn)
	}
	idx, relevant, err := bindCrossColumns(header, names...)
	if err != nil || !relevant {
		return 0, 0, false, err
	}
	idIdx := 0
	if len(idx) > 1 {
		idIdx = idx[1]
	}
	return idIdx, idx[0], true, nil
}

// checkCrossSections spiegelt den fail-closed Sektionsnamen-Guard der
// Coverage-Quellen (DC-FA-COV-001): ein Name ohne Heading-Treffer (Tippfehler,
// Kurzform statt vollem Klartext) blankte sonst still die ganze Sicht.
func checkCrossSections(field string, content []byte, include, exclude []string) error {
	names := append(append([]string{}, include...), exclude...)
	if len(names) == 0 {
		return nil
	}
	headings := map[string]bool{}
	for _, h := range rules.HeadingTexts(content) {
		headings[h] = true
	}
	for _, name := range names {
		if !headings[name] {
			return fmt.Errorf("%s: Abschnitt %q trifft keine Überschrift (voller Heading-Text erwartet)", field, name)
		}
	}
	return nil
}

// diffViews bildet je Anforderung die beiden Mengendifferenzen
// (DC-FA-XREF-001.a Schritte 4–5): `exclude-req`-Kennungen fallen vor dem Diff
// aus beiden Sichten (Ableitungssprünge in Mittelschichten); `mode: superset`
// meldet nur B\F. Sortiert nach (Anforderung, Artefakt, Richtung).
func diffViews(fwd, bwd crossView, cc *model.TraceCrossConsistency) []CrossFinding {
	reqs := map[string]bool{}
	for r := range fwd {
		reqs[r] = true
	}
	for r := range bwd {
		reqs[r] = true
	}
	out := []CrossFinding{}
	for r := range reqs {
		if cc.ExcludeReq != nil && cc.ExcludeReq.MatchString(r) {
			continue
		}
		if cc.Mode != model.TraceCrossModeSuperset {
			out = append(out, crossMissing(r, fwd[r], bwd[r], CrossDirForwardOnly)...)
		}
		out = append(out, crossMissing(r, bwd[r], fwd[r], CrossDirBackwardOnly)...)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Requirement != b.Requirement {
			return a.Requirement < b.Requirement
		}
		if a.Artifact != b.Artifact {
			return a.Artifact < b.Artifact
		}
		return a.Direction < b.Direction
	})
	return out
}

// crossMissing liefert je Artefakt aus have, das other nicht kennt, einen Befund
// mit der Fundstelle aus have.
func crossMissing(req string, have, other map[string]crossEdge, dir string) []CrossFinding {
	var out []CrossFinding
	for artifact, e := range have {
		if _, known := other[artifact]; known {
			continue
		}
		out = append(out, CrossFinding{
			Requirement: req, Artifact: artifact, Direction: dir,
			File: e.file, Line: e.line,
		})
	}
	return out
}
