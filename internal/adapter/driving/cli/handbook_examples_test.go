package cli_test

// Handbuch-E2E-Beispiel-Verankerung (slice-062, Dimension B — Schwester zu
// slice-061s Config-Beispiel-Verifikation). Das Benutzerhandbuch dokumentiert
// Kommando-Aufrufe MIT konkreter Ausgabe/Exit-Code (```bash + ```text/```json).
// Diese Verhaltensbehauptungen liefen nie gegen das echte Binary — driftet das
// CLI-Verhalten, bleibt das Handbuch still falsch. Dieser Harness koppelt die
// dokumentierten Ausgabeblöcke an die tatsächliche `cli.Run`-Ausgabe:
//
//   1. FAIL-CLOSED-Klassifikation (TestHandbook_OutputBlocksClassified): jeder
//      ```text-/```json-Ausgabeblock des Handbuchs ist entweder an ein Beispiel
//      verankert ODER trägt den Marker `<!-- d-check-test:not-replayable: … -->`
//      (kein stiller Ausschluss — dieselbe Opt-out-Philosophie wie slice-061).
//   2. REPLAY (TestHandbook_AnchoredExamplesReplay): je verankertem Beispiel
//      eine Fixture, die die Prämisse herstellt; das echte Binary läuft mit den
//      DOKUMENTIERTEN Flags; Exit-Code + Ausgabe-FORM werden geprüft. Form-Anker
//      müssen in BEIDEN stehen (Doku-Block UND echte Ausgabe) — so fängt der
//      Test Drift in beide Richtungen. Auf Form geprüft (Regex/Schlüssel-Menge),
//      nicht auf wörtliche Datei-Zahlen/Pfade (sonst Wartungsfalle).
//
// Beispiel-Auswahl (DoD „E2E-verankert vs. begründet ausgenommen"):
//   VERANKERT (9): §3/§4.1 sauberes Repo (Summary, Exit 0) · §3/§4.1 kaputter
//     Link (Befund-Zeilen-Schema + Exit 1) · §4.9 --doctor (Prosa-Diagnose) ·
//     §4.9 --doctor --json · §4.11 --json · §4.11 --yaml · §4.12 --trace
//     (RTM) · §4.12 --trace --require-complete (tabellarische Quelle,
//     Nullmenge/Exit 0 als dokumentierte v0.42.0-Grenze).
//   AUSGENOMMEN mit Grund:
//     · §4.13 --print-mk: ```text ist eine ABGEKÜRZTE Illustration (Elision
//       „# …"), nicht die wörtliche Ausgabe → not-replayable-Marker.
//     · §4.5–4.8 (enable/disable, ids, matrix, external): prosaische
//       Verhaltensaussagen ohne dokumentierten Ausgabeblock — schon von den
//       Modul-Akzeptanztests gedeckt (TestID001/TestMTX001/TestEXT001/…).
//     · Container-Mount- & `--network none`-Claims (§4.1/§4.2/§4.8) und
//       nativ==Container: schon von `tools/image-test.sh` verankert.
//     · Digest-/`docker pull`-Beispiele (§2): nicht replaybar (Netz/Registry).
//   ```yaml-Blöcke sind Config-Domäne (slice-061s Harness). Die BRÜCKE: ein
//   ```yaml-Block, den slice-061 per `not-config`-Marker als „kein Config-Input"
//   ausnimmt, IST eine CLI-Ausgabe/Illustration — er wird darum auch von DIESER
//   Sweep erfasst (anchored-oder-not-replayable). So rutscht eine künftige
//   yaml-CLI-Ausgabe (z. B. das in §4.9 zugesagte --doctor --yaml) nicht still
//   zwischen beiden Harnessen durch (R1-F-1). Die EINE heutige CLI-Ausgabe-YAML
//   (§4.11 --yaml) ist als E7 verankert.
//
// Richtungs-Semantik der Kopplung: die formTokens (text/json/yaml) sind
// BIDIREKTIONAL (Anker muss in Doku-Block UND echter Ausgabe stehen — Drift auf
// jeder Seite ⇒ rot). Die zusätzliche STRUKTURELLE json/yaml-Kopplung ist
// bewusst nur `dokumentiert ⊆ real`: sie fängt jede Umbenennung/Entfernung eines
// DOKUMENTIERTEN Schlüssels (in beide Richtungen), toleriert aber ein neues,
// undokumentiertes Binary-Feld — die Doku-Blöcke sind abgekürzte Illustrationen,
// ein `actual ⊆ documented` wäre eine False-Positive-Falle (R1-F-2).

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// handbookPath: repo-relativ zum Paket (cli → driving → adapter → internal →
// Wurzel), Präzedenz slice-061. Der Docker-`make test`-Build trägt den vollen
// Repo-Kontext.
const handbookPath = "../../../../docs/user/benutzerhandbuch.md"

