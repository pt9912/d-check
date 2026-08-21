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
		// Die unlesbare Roadmap meldet die Slice-Invariante bereits fail-closed;
		// hier kein Doppelbefund — dieselbe Disziplin wie beim unbestimmbaren
		// Aktiv-Status.
		return nil
	}
	sliceDir := path.Dir(cfg.Roadmap)
	headingNo, hasActive, fail := planningActiveStatus(content, cfg, sliceDir)
	if fail != nil {
		// Der Aktiv-Status trägt beide Fähigkeiten; ist er nicht bestimmbar,
		// meldet ihn die Slice-Invariante bereits — hier kein Doppelbefund.
		return nil
	}
	flach, ruhe, dirErrs := waveSets(fsys, w)
	if len(dirErrs) > 0 {
		// Ein Tippfehler im Pfad schaltete die Faehigkeit sonst still ab —
		// im Ruhe-Zustand waere die leere Menge konsistent und kein Befund
		// entstuende (dieselbe Disziplin wie closure.dir, C2).
		var out []model.Finding
		for _, d := range dirErrs {
			out = append(out, waveFinding(cfg.Roadmap, headingNo, d, model.ReasonWaveDrift,
				"Wellen-Verzeichnis "+d+" fehlt oder ist unlesbar (fail-closed)")...)
		}
		return out
	}
	out := waveDrift(cfg, w, headingNo, hasActive, flach)
	return append(out, waveRegisters(cfg, w, content, flach, ruhe)...)
}

