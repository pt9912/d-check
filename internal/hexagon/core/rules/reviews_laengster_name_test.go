package rules

import (
	"testing"
)

// Unter match: name deckt ein Report nur den längsten Slice-Basisnamen, der
// in seinem Namen steht (DC-FA-RVW-001): der Report zu slice-cache-warmup
// deckt slice-cache nicht, wenn slice-cache-warmup ein Slice ist.

const rvZusage = "## 2. Definition of Done\n\n- [x] Review durchgeführt\n"

// Zwei Slices, ein Report für den längeren: der kürzere meldet, bis er einen
// eigenen Report hat.
func TestReviewsMatchName_LaengsterNameGewinnt(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	files := map[string]string{
		rvDone + "/slice-cache.md":                     rvZusage,
		rvDone + "/slice-cache-warmup.md":              rvZusage,
		"docs/reviews/2026-10-09-slice-cache-warmup.md": "# Review\n",
	}
	f := rvRunCfg(files, cfg)
	if len(f) != 1 || f[0].File != rvDone+"/slice-cache.md" {
		t.Fatalf("slice-cache ohne eigenen Report meldet, slice-cache-warmup nicht, got %+v", f)
	}
	files["docs/reviews/2026-10-09-slice-cache.md"] = "# Review\n"
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("mit eigenem Report sind beide gedeckt, got %+v", f)
	}
}

// Ein archivierter Stub des längeren Namens zählt mit: sein Report deckt den
// kürzeren nicht, auch wenn skip-pattern den Stub aus der Prüfung nimmt.
func TestReviewsMatchName_StubDesLaengerenZaehlt(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	cfg.SkipPattern = stubPattern
	files := map[string]string{
		rvDone + "/slice-cache.md":                     rvZusage,
		rvDone + "/slice-cache-warmup.md":              "# slice-cache-warmup\n\n" + stubMarker,
		"docs/reviews/2026-10-09-slice-cache-warmup.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-cache.md" {
		t.Fatalf("der Stub von slice-cache-warmup zählt als längerer Name, got %+v", f)
	}
}

// Gegenprobe zur Wortgrenze: der Report zu slice-cachex deckt slice-cache
// schon deshalb nicht, weil nach slice-cache ein Buchstabe folgt.
func TestReviewsMatchName_KeinePraefixKollisionOhneGrenze(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	files := map[string]string{
		rvDone + "/slice-cache.md":                rvZusage,
		rvDone + "/slice-cachex.md":               rvZusage,
		"docs/reviews/2026-10-09-slice-cachex.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-cache.md" {
		t.Fatalf("slice-cachex deckt slice-cache nicht, got %+v", f)
	}
}

// Ein über exempt-paths ausgenommener längerer Slice zählt mit: sein Report
// deckt den kürzeren nicht.
func TestReviewsMatchName_ExemptDesLaengerenZaehlt(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	cfg.ExemptPaths = []string{rvDone + "/slice-cache-warmup.md"}
	files := map[string]string{
		rvDone + "/slice-cache.md":                      rvZusage,
		rvDone + "/slice-cache-warmup.md":               rvZusage,
		"docs/reviews/2026-10-09-slice-cache-warmup.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-cache.md" {
		t.Fatalf("der ausgenommene slice-cache-warmup zählt als längerer Name, got %+v", f)
	}
}

// Ein Stub des längeren Namens unter einem Unterverzeichnis zählt mit
// recursive mit; ohne recursive sieht der Lauf ihn nicht, und der Report
// deckt den kürzeren weiter.
func TestReviewsMatchName_StubImUnterverzeichnis(t *testing.T) {
	cfg := rvCfg()
	cfg.Match = "name"
	cfg.SkipPattern = stubPattern
	files := map[string]string{
		rvDone + "/slice-cache.md":                      rvZusage,
		rvDone + "/welle-1/slice-cache-warmup.md":       "# slice-cache-warmup\n\n" + stubMarker,
		"docs/reviews/2026-10-09-slice-cache-warmup.md": "# Review\n",
	}
	if f := rvRunCfg(files, cfg); f != nil {
		t.Fatalf("ohne recursive ist der Stub ungesehen, der Report deckt slice-cache, got %+v", f)
	}
	cfg.Recursive = true
	if f := rvRunCfg(files, cfg); len(f) != 1 || f[0].File != rvDone+"/slice-cache.md" {
		t.Fatalf("mit recursive zählt der Stub als längerer Name, got %+v", f)
	}
}
