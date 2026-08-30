package rules

// Chronologie-Monotonie (DC-FA-STRUCT-001 §Chronologie-Monotonie,
// §DC-FA-STRUCT-001.a Schritt 6, ADR-0057): die ERSTE von drei
// structure-Bedingungen, die die ROHEN Abschnitts-Zeilen lesen — die
// Bereinigung aus Schritt 5 leert
// Inline-Code, und reale Schlüsselspalten (Release-Register) stehen genau
// dort. Die Lexik-Fragen bleiben geteilt: ob eine Zeile Tabellenzeile,
// Kopf-/Trennzeile oder Zelle ist, beantwortet markdown.go — dieselbe Antwort
// wie bei targets und planning.waves.

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// tableDateRE/tableVersionRE sind die geschlossene Typ-Menge der
// Schlüsselspalte (ADR-0057 Entscheidung 3): ISO-Datum und Punkt-Version mit
// optionalem v-Präfix und mindestens zwei numerischen Segmenten. Getypt wird
// der ERSTE Treffer in der rohen Zelle — so überleben Inline-Code-Backticks
// und ein HTML-Anker neben dem Schlüssel.
var (
	tableDateRE    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	tableVersionRE = regexp.MustCompile(`v?\d+(?:\.\d+)+`)
)

// tableKey ist der typisierte Schlüssel einer Zelle. tok trägt den
// Fund-Token für Meldungen.
type tableKey struct {
	isDate bool
	date   string // ISO-Format — der String-Vergleich ist ordnungstreu
	segs   []int  // Versions-Segmente, numerisch
	tok    string
}

func (k tableKey) typeName() string {
	if k.isDate {
		return "Datum"
	}
	return "Version"
}

// typeTableKey typisiert eine rohe Zelle: der früheste Treffer beider Muster
// gewinnt (ein Datum enthält keine Punkte, eine Version keine Bindestriche —
// derselbe Index ist nicht konstruierbar; Gleichstand fiele ans Datum).
func typeTableKey(cell string) (tableKey, bool) {
	d := tableDateRE.FindStringIndex(cell)
	v := tableVersionRE.FindStringIndex(cell)
	switch {
	case d == nil && v == nil:
		return tableKey{}, false
	case v == nil || (d != nil && d[0] <= v[0]):
		tok := cell[d[0]:d[1]]
		return tableKey{isDate: true, date: tok, tok: tok}, true
	default:
		tok := cell[v[0]:v[1]]
		segs := versionSegments(tok)
		if segs == nil {
			return tableKey{}, false
		}
		return tableKey{segs: segs, tok: tok}, true
	}
}

// versionSegments zerlegt eine Punkt-Version in numerische Segmente; das
// optionale v-Präfix trägt keine Ordnung. Ein Segment jenseits des
// int-Bereichs ist ERREICHBAR (`\d+` ist unbegrenzt) und liefert nil —
// der Aufrufer behandelt das als untypisierbar, damit ein Zahlen-Monstrum
// nicht still als kleinstmögliche Version vergleicht — Befund statt
// Übersprung.
func versionSegments(tok string) []int {
	parts := strings.Split(strings.TrimPrefix(tok, "v"), ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		out = append(out, n)
	}
	return out
}

