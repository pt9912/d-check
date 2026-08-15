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
		return []model.Finding{structureFinding(r, file, 1, model.ReasonSectionMissing,
			"Datei ist unlesbar (fail-closed)")}
	}
	lines := splitLines(content)
	heads := FindSectionHeads(lines, structureMatcher(r))
	if len(heads) == 0 {
		return []model.Finding{structureFinding(r, file, 1, model.ReasonSectionMissing,
			"kein Abschnitt passt auf den Selektor")}
	}
	if r.EffectiveSections() == "one" && len(heads) > 1 {
		// Abbruch für DIESE Datei innerhalb DIESER Regel: ohne eindeutigen
		// Abschnitt sagt eine Messung nichts.
		return []model.Finding{structureFinding(r, file, heads[1].Line, model.ReasonSectionAmbiguous,
			"Abschnitt kommt mehrfach vor, erwartet ist genau einer (sections: one)")}
	}
	// In `one` ist nach der Kardinalitäts-Prüfung genau ein Treffer übrig,
	// `each` misst jeden.
	var out []model.Finding
	for _, h := range heads {
		body := SectionProse(content, lines, h.Line, h.Level)
		out = append(out, structureConditions(r, file, h.Line, body)...)
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

// structureConditions prüft die sechs Bedingungen auf dem bereinigten
// Abschnitts-Text. Jede hat einen eigenen Grund-Code, damit zwei Verletzungen
// desselben Abschnitts nicht unter der Befund-Deduplikation zusammenfallen.
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
		if got := countTaskItems(body); got > *r.MaxTasks {
			add(model.ReasonSectionOversized, "Abschnitt trägt "+strconv.Itoa(got)+
				" Task-Items, erlaubt sind "+strconv.Itoa(*r.MaxTasks))
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

func structureFinding(r model.StructureRule, file string, line int, reason, msg string) model.Finding {
	return model.Finding{
		File: file, Line: line, Rule: "structure",
		Target: r.Identity(), Reason: reason, Message: msg,
	}
}

// taskItemRE erkennt ein Task-Item: nach optionalem Whitespace ein
// Listen-Marker (-, *, + oder <ziffern>.), Whitespace und [ ] bzw. [x]/[X].
var taskItemRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[[ xX]\]`)

func countTaskItems(body string) int {
	n := 0
	for _, l := range strings.Split(body, "\n") {
		if taskItemRE.MatchString(l) {
			n++
		}
	}
	return n
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

