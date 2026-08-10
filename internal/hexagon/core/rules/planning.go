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

// CheckPlanningClosure ist die zweite planning-Fähigkeit (DC-FA-PLAN-001
// §Closure-Note-Struktur, Spez-Schritte C1–C5): sie prüft die **Struktur** der
// Closure-Notizen abgeschlossener Slices. Opt-in **innerhalb** des opt-in
// Moduls — ohne cfg.Closure.Dir wird keine Slice-Datei geöffnet und der
// Befundsatz ist byte-identisch. Läuft unabhängig von der Aktiv-Status-Prüfung:
// beide können nebeneinander Befunde liefern.
//
// Geprüft wird Struktur, nicht Bedeutung — ob eine formal ausreichende Notiz
// inhaltlich trägt, ist semantisch und ausdrücklich nicht zugesagt.
func CheckPlanningClosure(fsys driven.Filesystem, cfg model.PlanningConfig) []model.Finding {
	dir := cfg.Closure.Dir
	if dir == "" || fsys == nil {
		return nil // C1: inert — keine Slice-Datei wird geöffnet
	}
	// Das Muster ist am Config-Rand bereits auf Kompilierbarkeit geprüft
	// (Exit 2); hier ist ein Fehler nicht mehr erreichbar, wird aber
	// fail-closed behandelt statt ignoriert.
	re, err := regexp.Compile(cfg.Closure.EffectiveHeadingPattern())
	if err != nil {
		return closureFinding(dir, 1, dir, model.ReasonClosureNoteMissing,
			"heading-pattern "+cfg.Closure.EffectiveHeadingPattern()+" ist kein gültiges Regex (fail-closed)")
	}
	entries, err := fsys.List(dir)
	if err != nil {
		// C2 fail-closed: gesetztes, aber fehlendes/unlesbares Verzeichnis ist
		// kein stilles Grün — sonst schaltete ein Tippfehler im Pfad die ganze
		// Prüfung ab.
		return closureFinding(dir, 1, dir, model.ReasonClosureNoteMissing,
			"Closure-Verzeichnis "+dir+" fehlt oder ist unlesbar (fail-closed)")
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if ok, _ := path.Match(cfg.EffectiveClosureGlob(), e.Name); ok {
			names = append(names, e.Name)
		}
	}
	// C2 fail-closed, zweite Hälfte: `closure.dir` ist der Aktivierungs-Schalter —
	// ihn zu setzen IST die Behauptung, dass dort Closure-Notizen liegen. Null
	// Kandidaten heißt also nicht „nichts zu tun", sondern „die Behauptung stimmt
	// nicht mehr" (typischer Auslöser: der Bestand wandert in Unterordner, das
	// Gate liefe fortan leer und grün). Dieselbe Logik wie der Nullmengen-Guard
	// der RTM-Anforderungsquellen — R1-F-5.
	if len(names) == 0 {
		return closureFinding(dir, 1, dir, model.ReasonClosureNoteMissing,
			"Closure-Verzeichnis "+dir+" enthält keine Datei nach "+strconv.Quote(cfg.EffectiveClosureGlob())+
				" — das Gate liefe leer (fail-closed; ist der Bestand umgezogen?)")
	}
	sort.Strings(names) // stabile Reihenfolge (DC-QA-02)
	var out []model.Finding
	for _, name := range names {
		out = append(out, checkClosureNote(fsys, cfg, dir, name, re)...)
	}
	return out
}

