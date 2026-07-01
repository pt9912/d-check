package rules

import (
	"path"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckPlanning ist das Regelmodul planning (DC-FA-PLAN-001): es prüft
// **hermetisch** (nur der Filesystem-Port, **kein** git, **kein** Netz) die
// Planning-Lifecycle-Invariante — der `planning.marker` steht im
// `planning.heading`-Block der Roadmap genau dann, wenn **kein**
// `planning.slice-glob`-Slice im Roadmap-Verzeichnis liegt (`hasActive ==
// hasSlices`); sonst `planning-drift`. Fail-closed bei fehlender kanonischer
// Überschrift bzw. fehlender/unlesbarer Roadmap-Datei. Opt-in (leere Roadmap ⇒
// inert). Diagnose-only: kein `--repair`-Hunk.
func CheckPlanning(fsys driven.Filesystem, cfg model.PlanningConfig) []model.Finding {
	if cfg.Roadmap == "" || fsys == nil {
		return nil // inert (Modul ohne Roadmap oder ohne Port)
	}
	sliceDir := path.Dir(cfg.Roadmap)
	content, err := fsys.ReadFile(cfg.Roadmap)
	if err != nil {
		return planningDrift(cfg.Roadmap, 1, sliceDir,
			"Roadmap-Datei "+cfg.Roadmap+" fehlt oder ist unlesbar (fail-closed)")
	}
	lines := splitLines(content)
	headingNo, headingCount := planningHeadingLine(lines, cfg.EffectiveHeading())
	if headingCount == 0 {
		return planningDrift(cfg.Roadmap, 1, sliceDir,
			"kanonische Überschrift „"+cfg.EffectiveHeading()+"“ fehlt in "+cfg.Roadmap+
				" — Aktiv-Status nicht bestimmbar (fail-closed)")
	}
	if headingCount > 1 {
		return planningDrift(cfg.Roadmap, headingNo, sliceDir,
			"kanonische Überschrift „"+cfg.EffectiveHeading()+"“ kommt mehrfach in "+cfg.Roadmap+
				" vor — Aktiv-Status mehrdeutig (fail-closed)")
	}
	hasActive := !planningBlockHasMarker(lines, headingNo, cfg.EffectiveMarker())
	hasSlices := planningHasSlices(fsys, sliceDir, cfg.EffectiveSliceGlob())
	if hasActive == hasSlices {
		return nil
	}
	var msg string
	if hasSlices {
		msg = "Slice(s) in " + sliceDir + ", aber die Roadmap §Aktuelle Welle trägt den Ruhe-Marker „" +
			cfg.EffectiveMarker() + "“ — die Roadmap muss die aktive Welle benennen"
	} else {
		msg = "kein Slice in " + sliceDir + ", aber die Roadmap §Aktuelle Welle benennt eine aktive Welle" +
			" — die Roadmap muss den Ruhe-Marker „" + cfg.EffectiveMarker() + "“ tragen"
	}
	return planningDrift(cfg.Roadmap, headingNo, sliceDir, msg)
}

// planningDrift baut den (einzelnen) planning-drift-Befund.
func planningDrift(roadmap string, line int, sliceDir, msg string) []model.Finding {
	return []model.Finding{{
		File: roadmap, Line: line, Rule: "planning", Target: sliceDir,
		Reason: model.ReasonPlanningDrift, Message: msg,
	}}
}

// planningHeadingLine liefert die 1-basierte Zeilennummer der **ersten** kanonischen
// Überschrift (getrimmt exakt gleich; der Skript-Guard `^## Aktuelle Welle[[:space:]]*$`:
// keine führenden Zeichen, nur nachlaufender Whitespace) **und** die Anzahl der
// Vorkommen. count == 0 ⇒ fehlt (fail-closed); count > 1 ⇒ mehrdeutig (fail-closed —
// bei zwei Blöcken wäre der geprüfte Ausschnitt sonst uneindeutig, ein Silent-Grün-Risiko).
func planningHeadingLine(lines []string, heading string) (first, count int) {
	for i, raw := range lines {
		if strings.TrimRight(raw, " \t\r") == heading {
			count++
			if first == 0 {
				first = i + 1
			}
		}
	}
	return first, count
}

// planningBlockHasMarker sucht den Ruhe-Marker (literaler Teilstring) NUR im
// Aktiv-Status-Block: ab der Zeile **nach** der Überschrift (headingNo) bis zur
// nächsten `## `-H2 (exklusive). So verfälscht ein erklärender Verweis anderswo
// den Status nicht.
func planningBlockHasMarker(lines []string, headingNo int, marker string) bool {
	for i := headingNo; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			return false // nächste H2 erreicht
		}
		if strings.Contains(lines[i], marker) {
			return true
		}
	}
	return false
}

// planningHasSlices meldet, ob mindestens ein Verzeichnis-Eintrag von dir seinen
// Basisnamen gegen slice-glob (per path.Match) matcht.
func planningHasSlices(fsys driven.Filesystem, dir, glob string) bool {
	entries, err := fsys.List(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if ok, _ := path.Match(glob, e.Name); ok {
			return true
		}
	}
	return false
}
