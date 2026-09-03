// Zieht repo-weite Markdown-Verweise auf archivierte Dateien nach. Eigene,
// self-contained Pfad-Aufloesung -- dieses Werkzeug ist fuer jedes Repo mit
// demselben Regelwerk gedacht, nicht an d-checks interne Pakete gekoppelt,
// die sich ohnehin nicht modul-uebergreifend importieren liessen.
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var mdLinkRE = regexp.MustCompile(`\]\(([^)\s]+)\)`)

// Move beschreibt eine archivierte Datei: alter und neuer Repo-relativer
// Pfad, Forward-Slash-normalisiert wie in Markdown-Links.
type Move struct {
	Old string
	New string
}

// resolveLink loest ein Link-Ziel relativ zum Verzeichnis der linktragenden
// Datei (sourceRel, Repo-relativ) in einen bereinigten Repo-relativen Pfad
// auf. Ein Fragment (#anchor) wird abgetrennt und separat zurueckgegeben --
// es ist nicht Teil der Pfad-Identitaet, bleibt aber beim Umschreiben
// erhalten.
func resolveLink(sourceRel, target string) (resolved, fragment string) {
	if i := strings.IndexByte(target, '#'); i >= 0 {
		fragment = target[i:]
		target = target[:i]
	}
	if target == "" {
		return "", fragment
	}
	if strings.HasPrefix(target, "/") {
		return path.Clean(strings.TrimPrefix(target, "/")), fragment
	}
	return path.Clean(path.Join(path.Dir(sourceRel), target)), fragment
}

// relativize baut aus einem Repo-relativen Zielpfad das Link-Ziel, wie es
// von sourceRel aus geschrieben werden muesste (geschwister-relativ).
func relativize(sourceRel, targetRel string) string {
	rel, err := filepath.Rel(path.Dir(sourceRel), targetRel)
	if err != nil {
		return targetRel
	}
	return filepath.ToSlash(rel)
}

// RewriteFile durchsucht den Inhalt einer Markdown-Datei (Repo-relativer
// Pfad selfRel) nach Links, deren aufgeloestes Ziel einem Move.Old
// entspricht, und ersetzt sie durch das neue, wieder relativierte Ziel samt
// erhaltenem Fragment. Liefert den (ggf. unveraenderten) Inhalt und die
// Anzahl vorgenommener Ersetzungen.
func RewriteFile(selfRel, content string, moves []Move) (string, int) {
	byOld := make(map[string]string, len(moves))
	for _, m := range moves {
		byOld[m.Old] = m.New
	}
	n := 0
	out := mdLinkRE.ReplaceAllStringFunc(content, func(match string) string {
		target := match[2 : len(match)-1] // "](" ... ")" abtrennen
		resolved, fragment := resolveLink(selfRel, target)
		newTarget, ok := byOld[resolved]
		if !ok {
			return match
		}
		n++
		return "](" + relativize(selfRel, newTarget) + fragment + ")"
	})
	return out, n
}

// RewriteRepo wendet RewriteFile auf jede .md-Datei unter root an (rekursiv,
// `.git` ausgenommen) und schreibt geaenderte Dateien zurueck. Liefert die
// betroffenen Dateien (Repo-relativ) mit ihrer Treffer-Anzahl, sortiert.
func RewriteRepo(root string, moves []Move) (map[string]int, error) {
	hits := map[string]int{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".md") {
			return nil
		}
		rel := filepath.ToSlash(relOrSelf(root, p))
		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("%s lesen: %w", p, err)
		}
		newContent, n := RewriteFile(rel, string(b), moves)
		if n == 0 {
			return nil
		}
		hits[rel] = n
		if err := os.WriteFile(p, []byte(newContent), 0o644); err != nil {
			return fmt.Errorf("%s schreiben: %w", p, err)
		}
		return nil
	})
	return hits, err
}

// PreviewRewrites liefert dieselbe Treffer-Menge wie RewriteRepo, OHNE
// etwas zu schreiben -- die Grundlage von --dry-run.
func PreviewRewrites(root string, moves []Move) (map[string]int, error) {
	hits := map[string]int{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".md") {
			return nil
		}
		rel := filepath.ToSlash(relOrSelf(root, p))
		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("%s lesen: %w", p, err)
		}
		_, n := RewriteFile(rel, string(b), moves)
		if n > 0 {
			hits[rel] = n
		}
		return nil
	})
	return hits, err
}

func relOrSelf(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return rel
}

// SortedKeys liefert die Schluessel einer Treffer-Map sortiert -- fuer eine
// deterministische --dry-run-Ausgabe.
func SortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
