package rules

import "testing"

// Mit ignore-link-targets ist ein reiner Pfad-Nachzug keine Core-Drift; jede
// Änderung am Linktext, an der übrigen Zeile oder an der Zahl der Links bleibt
// eine. Ohne den Schlüssel bleibt jeder Nachzug Drift.
func TestVCSIgnoreLinkTargets(t *testing.T) {
	cases := []struct {
		name              string
		base, head        string
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
		{"spitzklammer-ziel mit leerraum nachgezogen",
			"[A](<alt/a b.md>)", "[A](<neu/a b.md>)", 0, 1},
		{"referenz-definition mit titel nachgezogen", `[rel]: a.md "T"`, `[rel]: x/a.md "T"`, 0, 1},
		{"referenz-titel geaendert", `[rel]: a.md "Alt"`, `[rel]: x/a.md "Neu"`, 1, 1},
		// keine Link-Syntax: das Wort bleibt Teil des Core
		{"fussnote geaendert", "[^1]: Nicht erlaubt.", "[^1]: Immer erlaubt.", 1, 1},
		{"fussnote mit pfadwort geaendert", "[^1]: a.md", "[^1]: x/a.md", 1, 1},
		{"prosa in referenz-form geaendert", "[Hinweis]: Verboten ist das.", "[Hinweis]: Erlaubt ist das.", 1, 1},
		{"einzelwort ohne pfadzeichen geaendert", "[Hinweis]: Verboten", "[Hinweis]: Erlaubt", 1, 1},
		{"klammer ohne link geaendert", "f](a.md und weiter", "f](x/a.md und weiter", 1, 1},
		{"offener link ohne schliessende klammer", "[f](a.md und weiter", "[f](x/a.md und weiter", 1, 1},
		// balancierte Klammern liest die links-Erkennung als Teil von Ziel und Linktext
		{"ziel mit eigener klammer nachgezogen", "[A](a(1).md)", "[A](x/a(1).md)", 0, 1},
		{"linktext mit eckiger klammer nachgezogen", "[a[0]](a.md)", "[a[0]](x/a.md)", 0, 1},
		// Code und Escapes: link-förmiger Text, der kein Link ist, bleibt
		{"ziel in inline-code", "`[R](a.md)`", "`[R](x/a.md)`", 1, 1},
		{"funktionsaufruf in inline-code geaendert", "`fns[k](true)`", "`fns[k](false)`", 1, 1},
		{"escapte klammern geaendert", `\[Ausnahme\](a.md)`, `\[Ausnahme\](x/a.md)`, 1, 1},
		{"escapte oeffnende klammer geaendert", `\[Ausnahme](a.md)`, `\[Ausnahme](x/a.md)`, 1, 1},
		{"link im codeblock nachgezogen", "```\n[R](a.md)\n```", "```\n[R](x/a.md)\n```", 1, 1},
		{"link neben inline-code nachgezogen", "`x` [R](a.md)", "`x` [R](x/a.md)", 0, 1},
		{"doppelter backslash vor link nachgezogen", `\\[R](a.md)`, `\\[R](x/a.md)`, 0, 1},
		{"code-span ab dem linktext geaendert", "[a `b](a.md)` c", "[a `b](x/a.md)` c", 1, 1},
		{"eingerueckter code geaendert", "Text\n\n    [a](a.md)", "Text\n\n    [a](x/a.md)", 1, 1},
		{"tab-eingerueckter code geaendert", "Text\n\n\t[a](a.md)", "Text\n\n\t[a](x/a.md)", 1, 1},
		{"bild hinter escaptem ausrufezeichen nachgezogen", `\![B](a.md)`, `\![B](x/a.md)`, 0, 1},
		{"code-span im linktext nachgezogen",
			"Siehe [`releasing.md`](../../user/releasing.md).",
			"Siehe [`releasing.md`](../../user/maintainer/releasing.md).", 0, 1},
		{"leerzeichen-tab-einzug geaendert", "Text\n\n \t[a](a.md)", "Text\n\n \t[a](x/a.md)", 1, 1},
		{"html-block geaendert", "<div>\n[a](a.md)\n</div>", "<div>\n[a](x/a.md)\n</div>", 1, 1},
		{"pre-block mit leerzeile geaendert", "<pre>\n\n[a](a.md)\n</pre>", "<pre>\n\n[a](x/a.md)\n</pre>", 1, 1},
		{"link nach html-block nachgezogen", "<div>x</div>\n\n[R](a.md)", "<div>x</div>\n\n[R](x/a.md)", 0, 1},
		{"escapte schliessende klammer geaendert", `[Ausnahme\](a.md)`, `[Ausnahme\](x/a.md)`, 1, 1},
		// fail-safe: was die Erkennung nicht als Link liest, bleibt Drift
		{"escapte klammer im linktext nachgezogen", `[a\]b](a.md)`, `[a\]b](x/a.md)`, 1, 1},
		{"ziel auf der folgezeile nachgezogen", "[R](\na.md)", "[R](\nx/a.md)", 1, 1},
		// Klammertext ohne gültigen Zielausdruck ist kein Link
		{"klammertext mit leerraum geaendert", "[B](a.md empfohlen)", "[B](x/a.md empfohlen)", 1, 1},
		{"offener titel geaendert", `[B](a.md "offen)`, `[B](x/a.md "offen)`, 1, 1},
		{"spitzklammer mit rest geaendert", "[B](<a.md> nicht)", "[B](<x/a.md> nicht)", 1, 1},
		{"tabellenzelle mit pipe geaendert", "| [a](b | nicht) |", "| [a](b | doch) |", 1, 1},
		{"escapte zielklammer geaendert", `[B](a.md\) gilt)`, `[B](x/a.md\) gilt)`, 1, 1},
		// Container und Code nach Markdown
		{"code im zitat geaendert", ">     [a](a.md)", ">     [a](x/a.md)", 1, 1},
		{"code im listenpunkt geaendert", "-     [a](a.md)", "-     [a](x/a.md)", 1, 1},
		{"tilde-fence im zitat geaendert", "> ~~~\n> [a](a.md)\n> ~~~", "> ~~~\n> [a](x/a.md)\n> ~~~", 1, 1},
		{"eingerueckter fence verschiebt nichts",
			"    ```\n\n```\n[a](a.md)\n```", "    ```\n\n```\n[a](x/a.md)\n```", 1, 1},
		{"backticks ueber listenzeilen geaendert", "- `a\n- `[b](a.md)`", "- `a\n- `[b](x/a.md)`", 1, 1},
		{"code-span bis in den linktext geaendert", "- `a\n- `[b`](a.md)", "- `a\n- `[b`](x/a.md)", 1, 1},
		{"html-kommentar mit leerzeile geaendert", "<!--\n\n[a](a.md)\n-->", "<!--\n\n[a](x/a.md)\n-->", 1, 1},
		{"link nach einzeiligem kommentar nachgezogen", "<!-- x -->\n\n[R](a.md)", "<!-- x -->\n\n[R](x/a.md)", 0, 1},
		{"listenpunkt-link nachgezogen", "- [R](a.md)", "- [R](x/a.md)", 0, 1},
		{"zitat-link nachgezogen (fail-safe)", "> [R](a.md)", "> [R](x/a.md)", 1, 1},
		// nur ein auflösendes Ziel wird normiert, auf Dateiname und Anker
		{"verzeichnis-ziel nachgezogen", "[P](../planning/)", "[P](../neu/planning/)", 0, 1},
		{"nachzug auf anderes dokument", "[R](a.md)", "[R](b.md)", 1, 1},
		{"neues ziel existiert nicht", "[R](a.md)", "[R](x/fehlt.md)", 1, 1},
		{"verzeichnis eines nicht aufloesenden ziels geaendert", "[B](nicht/empfohlen)", "[B](doch/empfohlen)", 1, 1},
		{"anker geaendert", "[R](a.md#alt)", "[R](x/a.md#neu)", 1, 1},
		{"reiner anker geaendert", "[A](#alt)", "[A](#neu)", 1, 1},
		{"externes ziel geaendert", "[G](https://a.example/x.md)", "[G](https://b.example/x.md)", 1, 1},
		{"klammer um inneren link geaendert", "[a [b](a.md) d](a.md)", "[a [b](a.md) d](x/a.md)", 1, 1},
		{"bild im linktext nachgezogen", "[![I](img/a.png)](a.md)", "[![I](img/a.png)](x/a.md)", 0, 1},
		{"pipe ohne leerraum geaendert", "| [a](abgelehnt|v2) |", "| [a](angenommen|v2) |", 1, 1},
		{"link-text im referenz-titel geaendert", `[l]: a.md "[x](a.md)"`, `[l]: a.md "[x](x/a.md)"`, 1, 1},
		{"fence-schliesser anderes zeichen", "```\n~~~\n[a](a.md)\n```", "```\n~~~\n[a](x/a.md)\n```", 1, 1},
		{"fence-schliesser zu kurz", "````\n```\n[a](a.md)\n````", "````\n```\n[a](x/a.md)\n````", 1, 1},
		// eine Absatz-Folgezeile in Referenz-Form: ein Inhaltswort löst nicht auf
		{"absatz-folgezeile in referenz-form", "Text\n[Status]: Abgelehnt.", "Text\n[Status]: Angenommen.", 1, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, mode := range []struct {
				schalter bool
				want     int
			}{{true, c.mitSchalter}, {false, c.ohne}} {
				cfg := adrConfig()
				cfg.IgnoreLinkTargets = mode.schalter
				files := refs(adr("Accepted", c.base), adr("Accepted", c.head))
				for _, ref := range []string{"BASE", "HEAD"} {
					for _, p := range nachzugBaum() {
						files[ref][p] = nil
					}
				}
				fv := &fakeVCS{files: files}
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

// nachzugBaum liefert die Dateien, auf die die Links der Fälle in beiden Ständen
// zeigen dürfen: ein Ziel löst nur auf, wenn es hier steht.
func nachzugBaum() []string {
	return []string{
		"docs/plan/adr/a.md", "docs/plan/adr/x/a.md",
		"docs/plan/adr/b.md", "docs/plan/adr/x/b.md",
		"docs/plan/adr/img/a.png", "docs/plan/adr/x/img/a.png",
		"docs/plan/adr/a(1).md", "docs/plan/adr/x/a(1).md",
		"docs/plan/adr/alt/a b.md", "docs/plan/adr/neu/a b.md",
		"docs/user/releasing.md", "docs/user/maintainer/releasing.md",
		"docs/plan/planning/p.md", "docs/plan/neu/planning/p.md",
	}
}

// Ein Umzug hat in BASE und HEAD verschiedene Bäume: der Nachzug geht durch,
// ein unveränderter Link auf eine gelöschte Datei ist keine Drift, und ein
// neues Ziel, das nur Dateiname und Anker eines aufgelösten alten trägt, ist
// eine.
func TestVCSIgnoreLinkTargetsUmzug(t *testing.T) {
	baseTree := []string{"docs/user/releasing.md", "spec/s.md"}
	headTree := []string{"docs/user/maintainer/releasing.md", "spec/s.md"}
	cases := []struct {
		name       string
		base, head string
		want       int
	}{
		{"nachzug nach umzug", "[R](../../user/releasing.md#prep)", "[R](../../user/maintainer/releasing.md#prep)", 0},
		{"unveraenderter link auf geloeschte datei", "[R](../../user/releasing.md)", "[R](../../user/releasing.md)", 0},
		{"neues ziel nur dateiname", "[S](../../../spec/s.md#a)", "[S](s.md#a)", 1},
		{"neues ziel fehlt in beiden staenden", "[R](../../user/releasing.md)", "[R](../../user/weg/releasing.md)", 1},
		{"neues ziel traegt die marke", "[R](../../user/releasing.md)", "[R](\x00releasing.md)", 1},
		// die benannte Grenze: ein Ziel, das nur noch in BASE existiert, löst auf
		{"nachzug auf datei nur in base", "[R](../../user/maintainer/releasing.md)", "[R](../../user/releasing.md)", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := adrConfig()
			cfg.IgnoreLinkTargets = true
			files := refs(adr("Accepted", c.base), adr("Accepted", c.head))
			for _, p := range baseTree {
				files["BASE"][p] = nil
			}
			for _, p := range headTree {
				files["HEAD"][p] = nil
			}
			got, err := CheckVCS(&fakeVCS{files: files}, cfg, "BASE", "HEAD")
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if len(got) != c.want {
				t.Fatalf("Befunde = %d, want %d (%v)", len(got), c.want, got)
			}
		})
	}
}
