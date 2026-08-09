package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
)

// DC-FA-SPAN-001: Happy/Boundary/Negative gegen In-Memory-FS.
func TestSpansModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		// Happy: balancierte Spans, auch mehrzeilig über den Umbruch
		"docs/ok.md": "ein `span` und ein (`u-boot init &&\nu-boot up`, `noch einer`) fertig",
		// Boundary: alleinstehende literale Backticks (beidseitig
		// Whitespace bzw. Zeilenrand)
		"docs/literal.md": "ein ` allein\nund am Ende `",
		// Negative: Opener klebt an Nicht-Whitespace, im Absatz
		// ungeschlossen (u-boot-Titel-Klasse)
		"docs/kaputt.md": "- `feat(cli): u-boot logs --json\n  ([`LH-X-001`](spec.md#x) /\n\nnächster Absatz ohne Befund",
		"docs/spec.md":   "# x",
	})

	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"spans"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("Befunde = %+v, want genau 1 span-unclosed", res.Findings)
	}
	// Das Pairing macht den Titel-Opener zum geschlossenen Span
	// (er paart mit dem nächsten Backtick); der Befund zeigt auf die
	// übrigbleibende ungeschlossene Folge des Absatzes — Zeile 2.
	f := res.Findings[0]
	if f.File != "docs/kaputt.md" || f.Line != 2 || f.Reason != model.ReasonSpanUnclosed {
		t.Fatalf("Befund = %+v", f)
	}
	// Ziel: Backtick-Folge + folgende Nicht-Whitespace-Zeichen, ≤ 30
	if !strings.HasPrefix(f.Target, "`](spec.md#x)") || len(f.Target) > 30 {
		t.Fatalf("Target = %q", f.Target)
	}
}

// span-nested-link: Link-Syntax im Linktext eines weiteren Links;
// benachbarte eigenständige Links sind kein Treffer.
func TestSpansNestedLink(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "[[innen](x.md)](y.md) kaputt\n[a](x.md)[b](y.md) ok\nin `code: ](x)](` kein Befund",
		"docs/x.md": "x",
		"docs/y.md": "y",
	})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"spans"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("Befunde = %+v, want genau 1 span-nested-link", res.Findings)
	}
	f := res.Findings[0]
	if f.Line != 1 || f.Reason != model.ReasonSpanNestedLink || !strings.HasPrefix(f.Target, "](x.md)](") {
		t.Fatalf("Befund = %+v", f)
	}
}

// Bildreferenz als Linktext (Badge-Muster) ist legales Markdown —
// kein span-nested-link (Lastenheft 0.7.1, Kalibrierungs-Befund).
func TestSpansBadgeKeinTreffer(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "[![Build](badge.svg)](https://ci.example.org) Status",
		"docs/badge.svg": "x",
	})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"spans"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Badge-Muster geflaggt: %+v", res.Findings)
	}
}

// Mehrzeilige Spans, deren Opener an Text klebt und die im Absatz
// geschlossen werden, sind KEIN Befund — nur die Parität zählt.
func TestSpansMehrzeiligGeschlossen(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Befehl (`u-boot stop\npostgres && u-boot up`, `x`) Ende",
	})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}}, []string{"spans"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("False-Positive auf legalem mehrzeiligem Span: %+v", res.Findings)
	}
}

