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
	"strings"
	"testing"

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
	for _, want := range []string{"Baseline v1.3.0", "matrix:", "from: spec-straten, to: adr", "exclude-sections", "Carveouts"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Vorlage ohne %q:\n%s", want, stdout)
		}
	}
	// read-only: das Werkzeug schreibt nichts (Aufrufer leitet um).
	if _, err := os.Stat(filepath.Join(root, ".d-check.yml")); !os.IsNotExist(err) {
		t.Fatalf(".d-check.yml wurde geschrieben — read-only verletzt (DC-QA-03)")
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
