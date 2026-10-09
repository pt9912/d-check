# Eine Zuordnung über den Namen deckt auch den Slice, dessen Name Präfix eines anderen ist

**Sub-Area:** `*`

`reviews.match: name` ordnet einen Report einem Slice zu, wenn der Dateiname
des Reports den Basisnamen des Slice enthält, gefolgt von einem Zeichen, das
weder Buchstabe noch Ziffer ist. Ist ein Basisname über einen Bindestrich,
Unterstrich oder Punkt Präfix eines anderen (`slice-a-foo` und
`slice-a-foo-bar`), deckt der Report des längeren auch den kürzeren — ein
fehlender Report fällt dann still nicht auf. Ableiter: bei einer
Namens-Zuordnung die Menge der Basisnamen auf Präfix-Paare zählen, bevor man
einer grünen Deckung traut.
