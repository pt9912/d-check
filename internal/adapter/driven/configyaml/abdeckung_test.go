package configyaml_test

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

// Die Abdeckungs-Dateien der RTM (trace.coverage, DC-FA-COV-001) sind eine
// Ableitung der Testquellen: je Nachweisart eine Tabelle Kennung → Test →
// Datei. Dieser Test hält jede committete Datei gegen ihre Ableitung; mit
// ABDECKUNG_ZIEL=<verzeichnis> schreibt er sie stattdessen dorthin
// (`make abdeckung`).
//
// GRENZE: Gezählt wird eine Deklaration, kein Beleg — eine Kennung im
// Doc-Kommentar eines Tests sagt, dass er die Anforderung prüfen soll, nicht,
// dass er es tut.

// abdeckungsZeile ist eine Zeile einer Abdeckungs-Tabelle.
type abdeckungsZeile struct {
	ids  []string
	test string
	file string
}

// abdeckungsDatei beschreibt eine Abdeckungs-Datei und ihre Ableitung.
type abdeckungsDatei struct {
	path    string
	title   string
	intro   string
	derive  func(root string, pat *regexp.Regexp) ([]abdeckungsZeile, error)
	testCol string
}

func abdeckungsDateien() []abdeckungsDatei {
	return []abdeckungsDatei{
		{
			path:  "docs/user/abdeckung-tests.md",
			title: "Test-Abdeckung je Anforderung (Go-Suite)",
			intro: "Abgeleitet aus den Go-Tests unter `internal/` und `cmd/`: eine Zeile je\n" +
				"Testfunktion, deren Doc-Kommentar unmittelbar über `func Test…` eine\n" +
				"Anforderungs-Kennung nennt. Eine Kennung an anderer Stelle des Tests zählt\n" +
				"nicht. Die Datei schreibt `make abdeckung`; `make test` hält sie gegen die\n" +
				"Testquellen.\n\n" +
				"**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass der Test die\n" +
				"Anforderung prüfen soll, nicht, dass er es tut.",
			derive:  goTestZeilen,
			testCol: "Test",
		},
		{
			path:  "docs/user/abdeckung-e2e.md",
			title: "E2E-Abdeckung je Anforderung (Image-Test)",
			intro: "Abgeleitet aus `tools/image-test.sh` (`make image-test`): eine Zeile je\n" +
				"Phase, die unter ihrer Kopfzeile einen Anker `# abdeckung:` mit ihren\n" +
				"Kennungen trägt. Geprüft wird das gebaute Image, nativ gegen Container. Die\n" +
				"Datei schreibt `make abdeckung`; `make test` hält sie gegen das Skript.\n\n" +
				"**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass die Phase\n" +
				"die Anforderung prüfen soll, nicht, dass sie es tut.",
			derive:  imageTestZeilen,
			testCol: "Phase",
		},
	}
}

// TestAbdeckungsDateienFolgenIhrerAbleitung hält die Abdeckungs-Dateien der
// RTM gegen ihre Ableitung aus den Testquellen (DC-FA-COV-001).
func TestAbdeckungsDateienFolgenIhrerAbleitung(t *testing.T) {
	root := repoRoot()
	pat := liveReqPattern(t, root)
	ziel := os.Getenv("ABDECKUNG_ZIEL")
	for _, d := range abdeckungsDateien() {
		rows, err := d.derive(root, pat)
		if err != nil {
			t.Fatalf("%s: Ableitung scheitert: %v", d.path, err)
		}
		want := renderAbdeckung(d, rows)
		if ziel != "" {
			out := filepath.Join(ziel, filepath.Base(d.path))
			if err := os.WriteFile(out, []byte(want), 0o600); err != nil {
				t.Fatalf("%s: schreiben: %v", out, err)
			}
			continue
		}
		got, err := os.ReadFile(filepath.Join(root, d.path))
		if err != nil {
			t.Fatalf("%s fehlt (make abdeckung): %v", d.path, err)
		}
		if string(got) != want {
			t.Errorf("%s weicht von der Ableitung ab — make abdeckung schreibt sie neu", d.path)
		}
	}
}

// liveReqPattern liest trace.requirements.id-pattern aus der .d-check.yml des
// Repos — dieselbe Erkennung wie die RTM, kein zweites Muster.
func liveReqPattern(t *testing.T, root string) *regexp.Regexp {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, ".d-check.yml"))
	if err != nil {
		t.Fatalf(".d-check.yml nicht lesbar (fail-closed): %v", err)
	}
	cfg, err := configyaml.Decode(raw)
	if err != nil {
		t.Fatalf(".d-check.yml dekodiert nicht (fail-closed): %v", err)
	}
	if cfg.Trace.ReqPattern == nil {
		t.Fatal("trace.requirements.id-pattern fehlt (fail-closed)")
	}
	return cfg.Trace.ReqPattern
}

