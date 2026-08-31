package rules

import (
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckVCS ist das Regelmodul vcs (DC-FA-VCS-001): es vergleicht den **Core**
// einer immutablen Datei über zwei git-Stände (`core(BASE)` vs. `core(HEAD)`)
// und meldet Körper-Drift, unzulässigen Status-Übergang (HeadAllow) oder
// Löschung/Umbenennung als `core-drift-vcs`. Opt-in (leere Paths ⇒ inert);
// es liest die git-Historie über den VCS-Port (nicht-hermetisch, aber lokal/
// lesend/deterministisch) — ein Port-Fehler (fehlendes `.git`/Range) wird als
// error zurückgegeben, den der Aufrufer auf Exit 2 abbildet (fail-closed).
// Diagnose-only: kein `--repair`-Hunk.
func CheckVCS(vcs driven.VCS, cfg model.VCSConfig, base, head string) ([]model.Finding, error) {
	if len(cfg.Paths) == 0 || vcs == nil {
		return nil, nil // inert (Modul ohne Klassen-Config oder ohne Port)
	}
	changes, err := vcs.ChangedPaths(base, head)
	if err != nil {
		return nil, err
	}
	var findings []model.Finding
	for _, c := range changes {
		if !ignored(c.Path, cfg.Paths) {
			continue // nicht in der geschützten Klasse (vcs.paths)
		}
		switch c.Status {
		case driven.VCSAdded:
			// neue Datei ist noch nicht immutabel (Proposed→Accepted-Reifung frei)
		case driven.VCSDeleted:
			f, err := vcsDeleted(vcs, cfg, base, c.Path)
			if err != nil {
				return nil, err
			}
			findings = append(findings, f...)
		case driven.VCSModified:
			f, err := vcsModified(vcs, cfg, base, head, c.Path)
			if err != nil {
				return nil, err
			}
			findings = append(findings, f...)
		}
	}
	return findings, nil
}

// vcsDeleted meldet die Löschung/Umbenennung einer immutablen BASE-Datei
// (der Pfad einer immutablen Datei ist stabil). Dass ein Rename überhaupt als
// Delete-Hälfte hier ankommt, hält die Diff-Übersetzung im VCS-Adapter, nicht
// diese Funktion.
func vcsDeleted(vcs driven.VCS, cfg model.VCSConfig, base, path string) ([]model.Finding, error) {
	content, ok, err := vcs.FileAt(base, path)
	if err != nil {
		return nil, err
	}
	if !ok || !baseImmutable(content, cfg.ImmutableWhen) {
		return nil, nil
	}
	return []model.Finding{{
		File: path, Line: 1, Rule: "vcs", Target: path,
		Reason:  model.ReasonCoreDriftVCS,
		Message: "immutable Datei gelöscht oder umbenannt — der Pfad einer immutablen Datei ist stabil",
	}}, nil
}

// vcsModified vergleicht den Core einer modifizierten immutablen Datei und
// prüft den Status-Übergang. BASE nicht immutabel (z. B. Proposed) ⇒ frei
// (Grandfathering, DC-FA-VCS-001.a Schritt 3).
func vcsModified(vcs driven.VCS, cfg model.VCSConfig, base, head, path string) ([]model.Finding, error) {
	baseContent, ok, err := vcs.FileAt(base, path)
	if err != nil {
		return nil, err
	}
	if !ok || !baseImmutable(baseContent, cfg.ImmutableWhen) {
		return nil, nil
	}
	headContent, ok, err := vcs.FileAt(head, path)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	line := vcsHeadStatusLineNo(headContent, splitLines(headContent), cfg.StatusLine)
	if line == 0 {
		line = 1
	}
	var findings []model.Finding
	if vcsCore(baseContent, cfg.StatusLine, cfg.ExcludeSections) !=
		vcsCore(headContent, cfg.StatusLine, cfg.ExcludeSections) {
		findings = append(findings, model.Finding{
			File: path, Line: line, Rule: "vcs", Target: path,
			Reason:  model.ReasonCoreDriftVCS,
			Message: "Core einer immutablen Datei hat sich über die Commit-Range geändert",
		})
	}
	if cfg.HeadAllow != nil && !cfg.HeadAllow.MatchString(metaStatusLine(headContent, cfg.StatusLine)) {
		findings = append(findings, model.Finding{
			File: path, Line: line, Rule: "vcs", Target: path,
			Reason:  model.ReasonCoreDriftVCS,
			Message: "unzulässiger Status-Übergang einer immutablen Datei",
		})
	}
	return findings, nil
}

// baseImmutable: die BASE-Version gilt als immutabel, wenn irgendeine Zeile
// ImmutableWhen matcht (erstes Vorkommen; DC-FA-VCS-001.a Schritt 3).
func baseImmutable(content []byte, when *regexp.Regexp) bool {
	if when == nil {
		return false
	}
	prose := proseLineSet(content)
	for i, ln := range splitLines(content) {
		if !prose[i+1] {
			continue // ein Accepted-Beispiel im Code-Block macht nichts immutabel
		}
		if when.MatchString(ln) {
			return true
		}
	}
	return false
}

// vcsCore bildet den normalisierten Core: roher Inhalt ohne die
// exclude-sections-Abschnitte und ohne die Kopf-Status-Zeile, mit der
// reflow-invarianten Whitespace-Normalisierung von immutable/pins
// (DC-FA-VCS-001.a Schritt 4).
func vcsCore(content []byte, statusLine *regexp.Regexp, excludeSections []string) string {
	excluded := excludedRanges(content, excludeSections)
	lines := splitLines(content)
	strip := vcsHeadStatusLineNo(content, lines, statusLine)
	var b strings.Builder
	for i, raw := range lines {
		no := i + 1
		if no == strip || inRanges(excluded, no) {
			continue
		}
		b.WriteString(raw)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(pinWhitespaceRE.ReplaceAllString(b.String(), " "))
}

// vcsHeadStatusLineNo liefert die 1-basierte Zeilennummer der **Kopf**-Status-
// Zeile: das erste statusLine-Vorkommen **vor** der ersten `## `-H2; 0 ohne
// Treffer (oder ohne statusLine). Eine gleichlautende Zeile im Körper (nach der
// ersten H2) bleibt damit Teil des Core (DC-FA-VCS-001.a Schritt 4).
func vcsHeadStatusLineNo(content []byte, lines []string, statusLine *regexp.Regexp) int {
	if statusLine == nil {
		return 0
	}
	prose := proseLineSet(content)
	for i, raw := range lines {
		if !prose[i+1] {
			continue // Fence-Inneres: ein Beispiel-Kopf ist kein Status
		}
		if level, _, ok := parseATXHeading(raw); ok && level >= 2 {
			return 0
		}
		if statusLine.MatchString(raw) {
			return i + 1
		}
	}
	return 0
}

// metaStatusLine liefert den Text der Kopf-Status-Zeile (für die HeadAllow-
// Übergangs-Prüfung); "" ohne Treffer.
func metaStatusLine(content []byte, statusLine *regexp.Regexp) string {
	lines := splitLines(content)
	no := vcsHeadStatusLineNo(content, lines, statusLine)
	if no == 0 {
		return ""
	}
	return lines[no-1]
}

// splitLines zerlegt rohen Datei-Inhalt in Zeilen ('\n'-getrennt).
func splitLines(content []byte) []string {
	return strings.Split(string(content), "\n")
}
