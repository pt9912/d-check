package configyaml_test

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Die Abdeckungs-Dateien der RTM (trace.coverage, ADR-0104) sind eine
// Ableitung der Testquellen: je Nachweisart eine Tabelle Kennung → Test →
// Datei. Der Wächter hält jede committete Datei gegen ihre Ableitung; mit
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

// abdeckungsDateien ist die Liste der Abdeckungs-Dateien; der Wächter hält sie
// gegen trace.coverage der .d-check.yml.
func abdeckungsDateien() []abdeckungsDatei {
	return []abdeckungsDatei{
		{
			path:  "docs/user/abdeckung-tests.md",
			title: "Test-Abdeckung je Anforderung (Go-Suite)",
			intro: "Abgeleitet aus den Go-Tests unter `internal/` und `cmd/`: eine Zeile je\n" +
				"Testfunktion, deren Doc-Kommentar unmittelbar über `func Test…` eine\n" +
				"Anforderungs-Kennung nennt. Gezählt wird, was `go test` unter linux/amd64\n" +
				"ohne `-tags` als Test ausführt; eine Kennung an anderer Stelle des Tests\n" +
				"oder im Datei-Kommentar zählt nicht. Die Datei schreibt `make abdeckung`;\n" +
				"`make test` hält sie gegen die Testquellen.\n\n" +
				"**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass der Test die\n" +
				"Anforderung prüfen soll, nicht, dass er es tut.",
			derive:  goTestZeilen,
			testCol: "Test",
		},
		{
			path:  "docs/user/abdeckung-e2e.md",
			title: "E2E-Abdeckung je Anforderung (Image-Test)",
			intro: "Abgeleitet aus `tools/image-test.sh` (`make image-test`): eine Zeile je\n" +
				"Phase. Eine Phase ist eine Kommentarzeile mit mindestens drei Strichen vor\n" +
				"einer Nummer in Klammern; unter ihr steht ein Anker `# abdeckung:` mit ihren\n" +
				"Kennungen, sonst ist `make test` rot. Geprüft wird das gebaute Image, nativ\n" +
				"gegen Container. Die Datei schreibt `make abdeckung`; `make test` hält sie\n" +
				"gegen das Skript.\n\n" +
				"**Grenze:** eine Deklaration, kein Beleg — die Zeile sagt, dass die Phase\n" +
				"die Anforderung prüfen soll, nicht, dass sie es tut.",
			derive:  imageTestZeilen,
			testCol: "Phase",
		},
	}
}

