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
	headingNo, hasActive, fail := planningActiveStatus(content, cfg, sliceDir)
	if fail != nil {
		return fail
	}
	hasSlices := planningHasSlices(fsys, sliceDir, cfg.EffectiveSliceGlob())
	if hasActive == hasSlices {
		return nil
	}
	var msg string
	if hasSlices {
		msg = "Slice(s) in " + sliceDir + ", aber die Roadmap-Sektion „" + cfg.EffectiveHeading() +
			"“ trägt den Ruhe-Marker „" + cfg.EffectiveMarker() + "“ — die Sektion muss die Arbeit benennen"
	} else {
		msg = "kein Slice in " + sliceDir + ", aber die Roadmap-Sektion „" + cfg.EffectiveHeading() +
			"“ trägt den Ruhe-Marker „" + cfg.EffectiveMarker() + "“ nicht — er gehört dorthin"
	}
	return planningDrift(cfg.Roadmap, headingNo, sliceDir, msg)
}

// planningActiveStatus bestimmt den Aktiv-Status EINMAL: Zeile der kanonischen
// Überschrift und ob ihr Block den Ruhe-Marker trägt. Jede Fähigkeit, die den
// Status braucht, ruft diese Stelle auf — eine zweite Antwort auf dieselbe
// Frage wäre ein Defekt (ADR-0054). Ist der Status nicht bestimmbar (Überschrift
// fehlt oder mehrfach), liefert fail den fail-closed-Befund.
func planningActiveStatus(
	content []byte, cfg model.PlanningConfig, sliceDir string,
) (headingNo int, hasActive bool, fail []model.Finding) {
	lines := splitLines(content)
	prose := proseLineSet(content)
	no, count := planningHeadingLine(lines, prose, cfg.EffectiveHeading())
	if count == 0 {
		return 0, false, planningDrift(cfg.Roadmap, 1, sliceDir,
			"kanonische Überschrift „"+cfg.EffectiveHeading()+"“ fehlt in "+cfg.Roadmap+
				" — Aktiv-Status nicht bestimmbar (fail-closed)")
	}
	if count > 1 {
		return 0, false, planningDrift(cfg.Roadmap, no, sliceDir,
			"kanonische Überschrift „"+cfg.EffectiveHeading()+"“ kommt mehrfach in "+cfg.Roadmap+
				" vor — Aktiv-Status mehrdeutig (fail-closed)")
	}
	return no, !planningBlockHasMarker(lines, prose, no, cfg.EffectiveMarker()), nil
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
	// der RTM-Anforderungsquellen.
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
	headingNo, level, ambiguousAt := closureHeadingLine(lines, re)
	if ambiguousAt != 0 {
		// Mehrdeutig ⇒ NICHT messen: ein zweiter Abschnitt ist der typische Rest
		// einer Vorlage, und eine Satzzahl über den falschen sagt nichts.
		return closureFinding(file, ambiguousAt, dir, model.ReasonClosureNoteAmbiguous,
			"mehrere Abschnitte passen auf "+cfg.Closure.EffectiveHeadingPattern()+
				" — ohne eindeutigen Abschnitt wird nicht gemessen")
	}
	if headingNo == 0 {
		// C3: kein passender Abschnitt — missing schließt thin/boilerplate aus,
		// ohne Abschnitt gibt es nichts zu messen.
		return closureFinding(file, 1, dir, model.ReasonClosureNoteMissing,
			"kein Closure-Notiz-Abschnitt (keine Überschrift passt auf "+
				cfg.Closure.EffectiveHeadingPattern()+")")
	}
	body := SectionProse(content, lines, headingNo, level)
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
		if containsWord(lower, strings.ToLower(phrase)) {
			out = append(out, closureFinding(file, headingNo, dir, model.ReasonClosureNoteBoilerplate,
				"Closure-Notiz enthält die Floskel „"+phrase+"“")...)
			break
		}
	}
	// C4b: Platzhalter — opt-in; ohne den Schalter entfällt der Schritt und der
	// Befundsatz ist byte-identisch.
	if cfg.Closure.Placeholder {
		out = append(out, checkClosurePlaceholder(
			file, dir, content, headingNo, SectionEnd(lines, headingNo, level))...)
	}
	return out
}

