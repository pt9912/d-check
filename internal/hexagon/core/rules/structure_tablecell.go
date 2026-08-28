package rules

// Zellenlaenge (DC-FA-STRUCT-001 §Zellenlaenge, §DC-FA-STRUCT-001.a Schritt 6,
// ADR-0069): die neunte structure-Bedingung. Sie adressiert ihre Spalte ueber
// den KOPFZEILEN-Namen, nicht ueber eine Position — eine eingefuegte Spalte
// verschoebe sonst still die Messung.
//
// Wie die Chronologie-Bedingung liest sie die ROHEN Abschnitts-Zeilen: die
// Bereinigung leert Inline-Code, und eine reale Zelle traegt ihn.
//
// Grenze: gemessen wird die Zelle, wie sie dasteht — Markdown-Syntax
// eingeschlossen. Eine Zelle aus einem einzigen langen Link ist lang, auch
// wenn ihr sichtbarer Text kurz ist.

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// structureCellMax prueft die Zellenlaenge in der benannten Spalte jeder
// zusammenhaengenden Tabelle des Abschnitts. Bindet keine Tabelle die Spalte,
// ist das ein Befund an der Ueberschrift — die Bedingung zu setzen IST die
// Behauptung, dass es diese Spalte gibt.
func structureCellMax(
	r model.StructureRule, file string, lines []string, prose map[int]bool, headingNo, level int,
) []model.Finding {
	end := SectionEnd(lines, headingNo, level)
	if end == 0 {
		end = len(lines) + 1
	}
	var out []model.Finding
	scan := cellMaxScan{rule: r, file: file, col: -1}
	for i := headingNo; i < end-1; i++ {
		if !tableRowLine(lines, prose, i) {
			scan.inTable, scan.col = false, -1
			continue
		}
		if !scan.inTable {
			scan.inTable = true
			out = append(out, scan.bind(lines, prose, i)...)
		}
		if tableHeaderOrSeparator(lines, prose, i) {
			continue // Kopf- und Trennzeile tragen keine Daten
		}
		if f := scan.measure(lines[i], i+1); f != nil {
			out = append(out, *f)
		}
	}
	if !scan.bound {
		out = append(out, structureFinding(r, file, headingNo, model.ReasonSectionColumnMissing,
			"keine Tabelle des Abschnitts traegt eine Kopfzelle "+strconv.Quote(r.CellMaxColumn)))
	}
	return out
}

// cellMaxScan traegt den Lauf-Zustand ueber die Tabellen eines Abschnitts:
// col ist die 0-basierte Position der benannten Spalte in der laufenden
// Tabelle (-1 = nicht gebunden, Tabelle irrelevant), bound merkt sich, ob
// IRGENDEINE Tabelle sie gebunden hat.
type cellMaxScan struct {
	rule    model.StructureRule
	file    string
	col     int
	inTable bool
	bound   bool
}

// bind adressiert die Spalte in der Kopfzeile der beginnenden Tabelle. Eine
// Tabelle ohne Kopf-/Trennzeile bindet nichts; ein doppelter Name ist ein
// Befund, kein Erst-Treffer — welche der beiden Spalten gemeint ist, sagt die
// Konfiguration nicht.
func (s *cellMaxScan) bind(lines []string, prose map[int]bool, i int) []model.Finding {
	s.col = -1
	if i+1 >= len(lines) || !tableRowLine(lines, prose, i+1) || !tableSeparatorRow(lines[i+1]) {
		return nil
	}
	var hits []int
	for n, cell := range tableCells(lines[i]) {
		if strings.TrimSpace(cell) == s.rule.CellMaxColumn {
			hits = append(hits, n)
		}
	}
	if len(hits) > 1 {
		return []model.Finding{structureFinding(s.rule, s.file, i+1, model.ReasonSectionColumnMissing,
			"Kopfzelle "+strconv.Quote(s.rule.CellMaxColumn)+" kommt "+strconv.Itoa(len(hits))+
				"-mal vor — die Spalte ist nicht adressierbar")}
	}
	if len(hits) == 1 {
		s.col, s.bound = hits[0], true
	}
	return nil
}

// measure prueft eine Datenzeile der laufenden Tabelle. Gezaehlt werden
// ZEICHEN, nicht Bytes: die Schwelle beschreibt einen Text, und ein Umlaut ist
// ein Zeichen.
func (s *cellMaxScan) measure(line string, lineNo int) *model.Finding {
	if s.col < 0 {
		return nil // Tabelle ohne die benannte Spalte
	}
	cells := tableCells(line)
	if len(cells) <= s.col {
		f := structureFinding(s.rule, s.file, lineNo, model.ReasonSectionColumnMissing,
			"Zeile hat "+strconv.Itoa(len(cells))+" Zelle(n) — Spalte "+
				strconv.Quote(s.rule.CellMaxColumn)+" (Position "+strconv.Itoa(s.col+1)+") fehlt")
		return &f
	}
	got := utf8.RuneCountInString(cells[s.col])
	if n := s.rule.CellMinChars; n != nil && got < *n {
		f := structureFinding(s.rule, s.file, lineNo, model.ReasonSectionCellUndersized,
			"Zelle der Spalte "+strconv.Quote(s.rule.CellMaxColumn)+" hat "+strconv.Itoa(got)+
				" Zeichen, verlangt sind "+strconv.Itoa(*n))
		return &f
	}
	if n := s.rule.CellMaxChars; n != nil && got > *n {
		f := structureFinding(s.rule, s.file, lineNo, model.ReasonSectionCellOversized,
			"Zelle der Spalte "+strconv.Quote(s.rule.CellMaxColumn)+" hat "+strconv.Itoa(got)+
				" Zeichen, erlaubt sind "+strconv.Itoa(*n))
		return &f
	}
	return nil
}
