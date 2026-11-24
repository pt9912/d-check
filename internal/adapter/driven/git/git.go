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
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Adapter implementiert driven.VCS über go-git auf einem geöffneten Repository.
type Adapter struct {
	repo *gogit.Repository
}

// Open öffnet das Repository unter root (erwartet root/.git). Fehlt oder ist
// `.git` unlesbar ⇒ Fehler (fail-closed, Exit 2; DC-FA-VCS-001). Die
// Objekt-Storage wird — anders als über das bequeme gogit.PlainOpen — über
// packAliasFS aufgebaut, damit ein Pack unter unkanonischem Präfix nicht
// unsichtbar bleibt (eingehender CR, ai-harness-init, 2026-09-17).
func Open(root string) (*Adapter, error) {
	wt := osfs.New(root)
	dot, err := wt.Chroot(gogit.GitDirName)
	if err != nil {
		return nil, fmt.Errorf("kein lesbares git-Repository unter %s: %w", root, err)
	}
	if _, err := dot.Stat(""); err != nil {
		return nil, fmt.Errorf("kein lesbares git-Repository unter %s: %w", root, err)
	}
	s := filesystem.NewStorage(newPackAliasFS(dot), cache.NewObjectLRUDefault())
	repo, err := gogit.Open(s, wt)
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

// AllPaths liefert für base..head (head == driven.IndexRef ⇒ staged:
// base-Tree gegen den Index) JEDEN Datei-Pfad an beiden Enden — direkt über
// die Baum-Struktur aufgelöst, nicht über einen Diff (DC-FA-VCS-001.a
// Schritt 2). Kein auflösbarer HEAD im staged-Modus (erster Commit) ⇒ nichts
// zu schützen — wie das abgelöste Skript, das im --staged-Modus
// `git rev-parse --verify HEAD` prüfte und bei Fehlschlag still blieb.
func (a *Adapter) AllPaths(base, head string) ([]string, []string, error) {
	if head == driven.IndexRef {
		if !a.hasHead() {
			return nil, nil, nil
		}
		baseAll, err := a.pathsAt(base)
		if err != nil {
			return nil, nil, fmt.Errorf("staged-Basis %q nicht auflösbar: %w", base, err)
		}
		idx, err := a.repo.Storer.Index()
		if err != nil {
			return nil, nil, fmt.Errorf("git-Index nicht lesbar: %w", err)
		}
		headAll := make([]string, 0, len(idx.Entries))
		for _, e := range idx.Entries {
			headAll = append(headAll, e.Name)
		}
		return baseAll, headAll, nil
	}
	// DC-FA-VCS-002: der Leerfall wird VOR dem Tree-Vergleich geprüft —
	// Basis und Spitze benennen denselben Commit, oder die Spitze ist von
	// der Basis aus erreichbar (rev-list --count base..head = 0). Dann ist
	// der Prüfbereich leer und das Grün wäre still: „geprüft, nichts
	// gefunden" hätte keinen Gegenstand. --staged löst den Check nicht aus
	// (der Index gegen HEAD braucht keine Range-Tiefe).
	baseHash, err := a.repo.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		return nil, nil, fmt.Errorf("Range-Basis %q nicht auflösbar: %w", base, err)
	}
	headHash, err := a.repo.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		return nil, nil, fmt.Errorf("Range-Spitze %q nicht auflösbar: %w", head, err)
	}
	if *baseHash == *headHash {
		return nil, nil, fmt.Errorf("Range-Leerfall %q..%q — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft", base, head)
	}
	baseAnc, err := a.ancestors(*baseHash)
	if err != nil {
		return nil, nil, fmt.Errorf("Range-Basis-Vorfahren nicht lesbar: %w", err)
	}
	if baseAnc[*headHash] {
		return nil, nil, fmt.Errorf("Range-Leerfall %q..%q — 0 Commits, es wurde nichts geprüft", base, head)
	}
	baseAll, err := a.pathsAt(base)
	if err != nil {
		return nil, nil, fmt.Errorf("Range-Basis %q nicht auflösbar: %w", base, err)
	}
	headAll, err := a.pathsAt(head)
	if err != nil {
		return nil, nil, fmt.Errorf("Range-Spitze %q nicht auflösbar: %w", head, err)
	}
	return baseAll, headAll, nil
}

