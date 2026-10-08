# Ein Wert hat mehrere Träger, und kein Gate hält sie gleich

**Sub-Area:** `*`

Eine Zusage (etwa eine Schwelle) steht in der Spezifikation, ihr Wert aber
zusätzlich im Build — als Default einer Make-Variable und als Build-Argument.
Ändert jemand einen Träger, bleiben die anderen still stehen: Die
Spezifikation sagt dann etwas anderes, als der Lauf prüft, und jedes Gate
bleibt grün. Ableiter: wer eine Zahl in die Spezifikation schreibt, zählt am
Code nach, wo dieselbe Zahl sonst noch steht, und benennt den Gleichhalter
oder seine Abwesenheit.
