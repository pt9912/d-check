package rules

import (
	"reflect"
	"strings"
	"testing"
)

func TestPreprocessMarkdown_FencesUndInlineCode(t *testing.T) {
	content := "Zeile1 [a](b.md)\n```\n[ignoriert](fehlt.md)\n```\nZeile5 `[code](x.md)` Ende\n~~~\nauch ignoriert\n~~~\nZeile9"
	lines := PreprocessMarkdown([]byte(content))

	var nums []int
	for _, l := range lines {
		nums = append(nums, l.No)
	}
	want := []int{1, 5, 9}
	if !reflect.DeepEqual(nums, want) {
		t.Fatalf("Zeilennummern = %v, want %v", nums, want)
	}
	// Inline-Code-Span ist geleert — positionserhaltend durch
	// Leerzeichen gleicher Länge (DC-FA-LINK-001 Boundary)
	if got := lines[1].Text; got != "Zeile5 "+strings.Repeat(" ", len("`[code](x.md)`"))+" Ende" {
		t.Fatalf("Inline-Code nicht positionserhaltend geleert: %q", got)
	}
}

// stripLine: Einzelzeilen-Strip über die absatzweise API (Test-Helfer).
func stripLine(s string) string {
	return stripInlineCodeByLine([]proseLine{{no: 1, raw: s}})[1]
}

func TestStripInlineCode_MehrfachBackticks(t *testing.T) {
	// Doppel-Backtick-Span mit einfachem Backtick im Inhalt —
	// ersetzt durch Leerzeichen gleicher Länge (positionserhaltend)
	got := stripLine("a ``x ` y`` b")
	if got != "a "+strings.Repeat(" ", len("``x ` y``"))+" b" {
		t.Fatalf("got %q", got)
	}
	// unbalancierte Backticks bleiben erhalten
	if got := stripLine("a ` b"); got != "a ` b" {
		t.Fatalf("unbalanciert: got %q", got)
	}
	// positionserhaltend: kein Verschmelzen angrenzenden Texts
	if got := stripLine("AD`x`R-0042"); got != "AD   R-0042" {
		t.Fatalf("verschmolzen: %q", got)
	}
	// ungeschlossene Folge ist literal, die Suche läuft dahinter
	// weiter: der spätere Doppel-Backtick-Span wird erkannt
	if got := stripLine("a ` b ``x`` c"); got != "a ` b "+strings.Repeat(" ", len("``x``"))+" c" {
		t.Fatalf("Scan nach literaler Folge abgebrochen: %q", got)
	}
}

// Mehrzeilige Code-Spans (CommonMark): ein über den Zeilenumbruch
// gebrochener Span invertiert NICHT die Backtick-Parität der
// Folgezeile — der DC-QA-04-False-Positive aus dem u-boot-Gegentest
// (Slice slice-012): [`ID`](ziel) nach Span-Fortsetzung bleibt ein
// intakter Link.
func TestPreprocessMarkdown_MehrzeiligerSpan(t *testing.T) {
	content := "vor (`u-boot init && u-boot add\n" +
		"postgres && u-boot up`, `x.md` [`LH-AK-002`](l.md#a)) ohne\n" +
		"\n" +
		"neuer Absatz ` bleibt literal"
	lines := PreprocessMarkdown([]byte(content))

	// Zeile 2: beide Einzel-Spans geleert, der Link intakt
	if !strings.Contains(lines[1].Text, "[") || !strings.Contains(lines[1].Text, "](l.md#a)") {
		t.Fatalf("Link nach Span-Fortsetzung zerstört: %q", lines[1].Text)
	}
	if strings.Contains(lines[1].Text, "x.md` ") {
		t.Fatalf("Einzel-Span nicht geleert: %q", lines[1].Text)
	}
	// Zeile 1: Span-Anteil ab Backtick geleert
	if strings.Contains(lines[0].Text, "u-boot init") {
		t.Fatalf("Span-Beginn nicht geleert: %q", lines[0].Text)
	}
	// Leerzeile beendet den Absatz: ungeschlossener Backtick im
	// neuen Absatz bleibt literal
	if lines[3].Text != "neuer Absatz ` bleibt literal" {
		t.Fatalf("Absatzgrenze verletzt: %q", lines[3].Text)
	}
}

// Fences unterbrechen Absätze: kein Span über eine Fence hinweg.
func TestPreprocessMarkdown_FenceUnterbrichtAbsatz(t *testing.T) {
	content := "a ` offen\n```\ncode\n```\nzu ` b"
	lines := PreprocessMarkdown([]byte(content))
	if lines[0].Text != "a ` offen" || lines[1].Text != "zu ` b" {
		t.Fatalf("Fence-Grenze verletzt: %q / %q", lines[0].Text, lines[1].Text)
	}
}

