package rules

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// versionsFixture: ein Release-Register (version.md, §Aktuell nennt v0.27.0)
// und ein Doc mit einem ghcr-Pin in einem bash-Fence.
func versionsFixture() map[string]string {
	return map[string]string{
		"version.md": "# Register\n\n## Aktuell\n\n" +
			"Aktuelle Version: [`v0.27.0`](#v0.27.0) — 2026-06-23.\n\n" +
			"## Verlauf\n\n| `v0.27.0` <a id=\"v0.27.0\"></a> | x |\n| `v0.26.0` | y |\n",
		"docs/use.md": "Beispiel:\n\n```bash\ndocker run ghcr.io/pt9912/d-check:v0.27.0\n```\n",
	}
}

func versionsCfg() model.Config {
	return model.Config{Versions: model.VersionsConfig{Patterns: []model.VersionPattern{{
		PinPattern:  regexp.MustCompile(`ghcr\.io/[^\s:]+:(v[0-9]+\.[0-9]+\.[0-9]+)`),
		CurrentFrom: "version.md#aktuell",
	}}}}
}

// DC-FA-VER-001 Happy: alle Pins tragen die aktuelle Version → kein Befund.
func TestVersionsHappy(t *testing.T) {
	m := coretest.NewMemFS(versionsFixture())
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("alle Pins == aktuell → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 Negative: ein veralteter Pin INNERHALB eines Fenced-Blocks →
// version-stale (Zeile im Fence).
func TestVersionsStaleInFence(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "```bash\ndocker run ghcr.io/pt9912/d-check:v0.26.0\n```\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%d %s %s %s", f.Line, f.Rule, f.Target, f.Reason))
	}
	want := []string{"2 versions v0.26.0 version-stale"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Befunde = %v\nwant %v", got, want)
	}
}

