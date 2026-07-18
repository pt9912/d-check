# Review R1 — slice-078: Geteiltes Referenz-Ventil `ignore-refs` mit Quell-Skopus

**Datum:** 2026-07-18 · **Reviewer:** unabhängig, kontext-getrennt ·
**Gegenstand:** [`slice-078`](../plan/planning/done/slice-078-ignore-refs-quell-skopus.md)
(Commit-Range `5223f55..HEAD` auf `main`, Tip `8f9535d`). **Typ:** Pre-Release-Review
(ein `BLOCK` verhindert das Tag `v0.49.0`).

## Eingangs-Kontext

- **Vertrag:** [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  (Lastenheft, neu; Bereich `REF`), [`DC-FA-REF-001.a`](../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)
  + §2-Schema (Spezifikation), [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Proposed).
- **Berührte Module:** `links`, `anchors`, `codepaths` (honorieren das Ventil),
  Adapter `configyaml` (Validierung), Kern-Modelle.
- **Prüftiefe:** Diff gelesen; Semantik gegen Lastenheft/Spec/ADR abgeglichen; vier
  gezielte Quell-Mutationen gefahren (`make test`), um die behaupteten Mutations-Pins
  zu verifizieren; Realdatenbeleg gelesen; Release-Prep-Currency geprüft.
- **Gate-Läufe:** `make test` grün (Exit 0); `make doc-check` grün (Whole-Repo,
  inkl. dieses Reports).

## Findings

### F-1 (MEDIUM · Maintainability / `DC-FA-REF-001` / `DC-QA-02`) — Top-Level-`ignore-refs`→`codepaths`-Verdrahtung ist von keinem Test gepinnt

**Pfad:** `internal/hexagon/core/rules/codepaths.go:36`

**Befund:** Die Zeile `refs := ignoreRefs` reicht das geteilte Top-Level-Ventil an
`codepaths` durch (unabhängig vom Alias). Mutiert man sie zu `refs := nil` (nur der
Alias erreicht dann `codepaths`), bleibt die **gesamte** Testsuite grün — kein Test
prüft, dass ein Top-Level-`ignore-refs`-Eintrag einen `codepaths`-Befund
unterdrückt. Der querschnittliche Pin-Test `TestRefsGeteiltesVentilLinksUndAnchors`
deckt nur `links`+`anchors`; `TestRefsAliasKoexistenz` nutzt Top-Level ausschließlich
für ein `links`-Ziel. Damit ist die DoD-Zusage „querschnittliches Wiring per Mutation
gepinnt (je genau der zuständige Test kippt)" für den `codepaths`-Zweig **nicht
erfüllt** — die Verdrahtung ist real-datenbelegt (der Realdatenbeleg schweigt die 5
`codepath-missing` aus `tools/check_*.py` über einen **Top-Level**-`refs`-Eintrag),
aber nicht regressions-gesichert.

**Verifizierbar:** ja — `make test` bleibt nach der Mutation `refs := nil` grün
(vom Reviewer gefahren).

### F-2 (LOW · Maintainability) — Alias-Globs `codepaths.ignore-refs` werden config-zeitig nicht validiert (fail-open ggü. fail-closed des neuen Ventils)

**Pfad:** `internal/adapter/driven/configyaml/configyaml.go:948` (`applyCodepaths`)

**Befund:** Das neue Top-Level-`ignore-refs` validiert jedes `in`/`refs`/`keep`-Glob
segmentweise und bricht bei einem ungültigen Muster mit Exit 2 ab (`validRefGlob`,
getestet in `TestRefs_UngueltigesGlobExit2`). Der **Alias** `codepaths.ignore-refs`
durchläuft diese Validierung **nicht** — `applyCodepaths` prüft nur `roots`, `IgnoreRefs`
wird ungeprüft durchgereicht. Ein ungültiges Alias-Glob (z. B. `["["]`) wird zur
Laufzeit von `matchGlob` still verschluckt (`path.Match`-Fehler ignoriert → matcht
nie), das Ventil wäre wirkungslos ohne Hinweis. Dieselbe Fähigkeit trägt damit zwei
verschiedene Fehlermodi. Der Punkt ist **vorbestehend** und durch die
Byte-Identitäts-Zusage des Alias ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus))
gedeckelt (nachträgliche Validierung wäre selbst eine Verhaltensänderung); daher nur
zur Kenntnis, kein durch diesen Slice eingeführter Defekt.

