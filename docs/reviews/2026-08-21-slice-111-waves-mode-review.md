# Review-Report: slice-111 — `planning.waves.mode` (Kennungs-Bijektion als opt-in)

**Datum:** 2026-08-21 · **Review-Art:** Code (geprüft gegen Lastenheft 0.62.0 §Wellen-Invariante, `spec/spezifikation.md` §`DC-FA-PLAN-001.a` W1/W3/§2/§4, ADR-0055 Entscheidung 6, Hard Rules `AGENTS.md` §3, `MR-025`)
**Gegenstand:** slice-111 / welle-79, Commit-Kette `fa3ac07` (CR/Lastenheft) → `fa5c4fe` (ADR-0055-Fortschreibung) → `f4146a3` (feat) → `f06fe8a` (Release-Prep), **vor** dem Release v0.62.0
**Skill:** `reviewer.md` @ 1.5.0 · **Modell-ID:** `claude-fable-5`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-111-wave-drift-zwei-haelften.md`; `DC-FA-PLAN-001`, `DC-QA-02`, `DC-QA-03`; ADR-0054/0055; der Konsumenten-CR liegt vereinbarungsgemäß **nicht** als Repo-Datei vor — seine Landung ist `fa3ac07`. Nicht erhalten: DoD-Abhakung (Verifikations-Rolle).

## Findings

### F-1 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** ADR-0055 / `MR-025` (Spiegel-Wortlaut)
- **pfad:** `docs/plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md:78–80` (dito `docs/plan/planning/in-progress/slice-111-wave-drift-zwei-haelften.md:73–75`)
- **befund:** Entscheidung 6 behauptet: „der Block wird nicht zweimal bestimmt, **die Bestimmung liefert neben dem Bool die Kennungs-Liste**". Der Code tut das nicht: `planningActiveStatus` liefert nur `(headingNo, hasActive, fail)`; die Kennungs-Menge liest ein **separater** Scan (`planningBlockIDs`, aufgerufen aus `waveBijection`), der `splitLines`/`proseLineSet` erneut berechnet und nur die Block-**Grenze** über `planningBlockEnd` teilt. Lastenheft und Spezifikation formulieren korrekt („der Aktiv-Status wird nicht ein zweites Mal bestimmt; beide Fähigkeiten lesen denselben Block" — trifft zu); die ADR beschreibt eine Implementierungs-Form, die nicht existiert. Failure-Szenario: ein Wartender verlässt sich auf die Ein-Lese-Form der ADR und der ID-Scan sieht seine Transformation nicht.
- **verifizierbar:** nein (kein Gate misst ADR-Prosa gegen Code; Beleg ist das Signatur-Zitat)
- **klasse:** adr-code-wortlaut-drift

### F-2 · LOW

- **kategorie:** LOW
- **quelle:** Maintainability (Release-Prep-Currency; `DC-FA-VER-001`-Blindspot)
- **pfad:** `docs/user/benutzerhandbuch.md:75`
- **befund:** Der bare-Tag-Listeneintrag in §„Versionen und Tags" steht noch auf `:v0.61.0`; jede vorherige Release-Prep hat genau diese Zeile mitgehoben. Die Commit-Botschaft von `f06fe8a` zählt „alle 24 bare-Tag-ghcr-Pins" — diese Zeile ist präfixfrei und fällt aus der Zählung wie aus dem `versions`-Gate (das nur ghcr-präfixierte Pins prüft).
- **verifizierbar:** nein durch Gate — genau der bekannte Blindspot; bestätigt per grep (eine Fundstelle, Zeile 75)
- **klasse:** release-prep-doc-currency-blindspot

### F-3 · LOW

- **kategorie:** LOW
- **quelle:** `DC-FA-PLAN-001` (Akzeptanzkriterien-Ränder)
- **pfad:** `internal/hexagon/core/rules/planning_waves_test.go`, `internal/adapter/driven/configyaml/configyaml_test.go`
- **befund:** Vier Rand-Proben des neuen Vertrags fehlen als Tests: (a) explizites `mode: one`; (b) `mode: many` bei fehlender/mehrdeutiger `planning.heading`-Überschrift (fail-Pfad vor der Verzweigung, nur für `one` getestet); (c) W4/W5-Registerzeilen unter `many` mit Inhalt; (d) die `mode`-Validierung bei **inerter** Fähigkeit (Live-Probe: Exit 2, aber ohne Testfall).
- **verifizierbar:** ja — `make test` bleibt grün, die Abwesenheit zeigt die Testliste
- **klasse:** randproben-luecke

### F-4 · LOW

- **kategorie:** LOW
- **quelle:** AGENTS §3.7 (Herkunft nur als **ein** auflösbares Feld)
- **pfad:** `internal/hexagon/core/rules/planning_waves.go`, `internal/hexagon/core/rules/planning_waves_test.go`
- **befund:** Der `waveBijection`-Kommentar (und das Test-Helper-Echo `wavesManyCfg`) trägt neben dem auflösbaren Feld zusätzlich den in-repo **nicht auflösbaren** CR-Titel „Bijektion statt Singleton" und wiederholt ihn über mehrere Kommentare; die Regel verlangt genau ein Feld. Kein HIGH: jede Stelle trägt klar eine Zusage-Klasse.
- **verifizierbar:** nein (kein Kommentar-Gate; Review-Anker)
- **klasse:** kommentar-herkunfts-prosa

### F-5 · INFO

- **kategorie:** INFO
- **quelle:** `DC-FA-PLAN-001` (Wortlaut „Menge der flachen Wellendokumente")
- **pfad:** `internal/hexagon/core/rules/planning_waves.go` (`waveSets`, Vorbestand), `spec/lastenheft.md` §Wellen-Invariante Zeile 1/2
- **befund:** `flach` ist eine **Kennungs**-Map: zwei Dateien derselben Kennung kollabieren zu einem Element — unter `many` genügt ein Block-Zeiger für beide, unter `one` gälten zwei solche Dateien als „genau eines". Vorbestand seit v0.59.0, in beiden Modi kohärent; nur liest sich der Lastenheft-Wortlaut als Datei-Menge, implementiert ist die Kennungs-Menge. Dokumentationswürdige Annahme, kein Verhaltens-Defekt.
- **verifizierbar:** ja — konstruierbares Fixture mit zwei gleich-kennigen Dateien
- **klasse:** mengen-vs-datei-wortlaut

## Negativ-Proben (geprüft, ohne Befund)

Alle Kommandos selbst gelaufen; Gate-Ausgaben in Datei umgeleitet, Exit explizit geprüft (BEO-007-Arbeitsregel).

1. **Gates:** `make test` → **Exit 0**; `make gates` → **Exit 0** (doc-check 411 Dateien / 0 Befunde, lint, test, arch-check, **Coverage 94,50 % ≥ 93 %**, semgrep, gate-consistency, planning-check); `make build` → Exit 0.
2. **Commit-Schnitt:** `git show --stat` aller vier Commits — `fa3ac07` ändert **ausschließlich** `spec/lastenheft.md` (27+/2−); `fa5c4fe` ausschließlich ADR-0055 (58+); `f4146a3` genau die acht deklarierten Code-/Spec-/Test-Dateien (Feat-Commit fasst README/Handbuch nicht an); `f06fe8a` genau die fünf Release-Prep-Dateien.
3. **Default-Treue byte-identisch (DC-QA-02):** der `waveDrift`-Körper ist unangetastet; `planningBlockHasMarker`-Refactor semantikerhaltend; **einziger** Mode-Lesepunkt im Kern via `EffectiveMode()`. **Live-Beweis:** gepinntes Alt-Image (v0.61.0-Digest) gegen HEAD-Build auf demselben Fixture ohne `mode`-Schlüssel — `cmp` **byte-identisch** im grünen (Exit 0/0) wie im roten Szenario (Exit 1/1, `wave-drift`-Zeile identisch).
4. **fail-closed-Wächter vor der Verzweigung:** `dirErrs`-Guard liegt **vor** dem Modus-Branch; Live-Probe: fehlendes `done`-Verzeichnis meldet unter `many` `wave-drift` (Exit 1); `TestWavesManyFailClosedBleibt` deckt `waves.dir`.
5. **Bijektions-Mechanik:** Block-Grenze ist **eine** Antwort (`planningBlockEnd`; Marker-Schleife und ID-Scan identisch geklammert — ADR-0054 gehalten); die Sektions-Überschriftszeile selbst ist ausgenommen; `regexp.QuoteMeta(prefix)` neutralisiert Glob-Metazeichen (Register-Konsistenz zu `waveRowsFrom`); leeres Präfix Exit-2-gesperrt; `\d+` maximal-munch — `welle-790` wird nie als `welle-79` gelesen; `hasActive` wird im `many`-Pfad nicht gelesen; Register-Kennungen liegen hinter der nächsten H2 und zählen nicht als Block-Zeiger.
6. **Dedup:** `SortFindings` dedupliziert über (Datei, Zeile, Regel, **Ziel**, Grund); die beiden Richtungen sind mengendisjunkt mit Kennungs-targets. **Live:** beide Richtungen an derselben Zeile → zwei `wave-drift`, targets `welle-1`/`welle-2`, beide ausgegeben, Exit 1.
7. **Akzeptanzkriterien:** Happy `many` (Test + Live) · Negative beide Richtungen mit target-Assertion (Test + Live) · Marker-Orthogonalität samt `one`-Abgrenzung (Test + Live) · fail-closed (`mode: viele` und `mode: ''` → **Exit 2**, Fehlertext nennt `planning.waves.mode` wörtlich) · Default (Kontrast-Tests + Byte-Probe 3). Fence-/Mehrfach-/Null-/Mehr-Wellen-Proben als Tests vorhanden.
8. **Config-Rand-Konsistenz:** `mode` wird auch bei **inerter** Fähigkeit validiert (Live: `waves:` nur mit `mode: viele`, ohne `dir` → Exit 2); `many` durchgereicht, Abwesenheit ⇒ Default `one`; print-config-Template dekodiert weiterhin.
9. **Release-Prep:** Handbuch-§11-Zeile 1.53 chronologisch **unter** 1.52; `version.md` v0.62.0-Zeile mit Anker; **24** Handbuch- + je **1** README-ghcr-Pin auf v0.62.0, kein ghcr-`v0.61.0`-Rest außerhalb Historie/`done/`; README DE/EN synchron; CHANGELOG deckt den Diff; die Handbuch-**Digest**-Zeile bleibt vereinbarungsgemäß auf dem v0.61.0-Digest (Backfill nach Tag, deklariert). Einzige Ausnahme: F-2.
10. **Sequenzierung der offen bleibenden Flächen:** Profil-Umstellung auf `many`, Grenz-Kommentare, AGENTS-§3.3-Hinweis — ausdrücklich **nach** Release + Digest-Backfill deklariert (Slice §2 Schritt 5, §5-Risiko, ADR-Konsequenz); Profile heute unverändert. Kein Befund.
11. **ADR-Disziplin:** ADR-0055 ist `Proposed` — Fortschreibung zulässig; `## Geschichte`-Eintrag vorhanden; kein zweiter Grund-Code (keine neue `ReasonWave*`-Konstante).
12. **Spec-Straten (§3.4):** keine Abwärts-Referenzen; `matrix` grün.
13. **DC-QA-03:** nur Roadmap + Listing über den bestehenden Port; kein Netz, kein Schreiben.
14. **Kommentar-Klassen:** alle neuen Stellen tragen eine der fünf Klassen; einziger Rest ist der Herkunfts-Zusatz (F-4).

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 (ADR-Wortlaut ≠ Code-Form) |
| LOW | 3 | F-2 (Handbuch-bare-Tag), F-3 (vier Rand-Proben fehlen), F-4 (Kommentar-Herkunfts-Zusatz) |
| INFO | 1 | F-5 (Kennungs- vs. Datei-Mengen-Wortlaut, Vorbestand) |

