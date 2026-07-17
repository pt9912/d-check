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
// cc nil ⇒ kein Abgleich (RTM byte-identisch, DC-QA-02). reqPat ist das
// `trace.requirements.id-pattern` — es erkennt die Anforderungen der
// Vorwärts-ID-Spalte, während die Rück-Sicht ihr eigenes `req-pattern` für die
// Kanten-Zelle mitbringt (freier Prosa-Text statt kuratierter ID-Spalte).
//
// Die Phasen folgen der Fehlerpräzedenz des Vertrags und sind bewusst **über
// beide Sichten** gestaffelt, nicht je Sicht durchlaufen: Quellen lesen →
// Header-Bindung → Range-Expansion → Diff. Der erste Fehler beendet den Lauf.
// Reiner Lese-Pfad (DC-QA-03).
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
	fwdSrc, err := spanCrossSource(crossForwardField, cc.Forward.File, fwdContent, cc.Forward.Sections, cc.Forward.ExcludeSections)
	if err != nil {
		return nil, err
	}
	bwdSrc, err := spanCrossSource(crossBackwardField, cc.Backward.File, bwdContent, cc.Backward.Sections, nil)
	if err != nil {
		return nil, err
	}
	fwdTables, err := bindForwardTables(fwdSrc, cc.Forward)
	if err != nil {
		return nil, err
	}
	bwdTables, err := bindBackwardTables(bwdSrc, cc.Backward)
	if err != nil {
		return nil, err
	}
	fwd, err := forwardEdges(fwdTables, cc.Forward, reqPat)
	if err != nil {
		return nil, err
	}
	bwd, err := backwardEdges(bwdTables, cc.Backward, cc.Forward.DesignPattern)
	if err != nil {
		return nil, err
	}
	fwd, bwd = excludeReqs(fwd, cc.ExcludeReq), excludeReqs(bwd, cc.ExcludeReq)
	if err := crossVacuity(fwd, bwd, cc); err != nil {
		return nil, err
	}
	return diffViews(fwd, bwd, cc.Mode), nil
}

// crossSource ist eine gelesene Sicht-Quelle samt ihrer Abschnitts-Maske.
type crossSource struct {
	file    string
	content []byte
	mask    []bool
}

// boundTable ist eine **relevante** Tabelle mit aufgelösten Spaltenindizes:
// vorwärts (Anforderungs-Spalte, Design-Spalte), rückwärts (Artefakt-ID-Spalte,
// Kanten-Spalte).
type boundTable struct {
	rows      []mdTableRow
	primary   int
	secondary int
}

// readCrossFile liest eine Sicht-Quelle (Stufe „Quellen lesen"). Fail-closed:
// eine fehlende Datei ist Exit 2 — eine stillschweigend leere Sicht meldete sonst
// jede Kante der Gegenseite als Differenz.
func readCrossFile(fsys driven.Filesystem, field, file string) ([]byte, error) {
	if !pathExists(fsys, file) {
		return nil, fmt.Errorf("%s.file: Datei %q fehlt", field, file)
	}
	return fsys.ReadFile(file)
}

// spanCrossSource spannt eine gelesene Quelle auf ihre Abschnitte (eigene Stufe
// NACH dem Lesen beider Quellen — sonst verdeckte ein Vorwärts-Sektionsfehler die
// fehlende Rückwärts-Datei, die eine Stufe früher liegt).
func spanCrossSource(field, file string, content []byte, sections, exclude []string) (crossSource, error) {
	if err := checkCrossSections(field, content, sections, exclude); err != nil {
		return crossSource{}, err
	}
	return crossSource{file: file, content: content, mask: rules.SectionMask(content, sections, exclude)}, nil
}

// bindForwardTables bindet die Vorwärts-Rollen an ihre Spalten (Header-Phase).
func bindForwardTables(src crossSource, fc model.TraceCrossForward) ([]boundTable, error) {
	var out []boundTable
	for _, t := range markdownTables(src.content, src.mask) {
		idx, relevant, err := bindCrossColumns(t.header, fc.ReqColumn, fc.DesignColumn)
		if err != nil {
			return nil, fmt.Errorf("%s: Tabelle ab Zeile %d: %w", crossForwardField, t.line, err)
		}
		if !relevant {
			continue
		}
		if err := crossBadRow(crossForwardField, t); err != nil {
			return nil, err
		}
		out = append(out, boundTable{rows: t.rows, primary: idx[0], secondary: idx[1]})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s: keine Tabelle mit den konfigurierten Headern %q, %q in %q",
			crossForwardField, fc.ReqColumn, fc.DesignColumn, src.file)
	}
	return out, nil
}

