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
