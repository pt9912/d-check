package rules

// Vierte planning-Fähigkeit (DC-FA-PLAN-001 §Register-Deckung): die eine
// maschinell entscheidbare Hälfte der Register-Paarung — eine zitierte
// Beobachtungs-Kennung hat eine Zeile im Register. Die Umkehrung ist
// ausgeschlossen: die meisten Zeilen stehen unter der Schwelle und sind
// nirgends zitiert.

import (
	"errors"
	"path"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckPlanningObservations meldet zitierte Kennungen ohne Registerzeile.
// Opt-in innerhalb des opt-in Moduls: ohne observations.register wird keine
// Datei geöffnet. Fail-closed bei unlesbarem Register oder Verzeichnis — ein
// Tippfehler im Pfad schaltete die Fähigkeit sonst still ab, und die leere
// Kennungs-Menge wäre widerspruchsfrei (dieselbe Disziplin wie waves.dir).
func CheckPlanningObservations(fsys driven.Filesystem, cfg model.PlanningConfig) []model.Finding {
	o := cfg.Observations
	if o.Register == "" || fsys == nil {
		return nil
	}
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
	dirs := o.Dirs
	if len(dirs) == 0 {
		dirs = []string{path.Dir(o.Register)}
	}
	var out []model.Finding
	for _, d := range dirs {
		files, err := markdownFilesUnder(fsys, d)
		if err != nil {
			out = append(out, observationFinding(o.Register, 1, d,
				"Zitier-Verzeichnis "+d+" fehlt oder ist unlesbar (fail-closed)")...)
			continue
		}
		for _, f := range files {
			out = append(out, citedWithoutRow(fsys, f, re, declared)...)
		}
	}
	return out
}

// declaredObservationIDs sammelt die Kennungen, die das Register als
// Tabellenzeile FÜHRT: erste Zelle einer Pipe-Zeile, die ganz aus der Kennung
// besteht. Eine Erwähnung im Fließtext einer anderen Zelle deklariert nichts —
// sonst deklarierte sich jede Quer-Referenz selbst.
func declaredObservationIDs(content string, re *regexp.Regexp) map[string]bool {
	out := map[string]bool{}
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
func citedWithoutRow(fsys driven.Filesystem, file string, re *regexp.Regexp, declared map[string]bool) []model.Finding {
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
			if declared[id] {
				continue
			}
			if inCodeOnly(codeByLine[pl.no], links, loc[0], loc[1]) {
				continue // Beispiel, keine Behauptung
			}
			out = append(out, observationFinding(file, pl.no, id,
				"Kennung "+id+" ist zitiert, hat aber keine Zeile im Register")...)
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
