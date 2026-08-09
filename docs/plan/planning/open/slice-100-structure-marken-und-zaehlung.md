# Slice slice-100: Modul `structure` — Marken und abschnitts-treue Zählung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-69-structure-schnitt](../welle-69-structure-schnitt.md) —
Folge-Slice aus [slice-096](../in-progress/slice-096-structure-modul-analyse.md).

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md)
(Entscheidung 5).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Die beiden Bedingungen, die kein bestehendes Modul kennt: `max-tasks`
(abschnitts-treue Zählung von Task-Items) sowie `require-all` und `require-any`
(benannte, zeilenverankert hervorgehobene Marken).

## 2. Warum das die heikleren zwei sind

Genau an ihnen sind die abgelösten Skripte gescheitert — und zwar an der
*Abschnitts*-Treue, nicht an der Zählung selbst: das Skript zählte **dateiweit**
und musste den Abschnitts-Schnitt nachbauen. Dieselbe Falle bei den Marken: ein
Teilstring-Vergleich hielte ein `Boundary` mitten im Satz für erfüllt und
erzeugte ein **Falsch-Grün** — teurer als ein Falsch-Rot, weil es unbemerkt
bleibt.

Die zwei Marken-Formen sind kein Komfort: eine Anforderung muss **alle**
Akzeptanz-Bausteine tragen, ein Lerneintrag genügt mit **einer** von mehreren
zulässigen Formen. Mit nur einer Form bliebe die zweite Prüfung ungedeckt und
der Adopter behielte ein Skript für einen einzigen Fall.

## 3. Definition of Done

- [ ] `max-tasks` abschnitts-treu und fence-treu; Test, der dateiweite Zählung
      ausschließt (Items **außerhalb** des Abschnitts dürfen nicht zählen).
- [ ] `require-all` und `require-any` mit zeilenverankerter, hervorgehobener
      Marke; Test, der den Teilstring-Fall als **verletzt** belegt.
- [ ] **Paritäts-Beleg** gegen die beigezogenen Adopter-Fixtures für die beiden
      Prüfungen; `make gates` grün; Release als Minor.

## 4. Risiken / offene Punkte

- **Die Marken-Syntax ist eine Konvention, keine Norm.** `**M:**` ist die Form
  des Antragstellers; ein anderes Repo schreibt vielleicht `__M__` oder eine
  Definition-List. — **Ausgang:** offen; zu entscheiden, ob die Marke
  konfigurierbar wird oder die Konvention Teil der Zusage bleibt.
- **Task-Item-Erkennung** ist Markdown-Lexik und damit ein Ort, an dem eigene
  Heuristik schon einmal teuer war. — **Ausgang:** offen; die vorhandene
  Listen-/Fence-Lexik nutzen, nicht neu schreiben.

## 5. Trigger

**Start** (`next` → `in-progress`): [slice-099](slice-099-structure-modul-kern.md)
in `done/` (der Kern trägt Abschnitts-Bestimmung und Bereinigung); WIP-Slot frei.

**Rückführungen:** `in-progress` → `open`, falls die Marken-Syntax-Frage eine
eigene Vertragsentscheidung verlangt.

## 6. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**; andere
  Klasse, nichts zu berücksichtigen.

## 7. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die Zusage steht bereits im Lastenheft, der
Code liefert sie.

## 8. Closure-Notiz (nach `done/`)

_Ausstehend._
