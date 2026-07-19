package rules

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls sources (spec/spezifikation.md §4).
const (
	ReasonSourceDrift       = "source-drift"
	ReasonSourceUnreachable = "source-unreachable"
)

// sourcesConfigFile ist die Befund-Datei für Config-Pins (die
// .d-check.yml führt sie; Marker-Pins tragen ihre eigene Datei).
const sourcesConfigFile = ".d-check.yml"

// Entpack-/Zip-Bomben-Grenzen des Content-Manifests (DC-FA-SRC-001.a Schritt 4).
const (
	maxUnpackBytes int64 = 256 << 20 // ≤ 256 MiB entpackte Gesamtgröße
	maxZipEntries        = 10000     // ≤ 10 000 Datei-Einträge
)

// sourcePinRE erkennt eine `source-pin`-Direktive und liefert ihren Körper
// (Gruppe 1) — der loose Detektor: ein Körper ohne gültiges `sha256:<hex>`
// ist eine malformte Direktive (fail-closed, DC-FA-SRC-001.a Schritt 2).
var sourcePinRE = regexp.MustCompile(`<!--\s*source-pin\s*:\s*(.*?)\s*-->`)

// sourcePinBodyRE parst den wohlgeformten Direktiven-Körper: optionales
// Schlüsselwort `zip` (Archiv) + `sha256:<hex>` (DC-FA-SRC-001.a Schritt 1).
var sourcePinBodyRE = regexp.MustCompile(`^(?:(zip)\s+)?sha256:([0-9a-fA-F]+)$`)

// SourceRef ist ein aufgelöster Quell-Pin (Marker oder Config): file/line für
// den Befund, url das http(s)-Ziel, sha256 der (kleingeschriebene) Pin-Hash,
// zip markiert ein Archiv-Ziel (unpack: zip).
type SourceRef struct {
	file   string
	line   int
	url    string
	sha256 string
	zip    bool
}

// CollectSourcePins sammelt die Marker-Pins einer Datei — nur bei aktivem
// Modul sources (der Aufrufer prüft applies). Ein wohlgeformter Marker bindet
// an den unmittelbar links stehenden http(s)-Link derselben Zeile; an einem
// repo-internen/Nicht-http(s)-Link oder ohne vorausgehenden Link ist er inert.
// Eine malformte Direktive (kein `sha256:<hex>`) ⇒ error (Aufrufer mappt auf
// Exit 2, fail-closed — DC-FA-SRC-001.a Schritt 2).
func CollectSourcePins(file string, lines []Line, content []byte) ([]SourceRef, error) {
	rawLines := strings.Split(string(content), "\n")
	var refs []SourceRef
	for _, ln := range lines {
		raw := ""
		if i := ln.No - 1; i >= 0 && i < len(rawLines) {
			raw = rawLines[i]
		}
		lineRefs, err := sourcePinsOnLine(file, ln, raw)
		if err != nil {
			return nil, err
		}
		refs = append(refs, lineRefs...)
	}
	return refs, nil
}

// sourcePinsOnLine parst die source-pin-Marker einer Zeile und bindet die
// wohlgeformten an ihren nächstgelegenen vorausgehenden http(s)-Link (auf der
// rohen Zeile geprüft — die Vorverarbeitung ist längenerhaltend, wie bei pins).
func sourcePinsOnLine(file string, ln Line, raw string) ([]SourceRef, error) {
	markers := sourcePinRE.FindAllStringSubmatchIndex(ln.Text, -1)
	if markers == nil {
		return nil, nil
	}
	links := nonImageLinkEnds(ln.Text)
	var refs []SourceRef
	for _, m := range markers {
		zipKind, hash, ok := parseSourcePinBody(ln.Text[m[2]:m[3]])
		if !ok {
			return nil, fmt.Errorf("%s:%d: malformte source-pin-Direktive (kein sha256:<hex>): %q",
				file, ln.No, strings.TrimSpace(ln.Text[m[0]:m[1]]))
		}
		target, bound := nearestLink(links, m[0], raw)
		if !bound || !hasHTTPScheme(target) {
			continue // inert: kein/Nicht-http(s)-Link (repo-intern bleibt pins/links-Domäne)
		}
		refs = append(refs, SourceRef{
			file: file, line: ln.No, url: stripFragment(target),
			sha256: strings.ToLower(hash), zip: zipKind,
		})
	}
	return refs, nil
}

// parseSourcePinBody prüft den wohlgeformten Körper (optional `zip` + gültiger
// sha256:<hex>); ok=false ⇒ malformt.
func parseSourcePinBody(body string) (zipKind bool, hash string, ok bool) {
	m := sourcePinBodyRE.FindStringSubmatch(strings.TrimSpace(body))
	if m == nil {
		return false, "", false
	}
	return m[1] != "", m[2], true
}

