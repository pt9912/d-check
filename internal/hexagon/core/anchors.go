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
		textEnd, ok := matchBracket(s, i, '[', ']')
		if !ok || textEnd+1 >= len(s) || s[textEnd+1] != '(' {
			b.WriteByte(s[i])
			continue
		}
		destEnd, ok := matchBracket(s, textEnd+1, '(', ')')
		if !ok {
			b.WriteByte(s[i])
			continue
		}
		b.WriteString(s[i+1 : textEnd]) // Linktext
		i = destEnd
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

// anchorRef ist ein Link mit Fragment, dessen Ziel auflösbar ist.
type anchorRef struct {
	line   int
	target string // Original-Zielausdruck (für den Befund)
	frag   string // dekodiertes Fragment
	rel    string // Zieldatei (relativ zur Wurzel)
	own    bool   // Ziel ist die prüfende Datei selbst
}

// resolveAnchorRef klassifiziert einen Link für das Modul anchors;
// ok=false bedeutet: nicht zuständig (kein Fragment, extern, Escape,
// Nicht-Markdown, fehlende Datei — Letzteres meldet `links`).
func resolveAnchorRef(fsys driven.Filesystem, file string, ref LinkRef) (anchorRef, bool) {
	t := ref.Target
	if t == "" || IsExternalScheme(t) {
		return anchorRef{}, false
	}
	idx := strings.IndexByte(t, '#')
	if idx == -1 || idx+1 >= len(t) {
		return anchorRef{}, false
	}
	frag := t[idx+1:]
	if dec, err := url.PathUnescape(frag); err == nil {
		frag = dec
	}
	out := anchorRef{line: ref.Line, target: t, frag: frag}
	if idx == 0 {
		out.rel, out.own = file, true
		return out, true
	}
	rel, escaped, _ := ResolveTarget(file, t[:idx])
	if escaped || !strings.HasSuffix(rel, ".md") {
		return anchorRef{}, false
	}
	if kind, err := fsys.Kind(rel); err != nil || kind != driven.KindFile {
		return anchorRef{}, false
	}
	out.rel = rel
	return out, true
}

// checkAnchors ist das Regelmodul `anchors` (DC-FA-ANCH-001): Links
// mit Fragment werden gegen die Heading-Slugs der Zieldatei geprüft.
func checkAnchors(fsys driven.Filesystem, file string, content []byte, lines []Line, cache map[string]map[string]bool) []Finding {
	var findings []Finding
	for _, ref := range ExtractLinks(lines) {
		a, ok := resolveAnchorRef(fsys, file, ref)
		if !ok {
			continue
		}
		slugs := slugsFor(fsys, cache, a, content)
		if slugs == nil || slugs[a.frag] {
			continue
		}
		findings = append(findings, Finding{
			File: file, Line: a.line, Rule: "anchors",
			Target: a.target, Reason: ReasonAnchorMissing,
			Message: "Anker entspricht keinem Heading-Slug der Zieldatei",
		})
	}
	return findings
}

// slugsFor liefert die Slug-Menge der Zieldatei aus dem Cache bzw.
// liest sie nach (nil = nicht lesbar → Modul schweigt).
func slugsFor(fsys driven.Filesystem, cache map[string]map[string]bool, a anchorRef, own []byte) map[string]bool {
	if s, ok := cache[a.rel]; ok {
		return s
	}
	content := own
	if !a.own {
		read, err := fsys.ReadFile(a.rel)
		if err != nil {
			cache[a.rel] = nil
			return nil
		}
		content = read
	}
	s := HeadingSlugs(content)
	cache[a.rel] = s
	return s
}
