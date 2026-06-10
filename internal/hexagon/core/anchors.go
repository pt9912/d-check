package core

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonAnchorMissing — Grund-Code des Moduls anchors
// (spec/spezifikation.md §4).
const ReasonAnchorMissing = "anchor-missing"

// ExtractHeadings liefert die Texte aller ATX-Headings (#–######)
// außerhalb von Fenced-Code-Blöcken, in Dokumentreihenfolge
// (DC-FA-ANCH-001; Setext wird in 0.x nicht unterstützt —
// spec/spezifikation.md §DC-FA-ANCH-001.a).
func ExtractHeadings(content []byte) []string {
	var out []string
	inFence := false
	for _, raw := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimLeft(raw, " \t")
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence || !strings.HasPrefix(trimmed, "#") {
			continue
		}
		level := 0
		for level < len(trimmed) && trimmed[level] == '#' {
			level++
		}
		if level > 6 || level >= len(trimmed) || (trimmed[level] != ' ' && trimmed[level] != '\t') {
			continue
		}
		out = append(out, strings.TrimSpace(trimmed[level+1:]))
	}
	return out
}

// stripHeadingLinks ersetzt Markdown-Links im Heading-Text durch ihren
// Linktext (Slug-Schritt 1).
func stripHeadingLinks(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '!' && i+1 < len(s) && s[i+1] == '[' {
			continue // Bild-Bang entfällt, [alt](…) folgt
		}
		if s[i] != '[' {
			b.WriteByte(s[i])
			continue
		}
		depth := 0
		j := i
		for ; j < len(s); j++ {
			if s[j] == '[' {
				depth++
			} else if s[j] == ']' {
				depth--
				if depth == 0 {
					break
				}
			}
		}
		if j >= len(s) || j+1 >= len(s) || s[j+1] != '(' {
			b.WriteByte(s[i])
			continue
		}
		k := strings.IndexByte(s[j+1:], ')')
		if k == -1 {
			b.WriteByte(s[i])
			continue
		}
		b.WriteString(s[i+1 : j]) // Linktext
		i = j + 1 + k
	}
	return b.String()
}

// Slugify bildet den GitHub-Slug eines Heading-Texts
// (spec/spezifikation.md §DC-FA-ANCH-001.a, Schritte 1–4): Markup
// entfernen (Backticks, Emphasis-Sterne; Links → Linktext),
// Unicode-Kleinschreibung, nur Buchstaben/Ziffern/Leerzeichen/-/_
// behalten, jedes Leerzeichen → '-'.
func Slugify(heading string) string {
	t := strings.ToLower(stripHeadingLinks(heading))
	var b strings.Builder
	for _, r := range t {
		switch {
		case r == '`' || r == '*':
			// Markup-Zeichen entfallen, der Text bleibt
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ' || r == '\t':
			b.WriteByte('-')
		}
	}
	return b.String()
}

// HeadingSlugs liefert die Slug-Menge einer Datei inkl.
// Duplikat-Suffixen `-1`, `-2`, … in Dokumentreihenfolge
// (Slug-Schritt 5).
func HeadingSlugs(content []byte) map[string]bool {
	counts := map[string]int{}
	set := map[string]bool{}
	for _, h := range ExtractHeadings(content) {
		base := Slugify(h)
		n := counts[base]
		counts[base]++
		s := base
		if n > 0 {
			s = fmt.Sprintf("%s-%d", base, n)
		}
		set[s] = true
	}
	return set
}

// checkAnchors ist das Regelmodul `anchors` (DC-FA-ANCH-001): Links
// mit Fragment werden gegen die Heading-Slugs der Zieldatei geprüft.
// Existiert die Zieldatei nicht (oder ist sie kein .md / ein
// Symlink/Escape-Fall), schweigt das Modul — diese Befunde gehören
// dem Modul `links`.
func checkAnchors(fsys driven.Filesystem, file string, content []byte, lines []Line, cache map[string]map[string]bool) []Finding {
	getSlugs := func(rel string, own []byte) map[string]bool {
		if s, ok := cache[rel]; ok {
			return s
		}
		c := own
		if c == nil {
			read, err := fsys.ReadFile(rel)
			if err != nil {
				cache[rel] = nil
				return nil
			}
			c = read
		}
		s := HeadingSlugs(c)
		cache[rel] = s
		return s
	}

	var findings []Finding
	for _, ref := range ExtractLinks(lines) {
		t := ref.Target
		if t == "" || IsExternalScheme(t) {
			continue
		}
		idx := strings.IndexByte(t, '#')
		if idx == -1 {
			continue
		}
		frag := t[idx+1:]
		if frag == "" {
			continue
		}
		if dec, err := url.PathUnescape(frag); err == nil {
			frag = dec
		}
		var rel string
		var own []byte
		if idx == 0 {
			rel, own = file, content
		} else {
			r, escaped, _ := ResolveTarget(file, t[:idx])
			if escaped || !strings.HasSuffix(r, ".md") {
				continue
			}
			if kind, err := fsys.Kind(r); err != nil || kind != driven.KindFile {
				continue // fehlende Datei/Symlink meldet `links`
			}
			rel = r
		}
		slugs := getSlugs(rel, own)
		if slugs == nil {
			continue
		}
		if !slugs[frag] {
			findings = append(findings, Finding{
				File: file, Line: ref.Line, Rule: "anchors",
				Target: t, Reason: ReasonAnchorMissing,
				Message: "Anker entspricht keinem Heading-Slug der Zieldatei",
			})
		}
	}
	return findings
}
