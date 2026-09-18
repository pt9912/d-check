# Slice slice-229: Instanz-Identitäts-Ausnahme für `matrix` (`allow-if-same-id`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus — kein repo-weiter Beleg, den dieser Slice allein nicht liefert
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix),
[ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md). Anlass: eingehender CR
`docs/plan/cr/2026-09-18-cr-eingehend-pg-change-feed-matrix-instanz-identitaet.md`
(Adopter `pg-change-feed`), mit zwei eigenen Nachtrags-Messungen
(Präzedenzfall Supersede-Lineage-Ausnahme, Empfehlung zur schmaleren
Umsetzung).

**Berührte Spec-Stellen:** [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) ·
[`DC-FA-MTX-001.a`](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)
Schritt 7 (bereits geschrieben, siehe §1-Ausschluss).

**Verantwortlich:** claude-sonnet-5.

**Autor:** pt9912. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Rule-Feld `matrix.rules[].allow-if-same-id` implementieren:
nimmt eine Token-Form-`matrix-forbidden`-Kante aus, wenn die (bereits
vorhandene) `token`-Capture-Gruppe der Ziel-Klasse am Fund und dieselbe
Capture-Gruppe der Quell-Klasse am Pfad der Quelldatei übereinstimmende
IDs liefern — wie in [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
spezifiziert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Spec-Texte selbst** (Lastenheft-Beschreibung, Spezifikations-Algorithmus,
  Config-Schema-Zeile, Grund-Code-Zeile, Historie-Einträge beider Dokumente,
  siehe [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix))
  — **bereits geschrieben**, vor diesem Slice, als Teil der Entscheidungsfindung
  zum CR (derselbe Ablauf wie bei anderen `matrix`-Erweiterungen: Spec zuerst,
  ADR + Implementierung „folgen separat", vgl. die Historie-Einträge zu den
  Struktur- und Span-Erweiterungen des Moduls `structure`). Dieser Slice
  liefert **ADR + Implementierung + Tests**, keine weiteren Spec-Änderungen
  außer denen, die der Review dieses Slice selbst auslöst.
- **Die Instanz-Identitäts-Ausnahme für die Link-Form** von `matrix-forbidden`
  — explizit als Out-of-Scope in der Spec selbst benannt (siehe
  [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)):
  eine andere Korrelations-Frage (Pfad-Pfad statt Pfad-Prosa-Token), nicht
  Teil der CR-Bitte.
- **Eine Erweiterung der Supersede-Lineage-Ausnahme** um dieselbe
  Pfad-Korrelation — der CR bittet nicht darum, und die beiden Mechanismen
  bleiben bewusst getrennt (deklariertes Feld vs. Pfad-Muster, siehe
  zweite CR-Messung).
- **Eine Antwort an den Absender** des CR — der CR-Prozess dieses Repos
  kennt keine Pflicht zur Rückmeldung an einen externen Adopter (anders als
  bei den Antwort-Dateien zu **ausgehenden** Change Requests); die Umsetzung
  selbst ist die Antwort.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
      erfüllt: alle **acht** Akzeptanzkriterien (die drei bestehenden —
      Happy Path, Boundary, Negative — unverändert grün, die fünf neuen —
      Instanz-Identität Happy, Boundary, keine Korrelation,
      Fehlkonfiguration, Default) als Tests referenziert:
      `TestMatrixAllowIfSameID` (`internal/hexagon/core/rules/matrix_test.go`)
      und `TestDecode_MatrixAllowIfSameID{Happy,FailClosed,DefaultAus}`
      (`internal/adapter/driven/configyaml/configyaml_test.go`).
- [x] [ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md) `Accepted`, referenziert von diesem Slice.
- [x] `make gates` grün (795 Dateien, Coverage 94,70 %, 0 Lint-/Semgrep-Befunde).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: `2026-09-18-slice-229-matrix-instanz-identitaet-review-r1.md`
      (2 MEDIUM, F-1/F-2 behoben, siehe Closure-Notiz).
- [x] Doku-Update: `harness/README.md` §Sensors unverändert (kein neues
      Modul, kein neues Gate) — kein weiterer öffentlicher Vertrag berührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — zwei neue
      Verzeichnisse (`BEO-ALL/adr-fitness-function-names-wrong-gate`,
      `BEO-ALL/negativtest-deckt-nur-eine-regelseite`), je 1× (Beleg
      `evidence/slice-229.md`).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      Repo ohne Wellen-Betrieb, hier geprüft: Anker-Paarung entfällt (kein
      `liegt in`-Feld, kein Schwellen-Übertritt); Folge-Slice-Paarung
      entfällt (§7 nennt „keine"); Register-Paarung bestätigt (beide neuen
      BEO-Verzeichnisse existieren mit `observation.md`/`state.md`/
      `evidence/slice-229.md`).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area?

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0087-matrix-instanz-identitaets-ausnahme.md` | neu | ADR: Design-Entscheid (Rule-Feld statt Klassen-Feld, Token-Wiederverwendung) |
| `docs/plan/adr/README.md` | update | Index-Eintrag für [ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md) |
| `internal/hexagon/core/model/config.go` | update | `MatrixRule.AllowIfSameID bool` |
| `internal/hexagon/core/rules/matrix.go` | update | Instanz-ID-Extraktion (Quell-Pfad + Ziel-Capture), Vergleich, Ausnahme in `tokenFindings` |
| `internal/adapter/driven/config/*.go` (Config-Adapter, genauer Pfad im ersten Lauf ermittelt) | update | Fail-closed-Validierung: `allow-if-same-id: true` verlangt genau eine Capture-Gruppe auf `from`- **und** `to`-Klasse |
| `internal/hexagon/core/rules/matrix_test.go` (bzw. Pfad-Äquivalent) | update | Fünf neue Testfälle nach den [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)-Akzeptanzkriterien |
| `internal/adapter/driven/config/*_test.go` | update | Fehlkonfigurations-Testfall (Exit 2) |

**Vorab geprüft:** Spec-Texte bereits vorhanden (siehe §1-Ausschluss) —
dieser Lauf beginnt bei der ADR, nicht bei der Spec.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Nutzer bestätigt die Empfehlung aus der
zweiten CR-Messung explizit („leg los", 2026-09-18); WIP-Limit frei
(`open/`, `next/`, `in-progress/` waren leer bis auf `roadmap.md`).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die
  Config-Adapter-Validierung mehr Umbau verlangt als erwartet (z. B. wenn
  Capture-Gruppen-Zählung nicht lokal am Regex, sondern über eine tiefere
  Schema-Änderung laufen muss).
- `in-progress` → `open` (blockiert): falls `make gates` auf dem
  Endstand strukturell rot bleibt und kein Carveout trägt.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, ADR `Accepted`, unabhängiger Review
abgeschlossen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Doppelte Nutzung des `token`-Feldes** (Fließtext-Suche **und**
  Pfad-Extraktion) könnte eine Konnotation verletzen, die
  [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
  bisher implizit trug — dass `token` nur gegen **fremden** Text läuft, nie
  gegen den eigenen Pfad. **Ausgang:** entfallen — der unabhängige Review
  prüfte ADR-Treue und Spec-Konformität gezielt (siehe Report, Negativbefunde)
  und fand keine Diskrepanz; die Spec erlaubt die Doppelnutzung explizit
  und die Implementierung folgt ihr wörtlich.
- **Capture-Gruppen-Zählung am Regex** (`regexp.Regexp.NumSubexp()`) könnte
  Fälle mit **benannten**, aber verschachtelten oder nicht-capturing Gruppen
  (`(?:...)`) falsch zählen. **Ausgang:** entfallen — `NumSubexp()` ist in
  Gos `regexp`-Paket als Zahl der **capturing** Gruppen dokumentiert;
  nicht-capturing Gruppen (`(?:...)`) zählen per Sprachdefinition nicht mit,
  unabhängig von Verschachtelung. Kein repo-eigener Test nötig, um eine
  Standardbibliotheks-Zusage zu bestätigen.
- **Bestehende Repos mit `token`-Mustern ohne Capture-Gruppe** — das Feature
  ist opt-in (`allow-if-same-id: false` per Default), daher sollte kein
  bestehendes `.d-check.yml` brechen. **Ausgang:** entfallen — bestätigt:
  dieses Repos eigene `.d-check.yml` (nutzt `token` ohne `allow-if-same-id`)
  blieb über den gesamten Slice hinweg unverändert grün (`make gates`,
  wiederholt), `TestMatrixTokenReferenz` (Klassen ohne Capture-Gruppe,
  unverändert im Diff) blieb grün.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — wird vor dem `git mv` nach done/ gefüllt. -->

- **Was hat funktioniert:** Spec-first (Lastenheft + Spezifikation vor Code)
  trug wie erwartet — die Implementierung folgte der bereits geschriebenen
  Schritt-7-Beschreibung ohne Diskrepanz (Risiko 1 aus §6 entfallen). Die
  Wiederverwendung von `token` für die Quell-Pfad-Extraktion ([ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md),
  Alternative D) hielt sich als eng begründete, lokale Änderung — nur
  `matrix.go`, `config.go`, `configyaml.go` und ihre Tests berührt, keine
  andere Modul-Datei. `golangci-lint`s `gocognit`-Schwelle zwang zu einer
  Zerlegung (`applyMatrix` → drei Helfer, `tokenFindings` →
  `tokenFindingsOnLine`), die den Code im Ergebnis lesbarer macht.
- **Was ging anders als geplant:** Zwei unabhängige Korrektur-Runden nach
  Review und Verifikation. Review (F-1/F-2, beide MEDIUM): eine
  Fitness-Function-Zeile in [ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md) nannte ein Gate, das die behauptete
  Eigenschaft nicht prüft (die ADR war zu diesem Zeitpunkt bereits
  `Accepted` — Korrektur per `## Geschichte`-Nachtrag statt Kern-Edit,
  `AGENTS.md` §3.5); ein Negativtest deckte nur eine Seite einer
  symmetrischen Validierung. Verifikation: ein Zählfehler in der eigenen
  DoD-Formulierung („sieben" statt tatsächlich acht Akzeptanzkriterien,
  3+5 statt der behaupteten 4+5) — eine Instanz der in `AGENTS.md` §5
  benannten Zählmethoden-Pflicht, hier gegen die eigene Planung statt
  gegen fremden Code angewandt.
- **Beobachtungs-Register (`../observations/`):** zwei neue Verzeichnisse
  angelegt, je mit `evidence/slice-229.md` (1×):
  `BEO-ALL/adr-fitness-function-names-wrong-gate` (Review-Finding-Klasse F-1)
  und `BEO-ALL/negativtest-deckt-nur-eine-regelseite` (Review-Finding-Klasse
  F-2).
- **Folge-Slices:** keine.
- **Risiken aus §6:** alle drei entfallen — siehe §6.
- **Drei Paarungen:** Anker entfällt (kein Schwellen-Übertritt), Folge-Slice
  entfällt (keine), Register bestätigt (siehe oben).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**;
bedingt ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Default) — `internal/hexagon/core/rules/matrix.go`,
`internal/hexagon/core/model/`, der Config-Adapter und `docs/plan/adr/`
tragen keine eigene Modus-Deklaration in `harness/conventions.md` und
fallen damit unter den Default.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 46 Verzeichnisse). **Keine**
Treffer für Sub-Area `*`, die dieses Feature (Config-Modell,
Regel-Auswertung) direkt beträfen.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-18 gelesen: kein Bezug zu diesem Slice
(letzter bekannter Stand: `image-scan.yml` grün, `upstream-drift.yml`
Fremd-Bestand-Freshness, siehe slice-228-Kontext vom Vortag) — wird vor der
Closure erneut geprüft, falls sich der Stand seither geändert hat.

**Modus-Begründung:** die einzige berührte Sub-Area ist GF (Default) — kein
Begründungsblock nötig.
