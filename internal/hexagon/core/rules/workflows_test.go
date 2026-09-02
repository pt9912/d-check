package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/workflowyaml"
	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// wfCfg ist die Minimal-Konfiguration: das Verzeichnis ist der
// Aktivierungs-Schalter.
func wfCfg() model.WorkflowsConfig { return model.WorkflowsConfig{Dir: ".github/workflows"} }

// wfRun laesst das Modul ueber eine Datei-Menge laufen.
func wfRun(files map[string]string) []model.Finding {
	return CheckWorkflows(coretest.NewMemFS(files), workflowyaml.New(), wfCfg())
}

// reasons sammelt die Grund-Codes eines Befundsatzes.
func reasons(fs []model.Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Reason)
	}
	return out
}

func hasReason(fs []model.Finding, want string) bool {
	for _, f := range fs {
		if f.Reason == want {
			return true
		}
	}
	return false
}

// Happy Path: gepinnte fremde Referenz plus lokale Referenz, deren Ziel
// existiert und deren Job die geforderten Rechte deklariert.
func TestWorkflowsHappyPath(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "permissions: {}\njobs:\n" +
			"  build:\n    permissions:\n      contents: read\n    steps:\n" +
			"      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1\n" +
			"  call:\n    permissions:\n      contents: read\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "permissions: {}\njobs:\n  s:\n    permissions:\n      contents: read\n",
	}
	if f := wfRun(files); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", reasons(f))
	}
}

// Negative (Pin): beweglicher Tag und SHA ohne Tag-Kommentar.
func TestWorkflowsPinForm(t *testing.T) {
	base := ".github/workflows/b.yml"
	target := "permissions: {}\njobs:\n  s:\n    runs-on: x\n"
	beweglich := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n      - uses: actions/checkout@v4\n",
		base:                      target,
	}
	f := wfRun(beweglich)
	if !hasReason(f, ReasonUsesPinMissing) {
		t.Fatalf("erwartet uses-pin-missing, got %+v", reasons(f))
	}
	if f[0].Line != 4 {
		t.Errorf("Line = %d, want 4 (die uses:-Zeile)", f[0].Line)
	}
	ohneTag := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n" +
			"      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1\n",
		base: target,
	}
	if f := wfRun(ohneTag); !hasReason(f, ReasonUsesPinUntagged) {
		t.Fatalf("erwartet uses-pin-untagged, got %+v", reasons(f))
	}
}

// Negative: die lokale Referenz zeigt ins Leere.
func TestWorkflowsLocalMissing(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    uses: ./.github/workflows/weg.yml\n",
	}
	f := wfRun(files)
	if !hasReason(f, ReasonUsesLocalMissing) {
		t.Fatalf("erwartet uses-local-missing, got %+v", reasons(f))
	}
}

// DER BELEGTE AUSFALL: der Job erbt permissions vom Kopf und deklariert selbst
// keines, das Ziel verlangt Rechte.
func TestWorkflowsPermsUndeclared(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "permissions: {}\njobs:\n  c:\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "permissions: {}\njobs:\n  s:\n    permissions:\n      contents: read\n",
	}
	f := wfRun(files)
	if !hasReason(f, ReasonUsesLocalPermsUndecl) {
		t.Fatalf("erwartet uses-local-perms-undeclared, got %+v", reasons(f))
	}
	if f[0].Line != 4 {
		t.Errorf("Line = %d, want 4 (die uses:-Zeile)", f[0].Line)
	}
}

// Negative: der Aufrufer fuehrt den Scope zu niedrig bzw. nennt ihn nicht.
func TestWorkflowsPermsNarrow(t *testing.T) {
	zuNiedrig := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    permissions:\n      contents: read\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n    permissions:\n      contents: write\n",
	}
	if f := wfRun(zuNiedrig); !hasReason(f, ReasonUsesLocalPermsNarrow) {
		t.Fatalf("read gegen write: erwartet uses-local-perms-narrow, got %+v", reasons(f))
	}
	// Ein vom Aufrufer NICHT genannter Scope ist `none` — nicht "egal".
	fehlt := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    permissions:\n      contents: read\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n    permissions:\n      packages: write\n",
	}
	if f := wfRun(fehlt); !hasReason(f, ReasonUsesLocalPermsNarrow) {
		t.Fatalf("fehlender Scope: erwartet uses-local-perms-narrow, got %+v", reasons(f))
	}
}

// Boundary: `read-all` deckt ein read-Verlangen und deckt ein write nicht.
func TestWorkflowsReadAll(t *testing.T) {
	deckt := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    permissions: read-all\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n    permissions:\n      contents: read\n",
	}
	if f := wfRun(deckt); f != nil {
		t.Fatalf("read-all gegen contents:read — erwartet befundfrei, got %+v", reasons(f))
	}
	deckNicht := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    permissions: read-all\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n    permissions:\n      contents: write\n",
	}
	if f := wfRun(deckNicht); !hasReason(f, ReasonUsesLocalPermsNarrow) {
		t.Fatalf("read-all gegen contents:write — erwartet narrow, got %+v", reasons(f))
	}
}

