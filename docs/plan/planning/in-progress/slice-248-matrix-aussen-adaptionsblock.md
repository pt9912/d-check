# slice-248: matrix-Klassen `aussen` + `adaptionsblock` adoptieren — samt Bestands-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`slice-246`](../done/welle-91/slice-246-dcheck-yml-positionen-evaluieren.md)
(dessen Probe: 29 Befunde gemessen, Kandidaten-Konfiguration nach dem
Snapshot `/tmp/aih-v6.13.0`),
[`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
(Referenz-Richtung — der Grundsatz, den die Klassen mechanisieren),
[`MR-034`](../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle)
(Vorbild-Form für die welle-Kante). Das Anforderungs-Delta (`DC-*`) entsteht
in diesem Slice, falls der Umfang es verlangt — nie per ADR.

**Berührte Spec-Stellen:** — *(Config- und Bestands-Arbeit ist kein
Spec-Stratum; entstehende Anforderungen trägt das Lastenheft)*.

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die zwei gemessenen Matrix-Positionen der Schwester werden
adoptiert — Klassen-Ordnung: `aussen` als **letzte** Klasse (First-Match),
`adaptionsblock` davor — samt Regeln `spec-straten→aussen`/`sicht→aussen` und
`spec-straten→adaptionsblock`/`sicht→adaptionsblock` — und der Bestand wird
auf den neuen Stand gezogen: die in slice-246 gemessenen Befunde (28 `aussen`-Links
aus den Straten — 20 im Lastenheft, 7 in der Spezifikation, 1 in der Sicht,
Ziele 14× Konventionsspeicher, 6× AGENTS.md, 2× Baseline-Zitat, je 1
Harness-README, Packaging, Carveout, Register, 2× CR — dazu 1 nacktes
MR-Token und 11 `matrix-inactive`-Funde der Status-Prüfung) sind je
entschieden — beseitigt oder per ADR gesichert ausgenommen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **ids-Weite `ADR-([A-Z]+-)?\d{4}`**: von slice-246 abgelehnt — keine
  segmentierte ADR-Form im Bestand (gemessen: null Treffer); die Weite wäre
  unbelegte Vorsorge. Sie wird erst bei erstem Vorkommen adoptiert.
- **exclude-sections-Scoping**: von slice-246 abgelehnt — kein eigener Fall
  (d-checks Spec-Historie trägt ein anderes Heading als die ADR-Geschichte,
  die globale Ausnahme `[Geschichte]` trifft sie nicht). Ein Scoping wäre
  Erweiterung ohne Gegenstand.
- **Status-Prüfungs-Weitung als stillschweigende Senkung**: die
  `aussen`-Klasse zieht jede Datei in `status: forbidden` — die nötigen
  Ausnahmen sind je ein §3.6-Fall (Senkung nur per ADR), nicht ein
  exempt-paths-Freifeld. Gemessen: 11 `matrix-inactive`-Funde — ADR-Index ×7
  (er führt superseded ADRs, das ist seine Aufgabe), CHANGELOG ×3,
  `docs/user/releasing.md` ×1 (Superseded-Referenz). Die Schwester nimmt den
  ADR-Index, `docs/reviews/**` und `done/welle-*.md` in status.exempt-paths —
  d-checks gemessene Ausnahme ([ADR-0097](../../adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md)) ist file-weit über
  `matrix.exempt-paths` geregelt; `docs/reviews/**` bleibt ohne Fund und
  darum ohne Ausnahme.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Beide Klassen samt Regeln stehen in `.d-check.yml` — `aussen` als
      letzte Klasse, `adaptionsblock` davor; die Straten- und Sicht-Regeln
      tragen jede eine Begründung am eigenen Bestand.
- [ ] Jeder der in slice-246 gemessenen Befunde ist entschieden — beseitigt
      (Link/Token entfernt oder umgestellt) oder per ADR gesichert
      ausgenommen; die Differenz zum Vorher-Lauf ist notiert und der Lauf
      über der Restmenge grün.
- [ ] Der Status-Seiteneffekt ist geklärt — jede nötige Ausnahme trägt
      eine ADR, oder der gemessene Fall ist auf andere Weise entschieden.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; jedes Risiko aus §6 mit Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` | update | beide Klassen + Regeln (Ordnung: adaptionsblock vor aussen, aussen zuletzt) |
| `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md` | update | Link-Nachzug: die gemessenen Verweise je entscheiden |
| Status-Fälle: `docs/plan/adr/README.md`, `CHANGELOG.md`, `docs/user/releasing.md` | update oder ADR | je `matrix-inactive`-Fund eine Fallentscheid (11 gemessen) |
| Testdatei | update | die .d-check.yml-Beispiele im Spec-Testkorpus gegen die neuen Klassen prüfen |

**Ansatz:** die Probe von slice-246 als Ausgang nehmen, die Klassen-Ordnung
korrigieren (`aussen` zuletzt — die Probe hatte sie davor, die Befund-Labels
`aussen`/`adaptionsblock` sind darum teilweise vertauscht, die Zählung
nicht), dann je Befund die Fallentscheidung dokumentieren.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt,
`Verantwortlich:` gesetzt; der Nachtlauf-Stand wird bei der Beanspruchung
gelesen ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): übersteigt der Bestands-Nachzug eine
  Sitzung (die gemessenen 40 Befunde entfalten je Link eine
  Fallentscheidung), wird der Zuschnitt auf eine Teilmenge (nur
  `adaptionsblock`, oder nur `aussen`) verengt und zurückgeführt.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Die Baseline-Zitate aus den Straten
  (Zitate in den `.harness/baseline`-Baum) stehen im Spannungsfeld: sie sind
  Zitat (citations/`d-check:cite`-Mechanik), aber `aussen`-Links. Entfernen
  bricht die Zitat-Mechanik; ausnehmen braucht eine ADR. — **Ausgang:** *(offen)*
- Die Link-Entfernung in den Spec-Straten berührt das Lastenheft
  (abnahmebindend, [`MR-032`](../../../../harness/conventions.md#mr-032)-Pflichten bei Änderung). — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug samt Spec-Stratum (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug samt Spec-Stratum)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)/[`MR-034`](../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle) tragen die
  Referenz-Richtung; die Klassen sind am Schwester-Stand belegt.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — der Bestands-Umfang ist
  gemessen (40 Befunde in slice-246, Probe mit Kandidaten-Konfiguration).
- **Reconciliation-Aufwand:** Mittel — die Fallentscheidungen sind je Link
  klein; ihre Zahl (40 gemessen) trägt der Rückführungs-Trigger in §4.
