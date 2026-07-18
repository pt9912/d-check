package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

// DC-FA-REF-001 (E2E): das geteilte ignore-refs-Ventil über die reale
// .d-check.yml. Ein Template-Verzeichnis mischt einen Ziel-Repo-Platzhalter
// (soll ignoriert werden) mit zwei per keep zurückgeholten Verweisen — einer
// real (grün), einer verfälscht (ERROR). Erwartet: Exit 1, der Tippfehler
// gemeldet, der Platzhalter NICHT. Ohne keep-Support wäre real-typo.md still
// (CR nicht abgenommen); der Tippfehler-Test verriegelt „Muster statt
// Heuristik".
func TestRefs_KeepUndTippfehlerEndToEnd(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, ".d-check.yml",
		"ignore-refs:\n"+
			"  - in: \"docs/tpl/**\"\n"+
			"    refs: [\"docs/tpl/**\"]\n"+
			"    keep: [\"docs/tpl/real-*\"]\n")
	write(t, dir, "docs/tpl/a.md",
		"Platzhalter: [p](ziel-repo.md) — real: [r](real-ok.md) — Tippfehler: [t](real-tippfehler.md)\n")
	write(t, dir, "docs/tpl/real-ok.md", "# Da\n")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "links", dir}, &stdout, &stderr)
	out := stdout.String()
	if code != 1 {
		t.Fatalf("Exit = %d, want 1\nstdout=%s\nstderr=%s", code, out, stderr.String())
	}
	if strings.Contains(out, "ziel-repo.md") {
		t.Fatalf("Platzhalter ziel-repo.md wurde gemeldet (refs sollte ihn ignorieren):\n%s", out)
	}
	if !strings.Contains(out, "real-tippfehler.md") {
		t.Fatalf("Tippfehler real-tippfehler.md fehlt im Befund (keep sollte ihn scharf lassen):\n%s", out)
	}
	if strings.Contains(out, "real-ok.md") {
		t.Fatalf("real-ok.md (existiert, per keep geprüft) sollte grün sein:\n%s", out)
	}
}

// DC-FA-REF-001 / DC-FA-CONF-001 Negative: ein ungültiges Glob in
// in/refs/keep ⇒ Exit 2 (config-zeitig, fail-closed — kein still
// wirkungsloses Ventil). Segmentweise Validierung wie die Laufzeit.
func TestRefs_UngueltigesGlobExit2(t *testing.T) {
	cases := []string{
		"ignore-refs:\n  - refs: [\"[\"]\n",
		"ignore-refs:\n  - in: \"[a/b]\"\n    refs: [\"docs/**\"]\n",
		"ignore-refs:\n  - refs: [\"docs/**\"]\n    keep: [\"\"]\n",
	}
	for _, yaml := range cases {
		dir := t.TempDir()
		write(t, dir, ".d-check.yml", yaml)
		write(t, dir, "docs/a.md", "x\n")

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"--enable", "links", dir}, &stdout, &stderr)
		if code != 2 || !strings.Contains(stderr.String(), "ignore-refs") {
			t.Fatalf("YAML %q: Exit = %d, stderr = %q, want 2 + ignore-refs-Hinweis", yaml, code, stderr.String())
		}
	}
}

// DC-FA-REF-001 / ADR-0044 Entsch. 7 — die Symlink-Ablehnung
// (DC-FA-LINK-002) bleibt unberührt: ein Symlink-Ziel, das ein ignore-refs-Glob
// matcht, erzeugt weiter einen `symlink`-Befund, weil das Ventil NACH der
// Symlink-Prüfung greift. Pinnt die Reihenfolge in links.go.
func TestRefs_SymlinkBleibtTrotzVentil(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "docs/real.md", "# Real\n")
	write(t, dir, "docs/a.md", "[x](slink.md)\n")
	if err := os.Symlink("real.md", filepath.Join(dir, "docs", "slink.md")); err != nil {
		t.Skip("Symlink nicht unterstützt: " + err.Error())
	}
	write(t, dir, ".d-check.yml", "ignore-refs:\n  - refs: [\"docs/slink.md\"]\n")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "links", dir}, &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "symlink") {
		t.Fatalf("Symlink trotz ignore-refs: Exit = %d, want 1 + symlink-Befund\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
}

// DC-FA-REF-001 Regression: ohne ignore-refs-Block ist das Verhalten
// byte-identisch (DC-QA-02) — ein fehlendes Ziel bleibt target-missing.
func TestRefs_DefaultAusByteIdentisch(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "docs/tpl/a.md", "[p](ziel-repo.md)\n")

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"--enable", "links", dir}, &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "ziel-repo.md") {
		t.Fatalf("ohne ignore-refs: Exit = %d, want 1 + target-missing für ziel-repo.md\nstdout=%s", code, stdout.String())
	}
}
