package main

import (
	"strings"
	"testing"
)

func hasSubheading(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "## ") {
			return true
		}
	}
	return false
}

func TestExtractTitle(t *testing.T) {
	cases := []struct{ content, want string }{
		{"# Slice slice-190: Das Werkzeug bauen\n\nRest", "Das Werkzeug bauen"},
		{"# Welle welle-87: Wellen-Archivierung nachruesten\n", "Wellen-Archivierung nachruesten"},
	}
	for _, c := range cases {
		if got := ExtractTitle(c.content); got != c.want {
			t.Errorf("ExtractTitle(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestExtractTitle_Keine(t *testing.T) {
	if got := ExtractTitle("kein Heading hier"); got != "" {
		t.Fatalf("erwartet leeren String, got %q", got)
	}
}

// TestSliceStub prueft die Kanon-Form: keine Abschnittsueberschriften (das
// ist die Kuerzung, form-pruefbar), der ARCHIVIERT-Marker mit dem korrekten
// unzip-Zeiger, und dass Welle-Feld und Archiviert-mit-Feld getrennt bleiben
// (zwei Tatsachen, kein Widerspruch -- Template-Vorgabe).
func TestSliceStub(t *testing.T) {
	got := SliceStub("slice-190", "Das Werkzeug bauen", "— wellenlos", "welle-87", "DC-FA-PLAN-001")
	if hasSubheading(got) {
		t.Fatal("Stub traegt eine Abschnittsueberschrift -- das ist keine Kuerzung mehr")
	}
	if !strings.Contains(got, "unzip -p done/welle-87/archiv.zip") {
		t.Fatal("ARCHIVIERT-Marker fehlt oder zeigt auf den falschen Pfad")
	}
	if !strings.Contains(got, "**Welle:** — wellenlos") {
		t.Fatal("urspruengliches Welle-Feld nicht unveraendert uebernommen")
	}
	if !strings.Contains(got, "**Archiviert mit:** welle-87") {
		t.Fatal("Archiviert-mit-Feld fehlt")
	}
	if !strings.Contains(got, "**Hervorgegangen:** DC-FA-PLAN-001") {
		t.Fatal("Hervorgegangen-Feld nicht uebernommen")
	}
}

// TestExtractSurvivingIDs belegt die an der Archivierung von welle-60/61/63
// gemessene Luecke: eine Anforderung, deren einziger zitierender Slice
// archiviert wird, wurde zur Trace-Waise, weil der Stub die Kennung nicht
// mehr trug. Erfasst DC-*-Requirement-IDs und ADR-NNNN, dedupliziert,
// sortiert.
func TestExtractSurvivingIDs(t *testing.T) {
	content := "Betrifft DC-FA-REF-001 und noch einmal DC-FA-REF-001, dazu ADR-0044 und ADR-0020."
	got := ExtractSurvivingIDs(content)
	want := []string{"ADR-0020", "ADR-0044", "DC-FA-REF-001"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestFormatHervorgegangen_Keine(t *testing.T) {
	if got := FormatHervorgegangen(nil); got != "— keine —" {
		t.Fatalf("erwartet Template-Leerwert, got %q", got)
	}
}

// TestExtractSurvivingIDs_IgnoriertInlineCodeUndFences belegt die an
// slice-075 gemessene Luecke: eine erfundene, aehnlich geformte Kennung
// als Illustration eines Parsing-Grenzfalls in Inline-Code
// ("`GG-QA-001, 007 Sekunden`") darf nicht als echter Beleg in den Stub
// wandern -- sonst behauptet er eine Anforderung, die es nie gab. Eine
// echte Kennung ausserhalb von Code bleibt dabei erfasst.
func TestExtractSurvivingIDs_IgnoriertInlineCodeUndFences(t *testing.T) {
	content := "Betrifft DC-FA-REF-001 (echt).\n" +
		"Beispiel: `GG-QA-001, 007 Sekunden` ist kein echter Beleg.\n" +
		"```\nADR-9999 steht nur im Fenced-Block\n```\n"
	got := ExtractSurvivingIDs(content)
	want := []string{"DC-FA-REF-001"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// TestExtractSurvivingIDs_ErhaeltLinkLabel belegt die Gegenprobe zur
// vorigen: die weitaus haeufigere, echte Zitat-Form dieses Repos ist ein
// Markdown-Link mit Kennung als Label ("[`DC-FA-XXX-001`](ziel)") -- die
// erste Fassung des Inline-Code-Ausschlusses hatte das mitgeloescht
// (gemessen: 3 von 4 echten Kennungen in slice-073 verschwanden).
func TestExtractSurvivingIDs_ErhaeltLinkLabel(t *testing.T) {
	content := "Siehe [`DC-FA-PLAN-001`](spec/lastenheft.md#anchor) und " +
		"[ADR-0048](docs/plan/adr/0048-x.md)."
	got := ExtractSurvivingIDs(content)
	want := []string{"ADR-0048", "DC-FA-PLAN-001"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestWelleStub(t *testing.T) {
	got := WelleStub("welle-87", "Wellen-Archivierung", "welle-87-results.md", 2, 1)
	if hasSubheading(got) {
		t.Fatal("Stub traegt eine Abschnittsueberschrift")
	}
	if !strings.Contains(got, "2 Slices, 1 Reviews") {
		t.Fatalf("Zaehlung fehlt oder falsch: %q", got)
	}
	if !strings.Contains(got, "welle-87-results.md") {
		t.Fatal("Ergebnisnotiz-Zeiger fehlt")
	}
}
