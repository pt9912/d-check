package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// mnCfg ist die Standard-Konfiguration der Tests: eine Soll-Menge aus zwei
// Globs, eine Ist-Menge aus einem Dokument, Default-Erkennungsform.
func mnCfg() model.MentionsConfig {
	return model.MentionsConfig{
		Artifacts: []string{"tools/*.sh"},
		Documents: []string{"docs/handbuch.md"},
	}
}

func mnRun(t *testing.T, files map[string]string, cfg model.MentionsConfig) MentionsResult {
	t.Helper()
	res, err := CheckMentions(coretest.NewMemFS(files), cfg, nil)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	return res
}

// Happy Path (DC-FA-MENT-001): jedes Mitglied der Soll-Menge steht in
// mindestens einem Dokument ⇒ kein Befund.
func TestMentionsHappyPath(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "#!/bin/sh\n",
		"tools/b.sh":       "#!/bin/sh\n",
		"docs/handbuch.md": "Siehe `tools/a.sh` und `tools/b.sh`.\n",
	}
	res := mnRun(t, files, mnCfg())
	if len(res.Findings) != 0 {
		t.Fatalf("erwartet befundfrei, got %+v", res.Findings)
	}
	if res.Artifacts != 2 || res.Mentioned != 2 || res.Documents != 1 {
		t.Fatalf("Bezugsmenge falsch: %d/%d über %d", res.Mentioned, res.Artifacts, res.Documents)
	}
}

// Negativ (DC-FA-MENT-001): ein Mitglied ohne Vorkommen ⇒ ein Befund
// artifact-unmentioned, dessen File der Artefakt-Pfad und dessen Line der
// VERTRAGS-PLATZHALTER 1 ist (DC-FA-CLI-004).
func TestMentionsFindet(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "#!/bin/sh\n",
		"tools/fehlt.sh":   "#!/bin/sh\n",
		"docs/handbuch.md": "Nur `tools/a.sh` ist beschrieben.\n",
	}
	res := mnRun(t, files, mnCfg())
	if len(res.Findings) != 1 {
		t.Fatalf("erwartet genau einen Befund, got %+v", res.Findings)
	}
	f := res.Findings[0]
	if f.File != "tools/fehlt.sh" || f.Line != 1 || f.Reason != ReasonArtifactUnmentioned {
		t.Fatalf("Befund-Form falsch: %+v", f)
	}
	if f.Target != "docs/handbuch.md" {
		t.Fatalf("Target soll die Ist-Menge nennen, got %q", f.Target)
	}
	if res.Mentioned != 1 || res.Artifacts != 2 {
		t.Fatalf("Bezugsmenge falsch: %d von %d", res.Mentioned, res.Artifacts)
	}
}

// Boundary (DC-FA-MENT-001): die Ist-Menge wird als VEREINIGUNG gelesen — ein
// Vorkommen in EINEM von mehreren Dokumenten genuegt.
func TestMentionsIstMengeIstVereinigung(t *testing.T) {
	cfg := mnCfg()
	cfg.Documents = []string{"docs/*.md"}
	files := map[string]string{
		"tools/a.sh":     "#!/bin/sh\n",
		"docs/eins.md":   "nichts hier\n",
		"docs/zwei.md":   "nichts hier\n",
		"docs/drei.md":   "hier steht `tools/a.sh`\n",
		"docs/vier.md":   "nichts hier\n",
	}
	res := mnRun(t, files, cfg)
	if len(res.Findings) != 0 {
		t.Fatalf("ein Vorkommen genuegt, got %+v", res.Findings)
	}
	if res.Documents != 4 {
		t.Fatalf("erwartet 4 Dokumente, got %d", res.Documents)
	}
}

