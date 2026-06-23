package rules

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// diagramFixture: arch.md mit einer ARC-Definitions-Tabelle (Fließtext)
// und einem mermaid-Diagramm; ARC-01/ARC-02 sind definiert, ARC-99 nicht.
func diagramFixture() map[string]string {
	return map[string]string{
		"spec/arch.md": "## Komponenten\n" +
			"| ARC-01 Foo | x |\n" +
			"| ARC-02 Bar | y |\n" +
			"\n" +
			"```mermaid\n" +
			"flowchart TB\n" +
			"  A[\"ARC-01\"] --> B[\"ARC-02\"]\n" +
			"  C[\"ARC-99\"]\n" +
			"```\n",
	}
}

func diagramsCfg(definedIn string, fences ...string) model.Config {
	return model.Config{Diagrams: model.DiagramsConfig{
		Fences:   fences,
		Patterns: []model.DiagramPattern{{Regex: regexp.MustCompile(`ARC-\d{2}`), DefinedIn: definedIn}},
	}}
}

// DC-FA-DIAG-001 Boundary: ARC-99 im mermaid-Fence ohne Tabellen-Definition
// → ein diagram-id-undefined (Zeile im Fence); die in der Tabelle
// definierten ARC-01/ARC-02 bleiben befundfrei.
func TestDiagramsUndefinedID(t *testing.T) {
	m := coretest.NewMemFS(diagramFixture())
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s %s", f.Line, f.Rule, f.Target, f.Reason))
	}
	want := []string{"8 diagrams ARC-99 diagram-id-undefined"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
}

// DC-FA-DIAG-001 Happy: alle Kennungen im Diagramm in defined-in definiert
// → kein Befund.
func TestDiagramsAllDefined(t *testing.T) {
	fix := diagramFixture()
	fix["spec/arch.md"] = "| ARC-99 Baz | z |\n" + fix["spec/arch.md"] // ARC-99 nun definiert
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("erwartete 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-DIAG-001 Negative: ohne aktives Modul kein Befund (byte-identisch
// zum Lauf ohne diagrams).
func TestDiagramsNichtAktiv(t *testing.T) {
	m := coretest.NewMemFS(diagramFixture())
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Modul nicht aktiv → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-DIAG-001: nicht-gelistete Fences (bash) bleiben opak — ARC-99 dort
// ist kein Befund; nur der mermaid-Fence wird geöffnet.
func TestDiagramsNichtGelisteterFence(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/arch.md": "| ARC-01 Foo | x |\n" +
			"```bash\necho ARC-99\n```\n" +
			"```mermaid\nA[\"ARC-01\"]\n```\n",
	})
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("bash-Fence darf nicht geprüft werden, got %v", res.Findings)
	}
}

// DC-FA-DIAG-001: eine eigene fences-Liste öffnet eine andere Sprache.
func TestDiagramsCustomFence(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/arch.md": "| ARC-01 |\n```dot\nARC-99\n```\n",
	})
	res, err := Run(m, nil, diagramsCfg("spec/arch.md", "dot"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "ARC-99" {
		t.Fatalf("custom fence dot nicht geöffnet, got %v", res.Findings)
	}
}

// DC-FA-DIAG-001: fehlendes bzw. escapendes defined-in → Exit 2 (Run-Fehler).
func TestDiagramsDefinedInValidierung(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{"spec/arch.md": "```mermaid\nARC-01\n```\n"})
	if _, err := Run(m, nil, diagramsCfg("spec/fehlt.md"), []string{"diagrams"}); err == nil {
		t.Fatal("fehlendes defined-in muss Fehler liefern")
	}
	if _, err := Run(m, nil, diagramsCfg("../escape.md"), []string{"diagrams"}); err == nil {
		t.Fatal("escapendes defined-in muss Fehler liefern")
	}
}
