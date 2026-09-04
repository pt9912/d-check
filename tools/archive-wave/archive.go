// Fuehrt die Archivierung aus: baut archiv.zip, schreibt beide Stub-Arten,
// entfernt die alten Volltexte. Liefert die vorgenommenen Moves, damit der
// Aufrufer anschliessend RewriteRepo darauf anwenden kann.
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Plan ist das Ergebnis des Sammeln-Schritts: alles, was zu welleID gehoert.
type Plan struct {
	WelleID   string
	WellePlan string // absoluter Pfad, leer fuer eine Welle ohne Plan-Datei
	Slices    []string
	Reviews   []string
}

// RelPath macht einen absoluten Pfad Repo-relativ, Forward-Slash-normalisiert
// -- die Form, in der Markdown-Links Ziele schreiben.
func RelPath(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return abs
	}
	return filepath.ToSlash(rel)
}

// Apply fuehrt die Archivierung aus: legt docs/plan/planning/done/<welleID>/
// an, baut archiv.zip, schreibt beide Stub-Arten, loescht die alten
// Volltexte (Review-Reports bekommen keinen Stub -- Kanon-Vorgabe). Liefert
// die Liste der Pfad-Aenderungen fuer den anschliessenden Verweis-Nachzug.
//
// Kein Rollback bei einem Fehler mittendrin -- derselbe Umgang wie bei jedem
// git-mv-Schritt: der Bediener sieht den Zwischenstand per `git status` und
// entscheidet.
func Apply(root string, p Plan) ([]Move, error) {
	archiveDir := filepath.Join(root, "docs/plan/planning/done", p.WelleID)
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return nil, fmt.Errorf("%s anlegen: %w", archiveDir, err)
	}

	var all []string
	if p.WellePlan != "" {
		all = append(all, p.WellePlan)
	}
	all = append(all, p.Slices...)
	all = append(all, p.Reviews...)
	if err := buildZip(root, filepath.Join(archiveDir, "archiv.zip"), all); err != nil {
		return nil, err
	}

	var moves []Move

	// Eine Welle ohne Plan-Datei (welle-60..66, vor der Plan-Datei-Konvention,
	// slice-191) hat nichts, das ein Stub ersetzen koennte -- ueberspringen
	// statt eine Datei zu erfinden.
	if p.WellePlan != "" {
		welleTitle, err := readTitle(p.WellePlan)
		if err != nil {
			return nil, err
		}
		resultsFile := p.WelleID + "-results.md"
		welleStub := WelleStub(p.WelleID, welleTitle, resultsFile, len(p.Slices), len(p.Reviews))
		welleNewAbs := filepath.Join(archiveDir, filepath.Base(p.WellePlan))
		if err := os.WriteFile(welleNewAbs, []byte(welleStub), 0o644); err != nil {
			return nil, fmt.Errorf("%s schreiben: %w", welleNewAbs, err)
		}
		if err := os.Remove(p.WellePlan); err != nil {
			return nil, fmt.Errorf("%s loeschen: %w", p.WellePlan, err)
		}
		moves = append(moves, Move{Old: RelPath(root, p.WellePlan), New: RelPath(root, welleNewAbs)})
	}

	for _, s := range p.Slices {
		raw, err := os.ReadFile(s)
		if err != nil {
			return nil, fmt.Errorf("%s lesen: %w", s, err)
		}
		title := ExtractTitle(string(raw))
		if title == "" {
			return nil, fmt.Errorf("%s: keine Ueberschrift gefunden", s)
		}
		hervorgegangen := FormatHervorgegangen(ExtractSurvivingIDs(string(raw)))
		field, err := ReadWelleField(s)
		if err != nil {
			return nil, err
		}
		id := SliceIDFromPath(s)
		newAbs := filepath.Join(archiveDir, filepath.Base(s))
		// Ein Feld-Link ist von der ALTEN Slice-Position aus geschrieben (oft
		// per "Ortsfeste Verweise"-Idiom "../done/X", das nur von einer der
		// vier Lifecycle-Wurzeln aus aufloest) -- der Stub liegt aber eine
		// Ebene TIEFER (done/<welle-id>/). Ein blosses Uebernehmen brach hier
		// gemessen (welle-70): RewriteRepo laeuft erst NACH dem Schreiben und
		// loest von der neuen, tieferen Position aus falsch auf.
		field = RewriteFieldForMove(RelPath(root, s), RelPath(root, newAbs), field, moves)
		stub := SliceStub(id, title, field, p.WelleID, hervorgegangen)
		if err := os.WriteFile(newAbs, []byte(stub), 0o644); err != nil {
			return nil, fmt.Errorf("%s schreiben: %w", newAbs, err)
		}
		if err := os.Remove(s); err != nil {
			return nil, fmt.Errorf("%s loeschen: %w", s, err)
		}
		moves = append(moves, Move{Old: RelPath(root, s), New: RelPath(root, newAbs)})
	}

	// Review-Reports bekommen KEINEN Stub -- sie haben keine Identitaet
	// jenseits ihres Slice und liegen nur noch im Archiv.
	for _, r := range p.Reviews {
		if err := os.Remove(r); err != nil {
			return nil, fmt.Errorf("%s loeschen: %w", r, err)
		}
	}

	sort.Slice(moves, func(i, j int) bool { return moves[i].Old < moves[j].Old })
	return moves, nil
}

