package rules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// readErrFS ist ein Filesystem, dessen ReadFile immer scheitert — MemFS kann
// "vorhanden, aber unlesbar" nicht ausdruecken, und genau die Unterscheidung
// traegt den fail-closed-Vertrag.
type readErrFS struct{ driven.Filesystem }

func (readErrFS) ReadFile(string) ([]byte, error) { return nil, errors.New("unlesbar") }

func ptr(n int) *int { return &n }

// Happy Path: eine Regel, deren Abschnitt alle Bedingungen erfuellt.
func TestStructureHappyPath(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Ergebnis\n\n- **Beleg:** vorhanden.\nEins. Zwei.\n",
	}
	r := model.StructureRule{
		Files: "docs/*.md", Section: "## Ergebnis", NonEmpty: true,
		MinSentences: ptr(2), RequireAll: []string{"Beleg"},
	}
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r}); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}

// Ohne Regeln ist das Modul inert — keine Datei wird geoeffnet.
func TestStructureInertOhneRegeln(t *testing.T) {
	if f := CheckStructure(coretest.NewMemFS(map[string]string{"a.md": "x"}), nil); f != nil {
		t.Fatalf("ohne Regeln erwartet inert, got %+v", f)
	}
}

// Nullmengen-Haerte: eine Regel zu setzen IST die Behauptung, dass sie Dateien
// trifft — auch dann, wenn erst exempt-paths die Menge geleert hat.
func TestStructureNullmengeFailClosed(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n\n## E\n\nx.\n"}
	for name, r := range map[string]model.StructureRule{
		"Glob trifft nichts": {Files: "andere/*.md", Section: "## E"},
		"exempt leert":       {Files: "docs/*.md", Section: "## E", ExemptPaths: []string{"docs/**"}},
	} {
		t.Run(name, func(t *testing.T) {
			f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
			if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
				t.Fatalf("erwartet section-missing, got %+v", f)
			}
		})
	}
}

// Kardinalitaet: one bricht bei mehreren Treffern ab, each prueft jeden.
func TestStructureKardinalitaet(t *testing.T) {
	body := "# T\n\n## E\n\n\n\n## E\n\n\n"
	files := map[string]string{"docs/a.md": body}
	one := model.StructureRule{Files: "docs/*.md", Section: "## E", NonEmpty: true}
	f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{one})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionAmbiguous {
		t.Fatalf("sections one ⇒ ambiguous erwartet, got %+v", f)
	}
	if f[0].Line != 7 {
		t.Errorf("Zeile = %d, want 7 (der ZWEITE Treffer)", f[0].Line)
	}
	each := one
	each.Sections = "each"
	f = CheckStructure(coretest.NewMemFS(files), []model.StructureRule{each})
	if len(f) != 2 {
		t.Fatalf("sections each ⇒ jeder Treffer geprueft, got %+v", f)
	}
	for _, x := range f {
		if x.Reason != model.ReasonSectionEmpty {
			t.Errorf("erwartet section-empty, got %+v", x)
		}
	}
}

// Die sechs Bedingungen, jede mit EIGENEM Grund-Code — sonst fielen zwei
// Verletzungen desselben Abschnitts unter der Deduplikation zusammen.
func TestStructureSechsBedingungen(t *testing.T) {
	cases := []struct {
		name string
		body string
		rule model.StructureRule
		want string
	}{
		{"non-empty", "# T\n\n## E\n\n \n", model.StructureRule{NonEmpty: true}, model.ReasonSectionEmpty},
		{"min-sentences", "# T\n\n## E\n\nNur eins.\n", model.StructureRule{MinSentences: ptr(3)}, model.ReasonSectionThin},
		{"max-tasks", "# T\n\n## E\n\n- [ ] a\n- [x] b\n", model.StructureRule{MaxTasks: ptr(1)}, model.ReasonSectionOversized},
		{"forbid-pattern", "# T\n\n## E\n\nTODO offen.\n", model.StructureRule{ForbidPattern: `TODO`}, model.ReasonSectionForbidden},
		{"require-pattern", "# T\n\n## E\n\nnichts.\n", model.StructureRule{RequirePattern: `Beleg`}, model.ReasonSectionPatternMissing},
		{"require-all", "# T\n\n## E\n\nnichts.\n", model.StructureRule{RequireAll: []string{"Beleg"}}, model.ReasonSectionMarkerMissing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.rule
			r.Files, r.Section = "docs/*.md", "## E"
			f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": tc.body}), []model.StructureRule{r})
			if len(f) != 1 || f[0].Reason != tc.want {
				t.Fatalf("erwartet %s, got %+v", tc.want, f)
			}
		})
	}
}

