package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// immDoc baut eine ADR-artige Fixture: der Pin-Marker steht stets auf Zeile 4,
// die Entscheidungs-Zeile (Core) auf Zeile 8, der ## Geschichte-Abschnitt am
// Ende (ausnehmbar). extraGeschichte hängt eine Zeile unter ## Geschichte an.
func immDoc(pin, decision, extraGeschichte string) string {
	ls := []string{
		"# ADR-0099 — X",       // 1
		"",                     // 2
		"**Status:** Accepted", // 3
		"<!-- immutable: sha256:" + pin + " -->", // 4 (Marker)
		"",               // 5
		"## Entscheidung", // 6
		"",                // 7
		decision,          // 8 (Core)
		"",                // 9
		"## Geschichte",   // 10
		"",                // 11
		"| Datum | Ereignis |",                  // 12
		"| --- | --- |",                         // 13
		"| 2026-01-01 | Proposed -> Accepted |", // 14
	}
	if extraGeschichte != "" {
		ls = append(ls, extraGeschichte)
	}
	return strings.Join(ls, "\n") + "\n"
}

const immMarkerLine = 4

func runImmutable(t *testing.T, content string, excl []string) []model.Finding {
	t.Helper()
	c := []byte(content)
	return CheckImmutable("adr.md", PreprocessMarkdown(c), c, model.ImmutableConfig{ExcludeSections: excl})
}

// correctPin liefert den für `decision` korrekten Pin (Geschichte ausgenommen).
func correctPin(decision string) string {
	return immutableCoreHash([]byte(immDoc("0", decision, "")), immMarkerLine, []string{"Geschichte"})
}

func TestImmutableHappyPath(t *testing.T) {
	good := immDoc(correctPin("Tue A."), "Tue A.", "")
	if fs := runImmutable(t, good, []string{"Geschichte"}); len(fs) != 0 {
		t.Fatalf("korrekter Pin: erwartet 0 Befunde, got %d: %+v", len(fs), fs)
	}
}

func TestImmutableReflowBoundary(t *testing.T) {
	// Nur-Whitespace-Änderung am Core (Wort-Inhalt identisch) → kein Befund.
	reflow := immDoc(correctPin("Tue A."), "Tue  A.", "")
	if fs := runImmutable(t, reflow, []string{"Geschichte"}); len(fs) != 0 {
		t.Fatalf("Reflow am Core: erwartet 0 Befunde, got %d: %+v", len(fs), fs)
	}
}

func TestImmutableExcludedSection(t *testing.T) {
	// Anhang unter ## Geschichte (ausgenommen) → kein Befund.
	appended := immDoc(correctPin("Tue A."), "Tue A.", "| 2026-02-02 | Notiz |")
	if fs := runImmutable(t, appended, []string{"Geschichte"}); len(fs) != 0 {
		t.Fatalf("Geschichte-Anhang: erwartet 0 Befunde, got %d: %+v", len(fs), fs)
	}
}

func TestImmutableNegativeCoreEdit(t *testing.T) {
	bad := immDoc(correctPin("Tue A."), "Tue B.", "") // Core geändert, Pin unverändert
	fs := runImmutable(t, bad, []string{"Geschichte"})
	if len(fs) != 1 {
		t.Fatalf("Core-Edit: erwartet 1 Befund, got %d: %+v", len(fs), fs)
	}
	if fs[0].Reason != model.ReasonCoreDrift || fs[0].Rule != "immutable" || fs[0].Line != immMarkerLine {
		t.Fatalf("Befund-Form falsch: %+v", fs[0])
	}
}

func TestImmutableNoMarker(t *testing.T) {
	noMarker := strings.Replace(immDoc("0", "Tue A.", ""), "<!-- immutable: sha256:0 -->", "", 1)
	if fs := runImmutable(t, noMarker, []string{"Geschichte"}); len(fs) != 0 {
		t.Fatalf("ohne Marker: erwartet 0 Befunde, got %d: %+v", len(fs), fs)
	}
}

func TestImmutableMarkerInFenceInert(t *testing.T) {
	// Marker NUR in einem Fenced-Code-Block (Syntax-Beispiel) → kein Live-Pin.
	fenced := strings.Join([]string{
		"# Doku", "", "Beispiel:", "",
		"```", "<!-- immutable: sha256:abc123 -->", "```",
		"", "Körper.",
	}, "\n") + "\n"
	if fs := runImmutable(t, fenced, nil); len(fs) != 0 {
		t.Fatalf("Marker in Fence: erwartet 0 Befunde (inert), got %d: %+v", len(fs), fs)
	}
}

func TestImmutableFirstMarkerWins(t *testing.T) {
	// Zwei Marker: nur der erste (Zeile 4) ist der Pin. Erster falsch ⇒ genau
	// ein Befund, verankert an Zeile 4 (der zweite zählt nicht als Pin).
	two := strings.Join([]string{
		"# ADR-0099 — X", "", "**Status:** Accepted",
		"<!-- immutable: sha256:0 -->", // 4: erster (Pin), falsch
		"", "## Entscheidung", "",
		"<!-- immutable: sha256:ff -->", // 8: zweiter, inert als Pin
		"", "Tue A.",
	}, "\n") + "\n"
	fs := runImmutable(t, two, nil)
	if len(fs) != 1 {
		t.Fatalf("zwei Marker: erwartet 1 Befund, got %d: %+v", len(fs), fs)
	}
	if fs[0].Line != 4 {
		t.Fatalf("Befund sollte am ersten Marker (Zeile 4) hängen, got Zeile %d", fs[0].Line)
	}
}
