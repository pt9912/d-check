// Package workflowyaml ist der Adapter des Workflow-Ports (DC-FA-WF-001): er
// übersetzt eine CI-Workflow-Datei in die Port-Struktur. Die yaml-Abhängigkeit
// lebt hier, nicht im Kern (Hexagon-Schnitt ADR-0004/ADR-0012; dieselbe
// Ansiedlung wie beim report-Adapter, ADR-0009).
//
// WARUM DER BAUM UND KEINE TEXTSUCHE: die Rechte-Frage braucht die
// Job-Zuordnung einer Referenz, und der Tag-Kommentar ist in YAML NICHT Teil
// des Werts (`uses: x@sha # v1` sind zwei Knoten-Eigenschaften). Eine
// zeilenweise Näherung verwechselt beides.
package workflowyaml

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Parser implementiert driven.WorkflowParser. Zustandslos.
type Parser struct{}

// New liefert den Adapter.
func New() *Parser { return &Parser{} }

// ParseWorkflow übersetzt den Datei-Inhalt in die Port-Struktur.
func (*Parser) ParseWorkflow(content []byte) (driven.WorkflowDoc, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return driven.WorkflowDoc{}, err
	}
	if len(doc.Content) == 0 {
		return driven.WorkflowDoc{}, fmt.Errorf("leeres Dokument")
	}
	root := doc.Content[0]
	out := driven.WorkflowDoc{HeadPerms: perms(mapValue(root, "permissions"))}
	jobs := mapValue(root, "jobs")
	if jobs == nil || jobs.Kind != yaml.MappingNode {
		return out, nil
	}
	for i := 0; i+1 < len(jobs.Content); i += 2 {
		job := jobs.Content[i+1]
		jp := perms(mapValue(job, "permissions"))
		out.JobPerms = append(out.JobPerms, jp)
		if u := mapValue(job, "uses"); u != nil {
			out.Uses = append(out.Uses, driven.WorkflowUses{
				Value: u.Value, Comment: u.LineComment, Line: u.Line,
				JobLevel: true, JobPerms: jp,
			})
		}
		steps := mapValue(job, "steps")
		if steps == nil || steps.Kind != yaml.SequenceNode {
			continue
		}
		for _, st := range steps.Content {
			if u := mapValue(st, "uses"); u != nil {
				out.Uses = append(out.Uses, driven.WorkflowUses{
					Value: u.Value, Comment: u.LineComment, Line: u.Line, JobPerms: jp,
				})
			}
		}
	}
	return out, nil
}

// perms liest einen `permissions:`-Knoten. Ein Skalar ist read-all/write-all,
// ein Mapping die Scope-Liste. Ein FEHLENDER Knoten ist nicht dasselbe wie ein
// leeres Mapping — nur der erste erbt still vom Workflow-Kopf.
func perms(n *yaml.Node) driven.WorkflowPerms {
	if n == nil {
		return driven.WorkflowPerms{}
	}
	out := driven.WorkflowPerms{Scopes: map[string]string{}, Declared: true}
	if n.Kind == yaml.ScalarNode {
		switch n.Value {
		case "read-all":
			out.All = "read"
		case "write-all":
			out.All = "write"
		}
		return out
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		out.Scopes[n.Content[i].Value] = n.Content[i+1].Value
	}
	return out
}

// mapValue liefert den Wert-Knoten zu einem Mapping-Schlüssel.
func mapValue(n *yaml.Node, key string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}
