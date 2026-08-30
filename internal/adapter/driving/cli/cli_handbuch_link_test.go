package cli_test

import (
	"strings"
	"testing"
)

// Die Handbuch-URL, wie sie in beiden Ausgaben stehen muss. Sie steht hier
// LITERAL und nicht als Import aus dem Produktionscode: ein Test, der die
// Konstante des Prueflings wiederverwendet, prueft nur, dass eine Variable
// gleich sich selbst ist — er ueberlebt jede Aenderung ihres Wertes.
const wantHandbuchURL = "https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md"

// DC-FA-CLI-001: die Hilfe nennt das Benutzerhandbuch. Wer d-check aus dem
// Werkzeug heraus kennenlernt, findet den aufgabenorientierten Einstieg sonst
// nicht — beide Ausgaben verwiesen vorher ausschliesslich auf ANDERE Ausgaben.
func TestCLI001_HilfeNenntHandbuch(t *testing.T) {
	code, _, stderr := run(t, "--help")
	if code != 0 {
		t.Fatalf("Exit = %d, want 0", code)
	}
	if !strings.Contains(stderr, wantHandbuchURL) {
		t.Fatalf("Hilfe ohne Handbuch-URL:\n%s", stderr)
	}
	if !strings.Contains(stderr, "Benutzerhandbuch") {
		t.Fatalf("Hilfe ohne Ueberschrift zum Handbuch-Zeiger:\n%s", stderr)
	}
}

// DC-FA-CLI-010: dasselbe im Kopf des erzeugten Fragments. Es reist in ein
// FREMDES Repo — dort ist der Kopf der einzige Ort, an dem ein Zeiger
// dauerhaft mitfaehrt.
func TestCLI010_PrintMKNenntHandbuch(t *testing.T) {
	code, stdout, stderr := run(t, "--print-mk")
	if code != 0 {
		t.Fatalf("Exit = %d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, wantHandbuchURL) {
		t.Fatalf("--print-mk ohne Handbuch-URL:\n%s", stdout)
	}
	// Im Kopf, nicht irgendwo: die Zeile steht VOR der ersten Variablen-
	// ZUWEISUNG, sonst stuende sie zwischen Targets statt in der Einleitung.
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
	if kopf := stdout[:i]; !strings.Contains(kopf, wantHandbuchURL) {
		t.Fatalf("Handbuch-URL steht nicht im Kopf des Fragments:\n%s", kopf)
	}
}

// DIE URL DARF KEINE VERSION TRAGEN. Dieser Test ist der Waechter ueber
// genau diese Zusage — er faellt, sobald jemand `blob/v0.68.0` schreibt.
//
// GRENZE: er ist fail-closed und damit ein OBERMENGEN-Waechter — er faellt
// auch, wenn die Zeile ganz fehlt. Die scharfe Probe fuer ihn allein ist
// eine Versionierung in BEIDEN Literalen, Produktion und Test.
func TestHandbuchURL_TraegtKeineVersion(t *testing.T) {
	_, _, hilfe := run(t, "--help")
	_, fragment, _ := run(t, "--print-mk")
	for name, ausgabe := range map[string]string{"--help": hilfe, "--print-mk": fragment} {
		zeile := zeileMit(ausgabe, "benutzerhandbuch.md")
		if zeile == "" {
			t.Fatalf("%s: keine Handbuch-Zeile gefunden", name)
		}
		if !strings.Contains(zeile, "/blob/main/") {
			t.Errorf("%s: Handbuch-URL zeigt nicht auf main: %q", name, zeile)
		}
		if strings.Contains(zeile, "/blob/v") || strings.Contains(zeile, "@v") {
			t.Errorf("%s: Handbuch-URL traegt eine Version: %q", name, zeile)
		}
	}
}

// zeileMit liefert die erste Zeile, die needle enthaelt (leer, wenn keine).
func zeileMit(text, needle string) string {
	for _, l := range strings.Split(text, "\n") {
		if strings.Contains(l, needle) {
			return strings.TrimSpace(l)
		}
	}
	return ""
}
