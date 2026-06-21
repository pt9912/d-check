// Package cli ist der einzige driving Adapter von d-check
// (spec/architecture.md §2; Ordnerkonvention gemäß ADR-0005):
// Argument-Parsing, Composition Root, Exit-Code-Mapping
// (DC-FA-CLI-001…004).
package cli

import (
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/core/app"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	configyaml "github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	fsadapter "github.com/pt9912/d-check/internal/adapter/driven/fs"
	"github.com/pt9912/d-check/internal/adapter/driven/httpcheck"
	"github.com/pt9912/d-check/internal/adapter/driven/report"
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
	root          string
	json          bool
	yaml          bool
	doctor        bool
	trace         bool
	repair        bool
	repairBroad   bool
	enable        []string
	disable       []string
	printConfig   bool
	printMK       bool
	suggestConfig string
	idPrefix      string
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
		"-suggest-config": true, "--suggest-config": true,
		"-id-prefix": true, "--id-prefix": true,
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

// writeUsage gibt die Hilfe aus (DC-FA-CLI-001.a): Kurzbeschreibung,
// Synopsis mit dem Pfad-Argument, Flag-Liste und ein Konfigurations-
// Hinweis, der auf --print-config/--suggest-config verweist (das
// Config-Format wird dort gezeigt, nicht hier dupliziert).
func writeUsage(flags *flag.FlagSet) {
	out := flags.Output()
	fmt.Fprintln(out, "d-check — prüft Markdown-Dokumentation auf kaputte Referenzen")
	fmt.Fprintln(out, "(lokale Links, Heading-Anker, Kennungs-Linkpflicht, Referenzmatrix u. a.).")
	fmt.Fprintln(out, "\nAufruf:")
	fmt.Fprintln(out, "  d-check [optionen] [pfad]")
	fmt.Fprintln(out, "\n  [pfad]   Scan-Wurzel (Default: aktuelles Verzeichnis); gilt als Repo-Wurzel.")
	fmt.Fprintln(out, "\nOptionen:")
	flags.PrintDefaults()
	fmt.Fprintln(out, "\nKonfiguration (optionale .d-check.yml in der Repo-Wurzel):")
	fmt.Fprintln(out, "  d-check --print-config             kommentiertes Start-Gerüst ausgeben")
	fmt.Fprintln(out, "  d-check --suggest-config <quelle>  Gerüst aus Autoritäts-Quellen vorschlagen")
}

// splitSources zerlegt den --suggest-config-Wert in einzelne Quellen
// (kommagetrennt, getrimmt, leere verworfen).
func splitSources(v string) []string {
	var out []string
	for _, s := range strings.Split(v, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// comboError meldet die erste unzulässige Modus-Kombination (DC-FA-CLI-003/
// 004); leerer String = zulässig. --doctor und --trace sind jeweils mit
// --json/--yaml kombinierbar (maschinenlesbar), aber nicht mit --repair
// bzw. untereinander.
func comboError(o options) string {
	switch {
	case o.json && o.yaml:
		return "--json und --yaml sind nicht kombinierbar"
	case o.repair && (o.json || o.yaml):
		return "--repair ist nicht mit --json/--yaml kombinierbar"
	case o.doctor && o.repair:
		return "--doctor und --repair sind nicht kombinierbar"
	case o.trace && (o.doctor || o.repair):
		return "--trace ist nicht mit --doctor/--repair kombinierbar"
	}
	return ""
}

// earlyGenerators behandelt die repo-freien stdout-Generatoren
// (--print-config/--print-mk) — Kurzschluss vor jedem Repo-Zugriff;
// true ⇒ behandelt (Run endet mit Exit 0; DC-FA-CLI-005/010).
func earlyGenerators(o options, stdout io.Writer) bool {
	switch {
	case o.printConfig:
		fmt.Fprint(stdout, configTemplate)
	case o.printMK:
		fmt.Fprint(stdout, makefileFragment())
	default:
		return false
	}
	return true
}

// runTrace gibt die Requirements Traceability Matrix aus (DC-FA-CLI-009) —
// read-only; ausgelagert aus Run, um dessen Komplexität niedrig zu halten.
func runTrace(fsys driven.Filesystem, opts options, stdout, stderr io.Writer) int {
	matrix, err := app.BuildTraceMatrix(fsys)
	if err == nil {
		switch {
		case opts.json:
			err = report.TraceJSON(stdout, matrix)
		case opts.yaml:
			err = report.TraceYAML(stdout, matrix)
		default:
			err = report.Trace(stdout, matrix)
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	return 0
}

// parseOptions parst die Argumente; code/done steuern den sofortigen
// Exit (Usage-Fehler → 2, -h → 0).
func parseOptions(args []string, stderr io.Writer) (options, int, bool) {
	flags := flag.NewFlagSet("d-check", flag.ContinueOnError)
	flags.SetOutput(io.Discard) // Fehlertexte einheitlich unten
	var enable, disable multiFlag
	jsonOut := flags.Bool("json", false, "maschinenlesbare JSON-Ausgabe")
	yamlOut := flags.Bool("yaml", false, "maschinenlesbare YAML-Ausgabe (gleiche Struktur wie --json)")
	doctorOut := flags.Bool("doctor", false, "erklärende, gruppierte Diagnose mit Fix-Kandidaten auf stdout (statt Befund-Zeilen)")
	traceOut := flags.Bool("trace", false, "Requirements Traceability Matrix (RTM) auf stdout — read-only, kein Dokument erzeugt; mit --json/--yaml maschinenlesbar")
	repairOut := flags.Bool("repair", false, "Reparatur-Patch (unified diff) auf stdout, git-apply-kompatibel; konservativ (nur eindeutige Fixes)")
	repairBroadOut := flags.Bool("repair-broad", false, "wie --repair, zusätzlich Best-Guess-Reparaturen (review-pflichtig, Marker auf stderr)")
	printConfig := flags.Bool("print-config", false, "Konfigurations-Startgerüst auf stdout ausgeben und beenden")
	printMK := flags.Bool("print-mk", false, "include-bares d-check.mk (doc-check-Target, gepinntes Image) auf stdout ausgeben und beenden")
	suggestConfig := flags.String("suggest-config", "", "Config aus Autoritäts-Quellen (kommagetrennt) vorschlagen und beenden")
	idPrefix := flags.String("id-prefix", "", "Kennungs-Präfix für --suggest-config ai-harness[-init] (z. B. AC); ohne Angabe Platzhalter <PREFIX>")
	flags.Var(&enable, "enable", "Regelmodul aktivieren (wiederholbar)")
	flags.Var(&disable, "disable", "Regelmodul deaktivieren (wiederholbar)")
	flags.Usage = func() { writeUsage(flags) }

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
	opts := options{root: ".", json: *jsonOut, yaml: *yamlOut, doctor: *doctorOut,
		repair: *repairOut || *repairBroadOut, repairBroad: *repairBroadOut,
		enable: enable, disable: disable, printConfig: *printConfig, suggestConfig: *suggestConfig,
		idPrefix: *idPrefix, trace: *traceOut, printMK: *printMK}
	if flags.NArg() == 1 {
		opts.root = flags.Arg(0)
	}
	// DC-FA-CLI-004: --json und --yaml sind dasselbe Dokument in zwei
	// Serialisierungen und schließen sich aus. --doctor IST mit --json
	// oder --yaml kombinierbar (maschinenlesbare Diagnose, DC-FA-CLI-007).
	// Nutzungsfehler (DC-FA-CLI-003, Exit 2): --json+--yaml,
	// --repair+--json/--yaml und --doctor+--repair.
	if msg := comboError(opts); msg != "" {
		fmt.Fprintln(stderr, "d-check: error: "+msg)
		return options{}, 2, true
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
func loadConfig(fsys *fsadapter.Adapter, stderr io.Writer) (model.Config, bool) {
	var content []byte
	if kind, err := fsys.Kind(configyaml.FileName); err == nil && kind == driven.KindFile {
		read, err := fsys.ReadFile(configyaml.FileName)
		if err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return model.Config{}, false
		}
		content = read
	}
	cfg, err := configyaml.Decode(content)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return model.Config{}, false
	}
	return cfg, true
}

// render gibt das Ergebnis aus und liefert den Exit-Code
// (DC-FA-CLI-003/004). Die Ausgabe-Modi --doctor und --json ersetzen das
// Default-Text-Format; der Exit-Code richtet sich allein nach dem
// Befund-Stand (0 = keine, 1 = mindestens einer).
func render(res rules.Result, opts options, cfg model.Config, fsys driven.Filesystem, stdout, stderr io.Writer) int {
	code, err := renderStdout(res, opts, cfg, fsys, stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	return code
}

// renderStdout wählt das Ausgabeformat und schreibt es auf stdout; es
// liefert den Befund-Exit-Code (0/1) und einen etwaigen Render-Fehler, den
// der Aufrufer auf Exit 2 abbildet. Die Modi --doctor/--repair sowie
// --json/--yaml ersetzen das Default-Text-Format (DC-FA-CLI-004).
func renderStdout(res rules.Result, opts options, cfg model.Config, fsys driven.Filesystem, stdout, stderr io.Writer) (int, error) {
	exit := 0
	if len(res.Findings) > 0 {
		exit = 1
	}
	sum := report.Summary{FilesChecked: res.FilesChecked, FindingCount: len(res.Findings)}
	switch {
	case opts.doctor && opts.json:
		return exit, report.DoctorJSON(stdout, res.Findings, sum, exit, cfg)
	case opts.doctor && opts.yaml:
		return exit, report.DoctorYAML(stdout, res.Findings, sum, exit, cfg)
	case opts.doctor:
		return exit, report.Doctor(stdout, stderr, res.Findings, sum, cfg)
	case opts.repair:
		edits, err := app.RepairEdits(fsys, res.Findings, cfg, opts.repairBroad)
		if err != nil {
			return exit, err
		}
		return exit, report.Repair(stdout, stderr, edits)
	case opts.json:
		return exit, report.JSON(stdout, res.Findings, sum, exit)
	case opts.yaml:
		return exit, report.YAML(stdout, res.Findings, sum, exit)
	default:
		return exit, report.Text(stdout, stderr, res.Findings, sum)
	}
}

// Run führt das CLI aus und liefert den Prozess-Exit-Code
// (DC-FA-CLI-003: 0 = keine Befunde, 1 = Befunde,
// 2 = Nutzungs-/Umgebungsfehler).
func Run(args []string, stdout, stderr io.Writer) int {
	opts, code, done := parseOptions(args, stderr)
	if done {
		return code
	}
	// DC-FA-CLI-005: Kurzschluss VOR jedem Repo-Zugriff (kein Scan,
	// kein Lesen/Schreiben) — statisches Gerüst auf stdout, Exit 0.
	// Repo-freie stdout-Generatoren (--print-config/--print-mk): Kurzschluss
	// vor jedem Repo-Zugriff (DC-FA-CLI-005/010), dann Exit 0.
	if earlyGenerators(opts, stdout) {
		return 0
	}
	fsys, ok := openRoot(opts.root, stderr)
	if !ok {
		return 2
	}
	// DC-FA-CLI-006: Lese-Durchgang, gibt ein vorgeschlagenes Gerüst auf
	// stdout aus (schreibt nie). openRoot oben hat die Wurzel validiert.
	if opts.suggestConfig != "" {
		sources := splitSources(opts.suggestConfig)
		if len(sources) == 0 {
			fmt.Fprintln(stderr, "d-check: error: --suggest-config braucht mindestens eine Quelle")
			return 2
		}
		// --id-prefix (slice-037): leer = Platzhalter/Ableitung; sonst muss
		// der Wert der Präfix-Gestalt entsprechen (Negative: Exit 2).
		if opts.idPrefix != "" && !app.ValidIDPrefix(opts.idPrefix) {
			fmt.Fprintf(stderr, "d-check: error: ungültiges --id-prefix %q (erwartet Großbuchstaben-Präfix wie AC)\n", opts.idPrefix)
			return 2
		}
		out, err := app.SuggestConfig(fsys, sources, opts.idPrefix)
		if err != nil {
			fmt.Fprintf(stderr, "d-check: error: %v\n", err)
			return 2
		}
		fmt.Fprint(stdout, out)
		return 0
	}
	// DC-FA-CLI-009 (slice-036): Requirements Traceability Matrix —
	// read-only, kein Dokument erzeugt (ausgelagert: hält Run schlank).
	if opts.trace {
		return runTrace(fsys, opts, stdout, stderr)
	}
	cfg, ok := loadConfig(fsys, stderr)
	if !ok {
		return 2
	}
	modules, err := model.EffectiveModules(cfg, opts.enable, opts.disable)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	// Composition Root: der HTTP-Adapter wird nur verdrahtet, wenn
	// das Modul external aktiv ist — ohne Aktivierung existiert
	// strukturell keine Netzwerk-Tür (DC-QA-03; der Kern behandelt
	// nil als No-op).
	var checker driven.HTTPChecker
	for _, m := range modules {
		if m == "external" {
			checker = httpcheck.New(time.Duration(cfg.External.EffectiveTimeoutSeconds()) * time.Second)
			break
		}
	}
	res, err := rules.Run(fsys, checker, cfg, modules)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	return render(res, opts, cfg, fsys, stdout, stderr)
}
