// Akzeptanztests des CLI gegen die Lastenheft-Kriterien
// (Happy/Boundary/Negative) — echtes Dateisystem via t.TempDir
// (u-boot-Konvention: package cli_test).
package cli_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	configyaml "github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// DC-FA-CLI-001 Happy: fehlerfreies Repo → Exit 0 + Zusammenfassung.
func TestCLI001_Happy(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[ok](b.md)")
	write(t, root, "docs/b.md", "x")
	code, _, stderr := run(t, "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stderr, "2 Datei(en) geprüft, 0 Befund(e)") {
		t.Fatalf("Zusammenfassung fehlt: %q", stderr)
	}
}

// DC-FA-CLI-001 Boundary: keine Markdown-Dateien (aber Inhalt) →
// 0 geprüft, Exit 0.
func TestCLI001_Boundary_KeineMarkdownDateien(t *testing.T) {
	root := t.TempDir()
	write(t, root, "src/main.txt", "kein markdown")
	code, _, stderr := run(t, "--disable", "anchors", root)
	if code != 0 || !strings.Contains(stderr, "0 Datei(en) geprüft") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-DIST-001 Negative: gänzlich leere Wurzel → Exit 2 mit Hinweis.
func TestDIST001_GaenzlichLeereWurzel(t *testing.T) {
	root := t.TempDir()
	code, _, stderr := run(t, root)
	if code != 2 || !strings.Contains(stderr, "leer") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-DIST-001/ADR-0002: Optionen dürfen NACH dem Pfad stehen
// (Container-Aufrufmuster: ENTRYPOINT setzt /repo, Optionen angehängt).
func TestDIST001_OptionenNachPfad(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[kaputt](fehlt.md)")
	code, stdout, _ := run(t, root, "--json", "--disable", "anchors")
	if code != 1 {
		t.Fatalf("Exit = %d, want 1", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
		t.Fatalf("kein JSON auf stdout: %q", stdout)
	}
}

// Review R2/A1: hängendes wertnehmendes Flag am Ende darf das
// Pfad-Argument nicht als Wert verschlucken → Nutzungsfehler Exit 2.
func TestR2A1_HaengendesWertFlag(t *testing.T) {
	root := t.TempDir()
	write(t, root, "links/doc.md", "[kaputt](fehlt.md)")
	code, _, stderr := run(t, "links", "--disable")
	if code != 2 || !strings.Contains(stderr, "needs an argument") {
		t.Fatalf("Exit = %d, stderr = %q (Nutzungsfehler erwartet)", code, stderr)
	}
}

// Flag-Fehler tragen den d-check: error:-Präfix; -h endet mit Exit 0.
func TestFlagFehlerUndHelp(t *testing.T) {
	code, _, stderr := run(t, "--foo")
	if code != 2 || !strings.Contains(stderr, "d-check: error:") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	code, stdout, stderr := run(t, "-h")
	if code != 0 || stdout != "" || !strings.Contains(stderr, "-json") ||
		!strings.Contains(stderr, "[pfad]") || !strings.Contains(stderr, "print-config") {
		t.Fatalf("-h: Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
	// mehr als ein Pfad-Argument → Exit 2 (DC-FA-CLI-001)
	code, _, stderr = run(t, "pfad1", "pfad2")
	if code != 2 || !strings.Contains(stderr, "höchstens ein Pfad-Argument") {
		t.Fatalf("zwei Pfade: Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-CONF-001: leere bzw. Nur-Kommentar-Config = Defaults.
func TestCONF001_LeereConfig(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", "")
	write(t, root, "docs/a.md", "[ok](b.md)")
	write(t, root, "docs/b.md", "x")
	code, _, stderr := run(t, "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("leere Config: Exit = %d, stderr = %q", code, stderr)
	}
	write(t, root, ".d-check.yml", "# nur ein Kommentar\n")
	code, _, stderr = run(t, "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("Kommentar-Config: Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-CLI-001 Negative: nicht existierendes Verzeichnis → Exit 2.
func TestCLI001_Negative_FehlendeWurzel(t *testing.T) {
	code, _, stderr := run(t, "/gibt/es/nicht")
	if code != 2 || !strings.Contains(stderr, "error") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-CLI-002 Happy: deaktiviertes Modul läuft nicht.
func TestCLI002_DisableLinks(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[kaputt](fehlt.md)")
	code, _, _ := run(t, "--disable", "links", "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (links deaktiviert)", code)
	}
}

// DC-FA-CLI-002 Boundary: CLI-Präzedenz vor Config.
func TestCLI002_CLIVorConfig(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", "modules: [links]\n")
	write(t, root, "docs/a.md", "[kaputt](fehlt.md)")
	code, _, _ := run(t, "--disable", "links", root)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (CLI schlägt Config)", code)
	}
}

// DC-FA-CLI-002 Negative: unbekanntes Modul → Exit 2 + gültige Namen.
func TestCLI002_UnbekanntesModul(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "x")
	code, _, stderr := run(t, "--enable", "foo", root)
	if code != 2 || !strings.Contains(stderr, "links, anchors, ids, matrix, external") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-FA-CLI-003 Boundary: genau ein Befund → Exit 1 (nicht 2).
func TestCLI003_EinBefund(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[kaputt](fehlt.md)")
	code, stdout, _ := run(t, "--disable", "anchors", root)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1", code)
	}
	if !strings.Contains(stdout, "docs/a.md:1\tfehlt.md\ttarget-missing") {
		t.Fatalf("Befund-Zeile fehlt: %q", stdout)
	}
}

// DC-FA-CLI-003 Negative: ungültige Option → Exit 2, keine Prüfausgabe.
func TestCLI003_UngueltigeOption(t *testing.T) {
	code, stdout, _ := run(t, "--foo")
	if code != 2 || stdout != "" {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-FA-CLI-004: Text-Format und --json.
func TestCLI004_Ausgabeformate(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[k1](f1.md)\n[k2](f2.md)")

	code, stdout, _ := run(t, "--disable", "anchors", root)
	if code != 1 || strings.Count(stdout, "\n") != 2 {
		t.Fatalf("Text: Exit = %d, stdout = %q", code, stdout)
	}

	code, stdout, _ = run(t, "--json", "--disable", "anchors", root)
	if code != 1 {
		t.Fatalf("JSON: Exit = %d", code)
	}
	var doc struct {
		Findings []map[string]any `json:"findings"`
		Summary  struct {
			FilesChecked int `json:"filesChecked"`
			FindingCount int `json:"findingCount"`
		} `json:"summary"`
		ExitCode int `json:"exitCode"`
	}
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("stdout ist kein JSON: %v\n%q", err, stdout)
	}
	if doc.Summary.FindingCount != 2 || doc.ExitCode != 1 || len(doc.Findings) != 2 {
		t.Fatalf("JSON-Inhalt: %+v", doc)
	}

	// unbekannte Option --format → Exit 2 (DC-FA-CLI-004 Negative)
	code, _, _ = run(t, "--json", "--format", "xml", root)
	if code != 2 {
		t.Fatalf("--format: Exit = %d, want 2", code)
	}
}

// DC-FA-CLI-004 YAML: --yaml gibt dieselbe Struktur wie --json aus, nur als
// YAML; parsbar, gleiche findingCount, camelCase-Schlüssel-Parität, Exit 1.
func TestCLI004_YAML(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[k1](f1.md)\n[k2](f2.md)")
	code, stdout, _ := run(t, "--yaml", "--disable", "anchors", root)
	if code != 1 {
		t.Fatalf("Exit = %d", code)
	}
	var doc struct {
		Findings []map[string]any `yaml:"findings"`
		Summary  struct {
			FilesChecked int `yaml:"filesChecked"`
			FindingCount int `yaml:"findingCount"`
		} `yaml:"summary"`
		ExitCode int `yaml:"exitCode"`
	}
	if err := yaml.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("stdout ist kein YAML: %v\n%q", err, stdout)
	}
	if doc.Summary.FindingCount != 2 || doc.ExitCode != 1 || len(doc.Findings) != 2 {
		t.Fatalf("YAML-Inhalt: %+v", doc)
	}
	// Schlüssel-Parität: camelCase wie JSON (nicht kleingeschrieben durch yaml.v3).
	if !strings.Contains(stdout, "findingCount:") || strings.Contains(stdout, "fileschecked:") {
		t.Fatalf("YAML-Schlüssel weichen von JSON ab:\n%s", stdout)
	}
}

// DC-FA-CLI-007 + DC-FA-CLI-004: --doctor --yaml = YAML-Diagnose mit
// reasonText/fixCandidate (analog --doctor --json), eingebettete
// Befund-Felder flach (yaml:",inline"), fixCandidate explizit null.
func TestCLI007_DoctorYAML(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids, links]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n[x](fehlt.md)\n")
	code, stdout, stderr := run(t, "--doctor", "--yaml", root)
	if code != 1 {
		t.Fatalf("Exit = %d (stderr: %s)", code, stderr)
	}
	var doc struct {
		Findings []struct {
			Reason       string `yaml:"reason"`
			ReasonText   string `yaml:"reasonText"`
			FixCandidate *struct {
				Replacement string `yaml:"replacement"`
			} `yaml:"fixCandidate"`
		} `yaml:"findings"`
		ExitCode int `yaml:"exitCode"`
	}
	if err := yaml.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("kein YAML: %v\n%q", err, stdout)
	}
	if doc.ExitCode != 1 || len(doc.Findings) != 2 {
		t.Fatalf("YAML-Diagnose: %+v", doc)
	}
	withCand, nullCand := 0, 0
	for _, f := range doc.Findings {
		if f.ReasonText == "" {
			t.Fatalf("reasonText fehlt (inline?): %+v", f)
		}
		switch f.Reason {
		case "id-unlinked":
			if f.FixCandidate == nil || !strings.Contains(f.FixCandidate.Replacement, "[`ADR-0042`](") {
				t.Fatalf("id-unlinked ohne fixCandidate: %+v", f)
			}
			withCand++
		default:
			if f.FixCandidate == nil {
				nullCand++
			}
		}
	}
	if withCand != 1 || nullCand != 1 {
		t.Fatalf("Kandidaten: with=%d null=%d", withCand, nullCand)
	}
	if !strings.Contains(stdout, "fixCandidate: null") {
		t.Fatalf("fixCandidate nicht explizit null:\n%s", stdout)
	}
}

// DC-FA-CLI-004 YAML Negative: --json+--yaml und --repair+--yaml → Exit 2.
func TestCLI004_YAML_Negative(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "x")
	code, stdout, stderr := run(t, "--json", "--yaml", root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "nicht kombinierbar") {
		t.Fatalf("--json --yaml: Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
	code, _, stderr = run(t, "--repair", "--yaml", root)
	if code != 2 || !strings.Contains(stderr, "nicht mit --json/--yaml") {
		t.Fatalf("--repair --yaml: Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-QA-02 für die YAML-Ausgabe: 10 Läufe byte-identisch.
func TestCLI004_YAML_Determinismus(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/z.md", "[f](f1.md)\n[g](f2.md)")
	write(t, root, "docs/a.md", "[h](f3.md)")
	var first string
	for i := 0; i < 10; i++ {
		_, stdout, _ := run(t, "--yaml", "--disable", "anchors", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab", i)
		}
	}
}

// DC-FA-SCAN-001: node_modules wird übersprungen; explizite Wurzel
// muss existieren; Ignore-Muster aus Config.
func TestSCAN001(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[ok](b.md)")
	write(t, root, "docs/b.md", "x")
	write(t, root, "node_modules/x/README.md", "[kaputt](fehlt.md)")
	code, _, _ := run(t, "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("node_modules nicht übersprungen: Exit = %d", code)
	}

	write(t, root, ".d-check.yml", "scan:\n  roots: [handbuch]\n")
	code, _, stderr := run(t, root)
	if code != 2 || !strings.Contains(stderr, "handbuch") {
		t.Fatalf("fehlende Wurzel: Exit = %d, stderr = %q", code, stderr)
	}

	write(t, root, ".d-check.yml", "scan:\n  ignore: [\"docs/archive/**\"]\n")
	write(t, root, "docs/archive/alt.md", "[kaputt](fehlt.md)")
	code, _, _ = run(t, "--disable", "anchors", root)
	if code != 0 {
		t.Fatalf("Ignore-Muster wirkt nicht: Exit = %d", code)
	}
}

// DC-FA-LINK-002: Symlinks erzeugen genau einen Befund mit Grund symlink.
func TestLINK002_Symlinks(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[s](ziel.md)")
	write(t, root, "docs/echt.md", "x")
	if err := os.Symlink(filepath.Join(root, "docs/echt.md"), filepath.Join(root, "docs/ziel.md")); err != nil {
		t.Skipf("Symlinks nicht verfügbar: %v", err)
	}
	code, stdout, _ := run(t, "--disable", "anchors", root)
	if code != 1 || strings.Count(stdout, "symlink") != 1 || strings.Count(stdout, "\n") != 1 {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-QA-02: wiederholter Lauf liefert byte-identische Ausgabe.
func TestQA02_Determinismus(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/z.md", "[f](f1.md)\n[g](f2.md)")
	write(t, root, "docs/a.md", "[h](f3.md)")
	var first string
	for i := 0; i < 10; i++ { // 10 Läufe gemäß DC-QA-02-Messmethode
		_, stdout, _ := run(t, "--json", "--disable", "anchors", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab", i)
		}
	}
}

// DC-FA-EXT-001: Modul external (opt-in) — Happy (erreichbar),
// Negative (404 → external-status), Boundary (nicht aktiviert →
// keine Befunde für externe Links).
func TestEXT001(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path == "/fehlt" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	root := t.TempDir()
	write(t, root, "docs/a.md", "[ok]("+srv.URL+"/da)\n[kaputt]("+srv.URL+"/fehlt)\n")

	// Boundary: Modul nicht aktiviert → kein Befund, kein Request
	code, _, _ := run(t, "--disable", "anchors", root)
	if code != 0 || requests != 0 {
		t.Fatalf("opt-in verletzt: Exit = %d, Requests = %d", code, requests)
	}

	// Happy + Negative: aktiviert → genau der 404-Link wird gemeldet
	code, stdout, _ := run(t, "--enable", "external", "--disable", "anchors", root)
	if code != 1 || !strings.Contains(stdout, "external-status") ||
		strings.Contains(stdout, "/da\t") {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-FA-CODE-001: Pfade in Inline-Code (Modul codepaths, opt-in) —
// Happy (existierender Pfad), Boundary (Marker), Negative
// (codepath-missing, Exit 1).
func TestCODE001(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", "modules: [codepaths]\ncodepaths:\n  roots: [docs]\n")
	write(t, root, "docs/a.md", "Siehe `docs/b.md`.\n")
	write(t, root, "docs/b.md", "x\n")

	code, _, stderr := run(t, root)
	if code != 0 {
		t.Fatalf("Happy: Exit = %d (stderr: %s)", code, stderr)
	}

	write(t, root, "docs/a.md", "Kaputt: `docs/fehlt.md`\n"+
		"Beispiel: `../../etc/passwd` <!-- d-check:ignore (Lehrtext) -->\n")
	code, stdout, _ := run(t, root)
	if code != 1 || !strings.Contains(stdout, "codepath-missing") ||
		strings.Contains(stdout, "repo-escape") {
		t.Fatalf("Exit = %d, stdout = %q (genau codepath-missing erwartet)", code, stdout)
	}
}

// DC-FA-ID-001: Linkpflicht für Kennungen (Modul ids) — Happy
// (verlinkt), Boundary (Inline-Code), Negative (nackt → id-unlinked).
func TestID001(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-beispiel.md", "# ADR")
	write(t, root, "docs/a.md",
		"[ADR-0042](plan/adr/0042-beispiel.md)\n`ADR-0042` in Inline-Code\n")
	code, _, stderr := run(t, root)
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (stderr: %s)", code, stderr)
	}

	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n")
	code, stdout, _ := run(t, root)
	if code != 1 || !strings.Contains(stdout, "docs/a.md:1\tADR-0042\tid-unlinked") {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-FA-CLI-007 Happy: --doctor gibt eine erklärende, gruppierte Diagnose
// mit Fix-Kandidat auf stdout aus (statt der knappen Befund-Zeile),
// Exit 1, Zusammenfassung auf stderr.
func TestCLI007_Doctor_Happy(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n")
	code, stdout, stderr := run(t, "--doctor", root)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1 (stderr: %s)", code, stderr)
	}
	for _, want := range []string{"Diagnose", "docs/a.md", "ohne Markdown-Link", "Fix-Kandidat", "[`ADR-0042`]("} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Diagnose ohne %q:\n%s", want, stdout)
		}
	}
	// Default-Format ist ERSETZT, nicht ergänzt (DC-FA-CLI-004).
	if strings.Contains(stdout, "docs/a.md:1\tADR-0042\tid-unlinked") {
		t.Fatalf("knappe Befund-Zeile trotz --doctor:\n%s", stdout)
	}
	if !strings.Contains(stderr, "1 Befund(e)") {
		t.Fatalf("Zusammenfassung auf stderr fehlt: %q", stderr)
	}
}

// DC-FA-CLI-007 Boundary: keine Befunde → Diagnose weist „keine Befunde"
// aus, Exit 0, keine Fix-Kandidaten.
func TestCLI007_Doctor_KeineBefunde(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[ok](b.md)")
	write(t, root, "docs/b.md", "x")
	code, stdout, _ := run(t, "--doctor", "--disable", "anchors", root)
	if code != 0 || !strings.Contains(stdout, "keine Befunde") || strings.Contains(stdout, "Fix-Kandidat") {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-FA-CLI-007 JSON-Variante: --doctor --json gibt die Diagnose
// maschinenlesbar aus — je findings-Eintrag reasonText und fixCandidate
// (Objekt mit replacement, sonst explizit null); stdout ist reines JSON,
// Exit 1. Zwei Befunde: id-unlinked (mit Kandidat), target-missing (null).
func TestCLI007_DoctorJSON_Happy(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids, links]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n[x](fehlt.md)\n")
	code, stdout, stderr := run(t, "--doctor", "--json", root)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1 (stderr: %s)", code, stderr)
	}
	var doc struct {
		Findings []struct {
			File         string `json:"file"`
			Reason       string `json:"reason"`
			ReasonText   string `json:"reasonText"`
			FixCandidate *struct {
				Original    string `json:"original"`
				Replacement string `json:"replacement"`
				Note        string `json:"note"`
			} `json:"fixCandidate"`
		} `json:"findings"`
		Summary struct {
			FindingCount int `json:"findingCount"`
		} `json:"summary"`
		ExitCode int `json:"exitCode"`
	}
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("stdout ist kein JSON: %v\n%q", err, stdout)
	}
	if doc.ExitCode != 1 || doc.Summary.FindingCount != 2 || len(doc.Findings) != 2 {
		t.Fatalf("JSON-Inhalt: %+v", doc)
	}
	withCand, nullCand := 0, 0
	for _, f := range doc.Findings {
		if f.ReasonText == "" {
			t.Fatalf("reasonText fehlt: %+v", f)
		}
		switch f.Reason {
		case "id-unlinked":
			if f.FixCandidate == nil || !strings.Contains(f.FixCandidate.Replacement, "[`ADR-0042`](") {
				t.Fatalf("id-unlinked ohne fixCandidate.replacement: %+v", f)
			}
			withCand++
		default:
			if f.FixCandidate == nil {
				nullCand++
			}
		}
	}
	if withCand != 1 || nullCand != 1 {
		t.Fatalf("Kandidaten-Verteilung: with=%d null=%d", withCand, nullCand)
	}
	// fixCandidate ist EXPLIZIT null (nicht weggelassen) — sonst verschwindet
	// die Aussage „kein eindeutiger Fix".
	if !strings.Contains(stdout, "\"fixCandidate\": null") {
		t.Fatalf("fixCandidate nicht explizit null:\n%s", stdout)
	}
	// stdout ist reines JSON: keine Prosa-Diagnose.
	if strings.Contains(stdout, "Diagnose") || strings.Contains(stdout, "Fix-Kandidat:") {
		t.Fatalf("stdout enthält Prosa trotz --json:\n%s", stdout)
	}
}

// DC-QA-02 für die JSON-Diagnose: 10 Läufe byte-identisch (feste
// Struct-Reihenfolge, keine Map-Iteration).
func TestCLI007_DoctorJSON_Determinismus(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/z.md", "ADR-0042 und nochmal ADR-0042\n")
	write(t, root, "docs/a.md", "ADR-0042\n")
	var first string
	for i := 0; i < 10; i++ {
		_, stdout, _ := run(t, "--doctor", "--json", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab:\n%q\n---\n%q", i, first, stdout)
		}
	}
}

// DC-QA-02 für den Diagnose-Modus: 10 Läufe byte-identisch.
func TestCLI007_Doctor_Determinismus(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/z.md", "ADR-0042 und nochmal ADR-0042\n")
	write(t, root, "docs/a.md", "ADR-0042\n")
	var first string
	for i := 0; i < 10; i++ {
		_, stdout, _ := run(t, "--doctor", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab:\n%q\n---\n%q", i, first, stdout)
		}
	}
}

// DC-FA-CLI-008 Happy: --repair (konservativ) erzeugt einen unified diff,
// der eine nackte Kennung verlinkt; git apply nimmt ihn sauber an und ein
// erneuter Lauf meldet den Befund nicht mehr (Round-Trip).
func TestCLI008_Repair_RoundTrip(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git nicht verfügbar")
	}
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n")

	code, patch, _ := run(t, "--repair", root)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1", code)
	}
	if !strings.Contains(patch, "--- a/docs/a.md") || !strings.Contains(patch, "+nacktes [`ADR-0042`](") {
		t.Fatalf("Patch unerwartet:\n%s", patch)
	}

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s\npatch:\n%s", args, err, out, patch)
		}
	}
	git("init")
	if err := os.WriteFile(filepath.Join(root, "fix.patch"), []byte(patch), 0o644); err != nil {
		t.Fatal(err)
	}
	git("apply", "fix.patch")

	code2, _, _ := run(t, root) // erneuter Lauf auf dem gepatchten Baum
	if code2 != 0 {
		t.Fatalf("nach git apply noch Befunde: Exit = %d", code2)
	}
}

// DC-FA-CLI-008 Boundary: nur best-guess-fähige Befunde (target-missing) →
// konservativ leerer Patch; breite Stufe Best-Guess-Hunk + stderr-Marker.
func TestCLI008_Repair_Boundary_Stufen(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "[x](alt.md)\n")
	write(t, root, "docs/sub/alt.md", "ok\n")

	code, stdout, _ := run(t, "--repair", "--disable", "anchors", root)
	if code != 1 || stdout != "" {
		t.Fatalf("konservativ: Exit = %d, stdout = %q (leerer Patch erwartet)", code, stdout)
	}

	code, stdout, stderr := run(t, "--repair-broad", "--disable", "anchors", root)
	if code != 1 || !strings.Contains(stdout, "+[x](sub/alt.md)") {
		t.Fatalf("breit: Exit = %d, stdout = %q", code, stdout)
	}
	if !strings.Contains(stderr, "review-pflichtig") {
		t.Fatalf("breit: stderr ohne Best-Guess-Marker: %q", stderr)
	}
}

// DC-FA-CLI-008 Negative: --repair mit --json bzw. --doctor → Exit 2,
// keine Ausgabe auf stdout.
func TestCLI008_Repair_Inkompatibel(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "x")
	code, stdout, stderr := run(t, "--repair", "--json", root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "nicht mit --json") {
		t.Fatalf("--repair --json: Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
	code, _, stderr = run(t, "--doctor", "--repair", root)
	if code != 2 || !strings.Contains(stderr, "nicht kombinierbar") {
		t.Fatalf("--doctor --repair: Exit = %d, stderr = %q", code, stderr)
	}
}

// DC-QA-02 für den Patch-Modus: 10 Läufe byte-identisch (auch mehrere
// gleiche Kennungen auf einer Zeile).
func TestCLI008_Repair_Determinismus(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"modules: [ids]\nids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR")
	write(t, root, "docs/z.md", "ADR-0042 und ADR-0042\n")
	write(t, root, "docs/a.md", "ADR-0042\n")
	var first string
	for i := 0; i < 10; i++ {
		_, stdout, _ := run(t, "--repair", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab:\n%q\n---\n%q", i, first, stdout)
		}
	}
}

// DC-FA-CONF-001 (Constraint ids.patterns[].target muss existieren):
// nicht existierendes Target → Exit 2 ohne Prüfung.
func TestCONF001_IDTargetFehlt(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml",
		"ids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: gibt/es/nicht\n")
	write(t, root, "docs/a.md", "x")
	code, stdout, stderr := run(t, root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "gibt/es/nicht") {
		t.Fatalf("Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
}

// DC-FA-CONF-001: ids-Target außerhalb der Repo-Wurzel → Exit 2,
// auch wenn der Pfad außerhalb tatsächlich existiert.
func TestCONF001_IDTargetVerlaesstWurzel(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "repo")
	write(t, base, "draussen/def.md", "x") // existiert — außerhalb der Wurzel
	write(t, root, ".d-check.yml",
		"ids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: ../draussen\n")
	write(t, root, "docs/a.md", "x")
	code, stdout, stderr := run(t, root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "verlässt") {
		t.Fatalf("Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
}

// DC-FA-MTX-001: Referenzmatrix — Happy (Slice → aktives ADR),
// Boundary (superseded ADR → matrix-inactive), Negative
// (Lastenheft → ADR → matrix-forbidden).
func TestMTX001(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", `modules: [matrix]
matrix:
  classes:
    - name: contract
      paths: [spec/lastenheft.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
    - name: slice
      paths: ["docs/plan/planning/**/slice-*.md"]
  rules:
    - {from: contract, to: adr, allow: false}
`)
	write(t, root, "docs/plan/adr/0001-x.md", "# X\n\n**Status:** Accepted\n")
	write(t, root, "docs/plan/adr/0002-y.md", "# Y\n\n**Status:** Superseded by ADR-0042\n")
	write(t, root, "docs/plan/planning/done/slice-001-a.md", "[ok](../../adr/0001-x.md)\n")
	write(t, root, "spec/lastenheft.md", "# LH\n")
	code, _, stderr := run(t, root)
	if code != 0 {
		t.Fatalf("Happy: Exit = %d (stderr: %s)", code, stderr)
	}

	write(t, root, "docs/plan/planning/done/slice-001-a.md", "[inaktiv](../../adr/0002-y.md)\n")
	code, stdout, _ := run(t, root)
	if code != 1 || !strings.Contains(stdout, "matrix-inactive") {
		t.Fatalf("Boundary: Exit = %d, stdout = %q", code, stdout)
	}

	write(t, root, "spec/lastenheft.md", "# LH\n[abwärts](../docs/plan/adr/0001-x.md)\n")
	code, stdout, _ = run(t, root)
	if code != 1 || !strings.Contains(stdout, "matrix-forbidden") {
		t.Fatalf("Negative: Exit = %d, stdout = %q", code, stdout)
	}
}

// DC-FA-CONF-001 Negative: ungültige Config → Exit 2 mit Zeilenangabe,
// keine Prüfung mit stillschweigenden Defaults.
func TestCONF001_UngueltigeConfig(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", "modul: kaputt\n") // unbekannter Schlüssel
	write(t, root, "docs/a.md", "x")
	code, stdout, stderr := run(t, root)
	if code != 2 || stdout != "" {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
	}
	if !strings.Contains(stderr, "line") && !strings.Contains(stderr, "Zeile") {
		t.Fatalf("Zeilenangabe fehlt: %q", stderr)
	}
}

// DC-FA-CLI-005 Happy: --print-config gibt gültiges, vom eigenen Parser
// akzeptiertes YAML auf stdout aus, Exit 0.
func TestCLI005_PrintConfig(t *testing.T) {
	code, stdout, stderr := run(t, "--print-config")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	if stdout == "" || !strings.Contains(stdout, "link-policy") || !strings.Contains(stdout, "modules:") {
		t.Fatalf("Gerüst unvollständig: %q", stdout)
	}
	if _, err := configyaml.Decode([]byte(stdout)); err != nil {
		t.Fatalf("eigener Parser lehnt das Gerüst ab: %v", err)
	}
}

// DC-FA-CLI-005 Boundary/Negative: kein Repo-Zugriff — eine nicht
// existierende Wurzel als Argument führt trotzdem zu Exit 0 (sonst
// Exit 2), und die Ausgabe ist deterministisch (DC-QA-02).
func TestCLI005_KeinRepoZugriff(t *testing.T) {
	code, stdout, _ := run(t, "--print-config", filepath.Join(t.TempDir(), "gibt-es-nicht"))
	if code != 0 {
		t.Fatalf("Repo-Zugriff trotz --print-config: Exit = %d", code)
	}
	_, stdout2, _ := run(t, "--print-config")
	if stdout != stdout2 {
		t.Fatal("Ausgabe nicht deterministisch")
	}
}

// slice-053 (Config-Surface-Bereinigung): das --print-config-Gerüst führt jetzt
// alle Module — die „Verfügbar"-Zeile nennt `vcs` (und die zuvor fehlenden
// `pins`/`immutable`), und es gibt kommentierte Blöcke für `immutable` und `vcs`.
// Das Gerüst dekodiert weiterhin fehlerfrei (alle Blöcke auskommentiert).
func TestCLI053_PrintConfig_VollesModulset(t *testing.T) {
	_, stdout, _ := run(t, "--print-config")
	for _, want := range []string{
		"versions, pins, immutable, vcs, commits, planning, tracked, targets, external", // vollständige Verfügbar-Liste
		"# --- immutable:",
		"# --- vcs:",
		"# --- commits:",
		"# --- planning:",
		"# --- tracked:",
		"# --- targets:",
		"# --- trace:", // slice-066: konfigurierbare RTM-Quellen im Gerüst
		"# --- pins:",
		"immutable-when:", // der vcs-Config-Key, nach dem gefragt wurde
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("--print-config ohne %q:\n%s", want, stdout)
		}
	}
	if _, err := configyaml.Decode([]byte(stdout)); err != nil {
		t.Fatalf("erweitertes Gerüst dekodiert nicht: %v", err)
	}
}

// DC-FA-CLI-006 Happy: --suggest-config leitet aus definierten Kennungen
// ein ids-Muster ab (Round-Trip: regex matcht die Quell-IDs), target =
// Quelle, gültiges vom eigenen Parser akzeptiertes YAML.
func TestCLI006_SuggestConfig(t *testing.T) {
	root := t.TempDir()
	// Doppelpunkt-Heading (ADR-0042:) deckt die Token-Bereinigung ab.
	write(t, root, "docs/plan/adr/0042-x.md", "# ADR-0042: Beispiel\ntext\n")
	write(t, root, "docs/plan/adr/0099-y.md", "# ADR-0099 — Beispiel\ntext\n")
	code, stdout, stderr := run(t, "--suggest-config", "docs/plan/adr/", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("Vorschlag dekodiert nicht: %v\n%s", err, stdout)
	}
	if len(cfg.IDPatterns) != 1 {
		t.Fatalf("erwartet 1 ids-Muster, got %d\n%s", len(cfg.IDPatterns), stdout)
	}
	re := cfg.IDPatterns[0].Regex
	if !re.MatchString("ADR-0042") || !re.MatchString("ADR-0099") {
		t.Fatalf("Round-Trip verletzt: %q matcht die Quell-IDs nicht", re.String())
	}
	if cfg.IDPatterns[0].Target != "docs/plan/adr/" {
		t.Fatalf("target = %q", cfg.IDPatterns[0].Target)
	}
	// das Modul ids MUSS in der Modul-Liste stehen, sonst sind die
	// abgeleiteten Muster im erzeugten Config inaktiv (Review R1).
	found := false
	for _, m := range cfg.Modules {
		if m == "ids" {
			found = true
		}
	}
	if !found {
		t.Fatalf("ids fehlt in modules %v — Muster wären inaktiv\n%s", cfg.Modules, stdout)
	}
}

// DC-FA-CLI-006: --suggest-config ohne reale Quelle (nur Trenner) ist
// ein Nutzungsfehler (Review R1), kein stilles leeres Gerüst.
func TestCLI006_LeereQuellen(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", ",", root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "mindestens eine Quelle") {
		t.Fatalf("Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
}

// DC-FA-CLI-006 Boundary: Quelle ohne Kennungs-Headings → kein ids-Muster,
// kein Absturz, gültiges YAML.
func TestCLI006_KeineKennungen(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/notes.md", "# Notizen\nkein ID-Heading hier\n")
	code, stdout, _ := run(t, "--suggest-config", "docs/", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("dekodiert nicht: %v\n%s", err, stdout)
	}
	if len(cfg.IDPatterns) != 0 {
		t.Fatalf("kein ids-Muster erwartet, got %d", len(cfg.IDPatterns))
	}
}

// DC-FA-CLI-006 Negative: nicht existierende Quelle → Exit 2.
func TestCLI006_QuelleFehlt(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, _, stderr := run(t, "--suggest-config", "gibt/es/nicht", root)
	if code != 2 || !strings.Contains(stderr, "existiert nicht") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// harnessRepo legt ein Repo mit den Harness-Artefakten an (für die
// ai-harness-Vorlage). Reihenfolge der Targets: ADR/MR/DC/slice + Scope.
func harnessRepo(t *testing.T, root string) {
	write(t, root, "spec/lastenheft.md", "# DC-FA-CLI-001 — x\n")
	write(t, root, "harness/conventions.md", "# MR-000 — x\n")
	write(t, root, "docs/plan/adr/0001-x.md", "# ADR-0001 — x\n")
	write(t, root, "docs/plan/planning/open/slice-001-x.md", "# slice-001\n")
	write(t, root, "docs/user/handbuch.md", "# Handbuch\n")
}

// DC-FA-CLI-006 ai-harness Happy: die reservierte Quelle erzeugt die
// Harness-Vorlage — kanonische ids-Muster (aktiv, da Targets existieren),
// matrix samt Referenzrichtung, Standard-Modulset; dekodiert über den
// eigenen Parser, Exit 0; read-only (kein .d-check.yml geschrieben).
func TestCLI006_AiHarness_Happy(t *testing.T) {
	root := t.TempDir()
	harnessRepo(t, root)
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("Vorlage dekodiert nicht: %v\n%s", err, stdout)
	}
	// alle vier Target-tragenden Muster aktiv (CO bleibt auskommentiert).
	if len(cfg.IDPatterns) != 4 {
		t.Fatalf("erwartet 4 aktive ids-Muster, got %d\n%s", len(cfg.IDPatterns), stdout)
	}
	var adr bool
	for _, p := range cfg.IDPatterns {
		if p.Regex.MatchString("ADR-0001") && p.Target == "docs/plan/adr/" {
			adr = true
		}
	}
	if !adr {
		t.Fatalf("ADR-Muster fehlt/aktiv-falsch\n%s", stdout)
	}
	hasMod := func(m string) bool {
		for _, x := range cfg.Modules {
			if x == m {
				return true
			}
		}
		return false
	}
	if !hasMod("ids") || !hasMod("matrix") {
		t.Fatalf("Modulset unvollständig: %v", cfg.Modules)
	}
	// spans/hostpaths sind Teil der fixen Aktiv-Menge (slice-065, DC-FA-CLI-006).
	if !hasMod("spans") || !hasMod("hostpaths") {
		t.Fatalf("spans/hostpaths fehlen im fixen Modulset: %v", cfg.Modules)
	}
	// planning ist repo-bewusst: harnessRepo legt KEINE Roadmap an → auskommentiert,
	// nicht in modules (das Modul fiele sonst fail-closed).
	if hasMod("planning") {
		t.Fatalf("planning aktiv trotz fehlender Roadmap: %v", cfg.Modules)
	}
	for _, want := range []string{"Baseline v1.3.0", "matrix:", "from: spec-straten, to: adr", "exclude-sections", "Carveouts", "external, diagrams, versions, pins, immutable, tracked, targets", "d-check --print-mk", "d-check --print-config", "Modul planning auskommentiert"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Vorlage ohne %q:\n%s", want, stdout)
		}
	}
	// read-only: das Werkzeug schreibt nichts (Aufrufer leitet um).
	if _, err := os.Stat(filepath.Join(root, ".d-check.yml")); !os.IsNotExist(err) {
		t.Fatalf(".d-check.yml wurde geschrieben — read-only verletzt (DC-QA-03)")
	}
}

// DC-FA-CLI-006 ai-harness planning repo-bewusst (slice-065): liegt die Roadmap
// im Baum, ist Modul planning aktiv (in modules) und der planning-Block gesetzt;
// fehlt sie, ist beides auskommentiert (siehe TestCLI006_AiHarness_Happy).
func TestCLI006_AiHarness_PlanningAktiv(t *testing.T) {
	root := t.TempDir()
	harnessRepo(t, root)
	write(t, root, "docs/plan/planning/in-progress/roadmap.md", "## Aktuelle Welle\n\nKeine aktive Welle.\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("dekodiert nicht: %v\n%s", err, stdout)
	}
	var planning bool
	for _, m := range cfg.Modules {
		if m == "planning" {
			planning = true
		}
	}
	if !planning {
		t.Fatalf("planning nicht in modules trotz vorhandener Roadmap: %v", cfg.Modules)
	}
	if !strings.Contains(stdout, "planning:\n  roadmap: docs/plan/planning/in-progress/roadmap.md") {
		t.Fatalf("planning-Block nicht aktiv trotz Roadmap:\n%s", stdout)
	}
	if strings.Contains(stdout, "Modul planning auskommentiert") {
		t.Fatalf("planning auskommentiert trotz vorhandener Roadmap:\n%s", stdout)
	}
}

// DC-FA-CLI-006 ai-harness Boundary: fehlt docs/plan/adr/, erscheint der
// ADR-Block (ids-Muster + matrix-Klasse) auskommentiert mit Hinweis statt
// aktiv; scan.roots enthält nur existierende Pfade.
func TestCLI006_AiHarness_Boundary(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# DC-FA-CLI-001 — x\n")
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, _ := run(t, "--suggest-config", "ai-harness", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("dekodiert nicht: %v\n%s", err, stdout)
	}
	// ADR-Muster darf NICHT aktiv sein (nur DC ist aktiv, da nur
	// spec/lastenheft.md existiert).
	for _, p := range cfg.IDPatterns {
		if p.Target == "docs/plan/adr/" {
			t.Fatalf("ADR-Muster aktiv trotz fehlendem docs/plan/adr/\n%s", stdout)
		}
	}
	for _, want := range []string{
		"Target docs/plan/adr/ fehlt im Repo",
		"docs/plan/adr fehlt im Repo — Klasse auskommentiert",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Hinweis %q fehlt:\n%s", want, stdout)
		}
	}
	// scan.roots nur existierende: harness/ fehlt → nicht in scan.roots
	// (das Wort "harness" steht zwar im ai-harness-Header, aber nicht als root).
	if !strings.Contains(stdout, "roots: [spec, docs]") || strings.Contains(stdout, "spec, docs, harness") {
		t.Fatalf("scan.roots nicht auf existierende Pfade zugeschnitten:\n%s", stdout)
	}
}

// DC-FA-CLI-006 ai-harness Abgrenzung: ai-harness ist ein reservierter
// Modus, KEINE fehlende Quelle → kein Exit 2 (obwohl keine Datei
// 'ai-harness' existiert), Exit 0.
func TestCLI006_AiHarness_KeinQuellenfehler(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# DC-FA-CLI-001 — x\n")
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", root)
	if code != 0 || stdout == "" {
		t.Fatalf("ai-harness als fehlende Quelle behandelt: Exit = %d, stderr = %q", code, stderr)
	}
	if _, err := configyaml.Decode([]byte(stdout)); err != nil {
		t.Fatalf("dekodiert nicht: %v\n%s", err, stdout)
	}
}

// DC-QA-02 für die ai-harness-Vorlage: 10 Läufe byte-identisch.
func TestCLI006_AiHarness_Determinismus(t *testing.T) {
	root := t.TempDir()
	harnessRepo(t, root)
	var first string
	for i := 0; i < 10; i++ {
		_, stdout, _ := run(t, "--suggest-config", "ai-harness", root)
		if i == 0 {
			first = stdout
		} else if stdout != first {
			t.Fatalf("Lauf %d weicht ab:\n%q\n---\n%q", i, first, stdout)
		}
	}
}

// DC-FA-CLI-006 ai-harness-init (Mode 1): Voll-Kanon — alle Blöcke aktiv,
// auch ohne vorhandene Targets (Zielbild fürs leere Repo); kein
// repo-bewusstes Auskommentieren. Decode prüft keine Target-Existenz.
func TestCLI006_AiHarnessInit_VollKanon(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n") // minimal: kein spec/, kein docs/plan/adr/
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness-init", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("Voll-Kanon dekodiert nicht: %v\n%s", err, stdout)
	}
	// alle vier Target-Muster aktiv, obwohl die Targets fehlen.
	if len(cfg.IDPatterns) != 4 {
		t.Fatalf("erwartet 4 aktive ids-Muster (Voll-Kanon), got %d\n%s", len(cfg.IDPatterns), stdout)
	}
	var adr bool
	for _, p := range cfg.IDPatterns {
		if p.Target == "docs/plan/adr/" {
			adr = true
		}
	}
	if !adr {
		t.Fatalf("ADR-Muster nicht aktiv trotz Voll-Kanon\n%s", stdout)
	}
	// slice-065: spans/hostpaths gehören zur fixen Aktiv-Menge; planning ist im
	// Voll-Kanon aktiv (Modul + Block), obwohl keine Roadmap existiert (Zielbild
	// fürs leere Repo — läuft, sobald die Struktur angelegt ist).
	hasMod := func(m string) bool {
		for _, x := range cfg.Modules {
			if x == m {
				return true
			}
		}
		return false
	}
	if !hasMod("spans") || !hasMod("hostpaths") || !hasMod("planning") {
		t.Fatalf("Voll-Kanon: spans/hostpaths/planning fehlen im Modulset: %v", cfg.Modules)
	}
	if !strings.Contains(stdout, "planning:\n  roadmap: docs/plan/planning/in-progress/roadmap.md") {
		t.Fatalf("Voll-Kanon: planning-Block nicht aktiv:\n%s", stdout)
	}
	// kein repo-bewusstes Auskommentieren; volle scan.roots.
	if strings.Contains(stdout, "fehlt im Repo") {
		t.Fatalf("Voll-Kanon kommentiert aus — unerwartet:\n%s", stdout)
	}
	if !strings.Contains(stdout, "roots: [spec, docs, harness]") {
		t.Fatalf("Voll-Kanon: scan.roots nicht vollständig:\n%s", stdout)
	}
	// read-only (DC-QA-03): auch der init-Pfad schreibt nichts.
	if _, err := os.Stat(filepath.Join(root, ".d-check.yml")); !os.IsNotExist(err) {
		t.Fatalf(".d-check.yml geschrieben — read-only verletzt (DC-QA-03)")
	}
}

// reqPattern liefert das Anforderungs-Muster (target spec/lastenheft.md)
// aus einem dekodierten Vorschlag.
func reqPattern(t *testing.T, stdout string) *regexp.Regexp {
	t.Helper()
	cfg, err := configyaml.Decode([]byte(stdout))
	if err != nil {
		t.Fatalf("dekodiert nicht: %v\n%s", err, stdout)
	}
	for _, p := range cfg.IDPatterns {
		if p.Target == "spec/lastenheft.md" {
			return p.Regex
		}
	}
	t.Fatalf("kein Anforderungs-Muster (target spec/lastenheft.md):\n%s", stdout)
	return nil
}

// slice-037 Happy: --id-prefix AC backt das Fremd-Präfix ins
// Anforderungs-Muster (statt d-checks DC-), Voll-Kanon.
func TestCLI037_IDPrefix_Explizit(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness-init", "--id-prefix", "AC", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	re := reqPattern(t, stdout)
	if !re.MatchString("AC-FA-CLI-001") {
		t.Fatalf("Muster matcht AC-FA-CLI-001 nicht: %q", re.String())
	}
	if re.MatchString("DC-FA-CLI-001") {
		t.Fatalf("Muster matcht noch DC- (Fremd-Präfix nicht angewandt): %q", re.String())
	}
}

// slice-037 Boundary: ai-harness-init ohne Präfix → markierter Platzhalter
// <PREFIX> + TODO, KEIN stiller DC--Default; Gerüst bleibt dekodierbar.
func TestCLI037_IDPrefix_Platzhalter(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, _ := run(t, "--suggest-config", "ai-harness-init", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	if !strings.Contains(stdout, "<PREFIX>-(FA-") || !strings.Contains(stdout, "TODO") {
		t.Fatalf("Platzhalter <PREFIX>/TODO fehlt:\n%s", stdout)
	}
	if strings.Contains(stdout, "DC-(FA-") {
		t.Fatalf("stiller DC--Default trotz fehlendem Präfix:\n%s", stdout)
	}
	if _, err := configyaml.Decode([]byte(stdout)); err != nil {
		t.Fatalf("Platzhalter-Gerüst dekodiert nicht: %v\n%s", err, stdout)
	}
}

// slice-037: ai-harness (repo-bewusst) leitet das Präfix aus dem
// Lastenheft ab (erste/eindeutige FA-Kennung), ohne --id-prefix.
func TestCLI037_IDPrefix_AbleitungAiHarness(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# AC-FA-CORE-001 — x\n")
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	re := reqPattern(t, stdout)
	if !re.MatchString("AC-FA-CORE-001") {
		t.Fatalf("abgeleitetes Präfix matcht AC- nicht: %q", re.String())
	}
}

// slice-037 Negative: mehrdeutiges Präfix im Lastenheft → Nutzungsfehler
// (Exit 2), der Mensch gibt --id-prefix explizit.
func TestCLI037_IDPrefix_KonfliktFehler(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# DC-FA-A-001 — x\n\n# AC-FA-B-001 — y\n")
	write(t, root, "docs/a.md", "# x\n")
	code, _, stderr := run(t, "--suggest-config", "ai-harness", root)
	if code != 2 || !strings.Contains(stderr, "mehrdeutig") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// slice-037 Negative: ungültiger --id-prefix-Wert → Exit 2.
func TestCLI037_IDPrefix_Ungueltig(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, _, stderr := run(t, "--suggest-config", "ai-harness-init", "--id-prefix", "ac", root)
	if code != 2 || !strings.Contains(stderr, "ungültiges --id-prefix") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// slice-037 (Review R1 MEDIUM-2): --id-prefix überschreibt die Ableitung im
// ai-harness-Modus (explizit gewinnt) trotz vorhandenem DC-Lastenheft.
func TestCLI037_IDPrefix_UeberschreibtAbleitung(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# DC-FA-CLI-001 — x\n")
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", "--id-prefix", "ZZ", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	re := reqPattern(t, stdout)
	if !re.MatchString("ZZ-FA-CLI-001") || re.MatchString("DC-FA-CLI-001") {
		t.Fatalf("--id-prefix überschreibt die DC-Ableitung nicht: %q", re.String())
	}
}

// slice-037 (Review R1 MEDIUM-2): --id-prefix übergeht den Mehrdeutigkeits-
// Fehler — explizites Präfix statt Ableitung, daher kein Konflikt.
func TestCLI037_IDPrefix_FlagUebergehtKonflikt(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "# DC-FA-A-001 — x\n\n# AC-FA-B-001 — y\n")
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, stderr := run(t, "--suggest-config", "ai-harness", "--id-prefix", "ZZ", root)
	if code != 0 {
		t.Fatalf("Exit = %d (Flag sollte den Konflikt übergehen), stderr = %q", code, stderr)
	}
	re := reqPattern(t, stdout)
	if !re.MatchString("ZZ-FA-A-001") {
		t.Fatalf("explizites Präfix nicht angewandt: %q", re.String())
	}
}

// traceDoc spiegelt die JSON-Struktur der RTM (DC-FA-CLI-009).
type traceDoc struct {
	Requirements []struct {
		ID       string   `json:"id"`
		Title    string   `json:"title"`
		ADRs     []string `json:"adrs"`
		Slices   []string `json:"slices"`
		Coverage []string `json:"coverage"`
		Orphan   bool     `json:"orphan"`
	} `json:"requirements"`
	Total   int `json:"total"`
	Orphans int `json:"orphans"`
}

// traceRepo: zwei Anforderungen — TRC-001 (von ADR + Slice referenziert),
// TRC-002 (Waise, kein Slice).
func traceRepo(t *testing.T, root string) {
	write(t, root, "spec/lastenheft.md",
		"### DC-FA-TRC-001 — Erste Anforderung\nText.\n\n### DC-FA-TRC-002 — Waise ohne Slice\nText.\n")
	write(t, root, "docs/plan/adr/0099-x.md", "# ADR-0099 — x\nBezug: DC-FA-TRC-001\n")
	write(t, root, "docs/plan/planning/done/slice-099-x.md", "# slice-099\nBezug: DC-FA-TRC-001\n")
}

// slice-036 Happy: --trace gibt eine Markdown-RTM mit Anforderung, ADR,
// Slice und Waisen-Status aus; read-only (DC-QA-03).
func TestCLI036_Trace_Markdown(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, stdout, stderr := run(t, "--trace", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	for _, want := range []string{"DC-FA-TRC-001", "ADR-0099", "slice-099", "DC-FA-TRC-002", "WAISE"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("RTM ohne %q:\n%s", want, stdout)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".d-check.yml")); !os.IsNotExist(err) {
		t.Fatalf("--trace hat geschrieben — read-only verletzt (DC-QA-03)")
	}
}

// slice-036: --trace --json — strukturgleich, Waisen-Erkennung.
func TestCLI036_Trace_JSON(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, stdout, _ := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON dekodiert nicht: %v\n%s", err, stdout)
	}
	if doc.Total != 2 || doc.Orphans != 1 {
		t.Fatalf("Total/Orphans = %d/%d (erwartet 2/1)\n%s", doc.Total, doc.Orphans, stdout)
	}
	for _, r := range doc.Requirements {
		switch r.ID {
		case "DC-FA-TRC-001":
			if r.Orphan || len(r.ADRs) != 1 || r.ADRs[0] != "ADR-0099" ||
				len(r.Slices) != 1 || r.Slices[0] != "slice-099" {
				t.Fatalf("TRC-001 falsch verlinkt: %+v", r)
			}
		case "DC-FA-TRC-002":
			if !r.Orphan {
				t.Fatalf("TRC-002 nicht als Waise erkannt: %+v", r)
			}
		}
	}
}

// slice-036: --trace --yaml — gültige YAML-Struktur, Waisen-Flag.
func TestCLI036_Trace_YAML(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, stdout, _ := run(t, "--trace", "--yaml", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	for _, want := range []string{"requirements:", "id: DC-FA-TRC-001", "orphan: true"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("YAML ohne %q:\n%s", want, stdout)
		}
	}
}

// slice-036 Negative: --trace + --repair ist ein Nutzungsfehler (Exit 2).
func TestCLI036_Trace_RepairKonflikt(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, _, stderr := run(t, "--trace", "--repair", root)
	if code != 2 || !strings.Contains(stderr, "nicht mit --doctor/--repair") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// slice-036 Boundary: kein Lastenheft → leere RTM, Exit 0.
func TestCLI036_Trace_KeinLastenheft(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/a.md", "# x\n")
	code, stdout, _ := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON dekodiert nicht: %v\n%s", err, stdout)
	}
	if doc.Total != 0 {
		t.Fatalf("erwartet 0 Anforderungen, got %d", doc.Total)
	}
}

// slice-036 (Review R1 LOW-1 / R2 LOW-2): backtick-umschlossene Heading-ID →
// Titel ohne Kennung/Backticks; ein Titel-initialer Code-Span bleibt aber
// intakt (führender Backtick darf nicht mit-gestrippt werden).
func TestCLI036_Trace_BacktickHeading(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md",
		"### `DC-FA-BTK-001` — Titel mit Backtick\nText.\n\n"+
			"### DC-FA-MOD-001 — `links`-Modul beschrieben\nText.\n")
	code, stdout, _ := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d", code)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON dekodiert nicht: %v\n%s", err, stdout)
	}
	want := map[string]string{
		"DC-FA-BTK-001": "Titel mit Backtick",        // backtick-umschlossene Kennung (LOW-1)
		"DC-FA-MOD-001": "`links`-Modul beschrieben", // Titel-initialer Code-Span bleibt (LOW-2)
	}
	got := map[string]string{}
	for _, r := range doc.Requirements {
		got[r.ID] = r.Title
	}
	for id, title := range want {
		if got[id] != title {
			t.Fatalf("Titel von %s = %q (erwartet %q)\n%s", id, got[id], title, stdout)
		}
	}
}

// traceReqIdx liefert den Index der Anforderung id in der RTM (oder -1).
func traceReqIdx(d traceDoc, id string) int {
	for i := range d.Requirements {
		if d.Requirements[i].ID == id {
			return i
		}
	}
	return -1
}

// slice-066 (DC-FA-CLI-009): ohne trace-Block liefert die RTM byte-identisch
// dasselbe wie mit einem trace-Block, der die Konventions-Defaults 1:1 restated —
// Regressionsschutz für resolveTrace (jeder Default-Override muss inert sein).
func TestCLI066_Trace_DefaultByteIdentisch(t *testing.T) {
	rootA := t.TempDir()
	traceRepo(t, rootA)
	_, stdoutA, _ := run(t, "--trace", rootA)

	rootB := t.TempDir()
	traceRepo(t, rootB)
	write(t, rootB, ".d-check.yml", `trace:
  requirements:
    source: spec/lastenheft.md
    id-pattern: '[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?'
  adrs:
    dir: docs/plan/adr
    file-pattern: '^(\d{4})-.*\.md$'
    id-prefix: 'ADR-'
  slices:
    dir: docs/plan/planning
    file-pattern: '^(slice-\d+)-.*\.md$'
    id-prefix: ''
`)
	_, stdoutB, _ := run(t, "--trace", rootB)

	if stdoutA != stdoutB {
		t.Fatalf("trace-Default nicht byte-identisch:\n--- ohne Block ---\n%s\n--- mit Default-Block ---\n%s", stdoutA, stdoutB)
	}
}

// slice-066 (DC-FA-CLI-009): ein trace-Block bildet eine Fremd-Konvention
// vollständig ab (grid-gym-Gestalt GG-<FAMILIE>-NNN, Slices NNN-titel.md).
// Vorher/Nachher: ohne Block trifft der Default-Regex nur die QA-Familie und
// keinen NNN-Slice; mit Block werden alle Familien + Slice-Owner erkannt.
func TestCLI066_Trace_FremdKonvention(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md",
		"### GG-ARCH-001\nHexagonale Grenzen.\n\n"+
			"### GG-QA-001\nStatische Analyse.\n\n"+
			"### GG-QA-002\nWaise ohne Slice.\n")
	write(t, root, "docs/plan/adr/0007-x.md", "# ADR-0007\nBezug: GG-ARCH-001\n")
	write(t, root, "docs/plan/planning/063-arch.md", "# 063\nBezug: GG-ARCH-001 und GG-QA-001\n")

	// Vorher: Default-Konvention.
	_, before, _ := run(t, "--trace", "--json", root)
	var db traceDoc
	if err := json.Unmarshal([]byte(before), &db); err != nil {
		t.Fatalf("JSON (vorher) dekodiert nicht: %v\n%s", err, before)
	}
	if traceReqIdx(db, "GG-ARCH-001") >= 0 {
		t.Fatalf("GG-ARCH-001 ohne Config sichtbar — der Default sollte die ARCH-Familie nicht treffen:\n%s", before)
	}
	if i := traceReqIdx(db, "GG-QA-001"); i < 0 || !db.Requirements[i].Orphan {
		t.Fatalf("GG-QA-001 ohne Config sollte Waise sein (NNN-Slice unerkannt):\n%s", before)
	}

	// Nachher: trace-Block für die Fremd-Konvention.
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  slices:
    file-pattern: '^(\d+)-.*\.md$'
`)
	_, after, _ := run(t, "--trace", "--json", root)
	var da traceDoc
	if err := json.Unmarshal([]byte(after), &da); err != nil {
		t.Fatalf("JSON (nachher) dekodiert nicht: %v\n%s", err, after)
	}
	if da.Total != 3 || da.Orphans != 1 {
		t.Fatalf("Total/Orphans = %d/%d (erwartet 3/1)\n%s", da.Total, da.Orphans, after)
	}
	arch := traceReqIdx(da, "GG-ARCH-001")
	if arch < 0 || da.Requirements[arch].Orphan ||
		len(da.Requirements[arch].Slices) != 1 || da.Requirements[arch].Slices[0] != "063" ||
		len(da.Requirements[arch].ADRs) != 1 || da.Requirements[arch].ADRs[0] != "ADR-0007" {
		t.Fatalf("GG-ARCH-001 nach Config falsch: %+v\n%s", da.Requirements, after)
	}
	if i := traceReqIdx(da, "GG-QA-001"); i < 0 || da.Requirements[i].Orphan {
		t.Fatalf("GG-QA-001 nach Config sollte via Slice 063 kein Waise sein\n%s", after)
	}
	if i := traceReqIdx(da, "GG-QA-002"); i < 0 || !da.Requirements[i].Orphan {
		t.Fatalf("GG-QA-002 sollte Waise sein\n%s", after)
	}
}

// slice-066 (DC-FA-CLI-009): ALLE acht Override-Achsen mit Nicht-Default-Werten
// (eigene Quell-Datei, Kennungs-Regex, ADR-/Slice-Verzeichnis, -Dateimuster und
// -Owner-Präfix) — mutations-hart: das Ergebnis (Owner-Kennungen mit den
// custom-Präfixen aus den custom-Verzeichnissen) fällt weg, sobald ein einzelner
// resolveTrace-Zweig auf den Default zurückfällt.
func TestCLI066_Trace_VollCustomKonvention(t *testing.T) {
	root := t.TempDir()
	// Custom Anforderungs-Quelle + -Kennung (nicht spec/lastenheft.md, nicht -FA-/-QA-).
	write(t, root, "spec/reqs.md", "### REQ-42 — Eine Anforderung\nText.\n")
	// Custom ADR-Verzeichnis + -Dateimuster (D-NNNN statt NNNN) + Präfix DEC-.
	write(t, root, "decisions/D-0007-x.md", "# Decision\nBezug: REQ-42\n")
	// Custom Slice-Verzeichnis + -Dateimuster (task-N) + Präfix T.
	write(t, root, "work/task-5-y.md", "# Task\nBezug: REQ-42\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    source: spec/reqs.md
    id-pattern: 'REQ-\d+'
  adrs:
    dir: decisions
    file-pattern: '^D-(\d{4})-.*\.md$'
    id-prefix: 'DEC-'
  slices:
    dir: work
    file-pattern: '^task-(\d+)-.*\.md$'
    id-prefix: 'T'
`)
	code, stdout, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON dekodiert nicht: %v\n%s", err, stdout)
	}
	if doc.Total != 1 {
		t.Fatalf("Total = %d (erwartet 1: REQ-42)\n%s", doc.Total, stdout)
	}
	i := traceReqIdx(doc, "REQ-42")
	if i < 0 {
		t.Fatalf("REQ-42 nicht gefunden (custom source/id-pattern greift nicht)\n%s", stdout)
	}
	r := doc.Requirements[i]
	if r.Title != "Eine Anforderung" {
		t.Fatalf("Titel = %q (erwartet 'Eine Anforderung')", r.Title)
	}
	if r.Orphan || len(r.ADRs) != 1 || r.ADRs[0] != "DEC-0007" || len(r.Slices) != 1 || r.Slices[0] != "T5" {
		t.Fatalf("REQ-42 falsch aufgelöst: %+v (erwartet ADRs=[DEC-0007], Slices=[T5], kein Waise)\n%s", r, stdout)
	}
}

// slice-066 Negative (Config): eine ungültige trace-Regex ⇒ Exit 2 (fail-closed).
func TestCLI066_Trace_UngueltigeRegex(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	write(t, root, ".d-check.yml", "trace:\n  requirements:\n    id-pattern: '['\n")
	code, _, stderr := run(t, "--trace", root)
	if code != 2 {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
	if !strings.Contains(stderr, "id-pattern") {
		t.Fatalf("stderr nennt das Feld nicht: %q", stderr)
	}
}

// slice-066 Negative (Config): eine file-pattern ohne Capture-Gruppe ⇒ Exit 2 —
// verhindert das m[1]-Panic der Owner-Ableitung (fail-closed).
func TestCLI066_Trace_FilePatternOhneCapture(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	write(t, root, ".d-check.yml", `trace:
  slices:
    file-pattern: '^slice-\d+-.*\.md$'
`)
	code, _, stderr := run(t, "--trace", root)
	if code != 2 {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
	if !strings.Contains(stderr, "Capture-Gruppe") {
		t.Fatalf("stderr nennt die Capture-Gruppe nicht: %q", stderr)
	}
}

// slice-066: --require-complete (DC-FA-CLI-011) erbt die konfigurierten Quellen —
// eine nur per trace.id-pattern sichtbare Waisen-Anforderung ⇒ Exit 1 (ohne
// Config wäre sie unsichtbar ⇒ Exit 0).
func TestCLI066_Trace_RequireCompleteErbtConfig(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "### GG-ARCH-001\nWaise ohne Slice.\n")

	// Ohne Config: der Default-Regex sieht GG-ARCH-001 nicht ⇒ 0 Waisen ⇒ Exit 0.
	if code, _, _ := run(t, "--trace", "--require-complete", root); code != 0 {
		t.Fatalf("ohne Config: Exit = %d (erwartet 0 — GG-ARCH-001 unsichtbar)", code)
	}

	write(t, root, ".d-check.yml", "trace:\n  requirements:\n    id-pattern: 'GG-[A-Z][A-Z0-9]*-\\d{3}'\n")
	if code, _, stderr := run(t, "--trace", "--require-complete", root); code != 1 {
		t.Fatalf("mit Config: Exit = %d (erwartet 1 — GG-ARCH-001 Waise), stderr = %q", code, stderr)
	}
}

// slice-066 Negative (Config): ein trace-Pfad, der die Repo-Wurzel verlässt
// (führendes '..'), ⇒ Exit 2 (fail-closed, analog planning.roadmap).
func TestCLI066_Trace_PfadEscape(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	write(t, root, ".d-check.yml", "trace:\n  requirements:\n    source: ../outside.md\n")
	code, _, stderr := run(t, "--trace", root)
	if code != 2 {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
	if !strings.Contains(stderr, "relativ zur Repo-Wurzel") {
		t.Fatalf("stderr nennt die Pfad-Regel nicht: %q", stderr)
	}
}

// covRepo: Fremd-Konvention (GG-<FAM>-NNN), keine Slices — Deckung liegt in der
// kuratierten traceability.md (slice-067, DC-FA-COV-001).
func covRepo(t *testing.T, root string) {
	write(t, root, "spec/lastenheft.md",
		"### GG-QA-001\nText.\n\n### GG-QA-002\nText.\n\n### GG-QA-003\nText.\n\n"+
			"### GG-QA-004\nText.\n\n### GG-QA-005\nText.\n\n### GG-QA-006\nText.\n\n### GG-RT-001\nText.\n")
}

// slice-067 (DC-FA-COV-001): eine trace.coverage-Quelle mit Range deckt
// slice-lose Anforderungen ab — Coverage-Spalte erscheint, gedeckte sind keine
// Waisen; die Range `001..006` deckt alle sechs.
func TestCLI067_Coverage_RangeDecktAb(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md", "# Trace\n\nGG-QA-001..006 sind statisch geprüft.\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      ranges: true
`)
	_, md, _ := run(t, "--trace", root)
	if !strings.Contains(md, "| Anforderung | Titel | ADRs | Slices | Coverage | Status |") {
		t.Fatalf("Coverage-Spalte fehlt:\n%s", md)
	}
	code, js, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(js), &doc); err != nil {
		t.Fatalf("JSON: %v\n%s", err, js)
	}
	if doc.Total != 7 || doc.Orphans != 1 {
		t.Fatalf("Total/Orphans = %d/%d (erwartet 7/1)\n%s", doc.Total, doc.Orphans, js)
	}
	for _, id := range []string{"GG-QA-001", "GG-QA-004", "GG-QA-006"} {
		i := traceReqIdx(doc, id)
		if i < 0 || doc.Requirements[i].Orphan ||
			len(doc.Requirements[i].Coverage) != 1 || doc.Requirements[i].Coverage[0] != "Trace" {
			t.Fatalf("%s nicht via Range gedeckt: %+v\n%s", id, doc.Requirements[i], js)
		}
	}
	if i := traceReqIdx(doc, "GG-RT-001"); i < 0 || !doc.Requirements[i].Orphan {
		t.Fatalf("GG-RT-001 (keine Coverage) sollte Waise sein\n%s", js)
	}
}

// slice-067: ohne trace.coverage keine Coverage-Spalte und kein coverage-Feld
// (byte-identisch); mit Quelle erscheinen beide.
func TestCLI067_Coverage_ByteIdentischOhneBlock(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md", "# Trace\n\nGG-QA-001..006.\n")
	write(t, root, ".d-check.yml", "trace:\n  requirements:\n    id-pattern: 'GG-[A-Z][A-Z0-9]*-\\d{3}'\n")
	_, md, _ := run(t, "--trace", root)
	if strings.Contains(md, "Coverage") {
		t.Fatalf("Coverage-Spalte trotz fehlendem Block:\n%s", md)
	}
	_, js, _ := run(t, "--trace", "--json", root)
	if strings.Contains(js, "coverage") {
		t.Fatalf("coverage-Feld trotz fehlendem Block:\n%s", js)
	}
}

// slice-067: exclude-sections nimmt eine „ohne Design-Artefakt"-Sektion aus —
// nur dort genannte IDs werden NICHT kreditiert (voller Heading-Text).
func TestCLI067_Coverage_ExcludeSection(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md",
		"## 27.1 Design\n\nGG-QA-001..003 geprüft.\n\n"+
			"### 27.1.1 Ohne Design\n\nGG-QA-004..006 offen.\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      exclude-sections: ["27.1.1 Ohne Design"]
`)
	code, js, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	_ = json.Unmarshal([]byte(js), &doc)
	if i := traceReqIdx(doc, "GG-QA-002"); i < 0 || doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-002 (§27.1) sollte gedeckt sein\n%s", js)
	}
	if i := traceReqIdx(doc, "GG-QA-005"); i < 0 || !doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-005 (nur §27.1.1, ausgeschlossen) sollte Waise sein\n%s", js)
	}
}

// slice-067 Negative: ein Sektionsname ohne Heading-Treffer (Kurzform statt
// vollem Text) ⇒ Exit 2 (fail-closed Guard gegen stilles Leer-Blanking).
func TestCLI067_Coverage_SektionOhneTreffer(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md", "## 27.1 Design\n\nGG-QA-001..003.\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      sections: ["27.1"]
`)
	code, _, stderr := run(t, "--trace", root)
	if code != 2 || !strings.Contains(stderr, "Abschnitt") {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
}

// slice-067 Negative: fehlende Coverage-Datei ⇒ Exit 2 (files sind explizit
// benannt, Fehlen ist Fehler, nicht Skip).
func TestCLI067_Coverage_FehlendeDatei(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/nicht-da.md]
      label: Trace
`)
	code, _, stderr := run(t, "--trace", root)
	if code != 2 || !strings.Contains(stderr, "fehlt") {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
}

// slice-067 Negative: eine ungültige Range (AAA>BBB) in der Coverage-Quelle
// ⇒ Exit 2.
func TestCLI067_Coverage_UngueltigeRange(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md", "# Trace\n\nGG-QA-009..003 (falsch).\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      ranges: true
`)
	code, _, stderr := run(t, "--trace", root)
	if code != 2 || !strings.Contains(stderr, "AAA>BBB") {
		t.Fatalf("Exit = %d (erwartet 2), stderr = %q", code, stderr)
	}
}

// slice-067: --require-complete berücksichtigt Coverage — sind alle
// Anforderungen via Slice ODER Coverage gedeckt, Exit 0 (ohne Coverage Exit 1).
func TestCLI067_Coverage_RequireComplete(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "### GG-QA-001\nText.\n\n### GG-QA-002\nText.\n")
	write(t, root, "docs/plan/traceability.md", "# Trace\n\nGG-QA-001..002 geprüft.\n")

	// Ohne Coverage-Block: beide Waisen ⇒ Exit 1.
	write(t, root, ".d-check.yml", "trace:\n  requirements:\n    id-pattern: 'GG-[A-Z][A-Z0-9]*-\\d{3}'\n")
	if code, _, _ := run(t, "--trace", "--require-complete", root); code != 1 {
		t.Fatalf("ohne Coverage: Exit = %d (erwartet 1)", code)
	}

	// Mit Coverage-Block: beide gedeckt ⇒ Exit 0.
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      ranges: true
`)
	if code, _, stderr := run(t, "--trace", "--require-complete", root); code != 0 {
		t.Fatalf("mit Coverage: Exit = %d (erwartet 0), stderr = %q", code, stderr)
	}
}

// slice-067 (R1-MEDIUM): die positive sections-Whitelist zählt NUR die gelisteten
// Abschnitte — IDs außerhalb werden nicht kreditiert (Include-Zweig SelectSections).
func TestCLI067_Coverage_IncludeSection(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md",
		"## 27.1 Design\n\nGG-QA-001..003 nur Design.\n\n"+
			"## 27.2 Impl\n\nGG-QA-004..006 implementiert.\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      sections: ["27.2 Impl"]
`)
	code, js, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	_ = json.Unmarshal([]byte(js), &doc)
	if i := traceReqIdx(doc, "GG-QA-005"); i < 0 || doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-005 (§27.2, Whitelist) sollte gedeckt sein\n%s", js)
	}
	if i := traceReqIdx(doc, "GG-QA-002"); i < 0 || !doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-002 (§27.1, NICHT in Whitelist) sollte Waise sein\n%s", js)
	}
}

// slice-067 (R1-LOW): ranges:false expandiert keine `..`-Range end-to-end —
// nur die exakte Kennung wird gedeckt (Pointer-Default-Auflösung im Decode).
func TestCLI067_Coverage_RangesFalse(t *testing.T) {
	root := t.TempDir()
	covRepo(t, root)
	write(t, root, "docs/plan/traceability.md", "# Trace\n\nGG-QA-001..006.\n")
	write(t, root, ".d-check.yml", `trace:
  requirements:
    id-pattern: 'GG-[A-Z][A-Z0-9]*-\d{3}'
  coverage:
    - files: [docs/plan/traceability.md]
      label: Trace
      ranges: false
`)
	code, js, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	_ = json.Unmarshal([]byte(js), &doc)
	if i := traceReqIdx(doc, "GG-QA-001"); i < 0 || doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-001 (exakt) sollte gedeckt sein\n%s", js)
	}
	if i := traceReqIdx(doc, "GG-QA-003"); i < 0 || !doc.Requirements[i].Orphan {
		t.Fatalf("GG-QA-003 sollte ohne Range-Expansion Waise sein\n%s", js)
	}
}

