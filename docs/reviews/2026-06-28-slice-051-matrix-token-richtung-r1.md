# Review-Report — slice-051 (matrix: token-basierte Referenz-Richtung) — R1

| Feld | Wert |
|---|---|
| Datum | 2026-06-28 |
| Gegenstand | Uncommittete Änderungen slice-051 — Modul `matrix`, `DC-FA-MTX-003` (Token-Referenz-Richtung, Provenance-Marker, `exempt-paths`-Grandfathering) |
| Reviewer | unabhängiger Reviewer (Skill 1.2.0) |
| Eingang | `git diff` + untracked (ADR-0022, slice-051); Lastenheft `DC-FA-MTX-003`; ADR-0022; Spezifikation §DC-FA-MTX-001.a Schritt 6 + §2-Schema + §4; Regelwerk §Referenz-Richtung (SDP) |
| Gate-Status laut Auftrag | `make gates` grün; Negativ-Probe (unmarkierter `slice-099` im nicht-grandfatherten ADR ⇒ `matrix-forbidden`) verifiziert |
| Code im Fokus | `internal/hexagon/core/rules/matrix.go`, `…/model/config.go`, `…/configyaml/configyaml.go`, `…/app/suggest.go`, `…/cli/config_template.go`, Tests `matrix_test.go` / `configyaml_test.go`, Dogfood `.d-check.yml` + Spec/Doku |

