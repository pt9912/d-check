package rules

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/d-check/internal/adapter/driven/httpcheck"
	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// fileOKBody ist der stabile Inhalt der Einzeldatei-Quelle /file-ok.
const fileOKBody = "hallo welt\nzeile zwei\n"

// zipFiles sind die kanonischen Archiv-Einträge (normalisierte Pfade →
// Inhalt); enthält einen verschachtelten Pfad, damit der VOLLE Pfad (kein
// Basisname) ins Manifest muss.
func zipFiles() map[string]string {
	return map[string]string{
		"a.txt":                "AAA\n",
		"b.txt":                "BBB\n",
		"unterordner/datei.md": "verschachtelt\n",
	}
}

// rawHash liefert den Kleinbuchstaben-Hex-sha256 der Bytes.
func rawHash(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// wantManifest errechnet den erwarteten Content-Manifest-Hash UNABHÄNGIG von
// der Produktions-Funktion (echtes Orakel): je Datei `<hex>  <pfad>\n`, nach
// Pfad (sekundär hex) sortiert, davon der sha256. Erzwingt damit die
// Produktions-Sortierung — ein unsortiertes Manifest kippt den Vergleich.
func wantManifest(files map[string]string) string {
	type e struct{ path, hexsum string }
	var es []e
	for name, content := range files {
		es = append(es, e{path: name, hexsum: rawHash([]byte(content))})
	}
	sort.Slice(es, func(i, j int) bool {
		if es[i].path != es[j].path {
			return es[i].path < es[j].path
		}
		return es[i].hexsum < es[j].hexsum
	})
	var buf bytes.Buffer
	for _, x := range es {
		buf.WriteString(x.hexsum)
		buf.WriteString("  ")
		buf.WriteString(x.path)
		buf.WriteByte('\n')
	}
	return rawHash(buf.Bytes())
}

// buildZip erzeugt ein Zip mit den Einträgen in der gegebenen Reihenfolge.
func buildZip(t *testing.T, order []string, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range order {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(files[name])); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// sourcesServer bedient die von den Tests referenzierten Routen (loopback).
func sourcesServer(t *testing.T, zipA, zipB []byte) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	write := func(w http.ResponseWriter, body []byte) {
		_, _ = w.Write(body)
	}
	mux.HandleFunc("/file-ok", func(w http.ResponseWriter, _ *http.Request) { write(w, []byte(fileOKBody)) })
	mux.HandleFunc("/file-changed", func(w http.ResponseWriter, _ *http.Request) { write(w, []byte("anderer inhalt\n")) })
	mux.HandleFunc("/notfound", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/html", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		write(w, []byte("<html><body>kein zip</body></html>"))
	})
	mux.HandleFunc("/zip-a", func(w http.ResponseWriter, _ *http.Request) { write(w, zipA) })
	mux.HandleFunc("/zip-b", func(w http.ResponseWriter, _ *http.Request) { write(w, zipB) })
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func realChecker() driven.HTTPChecker { return httpcheck.New(3 * time.Second) }

// sourceFindings filtert die sources-Befunde eines Laufs.
func sourceFindings(fs []model.Finding) []model.Finding {
	var out []model.Finding
	for _, f := range fs {
		if f.Rule == "sources" {
			out = append(out, f)
		}
	}
	return out
}

// DC-FA-SRC-001 Happy: eine per Marker gepinnte Einzeldatei UND ein per Config
// gepinntes Archiv (unpack: zip), beide unverändert ⇒ 0 Befunde.
func TestSourcesHappyMarkerUndConfigZip(t *testing.T) {
	files := zipFiles()
	srv := sourcesServer(t, buildZip(t, []string{"a.txt", "b.txt", "unterordner/datei.md"}, files), nil)
	m := coretest.NewMemFS(map[string]string{
		"docs/quelle.md": "[Regelwerk](" + srv.URL + "/file-ok) <!-- source-pin: sha256:" + rawHash([]byte(fileOKBody)) + " -->\n",
	})
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: srv.URL + "/zip-a", Sha256: wantManifest(files), Unpack: "zip", Line: 5},
	}}}
	res, err := Run(m, realChecker(), cfg, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("erwartet 0 Befunde, got %+v", got)
	}
}