// Boundary: verlangt das Ziel nichts, gibt es nichts zu vergleichen — auch
// ohne Deklaration beim Aufrufer.
func TestWorkflowsZielOhneForderung(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "permissions: {}\njobs:\n  c:\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n    runs-on: x\n",
	}
	if f := wfRun(files); f != nil {
		t.Fatalf("erwartet befundfrei, got %+v", reasons(f))
	}
}

// Boundary: eine Referenz in ein FREMDES Repository wird auf den Pin geprueft,
// nicht auf Rechte — ihr Inhalt liegt nicht vor.
func TestWorkflowsFremdesRepoNurPin(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n" +
			"    uses: other/repo/.github/workflows/x.yml@3d3c42e5aac5ba805825da76410c181273ba90b1 # v1\n",
	}
	if f := wfRun(files); f != nil {
		t.Fatalf("erwartet befundfrei (Pin ok, keine Rechte-Pruefung), got %+v", reasons(f))
	}
}

// fail-closed: leere Pruefmenge ist kein Gruen.
func TestWorkflowsLeereMenge(t *testing.T) {
	leer := map[string]string{"docs/a.md": "# x\n"}
	f := CheckWorkflows(coretest.NewMemFS(leer), workflowyaml.New(), wfCfg())
	if !hasReason(f, ReasonUsesPinMissing) {
		t.Fatalf("kein Kandidat: erwartet fail-closed-Befund, got %+v", reasons(f))
	}
	ohneRef := map[string]string{".github/workflows/a.yml": "jobs:\n  b:\n    runs-on: x\n"}
	if f := wfRun(ohneRef); !hasReason(f, ReasonUsesPinMissing) {
		t.Fatalf("keine Referenz: erwartet fail-closed-Befund, got %+v", reasons(f))
	}
}

// fail-closed: unlesbares YAML meldet, statt uebersprungen zu werden.
func TestWorkflowsUnparsable(t *testing.T) {
	kaputt := map[string]string{".github/workflows/a.yml": "jobs:\n  b:\n   - [unbalanced\n"}
	if f := wfRun(kaputt); !hasReason(f, ReasonWorkflowUnparsable) {
		t.Fatalf("erwartet workflow-unparsable, got %+v", reasons(f))
	}
	// Auch das ZIEL einer lokalen Referenz faellt unter dieselbe Zusage.
	zielKaputt := map[string]string{
		".github/workflows/a.yml": "jobs:\n  c:\n    permissions:\n      contents: read\n    uses: ./.github/workflows/b.yml\n",
		".github/workflows/b.yml": "jobs:\n  s:\n   - [unbalanced\n",
	}
	if f := wfRun(zielKaputt); !hasReason(f, ReasonWorkflowUnparsable) {
		t.Fatalf("Ziel unlesbar: erwartet workflow-unparsable, got %+v", reasons(f))
	}
}

// Modul-aus: ohne dir wird KEINE Datei geoeffnet und der Befundsatz ist
// byte-identisch (DC-QA-02).
func TestWorkflowsInert(t *testing.T) {
	files := map[string]string{".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n      - uses: actions/checkout@v4\n"}
	if f := CheckWorkflows(coretest.NewMemFS(files), workflowyaml.New(), model.WorkflowsConfig{}); f != nil {
		t.Fatalf("ohne dir erwartet inert, got %+v", reasons(f))
	}
}

// exempt-paths nimmt eine Datei heraus — hebt den Leerlauf-Befund aber nicht
// aus: bleiben null Kandidaten, ist das derselbe fail-closed-Befund.
func TestWorkflowsExemptPaths(t *testing.T) {
	files := map[string]string{".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n      - uses: actions/checkout@v4\n"}
	cfg := model.WorkflowsConfig{Dir: ".github/workflows", ExemptPaths: []string{".github/workflows/a.yml"}}
	f := CheckWorkflows(coretest.NewMemFS(files), workflowyaml.New(), cfg)
	if len(f) != 1 || f[0].Reason != ReasonUsesPinMissing || f[0].File != ".github/workflows" {
		t.Fatalf("erwartet genau den Leerlauf-Befund auf dem Verzeichnis, got %+v", f)
	}
}

// countReason zaehlt die Treffer eines Grund-Codes.
func countReason(fs []model.Finding, want string) int {
	n := 0
	for _, f := range fs {
		if f.Reason == want {
			n++
		}
	}
	return n
}

const sha1 = "3d3c42e5aac5ba805825da76410c181273ba90b1"

