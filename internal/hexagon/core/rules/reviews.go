package rules

// Modul reviews (DC-FA-RVW-001): Review-Report-Deckung. Ein `done/`-Slice mit
// Review-Zusage -- ein DoD-Item, dessen TEXT (Checkbox-Zeile plus lose
// Folgezeilen bis zur naechsten Checkbox/Leerzeile/Dateiende) auf
// reviews.promise-pattern passt (abwesend: model.DefaultPromisePattern),
// in JEDER der drei CommonMark-Bullet-Formen (`-`/`*`/`+`) und unabhaengig vom
// Haken-Zustand -- braucht mindestens einen Report in reviews.reviews-dir, der
// ihn deckt: ueber dieselbe slice-<NNN>-Kennung im Dateinamen (match: id) oder
// ueber den Basisnamen des Slice ohne .md als Teil des Dateinamens
// (match: name).
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): reviews-dir wird NICHT rekursiv
// gelesen; done-dir nur mit reviews.recursive. Ohne den Schluessel faellt ein
// archivierter Slice unter done/<welle-id>/ aus der Kandidatenmenge, mit ihm
// nimmt reviews.skip-pattern den Stub an seinem Inhalt aus. Geprueft wird die
// DECKUNG (ein Report existiert), nicht seine QUALITAET -- dieselbe Grenze wie
// beim DoD-Haken selbst: eine Selbstauskunft.

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonReviewMissing ist der Grund-Code des Moduls reviews
// (spec/spezifikation.md §4, SPEC-081).
const ReasonReviewMissing = "review-missing"

