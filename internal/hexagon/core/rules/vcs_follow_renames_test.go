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
