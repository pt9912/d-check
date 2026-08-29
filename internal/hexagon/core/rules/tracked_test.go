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
// ein fehlendes, ein Verzeichnis- und ein Bild-Ziel
// (DC-FA-TRK-001-Akzeptanzkriterien; Bild-Referenzen zählen — Schritt 3).
func trackedFiles() map[string]string {
	return map[string]string{
		"doc.md":    "[t](t.md) [u](u.md) [m](m.md) [d](sub) [out](../raus.md) ![i](img.png)\n",
		"t.md":      "x",
		"u.md":      "x",
		"img.png":   "px",
		"sub/in.md": "x",
	}
}

// Happy + Negative: getracktes Ziel still, existierendes untracktes Ziel
// ⇒ target-untracked (Datei, Zeile, Ziel).
func TestCheckTracked_HappyUndNegative(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	lines := PreprocessMarkdown([]byte(trackedFiles()["doc.md"]))
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{"t.md": true, "sub/in.md": true})
	// Zwei Befunde: das untrackte Link-Ziel u.md UND das untrackte
	// Bild-Ziel img.png (Bild-Referenzen zählen — verriegelt Schritt 3
	// gegen eine IsImage-Skip-Mutation, R2-M4).
	if len(f) != 2 {
		t.Fatalf("Findings = %d, want 2: %+v", len(f), f)
	}
	targets := map[string]bool{}
	for _, fd := range f {
		if fd.Reason != model.ReasonTargetUntracked || fd.Rule != "tracked" || fd.File != "doc.md" {
			t.Fatalf("unerwarteter Befund: %+v", fd)
		}
		targets[fd.Target] = true
	}
	if !targets["u.md"] || !targets["img.png"] {
		t.Fatalf("Targets = %v, want u.md + img.png", targets)
	}
	if !strings.Contains(f[0].Message, "frischen Klon") {
		t.Fatalf("Message ohne Klon-Hinweis: %q", f[0].Message)
	}
}

// Befund-target trägt den AUFGELÖSTEN Zielpfad (Spec .a Schritt 5) — die
// Form, die das exempt-targets-Ventil matcht; roh geschriebene "./"-Ziele
// werden für Lookup UND Befund normalisiert (R2-M1-Verriegelung).
func TestCheckTracked_TargetAufgeloest(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{
		"doc.md": "[roh](./v.md) [ok](./t.md)\n",
		"v.md":   "x",
		"t.md":   "x",
	})
	lines := PreprocessMarkdown([]byte("[roh](./v.md) [ok](./t.md)\n"))
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{"t.md": true})
	if len(f) != 1 {
		t.Fatalf("Findings = %d, want 1 (nur ./v.md; ./t.md ist über den aufgelösten Pfad getrackt): %+v", len(f), f)
	}
	if f[0].Target != "v.md" {
		t.Fatalf("Target = %q, want aufgelöst \"v.md\" (nicht die rohe Schreibweise)", f[0].Target)
	}
}

// Symlink-Referenzen sind die Domäne von DC-FA-LINK-002 (links meldet
// `symlink` kategorisch): tracked prüft weder ein Symlink-Ziel noch einen
// Pfad DURCH einen Verzeichnis-Symlink — sonst false-positive, denn der
// Index führt den realen Pfad, nicht den Alias (R2-M2/M3-Verriegelung).
func TestCheckTracked_SymlinkZieleUeberspringen(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{
		"doc.md":       "[s](slink.md) [durch](dlink/in.md)\n",
		"slink.md":     "x",
		"dlink/in.md":  "x",
	})
	fsys.AddSymlink("slink.md")
	fsys.AddSymlink("dlink")
	lines := PreprocessMarkdown([]byte("[s](slink.md) [durch](dlink/in.md)\n"))
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{})
	if len(f) != 0 {
		t.Fatalf("Findings = %d, want 0 (Symlink-Pfade sind links-Domäne): %+v", len(f), f)
	}
}

