package rules

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// immutableRE erkennt den Immutabilitäts-Pin `<!-- immutable: sha256:<hex> -->`
// und liefert den hinterlegten Hash in Gruppe 1 (spec/spezifikation.md
// §DC-FA-IMM-001.a Schritt 1).
var immutableRE = regexp.MustCompile(`<!--\s*immutable:\s*sha256:([0-9a-fA-F]+)\s*-->`)

// CheckImmutable ist das Regelmodul `immutable` (DC-FA-IMM-001): trägt eine
// Datei den Pin-Marker (erster Treffer auf einer **vorverarbeiteten** Zeile —
// in Fenced-/Inline-Code ist er inert), wird ihr whitespace-normalisierter
// **Core** (rohe Datei ohne die Marker-Zeile und ohne die exclude-sections-
// Abschnitte) gegen den Pin gehasht; Abweichung → `core-drift`. Hermetisch
// (kein git — nur die gescannte Datei selbst wird gelesen), diagnose-only;
// ohne Marker keine Prüfung (opt-in pro Datei). Default-aus byte-identisch.
func CheckImmutable(file string, lines []Line, content []byte, cfg model.ImmutableConfig) []model.Finding {
	markerLine, want, ok := firstImmutablePin(lines)
	if !ok {
		return nil
	}
	got := immutableCoreHash(content, markerLine, cfg.ExcludeSections)
	if got == want {
		return nil
	}
	return []model.Finding{{
		File: file, Line: markerLine, Rule: "immutable",
		Target: "sha256:" + shortHash(want), Reason: model.ReasonCoreDrift,
		Message: "Core gedriftet: erwartet sha256:" + shortHash(want) +
			", errechnet sha256:" + shortHash(got),
	}}
}

// firstImmutablePin liefert Zeilennummer und (kleingeschriebenen) Hash des
// **ersten** Pin-Markers auf einer vorverarbeiteten Zeile; ok=false ohne
// Marker. Weitere Marker sind inert (ein wirksamer Pin je Datei).
func firstImmutablePin(lines []Line) (int, string, bool) {
	for _, ln := range lines {
		if m := immutableRE.FindStringSubmatch(ln.Text); m != nil {
			return ln.No, strings.ToLower(m[1]), true
		}
	}
	return 0, "", false
}

// immutableCoreHash bildet den SHA-256 des normalisierten Core: die **rohe**
// Datei ohne die Marker-Zeile (sonst Selbstbezug) und ohne die
// exclude-sections-Abschnitte (Section-Abgrenzung wie matrix.exclude-sections),
// alle Whitespace-Folgen zu einem Leerzeichen kollabiert, Rand getrimmt —
// dieselbe reflow-invariante Normalisierung wie pins (spec/spezifikation.md
// §DC-FA-IMM-001.a Schritt 2–3).
func immutableCoreHash(content []byte, markerLine int, excludeSections []string) string {
	excluded := excludedRanges(content, excludeSections)
	var b strings.Builder
	for i, raw := range strings.Split(string(content), "\n") {
		no := i + 1
		if no == markerLine || inRanges(excluded, no) {
			continue
		}
		b.WriteString(raw)
		b.WriteByte('\n')
	}
	norm := strings.TrimSpace(pinWhitespaceRE.ReplaceAllString(b.String(), " "))
	sum := sha256.Sum256([]byte(norm))
	return hex.EncodeToString(sum[:])
}
