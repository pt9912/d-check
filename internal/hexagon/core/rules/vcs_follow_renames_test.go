package rules

import (
	"errors"
	"testing"
)

// Mit follow-renames gleicht der Core ein Ziel über die Umbenennung der Range
// gleich (DC-FA-VCS-001.a Schritt 4, ADR-0115); die Zuordnung Quelle → Ende
// liefert der Port, hier ein Doppelgänger. Je Fall steht die Zahl der Befunde
// in den vier Schlüssel-Kombinationen.
func TestVCSFollowRenames(t *testing.T) {
	type counts struct{ follow, both, ignore, off int }
	rename := map[string]string{"spec/s.md": "spec/t.md"}
	cases := []struct {
		name       string
		base, head string
		renames    map[string]string
		want       counts
	}{
		{"nachzug auf die umbenannte datei",
			"[S](../../../spec/s.md)", "[S](../../../spec/t.md)", rename, counts{0, 0, 1, 1}},
		{"nachzug mit anker auf die umbenannte datei",
			"[S](../../../spec/s.md#a)", "[S](../../../spec/t.md#a)", rename, counts{0, 0, 1, 1}},
		{"bild-nachzug und referenz-definition",
			"![S](../../../spec/s.md)\n\n[r]: ../../../spec/s.md", "![S](../../../spec/t.md)\n\n[r]: ../../../spec/t.md",
			rename, counts{0, 0, 1, 1}},
		{"gelöscht und anderes angelegt: keine umbenennung",
			"[S](../../../spec/s.md)", "[S](../../../spec/t.md)", nil, counts{1, 1, 1, 1}},
		{"umbiegen auf eine bestehende, nicht umbenannte datei",
			"[S](../../../spec/s.md)", "[S](../../../spec/c.md)", rename, counts{1, 1, 1, 1}},
		{"anker geändert",
			"[S](../../../spec/s.md#alt)", "[S](../../../spec/t.md#neu)", rename, counts{1, 1, 1, 1}},
		{"linktext geändert",
			"[S](../../../spec/s.md)", "[T](../../../spec/t.md)", rename, counts{1, 1, 1, 1}},
		{"übrige zeile geändert",
			"Siehe [S](../../../spec/s.md).", "Lies [S](../../../spec/t.md).", rename, counts{1, 1, 1, 1}},
		{"link entfernt",
			"Siehe [S](../../../spec/s.md).", "Siehe S.", rename, counts{1, 1, 1, 1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				label          string
				follow, ignore bool
				want           int
			}{
				{"follow-renames allein", true, false, c.want.follow},
				{"beide schlüssel", true, true, c.want.both},
				{"ignore-link-targets allein", false, true, c.want.ignore},
				{"ohne schlüssel", false, false, c.want.off},
			} {
				cfg := adrConfig()
				cfg.FollowRenames, cfg.IgnoreLinkTargets = mode.follow, mode.ignore
				fv := vcsRenameFake(adr("Accepted", c.base), adr("Accepted", c.head), c.renames)
				got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
				if err != nil {
					t.Fatalf("%s: unerwarteter Fehler: %v", mode.label, err)
				}
				if len(got) != mode.want {
					t.Fatalf("%s: Befunde = %d, want %d (%v)", mode.label, len(got), mode.want, got)
				}
				if mode.follow != (fv.renamesCall == 1) {
					t.Fatalf("%s: Renames-Aufrufe = %d, mit Schlüssel genau einer, ohne keiner", mode.label, fv.renamesCall)
				}
			}
		})
	}
}

// vcsRenameFake legt die Pfad-Bäume an, auf die die Links der Fälle zeigen:
// spec/s.md besteht nur in BASE, spec/t.md nur in HEAD, spec/c.md in beiden.
func vcsRenameFake(base, head []byte, renames map[string]string) *fakeVCS {
	files := refs(base, head)
	files["BASE"]["spec/s.md"], files["BASE"]["spec/c.md"] = nil, nil
	files["HEAD"]["spec/t.md"], files["HEAD"]["spec/c.md"] = nil, nil
	return &fakeVCS{files: files, renames: renames}
}

// Ein unveränderter Link auf eine umbenannte Datei bleibt gleich: die
// Gleichsetzung ist eine Funktion des Pfads und gilt auf beiden Seiten
// (ADR-0115 Entscheidung 9).
func TestVCSFollowRenamesUnveraenderterLink(t *testing.T) {
	cfg := adrConfig()
	cfg.FollowRenames = true
	body := adr("Accepted", "Siehe [S](../../../spec/s.md).")
	head := append(append([]byte{}, body...), []byte("| 2026-02-02 | Notiz |\n")...)
	fv := vcsRenameFake(body, head, map[string]string{"spec/s.md": "spec/t.md"})
	got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
	if err != nil || len(got) != 0 {
		t.Fatalf("unveränderter Link auf umbenannte Datei: kein Befund erwartet, got %v err=%v", got, err)
	}
}

