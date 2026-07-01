// Package git ist der driven VCS-Adapter (Modul vcs, DC-FA-VCS-001): er liest
// die git-Historie **rein lesend** aus dem read-only `.git` über die
// reine-Go-Bibliothek go-git — **ohne** externes git-Binary (das
// distroless-Image bleibt unangetastet, ADR-0002/ADR-0024) und **ohne** Netz.
// Einziger go-git-Ort des Repos (arch-check R2-Klasse, „eine Tür je externer
// Abhängigkeit") — und damit die einzige git-Tür von d-check
// (spec/architecture.md §2). Nicht auflösbares Ref/`.git` ⇒ Fehler (fail-closed).
package git

import (
	"errors"
	"fmt"
	"io"
	"sort"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Adapter implementiert driven.VCS über go-git auf einem geöffneten Repository.
type Adapter struct {
	repo *gogit.Repository
}

// Open öffnet das Repository unter root (erwartet root/.git). Fehlt oder ist
// `.git` unlesbar ⇒ Fehler (fail-closed, Exit 2; DC-FA-VCS-001).
func Open(root string) (*Adapter, error) {
	repo, err := gogit.PlainOpen(root)
	if err != nil {
		return nil, fmt.Errorf("kein lesbares git-Repository unter %s: %w", root, err)
	}
	return &Adapter{repo: repo}, nil
}

// treeAt löst ein Ref (z. B. "HEAD", ein Commit-SHA, "main~2") zum Tree auf.
func (a *Adapter) treeAt(ref string) (*object.Tree, error) {
	hash, err := a.repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return nil, err
	}
	commit, err := a.repo.CommitObject(*hash)
	if err != nil {
		return nil, err
	}
	return commit.Tree()
}

// hasHead meldet, ob das Repository einen auflösbaren HEAD trägt (vor dem
// ersten Commit nicht — dann gibt es im --staged-Modus nichts zu schützen).
func (a *Adapter) hasHead() bool {
	_, err := a.repo.Head()
	return err == nil
}

// ChangedPaths liefert die Änderungen zwischen base und head (DC-FA-VCS-001.a
// Schritt 2). head == driven.IndexRef ⇒ staged-Diff (base-Tree vs. Index).
func (a *Adapter) ChangedPaths(base, head string) ([]driven.VCSChange, error) {
	if head == driven.IndexRef {
		// Kein auflösbarer HEAD (erster Commit) ⇒ nichts zu schützen — wie das
		// abgelöste Skript, das im --staged-Modus `git rev-parse --verify HEAD`
		// prüfte und bei Fehlschlag still blieb (jeder Head()-Fehler ⇒ leer).
		if !a.hasHead() {
			return nil, nil
		}
		baseTree, err := a.treeAt(base)
		if err != nil {
			return nil, fmt.Errorf("staged-Basis %q nicht auflösbar: %w", base, err)
		}
		idx, err := a.repo.Storer.Index()
		if err != nil {
			return nil, fmt.Errorf("git-Index nicht lesbar: %w", err)
		}
		return diffTreeIndex(baseTree, idx)
	}
	baseTree, err := a.treeAt(base)
	if err != nil {
		return nil, fmt.Errorf("Range-Basis %q nicht auflösbar: %w", base, err)
	}
	headTree, err := a.treeAt(head)
	if err != nil {
		return nil, fmt.Errorf("Range-Spitze %q nicht auflösbar: %w", head, err)
	}
	return diffTrees(baseTree, headTree)
}

// diffTrees übersetzt einen go-git-Tree-Diff in VCSChanges. Renames erscheinen
// als Delete(alt) + Add(neu) — keine Inhalts-Ähnlichkeits-Erkennung
// (DC-FA-VCS-001 Out-of-Scope); der Delete des immutablen Pfads ist der Befund.
func diffTrees(base, head *object.Tree) ([]driven.VCSChange, error) {
	changes, err := base.Diff(head)
	if err != nil {
		return nil, fmt.Errorf("git-Diff fehlgeschlagen: %w", err)
	}
	out := make([]driven.VCSChange, 0, len(changes))
	for _, c := range changes {
		switch {
		case c.From.Name == "":
			out = append(out, driven.VCSChange{Status: driven.VCSAdded, Path: c.To.Name})
		case c.To.Name == "":
			out = append(out, driven.VCSChange{Status: driven.VCSDeleted, Path: c.From.Name})
		default:
			out = append(out, driven.VCSChange{Status: driven.VCSModified, Path: c.To.Name})
		}
	}
	return out, nil
}

