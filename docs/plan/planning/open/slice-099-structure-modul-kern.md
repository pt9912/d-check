# Slice slice-099: Modul `structure` — Kern (Abschnitt, Bereinigung, drei Grund-Codes)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-69-structure-schnitt](../welle-69-structure-schnitt.md) —
Folge-Slice aus [slice-096](../in-progress/slice-096-structure-modul-analyse.md).

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Das Modul `structure` als 20. Regelmodul: Regel-Liste, Abschnitts-Bestimmung
(inkl. Mehrdeutigkeit), Fence-Bereinigung und die Bedingungen `non-empty`,
`min-sentences`, `forbid-pattern`. Die Marken- und Zähl-Regeln folgen in
[slice-100](slice-100-structure-marken-und-zaehlung.md).

## 2. Vorgehen

1. **Gemeinsame Mechanik zuerst.** Abschnitts-Bestimmung und Fence-Bereinigung
   werden aus der Closure-Fähigkeit **herausgezogen**, nicht kopiert — die
   Spezifikation weist die eine als Preset der anderen aus, und zwei Kopien
   könnten driften, ohne dass ein Test es merkt.
2. **Drei Grund-Codes** (`section-missing`, `section-ambiguous`,
   `section-constraint`) im Lockstep: `AllReasons()`, die Doctor-Klartexte und
   die Spezifikations-§4-Tabelle im **selben** Commit — sonst bricht die
   Verriegelung.
3. **`closure-note-ambiguous`** additiv in der Closure-Fähigkeit, über dieselbe
   herausgezogene Mechanik.

## 3. Definition of Done

- [ ] Modul + Config-Schema + fail-closed Ränder; Abschnitts-Bestimmung und
      Bereinigung **geteilt** mit der Closure-Fähigkeit (ein Test fährt dieselbe
      Eingabe durch beide Oberflächen und vergleicht die Befund-Positionen).
- [ ] Vier neue Grund-Codes (drei `section-*` plus `closure-note-ambiguous`) im
      Lockstep mit `AllReasons()` und Spezifikation §4; Akzeptanzkriterien je
      Fall als Test, inklusive „Mehrdeutigkeit schlägt Messung".
- [ ] `make gates` grün; Release als **Minor** (neues Modul + additiver Code;
      ein Repo mit zwei Closure-Abschnitten wird danach rot — gehört in die
      Release-Notiz).

## 4. Risiken / offene Punkte

- **Das Herausziehen berührt ausgelieferten Code.** Die Closure-Fähigkeit ist
  seit v0.52.0 draußen; ein Refactor ihrer Abschnitts-Logik darf ihren
  Befundsatz nicht verändern. — **Ausgang:** offen; der Befund-Positions-Test
  aus DoD 1 ist die Absicherung.
- **`section-constraint` bündelt mehrere Bedingungen.** Anders als die drei
  Closure-Codes ist er ein Sammel-Code — die Unterscheidung liegt in der
  `message`, die **nicht** stabilitätsgarantiert ist. — **Ausgang:** offen; zu
  prüfen, ob das für Konsumenten reicht oder ob je Bedingung ein Code nötig ist.

## 5. Trigger

**Start** (`next` → `in-progress`): [slice-096](../in-progress/slice-096-structure-modul-analyse.md)
in `done/`; WIP-Slot frei.

**Rückführungen:** `in-progress` → `next`, falls das Herausziehen der
gemeinsamen Mechanik einen eigenen Refactor-Slice verlangt.

## 6. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). Andere
  Klasse — Referenz zwischen Dokumenten statt Form innerhalb eines; in
  [slice-096](../in-progress/slice-096-structure-modul-analyse.md) ausdrücklich
  als Nicht-Ziel dieses Moduls festgehalten.

## 7. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Anforderung, Algorithmus und ADR liegen
bereits; dieser Slice liefert, was sie versprechen.

## 8. Closure-Notiz (nach `done/`)

_Ausstehend._