// Die Port-Methode wird nur mit dem Schlüssel gerufen und nur, wenn das Modul
// überhaupt etwas schützt; ihr Fehler bricht den Lauf, und eine Zuordnung, in
// der eine Quelle zugleich Ende ist, ebenso.
func TestVCSFollowRenamesPortMethode(t *testing.T) {
	cfg := adrConfig()
	cfg.FollowRenames = true

	inert := vcsRenameFake(adr("Accepted", "x"), adr("Accepted", "x"), nil)
	leer := cfg
	leer.Paths = nil
	if _, err := CheckVCS(inert, leer, "BASE", "HEAD"); err != nil || inert.renamesCall != 0 {
		t.Fatalf("ohne vcs.paths inert und ohne Renames-Aufruf erwartet: err=%v Aufrufe=%d", err, inert.renamesCall)
	}

	fehler := vcsRenameFake(adr("Accepted", "x"), adr("Accepted", "x"), nil)
	fehler.renamesErr = errors.New("Baum nicht lesbar")
	if _, err := CheckVCS(fehler, cfg, "BASE", "HEAD"); err == nil {
		t.Fatal("Fehler der Port-Methode still übergangen — erwartet war ein Fehler")
	}

	kette := vcsRenameFake(adr("Accepted", "x"), adr("Accepted", "x"), map[string]string{"a.md": "b.md", "b.md": "c.md"})
	if _, err := CheckVCS(kette, cfg, "BASE", "HEAD"); err == nil {
		t.Fatal("Quelle zugleich Ende still übergangen — erwartet war ein Fehler")
	}
}

// Randfälle der Gleichsetzung (DC-FA-VCS-001): die Zuordnung des Ports enthält
// nur die teilnehmenden Umbenennungen; ein Ziel, das nicht teilnimmt, bleibt
// roh (follow-renames allein) oder wird auf Dateiname und Anker normiert
// (mit ignore-link-targets), und die beiden Formen gleichen sich nie.
func TestVCSFollowRenamesRandfaelle(t *testing.T) {
	type counts struct{ follow, both, ignore int }
	cases := []struct {
		name                 string
		baseLink, headLink   string
		baseTree, headTree   []string
		renames              map[string]string
		want                 counts
	}{
		{"beide Schlüssel: Umbiegen auf eine andere Datei mit dem Namen der Quelle",
			"spec/s.md", "other/s.md",
			[]string{"spec/s.md", "other/s.md"}, []string{"spec/t.md", "other/s.md"},
			map[string]string{"spec/s.md": "spec/t.md"}, counts{1, 1, 0}},
		{"Umbiegen auf den alten Namen: d2/c.md wird zu d3/a.md",
			"d1/a.md", "d2/c.md",
			[]string{"d1/a.md", "d2/c.md"}, []string{"d1/a.md", "d3/a.md"},
			map[string]string{"d2/c.md": "d3/a.md"}, counts{1, 1, 1}},
		{"Umbenennung ohne Namenswechsel: Umbiegen auf eine andere Datei gleichen Namens",
			"d1/x.md", "d3/x.md",
			[]string{"d1/x.md"}, []string{"d2/x.md", "d3/x.md"},
			map[string]string{"d1/x.md": "d2/x.md"}, counts{1, 1, 0}},
		{"zwei nicht teilnehmende Ziele mit gleichem Namen",
			"d1/x.md", "d3/x.md",
			[]string{"d1/x.md"}, []string{"d1/x.md", "d3/x.md"},
			nil, counts{1, 0, 0}},
		{"freigemachter Name: Nachzug von C auf A mit verschiedenen Namen",
			"d2/c.md", "d1/a.md",
			[]string{"d1/a.md", "d2/c.md"}, []string{"d1/a.md"},
			nil, counts{1, 1, 1}},
		{"freigemachter Name: Umbiegen von C auf eine andere Datei mit dem Namen von C",
			"d2/c.md", "d9/c.md",
			[]string{"d1/a.md", "d2/c.md", "d9/c.md"}, []string{"d1/a.md", "d9/c.md"},
			nil, counts{1, 0, 0}},
		{"zwei Quellen auf ein Ende: Umbiegen von A auf B",
			"d1/a.md", "d2/b.md",
			[]string{"d1/a.md", "d2/b.md"}, []string{"d3/e.md"},
			nil, counts{1, 1, 1}},
		{"zwei Quellen auf ein Ende: Umbiegen von A auf eine andere Datei mit dem Namen von A",
			"d1/a.md", "d9/a.md",
			[]string{"d1/a.md", "d2/b.md", "d9/a.md"}, []string{"d3/e.md", "d9/a.md"},
			nil, counts{1, 0, 0}},
		{"Rotation: Nachzug von A auf X folgt",
			"a.md", "x.md",
			[]string{"a.md", "c.md"}, []string{"x.md", "a.md"},
			map[string]string{"a.md": "x.md"}, counts{0, 0, 1}},
		{"Rotation: Nachzug von C auf A bleibt Drift",
			"c.md", "a.md",
			[]string{"a.md", "c.md"}, []string{"x.md", "a.md"},
			map[string]string{"a.md": "x.md"}, counts{1, 1, 1}},
		{"Verzeichnis umbenannt: ein Link auf das Verzeichnis folgt nicht",
			"d1/", "d2/",
			[]string{"d1/a.md"}, []string{"d2/a.md"},
			map[string]string{"d1/a.md": "d2/a.md"}, counts{1, 1, 1}},
		{"Verzeichnis umbenannt: ein Link auf eine Datei darin folgt",
			"d1/a.md", "d2/a.md",
			[]string{"d1/a.md"}, []string{"d2/a.md"},
			map[string]string{"d1/a.md": "d2/a.md"}, counts{0, 0, 0}},
		{"Zwischenpfad einer Kette neu belegt: der Nachzug auf ihn ist Drift",
			"a.md", "b.md",
			[]string{"a.md"}, []string{"d.md", "b.md"},
			map[string]string{"a.md": "d.md"}, counts{1, 1, 1}},
		{"Nachzug auf das Ende über die Kette",
			"a.md", "d.md",
			[]string{"a.md"}, []string{"d.md", "b.md"},
			map[string]string{"a.md": "d.md"}, counts{0, 0, 1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				label          string
				follow, ignore bool
				want           int
			}{
				{"follow-renames allein", true, false, c.want.follow},
				{"beide schlüssel", true, true, c.want.both},
				{"ignore-link-targets allein", false, true, c.want.ignore},
				{"ohne schlüssel", false, false, 1},
			} {
				cfg := adrConfig()
				cfg.FollowRenames, cfg.IgnoreLinkTargets = mode.follow, mode.ignore
				base := adr("Accepted", "Siehe [S](../../../"+c.baseLink+").")
				head := adr("Accepted", "Siehe [S](../../../"+c.headLink+").")
				fv := vcsTreeFake(base, head, c.baseTree, c.headTree, c.renames)
				got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
				if err != nil {
					t.Fatalf("%s: unerwarteter Fehler: %v", mode.label, err)
				}
				if len(got) != mode.want {
					t.Fatalf("%s: Befunde = %d, want %d (%v)", mode.label, len(got), mode.want, got)
				}
			}
		})
	}
}

