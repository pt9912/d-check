package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// dodBody ist die Gestalt, um die der Antrag gestellt wurde: drei
// Liefer-Punkte und vier Punkte, die in JEDEM Slice stehen.
const dodBody = "# T\n\n## 4. Definition of Done\n\n" +
	"- [ ] **Lieferung A** steht.\n" +
	"- [ ] **Lieferung B** steht.\n" +
	"- [ ] **Lieferung C** steht.\n" +
	"- [ ] **Konstante:** Closure-Notiz geschrieben.\n" +
	"- [ ] **Konstante:** Beobachtungs-Register fortgeschrieben.\n" +
	"- [ ] **Konstante:** jedes Risiko aus dem Plan mit Ausgang.\n" +
	"- [ ] **Konstante:** Gate-Lauf gruen, Review, Verifikation.\n"

func dodRule() model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", Section: "## 4. Definition of Done", MaxTasks: ptr(3),
	}
}

func laufe(t *testing.T, body string, r model.StructureRule) []model.Finding {
	t.Helper()
	return CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
}

// Die erste Probe des Antrags — und die Umkehr dazu in derselben Funktion:
// OHNE das Muster meldet dieselbe Eingabe rot. Ein Regressions-Test ohne
// belegte Regression ist keiner (BEO-023).
func TestTasksIgnorePattern_KonstantenZaehlenNichtMit(t *testing.T) {
	r := dodRule()
	if f := laufe(t, dodBody, r); len(f) != 1 || f[0].Reason != model.ReasonSectionOversized {
		t.Fatalf("VORZUSTAND: sieben Items bei max-tasks 3 muessen melden, got %+v", f)
	}
	r.TasksIgnorePattern = `^\*\*Konstante:`
	if f := laufe(t, dodBody, r); f != nil {
		t.Fatalf("drei Liefer-Punkte bei max-tasks 3 ⇒ befundfrei, got %+v", f)
	}
}

