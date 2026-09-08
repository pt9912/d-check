# Carveouts

Dokumentierte Ausnahmen gegenüber den eigenen Regeln (z. B. dauerhaft
rote Gates, temporär gelockerte Schwellen). Format:
`CO-<NNN>-<kurzer-titel>.md` mit Pflichtfeldern Auflösungs-Trigger
(oder expliziter Permanenz-Begründung) und Folge-Slice.

Aktuell: **ein Carveout.**

| ID | Titel | Betroffenes Gate | Auflösungs-Trigger |
| --- | --- | --- | --- |
| [CO-001](CO-001-vcs-range-stiller-skip.md) | `vcs --range` überspringt still, solange das publizierte Image den Fix nicht trägt | `adr-check` | nächstes Release veröffentlicht und `DCHECK_IMAGE` zeigt darauf |
