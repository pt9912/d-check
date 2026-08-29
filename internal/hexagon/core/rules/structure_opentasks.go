package rules

// Offene Task-Items (DC-FA-STRUCT-001 §Offene Task-Items, ADR-0074): die
// dritte structure-Bedingung, die die ROHEN Abschnitts-Zeilen liest.
//
// WARUM ROH: die Bereinigung aus Schritt 5 leert Inline-Code ABSATZWEISE. Ein
// einzelner ueberzaehliger Backtick im Absatz macht damit alles dahinter
// unsichtbar — auf dem bereinigten Text ist eine Zusage ueber offene Haken
// deshalb nicht haltbar. Der Preis der Paarung ist fuer Prosa richtig
// (ADR-0059) und kippt fuer eine Vorbedingung.
//
// WARUM DIE MODUL-LEXIK: ein forbid-pattern in der Konfiguration deckt nur die
// Bullet-Form, die sein Autor aufgeschrieben hat. taskItemRE kennt alle vier
// Listen-Marker; diese Bedingung teilt seine Form und kann den Fehler nicht
// machen.
//
// GRENZE: roh heisst roh. Ein Task-Item in einer Inline-Code-Spanne zaehlt mit;
// nur der FENCE bleibt aussen vor, ueber dieselbe fence-bewusste Zeilen-Auswahl,
// die auch die Tabellen-Bedingungen benutzen.

import (
	"regexp"
	"strconv"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// openTaskItemRE erkennt ein OFFENES Task-Item: dieselbe Listen-Lexik wie
// taskItemRE (Bindestrich, Stern, Plus, geordnete Liste; fuehrender Weissraum,
// Leerzeichen ODER Tab als Trenner), nur mit leerer Box.
var openTaskItemRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[ \]`)

// structureOpenTasks meldet je offenem Task-Item einen Befund auf DESSEN Zeile
// — die Reparatur ist dort, wo der Haken steht, nicht an der
// Abschnitts-Ueberschrift.
func structureOpenTasks(
	r model.StructureRule, file string, lines []string, prose map[int]bool, headingNo, level int,
) []model.Finding {
	end := SectionEnd(lines, headingNo, level)
	if end == 0 {
		end = len(lines) + 1
	}
	var hits []int
	for i := headingNo; i < end-1; i++ {
		if prose[i+1] && openTaskItemRE.MatchString(lines[i]) {
			hits = append(hits, i+1)
		}
	}
	if len(hits) <= *r.MaxOpenTasks {
		return nil
	}
	out := make([]model.Finding, 0, len(hits))
	for _, no := range hits {
		out = append(out, structureFinding(r, file, no, model.ReasonSectionTasksOpen,
			"offenes Task-Item; erlaubt sind "+strconv.Itoa(*r.MaxOpenTasks)+
				", gefunden "+strconv.Itoa(len(hits))))
	}
	return out
}
