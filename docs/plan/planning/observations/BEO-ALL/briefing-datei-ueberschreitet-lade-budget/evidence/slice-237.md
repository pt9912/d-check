**Vorgang:** slice-237
**Fund:** Der zweite Sensor existiert jetzt auch auf Abschnitts-Ebene:
`structure[].max-lines`
([ADR-0089](../../../../../adr/0089-structure-max-lines-zwoelfte-bedingung.md))
ergänzt `file[].max-lines` (slice-236) um
ein Zeilenbudget je Abschnitt statt für die ganze Datei. Die Beobachtung
bleibt trotzdem offen — dieser Slice setzt bewusst **keine** Schwelle für
`AGENTS.md` (Auftraggeber-Entscheidung mit Vertragswirkung, `AGENTS.md`
§3.6, gehört in `slice-239`). Beide Sensoren existieren jetzt, ihre Nutzung
für die Briefing-Datei nicht.
