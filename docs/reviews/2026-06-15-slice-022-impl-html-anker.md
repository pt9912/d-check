# Review-Report: slice-022 Implementierung — Inline-HTML-Anker — 2026-06-15

**Review-Art:** Code-Review der Implementierung (zweiter Lauf nach dem
Doc/CR-Review). Gegenstand ist der Code, nicht erneut der Vertrag.

**Gegenstand:** Commit `6d402b8` (`feat(anchors)`):
`internal/hexagon/core/anchors.go` (+67), `anchors_test.go` (+78),
`internal/hexagon/core/codepaths.go` (±8).

**Skill:** `.harness/skills/reviewer.md` · **Datum:** 2026-06-15

**Eingangs-Kontext:** `DC-FA-ANCH-001` / §`DC-FA-ANCH-001.b` /
`DC-FA-CODE-001`, `DC-QA-01`/`DC-QA-02`, `AGENTS.md` §3, `ADR-0005`,
slice-022, Doc-Review-Report (R1-Dispositionen).

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar |
|---|---|---|---|---|---|
| 1 | 🟡 LOW | `DC-FA-ANCH-001` / Robustheit (Falsch-Positiv-Richtung) | `internal/hexagon/core/anchors.go:99` | `htmlTagRE` begrenzt den Tag mit `[^>]*`; ein rohes `>` in einem früheren Attributwert (`<a title="x > y" name="z">`) beendet die Erfassung vorzeitig, ein danach stehendes id/name wird übersehen. Ein Link darauf meldet dann fälschlich `anchor-missing` — das ist die „Wolf-rufen"-Richtung im Gate-Pfad (kann eine grüne doc-check-CI brechen). Eintritt nur bei literalem `>` im Attribut (regulär `&gt;`). | ja — Unit-Test mit `<a title="a > b" name="z">` + Link `#z`; erwartet kein Befund, aktuell anchor-missing |
| 2 | 🟡 LOW | `DC-FA-ANCH-001` / Testabdeckung | `internal/hexagon/core/anchors_test.go` | Der Vertragsfall „`#anker` innerhalb derselben Datei" ist für HTML-Anker nicht direkt getestet; `TestAnchorsHTMLModul` prüft nur Cross-File-Ziele. Der Selbst-Datei-Pfad (`slugsFor` mit own-Content) teilt die Extraktion, ist aber für HTML-Anker ungeprüft. | ja — Test mit `<a name="x">` und Link `#x` in derselben Datei |
| 3 | 🟡 LOW | `DC-FA-CODE-001` / Testabdeckung | `internal/hexagon/core/codepaths.go:174` | Die im Doc-Review verlangte HTML-Anker-Konsistenz für `codepaths` ist verdrahtet (`codepathSlugs` → `AnchorSet`), aber kein `codepaths`-Test prüft einen Inline-Code-Pfad mit Fragment gegen einen HTML-Anker des Ziels. | ja — codepaths-Test mit Inline-Code-Pfad `ziel.md#html-id` und `<div id="html-id">` in ziel.md |
| 4 | 🔵 INFO | `DC-QA-01` / Maintainability (Perf) | `internal/hexagon/core/anchors.go:146` | `AnchorSet` ruft `HeadingSlugs` (extractHeadingLines-Scan) **und** `htmlAnchors` (PreprocessMarkdown-Scan) — zwei volle Zeilen-Durchläufe je Zieldatei plus Regex-Allokation je Prosa-Zeile. Pro distinktem Ziel via `slugCache` nur einmal (gedeckelt), bei sehr großen Korpora dennoch zusätzlicher Aufwand. | ja — `make bench` gegen das `DC-QA-01`-Fixture (Median), kein Gate |
| 5 | 🔵 INFO | dokumentationswürdige Annahme | `internal/hexagon/core/anchors.go:132` | `attrValue` unterscheidet „leerer Wert" (`id=""`) nicht von „kein Treffer" (beide → `""`); ein leerer id/name erzeugt bewusst keinen Anker. Sinnvoll, aber als Annahme nicht in der Spezifikation vermerkt. | nein — Designnotiz |

## Negativbefunde (geprüft, ohne Befund)

