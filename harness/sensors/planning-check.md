# `make planning-check` — hält den Ruhe-Marker der Roadmap gegen das Lifecycle-Verzeichnis

## Vertrag

Via Modul `planning` (Image, dogfood). Die Roadmap
(`docs/plan/planning/in-progress/roadmap.md`) §Offene Wellen darf den
Ruhe-Marker „Nichts in Arbeit" **genau dann** tragen, wenn kein `slice-*` in
`docs/plan/planning/in-progress/` liegt (`planning-drift`).

**Beide Richtungen zählen:** ein fehlender Marker bei leerem `in-progress/` und
ein stehengebliebener Marker bei beanspruchtem Slice sind derselbe Defekt. Das
ist die Kopplung, die [`MR-013`](../conventions/MR-013-lifecycle-move-buendelung.md)
atomar hält — ein byte-reiner Lifecycle-Move ist deshalb gate-rot.

## Grenze — was das Grün nicht abdeckt

1. **Er prüft die Marker-Hälfte, nicht die Listen-Hälfte** — die Bijektion
   „Zeiger ⟺ Welle-Datei" misst `wave-drift` unter `mode: many`. Wer die eine
   für die andere hält, hält einen halben Wächter für einen ganzen.
2. **Hermetisch — kein git.** Der Zustand ist das Verzeichnis im Arbeitsbaum,
   nicht der gemergte Stand.
3. **Der Negativ-Selbsttest lebt als Akzeptanztest im Modul** (`make test`).

## Bindung

Bestandteil von `make gates`.
[ADR-0028](../../docs/plan/adr/0028-planning-lifecycle-modul.md) ·
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) ·
[`AGENTS.md`](../../AGENTS.md) §3.3
