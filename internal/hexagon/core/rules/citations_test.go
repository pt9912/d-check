package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// citeSrc ist die dreizeilige Quell-Datei, gegen die die Direktiven
// zitieren (Zeile 2 trägt den langen, prüfbaren Zitatsatz).
const citeSrc = "Zeile eins mit Inhalt\n" +
	"Die zweite Zeile traegt das Zitat hier drin\n" +
	"Zeile drei\n"

func citeFindings(t *testing.T, m *coretest.MemFS) []string {
	t.Helper()
	res, err := Run(m, nil, model.Config{}, []string{"citations"})
	if err != nil {
		t.Fatalf("Run(citations) unerwarteter Fehler (kein fail-closed erwartet): %v", err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s", f.Target, f.Reason))
	}
	return got
}

// DC-FA-CITE-001 Happy inline: das per d-check:cite markierte inline-Zitat
// („…") ist ein Teilstring der normalisierten Quell-Spanne → kein Befund.
func TestCitationsInlineHappy(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:2-2 -->\n" +
			"Text davor: „Die zweite Zeile traegt das Zitat\" steht hier.\n",
	})
	if got := citeFindings(t, m); len(got) != 0 {
		t.Fatalf("korrektes inline-Zitat: %v, want 0 Befunde", got)
	}
}

// DC-FA-CITE-001 Negative: ein gedriftetes Wort bricht den Teilstring →
// citation-mismatch (die Kern-Fitness-Funktion).
func TestCitationsInlineMismatch(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:2-2 -->\n" +
			"Text davor: „Die zweite Zeile traegt DEN Zitat\" steht hier.\n",
	})
	got := citeFindings(t, m)
	want := []string{"docs/src.md:2-2 citation-mismatch"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("gedriftetes Zitat: %v, want %v", got, want)
	}
}

// DC-FA-CITE-001 Happy Block: ein >-Blockquote als Zitat, mehrzeilig,
// re-wrapped gegen die Quelle → Teilstring, kein Befund.
func TestCitationsBlockHappy(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:1-2 -->\n" +
			"> Zeile eins mit Inhalt\n" +
			"> Die zweite Zeile traegt das Zitat\n",
	})
	if got := citeFindings(t, m); len(got) != 0 {
		t.Fatalf("korrektes Block-Zitat: %v, want 0 Befunde", got)
	}
}

// DC-FA-CITE-001 Schritt 3: Bereich hinter dem Datei-Ende ⇒
// citation-out-of-range (kein Vergleich); invertierter Bereich ⇒
// citation-inverted-range; Repo-Escape ⇒ repo-escape.
func TestCitationsZitatFaeule(t *testing.T) {
	cases := []struct {
		name, directive, want string
	}{
		{"out-of-range", "docs/src.md:2-9", "docs/src.md:2-9 citation-out-of-range"},
		{"inverted", "docs/src.md:3-1", "docs/src.md:3-1 citation-inverted-range"},
		{"von-null (F-1: Befund statt Absturz)", "docs/src.md:0-3", "docs/src.md:0-3 citation-inverted-range"},
		{"escape", "../../etc/passwd:1-1", "../../etc/passwd:1-1 repo-escape"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := coretest.NewMemFS(map[string]string{
				"docs/src.md": citeSrc,
				"docs/citing.md": "<!-- d-check:cite " + c.directive + " -->\n" +
					"> irgendein folgendes Zitat als Rumpf\n",
			})
			got := citeFindings(t, m)
			if fmt.Sprint(got) != fmt.Sprint([]string{c.want}) {
				t.Fatalf("%s: %v, want [%s]", c.name, got, c.want)
			}
		})
	}
}

// DC-FA-CITE-001 Schritt 4: ein normalisiertes Zitat unter 16 Zeichen wird
// NICHT geprüft — auch wenn es kein Teilstring ist (dokumentierter
// Trade-off, keine Falsch-Rot-Gefahr).
func TestCitationsMindestlaenge(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:1-1 -->\n" +
			"Kurz: „nicht drin\" — zu kurz zum Prüfen.\n",
	})
	if got := citeFindings(t, m); len(got) != 0 {
		t.Fatalf("kurzes Zitat: %v, want 0 Befunde (ungeprüft)", got)
	}
}

// DC-FA-CITE-001 Schritt 1/2: eine strukturell unbrauchbare Direktive
// (malformter Span bzw. kein folgendes Zitat) ist fail-closed → error
// (Aufrufer mappt auf Exit 2), kein stiller Nicht-Vergleich.
func TestCitationsFailClosed(t *testing.T) {
	cases := map[string]string{
		"malformter Span": "<!-- d-check:cite docs/src.md -->\n> Zitat\n",
		"kein Zitat":      "<!-- d-check:cite docs/src.md:1-1 -->\n\n",
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			m := coretest.NewMemFS(map[string]string{
				"docs/src.md":    citeSrc,
				"docs/citing.md": body,
			})
			if _, err := Run(m, nil, model.Config{}, []string{"citations"}); err == nil {
				t.Fatalf("%s: kein Fehler — fail-closed (Exit 2) erwartet", name)
			}
		})
	}
}

