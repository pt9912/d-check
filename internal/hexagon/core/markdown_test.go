package core

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

func TestExtractLinks(t *testing.T) {
	lines := []Line{
		{No: 3, Text: "[a](x.md) und ![bild](img.png) und [b](y.md \"Titel\")"},
		{No: 7, Text: "[text [nested]](z.md) und [c](<mit leer.md>)"},
	}
	refs := ExtractLinks(lines)
	want := []LinkRef{
		{Line: 3, Target: "x.md"},
		{Line: 3, Target: "img.png", IsImage: true},
		{Line: 3, Target: "y.md"},
		{Line: 7, Target: "z.md"},
		{Line: 7, Target: "mit leer.md"},
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
		{Line: 2, Target: "ohne-ende.md"}, // <… ohne '>' → toleranter Prefix-Strip
		{Line: 2, Target: "a.md"},
	}
	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("refs = %+v\nwant  %+v", refs, want)
	}
}
