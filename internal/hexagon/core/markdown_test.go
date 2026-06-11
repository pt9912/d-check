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

func TestStripInlineCode_MehrfachBackticks(t *testing.T) {
	// Doppel-Backtick-Span mit einfachem Backtick im Inhalt —
	// ersetzt durch Leerzeichen gleicher Länge (positionserhaltend)
	got := stripInlineCode("a ``x ` y`` b")
	if got != "a "+strings.Repeat(" ", len("``x ` y``"))+" b" {
		t.Fatalf("got %q", got)
	}
	// unbalancierte Backticks bleiben erhalten
	if got := stripInlineCode("a ` b"); got != "a ` b" {
		t.Fatalf("unbalanciert: got %q", got)
	}
	// positionserhaltend: kein Verschmelzen angrenzenden Texts
	if got := stripInlineCode("AD`x`R-0042"); got != "AD   R-0042" {
		t.Fatalf("verschmolzen: %q", got)
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
