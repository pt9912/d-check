# Slice slice-040: Planning-Konsistenz-Gate (Roadmap ↔ Slice-State)

**Status:** next (geplant — noch nicht in Arbeit).

**Welle:** welle-29-planning-consistency (Trigger: Nutzer-Audit 2026-06-21 —
`roadmap.md` sagte „Keine aktive Welle", während ein `slice-*` in
`in-progress/` lag; die Lifecycle-/Roadmap-Konsistenz ist dokumentiert, aber
nicht maschinell erzwungen).

**Bezug:** Lifecycle-Konvention (`AGENTS.md` §3.3 / Slice-Lifecycle
`open → next → in-progress → done`) und die Roadmap-Format-Regel
(`docs/plan/planning/in-progress/roadmap.md` §Aktuelle Welle); enforced als
Meta-Gate analog `make gate-consistency`. Read-only, deterministisch.

**Autor:** pt9912. **Datum:** 2026-06-21.

---

## 1. Ziel

Ein Meta-Gate `make planning-check` (in `make gates`), das die
**aspirativ-vs-bindend-Lücke** der Planning-Konsistenz schließt: die Roadmap
darf „Keine aktive Welle" **nur** behaupten, wenn kein `slice-*` in
`in-progress/` liegt — und umgekehrt. Dieselbe Klasse wie das
Traceability-Gate ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)):
eine dokumentierte Governance-Regel bekommt einen mechanischen Wächter.

## 2. Regel (maschinell)

- `hasSlices` := es existiert ≥1 `docs/plan/planning/in-progress/slice-*.md`.
- `hasActive` := der Abschnitt `## Aktuelle Welle` der `roadmap.md` enthält
  **nicht** den Marker „Keine aktive Welle".
- **Konsistent ⟺ `hasActive == hasSlices`.** Fehlschlag sonst (fail-closed,
  Exit ≠ 0) mit erklärender Meldung.

Damit wird zugleich das bisherige Closure-only-Update-Muster korrigiert: ein
Slice in `in-progress/` erzwingt eine benannte aktive Welle in der Roadmap.

## 3. Definition of Done (vorläufig)

Geplante Artefakte (Pfad gefenced, da noch nicht existent):

```text
tools/planning-consistency.sh   # Wächter-Skript: Regel + Negativ-Selbsttest
```

- [ ] Das Wächter-Skript implementiert die Regel; **Negativ-Selbsttest**
  (beide Inkonsistenz-Richtungen feuern nachweislich) wie
  `tools/gate-consistency.sh`.
- [ ] `make planning-check` ruft das Skript; in `make gates` aufgenommen
  (vor `record-gates`); `.PHONY` + `help`.
- [ ] Doku-Sync: `harness/README.md` §Sensors **und** `AGENTS.md` §4
  dokumentieren das Target (sonst rotes `make gate-consistency`).
- [ ] `make gates` grün (inkl. `planning-check` + `gate-consistency` über das
  neue Target) — Dogfooding: während slice-040 in-progress ist, trägt die
  Roadmap eine aktive Welle.
- [ ] Akzeptanz: Wegwerf-Fixture mit inkonsistentem Stand → Exit ≠ 0.
- [ ] Unabhängiges Review R1; Closure. Kein ADR nötig (additives Meta-Gate,
  keine neuen Import-Kanten; enforct eine bestehende Konvention).

## 4. Risiken / offene Punkte

- **Roadmap-Parsing** ist string-basiert (Marker „Keine aktive Welle"); ein
  umformulierter Marker bricht den Check (fail-closed) — der Marker wird als
  Konvention festgehalten.
- **Welle-Name-Abgleich** (in-progress-Slice-Welle == benannte aktive Welle)
  ist v1 **out of scope** — v1 prüft nur das Vorhandensein/Fehlen, nicht die
  Namens-Gleichheit (Folgepunkt, falls nötig).

## 5. Trigger

Nutzer-Audit 2026-06-21: dokumentierte, aber nicht erzwungene Roadmap↔Slice-
State-Konsistenz; Schwester-Lücke ADR-immutable → slice-041.

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Harness-Mechanik/Doku; Greenfield-Default;
`tools/harness/`-nahe Konvention).
