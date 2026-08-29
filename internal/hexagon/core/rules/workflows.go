package rules

// Modul workflows (DC-FA-WF-001, §DC-FA-WF-001.a): Deklarations-Konsistenz der
// `uses:`-Referenzen von CI-Workflows. Zwei Zusagen an dieselbe Zeile -- die
// FREMDE Referenz ist unbeweglich gepinnt, die LOKALE loest auf und bekommt vom
// aufrufenden Job die Rechte, die sie verlangt.
//
// DAS PARSEN LIEGT IM ADAPTER (driven.WorkflowParser): der Kern traegt keine
// Fremd-Bibliothek (Hexagon-Schnitt, ADR-0004/ADR-0012). Hier steht nur die
// Regel -- die Struktur kommt ueber den Port herein.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): das Modul SCANNT die Dateien in
// workflows.dir und LIEST darueber hinaus die Ziele lokaler Referenzen, auch
// ausserhalb dieses Verzeichnisses. Dieselbe Parse-Zusage gilt dort. Was es
// NICHT deckt, ist jede weitere Actions-Semantik: ein gruener Lauf sagt "diese
// Deklarations-Klasse liegt nicht vor", nicht "der Workflow laeuft".

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls workflows (spec/spezifikation.md §4, SPEC-071..076).
const (
	ReasonUsesPinMissing       = "uses-pin-missing"
	ReasonUsesPinUntagged      = "uses-pin-untagged"
	ReasonUsesLocalMissing     = "uses-local-missing"
	ReasonUsesLocalPermsUndecl = "uses-local-perms-undeclared"
	ReasonUsesLocalPermsNarrow = "uses-local-perms-narrow"
	ReasonWorkflowUnparsable   = "workflow-unparsable"
)

// pinnedRe erkennt einen vollen 40-stelligen Commit-SHA hinter dem `@`.
var pinnedRe = regexp.MustCompile(`@[0-9a-f]{40}(\s|$)`)

// hasTagComment sagt, ob hinter dem SHA ein Kommentar mit Inhalt steht. Der
// Tag-Kommentar ist KEIN Teil des Werts -- YAML fuehrt ihn getrennt.
func hasTagComment(c string) bool {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(c), "#")) != ""
}

// permLevel ordnet die drei Stufen; ein nicht genannter Scope ist `none`.
// Das ist die Semantik des CI-Systems, kein Zugestaendnis.
func permLevel(v string) int {
	switch strings.TrimSpace(v) {
	case "write":
		return 2
	case "read":
		return 1
	default:
		return 0
	}
}

func levelName(l int) string {
	switch l {
	case 2:
		return "write"
	case 1:
		return "read"
	default:
		return "none"
	}
}

// permLevelOf liefert die Stufe eines Scopes in einem Block: der explizite
// Wert, sonst die All-Stufe eines read-all/write-all, sonst `none`.
func permLevelOf(p driven.WorkflowPerms, scope string) int {
	if v, ok := p.Scopes[scope]; ok {
		return permLevel(v)
	}
	return permLevel(p.All)
}