// slice-038 Happy: --print-mk gibt ein include-bares d-check.mk auf stdout
// (version-gepinntes Image + doc-check-Target), Exit 0; repo-frei (read-only).
func TestCLI038_PrintMK(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	for _, want := range []string{
		"DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v", // version-gepinnter, überschreibbarer Ref
		"DCHECK_DIGEST ?=",                         // Digest-Override-Variable
		"\ndoc-check: ## ",                         // ##-annotiertes Target
		"--network none",
		":/repo:ro",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("d-check.mk ohne %q:\n%s", want, stdout)
		}
	}
}

// slice-038 Negative: --print-mk mit unbekanntem Flag → Nutzungsfehler (Exit 2).
func TestCLI038_PrintMK_UnbekanntesFlag(t *testing.T) {
	code, _, stderr := run(t, "--print-mk", "--bogus")
	if code != 2 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// slice-044 Happy (DC-FA-CLI-011): --trace --require-complete bei 0 Waisen →
// Exit 0; die advisory RTM bleibt unverändert auf stdout.
func TestCLI044_RequireComplete_KeineWaisen(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", "### DC-FA-TRC-001 — Referenzierte Anforderung\nText.\n")
	write(t, root, "docs/plan/planning/done/slice-099-x.md", "# slice-099\nBezug: DC-FA-TRC-001\n")
	code, _, stderr := run(t, "--trace", "--require-complete", root)
	if code != 0 {
		t.Fatalf("Exit = %d (erwartet 0, keine Waisen), stderr = %q", code, stderr)
	}
}

// slice-044 Boundary (DC-FA-CLI-011): --trace --require-complete bei einer Waise
// → Exit 1 (Befund-Code, DC-FA-CLI-003); die RTM steht weiter vollständig auf
// stdout, und es wird nichts geschrieben (DC-QA-03).
func TestCLI044_RequireComplete_Waise(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root) // DC-FA-TRC-002 ist Waise
	code, stdout, stderr := run(t, "--trace", "--require-complete", root)
	if code != 1 {
		t.Fatalf("Exit = %d (erwartet 1 bei Waise), stderr = %q", code, stderr)
	}
	for _, want := range []string{"DC-FA-TRC-002", "WAISE"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("RTM fehlt auf stdout trotz Strict-Exit (%q):\n%s", want, stdout)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".d-check.yml")); !os.IsNotExist(err) {
		t.Fatalf("--require-complete hat geschrieben — read-only verletzt (DC-QA-03)")
	}
}

