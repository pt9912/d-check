# Review R2 — slice-271: Eine Menge, die erst `skip-pattern` leert, ist kein Befund

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-271, Commits `636fc084` (Plan-Änderung nach R1, Befund-Entscheidung Punkt 1)
und `d94bf010` (Umsetzung `skip-allows-empty`, ADR-0106); geprüft gegen die Befunde F-1 bis F-5 aus R1
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-271 nach der Plan-Änderung; R1 zu slice-271;
[ADR-0106](../plan/adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md) (Proposed),
[ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) (Entscheidung 8),
[ADR-0081](../plan/adr/0081-reviews-modul.md) (Entscheidung 5),
[ADR-0078](../plan/adr/0078-erklaerte-leermenge-mit-zahl.md),
[ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md) (Vergleich der Supersedes-Form);
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in);
[`MR-025`](../../harness/conventions.md#mr-025); `AGENTS.md` §3.5, §3.6, §3.7, §5 Regel 4/5.
Vorherige Findings am selben Modul: R1 zu slice-271; die Reviews zu slice-263 (dort ein Nachzug,
weil `skip-pattern` in der Exit-2-Aufzählung von `structure` fehlte — Lastenheft-Historie 0.99.1).

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | Baum „Umzug" aus R1 (flacher Stub, Volltext unter `done/wellenlos/`, ohne `recursive`), `d-check:latest` (gebaut 20:21, HEAD), alle drei Module mit `skip-pattern`, **ohne** Schlüssel | Exit 1, drei Befunde — wortgleich mit `v0.85.0` aus R1 M1 |
| M2 | derselbe Baum, je Block `skip-allows-empty: true` | `0 Befund(e)`, Exit 0 — die in ADR-0106 benannte Grenze |
| M3 | Schlüssel ohne `skip-pattern` | Exit 2, `structure[0]: skip-allows-empty ist ohne skip-pattern wirkungslos (halbe Aktivierung)` |
| M4 | Wegwerf-Klon auf `d94bf010`, Mutation: `&& …SkipAllowsEmpty` in allen drei Modulen gestrichen, halbe-Aktivierungs-Prüfung in `skipFehler` und `structureUeberschriftFehler` abgeschaltet; `make test IMAGE=d-check-r271` | rot: `TestClosureRuhezustand`, `TestClosureSkipPattern_LeereMengeFailClosed`, `TestReviewsRuhezustand_StubUnterWelle`, `…_FlacherStubMitRequirePromises`, `TestStructureRuhezustand`, `TestStructureSkipPattern_LeereMengeFailClosed`, `TestDecode_SkipAllowsEmpty` (alle drei Blöcke „got <nil>") — jeweils mit der erwarteten Meldung |
| M5 | `grep -n 'skip-allows-empty'` in `spec/` gegen die Exit-2-Aufzählungen der drei Anforderungen | §2-Schema und Lastenheft-`reviews` nennen den Fall; fünf Aufzählungen nicht (F-2) |

## Stand der R1-Befunde

| R1 | Stand | Beleg |
|---|---|---|
| F-1 HIGH (Lockerung ohne ADR) | **erledigt** | Default fail-closed (M1, M4); ADR-0106 trägt die Ausnahme mit Bezug auf ADR-0048 E8 und ADR-0081 E5; Index nachgezogen; Plan-Änderung als eigener Commit vor dem Code |
| F-2 MEDIUM (Grenzen-Liste) | **erledigt** | ohne Schlüssel ist der Umzug wieder rot (M1); mit Schlüssel steht er als Grenze in ADR-0106 §Konsequenzen, im Lastenheft (`DC-FA-RVW-001` Fließtext) und in Plan §6; der Kommentar in `planning.go` stimmt für den Default wieder |
| F-3 MEDIUM (Ungleichbehandlung zu `exempt-section-pattern`) | **teilweise** — die Leere ist jetzt deklariert; die Form der Deklaration siehe F-1 unten |
| F-4 LOW („vorbehalten") | **erledigt** | Wortlaut ersetzt; die Meldungen „nach Abzug von skip-pattern" stimmen im Default wieder |
| F-5 INFO (gemischte Leere) | **erledigt** | `TestStructureRuhezustand_GemischteAusnahmen` deckt beide Fälle (für `structure`) |

## Findings

### F-1 — MEDIUM: ADR-0106 beruft sich auf ADR-0078 und wählt die Form, die ADR-0078 verworfen hat

- **kategorie:** MEDIUM
- **quelle:** Skill-Frage 9 (Quelle über ihren Geltungsbereich zitiert); ADR-0078 Entscheidung 2
- **pfad:** `docs/plan/adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md` · „(das Muster: eine Leere, die ein generisches Muster erzeugt, wird deklariert)"
- **befund:** ADR-0078 hat für genau diese Klasse den beantragten Bool `exempt-may-empty: true` verworfen und eine **Zahl** gewählt, weil nur sie „die Prüfung im Gate-Lauf" hält und Drift laut macht; ADR-0106 nennt ADR-0078 als Muster, führt aber einen Bool derselben Art ein, und ihre Tabelle *Verglichene Alternativen* enthält die Zahl-Form nicht. Failure-Szenario (M2): mit gesetztem Schlüssel ist der Umzug neben einem Stub still — genau die Drift, die ADR-0078 mit der Zahl abfängt; ein Leser, der ADR-0106 als „nach ADR-0078" liest, hält den Wächter für gleich stark. Die Abwägung kann zugunsten des Bools ausgehen (eine Stub-Zahl wächst mit jeder Archivierung), sie steht aber nicht da.
- **verifizierbar:** nein (Begründungs-Inhalt einer ADR im Status Proposed)
- **klasse:** quelle-ueber-geltungsbereich

### F-2 — MEDIUM: Die Exit-2-Aufzählungen nennen die halbe Aktivierung nicht

- **kategorie:** MEDIUM
- **quelle:** `MR-025` (Spiegel vor dem Editieren); Skill-Frage 10; Wiederholung der Klasse aus slice-263 (Lastenheft 0.99.1)
- **pfad:** `spec/spezifikation.md` §DC-FA-PLAN-001.a C1 · „nicht kompilierendes `skip-pattern`. Ein **abwesendes**"; §DC-FA-STRUCT-001.a Schritt 1 · „`exempt-section-pattern`/`skip-pattern`;"; §DC-FA-RVW-001.a Schritt 1 · „kompilierendes `reviews.skip-pattern`. Die Meldung nennt den Schlüssel."; `spec/lastenheft.md` `DC-FA-PLAN-001` · „`boilerplate`-Eintrag sowie ein nicht kompilierendes `skip-pattern`"; `DC-FA-STRUCT-001` fail-closed-Liste · „`open-tasks-require-marker-section`/`skip-pattern`;"
- **befund:** Der Code bricht in allen drei Modulen mit Exit 2 ab, wenn `skip-allows-empty` ohne `skip-pattern` steht (M3, M4); die fünf geschlossenen Exit-2-Aufzählungen der Anforderungen und Algorithmen führen den Fall nicht — nur das §2-Schema und das Lastenheft unter `DC-FA-RVW-001` (dort auch als Kriterium „fail-closed (halbe Aktivierung)", für `planning` und `structure` fehlt ein solches Kriterium). Eine Abnahme gegen das Lastenheft sieht für `planning`/`structure` einen Exit 2, den keine Aufzählung zusagt; dieselbe Lücke wurde bei `skip-pattern` selbst schon einmal per Review nachgezogen.
- **verifizierbar:** ja — M5 (`grep`) gegen M3
- **klasse:** spiegel-exit2-aufzaehlung

### F-3 — INFO: ADR-0048 und ADR-0081 zeigen nicht auf ADR-0106

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `docs/plan/adr/0081-reviews-modul.md` · „**Entscheidung 4 abgelöst** durch"
- **befund:** Wer Entscheidung 8 bzw. 5 liest, findet die deklarierte Ausnahme nicht; ADR-0081 trägt für ADR-0105 einen `## Geschichte`-Eintrag, für ADR-0106 keinen. Der Präzedenzfall ADR-0078 → ADR-0075 hat ebenfalls keinen Rückverweis, ein `## Geschichte`-Anhang bliebe von `make adr-check` zugelassen.
- **verifizierbar:** nein
- **klasse:** adr-rueckverweis-fehlt

## Zur Zusatzfrage: ADR-0106 ohne `Supersedes`

**Richtig.** ADR-0048 E8 und ADR-0081 E5 gelten im Default unverändert (M1 byte-gleich zu `v0.85.0`);
der Schlüssel ist eine neue, opt-in Deklaration, die den Wächter nur für den erklärten Fall zurückschneidet —
dieselbe Lage wie ADR-0078 gegenüber ADR-0075, das ebenfalls ohne `Supersedes` steht. ADR-0105 brauchte
`Supersedes`, weil sie eine bestehende Entscheidung (E4: nicht rekursiv) für **dieses** Repo umkehrte; hier
wird nichts umgekehrt, und dieses Repo setzt den Schlüssel nicht. Die fehlende Sichtbarkeit von der alten
ADR aus ist F-3, kein Formfehler.

## Negativbefunde

- **Default byte-gleich:** ohne Schlüssel deckt sich der Befundsatz auf dem R1-Baum mit `v0.85.0` (M1); die beiden alten Tests zur leeren Menge sind wieder in ihrer Ursprungsform — geprüft, ohne Befund.
- **Tests aus dem richtigen Grund rot:** M4 fängt Kern- und Config-Rand-Mutationen mit den erwarteten Meldungen — ohne Befund.
- **Randfälle mit Schlüssel:** keine passende Datei, `exempt-paths` leert, unlesbares `reviews-dir`, unlesbares Unterverzeichnis bleiben Befunde (Reihenfolge der `case`-Zweige in `reviewLeerlauf` unverändert gegenüber R1) — ohne Befund.
- **`reviewLeerlauf`-Auslagerung:** reine Umformung des `switch`, Rückgabe statt `append`; Reihenfolge der Befunde unverändert (Leerlauf-Befund weiter zuletzt) — ohne Befund.
- **Plan-Änderung vor dem Code (`AGENTS.md` §6 Schritt 4):** `636fc084` liegt vor `d94bf010`, ändert Ziel, DoD, Plan-Tabelle und Risiko §6; Abgrenzung unverändert — ohne Befund.
- **ADR-Form:** Re-Evaluierungs-Trigger vorhanden, Index-Zeile nachgezogen, Status Proposed, `Schärft:` mit Kennungen, Provenance-Marker zeigt nur die Herkunft — ohne Befund. Keine Accepted-ADR berührt (§3.5).
- **Kommentare §3.7:** neue Kommentare in `configyaml.go`, `model/config.go`, `planning.go`, `reviews.go`, `structure.go` und den Tests tragen Zusage, Kopplung oder Grenze plus höchstens eine `DC-*`-Kennung — ohne Befund.
- **Sensor-Doku und Profile:** `review-coverage.md` und `verify-closure-notes.md` sagen, dass dieses Repo den Schlüssel nicht setzt; `.d-check.yml`/`.d-check.closure.yml` setzen ihn nicht — konsistent, ohne Befund.
- **`--print-config`-Vorlage:** alle drei Blöcke tragen den Schlüssel mit „nur mit skip-pattern" — ohne Befund.
- **Hexagon/Netz/Suppression:** keine neuen Importe über Schichtgrenzen, kein Netz, keine `//nolint` — ohne Befund.

## Kategorie-Summary

HIGH 0 · MEDIUM 2 (quelle-ueber-geltungsbereich, spiegel-exit2-aufzaehlung) · LOW 0 · INFO 1.
R1: F-1, F-2, F-4, F-5 erledigt; F-3 in der Deklarations-Form offen (hier F-1).

## Verdikt

**Mit Auflagen.** Kein HIGH mehr; der Gate-Pfad ist im Default wieder fail-closed und die Ausnahme hat ihre
ADR. Die beiden MEDIUM sind Text am Proposed-ADR und an der Spec, kein Verhalten — sie sollten vor dem
`Accepted`-Übergang von ADR-0106 bzw. vor der Closure nachgezogen sein.
