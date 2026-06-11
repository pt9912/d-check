// Akzeptanztests des CLI gegen die Lastenheft-Kriterien
// (Happy/Boundary/Negative) — echtes Dateisystem via t.TempDir
// (u-boot-Konvention: package cli_test).
package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	if code != 0 || stdout != "" || !strings.Contains(stderr, "-json") {
		t.Fatalf("-h: Exit = %d, stdout = %q, stderr = %q", code, stdout, stderr)
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
	if strings.Contains(stderr, "nicht implementiert") {
		t.Fatalf("Interim-Hinweis darf für ids nicht mehr erscheinen: %q", stderr)
	}

	write(t, root, "docs/a.md", "nacktes ADR-0042 im Fließtext\n")
	code, stdout, _ := run(t, root)
	if code != 1 || !strings.Contains(stdout, "docs/a.md:1\tADR-0042\tid-unlinked") {
		t.Fatalf("Exit = %d, stdout = %q", code, stdout)
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
