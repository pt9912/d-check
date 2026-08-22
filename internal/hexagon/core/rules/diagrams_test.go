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

// slice-076 (Review R-F-2): der geteilte FenceToggle wirkt auch auf
// diagramFenceLines — die dritte Aufrufstelle. Eine mermaid-Öffner-Zeile mit
// Backtick in der Infozeile ist nach CommonMark kein Fence-Öffner; ihr Inhalt
// bleibt Prosa und wird NICHT als Diagramm gelesen (0 Befunde). Ohne die Regel
// an dieser Stelle öffnete die Zeile einen mermaid-Fence, ARC-99 würde gelesen
// und — weil die defined-in-Quelle es nicht kennt — als diagram-id-undefined
// gemeldet. Die defined-in-Quelle ist bewusst eine **separate** Datei: läge das
// Diagramm in derselben Datei, „rettete" der proseLines-basierte Definitions-Scan
// (nicht mutiert) das ARC-99 als Prosa-Definition und maskierte die Mutation.
func TestDiagramsInfozeileMitBacktickKeinFence(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/diag.md": "```mermaid `x`\nC[\"ARC-99\"]\n```\n",
		"spec/defs.md": "ARC-01 ist definiert.\n",
	})
	res, err := Run(m, nil, diagramsCfg("spec/defs.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Infozeile mit Backtick fälschlich als mermaid-Fence gelesen: %+v", res.Findings)
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

// DC-FA-DIAG-001 (R1-impl F-4): Muster-Präzedenz — ein überlappendes Vorkommen
// gehört dem früher deklarierten Muster (wie ids). ARC-01 gehört dem ersten
// Muster (defined-in a.md, definiert) und wird NICHT dem zweiten zugeordnet
// (defined-in b.md, leer) — sonst gäbe es einen Befund.
func TestDiagramsMusterPraezedenz(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/a.md": "| ARC-01 |\n```mermaid\nARC-01\n```\n",
		"spec/b.md": "leer\n",
	})
	cfg := model.Config{Diagrams: model.DiagramsConfig{
		Patterns: []model.DiagramPattern{
			{Regex: regexp.MustCompile(`ARC-\d{2}`), DefinedIn: "spec/a.md"},
			{Regex: regexp.MustCompile(`ARC-\d+`), DefinedIn: "spec/b.md"},
		},
	}}
	res, err := Run(m, nil, cfg, []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Präzedenz: ARC-01 gehört dem ersten (definierten) Muster → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-DIAG-001 (R1-impl F-1): die Token-Grenze trägt die Nutzer-Regex (exakt
// wie ids); ARC-\d{2} matcht ARC-99 auch innerhalb eines größeren Worts. Pinnt
// das vertragskonforme Verhalten (§DC-FA-DIAG-001.a: „Token derselben Regex") —
// ein Bruch (still hinzugefügte Wortgrenzen) bliebe sonst grün.
func TestDiagramsTokenGrenzeRegexDelegiert(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"spec/a.md": "keine Definition\n```mermaid\nXARC-99Y\n```\n",
	})
	res, err := Run(m, nil, diagramsCfg("spec/a.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "ARC-99" {
		t.Fatalf("Token-Grenze ist regex-delegiert: erwartet 1× ARC-99, got %v", res.Findings)
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

// DC-FA-DIAG-001 Ventil (Datei): eine Datei in exempt-paths wird gar nicht
// geprueft — datei-weit, unabhaengig vom Scan-Bereich.
func TestDiagramsVentilDatei(t *testing.T) {
	fix := diagramFixture()
	fix["docs/report.md"] = "```mermaid\nflowchart TB\n  X[\"ARC-77\"]\n```\n"
	cfg := diagramsCfg("spec/arch.md")
	cfg.Diagrams.ExemptPaths = []string{"docs/**"}
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if f.File == "docs/report.md" {
			t.Fatalf("exempt-paths muss die Datei datei-weit ausnehmen, got %+v", res.Findings)
		}
	}
}

// DC-FA-DIAG-001 Ventil (Zeile): der Marker nimmt SEINE Zeile aus — eine
// Kennung auf einer anderen Zeile desselben Blocks meldet weiter. In Mermaid
// steht er als %%-Kommentar; das Modul sucht das Token, nicht den Kommentar.
func TestDiagramsVentilZeile(t *testing.T) {
	fix := diagramFixture()
	fix["spec/arch.md"] = "| ARC-01 Foo | x |\n\n```mermaid\n" +
		"  C[\"ARC-98\"] %% d-check:ignore\n" +
		"  D[\"ARC-99\"]\n```\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "ARC-99" {
		t.Fatalf("nur die markierte Zeile ist frei, got %+v", res.Findings)
	}
}

// DC-FA-DIAG-001 Ventil (ganzer Block): der Marker auf der OEFFNUNGSZEILE
// nimmt den kompletten Fence aus — genau der Fall, der welle-80 zum Scoping
// zwang (ein Beispiel-Diagramm mit mehreren erfundenen Kennungen).
func TestDiagramsVentilBlockUeberOeffnungszeile(t *testing.T) {
	fix := diagramFixture()
	fix["spec/arch.md"] = "| ARC-01 Foo | x |\n\n```mermaid <!-- d-check:ignore -->\n" +
		"  C[\"ARC-97\"]\n  D[\"ARC-98\"]\n  E[\"ARC-99\"]\n```\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Marker auf der Oeffnungszeile nimmt den ganzen Block aus, got %+v", res.Findings)
	}
}

// DC-FA-DIAG-001: der Marker auf der Oeffnungszeile wirkt nur fuer SEINEN
// Block — ein zweites Diagramm derselben Datei bleibt geprueft.
func TestDiagramsVentilBlockNurEigenerFence(t *testing.T) {
	fix := diagramFixture()
	fix["spec/arch.md"] = "| ARC-01 Foo | x |\n\n```mermaid <!-- d-check:ignore -->\n" +
		"  C[\"ARC-97\"]\n```\n\n```mermaid\n  D[\"ARC-96\"]\n```\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "ARC-96" {
		t.Fatalf("nur der markierte Block ist frei, got %+v", res.Findings)
	}
}

// DC-FA-DIAG-001 / DC-QA-02: ohne beide Ventile bleibt der Befundsatz, was er
// war — die Erweiterung ist opt-in.
func TestDiagramsOhneVentileUnveraendert(t *testing.T) {
	m := coretest.NewMemFS(diagramFixture())
	res, err := Run(m, nil, diagramsCfg("spec/arch.md"), []string{"diagrams"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "ARC-99" {
		t.Fatalf("ohne Ventile unveraendert ein Befund, got %+v", res.Findings)
	}
}
