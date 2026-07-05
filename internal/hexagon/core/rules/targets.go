package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls targets (DC-FA-TGT-001; spec/spezifikation.md §4).
const (
	ReasonGatePhantom      = "gate-phantom"
	ReasonGateUndocumented = "gate-undocumented"
)

var (
	// makefileRuleRe erkennt eine Regelzeile in Parität zum abgelösten
	// tools/gate-consistency.sh (`^[a-zA-Z][a-zA-Z0-9 _-]*:([^=]|$)`):
	// Regelname(n) am Zeilenanfang, gefolgt von `:` ohne unmittelbar folgendes
	// `=` (Zuweisungen `X :=`/`X ?=` ausgenommen); der erste Name ist ein
	// Buchstabe (`.PHONY`/`.DEFAULT_GOAL` und Pattern-Rules `%…` fallen weg).
	makefileRuleRe = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9 _-]*):([^=]|$)`)
	// docTargetRe erkennt ein `make X`-Token (X = [a-z][a-z0-9_-]*, in Parität
	// zum Skript-`grep -oE '`+"`"+`make [a-z][a-z0-9_-]*`+"`"+`'`).
	docTargetRe = regexp.MustCompile("`make ([a-z][a-z0-9_-]*)`")
)

// targetRef ist ein erkanntes Target mit Fundstelle.
type targetRef struct {
	name string
	file string
	line int
}

// CheckTargets ist das Regelmodul targets (DC-FA-TGT-001): es prüft **hermetisch**
// (nur der Filesystem-Port, **kein** git, **kein** Netz, **kein Ausführen** des
// Makefile) die Deklarations-Konsistenz Doku ↔ Build-Targets. Zwei Richtungen:
// ein in einer Doku-**Tabellenzeile** als `make X` behauptetes Target ohne
// Makefile-Regel ⇒ gate-phantom (Richtung 1); eine Makefile-Regel (minus
// exempt-targets) ohne Eintrag in der Autoritäts-Doku ⇒ gate-undocumented
// (Richtung 2). **fail-closed:** eine fehlende/unlesbare konfigurierte Datei ⇒
// error (Exit 2). Leeres Makefiles ⇒ inert; die Richtungen sind an ihre
// jeweilige Doku-Quelle gekoppelt und voneinander unabhängig. Diagnose-only.
func CheckTargets(fsys driven.Filesystem, cfg model.TargetsConfig) ([]model.Finding, error) {
	if len(cfg.Makefiles) == 0 || fsys == nil {
		return nil, nil // inert (keine Regelmenge / kein Port)
	}
	ruleRefs, ruleSet, err := collectMakefileRules(fsys, cfg.Makefiles)
	if err != nil {
		return nil, err
	}
	phantom, err := phantomFindings(fsys, cfg.DocTables, ruleSet)
	if err != nil {
		return nil, err
	}
	undoc, err := undocumentedFindings(fsys, cfg.Authority, cfg.ExemptTargets, ruleRefs)
	if err != nil {
		return nil, err
	}
	return append(phantom, undoc...), nil
}

// collectMakefileRules extrahiert die Regelnamen (mit Fundstelle) aus allen
// konfigurierten Makefile-Quellen. fail-closed bei fehlender/unlesbarer Datei.
func collectMakefileRules(fsys driven.Filesystem, makefiles []string) ([]targetRef, map[string]bool, error) {
	var refs []targetRef
	set := map[string]bool{}
	for _, mf := range makefiles {
		content, err := fsys.ReadFile(mf)
		if err != nil {
			return nil, nil, fmt.Errorf("das Modul targets kann das Makefile %q nicht lesen (DC-FA-TGT-001, fail-closed): %w", mf, err)
		}
		for i, line := range splitLines(content) {
			m := makefileRuleRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			for _, name := range strings.Fields(m[1]) {
				refs = append(refs, targetRef{name: name, file: mf, line: i + 1})
				set[name] = true
			}
		}
	}
	return refs, set, nil
}

// phantomFindings ist Richtung 1: jedes in einer doc-tables-Tabellenzeile
// dokumentierte `make X` ohne Makefile-Regel ⇒ gate-phantom.
func phantomFindings(fsys driven.Filesystem, docTables []string, ruleSet map[string]bool) ([]model.Finding, error) {
	var out []model.Finding
	for _, dt := range docTables {
		docs, err := extractDocTargets(fsys, dt)
		if err != nil {
			return nil, err
		}
		for _, d := range docs {
			if !ruleSet[d.name] {
				out = append(out, model.Finding{
					File: d.file, Line: d.line, Rule: "targets", Target: d.name,
					Reason:  ReasonGatePhantom,
					Message: "dokumentiertes Target `make " + d.name + "` ohne Makefile-Regel",
				})
			}
		}
	}
	return out, nil
}

// undocumentedFindings ist Richtung 2: jede nicht-exempte Makefile-Regel ohne
// Eintrag in der Autoritäts-Doku ⇒ gate-undocumented. Leere authority ⇒ Richtung
// 2 entfällt (kein Befund).
func undocumentedFindings(fsys driven.Filesystem, authority string, exemptTargets []string, ruleRefs []targetRef) ([]model.Finding, error) {
	if authority == "" {
		return nil, nil
	}
	authDocs, err := extractDocTargets(fsys, authority)
	if err != nil {
		return nil, err
	}
	authSet := map[string]bool{}
	for _, d := range authDocs {
		authSet[d.name] = true
	}
	exempt := map[string]bool{}
	for _, e := range exemptTargets {
		exempt[e] = true
	}
	var out []model.Finding
	for _, r := range ruleRefs {
		if exempt[r.name] || authSet[r.name] {
			continue
		}
		out = append(out, model.Finding{
			File: r.file, Line: r.line, Rule: "targets", Target: r.name,
			Reason:  ReasonGateUndocumented,
			Message: "Makefile-Regel `" + r.name + "` ohne Deklaration in der Autoritäts-Doku " + authority,
		})
	}
	return out, nil
}

// extractDocTargets liest eine Doku-Datei und extrahiert die `make X`-Tokens
// **ausschließlich aus Tabellenzeilen** — eine Tabellenzeile ist eine Zeile,
// deren **erstes Zeichen** ein Pipe `|` ist (Spalte 0, in Parität zu
// `grep -E '^\|'`; Einrückung zählt nicht). Prosa-Erwähnungen zählen nicht.
func extractDocTargets(fsys driven.Filesystem, file string) ([]targetRef, error) {
	content, err := fsys.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("das Modul targets kann die Doku-Datei %q nicht lesen (DC-FA-TGT-001, fail-closed): %w", file, err)
	}
	var out []targetRef
	for i, line := range splitLines(content) {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		for _, m := range docTargetRe.FindAllStringSubmatch(line, -1) {
			out = append(out, targetRef{name: m[1], file: file, line: i + 1})
		}
	}
	return out, nil
}
