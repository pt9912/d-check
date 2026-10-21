**Vorgang:** slice-278
**Fund:** Zwei Zahlen der Messungs-Tabelle waren aus dem falschen Lauf oder mit
falscher Menge entstanden: „42 Befunde" bei `match: name` (aus dem Lauf ohne Reports
übernommen, nie gezählt; gezählt sind 19) und „22 lebende Slices" (276 − 254, wobei die
254 aus einem `grep -rl` über **alle** Dateien unter `done/` kamen, nicht nur über
`slice-*.md`; richtig sind 47 und 229 Stubs). Der Review fand beides durch Nachzählen
in einer Kopie; die Entscheidungen blieben richtig, die Zahlen standen aber im
Sensor-Vertrag.
