// Package driven trägt die vom Hexagon definierten driven Ports
// (spec/architecture.md §2; Ordnerkonvention gemäß ADR-0005,
// u-boot-Vorbild). Adapter implementieren diese Interfaces; der Kern
// kennt nur sie.
package driven

// EntryKind klassifiziert einen Pfad aus Lstat-Sicht (keine
// Symlink-Auflösung) — Grundlage der Symlink-Ablehnung
// (DC-FA-LINK-002).
type EntryKind int

// Lstat-Klassifikationen eines Pfads.
const (
	KindMissing EntryKind = iota
	KindFile
	KindDir
	KindSymlink
)

// DirEntry ist ein Verzeichniseintrag des Filesystem-Ports.
type DirEntry struct {
	Name string
	Kind EntryKind
}

// Filesystem ist der driven Port für alle Dateisystem-Zugriffe; alle
// Pfade sind '/'-getrennt und relativ zur Repo-Wurzel. Der
// Filesystem-Adapter ist die einzige Dateisystem-Tür von d-check.
type Filesystem interface {
	// Kind liefert die Lstat-Klassifikation des Pfads.
	Kind(rel string) (EntryKind, error)
	// ReadFile liest den Datei-Inhalt.
	ReadFile(rel string) ([]byte, error)
	// List liefert die Einträge eines Verzeichnisses (nicht rekursiv).
	List(relDir string) ([]DirEntry, error)
}
