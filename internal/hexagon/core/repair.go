package core

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// RepairEdit ist eine vorgeschlagene, NICHT angewendete Zeilen-Ersetzung
// (slice-026, DC-FA-CLI-008). ReviewRequired markiert Best-Guess-Edits der
// breiten Stufe.
type RepairEdit struct {
	File           string
	Line           int
	OldLine        string
	NewLine        string
	ReviewRequired bool
}

// posRepl ist eine Ersetzung [start,end) → with innerhalb einer Zeile.
type posRepl struct {
	start, end int
	with       string
}

// RepairEdits leitet die Zeilen-Ersetzungen für --repair ab
// (spec/spezifikation.md §DC-FA-CLI-008.a). Konservativ (Default): nur
// eindeutige Fixes — `id-unlinked` → Markdown-Link auf die bekannte
// Definition ([`FixCandidateFor`]), angewandt nur auf nackte Prosa-
// Vorkommen (Inline-Code-Vorkommen sind im vorverarbeiteten Text geleert
// und werden nicht angefasst — kein zerrissener Code-Span). Mit broad
// zusätzlich Best-Guess: `target-missing` → eindeutige Datei gleichen
// Basisnamens, als ReviewRequired markiert. Read-only (liest betroffene
// Dateien über fsys); deterministisch (Eingabe sortiert, Ausgabe nach
// Datei/Zeile).
func RepairEdits(fsys driven.Filesystem, findings []Finding, cfg Config, broad bool) ([]RepairEdit, error) {
	var basenameIdx map[string]string
	if broad {
		idx, err := uniqueBasenames(fsys, cfg)
		if err != nil {
			return nil, err
		}
		basenameIdx = idx
	}
	var files []string
	perFile := map[string][]Finding{}
	for _, f := range findings {
		if _, ok := perFile[f.File]; !ok {
			files = append(files, f.File)
		}
		perFile[f.File] = append(perFile[f.File], f)
	}
	var out []RepairEdit
	for _, file := range files {
		edits, err := repairFile(fsys, file, perFile[file], cfg, broad, basenameIdx)
		if err != nil {
			return nil, err
		}
		out = append(out, edits...)
	}
	return out, nil
}

// repairFile liest die Datei einmal und erzeugt die Zeilen-Edits ihrer
// Befunde (Zeilen in aufsteigender Reihenfolge).
func repairFile(fsys driven.Filesystem, file string, findings []Finding, cfg Config, broad bool, basenameIdx map[string]string) ([]RepairEdit, error) {
	content, err := fsys.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%s nicht lesbar: %w", file, err)
	}
	rawLines := strings.Split(string(content), "\n")
	textByLine := map[int]string{}
	for _, ln := range PreprocessMarkdown(content) {
		textByLine[ln.No] = ln.Text
	}
	var lines []int
	perLine := map[int][]Finding{}
	for _, f := range findings {
		if _, ok := perLine[f.Line]; !ok {
			lines = append(lines, f.Line)
		}
		perLine[f.Line] = append(perLine[f.Line], f)
	}
	sort.Ints(lines)
	var out []RepairEdit
	for _, line := range lines {
		if line < 1 || line > len(rawLines) {
			continue
		}
		raw := rawLines[line-1]
		text, ok := textByLine[line]
		if !ok {
			text = raw
		}
		newRaw, review := repairLine(raw, text, file, perLine[line], cfg, broad, basenameIdx)
		if newRaw != raw {
			out = append(out, RepairEdit{File: file, Line: line, OldLine: raw, NewLine: newRaw, ReviewRequired: review})
		}
	}
	return out, nil
}

// repairLine wendet die anwendbaren Ersetzungen einer Zeile an: konservativ
// die id-unlinked-Kandidaten (alle nackten Prosa-Vorkommen je Token), mit
// broad zusätzlich die target-missing-Best-Guesses. Liefert die neue Zeile
// und ob ein Best-Guess (ReviewRequired) beteiligt war.
func repairLine(raw, text, file string, findings []Finding, cfg Config, broad bool, basenameIdx map[string]string) (string, bool) {
	spans := ExtractLinkSpans(text)
	repls := conservativeRepls(text, spans, findings, cfg)
	review := false
	if broad {
		if br := broadRepls(text, spans, file, findings, basenameIdx); len(br) > 0 {
			repls = append(repls, br...)
			review = true
		}
	}
	if len(repls) == 0 {
		return raw, false
	}
	return applyRepls(raw, repls), review
}