// pathsAt enumeriert jeden Datei-Pfad im Tree von ref (fail-closed über
// walkTree).
func (a *Adapter) pathsAt(ref string) ([]string, error) {
	tree, err := a.treeAt(ref)
	if err != nil {
		return nil, err
	}
	var out []string
	if err := walkTree(a.repo, tree, "", &out); err != nil {
		return nil, fmt.Errorf("nicht vollständig lesbarer Tree zu %q: %w", ref, err)
	}
	return out, nil
}

// walkTree sammelt jeden Datei-Pfad (regulär/ausführbar/symlink) unter tree,
// rekursiv über *jeden* Unterbaum. Anders als der geteilte Walker hinter
// Tree.Files() (der einen nicht ladbaren Unterbaum in ein stilles io.EOF
// verwandelt, CO-001 dritte Ausprägung) bricht dieser hier mit einem Fehler
// ab, sobald ein Unterbaum-Objekt nicht lesbar ist — der Aufruf über
// repo.TreeObject liest das Objekt tatsächlich, statt nur den Tree-Eintrag zu
// prüfen. Ein Gitlink (Submodul) trägt keinen Blob dieses Repos und wird
// übersprungen, nicht aufgelöst.
func walkTree(repo *gogit.Repository, tree *object.Tree, prefix string, out *[]string) error {
	return walkEntries(repo, tree, prefix, func(full string, _ object.TreeEntry) {
		*out = append(*out, full)
	})
}

// walkEntries ruft visit für jeden Datei-Eintrag unter tree (regulär,
// ausführbar, Symlink) mit seinem vollen Pfad; die Fehlerpolitik ist die von
// walkTree.
func walkEntries(repo *gogit.Repository, tree *object.Tree, prefix string, visit func(path string, e object.TreeEntry)) error {
	for _, e := range tree.Entries {
		full := joinPath(prefix, e.Name)
		switch {
		case e.Mode == filemode.Dir:
			sub, err := repo.TreeObject(e.Hash)
			if err != nil {
				return fmt.Errorf("nicht lesbarer Unterbaum %q: %w", full, err)
			}
			if err := walkEntries(repo, sub, full, visit); err != nil {
				return err
			}
		case isFileMode(e.Mode):
			visit(full, e)
		}
	}
	return nil
}