// vcsTreeFake legt die Pfad-Bäume beider Stände fest, auf die die Links zeigen.
func vcsTreeFake(base, head []byte, baseTree, headTree []string, renames map[string]string) *fakeVCS {
	files := refs(base, head)
	for _, p := range baseTree {
		files["BASE"][p] = nil
	}
	for _, p := range headTree {
		files["HEAD"][p] = nil
	}
	return &fakeVCS{files: files, renames: renames}
}

// Die Umbenennung der immutablen Datei selbst bleibt Drift, mit und ohne den
// Schlüssel: ihr Pfad ist stabil, der Schlüssel wirkt nur auf Link-Ziele.
func TestVCSFollowRenamesImmutableDateiSelbst(t *testing.T) {
	neu := "docs/plan/adr/0100-x.md"
	for _, follow := range []bool{false, true} {
		cfg := adrConfig()
		cfg.FollowRenames = follow
		fv := &fakeVCS{
			files: map[string]map[string][]byte{
				"BASE": {adrPath: adr("Accepted", "Tue A.")},
				"HEAD": {neu: adr("Accepted", "Tue A.")},
			},
			renames: map[string]string{adrPath: neu},
		}
		got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
		if err != nil || len(got) != 1 || got[0].Target != adrPath {
			t.Fatalf("follow-renames=%v: ein Befund am alten Pfad erwartet, got %v err=%v", follow, got, err)
		}
	}
}

// Beide Schlüssel zusammen lassen die Tests von ignore-link-targets allein
// unberührt: derselbe Nachzug auf eine umgezogene Datei gleichen Namens bleibt
// kein Befund, und eine Änderung am Linktext bleibt einer.
func TestVCSFollowRenamesBeideSchluesselIgnoreBleibt(t *testing.T) {
	cases := []struct {
		name       string
		base, head string
		want       int
	}{
		{"umzug gleichen namens", "[R](../../../docs/user/releasing.md)", "[R](../../../docs/user/maintainer/releasing.md)", 0},
		{"linktext geändert", "[R](../../../docs/user/releasing.md)", "[Neu](../../../docs/user/maintainer/releasing.md)", 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := adrConfig()
			cfg.FollowRenames, cfg.IgnoreLinkTargets = true, true
			fv := vcsTreeFake(adr("Accepted", c.base), adr("Accepted", c.head),
				[]string{"docs/user/releasing.md"}, []string{"docs/user/maintainer/releasing.md"}, nil)
			got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
			if err != nil || len(got) != c.want {
				t.Fatalf("Befunde = %d, want %d (%v) err=%v", len(got), c.want, got, err)
			}
		})
	}
}
