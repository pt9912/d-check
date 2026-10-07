package rules

import (
	"fmt"
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
		Makefiles: []string{tgtMakefile}, DocTables: []string{tgtDoc}, Authority: []string{tgtDoc},
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
	cfg2 := model.TargetsConfig{Makefiles: []string{tgtMakefile}, Authority: []string{tgtDoc}}
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

// DC-FA-TGT-001: eine Tabellenzeile IM Fence dokumentiert nichts. Ein
// Beispiel-Block, der eine Doku-Tabelle zeigt, liess sonst ein undokumentiertes
// Target als dokumentiert gelten — gate-undocumented entfiel STILL.
func TestTargetsTabellenzeileImFenceDokumentiertNicht(t *testing.T) {
	doc := "# Doc\n\nSo sieht eine Ziel-Tabelle aus:\n\n" +
		"```markdown\n| Target | Zweck |\n|---|---|\n| `make geist` | Beispiel |\n```\n\n" +
		"| Target | Zweck |\n|---|---|\n| `make build` | x |\n"
	files := map[string]string{
		tgtMakefile: "build:\n\techo\n\ngeist:\n\techo\n",
		tgtDoc:      doc,
	}
	got := mustCheck(t, files, tgtCfg())
	var undoc int
	for _, f := range got {
		if f.Reason == ReasonGateUndocumented {
			undoc++
		}
	}
	if undoc != 1 {
		t.Fatalf("die Makefile-Regel geist ist NICHT dokumentiert (nur im Fence) → 1 gate-undocumented, got %d (%+v)", undoc, got)
	}
}

// TestCheckTargetsMakefileGlob: ein Glob-Eintrag expandiert gegen die
// Repo-Wurzel; eine Regel in einer per Glob erfassten Datei ohne Doku-Zeile
// meldet gate-undocumented an DIESER Datei; eine wörtlich UND per Glob
// erfasste Datei zählt einmal (genau ein Befund je Regelzeile); `*` bleibt
// in seinem Segment (das Unterverzeichnis trifft es nicht).
func TestCheckTargetsMakefileGlob(t *testing.T) {
	files := map[string]string{
		tgtMakefile:              "help:\n\techo\n",
		"harness/mk/a.mk":        "alpha:\n\techo\ndelta:\n\techo\n",
		"harness/mk/b.mk":        "beta:\n\techo\n",
		"harness/mk/sub/c.mk":    "gamma:\n\techo\n",
		"harness/mk/readme.txt":  "omega:\n",
		tgtDoc:                   tgtDocTable("help", "alpha", "beta"),
	}
	cfg := tgtCfg()
	cfg.Makefiles = []string{tgtMakefile, "harness/mk/*.mk", "harness/mk/a.mk"}
	f := mustCheck(t, files, cfg)
	if len(f) != 1 || f[0].Reason != ReasonGateUndocumented || f[0].Target != "delta" ||
		f[0].File != "harness/mk/a.mk" || f[0].Line != 3 {
		t.Fatalf("erwartet genau gate-undocumented delta @ harness/mk/a.mk:3, bekam %+v", f)
	}
}

// TestCheckTargetsMakefileGlobDoppelstern: `**` erfasst Fragmente in
// Unterverzeichnissen.
func TestCheckTargetsMakefileGlobDoppelstern(t *testing.T) {
	files := map[string]string{
		tgtMakefile:           "help:\n\techo\n",
		"harness/mk/sub/c.mk": "gamma:\n\techo\n",
		tgtDoc:                tgtDocTable("help"),
	}
	cfg := tgtCfg()
	cfg.Makefiles = []string{tgtMakefile, "harness/**/*.mk"}
	f := mustCheck(t, files, cfg)
	if len(f) != 1 || f[0].Target != "gamma" || f[0].File != "harness/mk/sub/c.mk" {
		t.Fatalf("erwartet gate-undocumented gamma @ harness/mk/sub/c.mk, bekam %+v", f)
	}
}

// TestCheckTargetsMakefileGlobLeer: ein Glob ohne Treffer ist kein stilles
// Grün — fail-closed wie eine fehlende wörtliche Datei, auch wenn das
// Verzeichnis existiert oder fehlt.
func TestCheckTargetsMakefileGlobLeer(t *testing.T) {
	for _, files := range []map[string]string{
		{tgtMakefile: "help:\n\techo\n", tgtDoc: tgtDocTable("help")},
		{tgtMakefile: "help:\n\techo\n", "harness/mk/x.txt": "x", tgtDoc: tgtDocTable("help")},
	} {
		cfg := tgtCfg()
		cfg.Makefiles = []string{tgtMakefile, "harness/mk/*.mk"}
		_, err := CheckTargets(coretest.NewMemFS(files), cfg)
		if err == nil || !strings.Contains(err.Error(), "harness/mk/*.mk") || !strings.Contains(err.Error(), "keine Datei") {
			t.Fatalf("leerer Glob ⇒ fail-closed mit Muster und „keine Datei“ erwartet, bekam %v", err)
		}
	}
}

// TestCheckTargetsMakefileGlobGrenzen: die Grenzen aus Schritt 1a, je Fall
// die beobachtbare Folge — gemeldete Fundstellen (Datei:Ziel) oder ein
// fail-closed-Fehler mit dem erwarteten Wortlaut-Fragment.
func TestCheckTargetsMakefileGlobGrenzen(t *testing.T) {
	cases := []struct {
		name      string
		files     map[string]string
		symlinks  []string
		makefiles []string
		want      []string // "datei:ziel" je Befund
		wantErr   string
	}{
		{name: "SKIP_DIR unter dem Präfix wird nicht betreten",
			files:     map[string]string{"harness/build/x.mk": "skipped:\n", "harness/mk/a.mk": "alpha:\n"},
			makefiles: []string{"harness/**/*.mk"},
			want:      []string{"harness/mk/a.mk:alpha"}},
		{name: "SKIP_DIR im Muster unterhalb des Präfixes nennt den Grund",
			files:     map[string]string{"harness/x/build/a.mk": "alpha:\n"},
			makefiles: []string{"harness/*/build/*.mk"},
			wantErr:   "\"build\" unterhalb des Präfixes wird grundsätzlich übersprungen"},
		{name: "SKIP_DIR-Name im festen Präfix wird betreten",
			files:     map[string]string{"build/x.mk": "built:\n"},
			makefiles: []string{"build/*.mk"},
			want:      []string{"build/x.mk:built"}},
		{name: "Symlink im Präfix wird nicht verfolgt",
			files:     map[string]string{"harness/mk/a.mk": "alpha:\n"},
			symlinks:  []string{"harness"},
			makefiles: []string{"harness/mk/*.mk"},
			wantErr:   "keine Datei"},
		{name: "Symlink-Treffer ist laut",
			files:     map[string]string{"harness/mk/a.mk": "alpha:\n"},
			symlinks:  []string{"harness/mk/link.mk"},
			makefiles: []string{"harness/mk/*.mk"},
			wantErr:   "Symlink"},
		{name: "Präfix ist eine Datei",
			files:     map[string]string{"harness": "x"},
			makefiles: []string{"harness/*.mk"},
			wantErr:   "keine Datei"},
		{name: "Glob im ersten Segment",
			files:     map[string]string{"x.mk": "root:\n", "sub/y.mk": "nested:\n"},
			makefiles: []string{"*.mk"},
			want:      []string{"x.mk:root"}},
		{name: "Fragezeichen und Klasse",
			files:     map[string]string{"harness/mk/a.mk": "alpha:\n", "harness/mk/b.mk": "beta:\n", "harness/mk/cc.mk": "gamma:\n"},
			makefiles: []string{"harness/mk/?.mk", "harness/mk/[ab].mk"},
			want:      []string{"harness/mk/a.mk:alpha", "harness/mk/b.mk:beta"}},
		{name: "Dublette aus zwei Globs",
			files:     map[string]string{"harness/mk/a.mk": "alpha:\n"},
			makefiles: []string{"harness/mk/*.mk", "harness/mk/a*.mk"},
			want:      []string{"harness/mk/a.mk:alpha"}},
		{name: "Dublette über bereinigten Pfad, Form der ersten Nennung",
			files:     map[string]string{"Makefile": "top:\n"},
			makefiles: []string{"./Makefile", "Make*"},
			want:      []string{"./Makefile:top"}},
		{name: "Dublette über bereinigten Pfad, umgekehrte Reihenfolge",
			files:     map[string]string{"Makefile": "top:\n"},
			makefiles: []string{"Make*", "./Makefile"},
			want:      []string{"Makefile:top"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			files := map[string]string{tgtDoc: tgtDocTable()}
			for k, v := range c.files {
				files[k] = v
			}
			m := coretest.NewMemFS(files)
			for _, s := range c.symlinks {
				m.AddSymlink(s)
			}
			cfg := tgtCfg()
			cfg.Makefiles = c.makefiles
			f, err := CheckTargets(m, cfg)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("fail-closed mit %q erwartet, bekam err=%v, befunde=%+v", c.wantErr, err, f)
				}
				return
			}
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			var got []string
			for _, x := range f {
				got = append(got, x.File+":"+x.Target)
			}
			if strings.Join(got, ",") != strings.Join(c.want, ",") {
				t.Fatalf("Befunde = %q, want %q", got, c.want)
			}
		})
	}
}