// diffTreeIndex difft den HEAD-Tree gegen den staged Index (rein lesend, ohne
// Working-Tree-Zugriff): in Index, nicht in HEAD ⇒ Added; beide mit
// abweichendem Blob-Hash ⇒ Modified; in HEAD, nicht im Index ⇒ Deleted.
func diffTreeIndex(headTree *object.Tree, idx *index.Index) ([]driven.VCSChange, error) {
	headFiles := map[string]plumbing.Hash{}
	if err := headTree.Files().ForEach(func(f *object.File) error {
		headFiles[f.Name] = f.Hash
		return nil
	}); err != nil {
		return nil, fmt.Errorf("HEAD-Tree nicht lesbar: %w", err)
	}
	idxFiles := make(map[string]plumbing.Hash, len(idx.Entries))
	for _, e := range idx.Entries {
		idxFiles[e.Name] = e.Hash
	}
	var out []driven.VCSChange
	for name, ih := range idxFiles {
		hh, ok := headFiles[name]
		switch {
		case !ok:
			out = append(out, driven.VCSChange{Status: driven.VCSAdded, Path: name})
		case hh != ih:
			out = append(out, driven.VCSChange{Status: driven.VCSModified, Path: name})
		}
	}
	for name := range headFiles {
		if _, ok := idxFiles[name]; !ok {
			out = append(out, driven.VCSChange{Status: driven.VCSDeleted, Path: name})
		}
	}
	return out, nil
}

// CommitMessages liefert die rohen Messages der Nicht-Merge-Commits der Range
// base..head (git rev-list --no-merges-Parität; Modul commits, DC-FA-COMMITS-001).
// Vorgehen: die von base erreichbaren Commits (base + Vorfahren) bilden die
// Ausschluss-Menge; von head aus werden die übrigen Commits durchlaufen und die
// Nicht-Merges (≤ 1 Elternteil) gesammelt. Deterministisch nach Commit-SHA
// sortiert (DC-QA-02). head == driven.IndexRef ist nicht auflösbar ⇒ Fehler
// (fail-closed; die Pending-Message hat der --commit-msg-Kurzschluss-Modus).
func (a *Adapter) CommitMessages(base, head string) ([]driven.CommitMeta, error) {
	if head == driven.IndexRef {
		return nil, fmt.Errorf("commits: --staged wird nicht unterstützt — nutze --range oder --commit-msg (die Pending-Message ist kein Commit)")
	}
	baseHash, err := a.repo.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		return nil, fmt.Errorf("Range-Basis %q nicht auflösbar: %w", base, err)
	}
	headHash, err := a.repo.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		return nil, fmt.Errorf("Range-Spitze %q nicht auflösbar: %w", head, err)
	}
	excluded, err := a.ancestors(*baseHash)
	if err != nil {
		return nil, fmt.Errorf("Range-Basis-Vorfahren nicht lesbar: %w", err)
	}
	visited := map[plumbing.Hash]bool{}
	var metas []driven.CommitMeta
	stack := []plumbing.Hash{*headHash}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[cur] || excluded[cur] {
			continue
		}
		visited[cur] = true
		c, err := a.repo.CommitObject(cur)
		if err != nil {
			return nil, fmt.Errorf("commit %s nicht lesbar: %w", cur, err)
		}
		if len(c.ParentHashes) <= 1 { // Nicht-Merge (--no-merges)
			metas = append(metas, driven.CommitMeta{ShortSHA: cur.String()[:7], Message: c.Message})
		}
		for _, p := range c.ParentHashes {
			if !visited[p] && !excluded[p] {
				stack = append(stack, p)
			}
		}
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].ShortSHA < metas[j].ShortSHA })
	return metas, nil
}

// ancestors sammelt alle von h erreichbaren Commit-Hashes (inkl. h) — die
// Ausschluss-Menge für die Range base..head (der base-seitige rev-list-Ausschluss).
func (a *Adapter) ancestors(h plumbing.Hash) (map[plumbing.Hash]bool, error) {
	set := map[plumbing.Hash]bool{}
	stack := []plumbing.Hash{h}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if set[cur] {
			continue
		}
		set[cur] = true
		c, err := a.repo.CommitObject(cur)
		if err != nil {
			return nil, err
		}
		for _, p := range c.ParentHashes {
			if !set[p] {
				stack = append(stack, p)
			}
		}
	}
	return set, nil
}

// FileAt liefert den Inhalt von path an ref (ref == driven.IndexRef ⇒ staged
// Index); ok=false, wenn an ref abwesend.
func (a *Adapter) FileAt(ref, path string) ([]byte, bool, error) {
	if ref == driven.IndexRef {
		return a.fileFromIndex(path)
	}
	tree, err := a.treeAt(ref)
	if err != nil {
		return nil, false, fmt.Errorf("ref %q nicht auflösbar: %w", ref, err)
	}
	f, err := tree.File(path)
	if errors.Is(err, object.ErrFileNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	s, err := f.Contents()
	if err != nil {
		return nil, false, err
	}
	return []byte(s), true, nil
}

// fileFromIndex liest den staged Blob eines Pfads (DC-FA-VCS-001.a; --staged).
func (a *Adapter) fileFromIndex(path string) ([]byte, bool, error) {
	idx, err := a.repo.Storer.Index()
	if err != nil {
		return nil, false, fmt.Errorf("git-Index nicht lesbar: %w", err)
	}
	for _, e := range idx.Entries {
		if e.Name != path {
			continue
		}
		blob, err := a.repo.BlobObject(e.Hash)
		if err != nil {
			return nil, false, err
		}
		r, err := blob.Reader()
		if err != nil {
			return nil, false, err
		}
		defer func() { _ = r.Close() }()
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, false, err
		}
		return b, true, nil
	}
	return nil, false, nil
}
