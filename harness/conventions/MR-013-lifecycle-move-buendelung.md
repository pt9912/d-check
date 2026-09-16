# MR-013 — Lifecycle-Move-Commit bündelt gekoppelte Verweise

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. `AGENTS.template.md` §3.3 selbst deckt
  seit `v6.6.0` bereits zwei Fälle (Regelfall · Lifecycle-Übergang nach
  `done/`, Inhalt vor reinem `git mv`) und damit die **Slice**-Lifecycle-Moves
  (Beanspruchung, Closure) vollständig ab — dort ändert sich beim Move nur
  die bewegte Slice-Datei nicht, andere Dateien (Roadmap-Ruhe-Marker,
  Pfad-Verweise) dürfen mitreisen, reine Git-Semantik ohne Kanon-Bedarf. Die
  ursprüngliche Aussage dieses Eintrags nannte `modul-05-planning-harness.md`
  — das war die **falsche** Stelle, wie
  [slice-223](../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md)
  beim Auflösen der Slice-Lifecycle-Hälfte gezeigt hat; die tatsächlich
  einschlägige Kanon-Stelle ist
  [`grundlagen-traceability.md` §Herkunfts-Anker](../../.harness/baseline/v6.9.0/regelwerk/grundlagen-traceability.md#herkunfts-anker):
  *„Der `git mv` zieht die Pfad-Berichtigung nach sich, als eigener Commit
  nach dem Umzug."*
- **Datum:** 2026-06-21 (getrimmt 2026-09-16,
  [slice-223](../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md))
- **Geltungsbereich:** [`AGENTS.md` §3.3](../../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits),
  ausschließlich die MR-/Wellen-Lifecycle-Moves
  (`conventions/` → `conventions/done/`, flaches Wellendokument → `done/`).
  **Nicht mehr** Geltungsbereich: die Slice-Lifecycle-Moves selbst
  (Beanspruchung, Closure) — die deckt der Kanon jetzt direkt, siehe Feld
  `Ersetzt-Baseline-Regel`.
- **Adaption:** Eine nach `conventions/done/` bzw. `done/` wandernde Datei
  trägt **relative** Verweise, die vom neuen Ort eine Ebene tiefer auflösen
  müssen — ein byte-reiner Move-Commit wäre `doc-check`-rot. Der Kanon selbst
  kennt für genau diesen Fall eine Antwort (§Ersetzt-Baseline-Regel oben):
  Move und Korrektur als **zwei** Commits, beide im selben Push — der
  Zwischenstand ist zulässig, solange er nicht die Spitze eines Push wird.
  **Warum dieser Weg lokal trotzdem nicht gangbar ist:** Der lokale
  `pre-commit`-Hook fährt `doc-check` auf **jedem** Commit, nicht nur am
  Push-Tip ([ADR-0013](../../docs/plan/adr/0013-pr-ci-und-traceability-gate.md)/[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md)-Kette).
  Ein reiner Move-Commit hätte die bewegte Datei mit **alten**,
  jetzt falschen Link-Tiefen — der Hook lehnt genau diesen Commit lokal ab,
  bevor der zweite je entsteht. Adaption: Der Move-Commit trägt deshalb die
  **Link-Tiefen-Fixes der bewegten Datei selbst** mit; alles Übrige bleibt
  Commit 2. Drückt der Fix-Umfang den Rename-Score Richtung 50 %, deklariert
  die Commit-Botschaft den Move ausdrücklich als `git mv` — die Botschaft
  ersetzt dann, was `git log --follow` nicht mehr sicher zeigt.

  **Gemessen, nicht am Vorgänger vorbeigelesen** — eine Probe mit dieser
  getrimmten Fassung selbst
  ([slice-223](../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md)):
  fünf `MR`-Einträge (059/061/062/063/064) wanderten in genau dieser
  Zwei-Halbschritt-Form nach `conventions/done/` — ein Commit mit Index-
  Zeilen-Verschiebung plus dieser Datei, dann der reine `git mv` der fünf
  Dateien mit ihren Link-Tiefen-Fixes im selben Commit.
- **Begründung:** Sichtbar 2026-06-21 — die PR-/Push-CI prüft den Push-Tip,
  der ein Zwischen-Commit sein kann; sie lief auf dem reinen Move-Commit von
  slice-040 rot (`target-missing` + `planning-check`). Für die
  **Slice**-Lifecycle-Hälfte war das ein Missverständnis der eigenen
  Git-Semantik, nicht ein echter Kanon-Konflikt: Die Kopplung betrifft
  **fremde** Dateien (Roadmap, Sensoren-Index), nicht den Slice-Body selbst,
  und die Rename-Detection bleibt davon unberührt — genau deshalb konnte
  [slice-223](../../docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md)
  diese Hälfte auflösen. Für die **MR-/Wellen**-Hälfte bleibt der Konflikt
  real: dort ändert sich die bewegte Datei selbst, und der lokale Hook lässt
  den kanonischen Zwei-Commit-Weg nicht zu.
- **Auflösungs-Trigger:** permanent, solange der lokale `pre-commit`-Hook
  `doc-check` auf jedem Commit statt nur am Push-Tip fährt.
