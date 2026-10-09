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
		{"fussnote mit pfadwort geaendert", "[^1]: alt.md", "[^1]: neu.md", 1, 1},
		{"prosa in referenz-form geaendert", "[Hinweis]: Verboten ist das.", "[Hinweis]: Erlaubt ist das.", 1, 1},
		{"einzelwort ohne pfadzeichen geaendert", "[Hinweis]: Verboten", "[Hinweis]: Erlaubt", 1, 1},
		{"klammer ohne link geaendert", "f](alt und weiter", "f](neu und weiter", 1, 1},
		{"offener link ohne schliessende klammer", "[f](alt und weiter", "[f](neu und weiter", 1, 1},
		// balancierte Klammern liest die links-Erkennung als Teil von Ziel und Linktext
		{"ziel mit eigener klammer nachgezogen", "[A](a(1).md)", "[A](x/a(1).md)", 0, 1},
		{"linktext mit eckiger klammer nachgezogen", "[a[0]](a.md)", "[a[0]](x/a.md)", 0, 1},
		// Code und Escapes: link-förmiger Text, der kein Link ist, bleibt
		{"ziel in inline-code", "`[R](a.md)`", "`[R](x/a.md)`", 1, 1},
		{"funktionsaufruf in inline-code geaendert", "`fns[k](true)`", "`fns[k](false)`", 1, 1},
		{"escapte klammern geaendert", `\[Ausnahme\](keine)`, `\[Ausnahme\](alle)`, 1, 1},
		{"escapte oeffnende klammer geaendert", `\[Ausnahme](keine)`, `\[Ausnahme](alle)`, 1, 1},
		{"link im codeblock nachgezogen", "```\n[R](a.md)\n```", "```\n[R](x/a.md)\n```", 1, 1},
		{"link neben inline-code nachgezogen", "`x` [R](a.md)", "`x` [R](x/a.md)", 0, 1},
		{"doppelter backslash vor link nachgezogen", `\\[R](a.md)`, `\\[R](x/a.md)`, 0, 1},
		{"code-span ab dem linktext geaendert", "[a `b](Nicht)` c", "[a `b](Immer)` c", 1, 1},
		{"eingerueckter code geaendert", "Text\n\n    [a](Nicht)", "Text\n\n    [a](Immer)", 1, 1},
		{"tab-eingerueckter code geaendert", "Text\n\n\t[a](Nicht)", "Text\n\n\t[a](Immer)", 1, 1},
		{"bild hinter escaptem ausrufezeichen nachgezogen", `\![B](a.md)`, `\![B](x/a.md)`, 0, 1},
		{"code-span im linktext nachgezogen",
			"Siehe [`releasing.md`](../../user/releasing.md).",
			"Siehe [`releasing.md`](../../user/maintainer/releasing.md).", 0, 1},
		{"leerzeichen-tab-einzug geaendert", "Text\n\n \t[a](Nicht)", "Text\n\n \t[a](Immer)", 1, 1},
		{"html-block geaendert", "<div>\n[a](Nicht)\n</div>", "<div>\n[a](Immer)\n</div>", 1, 1},
		{"pre-block mit leerzeile geaendert", "<pre>\n\n[a](Nicht)\n</pre>", "<pre>\n\n[a](Immer)\n</pre>", 1, 1},
		{"link nach html-block nachgezogen", "<div>x</div>\n\n[R](a.md)", "<div>x</div>\n\n[R](x/a.md)", 0, 1},
		{"escapte schliessende klammer geaendert", `[Ausnahme\](keine)`, `[Ausnahme\](alle)`, 1, 1},
		// fail-safe: was die Erkennung nicht als Link liest, bleibt Drift
		{"escapte klammer im linktext nachgezogen", `[a\]b](a.md)`, `[a\]b](x/a.md)`, 1, 1},
		{"ziel auf der folgezeile nachgezogen", "[R](\na.md)", "[R](\nx/a.md)", 1, 1},
		// Klammertext ohne gültigen Zielausdruck ist kein Link
		{"klammertext mit leerraum geaendert", "[B](nicht empfohlen)", "[B](stark empfohlen)", 1, 1},
		{"offener titel geaendert", `[B](a.md "nicht)`, `[B](a.md "doch)`, 1, 1},
		{"spitzklammer mit rest geaendert", "[B](<a.md> nicht)", "[B](<a.md> doch)", 1, 1},
		{"tabellenzelle mit pipe geaendert", "| [a](b | nicht) |", "| [a](b | doch) |", 1, 1},
		{"escapte zielklammer geaendert", `[B](Verboten\) gilt)`, `[B](Erlaubt\) gilt)`, 1, 1},
		// Container und Code nach Markdown
		{"code im zitat geaendert", ">     [a](Nicht)", ">     [a](Immer)", 1, 1},
		{"code im listenpunkt geaendert", "-     [a](Nicht)", "-     [a](Immer)", 1, 1},
		{"tilde-fence im zitat geaendert", "> ~~~\n> [a](Nicht)\n> ~~~", "> ~~~\n> [a](Immer)\n> ~~~", 1, 1},
		{"eingerueckter fence verschiebt nichts",
			"    ```\n\n```\n[a](Nicht)\n```", "    ```\n\n```\n[a](Immer)\n```", 1, 1},
		{"backticks ueber listenzeilen geaendert", "- `a\n- `[b](Nicht)`", "- `a\n- `[b](Immer)`", 1, 1},
		{"code-span bis in den linktext geaendert", "- `a\n- `[b`](Nicht)", "- `a\n- `[b`](Immer)", 1, 1},
		{"html-kommentar mit leerzeile geaendert", "<!--\n\n[a](Nicht)\n-->", "<!--\n\n[a](Immer)\n-->", 1, 1},
		{"link nach einzeiligem kommentar nachgezogen", "<!-- x -->\n\n[R](a.md)", "<!-- x -->\n\n[R](x/a.md)", 0, 1},
		{"listenpunkt-link nachgezogen", "- [R](a.md)", "- [R](x/a.md)", 0, 1},
		{"zitat-link nachgezogen (fail-safe)", "> [R](a.md)", "> [R](x/a.md)", 1, 1},
		// die benannte Grenze: eine Absatz-Folgezeile in Referenz-Form wird geleert
		{"absatz-folgezeile in referenz-form", "Text\n[Status]: Abgelehnt.", "Text\n[Status]: Angenommen.", 0, 1},
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
