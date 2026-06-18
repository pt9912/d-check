// Package report ist der Reporter-Adapter: Text- und JSON-Rendering
// der Befunde (DC-FA-CLI-004, Formate gemäß spec/spezifikation.md
// §2). Die Writer (stdout/stderr) injiziert das CLI.
package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pt9912/d-check/internal/hexagon/core"
)

// Summary sind die Lauf-Kennzahlen.
type Summary struct {
	FilesChecked int `json:"filesChecked"`
	FindingCount int `json:"findingCount"`
}

// Text schreibt Befunde zeilenweise auf stdout
// (`<file>:<line>\t<target>\t<reason>`) und die Zusammenfassung auf
// stderr.
func Text(stdout, stderr io.Writer, findings []core.Finding, sum Summary) error {
	for _, f := range findings {
		if _, err := fmt.Fprintf(stdout, "%s:%d\t%s\t%s\n", f.File, f.Line, f.Target, f.Reason); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(stderr, "d-check: %d Datei(en) geprüft, %d Befund(e)\n",
		sum.FilesChecked, sum.FindingCount)
	return err
}

// Doctor schreibt die erklärende, nach Datei gruppierte Diagnose auf
// stdout (statt der knappen Befund-Zeilen — DC-FA-CLI-007) und die
// Zusammenfassung auf stderr. Je Befund: Grund-Klartext, Regel, Stelle
// und — wo eindeutig ableitbar — ein Fix-Kandidat (nicht angewendet).
// Die Befunde sind bereits stabil sortiert (DC-QA-02); die Ausgabe
// iteriert nur in dieser Reihenfolge und über die geordneten
// ids-Muster, leckt also keine Map-Reihenfolge.
func Doctor(stdout, stderr io.Writer, findings []core.Finding, sum Summary, cfg core.Config) error {
	if len(findings) == 0 {
		if _, err := fmt.Fprintln(stdout, "d-check Diagnose — keine Befunde."); err != nil {
			return err
		}
		return doctorSummary(stderr, sum)
	}
	if _, err := fmt.Fprintf(stdout, "d-check Diagnose — %d Befund(e) in %d Datei(en):\n",
		len(findings), distinctFiles(findings)); err != nil {
		return err
	}
	curFile := ""
	for _, f := range findings {
		if f.File != curFile {
			curFile = f.File
			if _, err := fmt.Fprintf(stdout, "\n%s\n", curFile); err != nil {
				return err
			}
		}
		if err := doctorFinding(stdout, f, cfg); err != nil {
			return err
		}
	}
	return doctorSummary(stderr, sum)
}

// doctorFinding rendert einen Befund samt optionalem Fix-Kandidaten
// (nil-Kandidat: nur die Befund-Zeilen).
func doctorFinding(stdout io.Writer, f core.Finding, cfg core.Config) error {
	if _, err := fmt.Fprintf(stdout, "  Z. %d · %s [%s]\n      Stelle: %s\n",
		f.Line, core.ReasonText(f.Reason), f.Rule, f.Target); err != nil {
		return err
	}
	c := core.FixCandidateFor(f, cfg)
	if c == nil {
		return nil
	}
	_, err := fmt.Fprintf(stdout, "      Fix-Kandidat: %s → %s\n        (%s)\n",
		c.Original, c.Replacement, c.Note)
	return err
}

// doctorSummary schreibt die Lauf-Zusammenfassung auf stderr (wie der
// Default-Reporter).
func doctorSummary(stderr io.Writer, sum Summary) error {
	_, err := fmt.Fprintf(stderr, "d-check: %d Datei(en) geprüft, %d Befund(e)\n",
		sum.FilesChecked, sum.FindingCount)
	return err
}

// distinctFiles zählt die Dateien mit mindestens einem Befund; nutzt die
// stabile Sortierung (gleiche Dateien sind zusammenhängend).
func distinctFiles(findings []core.Finding) int {
	n, cur := 0, ""
	for _, f := range findings {
		if f.File != cur {
			cur = f.File
			n++
		}
	}
	return n
}

// Repair schreibt den unified diff der Reparatur-Edits auf stdout
// (`git apply`-kompatibel) und die review-pflichtig-Markierungen samt
// Zusammenfassung auf stderr (DC-FA-CLI-008). Die Best-Guess-Marker
// bleiben damit AUSSERHALB des Patches — der stdout-Patch bleibt
// `git apply`-rein. Edits sind stabil sortiert (DC-QA-02).
func Repair(stdout, stderr io.Writer, edits []core.RepairEdit) error {
	curFile := ""
	review := 0
	for _, e := range edits {
		if e.File != curFile {
			curFile = e.File
			if _, err := fmt.Fprintf(stdout, "--- a/%s\n+++ b/%s\n", e.File, e.File); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(stdout, "@@ -%d,1 +%d,1 @@\n-%s\n+%s\n",
			e.Line, e.Line, e.OldLine, e.NewLine); err != nil {
			return err
		}
		if e.ReviewRequired {
			review++
			if _, err := fmt.Fprintf(stderr, "review-pflichtig (Best-Guess): %s:%d\n", e.File, e.Line); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(stderr, "d-check: %d Reparatur-Hunk(s), %d review-pflichtig\n", len(edits), review)
	return err
}

type jsonDoc struct {
	Findings []core.Finding `json:"findings"`
	Summary  Summary        `json:"summary"`
	ExitCode int            `json:"exitCode"`
}

// JSON schreibt das maschinenlesbare Dokument auf stdout; stdout
// enthält dann keine unstrukturierten Zeilen (DC-FA-CLI-004).
func JSON(stdout io.Writer, findings []core.Finding, sum Summary, exitCode int) error {
	if findings == nil {
		findings = []core.Finding{}
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonDoc{Findings: findings, Summary: sum, ExitCode: exitCode})
}