// Kein Doppelbefund: das fehlende m.md und das Verzeichnis-Ziel sub sind
// keine Kandidaten (target-missing bleibt links; der Index führt nur Dateien);
// ein escapetes Ziel bleibt Sache von links (repo-escape).
func TestCheckTracked_KeinDoppelbefundKeineDirsKeinEscape(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	lines := PreprocessMarkdown([]byte(trackedFiles()["doc.md"]))
	// Leere Index-Menge: NUR die existierenden Datei-Ziele t.md/u.md/img.png
	// dürfen feuern — m.md (fehlt), sub (Dir) und ../raus.md (escaped) nie.
	f := CheckTracked(fsys, "doc.md", lines, model.TrackedConfig{}, map[string]bool{})
	if len(f) != 3 {
		t.Fatalf("Findings = %d, want 3 (nur t.md+u.md+img.png): %+v", len(f), f)
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
	cfg := model.TrackedConfig{ExemptTargets: []string{"u.*", "img.png"}}
	f := CheckTracked(fsys, "doc.md", lines, cfg, map[string]bool{"t.md": true})
	if len(f) != 0 {
		t.Fatalf("Findings = %d, want 0 (u.md + img.png exempt): %+v", len(f), f)
	}
}

// fail-closed: aktives tracked ohne verdrahteten VCS-Port ⇒ error (Exit 2
// beim Aufrufer), kein stilles Grün.
func TestRunWithVCS_TrackedFailClosedOhnePort(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{"doc.md": "x"})
	_, err := RunWithVCS(fsys, nil, nil, nil, "", "", model.Config{Roots: []string{"."}}, []string{"tracked"})
	if err == nil || !strings.Contains(err.Error(), "tracked") {
		t.Fatalf("err = %v, want tracked-fail-closed", err)
	}
}

// fail-closed: ein Port-Fehler (unlesbarer Index) propagiert als error.
func TestRunWithVCS_TrackedPortFehler(t *testing.T) {
	fsys := coretest.NewMemFS(map[string]string{"doc.md": "x"})
	v := &fakeVCS{err: errIndex}
	_, err := RunWithVCS(fsys, nil, v, nil, "", "", model.Config{Roots: []string{"."}}, []string{"tracked"})
	if err == nil || !strings.Contains(err.Error(), errIndex.Error()) {
		t.Fatalf("err = %v, want Port-Fehler", err)
	}
}

// Integration über checkFile: links meldet das fehlende Ziel, tracked NUR das
// existierende untrackte — kein Doppelbefund, beide Module koexistieren.
func TestRunWithVCS_TrackedEndToEnd(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	v := &fakeVCS{tracked: map[string]bool{"t.md": true, "sub/in.md": true}}
	res, err := RunWithVCS(fsys, nil, v, nil, "", "", model.Config{Roots: []string{"."}}, []string{"links", "tracked"})
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
			if f.Target != "u.md" && f.Target != "img.png" {
				t.Fatalf("target-untracked auf %q, want u.md/img.png", f.Target)
			}
		}
	}
	// links: m.md fehlt (+ ../raus.md ist repo-escape, nicht missing);
	// tracked: genau u.md + img.png (auch die untracked Quell-Dateien
	// doc.md/u.md selbst sind KEIN Befund — Out-of-Scope der ersten Stufe).
	if missing != 1 || untracked != 2 {
		t.Fatalf("missing=%d untracked=%d, want 1/2: %+v", missing, untracked, res.Findings)
	}
}

// Modul-aus: ohne tracked ist der Befundsatz frei von target-untracked und
// der (nil-)Port bleibt ungenutzt — default-aus byte-identisch.
func TestRunWithVCS_TrackedAusByteIdentisch(t *testing.T) {
	fsys := coretest.NewMemFS(trackedFiles())
	res, err := RunWithVCS(fsys, nil, nil, nil, "", "", model.Config{Roots: []string{"."}}, []string{"links"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == model.ReasonTargetUntracked {
			t.Fatalf("target-untracked ohne aktives Modul: %+v", f)
		}
	}
}
