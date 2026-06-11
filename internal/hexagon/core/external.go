package core

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls external (spec/spezifikation.md §4).
const (
	ReasonExternalStatus    = "external-status"
	ReasonExternalTimeout   = "external-timeout"
	ReasonExternalRedirects = "external-redirects"
)

// externalRef ist ein Vorkommen einer externen URL: target ist das
// Original-Linkziel (Befund-Target), url die Prüf-URL ohne Fragment
// (Request- und Dedupe-Key — Fragmente werden nie übertragen).
type externalRef struct {
	file   string
	line   int
	target string
	url    string
}

// collectExternalURLs sammelt die http(s)-Linkziele einer Datei —
// nur bei explizit aktiviertem Modul external (DC-FA-EXT-001,
// opt-in). Der Schema-Vergleich ist case-insensitiv (RFC 3986);
// konsistent mit IsExternalScheme, das solche Ziele für `links`
// überspringt.
func collectExternalURLs(file string, lines []Line) []externalRef {
	var refs []externalRef
	for _, ref := range ExtractLinks(lines) {
		t := ref.Target
		if !hasHTTPScheme(t) {
			continue
		}
		url := t
		if idx := strings.IndexByte(url, '#'); idx != -1 {
			url = url[:idx]
		}
		refs = append(refs, externalRef{file: file, line: ref.Line, target: t, url: url})
	}
	return refs
}

// hasHTTPScheme prüft case-insensitiv auf http://- bzw.
// https://-Schema (DC-FA-EXT-001: nur diese Schemata werden geprüft).
func hasHTTPScheme(t string) bool {
	return (len(t) >= 7 && strings.EqualFold(t[:7], "http://")) ||
		(len(t) >= 8 && strings.EqualFold(t[:8], "https://"))
}

// checkExternal prüft die gesammelten URLs — genau eine Prüfung pro
// URL (Dedupe), begrenzte Parallelität — und mappt die Ergebnisse auf
// Befunde pro Vorkommen (spec/spezifikation.md §DC-FA-EXT-001.a).
// Die Parallelität darf die Ausgabe nicht beeinflussen (DC-QA-02.a) —
// Befunde entstehen aus der deterministischen Vorkommens-Liste.
func checkExternal(checker driven.HTTPChecker, refs []externalRef, parallel int) []Finding {
	if checker == nil || len(refs) == 0 {
		return nil
	}
	results := checkURLs(checker, uniqueURLs(refs), parallel)
	var findings []Finding
	for _, r := range refs {
		reason, msg := externalVerdict(results[r.url])
		if reason == "" {
			continue
		}
		findings = append(findings, Finding{
			File: r.file, Line: r.line, Rule: "external",
			Target: r.target, Reason: reason, Message: msg,
		})
	}
	return findings
}

// uniqueURLs liefert die URLs der Vorkommen dedupliziert und sortiert.
func uniqueURLs(refs []externalRef) []string {
	seen := map[string]bool{}
	var urls []string
	for _, r := range refs {
		if !seen[r.url] {
			seen[r.url] = true
			urls = append(urls, r.url)
		}
	}
	sort.Strings(urls)
	return urls
}

// checkURLs führt die Prüfungen mit begrenzter Parallelität aus
// (EXTERNAL_PARALLEL, spec/spezifikation.md §3).
func checkURLs(checker driven.HTTPChecker, urls []string, parallel int) map[string]driven.HTTPResult {
	if parallel < 1 {
		parallel = 1
	}
	results := make(map[string]driven.HTTPResult, len(urls))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)
	for _, u := range urls {
		wg.Add(1)
		sem <- struct{}{}
		go func(u string) {
			defer wg.Done()
			defer func() { <-sem }()
			res := checker.Check(u)
			mu.Lock()
			results[u] = res
			mu.Unlock()
		}(u)
	}
	wg.Wait()
	return results
}

// externalVerdict mappt ein Port-Ergebnis auf Grund-Code und Meldung:
// Status < 400 → kein Befund; Transportfehler zählen als
// `external-status` (spec/spezifikation.md §4).
func externalVerdict(res driven.HTTPResult) (reason, message string) {
	switch {
	case res.Timeout:
		return ReasonExternalTimeout, "Timeout überschritten"
	case res.TooManyRedirects:
		return ReasonExternalRedirects, "Redirect-Kette länger als 5 Stationen"
	case res.TransportError != "":
		return ReasonExternalStatus, "nicht erreichbar: " + res.TransportError
	case res.Status >= 400:
		return ReasonExternalStatus, fmt.Sprintf("HTTP-Status %d", res.Status)
	}
	return "", ""
}
