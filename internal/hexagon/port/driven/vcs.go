package driven

// VCSStatus klassifiziert eine Änderung im git-Diff einer Commit-Range
// (DC-FA-VCS-001.a Schritt 2). Renames werden pfadbasiert als Löschung des
// alten Pfads + Hinzufügung des neuen erkannt — eine immutable Datei behält
// ihren Pfad, ihr Verschwinden ist der Befund (kein Inhalts-Ähnlichkeits-
// Matching, spec/lastenheft.md §DC-FA-VCS-001 Out-of-Scope).
type VCSStatus rune

// Diff-Status-Klassen einer Commit-Range bzw. eines staged Index-Diffs.
const (
	VCSAdded    VCSStatus = 'A'
	VCSModified VCSStatus = 'M'
	VCSDeleted  VCSStatus = 'D'
)

// VCSChange ist ein geänderter Pfad zwischen zwei git-Refs.
type VCSChange struct {
	Status VCSStatus
	// Path ist der betroffene Pfad ('/'-getrennt, relativ zur Repo-Wurzel):
	// bei Added/Modified der Ziel-Pfad, bei Deleted der entfallene.
	Path string
}

// IndexRef ist der Sentinel-Ref für den staged Index (--staged-Modus,
// DC-FA-VCS-001.a Schritt 1): BASE = HEAD, HEAD = der staged Index.
const IndexRef = ":index:"

// CommitMeta ist ein Commit mit Kurz-SHA und roher (unbereinigter) Message —
// die Eingabe des Moduls commits (DC-FA-COMMITS-001). Die git-`strip`-Bereinigung
// und die Kennungs-Prüfung leistet der Kern, nicht der Adapter.
type CommitMeta struct {
	ShortSHA string
	Message  string
}

// VCS ist der driven Port des Moduls vcs (DC-FA-VCS-001): liest die
// git-Historie **rein lesend** aus dem read-only `.git` — ohne externes
// git-Binary, ohne Netz (spec/architecture.md §2). Der git-Adapter ist die
// einzige git-Tür von d-check; die Eingabe ist gegenüber den hermetischen
// Modulen um `.git` + Range erweitert, bleibt aber lokal/deterministisch
// (DC-QA-02/DC-QA-03). Ein nicht auflösbares Ref/`.git` ⇒ Fehler (fail-closed,
// Exit 2).
type VCS interface {
	// ChangedPaths liefert die Änderungen zwischen base und head (head ==
	// IndexRef ⇒ staged-Diff gegen die base-Version).
	ChangedPaths(base, head string) ([]VCSChange, error)
	// FileAt liefert den Inhalt von path an ref; ok=false, wenn an ref abwesend.
	FileAt(ref, path string) (content []byte, ok bool, err error)
	// CommitMessages liefert die rohen Messages der **Nicht-Merge**-Commits der
	// Range base..head (git rev-list --no-merges-Parität) in deterministischer
	// Reihenfolge (nach Commit-SHA) — die Eingabe des Moduls commits
	// (DC-FA-COMMITS-001). Eine gültige, leere Range ist kein Fehler (leere
	// Liste); ein nicht auflösbares Ref ⇒ Fehler (fail-closed, Exit 2). head ==
	// IndexRef wird nicht unterstützt (die Pending-Message existiert nicht als
	// Commit — dafür der --commit-msg-Kurzschluss-Modus) ⇒ Fehler.
	CommitMessages(base, head string) ([]CommitMeta, error)
	// TrackedPaths liefert die Menge der im git-Index getrackten Pfade
	// ('/'-getrennt, repo-relativ) — die Eingabe des Moduls tracked
	// (DC-FA-TRK-001). Der Index ist die Wahrheit (keine .gitignore-
	// Interpretation): eine gestagte, noch nie committete Datei ist
	// enthalten. Ein unlesbarer Index ⇒ Fehler (fail-closed, Exit 2);
	// ein leerer Index (frisches Repo) ist kein Fehler (leere Menge).
	TrackedPaths() (map[string]bool, error)
}
