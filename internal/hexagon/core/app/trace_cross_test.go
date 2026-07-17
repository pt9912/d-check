package app

import (
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Muster der Fremd-Konvention grid-gym (Trigger 088): Anforderungen `GG-<FAM>-NNN`,
// Design-Artefakte `GG-AR-*`. Die Namensräume sind disjunkt — Vorbedingung des
// Mengen-Diffs (DC-FA-XREF-001).
var (
	ggReqPat    = regexp.MustCompile(`GG-[A-Z]+-\d{3}`)
	ggDesignPat = regexp.MustCompile(`GG-AR-[A-Z0-9-]+`)
)

// crossCfg baut eine vollständige Abgleich-Konfiguration; die Tests variieren
// gezielt einzelne Felder.
func crossCfg() *model.TraceCrossConsistency {
	return &model.TraceCrossConsistency{
		Forward: model.TraceCrossForward{
			File:          "docs/traceability.md",
			ReqColumn:     "Anforderung",
			DesignColumn:  "Design-Artefakte",
			DesignPattern: ggDesignPat,
			Ranges:        true,
		},
		Backward: model.TraceCrossBackward{
			File:             "spec/architecture.md",
			ArtifactIDColumn: model.TraceCrossArtifactFirst,
			EdgeColumn:       "Bezug",
			ReqPattern:       ggReqPat,
			Ranges:           true,
		},
		Mode: model.TraceCrossModeEqual,
	}
}

// fwdDoc/bwdDoc bauen die beiden Sicht-Quellen aus je einer Tabellenzeile.
func fwdDoc(rows string) string {
	return "# Traceability\n\n## 27.1 Anforderung zu Design\n\n" +
		"| Anforderung | Design-Artefakte |\n|---|---|\n" + rows
}

func bwdDoc(rows string) string {
	return "# Architektur\n\n## 4 Komponenten\n\n" +
		"| Komponente | Bezug |\n|---|---|\n" + rows
}

func crossFS(fwdRows, bwdRows string) *coretest.MemFS {
	return coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwdDoc(fwdRows),
		"spec/architecture.md": bwdDoc(bwdRows),
	})
}

// Happy Path (DC-FA-XREF-001): eine Anforderung nennt vorwärts fünf Artefakte,
// dieselben fünf nennen sie rückwärts — keine Differenz.
func TestCrossConsistencyConsistent1N(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-P-005, GG-AR-P-009, GG-AR-COMP-CORE, GG-AR-COMP-SCHED, GG-AR-COMP-DOMAIN |\n",
		"| GG-AR-P-005 | GG-ARCH-006 |\n| GG-AR-P-009 | GG-ARCH-006 |\n| GG-AR-COMP-CORE | GG-ARCH-006 |\n"+
			"| GG-AR-COMP-SCHED | GG-ARCH-006 |\n| GG-AR-COMP-DOMAIN | GG-ARCH-006 |\n",
	)
	got, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("konsistentes 1:N ergab Differenzen: %+v", got)
	}
}

// Boundary (DC-FA-XREF-001): F(R)={A,B}, B(R)={B,C} ⇒ genau zwei Befunde, je mit
// Richtungslabel und Quell-Fundstelle. Zugleich der reale grid-gym-Drift-Typ.
func TestCrossConsistencyBothDirections(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN |\n",
		"| GG-AR-COMP-DOMAIN | GG-ARCH-006 |\n| GG-AR-COMP-SCHED | GG-ARCH-006 |\n",
	)
	got, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	want := []CrossFinding{
		{Requirement: "GG-ARCH-006", Artifact: "GG-AR-COMP-CORE", Direction: CrossDirForwardOnly, File: "docs/traceability.md", Line: 7},
		{Requirement: "GG-ARCH-006", Artifact: "GG-AR-COMP-SCHED", Direction: CrossDirBackwardOnly, File: "spec/architecture.md", Line: 8},
	}
	if len(got) != len(want) {
		t.Fatalf("erwartete %d Befunde, got %+v", len(want), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("Befund %d = %+v, want %+v", i, got[i], w)
		}
	}
}