// compareTableKeys vergleicht zwei Schlüssel DESSELBEN Typs: Datum als
// ISO-String, Version segmentweise numerisch — die kürzere Segmentfolge ist
// bei gleichem Präfix kleiner (1.9 < 1.9.1, und 1.9 < 1.10).
func compareTableKeys(a, b tableKey) int {
	if a.isDate {
		return strings.Compare(a.date, b.date)
	}
	for i := 0; i < len(a.segs) && i < len(b.segs); i++ {
		if a.segs[i] != b.segs[i] {
			if a.segs[i] < b.segs[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case len(a.segs) < len(b.segs):
		return -1
	case len(a.segs) > len(b.segs):
		return 1
	}
	return 0
}

// structureTableOrder prüft die Chronologie-Monotonie eines Abschnitts:
// Datenzeilen je ZUSAMMENHÄNGENDER Tabelle, nicht-strikt monoton in Richtung
// table-order. Null Datenzeilen ⇒ Leerlauf-Befund an der Überschrift — die
// Bedingung zu setzen IST die Behauptung, dass hier eine chronologische
// Tabelle steht.
func structureTableOrder(
	r model.StructureRule, file string, lines []string, prose map[int]bool, headingNo, level int,
) []model.Finding {
	end := SectionEnd(lines, headingNo, level)
	if end == 0 {
		end = len(lines) + 1
	}
	var out []model.Finding
	dataRows := 0
	var prev *tableKey
	prevLine := 0
	inTable := false
	for i := headingNo; i < end-1; i++ {
		if !tableRowLine(lines, prose, i) {
			// Tabellen-Ende: der Vergleich startet in der nächsten Tabelle neu.
			inTable, prev = false, nil
			continue
		}
		if !inTable {
			inTable, prev = true, nil
		}
		if tableHeaderOrSeparator(lines, prose, i) {
			continue // Kopf- und Trennzeile deklarieren keine Daten
		}
		dataRows++
		var f *model.Finding
		f, prev, prevLine = structureTableRow(r, file, lines[i], i+1, prev, prevLine)
		if f != nil {
			out = append(out, *f)
		}
	}
	if dataRows == 0 {
		out = append(out, structureFinding(r, file, headingNo, model.ReasonSectionUnordered,
			"Abschnitt trägt keine Tabellen-Datenzeile — die Chronologie-Zusage läuft leer"))
	}
	return out
}

// structureTableRow verarbeitet eine Datenzeile: Zelle adressieren, Schlüssel
// typisieren, gegen den Vorgänger vergleichen. Eine untypisierbare Zelle ist
// ein Befund, kein Übersprung — und setzt den Vergleich zurück (nur
// benachbarte typisierbare Paare werden verglichen).
func structureTableRow(
	r model.StructureRule, file, line string, lineNo int, prev *tableKey, prevLine int,
) (*model.Finding, *tableKey, int) {
	col := r.Table.EffectiveOrderColumn()
	cells := tableCells(line)
	if len(cells) < col {
		f := structureFinding(r, file, lineNo, model.ReasonSectionCellUntyped,
			"Zeile hat "+strconv.Itoa(len(cells))+" Zelle(n) — Schlüsselspalte "+
				strconv.Itoa(col)+" fehlt")
		return &f, nil, 0
	}
	key, ok := typeTableKey(cells[col-1])
	if !ok {
		f := structureFinding(r, file, lineNo, model.ReasonSectionCellUntyped,
			"Schlüsselzelle (Spalte "+strconv.Itoa(col)+") trägt kein Datums-/Versions-Token")
		return &f, nil, 0
	}
	if prev == nil {
		return nil, &key, lineNo
	}
	if prev.isDate != key.isDate {
		// Anker-Reset wie bei den beiden anderen Fehlerfaellen: die
		// Misch-Zelle meldet sich selbst, eine gesunde
		// Folge-Zeile dahinter meldet nicht.
		f := structureFinding(r, file, lineNo, model.ReasonSectionCellUntyped,
			"Typ-Mischung in der Schlüsselspalte: "+key.typeName()+" ("+key.tok+
				") neben "+prev.typeName()+" ("+prev.tok+", Zeile "+strconv.Itoa(prevLine)+")")
		return &f, nil, 0
	}
	c := compareTableKeys(key, *prev)
	if (r.Table.Order == "asc" && c < 0) || (r.Table.Order == "desc" && c > 0) {
		f := structureFinding(r, file, lineNo, model.ReasonSectionUnordered,
			"Schlüssel "+key.tok+" bricht die "+r.Table.Order+"-Ordnung gegenüber "+
				prev.tok+" (Zeile "+strconv.Itoa(prevLine)+")")
		return &f, &key, lineNo
	}
	return nil, &key, lineNo
}
