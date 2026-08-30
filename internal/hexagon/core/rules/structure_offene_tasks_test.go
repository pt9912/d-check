package rules

import (
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

func rohRule() model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", SectionPattern: `^## DoD`, Sections: "each",
		MaxOpenTasks: ptr(0),
	}
}

// DIE BLINDSTELLE IST GESCHLOSSEN — und der Vorzustand steht in derselben
// Funktion: die BEREINIGT lesende Form faellt an genau diesem Backtick auf
// null Befunde, die rohe nicht. Ein Regressions-Test ohne belegte Regression
// ist keiner (BEO-023); die Klasse selbst ist BEO-016.
func TestMaxOpenTasks_BacktickSchaltetNichtAb(t *testing.T) {
	// Ein Backtick VOR dem Haken und einer dahinter, im selben Absatz: die
	// absatzweise Paarung umschliesst das Item und leert es.
	body := "# D\n\n## DoD\n\nEin `Wort und noch eins\n- [ ] offener Punkt\nund `hier schliesst es.\n"

	vor := model.StructureRule{
		Files: "docs/*.md", SectionPattern: `^## DoD`, Sections: "each",
		ForbidPattern: `- \[ \]`,
	}
	if f := laufe(t, body, vor); f != nil {
		t.Fatalf("VORZUSTAND: die bereinigt lesende Form muss hier BLIND sein, got %+v", f)
	}
	f := laufe(t, body, rohRule())
	if len(f) != 1 || f[0].Reason != model.ReasonSectionTasksOpen {
		t.Fatalf("die rohe Bedingung muss trotz des Backticks melden, got %+v", f)
	}
	if f[0].Line != 6 {
		t.Fatalf("Befund gehoert auf die Zeile des Items (6), got %d", f[0].Line)
	}
}

// ALLE VIER LISTEN-MARKER, eingerueckt und mit Tab-Trenner — und die gehakte
// Box meldet nie. Ein Konfigurations-Muster deckt nur die Form, die sein Autor
// aufgeschrieben hat; die Modul-Lexik kann diesen Fehler nicht machen
// (BEO-003).
func TestMaxOpenTasks_AlleMarkerFormen(t *testing.T) {
	for name, tc := range map[string]struct {
		zeile string
		meldet bool
	}{
		"bindestrich":      {"- [ ] offen", true},
		"stern":            {"* [ ] offen", true},
		"plus":             {"+ [ ] offen", true},
		"geordnet":         {"1. [ ] offen", true},
		"eingerueckt":      {"    - [ ] offen", true},
		"tab-trenner":      {"-\t[ ] offen", true},
		"gehakt klein":     {"- [x] erledigt", false},
		"gehakt gross":     {"- [X] erledigt", false},
		"kein task-item":   {"- blosser Punkt", false},
		"blockquote":       {"> - [ ] zitiert", false},
		"tab in der box":   {"- [\t] exotisch", false},
	} {
		t.Run(name, func(t *testing.T) {
			f := laufe(t, "# D\n\n## DoD\n\n"+tc.zeile+"\n", rohRule())
			if tc.meldet && (len(f) != 1 || f[0].Reason != model.ReasonSectionTasksOpen) {
				t.Fatalf("erwartet genau ein section-tasks-open, got %+v", f)
			}
			if !tc.meldet && f != nil {
				t.Fatalf("erwartet befundfrei, got %+v", f)
			}
		})
	}
}

// FENCE-TREUE: ein Task-Item INNERHALB eines Fenced-Blocks zaehlt nicht —
// sonst meldete ein Dokument, das UEBER Task-Items schreibt, seine eigene
// Illustration.
func TestMaxOpenTasks_FenceBleibtAussen(t *testing.T) {
	body := "# D\n\n## DoD\n\n```md\n- [ ] nur Beispiel\n```\n"
	if f := laufe(t, body, rohRule()); f != nil {
		t.Fatalf("ein Item im Fence darf nicht melden, got %+v", f)
	}
	// Gegenprobe an derselben Datei: ausserhalb des Fence meldet dasselbe Item.
	if f := laufe(t, body+"\n- [ ] echt\n", rohRule()); len(f) != 1 {
		t.Fatalf("dasselbe Item ausserhalb des Fence muss melden, got %+v", f)
	}
}

// EIN BEFUND JE ITEM, auf SEINER Zeile — nicht einer je Datei und nicht einer
// auf der Abschnitts-Ueberschrift. Die Reparatur ist dort, wo der Haken steht.
func TestMaxOpenTasks_EinBefundJeItem(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n- [x] zwei\n- [ ] drei\n"
	f := laufe(t, body, rohRule())
	if len(f) != 2 {
		t.Fatalf("zwei offene Items ⇒ zwei Befunde, got %+v", f)
	}
	if f[0].Line != 5 || f[1].Line != 7 {
		t.Fatalf("Befunde gehoeren auf Zeile 5 und 7, got %d und %d", f[0].Line, f[1].Line)
	}
}

// DIE SCHWELLE ERLAUBT DIE ERSTEN N in Dokument-Reihenfolge und meldet nur den
// Ueberhang. Meldete sie alle, waere eine Verletzung drei Befunde, und keiner
// davon die Reparaturstelle.
func TestMaxOpenTasks_SchwelleMeldetNurDenUeberhang(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n- [ ] zwei\n- [ ] drei\n- [ ] vier\n"
	r := rohRule()
	r.MaxOpenTasks = ptr(2)
	f := laufe(t, body, r)
	if len(f) != 2 {
		t.Fatalf("vier offene bei Grenze 2 ⇒ zwei Befunde, got %+v", f)
	}
	if f[0].Line != 7 || f[1].Line != 8 {
		t.Fatalf("gemeldet gehoeren das dritte und vierte Item (Zeile 7, 8), got %d und %d",
			f[0].Line, f[1].Line)
	}
	r.MaxOpenTasks = ptr(4)
	if f := laufe(t, body, r); f != nil {
		t.Fatalf("vier offene bei Grenze 4 ⇒ befundfrei, got %+v", f)
	}
}