// Zwei Verletzungen desselben Abschnitts muessen NEBENEINANDER stehen.
func TestStructureZweiBefundeAmselbenAbschnitt(t *testing.T) {
	r := model.StructureRule{
		Files: "docs/*.md", Section: "## E", MinSentences: ptr(5), RequirePattern: `Beleg`,
	}
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": "# T\n\n## E\n\nEins.\n"}), []model.StructureRule{r})
	if len(f) != 2 {
		t.Fatalf("erwartet zwei Befunde mit eigenen Grund-Codes, got %+v", f)
	}
}

// Die Marken-Formen: `- **M:**`, `**M:**` und `- **M (Zusatz):**` treffen,
// `M` im Fliesstext nicht, und `**Marker**` nicht auf `M`.
func TestStructureMarkenFormen(t *testing.T) {
	cases := map[string]bool{
		"- **Beleg:** da":        true,
		"**Beleg:** da":          true,
		"- **Beleg (Lauf):** da": true,
		"1. **Beleg:** da":       true,
		"Der Beleg steht da":     false,
		"- **Belegschaft:** da":  false,
		"Text **Beleg:** mitten": false,
	}
	for zeile, want := range cases {
		t.Run(zeile, func(t *testing.T) {
			body := "# T\n\n## E\n\n" + zeile + "\n"
			r := model.StructureRule{Files: "docs/*.md", Section: "## E", RequireAll: []string{"Beleg"}}
			f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
			if (len(f) == 0) != want {
				t.Fatalf("Marke erkannt=%v erwartet, got %+v", want, f)
			}
		})
	}
}

// Die Regel-Identitaet steht im Ziel: ohne sie verloere man je Datei den
// Befund der zweiten Regel unter der Deduplikation.
func TestStructureIdentitaetImZiel(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n\n## E\n\n \n"}
	r1 := model.StructureRule{Files: "docs/*.md", Section: "## E", NonEmpty: true}
	r2 := model.StructureRule{Files: "docs/**", Section: "## E", NonEmpty: true}
	f := model.SortFindings(CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r1, r2}))
	if len(f) != 2 {
		t.Fatalf("zwei Regeln ⇒ zwei Befunde, got %+v", f)
	}
	if f[0].Target == f[1].Target {
		t.Errorf("die Regel-Identitaeten muessen sich unterscheiden: %q", f[0].Target)
	}
}

// section-pattern statt section, und der Vergleich schliesst die #-Folge ein.
func TestStructureSelektorFormen(t *testing.T) {
	body := "# T\n\n### E\n\n \n"
	klartext := model.StructureRule{Files: "docs/*.md", Section: "## E", NonEmpty: true}
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{klartext})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("die #-Folge gehoert zum Vergleich ⇒ missing erwartet, got %+v", f)
	}
	muster := model.StructureRule{Files: "docs/*.md", SectionPattern: `^#{2,3} E$`, NonEmpty: true}
	f = CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{muster})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionEmpty {
		t.Fatalf("section-pattern trifft ⇒ empty erwartet, got %+v", f)
	}
}

