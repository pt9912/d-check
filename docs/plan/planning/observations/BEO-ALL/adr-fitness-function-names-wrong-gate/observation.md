# Eine Fitness-Function-Zeile nennt ein Gate, das die behauptete Eigenschaft nicht prüft

**Sub-Area:** `*`

Eine ADR-Fitness-Function-Tabelle nennt ein reales, existierendes
Make-Target als Beleg für eine Eigenschaft, die dieses Target tatsächlich
nicht prüft (hier: `make doc-check` als Beleg für „YAML-Config-Beispiel
in der Doku bleibt syntaktisch gültig" — `doc-check` dekodiert kein
gefenctes YAML, und der einzige Mechanismus, der das täte, deckt die
betroffene Datei nicht ab). Das Gate existiert und läuft grün — es prüft
nur nicht das Behauptete.
