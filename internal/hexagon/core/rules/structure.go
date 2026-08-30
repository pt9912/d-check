package rules

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckStructure prüft Struktur-Invarianten innerhalb eines Dokuments
// (DC-FA-STRUCT-001, §DC-FA-STRUCT-001.a). Post-Pass wie planning: die Regeln
// benennen ihre Dateien selbst über eigene Globs, unabhängig von
// scan.roots/scan.ignore — deshalb kennt das Modul kein <modul>.scope.
func CheckStructure(fsys driven.Filesystem, rules []model.StructureRule) []model.Finding {
	if len(rules) == 0 {
		return nil
	}
	all, err := structureTree(fsys)
	if err != nil {
		// fail-closed wie jede andere Randbedingung des Moduls: ohne Dateibaum
		// ist keine Regel messbar, und ein leerer Befundsatz wäre von "alles
		// erfüllt" nicht zu unterscheiden. Je Regel ein Befund, damit die
		// Identität der unmessbaren Behauptung erhalten bleibt.
		out := make([]model.Finding, 0, len(rules))
		for _, r := range rules {
			out = append(out, model.Finding{
				File: r.Files, Line: 1, Rule: "structure", Target: r.Identity(),
				Reason:  model.ReasonSectionMissing,
				Message: "Dateibaum nicht lesbar (" + err.Error() + ") — Regel nicht messbar (fail-closed)",
			})
		}
		return out
	}
	var out []model.Finding
	for _, r := range rules {
		out = append(out, checkStructureRule(fsys, r, all)...)
	}
	return out
}

