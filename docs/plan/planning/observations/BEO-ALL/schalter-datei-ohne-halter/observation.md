# Eine Schalter-Datei, deren Fehlen ein Gate lockert, hat keinen Halter

**Sub-Area:** `*`

Ein Wächter wird durch eine Datei im Arbeitsbaum streng geschaltet; ohne sie
gilt der lockere Zweig. Dann lockert jedes Entfernen der Datei das Gate — ein
`rm`, ein `git stash`, ein Branch-Wechsel, ein toter Symlink — ohne Signal, und
kein Gate prüft, dass die Datei da ist. Die Richtung des Ausfalls ist verkehrt:
ein verlorener Schalter soll das Strengere ergeben, nicht das Losere. Ableiter:
vor dem Anlegen eines Schalters fragen, was bei seinem **Fehlen** passiert, und
den Default auf die strenge Seite legen, den Schalter auf die lockere.