// isFileMode meldet die drei Arten, die als Datei zählen: regulär, ausführbar,
// Symlink. Ein Verzeichnis und ein Gitlink zählen nicht.
func isFileMode(m filemode.FileMode) bool {
	return m == filemode.Regular || m == filemode.Executable || m == filemode.Symlink
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
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
		// Grenze der Bibliothek: go-git meldet ein UNLESBARES Blob mit
		// demselben Fehler wie eine im Tree fehlende Datei. Nur der zweite
		// Fall ist harmlos (Datei später angelegt) und darf befundfrei
		// bleiben; der erste ist ein Umgebungsfehler und muss den Lauf
		// abbrechen, sonst verschwindet eine echte Core-Änderung still.
		// Der Tree-Eintrag trennt die beiden: er trägt Name und Hash und
		// braucht das Blob nicht.
		if unreadable, ferr := entryUnreadable(tree, path); unreadable {
			return nil, false, fmt.Errorf("nicht lesbares Objekt zu %q in %q: %w", path, ref, err)
		} else if ferr != nil {
			return nil, false, fmt.Errorf("nicht lesbarer Tree-Eintrag zu %q in %q: %w", path, ref, ferr)
		}
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

// entryUnreadable trennt die zwei Zustände, die go-git beide als
// object.ErrFileNotFound meldet: existiert der Tree-Eintrag zu path, war das
// zugehörige Blob unlesbar (true); existiert er nicht, fehlt die Datei in
// diesem Stand (false, nil). Ein Fehler beim Nachsehen, der nicht "gibt es
// nicht" bedeutet, wird durchgereicht statt als Abwesenheit gelesen.
func entryUnreadable(tree *object.Tree, path string) (bool, error) {
	e, err := tree.FindEntry(path)
	switch {
	case err == nil:
		// Nur ein Eintrag, der einen Datei-Inhalt bezeichnet, kann unlesbar
		// sein. Ein Gitlink (Submodul) und ein Verzeichnis tragen keinen —
		// für sie ist die Abwesenheit eines Blobs der Normalfall und
		// befundfrei richtig, nicht ein Umgebungsfehler.
		return e.Mode == filemode.Regular || e.Mode == filemode.Executable ||
			e.Mode == filemode.Symlink, nil
	case errors.Is(err, object.ErrEntryNotFound), errors.Is(err, object.ErrDirectoryNotFound):
		return false, nil
	default:
		return false, err
	}
}

// TrackedPaths liefert die Menge der im git-Index getrackten Pfade
// (DC-FA-TRK-001.a Schritt 2). Index-Wahrheit: gestagte, noch nie
// committete Dateien sind enthalten; .gitignore wird nicht interpretiert.
// Ein frisches Repo ohne Einträge liefert eine leere Menge (kein Fehler);
// ein unlesbarer Index ⇒ Fehler (fail-closed, Exit 2).
func (a *Adapter) TrackedPaths() (map[string]bool, error) {
	idx, err := a.repo.Storer.Index()
	if err != nil {
		return nil, fmt.Errorf("git-Index nicht lesbar: %w", err)
	}
	set := make(map[string]bool, len(idx.Entries))
	for _, e := range idx.Entries {
		set[e.Name] = true
	}
	return set, nil
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

// LineChanges liefert die Zeilen-Summe der Änderungen zwischen base und head
// über die Pfade paths (DC-FA-MET-001): ohne Rename-Erkennung, über literale
// Pfade, ohne die Pfade aus excluded. Ein nicht auflösbares Ref oder eine
// Basis, die kein Vorfahr von head ist, ⇒ Fehler (fail-closed).
func (a *Adapter) LineChanges(base, head string, paths, excluded []string) (driven.LineStat, error) {
	var stat driven.LineStat
	baseTree, headTree, err := a.rangeTrees(base, head)
	if err != nil {
		return stat, err
	}
	changes, err := object.DiffTree(baseTree, headTree)
	if err != nil {
		return stat, fmt.Errorf("diff %s..%s nicht lesbar: %w", base, head, err)
	}
	for _, ch := range changes {
		name := changeName(ch)
		if !underAny(name, paths) || underAny(name, excluded) {
			continue
		}
		n, uncounted, err := changeLines(ch)
		if err != nil {
			return stat, fmt.Errorf("änderung an %s nicht lesbar: %w", name, err)
		}
		if uncounted {
			stat.Binary = append(stat.Binary, name)
		}
		stat.Lines += n
	}
	sort.Strings(stat.Binary)
	stat.Missing = missingPaths(paths, baseTree, headTree)
	return stat, nil
}

// rangeTrees löst base und head auf, verlangt base als Vorfahr von head
// (einschließlich Gleichheit) und liefert beide Trees.
func (a *Adapter) rangeTrees(base, head string) (*object.Tree, *object.Tree, error) {
	baseHash, err := a.repo.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		return nil, nil, fmt.Errorf("basis %q nicht auflösbar: %w", base, err)
	}
	headHash, err := a.repo.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		return nil, nil, fmt.Errorf("spitze %q nicht auflösbar: %w", head, err)
	}
	ok, err := a.isAncestor(*baseHash, *headHash)
	if err != nil {
		return nil, nil, fmt.Errorf("historie von %q nicht lesbar (flacher klon ohne die basis?): %w", head, err)
	}
	if !ok {
		return nil, nil, fmt.Errorf("basis %q ist kein vorfahr von %q", base, head)
	}
	baseTree, err := a.treeAt(base)
	if err != nil {
		return nil, nil, fmt.Errorf("tree von %q nicht lesbar: %w", base, err)
	}
	headTree, err := a.treeAt(head)
	if err != nil {
		return nil, nil, fmt.Errorf("tree von %q nicht lesbar: %w", head, err)
	}
	return baseTree, headTree, nil
}

