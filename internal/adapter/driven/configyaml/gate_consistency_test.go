package configyaml_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

// DC-QA-03 / ADR-0032 (slice-064): das netzlose `make doc-check` ist die
// DC-QA-03-Messmethode. Diese Prüfung — der Rest des entfernten
// tools/gate-consistency.sh, jetzt als getippter Go-Test statt Shell-grep —
// stellt sicher, dass die Live-.d-check.yml genau diesen Netzlos-Modulsatz führt:
// alle netzlosen Doku-Module präsent, kein Netz-/Range-Modul (external/vcs).

// netlessDocModules ist der Modulsatz, mit dem der netzlose doc-check DC-QA-03
// beweist. An die Messmethode gekoppelt (nicht die alte 5-Modul-Skript-Teilmenge,
// ADR-0032 R1-F-6) — fällt spans/hostpaths/versions aus der Liste, wird die
// Prüfung rot. Als Funktion (frische Slice je Aufruf) statt Package-Var —
// gochecknoglobals + kein geteilter Zustand.
func netlessDocModules() []string {
	return []string{
		"links", "anchors", "ids", "matrix", "codepaths",
		"spans", "hostpaths", "versions", "structure", "diagrams",
		"citations",
	}
}

// forbiddenInNetless: Module, die die Netzlos-/Baum-Scan-Beweisaussage brechen —
// external und sources (beide Netzzugriff) und vcs (braucht eine Commit-Range,
// kein Baum-Scan).
func forbiddenInNetless() []string {
	return []string{"external", "sources", "vcs"}
}

// assertNetlessModules kapselt die Invariante, damit Live-Prüfung und
// Mutations-Regressionstest denselben Guard treffen (slice-057-R3-Lehre: nur der
// Guard löst den Befund aus).
func assertNetlessModules(modules []string) error {
	set := map[string]bool{}
	for _, m := range modules {
		set[m] = true
	}
	for _, m := range netlessDocModules() {
		if !set[m] {
			return fmt.Errorf("modules ohne %q — der Netzlos-Lauf beweist DC-QA-03 nur mit allen netzlosen Doku-Modulen", m)
		}
	}
	for _, m := range forbiddenInNetless() {
		if set[m] {
			return fmt.Errorf("modules aktiviert %q — das Netzlos-Gate darf kein Netz-/Range-Modul tragen (DC-QA-03)", m)
		}
	}
	return nil
}

// TestQA03_NetlessModuleList_Live liest die Repo-.d-check.yml (Relativ-Pfad-Muster
// wie docexamples_test.go / diagnose_test.go), dekodiert typisiert und prüft den
// Netzlos-Modulsatz. Fail-closed: fehlende/undekodierbare Datei ⇒ Test rot.
func TestQA03_NetlessModuleList_Live(t *testing.T) {
	raw, err := os.ReadFile("../../../../.d-check.yml")
	if err != nil {
		t.Fatalf(".d-check.yml nicht lesbar (fail-closed): %v", err)
	}
	cfg, err := configyaml.Decode(raw)
	if err != nil {
		t.Fatalf(".d-check.yml dekodiert nicht (fail-closed): %v", err)
	}
	if err := assertNetlessModules(cfg.Modules); err != nil {
		t.Fatalf("DC-QA-03-Netzlos-Modulliste verletzt: %v\nmodules = %v", err, cfg.Modules)
	}
}

// TestQA03_ClosureProfil_KeineZweiteNetzTuer prüft das ZWEITE Prüf-Profil
// (.d-check.closure.yml, gefahren von `make verify-closure-notes` über --config,
// DC-FA-CLI-012/ADR-0048). Es trägt bewusst NICHT den vollen Netzlos-Doku-Satz —
// es ist ein fokussiertes Profil, das nur `planning` per Kommandozeile
// dazuschaltet. Die Invariante, die hier zählt, ist die andere Hälfte: eine
// zweite Config-Datei darf keine zweite **Netz-Tür** aufmachen. Ohne diesen Test
// bliebe genau das ungeprüft (ADR-0048 §Konsequenzen).
func TestQA03_ClosureProfil_KeineZweiteNetzTuer(t *testing.T) {
	raw, err := os.ReadFile("../../../../.d-check.closure.yml")
	if err != nil {
		t.Fatalf(".d-check.closure.yml nicht lesbar (fail-closed): %v", err)
	}
	cfg, err := configyaml.Decode(raw)
	if err != nil {
		t.Fatalf(".d-check.closure.yml dekodiert nicht (fail-closed): %v", err)
	}
	set := map[string]bool{}
	for _, m := range cfg.Modules {
		set[m] = true
	}
	for _, m := range forbiddenInNetless() {
		if set[m] {
			t.Errorf("Closure-Profil aktiviert %q — eine zweite Config darf keine zweite Netz-/Range-Tür sein (DC-QA-03)", m)
		}
	}
	// Und es muss die Fähigkeit tatsächlich scharf schalten — ein Profil ohne
	// closure.dir wäre ein stilles Grün: `make verify-closure-notes` liefe
	// erfolgreich, ohne irgendetwas zu prüfen.
	if cfg.Planning.Closure.Dir == "" {
		t.Error("Closure-Profil ohne planning.closure.dir — das Gate liefe leer (stilles Grün)")
	}
}

// TestQA03_NetlessModuleList_Guards verriegelt den Guard gegen die historischen
// Fehlermodi: ein fehlendes Netzlos-Modul und ein gesetztes Netz-/Range-Modul
// müssen die Invariante rot machen, der intakte Satz nicht. Synthetische
// Modul-Listen treffen ausschließlich den Guard.
func TestQA03_NetlessModuleList_Guards(t *testing.T) {
	full := netlessDocModules()
	cases := []struct {
		name    string
		modules []string
		wantErr bool
	}{
		{"intakt", full, false},
		{"links fehlt", full[1:], true},
		{"structure fehlt", full[:len(full)-1], true},
		{"external gesetzt", append(append([]string(nil), full...), "external"), true},
		{"sources gesetzt", append(append([]string(nil), full...), "sources"), true},
		{"vcs gesetzt", append(append([]string(nil), full...), "vcs"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := assertNetlessModules(tc.modules); (err != nil) != tc.wantErr {
				t.Fatalf("assertNetlessModules(%v) err=%v, wantErr=%v", tc.modules, err, tc.wantErr)
			}
		})
	}
}