// permDemands liefert die Scopes, die ein Block ueber `none` hinaus verlangt --
// aufsteigend sortiert, weil eine Map-Iteration die Befund-Reihenfolge sonst
// wuerfelt (DC-QA-02).
func permDemands(p driven.WorkflowPerms) []string {
	var out []string
	for k, v := range p.Scopes {
		if permLevel(v) > 0 {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// wfDemands ist, was eine Zieldatei VERLANGT: die Vereinigung aus Kopf und
// allen Jobs. Die Vereinigung ist die sichere Seite -- welcher Job laeuft,
// entscheidet die Laufzeit.
func wfDemands(doc driven.WorkflowDoc) driven.WorkflowPerms {
	out := driven.WorkflowPerms{Scopes: map[string]string{}}
	merge := func(p driven.WorkflowPerms) {
		if permLevel(p.All) > permLevel(out.All) {
			out.All = p.All
		}
		for k, v := range p.Scopes {
			if permLevel(v) > permLevel(out.Scopes[k]) {
				out.Scopes[k] = v
			}
		}
	}
	merge(doc.HeadPerms)
	for _, p := range doc.JobPerms {
		merge(p)
	}
	return out
}

// CheckWorkflows ist das Regelmodul workflows (DC-FA-WF-001): hermetisch (nur
// Filesystem- und Parser-Port, kein git, kein Netz, kein Ausfuehren), opt-in
// ueber Dir.
func CheckWorkflows(fsys driven.Filesystem, wp driven.WorkflowParser, cfg model.WorkflowsConfig) []model.Finding {
	if strings.TrimSpace(cfg.Dir) == "" || wp == nil {
		return nil // inert: keine Datei wird geoeffnet
	}
	files := workflowCandidates(fsys, cfg)
	var out []model.Finding
	refs := 0
	for _, f := range files {
		doc, err := readWorkflow(fsys, wp, f)
		if err != nil {
			out = append(out, model.Finding{File: f, Line: 1, Rule: "workflows", Target: f,
				Reason: ReasonWorkflowUnparsable, Message: "Datei ist kein gueltiges YAML: " + err.Error()})
			continue
		}
		refs += len(doc.Uses)
		for _, r := range doc.Uses {
			out = append(out, checkRef(fsys, wp, f, r)...)
		}
	}
	// FAIL-CLOSED: nichts gefunden und nichts zu pruefen sehen im Exit-Code
	// sonst identisch aus.
	if len(files) == 0 || refs == 0 {
		out = append(out, model.Finding{File: cfg.Dir, Line: 1, Rule: "workflows", Target: cfg.Dir,
			Reason: ReasonUsesPinMissing,
			Message: fmt.Sprintf("leere Pruefmenge: %d Workflow-Datei(en), %d uses:-Referenz(en) — fail-closed",
				len(files), refs)})
	}
	return out
}

// readWorkflow liest und uebersetzt eine Datei ueber die beiden Ports.
func readWorkflow(fsys driven.Filesystem, wp driven.WorkflowParser, rel string) (driven.WorkflowDoc, error) {
	b, err := fsys.ReadFile(rel)
	if err != nil {
		return driven.WorkflowDoc{}, err
	}
	return wp.ParseWorkflow(b)
}

// workflowCandidates liefert die Kandidaten stabil sortiert: Dateien
// UNMITTELBAR in Dir mit Endung .yml ODER .yaml (beide, weil CI-Systeme beide
// lesen), abzueglich exempt-paths.
func workflowCandidates(fsys driven.Filesystem, cfg model.WorkflowsConfig) []string {
	entries, err := fsys.List(cfg.Dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.Kind != driven.KindFile {
			continue
		}
		if !strings.HasSuffix(e.Name, ".yml") && !strings.HasSuffix(e.Name, ".yaml") {
			continue
		}
		rel := path.Join(cfg.Dir, e.Name)
		if matchAnyWorkflowGlob(cfg.ExemptPaths, rel) {
			continue
		}
		out = append(out, rel)
	}
	sort.Strings(out)
	return out
}

// matchAnyWorkflowGlob prueft die Ventil-Globs (Go path.Match, wie jedes
// andere d-check-Glob).
func matchAnyWorkflowGlob(globs []string, rel string) bool {
	for _, g := range globs {
		if ok, err := path.Match(g, rel); err == nil && ok {
			return true
		}
	}
	return false
}

// checkRef prueft EINE Referenz: fremde auf die Pin-Form, lokale auf Existenz
// und -- nur auf Job-Ebene -- auf Rechte-Deckung.
func checkRef(fsys driven.Filesystem, wp driven.WorkflowParser, file string, r driven.WorkflowUses) []model.Finding {
	find := func(reason, msg, target string) model.Finding {
		return model.Finding{File: file, Line: r.Line, Rule: "workflows", Target: target,
			Reason: reason, Message: msg}
	}
	if !strings.HasPrefix(strings.TrimSpace(r.Value), "./") {
		// FREMDE Referenz -- auch eine auf einen Workflow in einem anderen
		// Repository: ihr Inhalt liegt nicht vor, also nur die Pin-Form.
		if !pinnedRe.MatchString(r.Value) {
			return []model.Finding{find(ReasonUsesPinMissing, "kein voller 40-stelliger Commit-SHA: "+r.Value, r.Value)}
		}
		if !hasTagComment(r.Comment) {
			return []model.Finding{find(ReasonUsesPinUntagged, "SHA ohne Tag-Kommentar dahinter: "+r.Value, r.Value)}
		}
		return nil
	}
	target := path.Clean(strings.TrimPrefix(strings.TrimSpace(r.Value), "./"))
	kind, err := fsys.Kind(target)
	if err != nil || kind == driven.KindMissing {
		return []model.Finding{find(ReasonUsesLocalMissing, "lokale Referenz zeigt auf kein existierendes Ziel: ./"+target, target)}
	}
	// Nur der aufgerufene WORKFLOW (Job-Ebene) kennt eine Rechte-Weitergabe.
	// Eine lokale Action (Step-Ebene) erbt die Job-Rechte und deklariert
	// nichts -- dort gibt es nichts zu vergleichen.
	if !r.JobLevel {
		return nil
	}
	doc, perr := readWorkflow(fsys, wp, target)
	if perr != nil {
		return []model.Finding{find(ReasonWorkflowUnparsable,
			"Ziel der lokalen Referenz ist kein gueltiges YAML: "+perr.Error(), target)}
	}
	return checkPerms(find, target, r, wfDemands(doc))
}

// checkPerms vergleicht die geforderten mit den gefuehrten Rechten.
func checkPerms(find func(reason, msg, target string) model.Finding,
	target string, r driven.WorkflowUses, want driven.WorkflowPerms,
) []model.Finding {
	scopes := permDemands(want)
	if len(scopes) == 0 && permLevel(want.All) == 0 {
		return nil // das Ziel verlangt nichts -- nichts zu vergleichen
	}
	if !r.JobPerms.Declared {
		return []model.Finding{find(ReasonUsesLocalPermsUndecl,
			"Job traegt kein eigenes permissions:, das Ziel ./"+target+" verlangt aber Rechte — "+
				"er erbt den Workflow-Kopf und kann nichts weitergeben, was er nicht deklariert", target)}
	}
	if permLevel(want.All) > 0 {
		// Ein read-all/write-all des ZIELS fordert jeden Scope, den der
		// Aufrufer nennt -- mehr laesst sich von hier nicht bestimmen.
		for s := range r.JobPerms.Scopes {
			if !containsStr(scopes, s) {
				scopes = append(scopes, s)
			}
		}
		sort.Strings(scopes)
	}
	var out []model.Finding
	for _, s := range scopes {
		if have, need := permLevelOf(r.JobPerms, s), permLevelOf(want, s); have < need {
			out = append(out, find(ReasonUsesLocalPermsNarrow,
				fmt.Sprintf("Scope %s: Aufrufer fuehrt %s, Ziel ./%s verlangt %s",
					s, levelName(have), target, levelName(need)), target))
		}
	}
	return out
}

func containsStr(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
