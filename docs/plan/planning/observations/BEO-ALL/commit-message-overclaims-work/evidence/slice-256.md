**Vorgang:** slice-256
**Fund:** Richtung (b), zweimal, beide vom unabhängigen Review gefunden. Die
Botschaft des Release-Commits nannte „fehlende Referenz liefert leer
(UNGEPRUEFT-Pfad)"; unter `set -euo pipefail` brach der Schritt über den
ERR-Trap vorher ab (R1-F-3, nachgemessen). Die Botschaft des R1-Fix-Commits
schrieb, ein zweites Ziehen „derselben Referenz" scheitere — gemessen trifft
das nur Digest-Referenzen, nicht Tags (R2-2). Beide Male stand die Aussage
aus der Vorstellung des Codes, nicht aus einem Lauf gegen ihn.