// changeName ist der Pfad einer Änderung: der Zielpfad, bei einer Löschung der
// Ausgangspfad.
func changeName(ch *object.Change) string {
	if ch.To.Name != "" {
		return ch.To.Name
	}
	return ch.From.Name
}

// changeLines zählt Einfügungen plus Löschungen einer Änderung; uncounted
// meldet eine Datei, die nicht zählbar ist (Binärdatei mit NUL in den ersten
// 8000 Bytes auf einer der beiden Seiten wie bei git, oder ein Submodul) und
// 0 beiträgt. Eine hinzugefügte oder gelöschte leere Datei ist zählbar.
func changeLines(ch *object.Change) (lines int, uncounted bool, err error) {
	if isSubmodule(ch.From.TreeEntry) || isSubmodule(ch.To.TreeEntry) {
		return 0, true, nil
	}
	from, to, err := ch.Files()
	if err != nil {
		return 0, false, err
	}
	for _, f := range []*object.File{from, to} {
		if f == nil {
			continue
		}
		bin, berr := f.IsBinary()
		if berr != nil {
			return 0, false, berr
		}
		if bin {
			return 0, true, nil
		}
	}
	patch, err := ch.Patch()
	if err != nil {
		return 0, false, err
	}
	for _, st := range patch.Stats() {
		lines += st.Addition + st.Deletion
	}
	return lines, false, nil
}

func isSubmodule(e object.TreeEntry) bool {
	return e.Mode == filemode.Submodule
}

// underAny meldet, ob name unter einem der literalen Pfade liegt: gleich oder
// ein Eintrag darunter (ein Verzeichnis zählt rekursiv).
func underAny(name string, paths []string) bool {
	for _, p := range paths {
		p = strings.TrimSuffix(p, "/")
		if name == p || strings.HasPrefix(name, p+"/") {
			return true
		}
	}
	return false
}

