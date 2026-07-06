package app

import (
	"github.com/pt9912/d-check/internal/hexagon/core/rules"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// idShape ist die allgemeine Kennungs-Gestalt für die Extraktion aus
// Headings (DC-FA-CLI-006.a Schritt 2): Präfix aus Großbuchstaben-Segmenten,
// dann -NNN, optional ein Suffix-Buchstabe.
var idShape = regexp.MustCompile(`^([A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*)-(\d+)([A-Za-z]?)$`)

// reqShape erkennt eine Anforderungs-Kennung (`<PREFIX>-FA-<BEREICH>-NNN`
// bzw. `<PREFIX>-QA-NN`) und fängt das **Projekt-Präfix** (Segment vor
// `-FA-`/`-QA-`). Grundlage der Präfix-Ableitung im ai-harness-Modus
// (slice-037).
var reqShape = regexp.MustCompile(`^([A-Z][A-Z0-9]*)-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?$`)

// idPrefixShape validiert einen explizit angegebenen --id-prefix: ein
// Großbuchstaben-Segment (z. B. `DC`, `AC`, `BC`).
var idPrefixShape = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)

// ValidIDPrefix meldet, ob ein --id-prefix-Wert der erlaubten Gestalt
// entspricht (CLI prüft das vor dem Lauf — Negative: Exit 2).
func ValidIDPrefix(p string) bool { return idPrefixShape.MatchString(p) }

// suggestedPattern ist ein aus einer Autoritäts-Quelle abgeleitetes
// ids-Muster samt der Quell-Kennungen (für den Nachweis-Kommentar).
type suggestedPattern struct {
	target string
	regex  string
	ids    []string
}

// harnessSource ist das reservierte Schlüsselwort (kein Pfad), das die
// ai-harness-course-Vorlage statt einer Quellen-Ableitung anfordert
// (DC-FA-CLI-006.a). harnessBaseline ist der im Header genannte Pin.
const (
	harnessSource     = "ai-harness"      // Mode 2: repo-bewusst
	harnessInitSource = "ai-harness-init" // Mode 1: Voll-Kanon
	harnessBaseline   = "v1.3.0"
	harnessExempt     = `[CHANGELOG.md, "docs/reviews/**"]`
	// harnessRoadmap ist der Konventions-Pfad der Planning-Roadmap; der
	// repo-bewusste planning-Block der ai-harness-Vorlage hängt an seiner Existenz
	// (DC-FA-CLI-006, K1–K4).
	harnessRoadmap = "docs/plan/planning/in-progress/roadmap.md"
)

// SuggestConfig liest die benannten Autoritäts-Quellen, leitet je Quelle
// ein ids-Muster aus den dort definierten Kennungen ab und liefert ein
// kommentiertes .d-check.yml-Gerüst (DC-FA-CLI-006.a). Die reservierte
// Quelle `ai-harness` erzeugt stattdessen die Harness-Vorlage (repo-
// bewusst), kombinierbar mit echten Quellen. Reiner Lese-Pfad — es wird
// nie geschrieben (DC-QA-03). Eine fehlende oder die Repo-Wurzel
// verlassende (echte) Quelle ist ein Fehler (CLI: Exit 2).
func SuggestConfig(fsys driven.Filesystem, sources []string, idPrefix string) (string, error) {
	initMode, harness := false, false
	var realSrc []string
	for _, src := range sources {
		switch src {
		case harnessInitSource:
			initMode = true
		case harnessSource:
			harness = true
		default:
			realSrc = append(realSrc, src)
		}
	}
	var patterns []suggestedPattern
	for _, src := range realSrc {
		rel, escaped := rules.ResolveConfigPath(src)
		if escaped {
			return "", fmt.Errorf("Autoritäts-Quelle verlässt die Repository-Wurzel: %s", src)
		}
		kind, err := fsys.Kind(rel)
		if err != nil || kind == driven.KindMissing {
			return "", fmt.Errorf("Autoritäts-Quelle existiert nicht: %s", src)
		}
		ids, err := extractDefinedIDs(fsys, rel, kind)
		if err != nil {
			return "", err
		}
		patterns = append(patterns, suggestedPattern{target: src, regex: deriveRegex(ids), ids: ids})
	}
	if initMode || harness {
		// Anforderungs-Präfix (slice-037, ADR-0015): explizit via --id-prefix;
		// sonst im repo-bewussten Modus aus dem Lastenheft abgeleitet; sonst
		// Platzhalter — kein stiller DC--Default in Fremd-Repos.
		reqPrefix := idPrefix
		if reqPrefix == "" && harness {
			derived, err := deriveReqPrefix(fsys)
			if err != nil {
				return "", err
			}
			reqPrefix = derived
		}
		// initMode (Voll-Kanon) hat Vorrang: repoAware nur bei reinem ai-harness.
		return renderHarness(fsys, patterns, !initMode, reqPrefix), nil
	}
	return renderSuggestion(patterns, probeOptInModules(fsys)), nil
}

