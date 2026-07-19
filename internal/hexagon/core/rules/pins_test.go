package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// pinHashOf berechnet den erwarteten Span-Hash über denselben Produktionspfad
// (spanHash), damit Fixtures keinen hartkodierten Hash brauchen.
func pinHashOf(t *testing.T, fix map[string]string, rel, anchor string) string {
	t.Helper()
	h, ok := spanHash(coretest.NewMemFS(fix), rel, anchor, map[string]string{})
	if !ok {
		t.Fatalf("spanHash(%s#%s) nicht auflösbar", rel, anchor)
	}
	return h
}

// DC-FA-PIN-001 Happy: Pin == normalisierter Ziel-Span-Hash → kein Befund.
func TestPinsHappy(t *testing.T) {
	fix := map[string]string{"spec/ziel.md": "## Abschnitt\n\nStabiler Inhalt hier.\n"}
	h := pinHashOf(t, fix, "spec/ziel.md", "abschnitt")
	fix["docs/src.md"] = "Siehe [Abschnitt](../spec/ziel.md#abschnitt) <!-- dpin: sha256:" + h + " -->.\n"
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("matching pin → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 Boundary (Reflow): nur-Whitespace/Umbruch-Änderung am Ziel-Span
// (Wort-Inhalt gleich) → kein Befund (Normalisierung absorbiert Reflow).
func TestPinsReflow(t *testing.T) {
	compact := map[string]string{"spec/z.md": "## A\n\nHallo Welt foo bar.\n"}
	h := pinHashOf(t, compact, "spec/z.md", "a")
	fix := map[string]string{
		"spec/z.md": "## A\n\nHallo   Welt\nfoo    bar.\n", // reflowed, gleiche Worte
		"docs/s.md": "[A](../spec/z.md#a) <!-- dpin: sha256:" + h + " -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Reflow absorbiert → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 Negative: inhaltliche Änderung am Ziel-Span → link-stale.
func TestPinsStale(t *testing.T) {
	orig := map[string]string{"spec/z.md": "## A\n\nAlter Inhalt.\n"}
	h := pinHashOf(t, orig, "spec/z.md", "a")
	fix := map[string]string{
		"spec/z.md": "## A\n\nNeuer, geänderter Inhalt.\n",
		"docs/s.md": "[A](../spec/z.md#a) <!-- dpin: sha256:" + h + " -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Rule != "pins" || res.Findings[0].Reason != model.ReasonLinkStale {
		t.Fatalf("inhaltliche Drift → 1× link-stale, got %v", res.Findings)
	}
}

// slice-081 (dpin-Ergonomie): der link-stale-Befund trägt den VOLLEN errechneten
// Ist-Hash (Re-Pin-Vorlage), nicht nur shortHash — sonst ist dpin nicht adoptierbar.
// Mutations-echt: mit shortHash(got) enthält die Message den vollen Hash NICHT.
func TestPinsStaleMessageFullHash(t *testing.T) {
	orig := map[string]string{"spec/z.md": "## A\n\nAlter Inhalt.\n"}
	h := pinHashOf(t, orig, "spec/z.md", "a")
	fix := map[string]string{
		"spec/z.md": "## A\n\nNeuer, geänderter Inhalt.\n",
		"docs/s.md": "[A](../spec/z.md#a) <!-- dpin: sha256:" + h + " -->\n",
	}
	fullGot := pinHashOf(t, fix, "spec/z.md", "a") // voller Ist-Hash des gedrifteten Spans
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("Drift → 1 Befund, got %v", res.Findings)
	}
	if !strings.Contains(res.Findings[0].Message, fullGot) {
		t.Fatalf("link-stale-Message muss den vollen Ist-Hash %q tragen (Re-Pin-Vorlage), got %q", fullGot, res.Findings[0].Message)
	}
}

// DC-FA-PIN-001 Boundary (Marker-Bindung): zwei Links, je ein unmittelbar
// folgender Marker → beide gebunden; ein Folgezeilen-Marker ist inert. Hier:
// a ok, b stale → genau ein Befund (b).
func TestPinsMarkerBindung(t *testing.T) {
	fix := map[string]string{
		"spec/a.md": "## X\n\nInhalt A.\n",
		"spec/b.md": "## Y\n\nInhalt B.\n",
	}
	ha := pinHashOf(t, fix, "spec/a.md", "x")
	fix["docs/s.md"] = "[a](../spec/a.md#x) <!-- dpin: sha256:" + ha + " --> [b](../spec/b.md#y) <!-- dpin: sha256:deadbeef -->\n" +
		"<!-- dpin: sha256:cafebabe -->\n" // Folgezeile, keinem Link folgend → inert
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s", f.Target, f.Reason))
	}
	want := []string{"../spec/b.md#y link-stale"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Marker-Bindung: nur b stale (a ok, Folgezeile inert)\n got %v\nwant %v", got, want)
	}
}

// DC-FA-PIN-001 Boundary (Ziel weg): gepinnter Link mit fehlendem Ziel → im
// pins-only-Lauf kein Befund; mit aktivem links genau ein target-missing
// (kein Doppelbefund durch pins).
func TestPinsZielWeg(t *testing.T) {
	fix := map[string]string{"docs/s.md": "[x](../spec/fehlt.md#a) <!-- dpin: sha256:deadbeef -->\n"}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Ziel weg, pins-only → 0 Befunde, got %v", res.Findings)
	}
	res2, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"links", "pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2.Findings) != 1 || res2.Findings[0].Rule != "links" || res2.Findings[0].Reason != model.ReasonTargetMissing {
		t.Fatalf("Ziel weg + links: genau 1 target-missing (links), kein link-stale, got %v", res2.Findings)
	}
}

