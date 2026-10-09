# slice-262: Festlegungen der Netz- und Nachtlauf-Werkzeuge in die Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 2,
Rest); geteilt aus slice-260, dessen Abgrenzung neun Werkzeuge mit eigener
Festlegung ergab.

**Berührte Spec-Stellen:** `spec/spezifikation.md` §7.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Festlegungen der Werkzeuge, die gegen einen fremden Stand
prüfen — über das Netz, im Nachtlauf —, stehen in §7 der Spezifikation, und
ihre Sensor-Dateien verlinken die Kennung: `image-scan` (Plattformen aus dem
Index, Plattform-Nachweis, Entscheidungslauf), die Versions-Achsen samt der
Action-Pins (Gleich/Ungleich, Präfix, fail-open), die Digest-Achsen,
`baseline-freshness` (Currency und Content-Drift) und `nightly-state` (was
gelesen werden muss, was planmäßig ist). Beim Beanspruchen zu prüfen, ob sie
eine Festlegung treffen: `history-range-guard`, `selbstpruefung` und
`baseline-probe` — keinem der Slices zugeordnet, als slice-260 geteilt wurde.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die lokalen Wächter, Hooks und Prüfer** — slice-260.
- **Die Gates aus `make gates`** — slice-259.
- **Werkzeuge, die nur bewegen oder sagen** — wie in slice-260.

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
      [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 2 mit dem
      Anteil dieses Slice vermerkt (eingelöst, sobald auch slice-260 schließt).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §7-Einträge |
| `harness/sensors/*.md` der betroffenen Werkzeuge | update | Verweis auf die Kennung |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): ein Werkzeug trifft mehr als eine
  Festlegung, die sich nicht in einer §7-Zeile fassen lässt — dann je
  Werkzeug ein Slice.

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
Harness-Doku unter dem Default `*` (`ALL`); `tools/harness/` (`HARN`), soweit
die Achsen dort liegen — beim Beanspruchen neu prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
