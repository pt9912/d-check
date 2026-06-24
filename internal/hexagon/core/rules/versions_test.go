package rules

import (
	"fmt"
	"regexp"
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
	return model.Config{Versions: model.VersionsConfig{
		PinPattern:  regexp.MustCompile(`ghcr\.io/[^\s:]+:(v[0-9]+\.[0-9]+\.[0-9]+)`),
		CurrentFrom: "version.md#aktuell",
	}}
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
	cfg.Versions.ExemptPaths = []string{"docs/alt.md"}
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
	cfg.Versions.CurrentFrom = "version.md#fehlt"
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
	cfg.Versions.CurrentFrom = "version.md#cur"
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
	cfg.Versions.CurrentFrom = "reg.md"
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
	cfg.Versions.PinPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
	m := coretest.NewMemFS(fix)
	res, err := Run(m, nil, cfg, []string{"versions"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Target != "v0.26.0" {
		t.Fatalf("ohne Capture-Gruppe = ganzer Treffer, got %v", res.Findings)
	}
}
