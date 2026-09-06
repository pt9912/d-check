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
	"unicode"
	"unicode/utf8"

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
func CheckMentions(fsys driven.Filesystem, cfg model.MentionsConfig, ignore []string) (MentionsResult, error) {
	var res MentionsResult
	if len(cfg.Artifacts) == 0 || len(cfg.Documents) == 0 {
		return res, fmt.Errorf("das Modul mentions braucht mentions.artifacts UND mentions.documents (DC-FA-MENT-001, fail-closed)")
	}
	// EIN Baum-Durchlauf fuer beide Mengen: die gesammelte Datei-Liste ist
	// dieselbe, nur die Filter unterscheiden sich.
	var all []string
	if err := mentionsWalk(fsys, "", ignore, &all); err != nil {
		return res, err
	}
	artifacts := mentionsFilter(all, cfg.Artifacts)
	if len(artifacts) == 0 {
		return res, fmt.Errorf("mentions.artifacts %v trifft kein Artefakt — eine Deckungs-Aussage ueber null Mitglieder ist keine (DC-FA-MENT-001, fail-closed)", cfg.Artifacts)
	}
	documents := mentionsFilter(all, cfg.Documents)
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
		if mentionsOccurs(corpus, needle) {
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

// mentionsFilter waehlt aus der gesammelten Datei-Liste die Mitglieder einer
// Menge: alles, was mindestens eines der Globs trifft -- stabil sortiert und
// dedupliziert (DC-QA-02). Verzeichnisse kommen in der Liste nicht vor: die
// Erwaehnungs-Frage gilt einem Artefakt, und ein Verzeichnis ist keines.
func mentionsFilter(all, globs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, rel := range all {
		if seen[rel] || !mentionsMatchAny(globs, rel) {
			continue
		}
		seen[rel] = true
		out = append(out, rel)
	}
	sort.Strings(out)
	return out
}

// mentionsWalk laeuft den Baum ab der Repo-Wurzel ab und sammelt Dateien --
// unter DERSELBEN Pruen-Regel wie die Markdown-Discovery (scan.go): die feste
// Skip-Liste UND scan.ignore. Ohne die zweite bekaeme ein Adopter einen
// bewusst ausgenommenen Fremdbaum ueber ein weites artifacts-Glob als
// Soll-Mitglieder zurueck.
//
// NICHT auf scan.roots eingeschraenkt, und das ist eine Wahl: Die Soll-Menge
// ist kein Markdown-Scan, ihre Globs SIND ihr Geltungsbereich -- ein Artefakt
// liegt typisch ausserhalb der Doku-Wurzeln (tools/, harness/). "Relativ zur
// Scan-Wurzel" meint in DC-FA-MENT-001 die Pfad-FORM, nicht eine Beschraenkung
// auf scan.roots.
//
// Ein unlesbares Verzeichnis ist ein FEHLER, kein stilles Ueberspringen: es
// verkleinerte sonst die Soll-Menge, ohne dass ein Befund oder eine Zahl das
// zeigte. Die Wurzel selbst unlesbar zu finden ist derselbe Fall.
func mentionsWalk(fsys driven.Filesystem, dir string, ignore []string, out *[]string) error {
	entries, err := fsys.List(dir)
	if err != nil {
		return fmt.Errorf("mentions: Verzeichnis %q nicht lesbar (%v) — eine unvollständige Soll-Menge ist keine (DC-FA-MENT-001, fail-closed)", dir, err)
	}
	for _, e := range entries {
		rel := e.Name
		if dir != "" {
			rel = dir + "/" + e.Name
		}
		switch e.Kind {
		case driven.KindDir:
			if isSkipDir(e.Name) || dirIgnored(rel, ignore) {
				continue
			}
			if err := mentionsWalk(fsys, rel, ignore, out); err != nil {
				return err
			}
		case driven.KindFile:
			if !ignored(rel, ignore) {
				*out = append(*out, rel)
			}
		}
	}
	return nil
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

// mentionsMatchAny prueft die Mengen-Globs mit der SEMANTIK VON scan.ignore --
// matchGlob loest `**` segmentweise ueber beliebig viele Segmente auf, blankes
// path.Match nicht. Ohne das fielen bei `tools/**/*.sh` still fuenf von elf
// Mitgliedern aus der Soll-Menge, und fail-closed griffe nicht, weil die Menge
// nicht leer ist.
func mentionsMatchAny(globs []string, rel string) bool {
	for _, g := range globs {
		if matchGlob(g, rel) {
			return true
		}
	}
	return false
}

// mentionsOccurs prueft, ob needle im Korpus als EIGENSTAENDIGE Nennung steht.
//
// Blankes strings.Contains genuegt nicht, und der Fall ist gemessen: die
// Nennung von `image-test.md` deckte das Mitglied `test.md` ab, obwohl dieses
// nirgends genannt war. Genau die Kollision, gegen die der Default `path`
// steht.
//
// Die Grenze ist ASYMMETRISCH und ENG geschnitten, weil eine weite Grenze
// erwaehnte Artefakte als unerwaehnt meldet -- ebenfalls gemessen. Blockiert
// wird nur, was WIRKLICH Namensbestandteil sein kann:
//
//   links   Buchstabe oder Ziffer (jede Schrift), `-`, `.`
//   rechts  Buchstabe oder Ziffer (jede Schrift)
//
// Was damit ERLAUBT ist und warum -- alle vier Formen am Bestand gemessen:
// `/` links traegt die `../`-relative Verlinkung, die in Markdown der
// Regelfall ist; `_` links und rechts traegt die Kursivierung `_x.md_`; `.`
// rechts traegt den Satz-Schlusspunkt („siehe x.md."); `-` rechts traegt das
// deutsche Kompositum („die x.md-Datei").
//
// GRENZE, ausgesprochen: Damit deckt `image_test.md` das Mitglied `test.md`,
// `a.sh.bak` das Mitglied `a.sh` und ein fremdes Praefix (`x/docs/a.md`) das
// Mitglied `docs/a.md`. Das ist der Preis dafuer, dass die haeufigen Formen
// zaehlen -- eine rein textuelle Regel kann Namensteil und Satzzeichen nicht
// trennen, wo dasselbe Zeichen beides ist.
func mentionsOccurs(corpus, needle string) bool {
	for i := 0; i+len(needle) <= len(corpus); {
		j := strings.Index(corpus[i:], needle)
		if j < 0 {
			return false
		}
		at := i + j
		end := at + len(needle)
		if mentionsLeftFree(corpus[:at]) && mentionsRightFree(corpus[end:]) {
			return true
		}
		i = at + 1
	}
	return false
}

// mentionsLeftFree sagt, ob vor der Fundstelle kein Namensbestandteil steht.
func mentionsLeftFree(before string) bool {
	if before == "" {
		return true
	}
	r, _ := utf8.DecodeLastRuneInString(before)
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '.'
}

// mentionsRightFree sagt, ob nach der Fundstelle kein Namensbestandteil steht.
// Enger als links: `.` und `-` sind hier Satzzeichen bzw. Kompositum-Fuge.
func mentionsRightFree(after string) bool {
	if after == "" {
		return true
	}
	r, _ := utf8.DecodeRuneInString(after)
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}
