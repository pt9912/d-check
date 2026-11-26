# `make planning-check` — hält den Ruhe-Marker der Roadmap gegen das Lifecycle-Verzeichnis

## Vertrag

Via Modul `planning` (Image, dogfood). Die Roadmap
(`docs/plan/planning/in-progress/roadmap.md`) §Offene Wellen darf den
Ruhe-Marker „Nichts in Arbeit" **genau dann** tragen, wenn kein `slice-*` in
`docs/plan/planning/in-progress/` liegt (`planning-drift`).

**Beide Richtungen zählen:** ein fehlender Marker bei leerem `in-progress/` und
ein stehengebliebener Marker bei beanspruchtem Slice sind derselbe Defekt.
**Fail-closed via Heading-Guard:** Fehlt die §Offene-Wellen-Überschrift oder
ist sie umbenannt, meldet der Lauf — er passiert nicht still grün. Das
ist die Kopplung, die [`MR-013`](../conventions/done/MR-013-lifecycle-move-buendelung.md)
atomar hält — ein byte-reiner Lifecycle-Move ist deshalb gate-rot.

## Grenze — was das Grün nicht abdeckt

1. **Er prüft die Marker-Hälfte, nicht die Listen-Hälfte** — die Bijektion
   „Zeiger ⟺ Welle-Datei" misst `wave-drift` unter `mode: many`. Wer die eine
   für die andere hält, hält einen halben Wächter für einen ganzen.
2. **Hermetisch — kein git.** Der Zustand ist das Verzeichnis im Arbeitsbaum,
   nicht der gemergte Stand.
3. **Er liest EIN Wellen-Präfix: `wave-`.** `planning.waves.glob` und
   `results-glob` stehen auf `wave-*`
   ([`MR-082`](../conventions.md#mr-082)). Die Bestands-Wellen `welle-<NN>` sind
   für `wave-drift` und `wave-results-missing` **ungewächtert** — ihre
   Register-Zeilen (`grep -c '^| welle-' docs/plan/planning/in-progress/roadmap.md`)
   und ihre Ergebnisnotizen (`git ls-files | grep -c 'done/welle-[0-9]*-results.md$'`)
   melden keinen Befund, auch wo einer stünde. Ob das Modul mehrere Präfixe liest,
   ist ein Produkt-Merkmal und gehört dem Auftraggeber; die Senkung trägt
   [`ADR-0117`](../../docs/plan/adr/0117-wellen-namen-ueber-die-konfiguration.md).
## Bindung

Bestandteil von `make gates`.
[ADR-0028](../../docs/plan/adr/0028-planning-lifecycle-modul.md) ·
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) ·
[`AGENTS.md`](../../AGENTS.md) §3.3