// notReplayableMarker nimmt einen dokumentierten Ausgabeblock von der
// E2E-Verankerung aus (externer Zustand, konkrete Digests, Netz, oder eine
// abgekürzte Illustration). Er steht — wie slice-061s not-config — in der Zeile
// UNMITTELBAR vor dem öffnenden Fence, mit Grund. Opt-out, nicht Opt-in: ein
// neuer echter Ausgabeblock ist automatisch prüfpflichtig; ein neuer
// nicht-replaybarer Block bricht die Klassifikation LAUT, bis er annotiert ist.
const notReplayableMarker = "d-check-test:not-replayable"

// notConfigMarker ist slice-061s Marker (`docexamples_test.go`): er nimmt einen
// ```yaml-Block von der Config-Validierung aus, weil er KEIN `.d-check.yml`-Input
// ist — also eine CLI-Ausgabe/Illustration. Genau diese Blöcke gehören in DIESE
// Sweep (Brücke, s. Datei-Kopf / R1-F-1). Bewusste String-Dopplung: die beiden
// Harnesse liegen in verschiedenen Testpaketen; der Marker-Vertrag ist der Doku.
const notConfigMarker = "d-check-test:not-config"

// handbookExample verankert genau ein dokumentiertes Kommando-Beispiel.
type handbookExample struct {
	name       string                          // Label für Meldungen
	outputInfo string                          // Fence-Typ des Ausgabeblocks: text|json|yaml
	outputDisc string                          // Diskriminator: eindeutiger Teilstring GENAU des Ausgabeblocks (skoped auf outputInfo)
	flags      []string                        // Flags, mit denen das echte Binary läuft (der Pfad wird angehängt)
	wantFlags  []string                        // Flag-Token, die in IRGENDEINEM ```bash-Block dokumentiert sein müssen (Aufruf-Kopplung)
	wantExit   int                             // erwarteter Exit-Code
	fixture    func(t *testing.T, root string) // stellt die Prämisse her
	formTokens []string                        // Form-Anker: müssen in BEIDEM stehen — Doku-Block UND echter (stdout+stderr)-Ausgabe
	structured bool                            // zusätzlich: dokumentierte Schlüssel-Menge ⊆ echte (json/yaml)
}