// bindBackwardTables bindet die Rück-Rollen. Relevanz entsteht **allein** über die
// Kanten-Spalte (DC-FA-XREF-001.a Schritt 3: „zählt jede Tabelle mit einem
// edge-column-Header"); die Artefakt-ID-Spalte wird danach aufgelöst — fehlt ein
// konfigurierter ID-Header in einer relevanten Tabelle, ist das Exit 2 und kein
// stilles Überspringen, sonst verschwänden ihre Rück-Kanten lautlos.
func bindBackwardTables(src crossSource, bc model.TraceCrossBackward) ([]boundTable, error) {
	var out []boundTable
	for _, t := range markdownTables(src.content, src.mask) {
		idx, relevant, err := bindCrossColumns(t.header, bc.EdgeColumn)
		if err != nil {
			return nil, fmt.Errorf("%s: Tabelle ab Zeile %d: %w", crossBackwardField, t.line, err)
		}
		if !relevant {
			continue
		}
		idIdx, err := backwardIDColumn(t, bc)
		if err != nil {
			return nil, err
		}
		if err := crossBadRow(crossBackwardField, t); err != nil {
			return nil, err
		}
		out = append(out, boundTable{rows: t.rows, primary: idIdx, secondary: idx[0]})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s: keine Tabelle mit dem konfigurierten Header %q in %q",
			crossBackwardField, bc.EdgeColumn, src.file)
	}
	return out, nil
}

// backwardIDColumn löst die Artefakt-ID-Spalte einer bereits relevanten Tabelle
// auf: der Sentinel `first` nimmt die erste Spalte (heterogene ID-Header,
// ADR-0038); ein Header-Name muss genau einmal vorkommen (sonst Exit 2).
func backwardIDColumn(t mdTable, bc model.TraceCrossBackward) (int, error) {
	if bc.ArtifactIDColumn == model.TraceCrossArtifactFirst {
		return 0, nil
	}
	count, idx := 0, 0
	for i, cell := range t.header {
		if cell == bc.ArtifactIDColumn {
			if count == 0 {
				idx = i
			}
			count++
		}
	}
	if count != 1 {
		return 0, fmt.Errorf("%s: Tabelle ab Zeile %d: konfigurierter Header %q kommt %d-mal vor (genau einmal erwartet)",
			crossBackwardField, t.line, bc.ArtifactIDColumn, count)
	}
	return idx, nil
}

// forwardEdges extrahiert die Kanten der Vorwärts-Sicht (Range-Phase): je
// Anforderung der ID-Zelle (range-aware) jedes Artefakt der Design-Zelle.
func forwardEdges(tables []boundTable, fc model.TraceCrossForward, reqPat *regexp.Regexp) (crossView, error) {
	view := crossView{}
	for _, bt := range tables {
		for _, row := range bt.rows {
			reqs, err := rangeAwareIDs(crossForwardField, row.cells[bt.primary], reqPat, fc.Ranges)
			if err != nil {
				return nil, err
			}
			for _, req := range reqs {
				for _, artifact := range fc.DesignPattern.FindAllString(row.cells[bt.secondary], -1) {
					view.add(req, artifact, crossEdge{file: fc.File, line: row.line})
				}
			}
		}
	}
	return view, nil
}

// backwardEdges extrahiert die **invertierten** Kanten der Rück-Sicht: je
// Anforderung der Kanten-Zelle (range-aware) die Artefakt-ID der ID-Spalte. Die
// Artefakt-ID kommt bewusst über designPat — dasselbe Muster wie die
// Vorwärts-Sicht, denn nur ein gemeinsamer Namensraum macht den Diff bedeutsam.
func backwardEdges(tables []boundTable, bc model.TraceCrossBackward, designPat *regexp.Regexp) (crossView, error) {
	view := crossView{}
	for _, bt := range tables {
		for _, row := range bt.rows {
			artifact := designPat.FindString(row.cells[bt.primary])
			if artifact == "" {
				continue // Zeile ohne Artefakt-Kennung (Zwischenüberschrift o. Ä.)
			}
			reqs, err := rangeAwareIDs(crossBackwardField, row.cells[bt.secondary], bc.ReqPattern, bc.Ranges)
			if err != nil {
				return nil, err
			}
			for _, req := range reqs {
				view.add(req, artifact, crossEdge{file: bc.File, line: row.line})
			}
		}
	}
	return view, nil
}

