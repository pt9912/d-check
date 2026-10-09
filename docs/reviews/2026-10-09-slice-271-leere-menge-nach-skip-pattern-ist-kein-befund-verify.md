# Verifikation — slice-271: Eine Menge, die erst `skip-pattern` leert, ist kein Befund

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verification.md`): DoD und Plan gegen den tatsächlichen Stand, nicht Diff gegen Plan
**Gegenstand:** slice-271 auf `2382b83d` (Code-Stand `d94bf010`, Spec-Nachzug `2382b83d`); DoD-Punkte 1–3
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan §1, §2, §6 (Fassung nach der Plan-Änderung `636fc084`);
[Befund von `ai-harness-course`](../plan/cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md), Punkt 1;
Reviews R1 und R2 zu slice-271 (nur als Hinweis auf Messbäume, nicht als Beleg übernommen);
[ADR-0106](../plan/adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md);
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).

Alle Läufe in einem frischen Klon bzw. in Wegwerf-Bäumen unter dem Session-Scratchpad, nur über
`make`/Docker. Am Arbeitsbaum ist nichts geändert außer dieser Datei.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| V1 | `git clone` des Arbeits-Repos (HEAD `2382b83d`), `make build IMAGE=d-check-v271` | Image `sha256:73f12a08…` — dieselbe ID wie `d-check:latest`, damit entspricht `latest` HEAD |
| V2 | `docker pull ghcr.io/pt9912/d-check:v0.85.0` | gezogen; Vergleichs-Image für die Vorfassung |
| V3 | `make gates` im Klon | Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; `ok …/core/rules`, `ok …/configyaml`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `0 findings` |
| V4 | Wegwerf-Bäume (siehe unten), je `docker run --rm --network none -v <baum>:/repo:ro <image>`, HEAD und `v0.85.0` | Tabelle *Befund-Fälle* |
| V5 | dieselben Bäume **ohne** Schlüssel, stdout, stderr, Exit-Code und `--json` getrennt, HEAD gegen `v0.85.0` per `cmp` | sechs Bäume, alle vier Kanäle je `ident`; Umzugs-Baum stdout beider Images `sha256:c62a3e6c…` |
| V6 | Mutationen im Klon, je `make test IMAGE=d-check-vmut`, danach `git checkout -- internal/` | Tabelle *Bewusstes Brechen* |
| V7 | `docker run … d-check-v271 --print-config \| grep skip-allows-empty` | drei Zeilen (planning.closure, structure, reviews), je „nur mit skip-pattern" |

### Befund-Fälle (V4)

Stub-Inhalt: `# slice-001 — x` + `> **ARCHIVIERT** — Volltext: slice-001-archiv.zip`; Muster
`skip-pattern: '(?m)^> \*\*ARCHIVIERT'`; `docs/reviews/` leer; `match: name`.

| Baum | ohne Schlüssel, HEAD | ohne Schlüssel, `v0.85.0` | mit `skip-allows-empty: true`, HEAD |
|---|---|---|---|
| **Befund 1a, mit Wellen:** Stub unter `done/welle-1/`, `welle-1-results.md` flach, `recursive: true`, `skip-pattern` | Exit 1, `review-missing` „leere Pruefmenge: 0 Kandidat(en) … — fail-closed" | Exit 1, identisch | Exit 0, `0 Befund(e)` |
| dto. **ohne** `recursive` | Exit 1, leere Prüfmenge | Exit 1, identisch | **Exit 1**, leere Prüfmenge — siehe V-2 |
| **Befund 1b, flacher Stub, `require-promises: true`**, kein `skip-pattern` | Exit 1, „keine Review-Zusage unter 1 Kandidat(en) — require-promises …" | Exit 1, identisch | (Schlüssel ohne Muster ⇒ Exit 2, siehe unten) |
| dto. mit `skip-pattern` | Exit 1, leere Prüfmenge | Exit 1, identisch | Exit 0, `0 Befund(e)` |
| **Befund 1c, flacher Stub mit `skip-pattern`** | Exit 1, leere Prüfmenge | Exit 1, identisch | Exit 0, `0 Befund(e)` |
| **Umzug (R1 F-1):** flacher Stub `done/slice-001.md`, Volltext `done/wellenlos/slice-002.md` (Review-Zusage, leere Closure-Notiz), `planning.closure` + `reviews` + `structure` über `done/*.md`, je `skip-pattern`, **kein** `recursive` | Exit 1, drei Befunde: `closure-note-missing` „… ist der Bestand umgezogen?", `review-missing` leere Prüfmenge, `section-missing` „Regel trifft keine Datei …" | Exit 1, identisch (byte-gleich, V5) | Exit 0, `0 Befund(e)` — die in ADR-0106 §Konsequenzen benannte Grenze |
| **Halbe Aktivierung:** `skip-allows-empty: true` ohne `skip-pattern` | `reviews`: Exit 2 `reviews.skip-allows-empty ist ohne skip-pattern wirkungslos (halbe Aktivierung)`; `planning`: Exit 2 `planning.closure.skip-allows-empty …`; `structure`: Exit 2 `structure[0]: skip-allows-empty …` | `v0.85.0`: Exit 2 `field skip-allows-empty not found` (Schlüssel unbekannt) | — |

