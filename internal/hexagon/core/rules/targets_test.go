package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const (
	tgtMakefile = "Makefile"
	tgtDoc      = "docs/x.md"
)

// tgtDocTable baut eine Doku mit einer Tabelle, die je Name eine `make <name>`-
// Zeile (Spalte 0, Pipe-Präfix) trägt — plus eine **Prosa**-Zeile mit
// `make prose-ghost`, das NICHT zählen darf (Tabellen-Scoping, ^|-Spalte-0).
func tgtDocTable(names ...string) string {
	var b strings.Builder
	b.WriteString("# Doc\n\nRichtig: `make prose-ghost` (Prosa, keine Tabellenzeile).\n\n| Target | Zweck |\n|---|---|\n")
	for _, n := range names {
		b.WriteString("| `make " + n + "` | x |\n")
	}
	return b.String()
}

// tgtCfg: makefiles+doc-tables+authority = derselbe Doc (beide Richtungen aktiv).
func tgtCfg() model.TargetsConfig {
	return model.TargetsConfig{
		Makefiles: []string{tgtMakefile}, DocTables: []string{tgtDoc}, Authority: tgtDoc,
	}
}

func mustCheck(t *testing.T, files map[string]string, cfg model.TargetsConfig) []model.Finding {
	t.Helper()
	f, err := CheckTargets(coretest.NewMemFS(files), cfg)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	return f
}

// TestCheckTargetsHappy: Makefile-Regeln == dokumentierte Targets ⇒ kein Befund
// (die Prosa-`make prose-ghost` darf ebenfalls nichts auslösen — Scoping-Guard).
func TestCheckTargetsHappy(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build test: dep\n\techo\nlint:\n\techo\n",
		tgtDoc:      tgtDocTable("build", "test", "lint"),
	}
	if f := mustCheck(t, files, tgtCfg()); len(f) != 0 {
		t.Fatalf("konsistent ⇒ kein Befund, bekam %+v", f)
	}
}

// TestCheckTargetsPhantom: dokumentiertes `make ghost` ohne Makefile-Regel ⇒
// gate-phantom (Richtung 1), an Datei:Zeile der Doku-Behauptung.
func TestCheckTargetsPhantom(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\n",
		tgtDoc:      tgtDocTable("build", "ghost"),
	}
	f := mustCheck(t, files, tgtCfg())
	if len(f) != 1 || f[0].Reason != ReasonGatePhantom || f[0].Target != "ghost" ||
		f[0].Rule != "targets" || f[0].File != tgtDoc {
		t.Fatalf("gate-phantom für ghost erwartet, bekam %+v", f)
	}
}

// TestCheckTargetsUndocumented: Makefile-Regel `secret` ohne Autoritäts-Doku-
// Eintrag ⇒ gate-undocumented (Richtung 2), an Datei:Zeile der Makefile-Regel.
func TestCheckTargetsUndocumented(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\nsecret:\n\techo\n", // secret in Zeile 3
		tgtDoc:      tgtDocTable("build"),
	}
	f := mustCheck(t, files, tgtCfg())
	if len(f) != 1 || f[0].Reason != ReasonGateUndocumented || f[0].Target != "secret" ||
		f[0].Rule != "targets" || f[0].File != tgtMakefile || f[0].Line != 3 {
		t.Fatalf("gate-undocumented für secret (Makefile:3) erwartet, bekam %+v", f)
	}
}

// TestCheckTargetsProseNotCounted (Tabellen-Scoping, tötet die „Prosa zählt"-
// Mutation): ein Phantom-`make X` NUR in Prosa (keine Tabellenzeile) ⇒ **kein**
// gate-phantom. Ohne den ^|-Guard würde `make prose-ghost` extrahiert und (keine
// Regel) als gate-phantom feuern — dann wäre dieser Test rot.
func TestCheckTargetsProseNotCounted(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\n",
		// Doku: `make phantom-in-prose` nur in Prosa, kein Tabelleneintrag.
		tgtDoc: "# Doc\n\nSiehe `make phantom-in-prose` im Fließtext.\n\n| T | Z |\n|---|---|\n| `make build` | x |\n",
	}
	if f := mustCheck(t, files, tgtCfg()); len(f) != 0 {
		t.Fatalf("Prosa-`make X` darf kein gate-phantom auslösen, bekam %+v", f)
	}
}

// TestCheckTargetsColumnZero (tötet die „getrimmt"-Mutation): eine **eingerückte**
// `| make X |`-Zeile (Pipe nicht in Spalte 0) ist keine Tabellenzeile ⇒ **kein**
// gate-phantom. Skript-Parität (`grep -E '^\|'`); mit „getrimmt" würde `make ghost`
// extrahiert und feuern.
func TestCheckTargetsColumnZero(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\n",
		tgtDoc:      "# Doc\n\n| T | Z |\n|---|---|\n| `make build` | x |\n  | `make ghost` | eingerückt |\n",
	}
	if f := mustCheck(t, files, tgtCfg()); len(f) != 0 {
		t.Fatalf("eingerückte Tabellenzeile (Pipe nicht Spalte 0) darf nichts auslösen, bekam %+v", f)
	}
}

