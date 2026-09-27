# Commits/PRs nennen mindestens eine Kennung

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Commits/PRs müssen mindestens eine `DC-*`-, `ADR-*`-, `MR-*`- oder
`slice-*`-ID nennen (maschinell erzwungen: `make trace-check` /
`commit-msg`-Hook / PR-CI — über das Modul `commits`, dogfooded über das
eigene Image,
[ADR-0027](../../../docs/plan/adr/0027-commits-traceability-modul.md).
Ausnahme: Merge-/Revert-Commits). Vergeben werden IDs nur beim
Spec-/ADR-Schreiben nach dem deklarierten Schema
([`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage))
— nie ad hoc im Commit/PR; Agenten referenzieren IDs, sie erfinden keine.
Struktur-IDs (`SPEC-<NNN>`/`ARC-<NNN>`,
[`MR-000`](../../conventions.md#mr-000--baseline-aussage)) entstehen nur
beim Schreiben der Spec-Straten — fortlaufend je Datei — und gehören
**nicht** in Commit-Botschaften.
