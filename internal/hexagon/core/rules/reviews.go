package rules

// Modul reviews (DC-FA-RVW-001): Review-Report-Deckung. Ein `done/`-Slice mit
// Review-Zusage -- ein DoD-Haken, dessen Zeile "Review" nennt, in JEDER der
// drei CommonMark-Bullet-Formen (`-`/`*`/`+`) und unabhaengig vom Haken-Zustand
// -- braucht mindestens einen Report in reviews.reviews-dir, dessen Dateiname
// dieselbe slice-<NNN>-Kennung traegt.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): das Modul scannt done-dir und
// reviews-dir NICHT rekursiv -- ein bereits archivierter Slice (Stub unter
// done/<welle-id>/) traegt keine DoD mehr und faellt damit natuerlich aus der
// Kandidatenmenge, nicht durch Sonderfall. Geprueft wird die DECKUNG (ein
// Report existiert), nicht seine QUALITAET -- dieselbe Grenze wie beim
// DoD-Haken selbst: eine Selbstauskunft.

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonReviewMissing ist der Grund-Code des Moduls reviews
// (spec/spezifikation.md §4, SPEC-081).
const ReasonReviewMissing = "review-missing"

// reviewLineRE erkennt eine DoD-Zeile mit Review-Zusage: dieselbe
// bullet-Toleranz wie taskItemRE (structure.go, ADR-0074) -- Haken-Zustand
// zaehlt hier nicht, nur dass die Zeile die Phrase "unabhängiger Review"
// traegt (Groß-/Kleinschreibung am Wortanfang egal, gemessen an beiden Formen
// im Bestand). Bloßes "Review" ist ZU BREIT -- gemessen an slice-183: dessen
// "Adaptions-Review" ist ein anderes, in der Slice-Datei SELBST dokumentiertes
// Konzept ohne eigenen Report unter docs/reviews/, kein Review-Zusage im Sinne
// dieses Moduls.
var reviewLineRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[[ xX]\].*[Uu]nabhängiger Review`)

// sliceIDRE liest die slice-<NNN>-Kennung aus einem Dateinamen -- dieselbe
// Form wie tools/archive-wave/collect.go's sliceIDInNameRE, hier unabhaengig
// nachgebaut: das Modul liest keine Fremd-Werkzeuge (Hexagon-Schnitt).
var sliceIDRE = regexp.MustCompile(`slice-([0-9]+)`)

// CheckReviews ist das Regelmodul reviews (DC-FA-RVW-001): hermetisch (nur
// Filesystem-Port, kein git, kein Netz), opt-in ueber DoneDir.
func CheckReviews(fsys driven.Filesystem, cfg model.ReviewsConfig) []model.Finding {
	if strings.TrimSpace(cfg.DoneDir) == "" {
		return nil // inert: keine Datei wird geoeffnet
	}
	candidates := reviewCandidates(fsys, cfg)
	reviewNames, listErr := fsys.List(cfg.ReviewsDir)
	var out []model.Finding
	promises := 0
	for _, f := range candidates {
		id := sliceIDRE.FindString(path.Base(f))
		if id == "" {
			continue // reviewCandidates filtert schon auf "slice-"-Praefix; verteidigt gegen Drift
		}
		b, err := fsys.ReadFile(f)
		if err != nil {
			continue
		}
		line, ok := reviewPromise(string(b))
		if !ok {
			continue
		}
		promises++
		if !hasMatchingReview(reviewNames, id) {
			out = append(out, model.Finding{File: f, Line: line, Rule: "reviews", Target: cfg.ReviewsDir,
				Reason: ReasonReviewMissing,
				Message: fmt.Sprintf("Review-Zusage ohne Report unter %s fuer %s", cfg.ReviewsDir, id)})
		}
	}
	// FAIL-CLOSED: leere Kandidatenmenge oder unlesbares reviews-dir sehen
	// sonst identisch aus wie "alles gedeckt". Null Review-ZUSAGEN unter
	// vorhandenen Kandidaten ist dagegen ein legitimer Zustand (ein kleiner
	// oder junger Bestand kann das sein) und loest fuer sich allein KEIN
	// Fail-Closed aus -- anders als bei workflows' refs==0, wo eine
	// Workflow-Datei ohne jede uses:-Zeile ein Anomalie-Signal ist.
	if len(candidates) == 0 || listErr != nil {
		out = append(out, model.Finding{File: cfg.DoneDir, Line: 1, Rule: "reviews", Target: cfg.DoneDir,
			Reason: ReasonReviewMissing,
			Message: fmt.Sprintf("leere Pruefmenge: %d Kandidat(en), %d Review-Zusage(n), reviews-dir lesbar: %v — fail-closed",
				len(candidates), promises, listErr == nil)})
	}
	return out
}

// reviewPromise sucht die ERSTE Review-Zusage-Zeile im Dokument (roh, ueber
// alle Zeilen -- dieselbe Lexik-Entscheidung wie max-open-tasks/ADR-0074:
// Haken-Zustand zaehlt nicht, nur dass die Zeile "Review" nennt) und liefert
// ihre 1-basierte Zeilennummer fuer den Befund.
func reviewPromise(content string) (line int, ok bool) {
	lines := strings.Split(content, "\n")
	for i, l := range lines {
		if reviewLineRE.MatchString(l) {
			return i + 1, true
		}
	}
	return 0, false
}

// reviewCandidates liefert die Slice-Dateien unmittelbar in DoneDir, stabil
// sortiert, abzueglich exempt-paths -- nicht rekursiv (archivierte Stubs in
// Unterverzeichnissen sind damit ausgeschlossen, siehe Modul-Kommentar).
func reviewCandidates(fsys driven.Filesystem, cfg model.ReviewsConfig) []string {
	entries, err := fsys.List(cfg.DoneDir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.Kind != driven.KindFile {
			continue
		}
		if !strings.HasSuffix(e.Name, ".md") || !strings.HasPrefix(e.Name, "slice-") {
			continue
		}
		rel := path.Join(cfg.DoneDir, e.Name)
		if matchAnyGlob(cfg.ExemptPaths, rel) {
			continue
		}
		out = append(out, rel)
	}
	sort.Strings(out)
	return out
}

// hasMatchingReview prueft, ob mindestens ein Eintrag in reviewNames die
// slice-<NNN>-Kennung im Dateinamen traegt -- Substring-Match, dieselbe Form
// wie tools/archive-wave/collect.go's CollectReviews (1:N zulaessig, z. B.
// -r1/-r2-Suffixe).
func hasMatchingReview(reviewNames []driven.DirEntry, id string) bool {
	for _, e := range reviewNames {
		if e.Kind != driven.KindFile {
			continue
		}
		if m := sliceIDRE.FindString(e.Name); m == id {
			return true
		}
	}
	return false
}
