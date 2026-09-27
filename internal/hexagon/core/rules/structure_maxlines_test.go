package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// abschnittMitZeilen baut einen Abschnitt mit genau n gezaehlten Zeilen,
// gefolgt von einer WEITEREN Ueberschrift: SectionEnd begrenzt den Abschnitt
// dann auf eine konkrete Zeile statt auf "Dateiende", und body enthaelt keine
// zusaetzliche Phantom-Zeile aus dem schliessenden Zeilenumbruch der Datei
// (derselbe split-bedingte Leerstring, den JEDE letzte Ueberschrift einer
// Datei sonst als zusaetzliche "Zeile" trueg — eine Eigenschaft der
// Zeilennummerierung dieses Pakets, nicht dieser Bedingung).
func abschnittMitZeilen(n int) string {
	body := "## Harte Regeln\n"
	for i := 0; i < n; i++ {
		body += "Zeile.\n"
	}
	return body + "## Naechster Abschnitt\n\nx.\n"
}

// Grenzwert-Symmetrie: N ist gruen, N+1 ist rot -- wie bei jeder anderen
// zaehlenden Bedingung des Moduls (max-tasks, cell-max-chars).
func TestStructureMaxLines_Grenzwert(t *testing.T) {
	r := model.StructureRule{Files: "docs/*.md", Section: "## Harte Regeln", MaxLines: ptr(5)}
	fsGruen := coretest.NewMemFS(map[string]string{"docs/a.md": abschnittMitZeilen(5)})
	if f := CheckStructure(fsGruen, []model.StructureRule{r}); f != nil {
		t.Fatalf("5 Zeilen bei max-lines 5 ⇒ befundfrei, got %+v", f)
	}
	fsRot := coretest.NewMemFS(map[string]string{"docs/a.md": abschnittMitZeilen(6)})
	f := CheckStructure(fsRot, []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionLinesExceeded {
		t.Fatalf("6 Zeilen bei max-lines 5 muss section-lines-exceeded melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "6 Zeilen") || !strings.Contains(f[0].Message, "5") {
		t.Errorf("Meldung nennt Ist- und Soll-Zahl: %q", f[0].Message)
	}
}

// Fenced-Code zaehlt NICHT mit -- SectionProse entfernt seine Zeilen
// vollstaendig (ADR-0089, dieselbe Grundmenge wie min-sentences). Derselbe
// Abschnitt mit einem Codeblock bleibt darum unter der Schwelle, obwohl der
// ROHE Abschnitt mehr Zeilen hat, als max-lines erlaubt.
func TestStructureMaxLines_FencedCodeZaehltNicht(t *testing.T) {
	body := "## Harte Regeln\n" +
		"Zeile.\nZeile.\n" +
		"```\n" + strings.Repeat("code-zeile\n", 10) + "```\n"
	r := model.StructureRule{Files: "docs/*.md", Section: "## Harte Regeln", MaxLines: ptr(5)}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": body})
	if f := CheckStructure(fs, []model.StructureRule{r}); f != nil {
		t.Fatalf("zwei Prosa-Zeilen + zehn Fenced-Code-Zeilen bei max-lines 5 "+
			"muessen befundfrei bleiben (Fenced-Code zaehlt nicht mit), got %+v", f)
	}
}

// Ein fehlender Abschnitt bleibt section-missing -- max-lines fuehrt KEINE
// eigene Form fuer diesen Fall ein.
func TestStructureMaxLines_FehlenderAbschnittBleibtSectionMissing(t *testing.T) {
	r := model.StructureRule{Files: "docs/*.md", Section: "## Harte Regeln", MaxLines: ptr(5)}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": "# T\n\n## Anderer Abschnitt\n\nx.\n"})
	f := CheckStructure(fs, []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("fehlender Abschnitt muss weiter section-missing melden, got %+v", f)
	}
}

// sections: one mit mehreren Treffern bleibt section-ambiguous -- keine
// Messung, wie bei jeder anderen Bedingung.
func TestStructureMaxLines_MehrdeutigBleibtSectionAmbiguous(t *testing.T) {
	r := model.StructureRule{Files: "docs/*.md", Section: "## Harte Regeln", MaxLines: ptr(1)}
	body := "# T\n\n## Harte Regeln\n\nZeile.\n\n## Harte Regeln\n\nZeile.\n"
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": body})
	f := CheckStructure(fs, []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionAmbiguous {
		t.Fatalf("mehrdeutiger Abschnitt muss section-ambiguous melden, keine Messung, got %+v", f)
	}
}

// hint gewinnt gegen die modul-eigene Meldung -- dieselbe Form wie bei jeder
// anderen Bedingung (structureFinding statt structureRawFinding).
func TestStructureMaxLines_HintGewinnt(t *testing.T) {
	r := model.StructureRule{
		Files: "docs/*.md", Section: "## Harte Regeln", MaxLines: ptr(1),
		Hint: "AGENTS.md waechst -- Abschnitt kuerzen oder auslagern",
	}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": abschnittMitZeilen(2)})
	f := CheckStructure(fs, []model.StructureRule{r})
	if len(f) != 1 || f[0].Message != "AGENTS.md waechst -- Abschnitt kuerzen oder auslagern" {
		t.Fatalf("hint muss die Meldung ersetzen, got %+v", f)
	}
}

// Ohne max-lines ist der Befundsatz byte-identisch -- dieselbe Zusage wie
// bei jedem anderen abwesenden Schluessel (DC-QA-02).
func TestStructureMaxLines_AbwesendByteIdentisch(t *testing.T) {
	r := model.StructureRule{Files: "docs/*.md", Section: "## Harte Regeln", NonEmpty: true}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": abschnittMitZeilen(999)})
	if f := CheckStructure(fs, []model.StructureRule{r}); f != nil {
		t.Fatalf("ohne max-lines darf die Zeilenzahl keinen Befund erzeugen, got %+v", f)
	}
}
