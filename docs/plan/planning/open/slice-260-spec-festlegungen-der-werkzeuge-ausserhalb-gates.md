# slice-260: Festlegungen der Harness-Werkzeuge außerhalb von `make gates` in die Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 2,
Rest); Folge von slice-259, der den Abschnitt §7 der Spezifikation anlegt.

**Berührte Spec-Stellen:** `spec/spezifikation.md` §7.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Festlegungen der Harness-Werkzeuge, die nicht in `make gates`
laufen und keine eigene Anforderung verfeinern, stehen in §7 der
Spezifikation, und ihre Sensor-Dateien verlinken die Kennung: `image-scan`
(Plattform-Nachweis, Entscheidungslauf), `blackbox-probe` (Kanarienlauf,
Vergleich), der Tool-Call-Wächter und die Hooks, die Frische-Achsen,
`nightly-state`, `record-gates`, `archive-wave` — je nachdem, welche davon
eine Festlegung treffen statt nur zu bewegen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Abschnitt selbst und die Gates aus `make gates`** — slice-259.
- **Werkzeuge, die nur bewegen oder sagen** (`slice-mv`, `help`, `clean`,
  `versions`) — sie treffen keine Festlegung, die man fortschreiben müsste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Je Werkzeug mit eigener Festlegung ein §7-Eintrag, am Code geprüft.
- [ ] Die Sensor-Dateien verlinken die Kennung statt Schwelle und Randform zu
      führen; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft;
      [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 2 als
      eingelöst vermerkt.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §7-Einträge |
| `harness/sensors/*.md` der betroffenen Werkzeuge | update | Verweis auf die Kennung |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-259 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): mehr als etwa sechs Werkzeuge treffen
  eine eigene Festlegung — dann nach Wächter/Hooks und Nachtlauf-Werkzeugen
  teilen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Keine bekannt.

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

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** die Spezifikation und die
Harness-Doku unter dem Default `*` (`ALL`); `tools/harness/` (`HARN`) ist
berührt, soweit Wächter und Hooks dort liegen — beim Anlegen des Plans für
die Umsetzung neu prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen;
Stand beim Schnitt wie slice-259.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