// structureTree listet alle Markdown-Dateien des Baums (SKIP_DIRS gelten, die
// Scan-Wurzeln nicht) — die Grundmenge, gegen die die Regel-Globs laufen.
func structureTree(fsys driven.Filesystem) ([]string, error) {
	var out []string
	if err := walkMarkdown(fsys, "", nil, &out); err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// checkStructureRule wertet eine Regel aus: Kandidaten (Schritt 2), dann je
// Datei Abschnitts-Findung, Kardinalität und Bedingungen (Schritte 3–6).
func checkStructureRule(fsys driven.Filesystem, r model.StructureRule, all []string) []model.Finding {
	var cands []string
	for _, f := range all {
		if !matchGlob(r.Files, f) {
			continue
		}
		if structureExempt(r, f) {
			continue
		}
		cands = append(cands, f)
	}
	// Nullmengen-Härte: eine Regel zu setzen IST die Behauptung, dass sie Dateien
	// trifft — auch dann, wenn erst exempt-paths die Menge geleert hat.
	if len(cands) == 0 {
		return []model.Finding{{
			File: r.Files, Line: 1, Rule: "structure", Target: r.Identity(),
			Reason: model.ReasonSectionMissing,
			Message: "Regel trifft keine Datei (auch nach Abzug von exempt-paths) — " +
				"das Gate liefe leer",
		}}
	}
	var out []model.Finding
	for _, f := range cands {
		out = append(out, checkStructureFile(fsys, r, f)...)
	}
	return out
}

// structureExemptSections zieht die von exempt-section-pattern getroffenen
// Abschnitte ab. Das Muster sieht DIESELBE Zeichenkette wie section-pattern --
// die getrimmte Ueberschriften-Zeile EINSCHLIESSLICH der #-Folge. Zwei RE2 in
// einer Regel mit zwei verschiedenen Zielen waeren die Falle: das zweite
// Muster, analog zum ersten geschrieben, traefe still nichts (ADR-0075).
func structureExemptSections(r model.StructureRule, lines []string, heads []SectionHead) []SectionHead {
	re := regexp.MustCompile(r.ExemptSectionPattern)
	out := make([]SectionHead, 0, len(heads))
	for _, h := range heads {
		if re.MatchString(strings.TrimSpace(lines[h.Line-1])) {
			continue
		}
		out = append(out, h)
	}
	return out
}

func structureExempt(r model.StructureRule, file string) bool {
	for _, g := range r.ExemptPaths {
		if matchGlob(g, file) {
			return true
		}
	}
	return false
}

// checkStructureFile prüft eine Kandidaten-Datei gegen eine Regel.
func checkStructureFile(fsys driven.Filesystem, r model.StructureRule, file string) []model.Finding {
	content, err := fsys.ReadFile(file)
	if err != nil {
		// NICHT durch MessageFor: hier hat die Regel gar nicht gemessen. Ein
		// verfasster Hinweis auf die gehuetete Zusage verdraengte die
		// fail-closed-Ursache und zeigte auf eine Bedingung, die nie geprueft
		// wurde (ADR-0073, benannte Grenze).
		return []model.Finding{structureRawFinding(r, file, 1, model.ReasonSectionMissing,
			"Datei ist unlesbar (fail-closed)")}
	}
	lines := splitLines(content)
	heads := FindSectionHeads(lines, structureMatcher(r))
	if len(heads) == 0 {
		return []model.Finding{structureFinding(r, file, 1, model.ReasonSectionMissing,
			"kein Abschnitt passt auf den Selektor")}
	}
	// Die Ausnahme erklaert die GRUNDMENGE und laeuft deshalb vor der
	// Kardinalitaets-Pruefung: was ausgenommen ist, kann `sections: one` nicht
	// mehrdeutig machen. Dieselbe Reihenfolge wie bei exempt-paths, eine
	// Granularitaetsstufe tiefer (ADR-0075).
	if n := len(heads); r.ExemptSectionPattern != "" {
		heads = structureExemptSections(r, lines, heads)
		// Nullmengen-Haerte wie bei exempt-paths: ein Ventil, das die Menge
		// leert, schaltet die Regel NICHT still ab. Die Meldung nennt den
		// Schluessel, der es tat -- sonst suchte der Leser am Selektor.
		if len(heads) == 0 {
			// RAW, nicht ueber MessageFor: das ist die LEER LAUFENDE REGEL,
			// die ADR-0073 ausdruecklich vom verfassten Hinweis ausnimmt --
			// hier hat die Regel nicht gemessen, und ein Hinweis auf die
			// gehuetete Zusage verdraengte die Ursache. Ohne diesen Zug
			// loeschte ein `hint` genau die Diagnose, die das zu breite
			// Muster sichtbar macht.
			return []model.Finding{structureRawFinding(r, file, 1, model.ReasonSectionMissing,
				"alle "+strconv.Itoa(n)+" passenden Abschnitte sind von "+
					"exempt-section-pattern ausgenommen — die Regel liefe leer")}
		}
	}
	if r.EffectiveSections() == "one" && len(heads) > 1 {
		// Abbruch für DIESE Datei innerhalb DIESER Regel: ohne eindeutigen
		// Abschnitt sagt eine Messung nichts.
		return []model.Finding{structureFinding(r, file, heads[1].Line, model.ReasonSectionAmbiguous,
			"Abschnitt kommt mehrfach vor, erwartet ist genau einer (sections: one)")}
	}
	// In `one` ist nach der Kardinalitäts-Prüfung genau ein Treffer übrig,
	// `each` misst jeden.
	var prose map[int]bool
	if r.Table.Order != "" || len(r.Table.Columns) > 0 {
		// Die beiden Tabellen-Bedingungen brauchen die fence-bewusste
		// Tabellenzeilen-Auswahl auf den ROHEN Zeilen (ADR-0057, ADR-0069).
		prose = proseLineSet(content)
	}
	var out []model.Finding
	for _, h := range heads {
		body := SectionProse(content, lines, h.Line, h.Level)
		out = append(out, structureConditions(r, file, h.Line, body)...)
		if r.Table.Order != "" {
			out = append(out, structureTableOrder(r, file, lines, prose, h.Line, h.Level)...)
		}
		// Ein Lauf je Spalte: die Befunde bleiben damit spaltenweise
		// gruppiert, so wie sie es unter den frueheren Ein-Spalten-Regeln
		// waren (ADR-0070 aendert die Form, nicht die Ausgabe).
		for _, c := range r.Table.Columns {
			out = append(out, structureCellMax(r, c, file, lines, prose, h.Line, h.Level)...)
		}
		if r.HeadingsMatch != "" {
			out = append(out, structureHeadings(r, file, lines, h)...)
		}
	}
	return out
}

// structureMatcher liefert das Abschnitts-Prädikat: Klartext-Vergleich der
// getrimmten Überschriften-Zeile (einschließlich der `#`-Folge) oder RE2.
func structureMatcher(r model.StructureRule) func(string) bool {
	if r.Section != "" {
		want := r.Section
		return func(raw string) bool { return strings.TrimSpace(raw) == want }
	}
	re := regexp.MustCompile(r.SectionPattern)
	return func(raw string) bool { return re.MatchString(strings.TrimSpace(raw)) }
}

// structureConditions prüft die sechs Prosa-Bedingungen auf dem bereinigten
// Abschnitts-Text. Jede hat einen eigenen Grund-Code, damit zwei Verletzungen
// desselben Abschnitts nicht unter der Befund-Deduplikation zusammenfallen.
// Zwei Bedingungen lesen einen anderen Text und leben anderswo: die
// Chronologie-Monotonie auf den rohen Abschnitts-Zeilen
// (structure_tableorder.go, ADR-0057) und die Überschriften-Bedingung auf den
// Überschriften selbst (structureHeadings).
func structureConditions(r model.StructureRule, file string, line int, body string) []model.Finding {
	var out []model.Finding
	add := func(reason, msg string) {
		out = append(out, structureFinding(r, file, line, reason, msg))
	}
	if r.NonEmpty && strings.TrimSpace(body) == "" {
		add(model.ReasonSectionEmpty, "Abschnitt ist leer")
	}
	if r.MinSentences != nil {
		if got := countSentenceEnds(body); got < *r.MinSentences {
			add(model.ReasonSectionThin, "Abschnitt trägt "+strconv.Itoa(got)+
				" Satzende-Zeichen, verlangt sind "+strconv.Itoa(*r.MinSentences))
		}
	}
	if r.MaxTasks != nil {
		got, ignoriert := countTaskItems(body, taskIgnoreRE(r))
		if got > *r.MaxTasks {
			// Die Zahl der ignorierten Items steht NUR dabei, wenn ein Muster
			// gesetzt ist: sonst waere die Meldung jeder Bestandsregel eine
			// andere. Sie steht dann AUCH bei null -- ein Muster, das nichts
			// trifft, ist eine Zusage, die nicht wirkt, und das gehoert
			// sichtbar (ADR-0075).
			zusatz := ""
			if r.TasksIgnorePattern != "" {
				zusatz = " (" + strconv.Itoa(ignoriert) + " ignoriert)"
			}
			add(model.ReasonSectionOversized, "Abschnitt trägt "+strconv.Itoa(got)+
				" Task-Items"+zusatz+", erlaubt sind "+strconv.Itoa(*r.MaxTasks))
		}
	}
	if r.ForbidPattern != "" && regexp.MustCompile(r.ForbidPattern).MatchString(body) {
		add(model.ReasonSectionForbidden, "verbotenes Muster trifft: "+r.ForbidPattern)
	}
	if r.RequirePattern != "" && !regexp.MustCompile(r.RequirePattern).MatchString(body) {
		add(model.ReasonSectionPatternMissing, "gefordertes Muster fehlt: "+r.RequirePattern)
	}
	for _, m := range r.RequireAll {
		if !hasMarker(body, m) {
			add(model.ReasonSectionMarkerMissing, "geforderte Marke fehlt: "+m)
			break
		}
	}
	return out
}

// structureHeadings prüft die Überschriften INNERHALB des Abschnitts gegen
// headings-match (§DC-FA-STRUCT-001.a Schritt 6). Sie liest als zweite
// Bedingung neben der Chronologie nicht den bereinigten Abschnitts-Text,
// sondern die Überschriften selbst — über die geteilte Erkennung, nicht über
// ein eigenes Muster. Ein Befund je verletzender Überschrift, auf IHRER
// Zeile; dort ist die Reparatur.
func structureHeadings(r model.StructureRule, file string, lines []string, h SectionHead) []model.Finding {
	re := regexp.MustCompile(r.HeadingsMatch)
	var out []model.Finding
	for _, sh := range SectionHeadings(lines, h.Line, h.Level, r.EffectiveHeadingsLevel(h.Level)) {
		if re.MatchString(sh.Text) {
			continue
		}
		out = append(out, structureFinding(r, file, sh.Line, model.ReasonSectionHeadingMismatch,
			"Überschrift \""+sh.Text+"\" genügt dem geforderten Muster nicht: "+r.HeadingsMatch))
	}
	return out
}

// structureRawFinding traegt die MODUL-EIGENE Meldung, ohne den verfassten
// Hinweis der Regel. Sie ist den Befunden vorbehalten, die KEINE Bedingung
// verletzen — dort ist die Ursache die Nachricht (ADR-0073).
func structureRawFinding(r model.StructureRule, file string, line int, reason, msg string) model.Finding {
	return model.Finding{
		File: file, Line: line, Rule: "structure",
		Target: r.Identity(), Reason: reason, Message: msg,
	}
}

func structureFinding(r model.StructureRule, file string, line int, reason, msg string) model.Finding {
	return model.Finding{
		File: file, Line: line, Rule: "structure",
		Target: r.Identity(), Reason: reason, Message: r.MessageFor(msg),
	}
}

// taskItemRE erkennt ein Task-Item: nach optionalem Whitespace ein
// Listen-Marker (-, *, + oder <ziffern>.), Whitespace und [ ] bzw. [x]/[X].
var taskItemRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[[ xX]\]`)

// countTaskItems zaehlt die Task-Items des bereinigten Abschnitts-Textes und
// meldet daneben, wie viele davon `ignore` herausgenommen hat.
//
// ABGRENZUNG, welche Zeichenkette `ignore` sieht: den ITEM-TEXT hinter
// Listen-Marker und Checkbox, NICHT die rohe Zeile. Gegen die rohe Zeile
// bezeichnete `^` immer den Listen-Marker; die verankerte Muster-Form waere
// damit unschreibbar, und ein freies Muster nahm an 444 realen Items zwei
// echte Liefer-Zusagen still mit (ADR-0075). Das ist die bewusste Asymmetrie
// zu exempt-section-pattern, das die ROHE Ueberschriften-Zeile sieht: dort
// gibt es mit section-pattern bereits einen Bezugspunkt, hier gibt es keinen.
func countTaskItems(body string, ignore *regexp.Regexp) (gezaehlt, ignoriert int) {
	for _, l := range strings.Split(body, "\n") {
		marker := taskItemRE.FindString(l)
		if marker == "" {
			continue
		}
		if ignore != nil && ignore.MatchString(strings.TrimSpace(l[len(marker):])) {
			ignoriert++
			continue
		}
		gezaehlt++
	}
	return gezaehlt, ignoriert
}

// taskIgnoreRE liefert das kompilierte Ignorier-Muster der Regel oder nil.
// Das Kompilat ist am Config-Rand bereits geprueft (Exit 2).
func taskIgnoreRE(r model.StructureRule) *regexp.Regexp {
	if r.TasksIgnorePattern == "" {
		return nil
	}
	return regexp.MustCompile(r.TasksIgnorePattern)
}

// markerRE erkennt eine Auszeichnungs-Marke am Zeilen-Anfang: nach optionalem
// Listen-Marker und Whitespace ein hervorgehobener Textlauf.
var markerRE = regexp.MustCompile(`^[ \t]*(?:(?:[-*+]|[0-9]+\.)[ \t]+)?\*\*([^*]+)\*\*`)

// hasMarker meldet, ob die Marke m als hervorgehobener Textlauf am
// Zeilen-Anfang vorkommt. Der Inhalt muss mit m beginnen und dort enden oder
// mit einem nicht-alphanumerischen Zeichen weitergehen — `**M:**` und
// `**M (Zusatz):**` treffen, `M` im Fließtext nicht, `**Marker**` nicht auf `M`.
func hasMarker(body, m string) bool {
	for _, l := range strings.Split(body, "\n") {
		sub := markerRE.FindStringSubmatch(l)
		if sub == nil {
			continue
		}
		// CutPrefix statt HasPrefix+Slice: eine Operation, damit Prüfung und
		// Schnitt nicht auseinanderlaufen können.
		rest, ok := strings.CutPrefix(sub[1], m)
		if !ok {
			continue
		}
		if rest == "" || !isWordRune(rest) {
			return true
		}
	}
	return false
}

// isWordRune meldet, ob s mit einem Buchstaben oder einer Ziffer beginnt —
// unicode-weit. Anders als bei der Floskel-Wortgrenze ist hier die ASCII-Menge
// falsch: `**Ergebnisüberblick:**` setzt die Marke `Ergebnis` fort, und ein
// Umlaut ist genauso Wort-Fortsetzung wie ein `s`.
func isWordRune(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

