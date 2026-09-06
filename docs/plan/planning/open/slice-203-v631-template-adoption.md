# Slice slice-203: Adoption der v6.3.1-Template-Deltas

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist von der DoD dieses Slice
nicht verschieden.

**Bezug:** [`MR-065`](../../../../harness/conventions.md#mr-065) (der Bump, der
die Vorlagen brachte), [`MR-021`](../../../../harness/conventions.md#mr-021).
Ob die Adoption eigene `MR`-Einträge braucht, entscheidet der Slice — nicht
dieser Plan.

**Berührte Spec-Stellen:** — (Adoption von Harness-Form, keine Spec-Stelle).

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Die vier Delta-Stränge, die
[`MR-065`](../../../../harness/conventions.md#mr-065) auf Regelwerk-Ebene
adoptiert hat, in der **verkörperten Form** dieses Repos umsetzen — allen
voran die Auslagerung überladener Sensors-Einträge nach
`harness/sensors/<target>.md`. <!-- d-check:ignore (Kanon-Form aus v6.3.1, wird mit diesem Slice angelegt) -->

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und die
Closure-Pflichten darunter zählen nicht mit.

- [ ] **(1)** Sensors-Auslagerung: jedes Gate, dessen Vertrag mehr als **einen
      Satz** braucht, hat seine eigene Datei; die Tabellenzelle wird zum Link
      darauf. Kandidaten sind gemessen, nicht geschätzt — die Zellen-Länge in
      `harness/README.md` §Sensors entscheidet.
- [ ] **(2)** Nicht-Gates stehen in einer **zweiten Tabelle** und tragen
      `kein Gate` **in der Zeile**, nicht in Prosa daneben; die
      Vollständigkeits-Zeilen nennen ihr **Kommando** statt einer
      eingefrorenen Zahl.
- [ ] **(3)** `AGENTS.md` §6 trägt den Rollenwechsel-Zusatz nach Schritt 8
      (Handoff an Reviewer, kein Self-Review); die Slice-Vorlage-Nachfolge
      (fünfter Closure-Punkt *Review durchgeführt*) ist entschieden —
      template-forward oder mit benannter Bestands-Grenze.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/sensors/` (neu) <!-- d-check:ignore (wird mit diesem Slice angelegt) --> | neu | ein Gate je Datei, sobald sein Vertrag mehr als einen Satz braucht |
| `harness/README.md` | update | §Sensors auf Index-Form; zweite Tabelle für Nicht-Gates |
| `AGENTS.md` | update | §6 Rollenwechsel nach Schritt 8; §4-Tabelle folgt derselben Frage |
| Slice-Vorlage-Nachfolge | Entscheid | fünfter Closure-Punkt — template-forward oder mit Bestands-Grenze |

## 4. Trigger

**Start** (`open` → `in-progress`): [`MR-065`](../../../../harness/conventions.md#mr-065)
ist geschlossen, die Vorlagen liegen unter
`.harness/baseline/v6.3.1/templates/`. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt die Messung, dass die
  Sensors-Auslagerung allein mehr als drei Liefer-Punkte trägt — d-checks
  §Sensors-Tabelle ist groß —, wird sie ein eigener Slice und die beiden
  übrigen Stränge laufen getrennt.
- `in-progress` → `open` (blockiert): Erweist sich die Auslagerung als
  Formfrage, die einen eigenen Entscheid braucht (welche Zellen zählen als
  „mehr als ein Satz"), ruht der Slice bis zu diesem Entscheid.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien **und** ein Lerneintrag: (a) `make gates` grün und
jede ausgelagerte Sensor-Datei ist aus ihrer Tabellenzeile heraus verlinkt —
der Link-Sensor hält die Zuordnung; (b) kein Nicht-Gate steht mehr unmarkiert
in der Gate-Tabelle. Dazu die Closure-Notiz.

## 6. Risiken und offene Punkte

- **Die Auslagerung wird zur Umschichtung ohne Gewinn** — der Kanon warnt
  ausdrücklich davor, `harness/sensors/` <!-- d-check:ignore (Kanon-Form aus v6.3.1) --> zur Halde zu machen, die die
  Sektion vorher war: hinein gehört Grenze und Bedienvertrag, **nicht** der
  Deckungsnachweis des Werkzeugs. — **Ausgang:** <offen>
- **„Mehr als ein Satz" ist eine Urteilsgrenze, kein `grep`** — die
  Kandidatenmenge kann zwischen zwei Läufen anders ausfallen. — **Ausgang:** <offen>
- **Der Bestand ist groß:** d-checks §Sensors-Tabelle trägt über 40 Zeilen,
  viele davon absatzlang. Die Größenregel des Slice kann daran brechen. —
  **Ausgang:** <offen>

## 7. Closure-Notiz

<wird vor dem `git mv` nach `done/` gefüllt>

## 8. Sub-Area-Modus-Begründung

<die drei Vorprüfungen entstehen spätestens bei der Beanspruchung; ein Plan in
`open/` trägt den Nachtlauf-Block noch nicht>
