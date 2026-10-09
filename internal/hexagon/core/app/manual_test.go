package app

import (
	"strings"
	"testing"
)

const manualDoc = "# Handbuch\n\nEinleitung.\n\n" +
	"## Modul reviews\n\nText zu reviews.\n\n### reviews.match\n\nmatch-Text.\n\n" +
	"## Modul structure\n\n```\n## reviews im Code-Block\n```\n\nText zu structure.\n"

func manualDocs() []ManualDoc {
	return []ManualDoc{{Path: "docs/h.md", Content: manualDoc}, {Path: "spec/s.md", Content: "# Spec\n\n## Reviews\n\nSpec-Text.\n"}}
}

// Treffer in beiden Dokumenten, Groß-/Kleinschreibung egal; ein verschachtelter
// Treffer steht einmal, eine Überschrift im Code-Block ist kein Treffer.
func TestManualSections_TrefferUndVerschachtelung(t *testing.T) {
	out, found := ManualSections(manualDocs(), " REVIEWS ")
	want := "==> docs/h.md:5\n## Modul reviews\n\nText zu reviews.\n\n### reviews.match\n\nmatch-Text.\n" +
		"\n==> spec/s.md:3\n## Reviews\n\nSpec-Text.\n"
	if !found || out != want {
		t.Fatalf("got found=%v\n%q\nwant\n%q", found, out, want)
	}
}

// Ein Abschnitt endet nicht an einer Überschrift im Code-Block; der Titel
// eines Dokuments liefert es ganz.
func TestManualSections_CodeBlockUndGanzesDokument(t *testing.T) {
	out, found := ManualSections(manualDocs(), "structure")
	if !found || !strings.Contains(out, "## reviews im Code-Block") || !strings.HasSuffix(out, "Text zu structure.\n") {
		t.Fatalf("der Abschnitt reicht über den Code-Block bis zum Ende, got %q", out)
	}
	out, found = ManualSections(manualDocs(), "Handbuch")
	if !found || out != "==> docs/h.md:1\n"+strings.TrimRight(manualDoc, "\n")+"\n" {
		t.Fatalf("der Titel liefert das ganze Dokument, got %q", out)
	}
}

// Kein Treffer, wenn nur der Fließtext den Begriff nennt.
func TestManualSections_KeinTreffer(t *testing.T) {
	if out, found := ManualSections(manualDocs(), "Einleitung"); found || out != "" {
		t.Fatalf("Fließtext ist kein Treffer, got %v %q", found, out)
	}
}

// Der Titel ist die erste Überschrift erster Ebene außerhalb von Fenced-Code.
func TestManualTitle_UeberspringtCodeBlock(t *testing.T) {
	doc := "```\n# kein Titel\n```\n\n# Titel\n\n## Abschnitt\n"
	if got := ManualTitle(doc); got != "Titel" {
		t.Fatalf("got %q, want %q", got, "Titel")
	}
}

// Verglichen wird der Überschriften-Text ohne die führende #-Folge: der
// Begriff "#" trifft keine Überschrift.
func TestManualSections_RauteIstKeinText(t *testing.T) {
	if out, found := ManualSections(manualDocs(), "#"); found || out != "" {
		t.Fatalf("die #-Folge gehört nicht zum Text, got %v %q", found, out)
	}
}
