package configyaml_test

import (
	"fmt"
	"os"
	"regexp"
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
		"citations", "file",
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
	// BEO-ALL/modulliste-spiegel-ungegated: die
	// beiden Schleifen oben prüfen nur netlessDocModules() ⊆ modules und
	// forbiddenInNetless() ∩ modules = ∅ -- ein NEUES, weder gelistetes
	// noch verbotenes Modul in modules fiele durch beide Maschen. Diese
	// dritte Richtung verlangt eine bewusste Einordnung jedes Moduls.
	known := map[string]bool{}
	for _, m := range netlessDocModules() {
		known[m] = true
	}
	for _, m := range forbiddenInNetless() {
		known[m] = true
	}
	for _, m := range modules {
		if !known[m] {
			return fmt.Errorf("modules aktiviert %q — weder in netlessDocModules() noch in forbiddenInNetless() eingeordnet", m)
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
// ADR-0048). Es trägt bewusst NICHT den vollen Netzlos-Doku-Satz —
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

// disableTokenRE liest einen einzelnen "--disable <modul>"-Token.
var disableTokenRE = regexp.MustCompile(`--disable ([a-z]+)`)

// TestFocusDisable_DecktDCheckYmlModules (BEO-ALL/modulliste-spiegel-ungegated):
// FOCUS_DISABLE (Makefile) spiegelt bewusst NICHT model.ValidModules(), sondern die
// .d-check.yml-modules-Liste (Makefile-Kommentar: "Spiegelt die
// .d-check.yml-modules-Liste; wächst die dort, hier nachziehen") -- ein
// anderer Fundort-Typ als die drei ValidModules()-Spiegel in
// registry_mirror_test.go (Paket app). Fail-closed: unlesbares Makefile,
// keine FOCUS_DISABLE-Definition oder ein leerer Treffer brechen den Test.
func TestFocusDisable_DecktDCheckYmlModules(t *testing.T) {
	mk, err := os.ReadFile("../../../../Makefile")
	if err != nil {
		t.Fatalf("Makefile nicht lesbar (fail-closed): %v", err)
	}
	block := regexp.MustCompile(`(?s)FOCUS_DISABLE := (.*?)\nadr-check:`).FindSubmatch(mk)
	if block == nil {
		t.Fatal("FOCUS_DISABLE-Definition im Makefile nicht gefunden (fail-closed)")
	}
	var focus []string
	for _, m := range disableTokenRE.FindAllSubmatch(block[1], -1) {
		focus = append(focus, string(m[1]))
	}
	if len(focus) == 0 {
		t.Fatal("FOCUS_DISABLE traegt keinen --disable-Token (fail-closed)")
	}

	raw, err := os.ReadFile("../../../../.d-check.yml")
	if err != nil {
		t.Fatalf(".d-check.yml nicht lesbar (fail-closed): %v", err)
	}
	cfg, err := configyaml.Decode(raw)
	if err != nil {
		t.Fatalf(".d-check.yml dekodiert nicht (fail-closed): %v", err)
	}

	for _, issue := range focusDisableIssues(focus, cfg.Modules) {
		t.Error(issue)
	}
}

// focusDisableIssues vergleicht FOCUS_DISABLE-Tokens gegen die
// .d-check.yml-modules-Liste in beide Richtungen. Reine Funktion (kein
// *testing.T) -- so treffen Live-Prüfung und Guard-Test denselben Code.
func focusDisableIssues(focus, ymlModules []string) []string {
	var issues []string
	focusSet := map[string]bool{}
	for _, m := range focus {
		focusSet[m] = true
	}
	ymlSet := map[string]bool{}
	for _, m := range ymlModules {
		ymlSet[m] = true
	}
	for _, m := range ymlModules {
		if !focusSet[m] {
			issues = append(issues, fmt.Sprintf("FOCUS_DISABLE fehlt %q aus .d-check.yml modules", m))
		}
	}
	for _, m := range focus {
		if !ymlSet[m] {
			issues = append(issues, fmt.Sprintf("FOCUS_DISABLE nennt %q, das .d-check.yml modules nicht (mehr) kennt", m))
		}
	}
	return issues
}

// TestFocusDisableIssues_Guards verriegelt focusDisableIssues gegen die
// historischen Fehlermodi: ein in .d-check.yml neues, in FOCUS_DISABLE
// fehlendes Modul und ein in FOCUS_DISABLE verwaistes (in .d-check.yml
// entferntes) Modul müssen je einen Befund erzeugen, der intakte Satz
// keinen.
func TestFocusDisableIssues_Guards(t *testing.T) {
	yml := []string{"links", "anchors", "ids"}
	if issues := focusDisableIssues(yml, yml); len(issues) != 0 {
		t.Fatalf("intakter Satz erzeugt Befund: %v", issues)
	}
	if issues := focusDisableIssues(yml[:2], yml); len(issues) == 0 {
		t.Fatal("neues .d-check.yml-Modul ohne FOCUS_DISABLE-Nachzug erzeugte keinen Befund")
	}
	if issues := focusDisableIssues(append(append([]string(nil), yml...), "entfernt"), yml); len(issues) == 0 {
		t.Fatal("verwaister FOCUS_DISABLE-Token erzeugte keinen Befund")
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
		{"unbekanntes Modul gesetzt (slice-238)", append(append([]string(nil), full...), "mentions"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := assertNetlessModules(tc.modules); (err != nil) != tc.wantErr {
				t.Fatalf("assertNetlessModules(%v) err=%v, wantErr=%v", tc.modules, err, tc.wantErr)
			}
		})
	}
}
