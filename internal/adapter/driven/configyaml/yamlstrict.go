package configyaml

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// decodeStrictNode dekodiert einen Teilbaum und lehnt unbekannte Schlüssel ab,
// wie es KnownFields(true) am Dokument tut — yaml.Node.Decode erbt diese
// Strenge nicht. Aliase und Merge-Keys löst Node.Decode auf; die Schlüssel
// prüft unknownKeys am selben Knotenbaum, so dass die Zeilen aus dem Dokument
// kommen (DC-FA-REQ-001).
func decodeStrictNode(n *yaml.Node, out any) error {
	var errs []string
	if err := n.Decode(out); err != nil {
		var te *yaml.TypeError
		if !errors.As(err, &te) {
			return err
		}
		errs = append(errs, te.Errors...)
	}
	unknownKeys(n, reflect.TypeOf(out), &errs)
	if len(errs) == 0 {
		return nil
	}
	return &yaml.TypeError{Errors: errs}
}

// unknownKeys meldet jeden Schlüssel eines Mappings, den der Zieltyp nicht
// über sein yaml-Tag kennt. Map-Schlüssel sind frei und werden nicht geprüft;
// ein yaml.Node-Feld nimmt beliebigen Inhalt.
func unknownKeys(n *yaml.Node, t reflect.Type, errs *[]string) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	n = resolveAlias(n)
	switch {
	case t == reflect.TypeFor[yaml.Node]():
	case t.Kind() == reflect.Slice && n.Kind == yaml.SequenceNode:
		for _, c := range n.Content {
			unknownKeys(c, t.Elem(), errs)
		}
	case t.Kind() == reflect.Map && n.Kind == yaml.MappingNode:
		for _, kv := range mappingPairs(n) {
			unknownKeys(kv[1], t.Elem(), errs)
		}
	case t.Kind() == reflect.Struct && n.Kind == yaml.MappingNode:
		fields := yamlFieldTypes(t)
		for _, kv := range mappingPairs(n) {
			ft, ok := fields[kv[0].Value]
			if !ok {
				*errs = append(*errs, fmt.Sprintf("line %d: field %s not found in type configyaml.%s", kv[0].Line, kv[0].Value, t.Name()))
				continue
			}
			unknownKeys(kv[1], ft, errs)
		}
	}
}

// mappingPairs liefert die Schlüssel-Wert-Paare eines Mappings und löst dabei
// Merge-Keys (`<<: *anker`, auch als Liste von Ankern) in ihre Paare auf.
func mappingPairs(n *yaml.Node) [][2]*yaml.Node {
	var out [][2]*yaml.Node
	for i := 0; i+1 < len(n.Content); i += 2 {
		k, v := n.Content[i], n.Content[i+1]
		if k.Tag == "!!merge" {
			out = append(out, mergedPairs(v)...)
			continue
		}
		out = append(out, [2]*yaml.Node{k, v})
	}
	return out
}

// mergedPairs löst den Wert eines Merge-Keys auf: ein Mapping oder eine Liste
// von Mappings, jeweils auch als Alias.
func mergedPairs(v *yaml.Node) [][2]*yaml.Node {
	v = resolveAlias(v)
	switch v.Kind {
	case yaml.MappingNode:
		return mappingPairs(v)
	case yaml.SequenceNode:
		var out [][2]*yaml.Node
		for _, m := range v.Content {
			if m = resolveAlias(m); m.Kind == yaml.MappingNode {
				out = append(out, mappingPairs(m)...)
			}
		}
		return out
	}
	return nil
}

func resolveAlias(n *yaml.Node) *yaml.Node {
	for n.Kind == yaml.AliasNode && n.Alias != nil {
		n = n.Alias
	}
	return n
}

func yamlFieldTypes(t reflect.Type) map[string]reflect.Type {
	fields := map[string]reflect.Type{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = strings.ToLower(f.Name)
		}
		fields[name] = f.Type
	}
	return fields
}
