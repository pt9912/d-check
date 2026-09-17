package git

import (
	"os"
	"regexp"
	"strings"

	billy "github.com/go-git/go-billy/v5"
)

const (
	packDir    = "objects/pack"
	packPrefix = "pack-"
	packExt    = ".pack"
	idxExt     = ".idx"
)

// packHashSuffixRE erkennt einen SHA1- (40) oder SHA256-Hex-Hash (64) am
// Ende eines Pack-Dateistamms, optional durch "-"/"_" vom Rest getrennt —
// git benennt einen Pack nach dem Hash seines Inhalts, unabhängig vom
// Präfix, den ein Werkzeug davor setzt.
var packHashSuffixRE = regexp.MustCompile(`(?:^|[-_])([0-9a-f]{40}|[0-9a-f]{64})$`)

// packAliasFS macht Packs unter objects/pack unter ihrem kanonischen Namen
// pack-<hash>.{pack,idx} sichtbar, unabhängig vom tatsächlichen Präfix auf
// der Platte (z. B. loose-<hash>.pack, wie `git maintenance
// run --task=loose-objects` schreibt). go-gits DotGit-Schicht entdeckt und
// öffnet Packs ausschließlich über diesen kanonischen Namen
// (storage/filesystem/dotgit.ObjectPacks/objectPackPath, go-git v5.19.2) —
// ein Pack unter anderem Präfix ist für sie unsichtbar, obwohl das Objekt
// gültig ist und `git` selbst es anstandslos liest (eingehender CR,
// ai-harness-init, 2026-09-17). Rein lesend: es wird nie geschrieben, das
// gemountete Repository bleibt unverändert (DC-QA-03). Jede Methode außer
// Open/ReadDir wird über das eingebettete Interface unverändert
// durchgereicht.
type packAliasFS struct {
	billy.Filesystem
}

func newPackAliasFS(fs billy.Filesystem) billy.Filesystem {
	return &packAliasFS{Filesystem: fs}
}

// packAliases scannt objects/pack und bildet für jeden nicht-kanonisch
// benannten Pack mit gültigem Hash-Suffix UND vorhandener .idx-Datei den
// kanonischen Namen auf den realen ab. Ein bereits kanonisch benannter
// Eintrag — auch mit ungültigem oder fremdem Hash — wird nie überschrieben:
// die bestehende Fehldiagnose für einen wirklich kaputt benannten
// kanonischen Pack bleibt unverändert bestehen.
func (fs *packAliasFS) packAliases() map[string]string {
	entries, err := fs.Filesystem.ReadDir(packDir)
	if err != nil {
		return nil
	}
	onDisk := make(map[string]bool, len(entries))
	for _, e := range entries {
		onDisk[e.Name()] = true
	}
	aliases := make(map[string]string)
	for name := range onDisk {
		if strings.HasPrefix(name, packPrefix) || !strings.HasSuffix(name, packExt) {
			continue
		}
		stem := strings.TrimSuffix(name, packExt)
		m := packHashSuffixRE.FindStringSubmatch(stem)
		if m == nil {
			continue
		}
		idxName := stem + idxExt
		if !onDisk[idxName] {
			continue
		}
		canonicalPack := packPrefix + m[1] + packExt
		canonicalIdx := packPrefix + m[1] + idxExt
		if onDisk[canonicalPack] || onDisk[canonicalIdx] {
			continue
		}
		aliases[canonicalPack] = name
		aliases[canonicalIdx] = idxName
	}
	return aliases
}

// aliasFileInfo überschreibt Name(), damit ein virtueller Verzeichnis-
// Eintrag den kanonischen statt des realen Namens trägt.
type aliasFileInfo struct {
	os.FileInfo
	name string
}

func (fi aliasFileInfo) Name() string { return fi.name }

// ReadDir spiegelt jeden aliasfähigen Pack unter seinem kanonischen Namen
// zusätzlich in die Liste — die realen Einträge bleiben erhalten, DotGits
// eigener Präfix-Filter ignoriert sie ohnehin.
func (fs *packAliasFS) ReadDir(p string) ([]os.FileInfo, error) {
	entries, err := fs.Filesystem.ReadDir(p)
	if err != nil || strings.TrimSuffix(p, "/") != packDir {
		return entries, err
	}
	aliases := fs.packAliases()
	if len(aliases) == 0 {
		return entries, nil
	}
	byName := make(map[string]os.FileInfo, len(entries))
	for _, e := range entries {
		byName[e.Name()] = e
	}
	out := append([]os.FileInfo(nil), entries...)
	for canonical, realName := range aliases {
		if fi, ok := byName[realName]; ok {
			out = append(out, aliasFileInfo{FileInfo: fi, name: canonical})
		}
	}
	return out, nil
}

// Open löst einen kanonischen Pack-Pfad auf den realen Dateinamen auf, falls
// nötig — jeder andere Pfad (HEAD, Refs, lose Objekte, ein bereits kanonisch
// benannter Pack) geht unverändert an das eingebettete Filesystem.
func (fs *packAliasFS) Open(filename string) (billy.File, error) {
	dir, base := splitPackPath(filename)
	if dir == packDir && strings.HasPrefix(base, packPrefix) {
		if realName, ok := fs.packAliases()[base]; ok {
			return fs.Filesystem.Open(fs.Join(packDir, realName))
		}
	}
	return fs.Filesystem.Open(filename)
}

// splitPackPath zerlegt einen '/'-getrennten .git-relativen Pfad in
// Verzeichnis und Basisname — ohne path/filepath, weil DotGit seine eigenen
// Pfade unabhängig vom Host-Trennzeichen immer mit '/' zusammensetzt.
func splitPackPath(filename string) (dir, base string) {
	i := strings.LastIndex(filename, "/")
	if i < 0 {
		return "", filename
	}
	return filename[:i], filename[i+1:]
}
