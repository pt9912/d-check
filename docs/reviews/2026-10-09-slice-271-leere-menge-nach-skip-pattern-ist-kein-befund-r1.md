# Review R1 — slice-271: Eine Menge, die erst `skip-pattern` leert, ist kein Befund

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-271, Commit `5dab6d41` (`git show HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-271 (§1 Ziel und Abgrenzung, §6 Risiko);
eingehender Befund `docs/plan/cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md`,
Punkt 1 samt Entscheidung;
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
und ihre `.a`-Algorithmen, SPEC-039/049/081;
[ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) (Entscheidung 8),
[ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) (Leerlauf),
[ADR-0078](../plan/adr/0078-erklaerte-leermenge-mit-zahl.md),
[ADR-0081](../plan/adr/0081-reviews-modul.md) (Entscheidung 5),
[ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md);
[`MR-025`](../../harness/conventions.md#mr-025); `AGENTS.md` §3.6, §3.7, §3.8, §5 Regel 17.
Vorherige Findings am selben Modul: die Reviews zu slice-263 (Einführung von
`skip-pattern`, „Nullmengen-Regel nach beiden Abzügen") und slice-265.

## Messungen (Kommando und Ergebnis)

Alle im Wegwerf-Baum bzw. -Klon unter dem Scratchpad, nicht im Arbeitsbaum.

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | Baum „Umzug": `done/slice-001.md` = Stub (Marker), `done/wellenlos/slice-002.md` = Volltext mit Review-Zusage und leerer Closure-Notiz; `planning.closure`, `reviews`, eine `structure`-Regel über `done/*.md`, je mit `skip-pattern`, **ohne** `recursive`; `docker run --rm --network none … --enable planning --enable reviews --enable structure` mit `v0.85.0` und `d-check:latest` (HEAD) | `v0.85.0`: Exit 1, drei Befunde (`closure-note-missing` „ist der Bestand umgezogen?", `review-missing` leere Prüfmenge, `section-missing`). HEAD: `0 Befund(e)`, Exit 0 — der Volltext unter `wellenlos/` wird von keinem der drei Module gelesen |
| M2 | Baum „gemischt": `done/slice-001.md` Stub, `done/slice-002.md` per `exempt-paths` ausgenommen; `reviews` und `structure` | `v0.85.0`: zwei Leerlauf-Befunde. HEAD: `0 Befund(e)` — die gemischte Leere ist still |
| M3 | Klon, Mutation „Ausnahme nie wirksam" (`skipped > 0` → `skipped < 0` bzw. `found > 0` → `found < 0` in allen drei Modulen), `make test IMAGE=d-check-r271` | rot: `TestReviewsRuhezustand_StubUnterWelle`, `…_FlacherStubMitRequirePromises`, `TestClosureRuhezustand`, `TestStructureRuhezustand` — je mit der erwarteten Leerlauf-Meldung |
| M4 | Klon, Mutation „jede leere Menge still" (Bedingung ohne `skipped`/`found`) | rot u. a.: `TestReviewsEmptyScopeFailsClosed`, `…_GegenprobenBleibenRot`, `TestClosureRuhezustand` („keine passende Datei"), `TestStructureRuhezustand` („exempt-paths leert") |
| M5 | Klon, Mutation `&& listErr == nil` gestrichen (reviews) und `exempt-paths` als `skipped` gezählt (structure) | rot: `…_GegenprobenBleibenRot` („unlesbares reviews-dir bleibt ein Befund, got []"), `TestStructureRuhezustand`, `TestStructureNullmengeFailClosed/exempt_leert` |
| M6 | `grep` über `spec/`, `internal/`, `harness/sensors/`, `.d-check*.yml`, `docs/user/`, `README*` nach der alten Leere-Regel (`nach dem Abzug`, `Nullmenge`, `skip-pattern`, `hebelt … nicht aus`) | Spiegel bis auf die Stellen in F-2/F-4 nachgezogen; Handbuch trägt keine Aussage zur Leere nach `skip-pattern` |

## Findings

### F-1 — HIGH: Eine Prüfregel wird gelockert, ohne ADR, gegen drei akzeptierte Entscheidungen

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.6; ADR-0048 Entscheidung 8, ADR-0081 Entscheidung 5; Skill-Frage 4
- **pfad:** `internal/hexagon/core/rules/planning.go` · „if len(names) == 0 && found > 0 {"; ebenso `structure.go` · „if len(cands) == 0 && skipped > 0 {" und `reviews.go` · „case len(candidates) == 0 && skipped > 0 && listErr == nil:"
- **befund:** Der Diff nimmt einen Fall aus dem Nullmengen-Guard, den ADR-0048 Entscheidung 8 („Fail-closed auch bei null Kandidaten … zieht der Bestand in Unterordner um, liefe das Gate sonst fortan leer und grün") und ADR-0081 Entscheidung 5 („Fail-closed bei leerer Kandidatenmenge, wie jedes andere Modul dieser Familie") ausnahmslos festlegen; es gibt weder eine Folge-ADR noch einen ADR-Bezug im Plan, nur die Entscheidung im CR-Dokument. Failure-Szenario gemessen (M1): Ein Bestand, der neben einem flachen Stub in ein nicht gelesenes Unterverzeichnis umzieht — genau der Auslöser, den ADR-0048 benennt —, war mit `v0.85.0` in allen drei Modulen rot und ist mit HEAD in allen drei still grün. Das Repo hat die Pflicht für denselben Typ Abweichung zuletzt selbst gezogen (ADR-0105: „widerspricht einer akzeptierten Entscheidung und braucht nach `AGENTS.md` §3.6 eine eigene").
- **verifizierbar:** ja — M1 gegen beide Images; kein Gate (§3.6 ist ungedeckt)
- **klasse:** gate-lockerung-ohne-adr

### F-2 — MEDIUM: Die Grenzen-Liste nennt die kleinere Lücke, und ein Kommentar verspricht die größere als gedeckt

- **kategorie:** MEDIUM
- **quelle:** Skill-Frage 18 (Grenzen-Liste ohne größte Lücke); `AGENTS.md` §3.8; `DC-FA-PLAN-001`
- **pfad:** `internal/hexagon/core/rules/planning.go` · „(typischer Auslöser: der Bestand wandert in Unterordner, das Gate liefe fortan leer und grün)"; `spec/spezifikation.md` §DC-FA-PLAN-001.a C2 · „wandert in Unterordner) und das Gate liefe fortan leer und grün"; `harness/sensors/review-coverage.md` · „trifft das Muster auch Volltexte, wird der Lauf damit still"
- **befund:** Alle nachgezogenen Grenzen (Lastenheft, Spezifikation, beide Sensor-Dateien, Plan §6) nennen als Preis nur das zu breite Muster. Die gemessene größere Lücke (M1) — ein einziger übersprungener Stub schaltet den Umzugs-Wächter ab, auch bei korrektem Muster, sobald die Volltexte außerhalb der gelesenen Menge liegen (ohne `recursive`, oder unter `SKIP_DIRS`) — steht nirgends; der Kommentar direkt unter dem neuen Block und Spez-Schritt C2 behaupten den Umzugsfall weiterhin als gefangen.
- **verifizierbar:** ja — M1
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-3 — MEDIUM: Zwei Abzüge desselben Moduls, gleiche Gefahr, verschiedene Antwort ohne benannten Grund

- **kategorie:** MEDIUM
- **quelle:** Skill-Frage 11; ADR-0078 (Entscheidungen 1–2), ADR-0075
- **pfad:** `internal/hexagon/core/rules/structure.go` · „Leert erst skip-pattern sie, ist das der Ruhezustand"
- **befund:** In `structure` muss eine von `exempt-section-pattern` geleerte Menge mit `exempt-expect-count` **deklariert** werden, weil ADR-0078 für ein generisches Muster eine stille Erlaubnis ausdrücklich verworfen hat („Ohne diese Antwort schaltete ein zu breites Muster die Regel still ab"); eine von `skip-pattern` — ebenfalls ein generisches RE2-Muster — geleerte Menge ist nach diesem Diff ohne Deklaration still. Failure-Szenario: ein Muster ohne Zeilenanker oder ohne `— Volltext:`, das Volltexte trifft, die den Marker zitieren, leert die Regel; bei `exempt-section-pattern` wäre derselbe Fehler rot, hier ist er grün. Weder Spec noch Plan nennen den Grund für die Ungleichbehandlung.
- **verifizierbar:** ja — Testbaum mit einem Volltext, der `> **ARCHIVIERT** — Volltext:` zitiert, als einziger Kandidatin
- **klasse:** modul-inkonsistenz-gleiche-eingabeklasse

### F-4 — LOW: „bleibt dem Fall vorbehalten" schließt Fälle aus, die weiter melden

- **kategorie:** LOW
- **quelle:** `DC-FA-PLAN-001`, `DC-FA-STRUCT-001` (Doku-Drift gegen SPEC-039/SPEC-049)
- **pfad:** `spec/lastenheft.md` · „`closure-note-missing` auf dem Verzeichnis bleibt dem Fall vorbehalten, dass es keine einzige passende Datei enthält"; ebenso „`section-missing` auf dem Glob bleibt dem Fall vorbehalten, dass `skip-pattern` keine Datei ausgenommen hat"
- **befund:** Auf dem Verzeichnis meldet `closure-note-missing` auch ein fehlendes oder unlesbares `closure.dir` (SPEC-039), auf dem Glob meldet `section-missing` auch ohne gesetztes `skip-pattern`; „vorbehalten" liest sich als abschließend. Daneben nennen die unveränderten Meldungen „(nach Abzug von skip-pattern …)" bzw. „auch nach Abzug von exempt-paths und skip-pattern" einen Abzug, der in ihrem einzigen verbliebenen Auslösefall nichts genommen hat.
- **verifizierbar:** nein (Wortlaut)
- **klasse:** doku-drift-ausschliesslichkeit

### F-5 — INFO: Die gemischte Leere ist zugesagt, aber nicht getestet

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/skip_ruhezustand_test.go` · „leert exempt-paths die Menge, bleibt section-missing"
- **befund:** Nimmt `exempt-paths` eine Datei und `skip-pattern` die übrigen, ist der Lauf still (M2) — spec-konform („mindestens eine Datei ausgenommen"), aber ohne Test; ebenso, dass eine Datei, die beide trifft, als `exempt` und nicht als übersprungen zählt (Reihenfolge-Abhängigkeit in `reviewWalk.visit` und `checkStructureRule`).
- **verifizierbar:** ja — M2
- **klasse:** randfall-ohne-test

## Negativbefunde

- **Randfälle im Code:** unlesbare Datei bleibt Kandidatin (`closureSkip`, `structureSkipped`, `reviewSkipped` geben bei Lesefehler „nicht übersprungen"); `skip-pattern` gesetzt ohne Treffer meldet weiter (M4); unlesbares Unterverzeichnis unter `recursive` meldet in `planning` (Fehler aus `closureCandidates` vor der Ausnahme) und in `reviews` (`badDirs` vor der Ausnahme); unlesbares `reviews-dir` neben lauter Stubs meldet (M5); Zusagen nur in übersprungenen Stubs zählen nicht (`…_RequirePromisesZaehltNurUebrige`); `found` zählt in `planning` nur Glob-Treffer, eine `README.md` hebt die Ausnahme nicht aus — geprüft, ohne Befund.
- **Tests aus dem richtigen Grund rot:** drei Mutationen (M3–M5), jede vom passenden Test mit der erwarteten Meldung gefangen — geprüft, ohne Befund.
- **Spiegel-Liste (MR-025):** Lastenheft (drei Kriterien, Fließtext, Historie 0.103.0), Spezifikation (C2, STRUCT Schritt 2, RVW Schritte 2 und 5, §2-Schema vier Zeilen, SPEC-039/049/081, Historie), `model/config.go`, `config_template.go`, beide Sensor-Dateien, `.d-check.yml`/`.d-check.closure.yml` (Kommentare sagen nichts über die Leere) — nachgezogen bis auf F-2/F-4; Handbuch ohne Aussage dazu, Release-Prep laut Plan (`AGENTS.md` §5 Regel 17). ADRs: ADR-0048 und ADR-0081 sagen weiter das Gegenteil (immutabel, siehe F-1).
- **Kommentar-Regel §3.7:** neue Kommentare in `planning.go`, `reviews.go`, `structure.go`, `config.go` und im Testfile tragen Zusage/Grenze plus höchstens eine `DC-*`-Kennung; keine Review-Historie, keine Slice-Nummer — ohne Befund (der veraltete Nachbar-Kommentar ist F-2).
- **Abgrenzung des Plans:** `exempt-paths` bleibt Leer-Ursache mit Befund; `match: name` (slice-272) und Handbuch unberührt — ohne Befund.
- **Hexagon/Netz/Suppression:** keine neuen Importe, kein Netzzugriff, keine `//nolint` — ohne Befund.
- **Commit-Botschaft:** Behauptungen (drei Module, drei Fail-closed-Fälle, Spiegel-Liste) decken sich mit dem Diff; die Grenze aus F-2 fehlt dort ebenso — kein eigener Befund.

## Kategorie-Summary

HIGH 1 (gate-lockerung-ohne-adr) · MEDIUM 2 (grenzen-liste-ohne-groesste-luecke, modul-inkonsistenz-gleiche-eingabeklasse) · LOW 1 · INFO 1.

## Verdikt

**Blockiert.** F-1 verlangt eine Folge-ADR (Supersedes ADR-0048 Entscheidung 8 und ADR-0081 Entscheidung 5, für den Fall der durch `skip-pattern` geleerten Menge), die F-2 und F-3 als Konsequenz bzw. verglichene Alternative trägt — etwa die Deklarations-Form aus ADR-0078. Die Entscheidung des Auftraggebers wird damit nicht bestritten; ihr fehlt der Träger, den §3.6 verlangt. Bei Widerspruch des Implementers ist F-1 ein HIGH mit Rollen-Widerspruch und läuft über den Architect (`v6.17.0` · `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz).