// DC-FA-SRC-001 Boundary (Archiv-Determinismus, MUTATIONS-verifiziert): zwei
// Zips MIT UNTERSCHIEDLICHER Eintrags-Reihenfolge, inhaltlich identisch, beide
// auf denselben Manifest-Hash gepinnt ⇒ KEIN source-drift. Beide Zips sind
// bewusst NICHT pfad-sortiert eingelagert — ein Manifest ohne die
// Pfad-Sortierung ergäbe für mindestens eines einen anderen Hash und kippte den
// Test (erzwingt die Sortier-Invariante).
func TestSourcesZipReorderInvariant(t *testing.T) {
	files := zipFiles()
	zipA := buildZip(t, []string{"unterordner/datei.md", "b.txt", "a.txt"}, files) // rückwärts
	zipB := buildZip(t, []string{"b.txt", "unterordner/datei.md", "a.txt"}, files) // gemischt
	srv := sourcesServer(t, zipA, zipB)
	pin := wantManifest(files)
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: srv.URL + "/zip-a", Sha256: pin, Unpack: "zip", Line: 3},
		{URL: srv.URL + "/zip-b", Sha256: pin, Unpack: "zip", Line: 6},
	}}}
	res, err := Run(coretest.NewMemFS(nil), realChecker(), cfg, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("Reihenfolge-Invarianz verletzt (Sortierung fehlt?): %+v", got)
	}
}

// DC-FA-SRC-001 Boundary (Modul-aus / netzlos): ohne aktives sources bleibt der
// Befundsatz byte-identisch und es wird KEINE Netzverbindung geöffnet
// (panicChecker.Fetch würde sofort fatal).
func TestSourcesModuleOffKeinNetz(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/quelle.md": "[x](https://irrelevant.test/file-ok) <!-- source-pin: sha256:" +
			strings.Repeat("0", 64) + " -->\n",
	})
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: "https://irrelevant.test/x", Sha256: strings.Repeat("0", 64), Unpack: "none", Line: 1},
	}}}
	res, err := Run(m, panicChecker{t}, cfg, []string{"links", "anchors"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("sources-Befund ohne aktives Modul: %+v", got)
	}
}

// DC-FA-SRC-001 Negative (Drift): geänderter Inhalt ⇒ genau 1 source-drift mit
// dem VOLLEN Kleinbuchstaben-Ist-Hash in der Meldung (Re-Pin-Vorlage).
func TestSourcesDriftVollerHash(t *testing.T) {
	srv := sourcesServer(t, nil, nil)
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: srv.URL + "/file-ok", Sha256: strings.Repeat("0", 64), Unpack: "none", Line: 2},
	}}}
	res, err := Run(coretest.NewMemFS(nil), realChecker(), cfg, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	got := sourceFindings(res.Findings)
	if len(got) != 1 || got[0].Reason != ReasonSourceDrift {
		t.Fatalf("erwartet genau 1 source-drift, got %+v", got)
	}
	if !strings.Contains(got[0].Message, rawHash([]byte(fileOKBody))) {
		t.Fatalf("Meldung ohne vollen Ist-Hash: %q", got[0].Message)
	}
	if got[0].File != ".d-check.yml" || got[0].Line != 2 || got[0].Target != srv.URL+"/file-ok" {
		t.Fatalf("Befund-Koordinaten falsch: %+v", got[0])
	}
}

// DC-FA-SRC-001 Boundary (unerreichbar ≠ Drift): 404 sowie unpack: zip auf eine
// 200-HTML-Antwort ⇒ source-unreachable (nicht source-drift).
func TestSourcesUnreachableNichtDrift(t *testing.T) {
	srv := sourcesServer(t, nil, nil)
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: srv.URL + "/notfound", Sha256: strings.Repeat("a", 64), Unpack: "none", Line: 2},
		{URL: srv.URL + "/html", Sha256: strings.Repeat("a", 64), Unpack: "zip", Line: 4},
	}}}
	res, err := Run(coretest.NewMemFS(nil), realChecker(), cfg, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	got := sourceFindings(res.Findings)
	if len(got) != 2 {
		t.Fatalf("erwartet 2 Befunde, got %+v", got)
	}
	for _, f := range got {
		if f.Reason != ReasonSourceUnreachable {
			t.Fatalf("erwartet source-unreachable, got %+v", f)
		}
	}
}

// DC-FA-SRC-001 Boundary (kein Doppelbefund / repo-intern): ein source-pin an
// einem repo-internen Link ist inert — kein Befund, KEIN Netzzugriff
// (panicChecker.Fetch würde fatal).
func TestSourcesRepoInternMarkerInert(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/quelle.md": "[intern](./andere.md) <!-- source-pin: sha256:" + strings.Repeat("0", 64) + " -->\n",
	})
	res, err := Run(m, panicChecker{t}, model.Config{}, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("repo-interner Marker nicht inert: %+v", got)
	}
}

