// Package cli ist der einzige driving Adapter von d-check
// (spec/architecture.md §2; Ordnerkonvention gemäß ADR-0005):
// Argument-Parsing, Composition Root, Exit-Code-Mapping
// (DC-FA-CLI-001…004).
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	configyaml "github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	fsadapter "github.com/pt9912/d-check/internal/adapter/driven/fs"
	"github.com/pt9912/d-check/internal/adapter/driven/httpcheck"
	"github.com/pt9912/d-check/internal/adapter/driven/report"
	"github.com/pt9912/d-check/internal/hexagon/core"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// mountHint ergänzt im Container-Kontext den von DC-FA-DIST-001
// geforderten Hinweis auf den erwarteten Mount.
func mountHint(absRoot string) string {
	if absRoot == "/repo" {
		return " — wurde das Repository nach /repo gemountet? (docker run -v \"$PWD:/repo:ro\" …)"
	}
	return ""
}

type multiFlag []string

func (m *multiFlag) String() string { return fmt.Sprint([]string(*m)) }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// options sind die geparsten CLI-Eingaben.
type options struct {
	root    string
	json    bool
	enable  []string
	disable []string
}

// reorderArgs erlaubt Optionen auch NACH dem Pfad-Argument — nötig
// für das Image-Aufrufmuster aus DC-FA-DIST-001/ADR-0002
// (ENTRYPOINT ["/d-check","/repo"], Optionen werden angehängt).
// Ein wertnehmendes Flag ohne Wert ist ein Nutzungsfehler
// (Review R2/A1: sonst würde das Pfad-Argument als Wert verschluckt).
func reorderArgs(args []string) ([]string, error) {
	valueFlags := map[string]bool{
		"-enable": true, "--enable": true,
		"-disable": true, "--disable": true,
	}
	var flagArgs, positionals []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if len(a) > 1 && a[0] == '-' {
			flagArgs = append(flagArgs, a)
			if valueFlags[a] {
				if i+1 >= len(args) {
					return nil, fmt.Errorf("flag needs an argument: %s", a)
				}
				i++
				flagArgs = append(flagArgs, args[i])
			}
			continue
		}
		positionals = append(positionals, a)
	}
	return append(flagArgs, positionals...), nil
}

// parseOptions parst die Argumente; code/done steuern den sofortigen
// Exit (Usage-Fehler → 2, -h → 0).
func parseOptions(args []string, stderr io.Writer) (options, int, bool) {
	flags := flag.NewFlagSet("d-check", flag.ContinueOnError)
	flags.SetOutput(io.Discard) // Fehlertexte einheitlich unten
	var enable, disable multiFlag
	jsonOut := flags.Bool("json", false, "maschinenlesbare JSON-Ausgabe")
	flags.Var(&enable, "enable", "Regelmodul aktivieren (wiederholbar)")
	flags.Var(&disable, "disable", "Regelmodul deaktivieren (wiederholbar)")

	reordered, err := reorderArgs(args)
	if err == nil {
		err = flags.Parse(reordered)
	}
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			flags.SetOutput(stderr)
			flags.Usage()
			return options{}, 0, true
		}
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return options{}, 2, true
	}
	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "d-check: error: höchstens ein Pfad-Argument")
		return options{}, 2, true
	}
	opts := options{root: ".", json: *jsonOut, enable: enable, disable: disable}
	if flags.NArg() == 1 {
		opts.root = flags.Arg(0)
	}
	return opts, 0, false
}

// openRoot validiert die Scan-Wurzel (existiert, Verzeichnis, nicht
// gänzlich leer — DC-FA-DIST-001 Negative) und liefert den
// Filesystem-Adapter.
func openRoot(root string, stderr io.Writer) (*fsadapter.Adapter, bool) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return nil, false
	}
	if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "d-check: error: Scan-Wurzel nicht gefunden: %s%s\n",
			root, mountHint(absRoot))
		return nil, false
	}
	fsys := fsadapter.New(absRoot)
	if entries, err := fsys.List(""); err == nil && len(entries) == 0 {
		fmt.Fprintf(stderr, "d-check: error: Scan-Wurzel ist leer: %s%s\n",
			root, mountHint(absRoot))
		return nil, false
	}
	return fsys, true
}

// loadConfig beschafft den Config-Inhalt über den Filesystem-Adapter
// und lässt ihn vom Config-Adapter dekodieren/validieren
// (spec/architecture.md §2).
func loadConfig(fsys *fsadapter.Adapter, stderr io.Writer) (core.Config, bool) {
	var content []byte
	if kind, err := fsys.Kind(configyaml.FileName); err == nil && kind == driven.KindFile {
		read, err := fsys.ReadFile(configyaml.FileName)
		if err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return core.Config{}, false
		}
		content = read
	}
	cfg, err := configyaml.Decode(content)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return core.Config{}, false
	}
	return cfg, true
}

// render gibt das Ergebnis aus und liefert den Exit-Code
// (DC-FA-CLI-003/004).
func render(res core.Result, jsonOut bool, stdout, stderr io.Writer) int {
	exit := 0
	if len(res.Findings) > 0 {
		exit = 1
	}
	sum := report.Summary{FilesChecked: res.FilesChecked, FindingCount: len(res.Findings)}
	if jsonOut {
		if err := report.JSON(stdout, res.Findings, sum, exit); err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return 2
		}
		return exit
	}
	if err := report.Text(stdout, stderr, res.Findings, sum); err != nil {
		return 2
	}
	return exit
}

// Run führt das CLI aus und liefert den Prozess-Exit-Code
// (DC-FA-CLI-003: 0 = keine Befunde, 1 = Befunde,
// 2 = Nutzungs-/Umgebungsfehler).
func Run(args []string, stdout, stderr io.Writer) int {
	opts, code, done := parseOptions(args, stderr)
	if done {
		return code
	}
	fsys, ok := openRoot(opts.root, stderr)
	if !ok {
		return 2
	}
	cfg, ok := loadConfig(fsys, stderr)
	if !ok {
		return 2
	}
	modules, err := core.EffectiveModules(cfg, opts.enable, opts.disable)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	// Composition Root: der HTTP-Adapter wird nur vom explizit
	// aktivierten Modul external benutzt (DC-QA-03).
	checker := httpcheck.New(time.Duration(cfg.External.EffectiveTimeoutSeconds()) * time.Second)
	res, err := core.Run(fsys, checker, cfg, modules)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	return render(res, opts.json, stdout, stderr)
}