// TestCheckTargetsAuthorityListe: gate-undocumented misst gegen die
// Vereinigung aller Autoritäts-Dateien — ein Target in einer davon gilt als
// dokumentiert, eines in beiden ergibt keinen Befund, eines in keiner meldet
// an der Regelzeile im Fragment und nennt alle Autoritäts-Dateien.
func TestCheckTargetsAuthorityListe(t *testing.T) {
	files := map[string]string{
		tgtMakefile:       "own:\n\techo\n",
		"harness/mk/a.mk": "tool:\n\techo\nboth:\n\techo\nghost:\n\techo\n",
		"harness/README.md":  tgtDocTable("own", "both"),
		"harness/targets.md": tgtDocTable("tool", "both"),
	}
	cfg := model.TargetsConfig{
		Makefiles: []string{tgtMakefile, "harness/mk/a.mk"},
		Authority: []string{"harness/README.md", "harness/targets.md"},
	}
	f := mustCheck(t, files, cfg)
	if len(f) != 1 || f[0].Reason != ReasonGateUndocumented || f[0].Target != "ghost" ||
		f[0].File != "harness/mk/a.mk" || f[0].Line != 5 {
		t.Fatalf("genau gate-undocumented ghost @ harness/mk/a.mk:5 erwartet, bekam %+v", f)
	}
	want := "Makefile-Regel `ghost` ohne Deklaration in einer der Autoritäts-Dokus harness/README.md, harness/targets.md"
	if f[0].Message != want {
		t.Fatalf("Message = %q\nwant      %q", f[0].Message, want)
	}
}

