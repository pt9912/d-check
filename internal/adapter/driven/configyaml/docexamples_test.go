package configyaml_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

// docExampleFiles: die wartungsaktive Nutzer-Doku, deren gefencete
// `.d-check.yml`-Beispiele Leser kopieren (repo-relativ zum Paket:
// configyaml → driven → adapter → internal → Wurzel). Kein `done/`-Slice, kein
// Review-Report, kein Baseline-Cache — nur Doku, deren Config-Beispiele echt
// verbindlich sind. Der Docker-`make test`-Build trägt den vollen Repo-Kontext.
func docExampleFiles() []string {
	return []string{
		"../../../../docs/user/benutzerhandbuch.md",
		"../../../../README.md",
		"../../../../README.de.md",
	}
}

// notConfigMarker nimmt einen ```yaml-Block, der KEIN `.d-check.yml`-Input ist
// (z. B. ein `--yaml`-Ausgabe-Beispiel `findings:`), von der Validierung aus.
// Er steht in der Zeile unmittelbar vor dem öffnenden Fence, mit Grund. Bewusst
// Opt-out, nicht Opt-in: ein neues echtes Config-Beispiel ist automatisch
// abgedeckt; ein neuer Nicht-Config-Block bricht den Test LAUT, bis er annotiert
// ist (fail-closed zum Prüfen hin — dieselbe Philosophie wie `d-check:ignore`).
const notConfigMarker = "d-check-test:not-config"

// TestDocExamples_ConfigBeispieleValidieren extrahiert jeden ```yaml-Block der
// Live-Doku und validiert ihn gegen configyaml.Decode. Ein Config-Beispiel, das
// der Validator ablehnt und nicht als not-config markiert ist, macht den Test
// rot — genau der `hostpaths.prefixes: ["/home"]`-Fall (Fix c8c33a0), der bisher
// durch kein Gate fiel (d-check scannt Fenced-Code nicht). Schließt die
// Blindspot-Klasse „Doku-Config-Beispiel widerspricht dem Validator".
func TestDocExamples_ConfigBeispieleValidieren(t *testing.T) {
	total := 0
	for _, f := range docExampleFiles() {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("Doku-Datei nicht lesbar (fail-closed): %s: %v", f, err)
		}
		for _, b := range extractYAMLBlocks(t, f, string(raw)) {
			total++
			if b.skip {
				continue
			}
			if _, err := configyaml.Decode([]byte(b.body)); err != nil {
				t.Errorf("%s:%d: ```yaml-Config-Beispiel von configyaml.Decode abgelehnt: %v\n"+
					"(Nicht-Config-Block? Marker `<!-- %s: <Grund> -->` vor den Fence setzen)",
					filepath.Base(f), b.line, err, notConfigMarker)
			}
		}
	}
	if total == 0 {
		t.Fatal("keine ```yaml-Blöcke in der Doku-Menge gefunden (fail-closed — Extraktor kaputt?)")
	}
}

type yamlBlock struct {
	line int    // 1-basierte Zeile des öffnenden Fence
	body string // Inhalt zwischen den Fences
	skip bool   // per not-config-Marker ausgenommen
}

// extractYAMLBlocks liest ```yaml … ``` Blöcke zeilenweise. Fail-closed bei
// unbalanciertem öffnendem Fence (kein stilles Verschlucken). Ein
// not-config-Marker in der Zeile unmittelbar vor dem öffnenden Fence markiert
// den Block als skip.
func extractYAMLBlocks(t *testing.T, file, content string) []yamlBlock {
	t.Helper()
	lines := strings.Split(content, "\n")
	var blocks []yamlBlock
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "```yaml" && !strings.HasPrefix(trimmed, "```yaml ") {
			continue
		}
		openLine := i + 1
		skip := i > 0 && strings.Contains(lines[i-1], notConfigMarker)
		var body []string
		j := i + 1
		closed := false
		for ; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "```" {
				closed = true
				break
			}
			body = append(body, lines[j])
		}
		if !closed {
			t.Fatalf("%s:%d: unbalancierter ```yaml-Fence (fail-closed)", filepath.Base(file), openLine)
		}
		blocks = append(blocks, yamlBlock{line: openLine, body: strings.Join(body, "\n"), skip: skip})
		i = j
	}
	return blocks
}