// Negativ/Kern: derselbe SHA traegt in DERSELBEN Datei zwei verschiedene
// Tag-Kommentare ueber drei Zeilen -- eine Zeile je Fund, nicht einer je SHA
// (DoD: "SHA mit 2 Kommentaren ueber 3 Zeilen ⇒ 3 Befunde").
func TestWorkflowsPinTagConflictSameFile(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n" +
			"      - uses: actions/checkout@" + sha1 + " # v4.2.0\n" +
			"      - uses: actions/checkout@" + sha1 + " # v3.6.0\n" +
			"      - uses: actions/checkout@" + sha1 + " # v4.2.0\n",
	}
	f := wfRun(files)
	if got := countReason(f, ReasonUsesPinTagConflict); got != 3 {
		t.Fatalf("erwartet 3 Befunde uses-pin-tag-conflict, got %d in %+v", got, f)
	}
	for _, finding := range f {
		if finding.Reason != ReasonUsesPinTagConflict {
			continue
		}
		if !strings.Contains(finding.Message, "v3.6.0") || !strings.Contains(finding.Message, "v4.2.0") {
			t.Errorf("Message nennt nicht beide widerspruechlichen Tags: %q", finding.Message)
		}
	}
}

// Dateiuebergreifend: derselbe SHA in ZWEI Dateien mit unterschiedlichem
// Tag-Kommentar -- die erste Bedingung des Moduls, die nicht je Datei urteilt.
func TestWorkflowsPinTagConflictCrossFile(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n      - uses: docker/login-action@" + sha1 + " # v4.2.0\n",
		".github/workflows/b.yml": "jobs:\n  c:\n    steps:\n      - uses: docker/login-action@" + sha1 + " # v3.6.0\n",
	}
	f := wfRun(files)
	if got := countReason(f, ReasonUsesPinTagConflict); got != 2 {
		t.Fatalf("erwartet 2 Befunde (je Datei einen), got %d in %+v", got, f)
	}
	seen := map[string]bool{}
	for _, finding := range f {
		if finding.Reason == ReasonUsesPinTagConflict {
			seen[finding.File] = true
		}
	}
	if !seen[".github/workflows/a.yml"] || !seen[".github/workflows/b.yml"] {
		t.Fatalf("erwartet je einen Befund in a.yml UND b.yml, got %+v", f)
	}
}

// Wiederholung ist kein Befund: derselbe SHA mit IDENTISCHEM Tag-Kommentar an
// fuenf Stellen ueber zwei Dateien bleibt befundfrei fuer diesen Grund-Code.
func TestWorkflowsPinTagConflictWiederholungKeinBefund(t *testing.T) {
	line := "      - uses: actions/checkout@" + sha1 + " # v4.2.0\n"
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n" + line + line + line,
		".github/workflows/b.yml": "jobs:\n  c:\n    steps:\n" + line + line,
	}
	f := wfRun(files)
	if got := countReason(f, ReasonUsesPinTagConflict); got != 0 {
		t.Fatalf("identischer Kommentar an fuenf Stellen: erwartet 0 Befunde, got %d in %+v", got, f)
	}
}

// Verdrahtung (BEO-023): ein Konflikt, der NUR ueber den vollen Einstiegspunkt
// CheckWorkflows sichtbar wird, beweist, dass die Sammel-Schleife tatsaechlich
// an checkTagConflicts angeschlossen ist -- nicht nur, dass die Gruppierung
// fuer sich isoliert korrekt waere.
func TestWorkflowsPinTagConflictVerdrahtung(t *testing.T) {
	files := map[string]string{
		".github/workflows/a.yml": "jobs:\n  b:\n    steps:\n      - uses: actions/checkout@" + sha1 + " # v4.2.0\n",
		".github/workflows/b.yml": "jobs:\n  c:\n    steps:\n      - uses: actions/checkout@" + sha1 + " # v3.6.0\n",
	}
	if f := wfRun(files); !hasReason(f, ReasonUsesPinTagConflict) {
		t.Fatalf("CheckWorkflows liefert den Konflikt nicht -- checkTagConflicts ist nicht angeschlossen, got %+v", reasons(f))
	}
}

// checkTagConflicts direkt: Determinismus (DC-QA-02) -- zwei SHA-Gruppen,
// Ausgabe sortiert nach SHA, unabhaengig von der Einfuege-Reihenfolge.
func TestCheckTagConflictsSortiert(t *testing.T) {
	pinned := []pinnedTag{
		{file: "b.yml", line: 1, sha: "b000000000000000000000000000000000000b", value: "x@b", tag: "v2"},
		{file: "b.yml", line: 2, sha: "b000000000000000000000000000000000000b", value: "x@b", tag: "v3"},
		{file: "a.yml", line: 1, sha: "a000000000000000000000000000000000000a", value: "y@a", tag: "v1"},
		{file: "a.yml", line: 2, sha: "a000000000000000000000000000000000000a", value: "y@a", tag: "v9"},
	}
	got := checkTagConflicts(pinned)
	if len(got) != 4 {
		t.Fatalf("erwartet 4 Befunde, got %d: %+v", len(got), got)
	}
	if got[0].File != "a.yml" || got[1].File != "a.yml" {
		t.Fatalf("erwartet SHA a... zuerst (sortiert), got %+v", got)
	}
}
