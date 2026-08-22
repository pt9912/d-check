package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// versionsRepo legt ein Repo mit einem Release-Register, einem Baseline-Pin
// und einem Dokument an, das Pins beider Reihen trägt.
func versionsRepo(t *testing.T, cfg string) string {
	t.Helper()
	dir := t.TempDir()
	writeAt(t, dir, ".d-check.yml", cfg)
	writeAt(t, dir, "version.md", "# Register\n\n## Aktuell\n\nAktuelle Version: `v0.27.0`.\n")
	writeAt(t, dir, "pin.md", "# Konventionen\n\n## Baseline\n\nGepinnt auf `v5.9.0`.\n")
	writeAt(t, dir, "doc.md", "# D\n\n```bash\ndocker run ghcr.io/x/d-check:v0.26.0\n```\n\nBaum: baseline/v5.7.0\n")
	return dir
}

const versionsKurzform = "versions:\n" +
	"  pin-pattern: 'ghcr\\.io/[^\\s:]+:(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n" +
	"  current-from: version.md#aktuell\n"

const versionsEinPaarListe = "versions:\n" +
	"  patterns:\n" +
	"    - pin-pattern: 'ghcr\\.io/[^\\s:]+:(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n" +
	"      current-from: version.md#aktuell\n"

// runVersions läuft mit --json über das Repo und liefert Exit-Code und die
// Roh-Ausgabe (für den byteweisen Vergleich).
func runVersions(t *testing.T, dir string) (int, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "versions", "--json", dir}, &stdout, &stderr)
	if code == 2 {
		t.Fatalf("Nutzungsfehler: %s", stderr.String())
	}
	return code, stdout.String()
}

// DC-FA-VER-001 / DC-QA-02 (E2E): die Kurzform IST die einelementige
// patterns-Liste — beide Schreibweisen liefern eine BYTE-identische Ausgabe.
// Der Pfad-Anteil der JSON-Ausgabe ist repo-relativ, die beiden TempDirs
// unterscheiden sich also nicht in der Ausgabe.
func TestVersionsPatterns_KurzformByteIdentisch(t *testing.T) {
	codeK, outK := runVersions(t, versionsRepo(t, versionsKurzform))
	codeL, outL := runVersions(t, versionsRepo(t, versionsEinPaarListe))
	if codeK != 1 {
		t.Fatalf("Vorbedingung: ein veralteter Pin ⇒ Exit 1, got %d\n%s", codeK, outK)
	}
	if codeK != codeL || outK != outL {
		t.Fatalf("Kurzform und Ein-Paar-Liste müssen byte-identisch sein:\nExit %d/%d\n%s\n---\n%s", codeK, codeL, outK, outL)
	}
	if !strings.Contains(outK, "v0.26.0") || strings.Contains(outK, "v5.7.0") {
		t.Fatalf("ein Paar prüft nur seine Reihe: %s", outK)
	}
}

// DC-FA-VER-001 (E2E): zwei Paare prüfen zwei unabhängige Reihen gegen zwei
// Quellen — beide veralteten Pins werden gemeldet, je gegen ihre eigene
// erwartete Version.
func TestVersionsPatterns_ZweiReihen(t *testing.T) {
	dir := versionsRepo(t, "versions:\n"+
		"  patterns:\n"+
		"    - pin-pattern: 'ghcr\\.io/[^\\s:]+:(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n"+
		"      current-from: version.md#aktuell\n"+
		"    - pin-pattern: 'baseline/(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n"+
		"      current-from: pin.md#baseline\n")
	code, out := runVersions(t, dir)
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\n%s", code, out)
	}
	for _, want := range []string{"v0.26.0", "erwartet v0.27.0", "v5.7.0", "erwartet v5.9.0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("Ausgabe ohne %q: %s", want, out)
		}
	}
}

// DC-FA-VER-001 fail-closed (E2E): Kurzform und patterns zugleich ⇒ Exit 2
// ohne Prüfung.
func TestVersionsPatterns_MischformExit2(t *testing.T) {
	dir := versionsRepo(t, versionsKurzform+
		"  patterns:\n"+
		"    - pin-pattern: 'baseline/(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n")
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "versions", dir}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("Exit = %d, want 2\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "zugleich gesetzt") {
		t.Fatalf("stderr benennt die Mischform nicht: %s", stderr.String())
	}
}

// uncommentBlock nimmt aus der --print-config-Vorlage den Block, dessen
// Kopfzeile mit header beginnt, und entfernt je Zeile genau ein führendes
// "#" samt einem Leerzeichen — so, wie ein Nutzer einen Vorlagen-Block
// einkommentiert.
func uncommentBlock(t *testing.T, template, header string) string {
	t.Helper()
	lines := strings.Split(template, "\n")
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, header) {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatalf("Vorlagen-Block %q nicht gefunden", header)
	}
	var out []string
	for _, l := range lines[start:] {
		if !strings.HasPrefix(l, "#") {
			break
		}
		l = strings.TrimPrefix(strings.TrimPrefix(l, "#"), " ")
		out = append(out, l)
	}
	return strings.Join(out, "\n") + "\n"
}

// DC-FA-CLI-005 / DC-FA-VER-001: jeder der beiden versions-Blöcke der Vorlage
// ist für sich einkommentierbar und wird vom eigenen Parser angenommen —
// beide zusammen sind der benannte Nutzungsfehler.
func TestVersionsPatterns_VorlagenBloeckeEinkommentierbar(t *testing.T) {
	code, template, stderr := run(t, "--print-config")
	if code != 0 {
		t.Fatalf("--print-config: Exit = %d, stderr = %q", code, stderr)
	}
	kurz := uncommentBlock(t, template, "# --- versions:")
	liste := uncommentBlock(t, template, "# --- versions, ALTERNATIVE")
	for name, block := range map[string]string{"kurzform": kurz, "patterns": liste} {
		if !strings.Contains(block, "versions:") {
			t.Fatalf("%s: Block ohne versions-Schlüssel:\n%s", name, block)
		}
		if _, err := configyaml.Decode([]byte(block)); err != nil {
			t.Fatalf("%s: eigener Parser lehnt den einkommentierten Block ab: %v\n%s", name, err, block)
		}
	}
	// Beide Blöcke zugleich sind kein gültiger Zustand — hier fängt es schon
	// das YAML (zwei `versions:`-Schlüssel); die Mischform INNERHALB eines
	// Blocks fängt der Config-Rand (configyaml-Tests).
	if _, err := configyaml.Decode([]byte(kurz + liste)); err == nil {
		t.Fatal("beide Blöcke zugleich müssen abgewiesen werden")
	}
}
