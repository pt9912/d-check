package rules

import (
	"fmt"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// DC-FA-REF-001 — das geteilte Ventil erreicht `links` UND `anchors`
// (querschnittlich): das modul-lokale codepaths.ignore-refs konnte das
// NIE. Ohne Ventil: target-missing (links) + anchor-missing (anchors);
// mit Top-Level ignore-refs auf die Quelldatei sind beide still.
func TestRefsGeteiltesVentilLinksUndAnchors(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/tpl/a.md":    "Link: [x](fehlt.md) und Anker: [y](ziel.md#gibt-es-nicht).\n",
		"docs/tpl/ziel.md": "# Da\n",
	})
	plain, err := Run(m, nil, model.Config{}, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plain.Findings) != 2 {
		t.Fatalf("ohne Ventil: %d Befunde, want 2 (target-missing + anchor-missing)", len(plain.Findings))
	}
	cfg := model.Config{IgnoreRefs: []model.IgnoreRef{{In: "docs/tpl/**", Refs: []string{"docs/tpl/**"}}}}
	ign, err := Run(m, nil, cfg, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ign.Findings) != 0 {
		var got []string
		for _, f := range ign.Findings {
			got = append(got, fmt.Sprintf("%s %s %s", f.Rule, f.Target, f.Reason))
		}
		t.Fatalf("mit geteiltem Ventil: %d Befunde, want 0: %v", len(ign.Findings), got)
	}
}

// DC-FA-REF-001 — `keep` gewinnt reihenfolge-unabhängig: refs=tpl/x/**
// würde alle drei Ziele ignorieren, keep=tpl/x/keep-* holt zwei zurück.
// Ein zurückgeholtes REALES Ziel bleibt grün, ein zurückgeholter
// TIPPFEHLER feuert; das nicht-gekeepte Ziel ist still. Erwartet: genau
// ein target-missing (keep-typo). Ohne `keep` wäre der CR nicht
// abgenommen — dieser Test verriegelt es.
func TestRefsKeepGewinnt(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/tpl/a.md":         "[a](x/keep-me.md) [b](x/drop.md) [c](x/keep-typo.md)\n",
		"docs/tpl/x/keep-me.md": "real",
	})
	cfg := model.Config{IgnoreRefs: []model.IgnoreRef{{
		In:   "docs/tpl/**",
		Refs: []string{"docs/tpl/x/**"},
		Keep: []string{"docs/tpl/x/keep-*"},
	}}}
	res, err := Run(m, nil, cfg, []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s", f.Target, f.Reason))
	}
	want := []string{"x/keep-typo.md target-missing"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("keep: Befunde = %v\nwant %v (drop ignoriert, keep-me real→grün, keep-typo feuert)", got, want)
	}
}

// DC-FA-REF-001 — der Anker bleibt scharf: ein per `keep` zurückgeholtes,
// reales Ziel behält die Anker-Prüfung (ein toter Anker feuert), während der
// Platzhalter still ist. Verriegelt, dass das Ventil Anker nicht pauschal
// unterdrückt — der Kern des CR (umbenannte Überschrift → toter Anker im
// ausgelieferten Artefakt).
func TestRefsAnkerBleibtScharf(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/tpl/a.md":    "Platzhalter: [p](weg.md) — real, toter Anker: [k](real.md#gibt-es-nicht)\n",
		"docs/tpl/real.md": "# Titel\n",
	})
	cfg := model.Config{IgnoreRefs: []model.IgnoreRef{{
		In:   "docs/tpl/**",
		Refs: []string{"docs/tpl/**"},
		Keep: []string{"docs/tpl/real.md"},
	}}}
	res, err := Run(m, nil, cfg, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s %s", f.Rule, f.Target, f.Reason))
	}
	want := []string{"anchors real.md#gibt-es-nicht anchor-missing"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Anker-Test: Befunde = %v\nwant %v (Platzhalter ignoriert, kept-Ziel behält Anker-Prüfung)", got, want)
	}
}

