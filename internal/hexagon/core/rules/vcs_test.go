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
// Akzeptanztests von CheckVCS (DC-FA-VCS-001): AllPaths leitet seine Mengen aus
// files her (die Schlüssel je Ref sind die "im Tree vorhandenen" Pfade). err
// simuliert den fail-closed-Pfad (fehlendes .git/Range); allPathsErr simuliert
// gezielt einen unlesbaren Unterbaum an einem der beiden Enden (slice-220).
type fakeVCS struct {
	files       map[string]map[string][]byte // ref → pfad → inhalt (auch AllPaths-Quelle)
	allPathsErr map[string]error             // ref → Fehler von AllPaths an genau diesem Ende
	commits     []driven.CommitMeta          // Modul commits (DC-FA-COMMITS-001)
	tracked     map[string]bool              // Modul tracked (DC-FA-TRK-001)
	err         error
	fileErr     error // FileAt scheitert — unlesbares Objekt oder unlesbarer Tree-Eintrag
}

func (f *fakeVCS) AllPaths(base, head string) ([]string, []string, error) {
	if f.err != nil {
		return nil, nil, f.err
	}
	baseAll, err := f.pathsAt(base)
	if err != nil {
		return nil, nil, err
	}
	headAll, err := f.pathsAt(head)
	if err != nil {
		return nil, nil, err
	}
	return baseAll, headAll, nil
}

