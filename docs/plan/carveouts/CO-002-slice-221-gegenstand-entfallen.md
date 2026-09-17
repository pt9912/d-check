# CO-002: `slice-221` schließt ohne Lieferung — der DoD-Wächter kennt diese Form noch nicht

**Status:** Aktiv.

**Datum angelegt:** 2026-09-16. **Letzte Prüfung:** 2026-09-16.

**Betroffenes Gate:** `verify-closure-notes` (Modul `structure`, Regel
`max-open-tasks: 0` für `## N. Definition of Done`).

**Geltungsbereich:** genau eine Datei —
[`docs/plan/planning/done/slice-221-agents-md-tabellenzellen.md`](../planning/done/slice-221-agents-md-tabellenzellen.md).

**Folge-Slice:** [`slice-225`](../planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md)

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout — ein
Carveout braucht immer einen Auflösungs-Trigger **und** einen Folge-Slice.

---

## Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — technische Begründung, keine
„noch nicht geschafft"-Aussagen.

**Baseline `v6.9.0` führt einen vierten Slice-Lifecycle-Zweig, den dieses
Repos Gate-Konfiguration noch nicht kennt.** `modul-05-planning-harness.md`
(seit `v6.6.0`→`v6.9.0`-Delta, gemessen in [`MR-072`](../../../harness/conventions.md#mr-072)):
Ein Slice, dessen Gegenstand ein anderer übernimmt oder der ganz entfällt,
geht **ohne Lieferung** nach `done/` — §7 trägt eine `**Gegenstand:**`-Zeile
statt gelieferter DoD-Punkte, und die Liefer-Häkchen bleiben leer.

[`slice-221`](../planning/done/slice-221-agents-md-tabellenzellen.md) ist
genau dieser Fall: Sein Gegenstand (`AGENTS.md` §4s Zellen kürzen) verschwand
zwischen Anlage und Wiederaufnahme, weil
[slice-222](../planning/done/slice-222-baseline-v660-bump.md) §4 radikaler
gelöst hat (ganz gestrichen statt gekürzt) — eine andere Entscheidung, kein
Bezug auf diesen Slice. Keiner seiner drei Liefer-Punkte (Zellen kürzen,
`cell-max-chars`-Wächter, Inventur) ist erfüllbar, weil ihr gemeinsamer
Gegenstand (die §4-Tabelle) nicht mehr existiert.

`.d-check.closure.yml`s `structure`-Regel für die DoD-Sektion
(`max-open-tasks: 0`) kennt diese Ausnahme **nicht** — sie verlangt
ausnahmslos, dass jeder Haken gesetzt ist, unabhängig davon, ob ein
`Gegenstand:`-Ausgang im selben Slice erklärt, warum keine Lieferung
stattfand. `make verify-closure-notes` liefe deshalb bei einem regulären
`git mv` nach `done/` auf `section-tasks-open` — der lokale
`pre-commit`-Hook löst dieses Target automatisch bei jedem
Slice-Closure-Move aus und würde den Commit blockieren.

## Auflösungs-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist" ist kein Trigger.

[`slice-225`](../planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md)
ist geschlossen: `.d-check.closure.yml`s `structure`-Regel erkennt die
`**Gegenstand:**`-Zeile und lässt offene Liefer-Häkchen dafür zu; ein
Bruch-Test bestätigt beide Richtungen (mit `Gegenstand:`-Zeile grün, ohne
weiterhin rot). Prüfbar mit:

```bash
make verify-closure-notes
```

gegen `slice-221`s Datei **ohne** den `exempt-paths`-Eintrag unten — vorher
`section-tasks-open`, danach `0 Befund(e)`.

## Geltungs-Konfiguration

| Datei | Zeile/Section | Wert |
|---|---|---|
| [`.d-check.closure.yml`](../../../.d-check.closure.yml) | `structure`-Regel für `## N. Definition of Done`, `exempt-paths` | `docs/plan/planning/done/slice-221-agents-md-tabellenzellen.md` <!-- d-check:ignore (slice-221 liegt noch nicht in done/, wandert erst mit dem Closure-Move dorthin) --> |

## Verifikation (nach Auflösung)

- [ ] Gate ist für den Geltungsbereich aktiviert (`exempt-paths`-Eintrag entfernt).
- [ ] `make gates` grün ohne Ausnahme.
- [ ] Datei wird nach `docs/plan/carveouts/done/` bewegt (reiner `git mv`). <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [ ] Folge-Slice geschlossen oder explizit dokumentiert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-16 | Angelegt | [slice-221](../planning/done/slice-221-agents-md-tabellenzellen.md) |
