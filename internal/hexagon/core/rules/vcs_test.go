package rules

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// fakeVCS ist ein hermetischer driven.VCS-Doppelgänger (kein echtes git) für die
// Akzeptanztests von CheckVCS (DC-FA-VCS-001): er liefert vorgegebene Changes und
// Datei-Inhalte je Ref. err simuliert den fail-closed-Pfad (fehlendes .git/Range).
type fakeVCS struct {
	changes []driven.VCSChange
	files   map[string]map[string][]byte // ref → pfad → inhalt
	commits []driven.CommitMeta          // Modul commits (DC-FA-COMMITS-001)
	tracked map[string]bool              // Modul tracked (DC-FA-TRK-001)
	err     error
}

func (f *fakeVCS) ChangedPaths(_, _ string) ([]driven.VCSChange, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.changes, nil
}

func (f *fakeVCS) CommitMessages(_, _ string) ([]driven.CommitMeta, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.commits, nil
}

func (f *fakeVCS) TrackedPaths() (map[string]bool, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.tracked, nil
}

func (f *fakeVCS) FileAt(ref, path string) ([]byte, bool, error) {
	if m, ok := f.files[ref]; ok {
		if c, ok := m[path]; ok {
			return c, true, nil
		}
	}
	return nil, false, nil
}

// adrConfig ist die dogfood-Klasse: Accepted-ADRs unter docs/plan/adr/, Core
// ohne Geschichte, nur die Kopf-Status-Zeile gestrippt, erlaubter Übergang
// Accepted|Superseded (Skript-Parität, ADR-0024).
func adrConfig() model.VCSConfig {
	return model.VCSConfig{
		Paths:           []string{"docs/plan/adr/[0-9]*.md"},
		ImmutableWhen:   regexp.MustCompile(`^\*\*Status:\*\* Accepted`),
		ExcludeSections: []string{"Geschichte"},
		StatusLine:      regexp.MustCompile(`^\*\*Status:\*\*`),
		HeadAllow:       regexp.MustCompile(`^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})`),
	}
}

const adrPath = "docs/plan/adr/0099-x.md"

// adr baut ein ADR mit gegebenem Status + Entscheidungs-Körper (mit Geschichte).
func adr(status, decision string) []byte {
	return []byte("# ADR-0099 — X\n\n**Status:** " + status +
		"\n**Datum:** 2026-01-01\n\n## Entscheidung\n\n" + decision +
		"\n\n## Geschichte\n\n| Datum | Ereignis |\n|---|---|\n| 2026-01-01 | Proposed → Accepted |\n")
}

func refs(base, head []byte) map[string]map[string][]byte {
	m := map[string]map[string][]byte{"BASE": {}, "HEAD": {}}
	if base != nil {
		m["BASE"][adrPath] = base
	}
	if head != nil {
		m["HEAD"][adrPath] = head
	}
	return m
}

