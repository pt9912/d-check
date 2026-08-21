package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// tableOrderRule baut die Minimal-Regel der Chronologie-Bedingung.
func tableOrderRule(order string, col *int) model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", Section: "## Historie",
		TableOrder: order, TableColumn: col,
	}
}

func tableOrderFindings(t *testing.T, body, order string, col *int) []model.Finding {
	t.Helper()
	files := map[string]string{"docs/a.md": "# T\n\n## Historie\n\n" + body}
	return CheckStructure(coretest.NewMemFS(files),
		[]model.StructureRule{tableOrderRule(order, col)})
}

// Happy Path: absteigend sortierte Datums-Spalte, gleiche Schlüssel in Folge
// erlaubt (nicht-strikt — mehrere Releases an einem Tag sind der Normalfall).
func TestTableOrderHappyDesc(t *testing.T) {
	body := "| Datum | Was |\n|---|---|\n" +
		"| 2026-08-16 | c |\n| 2026-08-16 | b |\n| 2026-08-10 | a |\n"
	if f := tableOrderFindings(t, body, "desc", nil); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}

// Bruch: die brechende Datenzeile wird gemeldet — mit IHRER Zeilennummer.
func TestTableOrderBruch(t *testing.T) {
	body := "| Datum | Was |\n|---|---|\n" +
		"| 2026-08-10 | a |\n| 2026-08-16 | b |\n"
	f := tableOrderFindings(t, body, "desc", nil)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionUnordered {
		t.Fatalf("erwartet section-unordered, got %+v", f)
	}
	if f[0].Line != 8 {
		t.Errorf("Line = %d, want 8 (die brechende Zeile)", f[0].Line)
	}
}

// Typ-Pflicht: segmentweise numerisch ist 0.10.0 groesser als 0.9.0 und 1.10
// groesser als 1.9 — ein zeichenweiser Vergleich meldete beide Tabellen rot.
func TestTableOrderTypPflichtVersionen(t *testing.T) {
	desc := "| Version |\n|---|\n| 0.10.0 |\n| 0.9.0 |\n| 0.8.2 |\n"
	if f := tableOrderFindings(t, desc, "desc", nil); f != nil {
		t.Fatalf("desc-Versionen: erwartet befundfrei, got %+v", f)
	}
	asc := "| Handbuch |\n|---|\n| 1.9 |\n| 1.10 |\n| 1.10.1 |\n"
	if f := tableOrderFindings(t, asc, "asc", nil); f != nil {
		t.Fatalf("asc-Versionen: erwartet befundfrei, got %+v", f)
	}
}

// Roh-Lesung: die Schluesselzelle steht in Inline-Code, die aktuelle Zeile
// traegt zusaetzlich einen HTML-Anker (version.md-Form). Auf dem bereinigten
// Abschnitts-Text waere jede Zelle leer.
func TestTableOrderRoheZellen(t *testing.T) {
	body := "| Version | Datum |\n|---|---|\n" +
		"| `v0.60.0` <a id=\"v0.60.0\"></a> | 2026-08-16 |\n" +
		"| `v0.59.0` | 2026-08-16 |\n"
	if f := tableOrderFindings(t, body, "desc", nil); f != nil {
		t.Fatalf("erwartet befundfrei (rohe Zellen typisierbar), got %+v", f)
	}
	// Gegenprobe: derselbe Inhalt falsch herum meldet — die Zellen werden
	// also wirklich gelesen, nicht still uebersprungen.
	kaputt := "| Version | Datum |\n|---|---|\n" +
		"| `v0.59.0` | 2026-08-16 |\n" +
		"| `v0.60.0` <a id=\"v0.60.0\"></a> | 2026-08-16 |\n"
	f := tableOrderFindings(t, kaputt, "desc", nil)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionUnordered {
		t.Fatalf("Gegenprobe: erwartet section-unordered, got %+v", f)
	}
}

