package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// rfGroup ist die Lifecycle-Gruppe der Tests: drei wandernde Orte plus der
// ortsfeste Ruheort.
func rfGroup() []model.ResolveFromGroup {
	return []model.ResolveFromGroup{{
		Dirs:      []string{"plan/open", "plan/next", "plan/in-progress"},
		FixedDirs: []string{"plan/done"},
	}}
}

func rfFindings(t *testing.T, files map[string]string) []model.Finding {
	t.Helper()
	cfg := model.Config{Roots: []string{"."}, ResolveFrom: rfGroup()}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	var out []model.Finding
	for _, f := range res.Findings {
		if f.Reason == model.ReasonLinkPositionDependent {
			out = append(out, f)
		}
	}
	return out
}

// Happy Path: ein Geschwister-Praefix-Verweis loest von jedem Ort der Gruppe
// auf dasselbe Ziel auf.
func TestResolveFromHappy(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](../open/slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	if got := rfFindings(t, files); got != nil {
		t.Fatalf("ueberall aufloesbarer Verweis → 0 Befunde, got %+v", got)
	}
}

// Negative (Aufloesbarkeit): der praefixlose Nachbar-Verweis ist am Ist-Ort
// gruen und bricht beim naechsten Move — genau die 19er-Klasse.
func TestResolveFromPraefixloserNachbar(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	got := rfFindings(t, files)
	if len(got) != 1 || got[0].Reason != model.ReasonLinkPositionDependent {
		t.Fatalf("praefixloser Nachbar → genau ein link-position-dependent, got %+v", got)
	}
}

// Negative (Ziel-Identitaet): der Verweis loest ueberall auf — aber auf
// verschiedene Dateien. Die stille Haelfte der Klasse.
func TestResolveFromDivergierendeZiele(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md":        "# S1\n\nSiehe [N](notiz.md).\n",
		"plan/open/notiz.md":          "# Notiz A\n",
		"plan/next/notiz.md":          "# Notiz B\n",
		"plan/in-progress/notiz.md":   "# Notiz C\n",
		"plan/done/notiz.md":          "# Notiz D\n",
	}
	got := rfFindings(t, files)
	if len(got) != 1 {
		t.Fatalf("divergierende Ziele → genau ein Befund, got %+v", got)
	}
}

// Boundary (ortsfeste Datei): eine done/-Datei ist keine Quelle — ihr nur
// lokal aufloesender Verweis meldet nicht (sonst 108 Falsch-Positive am
// gemessenen Bestand).
func TestResolveFromOrtsfesteDateiIstKeineQuelle(t *testing.T) {
	files := map[string]string{
		"plan/done/slice-9.md":  "# S9\n\nSiehe [N](nachbar.md).\n",
		"plan/done/nachbar.md":  "# N\n",
	}
	if got := rfFindings(t, files); got != nil {
		t.Fatalf("ortsfeste Datei ist keine Quelle → 0 Befunde, got %+v", got)
	}
}

// Der fixed-dirs-Ort zaehlt als hypothetischer Quellort: ein Verweis, der von
// dort nicht aufloeste, meldet — die Datei kann dorthin wandern. Beobachtbar
// nur mit einem fixed-dirs-Ort ANDERER Tiefe (unter Geschwistern loesen
// Geschwister-Praefix-Verweise identisch).
func TestResolveFromFixedDirZaehltAlsOrt(t *testing.T) {
	cfg := model.Config{Roots: []string{"."}, ResolveFrom: []model.ResolveFromGroup{{
		Dirs:      []string{"plan/open", "plan/in-progress"},
		FixedDirs: []string{"plan/archiv/tief"},
	}}}
	files := map[string]string{
		// ../open/slice-2.md loest von beiden dirs auf — von plan/archiv/tief
		// aber als plan/archiv/open/slice-2.md, das nicht existiert.
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](../open/slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	var got []model.Finding
	for _, f := range res.Findings {
		if f.Reason == model.ReasonLinkPositionDependent {
			got = append(got, f)
		}
	}
	if len(got) != 1 {
		t.Fatalf("fixed-dirs-Ort muss als hypothetischer Quellort zaehlen → 1 Befund, got %+v", got)
	}
}

// Und der wurzel-stabile Verweis bleibt ueberall gruen.
func TestResolveFromWurzelStabilerVerweis(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [R](../../plan/open/slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	if got := rfFindings(t, files); got != nil {
		t.Fatalf("wurzel-stabiler Verweis → 0 Befunde, got %+v", got)
	}
}

// Das Ventil gilt: ein ignore-refs-Eintrag nimmt das Ziel auch hier aus.
func TestResolveFromVentil(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	cfg := model.Config{Roots: []string{"."}, ResolveFrom: rfGroup(), IgnoreRefs: []model.IgnoreRef{{Refs: []string{"plan/open/slice-2.md"}}}}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == model.ReasonLinkPositionDependent {
			t.Fatalf("ignore-refs muss auch hier ausnehmen, got %+v", f)
		}
	}
}

// Modul-aus: ohne Gruppen ist der Befundsatz byte-identisch.
func TestResolveFromInertOhneGruppen(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	res, err := Run(coretest.NewMemFS(files), nil, model.Config{Roots: []string{"."}}, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	for _, f := range res.Findings {
		if f.Reason == model.ReasonLinkPositionDependent {
			t.Fatalf("ohne Gruppen inert, got %+v", f)
		}
	}
}

// Realdatenbeleg im Kleinen: die 19er-Klasse — ein Slice mit praefixlosen
// Geschwister-Verweisen meldet VOR dem Move, einmal je Referenz.
func TestResolveFromRetroKlasse(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-a.md": "# A\n\nSiehe [B](slice-b.md) und [C](slice-c.md).\n",
		"plan/open/slice-b.md": "# B\n\nZurueck zu [A](slice-a.md).\n",
		"plan/open/slice-c.md": "# C\n",
	}
	got := rfFindings(t, files)
	if len(got) != 3 {
		t.Fatalf("drei praefixlose Geschwister-Verweise → drei Befunde, got %d (%+v)", len(got), got)
	}
}

// Die Meldung der Aufloesbarkeits-Haelfte nennt den brechenden Ort — ohne den
// Zweig fiele der Fall in die Divergenz-Meldung (falsche Diagnose).
func TestResolveFromMeldungNenntDenOrt(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [S2](slice-2.md).\n",
		"plan/open/slice-2.md": "# S2\n",
	}
	got := rfFindings(t, files)
	if len(got) != 1 || !strings.Contains(got[0].Message, "löst von plan/") {
		t.Fatalf("die Meldung muss den nicht aufloesenden Ort nennen, got %+v", got)
	}
}

// Ein Ziel, das schon am IST-Ort fehlt, meldet target-missing — resolve-from
// schweigt (kein Doppelbefund derselben Referenz).
func TestResolveFromKeinDoppelbefundBeiTargetMissing(t *testing.T) {
	files := map[string]string{
		"plan/open/slice-1.md": "# S1\n\nSiehe [T](tippfehler.md).\n",
	}
	cfg := model.Config{Roots: []string{"."}, ResolveFrom: rfGroup()}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"links"})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	var missing, posdep int
	for _, f := range res.Findings {
		switch f.Reason {
		case model.ReasonTargetMissing:
			missing++
		case model.ReasonLinkPositionDependent:
			posdep++
		}
	}
	if missing != 1 || posdep != 0 {
		t.Fatalf("fehlendes Ziel → genau ein target-missing, kein position-dependent (got missing=%d posdep=%d)", missing, posdep)
	}
}
