package core

import (
	"regexp"
	"strings"
	"testing"
)

// Vollständigkeits-Sicherung (slice-025 DoD): jeder Grund-Code aus
// AllReasons hat einen Klartext, und das Mapping trägt keine verwaisten
// Einträge. Ein neuer Reason*-Code ohne Klartext bricht hier — nicht
// still in der Diagnose.
func TestReasonTextDeckungGegenAllReasons(t *testing.T) {
	reasons := AllReasons()
	texts := reasonTexts()
	for _, r := range reasons {
		if _, ok := texts[r]; !ok {
			t.Errorf("Grund-Code %q ohne Klartext", r)
		}
		if ReasonText(r) == r {
			t.Errorf("ReasonText(%q) liefert den nackten Code zurück (Klartext fehlt)", r)
		}
	}
	if len(texts) != len(reasons) {
		t.Errorf("reasonTexts hat %d Einträge, AllReasons %d — verwaister oder fehlender Code",
			len(texts), len(reasons))
	}
}

// ReasonText: unbekannter Code → fail-safe der nackte Code.
func TestReasonTextUnbekannt(t *testing.T) {
	if got := ReasonText("gibt-es-nicht"); got != "gibt-es-nicht" {
		t.Fatalf("ReasonText(unbekannt) = %q, want den Code selbst", got)
	}
}

// FixCandidateFor: id-unlinked liefert einen Link-Kandidaten auf das
// Definitions-Target (relativ zum Verzeichnis der Befund-Datei).
func TestFixCandidateFor_IDUnlinked(t *testing.T) {
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"},
	}}
	f := Finding{File: "docs/a.md", Line: 1, Rule: "ids", Target: "ADR-0042", Reason: ReasonIDUnlinked}
	c := FixCandidateFor(f, cfg)
	if c == nil {
		t.Fatal("kein Kandidat für id-unlinked")
	}
	if c.Original != "ADR-0042" {
		t.Fatalf("Original = %q", c.Original)
	}
	if c.Replacement != "[`ADR-0042`](plan/adr)" {
		t.Fatalf("Replacement = %q, want [`ADR-0042`](plan/adr)", c.Replacement)
	}
}

// FixCandidateFor: andere Grund-Codes liefern (noch) keinen Kandidaten.
func TestFixCandidateFor_AndereCodesNil(t *testing.T) {
	cfg := Config{}
	for _, r := range []string{ReasonTargetMissing, ReasonSpanUnclosed, ReasonAnchorMissing} {
		f := Finding{File: "docs/a.md", Line: 1, Target: "x", Reason: r}
		if c := FixCandidateFor(f, cfg); c != nil {
			t.Errorf("Grund-Code %q sollte keinen Kandidaten liefern, got %+v", r, c)
		}
	}
}

// FixCandidateFor: id-unlinked ohne passendes Muster → kein Kandidat
// (kein Raten ohne bekanntes Definitions-Target).
func TestFixCandidateFor_OhneMusterNil(t *testing.T) {
	cfg := Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`MR-\d{3}`), Target: "harness/conventions.md"},
	}}
	f := Finding{File: "docs/a.md", Line: 1, Target: "ADR-0042", Reason: ReasonIDUnlinked}
	if c := FixCandidateFor(f, cfg); c != nil {
		t.Fatalf("ohne passendes Muster kein Kandidat erwartet, got %+v", c)
	}
}

// relLink: Datei in der Repo-Wurzel → Target unverändert relativ.
func TestRelLink_Wurzel(t *testing.T) {
	if got := relLink("README.md", "spec/lastenheft.md"); !strings.HasPrefix(got, "spec/") {
		t.Fatalf("relLink(Wurzel) = %q", got)
	}
}