// Erkennungsform (DC-FA-MENT-001): dieselbe Eingabe, zwei Urteile — der
// Unterschied steht in der Konfiguration. Unter dem Default `path` ist eine
// Nennung OHNE Pfad kein Vorkommen; unter `basename` ist sie eines.
func TestMentionsErkennungsform(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "#!/bin/sh\n",
		"docs/handbuch.md": "Siehe `a.sh` (ohne Pfad).\n",
	}
	if res := mnRun(t, files, mnCfg()); len(res.Findings) != 1 {
		t.Fatalf("Default path: erwartet einen Befund, got %+v", res.Findings)
	}
	cfg := mnCfg()
	cfg.Match = model.MentionsMatchBasename
	if res := mnRun(t, files, cfg); len(res.Findings) != 0 {
		t.Fatalf("basename: erwartet befundfrei, got %+v", res.Findings)
	}
}

// Negativ, fail-closed (DC-FA-MENT-001): eine Soll-Menge, deren Glob kein
// Artefakt trifft, ist Exit 2 und NICHT "0 Befunde". Der Test prueft, dass die
// Meldung das leere Glob NENNT — ohne das waere der Fehler nicht behebbar.
func TestMentionsLeereSollMengeFailClosed(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"gibtesnicht/*.sh"}
	_, err := CheckMentions(coretest.NewMemFS(map[string]string{"docs/handbuch.md": "x\n"}), cfg, nil)
	if err == nil {
		t.Fatal("erwartet fail-closed error, got nil")
	}
	if !strings.Contains(err.Error(), "gibtesnicht/*.sh") {
		t.Fatalf("Meldung nennt das leere Glob nicht: %v", err)
	}
}

// Negativ, fail-closed (DC-FA-MENT-001): eine Ist-Menge ohne Dokument ist
// ebenso Exit 2 — eine Deckungs-Aussage gegen null Dokumente ist keine.
// WICHTIG fuer die Aussagekraft dieses Tests: die Soll-Menge ist hier NICHT
// leer, sonst schluege schon die erste Bedingung zu und der Test pruefte die
// zweite gar nicht (BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt).
func TestMentionsLeereIstMengeFailClosed(t *testing.T) {
	cfg := mnCfg()
	cfg.Documents = []string{"gibtesnicht/*.md"}
	res, err := CheckMentions(coretest.NewMemFS(map[string]string{"tools/a.sh": "x\n"}), cfg, nil)
	if err == nil {
		t.Fatalf("erwartet fail-closed error, got %+v", res)
	}
	if !strings.Contains(err.Error(), "gibtesnicht/*.md") {
		t.Fatalf("Meldung nennt das leere Glob nicht: %v", err)
	}
	if strings.Contains(err.Error(), "artifacts") {
		t.Fatalf("die erste Bedingung hat zugeschlagen, nicht die geprueft gemeinte: %v", err)
	}
}

// Negativ, fail-closed (DC-FA-MENT-001): ein aktives Modul ohne eine der
// beiden Mengen bricht ab, statt still uebersprungen zu werden.
func TestMentionsFehlenderSchluesselBrichtAb(t *testing.T) {
	for name, cfg := range map[string]model.MentionsConfig{
		"ohne artifacts": {Documents: []string{"docs/*.md"}},
		"ohne documents": {Artifacts: []string{"tools/*.sh"}},
		"beide leer":     {},
	} {
		if _, err := CheckMentions(coretest.NewMemFS(map[string]string{"docs/a.md": "x"}), cfg, nil); err == nil {
			t.Fatalf("%s: erwartet fail-closed error, got nil", name)
		}
	}
}

// Das Modul OEFFNET die Soll-Artefakte nicht (AGENTS.md §3.8): ein Artefakt,
// dessen Inhalt den eigenen Pfad traegt, ist dadurch NICHT erwaehnt — nur die
// Ist-Menge zaehlt. Ohne diesen Test waere die Grenz-Aussage im
// Modul-Kommentar eine Behauptung.
func TestMentionsScanntDieSollMengeNicht(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "# ich bin tools/a.sh\n",
		"docs/handbuch.md": "leer\n",
	}
	res := mnRun(t, files, mnCfg())
	if len(res.Findings) != 1 {
		t.Fatalf("Selbstnennung darf nicht zaehlen, got %+v", res.Findings)
	}
}