- **HTML-Kommentare:** `<!-- … -->` (u. a. der `d-check:ignore`-Marker) matcht `htmlTagRE` nicht (`!` ist kein `[a-zA-Z]`) → keine Phantom-Anker; verifiziert.
- **Modul-Konsistenz (anchors ↔ codepaths):** beide nutzen `AnchorSet` → identische gültige Anker-Menge; die im Doc-Review (F2) verlangte Rang-1/Rang-2-Angleichung ist im Code eingelöst.
- **Determinismus (`DC-QA-02`):** Anker-Menge ist Mengen-Mitgliedschaft, Ausgabe-Reihenfolge über `SortFindings`; keine ausgabewirksame Map-Iteration → identische Eingabe ⇒ identische Ausgabe (Determinismus-Test grün).
- **Monotonie/Abwärtskompatibilität:** `AnchorSet` ⊇ `HeadingSlugs` (strikte Obermenge) → die Erweiterung kann Befunde nur entfernen, nie hinzufügen (Ausnahme: Finding 1); kein bestehender anchors-/codepaths-Test kippt (gates grün).
- **Case-Sensitivität:** Heading-Slugs kleingeschrieben, HTML-Anker wörtlich, Fragment unverändert verglichen — GitHub-konform (`TestAnchorsHTMLModul`: `#übersicht` ≠ `id="Übersicht"`).
- **Fence/Inline-Code:** Ausschluss über `PreprocessMarkdown` (etablierter Fence-Automat) — kein zweiter, divergierender Parser (`TestHTMLAnchors` deckt beide Ausschlüsse ab).
- **Tag-Exaktheit:** `EqualFold(tag[1], "a")` vergleicht den vollen Tag-Namen, kein Präfix → `name` an `<area>`/`<abbr>` ergibt keinen Anker (`TestHTMLAnchors` verifiziert).
- **`data-id`-Abgrenzung:** `(?:^|\s)id` verlangt Whitespace/Zeilenanfang vor `id`; in `data-id` steht davor `-` → kein Treffer (`TestHTMLAnchors` verifiziert).
- **ReDoS:** lineare Regexe ohne verschachtelte Quantoren.
- **Arch (`ADR-0005`):** nur `regexp` (stdlib) ergänzt, keine neue Cross-Layer-Abhängigkeit (`make arch-check` grün).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 3 | 2 |

## Verdikt

**Kein HIGH/MEDIUM → blockiert die Closure nicht.** Die drei LOW sind
nice-to-fix: Finding 1 ist eine konservative Regex-Grenze bei
malformiertem HTML (Maintainer-Entscheidung: dokumentieren oder
Mini-Härtung); Finding 2 und 3 sind günstige Test-Ergänzungen, die die
beiden ungeprüften Arme der geteilten Anker-Menge abdecken — sinnvoll
**vor** der Closure, da der Code dann offen ist. Die zwei INFO sind
Notizen. `make gates` ist auf dem Review-Stand grün (Coverage 94,80 %).

## Disposition (Review R1 — 2026-06-15)

- **F1 — gefixt:** `htmlTagRE` ist quote-bewusst (`(?:"[^"]*"|'[^']*'|[^>'"])*`); ein `>` in einem Attributwert beendet den Tag nicht mehr. Test `<a title="x > y" name="gtinattr">` in `TestHTMLAnchors` (scheiterte mit der alten Regex).
- **F2 — gefixt:** `TestAnchorsHTMLSelbeDatei` — HTML-Anker via `#frag` innerhalb derselben Datei plus Negativfall.
- **F3 — gefixt:** `TestCodepathsHTMLAnker` — `codepaths`-Anker-Prüfung gegen eine HTML-`id` der Zieldatei plus Negativfall.
- **F5 — gefixt:** `attrValue`-Kommentar dokumentiert die Annahme (leerer Attributwert ⇒ kein Anker).
- **F4 — akzeptiert:** der Doppel-Scan ist über den `slugCache` gedeckelt (eine `AnchorSet` je distinktem Ziel); kein riskanter Refactor zugunsten einer INFO-Notiz.

`make gates` nach R1 grün (50 Dateien / 0 Befunde, Coverage 94,80 %).
