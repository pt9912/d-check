package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// cellMaxRule baut die Minimal-Regel der Zellenlaengen-Bedingung: EINE Spalte
// in der Tabellen-Klammer.
func cellMaxRule(col string, obergrenze int) model.StructureRule {
	return cellMaxRuleN(model.TableColumnRule{Name: col, CellMaxChars: &obergrenze})
}

// cellMaxRuleN baut die Regel ueber BELIEBIG viele Spalten desselben
// Abschnitts — die Form, die ADR-0070 an die Stelle mehrerer Regeln mit
// wiederholtem Selektor setzt.
func cellMaxRuleN(cols ...model.TableColumnRule) model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", Section: "## Index",
		Table: model.TableRule{Columns: cols},
	}
}

// cellMaxFindings laesst die Bedingung ueber einen Abschnitts-Rumpf laufen.
// Der Rumpf beginnt bei Zeile 5.
func cellMaxFindings(t *testing.T, body, col string, obergrenze int) []model.Finding {
	t.Helper()
	files := map[string]string{"docs/a.md": "# T\n\n## Index\n\n" + body}
	return CheckStructure(coretest.NewMemFS(files),
		[]model.StructureRule{cellMaxRule(col, obergrenze)})
}

// Happy Path: jede Zelle der benannten Spalte bleibt unter der Schwelle. Die
// NACHBAR-Spalte ist laenger — die Bedingung ist spalten-gebunden, nicht
// zeilen-pauschal.
func TestCellMaxHappy(t *testing.T) {
	body := "| ID | Titel | Bezug |\n|---|---|---|\n" +
		"| A | kurz | " + strings.Repeat("x", 400) + " |\n" +
		"| B | auch kurz | " + strings.Repeat("y", 400) + " |\n"
	if f := cellMaxFindings(t, body, "Titel", 20); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}

// Zu lange Zelle: gemeldet wird IHRE Zeile, nicht die des Abschnitts — dort
// ist die Reparatur. Die Meldung nennt Ist und Soll.
func TestCellMaxZuLang(t *testing.T) {
	body := "| ID | Titel |\n|---|---|\n" +
		"| A | kurz |\n" +
		"| B | " + strings.Repeat("x", 30) + " |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellOversized {
		t.Fatalf("erwartet genau ein section-cell-oversized, got %+v", f)
	}
	if f[0].Line != 8 {
		t.Errorf("Line = %d, want 8 (die zu lange Zeile)", f[0].Line)
	}
	if !strings.Contains(f[0].Message, "30 Zeichen") || !strings.Contains(f[0].Message, "erlaubt sind 20") {
		t.Errorf("Meldung nennt Ist/Soll nicht: %q", f[0].Message)
	}
}

// Grenzfall: genau auf der Schwelle ist erlaubt, eines darueber nicht — die
// Schwelle ist die groesste zulaessige Laenge, nicht die kleinste verbotene.
func TestCellMaxGrenzfall(t *testing.T) {
	auf := "| Titel |\n|---|\n| " + strings.Repeat("x", 20) + " |\n"
	if f := cellMaxFindings(t, auf, "Titel", 20); f != nil {
		t.Fatalf("genau auf der Schwelle: erwartet befundfrei, got %+v", f)
	}
	drueber := "| Titel |\n|---|\n| " + strings.Repeat("x", 21) + " |\n"
	if f := cellMaxFindings(t, drueber, "Titel", 20); len(f) != 1 {
		t.Fatalf("eines darueber: erwartet einen Befund, got %+v", f)
	}
}

// Zeichen, nicht Bytes: 20 Umlaute sind 40 Byte und 20 Zeichen. Eine
// Byte-Zaehlung meldete diese Zelle rot.
func TestCellMaxZeichenNichtBytes(t *testing.T) {
	body := "| Titel |\n|---|\n| " + strings.Repeat("ä", 20) + " |\n"
	if f := cellMaxFindings(t, body, "Titel", 20); f != nil {
		t.Fatalf("erwartet befundfrei (20 Zeichen), got %+v", f)
	}
}