// DIE ABSCHNITTSGRENZE HAELT: ein offenes Item im Nachbar-Abschnitt zaehlt
// nicht mit. Ohne diese Zusage zaehlte die Bedingung ueber die ganze Datei,
// und die Mutation dazu liefe gruen durch.
func TestMaxOpenTasks_NurImAbschnitt(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [x] erledigt\n\n## Anderes\n\n- [ ] gehoert nicht dazu\n"
	if f := laufe(t, body, rohRule()); f != nil {
		t.Fatalf("das Item hinter der Abschnittsgrenze darf nicht melden, got %+v", f)
	}
	// Gegenprobe: dasselbe Item INNERHALB des Abschnitts meldet.
	drin := "# D\n\n## DoD\n\n- [x] erledigt\n- [ ] gehoert dazu\n\n## Anderes\n"
	if f := laufe(t, drin, rohRule()); len(f) != 1 || f[0].Line != 6 {
		t.Fatalf("dasselbe Item innerhalb muss auf Zeile 6 melden, got %+v", f)
	}
	// DIE GRENZE IST BEIDSEITIG SCHARF: das Item auf der LETZTEN Zeile des
	// Abschnitts gehoert noch dazu. Ohne diese Zusage liesse sich der Bereich
	// um eine Zeile verkuerzen, ohne dass ein Test rot wird.
	rand := "# D\n\n## DoD\n\n- [ ] letzte Zeile des Abschnitts\n## Anderes\n\nText.\n"
	if f := laufe(t, rand, rohRule()); len(f) != 1 || f[0].Line != 5 {
		t.Fatalf("das Item direkt vor der naechsten Ueberschrift zaehlt mit, got %+v", f)
	}
	// Und am DATEIENDE, wo keine naechste Ueberschrift folgt.
	ende := "# D\n\n## DoD\n\n- [ ] letzte Zeile der Datei\n"
	if f := laufe(t, ende, rohRule()); len(f) != 1 || f[0].Line != 5 {
		t.Fatalf("das Item auf der letzten Dateizeile zaehlt mit, got %+v", f)
	}
}

// EIN ABWESENDER SCHLUESSEL IST DIE BEDINGUNG AUS — und die explizite Null ist
// davon unterscheidbar. Sonst waere die Null-Schwelle, also der eigentliche
// Anwendungsfall, unerreichbar.
func TestMaxOpenTasks_AbwesendIstAus(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] offen\n"
	r := rohRule()
	r.MaxOpenTasks = nil
	if f := laufe(t, body, r); f != nil {
		t.Fatalf("ohne den Schluessel ist die Bedingung aus, got %+v", f)
	}
	if f := laufe(t, body, rohRule()); len(f) != 1 {
		t.Fatalf("mit expliziter Null meldet dieselbe Eingabe, got %+v", f)
	}
}

// DIE INLINE-CODE-GRENZE, wie sie IST und nicht wie sie klingt: eine
// EINZEILIGE Spanne meldet nicht, weil das Muster zeilen-verankert ist und der
// Backtick vor dem Listen-Marker steht. Nur die MEHRZEILIGE zaehlt mit — das
// ist der ausgewiesene Preis der rohen Lesung.
func TestMaxOpenTasks_InlineCodeGrenze(t *testing.T) {
	if f := laufe(t, "# D\n\n## DoD\n\nSo sieht es aus: `- [ ] offen`\n", rohRule()); f != nil {
		t.Fatalf("einzeilige Inline-Spanne meldet nicht, got %+v", f)
	}
	mehrzeilig := "# D\n\n## DoD\n\nSo `sieht es aus:\n- [ ] offen\nund` fertig.\n"
	if f := laufe(t, mehrzeilig, rohRule()); len(f) != 1 || f[0].Line != 6 {
		t.Fatalf("mehrzeilige Spanne zaehlt mit — das ist der Preis, got %+v", f)
	}
}

// DIE LEXIK IST GETEILT, NICHT KOPIERT — und diese Probe faengt die Doppelung,
// statt sie nur zu behaupten: sie ERWEITERT taskItemRE um eine Form, die es
// heute nicht gibt, und verlangt, dass offenerHaken ihr folgt. Ein zweites RE2
// neben taskItemRE — ein woertliches Praefix, wie es der Vorgaenger-Bau hatte —
// folgt ihr NICHT und macht diesen Test rot (BEO-003).
//
// Die naheliegende Form ("offenerHaken trifft nur, was taskItemRE trifft") kann
// das nicht: offenerHaken LIEST den Treffer von taskItemRE, die Bedingung ist
// per Konstruktion unerfuellbar, und der Test waere gruen ohne zu messen
// (BEO-023).
func TestOffenerHaken_FolgtDerErweitertenLexik(t *testing.T) {
	orig := taskItemRE
	t.Cleanup(func() { taskItemRE = orig })

	if offenerHaken("1) [ ] runde Klammer") {
		t.Fatalf("VORZUSTAND: die runde Klammer gehoert heute NICHT zur Lexik")
	}
	taskItemRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+[.)])[ \t]+\[[ xX]\]`)
	if !offenerHaken("1) [ ] runde Klammer") {
		t.Fatalf("offenerHaken folgt der erweiterten Lexik nicht — es gibt ein zweites Muster")
	}
	if offenerHaken("1) [x] runde Klammer") {
		t.Fatalf("die Verengung auf die LEERE Box haelt nicht mit")
	}
}
