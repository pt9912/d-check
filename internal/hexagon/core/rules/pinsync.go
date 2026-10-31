package rules

// Modul pinsync (DC-FA-PSYNC-001): ein benannter Wert muss an allen
// konfigurierten Dateistellen gleich sein. Das Modul ist hermetisch -- es
// liest benannte Dateien ueber den Filesystem-Port und fuehrt nichts aus.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): die Quelldateien liegen ausserhalb
// der Scan-Menge. scan.roots, scan.ignore und die Skip-Liste schraenken sie
// nicht ein, und die Zusage der Scan-Menge gilt dort nicht; sie gilt fuer die
// ausdruecklich benannten Dateien. Gelesen wird Text ohne Format: ob ein
// Ausdruck die richtige Stelle meint, sagt der Ausdruck, nicht das Modul.
//
// FAIL-CLOSED: ein Eingabefehler ist ein error (Exit 2) und beendet den Lauf
// beim ersten Fehler -- zuerst die Konfiguration aller Pins, dann das Lesen --
// ohne dass ein Befund erscheint.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls pinsync (spec/spezifikation.md §4).
const (
	ReasonPinsyncMismatch  = "pinsync-mismatch"
	ReasonPinsyncNoMatch   = "pinsync-no-match"
	ReasonPinsyncAmbiguous = "pinsync-ambiguous"
)

// PinsyncResult traegt die Befunde und die Bezugsmenge des Laufs.
type PinsyncResult struct {
	Findings []model.Finding
	Pins     int
	Sources  int
}

// Note liefert die Bezugsmenge fuer die Zusammenfassung: sie sagt, was
// gelesen wurde, nicht, dass jeder Pin des Repos erfasst ist.
func (r PinsyncResult) Note() string {
	return fmt.Sprintf("pinsync: %d Pin(s) geprüft, %d Quelle(n) gelesen", r.Pins, r.Sources)
}

// pinsyncRead ist das Ergebnis einer gelesenen Quelle: die Zeilen, die
// treffen, und der Wert der ersten treffenden Zeile.
type pinsyncRead struct {
	src   model.PinsyncSource
	lines []int
	value string
}

// CheckPinsync ist das Regelmodul pinsync. Der error-Rueckgabewert traegt die
// Eingabefehler; der Aufrufer mappt ihn auf Exit 2.
func CheckPinsync(fsys driven.Filesystem, cfg model.PinsyncConfig) (PinsyncResult, error) {
	var res PinsyncResult
	if len(cfg.Pins) == 0 {
		return res, fmt.Errorf("pinsync: pinsync.pins ist leer, das Modul hätte nichts zu prüfen (DC-FA-PSYNC-001, fail-closed)")
	}
	if err := cfg.Validate(); err != nil {
		return res, fmt.Errorf("pinsync: %v (DC-FA-PSYNC-001, fail-closed)", err)
	}
	res.Pins = len(cfg.Pins)
	for _, pin := range cfg.SortedPins() {
		reads := make([]pinsyncRead, 0, len(pin.Sources))
		for _, src := range pin.Sources {
			rd, err := readPinsyncSource(fsys, src)
			if err != nil {
				return PinsyncResult{}, fmt.Errorf("pinsync: pin %q, quelle %s: %v (DC-FA-PSYNC-001, fail-closed)", pin.Name, src.File, err)
			}
			res.Sources++
			reads = append(reads, rd)
		}
		res.Findings = append(res.Findings, pinsyncFindings(pin.Name, reads)...)
	}
	return res, nil
}

// readPinsyncSource prueft den Pfad Segment fuer Segment (kein Symlink, am
// Ende eine Datei) und faengt dann die Treffer-Zeilen.
func readPinsyncSource(fsys driven.Filesystem, src model.PinsyncSource) (pinsyncRead, error) {
	if err := checkPinsyncPath(fsys, src.File); err != nil {
		return pinsyncRead{}, err
	}
	content, err := fsys.ReadFile(src.File)
	if err != nil {
		return pinsyncRead{}, fmt.Errorf("nicht lesbar: %v", err)
	}
	re, err := model.CompilePinsyncPattern(src.Pattern)
	if err != nil {
		return pinsyncRead{}, err
	}
	rd := pinsyncRead{src: src}
	for i, line := range pinsyncLines(content) {
		if v, ok := pinsyncCapture(re, line); ok {
			if len(rd.lines) == 0 {
				rd.value = v
			}
			rd.lines = append(rd.lines, i+1)
		}
	}
	return rd, nil
}