// Escapte Pipe: `\|` ist ein Zeichen der Zelle, kein Zelltrenner. Ein
// naiver Split zerteilte die Zelle und meldete die zu lange Zelle NICHT —
// genau das Loch, das ein Zeichenketten-Waechter offenlaesst.
func TestCellMaxEscaptePipe(t *testing.T) {
	lang := strings.Repeat("x", 15) + ` \| ` + strings.Repeat("y", 15)
	body := "| ID | Titel |\n|---|---|\n| A | " + lang + " |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellOversized {
		t.Fatalf("erwartet section-cell-oversized, got %+v", f)
	}
	// 15 + " | " + 15 = 33: die escapte Pipe zaehlt als EIN Zeichen.
	if !strings.Contains(f[0].Message, "33 Zeichen") {
		t.Errorf("escapte Pipe falsch gezaehlt: %q", f[0].Message)
	}
}

// Fehlende Spalte: kein Kopf traegt den Namen ⇒ Befund an der Ueberschrift.
// Die Bedingung zu setzen IST die Behauptung, dass es die Spalte gibt.
func TestCellMaxSpalteFehlt(t *testing.T) {
	body := "| ID | Name |\n|---|---|\n| A | b |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionColumnMissing {
		t.Fatalf("erwartet section-column-missing, got %+v", f)
	}
	if f[0].Line != 3 {
		t.Errorf("Line = %d, want 3 (die Abschnitts-Ueberschrift)", f[0].Line)
	}
}

// Doppelter Spaltenname: nicht adressierbar — welche der beiden gemeint ist,
// sagt die Konfiguration nicht. Befund auf der Kopfzeile.
func TestCellMaxSpalteDoppelt(t *testing.T) {
	body := "| Titel | Titel |\n|---|---|\n| a | b |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 2 {
		t.Fatalf("erwartet Kopf-Befund + Leerlauf-Befund, got %+v", f)
	}
	if f[0].Reason != model.ReasonSectionColumnMissing || f[0].Line != 5 {
		t.Errorf("erwartet section-column-missing auf der Kopfzeile 5, got %+v", f[0])
	}
}

// Zu kurze Datenzeile: die Spalte ist dort nicht vorhanden — Befund statt
// stillem Uebersprung, sonst schaltete eine kaputte Zeile die Messung ab.
func TestCellMaxZeileZuKurz(t *testing.T) {
	body := "| ID | Titel |\n|---|---|\n| A |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionColumnMissing {
		t.Fatalf("erwartet section-column-missing, got %+v", f)
	}
	if f[0].Line != 7 {
		t.Errorf("Line = %d, want 7 (die zu kurze Datenzeile)", f[0].Line)
	}
}

// Leere Pruefmenge: der Abschnitt traegt gar keine Tabelle ⇒ Leerlauf-Befund.
// Eine Bedingung, die nichts misst, meldet nicht Erfolg.
func TestCellMaxLeerlauf(t *testing.T) {
	f := cellMaxFindings(t, "Nur Prosa, keine Tabelle.\n", "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionColumnMissing {
		t.Fatalf("erwartet section-column-missing (Leerlauf), got %+v", f)
	}
}

// Fence-Treue: eine Tabelle IN einem Fenced-Block deklariert nichts — sie
// bindet die Spalte nicht und ihre Zellen werden nicht gemessen.
func TestCellMaxFenceTreu(t *testing.T) {
	body := "```\n| Titel |\n|---|\n| " + strings.Repeat("x", 40) + " |\n```\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionColumnMissing {
		t.Fatalf("erwartet nur den Leerlauf-Befund, got %+v", f)
	}
}

// Zweite Tabelle im Abschnitt: die Spalte wird je Tabelle neu gebunden, und
// eine Tabelle OHNE die Spalte schaltet die Messung der anderen nicht ab.
func TestCellMaxZweiTabellen(t *testing.T) {
	body := "| Name |\n|---|\n| " + strings.Repeat("x", 40) + " |\n\n" +
		"| Titel |\n|---|\n| " + strings.Repeat("y", 40) + " |\n"
	f := cellMaxFindings(t, body, "Titel", 20)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellOversized {
		t.Fatalf("erwartet genau den Befund der zweiten Tabelle, got %+v", f)
	}
	if f[0].Line != 11 {
		t.Errorf("Line = %d, want 11 (Datenzeile der zweiten Tabelle)", f[0].Line)
	}
}

