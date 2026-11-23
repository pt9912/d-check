package rules

// Modul linelength (DC-FA-LLEN-001, spec/spezifikation.md §DC-FA-LLEN-001.a):
// keine Zeile einer Markdown-Datei der Eingabemenge trägt mehr Codepunkte, als
// linelength.max-chars erlaubt. Hermetisch: das Modul liest die Dateien, die
// der Lauf ohnehin liest, und schreibt nie.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): die Eingabemenge ist der effektive
// Scan-Scope des Moduls minus exempt-paths; eine andere Datei liest das Modul
// nicht. Was Fenced-Code und was eine Tabellenzeile ist, beantworten
// proseLineSet und tableRowLine — mit deren Grenzen (Umschalten bei jeder
// Zaun-Zeile, ein offener Zaun nimmt den Rest der Datei aus, keine Trennzeile
// verlangt). Den offenen Zaun meldet spans, nicht dieses Modul.

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// ReasonLineCharsExceeded meldet eine Zeile mit mehr Codepunkten, als max-chars
// erlaubt (spec/spezifikation.md §4).
const ReasonLineCharsExceeded = "line-chars-exceeded"

// LinelengthCount ist die Bezugsmenge des Moduls: geprüfte Dateien, geprüfte
// Zeilen (nach den Ausnahmen (a) und (b)) und die davon über (c)
// ausgenommenen.
type LinelengthCount struct {
	Files       int
	Lines       int
	Unbreakable int
}

// Note liefert die Zeile der Bezugsmenge für die Zusammenfassung (Schritt 8):
// sie sagt, was gemessen wurde, nicht, dass jede Zeile des Repos gemessen ist.
func (c LinelengthCount) Note() string {
	return fmt.Sprintf("linelength: %d Datei(en), %d Zeile(n) geprüft, %d davon mit einem Wort über der Grenze ausgenommen",
		c.Files, c.Lines, c.Unbreakable)
}

// checkLinelengthConfig prüft vor dem Lauf, dass ein aktives Modul eine
// Grenze trägt (Schritt 1); der übrige Block ist am Config-Rand geprüft.
func checkLinelengthConfig(cfg model.LinelengthConfig) error {
	if cfg.MaxChars == nil {
		return fmt.Errorf("linelength: linelength.max-chars fehlt — das Modul ist aktiv und hätte keine Grenze (DC-FA-LLEN-001, fail-closed)")
	}
	return nil
}

// linelengthExempt meldet, ob die Datei über exempt-paths aus der
// Eingabemenge fällt (Schritt 2; Glob-Auswertung wie scan.ignore).
func linelengthExempt(file string, cfg model.LinelengthConfig) bool {
	return ignored(file, cfg.ExemptPaths)
}

// CheckLinelength misst die Zeilen einer Datei der Eingabemenge (Schritte 3
// bis 7) und liefert die Befunde in Zeilen-Reihenfolge samt dem Beitrag der
// Datei zur Bezugsmenge. Ein d-check:ignore-Marker wirkt nicht.
func CheckLinelength(file string, content []byte, cfg model.LinelengthConfig) ([]model.Finding, LinelengthCount) {
	count := LinelengthCount{Files: 1}
	limit := *cfg.MaxChars
	lines := textLines(content)
	prose := proseLineSet(content)
	var out []model.Finding
	for i, line := range lines {
		if cfg.CodeBlocksExempt() && !prose[i+1] {
			continue
		}
		if cfg.TablesExempt() && tableRowLine(lines, prose, i) {
			continue
		}
		count.Lines++
		length := utf8.RuneCountInString(line)
		if length <= limit {
			continue
		}
		if cfg.UnbreakableExempt() && longestWord(line) > limit {
			count.Unbreakable++
			continue
		}
		out = append(out, model.Finding{
			File: file, Line: i + 1, Rule: "linelength",
			Target: "max-chars=" + strconv.Itoa(limit), Reason: ReasonLineCharsExceeded,
			Message: "Zeile hat " + strconv.Itoa(length) + " Zeichen, erlaubt sind " + strconv.Itoa(limit),
		})
	}
	return out, count
}

// longestWord liefert die Codepunkt-Zahl des längsten Worts der Zeile (Schritt
// 6): ein Wort ist eine längste Folge ohne U+0020 und U+0009 — kein anderes
// Zeichen trennt, auch kein geschütztes Leerzeichen. Eine ungültige
// UTF-8-Folge zählt je Byte als ein Codepunkt.
func longestWord(line string) int {
	longest := 0
	for _, w := range strings.FieldsFunc(line, func(r rune) bool { return r == ' ' || r == '\t' }) {
		if c := utf8.RuneCountInString(w); c > longest {
			longest = c
		}
	}
	return longest
}