// slice-044 (Konsumenten-Gate-Pfad, doc-complete mit TRACE_FLAGS=--json):
// --trace --require-complete --json → Exit 1, RTM-JSON intakt (orphans == 1).
func TestCLI044_RequireComplete_JSON(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, stdout, _ := run(t, "--trace", "--require-complete", "--json", root)
	if code != 1 {
		t.Fatalf("Exit = %d (erwartet 1)", code)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON dekodiert nicht: %v\n%s", err, stdout)
	}
	if doc.Orphans != 1 {
		t.Fatalf("Orphans = %d (erwartet 1)\n%s", doc.Orphans, stdout)
	}
}

// slice-044 Negative (DC-FA-CLI-011): --require-complete ohne --trace ist ein
// Nutzungsfehler (Exit 2), keine RTM.
func TestCLI044_RequireComplete_OhneTrace(t *testing.T) {
	root := t.TempDir()
	traceRepo(t, root)
	code, _, stderr := run(t, "--require-complete", root)
	if code != 2 || !strings.Contains(stderr, "--require-complete erfordert --trace") {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
}

// slice-044 (DC-FA-CLI-010-Erweiterung): das --print-mk-Fragment trägt zusätzlich
// die TRACE_FLAGS-Variable und die Targets doc-trace (advisory) sowie
// doc-complete (--trace --require-complete).
func TestCLI044_PrintMK_TraceTargets(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	for _, want := range []string{
		"TRACE_FLAGS ?=",
		"\ndoc-trace: ## ",
		"--trace $(TRACE_FLAGS)",
		"\ndoc-complete: ## ",
		"--trace --require-complete $(TRACE_FLAGS)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("d-check.mk ohne %q:\n%s", want, stdout)
		}
	}
}

