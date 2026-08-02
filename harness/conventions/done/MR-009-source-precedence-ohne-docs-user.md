# MR-009 — Source Precedence ohne `docs/user`-Rang

- **Status:** Accepted
- **Datum:** 2026-06-11
- **Geltungsbereich:** [`harness/README.md` §Source precedence](../../README.md#source-precedence),
  [`AGENTS.md` §2](../../../AGENTS.md#2-kanonische-quellen-source-precedence)
- **Adaption:** Der Template-Default führt neun Ränge inkl. Rang 6
  `docs/user/*` (Operations, Quality, Releasing); d-check führt acht
  Ränge ohne `docs/user`, weil kein Operations-Doku-Stratum existiert
  (CLI-Tool vor dem ersten Release). Sichtbar geworden im
  Template-Vergleich (User-Review, 2026-06-11) — bis dahin eine
  stille Abweichung.
- **Begründung:** Ein Rang für nicht existierende Dateien wäre ein
  halluzinierter Eintrag (gleiche Klasse wie ein behauptetes Gate);
  die Rangordnung ist laut Baseline projektspezifische Wahl, die hier
  deklariert wird.
- **Auflösungs-Trigger:** welle-04 — mit der Release-Pipeline
  entsteht Betriebs-/Releasing-Doku; der `docs/user`-Rang wird dann
  eingefügt und dieser Eintrag als aufgelöst markiert. *(Eingetreten
  mit slice-011, siehe
  [`MR-010`](../../conventions.md#mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt).)*