// closureHeadingLine liefert die Abschnitts-Überschrift der Closure-Notiz: die
// 1-basierte Zeilennummer der **ersten** Überschrift, deren getrimmte Fassung
// auf das Muster passt, samt ihrer Ebene (0 ⇒ kein Treffer), und in ambiguousAt
// die Zeile des **zweiten** Treffers (0 ⇒ eindeutig) — ohne eindeutigen
// Abschnitt sagt eine Messung nichts (§DC-FA-PLAN-001.a Schritt C3).
//
// Fenced-Code wird übersprungen (eine '#'-Zeile in einem Beispielblock ist
// keine Überschrift), und was als Überschrift zählt, entscheidet der
// **geteilte** ATX-Parser (derselbe, den `anchors`/`matrix` nutzen), nicht eine
// eigene '#'-Zählung: `#1 war ein Thema` ist Fließtext, keine H1.
func closureHeadingLine(lines []string, re *regexp.Regexp) (lineNo, level, ambiguousAt int) {
	heads := FindSectionHeads(lines, func(raw string) bool {
		return re.MatchString(strings.TrimSpace(raw))
	})
	if len(heads) == 0 {
		return 0, 0, 0
	}
	if len(heads) > 1 {
		return heads[0].Line, heads[0].Level, heads[1].Line
	}
	return heads[0].Line, heads[0].Level, 0
}


// placeholderRE erkennt die Auszeichnungs-Form eines Vorlagen-Platzhalters
// (§DC-FA-PLAN-001.a Schritt C4b): eine öffnende Winkelklammer ohne
// vorausgehendes Wortzeichen und ohne Schrägstrich, deren Inneres mit einem
// anderen Zeichen als `!` oder `/` beginnt und KEIN Whitespace enthält. Das
// Vorzeichen wird KONSUMIERT statt hineingeschaut — RE2 kennt keine
// Lookarounds.
var placeholderRE = regexp.MustCompile(`(^|[^\w/])<([^\s!/<>][^\s<>]{0,79})>`)

