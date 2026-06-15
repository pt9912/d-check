package core

import (
	"fmt"
	"testing"
)

// DC-FA-ANCH-001.a — Slug-Schritte gegen reale Heading-Formen dieses
// Projekts (inkl. Gedankenstrich → Doppel-Bindestrich, Umlaute,
// Code-Spans: Text bleibt, Backticks entfallen).
func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"DC-QA-01 — Performance":                  "dc-qa-01--performance",
		"DC-FA-CLI-002.a — Modul-Auflösung":       "dc-fa-cli-002a--modul-auflösung",
		"7. Closure-Notiz (nach `done/`)":         "7-closure-notiz-nach-done",
		"Die `.d-check.yml`":                      "die-d-checkyml",
		"Grund- und Fehler-Codes":                 "grund--und-fehler-codes",
		"[Titel](x.md) im Link":                   "titel-im-link",
		"snake_case bleibt":                       "snake_case-bleibt",
		"*Betont* und **fett**":                   "betont-und-fett",
		"MR-003 — Sensor `tools/verify-doc.sh`":   "mr-003--sensor-toolsverify-docsh",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHeadingSlugs_DuplikateUndFences(t *testing.T) {
	content := []byte("# Beispiel\n## Beispiel\n```\n# imfence\n```\n### Beispiel\n#ohneblank\n####### zu tief\n")
	slugs := HeadingSlugs(content)
	for _, want := range []string{"beispiel", "beispiel-1", "beispiel-2"} {
		if !slugs[want] {
			t.Errorf("Slug %q fehlt: %v", want, slugs)
		}
	}
	for _, no := range []string{"imfence", "ohneblank", "zu-tief", "beispiel-3"} {
		if slugs[no] {
			t.Errorf("Slug %q dürfte nicht existieren", no)
		}
	}
}

// DC-FA-ANCH-001 Happy/Boundary/Negative über Run.
func TestAnchorsModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "# Zweck und Geltungsbereich\n## Beispiel\n## Beispiel\n" +
			"[ok](#zweck-und-geltungsbereich)\n" +
			"[dup](#beispiel-1)\n" +
			"[fehlt](#gibt-es-nicht)\n" +
			"[quer](b.md#abschnitt-zwei)\n" +
			"[querfehlt](b.md#nope)\n" +
			"[dateifehlt](weg.md#x)\n" +
			"[umlaut](b.md#%C3%BCbersicht)\n" +
			"[nichtmd](c.txt#egal)\n",
		"docs/b.md": "# Abschnitt Zwei\n# Übersicht\n",
		"docs/c.txt": "kein markdown",
	})
	res, err := Run(m, nil, Config{}, []string{"anchors"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s", f.Line, f.Target, f.Reason))
	}
	want := []string{
		"6 #gibt-es-nicht anchor-missing",
		"8 b.md#nope anchor-missing",
	}
	if len(got) != len(want) {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Befund %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// Genau EIN Befund (aus links) bei fehlender Zieldatei mit Fragment —
// anchors schweigt (DC-FA-ANCH-001 Boundary).
func TestAnchorsSchweigtBeiFehlenderDatei(t *testing.T) {
	m := newMemFS(map[string]string{"docs/a.md": "[x](weg.md#frag)"})
	res, err := Run(m, nil, Config{}, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Rule != "links" {
		t.Fatalf("Befunde = %+v (genau einer, aus links, erwartet)", res.Findings)
	}
}

// scan.roots "." = gesamte Repo-Wurzel inkl. Nicht-Default-Verzeichnissen.
func TestDiscoverFiles_PunktWurzel(t *testing.T) {
	m := newMemFS(map[string]string{
		"harness/h.md":  "x",
		"docs/a.md":     "x",
		"README.md":     "x",
		"node_modules/x/y.md": "x",
	})
	files, err := DiscoverFiles(m, []string{"."}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"harness/h.md": true, "docs/a.md": true, "README.md": true}
	if len(files) != 3 {
		t.Fatalf("files = %v", files)
	}
	for _, f := range files {
		if !want[f] {
			t.Fatalf("unerwartete Datei %q in %v", f, files)
		}
	}
}

// DC-FA-ANCH-001.b — Inline-HTML-Anker: id an beliebigem Element,
// name nur an <a>; wörtlich; außerhalb Fence/Inline-Code; data-id und
// name an <area> zählen nicht.
func TestHTMLAnchors(t *testing.T) {
	content := []byte(
		"# Heading Eins\n" +
			"<a name=\"namea\"></a>\n" +
			"<div id=\"DivId\">x</div>\n" +
			"<span id='singleid'></span>\n" +
			"<a name='singlename'></a>\n" +
			"<area id=\"areaid\">\n" +
			"<area name=\"areaname\"></area>\n" +
			"<div data-id=\"datid\">x</div>\n" +
			"`<a id=\"inlineid\"></a>`\n" +
			"```\n<a id=\"fenceid\"></a>\n```\n",
	)
	got := htmlAnchors(content)
	for _, want := range []string{"namea", "DivId", "singleid", "singlename", "areaid"} {
		if !got[want] {
			t.Errorf("HTML-Anker %q fehlt: %v", want, got)
		}
	}
	for _, no := range []string{"areaname", "datid", "inlineid", "fenceid", "heading-eins"} {
		if got[no] {
			t.Errorf("HTML-Anker %q dürfte nicht existieren: %v", no, got)
		}
	}
	set := AnchorSet(content)
	if !set["heading-eins"] || !set["namea"] || !set["DivId"] {
		t.Errorf("AnchorSet vereinigt Heading-Slugs und HTML-Anker nicht: %v", set)
	}
}

// DC-FA-ANCH-001 (HTML-Anker) Happy/Boundary/Negative über Run:
// Heading-Slug, <a name>, id; Negativ: name an <area>, Inline-Code,
// Fence, Case-Mismatch (HTML wörtlich/case-sensitiv).
func TestAnchorsHTMLModul(t *testing.T) {
	m := newMemFS(map[string]string{
		"docs/a.md": "[h](b.md#echtes-heading)\n" +
			"[na](b.md#abschnitt-2)\n" +
			"[iddiv](b.md#Übersicht)\n" +
			"[aid](b.md#zonen-id)\n" +
			"[an](b.md#zone)\n" +
			"[inl](b.md#inline-phantom)\n" +
			"[fen](b.md#fence-phantom)\n" +
			"[low](b.md#übersicht)\n",
		"docs/b.md": "# Echtes Heading\n" +
			"<a name=\"abschnitt-2\"></a>\n" +
			"<div id=\"Übersicht\">Inhalt</div>\n" +
			"<area name=\"zone\"></area>\n" +
			"<area id=\"zonen-id\">\n" +
			"`<div id=\"inline-phantom\">`\n" +
			"```\n<div id=\"fence-phantom\">\n```\n",
	})
	res, err := Run(m, nil, Config{}, []string{"anchors"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s", f.Line, f.Target, f.Reason))
	}
	want := []string{
		"5 b.md#zone anchor-missing",
		"6 b.md#inline-phantom anchor-missing",
		"7 b.md#fence-phantom anchor-missing",
		"8 b.md#übersicht anchor-missing",
	}
	if len(got) != len(want) {
		t.Fatalf("Befunde = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Befund %d = %q, want %q", i, got[i], want[i])
		}
	}
}
