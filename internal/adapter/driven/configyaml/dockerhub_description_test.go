package configyaml_test

import (
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

// Invariante der Docker-Hub-Kurzbeschreibung (DC-FA-DIST-002, ADR-0065).
//
// ZUSAGE: die Datei, aus der das Release die Hub-Description setzt, ist nicht
// leer und passt in Docker Hubs Limit.
//
// KOPPLUNG: Docker Hub misst ZEICHEN, nicht Bytes. Eine Byte-Messung waere bei
// Umlauten strenger als die Regel und meldete einen Text rot, den Docker Hub
// annimmt.
//
// ABGRENZUNG: hier greift die Pruefung beim Schreiben. Im Release ist der
// Darstellungs-Schritt `continue-on-error` — das Bild ist die Zusage, der
// Beschreibungstext ist Praesentation, und das Lastenheft sagt zu, dass sein
// Fehlschlag das Release gruen laesst.
//
// GRENZE: fehlt die Datei, ueberspringt der Test. Er fordert das Artefakt nicht
// ein; er sagt nur, dass es wohlgeformt ist, WENN es da ist.
const dockerHubDescriptionLimit = 100

func dockerHubDescriptionFile() string {
	return repoRoot() + "/packaging/dockerhub/description.txt"
}

func TestDockerHubDescriptionLimit(t *testing.T) {
	path := dockerHubDescriptionFile()
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skipf("%s fehlt — Invariante nicht anwendbar", path)
		}
		t.Fatalf("%s nicht lesbar: %v", path, err)
	}

	desc := strings.TrimRight(string(raw), "\n")
	if strings.TrimSpace(desc) == "" {
		t.Fatalf("%s ist leer — die Hub-Seite zeigte dann keine Kurzbeschreibung", path)
	}
	if strings.Contains(desc, "\n") {
		t.Fatalf("%s hat mehr als eine Zeile — das Release liest nur die erste", path)
	}
	if n := utf8.RuneCountInString(desc); n > dockerHubDescriptionLimit {
		t.Fatalf("%s ist %d Zeichen lang, Docker Hub erlaubt %d", path, n, dockerHubDescriptionLimit)
	}
}

// Der Platzhalter, den das Release durch die Tag-Version ersetzt. Fehlt er,
// zeigte die Hub-Seite eine versionslose Anleitung — das soll beim Schreiben
// auffallen, nicht beim Veroeffentlichen.
func TestDockerHubOverviewCarriesVersionPlaceholder(t *testing.T) {
	path := repoRoot() + "/packaging/dockerhub/overview.md"
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skipf("%s fehlt — Invariante nicht anwendbar", path)
		}
		t.Fatalf("%s nicht lesbar: %v", path, err)
	}
	if !strings.Contains(string(raw), "__VERSION__") {
		t.Fatalf("%s enthaelt kein __VERSION__ — die Hub-Seite zeigte eine versionslose Anleitung", path)
	}
}
