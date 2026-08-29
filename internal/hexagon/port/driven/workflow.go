package driven

// Der Workflow-Port (DC-FA-WF-001): das Modul `workflows` braucht die
// STRUKTUR einer Workflow-Datei — Referenzen mit ihrer Zeile und der
// Rechte-Kontext ihres Jobs. Das Parsen selbst ist Adapter-Sache: der Kern
// trägt keine Fremd-Bibliothek (Hexagon-Schnitt, ADR-0004/ADR-0012; die
// yaml-Abhängigkeit lebt in den Adaptern, ADR-0009).
//
// Die Typen hier sind PORT-Typen wie DirEntry und VCSStatus, nicht
// Modell-Typen: der Port beschreibt, was von außen hereinkommt, und bleibt
// dadurch frei von der Kern-Konfiguration.

// WorkflowPerms ist ein `permissions:`-Block. Declared unterscheidet
// "kein permissions" von "leeres permissions" — nur die erste Form erbt
// still vom Workflow-Kopf, und genau das trägt den undeclared-Befund.
type WorkflowPerms struct {
	// Scopes bildet Scope-Namen auf "none"/"read"/"write" ab.
	Scopes map[string]string
	// All ist "read" bzw. "write" bei `read-all`/`write-all`, sonst leer;
	// es setzt JEDEN Scope, auch einen, den Scopes nicht nennt.
	All string
	// Declared sagt, ob der Block überhaupt dasteht.
	Declared bool
}

// WorkflowUses ist eine `uses:`-Referenz mit ihrer Fundstelle.
type WorkflowUses struct {
	// Value ist der Skalar OHNE Kommentar; Comment ist der YAML-LineComment.
	// Beide getrennt, weil `uses: x@sha # v1` in YAML zwei Dinge sind und
	// eine Textsuche das nicht unterscheidet.
	Value   string
	Comment string
	Line    int
	// JobLevel: true für den aufgerufenen Workflow (Job-Ebene), false für
	// eine Action (Step-Ebene). Nur der erste kennt eine Rechte-Weitergabe.
	JobLevel bool
	// JobPerms sind die Rechte des Jobs, in dem die Referenz steht.
	JobPerms WorkflowPerms
}

// WorkflowDoc ist die geparste Struktur einer Workflow-Datei.
type WorkflowDoc struct {
	// HeadPerms ist der `permissions:`-Block des Workflow-Kopfes.
	HeadPerms WorkflowPerms
	// JobPerms sind die Blöcke aller Jobs — für die Frage, was eine Datei
	// als Ziel VERLANGT (Vereinigung aus Kopf und Jobs).
	JobPerms []WorkflowPerms
	// Uses sind alle gefundenen Referenzen in Dokument-Reihenfolge.
	Uses []WorkflowUses
}

// WorkflowParser übersetzt den Inhalt einer Workflow-Datei in die Struktur.
// Ein Fehler bedeutet: nicht parsbar — das Modul macht daraus einen Befund,
// keinen Übersprung.
type WorkflowParser interface {
	ParseWorkflow(content []byte) (WorkflowDoc, error)
}
