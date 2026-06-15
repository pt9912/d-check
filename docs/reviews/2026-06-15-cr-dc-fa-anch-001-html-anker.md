# Review-Report: CR `DC-FA-ANCH-001` — Inline-HTML-Anker (slice-022) — 2026-06-15

**Review-Art:** Dokument-Review **vor Implementierung** (Change Request +
Spezifikations-Schärfung + Slice-Plan). Kein Code im Umfang — `anchors.go`
ist unangetastet.

**Gegenstand:** Uncommitteter Working-Tree-Stand: `spec/lastenheft.md`
(`DC-FA-ANCH-001` auf 0.12.0), `spec/spezifikation.md` (neu
§`DC-FA-ANCH-001.b`, fortgeschriebener §`DC-FA-CODE-001.a`-Anker-Bezug),
`docs/plan/planning/in-progress/slice-022-html-anker.md` (seither aus
`open/` dorthin verschoben), Roadmap.

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md)
· **Datum:** 2026-06-15

**Eingangs-Kontext:** [`spec/lastenheft.md`](../../spec/lastenheft.md)
`DC-FA-ANCH-001` / `DC-FA-CODE-001`,
[`spec/spezifikation.md`](../../spec/spezifikation.md) §`DC-FA-ANCH-001.a`/`.b`,
`DC-QA-02`/`DC-QA-03`, [`.d-check.yml`](../../.d-check.yml),
[`AGENTS.md`](../../AGENTS.md) §3, `MR-006`.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar |
|---|---|---|---|---|---|
| 1 | 🟠 MEDIUM | `DC-FA-ANCH-001` / Spec-Treue (GitHub-Render-Modell) | `spec/spezifikation.md` Z. 182 | §`DC-FA-ANCH-001.b` schließt nur Fenced-Code-Blöcke aus und schweigt zu Inline-Code-Spans. GitHub rendert HTML in Inline-Code (`` `id="x"` ``) nicht als Anker; bei wörtlicher Aufnahme würden Attributwerte aus Inline-Code — auch die Beispiele in Lastenheft und Spezifikation selbst — als reale Anker geführt, sodass ein anderswo kaputtes gleichnamiges Fragment unbemängelt bleibt. Das Negativkriterium testet nur den Fenced-Fall. | ja — Fixture mit `id="phantom"` in Inline-Code plus Link `#phantom` ohne realen Anker; `make doc-check` muss `anchor-missing` melden, bliebe bei Inline-Code-Einschluss aus |
| 2 | 🟠 MEDIUM | `DC-FA-CODE-001` / Konsistenz Rang 1 ↔ Rang 2 | `spec/lastenheft.md` Z. 404 vs. `spec/spezifikation.md` Z. 317 | Das Lastenheft beschreibt die `codepaths`-Anker-Prüfung als „gegen die Headings der Zieldatei … gleiches Slug-Verfahren"; die Spezifikation prüft jetzt „gegen die gültige Anker-Menge (Heading-Slugs und Inline-HTML-Anker)". Rang 2 erweitert `codepaths` über die unveränderte Rang-1-Formulierung hinaus. | nein — Dokument-Semantik Lastenheft↔Spezifikation; kein Gate prüft sie, nur Review |
| 3 | 🟡 LOW | Maintainability / Doku-Drift | `spec/lastenheft.md` Z. 608 | §5 listet „HTML" unverändert unter den nicht geprüften Nicht-Markdown-Formaten; die Abgrenzung trägt nur `DC-FA-ANCH-001` (Z. 292). Ein bei §5 einsteigender Leser sieht den scheinbaren Widerspruch zur neuen Fähigkeit ohne Querverweis. | nein — Prosa-Konsistenz |
| 4 | 🔵 INFO | undokumentierte Annahme (Implementierung) | `spec/spezifikation.md` Z. 186 | „`name` an einem `<a>`-Element" lässt offen, dass eine zeilenbasierte Erkennung den Tag-Namen exakt treffen muss (Wortgrenze). Ein Präfix-Match zündet sonst fälschlich bei `<area>`/`<abbr>`, wo GitHub `name` nicht als Anker honoriert. | ja — späterer Unit-Test: `<area name="x">` darf keinen Anker ergeben |
| 5 | 🔵 INFO | bewusste Won't-Fix-Designnotiz | `spec/spezifikation.md` Z. 188 | Die permissive Richtung (irrtümlich erkannter Anker ⇒ Modul schweigt) kann einen real kaputten gleichnamigen Anker maskieren; in der Spezifikation als bewusst akzeptiert und konsistent mit dem Schweige-Charakter des Moduls dokumentiert. Hier festgehalten, damit „gesehen" von „übersehen" unterscheidbar bleibt. | nein — Designentscheidung |

