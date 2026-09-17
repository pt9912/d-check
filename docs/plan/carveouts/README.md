# Carveouts

Dokumentierte Ausnahmen gegenüber den eigenen Regeln (z. B. dauerhaft
rote Gates, temporär gelockerte Schwellen). Format:
`CO-<NNN>-<kurzer-titel>.md` mit Pflichtfeldern Auflösungs-Trigger
(oder expliziter Permanenz-Begründung) und Folge-Slice.

Aktuell: **zwei Carveouts.**

| ID | Titel | Betroffenes Gate | Auflösungs-Trigger |
| --- | --- | --- | --- |
| [CO-001](CO-001-vcs-range-stiller-skip.md) | `vcs --range` überspringt still, solange das publizierte Image den Fix nicht trägt | `adr-check` | nächstes Release veröffentlicht und `DCHECK_IMAGE` zeigt darauf |
| [CO-002](CO-002-slice-221-gegenstand-entfallen.md) | `slice-221` schließt ohne Lieferung — der DoD-Wächter kennt den `Gegenstand:`-Ausgang noch nicht | `verify-closure-notes` | [slice-225](../planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md) geschlossen |