// Range-aware (DC-FA-XREF-001): eine Rück-Kante `..009` und eine Vorwärts-Zeile
// `..004` werden beide expandiert und je Einzel-ID verglichen — GG-SIM-005..009
// bleiben als Rück-Kante ohne RTM-Eintrag.
func TestCrossConsistencyRangeAware(t *testing.T) {
	fs := crossFS(
		"| GG-SIM-001..004 | GG-AR-COMP-SIM |\n",
		"| GG-AR-COMP-SIM | GG-SIM-001..009 |\n",
	)
	got, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	// 001..004 sind beidseitig bekannt; 005..009 nur rückwärts.
	if len(got) != 5 {
		t.Fatalf("erwartete 5 Differenzen (GG-SIM-005..009), got %+v", got)
	}
	for _, f := range got {
		if f.Direction != CrossDirBackwardOnly {
			t.Fatalf("unerwartete Richtung: %+v", f)
		}
	}
	if got[0].Requirement != "GG-SIM-005" || got[4].Requirement != "GG-SIM-009" {
		t.Fatalf("Sortierung/Expansion falsch: %+v", got)
	}
}

// Link-transparente Range über den Kreuzverweis-Abgleich (DC-FA-XREF-001,
// ADR-0039 „Ein Fix, zwei Konsumenten"): die Vorwärts-/Rück-Sicht speist ihre
// Zellen über rangeAwareIDs (row.cells) — ein anderer Vorlauf als der Prosa-Text
// der Coverage-Quelle; diese Achse verriegelt kein Coverage-Test (slice-073
// R2-F-1). Zwei Achsen, beide gegen die Mutation von LinkSuffixEnd gehärtet:
//   - „verlinkt": das Link-Suffix darf eine echte Range dahinter nicht brechen.
//   - „Klammer-URL": ein geklammertes Ziel, dessen Pfadsegmente wie ein Enum
//     aussehen, darf NICHT expandieren — sonst injizierte eine naive „bis zur
//     ersten `)`"-Abgrenzung falsche Rück-Kanten (die stille Richtung aus R1-F-1).
//     Asymmetrisch nur in der Rück-Sicht, damit die Falsch-Expansion als
//     Differenz sichtbar wird statt sich beidseitig aufzuheben.
func TestCrossConsistencyRangeAwareLinkTransparent(t *testing.T) {
	for _, tc := range []struct {
		name             string
		forward, backward string
		wantDiffs        int
	}{
		{
			"verlinkt, echte Range dahinter",
			"| [`GG-SIM-001`](../spec/x.md)..004 | GG-AR-COMP-SIM |\n",
			"| GG-AR-COMP-SIM | [`GG-SIM-001`](../spec/x.md)..009 |\n",
			5, // 001..004 beidseitig, 005..009 nur rückwärts
		},
		{
			"Klammer-URL, Pfadsegmente sind kein Enum",
			"| GG-SIM-001 | GG-AR-COMP-SIM |\n",
			"| GG-AR-COMP-SIM | [`GG-SIM-001`](../a(2)/002/003.md) |\n",
			0, // beide Sichten kennen nur 001; /002/003 ist URL-intern, kein Enum
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := crossConsistency(crossFS(tc.forward, tc.backward), crossCfg(), ggReqPat)
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if len(got) != tc.wantDiffs {
				t.Fatalf("Differenzen = %d, want %d: %+v", len(got), tc.wantDiffs, got)
			}
		})
	}
}

// Superset-Modus (DC-FA-XREF-001): F ⊋ B ⇒ F\B ist kein Befund.
func TestCrossConsistencySupersetMode(t *testing.T) {
	fwd := "| GG-ARCH-006 | GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN |\n"
	bwd := "| GG-AR-COMP-CORE | GG-ARCH-006 |\n"
	cfg := crossCfg()
	cfg.Mode = model.TraceCrossModeSuperset
	got, err := crossConsistency(crossFS(fwd, bwd), cfg, ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("superset meldete F\\B: %+v", got)
	}
	// Gegenprobe: derselbe Stand gatet unter equal.
	equalGot, err := crossConsistency(crossFS(fwd, bwd), crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(equalGot) != 1 || equalGot[0].Direction != CrossDirForwardOnly {
		t.Fatalf("equal-Gegenprobe: %+v", equalGot)
	}
}

// Ableitungssprung (DC-FA-XREF-001): eine Rück-Kante auf eine Mittelschicht-ID
// fällt per exclude-req aus dem Vergleich — weder Differenz noch Waise.
func TestCrossConsistencyExcludeReq(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE |\n",
		"| GG-AR-COMP-CORE | GG-ARCH-006 |\n| GG-AR-COMP-MID | GG-SPEC-042 |\n",
	)
	cfg := crossCfg()
	cfg.ExcludeReq = regexp.MustCompile(`^GG-SPEC-`)
	got, err := crossConsistency(fs, cfg, ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("exclude-req griff nicht: %+v", got)
	}
	// Gegenprobe: ohne Ventil ist die Mittelschicht-Kante eine Differenz.
	got2, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got2) != 1 || got2[0].Requirement != "GG-SPEC-042" {
		t.Fatalf("Gegenprobe ohne exclude-req: %+v", got2)
	}
}

