# Review-Report: slice-193 — 2026-09-03

**Review-Art:** Code — geprüft gegen Plan, ADR und Konventionen (Modul 10
§Drei Review-Arten).

**Gegenstand:** Commits `faac05c` (MR-058 → `done/`), `8bd9ab7`
(feat: Baseline-Pin auf v6.0.0)

**Skill:** `.harness/skills/reviewer.md`
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**

- [slice-193](../plan/planning/done/slice-193-baseline-v600-bump.md)
- [welle-88](../plan/planning/welle-88-baseline-v600-migration.md)
- `AGENTS.md` §3.3, §3.5 (MR-013, ADR-Immutabilität)
- [`MR-021`](../../harness/conventions.md#mr-021),
  [`MR-051`](../../harness/conventions.md#mr-051),
  [`MR-055`](../../harness/conventions.md#mr-055)

---

## Findings

### F-1 — MR-060s eigene Direktiven-/Vorkommen-Zählung war ungenau

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Selbstauskunft)
- `pfad`: `harness/conventions/MR-060-baseline-v600.md`
- `befund`: MR-060 behauptete „11 lebende Direktiven … in 10 Dateien … 7
  unverändert bestätigt" und „acht reale Fundstellen" für den Tombstone —
  beide Zahlen unterzählt. Reale Messung: 15 Direktiven in 11 Dateien (11
  unverändert bestätigt, nicht 7); der Tombstone deckt 12 reale
  Fundstellen (4 `done/`-Slices, 2 Reviews, 2 frozen MRs, 3 immutable
  ADRs, 1 CR), nicht acht. Zusätzlich: die Übersetzungsfehler-Notiz nannte
  „vier Vorkommen" für den `../baseline/`-Fix in `reviewer.md`, real waren
  es fünf.
- `verifizierbar`: ja — `grep -rc "d-check:cite \.harness/baseline/v6\.0\.0"`
  über die betroffenen Dateien, bzw. `grep -rl "\.harness/baseline/v5\.18\.0"`
  über die fünf Tombstone-Wurzeln.
- `klasse`: „Gezählt-statt-geschätzt"-Anspruch ohne tatsächliche Nachzählung
  (dieselbe Familie wie `BEO-009`, Richtung b: Schluss reicht weiter als
  die tatsächliche Messung).

## Negativbefunde

- geprüft, ohne Befund: alle 15 `d-check:cite`-Direktiven lösen gegen
  `.harness/baseline/v6.0.0/` korrekt auf, inkl. der vier
  zeilenspann-verschobenen in `modul-05-planning-harness.md`
  (`make doc-check --enable citations` — 0 Befunde).
- geprüft, ohne Befund: `.harness/baseline/v5.18.0/` vollständig entfernt,
  `v6.0.0` einziger vendorter Baum, `make baseline-verify` grün.
- geprüft, ohne Befund: die drei immutable ADRs (0080/0081/0082) sind im
  Diff beider Commits unverändert — AGENTS.md §3.5 respektiert.
- geprüft, ohne Befund: `.d-check.yml`-Tombstone deckt exakt die reale
  Fundmenge über die fünf benannten Wurzeln.
- geprüft, ohne Befund: `.claude/rules/`-Aliase (vier Symlinks) lösen auf
  `v6.0.0` auf.
- geprüft, ohne Befund: MR-013-Commit-Trennung — `faac05c` (reiner Move +
  ein externer Pointer-Fix) und `8bd9ab7` (alles Übrige) sind sauber
  getrennt, beide einzeln `make gates`-grün.
- geprüft, ohne Befund: `docker run … d-check:latest` auf dem Endstand —
  497 Dateien, 0 Befunde.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Gezählt-statt-geschätzt-Anspruch ohne
tatsächliche Nachzählung

## Verdikt

**Merge-blockierend:** nein — MEDIUM, vor der Closure behoben (die
korrigierten Zahlen sind bereits in `harness/conventions/MR-060-baseline-v600.md`
und `.d-check.yml` eingearbeitet, Commit `8bd9ab7`).

**Übergabe:** Finding ging an den Implementer; die Finding-Klasse geht in
die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüfte der Verifier separat, in
eigenem Kontext.
