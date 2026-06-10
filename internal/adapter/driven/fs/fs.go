// Package fs ist der driven Filesystem-Adapter — die einzige
// Dateisystem-Tür von d-check (spec/architecture.md §2; Pfade gemäß
// ADR-0005). Alle Port-Pfade sind '/'-getrennt und relativ zur
// Repo-Wurzel.
package fs

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Adapter implementiert driven.Filesystem über das echte Dateisystem.
type Adapter struct {
	Root string // absoluter Pfad der Repo-Wurzel
}

func New(root string) *Adapter { return &Adapter{Root: root} }

func (a *Adapter) abs(rel string) string {
	return filepath.Join(a.Root, filepath.FromSlash(rel))
}

// Kind klassifiziert per Lstat (keine Symlink-Auflösung —
// DC-FA-LINK-002).
func (a *Adapter) Kind(rel string) (driven.EntryKind, error) {
	info, err := os.Lstat(a.abs(rel))
	if err != nil {
		if os.IsNotExist(err) {
			return driven.KindMissing, nil
		}
		return driven.KindMissing, err
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return driven.KindSymlink, nil
	case info.IsDir():
		return driven.KindDir, nil
	default:
		return driven.KindFile, nil
	}
}

func (a *Adapter) ReadFile(rel string) ([]byte, error) {
	return os.ReadFile(a.abs(rel))
}

// List liefert Verzeichniseinträge deterministisch sortiert
// (DC-QA-02).
func (a *Adapter) List(relDir string) ([]driven.DirEntry, error) {
	entries, err := os.ReadDir(a.abs(relDir))
	if err != nil {
		return nil, err
	}
	out := make([]driven.DirEntry, 0, len(entries))
	for _, e := range entries {
		kind := driven.KindFile
		if e.Type()&os.ModeSymlink != 0 {
			kind = driven.KindSymlink
		} else if e.IsDir() {
			kind = driven.KindDir
		}
		out = append(out, driven.DirEntry{Name: e.Name(), Kind: kind})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