// slice-047: --print-mk-Fragment um doc-doctor/doc-repair/doc-help + die
// DCHECK_DIGEST-Variable (Digest-Override per ifeq, sticht den Tag) erweitert;
// doc-repair unterdrückt das Recipe-Echo (@), damit stdout git-apply-rein bleibt.
func TestCLI047_PrintMK_NeueTargets(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	for _, want := range []string{
		"DCHECK_DIGEST ?=",
		"ifeq ($(strip $(DCHECK_DIGEST)),)",
		"DCHECK_REF := $(DCHECK_IMAGE)",
		"d-check@$(DCHECK_DIGEST)",
		"\ndoc-doctor: ## ",
		"$(DCHECK_REF) --doctor",
		"\ndoc-repair: ## ",
		"\t@docker run", // doc-repair: Recipe-Echo unterdrückt (stdout-rein)
		"$(DCHECK_REF) --repair",
		"\ndoc-help: ## ",
		"$(MAKEFILE_LIST)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("d-check.mk ohne %q:\n%s", want, stdout)
		}
	}
}

// slice-053 (DC-FA-CLI-010-Erweiterung): das --print-mk-Fragment trägt zusätzlich
// das Target doc-immutable, das dem Konsumenten die git-Diff-Immutabilität (Modul
// vcs, DC-FA-VCS-001) verteilt — nur vcs (alle übrigen Module via --disable
// abgewählt, aus dem Modulsatz abgeleitet), RANGE/STAGED-getrieben.
// mkTargetRecipe liefert die (einzelne) Recipe-Zeile des --print-mk-Fragments, die
// marker enthält — nötig, seit mehrere Fokus-Targets (doc-immutable/doc-commits)
// je eigene `--disable`-Listen tragen: eine Negativ-Prüfung („X nicht abgewählt")
// muss auf DIE Recipe-Zeile schauen, nicht auf das ganze Fragment.
func mkTargetRecipe(t *testing.T, stdout, marker string) string {
	t.Helper()
	for _, ln := range strings.Split(stdout, "\n") {
		if strings.Contains(ln, marker) {
			return ln
		}
	}
	t.Fatalf("keine Recipe-Zeile mit %q:\n%s", marker, stdout)
	return ""
}

