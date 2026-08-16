package rules

import (
	"path"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckResolveFromDirs prüft einmal je Lauf, dass jeder Ort jeder Gruppe im
// Baum existiert (fail-closed): ein Tippfehler in einem dirs-Eintrag schaltete
// die Quellen-Rolle sonst still ab — der Fehlzustand wäre von Konsistenz nicht
// unterscheidbar, dieselbe Klasse wie das unlesbare waves.dir.
func CheckResolveFromDirs(fsys driven.Filesystem, groups []model.ResolveFromGroup) []model.Finding {
	var out []model.Finding
	for _, g := range groups {
		for _, d := range append(append([]string{}, g.Dirs...), g.FixedDirs...) {
			if kind, err := fsys.Kind(d); err != nil || kind != driven.KindDir {
				out = append(out, model.Finding{
					File: d, Line: 1, Rule: "links", Target: d,
					Reason: model.ReasonLinkPositionDependent,
					Message: "resolve-from-Verzeichnis " + d + " existiert nicht — " +
						"die Gruppe prüft dort keine Quelle, der Zustand sähe sonst wie Konsistenz aus (fail-closed)",
				})
			}
		}
	}
	return out
}

// CheckResolveFrom prüft ortsfeste Verweise (DC-FA-LINK-001 §Ortsfeste
// Verweise, Schritt 6): eine Datei in einem wandernden Verzeichnis einer
// resolve-from-Gruppe muss jedes relative Ziel von JEDEM Ort der Gruppe
// auflösen — und überall auf dasselbe Ziel. Dateien in fixed-dirs sind am
// Endzustand und keine Quellen. Ohne Gruppen byte-identisch.
func CheckResolveFrom(
	fsys driven.Filesystem, file string, lines []Line,
	groups []model.ResolveFromGroup, ignoreRefs []model.IgnoreRef,
) []model.Finding {
	group := resolveFromGroupOf(groups, path.Dir(file))
	if group == nil {
		return nil
	}
	orte := append(append([]string{}, group.Dirs...), group.FixedDirs...)
	var out []model.Finding
	for _, ref := range ExtractLinks(lines) {
		rel, escaped, ok := localTarget(file, ref.Target)
		if !ok || escaped {
			continue // externe/Anker-Ziele bzw. repo-escape meldet Schritt 4
		}
		if refIgnored(ignoreRefs, file, rel) {
			continue // dasselbe Ventil wie die bestehende Prüfung
		}
		if f := positionDependent(fsys, file, ref, orte); f != nil {
			out = append(out, *f)
		}
	}
	return out
}

// resolveFromGroupOf liefert die Gruppe, in deren Dirs das Verzeichnis liegt —
// nur dann ist die Datei eine Quelle. Der Config-Rand garantiert, dass ein
// Verzeichnis höchstens einer Gruppe als Dirs-Mitglied angehört.
func resolveFromGroupOf(groups []model.ResolveFromGroup, dir string) *model.ResolveFromGroup {
	for i := range groups {
		for _, d := range groups[i].Dirs {
			if path.Clean(d) == path.Clean(dir) {
				return &groups[i]
			}
		}
	}
	return nil
}

// positionDependent löst das Ziel der Referenz von jedem Ort der Gruppe auf.
// Nicht überall auflösbar ODER nicht überall dasselbe Ziel ⇒ ein Befund an der
// Referenz-Zeile; die Meldung nennt den ersten nicht auflösenden Ort bzw. die
// divergierenden Ziele.
func positionDependent(fsys driven.Filesystem, file string, ref LinkRef, orte []string) *model.Finding {
	pathPart := ref.Target
	if idx := strings.IndexByte(pathPart, '#'); idx != -1 {
		pathPart = pathPart[:idx]
	}
	// Vorbedingung: das Ziel löst vom IST-Ort sauber auf — sonst meldet
	// Schritt 4/5 bereits (target-missing, repo-escape oder symlink), und ein
	// zweiter Befund derselben Referenz wäre eine Doppel-Meldung. Die Klasse
	// dieses Schritts ist „am Ist-Ort grün".
	istRel, istEscaped, _ := ResolveTarget(file, pathPart)
	if istKind, err := fsys.Kind(istRel); istEscaped || err != nil || istKind == driven.KindMissing {
		return nil
	}
	if hit, err := symlinkInPath(fsys, istRel); err == nil && hit {
		return nil
	}
	ziele := map[string]bool{}
	for _, ort := range orte {
		hypothetisch := path.Join(ort, path.Base(file))
		rel, escaped, _ := ResolveTarget(hypothetisch, pathPart)
		// escaped zuerst: ein Pfad außerhalb der Wurzel wird nicht abgefragt.
		if !escaped {
			if kind, err := fsys.Kind(rel); err == nil && kind != driven.KindMissing {
				ziele[rel] = true
				continue
			}
		}
		return &model.Finding{
			File: file, Line: ref.Line, Rule: "links",
			Target: ref.Target, Reason: model.ReasonLinkPositionDependent,
			Message: "Verweis löst von " + ort + " nicht auf — er bricht, sobald die Datei dorthin wandert" +
				" (Reparatur: Pfad präfixieren, nicht Ziel anlegen)",
		}
	}
	if len(ziele) > 1 {
		namen := make([]string, 0, len(ziele))
		for z := range ziele {
			namen = append(namen, z)
		}
		sort.Strings(namen) // DC-QA-02: unabhängig von der Map-Iteration
		return &model.Finding{
			File: file, Line: ref.Line, Rule: "links",
			Target: ref.Target, Reason: model.ReasonLinkPositionDependent,
			Message: "Verweis löst je nach Ort auf verschiedene Ziele auf (" +
				strings.Join(namen, " · ") + ") — er meint je nach Lifecycle-Ort etwas anderes",
		}
	}
	return nil
}