// TestAbdeckungsDateienFolgenIhrerAbleitung hält die Abdeckungs-Dateien der
// RTM gegen ihre Ableitung aus den Testquellen (ADR-0104).
func TestAbdeckungsDateienFolgenIhrerAbleitung(t *testing.T) {
	root := repoRoot()
	cfg := liveConfig(t, root)
	pat := cfg.Trace.ReqPattern
	if err := abdeckungsListeDecktCoverage(abdeckungsDateien(), cfg.Trace.Coverage); err != nil {
		t.Fatal(err)
	}
	ziel := os.Getenv("ABDECKUNG_ZIEL")
	for _, d := range abdeckungsDateien() {
		rows, err := d.derive(root, pat)
		if err != nil {
			t.Fatalf("%s: Ableitung scheitert: %v", d.path, err)
		}
		want := renderAbdeckung(d, rows)
		if err := nurTabellenKennungen(want, rows, pat); err != nil {
			t.Fatalf("%s: %v", d.path, err)
		}
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

// liveConfig dekodiert die .d-check.yml des Repos — dieselbe Erkennung wie
// die RTM, kein zweites Muster.
func liveConfig(t *testing.T, root string) model.Config {
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
	return cfg
}

// abdeckungsListeDecktCoverage verlangt, dass trace.coverage jede Datei dieser
// Liste einbindet und keine weitere unter `docs/user/abdeckung-`.
// GRENZE: Eine Coverage-Quelle unter anderem Pfad prüft der Abgleich nicht,
// ebenso wenig Labels und die Aufteilung der Dateien auf Quellen.
func abdeckungsListeDecktCoverage(liste []abdeckungsDatei, coverage []model.TraceCoverage) error {
	want := map[string]bool{}
	for _, d := range liste {
		want[d.path] = true
	}
	got := map[string]bool{}
	for _, c := range coverage {
		for _, f := range c.Files {
			got[f] = true
		}
	}
	for _, p := range sortierteSchluessel(want) {
		if !got[p] {
			return fmt.Errorf("trace.coverage bindet %s nicht ein", p)
		}
	}
	for _, p := range sortierteSchluessel(got) {
		if strings.HasPrefix(p, "docs/user/abdeckung-") && !want[p] {
			return fmt.Errorf("trace.coverage bindet %s ein, die Ableitung kennt sie nicht", p)
		}
	}
	return nil
}

// sortierteSchluessel liefert die Schlüssel einer Menge sortiert — die
// Meldung bei mehreren Abweichungen ist so jedes Mal dieselbe.
func sortierteSchluessel(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// nurTabellenKennungen verlangt, dass jede Kennung der Datei aus der Tabelle
// stammt: die RTM liest jede Kennung der ganzen Datei, auch in Prosa.
func nurTabellenKennungen(doc string, rows []abdeckungsZeile, pat *regexp.Regexp) error {
	inTabelle := map[string]bool{}
	for _, r := range rows {
		for _, id := range r.ids {
			inTabelle[id] = true
		}
	}
	for _, id := range kennungen(doc, pat) {
		if !inTabelle[id] {
			return fmt.Errorf("Kennung %s steht in der Datei, aber in keiner Ableitungs-Zeile", id)
		}
	}
	return nil
}

// goTestZeilen liest alle *_test.go unter internal/ und cmd/ und überspringt,
// was `go test ./...` überspringt: testdata und Verzeichnisse mit `_` oder `.`
// am Anfang.
func goTestZeilen(root string, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	var rows []abdeckungsZeile
	for _, dir := range []string{"internal", "cmd"} {
		start := filepath.Join(root, dir)
		err := filepath.WalkDir(start, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if e.IsDir() {
				if p != start && (e.Name() == "testdata" || strings.HasPrefix(e.Name(), "_") || strings.HasPrefix(e.Name(), ".")) {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(p, "_test.go") {
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
// Testfunktion, die `go test` ausführt (Name nach seiner Regel, genau ein
// Parameter `*testing.T`, keine Methode, Datei im Standard-Build), deren
// Doc-Kommentar eine Kennung nennt.
func goTestZeilenAus(file string, src []byte, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	f, err := parser.ParseFile(token.NewFileSet(), file, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	ok, err := gebautVonGoTest(file, src)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	if !ok {
		return nil, nil
	}
	testingNames := testingImportNames(f)
	var rows []abdeckungsZeile
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Doc == nil || !istTestName(fd.Name.Name) || !nimmtTestingT(fd, testingNames) {
			continue
		}
		if ids := kennungen(fd.Doc.Text(), pat); len(ids) > 0 {
			rows = append(rows, abdeckungsZeile{ids: ids, test: fd.Name.Name, file: file})
		}
	}
	return rows, nil
}

// gebautVonGoTest fragt `go/build`, ob `go test` die Datei baut — dieselbe
// Antwort wie das Go-Werkzeug: Dateiname (GOOS/GOARCH-Suffix, `_` und `.` am
// Anfang), `//go:build` und `// +build`. Der Kontext ist fest linux/amd64
// (Architektur-Stufe v1) ohne cgo und ohne `-tags`, wie `make test` in der CI
// läuft; die Datei ist so auf jedem Rechner dieselbe.
// GRENZE: Eine Datei, die nur unter einer anderen Plattform oder mit `-tags`
// gebaut wird, zählt nicht — auch dann nicht, wenn `make test` auf einem
// arm64-Rechner sie baut.
func gebautVonGoTest(file string, src []byte) (bool, error) {
	ctx := build.Default
	ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled, ctx.BuildTags = "linux", "amd64", false, nil
	ctx.ToolTags = []string{"amd64.v1"}
	ctx.OpenFile = func(string) (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(src)), nil }
	return ctx.MatchFile(path.Dir(file), path.Base(file))
}

// testingImportNames liefert die Namen, unter denen die Datei `testing`
// importiert — eine Datei kann es mehrfach tun, auch unter einem Alias.
func testingImportNames(f *ast.File) map[string]bool {
	names := map[string]bool{}
	for _, imp := range f.Imports {
		if imp.Path.Value != `"testing"` {
			continue
		}
		if imp.Name != nil {
			names[imp.Name.Name] = true
		} else {
			names["testing"] = true
		}
	}
	return names
}

// nimmtTestingT verlangt genau einen Parameter vom Typ `*testing.T`.
func nimmtTestingT(fd *ast.FuncDecl, testingNames map[string]bool) bool {
	params := fd.Type.Params.List
	if len(params) != 1 || len(params[0].Names) > 1 {
		return false
	}
	star, ok := params[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "T" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && testingNames[pkg.Name]
}

// istTestName folgt der Namensregel von `go test`: `Test`, dahinter nichts
// oder ein Zeichen, das kein Kleinbuchstabe ist.
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
	imageTestPhaseRE = regexp.MustCompile(`^\s*#\s*-{3,}\s*(\(\d+[a-z]?\).*?)[\s-]*$`)
	imageTestAnkerRE = regexp.MustCompile(`^\s*#\s*abdeckung:(.*)$`)
)

// imageTestZeilen liest die Phasen von tools/image-test.sh.
func imageTestZeilen(root string, pat *regexp.Regexp) ([]abdeckungsZeile, error) {
	src, err := os.ReadFile(filepath.Join(root, "tools/image-test.sh"))
	if err != nil {
		return nil, err
	}
	return imageTestZeilenAus("tools/image-test.sh", string(src), pat)
}

// imageTestZeilenAus verlangt unter jeder Phasen-Kopfzeile — ein Kommentar mit
// mindestens drei Strichen vor einer Nummer in Klammern — einen Anker mit
// mindestens einer Kennung; eine Phase ohne Anker und ein Anker ohne Phase
// darüber sind ein Fehler (fail-closed).
// GRENZE: Eine Kopfzeile, die dem Muster nicht folgt — keine drei Striche,
// eine Klammer wie `(4.1)` oder `(4B)`, eine Phase ohne Klammer —, ist keine
// Phase; ohne Anker fällt sie still aus der Ableitung.
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