// missingPaths liefert die Pfade, die weder im Basis- noch im Spitzen-Tree als
// Datei oder Verzeichnis existieren, sortiert.
func missingPaths(paths []string, trees ...*object.Tree) []string {
	var out []string
	for _, p := range paths {
		clean := strings.TrimSuffix(p, "/")
		found := false
		for _, t := range trees {
			if _, err := t.FindEntry(clean); err == nil {
				found = true
				break
			}
		}
		if !found {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

// isAncestor meldet, ob base von head aus erreichbar ist (einschließlich
// Gleichheit); die Suche endet beim Fund. Ein nicht lesbarer Commit auf dem Weg
// ist ein Fehler.
func (a *Adapter) isAncestor(base, head plumbing.Hash) (bool, error) {
	seen := map[plumbing.Hash]bool{}
	stack := []plumbing.Hash{head}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur == base {
			return true, nil
		}
		if seen[cur] {
			continue
		}
		seen[cur] = true
		c, err := a.repo.CommitObject(cur)
		if err != nil {
			return false, err
		}
		stack = append(stack, c.ParentHashes...)
	}
	return false, nil
}

// Renames liefert die Zuordnung Quelle → Ende der Umbenennungen, die an
// Link-Zielen teilnehmen (driven.VCS.Renames, DC-FA-VCS-001.a Schritt 2,
// ADR-0115). Eine Umbenennung ist ein Schritt, in dem ein Pfad verschwindet und
// ein anderer mit gleichem Eintrag erscheint; die Bibliotheks-Erkennung
// bleibt ungenutzt.
func (a *Adapter) Renames(base, head string) (map[string]string, error) {
	if head == driven.IndexRef {
		return a.stagedRenames(base)
	}
	return a.rangeRenames(base, head)
}

// fileEntry ist, was die Erkennung vergleicht: Blob-Kennung und Art (Datei
// oder Symlink); das Ausführbar-Bit zählt nicht.
type fileEntry struct {
	hash    plumbing.Hash
	symlink bool
}

func newFileEntry(e object.TreeEntry) fileEntry {
	return fileEntry{hash: e.Hash, symlink: e.Mode == filemode.Symlink}
}

// stepEvents sind die Pfad-Ereignisse eines Schritts: ein verschwundener Pfad
// trägt sein Ende (leer: ohne Ende, gelöscht), ein neuer Pfad steht in created.
// Ein Pfad, der vor und nach dem Schritt als Datei besteht, hat kein Ereignis.
type stepEvents struct {
	gone    map[string]string
	created map[string]bool
}

// eventsOf bildet die Ereignisse aus den verschwundenen und den neuen Pfaden
// eines Schritts: A → B nur, wenn genau ein verschwundener und genau ein neuer
// Pfad den Eintrag tragen.
func eventsOf(gone, added map[string]fileEntry) stepEvents {
	goneBy, addedBy := groupByEntry(gone), groupByEntry(added)
	ev := stepEvents{gone: make(map[string]string, len(gone)), created: make(map[string]bool, len(added))}
	for p, e := range gone {
		ev.gone[p] = ""
		if g, n := goneBy[e], addedBy[e]; len(g) == 1 && len(n) == 1 {
			ev.gone[p] = n[0]
		}
	}
	for q := range added {
		ev.created[q] = true
	}
	return ev
}

func groupByEntry(paths map[string]fileEntry) map[fileEntry][]string {
	out := make(map[fileEntry][]string, len(paths))
	for p, e := range paths {
		out[e] = append(out[e], p)
	}
	return out
}

// lineSet fasst zusammen, wo eine Quelle auf allen Abstammungslinien steht, die
// einen Commit erreichen: unchanged — eine Linie lässt sie in jedem Schritt
// bestehen, deleted — eine Linie lässt sie ohne Ende verschwinden, at — die
// Pfade, an denen Linien stehen, die sie umbenannt haben.
type lineSet struct {
	unchanged bool
	deleted   bool
	at        map[string]bool
}

func (ls lineSet) with(p string) lineSet {
	at := make(map[string]bool, len(ls.at)+1)
	for k := range ls.at {
		at[k] = true
	}
	at[p] = true
	ls.at = at
	return ls
}

func (ls lineSet) union(o lineSet) lineSet {
	ls.unchanged = ls.unchanged || o.unchanged
	ls.deleted = ls.deleted || o.deleted
	for p := range o.at {
		ls = ls.with(p)
	}
	return ls
}

// land setzt eine Linie, die vor dem Schritt am Pfad p steht, hinter den
// Schritt; unchanged sagt, dass sie die Quelle noch nie verlassen hat. Ein Pfad,
// den der Schritt neu anlegt, kann vorher keine Quelle getragen haben: die
// Linie fällt weg.
func (ls lineSet) land(p string, unchanged bool, ev stepEvents) lineSet {
	if ev.created[p] {
		return ls
	}
	to, moved := ev.gone[p]
	switch {
	case !moved && unchanged:
		ls.unchanged = true
	case !moved:
		ls = ls.with(p)
	case to == "":
		ls.deleted = true
	default:
		ls = ls.with(to)
	}
	return ls
}

func (ls lineSet) step(source string, ev stepEvents) lineSet {
	out := lineSet{deleted: ls.deleted}
	if ls.unchanged {
		out = out.land(source, true, ev)
	}
	for p := range ls.at {
		out = out.land(p, false, ev)
	}
	return out
}

// lineStates ordnet jeder Quelle ihr lineSet an einem Commit zu; eine Quelle
// ohne Eintrag steht auf jeder Linie unverändert.
type lineStates map[string]lineSet

func (s lineStates) get(source string) lineSet {
	if ls, ok := s[source]; ok {
		return ls
	}
	return lineSet{unchanged: true}
}

// advance führt die Zustände über einen Schritt fort. Neu in die Zuordnung
// kommen nur Quellen, auf deren Pfad der Schritt ein Ereignis hat.
func advance(parent lineStates, ev stepEvents) lineStates {
	if len(ev.gone) == 0 && len(ev.created) == 0 {
		return parent
	}
	out := make(lineStates, len(parent))
	for source, ls := range parent {
		out[source] = ls.step(source, ev)
	}
	touch := func(p string) {
		if _, ok := parent[p]; !ok {
			out[p] = parent.get(p).step(p, ev)
		}
	}
	for p := range ev.gone {
		touch(p)
	}
	for p := range ev.created {
		touch(p)
	}
	return out
}

// mergeStates vereint die Zustände der Eltern-Commits eines Merges; ein
// Merge-Commit trägt selbst keine Umbenennung bei.
func mergeStates(parents []lineStates) lineStates {
	sources := map[string]bool{}
	for _, p := range parents {
		for source := range p {
			sources[source] = true
		}
	}
	out := make(lineStates, len(sources))
	for source := range sources {
		var acc lineSet
		for _, p := range parents {
			acc = acc.union(p.get(source))
		}
		out[source] = acc
	}
	return out
}

// finalRenames filtert die Zustände an HEAD auf die teilnehmenden
// Umbenennungen (ADR-0115 Entscheidung 9): die Quelle liegt in BASE, genau ein
// Ende und keine Linie ohne Ende, das Ende liegt in HEAD und nicht in BASE und
// ist Ende keiner zweiten Quelle.
func finalRenames(states lineStates, baseSet, headSet map[string]bool) map[string]string {
	ends := map[string]string{}
	count := map[string]int{}
	for source, ls := range states {
		if !baseSet[source] || ls.deleted || len(ls.at) != 1 {
			continue
		}
		for end := range ls.at {
			if headSet[end] && !baseSet[end] {
				ends[source] = end
				count[end]++
			}
		}
	}
	out := make(map[string]string, len(ends))
	for source, end := range ends {
		if count[end] == 1 {
			out[source] = end
		}
	}
	return out
}

func entryIndex(t *object.Tree) map[string]object.TreeEntry {
	out := make(map[string]object.TreeEntry, len(t.Entries))
	for _, e := range t.Entries {
		out[e.Name] = e
	}
	return out
}

// diffFiles sammelt die Datei-Pfade, die zwischen left und right verschwinden
// (gone) oder entstehen (added). Ein Teilbaum gleicher Kennung auf beiden
// Seiten bleibt ungelesen; jeder andere Unterbaum wird gelesen, und ein
// unlesbarer ist ein Fehler — kein Diff-Walker, der ihn als Ende der Aufzählung
// liest (DC-FA-VCS-001.a Schritt 2).
func diffFiles(repo *gogit.Repository, left, right *object.Tree, prefix string, gone, added map[string]fileEntry) error {
	l, r := entryIndex(left), entryIndex(right)
	for _, name := range entryNames(l, r) {
		le, lok := l[name]
		re, rok := r[name]
		pair := entryPair{left: le, right: re, hasLeft: lok, hasRight: rok}
		if err := diffEntry(repo, pair, joinPath(prefix, name), gone, added); err != nil {
			return err
		}
	}
	return nil
}

// entryNames liefert die Namen beider Seiten, sortiert und ohne Doppelte.
func entryNames(l, r map[string]object.TreeEntry) []string {
	names := make([]string, 0, len(l)+len(r))
	for n := range l {
		names = append(names, n)
	}
	for n := range r {
		if _, ok := l[n]; !ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}

// entryPair sind die Einträge eines Namens auf beiden Seiten eines Schritts.
type entryPair struct {
	left, right       object.TreeEntry
	hasLeft, hasRight bool
}

func diffEntry(repo *gogit.Repository, p entryPair, full string, gone, added map[string]fileEntry) error {
	if p.hasLeft && p.hasRight && p.left.Mode == filemode.Dir && p.right.Mode == filemode.Dir {
		if p.left.Hash == p.right.Hash {
			return nil
		}
		return diffSubtrees(repo, p.left.Hash, p.right.Hash, full, gone, added)
	}
	if p.hasLeft {
		if err := collectSide(repo, p.left, full, p.hasRight && isFileMode(p.right.Mode), gone); err != nil {
			return err
		}
	}
	if p.hasRight {
		return collectSide(repo, p.right, full, p.hasLeft && isFileMode(p.left.Mode), added)
	}
	return nil
}

func diffSubtrees(repo *gogit.Repository, left, right plumbing.Hash, full string, gone, added map[string]fileEntry) error {
	lt, err := repo.TreeObject(left)
	if err != nil {
		return fmt.Errorf("nicht lesbarer Unterbaum %q: %w", full, err)
	}
	rt, err := repo.TreeObject(right)
	if err != nil {
		return fmt.Errorf("nicht lesbarer Unterbaum %q: %w", full, err)
	}
	return diffFiles(repo, lt, rt, full, gone, added)
}

// collectSide trägt die Datei-Pfade eines Eintrags in out ein. otherIsFile sagt,
// dass die Gegenseite am selben Pfad eine Datei trägt: dann besteht der Pfad
// fort, auch mit geändertem Inhalt oder geänderter Art.
func collectSide(repo *gogit.Repository, e object.TreeEntry, full string, otherIsFile bool, out map[string]fileEntry) error {
	switch {
	case e.Mode == filemode.Dir:
		sub, err := repo.TreeObject(e.Hash)
		if err != nil {
			return fmt.Errorf("nicht lesbarer Unterbaum %q: %w", full, err)
		}
		return walkEntries(repo, sub, full, func(p string, fe object.TreeEntry) {
			out[p] = newFileEntry(fe)
		})
	case isFileMode(e.Mode) && !otherIsFile:
		out[full] = newFileEntry(e)
	}
	return nil
}

// rangeWalk hält die Commits der Range (von HEAD erreichbar, von BASE nicht),
// Eltern vor Kindern.
type rangeWalk struct {
	repo     *gogit.Repository
	excluded map[plumbing.Hash]bool
	commits  map[plumbing.Hash]*object.Commit
	order    []plumbing.Hash
}

func (w *rangeWalk) visit(h plumbing.Hash) error {
	if w.excluded[h] || w.commits[h] != nil {
		return nil
	}
	c, err := w.repo.CommitObject(h)
	if err != nil {
		return fmt.Errorf("commit %s nicht lesbar: %w", h, err)
	}
	w.commits[h] = c
	for _, p := range c.ParentHashes {
		if err := w.visit(p); err != nil {
			return err
		}
	}
	w.order = append(w.order, h)
	return nil
}

// step liest den Schritt eines Commits mit genau einem Eltern-Commit: die
// Bäume beider, auch wenn der Eltern-Commit vor BASE liegt.
func (w *rangeWalk) step(c *object.Commit) (stepEvents, error) {
	parent, err := w.repo.CommitObject(c.ParentHashes[0])
	if err != nil {
		return stepEvents{}, fmt.Errorf("commit %s nicht lesbar: %w", c.ParentHashes[0], err)
	}
	pt, err := parent.Tree()
	if err != nil {
		return stepEvents{}, fmt.Errorf("baum von commit %s nicht lesbar: %w", parent.Hash, err)
	}
	ct, err := c.Tree()
	if err != nil {
		return stepEvents{}, fmt.Errorf("baum von commit %s nicht lesbar: %w", c.Hash, err)
	}
	gone, added := map[string]fileEntry{}, map[string]fileEntry{}
	if err := diffFiles(w.repo, pt, ct, "", gone, added); err != nil {
		return stepEvents{}, fmt.Errorf("schritt %s nicht lesbar: %w", c.Hash, err)
	}
	return eventsOf(gone, added), nil
}

func (w *rangeWalk) statesOf(h plumbing.Hash, states map[plumbing.Hash]lineStates) (lineStates, error) {
	c := w.commits[h]
	switch len(c.ParentHashes) {
	case 0:
		return nil, nil
	case 1:
		ev, err := w.step(c)
		if err != nil {
			return nil, err
		}
		return advance(states[c.ParentHashes[0]], ev), nil
	}
	parents := make([]lineStates, len(c.ParentHashes))
	for i, p := range c.ParentHashes {
		parents[i] = states[p]
	}
	return mergeStates(parents), nil
}

func (a *Adapter) rangeRenames(base, head string) (map[string]string, error) {
	baseHash, err := a.repo.ResolveRevision(plumbing.Revision(base))
	if err != nil {
		return nil, fmt.Errorf("Range-Basis %q nicht auflösbar: %w", base, err)
	}
	headHash, err := a.repo.ResolveRevision(plumbing.Revision(head))
	if err != nil {
		return nil, fmt.Errorf("Range-Spitze %q nicht auflösbar: %w", head, err)
	}
	isAnc, err := a.isAncestor(*baseHash, *headHash)
	if err != nil {
		return nil, fmt.Errorf("historie von %q nicht lesbar: %w", head, err)
	}
	if !isAnc {
		return map[string]string{}, nil
	}
	excluded, err := a.ancestors(*baseHash)
	if err != nil {
		return nil, fmt.Errorf("Range-Basis-Vorfahren nicht lesbar: %w", err)
	}
	w := &rangeWalk{repo: a.repo, excluded: excluded, commits: map[plumbing.Hash]*object.Commit{}}
	if err := w.visit(*headHash); err != nil {
		return nil, err
	}
	states := make(map[plumbing.Hash]lineStates, len(w.order))
	for _, h := range w.order {
		st, err := w.statesOf(h, states)
		if err != nil {
			return nil, err
		}
		states[h] = st
	}
	baseSet, err := a.pathSet(base)
	if err != nil {
		return nil, err
	}
	headSet, err := a.pathSet(head)
	if err != nil {
		return nil, err
	}
	return finalRenames(states[*headHash], baseSet, headSet), nil
}

func (a *Adapter) pathSet(ref string) (map[string]bool, error) {
	paths, err := a.pathsAt(ref)
	if err != nil {
		return nil, err
	}
	return keySet(paths), nil
}

func keySet(paths []string) map[string]bool {
	out := make(map[string]bool, len(paths))
	for _, p := range paths {
		out[p] = true
	}
	return out
}

// stagedRenames liest den einen Schritt HEAD → Index; „BASE“ der Bedingungen ist
// HEAD. Ohne HEAD (erster Commit) ist die Zuordnung leer.
func (a *Adapter) stagedRenames(base string) (map[string]string, error) {
	if !a.hasHead() {
		return map[string]string{}, nil
	}
	tree, err := a.treeAt(base)
	if err != nil {
		return nil, fmt.Errorf("staged-Basis %q nicht auflösbar: %w", base, err)
	}
	left := map[string]fileEntry{}
	if err := walkEntries(a.repo, tree, "", func(p string, e object.TreeEntry) {
		left[p] = newFileEntry(e)
	}); err != nil {
		return nil, fmt.Errorf("nicht vollständig lesbarer Tree zu %q: %w", base, err)
	}
	idx, err := a.repo.Storer.Index()
	if err != nil {
		return nil, fmt.Errorf("git-Index nicht lesbar: %w", err)
	}
	right := map[string]fileEntry{}
	for _, e := range idx.Entries {
		if isFileMode(e.Mode) {
			right[e.Name] = fileEntry{hash: e.Hash, symlink: e.Mode == filemode.Symlink}
		}
	}
	gone, added := map[string]fileEntry{}, map[string]fileEntry{}
	for p, e := range left {
		if _, ok := right[p]; !ok {
			gone[p] = e
		}
	}
	for p, e := range right {
		if _, ok := left[p]; !ok {
			added[p] = e
		}
	}
	states := advance(nil, eventsOf(gone, added))
	return finalRenames(states, pathKeys(left), pathKeys(right)), nil
}

func pathKeys(m map[string]fileEntry) map[string]bool {
	out := make(map[string]bool, len(m))
	for p := range m {
		out[p] = true
	}
	return out
}
