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
		"docs/traceability.md":  fwdDoc(fwdRows),
		"spec/architecture.md":  bwdDoc(bwdRows),
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
