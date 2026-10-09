package rules

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// stubMarker ist die Kopfzeile eines archivierten Stubs; das skip-pattern der
// Tests trifft genau sie.
const (
	stubMarker  = "> **ARCHIVIERT** — Volltext im Archiv.\n"
	stubPattern = `(?m)^> \*\*ARCHIVIERT\*\*`
	kaputteNote = "# Slice\n\n## 7. Closure-Notiz\n\nKurz.\n"
)

// readErrOneFS lässt ReadFile für genau eine Datei scheitern — MemFS kennt nur
// "fehlt", nicht "unlesbar", und die Unterscheidung trägt den fail-closed Rand.
type readErrOneFS struct {
	driven.Filesystem
	bad string
}

func (f readErrOneFS) ReadFile(rel string) ([]byte, error) {
	if rel == f.bad {
		return nil, errors.New("unlesbar")
	}
	return f.Filesystem.ReadFile(rel)
}

// listErrSubFS lässt List für genau ein Verzeichnis scheitern.
type listErrSubFS struct {
	driven.Filesystem
	bad string
}

func (f listErrSubFS) List(rel string) ([]driven.DirEntry, error) {
	if rel == f.bad {
		return nil, errors.New("unlesbar")
	}
	return f.Filesystem.List(rel)
}

func filesOf(fs []model.Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.File)
	}
	return out
}

// Ohne recursive bleibt ein Unterverzeichnis unsichtbar, und
// in derselben Funktion die Umkehr: mit recursive meldet derselbe Baum die
// dünne Notiz unter wellenlos/.
func TestClosureRecursive_SiehtUnterverzeichnis(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md":           "# Slice\n\n" + richNote,
		closureDir + "/wellenlos/slice-002-b.md": kaputteNote,
	}
	cfg := closureCfg()
	if f := CheckPlanningClosure(coretest.NewMemFS(files), cfg); f != nil {
		t.Fatalf("ohne recursive ist wellenlos/ unsichtbar, got %+v", f)
	}
	cfg.Closure.Recursive = true
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	want := closureDir + "/wellenlos/slice-002-b.md"
	if len(f) != 1 || f[0].File != want || f[0].Reason != model.ReasonClosureNoteThin {
		t.Fatalf("mit recursive: genau closure-note-thin auf %s erwartet, got %+v", want, f)
	}
}

// Ohne recursive ist ein Verzeichnis, dessen Name den Filter trifft,
// Kandidat und meldet sich als unlesbar — mit recursive wird es betreten.
func TestClosureRecursive_VerzeichnisNameOhneSchalterIstKandidat(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md":          "# Slice\n\n" + richNote,
		closureDir + "/slice-900-dir.md/slice-x.md": "# Slice\n\n" + richNote,
	}
	cfg := closureCfg()
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].File != closureDir+"/slice-900-dir.md" {
		t.Fatalf("ohne recursive: das Verzeichnis ist Kandidat und unlesbar, got %+v", f)
	}
	cfg.Closure.Recursive = true
	if f := CheckPlanningClosure(coretest.NewMemFS(files), cfg); f != nil {
		t.Fatalf("mit recursive: das Verzeichnis wird betreten, sein Inhalt ist sauber, got %+v", f)
	}
}

// Ein unlesbares Unterverzeichnis ist unter recursive kein stilles Grün, und
// die Meldung nennt das Unterverzeichnis.
func TestClosureRecursive_UnlesbaresUnterverzeichnisFailClosed(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md":           "# Slice\n\n" + richNote,
		closureDir + "/wellenlos/slice-002-b.md": "# Slice\n\n" + richNote,
	}
	cfg := closureCfg()
	cfg.Closure.Recursive = true
	fsys := listErrSubFS{Filesystem: coretest.NewMemFS(files), bad: closureDir + "/wellenlos"}
	f := CheckPlanningClosure(fsys, cfg)
	if len(f) != 1 || f[0].Reason != model.ReasonClosureNoteMissing ||
		!strings.Contains(f[0].Message, closureDir+"/wellenlos") {
		t.Fatalf("unlesbares Unterverzeichnis ⇒ closure-note-missing mit seinem Pfad, got %+v", f)
	}
}

