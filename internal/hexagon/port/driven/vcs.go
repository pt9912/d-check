package driven

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
	// AllPaths liefert für die Range base..head (head == IndexRef ⇒ staged:
	// base-Tree gegen den Index) JEDEN Datei-Pfad ('/'-getrennt, repo-relativ)
	// an BEIDEN Enden — direkt über die Baum-Struktur aufgelöst, nicht aus
	// einem Diff hergeleitet (DC-FA-VCS-001.a Schritt 2). Ein Unterbaum, dessen
	// Objekt nicht lesbar ist, ist ein Umgebungsfehler und wird zum Fehler
	// (fail-closed, Exit 2) — nicht zu einer Lücke, die den Rest des Trees
	// wortlos verschluckt (CO-001, dritte Ausprägung). Eine Ausnahme bleibt:
	// staged VOR dem ersten Commit hat nichts, das immutabel sein könnte —
	// beide Mengen sind dann leer, kein Fehler (Parität zum abgelösten
	// ChangedPaths).
	AllPaths(base, head string) (baseAll, headAll []string, err error)
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
