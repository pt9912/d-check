# Impl-Review R1 — slice-065 (ai-harness-Vorlage-Modulset)

**Datum:** 2026-07-06 · **Reviewer:** unabhängiger Subagent · **Verdikt: ACCEPT**
(2 Befunde, keine blockierend — beide eingearbeitet).

**Gegenstand:** Commits 4a90b13 (Move), 571ce58 (Spec-CR), eee472f (feat),
9567910 (release-prep) auf main; Basis 27168eb. DC-FA-CLI-006-CR / ADR-0033.

## Geprüft (real gelesen)

`internal/hexagon/core/app/suggest.go` (`renderHarness`, `renderHarnessPlanning`),
`internal/adapter/driving/cli/cli_acceptance_test.go` (Happy / PlanningAktiv /
VollKanon), `spec/lastenheft.md` §DC-FA-CLI-006, `spec/spezifikation.md`
§DC-FA-CLI-006.a; zur Verifikation `model/config.go`, `rules/{spans,hostpaths,codepaths,planning}.go`,
`rules/run.go`, `configyaml.go`.

## Ergebnis je Prüf-Frage

1. **Kopplung modules ↔ planning-Block: sauber.** `planning` ins Set bei
   `!repoAware || pathExists(roadmap)`; Block aktiv bei `!(repoAware && !pathExists(roadmap))`.
   Per De Morgan logisch identisch — Wahrheitstabelle über (repoAware, roadmapExists)
   deckt in allen vier Fällen. Kein Auseinanderfallen möglich (gleiche
   `harnessRoadmap`-Konstante).
2. **Decode-Vollständigkeit: sauber.** Kein emittierter Pfad scheitert an
   `configyaml.Decode` — planning-only (`roadmap:` genügt), spans/hostpaths
   konfigfrei (Default-Fallback), codepaths ohne `roots` dekodiert (unverändert
   ggü. Vor-Stand). Nuance: `ai-harness-init` mit planning aktiv gäbe im leeren
   Repo beim Lauf `planning-drift` — konsistent mit der „Zielbild/wächst-hinein"-
   Semantik (wie fehlende ids-`target`s), kein Defekt.
3. **Spec↔Code-Parität: sauber.** Fix (7) + planning bedingt (8.) + inaktiv via
   `--print-config` (7) + inaktiv via `--print-mk` (2) = **17** = `validModules()`.
   Keine Divergenz zwischen Renderer, Lastenheft-AK/Out-of-Scope und
   Spezifikations-Vorlage.
4. **Test-Härte: sauber.** Alle drei geprüften Mutationen (planning immer aktiv;
   spans/hostpaths nicht aufgenommen; Block immer aktiv) sterben an je einem
   benannten Test; `PlanningAktiv` sperrt die Gegenrichtung.
5. **Negativbefunde:** Logik-Kopplung, Decode, Enumeration, Determinismus,
   Read-only-Vertrag (DC-QA-03) — jeweils sauber.

## Befunde & Auflösung

- **F-1 (LOW) — Kommentar-Präzision:** `renderHarness`/`renderHarnessPlanning`
  begründeten die repo-bewusste planning-Behandlung mit „fail-closed (Exit 2 ohne
  Roadmap)". Real meldet das Modul zur Laufzeit `planning-drift` → **Exit 1**
  (Exit 2 ist Nutzungs-/Config-Fehlern vorbehalten). **Eingearbeitet:** Kommentare
  auf „planning-drift (Exit 1)" korrigiert.
- **F-2 (INFO) — 1:1-Parität:** der Renderer emittiert im `codepaths`-Block eine
  `# ignore-refs`-Kommentarzeile, die im kanonischen Spec-Beispiel fehlte —
  Widerspruch zum „1:1"-Anspruch der 0.39.0-Historie. **Eingearbeitet:** die
  `# ignore-refs`-Zeile in die Spezifikations-Vorlage aufgenommen.

Beide Befunde betreffen Kommentar-/Doku-Präzision, kein Laufzeitverhalten; die
technische Umsetzung ist korrekt und test-hart.