// checkClosureNote prüft eine einzelne Slice-Datei (Spez-Schritte C3–C5).
func checkClosureNote(
	fsys driven.Filesystem, cfg model.PlanningConfig, dir, name string, re *regexp.Regexp,
) []model.Finding {
	file := path.Join(dir, name)
	content, err := fsys.ReadFile(file)
	if err != nil {
		return closureFinding(file, 1, dir, model.ReasonClosureNoteMissing,
			"Slice-Datei "+file+" ist unlesbar (fail-closed)")
	}
	lines := splitLines(content)
	headingNo, level := closureHeadingLine(lines, re)
	if headingNo == 0 {
		// C3: kein passender Abschnitt — missing schließt thin/boilerplate aus,
		// ohne Abschnitt gibt es nichts zu messen.
		return closureFinding(file, 1, dir, model.ReasonClosureNoteMissing,
			"kein Closure-Notiz-Abschnitt (keine Überschrift passt auf "+
				cfg.Closure.EffectiveHeadingPattern()+")")
	}
	body := closureSectionProse(lines, headingNo, level)
	var out []model.Finding
	// C4a: Substanz — Satzende-Zeichen außerhalb der Fenced-Code-Blöcke.
	want := cfg.Closure.EffectiveMinSentences()
	if got := countSentenceEnds(body); got < want {
		out = append(out, closureFinding(file, headingNo, dir, model.ReasonClosureNoteThin,
			"Closure-Notiz trägt "+strconv.Itoa(got)+" Satzende-Zeichen außerhalb von Code-Blöcken, "+
				"verlangt sind "+strconv.Itoa(want))...)
	}
	// C4b: Floskel — literaler Teilstring, case-insensitiv; der erste Treffer
	// benennt die Meldung.
	lower := strings.ToLower(body)
	for _, phrase := range cfg.Closure.Boilerplate {
		if strings.Contains(lower, strings.ToLower(phrase)) {
			out = append(out, closureFinding(file, headingNo, dir, model.ReasonClosureNoteBoilerplate,
				"Closure-Notiz enthält die Floskel „"+phrase+"“")...)
			break
		}
	}
	// C4b: Platzhalter — opt-in; ohne den Schalter entfällt der Schritt und der
	// Befundsatz ist byte-identisch.
	if cfg.Closure.Placeholder {
		out = append(out, checkClosurePlaceholder(
			file, dir, content, headingNo, closureSectionEnd(lines, headingNo, level))...)
	}
	return out
}

// closureHeadingLine liefert die 1-basierte Zeilennummer der **ersten**
// Überschrift, deren getrimmte Fassung auf das Muster passt, samt ihrer Ebene.
// 0 ⇒ kein Treffer. Fenced-Code wird übersprungen — eine '#'-Zeile in einem
// Beispielblock ist keine Überschrift.
//
// Was als Überschrift zählt, entscheidet der **geteilte** ATX-Parser (derselbe,
// den `anchors`/`matrix` nutzen), nicht eine eigene '#'-Zählung: `#1 war ein
// Thema` ist Fließtext, keine H1. Eine eigene Heuristik hier hätte solche
// Zeilen als Überschrift gelesen und den Abschnitt vorzeitig beendet — R1-F-1.
func closureHeadingLine(lines []string, re *regexp.Regexp) (lineNo, level int) {
	inFence := false
	for i, raw := range lines {
		if FenceToggle(TrimFenceIndent(raw)) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lvl, _, ok := parseATXHeading(raw)
		if !ok {
			continue
		}
		if re.MatchString(strings.TrimSpace(raw)) {
			return i + 1, lvl
		}
	}
	return 0, 0
}

// closureSectionProse liefert den Abschnitts-Text ab der Zeile **nach** der
// Überschrift bis zur nächsten Überschrift **gleicher oder höherer** Ebene
// (exklusive) bzw. bis zum Dateiende — bereinigt um die Fenced-Code-Blöcke
// (Spez-Schritt C4). Eine tiefere Überschrift gehört noch zum Abschnitt.
// closureSectionEnd liefert die 1-basierte Zeilennummer der Überschrift, die den
// Abschnitt beendet (0 ⇒ Dateiende): die nächste echte ATX-Überschrift gleicher
// oder höherer Ebene außerhalb von Fenced-Code.
func closureSectionEnd(lines []string, headingNo, level int) int {
	inFence := false
	for i := headingNo; i < len(lines); i++ {
		if FenceToggle(TrimFenceIndent(lines[i])) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if lvl, _, ok := parseATXHeading(lines[i]); ok && lvl <= level {
			return i + 1
		}
	}
	return 0
}

// placeholderRE erkennt die Auszeichnungs-Form eines Vorlagen-Platzhalters
// (§DC-FA-PLAN-001.a Schritt C4b): eine öffnende Winkelklammer ohne
// vorausgehendes Wortzeichen und ohne Schrägstrich, deren erstes Inneres kein
// Whitespace, `!` oder `/` ist. Das Vorzeichen wird KONSUMIERT statt
// hineingeschaut — RE2 kennt keine Lookarounds.
var placeholderRE = regexp.MustCompile(`(^|[^\w/])<([^\s!/][^<>]{0,79})>`)