Hinweis zur Lösungsangabe: das `befund`-Feld bleibt beobachtend (Skill-Regel
„kein Lösungsvorschlag im Befund"); Fix-Hinweise stehen separat unter
„Übergabe".

---

## Findings

### F-1 — Nur der erste verbotene Token je Zeile/Klasse wird gemeldet
- **Kategorie:** LOW
- **Quelle:** Maintainability / DC-FA-MTX-003
- **Pfad:** `internal/hexagon/core/rules/matrix.go:112`
- **Befund:** `tokenFindings` nutzt `c.Token.FindStringIndex(stripped)` (erster
  Treffer) statt `FindAllStringIndex`. Eine Prosa-Zeile mit zwei verbotenen
  Token derselben Klasse (z. B. `… begründet durch slice-042 und slice-099.`)
  erzeugt nur **einen** `matrix-forbidden`-Befund; der zweite Token wird
  verschwiegen. Kein False-Green (das Gate bleibt rot, solange der erste Token
  steht), aber unvollständige Diagnose: nach dem Fix des ersten Tokens zündet
  der zweite erst im Folgelauf.
- **Verifizierbar:** ja — Unit-Test mit einer Zeile, die zwei verbotene Token
  trägt, erwartet zwei Befunde.
- **Übergabe:** `FindAllStringIndex` iterieren oder pro Treffer einen Befund
  emittieren.

### F-2 — Provenance-Marker ist ein nackter Substring auf der rohen Zeile
- **Kategorie:** LOW
- **Quelle:** DC-FA-MTX-003 / Maintainability
- **Pfad:** `internal/hexagon/core/rules/matrix.go:101` (`strings.Contains(pl.raw, provenanceMarker)`)
- **Befund:** Zwei gekoppelte Effekte. (a) Ein Marker nimmt **alle** Token der
  Zeile aus, nicht nur die deklarierte Referenz — ein legitim markierter
  Provenance-Token kann auf derselben Zeile einen getarnten Entscheidungs-Token
  mit-tarnen. (b) Die Erkennung ist ein reiner Substring-Test auf `pl.raw`:
  taucht der Literalstring `d-check:status-provenance` in **Prosa oder
  Inline-Code** auf (nicht als HTML-Kommentar — etwa eine Zeile, die den Marker
  erklärt), wird die ganze Zeile von der Token-Prüfung ausgenommen. Aktuell
  keine schädliche Ko-Vorkommnis im Dogfood (geprüft: ADR-0022 Z.42 nennt den
  Marker-Namen, trägt aber keinen unmarkierten Token).
- **Verifizierbar:** ja — Fixture mit einer Zeile, die den Marker-Namen als
  Fließtext **und** einen separaten unmarkierten verbotenen Token enthält;
  erwartet einen Befund.
- **Übergabe:** Effekt (a)/(b) bewusst akzeptieren (Marker-Ehrlichkeit ist laut
  Skill Reviewer-Backstop) oder die Erkennung auf den HTML-Kommentar-Token
  verengen; Designentscheid in ADR-0022 festhalten.

### F-3 — `linkSpanRe` strippt Badge-/verschachtelte Links nur partiell
- **Kategorie:** LOW
- **Quelle:** DC-FA-MTX-003
- **Pfad:** `internal/hexagon/core/rules/matrix.go:25` + `:104`
- **Befund:** `linkSpanRe = \[[^\]]*\]\([^)]*\)` matcht beim Badge-Muster
  `[![alt](img)](ziel)` nur `[![alt](img)` und lässt den Rest `](ziel)` stehen.
  Enthält `ziel` einen Slice-Token (z. B. `](../planning/done/slice-045-x.md)`),
  findet der Token-Scan dort einen Treffer und erzeugt — zusätzlich zum
  Link-Befund desselben Links — einen **Doppelbefund** in Token-Form (anderes
  `Target`, daher nicht von `SortFindings` dedupliziert). Einfache Links
  `[text](ziel)` werden vollständig entfernt; das Risiko betrifft nur
  verschachtelte/Badge-Links und Ziele mit `)` im Pfad. Aktuell kein solcher
  Link in den geprüften matrix-Klassen-Dateien.
- **Verifizierbar:** ja — Unit-Test mit einem Badge-Link auf eine `slice-NNN`-Datei
  in einem nicht-grandfatherten ADR; erwartet genau einen Befund.
- **Übergabe:** Link-Spans über den vorverarbeiteten Link-Extraktor entfernen
  (statt eines flachen Regex) oder Badge-/verschachtelte Spans rekursiv strippen.

### F-4 — Spec'd Carve-outs „Token in Fence/exclude-section zählt nicht" im Token-Pfad ungetestet
- **Kategorie:** MEDIUM
- **Quelle:** DC-FA-MTX-003 (fehlende Negativtests bei neuem öffentlichem Vertrag)
- **Pfad:** `internal/hexagon/core/rules/matrix_test.go:356` (`TestMatrixTokenReferenz`, 356–394)
- **Befund:** Der Test pinnt Boundary (unmarkiert → Befund), Happy (Marker → frei),
  Grandfathering (`exempt-paths` → frei), Link-Kein-Doppel und Self-Class. Die
  ebenfalls in Spezifikation §DC-FA-MTX-001.a Schritt 6 und im Slice-DoD genannten
  Token-Carve-outs **Fenced-Code** und **`exclude-sections`** werden über den
  Token-Pfad von keinem Fixture ausgeübt (kein Fixture enthält einen Slice-Token
  in einer Fence oder in einer Historie-/Geschichte-Sektion). Ein künftiger Umbau
  von `tokenFindings` (Iteration über rohe statt `proseLines`, oder Wegfall des
  `inRanges(excluded, …)`-Checks) bräche die Ausnahmen, ohne dass ein Test rot wird.
  Residualrisiko aktuell niedrig, weil `proseLines`/`inRanges` geteilte, im
  Link-Pfad getestete Helfer sind.
- **Verifizierbar:** ja — Mutationsprobe (Token-Scan auf rohe Zeilen umstellen)
  bleibt grün; fehlt das Regressions-Fixture.
- **Übergabe:** Zwei Fixtures zu `TestMatrixTokenReferenz` ergänzen — ein
  Slice-Token in einer ```` ``` ````-Fence und einer unter `## Historie`/`## Geschichte`,
  jeweils erwartet kein Befund.

### F-5 — Token-Scan strippt keinen Inline-Code (Divergenz zu Link-Scan und `ids`)
- **Kategorie:** LOW
- **Quelle:** DC-FA-MTX-003
- **Pfad:** `internal/hexagon/core/rules/matrix.go:100`–`:104`
- **Befund:** Der Token-Scan arbeitet auf `pl.raw` und entfernt nur Fences
  (via `proseLines`) und Link-Spans (via `linkSpanRe`) — **nicht** Inline-Code.
  Ein Slice-Token in Inline-Code (z. B. `` `docs/plan/planning/done/slice-042-x.md` ``)
  im Körper eines ADR/Spec-Dokuments außerhalb der Historie erzeugt
  `matrix-forbidden` und verlangt einen Marker. Das ist spec-treu (Schritt 6
  nennt als Ausnahmen nur Fences/`exclude-sections`/Links, nicht Inline-Code),
  divergiert aber vom Link-Pfad (`PreprocessMarkdown` leert Inline-Code) und vom
  Modul `ids`. Latente False-Positive-Fläche; aktuell kein Dogfood-Vorkommen
  (geprüft: alle Spec-/ADR-Body-Slice-Token liegen in der Historie oder tragen
  den Marker).
- **Verifizierbar:** ja — Fixture mit einem backtick-umschlossenen Slice-Pfad in
  einem nicht-grandfatherten ADR-Körper; beobachtet einen `matrix-forbidden`-Befund.
- **Übergabe:** entweder Inline-Code vor dem Token-Scan leeren (Parität zum
  Link-Pfad) oder die Spec-Ausnahmeliste um Inline-Code ergänzen, je nach
  gewünschter Semantik.

### F-6 — `exempt-paths` überspringt die Datei komplett, nicht nur die Token-Prüfung
- **Kategorie:** INFO
- **Quelle:** ADR-0022 / DC-FA-MTX-003
- **Pfad:** `internal/hexagon/core/rules/matrix.go:34`
- **Befund:** Der `exempt-paths`-Guard steht ganz oben in `CheckMatrix` und
  überspringt die Datei vollständig — auch die Link-`matrix-forbidden`-, die
  `matrix-inactive`- und die `matrix-downward`-Prüfung, nicht nur die
  Token-Prüfung, die das Grandfathering motiviert. Für die 21 immutablen Alt-ADRs
  zurzeit verlustfrei (sie waren grün), und ADR-0022 deklariert „ganz
  übersprungen" bewusst. Folge: eine künftig in einem grandfatherten Alt-ADR
  auftauchende `matrix-inactive`-Verletzung bliebe ungeprüft — praktisch durch die
  `adr-check`-Immutabilität ausgeschlossen.
- **Verifizierbar:** ja — Fixture: exempt-ADR mit Link auf ein `superseded`-Dokument
  erzeugt kein `matrix-inactive`.
- **Übergabe:** als bewusste, dokumentierte Entscheidung belassen (kein Eingriff
  nötig); ggf. in ADR-0022 die Breite explizit als gewollt benennen.

---

## Negativbefunde (geprüft, ohne Befund)

- **Fail-closed-Config:** `compileMatrixToken` lehnt nicht kompilierende Regexe
  und Leerstring-Matcher ab (Exit 2); beide durch `TestDecode_MatrixTokenFailClosed`
  (`[unclosed`, `x*`) abgedeckt. Für jede Klasse aufgerufen. Ohne Befund.
- **Determinismus (DC-QA-02):** `SortFindings` sortiert total (File, Line, Rule,
  Target, Reason) + dedupliziert; `tokenFindings` iteriert die Slices `cfg.Classes`
  und `proseLines` (keine Map). Default-aus byte-identisch: `Token == nil` ⇒
  `tokenFindings` no-op, `ignored(file, nil) == false`, `lineageValues` unverändert.
  Ohne Befund.
- **Grandfathering-Globs:** `00[01][0-9]-*.md` + `002[01]-*.md` decken konkret
  0001–0021 und schließen 0022+ aus (per-Segment `path.Match`-Probe verifiziert);
  ADR-0022 ist die einzige zurzeit nicht-grandfatherte ADR-Datei. Ohne Befund.
- **Self-Class- und Regel-Gating:** Token-Scan überspringt die Quell-Klasse
  (`c.Name == srcClass`) und Klassen ohne verbotene Regel (`!found || rule.Allow`);
  durch `slice-045`/`slice-046`-Fixture (slice→slice kein Befund) belegt. Ohne Befund.
- **exclude-sections der Spec-Historie:** Beide Spec-Straten führen die
  Versions-/Änderungstabelle unter `## 7. Historie`, das mit
  `exclude-sections: ["7. Historie", …]` exakt matcht; alle Spec-Body-Slice-Token
  liegen ausschließlich dort → keine False-Positives. Ohne Befund.
- **ADR-0022-Dogfood + Richtung:** der `slice-051`-Beleg (Z.77) trägt den Marker
  auf der rohen Zeile, der Historie-Eintrag (Z.83) liegt unter `## Geschichte`
  (excluded), keine weiteren Body-Token. `Bezug:`/`Schärft:` verweisen aufwärts
  (Lastenheft/Spezifikation) bzw. seitlich (ADR-0021/MR-006/Regelwerk) — kein
  Abwärtsverweis. Ohne Befund.
- **suggest.go / config_template.go:** der emittierte `token:`-/`adr→slice`-/
  `# exempt-paths`-Block dekodiert (Round-Trip `configyaml.Decode` in
  `cli_acceptance_test.go`, ai-harness + ai-harness-init) und ist 10× byte-identisch;
  `slice-\d{3}` kompiliert und matcht den Leerstring nicht. Ohne Befund.
- **Spec↔Code-Treue:** `DC-FA-MTX-003` + §DC-FA-MTX-001.a Schritt 6 + §2-Schema
  (`matrix.classes[].token`, `matrix.exempt-paths`) + §4 (`matrix-forbidden`
  Token-Form) sind im Code vollständig kodiert; `matrix-forbidden` wiederverwendet
  (kein neuer Grund-Code). AC-Trio (Happy/Boundary/Negative-Grandfathering)
  getestet. Ohne Befund.
- **Reviewer-Skill-Abspeckung:** Version 1.1.0→1.2.0, MEDIUM-Anker von
  „Referenz-Richtung beurteilen" auf „Marker-Ehrlichkeit" umgestellt, Link auf
  `DC-FA-MTX-003` (Slug korrekt). Ohne Befund.
- **Lint/Komplexität:** `lineageValues`/`tokenFindings`-Extraktion sauber,
  `CheckMatrix` flach; `regexp`-Import in `rules` bereits etabliert, keine neue
  verbotene Import-Kante. Ohne Befund.
- **Konsistenz Version/Changelog/Roadmap/Index:** 0.31.0 in Lastenheft, CHANGELOG,
  ADR-0022 (Accepted) + README-Indexzeile, Roadmap-Flip auf aktive welle-40 —
  konsistent. Ohne Befund.

---

## Kategorie-Summary

| Severity | Anzahl |
|---|---|
| BLOCKER | 0 |
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 1 |

---

## Verdikt

**Mergebar: ja (mit einer empfohlenen Test-Ergänzung).**

Es gibt keinen BLOCKER/HIGH: das Feature ist additiv und default-aus byte-identisch,
fail-closed bei Fehlkonfiguration, deterministisch, dogfood-sauber (Spec-Historie
excluded, ADR-0022 markiert/grandfathered), und der Negativ-Probe-Pfad feuert
korrekt. Kein Befund ist ein Stilles-Grün- oder Exit-Code-Fehler.

Das einzige MEDIUM (F-4, fehlende Fence-/exclude-section-Regressionstests im
Token-Pfad) blockiert nach Skill-Regel „typischerweise", wird hier aber als
**nicht-blockierend** eingestuft: das Verhalten ist über die geteilten,
anderswo getesteten Helfer (`proseLines`/`inRanges`) korrekt; F-4 schließt nur
eine Regressions-Lücke. Empfehlung: die zwei Fixtures vor dem Release nachziehen.
Die LOW-Befunde (F-1 Erst-Treffer-only, F-2 Marker-Substring, F-3 Badge-Link-Leak,
F-5 Inline-Code) sind latente Kanten ohne aktuellen Dogfood-Trigger und können als
Folge-CR adressiert werden.