// Untypisierbare Zelle, fehlende Spalte und Typ-Mischung sind Befunde, kein
// stiller Uebersprung; benachbarte typisierbare Paare werden weiter geprueft.
func TestTableOrderUntypisierbar(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []string
	}{
		{"kein Token", "| Datum |\n|---|\n| 2026-08-16 |\n| kaputt |\n| 2026-08-10 |\n",
			[]string{model.ReasonSectionCellUntyped}},
		{"zu wenige Zellen", "| A | B |\n|---|---|\n| x | 2026-08-16 |\n| nur-eine |\n",
			[]string{model.ReasonSectionCellUntyped}},
		{"Typ-Mischung", "| Schluessel |\n|---|\n| 2026-08-16 |\n| v0.60.0 |\n",
			[]string{model.ReasonSectionCellUntyped}},
	}
	zwei := 2
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var col *int
			if c.name == "zu wenige Zellen" {
				col = &zwei
			}
			f := tableOrderFindings(t, c.body, "desc", col)
			var got []string
			for _, x := range f {
				got = append(got, x.Reason)
			}
			if len(got) != len(c.want) || got[0] != c.want[0] {
				t.Fatalf("erwartet %v, got %+v", c.want, f)
			}
		})
	}
}

// F-1-Kaskadenfall (Review): der Anker wird auch nach einer Typ-Mischung
// zurueckgesetzt — die Misch-Zelle meldet sich selbst, die GESUNDE
// Folge-Zeile dahinter (die den Spalten-Typ traegt) meldet NICHT. Genau eine
// Lesart, gepinnt.
func TestTableOrderTypMischungKaskade(t *testing.T) {
	body := "| Schluessel |\n|---|\n" +
		"| 2026-08-16 |\n| v0.60.0 |\n| 2026-08-10 |\n| 2026-08-01 |\n"
	f := tableOrderFindings(t, body, "desc", nil)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellUntyped {
		t.Fatalf("genau EIN cell-untyped an der Misch-Zelle erwartet, got %+v", f)
	}
	if f[0].Line != 8 {
		t.Errorf("Line = %d, want 8 (die Misch-Zelle, nicht die gesunde Folge-Zeile)", f[0].Line)
	}
}

// F-2 (Review): ein Versions-Segment jenseits des int-Bereichs ist
// untypisierbar — Befund statt stillem Vergleich als kleinstmoegliche Version.
func TestTableOrderVersionsSegmentUeberlauf(t *testing.T) {
	body := "| V |\n|---|\n" +
		"| 0.10.0 |\n| 12345678901234567890.1 |\n| 0.8.0 |\n"
	f := tableOrderFindings(t, body, "desc", nil)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellUntyped {
		t.Fatalf("Ueberlauf-Segment ⇒ genau ein cell-untyped erwartet, got %+v", f)
	}
	if f[0].Line != 8 {
		t.Errorf("Line = %d, want 8 (die Ueberlauf-Zelle)", f[0].Line)
	}
}

// Nach einer untypisierbaren Zelle setzt der Vergleich beim naechsten
// typisierbaren NACHBAR-Paar wieder auf — die kaputte Zelle meldet sich
// selbst, macht aber nicht die restliche Tabelle unpruefbar.
func TestTableOrderSetztNachUntypisierbarNeuAuf(t *testing.T) {
	body := "| Datum |\n|---|\n" +
		"| 2026-08-16 |\n| kaputt |\n| 2026-08-20 |\n| 2026-08-10 |\n"
	f := tableOrderFindings(t, body, "desc", nil)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionCellUntyped {
		t.Fatalf("genau ein cell-untyped erwartet (2026-08-20 ist kein Nachbar von 2026-08-16), got %+v", f)
	}
}