// Abschnitts-Span (DC-FA-XREF-001): die Whitelist scopt die Vorwärts-Sicht; die
// Zeilennummern der Befunde bleiben die der Original-Datei.
func TestCrossConsistencySections(t *testing.T) {
	fwd := fwdDoc("| GG-ARCH-006 | GG-AR-COMP-CORE |\n") +
		"\n## 27.1.1 Anforderungen ohne Design-Artefakt\n\n" +
		"| Anforderung | Design-Artefakte |\n|---|---|\n| GG-ARCH-099 | GG-AR-COMP-GHOST |\n"
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwd,
		"spec/architecture.md": bwdDoc("| GG-AR-COMP-CORE | GG-ARCH-006 |\n"),
	})
	cfg := crossCfg()
	cfg.Forward.ExcludeSections = []string{"27.1.1 Anforderungen ohne Design-Artefakt"}
	got, err := crossConsistency(fs, cfg, ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("ausgeschlossener Abschnitt zählte mit: %+v", got)
	}
	// Ohne Blacklist ist die Geister-Zeile eine Differenz — mit ihrer echten Zeile.
	got2, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got2) != 1 || got2[0].Artifact != "GG-AR-COMP-GHOST" {
		t.Fatalf("Gegenprobe: %+v", got2)
	}
	if want := strings.Count(fwd[:strings.Index(fwd, "| GG-ARCH-099")], "\n") + 1; got2[0].Line != want {
		t.Fatalf("Zeilennummer %d, want %d (Original-Datei, nicht gefilterter Text)", got2[0].Line, want)
	}
}

// Negative (DC-FA-XREF-001): die Config-/Quell-Fehler sind fail-closed — jeder
// beendet den Lauf, statt still eine leere Sicht zu vergleichen.
func TestCrossConsistencyFailClosed(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*model.TraceCrossConsistency)
		fs   *coretest.MemFS
		want string
	}{
		{
			name: "fehlende Vorwärts-Datei",
			mut:  func(c *model.TraceCrossConsistency) { c.Forward.File = "docs/fehlt.md" },
			want: "fehlt",
		},
		{
			name: "fehlende Rückwärts-Datei",
			mut:  func(c *model.TraceCrossConsistency) { c.Backward.File = "spec/fehlt.md" },
			want: "fehlt",
		},
		{
			name: "Vorwärts-Spalte trifft keine Tabelle",
			mut:  func(c *model.TraceCrossConsistency) { c.Forward.ReqColumn = "Kennung" },
			want: "keine Tabelle",
		},
		{
			name: "Kanten-Spalte trifft keine Tabelle",
			mut:  func(c *model.TraceCrossConsistency) { c.Backward.EdgeColumn = "Verweis" },
			want: "keine Tabelle",
		},
		{
			name: "Abschnittsname ohne Heading-Treffer",
			mut:  func(c *model.TraceCrossConsistency) { c.Forward.Sections = []string{"27.1"} },
			want: "trifft keine Überschrift",
		},
		{
			name: "ungültige Range in der Kanten-Zelle",
			fs:   crossFS("| GG-SIM-001 | GG-AR-COMP-SIM |\n", "| GG-AR-COMP-SIM | GG-SIM-009..003 |\n"),
			want: "AAA>BBB",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := tc.fs
			if fs == nil {
				fs = crossFS("| GG-ARCH-006 | GG-AR-COMP-CORE |\n", "| GG-AR-COMP-CORE | GG-ARCH-006 |\n")
			}
			cfg := crossCfg()
			if tc.mut != nil {
				tc.mut(cfg)
			}
			_, err := crossConsistency(fs, cfg, ggReqPat)
			if err == nil {
				t.Fatal("erwartete Fehler (fail-closed), got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Fehlertext %q enthält %q nicht", err.Error(), tc.want)
			}
		})
	}
}

