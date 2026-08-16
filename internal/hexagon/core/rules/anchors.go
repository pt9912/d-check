// Package rules trägt die I/O-freie Prüf-Engine des Kerns: die Regelmodule,
// die Markdown-/Pfad-Primitive sowie die Prüf-Orchestrierung (Run, Discovery).
// Importiert nur model und die Ports (spec/architecture.md §Kern; ADR-0012).
package rules

import (
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// ReasonAnchorMissing — Grund-Code des Moduls anchors
// (spec/spezifikation.md §4).
const ReasonAnchorMissing = "anchor-missing"

// ExtractHeadings liefert die Texte aller ATX-Headings (#–######)
// außerhalb von Fenced-Code-Blöcken, in Dokumentreihenfolge
// (DC-FA-ANCH-001; gemeinsamer Scanner: extractHeadingLines).
func ExtractHeadings(content []byte) []string {
	var out []string
	for _, h := range extractHeadingLines(content) {
		out = append(out, h.text)
	}
	return out
}

// StripHeadingLinks ersetzt Markdown-Links im Heading-Text durch ihren
// Linktext (Slug-Schritt 1).
func StripHeadingLinks(s string) string {
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
	t := strings.ToLower(StripHeadingLinks(heading))
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

// headingSlugsOrdered liefert je Überschrift ihren **effektiven** Slug in
// Dokumentreihenfolge, inklusive Duplikat-Suffix `-1`, `-2`, … (Slug-Schritt 5).
// Dies ist die EINE Slug-Antwort des Produkts; HeadingSlugs ist ihre
// Mengen-Sicht, slugSection ihre Zeilen-Sicht.
func headingSlugsOrdered(headings []string) []string {
	counts := map[string]int{}
	out := make([]string, 0, len(headings))
	for _, h := range headings {
		base := Slugify(h)
		n := counts[base]
		counts[base]++
		if n > 0 {
			out = append(out, fmt.Sprintf("%s-%d", base, n))
			continue
		}
		out = append(out, base)
	}
	return out
}

// HeadingSlugs liefert dieselben Slugs als Menge.
func HeadingSlugs(content []byte) map[string]bool {
	set := map[string]bool{}
	for _, s := range headingSlugsOrdered(ExtractHeadings(content)) {
		set[s] = true
	}
	return set
}

// DecodeFragment macht aus einem Link-Fragment die Zeichenkette, gegen die
// verglichen wird: prozent-kodierte Sequenzen werden aufgelöst (`#a%20b`
// adressiert den Anker `a b`). Dieselbe Auflösung für jeden Konsumenten —
// sonst trifft derselbe Anker im selben Lauf einmal und einmal nicht.
func DecodeFragment(frag string) string {
	if dec, err := url.PathUnescape(frag); err == nil {
		return dec
	}
	return frag
}

var (
	// htmlTagRE erfasst öffnende HTML-Tags (Tag-Name + Attribut-Teil)
	// auf vorverarbeiteten Zeilen; Inline-Code ist dort bereits geleert.
	// Der Attribut-Teil erlaubt gequotete Werte (die ein `>` enthalten
	// dürfen), damit ein literales `>` im Attribut den Tag nicht
	// vorzeitig beendet.
	htmlTagRE = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)((?:"[^"]*"|'[^']*'|[^>'"])*)>`)
	// htmlAttrIDRE/htmlAttrNameRE lesen den Wert eines id- bzw.
	// name-Attributs aus dem Attribut-Teil (doppelte oder einfache
	// Anführungszeichen; Attributname an Wortgrenze, kein Treffer in
	// `data-id`).
	htmlAttrIDRE   = regexp.MustCompile(`(?i)(?:^|\s)id\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	htmlAttrNameRE = regexp.MustCompile(`(?i)(?:^|\s)name\s*=\s*(?:"([^"]*)"|'([^']*)')`)
)