// linkEnd ist das schließende `)` und Ziel eines Nicht-Bild-Links.
type linkEnd struct {
	end    int
	target string
}

// nonImageLinkEnds liefert die Nicht-Bild-Links einer vorverarbeiteten Zeile.
func nonImageLinkEnds(text string) []linkEnd {
	var out []linkEnd
	forEachLink(text, func(ref LinkRef, span LinkSpan) {
		if !span.IsImage {
			out = append(out, linkEnd{end: span.End, target: ref.Target})
		}
	})
	return out
}

// nearestLink bindet an den Link, dessen schließendes `)` dem Marker
// unmittelbar (nur Whitespace dazwischen, auf der rohen Zeile) vorausgeht; bei
// mehreren den nächstgelegenen (größtes end). Ohne Treffer inert.
func nearestLink(links []linkEnd, markerStart int, raw string) (string, bool) {
	best := -1
	for i, l := range links {
		if l.end <= markerStart && markerStart <= len(raw) && strings.TrimSpace(raw[l.end:markerStart]) == "" {
			if best == -1 || l.end > links[best].end {
				best = i
			}
		}
	}
	if best == -1 {
		return "", false
	}
	return links[best].target, true
}

// stripFragment entfernt ein #fragment (wie external — Fragmente sind für den
// Inhalts-Abruf bedeutungslos).
func stripFragment(u string) string {
	if i := strings.IndexByte(u, '#'); i != -1 {
		return u[:i]
	}
	return u
}

// CheckSources ist der Netz-Post-Pass des Moduls sources (DC-FA-SRC-001.a):
// über die gesammelten Marker-Pins und die Config-Pins wird je gepinnte Quelle
// geholt, gehasht und gegen den Pin verglichen. Gleich ⇒ kein Befund; ungleich
// ⇒ source-drift (voller Ist-Hash); nicht materialisierbar ⇒ source-unreachable.
// checker == nil (Modul nicht verdrahtet) ist No-op (DC-QA-03). Die Findung ist
// deterministisch aus der Pin-Liste (DC-QA-02); je (url, unpack) wird genau
// einmal geholt.
func CheckSources(checker driven.HTTPChecker, markerRefs []SourceRef, cfg model.SourcesConfig) []model.Finding {
	if checker == nil {
		return nil
	}
	refs := allSourceRefs(markerRefs, cfg)
	if len(refs) == 0 {
		return nil
	}
	results := fetchSources(checker, refs)
	var findings []model.Finding
	for _, r := range refs {
		if f, ok := sourceVerdict(r, results[fetchKey{r.url, r.zip}]); ok {
			findings = append(findings, f)
		}
	}
	return findings
}

// allSourceRefs vereint Marker-Pins und Config-Pins (Config-Befunde tragen
// .d-check.yml + die url-Feld-Zeile).
func allSourceRefs(markerRefs []SourceRef, cfg model.SourcesConfig) []SourceRef {
	refs := make([]SourceRef, 0, len(markerRefs)+len(cfg.Pins))
	refs = append(refs, markerRefs...)
	for _, p := range cfg.Pins {
		refs = append(refs, SourceRef{
			file: sourcesConfigFile, line: p.Line, url: p.URL,
			sha256: strings.ToLower(p.Sha256), zip: p.Unpack == model.SourceUnpackZip,
		})
	}
	return refs
}

type fetchKey struct {
	url string
	zip bool
}

// sourceResult trägt entweder den errechneten Hash (reason == "") oder einen
// source-unreachable-Grund samt Meldung.
type sourceResult struct {
	hash    string
	reason  string
	message string
}

// fetchSources holt jede eindeutige (url, unpack)-Kombination genau einmal.
func fetchSources(checker driven.HTTPChecker, refs []SourceRef) map[fetchKey]sourceResult {
	out := map[fetchKey]sourceResult{}
	for _, r := range refs {
		k := fetchKey{r.url, r.zip}
		if _, done := out[k]; !done {
			out[k] = fetchAndHash(checker, r.url, r.zip)
		}
	}
	return out
}

// fetchAndHash holt eine Quelle und errechnet ihren Content-Hash (Roh-Bytes
// bzw. Zip-Content-Manifest). Jede nicht materialisierbare Antwort ⇒
// source-unreachable, ohne Hash-Schritt.
func fetchAndHash(checker driven.HTTPChecker, url string, zipKind bool) sourceResult {
	res := checker.Fetch(url)
	if reason, msg, bad := fetchFailure(res); bad {
		return sourceResult{reason: reason, message: msg}
	}
	if zipKind {
		h, ok := zipManifestHash(res.Body)
		if !ok {
			return sourceResult{
				reason:  ReasonSourceUnreachable,
				message: "kein gültiges Zip oder Entpack-Limit überschritten (≤ 256 MiB, ≤ 10000 Einträge)",
			}
		}
		return sourceResult{hash: h}
	}
	sum := sha256.Sum256(res.Body)
	return sourceResult{hash: hex.EncodeToString(sum[:])}
}