// TestCheckTargetsExempt (Boundary): eine Makefile-Regel in exempt-targets, die in
// der Autoritäts-Doku fehlt ⇒ **kein** gate-undocumented (Utility-Ausnahme).
func TestCheckTargetsExempt(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\nclean:\n\techo\n",
		tgtDoc:      tgtDocTable("build"), // clean NICHT dokumentiert
	}
	cfg := tgtCfg()
	cfg.ExemptTargets = []string{"clean"}
	if f := mustCheck(t, files, cfg); len(f) != 0 {
		t.Fatalf("exemptes clean ⇒ kein gate-undocumented, bekam %+v", f)
	}
	// Gegenprobe: ohne exempt ⇒ gate-undocumented für clean.
	if f := mustCheck(t, files, tgtCfg()); len(f) != 1 || f[0].Target != "clean" || f[0].Reason != ReasonGateUndocumented {
		t.Fatalf("ohne exempt ⇒ gate-undocumented für clean erwartet, bekam %+v", f)
	}
}

// TestCheckTargetsMakefileHeuristic (Skript-Parität): Zuweisungen (`X :=`/`X ?=`),
// `.PHONY`/`.DEFAULT_GOAL` und Pattern-Rules (`%.o:`) sind **keine** Regeln;
// Mehrfach-Target-Zeilen liefern beide Namen. Geprüft über Richtung 2: nur die
// echten, undokumentierten Regeln feuern.
func TestCheckTargetsMakefileHeuristic(t *testing.T) {
	mk := "X := val\n" +
		"Y ?= val2\n" +
		".PHONY: phony-prereq\n" +
		".DEFAULT_GOAL := build\n" +
		"%.o: %.c\n\tcc\n" +
		"build test: dep\n\techo\n" +
		"lint:\n\techo\n"
	files := map[string]string{tgtMakefile: mk, tgtDoc: tgtDocTable("build", "test")} // lint fehlt
	f := mustCheck(t, files, tgtCfg())
	// Genau `lint` ist eine echte, undokumentierte Regel; X/Y/.PHONY/.DEFAULT_GOAL/%.o
	// und der .PHONY-Prereq `phony-prereq` sind KEINE Regeln.
	if len(f) != 1 || f[0].Target != "lint" || f[0].Reason != ReasonGateUndocumented {
		t.Fatalf("nur lint als undokumentierte Regel erwartet, bekam %+v", f)
	}
}

// TestCheckTargetsFailClosed: aktives Modul mit fehlender konfigurierter Datei
// ⇒ error (Exit 2), kein stilles Grün.
func TestCheckTargetsFailClosed(t *testing.T) {
	// fehlendes Makefile
	if _, err := CheckTargets(coretest.NewMemFS(map[string]string{tgtDoc: tgtDocTable("build")}), tgtCfg()); err == nil {
		t.Fatal("fehlendes Makefile ⇒ error erwartet (fail-closed)")
	}
	// fehlende Doku-Datei
	if _, err := CheckTargets(coretest.NewMemFS(map[string]string{tgtMakefile: "build:\n\techo\n"}), tgtCfg()); err == nil {
		t.Fatal("fehlende Doku-Datei ⇒ error erwartet (fail-closed)")
	}
}

// TestCheckTargetsInert: leeres makefiles bzw. nil-fs ⇒ inert (kein Befund, kein
// Fehler, byte-identisch).
func TestCheckTargetsInert(t *testing.T) {
	f, err := CheckTargets(coretest.NewMemFS(map[string]string{tgtMakefile: "build:\n"}), model.TargetsConfig{})
	if err != nil || f != nil {
		t.Fatalf("leeres makefiles ⇒ inert, bekam f=%+v err=%v", f, err)
	}
	if f, err := CheckTargets(nil, tgtCfg()); err != nil || f != nil {
		t.Fatalf("nil-fs ⇒ inert, bekam f=%+v err=%v", f, err)
	}
}

// TestCheckTargetsDirectionDecoupling: leeres doc-tables ⇒ nur Richtung 2; leere
// authority ⇒ nur Richtung 1 (die Richtungen sind unabhängig).
func TestCheckTargetsDirectionDecoupling(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\nsecret:\n\techo\n",
		tgtDoc:      tgtDocTable("build", "ghost"), // ghost=Phantom, secret=undokumentiert
	}
	// Nur Richtung 2 (kein doc-tables): nur gate-undocumented (secret), kein Phantom.
	cfg2 := model.TargetsConfig{Makefiles: []string{tgtMakefile}, Authority: tgtDoc}
	f := mustCheck(t, files, cfg2)
	if len(f) != 1 || f[0].Reason != ReasonGateUndocumented || f[0].Target != "secret" {
		t.Fatalf("nur Richtung 2 erwartet (secret undokumentiert), bekam %+v", f)
	}
	// Nur Richtung 1 (keine authority): nur gate-phantom (ghost), kein undocumented.
	cfg1 := model.TargetsConfig{Makefiles: []string{tgtMakefile}, DocTables: []string{tgtDoc}}
	f = mustCheck(t, files, cfg1)
	if len(f) != 1 || f[0].Reason != ReasonGatePhantom || f[0].Target != "ghost" {
		t.Fatalf("nur Richtung 1 erwartet (ghost phantom), bekam %+v", f)
	}
}
