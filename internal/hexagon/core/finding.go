// Package core trägt die I/O-freie Prüf-Logik von d-check: Markdown-
// Analyse, Regelmodule, Befund-Modell und die Port-Interfaces
// (spec/architecture.md §2 — Kern; Modul-Pfade gemäß ADR-0004).
package core

import "sort"

// Grund-Codes der Befunde (spec/spezifikation.md §4 — stabil,
// maschinenlesbar).
const (
	ReasonTargetMissing  = "target-missing"
	ReasonRepoEscape     = "repo-escape"
	ReasonSymlink        = "symlink"
	ReasonIDUnlinked     = "id-unlinked"
	ReasonSpanUnclosed   = "span-unclosed"
	ReasonSpanNestedLink = "span-nested-link"
)

// Finding ist ein einzelner Befund (spec/spezifikation.md §2).
type Finding struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Rule    string `json:"rule"`
	Target  string `json:"target"`
	Reason  string `json:"reason"`
	Message string `json:"message,omitempty"`
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