func (f *fakeVCS) pathsAt(ref string) ([]string, error) {
	if err, ok := f.allPathsErr[ref]; ok {
		return nil, err
	}
	m := f.files[ref]
	out := make([]string, 0, len(m))
	for p := range m {
		out = append(out, p)
	}
	return out, nil
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
	if f.fileErr != nil {
		return nil, false, f.fileErr
	}
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
			fv := &fakeVCS{files: refs(c.base, c.head)}
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
	drift := &fakeVCS{files: refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B."))}
	if got, err := CheckVCS(drift, cfg, "BASE", "HEAD"); err != nil || len(got) != 1 || got[0].Reason != model.ReasonCoreDriftVCS {
		t.Fatalf("ohne status-line: ein core-drift-vcs (Body) erwartet, got %v err=%v", got, err)
	}
	clean := &fakeVCS{files: refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue  A."))}
	if got, err := CheckVCS(clean, cfg, "BASE", "HEAD"); err != nil || len(got) != 0 {
		t.Fatalf("ohne status-line + reflow: kein Befund erwartet, got %v err=%v", got, err)
	}
}

// TestVCSDeleteAddClass deckt Löschung, Hinzufügung und die Klassen-Filterung ab.
func TestVCSDeleteAddClass(t *testing.T) {
	cfg := adrConfig()
	cases := []struct {
		name  string
		files map[string]map[string][]byte
		want  int
	}{
		{"geloeschte accepted feuert", refs(adr("Accepted", "Tue A."), nil), 1},
		{"geloeschte proposed ist frei", refs(adr("Proposed", "Tue A."), nil), 0},
		{"neue datei ist frei", refs(nil, adr("Accepted", "Tue A.")), 0},
		{"unveraenderte datei bleibt unveraendert",
			refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue A.")), 0},
		{"pfad ausserhalb der klasse ignoriert",
			map[string]map[string][]byte{"BASE": {"docs/other.md": adr("Accepted", "Tue A.")}, "HEAD": {"docs/other.md": adr("Accepted", "Tue B.")}}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fv := &fakeVCS{files: c.files}
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
	fv := &fakeVCS{files: refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B."))}
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
	driftPort := &fakeVCS{files: refs(adr("Accepted", "Tue A."), adr("Accepted", "Tue B."))}

	res, err := RunWithVCS(m, nil, driftPort, nil, "BASE", "HEAD", cfg, []string{"vcs"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != model.ReasonCoreDriftVCS {
		t.Fatalf("aktiv: ein core-drift-vcs erwartet, got %v", res.Findings)
	}

	// Modul-aus: derselbe Port, aber vcs nicht im Modulsatz ⇒ keine vcs-Befunde.
	resOff, err := RunWithVCS(m, nil, driftPort, nil, "BASE", "HEAD", cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	for _, f := range resOff.Findings {
		if f.Reason == model.ReasonCoreDriftVCS {
			t.Fatal("Modul-aus: vcs hätte nicht laufen dürfen")
		}
	}

	// fail-closed über RunWithVCS.
	if _, err := RunWithVCS(m, nil, &fakeVCS{err: errors.New("kein .git")}, nil, "BASE", "HEAD", cfg, []string{"vcs"}); err == nil {
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
	fv := &fakeVCS{files: refs(base, head)}
	got, err := CheckVCS(fv, adrConfig(), "BASE", "HEAD")
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Proposed-Datei mit Accepted-Beispiel im Fence ist nicht immutabel → 0 Befunde, got %+v", got)
	}
}

// Ein Rename AUS der immutablen Klasse heraus: die Prüfung hängt am ALTEN Pfad.
// Der Delete des in-Klasse-Pfads trägt den Befund, der Add auf einem Pfad
// außerhalb von vcs.paths ist frei. Grenze, die dieser Test hält: wer die
// Klassen-Prüfung je auf den NEUEN Pfad umstellte, machte genau diesen Fall
// still — und der Rename innerhalb der Klasse bliebe grün, also unauffällig.
func TestVCSRenameOutOfClass(t *testing.T) {
	fv := &fakeVCS{files: map[string]map[string][]byte{
		"BASE": {adrPath: adr("Accepted", "Tue A.")},
		"HEAD": {"docs/notes/kern.md": adr("Accepted", "Tue A.")},
	}}
	got, err := CheckVCS(fv, adrConfig(), "BASE", "HEAD")
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Rename aus der Klasse: ein core-drift-vcs erwartet, got %d (%v)", len(got), got)
	}
	if got[0].Target != adrPath {
		t.Fatalf("Befund-Ziel = %q, erwartet der ALTE Pfad %q", got[0].Target, adrPath)
	}
}

// TestVCSAllPathsFehlerAnBeidenEnden: ein unlesbarer Unterbaum an BASE ODER
// an HEAD bricht CheckVCS fail-closed ab — unabhängig davon, ob es zu einem
// betroffenen Pfad ein "Pendant" auf der Gegenseite gibt (CO-001s zweite und
// dritte Ausprägung fallen im neuen Entwurf auf denselben Codepfad: die
// AllPaths-Fehlerprüfung, nicht mehr ein Sonderfall im "Added"-Zweig).
func TestVCSAllPathsFehlerAnBeidenEnden(t *testing.T) {
	cfg := adrConfig()
	unlesbarBase := &fakeVCS{allPathsErr: map[string]error{"BASE": errors.New("Unterbaum nicht lesbar")}}
	if _, err := CheckVCS(unlesbarBase, cfg, "BASE", "HEAD"); err == nil {
		t.Fatal("unlesbarer BASE-Tree still passiert — erwartet war ein Fehler")
	}

	// Die vierte, bisher fehldiagnostizierte Ausprägung (CO-001): ein
	// unlesbarer HEAD-Tree darf nicht als core-drift-vcs "gelöscht" erscheinen,
	// sondern muss den Umgebungsfehler selbst melden.
	unlesbarHead := &fakeVCS{
		files:       refs(adr("Accepted", "Tue A."), nil),
		allPathsErr: map[string]error{"HEAD": errors.New("Unterbaum nicht lesbar")},
	}
	if _, err := CheckVCS(unlesbarHead, cfg, "BASE", "HEAD"); err == nil {
		t.Fatal("unlesbarer HEAD-Tree still als Loeschung gemeldet — erwartet war ein Fehler")
	}

	// Gegenrichtung: eine wirklich neu angelegte Datei bleibt befundfrei.
	echtNeu := &fakeVCS{
		files: map[string]map[string][]byte{"HEAD": {"docs/plan/adr/0001-a.md": adr("Accepted", "Neu.")}},
	}
	if got, err := CheckVCS(echtNeu, cfg, "BASE", "HEAD"); err != nil || len(got) != 0 {
		t.Errorf("neu angelegte Datei falsch behandelt: %d Befund(e), err=%v", len(got), err)
	}
}
