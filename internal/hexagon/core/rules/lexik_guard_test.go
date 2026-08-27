package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Der Wächter, den ADR-0054s Re-Evaluierungs-Trigger verlangt: *kann ein Gate
// die Klasse prüfen statt eines Reviews?* Die Klasse ist „ein Konsument
// beantwortet eine geteilte Lexik-Frage selbst und anders" — ein stiller
// Grün-Pfad, den kein Doku-Gate sieht, weil jedes Modul für sich konsistent
// ist. Bisher hing die Prüfung an einem Review; hier hängt sie an einem Test.
//
// GEGENSTAND ist die Frage „ist diese Zeile eine Direktive": jeder Zugriff auf
// den Direktiven-Marker auf einer ROHEN Zeile. Die geteilte Antwort steht in
// stripInlineCodeByLine; wer sie umgeht, steht unten mit Grund.
//
// GRENZE, benannt: der Wächter liest den QUELLTEXT dieses Pakets, nicht das
// Verhalten. Er fängt die Schreibweise `pl.raw`/`raw` neben dem Marker — nicht
// eine Umgehung über eine Zwischenvariable oder ein anderes Paket. Er ist ein
// Stolperdraht für den nächsten Konsumenten, keine Beweisführung.
const (
	markerName    = "ignoreMarker"
	markerLiteral = "d-check:ignore"
)

// rohLesenErlaubt nennt je Datei, warum der Marker dort auf der ROHEN Zeile
// gelesen wird — und die Begründung ist in beiden Fällen dieselbe Klasse: die
// Eingabe ist keine Prosa, also gibt es keine geteilte Prosa-Antwort, die zu
// übernehmen wäre (ADR-0054 §Entscheidung 2, „andere Frage").
func rohLesenErlaubt() map[string]string {
	return map[string]string{
		// versions liest ALLE Zeilen, einschliesslich Fenced-Code
		// (ADR-0019) — dort ist ein Backtick literaler Inhalt.
		"versions.go": "liest alle Zeilen inkl. Fences (ADR-0019)",
		// diagrams liest die Zeilen INNERHALB eines Fence; dort gibt es
		// kein Inline-Code.
		"diagrams.go": "liest Zeilen innerhalb eines Fence",
	}
}

func TestLexikGuard_KeinKonsumentLiestRohOhneGrund(t *testing.T) {
	eintraege, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Paket-Verzeichnis nicht lesbar: %v", err)
	}
	erlaubt := rohLesenErlaubt()
	gesehen := map[string]bool{}
	for _, e := range eintraege {
		name := e.Name()
		if e.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("%s nicht lesbar: %v", name, err)
		}
		for i, zeile := range strings.Split(string(b), "\n") {
			// Der Marker zaehlt unter BEIDEN Schreibweisen — als Konstante
			// und als Literal. Nur das Literal zu vergessen waere die
			// billigste Umgehung ueberhaupt.
			if !strings.Contains(zeile, markerName) && !strings.Contains(zeile, markerLiteral) {
				continue
			}
			// Nur der LESENDE Zugriff zaehlt, nicht die Deklaration der
			// Konstanten und nicht Kommentare. Gefangen sind die zwei
			// gaengigen Lese-Aufrufe; wer den Marker anders sucht (regexp,
			// eigene Schleife), steht nicht darin — siehe GRENZE oben.
			if !strings.Contains(zeile, "strings.Contains(") && !strings.Contains(zeile, "strings.Index(") {
				continue
			}
			if !strings.Contains(zeile, ".raw") && !strings.Contains(zeile, "(raw,") {
				continue
			}
			grund, ok := erlaubt[name]
			if !ok {
				t.Errorf("%s:%d liest %s auf der rohen Zeile ohne Eintrag in rohLesenErlaubt() — "+
					"„ist das eine Direktive\" ist eine Prosa-Frage und bekommt die geteilte "+
					"Antwort (stripInlineCodeByLine, ADR-0054/ADR-0062). Ist die Eingabe "+
					"KEINE Prosa, gehoert die Datei mit Grund in die Liste", name, i+1, markerName)
				continue
			}
			gesehen[name] = true
			_ = grund
		}
	}
	// Gegenrichtung: ein Eintrag, den niemand mehr braucht, ist eine
	// Erlaubnis ohne Gegenstand — sie verdeckt den naechsten Fall.
	for name := range erlaubt {
		if !gesehen[name] {
			t.Errorf("rohLesenErlaubt() nennt %q, aber dort liest niemand mehr roh — "+
				"Eintrag entfernen, sonst deckt er kuenftige Umgehungen", name)
		}
	}
}