// DC-FA-CITE-001 Determinismus: ohne d-check:cite-Direktive prüft das
// Modul nichts (kein Prosa-Scanning) — auch eine bloße Erwähnung des
// Marker-Strings ohne Kommentar-Form löst nichts aus.
func TestCitationsOhneDirektiveStill(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "Der Marker heißt d-check:cite und braucht die Kommentar-Form.\n" +
			"„Ein Zitat ohne Direktive\" wird nicht geprüft.\n",
	})
	if got := citeFindings(t, m); len(got) != 0 {
		t.Fatalf("ohne Direktive: %v, want 0 Befunde", got)
	}
}

// inlineQuoteSpan-Einheit: frühester Öffner gewinnt; „ paart mit ", das
// erste " paart mit dem nächsten "; ohne schließendes Paar kein Span.
func TestInlineQuoteSpanEinheit(t *testing.T) {
	cases := []struct {
		in, want string
		ok       bool
	}{
		{"Prosa davor „Zitat hier\" danach", "Zitat hier", true},
		{"gerade \"quoted text\" hier", "quoted text", true},
		{"kein Paar „offen ohne Ende", "", false},
		{"gar keine Anführung", "", false},
	}
	for _, c := range cases {
		got, ok := inlineQuoteSpan(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("inlineQuoteSpan(%q) = (%q,%v), want (%q,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

// normalizeWhitespace-Einheit: jeder Whitespace-Lauf (inkl. Umbruch) zu
// einem Leerzeichen, getrimmt; sonst keine Änderung.
func TestNormalizeWhitespaceEinheit(t *testing.T) {
	if got := normalizeWhitespace("  a\t\nb   c \n"); got != "a b c" {
		t.Fatalf("normalizeWhitespace = %q, want %q", got, "a b c")
	}
	if got := normalizeWhitespace("Satz-Zeichen, bleiben!"); got != "Satz-Zeichen, bleiben!" {
		t.Fatalf("normalizeWhitespace darf Nicht-Whitespace nicht ändern: %q", got)
	}
	// Der reale Umbruch-Fall: eine über zwei Zeilen re-wrappte Quelle
	// normalisiert zum selben String wie das einzeilige Zitat.
	if !strings.Contains(normalizeWhitespace("Die zweite Zeile\ntraegt das Zitat hier"),
		normalizeWhitespace("zweite Zeile traegt das")) {
		t.Fatal("re-wrapped Teilstring nach Normalisierung nicht gefunden")
	}
}

// DC-FA-CITE-001.a Schritt 2: ein Fenced-Block zwischen Direktive und Zitat
// trennt wie eine Leerzeile den Absatz — der Direktive folgt damit KEIN Zitat,
// und das ist fail-closed. Ohne die geteilte Absatzgrenze paart das Modul die
// Direktive mit einem Zitat hinter dem Block und meldet nichts.
func TestCitationsFenceTrenntDirektiveVomZitat(t *testing.T) {
	block := "```text\nBeispiel im Code-Block\n```\n"
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md":    citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:2-2 -->\n\n" + block + "\nText davor: „Die zweite Zeile traegt das Zitat\" steht hier.\n",
	})
	if _, err := Run(m, nil, model.Config{}, []string{"citations"}); err == nil {
		t.Fatal("Fence zwischen Direktive und Zitat: fail-closed (Exit 2) erwartet, kein Fehler geliefert")
	}
}

// Gegenprobe zur Trennung: dieselbe Datei OHNE den Block prueft das Zitat und
// bleibt befundfrei — die Trennung liegt am Fence, nicht am Zitat.
func TestCitationsOhneFencePaartNormal(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md":    citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:2-2 -->\n\nText davor: „Die zweite Zeile traegt das Zitat\" steht hier.\n",
	})
	if got := citeFindings(t, m); len(got) != 0 {
		t.Fatalf("ohne Fence: %v, want 0 Befunde", got)
	}
}

// Der Blockquote-Zweig traegt dieselbe Grenze: ein Fence zwischen Direktive und
// `>`-Block trennt ebenso.
func TestCitationsFenceTrenntVomBlockquote(t *testing.T) {
	block := "```\nBeispiel\n```\n"
	m := coretest.NewMemFS(map[string]string{
		"docs/src.md": citeSrc,
		"docs/citing.md": "<!-- d-check:cite docs/src.md:2-2 -->\n\n" + block +
			"\n> Die zweite Zeile traegt das Zitat hier drin\n",
	})
	if _, err := Run(m, nil, model.Config{}, []string{"citations"}); err == nil {
		t.Fatal("Fence vor dem Blockquote: fail-closed (Exit 2) erwartet")
	}
}
