package rules

import (
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckPlanningWaves ist die dritte planning-Fähigkeit (DC-FA-PLAN-001
// §Wellen-Invariante, Spez-Schritte W1–W5): dieselbe Lifecycle-Invariante eine
// Ebene höher — die Wellen-Abschnitte der Roadmap gegen die Wellen-Dateien.
// Opt-in innerhalb des opt-in Moduls: ohne waves.dir wird kein Wellendokument
// geöffnet. Diagnose-only.
func CheckPlanningWaves(fsys driven.Filesystem, cfg model.PlanningConfig) []model.Finding {
	w := cfg.Waves
	if w.Dir == "" || cfg.Roadmap == "" || fsys == nil {
		return nil // W1: inert
	}
	content, err := fsys.ReadFile(cfg.Roadmap)
	if err != nil {
		return waveFinding(cfg.Roadmap, 1, w.Dir, model.ReasonWaveDrift,
			"Roadmap-Datei "+cfg.Roadmap+" fehlt oder ist unlesbar (fail-closed)")
	}
	sliceDir := path.Dir(cfg.Roadmap)
	headingNo, hasActive, fail := planningActiveStatus(content, cfg, sliceDir)
	if fail != nil {
		// Der Aktiv-Status trägt beide Fähigkeiten; ist er nicht bestimmbar,
		// meldet ihn die Slice-Invariante bereits — hier kein Doppelbefund.
		return nil
	}
	flach, ruhe := waveSets(fsys, w)
	out := waveDrift(cfg, w, headingNo, hasActive, flach)
	return append(out, waveRegisters(cfg, w, content, flach, ruhe)...)
}

// waveSets liefert die beiden Wellen-Mengen als Kennungs-Mengen (W2): die
// **flachen** Plan-Dokumente aus waves.dir (glob abzüglich results-glob) und
// die **Ergebnisnotizen** aus dem Ruheort. Ein unlesbares Verzeichnis liefert
// eine leere Menge — der daraus folgende Befund ist fail-closed, nicht still.
func waveSets(fsys driven.Filesystem, w model.WavesConfig) (flach, ruhe map[string]bool) {
	prefix := wavePrefix(w.EffectiveGlob())
	flach = map[string]bool{}
	for _, name := range dirNames(fsys, w.Dir) {
		if !matchBase(w.EffectiveGlob(), name) || matchBase(w.EffectiveResultsGlob(), name) {
			continue
		}
		if id := waveID(prefix, name); id != "" {
			flach[id] = true
		}
	}
	ruhe = map[string]bool{}
	for _, name := range dirNames(fsys, w.EffectiveDoneDir()) {
		if !matchBase(w.EffectiveResultsGlob(), name) {
			continue
		}
		if id := waveID(prefix, name); id != "" {
			ruhe[id] = true
		}
	}
	return flach, ruhe
}

// waveDrift ist W3: der Aktiv-Status gegen die Zahl der flachen Wellendokumente.
func waveDrift(
	cfg model.PlanningConfig, w model.WavesConfig, headingNo int, hasActive bool, flach map[string]bool,
) []model.Finding {
	if hasActive == (len(flach) == 1) {
		return nil
	}
	var msg string
	switch {
	case hasActive && len(flach) == 0:
		msg = "die Roadmap benennt eine aktive Welle, aber in " + w.Dir +
			" liegt kein flaches Wellendokument"
	case hasActive:
		msg = "die Roadmap benennt eine aktive Welle, aber in " + w.Dir +
			" liegen " + strconv.Itoa(len(flach)) + " flache Wellendokumente — genau eines ist erwartet"
	default:
		msg = "die Roadmap trägt den Ruhe-Marker „" + cfg.EffectiveMarker() + "“, aber in " +
			w.Dir + " liegt " + strconv.Itoa(len(flach)) + " flaches Wellendokument"
	}
	return waveFinding(cfg.Roadmap, headingNo, w.Dir, model.ReasonWaveDrift, msg)
}

// waveRegisters ist W4/W5: die drei Register-Aussagen über Vorschau und
// Abschluss-Register.
func waveRegisters(
	cfg model.PlanningConfig, w model.WavesConfig, content []byte, flach, ruhe map[string]bool,
) []model.Finding {
	lines := splitLines(content)
	prose := proseLineSet(content)
	prefix := wavePrefix(w.EffectiveGlob())
	var out []model.Finding

	for _, r := range waveTableRows(lines, prose, w.EffectiveNextHeading(), prefix) {
		if flach[r.id] || ruhe[r.id] {
			out = append(out, waveFinding(cfg.Roadmap, r.line, r.id, model.ReasonWavePreviewExists,
				"Vorschau-Zeile nennt "+r.id+", für die bereits eine Datei existiert — "+
					"eine geplante Welle steht in der Vorschau und nirgends sonst")...)
		}
	}

	closedNo, rowsClosed := waveClosedRows(lines, prose, w.EffectiveClosedHeading(), prefix)
	registriert := map[string]bool{}
	for _, r := range rowsClosed {
		registriert[r.id] = true
		if !ruhe[r.id] {
			out = append(out, waveFinding(cfg.Roadmap, r.line, r.id, model.ReasonWaveResultsMissing,
				"Abschluss-Register nennt "+r.id+", aber in "+w.EffectiveDoneDir()+
					" liegt keine Ergebnisnotiz nach "+strconv.Quote(w.EffectiveResultsGlob()))...)
		}
	}
	if closedNo != 0 {
		var fehlend []string
		for id := range ruhe {
			if !registriert[id] {
				fehlend = append(fehlend, id)
			}
		}
		sort.Strings(fehlend) // DC-QA-02
		for _, id := range fehlend {
			out = append(out, waveFinding(cfg.Roadmap, closedNo, id, model.ReasonWaveUnregistered,
				"Ergebnisnotiz zu "+id+" liegt in "+w.EffectiveDoneDir()+
					", das Abschluss-Register nennt sie nicht")...)
		}
	}
	return out
}

// waveRow ist eine Registerzeile mit ihrer Zeilennummer und der Kennung aus der
// ERSTEN Spalte.
type waveRow struct {
	line int
	id   string
}

// waveTableRows liefert die Registerzeilen unter der gegebenen Überschrift:
// Tabellenzeilen außerhalb von Fenced-Code (dieselbe Lexik wie targets), gelesen
// wird die erste Spalte. Zeilen ohne Kennung werden übersprungen — eine geplante
// Welle trägt einen Namen, und die Trigger-Spalte darf andere Wellen nennen.
func waveTableRows(lines []string, prose map[int]bool, heading, prefix string) []waveRow {
	no, _ := planningHeadingLine(lines, prose, heading)
	if no == 0 {
		return nil
	}
	return waveRowsFrom(lines, prose, no, prefix)
}

// waveClosedRows liefert zusätzlich die Zeile der Überschrift — sie trägt die
// Befunde der Gegenrichtung (Ergebnisnotiz ohne Registerzeile).
func waveClosedRows(lines []string, prose map[int]bool, heading, prefix string) (int, []waveRow) {
	no, _ := planningHeadingLine(lines, prose, heading)
	if no == 0 {
		return 0, nil
	}
	return no, waveRowsFrom(lines, prose, no, prefix)
}

// waveRowsFrom sammelt die Registerzeilen ab der Überschrift bis zur nächsten
// Überschrift gleicher oder höherer Ordnung (geteilte Abschnitts-Mechanik).
func waveRowsFrom(lines []string, prose map[int]bool, headingNo int, prefix string) []waveRow {
	level, _, ok := parseATXHeading(lines[headingNo-1])
	if !ok {
		level = 2
	}
	end := SectionEnd(lines, headingNo, level)
	if end == 0 || end > len(lines)+1 {
		end = len(lines) + 1
	}
	idRE := regexp.MustCompile(regexp.QuoteMeta(prefix) + `\d+`)
	var out []waveRow
	for i := headingNo; i < end-1; i++ {
		if !prose[i+1] || !strings.HasPrefix(lines[i], "|") {
			continue
		}
		cells := strings.Split(strings.Trim(lines[i], "|"), "|")
		if len(cells) == 0 {
			continue
		}
		if id := idRE.FindString(cells[0]); id != "" {
			out = append(out, waveRow{line: i + 1, id: id})
		}
	}
	return out
}

// wavePrefix liefert das Kennungs-Präfix aus dem Plan-Glob (der Teil vor dem
// ersten Platzhalter) — daran hängen Datei-Zuordnung und Zeilen-Erkennung.
func wavePrefix(glob string) string {
	if i := strings.IndexAny(glob, "*?["); i > 0 {
		return glob[:i]
	}
	return "welle-"
}

// waveID liefert die Kennung eines Basisnamens: Präfix plus die unmittelbar
// folgende Ziffernfolge (`welle-74-geteilte-lexik.md` ⇒ `welle-74`).
func waveID(prefix, name string) string {
	if !strings.HasPrefix(name, prefix) {
		return ""
	}
	rest := name[len(prefix):]
	n := 0
	for n < len(rest) && rest[n] >= '0' && rest[n] <= '9' {
		n++
	}
	if n == 0 {
		return ""
	}
	return prefix + rest[:n]
}

// dirNames liefert die Basisnamen eines Verzeichnisses; ein unlesbares
// Verzeichnis liefert nichts (der Folgebefund ist fail-closed).
func dirNames(fsys driven.Filesystem, dir string) []string {
	entries, err := fsys.List(dir)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name)
	}
	sort.Strings(out) // DC-QA-02
	return out
}

// matchBase prüft einen Basisnamen gegen einen Glob (path.Match, wie überall).
func matchBase(glob, name string) bool {
	ok, err := path.Match(glob, name)
	return err == nil && ok
}

func waveFinding(file string, line int, target, reason, msg string) []model.Finding {
	return []model.Finding{{
		File: file, Line: line, Rule: "planning", Target: target,
		Reason: reason, Message: msg,
	}}
}