// checkboxLineRE erkennt den BEGINN eines DoD-Items: dieselbe bullet-Toleranz
// wie taskItemRE (structure.go, ADR-0074) -- Haken-Zustand zaehlt nicht.
var checkboxLineRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[[ xX]\]`)

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
	// Das Muster ist am Config-Rand geprueft (Exit 2), MustCompile trifft hier
	// gueltige Muster.
	promiseRE := regexp.MustCompile(cfg.EffectivePromisePattern())
	menge := reviewCandidates(fsys, cfg)
	candidates, badDirs, skipped := menge.out, menge.badDirs, menge.skipped
	reviewNames, listErr := fsys.List(cfg.ReviewsDir)
	var out []model.Finding
	for _, d := range badDirs {
		out = append(out, model.Finding{File: cfg.DoneDir, Line: 1, Rule: "reviews", Target: d,
			Reason:  ReasonReviewMissing,
			Message: "Verzeichnis " + d + " unlesbar — fail-closed"})
	}
	promises := 0
	for _, f := range candidates {
		finding, isPromise := reviewFinding(fsys, cfg, promiseRE, reviewNames, menge.names, f)
		if isPromise {
			promises++
		}
		if finding != nil {
			out = append(out, *finding)
		}
	}
	lage := reviewLage{candidates: len(candidates), badDirs: len(badDirs), skipped: skipped,
		promises: promises, listOK: listErr == nil}
	if lf := reviewLeerlauf(cfg, lage); lf != nil {
		out = append(out, *lf)
	}
	return out
}

// reviewLage fasst zusammen, worueber der Leerlauf-Befund entscheidet.
type reviewLage struct {
	candidates, badDirs, skipped, promises int
	listOK                                 bool
}

// reviewLeerlauf liefert den Leerlauf-Befund auf reviews.done-dir (nil ⇒
// keiner).
func reviewLeerlauf(cfg model.ReviewsConfig, l reviewLage) *model.Finding {
	// FAIL-CLOSED: leere Kandidatenmenge, oder unlesbares reviews-dir OHNE dass
	// eine einzige Review-Zusage vorliegt, saehen sonst identisch aus wie
	// "alles gedeckt". Ein unlesbares reviews-dir MIT vorhandenen Zusagen
	// braucht diese Zeile NICHT zusaetzlich: jede Zusage hat dann bereits ihren
	// eigenen `review-missing`-Befund oben ausgeloest (die Zuordnung kann gegen
	// eine leere Namens-Liste nie treffen). Null Review-ZUSAGEN unter
	// vorhandenen Kandidaten ist ohne require-promises ein legitimer Zustand
	// (ein kleiner oder junger Bestand kann das sein); mit dem Schluessel ist
	// er ein Befund, weil die Pruefung sonst ueber nichts gruen meldet. Neben
	// einem unlesbaren Unterverzeichnis entfaellt der Leerlauf-Befund: die
	// Menge ist dann nicht leer, sondern unvollstaendig, und das ist gemeldet.
	// Leert erst skip-pattern die Menge und erklaert skip-allows-empty das zum
	// Ruhezustand eines Repos, dessen Slices alle archiviert sind, ist das kein
	// Befund; require-promises zaehlt ohnehin nur die uebrig gelassenen
	// Kandidaten.
	switch {
	case l.badDirs > 0:
		return nil
	case l.candidates == 0 && l.skipped > 0 && cfg.SkipAllowsEmpty && l.listOK:
		return nil
	case l.candidates == 0 || (!l.listOK && l.promises == 0):
		return &model.Finding{File: cfg.DoneDir, Line: 1, Rule: "reviews", Target: cfg.DoneDir,
			Reason: ReasonReviewMissing,
			Message: fmt.Sprintf("leere Pruefmenge: %d Kandidat(en), %d Review-Zusage(n), reviews-dir lesbar: %v — fail-closed",
				l.candidates, l.promises, l.listOK)}
	case cfg.RequirePromises && l.promises == 0:
		return &model.Finding{File: cfg.DoneDir, Line: 1, Rule: "reviews", Target: cfg.DoneDir,
			Reason: ReasonReviewMissing,
			Message: fmt.Sprintf("keine Review-Zusage unter %d Kandidat(en) — require-promises: das Muster %q trifft keinen DoD-Punkt",
				l.candidates, cfg.EffectivePromisePattern())}
	}
	return nil
}

// reviewFinding prueft einen Kandidaten: traegt er eine Zusage, und deckt ein
// Report sie? Liefert den Befund (nil ⇒ keiner) und ob eine Zusage vorliegt.
// Eine unlesbare Datei faellt still aus -- Bestand, nicht Gegenstand dieser
// Zuordnung.
func reviewFinding(
	fsys driven.Filesystem, cfg model.ReviewsConfig, promiseRE *regexp.Regexp,
	reviewNames []driven.DirEntry, names []string, f string,
) (*model.Finding, bool) {
	base := path.Base(f)
	var id string
	if !cfg.MatchByName() {
		id = sliceIDRE.FindString(base)
	}
	b, err := fsys.ReadFile(f)
	if err != nil {
		return nil, false
	}
	line, ok := reviewPromise(string(b), promiseRE)
	if !ok {
		return nil, false
	}
	switch {
	case cfg.MatchByName():
		key := strings.TrimSuffix(base, ".md")
		if hasReviewContaining(reviewNames, key, longerNames(names, key)) {
			return nil, true
		}
		return &model.Finding{File: f, Line: line, Rule: "reviews", Target: cfg.ReviewsDir,
			Reason:  ReasonReviewMissing,
			Message: fmt.Sprintf("Review-Zusage ohne Report unter %s fuer %s", cfg.ReviewsDir, key)}, true
	case id == "":
		// Eine Zusage, deren Kennung nicht lesbar ist, faellt nicht still aus:
		// sie waere sonst eine ungepruefte Zusage.
		return &model.Finding{File: f, Line: line, Rule: "reviews", Target: cfg.ReviewsDir,
			Reason: ReasonReviewMissing,
			Message: "Review-Zusage, aber keine slice-<NNN>-Kennung im Dateinamen " + base +
				" — match: name ordnet über den Basisnamen zu"}, true
	case hasMatchingReview(reviewNames, id):
		return nil, true
	default:
		return &model.Finding{File: f, Line: line, Rule: "reviews", Target: cfg.ReviewsDir,
			Reason:  ReasonReviewMissing,
			Message: fmt.Sprintf("Review-Zusage ohne Report unter %s fuer %s", cfg.ReviewsDir, id)}, true
	}
}

// reviewPromise sucht das ERSTE DoD-Item, dessen TEXT das Zusage-Muster
// traegt, und liefert die 1-basierte Zeilennummer seines Checkbox-Starts.
//
// Ein Item ist NICHT auf seine Checkbox-Zeile beschraenkt: lange DoD-Punkte
// laufen ueber mehrere Zeilen, und die Zusage kann auf einer FOLGEZEILE
// stehen. Eine Item-Grenze ist deshalb der Bereich von einer Checkbox-Zeile
// bis ausschließlich der naechsten Checkbox-Zeile, einer Leerzeile oder dem
// Dateiende -- dieselbe Grenze, an der ein loses Markdown-Listenelement endet.
// GRENZE: jeder Checkbox-Punkt der Datei zaehlt, nicht nur der im
// DoD-Abschnitt, auch einer in einem Codeblock; das Muster sieht den Punkt ab
// dem Bullet, ein ^ verankert also vor "- [ ]".
func reviewPromise(content string, promiseRE *regexp.Regexp) (line int, ok bool) {
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		if !checkboxLineRE.MatchString(lines[i]) {
			continue
		}
		item := lines[i]
		for j := i + 1; j < len(lines) && strings.TrimSpace(lines[j]) != "" && !checkboxLineRE.MatchString(lines[j]); j++ {
			item += "\n" + lines[j]
		}
		if promiseRE.MatchString(item) {
			return i + 1, true
		}
	}
	return 0, false
}

// reviewCandidates liefert die Slice-Dateien in DoneDir -- mit recursive auch
// in seinen Unterverzeichnissen (SKIP_DIRS ausgenommen, ein Symlink auf ein
// Verzeichnis wird nicht verfolgt) --, stabil sortiert, abzueglich
// exempt-paths und der Dateien, deren Inhalt skip-pattern trifft. badDirs
// nennt die unlesbaren UNTERverzeichnisse, sortiert; die uebrigen Eintraege
// werden trotzdem gelesen; skipped zaehlt die Dateien, die skip-pattern
// ausgenommen hat; names traegt die Basisnamen ohne .md ALLER gesehenen
// slice-*.md, vor beiden Abzuegen (die Zuordnung unter match: name). Ein
// unlesbares DoneDir selbst ergibt eine leere
// Menge und damit den Leerlauf-Befund.
func reviewCandidates(fsys driven.Filesystem, cfg model.ReviewsConfig) reviewWalk {
	w := reviewWalk{fsys: fsys, cfg: cfg}
	if cfg.SkipPattern != "" {
		w.skipRE = regexp.MustCompile(cfg.SkipPattern)
	}
	if entries, err := fsys.List(cfg.DoneDir); err == nil {
		w.visit(cfg.DoneDir, entries)
	}
	sort.Strings(w.out)
	sort.Strings(w.badDirs)
	return w
}

// reviewWalk sammelt die Kandidaten und die unlesbaren Unterverzeichnisse.
type reviewWalk struct {
	fsys    driven.Filesystem
	cfg     model.ReviewsConfig
	skipRE  *regexp.Regexp
	out     []string
	badDirs []string
	skipped int
	names   []string
}

func (w *reviewWalk) visit(dir string, entries []driven.DirEntry) {
	for _, e := range entries {
		rel := path.Join(dir, e.Name)
		if e.Kind == driven.KindDir {
			w.descend(rel, e.Name)
			continue
		}
		if e.Kind != driven.KindFile || !strings.HasSuffix(e.Name, ".md") || !strings.HasPrefix(e.Name, "slice-") {
			continue
		}
		w.names = append(w.names, strings.TrimSuffix(e.Name, ".md"))
		if matchAnyGlob(w.cfg.ExemptPaths, rel) {
			continue
		}
		if reviewSkipped(w.fsys, w.skipRE, rel) {
			w.skipped++
			continue
		}
		w.out = append(w.out, rel)
	}
}

func (w *reviewWalk) descend(rel, name string) {
	if !w.cfg.Recursive || isSkipDir(name) {
		return
	}
	entries, err := w.fsys.List(rel)
	if err != nil {
		w.badDirs = append(w.badDirs, rel)
		return
	}
	w.visit(rel, entries)
}

// reviewSkipped nimmt eine Datei nach ihrem Inhalt aus. Eine unlesbare Datei
// bleibt Kandidatin.
func reviewSkipped(fsys driven.Filesystem, skipRE *regexp.Regexp, rel string) bool {
	if skipRE == nil {
		return false
	}
	b, err := fsys.ReadFile(rel)
	return err == nil && skipRE.Match(b)
}

// hasMatchingReview prueft, ob mindestens ein Eintrag in reviewNames die
// slice-<NNN>-Kennung im Dateinamen traegt -- dieselbe Form wie
// tools/archive-wave/collect.go's CollectReviews (1:N zulaessig, z. B.
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

// hasReviewContaining prueft, ob ein Report-Dateiname den Basisnamen des Slice
// enthaelt (match: name) -- mit Datums-Praefix und beliebigem Suffix -- und
// keinen der laengeren Basisnamen, die ihn enthalten: der Report zu
// slice-a-foo-bar deckt slice-a-foo nicht, wenn slice-a-foo-bar ein Slice ist.
// GRENZE: ein Report, der beide Namen traegt, deckt nur den laengeren; vor dem
// Basisnamen gilt keine Grenze.
func hasReviewContaining(reviewNames []driven.DirEntry, key string, longer []string) bool {
	for _, e := range reviewNames {
		if e.Kind != driven.KindFile || !nameInReport(e.Name, key) {
			continue
		}
		covered := true
		for _, l := range longer {
			if nameInReport(e.Name, l) {
				covered = false
				break
			}
		}
		if covered {
			return true
		}
	}
	return false
}

// nameInReport prueft, ob report den Basisnamen key traegt: nach ihm steht ein
// Zeichen, das weder Buchstabe noch Ziffer ist (Unicode), oder das Ende --
// slice-x1 steht nicht in slice-x12-y.
func nameInReport(report, key string) bool {
	for rest := report; ; {
		i := strings.Index(rest, key)
		if i < 0 {
			return false
		}
		if r, _ := utf8.DecodeRuneInString(rest[i+len(key):]); r == utf8.RuneError ||
			!unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
		rest = rest[i+1:]
	}
}

// longerNames liefert die Basisnamen, die key enthalten und laenger sind --
// die Namen, deren Report key nicht mitdeckt.
func longerNames(names []string, key string) []string {
	var out []string
	for _, n := range names {
		if len(n) > len(key) && strings.Contains(n, key) {
			out = append(out, n)
		}
	}
	return out
}