// htmlTagNames sind die Tag-Namen, die der Nachfilter aus C4b verwirft. Als
// Liste statt als Regex-Alternation: eine Liste liest und pflegt sich.
func htmlTagNames() map[string]bool {
	return map[string]bool{
		"a": true, "b": true, "i": true, "p": true, "br": true, "hr": true,
		"em": true, "strong": true, "code": true, "pre": true, "div": true,
		"span": true, "ul": true, "ol": true, "li": true, "h1": true, "h2": true,
		"h3": true, "h4": true, "h5": true, "h6": true, "td": true, "tr": true,
		"th": true, "img": true, "sup": true, "sub": true, "table": true,
		"thead": true, "tbody": true, "kbd": true, "blockquote": true,
		"details": true, "summary": true, "figure": true, "figcaption": true,
	}
}

// placeholderRejected meldet, ob ein Treffer von den Nachfiltern verworfen wird.
// Nachfilter sind bewusst Code und nicht Teil des Musters (ADR-0052): in die
// Regex gepresst wären sie unlesbar und einzeln unprüfbar.
func placeholderRejected(inner string) bool {
	if strings.Contains(inner, "://") || strings.Contains(inner, "@") {
		return true
	}
	first := inner
	if i := strings.IndexAny(first, " \t/"); i != -1 {
		first = first[:i]
	}
	return htmlTagNames()[strings.ToLower(first)]
}

// checkClosurePlaceholder sucht den ERSTEN Platzhalter im Abschnitt
// (§DC-FA-PLAN-001.a Schritt C4b). Gearbeitet wird auf den vorverarbeiteten
// Zeilen: Fences sind entfernt und Inline-Code-Spans positionserhaltend
// geleert — dieselbe geteilte Lexik wie überall, kein Nachbau.
func checkClosurePlaceholder(
	file, dir string, content []byte, headingNo, endNo int,
) []model.Finding {
	for _, ln := range PreprocessMarkdown(content) {
		if ln.No <= headingNo || (endNo != 0 && ln.No >= endNo) {
			continue
		}
		// Schleife statt FindAll, weil ein vom Nachfilter verworfener Treffer die
		// Suche nicht beenden darf. Weitergesucht wird auf dem Reststring; dessen
		// Anfang erfüllt die Vorzeichen-Bedingung über die ^-Alternative, sodass
		// ein unmittelbar folgender Platzhalter sichtbar bleibt.
		for rest := ln.Text; rest != ""; {
			loc := placeholderRE.FindStringSubmatchIndex(rest)
			if loc == nil {
				break
			}
			inner := rest[loc[4]:loc[5]]
			if placeholderRejected(inner) {
				rest = rest[loc[1]:]
				continue
			}
			hit := "<" + inner + ">"
			return closureFinding(file, ln.No, dir, model.ReasonClosureNotePlaceholder,
				"Closure-Notiz trägt einen unausgefüllten Platzhalter "+clipRunes(hit, 40))
		}
	}
	return nil
}

func closureSectionProse(lines []string, headingNo, level int) string {
	var b strings.Builder
	inFence := false
	for i := headingNo; i < len(lines); i++ {
		trimmed := TrimFenceIndent(lines[i])
		if FenceToggle(trimmed) {
			inFence = !inFence
			continue // Fence-Zeilen selbst zählen nicht
		}
		if inFence {
			continue
		}
		// Abschnitts-Ende nur an einer ECHTEN ATX-Überschrift (geteilter Parser)
		// gleicher oder höherer Ebene — `#1 war ein Thema` beendet nichts.
		if lvl, _, ok := parseATXHeading(lines[i]); ok && lvl <= level {
			break
		}
		b.WriteString(lines[i])
		b.WriteByte('\n')
	}
	return b.String()
}

// countSentenceEnds zählt die Satzende-Zeichen '.', '!' und '?' im bereits
// bereinigten Text. Bewusst simpel: die Prüfung ist eine Substanz-Untergrenze,
// keine Sprachanalyse — ein Abkürzungspunkt zählt mit, und das ist billiger und
// ehrlicher als eine Heuristik, die manchmal danebenliegt.
func countSentenceEnds(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' || s[i] == '!' || s[i] == '?' {
			n++
		}
	}
	return n
}

// closureFinding baut einen Closure-Note-Befund (Spez-Schritt C5).
func closureFinding(file string, line int, dir, reason, msg string) []model.Finding {
	return []model.Finding{{
		File: file, Line: line, Rule: "planning", Target: dir,
		Reason: reason, Message: msg,
	}}
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
