package rules

// Modul mentions (DC-FA-MENT-001): Erwähnungs-Deckung. Eine über Pfad-Globs
// konfigurierte SOLL-Menge von Artefakten wird gegen eine ebenso konfigurierte
// IST-Menge von Dokumenten gehalten; gemeldet wird jedes Mitglied, das in
// KEINEM Dokument vorkommt.
//
// ACHSE, abgegrenzt (ADR-0084): Die RTM misst, ob eine Anforderung VERFOLGT
// ist, und arbeitet ueber Kennungen. Dieses Modul misst, ob ein Artefakt
// ERWAEHNT ist, und arbeitet ueber Pfade -- es braucht auf keiner Seite eine
// Kennung und ist damit unabhaengig vom ID-Schema des Repos. Ebenso wenig ist
// es targets: jenes vergleicht Doku gegen Build-Regeln, dieses gegen den
// Datei-Bestand.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): Das Modul liest zwei Eingaben und
// scannt nur eine. Die IST-Dokumente werden als Text gelesen -- dort gilt die
// Zusage. Die SOLL-Artefakte werden NIE GEOEFFNET: sie werden aus dem
// Dateibaum als Pfade aufgesammelt, und ihre Mitgliedschaft ist eine Aussage
// ueber das Verzeichnis, nicht ueber ihren Inhalt. Ein Artefakt, das existiert,
// aber leer oder unlesbar ist, bleibt Mitglied der SOLL-Menge.
//
// FAIL-CLOSED auf beiden Seiten: eine leere SOLL- oder IST-Menge ist Exit 2,
// nicht "0 Befunde" -- ein Lauf ueber null Mitglieder oder gegen null Dokumente
// behauptete Deckung, ohne etwas geprueft zu haben.

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonArtifactUnmentioned ist der Grund-Code des Moduls mentions
// (spec/spezifikation.md §4).
const ReasonArtifactUnmentioned = "artifact-unmentioned"

// MentionsResult traegt die Befunde und die Bezugsmenge des Laufs -- letztere
// gehoert in die ZUSAMMENFASSUNG, nicht in ein Befund-Feld: DC-FA-CLI-004 sagt
// den Wortlaut des message-Feldes ausdruecklich nicht zu.
type MentionsResult struct {
	Findings  []model.Finding
	Artifacts int
	Mentioned int
	Documents int
}

// Note formuliert die Bezugsmenge fuer die Zusammenfassung (DC-FA-MENT-001).
func (r MentionsResult) Note() string {
	return fmt.Sprintf("mentions: %d von %d Artefakt(en) erwähnt, über %d Dokument(e)",
		r.Mentioned, r.Artifacts, r.Documents)
}

// CheckMentions ist das Regelmodul mentions (DC-FA-MENT-001): hermetisch (nur
// Filesystem-Port, kein git, kein Netz). Der error-Rueckgabewert traegt die
// fail-closed-Faelle; der Aufrufer mappt ihn auf Exit 2.
func CheckMentions(fsys driven.Filesystem, cfg model.MentionsConfig) (MentionsResult, error) {
	var res MentionsResult
	if len(cfg.Artifacts) == 0 || len(cfg.Documents) == 0 {
		return res, fmt.Errorf("das Modul mentions braucht mentions.artifacts UND mentions.documents (DC-FA-MENT-001, fail-closed)")
	}
	artifacts := mentionsResolve(fsys, cfg.Artifacts)
	if len(artifacts) == 0 {
		return res, fmt.Errorf("mentions.artifacts %v trifft kein Artefakt — eine Deckungs-Aussage ueber null Mitglieder ist keine (DC-FA-MENT-001, fail-closed)", cfg.Artifacts)
	}
	documents := mentionsResolve(fsys, cfg.Documents)
	if len(documents) == 0 {
		return res, fmt.Errorf("mentions.documents %v trifft kein Dokument — eine Deckungs-Aussage gegen null Dokumente ist keine (DC-FA-MENT-001, fail-closed)", cfg.Documents)
	}

	corpus := mentionsCorpus(fsys, documents)
	target := strings.Join(cfg.Documents, ",")
	basename := cfg.EffectiveMatch() == model.MentionsMatchBasename
	res.Artifacts, res.Documents = len(artifacts), len(documents)
	for _, a := range artifacts {
		needle := a
		if basename {
			needle = path.Base(a)
		}
		if strings.Contains(corpus, needle) {
			res.Mentioned++
			continue
		}
		// Line 1 ist ein VERTRAGS-PLATZHALTER, keine Fundstelle: das Artefakt
		// wird nicht geoeffnet, der Befund hat keine Zeile, und diese Regel
		// sagt auch keine zu (DC-FA-CLI-004 verlangt <pfad>:<zeile> fuer jeden
		// Befund).
		res.Findings = append(res.Findings, model.Finding{
			File: a, Line: 1, Rule: "mentions", Target: target,
			Reason:  ReasonArtifactUnmentioned,
			Message: fmt.Sprintf("kein Vorkommen als %q in der Ist-Menge", needle),
		})
	}
	return res, nil
}

// mentionsResolve sammelt die Dateien, die mindestens eines der Globs treffen
// -- rekursiv ueber den Filesystem-Port, stabil sortiert und dedupliziert
// (DC-QA-02). Verzeichnisse sind keine Mitglieder: die Erwaehnungs-Frage gilt
// einem Artefakt, und ein Verzeichnis ist keines.
func mentionsResolve(fsys driven.Filesystem, globs []string) []string {
	seen := map[string]bool{}
	var all []string
	mentionsWalk(fsys, "", &all)
	var out []string
	for _, rel := range all {
		if seen[rel] || !matchAnyGlob(globs, rel) {
			continue
		}
		seen[rel] = true
		out = append(out, rel)
	}
	sort.Strings(out)
	return out
}

// mentionsWalk laeuft den Baum ab der Scan-Wurzel ab und sammelt Dateien --
// dieselbe Skip-Liste wie die Markdown-Discovery (scan.go), damit .git und
// Build-Verzeichnisse nicht zu Mitgliedern werden.
func mentionsWalk(fsys driven.Filesystem, dir string, out *[]string) {
	entries, err := fsys.List(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		rel := e.Name
		if dir != "" {
			rel = dir + "/" + e.Name
		}
		switch e.Kind {
		case driven.KindDir:
			if isSkipDir(e.Name) {
				continue
			}
			mentionsWalk(fsys, rel, out)
		case driven.KindFile:
			*out = append(*out, rel)
		}
	}
}

// mentionsCorpus liest die Ist-Menge EINMAL und haelt sie als einen String --
// die Frage lautet "kommt in IRGENDEINEM Dokument vor", also ist die
// Vereinigung die richtige Form. Ein unlesbares Dokument traegt nichts bei; es
// zaehlt aber als Mitglied der Ist-Menge, weil seine Existenz die
// Mengen-Aussage traegt und nicht sein Inhalt.
func mentionsCorpus(fsys driven.Filesystem, documents []string) string {
	var b strings.Builder
	for _, d := range documents {
		content, err := fsys.ReadFile(d)
		if err != nil {
			continue
		}
		b.Write(content)
		b.WriteByte('\n')
	}
	return b.String()
}