// conservativeRepls liefert die eindeutigen Fixes einer Zeile: je
// distinktem id-unlinked-Token alle nackten Prosa-Vorkommen → Link.
func conservativeRepls(text string, spans []LinkSpan, findings []Finding, cfg Config) []posRepl {
	var repls []posRepl
	seen := map[string]bool{}
	for _, f := range findings {
		if f.Reason != ReasonIDUnlinked || seen[f.Target] {
			continue
		}
		seen[f.Target] = true
		c := FixCandidateFor(f, cfg)
		if c == nil {
			continue
		}
		for _, occ := range bareOccurrences(text, spans, f.Target) {
			repls = append(repls, posRepl{occ[0], occ[1], c.Replacement})
		}
	}
	return repls
}

// broadRepls liefert die Best-Guess-Fixes einer Zeile: target-missing →
// eindeutige Datei gleichen Basisnamens (relativ zur Befund-Datei).
func broadRepls(text string, spans []LinkSpan, file string, findings []Finding, basenameIdx map[string]string) []posRepl {
	var repls []posRepl
	for _, f := range findings {
		if f.Reason != ReasonTargetMissing {
			continue
		}
		uniq, ok := basenameIdx[path.Base(f.Target)]
		if !ok {
			continue
		}
		if occ, found := destSpan(text, spans, f.Target); found {
			repls = append(repls, posRepl{occ[0], occ[1], relLink(file, uniq)})
		}
	}
	return repls
}

// bareOccurrences liefert die Byte-Spannen aller Vorkommen von token im
// (vorverarbeiteten) text, die NICHT in einem Link liegen — Positionen
// sind offset-gleich zur Rohzeile (Inline-Code ist positionserhaltend
// geleert).
func bareOccurrences(text string, spans []LinkSpan, token string) [][2]int {
	var out [][2]int
	for off := 0; off <= len(text)-len(token); {
		i := strings.Index(text[off:], token)
		if i < 0 {
			break
		}
		start := off + i
		end := start + len(token)
		// Wortgrenze: das Vorkommen darf nicht Teil einer längeren Kennung
		// sein (z. B. ADR-0001 in ADR-00012) — sonst landet die Ersetzung
		// falsch; konservativ heißt eindeutig.
		if wholeToken(text, start, end) && !idOccurrenceExempt(spans, start, end) {
			out = append(out, [2]int{start, end})
		}
		off = end
	}
	return out
}

// wholeToken prüft, dass [start,end) nicht von alphanumerischen Zeichen
// umgeben ist (keine längere Kennung, in die das Token eingebettet ist).
func wholeToken(s string, start, end int) bool {
	if start > 0 && isWordByte(s[start-1]) {
		return false
	}
	if end < len(s) && isWordByte(s[end]) {
		return false
	}
	return true
}

func isWordByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

// destSpan findet die Ziel-Spanne [start,end) des Links, dessen Ziel auf
// target normalisiert — für die Best-Guess-Ersetzung von target-missing.
func destSpan(text string, spans []LinkSpan, target string) ([2]int, bool) {
	for _, sp := range spans {
		start := sp.TextEnd + 2 // hinter "]("
		end := sp.End - 1       // das ")"
		if start > end || end > len(text) {
			continue
		}
		if normalizeTarget(text[start:end]) == target {
			return [2]int{start, end}, true
		}
	}
	return [2]int{}, false
}

// applyRepls wendet die Ersetzungen rechts-nach-links an (stabile
// Offsets); überlappende Ersetzungen werden verworfen.
func applyRepls(raw string, repls []posRepl) string {
	sort.Slice(repls, func(i, j int) bool { return repls[i].start > repls[j].start })
	lastStart := len(raw)
	for _, r := range repls {
		if r.start < 0 || r.start >= r.end || r.end > len(raw) || r.end > lastStart {
			continue
		}
		raw = raw[:r.start] + r.with + raw[r.end:]
		lastStart = r.start
	}
	return raw
}

// uniqueBasenames bildet jeden Datei-Basisnamen, der im Scan-Bestand
// genau einmal vorkommt, auf seinen (repo-relativen) Pfad ab — Grundlage
// des target-missing-Best-Guess.
func uniqueBasenames(fsys driven.Filesystem, cfg Config) (map[string]string, error) {
	files, err := DiscoverFiles(fsys, cfg.Roots, cfg.Ignore)
	if err != nil {
		return nil, err
	}
	byBase := map[string][]string{}
	for _, f := range files {
		b := path.Base(f)
		byBase[b] = append(byBase[b], f)
	}
	out := map[string]string{}
	for b, ps := range byBase {
		if len(ps) == 1 {
			out[b] = ps[0]
		}
	}
	return out, nil
}
