package rules

// Tabelle des Plans (DC-FA-MET-001): genau eine Markdown-Tabelle außerhalb
// umzäunter Blöcke mit konfiguriertem Kopf; gezählt werden die Zeilen je
// Klasse. Zäune erkennt der Leser wie der Tabellen-Leser der RTM (Öffner-Test
// des Produkts nach Entfernen von Space und Tab vorn, geschlossen mit
// mindestens gleicher Länge), Zellen teilt der gemeinsame Zeilen-Teiler
// SplitPipeTableLine. Kopf und Datenzeilen beginnen an Spalte 0; Einrückung,
// Zitat-Zeichen und Listen machen eine Zeile zu keiner Tabellenzeile.
// Ein Tabellenfehler ist ein Eingabefehler des Plans.

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// readMetricsTable liefert die Zeilenzahl je Klasse der einen Tabelle des Plans.
func readMetricsTable(content string, t model.MetricsTable) (map[string]int, error) {
	lines := metricsOutsideFence(strings.Split(content, "\n"))
	starts := metricsTableStarts(lines, t.Header)
	switch len(starts) {
	case 0:
		return nil, fmt.Errorf("tabelle fehlt (kopf %q)", strings.Join(t.Header, " | "))
	case 1:
	default:
		return nil, fmt.Errorf("tabelle mehrdeutig: %d Tabellen mit dem Kopf %q", len(starts), strings.Join(t.Header, " | "))
	}
	return countMetricsRows(lines, starts[0], t)
}

// metricsLine ist eine Plan-Zeile mit der Angabe, ob sie außerhalb eines Zauns
// liegt (Zaun-Zeilen selbst zählen als innerhalb).
type metricsLine struct {
	text    string
	outside bool
}

// metricsOutsideFence markiert die Zeilen außerhalb umzäunter Blöcke.
func metricsOutsideFence(raw []string) []metricsLine {
	out := make([]metricsLine, len(raw))
	var char byte
	run := 0
	for i, r := range raw {
		text := strings.TrimSuffix(r, "\r")
		trimmed := TrimFenceIndent(text)
		switch {
		case run > 0:
			if FenceCloses(trimmed, char, run) {
				run = 0
			}
		case FenceToggle(trimmed):
			char, run = FenceRun(trimmed)
		default:
			out[i].outside = true
		}
		out[i].text = text
	}
	return out
}

// metricsTableStarts liefert die Indizes der Kopfzeilen, deren Zellen genau
// header sind (außerhalb von Zäunen).
func metricsTableStarts(lines []metricsLine, header []string) []int {
	var starts []int
	for i, l := range lines {
		if !l.outside || !atColumnZero(l.text) {
			continue
		}
		cells, ok := SplitPipeTableLine(l.text)
		if ok && equalStrings(cells, header) {
			starts = append(starts, i)
		}
	}
	return starts
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// countMetricsRows prüft Trennzeile und Datenzeilen ab dem Kopf bei index und
// zählt die Zeilen je Klasse.
func countMetricsRows(lines []metricsLine, head int, t model.MetricsTable) (map[string]int, error) {
	if head+1 >= len(lines) || !lines[head+1].outside || !isMetricsSeparator(lines[head+1].text, len(t.Header)) {
		return nil, fmt.Errorf("trennzeile nach dem tabellenkopf fehlt")
	}
	keyIdx, classIdx := indexOfString(t.Header, t.KeyColumn), indexOfString(t.Header, t.ClassColumn)
	counts := map[string]int{}
	keys := map[string]bool{}
	n := 0
	leading := strings.HasPrefix(lines[head].text, "|")
	for i := head + 2; i < len(lines) && lines[i].outside; i++ {
		if !atColumnZero(lines[i].text) || leading && !strings.HasPrefix(lines[i].text, "|") {
			break
		}
		cells, ok := SplitPipeTableLine(lines[i].text)
		if !ok {
			break
		}
		n++
		if len(cells) != len(t.Header) {
			return nil, fmt.Errorf("zeile %d der tabelle: %d zellen statt %d", n, len(cells), len(t.Header))
		}
		if err := checkMetricsRow(cells[keyIdx], cells[classIdx], n, keys, t); err != nil {
			return nil, err
		}
		counts[cells[classIdx]]++
	}
	if n == 0 {
		return nil, fmt.Errorf("tabelle ohne zeilen")
	}
	return counts, nil
}

// checkMetricsRow prüft Schlüssel (nicht leer, eindeutig) und Klasse (eine der
// konfigurierten).
func checkMetricsRow(key, class string, n int, keys map[string]bool, t model.MetricsTable) error {
	if key == "" {
		return fmt.Errorf("zeile %d der tabelle: %s leer", n, t.KeyColumn)
	}
	if keys[key] {
		return fmt.Errorf("zeile %d der tabelle: %s %q doppelt", n, t.KeyColumn, key)
	}
	keys[key] = true
	if _, ok := t.Classes[class]; !ok {
		return fmt.Errorf("zeile %s: klasse %q unbekannt", key, class)
	}
	return nil
}

func indexOfString(list []string, s string) int {
	for i, v := range list {
		if v == s {
			return i
		}
	}
	return -1
}

// isMetricsSeparator meldet, ob die Zeile eine Trennzeile mit want Zellen ist:
// jede Zelle aus '-' mit optionalem ':' an einem oder beiden Enden.
func isMetricsSeparator(line string, want int) bool {
	cells, ok := SplitPipeTableLine(line)
	if !ok || len(cells) != want {
		return false
	}
	for _, c := range cells {
		core := strings.TrimSuffix(strings.TrimPrefix(c, ":"), ":")
		if core == "" || strings.Trim(core, "-") != "" {
			return false
		}
	}
	return true
}

// metricsEstimate bildet die Schätzung: Summe aus Klassenfaktor mal Zeilen und
// Feldfaktor mal Wert des Block-Feldes, mit Überlauf-Prüfung.
func metricsEstimate(counts map[string]int, b *metricsBlock, cfg model.MetricsConfig) (int, error) {
	sum := 0
	add := func(factor, n int) error {
		if factor != 0 && n > math.MaxInt/factor {
			return fmt.Errorf("schätzung: %d · %d läuft über", factor, n)
		}
		p := factor * n
		if sum > math.MaxInt-p {
			return fmt.Errorf("schätzung läuft über")
		}
		sum += p
		return nil
	}
	for _, class := range sortedKeys(counts) {
		n := counts[class]
		if err := add(cfg.Table.Classes[class], n); err != nil {
			return 0, err
		}
	}
	for _, name := range sortedKeys(cfg.EstimateInputs) {
		factor := cfg.EstimateInputs[name]
		if err := add(factor, b.inputs[name]); err != nil {
			return 0, err
		}
	}
	return sum, nil
}

// sortedKeys liefert die Schlüssel einer Abbildung in fester Reihenfolge
// (DC-QA-02).
func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// atColumnZero meldet, ob die Zeile nicht mit Leerraum oder einem Zitat-Zeichen
// beginnt.
func atColumnZero(text string) bool {
	return text != "" && text[0] != ' ' && text[0] != '\t' && text[0] != '>'
}
