// Package app trägt die Anwendungs-Modi auf den Befunden: Diagnose
// (--doctor), Reparatur-Patch (--repair) und Konfigurations-Vorschlag
// (--suggest-config). Importiert model und rules (spec/architecture.md
// §Kern; ADR-0012).
package app

import (
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"path/filepath"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// FixCandidate beschreibt eine vorgeschlagene, NICHT angewendete Änderung
// für einen Befund (spec/spezifikation.md §DC-FA-CLI-007.a). Sie ist die
// gemeinsame Quelle für den Diagnose-Modus (slice-025, --doctor) und den
// späteren Patch-Modus (slice-026, --repair): eine Ableitung, zwei
// Ausgaben.
type FixCandidate struct {
	Original    string // exakter Token, der ersetzt würde
	Replacement string // vorgeschlagener Ersatz (gerendert, nicht angewendet)
	Note        string // kurze Erläuterung für den Menschen
}

// FixCandidateFor leitet — wo aus dem Befund EINDEUTIG ableitbar — einen
// Fix-Kandidaten ab. In dieser Version liefert nur `id-unlinked` einen
// Kandidaten: die nackte Kennung wird zu einem Markdown-Link auf das in
// der passenden ids-Regel deklarierte Definitions-Target ausgeführt
// (Datei-Ebene; den genauen Anker setzt der Mensch). Best-Guess-Fälle
// (`target-missing`, `span-*` …) liefern hier bewusst KEINEN Kandidaten
// — sie gehören in die breite Stufe von --repair (slice-026). Gibt nil
// zurück, wenn kein eindeutiger Kandidat existiert.
func FixCandidateFor(f model.Finding, cfg model.Config) *FixCandidate {
	if f.Reason != model.ReasonIDUnlinked {
		return nil
	}
	id := f.Target
	// Erstes Muster in Deklarationsreihenfolge, das die Kennung VOLL
	// matcht — dieselbe Präzedenz wie rules.CheckIDLine, das den Befund erzeugt.
	for _, p := range cfg.IDPatterns {
		if loc := p.Regex.FindStringIndex(id); loc != nil && loc[0] == 0 && loc[1] == len(id) {
			return &FixCandidate{
				Original:    id,
				Replacement: "[`" + id + "`](" + relLink(f.File, p.Target) + ")",
				Note:        "Kennung als Markdown-Link auf ihre Definition (" + p.Target + ") ausführen; Anker ggf. ergänzen",
			}
		}
	}
	return nil
}

// relLink baut den repo-relativen Link vom Verzeichnis der Befund-Datei
// zum (repo-relativen) Target. Slash-normalisiert und deterministisch.
func relLink(fromFile, target string) string {
	rel, err := filepath.Rel(filepath.Dir(fromFile), target)
	if err != nil {
		return target
	}
	return filepath.ToSlash(rel)
}

// AllReasons ist die kanonische Liste aller Grund-Codes (über die
// Reason*-Konstanten, nicht über String-Literale) — die Quelle der
// Wahrheit für die Vollständigkeits-Prüfung des Klartext-Mappings
// (diagnose_test.go). Als Funktion statt Paket-Global (gochecknoglobals).
// Ein neuer Grund-Code wird hier ergänzt; fehlt der Klartext, bricht der
// Test, nicht still die Diagnose.
func AllReasons() []string {
	return []string{
		model.ReasonTargetMissing, model.ReasonRepoEscape, model.ReasonSymlink,
		rules.ReasonAnchorMissing, model.ReasonIDUnlinked, rules.ReasonCodepathMissing,
		rules.ReasonMatrixInactive, rules.ReasonMatrixForbidden, rules.ReasonMatrixDownward,
		rules.ReasonExternalStatus, rules.ReasonExternalTimeout, rules.ReasonExternalRedirects,
		model.ReasonSpanUnclosed, model.ReasonSpanNestedLink, model.ReasonHostpathForbidden,
	}
}

// reasonTexts bildet jeden Grund-Code auf einen Klartext ab
// (spec/spezifikation.md §4). Als Funktion statt Paket-Global
// (gochecknoglobals); die Diagnose zeigt den Klartext statt des nackten
// Codes.
func reasonTexts() map[string]string {
	return map[string]string{
		model.ReasonTargetMissing:     "Linkziel existiert nicht",
		model.ReasonRepoEscape:        "Aufgelöstes Ziel verlässt die Repository-Wurzel",
		model.ReasonSymlink:           "Ziel ist oder enthält einen Symlink",
		rules.ReasonAnchorMissing:     "Anker entspricht keinem Heading-Slug",
		model.ReasonIDUnlinked:        "Kennung im Fließtext ohne Markdown-Link auf ihre Definition",
		rules.ReasonCodepathMissing:   "Ziel eines Inline-Code-Pfads existiert nicht",
		rules.ReasonMatrixInactive:    "Referenz auf ein Dokument mit verbotenem Status (z. B. superseded)",
		rules.ReasonMatrixForbidden:   "Referenz zwischen Dokumentklassen nicht erlaubt (Referenzrichtung)",
		rules.ReasonMatrixDownward:    "Klasseninterner Abwärtsverweis gegen die deklarierte Rangordnung (order/direction)",
		rules.ReasonExternalStatus:    "Externer Link: HTTP-Status ≥ 400 oder Transportfehler",
		rules.ReasonExternalTimeout:   "Externer Link: Zeitüberschreitung",
		rules.ReasonExternalRedirects: "Externer Link: zu viele Redirects",
		model.ReasonSpanUnclosed:      "Ungeschlossene Code-Span-Öffnung (klebt an Nicht-Whitespace)",
		model.ReasonSpanNestedLink:    "Verschachtelte Link-Syntax im Linktext (rendert zerrissen)",
		model.ReasonHostpathForbidden: "Host-lokaler absoluter Pfad (Maschinen-Layout-Leak)",
	}
}

// ReasonText liefert den Klartext eines Grund-Codes; bei unbekanntem Code
// den Code selbst (fail-safe — die Vollständigkeits-Prüfung verhindert,
// dass das in der Praxis auftritt).
func ReasonText(reason string) string {
	if t, ok := reasonTexts()[reason]; ok {
		return t
	}
	return reason
}
