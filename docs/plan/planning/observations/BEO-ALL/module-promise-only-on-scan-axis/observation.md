# Eine Modul-Grenze wird nur auf der Quell-Achse gedacht

**Sub-Area:** `*`

Ein Modul prüft, was es scannt — liest aber Dateien, die es nie scannt:
Zieldateien außerhalb der Scan-Wurzeln, Post-Pässe über selbst benannte
Verzeichnisse, git-Revisionen. In diesen Eingaben gilt keine Zusage, die das
scannende Modul gibt, und die Folge kann still sein. Verkörpert seit dem
Vollständigkeits-Zensus in `AGENTS.md` §3.8; der Reviewer-Skill trägt dazu
den MEDIUM-Anker mit der Prüffrage: welche Eingaben liest dieses Modul, die
es nicht scannt — und gilt dort dieselbe Zusage?

## Benannt, nicht gezählt

Zwei Vorkommen ohne abgeschlossenen Vorgang (die Klasse „war dreimal
unvollständig") wurden bei der `v6.0.0`-Korrektur aus dem Zähler entfernt,
bleiben aber genannt, statt stillschweigend zu verschwinden.
