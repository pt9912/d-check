package rules

import (
	"path"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// CheckResolveFromDirs prüft einmal je Lauf die Gruppen-Orte (fail-closed):
// existiert KEIN dirs-Ort einer Gruppe, zeigt sie sicher ins Leere (Tippfehler
// im Stamm-Pfad), und ein Ort, der als DATEI existiert, ist sicher falsch.
// Ein EINZELNER fehlender Ort meldet bewusst nicht: git überträgt leere
// Verzeichnisse nicht — ein legitim geleertes Lifecycle-Verzeichnis fehlt auf
// jedem frischen Klon und wäre von einem Tippfehler nicht unterscheidbar
// (benannte Grenze im Vertrag).
func CheckResolveFromDirs(fsys driven.Filesystem, groups []model.ResolveFromGroup) []model.Finding {
	var out []model.Finding
	for _, g := range groups {
		existiert := 0
		for _, d := range g.Dirs {
			kind, err := fsys.Kind(d)
			if err == nil && kind == driven.KindDir {
				existiert++
			}
			if err == nil && kind == driven.KindFile {
				out = append(out, resolveDirFinding(d,
					"resolve-from-Ort "+d+" existiert als Datei, nicht als Verzeichnis (fail-closed)"))
			}
		}
		if existiert == 0 {
			out = append(out, resolveDirFinding(g.Dirs[0],
				"kein dirs-Verzeichnis der resolve-from-Gruppe existiert ("+strings.Join(g.Dirs, ", ")+
					") — die Gruppe prüft keine einzige Quelle (fail-closed)"))
		}
	}
	return out
}

// resolveDirFinding ist die Befund-Form der Gruppen-Ort-Prüfung.
func resolveDirFinding(target, msg string) model.Finding {
	return model.Finding{
		File: target, Line: 1, Rule: "links", Target: target,
		Reason: model.ReasonLinkPositionDependent, Message: msg,
	}
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