// PRESET-KOPPLUNG: dieselbe Eingabe, beide Oberflaechen, gleiche
// Abschnitts-Grenze. Zwei Kopien der Mechanik koennten driften, ohne dass ein
// Test es merkt — dieser Test ist der Beleg, dass es eine ist.
func TestStructurePresetKopplungMitClosure(t *testing.T) {
	body := "# S\n\n## 7. Closure-Notiz\n\nEins `code` zwei.\n\n### Tiefer\n\nDrei.\n\n## 8. Danach\n\nVier. Fuenf. Sechs.\n"
	files := map[string]string{closureDir + "/slice-001-a.md": body}
	cfg := closureCfg()
	cfg.Closure.MinSentences = 9
	viaClosure := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(viaClosure) != 1 || viaClosure[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("Closure-Pfad: thin erwartet, got %+v", viaClosure)
	}
	r := model.StructureRule{
		Files: "**/slice-001-a.md", SectionPattern: `^#{2,3} .*Closure-Notiz`, MinSentences: ptr(9),
	}
	viaStructure := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(viaStructure) != 1 || viaStructure[0].Reason != model.ReasonSectionThin {
		t.Fatalf("structure-Pfad: thin erwartet, got %+v", viaStructure)
	}
	if viaClosure[0].Line != viaStructure[0].Line {
		t.Errorf("beide Oberflaechen muessen dieselbe Zeile melden: %d vs %d",
			viaClosure[0].Line, viaStructure[0].Line)
	}
	// Die Zaehl-Aussage muss identisch sein — sie steht in beiden Meldungen.
	cn := strings.SplitN(viaClosure[0].Message, "trägt ", 2)[1]
	sn := strings.SplitN(viaStructure[0].Message, "trägt ", 2)[1]
	if strings.Fields(cn)[0] != strings.Fields(sn)[0] {
		t.Errorf("beide Oberflaechen muessen dieselbe Satzzahl messen: %q vs %q", cn, sn)
	}
}

// Ein Task-Item braucht den LISTEN-MARKER. Eine Kastenklammer im Fliesstext
// ist keins — ohne diesen Fall bliebe der Rueckbau auf ein blosses `[ ]`
// unbemerkt.
func TestStructureTaskItemBrauchtListenMarker(t *testing.T) {
	body := "# T\n\n## E\n\nDer Zustand [x] wurde geprueft und [ ] blieb offen.\n\n- [ ] echt\n"
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", MaxTasks: ptr(1)}
	if f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r}); f != nil {
		t.Fatalf("nur EIN echtes Task-Item ⇒ befundfrei bei max-tasks 1, got %+v", f)
	}
	r.MaxTasks = ptr(0)
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionOversized {
		t.Fatalf("das eine echte Task-Item muss zaehlen, got %+v", f)
	}
}

// Die Schwelle an ihrem RAND: genau erreicht ist erfuellt, eins darunter nicht.
func TestStructureMinSentencesAmRand(t *testing.T) {
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", MinSentences: ptr(3)}
	for body, want := range map[string]int{
		"# T\n\n## E\n\nEins. Zwei. Drei.\n": 0,
		"# T\n\n## E\n\nEins. Zwei.\n":       1,
	} {
		f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
		if len(f) != want {
			t.Errorf("%q: erwartet %d Befund(e), got %+v", body, want, f)
		}
	}
}

// Die Marken-Grenze ist UNICODE-weit, nicht ASCII: eine Ziffer und ein Umlaut
// setzen die Marke gleichermassen fort.
func TestStructureMarkenGrenzeUnicode(t *testing.T) {
	for zeile, erfuellt := range map[string]bool{
		"- **Beleg:** da":            true,
		"- **Beleg2:** da":           false,
		"- **Belegüberblick:** da":   false,
		"- **Beleg (Lauf):** da":     true,
		"- **Beleg-Nachweis:** da":   true,
	} {
		body := "# T\n\n## E\n\n" + zeile + "\n"
		r := model.StructureRule{Files: "docs/*.md", Section: "## E", RequireAll: []string{"Beleg"}}
		f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
		if (len(f) == 0) != erfuellt {
			t.Errorf("%q: Marke erfuellt=%v erwartet, got %+v", zeile, erfuellt, f)
		}
	}
}

// Eine Ueberschrift mit nachlaufendem Whitespace muss ein auf `$` verankertes
// section-pattern weiter treffen — der Vergleich laeuft auf der GETRIMMTEN Zeile.
func TestStructureUeberschriftGetrimmt(t *testing.T) {
	body := "# T\n\n## E   \n\n \n"
	r := model.StructureRule{Files: "docs/*.md", SectionPattern: `^## E$`, NonEmpty: true}
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionEmpty {
		t.Fatalf("nachlaufender Whitespace darf den Selektor nicht kippen, got %+v", f)
	}
}

