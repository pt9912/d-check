package configyaml_test

import (
	"regexp"
	"strings"
	"testing"
)

func abdeckungTestPattern() *regexp.Regexp {
	return regexp.MustCompile(`[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?`)
}

// Gezählt wird nur die Kennung im Doc-Kommentar einer Testfunktion; jede
// andere Stelle und jede andere Funktion zählt nicht (DC-FA-COV-001).
func TestGoTestZeilenZaehltNurDenDocKommentarEinerTestfunktion(t *testing.T) {
	src := `package x

import "testing"

// TestA prüft DC-FA-AAA-001 und DC-QA-02, doppelt DC-QA-02.
func TestA(t *testing.T) {}

func TestRumpf(t *testing.T) {
	// DC-FA-RUM-001 im Rumpf
	t.Run("DC-FA-RUN-001", func(t *testing.T) {})
	t.Fatal("DC-FA-MSG-001")
}

// DC-FA-LEE-001 mit Leerzeile davor

func TestLeerzeile(t *testing.T) {}

func TestDC_FA_NAM_001(t *testing.T) {}

// hilfe nennt DC-FA-HLF-001
func hilfe() {}

type s struct{}

// TestMethode nennt DC-FA-MTH-001
func (s) TestMethode(t *testing.T) {}

// Testing nennt DC-FA-ING-001, ist aber kein Testname.
func Testing(t *testing.T) {}

// Test_unterstrich nennt DC-FA-UNT-001.
func Test_unterstrich(t *testing.T) {}

// TestOhneKennung nennt nichts.
func TestOhneKennung(t *testing.T) {}
`
	rows, err := goTestZeilenAus("x_test.go", []byte(src), abdeckungTestPattern())
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(rows))
	for i, r := range rows {
		got[i] = r.test + "=" + strings.Join(r.ids, ",")
	}
	want := []string{"TestA=DC-FA-AAA-001,DC-QA-02", "Test_unterstrich=DC-FA-UNT-001"}
	if strings.Join(got, ";") != strings.Join(want, ";") {
		t.Fatalf("Zeilen = %v, want %v", got, want)
	}
}

// Jede Phase von image-test.sh trägt einen Anker mit Kennung; fehlt er oder
// steht er ohne Phase, scheitert die Ableitung (fail-closed).
func TestImageTestZeilenVerlangtDenAnkerJederPhase(t *testing.T) {
	pat := abdeckungTestPattern()
	gut := "# --- (1) Happy: a ----\n# abdeckung: DC-FA-DIST-001, DC-QA-02\ncode\n" +
		"# --- (2) Boundary: b ---\n# abdeckung: DC-QA-03\n"
	rows, err := imageTestZeilenAus("s.sh", gut, pat)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].test != "(1) Happy: a" || strings.Join(rows[0].ids, ",") != "DC-FA-DIST-001,DC-QA-02" {
		t.Fatalf("Zeilen = %+v", rows)
	}
	for name, src := range map[string]string{
		"phase ohne anker":         "# --- (1) Happy ---\ncode\n",
		"phase ohne anker am ende": "# --- (1) Happy ---\n",
		"anker ohne phase":         "code\n# abdeckung: DC-QA-02\n",
		"anker ohne kennung":       "# --- (1) Happy ---\n# abdeckung: keine\n",
	} {
		if _, err := imageTestZeilenAus("s.sh", src, pat); err == nil {
			t.Errorf("%s: kein Fehler", name)
		}
	}
}
