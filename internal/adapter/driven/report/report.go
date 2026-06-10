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