// slice-076 (DC-FA-LINK-001.a Schritt 1): eine ```-Zeile mit Backtick in der
// Infozeile ist KEIN Fence-Öffner (CommonMark), sondern Fließtext — sie schaltet
// den Fence-Zustand nicht um; ein echter ```-Öffner dahinter verdeckt weiterhin.
// Für ~~~ gilt die Regel nicht (CommonMark-Asymmetrie).
func TestPreprocessMarkdown_FenceInfozeileMitBacktickIstFliesstext(t *testing.T) {
	content := "```yaml-Fence (`x.md`) — Satz.\n[nach](ziel.md)\n```echt\nverdeckt\n```\nnach-echtem-Fence"
	lines := PreprocessMarkdown([]byte(content))
	var nums []int
	for _, l := range lines {
		nums = append(nums, l.No)
	}
	// Zeile 1 (Infozeile mit Backtick) + Zeile 2 (Prosa dahinter) bleiben sichtbar;
	// Zeile 3 ```echt öffnet einen echten Fence, Zeile 4 wird verdeckt, Zeile 5
	// schließt; Zeile 6 ist wieder sichtbar.
	want := []int{1, 2, 6}
	if !reflect.DeepEqual(nums, want) {
		t.Fatalf("Zeilennummern = %v, want %v", nums, want)
	}
	// ~~~ mit Backtick in der Infozeile bleibt ein Öffner (Regel gilt nur für ```).
	tilde := PreprocessMarkdown([]byte("~~~info (`x`)\nverdeckt\n~~~\nsichtbar"))
	if len(tilde) != 1 || tilde[0].No != 4 {
		t.Fatalf("~~~-Infozeile fälschlich als Fließtext behandelt: %+v", tilde)
	}
}

func TestExtractLinks(t *testing.T) {
	lines := []Line{
		{No: 3, Text: "[a](x.md) und ![bild](img.png) und [b](y.md \"Titel\")"},
		{No: 7, Text: "[text [nested]](z.md) und [c](<mit leer.md>)"},
	}
	refs := ExtractLinks(lines)
	want := []LinkRef{
		{Line: 3, Target: "x.md", Text: "a"},
		{Line: 3, Target: "img.png", Text: "bild", IsImage: true},
		{Line: 3, Target: "y.md", Text: "b"},
		{Line: 7, Target: "z.md", Text: "text [nested]"},
		{Line: 7, Target: "mit leer.md", Text: "c"},
	}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}

// Parser-Kanten: unbalancierte Klammern, fehlende Ziel-Klammer,
// unterminiertes <…>-Quoting — kein Link bzw. tolerantes Entquoten.
func TestExtractLinks_Kanten(t *testing.T) {
	lines := []Line{
		{No: 1, Text: "[offen](ohne-ende und [kein-ziel] sowie [x]y"},
		{No: 2, Text: "[q](<ohne-ende.md) und [ok](a.md)"},
	}
	refs := ExtractLinks(lines)
	want := []LinkRef{
		{Line: 2, Target: "ohne-ende.md", Text: "q"}, // <… ohne '>' → toleranter Prefix-Strip
		{Line: 2, Target: "a.md", Text: "ok"},
	}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}

// ADR-0092: ein unbalanciertes `[` in gewöhnlicher Prosa
// (keine schließende `]` auf seiner Zeile) darf NICHT mit einer späteren,
// unabhängigen `](…)`-Sequenz im selben Absatz zu einem erfundenen Link
// verschmelzen — die Linktext-Klammer bleibt strikt zeilenlokal, nur die
// Adress-Klammer darf um genau eine Zeile verlängert werden.
func TestExtractLinks_UnbalancierteKlammerVerschmilztNicht(t *testing.T) {
	lines := []Line{
		{No: 1, Text: "Das Array a[i steht in der Doku."},
		{No: 2, Text: "Ein Verweis](echt.md) wird erwähnt."},
	}
	refs := ExtractLinks(lines)
	if len(refs) != 0 {
		t.Fatalf("erfundener Link entstand aus unabhängigen Zeilen: %+v", refs)
	}
}

// ADR-0092: eine Zieladresse hinter `](` mit genau einem Zeilenumbruch wird
// erkannt (Grenzwert: der Lookahead reicht exakt eine Zeile weit), der Fund
// wird der öffnenden Zeile zugeschrieben.
func TestExtractLinks_ZeilenumbruchHinterKlammer(t *testing.T) {
	lines := []Line{
		{No: 1, Text: "[text]("},
		{No: 2, Text: "ziel.md) danach"},
	}
	refs := ExtractLinks(lines)
	want := []LinkRef{{Line: 1, Target: "ziel.md", Text: "text"}}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}

// ADR-0092: zwei vollständige, einzeilige Links auf unmittelbar
// aufeinanderfolgenden Zeilen desselben Absatzes bleiben unabhängig — der
// Ein-Zeilen-Lookahead greift nur, wenn die Adress-Klammer auf der
// aktuellen Zeile NICHT schließt.
func TestExtractLinks_ZweiVollstaendigeLinksImSelbenAbsatz(t *testing.T) {
	lines := []Line{
		{No: 1, Text: "[a](b.md)"},
		{No: 2, Text: "[c](d.md)"},
	}
	refs := ExtractLinks(lines)
	want := []LinkRef{
		{Line: 1, Target: "b.md", Text: "a"},
		{Line: 2, Target: "d.md", Text: "c"},
	}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}