// DC-FA-VER-001 Boundary: ein Pin in einer exempt-paths-Datei und ein zweiter
// auf einer d-check:ignore-Zeile → kein Befund für beide.
func TestVersionsVentile(t *testing.T) {
	fix := versionsFixture()
	fix["docs/alt.md"] = "ghcr.io/pt9912/d-check:v0.1.0\n"
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.1.0 <!-- d-check:ignore (Beispiel) -->\n"
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].ExemptPaths = []string{"docs/alt.md"}
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("exempt-paths + d-check:ignore → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 Modul-aus: ohne aktives Modul kein Befund (byte-identisch).
func TestVersionsNichtAktiv(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.26.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, versionsCfg(), []string{"links"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("Modul nicht aktiv → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001: die current-from-Datei ist selbst von der Pin-Prüfung
// ausgenommen (ihr Verlauf darf alte Pins listen).
func TestVersionsCurrentFromSelbstAusgenommen(t *testing.T) {
	fix := versionsFixture()
	fix["version.md"] += "\nghcr.io/pt9912/d-check:v0.1.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("current-from-Datei selbst-ausgenommen → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 fail-closed: ein nicht auflösbarer current-from-Anker → Exit 2.
func TestVersionsCurrentFromAnkerFehlt(t *testing.T) {
	m := coretest.NewMemFS(versionsFixture())
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].CurrentFrom = "version.md#fehlt"
	if _, err := Run(m, nil, cfg, []string{"versions"}); err == nil {
		t.Fatal("nicht auflösbarer current-from-Anker muss Fehler liefern")
	}
}

// DC-FA-VER-001 fail-closed: ein Span ohne erkennbare Version → Exit 2.
func TestVersionsCurrentFromKeineVersion(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md":  "## Aktuell\n\nkeine Versionsnummer hier\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	if _, err := Run(m, nil, versionsCfg(), []string{"versions"}); err == nil {
		t.Fatal("current-from-Span ohne Version muss Fehler liefern")
	}
}

// DC-FA-VER-001: current-from über einen Inline-HTML-Anker (kein Heading-Slug).
func TestVersionsCurrentFromHTMLAnker(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md":  "<a id=\"cur\"></a> aktuell ist v0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].CurrentFrom = "version.md#cur"
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("HTML-Anker-current-from → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001: current-from ohne Anker = ganze Datei.
func TestVersionsCurrentFromGanzeDatei(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"reg.md":      "v0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].CurrentFrom = "reg.md"
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("ganze-Datei current-from → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001: ohne Capture-Gruppe zählt der ganze Treffer als Version.
func TestVersionsPinOhneGruppe(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "tag v0.26.0 irgendwo\n"
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].PinPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "v0.26.0" {
		t.Fatalf("ohne Capture-Gruppe = ganzer Treffer, got %v", res.Findings)
	}
}

// DC-FA-VER-001 / DC-FA-ANCH-001.b: ein HTML-Anker INNERHALB eines Fence ist
// keiner — sonst loeste current-from auf, waehrend anchors denselben Anker im
// selben Lauf fuer nicht existent haelt. Erwartet ist der fail-closed-Abbruch.
func TestVersionsCurrentFromAnkerNurImFence(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md": "# Versionen\n\n```html\n<a id=\"aktuell\"></a> v0.27.0\n```\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	if _, err := Run(m, nil, versionsCfg(), []string{"versions"}); err == nil {
		t.Fatal("HTML-Anker nur im Fence: fail-closed (Exit 2) erwartet")
	}
}

// Gegenprobe: derselbe Anker AUSSERHALB des Fence loest auf und die Pruefung
// laeuft — die Ablehnung liegt am Fence, nicht am Anker.
func TestVersionsCurrentFromAnkerAusserhalbFence(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md": "# Versionen\n\n<a id=\"aktuell\"></a> v0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatalf("Anker ausserhalb des Fence muss aufloesen: %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("passender Pin → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 / DC-FA-ANCH-001.b: die Anker-Frage hat EINE Antwort. Vier
// Zeichenfolgen, die wie ein Anker aussehen und keiner sind — fuer anchors
// waren sie nie welche, fuer current-from bis zur Vereinheitlichung schon.
// Erwartet ist jedes Mal der fail-closed-Abbruch.
func TestVersionsCurrentFromAnkerParitaet(t *testing.T) {
	faelle := map[string]string{
		"Inline-Code":     "# V\n\nBeispiel: `<a id=\"aktuell\"></a>` v0.27.0\n",
		"data-id":         "# V\n\n<span data-id=\"aktuell\"></span> v0.27.0\n",
		"ohne Tag":        "# V\n\nid=\"aktuell\" v0.27.0\n",
		"name an area":    "# V\n\n<area name=\"aktuell\"> v0.27.0\n",
	}
	for name, version := range faelle {
		t.Run(name, func(t *testing.T) {
			m := coretest.NewMemFS(map[string]string{
				"version.md":  version,
				"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
			})
			if _, err := Run(m, nil, versionsCfg(), []string{"versions"}); err == nil {
				t.Fatalf("%s ist kein Anker → fail-closed (Exit 2) erwartet", name)
			}
		})
	}
}

// Gegenprobe zur Paritaet: das echte <a id=...> loest auf.
func TestVersionsCurrentFromEchterHTMLAnker(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md":  "# V\n\n<a id=\"aktuell\"></a> v0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatalf("echter HTML-Anker muss aufloesen: %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("passender Pin → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 / DC-FA-ANCH-001: der Duplikat-Slug ist Teil der Anker-Antwort.
// Ohne den Zaehler loest #alt-1 nicht auf, waehrend anchors ihn fuer gueltig
// haelt — dieselbe Datei, dieselbe Frage, zwei Antworten.
func TestVersionsCurrentFromDuplikatSlug(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md":  "# V\n\n## Alt\n\nv0.1.0\n\n## Alt\n\nv0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].CurrentFrom = "version.md#alt-1"
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatalf("Duplikat-Slug muss aufloesen: %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("zweiter Alt-Abschnitt traegt v0.27.0 → 0 Befunde, got %v", res.Findings)
	}
}

// Prozent-kodierte Fragmente adressieren den dekodierten Anker — anchors
// dekodiert seit jeher, current-from bis zur Vereinheitlichung nicht.
func TestVersionsCurrentFromKodiertesFragment(t *testing.T) {
	m := coretest.NewMemFS(map[string]string{
		"version.md":  "# V\n\n<a id=\"a b\"></a> v0.27.0\n",
		"docs/use.md": "ghcr.io/pt9912/d-check:v0.27.0\n",
	})
	cfg := versionsCfg()
	cfg.Versions.Patterns[0].CurrentFrom = "version.md#a%20b"
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatalf("kodiertes Fragment muss aufloesen: %v", err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("passender Pin → 0 Befunde, got %v", res.Findings)
	}
}

// zweiReihenFixture: zwei unabhängige Versions-Reihen — das eigene Release
// (version.md, v0.27.0) und ein fremder Baseline-Stand (pin.md, v5.9.0).
func zweiReihenFixture() map[string]string {
	fix := versionsFixture()
	fix["pin.md"] = "# Konventionen\n\n## Baseline\n\nGepinnt auf `v5.9.0`.\n"
	return fix
}

// zweiReihenCfg: Paar 1 prüft ghcr-Pins gegen version.md#aktuell, Paar 2
// baseline/<tag>-Pfade gegen pin.md#baseline.
func zweiReihenCfg() model.Config {
	return model.Config{Versions: model.VersionsConfig{Patterns: []model.VersionPattern{
		{
			PinPattern:  regexp.MustCompile(`ghcr\.io/[^\s:]+:(v[0-9]+\.[0-9]+\.[0-9]+)`),
			CurrentFrom: "version.md#aktuell",
		},
		{
			PinPattern:  regexp.MustCompile(`baseline/(v[0-9]+\.[0-9]+\.[0-9]+)`),
			CurrentFrom: "pin.md#baseline",
		},
	}}}
}

// DC-FA-VER-001 Happy Path (zwei Paare): beide Reihen aktuell → kein Befund.
func TestVersionsZweiPaareHappy(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.27.0 und baseline/v5.9.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, zweiReihenCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("beide Reihen aktuell → 0 Befunde, got %v", res.Findings)
	}
}

// DC-FA-VER-001 Negative (je Paar): nur die zweite Reihe ist veraltet → genau
// ein Befund, dessen erwartete Version aus der Quelle des ZWEITEN Paares
// stammt.
func TestVersionsZweitesPaarMeldet(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.27.0 und baseline/v5.7.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, zweiReihenCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "v5.7.0" {
		t.Fatalf("nur die zweite Reihe veraltet → 1 Befund v5.7.0, got %v", res.Findings)
	}
	if !strings.Contains(res.Findings[0].Message, "erwartet v5.9.0") {
		t.Fatalf("erwartete Version muss aus der Quelle des zweiten Paares stammen: %q", res.Findings[0].Message)
	}
}

// DC-FA-VER-001 Boundary: die Datei-Ventile sind paar-lokal — die
// exempt-paths des ersten Paares schirmen das zweite nicht ab, und die
// Quell-Datei des ersten Paares wird vom zweiten geprüft.
func TestVersionsVentilePaarLokal(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/alt.md"] = "ghcr.io/pt9912/d-check:v0.1.0 und baseline/v5.7.0\n"
	fix["version.md"] += "\nbaseline/v5.6.0\n"
	fix["docs/use.md"] = ""
	cfg := zweiReihenCfg()
	cfg.Versions.Patterns[0].ExemptPaths = []string{"docs/alt.md"}
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Target))
	}
	want := []string{"docs/alt.md:1 v5.7.0", "version.md:12 v5.6.0"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ventile müssen paar-lokal wirken: got %v, want %v", got, want)
	}
}

// DC-FA-VER-001 Boundary: der ZEILEN-Marker ist nicht paar-lokal — er nimmt
// die Zeile allen Paaren aus.
func TestVersionsZeilenMarkerGiltAllenPaaren(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.1.0 baseline/v5.7.0 <!-- d-check:ignore -->\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, zweiReihenCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("d-check:ignore nimmt die Zeile allen Paaren aus, got %v", res.Findings)
	}
}

// DC-FA-VER-001 Boundary: treffen zwei Paare dieselbe Zeile mit demselben
// Pin-Wert gegen dieselbe Erwartung, entsteht genau EIN Befund — und die
// Erwartung steht darin genau einmal.
func TestVersionsDedupGleichesTupel(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.26.0\n"
	cfg := versionsCfg()
	cfg.Versions.Patterns = append(cfg.Versions.Patterns, model.VersionPattern{
		PinPattern:  regexp.MustCompile(`d-check:(v[0-9]+\.[0-9]+\.[0-9]+)`),
		CurrentFrom: "version.md#aktuell",
	})
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("identisches Befund-Tupel zweier Paare → 1 Befund, got %v", res.Findings)
	}
	want := "Versions-Pin trägt v0.26.0, erwartet v0.27.0 (versions.patterns[0].current-from) sowie v0.27.0 (versions.patterns[1].current-from)"
	if res.Findings[0].Message != want {
		t.Fatalf("beide Paare nennen ihre Erwartung:\ngot  %q\nwant %q", res.Findings[0].Message, want)
	}
}

// DC-FA-VER-001: trifft DASSELBE Paar denselben Wert mehrfach auf einer Zeile,
// steht seine Erwartung trotzdem nur einmal in der Nachricht.
func TestVersionsGleichesPaarZweiTrefferEineErwartung(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "ghcr.io/a/d-check:v0.26.0 ghcr.io/b/d-check:v0.26.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	want := "Versions-Pin trägt v0.26.0, erwartet v0.27.0 (versions.current-from)"
	if len(res.Findings) != 1 || res.Findings[0].Message != want {
		t.Fatalf("zwei Treffer desselben Paares → eine Erwartung, got %v", res.Findings)
	}
}

// DC-FA-VER-001 / DC-QA-02: bei genau EINEM Paar nennt die Nachricht den
// Kurzform-Schlüssel — der Wortlaut bestehender Konfigurationen bleibt
// unberührt.
func TestVersionsEinPaarBehaeltWortlaut(t *testing.T) {
	fix := versionsFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.26.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, versionsCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	want := "Versions-Pin trägt v0.26.0, erwartet v0.27.0 (versions.current-from)"
	if len(res.Findings) != 1 || res.Findings[0].Message != want {
		t.Fatalf("Ein-Paar-Wortlaut: got %v", res.Findings)
	}
}

// DC-FA-VER-001 fail-closed: bei mehreren Paaren benennt auch die
// Quellauflösungs-Meldung das betroffene Paar.
func TestVersionsFehlermeldungNenntDasPaar(t *testing.T) {
	m := coretest.NewMemFS(zweiReihenFixture())
	cfg := zweiReihenCfg()
	cfg.Versions.Patterns[1].CurrentFrom = "pin.md#fehlt"
	_, err := Run(m, nil, cfg, []string{"versions"})
	if err == nil {
		t.Fatal("nicht auflösbare Quelle muss den Lauf brechen")
	}
	if !strings.Contains(err.Error(), "versions.patterns[1].current-from") {
		t.Fatalf("Meldung nennt das Paar nicht: %v", err)
	}
}

// DC-FA-VER-001 Boundary: zwei Paare, zwei verschiedene Pin-Werte auf einer
// Zeile → zwei Befunde, je gegen die eigene Erwartung.
func TestVersionsZweiWerteZweiBefunde(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v0.26.0 baseline/v5.7.0\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, zweiReihenCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.Target)
	}
	want := []string{"v0.26.0", "v5.7.0"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("zwei Werte → zwei Befunde: got %v, want %v", got, want)
	}
}