// DC-FA-REF-001 — das TOP-LEVEL-Ventil (nicht nur der Alias
// `codepaths.ignore-refs`) erreicht auch `codepaths`. Pinnt die Kombination in
// CheckCodepaths (`refs := ignoreRefs` + Alias-Anhang): die Mutation `refs := nil`
// ließ sonst die ganze Suite grün (Review R1-F-1).
func TestRefsGeteiltesVentilCodepaths(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/tpl/a.md": "Inline-Pfad-Platzhalter: `docs/tpl/weg.py`.\n",
	})
	cfg := model.Config{
		Codepaths:  model.CodepathsConfig{Roots: []string{"docs"}},
		IgnoreRefs: []model.IgnoreRef{{In: "docs/tpl/**", Refs: []string{"docs/tpl/**"}}},
	}
	res, err := Run(m, nil, cfg, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Top-Level ignore-refs muss codepaths erreichen: %d Befunde, want 0", len(res.Findings))
	}
	// Gegenprobe: ohne Top-Level-Ventil fällt der codepath-missing.
	res2, err := Run(m, nil, model.Config{Codepaths: model.CodepathsConfig{Roots: []string{"docs"}}}, []string{"codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2.Findings) != 1 {
		t.Fatalf("Gegenprobe ohne Ventil: want 1 codepath-missing, got %d", len(res2.Findings))
	}
}

// DC-FA-REF-001 — der Quell-Skopus `in` isoliert: dasselbe Ziel-Muster
// bleibt in einer Datei AUSSERHALB des in-Globs voll geprüft. Beide
// Dateien referenzieren dasselbe aufgelöste Ziel (ziel/fehlt.md); nur
// die in `tpl/` ist still.
func TestRefsSkopusIsolation(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/tpl/a.md":   "[a](../ziel/fehlt.md)\n",
		"docs/other/b.md": "[a](../ziel/fehlt.md)\n",
	})
	cfg := model.Config{IgnoreRefs: []model.IgnoreRef{{In: "docs/tpl/**", Refs: []string{"docs/ziel/**"}}}}
	res, err := Run(m, nil, cfg, []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s %s", f.File, f.Reason))
	}
	want := []string{"docs/other/b.md target-missing"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Skopus-Isolation: Befunde = %v\nwant %v (tpl still, other geprüft)", got, want)
	}
}

// DC-FA-REF-001 — Top-Level ignore-refs und der Alias codepaths.ignore-refs
// koexistieren additiv: der Alias wirkt weiter (nur codepaths), das
// Top-Level-Ventil zusätzlich für links. Verriegelt, dass der Alias nicht
// aus Versehen auf links durchschlägt (er ist codepaths-skopiert).
func TestRefsAliasKoexistenz(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "Code-Pfad: `docs/weg.md` und Link: [l](weg2.md)\n",
	})
	cfg := model.Config{
		Codepaths:  model.CodepathsConfig{Roots: []string{"docs"}, IgnoreRefs: []string{"docs/weg.md"}},
		IgnoreRefs: []model.IgnoreRef{{Refs: []string{"docs/weg2.md"}}},
	}
	res, err := Run(m, nil, cfg, []string{"links", "codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	// Alias unterdrückt den codepaths-Befund für docs/weg.md; Top-Level
	// unterdrückt den links-Befund für docs/weg2.md → 0 Befunde.
	if len(res.Findings) != 0 {
		var got []string
		for _, f := range res.Findings {
			got = append(got, fmt.Sprintf("%s %s %s", f.Rule, f.Target, f.Reason))
		}
		t.Fatalf("Alias+Top-Level: %d Befunde, want 0: %v", len(res.Findings), got)
	}
	// Gegenprobe: der Alias ist codepaths-skopiert — er unterdrückt den
	// links-Befund NICHT. Ohne Top-Level bleibt der links-Befund.
	cfgAliasOnly := model.Config{
		Codepaths: model.CodepathsConfig{Roots: []string{"docs"}, IgnoreRefs: []string{"docs/weg.md", "docs/weg2.md"}},
	}
	res2, err := Run(m, nil, cfgAliasOnly, []string{"links", "codepaths"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2.Findings) != 1 || res2.Findings[0].Rule != "links" {
		t.Fatalf("Alias-Skopus: want genau 1 links-Befund (Alias erreicht links nicht), got %d", len(res2.Findings))
	}
}