// waveSets liefert die beiden Wellen-Mengen (W2): die **flachen**
// Plan-Dokumente aus waves.dir (glob abzüglich results-glob) als
// Kennungs-Menge und die **Ergebnisnotizen** aus dem Ruheort als
// Kennung→Pfad (der Pfad wird das Befund-`target` der Gegenrichtung).
// Unlesbare Verzeichnisse kommen als dirErrs zurück — der Aufrufer meldet sie
// fail-closed.
func waveSets(fsys driven.Filesystem, w model.WavesConfig) (flach map[string]bool, ruhe map[string]string, dirErrs []string) {
	prefix := wavePrefix(w.EffectiveGlob())
	flach = map[string]bool{}
	names, err := dirNames(fsys, w.Dir)
	if err != nil {
		dirErrs = append(dirErrs, w.Dir)
	}
	for _, name := range names {
		if !matchBase(w.EffectiveGlob(), name) || matchBase(w.EffectiveResultsGlob(), name) {
			continue
		}
		if id := waveID(prefix, name); id != "" {
			flach[id] = true
		}
	}
	ruhe = map[string]string{}
	names, err = dirNames(fsys, w.EffectiveDoneDir())
	if err != nil && w.EffectiveDoneDir() != w.Dir {
		dirErrs = append(dirErrs, w.EffectiveDoneDir())
	}
	for _, name := range names {
		if !matchBase(w.EffectiveResultsGlob(), name) {
			continue
		}
		if id := waveID(prefix, name); id != "" {
			ruhe[id] = path.Join(w.EffectiveDoneDir(), name)
		}
	}
	return flach, ruhe, dirErrs
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
// Abschluss-Register. Eine fehlende Register-Überschrift ist fail-closed —
// die Fähigkeit zu aktivieren IST die Behauptung, dass die Roadmap beide
// Register führt; sonst schaltete ein Tippfehler in der Überschrift alle drei
// Aussagen wortlos ab (dieselbe Disziplin wie der Heading-Guard der
// Slice-Invariante).
func waveRegisters(
	cfg model.PlanningConfig, w model.WavesConfig, content []byte, flach map[string]bool, ruhe map[string]string,
) []model.Finding {
	lines := splitLines(content)
	prose := proseLineSet(content)
	prefix := wavePrefix(w.EffectiveGlob())
	var out []model.Finding

	// Das Befund-Ziel ist je Bedeutung eindeutig (die fehlende Überschrift bzw.
	// das Verzeichnis) — die Deduplikation über (Datei, Zeile, Regel, Ziel,
	// Grund) ließe zwei wave-drift-Befunde desselben Tupels sonst
	// zusammenfallen, und eine der beiden Reparaturen verschwände.
	nextNo, rowsNext := waveRegisterRows(lines, prose, w.EffectiveNextHeading(), prefix)
	if nextNo == 0 {
		out = append(out, waveFinding(cfg.Roadmap, 1, w.EffectiveNextHeading(), model.ReasonWaveDrift,
			"Register-Überschrift „"+w.EffectiveNextHeading()+"“ fehlt in "+cfg.Roadmap+
				" — die Vorschau-Aussage ist nicht prüfbar (fail-closed)")...)
	}
	for _, r := range rowsNext {
		if flach[r.id] || ruhe[r.id] != "" {
			out = append(out, waveFinding(cfg.Roadmap, r.line, r.id, model.ReasonWavePreviewExists,
				"Vorschau-Zeile nennt "+r.id+", für die bereits eine Datei existiert — "+
					"eine geplante Welle steht in der Vorschau und nirgends sonst")...)
		}
	}

	closedNo, rowsClosed := waveRegisterRows(lines, prose, w.EffectiveClosedHeading(), prefix)
	if closedNo == 0 {
		return append(out, waveFinding(cfg.Roadmap, 1, w.EffectiveClosedHeading(), model.ReasonWaveDrift,
			"Register-Überschrift „"+w.EffectiveClosedHeading()+"“ fehlt in "+cfg.Roadmap+
				" — die Abschluss-Aussagen sind nicht prüfbar (fail-closed)")...)
	}
	registriert := map[string]bool{}
	for _, r := range rowsClosed {
		registriert[r.id] = true
		if ruhe[r.id] == "" {
			out = append(out, waveFinding(cfg.Roadmap, r.line, r.id, model.ReasonWaveResultsMissing,
				"Abschluss-Register nennt "+r.id+", aber in "+w.EffectiveDoneDir()+
					" liegt keine Ergebnisnotiz nach "+strconv.Quote(w.EffectiveResultsGlob()))...)
		}
	}
	var fehlend []string
	for id := range ruhe {
		if !registriert[id] {
			fehlend = append(fehlend, id)
		}
	}
	sort.Strings(fehlend) // DC-QA-02
	for _, id := range fehlend {
		out = append(out, waveFinding(cfg.Roadmap, closedNo, ruhe[id], model.ReasonWaveUnregistered,
			"Ergebnisnotiz "+ruhe[id]+" liegt im Ruheort, das Abschluss-Register nennt sie nicht")...)
	}
	return out
}

// waveRow ist eine Registerzeile mit ihrer Zeilennummer und der Kennung aus der
// ERSTEN Spalte.
type waveRow struct {
	line int
	id   string
}

// waveRegisterRows liefert die Zeile der Register-Überschrift (0 ⇒ fehlt,
// fail-closed beim Aufrufer) und die Registerzeilen darunter: Tabellenzeilen
// nach der geteilten Lexik, gelesen wird die erste Spalte. Zeilen ohne Kennung
// werden übersprungen — eine geplante Welle trägt einen Namen, und die
// Trigger-Spalte darf andere Wellen nennen.
func waveRegisterRows(lines []string, prose map[int]bool, heading, prefix string) (int, []waveRow) {
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
		if !tableRowLine(lines, prose, i) || tableHeaderOrSeparator(lines, prose, i) {
			continue
		}
		cells := tableCells(lines[i])
		if len(cells) == 0 {
			continue
		}
		if id := idRE.FindString(cells[0]); id != "" {
			out = append(out, waveRow{line: i + 1, id: id})
		}
	}
	return out
}

// wavePrefix liefert das Kennungs-Präfix eines Globs: der literale Teil vor
// dem ersten Platzhalter. Daran hängen Datei-Zuordnung und Zeilen-Erkennung;
// dass er nicht leer ist und results-glob dasselbe Präfix trägt, stellt der
// Config-Rand sicher (Exit 2, Spez-Schritt W1).
func wavePrefix(glob string) string {
	if i := strings.IndexAny(glob, "*?["); i >= 0 {
		return glob[:i]
	}
	return glob
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

// dirNames liefert die Basisnamen eines Verzeichnisses, stabil sortiert
// (DC-QA-02); der Fehler eines unlesbaren Verzeichnisses geht an den Aufrufer.
func dirNames(fsys driven.Filesystem, dir string) ([]string, error) {
	entries, err := fsys.List(dir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name)
	}
	sort.Strings(out)
	return out, nil
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
