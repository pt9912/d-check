package rules

import "testing"

// Mit ignore-link-targets ist ein reiner Pfad-Nachzug keine Core-Drift; jede
// Änderung am Linktext, an der übrigen Zeile oder an der Zahl der Links bleibt
// eine. Ohne den Schlüssel meldet derselbe Nachzug — der Vergleich davor.
func TestVCSIgnoreLinkTargets(t *testing.T) {
	cases := []struct {
		name       string
		base, head string
		mitSchalter, ohne int
	}{
		{"inline-ziel nachgezogen",
			"Siehe [Releasing](../../user/releasing.md).",
			"Siehe [Releasing](../../user/maintainer/releasing.md).", 0, 1},
		{"ziel mit anker nachgezogen",
			"Siehe [Prep](../../user/releasing.md#release-prep-vor-dem-tag).",
			"Siehe [Prep](../../user/maintainer/releasing.md#release-prep-vor-dem-tag).", 0, 1},
		{"zwei links in einer zeile nachgezogen",
			"[A](a.md) und [B](b.md)", "[A](x/a.md) und [B](x/b.md)", 0, 1},
		{"ziel mit titel nachgezogen, titel gleich",
			`[A](a.md "Titel")`, `[A](x/a.md "Titel")`, 0, 1},
		{"bild nachgezogen", "![Bild](img/a.png)", "![Bild](x/img/a.png)", 0, 1},
		{"referenz-definition nachgezogen", "[rel]: ../../user/releasing.md",
			"[rel]: ../../user/maintainer/releasing.md", 0, 1},
		{"linktext geaendert", "[Releasing](a.md)", "[Release-Prozess](a.md)", 1, 1},
		{"text neben dem link geaendert", "Siehe [R](a.md).", "Lies [R](x/a.md).", 1, 1},
		{"titel geaendert", `[A](a.md "Alt")`, `[A](x/a.md "Neu")`, 1, 1},
		{"link entfernt", "Siehe [R](a.md).", "Siehe R.", 1, 1},
		{"link hinzugefuegt", "Siehe R.", "Siehe [R](a.md).", 1, 1},
		{"referenz-label geaendert", "[rel]: a.md", "[rel2]: x/a.md", 1, 1},
		// die benannte Grenze: ohne Code-Kontext wird auch ein Ziel in Inline-Code
		// geleert
		{"ziel in inline-code", "`[R](a.md)`", "`[R](x/a.md)`", 0, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				schalter bool
				want     int
			}{{true, c.mitSchalter}, {false, c.ohne}} {
				cfg := adrConfig()
				cfg.IgnoreLinkTargets = mode.schalter
				fv := &fakeVCS{files: refs(adr("Accepted", c.base), adr("Accepted", c.head))}
				got, err := CheckVCS(fv, cfg, "BASE", "HEAD")
				if err != nil {
					t.Fatalf("unerwarteter Fehler: %v", err)
				}
				if len(got) != mode.want {
					t.Fatalf("ignore-link-targets=%v: Befunde = %d, want %d (%v)", mode.schalter, len(got), mode.want, got)
				}
			}
		})
	}
}
