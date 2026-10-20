package cli_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// listConfig ist die Listen-Form aus DC-FA-REQ-001: Überschriften für FA, eine
// Tabelle für QA, beide aus derselben Datei.
const listConfig = `trace:
  requirements:
    - source: spec/lastenheft.md
      id-pattern: 'SFC-FA-\d{2}'
      format: headings
    - source: spec/lastenheft.md
      id-pattern: 'SFC-QA-\d{2}'
      format: table
      table:
        id-column: ID
        text-column: Anforderung
        modality-column: Priorität
`

// mixedLastenheft führt 21 Überschriften- und 6 Tabellen-Anforderungen in einer
// Datei; extra hängt weitere Zeilen an.
func mixedLastenheft(extra string) string {
	var b strings.Builder
	b.WriteString("# Lastenheft\n\n")
	for i := 1; i <= 21; i++ {
		fmt.Fprintf(&b, "#### SFC-FA-%02d — Funktion %d\nDas System MUSS Funktion %d bieten.\n\n", i, i, i)
	}
	b.WriteString("| ID | Aspekt | Anforderung | Priorität |\n|---|---|---|---|\n")
	for i := 1; i <= 6; i++ {
		fmt.Fprintf(&b, "| SFC-QA-%02d | Aspekt | Qualität %d. | MUSS |\n", i, i)
	}
	return b.String() + extra
}

// AK-1: 21 Überschriften-Zeilen, danach 6 Tabellen-Zeilen, Summenzeile über alle.
func TestCLI_TraceListe_VereinigungInListenreihenfolge(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", mixedLastenheft(""))
	write(t, root, ".d-check.yml", listConfig)
	code, stdout, stderr := run(t, "--trace", "--json", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q", code, stderr)
	}
	var doc traceDoc
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		t.Fatalf("JSON: %v\n%s", err, stdout)
	}
	if doc.Total != 27 {
		t.Fatalf("Total = %d, want 27", doc.Total)
	}
	for i, r := range doc.Requirements {
		want := fmt.Sprintf("SFC-FA-%02d", i+1)
		if i >= 21 {
			want = fmt.Sprintf("SFC-QA-%02d", i-20)
		}
		if r.ID != want {
			t.Fatalf("Zeile %d = %s, want %s", i, r.ID, want)
		}
	}
	code, stdout, _ = run(t, "--trace", root)
	if code != 0 || !strings.Contains(stdout, "27 Anforderung(en), 27 Waise(n).") {
		t.Fatalf("Text-Ausgabe: Exit %d, Summenzeile fehlt:\n%s", code, stdout)
	}
}

// AK-2: dieselbe Kennung als Überschrift und als Tabellenzeile ⇒ Exit 2, die
// Meldung nennt Kennung und beide Einträge.
func TestCLI_TraceListe_DoppelteKennungExit2(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", mixedLastenheft("\n#### SFC-QA-03 — Auch als Überschrift\nText.\n"))
	write(t, root, ".d-check.yml", strings.ReplaceAll(strings.ReplaceAll(listConfig, `SFC-FA-\d{2}`, `SFC-(?:FA|QA)-\d{2}`), `SFC-QA-\d{2}`, `SFC-(?:FA|QA)-\d{2}`))
	code, stdout, stderr := run(t, "--trace", root)
	if code != 2 || stdout != "" {
		t.Fatalf("Exit/stdout = %d/%q", code, stdout)
	}
	for _, want := range []string{"SFC-QA-03", "Eintrag 0", "Eintrag 1"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr %q trägt %q nicht", stderr, want)
		}
	}
}

// AK-3 und AK-4: ein Fehler im zweiten Eintrag nennt dessen Index; ein Eintrag,
// der nichts findet, verschwindet nicht in der Vereinigung.
func TestCLI_TraceListe_FehlerNenntDenEintrag(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", mixedLastenheft(""))
	for name, c := range map[string]struct{ old, new, want string }{
		"header fehlt (AK-3)": {"id-column: ID", "id-column: Kennung", "trace.requirements[1]"},
		"nullmenge (AK-4)": {`'SFC-QA-\d{2}'
      format: table
      table:
        id-column: ID
        text-column: Anforderung
        modality-column: Priorität`, `'SFC-RB-\d{2}'`, "trace.requirements[1]"},
	} {
		t.Run(name, func(t *testing.T) {
			write(t, root, ".d-check.yml", strings.Replace(listConfig, c.old, c.new, 1))
			code, stdout, stderr := run(t, "--trace", root)
			if code != 2 || stdout != "" || !strings.Contains(stderr, c.want) {
				t.Fatalf("Exit/stdout/stderr = %d/%q/%q, want 2/\"\"/enthält %q", code, stdout, stderr, c.want)
			}
		})
	}
}

// AK-6: --require-complete wertet die Vereinigung aus — Exit 1, solange eine der
// 27 weder Slice- noch Coverage-Referenz trägt, Exit 0, wenn keine.
func TestCLI_TraceListe_RequireCompleteWertetVereinigungAus(t *testing.T) {
	root := t.TempDir()
	write(t, root, "spec/lastenheft.md", mixedLastenheft(""))
	write(t, root, ".d-check.yml", listConfig)
	var ids []string
	for i := 1; i <= 21; i++ {
		ids = append(ids, fmt.Sprintf("SFC-FA-%02d", i))
	}
	for i := 1; i <= 5; i++ { // SFC-QA-06 bleibt Waise
		ids = append(ids, fmt.Sprintf("SFC-QA-%02d", i))
	}
	write(t, root, "docs/plan/planning/slice-001-x.md", "Liefert "+strings.Join(ids, " ")+".\n")
	code, _, stderr := run(t, "--trace", "--require-complete", root)
	if code != 1 || !strings.Contains(stderr, "1") {
		t.Fatalf("Exit = %d, stderr = %q, want 1", code, stderr)
	}
	write(t, root, "docs/plan/planning/slice-002-y.md", "Liefert SFC-QA-06.\n")
	code, _, stderr = run(t, "--trace", "--require-complete", root)
	if code != 0 {
		t.Fatalf("Exit = %d, stderr = %q, want 0", code, stderr)
	}
}

// AK-7: eine leere Liste ist ein Konfigurationsfehler.
func TestCLI_TraceListe_LeereListeExit2(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".d-check.yml", "trace:\n  requirements: []\n")
	code, stdout, stderr := run(t, "--trace", root)
	if code != 2 || stdout != "" || !strings.Contains(stderr, "nicht leer") {
		t.Fatalf("Exit/stdout/stderr = %d/%q/%q", code, stdout, stderr)
	}
}
