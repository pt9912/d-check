package rules

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// errIndex simuliert einen unlesbaren git-Index (fail-closed-Pfad).
var errIndex = errors.New("git-Index nicht lesbar (Test)")

// trackedFiles: eine Quell-Datei verlinkt ein getracktes, ein untracktes,
// ein fehlendes und ein Verzeichnis-Ziel (DC-FA-TRK-001-Akzeptanzkriterien).
func trackedFiles() map[string]string {
	return map[string]string{
		"doc.md":    "[t](t.md) [u](u.md) [m](m.md) [d](sub) [out](../raus.md)\n",
		"t.md":      "x",
		"u.md":      "x",
		"sub/in.md": "x",
	}
}

// Happy + Negative: getracktes Ziel still, existierendes untracktes Ziel
// ⇒ target-untracked (Datei, Zeile, Ziel).
func TestCheckTracked_HappyUndNegative(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	lines := PreprocessMarkdown([]byte(trackedFiles()["doc.md"]))
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{"t.md": true, "sub/in.md": true})
	if len(f) != 1 {
		t.Fatalf("Findings = %d, want 1: %+v", len(f), f)
	}
	got := f[0]
	if got.Reason != model.ReasonTargetUntracked || got.Rule != "tracked" ||
		got.File != "doc.md" || got.Line != 1 || got.Target != "u.md" {
		t.Fatalf("unerwarteter Befund: %+v", got)
	}
	if !strings.Contains(got.Message, "frischen Klon") {
		t.Fatalf("Message ohne Klon-Hinweis: %q", got.Message)
	}
}

// Kein Doppelbefund: das fehlende m.md und das Verzeichnis-Ziel sub sind
// keine Kandidaten (target-missing bleibt links; der Index führt nur Dateien);
// ein escapetes Ziel bleibt Sache von links (repo-escape).
func TestCheckTracked_KeinDoppelbefundKeineDirsKeinEscape(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	lines := PreprocessMarkdown([]byte(trackedFiles()["doc.md"]))
	// Leere Index-Menge: NUR die existierenden Datei-Ziele t.md/u.md dürfen
	// feuern — m.md (fehlt), sub (Dir) und ../raus.md (escaped) nie.
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{})
	if len(f) != 2 {
		t.Fatalf("Findings = %d, want 2 (nur t.md+u.md): %+v", len(f), f)
	}
	for _, fd := range f {
		if fd.Target == "m.md" || fd.Target == "sub" || fd.Target == "../raus.md" {
			t.Fatalf("Nicht-Kandidat gemeldet: %+v", fd)
		}
	}
}

// Ventil: ein exempt-targets-Glob über den AUFGELÖSTEN Zielpfad nimmt die
// Referenz aus (referenz-weit).
func TestCheckTracked_ExemptTargets(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	lines := PreprocessMarkdown([]byte(trackedFiles()["doc.md"]))
	cfg := model.TrackedConfig{ExemptTargets: []string{"u.*"}}
	f := CheckTracked(fsys, "doc.md", lines, cfg, map[string]bool{"t.md": true})
	if len(f) != 0 {
		t.Fatalf("Findings = %d, want 0 (u.md exempt): %+v", len(f), f)
	}
}

// fail-closed: aktives tracked ohne verdrahteten VCS-Port ⇒ error (Exit 2
// beim Aufrufer), kein stilles Grün.
func TestRunWithVCS_TrackedFailClosedOhnePort(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{"doc.md": "x"})
	_, err := RunWithVCS(fsys, nil, nil, "", "", model.Config{Roots: []string{"."}}, []string{"tracked"})
	if err == nil || !strings.Contains(err.Error(), "tracked") {
		t.Fatalf("err = %v, want tracked-fail-closed", err)
	}
}

// fail-closed: ein Port-Fehler (unlesbarer Index) propagiert als error.
func TestRunWithVCS_TrackedPortFehler(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{"doc.md": "x"})
	v := &fakeVCS{err: errIndex}
	_, err := RunWithVCS(fsys, nil, v, "", "", model.Config{Roots: []string{"."}}, []string{"tracked"})
	if err == nil || !strings.Contains(err.Error(), errIndex.Error()) {
		t.Fatalf("err = %v, want Port-Fehler", err)
	}
}

// Integration über checkFile: links meldet das fehlende Ziel, tracked NUR das
// existierende untrackte — kein Doppelbefund, beide Module koexistieren.
func TestRunWithVCS_TrackedEndToEnd(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	v := &fakeVCS{tracked: map[string]bool{"t.md": true, "sub/in.md": true}}
	res, err := RunWithVCS(fsys, nil, v, "", "", model.Config{Roots: []string{"."}}, []string{"links", "tracked"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	var missing, untracked int
	for _, f := range res.Findings {
		switch f.Reason {
		case model.ReasonTargetMissing:
			missing++
		case model.ReasonTargetUntracked:
			untracked++
			if f.Target != "u.md" {
				t.Fatalf("target-untracked auf %q, want u.md", f.Target)
			}
		}
	}
	// links: m.md fehlt (+ ../raus.md ist repo-escape, nicht missing);
	// tracked: genau u.md (auch die untracked Quell-Dateien doc.md/u.md selbst
	// sind KEIN Befund — Out-of-Scope der ersten Ausbaustufe).
	if missing != 1 || untracked != 1 {
		t.Fatalf("missing=%d untracked=%d, want 1/1: %+v", missing, untracked, res.Findings)
	}
}

// Modul-aus: ohne tracked ist der Befundsatz frei von target-untracked und
// der (nil-)Port bleibt ungenutzt — default-aus byte-identisch.
func TestRunWithVCS_TrackedAusByteIdentisch(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	res, err := RunWithVCS(fsys, nil, nil, "", "", model.Config{Roots: []string{"."}}, []string{"links"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == model.ReasonTargetUntracked {
			t.Fatalf("target-untracked ohne aktives Modul: %+v", f)
		}
	}
}