func handbookExamples() []handbookExample {
	idsFixture := func(t *testing.T, root string) {
		write(t, root, ".d-check.yml",
			"ids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/plan/adr/\n")
		write(t, root, "docs/plan/adr/0007-x.md", "# ADR\n")
		write(t, root, "docs/a.md", "nacktes ADR-0007 im Fließtext\n")
	}
	brokenLink := func(t *testing.T, root string) {
		write(t, root, "docs/a.md", "[kaputt](fehlt.md)\n")
	}
	tableOnlyRequirements := func(t *testing.T, root string) {
		write(t, root, ".d-check.yml", "trace:\n  requirements:\n    source: spec/lastenheft.md\n    id-pattern: 'F-[0-9]+'\n")
		write(t, root, "spec/lastenheft.md", "| ID | Modalität | Anforderung |\n| --- | --- | --- |\n| F-1 | Muss | Das Repository enthält alles. |\n")
	}
	crossConsistency := func(t *testing.T, root string) {
		write(t, root, ".d-check.yml", "trace:\n  requirements:\n    id-pattern: 'GG-[A-Z]+-\\d{3}'\n"+
			"  cross-consistency:\n    forward:\n      file: docs/traceability.md\n"+
			"      sections: [\"27.1 Anforderung zu Design\"]\n      req-column: Anforderung\n"+
			"      design-column: Design-Artefakte\n      design-pattern: 'GG-AR-[A-Z0-9-]+'\n"+
			"    backward:\n      file: spec/architecture.md\n      edge-column: Bezug\n"+
			"      req-pattern: 'GG-[A-Z]+-\\d{3}'\n    exclude-req: '^GG-SPEC-'\n")
		write(t, root, "spec/lastenheft.md", "### GG-ARCH-006\nDer Scheduler MUSS entkoppelt sein.\n")
		write(t, root, "docs/plan/planning/slice-071-x.md", "# slice\nBezug: GG-ARCH-006\n")
		write(t, root, "docs/traceability.md", "# Traceability\n\n## 27.1 Anforderung zu Design\n\n"+
			"| Anforderung | Design-Artefakte |\n|---|---|\n| GG-ARCH-006 | GG-AR-COMP-CORE |\n")
		write(t, root, "spec/architecture.md", "# Architektur\n\n## 4 Komponenten\n\n"+
			"| Komponente | Bezug |\n|---|---|\n| GG-AR-COMP-SCHED | GG-ARCH-006 |\n"+
			"| GG-AR-COMP-MID | GG-SPEC-042 |\n")
	}
	return []handbookExample{
		{
			name:       "§3/§4.1 sauberes Repo (Summary, Exit 0)",
			outputInfo: "text",
			outputDisc: "geprüft, 0 Befund(e)",
			flags:      nil,
			wantExit:   0,
			fixture: func(t *testing.T, root string) {
				write(t, root, "docs/a.md", "[ok](b.md)\n")
				write(t, root, "docs/b.md", "x\n")
			},
			formTokens: []string{"Datei(en) geprüft,", "Befund(e)"},
		},
		{
			name:       "§3/§4.1 kaputter Link (Befund-Zeile, Exit 1)",
			outputInfo: "text",
			outputDisc: "target-missing",
			flags:      nil,
			wantExit:   1,
			fixture:    brokenLink,
			// Der Tab-separierte Befund-Zeilen-Schema-Anker + der Grund-Code;
			// Pfad/Zeile bewusst NICHT gekoppelt (Form statt Wortlaut).
			formTokens: []string{"\tfehlt.md\ttarget-missing", "Datei(en) geprüft,", "Befund(e)"},
		},
		{
			name:       "§4.9 --doctor (Prosa-Diagnose)",
			outputInfo: "text",
			outputDisc: "Fix-Kandidat:",
			flags:      []string{"--enable", "ids", "--doctor"},
			wantFlags:  []string{"--enable ids", "--doctor"},
			wantExit:   1,
			fixture:    idsFixture,
			formTokens: []string{"d-check Diagnose —", "ohne Markdown-Link", "Fix-Kandidat:"},
		},
		{
			name:       "§4.12 --trace (RTM-Markdown)",
			outputInfo: "text",
			outputDisc: "# Requirements Traceability Matrix",
			flags:      []string{"--trace"},
			wantFlags:  []string{"--trace"},
			wantExit:   0,
			fixture:    traceRepo,
			formTokens: []string{"# Requirements Traceability Matrix", "| Anforderung", "Anforderung(en),", "Waise(n)."},
		},
		{
			name:       "§4.12 --trace --require-complete (explizite Quelle bleibt leer)",
			outputInfo: "text",
			outputDisc: "ergab 0 Anforderungen",
			flags:      []string{"--trace", "--require-complete"},
			wantFlags:  []string{"--trace", "--require-complete"},
			wantExit:   2,
			fixture:    tableOnlyRequirements,
			formTokens: []string{"d-check: error:", "ergab 0 Anforderungen"},
		},
		{
			// slice-071: der Abgleich erscheint als EIGENER Abschnitt unter der RTM
			// (keine RTM-Spalte) und gatet über das globale --require-complete. Die
			// Form-Anker koppeln beide Richtungslabels + die stderr-Zählzeile.
			name:       "§4.12 --trace --require-complete (Kreuzverweis-Differenzen)",
			outputInfo: "text",
			outputDisc: "## Kreuzverweis-Konsistenz",
			flags:      []string{"--trace", "--require-complete"},
			wantFlags:  []string{"--trace", "--require-complete"},
			wantExit:   1,
			fixture:    crossConsistency,
			formTokens: []string{
				"## Kreuzverweis-Konsistenz",
				"in RTM, ohne Rück-Kante",
				"Rück-Kante, ohne RTM-Eintrag",
				"Differenz(en).",
				"Kreuzverweis-Differenz(en) zwischen Vorwärts- und Rück-Sicht",
			},
		},
		{
			name:       "§4.11 --json (Befundliste)",
			outputInfo: "json",
			outputDisc: "target-missing",
			flags:      []string{"--json"},
			wantFlags:  []string{"--json"},
			wantExit:   1,
			fixture:    brokenLink,
			structured: true,
		},
		{
			name:       "§4.9 --doctor --json (maschinenlesbare Diagnose)",
			outputInfo: "json",
			outputDisc: "reasonText",
			flags:      []string{"--enable", "ids", "--doctor", "--json"},
			wantFlags:  []string{"--enable ids", "--doctor", "--json"},
			wantExit:   1,
			fixture:    idsFixture,
			structured: true,
		},
		{
			name:       "§4.11 --yaml (Befundliste als YAML)",
			outputInfo: "yaml",
			outputDisc: "exitCode:",
			flags:      []string{"--yaml"},
			wantFlags:  []string{"--yaml"},
			wantExit:   1,
			fixture:    brokenLink,
			structured: true,
		},
	}
}