// Leere Zelle: unter einer Obergrenze ALLEIN passiert sie gruen — null
// Zeichen liegen unter jeder Schwelle. Erst die Untergrenze faengt sie, und
// zwar mit einem ANDEREN Grund-Code: die Reparatur ist ausfuellen, nicht
// kuerzen.
func TestCellMaxLeereZelle(t *testing.T) {
	body := "| ID | Titel |\n|---|---|\n| A |  |\n"
	if f := cellMaxFindings(t, body, "Titel", 20); f != nil {
		t.Fatalf("nur Obergrenze: erwartet befundfrei, got %+v", f)
	}
	untergrenze := 1
	r := cellMaxRule("Titel", 20)
	r.Table.Columns[0].CellMinChars = &untergrenze
	files := map[string]string{"docs/a.md": "# T\n\n## Index\n\n" + body}
	f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r})
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellUndersized {
		t.Fatalf("mit Untergrenze: erwartet section-cell-undersized, got %+v", f)
	}
	if f[0].Line != 7 {
		t.Errorf("Line = %d, want 7 (die leere Zelle)", f[0].Line)
	}
}

// Zwei Spalten, eine Zeile, EINE Regel: das Befund-Ziel traegt die Spalte,
// sonst faenden die beiden Befunde dieselbe Adresse (Datei, Zeile, Regel,
// Ziel, Grund) und die Deduplikation verloere einen.
func TestCellMaxZweiSpaltenEineZeile(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n\n## Index\n\n" +
		"| Titel | Status |\n|---|---|\n| " + strings.Repeat("x", 30) + " | " +
		strings.Repeat("y", 30) + " |\n"}
	zwanzig := 20
	f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{cellMaxRuleN(
		model.TableColumnRule{Name: "Titel", CellMaxChars: &zwanzig},
		model.TableColumnRule{Name: "Status", CellMaxChars: &zwanzig},
	)})
	if len(f) != 2 {
		t.Fatalf("erwartet zwei Befunde (eine Zeile, zwei Spalten), got %+v", f)
	}
	if f[0].Target == f[1].Target {
		t.Errorf("Ziele identisch (%q) — die Deduplikation faellt sie zusammen", f[0].Target)
	}
}

// Die Zellen-Spalten stehen seit ADR-0070 NICHT mehr in der Regel-Identitaet:
// sie sind eine Liste INNERHALB einer Regel und koennen keine zwei Regeln mehr
// kollidieren lassen. Getrennt wird je BEFUND ueber ColumnTarget — und dessen
// Form ist die des frueheren targets, damit sich fuer keine Bestandsregel ein
// Befund-Ziel aendert.
func TestStructureIdentitaetOhneSpalte(t *testing.T) {
	r := cellMaxRuleN(model.TableColumnRule{Name: "Titel"})
	r.Files, r.Section = "docs/*.md", "## H"
	if got, want := r.Identity(), "docs/*.md :: ## H"; got != want {
		t.Errorf("Identity() = %q, want %q", got, want)
	}
	if got, want := r.ColumnTarget("Titel"), "docs/*.md :: ## H :: Spalte Titel"; got != want {
		t.Errorf("ColumnTarget() = %q, want %q", got, want)
	}
}

// Die CHRONOLOGIE-Spalte bleibt in der Identitaet: zwei Chronologie-Zusagen
// ueber denselben Abschnitt sind zwei Regeln und brauchen zwei Identitaeten.
func TestStructureIdentitaetMitOrderColumn(t *testing.T) {
	zwei := 2
	r := model.StructureRule{Files: "docs/*.md", Section: "## H",
		Table: model.TableRule{Order: "desc", OrderColumn: &zwei}}
	if got, want := r.Identity(), "docs/*.md :: ## H :: Spalte 2"; got != want {
		t.Errorf("Identity() = %q, want %q", got, want)
	}
	r.Table.OrderColumn = nil
	if got, want := r.Identity(), "docs/*.md :: ## H"; got != want {
		t.Errorf("Identity() ohne explizite Spalte = %q, want %q", got, want)
	}
}

// Modul-aus: ohne cell-max-column ist der Befundsatz byte-identisch
// (DC-QA-02) — die Bedingung liest dann keine Zeile.
func TestCellMaxOhneSchluesselInert(t *testing.T) {
	files := map[string]string{"docs/a.md": "# T\n\n## Index\n\n| Titel |\n|---|\n| " +
		strings.Repeat("x", 400) + " |\n"}
	r := model.StructureRule{Files: "docs/*.md", Section: "## Index"}
	if f := CheckStructure(coretest.NewMemFS(files), []model.StructureRule{r}); f != nil {
		t.Fatalf("erwartet befundfrei ohne die Bedingung, got %+v", f)
	}
}