### Bewusstes Brechen (V6)

| Mutation (alle drei Module bzw. beide Config-Ränder) | rot | Meldung (Auszug) |
|---|---|---|
| **A** Ausnahme nie wirksam (`… && false`) | `TestReviewsRuhezustand_StubUnterWelle`, `…_FlacherStubMitRequirePromises`, `TestClosureRuhezustand`, `TestStructureRuhezustand`, `…_GemischteAusnahmen` | „erklärter Ruhezustand ⇒ kein Befund, got [… leere Pruefmenge …]" bzw. `closure-note-missing` / `section-missing` |
| **B** Bedingung „Schlüssel gesetzt" entfernt (= erste, stille Fassung) | `TestClosureRuhezustand`, `TestClosureSkipPattern_LeereMengeFailClosed`, `TestReviewsRuhezustand_StubUnterWelle`, `…_FlacherStubMitRequirePromises`, `TestStructureRuhezustand`, `TestStructureSkipPattern_LeereMengeFailClosed` | „ohne skip-allows-empty ⇒ closure-note-missing, got []", „skip-pattern ohne skip-allows-empty ⇒ leere Menge als Befund, got []", „alles ausgenommen ⇒ section-missing auf dem Glob …, got []" |
| **C** Zähl-Bedingung entfernt (mit Schlüssel jede leere Menge still; `found > 0` → `found >= 0`, damit es kompiliert) | `TestClosureRuhezustand`, `TestReviewsRuhezustand_GegenprobenBleibenRot`, `TestStructureRuhezustand`, `…_GemischteAusnahmen` | „keine passende Datei ⇒ closure-note-missing, got []", „keine Slice-Datei ⇒ Befund, auch mit skip-allows-empty, got []", „exempt-paths leert die Menge ⇒ section-missing, auch mit skip-allows-empty, got []" |
| **D** halbe Aktivierung nicht abgewiesen (`skipFehler`, `structureUeberschriftFehler`) | `TestDecode_SkipAllowsEmpty` (alle drei Blöcke) | „halbe Aktivierung: Fehler mit "reviews.skip-allows-empty" erwartet, got <nil>" |

Jede Mutation trifft genau die Tests, deren Aussage sie bricht, mit der zur Aussage passenden Meldung,
nicht über einen Nebenpfad. Ein erster Anlauf von C brach mit `declared and not used: found` ab — kein
Beleg, deshalb wiederholt.

## DoD-Abgleich

| DoD | Stand | Beleg |
|---|---|---|
| **1** Lastenheft und Spezifikation sagen den Opt-in-Schlüssel zu, Default fail-closed | **erfüllt** | Lastenheft 0.103.0: `DC-FA-PLAN-001` Kriterium „Boundary (Stub ausgenommen)" und Config-Rand, `DC-FA-STRUCT-001` Fließtext, fail-closed-Liste und Kriterien „Boundary/fail-closed (Inhalts-Ausnahme)", `DC-FA-RVW-001` Fließtext, Exit-2-Absatz, Kriterien „Boundary (alles archiviert)" und „fail-closed (halbe Aktivierung)", Historie-Zeile 0.103.0. Spezifikation: PLAN C1/C2, STRUCT Schritte 1 und 2, RVW Schritte 1 und 2, §2-Schema (drei neue Zeilen, drei `skip-pattern`-Zeilen nachgezogen), `SPEC-039`/`SPEC-049`/`SPEC-081`, Historie. Default fail-closed belegt durch V5 (byte-gleich zu `v0.85.0`). |
| **2** Module folgen, Exit 2 ohne `skip-pattern`, Tests je Modul und für die Befund-Fälle, `make gates` grün | **erfüllt** (mit V-3) | Verhalten V4; Exit 2 in allen drei Modulen V4 (letzte Zeile); je Modul „mit Schlüssel still" / „ohne Schlüssel Befund" / „keine passende Datei Befund" in `skip_ruhezustand_test.go`, rot aus dem richtigen Grund (V6 A, B, C); die drei Befund-Fälle als `TestReviewsRuhezustand_StubUnterWelle` (1a) und `…_FlacherStubMitRequirePromises` (1b und 1c in einem Test); Config-Rand `TestDecode_SkipAllowsEmpty` (V6 D); `make gates` grün selbst gefahren (V3). |
| **3** Folge-ADR, `Proposed` bis zur Closure, Index nachgezogen | **erfüllt** | `docs/plan/adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md`, `**Status:** Proposed`, Bezug auf ADR-0048 E8 und ADR-0081 E5, Re-Evaluierungs-Trigger vorhanden, verglichene Alternativen einschließlich der Zahl-Form; Index `docs/plan/adr/README.md` Zeile 116 mit Status Proposed. |
| 4, 5 | nicht Gegenstand | Review liegt vor (R1, R2); Closure folgt. |