// Doppelter Rollen-Header ist fail-closed — die Spaltenbindung wäre sonst geraten
// (DC-FA-XREF-001: ID-Header nicht genau einmal ⇒ Exit 2).
func TestCrossConsistencyDuplicateHeader(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": "# T\n\n| Anforderung | Design-Artefakte | Anforderung |\n|---|---|---|\n" +
			"| GG-ARCH-006 | GG-AR-COMP-CORE | GG-ARCH-006 |\n",
		"spec/architecture.md": bwdDoc("| GG-AR-COMP-CORE | GG-ARCH-006 |\n"),
	})
	_, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err == nil || !strings.Contains(err.Error(), "mehrfach") {
		t.Fatalf("erwartete Mehrfach-Header-Fehler, got %v", err)
	}
}

// Default byte-identisch (DC-QA-02): ohne trace.cross-consistency-Block bleibt die
// RTM unverändert und trägt weder Befunde noch den Abgleich-Abschnitt.
func TestBuildTraceMatrixWithoutCrossBlock(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md": "# L\n\n### DC-FA-LINK-001 — Links\n\nText.\n",
	})
	m, err := BuildTraceMatrix(fs, model.TraceConfig{})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if m.CrossActive {
		t.Fatal("CrossActive ohne Block gesetzt")
	}
	if m.CrossConsistency != nil {
		t.Fatalf("Befunde ohne Block: %+v", m.CrossConsistency)
	}
}

// Vakuum a (DC-FA-XREF-001.a Schritt 5): ein design-pattern, das kompiliert, aber
// am Artefakt-Namensraum vorbeigreift, räumt BEIDE Sichten leer (es ist geteilt) —
// `0 Differenz(en)`/Exit 0 behauptete eine nie geprüfte Konsistenz (R1-F-1).
func TestCrossConsistencyVakuumBeideSichtenLeer(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN |\n",
		"| GG-AR-COMP-SCHED | GG-ARCH-006 |\n",
	)
	cfg := crossCfg()
	// Tippfehler-Klasse: kompiliert, trifft aber kein Artefakt (GG-ARCH- statt GG-AR-).
	cfg.Forward.DesignPattern = regexp.MustCompile(`GG-ARCH-COMP-[A-Z]+`)
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("Namensraum-Fehlgriff ergab stilles Grün statt Exit 2")
	}
	if !strings.Contains(err.Error(), "beide Sichten ergaben 0 Kanten") {
		t.Fatalf("Fehlertext benennt das Vakuum nicht: %v", err)
	}
	// Hier greifen die Muster wirklich vorbei (schon vor jedem Ausschluss) — die
	// Diagnose zeigt dorthin und nicht auf ein Ventil.
	if strings.Contains(err.Error(), "exclude-req") {
		t.Fatalf("Diagnose verdächtigt ein Ventil statt der Muster: %v", err)
	}
	// Der Hint muss auf DIESE Ursache zeigen (design-pattern), nicht auf die der
	// anderen Vakuum-Art — sonst wären die Hints vertauschbar, ohne dass ein Test
	// kippt, und die Meldung schickte in die falsche Config-Ecke.
	if !strings.Contains(err.Error(), "design-pattern") || strings.Contains(err.Error(), "edge-column") {
		t.Fatalf("Hint zeigt nicht auf das design-pattern: %v", err)
	}
}

