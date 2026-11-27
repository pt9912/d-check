# Ein Closure-Trigger misst etwas, das ein späterer Slice derselben Welle verändert

**Sub-Area:** `*`

Ein Wellen-Closure-Trigger wird bei der Eröffnung als Messung festgeschrieben: eine Menge, die „genau“ einer an anderer
Stelle geführten Liste gleichen soll, oder ein Verweis auf „das Kommando“ eines Slice-Plans. Danach erzeugt ein Slice der
Welle ein Artefakt, das in die Menge fällt, oder ändert den Plan, sodass dort zwei Kommandos stehen. Der Trigger ist dann
nach seinem Wortlaut unerfüllbar oder mehrdeutig, und sichtbar wird das erst bei der Closure oder durch eine Rolle, die
den Trigger gegen den Zwischenstand misst. Ableiter: der Trigger nennt Kommando und erwarteten Wert selbst; wo er
Mengen-Gleichheit verlangt, nennt er die Klassen der erlaubten Treffer statt einer Liste vom Stand eines Commits, und
jeder Slice der Welle misst vor seiner Closure, ob seine Artefakte in eine solche Menge fallen.
