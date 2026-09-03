# Eine in Lastenheft/Spezifikation zugesagte Randbedingung hat keinen eigenen Test

**Sub-Area:** `*`

Eine Regel verspricht in Prosa mehr, als ihre Test-Datei belegt: der Code
setzt die Zusage korrekt um, aber kein Test dreht genau diese Bedingung um —
eine künftige Implementierungsänderung könnte sie stillschweigend brechen,
ohne dass ein Test rot würde. Unterschied zu [`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../wortlaut-behauptet-pruefung-die-fehlt/observation.md):
dort behauptet ein vorhandener Test eine Prüfung, die er nicht leistet; hier
fehlt der Test ganz.
