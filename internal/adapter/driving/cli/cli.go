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
	"path"
	"path/filepath"
	"strings"
	"time"

	configyaml "github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	fsadapter "github.com/pt9912/d-check/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/d-check/internal/adapter/driven/git"
	workflowyaml "github.com/pt9912/d-check/internal/adapter/driven/workflowyaml"
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
	root            string
	json            bool
	yaml            bool
	doctor          bool
	trace           bool
	repair          bool
	repairBroad     bool
	enable          []string
	disable         []string
	printConfig     bool
	printMK         bool
	suggestConfig   string
	idPrefix        string
	requireComplete bool
	vcsRange        string
	vcsStaged       bool
	commitMsg       string
	configPath      string
}

// reorderArgs erlaubt Optionen auch NACH dem Pfad-Argument — nötig
// für das Image-Aufrufmuster aus DC-FA-DIST-001/ADR-0002
// (ENTRYPOINT ["/d-check","/repo"], Optionen werden angehängt).
// Ein wertnehmendes Flag ohne Wert ist ein Nutzungsfehler
// (sonst würde das Pfad-Argument als Wert verschluckt).
func reorderArgs(args []string) ([]string, error) {
	valueFlags := map[string]bool{
		"-enable": true, "--enable": true,
		"-disable": true, "--disable": true,
		"-suggest-config": true, "--suggest-config": true,
		"-id-prefix": true, "--id-prefix": true,
		"-range": true, "--range": true,
		"-commit-msg": true, "--commit-msg": true,
		"-config": true, "--config": true,
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

// flagWasSet meldet, ob ein Flag auf der Kommandozeile EXPLIZIT gesetzt wurde —
// die Unterscheidung, die der Wert allein nicht trägt (ein gesetztes `--config ""`
// sieht aus wie ein abwesendes Flag).
func flagWasSet(flags *flag.FlagSet, name string) bool {
	set := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

// handbuchURL und handbuchURLRaw zeigen beide auf den HAUPTZWEIG. Damit
// stehen in beiden Ausgaben dieselben zwei Zeiger, und ein Release muss sie
// nicht anfassen. Zwei Formen statt einer, seit slice-187: `blob` rendert
// GitHub-seitig (fuer den Menschen, der klickt), `raw` liefert den
// Quelltext roh (fuer den Code-Agenten, der ihn zieht) -- gemessen 174,6 KB
// gegen 1,2 MB fuer denselben Inhalt, und `blob` verliert dabei jedes
// Markdown-Linkziel. Die Reihenfolge in der Ausgabe haelt `blob` zuerst
// (unveraenderte Gewohnheit), `raw` danach mit eigener Beschriftung.
//
// GRENZE: kein Sensor dieses Repos loest eine der beiden auf -- d-check
// scannt nur Markdown (DC-FA-STRUCT-001), diese Datei also nie, und
// `external` ist strikt opt-in und nirgends aktiviert. Zeigt eine URL ins
// Leere, faellt es keinem Lauf auf. Zweiter Preis: wer ein aelteres Image
// faehrt, liest ein neueres Handbuch.
//
// NUR handbuchURL (die blob-Form) STEHT EIN ZWEITES MAL in
// packaging/dockerhub/overview.md, ohne Kopplung -- eine Umbenennung braeche
// beide still. Die raw-Form hat dort KEIN Gegenstueck: die Docker-Hub-Seite
// ist ein handkuratierter Link fuer einen Menschen, keine Illustration von
// --help/--print-mk-Output, und braucht die maschinenfreundliche Form nicht.
// Eine Kopplung ueber das eigene `versions`-Modul traegt hier ohnehin nicht: der Fundort ist eine
// Go-Konstante, d-check scannt ausschliesslich Markdown (BEO-ALL/semantic-change-body-only-edges-stale).
const (
	handbuchURL    = "https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md"
	handbuchURLRaw = "https://raw.githubusercontent.com/pt9912/d-check/refs/heads/main/docs/user/benutzerhandbuch.md"
)

// writeUsage gibt die Hilfe aus (DC-FA-CLI-001.a): Kurzbeschreibung,
// Synopsis mit dem Pfad-Argument, Flag-Liste, ein Konfigurations-Hinweis,
// der auf --print-config/--suggest-config verweist (das Config-Format wird
// dort gezeigt, nicht hier dupliziert), und die beiden Handbuch-Zeiger.
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
	fmt.Fprintln(out, "\nBenutzerhandbuch (aufgabenorientiert, deutsch):")
	fmt.Fprintln(out, "  "+handbuchURL)
	fmt.Fprintln(out, "  "+handbuchURLRaw+"  (roh, für Werkzeuge/Agenten)")
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
	case o.requireComplete && !o.trace:
		return "--require-complete erfordert --trace"
	}
	return gitComboError(o)
}

// gitComboError prüft die Modus-Exklusivität der git-nutzenden Flags
// (--range/--staged für vcs, --commit-msg für commits) — ausgelagert, damit
// comboError unter der gocyclo-Schwelle bleibt.
func gitComboError(o options) string {
	switch {
	case o.vcsStaged && o.vcsRange != "":
		return "--range und --staged sind nicht kombinierbar"
	case o.commitMsg != "" && (o.vcsRange != "" || o.vcsStaged || o.trace || o.doctor || o.repair):
		return "--commit-msg ist ein eigener Modus (nicht mit --range/--staged/--trace/--doctor/--repair kombinierbar)"
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
// Mit opts.requireComplete (DC-FA-CLI-011) bindet es Waisen an Exit 1.
func runTrace(fsys driven.Filesystem, cfg model.Config, opts options, stdout, stderr io.Writer) int {
	matrix, err := app.BuildTraceMatrix(fsys, cfg.Trace)
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
	if !opts.requireComplete {
		return 0
	}
	return requireCompleteExit(matrix, stderr)
}

// requireCompleteExit wertet die opt-in Strict-Exit-Invariante (DC-FA-CLI-011):
// ≥1 gatende Requirements-Waise **oder** — bei aktivem trace.cross-consistency —
// ≥1 Kreuzverweis-Differenz (DC-FA-XREF-001) ⇒ Exit 1 (Befund-Code,
// DC-FA-CLI-003). Der Default-Lauf ohne --require-complete bleibt advisory Exit 0
// (DC-FA-CLI-009 unangetastet); die RTM steht bereits auf stdout, die Zähl-Zeilen
// gehen auf stderr. Beide Ursachen melden getrennt — sonst verschwiege die eine
// die andere.
func requireCompleteExit(matrix app.TraceMatrix, stderr io.Writer) int {
	exit := 0
	// Gatend sind ohne modality alle Waisen, mit modality nur die
	// require-levels-Stufen (DC-FA-MOD-001, slice-068).
	if matrix.GatingOrphans > 0 {
		// Waise = ohne Slice, und — bei aktiver trace.coverage — auch ohne
		// Coverage (DC-FA-COV-001/DC-FA-CLI-011).
		lack := "ohne referenzierenden Slice"
		if matrix.CoverageActive {
			lack = "ohne referenzierenden Slice und ohne Coverage"
		}
		if matrix.ModalityActive {
			// Bei aktivem modality: gatende von der Gesamt-Waisenzahl.
			fmt.Fprintf(stderr, "d-check: %d gatende von %d Requirements-Waise(n) %s (--require-complete, require-levels)\n",
				matrix.GatingOrphans, matrix.Orphans, lack)
		} else {
			fmt.Fprintf(stderr, "d-check: %d Requirements-Waise(n) %s (--require-complete)\n", matrix.Orphans, lack)
		}
		exit = 1
	}
	if n := len(matrix.CrossConsistency); n > 0 {
		fmt.Fprintf(stderr, "d-check: %d Kreuzverweis-Differenz(en) zwischen Vorwärts- und Rück-Sicht (--require-complete)\n", n)
		exit = 1
	}
	return exit
}

// runCommitMsg prüft eine einzelne Commit-Message (Modul commits,
// --commit-msg <datei|->) — der Kurzschluss-Modus für den commit-msg-Hook: ohne
// Repo-Scan und ohne VCS-Port (die Pending-Message existiert noch nicht als
// Commit). Fail-closed ohne konfigurierte commits.id-patterns oder bei nicht
// lesbarer Message (Exit 2); ein `commit-untraceable`-Befund ⇒ Exit 1 (Befund-Zeile
// im Text-Format auf stdout), sonst Exit 0 (DC-FA-COMMITS-001.a).
func runCommitMsg(opts options, cfg model.Config, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(cfg.Commits.IDPatterns) == 0 {
		fmt.Fprintln(stderr, "d-check: error: --commit-msg braucht konfigurierte commits.id-patterns (.d-check.yml)")
		return 2
	}
	raw, err := readCommitMsg(opts.commitMsg, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: Commit-Message nicht lesbar: %v\n", err)
		return 2
	}
	f, bad := rules.CheckCommitMessage("pending", raw, cfg.Commits)
	if !bad {
		return 0
	}
	fmt.Fprintf(stdout, "%s:%d\t%s\t%s\n", f.File, f.Line, f.Target, f.Reason)
	fmt.Fprintf(stderr, "d-check: commit-untraceable — Commit-Message ohne DC-/ADR-/MR-/slice-ID: %q\n", f.Message)
	return 1
}

// readCommitMsg liest die Commit-Message aus <datei> (oder stdin bei "-").
func readCommitMsg(path string, stdin io.Reader) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(stdin)
	}
	return os.ReadFile(path)
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
	printMK := flags.Bool("print-mk", false, "include-bares d-check.mk (doc-check/doc-trace/doc-complete-Targets, gepinntes Image) auf stdout ausgeben und beenden")
	requireComplete := flags.Bool("require-complete", false, "mit --trace: ≥1 Requirements-Waise ⇒ Exit 1 statt 0 (opt-in Vollständigkeits-Gate); ohne --trace Nutzungsfehler")
	vcsRange := flags.String("range", "", "Module vcs/commits: git-Commit-Range <base>..<head> (vcs: git-Diff-Immutabilität DC-FA-VCS-001; commits: Traceability-Kennung DC-FA-COMMITS-001; fail-closed ohne .git)")
	vcsStaged := flags.Bool("staged", false, "Modul vcs: staged-Diff (lokaler pre-commit) statt --range")
	commitMsg := flags.String("commit-msg", "", "Modul commits: eine Commit-Message aus <datei> (oder - für stdin) auf eine Traceability-Kennung prüfen (commit-msg-Hook; DC-FA-COMMITS-001)")
	suggestConfig := flags.String("suggest-config", "", "Config aus Autoritäts-Quellen (kommagetrennt) vorschlagen und beenden")
	idPrefix := flags.String("id-prefix", "", "Kennungs-Präfix für --suggest-config ai-harness[-init] (z. B. AC); ohne Angabe Platzhalter <PREFIX>")
	configPath := flags.String("config", "", "Konfigurationsdatei statt der konventionellen "+configyaml.FileName+" der Scan-Wurzel (Pfad relativ zur Wurzel; ersetzt, ergänzt nicht; fehlt sie ⇒ Exit 2)")
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
	// Ein LEERER --config-Wert ist vom abwesenden Flag nicht am Wert
	// unterscheidbar (beides ""), wäre aber ein stiller Rückfall auf die
	// konventionelle Datei — genau das, was DC-FA-CLI-012 verbietet. Realer
	// Auslöser: eine nicht expandierte Make-/CI-Variable. Also am
	// gesetzt-Zustand unterscheiden, nicht am Wert.
	if flagWasSet(flags, "config") && *configPath == "" {
		fmt.Fprintln(stderr, "d-check: error: --config braucht einen Pfad (leerer Wert — kein Rückfall auf "+
			configyaml.FileName+")")
		return options{}, 2, true
	}
	opts := options{root: ".", json: *jsonOut, yaml: *yamlOut, doctor: *doctorOut,
		repair: *repairOut || *repairBroadOut, repairBroad: *repairBroadOut,
		enable: enable, disable: disable, printConfig: *printConfig, suggestConfig: *suggestConfig,
		idPrefix: *idPrefix, trace: *traceOut, printMK: *printMK, requireComplete: *requireComplete,
		vcsRange: *vcsRange, vcsStaged: *vcsStaged, commitMsg: *commitMsg,
		configPath: *configPath}
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
// Ein per --config gesetzter Pfad (DC-FA-CLI-012) **ersetzt** die
// konventionelle Datei: er wird relativ zur Scan-Wurzel aufgelöst, muss
// innerhalb der Wurzel liegen und existieren — beides fail-closed (Exit 2),
// ohne Rückfall auf Defaults oder auf die konventionelle Datei. Ein
// vertipptes Profil darf nicht unbemerkt einen anderen Prüfumfang fahren.
func loadConfig(fsys *fsadapter.Adapter, configPath string, stderr io.Writer) (model.Config, bool) {
	var content []byte
	switch {
	case configPath != "":
		read, ok := readConfigOverride(fsys, configPath, stderr)
		if !ok {
			return model.Config{}, false
		}
		content = read
	default:
		if kind, err := fsys.Kind(configyaml.FileName); err == nil && kind == driven.KindFile {
			read, err := fsys.ReadFile(configyaml.FileName)
			if err != nil {
				fmt.Fprintf(stderr, "d-check: error: %v\n", err)
				return model.Config{}, false
			}
			content = read
		}
	}
	loaded := configyaml.FileName
	if configPath != "" {
		loaded = normalizeConfigPath(fsys, configPath)
	}
	cfg, err := configyaml.Decode(content)
	if err != nil {
		// Die Meldung muss die geladene Datei nennen, nicht die konventionelle.
		fmt.Fprintf(stderr, "d-check: error: %v\n", configyaml.RetargetFileName(err, loaded))
		return model.Config{}, false
	}
	if len(content) > 0 {
		cfg.ConfigFile = loaded
	}
	return cfg, true
}

// readConfigOverride liest die per --config genannte Datei (DC-FA-CLI-012):
// Auflösung relativ zur Scan-Wurzel, Wurzel-Constraint und Existenz jeweils
// fail-closed. Kein Rückfall — weder auf Defaults noch auf die konventionelle
// Datei.
func readConfigOverride(fsys *fsadapter.Adapter, configPath string, stderr io.Writer) ([]byte, bool) {
	rel := normalizeConfigPath(fsys, configPath)
	if rel == "" {
		fmt.Fprintf(stderr,
			"d-check: error: --config %s liegt außerhalb der Scan-Wurzel\n", configPath)
		return nil, false
	}
	// Die lexikalische Prüfung oben hält nur `..` und absolute Pfade. Ein
	// Verzeichnis-**Symlink** entkäme ihr (`esc/etc/passwd` mit `esc -> /`), denn
	// Kind lstatet nur die letzte Komponente. Also segmentweise prüfen — dieselbe
	// Alias-Disziplin wie bei den Link-Zielen.
	if seg, ok := firstSymlinkSegment(fsys, rel); ok {
		fmt.Fprintf(stderr,
			"d-check: error: --config %s führt über den Symlink %s aus der Scan-Wurzel heraus\n",
			configPath, seg)
		return nil, false
	}
	kind, err := fsys.Kind(rel)
	if err != nil || kind != driven.KindFile {
		fmt.Fprintf(stderr,
			"d-check: error: --config %s nicht gefunden oder keine Datei (kein Rückfall auf %s)\n",
			configPath, configyaml.FileName)
		return nil, false
	}
	content, err := fsys.ReadFile(rel)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return nil, false
	}
	return content, true
}

// firstSymlinkSegment prüft jedes Pfad-Präfix von rel (inklusive rel selbst) auf
// Symlink und liefert das erste gefundene. Ein Symlink irgendwo im Pfad macht die
// Wurzel-Constraint wertlos, weil das reale Ziel beliebig woanders liegen kann.
func firstSymlinkSegment(fsys *fsadapter.Adapter, rel string) (string, bool) {
	segs := strings.Split(rel, "/")
	for i := range segs {
		prefix := strings.Join(segs[:i+1], "/")
		if kind, err := fsys.Kind(prefix); err == nil && kind == driven.KindSymlink {
			return prefix, true
		}
	}
	return "", false
}

// normalizeConfigPath bildet den --config-Wert auf einen Wurzel-relativen,
// '/'-getrennten Pfad ab. Ein absoluter Pfad wird auf die Wurzel bezogen; was
// nach Auflösung der `..`-Segmente außerhalb liegt, liefert "" (der Aufrufer
// bricht ab). Die Wurzel ist die Grenze, weil der Lauf im read-only-Mount
// stattfindet — eine Datei außerhalb wäre dort ohnehin unerreichbar.
func normalizeConfigPath(fsys *fsadapter.Adapter, configPath string) string {
	p := filepath.FromSlash(configPath)
	if filepath.IsAbs(p) {
		relToRoot, err := filepath.Rel(fsys.Root, p)
		if err != nil {
			return ""
		}
		p = relToRoot
	}
	clean := path.Clean(filepath.ToSlash(p))
	if clean == ".." || strings.HasPrefix(clean, "../") || clean == "." {
		return ""
	}
	return clean
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
	sum := report.Summary{FilesChecked: res.FilesChecked, FindingCount: len(res.Findings), Notes: res.Notes}
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
	cfg, ok := loadConfig(fsys, opts.configPath, stderr)
	if !ok {
		return 2
	}
	// DC-FA-CLI-009 (slice-036/066): Requirements Traceability Matrix —
	// read-only, kein Dokument erzeugt (ausgelagert: hält Run schlank). Braucht
	// cfg (den opt-in trace-Block, slice-066), daher nach loadConfig; ein
	// trace-Config-Fehler ⇒ Exit 2 (fail-closed).
	if opts.trace {
		return runTrace(fsys, cfg, opts, stdout, stderr)
	}
	// DC-FA-COMMITS-001: --commit-msg-Kurzschluss (Modul commits) — prüft eine
	// einzelne Pending-Message ohne Scan/VCS-Port (commit-msg-Hook). Braucht cfg
	// (commits.id-patterns), daher nach loadConfig.
	if opts.commitMsg != "" {
		return runCommitMsg(opts, cfg, os.Stdin, stdout, stderr)
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
	checker := httpChecker(modules, cfg)
	// Der VCS-/git-Adapter wird analog nur verdrahtet, wenn das Modul vcs
	// aktiv ist — ohne Aktivierung existiert strukturell keine git-Tür
	// (DC-FA-VCS-001; fail-closed ohne .git/Range → Exit 2).
	vcsPort, vcsBase, vcsHead, code := resolveVCS(opts, modules, stderr)
	if code != 0 {
		return code
	}
	res, err := rules.RunWithVCS(fsys, checker, vcsPort, workflowyaml.New(), vcsBase, vcsHead, cfg, modules)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return 2
	}
	return render(res, opts, cfg, fsys, stdout, stderr)
}

// httpChecker verdrahtet den HTTP-Adapter nur, wenn ein Netz-Modul aktiv ist
// (external oder sources — DC-QA-03: ohne Aktivierung keine Netzwerk-Tür; nil
// ist im Kern No-op). Beide teilen Client und Timeout (external.timeout-seconds).
func httpChecker(modules []string, cfg model.Config) driven.HTTPChecker {
	for _, m := range modules {
		if m == "external" || m == "sources" {
			return httpcheck.New(time.Duration(cfg.External.EffectiveTimeoutSeconds()) * time.Second)
		}
	}
	return nil
}

// gitModuleActive prüft, ob ein range-nutzendes git-Modul (vcs oder commits)
// im effektiven Modulsatz ist — beide teilen den VCS-/git-Adapter und die Range.
func gitModuleActive(modules []string) bool {
	return moduleActive(modules, "vcs") || moduleActive(modules, "commits")
}

// moduleActive prüft die Mitgliedschaft eines Moduls im effektiven Satz.
func moduleActive(modules []string, name string) bool {
	for _, m := range modules {
		if m == name {
			return true
		}
	}
	return false
}

// vcsRefs leitet base/head aus --range/--staged ab (DC-FA-VCS-001.a Schritt 1);
// msg != "" ist ein fail-closed Nutzungsfehler (Exit 2). --staged: BASE=HEAD,
// HEAD=staged Index. --range: <base>..<head> mit Pflicht-Separator und nicht
// leeren Seiten.
func vcsRefs(opts options) (base, head, msg string) {
	switch {
	case opts.vcsStaged:
		return "HEAD", driven.IndexRef, ""
	case opts.vcsRange != "":
		b, h, ok := strings.Cut(opts.vcsRange, "..")
		if !ok || b == "" || h == "" {
			return "", "", "--range braucht <base>..<head>"
		}
		return b, h, ""
	default:
		return "", "", "ein git-nutzendes Modul (vcs/commits) braucht --range <base>..<head> oder --staged (commits: alternativ --commit-msg)"
	}
}

// resolveVCS verdrahtet den VCS-/git-Adapter und löst die Range auf, wenn das
// Modul vcs aktiv ist (DC-FA-VCS-001). code: 0 = ok (Port nil, falls vcs
// inaktiv), 2 = fail-closed (Nutzungs-/git-Fehler, bereits auf stderr gemeldet —
// fehlende/kaputte Range oder fehlendes/unlesbares .git).
func resolveVCS(opts options, modules []string, stderr io.Writer) (driven.VCS, string, string, int) {
	// tracked braucht den Adapter, aber KEINE Range (nur der Index-Stand,
	// DC-FA-TRK-001); vcs/commits brauchen Adapter + Range.
	rangeNeeded := gitModuleActive(modules)
	indexNeeded := moduleActive(modules, "tracked")
	if !rangeNeeded && !indexNeeded {
		return nil, "", "", 0
	}
	var base, head string
	if rangeNeeded {
		var msg string
		base, head, msg = vcsRefs(opts)
		if msg != "" {
			fmt.Fprintln(stderr, "d-check: error: "+msg)
			return nil, "", "", 2
		}
	}
	absRoot, err := filepath.Abs(opts.root)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return nil, "", "", 2
	}
	adapter, err := gitadapter.Open(absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "d-check: error: %v\n", err)
		return nil, "", "", 2
	}
	return adapter, base, head, 0
}
