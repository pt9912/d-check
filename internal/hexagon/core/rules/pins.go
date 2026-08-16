package rules

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// dpinRE erkennt den Content-Pin-Marker `<!-- dpin: sha256:<hex> -->` und
// liefert den hinterlegten Hash in Gruppe 1 (spec/spezifikation.md
// §DC-FA-PIN-001.a Schritt 1).
var dpinRE = regexp.MustCompile(`<!--\s*dpin:\s*sha256:([0-9a-fA-F]+)\s*-->`)

// pinWhitespaceRE kollabiert Whitespace-Folgen für die reflow-invariante
// Normalisierung des Ziel-Spans (§DC-FA-PIN-001.a Schritt 3).
var pinWhitespaceRE = regexp.MustCompile(`\s+`)

// CheckPins ist das Regelmodul `pins` (DC-FA-PIN-001): ein Link mit Content-Pin
// (`<!-- dpin: sha256:<hex> -->`, unmittelbar nach dem Link) wird gegen den
// whitespace-normalisierten **rohen** Ziel-Span (ganze Datei oder
// Heading-Section, inkl. Fenced-Code) gehasst; Mismatch → `link-stale`. Nur
// auflösbare, repo-interne Ziele — der strukturelle Befund bleibt
// `links`/`anchors` (kein eigener/doppelter Befund). Diagnose-only. spanCache
// cached den errechneten Hash je aufgelöstem (Datei, Anker)-Span.
func CheckPins(fsys driven.Filesystem, file string, lines []Line, content []byte, spanCache map[string]string) []model.Finding {
	rawLines := strings.Split(string(content), "\n")
	var findings []model.Finding
	for _, ln := range lines {
		raw := ""
		if i := ln.No - 1; i >= 0 && i < len(rawLines) {
			raw = rawLines[i]
		}
		for _, b := range bindPins(ln.Text, raw) {
			if f, ok := pinFinding(fsys, file, ln.No, b.target, b.hash, spanCache); ok {
				findings = append(findings, f)
			}
		}
	}
	return findings
}

type pinBinding struct {
	target string
	hash   string
}

// bindPins findet die dpin-Marker einer Zeile (Links + Marker auf der
// vorverarbeiteten Zeile — so zählen nur echte, nicht in Inline-Code/Fences
// liegende Vorkommen) und bindet jeden Marker an den Nicht-Bild-Link, dessen
// schließendes `)` ihm **unmittelbar** vorausgeht. Die „nur Whitespace
// dazwischen"-Prüfung läuft auf der **rohen** Zeile (`raw`): die Vorverarbeitung
// leert Inline-Code positionserhaltend zu Leerzeichen, was den Vertrag sonst
// verfälschte (Impl-R1 F-1). Positionen stimmen, weil das Stripping längen-
// erhaltend ist. Ein Marker ohne unmittelbar vorausgehenden Link ist inert
// (§DC-FA-PIN-001.a Schritt 1).
func bindPins(text, raw string) []pinBinding {
	markers := dpinRE.FindAllStringSubmatchIndex(text, -1)
	if markers == nil {
		return nil
	}
	type lk struct {
		end    int
		target string
	}
	var links []lk
	forEachLink(text, func(ref LinkRef, span LinkSpan) {
		if !span.IsImage {
			links = append(links, lk{end: span.End, target: ref.Target})
		}
	})
	var out []pinBinding
	for _, m := range markers {
		start := m[0]
		best := -1
		for i, l := range links {
			if l.end <= start && start <= len(raw) && strings.TrimSpace(raw[l.end:start]) == "" {
				if best == -1 || l.end > links[best].end {
					best = i
				}
			}
		}
		if best != -1 {
			out = append(out, pinBinding{target: links[best].target, hash: strings.ToLower(text[m[2]:m[3]])})
		}
	}
	return out
}

// pinFinding löst das Ziel auf, hasst den normalisierten Span und vergleicht.
// ok=false (kein Befund), wenn das Ziel extern, nicht auflösbar oder außerhalb
// der Repo-Wurzel ist (struktureller Befund bleibt links/anchors) oder der Hash
// stimmt.
func pinFinding(fsys driven.Filesystem, file string, line int, target, want string, cache map[string]string) (model.Finding, bool) {
	if IsExternalScheme(target) {
		return model.Finding{}, false
	}
	pathPart, anchor := target, ""
	if i := strings.IndexByte(target, '#'); i != -1 {
		pathPart, anchor = target[:i], DecodeFragment(target[i+1:])
	}
	rel := file // pathPart == "" → Same-file-Anker
	if pathPart != "" {
		r, escaped, ok := ResolveTarget(file, pathPart)
		if !ok || escaped {
			return model.Finding{}, false
		}
		rel = r
	}
	got, ok := spanHash(fsys, rel, anchor, cache)
	if !ok || got == want {
		return model.Finding{}, false
	}
	return model.Finding{
		File: file, Line: line, Rule: "pins",
		Target: target, Reason: model.ReasonLinkStale,
		Message: "Ziel-Span gedriftet: erwartet sha256:" + shortHash(want) + ", errechnet sha256:" + got + " (voller Ist-Hash zum Re-Pinnen)",
	}, true
}

// spanHash liefert (gecacht) den whitespace-normalisierten SHA-256 des
// Ziel-Spans: ganze Datei ohne Anker, sonst die Heading-Section. ok=false, wenn
// die Datei fehlt oder der Anker keine Section trifft (§DC-FA-PIN-001.a
// Schritt 2–3).
func spanHash(fsys driven.Filesystem, rel, anchor string, cache map[string]string) (string, bool) {
	key := rel + "\x00" + anchor
	if h, hit := cache[key]; hit {
		return h, h != ""
	}
	h := ""
	if content, err := fsys.ReadFile(rel); err == nil {
		span := string(content)
		ok := true
		if anchor != "" {
			span, ok = headingSection(content, anchor)
		}
		if ok {
			norm := strings.TrimSpace(pinWhitespaceRE.ReplaceAllString(span, " "))
			sum := sha256.Sum256([]byte(norm))
			h = hex.EncodeToString(sum[:])
		}
	}
	cache[key] = h
	return h, h != ""
}

// shortHash kürzt einen Hash für die Befund-Meldung.
func shortHash(h string) string {
	if len(h) > 12 {
		return h[:12] + "…"
	}
	return h
}