// fetchFailure mappt eine Fehler-Antwort auf source-unreachable (bad=true).
func fetchFailure(res driven.FetchResult) (reason, message string, bad bool) {
	switch {
	case res.Timeout:
		return ReasonSourceUnreachable, "Timeout überschritten", true
	case res.TooManyRedirects:
		return ReasonSourceUnreachable, "Redirect-Kette länger als 5 Stationen", true
	case res.TransportError != "":
		return ReasonSourceUnreachable, "nicht erreichbar: " + res.TransportError, true
	case res.TooLarge:
		return ReasonSourceUnreachable, "Body-Limit überschritten (> 64 MiB)", true
	case res.Status >= 400:
		return ReasonSourceUnreachable, fmt.Sprintf("HTTP-Status %d", res.Status), true
	}
	return "", "", false
}

// sourceVerdict vergleicht Ist-Hash und Pin (beide kleingeschrieben) und
// liefert den Befund (ok=true) bzw. keinen (Hash gleich).
func sourceVerdict(r SourceRef, res sourceResult) (model.Finding, bool) {
	switch {
	case res.reason != "":
		return newSourceFinding(r, res.reason, res.message), true
	case res.hash == r.sha256:
		return model.Finding{}, false
	default:
		return newSourceFinding(r, ReasonSourceDrift,
			"Upstream-Inhalt gedriftet: errechnet sha256:"+res.hash+" (Pin sha256:"+r.sha256+")"), true
	}
}

func newSourceFinding(r SourceRef, reason, message string) model.Finding {
	return model.Finding{
		File: r.file, Line: r.line, Rule: "sources",
		Target: r.url, Reason: reason, Message: message,
	}
}

// zipEntry ist ein Content-Manifest-Eintrag: normalisierter Pfad + Inhalts-Hash.
type zipEntry struct {
	path   string
	hexsum string
}

// zipManifestHash errechnet den byte-genauen Content-Manifest-Hash eines Zip
// (DC-FA-SRC-001.a Schritt 4): je regulärer Datei-Eintrag `<hex>  <pfad>\n`,
// nach `<pfad>` (sekundär `<hex>`) byteweise sortiert, davon der sha256. ok=false
// bei ungültigem Zip oder überschrittenem Entpack-Limit.
func zipManifestHash(body []byte) (string, bool) {
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", false
	}
	entries, ok := zipEntries(zr)
	if !ok {
		return "", false
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].path != entries[j].path {
			return entries[i].path < entries[j].path
		}
		return entries[i].hexsum < entries[j].hexsum
	})
	var buf bytes.Buffer
	for _, e := range entries {
		buf.WriteString(e.hexsum)
		buf.WriteString("  ")
		buf.WriteString(e.path)
		buf.WriteByte('\n')
	}
	sum := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(sum[:]), true
}

// zipEntries hasht die regulären Datei-Einträge unter den Ressourcen-Grenzen
// (Verzeichnis-Einträge — Name endet auf `/` — ausgenommen).
func zipEntries(zr *zip.Reader) ([]zipEntry, bool) {
	var entries []zipEntry
	remaining := maxUnpackBytes
	count := 0
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}
		count++
		if count > maxZipEntries {
			return nil, false
		}
		sum, n, ok := hashZipFile(f, remaining)
		if !ok {
			return nil, false
		}
		remaining -= n
		entries = append(entries, zipEntry{path: normalizeZipPath(f.Name), hexsum: sum})
	}
	return entries, true
}

// hashZipFile liest einen Eintrag bis zum verbleibenden Budget und hasht ihn;
// ok=false bei Lese-Fehler oder Budget-Überschreitung.
func hashZipFile(f *zip.File, budget int64) (sum string, n int64, ok bool) {
	rc, err := f.Open()
	if err != nil {
		return "", 0, false
	}
	defer func() { _ = rc.Close() }()
	h := sha256.New()
	n, err = io.Copy(h, io.LimitReader(rc, budget+1))
	if err != nil {
		return "", 0, false
	}
	if n > budget {
		return "", 0, false
	}
	return hex.EncodeToString(h.Sum(nil)), n, true
}

// normalizeZipPath normalisiert einen Zip-internen Pfad: Backslash → `/`,
// führendes `./` und `/` entfernt; der volle Pfad bleibt (kein Basisname).
func normalizeZipPath(name string) string {
	p := strings.ReplaceAll(name, "\\", "/")
	for {
		switch {
		case strings.HasPrefix(p, "./"):
			p = p[2:]
		case strings.HasPrefix(p, "/"):
			p = p[1:]
		default:
			return p
		}
	}
}
