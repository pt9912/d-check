# Slice slice-204: Auflösung für den Slice, der während der Arbeit wächst

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
(3×, Schwelle mit slice-203 erreicht). Ob die Auflösung eine Hard Rule, ein
`MR`-Eintrag oder eine Sensor-Regel wird, entscheidet der Slice — dieser Plan
setzt es nicht voraus.

**Berührte Spec-Stellen:** — (Planungs-Form, keine Spec-Stelle).

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Die Spannung auflösen, die drei Slices unabhängig voneinander getroffen hat:
**„einen Slice nicht über die Ein-Sitzungs-Review-Grenze wachsen lassen"**
gegen **„ihn nicht mitten in einem Beleg teilen, der nur im Zusammenhang
trägt"**. Der Kanon kennt für den Konflikt keine Antwort; er kennt nur die
Rückführung, und die wurde dreimal aus je gutem Grund nicht gezogen.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Entscheid | neu | welche Artefaktklasse die Auflösung trägt — Hard Rule, `MR` oder Sensor |
| Träger der Entscheidung | neu | erst nach dem Entscheid benennbar |

**Die drei Instanzen, gemessen** — sie teilen die Folge, nicht den Grund:

| Vorgang | Warum nicht zurückgeführt |
|---|---|
| slice-195 | Eine Teilung hätte den Zähler-Diff-Beleg der Migration zerrissen |
| slice-197 | Eine Teilung hätte die Werkzeug-Korrektur mehrfach nachziehen müssen |
| slice-203 | Kein einzelner Nachsteuerungs-Schritt sprengte die Grenze; die Summe schon |

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Rückführung selbst abschaffen oder aufweichen.** Sie ist Kanon
  (`modul-05` §Trigger je Lifecycle-Übergang) und bleibt es; gesucht ist die
  Antwort für den Fall, dass ihr etwas entgegensteht — nicht ihre Ersetzung.
- **Eine Größen-Metrik erfinden.** „In einer Review-Sitzung prüfbar" ist ein
  Urteil; wer daraus eine Zahl macht, tauscht ein ehrliches Urteil gegen eine
  falsche Genauigkeit.
- **Rückwirkende Anwendung** auf slice-195, slice-197 und slice-203. Ihre
  Lauf-Belege bleiben, wie sie sind.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Der Entscheid steht: Welche Artefaktklasse trägt die Auflösung,
      und **warum diese** — mit der Alternative, die verworfen wurde.
- [ ] **(2)** Die Auflösung ist geschrieben und trägt den Herkunfts-Anker
      `seit slice-204`; das Beobachtungs-Register führt den Eintrag von
      `geplant` auf `verkörpert` mit Zielort.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Auflösung wird eine Erlaubnis statt einer Regel.** „Wachsen ist okay,
  wenn man es benennt" wäre keine Antwort, sondern die Aufgabe der Grenze —
  und würde die Beobachtung nicht schließen, sondern legalisieren.
  — **Ausgang:** <offen>
- **Der Slice wächst selbst über seine Grenze.** Die Ironie ist real und der
  Grund, warum die Auflösung nicht in slice-203 geschrieben wurde.
  — **Ausgang:** <offen>
- **Drei Instanzen mit drei verschiedenen Gründen tragen womöglich keine
  gemeinsame Regel.** Dann ist das Ergebnis eine benannte Grenze statt einer
  Auflösung — auch das ist ein Ausgang, aber er gehört ausgeschrieben.
  — **Ausgang:** <offen>

## 6. Trigger

**Start** (`open` → `in-progress`): slice-203 ist geschlossen, der WIP-Slot
frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass der Entscheid eine ADR
  braucht statt einer Regelzeile, wird die ADR ein eigener Vorgang.
- `in-progress` → `open` (blockiert): Trägt keine gemeinsame Regel (drittes
  Risiko), ruht der Slice, bis eine vierte Instanz die Klasse schärft.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) der
Registereintrag steht auf `verkörpert` mit auflösendem Zielort; (b) der
Zielort trägt den Herkunfts-Anker.

## 7. Vorgelagert (vor der Modus-Begründung)

<die drei Vorprüfungen entstehen spätestens bei der Beanspruchung; ein Plan in
`open/` trägt den Nachtlauf-Block noch nicht>

## 8. Sub-Area-Modus-Begründung

<entsteht mit der Beanspruchung>

## 9. Closure-Notiz (nach `done/`)

<wird vor dem `git mv` nach `done/` gefüllt>