// fence-unclosed (§DC-FA-SPAN-001.a Schritt 3): die Block-Ebene von
// span-unclosed. Geprüft werden beide Schluss-Lesarten, die Wahl der gemeldeten
// Zeile und die Fälle, die kein Befund sind.
func TestFenceUnclosed(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantLine int // 0 = kein Befund
	}{
		{"balanciert", "# T\n\n```\ncode\n```\n\nProsa.\n", 0},
		{"mehrere balancierte Fences", "```\na\n```\ntext\n```go\nb\n```\n", 0},
		{"gar kein Fence", "# T\n\nNur Prosa.\n", 0},
		{"offen am Dateiende", "# T\n\n```\ncode\n", 3},
		{"offen, mit Inhalt dahinter", "# T\n\n```\ncode\n\nnoch mehr\n", 3},
		{"zweiter Fence offen", "```\na\n```\ntext\n```go\nb\n", 5},
		{"Tilde-Fence offen", "# T\n\n~~~\ncode\n", 3},
		// Die Infozeilen-Regel aus ADR-0042 gilt weiter: eine ```-Zeile mit
		// Backtick im Rest ist KEIN Öffner, also auch kein offener Fence.
		{"Backtick in der Infozeile ist kein Öffner", "# T\n\n``` `x` ```\n\nProsa.\n", 0},
		// Space/Tab-Einrückung ist für proseLines kein Hindernis, also
		// zählt die Zeile auch hier als Fence.
		{"eingerückter Fence offen", "# T\n\n  ```\ncode\n", 3},
		// Unicode-Whitespace ist für proseLines keine Einrückung und darf
		// deshalb auch hier keine sein; gemeldet wird der echte Öffner in
		// Zeile 7.
		{"U+00A0 vor der Fence ist kein Toggle", "a\n\n\u00a0```\n\nb\n\n```yaml\nx\n", 7},
		// Naiver Toggle balanciert, CommonMark-Schluss offen: ein ```-Block
		// wird von ~~~ nicht geschlossen (und umgekehrt), ein ````-Block
		// nicht von einer kürzeren Folge.
		{"Tilde schließt keinen Backtick-Block", "# T\n\n```\nx\n~~~\n\nProsa.\n", 3},
		{"Backtick schließt keinen Tilde-Block", "# T\n\n~~~\nx\n```\n\nProsa.\n", 3},
		{"kürzerer Schluss schließt nicht", "# T\n\n````\nx\n```\n\nProsa.\n", 3},
		// Legale Verschachtelung: der äußere ````-Block zeigt einen
		// ```-Block; unter der strengen Lesart ist alles geschlossen.
		{"legale Verschachtelung bleibt grün", "````markdown\n```yaml\nx\n```\n````\n\nProsa.\n", 0},
		// Umgekehrte Richtung: der ````-Block zeigt EINE ```-Zeile. Für die
		// strenge Lesart ist er geschlossen, für den naiven Toggle bleibt die
		// Parität gekippt — und die Vorverarbeitung überspringt den Rest.
		{"nur die naive Lesart kippt", "````markdown\n```\n````\n\nProsa.\n", 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := checkUnclosedFence("d.md", []byte(tc.content))
			if tc.wantLine == 0 {
				if len(got) != 0 {
					t.Fatalf("erwartet kein Befund, got %+v", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("erwartet GENAU EINEN Befund je Datei, got %+v", got)
			}
			if got[0].Reason != model.ReasonFenceUnclosed {
				t.Errorf("Grund = %q, want %q", got[0].Reason, model.ReasonFenceUnclosed)
			}
			if got[0].Line != tc.wantLine {
				t.Errorf("Zeile = %d, want %d", got[0].Line, tc.wantLine)
			}
		})
	}
}

// Der Reproduktionsfall, der die Klasse ausgelöst hat: ein offener Fence im
// Closure-Notiz-Abschnitt eines abgeschlossenen Slice.
func TestFenceUnclosedDeckDenReviewFall(t *testing.T) {
	content := "# S\n\n## 6. Closure-Notiz\n\nEins. Zwei.\n\n```\nunbalanciert\n\nalles gut\n"
	got := checkUnclosedFence("done/slice-001-x.md", []byte(content))
	if len(got) != 1 || got[0].Line != 7 {
		t.Fatalf("erwartet ein fence-unclosed an Zeile 7, got %+v", got)
	}
}

// Das Befund-Ziel nach §DC-FA-SPAN-001.a Schritt 3: die getrimmte
// Öffnungszeile, gekappt auf 30 Runen.
func TestFenceUnclosedZiel(t *testing.T) {
	cases := []struct{ name, line, want string }{
		{"kurze Zeile unverändert", "```yaml", "```yaml"},
		{"eingerückt: links getrimmt", "    ```yaml", "```yaml"},
		{"Trailing-Whitespace weg", "```yaml   ", "```yaml"},
		// Runen, nicht Bytes — Umlaute werden nicht mitten im Zeichen
		// abgeschnitten.
		{"auf 30 Runen gekappt", "```jsonc-mit-ähnlich-länglicher-Infozeile", "```jsonc-mit-ähnlich-längliche"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := checkUnclosedFence("d.md", []byte("# T\n\n"+tc.line+"\ncode\n"))
			if len(got) != 1 {
				t.Fatalf("erwartet genau einen Befund, got %+v", got)
			}
			if got[0].Target != tc.want {
				t.Errorf("Ziel = %q, want %q", got[0].Target, tc.want)
			}
			if n := len([]rune(got[0].Target)); n > 30 {
				t.Errorf("Ziel = %d Runen, want ≤ 30", n)
			}
		})
	}
}

// Verdrahtung: der Befund muss über die Modul-Oberfläche CheckSpans
// herauskommen, nicht nur aus checkUnclosedFence direkt.
func TestCheckSpansVerdrahtetFenceUnclosed(t *testing.T) {
	content := []byte("# T\n\n```\ncode\n")
	got := CheckSpans("d.md", content, PreprocessMarkdown(content))
	for _, f := range got {
		if f.Reason == model.ReasonFenceUnclosed {
			return
		}
	}
	t.Fatalf("CheckSpans meldet kein fence-unclosed — ist der Aufruf verdrahtet? got %+v", got)
}