// excludeReqs entfernt die auf exclude-req passenden Anforderungen aus einer Sicht
// (DC-FA-XREF-001.a Schritt 4: Ableitungssprünge in Mittelschichten). Eigene Stufe
// **vor** der Vakuitäts-Prüfung, denn maßgeblich ist, was am Ende tatsächlich
// verglichen wird: ein Ventil, das alle Anforderungen verschluckt, schaltet das
// Gate ebenso still ab wie ein fehlgreifendes Muster. nil ⇒ keine Ausnahme.
func excludeReqs(view crossView, exclude *regexp.Regexp) crossView {
	if exclude == nil {
		return view
	}
	out := crossView{}
	for req, artifacts := range view {
		if !exclude.MatchString(req) {
			out[req] = artifacts
		}
	}
	return out
}

// crossVacuity prüft, ob der Abgleich überhaupt etwas vergleicht
// (DC-FA-XREF-001.a Schritt 5). Ein Abgleich ohne eine einzige Kante ist **kein
// bestandener** Abgleich: die Muster greifen am Inhalt vorbei, und ein
// `0 Differenz(en)`/Exit 0 behauptete eine nie geprüfte Konsistenz.
//
// Vakuum ist genau zweierlei:
//   - **beide** Sichten kantenleer — typischer Anlass ist ein `design-pattern`,
//     das kompiliert, aber am Artefakt-Namensraum vorbeigreift; weil es zwischen
//     den Sichten geteilt ist (Schritt 3), räumt es beide zugleich leer;
//   - die **Rück**-Sicht kantenleer unter `mode: superset` — dort gatet allein
//     `B \ F`, also kann konstruktionsbedingt nie ein Befund entstehen.
//
// **Nicht** vakuum ist eine einseitig leere Sicht mit nicht-leerer Gegenseite:
// der Diff läuft über `keys(F) ∪ keys(B)` und ist dafür wohldefiniert. Eine noch
// unrestrukturierte Vorwärts-Sicht bei gepflegten Rück-Kanten ist der erwartete
// Bootstrap-Zustand (ADR-0038 Entscheidungen 3/7) — sie meldet `B \ F` laut. Ein
// symmetrisch je Sicht feuernder Guard würgte sie mit einer Config-Fehldiagnose ab.
func crossVacuity(fwd, bwd crossView, cc *model.TraceCrossConsistency) error {
	if len(fwd) == 0 && len(bwd) == 0 {
		return fmt.Errorf("trace.cross-consistency: beide Sichten ergaben 0 Kanten — der Abgleich verglich nichts; %s",
			vacuityHint(cc, crossForwardField+".design-pattern trifft den Artefakt-Namensraum beider Sichten"))
	}
	if len(bwd) == 0 && cc.Mode == model.TraceCrossModeSuperset {
		return fmt.Errorf("trace.cross-consistency: die Rück-Sicht ergab 0 Kanten und mode: superset gatet allein "+
			"B \\ F — der Abgleich könnte nie eine Differenz melden; %s",
			vacuityHint(cc, crossBackwardField+".edge-column/req-pattern treffen die Rück-Kanten"))
	}
	return nil
}

// vacuityHint benennt die möglichen Ursachen — und nennt das Ventil nur, wenn es
// überhaupt gesetzt ist. Eine Meldung, die pauschal in die Muster-Config schickt,
// wäre bei einem übergriffigen `exclude-req` eine Fehldiagnose (dieselbe Klasse,
// die den symmetrischen Guard aus dem Erst-Entwurf disqualifizierte).
func vacuityHint(cc *model.TraceCrossConsistency, patternHint string) string {
	if cc.ExcludeReq != nil {
		return "prüfe, ob " + patternHint + " — und ob exclude-req nicht alle Anforderungen verschluckt"
	}
	return "prüfe, ob " + patternHint
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
// (DC-FA-XREF-001.a Schritt 6) über `keys(F) ∪ keys(B)`; `mode: superset` meldet
// nur B\F. Die Ausschluss-Stufe ist bereits gelaufen (Schritt 4). Sortiert nach
// (Anforderung, Artefakt, Richtung).
func diffViews(fwd, bwd crossView, mode string) []CrossFinding {
	reqs := map[string]bool{}
	for r := range fwd {
		reqs[r] = true
	}
	for r := range bwd {
		reqs[r] = true
	}
	out := []CrossFinding{}
	for r := range reqs {
		if mode != model.TraceCrossModeSuperset {
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
