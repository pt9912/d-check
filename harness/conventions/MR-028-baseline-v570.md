# MR-028 — Baseline-Pin-Hebung auf v5.7.0 (vierter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-21
- **Geltungsbereich:** [§Baseline](../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../AGENTS.md), [`harness/README.md`](../README.md),
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), den
  Planning-Docs und den aktiven `MR-*`-Dateien
- **Adaption:** Der Baseline-Pin ist von `v5.6.0` auf **`v5.7.0`** gehoben
  (Kurs-Tag 2026-08-21, Kurs-Welle 81 „Zwei Hälften, ein Wächter") — die von
  [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, vierter Nachtrag der Serie; ersetzt
  [`MR-026`](done/MR-026-baseline-v560.md) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle
  (`lab-regelwerk.zip`, beide Bäume, `SHA256SUMS`, 51 Dateien — gemessen,
  die team-sim-Erweiterungen der Kurs-Welle liegen außerhalb des Bundles),
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema
  `.harness/baseline/<tag>/{regelwerk,templates}/`. Das Delta ist **eine
  Stufe** (drei Regelwerks-Dateien, +5/−3 Zeilen): `modul-06-roadmap.md`
  fasst §Offene Wellen als zwei unabhängige Aussagen (Ruhe-Marker
  **zusätzlich** zur Liste; gewächtert nur die Marker-Hälfte),
  `modul-10-review-harness.md` führt `klasse` als sechstes Output-Feld,
  dazu der README-Stand.
  - **Alter Baum entfernt, Historie via Tombstone.** `.harness/baseline/v5.6.0/`
    ist entfernt (ein Pin, eine netzlose Lese-Form — dieselbe Entscheidung
    wie bei den Vorgänger-Ablösungen). Die **eingefrorenen** Verweise darauf
    — der `done/`-Slice der welle-78-Etappe C und die aufgelöste
    [`MR-024`](done/MR-024-aktuelle-welle-ruhe-marker-form.md) — sind über
    das geteilte Referenz-Ventil `ignore-refs`
    ([`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
    quell-skopiert ausgenommen; **lebende** Verweise sind pin-gebunden
    retargetet
    ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)).
  - **Der Konformitäts-Abgleich ist nicht Teil dieser MR:** das
    Ein-Stufen-Delta wird je geänderter Regel auditiert (Delta-Audit in
    slice-110); die eine Produkt-Folge — `planning.waves.mode` nach dem
    Konsumenten-CR — ist ein eigener Slice mit eigenem Release-Punkt.
    Diese MR hebt den Pin, sie behauptet keine Konformität.
- **Begründung:** Auftraggeber-Anstoß 2026-08-21 („Jetzt gibt es v5.7.0");
  der Baseline-Default sticht die repo-lokale Adaption. Vendored wird das
  **Release-Asset am Tag** (`--check-latest` ist die
  Currency-/Authentizitäts-Gegenprobe), nicht der Kurs-Arbeitsbaum.
- **Auflösungs-Trigger:** die nächste Pin-Hebung ersetzt diesen Eintrag durch
  ihren Nachfolger — wie
  [`MR-026`](done/MR-026-baseline-v560.md) durch diesen Eintrag ersetzt
  wurde.
  [`MR-023`](../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  bleibt daneben **aktiv** stehen: es trägt das Bundle-Layout (dort
  „permanent"), nicht den Pin-Wert.
