package app

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// BEO-ALL/modulliste-spiegel-ungegated (seit slice-238): mehrere Prosa-/
// Doku-Fundorte spiegeln model.ValidModules() als wörtliche, vollständige
// Liste -- diese Datei hält sie mechanisch nach. Zwei Erkennungsformen: eine
// ankernde Umgebungs-Phrase mit eingegrenzter Backtick-Liste (Lastenheft-Sätze,
// Optionen-Tabellenzelle -- ein zweiter Backtick-Token auf derselben Zeile,
// etwa eine referenzierte DC-ID, würde sonst mitgezählt) und ein
// whole-file-Zeilenmuster für Listen, in denen jede Zeile genau einen
// Modulnamen einleitet (README-Bullets, Handbuch-Tabelle) -- dort kollidiert
// kein anderer Backtick-Token mit dem Muster (geprüft: Trefferzahl je Datei
// entspricht genau der Modul-Anzahl). Bewusst NICHT gedeckt: die
// Bereichskürzel-Liste (Abkürzungen, keine wörtlichen Modulnamen) und die
// "fixe Standard-Modulset"-Stellen (ai-harness-Gerüst, eine bewusste
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
// Guard-Test denselben Code.
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

// bulletModuleLineRE liest eine README-Bullet-Zeile, die genau einen
// Modulnamen einleitet (`- `modul` — …`). Whole-file-Muster: die Trefferzahl
// in README.md/README.de.md entspricht genau len(model.ValidModules()) --
// kein anderer Bullet im Dokument trifft dieselbe Form.
var bulletModuleLineRE = regexp.MustCompile(`(?m)^- ` + "`" + `([a-z]+)` + "`" + ` — `)

func TestModulRegistrySpiegel_ReadmeEn(t *testing.T) {
	content := readRepoFile(t, "README.md")
	assertModuleMirrorComplete(t, "README.md Modul-Bullets", moduleLinesIn(content, bulletModuleLineRE))
}

func TestModulRegistrySpiegel_ReadmeDe(t *testing.T) {
	content := readRepoFile(t, "README.de.md")
	assertModuleMirrorComplete(t, "README.de.md Modul-Bullets", moduleLinesIn(content, bulletModuleLineRE))
}

// tableModuleRowRE liest eine Handbuch-Tabellenzeile, deren erste Zelle
// genau einen Modulnamen trägt (`| `modul` …`).
var tableModuleRowRE = regexp.MustCompile("(?m)^\\| `([a-z]+)`")

func TestModulRegistrySpiegel_BenutzerhandbuchTabelle(t *testing.T) {
	content := readRepoFile(t, filepath.Join("docs", "user", "benutzerhandbuch.md"))
	assertModuleMirrorComplete(t, "docs/user/benutzerhandbuch.md §6 Regelmodule", moduleLinesIn(content, tableModuleRowRE))
}

// moduleLinesIn wendet re (genau eine Capture-Gruppe je Zeile) über den
// gesamten Text an.
func moduleLinesIn(content string, re *regexp.Regexp) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		out = append(out, m[1])
	}
	return out
}
