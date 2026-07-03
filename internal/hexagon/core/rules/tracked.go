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
// und die werden je einzeln geprüft). Die Index-Menge `tracked` lädt der
// Lauf einmal über den VCS-Port (Schritt 2); `exempt-targets` nimmt
// aufgelöste Ziel-Pfade referenz-weit aus (Schritt 4).
func CheckTracked(fsys driven.Filesystem, file string, lines []Line, cfg model.TrackedConfig, tracked map[string]bool) []model.Finding {
	var findings []model.Finding
	for _, ref := range ExtractLinks(lines) {
		rel, escaped, ok := localTarget(file, ref.Target)
		if !ok || escaped {
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
			Target: ref.Target, Reason: model.ReasonTargetUntracked,
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