// TestVCSModified deckt die Core-/Status-Semantik einer modifizierten ADR ab —
// die sieben Selbsttest-Klassen des abgelösten Skripts plus die Reflow-Boundary.
func TestVCSModified(t *testing.T) {
	afterGeschichte := func(detail string) []byte {
		return []byte("# ADR-0099 — X\n\n**Status:** Accepted\n\n## Entscheidung\n\nTue A.\n\n" +
			"## Geschichte\n\n| Datum | Ereignis |\n|---|---|\n| 2026-01-01 | Proposed → Accepted |\n\n" +
			"## Anhang\n\n" + detail + "\n")
	}
	bodyStatus := func(tail string) []byte {
		return []byte("# ADR-0099 — X\n\n**Status:** Accepted\n\n## Konsequenzen\n\n" +
			"**Status:** der Migration " + tail + "\n\n" +
			"## Geschichte\n\n| Datum | Ereignis |\n|---|---|\n| 2026-01-01 | Proposed → Accepted |\n")
	}
	cases := []struct {
		name       string
		base, head []byte
		want       int
	}{
		{"geschichte-anhang feuert nicht", adr("Accepted", "Tue A."),
			append(adr("Accepted", "Tue A."), []byte("| 2026-02-02 | Notiz |\n")...), 0},
		{"superseded-uebergang feuert nicht", adr("Accepted", "Tue A."),
			adr("Superseded by ADR-0100", "Tue A."), 0},
		{"reflow am core feuert nicht", adr("Accepted", "Tue A."),
			adr("Accepted", "Tue   A."), 0},
		{"koerper-edit feuert", adr("Accepted", "Tue A."),
			adr("Accepted", "Tue B."), 1},
		{"status-rueckfall feuert", adr("Accepted", "Tue A."),
			adr("Proposed", "Tue A."), 1},
		{"proposed-base ist frei", adr("Proposed", "Tue A."),
			adr("Proposed", "Tue B."), 0},
		{"edit nach geschichte feuert", afterGeschichte("Detail A."),
			afterGeschichte("Detail B."), 1},
		{"koerper-status-zeile feuert", bodyStatus("bleibt offen."),
			bodyStatus("ist geklärt."), 1},
		// Körper-Drift UND unzulässiger Status-Übergang zugleich: CheckVCS liefert
		// beide core-drift-vcs-Befunde roh (2); sie teilen den Sort-Key und werden
		// erst in RunWithVCS via SortFindings auf einen dedupliziert (R2-F-3,
		// dokumentiert) — das Gate feuert so oder so.
		{"koerper+status zugleich (CheckVCS roh: 2)", adr("Accepted", "Tue A."),
			adr("Proposed", "Tue B."), 2},
	}
	cfg := adrConfig()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fv := &fakeVCS{
				changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
				files:   refs(c.base, c.head),
			}
			got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if len(got) != c.want {
				t.Fatalf("Befunde = %d, want %d (%v)", len(got), c.want, got)
			}
			for _, f := range got {
				if f.Reason != model.ReasonCoreDriftVCS || f.Rule != "vcs" || f.File != adrPath {
					t.Fatalf("unerwarteter Befund: %+v", f)
				}
			}
		})
	}
}

// TestVCSNoStatusLine: ohne status-line wird keine Kopf-Status-Zeile gestrippt
// und (ohne head-allow) keine Übergangs-Prüfung gefahren — reiner Core-Vergleich
// (deckt den statusLine==nil-Zweig ab, R2-F-4).
func TestVCSNoStatusLine(t *testing.T) {
	cfg := model.VCSConfig{
		Paths:           []string{"docs/plan/adr/[0-9]*.md"},
		ImmutableWhen:   regexp.MustCompile(`^\*\*Status:\*\* Accepted`),
		ExcludeSections: []string{"Geschichte"},
		// StatusLine + HeadAllow bewusst nil
	}
	drift := &fakeVCS{
		changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
		files:   refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B.")),
	}
	if got, err := CheckVCS(drift, cfg, "BASE", "HEAD"); err != nil || len(got) != 1 || got[0].Reason != model.ReasonCoreDriftVCS {
		t.Fatalf("ohne status-line: ein core-drift-vcs (Body) erwartet, got %v err=%v", got, err)
	}
	clean := &fakeVCS{
		changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
		files:   refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue  A.")),
	}
	if got, err := CheckVCS(clean, cfg, "BASE", "HEAD"); err != nil || len(got) != 0 {
		t.Fatalf("ohne status-line + reflow: kein Befund erwartet, got %v err=%v", got, err)
	}
}

// TestVCSDeleteAddClass deckt Löschung, Hinzufügung und die Klassen-Filterung ab.
func TestVCSDeleteAddClass(t *testing.T) {
	cfg := adrConfig()
	cases := []struct {
		name   string
		change driven.VCSChange
		files  map[string]map[string][]byte
		want   int
	}{
		{"geloeschte accepted feuert", driven.VCSChange{Status: driven.VCSDeleted, Path: adrPath},
			refs(adr("Accepted", "Tue A."), nil), 1},
		{"geloeschte proposed ist frei", driven.VCSChange{Status: driven.VCSDeleted, Path: adrPath},
			refs(adr("Proposed", "Tue A."), nil), 0},
		{"neue datei ist frei", driven.VCSChange{Status: driven.VCSAdded, Path: adrPath},
			refs(nil, adr("Accepted", "Tue A.")), 0},
		{"pfad ausserhalb der klasse ignoriert", driven.VCSChange{Status: driven.VCSModified, Path: "docs/other.md"},
			map[string]map[string][]byte{"BASE": {"docs/other.md": adr("Accepted", "Tue A.")}, "HEAD": {"docs/other.md": adr("Accepted", "Tue B.")}}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fv := &fakeVCS{changes: []driven.VCSChange{c.change}, files: c.files}
			got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if len(got) != c.want {
				t.Fatalf("Befunde = %d, want %d (%v)", len(got), c.want, got)
			}
		})
	}
}