// Kopf- und Trennzeile deklarieren keine Daten — auch wenn die Kopfzelle wie
// ein Schluessel aussieht oder ein Wort traegt.
func TestTableOrderKopfUndTrennzeile(t *testing.T) {
	body := "| Datum |\n|---|\n| 2026-08-16 |\n| 2026-08-10 |\n"
	if f := tableOrderFindings(t, body, "desc", nil); f != nil {
		t.Fatalf("Kopfzeile 'Datum' darf nicht als untypisierbare Datenzeile melden, got %+v", f)
	}
}

// Schluesselspalte per table-column adressieren (Roadmap-Abschlussregister:
// Spalte 2 traegt das Datum, Spalte 1 einen Namen).
func TestTableOrderSpaltenAdresse(t *testing.T) {
	body := "| Welle | Abschluss |\n|---|---|\n" +
		"| welle-76-x | 2026-08-16 |\n| welle-75-y | 2026-08-10 |\n"
	col := 2
	if f := tableOrderFindings(t, body, "desc", &col); f != nil {
		t.Fatalf("erwartet befundfrei ueber Spalte 2, got %+v", f)
	}
}

// Leerlauf: kein Datenbestand — ob gar keine Tabellenzeile oder nur eine im
// Fence — meldet an der Abschnitts-Ueberschrift. Die Bedingung zu setzen IST
// die Behauptung, dass hier eine chronologische Tabelle steht.
func TestTableOrderLeerlauf(t *testing.T) {
	fence := "```"
	cases := map[string]string{
		"ohne Tabelle": "Nur Prosa.\n",
		"Tabelle im Fence": fence + "\n| Datum |\n|---|\n| 2026-08-16 |\n| 2026-08-20 |\n" +
			fence + "\n",
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			f := tableOrderFindings(t, body, "desc", nil)
			if len(f) != 1 || f[0].Reason != model.ReasonSectionUnordered {
				t.Fatalf("erwartet Leerlauf-section-unordered, got %+v", f)
			}
			if f[0].Line != 3 {
				t.Errorf("Line = %d, want 3 (die Abschnitts-Ueberschrift)", f[0].Line)
			}
			if !strings.Contains(f[0].Message, "keine Tabellen-Datenzeile") {
				t.Errorf("Leerlauf-Meldung fehlt, got %q", f[0].Message)
			}
		})
	}
}

// Zwei zusammenhaengende Tabellen werden je fuer sich geprueft: der Vergleich
// startet an der zweiten neu (benannte Grenze aus ADR-0057 — zwei einzeln
// sortierte, gegenlaeufig liegende Tabellen sind KEIN Befund).
func TestTableOrderJeTabelle(t *testing.T) {
	body := "| Datum |\n|---|\n| 2026-08-16 |\n| 2026-08-10 |\n" +
		"\nProsa dazwischen.\n\n" +
		"| Datum |\n|---|\n| 2026-09-01 |\n| 2026-08-30 |\n"
	if f := tableOrderFindings(t, body, "desc", nil); f != nil {
		t.Fatalf("zwei einzeln sortierte Tabellen ⇒ befundfrei (benannte Grenze), got %+v", f)
	}
}

// Ohne table-order ist der Befundsatz unveraendert — die Bedingung ist strikt
// opt-in je Regel (Byte-Identitaets-Haelfte der Chronologie).
func TestTableOrderOptIn(t *testing.T) {
	body := "| Datum |\n|---|\n| 2026-08-10 |\n| 2026-08-16 |\n"
	if f := tableOrderFindings(t, body, "", nil); f != nil {
		t.Fatalf("ohne table-order erwartet befundfrei, got %+v", f)
	}
}

// Versionsvergleich am Rand: kuerzere Segmentfolge ist bei gleichem Praefix
// kleiner (1.9 < 1.9.1), das v-Praefix traegt keine Ordnung.
func TestTableOrderVersionsRaender(t *testing.T) {
	asc := "| V |\n|---|\n| 1.9 |\n| 1.9.1 |\n| v1.10 |\n"
	if f := tableOrderFindings(t, asc, "asc", nil); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", f)
	}
}
