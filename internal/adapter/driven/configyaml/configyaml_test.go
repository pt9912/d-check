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
