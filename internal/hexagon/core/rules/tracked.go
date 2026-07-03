package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckTracked ist das Regelmodul `tracked` (DC-FA-TRK-001): jedes
// auflösbare, **existierende** repo-interne Link-/Bild-Ziel (Datei-Ebene
// der links-Auflösung, Fragment verworfen) muss im git-Index getrackt
// sein — sonst wäre die Referenz auf jedem frischen Klon ein
// `target-missing` (Umgebungs-Drift zwischen Arbeitsbäumen).
//
// Kein Doppelbefund (DC-FA-TRK-001.a Schritt 3): nicht existierende oder
// escapte Ziele bleiben Sache von `links`; Verzeichnis-Ziele sind kein
// Kandidat (der git-Index führt nur Dateien — ein Verzeichnis existiert
// auf dem frischen Klon genau dann, wenn es getrackte Dateien enthält,
// und die werden je einzeln geprüft); Symlink-Referenzen (Ziel ist oder
// durchläuft einen Symlink) sind die Domäne von DC-FA-LINK-002 — `links`
// meldet sie kategorisch als `symlink`, `tracked` prüft sie nicht (sonst
// false-positive hinter getrackten Verzeichnis-Symlinks: der Index führt
// den realen Pfad, nicht den Symlink-Alias). Die Index-Menge `tracked`
// lädt der Lauf einmal über den VCS-Port (Schritt 2); `exempt-targets`
// nimmt aufgelöste Ziel-Pfade referenz-weit aus (Schritt 4). Der Befund
// nennt als target den AUFGELÖSTEN Zielpfad (Schritt 5) — dieselbe Form,
// die das Ventil matcht.
func CheckTracked(fsys driven.Filesystem, file string, lines []Line, cfg model.TrackedConfig, tracked map[string]bool) []model.Finding {
	var findings []model.Finding
	for _, ref := range ExtractLinks(lines) {
		rel, escaped, ok := localTarget(file, ref.Target)
		if !ok || escaped {
			continue
		}
		if hit, serr := symlinkInPath(fsys, rel); serr != nil || hit {
			continue
		}
		kind, err := fsys.Kind(rel)
		if err != nil || kind != driven.KindFile {
			continue
		}
		if exemptTarget(rel, cfg.ExemptTargets) {
			continue
		}
		if tracked[rel] {
			continue
		}
		findings = append(findings, model.Finding{
			File: file, Line: ref.Line, Rule: "tracked",
			Target: rel, Reason: model.ReasonTargetUntracked,
			Message: "Linkziel ist nicht im git-Index getrackt (untracked/gitignoriert — fehlt auf jedem frischen Klon)",
		})
	}
	return findings
}

// exemptTarget prüft den aufgelösten Ziel-Pfad gegen die
// exempt-targets-Globs (referenz-weit, matchGlob wie scan.ignore).
func exemptTarget(rel string, globs []string) bool {
	for _, g := range globs {
		if matchGlob(g, rel) {
			return true
		}
	}
	return false
}