// TestHandbook_TraceParsergrenzenDokumentiert koppelt die nicht aus der
// Konfiguration ableitbaren Parsergrenzen an die Nutzer-Doku. Die eigentliche
// Verhaltensprobe (explizite Quelle + falsches Format + --require-complete =>
// Exit 2) ist als achtes Replay-Beispiel oben verankert; diese Marker verhindern,
// dass Warnung, Definitionsgrammatik oder Waisen-/Owner-Semantik unabhängig
// davon still aus dem Handbuch verschwinden (slice-069).
func TestHandbook_TraceParsergrenzenDokumentiert(t *testing.T) {
	raw, err := os.ReadFile(handbookPath)
	if err != nil {
		t.Fatalf("Handbuch nicht lesbar: %v", err)
	}
	text := strings.Join(strings.Fields(string(raw)), " ")
	for _, want := range []string{
		"ATX-Markdown-Überschrift",
		"erste vollständige Token",
		"Tabellenzeilen, Listen, Fließtext und Setext-Überschriften",
		"Eine ADR-Referenz allein verhindert den Waisenstatus nicht",
		"konditional `coverage` und `modality`",
		"Ohne aktive `modality` gatet jede Waise",
		"gaten nur Waisen der in `require-levels` gelisteten Stufen",
		"Basisdateinamen",
		"Capture-Gruppe 1 zusammen mit",
		"`table.id-column`, genau eine von `table.text-column` oder `table.text-columns`",
		"`duplicate-ids: error`",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("Trace-Parsergrenze %q fehlt im Benutzerhandbuch", want)
		}
	}
	for _, forbidden := range []string{
		"**niemand** referenziert (Waise)",
		"mindestens eine Waise ⇒ Exit 1 statt 0",
	} {
		if strings.Contains(text, forbidden) {
			t.Errorf("widersprüchliche Trace-Aussage %q steht noch im Benutzerhandbuch", forbidden)
		}
	}
}

// mustBlocks liest das Handbuch und extrahiert alle Fenced-Blöcke — fail-closed
// bei Lesefehler, unbalanciertem Fence oder leerer Menge (Extraktor kaputt).
func mustBlocks(t *testing.T) []fencedBlock {
	t.Helper()
	raw, err := os.ReadFile(handbookPath)
	if err != nil {
		t.Fatalf("Handbuch nicht lesbar (fail-closed): %s: %v", handbookPath, err)
	}
	blocks, err := extractFencedBlocks(string(raw))
	if err != nil {
		t.Fatalf("%s: %v (fail-closed)", filepath.Base(handbookPath), err)
	}
	if len(blocks) == 0 {
		t.Fatal("kein einziger Fenced-Block extrahiert (fail-closed — Extraktor kaputt?)")
	}
	return blocks
}

