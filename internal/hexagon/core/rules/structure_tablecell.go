package rules

// Zellenlaenge (DC-FA-STRUCT-001 §Zellenlaenge, §DC-FA-STRUCT-001.a Schritt 6,
// ADR-0069): die neunte structure-Bedingung. Sie adressiert ihre Spalte ueber
// den KOPFZEILEN-Namen, nicht ueber eine Position — eine eingefuegte Spalte
// verschoebe sonst still die Messung.
//
// Seit ADR-0070 traegt eine Regel eine LISTE solcher Spalten; dieser Lauf gilt
// genau EINER davon. Die Trennung der Befunde leistet deshalb nicht mehr die
// Regel-Identitaet, sondern ColumnTarget je Befund: zwei zu lange Zellen
// derselben Zeile fielen sonst unter die Deduplikation zusammen.
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
// ist das ein Befund an der Ueberschrift — die Spalte zu nennen IST die
// Behauptung, dass es sie gibt.
func structureCellMax(
	r model.StructureRule, col model.TableColumnRule, file string,
	lines []string, prose map[int]bool, headingNo, level int,
) []model.Finding {
	end := SectionEnd(lines, headingNo, level)
	if end == 0 {
		end = len(lines) + 1
	}
	var out []model.Finding
	scan := cellMaxScan{rule: r, col: col, file: file, pos: -1}
	for i := headingNo; i < end-1; i++ {
		if !tableRowLine(lines, prose, i) {
			scan.inTable, scan.pos = false, -1
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
		out = append(out, scan.finding(headingNo, model.ReasonSectionColumnMissing,
			"keine Tabelle des Abschnitts traegt eine Kopfzelle "+strconv.Quote(col.Name)))
	}
	return out
}

// cellMaxScan traegt den Lauf-Zustand ueber die Tabellen eines Abschnitts:
// pos ist die 0-basierte Position der benannten Spalte in der laufenden
// Tabelle (-1 = nicht gebunden, Tabelle irrelevant), bound merkt sich, ob
// IRGENDEINE Tabelle sie gebunden hat.
type cellMaxScan struct {
	rule    model.StructureRule
	col     model.TableColumnRule
	file    string
	pos     int
	inTable bool
	bound   bool
}

// finding bindet jeden Befund dieses Laufs an die SPALTE, nicht nur an die
// Regel — sonst traefen zwei Spalten derselben Zeile dieselbe Befund-Adresse.
func (s *cellMaxScan) finding(line int, reason, msg string) model.Finding {
	return model.Finding{
		File: s.file, Line: line, Rule: "structure",
		Target: s.rule.ColumnTarget(s.col.Name), Reason: reason, Message: s.rule.MessageFor(msg),
	}
}

// bind adressiert die Spalte in der Kopfzeile der beginnenden Tabelle. Eine
// Tabelle ohne Kopf-/Trennzeile bindet nichts; ein doppelter Name ist ein
// Befund, kein Erst-Treffer — welche der beiden Spalten gemeint ist, sagt die
// Konfiguration nicht.
func (s *cellMaxScan) bind(lines []string, prose map[int]bool, i int) []model.Finding {
	s.pos = -1
	if i+1 >= len(lines) || !tableRowLine(lines, prose, i+1) || !tableSeparatorRow(lines[i+1]) {
		return nil
	}
	var hits []int
	for n, cell := range tableCells(lines[i]) {
		if strings.TrimSpace(cell) == s.col.Name {
			hits = append(hits, n)
		}
	}
	if len(hits) > 1 {
		return []model.Finding{s.finding(i+1, model.ReasonSectionColumnMissing,
			"Kopfzelle "+strconv.Quote(s.col.Name)+" kommt "+strconv.Itoa(len(hits))+
				"-mal vor — die Spalte ist nicht adressierbar")}
	}
	if len(hits) == 1 {
		s.pos, s.bound = hits[0], true
	}
	return nil
}

// measure prueft eine Datenzeile der laufenden Tabelle. Gezaehlt werden
// ZEICHEN, nicht Bytes: die Schwelle beschreibt einen Text, und ein Umlaut ist
// ein Zeichen.
func (s *cellMaxScan) measure(line string, lineNo int) *model.Finding {
	if s.pos < 0 {
		return nil // Tabelle ohne die benannte Spalte
	}
	cells := tableCells(line)
	if len(cells) <= s.pos {
		f := s.finding(lineNo, model.ReasonSectionColumnMissing,
			"Zeile hat "+strconv.Itoa(len(cells))+" Zelle(n) — Spalte "+
				strconv.Quote(s.col.Name)+" (Position "+strconv.Itoa(s.pos+1)+") fehlt")
		return &f
	}
	got := utf8.RuneCountInString(cells[s.pos])
	if n := s.col.CellMinChars; n != nil && got < *n {
		f := s.finding(lineNo, model.ReasonSectionCellUndersized,
			"Zelle der Spalte "+strconv.Quote(s.col.Name)+" hat "+strconv.Itoa(got)+
				" Zeichen, verlangt sind "+strconv.Itoa(*n))
		return &f
	}
	if n := s.col.CellMaxChars; n != nil && got > *n {
		f := s.finding(lineNo, model.ReasonSectionCellOversized,
			"Zelle der Spalte "+strconv.Quote(s.col.Name)+" hat "+strconv.Itoa(got)+
				" Zeichen, erlaubt sind "+strconv.Itoa(*n))
		return &f
	}
	return nil
}