// htmlAnchorLines liefert je Inline-HTML-Anker die 1-basierte Zeile seines
// ersten Vorkommens (DC-FA-ANCH-001.b): id-Werte an beliebigen Elementen und
// name-Werte an <a>-Elementen, wörtlich. Erkennung auf den vorverarbeiteten
// Zeilen — Fenced-Code-Blöcke und Inline-Code-Spans sind dort entfernt (GitHub
// rendert HTML in Code-Auszeichnung nicht als Sprungziel).
//
// Dies ist die EINE Anker-Erkennung des Produkts; htmlAnchors ist ihre
// Namens-Sicht, headingSection ihre Zeilen-Sicht.
func htmlAnchorLines(content []byte) map[string]int {
	out := map[string]int{}
	add := func(v string, no int) {
		if v == "" {
			return
		}
		if _, seen := out[v]; !seen {
			out[v] = no
		}
	}
	for _, ln := range PreprocessMarkdown(content) {
		for _, tag := range htmlTagRE.FindAllStringSubmatch(ln.Text, -1) {
			add(attrValue(htmlAttrIDRE, tag[2]), ln.No)
			if strings.EqualFold(tag[1], "a") {
				add(attrValue(htmlAttrNameRE, tag[2]), ln.No)
			}
		}
	}
	return out
}

// htmlAnchors liefert dieselben Anker als Namens-Menge.
func htmlAnchors(content []byte) map[string]bool {
	set := map[string]bool{}
	for a := range htmlAnchorLines(content) {
		set[a] = true
	}
	return set
}

// attrValue liefert den erfassten Attributwert (doppelte vor einfachen
// Anführungszeichen) oder "". Ein leerer Wert (`id=""`) ist von „kein
// Treffer" nicht unterscheidbar und erzeugt bewusst keinen Anker —
// ein leeres Fragment ist als Sprungziel ohnehin nutzlos.
func attrValue(re *regexp.Regexp, attrs string) string {
	m := re.FindStringSubmatch(attrs)
	if m == nil {
		return ""
	}
	if m[1] != "" {
		return m[1]
	}
	return m[2]
}

// AnchorSet ist die gültige Anker-Menge einer Datei: Heading-Slugs
// (DC-FA-ANCH-001.a) vereinigt mit Inline-HTML-Ankern
// (DC-FA-ANCH-001.b). Geteilt von den Modulen anchors und codepaths.
func AnchorSet(content []byte) map[string]bool {
	set := HeadingSlugs(content)
	for a := range htmlAnchors(content) {
		set[a] = true
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
	frag := DecodeFragment(t[idx+1:])
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

// CheckAnchors ist das Regelmodul `anchors` (DC-FA-ANCH-001): Links
// mit Fragment werden gegen die Heading-Slugs der Zieldatei geprüft.
func CheckAnchors(fsys driven.Filesystem, file string, content []byte, lines []Line, cache map[string]map[string]bool, ignoreRefs []model.IgnoreRef) []model.Finding {
	var findings []model.Finding
	for _, ref := range ExtractLinks(lines) {
		a, ok := resolveAnchorRef(fsys, file, ref)
		if !ok {
			continue
		}
		// Geteiltes Referenz-Ventil (DC-FA-REF-001): ist das Ziel
		// ignoriert, entfällt die Anker-Prüfung.
		if refIgnored(ignoreRefs, file, a.rel) {
			continue
		}
		slugs := slugsFor(fsys, cache, a, content)
		if slugs == nil || slugs[a.frag] {
			continue
		}
		findings = append(findings, model.Finding{
			File: file, Line: a.line, Rule: "anchors",
			Target: a.target, Reason: ReasonAnchorMissing,
			Message: "Anker entspricht keinem Heading-Slug und keinem HTML-Anker der Zieldatei",
		})
	}
	return findings
}

// slugsFor liefert die gültige Anker-Menge der Zieldatei aus dem Cache
// bzw. liest sie nach (nil = nicht lesbar → Modul schweigt).
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
	s := AnchorSet(content)
	cache[a.rel] = s
	return s
}