// Determinismus (DC-QA-02): identische Eingabe ⇒ identische Befund-Reihenfolge.
func TestMentionsDeterministisch(t *testing.T) {
	files := map[string]string{
		"tools/z.sh": "x", "tools/a.sh": "x", "tools/m.sh": "x",
		"docs/handbuch.md": "leer\n",
	}
	var first []string
	for i := 0; i < 5; i++ {
		res := mnRun(t, files, mnCfg())
		var got []string
		for _, f := range res.Findings {
			got = append(got, f.File)
		}
		if i == 0 {
			first = got
			continue
		}
		if strings.Join(got, "|") != strings.Join(first, "|") {
			t.Fatalf("Lauf %d weicht ab: %v vs. %v", i, got, first)
		}
	}
	if strings.Join(first, "|") != "tools/a.sh|tools/m.sh|tools/z.sh" {
		t.Fatalf("erwartet stabile Sortierung, got %v", first)
	}
}

// Die Bezugsmenge steht in der ZUSAMMENFASSUNG (DC-FA-MENT-001): ein sauberer
// Lauf sagt, worueber er nichts gefunden hat.
func TestMentionsNoteNenntDieBezugsmenge(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "x",
		"docs/handbuch.md": "`tools/a.sh`\n",
	}
	note := mnRun(t, files, mnCfg()).Note()
	for _, want := range []string{"1 von 1", "1 Dokument"} {
		if !strings.Contains(note, want) {
			t.Fatalf("Note %q enthaelt %q nicht", note, want)
		}
	}
}

// EffectiveMatch: Default ist die strengere Form (DC-FA-MENT-001).
func TestMentionsEffectiveMatch(t *testing.T) {
	if got := (model.MentionsConfig{}).EffectiveMatch(); got != model.MentionsMatchPath {
		t.Fatalf("Default soll %q sein, got %q", model.MentionsMatchPath, got)
	}
	if got := (model.MentionsConfig{Match: model.MentionsMatchBasename}).EffectiveMatch(); got != model.MentionsMatchBasename {
		t.Fatalf("gesetzte Form soll durchgereicht werden, got %q", got)
	}
}

// Unter `basename` darf eine TEILZEICHENKETTEN-Kollision kein Mitglied
// abdecken: `test.md` gilt nicht als erwaehnt, nur weil `image-test.md` im
// Korpus steht. Genau die Kollision, gegen die der
// Default `path` steht. Der Test faellt, sobald die Grenz-Pruefung entfaellt.
func TestMentionsBasenameKeineTeilzeichenkette(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"s/*.md"}
	cfg.Match = model.MentionsMatchBasename
	files := map[string]string{
		"s/test.md":        "x",
		"s/image-test.md":  "x",
		"docs/handbuch.md": "nur [image-test](s/image-test.md) ist verlinkt\n",
	}
	res := mnRun(t, files, cfg)
	if len(res.Findings) != 1 || res.Findings[0].File != "s/test.md" {
		t.Fatalf("erwartet genau s/test.md als Fund, got %+v", res.Findings)
	}
	// Gegenprobe: der Pfad-Praefix vor dem Basisnamen ist LEGITIM und darf
	// nicht mitblockieren -- sonst waere die Regel unbrauchbar.
	if res.Mentioned != 1 {
		t.Fatalf("s/image-test.md soll ueber seinen relativen Link zaehlen, got %d erwaehnt", res.Mentioned)
	}
}