// TestHandbook_OutputBlocksClassified verriegelt: JEDER dokumentierte
// ```text-/```json-Ausgabeblock ist entweder an genau ein Beispiel verankert
// (Diskriminator im Block, skoped auf den Fence-Typ) oder trägt den
// not-replayable-Marker. Ein neuer, unklassifizierter Ausgabeblock ⇒ rot
// (Silent-Exclusion-Guard, slice-061-Philosophie).
func TestHandbook_OutputBlocksClassified(t *testing.T) {
	blocks := mustBlocks(t)
	examples := handbookExamples()
	outputBlocks := 0
	for _, b := range blocks {
		info := firstToken(b.info)
		// Sweep-Menge: alle ```text/```json plus die BRÜCKEN-yaml (ein
		// ```yaml, das slice-061 per not-config-Marker als Nicht-Config
		// ausnimmt = CLI-Ausgabe/Illustration). Config-```yaml (ohne Marker)
		// bleibt slice-061s Domäne. R1-F-1: sonst rutschte eine künftige
		// yaml-CLI-Ausgabe still zwischen beiden Harnessen durch.
		swept := info == "text" || info == "json" ||
			(info == "yaml" && strings.Contains(b.prev, notConfigMarker))
		if !swept {
			continue
		}
		outputBlocks++
		if strings.Contains(b.prev, notReplayableMarker) {
			continue // begründet ausgenommen
		}
		matches := 0
		for _, e := range examples {
			if e.outputInfo == info && strings.Contains(b.body, e.outputDisc) {
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("%s:%d: ```%s-Ausgabeblock nicht eindeutig klassifiziert (%d passende Beispiele) — "+
				"entweder an ein Beispiel verankern oder Marker `<!-- %s: <Grund> -->` in die Zeile UNMITTELBAR vor den Fence setzen",
				filepath.Base(handbookPath), b.line, info, matches, notReplayableMarker)
		}
	}
	if outputBlocks == 0 {
		t.Fatal("keine text/json/not-config-yaml-Ausgabeblöcke gefunden (fail-closed — Extraktor kaputt?)")
	}
}

// TestHandbook_AnchoredExamplesReplay fährt je verankertem Beispiel das echte
// Binary gegen eine Prämissen-Fixture und koppelt Exit-Code + Ausgabe-Form an
// den dokumentierten Block. Adversariale Probe: verletzt man einen Form-Anker
// im Doku-Block, fehlt er dort ⇒ rot; driftet das Binary, fehlt er in der
// echten Ausgabe ⇒ rot.
func TestHandbook_AnchoredExamplesReplay(t *testing.T) {
	blocks := mustBlocks(t)
	var bashBodies []string
	for _, b := range blocks {
		if firstToken(b.info) == "bash" {
			bashBodies = append(bashBodies, b.body)
		}
	}

	for _, e := range handbookExamples() {
		// (a) den dokumentierten Ausgabeblock lokalisieren — genau einer
		// (Uniqueness-Guard: Diskriminator-Tippfehler / entfernter Block ⇒ rot).
		var docBlocks []fencedBlock
		for _, b := range blocks {
			if firstToken(b.info) == e.outputInfo && strings.Contains(b.body, e.outputDisc) {
				docBlocks = append(docBlocks, b)
			}
		}
		if len(docBlocks) != 1 {
			t.Errorf("%s: erwartet genau 1 Doku-```%s-Block mit Disc %q, gefunden %d (fail-closed)",
				e.name, e.outputInfo, e.outputDisc, len(docBlocks))
			continue
		}
		docBlock := docBlocks[0]

		// (b) Aufruf-Kopplung: die Flags, mit denen wir replayen, sind auch
		// im Handbuch dokumentiert.
		for _, ft := range e.wantFlags {
			if !anyContains(bashBodies, ft) {
				t.Errorf("%s: Flag %q in keinem ```bash-Block dokumentiert (Aufruf driftet)", e.name, ft)
			}
		}

		// (c) echtes Binary gegen die Prämissen-Fixture.
		root := t.TempDir()
		e.fixture(t, root)
		args := append(append([]string{}, e.flags...), root)
		code, stdout, stderr := run(t, args...)
		if code != e.wantExit {
			t.Errorf("%s: Exit = %d, want %d (stderr: %s)", e.name, code, e.wantExit, stderr)
			continue
		}
		combined := stdout + "\n" + stderr

		// (d) Form-Anker in BEIDEM.
		for _, tok := range e.formTokens {
			if !strings.Contains(docBlock.body, tok) {
				t.Errorf("%s: Form-Anker %q fehlt im Doku-Ausgabeblock (Zeile %d) — Doku driftet oder Disc falsch",
					e.name, tok, docBlock.line)
			}
			if !strings.Contains(combined, tok) {
				t.Errorf("%s: Form-Anker %q fehlt in der echten Ausgabe — Binary driftet von der Doku\n%s",
					e.name, tok, combined)
			}
		}

		// (e) strukturelle Kopplung (json/yaml): jeder dokumentierte Schlüssel
		// erscheint in der echten Ausgabe (fängt Umbenennungen in beide Rtg.).
		if e.structured {
			docKeys := keysOf(t, e.outputInfo, docBlock.body)
			actKeys := keysOf(t, e.outputInfo, stdout)
			if len(docKeys) == 0 {
				t.Errorf("%s: Doku-Block trägt keine Schlüssel (Parser/Disc kaputt?)", e.name)
			}
			for k := range docKeys {
				if !actKeys[k] {
					t.Errorf("%s: dokumentierter Schlüssel %q fehlt in der echten %s-Ausgabe", e.name, k, e.outputInfo)
				}
			}
		}
	}
}

func anyContains(haystacks []string, needle string) bool {
	for _, h := range haystacks {
		if strings.Contains(h, needle) {
			return true
		}
	}
	return false
}

type fencedBlock struct {
	info string // Info-String des öffnenden Fence (z. B. "text", "json", "yaml title=x")
	line int    // 1-basierte Zeile des öffnenden Fence
	body string // Inhalt zwischen den Fences
	prev string // Quellzeile UNMITTELBAR vor dem öffnenden Fence (für Marker), "" bei Zeile 1
}

// extractFencedBlocks liest alle Fenced-Blöcke aus Markdown mit korrekter
// Fence-Zustandsverfolgung: ein Fence öffnet auf ``` (mit beliebigem
// Info-String) und schließt erst auf ein nacktes ``` — ein ```-Zeile INNERHALB
// eines anderen Fences (z. B. ```json in einem ```markdown-Beispiel) ist damit
// Body, kein eigener Block. Fail-closed bei unbalanciertem Fence. Reine
// Funktion (kein *testing.T) — Extraktion und Fail-closed-Pfad sind
// unit-testbar. Analog zu slice-061s extractYAMLBlocks, hier fence-typ-agnostisch.
func extractFencedBlocks(content string) ([]fencedBlock, error) {
	lines := strings.Split(content, "\n")
	var blocks []fencedBlock
	inFence := false
	var info, prev string
	openLine := 0
	var body []string
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !inFence {
			if strings.HasPrefix(trimmed, "```") {
				info = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				inFence = true
				openLine = i + 1
				body = nil
				prev = ""
				if i > 0 {
					prev = lines[i-1]
				}
			}
			continue
		}
		if trimmed == "```" { // nacktes ``` schließt den offenen Fence
			blocks = append(blocks, fencedBlock{info: info, line: openLine, body: strings.Join(body, "\n"), prev: prev})
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

// firstToken gibt das erste Wort des Info-Strings zurück (Kleinbuchstaben) —
// „yaml title=x" → „yaml".
func firstToken(info string) string {
	f := strings.Fields(info)
	if len(f) == 0 {
		return ""
	}
	return strings.ToLower(f[0])
}

// keysOf parst einen json-/yaml-Block und sammelt seine Schlüssel-Pfade
// (`findings[].file`, `summary.filesChecked`, `exitCode`, …).
func keysOf(t *testing.T, info, body string) map[string]bool {
	t.Helper()
	var v any
	var err error
	switch info {
	case "json":
		err = json.Unmarshal([]byte(body), &v)
	case "yaml":
		err = yaml.Unmarshal([]byte(body), &v)
	default:
		t.Fatalf("keysOf: unerwarteter Info-Typ %q", info)
	}
	if err != nil {
		t.Fatalf("keysOf(%s): %v\n%s", info, err, body)
	}
	set := map[string]bool{}
	collectKeys(v, "", set)
	return set
}

func collectKeys(v any, prefix string, set map[string]bool) {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			set[p] = true
			collectKeys(vv, p, set)
		}
	case []any:
		for _, e := range t {
			collectKeys(e, prefix+"[]", set)
		}
	}
}

// TestExtractFencedBlocks verriegelt die Extraktor-Semantik (Slice-Plan §3):
// einfache Extraktion, Fence-Zustand (json-in-markdown ist kein Block),
// prev-Zeile für den Marker, Fail-closed bei unbalanciertem Fence,
// Info-String-Varianten.
func TestExtractFencedBlocks(t *testing.T) {
	// (1) einfacher Block mit korrekter Zeile + prev
	b, err := extractFencedBlocks("vor\n```text\nZeile\n```\nnach\n")
	if err != nil || len(b) != 1 || strings.TrimSpace(b[0].body) != "Zeile" || b[0].line != 2 || b[0].prev != "vor" {
		t.Fatalf("(1) blocks=%+v err=%v", b, err)
	}
	// (2) ```json INNERHALB eines ```markdown-Blocks ist KEIN eigener Block
	b, err = extractFencedBlocks("```markdown\n# x\n```json\n{\"a\":1}\n```\n")
	if err != nil || len(b) != 1 || firstToken(b[0].info) != "markdown" {
		t.Fatalf("(2) json im markdown-Block darf nicht extrahiert werden: blocks=%+v err=%v", b, err)
	}
	// (3) prev ist die Zeile unmittelbar vor dem Fence (Marker-Regel)
	b, _ = extractFencedBlocks("<!-- d-check-test:not-replayable: x -->\n```text\nA\n```\n")
	if len(b) != 1 || !strings.Contains(b[0].prev, notReplayableMarker) {
		t.Fatalf("(3) prev trägt den Marker nicht: %+v", b)
	}
	// (3b) eine Leerzeile zwischen Marker und Fence ⇒ prev ist leer (Marker wirkungslos)
	b, _ = extractFencedBlocks("<!-- d-check-test:not-replayable: x -->\n\n```text\nA\n```\n")
	if len(b) != 1 || strings.Contains(b[0].prev, notReplayableMarker) {
		t.Fatalf("(3b) Leerzeile zwischen Marker und Fence ⇒ prev ohne Marker: %+v", b)
	}
	// (4) unbalancierter Fence ⇒ Fehler (fail-closed)
	if _, err := extractFencedBlocks("```text\nA\n"); err == nil {
		t.Fatal("(4) unbalancierter Fence: Fehler erwartet")
	}
	// (5) Info-String-Varianten: erstes Token zählt
	got, err := extractFencedBlocks("```yaml title=x\na: 1\n```\n")
	if err != nil || len(got) != 1 || firstToken(got[0].info) != "yaml" {
		t.Fatalf("(5) Info-Suffix: %+v err=%v", got, err)
	}
}

// TestHandbookVerifiers verriegelt die Kopplungs-Helfer auf synthetischen Daten
// — encodiert die adversariale Probe permanent: fehlt ein dokumentierter
// Schlüssel in der „echten" Ausgabe, meldet die Mengen-Differenz das.
func TestHandbookVerifiers(t *testing.T) {
	// collectKeys sammelt verschachtelte Pfade inkl. Array-Marker.
	set := map[string]bool{}
	var v any
	if err := json.Unmarshal([]byte(`{"findings":[{"file":"x","fix":{"note":"n"}}],"exitCode":1}`), &v); err != nil {
		t.Fatal(err)
	}
	collectKeys(v, "", set)
	for _, want := range []string{"findings[].file", "findings[].fix.note", "exitCode"} {
		if !set[want] {
			t.Fatalf("Schlüssel %q fehlt: %v", want, set)
		}
	}
	// json == yaml Schlüssel-Ernte (Struktur-Parität), damit die Kopplung
	// fence-typ-übergreifend gilt.
	jk := keysOf(t, "json", `{"summary":{"findingCount":1}}`)
	yk := keysOf(t, "yaml", "summary:\n  findingCount: 1\n")
	if !jk["summary.findingCount"] || !yk["summary.findingCount"] {
		t.Fatalf("json/yaml-Schlüssel-Parität verletzt: json=%v yaml=%v", jk, yk)
	}
	// Adversariale Probe: dokumentierter Schlüssel fehlt in der „echten" Menge.
	doc := keysOf(t, "json", `{"reasonText":"x"}`)
	act := keysOf(t, "json", `{"reason":"x"}`)
	if act["reasonText"] {
		t.Fatal("Setup falsch: reasonText darf in act nicht sein")
	}
	missing := false
	for k := range doc {
		if !act[k] {
			missing = true
		}
	}
	if !missing {
		t.Fatal("adversariale Probe: umbenannter Schlüssel müsste als fehlend auffallen")
	}
}
