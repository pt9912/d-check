package rules

// Vierte planning-Fähigkeit (DC-FA-PLAN-001 §Register-Deckung, Tabellen-Modus)
// und ihre additive fünfte (Verzeichnis-Modus, ADR-0083): die eine maschinell
// entscheidbare Hälfte der Register-Paarung — eine zitierte
// Beobachtungs-Kennung hat einen Nachweis (Tabellenzeile oder Verzeichnis).
// Die Umkehrung ist ausgeschlossen: die meisten Kennungen stehen unter der
// Schwelle und sind nirgends zitiert.

import (
	"errors"
	"path"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckPlanningObservations meldet zitierte Kennungen ohne Nachweis — als
// Registerzeile (Tabellen-Modus, `Register`) oder als Verzeichnis
// (Verzeichnis-Modus, `Dir`, ADR-0083). Opt-in innerhalb des opt-in Moduls:
// ohne beides wird keine Datei geöffnet. `Register` und `Dir` sind mutually
// exclusive — durchgesetzt am Systemrand beim Config-Laden
// (configyaml.applyObservations, Exit 2), hier als Invariante vertraut, nicht
// erneut geprüft. Fail-closed bei unlesbarem Register/Verzeichnis — ein
// Tippfehler im Pfad schaltete die Fähigkeit sonst still ab, und die leere
// Kennungs-Menge wäre widerspruchsfrei (dieselbe Disziplin wie waves.dir).
func CheckPlanningObservations(fsys driven.Filesystem, cfg model.PlanningConfig) []model.Finding {
	o := cfg.Observations
	if fsys == nil {
		return nil
	}
	switch {
	case o.Register != "":
		return checkObservationsTable(fsys, o)
	case o.Dir != "":
		return checkObservationsDir(fsys, o)
	default:
		return nil
	}
}

// checkObservationsTable ist die vierte planning-Fähigkeit: Tabellen-Modus.
func checkObservationsTable(fsys driven.Filesystem, o model.ObservationsConfig) []model.Finding {
	re, err := regexp.Compile(o.EffectivePattern())
	if err != nil {
		return observationFinding(o.Register, 1, o.Register,
			"Kennungs-Muster "+o.EffectivePattern()+" ist kein gültiger Ausdruck (fail-closed)")
	}
	content, err := fsys.ReadFile(o.Register)
	if err != nil {
		return observationFinding(o.Register, 1, o.Register,
			"Register "+o.Register+" fehlt oder ist unlesbar (fail-closed)")
	}
	declared := declaredObservationIDs(string(content), re)
	return scanCitedWithoutDeclaration(fsys, o, o.Register, declared, re)
}

// checkObservationsDir ist die fünfte planning-Fähigkeit (ADR-0083):
// Verzeichnis-Modus. Eine zitierte Kennung <pfad> gilt als nachgewiesen, wenn
// <Dir>/<pfad>/observation.md existiert — kein Tabellenzeilen-Parsing.
func checkObservationsDir(fsys driven.Filesystem, o model.ObservationsConfig) []model.Finding {
	re, err := regexp.Compile(o.EffectivePattern())
	if err != nil {
		return observationFinding(o.Dir, 1, o.Dir,
			"Kennungs-Muster "+o.EffectivePattern()+" ist kein gültiger Ausdruck (fail-closed)")
	}
	if k, err := fsys.Kind(o.Dir); err != nil || k != driven.KindDir {
		return observationFinding(o.Dir, 1, o.Dir,
			"Verzeichnis "+o.Dir+" fehlt oder ist unlesbar (fail-closed)")
	}
	declared := declaredObservationIDsFromDir{fsys: fsys, root: o.Dir}
	return scanCitedWithoutDeclaration(fsys, o, o.Dir, declared, re)
}

// scanCitedWithoutDeclaration ist der modusunabhängige Rest: Zitier-Dateien
// sammeln, jede gegen die deklarierte Menge halten. declaredSet abstrahiert
// "hat diese Kennung einen Nachweis" — Tabellenzeile oder Verzeichnis.
func scanCitedWithoutDeclaration(fsys driven.Filesystem, o model.ObservationsConfig, self string, declared declaredSet, re *regexp.Regexp) []model.Finding {
	dirs := o.Dirs
	if len(dirs) == 0 {
		dirs = []string{path.Dir(self)}
	}
	var out []model.Finding
	for _, d := range dirs {
		files, err := markdownFilesUnder(fsys, d)
		if err != nil {
			out = append(out, observationFinding(self, 1, d,
				"Zitier-Verzeichnis "+d+" fehlt oder ist unlesbar (fail-closed)")...)
			continue
		}
		for _, f := range files {
			out = append(out, citedWithoutRow(fsys, f, re, declared)...)
		}
	}
	return out
}

// declaredSet beantwortet "hat diese Kennung einen Nachweis" — modusunabhängig.
type declaredSet interface {
	Has(id string) bool
}

// mapDeclared adaptiert die Tabellen-Modus-Menge (map) auf declaredSet.
type mapDeclared map[string]bool

func (m mapDeclared) Has(id string) bool { return m[id] }

// declaredObservationIDsFromDir prüft je Kennung die Existenz von
// <root>/<id>/observation.md — ohne Vorab-Scan, jede Prüfung ist ein
// einzelner Filesystem-Zugriff (dieselbe Kosten-Form wie ein Kartei-Lookup).
type declaredObservationIDsFromDir struct {
	fsys driven.Filesystem
	root string
}

func (d declaredObservationIDsFromDir) Has(id string) bool {
	k, err := d.fsys.Kind(path.Join(d.root, id, "observation.md"))
	return err == nil && k == driven.KindFile
}

// declaredObservationIDs sammelt die Kennungen, die das Register als
// Tabellenzeile FÜHRT: erste Zelle einer Pipe-Zeile, die ganz aus der Kennung
// besteht. Eine Erwähnung im Fließtext einer anderen Zelle deklariert nichts —
// sonst deklarierte sich jede Quer-Referenz selbst.
func declaredObservationIDs(content string, re *regexp.Regexp) mapDeclared {
	out := mapDeclared{}
	for _, ln := range strings.Split(content, "\n") {
		t := strings.TrimSpace(ln)
		if !strings.HasPrefix(t, "|") {
			continue
		}
		cell := strings.TrimSpace(strings.SplitN(strings.TrimPrefix(t, "|"), "|", 2)[0])
		cell = strings.Trim(cell, "`")
		if m := re.FindString(cell); m != "" && m == cell {
			out[cell] = true
		}
	}
	return out
}

// citedWithoutRow meldet je Vorkommen, das eine Registerzeile beansprucht.
// Gezählt werden Prosa UND Linktext; ein reines Inline-Code-Span zählt nicht —
// gemessen ist das die Trennlinie zwischen Zitat und Beispiel, und sie ist
// zwingend, weil die verbreitete Zitier-Form ([`BEO-NNN`](…)) ihre Kennung
// gerade IN Backticks führt. Fenced Code bleibt außen (Vorverarbeitung).
func citedWithoutRow(fsys driven.Filesystem, file string, re *regexp.Regexp, declared declaredSet) []model.Finding {
	content, err := fsys.ReadFile(file)
	if err != nil {
		return nil // unlesbare Einzeldatei: der Scan meldet sie an anderer Stelle
	}
	prose := proseLines(content)
	codeByLine := inlineSpansByLine(prose)
	var out []model.Finding
	for _, pl := range prose {
		links := ExtractLinkSpans(pl.raw)
		for _, loc := range re.FindAllStringIndex(pl.raw, -1) {
			id := pl.raw[loc[0]:loc[1]]
			if declared.Has(id) {
				continue
			}
			if inCodeOnly(codeByLine[pl.no], links, loc[0], loc[1]) {
				continue // Beispiel, keine Behauptung
			}
			out = append(out, observationFinding(file, pl.no, id,
				"Kennung "+id+" ist zitiert, hat aber keinen Nachweis im Register")...)
		}
	}
	return out
}

// inCodeOnly meldet, ob [start,end) in einem Inline-Code-Span liegt, OHNE
// zugleich verlinkt zu sein. Die Verlinkungs-Frage beantwortet
// IDOccurrenceExempt — dieselbe Stelle, die ids fuer seine Linkpflicht ruft;
// eine zweite Antwort auf dieselbe Frage waere ein Defekt.
func inCodeOnly(code []inlineSpan, links []LinkSpan, start, end int) bool {
	for _, sp := range code {
		if start >= sp.start && end <= sp.end {
			return !IDOccurrenceExempt(links, start, end)
		}
	}
	return false
}

// markdownFilesUnder sammelt die Markdown-Dateien eines Verzeichnisbaums.
// Die Art des Wurzelpfads wird ausdruecklich geprueft: ein vertipptes
// Verzeichnis liefert bei manchen Adaptern eine LEERE Liste statt eines
// Fehlers, und die Faehigkeit ginge damit still inert — gemessen am
// In-Memory-Adapter der Kern-Tests.
func markdownFilesUnder(fsys driven.Filesystem, dir string) ([]string, error) {
	if k, err := fsys.Kind(dir); err != nil || k != driven.KindDir {
		return nil, errNotADir
	}
	entries, err := fsys.List(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		p := path.Join(dir, e.Name)
		if e.Kind == driven.KindDir {
			sub, err := markdownFilesUnder(fsys, p)
			if err != nil {
				return nil, err
			}
			out = append(out, sub...)
			continue
		}
		if strings.HasSuffix(e.Name, ".md") {
			out = append(out, p)
		}
	}
	return out, nil
}

func observationFinding(file string, line int, target, msg string) []model.Finding {
	return []model.Finding{{
		File:   file,
		Line:   line,
		Rule:   "planning",
		Target: target,
		Reason: model.ReasonObservationUnregistered,
		Message: msg,
	}}
}

// errNotADir traegt den fail-closed-Fall eines Zitier-Verzeichnisses, das
// keines ist (fehlend, Datei oder unlesbar).
var errNotADir = errors.New("kein lesbares Verzeichnis")
