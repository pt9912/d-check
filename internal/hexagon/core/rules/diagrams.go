package rules

import (
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckDiagrams ist das Regelmodul `diagrams` (DC-FA-DIAG-001): es öffnet
// gezielt die in cfg.Fences gelisteten Diagramm-Fences (Default mermaid)
// und prüft jede dort gefundene Kennung auf Existenz in ihrer
// defined-in-Quelle. Anders als ids gilt KEINE Link-Policy (in Fences
// kein Markdown-Link möglich); geprüft wird, ob das Token in defined-in
// als eigenständiges Kennungs-Token außerhalb von Fences vorkommt. Reine
// Token-Extraktion über Rohtext (kein Mermaid-Parser, DC-QA-03); der
// Fence-Zustand der übrigen Module bleibt unberührt. defCache cached je
// (defined-in, Regex) die Token-Menge der Quelle.
func CheckDiagrams(fsys driven.Filesystem, file string, content []byte, cfg model.DiagramsConfig, defCache map[string]map[string]bool) []model.Finding {
	if len(cfg.Patterns) == 0 {
		return nil
	}
	langs := map[string]bool{}
	for _, l := range cfg.EffectiveFences() {
		langs[strings.ToLower(l)] = true
	}
	var findings []model.Finding
	for _, pl := range diagramFenceLines(content, langs) {
		findings = append(findings, diagramLineFindings(fsys, file, pl, cfg.Patterns, defCache)...)
	}
	return findings
}

// diagramLineFindings prüft eine Diagramm-Fence-Zeile gegen alle Muster
// (Deklarationsreihenfolge = Präzedenz; überlappende Vorkommen gehören
// dem früheren Muster, wie im Modul ids).
func diagramLineFindings(fsys driven.Filesystem, file string, pl proseLine, patterns []model.DiagramPattern, defCache map[string]map[string]bool) []model.Finding {
	var claimed [][2]int
	var findings []model.Finding
	for _, p := range patterns {
		for _, m := range p.Regex.FindAllStringIndex(pl.raw, -1) {
			if overlapsClaimed(claimed, m[0], m[1]) {
				continue
			}
			claimed = append(claimed, [2]int{m[0], m[1]})
			id := pl.raw[m[0]:m[1]]
			if definedTokens(fsys, p, defCache)[id] {
				continue // im defined-in definiert
			}
			findings = append(findings, model.Finding{
				File: file, Line: pl.no, Rule: "diagrams",
				Target: id, Reason: model.ReasonDiagramIDUndefined,
				Message: "Kennung im Diagramm ohne Definition in defined-in",
			})
		}
	}
	return findings
}

// definedTokens liefert (gecacht) die Menge der Kennungs-Token, die das
// Muster in seiner defined-in-Quelle außerhalb von Fences findet —
// Heading wie Fließtext/Tabelle (proseLines entfernt nur Fences). Ein
// Lesefehler liefert eine leere Menge (fail-closed: das Token gilt dann
// als undefiniert; die Existenz von defined-in wird beim Lauf-Start
// erzwungen, ensureDiagramsDefinedInExist).
func definedTokens(fsys driven.Filesystem, p model.DiagramPattern, cache map[string]map[string]bool) map[string]bool {
	key := p.DefinedIn + "\x00" + p.Regex.String()
	if set, ok := cache[key]; ok {
		return set
	}
	set := map[string]bool{}
	if content, err := fsys.ReadFile(p.DefinedIn); err == nil {
		var sb strings.Builder
		for _, pl := range proseLines(content) {
			sb.WriteString(pl.raw)
			sb.WriteByte('\n')
		}
		for _, tok := range p.Regex.FindAllString(sb.String(), -1) {
			set[tok] = true
		}
	}
	cache[key] = set
	return set
}