// Die zweite Probe des Antrags: das Muster nimmt die Konstanten heraus und
// deckt den vierten Liefer-Punkt NICHT zu.
func TestTasksIgnorePattern_VierterLieferPunktMeldet(t *testing.T) {
	r := dodRule()
	r.TasksIgnorePattern = `^\*\*Konstante:`
	body := strings.Replace(dodBody, "- [ ] **Konstante:** Closure-Notiz geschrieben.\n",
		"- [ ] **Lieferung D** steht.\n- [ ] **Konstante:** Closure-Notiz geschrieben.\n", 1)
	f := laufe(t, body, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionOversized {
		t.Fatalf("vier Liefer-Punkte bei max-tasks 3 muessen melden, got %+v", f)
	}
}

// Die dritte Probe des Antrags: ohne den Schluessel ist die Meldung
// BYTE-IDENTISCH zu der vor dieser Faehigkeit. Der Vergleich steht als
// Literal da, weil genau er die Zusage ist — jede Ergaenzung an dieser
// Zeichenkette bricht ihn.
func TestTasksIgnorePattern_AbwesendIstByteIdentisch(t *testing.T) {
	f := laufe(t, dodBody, dodRule())
	const want = "Abschnitt trägt 7 Task-Items, erlaubt sind 3"
	if len(f) != 1 || f[0].Message != want {
		t.Fatalf("ohne Muster erwartet %q, got %+v", want, f)
	}
}

// Die Ueberdeckung ist SICHTBAR: die Meldung nennt die Zahl der ignorierten
// Items. Beide Raender sind hier drin — ein Muster, das zu viel nimmt, und
// eines, das nichts nimmt.
func TestTasksIgnorePattern_MeldungNenntIgnorierte(t *testing.T) {
	for name, tc := range map[string]struct{ pat, want string }{
		"zu breit": {`^\*\*`, "Abschnitt trägt 0 Task-Items (7 ignoriert), erlaubt sind -1"},
		"trifft nichts": {`^ZZZ`,
			"Abschnitt trägt 7 Task-Items (0 ignoriert), erlaubt sind -1"},
	} {
		t.Run(name, func(t *testing.T) {
			r := dodRule()
			// Schwelle -1: erzwingt die Meldung auch dort, wo das Muster ALLES
			// nimmt. Ohne sie bliebe genau der interessante Fall stumm — und
			// das ist die benannte Grenze der Sichtbarkeits-Zusage.
			r.MaxTasks = ptr(-1)
			r.TasksIgnorePattern = tc.pat
			f := laufe(t, dodBody, r)
			if len(f) != 1 || f[0].Message != tc.want {
				t.Fatalf("erwartet %q, got %+v", tc.want, f)
			}
		})
	}
}

// Das Muster sieht den ITEM-TEXT, nicht die rohe Zeile. Beide Richtungen, weil
// nur das Paar die Entscheidung belegt: verankert am Text trifft, verankert am
// Listen-Marker trifft nicht.
func TestTasksIgnorePattern_SiehtDenItemText(t *testing.T) {
	for name, tc := range map[string]struct {
		pat  string
		will int
	}{
		"am Item-Text verankert":   {`^\*\*Konstante:`, 0},
		"am Listen-Marker verankert": {`^- \[ \] \*\*Konstante:`, 1},
	} {
		t.Run(name, func(t *testing.T) {
			r := dodRule()
			r.TasksIgnorePattern = tc.pat
			if f := laufe(t, dodBody, r); len(f) != tc.will {
				t.Fatalf("erwartet %d Befund(e), got %+v", tc.will, f)
			}
		})
	}
}

// FENCE- UND INLINE-CODE-TREUE, und die zweite Haelfte ist die teure: das
// Muster sieht den BEREINIGTEN Text. Was in Backticks steht, ist beim Zaehlen
// durch Leerzeichen ersetzt — ein Muster auf `make gates` trifft ein Item
// NICHT, das die Wendung in Inline-Code schreibt. Genau so steht sie in den
// Slice-Plaenen, die den Antrag ausgeloest haben.
func TestTasksIgnorePattern_LiestDenBereinigtenText(t *testing.T) {
	body := "# T\n\n## E\n\n" +
		"- [ ] `make gates` gruen.\n" +
		"- [ ] noch etwas.\n\n" +
		"```\n- [ ] im Fence, zaehlt nie\n```\n"
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", MaxTasks: ptr(2)}
	if f := laufe(t, body, r); f != nil {
		t.Fatalf("das Item im Fence darf nicht zaehlen, got %+v", f)
	}
	r.MaxTasks = ptr(0)
	r.TasksIgnorePattern = `make gates`
	f := laufe(t, body, r)
	if len(f) != 1 || f[0].Message != "Abschnitt trägt 2 Task-Items (0 ignoriert), erlaubt sind 0" {
		t.Fatalf("Inline-Code ist beim Zaehlen geleert — das Muster darf ihn nicht sehen, got %+v", f)
	}
	// Gegenprobe: dieselbe Zusage ueber den Text, der die Bereinigung ueberlebt.
	r.TasksIgnorePattern = `^gruen\.$`
	if f := laufe(t, body, r); len(f) != 1 ||
		f[0].Message != "Abschnitt trägt 1 Task-Items (1 ignoriert), erlaubt sind 0" {
		t.Fatalf("ueber den ueberlebenden Text muss es treffen, got %+v", f)
	}
}

// --- exempt-section-pattern ---

// acBody traegt drei gleichartige Abschnitte, von denen die beiden ersten den
// Bestand vor einem Stichtag darstellen.
const acBody = "# Lastenheft\n\n" +
	"### AC-001 — alt\n\nHappy: ja.\n\n" +
	"### AC-002 — alt\n\nHappy: ja.\n\n" +
	"### AC-042 — neu\n\nHappy: ja.\n"

func acRule() model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", SectionPattern: `^### AC-`, Sections: "each",
		RequireAll: []string{"Boundary"},
	}
}

