// Package cli ist der einzige driving Adapter von d-check
// (spec/architecture.md §2; Ordnerkonvention gemäß ADR-0005):
// Argument-Parsing, Composition Root, Exit-Code-Mapping
// (DC-FA-CLI-001…004).
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	configyaml "github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	fsadapter "github.com/pt9912/d-check/internal/adapter/driven/fs"
	"github.com/pt9912/d-check/internal/adapter/driven/report"
	"github.com/pt9912/d-check/internal/hexagon/core"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

type multiFlag []string

func (m *multiFlag) String() string { return fmt.Sprint([]string(*m)) }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// Run führt das CLI aus und liefert den Prozess-Exit-Code
// (DC-FA-CLI-003: 0 = keine Befunde, 1 = Befunde,
// 2 = Nutzungs-/Umgebungsfehler).
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("d-check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var enable, disable multiFlag
	jsonOut := flags.Bool("json", false, "maschinenlesbare JSON-Ausgabe")
	flags.Var(&enable, "enable", "Regelmodul aktivieren (wiederholbar)")
	flags.Var(&disable, "disable", "Regelmodul deaktivieren (wiederholbar)")
	if err := flags.Parse(args); err != nil {
		return 2 // ungültige Nutzung
	}

	root := "."
	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "d-check: error: höchstens ein Pfad-Argument")
		return 2
	}
	if flags.NArg() == 1 {
		root = flags.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "d-check: error: Scan-Wurzel nicht gefunden: %s\n", root)
		return 2
	}

	fsys := fsadapter.New(absRoot)

	// Config-Inhalt beschafft das CLI über den Filesystem-Adapter;
	// der Config-Adapter dekodiert/validiert nur
	// (spec/architecture.md §2).
	var content []byte
	if kind, kerr := fsys.Kind(configyaml.FileName); kerr == nil && kind == driven.KindFile {
		content, err = fsys.ReadFile(configyaml.FileName)
		if err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return 2
		}
	}
	cfg, err := configyaml.Decode(content)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}

	modules, err := core.EffectiveModules(cfg, enable, disable)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}

	res, err := core.Run(fsys, cfg, modules)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	for _, m := range res.SkippedModules {
		fmt.Fprintf(stderr, "d-check: Hinweis: Modul %q ist noch nicht implementiert — übersprungen\n", m)
	}

	exit := 0
	if len(res.Findings) > 0 {
		exit = 1
	}
	sum := report.Summary{FilesChecked: res.FilesChecked, FindingCount: len(res.Findings)}
	if *jsonOut {
		if err := report.JSON(stdout, res.Findings, sum, exit); err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return 2
		}
	} else {
		if err := report.Text(stdout, stderr, res.Findings, sum); err != nil {
			return 2
		}
	}
	return exit
}