func TestCLI053_PrintMK_ImmutableTarget(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "\ndoc-immutable: ## ") {
		t.Fatalf("d-check.mk ohne doc-immutable-Target:\n%s", stdout)
	}
	line := mkTargetRecipe(t, stdout, "--enable vcs ")
	for _, want := range []string{
		"--disable links",     // Fokus auf vcs (alle Doc-Module abgewählt)
		"--disable immutable", // auch das hermetische Schwester-Modul
		"$(if $(STAGED),--staged,--range $(RANGE))",
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("doc-immutable-Recipe ohne %q:\n%s", want, line)
		}
	}
	// doc-immutable darf vcs NICHT abwählen (in SEINER Recipe-Zeile).
	if strings.Contains(line, "--disable vcs") {
		t.Fatalf("doc-immutable wählt vcs fälschlich ab:\n%s", line)
	}
}

// slice-056 (DC-FA-CLI-010-Erweiterung): das --print-mk-Fragment trägt zusätzlich
// das Target doc-commits, das dem Konsumenten die Commit-Message-Traceability
// (Modul commits, DC-FA-COMMITS-001) verteilt — nur commits (alle übrigen Module
// via --disable abgewählt, aus dem Modulsatz abgeleitet), RANGE-getrieben.
func TestCLI056_PrintMK_CommitsTarget(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "\ndoc-commits: ## ") {
		t.Fatalf("d-check.mk ohne doc-commits-Target:\n%s", stdout)
	}
	line := mkTargetRecipe(t, stdout, "--enable commits ")
	for _, want := range []string{
		"--disable links", // Fokus auf commits (alle Doc-Module abgewählt)
		"--disable vcs",   // auch das git-Schwester-Modul
		"--range $(RANGE)",
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("doc-commits-Recipe ohne %q:\n%s", want, line)
		}
	}
	// doc-commits darf commits NICHT abwählen (in SEINER Recipe-Zeile).
	if strings.Contains(line, "--disable commits") {
		t.Fatalf("doc-commits wählt commits fälschlich ab:\n%s", line)
	}
}