// Die vierte und fuenfte Probe des Antrags in einem Lauf — mit der Umkehr
// davor: OHNE den Schluessel melden alle drei Abschnitte (BEO-023).
func TestExemptSectionPattern_NimmtDenBestandHeraus(t *testing.T) {
	r := acRule()
	if f := laufe(t, acBody, r); len(f) != 3 {
		t.Fatalf("VORZUSTAND: alle drei Abschnitte muessen melden, got %+v", f)
	}
	r.ExemptSectionPattern = `^### AC-0(0[1-9]|1[0-9])\b`
	f := laufe(t, acBody, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMarkerMissing {
		t.Fatalf("nur der neue Abschnitt darf melden, got %+v", f)
	}
	if f[0].Line != 11 {
		t.Fatalf("der Befund muss auf AC-042 zeigen (Zeile 11), got Zeile %d", f[0].Line)
	}
}

// Das Muster sieht DIESELBE Zeile wie section-pattern — die rohe
// Ueberschriften-Zeile SAMT #-Folge. Beide Richtungen, weil nur das Paar die
// Entscheidung belegt; die Form ohne #-Folge ist die aus dem Antrag.
func TestExemptSectionPattern_SiehtDieRoheUeberschriftenZeile(t *testing.T) {
	for name, tc := range map[string]struct {
		pat  string
		will int
	}{
		"mit #-Folge, wie section-pattern": {`^### AC-00`, 1},
		"ohne #-Folge, wie im Antrag":      {`^AC-00`, 3},
	} {
		t.Run(name, func(t *testing.T) {
			r := acRule()
			r.ExemptSectionPattern = tc.pat
			if f := laufe(t, acBody, r); len(f) != tc.will {
				t.Fatalf("erwartet %d Befund(e), got %+v", tc.will, f)
			}
		})
	}
}

// Nullmengen-Haerte eine Granularitaetsstufe tiefer: leert das Ventil die
// Menge, meldet die Regel — und die Meldung nennt den Schluessel, der es tat.
func TestExemptSectionPattern_LeereMengeMeldet(t *testing.T) {
	r := acRule()
	r.ExemptSectionPattern = `^### AC-`
	f := laufe(t, acBody, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("alle Abschnitte ausgenommen ⇒ section-missing, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "exempt-section-pattern") ||
		!strings.Contains(f[0].Message, "alle 3 passenden") {
		t.Fatalf("die Meldung muss Schluessel und Zahl nennen, got %q", f[0].Message)
	}
}

// Die Ausnahme erklaert die GRUNDMENGE und laeuft deshalb VOR der
// Kardinalitaets-Pruefung: zwei Treffer, einer ausgenommen, sections: one ⇒
// kein section-ambiguous.
func TestExemptSectionPattern_LaeuftVorDerKardinalitaet(t *testing.T) {
	body := "# T\n\n### AC-001 — alt\n\n**Boundary:** ja.\n\n### AC-042 — neu\n\n**Boundary:** ja.\n"
	r := model.StructureRule{
		Files: "docs/*.md", SectionPattern: `^### AC-`, RequireAll: []string{"Boundary"},
	}
	if f := laufe(t, body, r); len(f) != 1 || f[0].Reason != model.ReasonSectionAmbiguous {
		t.Fatalf("VORZUSTAND: zwei Treffer bei sections one ⇒ ambiguous, got %+v", f)
	}
	r.ExemptSectionPattern = `^### AC-001\b`
	if f := laufe(t, body, r); f != nil {
		t.Fatalf("nach Abzug bleibt genau einer ⇒ befundfrei, got %+v", f)
	}
}

// Eine Ueberschrift im Fenced-Block ist keine — sie wird weder ausgewaehlt
// noch ausgenommen. Die Zusage ist geerbt, nicht neu; sie steht hier, weil
// eine geerbte Zusage ohne Probe eine Erinnerung ist.
//
// DAS MUSTER ZEIGT AUF DIE ECHTE UEBERSCHRIFT, nicht auf die gefencte, und
// das ist der ganze Trick: zeigte es auf die gefencte, waere "nie gewaehlt"
// von "gewaehlt und dann ausgenommen" am Ergebnis nicht unterscheidbar --
// ein fence-blinder Scanner liefe gruen durch. So bleibt die gefencte
// Ueberschrift die EINZIGE nicht ausgenommene, und Fence-Blindheit macht aus
// dem Nullmengen-Befund einen Marken-Befund.
func TestExemptSectionPattern_FenceTreu(t *testing.T) {
	body := "# T\n\n```\n### AC-001 — im Fence\n```\n\n### AC-042 — echt\n\nHappy: ja.\n"
	r := acRule()
	r.ExemptSectionPattern = `^### AC-042\b`
	f := laufe(t, body, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("die gefencte Ueberschrift ist keine ⇒ Nullmenge, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "alle 1 passenden") {
		t.Fatalf("genau EIN Abschnitt darf gezaehlt worden sein, got %q", f[0].Message)
	}
}

// Der verfasste Hinweis (ADR-0073) gewinnt gegen die modul-eigene Meldung —
// ABER NICHT bei der leer laufenden Regel: dort hat die Regel nicht gemessen,
// und die Diagnose, dass das Ventil die Menge geleert hat, ist die einzige
// Spur. Ohne diese Ausnahme loeschte ein hint genau den Text, der das zu
// breite Muster sichtbar macht.
func TestExemptSectionPattern_NullmengeBehaeltIhreMeldung(t *testing.T) {
	r := acRule()
	r.ExemptSectionPattern = `^### AC-`
	r.Hint = "Jede Anforderung braucht eine Boundary-Marke"
	f := laufe(t, acBody, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("erwartet section-missing, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "exempt-section-pattern") {
		t.Fatalf("der hint darf die Nullmengen-Diagnose nicht verdraengen, got %q", f[0].Message)
	}
}

// Die Kehrseite, und sie ist eine BENANNTE GRENZE, kein Defekt: bei einer
// verletzten BEDINGUNG gewinnt der Hinweis wie ueberall sonst — samt der Zahl
// der ignorierten Items. Wer die Ueberdeckung sehen will, laesst den Hinweis
// an dieser Regel weg. Der Test haelt die Grenze fest, damit sie nicht
// stillschweigend kippt.
func TestTasksIgnorePattern_HintVerdraengtDieZahl(t *testing.T) {
	r := dodRule()
	r.TasksIgnorePattern = `^ZZZ`
	r.Hint = "Hoechstens drei Liefer-Punkte"
	f := laufe(t, dodBody, r)
	if len(f) != 1 || f[0].Message != "Hoechstens drei Liefer-Punkte" {
		t.Fatalf("bei verletzter Bedingung gewinnt der Hinweis, got %+v", f)
	}
}

// Zwei Regeln ueber denselben Selektor, die sich nur durch die Ausnahme
// unterscheiden, sind VERSCHIEDENE Zusagen — und genau diese Paarung braucht
// Grandfathering. Ohne die Ausnahme in der Identitaet waere sie ein
// Konfigurations-Duplikat und damit unschreibbar (ADR-0075).
func TestExemptSectionPattern_TraegtDieRegelIdentitaet(t *testing.T) {
	a := acRule()
	b := acRule()
	b.ExemptSectionPattern = `^### AC-001\b`
	if a.Identity() == b.Identity() {
		t.Fatalf("zwei Regeln mit verschiedener Ausnahme brauchen verschiedene Identitaeten, beide %q", a.Identity())
	}
	// Ein LEERES Muster darf die Identitaet NICHT veraendern: sonst aenderte
	// sich das target jeder Bestandsregel.
	if a.Identity() != "docs/*.md :: ^### AC-" {
		t.Fatalf("die Identitaet einer Regel ohne Ausnahme muss unveraendert bleiben, got %q", a.Identity())
	}
}
