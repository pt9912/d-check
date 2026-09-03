# Ein Gate-Aufruf hinter einer Pipe meldet den Exit des letzten Pipe-Glieds

**Sub-Area:** `*`

`make gates 2>&1 | tail` liefert den Exit von `tail` — die `&&`-Kette läuft
weiter, ein roter Lauf wird als grün behandelt, und die Botschaft behauptet
Gates, die nie bestanden wurden. Der Stop-Hook fängt nur das Sitzungsende,
nicht den Commit dazwischen.