// htmlTagNames sind die Tag-Namen, die der Nachfilter aus C4b verwirft. Als
// Liste statt als Regex-Alternation: eine Liste liest und pflegt sich.
func htmlTagNames() map[string]bool {
	names := []string{
		"a", "abbr", "address", "area", "article", "aside", "audio",
		"b", "base", "bdi", "bdo", "blockquote", "body", "br", "button",
		"canvas", "caption", "cite", "code", "col", "colgroup",
		"data", "datalist", "dd", "del", "details", "dfn", "dialog", "div", "dl", "dt",
		"em", "embed", "fieldset", "figcaption", "figure", "footer", "form",
		"h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hgroup", "hr", "html",
		"i", "iframe", "img", "input", "ins", "kbd", "label", "legend", "li", "link",
		"main", "map", "mark", "menu", "meta", "meter", "nav", "noscript",
		"object", "ol", "optgroup", "option", "output",
		"p", "param", "picture", "pre", "progress", "q",
		"rp", "rt", "ruby", "s", "samp", "script", "section", "select", "slot",
		"small", "source", "span", "strong", "style", "sub", "summary", "sup", "svg",
		"table", "tbody", "td", "template", "textarea", "tfoot", "th", "thead",
		"time", "title", "tr", "track", "u", "ul", "var", "video", "wbr",
	}
	out := make(map[string]bool, len(names))
	for _, n := range names {
		out[n] = true
	}
	return out
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

// containsWord meldet, ob phrase in s als eigenständiger Wortlaut vorkommt:
// unmittelbar davor und dahinter steht kein Wortzeichen (§DC-FA-PLAN-001.a
// Schritt C4). Index-Suche mit Nachbar-Prüfung statt einer Regex je Phrase —
// das kommt ohne Kompilierung aus und braucht keine Maskierung der Phrase.
// Beide Argumente sind bereits kleingeschrieben.
func containsWord(s, phrase string) bool {
	if phrase == "" {
		return false
	}
	for i := 0; ; {
		j := strings.Index(s[i:], phrase)
		if j == -1 {
			return false
		}
		start := i + j
		end := start + len(phrase)
		if !isWordByte(s, start-1) && !isWordByte(s, end) {
			return true
		}
		i = start + 1
	}
}

// isWordByte meldet, ob an Position i ein ASCII-Wortzeichen steht; außerhalb
// des Strings ist die Antwort nein (Zeilen- und Textränder sind Grenzen).
func isWordByte(s string, i int) bool {
	if i < 0 || i >= len(s) {
		return false
	}
	c := s[i]
	return c == '_' || ('0' <= c && c <= '9') ||
		('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
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
			// Winkelklammer-Linkziel `](<ziel>)` ist Markdown-Syntax, keine
			// Vorlage — erkennbar an der oeffnenden Klammer als Vorzeichen.
			if rest[loc[2]:loc[3]] == "(" {
				rest = rest[loc[1]:]
				continue
			}
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


// countSentenceEnds zählt die Satzende-Zeichen '.', '!' und '?' im bereits
// bereinigten Text. Bewusst simpel: die Prüfung ist eine Substanz-Untergrenze,
// keine Sprachanalyse — ein Abkürzungspunkt zählt mit, und das ist billiger und
// ehrlicher als eine Heuristik, die manchmal danebenliegt.
func countSentenceEnds(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '.' && s[i] != '!' && s[i] != '?' {
			continue
		}
		// \r gehoert dazu: bei CRLF steht es zwischen Satzende und Zeilenumbruch,
		// und ohne diesen Fall zaehlte in einer CRLF-Arbeitskopie kein einziges
		// zeilenschliessendes Satzende.
		if i+1 < len(s) && s[i+1] != ' ' && s[i+1] != '\t' && s[i+1] != '\n' && s[i+1] != '\r' {
			continue
		}
		n++
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
func planningHeadingLine(lines []string, prose map[int]bool, heading string) (first, count int) {
	for i, raw := range lines {
		if !prose[i+1] {
			continue // Fence-Inneres: eine Beispiel-Überschrift ist keine
		}
		if strings.TrimRight(raw, " \t\r") == heading {
			count++
			if first == 0 {
				first = i + 1
			}
		}
	}
	return first, count
}

// planningBlockEnd liefert die EINE Block-Grenze des Aktiv-Status-Abschnitts
// (exklusives Zeilen-Ende): Marker-Suche und Kennungs-Scan (mode: many) lesen
// denselben Ausschnitt — zwei Grenzen wären zwei Antworten auf dieselbe Frage
// (ADR-0054). Wo der Block endet, entscheidet die geteilte
// Abschnitts-Mechanik — nicht ein roher Präfix-Vergleich: eine eingerückte
// oder tab-getrennte H2 ist eine Überschrift, eine H1 beendet einen
// H2-Abschnitt ebenfalls.
func planningBlockEnd(lines []string, headingNo int) int {
	level, _, ok := parseATXHeading(lines[headingNo-1])
	if !ok {
		level = 2
	}
	end := SectionEnd(lines, headingNo, level)
	if end == 0 || end > len(lines)+1 {
		end = len(lines) + 1
	}
	return end
}

// planningBlockHasMarker sucht den Ruhe-Marker (literaler Teilstring) NUR im
// Aktiv-Status-Block: ab der Zeile **nach** der Überschrift (headingNo) bis zur
// Block-Grenze. So verfälscht ein erklärender Verweis anderswo den Status nicht.
func planningBlockHasMarker(lines []string, prose map[int]bool, headingNo int, marker string) bool {
	end := planningBlockEnd(lines, headingNo)
	for i := headingNo; i < end-1; i++ {
		if !prose[i+1] {
			continue // Fence-Inneres beendet den Block nicht und trägt keinen Marker
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