// goTestZeilen liest alle *_test.go unter internal/ und cmd/ (ohne testdata).
func goTestZeilen(root string, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	var rows []abdeckungsZeile
	for _, dir := range []string{"internal", "cmd"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if e.IsDir() && e.Name() == "testdata" {
				return filepath.SkipDir
			}
			if e.IsDir() || !strings.HasSuffix(p, "_test.go") {
				return nil
			}
			src, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}
			r, err := goTestZeilenAus(filepath.ToSlash(rel), src, pat)
			if err != nil {
				return err
			}
			rows = append(rows, r...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

// goTestZeilenAus liest eine Testdatei mit dem Go-Parser: eine Zeile je
// Testfunktion (`func TestXxx(t *testing.T)`, keine Methode), deren
// Doc-Kommentar eine Kennung nennt.
func goTestZeilenAus(file string, src []byte, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	f, err := parser.ParseFile(token.NewFileSet(), file, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	var rows []abdeckungsZeile
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Doc == nil || !istTestName(fd.Name.Name) {
			continue
		}
		if ids := kennungen(fd.Doc.Text(), pat); len(ids) > 0 {
			rows = append(rows, abdeckungsZeile{ids: ids, test: fd.Name.Name, file: file})
		}
	}
	return rows, nil
}

// istTestName folgt der Regel von `go test`: `Test`, dahinter nichts oder ein
// Zeichen, das kein Kleinbuchstabe ist.
func istTestName(name string) bool {
	if !strings.HasPrefix(name, "Test") {
		return false
	}
	rest := name[len("Test"):]
	if rest == "" {
		return true
	}
	r, _ := utf8.DecodeRuneInString(rest)
	return !unicode.IsLower(r)
}

var (
	imageTestPhaseRE = regexp.MustCompile(`^# --- (\(\d+\) .*?) -+$`)
	imageTestAnkerRE = regexp.MustCompile(`^# abdeckung: (.+)$`)
)

// imageTestZeilen liest die Phasen von tools/image-test.sh.
func imageTestZeilen(root string, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	src, err := os.ReadFile(filepath.Join(root, "tools/image-test.sh"))
	if err != nil {
		return nil, err
	}
	return imageTestZeilenAus("tools/image-test.sh", string(src), pat)
}

// imageTestZeilenAus verlangt unter jeder Phasen-Kopfzeile einen Anker mit
// mindestens einer Kennung; eine Phase ohne Anker und ein Anker ohne Phase
// darüber sind ein Fehler (fail-closed).
func imageTestZeilenAus(file, src string, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	var rows []abdeckungsZeile
	phase := ""
	sc := bufio.NewScanner(strings.NewReader(src))
	for no := 1; sc.Scan(); no++ {
		line := sc.Text()
		if m := imageTestAnkerRE.FindStringSubmatch(line); m != nil {
			if phase == "" {
				return nil, fmt.Errorf("%s:%d: Anker ohne Phasen-Kopfzeile darüber", file, no)
			}
			ids := kennungen(m[1], pat)
			if len(ids) == 0 {
				return nil, fmt.Errorf("%s:%d: Anker ohne Kennung", file, no)
			}
			rows = append(rows, abdeckungsZeile{ids: ids, test: phase, file: file})
			phase = ""
			continue
		}
		if phase != "" {
			return nil, fmt.Errorf("%s:%d: Phase %q ohne Anker darunter", file, no-1, phase)
		}
		if m := imageTestPhaseRE.FindStringSubmatch(line); m != nil {
			phase = m[1]
		}
	}
	if phase != "" {
		return nil, fmt.Errorf("%s: Phase %q ohne Anker darunter", file, phase)
	}
	return rows, sc.Err()
}

// kennungen liefert die Kennungen eines Texts, sortiert und ohne Dubletten.
func kennungen(text string, pat *regexp.Regexp) []string {
	seen := map[string]bool{}
	var ids []string
	for _, id := range pat.FindAllString(text, -1) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// renderAbdeckung schreibt die Tabelle, sortiert nach Datei und Test.
func renderAbdeckung(d abdeckungsDatei, rows []abdeckungsZeile) string {
	sort.SliceStable(rows, func(a, b int) bool {
		if rows[a].file != rows[b].file {
			return rows[a].file < rows[b].file
		}
		return rows[a].test < rows[b].test
	})
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n%s\n\n", d.title, d.intro)
	fmt.Fprintf(&b, "| Kennung | %s | Datei |\n| --- | --- | --- |\n", d.testCol)
	for _, r := range rows {
		links := make([]string, len(r.ids))
		for i, id := range r.ids {
			links[i] = fmt.Sprintf("[`%s`](../../spec/lastenheft.md)", id)
		}
		fmt.Fprintf(&b, "| %s | `%s` | [`%s`](../../%s) |\n", strings.Join(links, ", "), r.test, r.file, r.file)
	}
	return b.String()
}