**Plan-Abgrenzung (§1):** `exempt-paths` bleibt Leer-Ursache mit Befund (Test `TestStructureRuhezustand`,
Mutation C); `match: name` unberührt; Handbuch und README tragen den Schlüssel nicht
(`grep -rn skip-allows-empty docs/user README*` leer) — eingehalten.

## Befunde

### V-1 — INFO: Befund-Fall 1b wird nicht durch den Schlüssel allein grün

- **befund:** Der Fall „flacher Stub, `require-promises: true`" ohne `skip-pattern` bleibt rot. Der Schlüssel allein wäre Exit 2. Grün wird er erst mit `skip-pattern` **und** Schlüssel (V4). Das entspricht der Spec („Gezählt werden dabei nur die Kandidaten, die `skip-pattern` übrig lässt") und der Entscheidung im Befund („`skip-pattern` **und** `skip-allows-empty: true` setzen"). Es ist keine DoD-Verletzung. Der Test deckt beide Stufen ab.
- **verifizierbar:** ja, V4

### V-2 — LOW: Die Antwort an den Kurs nennt `recursive` für die Wellen-Form nicht

- **befund:** In der Wellen-Form (Stub unter `done/<welle-id>/`) bleibt das Gate mit `skip-pattern` und Schlüssel rot, solange `recursive` nicht gesetzt ist (V4, Zeile 2). Ohne `recursive` ist der Stub kein Kandidat, und `skip-pattern` nimmt nichts heraus. Das Verhalten ist richtig, denn es ist fail-closed nach ADR-0106 Entscheidung 4. Die Empfehlung an den Kurs nennt aber nur zwei Schlüssel: In der Befund-Datei §Entscheidung Punkt 1 steht „`skip-pattern` **und** `skip-allows-empty: true` setzen". Das Lastenheft-Kriterium „Boundary (alles archiviert)" schreibt „flach oder unter Unterverzeichnissen", ohne `recursive` zu nennen. Folgt der Kurs der Empfehlung für Befund-Zeile 1, bleibt sein Repo rot.
- **vorschlag:** Vor der Closure oder in der Release-Notiz `recursive: true` für die Wellen-Form nennen. Das betrifft die Befund-Entscheidung und gegebenenfalls das Kriterium. Das Verhalten muss nicht geändert werden.
- **verifizierbar:** ja, V4

### V-3 — INFO: `structure` hat keinen eigenen Test für „Glob trifft nichts" bei gesetztem Schlüssel

- **befund:** Für `planning` (`README.md` allein) und `reviews` (keine Slice-Datei) prüft je ein Test, dass bei gesetztem Schlüssel eine Menge ohne passende Datei rot bleibt. Bei `structure` deckt der Test mit Schlüssel diesen Fall nur in der Variante „`exempt-paths` leert die Menge" ab. Der Fall „Glob trifft nichts" ist nur ohne Schlüssel getestet (`TestStructureNullmengeFailClosed`). Die Fehlerklasse ist trotzdem bewacht: Mutation C wird über die `exempt-paths`-Variante rot. Ein eigener Test fehlt nur für die wörtliche Lesart des DoD-Punkts.
- **verifizierbar:** ja, V6 C

### Offen aus R2, nicht Gegenstand der DoD

- R2 F-3 (INFO): ADR-0048 und ADR-0081 zeigen nicht auf ADR-0106 (`grep -c 0106` = 0 in beiden).
- Das Risiko in Plan §6 ist eingetreten, wie geplant und deklariert: Mit Schlüssel ist der Umzugs-Baum still (V4, letzte Befund-Zeile). Bei der Closure braucht es seinen Ausgang.

## Negativbefunde

- **Default byte-identisch:** sechs Bäume, stdout, stderr, Exit und `--json` gleich zu `v0.85.0` (V5). Ohne Befund.
- **Umzugs-Wächter ohne Schlüssel:** in allen drei Modulen rot wie unter `v0.85.0` (V4). Ohne Befund.
- **Halbe Aktivierung:** in allen drei Modulen Exit 2, die Meldung nennt den Block und den Schlüssel (V4). Ohne Befund.
- **Tests aus dem richtigen Grund rot:** vier Mutationsklassen, jede mit passender Meldung (V6). Ohne Befund.
- **Gate-Lauf:** `make gates` selbst gefahren, Exit 0, zehn Glieder grün (V3). Ohne Befund.
- **`--print-config`:** trägt den Schlüssel in allen drei Blöcken (V7). Ohne Befund.

## Verdikt

**DoD 1–3 bestätigt.** Keine DoD-Verletzung. V-2 (LOW) sollte vor der Closure in der Befund-Entscheidung
nachgezogen werden, weil die Empfehlung an den Adopter für die Wellen-Form sonst nicht trägt. V-1 und V-3
sind INFO.
