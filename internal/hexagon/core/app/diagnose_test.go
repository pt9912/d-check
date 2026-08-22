package app

import (
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"os"
	"path/filepath"
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

// Verriegelung AllReasons ↔ Spezifikation §4 (slice-060): jeder in der
// §4-Grund-Code-Tabelle dokumentierte Code steht in AllReasons, und
// AllReasons führt keinen Code ohne §4-Zeile. Die hand-gepflegte Liste
// kann der doc-first gepflegten Spec damit nicht mehr still
// hinterherhinken (von v0.25 bis v0.37 fehlten so sieben Codes, die der
// Deckungs-Test oben nicht sehen konnte — er prüft nur das Paar
// reasonTexts ↔ AllReasons, nicht die Liste selbst).
func TestAllReasonsDeckungGegenSpezifikationGrundCodes(t *testing.T) {
	specCodes := grundCodesAusSpezifikation(t)
	inAll := map[string]bool{}
	for _, r := range AllReasons() {
		inAll[r] = true
	}
	for code := range specCodes {
		if !inAll[code] {
			t.Errorf("§4-Grund-Code %q fehlt in AllReasons()", code)
		}
	}
	for r := range inAll {
		if !specCodes[r] {
			t.Errorf("AllReasons()-Code %q hat keine Zeile in Spezifikation §4", r)
		}
	}
}

// grundCodesAusSpezifikation extrahiert die Code-Spalte der Grund-Code-
// Tabelle aus spec/spezifikation.md §4 (Pfad relativ zum Paket; der
// Docker-Build-Kontext trägt die gesamte Repo-Wurzel). Fail-closed:
// unlesbare Spec, fehlende oder mehrdeutige Überschrift, eine leere
// Tabelle und jede Tabellen-Body-Zeile ohne Backtick-Code sind Fehler —
// kein stilles Grün mit leerer Menge, kein stilles Auslassen einzelner
// Zeilen. Annahme: alle Tabellenzeilen unter §4 sind Grund-Codes; eine
// künftige zweite Tabelle (z. B. Fehler-Codes) macht den Test laut rot
// und verlangt eine bewusste Parser-Anpassung. Die erste Spalte trägt
// seit der Struktur-ID-Vergabe die SPEC-Kennung (MR-000): sie wird
// mitgelesen und auf Eindeutigkeit geprüft — dieselbe Kopplung bindet
// damit auch die Vergabe an die Tabelle.
func grundCodesAusSpezifikation(t *testing.T) map[string]bool {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "spec", "spezifikation.md"))
	if err != nil {
		t.Fatalf("Spezifikation nicht lesbar (fail-closed): %v", err)
	}
	heading := regexp.MustCompile(`^## \d+\. Grund- und Fehler-Codes$`)
	row := regexp.MustCompile("^\\|\\s*`(SPEC-[0-9]{3})`\\s*\\|\\s*`([a-z0-9-]+)`\\s*\\|")
	divider := regexp.MustCompile(`^\|[\s:|-]+$`)
	lines := strings.Split(string(raw), "\n")
	start := -1
	for i, l := range lines {
		if heading.MatchString(strings.TrimSpace(l)) {
			if start >= 0 {
				t.Fatalf("Grund-Code-Überschrift mehrdeutig (Zeilen %d und %d, fail-closed)", start+1, i+1)
			}
			start = i
		}
	}
	if start < 0 {
		t.Fatal("Grund-Code-Überschrift in der Spezifikation nicht gefunden (fail-closed)")
	}
	codes := map[string]bool{}
	kennungen := map[string]bool{}
	for _, l := range lines[start+1:] {
		if strings.HasPrefix(l, "## ") {
			break
		}
		if !strings.HasPrefix(l, "|") || divider.MatchString(l) || strings.HasPrefix(l, "| Kennung |") {
			continue // keine Tabellenzeile, Trennzeile oder Kopfzeile
		}
		m := row.FindStringSubmatch(l)
		if m == nil {
			t.Fatalf("§4-Tabellenzeile ohne SPEC-Kennung und Backtick-Code (fail-closed): %q", l)
		}
		if kennungen[m[1]] {
			t.Fatalf("§4-Kennung %s doppelt vergeben (fail-closed)", m[1])
		}
		kennungen[m[1]] = true
		codes[m[2]] = true
	}
	if len(codes) == 0 {
		t.Fatal("keine Grund-Codes in der §4-Tabelle gefunden (fail-closed)")
	}
	return codes
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
	cfg := model.Config{IDPatterns: []model.IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"},
	}}
	f := model.Finding{File: "docs/a.md", Line: 1, Rule: "ids", Target: "ADR-0042", Reason: model.ReasonIDUnlinked}
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
	cfg := model.Config{}
	for _, r := range []string{model.ReasonTargetMissing, model.ReasonSpanUnclosed, rules.ReasonAnchorMissing} {
		f := model.Finding{File: "docs/a.md", Line: 1, Target: "x", Reason: r}
		if c := FixCandidateFor(f, cfg); c != nil {
			t.Errorf("Grund-Code %q sollte keinen Kandidaten liefern, got %+v", r, c)
		}
	}
}

// FixCandidateFor: id-unlinked ohne passendes Muster → kein Kandidat
// (kein Raten ohne bekanntes Definitions-Target).
func TestFixCandidateFor_OhneMusterNil(t *testing.T) {
	cfg := model.Config{IDPatterns: []model.IDPattern{
		{Regex: regexp.MustCompile(`MR-\d{3}`), Target: "harness/conventions.md"},
	}}
	f := model.Finding{File: "docs/a.md", Line: 1, Target: "ADR-0042", Reason: model.ReasonIDUnlinked}
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
