# slice-266: Die Doku der Module ist netzlos aus dem Image lesbar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** eingehender CR von `sf-connector`
([`docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md`](../../cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md),
Punkt 4): ein netzlos arbeitendes Repo kann die Erkennungsregel eines Moduls
heute nur durch Probieren ermitteln.

**Berührte Spec-Stellen:** CLI-Anforderungen im Lastenheft (Hilfe-Ausgabe) —
welche, entscheidet der Schnitt beim Beanspruchen.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Wer nur das Image hat, liest aus ihm, was ein Modul prüft und welche
Schlüssel es nimmt — über eine Ausgabe des Werkzeugs (etwa `--help <modul>`)
oder eine mitgelieferte Datei. Beim Beanspruchen wird entschieden, welche Form
und welche Quelle; die Quelle darf nicht zu einer zweiten Beschreibung neben
Spezifikation und Handbuch werden, die driftet.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Erkennungsregel von `reviews` selbst** — slice-265.
- **Eine Übersetzung oder Kürzung des Handbuchs** — der Gegenstand ist die
  Erreichbarkeit, nicht ein neuer Text.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Die Doku der Module ist aus dem Image ohne Netz lesbar; die Form ist im
      Lastenheft zugesagt und durch einen Test gehalten.
- [ ] Ein Sensor oder Test hält die mitgelieferte Doku gegen ihre Quelle, damit
      sie nicht driftet; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft; der CR
      trägt seine Entscheidung zu Punkt 4.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| CLI und Image | update | Ausgabe bzw. mitgelieferte Doku |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Zusage |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die gewählte Quelle verlangt einen Umbau
  des Handbuchs — dann erst die Quelle, dann die Ausgabe.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Zweite Beschreibung** — eine eigene Modul-Doku im Binary driftete gegen
  Spezifikation und Handbuch. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Produkt unter dem Default `*`
(`ALL`); beim Beanspruchen neu prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
