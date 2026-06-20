package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckLinks ist das Regelmodul `links` (DC-FA-LINK-001/002): lokale
// Link- und Bildziele müssen existieren, innerhalb der Repo-Wurzel
// liegen und symlink-frei sein. Fragment-Teile gemischter Ziele
// ignoriert das Modul (Aufgabe von `anchors`); reine Anker-Ziele und
// externe Schemata werden übersprungen.
func CheckLinks(fsys driven.Filesystem, file string, lines []Line) []model.Finding {
	var findings []model.Finding
	for _, ref := range ExtractLinks(lines) {
		target := ref.Target
		rel, escaped, ok := localTarget(file, target)
		if !ok {
			continue
		}

		// Symlink-Prüfung hat Vorrang: alle Komponenten des lexikalisch
		// aufgelösten Ziels innerhalb der Repo-Wurzel per Lstat
		// (DC-FA-LINK-002; genau ein Befund pro Linkziel).
		if !escaped {
			if hit, err := symlinkInPath(fsys, rel); err == nil && hit {
				findings = append(findings, model.Finding{
					File: file, Line: ref.Line, Rule: "links",
					Target: target, Reason: model.ReasonSymlink,
					Message: "Linkziel ist oder enthält einen Symlink",
				})
				continue
			}
		}
		if escaped {
			findings = append(findings, model.Finding{
				File: file, Line: ref.Line, Rule: "links",
				Target: target, Reason: model.ReasonRepoEscape,
				Message: "Linkziel verlässt die Repository-Wurzel",
			})
			continue
		}
		kind, err := fsys.Kind(rel)
		if err != nil || kind == driven.KindMissing {
			findings = append(findings, model.Finding{
				File: file, Line: ref.Line, Rule: "links",
				Target: target, Reason: model.ReasonTargetMissing,
				Message: "Linkziel existiert nicht",
			})
		}
	}
	return findings
}

// symlinkInPath prüft jede Pfad-Komponente per Lstat auf Symlink.
func symlinkInPath(fsys driven.Filesystem, rel string) (bool, error) {
	if rel == "" {
		return false, nil
	}
	segs := strings.Split(rel, "/")
	for i := 1; i <= len(segs); i++ {
		prefix := strings.Join(segs[:i], "/")
		kind, err := fsys.Kind(prefix)
		if err != nil {
			return false, err
		}
		if kind == driven.KindSymlink {
			return true, nil
		}
		if kind == driven.KindMissing {
			return false, nil
		}
	}
	return false, nil
}