// checkPinsyncPath klassifiziert jedes Praefix des Pfades mit Kind (Lstat):
// fehlend, Symlink oder -- vor dem Ende -- kein Verzeichnis sowie am Ende ein
// Verzeichnis sind Eingabefehler.
func checkPinsyncPath(fsys driven.Filesystem, file string) error {
	segs := strings.Split(file, "/")
	for i := range segs {
		prefix := strings.Join(segs[:i+1], "/")
		kind, err := fsys.Kind(prefix)
		if err != nil {
			return fmt.Errorf("%q nicht prüfbar: %v", prefix, err)
		}
		switch {
		case kind == driven.KindMissing:
			return fmt.Errorf("%q fehlt", prefix)
		case kind == driven.KindSymlink:
			return fmt.Errorf("%q ist ein Symlink (nicht gelesen)", prefix)
		case i < len(segs)-1 && kind != driven.KindDir:
			return fmt.Errorf("%q ist kein Verzeichnis", prefix)
		case i == len(segs)-1 && kind != driven.KindFile:
			return fmt.Errorf("%q ist keine Datei", prefix)
		}
	}
	return nil
}

// pinsyncLines teilt am '\n', schneidet ein einzelnes '\r' am Zeilenende ab
// und kennt keine leere Schlusszeile hinter einem abschliessenden '\n'.
func pinsyncLines(content []byte) []string {
	if len(content) == 0 {
		return nil
	}
	text := strings.TrimSuffix(string(content), "\n")
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSuffix(l, "\r")
	}
	return lines
}

// pinsyncCapture liefert den Wert der Fanggruppe in der ersten Passung der
// Zeile; eine nicht beteiligte oder leere Gruppe ist kein Treffer.
func pinsyncCapture(re *regexp.Regexp, line string) (string, bool) {
	idx := re.FindStringSubmatchIndex(line)
	if len(idx) < 4 || idx[2] < 0 || idx[3] <= idx[2] {
		return "", false
	}
	return line[idx[2]:idx[3]], true
}

// pinsyncFindings bildet die Befunde eines Pins: erst je Quelle ohne oder mit
// mehreren Treffern, dann die Quellen, deren Wert vom Bezug abweicht.
func pinsyncFindings(name string, reads []pinsyncRead) []model.Finding {
	var out []model.Finding
	var exact []pinsyncRead
	for _, rd := range reads {
		switch len(rd.lines) {
		case 0:
			out = append(out, pinsyncFinding(name, rd.src.File, 1, ReasonPinsyncNoMatch,
				fmt.Sprintf("Pin %q: der Ausdruck %q trifft in %s keine Zeile mit nicht leerem Wert", name, rd.src.Pattern, rd.src.File)))
		case 1:
			exact = append(exact, rd)
		default:
			out = append(out, pinsyncFinding(name, rd.src.File, rd.lines[0], ReasonPinsyncAmbiguous,
				fmt.Sprintf("Pin %q: der Ausdruck %q trifft in %s mehrere Zeilen (%s)", name, rd.src.Pattern, rd.src.File, joinInts(rd.lines))))
		}
	}
	if len(exact) < 2 {
		return out
	}
	ref := exact[0]
	for _, rd := range exact[1:] {
		if rd.value != ref.value {
			out = append(out, pinsyncFinding(name, rd.src.File, rd.lines[0], ReasonPinsyncMismatch,
				fmt.Sprintf("Pin %q: %s hat %q, der Bezug %s hat %q (Werte: %s)", name, rd.src.File, rd.value, ref.src.File, ref.value, pinsyncValues(reads))))
		}
	}
	return out
}

// pinsyncValues nennt alle Quellen des Pins mit ihrem Wert; eine Quelle ohne
// genau einen Treffer steht als `?`.
func pinsyncValues(reads []pinsyncRead) string {
	parts := make([]string, 0, len(reads))
	for _, rd := range reads {
		v := "?"
		if len(rd.lines) == 1 {
			v = rd.value
		}
		parts = append(parts, rd.src.File+"="+v)
	}
	return strings.Join(parts, ", ")
}

func joinInts(ns []int) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ", ")
}

func pinsyncFinding(pin, file string, line int, reason, msg string) model.Finding {
	return model.Finding{File: file, Line: line, Rule: "pinsync", Target: pin, Reason: reason, Message: msg}
}
