package configyaml_test

import (
	"fmt"
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
		"../../../../docs/user/operations.md",
		"../../../../README.md",
		"../../../../README.de.md",
	}
}

// notConfigMarker nimmt einen ```yaml-Block, der KEIN `.d-check.yml`-Input ist
// (z. B. ein `--yaml`-Ausgabe-Beispiel `findings:`), von der Validierung aus.
// Er steht in der Zeile UNMITTELBAR vor dem öffnenden Fence, mit Grund (eine
// Leerzeile dazwischen macht ihn wirkungslos). Bewusst Opt-out, nicht Opt-in:
// ein neues echtes Config-Beispiel ist automatisch abgedeckt; ein neuer
// Nicht-Config-Block bricht den Test LAUT, bis er annotiert ist (fail-closed
// zum Prüfen hin — dieselbe Philosophie wie `d-check:ignore`).
const notConfigMarker = "d-check-test:not-config"

// TestDocExamples_ConfigBeispieleValidieren extrahiert jeden ```yaml-Block der
// Live-Doku und validiert ihn gegen configyaml.Decode. Ein Config-Beispiel, das
// der Validator ablehnt und nicht als not-config markiert ist, macht den Test
// rot — genau der `hostpaths.prefixes: ["/home"]`-Fall (Fix c8c33a0), der bisher
// durch kein Gate fiel (d-check scannt Fenced-Code nicht). Schließt die
// Blindspot-Klasse „Doku-Config-Beispiel widerspricht dem Validator".
func TestDocExamples_ConfigBeispieleValidieren(t *testing.T) {
	validated := 0
	for _, f := range docExampleFiles() {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("Doku-Datei nicht lesbar (fail-closed): %s: %v", f, err)
		}
		blocks, err := extractYAMLBlocks(string(raw))
		if err != nil {
			t.Fatalf("%s: %v (fail-closed)", filepath.Base(f), err)
		}
		for _, b := range blocks {
			if b.skip {
				continue
			}
			validated++
			if _, err := configyaml.Decode([]byte(b.body)); err != nil {
				t.Errorf("%s:%d: ```yaml-Config-Beispiel von configyaml.Decode abgelehnt: %v\n"+
					"(Nicht-Config-Block? Marker `<!-- %s: <Grund> -->` in die Zeile UNMITTELBAR vor den Fence setzen)",
					filepath.Base(f), b.line, err, notConfigMarker)
			}
		}
	}
	if validated == 0 {
		t.Fatal("kein einziges ```yaml-Config-Beispiel validiert (fail-closed — Extraktor kaputt oder alles übersprungen?)")
	}
}

type yamlBlock struct {
	line int    // 1-basierte Zeile des öffnenden Fence
	body string // Inhalt zwischen den Fences
	skip bool   // per not-config-Marker ausgenommen
}

// extractYAMLBlocks liest die ```yaml/```yml-Blöcke aus Markdown mit korrekter
// **Fence-Zustandsverfolgung**: ein Fence öffnet auf ``` (mit beliebigem
// Info-String) und schließt erst auf ein nacktes ``` — eine ```yaml-Zeile
// INNERHALB eines anderen Fences (z. B. eines ```markdown-Beispiels) ist damit
// Body, kein eigener Block. Fail-closed bei unbalanciertem Fence (Fehler statt
// stillem Verschlucken). Ein not-config-Marker in der Zeile unmittelbar vor dem
// öffnenden Fence markiert den Block als skip. Reine Funktion (kein *testing.T)
// — dadurch sind Extraktion und Fail-closed-Pfad unit-testbar.
func extractYAMLBlocks(content string) ([]yamlBlock, error) {
	lines := strings.Split(content, "\n")
	var blocks []yamlBlock
	inFence := false
	var info string
	openLine := 0
	var body []string
	markerBefore := false
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !inFence {
			if strings.HasPrefix(trimmed, "```") {
				info = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				inFence = true
				openLine = i + 1
				body = nil
				markerBefore = i > 0 && strings.Contains(lines[i-1], notConfigMarker)
			}
			continue
		}
		if trimmed == "```" { // nacktes ``` schließt den offenen Fence
			if isYAMLInfo(info) {
				blocks = append(blocks, yamlBlock{line: openLine, body: strings.Join(body, "\n"), skip: markerBefore})
			}
			inFence = false
			continue
		}
		body = append(body, lines[i])
	}
	if inFence {
		return nil, fmt.Errorf("unbalancierter Fence ab Zeile %d", openLine)
	}
	return blocks, nil
}

// isYAMLInfo erkennt einen YAML-Fence am ersten Token des Info-Strings —
// `yaml`, `yml` (case-insensitiv), auch mit Zusatz-Attributen (`yaml title=…`).
func isYAMLInfo(info string) bool {
	f := strings.Fields(info)
	if len(f) == 0 {
		return false
	}
	switch strings.ToLower(f[0]) {
	case "yaml", "yml":
		return true
	}
	return false
}

// TestExtractYAMLBlocks verriegelt die Extraktor-Semantik (Slice-Plan §4):
// einfache Extraktion, Fence-Zustand (yaml-in-markdown ist kein Block),
// Marker-Skip inkl. der „unmittelbare Vorzeile"-Regel, Fail-closed bei
// unbalanciertem Fence und die Info-String-Varianten.
func TestExtractYAMLBlocks(t *testing.T) {
	// (1) einfacher yaml-Block wird mit korrekter Zeile extrahiert
	b, err := extractYAMLBlocks("vor\n```yaml\na: 1\n```\nnach\n")
	if err != nil || len(b) != 1 || strings.TrimSpace(b[0].body) != "a: 1" || b[0].line != 2 {
		t.Fatalf("(1) blocks=%+v err=%v", b, err)
	}
	// (2) ```yaml INNERHALB eines ```markdown-Blocks ist KEIN eigener Block
	b, err = extractYAMLBlocks("```markdown\n# x\n```yaml\nnicht: extrahiert\n```\n")
	if err != nil || len(b) != 0 {
		t.Fatalf("(2) yaml im markdown-Block darf nicht extrahiert werden: blocks=%+v err=%v", b, err)
	}
	// (3a) Marker in der Vorzeile ⇒ skip
	b, _ = extractYAMLBlocks("<!-- d-check-test:not-config: x -->\n```yaml\na: 1\n```\n")
	if len(b) != 1 || !b[0].skip {
		t.Fatalf("(3a) Marker-Vorzeile: skip erwartet: %+v", b)
	}
	// (3b) Leerzeile zwischen Marker und Fence ⇒ KEIN skip
	b, _ = extractYAMLBlocks("<!-- d-check-test:not-config: x -->\n\n```yaml\na: 1\n```\n")
	if len(b) != 1 || b[0].skip {
		t.Fatalf("(3b) Leerzeile zwischen Marker und Fence ⇒ kein skip: %+v", b)
	}
	// (4) unbalancierter Fence ⇒ Fehler (fail-closed)
	if _, err := extractYAMLBlocks("```yaml\na: 1\n"); err == nil {
		t.Fatal("(4) unbalancierter Fence: Fehler erwartet")
	}
	// (5) yml/YAML/Info-Suffix werden als YAML erkannt
	for _, open := range []string{"```yml", "```YAML", "```yaml title=x"} {
		got, err := extractYAMLBlocks(open + "\na: 1\n```\n")
		if err != nil || len(got) != 1 {
			t.Fatalf("(5) %q: 1 Block erwartet: %+v err=%v", open, got, err)
		}
	}
}