// DC-FA-VER-001 Boundary: treffen zwei Paare auf derselben Zeile DENSELBEN
// Wert mit VERSCHIEDENEN Erwartungen, entsteht ein Befund, dessen Nachricht
// beide Erwartungen in Deklarationsreihenfolge nennt. Die Befund-Adresse kann
// zwei Befunde an derselben Stelle nicht unterscheiden — ohne die Bündelung
// verschwände die zweite Erwartung still.
func TestVersionsGleicherWertBeideErwartungen(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "werkzeug:v1.0.0\n"
	cfg := model.Config{Versions: model.VersionsConfig{Patterns: []model.VersionPattern{
		{PinPattern: regexp.MustCompile(`werkzeug:(v[0-9]+\.[0-9]+\.[0-9]+)`), CurrentFrom: "version.md#aktuell"},
		{PinPattern: regexp.MustCompile(`erkzeug:(v[0-9]+\.[0-9]+\.[0-9]+)`), CurrentFrom: "pin.md#baseline"},
	}}}
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("gleiche Befund-Adresse → 1 Befund, got %v", res.Findings)
	}
	want := "Versions-Pin trägt v1.0.0, erwartet v0.27.0 (versions.patterns[0].current-from) sowie v5.9.0 (versions.patterns[1].current-from)"
	if res.Findings[0].Message != want {
		t.Fatalf("Nachricht muss beide Erwartungen nennen:\ngot  %q\nwant %q", res.Findings[0].Message, want)
	}
}

