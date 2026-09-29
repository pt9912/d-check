# Ein history-lesendes Modul meldet stilles Grün, wenn die angeforderte Range aufgelöst, aber leer ist

**Sub-Area:** `*`

Der Prüfbereich einer Range-Prüfung kann leer sein, ohne dass die Range
ungültig wäre — typisch im shallow-Klon, in dem `HEAD..HEAD` null Commits
zählt. Meldet das Modul dafür „0 Befunde" mit Exit 0, unterscheidet sich das
Ergebnis nicht von „wirklich nichts zu melden" — das Grün ist über einem
Prüfbereich zustande gekommen, der nie geprüft hat. Fail-closed wäre der
laut-Abruch (Exit ≠ 0) mit benannter Meldung oder eine ausdrücklich
deklarierte Ausnahme; das Vorlauf-Wächter-Muster fängt die Lage eine Stufe
vor dem Modul, entbindet es aber nicht von der eigenen Semantik.
