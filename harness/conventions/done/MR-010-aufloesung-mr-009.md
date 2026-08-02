# MR-010 — Auflösung von MR-009: `docs/user`-Rang eingefügt

- **Status:** Accepted
- **Aufgelöst durch:** Baseline-Stand v5.0.0 (die 9-Rang-Default-Liste inkl. docs/user ist Baseline-konform)
- **Datum:** 2026-06-11
- **Geltungsbereich:** [`harness/README.md` §Source precedence](../../README.md#source-precedence),
  [`AGENTS.md` §2](../../../AGENTS.md#2-kanonische-quellen-source-precedence),
  `docs/user/`
- **Adaption:** Der Auflösungs-Trigger von
  [`MR-009`](../../conventions.md#mr-009--source-precedence-ohne-docsuser-rang) ist
  eingetreten: mit der GHCR-Release-Pipeline (slice-011) existiert
  Betriebs-/Releasing-Doku (`docs/user/releasing.md`,
  `docs/user/operations.md`). Der `docs/user`-Rang ist als Rang 6 in
  beide Source-Precedence-Tabellen eingefügt (Template-Default
  wiederhergestellt, neun Ränge); die nachfolgenden Ränge rücken um
  eins.
- **Begründung:** Baseline-Konformität, sobald die Dateien real
  existieren — kein Rang für Phantome, kein Phantom für Ränge.
- **Auflösungs-Trigger:** permanent (Baseline-Konformität).
