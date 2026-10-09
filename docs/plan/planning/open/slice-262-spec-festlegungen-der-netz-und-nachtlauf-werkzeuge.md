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

**Verantwortlich:** pt9912

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
| `harness/sensors/*.md` der betroffenen Werkzeuge | update | Verweis auf die Kennung; die gemessenen Abweichungen zur Wirklichkeit des Codes nachziehen |

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Spezifikation und die
Harness-Doku unter dem Default `*` (`ALL`). Die Skripte der Achsen liegen unter
`tools/harness/` (`HARN`), werden hier aber nur gelesen, nicht geändert — die
Sub-Area ist nicht berührt.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Eine berührt den Slice.
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(verkörpert als Schritt 17 im Workflow-Skelett): Die Sensor-Dateien
beschreiben die Werkzeuge in eigener Prosa, und genau dort liegen die
Abweichungen, die die Messung fand.

**Messung beim Beanspruchen** (am Code, nicht an der Doku):
- `harness/sensors/freshness-go.md` nennt „vier Versions-Achsen“, im Code sind
  es fünf (mit Trivy); `runtime-base-digest.md` nennt „fünf Achsen“, es sind
  sechs. `freshness-trivy` und `trivy-digest` stehen nicht im Nachtlauf
  `upstream-drift.yml`, obwohl der Gate-Index „Nachtlauf“ sagt.
- `baseline-freshness.md` sagt „Werkzeug-Ausfall → SKIP“; ein fehlendes `curl`
  endet mit Exit 1.
- `nightly-state.md` sagt, eine planmäßige Meldung werde anders behandelt; der
  Code unterscheidet nicht, er gibt einen festen Hinweis aus.
- `history-range-guard` und `selbstpruefung` treffen eigene Festlegungen und
  bekommen einen §7-Eintrag; `baseline-probe` ist der Selbsttest der
  Alias-Frage aus [`SPEC-092`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und bekommt keinen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`upstream-drift` grün (2026-10-09T07:17Z). `image-scan` rot (2026-10-09T10:37Z)
— der Lauf liegt vor dem Release v0.85.0, das die gemeldeten CVEs behebt.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