// DC-FA-PIN-001 Modul-aus: ohne aktives pins kein link-stale.
func TestPinsNichtAktiv(t *testing.T) {
	fix := map[string]string{
		"spec/z.md": "## A\n\nInhalt.\n",
		"docs/s.md": "[A](../spec/z.md#a) <!-- dpin: sha256:deadbeef -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("pins nicht aktiv → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 / DC-FA-CONF-002: pins.scope begrenzt die gescannten
// Quell-Dateien — nur die in-scope-Datei wird auf Pins geprüft.
func TestPinsScope(t *testing.T) {
	fix := map[string]string{
		"spec/z.md":     "## A\n\nInhalt.\n",
		"docs/in/s.md":  "[A](../../spec/z.md#a) <!-- dpin: sha256:deadbeef -->\n",
		"docs/out/s.md": "[A](../../spec/z.md#a) <!-- dpin: sha256:deadbeef -->\n",
	}
	cfg := model.Config{Scopes: map[string]*model.ScopeConfig{"pins": {Roots: []string{"docs/in"}}}}
	res, err := Run(coretest.NewMemFS(fix), nil, cfg, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].File != "docs/in/s.md" {
		t.Fatalf("pins.scope: nur docs/in geprüft, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 (Impl-R1 F-1): Inline-Code zwischen Link und Marker ist
// Nicht-Whitespace → der Marker ist inert (Prüfung auf der ROHEN Zeile, nicht
// der vorverarbeiteten, wo der Code zu Leerzeichen würde).
func TestPinsInlineCodeZwischenInert(t *testing.T) {
	fix := map[string]string{
		"spec/z.md": "Inhalt.\n",
		"docs/s.md": "[a](../spec/z.md) `code` <!-- dpin: sha256:deadbeef -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Code zwischen Link und Marker → inert, 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 (Impl-R1 F-2): ein Pin an einem externen Link ist inert.
func TestPinsExternInert(t *testing.T) {
	fix := map[string]string{"docs/s.md": "[x](https://example.com/p) <!-- dpin: sha256:deadbeef -->\n"}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("externes Ziel → kein pins-Befund, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 (Impl-R1 F-2): ein Ziel außerhalb der Repo-Wurzel (repo-escape)
// erzeugt keinen pins-Befund (Out-of-Scope; struktureller Befund bleibt links).
func TestPinsRepoEscape(t *testing.T) {
	fix := map[string]string{"docs/s.md": "[x](../../outside.md) <!-- dpin: sha256:deadbeef -->\n"}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("repo-escape-Ziel → kein pins-Befund, got %v", res.Findings)
	}
}

// DC-FA-PIN-001 (Impl-R1 F-2): Same-file-Anker-Pin (Linkziel `#anker`) → der Span
// ist die Heading-Section derselben Datei; Drift → link-stale.
func TestPinsSameFileAnker(t *testing.T) {
	fix := map[string]string{
		"docs/s.md": "## A\n\nInhalt A.\n\n## B\n\n[x](#a) <!-- dpin: sha256:deadbeef -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "#a" || res.Findings[0].Reason != model.ReasonLinkStale {
		t.Fatalf("Same-file-Anker-Drift → 1× link-stale (#a), got %v", res.Findings)
	}
}

// DC-FA-PIN-001 (Impl-R1 F-2): ein Anker, der keine Section trifft, ist nicht
// auflösbar → kein pins-Befund (struktureller Befund bleibt anchors).
func TestPinsAnkerOhneSection(t *testing.T) {
	fix := map[string]string{
		"spec/z.md": "## A\n\nInhalt.\n",
		"docs/s.md": "[x](../spec/z.md#fehlt) <!-- dpin: sha256:deadbeef -->\n",
	}
	res, err := Run(coretest.NewMemFS(fix), nil, model.Config{}, []string{"pins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Anker ohne Section → kein pins-Befund, got %v", res.Findings)
	}
}
