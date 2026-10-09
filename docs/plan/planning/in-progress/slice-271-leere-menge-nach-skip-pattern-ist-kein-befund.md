# slice-271: Eine Menge, die erst `skip-pattern` leert, ist kein Befund

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [Befund von `ai-harness-course`](../../cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md),
Punkt 1 (Auftraggeber-Entscheid: alle drei Module).

**Berührte Spec-Stellen:** [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
samt ihrer `.a`-Algorithmen und der Grund-Code-Zeilen `review-missing`,
`closure-note-missing`, `section-missing`.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein Repo, dessen abgeschlossene Slices alle archiviert sind, bleibt
grün. In `reviews`, `planning.closure` und `structure` gilt gleich: Ist die
Kandidatenmenge nur deshalb leer, weil `skip-pattern` Dateien ausgenommen hat,
gibt es keinen Befund. `reviews.require-promises` urteilt nur über die nicht
übersprungenen Kandidaten. Fail-closed bleibt, wenn keine passende Datei
existiert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`exempt-paths` als Leer-Ursache** — eine namentliche Ausnahme ist eine
  Entscheidung über einzelne Dateien, kein Ruhezustand; leert sie die Menge,
  bleibt das ein Befund. Der Befund des Kurses betrifft nur `skip-pattern`.
- **Die Präfix-Deckung unter `match: name`** — slice-272.
- **Handbuch und Release-Notiz** — Release-Prep (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Lastenheft (Akzeptanzkriterien) und Spezifikation (Algorithmen,
      Grund-Code-Zeilen) der drei Anforderungen sagen die neue Regel zu.
- [ ] Die drei Module folgen ihr, je Modul ein Test für „nur übersprungen ⇒
      still" und einer für „keine passende Datei ⇒ Befund"; die drei Fälle des
      Befunds als Tests; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft; der
      Befund trägt seine Entscheidung zu Punkt 1.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Zusage, Algorithmus, Grund-Codes |
| `internal/hexagon/core/rules/` (`reviews`, `planning`, `structure`) samt Tests | update | Regel |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die drei Module brauchen je eine eigene
  Zähl-Mechanik, die sich nicht teilen lässt — dann je Modul ein Slice.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Ein zu breites `skip-pattern` macht das Gate still** — trifft das Muster
  auch Volltexte, sieht der Lauf nichts mehr und meldet es nicht. Bisher fing
  die Leere-Regel diesen Fall. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind drei Module der
Kern-Regeln und die beiden Spec-Straten — alle unter dem Default `*` (`ALL`).
Keines trägt eine eigene Konvention, einen eigenen Modus oder eine eigene
Inventur-Linie; eine eigene Sub-Area erfüllt das Kriterium nicht.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Eine berührt den Slice.
[`BEO-ALL/stilles-gruen-ueber-leerer-range`](../observations/BEO-ALL/stilles-gruen-ueber-leerer-range/observation.md)
(1×) ist dieselbe Familie mit umgekehrtem Vorzeichen: Dort ist Grün über
einer leeren Menge der Fehler, hier wird es für die durch `skip-pattern`
geleerte Menge gewollt. Der Slice muss die beiden Fälle trennbar halten — das
Risiko in §6. Keine Beobachtung erreicht mit diesem Slice 3×.

**Messung beim Beanspruchen:** Mit `v0.85.0` nachgestellt sind alle drei Fälle
des Befunds rot (`review-missing` zweimal als leere Menge, einmal als
Kandidat ohne Zusage), auf demselben Baum `closure-note-missing` und
`section-missing` (Befund-Datei §Einordnung).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`upstream-drift` grün (2026-10-09T07:17Z). `image-scan` rot (2026-10-09T10:37Z)
— der Lauf liegt vor dem Release v0.85.0, das die gemeldeten CVEs behebt;
der Slice berührt das Image nicht.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
