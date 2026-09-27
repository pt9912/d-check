package app

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// slice-238 (BEO-ALL/modulliste-spiegel-ungegated, 3. Evidenz): drei Prosa-
// Fundorte spiegeln model.ValidModules() als wörtliche, vollständige
// Backtick-Liste -- gefunden über eine ankernde Umgebungs-Phrase je Ort,
// damit ein zweiter Backtick-Token auf derselben Zeile (z. B. eine
// referenzierte DC-ID) nicht mitgezählt wird. Deckt genau die drei Orte, die
// slice-236s Review nachtragen musste (F-2/F-3/F-4) -- nicht die
// Bereichskürzel-Liste (Abkürzungen, keine wörtlichen Modulnamen) und nicht
// die "fixe Standard-Modulset"-Stellen (ai-harness-Gerüst, eine bewusste
// Teilmenge, kein Vollständigkeits-Anspruch).

// backtickTokenRE liest einzelne `wort`-Tokens aus einem bereits
// eingegrenzten Textabschnitt.
var backtickTokenRE = regexp.MustCompile("`([a-z]+)`")

// moduleTokensIn extrahiert die Backtick-Tokens aus s.
func moduleTokensIn(s string) []string {
	var out []string
	for _, m := range backtickTokenRE.FindAllStringSubmatch(s, -1) {
		out = append(out, m[1])
	}
	return out
}

// moduleMirrorIssues prüft in beide Richtungen: jedes Modul aus
// model.ValidModules() steht in got, und got trägt keinen Namen, den die
// Registry nicht kennt (verwaister Fundort-Eintrag nach einer Entfernung).
// Reine Funktion (kein *testing.T) -- so treffen Live-Prüfung und
// Guard-Test denselben Code (slice-057-R3-Lehre: nur der Guard löst den
// Befund aus).
func moduleMirrorIssues(got []string) []string {
	if len(got) == 0 {
		return []string{"kein Modulname gefunden (Anker-Phrase verschoben oder Regex verfehlt die Zeile)"}
	}
	var issues []string
	set := map[string]bool{}
	for _, m := range got {
		set[m] = true
	}
	for _, m := range model.ValidModules() {
		if !set[m] {
			issues = append(issues, "Modul \""+m+"\" fehlt")
		}
	}
	valid := map[string]bool{}
	for _, m := range model.ValidModules() {
		valid[m] = true
	}
	for _, m := range got {
		if !valid[m] {
			issues = append(issues, "Eintrag \""+m+"\" ist kein bekanntes Modul (verwaist?)")
		}
	}
	return issues
}

func assertModuleMirrorComplete(t *testing.T, ort string, got []string) {
	t.Helper()
	for _, issue := range moduleMirrorIssues(got) {
		t.Errorf("%s: %s", ort, issue)
	}
}

// TestModuleMirrorIssues_Guards verriegelt moduleMirrorIssues gegen die
// historischen Fehlermodi (slice-236: ein fehlendes Modul in drei Prosa-
// Fundorten) -- ein fehlendes und ein verwaistes Modul müssen je einen
// Befund erzeugen, der intakte Satz keinen.
func TestModuleMirrorIssues_Guards(t *testing.T) {
	full := model.ValidModules()
	if issues := moduleMirrorIssues(full); len(issues) != 0 {
		t.Fatalf("intakter Satz erzeugt Befund: %v", issues)
	}
	fehlend := full[1:]
	if issues := moduleMirrorIssues(fehlend); len(issues) == 0 {
		t.Fatalf("fehlendes Modul %q erzeugte keinen Befund", full[0])
	}
	verwaist := append(append([]string(nil), full...), "kein-echtes-modul")
	if issues := moduleMirrorIssues(verwaist); len(issues) == 0 {
		t.Fatal("verwaister Eintrag erzeugte keinen Befund")
	}
	if issues := moduleMirrorIssues(nil); len(issues) == 0 {
		t.Fatal("leere Liste erzeugte keinen Befund (fail-closed erwartet)")
	}
}

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", rel))
	if err != nil {
		t.Fatalf("%s nicht lesbar (fail-closed): %v", rel, err)
	}
	return string(raw)
}

func TestModulRegistrySpiegel_LastenheftCLI002(t *testing.T) {
	content := readRepoFile(t, filepath.Join("spec", "lastenheft.md"))
	re := regexp.MustCompile(`(?s)gegliedert:(.*?)\. Ohne Konfiguration`)
	m := re.FindStringSubmatch(content)
	if m == nil {
		t.Fatal("DC-FA-CLI-002-Beschreibungssatz nicht gefunden (fail-closed)")
	}
	assertModuleMirrorComplete(t, "spec/lastenheft.md DC-FA-CLI-002", moduleTokensIn(m[1]))
}

func TestModulRegistrySpiegel_LastenheftGlossar(t *testing.T) {
	content := readRepoFile(t, filepath.Join("spec", "lastenheft.md"))
	re := regexp.MustCompile(`Regelmodul \| Benannte, einzeln aktivierbare Prüf-Einheit \(([^)]+)\)\.`)
	m := re.FindStringSubmatch(content)
	if m == nil {
		t.Fatal("Glossar-Zeile 'Regelmodul' nicht gefunden (fail-closed)")
	}
	assertModuleMirrorComplete(t, "spec/lastenheft.md §6 Glossar", moduleTokensIn(m[1]))
}

func TestModulRegistrySpiegel_OperationsMdOptionen(t *testing.T) {
	content := readRepoFile(t, filepath.Join("docs", "user", "operations.md"))
	re := regexp.MustCompile(`Regelmodule zu-/abschalten \(([^)]+)\); CLI schlägt Konfiguration`)
	m := re.FindStringSubmatch(content)
	if m == nil {
		t.Fatal("Optionen-Tabellenzeile '--enable/--disable' nicht gefunden (fail-closed)")
	}
	assertModuleMirrorComplete(t, "docs/user/operations.md Optionen-Tabelle", moduleTokensIn(m[1]))
}
