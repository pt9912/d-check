package rules

import (
	"testing"
)

// Unter match: name deckt ein Report nur den längsten Slice-Basisnamen, der
// in seinem Namen steht (DC-FA-RVW-001): der Report zu slice-cache-warmup
// deckt slice-cache nicht, wenn slice-cache-warmup ein Slice ist.

const rvZusage = "## 2. Definition of Done\n\n- [x] Review durchgeführt\n"

// Der Fall des Befunds: zwei Slices, ein Report für den längeren.
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

// Gegenprobe: slice-cachex ist kein längerer Name im Sinn der Zuordnung — er
// steht nicht vor einer Wortgrenze, sein Report deckt slice-cache ohnehin nicht.
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
