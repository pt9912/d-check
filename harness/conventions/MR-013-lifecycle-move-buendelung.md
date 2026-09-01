# MR-013 — Lifecycle-Move-Commit bündelt gekoppelte Verweise

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`modul-05-planning-harness.md` §Lifecycle als State Machine](../../.harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine)
- **Datum:** 2026-06-21
- **Geltungsbereich:** [`AGENTS.md` §3.3](../../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits),
  der Slice-Lifecycle in **beide** gekoppelten Richtungen — `docs/plan/planning/open/`
  → `…/in-progress/` (Beanspruchung) und `…/in-progress/` → `…/done/` (Closure) —,
  `make planning-check`; die MR-/Wellen-Lifecycle-Moves
  (`conventions/` → `conventions/done/`, flaches Wellendokument → `done/`)
- **Adaption:** Erste Adaption einer **Lifecycle-Regel** (Nachtrag zu
  [`MR-000`](../conventions.md#mr-000--baseline-aussage), das „keine Adaptionen … für
  Lifecycle-Regeln" vermerkt; der Baseline-Aussage-Eintrag bleibt unverändert). Das Baseline-§3.3
  trennt reinen Move (Commit 1, `R100`-Rename) vom Inhalt (Commit 2). Seit
  `make planning-check` (slice-040) den Roadmap-Zustand **atomar** an den
  in-progress-Stand koppelt (kein `slice-*` in `…/in-progress/` ⟺ Roadmap
  §Offene Wellen trägt den Ruhe-Marker), wäre ein *byte*-reiner
  Move-Commit beim Lifecycle-Move zwangsläufig gate-rot (leeres
  `in-progress/` bei noch aktiver Roadmap; zusätzlich `target-missing` auf
  jeden Verweis, der den Slice über seinen `in-progress/`-Pfad verlinkt).
  Adaption: der `git mv`-Commit dieses Moves trägt **zusätzlich** (a) den
  Roadmap-Flip §Offene Wellen und (b) alle Pfad-Verweise auf den Slice
  (Roadmap, `AGENTS.md` §4, `harness/README.md` §Sensors) von
  `in-progress/` nach `done/`. Der **Slice-Body** (DoD-Haken +
  Closure-Notiz; neue Slices ohne Status-Zeile [slice-091/D-5], historische
  `done/`-Slices behalten ihre) bleibt im zweiten Commit; weil die Slice-Datei im
  Move-Commit unverändert ist, hält die Rename-Detection (`R100`) und damit
  die `git log --follow`-Begründung des Baseline-§3.3.
  - **Dieselbe Bündelung gilt für die Beanspruchung** `open/` → `in-progress/`
    (seit slice-185), nur mit umgekehrtem Vorzeichen: dort **verlässt** der
    Ruhe-Marker §Offene Wellen, statt zurückzukehren, und die Pfad-Verweise auf
    den Slice wandern von `open/` nach `in-progress/`. `make planning-check`
    hält den Marker in **beide** Richtungen, ein byte-reiner Move-Commit ist
    deshalb auch hier gate-rot — gemessen: Marker stehengelassen, Slice in
    `in-progress/`, Ergebnis `planning-drift` („Slice(s) in
    `docs/plan/planning/in-progress`, aber die Roadmap-Sektion … trägt den
    Ruhe-Marker") und Exit 2. **Das Feld nannte bis dahin nur die
    Closure-Richtung**, während die Übersichtstabelle in
    [`conventions.md`](../conventions.md#mr-013) längst
    „Slice-Lifecycle-Moves" führte — der Widerspruch lag zwischen zwei Stellen
    desselben Dokuments, und die gelebte Praxis folgte der Tabelle.
  - **Dieselbe Bündelung gilt für MR- und Wellen-Lifecycle-Moves** (seit
    welle-79): eine nach `conventions/done/` bzw. `done/` wandernde Datei
    trägt **relative** Verweise, die vom neuen Ort eine Ebene tiefer
    auflösen müssen — ein byte-reiner Move-Commit wäre `doc-check`-rot.
    Der Move-Commit trägt deshalb die **Link-Tiefen-Fixes der bewegten
    Datei selbst** mit; alles Übrige bleibt Commit 2. Drückt der
    Fix-Umfang den Rename-Score Richtung 50 %, deklariert die
    Commit-Botschaft den Move ausdrücklich als `git mv` — die Botschaft
    ersetzt dann, was `git log --follow` nicht mehr sicher zeigt.
- **Begründung:** Sichtbar 2026-06-21 — die PR-/Push-CI prüft den Push-Tip,
  der ein Zwischen-Commit sein kann; sie lief auf dem reinen Move-Commit von
  slice-040 rot (`target-missing` + `planning-check`). Die
  Per-Commit-Grün-Regel (grün = Boden, nicht Decke) und die
  Rename-Detection schließen sich nur scheinbar aus: die Kopplung betrifft
  **fremde** Dateien, nicht den Slice-Body. slice-040 führte
  `make planning-check` ein und löste damit die Kollision aus.
- **Auflösungs-Trigger:** permanent, solange `make planning-check` den
  Roadmap-↔-in-progress-Invariant erzwingt.
