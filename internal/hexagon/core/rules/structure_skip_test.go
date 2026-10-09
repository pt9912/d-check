package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

func skipRule() model.StructureRule {
	return model.StructureRule{
		Files: "done/**/slice-*.md", SectionPattern: `^## [0-9]+\. Definition of Done`, NonEmpty: true,
	}
}

// skip-pattern nimmt den Stub aus der Kandidatenmenge, der Volltext daneben
// bleibt geprüft — und ohne das Muster meldet derselbe Baum den Stub.
func TestStructureSkipPattern_NimmtStubAus(t *testing.T) {
	files := map[string]string{
		"done/wellenlos/slice-001-a.md": "# slice-001\n\n" + stubMarker,
		"done/wellenlos/slice-002-b.md": "# Slice\n\n## 2. Definition of Done\n\n- [x] fertig\n",
		"done/wellenlos/slice-003-c.md": "# Slice\n\nohne DoD\n",
	}
	r := skipRule()
	f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if got := strings.Join(filesOf(f), ","); !strings.Contains(got, "slice-001-a.md") {
		t.Fatalf("VORZUSTAND: ohne skip-pattern meldet der Stub, got %v", got)
	}
	r.SkipPattern = stubPattern
	f = CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(f) != 1 || f[0].File != "done/wellenlos/slice-003-c.md" || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("mit skip-pattern: nur der Volltext ohne DoD meldet, got %+v", f)
	}
}

// Nimmt skip-pattern alle Dateien, gilt die Nullmengen-Härte — die Meldung
// nennt die Ausnahme.
func TestStructureSkipPattern_LeereMengeFailClosed(t *testing.T) {
	files := map[string]string{"done/wellenlos/slice-001-a.md": "# slice-001\n\n" + stubMarker}
	r := skipRule()
	r.SkipPattern = stubPattern
	f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(f) != 1 || f[0].File != r.Files || !strings.Contains(f[0].Message, "skip-pattern") {
		t.Fatalf("alles ausgenommen ⇒ section-missing auf dem Glob mit Nennung der Ausnahme, got %+v", f)
	}
}

// Ohne skip-pattern bleibt die Nullmengen-Meldung byte-identisch.
func TestStructureSkipPattern_OhneSchluesselMeldungUnveraendert(t *testing.T) {
	r := skipRule()
	f := CheckStructure(coretest.NewMemFS(map[string]string{}), []model.StructureRule{r})
	want := "Regel trifft keine Datei (auch nach Abzug von exempt-paths) — das Gate liefe leer"
	if len(f) != 1 || f[0].Message != want {
		t.Fatalf("Meldung ohne skip-pattern verändert, got %+v", f)
	}
}

// Eine unlesbare Datei verschluckt die Ausnahme nicht: sie bleibt Kandidatin
// und meldet sich fail-closed.
func TestStructureSkipPattern_UnlesbareDateiBleibtKandidat(t *testing.T) {
	bad := "done/wellenlos/slice-001-a.md"
	files := map[string]string{bad: "# slice-001\n\n" + stubMarker}
	r := skipRule()
	r.SkipPattern = stubPattern
	f := CheckStructure(readErrOneFS{Filesystem: coretest.NewMemFS(files), bad: bad}, []model.StructureRule{r})
	if len(f) != 1 || f[0].File != bad {
		t.Fatalf("unlesbar ⇒ Kandidatin, Befund auf %s, got %+v", bad, f)
	}
}
