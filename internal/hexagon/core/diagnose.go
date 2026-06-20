package core

import "path/filepath"

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
func FixCandidateFor(f Finding, cfg Config) *FixCandidate {
	if f.Reason != ReasonIDUnlinked {
		return nil
	}
	id := f.Target
	// Erstes Muster in Deklarationsreihenfolge, das die Kennung VOLL
	// matcht — dieselbe Präzedenz wie CheckIDLine, das den Befund erzeugt.
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
		ReasonTargetMissing, ReasonRepoEscape, ReasonSymlink,
		ReasonAnchorMissing, ReasonIDUnlinked, ReasonCodepathMissing,
		ReasonMatrixInactive, ReasonMatrixForbidden,
		ReasonExternalStatus, ReasonExternalTimeout, ReasonExternalRedirects,
		ReasonSpanUnclosed, ReasonSpanNestedLink, ReasonHostpathForbidden,
	}
}

// reasonTexts bildet jeden Grund-Code auf einen Klartext ab
// (spec/spezifikation.md §4). Als Funktion statt Paket-Global
// (gochecknoglobals); die Diagnose zeigt den Klartext statt des nackten
// Codes.
func reasonTexts() map[string]string {
	return map[string]string{
		ReasonTargetMissing:     "Linkziel existiert nicht",
		ReasonRepoEscape:        "Aufgelöstes Ziel verlässt die Repository-Wurzel",
		ReasonSymlink:           "Ziel ist oder enthält einen Symlink",
		ReasonAnchorMissing:     "Anker entspricht keinem Heading-Slug",
		ReasonIDUnlinked:        "Kennung im Fließtext ohne Markdown-Link auf ihre Definition",
		ReasonCodepathMissing:   "Ziel eines Inline-Code-Pfads existiert nicht",
		ReasonMatrixInactive:    "Referenz auf ein Dokument mit verbotenem Status (z. B. superseded)",
		ReasonMatrixForbidden:   "Referenz zwischen Dokumentklassen nicht erlaubt (Referenzrichtung)",
		ReasonExternalStatus:    "Externer Link: HTTP-Status ≥ 400 oder Transportfehler",
		ReasonExternalTimeout:   "Externer Link: Zeitüberschreitung",
		ReasonExternalRedirects: "Externer Link: zu viele Redirects",
		ReasonSpanUnclosed:      "Ungeschlossene Code-Span-Öffnung (klebt an Nicht-Whitespace)",
		ReasonSpanNestedLink:    "Verschachtelte Link-Syntax im Linktext (rendert zerrissen)",
		ReasonHostpathForbidden: "Host-lokaler absoluter Pfad (Maschinen-Layout-Leak)",
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
