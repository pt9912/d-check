package rules

import (
	"net/url"
	"path"
	"strings"
)

// ResolveTarget löst ein Linkziel relativ zur enthaltenden Datei auf:
// vollständige Prozent-Dekodierung (RFC 3986) VOR der
// Repo-Escape-Prüfung, dann lexikalische Normalisierung
// (DC-FA-LINK-001, spec/spezifikation.md §DC-FA-LINK-001.a Schritt 4).
// escaped=true bedeutet: der Pfad verlässt die Repo-Wurzel lexikalisch.
func ResolveTarget(fromFile, target string) (rel string, escaped bool, ok bool) {
	decoded, err := url.PathUnescape(target)
	if err != nil {
		// nicht dekodierbar → wie roh behandeln (Befund entsteht
		// dann über target-missing)
		decoded = target
	}
	var joined string
	if strings.HasPrefix(decoded, "/") {
		// absolute Ziele werden relativ zur Repo-Wurzel interpretiert
		joined = path.Clean(decoded)
		joined = strings.TrimPrefix(joined, "/")
	} else {
		joined = path.Join(path.Dir(fromFile), decoded)
	}
	if joined == ".." || strings.HasPrefix(joined, "../") {
		return joined, true, true
	}
	if joined == "." {
		joined = ""
	}
	return joined, false, true
}

// ResolveConfigPath normalisiert einen in der Konfiguration
// deklarierten, zur Repo-Wurzel relativen Pfad (lexikalisch;
// DC-FA-CONF-001). rel == "" steht für die Repo-Wurzel selbst;
// escaped meldet Pfade, die die Repo-Wurzel verlassen — gemeinsamer
// Helfer für Scan-Wurzeln und ids-Targets.
func ResolveConfigPath(p string) (rel string, escaped bool) {
	rel = path.Clean(strings.Trim(p, "/"))
	if rel == "." {
		return "", false
	}
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return rel, true
	}
	return rel, false
}

// localTarget zerlegt ein Linkziel in seinen Pfad-Teil (Fragment
// abgetrennt — Anker sind Aufgabe von `anchors`) und löst ihn relativ
// zur enthaltenden Datei auf. ok=false für leere, reine Anker- und
// externe Ziele. Gemeinsamer Kern der Ziel-Auflösung von `links` und
// `matrix`.
func localTarget(file, target string) (rel string, escaped, ok bool) {
	if target == "" || strings.HasPrefix(target, "#") || IsExternalScheme(target) {
		return "", false, false
	}
	pathPart := target
	if idx := strings.IndexByte(pathPart, '#'); idx != -1 {
		pathPart = pathPart[:idx]
	}
	if pathPart == "" {
		return "", false, false
	}
	rel, escaped, _ = ResolveTarget(file, pathPart)
	return rel, escaped, true
}

// IsExternalScheme prüft auf externes URL-Schema (http:, mailto:, …).
func IsExternalScheme(target string) bool {
	for i := 0; i < len(target); i++ {
		c := target[i]
		if c == ':' {
			return i > 0
		}
		isAlpha := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isExtra := c == '+' || c == '.' || c == '-' || (c >= '0' && c <= '9' && i > 0)
		if i == 0 && !isAlpha {
			return false
		}
		if !isAlpha && !isExtra {
			return false
		}
	}
	return false
}

// matchGlob prüft ein Ignore-Muster gegen einen '/'-relativen Pfad.
// Unterstützt '*' und '?' segmentweise sowie '**' für beliebig viele
// Segmente (spec/spezifikation.md §2, scan.ignore).
func matchGlob(pattern, rel string) bool {
	return matchSegs(strings.Split(pattern, "/"), strings.Split(rel, "/"))
}

func matchSegs(pat, segs []string) bool {
	if len(pat) == 0 {
		return len(segs) == 0
	}
	if pat[0] == "**" {
		for i := 0; i <= len(segs); i++ {
			if matchSegs(pat[1:], segs[i:]) {
				return true
			}
		}
		return false
	}
	if len(segs) == 0 {
		return false
	}
	if ok, _ := path.Match(pat[0], segs[0]); !ok {
		return false
	}
	return matchSegs(pat[1:], segs[1:])
}

// ignored prüft, ob ein '/'-relativer Pfad von einem der Ignore-Muster
// getroffen wird (scan.ignore; geteilt von ids und der Discovery).
func ignored(rel string, ignore []string) bool {
	for _, pat := range ignore {
		if matchGlob(pat, rel) {
			return true
		}
	}
	return false
}