// ApplySlice archiviert einen einzelnen wellenlosen Slice (Baseline-Regelwerk
// `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht, Zeile
// "Zeitdokumente archivieren ... Ohne Wellen tut es die Slice-Closure
// selbst"). Anders als Apply (Wellen-Modus) WANDERT die Slice-Datei nicht --
// ihr Pfad bleibt, nur ihr Inhalt wird durch den Stub ersetzt, und das
// Archiv liegt flach daneben (`<sliceID>-archiv.zip`), nicht in einem
// Wellen-Unterverzeichnis. Es gibt deshalb keinen Move fuer den Slice
// selbst -- nur seine Review-Reports verschwinden ersatzlos, wie im
// Wellen-Modus auch.
func ApplySlice(root, sliceID, slicePath string, reviews []string) error {
	zipPath := filepath.Join(filepath.Dir(slicePath), sliceID+"-archiv.zip")
	all := append([]string{slicePath}, reviews...)
	if err := buildZip(root, zipPath, all); err != nil {
		return err
	}

	raw, err := os.ReadFile(slicePath)
	if err != nil {
		return fmt.Errorf("%s lesen: %w", slicePath, err)
	}
	title := ExtractTitle(string(raw))
	if title == "" {
		return fmt.Errorf("%s: keine Ueberschrift gefunden", slicePath)
	}
	hervorgegangen := FormatHervorgegangen(ExtractSurvivingIDs(string(raw)))
	field, err := ReadWelleField(slicePath)
	if err != nil {
		return err
	}
	stub := SliceStubStandalone(sliceID, title, field, hervorgegangen)
	if err := os.WriteFile(slicePath, []byte(stub), 0o644); err != nil {
		return fmt.Errorf("%s schreiben: %w", slicePath, err)
	}

	for _, r := range reviews {
		if err := os.Remove(r); err != nil {
			return fmt.Errorf("%s loeschen: %w", r, err)
		}
	}
	return nil
}

func readTitle(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s lesen: %w", path, err)
	}
	title := ExtractTitle(string(b))
	if title == "" {
		return "", fmt.Errorf("%s: keine Ueberschrift gefunden", path)
	}
	return title, nil
}

func buildZip(root, zipPath string, files []string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("%s anlegen: %w", zipPath, err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, abs := range files {
		if err := addToZip(zw, root, abs); err != nil {
			_ = zw.Close()
			return err
		}
	}
	return zw.Close()
}

func addToZip(zw *zip.Writer, root, abs string) error {
	rel := RelPath(root, abs)
	w, err := zw.Create(rel)
	if err != nil {
		return fmt.Errorf("zip-Eintrag %s: %w", rel, err)
	}
	src, err := os.Open(abs)
	if err != nil {
		return fmt.Errorf("%s lesen: %w", abs, err)
	}
	defer src.Close()
	_, err = io.Copy(w, src)
	return err
}