// Vakuum b (DC-FA-XREF-001.a Schritt 5): unter `superset` gatet allein B\F — eine
// kantenleere Rück-Sicht kann konstruktionsbedingt nie einen Befund melden, egal
// wie voll die Vorwärts-Sicht ist. Mutations-hart getrennt von Vakuum a: hier ist
// F NICHT leer, der (a)-Zweig also nachweislich nicht die Ursache.
func TestCrossConsistencyVakuumRueckSichtLeerUnterSuperset(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE |\n",
		"| GG-AR-COMP-CORE | kein Bezug gepflegt |\n",
	)
	cfg := crossCfg()
	cfg.Mode = model.TraceCrossModeSuperset
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("superset mit kantenleerer Rück-Sicht ergab stilles Grün statt Exit 2")
	}
	if !strings.Contains(err.Error(), "Rück-Sicht ergab 0 Kanten") {
		t.Fatalf("Fehlertext benennt die leere Rück-Sicht nicht: %v", err)
	}
	// Gegenstück zur Hint-Pinnung oben: hier muss die Rück-Kanten-Config genannt
	// sein, nicht das design-pattern.
	if !strings.Contains(err.Error(), "edge-column") || strings.Contains(err.Error(), "design-pattern") {
		t.Fatalf("Hint zeigt nicht auf die Rück-Kanten-Config: %v", err)
	}
	// Gegenprobe: unter `equal` ist derselbe Stand KEIN Vakuum — F\B gatet und
	// meldet laut. Ohne diese Zeile wäre der mode-Zweig oben nicht belegt.
	got, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("equal mit leerer Rück-Sicht darf kein Fehler sein: %v", err)
	}
	if len(got) != 1 || got[0].Direction != CrossDirForwardOnly {
		t.Fatalf("equal-Gegenprobe: %+v", got)
	}
}

// KEIN Vakuum (DC-FA-XREF-001 Boundary): eine einseitig leere VORWÄRTS-Sicht ist
// ein wohldefiniertes Ergebnis, kein Fehler — der Diff läuft über keys(F) ∪ keys(B).
// Das ist der Bootstrap-Zustand (unrestrukturierte RTM, gepflegte Rück-Kanten), den
// ADR-0038 Entscheidung 3 dem Konsumenten aufträgt und Entscheidung 7 als
// Generator-Eingang braucht; ein symmetrischer Guard würgte ihn ab (R2-F-1).
func TestCrossConsistencyLeereVorwaertsSichtIstKeinVakuum(t *testing.T) {
	fs := crossFS(
		// Noch Prosa statt konkreter IDs — genau die Vorarbeit, die aussteht.
		"| GG-ARCH-006 | alle Scheduler-Komponenten (siehe Architektur) |\n",
		"| GG-AR-COMP-SCHED | GG-ARCH-006 |\n| GG-AR-P-005 | GG-ARCH-006 |\n",
	)
	got, err := crossConsistency(fs, crossCfg(), ggReqPat)
	if err != nil {
		t.Fatalf("leere Vorwärts-Sicht darf kein Exit 2 sein (Bootstrap-Zustand): %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("erwartete 2 laute B\\F-Befunde, got %+v", got)
	}
	for _, f := range got {
		if f.Direction != CrossDirBackwardOnly {
			t.Fatalf("unerwartete Richtung: %+v", f)
		}
	}
}

// R1-F-2 (MEDIUM, DC-FA-XREF-001.a Schritt 3): Relevanz entsteht allein über die
// Kanten-Spalte — „zählt jede Tabelle mit einem edge-column-Header". Fehlt in einer
// so relevanten Tabelle der konfigurierte ID-Header, ist das Exit 2; ein stilles
// Überspringen ließe ihre Rück-Kanten lautlos verschwinden.
func TestCrossConsistencyBackwardIDHeaderFehlt(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwdDoc("| GG-ARCH-006 | GG-AR-COMP-CORE |\n"),
		// ZWEI Bezug-Tabellen mit heterogenen ID-Headern (der Realfall aus ADR-0038).
		// Tabelle 1 trägt „Kennung" und deckt die Vorwärts-Sicht; Tabelle 2 trägt ihn
		// nicht. Solange die Relevanz an der ID-Spalte hing, verschwand Tabelle 2 samt
		// ihrer Kante still — B blieb nicht-leer (Tabelle 1), also griff auch kein
		// Vakuitäts-Guard, und `superset` meldete 0 Differenzen/Exit 0, obwohl die
		// echte Kante GG-AR-P-005 → GG-ARCH-006 keinen RTM-Eintrag hat.
		"spec/architecture.md": "# A\n\n## 4 Komponenten\n\n| Kennung | Bezug |\n|---|---|\n" +
			"| GG-AR-COMP-CORE | GG-ARCH-006 |\n\n" +
			"## 5 Ports\n\n| Port-ID | Bezug |\n|---|---|\n" +
			"| GG-AR-P-005 | GG-ARCH-006 |\n",
	})
	cfg := crossCfg()
	cfg.Backward.ArtifactIDColumn = "Kennung"
	cfg.Mode = model.TraceCrossModeSuperset
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("relevante Bezug-Tabelle ohne ID-Header wurde still übersprungen (Kante verloren, Exit 0)")
	}
	if !strings.Contains(err.Error(), "Kennung") || !strings.Contains(err.Error(), "0-mal") {
		t.Fatalf("Fehlertext benennt den fehlenden Header nicht: %v", err)
	}
}