// DC-FA-VER-001: die Ausgabe-Reihenfolge ist die geteilte Sortierung
// (Datei, Zeile, Regel, Ziel, Grund) — nicht die Deklarationsreihenfolge der
// Paare. Hier findet das ZUERST deklarierte Paar den lexikografisch
// GRÖSSEREN Wert; gemeldet wird der kleinere zuerst.
func TestVersionsAusgabeFolgtDerSortierung(t *testing.T) {
	fix := zweiReihenFixture()
	fix["docs/use.md"] = "ghcr.io/pt9912/d-check:v9.9.9 baseline/v1.1.1\n"
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, zweiReihenCfg(), []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.Target)
	}
	want := []string{"v1.1.1", "v9.9.9"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ausgabe folgt der Sortierung: got %v, want %v", got, want)
	}
}

// DC-FA-VER-001 fail-closed: löst die Quelle des ZWEITEN Paares nicht auf,
// bricht der Lauf — ein unvollständiges Paar schaltet die übrigen nicht still
// scharf.
func TestVersionsZweitesPaarQuelleFehltFailClosed(t *testing.T) {
	m := coretest.NewMemFS(zweiReihenFixture())
	cfg := zweiReihenCfg()
	cfg.Versions.Patterns[1].CurrentFrom = "pin.md#fehlt"
	if _, err := Run(m, nil, cfg, []string{"versions"}); err == nil {
		t.Fatal("nicht auflösbare Quelle eines Paares muss den Lauf brechen")
	}
}