## Negativbefunde (geprüft, ohne Befund)

- **Referenzrichtung (`MR-006`/matrix):** kein Spec→Slice- oder Spec→ADR-Link eingeführt; „slice-022" steht in beiden Historien als Klartext; `make doc-check` matrix grün.
- **Linkpflicht (`ids`, `link-policy: always`):** alle `DC-`-Kennungen in Spezifikation, Slice und Roadmap verlinkt; im Lastenheft target-exempt; doc-check 0 Befunde über 48 Dateien.
- **Anker-Auflösung:** der neue `.b`-Anker und alle Querverweise lösen auf; doc-check anchors grün (das eigentliche Dogfooding-Risiko der Änderung).
- **Akzeptanzkriterien-Form:** Happy/Boundary/Negative für das neue HTML-Verhalten plus neu gefasste Out-of-Scope-Liste vorhanden (Anforderungs-Anlege-Prozess erfüllt).
- **ADR-Bedarf:** kein neuer Port/Adapter, keine Dependency; `make arch-check` grün (Import-Regeln R1–R5, `ADR-0005`).
- **Abwärtskompatibilität:** die Änderung vergrößert ausschließlich die gültige Anker-Menge ⇒ reduziert Falsch-Befunde, erzeugt nie neue; kein bestehendes Akzeptanzkriterium kippt.
- **Determinismus (`DC-QA-02`):** Anker-Menge als ungeordnete Mengenmitgliedschaft, reihenfolgeunabhängig; keine neue Ordnungsabhängigkeit gegenüber dem heutigen Slug-Cache.
- **Case-Sensitivität:** Slug-Menge kleingeschrieben, HTML-Anker wörtlich, Fragment unverändert verglichen — entspricht GitHub (Heading-Link muss klein sein, HTML-Anker case-genau); kein Widerspruch zwischen den beiden Teilmengen.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 1 | 2 |

## Verdikt

**MEDIUM blockiert typischerweise** — hier ist „Merge" der Übergang in die
Implementierung. Empfehlung: Finding 1 und 2 als **Review R1 am Change
Request** (reiner Doc-Fix) vor dem Code lösen:

- **F1** — §`DC-FA-ANCH-001.b` um die Inline-Code-Span-Behandlung
  ergänzen (GitHub rendert nur Block-/Fließtext-HTML als Anker) und das
  HTML-Negativkriterium um den Inline-Code-Fall erweitern. Die
  Erkennungs-Infrastruktur existiert bereits (`stripInlineCodeByLine` in
  `internal/hexagon/core/markdown.go`), ist also kostenneutral.
- **F2** — den Lastenheft-Wortlaut `DC-FA-CODE-001` an die gemeinsame
  Anker-Menge angleichen (statt „gegen die Headings"), damit Rang 1 die
  Rang-2-Schärfung deckt.

**F3** (LOW) läuft idealerweise in derselben R1 mit (Querverweis §5 →
`DC-FA-ANCH-001`). **F4/F5** (INFO) sind Implementierungs- bzw.
Designnotizen, kein Blocker. **Kein HIGH:** keine Harness-Lüge, kein
Korrektheits- oder Import-Verstoß; `make gates` ist auf dem Review-Stand
grün.

## Disposition (Review R1 — 2026-06-15)

- **F1 — gefixt:** §`DC-FA-ANCH-001.b` schließt Inline-Code-Spans zusätzlich zu Fenced-Blöcken aus (GitHub rendert HTML in Code-Auszeichnung nicht); das HTML-Negativkriterium in `DC-FA-ANCH-001` deckt den Inline-Code-Fall mit ab.
- **F2 — gefixt:** Lastenheft `DC-FA-CODE-001` prüft den Anker jetzt „gegen die gültige Anker-Menge (Heading-Slugs und Inline-HTML-Anker)" — Rang 1 deckt die Rang-2-Schärfung.
- **F3 — gefixt:** §5 grenzt „Nicht-Markdown-Formate als eigenständige Dateien" ab und verweist für Inline-HTML-Anker auf `DC-FA-ANCH-001`.
- **F4 — gefixt:** §`DC-FA-ANCH-001.b` verlangt den exakten `<a>`-Tag-Vergleich (Wortgrenze, kein Präfix-Treffer auf `<area>`/`<abbr>`).
- **F5 — akzeptiert:** durch F1 deutlich entschärft (Code-Auszeichnung ausgeschlossen); der verbleibende Prosa-Restfall bleibt bewusste Won't-Fix-Designnotiz (Schweige-Charakter des Moduls).

Verifikation: erneuter `make gates`-Lauf nach R1 (doc-check über alle
geänderten Stellen); kein offenes MEDIUM/HIGH → Übergang in die
Implementierung frei.
