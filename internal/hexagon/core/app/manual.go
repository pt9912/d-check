package app

import (
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/rules"
)

// ManualDoc ist ein mitgeliefertes Dokument: sein Pfad im Repo und sein
// Inhalt.
type ManualDoc struct {
	Path    string
	Content string
}

// ManualSections liefert die Abschnitte der Dokumente, deren Überschrift den
// Begriff nennt (DC-FA-CLI-013.a Schritte 2–5): Treffer sind echte
// ATX-Überschriften außerhalb von Fenced-Code — dieselbe Erkennung wie bei
// structure —, verglichen in Kleinbuchstaben; ein Abschnitt reicht bis vor die
// nächste Überschrift gleicher oder höherer Ebene. Ein Treffer innerhalb eines
// bereits ausgegebenen Abschnitts wird nicht wiederholt. found ist false, wenn
// keine Überschrift den Begriff nennt.
func ManualSections(docs []ManualDoc, term string) (out string, found bool) {
	want := strings.ToLower(strings.TrimSpace(term))
	var parts []string
	for _, d := range docs {
		lines := strings.Split(d.Content, "\n")
		heads := rules.FindSectionHeads(lines, func(raw string) bool {
			return strings.Contains(strings.ToLower(manualHeadingText(raw)), want)
		})
		until := 0
		for _, h := range heads {
			if h.Line < until {
				continue
			}
			until = rules.SectionEnd(lines, h.Line, h.Level)
			stop := until - 1
			if until == 0 {
				until, stop = len(lines)+1, len(lines)
			}
			for stop > h.Line && strings.TrimSpace(lines[stop-1]) == "" {
				stop--
			}
			body := strings.Join(lines[h.Line-1:stop], "\n")
			parts = append(parts, "==> "+d.Path+":"+strconv.Itoa(h.Line)+"\n"+body+"\n")
		}
	}
	return strings.Join(parts, "\n"), len(parts) > 0
}

// manualHeadingText ist der Text einer ATX-Überschrift ohne die führende
// #-Folge, mit der Erkennung des Kerns.
func manualHeadingText(raw string) string {
	text, _ := rules.HeadingText(raw)
	return text
}

// ManualTitle ist der Text der ersten Überschrift erster Ebene außerhalb von
// Fenced-Code — der Begriff, der das ganze Dokument liefert.
func ManualTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, h := range rules.FindSectionHeads(lines, func(string) bool { return true }) {
		if h.Level == 1 {
			return manualHeadingText(lines[h.Line-1])
		}
	}
	return ""
}
