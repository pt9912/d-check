package app

import (
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

const sliceNamesReq = "XX-FA-ABC-001"

// sliceRefsOf liefert die Slice-Kennungen, die sliceNamesReq decken, für die
// Dateien files (Pfad → Inhalt) unter tc (DC-FA-CLI-009).
func sliceRefsOf(t *testing.T, tc model.TraceConfig, files map[string]string) []string {
	t.Helper()
	rt, err := resolveTrace(tc)
	if err != nil {
		t.Fatalf("resolveTrace: %v", err)
	}
	refs, err := traceRefs(coretest.NewMemFS(files), rt.sliceDir, rt.sliceOwner(), regexp.MustCompile(`XX-FA-[A-Z]+-\d+`))
	if err != nil {
		t.Fatalf("traceRefs: %v", err)
	}
	return refs[sliceNamesReq]
}

const sliceBody = "# Slice\n\nBezug: " + sliceNamesReq + "\n"

// Happy: ein Slice mit Namens-Kennung deckt die Anforderung ohne Konfiguration.
func TestTraceSliceNameDefaultCovers(t *testing.T) {
	got := sliceRefsOf(t, model.TraceConfig{}, map[string]string{
		"docs/plan/planning/slice-qa-sicherheit.md": sliceBody,
	})
	if !equalStrings(got, []string{"slice-qa-sicherheit"}) {
		t.Fatalf("Namens-Slice: got %v", got)
	}
}

// Boundary: der Ziffern-Fall behält die Kennung `slice-NNN`, auch neben
// Namens-Slices; ein Slice ohne Titel ist ein Namens-Fall mit derselben Kennung.
func TestTraceSliceDigitCaseUnchanged(t *testing.T) {
	got := sliceRefsOf(t, model.TraceConfig{}, map[string]string{
		"docs/plan/planning/slice-039-foo.md":       sliceBody,
		"docs/plan/planning/slice-273.md":           sliceBody,
		"docs/plan/planning/slice-qa-sicherheit.md": sliceBody,
	})
	want := []string{"slice-039", "slice-273", "slice-qa-sicherheit"}
	if !equalStrings(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Negative: kein Treffer für Namens-Belege unter observations, Punkt im Stamm, leeren Namen,
// Großschreibung und Dateien ohne `slice-`-Präfix.
func TestTraceSliceNameNegatives(t *testing.T) {
	got := sliceRefsOf(t, model.TraceConfig{}, map[string]string{
		"docs/plan/planning/observations/BEO-X/y/evidence/slice-qa-sicherheit.md": sliceBody,
		"docs/plan/planning/slice-a.b.md":                                          sliceBody,
		"docs/plan/planning/slice-.md":                                             sliceBody,
		"docs/plan/planning/Slice-X.md":                                            sliceBody,
		"docs/plan/planning/README.md":                                             sliceBody,
	})
	if len(got) != 0 {
		t.Fatalf("keine Kennung erwartet, got %v", got)
	}
}

// Ein Stub und ein Plan mit gleichem Basisnamen sind eine Kennung, und
// `slice-a-foo` neben `slice-a-foo-bar` sind zwei — kein Präfix-Treffer.
func TestTraceSliceNameIdentity(t *testing.T) {
	got := sliceRefsOf(t, model.TraceConfig{}, map[string]string{
		"docs/plan/planning/done/slice-a-foo.md":     sliceBody,
		"docs/plan/planning/done/w/slice-a-foo.md":   sliceBody,
		"docs/plan/planning/open/slice-a-foo-bar.md": sliceBody,
	})
	want := []string{"slice-a-foo", "slice-a-foo-bar"}
	if !equalStrings(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Ein gesetztes file-pattern gilt unverändert und ohne `observations`-Ausschluss;
// der Default-Namens-Fall greift dann nicht.
func TestTraceSliceFilePatternOverrides(t *testing.T) {
	tc := model.TraceConfig{SliceFile: regexp.MustCompile(`^(slice-[a-z-]+)\.md$`)}
	got := sliceRefsOf(t, tc, map[string]string{
		"docs/plan/planning/observations/BEO-X/y/evidence/slice-qa-sicherheit.md": sliceBody,
		"docs/plan/planning/slice-039-foo.md":                                      sliceBody,
	})
	if !equalStrings(got, []string{"slice-qa-sicherheit"}) {
		t.Fatalf("file-pattern muss ohne Ausschluss und ohne Default-Ziffernfall gelten: %v", got)
	}
}

// Eine fremde slices.dir, die selbst `observations` heißt, bleibt scanbar; nur
// ein Segment unterhalb von dir schließt aus.
func TestTraceSliceDirNamedObservations(t *testing.T) {
	tc := model.TraceConfig{SliceDir: "observations"}
	got := sliceRefsOf(t, tc, map[string]string{
		"observations/slice-qa-sicherheit.md":         sliceBody,
		"observations/sub/observations/slice-x-y.md": sliceBody,
	})
	if !equalStrings(got, []string{"slice-qa-sicherheit"}) {
		t.Fatalf("got %v", got)
	}
}

// Der observations-Ausschluss trifft nur den Namens-Fall: ein Beleg mit
// Ziffern-Kennung und Titel blieb schon vor der Namens-Erkennung ein Slice und
// bleibt es.
func TestTraceSliceDigitCaseIgnoresObservationsExclusion(t *testing.T) {
	got := sliceRefsOf(t, model.TraceConfig{}, map[string]string{
		"docs/plan/planning/observations/BEO-X/y/evidence/slice-007-foo.md": sliceBody,
	})
	if !equalStrings(got, []string{"slice-007"}) {
		t.Fatalf("Ziffern-Fall unter observations bleibt Slice: %v", got)
	}
}

// Eine nicht bereinigte slices.dir (`observations/`) ändert
// nichts: das Segment in dir selbst schließt nicht aus.
func TestTraceSliceDirUncleaned(t *testing.T) {
	for _, dir := range []string{"observations/"} {
		got := sliceRefsOf(t, model.TraceConfig{SliceDir: dir}, map[string]string{
			"observations/slice-qa-sicherheit.md":        sliceBody,
			"observations/sub/observations/slice-x-y.md": sliceBody,
		})
		if !equalStrings(got, []string{"slice-qa-sicherheit"}) {
			t.Fatalf("dir %q: got %v", dir, got)
		}
	}
}
