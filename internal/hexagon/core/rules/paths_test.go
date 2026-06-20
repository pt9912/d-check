package rules

import "testing"

func TestResolveTarget(t *testing.T) {
	cases := []struct {
		from, target string
		wantRel      string
		wantEscaped  bool
	}{
		{"docs/a.md", "b.md", "docs/b.md", false},
		{"docs/a.md", "../README.md", "README.md", false},
		{"docs/a.md", "../../etc/passwd", "../etc/passwd", true},
		{"a.md", "docs/x.md", "docs/x.md", false},
		// Prozent-Dekodierung VOR der Escape-Prüfung (DC-FA-LINK-001)
		{"docs/a.md", "mit%20leer.md", "docs/mit leer.md", false},
		{"docs/a.md", "..%2F..%2Fetc", "../etc", true},
		// absolute Ziele relativ zur Repo-Wurzel
		{"docs/a.md", "/spec/x.md", "spec/x.md", false},
	}
	for _, c := range cases {
		rel, escaped, _ := ResolveTarget(c.from, c.target)
		if rel != c.wantRel || escaped != c.wantEscaped {
			t.Errorf("ResolveTarget(%q,%q) = (%q,%v), want (%q,%v)",
				c.from, c.target, rel, escaped, c.wantRel, c.wantEscaped)
		}
	}
}

func TestIsExternalScheme(t *testing.T) {
	for target, want := range map[string]bool{
		"https://example.org": true,
		"mailto:a@b":          true,
		"ftp://x":             true,
		"docs/a.md":           false,
		"a:b.md":              true, // generisches Schema
		"./a:b":               false,
	} {
		if got := IsExternalScheme(target); got != want {
			t.Errorf("IsExternalScheme(%q) = %v, want %v", target, got, want)
		}
	}
}

func TestMatchGlob(t *testing.T) {
	for _, c := range []struct {
		pattern, rel string
		want         bool
	}{
		{"docs/archive/**", "docs/archive/alt.md", true},
		{"docs/archive/**", "docs/archive/tief/x.md", true},
		{"docs/archive/**", "docs/aktuell.md", false},
		{"*.md", "README.md", true},
		{"*.md", "docs/a.md", false},
		{"docs/*/intern.md", "docs/x/intern.md", true},
		// '**' matcht auch NULL Segmente (DC-FA-SCAN-001.a)
		{"docs/**/slice-*.md", "docs/slice-001.md", true},
		{"docs/**/slice-*.md", "docs/a/b/slice-001.md", true},
		{"docs/**/slice-*.md", "docs/a/b/notiz.md", false},
		// '?' = genau ein Zeichen; Muster kürzer als Pfad
		{"docs/f?.md", "docs/f1.md", true},
		{"docs/f?.md", "docs/f12.md", false},
		{"docs", "docs/a.md", false},
	} {
		if got := MatchGlob(c.pattern, c.rel); got != c.want {
			t.Errorf("MatchGlob(%q,%q) = %v, want %v", c.pattern, c.rel, got, c.want)
		}
	}
}