// TestVCSFailClosed: ein Port-Fehler (fehlendes .git/Range) wird als error
// durchgereicht — der Aufrufer mappt auf Exit 2 (DC-FA-VCS-001 fail-closed).
func TestVCSFailClosed(t *testing.T) {
	fv := &fakeVCS{err: errors.New("kein .git")}
	if _, err := CheckVCS(fv, adrConfig(), "BASE", "HEAD"); err == nil {
		t.Fatal("fail-closed erwartet: CheckVCS hätte den Port-Fehler durchreichen müssen")
	}
}

// TestVCSInert: ohne paths-Klasse oder ohne Port ist das Modul wirkungslos.
func TestVCSInert(t *testing.T) {
	fv := &fakeVCS{changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
		files: refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B."))}
	if got, err := CheckVCS(fv, model.VCSConfig{}, "BASE", "HEAD"); err != nil || got != nil {
		t.Fatalf("ohne paths inert erwartet: got=%v err=%v", got, err)
	}
	if got, err := CheckVCS(nil, adrConfig(), "BASE", "HEAD"); err != nil || got != nil {
		t.Fatalf("ohne Port inert erwartet: got=%v err=%v", got, err)
	}
}

// TestVCSDispatch deckt den Post-Pass in RunWithVCS ab: aktiv ⇒ vcs-Befunde
// erscheinen; ohne aktives vcs bleibt ein bereitgestellter Port ungenutzt
// (Modul-aus byte-identisch); ein Port-Fehler bricht den Lauf (Exit 2).
func TestVCSDispatch(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{"docs/x.md": "x"})
	cfg := model.Config{VCS: adrConfig()}
	driftPort := &fakeVCS{
		changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
		files:   refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B.")),
	}

	res, err := RunWithVCS(m, nil, driftPort, "BASE", "HEAD", cfg, []string{"vcs"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != model.ReasonCoreDriftVCS {
		t.Fatalf("aktiv: ein core-drift-vcs erwartet, got %v", res.Findings)
	}

	// Modul-aus: derselbe Port, aber vcs nicht im Modulsatz ⇒ keine vcs-Befunde.
	resOff, err := RunWithVCS(m, nil, driftPort, "BASE", "HEAD", cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	for _, f := range resOff.Findings {
		if f.Reason == model.ReasonCoreDriftVCS {
			t.Fatal("Modul-aus: vcs hätte nicht laufen dürfen")
		}
	}

	// fail-closed über RunWithVCS.
	if _, err := RunWithVCS(m, nil, &fakeVCS{err: errors.New("kein .git")}, "BASE", "HEAD", cfg, []string{"vcs"}); err == nil {
		t.Fatal("fail-closed über RunWithVCS erwartet")
	}
}

// DC-FA-VCS-001: eine Kopfzeile IM Code-Block ist keine. Eine Vorlagen- oder
// Konventionsdatei zeigt ihren eigenen Kopf als Beispiel — vorher galt sie
// dadurch als immutabel bzw. verschob die gestrippte Status-Zeile, und eine
// echte Core-Aenderung passierte das Gate mit Exit 0 ohne Ausgabe.
func TestVCSStatusImFenceZaehltNicht(t *testing.T) {
	beispiel := "# ADR-0099 — Vorlage\n\n**Status:** Proposed\n\n" +
		"```markdown\n**Status:** Accepted\n```\n\n## Entscheidung\n\n"
	base := []byte(beispiel + "Tue A.\n")
	head := []byte(beispiel + "Tue B.\n")
	fv := &fakeVCS{
		changes: []driven.VCSChange{{Status: driven.VCSModified, Path: adrPath}},
		files:   refs(base, head),
	}
	got, err := CheckVCS(fv, adrConfig(), "BASE", "HEAD")
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Proposed-Datei mit Accepted-Beispiel im Fence ist nicht immutabel → 0 Befunde, got %+v", got)
	}
}
