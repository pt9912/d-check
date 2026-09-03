# Ein Liefer-Punkt landet im Commit eines fremden Slice

**Sub-Area:** `*`

Die Arbeit ist getan und geprüft, aber die Commit-Botschaft trägt eine
andere Kennung — `git log --grep=slice-<NNN>` findet den eigenen
Liefer-Punkt nicht mehr, und die Closure-Gegenprüfung sieht eine Zusage
ohne Commit. Unterschied zu [`BEO-ALL/commit-message-overclaims-work`](../commit-message-overclaims-work/observation.md):
dort behauptet die Botschaft zu viel, hier ist sie inhaltlich richtig und
nur am falschen Vorgang befestigt. Nicht reparierbar, sobald gepusht.
