package rules

import (
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