## Verdikt

**APPROVE mit Auflagen.**

Der Kern des Slice ist sauber: das Vertrags-Viereck stimmt in der Verhaltens-Semantik auf allen Flächen überein, der `one`-Pfad ist nachweislich unangetastet (byte-identische Live-Gegenprobe gegen das gepinnte v0.61.0-Image, grün **und** rot), die Bijektion trägt die zugesagten Eigenschaften (Kennung als target, Marker-Orthogonalität, fail-closed vor der Verzweigung, Dedup-Trennung über das Ziel), der Config-Rand ist fail-closed inklusive Inert-Fall, und die Commit-Schnitte tragen exakt ihre Botschaften. Kein HIGH.

Auflagen vor dem Tag v0.62.0: **F-1** beheben (ein Satz in ADR-0055 Entscheidung 6 — die Kennungs-Liste kommt aus einem eigenen Scan über die **geteilte Block-Grenze**, nicht aus dem Rückgabewert der Aktiv-Status-Bestimmung; ohne Korrektur friert der falsche Wortlaut mit einer späteren `Accepted`-Setzung ein) und **F-2** beheben (`docs/user/benutzerhandbuch.md:75` auf `:v0.62.0`). F-3/F-4/F-5 sind nicht release-blockierend: F-3 als Nachzug-Kandidat, F-4 ggf. als Regel-Schärfung statt Edit, F-5 als Präzisierungs-Kandidat für eine spätere Lastenheft-Redaktion.
