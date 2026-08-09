# Slice slice-099: Modul `structure` — Implementierung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** noch keiner Welle zugeordnet — **nicht** welle-69: deren
Out-of-Scope schließt die Implementierung ausdrücklich aus, und sie ist mit
[slice-096](../done/slice-096-structure-modul-analyse.md) geschlossen. Dieser
Slice ist ihr **derivatives** Ergebnis und bekommt seine Welle bei der
Eröffnung der Umsetzungs-Welle.

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Preset-Kopplung + `closure-note-ambiguous`),
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Das Modul `structure` vollständig — Regel-Liste, Kandidaten-Menge,
Abschnitts-Findung mit beiden Kardinalitäts-Modi, Bereinigung und **alle sechs**
Bedingungen — plus die Preset-Kopplung der Closure-Fähigkeit samt
`closure-note-ambiguous`.

## 2. Warum in **einem** Slice, nicht in zweien

Der erste Zuschnitt trennte Kern und Marken/Zählung. Das ist an der
**Release-Grenze** gescheitert (Befund des Umsetzbarkeits-Reviews): das
veröffentlichte Schema führt alle Schlüssel, und die Config-Dekodierung ist
strikt — ein Release nach dem Kern allein wiese eine spezifikationskonforme
Konfiguration mit Exit 2 ab. Ein Modul, das die Hälfte seines eigenen Schemas
ablehnt, ist kein lieferbarer Zwischenstand.

## 3. Vorgehen

1. **Gemeinsame Mechanik zuerst herausziehen**, nicht kopieren:
   Abschnitts-Findung, Kardinalitäts-Behandlung, Fence-Bereinigung. Die
   Spezifikation weist die Closure-Fähigkeit als Preset aus — zwei Kopien
   könnten driften, ohne dass ein Test es merkt.
2. **Grund-Codes im Lockstep — es sind neun, nicht sieben:** die sechs
   Bedingungs-Codes (`section-empty`/`-thin`/`-oversized`/`-forbidden`/
   `-pattern-missing`/`-marker-missing`) **plus** die zwei Struktur-Codes
   (`section-missing`, `section-ambiguous`) **plus** `closure-note-ambiguous`.
   Alle neun gehören mit `AllReasons()`, den Doctor-Klartexten und der
   Spezifikations-§4-Tabelle in **denselben** Commit.
3. **CLI-Mit-Modifikation**, die jedes bisherige Modul mitgebracht hat:
   `--print-config`-Gerüst, `--suggest-config`-Vorlage und die
   `--print-mk`-Target-Liste.

## 4. Definition of Done

- [ ] Modul vollständig (beide `sections`-Modi, alle sechs Bedingungen,
      Kandidaten-Menge, fail-closed Ränder inkl. Leerlauf trotz `exempt-paths`);
      Abschnitts-Mechanik **geteilt** mit der Closure-Fähigkeit — belegt durch
      den Preset-Kopplungs-Test (dieselbe Eingabe, beide Oberflächen, gleiche
      Befund-Zeilen).
- [ ] **Neun** neue Grund-Codes im Lockstep mit `AllReasons()` und Spezifikation
      §4; jedes Akzeptanzkriterium als Test, insbesondere die drei Marken-Formen,
      `sections: each` und „Mehrdeutigkeit schlägt Messung".
- [ ] CLI-Enumerationen nachgezogen — einschließlich der Target-Liste in
      [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
      deren Out-of-Scope-Zeile die **Zahl** der Targets festschreibt und daher
      mitzuändern ist. `make gates` grün; **Paritäts-Beleg** gegen die
      beigezogenen Adopter-Fixtures; Release als **Minor** (neues Modul +
      additiver Code — ein Repo mit zwei Closure-Abschnitten wird danach rot,
      das gehört in die Release-Notiz).

## 5. Risiken / offene Punkte

- **Das Herausziehen berührt ausgelieferten Code.** Die Closure-Fähigkeit ist
  seit v0.52.0 draußen; ein Refactor ihrer Abschnitts-Logik darf ihren Befundsatz
  nicht verändern. — **Ausgang:** offen; der Preset-Kopplungs-Test ist die
  Absicherung.
- **Die Marken-Syntax ist eine Konvention, keine Norm.** `**M:**` ist die Form
  der beiden vermessenen Repos; ein drittes schreibt vielleicht anders.
  — **Ausgang:** offen; zu entscheiden, ob die Marke konfigurierbar wird oder
  die Konvention Teil der Zusage bleibt.
- **Die Paritäts-Fixtures liegen im Adopter-Repo**, nicht hier. — **Ausgang:**
  offen; beizuziehen, nicht nachzubauen.
- **Die Fence-Lexik trägt einen bekannten Defekt**
  ([slice-101](../in-progress/slice-101-fence-unbalanciert.md)): ein ungeschlossener Fence
  verschluckt still den Rest. `structure` erbt ihn über die geteilte Mechanik.
  — **Ausgang: entschieden.** Der Fix läuft zuerst; §6 macht ihn zur **bindenden**
  Start-Bedingung, nicht zur Empfehlung.

## 6. Trigger

**Start** (`next` → `in-progress`): [slice-096](../done/slice-096-structure-modul-analyse.md)
in `done/` **und** [slice-101](../in-progress/slice-101-fence-unbalanciert.md) in `done/`;
WIP-Slot frei. Die zweite Bedingung ist **bindend**, nicht bevorzugt: liefe
dieser Slice zuerst, erbte das neue Modul über die geteilte Mechanik einen
**bekannten** stillen Grün-Pfad — und läge damit ab Release als zugesagte
Fähigkeit vor, die ihre Zusage an einer bekannten Stelle nicht hält.

**Rückführungen:** `in-progress` → `next`, falls das Herausziehen der
gemeinsamen Mechanik einen eigenen Refactor-Slice verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001** (andere
  Klasse — Referenz zwischen Dokumenten statt Form innerhalb eines; in
  [slice-096](../done/slice-096-structure-modul-analyse.md) ausdrücklich
  als Nicht-Ziel festgehalten) und **BEO-002** (Semantik-Änderungen werden nur im
  Dokumentkörper nachgezogen). BEO-002 **betrifft diesen Slice unmittelbar**: er
  ändert die Grund-Code-Menge und damit gleich mehrere Spiegel — `AllReasons()`,
  Doctor-Klartexte, Spezifikation §4, drei CLI-Enumerationen und die
  Handbuch-Modulliste. Sie gehören vor dem Editieren aufgelistet.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Anforderung, Algorithmus und ADR liegen
bereits; dieser Slice liefert, was sie versprechen.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
