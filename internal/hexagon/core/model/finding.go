// Package model trägt das Befund-Modell und die Konfigurations-Typen des
// Kerns — der innerste, importfreie Daten-Ring (spec/architecture.md §Kern;
// ADR-0012).
package model

import "sort"

// Grund-Codes der Befunde (spec/spezifikation.md §4 — stabil,
// maschinenlesbar).
const (
	ReasonTargetMissing  = "target-missing"
	ReasonRepoEscape     = "repo-escape"
	ReasonSymlink        = "symlink"
	ReasonIDUnlinked     = "id-unlinked"
	ReasonSpanUnclosed   = "span-unclosed"
	// ReasonFenceUnclosed ist die Block-Ebene von ReasonSpanUnclosed: eine
	// Fence-Öffnung ohne Schluss bis zum DATEI-Ende (nicht Absatz-Ende — ein
	// Fence *ist* eine Absatzgrenze). Alles dahinter gilt jeder
	// Vorverarbeitung als Code und wird von allen Modulen übersprungen.
	ReasonFenceUnclosed  = "fence-unclosed"
	ReasonSpanNestedLink = "span-nested-link"
	ReasonHostpathForbidden = "hostpath-forbidden"
	ReasonDiagramIDUndefined = "diagram-id-undefined"
	ReasonVersionStale       = "version-stale"
	ReasonLinkStale          = "link-stale"
	ReasonCoreDrift          = "core-drift"
	ReasonCoreDriftVCS       = "core-drift-vcs"
	ReasonCommitUntraceable  = "commit-untraceable"
	ReasonPlanningDrift      = "planning-drift"
	ReasonTargetUntracked    = "target-untracked"
	// Die drei Closure-Note-Struktur-Codes des Moduls planning
	// (DC-FA-PLAN-001 §Closure-Note-Struktur). Drei statt einem, weil sie drei
	// verschiedene Reparaturen verlangen: Abschnitt schreiben, Substanz
	// ergänzen, Floskel ersetzen.
	ReasonClosureNoteMissing     = "closure-note-missing"
	ReasonClosureNoteThin        = "closure-note-thin"
	ReasonClosureNoteBoilerplate = "closure-note-boilerplate"
	// ReasonClosureNotePlaceholder ist die vierte, OPT-IN Struktur-Bedingung
	// (DC-FA-PLAN-001, Schalter planning.closure.placeholder): der unausgefuellte
	// Rumpf einer Vorlage. Er passiert die drei anderen Bedingungen, weil er
	// syntaktisch vollstaendig ist.
	ReasonClosureNotePlaceholder = "closure-note-placeholder"
	// ReasonClosureNoteAmbiguous: mehrere passende Ueberschriften — ohne
	// eindeutigen Abschnitt sagt eine Satzzahl nichts, und ein zweiter Abschnitt
	// ist der typische Rest einer Vorlage.
	ReasonClosureNoteAmbiguous = "closure-note-ambiguous"
	// Die dreizehn Grund-Codes des Moduls structure (DC-FA-STRUCT-001). Zwei
	// Struktur-Codes fuer die Abschnitts-Findung und elf Bedingungs-Codes —
	// je eigener Code, weil die Befund-Deduplikation zwei Verletzungen
	// desselben Abschnitts sonst zusammenfallen liesse.
	ReasonLinkPositionDependent = "link-position-dependent"
	ReasonWaveDrift             = "wave-drift"
	ReasonWavePreviewExists     = "wave-preview-exists"
	ReasonWaveResultsMissing    = "wave-results-missing"
	ReasonWaveUnregistered      = "wave-unregistered"
	ReasonSectionMissing        = "section-missing"
	ReasonSectionAmbiguous      = "section-ambiguous"
	ReasonSectionEmpty          = "section-empty"
	ReasonSectionThin           = "section-thin"
	ReasonSectionOversized      = "section-oversized"
	ReasonSectionForbidden      = "section-forbidden"
	ReasonSectionPatternMissing = "section-pattern-missing"
	ReasonSectionMarkerMissing  = "section-marker-missing"
	// ReasonSectionHeadingMismatch: eine Ueberschrift INNERHALB des Abschnitts
	// genuegt headings-match nicht; der Befund steht auf ihrer Zeile.
	ReasonSectionHeadingMismatch = "section-heading-mismatch"
	// Chronologie-Monotonie (ADR-0057): section-unordered meldet die
	// brechende Datenzeile (bzw. den Leerlauf ohne Datenzeile),
	// section-cell-untyped die untypisierbare Schluesselzelle — Befund statt
	// stillem Uebersprung, sonst schaltete ein Tippfehler die Pruefung der
	// restlichen Tabelle wortlos ab.
	ReasonSectionUnordered   = "section-unordered"
	ReasonSectionCellUntyped = "section-cell-untyped"
	// Zellenlaenge (ADR-0069), drei Codes fuer drei Reparaturen:
	// section-cell-oversized (kuerzen) und section-cell-undersized
	// (ausfuellen) melden auf der Zeile IHRER Zelle; eine Obergrenze allein
	// liesse die leere Zelle gruen passieren, weil 0 unter jeder Schwelle
	// liegt. section-column-missing meldet, dass die benannte Spalte nicht
	// adressierbar ist (kein Kopf traegt den Namen, er kommt mehrfach vor,
	// oder die Datenzeile reicht nicht bis zu ihr) — die Bedingung zu setzen
	// IST die Behauptung, dass es diese Spalte gibt.
	ReasonSectionCellOversized  = "section-cell-oversized"
	ReasonSectionCellUndersized = "section-cell-undersized"
	ReasonSectionColumnMissing  = "section-column-missing"
	// ReasonSectionTasksOpen meldet ein OFFENES Task-Item im Abschnitt. Eigener
	// Code neben section-oversized: jener zaehlt ALLE Items auf dem bereinigten
	// Text, dieser die offenen auf den ROHEN Zeilen (ADR-0074).
	ReasonSectionTasksOpen      = "section-tasks-open"
)

// Finding ist ein einzelner Befund (spec/spezifikation.md §2). json- und
// yaml-Tags halten dieselben Schlüssel für beide Serialisierungen
// (DC-FA-CLI-004); ohne yaml-Tag würde yaml.v3 die Feldnamen kleinschreiben.
type Finding struct {
	File    string `json:"file" yaml:"file"`
	Line    int    `json:"line" yaml:"line"`
	Rule    string `json:"rule" yaml:"rule"`
	Target  string `json:"target" yaml:"target"`
	Reason  string `json:"reason" yaml:"reason"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// SortFindings sortiert stabil nach (Datei, Zeile, Regel, Ziel, Grund)
// und entfernt identische Tupel — DC-QA-02 (spec/spezifikation.md
// §DC-QA-02.a).
func SortFindings(fs []Finding) []Finding {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		if a.Target != b.Target {
			return a.Target < b.Target
		}
		return a.Reason < b.Reason
	})
	out := fs[:0]
	var prev *Finding
	for i := range fs {
		f := fs[i]
		if prev != nil && f.File == prev.File && f.Line == prev.Line &&
			f.Rule == prev.Rule && f.Target == prev.Target && f.Reason == prev.Reason {
			continue
		}
		out = append(out, f)
		prev = &out[len(out)-1]
	}
	return out
}
