# Carveouts

Dokumentierte Ausnahmen gegenüber den eigenen Regeln (z. B. dauerhaft
rote Gates, temporär gelockerte Schwellen). Format:
`CO-<NNN>-<kurzer-titel>.md` mit Pflichtfeldern Auflösungs-Trigger
(oder expliziter Permanenz-Begründung) und Folge-Slice.

Aktuell: **keine aktiven Carveouts, zwei aufgelöst.** Ein aufgelöster
Carveout wandert per `git mv` nach `carveouts/done/` — derselbe
Verzeichnis-Zustand wie beim Slice-Lifecycle, kein Status-Feld hier im
Index.

## Aktive Carveouts

Keine.

## Aufgelöste Carveouts

| ID | Titel | Betroffenes Gate | aufgelöst durch |
| --- | --- | --- | --- |
| [CO-001](done/CO-001-vcs-range-stiller-skip.md) | `vcs --range` überspringt still, solange das publizierte Image den Fix nicht trägt | `adr-check` | Klassen-Fix `slice-220`, Release `v0.76.1` (`slice-219`) |
| [CO-002](done/CO-002-slice-221-gegenstand-entfallen.md) | `slice-221` schließt ohne Lieferung — der DoD-Wächter kannte den `Gegenstand:`-Ausgang noch nicht | `verify-closure-notes` | `open-tasks-require-marker`/`-section` ([ADR-0085](../adr/0085-bedingte-pflicht-marke-open-tasks.md)), Folge-Slice `slice-225` |
