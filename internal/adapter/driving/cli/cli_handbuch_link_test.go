package cli_test

import (
	"strings"
	"testing"
)

// Die zwei Handbuch-URLs, wie sie in beiden Ausgaben stehen muessen (seit
// slice-187: `blob` fuer den klickenden Menschen, `raw` fuer den ziehenden
// Agenten). Sie stehen hier LITERAL und nicht als Import aus dem
// Produktionscode: ein Test, der die Konstante des Prueflings
// wiederverwendet, prueft nur, dass eine Variable gleich sich selbst ist —
// er ueberlebt jede Aenderung ihres Wertes.
const (
	wantHandbuchURL    = "https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md"
	wantHandbuchURLRaw = "https://raw.githubusercontent.com/pt9912/d-check/refs/heads/main/docs/user/benutzerhandbuch.md"
)

// DC-FA-CLI-001: die Hilfe nennt das Benutzerhandbuch, in BEIDEN Formen. Wer
// d-check aus dem Werkzeug heraus kennenlernt, findet den
// aufgabenorientierten Einstieg sonst nicht — beide Ausgaben verwiesen vorher
// ausschliesslich auf ANDERE Ausgaben, und vor slice-187 fehlte die rohe Form
// ganz (der Hauptleser ist ein Code-Agent, kein Mensch im Terminal).
func TestCLI001_HilfeNenntHandbuch(t *testing.T) {
	code, _, stderr := run(t, "--help")
	if code != 0 {
		t.Fatalf("Exit = %d, want 0", code)
	}
	if !strings.Contains(stderr, wantHandbuchURL) {
		t.Fatalf("Hilfe ohne Handbuch-URL (blob):\n%s", stderr)
	}
	if !strings.Contains(stderr, wantHandbuchURLRaw) {
		t.Fatalf("Hilfe ohne Handbuch-URL (raw):\n%s", stderr)
	}
	if !strings.Contains(stderr, "Benutzerhandbuch") {
		t.Fatalf("Hilfe ohne Ueberschrift zum Handbuch-Zeiger:\n%s", stderr)
	}
}

// DC-FA-CLI-010: dasselbe im Kopf des erzeugten Fragments, ebenfalls beide
// Formen. Es reist in ein FREMDES Repo — dort ist der Kopf der einzige Ort,
// an dem ein Zeiger dauerhaft mitfaehrt.
func TestCLI010_PrintMKNenntHandbuch(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, wantHandbuchURL) {
		t.Fatalf("--print-mk ohne Handbuch-URL (blob):\n%s", stdout)
	}
	if !strings.Contains(stdout, wantHandbuchURLRaw) {
		t.Fatalf("--print-mk ohne Handbuch-URL (raw):\n%s", stdout)
	}
	// Im Kopf, nicht irgendwo: beide Zeilen stehen VOR der ersten Variablen-
	// ZUWEISUNG, sonst stuenden sie zwischen Targets statt in der Einleitung.
	// Der Anker ist "\nDCHECK_IMAGE ?=", nicht der blosse Name — der kommt
	// schon im Digest-Hinweis darueber vor.
	//
	// FEHLT DER ANKER, BRICHT DER TEST, statt auf die schwaechere Pruefung
	// oben zurueckzufallen: ein Waechter, der still seine Haelfte verliert,
	// liest sich wie einer, der sie noch haelt (BEO-023).
	i := strings.Index(stdout, "\nDCHECK_IMAGE ?=")
	if i <= 0 {
		t.Fatalf("Anker \"DCHECK_IMAGE ?=\" fehlt — die Kopf-Zusage ist nicht mehr pruefbar:\n%s", stdout)
	}
	kopf := stdout[:i]
	if !strings.Contains(kopf, wantHandbuchURL) {
		t.Fatalf("Handbuch-URL (blob) steht nicht im Kopf des Fragments:\n%s", kopf)
	}
	if !strings.Contains(kopf, wantHandbuchURLRaw) {
		t.Fatalf("Handbuch-URL (raw) steht nicht im Kopf des Fragments:\n%s", kopf)
	}
}

// KEINE DER BEIDEN URLS DARF EINE VERSION TRAGEN. Dieser Test ist der
// Waechter ueber genau diese Zusage — er faellt, sobald jemand
// `blob/v0.68.0` oder eine versionierte `raw`-Form schreibt.
//
// GRENZE: er ist fail-closed und damit ein OBERMENGEN-Waechter — er faellt
// auch, wenn eine Zeile ganz fehlt. Die scharfe Probe fuer ihn allein ist
// eine Versionierung in BEIDEN Literalen je Form, Produktion und Test.
func TestHandbuchURL_TraegtKeineVersion(t *testing.T) {
	_, _, hilfe := run(t, "--help")
	_, fragment, _ := run(t, "--print-mk")
	for name, ausgabe := range map[string]string{"--help": hilfe, "--print-mk": fragment} {
		zeilen := zeilenMit(ausgabe, "benutzerhandbuch.md")
		if len(zeilen) != 2 {
			t.Fatalf("%s: %d Handbuch-Zeilen gefunden, will 2:\n%v", name, len(zeilen), zeilen)
		}
		var sahBlob, sahRaw bool
		for _, zeile := range zeilen {
			switch {
			case strings.Contains(zeile, "/blob/"):
				sahBlob = true
				if !strings.Contains(zeile, "/blob/main/") {
					t.Errorf("%s: blob-URL zeigt nicht auf main: %q", name, zeile)
				}
				if strings.Contains(zeile, "/blob/v") {
					t.Errorf("%s: blob-URL traegt eine Version: %q", name, zeile)
				}
			case strings.Contains(zeile, "raw.githubusercontent.com"):
				sahRaw = true
				if !strings.Contains(zeile, "/refs/heads/main/") {
					t.Errorf("%s: raw-URL zeigt nicht auf main: %q", name, zeile)
				}
				if strings.Contains(zeile, "/refs/heads/v") || strings.Contains(zeile, "/refs/tags/") {
					t.Errorf("%s: raw-URL traegt eine Version: %q", name, zeile)
				}
			default:
				t.Errorf("%s: unerwartete Handbuch-Zeile ohne blob/raw-Form: %q", name, zeile)
			}
		}
		if !sahBlob || !sahRaw {
			t.Errorf("%s: nicht beide Formen gesehen (blob=%v, raw=%v)", name, sahBlob, sahRaw)
		}
	}
}

// zeilenMit liefert alle Zeilen, die needle enthalten (leer, wenn keine).
func zeilenMit(text, needle string) []string {
	var out []string
	for _, l := range strings.Split(text, "\n") {
		if strings.Contains(l, needle) {
			out = append(out, strings.TrimSpace(l))
		}
	}
	return out
}