// DC-FA-SRC-001 Boundary (Marker-Bindung): bei mehreren Links bindet der Marker
// an den unmittelbar links stehenden (hier /file-changed, dessen Hash gepinnt
// ist ⇒ grün); bände er an /file-ok, kippte es zu source-drift. Ein Marker ohne
// vorausgehenden Link auf derselben Zeile ist inert.
func TestSourcesMarkerBindungNaechsterLink(t *testing.T) {
	srv := sourcesServer(t, nil, nil)
	changed := rawHash([]byte("anderer inhalt\n"))
	m := coretest.NewMemFS(map[string]string{
		"docs/quelle.md": "[a](" + srv.URL + "/file-ok) [b](" + srv.URL + "/file-changed) <!-- source-pin: sha256:" + changed + " -->\n" +
			"kein link hier <!-- source-pin: sha256:" + strings.Repeat("f", 64) + " -->\n",
	})
	res, err := Run(m, realChecker(), model.Config{}, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("Marker band nicht an den nächsten Link (oder freier Marker nicht inert): %+v", got)
	}
}

// DC-FA-SRC-001: Case-Insensitivität — großgeschriebener Pin bei korrektem
// Inhalt ⇒ grün.
func TestSourcesCaseInsensitiv(t *testing.T) {
	srv := sourcesServer(t, nil, nil)
	cfg := model.Config{Sources: model.SourcesConfig{Pins: []model.SourcePin{
		{URL: srv.URL + "/file-ok", Sha256: strings.ToUpper(rawHash([]byte(fileOKBody))), Unpack: "none", Line: 2},
	}}}
	res, err := Run(coretest.NewMemFS(nil), realChecker(), cfg, []string{"sources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := sourceFindings(res.Findings); len(got) != 0 {
		t.Fatalf("Groß-Pin bei korrektem Inhalt nicht grün: %+v", got)
	}
}

// DC-FA-SRC-001 fail-closed: eine malformte source-pin-Direktive (kein
// sha256:<hex>) an einem http(s)-Link ⇒ Exit-2-Pfad (Run gibt error).
func TestSourcesMalformierterMarkerFailClosed(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"docs/quelle.md": "[x](https://example.test/f) <!-- source-pin: kein-hash -->\n",
	})
	if _, err := Run(m, panicChecker{t}, model.Config{}, []string{"sources"}); err == nil {
		t.Fatal("malformte source-pin-Direktive: error (Exit 2) erwartet")
	}
}

// DC-FA-SRC-001 Boundary (Fetch-Fehlerklassen, deterministisch über einen Stub):
// Timeout, > 5 Redirects, Transportfehler und überschrittenes Body-Limit ⇒ je
// source-unreachable (getrennt von source-drift).
func TestSourcesFetchFehlerklassen(t *testing.T) {
	stub := stubFetcher{fetch: map[string]driven.FetchResult{
		"http://x/timeout":   {Timeout: true},
		"http://x/redirects": {TooManyRedirects: true},
		"http://x/transport": {TransportError: "no such host"},
		"http://x/toolarge":  {TooLarge: true},
	}}
	var refs []SourceRef
	for _, u := range []string{"http://x/timeout", "http://x/redirects", "http://x/transport", "http://x/toolarge"} {
		refs = append(refs, SourceRef{file: "docs/a.md", line: 1, url: u, sha256: strings.Repeat("0", 64)})
	}
	got := CheckSources(stub, refs, model.SourcesConfig{})
	if len(got) != 4 {
		t.Fatalf("erwartet 4 Befunde, got %+v", got)
	}
	for _, f := range got {
		if f.Reason != ReasonSourceUnreachable {
			t.Fatalf("erwartet source-unreachable, got %+v", f)
		}
	}
}

// stubFetcher liefert kanonisierte Fetch-/Check-Ergebnisse aus Tabellen (für die
// Fehlerklassen-Tests ohne echtes Netz).
type stubFetcher struct {
	fetch map[string]driven.FetchResult
	check map[string]driven.HTTPResult
}

func (s stubFetcher) Check(url string) driven.HTTPResult { return s.check[url] }

func (s stubFetcher) Fetch(url string) driven.FetchResult { return s.fetch[url] }

// CheckSources ist No-op ohne Checker (DC-QA-03: nil = keine Netz-Tür).
func TestSourcesNilCheckerNoop(t *testing.T) {
	got := CheckSources(nil, []SourceRef{{url: "http://x/a", sha256: "x"}}, model.SourcesConfig{})
	if got != nil {
		t.Fatalf("nil-Checker muss No-op sein, got %+v", got)
	}
}
