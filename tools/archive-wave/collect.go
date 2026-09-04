// Sammelt die zu einer Welle gehoerenden Vorgaenge: den flachen Welle-Plan,
// die Slices mit passendem **Welle:**-Feld, und ihre Review-Reports.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var welleFieldRE = regexp.MustCompile(`(?m)^\*\*Welle:\*\*\s*(.*)$`)
var welleFieldStartRE = regexp.MustCompile(`(?m)^\*\*Welle:\*\*[ \t]*`)
var welleIDInFieldRE = regexp.MustCompile(`\bwelle-(\d+)\b`)
var sliceIDInNameRE = regexp.MustCompile(`slice-(\d+)`)

// ReadWelleField liest den rohen Wert des **Welle:**-Feldes einer
// Slice-Datei (z. B. "welle-87" oder ein mehrzeiliger wellenloser
// Begruendungs-Absatz), unveraendert fuer die Uebernahme in den Stub. Liest
// bis zur ersten LEERZEILE -- die Haus-Stil-Form dieses Repos schreibt das
// Feld haeufig als mehrsaetziger, umgebrochener Absatz (gemessen: 44 von 45
// wellenlosen Slices im ersten Anwendungslauf von slice-197), nicht als
// einzeilige Kennung. Ein Ein-Zeilen-Capture schnitt den Absatz mitten im
// Satz ab. Leerer String, wenn das Feld fehlt.
func ReadWelleField(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s lesen: %w", path, err)
	}
	loc := welleFieldStartRE.FindStringIndex(string(b))
	if loc == nil {
		return "", nil
	}
	rest := string(b)[loc[1]:]
	var lines []string
	for _, ln := range strings.Split(rest, "\n") {
		if strings.TrimSpace(ln) == "" {
			break
		}
		lines = append(lines, ln)
	}
	return strings.TrimSpace(strings.Join(lines, "\n")), nil
}

// SliceIDFromPath liest die slice-<NNN>-Kennung aus einem Dateinamen.
func SliceIDFromPath(path string) string {
	m := sliceIDInNameRE.FindStringSubmatch(filepath.Base(path))
	if m == nil {
		return ""
	}
	return "slice-" + m[1]
}

// FindWellePlan findet die flache Welle-Plan-Datei in
// <root>/docs/plan/planning/done/, die nicht auf -results.md endet. Genau
// eine Datei ist der Regelfall; **keine** ist gueltig fuer eine Welle von vor
// der Plan-Datei-Konvention (nur eine retroaktive `-results.md` existiert,
// slice-191) -- leerer String ohne Fehler signalisiert das. Mehrdeutig
// (>1 Kandidat) bleibt ein Fehler, kein Uebersprung.
func FindWellePlan(root, welleID string) (string, error) {
	dir := filepath.Join(root, "docs/plan/planning/done")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("done-Verzeichnis lesen: %w", err)
	}
	prefix := welleID + "-"
	var found []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ".md") {
			continue
		}
		if strings.HasSuffix(name, "-results.md") {
			continue
		}
		found = append(found, name)
	}
	switch len(found) {
	case 0:
		return "", nil
	case 1:
		return filepath.Join(dir, found[0]), nil
	default:
		sort.Strings(found)
		return "", fmt.Errorf("mehrdeutig: %d Kandidaten fuer %s: %v", len(found), welleID, found)
	}
}

// CollectSlices findet alle Slice-Dateien in <root>/docs/plan/planning/done/,
// deren **Welle:**-Feld die Kennung welleID nennt -- als eigenstaendige
// welle-<N>-Ziffernfolge im Feldwert, unabhaengig von Link- oder Freitextform.
func CollectSlices(root, welleID string) ([]string, error) {
	dir := filepath.Join(root, "docs/plan/planning/done")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("done-Verzeichnis lesen: %w", err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "slice-") || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s lesen: %w", path, err)
		}
		field := welleFieldRE.FindString(string(b))
		if field == "" {
			continue
		}
		m := welleIDInFieldRE.FindStringSubmatch(field)
		if m == nil {
			continue
		}
		if "welle-"+m[1] == welleID {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out, nil
}

// FindSlice findet die Slice-Datei in <root>/docs/plan/planning/done/, deren
// Name mit "<sliceID>-" beginnt. Genau eine Datei ist der Regelfall; keine
// oder mehrdeutig sind beides Fehler -- anders als FindWellePlan gibt es
// keinen legitimen Nullfall (ein wellenloser Einzel-Slice OHNE eigene Datei
// existiert nicht, ADR/Slice-Existenz wird zwar nicht erzwungen, aber der
// Aufrufer nennt hier ausdruecklich eine Kennung, die auflösen muss).
func FindSlice(root, sliceID string) (string, error) {
	dir := filepath.Join(root, "docs/plan/planning/done")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("done-Verzeichnis lesen: %w", err)
	}
	prefix := sliceID + "-"
	var found []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".md") {
			found = append(found, name)
		}
	}
	switch len(found) {
	case 0:
		return "", fmt.Errorf("%s nicht gefunden unter %s", sliceID, RelPath(root, dir))
	case 1:
		return filepath.Join(dir, found[0]), nil
	default:
		sort.Strings(found)
		return "", fmt.Errorf("mehrdeutig: %d Kandidaten fuer %s: %v", len(found), sliceID, found)
	}
}

// FindReview findet einen eigenstaendigen Review-Report unter
// <root>/docs/reviews/ per EXAKTEM Dateinamen -- anders als FindSlice
// (Praefix-Suche ueber die Slice-Nummer) kennt der Aufrufer hier bereits den
// vollen Dateinamen (aus `ls docs/reviews/`), ein Praefix-Match waere
// unnoetige Mehrdeutigkeit.
func FindReview(root, filename string) (string, error) {
	p := filepath.Join(root, "docs/reviews", filename)
	if k, err := os.Stat(p); err != nil || k.IsDir() {
		return "", fmt.Errorf("%s nicht gefunden unter docs/reviews/", filename)
	}
	return p, nil
}

// CollectReviews findet Review-Reports in <root>/docs/reviews/, deren
// Dateiname eine der Slice-Kennungen aus slicePaths traegt. 1:N zulaessig
// (mehrere Reviews desselben Slice, z. B. -r1/-r2-Suffixe).
func CollectReviews(root string, slicePaths []string) ([]string, error) {
	ids := map[string]bool{}
	for _, p := range slicePaths {
		m := sliceIDInNameRE.FindStringSubmatch(filepath.Base(p))
		if m != nil {
			ids["slice-"+m[1]] = true
		}
	}
	dir := filepath.Join(root, "docs/reviews")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reviews-Verzeichnis lesen: %w", err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := sliceIDInNameRE.FindStringSubmatch(e.Name())
		if m == nil || !ids["slice-"+m[1]] {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	sort.Strings(out)
	return out, nil
}