// Gegenpol zu R1-F-2: ein benannter ID-Header wird korrekt gebunden (nicht nur
// der `first`-Sentinel) — sonst wäre der Fehlpfad oben trivial grün.
func TestCrossConsistencyBackwardIDHeaderBenannt(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwdDoc("| GG-ARCH-006 | GG-AR-COMP-CORE |\n"),
		// ID-Spalte NICHT die erste — nur die Header-Bindung findet sie.
		"spec/architecture.md": "# A\n\n## 4 Komponenten\n\n| Schicht | Kennung | Bezug |\n|---|---|---|\n" +
			"| Kern | GG-AR-COMP-CORE | GG-ARCH-006 |\n",
	})
	cfg := crossCfg()
	cfg.Backward.ArtifactIDColumn = "Kennung"
	got, err := crossConsistency(fs, cfg, ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("benannte ID-Spalte nicht gebunden (Differenzen statt Deckung): %+v", got)
	}
}

// Vakuum nach dem Ventil (DC-FA-XREF-001.a Schritte 4–5): `exclude-req` läuft VOR
// der Vakuitäts-Prüfung — ein Ventil, das alle Anforderungen verschluckt, schaltet
// das Gate ebenso still ab wie ein fehlgreifendes Muster. Ohne die Stufenfolge
// meldete echter Drift `0 Differenz(en)`/Exit 0 (R3-F-1).
func TestCrossConsistencyVakuumDurchUebergriffigesVentil(t *testing.T) {
	// Echter Drift: F={CORE}, B={SCHED}, Schnittmenge null.
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE |\n",
		"| GG-AR-COMP-SCHED | GG-ARCH-006 |\n",
	)
	cfg := crossCfg()
	cfg.ExcludeReq = regexp.MustCompile(`.`) // verschluckt jede Anforderung
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("übergriffiges exclude-req schaltete das Gate still ab (Exit 0 trotz Drift)")
	}
	if !strings.Contains(err.Error(), "beide Sichten ergaben 0 Kanten") {
		t.Fatalf("Fehlertext benennt das Vakuum nicht: %v", err)
	}
	// Die Diagnose muss das Ventil benennen — und die Muster NICHT: sie sind hier
	// nachweislich korrekt (vor dem Ausschluss lagen Kanten vor). Eine Meldung, die
	// trotzdem in die Muster-Config schickt, ist eine Fehldiagnose
	// (DC-FA-XREF-001: „die Meldung benennt das Ventil, nicht die korrekten Muster").
	if !strings.Contains(err.Error(), "exclude-req hat jede Anforderung") {
		t.Fatalf("Fehldiagnose: Meldung nennt das übergriffige Ventil nicht: %v", err)
	}
	if strings.Contains(err.Error(), "design-pattern") {
		t.Fatalf("Fehldiagnose: Meldung verdächtigt die korrekten Muster: %v", err)
	}
	// Gegenprobe: ein eng gefasstes Ventil lässt den Drift stehen und meldet ihn.
	cfg.ExcludeReq = regexp.MustCompile(`^GG-SPEC-`)
	got, err := crossConsistency(fs, cfg, ggReqPat)
	if err != nil {
		t.Fatalf("enges Ventil darf kein Vakuum sein: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("erwartete 2 Differenzen trotz Ventil, got %+v", got)
	}
}

// Fehlerpräzedenz (DC-FA-XREF-001.a): die Stufe „Quellen lesen" läuft über BEIDE
// Sichten, bevor die Abschnitts-Spannung beginnt. Sonst verdeckte ein
// Vorwärts-Sektionsfehler die eine Stufe früher liegende fehlende Rück-Datei
// (R2-F-2/R3-F-3).
func TestCrossConsistencyPraezedenzQuelleVorAbschnitt(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwdDoc("| GG-ARCH-006 | GG-AR-COMP-CORE |\n"),
		// spec/architecture.md fehlt absichtlich.
	})
	cfg := crossCfg()
	cfg.Forward.Sections = []string{"Tippfehler-Abschnitt"} // Fehler einer SPÄTEREN Stufe
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("erwartete Fehler")
	}
	if !strings.Contains(err.Error(), "backward.file") {
		t.Fatalf("Vorwärts-Sektionsfehler verdeckte die fehlende Rück-Datei: %v", err)
	}
}