// GRENZE, festgehalten statt behauptet: Ein fremdes Pfad-Praefix deckt das
// Mitglied darunter -- `x/docs/a.md` gilt als Nennung von `docs/a.md`. Das ist
// der ausgewiesene Preis dafuer, dass `/` links erlaubt ist, und `/` links
// muss erlaubt sein, damit die `../`-relative Verlinkung zaehlt. Der Test
// haelt die Wahl fest: kippt sie, faellt er, und die Grenze ist neu zu
// bewerten statt still zu wandern.
func TestMentionsPfadPraefixDecktMitglied(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"docs/a.md"}
	cfg.Documents = []string{"h.md"}
	files := map[string]string{
		"docs/a.md": "x",
		"h.md":      "hier steht nur x/docs/a.md\n",
	}
	if res := mnRun(t, files, cfg); len(res.Findings) != 0 {
		t.Fatalf("die Grenze ist gewandert: Praefix-Treffer deckt nicht mehr, got %+v", res.Findings)
	}
}

// Und die rechte Grenze: `a.sh` darf nicht durch `a.shx` gedeckt werden.
func TestMentionsRechteGrenze(t *testing.T) {
	files := map[string]string{
		"tools/a.sh":       "x",
		"docs/handbuch.md": "hier steht tools/a.shx\n",
	}
	if res := mnRun(t, files, mnCfg()); len(res.Findings) != 1 {
		t.Fatalf("rechte Grenze traegt nicht, got %+v", res.Findings)
	}
}

// Die Mengen-Globs folgen der
// SEMANTIK VON scan.ignore -- `**` loest ueber beliebig viele Segmente auf.
// Mit blankem path.Match fielen die tiefer liegenden Mitglieder still aus der
// Soll-Menge, und fail-closed griffe nicht, weil die Menge nicht leer ist.
func TestMentionsGlobDoppelsternUeberMehrereSegmente(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"tools/**/*.sh"}
	files := map[string]string{
		"tools/flach.sh":            "x",
		"tools/harness/tiefer.sh":   "x",
		"tools/a/b/noch-tiefer.sh":  "x",
		"docs/handbuch.md":          "leer\n",
	}
	res := mnRun(t, files, cfg)
	if res.Artifacts != 3 {
		t.Fatalf("erwartet 3 Mitglieder ueber beliebig viele Segmente, got %d", res.Artifacts)
	}
}

// scan.ignore prunt die Soll-Menge. Die Fixture benutzt bewusst KEINEN Namen
// aus der festen Skip-Liste -- sonst pruefte der Test jene und nicht das
// scan.ignore-Verhalten, und er koennte nicht fallen. Ohne die Auswertung
// bekaeme ein Adopter einen
// bewusst ausgenommenen Fremdbaum ueber ein weites Glob zurueck.
func TestMentionsHonoriertScanIgnore(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"**/*.sh"}
	files := map[string]string{
		"tools/eigen.sh":    "x",
		"fremdbaum/fremd.sh": "x",
		"docs/handbuch.md":  "leer\n",
	}
	res, err := CheckMentions(coretest.NewMemFS(files), cfg, []string{"fremdbaum/**"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if res.Artifacts != 1 || len(res.Findings) != 1 || res.Findings[0].File != "tools/eigen.sh" {
		t.Fatalf("fremdbaum/ soll geprunt sein, got %d Mitglieder %+v", res.Artifacts, res.Findings)
	}
}

// Die Sortierung ist real lasttragend, und die Fixture muss sie pruefen
// koennen: MemFS.List sortiert selbst und der Walk laeuft tiefen-erst. Diese
// Fixture TRENNT die beiden Ordnungen: lexikografisch steht "tools.md" VOR
// "tools/a.sh" ('.' < '/'), die Tiefensuche liefert es danach. Ohne
// sort.Strings faellt der Test.
func TestMentionsSortierungTrenntDFSVonLexikografisch(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"tools.md", "tools/*.sh"}
	files := map[string]string{
		"tools.md":         "x",
		"tools/a.sh":       "x",
		"docs/handbuch.md": "leer\n",
	}
	res := mnRun(t, files, cfg)
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.File)
	}
	if strings.Join(got, "|") != "tools.md|tools/a.sh" {
		t.Fatalf("erwartet lexikografische Ordnung, got %v", got)
	}
}