// slice-057 (DC-FA-CLI-010-Erweiterung): das --print-mk-Fragment trägt zusätzlich
// das Target doc-planning (Modul planning, DC-FA-PLAN-001) — **hermetisch**, ohne
// --range (nur der Arbeitsbaum), nur planning (alle übrigen Module abgewählt).
func TestCLI057_PrintMK_PlanningTarget(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "\ndoc-planning: ## ") {
		t.Fatalf("d-check.mk ohne doc-planning-Target:\n%s", stdout)
	}
	line := mkTargetRecipe(t, stdout, "--enable planning ")
	for _, want := range []string{
		"--disable links",   // Fokus auf planning (alle übrigen abgewählt)
		"--disable commits", // auch das commits-Schwester-Modul
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("doc-planning-Recipe ohne %q:\n%s", want, line)
		}
	}
	// doc-planning ist hermetisch: KEINE Range, und planning NICHT abgewählt.
	if strings.Contains(line, "--range") || strings.Contains(line, "--staged") {
		t.Fatalf("doc-planning trägt fälschlich eine Range:\n%s", line)
	}
	if strings.Contains(line, "--disable planning") {
		t.Fatalf("doc-planning wählt planning fälschlich ab:\n%s", line)
	}
}

// slice-063 (DC-FA-CLI-010-Erweiterung, 10→11): das --print-mk-Fragment trägt
// zusätzlich das Target doc-targets (Modul targets, DC-FA-TGT-001) —
// **hermetisch**, ohne --range (nur Makefile + Doku-Tabellen im Mount), nur
// targets (alle übrigen Module abgewählt).
func TestCLI063_PrintMK_TargetsTarget(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "\ndoc-targets: ## ") {
		t.Fatalf("d-check.mk ohne doc-targets-Target:\n%s", stdout)
	}
	line := mkTargetRecipe(t, stdout, "--enable targets ")
	for _, want := range []string{
		"--disable links",    // Fokus auf targets (alle übrigen abgewählt)
		"--disable planning", // auch das hermetische Schwester-Modul
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("doc-targets-Recipe ohne %q:\n%s", want, line)
		}
	}
	// doc-targets ist hermetisch: KEINE Range, und targets NICHT abgewählt.
	if strings.Contains(line, "--range") || strings.Contains(line, "--staged") {
		t.Fatalf("doc-targets trägt fälschlich eine Range:\n%s", line)
	}
	if strings.Contains(line, "--disable targets") {
		t.Fatalf("doc-targets wählt targets fälschlich ab:\n%s", line)
	}
}