// ADR-0093: eine Link-Referenz-Definition wird unabhängig von ihrer
// Verwendung erkannt, mit IsDefinition markiert.
func TestExtractLinks_ReferenzDefinition_Erkannt(t *testing.T) {
	lines := []Line{{No: 3, Text: "[label]: ziel.md"}}
	refs := ExtractLinks(lines)
	want := []LinkRef{{Line: 3, Target: "ziel.md", Text: "label", IsDefinition: true}}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}

// ADR-0093: ein Titel-Suffix auf derselben Zeile wird wie bei einem
// Inline-Link abgetrennt (NormalizeTarget).
func TestExtractLinks_ReferenzDefinition_TitelAbgetrennt(t *testing.T) {
	lines := []Line{{No: 1, Text: `[label]: ziel.md "Titel"`}}
	refs := ExtractLinks(lines)
	if len(refs) != 1 || refs[0].Target != "ziel.md" {
		t.Fatalf("Titel nicht korrekt abgetrennt: %+v", refs)
	}
}

// ADR-0093: eine Definition mit Blockquote-Präfix wird nicht erkannt —
// benannte Grenze, kein Fund entsteht.
func TestExtractLinks_ReferenzDefinition_BlockquotePraefixNichtErkannt(t *testing.T) {
	lines := []Line{{No: 1, Text: "> [label]: ziel.md"}}
	if refs := ExtractLinks(lines); len(refs) != 0 {
		t.Fatalf("Blockquote-Praefix haette die Definition nicht erkennen duerfen: %+v", refs)
	}
}

// ADR-0094: ein Backslash vor dem eigentlichen Label-Ende lässt die GANZE
// Zeile unerkannt — nicht nur eine falsche, aber unschädliche Label-Grenze.
func TestExtractLinks_ReferenzDefinition_BackslashLaesstZeileUnerkannt(t *testing.T) {
	lines := []Line{{No: 1, Text: `[foo\]bar]: ziel.md`}}
	if refs := ExtractLinks(lines); len(refs) != 0 {
		t.Fatalf("Backslash-Escape haette die ganze Zeile unerkannt lassen muessen: %+v", refs)
	}
}

// ADR-0095: eine gewoehnliche Prosazeile ohne delimitierten Titel ist KEINE
// Definition — Rot-Beleg gegen die ADR-0093/0094-Fassung, die hier faelschlich
// das erfundene Ziel "First" erkannt haette.
func TestExtractLinks_ReferenzDefinition_ProsaOhneTitelDelimiterNichtErkannt(t *testing.T) {
	lines := []Line{{No: 1, Text: "[TERM]: First In, First Out"}}
	if refs := ExtractLinks(lines); len(refs) != 0 {
		t.Fatalf("Prosazeile ohne Titel-Delimiter haette unerkannt bleiben muessen: %+v", refs)
	}
}

// ADR-0095: ein korrekt delimitierter Titel mit einfachen Anfuehrungszeichen
// bzw. Klammern wird weiterhin erkannt (keine Regression durch die
// Titel-Delimiter-Pflicht).
func TestExtractLinks_ReferenzDefinition_AlternativeTitelDelimiter(t *testing.T) {
	cases := []string{
		`[label]: ziel.md 'Titel'`,
		`[label]: ziel.md (Titel)`,
	}
	for _, text := range cases {
		refs := ExtractLinks([]Line{{No: 1, Text: text}})
		if len(refs) != 1 || refs[0].Target != "ziel.md" {
			t.Fatalf("%q: refs = %+v, want genau einen LinkRef mit Target ziel.md", text, refs)
		}
	}
}

// Die Invariante aus TrimFenceIndent — Space und Tab, nicht unicode-weit —
// gilt fuer JEDEN Konsumenten der Fence-Lexik, nicht nur fuer den Waechter aus
// dem Modul spans. Ohne Assertion je Konsument liesse sich einer davon
// zurueckdrehen, ohne dass ein Test rot wird; der Silent-Gruen-Pfad waere
// wieder offen, waehrend fence-unclosed schweigt.
func TestProseLinesUnicodeWhitespaceIstKeineFenceEinrueckung(t *testing.T) {
	content := []byte("a\n\n\u00a0```\n\nsichtbar\n")
	lines := PreprocessMarkdown(content)
	for _, ln := range lines {
		if ln.No == 5 && strings.Contains(ln.Text, "sichtbar") {
			return
		}
	}
	t.Fatalf("U+00A0 hat die Fence-Paritaet gekippt, Zeile 5 verschwand: %+v", lines)
}

func TestDiagramFenceLinesUnicodeWhitespaceIstKeineFenceEinrueckung(t *testing.T) {
	// Die U+00A0-Zeile darf den mermaid-Fence nicht vorzeitig oeffnen — sonst
	// gaelte sein Inhalt als Prosa und der echte Diagramm-Inhalt als ausserhalb.
	content := []byte("\u00a0```\n\n```mermaid\nA --> B\n```\n")
	got := diagramFenceLines(content, map[string]bool{"mermaid": true})
	if len(got) != 1 || got[0].raw != "A --> B" {
		t.Fatalf("Diagramm-Inhalt nicht erkannt: %+v", got)
	}
}