// extractDefinedIDs sammelt die in den Headings einer Quelle (Datei oder
// Verzeichnis) definierten Kennungen — das führende Token, das idShape
// erfüllt (kein Fließtext-Mining). Ergebnis bytewise sortiert.
func extractDefinedIDs(fsys driven.Filesystem, rel string, kind driven.EntryKind) ([]string, error) {
	files := []string{rel}
	if kind == driven.KindDir {
		discovered, err := rules.DiscoverFiles(fsys, []string{rel}, nil)
		if err != nil {
			return nil, err
		}
		files = discovered
	}
	seen := map[string]bool{}
	for _, f := range files {
		content, err := fsys.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, h := range rules.ExtractHeadings(content) {
			fields := strings.Fields(rules.StripHeadingLinks(h))
			if len(fields) == 0 {
				continue
			}
			// führendes Token, von Markup/Satzzeichen befreit
			// (`DC-…`, ADR-0001:, [ADR-0001](…) → die nackte Kennung)
			tok := strings.Trim(fields[0], "`.,:;")
			if idShape.MatchString(tok) {
				seen[tok] = true
			}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

// deriveRegex bildet aus den Kennungen einer Quelle die Präfix-Alternation
// (?:p₁|p₂|…)-\d+[A-Za-z]? (DC-FA-CLI-006.a Schritt 3). Round-Trip-Invariante:
// das Ergebnis matcht jede Eingabe-Kennung. Leer, wenn keine Kennung vorliegt.
func deriveRegex(ids []string) string {
	prefixes := map[string]bool{}
	hasLetter := false
	for _, id := range ids {
		if m := idShape.FindStringSubmatch(id); m != nil {
			prefixes[m[1]] = true
			hasLetter = hasLetter || m[3] != ""
		}
	}
	if len(prefixes) == 0 {
		return ""
	}
	ps := make([]string, 0, len(prefixes))
	for p := range prefixes {
		ps = append(ps, regexp.QuoteMeta(p))
	}
	sort.Strings(ps)
	body := ps[0]
	if len(ps) > 1 {
		body = "(?:" + strings.Join(ps, "|") + ")"
	}
	body += `-\d+`
	if hasLetter {
		body += `[A-Za-z]?`
	}
	return body
}

// deriveReqPrefix leitet das Anforderungs-Präfix aus dem Lastenheft ab
// (ai-harness-Modus, slice-037): das **eindeutige** Projekt-Präfix aller
// FA-/QA-Kennungs-Headings in spec/lastenheft.md. Liefert "" wenn die Datei
// fehlt/keine Anforderungs-Kennung trägt (→ Platzhalter); Fehler bei
// **mehreren** verschiedenen Präfixen (der Mensch gibt --id-prefix explizit).
func deriveReqPrefix(fsys driven.Filesystem) (string, error) {
	const lh = "spec/lastenheft.md"
	if !pathExists(fsys, lh) {
		return "", nil
	}
	content, err := fsys.ReadFile(lh)
	if err != nil {
		return "", nil
	}
	prefixes := map[string]bool{}
	for _, h := range rules.ExtractHeadings(content) {
		fields := strings.Fields(rules.StripHeadingLinks(h))
		if len(fields) == 0 {
			continue
		}
		tok := strings.Trim(fields[0], "`.,:;")
		if m := reqShape.FindStringSubmatch(tok); m != nil {
			prefixes[m[1]] = true
		}
	}
	switch len(prefixes) {
	case 0:
		return "", nil
	case 1:
		var only string
		for p := range prefixes {
			only = p
		}
		return only, nil
	default:
		ps := make([]string, 0, len(prefixes))
		for p := range prefixes {
			ps = append(ps, p)
		}
		sort.Strings(ps)
		return "", fmt.Errorf("mehrdeutiges Anforderungs-Präfix im Lastenheft (%s) — --id-prefix explizit angeben", strings.Join(ps, ", "))
	}
}

// probeOptInModules lässt die opt-in-Module probeweise laufen und liefert
// (in fester Reihenfolge) jene mit ≥1 Befund (DC-FA-CLI-006.a Schritt 4).
func probeOptInModules(fsys driven.Filesystem) []string {
	optIn := []string{"codepaths", "spans", "hostpaths"}
	// Roots ["."] = derselbe Scope wie das gerenderte Gerüst (sonst
	// misst die Probe einen anderen Datei-Satz als die vorgeschlagene
	// model.Config; rules.Run nimmt die Module aus dem 4. Argument, nicht aus cfg).
	res, err := rules.Run(fsys, nil, model.Config{Roots: []string{"."}}, optIn)
	if err != nil {
		return nil
	}
	active := map[string]bool{}
	for _, f := range res.Findings {
		active[f.Rule] = true
	}
	var out []string
	for _, m := range optIn {
		if active[m] {
			out = append(out, m)
		}
	}
	return out
}

// renderSuggestion baut das kommentierte, dekodierbare Gerüst. Quellen
// ohne abgeleitete Kennungen werden als Hinweis vermerkt; gibt es
// ids-Muster, wird `ids` in die Modul-Liste aufgenommen (sonst wären
// die Muster im erzeugten model.Config inaktiv).
func renderSuggestion(patterns []suggestedPattern, probed []string) string {
	var withIDs []suggestedPattern
	var noIDs []string
	for _, p := range patterns {
		if p.regex != "" {
			withIDs = append(withIDs, p)
		} else {
			noIDs = append(noIDs, p.target)
		}
	}

	var b strings.Builder
	b.WriteString("# .d-check.yml — Vorschlag aus `d-check --suggest-config` (advisory).\n")
	b.WriteString("# Prüfen und verengen: die ids-Muster sind eine Best-Guess-Ableitung\n")
	b.WriteString("# aus den definierten Kennungen der benannten Quellen.\n\n")
	b.WriteString("scan:\n  roots: [\".\"]\n\n")

	modules := append([]string{"links", "anchors"}, probed...)
	if len(withIDs) > 0 {
		modules = append(modules, "ids")
	}
	b.WriteString("modules: [" + strings.Join(modules, ", ") + "]\n")

	for _, t := range noIDs {
		fmt.Fprintf(&b, "# Hinweis: in %q keine definierten Kennungen gefunden.\n", t)
	}
	if len(withIDs) == 0 {
		return b.String()
	}
	b.WriteString("\nids:\n  patterns:\n")
	for _, p := range withIDs {
		fmt.Fprintf(&b, "    # abgeleitet aus %d Kennung(en) in %s: %s\n",
			len(p.ids), p.target, strings.Join(p.ids, ", "))
		// target gequotet: ein Quellpfad mit YAML-Sonderzeichen (`:`, `#`)
		// soll das Gerüst nicht brechen oder still verfälschen.
		fmt.Fprintf(&b, "    - regex: '%s'\n      target: %q\n", p.regex, p.target)
		b.WriteString("      # link-policy: always   # einkommentieren für strenge Linkpflicht\n")
	}
	return b.String()
}

// pathExists meldet, ob ein (repo-relativer) Pfad im gescannten Baum
// existiert — Grundlage des repo-bewussten Zuschnitts der ai-harness-
// Vorlage. Trailing-Slash (Verzeichnis-Schreibweise) wird toleriert.
func pathExists(fsys driven.Filesystem, p string) bool {
	kind, err := fsys.Kind(strings.TrimSuffix(p, "/"))
	return err == nil && kind != driven.KindMissing
}

// existingRoots filtert eine feste Pfadliste deterministisch auf die im
// Baum vorhandenen (Reihenfolge erhalten).
func existingRoots(fsys driven.Filesystem, roots []string) []string {
	var out []string
	for _, r := range roots {
		if pathExists(fsys, r) {
			out = append(out, r)
		}
	}
	return out
}

// harnessIDPattern ist ein kanonisches ids-Muster der Baseline-Konvention.
// always: ohne festes Definitions-target → stets auskommentiert.
type harnessIDPattern struct {
	regex  string
	target string
	always bool
	hint   string
	todo   string // optionaler TODO-Kommentar vor dem (aktiven) Muster
}

// harnessIDPatterns spiegelt die ids-Konvention von .d-check.yml in fester
// Reihenfolge (DC-QA-02; DC-FA-CLI-006.a). Nur das **Anforderungs**-Muster
// trägt ein projektspezifisches Präfix (reqPrefix); ADR/MR/slice/CO sind
// konventions-fest. Leeres reqPrefix → markierter Platzhalter `<PREFIX>`
// plus TODO statt eines stillen `DC-` (slice-037, ADR-0015).
func harnessIDPatterns(reqPrefix string) []harnessIDPattern {
	todo := ""
	if reqPrefix == "" {
		reqPrefix = "<PREFIX>"
		todo = "TODO: <PREFIX> durch das Projekt-Kennungs-Präfix ersetzen (oder --id-prefix angeben)"
	}
	return []harnessIDPattern{
		{regex: `ADR-\d{4}`, target: "docs/plan/adr/"},
		{regex: `MR-\d{3}`, target: "harness/conventions.md"},
		{regex: reqPrefix + `-(FA-[A-Z]+|QA)-\d+`, target: "spec/lastenheft.md", todo: todo},
		{regex: `slice-\d{3}`, target: "docs/plan/planning/"},
		{regex: `CO-\d{3}`, always: true, hint: "Carveouts: target setzen, falls genutzt"},
	}
}

// harnessClass ist eine kanonische matrix-Klasse; probe ist der
// Existenz-Prüfpfad, paths sind YAML-fertig (Globs gequotet).
type harnessClass struct {
	name  string
	paths []string
	probe string
	// order/direction (DC-FA-MTX-002): klasseninterne Verweisrichtung; nur
	// gesetzt, wo eine geschichtete Klasse die Source-Precedence kodiert.
	order     []string
	direction string
	// token (DC-FA-MTX-003): Regex, das Referenzen auf diese Klasse als bare
	// ID-Token im Körper erkennt; nur für nicht-link-pflichtige Kennungen.
	token string
}

func harnessClasses() []harnessClass {
	return []harnessClass{
		{name: "spec-straten", paths: []string{"spec/lastenheft.md", "spec/spezifikation.md", "spec/architecture.md"}, probe: "spec",
			order: []string{"spec/lastenheft.md", "spec/spezifikation.md", "spec/architecture.md"}, direction: "no-downward"},
		{name: "adr", paths: []string{`"docs/plan/adr/[0-9]*.md"`}, probe: "docs/plan/adr"},
		{name: "slice", paths: []string{`"docs/plan/planning/**/slice-*.md"`}, probe: "docs/plan/planning",
			token: `slice-\d{3}`},
	}
}

// renderHarness baut die ai-harness-course-Vorlage (Baseline
// harnessBaseline). repoAware=true (Quelle `ai-harness`): nur im Baum
// vorhandene Pfade aktiv, fehlende auskommentiert mit Hinweis.
// repoAware=false (`ai-harness-init`): Voll-Kanon, alle Blöcke aktiv
// (Zielbild fürs leere Repo). extra sind aus echten Quellen abgeleitete
// Muster (Kombi-Aufruf). Deterministisch: feste Reihenfolge, keine
// Map-Iteration für die Ausgabe.
func renderHarness(fsys driven.Filesystem, extra []suggestedPattern, repoAware bool, reqPrefix string) string {
	var b strings.Builder
	b.WriteString("# .d-check.yml — Vorschlag aus `d-check --suggest-config` (advisory).\n")
	if repoAware {
		fmt.Fprintf(&b, "# Quelle ai-harness (repo-bewusst), Baseline %s — fehlende Artefakte auskommentiert. Prüfen und verengen.\n\n", harnessBaseline)
	} else {
		fmt.Fprintf(&b, "# Quelle ai-harness-init (Voll-Kanon), Baseline %s — Zielbild; läuft, sobald die Struktur (Scan-Wurzeln, ids-Targets) existiert. Prüfen und verengen.\n\n", harnessBaseline)
	}

	roots := []string{"spec", "docs", "harness"}
	if repoAware {
		if roots = existingRoots(fsys, roots); len(roots) == 0 {
			roots = []string{"."}
		}
	}
	fmt.Fprintf(&b, "scan:\n  roots: [%s]\n\n", strings.Join(roots, ", "))
	// Fixe Aktiv-Menge (DC-FA-CLI-006, Eignungs-Kriterium K1–K4): hermetische,
	// konfig-freie Baum-Scan-Module, die die adoptierte Konvention führt. planning
	// kommt hinzu, wenn seine Roadmap existiert (repo-bewusst) bzw. im Voll-Kanon —
	// sonst fiele es fail-closed (Exit 2 ohne Roadmap).
	modules := []string{"links", "anchors", "ids", "matrix", "codepaths", "spans", "hostpaths"}
	if !repoAware || pathExists(fsys, harnessRoadmap) {
		modules = append(modules, "planning")
	}
	fmt.Fprintf(&b, "modules: [%s]\n", strings.Join(modules, ", "))
	// Auffindbarkeit (DC-FA-CLI-006): die nicht qualifizierten situativen Module
	// werden nicht vorab aktiviert — Verweis aufs Voll-Schema statt stiller
	// Aktivierung eines inerten Moduls. vcs/commits brauchen eine Commit-Range
	// (kein .d-check.yml-Material) und gehen über --print-mk; versions/targets sind
	// repo-spezifisch (pin-pattern/authority) und bewusst vertagt.
	b.WriteString("# Weitere opt-in-Module sind situativ und hier nicht vorab aktiviert:\n")
	b.WriteString("# external, diagrams, versions, pins, immutable, tracked, targets — Voll-Schema: d-check --print-config.\n")
	b.WriteString("# vcs/commits brauchen eine Commit-Range und werden als Makefile-Target verteilt: d-check --print-mk.\n\n")
	b.WriteString(renderHarnessIDs(fsys, extra, repoAware, reqPrefix))
	b.WriteString("\n")
	b.WriteString(renderHarnessMatrix(fsys, repoAware))
	b.WriteString("\n")
	b.WriteString(renderHarnessPlanning(fsys, repoAware))
	// codepaths: das Datei-Ventil der Konvention — Review-Reports/Changelog
	// zitieren naturgemäß Datei:Zeile/Pfade und sollen kein codepath-missing
	// auslösen (Parität zum ids-exempt; roots bleiben repo-spezifisch).
	// ignore-refs startet leer (Tombstone-Register entfernter Artefakte ist
	// repo-historie-spezifisch) — nur als Hinweis auskommentiert.
	b.WriteString("\ncodepaths:\n  exempt-paths: " + harnessExempt + "\n")
	b.WriteString("  # ignore-refs: [\"tools/altes-skript.sh\"]  # Ziel-Pfade entfernter Artefakte (referenz-weit; leer starten, DC-FA-CODE-001)\n")
	return b.String()
}

// renderHarnessIDs rendert den ids-Block: scope und kanonische Muster.
// repoAware=true: Muster aktiv nur bei vorhandenem Target (sonst
// auskommentiert mit Hinweis), scope auf existierende roots verengt.
// repoAware=false: alle Muster aktiv (außer Carveout ohne Target), voller
// scope. Danach die extra-Muster echter Quellen.
func renderHarnessIDs(fsys driven.Filesystem, extra []suggestedPattern, repoAware bool, reqPrefix string) string {
	var b strings.Builder
	b.WriteString("ids:\n")
	scope := []string{"spec", "docs/user"}
	if repoAware {
		scope = existingRoots(fsys, scope)
	}
	if len(scope) > 0 {
		fmt.Fprintf(&b, "  scope:\n    roots: [%s]\n", strings.Join(scope, ", "))
	}
	b.WriteString("  patterns:\n")
	for _, p := range harnessIDPatterns(reqPrefix) {
		if !p.always && (!repoAware || pathExists(fsys, p.target)) {
			if p.todo != "" {
				fmt.Fprintf(&b, "    # %s\n", p.todo)
			}
			fmt.Fprintf(&b, "    - regex: '%s'\n      target: %s\n      link-policy: always\n      exempt-paths: %s\n",
				p.regex, p.target, harnessExempt)
			continue
		}
		if p.always {
			fmt.Fprintf(&b, "    # %s\n", p.hint)
		} else {
			fmt.Fprintf(&b, "    # Hinweis: Target %s fehlt im Repo — Muster auskommentiert.\n", p.target)
		}
		tgt := p.target
		if tgt == "" {
			tgt = "<definition>"
		}
		fmt.Fprintf(&b, "    # - regex: '%s'\n    #   target: %s\n    #   link-policy: always\n    #   exempt-paths: %s\n",
			p.regex, tgt, harnessExempt)
	}
	for _, p := range extra {
		if p.regex == "" {
			fmt.Fprintf(&b, "    # Hinweis: in %q keine definierten Kennungen gefunden.\n", p.target)
			continue
		}
		fmt.Fprintf(&b, "    # abgeleitet aus %d Kennung(en) in %s: %s\n",
			len(p.ids), p.target, strings.Join(p.ids, ", "))
		fmt.Fprintf(&b, "    - regex: '%s'\n      target: %q\n", p.regex, p.target)
	}
	return b.String()
}

// renderHarnessMatrix rendert den matrix-Block. repoAware=true: Klassen
// aktiv bei vorhandenem Pfad (sonst auskommentiert), Regeln nur aktiv,
// wenn beide Endpunkt-Klassen aktiv sind (keine baumelnde Referenz).
// repoAware=false: alle Klassen und Regeln aktiv. status und
// exclude-sections sind pfad-unabhängig und immer aktiv.
func renderHarnessMatrix(fsys driven.Filesystem, repoAware bool) string {
	classes := harnessClasses()
	active := make(map[string]bool, len(classes))
	for _, c := range classes {
		active[c.name] = !repoAware || pathExists(fsys, c.probe)
	}
	var b strings.Builder
	b.WriteString("matrix:\n  classes:\n")
	for _, c := range classes {
		pre := ""
		if !active[c.name] {
			fmt.Fprintf(&b, "    # Hinweis: %s fehlt im Repo — Klasse auskommentiert.\n", c.probe)
			pre = "# "
		}
		if len(c.order) > 0 || c.token != "" {
			// Blockform: order/direction (DC-FA-MTX-002) und/oder token (DC-FA-MTX-003).
			fmt.Fprintf(&b, "    %s- name: %s\n", pre, c.name)
			fmt.Fprintf(&b, "    %s  paths: [%s]\n", pre, strings.Join(c.paths, ", "))
			if len(c.order) > 0 {
				fmt.Fprintf(&b, "    %s  order: [%s]  # autoritativste Schicht zuerst\n", pre, strings.Join(c.order, ", "))
				fmt.Fprintf(&b, "    %s  direction: %s  # Abwärtsverweis ⇒ matrix-downward\n", pre, c.direction)
			}
			if c.token != "" {
				fmt.Fprintf(&b, "    %s  token: '%s'  # ID-Token im Körper ⇒ matrix-forbidden, außer per Marker deklariert\n", pre, c.token)
			}
			continue
		}
		fmt.Fprintf(&b, "    %s- {name: %s, paths: [%s]}\n", pre, c.name, strings.Join(c.paths, ", "))
	}
	b.WriteString("  rules:\n")
	for _, r := range [][2]string{{"spec-straten", "adr"}, {"spec-straten", "slice"}, {"adr", "slice"}} {
		pre := ""
		if !active[r[0]] || !active[r[1]] {
			pre = "# "
		}
		fmt.Fprintf(&b, "    %s- {from: %s, to: %s, allow: false}\n", pre, r[0], r[1])
	}
	b.WriteString("  status:\n    forbidden: [superseded, deprecated]\n")
	b.WriteString("  exclude-sections: [Historie, \"7. Historie\", Geschichte]\n")
	// exempt-paths ist repo-spezifisch (Grandfathering der vor Einführung
	// Accepted-ADRs) — als Kommentar, der Adopter trägt die konkrete Liste ein.
	b.WriteString("  # exempt-paths: [\"docs/plan/adr/0001-*.md\"]  # immutable Alt-ADRs grandfathern (DC-FA-MTX-003)\n")
	return b.String()
}

// renderHarnessPlanning rendert den planning-Block (DC-FA-CLI-006, K1–K4).
// repoAware=true: aktiv nur bei vorhandener Roadmap (sonst auskommentiert mit
// Hinweis — das Modul fällt ohne Roadmap fail-closed, DC-FA-PLAN-001). repoAware=
// false (Voll-Kanon): aktiv. Nur roadmap wird gesetzt; heading/marker/slice-glob
// sind Konventions-Defaults. Der Aktiv-Zustand deckt sich mit der Aufnahme von
// planning ins modules-Set in renderHarness (gleiche Bedingung).
func renderHarnessPlanning(fsys driven.Filesystem, repoAware bool) string {
	if repoAware && !pathExists(fsys, harnessRoadmap) {
		return "# Hinweis: " + harnessRoadmap + " fehlt im Repo — Modul planning auskommentiert.\n" +
			"# planning:\n#   roadmap: " + harnessRoadmap + "\n"
	}
	return "planning:\n  roadmap: " + harnessRoadmap + "\n"
}
