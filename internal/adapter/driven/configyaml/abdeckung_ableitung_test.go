package configyaml_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Gezählt wird nur die Kennung im Doc-Kommentar einer Funktion, die `go test`
// als Test ausführt; jede andere Stelle und jede andere Funktion zählt nicht
// (ADR-0104).
func TestGoTestZeilenZaehltNurDenDocKommentarEinerTestfunktion(t *testing.T) {
	src := `package x

import (
	"testing"
	tt "testing"

	fremd "example.com/fremd"
)

// TestFremdesT nennt DC-FA-FRE-001, nimmt aber kein *testing.T.
func TestFremdesT(t *fremd.T) {}

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

// TestMain nennt DC-FA-MAI-001, ist aber kein Test.
func TestMain(m *testing.M) {}

// TestOhneParameter nennt DC-FA-OHN-001.
func TestOhneParameter() {}

// TestBenchmarkTyp nennt DC-FA-BTY-001.
func TestBenchmarkTyp(b *testing.B) {}

// TestAlias nennt DC-FA-ALI-001 über einen umbenannten Import.
func TestAlias(t *tt.T) {}

// Test_unterstrich nennt DC-FA-UNT-001.
func Test_unterstrich(t *testing.T) {}

// TestOhneKennung nennt nichts.
func TestOhneKennung(t *testing.T) {}
`
	rows, err := goTestZeilenAus("x_test.go", []byte(src), liveConfig(t, repoRoot()).Trace.ReqPattern)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(rows))
	for i, r := range rows {
		got[i] = r.test + "=" + strings.Join(r.ids, ",")
	}
	want := []string{"TestA=DC-FA-AAA-001,DC-QA-02", "TestAlias=DC-FA-ALI-001", "Test_unterstrich=DC-FA-UNT-001"}
	if strings.Join(got, ";") != strings.Join(want, ";") {
		t.Fatalf("Zeilen = %v, want %v", got, want)
	}
}

// Eine Datei, die ein Build-Constraint aus dem Standard-Lauf nimmt, zählt
// nicht; eine, die es nicht tut, zählt.
func TestGoTestZeilenFolgtDemBuildConstraint(t *testing.T) {
	pat := liveConfig(t, repoRoot()).Trace.ReqPattern
	body := "package x\n\nimport \"testing\"\n\n// TestX nennt DC-FA-BLD-001.\nfunc TestX(t *testing.T) {}\n"
	for _, tc := range []struct {
		constraint string
		want       int
	}{
		{"", 1},
		{"//go:build integration\n\n", 0},
		{"//go:build ignore\n\n", 0},
		{"//go:build !integration\n\n", 1},
	} {
		rows, err := goTestZeilenAus("x_test.go", []byte(tc.constraint+body), pat)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != tc.want {
			t.Errorf("Constraint %q: Zeilen = %d, want %d", tc.constraint, len(rows), tc.want)
		}
	}
}

// Der Walk überspringt, was `go test ./...` überspringt.
func TestGoTestZeilenUeberspringtVerzeichnisseWieGoTest(t *testing.T) {
	root := t.TempDir()
	test := "package x\n\nimport \"testing\"\n\n// TestX nennt DC-FA-WLK-001.\nfunc TestX(t *testing.T) {}\n"
	for _, dir := range []string{"internal/a", "internal/testdata", "internal/_skip", "internal/.hidden", "cmd/c"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, dir, "x_test.go"), []byte(test), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := goTestZeilen(root, liveConfig(t, repoRoot()).Trace.ReqPattern)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, r := range rows {
		files = append(files, r.file)
	}
	if strings.Join(files, ";") != "internal/a/x_test.go;cmd/c/x_test.go" {
		t.Fatalf("gelesen = %v", files)
	}
}

// Jede Phase von image-test.sh trägt einen Anker mit Kennung, auch wenn ihre
// Kopfzeile anders umbrochen ist; fehlt er oder steht er ohne Phase, scheitert
// die Ableitung (fail-closed).
func TestImageTestZeilenVerlangtDenAnkerJederPhase(t *testing.T) {
	pat := liveConfig(t, repoRoot()).Trace.ReqPattern
	gut := "# --- (1) Happy: a ----\n# abdeckung: DC-FA-DIST-001, DC-QA-02\ncode\n" +
		"# --- (2) Boundary: b ---\n# abdeckung: DC-QA-03\n"
	rows, err := imageTestZeilenAus("s.sh", gut, pat)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].test != "(1) Happy: a" || strings.Join(rows[0].ids, ",") != "DC-FA-DIST-001,DC-QA-02" {
		t.Fatalf("Zeilen = %+v", rows)
	}
	for _, src := range []string{
		"# --- (1) Happy ---\ncode\n",
		"# --- (1) Happy ---\n",
		"# --- (1) Happy\ncode\n",
		"# --- (1) Happy --- \ncode\n",
		"#--- (1) Happy ---\ncode\n",
		"code\n# abdeckung: DC-QA-02\n",
		"# --- (1) Happy ---\n# abdeckung: keine\n",
	} {
		if _, err := imageTestZeilenAus("s.sh", src, pat); err == nil {
			t.Errorf("%q: kein Fehler", src)
		}
	}
}

// trace.coverage bindet genau die Abdeckungs-Dateien der Liste ein.
func TestAbdeckungsListeDecktCoverage(t *testing.T) {
	liste := []abdeckungsDatei{{path: "docs/user/abdeckung-a.md"}}
	for i, tc := range []struct {
		cov  []model.TraceCoverage
		fehl bool
	}{
		{[]model.TraceCoverage{{Files: []string{"docs/user/abdeckung-a.md"}}}, false},
		{nil, true},
		{[]model.TraceCoverage{{Files: []string{"docs/user/abdeckung-a.md", "docs/user/abdeckung-b.md"}}}, true},
		{[]model.TraceCoverage{{Files: []string{"docs/user/abdeckung-a.md"}}, {Files: []string{"spec/x.md"}}}, false},
	} {
		if err := abdeckungsListeDecktCoverage(liste, tc.cov); (err != nil) != tc.fehl {
			t.Errorf("Fall %d: err = %v, Fehler erwartet %v", i, err, tc.fehl)
		}
	}
}

// Eine Kennung außerhalb der Tabelle würde die RTM still mitzählen.
func TestNurTabellenKennungen(t *testing.T) {
	pat := liveConfig(t, repoRoot()).Trace.ReqPattern
	rows := []abdeckungsZeile{{ids: []string{"DC-QA-02"}}}
	if err := nurTabellenKennungen("Tabelle DC-QA-02", rows, pat); err != nil {
		t.Errorf("nur Tabellen-Kennung: %v", err)
	}
	if err := nurTabellenKennungen("Prosa nennt DC-QA-03, Tabelle DC-QA-02", rows, pat); err == nil {
		t.Error("Kennung außerhalb der Tabelle: kein Fehler")
	}
}
