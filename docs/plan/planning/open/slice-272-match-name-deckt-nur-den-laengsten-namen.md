# slice-272: Unter `match: name` deckt ein Report nur den längsten passenden Slice

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [Befund von `ai-harness-course`](../../cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md),
Punkt 2 (Auftraggeber-Entscheid: längster Name gewinnt).

**Berührte Spec-Stellen:** [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
samt `.a`-Algorithmus (Zuordnung unter `match: name`) und seiner Grenze zur
Präfix-Deckung.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Unter `reviews.match: name` deckt ein Report nur den längsten
Slice-Basisnamen, der in seinem Namen passt. Gezählt werden alle Slice-Dateien
unter `done/`, auch die übersprungenen Stubs — sonst deckte der Report eines
archivierten längeren Namens wieder den kürzeren. Die Grenze zur
Präfix-Deckung entfällt aus der Spezifikation.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`match: id`** — dort vergleicht der Abgleich die Kennung auf Gleichheit,
  eine Präfix-Deckung gibt es nicht.
- **Die Leere-Regel nach `skip-pattern`** — slice-271.
- **Handbuch und Release-Notiz** — Release-Prep (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Lastenheft und Spezifikation sagen die Zuordnung zu, die Grenze zur
      Präfix-Deckung ist zurückgenommen.
- [ ] Das Modul folgt ihr; Tests: der Fall des Befunds meldet `slice-cache`,
      ein archivierter längerer Stub deckt den kürzeren nicht, `slice-cachex`
      deckt `slice-cache` weiter nicht; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft; der
      Befund trägt seine Entscheidung zu Punkt 2.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Zuordnung, Grenze |
| `internal/hexagon/core/rules/reviews.go` samt Tests | update | Zuordnung |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-271 in `done/` (beide ändern
`reviews.go`); `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): der längste Name ist nicht eindeutig
  bestimmbar, etwa bei zwei gleich langen Treffern — dann zuerst die Regel
  dafür entscheiden.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Bestehende Repos werden rot** — ein Slice, den bisher nur der Report eines
  längeren Namens deckte, meldet jetzt `review-missing`. Gewollt, aber eine
  Verhaltensänderung, die die Release-Notiz nennen muss. — **Ausgang:** *(offen)*

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
