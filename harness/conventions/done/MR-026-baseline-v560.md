# MR-026 — Baseline-Pin-Hebung auf v5.6.0 (dritter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-21
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md),
  [`.harness/skills/reviewer.md`](../../../.harness/skills/reviewer.md), den
  Planning-Docs und den aktiven `MR-*`-Dateien
- **Adaption:** Der Baseline-Pin ist von `v5.0.0` auf **`v5.6.0`** gehoben
  (Kurs-Tag 2026-08-16) — die von
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, dritter Nachtrag nach
  [`MR-012`](../../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)
  und
  [`MR-016`](../../conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011)
  (Zählung der Titel-Serie; die v5.0.0-Hebung lief unter dem eigenen Titel
  des Layout-Eintrags außerhalb dieser Serie). Anders als bei
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  ändert sich **kein** Layout: dasselbe self-contained Bundle
  (`lab-regelwerk.zip`, beide Bäume, `SHA256SUMS`), dasselbe
  Materialisierungs-Skript, unverändertes Pfadschema
  `.harness/baseline/<tag>/{regelwerk,templates}/`. Das Delta ist rein
  **additiv** (sechs Stufen v5.1.0–v5.6.0, 20 Regelwerks-Dateien,
  +902/−152 Zeilen; u. a. §Vergabe, Straten-IDs/Reconciliation,
  Kommentar-Regel, zwei Korrektur-Stufen, Team-Fähigkeit, TA-7-Wirkung).
  - **Alter Baum entfernt, Historie via Tombstone.** `.harness/baseline/v5.0.0/`
    ist entfernt (ein Pin, eine netzlose Lese-Form — dieselbe Entscheidung wie
    bei der v1.4.0-Ablösung). Die **eingefrorenen** Verweise darauf — die
    immutablen `Accepted`-ADRs
    [`ADR-0047`](../../../docs/plan/adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
    und
    [`ADR-0048`](../../../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md)
    sowie drei Review-Reports (Lauf-Belege) — sind über das geteilte
    Referenz-Ventil `ignore-refs`
    ([`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
    quell-skopiert ausgenommen; **lebende** Verweise sind pin-gebunden
    retargetet ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)).
  - **Der Konformitäts-Abgleich ist nicht Teil dieser MR:** das Stufen-Delta
    wird je Regel auditiert (Etappe B der Migrations-Welle), Anpassungen
    folgen als eigene Slices/Adaptionen — diese MR hebt den Pin, sie
    behauptet keine Konformität.
- **Begründung:** Auftraggeber-Vorgabe (angekündigt 2026-08-16, Freigabe
  2026-08-21): Migration auf den neuesten Kurs-Stand; der Baseline-Default
  sticht die repo-lokale Adaption. Vendored wird das **Release-Asset am Tag**
  (`--check-latest` ist die Currency-/Authentizitäts-Gegenprobe), nicht der
  Kurs-Arbeitsbaum.
- **Auflösungs-Trigger:** die nächste Pin-Hebung ersetzt diesen Eintrag durch
  ihren Nachfolger — wie
  [`MR-012`](../../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)
  und
  [`MR-016`](../../conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011)
  durch ihre Nachfolger ersetzt wurden.
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  bleibt daneben **aktiv** stehen: es trägt das Bundle-Layout (dort
  „permanent"), nicht den Pin-Wert.
