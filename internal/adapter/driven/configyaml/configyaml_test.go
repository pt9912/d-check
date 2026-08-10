package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

func TestDecode_Vollbeispiel(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(`
scan:
  roots: [docs, spec]
  ignore: ["docs/archive/**"]
modules: [links, anchors, ids]
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
matrix:
  classes:
    - name: contract
      paths: [spec/lastenheft.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
  rules:
    - {from: contract, to: adr, allow: false}
  status:
    forbidden: [superseded, deprecated]
  exclude-sections: [Historie]
external:
  timeout-seconds: 10
  parallel: 4
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Roots) != 2 || len(cfg.IDPatterns) != 1 || cfg.IDPatterns[0].Target != "docs/plan/adr/" {
		t.Fatalf("cfg = %+v", cfg)
	}
	// matrix-Konfiguration wird in den Kern durchgereicht (DC-FA-MTX-001)
	m := cfg.Matrix
	if len(m.Classes) != 2 || m.Classes[0].Name != "contract" ||
		len(m.Rules) != 1 || m.Rules[0].From != "contract" || m.Rules[0].Allow ||
		len(m.StatusForbidden) != 2 || len(m.ExcludeSections) != 1 {
		t.Fatalf("Matrix = %+v", m)
	}
}

// matrix.status fehlt → Default [superseded, deprecated]
// (spec/spezifikation.md §2).
func TestDecode_MatrixStatusDefault(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x.md]\n"))
	if err != nil {
		t.Fatal(err)
	}
	got := cfg.Matrix.StatusForbidden
	if len(got) != 2 || got[0] != "superseded" || got[1] != "deprecated" {
		t.Fatalf("StatusForbidden = %v", got)
	}
}

// matrix.status.allow-supersede-lineage + supersede-fields werden
// durchgereicht; leerer Feldname ist Konfigurationsfehler (Exit 2);
// Default (kein Flag) ist false/leer (DC-FA-MTX-001 0.14.0).
func TestDecode_MatrixSupersedeLineage(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x.md]\n" +
		"  status:\n    forbidden: [superseded]\n    allow-supersede-lineage: true\n" +
		"    supersede-fields: [Supersedes, Aenderungstyp]\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Matrix.AllowSupersedeLineage || len(cfg.Matrix.SupersedeFields) != 2 {
		t.Fatalf("Lineage = %+v", cfg.Matrix)
	}
	// leerer Feldname → Exit 2
	if _, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x]\n" +
		"  status:\n    supersede-fields: [\"\"]\n")); err == nil || !strings.Contains(err.Error(), "supersede-fields") {
		t.Fatalf("leerer Feldname: Fehler erwartet, got %v", err)
	}
	// Default: ohne Flag false/leer
	def, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x.md]\n"))
	if err != nil || def.Matrix.AllowSupersedeLineage || len(def.Matrix.SupersedeFields) != 0 {
		t.Fatalf("Default-Lineage = %+v, err = %v", def.Matrix, err)
	}
}

func TestDecode_UnbekannterSchluessel(t *testing.T) {
	_, err := configyaml.Decode([]byte("modul: x\n"))
	if err == nil || !strings.Contains(err.Error(), "line") {
		t.Fatalf("err = %v (Zeilenangabe erwartet)", err)
	}
	if strings.Contains(err.Error(), "configyaml.raw") {
		t.Fatalf("interner Typname leakt: %v", err)
	}
}

func TestDecode_LeereUndKommentarDatei(t *testing.T) {
	for _, content := range []string{"", "# nur kommentar\n", "\n\n"} {
		cfg, err := configyaml.Decode([]byte(content))
		if err != nil || cfg.Roots != nil || cfg.Modules != nil {
			t.Fatalf("configyaml.Decode(%q): cfg=%+v err=%v (Defaults erwartet)", content, cfg, err)
		}
	}
}

func TestDecode_UngueltigerRegex(t *testing.T) {
	_, err := configyaml.Decode([]byte("ids:\n  patterns:\n    - regex: '[unclosed'\n      target: docs/\n"))
	if err == nil || !strings.Contains(err.Error(), "regex") {
		t.Fatalf("err = %v", err)
	}
}

// DC-FA-ID-001 (0.8.0): link-policy akzeptiert nur prose|always;
// alles andere ist Konfigurationsfehler (Exit 2). Leerwert = Default
// prose und wird akzeptiert.
func TestDecode_LinkPolicy(t *testing.T) {
	base := "ids:\n  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs/\n      link-policy: "
	if _, err := configyaml.Decode([]byte(base + "strict\n")); err == nil || !strings.Contains(err.Error(), "link-policy") {
		t.Fatalf("ungültige link-policy: err = %v", err)
	}
	for _, ok := range []string{"prose", "always"} {
		if _, err := configyaml.Decode([]byte(base + ok + "\n")); err != nil {
			t.Fatalf("link-policy %q abgelehnt: %v", ok, err)
		}
	}
	cfg, err := configyaml.Decode([]byte(base + "always\n      exempt-paths: [CHANGELOG.md]\n"))
	if err != nil {
		t.Fatalf("exempt-paths abgelehnt: %v", err)
	}
	if len(cfg.IDPatterns) != 1 || cfg.IDPatterns[0].LinkPolicy != "always" ||
		len(cfg.IDPatterns[0].ExemptPaths) != 1 {
		t.Fatalf("link-policy/exempt-paths nicht übernommen: %+v", cfg.IDPatterns)
	}
}

// ids-Regexe, die den Leerstring matchen (X*, (ADR)?), würden an
// jeder Position Leer-Matches und damit unbrauchbare Befunde mit
// leerem Target erzeugen — Konfigurationsfehler statt Befund-Flut.
func TestDecode_LeerstringRegex(t *testing.T) {
	for _, re := range []string{"X*", "(ADR)?", ""} {
		_, err := configyaml.Decode([]byte("ids:\n  patterns:\n    - regex: '" + re + "'\n      target: docs/\n"))
		if err == nil || !strings.Contains(err.Error(), "Leerstring") {
			t.Fatalf("regex %q: err = %v (Leerstring-Ablehnung erwartet)", re, err)
		}
	}
}

// TestDecode_CommitsFehler deckt die drei applyCommits-Ablehnungs-Guards ab
// (DC-FA-COMMITS-001) — dieselbe Klasse wie TestDecode_LeerstringRegex für ids:
// ohne diese Negativtests machte ein Refactor, der den Leerstring-Guard (das
// Silent-Grün-Schutz-Ventil) entfernt, das Gate still grün, ohne dass ein Test fällt.
func TestDecode_CommitsFehler(t *testing.T) {
	// leerer id-patterns-Eintrag
	if _, err := configyaml.Decode([]byte("commits:\n  id-patterns:\n    - ''\n")); err == nil {
		t.Fatal("leeres commits.id-patterns[0]: Fehler erwartet")
	}
	// Leerstring-matchendes Muster ⇒ jeder Commit gälte als getraced (Silent-Grün)
	for _, re := range []string{".*", "X*", "(ADR)?"} {
		_, err := configyaml.Decode([]byte("commits:\n  id-patterns:\n    - '" + re + "'\n"))
		if err == nil || !strings.Contains(err.Error(), "Leerstring") {
			t.Fatalf("commits.id-patterns %q: err = %v (Leerstring-Ablehnung erwartet)", re, err)
		}
	}
	// nicht kompilierbares id-pattern
	if _, err := configyaml.Decode([]byte("commits:\n  id-patterns:\n    - '[unclosed'\n")); err == nil {
		t.Fatal("ungültiges commits.id-patterns-Regex: Fehler erwartet")
	}
	// nicht kompilierbares exempt-pattern
	if _, err := configyaml.Decode([]byte("commits:\n  id-patterns:\n    - 'ADR-\\d{4}'\n  exempt-pattern: '(unclosed'\n")); err == nil {
		t.Fatal("ungültiges commits.exempt-pattern-Regex: Fehler erwartet")
	}
}

// TestDecode_CommitsHappy: eine gültige commits-Config wird übernommen und die
// Regexe kompiliert (der positive Gegenpol zu TestDecode_CommitsFehler).
func TestDecode_CommitsHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("commits:\n  id-patterns:\n    - 'ADR-\\d{4}'\n    - 'slice-\\d+'\n  exempt-pattern: '^(Merge |Revert )'\n"))
	if err != nil {
		t.Fatalf("gültige commits-Config abgelehnt: %v", err)
	}
	if len(cfg.Commits.IDPatterns) != 2 {
		t.Fatalf("commits.id-patterns nicht übernommen: %+v", cfg.Commits)
	}
	if cfg.Commits.ExemptPattern == nil || !cfg.Commits.ExemptPattern.MatchString("Merge branch 'x'") {
		t.Fatalf("commits.exempt-pattern nicht kompiliert/wirksam: %+v", cfg.Commits.ExemptPattern)
	}
}

// TestDecode_HostpathsPrefixesFehler deckt die beiden applyHostpaths-Ablehnungs-
// Guards ab (DC-FA-HOST-001 / DC-FA-CONF-001) — dieselbe Silent-Grün-Klasse wie
// TestDecode_CommitsFehler: ohne diese Negativtests machte ein Refactor, der den
// leerer-Name- oder den '/'-Guard entfernt, die Validierung still nachgiebig, ohne
// dass ein Test fällt. Der '/'-Guard ist genau die Regel, gegen die das frühere
// Handbuch-Beispiel `prefixes: ["/home"]` verstieß (führender / ⇒ Exit 2, Fix c8c33a0).
func TestDecode_HostpathsPrefixesFehler(t *testing.T) {
	// Eintrag mit '/' ⇒ Konfigurationsfehler (bare Verzeichnisnamen; die Regex
	// klammert die umgebenden Slashes selbst — ein führender/innerer / matchte nie).
	for _, p := range []string{"/home", "home/x", "/"} {
		if _, err := configyaml.Decode([]byte("hostpaths:\n  prefixes: [\"" + p + "\"]\n")); err == nil ||
			!strings.Contains(err.Error(), "ohne '/'") {
			t.Fatalf("hostpaths.prefixes %q: err = %v (Ablehnung „ohne /“ erwartet)", p, err)
		}
	}
	// leerer Name ⇒ Konfigurationsfehler
	if _, err := configyaml.Decode([]byte("hostpaths:\n  prefixes: [\"\"]\n")); err == nil ||
		!strings.Contains(err.Error(), "leeren Namen") {
		t.Fatalf("leerer hostpaths.prefixes-Name: Fehler erwartet, got %v", err)
	}
}

// TestDecode_HostpathsPrefixesHappy: bare Verzeichnisnamen (ohne /) werden
// akzeptiert und in den Kern durchgereicht (der positive Gegenpol; zugleich die
// korrekte Form des Handbuch-Beispiels `prefixes: [home, Users]`).
func TestDecode_HostpathsPrefixesHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("hostpaths:\n  prefixes: [home, Users]\n"))
	if err != nil {
		t.Fatalf("gültige hostpaths.prefixes abgelehnt: %v", err)
	}
	if len(cfg.Hostpaths.Prefixes) != 2 ||
		cfg.Hostpaths.Prefixes[0] != "home" || cfg.Hostpaths.Prefixes[1] != "Users" {
		t.Fatalf("hostpaths.prefixes nicht durchgereicht: %+v", cfg.Hostpaths)
	}
}

// TestDecode_PlanningFehler deckt den applyPlanning-Escape-Guard + die
// scope-Ablehnung ab (DC-FA-PLAN-001; planning ist ein Post-Pass ohne Scan).
func TestDecode_PlanningFehler(t *testing.T) {
	for _, bad := range []string{
		"planning:\n  roadmap: /abs/roadmap.md\n",
		"planning:\n  roadmap: ../raus/roadmap.md\n",
		"planning:\n  roadmap: r.md\n  slice-glob: 'slice-[.md'\n", // ungültiges path.Match-Glob (fail-open-Schutz)
	} {
		if _, err := configyaml.Decode([]byte(bad)); err == nil {
			t.Fatalf("ungültige planning-Config akzeptiert: %q", bad)
		}
	}
	// planning.scope existiert nicht (Post-Pass ohne Datei-Scan) ⇒ strikter Decoder lehnt ab
	if _, err := configyaml.Decode([]byte("planning:\n  scope:\n    roots: [x]\n")); err == nil {
		t.Fatal("unbekanntes planning.scope muss abgelehnt werden")
	}
}

// TestDecode_ClosureFehler deckt den Config-Rand der Closure-Fähigkeit ab
// (DC-FA-PLAN-001): was der Kern nur schlucken könnte, bricht hier laut ab —
// jeder dieser Werte wäre sonst ein stilles Grün.
func TestDecode_ClosureFehler(t *testing.T) {
	for name, bad := range map[string]string{
		"dir absolut":         "planning:\n  closure:\n    dir: /abs/done\n",
		"dir mit ..":          "planning:\n  closure:\n    dir: ../raus\n",
		"heading-pattern RE2": "planning:\n  closure:\n    dir: done\n    heading-pattern: '^(['\n",
		"min-sentences 0":     "planning:\n  closure:\n    dir: done\n    min-sentences: 0\n",
		"min-sentences neg":   "planning:\n  closure:\n    dir: done\n    min-sentences: -1\n",
		"glob leer":           "planning:\n  closure:\n    dir: done\n    glob: ''\n",
		"glob ungültig":       "planning:\n  closure:\n    dir: done\n    glob: '[a-'\n",
		"leere Floskel":       "planning:\n  closure:\n    dir: done\n    boilerplate: ['']\n",
		"Floskel nur Space":   "planning:\n  closure:\n    dir: done\n    boilerplate: ['   ']\n",
	} {
		if _, err := configyaml.Decode([]byte(bad)); err == nil {
			t.Errorf("%s: ungültige closure-Config akzeptiert: %q", name, bad)
		}
	}
}

// TestDecode_ClosureHappy: gültige Werte werden durchgereicht, und ein
// ABWESENDES min-sentences bleibt vom explizit gesetzten unterscheidbar (der
// Kern-Default greift nur im ersten Fall).
func TestDecode_ClosureHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"planning:\n  roadmap: r.md\n  closure:\n    dir: docs/done\n    boilerplate: ['Fertig.']\n"))
	if err != nil {
		t.Fatalf("gültige closure-Config abgelehnt: %v", err)
	}
	c := cfg.Planning.Closure
	if c.Dir != "docs/done" || len(c.Boilerplate) != 1 {
		t.Fatalf("closure-Config nicht durchgereicht: %+v", c)
	}
	if got := c.EffectiveMinSentences(); got != 4 {
		t.Errorf("Default min-sentences = %d, want 4", got)
	}
	if got := c.EffectiveHeadingPattern(); got == "" {
		t.Error("Default heading-pattern darf nicht leer sein")
	}
	explicit, err := configyaml.Decode([]byte(
		"planning:\n  roadmap: r.md\n  closure:\n    dir: d\n    min-sentences: 2\n"))
	if err != nil {
		t.Fatalf("explizites min-sentences abgelehnt: %v", err)
	}
	if got := explicit.Planning.Closure.EffectiveMinSentences(); got != 2 {
		t.Errorf("explizites min-sentences = %d, want 2", got)
	}
}

// TestDecode_PlanningHappy: eine gültige planning-Config wird übernommen, die
// Defaults greifen.
func TestDecode_PlanningHappy(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("planning:\n  roadmap: docs/plan/planning/in-progress/roadmap.md\n  marker: IDLE\n"))
	if err != nil {
		t.Fatalf("gültige planning-Config abgelehnt: %v", err)
	}
	if cfg.Planning.Roadmap == "" || cfg.Planning.Marker != "IDLE" {
		t.Fatalf("planning-Config nicht übernommen: %+v", cfg.Planning)
	}
	if cfg.Planning.EffectiveHeading() != "## Aktuelle Welle" || cfg.Planning.EffectiveSliceGlob() != "slice-*.md" {
		t.Fatalf("planning-Defaults falsch: heading=%q glob=%q", cfg.Planning.EffectiveHeading(), cfg.Planning.EffectiveSliceGlob())
	}
}

func TestDecode_MatrixUndExternalConstraints(t *testing.T) {
	if _, err := configyaml.Decode([]byte("matrix:\n  classes:\n    - name: a\n      paths: [x]\n  rules:\n    - {from: a, to: fehlt, allow: false}\n")); err == nil {
		t.Fatal("undeklarierte Klasse: Fehler erwartet")
	}
	if _, err := configyaml.Decode([]byte("external:\n  timeout-seconds: 999\n")); err == nil {
		t.Fatal("timeout außerhalb 1–300: Fehler erwartet")
	}
	// explizit gesetzte 0 ist KEIN stiller Default, sondern Fehler
	// (Review R1 zu slice-008; Constraint 1–300 bzw. 1–16)
	if _, err := configyaml.Decode([]byte("external:\n  timeout-seconds: 0\n")); err == nil {
		t.Fatal("timeout-seconds: 0 muss Konfigurationsfehler sein")
	}
	if _, err := configyaml.Decode([]byte("external:\n  parallel: 0\n")); err == nil {
		t.Fatal("parallel: 0 muss Konfigurationsfehler sein")
	}
	if _, err := configyaml.Decode([]byte("external:\n  parallel: 99\n")); err == nil {
		t.Fatal("parallel außerhalb 1–16: Fehler erwartet")
	}
	// ids: Muster ohne Target ist unvollständig (DC-FA-CONF-001)
	if _, err := configyaml.Decode([]byte("ids:\n  patterns:\n    - regex: \"X-\\\\d\"\n")); err == nil {
		t.Fatal("ids-Muster ohne target: Fehler erwartet")
	}
	// matrix: doppelter Klassen-Name ist mehrdeutig
	doppelt := "matrix:\n  classes:\n    - name: a\n      paths: [x]\n    - name: a\n      paths: [y]\n"
	if _, err := configyaml.Decode([]byte(doppelt)); err == nil {
		t.Fatal("doppelter Klassen-Name: Fehler erwartet")
	}
	// codepaths: Präfixe müssen relativ und ..-frei sein (DC-FA-CONF-001)
	for _, bad := range []string{
		"codepaths:\n  roots: [\"\"]\n",
		"codepaths:\n  roots: [/abs]\n",
		"codepaths:\n  roots: [\"../raus\"]\n",
	} {
		if _, err := configyaml.Decode([]byte(bad)); err == nil {
			t.Fatalf("ungültiges codepaths-Präfix akzeptiert: %q", bad)
		}
	}
	cfg, err := configyaml.Decode([]byte("codepaths:\n  roots: [docs, tools]\n"))
	if err != nil || len(cfg.Codepaths.Roots) != 2 {
		t.Fatalf("codepaths.roots nicht übernommen: %+v (%v)", cfg.Codepaths, err)
	}
	// codepaths.exempt-paths wird übernommen (Datei-Ventil wie ids; slice-043)
	cfgEx, err := configyaml.Decode([]byte("codepaths:\n  roots: [docs]\n  exempt-paths: [\"docs/reviews/**\"]\n"))
	if err != nil || len(cfgEx.Codepaths.ExemptPaths) != 1 || cfgEx.Codepaths.ExemptPaths[0] != "docs/reviews/**" {
		t.Fatalf("codepaths.exempt-paths nicht übernommen: %+v (%v)", cfgEx.Codepaths, err)
	}
}

func TestDecode_OhneDatei(t *testing.T) {
	cfg, err := configyaml.Decode(nil)
	if err != nil || cfg.Roots != nil || cfg.Modules != nil {
		t.Fatalf("Defaults erwartet: %+v (%v)", cfg, err)
	}
}

// DC-FA-CONF-002: Modul-lokaler Scan-Scope — Decoding und Constraints.
func TestDecode_ModulScope(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"scan:\n  roots: [\".\"]\n" +
			"modules: [links, anchors, ids]\n" +
			"links:\n  scope:\n    roots: [docs]\n" +
			"ids:\n" +
			"  scope:\n    roots: [spec, docs]\n    ignore: [\"docs/archive/**\"]\n" +
			"  patterns:\n    - regex: 'ADR-\\d{4}'\n      target: docs\n"))
	if err != nil {
		t.Fatal(err)
	}
	ids := cfg.Scopes["ids"]
	if ids == nil || len(ids.Roots) != 2 || len(ids.Ignore) != 1 {
		t.Fatalf("ids.scope nicht übernommen: %+v", cfg.Scopes)
	}
	links := cfg.Scopes["links"]
	if links == nil || len(links.Roots) != 1 || links.Ignore != nil {
		t.Fatalf("links.scope nicht übernommen: %+v", cfg.Scopes)
	}
	if cfg.Scopes["anchors"] != nil {
		t.Fatalf("anchors ohne scope-Schlüssel darf keinen Scope tragen")
	}

	// explizit leere roots-Liste ist gültig (prüft nichts)
	cfg, err = configyaml.Decode([]byte("ids:\n  scope:\n    roots: []\n  patterns:\n    - regex: 'X-\\d'\n      target: docs\n"))
	if err != nil || cfg.Scopes["ids"] == nil || cfg.Scopes["ids"].Roots == nil || len(cfg.Scopes["ids"].Roots) != 0 {
		t.Fatalf("leere scope.roots-Liste nicht übernommen: %+v (%v)", cfg.Scopes, err)
	}
}

// scope ohne roots ist ein Konfigurationsfehler (keine stille
// Vererbung); unbekannte Schlüssel im scope sind Fehler (strikt).
func TestDecode_ModulScopeFehler(t *testing.T) {
	for _, bad := range []string{
		"links:\n  scope:\n    ignore: [\"x/**\"]\n",
		"anchors:\n  scope: {}\n",
		"matrix:\n  scope:\n    wurzeln: [docs]\n",
		"links:\n  unbekannt: 1\n",
	} {
		if _, err := configyaml.Decode([]byte(bad)); err == nil {
			t.Fatalf("ungültiger scope akzeptiert: %q", bad)
		}
	}
}

// DC-FA-VER-001: Modul versions — Decoding des öffentlichen YAML-Vertrags
// (pin-pattern, current-from, exempt-paths übernommen; current-from-Default).
func TestDecode_Versions(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"versions:\n" +
			"  pin-pattern: 'ghcr\\.io/[^\\s:]+:(v[0-9]+\\.[0-9]+\\.[0-9]+)'\n" +
			"  current-from: version.md#aktuell\n" +
			"  exempt-paths: [CHANGELOG.md, \"docs/plan/planning/done/**\"]\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Versions.PinPattern == nil ||
		cfg.Versions.CurrentFrom != "version.md#aktuell" ||
		len(cfg.Versions.ExemptPaths) != 2 {
		t.Fatalf("versions nicht übernommen: %+v", cfg.Versions)
	}
	// current-from optional → der Kern-Default greift (EffectiveCurrentFrom).
	cfgD, err := configyaml.Decode([]byte("versions:\n  pin-pattern: 'x:(v\\d+\\.\\d+\\.\\d+)'\n"))
	if err != nil || cfgD.Versions.EffectiveCurrentFrom() != "version.md#aktuell" {
		t.Fatalf("current-from-Default nicht wirksam: %+v (%v)", cfgD.Versions, err)
	}
}

// DC-FA-VER-001: pin-pattern ist Pflicht und muss ein gültiger, nicht den
// Leerstring matchender Regex sein; unbekannte Schlüssel sind Fehler (Exit 2).
func TestDecode_VersionsFehler(t *testing.T) {
	for _, bad := range []string{
		"versions:\n  current-from: version.md#aktuell\n", // pin-pattern fehlt
		"versions:\n  pin-pattern: ''\n",                  // leer
		"versions:\n  pin-pattern: '[unclosed'\n",         // Regex-Fehler
		"versions:\n  pin-pattern: 'x*'\n",                // matcht den Leerstring
		"versions:\n  pin-pattern: 'x'\n  unbekannt: 1\n", // unbekannter Schlüssel
	} {
		if _, err := configyaml.Decode([]byte(bad)); err == nil {
			t.Fatalf("ungültige versions-Config akzeptiert: %q", bad)
		}
	}
}

// DC-FA-VER-001 / DC-FA-CONF-002: versions.scope wird übernommen.
func TestDecode_VersionsScope(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"versions:\n  scope:\n    roots: [docs]\n  pin-pattern: 'x:(v\\d+\\.\\d+\\.\\d+)'\n"))
	if err != nil || cfg.Scopes["versions"] == nil || len(cfg.Scopes["versions"].Roots) != 1 {
		t.Fatalf("versions.scope nicht übernommen: %+v (%v)", cfg.Scopes, err)
	}
}

// DC-FA-MTX-002: order/direction werden übernommen; gültige Kombination ok.
func TestDecode_MatrixOrderDirection(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"matrix:\n  classes:\n    - name: a\n      paths: [x.md, y.md]\n" +
			"      order: [x.md, y.md]\n      direction: no-downward\n"))
	if err != nil {
		t.Fatalf("gültige order/direction: %v", err)
	}
	c := cfg.Matrix.Classes[0]
	if len(c.Order) != 2 || c.Direction != "no-downward" {
		t.Fatalf("order/direction nicht übernommen: %+v", c)
	}
}

// DC-FA-MTX-002 Negative (fail-closed): order/direction nur gemeinsam,
// direction nur no-downward — sonst Konfigurationsfehler.
func TestDecode_MatrixDirectionFailClosed(t *testing.T) {
	base := "matrix:\n  classes:\n    - name: a\n      paths: [x.md]\n"
	cases := map[string]string{
		"order ohne direction":  base + "      order: [x.md]\n",
		"direction ohne order":  base + "      direction: no-downward\n",
		"unbekannte direction":  base + "      order: [x.md]\n      direction: sideways\n",
	}
	for name, in := range cases {
		if _, err := configyaml.Decode([]byte(in)); err == nil {
			t.Errorf("%s: Konfigurationsfehler erwartet, kein Fehler", name)
		}
	}
}

// DC-FA-MTX-003: matrix.classes[].token wird kompiliert, exempt-paths übernommen.
func TestDecode_MatrixTokenUndExemptPaths(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"matrix:\n  classes:\n    - {name: slice, paths: [s.md], token: 'slice-\\d{3}'}\n" +
			"  exempt-paths: [\"docs/plan/adr/0001-*.md\"]\n"))
	if err != nil {
		t.Fatalf("gültige token/exempt-paths: %v", err)
	}
	c := cfg.Matrix.Classes[0]
	if c.Token == nil || !c.Token.MatchString("slice-042") {
		t.Fatalf("token nicht kompiliert/übernommen: %+v", c)
	}
	if len(cfg.Matrix.ExemptPaths) != 1 {
		t.Fatalf("exempt-paths nicht übernommen: %+v", cfg.Matrix.ExemptPaths)
	}
}

// DC-FA-MTX-003 Negative (fail-closed): token kompiliert nicht / matcht Leerstring.
func TestDecode_MatrixTokenFailClosed(t *testing.T) {
	for _, tok := range []string{"[unclosed", "x*"} { // x* matcht den Leerstring
		in := "matrix:\n  classes:\n    - {name: a, paths: [x.md], token: '" + tok + "'}\n"
		if _, err := configyaml.Decode([]byte(in)); err == nil {
			t.Errorf("token %q: Konfigurationsfehler erwartet", tok)
		}
	}
}

// DC-FA-SRC-001 Happy: gültige sources[]-Pins werden übernommen; sha256 wird
// case-insensitiv zu Kleinbuchstaben normalisiert, unpack-Default ist none, und
// die Befund-Zeile ist die yaml-Node-Zeile des url-Feldes.
func TestDecode_SourcesHappy(t *testing.T) {
	const hex64 = "ABCDEF0123456789abcdef0123456789ABCDEF0123456789abcdef0123456789"
	cfg, err := configyaml.Decode([]byte("modules: [links]\n" +
		"sources:\n" +
		"  - url: https://example.org/regelwerk.md\n" +
		"    sha256: " + hex64 + "\n" +
		"  - url: http://example.org/bundle.zip\n" +
		"    sha256: " + strings.ToLower(hex64) + "\n" +
		"    unpack: zip\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Sources.Pins) != 2 {
		t.Fatalf("Pins = %+v", cfg.Sources.Pins)
	}
	p0 := cfg.Sources.Pins[0]
	if p0.URL != "https://example.org/regelwerk.md" || p0.Unpack != "none" || p0.Sha256 != strings.ToLower(hex64) {
		t.Fatalf("Pin[0] = %+v (Default unpack none, sha256 kleingeschrieben erwartet)", p0)
	}
	if p0.Line != 3 { // url-Feld-Zeile im obigen Dokument (1: modules, 2: sources, 3: url)
		t.Fatalf("Pin[0].Line = %d, erwartet 3", p0.Line)
	}
	if cfg.Sources.Pins[1].Unpack != "zip" {
		t.Fatalf("Pin[1].unpack = %q, erwartet zip", cfg.Sources.Pins[1].Unpack)
	}
}

// DC-FA-SRC-001 fail-closed: jeder ungültige sources[]-Eintrag ⇒ Exit 2
// (fehlende/nicht-http(s)-url, sha256-Länge ≠ 64, unbekanntes unpack).
func TestDecode_SourcesFailClosed(t *testing.T) {
	const hex64 = "0000000000000000000000000000000000000000000000000000000000000000"
	cases := map[string]string{
		"url fehlt":       "sources:\n  - sha256: " + hex64 + "\n",
		"url nicht http":  "sources:\n  - url: ftp://example.org/x\n    sha256: " + hex64 + "\n",
		"url repo-intern": "sources:\n  - url: docs/x.md\n    sha256: " + hex64 + "\n",
		"sha256 zu kurz":  "sources:\n  - url: https://example.org/x\n    sha256: abcdef\n",
		"sha256 fehlt":    "sources:\n  - url: https://example.org/x\n",
		"unpack unbekannt": "sources:\n  - url: https://example.org/x\n    sha256: " + hex64 +
			"\n    unpack: foo\n",
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := configyaml.Decode([]byte(in)); err == nil {
				t.Fatalf("%s: Konfigurationsfehler (Exit 2) erwartet", name)
			}
		})
	}
}

// Der Default von closure.glob ist ein VERWEIS auf slice-glob, kein kopiertes
// Literal: abwesend ⇒ es gilt slice-glob, gesetzt ⇒ es gilt der eigene Wert.
// Beides über die Kern-Methode, damit der Verweis nicht nur im Adapter lebt.
func TestDecode_ClosureGlobDefaultIstVerweis(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"planning:\n  roadmap: r.md\n  slice-glob: 'paket-*.md'\n  closure:\n    dir: docs/done\n"))
	if err != nil {
		t.Fatalf("gültige Config abgelehnt: %v", err)
	}
	if got := cfg.Planning.EffectiveClosureGlob(); got != "paket-*.md" {
		t.Errorf("ohne closure.glob = %q, want den slice-glob 'paket-*.md'", got)
	}
	cfg, err = configyaml.Decode([]byte(
		"planning:\n  roadmap: r.md\n  slice-glob: 'paket-*.md'\n  closure:\n    dir: docs/done\n    glob: '*.md'\n"))
	if err != nil {
		t.Fatalf("gültige Config abgelehnt: %v", err)
	}
	if got := cfg.Planning.EffectiveClosureGlob(); got != "*.md" {
		t.Errorf("mit closure.glob = %q, want '*.md'", got)
	}
	if got := cfg.Planning.EffectiveSliceGlob(); got != "paket-*.md" {
		t.Errorf("slice-glob wurde mitgezogen: %q", got)
	}
}

