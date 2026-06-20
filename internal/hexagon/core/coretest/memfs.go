// Package coretest stellt Test-Helfer für die Kern-Pakete bereit —
// insbesondere ein In-Memory-Filesystem, das den `driven.Filesystem`-Port
// ohne echtes Dateisystem erfüllt (ADR-0004-Konsequenz: Akzeptanztests
// ohne I/O). Es wird von den White-box-Tests aller Kern-Pakete genutzt und
// liegt deshalb in einem eigenen, von ihnen importierbaren Paket.
package coretest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// MemFS ist ein In-Memory-Filesystem für Kern-Tests.
type MemFS struct {
	files    map[string]string // rel → Inhalt
	symlinks map[string]bool   // rel → ist Symlink
}

// NewMemFS erzeugt ein MemFS aus einer rel→Inhalt-Abbildung.
func NewMemFS(files map[string]string) *MemFS {
	return &MemFS{files: files, symlinks: map[string]bool{}}
}

// AddFile fügt nach der Konstruktion eine Datei (rel→Inhalt) hinzu.
func (m *MemFS) AddFile(rel, content string) { m.files[rel] = content }

// AddSymlink markiert einen Pfad als Symlink.
func (m *MemFS) AddSymlink(rel string) { m.symlinks[rel] = true }

func (m *MemFS) dirs() map[string]bool {
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

// Kind erfüllt den driven.Filesystem-Port.
func (m *MemFS) Kind(rel string) (driven.EntryKind, error) {
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

// ReadFile erfüllt den driven.Filesystem-Port.
func (m *MemFS) ReadFile(rel string) ([]byte, error) {
	if c, ok := m.files[rel]; ok {
		return []byte(c), nil
	}
	return nil, fmt.Errorf("not found: %s", rel)
}

// List erfüllt den driven.Filesystem-Port.
func (m *MemFS) List(relDir string) ([]driven.DirEntry, error) {
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
