# Planning

Slice-Pläne und Roadmap dieses Repos.

## Lifecycle

`open/` → `next/` → `in-progress/` → `done/`

- Die Status-Bewegung eines Slice ist eine **reine Datei-Bewegung**
  (`git mv`), keine inhaltliche Änderung im selben Commit
  (siehe `AGENTS.md` §3.3).
- Slice-Dateinamen: `slice-<NNN>-<kurzer-titel>.md`.
- Nach Abschluss erhält der Slice in `done/` seine Closure-Notiz
  (§7 des Slice-Plans) mit dem Commit-Hash der Umsetzung. Der
  **Steering-Loop-Eintrag** der Closure-Notiz verweist auf die kanonische
  Definition im vendorten Regelwerk:
  [`grundlagen-klassifikation.md` §Steering Loop](../../../.harness/baseline/v5.0.0/regelwerk/grundlagen-klassifikation.md#steering-loop).
- Die aktuelle Welle und die Wellen-Reihenfolge stehen in
  [`in-progress/roadmap.md`](in-progress/roadmap.md).

## Carveouts

Dokumentierte Ausnahmen (z. B. dauerhaft rote Gates, temporäre
Aufweichungen) leben in [`../carveouts/`](../carveouts/) als
`CO-<NNN>-<titel>.md` — immer mit Auflösungs-Trigger oder expliziter
Permanenz-Begründung und Folge-Slice.
