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
// ist keiner (BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt); die Klasse selbst ist BEO-ALL/prosa-verschwindet-in-inline-code-spanne.
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
// (BEO-ALL/shared-lexicon-drifts-at-edges).
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
// folgt ihr NICHT und macht diesen Test rot (BEO-ALL/shared-lexicon-drifts-at-edges).
//
// Die naheliegende Form ("offenerHaken trifft nur, was taskItemRE trifft") kann
// das nicht: offenerHaken LIEST den Treffer von taskItemRE, die Bedingung ist
// per Konstruktion unerfuellbar, und der Test waere gruen ohne zu messen
// (BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt).
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

// gegenstandRule ist rohRule() plus die Marken-Kopplung (ADR-0085).
func gegenstandRule() model.StructureRule {
	r := rohRule()
	r.OpenTasksRequireMarker = "Gegenstand"
	return r
}

// ERLAUBNIS: eine vorhandene Marke ersetzt ALLE Einzelbefunde durch keinen —
// die Ablösung von CO-002s exempt-paths-Eintrag durch eine Inhalts-Pruefung.
func TestOpenTasksRequireMarker_VorhandeneMarkeErlaubtAlles(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n- [ ] zwei\n- [ ] drei\n\n**Gegenstand:** entfallen: Beispiel.\n"
	if f := laufe(t, body, gegenstandRule()); f != nil {
		t.Fatalf("vorhandene Marke muss alle Einzelbefunde tilgen, got %+v", f)
	}
}

// PFLICHT: fehlt die Marke, ersetzt EIN section-open-tasks-marker-missing die
// sonst drei section-tasks-open — nicht vier Befunde, einer.
func TestOpenTasksRequireMarker_FehlendeMarkeErgibtEinenBefund(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n- [ ] zwei\n- [ ] drei\n"
	f := laufe(t, body, gegenstandRule())
	if len(f) != 1 {
		t.Fatalf("erwartet GENAU einen Befund statt der drei Einzelbefunde, got %+v", f)
	}
	if f[0].Reason != model.ReasonSectionOpenTasksMarkerMissing {
		t.Fatalf("falscher Grund-Code, got %s", f[0].Reason)
	}
	if f[0].Line != 3 {
		t.Fatalf("Befund gehoert auf die Ueberschriftszeile (3), got %d", f[0].Line)
	}
}

// NORMALFALL UNBERUEHRT: ohne Ueberschuss-Items greift die Kopplung nicht --
// unabhaengig davon, ob die Marke dasteht.
func TestOpenTasksRequireMarker_NormalfallOhneOffeneItems(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [x] eins\n- [x] zwei\n"
	if f := laufe(t, body, gegenstandRule()); f != nil {
		t.Fatalf("keine offenen Items ⇒ befundfrei, unabhaengig von der Marke, got %+v", f)
	}
}

// ABWESENDER SCHLUESSEL: byte-identisches Verhalten zum Vorzustand -- die
// Kopplung greift nur, wenn OpenTasksRequireMarker gesetzt ist.
func TestOpenTasksRequireMarker_AbwesenderSchluesselByteIdentisch(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n- [ ] zwei\n- [ ] drei\n\n**Gegenstand:** entfallen: Beispiel.\n"
	f := laufe(t, body, rohRule())
	if len(f) != 3 {
		t.Fatalf("ohne den Schluessel muessen weiterhin alle drei melden, got %+v", f)
	}
	for _, x := range f {
		if x.Reason != model.ReasonSectionTasksOpen {
			t.Fatalf("ohne den Schluessel bleibt der alte Grund-Code, got %s", x.Reason)
		}
	}
}

// gegenstandSectionRule ist gegenstandRule() plus die Sektions-Verlegung
// (ADR-0085, Geschichte 2026-09-17): die Marke wird in "Closure-Notiz"
// gesucht, nicht im gezaehlten DoD-Abschnitt.
func gegenstandSectionRule() model.StructureRule {
	r := gegenstandRule()
	r.OpenTasksRequireMarkerSection = `^#{1,3} [0-9]+\. Closure-Notiz`
	return r
}

// R1-F-1 REPRODUKTION: die Marke steht NUR in "## 7. Closure-Notiz", nicht
// im gezaehlten DoD-Abschnitt -- exakt Baseline v6.9.0s eigene Ziel-Form.
// Ohne OpenTasksRequireMarkerSection findet die Kopplung sie NICHT (Vorzustand,
// vom Review R1-F-1 empirisch belegt); MIT gesetzter Sektion findet sie sie.
func TestOpenTasksRequireMarkerSection_BaselineZielFormWirdErkannt(t *testing.T) {
	body := "# D\n\n## 2. Definition of Done\n\n- [ ] eins\n- [ ] zwei\n\n" +
		"## 3. Plan (vor Code)\n\nText.\n\n## 7. Closure-Notiz\n\n" +
		"**Gegenstand:** entfallen: Beispiel.\n"
	rule := rohRule()
	rule.Section = ""
	rule.SectionPattern = `^## [0-9]+\. Definition of Done`
	rule.OpenTasksRequireMarker = "Gegenstand"

	// VORZUSTAND (R1-F-1): ohne die Sektions-Verlegung bleibt die Marke im
	// falschen Abschnitt unsichtbar -- die Pflicht-Haelfte feuert faelschlich.
	f := laufe(t, body, rule)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionOpenTasksMarkerMissing {
		t.Fatalf("VORZUSTAND: ohne Sektions-Verlegung muss die Marke im falschen Abschnitt unsichtbar bleiben, got %+v", f)
	}

	// FIX: mit OpenTasksRequireMarkerSection wird "Closure-Notiz" durchsucht.
	rule.OpenTasksRequireMarkerSection = `^#{1,3} [0-9]+\. Closure-Notiz`
	if f := laufe(t, body, rule); f != nil {
		t.Fatalf("mit Sektions-Verlegung muss die Marke in Closure-Notiz gefunden werden, got %+v", f)
	}
}

// FEHLENDER ABSCHNITT: existiert der benannte Abschnitt in der Datei gar
// nicht, gilt die Marke als fehlend -- dieselbe Lesart wie "Marke nicht da".
func TestOpenTasksRequireMarkerSection_FehlenderAbschnittGiltAlsFehlend(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n"
	if f := laufe(t, body, gegenstandSectionRule()); len(f) != 1 || f[0].Reason != model.ReasonSectionOpenTasksMarkerMissing {
		t.Fatalf("kein Closure-Notiz-Abschnitt ⇒ Marke fehlt, got %+v", f)
	}
}

// ANDERE NUMMERIERUNG: die Baseline-Form variiert die Ueberschriften-Nummer
// ("## 9. Closure-Notiz (nach `done/`)" ist realer Bestand dieses Repos) --
// das RE2-Praefix trifft trotzdem, weil es nicht auf das Zeilenende ankert.
func TestOpenTasksRequireMarkerSection_AndereNummerierungTrifftTrotzdem(t *testing.T) {
	body := "# D\n\n## DoD\n\n- [ ] eins\n\n## 9. Closure-Notiz (nach `done/`)\n\n" +
		"**Gegenstand:** entfallen: Beispiel.\n"
	if f := laufe(t, body, gegenstandSectionRule()); f != nil {
		t.Fatalf("Praefix-Muster muss abweichende Nummerierung treffen, got %+v", f)
	}
}
