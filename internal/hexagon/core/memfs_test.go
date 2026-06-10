package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// memFS ist ein In-Memory-Filesystem für Kern-Tests (ADR-0004
// Konsequenz: Akzeptanztests ohne echtes Dateisystem).
type memFS struct {
	files    map[string]string // rel → Inhalt
	symlinks map[string]bool   // rel → ist Symlink
}

func newMemFS(files map[string]string) *memFS {
	return &memFS{files: files, symlinks: map[string]bool{}}
}

func (m *memFS) dirs() map[string]bool {
	out := map[string]bool{"": true}
	add := func(p string) {
		segs := strings.Split(p, "/")
		for i := 1; i < len(segs); i++ {
			out[strings.Join(segs[:i], "/")] = true
		}
	}
	for f := range m.files {
		add(f)
	}
	for s := range m.symlinks {
		add(s)
	}
	return out
}

func (m *memFS) Kind(rel string) (driven.EntryKind, error) {
	if m.symlinks[rel] {
		return driven.KindSymlink, nil
	}
	if _, ok := m.files[rel]; ok {
		return driven.KindFile, nil
	}
	if m.dirs()[rel] {
		return driven.KindDir, nil
	}
	return driven.KindMissing, nil
}

func (m *memFS) ReadFile(rel string) ([]byte, error) {
	if c, ok := m.files[rel]; ok {
		return []byte(c), nil
	}
	return nil, fmt.Errorf("not found: %s", rel)
}

func (m *memFS) List(relDir string) ([]driven.DirEntry, error) {
	prefix := ""
	if relDir != "" {
		prefix = relDir + "/"
	}
	seen := map[string]driven.EntryKind{}
	addChild := func(p string, kind driven.EntryKind) {
		if !strings.HasPrefix(p, prefix) || p == relDir {
			return
		}
		rest := strings.TrimPrefix(p, prefix)
		if rest == "" {
			return
		}
		if idx := strings.IndexByte(rest, '/'); idx != -1 {
			seen[rest[:idx]] = driven.KindDir
			return
		}
		if existing, ok := seen[rest]; !ok || existing != driven.KindDir {
			seen[rest] = kind
		}
	}
	for f := range m.files {
		addChild(f, driven.KindFile)
	}
	for s := range m.symlinks {
		addChild(s, driven.KindSymlink)
	}
	var out []driven.DirEntry
	for name, kind := range seen {
		out = append(out, driven.DirEntry{Name: name, Kind: kind})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
