# MR-017 — Lokale Baseline-Lese-Form (Cache) aus dem Selbst-Scan ausgenommen

- **Status:** Accepted
- **Aufgelöst durch:** MR-019 (Regelwerk committet vendored statt gecacht)
- **Datum:** 2026-06-25
- **Geltungsbereich:** [`.d-check.yml`](../../../.d-check.yml) `scan.ignore`,
  das Materialisierungs-Skript `tools/harness/fetch-baseline-cache.sh`,
  [§Adoptierte Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen),
  [`AGENTS.md`](../../../AGENTS.md) §1, der gitignorierte Pfad `.harness/cache/`
- **Adaption:** Die in [`AGENTS.md`](../../../AGENTS.md) §1 vorgesehene Lese-Form der
  adoptierten Baseline (Bundle „herunterladen, entpacken, nur den benötigten
  Abschnitt laden") wird lokal **materialisiert** nach dem Pfadschema
  `.harness/cache/<tag>/regelwerk/` (entpacktes `lab-regelwerk.zip`) und
  `.harness/cache/<tag>/templates/` (entpacktes `lab-templates.zip`); aktueller
  `<tag>` = `v1.4.0`; materialisiert wird reproduzierbar per
  `tools/harness/fetch-baseline-cache.sh` (zieht die beiden Release-Assets und
  entpackt; Tag ohne Argument aus dieser §Baseline-Stand-Zeile abgeleitet —
  kein Drift). Der Cache ist **gitignored** ([`.gitignore`](../../../.gitignore)
  `.harness/cache/`) — ephemer, kein Repo-Vertrag — und wird über
  `scan.ignore: [".harness/cache/**"]` in [`.d-check.yml`](../../../.d-check.yml) aus
  dem Dogfooding-Selbst-Scan **ausgenommen**. Grund: der Cache trägt
  Fremdinhalt (die Kurs-Docs referenzieren ihre eigenen `ADR-`/`MR-`-IDs und
  Modulpfade, die in *diesem* Repo nicht existieren); ohne Ausnahme meldete
  `make doc-check` sie als `id-unlinked`/`codepath-missing`. Selbe Klasse wie
  die eingebauten `SKIP_DIRS` (`.git`, `vendor`, `node_modules` —
  Fremd-/Generiertes), daher **keine** Gate-Lockerung im Sinne von
  [`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden):
  ausgenommen wird Nicht-Repo-Inhalt, keine Repo-Doku verliert Deckung.
- **Begründung:** Sichtbar geworden bei der v1.4.0-Adoption
  ([`MR-016`](../../conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011)), als der
  Cache erstmals befüllt wurde (Nutzer-Auftrag 2026-06-25). Die Konvention ist
  tag-generisch (`<tag>`), nicht v1.4.0-spezifisch — darum ein eigener Eintrag
  statt Bündelung in den Pin-Bump-Nachtrag.
- **Auflösungs-Trigger:** permanent, solange die Baseline-Lese-Form lokal
  gecacht wird.