// Ein unlesbares Verzeichnis darf die Soll-Menge nicht still verkleinern.
// Geprueft ueber die WURZEL, weil MemFS keine gezielt
// unlesbaren Unterverzeichnisse kennt -- der Fehlerpfad ist derselbe.
func TestMentionsUnlesbaresVerzeichnisFailClosed(t *testing.T) {
	_, err := CheckMentions(brokenFS{}, mnCfg(), nil)
	if err == nil {
		t.Fatal("erwartet fail-closed error, got nil")
	}
	if !strings.Contains(err.Error(), "nicht lesbar") {
		t.Fatalf("Meldung nennt die Ursache nicht: %v", err)
	}
}

// brokenFS liefert fuer jedes Verzeichnis einen Lesefehler.
type brokenFS struct{}

func (brokenFS) Kind(string) (driven.EntryKind, error)      { return driven.KindMissing, errBroken }
func (brokenFS) ReadFile(string) ([]byte, error)            { return nil, errBroken }
func (brokenFS) List(string) ([]driven.DirEntry, error)     { return nil, errBroken }

var errBroken = fmt.Errorf("probe: nicht lesbar")

// Die Grenz-Pruefung darf ERWAEHNTE Artefakte nicht als unerwaehnt melden.
// Alle vier Formen sind am Bestand gemessen und in gewoehnlicher Prosa der
// Regelfall; eine zu weite Grenze machte das Modul im Default unbrauchbar.
func TestMentionsGrenzeMeldetHaeufigeFormenNicht(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"a/*.md"}
	cfg.Match = model.MentionsMatchBasename
	for name, text := range map[string]string{
		"Satz-Schlusspunkt":  "siehe x.md.\n",
		"deutsches Kompositum": "die x.md-Datei traegt es\n",
		"Kursivierung":       "_x.md_ ist kursiv\n",
		"Klammer":            "(x.md) in Klammern\n",
		"Inline-Code":        "`x.md` als Code\n",
		"deutsche Anfuehrung": "„x.md\" in Anfuehrung\n",
		"Doppelpunkt-Zeile":  "x.md:12 mit Zeilennummer\n",
		"Zeilenende":         "am Ende steht x.md",
	} {
		files := map[string]string{"a/x.md": "y", "docs/handbuch.md": text}
		if res := mnRun(t, files, cfg); len(res.Findings) != 0 {
			t.Errorf("%s: Nennung %q zaehlt nicht, got %+v", name, text, res.Findings)
		}
	}
}

// Und unter dem DEFAULT `path` traegt die `../`-relative Verlinkung, die in
// Markdown der Regelfall ist. Ohne das meldete ein Adopter mit dem ueblichen
// Layout seine gesamte Soll-Menge als Befund.
func TestMentionsPathTraegtRelativeVerlinkung(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"docs/plan/*.md"}
	cfg.Documents = []string{"h.md"}
	files := map[string]string{
		"docs/plan/a.md": "y",
		"h.md":           "[A](../docs/plan/a.md)\n",
	}
	if res := mnRun(t, files, cfg); len(res.Findings) != 0 {
		t.Fatalf("../-relative Verlinkung muss zaehlen, got %+v", res.Findings)
	}
}

// Nicht-ASCII links ist ein Namensbestandteil und muss blockieren -- eine
// byte-basierte Pruefung sah hier das zweite Byte eines Umlauts und liess
// durch. Deutsche Anfuehrungszeichen sind dagegen Satzzeichen und zaehlen
// (oben mitgeprueft).
func TestMentionsGrenzeIstRunenBasiert(t *testing.T) {
	cfg := mnCfg()
	cfg.Artifacts = []string{"a/test.md"}
	cfg.Documents = []string{"h.md"}
	files := map[string]string{
		"a/test.md": "y",
		"h.md":      "hier steht Ätest.md und sonst nichts\n",
	}
	cfg.Match = model.MentionsMatchBasename
	if res := mnRun(t, files, cfg); len(res.Findings) != 1 {
		t.Fatalf("Ätest.md darf test.md nicht decken, got %+v", res.Findings)
	}
}