// Eine vorhandene, aber unlesbare Kandidaten-Datei ist FAIL-CLOSED: sie still
// zu ueberspringen waere genau der Gruen-Pfad, den die Anforderung ausschliesst.
func TestStructureUnlesbareDateiFailClosed(t *testing.T) {
	fs := readErrFS{coretest.NewMemFS(map[string]string{"docs/a.md": "# T\n\n## E\n\nx.\n"})}
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", NonEmpty: true}
	f := CheckStructure(fs, []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("unlesbare Datei muss fail-closed melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "unlesbar") {
		t.Errorf("die Meldung muss die Ursache nennen: %q", f[0].Message)
	}
}

// Ist schon der Dateibaum nicht lesbar, meldet jede Regel — ein leerer
// Befundsatz waere von "alle Regeln erfuellt" nicht zu unterscheiden.
func TestStructureUnlesbarerBaumFailClosed(t *testing.T) {
	fs := listErrFS{coretest.NewMemFS(map[string]string{"docs/a.md": "# T\n"})}
	rules := []model.StructureRule{
		{Files: "docs/*.md", Section: "## E", NonEmpty: true},
		{Files: "spec/*.md", Section: "## F", NonEmpty: true},
	}
	f := CheckStructure(fs, rules)
	if len(f) != len(rules) {
		t.Fatalf("erwartet je Regel einen Befund, got %+v", f)
	}
	for _, fi := range f {
		if fi.Reason != model.ReasonSectionMissing || !strings.Contains(fi.Message, "fail-closed") {
			t.Errorf("erwartet fail-closed-Befund, got %+v", fi)
		}
	}
	if f[0].Target == f[1].Target {
		t.Errorf("die Regel-Identitaet muss erhalten bleiben: %q", f[0].Target)
	}
}

// Nummerierte Listen sind ebenfalls Task-Item-Traeger — ohne diesen Fall liesse
// sich der Ziffern-Zweig aus dem Muster entfernen, ohne dass ein Test rot wird.
func TestStructureTaskItemNummerierteListe(t *testing.T) {
	body := "# T\n\n## E\n\n1. [ ] eins\n2. [x] zwei\n"
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", MaxTasks: ptr(1)}
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionOversized {
		t.Fatalf("nummerierte Task-Items muessen zaehlen, got %+v", f)
	}
}

// Die Marke muss am ANFANG des ausgezeichneten Laufs stehen, nicht irgendwo
// darin: `**Zwischenbeleg:**` erfuellt `Beleg` NICHT.
func TestStructureMarkeNurAlsPraefix(t *testing.T) {
	body := "# T\n\n## E\n\n- **Zwischenbeleg:** da\n"
	r := model.StructureRule{Files: "docs/*.md", Section: "## E", RequireAll: []string{"Beleg"}}
	f := CheckStructure(coretest.NewMemFS(map[string]string{"docs/a.md": body}), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMarkerMissing {
		t.Fatalf("die Marke steht nicht am Anfang ⇒ fehlend erwartet, got %+v", f)
	}
}


// headingRule: eine Regel, die jede Unterabschnitts-Ueberschrift des
// Abschnitts auf eine Kennung verpflichtet.
func headingRule() model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", Section: "## Schema", HeadingPattern: `^ID-[0-9]{3} `,
	}
}

// DC-FA-STRUCT-001 Happy: alle Unterabschnitte genuegen dem Muster.
func TestStructureHeadingHappyPath(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### ID-001 Erster\n\nx.\n\n### ID-002 Zweiter\n\ny.\n",
	}
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{headingRule()}); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}

