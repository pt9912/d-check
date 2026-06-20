package core

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// fakeChecker beantwortet URLs aus einer Tabelle und zählt Aufrufe
// (thread-sicher — checkURLs läuft parallel).
type fakeChecker struct {
	mu      sync.Mutex
	results map[string]driven.HTTPResult
	calls   map[string]int
}

func newFakeChecker(results map[string]driven.HTTPResult) *fakeChecker {
	return &fakeChecker{results: results, calls: map[string]int{}}
}

func (f *fakeChecker) Check(url string) driven.HTTPResult {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls[url]++
	return f.results[url]
}

// panicChecker schlägt fehl, sobald er aufgerufen wird — Beleg der
// Opt-in-Garantie (DC-FA-EXT-001 Boundary / DC-QA-03).
type panicChecker struct{ t *testing.T }

func (p panicChecker) Check(url string) driven.HTTPResult {
	p.t.Fatalf("Netzwerkzugriff ohne aktiviertes Modul external: %s", url)
	return driven.HTTPResult{}
}

// DC-FA-EXT-001 Happy/Negative über Run: Status < 400 ok; 404 →
// external-status; Timeout → external-timeout; > REDIRECT_MAX →
// external-redirects; Transportfehler → external-status. Dedupe:
// genau eine Prüfung pro URL, Befund an jedem Vorkommen.
func TestExternalModul(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "[ok](https://ok.test)\n[kaputt](https://vierxx.test)\n" +
			"[zeit](https://langsam.test)\n[kreis](https://kreis.test)\n" +
			"[weg](https://dns.test)\n",
		"docs/b.md": "[nochmal kaputt](https://vierxx.test)\n",
	})
	checker := newFakeChecker(map[string]driven.HTTPResult{
		"https://ok.test":      {Status: 200},
		"https://vierxx.test":  {Status: 404},
		"https://langsam.test": {Timeout: true},
		"https://kreis.test":   {TooManyRedirects: true},
		"https://dns.test":     {TransportError: "no such host"},
	})
	res, err := Run(m, checker, Config{}, []string{"external"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s %s", f.File, f.Line, f.Target, f.Reason))
	}
	want := []string{
		"docs/a.md:2 https://vierxx.test external-status",
		"docs/a.md:3 https://langsam.test external-timeout",
		"docs/a.md:4 https://kreis.test external-redirects",
		"docs/a.md:5 https://dns.test external-status",
		"docs/b.md:1 https://vierxx.test external-status",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
	// Dedupe: zwei Vorkommen, genau EINE Prüfung der URL
	if checker.calls["https://vierxx.test"] != 1 {
		t.Fatalf("Aufrufe für vierxx = %d, want 1 (Dedupe)", checker.calls["https://vierxx.test"])
	}
}

// Review R1 zu slice-008: Fragmente werden vor Prüfung und Dedupe
// entfernt (eine Prüfung pro Ressource), der Befund nennt das
// Original-Linkziel; der Schema-Vergleich ist case-insensitiv.
func TestExternalFragmenteUndSchemaCase(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "[a](https://seite.test/doc#kapitel-1)\n" +
			"[b](https://seite.test/doc#kapitel-2)\n" +
			"[gross](HTTPS://gross.test)\n",
	})
	checker := newFakeChecker(map[string]driven.HTTPResult{
		"https://seite.test/doc": {Status: 404},
		"HTTPS://gross.test":     {Status: 404},
	})
	res, err := Run(m, checker, Config{}, []string{"external"})
	if err != nil {
		t.Fatal(err)
	}
	// genau EIN Request für beide Fragment-Varianten
	if checker.calls["https://seite.test/doc"] != 1 {
		t.Fatalf("Aufrufe = %v, want 1 für die fragmentfreie URL", checker.calls)
	}
	// Befunde an beiden Vorkommen, Target = Original-Linkziel
	var targets []string
	for _, f := range res.Findings {
		targets = append(targets, f.Target)
	}
	want := []string{ // sortiert nach Zeile (DC-QA-02.a)
		"https://seite.test/doc#kapitel-1",
		"https://seite.test/doc#kapitel-2",
		"HTTPS://gross.test", // case-insensitives Schema wird geprüft
	}
	if fmt.Sprint(targets) != fmt.Sprint(want) {
		t.Fatalf("Targets = %v, want %v", targets, want)
	}
}

// DC-FA-EXT-001 Boundary: ohne aktiviertes Modul erfolgt kein einziger
// Netzwerkzugriff — auch wenn externe Links existieren (DC-QA-03).
func TestExternalOptIn(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/a.md": "[extern](https://example.org)\n[kaputt](fehlt.md)\n",
	})
	res, err := Run(m, panicChecker{t}, Config{}, []string{"links", "anchors", "ids", "matrix"})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.Findings {
		if f.Rule == "external" {
			t.Fatalf("external-Befund ohne aktiviertes Modul: %+v", f)
		}
	}
}

// Konfigurations-Defaults: 0 = EXTERNAL_TIMEOUT/EXTERNAL_PARALLEL,
// gesetzte Werte gewinnen; checkURLs klemmt parallel < 1 auf 1.
func TestExternalConfigDefaultsUndParallelKlammer(t *testing.T) {
	var e ExternalConfig
	if e.EffectiveTimeoutSeconds() != 10 || e.EffectiveParallel() != 4 {
		t.Fatalf("Defaults = %d/%d, want 10/4",
			e.EffectiveTimeoutSeconds(), e.EffectiveParallel())
	}
	e = ExternalConfig{TimeoutSeconds: 30, Parallel: 2}
	if e.EffectiveTimeoutSeconds() != 30 || e.EffectiveParallel() != 2 {
		t.Fatalf("gesetzte Werte = %d/%d", e.EffectiveTimeoutSeconds(), e.EffectiveParallel())
	}
	checker := newFakeChecker(map[string]driven.HTTPResult{"https://x.test": {Status: 200}})
	if res := checkURLs(checker, []string{"https://x.test"}, 0); len(res) != 1 {
		t.Fatalf("parallel=0 muss auf 1 geklemmt werden: %v", res)
	}
}

// DC-QA-02: identische Eingabe ⇒ identische Ausgabe trotz interner
// Parallelität (Sortierung gilt auch für external).
func TestExternalDeterminismus(t *testing.T) {
	files := map[string]string{}
	results := map[string]driven.HTTPResult{}
	for i := 0; i < 20; i++ {
		url := fmt.Sprintf("https://h%02d.test", i)
		files[fmt.Sprintf("docs/f%02d.md", i)] = "[x](" + url + ")"
		results[url] = driven.HTTPResult{Status: 404}
	}
	m := coretest.NewMemFS(files)
	var prev []Finding
	for i := 0; i < 10; i++ {
		res, err := Run(m, newFakeChecker(results), Config{External: ExternalConfig{Parallel: 4}}, []string{"external"})
		if err != nil {
			t.Fatal(err)
		}
		if prev != nil && fmt.Sprint(prev) != fmt.Sprint(res.Findings) {
			t.Fatalf("Lauf %d weicht ab", i)
		}
		prev = res.Findings
	}
	if len(prev) != 20 {
		t.Fatalf("Befunde = %d, want 20", len(prev))
	}
}