// Doppelter ID-Header ist fail-closed (DC-FA-XREF-001: ID-Header nicht genau
// einmal ⇒ Exit 2) — sonst bände er still an sein erstes Vorkommen (R3-F-4).
func TestCrossConsistencyBackwardIDHeaderDoppelt(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"docs/traceability.md": fwdDoc("| GG-ARCH-006 | GG-AR-COMP-CORE |\n"),
		"spec/architecture.md": "# A\n\n## 4 K\n\n| Kennung | Bezug | Kennung |\n|---|---|---|\n" +
			"| GG-AR-COMP-CORE | GG-ARCH-006 | GG-AR-COMP-ALT |\n",
	})
	cfg := crossCfg()
	cfg.Backward.ArtifactIDColumn = "Kennung"
	_, err := crossConsistency(fs, cfg, ggReqPat)
	if err == nil {
		t.Fatal("doppelter ID-Header band still an das erste Vorkommen")
	}
	if !strings.Contains(err.Error(), "2-mal") {
		t.Fatalf("Fehlertext benennt die Mehrfach-Bindung nicht: %v", err)
	}
}

// slice-071/Defekt 1 (DC-FA-XREF-001): die Vergleichs-Schlüsselmenge ist NICHT die
// RTM-Anforderungsmenge. Scopt ein Repo seine RTM bewusst (Architektur-Meta ist
// keine Anforderung), holt forward.req-pattern die Familie in den Abgleich zurück
// — ohne es verschwindet die F\B-Richtung, und der verbleibende Befund sieht wie
// ein Treffer aus, ist aber ein Nebeneffekt von F = ∅.
func TestCrossConsistencyForwardReqPattern(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-007 | GG-AR-COMP-CORE |\n",
		"| GG-AR-P-006 | GG-ARCH-007 |\n",
	)
	// Die RTM scopt ARCH bewusst aus.
	rtmPat := regexp.MustCompile(`GG-(SIM|UI)-\d{3}`)

	// Ohne forward.req-pattern: F ist leer, nur die Rück-Kante feuert.
	got, err := crossConsistency(fs, crossCfg(), rtmPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 1 || got[0].Direction != CrossDirBackwardOnly {
		t.Fatalf("Vorbedingung des Tests trifft nicht zu: %+v", got)
	}

	// Mit forward.req-pattern: beide Richtungen erscheinen — der echte Drift.
	cfg := crossCfg()
	cfg.Forward.ReqPattern = regexp.MustCompile(`GG-[A-Z]+-\d{3}`)
	got2, err := crossConsistency(fs, cfg, rtmPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got2) != 2 {
		t.Fatalf("erwartete beide Richtungen, got %+v", got2)
	}
	if got2[0].Artifact != "GG-AR-COMP-CORE" || got2[0].Direction != CrossDirForwardOnly {
		t.Fatalf("die F\\B-Richtung fehlt — genau der Defekt: %+v", got2)
	}
	if got2[1].Artifact != "GG-AR-P-006" || got2[1].Direction != CrossDirBackwardOnly {
		t.Fatalf("B\\F-Richtung falsch: %+v", got2)
	}
}

// Der Default hält die Abwärtskompatibilität: ohne forward.req-pattern gilt das
// RTM-Muster (DC-FA-XREF-001) — sonst wäre der neue Schlüssel ein Breaking Change.
func TestCrossConsistencyForwardReqPatternDefault(t *testing.T) {
	fs := crossFS(
		"| GG-ARCH-006 | GG-AR-COMP-CORE |\n",
		"| GG-AR-COMP-CORE | GG-ARCH-006 |\n",
	)
	cfg := crossCfg()
	if cfg.Forward.ReqPattern != nil {
		t.Fatal("Testvorbedingung: ReqPattern muss ungesetzt sein")
	}
	got, err := crossConsistency(fs, cfg, ggReqPat)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Default-Fallback auf das RTM-Muster greift nicht: %+v", got)
	}
}