// DC-FA-STRUCT-001 Negative: ein Befund JE verletzender Ueberschrift, auf
// IHRER Zeile — nicht einer am Abschnittskopf.
func TestStructureHeadingBefundJeUeberschrift(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### Ohne Kennung\n\nx.\n\n### ID-002 Zweiter\n\ny.\n\n### Auch ohne\n\nz.\n",
	}
	got := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{headingRule()})
	var lines []int
	for _, f := range got {
		if f.Reason != model.ReasonSectionHeadingMismatch {
			t.Fatalf("unerwarteter Grund-Code: %+v", f)
		}
		lines = append(lines, f.Line)
	}
	if fmt.Sprint(lines) != fmt.Sprint([]int{5, 13}) {
		t.Fatalf("Befunde auf den Zeilen der Ueberschriften erwartet, got %+v", got)
	}
	if !strings.Contains(got[0].Message, "Ohne Kennung") {
		t.Fatalf("Meldung nennt die Ueberschrift nicht: %q", got[0].Message)
	}
}

// DC-FA-STRUCT-001: die Bedingung benutzt die Heading-Lexik des Moduls —
// eingerueckt und Tab-getrennt zaehlen, eine gleichlautende Zeile im
// Fenced-Block nicht.
func TestStructureHeadingLexikDesModuls(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n  ### Eingerueckt\n\nx.\n\n###\tTab-getrennt\n\ny.\n\n" +
			"```md\n### Im Fence\n```\n\n### ID-003 Gut\n\nz.\n",
	}
	got := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{headingRule()})
	var texts []string
	for _, f := range got {
		texts = append(texts, strconv.Itoa(f.Line))
	}
	if len(got) != 2 {
		t.Fatalf("eingerueckt + Tab melden, Fence nicht — got %+v", got)
	}
	if fmt.Sprint(texts) != fmt.Sprint([]string{"5", "9"}) {
		t.Fatalf("Zeilen der beiden echten Ueberschriften erwartet, got %+v", got)
	}
}

// DC-FA-STRUCT-001: Default-Ebene ist Abschnitts-Ebene + 1; heading-level
// waehlt eine andere.
func TestStructureHeadingEbene(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### Dritte ohne\n\nx.\n\n#### Vierte ohne\n\ny.\n",
	}
	r := headingRule()
	got := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(got) != 1 || got[0].Line != 5 {
		t.Fatalf("Default = Abschnitts-Ebene + 1 (nur die dritte Ebene), got %+v", got)
	}
	r.HeadingLevel = ptr(4)
	got = CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(got) != 1 || got[0].Line != 9 {
		t.Fatalf("heading-level 4 meldet nur die vierte Ebene, got %+v", got)
	}
}

// DC-FA-STRUCT-001 Grenze: ohne Ueberschrift der geprueften Ebene ist die
// Bedingung wirkungslos — benannt, nicht zugesagt.
func TestStructureHeadingOhneUeberschriftWirkungslos(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n\n## Schema\n\nnur Prosa.\n"}
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{headingRule()}); f != nil {
		t.Fatalf("ohne Ueberschrift kein Befund erwartet, got %+v", f)
	}
}

// DC-FA-STRUCT-001 Grenze: eine Ebene flacher als der Abschnitt kann in ihm
// nicht vorkommen — der Abschnitt endet dort.
func TestStructureHeadingEbeneFlacherAlsAbschnitt(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### ID-001 Gut\n\nx.\n\n## Ohne Kennung\n\ny.\n",
	}
	r := headingRule()
	r.HeadingLevel = ptr(2)
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r}); f != nil {
		t.Fatalf("flachere Ebene kann im Abschnitt nicht vorkommen, got %+v", f)
	}
}

// DC-FA-STRUCT-001 / DC-QA-02: ohne heading-pattern aendert sich nichts.
func TestStructureHeadingModulAus(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### Ohne Kennung\n\nx.\n",
	}
	r := headingRule()
	r.HeadingPattern = ""
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r}); f != nil {
		t.Fatalf("ohne heading-pattern kein Befund erwartet, got %+v", f)
	}
}

// DC-FA-STRUCT-001: unter sections: each wird JEDER Abschnitt geprueft.
func TestStructureHeadingSectionsEach(t *testing.T) {
	files := map[string]string{
		"docs/a.md": "# T\n\n## Schema\n\n### Ohne A\n\nx.\n\n## Schema\n\n### Ohne B\n\ny.\n",
	}
	r := headingRule()
	r.Sections = "each"
	got := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(got) != 2 {
		t.Fatalf("each prueft jeden Abschnitt, got %+v", got)
	}
}