// TestCheckTargetsAuthorityEinzelnWortlaut: eine einzelne Autoritäts-Datei
// behält den bisherigen Meldungstext (byte-identisch zur String-Form).
func TestCheckTargetsAuthorityEinzelnWortlaut(t *testing.T) {
	files := map[string]string{
		tgtMakefile: "build:\n\techo\nsecret:\n\techo\n",
		tgtDoc:      tgtDocTable("build"),
	}
	f := mustCheck(t, files, tgtCfg())
	want := "Makefile-Regel `secret` ohne Deklaration in der Autoritäts-Doku " + tgtDoc
	if len(f) != 1 || f[0].Message != want {
		t.Fatalf("Message = %+v\nwant %q", f, want)
	}
}

// TestCheckTargetsAuthorityListeFehlend: ein fehlender Listeneintrag ist
// fail-closed, auch wenn der andere existiert.
func TestCheckTargetsAuthorityListeFehlend(t *testing.T) {
	files := map[string]string{
		tgtMakefile:         "own:\n\techo\n",
		"harness/README.md": tgtDocTable("own"),
	}
	cfg := model.TargetsConfig{
		Makefiles: []string{tgtMakefile},
		Authority: []string{"harness/README.md", "harness/targets.md"},
	}
	if _, err := CheckTargets(coretest.NewMemFS(files), cfg); err == nil ||
		!strings.Contains(err.Error(), "harness/targets.md") {
		t.Fatalf("fehlende Autoritäts-Datei ⇒ fail-closed mit ihrem Namen erwartet, bekam %v", err)
	}
}