**Verifizierbar:** ja — eine `.d-check.yml` mit `codepaths.ignore-refs: ["["]` läuft
ohne Exit 2 durch (im Gegensatz zum Top-Level-Pendant).

### F-3 (LOW · `DC-FA-LINK-002` / ADR-0044 Entscheidung 7) — „Symlink bleibt" ist korrekt, aber ungetestet

**Pfad:** `internal/hexagon/core/rules/links.go:27`–`42`

**Befund:** Die Spezifikation und [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
(Entscheidung 7) fordern, dass die Symlink-Ablehnung trotz `ignore-refs`-Treffer
greift. Der Code stellt das strukturell sicher — die Symlink-Prüfung steht in
`CheckLinks` **vor** `refIgnored`, mit `continue` bei Treffer. Es gibt jedoch keinen
Test, der ein ignoriertes **und** per Symlink aufgelöstes Ziel prüft (weiterhin
`symlink`-Befund). In `anchors`/`codepaths` ist die Kante gegenstandslos (kein
Symlink-Check dort), also betrifft es nur `links`, wo die Reihenfolge die Invariante
garantiert. Nice-to-fix-Pin, keine Verhaltenslücke.

**Verifizierbar:** ja — ein `anchors`-freier Testfall mit Symlink-Ziel + passendem
`refs`-Glob würde die Invariante festnageln.

### F-4 (LOW/INFO · Maintainability) — DoD-Häkchen „Nutzerdoku" nicht nachgezogen

**Pfad:** `docs/plan/planning/done/slice-078-ignore-refs-quell-skopus.md:200`

**Befund:** Das DoD-Item „Nutzerdoku (Handbuch §5/§6 Ventil-Achsen) + CHANGELOG" ist
inhaltlich erledigt (Release-Prep-Commit `8f9535d`: Handbuch §5 vier-Achsen-Erklärung,
CHANGELOG, version.md), das Häkchen steht aber noch auf `[ ]`. Reine Doc-Drift; die
verbleibenden offenen Punkte (Nutzerdoku, Qualitäts-Review) werden regulär bei der
Closure abgehakt.

**Verifizierbar:** nein (Doku-Konsistenz, nicht gate-gebunden).

## Negativbefunde (geprüft, ohne Befund)

- **Semantik `refIgnored`** (`paths.go:140`): das Prädikat
  `(In=="" ∨ matchGlob(In,file)) ∧ ignored(rel,Refs) ∧ ¬ignored(rel,Keep)` ist
  deckungsgleich mit dem Spec-Match-Prädikat, dem Lastenheft-`refs ∧ ¬keep` und
  ADR-0044 Entscheidung 2/4. `keep` wirkt korrekt **pro Eintrag** (Spec: „desselben
  Eintrags"); mehrere Einträge additiv (Union). Reihenfolge-Unabhängigkeit gegeben.
- **Mutations-Pins `keep`/`in`/`links`+`anchors` sind echt** (nicht transitiv/blind):
  vom Reviewer verifiziert — `!ignored(rel,Keep)` entfernen kippt genau
  `TestRefsKeepGewinnt` + `TestRefsAnkerBleibtScharf` + `TestRefs_KeepUndTippfehlerEndToEnd`;
  den `In`-Skopus deaktivieren kippt genau `TestRefsSkopusIsolation`; die
  `anchors`-Verdrahtung deaktivieren kippt genau `TestRefsGeteiltesVentilLinksUndAnchors`.
- **Alias-Byte-Identität** (`codepaths.go:36`–`39`): für Alias-only reduziert sich
  `refIgnored` auf das alte `ignored(rel, cfg.IgnoreRefs)` (In leer, Keep leer); die
  unveränderten `TestCodepathsIgnoreRefs`/`…UnterdruecktEscapeUndAnker` bleiben grün.
  Der Alias ist auf `codepaths` skopiert und schlägt nicht auf `links` durch
  (`TestRefsAliasKoexistenz`-Gegenprobe).
- **Achsen-Präzedenz** stimmt mit der Spec: `scan.ignore` prunt vor jedem Modul, der
  Zeilen-Marker wirkt in `codepaths` Schritt 1 vor der Pfad-Erkennung, `refIgnored`
  greift bei der Ziel-Auflösung (vor Escape/Existenz/Anker — `links.go:40`,
  `codepaths.go:155`). Keine Achse überschreibt eine andere.
- **Fail-closed** (`configyaml.go` `applyIgnoreRefs`/`validIgnoreRefEntry`/`validRefGlob`):
  segmentweise Glob-Validierung wie zur Laufzeit; leeres `in` = repo-weit (erlaubt),
  leeres `refs`/`keep`-Element = Exit 2, leere `refs`-Liste = inert. Alle drei Kanten
  von `TestRefs_UngueltigesGlobExit2` abgedeckt.
- **Realdatenbeleg** ([`2026-07-18-…-realdatenbeleg-…`](2026-07-18-slice-078-realdatenbeleg-ai-harness-course.md)):
  Baseline 42 = 37 `target-missing` + 5 `codepath-missing` reproduziert **exakt** die
  CR-Messung (Aufbau treu); Ventil-Lauf 0 Befunde bei 158 Dateien; zwei in echte Ziele
  (gekeepter `.template.md`, Kurs-Verweis) injizierte Tippfehler feuern beide als
  ERROR, kein Platzhalter kommt hinzu — nicht durch Wegschauen. Das ausgelieferte
  v0.48.0-Image lehnt die Ventil-Config ab (Top-Level `ignore-refs` unbekannt ⇒ Exit
  2): Fähigkeit nachweislich neu.
- **Doc↔Code-Konsistenz:** Lastenheft-AKs (6), Spec-Algorithmus, ADR-Entscheidungen
  und Handbuch §5 (vier Achsen gegeneinander erklärt) stimmen überein. Das
  Handbuch-Config-Beispiel (`in: "lab/templates/**"` / `refs` / `keep:
  ["lab/templates/**/*.template.md"]`) besteht die segmentweise Glob-Validierung; das
  Alias-Beispiel `codepaths.ignore-refs: ["tools/altes-skript.sh"]` ebenso.
- **Release-Prep-Currency:** version.md-Anker-Kaskade korrekt (nur `#v0.49.0`
  referenziert, Anker von `v0.48.1` auf `v0.49.0` verschoben — Konvention eingehalten,
  keine hängende Referenz); ghcr-Tag-Pins in `README.md`/`README.de.md`/Handbuch auf
  `v0.49.0` gehoben; CHANGELOG-`[0.49.0]`, Handbuch-Versionshistorie (1.38) und
  Lastenheft-/Spezifikation-Historie konsistent. Lastenheft-Doc-Version 0.46.0→0.47.0
  (unabhängig vom Release-Tag). ADR-Index um ADR-0044 ergänzt.
- **ADR-Immutabilität / Referenz-Richtung:** ADR-0044 ist `Proposed` (Accepted-Übergang
  gehört zur Closure); Spec-Straten verweisen nicht abwärts auf ADRs/Slices; das
  gemeinsame Kürzel-Kriterium ist mit slice-079 widerspruchsfrei verankert.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 2 | F-2, F-3 |
| INFO | 1 | F-4 |

## Verdikt

**ACCEPT-WITH-NITS.** Kern-Semantik, Alias-Byte-Identität, Fail-closed-Validierung,
Achsen-Präzedenz und Realdatenbeleg sind solide und — für `keep`/`in`/`links`+`anchors`
— echt mutations-gepinnt (vom Reviewer nachgefahren). Kein HIGH.

Das einzige MEDIUM (F-1) ist eine **Test-Härte-Lücke**, kein Code-Defekt: die
`codepaths`-Honorierung des Top-Level-Ventils ist korrekt implementiert und durch den
Realdatenbeleg empirisch belegt (5 `codepath-missing` über einen Top-Level-`refs`-Eintrag
stillgelegt), aber von keinem automatisierten Test regressions-gesichert. Ich
**blockiere das Tag deswegen nicht** (dokumentierte Abweichung nach `reviewer.md`:
Verhalten real-datenbelegt + Code korrekt), empfehle aber, vor oder bei der Closure
einen Pin nachzuziehen, der einen Top-Level-`ignore-refs`-Eintrag gegen einen
`codepaths`-Befund verriegelt — dann trägt die DoD-Zusage „querschnittliches Wiring
per Mutation gepinnt" auch für `codepaths`. F-2/F-3 sind Nice-to-fix, F-4 reine
Doc-Drift.
