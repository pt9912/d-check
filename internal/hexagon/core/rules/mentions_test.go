package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
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
	res, err := CheckMentions(coretest.NewMemFS(files), cfg)
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
	_, err := CheckMentions(coretest.NewMemFS(map[string]string{"docs/handbuch.md": "x\n"}), cfg)
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
	res, err := CheckMentions(coretest.NewMemFS(map[string]string{"tools/a.sh": "x\n"}), cfg)
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
		if _, err := CheckMentions(coretest.NewMemFS(map[string]string{"docs/a.md": "x"}), cfg); err == nil {
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