// TestCheckTargetsAuthorityDublette: dieselbe Datei in zwei Schreibweisen ist
// EINE Autoritäts-Datei — gelesen einmal, gemeldet im Singular-Wortlaut der
// String-Form.
func TestCheckTargetsAuthorityDublette(t *testing.T) {
	files := map[string]string{
		tgtMakefile:         "own:\n\techo\nghost:\n\techo\n",
		"harness/README.md": tgtDocTable("own"),
	}
	cfg := model.TargetsConfig{
		Makefiles: []string{tgtMakefile},
		Authority: []string{"harness/README.md", "./harness/README.md"},
	}
	f := mustCheck(t, files, cfg)
	want := "Makefile-Regel `ghost` ohne Deklaration in der Autoritäts-Doku harness/README.md"
	if len(f) != 1 || f[0].Message != want {
		t.Fatalf("Message = %+v\nwant %q", f, want)
	}
}

// TestCheckTargetsAuthorityDisjunkt: mit authority-disjoint meldet ein
// Target, das in mehr als einer Autoritäts-Datei als Tabellenzeile steht,
// gate-declared-twice an der Zeile jeder weiteren Datei und nennt die Datei
// der ersten Nennung; exempt-targets nimmt davon nicht aus.
func TestCheckTargetsAuthorityDisjunkt(t *testing.T) {
	files := map[string]string{
		tgtMakefile:          "own:\n\techo\ntool:\n\techo\nboth:\n\techo\nclean:\n\techo\n",
		"harness/README.md":  tgtDocTable("own", "both", "clean"),
		"harness/targets.md": tgtDocTable("tool", "both", "clean"),
		"harness/third.md":   tgtDocTable("both"),
	}
	cfg := model.TargetsConfig{
		Makefiles:         []string{tgtMakefile},
		Authority:         []string{"harness/README.md", "harness/targets.md", "harness/third.md"},
		ExemptTargets:     []string{"clean"},
		AuthorityDisjoint: true,
	}
	var got []string
	for _, f := range mustCheck(t, files, cfg) {
		if f.Reason != ReasonGateDeclaredTwice {
			t.Fatalf("unerwarteter Befund %+v", f)
		}
		got = append(got, fmt.Sprintf("%s:%d:%s|%s", f.File, f.Line, f.Target, f.Message))
	}
	want := []string{
		"harness/targets.md:8:both|Target `both` steht in mehr als einer Autoritäts-Doku (zuerst in harness/README.md)",
		"harness/targets.md:9:clean|Target `clean` steht in mehr als einer Autoritäts-Doku (zuerst in harness/README.md)",
		"harness/third.md:7:both|Target `both` steht in mehr als einer Autoritäts-Doku (zuerst in harness/README.md)",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("Befunde =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// TestCheckTargetsAuthorityDisjunktGrenzen: ohne Schalter still (bisheriges
// Verhalten); eine Doppelung innerhalb EINER Datei und dieselbe Datei in zwei
// Schreibweisen sind keine zwei Teile.
func TestCheckTargetsAuthorityDisjunktGrenzen(t *testing.T) {
	files := map[string]string{
		tgtMakefile:          "both:\n\techo\n",
		"harness/README.md":  tgtDocTable("both", "both"),
		"harness/targets.md": tgtDocTable("both"),
	}
	cases := []model.TargetsConfig{
		{Makefiles: []string{tgtMakefile}, Authority: []string{"harness/README.md", "harness/targets.md"}},
		{Makefiles: []string{tgtMakefile}, Authority: []string{"harness/README.md"}, AuthorityDisjoint: true},
		{Makefiles: []string{tgtMakefile}, Authority: []string{"harness/README.md", "./harness/README.md"}, AuthorityDisjoint: true},
	}
	for i, cfg := range cases {
		if f := mustCheck(t, files, cfg); len(f) != 0 {
			t.Errorf("Fall %d: kein Befund erwartet, bekam %+v", i, f)
		}
	}
}