// skip-pattern nimmt den Stub aus, der Volltext daneben bleibt geprüft — und
// ohne das Muster meldet derselbe Baum den Stub.
func TestClosureSkipPattern_NimmtStubAus(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md":           "# Slice\n\n" + richNote,
		closureDir + "/wellenlos/slice-002-b.md": "# slice-002\n\n" + stubMarker,
		closureDir + "/wellenlos/slice-003-c.md": kaputteNote,
	}
	cfg := closureCfg()
	cfg.Closure.Recursive = true
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if got := strings.Join(filesOf(f), ","); !strings.Contains(got, "slice-002-b.md") {
		t.Fatalf("ohne skip-pattern meldet der Stub, got %v", got)
	}
	cfg.Closure.SkipPattern = stubPattern
	f = CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	want := closureDir + "/wellenlos/slice-003-c.md"
	if len(f) != 1 || f[0].File != want {
		t.Fatalf("mit skip-pattern: nur der dünne Volltext %s meldet, got %+v", want, f)
	}
}

// Nimmt skip-pattern alle Kandidaten, gilt die Nullmengen-Regel — die Meldung
// nennt die Ausnahme.
func TestClosureSkipPattern_LeereMengeFailClosed(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# slice-001\n\n" + stubMarker,
	}
	cfg := closureCfg()
	cfg.Closure.SkipPattern = stubPattern
	f := CheckPlanningClosure(coretest.NewMemFS(files), cfg)
	if len(f) != 1 || f[0].File != closureDir || f[0].Reason != model.ReasonClosureNoteMissing ||
		!strings.Contains(f[0].Message, "skip-pattern") {
		t.Fatalf("alles ausgenommen ⇒ closure-note-missing auf dem Verzeichnis mit Nennung der Ausnahme, got %+v", f)
	}
}

// Eine unlesbare Datei verschluckt die Ausnahme nicht: sie bleibt Kandidatin
// und meldet sich fail-closed.
func TestClosureSkipPattern_UnlesbareDateiBleibtKandidat(t *testing.T) {
	bad := closureDir + "/slice-002-b.md"
	files := map[string]string{
		closureDir + "/slice-001-a.md": "# Slice\n\n" + richNote,
		bad:                            "# slice-002\n\n" + stubMarker,
	}
	cfg := closureCfg()
	cfg.Closure.SkipPattern = stubPattern
	f := CheckPlanningClosure(readErrOneFS{Filesystem: coretest.NewMemFS(files), bad: bad}, cfg)
	if len(f) != 1 || f[0].File != bad || f[0].Reason != model.ReasonClosureNoteMissing {
		t.Fatalf("unlesbar ⇒ Kandidat, closure-note-missing auf %s, got %+v", bad, f)
	}
}

// Unter recursive bleiben die immer übersprungenen Verzeichnisse
// unbetreten — und ohne die Ausnahme meldete dieselbe Notiz darin.
func TestClosureRecursive_SkipDirsBleibenUnbetreten(t *testing.T) {
	files := map[string]string{
		closureDir + "/slice-001-a.md":              "# Slice\n\n" + richNote,
		closureDir + "/node_modules/slice-002-b.md": kaputteNote,
		closureDir + "/wellenlos/slice-003-c.md":    "# Slice\n\n" + richNote,
	}
	cfg := closureCfg()
	cfg.Closure.Recursive = true
	if f := CheckPlanningClosure(coretest.NewMemFS(files), cfg); f != nil {
		t.Fatalf("node_modules/ darf nicht betreten werden, got %+v", f)
	}
}

// Ohne die neuen Schlüssel bleibt die Meldung eines unlesbaren closure.dir
// byte-identisch — auch für ein ungereinigtes dir wie "x/".
func TestClosureDir_UngereinigtMeldungUnveraendert(t *testing.T) {
	cfg := closureCfg()
	cfg.Closure.Dir = "docs/fehlt/"
	f := CheckPlanningClosure(listErrFS{}, cfg)
	want := "Closure-Verzeichnis docs/fehlt/ fehlt oder ist unlesbar (fail-closed)"
	if len(f) != 1 || f[0].Message != want || f[0].File != "docs/fehlt/" {
		t.Fatalf("Meldung verändert: want %q, got %+v", want, f)
	}
}
