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

- [ ] [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
      erfüllt: alle sieben Akzeptanzkriterien (die vier bestehenden
      unverändert grün, die fünf neuen — Happy, Boundary, keine Korrelation,
      Fehlkonfiguration, Default) als Tests referenziert.
- [ ] [ADR-0087](../../adr/0087-matrix-instanz-identitaets-ausnahme.md) `Accepted`, referenziert von diesem Slice.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors unverändert (kein neues
      Modul, kein neues Gate) — kein weiterer öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — falls
      eine Beobachtung anfällt, sonst „keine" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      Repo ohne Wellen-Betrieb, hier geprüft.

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
  gegen den eigenen Pfad. Die Spec ist bereits geschrieben und erlaubt es
  explizit; das Risiko ist, dass die Implementierung eine Diskrepanz zur
  Spec-Absicht aufdeckt, die eine Nacharbeit an der Spec selbst verlangt.
  **Ausgang:** wird bei Closure eingetragen.
- **Capture-Gruppen-Zählung am Regex** (`regexp.Regexp.NumSubexp()`) könnte
  Fälle mit **benannten**, aber verschachtelten oder nicht-capturing Gruppen
  (`(?:...)`) falsch zählen. **Ausgang:** wird bei Closure eingetragen.
- **Bestehende Repos mit `token`-Mustern ohne Capture-Gruppe** — das Feature
  ist opt-in (`allow-if-same-id: false` per Default), daher sollte kein
  bestehendes `.d-check.yml` brechen. **Ausgang:** wird bei Closure
  eingetragen (Gegenprobe: eigenes `.d-check.yml` bleibt unverändert grün).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — wird vor dem `git mv` nach done/ gefüllt. -->

- **Was hat funktioniert:** <wird bei Closure gefüllt>
- **Was ging anders als geplant:** <wird bei Closure gefüllt>
- **Beobachtungs-Register (`../observations/`):** <wird bei Closure gefüllt>
- **Folge-Slices:** <falls welche entstehen>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <Anker · Folge-Slice · Register, Ergebnis>

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
